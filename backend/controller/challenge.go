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
	var req model.ChallengeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
		return
	}

	// 获取创建者钱包地址
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 验证时间
	if req.EndTime.Before(req.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间必须晚于开始时间"})
		return
	}

	if req.EndTime.Before(time.Now()) {
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
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
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
		"msg":   "挑战创建成功",
		"data":  challenge,
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

	// 只显示进行中或即将开始的挑战
	query = query.Where("status IN ?", []string{model.ChallengeStatusActive, model.ChallengeStatusDraft})
	query = query.Where("end_time >= ?", time.Now())

	var total int64
	query.Count(&total)

	// 分页
	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&challenges)

	c.JSON(http.StatusOK, gin.H{
		"data": challenges,
		"total": total,
		"page": parseInt(page),
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
		"data": challenge,
		"participant_count": participantCount,
		"participants": participants,
	})
}

// JoinChallenge 参与挑战
func JoinChallenge(c *gin.Context) {
	var req model.JoinChallengeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 检查挑战是否存在
	var challenge model.Challenge
	if err := config.DB.Where("id = ?", req.ChallengeID).First(&challenge).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "挑战不存在"})
		return
	}

	// 检查挑战状态
	if challenge.Status != model.ChallengeStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "挑战未开始或已结束"})
		return
	}

	// 检查时间
	now := time.Now()
	if now.Before(challenge.StartTime) || now.After(challenge.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不在挑战时间范围内"})
		return
	}

	// 检查是否已参与
	var existing model.ChallengeParticipant
	if err := config.DB.Where("challenge_id = ? AND wallet_addr = ?", req.ChallengeID, walletAddr).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已参与此挑战"})
		return
	}

	// 检查行为记录是否存在且已通过
	var behavior model.GreenBehavior
	if err := config.DB.Where("id = ? AND wallet_addr = ? AND status = ?", req.BehaviorID, walletAddr, model.BehaviorStatusApproved).First(&behavior).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "行为记录不存在或未通过审核"})
		return
	}

	// 创建参与者记录
	participant := &model.ChallengeParticipant{
		ID:           uuid.New().String(),
		ChallengeID:  req.ChallengeID,
		WalletAddr:   walletAddr,
		BehaviorID:   req.BehaviorID,
		JoinedAt:     time.Now(),
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
		"msg": "成功参与挑战",
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

	// 发放奖励
	txHash, err := service.MintReward(walletAddr, challenge.BehaviorType)
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
		"msg": "挑战已激活",
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
