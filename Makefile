COMPOSE := $(shell if docker compose version >/dev/null 2>&1; then echo "docker compose"; else echo "docker-compose"; fi)

.PHONY: build run test migrate frontend-install frontend-build frontend-test docker-build docker-up docker-down logs

build:
	cd backend && go build -o ../bin/server ./cmd/server

run:
	cd backend && go run ./cmd/server

test:
	cd backend && go test ./...

migrate:
	cd backend && go run ./cmd/migrate

frontend-install:
	npm --prefix frontend install

frontend-build:
	npm --prefix frontend run build

frontend-test:
	npm --prefix frontend test

docker-build:
	$(COMPOSE) build

docker-up:
	$(COMPOSE) up --build

docker-down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f backend frontend
