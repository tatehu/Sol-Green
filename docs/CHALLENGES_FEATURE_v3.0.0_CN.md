# 挑战活动功能 v3.0.0 / Challenge Activity Features v3.0.0

## 📋 功能概述 / Feature Overview

挑战活动功能允许用户发起和参与环保挑战，通过社交化和游戏化的方式提高用户参与度，增强平台的活跃度和影响力。  
The challenge activity feature allows users to initiate and participate in environmental challenges, increasing user engagement and platform vitality through social and gamified approaches.

## 🎯 功能特点 / Feature Characteristics

### 1. 挑战类型 / Challenge Types

- **个人挑战 / Individual Challenge**: 个人设定的环保目标 / Personal environmental goals
  - 示例: "本月垃圾分类 30 次" / Example: "Sort waste 30 times this month"
  - 适用场景: 个人习惯养成 / Suitable for: Personal habit formation

- **团队挑战 / Team Challenge**: 团队协作完成环保任务 / Team collaboration on environmental tasks
  - 示例: "团队植树 100 棵" / Example: "Team plants 100 trees"
  - 适用场景: 企业/社区活动 / Suitable for: Corporate/community activities

- **社区挑战 / Community Challenge**: 社区级别的环保活动 / Community-level environmental activities
  - 示例: "社区低碳出行周" / Example: "Community low-carbon travel week"
  - 适用场景: 社区共建 / Suitable for: Community co-construction

### 2. 挑战流程 / Challenge Process

```
发起挑战 → 激活挑战 → 用户参与 → 达成目标 → 发放奖励
Create Challenge → Activate Challenge → User Participation → Achieve Goal → Distribute Rewards
```

#### 详细流程 / Detailed Process

1. **挑战创建 / Challenge Creation**
   - 用户填写挑战信息 / User fills in challenge information
   - 设置奖励和参与条件 / Set rewards and participation conditions
   - 挑战初始状态为"草稿" / Challenge initial status is "draft"

2. **挑战激活 / Challenge Activation**
   - 创建者或管理员激活挑战 / Creator or admin activates challenge
   - 挑战状态变为"进行中" / Challenge status becomes "active"
   - 开始接受用户参与 / Start accepting user participation

3. **用户参与 / User Participation**
   - 用户浏览并选择挑战 / Users browse and select challenges
   - 提交符合条件的环保行为 / Submit qualified environmental behaviors
   - 系统自动关联行为和挑战 / System automatically associates behavior with challenge

4. **目标达成 / Goal Achievement**
   - 实时统计参与人数 / Real-time statistics of participants
   - 进度条显示完成情况 / Progress bar shows completion status
   - 达成目标时自动标记完成 / Automatically mark as completed when goal achieved

5. **奖励发放 / Reward Distribution**
   - 挑战完成后参与者领取奖励 / Participants claim rewards after challenge completion
   - 支持多种奖励类型 / Support multiple reward types
   - 链上奖励记录 / On-chain reward records

### 3. 奖励机制 / Reward Mechanism

#### 基础奖励 / Basic Rewards
- 根据行为类型发放基础奖励 / Distribute basic rewards based on behavior type
- 参考标准奖励配置 / Reference standard reward configuration

#### 挑战奖励 / Challenge Rewards
- 挑战完成后的额外奖励 / Additional rewards after challenge completion
- 奖励倍率可配置（默认 1.5 倍）/ Reward multiplier configurable (default 1.5x)

#### 奖励计算公式 / Reward Calculation Formula

```go
// 挑战奖励计算 / Challenge reward calculation
challengeReward = baseReward * challengeBonusRate
totalReward = baseReward + challengeReward

// 示例 / Example
// 垃圾分类基础奖励: 1000 SOLGREEN
// 挑战奖励: 1000 * 1.5 = 1500 SOLGREEN
// 总奖励: 1000 + 1500 = 2500 SOLGREEN
```

## 📡 API 接口 / API Interface

### 创建挑战 / Create Challenge

**POST** `/api/v1/challenges`

创建新的环保挑战活动。  
Create a new environmental challenge activity.

#### 请求参数 / Request Parameters

