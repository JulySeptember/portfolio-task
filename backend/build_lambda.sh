#!/usr/bin/env sh

set -euo pipefail

OUT=out
BIN=bootstrap

rm -rf "$OUT" lambda.zip

mkdir -p "$OUT"

echo "Building linux/arm64 binary..."

CGO_ENABLED=0 \
GOOS=linux \
GOARCH=arm64 \
go build \
  -ldflags="-s -w" \
  -o "$OUT/$BIN" \
  ./cmd/api

cd "$OUT"

zip -q -r ../lambda.zip "$BIN"

cd ..

echo "lambda.zip created at $(pwd)/lambda.zip"