package controller

import (
	"fmt"
	"net/http"
	"sol-green/config"
	"sol-green/model"
	"sol-green/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateMarketingActivity 创建营销活动（仅管理员）
func CreateMarketingActivity(c *gin.Context) {
	var req struct {
		Title           string    `json:"title" binding:"required"`
		Description     string    `json:"description" binding:"required"`
		ActivityType    string    `json:"activity_type" binding:"required"`
		StartTime       time.Time `json:"start_time" binding:"required"`
		EndTime         time.Time `json:"end_time" binding:"required"`
		RewardAmount    uint64    `json:"reward_amount" binding:"required,min=1"`
		MaxParticipants int       `json:"max_participants"`
		Rules           string    `json:"rules"`
		ImageURL        string    `json:"image_url"`
		Config          string    `json:"config"`
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

	// TODO: 检查管理员权限
	// if !isAdmin(walletAddr) {
	//     c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
	//     return
	// }

	activity := &model.MarketingActivity{
		ID:                uuid.New().String(),
		Title:             req.Title,
		Description:       req.Description,
		ActivityType:      req.ActivityType,
		Status:            model.MarketingStatusScheduled,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		RewardAmount:      req.RewardAmount,
		MaxParticipants:   req.MaxParticipants,
		CurrentParticipants: 0,
		Rules:             req.Rules,
		ImageURL:          req.ImageURL,
		Config:            req.Config,
		CreatedBy:         walletAddr,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := config.DB.Create(activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建活动失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "活动创建成功",
		"data": activity,
	})
}

// GetMarketingActivities 获取营销活动列表
func GetMarketingActivities(c *gin.Context) {
	status := c.Query("status")
	activityType := c.Query("activity_type")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	var activities []model.MarketingActivity
	query := config.DB.Model(&model.MarketingActivity{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if activityType != "" {
		query = query.Where("activity_type = ?", activityType)
	}

	// 只显示进行中或即将开始的活动
	now := time.Now()
	query = query.Where("(status = ? OR status = ?) AND end_time >= ?", 
		model.MarketingStatusActive, model.MarketingStatusScheduled, now)

	var total int64
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Offset(offset).Limit(parseInt(pageSize)).Order("start_time ASC").Find(&activities)

	c.JSON(http.StatusOK, gin.H{
		"data": activities,
		"total": total,
		"page": parseInt(page),
		"page_size": parseInt(pageSize),
	})
}

// JoinMarketingActivity 参与营销活动
func JoinMarketingActivity(c *gin.Context) {
	activityID := c.Param("id")
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 检查活动
	var activity model.MarketingActivity
	if err := config.DB.Where("id = ?", activityID).First(&activity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
		return
	}

	if activity.Status != model.MarketingStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "活动未开始或已结束"})
		return
	}

	now := time.Now()
	if now.Before(activity.StartTime) || now.After(activity.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不在活动时间范围内"})
		return
	}

	// 检查是否已参与
	var existing model.MarketingParticipant
	if err := config.DB.Where("activity_id = ? AND wallet_addr = ?", activityID, walletAddr).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您已参与此活动"})
		return
	}

	// 检查参与人数限制
	if activity.MaxParticipants > 0 && activity.CurrentParticipants >= activity.MaxParticipants {
		c.JSON(http.StatusBadRequest, gin.H{"error": "活动参与人数已满"})
		return
	}

	// 根据活动类型处理
	var participant *model.MarketingParticipant
	var err error

	switch activity.ActivityType {
	case model.MarketingTypeSignIn:
		participant, err = handleSignIn(activityID, walletAddr)
	case model.MarketingTypeInvite:
		participant, err = handleInvite(activityID, walletAddr, c)
	case model.MarketingTypeDailyTask:
		participant, err = handleDailyTask(activityID, walletAddr, c)
	default:
		participant = &model.MarketingParticipant{
			ID:         uuid.New().String(),
			ActivityID: activityID,
			WalletAddr: walletAddr,
			JoinedAt:   time.Now(),
		}
		err = config.DB.Create(participant).Error
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "参与活动失败", "detail": err.Error()})
		return
	}

	// 更新活动参与人数
	config.DB.Model(&activity).Update("current_participants", activity.CurrentParticipants+1)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "成功参与活动",
		"data": participant,
	})
}

