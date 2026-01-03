#!/bin/bash

# Sugu CLI ビルドスクリプト

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="$PROJECT_ROOT/dist"

echo "Building Sugu CLI..."

# 出力ディレクトリを作成
mkdir -p "$OUTPUT_DIR"

cd "$PROJECT_ROOT"
go build -o "$OUTPUT_DIR/sugu" .

echo "Created: $OUTPUT_DIR/sugu"
