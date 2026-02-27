# ============================
# Global variables
# ============================
FRONTEND_DIR=frontend
BACKEND_DIR=backend
INFRA_DIR=infra

AWS_REGION=ap-northeast-1
S3_BUCKET=your-frontend-bucket
CLOUDFRONT_ID=YOUR_CLOUDFRONT_DISTRIBUTION_ID

LAMBDA_FUNCTION_NAME=portfolio-task-backend


# ============================
# Frontend
# ============================

frontend-install:
    cd $(FRONTEND_DIR) && npm install

frontend-build:
    cd $(FRONTEND_DIR) && npm run build

frontend-deploy: frontend-build
    aws s3 sync $(FRONTEND_DIR)/out s3://$(S3_BUCKET) --delete
    aws cloudfront create-invalidation \
        --distribution-id $(CLOUDFRONT_ID) \
        --paths "/*"
    @echo "🚀 Frontend deployed!"


# ============================
# Backend (Go Lambda)
# ============================

backend-build:
    cd $(BACKEND_DIR) && \
    GOOS=linux GOARCH=amd64 go build -o main ./cmd/api && \
    zip lambda.zip main

backend-deploy: backend-build
    aws lambda update-function-code \
        --function-name $(LAMBDA_FUNCTION_NAME) \
        --zip-file fileb://$(BACKEND_DIR)/lambda.zip
    @echo "🚀 Backend Lambda deployed!"


# ============================
# Terraform (Infra)
# ============================

tf-init:
    cd $(INFRA_DIR) && terraform init

tf-plan:
    cd $(INFRA_DIR) && terraform plan

tf-apply:
    cd $(INFRA_DIR) && terraform apply -auto-approve

tf-destroy:
    cd $(INFRA_DIR) && terraform destroy -auto-approve


# ============================
# Utility
# ============================

clean:
    rm -f $(BACKEND_DIR)/main
    rm -f $(BACKEND_DIR)/lambda.zip
    @echo "🧹 Cleaned build artifacts"

all: frontend-deploy backend-deploy
