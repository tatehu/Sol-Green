# API 文档 v3.0.0

## 基础信息 / Basic Information

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式 / Authentication**: Bearer Token (JWT)
- **数据格式 / Data Format**: JSON

## 通用响应格式 / Common Response Format

### 成功响应 / Success Response

```json
{
  "data": {},
  "msg": "操作成功" // "Operation successful"
}
```

### 错误响应 / Error Response

```json
{
  "error": "错误信息", // "Error message"
  "detail": "详细错误信息（可选）" // "Detailed error message (optional)"
}
```

## 认证接口 / Authentication API

### 钱包登录 / Wallet Login

**POST** `/auth/wallet`

用户使用 Solana 钱包签名进行登录认证。  
User uses Solana wallet signature for login authentication.

#### 请求参数 / Request Parameters

```json
{
  "wallet_addr": "string",  // Solana 钱包地址 / Solana wallet address
  "signature": "string",     // 签名（hex 编码）/ Signature (hex encoded)
  "message": "string"        // 签名消息 / Signature message
}
```

#### 响应示例 / Response Example

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "wallet_addr": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"
}
```

#### 状态码 / Status Codes

- `200`: 登录成功 / Login successful
- `400`: 参数错误 / Parameter error
- `401`: 签名验证失败 / Signature verification failed

## 环保行为接口 / Green Behavior API

### 提交环保行为 / Submit Green Behavior

**POST** `/green/behavior/submit`

提交环保行为申报，系统将进行 AI 反欺诈检测并自动发放奖励。  
Submit environmental behavior application, system will perform AI fraud detection and automatically distribute rewards.

#### 请求参数 / Request Parameters

```json
{
  "behavior_type": "waste_sorting",  // 行为类型 / Behavior type: waste_sorting | tree_planting | low_carbon_travel
  "media_urls": [                    // 媒体文件 URL 数组（1-5 个）/ Media file URLs array (1-5)
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "location": "上海张江软件园"        // 行为地点（可选）/ Location (optional)
}
```

#### 响应示例 / Response Example

**自动通过（低风险）/ Auto Approved (Low Risk)**
```json
{
  "msg": "认证通过，奖励已发放", // "Verification passed, reward distributed"
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "approved",
  "tx_hash": "5j7s8K9...",
  "proof": "a1b2c3d4..."
}
```

**待审核（高风险）/ Pending Review (High Risk)**
```json
{
  "msg": "提交成功，进入人工审核", // "Submission successful, pending manual review"
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending_review"
}
```

### 查询行为状态 / Get Behavior Status

**GET** `/green/behavior/:id`

查询指定行为记录的详细状态。  
Query detailed status of specified behavior record.

### 第三方认证 / Third-Party Verification

**POST** `/green/partner/verify`

通过第三方认证机构验证环保行为。  
Verify environmental behavior through third-party verification organizations.

#### 请求参数 / Request Parameters

```json
{
  "behavior_id": "550e8400-e29b-41d4-a716-446655440000",
  "partner_id": "alashan_see",  // 认证机构 ID / Partner ID: alashan_see | gov_environment | unep | wwf | greenpeace
  "verify_code": "SEE123456"    // 认证码 / Verification code
}
```

## 挑战活动接口 / Challenge Activity API

### 创建挑战 / Create Challenge

**POST** `/challenges`

创建新的环保挑战活动。  
Create a new environmental challenge activity.

#### 请求参数 / Request Parameters

```json
{
  "title": "垃圾分类挑战周", // "Waste Sorting Challenge Week"
  "description": "本周完成 10 次垃圾分类", // "Complete 10 waste sorting this week"
  "challenge_type": "individual",  // 挑战类型 / Challenge type: individual | team | community
  "behavior_type": "waste_sorting",
  "reward_amount": 5000,
  "target_count": 50,
  "start_time": "2024-01-25T00:00:00Z",
  "end_time": "2024-02-01T23:59:59Z"
}
```

### 获取挑战列表 / Get Challenges

**GET** `/challenges?status=active&challenge_type=individual&page=1&page_size=20`

### 参与挑战 / Join Challenge

**POST** `/challenges/:id/join`

使用已审核通过的环保行为参与挑战。  
Join challenge using approved environmental behavior.

### 领取挑战奖励 / Claim Challenge Reward

**POST** `/challenges/:id/claim`

挑战完成后领取奖励。  
Claim reward after challenge completion.

## 营销活动接口 / Marketing Activity API

### 创建营销活动 / Create Marketing Activity

**POST** `/marketing/activities`

创建新的营销活动（仅管理员）。  
Create new marketing activity (admin only).

#### 请求参数 / Request Parameters

```json
{
  "title": "每日签到活动", // "Daily Sign-in Activity"
  "description": "连续签到7天获得额外奖励", // "Sign in for 7 consecutive days to get extra rewards"
  "activity_type": "sign_in",  // 活动类型 / Activity type: sign_in | invite | daily_task | lucky_draw | flash_sale | festival
  "reward_amount": 1000,
  "start_time": "2024-01-25T00:00:00Z",
  "end_time": "2024-02-25T23:59:59Z"
}
```

### 获取活动列表 / Get Marketing Activities

**GET** `/marketing/activities?status=active&activity_type=sign_in`

### 参与活动 / Join Marketing Activity

**POST** `/marketing/activities/:id/join`

参与营销活动。  
Join marketing activity.

### 领取奖励 / Claim Marketing Reward

**POST** `/marketing/activities/:id/claim`

领取营销活动奖励。  
Claim marketing activity reward.

## 运营管理接口 / Operations Management API

### 获取统计数据 / Get Statistics

**GET** `/admin/stats?days=7`

获取运营统计数据。  
Get operations statistics.

#### 响应示例 / Response Example

```json
{
  "period": {
    "start_date": "2024-01-17",
    "end_date": "2024-01-24",
    "days": 7
  },
  "users": {
    "total": 1000,
    "active": 500
  },
  "behaviors": {
    "total": 5000
  },
  "rewards": {
    "total": 10000000
  },
  "chain_stats": {
    "current_slot": 12345678,
    "program_id": "..."
  }
}
```

### 获取分析数据 / Get Analytics

**GET** `/admin/analytics?start_date=2024-01-01&end_date=2024-01-31`

获取详细数据分析。  
Get detailed data analytics.

### 激活营销活动 / Activate Marketing Activity

**POST** `/admin/activities/:id/activate`

激活营销活动（仅管理员）。  
Activate marketing activity (admin only).

## 合约管理接口 / Contract Management API

### 获取合约信息 / Get Contract Info

**GET** `/admin/contract/info?program_id=<program-id>`

获取智能合约信息。  
Get smart contract information.

### 检查升级 / Check Upgrade

**GET** `/admin/contract/upgrade/check?program_id=<program-id>`

检查是否有可用升级。  
Check if upgrade is available.

### 准备升级 / Prepare Upgrade

**POST** `/admin/contract/upgrade/prepare`

准备合约升级。  
Prepare contract upgrade.

### 执行升级 / Execute Upgrade

**POST** `/admin/contract/upgrade/execute`

执行合约升级。  
Execute contract upgrade.

### 升级历史 / Get Upgrade History

**GET** `/admin/contract/upgrade/history?program_id=<program-id>`

获取合约升级历史。  
Get contract upgrade history.

## 错误码说明 / Error Codes

| 状态码 / Status Code | 说明 / Description |
|---------------------|-------------------|
| 200 | 请求成功 / Request successful |
| 400 | 请求参数错误 / Request parameter error |
| 401 | 未认证或认证失败 / Unauthenticated or authentication failed |
| 403 | 权限不足 / Insufficient permissions |
| 404 | 资源不存在 / Resource not found |
| 429 | 请求过于频繁 / Too many requests |
| 500 | 服务器内部错误 / Internal server error |

## 限流说明 / Rate Limiting

- **IP 限流 / IP Rate Limit**: 每个 IP 每分钟最多 60 个请求 / Maximum 60 requests per minute per IP
- **防重复提交 / Duplicate Prevention**: 同一用户同一行为类型 5 分钟内只能提交一次 / Same user same behavior type can only submit once within 5 minutes
