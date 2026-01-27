# 部署文档 v3.0.0 / Deployment Documentation v3.0.0

## 部署环境要求 / Deployment Environment Requirements

### 服务器要求 / Server Requirements

- **CPU**: 2 核以上 / 2 cores or more
- **内存**: 4GB 以上 / 4GB RAM or more
- **存储**: 50GB 以上 / 50GB storage or more
- **操作系统**: Linux (Ubuntu 20.04+ 推荐 / recommended)
- **网络**: 稳定的互联网连接 / Stable internet connection

### 软件要求 / Software Requirements

- **Docker**: 20.10+ / Container runtime
- **Docker Compose**: 2.0+ / Container orchestration
- **Nginx**: 前端部署 / Frontend deployment
- **SSL 证书**: HTTPS / SSL certificate for HTTPS
- **Solana CLI**: 合约部署 (可选) / Smart contract deployment (optional)
- **Anchor CLI**: 合约开发 (可选) / Smart contract development (optional)

### 区块链网络要求 / Blockchain Network Requirements

- **Solana RPC 节点**: 测试网或主网 / Testnet or mainnet
- **钱包私钥**: 部署和管理合约 / For contract deployment and management
- **测试 SOL**: 合约部署和交易费用 / For contract deployment and transaction fees

## 部署方式 / Deployment Methods

### 方式一：Docker Compose 部署（推荐）/ Method 1: Docker Compose Deployment (Recommended)

#### 1. 准备服务器 / Prepare Server

```bash
# 更新系统 / Update system
sudo apt update && sudo apt upgrade -y

# 安装 Docker / Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose / Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装 / Verify installation
docker --version
docker-compose --version
```

#### 2. 克隆项目 / Clone Project

```bash
git clone <repository-url>
cd trends_solana
```

#### 3. 配置环境变量 / Configure Environment Variables

```bash
# 复制环境变量模板 / Copy environment template
cp .env.example .env

# 编辑配置文件 / Edit configuration file
nano .env

# 重要配置项 / Important configuration items
ENVIRONMENT=test                    # 环境类型 / Environment type
DB_TYPE=postgres                   # 数据库类型 / Database type
DATABASE_URL=postgres://...        # 数据库连接 / Database connection
SOLANA_RPC_URL=https://api.devnet.solana.com  # Solana RPC
SOLANA_ADMIN_PRIVKEY=your-private-key         # 管理员私钥 / Admin private key
JWT_SECRET=your-jwt-secret         # JWT 密钥 / JWT secret
```

#### 4. 构建和启动 / Build and Start

```bash
# 构建镜像 / Build images
docker-compose build

# 启动服务 / Start services
docker-compose up -d

# 查看服务状态 / Check service status
docker-compose ps

# 查看日志 / View logs
docker-compose logs -f backend
```

#### 5. 验证部署 / Verify Deployment

```bash
# 健康检查 / Health check
curl http://localhost:8080/api/v1/health

# 访问前端 / Access frontend
# http://your-server-ip:3000 (开发模式)
# 或配置 Nginx 反向代理 / Or configure Nginx reverse proxy
```

### 方式二：Kubernetes 部署 / Method 2: Kubernetes Deployment

#### 1. 准备 Kubernetes 集群 / Prepare Kubernetes Cluster

```bash
# 使用 kubeadm 或云服务商提供的 K8s 集群
# Use kubeadm or cloud provider's K8s cluster

# 创建命名空间 / Create namespace
kubectl create namespace sol-green

# 配置存储 / Configure storage
kubectl apply -f k8s/pv-pvc.yaml
```

#### 2. 配置 ConfigMap 和 Secret / Configure ConfigMap and Secret

```bash
# 创建 ConfigMap / Create ConfigMap
kubectl create configmap sol-green-config \
  --from-file=.env \
  -n sol-green

# 创建 Secret / Create Secret
kubectl create secret generic sol-green-secrets \
  --from-literal=jwt-secret=<your-jwt-secret> \
  --from-literal=solana-privkey=<your-private-key> \
  -n sol-green
```

