# 技术文档 v3.0.0 / Technical Documentation v3.0.0

## 系统架构 / System Architecture

### 整体架构 / Overall Architecture

```
┌─────────────┐
│  前端 (React) │ / Frontend (React)
└──────┬──────┘
       │ HTTP/HTTPS
       ▼
┌─────────────┐
│ Go 后端 API  │ / Go Backend API
└──────┬──────┘
       │
       ├──► PostgreSQL ── 数据存储 / Data Storage
       ├──► Redis ─────── 缓存/限流 / Cache/Rate Limiting
       └──► Solana ────── 区块链交互 / Blockchain Interaction
       └──► 多链支持 ──── 多链交互 / Multi-Chain Support
```

### 技术栈详解 / Technology Stack Details

#### 后端技术栈 / Backend Technology Stack

- **Gin**: 高性能 HTTP Web 框架 / High-performance HTTP web framework
- **GORM**: ORM 框架，支持多种数据库 / ORM framework supporting multiple databases
- **Redis**: 缓存和限流 / Caching and rate limiting
- **JWT**: 无状态认证 / Stateless authentication
- **Solana Go SDK**: Solana 区块链交互 / Solana blockchain interaction

#### 前端技术栈 / Frontend Technology Stack

- **React 18**: UI 框架 / UI framework
- **Solana Wallet Adapter**: 钱包连接 / Wallet connection
- **Axios**: HTTP 客户端 / HTTP client
- **React Router**: 路由管理 / Route management

#### 区块链技术栈 / Blockchain Technology Stack

- **Anchor**: Solana 智能合约框架 / Solana smart contract framework
- **Rust**: 合约开发语言 / Contract development language
- **SPL Token**: 代币标准 / Token standard

#### 运维技术栈 / Operations Technology Stack

- **Docker**: 容器化部署 / Containerized deployment
- **WebSocket**: 实时数据监听 / Real-time data monitoring
- **监控告警**: 系统监控 / System monitoring
- **数据分析**: 运营数据分析 / Operational data analysis

## 核心功能实现 / Core Feature Implementation

### 1. 钱包认证 / Wallet Authentication

#### 流程 / Process

1. 用户在前端连接 Solana 钱包 / User connects Solana wallet on frontend
2. 前端生成签名消息 / Frontend generates signature message
3. 用户使用钱包签名 / User signs with wallet
4. 后端验证签名 / Backend verifies signature
5. 生成 JWT Token 返回 / Generate JWT token and return

#### 代码位置 / Code Location

- 前端: `frontend/src/components/WalletConnect.jsx`
- 后端: `backend/controller/auth.go`

### 2. 环保行为提交 / Environmental Behavior Submission

#### 流程 / Process

1. 用户提交行为信息 / User submits behavior information
2. 后端验证参数 / Backend validates parameters
3. Redis 防重复提交 / Redis prevents duplicate submission
4. AI 反欺诈检测 / AI fraud detection
5. 根据风险等级处理 / Process based on risk level
6. 发放奖励 / Distribute rewards
7. 链上存证 / On-chain proof storage

#### 关键技术 / Key Technologies

- **AI 反欺诈检测**: 集成百度 AI、阿里云、AWS、Google Cloud
- **链上存证**: SHA256 哈希永久记录在 Solana 链上
- **多链支持**: 支持 Solana、Ethereum、Polygon、BSC

### 3. 挑战活动系统 / Challenge Activity System

#### 功能特点 / Features

- **创建挑战**: 用户发起个人/团队/社区挑战
- **参与挑战**: 提交行为自动参与
- **进度追踪**: 实时显示完成进度
- **奖励发放**: 挑战完成自动发放奖励

#### 数据模型 / Data Models

```go
type Challenge struct {
    ID              string
    Title           string
    ChallengeType   string  // individual/team/community
    BehaviorType    string  // waste_sorting/tree_planting/etc
    RewardAmount    uint64
    TargetCount     int
    CurrentCount    int
    Status          string  // draft/active/completed/cancelled
}
```

### 4. 营销活动系统 / Marketing Activity System

#### 活动类型 / Activity Types

