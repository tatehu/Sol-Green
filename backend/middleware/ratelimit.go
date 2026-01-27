package middleware

import (
	"sol-green/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit 请求限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 Redis 实现限流
		// 每个 IP 每分钟最多 60 个请求
		ip := c.ClientIP()
		key := "rate_limit:" + ip
		ctx := c.Request.Context()

		// 检查当前计数
		count, err := config.RedisClient.Get(ctx, key).Int()
		if err != nil && err.Error() != "redis: nil" {
			// Redis 不可用时，跳过限流
			c.Next()
			return
		}

		if count >= 60 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}

		// 增加计数
		if count == 0 {
			config.RedisClient.Set(ctx, key, 1, time.Minute)
		} else {
			config.RedisClient.Incr(ctx, key)
		}

		c.Next()
	}
}
