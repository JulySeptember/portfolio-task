# 📌 Serverless Task Management App

Next.js × Go × AWS × Terraform × MySQL を用いた  
フルスタック Serverless Task Management App です。

認証・API・Infrastructure を含めた  
実践的なモダン Web アプリ構成を採用しています。

---

<p align="center">
  <img src="./docs/architecture_and_erd_v1.png" width="900">
</p>

---

# ✨ Features

```text
- AWS Cognito Authentication
- API Gateway JWT Authorizer
- Serverless Go API (AWS Lambda)
- Task CRUD APIs
- Owner-isolated authorization
- Terraform Infrastructure as Code
- Swagger / OpenAPI
- Structured Logging
- Layered Architecture
- Private RDS MySQL
- CloudFront + S3 Frontend Hosting
```

---

# 🌍 Live Demo

| Service | URL |
| --- | --- |
| Frontend | https://xxxxx.cloudfront.net |
| Swagger UI | https://xxxxx.execute-api.ap-northeast-1.amazonaws.com/api/docs |

---

# 🧩 Tech Stack

| Layer | Technology |
| --- | --- |
| Frontend | Next.js + TypeScript |
| Backend | Go |
| Infrastructure | Terraform |
| Authentication | AWS Cognito |
| API | API Gateway HTTP API |
| Runtime | AWS Lambda |
| Database | MySQL (RDS) |
| Hosting | S3 + CloudFront |

---

# 🏗 System Architecture

```text
Client
  ↓
CloudFront
  ↓
S3 (Next.js Frontend)

Client
  ↓
API Gateway (JWT Authorizer)
  ↓
Lambda (Go)
  ↓
Handler
  ↓
Service
  ↓
Repository
  ↓
RDS MySQL
```

---

# 🔐 Authentication

認証は AWS Cognito Hosted UI を利用。

JWT 検証は API Gateway JWT Authorizer に委譲しています。

```text
Login
  ↓
JWT 発行
  ↓
Authorization: Bearer <token>
  ↓
API Gateway JWT validation
  ↓
validated claims → Lambda
```

Lambda 側では検証済み claims のみを利用します。

---

# 👤 User Bootstrap

ログイン後に bootstrap API を呼び出し、  
Cognito user と users table を同期します。

仕様:

```text
- 初回ログイン時 INSERT
- 既存ユーザー UPDATE
- Cognito sub を auth_user_id として利用
```

---

# 🧱 Backend Design

Backend は Layered Architecture を採用。

```text
Handler
  ↓
Service
  ↓
Repository
```

特徴:

```text
- Handler / Service / Repository separation
- Context timeout
- Structured logging
- Owner isolation
- Strict JSON decode
- API Gateway JWT delegation
- Private RDS architecture
```

---

# 🎨 Frontend Design

Frontend は feature-based architecture を採用。

```text
features/
├── auth
└── tasks
```

各 feature は:

```text
- api
- hooks
- components
- schemas
- queries
- utils
```

を内部管理。

特徴:

```text
- React Query cache management
- optimistic update
- Dialog / Full Page editor
- Hashids URL obfuscation
- Responsive UI
- shadcn/ui + Radix UI
```

---

# 🔒 Security

```text
- Cognito Authentication
- API Gateway JWT Authorizer
- Request timeout
- SQL timeout
- Panic recovery
- Strict JSON decode
- Unknown field reject
- Body size limit
- Owner isolation
- Private RDS
- IMDSv2 required
- Encrypted EBS
```

---

# 📡 API

## Public Endpoints

```http
GET /health
GET /api/docs
GET /api/spec/swagger.yml
```

---

## Protected Endpoints

```text
/api/v1/*
```

JWT authentication required.

---

## Main APIs

### User

```text
POST   /api/v1/auth/bootstrap
GET    /api/v1/users/me
DELETE /api/v1/users/me
```

### Tasks

```text
POST   /api/v1/tasks
GET    /api/v1/tasks
GET    /api/v1/tasks/{id}
PUT    /api/v1/tasks/{id}
PATCH  /api/v1/tasks/{id}/status
DELETE /api/v1/tasks/{id}
```

Detailed request/response schemas are available in Swagger/OpenAPI.

---

# 🗄 Database

## users

| column | type |
| --- | --- |
| id | bigint |
| auth_user_id | varchar |
| email | varchar |
| created_at | datetime |
| updated_at | datetime |

---

## tasks

| column | type |
| --- | --- |
| id | bigint |
| user_id | bigint |
| title | varchar |
| description | text |
| status | varchar |
| due_date | datetime |
| created_at | datetime |
| updated_at | datetime |

---

# 📄 Swagger / OpenAPI

local:

```text
http://localhost:8080/api/docs
```

OpenAPI spec:

```text
http://localhost:8080/api/spec/swagger.yml
```

---

# 📁 Directory Structure

```text
.
├── backend
│   ├── cmd
│   │   └── api
│   │
│   ├── internal
│   │   ├── apperr
│   │   ├── auth
│   │   ├── config
│   │   ├── container
│   │   ├── dto
│   │   ├── handlers
│   │   ├── httpx
│   │   ├── middleware
│   │   ├── models
│   │   ├── repository
│   │   ├── router
│   │   └── service
│   │
│   ├── migrations
│   ├── swagger
│   │
│   └── Makefile
│
├── frontend
│   └── src
│       ├── app
│       │
│       ├── components
│       ├── features
│       ├── lib
│       └── providers
│
├── infra
│   ├── bootstrap
│   └── main
│       └── modules
│
├── scripts
│
└── .github/workflows
```
---

# 🏗 Infrastructure

Terraform により以下を構築。

```text
- VPC
- Public / Private Subnets
- Security Groups
- RDS MySQL
- Lambda
- API Gateway HTTP API
- Cognito User Pool
- S3
- CloudFront
- Bastion EC2
```

---

# 🚀 Future Improvements

```text
- GitHub Actions CI/CD
- Automated Lambda migration
- Bastion removal
- Secrets Manager
- SSM Parameter Store
- Integration tests
- Refresh token rotation
- HttpOnly cookie authentication
```