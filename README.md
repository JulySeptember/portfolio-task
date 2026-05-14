# 📌 タスク管理アプリ（Portfolio）

Next.js × Go × AWS × Terraform × MySQL を用いたフルスタックWebアプリです。

単なるCRUDではなく、**実務レベルの設計思想（認証・マルチテナント・レイヤード構造・IaC）**を意識して構築しています。

---

## アーキテクチャ図 & ER 図
<img src="./docs/architecture_and_erd_v1.png" width="500">

---

# 🌐 システム構成

- Frontend: Next.js（S3 + CloudFront）
- Backend: Go（AWS Lambda）
- API Gateway（HTTP API）
- Database: RDS MySQL（Private Subnet）
- Auth: AWS Cognito（JWT / RS256）

---

# 🏗 アーキテクチャ

Client
  ↓
CloudFront
  ↓
API Gateway
  ↓
Lambda（Go）
  ↓
Service → Repository → MySQL

---

# 🧠 アーキテクチャ設計思想

- レイヤードアーキテクチャ採用
- Handler / Service / Repository 分離
- DI（Interface）による依存逆転
- AWSサーバレス前提設計
- テスト容易性を重視した構造

---

# 🔐 認証設計

- AWS Cognitoを利用
- JWT（ID Token）による認証
- RS256 + JWKS検証
- Middlewareで認証処理を一元化

### 認証フロー

1. Cognitoでログイン
2. ID Token取得
3. APIリクエストにBearer付与
4. LambdaでJWT検証
5. usersテーブルへ自動同期（Ensure）

---

# 👤 ユーザー管理

本システムは**ユーザー作成APIを持たない設計**です。

ログイン時に自動でユーザーが作成されます。

### Ensure方式

- Cognito sub（auth_user_id）をキーに管理
- 初回ログイン時にINSERT
- 既存ユーザーはUPSERTで更新

👉 ユーザー登録 = ログイン処理

---

# 📝 タスク管理機能

本アプリの中核機能であり、
ユーザーごとにタスクを完全分離して管理します（マルチテナント設計）。

---

## 機能一覧

### ■ タスク作成（Create）
- タイトル / 説明 / ステータス / 期限を登録
- user_idは認証コンテキストから自動付与

---

### ■ タスク一覧取得（List）
- 自分のタスクのみ取得
- ページネーション対応（limit / offset）
- created_at降順

---

### ■ タスク詳細取得（Get）
- task_id + user_idで取得
- 他ユーザーのデータはアクセス不可

---

### ■ タスク更新（Update）
- タイトル / 説明 / ステータス / 期限更新
- 所有者チェックあり（user_id制御）

---

### ■ タスク削除（Delete）
- 物理削除
- user_idによる厳密制御

---

### ■ ステータス管理
- TODO / DOING / DONE
- シンプルな3状態管理で進捗を可視化

---

## 技術的ポイント

- RepositoryパターンでDB層を分離
- Service層でビジネスロジック制御
- SQLレベルで user_id 制約を強制
- マルチテナント設計

---

# 🗄 データベース設計

## users

- id (PK)
- auth_user_id (Cognito sub)
- email
- created_at
- updated_at

## tasks

- id (PK)
- user_id (FK)
- title
- description
- status
- due_date
- created_at
- updated_at

---

# 📊 インデックス設計

CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_user_status ON tasks(user_id, status);
CREATE INDEX idx_tasks_user_created_at ON tasks(user_id, created_at DESC, id DESC);

---

# 🔌 Middleware構成

- CORS
- Recovery（panic制御）
- Logging（リクエスト追跡）
- Auth（JWT検証 + ユーザー同期）

実行順：

CORS → Recovery → Logging → Auth → Router

---

# 📡 API一覧

## Users

- GET /api/v1/users/me
- DELETE /api/v1/users/me

## Tasks

- GET /api/v1/tasks
- POST /api/v1/tasks
- GET /api/v1/tasks/{id}
- PUT /api/v1/tasks/{id}
- DELETE /api/v1/tasks/{id}

---

# 📄 API仕様（Swagger / OpenAPI）

本プロジェクトではAPI仕様管理に **Swagger（OpenAPI 3.0）** を採用しています。

## 🌐 Swagger UI

http://localhost:8080/api/v1/docs/


# 🔒 セキュリティ設計

- RDSはPrivate Subnet配置
- Lambdaのみアクセス許可
- IAM最小権限
- JWT検証（RS256 + JWKS）
- request size制限（1MB）
- unknown field拒否
- SQL timeout / context timeout
- graceful shutdown対応

---

# ⚙️ ローカル開発

make run
make migrate-up

---

# 🚀 CI/CD

## Frontend
- S3 deploy
- CloudFront invalidation

## Backend
- Lambda deploy

## Infrastructure
- Terraform apply

---

# 🧠 このプロジェクトの特徴

- AWSサーバレス構成
- 認証付きマルチテナント設計
- DI + レイヤードアーキテクチャ
- 実務想定のセキュリティ設計
- TerraformによるIaC

---

# ⚠️ 現状の制約

- Access Tokenベース認可は未実装（ID Tokenのみ）
- ユーザー登録APIなし（ログイン時自動作成）
- 全APIは認証必須（/health除く）