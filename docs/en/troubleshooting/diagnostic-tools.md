# Diagnostic Tools and Debugging

Comprehensive guide to diagnostic tools and debugging procedures for the Sports Prediction Contests platform.

## Docker Diagnostics

### Container Health Monitoring

#### Check Container Status
```bash
# List all containers with status
docker-compose ps

# Show detailed container information
docker inspect <container_name>

# Check container resource usage
docker stats --no-stream

# Monitor resource usage in real-time
docker stats
```

#### Container Logs Analysis
```bash
# View logs for all services
docker-compose logs

# View logs for specific service
docker-compose logs user-service

# Follow logs in real-time
docker-compose logs -f api-gateway

# View last N lines of logs
docker-compose logs --tail=50 postgres

# Search logs for specific patterns
docker-compose logs | grep -i "error\|warning\|failed"

# Export logs to file
docker-compose logs > system-logs.txt
```

### Network Diagnostics

#### Docker Network Analysis
```bash
# List Docker networks
docker network ls

# Inspect network configuration
docker network inspect sports-network

# Check network connectivity between containers
docker-compose exec api-gateway ping user-service
docker-compose exec user-service ping postgres
```

#### Port and Service Discovery
```bash
# Check which ports are exposed
docker-compose port api-gateway 8080
docker-compose port postgres 5432

# Test service connectivity
docker-compose exec api-gateway curl http://user-service:8084/health
docker-compose exec user-service curl http://postgres:5432
```

## Database Diagnostics

### PostgreSQL Health Checks

#### Connection Testing
```bash
# Test database connectivity
docker-compose exec postgres pg_isready -U sports_user -d sports_prediction

# Connect to database shell
docker-compose exec postgres psql -U sports_user -d sports_prediction

# Check database version and status
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "SELECT version();"
```

#### Database Performance Analysis
```bash
# Check active connections
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT count(*) as active_connections 
FROM pg_stat_activity 
WHERE state = 'active';"

# Check database size
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT pg_size_pretty(pg_database_size('sports_prediction')) as db_size;"

# Check table sizes
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT schemaname,tablename,pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size 
FROM pg_tables 
WHERE schemaname = 'public' 
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"

# Check slow queries
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT query, mean_time, calls, total_time 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;"
```

### Redis Diagnostics

#### Redis Health and Performance
```bash
# Test Redis connectivity
docker-compose exec redis redis-cli ping

# Check Redis info
docker-compose exec redis redis-cli info

# Monitor Redis commands in real-time
docker-compose exec redis redis-cli monitor

# Check Redis memory usage
docker-compose exec redis redis-cli info memory

# List all keys (use carefully in production)
docker-compose exec redis redis-cli keys "*"

# Check specific key types and values
docker-compose exec redis redis-cli type "leaderboard:contest:<CONTEST_ID>"
docker-compose exec redis redis-cli get "session:user:<USER_ID>"
```

## Service-Specific Diagnostics

### API Gateway Diagnostics

#### Health and Status Checks
```bash
# Check API Gateway health
curl -v http://localhost:8080/health

# Test routing to services
curl -v http://localhost:8080/v1/auth/health
curl -v http://localhost:8080/v1/contests/health

# Check API Gateway configuration
docker-compose exec api-gateway env | grep -E "(PORT|ENDPOINT|JWT)"
```

#### Request Tracing
```bash
# Enable verbose curl output for debugging
curl -v -H "Authorization: Bearer $JWT_TOKEN" \
     http://localhost:8080/v1/contests

# Test with different HTTP methods
curl -X POST -v -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/v1/test-endpoint
```

### Microservice Health Checks

#### Individual Service Testing
```bash
# Test each service health endpoint
services=("auth" "contests" "predictions" "scores" "sports" "notifications")

for service in "${services[@]}"; do
    echo "Testing $service service:"
    curl -s http://localhost:8080/v1/$service/health | jq '.'
    echo "---"
done
```

#### Service Configuration Verification
```bash
# Check service environment variables
docker-compose exec user-service env | grep -E "(DATABASE|REDIS|JWT|PORT)"
docker-compose exec contest-service env | grep -E "(DATABASE|REDIS|JWT|PORT)"

# Verify service can connect to dependencies
docker-compose exec user-service nc -zv postgres 5432
docker-compose exec scoring-service nc -zv redis 6379
```

## Performance Diagnostics

### System Resource Monitoring

#### CPU and Memory Analysis
```bash
# Check system resources
free -h
top -bn1 | head -20

# Monitor Docker container resources
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

# Check disk usage
df -h
docker system df
```

#### Application Performance Metrics
```bash
# Check response times for API endpoints
time curl -s http://localhost:8080/v1/contests > /dev/null

# Monitor database query performance
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT query, mean_time, calls 
FROM pg_stat_statements 
WHERE mean_time > 100 
ORDER BY mean_time DESC;"

# Check Redis performance
docker-compose exec redis redis-cli --latency-history -i 1
```

### Load Testing Diagnostics

#### Simple Load Testing
```bash
# Install Apache Bench (if not available)
# sudo apt-get install apache2-utils  # Ubuntu/Debian
# brew install httpie  # macOS

# Simple load test
ab -n 100 -c 10 http://localhost:8080/health

# Test with authentication
ab -n 50 -c 5 -H "Authorization: Bearer $JWT_TOKEN" \
   http://localhost:8080/v1/contests
```

## Log Analysis Tools

### Structured Log Analysis

