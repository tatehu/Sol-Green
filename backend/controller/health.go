package controller

import (
	"net/http"
	"sol-green/config"

	"github.com/gin-gonic/gin"
)

// Health 健康检查接口
// Health health check endpoint
func Health(c *gin.Context) {
	// 检查数据库连接
	// Check database connection
	sqlDB, err := config.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "数据库连接失败", // "Database connection failed"
			"error_en": "Database connection failed",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "数据库 ping 失败", // "Database ping failed"
			"error_en": "Database ping failed",
		})
		return
	}

	// 检查 Redis 连接（可选）
	// Check Redis connection (optional)
	ctx := c.Request.Context()
	if config.RedisClient != nil {
		if err := config.RedisClient.Ping(ctx).Err(); err != nil {
			// Redis 不可用不影响健康检查，只记录警告
			// Redis unavailability does not affect health check, just log warning
			config.Log.Warnf("Redis 连接失败: %v", err) // "Redis connection failed"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "sol-green",
		"version": "3.0.0",
	})
}
