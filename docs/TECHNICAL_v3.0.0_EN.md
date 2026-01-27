# Technical Documentation v3.0.0

## System Architecture

### Overall Architecture

```
┌─────────────┐
│ Frontend (React) │
└──────┬──────┘
       │ HTTP/HTTPS
       ▼
┌─────────────┐
│ Go Backend API │
└──────┬──────┘
       │
       ├──► PostgreSQL ── Data Storage
       ├──► Redis ─────── Cache/Rate Limiting
       └──► Solana ────── Blockchain Interaction
       └──► Multi-Chain ── Multi-Chain Support
```

### Technology Stack Details

#### Backend Technology Stack

- **Gin**: High-performance HTTP web framework
- **GORM**: ORM framework supporting multiple databases
- **Redis**: Caching and rate limiting
- **JWT**: Stateless authentication
- **Solana Go SDK**: Solana blockchain interaction

#### Frontend Technology Stack

- **React 18**: UI framework
- **Solana Wallet Adapter**: Wallet connection
- **Axios**: HTTP client
- **React Router**: Route management

#### Blockchain Technology Stack

- **Anchor**: Solana smart contract framework
- **Rust**: Contract development language
- **SPL Token**: Token standard

#### Operations Technology Stack

- **Docker**: Containerized deployment
- **WebSocket**: Real-time data monitoring
- **Monitoring & Alerting**: System monitoring
- **Data Analytics**: Operational data analysis

## Core Feature Implementation

### 1. Wallet Authentication

#### Process

1. User connects Solana wallet on frontend
2. Frontend generates signature message
3. User signs with wallet
4. Backend verifies signature
5. Generate JWT token and return

#### Code Location

- Frontend: `frontend/src/components/WalletConnect.jsx`
- Backend: `backend/controller/auth.go`

### 2. Environmental Behavior Submission

#### Process

1. User submits behavior information
2. Backend validates parameters
3. Redis prevents duplicate submission
4. AI fraud detection
5. Process based on risk level
6. Distribute rewards
7. On-chain proof storage

#### Key Technologies

- **AI Fraud Detection**: Integrated with Baidu AI, Alibaba Cloud, AWS, Google Cloud
- **On-chain Proof Storage**: SHA256 hash permanently recorded on Solana chain
- **Multi-Chain Support**: Support for Solana, Ethereum, Polygon, BSC

### 3. Challenge Activity System

#### Features

- **Create Challenges**: Users initiate individual/team/community challenges
- **Join Challenges**: Submit behavior to automatically participate
- **Progress Tracking**: Real-time display of completion progress
- **Reward Distribution**: Automatic reward distribution upon challenge completion

#### Data Models

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

### 4. Marketing Activity System

#### Activity Types

1. **Sign-in Activities**: Consecutive sign-in rewards
2. **Invite Activities**: Invite friends rewards
3. **Daily Tasks**: Complete tasks for rewards
4. **Lucky Draws**: Participate in lottery for random rewards
5. **Flash Sales**: Limited-time high rewards
6. **Festival Activities**: Thematic festival activities

### 5. Blockchain Data Monitoring

#### Implementation

- **WebSocket Subscription**: Real-time monitoring of Solana account changes
- **Polling Mode**: Backup data retrieval method
- **Event Filtering**: Only monitor relevant transaction events
- **Data Storage**: Structured storage to database

#### Monitoring Content

- **Transaction Events**: Reward distribution, proof storage transactions
- **Account Changes**: User balances, contract status
- **Program Logs**: Smart contract execution logs

### 6. Data Analytics Service

#### Analysis Dimensions

- **User Analysis**: Active users, new users, retention rate
- **Behavior Analysis**: Behavior type distribution, time trends
- **Reward Analysis**: Reward distribution statistics, cost analysis
- **Activity Analysis**: Participation rate, conversion rate, ROI

#### Report Types

- **Daily Reports**: Daily data summary
- **Weekly Reports**: Weekly trend analysis
- **Monthly Reports**: Monthly summary
- **Custom Reports**: On-demand analysis

## Database Design

### Core Data Tables

#### User Related
- `users`: User basic information
- `user_wallets`: Wallet address mapping

#### Behavior Related
- `green_behaviors`: Environmental behavior records
- `behavior_rewards`: Reward distribution records

#### Activity Related
- `challenges`: Challenge activities
- `challenge_participants`: Challenge participants
- `marketing_activities`: Marketing activities
- `marketing_participants`: Marketing activity participants

#### System Related
- `sign_in_records`: Sign-in records
- `invite_records`: Invitation records
- `analytics_data`: Data analytics cache

### Index Design

```sql
-- User wallet address index
CREATE INDEX idx_users_wallet_addr ON users(wallet_addr);

-- Behavior type and time index
CREATE INDEX idx_behaviors_type_time ON green_behaviors(behavior_type, submit_time);

-- Challenge status index
CREATE INDEX idx_challenges_status ON challenges(status);

-- Time range query index
CREATE INDEX idx_rewards_created_at ON behavior_rewards(created_at);
```

