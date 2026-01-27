# 快速开始指南 v3.0.0 / Quick Start Guide v3.0.0

## 5 分钟快速体验 / 5-Minute Quick Experience

### 1. 环境准备 / Environment Setup

确保已安装 / Ensure installed:
- Go 1.21+ (或使用 Docker / or use Docker)
- Node.js 18+
- Docker & Docker Compose (推荐 / Recommended)

### 2. 克隆项目 / Clone Project

```bash
git clone <repository-url>
cd trends_solana
```

### 3. 配置环境变量 / Configure Environment

```bash
cp .env.example .env
# 编辑 .env 文件（可选，开发环境可使用默认值）
# Edit .env file (optional, defaults can be used for development)
```

**详细配置说明 / Detailed Configuration Guide**:  
请参考 [配置获取指南](./CONFIG_SETUP_GUIDE_v3.0.0_CN.md) 了解如何获取各项配置的值。  
Please refer to [Configuration Setup Guide](./CONFIG_SETUP_GUIDE_v3.0.0_EN.md) for details on how to obtain configuration values.

### 4. 启动服务（最简单方式）/ Start Services (Easiest Way)

```bash
# 使用 Docker Compose 一键启动
# Start with Docker Compose
docker-compose up -d

# 查看服务状态
# Check service status
docker-compose ps

# 查看日志
# View logs
docker-compose logs -f backend
```

### 5. 访问应用 / Access Application

- **后端 API / Backend API**: http://localhost:8080
- **健康检查 / Health Check**: http://localhost:8080/api/v1/health
- **前端应用 / Frontend App**: 需要单独启动（见下方 / needs separate start, see below）

### 6. 启动前端 / Start Frontend

```bash
cd frontend
npm install
npm start
```

前端将在 http://localhost:3000 启动  
Frontend will start at http://localhost:3000

### 7. 测试 API / Test API

```bash
# 健康检查 / Health check
curl http://localhost:8080/api/v1/health

# 钱包登录（示例）/ Wallet login (example)
curl -X POST http://localhost:8080/api/v1/auth/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_addr": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
    "signature": "test_signature",
    "message": "test_message"
  }'
```

## 本地开发模式 / Local Development Mode

### 后端开发 / Backend Development

```bash
# 进入后端目录
# Enter backend directory
cd backend

# 安装依赖
# Install dependencies
go mod download

# 启动 PostgreSQL 和 Redis（使用 Docker）
# Start PostgreSQL and Redis (using Docker)
cd ..
docker-compose up -d postgres redis

# 配置环境变量（使用 SQLite 可跳过数据库）
# Configure environment (can skip database if using SQLite)
export DB_TYPE=sqlite

# 运行后端
# Run backend
cd backend
go run main.go
```

后端将在 `http://localhost:8080` 启动  
Backend will start at `http://localhost:8080`

### 前端开发 / Frontend Development

```bash
cd frontend
npm install
npm start
```

前端将在 `http://localhost:3000` 启动  
Frontend will start at `http://localhost:3000`

### 智能合约开发 / Smart Contract Development

```bash
# 安装 Anchor CLI
# Install Anchor CLI
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install latest
avm use latest

# 进入合约目录
# Enter contract directory
cd contracts/sol-green

# 构建合约
# Build contract
anchor build

# 运行测试
# Run tests
anchor test
```

## 验证项目 / Verify Project

运行验证脚本检查项目完整性：  
Run verification script to check project integrity:

```bash
./scripts/verify.sh
```

## 常见问题 / Common Issues

### Q: 后端启动失败？/ Backend startup failed?

**A:** 检查 / Check:
1. 端口 8080 是否被占用 / Port 8080 occupied
2. 数据库/Redis 是否正常启动 / Database/Redis running
3. 环境变量配置是否正确 / Environment variables correct

解决方案 / Solution:
```bash
# 检查端口占用
# Check port usage
lsof -i :8080

# 检查 Docker 服务
# Check Docker services
docker-compose ps

# 查看后端日志
# View backend logs
docker-compose logs backend
```

### Q: 前端无法连接后端？/ Frontend cannot connect to backend?

**A:** 
1. 检查后端是否运行在 http://localhost:8080 / Check backend running at http://localhost:8080
2. 检查 CORS 配置 / Check CORS configuration
3. 查看浏览器控制台错误信息 / Check browser console errors

### Q: Solana 交易失败？/ Solana transaction failed?

**A:**
1. 检查钱包是否有足够的 SOL / Check wallet has enough SOL
2. 确认网络配置（devnet/mainnet）/ Confirm network configuration
3. 验证 RPC 节点是否可用 / Verify RPC node available

## 下一步 / Next Steps

- 阅读 [README.md](../README.md) 了解项目详情 / Read [README.md](../README.md) for project details
- 查看 [开发文档](./DEVELOPMENT_v3.0.0_CN.md) 开始开发 / View [Development Guide](./DEVELOPMENT_v3.0.0_CN.md) to start development
- 参考 [API 文档](./API_v3.0.0_CN.md) 了解接口 / Refer to [API Documentation](./API_v3.0.0_CN.md) for API details
- 查看 [部署文档](./DEPLOYMENT_v3.0.0_CN.md) 进行生产部署 / View [Deployment Guide](./DEPLOYMENT_v3.0.0_CN.md) for production deployment

---

**版本 / Version**: v3.0.0  
**更新日期 / Update Date**: 2024-01-24