```json
{
  "title": "垃圾分类挑战周", // "Waste Sorting Challenge Week"
  "description": "本周完成 10 次垃圾分类", // "Complete 10 waste sorting this week"
  "challenge_type": "individual", // 挑战类型 / Challenge type: individual | team | community
  "behavior_type": "waste_sorting", // 关联行为类型 / Associated behavior type
  "reward_amount": 5000, // 挑战奖励 / Challenge reward amount
  "target_count": 50, // 目标参与人数 / Target participant count
  "start_time": "2024-01-25T00:00:00Z", // 开始时间 / Start time
  "end_time": "2024-02-01T23:59:59Z", // 结束时间 / End time
  "rules": "每次垃圾分类需上传照片证明" // 挑战规则 / Challenge rules
}
```

#### 响应示例 / Response Example

```json
{
  "msg": "挑战创建成功", // "Challenge created successfully"
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "垃圾分类挑战周",
    "status": "draft",
    "current_count": 0,
    "target_count": 50,
    "created_at": "2024-01-24T10:00:00Z"
  }
}
```

### 获取挑战列表 / Get Challenges List

**GET** `/api/v1/challenges?status=active&challenge_type=individual&page=1&page_size=20`

获取挑战活动列表。  
Get challenge activities list.

#### 查询参数 / Query Parameters

| 参数 / Parameter | 类型 / Type | 说明 / Description | 默认值 / Default |
|----------------|------------|-------------------|---------------|
| `status` | string | 挑战状态 / Challenge status | - |
| `challenge_type` | string | 挑战类型 / Challenge type | - |
| `page` | int | 页码 / Page number | 1 |
| `page_size` | int | 每页数量 / Page size | 20 |

#### 响应示例 / Response Example

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "垃圾分类挑战周",
      "challenge_type": "individual",
      "behavior_type": "waste_sorting",
      "reward_amount": 5000,
      "current_count": 25,
      "target_count": 50,
      "status": "active",
      "start_time": "2024-01-25T00:00:00Z",
      "end_time": "2024-02-01T23:59:59Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 20
}
```

### 获取挑战详情 / Get Challenge Details

**GET** `/api/v1/challenges/:id`

获取指定挑战的详细信息。  
Get detailed information of specified challenge.

#### 响应示例 / Response Example

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "垃圾分类挑战周",
    "description": "本周完成 10 次垃圾分类",
    "challenge_type": "individual",
    "behavior_type": "waste_sorting",
    "reward_amount": 5000,
    "current_count": 25,
    "target_count": 50,
    "status": "active",
    "start_time": "2024-01-25T00:00:00Z",
    "end_time": "2024-02-01T23:59:59Z",
    "rules": "每次垃圾分类需上传照片证明"
  },
  "participant_count": 25, // 参与者总数 / Total participants
  "participants": [ // 最近参与者 / Recent participants
    {
      "wallet_addr": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
      "joined_at": "2024-01-25T08:00:00Z"
    }
  ]
}
```

### 参与挑战 / Join Challenge

**POST** `/api/v1/challenges/:id/join`

使用已审核通过的环保行为参与挑战。  
Join challenge using approved environmental behavior.

#### 请求参数 / Request Parameters

```json
{
  "behavior_id": "550e8400-e29b-41d4-a716-446655440001" // 关联的环保行为ID / Associated environmental behavior ID
}
```

#### 响应示例 / Response Example

```json
{
  "msg": "成功参与挑战", // "Successfully joined challenge"
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "challenge_id": "550e8400-e29b-41d4-a716-446655440000",
    "wallet_addr": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
    "joined_at": "2024-01-25T08:00:00Z"
  }
}
```

### 激活挑战 / Activate Challenge

**POST** `/api/v1/challenges/:id/activate`

激活挑战（仅创建者或管理员）。  
Activate challenge (creator or admin only).

#### 响应示例 / Response Example

```json
{
  "msg": "挑战已激活", // "Challenge activated"
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "active",
    "updated_at": "2024-01-25T00:00:00Z"
  }
}
```

### 领取挑战奖励 / Claim Challenge Reward