1. **签到活动 / Sign-in Activities**: 连续签到获得奖励
2. **邀请活动 / Invite Activities**: 邀请好友获得奖励
3. **每日任务 / Daily Tasks**: 完成任务获得奖励
4. **抽奖活动 / Lucky Draws**: 参与抽奖获得随机奖励
5. **限时抢购 / Flash Sales**: 限时高额奖励
6. **节日活动 / Festival Activities**: 主题化节日活动

### 5. 链上数据监听 / Blockchain Data Monitoring

#### 实现方式 / Implementation

- **WebSocket 订阅**: 实时监听 Solana 账户变化
- **轮询模式**: 备用数据获取方式
- **事件过滤**: 只监听相关交易事件
- **数据存储**: 结构化存储到数据库

#### 监听内容 / Monitoring Content

- **交易事件**: 奖励发放、存证交易
- **账户变化**: 用户余额、合约状态
- **程序日志**: 智能合约执行日志

### 6. 数据分析服务 / Data Analytics Service

#### 分析维度 / Analysis Dimensions

- **用户分析**: 活跃用户、新增用户、留存率
- **行为分析**: 行为类型分布、时间趋势
- **奖励分析**: 奖励发放统计、成本分析
- **活动分析**: 参与度、转化率、ROI

#### 报告类型 / Report Types

- **日报 / Daily Reports**: 当日数据汇总
- **周报 / Weekly Reports**: 周趋势分析
- **月报 / Monthly Reports**: 月度总结
- **自定义报告 / Custom Reports**: 按需分析

## 数据库设计 / Database Design

### 核心数据表 / Core Data Tables

#### 用户相关 / User Related
- `users`: 用户基本信息
- `user_wallets`: 钱包地址映射

#### 行为相关 / Behavior Related
- `green_behaviors`: 环保行为记录
- `behavior_rewards`: 奖励发放记录

#### 活动相关 / Activity Related
- `challenges`: 挑战活动
- `challenge_participants`: 挑战参与者
- `marketing_activities`: 营销活动
- `marketing_participants`: 营销活动参与者

#### 系统相关 / System Related
- `sign_in_records`: 签到记录
- `invite_records`: 邀请记录
- `analytics_data`: 数据分析缓存

### 索引设计 / Index Design

```sql
-- 用户钱包地址索引
CREATE INDEX idx_users_wallet_addr ON users(wallet_addr);

-- 行为类型和时间索引
CREATE INDEX idx_behaviors_type_time ON green_behaviors(behavior_type, submit_time);

-- 挑战状态索引
CREATE INDEX idx_challenges_status ON challenges(status);

-- 时间范围查询索引
CREATE INDEX idx_rewards_created_at ON behavior_rewards(created_at);
```

## 性能优化 / Performance Optimization

### 后端优化 / Backend Optimization

#### 数据库优化 / Database Optimization
- **索引策略**: 为常用查询字段添加索引
- **连接池**: 使用 GORM 连接池管理
- **读写分离**: 支持主从数据库配置

#### 缓存策略 / Caching Strategy
- **Redis 缓存**: 热点数据缓存
- **应用缓存**: 本地内存缓存
- **CDN 加速**: 静态资源加速

#### 异步处理 / Asynchronous Processing
- **消息队列**: 使用 Redis 作为消息队列
- **后台任务**: 区块链交易异步处理
- **定时任务**: 数据分析定时执行

### 前端优化 / Frontend Optimization

#### 代码分割 / Code Splitting
- **路由分割**: 按页面分割代码
- **组件懒加载**: 动态导入组件
- **第三方库分割**: 单独打包大型库

#### 资源优化 / Resource Optimization
- **图片压缩**: WebP 格式和压缩
- **字体优化**: 使用系统字体
- **Bundle 分析**: 分析包大小并优化

### 区块链优化 / Blockchain Optimization

#### 交易优化 / Transaction Optimization
- **批量交易**: 合并多个操作
- **Gas 优化**: 合理设置交易参数
- **网络选择**: 根据成本选择网络

