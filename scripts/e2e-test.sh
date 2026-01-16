#!/bin/bash
set -e

echo "=========================================="
echo "Sports Prediction Contests - E2E Tests"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Change to project root
cd "$(dirname "$0")/.."

echo -e "${YELLOW}Starting infrastructure services...${NC}"
docker-compose up -d postgres redis

echo -e "${YELLOW}Waiting for database to be ready...${NC}"
MAX_DB_RETRIES=30
DB_RETRY_COUNT=0
until docker-compose exec -T postgres pg_isready -U sports_user -d sports_prediction > /dev/null 2>&1; do
    DB_RETRY_COUNT=$((DB_RETRY_COUNT + 1))
    if [ $DB_RETRY_COUNT -ge $MAX_DB_RETRIES ]; then
        echo -e "${RED}Database not ready after ${MAX_DB_RETRIES} attempts${NC}"
        docker-compose down
        exit 1
    fi
    echo "Waiting for PostgreSQL... (attempt $DB_RETRY_COUNT/$MAX_DB_RETRIES)"
    sleep 1
done
echo -e "${GREEN}Database is ready!${NC}"

echo -e "${YELLOW}Starting all microservices...${NC}"
docker-compose --profile services up -d

echo -e "${YELLOW}Waiting for services to start...${NC}"
sleep 15

# Wait for API Gateway to be healthy
echo -e "${YELLOW}Checking API Gateway health...${NC}"
MAX_RETRIES=30
RETRY_COUNT=0

until curl -s http://localhost:8080/health > /dev/null 2>&1; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        echo -e "${RED}API Gateway not ready after ${MAX_RETRIES} attempts${NC}"
        echo "Showing service logs:"
        docker-compose logs --tail=50 api-gateway
        docker-compose --profile services down
        exit 1
    fi
    echo "Waiting for API Gateway... (attempt $RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

echo -e "${GREEN}API Gateway is healthy!${NC}"

echo -e "${YELLOW}Running E2E tests...${NC}"
cd tests/e2e

# Run tests with verbose output
if go test -tags=e2e -v -timeout 5m ./...; then
    TEST_EXIT=0
    echo -e "${GREEN}=========================================="
    echo "All E2E tests passed!"
    echo -e "==========================================${NC}"
else
    TEST_EXIT=1
    echo -e "${RED}=========================================="
    echo "Some E2E tests failed!"
    echo -e "==========================================${NC}"
fi

cd ../..

echo -e "${YELLOW}Stopping services...${NC}"
docker-compose --profile services down

echo -e "${YELLOW}Stopping infrastructure...${NC}"
docker-compose down

exit $TEST_EXIT
