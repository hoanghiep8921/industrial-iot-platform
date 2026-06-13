# Deployment

> Infrastructure as Code — Kubernetes, Terraform, Docker Compose & Edge Deployment

## 📁 Cấu Trúc

```
deployment/
├── kubernetes/              # Kubernetes manifests & Helm charts
│   ├── base/                #   Base configurations
│   │   ├── namespace.yaml
│   │   └── rbac.yaml
│   ├── infrastructure/      #   Infrastructure services
│   │   ├── emqx/            #   EMQX MQTT Broker
│   │   ├── kafka/           #   Apache Kafka (Strimzi)
│   │   ├── timescaledb/     #   TimescaleDB
│   │   ├── postgres/        #   PostgreSQL
│   │   └── redis/           #   Redis
│   ├── services/            #   Application services
│   │   ├── gateway/
│   │   ├── telemetry/
│   │   ├── device/
│   │   ├── command/
│   │   ├── alarm/
│   │   └── analytics/
│   ├── observability/       #   Monitoring stack
│   │   ├── prometheus.yml
│   │   ├── grafana/
│   │   └── loki/
│   └── ingress/             #   Ingress & API Gateway
│       └── kong/
├── terraform/               # Cloud infrastructure
│   ├── modules/
│   │   ├── vpc/             #   VPC & networking
│   │   ├── eks/             #   EKS / GKE cluster
│   │   ├── rds/             #   Managed databases
│   │   └── mqtt-lb/         #   Load balancer cho MQTT
│   ├── environments/
│   │   ├── dev/
│   │   ├── staging/
│   │   └── production/
│   └── variables.tf
├── docker-compose/          # Docker Compose files
│   ├── dev.yml              #   Development environment
│   ├── test.yml             #   Integration test environment
│   └── build.yml            #   Build images
└── edge-deployment/         # Edge Gateway deployment
    ├── ansible/             #   Ansible playbooks
    │   ├── deploy.yaml
    │   ├── update.yaml
    │   └── rollback.yaml
    ├── inventory/           #   Inventory files
    │   ├── hanoi.yaml
    │   ├── hcm.yaml
    │   └── danang.yaml
    └── edge-configs/        #   Edge device configurations
        ├── hanoi/
        ├── hcm/
        └── danang/
```

## 🚀 Triển Khai

### Local Development
```bash
docker compose -f deployment/docker-compose/dev.yml up -d
```

### Kubernetes Production
```bash
# Apply infrastructure first
kubectl apply -f deployment/kubernetes/infrastructure/

# Then application services
kubectl apply -f deployment/kubernetes/services/

# Apply ingress
kubectl apply -f deployment/kubernetes/ingress/
```

### Terraform (Cloud Provisioning)
```bash
cd deployment/terraform/environments/production
terraform init
terraform plan
terraform apply
```

### Edge Gateway Deployment
```bash
cd deployment/edge-deployment

# Deploy to Hanoi factory edge gateways
ansible-playbook -i inventory/hanoi.yaml ansible/deploy.yaml

# Deploy to HCM factory
ansible-playbook -i inventory/hcm.yaml ansible/deploy.yaml

# Deploy to Da Nang factory
ansible-playbook -i inventory/danang.yaml ansible/deploy.yaml
```

## 📊 Monitoring Stack

```
Prometheus (metrics)──┐
                       ├──► Grafana (dashboards + alerting)
Loki (logs)────────────┘

Key metrics:
- MQTT: connections, messages/sec, bytes/sec
- Kafka: consumer lag, throughput
- Services: request rate, latency, error rate
- Edge: uptime, bandwidth, cache size
- Database: query performance, connection pool
```
