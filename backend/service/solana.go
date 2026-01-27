package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sol-green/config"
	"sol-green/model"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// MintReward 发放奖励（平台代币 SOLGREEN）
func MintReward(walletAddr string, behaviorType string) (string, error) {
	// 1. 连接 Solana 主网/测试网
	rpcURL := config.GetSolanaRPCURL()
	client := rpc.New(rpcURL)

	// 2. 加载管理员钱包（私钥需从环境变量读取，严禁硬编码）
	adminPrivKeyStr := config.SolanaAdminPrivKey()
	if adminPrivKeyStr == "" {
		return "", fmt.Errorf("管理员私钥未配置")
	}

	adminPrivKey, err := solana.PrivateKeyFromBase58(adminPrivKeyStr)
	if err != nil {
		return "", fmt.Errorf("私钥解析失败: %v", err)
	}

	// 3. 计算奖励数量（按行为类型）
	var rewardAmount uint64
	switch behaviorType {
	case "waste_sorting":
		rewardAmount = 1000 // 1000 SOLGREEN 代币
	case "tree_planting":
		rewardAmount = 5000
	case "low_carbon_travel":
		rewardAmount = 2000
	default:
		rewardAmount = 1000
	}

	// 4. 验证接收者地址
	receiverPubkey, err := solana.PublicKeyFromBase58(walletAddr)
	if err != nil {
		return "", fmt.Errorf("钱包地址无效: %v", err)
	}

	// 5. 构建 SOL 转账交易（简化版，实际应使用 SPL Token）
	// 注意：实际项目中应使用 SPL Token 转账，这里简化为 SOL 转账示例
	ctx := context.Background()
	recent, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("获取最新区块哈希失败: %v", err)
	}

	// 构建转账指令（使用系统程序）
	instruction := solana.NewInstruction(
		solana.SystemProgramID,
		solana.AccountMetaSlice{
			{PublicKey: adminPrivKey.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: receiverPubkey, IsSigner: false, IsWritable: true},
		},
		[]byte{0x02}, // Transfer instruction discriminator
		// 注意：实际应包含 lamports 数量，这里简化处理
	)

	// 6. 创建并签名交易
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

	// 7. 发送交易到链上
	sig, err := client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}

	config.Log.Infof("奖励发放成功: 交易哈希=%s, 金额=%d, 接收者=%s", sig, rewardAmount, walletAddr)

	return sig.String(), nil
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

	adminPrivKeyStr := config.SolanaAdminPrivKey()
	if adminPrivKeyStr == "" {
		return "", fmt.Errorf("管理员私钥未配置")
	}

	adminPrivKey, err := solana.PrivateKeyFromBase58(adminPrivKeyStr)
	if err != nil {
		return "", fmt.Errorf("私钥解析失败: %v", err)
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
	ctx := context.Background()
	recent, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("获取最新区块哈希失败: %v", err)
	}

	instructions := []solana.Instruction{
		solana.NewInstruction(
			proofContract,
			[][]byte{[]byte(proofHash)},
			solana.AccountMetaSlice{
				{PublicKey: adminPrivKey.PublicKey(), IsSigner: true, IsWritable: true},
				{PublicKey: receiverPubkey, IsSigner: false, IsWritable: false},
			},
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

	txHash, err := client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}

	return txHash.String(), nil
}
