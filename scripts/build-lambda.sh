#!/bin/bash

# Lambda 用ビルドスクリプト
# Linux x86_64 向けにクロスコンパイルし、ZIP パッケージを作成する

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="$PROJECT_ROOT/dist"

echo "Building Sugu Lambda..."

# 出力ディレクトリを作成
mkdir -p "$OUTPUT_DIR"

# Linux x86_64 向けにビルド
cd "$PROJECT_ROOT"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$OUTPUT_DIR/bootstrap" ./lambda

echo "Created: $OUTPUT_DIR/bootstrap"

# ZIP パッケージを作成
cd "$OUTPUT_DIR"
zip -j sugu-lambda.zip bootstrap

echo "Created: $OUTPUT_DIR/sugu-lambda.zip"
echo ""
echo "Deploy to AWS Lambda:"
echo "  Runtime: provided.al2023"
echo "  Architecture: x86_64"
echo "  Handler: bootstrap"
