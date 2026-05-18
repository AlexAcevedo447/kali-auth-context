# Makefile para Kali Auth Context

# Variables
APP_NAME=kali-auth-context
DOCKER_COMPOSE_DEV=docker-compose.dev.yml
DOCKER_COMPOSE_PROD=docker-compose.prod.yml

# ----------- Docker (Desarrollo) -----------

.PHONY: docker-up-dev
docker-up-dev:
	docker compose -f $(DOCKER_COMPOSE_DEV) up -d --build

.PHONY: docker-down-dev
docker-down-dev:
	docker compose -f $(DOCKER_COMPOSE_DEV) down

.PHONY: docker-logs-dev
docker-logs-dev:
	docker compose -f $(DOCKER_COMPOSE_DEV) logs -f

.PHONY: docker-restart-dev
docker-restart-dev:
	docker compose -f $(DOCKER_COMPOSE_DEV) restart

# ----------- Docker (Producción) -----------

.PHONY: docker-up-prod
docker-up-prod:
	docker compose -f $(DOCKER_COMPOSE_PROD) up -d --build

.PHONY: docker-down-prod
docker-down-prod:
	docker compose -f $(DOCKER_COMPOSE_PROD) down

.PHONY: docker-logs-prod
docker-logs-prod:
	docker compose -f $(DOCKER_COMPOSE_PROD) logs -f

.PHONY: docker-restart-prod
docker-restart-prod:
	docker compose -f $(DOCKER_COMPOSE_PROD) restart

# ----------- Go (Local) -----------

.PHONY: run
default: run

run:
	go run ./cmd/api/main.go

.PHONY: build
build:
	go build -o $(APP_NAME) ./cmd/api/main.go

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run || true

# ----------- Utilidades -----------

.PHONY: clean
clean:
	rm -rf $(APP_NAME)

.PHONY: help
help:
	@echo "Comandos útiles para Kali Auth Context:"
	@echo "  make docker-up-dev      # Levanta entorno dev en Docker (con build)"
	@echo "  make docker-down-dev    # Detiene entorno dev en Docker"
	@echo "  make docker-logs-dev    # Logs de entorno dev en Docker"
	@echo "  make docker-up-prod     # Levanta entorno prod en Docker (con build)"
	@echo "  make docker-down-prod   # Detiene entorno prod en Docker"
	@echo "  make docker-logs-prod   # Logs de entorno prod en Docker"
	@echo "  make run                # Ejecuta la app localmente (Go)"
	@echo "  make build              # Compila la app localmente (Go)"
	@echo "  make test               # Ejecuta tests (Go)"
	@echo "  make tidy               # Limpia dependencias (Go)"
	@echo "  make lint               # Lint (requiere golangci-lint)"
	@echo "  make clean              # Elimina binario generado"
