# 多链支持系统

## 📋 概述

Sol-Green 支持多链部署，用户可以选择在不同区块链上参与环保活动并获得奖励。

## 🌐 支持的区块链

### 1. Solana（主链）

**特点：**
- 高吞吐量
- 低交易费用
- 快速确认

**配置：**
```bash
DEFAULT_CHAIN=solana
SOLANA_RPC_URL=https://api.devnet.solana.com
```

### 2. Ethereum

**特点：**
- 生态成熟
- 用户基数大
- 智能合约丰富

**配置：**
```bash
ETHEREUM_RPC_URL=https://eth-goerli.g.alchemy.com/v2/YOUR_API_KEY
ETHEREUM_PRIVATE_KEY=your-private-key
ETHEREUM_TOKEN_ADDRESS=0x...
```

### 3. Polygon

**特点：**
- 低 Gas 费用
- 与 Ethereum 兼容
- 快速交易

**配置：**
```bash
POLYGON_RPC_URL=https://polygon-mumbai.g.alchemy.com/v2/YOUR_API_KEY
POLYGON_PRIVATE_KEY=your-private-key
POLYGON_TOKEN_ADDRESS=0x...
```

### 4. BSC (Binance Smart Chain)

**特点：**
- 低交易费用
- 高性能
- 币安生态

**配置：**
```bash
BSC_RPC_URL=https://data-seed-prebsc-1-s1.binance.org:8545
BSC_PRIVATE_KEY=your-private-key
BSC_TOKEN_ADDRESS=0x...
```

## 🔄 跨链功能

### 链选择

用户可以在提交行为时选择奖励发放的链：

```json
{
  "behavior_type": "waste_sorting",
  "media_urls": ["..."],
  "chain": "polygon"  // 选择链
}
```

### 多链余额查询

```bash
GET /api/v1/wallet/balance?chain=solana
GET /api/v1/wallet/balance?chain=polygon
```

### 跨链转账（未来）

- 支持在不同链间转移代币
- 使用跨链桥接协议
- 统一代币标准

## 💰 代币标准

### Solana

- **SPL Token**: Solana 代币标准
- **Token Mint**: 代币铸造地址
- **Token Account**: 代币账户

### EVM 链 (Ethereum/Polygon/BSC)

- **ERC-20**: 标准代币接口
- **Token Contract**: 代币合约地址
- **Balance**: 账户余额

## 🔧 配置管理

### 环境变量

```bash
# 默认链
DEFAULT_CHAIN=solana

# Solana 配置
SOLANA_RPC_URL=...
SOLANA_ADMIN_PRIVKEY=...

# Ethereum 配置
ETHEREUM_RPC_URL=...
ETHEREUM_PRIVATE_KEY=...
ETHEREUM_TOKEN_ADDRESS=...

# Polygon 配置
POLYGON_RPC_URL=...
POLYGON_PRIVATE_KEY=...
POLYGON_TOKEN_ADDRESS=...

# BSC 配置
BSC_RPC_URL=...
BSC_PRIVATE_KEY=...
BSC_TOKEN_ADDRESS=...
```

## 📊 链上数据统计

### 按链统计

- 各链交易数量
- 各链用户数量
- 各链奖励发放
- 各链活跃度

### 跨链对比

- 交易费用对比
- 确认时间对比
- 用户体验对比

## 🚀 部署指南

### Solana 部署

```bash
cd contracts/sol-green
anchor build
anchor deploy --provider.cluster devnet
```

### EVM 链部署

```bash
# 使用 Hardhat 或 Truffle
npx hardhat deploy --network polygon
```

## 🔒 安全考虑

1. **私钥管理**: 每个链使用独立的私钥
2. **权限控制**: 限制跨链操作权限
3. **审计要求**: 多链合约需要分别审计
4. **监控告警**: 监控各链异常情况

## 📚 参考资料

- [Solana 文档](https://docs.solana.com/)
- [Ethereum 文档](https://ethereum.org/developers/)
- [Polygon 文档](https://docs.polygon.technology/)
- [BSC 文档](https://docs.binance.org/)
