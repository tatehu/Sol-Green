# Challenge Activity Features v3.0.0

## 📋 Feature Overview

The challenge activity feature allows users to initiate and participate in environmental challenges, increasing user engagement and platform vitality through social and gamified approaches.

## 🎯 Feature Characteristics

### 1. Challenge Types

- **Individual Challenge**: Personal environmental goals
  - Example: "Sort waste 30 times this month"
  - Suitable for: Personal habit formation

- **Team Challenge**: Team collaboration on environmental tasks
  - Example: "Team plants 100 trees"
  - Suitable for: Corporate/community activities

- **Community Challenge**: Community-level environmental activities
  - Example: "Community low-carbon travel week"
  - Suitable for: Community co-construction

### 2. Challenge Process

```
Create Challenge → Activate Challenge → User Participation → Achieve Goal → Distribute Rewards
```

#### Detailed Process

1. **Challenge Creation**
   - User fills in challenge information
   - Set rewards and participation conditions
   - Challenge initial status is "draft"

2. **Challenge Activation**
   - Creator or admin activates challenge
   - Challenge status becomes "active"
   - Start accepting user participation

3. **User Participation**
   - Users browse and select challenges
   - Submit qualified environmental behaviors
   - System automatically associates behavior with challenge

4. **Goal Achievement**
   - Real-time statistics of participants
   - Progress bar shows completion status
   - Automatically mark as completed when goal achieved

5. **Reward Distribution**
   - Participants claim rewards after challenge completion
   - Support multiple reward types
   - On-chain reward records

### 3. Reward Mechanism

#### Basic Rewards
- Distribute basic rewards based on behavior type
- Reference standard reward configuration

#### Challenge Rewards
- Additional rewards after challenge completion
- Reward multiplier configurable (default 1.5x)

#### Reward Calculation Formula

```go
// Challenge reward calculation
challengeReward = baseReward * challengeBonusRate
totalReward = baseReward + challengeReward

// Example
// Waste sorting base reward: 1000 SOLGREEN
// Challenge reward: 1000 * 1.5 = 1500 SOLGREEN
// Total reward: 1000 + 1500 = 2500 SOLGREEN
```

## 📡 API Interface

### Create Challenge

**POST** `/api/v1/challenges`

Create a new environmental challenge activity.

#### Request Parameters

```json
{
  "title": "Waste Sorting Challenge Week",
  "description": "Complete 10 waste sorting this week",
  "challenge_type": "individual",
  "behavior_type": "waste_sorting",
  "reward_amount": 5000,
  "target_count": 50,
  "start_time": "2024-01-25T00:00:00Z",
  "end_time": "2024-02-01T23:59:59Z",
  "rules": "Upload photo proof for each waste sorting"
}
```

### Get Challenges List

**GET** `/api/v1/challenges?status=active&challenge_type=individual&page=1&page_size=20`

### Get Challenge Details

**GET** `/api/v1/challenges/:id`

### Join Challenge

**POST** `/api/v1/challenges/:id/join`

Join challenge using approved environmental behavior.

#### Request Parameters

```json
{
  "behavior_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

### Activate Challenge

**POST** `/api/v1/challenges/:id/activate`

Activate challenge (creator or admin only).

### Claim Challenge Reward

**POST** `/api/v1/challenges/:id/claim`

Claim reward after challenge completion.

## 📊 Data Models

### Challenge Table

| Field | Type | Description |
|-------|------|-------------|
| `id` | VARCHAR(36) | Primary key |
| `title` | VARCHAR(200) | Challenge title |
| `description` | TEXT | Challenge description |
| `challenge_type` | VARCHAR(20) | Challenge type |
| `behavior_type` | VARCHAR(20) | Associated behavior type |
| `reward_amount` | BIGINT | Reward amount |
| `target_count` | INT | Target participant count |
| `current_count` | INT | Current participant count |
| `start_time` | TIMESTAMP | Start time |
| `end_time` | TIMESTAMP | End time |
| `status` | VARCHAR(20) | Status |
| `rules` | TEXT | Challenge rules |
| `created_by` | VARCHAR(44) | Creator wallet |
| `created_at` | TIMESTAMP | Created time |
| `updated_at` | TIMESTAMP | Updated time |

### Challenge Participant Table

| Field | Type | Description |
|-------|------|-------------|
| `id` | VARCHAR(36) | Primary key |
| `challenge_id` | VARCHAR(36) | Challenge ID |
| `wallet_addr` | VARCHAR(44) | Wallet address |
| `behavior_id` | VARCHAR(36) | Associated behavior ID |
| `joined_at` | TIMESTAMP | Join time |
| `reward_claimed` | BOOLEAN | Reward claimed |
| `tx_hash` | VARCHAR(90) | Transaction hash |

## 🎮 Usage Scenarios

### Scenario 1: Individual Challenge
**Target Users**: Individuals wanting to develop environmental habits  
**Example**: User initiates "Sort waste 30 times continuously" challenge  
**Effect**: Increase user retention, cultivate long-term habits  

### Scenario 2: Corporate Challenge
**Target Users**: Companies wanting to raise employee environmental awareness  
**Example**: Company initiates "All employees plant trees" challenge  
**Effect**: Enhance corporate social responsibility image, strengthen team cohesion  

### Scenario 3: Community Challenge
**Target Users**: Community residents wanting to improve community environment  
**Example**: Community initiates "Low-carbon travel month" activity  
**Effect**: Promote community co-construction, raise resident environmental awareness  

## 📈 Operations Strategy

### Challenge Design Principles

1. **Clear Goals**: Challenge goals are specific and quantifiable
2. **Reasonable Time**: Challenge period is appropriate (1 week-1 month)
3. **Attractive Rewards**: Reward amounts are attractive
4. **Easy Participation**: Low participation threshold

### Publishing Strategy

1. **Warm-up Promotion**: Warm-up 1 week before challenge starts
2. **Staggered Release**: Different types of challenges released at different times
3. **Social Sharing**: Encourage users to share challenge progress
4. **Data Tracking**: Real-time monitoring of challenge effectiveness

## 🔧 Configuration Management

### Challenge Configuration

```bash
# Challenge reward multiplier
CHALLENGE_BONUS_RATE=1.5

# Maximum challenge participants
MAX_CHALLENGE_PARTICIPANTS=10000

# Challenge expiration days
CHALLENGE_EXPIRY_DAYS=30
```

---

**Document Version**: v3.0.0  
**Last Update**: 2024-01-24  
**Related Links**: [API Documentation](./API_v3.0.0_EN.md), [Technical Documentation](./TECHNICAL_v3.0.0_EN.md)
