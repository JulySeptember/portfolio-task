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
terraform apply
```

---

# 2. Lambda Artifact Upload

Lambda 用 ZIP を S3 にアップロードします。

```bash
aws s3 cp build.zip s3://<artifact-bucket>/lambda/app.zip
```

例:

```bash
aws s3 cp build.zip \
s3://portfolio-task-july-dev-backend-artifacts/lambda/app.zip
```

---

# 3. main

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

# 注意

## tfvars に Secrets を直接書かない

本番環境では以下を利用してください。

- AWS Secrets Manager
- SSM Parameter Store

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