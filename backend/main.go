package main

import (
	"sol-green/config"
	"sol-green/controller"
	"sol-green/middleware"
	"sol-green/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置、DB、Redis、AI 反欺诈
	// Initialize configuration, DB, Redis, AI fraud detection
	config.InitConfig()
	config.InitDB()
	config.InitRedis()
	config.InitAIFraudDetector()

	// 启动区块链监听服务（后台运行）
	// Start blockchain listener service (background)
	go func() {
		listener, err := service.NewBlockchainListener()
		if err != nil {
			config.Log.Warnf("区块链监听服务启动失败: %v", err) // "Blockchain listener service startup failed"
			return
		}
		if err := listener.StartListening(); err != nil {
			config.Log.Warnf("区块链监听服务启动失败: %v", err) // "Blockchain listener service startup failed"
		}
	}()

	r := gin.Default()
	// 全局中间件：日志、跨域、请求限流
	// Global middleware: logger, CORS, rate limiting
	r.Use(middleware.Logger(), middleware.Cors(), middleware.RateLimit())

	// 健康检查（无需认证）
	// Health check (no authentication required)
	r.GET("/api/v1/health", controller.Health)

	// API 路由
	// API routes
	api := r.Group("/api/v1")
	{
		// 用户认证（Solana 钱包登录）
		// User authentication (Solana wallet login)
		api.POST("/auth/wallet", controller.WalletAuth)
		
		// 环保行为相关
		// Environmental behavior related
		api.POST("/green/behavior/submit", middleware.Auth(), controller.SubmitGreenBehavior)
		api.GET("/green/behavior/:id", middleware.Auth(), controller.GetBehaviorStatus)
		api.POST("/green/reward/claim", middleware.Auth(), controller.ClaimReward)
		api.POST("/green/partner/verify", middleware.Auth(), controller.PartnerVerify)
		
		// 挑战活动相关
		// Challenge activity related
		api.POST("/challenges", middleware.Auth(), controller.CreateChallenge)
		api.GET("/challenges", controller.GetChallenges)
		api.GET("/challenges/:id", controller.GetChallengeDetail)
		api.POST("/challenges/:id/join", middleware.Auth(), controller.JoinChallenge)
		api.POST("/challenges/:id/activate", middleware.Auth(), controller.ActivateChallenge)
		api.POST("/challenges/:id/claim", middleware.Auth(), controller.ClaimChallengeReward)
		
		// 营销活动相关
		// Marketing activity related
		api.POST("/marketing/activities", middleware.Auth(), controller.CreateMarketingActivity)
		api.GET("/marketing/activities", controller.GetMarketingActivities)
		api.POST("/marketing/activities/:id/join", middleware.Auth(), controller.JoinMarketingActivity)
		api.POST("/marketing/activities/:id/claim", middleware.Auth(), controller.ClaimMarketingReward)
		
		// 运营管理相关
		// Operations management related
		api.GET("/admin/stats", middleware.Auth(), controller.GetAdminStats)
		api.GET("/admin/analytics", middleware.Auth(), controller.GetAnalytics)
		api.POST("/admin/activities/:id/activate", middleware.Auth(), controller.ActivateMarketingActivity)
		
		// 合约管理相关
		// Contract management related
		api.GET("/admin/contract/info", middleware.Auth(), controller.GetContractInfo)
		api.GET("/admin/contract/upgrade/check", middleware.Auth(), controller.CheckUpgrade)
		api.POST("/admin/contract/upgrade/prepare", middleware.Auth(), controller.PrepareUpgrade)
		api.POST("/admin/contract/upgrade/execute", middleware.Auth(), controller.ExecuteUpgrade)
		api.GET("/admin/contract/upgrade/history", middleware.Auth(), controller.GetUpgradeHistory)
	}

	// 启动服务
	// Start service
	r.Run(":8080")
}
