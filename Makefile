# Makefile
.SILENT:

.PHONY: run build test lint clean
# Variables

APP_NAME=SQL_Converter_api

# MAKE Commands

run:
	go run cmd/api/main.go

build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

lint:
	golangci-lint run --timeout 5m

test:
	go test -v ./...

up:
	docker-compose up --build -d

down:
	docker-compose down

down-and-clean:
	docker-compose down -v

# Debugging
debug:
	@echo "Current directory: $(shell pwd)"
	@docker-compose config