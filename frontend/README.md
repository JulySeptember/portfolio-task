# Frontend

Serverless Task Management App の
Next.js App Router frontend です。

---

# 技術スタック

| Layer | Technology |
| --- | --- |
| Framework | Next.js |
| Language | TypeScript |
| Styling | Tailwind CSS v4 |
| UI Components | shadcn/ui |
| API State Management | React Query |
| Global State Management | Zustand |
| Icons | Lucide React |

---

# 開発

依存 package install:

```bash
npm install
```

開発 server 起動:

```bash
npm run dev
```

ブラウザ:

```text
http://localhost:3000
```

---

# 環境変数

`.env.local` を作成。

例:

```env
NEXT_PUBLIC_API_URL=https://xxxxxxxx.execute-api.ap-northeast-1.amazonaws.com

NEXT_PUBLIC_COGNITO_DOMAIN=xxxxxxxx.auth.ap-northeast-1.amazoncognito.com

NEXT_PUBLIC_COGNITO_CLIENT_ID=xxxxxxxx

NEXT_PUBLIC_COGNITO_REDIRECT_URI=http://localhost:3000/auth/callback

NEXT_PUBLIC_COGNITO_LOGOUT_URI=http://localhost:3000/login
```

---

# Directory Structure

```text
src/
├── app
├── components
├── features
├── lib
├── providers
└── store
```

---

# App Structure

```text
app/
├── (auth)
├── (protected)
├── api
└── auth
```

| Directory | Role |
| --- | --- |
| (auth) | login page |
| (protected) | authenticated pages |
| api | route handlers |
| auth | Cognito callback handling |

---

# UI

本 project は以下を利用。

```text
- Tailwind CSS
- shadcn/ui
- Radix UI
```

導入済み base components:

```text
- button
- input
- card
- dialog
- dropdown-menu
- table
- textarea
- badge
- skeleton
```

---

# State Management

## React Query

API 通信管理に利用。

用途:

```text
- API fetch
- cache management
- loading state
- mutation
- refetch
```

---

## Zustand

frontend global state 管理。

用途:

```text
- auth state
```

---

# Authentication

認証は AWS Cognito Hosted UI を利用。

frontend flow:

```text
Login button
  ↓
Cognito Hosted UI
  ↓
redirect (/auth/callback)
  ↓
token exchange
  ↓
localStorage 保存
  ↓
Authorization: Bearer <id_token>
```

---

# JWT Validation

JWT validation は frontend ではなく
API Gateway JWT Authorizer に委譲。

```text
Client
  ↓
API Gateway JWT Authorizer
  ↓
validated claims
  ↓
Lambda
```

---

# User Bootstrap

ログイン後に bootstrap API を呼び出し、
Cognito user と users table を同期。

仕様:

```text
- 初回 login 時 INSERT
- 既存 user は UPDATE
- Cognito sub を auth_user_id として利用
```

---

# API

backend API:

```text
/api/v1/*
```

health check:

```text
/health
```

Swagger:

```text
/api/docs
```

---

# Features

## Tasks

対応予定:

```text
- create task
- update task
- update task status
- delete task
- pagination
- filtering
- sorting
```

---

# Frontend Architecture

feature based structure を採用。

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
- types
- utils
```

を内部に持つ。

---

# Deployment

frontend は以下構成で deploy:

```text
CloudFront
  ↓
S3
  ↓
Next.js static assets
```