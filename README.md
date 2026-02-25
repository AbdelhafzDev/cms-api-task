# CMS API

Modular monolith backend for content management, built with Go.

## Tech Stack

- **Go 1.25** with **Chi v5** (router) and **Uber Fx** (DI)
- **PostgreSQL 17** with sqlx
- **Meilisearch 1.13** for full-text search
- **Redis 7** for caching
- **Zap** for structured logging
- **JWT (RSA)** for authentication
- **OpenTelemetry** for tracing
- **Docker + Docker Compose** for containerization

## Getting Started

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Build and start all services
make build && make up

# 3. Run database migrations
make migrate-up
```

The API is now running at `http://localhost:8080`.

To enable optional tools (Adminer, Jaeger): `docker compose --profile tools up -d`

## Services

- **cms-api** — Application server (HTTP `:8080`, gRPC `:9090`)
- **cms-postgres** — PostgreSQL database (`:5432`)
- **cms-meilisearch** — Full-text search engine (`:7700`)
- **cms-redis** — Cache layer (`:6379`)
- **cms-jaeger** — Tracing UI (`:16686`, optional)
- **cms-adminer** — DB admin UI (`:8081`, optional)

## Commands

```bash
make help             # Show all commands
make build            # Build Docker dev image
make up               # Start all services
make down             # Stop all services
make test             # Run tests
make lint             # Run golangci-lint
make logs             # Tail container logs
make shell            # Enter container shell
make restart          # Restart container
make migrate-up       # Run pending migrations
make migrate-down     # Rollback last migration
```

## API Documentation

- **Swagger UI** — http://localhost:8080/swagger/
- **OpenAPI Spec** — `docs/http/openapi.yaml`

## Wiki

The `wiki/` folder contains:

- **`CMS_API_collection.json`** — Postman / Apidog collection with all endpoints, auto-token management, and i18n support (`Accept-Language: en|ar`)
- **`architecture.excalidraw`** — Architecture user journey diagram (open at [excalidraw.com](https://excalidraw.com))

## Project Structure

```
internal/
  modules/          # Feature modules (auth, program, discovery, worker, importer)
  transport/        # HTTP (Chi) and gRPC servers
  shared/           # Authorization, i18n, CQRS decorators
  infra/            # Database, Redis, Meilisearch, HTTP client
  pkg/              # Utilities (apperror, httputil, validator, cursor, etc.)
migrations/         # PostgreSQL migrations
docs/http/          # OpenAPI spec + Swagger UI
wiki/               # Postman collection + architecture diagram
```
