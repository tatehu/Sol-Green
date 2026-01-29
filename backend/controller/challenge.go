package controller

import (
	"net/http"
	"sol-green/config"
	"sol-green/model"
	"sol-green/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateChallenge 创建挑战活动
func CreateChallenge(c *gin.Context) {
	// 使用 string 接收时间字段，避免 Gin 直接按 RFC3339 解析导致 datetime-local（如 2026-01-25T02:07）报错
	var req struct {
		Title         string `json:"title" binding:"required"`
		Description   string `json:"description" binding:"required"`
		ChallengeType string `json:"challenge_type" binding:"required,oneof=individual team community"`
		BehaviorType  string `json:"behavior_type" binding:"required,oneof=waste_sorting tree_planting low_carbon_travel"`
		RewardAmount  uint64 `json:"reward_amount" binding:"required,min=1"`
		TargetCount   int    `json:"target_count" binding:"required,min=1"`
		StartTime     string `json:"start_time" binding:"required"`
		EndTime       string `json:"end_time" binding:"required"`
		ImageURL      string `json:"image_url,omitempty"`
		Rules         string `json:"rules,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	startTime, err := parseFlexibleTime(req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": "start_time: " + err.Error()})
		return
	}
	endTime, err := parseFlexibleTime(req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": "end_time: " + err.Error()})
		return
	}

	// 获取创建者钱包地址
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 验证时间
	if endTime.Before(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间必须晚于开始时间"})
		return
	}

	if endTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于当前时间"})
		return
	}

	// 创建挑战
	challenge := &model.Challenge{
		ID:            uuid.New().String(),
		Title:         req.Title,
		Description:   req.Description,
		ChallengeType: req.ChallengeType,
		BehaviorType:  req.BehaviorType,
		RewardAmount:  req.RewardAmount,
		TargetCount:   req.TargetCount,
		CurrentCount:  0,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        model.ChallengeStatusDraft,
		CreatorWallet: walletAddr,
		ImageURL:      req.ImageURL,
		Rules:         req.Rules,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := config.DB.Create(challenge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建挑战失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "挑战创建成功",
		"data": challenge,
	})
}

// GetChallenges 获取挑战列表
func GetChallenges(c *gin.Context) {
	status := c.Query("status")
	challengeType := c.Query("challenge_type")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	var challenges []model.Challenge
	query := config.DB.Model(&model.Challenge{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if challengeType != "" {
		query = query.Where("challenge_type = ?", challengeType)
	}

	// 默认只显示进行中或草稿（即将开始）；如果显式传了 status，则尊重调用方
	if status == "" {
		query = query.Where("status IN ?", []string{model.ChallengeStatusActive, model.ChallengeStatusDraft})
		// 默认过滤掉已结束的
		query = query.Where("end_time >= ?", time.Now())
	}

	var total int64
	query.Count(&total)

	// 分页
	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&challenges)

	c.JSON(http.StatusOK, gin.H{
		"data":      challenges,
		"total":     total,
		"page":      parseInt(page),
		"page_size": parseInt(pageSize),
	})
}

// GetChallengeDetail 获取挑战详情
func GetChallengeDetail(c *gin.Context) {
	id := c.Param("id")
	var challenge model.Challenge
	if err := config.DB.Where("id = ?", id).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "挑战不存在"})
		return
	}

	// 获取参与者数量
	var participantCount int64
	config.DB.Model(&model.ChallengeParticipant{}).Where("challenge_id = ?", id).Count(&participantCount)

	// 获取参与者列表（前10名）
	var participants []model.ChallengeParticipant
	config.DB.Where("challenge_id = ?", id).Order("joined_at DESC").Limit(10).Find(&participants)

	c.JSON(http.StatusOK, gin.H{
		"data":              challenge,
		"participant_count": participantCount,
		"participants":      participants,
	})
}

// JoinChallenge 参与挑战
func JoinChallenge(c *gin.Context) {
	challengeID := c.Param("id")
	var req struct {
		BehaviorID string `json:"behavior_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 检查挑战是否存在
	var challenge model.Challenge
	if err := config.DB.Where("id = ?", challengeID).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "挑战不存在"})
		return
	}

	// 检查时间
	now := time.Now()
	if now.Before(challenge.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "挑战尚未开始"})
		return
	}
	if now.After(challenge.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "挑战已结束"})
		return
	}

	// 检查挑战状态：允许active或draft（但时间已到）的挑战参与
	if challenge.Status == model.ChallengeStatusCompleted || challenge.Status == model.ChallengeStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "挑战已结束或已取消"})
		return
	}
	// 如果状态是draft但时间已到，自动激活
	if challenge.Status == model.ChallengeStatusDraft && now.After(challenge.StartTime) && now.Before(challenge.EndTime) {
		config.DB.Model(&challenge).Update("status", model.ChallengeStatusActive)
		challenge.Status = model.ChallengeStatusActive
	}

	// 检查是否已参与
	var existing model.ChallengeParticipant
	if err := config.DB.Where("challenge_id = ? AND wallet_addr = ?", challengeID, walletAddr).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已参与此挑战"})
		return
	}

	// 检查行为记录是否存在且已通过
	var behavior model.GreenBehavior
	if err := config.DB.Where("id = ? AND wallet_addr = ? AND status = ?", req.BehaviorID, walletAddr, model.BehaviorStatusApproved).First(&behavior).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "行为记录不存在或未通过审核"})
		return
	}

	// 检查行为类型是否匹配挑战要求
	if behavior.BehaviorType != challenge.BehaviorType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "行为类型不匹配挑战要求"})
		return
	}

	// 创建参与者记录
	participant := &model.ChallengeParticipant{
		ID:            uuid.New().String(),
		ChallengeID:   challengeID,
		WalletAddr:    walletAddr,
		BehaviorID:    req.BehaviorID,
		JoinedAt:      time.Now(),
		RewardClaimed: false,
	}

	if err := config.DB.Create(participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "参与挑战失败"})
		return
	}

	// 更新挑战参与人数
	config.DB.Model(&challenge).Update("current_count", challenge.CurrentCount+1)

	// 检查是否达到目标人数
	if challenge.CurrentCount+1 >= challenge.TargetCount {
		config.DB.Model(&challenge).Update("status", model.ChallengeStatusCompleted)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "成功参与挑战",
		"data": participant,
	})
}