#### Log Parsing and Filtering
```bash
# Parse JSON logs (if using structured logging)
docker-compose logs user-service | jq -r 'select(.level == "error") | .message'

# Filter logs by timestamp
docker-compose logs --since="2026-01-18T10:00:00" api-gateway

# Count error occurrences
docker-compose logs | grep -c "ERROR"

# Extract unique error messages
docker-compose logs | grep "ERROR" | sort | uniq -c | sort -nr
```

#### Log Aggregation
```bash
# Collect all service logs with timestamps
docker-compose logs -t > full-system-logs.txt

# Create service-specific log files
services=("api-gateway" "user-service" "contest-service" "prediction-service" "scoring-service" "sports-service" "notification-service")

for service in "${services[@]}"; do
    docker-compose logs $service > logs/$service.log
done
```

## Debugging Workflows

### Issue Investigation Process

#### Step 1: Initial Assessment
```bash
# Quick system health check
echo "=== System Health Check ==="
docker-compose ps
curl -s http://localhost:8080/health | jq '.'
docker stats --no-stream | head -10
```

#### Step 2: Service-Specific Investigation
```bash
# Check specific service that's failing
SERVICE_NAME="user-service"

echo "=== $SERVICE_NAME Investigation ==="
docker-compose logs --tail=50 $SERVICE_NAME
docker-compose exec $SERVICE_NAME env | grep -E "(DATABASE|REDIS|JWT)"
docker inspect $(docker-compose ps -q $SERVICE_NAME)
```

#### Step 3: Dependency Verification
```bash
# Check service dependencies
echo "=== Dependency Check ==="
docker-compose exec postgres pg_isready -U sports_user
docker-compose exec redis redis-cli ping
docker-compose exec api-gateway nc -zv postgres 5432
docker-compose exec api-gateway nc -zv redis 6379
```

### Common Debugging Scenarios

#### Database Connection Issues
```bash
# Debug database connectivity
echo "=== Database Debug ==="
docker-compose logs postgres | tail -20
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "\l"
docker-compose exec user-service nc -zv postgres 5432

# Check database locks
docker-compose exec postgres psql -U sports_user -d sports_prediction -c "
SELECT blocked_locks.pid AS blocked_pid,
       blocked_activity.usename AS blocked_user,
       blocking_locks.pid AS blocking_pid,
       blocking_activity.usename AS blocking_user,
       blocked_activity.query AS blocked_statement
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;"
```

#### Authentication Issues
```bash
# Debug JWT authentication
echo "=== JWT Debug ==="
echo "JWT_SECRET check:"
docker-compose exec api-gateway env | grep JWT_SECRET
docker-compose exec user-service env | grep JWT_SECRET

# Test token generation
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' \
  -v
```

## Automated Diagnostic Scripts

### Health Check Script
```bash
#!/bin/bash
# save as scripts/health-check.sh

echo "=== Sports Prediction Platform Health Check ==="
echo "Timestamp: $(date)"
echo

# Check Docker
echo "Docker Status:"
docker --version
docker-compose --version
echo

# Check containers
echo "Container Status:"
docker-compose ps
echo

# Check services
echo "Service Health:"
services=("health" "v1/auth/health" "v1/contests/health" "v1/predictions/health" "v1/scores/health" "v1/sports/health" "v1/notifications/health")

for endpoint in "${services[@]}"; do
    echo -n "$endpoint: "
    if curl -s -f http://localhost:8080/$endpoint > /dev/null; then
        echo "✅ OK"
    else
        echo "❌ FAIL"
    fi
done
echo

# Check database
echo "Database Status:"
if docker-compose exec -T postgres pg_isready -U sports_user -d sports_prediction > /dev/null 2>&1; then
    echo "PostgreSQL: ✅ OK"
else
    echo "PostgreSQL: ❌ FAIL"
fi

# Check Redis
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "Redis: ✅ OK"
else
    echo "Redis: ❌ FAIL"
fi

echo
echo "=== Health Check Complete ==="
```

### Performance Monitor Script
```bash
#!/bin/bash
# save as scripts/performance-monitor.sh

echo "=== Performance Monitor ==="
echo "Timestamp: $(date)"
echo

# System resources
echo "System Resources:"
echo "Memory: $(free -h | grep Mem | awk '{print $3 "/" $2}')"
echo "Disk: $(df -h / | tail -1 | awk '{print $3 "/" $2 " (" $5 " used)"}')"
echo

# Docker resources
echo "Docker Container Resources:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"
echo

# Database performance
echo "Database Performance:"
docker-compose exec -T postgres psql -U sports_user -d sports_prediction -c "
SELECT count(*) as active_connections FROM pg_stat_activity WHERE state = 'active';" 2>/dev/null
echo

# Redis performance
echo "Redis Performance:"
docker-compose exec -T redis redis-cli info stats | grep -E "(total_commands_processed|used_memory_human)" 2>/dev/null
```

## Troubleshooting Checklist

### Pre-Debugging Checklist
- [ ] All containers are running (`docker-compose ps`)
- [ ] No port conflicts (`lsof -i :8080`)
- [ ] Sufficient system resources (`free -h`, `df -h`)
- [ ] Environment variables are set (`.env` file exists)
- [ ] Network connectivity between services
- [ ] Database and Redis are accessible

### During Investigation
- [ ] Check service logs for errors
- [ ] Verify service configuration
- [ ] Test individual service endpoints
- [ ] Check database connectivity and performance
- [ ] Monitor system resources
- [ ] Trace request flow through services

### Post-Resolution
- [ ] Document the issue and solution
- [ ] Update monitoring/alerting if needed
- [ ] Consider preventive measures
- [ ] Update troubleshooting documentation

---

*These diagnostic tools help identify and resolve issues quickly. Use them systematically when investigating problems with the Sports Prediction Contests platform.*
