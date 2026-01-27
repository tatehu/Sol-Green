# Sol-Green Configuration Setup Guide v3.0.0

This document provides detailed instructions on how to obtain values for all configuration items in the `.env.example` file.

## Table of Contents

1. [Environment Type Configuration](#environment-type-configuration)
2. [Database Configuration](#database-configuration)
3. [Redis Configuration](#redis-configuration)
4. [JWT Secret Key](#jwt-secret-key)
5. [Solana Configuration](#solana-configuration)
6. [Reward Configuration](#reward-configuration)
7. [AI Fraud Detection Configuration](#ai-fraud-detection-configuration)
8. [Third-Party Verification Configuration](#third-party-verification-configuration)
9. [Rate Limiting Configuration](#rate-limiting-configuration)
10. [Multi-Chain Configuration](#multi-chain-configuration)
11. [Contract Upgrade Configuration](#contract-upgrade-configuration)
12. [Frontend API Address](#frontend-api-address)

---

## Environment Type Configuration

### ENVIRONMENT

**Configuration Item**: `ENVIRONMENT`

**Valid Values**:
- `dev`: Development environment (local development)
- `test`: Test environment (Solana Devnet)
- `mainnet`: Mainnet environment (Solana Mainnet)

**Setup Steps**:
1. Choose `dev` or `test` for development phase
2. Choose `mainnet` for production deployment

**Example**:
```bash
ENVIRONMENT=test
```

---

## Database Configuration

### DB_TYPE

**Configuration Item**: `DB_TYPE`

**Valid Values**:
- `sqlite`: SQLite database (default, no additional configuration needed)
- `postgres`: PostgreSQL database

**Setup Steps**:

#### SQLite (Recommended for Development)
No additional configuration needed:
```bash
DB_TYPE=sqlite
```

#### PostgreSQL (Recommended for Production)
1. **Install PostgreSQL**
   ```bash
   # macOS
   brew install postgresql
   brew services start postgresql
   
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   sudo systemctl start postgresql
   
   # Windows
   # Download and install PostgreSQL: https://www.postgresql.org/download/windows/
   ```

2. **Create Database and User**
   ```bash
   # Login to PostgreSQL
   psql postgres
   
   # Create database
   CREATE DATABASE solgreen;
   
   # Create user (optional)
   CREATE USER solgreen_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE solgreen TO solgreen_user;
   ```

3. **Configure Connection String**
   ```bash
   DB_TYPE=postgres
   DATABASE_URL=postgres://solgreen_user:your_password@localhost:5432/solgreen?sslmode=disable
   ```

---

## Redis Configuration

### REDIS_URL

**Configuration Item**: `REDIS_URL`

**Format**: `host:port`

**Setup Steps**:

1. **Install Redis**
   ```bash
   # macOS
   brew install redis
   brew services start redis
   
   # Ubuntu/Debian
   sudo apt-get install redis-server
   sudo systemctl start redis-server
   
   # Windows
   # Download and install Redis: https://github.com/microsoftarchive/redis/releases
   ```

2. **Check Redis Status**
   ```bash
   redis-cli ping
   # Should return: PONG
   ```

3. **Configure Connection Address**
   - Local default: `localhost:6379`
   - Remote server: `your-redis-host:6379`

**Example**:
```bash
REDIS_URL=localhost:6379
```

### REDIS_PASSWORD

**Configuration Item**: `REDIS_PASSWORD`

**Setup Steps**:

1. **Set Redis Password** (Optional, recommended for production)
   ```bash
   # Edit Redis configuration file
   # macOS: /usr/local/etc/redis.conf
   # Linux: /etc/redis/redis.conf
   
   # Find and modify:
   # requirepass your_redis_password
   
   # Restart Redis
   # macOS
   brew services restart redis
   
   # Linux
   sudo systemctl restart redis-server
   ```

2. **Configure Password**
   ```bash
   REDIS_PASSWORD=your_redis_password
   ```

**Note**: If no password is set, leave it empty:
```bash
REDIS_PASSWORD=
```

---

## JWT Secret Key

### JWT_SECRET

**Configuration Item**: `JWT_SECRET`

**Setup Steps**:

1. **Generate Random Secret** (Recommended: 32+ character random string)
   ```bash
   # Method 1: Using openssl
   openssl rand -base64 32
   
   # Method 2: Using Python
   python3 -c "import secrets; print(secrets.token_urlsafe(32))"
   
   # Method 3: Using Node.js
   node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
   ```

2. **Configure Secret**
   ```bash
   JWT_SECRET=generated_secret_string
   ```

**Security Tips**:
- Production environment must use strong random secrets
- Do not commit secrets to version control
- Rotate secrets regularly

**Example**:
```bash
JWT_SECRET=your-secret-key-change-in-production-minimum-32-characters
```

---

## Solana Configuration

### SOLANA_RPC_URL

**Configuration Item**: `SOLANA_RPC_URL`

**Setup Steps**:

#### Development Environment (Local Validator)
1. **Install Solana CLI**
   ```bash
   sh -c "$(curl -sSfL https://release.solana.com/stable/install)"
   ```

2. **Start Local Validator**
   ```bash
   solana-test-validator
   ```

3. **Configure RPC URL**
   ```bash
   SOLANA_RPC_URL=http://localhost:8899
   ```

#### Test Environment (Devnet)
```bash
SOLANA_RPC_URL=https://api.devnet.solana.com
```

#### Mainnet Environment
```bash
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
```

#### Using Third-Party RPC Providers (Recommended for Production)
1. **Register Account**
   - QuickNode: https://www.quicknode.com/
   - Alchemy: https://www.alchemy.com/
   - Helius: https://www.helius.dev/

2. **Create RPC Endpoint**
   - Login and create new Solana RPC endpoint
   - Select network (Devnet/Mainnet)
   - Get RPC URL (usually includes API Key)

3. **Configure RPC URL**
   ```bash
   SOLANA_RPC_URL=https://your-api-key.solana-mainnet.quiknode.pro/your-endpoint-id/
   ```

### SOLANA_WS_URL

**Configuration Item**: `SOLANA_WS_URL`

**Setup Steps**:

#### Development Environment (Local Validator)
```bash
SOLANA_WS_URL=ws://localhost:8900
```

#### Test Environment (Devnet)
```bash
SOLANA_WS_URL=wss://api.devnet.solana.com
```

#### Mainnet Environment
```bash
SOLANA_WS_URL=wss://api.mainnet-beta.solana.com
```

#### Using Third-Party WebSocket Service
- Use the same provider as RPC URL
- WebSocket URL is usually provided by the same service provider

### SOLANA_ADMIN_PRIVKEY

**Configuration Item**: `SOLANA_ADMIN_PRIVKEY`

**Setup Steps**:

1. **Generate New Solana Keypair**
   ```bash
   solana-keygen new --outfile ~/.config/solana/admin-keypair.json
   ```

2. **Export Private Key (Base58 Format)**
   ```bash
   solana-keygen pubkey ~/.config/solana/admin-keypair.json
   # Get public key address
   
   # Export private key (Base58)
   cat ~/.config/solana/admin-keypair.json | jq -r '.[]' | base58
   # Or use Solana CLI
   solana-keygen new --no-bip39-passphrase --outfile /tmp/temp-keypair.json
   # Then read the array from JSON file and convert to Base58
   ```

3. **Using Go Code to Export** (Recommended)
   Create temporary script `export_key.go`:
   ```go
   package main
   
   import (
       "encoding/json"
       "fmt"
       "io/ioutil"
       "github.com/mr-tron/base58"
   )
   
   func main() {
       data, _ := ioutil.ReadFile("~/.config/solana/admin-keypair.json")
       var keypair []uint8
       json.Unmarshal(data, &keypair)
       fmt.Println(base58.Encode(keypair))
   }
   ```

4. **Configure Private Key**
   ```bash
   SOLANA_ADMIN_PRIVKEY=exported_base58_private_key_string
   ```

**Security Warning**:
- ⚠️ **Never commit private keys to version control**
- Use environment variables or key management services
- Production environment should use hardware wallets or multisig wallets

### SOLANA_FEE_PAYER

**Configuration Item**: `SOLANA_FEE_PAYER`

**Setup Steps**:

1. **Use Admin Account Public Key**
   ```bash
   solana-keygen pubkey ~/.config/solana/admin-keypair.json
   ```

2. **Or Create Dedicated Fee Payer Account**
   ```bash
   solana-keygen new --outfile ~/.config/solana/fee-payer-keypair.json
   solana-keygen pubkey ~/.config/solana/fee-payer-keypair.json
   ```

3. **Fund Account with SOL** (Test Environment)
   ```bash
   # Devnet
   solana airdrop 2 $(solana-keygen pubkey ~/.config/solana/fee-payer-keypair.json) --url devnet
   ```

4. **Configure Fee Payer Address**
   ```bash
   SOLANA_FEE_PAYER=your_solana_address_base58_format
   ```

### SOLGREEN_TOKEN_MINT

**Configuration Item**: `SOLGREEN_TOKEN_MINT`

**Setup Steps**:

1. **Deploy SPL Token Program** (If not already deployed)
   ```bash
   # Create token using Solana CLI
   spl-token create-token --url devnet
   # Output: Creating token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   ```

2. **Or Get from Anchor Contract Deployment**
   ```bash
   cd contracts/sol-green
   anchor build
   anchor deploy --provider.cluster devnet
   # Program ID will be output after deployment
   ```

3. **Configure Token Mint Address**
   ```bash
   SOLGREEN_TOKEN_MINT=your_token_mint_address
   ```

### SOLGREEN_TOKEN_ADMIN_ATA

**Configuration Item**: `SOLGREEN_TOKEN_ADMIN_ATA`

**Setup Steps**:

1. **Create Associated Token Account (ATA)**
   ```bash
   # Create ATA using admin address
   spl-token create-account <TOKEN_MINT_ADDRESS> --owner <ADMIN_PUBKEY> --url devnet
   ```

2. **Or Calculate ATA Address Using Code**
   ```bash
   # ATA address can be calculated as Program Derived Address (PDA)
   # Using Solana CLI
   spl-token accounts --owner <ADMIN_PUBKEY> --url devnet
   ```

3. **Configure ATA Address**
   ```bash
   SOLGREEN_TOKEN_ADMIN_ATA=your_admin_token_account_address
   ```

### SOLANA_PROOF_CONTRACT

**Configuration Item**: `SOLANA_PROOF_CONTRACT`

**Setup Steps**:

1. **Deploy Sol-Green Smart Contract**
   ```bash
   cd contracts/sol-green
   anchor build
   anchor deploy --provider.cluster devnet
   ```

2. **Get Program ID**
   - Program ID will be displayed in terminal after successful deployment
   - Or check in `contracts/sol-green/target/deploy/sol_green-keypair.json`

3. **Configure Contract Address**
   ```bash
   SOLANA_PROOF_CONTRACT=your_program_id_contract_address
   ```

### SOLANA_PROOF_ACCOUNT

**Configuration Item**: `SOLANA_PROOF_ACCOUNT`

**Setup Steps**:

1. **Initialize Contract Account**
   ```bash
   # Initialize using Anchor client or CLI
   anchor run initialize
   ```

2. **Or Get PDA Through Code**
   - PDA (Program Derived Address) in contract can be calculated from seeds and program ID
   - Check `initialize` function in contract code

3. **Configure Proof Account Address**
   ```bash
   SOLANA_PROOF_ACCOUNT=your_proof_account_address_pda
   ```

---

## Reward Configuration

### REWARD_WASTE_SORTING

**Configuration Item**: `REWARD_WASTE_SORTING`

**Description**: Token reward amount for waste sorting behavior

**Setup Steps**:
Set value directly (unit: token minimum unit, e.g., 1000 = 1000 tokens)

**Example**:
```bash
REWARD_WASTE_SORTING=1000
```

### REWARD_TREE_PLANTING

**Configuration Item**: `REWARD_TREE_PLANTING`

**Description**: Token reward amount for tree planting behavior

**Example**:
```bash
REWARD_TREE_PLANTING=5000
```

### REWARD_LOW_CARBON_TRAVEL

**Configuration Item**: `REWARD_LOW_CARBON_TRAVEL`

**Description**: Token reward amount for low-carbon travel behavior

**Example**:
```bash
REWARD_LOW_CARBON_TRAVEL=2000
```

### CHALLENGE_BONUS_RATE

**Configuration Item**: `CHALLENGE_BONUS_RATE`

**Description**: Reward multiplier for challenge activities

**Example**:
```bash
CHALLENGE_BONUS_RATE=1.5  # Means 1.5x reward
```

---

## AI Fraud Detection Configuration

### AI_FRAUD_USE_REAL_API

**Configuration Item**: `AI_FRAUD_USE_REAL_API`

**Valid Values**: `true` | `false`

**Description**: 
- `false`: Use mock mode (recommended for development/testing)
- `true`: Use real API (for production)

**Configuration**:
```bash
AI_FRAUD_USE_REAL_API=false  # Development environment
AI_FRAUD_USE_REAL_API=true   # Production environment
```

### AI_FRAUD_PROVIDER

**Configuration Item**: `AI_FRAUD_PROVIDER`

**Valid Values**: `baidu` | `aliyun` | `aws` | `google`

**Configuration**:
```bash
AI_FRAUD_PROVIDER=baidu
```

### AI_FRAUD_THRESHOLD

**Configuration Item**: `AI_FRAUD_THRESHOLD`

**Description**: Fraud detection threshold (0-1), values above this are considered suspicious

**Configuration**:
```bash
AI_FRAUD_THRESHOLD=0.7
```

### Baidu AI Configuration

#### AI_FRAUD_API_KEY & AI_FRAUD_API_SECRET

**Setup Steps**:

1. **Register Baidu AI Open Platform Account**
   - Visit: https://ai.baidu.com/
   - Register and login

2. **Create Application**
   - Go to console: https://console.bce.baidu.com/ai/
   - Select "Image Recognition" or "Content Moderation"
   - Create new application

3. **Get API Key and Secret Key**
   - Find `API Key` and `Secret Key` in application details page
   - Copy and save

4. **Configure**
   ```bash
   AI_FRAUD_API_KEY=your_baidu_api_key
   AI_FRAUD_API_SECRET=your_baidu_secret_key
   ```

### Alibaba Cloud Content Security Configuration

#### ALIYUN_ACCESS_KEY_ID & ALIYUN_ACCESS_KEY_SECRET

**Setup Steps**:

1. **Register Alibaba Cloud Account**
   - Visit: https://www.aliyun.com/
   - Register and complete real-name verification

2. **Enable Content Security Service**
   - Login to console: https://home.console.aliyun.com/
   - Search for "Content Security" or "Green"
   - Enable service

3. **Create AccessKey**
   - Go to "Access Control" -> "Users" -> "Create User"
   - Create user and grant content security permissions
   - Create AccessKey, save `AccessKey ID` and `AccessKey Secret`

4. **Configure**
   ```bash
   ALIYUN_ACCESS_KEY_ID=your_access_key_id
   ALIYUN_ACCESS_KEY_SECRET=your_access_key_secret
   ```

### AWS Rekognition Configuration

#### AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY & AWS_REGION

**Setup Steps**:

1. **Register AWS Account**
   - Visit: https://aws.amazon.com/
   - Register account (credit card required)

2. **Create IAM User**
   - Login to AWS Console
   - Go to IAM service
   - Create new user, grant `AmazonRekognitionFullAccess` permission

3. **Create Access Key**
   - In user details page, create access key
   - Download and save `Access Key ID` and `Secret Access Key`

4. **Select Region**
   - Choose region where Rekognition service is available (e.g., `us-east-1`, `ap-southeast-1`)

5. **Configure**
   ```bash
   AWS_ACCESS_KEY_ID=your_aws_access_key_id
   AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
   AWS_REGION=us-east-1
   ```

### Google Cloud Vision Configuration

#### GOOGLE_CLOUD_PROJECT_ID & GOOGLE_APPLICATION_CREDENTIALS

**Setup Steps**:

1. **Register Google Cloud Account**
   - Visit: https://cloud.google.com/
   - Register and create project

2. **Enable Vision API**
   - Enable Vision API in Google Cloud Console
   - Create service account

3. **Create Service Account Key**
   - Go to "IAM & Admin" -> "Service Accounts"
   - Create service account, grant Vision API permissions
   - Create JSON key file and download

4. **Configure Project ID and Credentials Path**
   ```bash
   GOOGLE_CLOUD_PROJECT_ID=your_project_id
   GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
   ```

---

## Third-Party Verification Configuration

### Alashan SEE Ecological Association

#### PARTNER_ALASHAN_ENABLED

**Configuration Item**: `PARTNER_ALASHAN_ENABLED`

**Valid Values**: `true` | `false`

**Setup Steps**:
1. Contact Alashan SEE Ecological Association to obtain API access
2. Get API endpoint and authentication information

**Configuration**:
```bash
PARTNER_ALASHAN_ENABLED=false  # Set to false when not integrated
PARTNER_ALASHAN_API=https://api.alashansee.org/v1/verify
PARTNER_ALASHAN_API_KEY=your_api_key
PARTNER_ALASHAN_API_SECRET=your_api_secret
```

### Government Environmental Department

#### PARTNER_GOV_ENABLED

**Configuration Item**: `PARTNER_GOV_ENABLED`

**Setup Steps**:
1. Contact local environmental department to obtain API interface
2. Apply for API access and authentication information

**Configuration**:
```bash
PARTNER_GOV_ENABLED=false
PARTNER_GOV_API=https://api.gov-environment.gov.cn/v1/verify
PARTNER_GOV_API_KEY=your_api_key
PARTNER_GOV_API_SECRET=your_api_secret
```

### United Nations Environment Programme (UNEP)

#### PARTNER_UNEP_ENABLED

**Configuration Item**: `PARTNER_UNEP_ENABLED`

**Setup Steps**:
1. Visit UNEP official website to learn about API services
2. Apply for API access

**Configuration**:
```bash
PARTNER_UNEP_ENABLED=false
PARTNER_UNEP_API=https://api.unep.org/v1/verify
PARTNER_UNEP_API_KEY=your_api_key
PARTNER_UNEP_API_SECRET=your_api_secret
```

### World Wide Fund for Nature (WWF)

#### PARTNER_WWF_ENABLED

**Configuration Item**: `PARTNER_WWF_ENABLED`

**Setup Steps**:
1. Contact WWF to obtain partner API access
2. Get authentication information

**Configuration**:
```bash
PARTNER_WWF_ENABLED=false
PARTNER_WWF_API=https://api.wwf.org/v1/verify
PARTNER_WWF_API_KEY=your_api_key
PARTNER_WWF_API_SECRET=your_api_secret
```

### Greenpeace

#### PARTNER_GREENPEACE_ENABLED

**Configuration Item**: `PARTNER_GREENPEACE_ENABLED`

**Setup Steps**:
1. Contact Greenpeace to obtain API access
2. Get authentication information

**Configuration**:
```bash
PARTNER_GREENPEACE_ENABLED=false
PARTNER_GREENPEACE_API=https://api.greenpeace.org/v1/verify
PARTNER_GREENPEACE_API_KEY=your_api_key
PARTNER_GREENPEACE_API_SECRET=your_api_secret
```

**Note**: The API endpoints and authentication methods for the above third-party organizations are examples. In actual use, you need to contact the corresponding organizations to obtain real API documentation and access permissions.

---

## Rate Limiting Configuration

### RATE_LIMIT_REQUESTS

**Configuration Item**: `RATE_LIMIT_REQUESTS`

**Description**: Number of requests allowed per minute

**Configuration**:
```bash
RATE_LIMIT_REQUESTS=60  # 60 requests per minute
```

### RATE_LIMIT_BURST

**Configuration Item**: `RATE_LIMIT_BURST`

**Description**: Maximum number of burst requests allowed

**Configuration**:
```bash
RATE_LIMIT_BURST=10  # Allow burst of 10 requests
```

---

## Multi-Chain Configuration

### DEFAULT_CHAIN

**Configuration Item**: `DEFAULT_CHAIN`

**Valid Values**: `solana` | `ethereum` | `polygon` | `bsc`

**Configuration**:
```bash
DEFAULT_CHAIN=solana
```

### Ethereum Configuration

#### ETHEREUM_RPC_URL

**Setup Steps**:

1. **Use Public RPC** (Testing)
   ```bash
   ETHEREUM_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY
   ```

2. **Register Third-Party Service** (Recommended)
   - Infura: https://infura.io/
   - Alchemy: https://www.alchemy.com/
   - QuickNode: https://www.quicknode.com/
   - Create project and get RPC URL

3. **Configure**
   ```bash
   ETHEREUM_RPC_URL=your_ethereum_rpc_url
   ```

#### ETHEREUM_PRIVATE_KEY

**Setup Steps**:

1. **Create Ethereum Wallet**
   ```bash
   # Use MetaMask or command-line tools to create
   # Export private key (be careful with security)
   ```

2. **Configure Private Key** (Hex format, starts with 0x)
   ```bash
   ETHEREUM_PRIVATE_KEY=0xyour_private_key
   ```

#### ETHEREUM_TOKEN_ADDRESS

**Setup Steps**:

1. **Deploy ERC-20 Token Contract** (If not already deployed)
2. **Get Token Contract Address**

**Configuration**:
```bash
ETHEREUM_TOKEN_ADDRESS=your_erc20_token_contract_address
```

### Polygon Configuration

#### POLYGON_RPC_URL

**Setup Steps**:

1. **Use Public RPC** (Testing)
   ```bash
   POLYGON_RPC_URL=https://polygon-rpc.com
   ```

2. **Register Third-Party Service** (Recommended)
   - Alchemy: https://www.alchemy.com/
   - QuickNode: https://www.quicknode.com/
   - Infura: https://infura.io/

**Configuration**:
```bash
POLYGON_RPC_URL=your_polygon_rpc_url
```

#### POLYGON_PRIVATE_KEY

**Setup Steps**: Same as Ethereum, create Polygon wallet and export private key

**Configuration**:
```bash
POLYGON_PRIVATE_KEY=0xyour_private_key
```

#### POLYGON_TOKEN_ADDRESS

**Setup Steps**: Deploy ERC-20 token contract on Polygon

**Configuration**:
```bash
POLYGON_TOKEN_ADDRESS=your_polygon_token_contract_address
```

### BSC Configuration

#### BSC_RPC_URL

**Setup Steps**:

1. **Use Public RPC** (Testing)
   ```bash
   BSC_RPC_URL=https://bsc-dataseed.binance.org
   ```

2. **Register Third-Party Service** (Recommended)
   - QuickNode: https://www.quicknode.com/
   - Ankr: https://www.ankr.com/

**Configuration**:
```bash
BSC_RPC_URL=your_bsc_rpc_url
```

#### BSC_PRIVATE_KEY

**Setup Steps**: Create BSC wallet and export private key

**Configuration**:
```bash
BSC_PRIVATE_KEY=0xyour_private_key
```

#### BSC_TOKEN_ADDRESS

**Setup Steps**: Deploy BEP-20 token contract on BSC

**Configuration**:
```bash
BSC_TOKEN_ADDRESS=your_bsc_token_contract_address
```

---

## Contract Upgrade Configuration

### LATEST_PROGRAM_VERSION

**Configuration Item**: `LATEST_PROGRAM_VERSION`

**Description**: Latest version number of current smart contract

**Configuration**:
```bash
LATEST_PROGRAM_VERSION=1.0.0
```

**Update Steps**:
1. Update this version number after contract upgrade
2. Follow Semantic Versioning

---

## Frontend API Address

### REACT_APP_API_URL

**Configuration Item**: `REACT_APP_API_URL`

**Description**: Backend API address for frontend application connection

**Setup Steps**:

1. **Development Environment**
   ```bash
   REACT_APP_API_URL=http://localhost:8080
   ```

2. **Test Environment**
   ```bash
   REACT_APP_API_URL=https://api-test.solgreen.com
   ```

3. **Production Environment**
   ```bash
   REACT_APP_API_URL=https://api.solgreen.com
   ```

**Note**: 
- This configuration item needs to be set in `frontend/.env` file
- React applications can only access environment variables prefixed with `REACT_APP_` in frontend code

---

## Quick Configuration Checklist

### Required Configuration (Minimum to Run)

- [ ] `ENVIRONMENT` - Environment type
- [ ] `DB_TYPE` - Database type (SQLite needs no additional config)
- [ ] `REDIS_URL` - Redis address
- [ ] `JWT_SECRET` - JWT secret key
- [ ] `SOLANA_RPC_URL` - Solana RPC address
- [ ] `SOLANA_ADMIN_PRIVKEY` - Solana admin private key
- [ ] `SOLANA_FEE_PAYER` - Fee payer address
- [ ] `SOLGREEN_TOKEN_MINT` - Token mint address
- [ ] `SOLANA_PROOF_CONTRACT` - Contract address

### Optional Configuration (Feature Enhancement)

- [ ] AI fraud detection configuration (can skip in development)
- [ ] Third-party verification configuration (can skip in development)
- [ ] Multi-chain configuration (can skip if only using Solana)
- [ ] PostgreSQL database (can skip if using SQLite)

---

## Security Tips

1. **Private Key Security**
   - ⚠️ Never commit private keys to version control
   - Use environment variables or key management services (e.g., AWS Secrets Manager, HashiCorp Vault)
   - Production environment should use hardware wallets or multisig wallets

2. **API Key Security**
   - Do not hardcode API keys in code
   - Use `.env` file (already added to `.gitignore`)
   - Rotate keys regularly

3. **Database Security**
   - Use strong passwords in production
   - Enable SSL/TLS connections
   - Restrict database access IP

4. **Redis Security**
   - Set password in production
   - Restrict Redis access IP
   - Use Redis AUTH

---

## Frequently Asked Questions

### Q: How to quickly start development?

A: Use minimal configuration:
- `ENVIRONMENT=dev`
- `DB_TYPE=sqlite`
- `REDIS_URL=localhost:6379`
- `JWT_SECRET` (randomly generated)
- `SOLANA_RPC_URL=http://localhost:8899` (local validator)
- `AI_FRAUD_USE_REAL_API=false` (use mock mode)

### Q: What is the Solana private key format?

A: Solana private key is a Base58-encoded 64-byte array (JSON format) or direct Base58 string.

### Q: How to test multi-chain functionality?

A: Use test networks:
- Ethereum: Sepolia Testnet
- Polygon: Mumbai Testnet
- BSC: BSC Testnet

### Q: How to obtain third-party verification organization API?

A: Need to contact corresponding organizations to apply for API access. During development, you can set `PARTNER_*_ENABLED=false` to skip.

---

## Related Documentation

- [Quick Start Guide](./QUICKSTART_v3.0.0_EN.md)
- [Development Documentation](./DEVELOPMENT_v3.0.0_EN.md)
- [Deployment Documentation](./DEPLOYMENT_v3.0.0_EN.md)
- [Technical Documentation](./TECHNICAL_v3.0.0_EN.md)

---

**Document Version**: v3.0.0  
**Last Updated**: 2026-01-24
