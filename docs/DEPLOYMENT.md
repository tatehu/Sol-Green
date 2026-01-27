# 部署文档

## 部署环境要求

### 服务器要求

- **CPU**: 2 核以上
- **内存**: 4GB 以上
- **存储**: 50GB 以上
- **操作系统**: Linux (Ubuntu 20.04+ 推荐)

### 软件要求

- Docker 20.10+
- Docker Compose 2.0+
- Nginx (前端部署)
- SSL 证书 (HTTPS)

## 部署方式

### 方式一：Docker Compose 部署（推荐）

#### 1. 准备服务器

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. 克隆项目

```bash
git clone <repository-url>
cd trends_solana
```

#### 3. 配置环境变量

```bash
cp .env.example .env
nano .env  # 编辑配置文件
```

重要配置项：
- `JWT_SECRET`: 生成强随机密钥
- `SOLANA_ADMIN_PRIVKEY`: 管理员私钥
- `DATABASE_URL`: 生产数据库连接
- `REDIS_PASSWORD`: Redis 密码

#### 4. 部署

```bash
# 构建并启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 检查服务状态
docker-compose ps
```

#### 5. 验证部署

```bash
# 检查后端健康状态
curl http://localhost:8080/api/v1/health

# 检查数据库连接
docker-compose exec backend go run scripts/check_db.go
```

### 方式二：Kubernetes 部署

#### 1. 准备 Kubernetes 集群

```bash
# 使用 kubeadm 或云服务商提供的 K8s 集群
```

#### 2. 创建命名空间

```bash
kubectl create namespace sol-green
```

#### 3. 创建 ConfigMap 和 Secret

```bash
# 创建 ConfigMap
kubectl create configmap sol-green-config \
  --from-file=.env \
  -n sol-green

# 创建 Secret（敏感信息）
kubectl create secret generic sol-green-secrets \
  --from-literal=jwt-secret=<your-jwt-secret> \
  --from-literal=solana-privkey=<your-private-key> \
  -n sol-green
```

#### 4. 部署应用

```bash
# 应用 Kubernetes 配置
kubectl apply -f k8s/ -n sol-green

# 检查部署状态
kubectl get pods -n sol-green
```

### 方式三：传统部署

#### 1. 后端部署

```bash
# 编译二进制
go build -o sol-green main.go

# 创建 systemd 服务
sudo nano /etc/systemd/system/sol-green.service
```

服务文件内容：
```ini
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
```

```bash
# 启动服务
sudo systemctl start sol-green
sudo systemctl enable sol-green
```

#### 2. 前端部署

```bash
# 构建前端
cd frontend
npm install
npm run build

# 部署到 Nginx
sudo cp -r build/* /var/www/sol-green/
```

Nginx 配置：
```nginx
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
    }
}
```

## 数据库部署

### PostgreSQL 部署

#### Docker 方式

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_USER=solgreen \
  -e POSTGRES_PASSWORD=<password> \
  -e POSTGRES_DB=solgreen \
  -v postgres_data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:15-alpine
```

#### 传统方式

```bash
# 安装 PostgreSQL
sudo apt install postgresql postgresql-contrib

# 创建数据库
sudo -u postgres psql
CREATE DATABASE solgreen;
CREATE USER solgreen WITH PASSWORD '<password>';
GRANT ALL PRIVILEGES ON DATABASE solgreen TO solgreen;
```

### Redis 部署

#### Docker 方式

```bash
docker run -d \
  --name redis \
  -p 6379:6379 \
  -v redis_data:/data \
  redis:7-alpine
```

#### 传统方式

```bash
# 安装 Redis
sudo apt install redis-server

# 配置 Redis
sudo nano /etc/redis/redis.conf
# 设置 requirepass <password>

# 重启服务
sudo systemctl restart redis
```

## 智能合约部署

### 1. 准备环境

```bash
# 安装 Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/stable/install)"

# 安装 Anchor
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install latest
avm use latest
```

### 2. 配置网络

```bash
# 测试网
solana config set --url devnet

# 主网
solana config set --url mainnet
```

### 3. 创建钱包

```bash
solana-keygen new --outfile ~/.config/solana/deployer.json
solana config set --keypair ~/.config/solana/deployer.json
```

### 4. 获取测试币（测试网）

```bash
solana airdrop 2 $(solana address) --url devnet
```

### 5. 构建和部署

```bash
cd contracts/sol-green

# 更新程序 ID
anchor keys list
# 复制程序 ID 到 Anchor.toml 和 lib.rs

# 构建
anchor build

# 部署
anchor deploy --provider.cluster devnet
```

### 6. 验证部署

```bash
# 查看程序信息
solana program show <program-id> --url devnet
```

## 域名和 SSL 配置

### 1. 配置域名

```bash
# 添加 DNS A 记录指向服务器 IP
```

### 2. 安装 SSL 证书（Let's Encrypt）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

### 3. 更新 Nginx 配置

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # ... 其他配置
}
```

## 监控和日志

### 1. 日志收集

```bash
# 使用 Docker 日志
docker-compose logs -f backend

# 或使用 systemd 日志
journalctl -u sol-green -f
```

### 2. 监控设置

推荐使用：
- Prometheus + Grafana
- ELK Stack
- 云服务商监控服务

## 备份策略

### 数据库备份

```bash
# 创建备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U solgreen solgreen > /backup/solgreen_$DATE.sql

# 添加到 crontab（每天凌晨 2 点）
0 2 * * * /path/to/backup.sh
```

### 配置文件备份

```bash
# 备份 .env 和配置文件
tar -czf config_backup_$(date +%Y%m%d).tar.gz .env docker-compose.yml
```

## 更新部署

### 1. 代码更新

```bash
# 拉取最新代码
git pull origin main

# 重新构建
docker-compose build

# 滚动更新
docker-compose up -d
```

### 2. 数据库迁移

```bash
# 数据库迁移在应用启动时自动执行
# 如需手动迁移：
docker-compose exec backend go run scripts/migrate.go
```

## 故障排查

### 常见问题

1. **服务无法启动**
   - 检查端口占用
   - 查看日志
   - 验证环境变量

2. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接字符串
   - 检查防火墙

3. **前端无法访问后端**
   - 检查 CORS 配置
   - 验证 Nginx 代理配置
   - 检查防火墙规则

### 回滚操作

```bash
# Docker Compose 回滚
docker-compose down
docker-compose up -d --scale backend=0
# 恢复之前的镜像
docker-compose up -d
```

## 安全建议

1. **防火墙配置**
   ```bash
   sudo ufw allow 22/tcp
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   ```

2. **定期更新**
   - 系统更新
   - 依赖更新
   - 安全补丁

3. **访问控制**
   - 使用 SSH 密钥认证
   - 限制 SSH 访问 IP
   - 禁用 root 登录

4. **数据加密**
   - 数据库连接使用 SSL
   - 敏感数据加密存储
