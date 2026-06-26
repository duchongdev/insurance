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
bash scripts/sync-redis.sh

if [[ ! -f deploy/admin-dist/index.html ]]; then
  if [[ -f web/admin/dist/index.html ]]; then
    echo "[start] 同步 web/admin/dist → deploy/admin-dist ..."
    rm -rf deploy/admin-dist
    cp -a web/admin/dist deploy/admin-dist
  elif command -v npm >/dev/null 2>&1; then
    echo "[start] 未找到 deploy/admin-dist，正在构建管理后台 ..."
    bash scripts/build-admin.sh
  else
    echo "[start] 错误: 未找到 deploy/admin-dist/index.html"
    echo "[start] 请在开发机执行 bash scripts/build-admin.sh 后重新打包部署"
    exit 1
  fi
fi

export DOCKER_API_VERSION="${DOCKER_API_VERSION:-1.43}"

if docker compose version >/dev/null 2>&1; then
  COMPOSE="docker compose"
else
  COMPOSE="docker-compose"
fi

configure_nginx() {
  local ssl_dir="$ROOT/deploy/ssl"
  local conf_dir="$ROOT/deploy/nginx/conf.d"
  if [[ -f "$ssl_dir/fullchain.pem" && -f "$ssl_dir/privkey.pem" ]]; then
    cp "$conf_dir/bridge-ssl.conf.example" "$conf_dir/bridge-ssl.conf"
    cp "$conf_dir/bridge-http-redirect.conf.example" "$conf_dir/bridge.conf"
    echo "[start] 已启用 HTTPS :443（证书 deploy/ssl/）"
  else
    rm -f "$conf_dir/bridge-ssl.conf"
    cp "$conf_dir/bridge-http-serve.conf.example" "$conf_dir/bridge.conf"
    echo "[start] 未配置 SSL 证书，仅 HTTP :80（投产请将 fullchain.pem / privkey.pem 放入 deploy/ssl/）"
  fi
}

configure_nginx

BRIDGE_IMAGE="$(read_env BRIDGE_IMAGE)"
NGINX_IMAGE="$(read_env NGINX_IMAGE)"
NGINX_IMAGE="${NGINX_IMAGE:-nginx:1.26-alpine}"

images_ready=true
if [[ -n "$BRIDGE_IMAGE" ]] && ! docker image inspect "$BRIDGE_IMAGE" >/dev/null 2>&1; then
  images_ready=false
fi
if ! docker image inspect "$NGINX_IMAGE" >/dev/null 2>&1; then
  echo "[start] 警告: 本地无 Nginx 镜像 ${NGINX_IMAGE}，compose 可能尝试拉取"
  images_ready=false
fi

if [[ "$images_ready" == true && -n "$BRIDGE_IMAGE" ]]; then
  echo "[start] 使用镜像 ${BRIDGE_IMAGE}、${NGINX_IMAGE} 启动 (--no-build) ..."
  $COMPOSE up -d --no-build
else
  echo "[start] 构建并启动服务 ..."
  $COMPOSE up -d --build
fi

echo
echo "[start] 启动完成"
if [[ -f deploy/ssl/fullchain.pem && -f deploy/ssl/privkey.pem ]]; then
  echo "  管理后台: https://<域名或IP>/"
  echo "  健康检查: https://<域名或IP>/health/ready"
else
  echo "  管理后台: http://<服务器IP>/"
  echo "  健康检查: http://<服务器IP>/health/ready"
fi
echo "  查看日志: $COMPOSE logs -f nginx bridge"
echo "  停止服务: ./scripts/stop.sh"
