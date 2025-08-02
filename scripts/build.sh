#!/bin/bash

# Build script for Portfolio Backend
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="portfolio-backend"
BUILD_DIR="build"
VERSION=${VERSION:-"1.0.0"}
COMMIT_SHA=${GITHUB_SHA:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo -e "${GREEN}Building Portfolio Backend...${NC}"
echo "Version: $VERSION"
echo "Commit: $COMMIT_SHA"
echo "Build Time: $BUILD_TIME"

# Create build directory
mkdir -p $BUILD_DIR

# Build flags
LDFLAGS="-X main.version=$VERSION -X main.commit=$COMMIT_SHA -X main.buildTime=$BUILD_TIME -w -s"

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -f $BUILD_DIR/$BINARY_NAME*

# Lint code
echo -e "${YELLOW}Running linters...${NC}"
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo -e "${YELLOW}golangci-lint not found, skipping linting${NC}"
fi

# Run tests
echo -e "${YELLOW}Running tests...${NC}"
go test -v -race -coverprofile=coverage.out ./...

# Build for local platform
echo -e "${YELLOW}Building for local platform...${NC}"
CGO_ENABLED=0 go build \
    -ldflags="$LDFLAGS" \
    -o $BUILD_DIR/$BINARY_NAME \
    cmd/api/main.go

# Build for Linux (common deployment target)
echo -e "${YELLOW}Building for Linux...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="$LDFLAGS" \
    -o $BUILD_DIR/${BINARY_NAME}-linux-amd64 \
    cmd/api/main.go

# Build for macOS
echo -e "${YELLOW}Building for macOS...${NC}"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
    -ldflags="$LDFLAGS" \
    -o $BUILD_DIR/${BINARY_NAME}-darwin-amd64 \
    cmd/api/main.go

# Build for Windows
echo -e "${YELLOW}Building for Windows...${NC}"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
    -ldflags="$LDFLAGS" \
    -o $BUILD_DIR/${BINARY_NAME}-windows-amd64.exe \
    cmd/api/main.go

echo -e "${GREEN}Build completed successfully!${NC}"
echo "Binaries available in $BUILD_DIR/"
ls -la $BUILD_DIR/

# Display binary information
echo -e "${GREEN}Binary information:${NC}"
file $BUILD_DIR/$BINARY_NAME