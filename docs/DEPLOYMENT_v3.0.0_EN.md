# Deployment Documentation v3.0.0

## Deployment Environment Requirements

### Server Requirements

- **CPU**: 2 cores or more
- **Memory**: 4GB RAM or more
- **Storage**: 50GB storage or more
- **Operating System**: Linux (Ubuntu 20.04+ recommended)
- **Network**: Stable internet connection

### Software Requirements

- **Docker**: 20.10+ (Container runtime)
- **Docker Compose**: 2.0+ (Container orchestration)
- **Nginx**: Frontend deployment
- **SSL Certificate**: For HTTPS
- **Solana CLI**: Smart contract deployment (optional)
- **Anchor CLI**: Smart contract development (optional)

### Blockchain Network Requirements

- **Solana RPC Node**: Testnet or mainnet
- **Wallet Private Key**: For contract deployment and management
- **Test SOL**: For contract deployment and transaction fees

## Deployment Methods

### Method 1: Docker Compose Deployment (Recommended)

#### 1. Prepare Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

#### 2. Clone Project

```bash
git clone <repository-url>
cd trends_solana
```

#### 3. Configure Environment Variables

```bash
# Copy environment template
cp .env.example .env

# Edit configuration file
nano .env

# Important configuration items
ENVIRONMENT=test                    # Environment type
DB_TYPE=postgres                   # Database type
DATABASE_URL=postgres://...        # Database connection
SOLANA_RPC_URL=https://api.devnet.solana.com  # Solana RPC
SOLANA_ADMIN_PRIVKEY=your-private-key         # Admin private key
JWT_SECRET=your-jwt-secret         # JWT secret
```

#### 4. Build and Start

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f backend
```

#### 5. Verify Deployment

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Access frontend
# http://your-server-ip:3000 (development mode)
# Or configure Nginx reverse proxy
```

### Method 2: Kubernetes Deployment

#### 1. Prepare Kubernetes Cluster

```bash
# Use kubeadm or cloud provider's K8s cluster

# Create namespace
kubectl create namespace sol-green

# Configure storage
kubectl apply -f k8s/pv-pvc.yaml
```

#### 2. Configure ConfigMap and Secret

```bash
# Create ConfigMap
kubectl create configmap sol-green-config \
  --from-file=.env \
  -n sol-green

# Create Secret
kubectl create secret generic sol-green-secrets \
  --from-literal=jwt-secret=<your-jwt-secret> \
  --from-literal=solana-privkey=<your-private-key> \
  -n sol-green
```

#### 3. Deploy Application

```bash
# Deploy database
kubectl apply -f k8s/postgres-deployment.yaml

# Deploy Redis
kubectl apply -f k8s/redis-deployment.yaml

# Deploy backend
kubectl apply -f k8s/backend-deployment.yaml

# Deploy frontend
kubectl apply -f k8s/frontend-deployment.yaml

# Check deployment status
kubectl get pods -n sol-green
kubectl get services -n sol-green
```

#### 4. Configure Ingress

```bash
# Install Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Create Ingress rules
kubectl apply -f k8s/ingress.yaml
```

### Method 3: Traditional Deployment

#### 1. Backend Deployment

```bash
# Compile binary
cd backend
go mod download
CGO_ENABLED=0 GOOS=linux go build -o sol-green .

# Create systemd service
sudo nano /etc/systemd/system/sol-green.service

# Service configuration
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

# Start service
sudo systemctl daemon-reload
sudo systemctl start sol-green
sudo systemctl enable sol-green
```

#### 2. Frontend Deployment

```bash
# Build frontend
cd frontend
npm install
npm run build

# Deploy to Nginx
sudo mkdir -p /var/www/sol-green
sudo cp -r build/* /var/www/sol-green/

# Configure Nginx
sudo nano /etc/nginx/sites-available/sol-green

# Nginx configuration
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

# Enable site
sudo ln -s /etc/nginx/sites-available/sol-green /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## Smart Contract Deployment

### 1. Prepare Environment

```bash
# Install Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/stable/install)"

# Install Anchor CLI
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install latest
avm use latest

# Configure network
solana config set --url devnet  # or mainnet-beta
```

### 2. Create Wallet

```bash
# Create new wallet
solana-keygen new --outfile ~/.config/solana/deployer.json

# Configure default wallet
solana config set --keypair ~/.config/solana/deployer.json

# Get test tokens
solana airdrop 2 $(solana address) --url devnet
```

### 3. Build and Deploy Contract

```bash
# Enter contract directory
cd contracts/sol-green

# Build contract
anchor build

# Deploy to testnet
anchor deploy --provider.cluster devnet

# Verify deployment
solana program show <program-id> --url devnet
```

### 4. Update Configuration

```bash
# Update contract address in .env
SOLANA_PROOF_CONTRACT=<program-id>
SOLGREEN_TOKEN_MINT=<token-mint-address>
```

## Domain and SSL Configuration

### 1. Configure Domain

```bash
# Add DNS A record
# your-domain.com -> your-server-ip
```

### 2. Install SSL Certificate

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Get Let's Encrypt certificate
sudo certbot --nginx -d your-domain.com

# Set auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 3. Configure Nginx HTTPS

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # SSL security configuration
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

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

## Monitoring and Logging

### 1. Application Monitoring

```bash
# Check Docker container status
docker-compose ps