## Performance Optimization

### Backend Optimization

#### Database Optimization
- **Index Strategy**: Add indexes for commonly queried fields
- **Connection Pool**: Use GORM connection pool management
- **Read-Write Separation**: Support master-slave database configuration

#### Caching Strategy
- **Redis Cache**: Hotspot data caching
- **Application Cache**: Local memory cache
- **CDN Acceleration**: Static resource acceleration

#### Asynchronous Processing
- **Message Queue**: Use Redis as message queue
- **Background Tasks**: Blockchain transaction asynchronous processing
- **Scheduled Tasks**: Data analysis scheduled execution

### Frontend Optimization

#### Code Splitting
- **Route Splitting**: Split code by pages
- **Component Lazy Loading**: Dynamic component import
- **Third-party Library Splitting**: Package large libraries separately

#### Resource Optimization
- **Image Compression**: WebP format and compression
- **Font Optimization**: Use system fonts
- **Bundle Analysis**: Analyze bundle size and optimize

### Blockchain Optimization

#### Transaction Optimization
- **Batch Transactions**: Merge multiple operations
- **Gas Optimization**: Reasonably set transaction parameters
- **Network Selection**: Choose network based on cost

#### Contract Optimization
- **Storage Optimization**: Reduce on-chain storage costs
- **Computation Optimization**: Optimize contract logic complexity
- **Event Optimization**: Reasonable use of event logs

## Security Mechanisms

### Authentication Security

- **JWT Token**: Stateless authentication with expiration time
- **Signature Verification**: Solana wallet signature verification
- **Access Control**: Role-based access control

### Data Security

- **Password Encryption**: User sensitive information encrypted storage
- **Transmission Encryption**: HTTPS transmission encryption
- **API Rate Limiting**: Prevent malicious requests

### Blockchain Security

- **Private Key Management**: Environment variable storage, no hardcoding
- **Transaction Verification**: Multi-signature verification
- **Contract Auditing**: Professional audit agency auditing

## Monitoring and Alerting

### System Monitoring

- **Application Monitoring**: CPU, memory, disk usage
- **Service Monitoring**: API response time, error rate
- **Database Monitoring**: Connection count, slow queries

### Business Monitoring

- **User Monitoring**: Active users, registration conversion
- **Transaction Monitoring**: Transaction success rate, confirmation time
- **Activity Monitoring**: Activity participation, reward distribution

### Alerting Mechanism

- **Threshold Alerting**: Automatic alerting when exceeding set thresholds
- **Anomaly Detection**: Automatic detection of abnormal patterns
- **Notification Channels**: Email, SMS, Webhook

## Scalability Design

### Horizontal Scaling

- **Stateless Design**: Backend services stateless, support multiple instances
- **Load Balancing**: Nginx load balancing
- **Database Sharding**: Support data sharding

### Vertical Scaling

- **Cache Layer**: Multi-level cache architecture
- **Asynchronous Processing**: Message queue asynchronous processing
- **CDN Acceleration**: Global content distribution

### Feature Extension

- **Plugin Architecture**: Support feature modularization
- **API Version Control**: Support multi-version API
- **Configuration Hot Update**: Support dynamic configuration updates

## Deployment Architecture

### Production Architecture

```
┌─────────────────┐
│   Load Balancer │  Nginx
│   (Nginx)       │
└─────────┬───────┘
          │
    ┌─────┴─────┐
    │  API Gateway │  Routing, Authentication, Rate Limiting
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
    │ Blockchain │  Solana + Multi-Chain
    │  Networks │
    └───────────┘
```

### Containerized Deployment

- **Docker**: Application containerization
- **Kubernetes**: Container orchestration
- **Helm**: Application package management
- **CI/CD**: Automated deployment pipeline

## Fault Handling

### Common Faults

#### Database Connection Failed
- **Cause**: Network issues, configuration errors, insufficient resources
- **Handling**: Check configuration, restart service, scale resources

#### Redis Cache Failure
- **Cause**: Insufficient memory, service downtime
- **Handling**: Scale memory, cluster deployment, master-slave switching

#### Blockchain Network Exception
- **Cause**: Network congestion, RPC node issues
- **Handling**: Switch nodes, retry mechanism, multi-node backup

### Rollback Strategy

1. **Code Rollback**: Git version rollback
2. **Database Rollback**: Backup restoration
3. **Contract Rollback**: Deploy old version contract

### Emergency Plan

1. **Monitoring Alerting**: 24/7 monitoring
2. **Backup Strategy**: Multiple backups
3. **Recovery Process**: Detailed recovery steps
4. **Communication Mechanism**: Timely user notification

---

**Document Version**: v3.0.0  
**Last Update**: 2024-01-24
