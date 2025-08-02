#!/bin/bash

# Database migration script for Portfolio Backend
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"portfolio_user"}
DB_PASSWORD=${DB_PASSWORD:-""}
DB_NAME=${DB_NAME:-"portfolio_db"}
MIGRATIONS_PATH=${MIGRATIONS_PATH:-"migrations"}

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if migrate tool is installed
    if ! command -v migrate &> /dev/null; then
        log_error "golang-migrate is not installed."
        log_error "Install it with: go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
    
    # Check if migrations directory exists
    if [ ! -d "$MIGRATIONS_PATH" ]; then
        log_error "Migrations directory '$MIGRATIONS_PATH' does not exist."
        exit 1
    fi
    
    # Check if database credentials are provided
    if [ -z "$DB_PASSWORD" ]; then
        log_error "DB_PASSWORD is not set. Please provide database password."
        exit 1
    fi
}

build_connection_string() {
    CONNECTION_STRING="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"
}

migrate_up() {
    log_info "Running migrations up..."
    migrate -path $MIGRATIONS_PATH -database "$CONNECTION_STRING" up
    log_info "Migrations completed successfully."
}

migrate_down() {
    log_info "Rolling back migrations..."
    log_warn "This will roll back the last migration. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        migrate -path $MIGRATIONS_PATH -database "$CONNECTION_STRING" down 1
        log_info "Migration rollback completed."
    else
        log_info "Migration rollback cancelled."
    fi
}

migrate_down_all() {
    log_error "WARNING: This will roll back ALL migrations and destroy all data!"
    log_error "Are you absolutely sure? Type 'YES' to confirm:"
    read -r response
    if [[ "$response" == "YES" ]]; then
        migrate -path $MIGRATIONS_PATH -database "$CONNECTION_STRING" down -all
        log_info "All migrations rolled back."
    else
        log_info "Migration rollback cancelled."
    fi
}

migrate_force() {
    if [ -z "$2" ]; then
        log_error "Please provide a migration version to force."
        exit 1
    fi
    
    log_warn "Forcing migration to version $2..."
    migrate -path $MIGRATIONS_PATH -database "$CONNECTION_STRING" force $2
    log_info "Migration forced to version $2."
}

migrate_version() {
    log_info "Current migration version:"
    migrate -path $MIGRATIONS_PATH -database "$CONNECTION_STRING" version
}

migrate_create() {
    if [ -z "$2" ]; then
        log_error "Please provide a migration name."
        exit 1
    fi
    
    log_info "Creating new migration: $2"
    migrate create -ext sql -dir $MIGRATIONS_PATH -seq $2
    log_info "Migration files created successfully."
}

show_help() {
    echo "Database Migration Script for Portfolio Backend"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up              Run all pending migrations"
    echo "  down            Roll back the last migration"
    echo "  down-all        Roll back all migrations (DESTRUCTIVE)"
    echo "  force VERSION   Force migration to specific version"
    echo "  version         Show current migration version"
    echo "  create NAME     Create new migration files"
    echo "  help            Show this help message"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST         Database host (default: localhost)"
    echo "  DB_PORT         Database port (default: 3306)"
    echo "  DB_USER         Database user (default: portfolio_user)"
    echo "  DB_PASSWORD     Database password (required)"
    echo "  DB_NAME         Database name (default: portfolio_db)"
    echo "  MIGRATIONS_PATH Migration files path (default: migrations)"
}

# Main script logic
main() {
    case $1 in
        up)
            check_prerequisites
            build_connection_string
            migrate_up
            ;;
        down)
            check_prerequisites
            build_connection_string
            migrate_down
            ;;
        down-all)
            check_prerequisites
            build_connection_string
            migrate_down_all
            ;;
        force)
            check_prerequisites
            build_connection_string
            migrate_force $@
            ;;
        version)
            check_prerequisites
            build_connection_string
            migrate_version
            ;;
        create)
            migrate_create $@
            ;;
        help|--help|-h)
            show_help
            ;;
        "")
            log_error "No command provided."
            show_help
            exit 1
            ;;
        *)
            log_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main $@