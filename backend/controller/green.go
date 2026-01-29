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

// SubmitGreenBehavior 提交环保行为（垃圾分类/植树/低碳出行等）
// SubmitGreenBehavior submits environmental behavior (waste sorting/tree planting/low-carbon travel, etc.)
func SubmitGreenBehavior(c *gin.Context) {
	var req model.GreenBehaviorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "参数错误", // "Invalid parameters"
			"error_en": "Invalid parameters",
		})
		return
	}

	// 1. 获取当前钱包地址（从 JWT/Token 解析）
	// 1. Get current wallet address (from JWT/Token)
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":    "未登录", // "Unauthorized"
			"error_en": "Unauthorized",
		})
		return
	}

	// 2. 防重复提交（Redis 锁：用户+行为类型+时间窗口）
	// 2. Prevent duplicate submission (Redis lock: user + behavior type + time window)
	lockKey := "green:submit:lock:" + walletAddr + ":" + req.BehaviorType
	if config.RedisClient.Exists(c, lockKey).Val() > 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":    "请勿重复提交", // "Please do not submit duplicate"
			"error_en": "Please do not submit duplicate",
		})
		return
	}
	config.RedisClient.Set(c, lockKey, "1", 300) // 5分钟锁 / 5-minute lock

	// 3. AI 反欺诈检测（图片/视频真实性）
	// 3. AI fraud detection (image/video authenticity)
	fraudScore, isFraud := config.AIFraudDetector.Detect(req.MediaURLs)
	threshold := config.GetAIFraudThreshold()
	rewardCfg := config.GetRewardConfig()
	var rewardAmount uint64
	switch req.BehaviorType {
	case "waste_sorting":
		rewardAmount = rewardCfg.WasteSortingReward
	case "tree_planting":
		rewardAmount = rewardCfg.TreePlantingReward
	case "low_carbon_travel":
		rewardAmount = rewardCfg.LowCarbonTravelReward
	default:
		rewardAmount = rewardCfg.WasteSortingReward
	}
	if isFraud || fraudScore >= threshold { // 对齐文档：阈值由环境控制（主网更严格）
		// 记录可疑行为，进入人工审核
		// Record suspicious behavior, enter manual review
		behavior := &model.GreenBehavior{
			ID:           uuid.New().String(),
			WalletAddr:   walletAddr,
			BehaviorType: req.BehaviorType,
			MediaURLs:    req.MediaURLs,
			Status:       model.BehaviorStatusPendingReview,
			FraudScore:   fraudScore,
			SubmitTime:   time.Now(),
			Location:     req.Location,
		}
		config.DB.Create(behavior)
		c.JSON(http.StatusOK, gin.H{
			"msg":           "提交成功，进入人工审核", // "Submission successful, pending manual review"
			"msg_en":        "Submission successful, pending manual review",
			"id":            behavior.ID,
			"status":        behavior.Status,
			"fraud_score":   fraudScore,
			"reward_amount": rewardAmount,
		})
		return
	}

	// 4. 自动认证通过（低欺诈风险）
	// 4. Auto approve (low fraud risk)
	now := time.Now()
	behavior := &model.GreenBehavior{
		ID:           uuid.New().String(),
		WalletAddr:   walletAddr,
		BehaviorType: req.BehaviorType,
		MediaURLs:    req.MediaURLs,
		Status:       model.BehaviorStatusApproved,
		FraudScore:   fraudScore,
		SubmitTime:   now,
		ApproveTime:  &now,
		Location:     req.Location,
	}
	config.DB.Create(behavior)

	// 5. 调用 Solana 合约，发放奖励（SOL/平台代币）
	// 5. Call Solana contract to distribute rewards (SOL/platform tokens)
	txHash, err := service.MintReward(walletAddr, req.BehaviorType)
	if err != nil {
		config.Log.Errorf("奖励发放失败: %v", err) // "Reward distribution failed"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "奖励发放失败", // "Reward distribution failed"
			"error_en": "Reward distribution failed",
			"detail":   err.Error(),
		})
		return
	}

	// 6. 链上存证（行为哈希上链）
	// 6. On-chain proof storage (behavior hash on-chain)
	proofHash := service.GenerateBehaviorHash(behavior)
	_, err = service.RecordOnChain(proofHash, walletAddr)
	if err != nil {
		// 存证失败不影响奖励，记录日志即可
		// Proof storage failure does not affect rewards, just log
		config.Log.Errorf("链上存证失败: %v", err) // "On-chain proof storage failed"
	}

	// 更新行为记录
	behavior.TxHash = txHash
	behavior.ProofHash = proofHash
	config.DB.Save(behavior)

	c.JSON(http.StatusOK, gin.H{
		"msg":           "认证通过，奖励已发放", // "Verification passed, reward distributed"
		"msg_en":        "Verification passed, reward distributed",
		"id":            behavior.ID,
		"status":        behavior.Status,
		"fraud_score":   fraudScore,
		"reward_amount": rewardAmount,
		"tx_hash":       txHash,
		"proof":         proofHash,
	})
}

