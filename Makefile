# ============================
# Global variables
# ============================
include aws.env

FRONTEND_DIR := frontend
BACKEND_DIR := backend

BOOTSTRAP_DIR := infra/bootstrap
INFRA_DIR := infra/main

TF_ENV_FILE := envs/dev.tfvars

AWS_REGION := ap-northeast-1

LAMBDA_ARTIFACT_BUCKET := portfolio-task-july-dev-backend-artifacts
LAMBDA_ARTIFACT_KEY := lambda/portfolio-dev.zip

LAMBDA_BINARY := bootstrap
LAMBDA_ZIP := lambda.zip

FRONTEND_BUCKET = portfolio-task-july-dev-frontend-assets
CLOUDFRONT_ID := $(CLOUDFRONT_DISTRIBUTION_ID)

# ============================
# Frontend
# ============================

frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

frontend-deploy: 
	aws s3 sync $(FRONTEND_DIR)/out s3://$(FRONTEND_BUCKET) --delete
	aws cloudfront create-invalidation \
		--distribution-id $(CLOUDFRONT_ID) \
		--paths "/*"
	@echo "🚀 Frontend deployed"


# ============================
# Backend (Go Lambda)
# ============================

backend-build:
	cd $(BACKEND_DIR) && \
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
	go build \
	-o $(LAMBDA_BINARY) \
	./cmd/api

backend-package: 
	cd $(BACKEND_DIR) && \
	rm -f $(LAMBDA_ZIP) && \
	zip $(LAMBDA_ZIP) $(LAMBDA_BINARY)

backend-upload:
	aws s3 cp \
		$(BACKEND_DIR)/$(LAMBDA_ZIP) \
		s3://$(LAMBDA_ARTIFACT_BUCKET)/$(LAMBDA_ARTIFACT_KEY)

backend-migrate-up:
	@set -a && . ./backend/.env.production && set +a && \
	migrate \
		-path backend/migrations \
		-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true&multiStatements=true" \
		up

# ============================
# Terraform Bootstrap
# ============================

tf-bootstrap-init:
	cd $(BOOTSTRAP_DIR) && terraform init

tf-bootstrap-fmt:
	cd $(BOOTSTRAP_DIR) && terraform fmt -recursive

tf-bootstrap-validate:
	cd $(BOOTSTRAP_DIR) && terraform validate

tf-bootstrap-plan:
	cd $(BOOTSTRAP_DIR) && terraform plan

tf-bootstrap-apply:
	cd $(BOOTSTRAP_DIR) && terraform apply -auto-approve

tf-bootstrap-destroy:
	cd $(BOOTSTRAP_DIR) && terraform destroy -auto-approve


# ============================
# Terraform Main Infra
# ============================

tf-init:
	cd $(INFRA_DIR) && terraform init

tf-fmt:
	cd $(INFRA_DIR) && terraform fmt -recursive

tf-validate:
	cd $(INFRA_DIR) && terraform validate

tf-plan:
	cd $(INFRA_DIR) && terraform plan \
		-var-file=$(TF_ENV_FILE)

tf-apply: 
	cd $(INFRA_DIR) && terraform apply \
		-auto-approve \
		-var-file=$(TF_ENV_FILE)

tf-destroy:
	cd $(INFRA_DIR) && terraform destroy \
		-auto-approve \
		-var-file=$(TF_ENV_FILE)


# ============================
# Utility
# ============================

clean:
	rm -f $(BACKEND_DIR)/$(LAMBDA_BINARY)
	rm -f $(BACKEND_DIR)/$(LAMBDA_ZIP)
	@echo "🧹 Cleaned build artifacts"