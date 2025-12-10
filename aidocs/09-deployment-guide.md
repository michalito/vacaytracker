# 09 - Deployment Guide

> Docker setup, environment configuration, and production deployment

## Table of Contents

1. [Development Setup](#development-setup)
2. [Docker Configuration](#docker-configuration)
3. [Environment Variables](#environment-variables)
4. [Production Deployment](#production-deployment)
5. [Database Management](#database-management)
6. [Monitoring & Health Checks](#monitoring--health-checks)

---

## Development Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Docker | 24+ | Container runtime |
| Docker Compose | 2.20+ | Multi-container orchestration |
| Node.js | 20+ | Frontend development |
| Go | 1.23+ | Backend development (optional, for non-Docker) |

### Quick Start

```bash
# Clone repository
git clone <repository-url>
cd vacaytracker

# Start full stack with Docker
docker-compose -f docker-compose.dev.yml up

# Access:
# - Frontend: http://localhost:5173
# - Backend:  http://localhost:3000
# - API:      http://localhost:3000/api
```

### Local Development (Without Docker)

```bash
# Backend
cd vacaytracker-api
cp .env.example .env
# Edit .env with your settings
make run

# Frontend (new terminal)
cd vacaytracker-frontend
npm install
npm run dev
```

---

## Docker Configuration

### Backend Dockerfile

- [ ] **Create Dockerfile** `vacaytracker-api/Dockerfile`

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary
COPY --from=builder /app/server /app/server

# Copy migrations
COPY --from=builder /app/migrations /app/migrations

# Create data directory
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run server
CMD ["./server"]
```

### Frontend Dockerfile

- [ ] **Create Dockerfile** `vacaytracker-frontend/Dockerfile`

```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build application
RUN npm run build

# Runtime stage
FROM node:20-alpine

WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy built application
COPY --from=builder /app/build /app/build
COPY --from=builder /app/package*.json /app/

# Install production dependencies only
RUN npm ci --production

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Run application
CMD ["node", "build"]
```

### Development Docker Compose

- [ ] **Create docker-compose.dev.yml** `docker-compose.dev.yml`

```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./vacaytracker-api
      dockerfile: Dockerfile.dev
    ports:
      - "3000:3000"
    volumes:
      - ./vacaytracker-api:/app
      - backend-data:/app/data
    environment:
      - PORT=3000
      - ENV=development
      - DB_PATH=/app/data/vacaytracker.db
      - JWT_SECRET=${JWT_SECRET:-dev-secret-key-minimum-32-characters}
      - ADMIN_PASSWORD=${ADMIN_PASSWORD:-admin123}
      - ADMIN_EMAIL=${ADMIN_EMAIL:-admin@company.com}
      - ADMIN_NAME=${ADMIN_NAME:-Captain Admin}
    command: ["air", "-c", ".air.toml"]

  frontend:
    build:
      context: ./vacaytracker-frontend
      dockerfile: Dockerfile.dev
    ports:
      - "5173:5173"
    volumes:
      - ./vacaytracker-frontend:/app
      - /app/node_modules
    environment:
      - NODE_ENV=development
    depends_on:
      - backend

volumes:
  backend-data:
```

### Development Dockerfiles

- [ ] **Create Dockerfile.dev** `vacaytracker-api/Dockerfile.dev`

```dockerfile
FROM golang:1.23-alpine

WORKDIR /app

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]
```

- [ ] **Create Dockerfile.dev** `vacaytracker-frontend/Dockerfile.dev`

```dockerfile
FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
```

### Production Docker Compose

- [ ] **Create docker-compose.yml** `docker-compose.yml`

```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./vacaytracker-api
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    volumes:
      - backend-data:/app/data
    environment:
      - PORT=3000
      - ENV=production
      - DB_PATH=/app/data/vacaytracker.db
      - JWT_SECRET=${JWT_SECRET}
      - ADMIN_PASSWORD=${ADMIN_PASSWORD}
      - ADMIN_EMAIL=${ADMIN_EMAIL}
      - ADMIN_NAME=${ADMIN_NAME}
      - RESEND_API_KEY=${RESEND_API_KEY}
      - EMAIL_FROM_ADDRESS=${EMAIL_FROM_ADDRESS}
      - EMAIL_FROM_NAME=${EMAIL_FROM_NAME}
      - APP_URL=${APP_URL}
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  frontend:
    build:
      context: ./vacaytracker-frontend
      dockerfile: Dockerfile
    ports:
      - "80:3000"
    environment:
      - NODE_ENV=production
      - ORIGIN=${APP_URL}
    depends_on:
      backend:
        condition: service_healthy
    restart: unless-stopped

volumes:
  backend-data:
```

---

## Environment Variables

### Complete Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| **Server** |
| `PORT` | No | `3000` | Server port |
| `ENV` | No | `development` | Environment (development/production) |
| `APP_URL` | Production | `http://localhost:3000` | Public URL |
| **Database** |
| `DB_PATH` | No | `./data/vacaytracker.db` | SQLite database path |
| **Authentication** |
| `JWT_SECRET` | Yes | - | JWT signing secret (32+ chars) |
| `ADMIN_PASSWORD` | Yes | - | Initial admin password |
| `ADMIN_EMAIL` | No | `admin@company.com` | Initial admin email |
| `ADMIN_NAME` | No | `Admin` | Initial admin name |
| **Email** |
| `RESEND_API_KEY` | No | - | Resend API key for emails |
| `EMAIL_FROM_ADDRESS` | No | - | Sender email address |
| `EMAIL_FROM_NAME` | No | `VacayTracker` | Sender name |

### Development .env

```bash
# Server
PORT=3000
ENV=development
APP_URL=http://localhost:3000

# Database
DB_PATH=./data/vacaytracker.db

# Authentication
JWT_SECRET=development-secret-key-must-be-32-chars
ADMIN_PASSWORD=admin123
ADMIN_EMAIL=admin@company.com
ADMIN_NAME=Captain Admin

# Email (optional for development)
# RESEND_API_KEY=
# EMAIL_FROM_ADDRESS=
# EMAIL_FROM_NAME=VacayTracker
```

### Production .env

```bash
# Server
PORT=3000
ENV=production
APP_URL=https://vacaytracker.yourcompany.com

# Database
DB_PATH=/app/data/vacaytracker.db

# Authentication (use secure values!)
JWT_SECRET=<generate-with-openssl-rand-base64-32>
ADMIN_PASSWORD=<secure-admin-password>
ADMIN_EMAIL=admin@yourcompany.com
ADMIN_NAME=System Administrator

# Email
RESEND_API_KEY=re_xxxxxxxxxxxx
EMAIL_FROM_ADDRESS=noreply@yourcompany.com
EMAIL_FROM_NAME=VacayTracker
```

### Generating Secure Secrets

```bash
# Generate JWT secret
openssl rand -base64 32

# Generate admin password
openssl rand -base64 16

# Or using Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

---

## Production Deployment

### Pre-Deployment Checklist

- [ ] All environment variables configured
- [ ] JWT_SECRET is cryptographically random
- [ ] ADMIN_PASSWORD is secure
- [ ] SSL/TLS certificate configured
- [ ] Database backup strategy in place
- [ ] Monitoring configured
- [ ] Log aggregation configured

### Deployment Steps

#### 1. Prepare Server

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. Clone and Configure

```bash
# Clone repository
git clone <repository-url> /opt/vacaytracker
cd /opt/vacaytracker

# Create environment file
cp .env.example .env
nano .env  # Edit with production values
```

#### 3. Build and Deploy

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Check status
docker-compose ps
docker-compose logs -f
```

#### 4. Configure Reverse Proxy (Nginx)

```nginx
# /etc/nginx/sites-available/vacaytracker
server {
    listen 80;
    server_name vacaytracker.yourcompany.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name vacaytracker.yourcompany.com;

    ssl_certificate /etc/letsencrypt/live/vacaytracker.yourcompany.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/vacaytracker.yourcompany.com/privkey.pem;

    # Security headers
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Frontend
    location / {
        proxy_pass http://localhost:80;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # API
    location /api {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Health check
    location /health {
        proxy_pass http://localhost:3000;
    }
}
```

#### 5. SSL Certificate (Let's Encrypt)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d vacaytracker.yourcompany.com

# Auto-renewal (already configured by certbot)
sudo certbot renew --dry-run
```

### Deployment Automation

- [ ] **Create deploy.sh** `deploy.sh`

```bash
#!/bin/bash
set -e

echo "Deploying VacayTracker..."

# Pull latest changes
git pull origin main

# Build images
docker-compose build

# Stop current containers
docker-compose down

# Start new containers
docker-compose up -d

# Wait for health check
echo "Waiting for services to be healthy..."
sleep 10

# Check health
curl -f http://localhost:3000/health || exit 1

echo "Deployment complete!"
```

---

## Database Management

### Backup Strategy

```bash
# Manual backup
docker exec vacaytracker_backend_1 cp /app/data/vacaytracker.db /app/data/backup-$(date +%Y%m%d).db

# Automated daily backup (cron)
0 2 * * * /opt/vacaytracker/scripts/backup.sh
```

- [ ] **Create backup.sh** `scripts/backup.sh`

```bash
#!/bin/bash
BACKUP_DIR="/opt/vacaytracker/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Create backup directory
mkdir -p $BACKUP_DIR

# Copy database
docker cp vacaytracker_backend_1:/app/data/vacaytracker.db "$BACKUP_DIR/vacaytracker_$DATE.db"

# Compress
gzip "$BACKUP_DIR/vacaytracker_$DATE.db"

# Remove old backups
find $BACKUP_DIR -name "*.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: vacaytracker_$DATE.db.gz"
```

### Restore from Backup

```bash
# Stop backend
docker-compose stop backend

# Decompress backup
gunzip backups/vacaytracker_20240115.db.gz

# Copy to container volume
docker cp backups/vacaytracker_20240115.db vacaytracker_backend_1:/app/data/vacaytracker.db

# Start backend
docker-compose start backend
```

### Database Migrations

```bash
# Migrations run automatically on startup
# To run manually:
docker exec vacaytracker_backend_1 ./server -migrate
```

---

## Monitoring & Health Checks

### Health Check Endpoint

```bash
# Check application health
curl http://localhost:3000/health

# Response:
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

### Docker Health Checks

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' vacaytracker_backend_1

# View health check logs
docker inspect --format='{{json .State.Health}}' vacaytracker_backend_1 | jq
```

### Log Management

```bash
# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# View last 100 lines
docker-compose logs --tail=100 backend

# Export logs
docker-compose logs backend > backend.log 2>&1
```

### Monitoring Recommendations

| Tool | Purpose | Setup |
|------|---------|-------|
| **Uptime Kuma** | Uptime monitoring | Self-hosted, check /health |
| **Prometheus** | Metrics collection | Add metrics endpoint |
| **Grafana** | Visualization | Connect to Prometheus |
| **Loki** | Log aggregation | Docker log driver |

### Alerting

```bash
# Example: Simple health check with notification
#!/bin/bash
# /opt/vacaytracker/scripts/healthcheck.sh

if ! curl -sf http://localhost:3000/health > /dev/null; then
    echo "VacayTracker is DOWN!" | mail -s "Alert: VacayTracker" admin@company.com
fi
```

---

## Related Documents

- [04-backend-tasks.md](./04-backend-tasks.md) - Backend setup
- [05-frontend-tasks.md](./05-frontend-tasks.md) - Frontend setup
- [08-security-checklist.md](./08-security-checklist.md) - Security configuration
- [10-development-workflow.md](./10-development-workflow.md) - CI/CD pipeline
