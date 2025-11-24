#!/bin/bash

# Neonex Core Development Helper Script

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Functions
dev() {
    echo -e "${CYAN}🔥 Starting development server with hot reload...${NC}"
    air
}

serve() {
    echo -e "${GREEN}🚀 Starting server...${NC}"
    go run main.go
}

build() {
    echo -e "${YELLOW}🔨 Building application...${NC}"
    go build -o neonex .
    echo -e "${GREEN}✅ Build complete: ./neonex${NC}"
}

build-cli() {
    echo -e "${YELLOW}🔨 Building CLI...${NC}"
    go build -o neonex ./cmd/neonex
    echo -e "${GREEN}✅ CLI built: ./neonex${NC}"
}

install-air() {
    echo -e "${CYAN}📦 Installing Air...${NC}"
    go install github.com/air-verse/air@latest
    echo -e "${GREEN}✅ Air installed${NC}"
}

install-deps() {
    echo -e "${CYAN}📦 Installing dependencies...${NC}"
    go mod download
    go mod tidy
    echo -e "${GREEN}✅ Dependencies installed${NC}"
}

run-tests() {
    echo -e "${CYAN}🧪 Running tests...${NC}"
    go test -v ./...
}

clean() {
    echo -e "${YELLOW}🧹 Cleaning...${NC}"
    rm -rf tmp
    rm -f neonex neonex.exe
    rm -f build-errors.log
    echo -e "${GREEN}✅ Clean complete${NC}"
}

fmt() {
    echo -e "${CYAN}🎨 Formatting code...${NC}"
    go fmt ./...
    echo -e "${GREEN}✅ Code formatted${NC}"
}

new-module() {
    if [ -z "$1" ]; then
        echo -e "${RED}❌ Module name required: new-module modulename${NC}"
        return 1
    fi
    echo -e "${CYAN}📦 Creating module: $1${NC}"
    ./neonex module create "$1"
}

list-modules() {
    echo -e "${CYAN}📦 Available modules:${NC}"
    ./neonex module list
}

show-help() {
    cat << EOF

Neonex Core - Development Commands
===================================

Quick Start:
  dev              - Start with hot reload (air)
  serve            - Start without hot reload
  
Build:
  build            - Build main application
  build-cli        - Build CLI tool
  clean            - Clean build artifacts
  
Development:
  install-air      - Install Air for hot reload
  install-deps     - Install Go dependencies
  fmt              - Format all Go code
  run-tests        - Run all tests
  
Modules:
  new-module name  - Create new module
  list-modules     - List all modules
  
Examples:
  $ dev
  $ new-module product
  $ build-cli
  $ list-modules

EOF
}

# Parse command
case "$1" in
    dev)
        dev
        ;;
    serve)
        serve
        ;;
    build)
        build
        ;;
    build-cli)
        build-cli
        ;;
    install-air)
        install-air
        ;;
    install-deps)
        install-deps
        ;;
    test)
        run-tests
        ;;
    clean)
        clean
        ;;
    fmt)
        fmt
        ;;
    new-module)
        new-module "$2"
        ;;
    list-modules)
        list-modules
        ;;
    help|--help|-h)
        show-help
        ;;
    *)
        show-help
        ;;
esac
