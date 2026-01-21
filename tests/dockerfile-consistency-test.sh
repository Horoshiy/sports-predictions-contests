#!/bin/bash
# Test script to verify Dockerfile consistency

set -e

echo "=== Dockerfile Consistency Test ==="
echo ""

# Test 1: Check Go version consistency
echo "Test 1: Verifying Go version is 1.24 in all Dockerfiles..."
GO_VERSION_COUNT=$(grep -h "^FROM golang:" backend/*/Dockerfile bots/telegram/Dockerfile | grep -c "1.24-alpine")
EXPECTED_COUNT=9

if [ "$GO_VERSION_COUNT" -eq "$EXPECTED_COUNT" ]; then
    echo "✅ PASS: All $EXPECTED_COUNT Dockerfiles use golang:1.24-alpine"
else
    echo "❌ FAIL: Expected $EXPECTED_COUNT Dockerfiles with golang:1.24-alpine, found $GO_VERSION_COUNT"
    exit 1
fi

# Test 2: Check Alpine version consistency
echo ""
echo "Test 2: Verifying Alpine version is 3.19 in all Dockerfiles..."
ALPINE_VERSION_COUNT=$(grep -h "^FROM alpine:" backend/*/Dockerfile bots/telegram/Dockerfile | grep -c "3.19")

if [ "$ALPINE_VERSION_COUNT" -eq "$EXPECTED_COUNT" ]; then
    echo "✅ PASS: All $EXPECTED_COUNT Dockerfiles use alpine:3.19"
else
    echo "❌ FAIL: Expected $EXPECTED_COUNT Dockerfiles with alpine:3.19, found $ALPINE_VERSION_COUNT"
    exit 1
fi

# Test 3: Check for any :latest tags (should be none)
echo ""
echo "Test 3: Verifying no :latest tags are used..."
LATEST_COUNT=$(grep -h "^FROM" backend/*/Dockerfile bots/telegram/Dockerfile | grep -c ":latest" || true)

if [ "$LATEST_COUNT" -eq 0 ]; then
    echo "✅ PASS: No :latest tags found"
else
    echo "❌ FAIL: Found $LATEST_COUNT :latest tags (should be 0)"
    grep -n "^FROM.*:latest" backend/*/Dockerfile bots/telegram/Dockerfile
    exit 1
fi

# Test 4: Verify multi-stage build pattern
echo ""
echo "Test 4: Verifying multi-stage build pattern..."
BUILDER_COUNT=$(grep -h "AS builder" backend/*/Dockerfile bots/telegram/Dockerfile | wc -l | tr -d ' ')

if [ "$BUILDER_COUNT" -eq "$EXPECTED_COUNT" ]; then
    echo "✅ PASS: All $EXPECTED_COUNT Dockerfiles use multi-stage builds"
else
    echo "❌ FAIL: Expected $EXPECTED_COUNT multi-stage builds, found $BUILDER_COUNT"
    exit 1
fi

# Test 5: Verify CGO is disabled
echo ""
echo "Test 5: Verifying CGO_ENABLED=0 in all builds..."
CGO_COUNT=$(grep -h "CGO_ENABLED=0" backend/*/Dockerfile bots/telegram/Dockerfile | wc -l | tr -d ' ')

if [ "$CGO_COUNT" -eq "$EXPECTED_COUNT" ]; then
    echo "✅ PASS: All $EXPECTED_COUNT Dockerfiles disable CGO"
else
    echo "❌ FAIL: Expected $EXPECTED_COUNT with CGO_ENABLED=0, found $CGO_COUNT"
    exit 1
fi

echo ""
echo "=== All Tests Passed ✅ ==="
echo ""
echo "Summary:"
echo "  - Go version: golang:1.24-alpine (consistent across all services)"
echo "  - Alpine version: alpine:3.19 (consistent across all services)"
echo "  - No :latest tags (reproducible builds)"
echo "  - Multi-stage builds (optimized image size)"
echo "  - CGO disabled (static binaries)"
