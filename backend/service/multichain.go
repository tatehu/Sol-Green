package service

import (
	"context"
	"fmt"
	"sol-green/config"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// ChainType 链类型
type ChainType string

const (
	ChainSolana   ChainType = "solana"
	ChainEthereum ChainType = "ethereum"
	ChainPolygon  ChainType = "polygon"
	ChainBSC      ChainType = "bsc"
)

// GetDefaultChain 获取默认链（根据环境）
func GetDefaultChain() ChainType {
	env := config.GetEnvironment()
	
	// 测试环境默认使用 Solana Devnet
	if env == config.EnvironmentTest || env == config.EnvironmentDev {
		return ChainSolana
	}
	
	// 主网环境可以从配置读取
	chainStr := config.GetEnv("DEFAULT_CHAIN", "solana")
	return ChainType(strings.ToLower(chainStr))
}

// MintRewardMultiChain 多链奖励发放
func MintRewardMultiChain(walletAddr, behaviorType string, amount uint64) (string, error) {
	chain := GetDefaultChain()
	
	switch chain {
	case ChainSolana:
		return MintReward(walletAddr, behaviorType)
	case ChainEthereum:
		return MintRewardEthereum(walletAddr, amount)
	case ChainPolygon:
		return MintRewardPolygon(walletAddr, amount)
	case ChainBSC:
		return MintRewardBSC(walletAddr, amount)
	default:
		return MintReward(walletAddr, behaviorType)
	}
}

// MintRewardEthereum 在 Ethereum 上发放奖励
func MintRewardEthereum(walletAddr string, amount uint64) (string, error) {
	// TODO: 实现 Ethereum 代币转账
	// 需要使用 go-ethereum 或 ethers-go
	
	config.Log.Infof("Ethereum 奖励发放: wallet=%s, amount=%d", walletAddr, amount)
	
	// 模拟返回
	return "0x" + fmt.Sprintf("%064x", amount), nil
}

// MintRewardPolygon 在 Polygon 上发放奖励
func MintRewardPolygon(walletAddr string, amount uint64) (string, error) {
	// TODO: 实现 Polygon 代币转账
	// Polygon 使用与 Ethereum 相同的接口
	
	config.Log.Infof("Polygon 奖励发放: wallet=%s, amount=%d", walletAddr, amount)
	
	// 模拟返回
	return "0x" + fmt.Sprintf("%064x", amount), nil
}

// MintRewardBSC 在 BSC 上发放奖励
func MintRewardBSC(walletAddr string, amount uint64) (string, error) {
	// TODO: 实现 BSC 代币转账
	// BSC 也使用与 Ethereum 相同的接口
	
	config.Log.Infof("BSC 奖励发放: wallet=%s, amount=%d", walletAddr, amount)
	
	// 模拟返回
	return "0x" + fmt.Sprintf("%064x", amount), nil
}

// GetChainBalance 获取链上余额
func GetChainBalance(chain ChainType, walletAddr string) (uint64, error) {
	switch chain {
	case ChainSolana:
		return GetSolanaBalance(walletAddr)
	case ChainEthereum, ChainPolygon, ChainBSC:
		return GetEVMBalance(chain, walletAddr)
	default:
		return 0, fmt.Errorf("不支持的链类型: %s", chain)
	}
}

// GetSolanaBalance 获取 Solana 余额
func GetSolanaBalance(walletAddr string) (uint64, error) {
	client := rpc.New(config.GetSolanaRPCURL())
	ctx := context.Background()
	
	pubkey, err := solana.PublicKeyFromBase58(walletAddr)
	if err != nil {
		return 0, err
	}
	
	accountInfo, err := client.GetAccountInfo(ctx, pubkey)
	if err != nil {
		return 0, err
	}
	
	if accountInfo.Value == nil {
		return 0, nil
	}
	
	return accountInfo.Value.Lamports, nil
}

// GetEVMBalance 获取 EVM 链余额
func GetEVMBalance(chain ChainType, walletAddr string) (uint64, error) {
	// TODO: 实现 EVM 链余额查询
	// 需要使用 go-ethereum 或相应的 RPC 客户端
	
	config.Log.Infof("查询 %s 链余额: wallet=%s", chain, walletAddr)
	
	// 模拟返回
	return 0, nil
}

// GetSupportedChains 获取支持的链列表
func GetSupportedChains() []ChainType {
	return []ChainType{
		ChainSolana,
		ChainEthereum,
		ChainPolygon,
		ChainBSC,
	}
}
