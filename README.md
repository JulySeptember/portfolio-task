# 📌 タスク管理アプリ（Portfolio）

Next.js × Go × AWS × Terraform × MySQL を用いたフルスタックWebアプリです。

<img src="./docs/architecture_and_erd_v1.png" width="700">

---

# 🌐 システム構成

- Frontend: Next.js + S3 + CloudFront
- Backend: Go + AWS Lambda
- API: API Gateway HTTP API
- Database: RDS MySQL
- Auth: AWS Cognito + JWT Authorizer
- IaC: Terraform

---

# 🏗 アーキテクチャ

```text
Client
  ↓
CloudFront
  ↓
S3 (Next.js Hosting)

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

# 🧠 設計方針

- レイヤードアーキテクチャ
- Handler / Service / Repository 分離
- AWS サーバレス前提設計
- JWT 検証を API Gateway 側へ分離
- Context timeout によるリクエスト制御
---

# 🔐 認証設計

本システムは AWS Cognito + API Gateway JWT Authorizer を利用します。

JWT 検証は API Gateway 側で実施し、
Lambda 側では検証済み claims を利用します。

## 認証フロー

```text
1. Cognito Login
2. JWT 発行
3. Authorization: Bearer <JWT>
4. API Gateway JWT Authorizer が JWT を検証
5. claims を Lambda に転送
6. Middleware が AuthUser を Context に格納
7. users テーブルへ自動同期
```

## ユーザー同期

ユーザー作成 API は持たず、
ログイン時に自動同期します。

- Cognito sub を auth_user_id として利用
- 初回ログイン時に INSERT
- 既存ユーザーは UPDATE

---

# 📝 API

## Create Task

```http
POST /api/v1/tasks
```

- タイトル
- 説明
- ステータス
- 期限日

## List Tasks

```http
GET /api/v1/tasks
```

- pagination
- sorting
- status filtering
- owner isolation

デフォルトソート:

```text
created_at DESC
```

## Get Task

```http
GET /api/v1/tasks/{id}
```

- task_id + user_id で取得
- 他ユーザーアクセス不可

## Update Task

```http
PUT /api/v1/tasks/{id}
```

- 完全更新
- 所有者チェックあり

## Update Status

```http
PATCH /api/v1/tasks/{id}/status
```

## Delete Task

```http
DELETE /api/v1/tasks/{id}
```

## Task Status

```text
TODO
DOING
DONE
```

---

# 🗄 データベース設計

## users

| column | type |
|---|---|
| id | bigint |
| auth_user_id | varchar |
| email | varchar |
| created_at | datetime |
| updated_at | datetime |

## tasks

| column | type |
|---|---|
| id | bigint |
| user_id | bigint |
| title | varchar |
| description | text |
| status | varchar |
| due_date | datetime |
| created_at | datetime |
| updated_at | datetime |

---

# 📊 インデックス設計

```sql
CREATE INDEX idx_tasks_user_id
    ON tasks(user_id);

CREATE INDEX idx_tasks_user_created_at
    ON tasks(user_id, created_at DESC, id DESC);

CREATE INDEX idx_tasks_user_status_created
    ON tasks(user_id, status, created_at DESC, id DESC);

CREATE INDEX idx_tasks_user_due_date
    ON tasks(user_id, due_date, id DESC);
```

対応用途:

- ユーザー単位取得
- ステータス絞り込み
- created_at ソート
- due_date ソート
- pagination 最適化

---

# ⏰ due_date / Timezone

API の日時は RFC3339 UTC を使用します。

例:

```json
{
  "due_date": "2026-05-20T00:00:00Z"
}
```

- Backend は UTC で保存
- Frontend 側でローカルタイムへ変換
- timezone 差異による日付ズレを防止

---

# 🔌 Middleware

```text
CORS
  ↓
Recovery
  ↓
RequestID
  ↓
Logging
  ↓
Auth
  ↓
Router
```

---

# 🔒 セキュリティ

- API Gateway JWT Authorizer
- Cognito 認証
- request / SQL timeout
- strict JSON decode
- unknown field reject
- body size limit
---

# ⚙️ ローカル開発

## 起動

```bash
make run
```

## Migration

```bash
make migrate-up
```

## Swagger

```text
http://localhost:8080/api/v1/docs/
```

## ローカル認証バイパス

ローカル開発時のみ利用可能です。

```env
RUN_MODE=local
ENABLE_DEV_AUTH_BYPASS=true
```

production では利用不可
---

# 📄 OpenAPI

```text
/swagger/swagger.yml
```

---

# 📁 ディレクトリ構成

```text
internal/
├── apperr/
├── auth/
├── config/
├── container/
├── dto/
├── handlers/
├── httpx/
├── middleware/
├── models/
├── repository/
├── router/
└── service/
```

---

# ⚠️ 制約

- RBAC 未実装
- 全 API 認証必須（/health 除く）
- WebSocket 未対応
---

# 📚 技術スタック

| Layer | Technology |
|---|---|
| Frontend | Next.js |
| Backend | Go |
| Infra | Terraform |
| API | API Gateway HTTP API |
| Auth | AWS Cognito |
| Runtime | AWS Lambda |
| DB | MySQL (RDS) |
| Hosting | CloudFront + S3 |