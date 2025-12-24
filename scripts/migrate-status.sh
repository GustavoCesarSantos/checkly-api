#!/bin/sh

# Falha imediata em qualquer erro, variável não definida ou erro em pipe
set -eu

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
MIGRATIONS_DIR="$ROOT_DIR/migrations"
ENV_FILE="$ROOT_DIR/.env"

if [ -f "$ENV_FILE" ]; then
  export $(grep -v '^#' "$ENV_FILE" | xargs)
fi

if ! command -v migrate >/dev/null 2>&1; then
  echo "❌ ERRO: 'migrate' não está instalado"
  exit 1
fi

if [ -z "$DB_DSN" ]; then
  echo "❌ ERRO: variável DB_DSN não definida"
  exit 1
fi

echo "📊 Status das migrations:"

migrate \
  -path="$MIGRATIONS_DIR" \
  -database="$DB_DSN" \
  version
