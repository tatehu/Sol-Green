package service

import (
	"context"
	"fmt"
	"os"
	"sol-green/config"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// BlockchainListener 区块链监听服务
type BlockchainListener struct {
	solanaClient *rpc.Client
	wsClient     *ws.Client
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
	// 监听程序账户变化
	programID := solana.MustPublicKeyFromBase58(config.SolanaProofContract())
	if programID.IsZero() {
		return fmt.Errorf("程序 ID 未配置")
	}

	// 使用 WebSocket 订阅账户变化
	wsURL := config.GetSolanaWSURL()
	wsClient, err := ws.Connect(bl.ctx, wsURL)
	if err != nil {
		config.Log.Warnf("WebSocket 连接失败，使用轮询模式: %v", err)
		return bl.startPolling()
	}

	bl.wsClient = wsClient

	// 订阅账户通知
	sub, err := wsClient.AccountSubscribe(programID, rpc.CommitmentFinalized)
	if err != nil {
		config.Log.Warnf("账户订阅失败，使用轮询模式: %v", err)
		return bl.startPolling()
	}

	go bl.handleAccountUpdates(sub)

	config.Log.Info("区块链监听服务已启动")
	return nil
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

// handleAccountUpdates 处理账户更新
func (bl *BlockchainListener) handleAccountUpdates(sub *ws.AccountSubscription) {
	defer sub.Unsubscribe()
	for {
		select {
		case <-bl.ctx.Done():
			return
		case update, ok := <-sub.Recv():
			if !ok {
				config.Log.Warn("账户订阅通道已关闭")
				return
			}
			if update != nil {
				bl.processAccountUpdate(update)
			}
		}
	}
}

// processAccountUpdate 处理账户更新
func (bl *BlockchainListener) processAccountUpdate(update *ws.AccountResult) {
	if update == nil || update.Value == nil {
		return
	}
	// 解析账户数据，提取交易信息
	config.Log.Infof("收到账户更新: %v", update.Value.Account.Owner.String())
	
	// 记录到数据库
	// TODO: 保存交易数据到数据库
}

// pollTransactions 轮询交易
func (bl *BlockchainListener) pollTransactions() {
	// 获取最近的交易
	signatures, err := bl.solanaClient.GetSignaturesForAddress(
		bl.ctx,
		solana.MustPublicKeyFromBase58(config.SolanaProofContract()),
		&rpc.GetSignaturesForAddressOpts{
			Limit: rpc.UintPtr(10),
		},
	)
	if err != nil {
		config.Log.Errorf("获取交易签名失败: %v", err)
		return
	}

	for _, sig := range signatures {
		bl.processTransaction(sig.Signature.String())
	}
}

// processTransaction 处理交易
func (bl *BlockchainListener) processTransaction(signature string) {
	// 获取交易详情
	tx, err := bl.solanaClient.GetTransaction(
		bl.ctx,
		solana.MustSignatureFromBase58(signature),
		&rpc.GetTransactionOpts{},
	)
	if err != nil {
		config.Log.Errorf("获取交易详情失败: %v", err)
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
	if bl.wsClient != nil {
		bl.wsClient.Close()
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

	// 获取程序账户信息
	programID := solana.MustPublicKeyFromBase58(config.SolanaProofContract())
	accountInfo, err := client.GetAccountInfo(ctx, programID)
	if err != nil {
		config.Log.Warnf("获取程序账户信息失败: %v", err)
	}

	stats := map[string]interface{}{
		"current_slot":    slot,
		"program_id":      programID.String(),
		"account_exists":  accountInfo != nil,
		"timestamp":       time.Now().Unix(),
	}

	return stats, nil
}
