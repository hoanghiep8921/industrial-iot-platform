# ============================================================
# Gateway Service - Makefile
# ============================================================
.PHONY: run build test clean

run: ## Chạy gateway service (kết nối HiveMQ Cloud)
	go run ./cmd/gateway/

build: ## Build binary
	go build -o bin/gateway ./cmd/gateway/

test: ## Chạy tests
	go test ./...

lint: ## Chạy linter
	golangci-lint run ./...

clean: ## Dọn dẹp
	rm -rf bin/

docker-build: ## Build Docker image
	docker build -t industrial-iot/gateway-service .
