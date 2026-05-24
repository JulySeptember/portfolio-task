.PHONY: \
	run \
	db-wait \
	migrate-up-local \
	migrate-up-prod \
	migrate-down-local \
	migrate-down-prod \
	migrate-force-local \
	migrate-force-prod \
	db-reset \
	bootstrap-dev

# =========================
# app
# =========================

run:
	./run_local.sh

# =========================
# db
# =========================

db-wait:
	@echo "waiting for mysql..."
	@until docker exec taskdb_mysql mysqladmin ping -h"127.0.0.1" --silent; do \
		sleep 1; \
	done
	@echo "mysql ready"

# =========================
# migrations local
# =========================

migrate-up-local: db-wait
	@set -a && . ./.env.local && set +a && \
	migrate \
		-path migrations \
		-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true&multiStatements=true" \
		up

migrate-down-local:
	@set -a && . ./.env.local && set +a && \
	migrate \
		-path migrations \
		-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true&multiStatements=true" \
		down 1

migrate-force-local:
	@set -a && . ./.env.local && set +a && \
	migrate \
		-path migrations \
		-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true&multiStatements=true" \
		force 1

# =========================
# migrations production
# =========================

migrate-up-prod:
	@set -a && . ./.env.production && set +a && \
	migrate \
		-path migrations \
		-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true&multiStatements=true" \
		up

migrate-down-prod:
	@set -a && . ./.env.production && set +a && \
	migrate \
		-path migrations \
		-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true&multiStatements=true" \
		down 1

migrate-force-prod:
	@set -a && . ./.env.production && set +a && \
	migrate \
		-path migrations \
		-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true&multiStatements=true" \
		force 1

# =========================
# reset
# =========================

db-reset:
	docker compose down -v
	docker compose up -d
	$(MAKE) migrate-up-local

# =========================
# local dev bootstrap
# =========================

bootstrap-dev:
	curl -X POST \
		http://localhost:8080/api/v1/auth/bootstrap