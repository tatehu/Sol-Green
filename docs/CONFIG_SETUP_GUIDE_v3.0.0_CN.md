# Sol-Green 配置获取指南 v3.0.0

本文档详细说明如何获取 `.env.example` 文件中各项配置的值。

## 目录

1. [环境类型配置](#环境类型配置)
2. [数据库配置](#数据库配置)
3. [Redis 配置](#redis-配置)
4. [JWT 密钥](#jwt-密钥)
5. [Solana 配置](#solana-配置)
6. [奖励配置](#奖励配置)
7. [AI 反欺诈检测配置](#ai-反欺诈检测配置)
8. [第三方认证机构配置](#第三方认证机构配置)
9. [限流配置](#限流配置)
10. [多链配置](#多链配置)
11. [合约升级配置](#合约升级配置)
12. [前端 API 地址](#前端-api-地址)

---

## 环境类型配置

### ENVIRONMENT

**配置项**: `ENVIRONMENT`

**可选值**:
- `dev`: 开发环境（本地开发）
- `test`: 测试环境（Solana Devnet）
- `mainnet`: 主网环境（Solana Mainnet）

**获取步骤**:
1. 开发阶段选择 `dev` 或 `test`
2. 生产部署选择 `mainnet`

**示例**:
```bash
ENVIRONMENT=test
```

---

## 数据库配置

### DB_TYPE

**配置项**: `DB_TYPE`

**可选值**:
- `sqlite`: SQLite 数据库（默认，无需额外配置）
- `postgres`: PostgreSQL 数据库

**获取步骤**:

#### SQLite（推荐用于开发）
无需额外配置，直接使用：
```bash
DB_TYPE=sqlite
```

#### PostgreSQL（生产环境推荐）
1. **安装 PostgreSQL**
   ```bash
   # macOS
   brew install postgresql
   brew services start postgresql
   
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   sudo systemctl start postgresql
   
   # Windows
   # 下载并安装 PostgreSQL: https://www.postgresql.org/download/windows/
   ```

2. **创建数据库和用户**
   ```bash
   # 登录 PostgreSQL
   psql postgres
   
   # 创建数据库
   CREATE DATABASE solgreen;
   
   # 创建用户（可选）
   CREATE USER solgreen_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE solgreen TO solgreen_user;
   ```

3. **配置连接字符串**
   ```bash
   DB_TYPE=postgres
   DATABASE_URL=postgres://solgreen_user:your_password@localhost:5432/solgreen?sslmode=disable
   ```

---

## Redis 配置

### REDIS_URL

**配置项**: `REDIS_URL`

**格式**: `host:port`

**获取步骤**:

1. **安装 Redis**
   ```bash
   # macOS
   brew install redis
   brew services start redis
   
   # Ubuntu/Debian
   sudo apt-get install redis-server
   sudo systemctl start redis-server
   
   # Windows
   # 下载并安装 Redis: https://github.com/microsoftarchive/redis/releases
   ```

2. **检查 Redis 运行状态**
   ```bash
   redis-cli ping
   # 应返回: PONG
   ```

3. **配置连接地址**
   - 本地默认: `localhost:6379`
   - 远程服务器: `your-redis-host:6379`

**示例**:
```bash
REDIS_URL=localhost:6379
```

### REDIS_PASSWORD

**配置项**: `REDIS_PASSWORD`

**获取步骤**:

1. **设置 Redis 密码**（可选，生产环境推荐）
   ```bash
   # 编辑 Redis 配置文件
   # macOS: /usr/local/etc/redis.conf
   # Linux: /etc/redis/redis.conf
   
   # 找到并修改:
   # requirepass your_redis_password
   
   # 重启 Redis
   # macOS
   brew services restart redis
   
   # Linux
   sudo systemctl restart redis-server
   ```

2. **配置密码**
   ```bash
   REDIS_PASSWORD=your_redis_password
   ```

**注意**: 如果未设置密码，留空即可：
```bash
REDIS_PASSWORD=
```

---

## JWT 密钥

### JWT_SECRET

**配置项**: `JWT_SECRET`

**获取步骤**:

1. **生成随机密钥**（推荐使用 32 字符以上的随机字符串）
   ```bash
   # 方法 1: 使用 openssl
   openssl rand -base64 32
   
   # 方法 2: 使用 Python
   python3 -c "import secrets; print(secrets.token_urlsafe(32))"
   
   # 方法 3: 使用 Node.js
   node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
   ```

2. **配置密钥**
   ```bash
   JWT_SECRET=生成的密钥字符串
   ```

**安全提示**:
- 生产环境必须使用强随机密钥
- 不要将密钥提交到版本控制系统
- 定期轮换密钥

**示例**:
```bash
JWT_SECRET=your-secret-key-change-in-production-至少32字符
```

---

## Solana 配置

### SOLANA_RPC_URL

**配置项**: `SOLANA_RPC_URL`

**获取步骤**:

#### 开发环境（本地验证器）
1. **安装 Solana CLI**
   ```bash
   sh -c "$(curl -sSfL https://release.solana.com/stable/install)"
   ```

2. **启动本地验证器**
   ```bash
   solana-test-validator
   ```

3. **配置 RPC URL**
   ```bash
   SOLANA_RPC_URL=http://localhost:8899
   ```

#### 测试环境（Devnet）
```bash
SOLANA_RPC_URL=https://api.devnet.solana.com
```

#### 主网环境（Mainnet）
```bash
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
```

#### 使用第三方 RPC 服务商（推荐生产环境）
1. **注册账号**
   - QuickNode: https://www.quicknode.com/
   - Alchemy: https://www.alchemy.com/
   - Helius: https://www.helius.dev/

2. **创建 RPC 端点**
   - 登录后创建新的 Solana RPC 端点
   - 选择网络（Devnet/Mainnet）
   - 获取 RPC URL（通常包含 API Key）

3. **配置 RPC URL**
   ```bash
   SOLANA_RPC_URL=https://your-api-key.solana-mainnet.quiknode.pro/your-endpoint-id/
   ```

### SOLANA_WS_URL

**配置项**: `SOLANA_WS_URL`

**获取步骤**:

#### 开发环境（本地验证器）
```bash
SOLANA_WS_URL=ws://localhost:8900
```

#### 测试环境（Devnet）
```bash
SOLANA_WS_URL=wss://api.devnet.solana.com
```

#### 主网环境（Mainnet）
```bash
SOLANA_WS_URL=wss://api.mainnet-beta.solana.com
```

#### 使用第三方 WebSocket 服务
- 使用与 RPC URL 相同的服务商
- WebSocket URL 通常与 RPC URL 在同一服务商处提供

### SOLANA_ADMIN_PRIVKEY

**配置项**: `SOLANA_ADMIN_PRIVKEY`

**获取步骤**:

1. **生成新的 Solana 密钥对**
   ```bash
   solana-keygen new --outfile ~/.config/solana/admin-keypair.json
   ```

2. **导出私钥（Base58 格式）**
   ```bash
   solana-keygen pubkey ~/.config/solana/admin-keypair.json
   # 获取公钥地址
   
   # 导出私钥（Base58）
   cat ~/.config/solana/admin-keypair.json | jq -r '.[]' | base58
   # 或使用 Solana CLI
   solana-keygen new --no-bip39-passphrase --outfile /tmp/temp-keypair.json
   # 然后读取 JSON 文件中的数组，转换为 Base58
   ```

3. **使用 Go 代码导出**（推荐）
   创建临时脚本 `export_key.go`:
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

4. **配置私钥**
   ```bash
   SOLANA_ADMIN_PRIVKEY=导出的Base58私钥字符串
   ```

**安全提示**:
- ⚠️ **绝不要将私钥提交到版本控制系统**
- 使用环境变量或密钥管理服务
- 生产环境使用硬件钱包或多签钱包

### SOLANA_FEE_PAYER

**配置项**: `SOLANA_FEE_PAYER`

**获取步骤**:

1. **使用管理员账户的公钥**
   ```bash
   solana-keygen pubkey ~/.config/solana/admin-keypair.json
   ```

2. **或创建专门的费用支付账户**
   ```bash
   solana-keygen new --outfile ~/.config/solana/fee-payer-keypair.json
   solana-keygen pubkey ~/.config/solana/fee-payer-keypair.json
   ```

3. **为账户充值 SOL**（测试环境）
   ```bash
   # Devnet
   solana airdrop 2 $(solana-keygen pubkey ~/.config/solana/fee-payer-keypair.json) --url devnet
   ```

4. **配置费用支付地址**
   ```bash
   SOLANA_FEE_PAYER=你的Solana地址（Base58格式）
   ```

### SOLGREEN_TOKEN_MINT

**配置项**: `SOLGREEN_TOKEN_MINT`

**获取步骤**:

1. **部署 SPL Token 程序**（如果尚未部署）
   ```bash
   # 使用 Solana CLI 创建代币
   spl-token create-token --url devnet
   # 输出: Creating token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   ```

2. **或使用 Anchor 部署合约后获取**
   ```bash
   cd contracts/sol-green
   anchor build
   anchor deploy --provider.cluster devnet
   # 合约部署后会输出程序 ID
   ```

3. **配置代币 Mint 地址**
   ```bash
   SOLGREEN_TOKEN_MINT=你的代币Mint地址
   ```

### SOLGREEN_TOKEN_ADMIN_ATA

**配置项**: `SOLGREEN_TOKEN_ADMIN_ATA`

**获取步骤**:

1. **创建关联代币账户（ATA）**
   ```bash
   # 使用管理员地址创建 ATA
   spl-token create-account <TOKEN_MINT_ADDRESS> --owner <ADMIN_PUBKEY> --url devnet
   ```

2. **或使用代码计算 ATA 地址**
   ```bash
   # ATA 地址可以通过程序派生地址（PDA）计算得出
   # 使用 Solana CLI
   spl-token accounts --owner <ADMIN_PUBKEY> --url devnet
   ```

3. **配置 ATA 地址**
   ```bash
   SOLGREEN_TOKEN_ADMIN_ATA=你的管理员代币账户地址
   ```

### SOLANA_PROOF_CONTRACT

**配置项**: `SOLANA_PROOF_CONTRACT`

**获取步骤**:

1. **部署 Sol-Green 智能合约**
   ```bash
   cd contracts/sol-green
   anchor build
   anchor deploy --provider.cluster devnet
   ```

2. **获取程序 ID**
   - 部署成功后，程序 ID 会显示在终端
   - 或在 `contracts/sol-green/target/deploy/sol_green-keypair.json` 中查看

3. **配置合约地址**
   ```bash
   SOLANA_PROOF_CONTRACT=你的程序ID（合约地址）
   ```

### SOLANA_PROOF_ACCOUNT

**配置项**: `SOLANA_PROOF_ACCOUNT`

**获取步骤**:

1. **初始化合约账户**
   ```bash
   # 使用 Anchor 客户端或 CLI 初始化
   anchor run initialize
   ```

2. **或通过代码获取 PDA**
   - 合约中的 PDA（程序派生地址）可以通过种子和程序 ID 计算得出
   - 查看合约代码中的 `initialize` 函数

3. **配置证明账户地址**
   ```bash
   SOLANA_PROOF_ACCOUNT=你的证明账户地址（PDA）
   ```

---

## 奖励配置

### REWARD_WASTE_SORTING

**配置项**: `REWARD_WASTE_SORTING`

**说明**: 垃圾分类行为的奖励代币数量

**配置步骤**:
直接设置数值（单位：代币最小单位，如 1000 = 1000 tokens）

**示例**:
```bash
REWARD_WASTE_SORTING=1000
```

### REWARD_TREE_PLANTING

**配置项**: `REWARD_TREE_PLANTING`

**说明**: 植树行为的奖励代币数量

**示例**:
```bash
REWARD_TREE_PLANTING=5000
```

### REWARD_LOW_CARBON_TRAVEL

**配置项**: `REWARD_LOW_CARBON_TRAVEL`

**说明**: 低碳出行行为的奖励代币数量

**示例**:
```bash
REWARD_LOW_CARBON_TRAVEL=2000
```

### CHALLENGE_BONUS_RATE

**配置项**: `CHALLENGE_BONUS_RATE`

**说明**: 挑战活动的奖励加成倍数

**示例**:
```bash
CHALLENGE_BONUS_RATE=1.5  # 表示 1.5 倍奖励
```

---

## AI 反欺诈检测配置

### AI_FRAUD_USE_REAL_API

**配置项**: `AI_FRAUD_USE_REAL_API`

**可选值**: `true` | `false`

**说明**: 
- `false`: 使用模拟模式（开发测试推荐）
- `true`: 使用真实 API（生产环境）

**配置**:
```bash
AI_FRAUD_USE_REAL_API=false  # 开发环境
AI_FRAUD_USE_REAL_API=true   # 生产环境
```

### AI_FRAUD_PROVIDER

**配置项**: `AI_FRAUD_PROVIDER`

**可选值**: `baidu` | `aliyun` | `aws` | `google`

**配置**:
```bash
AI_FRAUD_PROVIDER=baidu
```

### AI_FRAUD_THRESHOLD

**配置项**: `AI_FRAUD_THRESHOLD`

**说明**: 欺诈检测阈值（0-1），超过此值判定为可疑

**配置**:
```bash
AI_FRAUD_THRESHOLD=0.7
```

### 百度 AI 配置

#### AI_FRAUD_API_KEY & AI_FRAUD_API_SECRET

**获取步骤**:

1. **注册百度 AI 开放平台账号**
   - 访问: https://ai.baidu.com/
   - 注册并登录

2. **创建应用**
   - 进入控制台: https://console.bce.baidu.com/ai/
   - 选择"图像识别"或"内容审核"
   - 创建新应用

3. **获取 API Key 和 Secret Key**
   - 在应用详情页面找到 `API Key` 和 `Secret Key`
   - 复制保存

4. **配置**
   ```bash
   AI_FRAUD_API_KEY=你的百度API_Key
   AI_FRAUD_API_SECRET=你的百度Secret_Key
   ```

### 阿里云内容安全配置

#### ALIYUN_ACCESS_KEY_ID & ALIYUN_ACCESS_KEY_SECRET

**获取步骤**:

1. **注册阿里云账号**
   - 访问: https://www.aliyun.com/
   - 注册并完成实名认证

2. **开通内容安全服务**
   - 登录控制台: https://home.console.aliyun.com/
   - 搜索"内容安全"或"Green"
   - 开通服务

3. **创建 AccessKey**
   - 进入"访问控制" -> "用户" -> "创建用户"
   - 创建用户并授予内容安全权限
   - 创建 AccessKey，保存 `AccessKey ID` 和 `AccessKey Secret`

4. **配置**
   ```bash
   ALIYUN_ACCESS_KEY_ID=你的AccessKey_ID
   ALIYUN_ACCESS_KEY_SECRET=你的AccessKey_Secret
   ```

### AWS Rekognition 配置

#### AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY & AWS_REGION

**获取步骤**:

1. **注册 AWS 账号**
   - 访问: https://aws.amazon.com/
   - 注册账号（需要信用卡）

2. **创建 IAM 用户**
   - 登录 AWS 控制台
   - 进入 IAM 服务
   - 创建新用户，授予 `AmazonRekognitionFullAccess` 权限

3. **创建访问密钥**
   - 在用户详情页面，创建访问密钥
   - 下载并保存 `Access Key ID` 和 `Secret Access Key`

4. **选择区域**
   - 选择 Rekognition 服务可用的区域（如 `us-east-1`, `ap-southeast-1`）

5. **配置**
   ```bash
   AWS_ACCESS_KEY_ID=你的AWS_Access_Key_ID
   AWS_SECRET_ACCESS_KEY=你的AWS_Secret_Access_Key
   AWS_REGION=us-east-1
   ```

### Google Cloud Vision 配置

#### GOOGLE_CLOUD_PROJECT_ID & GOOGLE_APPLICATION_CREDENTIALS

**获取步骤**:

1. **注册 Google Cloud 账号**
   - 访问: https://cloud.google.com/
   - 注册并创建项目

2. **启用 Vision API**
   - 在 Google Cloud Console 中启用 Vision API
   - 创建服务账号

3. **创建服务账号密钥**
   - 进入"IAM 和管理" -> "服务账号"
   - 创建服务账号，授予 Vision API 权限
   - 创建 JSON 密钥文件并下载

4. **配置项目 ID 和凭证路径**
   ```bash
   GOOGLE_CLOUD_PROJECT_ID=你的项目ID
   GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
   ```

---

## 第三方认证机构配置

### 阿拉善 SEE 生态协会

#### PARTNER_ALASHAN_ENABLED

**配置项**: `PARTNER_ALASHAN_ENABLED`

**可选值**: `true` | `false`

**获取步骤**:
1. 联系阿拉善 SEE 生态协会获取 API 访问权限
2. 获取 API 端点和认证信息

**配置**:
```bash
PARTNER_ALASHAN_ENABLED=false  # 未集成时设为 false
PARTNER_ALASHAN_API=https://api.alashansee.org/v1/verify
PARTNER_ALASHAN_API_KEY=你的API密钥
PARTNER_ALASHAN_API_SECRET=你的API密钥
```

### 政府环保部门

#### PARTNER_GOV_ENABLED

**配置项**: `PARTNER_GOV_ENABLED`

**获取步骤**:
1. 联系当地环保部门获取 API 接口
2. 申请 API 访问权限和认证信息

**配置**:
```bash
PARTNER_GOV_ENABLED=false
PARTNER_GOV_API=https://api.gov-environment.gov.cn/v1/verify
PARTNER_GOV_API_KEY=你的API密钥
PARTNER_GOV_API_SECRET=你的API密钥
```

### 联合国环境规划署 (UNEP)

#### PARTNER_UNEP_ENABLED

**配置项**: `PARTNER_UNEP_ENABLED`

**获取步骤**:
1. 访问 UNEP 官方网站了解 API 服务
2. 申请 API 访问权限

**配置**:
```bash
PARTNER_UNEP_ENABLED=false
PARTNER_UNEP_API=https://api.unep.org/v1/verify
PARTNER_UNEP_API_KEY=你的API密钥
PARTNER_UNEP_API_SECRET=你的API密钥
```

### 世界自然基金会 (WWF)

#### PARTNER_WWF_ENABLED

**配置项**: `PARTNER_WWF_ENABLED`

**获取步骤**:
1. 联系 WWF 获取合作伙伴 API 访问权限
2. 获取认证信息

**配置**:
```bash
PARTNER_WWF_ENABLED=false
PARTNER_WWF_API=https://api.wwf.org/v1/verify
PARTNER_WWF_API_KEY=你的API密钥
PARTNER_WWF_API_SECRET=你的API密钥
```

### 绿色和平 (Greenpeace)

#### PARTNER_GREENPEACE_ENABLED

**配置项**: `PARTNER_GREENPEACE_ENABLED`

**获取步骤**:
1. 联系 Greenpeace 获取 API 访问权限
2. 获取认证信息

**配置**:
```bash
PARTNER_GREENPEACE_ENABLED=false
PARTNER_GREENPEACE_API=https://api.greenpeace.org/v1/verify
PARTNER_GREENPEACE_API_KEY=你的API密钥
PARTNER_GREENPEACE_API_SECRET=你的API密钥
```

**注意**: 以上第三方机构的 API 端点和认证方式为示例，实际使用时需要联系对应机构获取真实的 API 文档和访问权限。

---

## 限流配置

### RATE_LIMIT_REQUESTS

**配置项**: `RATE_LIMIT_REQUESTS`

**说明**: 每分钟允许的请求数

**配置**:
```bash
RATE_LIMIT_REQUESTS=60  # 每分钟 60 个请求
```

### RATE_LIMIT_BURST

**配置项**: `RATE_LIMIT_BURST`

**说明**: 突发请求允许的最大数量

**配置**:
```bash
RATE_LIMIT_BURST=10  # 允许突发 10 个请求
```

---

## 多链配置

### DEFAULT_CHAIN

**配置项**: `DEFAULT_CHAIN`

**可选值**: `solana` | `ethereum` | `polygon` | `bsc`

**配置**:
```bash
DEFAULT_CHAIN=solana
```

### Ethereum 配置

#### ETHEREUM_RPC_URL

**获取步骤**:

1. **使用公共 RPC**（测试）
   ```bash
   ETHEREUM_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY
   ```

2. **注册第三方服务**（推荐）
   - Infura: https://infura.io/
   - Alchemy: https://www.alchemy.com/
   - QuickNode: https://www.quicknode.com/
   - 创建项目并获取 RPC URL

3. **配置**
   ```bash
   ETHEREUM_RPC_URL=你的Ethereum_RPC_URL
   ```

#### ETHEREUM_PRIVATE_KEY

**获取步骤**:

1. **创建 Ethereum 钱包**
   ```bash
   # 使用 MetaMask 或命令行工具创建
   # 导出私钥（注意安全）
   ```

2. **配置私钥**（Hex 格式，0x 开头）
   ```bash
   ETHEREUM_PRIVATE_KEY=0x你的私钥
   ```

#### ETHEREUM_TOKEN_ADDRESS

**获取步骤**:

1. **部署 ERC-20 代币合约**（如果尚未部署）
2. **获取代币合约地址**

**配置**:
```bash
ETHEREUM_TOKEN_ADDRESS=你的ERC20代币合约地址
```

### Polygon 配置

#### POLYGON_RPC_URL

**获取步骤**:

1. **使用公共 RPC**（测试）
   ```bash
   POLYGON_RPC_URL=https://polygon-rpc.com
   ```

2. **注册第三方服务**（推荐）
   - Alchemy: https://www.alchemy.com/
   - QuickNode: https://www.quicknode.com/
   - Infura: https://infura.io/

**配置**:
```bash
POLYGON_RPC_URL=你的Polygon_RPC_URL
```

#### POLYGON_PRIVATE_KEY

**获取步骤**: 同 Ethereum，创建 Polygon 钱包并导出私钥

**配置**:
```bash
POLYGON_PRIVATE_KEY=0x你的私钥
```

#### POLYGON_TOKEN_ADDRESS

**获取步骤**: 部署 Polygon 上的 ERC-20 代币合约

**配置**:
```bash
POLYGON_TOKEN_ADDRESS=你的Polygon代币合约地址
```

### BSC 配置

#### BSC_RPC_URL

**获取步骤**:

1. **使用公共 RPC**（测试）
   ```bash
   BSC_RPC_URL=https://bsc-dataseed.binance.org
   ```

2. **注册第三方服务**（推荐）
   - QuickNode: https://www.quicknode.com/
   - Ankr: https://www.ankr.com/

**配置**:
```bash
BSC_RPC_URL=你的BSC_RPC_URL
```

#### BSC_PRIVATE_KEY

**获取步骤**: 创建 BSC 钱包并导出私钥

**配置**:
```bash
BSC_PRIVATE_KEY=0x你的私钥
```

#### BSC_TOKEN_ADDRESS

**获取步骤**: 部署 BSC 上的 BEP-20 代币合约

**配置**:
```bash
BSC_TOKEN_ADDRESS=你的BSC代币合约地址
```

---

## 合约升级配置

### LATEST_PROGRAM_VERSION

**配置项**: `LATEST_PROGRAM_VERSION`

**说明**: 当前智能合约的最新版本号

**配置**:
```bash
LATEST_PROGRAM_VERSION=1.0.0
```

**更新步骤**:
1. 合约升级后，更新此版本号
2. 遵循语义化版本控制（Semantic Versioning）

---

## 前端 API 地址

### REACT_APP_API_URL

**配置项**: `REACT_APP_API_URL`

**说明**: 前端应用连接的后端 API 地址

**配置步骤**:

1. **开发环境**
   ```bash
   REACT_APP_API_URL=http://localhost:8080
   ```

2. **测试环境**
   ```bash
   REACT_APP_API_URL=https://api-test.solgreen.com
   ```

3. **生产环境**
   ```bash
   REACT_APP_API_URL=https://api.solgreen.com
   ```

**注意**: 
- 此配置项需要在 `frontend/.env` 文件中设置
- React 应用需要以 `REACT_APP_` 前缀开头的环境变量才能在前端代码中访问

---

## 快速配置检查清单

### 必需配置（最小化运行）

- [ ] `ENVIRONMENT` - 环境类型
- [ ] `DB_TYPE` - 数据库类型（SQLite 无需额外配置）
- [ ] `REDIS_URL` - Redis 地址
- [ ] `JWT_SECRET` - JWT 密钥
- [ ] `SOLANA_RPC_URL` - Solana RPC 地址
- [ ] `SOLANA_ADMIN_PRIVKEY` - Solana 管理员私钥
- [ ] `SOLANA_FEE_PAYER` - 费用支付地址
- [ ] `SOLGREEN_TOKEN_MINT` - 代币 Mint 地址
- [ ] `SOLANA_PROOF_CONTRACT` - 合约地址

### 可选配置（功能增强）

- [ ] AI 反欺诈检测配置（开发环境可跳过）
- [ ] 第三方认证机构配置（开发环境可跳过）
- [ ] 多链配置（仅使用 Solana 时可跳过）
- [ ] PostgreSQL 数据库（使用 SQLite 时可跳过）

---

## 安全提示

1. **私钥安全**
   - ⚠️ 绝不要将私钥提交到版本控制系统
   - 使用环境变量或密钥管理服务（如 AWS Secrets Manager, HashiCorp Vault）
   - 生产环境使用硬件钱包或多签钱包

2. **API 密钥安全**
   - 不要在代码中硬编码 API 密钥
   - 使用 `.env` 文件（已添加到 `.gitignore`）
   - 定期轮换密钥

3. **数据库安全**
   - 生产环境使用强密码
   - 启用 SSL/TLS 连接
   - 限制数据库访问 IP

4. **Redis 安全**
   - 生产环境设置密码
   - 限制 Redis 访问 IP
   - 使用 Redis AUTH

---

## 常见问题

### Q: 如何快速开始开发？

A: 使用最小化配置：
- `ENVIRONMENT=dev`
- `DB_TYPE=sqlite`
- `REDIS_URL=localhost:6379`
- `JWT_SECRET`（随机生成）
- `SOLANA_RPC_URL=http://localhost:8899`（本地验证器）
- `AI_FRAUD_USE_REAL_API=false`（使用模拟模式）

### Q: Solana 私钥格式是什么？

A: Solana 私钥是 Base58 编码的 64 字节数组（JSON 格式）或直接 Base58 字符串。

### Q: 如何测试多链功能？

A: 使用测试网络：
- Ethereum: Sepolia Testnet
- Polygon: Mumbai Testnet
- BSC: BSC Testnet

### Q: 第三方认证机构 API 如何获取？

A: 需要联系对应机构申请 API 访问权限。开发阶段可以设置 `PARTNER_*_ENABLED=false` 跳过。

---

## 相关文档

- [快速开始指南](./QUICKSTART_v3.0.0_CN.md)
- [开发文档](./DEVELOPMENT_v3.0.0_CN.md)
- [部署文档](./DEPLOYMENT_v3.0.0_CN.md)
- [技术文档](./TECHNICAL_v3.0.0_CN.md)

---

**文档版本**: v3.0.0  
**最后更新**: 2026-01-24
