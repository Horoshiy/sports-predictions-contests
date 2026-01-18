# Common Issues and Solutions

This guide covers the most frequently encountered issues when working with the Sports Prediction Contests platform and their solutions.

## Service Startup Issues

### Problem: Services fail to start
**Symptoms**: Docker containers exit immediately or show error messages

**Solutions**:
```bash
# Check Docker daemon is running
docker --version
sudo systemctl start docker  # Linux
# or restart Docker Desktop on macOS/Windows

# Check port availability
lsof -i :8080  # Check if port is in use
kill -9 <PID>  # Kill process using the port

# Restart services
make clean
make dev
```

### Problem: Database connection failed
**Symptoms**: Services can't connect to PostgreSQL

**Solutions**:
```bash
# Check PostgreSQL status
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres

# Reset database completely
docker-compose down -v
docker-compose up -d postgres
```

## Configuration Issues

### Problem: Environment variables not loaded
**Symptoms**: Services use default values instead of configured ones

**Solutions**:
```bash
# Verify .env file exists
ls -la .env

# Copy from example if missing
cp .env.example .env

# Check environment variables are loaded
docker-compose config
```

### Problem: JWT authentication fails
**Symptoms**: 401 Unauthorized errors for authenticated requests

**Solutions**:
```bash
# Check JWT secret is set
grep JWT_SECRET .env

# Regenerate JWT secret
openssl rand -base64 32

# Update .env file with new secret
JWT_SECRET=your_new_secret_here
```

## Network and Connectivity Issues

### Problem: API Gateway not responding
**Symptoms**: Connection refused on port 8080

**Solutions**:
```bash
# Check API Gateway status
curl http://localhost:8080/health

# Check container status
docker-compose ps api-gateway

# View API Gateway logs
docker-compose logs api-gateway

# Restart API Gateway
docker-compose restart api-gateway
```

### Problem: Service discovery issues
**Symptoms**: Services can't communicate with each other

**Solutions**:
```bash
# Check Docker network
docker network ls
docker network inspect sports-network

# Restart all services
docker-compose down
docker-compose up -d
```

## Performance Issues

### Problem: Slow API responses
**Symptoms**: Requests take longer than expected

**Solutions**:
```bash
# Check system resources
docker stats

# Check database performance
docker-compose exec postgres pg_stat_activity

# Restart Redis cache
docker-compose restart redis
```

### Problem: High memory usage
**Symptoms**: System becomes slow, containers restart

**Solutions**:
```bash
# Check memory usage
free -h
docker stats --no-stream

# Increase Docker memory limits
# Edit docker-compose.yml to add memory limits

# Clean up unused Docker resources
docker system prune -f
```

## Development Issues

### Problem: Hot reload not working
**Symptoms**: Code changes don't reflect immediately

**Solutions**:
```bash
# For Go services - restart specific service
docker-compose restart user-service

# For frontend - check Vite dev server
cd frontend
npm run dev

# Check file watching limits (Linux)
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

### Problem: Build failures
**Symptoms**: Docker build or Go compilation errors

**Solutions**:
```bash
# Clean build cache
docker builder prune -f

# Rebuild from scratch
docker-compose build --no-cache

# Check Go module issues
cd backend/user-service
go mod tidy
go mod download
```

## Testing Issues

### Problem: E2E tests fail
**Symptoms**: Test suite reports failures

**Solutions**:
```bash
# Ensure services are healthy before testing
./scripts/e2e-test.sh

# Check test database is clean
docker-compose down -v
docker-compose up -d postgres redis

# Run tests with verbose output
cd tests/e2e
go test -v -timeout 10m ./...
```

### Problem: Database tests interfere with each other
**Symptoms**: Tests pass individually but fail when run together

**Solutions**:
```bash
# Use test database isolation
export DATABASE_URL="postgres://sports_user:sports_password@localhost:5432/sports_prediction_test?sslmode=disable"

# Run tests sequentially
go test -p 1 ./...
```

## Production Issues

### Problem: SSL certificate errors
**Symptoms**: HTTPS connections fail

**Solutions**:
```bash
# Check certificate validity
openssl x509 -in /path/to/cert.pem -text -noout

# Renew Let's Encrypt certificate
sudo certbot renew

# Check nginx configuration
nginx -t
sudo systemctl reload nginx
```

### Problem: High load issues
**Symptoms**: Services become unresponsive under load

**Solutions**:
```bash
# Scale services horizontally
docker-compose up -d --scale api-gateway=3

# Check database connections
docker-compose exec postgres psql -U sports_user -c "SELECT count(*) FROM pg_stat_activity;"

# Implement connection pooling
# Add PgBouncer configuration
```

## Diagnostic Commands

### Health Checks
```bash
# Check all service health
curl http://localhost:8080/health
curl http://localhost:8080/v1/auth/health
curl http://localhost:8080/v1/contests/health

# Check database connectivity
docker-compose exec postgres pg_isready -U sports_user

# Check Redis connectivity
docker-compose exec redis redis-cli ping
```

### Log Analysis
```bash
# View all service logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f user-service

# Search logs for errors
docker-compose logs | grep -i error

# Follow logs in real-time
docker-compose logs -f --tail=100
```

### Resource Monitoring
```bash
# Monitor container resources
docker stats

# Check disk usage
df -h
docker system df

# Monitor network traffic
docker-compose exec api-gateway netstat -tuln
```

## Getting Additional Help

If you encounter issues not covered here:

1. **Check service logs**: `docker-compose logs [service-name]`
2. **Verify configuration**: `docker-compose config`
3. **Test connectivity**: Use curl commands to test API endpoints
4. **Check system resources**: Ensure adequate memory and disk space
5. **Review documentation**: Check the [API documentation](../api/services-overview.md)
6. **Search existing issues**: Look for similar problems in project issues

## Emergency Recovery

### Complete System Reset
```bash
# Stop all services
docker-compose down -v

# Remove all containers and volumes
docker system prune -a -f --volumes

# Rebuild everything
make clean
make setup
make dev
```

### Database Recovery
```bash
# Backup current database (if possible)
docker-compose exec postgres pg_dump -U sports_user sports_prediction > backup.sql

# Reset database
docker-compose down -v
docker-compose up -d postgres

# Restore from backup (if needed)
cat backup.sql | docker-compose exec -T postgres psql -U sports_user -d sports_prediction
```

---

*For additional diagnostic tools and debugging procedures, see [Diagnostic Tools](diagnostic-tools.md).*
