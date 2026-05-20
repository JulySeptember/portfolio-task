include .env

.PHONY: \
	run \
	db-wait \
	migrate-up \
	migrate-down \
	migrate-force \
	db-reset \
	bootstrap-dev

DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&multiStatements=true

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
# migrations
# =========================

migrate-up: db-wait
	migrate \
		-path migrations \
		-database "$(DB_URL)" \
		up

migrate-down: db-wait
	migrate \
		-path migrations \
		-database "$(DB_URL)" \
		down 1

migrate-force: db-wait
	migrate \
		-path migrations \
		-database "$(DB_URL)" \
		force 1

# =========================
# reset
# =========================

db-reset:
	docker compose down -v
	docker compose up -d
	$(MAKE) migrate-up

# =========================
# local dev bootstrap
# =========================

bootstrap-dev:
	curl -X POST \
		http://localhost:8080/api/v1/auth/bootstrap