#### 3. 部署应用 / Deploy Application

```bash
# 部署数据库 / Deploy database
kubectl apply -f k8s/postgres-deployment.yaml

# 部署 Redis / Deploy Redis
kubectl apply -f k8s/redis-deployment.yaml

# 部署后端 / Deploy backend
kubectl apply -f k8s/backend-deployment.yaml

# 部署前端 / Deploy frontend
kubectl apply -f k8s/frontend-deployment.yaml

# 检查部署状态 / Check deployment status
kubectl get pods -n sol-green
kubectl get services -n sol-green
```

#### 4. 配置 Ingress / Configure Ingress

```bash
# 安装 Ingress Controller / Install Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# 创建 Ingress 规则 / Create Ingress rules
kubectl apply -f k8s/ingress.yaml
```

### 方式三：传统部署 / Method 3: Traditional Deployment

#### 1. 后端部署 / Backend Deployment

```bash
# 编译二进制文件 / Compile binary
cd backend
go mod download
CGO_ENABLED=0 GOOS=linux go build -o sol-green .

# 创建 systemd 服务 / Create systemd service
sudo nano /etc/systemd/system/sol-green.service

# 服务配置 / Service configuration
[Unit]
Description=Sol-Green Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/sol-green
ExecStart=/opt/sol-green/sol-green
Restart=always
EnvironmentFile=/opt/sol-green/.env

[Install]
WantedBy=multi-user.target

# 启动服务 / Start service
sudo systemctl daemon-reload
sudo systemctl start sol-green
sudo systemctl enable sol-green
```

#### 2. 前端部署 / Frontend Deployment

```bash
# 构建前端 / Build frontend
cd frontend
npm install
npm run build

# 部署到 Nginx / Deploy to Nginx
sudo mkdir -p /var/www/sol-green
sudo cp -r build/* /var/www/sol-green/

# 配置 Nginx / Configure Nginx
sudo nano /etc/nginx/sites-available/sol-green

# Nginx 配置 / Nginx configuration
server {
    listen 80;
    server_name your-domain.com;

    root /var/www/sol-green;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# 启用站点 / Enable site
sudo ln -s /etc/nginx/sites-available/sol-green /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 智能合约部署 / Smart Contract Deployment

### 1. 准备环境 / Prepare Environment

```bash
# 安装 Solana CLI / Install Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/stable/install)"

# 安装 Anchor CLI / Install Anchor CLI
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install latest
avm use latest

# 配置网络 / Configure network
solana config set --url devnet  # 或 mainnet-beta / or mainnet-beta
```

### 2. 创建钱包 / Create Wallet

```bash
# 创建新钱包 / Create new wallet
solana-keygen new --outfile ~/.config/solana/deployer.json

# 配置默认钱包 / Configure default wallet
solana config set --keypair ~/.config/solana/deployer.json

# 获取测试代币 / Get test tokens
solana airdrop 2 $(solana address) --url devnet
```

### 3. 构建和部署合约 / Build and Deploy Contract

```bash
# 进入合约目录 / Enter contract directory
cd contracts/sol-green

# 构建合约 / Build contract
anchor build

# 部署到测试网 / Deploy to testnet
anchor deploy --provider.cluster devnet

# 验证部署 / Verify deployment
solana program show <program-id> --url devnet
```

### 4. 更新配置 / Update Configuration

```bash
# 在 .env 中更新合约地址 / Update contract address in .env
SOLANA_PROOF_CONTRACT=<program-id>
SOLGREEN_TOKEN_MINT=<token-mint-address>
```

## 域名和 SSL 配置 / Domain and SSL Configuration

### 1. 配置域名 / Configure Domain

```bash
# 添加 DNS A 记录 / Add DNS A record
# your-domain.com -> your-server-ip
```

### 2. 安装 SSL 证书 / Install SSL Certificate

```bash
# 安装 Certbot / Install Certbot
sudo apt install certbot python3-certbot-nginx

