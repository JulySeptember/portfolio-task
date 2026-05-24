# Terraform Deployment Guide

## Deployment Flow

```text
bootstrap apply
  ↓
Lambda build
  ↓
Lambda package
  ↓
Lambda upload
  ↓
main apply
  ↓
SSM bastion connect
  ↓
migration
  ↓
health check
  ↓
JWT auth test
```

---

# 1. Bootstrap Infrastructure

Terraform backend resources を作成します。

```bash
cd infra/bootstrap

terraform init

terraform apply
```

作成されるリソース:

```text
- Terraform state S3 bucket
- Terraform lock DynamoDB table
- Lambda artifact S3 bucket
```

---

# 2. Lambda Build

```bash
cd backend

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
go build \
  -o bootstrap \
  ./cmd/api
```

---

# 3. Lambda Package

```bash
zip lambda.zip bootstrap
```

---

# 4. Lambda Upload

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

# 5. Main Infrastructure Apply

```bash
cd infra/main

terraform init

terraform apply \
  -var-file=envs/dev.tfvars
```

作成される主なリソース:

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

# 6. Bastion SSM Connection

Bastion instance ID を取得:

```bash
cd infra/main

terraform output bastion_instance_id
```

SSM 接続:

```bash
aws ssm start-session \
  --target <bastion-instance-id>
```

---

# 7. Database Migration

migration 実行前に `.env.production` を設定します。

```env
DB_HOST=<rds-endpoint>
DB_PORT=13306
DB_NAME=taskdb
DB_USER=admin
DB_PASSWORD=xxxxxxxx
```

補足:

```text
13306 is local forwarded port via SSM tunnel
```

migration 実行:

```bash
make migrate-up-prod
```

---

# 8. Health Check

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

# 9. Cognito Authentication Test

```text
signup
  ↓
login
  ↓
JWT token acquisition
  ↓
authenticated API call
```

example:

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

# Notes

```text
- Lambda ZIP upload is required before terraform apply
- Database migration is required before API usage
- /api/v1/* requires JWT authentication
- JWT validation is handled by API Gateway
- Lambda only trusts validated claims
- Production environments should use Secrets Manager or SSM Parameter Store
```