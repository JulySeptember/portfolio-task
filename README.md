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

- AWS Cognito Hosted UI Authentication (Implicit Flow)
- API Gateway JWT Authorizer
- Serverless Go API (AWS Lambda)
- Task CRUD APIs
- User Bootstrap API
- Owner-Isolated Authorization
- Public ID Based Resource Access
- Terraform Infrastructure as Code
- Swagger / OpenAPI
- Structured Logging
- Layered Architecture
- Private RDS MySQL
- CloudFront + S3 Frontend Hosting

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
| UI | Tailwind CSS + shadcn/ui |
| State Management | React Query + Zustand |
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
Browser
   │
   ▼
CloudFront
   │
   ▼
S3 (Frontend)

Browser
   │
   ▼
Cognito Hosted UI
   │
   ▼
Implicit Flow
   │
   ▼
ID Token / Access Token
   │
   ▼
API Gateway JWT Authorizer
   │
   ▼
Lambda (Go)
   │
   ▼
Handler
   │
   ▼
Service
   │
   ▼
Repository
   │
   ▼
RDS MySQL
```

---

# 🔐 Authentication

AWS Cognito Hosted UI を利用した Implicit Flow を採用しています。

```text
Login
  ↓
Hosted UI
  ↓
ID Token
Access Token
(URL Hash)
  ↓
Frontend Storage
(localStorage)
  ↓
POST /api/v1/auth/bootstrap
  ↓
Authorization: Bearer <token>
  ↓
API Gateway JWT Validation
  ↓
Lambda
```

特徴:

```text
- Cognito Hosted UI
- Implicit Flow
- JWT Authentication
- URL Hash Token Retrieval
- User Bootstrap Synchronization
- API Gateway JWT Validation
```

---

# 👤 User Bootstrap

ログイン後に bootstrap API を実行し、
Cognito User と users テーブルを同期します。

```text
POST /api/v1/auth/bootstrap
```

仕様:

```text
- First Login → INSERT
- Existing User → UPDATE
- Cognito sub → auth_user_id
- Cognito email → email
```

---

# 🧱 Backend Design

Layered Architecture を採用しています。

```text
Handler
  ↓
Service
  ↓
Repository
```

特徴:

```text
- Handler / Service / Repository Separation
- Context Timeout
- Structured Logging
- Strict JSON Decode
- Owner Isolation
- JWT Delegation to API Gateway
- Private RDS Architecture
```

---

# 🎨 Frontend Design

Feature-Based Architecture を採用しています。

```text
- React Query
- Zustand
- Responsive UI
- Cognito Hosted UI Authentication
- Implicit Flow
- User Bootstrap
- shadcn/ui
```

---

# 🔒 Security

```text
- Cognito Hosted UI
- Implicit Flow
- API Gateway JWT Authorizer
- Owner Isolation
- Request Timeout
- SQL Timeout
- Panic Recovery
- Unknown Field Rejection
- Strict JSON Decode
- Body Size Limit
- Private RDS
- Encrypted EBS
- IMDSv2 Required
```

---

# 📡 API

## Public

```http
GET /health
GET /api/docs
GET /api/spec/swagger.yml
```

## Protected

```http
POST   /api/v1/auth/bootstrap

GET    /api/v1/users/me
DELETE /api/v1/users/me

POST   /api/v1/tasks
GET    /api/v1/tasks
GET    /api/v1/tasks/{publicId}
PUT    /api/v1/tasks/{publicId}
PATCH  /api/v1/tasks/{publicId}/status
DELETE /api/v1/tasks/{publicId}
```

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

## tasks

| column | type |
| --- | --- |
| id | bigint |
| public_id | varchar |
| user_id | bigint |
| title | varchar |
| description | text |
| status | varchar |
| due_date | datetime |
| created_at | datetime |
| updated_at | datetime |

---

# 📄 Swagger / OpenAPI

Local:

```text
http://localhost:8080/api/docs
```

OpenAPI Spec:

```text
http://localhost:8080/api/spec/swagger.yml
```

---

# 📁 Directory Structure

```text
.
├── backend
│   ├── cmd
│   ├── internal
│   │   ├── auth
│   │   ├── handlers
│   │   ├── service
│   │   ├── repository
│   │   ├── middleware
│   │   └── router
│   ├── migrations
│   └── swagger
│
├── frontend
│   └── src
│       ├── app
│       ├── components
│       ├── features
│       │   ├── auth
│       │   └── tasks
│       ├── lib
│       └── providers
│
├── infra
│   ├── bootstrap
│   └── main
│
└── .github
```

---

# 🏗 Infrastructure

Terraform により以下を構築。

```text
- VPC
- Public Subnets
- Private Subnets
- Security Groups
- RDS MySQL
- Lambda
- API Gateway HTTP API
- Cognito User Pool
- Cognito Hosted UI
- S3
- CloudFront
- Bastion EC2
```

---

# 🚀 Future Improvements

```text
- Migration to Authorization Code Flow + PKCE
- Refresh Token Support
- GitHub Actions CI/CD
- Lambda Deployment Pipeline
- Bastion Removal
- Secrets Manager
- SSM Parameter Store
- Integration Tests
- Rate Limiting
- CloudWatch Dashboard
```