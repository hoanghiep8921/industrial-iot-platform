# ============================================================
# Device Service - Makefile
# ============================================================
.PHONY: run build test lint clean docker-build

run: ## Chạy device service ở local
	go run ./cmd/device/

build: ## Build binary
	go build -o bin/device ./cmd/device/

test: ## Chạy tests
	go test ./...

test-cover: ## Chạy tests với coverage report
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint: ## Chạy linter
	golangci-lint run ./...

vet: ## Chạy go vet
	go vet ./...

clean: ## Dọn dẹp
	rm -rf bin/ coverage.out

docker-build: ## Build Docker image
	docker build -t industrial-iot/device-service .
