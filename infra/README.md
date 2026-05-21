# Terraform デプロイ手順

本プロジェクトの Terraform は以下の 2 段階で構成されています。

```text
infra/
├── bootstrap/
└── main/
```

---

# 1. bootstrap

Terraform backend 用リソースを作成します。

作成されるもの:

- tfstate 用 S3 バケット
- terraform lock 用 DynamoDB
- Lambda artifact 用 S3 バケット

## 実行

```bash
cd infra/bootstrap

terraform init
terraform apply -var-file=envs/dev.tfvars
```

---

# 2. Lambda Build

Go custom runtime (`provided.al2023`) 用の
Lambda binary を build します。

```bash
GOOS=linux GOARCH=arm64 go build \
-o bootstrap \
../../backend/cmd/api

zip lambda.zip bootstrap
```

---

# 3. Lambda Artifact Upload

Lambda 用 ZIP を S3 にアップロードします。

```bash
aws s3 cp lambda.zip \
s3://<artifact-bucket>/lambda/<project>-dev.zip
```

例:

```bash
aws s3 cp lambda.zip \
s3://portfolio-task-july-dev-backend-artifacts/lambda/portfolio-dev.zip
```

---

# 4. main

アプリケーション本体の AWS リソースを作成します。

作成されるもの:

- VPC
- Lambda
- API Gateway
- Cognito
- RDS MySQL
- CloudFront
- Frontend S3

## 実行

```bash
cd infra/main

terraform init

terraform apply \
-var-file=envs/dev.tfvars
```

---

# 🔐 API Gateway 認証設計

以下は public route です。

```text
GET /health
GET /docs/*
GET /docs/swagger.yml
```

以下は JWT 認証必須です。

```text
/api/v1/*
```

認証は Cognito JWT Authorizer により実施されます。

---

# 🔑 Cognito Domain

Cognito domain は AWS グローバルで一意である必要があります。

本プロジェクトでは AWS Account ID を suffix に付与しています。

例:

```text
portfolio-dev-123456789012-auth
```

---

# 👤 terraform_user_arn

`dev.tfvars` の IAM ARN はサンプル値です。

実際の AWS IAM User / Role ARN に変更してください。

例:

```tfvars
terraform_user_arn = "arn:aws:iam::123456789012:user/terraform-user"
```

---

# ⚠️ 注意

## tfvars に Secrets を直接書かない

本番環境では以下を利用してください。

- SSM Parameter Store
- Secrets Manager

---

## Lambda ZIP が必要

`terraform apply` 前に
Lambda ZIP が S3 に存在している必要があります。

---

## 開発環境設定

現在の設定は開発用です。

以下は本番で見直してください。

- skip_final_snapshot
- DB password
- IAM permissions
- CloudWatch retention
- Security Group egress rules

---

# 📦 Lambda Runtime

Lambda runtime:

```text
provided.al2023
```

handler:

```text
bootstrap
```

Go custom runtime を前提としています。