**POST** `/api/v1/challenges/:id/claim`

挑战完成后领取奖励。  
Claim reward after challenge completion.

#### 响应示例 / Response Example

```json
{
  "msg": "奖励领取成功", // "Reward claimed successfully"
  "tx_hash": "5j7s8K9..."
}
```

## 📊 数据模型 / Data Models

### Challenge 表 / Challenge Table

| 字段 / Field | 类型 / Type | 说明 / Description |
|-------------|------------|-------------------|
| `id` | VARCHAR(36) | 主键 / Primary key |
| `title` | VARCHAR(200) | 挑战标题 / Challenge title |
| `description` | TEXT | 挑战描述 / Challenge description |
| `challenge_type` | VARCHAR(20) | 挑战类型 / Challenge type |
| `behavior_type` | VARCHAR(20) | 关联行为类型 / Associated behavior type |
| `reward_amount` | BIGINT | 奖励数量 / Reward amount |
| `target_count` | INT | 目标参与人数 / Target participant count |
| `current_count` | INT | 当前参与人数 / Current participant count |
| `start_time` | TIMESTAMP | 开始时间 / Start time |
| `end_time` | TIMESTAMP | 结束时间 / End time |
| `status` | VARCHAR(20) | 状态 / Status |
| `rules` | TEXT | 挑战规则 / Challenge rules |
| `created_by` | VARCHAR(44) | 创建者钱包 / Creator wallet |
| `created_at` | TIMESTAMP | 创建时间 / Created time |
| `updated_at` | TIMESTAMP | 更新时间 / Updated time |

### ChallengeParticipant 表 / Challenge Participant Table

| 字段 / Field | 类型 / Type | 说明 / Description |
|-------------|------------|-------------------|
| `id` | VARCHAR(36) | 主键 / Primary key |
| `challenge_id` | VARCHAR(36) | 挑战ID / Challenge ID |
| `wallet_addr` | VARCHAR(44) | 钱包地址 / Wallet address |
| `behavior_id` | VARCHAR(36) | 关联行为ID / Associated behavior ID |
| `joined_at` | TIMESTAMP | 参与时间 / Join time |
| `reward_claimed` | BOOLEAN | 是否已领取奖励 / Reward claimed |
| `tx_hash` | VARCHAR(90) | 交易哈希 / Transaction hash |

## 🎮 使用场景 / Usage Scenarios

### 场景 1: 个人挑战 / Scenario 1: Individual Challenge

**目标用户**: 想要养成环保习惯的个人用户  
**示例**: 用户发起"连续30天垃圾分类"挑战  
**效果**: 提高用户粘性，培养长期习惯  

### 场景 2: 企业挑战 / Scenario 2: Corporate Challenge

**目标用户**: 企业希望提升员工环保意识  
**示例**: 企业发起"全员植树挑战"  
**效果**: 提升企业社会责任形象，增强团队凝聚力  

### 场景 3: 社区挑战 / Scenario 3: Community Challenge

**目标用户**: 社区居民希望改善社区环境  
**示例**: 社区发起"低碳出行月"活动  
**效果**: 促进社区共建，提升居民环保意识  

### 场景 4: 节日挑战 / Scenario 4: Festival Challenge

**目标用户**: 节日期间吸引用户参与  
**示例**: 春节期间发起"家庭环保挑战"  
**效果**: 增加节日活跃度，提升品牌曝光  

## 📈 运营策略 / Operations Strategy

### 挑战设计原则 / Challenge Design Principles

1. **目标明确**: 挑战目标具体可量化  
2. **时间合理**: 挑战周期适中（1周-1月）  
3. **奖励吸引人**: 奖励金额具有吸引力  
4. **易于参与**: 参与门槛不高  

### 发布策略 / Publishing Strategy

1. **预热宣传**: 挑战开始前1周预热  
2. **分阶段发布**: 不同类型挑战错峰发布  
3. **社交分享**: 鼓励用户分享挑战进度  
4. **数据跟踪**: 实时监控挑战效果  

### 激励机制 / Incentive Mechanism

