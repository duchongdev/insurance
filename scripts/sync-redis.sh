#!/usr/bin/env bash
# 根据 .env 中 REDIS_PASSWORD 写入 BRIDGE_REDIS_PASSWORD（供 install.sh / start.sh 调用）
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

[[ -f .env ]] || exit 0

read_env() {
  local key=$1
  grep -E "^${key}=" .env | head -1 | cut -d= -f2- | sed 's/^"//;s/"$//'
}

REDIS_PASSWORD="$(read_env REDIS_PASSWORD)"
[[ -n "$REDIS_PASSWORD" ]] || exit 0

grep -v '^BRIDGE_REDIS_PASSWORD=' .env > .env.tmp
mv .env.tmp .env
echo "BRIDGE_REDIS_PASSWORD=\"${REDIS_PASSWORD}\"" >> .env
