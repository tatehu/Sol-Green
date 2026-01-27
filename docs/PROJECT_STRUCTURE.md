# 项目结构说明

## 📁 目录结构

```
trends_solana/
├── backend/                    # Go 后端服务
│   ├── config/                 # 配置管理
│   │   ├── config.go          # 主配置文件（数据库、Redis、日志）
│   │   └── fraud_detector.go  # AI 反欺诈检测器
│   ├── controller/            # 控制器层（处理 HTTP 请求）
│   │   ├── auth.go           # 认证控制器（钱包登录）
│   │   ├── green.go           # 环保行为控制器
│   │   └── health.go          # 健康检查
│   ├── middleware/            # 中间件
│   │   ├── auth.go           # JWT 认证中间件
│   │   ├── cors.go           # 跨域中间件
│   │   ├── logger.go         # 日志中间件
│   │   └── ratelimit.go      # 限流中间件
│   ├── model/                 # 数据模型
│   │   └── green.go          # 环保行为、用户、奖励记录模型
│   ├── service/               # 业务逻辑层
│   │   ├── solana.go         # Solana 区块链交互（奖励发放、存证）
│   │   └── partner.go        # 第三方认证服务
│   ├── tests/                 # 测试文件
│   │   ├── controller_test.go
│   │   └── service_test.go
│   ├── main.go                # 应用入口
│   ├── go.mod                 # Go 模块依赖
│   ├── go.sum                 # Go 依赖校验
│   └── Dockerfile             # 后端 Docker 镜像
│
├── frontend/                   # React 前端应用
│   ├── public/                # 静态文件
│   │   └── index.html
│   ├── src/
│   │   ├── components/        # React 组件
│   │   │   └── WalletConnect.jsx
│   │   ├── pages/            # 页面组件
│   │   │   ├── GreenSubmit.jsx
│   │   │   ├── GreenSubmit.css
│   │   │   ├── BehaviorStatus.jsx
│   │   │   └── BehaviorStatus.css
│   │   ├── App.js            # 主应用组件
│   │   ├── App.css
│   │   ├── index.js          # 入口文件
│   │   └── index.css
│   ├── package.json
│   ├── Dockerfile             # 前端 Docker 镜像
│   └── nginx.conf             # Nginx 配置
│
├── contracts/                 # Solana 智能合约
│   └── sol-green/             # 合约项目
│       ├── src/
│       │   └── lib.rs         # 合约主文件
│       ├── Cargo.toml        # Rust 依赖配置
│       ├── Anchor.toml        # Anchor 框架配置
│       └── Xargo.toml         # 交叉编译配置
│
├── docs/                      # 项目文档
│   ├── DEVELOPMENT.md         # 开发文档
│   ├── TECHNICAL.md           # 技术文档
│   ├── DEPLOYMENT.md          # 部署文档
│   ├── API.md                 # API 文档
│   ├── QUICKSTART.md          # 快速开始
│   ├── PROJECT_STRUCTURE.md   # 本文件
│   ├── CHECKLIST.md           # 项目验证清单
│   └── SUMMARY.md             # 项目总结
│
├── scripts/                   # 部署和工具脚本
│   ├── deploy.sh              # 部署脚本
│   ├── test.sh                # 测试脚本
│   ├── build-contract.sh      # 合约构建脚本
│   └── verify.sh              # 项目验证脚本
│
├── docker-compose.yml         # Docker Compose 配置
├── .env.example              # 环境变量示例
├── .gitignore                 # Git 忽略文件
└── README.md                  # 项目说明
```

## 📦 模块说明

### 后端模块 (backend/)

#### config/
- **config.go**: 初始化数据库、Redis、日志等全局配置
- **fraud_detector.go**: AI 反欺诈检测实现（当前为模拟，可扩展为真实 API）

#### controller/
- **auth.go**: 处理钱包登录认证
- **green.go**: 处理环保行为相关请求（提交、查询、第三方认证）
- **health.go**: 健康检查接口