# 获取 Let's Encrypt 证书 / Get Let's Encrypt certificate
sudo certbot --nginx -d your-domain.com

# 设置自动续期 / Set auto-renewal
sudo crontab -e
# 添加: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 3. 配置 Nginx HTTPS / Configure Nginx HTTPS

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # SSL 安全配置 / SSL security configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    root /var/www/sol-green;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# 重定向 HTTP 到 HTTPS / Redirect HTTP to HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

## 监控和日志 / Monitoring and Logging

### 1. 应用监控 / Application Monitoring

```bash
# 查看 Docker 容器状态 / Check Docker container status
docker-compose ps

# 查看应用日志 / View application logs
docker-compose logs -f backend
docker-compose logs -f frontend

# 查看系统资源使用 / Check system resource usage
docker stats
```

### 2. 区块链监控 / Blockchain Monitoring

```bash
# 查看 Solana 网络状态 / Check Solana network status
solana cluster-version
solana epoch-info

# 查看钱包余额 / Check wallet balance
solana balance

# 查看交易确认 / Check transaction confirmation
solana confirm <transaction-signature>
```

### 3. 设置日志轮转 / Setup Log Rotation

```bash
# 配置 logrotate / Configure logrotate
sudo nano /etc/logrotate.d/sol-green

# 添加配置 / Add configuration
/var/log/sol-green/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    create 644 www-data www-data
    postrotate
        systemctl reload sol-green
    endscript
}
```

## 备份策略 / Backup Strategy

### 1. 数据库备份 / Database Backup

```bash
# PostgreSQL 备份 / PostgreSQL backup
docker-compose exec postgres pg_dump -U solgreen solgreen > backup_$(date +%Y%m%d_%H%M%S).sql

# 自动化备份脚本 / Automated backup script
#!/bin/bash
BACKUP_DIR="/opt/sol-green/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# 数据库备份 / Database backup
docker-compose exec postgres pg_dump -U solgreen solgreen > $BACKUP_DIR/db_$DATE.sql

# 配置文件备份 / Configuration backup
tar -czf $BACKUP_DIR/config_$DATE.tar.gz .env

# 删除7天前的备份 / Delete backups older than 7 days
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### 2. 合约备份 / Contract Backup

```bash
# 导出合约数据 / Export contract data
solana program dump <program-id> contract_backup.so

# 备份私钥 / Backup private keys
cp ~/.config/solana/id.json /secure/backup/wallet_backup.json
```

## 性能优化 / Performance Optimization

### 1. Docker 优化 / Docker Optimization

```yaml
# docker-compose.yml 优化配置 / Optimized docker-compose.yml
version: '3.8'
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

### 2. Nginx 优化 / Nginx Optimization

```nginx
# Nginx 性能配置 / Nginx performance configuration
worker_processes auto;
worker_connections 1024;

# Gzip 压缩 / Gzip compression
gzip on;
gzip_types text/plain text/css application/json application/javascript;

# 缓存配置 / Cache configuration
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

### 3. 数据库优化 / Database Optimization

```sql
-- 创建索引 / Create indexes
CREATE INDEX CONCURRENTLY idx_green_behaviors_wallet_time ON green_behaviors(wallet_addr, submit_time);
CREATE INDEX CONCURRENTLY idx_challenges_status_time ON challenges(status, created_at);

-- 数据库配置优化 / Database configuration optimization
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
```

## 安全配置 / Security Configuration

### 1. 防火墙配置 / Firewall Configuration

```bash
# 配置 UFW 防火墙 / Configure UFW firewall
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw --force enable    # 启用防火墙 / Enable firewall
```

### 2. SSL/TLS 配置 / SSL/TLS Configuration

```nginx
# SSL 配置 / SSL configuration
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
ssl_prefer_server_ciphers off;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
```

### 3. API 安全 / API Security

```bash
# 配置 API 限流 / Configure API rate limiting
RATE_LIMIT_REQUESTS=60
RATE_LIMIT_BURST=10

