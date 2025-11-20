#!/bin/bash

set -euo pipefail

# Set logging variables
LOG_DIR="./build/logs"
mkdir -p "$LOG_DIR"
LOG_FILE="$LOG_DIR/build.log"

# Other variables
BINARY_NAME="sidebar"
BACKEND_DIR="./ui/backend"
ELECTRON_DIR="./ui"

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

# Bunch of stuff for DBUS issues
sudo service dbus start
export XDG_RUNTIME_DIR=/run/user/$(id -u)
sudo mkdir -p $XDG_RUNTIME_DIR
sudo chmod 700 $XDG_RUNTIME_DIR
sudo chown $(id -un):$(id -gn) $XDG_RUNTIME_DIR
export DBUS_SESSION_BUS_ADDRESS=unix:path=$XDG_RUNTIME_DIR/bus
dbus-daemon --session --address=$DBUS_SESSION_BUS_ADDRESS --nofork --nopidfile --syslog-only &

if [ -z $1 ]; then
  echo "Must provide build mode as argument"
  echo "Options are:"
  echo "  build              - build backend + vite + copy binary, no packaging"
  echo "  dev                - backend + vite dev server + electron"
  echo "  release-build      - full backend + vite build + electron-builder"
  exit 1
fi

# Create Build directory if it doesn't exist
mkdir -p ./build
mkdir -p ./ui/backend
mkdir -p ./ui/dist

case $1 in
  build )
    log INFO "Cleaning old builds"
    rm -f build/$BINARY_NAME
    rm -f $BACKEND_DIR/$BINARY_NAME

    log INFO "Exporting NO_AT_BRIDGE"
    export NO_AT_BRIDGE=1

    log INFO "Running Go tests"
    run_cmd go test ./... -v

    log INFO "Building Go backend"
    run_cmd GOOS=linux GOARCH=amd64 go build -o build/$BINARY_NAME .
    chmod +x build/$BINARY_NAME

    log INFO "Copying backend binary to Electron"
    mkdir -p $BACKEND_DIR
    cp build/$BINARY_NAME $BACKEND_DIR/$BINARY_NAME

    log INFO "Installing UI dependencies"
    run_cmd npm --prefix $ELECTRON_DIR install

    log INFO "Running Vite frontend build"
    run_cmd npm --prefix $ELECTRON_DIR run build

    log INFO "UI built successfully"
    ;;

  dev )
    log INFO "Cleaning old builds"
    rm -f build/$BINARY_NAME
    rm -f $BACKEND_DIR/$BINARY_NAME

    log INFO "Exporting NO_AT_BRIDGE"
    export NO_AT_BRIDGE=1

    log INFO "Running Go tests"
    run_cmd go test ./... -v

    log INFO "Building Go backend"
    run_cmd GOOS=linux GOARCH=amd64 go build -o build/$BINARY_NAME .
    chmod +x build/$BINARY_NAME

    log INFO "Copying backend binary to Electron..."
    mkdir -p $BACKEND_DIR
    cp build/$BINARY_NAME $BACKEND_DIR/$BINARY_NAME

    log INFO "Copying supporting directories"
    cp -r templates/ $BACKEND_DIR/

    log INFO "Starting Vite + Electron dev mode..."

    log INFO "Setting new LOG_DIR for npm commands"
    LOG_FILE=../build/logs/build.log

    cd $ELECTRON_DIR
    run_cmd npm install electron-is-dev

    # Run electron
    run_cmd NODE_ENV=development npm start

    LOG_FILE="$LOG_DIR/build.log"
    log INFO "Setting LOG_DIR back"
    ;;
    
  release-build )
    log INFO "Cleaning old builds"
    rm -f build/$BINARY_NAME build/$BINARY_NAME.exe

    log INFO "Running Go tests"
    run_cmd go test ./... -v

    VERSION=$(git describe --tags --always)
    COMMIT=$(git rev-parse --short HEAD)
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

    log INFO "Building backend with metadata"
    run_cmd go build -ldflags="-X main.version=$VERSION -X main.commit=$COMMIT -X main.buildTime=$BUILD_TIME"

    log INFO "Building Linux backend binary"
    run_cmd GOOS=linux GOARCH=amd64 go build -o build/$BINARY_NAME

    log INFO "Building Windows backend binary"
    run_cmd GOOS=windows GOARCH=amd64 go build -o build/$BINARY_NAME.exe

    log INFO "Copying backend binary into Electron project"
    mkdir -p $BACKEND_DIR
    cp build/$BINARY_NAME $BACKEND_DIR/$BINARY_NAME

    log INFO "Installing UI dependencies"
    run_cmd npm --prefix $ELECTRON_DIR install

    log INFO "Running Vite frontend build"
    run_cmd npm --prefix $ELECTRON_DIR run build

    log INFO "Packaging Electron app"
    run_cmd npm --prefix $ELECTRON_DIR run dist

    log INFO "Creating release archives"
    run_cmd zip build/backend-linux.zip build/$BINARY_NAME
    run_cmd zip build/backend-windows.zip build/$BINARY_NAME.exe

    # Electron outputs appear inside ui/dist/
    run_cmd zip -r build/electron-dist.zip ui/dist/
    ;;
  
esac

log INFO "Build finished"