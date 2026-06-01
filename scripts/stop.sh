#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if docker compose version >/dev/null 2>&1; then
  COMPOSE="docker compose"
else
  COMPOSE="docker-compose"
fi

echo "[stop] 停止服务 ..."
$COMPOSE down
echo "[stop] 已停止（数据库数据保留在 Docker 卷 mysql_data 中）"
