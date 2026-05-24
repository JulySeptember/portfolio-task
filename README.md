# 📌 Serverless Task Management App

Next.js × Go × AWS × Terraform × MySQL を用いた  
フルスタック タスク管理アプリです。

---

<img src="./docs/architecture_and_erd_v1.png" width="900">

---

# ✨ Features

- JWT Authentication (AWS Cognito)
- Serverless Go API on AWS Lambda
- API Gateway JWT Authorizer
- Owner-isolated Task APIs
- Terraform Infrastructure as Code
- Swagger / OpenAPI documentation
- Structured logging
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

# 🌐 System Architecture

- Frontend: Next.js + S3 + CloudFront
- Backend: Go + AWS Lambda
- API: API Gateway HTTP API
- Database: RDS MySQL
- Authentication: AWS Cognito
- Infrastructure: Terraform

---

# 🏗 Architecture

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

# 🧩 Tech Stack

| Layer | Technology |
| --- | --- |
| Frontend | Next.js |
| Backend | Go |
| Infrastructure | Terraform |
| API | API Gateway HTTP API |
| Authentication | AWS Cognito |
| Runtime | AWS Lambda |
| Database | MySQL (RDS) |
| Hosting | S3 + CloudFront |

---

# 🔐 Authentication

認証は AWS Cognito + API Gateway JWT Authorizer を利用します。

JWT の検証は API Gateway 側で実施し、  
Lambda 側では検証済み claims のみを利用します。

```text
Login
  ↓
JWT 発行
  ↓
Authorization: Bearer <id_token>
  ↓
API Gateway JWT validation
  ↓
claims → Lambda
```

---

# 👤 User Bootstrap

ログイン後に bootstrap API を呼び出し、  
Cognito user と users table を同期します。

仕様:

- 初回ログイン時に INSERT
- 既存ユーザーは UPDATE
- Cognito sub を auth_user_id として利用

---

# 🧱 Backend Design

- Layered Architecture
- Handler / Service / Repository separation
- Context timeout
- Owner isolation
- Structured logging
- JWT verification offloaded to API Gateway
- Private RDS architecture

---

# 🔒 Security

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

---

# 📡 API

## Public Routes

```http
GET /health

GET /api/docs
GET /api/spec/swagger.yml
```

---

## JWT Required

```text
/api/v1/*
```

---

# 👤 User APIs

## Bootstrap User

認証済みユーザーを users table に同期します。

```http
POST /api/v1/auth/bootstrap
```

---

## Get Current User

現在ログイン中ユーザー情報を取得します。

```http
GET /api/v1/users/me
```

---

## Delete Current User

自分自身のアカウントを削除します。

関連 task は cascade delete されます。

```http
DELETE /api/v1/users/me
```

---

# ✅ Task APIs

## Create Task

新しい task を作成します。

```http
POST /api/v1/tasks
```

example:

```json
{
  "title": "Buy milk",
  "description": "Go to supermarket",
  "status": "TODO",
  "due_date": "2026-05-30T00:00:00Z"
}
```

---

## List Tasks

自分の task 一覧を取得します。

対応:

- pagination
- sorting
- status filtering
- owner isolation

```http
GET /api/v1/tasks?limit=20&status=TODO&sort=created_at&order=DESC
```

---

## Get Task

指定 task を取得します。

```http
GET /api/v1/tasks/{id}
```

仕様:

- task_id + user_id で取得
- 他ユーザー task は取得不可

---

## Update Task

task 情報を更新します。

```http
PUT /api/v1/tasks/{id}
```

---

## Update Task Status

task の status のみ更新します。

```http
PATCH /api/v1/tasks/{id}/status
```

status:

```text
TODO
DOING
DONE
```

---

## Delete Task

task を削除します。

```http
DELETE /api/v1/tasks/{id}
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
│   └── Makefile
│
├── frontend
│   └── src/app
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

# 📚 Infrastructure

Terraform により以下を構築しています。

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
- Secrets Manager / SSM Parameter Store
- Integration tests
```