# JWT 配置 / JWT configuration
JWT_SECRET=your-very-long-random-secret-key
```

## 故障排查 / Troubleshooting

### 常见问题 / Common Issues

#### 容器无法启动 / Container Won't Start
```bash
# 查看容器日志 / Check container logs
docker-compose logs backend

# 检查端口冲突 / Check port conflicts
netstat -tlnp | grep :8080

# 检查磁盘空间 / Check disk space
df -h
```

#### 数据库连接失败 / Database Connection Failed
```bash
# 检查数据库容器 / Check database container
docker-compose ps postgres

# 测试数据库连接 / Test database connection
docker-compose exec postgres psql -U solgreen -d solgreen -c "SELECT 1;"

# 查看数据库日志 / Check database logs
docker-compose logs postgres
```

#### 前端构建失败 / Frontend Build Failed
```bash
# 检查 Node.js 版本 / Check Node.js version
node --version
npm --version

# 清理缓存 / Clean cache
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

#### 智能合约部署失败 / Smart Contract Deployment Failed
```bash
# 检查 Solana 配置 / Check Solana configuration
solana config get

# 检查钱包余额 / Check wallet balance
solana balance

# 检查程序大小 / Check program size
ls -lh target/deploy/*.so
```

### 日志分析 / Log Analysis

```bash
# 查看错误日志 / View error logs
docker-compose logs backend | grep -i error

# 查看性能日志 / View performance logs
docker-compose logs backend | grep -i "took\|latency"

# 查看区块链相关日志 / View blockchain related logs
docker-compose logs backend | grep -i "solana\|blockchain"
```

## 更新部署 / Update Deployment

### 1. 代码更新 / Code Update

```bash
# 拉取最新代码 / Pull latest code
git pull origin main

# 重建镜像 / Rebuild images
docker-compose build --no-cache

# 滚动更新 / Rolling update
docker-compose up -d --scale backend=2
docker-compose up -d --scale backend=1
```

### 2. 零停机部署 / Zero-Downtime Deployment

```bash
# 使用蓝绿部署 / Blue-green deployment
docker-compose up -d backend-green
# 切换流量 / Switch traffic
docker-compose stop backend-blue
docker-compose rm backend-blue
```

### 3. 数据库迁移 / Database Migration

```bash
# 创建迁移脚本 / Create migration script
docker-compose exec backend ./migrate

# 验证迁移 / Verify migration
docker-compose exec postgres psql -U solgreen -d solgreen -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 5;"
```

## 扩展部署 / Scaling Deployment

### 水平扩展 / Horizontal Scaling

```bash
# 扩展后端实例 / Scale backend instances
docker-compose up -d --scale backend=3

# 配置负载均衡 / Configure load balancing
# 使用 Nginx 或云负载均衡器
```

### 垂直扩展 / Vertical Scaling

```bash
# 增加资源限制 / Increase resource limits
docker-compose.yml 更新资源配置 / Update resource configuration in docker-compose.yml

services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
```

### 全球部署 / Global Deployment

```bash
# 使用 CDN / Use CDN
# Cloudflare / AWS CloudFront / Aliyun CDN

# 多区域部署 / Multi-region deployment
# AWS: us-east-1, eu-west-1, ap-southeast-1
# Aliyun: cn-hangzhou, cn-shanghai, cn-shenzhen
```

---

**文档版本 / Document Version**: v3.0.0  
**最后更新 / Last Update**: 2024-01-24  
**兼容版本 / Compatible Version**: Docker 20.10+, Docker Compose 2.0+, Kubernetes 1.24+
