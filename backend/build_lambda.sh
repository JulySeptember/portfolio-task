#!/usr/bin/env sh
set -euo pipefail
OUT=out
BIN=bootstrap
rm -rf "$OUT" lambda.zip
mkdir -p "$OUT"
echo "Building linux/amd64 binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$OUT/$BIN" ./cmd/api
cd "$OUT"
zip -r ../lambda.zip "$BIN"
cd ..
echo "lambda.zip created at $(pwd)/lambda.zip"
