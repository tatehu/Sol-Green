package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sol-green/config"
	"sol-green/model"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	solanaRPCTimeout = 6 * time.Second
	solanaRPCRetries = 2
)

func rewardAmountByBehavior(behaviorType string) uint64 {
	cfg := config.GetRewardConfig()
	switch behaviorType {
	case "waste_sorting":
		return cfg.WasteSortingReward
	case "tree_planting":
		return cfg.TreePlantingReward
	case "low_carbon_travel":
		return cfg.LowCarbonTravelReward
	default:
		return cfg.WasteSortingReward
	}
}

// MintRewardAmount 发放奖励（Solana：当前实现为 SOL 转账示例）
// 文档目标是 SPL Token（SOLGREEN），后续可在此替换为 SPL Token transfer。
func MintRewardAmount(walletAddr string, amount uint64) (string, error) {
	if amount == 0 {
		return "", fmt.Errorf("奖励数量必须大于 0")
	}

	rpcURL := config.GetSolanaRPCURL()
	client := rpc.New(rpcURL)

	adminPrivKey, err := loadAdminPrivateKey()
	if err != nil {
		if config.IsDev() || config.IsTest() {
			config.Log.Warnf("MintRewardAmount 降级为 mock（%v）", err)
			return fmt.Sprintf("mock_tx_%d", time.Now().UnixNano()), nil
		}
		return "", err
	}

	adminInfo, err := rpcRetry("GetAccountInfo(admin)", func(ctx context.Context) (*rpc.GetAccountInfoResult, error) {
		return client.GetAccountInfo(ctx, adminPrivKey.PublicKey())
	})
	if err != nil || adminInfo == nil || adminInfo.Value == nil {
		return "", fmt.Errorf("管理员账户在当前网络不存在（请确认 SOLANA_RPC_URL 对应网络，并给管理员地址先空投/转入少量 SOL）: %v", err)
	}
	if adminInfo.Value.Executable {
		return "", fmt.Errorf("管理员地址是可执行程序账户（program），不能作为支付手续费的钱包；请配置普通钱包私钥到 SOLANA_ADMIN_PRIVKEY，并给该钱包空投 SOL")
	}

	receiverPubkey, err := solana.PublicKeyFromBase58(walletAddr)
	if err != nil {
		return "", fmt.Errorf("钱包地址无效: %v", err)
	}

	receiverInfo, err := rpcRetry("GetAccountInfo(receiver)", func(ctx context.Context) (*rpc.GetAccountInfoResult, error) {
		return client.GetAccountInfo(ctx, receiverPubkey)
	})
	if err != nil || receiverInfo == nil || receiverInfo.Value == nil {
		return "", fmt.Errorf("接收者账户在当前网络不存在（可能从未收到过 SOL/空投）。请先给该地址空投/转入少量 SOL 激活账户后再试: %v", err)
	}

	recent, err := rpcRetry("GetLatestBlockhash", func(ctx context.Context) (*rpc.GetLatestBlockhashResult, error) {
		return client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	})
	if err != nil {
		return "", fmt.Errorf("获取最新区块哈希失败: %v", err)
	}

	adminBal, balErr := rpcRetry("GetBalance(admin)", func(ctx context.Context) (*rpc.GetBalanceResult, error) {
		return client.GetBalance(ctx, adminPrivKey.PublicKey(), rpc.CommitmentProcessed)
	})
	if balErr == nil {
		const feeBuffer uint64 = 10_000
		if adminBal != nil && adminBal.Value < amount+feeBuffer {
			return "", fmt.Errorf("管理员 SOL 余额不足：balance=%d lamports，需要至少 %d lamports（含手续费缓冲）。请先给管理员地址空投/转入 SOL",
				adminBal.Value, amount+feeBuffer)
		}
	}

	instruction := system.NewTransferInstruction(
		amount,
		adminPrivKey.PublicKey(),
		receiverPubkey,
	).Build()

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(adminPrivKey.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("创建交易失败: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(adminPrivKey.PublicKey()) {
			return &adminPrivKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("签名交易失败: %v", err)
	}

	sig, err := rpcRetry("SendTransaction", func(ctx context.Context) (solana.Signature, error) {
		return client.SendTransaction(ctx, tx)
	})
	if err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}

	return sig.String(), nil
}

func isRetryableNetErr(err error) bool {
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
	// 兜底：部分库会把网络错误包装成字符串
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "timed out") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "tls handshake timeout")
}

func withRPCTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, solanaRPCTimeout)
}

func rpcRetry[T any](opName string, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	var lastErr error
	for attempt := 0; attempt <= solanaRPCRetries; attempt++ {
		ctx, cancel := withRPCTimeout(context.Background())
		v, err := fn(ctx)
		cancel()
		if err == nil {
			return v, nil
		}
		lastErr = err
		if !isRetryableNetErr(err) || attempt == solanaRPCRetries {
			return zero, fmt.Errorf("%s: %w", opName, err)
		}
		time.Sleep(time.Duration(attempt+1) * 300 * time.Millisecond)
	}
	return zero, fmt.Errorf("%s: %w", opName, lastErr)
}

