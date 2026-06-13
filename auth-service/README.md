# Auth Service

> Xác thực & Phân quyền — Identity Provider, RBAC & Audit Log

## 🎯 Chức Năng

- **Identity Provider**: Tích hợp Keycloak — OAuth2/OIDC, SSO
- **RBAC**: Role-based access control với các role: Admin, Supervisor, Operator, Viewer
- **Audit Log**: Ghi lại mọi hoạt động của người dùng
- **API Key Management**: API key cho service-to-service communication
- **Device Identity**: Quản lý danh tính thiết bị (X.509 certs)

## 📁 Cấu Trúc

```
auth-service/
├── identity-provider/       # Keycloak integration
│   ├── keycloak-config/     #   Realm, client, role configuration
│   ├── user-sync/           #   Đồng bộ user từ LDAP/AD
│   └── sso-handler/         #   Single Sign-On
├── rbac/                    # Role-based access control
│   ├── role-definition/     #   Định nghĩa roles & permissions
│   ├── policy-engine/       #   Policy evaluation (OPA/Rego)
│   └── tenant-isolation/    #   Phân tách dữ liệu theo tenant/nhà máy
└── audit-log/               # User activity audit
    ├── event-collector/     #   Thu thập sự kiện
    ├── audit-store/         #   Lưu trữ audit log
    └── compliance-report/   #   Báo cáo cho compliance
```

## 🔧 Công Nghệ

- **Identity**: Keycloak 24+ (OAuth2 / OIDC)
- **Policy Engine**: Open Policy Agent (OPA)
- **Database**: PostgreSQL
- **Protocol**: OAuth2 Authorization Code Flow + PKCE

## 👥 Role Definitions

```yaml
roles:
  admin:
    permissions:
      - users:manage
      - devices:manage
      - config:manage
      - audit:view
    scope: global

  supervisor:
    permissions:
      - devices:view
      - devices:control
      - alarms:acknowledge
      - alarms:escalate
      - reports:view
      - dashboards:view
    scope: factory  # Chỉ trong nhà máy được phân công

  operator:
    permissions:
      - devices:view
      - devices:control_basic
      - alarms:acknowledge
      - dashboards:view
    scope: area     # Chỉ trong khu vực được phân công

  viewer:
    permissions:
      - devices:view
      - dashboards:view
    scope: factory
```
