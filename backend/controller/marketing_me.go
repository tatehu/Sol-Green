package controller

import (
	"net/http"
	"sol-green/config"
	"sol-green/model"

	"github.com/gin-gonic/gin"
)

// GetMyMarketingActivityStatus 获取当前用户对某营销活动的参与/领取状态
// GET /api/v1/marketing/activities/:id/me
func GetMyMarketingActivityStatus(c *gin.Context) {
	activityID := c.Param("id")
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var participant model.MarketingParticipant
	if err := config.DB.Where("activity_id = ? AND wallet_addr = ?", activityID, walletAddr).First(&participant).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"joined":         false,
			"reward_claimed": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"joined":         true,
		"reward_claimed": participant.RewardClaimed,
		"tx_hash":        participant.TxHash,
		"joined_at":      participant.JoinedAt,
	})
}

