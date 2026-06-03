# 📌 Go × AWS Serverless Task Management App

Next.js × Go × AWS × Terraform × MySQL を用いた
フルスタック Serverless Task Management App です。

---

## 🌍 Live Demo

| Service | URL |
| --- | --- |
| Frontend | https://dgw03czfpoc25.cloudfront.net |
| Swagger UI | https://h5kvlgfwv1.execute-api.ap-northeast-1.amazonaws.com/api/docs |

---

# ✨ Features

- Serverless Go API (AWS Lambda)
- Task CRUD API
- User Authentication (Cognito)
- Owner-Isolated Authorization
- Public ID Based Resource Access
- Terraform Infrastructure as Code
- API Gateway JWT Authorizer
- CloudFront + S3 Frontend Hosting
- Private RDS MySQL
- CI/CD with GitHub Actions

---

<p align="center">
  <img src="./docs/architecture_and_erd_v1.png" width="900">
</p>

### Login (AWS Cognito Hosted UI)

<p align="center">
  <img src="./docs/rogin.jpg" width="90%" />
</p>

### Task Dialog

<p align="center">
  <img src="./docs/dialog.jpg" width="90%" />
</p>

### Task List

<p align="center">
  <img src="./docs/list.jpg" width="90%" />
</p>

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

# 🔐 Authentication

AWS Cognito Hosted UI を利用した JWT 認証を採用しています。

```text
Cognito Hosted UI
  ↓
JWT Token
  ↓
API Gateway JWT Authorizer
  ↓
Lambda
```

特徴:

- JWTベース認証
- API Gatewayによる認証分離
- Owner Isolation による認可制御

---

# 🧱 Backend Architecture

```text
Handler
  ↓
Service
  ↓
Repository
  ↓
RDS MySQL
```

特徴:

- レイヤードアーキテクチャ
- Context Timeout
- Structured Logging
- Owner Isolation
- Strict JSON Validation

---

# 🎨 Frontend

- Next.js (App Router)
- TypeScript
- React Query
- Zustand
- shadcn/ui

---

# 🔒 Security

- Cognito Authentication
- API Gateway JWT Authorizer
- Owner Isolation
- Private RDS (No Public Access)
- SQL Timeout
- Request Timeout
- Panic Recovery
- IMDSv2 Required
- Body Size Limitation

---

# ⚙ CI/CD

- GitHub Actions CI
- Go Test / Go Vet
- Lambda Build & Deploy
- Pull Request Validation

---

# 🏗 Infrastructure

Terraform により構築:

- VPC
- Public / Private Subnets
- Security Groups
- RDS MySQL
- Lambda
- API Gateway HTTP API
- Cognito User Pool
- S3 + CloudFront
- Bastion EC2

---

# 💡 Design Highlights

- AWSフルマネージド構成によるサーバーレス設計
- 認証と認可の分離（Cognito + API Gateway）
- Public ID によるセキュリティ強化
- Lambda + RDS のシンプルな構成設計
- Terraform による完全な IaC 管理

---

# 🚧 Engineering Background

本プロジェクトは、情報システム部門での実務経験を基盤に、
レガシー環境からクラウドネイティブ開発への技術転換を目的として構築しています。

- 情報システム部門での業務経験（約1.5年）
  - VBAによる業務改善ツール開発
  - AS400の運用・保守
- 基本情報技術者 / 応用情報技術者 取得
- レガシー環境（業務システム）からクラウド（AWS + Go）への移行経験を意識した学習・実装

---
