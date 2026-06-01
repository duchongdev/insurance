#!/usr/bin/env bash
# 根据 .env 中 MYSQL_* 写入 BRIDGE_DATABASE_DSN（供 install.sh / start.sh 调用）
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

[[ -f .env ]] || exit 0

read_env() {
  local key=$1
  grep -E "^${key}=" .env | head -1 | cut -d= -f2- | sed 's/^"//;s/"$//'
}

MYSQL_USER="$(read_env MYSQL_USER)"
MYSQL_PASSWORD="$(read_env MYSQL_PASSWORD)"
MYSQL_DATABASE="$(read_env MYSQL_DATABASE)"

[[ -n "$MYSQL_USER" && -n "$MYSQL_PASSWORD" && -n "$MYSQL_DATABASE" ]] || exit 0

dsn="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(mysql:3306)/${MYSQL_DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"

# 删除旧行后追加，避免 sed 对 DSN 中 & 等特殊字符误替换
grep -v '^BRIDGE_DATABASE_DSN=' .env > .env.tmp
mv .env.tmp .env
echo "BRIDGE_DATABASE_DSN=\"${dsn}\"" >> .env
