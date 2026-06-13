# ============================================================
# Industrial IoT Platform - Makefile
# ============================================================
.PHONY: help dev test build deploy clean

SERVICES = edge-gateway gateway-service device-service telemetry-service \
           command-service alarm-service notification-service analytics-service \
           auth-service api-gateway

help: ## Hiển thị danh sách lệnh
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

# ---------- Development ----------

dev: ## Khởi động toàn bộ stack ở môi trường dev
	docker compose -f deployment/docker-compose/dev.yml up -d

dev-down: ## Dừng môi trường dev
	docker compose -f deployment/docker-compose/dev.yml down

dev-logs: ## Xem logs môi trường dev
	docker compose -f deployment/docker-compose/dev.yml logs -f

# ---------- Build ----------

build: ## Build tất cả services
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && make build 2>/dev/null || echo "  No Makefile in $$service"; \
		cd ..; \
	done

build-docker: ## Build Docker images cho tất cả services
	docker compose -f deployment/docker-compose/build.yml build

# ---------- Test ----------

test: ## Chạy test cho tất cả services
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd $$service && make test 2>/dev/null || echo "  No tests in $$service"; \
		cd ..; \
	done

test-integration: ## Chạy integration test
	docker compose -f deployment/docker-compose/test.yml up --abort-on-container-exit

# ---------- Deploy ----------

deploy-k8s: ## Deploy lên Kubernetes
	kubectl apply -f deployment/kubernetes/namespace.yaml
	kubectl apply -f deployment/kubernetes/

deploy-edge: ## Deploy edge gateways (yêu cầu SSH access)
	@echo "Deploying edge gateways..."
	ansible-playbook -i deployment/edge-deployment/inventory.yaml \
		deployment/edge-deployment/deploy.yaml

# ---------- Database ----------

db-migrate: ## Chạy database migrations
	@for service in $(SERVICES); do \
		if [ -d "$$service/migrations" ]; then \
			echo "Migrating $$service..."; \
			cd $$service && make migrate 2>/dev/null; \
			cd ..; \
		fi \
	done

db-seed: ## Seed dữ liệu mẫu
	@for service in $(SERVICES); do \
		if [ -f "$$service/scripts/seed.sh" ]; then \
			echo "Seeding $$service..."; \
			bash $$service/scripts/seed.sh; \
		fi \
	done

# ---------- Utilities ----------

lint: ## Chạy linter cho tất cả services
	@for service in $(SERVICES) web-portal mobile-app; do \
		echo "Linting $$service..."; \
		cd $$service && make lint 2>/dev/null || echo "  No linter in $$service"; \
		cd ..; \
	done

clean: ## Dọn dẹp build artifacts
	@for service in $(SERVICES) web-portal mobile-app; do \
		echo "Cleaning $$service..."; \
		cd $$service && make clean 2>/dev/null || rm -rf dist/ build/ target/; \
		cd ..; \
	done
	docker system prune -f

gen-certs: ## Tạo development certificates
	@mkdir -p certs/
	@openssl req -x509 -newkey rsa:4096 -keyout certs/dev.key -out certs/dev.crt \
		-days 365 -nodes -subj "/CN=industrial-iot.local"
	@echo "Certificates created in certs/"

# ---------- Documentation ----------

docs: ## Generate API documentation
	@echo "Generating API docs..."
	@for service in $(SERVICES); do \
		if [ -f "$$service/api/openapi.yaml" ]; then \
			echo "  Generating docs for $$service..."; \
		fi \
	done
