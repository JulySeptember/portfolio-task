#!/usr/bin/env zsh

setopt PIPE_FAIL ERR_EXIT NO_UNSET

# =========================
# base environment
# =========================

export RUN_MODE=local

# dev auth by default (local only)
export AUTH_MODE=${AUTH_MODE:-dev}

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
# defaults
# =========================

export PORT=${PORT:-8080}

# =========================
# start docker
# =========================

echo "Starting Docker MySQL..."
docker compose up -d

# wait mysql
echo "Waiting for MySQL..."
sleep 5

# =========================
# DB DSN build
# =========================

export DB_DSN="${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"

# =========================
# debug info
# =========================

echo ""
echo "=============================="
echo "RUN_MODE   = ${RUN_MODE}"
echo "AUTH_MODE  = ${AUTH_MODE}"
echo "PORT       = ${PORT}"
echo "=============================="
echo ""

# =========================
# start api
# =========================

echo "Starting API server..."
echo "Swagger Docs: http://localhost:${PORT}/api/v1/docs/"
echo "Swagger YAML: http://localhost:${PORT}/api/v1/spec/swagger.yml"
echo ""

go run ./cmd/api

# ローカル用
# #!/usr/bin/env zsh
# setopt PIPE_FAIL ERR_EXIT NO_UNSET

# export RUN_MODE=local

# if [ -f .env ]; then
#   setopt allexport
#   source .env
#   unsetopt allexport
# fi

# echo "Starting MySQL (WSL)..."
# sudo service mysql start

# sleep 2

# export DB_DSN="${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"

# echo "Starting API server..."
# exec go run ./cmd/api
