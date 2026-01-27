# Quick Start Guide v3.0.0

## 5-Minute Quick Experience

### 1. Environment Setup

Ensure installed:
- Go 1.21+ (or use Docker)
- Node.js 18+
- Docker & Docker Compose (Recommended)

### 2. Clone Project

```bash
git clone <repository-url>
cd trends_solana
```

### 3. Configure Environment

```bash
cp .env.example .env
# Edit .env file (optional, defaults can be used for development)
```

**Detailed Configuration Guide**:  
Please refer to [Configuration Setup Guide](./CONFIG_SETUP_GUIDE_v3.0.0_EN.md) for details on how to obtain configuration values.

### 4. Start Services (Easiest Way)

```bash
# Start with Docker Compose
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f backend
```

### 5. Access Application

- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/api/v1/health
- **Frontend App**: needs separate start (see below)

### 6. Start Frontend

```bash
cd frontend
npm install
npm start
```

Frontend will start at http://localhost:3000

### 7. Test API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Wallet login (example)
curl -X POST http://localhost:8080/api/v1/auth/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_addr": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
    "signature": "test_signature",
    "message": "test_message"
  }'
```

## Local Development Mode

### Backend Development

```bash
# Enter backend directory
cd backend

# Install dependencies
go mod download

# Start PostgreSQL and Redis (using Docker)
cd ..
docker-compose up -d postgres redis

# Configure environment (can skip database if using SQLite)
export DB_TYPE=sqlite

# Run backend
cd backend
go run main.go
```

Backend will start at `http://localhost:8080`

### Frontend Development

```bash
cd frontend
npm install
npm start
```

Frontend will start at `http://localhost:3000`

### Smart Contract Development

```bash
# Install Anchor CLI
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install latest
avm use latest

# Enter contract directory
cd contracts/sol-green

# Build contract
anchor build

# Run tests
anchor test
```

## Verify Project

Run verification script to check project integrity:

```bash
./scripts/verify.sh
```

## Common Issues

### Q: Backend startup failed?

**A:** Check:
1. Port 8080 occupied
2. Database/Redis running
3. Environment variables correct

Solution:
```bash
# Check port usage
lsof -i :8080

# Check Docker services
docker-compose ps

# View backend logs
docker-compose logs backend
```

### Q: Frontend cannot connect to backend?

**A:** 
1. Check backend running at http://localhost:8080
2. Check CORS configuration
3. Check browser console errors

### Q: Solana transaction failed?

**A:**
1. Check wallet has enough SOL
2. Confirm network configuration (devnet/mainnet)
3. Verify RPC node available

## Next Steps

- Read [README.md](../README.md) for project details
- View [Development Guide](./DEVELOPMENT_v3.0.0_EN.md) to start development
- Refer to [API Documentation](./API_v3.0.0_EN.md) for API details
- View [Deployment Guide](./DEPLOYMENT_v3.0.0_EN.md) for production deployment

---

**Version**: v3.0.0  
**Update Date**: 2024-01-24
