# 🚀 Local / Terraform / AWS Deploy Guide

Serverless Task Management App の  
ローカル開発 / AWS deploy 手順。

対象:

```text
- Frontend (Next.js)
- Backend (Go Lambda)
- Infrastructure (Terraform)
- RDS Migration
- Cognito Authentication
```

---

# 🖥 Local Development

## 1. Frontend 起動

```bash
cd frontend

npm install
npm run dev
```

Frontend:

```text
http://localhost:3000
```

---

## 2. Backend 起動

```bash
cd backend

make db-reset
make run
```

別 terminal:

```bash
make bootstrap-dev
```

実行内容:

```text
- DB 初期化
- backend 起動
- mock user 作成
```

---

## 3. Access

```text
http://localhost:3000
```

---

# ☁ AWS Deploy

# 1. Terraform Bootstrap

Terraform backend 用リソース作成。

```bash
make tf-bootstrap-init
make tf-bootstrap-plan
make tf-bootstrap-apply
```

作成されるもの:

```text
- tfstate S3 bucket
- terraform lock DynamoDB
- Lambda artifact S3 bucket
```

---

# 2. Backend Lambda Build

Lambda 用 Go binary build。

```bash
make backend-build
```

---

# 3. Lambda Package

Lambda upload 用 ZIP 作成。

```bash
make backend-package
```

---

# 4. Lambda Artifact Upload

Lambda artifact bucket に upload。

```bash
make backend-upload
```

または:

```bash
make backend-build && \
make backend-package && \
make backend-upload
```

---

# 5. Terraform Infrastructure Apply

AWS infrastructure 作成。

```bash
make tf-init
make tf-plan
make tf-apply
```

作成対象:

```text
- VPC
- Public / Private Subnets
- Security Groups
- Lambda
- API Gateway
- Cognito
- RDS MySQL
- S3
- CloudFront
- Bastion EC2
```

---

# 6. Terraform Output 確認

```bash
cd infra/main

terraform output
```

確認するもの:

```text
- frontend_url
- backend_api_url
- bastion_instance_id
- rds_endpoint
- cognito_client_id
```

---

# 7. Frontend Build & Deploy

Next.js frontend を build して
S3 + CloudFront に deploy。
ただし、local.envが存在している場合リネームするか削除してからbuild

```bash
make frontend-build
make frontend-deploy
```

実行内容:

```text
- Next.js build
- static export
- S3 upload
- CloudFront invalidation
```

---

# 8. aws.env 更新

terraform output の値を反映。

```bash
export AWS_REGION=ap-northeast-1

export API_URL=https://xxxxxxxx.execute-api.ap-northeast-1.amazonaws.com

export BASTION_INSTANCE_ID=i-xxxxxxxxxxxxx

export RDS_ENDPOINT=xxxxxxxx.ap-northeast-1.rds.amazonaws.com

export FRONTEND_URL=https://xxxxxxxx.cloudfront.net

export COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxxxx

export COGNITO_USER_POOL_ID=ap-northeast-1_xxxxxxxxx

export CLOUDFRONT_DISTRIBUTION_ID=EXXXXXXXXXXXXX

export ID_TOKEN=xxxxxxxxxxxxxxxx
```

読み込み:

```bash
source aws.env
```

terraform destroy/apply 後に値が変わる場合は更新。

---

# 9. Bastion 起動

```bash
aws ec2 start-instances \
  --instance-ids $BASTION_INSTANCE_ID \
  --region $AWS_REGION
```

---

# 10. SSM Port Forward

別 terminal を開く。

```bash
aws ssm start-session \
  --target $BASTION_INSTANCE_ID \
  --document-name AWS-StartPortForwardingSessionToRemoteHost \
  --parameters "{\"host\":[\"$RDS_ENDPOINT\"],\"portNumber\":[\"3306\"],\"localPortNumber\":[\"13306\"]}" \
  --region $AWS_REGION
```

