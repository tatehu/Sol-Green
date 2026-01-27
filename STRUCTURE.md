# Sol-Green 项目结构

## 📁 完整目录树

```
trends_solana/
├── backend/                    # Go 后端服务
│   ├── config/                # 配置管理
│   ├── controller/            # 控制器
│   ├── middleware/            # 中间件
│   ├── model/                 # 数据模型
│   ├── service/               # 业务逻辑
│   ├── tests/                 # 测试文件
│   ├── main.go                # 应用入口
│   ├── go.mod                 # Go 依赖
│   ├── go.sum                 # 依赖校验
│   └── Dockerfile             # 后端镜像
│
├── frontend/                   # React 前端
│   ├── src/                   # 源代码
│   │   ├── components/        # 组件
│   │   ├── pages/            # 页面
│   │   ├── App.js
│   │   └── index.js
│   ├── public/                # 静态文件
│   ├── package.json
│   ├── Dockerfile
│   └── nginx.conf
│
├── contracts/                  # Solana 智能合约
│   └── sol-green/             # 合约项目
│       ├── src/
│       │   └── lib.rs
│       ├── Cargo.toml
│       ├── Anchor.toml
│       └── Xargo.toml
│
├── docs/                       # 项目文档
│   ├── DEVELOPMENT.md
│   ├── TECHNICAL.md
│   ├── DEPLOYMENT.md
│   ├── API.md
│   ├── QUICKSTART.md
│   ├── PROJECT_STRUCTURE.md
│   ├── CHECKLIST.md
│   └── SUMMARY.md
│
├── scripts/                    # 脚本工具
│   ├── deploy.sh
│   ├── test.sh
│   ├── build-contract.sh
│   └── verify.sh
│
├── docker-compose.yml          # Docker 编排
├── .env.example                # 环境变量示例
├── .gitignore                  # Git 忽略
└── README.md                   # 项目说明
```

## 🎯 项目组织原则

1. **前后端分离**: 后端在 `backend/`，前端在 `frontend/`
2. **合约独立**: 智能合约在 `contracts/`
3. **文档集中**: 所有文档在 `docs/`
4. **脚本工具**: 部署和工具脚本在 `scripts/`

## 📊 统计信息

- **后端 Go 文件**: 15 个
- **前端 React 文件**: 10+ 个
- **智能合约 Rust 文件**: 1 个主文件
- **测试文件**: 2 个
- **文档文件**: 8 个
- **脚本文件**: 4 个

## ✅ 验证状态

运行 `./scripts/verify.sh` 已验证：
- ✅ 目录结构完整
- ✅ 关键文件存在
- ✅ 脚本可执行
- ✅ Docker 配置有效

## 🚀 快速启动

```bash
# 1. 验证项目
./scripts/verify.sh

# 2. 配置环境
cp .env.example .env

# 3. 启动服务
docker-compose up -d

# 4. 启动前端（开发模式）
cd frontend && npm install && npm start
```
