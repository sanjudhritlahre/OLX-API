.PHONY: build run migrate-build migrate-up migrate-down

build:
	@go build -o bin/api ./cmd/api

run: build
	@./bin/api

migrate-build:
	@go build -o bin/migrate ./cmd/migrate

migrate-up: migrate-build
	@./bin/migrate up

migrate-down: migrate-build
	@./bin/migrate down