#### 合约优化 / Contract Optimization
- **存储优化**: 减少链上存储成本
- **计算优化**: 优化合约逻辑复杂度
- **事件优化**: 合理使用事件日志

## 安全机制 / Security Mechanisms

### 认证安全 / Authentication Security

- **JWT Token**: 无状态认证，设置过期时间
- **签名验证**: Solana 钱包签名验证
- **权限控制**: 基于角色的访问控制

### 数据安全 / Data Security

- **密码加密**: 用户敏感信息加密存储
- **传输加密**: HTTPS 传输加密
- **API 限流**: 防止恶意请求

### 区块链安全 / Blockchain Security

- **私钥管理**: 环境变量存储，不硬编码
- **交易验证**: 多重签名验证
- **合约审计**: 专业审计机构审计

## 监控告警 / Monitoring and Alerting

### 系统监控 / System Monitoring

- **应用监控**: CPU、内存、磁盘使用率
- **服务监控**: API 响应时间、错误率
- **数据库监控**: 连接数、慢查询

### 业务监控 / Business Monitoring

- **用户监控**: 活跃用户、注册转化
- **交易监控**: 交易成功率、确认时间
- **活动监控**: 活动参与度、奖励发放

### 告警机制 / Alerting Mechanism

- **阈值告警**: 超过设定阈值自动告警
- **异常检测**: 自动检测异常模式
- **通知渠道**: 邮件、短信、Webhook

## 扩展性设计 / Scalability Design

### 水平扩展 / Horizontal Scaling

- **无状态设计**: 后端服务无状态，支持多实例
- **负载均衡**: Nginx 负载均衡
- **数据库分库分表**: 支持数据分片

### 垂直扩展 / Vertical Scaling

- **缓存层**: 多级缓存架构
- **异步处理**: 消息队列异步处理
- **CDN 加速**: 全球内容分发

### 功能扩展 / Feature Extension

- **插件架构**: 支持功能插件化
- **API 版本控制**: 支持多版本 API
- **配置热更新**: 支持配置动态更新

## 部署架构 / Deployment Architecture

### 生产环境架构 / Production Architecture

```
┌─────────────────┐
│   Load Balancer │  Nginx
│   (Nginx)       │
└─────────┬───────┘
          │
    ┌─────┴─────┐
    │  API Gateway │  路由、认证、限流
    └─────┬─────┘
          │
    ┌─────┼─────┐
    │     │     │
┌───▼──┐  ▼  ┌───▼──┐
│ Backend │    │ Frontend │
│ Service │    │ Service │
└───┬───┘    └────┬───┘
    │            │
    └─────┬──────┘
          │
    ┌─────┴─────┐
    │  Database │  PostgreSQL + Redis
    │  Cache    │
    └─────┬─────┘
          │
    ┌─────┴─────┐
    │ Blockchain │  Solana + 多链
    │  Networks │
    └───────────┘
```

### 容器化部署 / Containerized Deployment

- **Docker**: 应用容器化
- **Kubernetes**: 容器编排
- **Helm**: 应用包管理
- **CI/CD**: 自动化部署管道

## 故障处理 / Fault Handling

### 常见故障 / Common Faults

#### 数据库连接失败 / Database Connection Failed
- **原因**: 网络问题、配置错误、资源不足
- **处理**: 检查配置、重启服务、扩容资源

#### Redis 缓存失效 / Redis Cache Failure
- **原因**: 内存不足、服务宕机
- **处理**: 扩容内存、集群部署、主从切换

#### 区块链网络异常 / Blockchain Network Exception
- **原因**: 网络拥堵、RPC 节点问题
- **处理**: 切换节点、重试机制、多节点备份

### 回滚策略 / Rollback Strategy

1. **代码回滚**: Git 版本回滚
2. **数据库回滚**: 备份恢复
3. **合约回滚**: 部署旧版本合约

### 应急预案 / Emergency Plan

1. **监控告警**: 24/7 监控
2. **备份策略**: 多重备份
3. **恢复流程**: 详细的恢复步骤
4. **沟通机制**: 及时通知用户

---

**文档版本 / Document Version**: v3.0.0  
**最后更新 / Last Update**: 2024-01-24
