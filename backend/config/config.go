package config

import (
	"context"
	"fmt"
	"os"
	"sol-green/model"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB              *gorm.DB
	RedisClient     *redis.Client
	AIFraudDetector *FraudDetector
	Log             *logrus.Logger
)

// InitConfig 初始化配置
func InitConfig() {
	// 加载 .env 文件（覆盖常见几种运行目录：项目根目录、backend、backend/config）
	_ = godotenv.Load(
		".env",         // 当前目录
		"../.env",      // 上一级目录
		"../../.env",   // 再上一级（例如从 backend/config 运行）
		"config/.env",  // 项目根目录下的 config/.env（如果存在）
		"backend/.env", // 项目根目录下的 backend/.env（如果存在）
	)

	// 初始化日志
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)

	// 初始化环境配置
	InitEnvironment()
}

// InitDB 初始化数据库
func InitDB() {
	var err error
	var dialector gorm.Dialector

	dbType := GetEnv("DB_TYPE", "sqlite")
	dbURL := GetEnv("DATABASE_URL", "")

	if dbType == "postgres" && dbURL != "" {
		dialector = postgres.Open(dbURL)
	} else {
		// 默认使用 SQLite（开发环境）
		dialector = sqlite.Open("sol-green.db")
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	// 自动迁移
	DB.AutoMigrate(
		&model.GreenBehavior{},
		&model.User{},
		&model.RewardRecord{},
		&model.Challenge{},
		&model.ChallengeParticipant{},
		&model.MarketingActivity{},
		&model.MarketingParticipant{},
		&model.SignInRecord{},
		&model.InviteRecord{},
	)

	Log.Info("数据库初始化成功")
}

// InitRedis 初始化 Redis
func InitRedis() {
	redisURL := GetEnv("REDIS_URL", "localhost:6379")
	redisPassword := GetEnv("REDIS_PASSWORD", "")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		Log.Warnf("Redis 连接失败，将使用内存缓存: %v", err)
		// 开发环境可以继续运行，生产环境应该 panic
	} else {
		Log.Info("Redis 初始化成功")
	}
}

// InitAIFraudDetector 初始化 AI 反欺诈检测器
func InitAIFraudDetector() {
	AIFraudDetector = NewFraudDetector()
	Log.Info("AI 反欺诈检测器初始化成功")
}

// GetEnv 获取环境变量
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		// 当环境变量未设置时，返回调用方提供的默认值
		return defaultValue
	}
	return value
}

// Now 返回当前时间
func Now() time.Time {
	return time.Now()
}

// SolanaRPCURL Solana RPC 地址（已废弃，使用 GetSolanaRPCURL）
func SolanaRPCURL() string {
	return GetSolanaRPCURL()
}

// SolanaAdminPrivKey 管理员私钥
func SolanaAdminPrivKey() string {
	return GetEnv("SOLANA_ADMIN_PRIVKEY", "")
}

// SolanaFeePayer 费用支付者地址
func SolanaFeePayer() string {
	return GetEnv("SOLANA_FEE_PAYER", "")
}

// SolGreenTokenMint 代币 mint 地址
func SolGreenTokenMint() string {
	return GetEnv("SOLGREEN_TOKEN_MINT", "")
}

// SolGreenTokenAdminATA 管理员代币账户
func SolGreenTokenAdminATA() string {
	return GetEnv("SOLGREEN_TOKEN_ADMIN_ATA", "")
}

// SolanaProofContract 存证合约地址
func SolanaProofContract() string {
	return GetEnv("SOLANA_PROOF_CONTRACT", "")
}

// SolanaProofAccount 存证账户地址
func SolanaProofAccount() string {
	return GetEnv("SOLANA_PROOF_ACCOUNT", "")
}
