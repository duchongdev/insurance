#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

read_env() {
  local key=$1
  grep -E "^${key}=" .env 2>/dev/null | head -1 | cut -d= -f2- | sed 's/^"//;s/"$//'
}

if [[ ! -f .env || ! -f config/config.yaml ]]; then
  echo "[start] 未找到配置文件，正在运行 install.sh ..."
  bash scripts/install.sh
  if [[ ! -f .env || ! -f config/config.yaml ]]; then
    exit 1
  fi
  if grep -q '请' .env 2>/dev/null; then
    echo "[start] 请先修改 .env 与 config/config.yaml 中的占位符后再启动"
    exit 1
  fi
fi

mkdir -p logs
bash scripts/sync-dsn.sh

export DOCKER_API_VERSION="${DOCKER_API_VERSION:-1.43}"

if docker compose version >/dev/null 2>&1; then
  COMPOSE="docker compose"
else
  COMPOSE="docker-compose"
fi

APP_PORT="$(read_env APP_PORT)"
APP_PORT="${APP_PORT:-5051}"
BRIDGE_IMAGE="$(read_env BRIDGE_IMAGE)"

if [[ -n "$BRIDGE_IMAGE" ]] && docker image inspect "$BRIDGE_IMAGE" >/dev/null 2>&1; then
  echo "[start] 使用镜像 ${BRIDGE_IMAGE} 启动 (--no-build) ..."
  $COMPOSE up -d --no-build
else
  echo "[start] 构建并启动服务 ..."
  $COMPOSE up -d --build
fi

echo
echo "[start] 启动完成"
echo "  管理后台: http://<服务器IP>:${APP_PORT}/admin/"
echo "  健康检查: http://<服务器IP>:${APP_PORT}/health/ready"
echo "  查看日志: $COMPOSE logs -f bridge"
echo "  停止服务: ./scripts/stop.sh"
