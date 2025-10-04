# Deployment Guide - Expensio

Production deployment guide for Expensio expense management system.

## Table of Contents

1. [Production Requirements](#production-requirements)
2. [Backend Deployment](#backend-deployment)
3. [Frontend Deployment](#frontend-deployment)
4. [Database Setup](#database-setup)
5. [Environment Variables](#environment-variables)
6. [Docker Deployment](#docker-deployment)
7. [Monitoring & Logging](#monitoring--logging)
8. [Security Checklist](#security-checklist)

## Production Requirements

### Minimum Server Specifications

**Backend Server:**

- CPU: 2 cores
- RAM: 2 GB
- Storage: 20 GB SSD
- OS: Linux (Ubuntu 22.04 recommended)

**Database Server:**

- MongoDB: 4 GB RAM, 50 GB storage
- Redis: 1 GB RAM, 10 GB storage

### Software Requirements

- Go 1.21+
- Node.js 18+
- MongoDB 7.0+
- Redis 7.0+
- Nginx (reverse proxy)
- SSL Certificate (Let's Encrypt recommended)

## Backend Deployment

### 1. Build Binary

```bash
cd backend

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o expensio cmd/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o expensio.exe cmd/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o expensio cmd/main.go
```

### 2. Create Systemd Service (Linux)

Create `/etc/systemd/system/expensio.service`:

```ini
[Unit]
Description=Expensio Backend API
After=network.target mongodb.service redis.service

[Service]
Type=simple
User=expensio
WorkingDirectory=/opt/expensio
ExecStart=/opt/expensio/expensio
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=expensio

Environment="PORT=8080"
Environment="MONGODB_URI=mongodb://localhost:27017"
Environment="MONGODB_DATABASE=expensio"

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable expensio
sudo systemctl start expensio
sudo systemctl status expensio
```

### 3. Nginx Configuration

Create `/etc/nginx/sites-available/expensio-api`:

```nginx
server {
    listen 80;
    server_name api.expensio.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

Enable:

```bash
sudo ln -s /etc/nginx/sites-available/expensio-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 4. SSL Certificate (Let's Encrypt)

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.expensio.com
```

## Frontend Deployment

### 1. Build for Production

```bash
cd frontend

# Install dependencies
pnpm install

# Build
pnpm build
```

### 2. Deploy to Vercel (Recommended)

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel --prod
```

Set environment variables in Vercel dashboard:

- `NEXTAUTH_URL`: https://expensio.com
- `NEXTAUTH_SECRET`: (generate secure random string)
- `NEXT_PUBLIC_API_URL`: https://api.expensio.com/api/v1

### 3. Deploy to Custom Server

#### Using PM2

```bash
# Install PM2
npm install -g pm2

# Start application
pm2 start npm --name "expensio-frontend" -- start

# Save PM2 configuration
pm2 save

# Setup startup script
pm2 startup
```

#### Nginx Configuration

Create `/etc/nginx/sites-available/expensio-frontend`:

```nginx
server {
    listen 80;
    server_name expensio.com www.expensio.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 4. Docker Deployment

```bash
# Build frontend image
docker build -t expensio-frontend:latest .

# Run container
docker run -d \
  -p 3000:3000 \
  -e NEXTAUTH_URL=https://expensio.com \
  -e NEXTAUTH_SECRET=your-secret \
  -e NEXT_PUBLIC_API_URL=https://api.expensio.com/api/v1 \
  --name expensio-frontend \
  expensio-frontend:latest
```

## Database Setup

### MongoDB Production Setup

#### 1. Install MongoDB

```bash
# Ubuntu/Debian
wget -qO - https://www.mongodb.org/static/pgp/server-7.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt update
sudo apt install -y mongodb-org
```

#### 2. Configure MongoDB

Edit `/etc/mongod.conf`:

```yaml
net:
  port: 27017
  bindIp: 127.0.0.1

security:
  authorization: enabled

storage:
  dbPath: /var/lib/mongodb
  journal:
    enabled: true
```

#### 3. Create Database User

```bash
mongosh

use admin
db.createUser({
  user: "expensio_admin",
  pwd: "secure-password-here",
  roles: [
    { role: "readWrite", db: "expensio" },
    { role: "dbAdmin", db: "expensio" }
  ]
})
```

Update backend `.env`:

```env
MONGODB_URI=mongodb://expensio_admin:secure-password-here@localhost:27017/expensio?authSource=admin
```

### Redis Production Setup

#### 1. Install Redis

```bash
sudo apt install redis-server
```

#### 2. Configure Redis

Edit `/etc/redis/redis.conf`:

```conf
bind 127.0.0.1
port 6379
requirepass your-redis-password
maxmemory 256mb
maxmemory-policy allkeys-lru
```

Restart Redis:

```bash
sudo systemctl restart redis-server
```

Update backend `.env`:

```env
REDIS_URL=localhost:6379
REDIS_PASSWORD=your-redis-password
```

## Environment Variables

### Backend Production `.env`

```env
# Server
PORT=8080

# MongoDB
MONGODB_URI=mongodb://user:pass@localhost:27017/expensio?authSource=admin
MONGODB_DATABASE=expensio

# Redis
REDIS_URL=localhost:6379
REDIS_PASSWORD=your-redis-password

# JWT - Use strong random strings!
JWT_SECRET=generate-with-openssl-rand-base64-32
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# External APIs
EXCHANGE_RATE_API_KEY=your-api-key

# OCR
TESSERACT_PATH=/usr/bin/tesseract
```

### Frontend Production `.env.production`

```env
NEXTAUTH_URL=https://expensio.com
NEXTAUTH_SECRET=generate-with-openssl-rand-base64-32
NEXT_PUBLIC_API_URL=https://api.expensio.com/api/v1
```

## Docker Deployment

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: "3.8"

services:
  mongodb:
    image: mongo:7.0
    container_name: expensio-mongo
    restart: always
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD}
      MONGO_INITDB_DATABASE: expensio
    volumes:
      - mongo-data:/data/db

  redis:
    image: redis:7.0-alpine
    container_name: expensio-redis
    restart: always
    ports:
      - "6379:6379"
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data

  backend:
    build: ./backend
    container_name: expensio-backend
    restart: always
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
      - redis
    environment:
      PORT: 8080
      MONGODB_URI: mongodb://admin:${MONGO_PASSWORD}@mongodb:27017/expensio?authSource=admin
      MONGODB_DATABASE: expensio
      REDIS_URL: redis:6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
      JWT_ACCESS_EXPIRY: 15m
      JWT_REFRESH_EXPIRY: 168h

  frontend:
    build: ./frontend
    container_name: expensio-frontend
    restart: always
    ports:
      - "3000:3000"
    depends_on:
      - backend
    environment:
      NEXTAUTH_URL: ${NEXTAUTH_URL}
      NEXTAUTH_SECRET: ${NEXTAUTH_SECRET}
      NEXT_PUBLIC_API_URL: ${NEXT_PUBLIC_API_URL}

volumes:
  mongo-data:
  redis-data:
```

Create `.env` for Docker Compose:

```env
MONGO_PASSWORD=your-mongo-password
REDIS_PASSWORD=your-redis-password
JWT_SECRET=your-jwt-secret
NEXTAUTH_URL=https://expensio.com
NEXTAUTH_SECRET=your-nextauth-secret
NEXT_PUBLIC_API_URL=https://api.expensio.com/api/v1
```

Deploy:

```bash
docker-compose up -d
```

## Monitoring & Logging

### Backend Logging

Use structured logging:

```go
// Already configured in the application
// Logs are output to stdout/stderr
```

View logs:

```bash
# Systemd
sudo journalctl -u expensio -f

# Docker
docker logs -f expensio-backend

# PM2
pm2 logs expensio-backend
```

### Monitoring Tools

#### 1. Application Metrics

Use Prometheus + Grafana:

```bash
# Install Prometheus
docker run -d -p 9090:9090 prom/prometheus

# Install Grafana
docker run -d -p 3001:3000 grafana/grafana
```

#### 2. Health Checks

Create health check endpoint (already available):

```
GET /api/v1/health
```

#### 3. Uptime Monitoring

Use external services:

- UptimeRobot
- Pingdom
- StatusCake

## Security Checklist

### Backend Security

- [ ] Change JWT_SECRET to strong random string
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS properly (whitelist domains)
- [ ] Use environment variables for secrets
- [ ] Enable MongoDB authentication
- [ ] Enable Redis password
- [ ] Set up firewall rules
- [ ] Regular security updates
- [ ] Rate limiting enabled
- [ ] Input validation on all endpoints
- [ ] SQL injection prevention (MongoDB)
- [ ] XSS prevention

### Frontend Security

- [ ] Change NEXTAUTH_SECRET to strong random string
- [ ] Enable HTTPS
- [ ] Set secure cookie flags
- [ ] Implement CSP headers
- [ ] XSS prevention in React
- [ ] Validate all user inputs
- [ ] Sanitize file uploads

### Infrastructure Security

- [ ] Use strong passwords
- [ ] Enable SSH key authentication
- [ ] Disable root SSH login
- [ ] Configure firewall (UFW/iptables)
- [ ] Regular backups
- [ ] Update packages regularly
- [ ] Use fail2ban for brute force protection
- [ ] Monitor logs for suspicious activity

## Backup Strategy

### MongoDB Backup

```bash
# Create backup script
#!/bin/bash
BACKUP_DIR="/backups/mongodb"
DATE=$(date +%Y%m%d_%H%M%S)

mongodump --uri="mongodb://user:pass@localhost:27017/expensio" \
  --out="$BACKUP_DIR/$DATE"

# Keep only last 7 days
find $BACKUP_DIR -type d -mtime +7 -exec rm -rf {} \;
```

Add to crontab:

```bash
0 2 * * * /usr/local/bin/mongo-backup.sh
```

### Application Backup

```bash
# Backup configuration and code
tar -czf /backups/expensio-$(date +%Y%m%d).tar.gz \
  /opt/expensio \
  /etc/nginx/sites-available/expensio-* \
  /etc/systemd/system/expensio.service
```

## Scaling

### Horizontal Scaling

1. **Load Balancer:** Use Nginx or HAProxy
2. **Multiple Backend Instances:** Run multiple Go servers
3. **Shared Session Store:** Use Redis for sessions
4. **Database Replication:** MongoDB replica set

### Vertical Scaling

1. Increase server resources (CPU, RAM)
2. Optimize database queries
3. Add database indexes
4. Implement caching strategy

## Troubleshooting

### Backend Issues

```bash
# Check service status
sudo systemctl status expensio

# View logs
sudo journalctl -u expensio -n 100

# Test API
curl http://localhost:8080/api/v1/health
```

### Frontend Issues

```bash
# Check PM2 status
pm2 status

# View logs
pm2 logs expensio-frontend

# Restart
pm2 restart expensio-frontend
```

### Database Issues

```bash
# MongoDB status
sudo systemctl status mongod

# Redis status
sudo systemctl status redis-server

# Test connections
mongosh --host localhost --port 27017
redis-cli -h localhost -p 6379 ping
```

## Performance Optimization

1. **Enable Gzip Compression** (Nginx)
2. **Use CDN** for static assets
3. **Implement Database Indexing**
4. **Use Redis Caching**
5. **Optimize Images**
6. **Lazy Loading** for frontend
7. **Connection Pooling** (already configured)

## Support

For deployment issues:

- Check logs first
- Review this guide
- Contact: devops@expensio.com

---

**Deployment Checklist:**

- [ ] Backend deployed and running
- [ ] Frontend deployed and running
- [ ] MongoDB secured and backed up
- [ ] Redis secured
- [ ] SSL certificates installed
- [ ] Environment variables set
- [ ] Monitoring configured
- [ ] Backups scheduled
- [ ] Security checklist completed
- [ ] Health checks passing
