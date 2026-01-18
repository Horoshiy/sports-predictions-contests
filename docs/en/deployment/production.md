# Production Deployment Guide

This guide covers deploying the Sports Prediction Contests platform to a production environment with security, scalability, and reliability considerations.

## Production Architecture

### Infrastructure Requirements

#### Minimum Hardware Specifications
- **CPU**: 4 cores (8 recommended)
- **RAM**: 8GB (16GB recommended)
- **Storage**: 50GB SSD (100GB+ recommended)
- **Network**: 1Gbps connection

#### Recommended Cloud Resources
- **AWS**: t3.large or larger instances
- **Google Cloud**: n1-standard-4 or larger
- **Azure**: Standard_D4s_v3 or larger
- **DigitalOcean**: 4GB+ droplets

### Security Considerations

#### Environment Variables Security
```bash
# Use strong, unique secrets
JWT_SECRET=$(openssl rand -base64 32)
POSTGRES_PASSWORD=$(openssl rand -base64 32)

# Never use default values in production
# ❌ Bad
JWT_SECRET=your_jwt_secret_key_here

# ✅ Good  
JWT_SECRET=xK9mP2vN8qR5tY7uI3oP6sA1dF4gH8jL9nM2bV5cX8z
```

#### Database Security
```bash
# Use SSL connections
DATABASE_URL=postgres://sports_user:${POSTGRES_PASSWORD}@postgres:5432/sports_prediction?sslmode=require

# Restrict database access
# Configure PostgreSQL pg_hba.conf for specific IP ranges
# Use connection pooling with PgBouncer
```

#### Network Security
```bash
# Restrict CORS origins
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# Use HTTPS only
# Configure reverse proxy (nginx/traefik) with SSL termination
# Implement rate limiting
```

## Production Docker Compose

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  # Reverse Proxy with SSL
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - api-gateway
    networks:
      - sports-network
    restart: unless-stopped

  # Database with persistence
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    networks:
      - sports-network
    restart: unless-stopped
    command: >
      postgres
      -c max_connections=200
      -c shared_buffers=256MB
      -c effective_cache_size=1GB

  # Redis with persistence
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    networks:
      - sports-network
    restart: unless-stopped
    command: redis-server --appendonly yes --maxmemory 512mb --maxmemory-policy allkeys-lru

  # API Gateway
  api-gateway:
    build:
      context: ./backend/api-gateway
      dockerfile: Dockerfile
    environment:
      - API_GATEWAY_PORT=8080
      - USER_SERVICE_ENDPOINT=user-service:8084
      - CONTEST_SERVICE_ENDPOINT=contest-service:8085
      - PREDICTION_SERVICE_ENDPOINT=prediction-service:8086
      - SCORING_SERVICE_ENDPOINT=scoring-service:8087
      - SPORTS_SERVICE_ENDPOINT=sports-service:8088
      - NOTIFICATION_SERVICE_ENDPOINT=notification-service:8089
      - JWT_SECRET=${JWT_SECRET}
      - CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS}
      - LOG_LEVEL=info
    depends_on:
      - postgres
      - redis
    networks:
      - sports-network
    restart: unless-stopped
    deploy:
      replicas: 2
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M

  # Microservices (similar configuration for each)
  user-service:
    build:
      context: ./backend/user-service
      dockerfile: Dockerfile
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - JWT_SECRET=${JWT_SECRET}
      - USER_SERVICE_PORT=8084
    depends_on:
      - postgres
      - redis
    networks:
      - sports-network
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 256M

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  sports-network:
    driver: bridge
```

## SSL/TLS Configuration

### Nginx Configuration (`nginx/nginx.conf`)

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api_gateway {
        server api-gateway:8080;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=auth:10m rate=5r/s;

    # SSL Configuration
    server {
        listen 80;
        server_name yourdomain.com www.yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name yourdomain.com www.yourdomain.com;

        # SSL Certificates
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        
        # SSL Security
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
        ssl_prefer_server_ciphers off;
        ssl_session_cache shared:SSL:10m;

        # Security Headers
        add_header Strict-Transport-Security "max-age=63072000" always;
        add_header X-Frame-Options DENY always;
        add_header X-Content-Type-Options nosniff always;

        # API Proxy
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://api_gateway/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Auth endpoints with stricter rate limiting
        location /api/v1/auth/ {
            limit_req zone=auth burst=10 nodelay;
            proxy_pass http://api_gateway/v1/auth/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Frontend
        location / {
            root /usr/share/nginx/html;
            index index.html;
            try_files $uri $uri/ /index.html;
        }
    }
}
```

## Environment Configuration

### Production Environment Variables (`.env.prod`)

```bash
# Environment
NODE_ENV=production
GO_ENV=production

# Database (use managed database service)
DATABASE_URL=postgres://sports_user:${POSTGRES_PASSWORD}@your-db-host:5432/sports_prediction?sslmode=require
POSTGRES_DB=sports_prediction
POSTGRES_USER=sports_user
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# Redis (use managed Redis service)
REDIS_URL=redis://your-redis-host:6379
REDIS_PASSWORD=${REDIS_PASSWORD}

# Services
API_GATEWAY_PORT=8080
USER_SERVICE_PORT=8084
CONTEST_SERVICE_PORT=8085
PREDICTION_SERVICE_PORT=8086
SCORING_SERVICE_PORT=8087
SPORTS_SERVICE_PORT=8088
NOTIFICATION_SERVICE_PORT=8089

# Security
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRATION=24h
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# External Services
TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
TELEGRAM_ENABLED=true
SMTP_HOST=${SMTP_HOST}
SMTP_PORT=587
SMTP_USER=${SMTP_USER}
SMTP_PASSWORD=${SMTP_PASSWORD}
EMAIL_ENABLED=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Monitoring
METRICS_ENABLED=true
HEALTH_CHECK_INTERVAL=30s
```

