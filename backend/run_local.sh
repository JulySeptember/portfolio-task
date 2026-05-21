#!/usr/bin/env zsh

setopt PIPE_FAIL ERR_EXIT NO_UNSET

# =========================
# cleanup
# =========================

cleanup() {
  echo ""
  echo "Stopping Docker containers..."
  docker compose stop
}

trap cleanup EXIT INT TERM

# =========================
# load .env
# =========================

if [ -f .env ]; then
  setopt allexport
  source .env
  unsetopt allexport
fi

# =========================
# runtime defaults
# =========================

export RUN_MODE=${RUN_MODE:-local}
export APP_ENV=${APP_ENV:-development}

# local dev auth bypass
export ENABLE_DEV_AUTH_BYPASS=${ENABLE_DEV_AUTH_BYPASS:-true}

# server
export PORT=${PORT:-8080}

# =========================
# docker
# =========================

echo "Starting Docker MySQL..."
docker compose up -d

echo "Waiting for MySQL..."
sleep 5

# =========================
# DB_DSN validation
# =========================

if [ -z "${DB_DSN:-}" ]; then
  echo ""
  echo "ERROR: DB_DSN is not set"
  echo ""
  exit 1
fi

# =========================
# debug info
# =========================

echo ""
echo "=============================="
echo "RUN_MODE                 = ${RUN_MODE}"
echo "APP_ENV                  = ${APP_ENV}"
echo "ENABLE_DEV_AUTH_BYPASS   = ${ENABLE_DEV_AUTH_BYPASS}"
echo "PORT                     = ${PORT}"
echo "=============================="
echo ""

# =========================
# start api
# =========================

echo "Starting API server..."
echo "Swagger Docs: http://localhost:${PORT}/api/docs/"
echo "Swagger YAML: http://localhost:${PORT}/api/spec/swagger.yml"
echo ""

go run ./cmd/api