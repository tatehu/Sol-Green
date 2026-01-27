package config

import (
	"fmt"
	"os"
)

// Environment 环境类型
type Environment string

const (
	EnvironmentDev    Environment = "dev"     // 开发环境
	EnvironmentTest   Environment = "test"    // 测试环境（Solana Devnet）
	EnvironmentMainnet Environment = "mainnet" // 主网环境（Solana Mainnet）
)

var (
	CurrentEnv Environment
)

// InitEnvironment 初始化环境配置
func InitEnvironment() {
	envStr := GetEnv("ENVIRONMENT", "test")
	CurrentEnv = Environment(envStr)

	// 验证环境值
	switch CurrentEnv {
	case EnvironmentDev, EnvironmentTest, EnvironmentMainnet:
		Log.Infof("当前环境: %s", CurrentEnv)
	default:
		Log.Warnf("未知环境: %s，使用默认测试环境", envStr)
		CurrentEnv = EnvironmentTest
	}
}

// GetEnvironment 获取当前环境
func GetEnvironment() Environment {
	if CurrentEnv == "" {
		InitEnvironment()
	}
	return CurrentEnv
}

// IsDev 是否为开发环境
func IsDev() bool {
	return GetEnvironment() == EnvironmentDev
}

// IsTest 是否为测试环境
func IsTest() bool {
	return GetEnvironment() == EnvironmentTest
}

// IsMainnet 是否为主网环境
func IsMainnet() bool {
	return GetEnvironment() == EnvironmentMainnet
}

// GetSolanaRPCURL 根据环境获取 Solana RPC URL
func GetSolanaRPCURL() string {
	env := GetEnvironment()
	
	// 优先使用环境变量配置
	if url := os.Getenv("SOLANA_RPC_URL"); url != "" {
		return url
	}

	// 根据环境返回默认 RPC
	switch env {
	case EnvironmentDev:
		return GetEnv("SOLANA_RPC_URL", "http://localhost:8899") // 本地验证器
	case EnvironmentTest:
		return GetEnv("SOLANA_RPC_URL", "https://api.devnet.solana.com") // Devnet
	case EnvironmentMainnet:
		return GetEnv("SOLANA_RPC_URL", "https://api.mainnet-beta.solana.com") // Mainnet
	default:
		return "https://api.devnet.solana.com"
	}
}

// GetSolanaWSURL 根据环境获取 Solana WebSocket URL
func GetSolanaWSURL() string {
	env := GetEnvironment()
	
	if url := os.Getenv("SOLANA_WS_URL"); url != "" {
		return url
	}

	switch env {
	case EnvironmentDev:
		return GetEnv("SOLANA_WS_URL", "ws://localhost:8900")
	case EnvironmentTest:
		return GetEnv("SOLANA_WS_URL", "wss://api.devnet.solana.com")
	case EnvironmentMainnet:
		return GetEnv("SOLANA_WS_URL", "wss://api.mainnet-beta.solana.com")
	default:
		return "wss://api.devnet.solana.com"
	}
}

// GetSolanaCluster 获取 Solana 集群名称
func GetSolanaCluster() string {
	env := GetEnvironment()
	switch env {
	case EnvironmentDev:
		return "localnet"
	case EnvironmentTest:
		return "devnet"
	case EnvironmentMainnet:
		return "mainnet-beta"
	default:
		return "devnet"
	}
}

// GetRewardConfig 根据环境获取奖励配置
func GetRewardConfig() RewardConfig {
	env := GetEnvironment()
	
	// 从环境变量读取，如果没有则使用默认值
	config := RewardConfig{
		WasteSortingReward:   uint64(GetEnvInt("REWARD_WASTE_SORTING", 1000)),
		TreePlantingReward:  uint64(GetEnvInt("REWARD_TREE_PLANTING", 5000)),
		LowCarbonTravelReward: uint64(GetEnvInt("REWARD_LOW_CARBON_TRAVEL", 2000)),
		ChallengeBonusRate:  float64(GetEnvFloat("CHALLENGE_BONUS_RATE", 1.5)), // 挑战奖励倍率
	}

	// 主网环境可以有不同的配置
	if env == EnvironmentMainnet {
		// 主网可以使用更高的奖励
		if config.WasteSortingReward == 1000 {
			config.WasteSortingReward = 2000
		}
		if config.TreePlantingReward == 5000 {
			config.TreePlantingReward = 10000
		}
		if config.LowCarbonTravelReward == 2000 {
			config.LowCarbonTravelReward = 4000
		}
	}

	return config
}

// RewardConfig 奖励配置
type RewardConfig struct {
	WasteSortingReward    uint64
	TreePlantingReward    uint64
	LowCarbonTravelReward uint64
	ChallengeBonusRate    float64
}

// GetEnvInt 获取整数环境变量
func GetEnvInt(key string, defaultValue int) int {
	value := GetEnv(key, fmt.Sprintf("%d", defaultValue))
	var result int
	fmt.Sscanf(value, "%d", &result)
	if result == 0 && value != "0" {
		return defaultValue
	}
	return result
}

// GetEnvFloat 获取浮点数环境变量
func GetEnvFloat(key string, defaultValue float64) float64 {
	value := GetEnv(key, fmt.Sprintf("%f", defaultValue))
	var result float64
	fmt.Sscanf(value, "%f", &result)
	if result == 0 && value != "0" {
		return defaultValue
	}
	return result
}

// GetAIFraudThreshold 获取 AI 反欺诈阈值
func GetAIFraudThreshold() float64 {
	env := GetEnvironment()
	
	// 主网使用更严格的阈值
	switch env {
	case EnvironmentMainnet:
		return GetEnvFloat("AI_FRAUD_THRESHOLD", 0.5) // 主网更严格
	case EnvironmentTest:
		return GetEnvFloat("AI_FRAUD_THRESHOLD", 0.7) // 测试环境
	default:
		return GetEnvFloat("AI_FRAUD_THRESHOLD", 0.7) // 开发环境
	}
}

// GetRateLimitConfig 获取限流配置
func GetRateLimitConfig() RateLimitConfig {
	env := GetEnvironment()
	
	config := RateLimitConfig{
		RequestsPerMinute: GetEnvInt("RATE_LIMIT_REQUESTS", 60),
		BurstSize:         GetEnvInt("RATE_LIMIT_BURST", 10),
	}

	// 主网可以更宽松
	if env == EnvironmentMainnet {
		if config.RequestsPerMinute == 60 {
			config.RequestsPerMinute = 120
		}
	}

	return config
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
}
