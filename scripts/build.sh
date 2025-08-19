#!/bin/bash

set -e

if [ -z $1 ]; then
  echo "Must provide build mode as argument"
  echo "Options are:"
  echo "  - build               Builds the main binary and runs its tests"
  echo "  - build-testless      Build the main binary without running tests"
  echo "  - release-build       Builds the main binary for windows and linux"
  exit 1
fi

# Create Build directory if it doesn't exist
mkdir -p ./build

case $1 in
  build )
    # Clean previous builds
    rm -f ./build/chatctl

    # Run the tests first so the binary wont build if tests fail
    echo "Running ChatWrapper tests"
    # go test ./... -v
    echo "Just kidding, there aren't any tests yet"
    echo "Will implement test-scripts soon... link to test-scripts here: https://bitfieldconsulting.com/posts/test-scripts"

    echo "Building ChatWrapper Binary"
    GOOS=linux GOARCH=amd64 go build -o build/chatctl .
    chmod +x build/chatctl
    echo "Binary built at: build/chatctl"
    ;;
  build-testless )
    # Clean previous builds
    rm -f ./build/chatctl

    echo "Building ChatWrapper Binary"
    GOOS=linux GOARCH=amd64 go build -o build/chatctl .
    echo "Binary built at: build/chatctl"
    ;;
  release-build )
    # Clean previous builds
    rm -f ./build/chatctl.exe build/chatctl-linux

    # Run the tests first so the binary wont build if tests fail
    echo "Running ChatWrapper tests"
    # go test ./... -v
    echo "Just kidding, there aren't any tests yet"
    echo "Will implement test-scripts soon... link to test-scripts here: https://bitfieldconsulting.com/posts/test-scripts"

    echo "Building ChatWrapper Binary"

    # Including Build Metadata
    VERSION=$(git describe --tags --always)
    COMMIT=$(git rev-parse --short HEAD)
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

    go build -ldflags="-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME"

    # Linux 64-bit
    GOOS=linux GOARCH=amd64 go build -o build/chatctl-linux
    chmod +x build/chatctl-linux
    echo "Binary built at: build/chatctl-linux"

    # Windows 64-bit
    GOOS=windows GOARCH=amd64 go build -o build/chatctl.exe
    chmod +x build/chatctl.exe
    echo "Binary built at: build/chatctl.exe"

    echo "Compressing builds..."
    zip build/chatctl-linux.zip build/chatctl-linux
    zip build/chatctl-windows.zip build/chatctl.exe
    ;;
  
esac

echo "Build finished"