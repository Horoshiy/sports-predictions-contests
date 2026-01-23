#!/bin/bash

# Script to generate Protocol Buffer files for all services
# This must be run when protoc is available

set -e

echo "Generating Protocol Buffer files..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please install Protocol Buffers compiler:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install -y protobuf-compiler"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Navigate to backend directory
cd "$(dirname "$0")/../backend"

# Clean old generated files
echo "Cleaning old generated files..."
find shared/proto -name "*.pb.go" -type f -delete

# Generate proto files
echo "Generating proto files..."
protoc --proto_path=proto \
    --go_out=. \
    --go_opt=module=github.com/sports-prediction-contests/shared \
    --go-grpc_out=. \
    --go-grpc_opt=module=github.com/sports-prediction-contests/shared \
    proto/*.proto

# Move generated files to correct location
if [ -d "github.com" ]; then
    cp -r github.com/sports-prediction-contests/shared/proto/* shared/proto/
    rm -rf github.com
fi

echo "Proto generation complete!"
echo ""
echo "Generated files:"
find shared/proto -name "*.pb.go" -type f | sort
