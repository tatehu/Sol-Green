# 技术文档

## 系统架构

### 整体架构

```
┌─────────────┐
│  前端 (React) │
└──────┬──────┘
       │ HTTP/HTTPS
       ▼
┌─────────────┐
│ Go 后端 API  │
└──────┬──────┘
       │
       ├──► PostgreSQL ── 数据存储
       ├──► Redis ─────── 缓存/限流
       └──► Solana ────── 区块链交互
```

### 技术栈详解

#### 后端技术栈

- **Gin**: 高性能 HTTP Web 框架
- **GORM**: ORM 框架，支持多种数据库
- **Redis**: 缓存和限流
- **JWT**: 无状态认证
- **Solana Go SDK**: Solana 区块链交互

#### 前端技术栈

- **React 18**: UI 框架
- **Solana Wallet Adapter**: 钱包连接
- **Axios**: HTTP 客户端
- **React Router**: 路由管理

#### 区块链技术栈

- **Anchor**: Solana 智能合约框架
- **Rust**: 合约开发语言
- **SPL Token**: 代币标准

## 核心功能实现

### 1. 钱包认证

#### 流程

1. 用户在前端连接 Solana 钱包
2. 前端生成签名消息
3. 用户使用钱包签名
4. 后端验证签名
5. 生成 JWT Token 返回前端

#### 代码位置

- 前端: `frontend/src/components/WalletConnect.jsx`
- 后端: `controller/auth.go`

### 2. 环保行为提交

#### 流程

1. 用户提交行为信息（类型、媒体文件、地点）
2. 后端验证请求参数
3. Redis 防重复提交检查
4. AI 反欺诈检测
5. 根据检测结果：
   - 低风险：自动通过，发放奖励
   - 高风险：进入人工审核
6. 链上存证
7. 返回结果

#### 代码位置

- 前端: `frontend/src/pages/GreenSubmit.jsx`
- 后端: `controller/green.go`
- 服务: `service/solana.go`

### 3. AI 反欺诈检测

#### 实现方式

当前为模拟实现，实际应集成：
- 百度 AI 图像审核
- 阿里云内容安全
- 深度伪造检测模型

#### 代码位置

- `config/fraud_detector.go`

### 4. 链上奖励发放

#### 流程

1. 后端调用 Solana 合约
2. 合约验证管理员权限
3. 执行代币转账
4. 记录奖励事件
5. 返回交易哈希

#### 代码位置

- 合约: `contracts/sol-green/src/lib.rs`
- 后端: `service/solana.go`

### 5. 链上存证

#### 流程

1. 生成行为哈希（SHA256）
2. 调用存证合约
3. 将哈希写入链上账户
4. 返回存证交易哈希

#### 代码位置

- 合约: `contracts/sol-green/src/lib.rs` (record_proof)
- 后端: `service/solana.go` (RecordOnChain)

## 数据库设计

### GreenBehavior 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键 |
| wallet_addr | VARCHAR(44) | 钱包地址 |
| behavior_type | VARCHAR(20) | 行为类型 |
| media_urls | JSON | 媒体文件 URL |
| location | VARCHAR(100) | 地点 |
| status | VARCHAR(20) | 状态 |
| fraud_score | DECIMAL(3,2) | 欺诈分数 |
| submit_time | TIMESTAMP | 提交时间 |
| approve_time | TIMESTAMP | 通过时间 |
| tx_hash | VARCHAR(90) | 交易哈希 |
| proof_hash | VARCHAR(64) | 存证哈希 |

### User 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键 |
| wallet_addr | VARCHAR(44) | 钱包地址（唯一） |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### RewardRecord 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | 主键 |
| wallet_addr | VARCHAR(44) | 钱包地址 |
| behavior_id | VARCHAR(36) | 行为 ID |
| amount | BIGINT | 奖励数量 |
| tx_hash | VARCHAR(90) | 交易哈希（唯一） |
| created_at | TIMESTAMP | 创建时间 |

## 智能合约设计

### 合约结构

```rust
pub mod sol_green {
    // 初始化奖励配置
    pub fn initialize_reward_config(...)
    
    // 发放环保奖励
    pub fn mint_green_reward(...)
    
    // 链上存证
    pub fn record_proof(...)
}
```

### 账户结构

- **RewardConfig**: 奖励配置账户
- **ProofRecord**: 存证记录账户

### 事件

- **RewardMinted**: 奖励发放事件
- **ProofRecorded**: 存证记录事件

## 安全机制

### 1. 认证与授权

- JWT Token 认证
- Solana 钱包签名验证
- 管理员权限校验

### 2. 防刷机制

- Redis 限流（IP 级别）
- 防重复提交（用户+行为类型+时间窗口）
- 请求频率限制

### 3. 数据安全

- 敏感信息加密存储
- SQL 注入防护（GORM）
- XSS 防护（前端转义）

### 4. 区块链安全

- 私钥安全存储（环境变量）
- 交易签名验证
- 合约权限控制

## 性能优化

### 后端优化

1. **数据库优化**
   - 添加索引（wallet_addr, behavior_type, status）
   - 使用连接池
   - 查询优化

2. **缓存策略**
   - Redis 缓存热点数据
   - 缓存过期时间设置

3. **异步处理**
   - 文件上传异步处理
   - 链上操作异步执行

### 前端优化

1. **代码分割**
   - 路由级别代码分割
   - 组件懒加载

2. **资源优化**
   - 图片压缩
   - CDN 加速

3. **渲染优化**
   - React.memo
   - useMemo/useCallback

## 监控与日志

### 日志系统

- 使用 `logrus` 记录结构化日志
- 日志级别：DEBUG, INFO, WARN, ERROR
- 日志输出：控制台 + 文件

### 监控指标

- API 响应时间
- 错误率
- 数据库连接数
- Redis 命中率
- Solana 交易成功率

## 扩展性设计

### 水平扩展

- 无状态后端设计
- 负载均衡支持
- 数据库读写分离

### 功能扩展

- 插件化 AI 检测
- 多链支持
- 多语言支持

## 故障处理

### 常见故障

1. **数据库连接失败**
   - 自动重试机制
   - 降级方案（使用缓存）

2. **Redis 不可用**
   - 降级到内存缓存
   - 限流功能降级

3. **Solana 网络问题**
   - 交易重试机制
   - 异步队列处理

### 容错机制

- 优雅降级
- 熔断器模式
- 超时控制
