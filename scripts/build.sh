#!/bin/bash

set -euo pipefail

# Set logging variables
LOG_DIR="./logs"
mkdir -p "$LOG_DIR"
LOG_FILE="$LOG_DIR/build.log"

# Logging function
log() {
  local level="$1"; shift
  local msg="$*"
  local timestamp
  timestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)
  echo "[$timestamp] [$level] $msg" | tee -a "$LOG_FILE"
}

# Run command and log both stdout + stderr
run_cmd() {
  local cmd="$*"
  log INFO "Running: $cmd"
  # Pipe both stdout and stderr to tee
  eval "$cmd" 2>&1 | tee -a "$LOG_FILE"
  local status=${PIPESTATUS[0]}
  if [[ $status -ne 0 ]]; then
    log ERROR "Command failed: $cmd (exit $status)"
    exit $status
  fi
}

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
    log INFO "Clean previous builds"
    rm -f ./build/sidebar

    # Run the linter here first; don't need CI yet
    log INFO "Linting codebase"
    run_cmd golangci-lint run

    # Run the tests first so the binary wont build if tests fail
    # TODO: Split tests into app tests and cli tests
    log INFO "Running Sidebar tests"
    log INFO "There aren't any CLI tests yet"
    log INFO "Will implement test-scripts soon... link to test-scripts here: https://bitfieldconsulting.com/posts/test-scripts"
    log INFO "Running the other tests now"
    run_cmd go test ./... -v

    log INFO "Building Sidebar Binary"
    run_cmd GOOS=linux GOARCH=amd64 go build -o build/sidebar .
    chmod +x build/sidebar
    log INFO "Binary built at: build/sidebar"
    ;;
  build-testless )
    log INFO "Clean previous builds"
    rm -f ./build/sidebar

    # Run the linter here first; don't need CI yet
    log INFO "Linting codebase"
    run_cmd golangci-lint run

    log INFO "Building Sidebar Binary"
    run_cmd GOOS=linux GOARCH=amd64 go build -o build/sidebar .
    chmod +x build/sidebar
    log INFO "Binary built at: build/sidebar"
    ;;
  release-build )
    log INFO "Clean previous builds"
    rm -f ./build/sidebar.exe build/sidebar-linux

    # Run the linter here first; don't need CI yet
    log INFO "Linting codebase"
    run_cmd golangci-lint run

    # Run the tests first so the binary wont build if tests fail
    # TODO: Split tests into app tests and cli tests
    log INFO "Running Sidebar tests"
    log INFO "Just kidding, there aren't any tests yet"
    log INFO "Will implement test-scripts soon... link to test-scripts here: https://bitfieldconsulting.com/posts/test-scripts"
    log INFO "Running the other tests now"
    run_cmd go test ./... -v

    log INFO "Building Sidebar Binary"

    # Including Build Metadata
    VERSION=$(git describe --tags --always)
    COMMIT=$(git rev-parse --short HEAD)
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

    run_cmd go build -ldflags="-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME"

    # Linux 64-bit
    run_cmd GOOS=linux GOARCH=amd64 go build -o build/sidebar-linux
    chmod +x build/sidebar-linux
    log INFO "Binary built at: build/sidebar-linux"

    # Windows 64-bit
    run_cmd GOOS=windows GOARCH=amd64 go build -o build/sidebar.exe
    chmod +x build/sidebar.exe
    log INFO "Binary built at: build/sidebar.exe"

    log INFO "Compressing builds..."
    run_cmd zip build/sidebar-linux.zip build/sidebar-linux
    run_cmd zip build/sidebar-windows.zip build/sidebar.exe
    ;;
  
esac

log INFO "Build finished"