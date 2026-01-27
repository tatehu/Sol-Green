package middleware

import (
	"sol-green/config"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		config.Log.WithFields(map[string]interface{}{
			"status":  statusCode,
			"method":  method,
			"path":    path,
			"latency": latency,
			"ip":      c.ClientIP(),
		}).Info("HTTP Request")
	}
}