func loadAdminPrivateKey() (solana.PrivateKey, error) {
	adminPrivKeyStr := strings.TrimSpace(config.SolanaAdminPrivKey())
	if adminPrivKeyStr == "" {
		return nil, fmt.Errorf("管理员私钥未配置（SOLANA_ADMIN_PRIVKEY）")
	}

	// 兼容 Solana CLI 导出的 JSON 数组私钥（~/.config/solana/id.json 的内容）
	if strings.HasPrefix(adminPrivKeyStr, "[") {
		var arr []int
		if err := json.Unmarshal([]byte(adminPrivKeyStr), &arr); err != nil {
			return nil, fmt.Errorf("管理员私钥 JSON 解析失败: %v", err)
		}
		b := make([]byte, len(arr))
		for i, v := range arr {
			if v < 0 || v > 255 {
				return nil, fmt.Errorf("管理员私钥 JSON 数组包含非法值: index=%d value=%d", i, v)
			}
			b[i] = byte(v)
		}
		if len(b) != 64 {
			return nil, fmt.Errorf("管理员私钥长度不正确: got=%d want=64（请提供 64 字节 secret key JSON 数组）", len(b))
		}
		return solana.PrivateKey(b), nil
	}

	// Base58 格式私钥（应为 64 字节）
	priv, err := solana.PrivateKeyFromBase58(adminPrivKeyStr)
	if err != nil {
		return nil, fmt.Errorf("管理员私钥 Base58 解析失败: %v", err)
	}
	if len(priv) != 64 {
		// 32 字节通常意味着你传的是公钥/ProgramID/seed，而不是 secret key
		return nil, fmt.Errorf("管理员私钥长度不正确: got=%d want=64（你可能把公钥/ProgramID 配到了 SOLANA_ADMIN_PRIVKEY）", len(priv))
	}
	return priv, nil
}

// MintReward 发放奖励（平台代币 SOLGREEN）
func MintReward(walletAddr string, behaviorType string) (string, error) {
	// 按环境配置计算奖励（对齐文档的奖励配置入口）
	rewardAmount := rewardAmountByBehavior(behaviorType)
	return MintRewardAmount(walletAddr, rewardAmount)
}

// GenerateBehaviorHash 生成行为哈希（用于链上存证）
func GenerateBehaviorHash(behavior *model.GreenBehavior) string {
	data := fmt.Sprintf("%s:%s:%s", behavior.WalletAddr, behavior.BehaviorType, behavior.SubmitTime.Format("2006-01-02T15:04:05Z07:00"))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// RecordOnChain 链上存证（将哈希写入 Solana 账户）
func RecordOnChain(proofHash, walletAddr string) (string, error) {
	// 简化版：调用自定义存证合约，将哈希写入链上账户
	// 实际项目需部署 Solana 存证合约，此处为示例
	rpcURL := config.GetSolanaRPCURL()
	client := rpc.New(rpcURL)

	adminPrivKey, err := loadAdminPrivateKey()
	if err != nil {
		if config.IsDev() || config.IsTest() {
			config.Log.Warnf("RecordOnChain 降级为 mock（%v）", err)
			return "mock_proof_hash_" + proofHash[:16], nil
		}
		return "", err
	}

	proofContractStr := config.SolanaProofContract()
	if proofContractStr == "" {
		// 如果没有配置存证合约，返回模拟哈希
		config.Log.Warn("存证合约未配置，跳过链上存证")
		return "mock_proof_hash_" + proofHash[:16], nil
	}

	proofContract, err := solana.PublicKeyFromBase58(proofContractStr)
	if err != nil {
		return "", fmt.Errorf("存证合约地址无效: %v", err)
	}

	receiverPubkey, err := solana.PublicKeyFromBase58(walletAddr)
	if err != nil {
		return "", fmt.Errorf("钱包地址无效: %v", err)
	}

	// 构建存证指令（自定义合约）
	recent, err := rpcRetry("GetLatestBlockhash", func(ctx context.Context) (*rpc.GetLatestBlockhashResult, error) {
		return client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	})
	if err != nil {
		return "", fmt.Errorf("获取最新区块哈希失败: %v", err)
	}

	instructions := []solana.Instruction{
		solana.NewInstruction(
			proofContract,
			solana.AccountMetaSlice{
				{PublicKey: adminPrivKey.PublicKey(), IsSigner: true, IsWritable: true},
				{PublicKey: receiverPubkey, IsSigner: false, IsWritable: false},
			},
			[]byte(proofHash),
		),
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(adminPrivKey.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("创建交易失败: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(adminPrivKey.PublicKey()) {
			return &adminPrivKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("签名交易失败: %v", err)
	}

	txHash, err := rpcRetry("SendTransaction", func(ctx context.Context) (solana.Signature, error) {
		return client.SendTransaction(ctx, tx)
	})
	if err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}

	return txHash.String(), nil
}