// handleSignIn 处理签到活动
func handleSignIn(activityID, walletAddr string) (*model.MarketingParticipant, error) {
	today := time.Now().Format("2006-01-02")
	
	// 检查今日是否已签到
	var todayRecord model.SignInRecord
	if err := config.DB.Where("activity_id = ? AND wallet_addr = ? AND sign_in_date = ?", 
		activityID, walletAddr, today).First(&todayRecord).Error; err == nil {
		return nil, fmt.Errorf("今日已签到")
	}

	// 获取连续签到天数
	var lastRecord model.SignInRecord
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	consecutiveDays := 1
	
	if err := config.DB.Where("activity_id = ? AND wallet_addr = ? AND sign_in_date = ?", 
		activityID, walletAddr, yesterday).First(&lastRecord).Error; err == nil {
		consecutiveDays = lastRecord.ConsecutiveDays + 1
	}

	// 创建签到记录
	signInRecord := &model.SignInRecord{
		ID:              uuid.New().String(),
		WalletAddr:      walletAddr,
		ActivityID:      activityID,
		SignInDate:      time.Now(),
		ConsecutiveDays: consecutiveDays,
	}

	if err := config.DB.Create(signInRecord).Error; err != nil {
		return nil, err
	}

	// 创建参与者记录
	participant := &model.MarketingParticipant{
		ID:         uuid.New().String(),
		ActivityID: activityID,
		WalletAddr: walletAddr,
		JoinedAt:   time.Now(),
		ExtraData:  fmt.Sprintf(`{"consecutive_days": %d}`, consecutiveDays),
	}

	return participant, config.DB.Create(participant).Error
}

// handleInvite 处理邀请活动
func handleInvite(activityID, walletAddr string, c *gin.Context) (*model.MarketingParticipant, error) {
	var req struct {
		InviteeAddr string `json:"invitee_addr" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	// 检查是否已邀请过此人
	var existing model.InviteRecord
	if err := config.DB.Where("inviter_addr = ? AND invitee_addr = ? AND activity_id = ?", 
		walletAddr, req.InviteeAddr, activityID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("已邀请过此用户")
	}

	// 创建邀请记录
	inviteRecord := &model.InviteRecord{
		ID:          uuid.New().String(),
		InviterAddr: walletAddr,
		InviteeAddr: req.InviteeAddr,
		ActivityID:  activityID,
		InvitedAt:   time.Now(),
	}

	if err := config.DB.Create(inviteRecord).Error; err != nil {
		return nil, err
	}

	// 创建参与者记录
	participant := &model.MarketingParticipant{
		ID:         uuid.New().String(),
		ActivityID: activityID,
		WalletAddr: walletAddr,
		JoinedAt:   time.Now(),
	}

	return participant, config.DB.Create(participant).Error
}

// handleDailyTask 处理每日任务
func handleDailyTask(activityID, walletAddr string, c *gin.Context) (*model.MarketingParticipant, error) {
	var req struct {
		TaskID string `json:"task_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	// TODO: 验证任务完成情况
	// 这里简化处理，实际应检查任务是否真正完成

	participant := &model.MarketingParticipant{
		ID:         uuid.New().String(),
		ActivityID: activityID,
		WalletAddr: walletAddr,
		JoinedAt:   time.Now(),
		ExtraData:  fmt.Sprintf(`{"task_id": "%s"}`, req.TaskID),
	}

	return participant, config.DB.Create(participant).Error
}

// ClaimMarketingReward 领取营销活动奖励
func ClaimMarketingReward(c *gin.Context) {
	activityID := c.Param("id")
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 检查参与者记录
	var participant model.MarketingParticipant
	if err := config.DB.Where("activity_id = ? AND wallet_addr = ?", activityID, walletAddr).First(&participant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "您未参与此活动"})
		return
	}

	if participant.RewardClaimed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "奖励已领取"})
		return
	}

	// 获取活动信息
	var activity model.MarketingActivity
	if err := config.DB.Where("id = ?", activityID).First(&activity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
		return
	}

	// 发放奖励
	txHash, err := service.MintRewardMultiChain(walletAddr, "marketing_reward", activity.RewardAmount)
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