// GetBehaviorStatus 查询行为认证状态
// GetBehaviorStatus queries behavior verification status
func GetBehaviorStatus(c *gin.Context) {
	id := c.Param("id")
	var behavior model.GreenBehavior
	if err := config.DB.Where("id = ?", id).First(&behavior).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":    "行为记录不存在", // "Behavior record not found"
			"error_en": "Behavior record not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": behavior})
}

// GetUserBehaviors 获取用户的行为记录列表
// GetUserBehaviors gets user's behavior records list
func GetUserBehaviors(c *gin.Context) {
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	status := c.Query("status")              // 可选：过滤状态，如 "approved"
	behaviorType := c.Query("behavior_type") // 可选：过滤行为类型

	var behaviors []model.GreenBehavior
	query := config.DB.Where("wallet_addr = ?", walletAddr)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if behaviorType != "" {
		query = query.Where("behavior_type = ?", behaviorType)
	}

	// GreenBehavior 模型没有 created_at 字段；用 submit_time 作为排序依据
	if err := query.Order("submit_time DESC").Find(&behaviors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  behaviors,
		"total": len(behaviors),
	})
}

// ClaimReward 领取奖励
func ClaimReward(c *gin.Context) {
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// TODO: 实现奖励领取逻辑
	c.JSON(http.StatusOK, gin.H{"msg": "奖励领取成功"})
}

// PartnerVerifyReq 第三方认证请求
type PartnerVerifyReq struct {
	BehaviorID string `json:"behavior_id" binding:"required"`
	PartnerID  string `json:"partner_id" binding:"required"`  // 认证机构 ID
	VerifyCode string `json:"verify_code" binding:"required"` // 认证码
}

// PartnerVerify 第三方机构认证
func PartnerVerify(c *gin.Context) {
	var req PartnerVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 查询行为记录
	var behavior model.GreenBehavior
	if err := config.DB.Where("id = ?", req.BehaviorID).First(&behavior).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "行为记录不存在"})
		return
	}

	// 调用第三方认证接口
	verified, err := service.VerifyWithPartner(req.PartnerID, req.VerifyCode, &behavior)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "认证失败", "detail": err.Error()})
		return
	}

	if verified {
		behavior.Status = model.BehaviorStatusApproved
		now := time.Now()
		behavior.ApproveTime = &now
		behavior.PartnerVerify = req.PartnerID
		config.DB.Save(&behavior)

		// 发放奖励
		txHash, _ := service.MintReward(behavior.WalletAddr, behavior.BehaviorType)
		behavior.TxHash = txHash
		config.DB.Save(&behavior)

		c.JSON(http.StatusOK, gin.H{
			"msg":     "第三方认证通过，奖励已发放",
			"tx_hash": txHash,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "第三方认证未通过"})
	}
}
