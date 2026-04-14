#!/bin/bash
set -euo pipefail

echo "==> Generando templates templ..."
~/go/bin/templ generate

echo "==> Compilando para Raspberry Pi (linux/arm64)..."
mkdir -p dist
GOOS=linux GOARCH=arm64 go build -o dist/porfolio-amelia .

echo "==> Listo: dist/porfolio-amelia ($(du -h dist/porfolio-amelia | cut -f1))"
