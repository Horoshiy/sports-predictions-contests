#!/bin/bash
set -e

echo "=========================================="
echo "Frontend E2E Tests - Playwright"
echo "=========================================="

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

cd "$(dirname "$0")/.."

# Parse arguments
HEADED=false
UI_MODE=false
DEBUG=false
TEST_FILE=""

while [[ $# -gt 0 ]]; do
  case $1 in
    --headed) HEADED=true; shift ;;
    --ui) UI_MODE=true; shift ;;
    --debug) DEBUG=true; shift ;;
    *) TEST_FILE="$1"; shift ;;
  esac
done

echo -e "${YELLOW}Starting infrastructure services...${NC}"
docker-compose up -d postgres redis

echo -e "${YELLOW}Waiting for database...${NC}"
MAX_RETRIES=30
RETRY_COUNT=0
until docker-compose exec -T postgres pg_isready -U sports_user -d sports_prediction > /dev/null 2>&1; do
  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo -e "${RED}Database not ready${NC}"
    docker-compose down
    exit 1
  fi
  sleep 1
done

echo -e "${YELLOW}Starting microservices...${NC}"
docker-compose --profile services up -d
sleep 15

echo -e "${YELLOW}Checking API Gateway health...${NC}"
RETRY_COUNT=0
until curl -s http://localhost:8080/health > /dev/null 2>&1; do
  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo -e "${RED}API Gateway not ready after ${MAX_RETRIES} attempts${NC}"
    echo "Showing API Gateway logs:"
    docker-compose logs --tail=50 api-gateway
    docker-compose --profile services down
    exit 1
  fi
  sleep 2
done

echo -e "${GREEN}Services ready!${NC}"

# Run Playwright tests
cd frontend

if [ "$UI_MODE" = true ]; then
  echo -e "${YELLOW}Running Playwright in UI mode...${NC}"
  npm run test:e2e:ui
  TEST_EXIT=$?
elif [ "$DEBUG" = true ]; then
  echo -e "${YELLOW}Running Playwright in debug mode...${NC}"
  npm run test:e2e:debug $TEST_FILE
  TEST_EXIT=$?
elif [ "$HEADED" = true ]; then
  echo -e "${YELLOW}Running Playwright in headed mode...${NC}"
  npm run test:e2e:headed $TEST_FILE
  TEST_EXIT=$?
else
  echo -e "${YELLOW}Running Playwright tests...${NC}"
  if npm run test:e2e $TEST_FILE; then
    TEST_EXIT=0
    echo -e "${GREEN}All tests passed!${NC}"
  else
    TEST_EXIT=1
    echo -e "${RED}Some tests failed!${NC}"
  fi
fi

cd ..

echo -e "${YELLOW}Stopping services...${NC}"
docker-compose --profile services down
docker-compose down

exit $TEST_EXIT
