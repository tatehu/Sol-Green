# 🌱 Sol-Green - 环保行为奖励平台

基于 Solana 区块链的环保行为记录与奖励系统，通过 Web3 技术激励用户参与环保活动。

## 📋 项目简介

Sol-Green 是一个去中心化的环保奖励平台，用户可以通过提交环保行为（垃圾分类、植树、低碳出行等）获得代币奖励。项目采用 AI 反欺诈检测、链上存证、第三方机构认证等技术，确保行为的真实性和可信度。

## ✨ 核心功能

- 🌱 **环保行为记录**: 垃圾分类、植树、低碳出行等行为上链存证
- 🎯 **挑战活动**: 发起和参与环保挑战，获得额外奖励
- 🎁 **营销活动**: 签到、邀请、每日任务、抽奖等多种营销活动
- 🤖 **AI 反欺诈**: 集成全球领先的 AI 检测技术，确保行为真实性
- 🌍 **第三方认证**: 支持联合国环境规划署、WWF、绿色和平等全球知名机构认证
- 💰 **代币奖励**: 基于 Solana 区块链的即时奖励发放
- 🔒 **链上存证**: 所有行为永久记录在 Solana 链上
- 🔍 **链上监听**: 实时监听链上事件，获取运行数据
- 📊 **数据分析**: 完整的运营数据分析和报告
- 🌐 **多链支持**: 支持 Solana、Ethereum、Polygon、BSC 等多链
- 🔄 **合约升级**: 完整的合约升级管理系统

## 🏗️ 项目结构

```
trends_solana/
├── backend/          # Go 后端服务
│   ├── config/      # 配置管理
│   ├── controller/  # 控制器
│   ├── middleware/  # 中间件
│   ├── model/       # 数据模型
│   ├── service/     # 业务逻辑
│   ├── tests/       # 测试文件
│   ├── main.go      # 应用入口
│   ├── go.mod       # Go 依赖
│   └── Dockerfile   # 后端镜像
│
├── frontend/        # React 前端应用
│   ├── src/        # 源代码
│   ├── public/      # 静态文件
│   ├── package.json
│   └── Dockerfile   # 前端镜像
│
├── contracts/       # Solana 智能合约
│   └── sol-green/   # 合约项目
│
├── docs/           # 项目文档
│   ├── DEVELOPMENT.md
│   ├── TECHNICAL.md
│   ├── DEPLOYMENT.md
│   └── API.md
│
├── scripts/        # 部署脚本
│   ├── deploy.sh
│   ├── test.sh
│   └── build-contract.sh
│
├── docker-compose.yml  # Docker 编排
└── .env.example        # 环境变量示例
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose (可选)
- Solana CLI (合约部署)
- Anchor CLI (合约开发)

### 1. 克隆项目

```bash
git clone <repository-url>
cd trends_solana
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入必要的配置
```

### 3. 启动后端服务

#### 方式一：使用 Docker Compose（推荐）

```bash
docker-compose up -d
```

#### 方式二：本地运行

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 启动 PostgreSQL 和 Redis（如果使用）
cd ..
docker-compose up -d postgres redis

# 运行后端
cd backend
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 4. 启动前端

```bash
cd frontend
npm install
npm start
```

前端应用将在 `http://localhost:3000` 启动。

### 5. 部署智能合约

```bash
# 安装 Anchor CLI
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install latest
avm use latest

# 构建合约
./scripts/build-contract.sh

# 部署到测试网
cd contracts/sol-green
anchor deploy --provider.cluster devnet
```

## 📡 API 文档

详细 API 文档请查看 [docs/API.md](./docs/API.md)

### 主要接口

- `POST /api/v1/auth/wallet` - 钱包登录
- `POST /api/v1/green/behavior/submit` - 提交环保行为
- `GET /api/v1/green/behavior/:id` - 查询行为状态
- `POST /api/v1/green/partner/verify` - 第三方认证
- `GET /api/v1/health` - 健康检查

## 🧪 测试

```bash
# 运行所有测试
./scripts/test.sh

# 运行后端测试
cd backend
go test -v ./tests/...

# 运行前端测试
cd frontend
npm test
```

## 🚢 部署

### 使用 Docker Compose

```bash
./scripts/deploy.sh
```

### 手动部署

详细部署文档请查看 [docs/DEPLOYMENT.md](./docs/DEPLOYMENT.md)

## 📚 文档 / Documentation

### v3.0.0 文档 / v3.0.0 Documentation

