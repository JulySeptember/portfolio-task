#!/usr/bin/env zsh
setopt PIPE_FAIL ERR_EXIT NO_UNSET

export RUN_MODE=local
# cleanup
cleanup() {
echo ""
echo "Stopping Docker containers..."
docker compose stop
}
trap cleanup EXIT INT TERM
# load .env
if [ -f .env ]; then
setopt allexport
source .env
unsetopt allexport
fi
# start docker
echo "Starting Docker MySQL..."
docker compose up -d
# wait mysql
echo "Waiting for MySQL..."
sleep 5
# DB DSN
export DB_DSN="${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"
# start api
echo "Starting API server..."
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
