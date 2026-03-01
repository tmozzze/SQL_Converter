# Makefile
.SILENT:

.PHONY: run build test lint clean
# Variables

APP_NAME=SQL_Converter_api

# MAKE Commands

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

# Swagger docs gen
swagger-gen:
	swag init -g cmd/api/main.go -o docs

# Debugging
debug:
	@echo "Current directory: $(shell pwd)"
	@docker-compose config