# View application logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Check system resource usage
docker stats
```

### 2. Blockchain Monitoring

```bash
# Check Solana network status
solana cluster-version
solana epoch-info

# Check wallet balance
solana balance

# Check transaction confirmation
solana confirm <transaction-signature>
```

### 3. Setup Log Rotation

```bash
# Configure logrotate
sudo nano /etc/logrotate.d/sol-green

# Add configuration
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

## Backup Strategy

### 1. Database Backup

```bash
# PostgreSQL backup
docker-compose exec postgres pg_dump -U solgreen solgreen > backup_$(date +%Y%m%d_%H%M%S).sql

# Automated backup script
#!/bin/bash
BACKUP_DIR="/opt/sol-green/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Database backup
docker-compose exec postgres pg_dump -U solgreen solgreen > $BACKUP_DIR/db_$DATE.sql

# Configuration backup
tar -czf $BACKUP_DIR/config_$DATE.tar.gz .env

# Delete backups older than 7 days
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### 2. Contract Backup

```bash
# Export contract data
solana program dump <program-id> contract_backup.so

# Backup private keys
cp ~/.config/solana/id.json /secure/backup/wallet_backup.json
```

## Performance Optimization

### 1. Docker Optimization

```yaml
# Optimized docker-compose.yml
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

### 2. Nginx Optimization

```nginx
# Nginx performance configuration
worker_processes auto;
worker_connections 1024;

# Gzip compression
gzip on;
gzip_types text/plain text/css application/json application/javascript;

# Cache configuration
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

### 3. Database Optimization

```sql
-- Create indexes
CREATE INDEX CONCURRENTLY idx_green_behaviors_wallet_time ON green_behaviors(wallet_addr, submit_time);
CREATE INDEX CONCURRENTLY idx_challenges_status_time ON challenges(status, created_at);

-- Database configuration optimization
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
```

## Security Configuration

### 1. Firewall Configuration

```bash
# Configure UFW firewall
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw --force enable    # Enable firewall
```

### 2. SSL/TLS Configuration

```nginx
# SSL configuration
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
ssl_prefer_server_ciphers off;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
```

### 3. API Security

```bash
# Configure API rate limiting
RATE_LIMIT_REQUESTS=60
RATE_LIMIT_BURST=10

# JWT configuration
JWT_SECRET=your-very-long-random-secret-key
```

## Troubleshooting

### Common Issues

#### Container Won't Start
```bash
# Check container logs
docker-compose logs backend

# Check port conflicts
netstat -tlnp | grep :8080

# Check disk space
df -h
```

#### Database Connection Failed
```bash
# Check database container
docker-compose ps postgres

# Test database connection
docker-compose exec postgres psql -U solgreen -d solgreen -c "SELECT 1;"

# Check database logs
docker-compose logs postgres
```

#### Frontend Build Failed
```bash
# Check Node.js version
node --version
npm --version

# Clean cache
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

#### Smart Contract Deployment Failed
```bash
# Check Solana configuration
solana config get

# Check wallet balance
solana balance

# Check program size
ls -lh target/deploy/*.so
```

### Log Analysis

```bash
# View error logs
docker-compose logs backend | grep -i error

# View performance logs
docker-compose logs backend | grep -i "took\|latency"

# View blockchain related logs
docker-compose logs backend | grep -i "solana\|blockchain"
```

## Update Deployment

### 1. Code Update

```bash
# Pull latest code
git pull origin main

# Rebuild images
docker-compose build --no-cache

# Rolling update
docker-compose up -d --scale backend=2
docker-compose up -d --scale backend=1
```

### 2. Zero-Downtime Deployment

```bash
# Blue-green deployment
docker-compose up -d backend-green
# Switch traffic
docker-compose stop backend-blue
docker-compose rm backend-blue
```

### 3. Database Migration

```bash
# Create migration script
docker-compose exec backend ./migrate

# Verify migration
docker-compose exec postgres psql -U solgreen -d solgreen -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 5;"
```

## Scaling Deployment

### Horizontal Scaling

```bash
# Scale backend instances
docker-compose up -d --scale backend=3

# Configure load balancing
# Use Nginx or cloud load balancer
```

### Vertical Scaling

```bash
# Increase resource limits
# Update resource configuration in docker-compose.yml

services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
```

### Global Deployment

```bash
# Use CDN
# Cloudflare / AWS CloudFront / Aliyun CDN

# Multi-region deployment
# AWS: us-east-1, eu-west-1, ap-southeast-1
# Aliyun: cn-hangzhou, cn-shanghai, cn-shenzhen
```

---

**Document Version**: v3.0.0  
**Last Update**: 2024-01-24  
**Compatible Version**: Docker 20.10+, Docker Compose 2.0+, Kubernetes 1.24+
