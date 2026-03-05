#!/bin/bash

# Scaffold Deployment Script
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$SCRIPT_DIR"

echo "🚀 Starting Scaffold Deployment..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Determine which docker compose command to use
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
else
    echo "❌ Cannot find docker compose command"
    exit 1
fi

# Copy environment file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "📋 Creating .env file from example..."
    cp .env.example .env
    echo "✅ Environment file created. Please review and modify .env as needed."
fi

# Function to show usage
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up          Start all services"
    echo "  down        Stop all services"
    echo "  rebuild     Rebuild and start services"
    echo "  logs        Show service logs"
    echo "  status      Show service status"
    echo "  clean       Clean up containers and volumes"
    echo "  config      Manage configuration files"
    echo "  help        Show this help message"
    echo ""
    echo "Configuration Management:"
    echo "  The application now supports YAML configuration files."
    echo "  Default config file: ../configs/config.yaml"
    echo "  You can specify custom config with SCAFFOLD_CONFIG environment variable."
}

# Function to manage configuration
manage_config() {
    echo "🔧 Configuration Management"
    echo ""
    
    CONFIG_FILE="../configs/config.yaml"
    
    case "${1:-}" in
        "show")
            if [ -f "$CONFIG_FILE" ]; then
                echo "📄 Current configuration ($CONFIG_FILE):"
                echo "----------------------------------------"
                cat "$CONFIG_FILE"
                echo "----------------------------------------"
            else
                echo "❌ Configuration file not found: $CONFIG_FILE"
                echo "💡 Run '$0 config init' to create default configuration"
            fi
            ;;
        "init")
            if [ -f "$CONFIG_FILE" ]; then
                echo "⚠️  Configuration file already exists: $CONFIG_FILE"
                echo "💡 Use '$0 config show' to view current config"
                echo "💡 Use '$0 config edit' to modify config"
            else
                echo "📋 Creating default configuration file..."
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
                echo "✅ Default configuration created: $CONFIG_FILE"
                echo "💡 You can now edit this file to customize your deployment"
            fi
            ;;
        "edit")
            if [ ! -f "$CONFIG_FILE" ]; then
                echo "❌ Configuration file not found. Creating first..."
                manage_config init
            fi
            
            # Use default editor or nano/vim
            if command -v nano &> /dev/null; then
                nano "$CONFIG_FILE"
            elif command -v vim &> /dev/null; then
                vim "$CONFIG_FILE"
            else
                echo "📝 Please edit the configuration file manually:"
                echo "   $CONFIG_FILE"
            fi
            ;;
        "validate")
            if [ -f "$CONFIG_FILE" ]; then
                echo "🔍 Validating configuration..."
                # Here you could add validation logic if needed
                echo "✅ Configuration file exists and is readable"
            else
                echo "❌ Configuration file not found: $CONFIG_FILE"
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
        echo "🐳 Starting Scaffold services..."
        shift
        $DOCKER_COMPOSE up -d "$@"
        echo "✅ Services started successfully!"
        echo "Application: http://localhost:${FRONTEND_PORT:-80}"
        ;;
    down)
        echo "🛑 Stopping Scaffold services..."
        $DOCKER_COMPOSE down
        echo "✅ Services stopped."
        ;;
    rebuild)
        echo "🔨 Rebuilding Scaffold services..."
        $DOCKER_COMPOSE down
        $DOCKER_COMPOSE build --no-cache
        $DOCKER_COMPOSE up -d
        echo "✅ Services rebuilt and started!"
        ;;
    logs)
        echo "📋 Showing service logs..."
        $DOCKER_COMPOSE logs -f
        ;;
    status)
        echo "📊 Service status:"
        $DOCKER_COMPOSE ps
        ;;
    clean)
        echo "🧹 Cleaning up containers and volumes..."
        $DOCKER_COMPOSE down -v --remove-orphans
        docker system prune -f
        echo "✅ Cleanup completed."
        ;;
    config)
        shift
        manage_config "$@"
        ;;
    help|"")
        show_usage
        ;;
    *)
        echo "❌ Unknown command: $1"
        show_usage
        exit 1
        ;;
esac