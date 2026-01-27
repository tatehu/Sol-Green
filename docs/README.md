# 🌱 Sol-Green - 环保行为奖励平台

基于 Solana 区块链的环保行为记录与奖励系统，通过 Web3 技术激励用户参与环保活动。

## 📋 项目简介

Sol-Green 是一个去中心化的环保奖励平台，用户可以通过提交环保行为（垃圾分类、植树、低碳出行等）获得代币奖励。项目采用 AI 反欺诈检测、链上存证、第三方机构认证等技术，确保行为的真实性和可信度。

## 🏗️ 技术架构

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL / SQLite
- **缓存**: Redis
- **认证**: JWT

### 区块链
- **平台**: Solana
- **合约语言**: Rust (Anchor Framework)
- **网络**: Devnet / Mainnet

### 前端
- **框架**: React 18
- **钱包**: Solana Wallet Adapter
- **UI**: 自定义组件

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
# 安装依赖
go mod download

# 启动 PostgreSQL 和 Redis（如果使用）
docker-compose up -d postgres redis

# 运行后端
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
cd contracts/sol-green
anchor build

# 部署到测试网
anchor deploy --provider.cluster devnet
```

## 📁 项目结构

```
trends_solana/
├── config/              # 配置管理
│   ├── config.go
│   └── fraud_detector.go
├── controller/          # 控制器
│   ├── auth.go
│   └── green.go
├── middleware/          # 中间件
│   ├── auth.go
│   ├── cors.go
│   ├── logger.go
│   └── ratelimit.go
├── model/              # 数据模型
│   └── green.go
├── service/            # 业务逻辑
│   ├── solana.go
│   └── partner.go
├── contracts/          # Solana 智能合约
│   └── sol-green/
│       ├── src/
│       │   └── lib.rs
│       ├── Cargo.toml
│       └── Anchor.toml
├── frontend/           # React 前端
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   └── App.js
│   └── package.json
├── tests/              # 测试文件
├── scripts/            # 部署脚本
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 🔧 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DB_TYPE` | 数据库类型 (sqlite/postgres) | sqlite |
| `DATABASE_URL` | PostgreSQL 连接字符串 | - |
| `REDIS_URL` | Redis 地址 | localhost:6379 |
| `JWT_SECRET` | JWT 密钥 | - |
| `SOLANA_RPC_URL` | Solana RPC 地址 | https://api.devnet.solana.com |
| `SOLANA_ADMIN_PRIVKEY` | 管理员私钥 (Base58) | - |

### Solana 配置

1. 创建钱包并获取私钥：
```bash
solana-keygen new --outfile ~/.config/solana/admin.json
```

2. 配置环境变量中的 `SOLANA_ADMIN_PRIVKEY`

3. 获取测试网 SOL：
```bash
solana airdrop 1 <your-wallet-address> --url devnet
```

## 📡 API 文档

### 认证接口

#### 钱包登录
```
POST /api/v1/auth/wallet
Content-Type: application/json

{
  "wallet_addr": "string",
  "signature": "string",
  "message": "string"
}
```

### 环保行为接口

#### 提交环保行为
```
POST /api/v1/green/behavior/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "behavior_type": "waste_sorting|tree_planting|low_carbon_travel",
  "media_urls": ["string"],
  "location": "string"
}
```

#### 查询行为状态
```
GET /api/v1/green/behavior/:id
Authorization: Bearer <token>
```

#### 第三方认证
```
POST /api/v1/green/partner/verify
Authorization: Bearer <token>
Content-Type: application/json

{
  "behavior_id": "string",
  "partner_id": "string",
  "verify_code": "string"
}
```

## 🧪 测试

```bash
# 运行所有测试
./scripts/test.sh

# 运行 Go 测试
go test -v ./tests/...

# 运行前端测试
cd frontend && npm test
```

## 🚢 部署

### 使用 Docker Compose

```bash
./scripts/deploy.sh
```

### 手动部署

1. **后端部署**
   - 构建镜像: `docker build -t sol-green-backend .`
   - 运行容器: `docker run -d -p 8080:8080 sol-green-backend`

2. **前端部署**
   - 构建: `cd frontend && npm run build`
   - 部署到 Nginx 或静态托管服务

3. **合约部署**
   ```bash
   cd contracts/sol-green
   anchor deploy --provider.cluster mainnet
   ```

## 📚 文档

- [开发文档](./docs/DEVELOPMENT.md)
- [技术文档](./docs/TECHNICAL.md)
- [部署文档](./docs/DEPLOYMENT.md)
- [API 文档](./docs/API.md)

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
