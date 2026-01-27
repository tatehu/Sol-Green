# API Documentation v3.0.0

## Basic Information

- **Base URL**: `http://localhost:8080/api/v1`
- **Authentication**: Bearer Token (JWT)
- **Data Format**: JSON

## Common Response Format

### Success Response

```json
{
  "data": {},
  "msg": "Operation successful"
}
```

### Error Response

```json
{
  "error": "Error message",
  "detail": "Detailed error message (optional)"
}
```

## Authentication API

### Wallet Login

**POST** `/auth/wallet`

User uses Solana wallet signature for login authentication.

#### Request Parameters

```json
{
  "wallet_addr": "string",  // Solana wallet address
  "signature": "string",     // Signature (hex encoded)
  "message": "string"        // Signature message
}
```

#### Response Example

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "wallet_addr": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"
}
```

#### Status Codes

- `200`: Login successful
- `400`: Parameter error
- `401`: Signature verification failed

## Green Behavior API

### Submit Green Behavior

**POST** `/green/behavior/submit`

Submit environmental behavior application, system will perform AI fraud detection and automatically distribute rewards.

#### Request Parameters

```json
{
  "behavior_type": "waste_sorting",  // Behavior type: waste_sorting | tree_planting | low_carbon_travel
  "media_urls": [                    // Media file URLs array (1-5)
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "location": "Shanghai Zhangjiang Software Park"  // Location (optional)
}
```

#### Response Example

**Auto Approved (Low Risk)**
```json
{
  "msg": "Verification passed, reward distributed",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "approved",
  "tx_hash": "5j7s8K9...",
  "proof": "a1b2c3d4..."
}
```

**Pending Review (High Risk)**
```json
{
  "msg": "Submission successful, pending manual review",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending_review"
}
```

### Get Behavior Status

**GET** `/green/behavior/:id`

Query detailed status of specified behavior record.

### Third-Party Verification

**POST** `/green/partner/verify`

Verify environmental behavior through third-party verification organizations.

#### Request Parameters

```json
{
  "behavior_id": "550e8400-e29b-41d4-a716-446655440000",
  "partner_id": "alashan_see",  // Partner ID: alashan_see | gov_environment | unep | wwf | greenpeace
  "verify_code": "SEE123456"    // Verification code
}
```

## Challenge Activity API

### Create Challenge

**POST** `/challenges`

Create a new environmental challenge activity.

#### Request Parameters

```json
{
  "title": "Waste Sorting Challenge Week",
  "description": "Complete 10 waste sorting this week",
  "challenge_type": "individual",  // Challenge type: individual | team | community
  "behavior_type": "waste_sorting",
  "reward_amount": 5000,
  "target_count": 50,
  "start_time": "2024-01-25T00:00:00Z",
  "end_time": "2024-02-01T23:59:59Z"
}
```

### Get Challenges

**GET** `/challenges?status=active&challenge_type=individual&page=1&page_size=20`

### Join Challenge

**POST** `/challenges/:id/join`

Join challenge using approved environmental behavior.

### Claim Challenge Reward

**POST** `/challenges/:id/claim`

Claim reward after challenge completion.

## Marketing Activity API

### Create Marketing Activity

**POST** `/marketing/activities`

Create new marketing activity (admin only).

#### Request Parameters

```json
{
  "title": "Daily Sign-in Activity",
  "description": "Sign in for 7 consecutive days to get extra rewards",
  "activity_type": "sign_in",  // Activity type: sign_in | invite | daily_task | lucky_draw | flash_sale | festival
  "reward_amount": 1000,
  "start_time": "2024-01-25T00:00:00Z",
  "end_time": "2024-02-25T23:59:59Z"
}
```

### Get Marketing Activities

**GET** `/marketing/activities?status=active&activity_type=sign_in`

### Join Marketing Activity

**POST** `/marketing/activities/:id/join`

Join marketing activity.

### Claim Marketing Reward

**POST** `/marketing/activities/:id/claim`

Claim marketing activity reward.

## Operations Management API

### Get Statistics

**GET** `/admin/stats?days=7`

Get operations statistics.

#### Response Example

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

### Get Analytics

**GET** `/admin/analytics?start_date=2024-01-01&end_date=2024-01-31`

Get detailed data analytics.

### Activate Marketing Activity

**POST** `/admin/activities/:id/activate`

Activate marketing activity (admin only).

## Contract Management API

### Get Contract Info

**GET** `/admin/contract/info?program_id=<program-id>`

Get smart contract information.

### Check Upgrade

**GET** `/admin/contract/upgrade/check?program_id=<program-id>`

Check if upgrade is available.

### Prepare Upgrade

**POST** `/admin/contract/upgrade/prepare`

Prepare contract upgrade.

### Execute Upgrade

**POST** `/admin/contract/upgrade/execute`

Execute contract upgrade.

### Get Upgrade History

**GET** `/admin/contract/upgrade/history?program_id=<program-id>`

Get contract upgrade history.

## Error Codes

| Status Code | Description |
|------------|-------------|
| 200 | Request successful |
| 400 | Request parameter error |
| 401 | Unauthenticated or authentication failed |
| 403 | Insufficient permissions |
| 404 | Resource not found |
| 429 | Too many requests |
| 500 | Internal server error |

## Rate Limiting

- **IP Rate Limit**: Maximum 60 requests per minute per IP
- **Duplicate Prevention**: Same user same behavior type can only submit once within 5 minutes
