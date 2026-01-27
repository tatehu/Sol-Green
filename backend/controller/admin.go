package controller

import (
	"net/http"
	"sol-green/config"
	"sol-green/model"
	"sol-green/service"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAdminStats 获取运营统计数据
func GetAdminStats(c *gin.Context) {
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

	// 获取时间范围（默认最近7天）
	days := c.DefaultQuery("days", "7")
	daysInt := parseInt(days)
	if daysInt < 1 || daysInt > 365 {
		daysInt = 7
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -daysInt)

	// 基础统计
	var totalUsers, totalBehaviors, activeUsers int64
	var totalRewards uint64

	config.DB.Model(&model.User{}).Count(&totalUsers)
	config.DB.Model(&model.GreenBehavior{}).
		Where("submit_time >= ?", startDate).
		Count(&totalBehaviors)
	
	config.DB.Model(&model.GreenBehavior{}).
		Where("submit_time >= ?", startDate).
		Distinct("wallet_addr").
		Count(&activeUsers)

	var rewards struct {
		Total uint64
	}
	config.DB.Model(&model.RewardRecord{}).
		Select("COALESCE(SUM(amount), 0) as total").
		Where("created_at >= ?", startDate).
		Scan(&rewards)
	totalRewards = rewards.Total

	// 获取链上统计
	chainStats, _ := service.GetChainStats()

	c.JSON(http.StatusOK, gin.H{
		"period": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
			"days":       daysInt,
		},
		"users": gin.H{
			"total":  totalUsers,
			"active": activeUsers,
		},
		"behaviors": gin.H{
			"total": totalBehaviors,
		},
		"rewards": gin.H{
			"total": totalRewards,
		},
		"chain_stats": chainStats,
	})
}

// GetAnalytics 获取详细数据分析
func GetAnalytics(c *gin.Context) {
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// TODO: 检查管理员权限

	// 获取时间范围
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始日期格式错误"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束日期格式错误"})
		return
	}

	// 获取分析数据
	analytics, err := service.GetAnalytics(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分析数据失败"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// ActivateMarketingActivity 激活营销活动
func ActivateMarketingActivity(c *gin.Context) {
	activityID := c.Param("id")
	walletAddr := c.GetString("wallet_addr")
	if walletAddr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// TODO: 检查管理员权限

	var activity model.MarketingActivity
	if err := config.DB.Where("id = ?", activityID).First(&activity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
		return
	}

	if activity.Status != model.MarketingStatusScheduled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只能激活已排期的活动"})
		return
	}

	// 检查时间
	now := time.Now()
	if now.Before(activity.StartTime) {
		activity.Status = model.MarketingStatusScheduled
	} else if now.After(activity.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "活动已过期"})
		return
	} else {
		activity.Status = model.MarketingStatusActive
	}

	config.DB.Save(&activity)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "活动已激活",
		"data": activity,
	})
}
