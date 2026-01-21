#!/bin/bash

# Dependency Version Consistency Checker
# Ensures critical shared dependencies maintain consistent versions across all microservices

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Validate we're in the project root
if [ ! -d "backend" ]; then
    echo -e "${RED}❌ Error: backend directory not found${NC}"
    echo "Please run this script from the project root directory"
    exit 1
fi

echo "🔍 Checking dependency version consistency across services..."
echo ""

# Critical dependencies that should be consistent
CRITICAL_DEPS=(
    "gorm.io/gorm"
    "google.golang.org/grpc"
    "google.golang.org/protobuf"
)

ERRORS=0

for dep in "${CRITICAL_DEPS[@]}"; do
    echo "Checking: $dep"
    
    # Get all versions of this dependency (field 3 contains version when grep includes filename)
    versions=$(grep "$dep " backend/*/go.mod 2>/dev/null | grep -v "^Binary" | awk '{print $3}' | sort -u)
    version_count=$(echo "$versions" | grep -c "^v" || echo "0")
    
    if [ "$version_count" -gt 1 ]; then
        echo -e "${RED}❌ INCONSISTENT VERSIONS FOUND${NC}"
        echo "   Versions in use:"
        grep "$dep " backend/*/go.mod | grep -v "^Binary" | while read -r line; do
            echo "   - $line"
        done
        echo ""
        ERRORS=$((ERRORS + 1))
    else
        # Get the actual version for display (field after dependency name, before any comments)
        actual_version=$(grep "$dep " backend/*/go.mod 2>/dev/null | head -1 | awk '{for(i=1;i<=NF;i++) if($i ~ /^v[0-9]/) print $i}')
        if [ -z "$actual_version" ]; then
            actual_version="$versions"
        fi
        echo -e "${GREEN}✅ Consistent: $actual_version${NC}"
    fi
    echo ""
done

if [ $ERRORS -gt 0 ]; then
    echo -e "${RED}❌ Found $ERRORS dependency inconsistencies${NC}"
    echo ""
    echo "To fix, update all services to use the same version:"
    echo "  cd backend/<service-name>"
    echo "  go get <dependency>@<version>"
    echo "  go mod tidy"
    exit 1
else
    echo -e "${GREEN}✅ All critical dependencies are consistent across services${NC}"
    exit 0
fi
