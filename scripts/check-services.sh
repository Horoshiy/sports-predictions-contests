#!/bin/bash

# Check and restart failed services

echo "Checking service health..."

services=(
    "sports-api-gateway:8080"
    "sports-user-service:8084"
    "sports-contest-service:8085"
    "sports-prediction-service:8086"
    "sports-scoring-service:8087"
    "sports-service:8088"
    "sports-notification-service:8089"
    "sports-challenge-service:8090"
)

for service_port in "${services[@]}"; do
    service="${service_port%%:*}"
    port="${service_port##*:}"
    
    # Check if container is running
    if ! docker ps --format '{{.Names}}' | grep -q "^${service}$"; then
        echo "❌ $service is not running - starting..."
        docker-compose up -d "${service#sports-}"
    else
        # Check if service is responding (for services with health endpoints)
        if [ "$service" = "sports-api-gateway" ]; then
            if ! curl -sf http://localhost:$port/health > /dev/null 2>&1; then
                echo "⚠️  $service is not responding - restarting..."
                docker restart "$service"
            else
                echo "✅ $service is healthy"
            fi
        else
            echo "✅ $service is running"
        fi
    fi
done

echo ""
echo "Service check complete!"
