#!/usr/bin/env bash
# Generate ERD from PostgreSQL database using Atlas
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$ROOT_DIR"

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Override Docker-internal hostname to localhost for host-side access
DB_HOST=localhost
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-cms}
DB_PASSWORD=${DB_PASSWORD:-cms_secret}
DB_NAME=${DB_NAME:-cms}
DB_SSL_MODE=${DB_SSL_MODE:-disable}

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

echo "Generating ERD from database..."

docker run --rm --network host \
    arigaio/atlas \
    schema inspect \
    --url "$DB_URL" \
    --format '{{ mermaid . }}' \
    > wiki/erd.mmd

echo "ERD exported to wiki/erd.mmd"
echo "View at: https://mermaid.live"
