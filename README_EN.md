# 🌱 Sol-Green - Environmental Behavior Reward Platform

A Solana blockchain-based environmental behavior recording and reward system that incentivizes users to participate in environmental activities through Web3 technology.

## 🎥 Demo Videos

- [Watch Sol-Green Demo Part 1](./assets/Sol%20Green%20Part%201.mov)
- [Watch Sol-Green Demo Part 2](./assets/Sol%20Green%20Part%202.mov)

<video src="./assets/Sol%20Green%20Part%201.mov" controls width="640">
  Your browser does not support the video tag. Please click the link above to watch Part 1.
</video>

<video src="./assets/Sol%20Green%20Part%202.mov" controls width="640">
  Your browser does not support the video tag. Please click the link above to watch Part 2.
</video>

## 📋 Project Overview

Sol-Green is a decentralized environmental reward platform where users can earn token rewards by submitting environmental behaviors (waste sorting, tree planting, low-carbon travel, etc.). The project uses AI fraud detection, on-chain proof storage, and third-party organization verification to ensure behavior authenticity and credibility.

## ✨ Core Features

- 🌱 **Environmental Behavior Recording**: Waste sorting, tree planting, low-carbon travel behaviors recorded on-chain
- 🎯 **Challenge Activities**: Create and participate in environmental challenges for extra rewards
- 🎁 **Marketing Activities**: Sign-in, invite, daily tasks, lucky draw, and various marketing activities
- 🤖 **AI Fraud Detection**: Integrated with leading global AI detection technologies to ensure behavior authenticity
- 🌍 **Third-Party Verification**: Supports global organizations like UNEP, WWF, Greenpeace for verification
- 💰 **Token Rewards**: Instant reward distribution based on Solana blockchain
- 🔒 **On-Chain Proof**: All behaviors permanently recorded on Solana chain
- 🔍 **Blockchain Monitoring**: Real-time monitoring of on-chain events and operational data
- 📊 **Data Analytics**: Complete operational data analysis and reporting
- 🌐 **Multi-Chain Support**: Supports Solana, Ethereum, Polygon, BSC, and more
- 🔄 **Contract Upgrade**: Complete contract upgrade management system

## 🏗️ Project Structure

```
trends_solana/
├── backend/          # Go backend service
├── frontend/         # React frontend application
├── contracts/        # Solana smart contracts
├── docs/             # Project documentation
├── scripts/          # Deployment scripts
└── docker-compose.yml
```

## 🚀 Quick Start

### Environment Requirements

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose (optional)
- Solana CLI (for contract deployment)
- Anchor CLI (for contract development)

### 1. Clone Project

```bash
git clone <repository-url>
cd trends_solana
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env file with necessary configurations
```

### 3. Start Backend Service

#### Option 1: Docker Compose (Recommended)

```bash
docker-compose up -d
```

#### Option 2: Local Development

```bash
cd backend
go mod download
go run main.go
```

Backend service will start at `http://localhost:8080`.

### 4. Start Frontend

```bash
cd frontend
npm install
npm start
```

Frontend application will start at `http://localhost:3000`.

## 📡 API Documentation

Detailed API documentation: [API Documentation v3.0.0](./docs/API_v3.0.0_EN.md)

### Main Endpoints

- `POST /api/v1/auth/wallet` - Wallet login
- `POST /api/v1/green/behavior/submit` - Submit environmental behavior
- `GET /api/v1/green/behavior/:id` - Query behavior status
- `POST /api/v1/challenges` - Create challenge
- `GET /api/v1/challenges` - Get challenges list
- `POST /api/v1/marketing/activities/:id/join` - Join marketing activity
- `GET /api/v1/admin/stats` - Get statistics (admin)

## 🧪 Testing

```bash
# Run all tests
./scripts/test.sh

# Run backend tests
cd backend
go test -v ./tests/...

# Run frontend tests
cd frontend
npm test
```

## 🚢 Deployment

### Using Docker Compose

```bash
./scripts/deploy.sh
```

### Manual Deployment

See [Deployment Documentation](./docs/DEPLOYMENT_v3.0.0_EN.md) for detailed instructions.

## 📚 Documentation

- [Quick Start Guide](./docs/QUICKSTART_v3.0.0_EN.md)
- [Development Guide](./docs/DEVELOPMENT_v3.0.0_EN.md)
- [Technical Documentation](./docs/TECHNICAL_v3.0.0_EN.md)
- [Deployment Guide](./docs/DEPLOYMENT_v3.0.0_EN.md)
- [API Documentation](./docs/API_v3.0.0_EN.md)
- [Changelog](./docs/CHANGELOG_v3.0.0_EN.md)

## 🔧 Technology Stack

### Backend
- Go 1.21+
- Gin Web Framework
- GORM ORM
- Redis Cache
- JWT Authentication
- Solana Go SDK

### Blockchain
- Solana
- Anchor Framework
- Rust
- SPL Token

### Frontend
- React 18
- Solana Wallet Adapter
- Axios
- React Router

## ✨ Key Features

1. **Challenge Activity System**: Create and participate in environmental challenges for extra rewards
2. **AI Fraud Detection**: Integrated with Baidu AI, Alibaba Cloud, AWS, Google Cloud, and other leading AI detection technologies
3. **Third-Party Verification**: Supports UNEP, WWF, Greenpeace, Alashan SEE, and other global environmental organizations
4. **Multi-Environment Support**: Development, test (Solana Devnet), and mainnet environments with one-click switching
5. **On-Chain Proof**: Behavior hashes permanently recorded on Solana chain
6. **Anti-Fraud Mechanism**: Redis rate limiting + duplicate prevention
7. **Security Mechanism**: JWT authentication, CORS, SQL injection protection

## 🔒 Security Notes

1. **Private Key Management**: Never commit private keys to code repository, use environment variables
2. **JWT Secret**: Production environment must use strong random secret
3. **API Rate Limiting**: Redis-based request rate limiting implemented
4. **AI Fraud Detection**: Integrated third-party content security APIs for image/video detection

## 🤝 Contributing

Welcome to submit Issues and Pull Requests!

## 📄 License

MIT License

## 📞 Contact

For questions, please submit an Issue or contact project maintainers.