// ClaimChallengeReward 领取挑战奖励
func ClaimChallengeReward(c *gin.Context) {
	challengeID := c.Param("id")
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 检查参与者记录
	var participant model.ChallengeParticipant
	if err := config.DB.Where("challenge_id = ? AND wallet_addr = ?", challengeID, walletAddr).First(&participant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "您未参与此挑战"})
		return
	}

	if participant.RewardClaimed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "奖励已领取"})
		return
	}

	// 获取挑战信息
	var challenge model.Challenge
	if err := config.DB.Where("id = ?", challengeID).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "挑战不存在"})
		return
	}

	// 检查挑战是否完成
	if challenge.Status != model.ChallengeStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "挑战尚未完成"})
		return
	}

	// 发放奖励（对齐文档：baseReward + baseReward*challengeBonusRate）
	rewardCfg := config.GetRewardConfig()
	var baseReward uint64
	switch challenge.BehaviorType {
	case "waste_sorting":
		baseReward = rewardCfg.WasteSortingReward
	case "tree_planting":
		baseReward = rewardCfg.TreePlantingReward
	case "low_carbon_travel":
		baseReward = rewardCfg.LowCarbonTravelReward
	default:
		baseReward = rewardCfg.WasteSortingReward
	}

	bonus := uint64(float64(baseReward) * rewardCfg.ChallengeBonusRate)
	totalReward := baseReward + bonus

	txHash, err := service.MintRewardMultiChain(walletAddr, "challenge_reward", totalReward)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "奖励发放失败", "detail": err.Error()})
		return
	}

	// 更新参与者记录
	participant.RewardClaimed = true
	participant.TxHash = txHash
	config.DB.Save(&participant)

	c.JSON(http.StatusOK, gin.H{
		"msg":     "奖励领取成功",
		"tx_hash": txHash,
	})
}

// ActivateChallenge 激活挑战（创建者）
func ActivateChallenge(c *gin.Context) {
	id := c.Param("id")
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var challenge model.Challenge
	if err := config.DB.Where("id = ?", id).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "挑战不存在"})
		return
	}

	// 检查权限
	if challenge.CreatorWallet != walletAddr {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有创建者可以激活挑战"})
		return
	}

	if challenge.Status != model.ChallengeStatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能激活草稿状态的挑战"})
		return
	}

	// 激活挑战
	challenge.Status = model.ChallengeStatusActive
	config.DB.Save(&challenge)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "挑战已激活",
		"data": challenge,
	})
}

func parseInt(s string) int {
	var result int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 1
		}
	}
	if result < 1 {
		return 1
	}
	return result
}
