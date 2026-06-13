# API Gateway

> API Gateway / Backend-for-Frontend — Entry point thống nhất cho Web & Mobile

## 🎯 Chức Năng

- **Web BFF**: Backend-for-Frontend cho Web Portal — aggregate data từ nhiều services
- **Mobile BFF**: Backend-for-Frontend cho Mobile App — tối ưu payload cho mobile
- **Rate Limiting**: Bảo vệ backend services khỏi quá tải
- **Circuit Breaker**: Ngăn lỗi lan truyền giữa các services
- **Request/Response Transform**: Adapt response format cho từng client
- **API Versioning**: Quản lý API versions

## 📁 Cấu Trúc

```
api-gateway/
├── web-bff/                 # Backend-for-Frontend cho Web Portal
│   ├── dashboard-api/       #   Dashboard data aggregation
│   ├── device-api/          #   Device management & control
│   ├── alarm-api/           #   Alarm center API
│   ├── report-api/          #   Reports & analytics
│   └── admin-api/           #   Admin functionality
├── mobile-bff/              # Backend-for-Frontend cho Mobile App
│   ├── monitoring-api/      #   Mobile monitoring endpoints
│   ├── control-api/         #   Quick control endpoints
│   ├── alarm-api/           #   Mobile alarm endpoints
│   └── sync-api/            #   Offline data sync
└── rate-limiter/            # Rate limiting & resilience
    ├── rate-limit-config/   #   Rate limit rules
    ├── circuit-breaker/     #   Circuit breaker patterns
    └── request-validator/   #   Request validation & sanitization
```

## 🔧 Công Nghệ

- **Reverse Proxy**: Kong / Traefik / Envoy
- **BFF Logic**: Node.js (Express/Fastify) hoặc Go
- **API Documentation**: OpenAPI 3.0 / Swagger
- **Service Discovery**: Kubernetes DNS / Consul

## 🔄 Request Flow

```
Client (Browser/Mobile)
    │
    ▼
┌──────────────────┐
│  Load Balancer   │ (TLS termination)
└────────┬─────────┘
         │
    ┌────▼─────┐     ┌──────────┐
    │   Kong    │────▶│ Rate     │
    │  Gateway  │     │ Limiter  │
    └────┬─────┘     └──────────┘
         │
    ┌────▼─────┐
    │ Auth      │ ← Keycloak (OAuth2 token validation)
    │ Plugin    │
    └────┬─────┘
         │
    ┌────▼──────────┐
    │  Web BFF /    │
    │  Mobile BFF   │ ← Aggregates data from multiple services
    └────┬──────────┘
         │
    ┌────▼──────────┐
    │  Backend      │
    │  Services     │ → Telemetry, Device, Alarm, Command, Analytics
    └───────────────┘
```