## Deployment Process

### 1. Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create application directory
sudo mkdir -p /opt/sports-prediction-contests
sudo chown $USER:$USER /opt/sports-prediction-contests
cd /opt/sports-prediction-contests
```

### 2. Code Deployment

```bash
# Clone repository
git clone https://github.com/coleam00/dynamous-kiro-hackathon .

# Checkout production branch/tag
git checkout v1.0.0

# Set up environment
cp .env.example .env.prod
# Edit .env.prod with production values
```

### 3. SSL Certificate Setup

```bash
# Using Let's Encrypt with Certbot
sudo apt install certbot
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# Copy certificates
sudo mkdir -p nginx/ssl
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem nginx/ssl/key.pem
sudo chown -R $USER:$USER nginx/ssl
```

### 4. Database Setup

```bash
# For managed database services (recommended)
# Create database instance on AWS RDS, Google Cloud SQL, etc.
# Update DATABASE_URL in .env.prod

# For self-hosted PostgreSQL
docker-compose -f docker-compose.prod.yml up -d postgres
sleep 10
docker-compose -f docker-compose.prod.yml exec postgres psql -U sports_user -d sports_prediction -c "SELECT version();"
```

### 5. Application Deployment

```bash
# Build and start services
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# Verify deployment
docker-compose -f docker-compose.prod.yml ps
curl -k https://yourdomain.com/api/health
```

## Monitoring and Logging

### Health Checks

```bash
# Create health check script
cat > health-check.sh << 'EOF'
#!/bin/bash
SERVICES=("api-gateway" "user-service" "contest-service" "prediction-service" "scoring-service" "sports-service" "notification-service")

for service in "${SERVICES[@]}"; do
    if ! curl -f -s http://localhost:8080/v1/${service}/health > /dev/null; then
        echo "❌ $service is unhealthy"
        exit 1
    else
        echo "✅ $service is healthy"
    fi
done
echo "All services are healthy"
EOF

chmod +x health-check.sh
```

### Log Management

```bash
# Configure log rotation
sudo tee /etc/logrotate.d/docker-containers << 'EOF'
/var/lib/docker/containers/*/*.log {
    rotate 7
    daily
    compress
    size=1M
    missingok
    delaycompress
    copytruncate
}
EOF

# View logs
docker-compose -f docker-compose.prod.yml logs -f --tail=100
```

### Monitoring Setup

```bash
# Add Prometheus monitoring (optional)
# Create monitoring/prometheus.yml
# Add Grafana dashboards
# Set up alerting rules
```

## Backup and Recovery

### Database Backup

```bash
# Create backup script
cat > backup-db.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

# Database backup
docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U sports_user sports_prediction | gzip > $BACKUP_DIR/db_backup_$DATE.sql.gz

# Keep only last 7 days
find $BACKUP_DIR -name "db_backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: db_backup_$DATE.sql.gz"
EOF

chmod +x backup-db.sh

# Add to crontab for daily backups
echo "0 2 * * * /opt/sports-prediction-contests/backup-db.sh" | crontab -
```

### Recovery Process

```bash
# Restore from backup
gunzip -c /opt/backups/db_backup_YYYYMMDD_HHMMSS.sql.gz | \
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U sports_user -d sports_prediction
```

## Scaling Considerations

### Horizontal Scaling

```bash
# Scale specific services
docker-compose -f docker-compose.prod.yml up -d --scale api-gateway=3 --scale user-service=2

# Use Docker Swarm or Kubernetes for advanced orchestration
```

### Database Scaling

```bash
# Use read replicas for read-heavy workloads
# Implement connection pooling with PgBouncer
# Consider database sharding for very large datasets
```

### Caching Strategy

```bash
# Redis clustering for high availability
# CDN for static assets
# Application-level caching for frequently accessed data
```

## Security Checklist

- [ ] Strong, unique passwords and secrets
- [ ] SSL/TLS certificates properly configured
- [ ] Database connections encrypted
- [ ] CORS properly configured
- [ ] Rate limiting implemented
- [ ] Security headers configured
- [ ] Regular security updates applied
- [ ] Firewall rules configured
- [ ] Access logs monitored
- [ ] Backup encryption enabled

## Maintenance

### Regular Tasks

```bash
# Weekly maintenance script
cat > maintenance.sh << 'EOF'
#!/bin/bash
echo "Starting weekly maintenance..."

# Update Docker images
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d

# Clean up unused Docker resources
docker system prune -f

# Rotate logs
logrotate -f /etc/logrotate.d/docker-containers

# Run health checks
./health-check.sh

echo "Maintenance completed"
EOF

chmod +x maintenance.sh
```

### Updates and Rollbacks

```bash
# Update process
git fetch origin
git checkout v1.1.0
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# Rollback process
git checkout v1.0.0
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d
```

## Troubleshooting Production Issues

### Common Production Problems

1. **High Memory Usage**: Scale services or optimize queries
2. **Database Connection Limits**: Implement connection pooling
3. **SSL Certificate Expiry**: Set up automatic renewal
4. **Service Discovery Issues**: Check Docker network configuration
5. **Performance Degradation**: Monitor and optimize database queries

### Emergency Procedures

```bash
# Quick service restart
docker-compose -f docker-compose.prod.yml restart [service-name]

# Emergency rollback
git checkout [previous-stable-tag]
docker-compose -f docker-compose.prod.yml up -d

# Database emergency access
docker-compose -f docker-compose.prod.yml exec postgres psql -U sports_user -d sports_prediction
```

---

**Production deployment requires careful planning and testing. Always test deployment procedures in a staging environment first.**
