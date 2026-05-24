# Terraform デプロイ手順

## 全体フロー

```text
bootstrap apply
  ↓
Lambda build
  ↓
Lambda upload
  ↓
main apply
  ↓
migration
  ↓
health check
  ↓
JWT auth test
```

---

# 1. Bootstrap Infrastructure

Terraform backend 用リソースを作成します。

```bash
cd infra/bootstrap

terraform init

terraform apply \
  -var-file=envs/dev.tfvars
```

作成:

```text
- tfstate S3 bucket
- terraform lock DynamoDB
- Lambda artifact S3 bucket
```

---

# 2. Lambda Build

```bash
GOOS=linux GOARCH=arm64 go build \
  -o bootstrap \
  ../../backend/cmd/api
```

```bash
zip lambda.zip bootstrap
```

---

# 3. Lambda Upload

```bash
aws s3 cp lambda.zip \
  s3://portfolio-task-july-dev-backend-artifacts/lambda/portfolio-dev.zip
```

確認:

```bash
aws s3 ls \
  s3://portfolio-task-july-dev-backend-artifacts/lambda/
```

---

# 4. Main Infrastructure Apply

```bash
cd infra/main

terraform init

terraform apply \
  -var-file=envs/dev.tfvars
```

作成される主なリソース:

```text
- VPC
- Security Group
- RDS MySQL
- Lambda
- API Gateway
- Cognito
- S3
- CloudFront
- Bastion
```

---

# 5. Migration

Bastion に SSM 接続:

```bash
aws ssm start-session \
  --target <bastion-instance-id>
```

instance id:

```bash
terraform output bastion_instance_id
```

migration 実行:

```bash
make migrate-up-prod
```

`.env.production`

```env
DB_HOST=<rds-endpoint>
DB_PORT=13306
DB_NAME=taskdb
DB_USER=admin
DB_PASSWORD=xxxxxxxx
```

---

# 6. Health Check

```bash
curl https://<api-domain>/health
```

成功例:

```json
{
  "status": "ok"
}
```

---

# 7. Cognito Login Test

```text
signup
  ↓
login
  ↓
JWT token取得
  ↓
authenticated API call
```

例:

```bash
curl \
  -H "Authorization: Bearer <jwt-token>" \
  https://<api-domain>/api/v1/users/me
```

---

# Lambda Runtime

```text
runtime      = provided.al2023
handler      = bootstrap
architecture = arm64
```

---

# 注意

```text
- Lambda ZIP は main apply 前に S3 upload 必須
- migration 未実行だと API は正常動作しない
- /api/v1/* は JWT 必須
- production では Secrets Manager / SSM 推奨
```