.PHONY: build run test migrate frontend-install frontend-build frontend-test docker-build docker-up docker-down

build:
	go build -o server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

migrate:
	go run ./cmd/migrate

frontend-install:
	npm --prefix frontend install

frontend-build:
	npm --prefix frontend run build

frontend-test:
	npm --prefix frontend test

docker-build:
	docker compose build

docker-up:
	docker compose up --build

docker-down:
	docker compose down