#### middleware/
- **auth.go**: JWT Token 验证中间件
- **cors.go**: 跨域资源共享配置
- **logger.go**: HTTP 请求日志记录
- **ratelimit.go**: 基于 Redis 的请求限流

#### model/
- **green.go**: 数据模型定义（GreenBehavior, User, RewardRecord）

#### service/
- **solana.go**: Solana 区块链交互（代币转账、链上存证）
- **partner.go**: 第三方认证机构接口（阿拉善 SEE、政府环保部门）

### 智能合约 (contracts/sol-green/)

- **lib.rs**: Anchor 框架编写的 Solana 智能合约
  - `initialize_reward_config`: 初始化奖励配置
  - `mint_green_reward`: 发放环保奖励
  - `record_proof`: 链上存证

### 前端模块 (frontend/)

- **components/WalletConnect.jsx**: Solana 钱包连接组件
- **pages/GreenSubmit.jsx**: 环保行为提交页面
- **pages/BehaviorStatus.jsx**: 行为状态查询页面
- **App.js**: 主应用组件（路由配置）

### 测试 (backend/tests/)

- **controller_test.go**: 控制器单元测试
- **service_test.go**: 服务层单元测试

### 部署 (scripts/)

- **deploy.sh**: 一键部署脚本（Docker Compose）
- **test.sh**: 运行所有测试
- **build-contract.sh**: 构建 Solana 智能合约
- **verify.sh**: 验证项目完整性

## 🔄 数据流

### 用户提交环保行为流程

```
前端 (frontend/src/pages/GreenSubmit.jsx)
  ↓ HTTP POST
后端控制器 (backend/controller/green.go)
  ↓
中间件验证 (backend/middleware/auth.go)
  ↓
业务逻辑 (backend/service/solana.go)
  ↓
AI 检测 (backend/config/fraud_detector.go)
  ↓
数据库存储 (backend/model/green.go)
  ↓
区块链交互 (backend/service/solana.go)
  ↓
智能合约 (contracts/sol-green/src/lib.rs)
  ↓
返回结果
```

## 📝 关键文件说明

### backend/main.go
应用入口，初始化所有服务并启动 HTTP 服务器

### backend/go.mod
Go 模块依赖定义

### .env.example
环境变量配置模板，包含数据库、Redis、Solana 等配置

### docker-compose.yml
Docker Compose 配置，定义所有服务（PostgreSQL、Redis、后端）

### backend/Dockerfile
后端应用的 Docker 镜像构建文件

### frontend/Dockerfile
前端应用的 Docker 镜像构建文件（多阶段构建）

## 🔧 扩展指南

### 添加新的 API 接口

1. 在 `backend/controller/` 创建新的控制器文件
2. 在 `backend/main.go` 注册路由
3. 在 `docs/API.md` 添加文档

### 添加新的数据模型

1. 在 `backend/model/` 定义模型结构
2. 在 `backend/config/config.go` 的 `InitDB()` 中添加自动迁移

### 添加新的智能合约功能

1. 在 `contracts/sol-green/src/lib.rs` 添加新函数
2. 更新账户结构和事件
3. 在 `backend/service/solana.go` 添加调用逻辑

### 添加新的前端页面

1. 在 `frontend/src/pages/` 创建页面组件
2. 在 `frontend/src/App.js` 添加路由
3. 更新导航菜单

## 🚀 运行方式

### 本地开发

```bash
# 后端
cd backend
go run main.go

# 前端
cd frontend
npm start
```

### Docker 部署

```bash
# 从项目根目录
docker-compose up -d
```

### 验证项目

```bash
# 运行验证脚本
./scripts/verify.sh
```

## 📊 项目统计

- **后端代码**: Go 语言，约 15 个文件
- **前端代码**: React + JavaScript，约 10 个文件
- **智能合约**: Rust + Anchor，1 个主文件
- **测试文件**: 2 个测试文件
- **文档**: 8 个文档文件
- **脚本**: 4 个部署/工具脚本
