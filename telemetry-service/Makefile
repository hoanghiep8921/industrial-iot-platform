# ============================================================
# Telemetry Service - Makefile
# ============================================================
.PHONY: run build test migrate clean

run: ## Chạy telemetry service (Kafka → TimescaleDB)
	go run ./cmd/telemetry/

build: ## Build binary
	go build -o bin/telemetry ./cmd/telemetry/

test: ## Chạy tests
	go test ./...

migrate: ## Chạy database migrations
	psql -h localhost -p 5432 -U iiot -d telemetry -f migrations/001_init.sql

lint: ## Chạy linter
	golangci-lint run ./...

clean: ## Dọn dẹp
	rm -rf bin/

docker-build: ## Build Docker image
	docker build -t industrial-iot/telemetry-service .
