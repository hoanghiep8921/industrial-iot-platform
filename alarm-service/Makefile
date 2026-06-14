# ============================================================
# Alarm Service - Makefile
# ============================================================
.PHONY: run build test migrate clean lint docker-build

run: ## Chạy alarm service local
	go run ./cmd/alarm/

build: ## Build binary
	go build -o bin/alarm ./cmd/alarm/

test: ## Chạy tests
	go test ./...

migrate: ## Chạy database migrations qua CLI psql
	psql -h localhost -p 5433 -U iiot -d industrial_iot -f migrations/001_init.sql

lint: ## Chạy linter
	golangci-lint run ./...

clean: ## Dọn dẹp build artifacts
	rm -rf bin/

docker-build: ## Build Docker image
	docker build -t industrial-iot/alarm-service .