1. **进度奖励**: 完成一定进度给予奖励  
2. **排名奖励**: 根据参与度给予额外奖励  
3. **分享奖励**: 分享挑战获得奖励  
4. **复参与奖励**: 重复参与获得更多奖励  

## 🔧 配置管理 / Configuration Management

### 挑战配置 / Challenge Configuration

```bash
# 挑战奖励倍率 / Challenge reward multiplier
CHALLENGE_BONUS_RATE=1.5

# 最大挑战参与人数 / Maximum challenge participants
MAX_CHALLENGE_PARTICIPANTS=10000

# 挑战过期时间（天）/ Challenge expiration days
CHALLENGE_EXPIRY_DAYS=30
```

### 数据库配置 / Database Configuration

```sql
-- 挑战相关索引 / Challenge related indexes
CREATE INDEX idx_challenges_status_time ON challenges(status, created_at);
CREATE INDEX idx_challenge_participants_challenge_wallet ON challenge_participants(challenge_id, wallet_addr);

-- 挑战统计视图 / Challenge statistics view
CREATE VIEW challenge_stats AS
SELECT
    c.id,
    c.title,
    c.current_count,
    c.target_count,
    ROUND(c.current_count::numeric / c.target_count * 100, 2) as completion_rate,
    COUNT(cp.id) as participant_count
FROM challenges c
LEFT JOIN challenge_participants cp ON c.id = cp.challenge_id
GROUP BY c.id, c.title, c.current_count, c.target_count;
```

## 📊 监控指标 / Monitoring Metrics

### 挑战效果指标 / Challenge Effectiveness Metrics

- **参与率**: 参与人数 / 目标人数  
- **完成率**: 完成挑战人数 / 参与人数  
- **平均参与时间**: 用户参与挑战的平均时间  
- **奖励发放率**: 成功领取奖励人数 / 应领取人数  

### 用户行为指标 / User Behavior Metrics

- **挑战发起数**: 每个用户发起的挑战数量  
- **挑战参与数**: 每个用户参与的挑战数量  
- **挑战完成数**: 每个用户完成的挑战数量  
- **挑战转化率**: 参与挑战后继续使用平台的用户比例  

### 运营效果指标 / Operations Effectiveness Metrics

- **挑战活跃度**: 每日活跃挑战数  
- **用户增长**: 通过挑战吸引的新用户数  
- **内容产出**: 用户通过挑战产生的环保行为数  
- **社区互动**: 用户间的挑战分享和互动数量  

## 🎯 最佳实践 / Best Practices

### 挑战创建 / Challenge Creation

1. **目标设定**: 基于用户调研设定合理目标  
2. **奖励设计**: 结合用户偏好设计奖励机制  
3. **规则简化**: 确保挑战规则简单易懂  
4. **视觉设计**: 使用吸引人的视觉元素  

### 挑战运营 / Challenge Operations

1. **进度跟踪**: 实时更新挑战进度  
2. **用户通知**: 重要节点主动通知用户  
3. **内容互动**: 鼓励用户分享挑战经历  
4. **数据分析**: 持续分析挑战效果并优化  

### 挑战维护 / Challenge Maintenance

1. **定期清理**: 清理过期和无效挑战  
2. **数据备份**: 定期备份挑战数据  
3. **性能监控**: 监控挑战系统性能  
4. **用户反馈**: 收集用户反馈并改进  

## 🚀 扩展功能 / Extended Features

### 未来规划 / Future Planning

1. **挑战模板**: 提供常用挑战模板  
2. **智能推荐**: 基于用户行为推荐合适挑战  
3. **挑战组合**: 支持多个挑战组合参与  
4. **挑战积分**: 挑战完成获得积分系统  

### 高级功能 / Advanced Features

1. **地理位置**: 基于地理位置的本地挑战  
2. **时间限制**: 限时挑战增加紧迫感  
3. **团队协作**: 团队成员协作完成挑战  
4. **NFT 奖励**: 挑战完成获得 NFT 徽章  

---

**文档版本**: v3.0.0  
**最后更新**: 2024-01-24  
**相关链接**: [API 文档](./API_v3.0.0_CN.md), [技术文档](./TECHNICAL_v3.0.0_CN.md)
