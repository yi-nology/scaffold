#!/bin/bash

# Scaffold Deployment Script with Chinese Mirrors Optimization
# Author: murphy-hz-init
# Description: Deploy scaffold application with optimized Chinese mirror sources

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$SCRIPT_DIR"

echo "🚀 Starting Scaffold Deployment with Chinese Mirrors..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    print_success "Docker and Docker Compose are installed"
}

# Determine which docker compose command to use
get_docker_compose_cmd() {
    if command -v docker-compose &> /dev/null; then
        echo "docker-compose"
    elif docker compose version &> /dev/null; then
        echo "docker compose"
    else
        print_error "Cannot find docker compose command"
        exit 1
    fi
}

DOCKER_COMPOSE=$(get_docker_compose_cmd)

# Configure Docker daemon for Chinese mirrors (if not already configured)
configure_docker_mirrors() {
    print_status "Checking Docker mirror configuration..."
    
    # Check if Docker daemon.json exists
    if [ ! -f "/etc/docker/daemon.json" ]; then
        print_warning "Docker daemon.json not found. Creating with Chinese mirrors..."
        sudo mkdir -p /etc/docker
        sudo tee /etc/docker/daemon.json > /dev/null <<EOF
{
    "registry-mirrors": [
        "https://hub-mirror.c.163.com",
        "https://mirror.baidubce.com",
        "https://docker.mirrors.ustc.edu.cn"
    ]
}
EOF
        print_status "Restarting Docker daemon..."
        sudo systemctl restart docker
        print_success "Docker mirrors configured successfully"
    else
        print_status "Docker daemon.json already exists"
    fi
}

# Pull base images with Chinese mirrors
pull_base_images() {
    print_status "Pre-pulling base images with Chinese mirrors..."
    
    # Array of base images to pull
    base_images=(
        "golang:1.25-alpine"
        "node:18-alpine" 
        "alpine:latest"
        "nginx:alpine"
    )
    
    for image in "${base_images[@]}"; do
        print_status "Pulling $image..."
        docker pull "$image" || print_warning "Failed to pull $image"
    done
    
    print_success "Base images pulled"
}

# Copy environment file if it doesn't exist
setup_environment() {
    if [ ! -f ".env" ]; then
        print_status "Creating .env file from example..."
        cp .env.example .env
        print_success "Environment file created. Please review and modify .env as needed."
    fi
}

# Enhanced build function with Chinese mirrors
enhanced_build() {
    print_status "Building services with Chinese mirrors optimization..."
    
    # Enable BuildKit for better caching
    export DOCKER_BUILDKIT=1
    
    # Build with parallel builds and Chinese mirrors
    $DOCKER_COMPOSE build \
        --parallel \
        --build-arg GOPROXY=https://goproxy.cn,direct \
        --build-arg GOSUMDB=sum.golang.google.cn
    
    print_success "Services built successfully"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up          Start all services"
    echo "  down        Stop all services"
    echo "  rebuild     Rebuild and start services with Chinese mirrors"
    echo "  logs        Show service logs"
    echo "  status      Show service status"
    echo "  clean       Clean up containers and volumes"
    echo "  config      Manage configuration files"
    echo "  optimize    Configure Docker mirrors and pre-pull images"
    echo "  help        Show this help message"
    echo ""
    echo "Optimization Features:"
    echo "  - Chinese Docker registry mirrors"
    echo "  - Go proxy configuration (goproxy.cn)"
    echo "  - NPM registry mirror (npmmirror.com)"
    echo "  - Alpine package mirror (aliyun.com)"
    echo "  - Parallel builds with BuildKit"
    echo ""
    echo "Configuration Management:"
    echo "  The application now supports YAML configuration files."
    echo "  Default config file: ../configs/config.yaml"
    echo "  You can specify custom config with SCAFFOLD_CONFIG environment variable."
}

