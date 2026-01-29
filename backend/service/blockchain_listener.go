package service

import (
	"context"
	"errors"
	"net"
	"sol-green/config"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	listenerRPCTimeout = 5 * time.Second
	listenerRPCRetries = 2
)

func listenerIsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout() || ne.Temporary()
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") || strings.Contains(msg, "timed out")
}

func (bl *BlockchainListener) getSignaturesWithRetry(ctx context.Context, programID solana.PublicKey) ([]*rpc.TransactionSignature, error) {
	var lastErr error
	for attempt := 0; attempt <= listenerRPCRetries; attempt++ {
		callCtx, cancel := context.WithTimeout(ctx, listenerRPCTimeout)
		limit := 20
		sigs, err := bl.solanaClient.GetSignaturesForAddressWithOpts(
			callCtx,
			programID,
			&rpc.GetSignaturesForAddressOpts{
				Limit:      &limit,
				Commitment: rpc.CommitmentConfirmed,
			},
		)
		cancel()
		if err == nil {
			return sigs, nil
		}
		lastErr = err
		if !listenerIsRetryable(err) || attempt == listenerRPCRetries {
			return nil, err
		}
		time.Sleep(time.Duration(attempt+1) * 300 * time.Millisecond)
	}
	return nil, lastErr
}

// BlockchainListener 区块链监听服务
type BlockchainListener struct {
	solanaClient *rpc.Client
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewBlockchainListener 创建区块链监听器
func NewBlockchainListener() (*BlockchainListener, error) {
	rpcURL := config.GetSolanaRPCURL()
	client := rpc.New(rpcURL)

	ctx, cancel := context.WithCancel(context.Background())

	return &BlockchainListener{
		solanaClient: client,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// StartListening 开始监听链上事件
func (bl *BlockchainListener) StartListening() error {
	// 简化版本：当前只启用轮询模式，避免依赖 WebSocket API 变更导致编译错误。
	// Simple version: only use polling for now.
	return bl.startPolling()
}

// startPolling 启动轮询模式
func (bl *BlockchainListener) startPolling() error {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-bl.ctx.Done():
				return
			case <-ticker.C:
				bl.pollTransactions()
			}
		}
	}()

	config.Log.Info("区块链轮询服务已启动")
	return nil
}

// pollTransactions 轮询交易
func (bl *BlockchainListener) pollTransactions() {
	contractStr := config.SolanaProofContract()
	if contractStr == "" {
		config.Log.Warn("存证合约地址未配置，跳过交易轮询")
		return
	}

	programID := solana.MustPublicKeyFromBase58(contractStr)
	if programID.IsZero() {
		config.Log.Warn("程序 ID 无效，跳过交易轮询")
		return
	}

	// 获取最近的交易（超时 + 小范围重试 + 数量限制）
	signatures, err := bl.getSignaturesWithRetry(bl.ctx, programID)
	if err != nil {
		// devnet 网络不稳定时避免刷屏 error
		if listenerIsRetryable(err) {
			config.Log.Warnf("获取交易签名失败（可能网络/RPC 不可达）: %v", err)
		} else {
			config.Log.Errorf("获取交易签名失败: %v", err)
		}
		return
	}

	if len(signatures) == 0 {
		config.Log.Debug("未获取到新的交易签名")
		return
	}

	for _, sig := range signatures {
		bl.processTransaction(sig.Signature.String())
	}
}

// processTransaction 处理交易
func (bl *BlockchainListener) processTransaction(signature string) {
	// 获取交易详情
	ctx, cancel := context.WithTimeout(bl.ctx, listenerRPCTimeout)
	defer cancel()
	tx, err := bl.solanaClient.GetTransaction(
		ctx,
		solana.MustSignatureFromBase58(signature),
		&rpc.GetTransactionOpts{},
	)
	if err != nil {
		if listenerIsRetryable(err) {
			config.Log.Warnf("获取交易详情失败（可能网络/RPC 不可达）: %v", err)
		} else {
			config.Log.Errorf("获取交易详情失败: %v", err)
		}
		return
	}

	// 解析交易，提取事件
	// TODO: 解析交易数据，提取奖励发放、存证等事件

	config.Log.Infof("处理交易: %s, 区块时间: %v", signature, tx.BlockTime)
}

// Stop 停止监听
func (bl *BlockchainListener) Stop() {
	if bl.cancel != nil {
		bl.cancel()
	}
	config.Log.Info("区块链监听服务已停止")
}

// GetChainStats 获取链上统计
func GetChainStats() (map[string]interface{}, error) {
	client := rpc.New(config.GetSolanaRPCURL())
	ctx := context.Background()

	// 获取最新区块高度
	slot, err := client.GetSlot(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"current_slot": slot,
		"timestamp":    time.Now().Unix(),
	}

	// 获取程序账户信息（如果已配置存证合约地址）
	contractStr := config.SolanaProofContract()
	if contractStr != "" {
		programID := solana.MustPublicKeyFromBase58(contractStr)
		accountInfo, err := client.GetAccountInfo(ctx, programID)
		if err != nil {
			config.Log.Warnf("获取程序账户信息失败: %v", err)
		}
		stats["program_id"] = programID.String()
		stats["account_exists"] = accountInfo != nil
	}

	return stats, nil
}