- [文档索引](./docs/DOCUMENTATION_INDEX_v3.0.0.md) / [Documentation Index](./docs/DOCUMENTATION_INDEX_v3.0.0.md)
- [更新日志](./docs/CHANGELOG_v3.0.0_CN.md) / [Changelog](./docs/CHANGELOG_v3.0.0_EN.md)
- [快速开始指南](./docs/QUICKSTART_v3.0.0_CN.md) / [Quick Start](./docs/QUICKSTART_v3.0.0_EN.md)
- [API 文档](./docs/API_v3.0.0_CN.md) / [API Documentation](./docs/API_v3.0.0_EN.md)
- [开发文档](./docs/DEVELOPMENT_v3.0.0_CN.md) / [Development Guide](./docs/DEVELOPMENT_v3.0.0_EN.md)
- [技术文档](./docs/TECHNICAL_v3.0.0_CN.md) / [Technical Documentation](./docs/TECHNICAL_v3.0.0_EN.md)
- [部署文档](./docs/DEPLOYMENT_v3.0.0_CN.md) / [Deployment Guide](./docs/DEPLOYMENT_v3.0.0_EN.md)
- [完整功能说明](./docs/COMPLETE_FEATURES_v3.0.0_CN.md) / [Complete Features](./docs/COMPLETE_FEATURES_v3.0.0_EN.md)

### 功能文档 / Feature Documentation

- [挑战活动](./docs/CHALLENGES_FEATURE_v3.0.0_CN.md) / [Challenge Activities](./docs/CHALLENGES_FEATURE_v3.0.0_EN.md)
- [营销活动系统](./docs/MARKETING_ACTIVITIES_v3.0.0_CN.md) / [Marketing Activities](./docs/MARKETING_ACTIVITIES_v3.0.0_EN.md)
- [AI 反欺诈检测](./docs/AI_FRAUD_DETECTION_v3.0.0_CN.md) / [AI Fraud Detection](./docs/AI_FRAUD_DETECTION_v3.0.0_EN.md)
- [第三方机构认证](./docs/THIRD_PARTY_VERIFICATION_v3.0.0_CN.md) / [Third-Party Verification](./docs/THIRD_PARTY_VERIFICATION_v3.0.0_EN.md)
- [链上数据监听](./docs/BLOCKCHAIN_MONITORING_v3.0.0_CN.md) / [Blockchain Monitoring](./docs/BLOCKCHAIN_MONITORING_v3.0.0_EN.md)
- [多链支持](./docs/MULTICHAIN_SUPPORT_v3.0.0_CN.md) / [Multi-Chain Support](./docs/MULTICHAIN_SUPPORT_v3.0.0_EN.md)
- [合约升级管理](./docs/CONTRACT_UPGRADE_v3.0.0_CN.md) / [Contract Upgrade](./docs/CONTRACT_UPGRADE_v3.0.0_EN.md)
- [运营管理系统](./docs/OPERATIONS_v3.0.0_CN.md) / [Operations Management](./docs/OPERATIONS_v3.0.0_EN.md)

## 🔧 技术栈

### 后端
- Go 1.21+
- Gin Web 框架
- GORM ORM
- Redis 缓存
- JWT 认证
- Solana Go SDK

### 区块链
- Solana
- Anchor 框架
- Rust
- SPL Token

### 前端
- React 18
- Solana Wallet Adapter
- Axios
- React Router

## ✨ 核心特性

1. **挑战活动系统**: 发起和参与环保挑战，提高用户参与度
2. **AI 反欺诈检测**: 集成百度 AI、阿里云、AWS、Google Cloud 等全球领先的 AI 检测技术
3. **第三方机构认证**: 支持 UNEP、WWF、Greenpeace、阿拉善 SEE 等全球知名环保组织认证
4. **多环境支持**: 开发、测试（Solana Devnet）、主网环境一键切换
5. **链上存证**: 行为哈希永久记录在 Solana 链上
6. **防刷机制**: Redis 限流 + 防重复提交
7. **安全机制**: JWT 认证、CORS、SQL 注入防护

## 🎯 挑战活动

用户可以发起和参与环保挑战活动：

- **个人挑战**: 个人环保目标挑战
- **团队挑战**: 团队协作完成环保任务
- **社区挑战**: 社区级别的环保活动

挑战完成后，参与者可获得额外奖励。

详细文档：[挑战活动 API](./docs/API.md#挑战活动接口)

## 🔒 安全说明

1. **私钥管理**: 严禁将私钥提交到代码仓库，使用环境变量管理
2. **JWT 密钥**: 生产环境必须使用强随机密钥
3. **API 限流**: 已实现基于 Redis 的请求限流
4. **AI 反欺诈**: 集成第三方内容安全 API 进行图片/视频检测

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 📞 联系方式

如有问题，请提交 Issue 或联系项目维护者。
