#!/bin/sh

# Database migration script
# Usage: ./scripts/migrate.sh [up|down|create|version|force|drop]

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Defaults
MIGRATE_IMAGE=${MIGRATE_IMAGE:-migrate/migrate:v4.17.1}
MIGRATE_NETWORK=${MIGRATE_NETWORK:-cms-api_cms-network}

DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-cms}
DB_PASSWORD=${DB_PASSWORD:-cms_secret}
DB_NAME=${DB_NAME:-cms}
DB_SSL_MODE=${DB_SSL_MODE:-disable}

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"
MIGRATIONS_DIR="$(pwd)/migrations"

run_migrate() {
    docker run --rm \
        --network "$MIGRATE_NETWORK" \
        -v "$MIGRATIONS_DIR:/migrations" \
        "$MIGRATE_IMAGE" \
        -path /migrations \
        -database "$DATABASE_URL" "$@"
}

case "$1" in
    up)
        echo "Running migrations up..."
        run_migrate up
        ;;
    down)
        echo "Running migrations down..."
        run_migrate down ${2:-1}
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Error: migration name required"
            echo "Usage: ./scripts/migrate.sh create <migration_name>"
            exit 1
        fi
        echo "Creating migration: $2"
        run_migrate create -ext sql -dir /migrations "$2"
        ;;
    version)
        echo "Current migration version:"
        run_migrate version
        ;;
    force)
        if [ -z "$2" ]; then
            echo "Error: version number required"
            echo "Usage: ./scripts/migrate.sh force <version>"
            exit 1
        fi
        echo "Forcing version to: $2"
        run_migrate force "$2"
        ;;
    drop)
        echo "Dropping all tables..."
        printf "Are you sure? This will delete all data! (y/n) "
        read -r REPLY
        case "$REPLY" in
            y|Y)
                run_migrate drop -f
                ;;
            *)
                echo "Aborted."
                ;;
        esac
        ;;
    *)
        echo "Usage: ./scripts/migrate.sh [up|down|create|version|force|drop]"
        echo ""
        echo "Commands:"
        echo "  up              Run all pending migrations"
        echo "  down [n]        Rollback migrations (default: 1)"
        echo "  create <name>   Create a new migration"
        echo "  version         Show current migration version"
        echo "  force <v>       Force set version (use with caution)"
        echo "  drop            Drop all tables"
        exit 1
        ;;
esac
