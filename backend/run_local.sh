#!/usr/bin/env zsh
setopt PIPE_FAIL ERR_EXIT NO_UNSET

export RUN_MODE=local

if [ -f .env ]; then
  setopt allexport
  source .env
  unsetopt allexport
fi

echo "Starting MySQL (WSL)..."
sudo service mysql start

sleep 2

export DB_DSN="${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local"

echo "Starting API server..."
exec go run ./cmd/api
