#!/bin/sh

# Falha imediata em qualquer erro, variável não definida ou erro em pipe
set -eu

# ----------------------------
# Configurações
# ----------------------------
MIGRATIONS_DIR="./migrations"

# ----------------------------
# Validações iniciais
# ----------------------------

# Verifica se o comando migrate está instalado
if ! command -v migrate >/dev/null 2>&1; then
  echo "❌ ERRO: 'migrate' não está instalado ou não está no PATH"
  exit 1
fi

# Verifica se foi passado ao menos um argumento
if [ "$#" -eq 0 ]; then
  echo "❌ Uso correto:"
  echo "  ./create-migration.sh <nome da migration>"
  echo "Exemplo:"
  echo "  ./create-migration.sh create movies table"
  exit 1
fi

# ----------------------------
# Normalização do nome
# ----------------------------

# Converte:
# - múltiplos argumentos → um nome
# - espaços → _
# - letras maiúsculas → minúsculas
# - remove caracteres inválidos
MIGRATION_NAME="$(echo "$*" \
  | tr '[:upper:]' '[:lower:]' \
  | sed -E 's/[^a-z0-9 ]+//g' \
  | tr ' ' '_' \
  | sed -E 's/_+/_/g' \
  | sed -E 's/^_|_$//g')"

# Validação final do nome
if [ -z "$MIGRATION_NAME" ]; then
  echo "❌ ERRO: nome da migration inválido após normalização"
  exit 1
fi

# ----------------------------
# Execução
# ----------------------------

echo "📄 Criando migration: $MIGRATION_NAME"

migrate create \
  -seq \
  -ext=".sql" \
  -dir="$MIGRATIONS_DIR" \
  "$MIGRATION_NAME"

echo "✅ Migration '$MIGRATION_NAME' criada com sucesso em ./migrations"