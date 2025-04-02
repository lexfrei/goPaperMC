#!/bin/bash

set -e

echo "Building goPaperMC CLI..."

# Determine current directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Go to project root directory
cd "$SCRIPT_DIR"

# Get Go version
GO_VERSION=$(go version | awk '{print $3}')
echo "Using $GO_VERSION"

# Set build variables
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Compile CLI client
echo "Building Linux amd64 binary..."
go build -o bin/papermc-linux-amd64 ./cmd/papermc

# For macOS
if command -v go &> /dev/null; then
    echo "Building macOS amd64 binary..."
    export GOOS=darwin
    export GOARCH=amd64
    go build -o bin/papermc-darwin-amd64 ./cmd/papermc
    
    echo "Building macOS arm64 binary..."
    export GOOS=darwin
    export GOARCH=arm64
    go build -o bin/papermc-darwin-arm64 ./cmd/papermc
fi

# For Windows
echo "Building Windows amd64 binary..."
export GOOS=windows
export GOARCH=amd64
go build -o bin/papermc-windows-amd64.exe ./cmd/papermc

echo "Build complete!"
echo "Binaries are available in the bin/ directory"
