## 📌 タスク管理アプリ（ポートフォリオ）

このプロジェクトは Next.js × Go × AWS × Terraform × RDS(MySQL) を用いて構築した、
タスク管理 Web アプリです。

フロント・バックエンド・インフラをすべて自前で設計・実装し、
AWS 上で本番運用可能な構成を再現しています。

---

## アーキテクチャ図
<img src="./docs/architecture.png" width="500">

## ER 図
<img src="./docs/erd.png" width="400">


---

🟦 技術スタック（Tech Stack）

フロントエンド
- Next.js（App Router）
- TypeScript
- React
- Tailwind CSS

バックエンド
- Go（標準ライブラリ）
- AWS Lambda（Go Runtime）
- API Gateway（HTTP API）
- Swagger（OpenAPI 自動生成）

インフラ（AWS）
- VPC（RDS 用プライベートサブネット）
- RDS for MySQL（無料利用枠）
- Lambda（VPC 内接続）
- API Gateway
- Cognito（ユーザー認証）
- S3（フロントホスティング）
- CloudFront（CDN）
- IAM（権限管理）
- Terraform（IaC）

その他
- Git / GitHub
- GitHub Actions（CI/CD）
- Docker（Go Lambda ビルド用）

---

🟩 アプリ概要（What）

機能一覧
- ユーザー認証（Cognito）
- タスク CRUD（作成 / 取得 / 更新 / 削除）
- ステータス管理（TODO / DOING / DONE）
- レスポンシブ対応（スマホ / PC）

画面
- ログイン / サインアップ
- タスク一覧
- タスク作成
- ステータス変更

---

🟧 AWS アーキテクチャ（Architecture）

Next.js (S3 + CloudFront)
        ↓
API Gateway（JWT 検証）
        ↓
Lambda（Go）
        ↓
RDS for MySQL（VPC 内）

RDS を安全に利用するための構成
- VPC（/16）
- パブリックサブネット × 2（ALB / Lambda / IGW 接続用）
- プライベートサブネット × 2（RDS 用）
- インターネットゲートウェイ
- セキュリティグループ（Lambda → RDS のみ許可）

---

🟨 開発ワークフロー（How）

STEP 1：企画・要件定義
- アプリの目的を明確化
- 必要機能の洗い出し
- 画面一覧（Figma）

STEP 2：API 設計
- エンドポイント一覧
- リクエスト / レスポンス定義
- 認証方式（Cognito JWT）
- RDS MySQL の ER 図作成

STEP 3：AWS 構成設計
- API Gateway → Lambda → RDS
- Cognito 認証
- S3 + CloudFront
- Terraform による IaC

STEP 4：リポジトリ構成
```
portfolio/
  frontend/
  backend/
  infra/
  docs/
  scripts/
```
STEP 5：インフラ構築（Terraform）
- VPC
- サブネット
- セキュリティグループ
- RDS（MySQL）
- Lambda（VPC 内）
- API Gateway
- Cognito
- S3
- CloudFront

STEP 6：バックエンド（Go × Lambda）
- タスク CRUD
- MySQL 接続（RDS）
- API Gateway 連携
- Swagger（OpenAPI）自動生成

STEP 7：フロント（Next.js）
- 認証（Cognito Hosted UI / Amplify）
- タスク一覧 / 作成 / 更新
- API と接続

STEP 8：デプロイ
- フロント → S3 + CloudFront
- API → Lambda
- IaC → Terraform apply

STEP 9：README / ドキュメント整備
- 技術構成
- アーキテクチャ図
- ER 図
- API 一覧
- 工夫した点
- 今後の改善点

---

📚 ドキュメント（docs/）
- architecture.png（アーキテクチャ図）
- erd.png（ER 図）
- api-design.md（API 設計）
- infra-design.md（インフラ設計）

---

🚀 CI/CD（GitHub Actions）
- フロント：S3 + CloudFront に自動デプロイ
- バックエンド：Lambda に自動デプロイ
- Terraform：main ブランチで自動 apply

---

🧠 工夫した点
- RDS を VPC 内に配置し、Lambda からのみアクセス可能にした安全設計
- Cognito を使った JWT 認証
- Terraform による完全 IaC 化
- GitHub Actions による自動デプロイ
- Next.js App Router + Tailwind によるモダンな UI

---

🔧 今後の改善点
- RDS Proxy の導入（コネクション最適化）
- OpenAPI（Swagger）自動生成の CI 化
- E2E テスト（Playwright）
- ダークモード対応
- タスク期限の通知
