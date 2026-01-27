# 营销活动系统

## 📋 概述

营销活动系统通过多样化的活动形式提高用户活跃度和项目知名度，包括签到、邀请、每日任务、抽奖等多种活动类型。

## 🎯 活动类型

### 1. 签到活动 (Sign In)

**特点：**
- 每日签到获得奖励
- 连续签到获得额外奖励
- 签到天数越多，奖励越丰厚

**配置示例：**
```json
{
  "consecutive_bonus": true,
  "bonus_multiplier": 1.2,
  "max_consecutive_days": 30
}
```

### 2. 邀请活动 (Invite)

**特点：**
- 邀请新用户注册
- 被邀请人完成首次行为后，邀请人获得奖励
- 支持多级邀请奖励

**配置示例：**
```json
{
  "inviter_reward": 500,
  "invitee_reward": 200,
  "max_invites": 10
}
```

### 3. 每日任务 (Daily Task)

**特点：**
- 每日完成指定任务
- 任务类型多样（垃圾分类、低碳出行等）
- 完成任务获得奖励

**配置示例：**
```json
{
  "tasks": [
    {"type": "waste_sorting", "count": 3, "reward": 300},
    {"type": "low_carbon_travel", "count": 1, "reward": 200}
  ]
}
```

### 4. 抽奖活动 (Lucky Draw)

**特点：**
- 参与活动获得抽奖机会
- 不同奖品和概率
- 增加用户参与趣味性

**配置示例：**
```json
{
  "prizes": [
    {"name": "一等奖", "reward": 10000, "probability": 0.01},
    {"name": "二等奖", "reward": 5000, "probability": 0.05},
    {"name": "三等奖", "reward": 1000, "probability": 0.2}
  ]
}
```

### 5. 限时抢购 (Flash Sale)

**特点：**
- 限时开放
- 限量参与
- 高额奖励

### 6. 节日活动 (Festival)

**特点：**
- 特定节日期间开放
- 主题化活动
- 特殊奖励

## 📡 API 接口

### 创建营销活动

```bash
POST /api/v1/marketing/activities
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "title": "每日签到活动",
  "description": "连续签到7天获得额外奖励",
  "activity_type": "sign_in",
  "start_time": "2024-01-25T00:00:00Z",
  "end_time": "2024-02-25T23:59:59Z",
  "reward_amount": 1000,
  "max_participants": 0,
  "rules": "每日签到一次，连续签到7天获得双倍奖励",
  "config": "{\"consecutive_bonus\": true}"
}
```

### 获取活动列表

```bash
GET /api/v1/marketing/activities?status=active&activity_type=sign_in
```

### 参与活动

```bash
POST /api/v1/marketing/activities/:id/join
Authorization: Bearer <token>
```

### 领取奖励

```bash
POST /api/v1/marketing/activities/:id/claim
Authorization: Bearer <token>
```

## 📊 活动效果分析

### 关键指标

- **参与率**: 参与人数 / 总用户数
- **完成率**: 完成活动人数 / 参与人数
- **奖励发放**: 总奖励数量
- **用户留存**: 活动期间用户留存率

### 运营建议

1. **定期活动**: 每周/每月固定活动，培养用户习惯
2. **节日活动**: 结合节日主题，提高参与度
3. **限时活动**: 营造紧迫感，促进参与
4. **奖励梯度**: 设置合理的奖励梯度，激励持续参与

## 🎮 活动策略

### 新用户激活

- 注册奖励
- 首次行为奖励
- 新手任务

### 用户留存

- 每日签到
- 连续奖励
- 成长体系

### 用户增长

- 邀请奖励
- 分享奖励
- 推荐奖励

### 活跃度提升

- 每日任务
- 挑战活动
- 排行榜

## 📈 数据追踪

系统自动追踪：
- 活动参与人数
- 奖励发放数量
- 用户行为数据
- 活动效果分析

## 🔧 配置管理

活动配置支持：
- 动态调整奖励
- 修改活动规则
- 暂停/恢复活动
- 活动时间调整
