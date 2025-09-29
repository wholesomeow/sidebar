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
    rm -f ./build/sidebar

    # Run the linter here first; don't need CI yet
    echo "Linting codebase"
    golangci-lint run

    # Run the tests first so the binary wont build if tests fail
    # TODO: Split tests into app tests and cli tests
    echo "Running Sidebar tests"
    # go test ./... -v
    echo "Just kidding, there aren't any tests yet"
    echo "Will implement test-scripts soon... link to test-scripts here: https://bitfieldconsulting.com/posts/test-scripts"

    echo "Building Sidebar Binary"
    GOOS=linux GOARCH=amd64 go build -o build/sidebar .
    chmod +x build/sidebar
    echo "Binary built at: build/sidebar"
    ;;
  build-testless )
    # Clean previous builds
    rm -f ./build/sidebar

    # Run the linter here first; don't need CI yet
    echo "Linting codebase"
    golangci-lint run

    echo "Building Sidebar Binary"
    GOOS=linux GOARCH=amd64 go build -o build/sidebar .
    echo "Binary built at: build/sidebar"
    ;;
  release-build )
    # Clean previous builds
    rm -f ./build/sidebar.exe build/sidebar-linux

    # Run the linter here first; don't need CI yet
    echo "Linting codebase"
    golangci-lint run

    # Run the tests first so the binary wont build if tests fail
    # TODO: Split tests into app tests and cli tests
    echo "Running Sidebar tests"
    # go test ./... -v
    echo "Just kidding, there aren't any tests yet"
    echo "Will implement test-scripts soon... link to test-scripts here: https://bitfieldconsulting.com/posts/test-scripts"

    echo "Building Sidebar Binary"

    # Including Build Metadata
    VERSION=$(git describe --tags --always)
    COMMIT=$(git rev-parse --short HEAD)
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

    go build -ldflags="-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME"

    # Linux 64-bit
    GOOS=linux GOARCH=amd64 go build -o build/sidebar-linux
    chmod +x build/sidebar-linux
    echo "Binary built at: build/sidebar-linux"

    # Windows 64-bit
    GOOS=windows GOARCH=amd64 go build -o build/sidebar.exe
    chmod +x build/sidebar.exe
    echo "Binary built at: build/sidebar.exe"

    echo "Compressing builds..."
    zip build/sidebar-linux.zip build/sidebar-linux
    zip build/sidebar-windows.zip build/sidebar.exe
    ;;
  
esac

echo "Build finished"