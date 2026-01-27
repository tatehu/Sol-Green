# 开发文档

## 开发环境搭建

### 1. 安装依赖

#### Go 依赖
```bash
go mod download
```

#### Node.js 依赖
```bash
cd frontend
npm install
```

#### Solana 工具链
```bash
# 安装 Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/stable/install)"

# 安装 Anchor CLI
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install latest
avm use latest
```

### 2. 配置开发环境

1. 复制环境变量文件：
```bash
cp .env.example .env
```

2. 配置数据库（开发环境可使用 SQLite）：
```bash
DB_TYPE=sqlite
```

3. 配置 Solana 测试网：
```bash
solana config set --url devnet
```

### 3. 启动开发服务

#### 后端
```bash
go run main.go
```

#### 前端
```bash
cd frontend
npm start
```

## 代码规范

### Go 代码规范

1. 使用 `gofmt` 格式化代码
2. 遵循 Go 官方代码规范
3. 使用有意义的变量和函数名
4. 添加必要的注释

```bash
# 格式化代码
go fmt ./...

# 运行 linter
golangci-lint run
```

### React 代码规范

1. 使用 ESLint 检查代码
2. 遵循 React Hooks 最佳实践
3. 组件命名使用 PascalCase
4. 文件命名使用 PascalCase（组件）或 camelCase（工具函数）

```bash
# 检查代码
cd frontend
npm run lint
```

## 开发流程

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 开发功能

- 编写代码
- 添加测试
- 更新文档

### 3. 提交代码

```bash
git add .
git commit -m "feat: 添加新功能"
```

### 4. 创建 Pull Request

在 GitHub 上创建 PR，等待代码审查。

## 测试

### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./controller

# 显示测试覆盖率
go test -cover ./...
```

### 集成测试

```bash
# 启动测试环境
docker-compose -f docker-compose.test.yml up -d

# 运行集成测试
go test -tags=integration ./tests/...
```

### E2E 测试

```bash
cd frontend
npm run test:e2e
```

## 调试

### 后端调试

1. 使用 `logrus` 记录日志
2. 使用 Go 调试器：
```bash
dlv debug main.go
```

### 前端调试

1. 使用 React DevTools
2. 使用浏览器开发者工具
3. 查看控制台日志

### Solana 合约调试

```bash
# 启动本地验证器
solana-test-validator

# 部署到本地
anchor deploy --provider.cluster localnet

# 查看日志
solana logs
```

## 常见问题

### 1. 数据库连接失败

- 检查数据库服务是否运行
- 验证连接字符串是否正确
- 检查防火墙设置

### 2. Redis 连接失败

- 检查 Redis 服务是否运行
- 验证地址和端口
- 开发环境可跳过 Redis（代码已处理）

### 3. Solana 交易失败

- 检查钱包余额是否充足
- 验证 RPC 节点是否可用
- 检查网络配置（devnet/mainnet）

### 4. 前端无法连接后端

- 检查 CORS 配置
- 验证 API 地址是否正确
- 检查后端服务是否运行

## 性能优化

### 后端优化

1. 使用 Redis 缓存热点数据
2. 数据库查询优化（索引、分页）
3. 使用连接池
4. 异步处理耗时操作

### 前端优化

1. 代码分割和懒加载
2. 图片优化和 CDN
3. 使用 React.memo 减少重渲染
4. 虚拟滚动（长列表）

## 代码审查清单

- [ ] 代码符合项目规范
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
- [ ] 没有硬编码的敏感信息
- [ ] 错误处理完善
- [ ] 日志记录适当
- [ ] 性能考虑合理
