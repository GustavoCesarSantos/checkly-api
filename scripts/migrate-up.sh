#!/bin/sh

# Falha imediata em qualquer erro, variável não definida ou erro em pipe
set -eu

# ----------------------------
# Configurações
# ----------------------------
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
MIGRATIONS_DIR="$ROOT_DIR/migrations"
ENV_FILE="$ROOT_DIR/.env"

# ----------------------------
# Carrega .env (se existir)
# ----------------------------
if [ -f "$ENV_FILE" ]; then
  export $(grep -v '^#' "$ENV_FILE" | xargs)
fi

# ----------------------------
# Validações
# ----------------------------
if ! command -v migrate >/dev/null 2>&1; then
  echo "❌ ERRO: 'migrate' não está instalado"
  exit 1
fi

if [ -z "$DB_DSN" ]; then
  echo "❌ ERRO: variável DB_DSN não definida"
  exit 1
fi

# ----------------------------
# Execução
# ----------------------------
echo "🚀 Aplicando migrations..."

migrate \
  -path="$MIGRATIONS_DIR" \
  -database="$DB_DSN" \
  up

echo "✅ Migrations aplicadas com sucesso"