この terminal は閉じない。

---

# 11. DB Connection Check

```bash
mysql \
  -h 127.0.0.1 \
  -P 13306 \
  -u admin \
  -p
```

---

# 12. Production ENV 確認

backend/.env.production:

```env
DB_HOST=127.0.0.1
DB_PORT=13306
DB_NAME=taskdb
DB_USER=admin
DB_PASSWORD=xxxxxxxx
```

---

# 13. Migration 実行

backend directory へ移動。

```bash
make migrate-up-prod
```

成功例:

```text
1/u create_users_table
2/u create_tasks_table
```

---

# 14. SSM Session 終了

SSM terminal を Ctrl+C。

---

# 15. Bastion 停止

```bash
aws ec2 stop-instances \
  --instance-ids $BASTION_INSTANCE_ID \
  --region $AWS_REGION
```

停止確認:

```bash
aws ec2 describe-instances \
  --instance-ids $BASTION_INSTANCE_ID \
  --region $AWS_REGION \
  --query 'Reservations[0].Instances[0].State.Name'
```

---

# 16. Lambda Logs

CloudWatch Logs 確認。

```bash
aws logs tail \
/aws/lambda/portfolio-task-july-api-dev \
--follow
```

---

# 17. Health Check

```bash
curl $API_URL/health
```

期待値:

```json
{"status":"ok"}
```

---

# 18. Swagger UI

```text
${API_URL}/api/docs
```

OpenAPI spec:

```text
${API_URL}/api/spec/swagger.yml
```

---

# 19. Frontend Access

```text
$FRONTEND_URL
```

確認項目:

```text
- signup
- signin
- JWT login
- task CRUD
- pagination
- filtering
- logout
```

---

# 20. API Functional Test

## User Bootstrap

```bash
curl -X POST $API_URL/api/v1/auth/bootstrap \
  -H "Authorization: Bearer $ID_TOKEN"
```

---

## Create Task

```bash
curl -X POST $API_URL/api/v1/tasks \
  -H "Authorization: Bearer $ID_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"first task",
    "description":"hello world",
    "status":"TODO",
    "due_date":"2026-06-01T12:00:00Z"
  }'
```

---

## Update Task

```bash
curl -X PUT $API_URL/api/v1/tasks/1 \
  -H "Authorization: Bearer $ID_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"updated title",
    "description":"updated description",
    "status":"DOING",
    "due_date":"2026-06-10T15:00:00Z"
  }'
```

---

## Update Task Status

```bash
curl -X PATCH $API_URL/api/v1/tasks/1/status \
  -H "Authorization: Bearer $ID_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"DONE"}'
```

---

## List Tasks

```bash
curl $API_URL/api/v1/tasks \
  -H "Authorization: Bearer $ID_TOKEN"
```

---

## Get Task

```bash
curl $API_URL/api/v1/tasks/1 \
  -H "Authorization: Bearer $ID_TOKEN"
```

---

## Delete Task

```bash
curl -X DELETE $API_URL/api/v1/tasks/1 \
  -H "Authorization: Bearer $ID_TOKEN"
```

---

## Get Current User

```bash
curl $API_URL/api/v1/users/me \
  -H "Authorization: Bearer $ID_TOKEN"
```

---

## Delete Current User

```bash
curl -X DELETE $API_URL/api/v1/users/me \
  -H "Authorization: Bearer $ID_TOKEN"
```

---

# 21. DB Connection Verification

CloudWatch Logs に:

```text
connected to db (mode=lambda ...)
```

が出れば成功。

---

# 22. Troubleshooting

最初に確認する場所:

```text
1. Lambda CloudWatch Logs
2. API Gateway Logs
3. Lambda Security Group
4. RDS Security Group
5. subnet route table
6. S3 artifact path
7. Cognito callback URL
8. CloudFront invalidation status
9. frontend environment variables
```