.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'


CONTAINER_NAME := cms-api


.PHONY: status
status: ## Check if container is running
	@docker ps -f name=$(CONTAINER_NAME) --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

.PHONY: logs
logs: ## Show container logs
	docker logs -f $(CONTAINER_NAME)

.PHONY: shell
shell: ## Enter container shell
	docker exec -it $(CONTAINER_NAME) sh

.PHONY: restart
restart: ## Restart the container
	@echo "Restarting container..."
	docker restart $(CONTAINER_NAME)
	@echo "Container restarted!"

.PHONY: stop
stop: ## Stop the container
	@echo "Stopping container..."
	docker stop $(CONTAINER_NAME)
	@echo "Container stopped!"

.PHONY: test
test: ## Run tests inside container
	docker exec $(CONTAINER_NAME) go test -v -race ./...

.PHONY: lint
lint: ## Run linter inside container
	docker exec $(CONTAINER_NAME) golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code inside container
	docker exec $(CONTAINER_NAME) go fmt ./...
	docker exec $(CONTAINER_NAME) goimports -w .

.PHONY: mocks
mocks: ## Generate mocks inside container
	docker exec $(CONTAINER_NAME) go generate ./...

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	./scripts/migrate.sh up

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	./scripts/migrate.sh down ${N:-1}

.PHONY: up
up: ## Start all services (including tools: Jaeger, Adminer)
	docker compose --profile tools up -d

.PHONY: down
down: ## Stop all services
	docker compose --profile tools down

.PHONY: erd
erd: ## Generate ERD diagram to wiki/erd.mmd
	./scripts/generate-erd.sh

.PHONY: build
build: ## Build docker image
	docker build --target dev -t cms-api:dev .

.PHONY: clean
clean: ## Remove container and clean tmp
	docker rm -f $(CONTAINER_NAME) 2>/dev/null || true
	rm -rf tmp/
