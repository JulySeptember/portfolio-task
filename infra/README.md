# Terraform 実行前の AWS 手動準備ガイド（コピペ用・一括）

## 要点
- 必須リソース: (A) tfstate 用 S3 バケット（Terraform backend 用）、(B) State ロック用 DynamoDB テーブル、(C) Lambda アーティファクト格納用 S3 バケット（tfstate バケットと同一でも可）、(D) Terraform 実行用 IAM 資格情報（CI/実行者）。
- 注意: Lambda ZIP は Lambda 作成時に S3 に存在している必要がある（アップロードのタイミングを運用で確実にすること）。
- 推奨運用: S3 バケットは Terraform 管理に含めるが、Lambda オブジェクトは CI が配置する「2 段階 apply」を推奨。

## Usage
- CI の動作: ビルド成果物を s3://{artifact-bucket}/{key} にアップロードする。
- デフォルト key: lambda/your-project-dev.zip
- モジュール連携: ルートの module.lambda に backend_bucket_name = module.backend_s3.backend_bucket_id を渡す。

## 1. tfstate 用 S3 バケット（必須）
- 目的: Terraform の状態ファイル（tfstate）を保存するバックエンド。
- 要件:
  - バージョニングを有効にする（必須）。
  - パブリックアクセスをすべてブロックする（必須）。
  - 暗号化を有効にする（SSE-S3 または KMS）。
  - バケットポリシーでアクセスを限定（特定の IAM ロール/アカウントのみ許可）。
- 運用: このバケットは modules/state_bucket で Terraform 管理するか、事前に手動で作成する。
- 命名例: my-project-tfstate-2026

## 2. DynamoDB テーブル（State ロック用、必須）
- 目的: terraform apply の同時実行による競合を防ぐ。
- 要件:
  - パーティションキー = LockID（型: String）
  - 課金モード = ON_DEMAND
- テーブル名例: terraform-lock-table
- 運用: Terraform 管理に含めるか事前作成して backend ブロックで指定する。

## 3. Lambda アーティファクト格納用 S3（必須）
- 目的: Lambda のデプロイ用 ZIP を置くバケット。
- 要件:
  - バージョニング推奨、パブリックアクセスブロック、暗号化。
  - CI に s3:PutObject 権限を付与。
- 運用パターン（選択）:
  - (A) 事前アップロード: バケット作成後、CI が ZIP をアップロード → その後 Terraform を実行。
  - (B) 2 段階 apply（推奨）: まず S3 バケットだけを Terraform で apply → CI がアップロード → 残りを apply。
  - (C) Terraform によるアップロード: CI がアーティファクトを Terraform 実行環境へ渡し、aws_s3_bucket_object で配置 → Lambda 作成。
- 推奨: 2 段階 apply（バケットは Terraform 管理、オブジェクトは CI が配置）。

## 4. IAM ユーザー / ロール（必須）
- CI 用と開発者用は分離する（別ロール/ユーザー）。
- 初期は広めでも可だが最終的に最小権限化すること。
- CI が S3 にアップロードするなら CI に s3:PutObject 権限を付与。
- ローカル実行例: aws configure

## 5. backend ブロック（providers.tf に追加する例）
- 注意: backend ブロックは tfstate バケットと DynamoDB が存在することを前提に有効化する（先に作成するか state_bucket モジュールを先に apply する）。

terraform {
  backend "s3" {
  bucket         = "portfolio-task-july-tfstate-dev"
  key            = "dev/terraform.tfstate"
  region         = "ap-northeast-1"
  dynamodb_table = "portfolio-task-july-tf-lock-dev"
  encrypt        = true
  }
}

## 6. 実行手順（順序どおり・コピペ可）
1) state_bucket モジュール（または手動）で tfstate 用 S3 バケットを作成する（versioning, encryption, block public access, bucket policy）。  
2) DynamoDB テーブル（LockID）を作成する（ON_DEMAND）。  
3) providers.tf に backend ブロックを追加する（上記 HCL を貼る）。  
4) `terraform init` を実行して backend を初期化する（tfstate がリモートに切り替わる）。  
5) modules/backend_s3 を apply して Lambda アーティファクト用バケットを作成する（または事前に手動で作成）。  
6) CI による Lambda ZIP のアップロードを実行（例: `aws s3 cp build.zip s3://{artifact-bucket}/lambda/your-project-dev.zip`）。  
7) `terraform plan` を実行して差分を確認（出力をレビュー）。  
8) `terraform apply` を実行（まずは dev 環境で検証）。

## 7. 最短チェックリスト（コピペ用）
- [ ] tfstate S3 バケット作成（versioning, encryption, block public access, bucket policy）  
- [ ] DynamoDB テーブル作成（Partition key = LockID, ON_DEMAND）  
- [ ] providers.tf に backend ブロック追加（bucket/key/region/dynamodb_table/encrypt）  
- [ ] IAM: CI に s3:PutObject 権限、Terraform 実行者に必要権限付与  
- [ ] modules/backend_s3 を apply（artifact バケット作成）  
- [ ] CI が Lambda ZIP を s3://{artifact-bucket}/lambda/your-project-dev.zip にアップロード  
- [ ] terraform init → terraform plan（レビュー） → terraform apply

## 8. 重要な注意点（短く）
- Secrets（DB パスワード等）は平文 tfvars に置かない（Secrets Manager / SSM を使用）。  
- RDS の `skip_final_snapshot` は開発向け設定なので本番では見直す。  
- S3/KMS のアクセス監査（CloudTrail, S3 アクセスログ）を有効にする。