# Function to manage configuration
manage_config() {
    print_status "Configuration Management"
    echo ""
    
    CONFIG_FILE="../configs/config.yaml"
    
    case "${1:-}" in
        "show")
            if [ -f "$CONFIG_FILE" ]; then
                print_status "Current configuration ($CONFIG_FILE):"
                echo "----------------------------------------"
                cat "$CONFIG_FILE"
                echo "----------------------------------------"
            else
                print_error "Configuration file not found: $CONFIG_FILE"
                print_status "Run '$0 config init' to create default configuration"
            fi
            ;;
        "init")
            if [ -f "$CONFIG_FILE" ]; then
                print_warning "Configuration file already exists: $CONFIG_FILE"
                print_status "Use '$0 config show' to view current config"
                print_status "Use '$0 config edit' to modify config"
            else
                print_status "Creating default configuration file..."
                mkdir -p "$(dirname "$CONFIG_FILE")"
                
                # Create default config from the application's default config
                cat > "$CONFIG_FILE" <<EOF
# Scaffold Application Configuration
# This file controls the application behavior when running in Docker

server:
  host: "0.0.0.0"
  port: 9090
  cache_dir: "./cache"
  debug: false
  read_timeout: 30
  write_timeout: 30

database:
  enable: true
  driver: "sqlite"
  host: "localhost"
  port: 3306
  user: "root"
  password: ""
  db_name: "/app/data/scaffold.db"
  ssl_mode: "disable"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

redis:
  enable: false
  host: "redis"
  port: 6379
  password: ""
  db: 0

security:
  access_key: ""

log:
  level: "info"
  filename: ""
  max_size: 100
  max_backups: 3
  max_age: 7
  compress: true

template:
  storage_path: "./templates"
  allowed_dirs:
    - "./templates"

app:
  name: "scaffold"
  version: "1.0.0"
EOF
                print_success "Default configuration created: $CONFIG_FILE"
                print_status "You can now edit this file to customize your deployment"
            fi
            ;;
        "edit")
            if [ ! -f "$CONFIG_FILE" ]; then
                print_error "Configuration file not found. Creating first..."
                manage_config init
            fi
            
            # Use default editor or nano/vim
            if command -v nano &> /dev/null; then
                nano "$CONFIG_FILE"
            elif command -v vim &> /dev/null; then
                vim "$CONFIG_FILE"
            else
                print_status "Please edit the configuration file manually:"
                echo "   $CONFIG_FILE"
            fi
            ;;
        "validate")
            if [ -f "$CONFIG_FILE" ]; then
                print_status "Validating configuration..."
                # Here you could add validation logic if needed
                print_success "Configuration file exists and is readable"
            else
                print_error "Configuration file not found: $CONFIG_FILE"
            fi
            ;;
        *)
            echo "Usage: $0 config [show|init|edit|validate]"
            echo ""
            echo "Commands:"
            echo "  show      Display current configuration"
            echo "  init      Create default configuration file"
            echo "  edit      Edit configuration file"
            echo "  validate  Validate configuration file"
            ;;
    esac
}

# Parse arguments
case "${1:-}" in
    up)
        print_status "Starting Scaffold services..."
        shift
        $DOCKER_COMPOSE up -d "$@"
        print_success "Services started successfully!"
        echo "Application: http://localhost:${FRONTEND_PORT:-80}"
        ;;
    down)
        print_status "Stopping Scaffold services..."
        $DOCKER_COMPOSE down
        print_success "Services stopped."
        ;;
    rebuild)
        print_status "Rebuilding Scaffold services with Chinese mirrors..."
        $DOCKER_COMPOSE down
        enhanced_build
        $DOCKER_COMPOSE up -d
        print_success "Services rebuilt and started!"
        ;;
    logs)
        print_status "Showing service logs..."
        $DOCKER_COMPOSE logs -f
        ;;
    status)
        print_status "Service status:"
        $DOCKER_COMPOSE ps
        ;;
    clean)
        print_status "Cleaning up containers and volumes..."
        $DOCKER_COMPOSE down -v --remove-orphans
        docker system prune -f
        print_success "Cleanup completed."
        ;;
    optimize)
        check_docker
        configure_docker_mirrors
        pull_base_images
        print_success "Docker optimization completed!"
        ;;
    config)
        shift
        manage_config "$@"
        ;;
    help|"")
        show_usage
        ;;
    *)
        print_error "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac