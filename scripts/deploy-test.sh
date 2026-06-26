#!/usr/bin/env bash
# 打包并部署到测试环境：本机构建镜像 → 停服 → 上传 → 解压 → 加载镜像 → 启动
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

DEPLOY_HOST="${DEPLOY_HOST:-root@10.41.61.41}"
DEPLOY_DIR="${DEPLOY_DIR:-/home/ins}"
BRIDGE_IMAGE="${BRIDGE_IMAGE:-insurance-bridge:deploy}"
# 测试服务器已有 nginx 镜像，勿从 Docker Hub 拉取；可通过环境变量覆盖
NGINX_IMAGE="${NGINX_IMAGE:-nginx:1.26.1-alpine}"
REDIS_IMAGE="${REDIS_IMAGE:-redis:7-alpine}"

info() { echo "[deploy] $*"; }
die() { echo "[deploy] 错误: $*" >&2; exit 1; }

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "未找到命令: $1"
}

need_cmd ssh
need_cmd scp
need_cmd docker

info "目标: ${DEPLOY_HOST}:${DEPLOY_DIR}"
info "Nginx 镜像: ${NGINX_IMAGE}（使用服务器已有镜像，不上传）"

# 1. 打包源码
info "1/8 打包源码 ..."
bash scripts/package.sh

ARCHIVE="$(ls -t dist/insurance-bridge-*.tar.gz 2>/dev/null | head -1)"
[[ -n "$ARCHIVE" && -f "$ARCHIVE" ]] || die "未找到 dist/insurance-bridge-*.tar.gz"
ARCHIVE_NAME="$(basename "$ARCHIVE")"
IMAGE_TAR="dist/bridge-image.tar.gz"
REDIS_IMAGE_TAR="dist/redis-image.tar.gz"
info "压缩包: ${ARCHIVE}"

# 2. 本机构建 linux/amd64 镜像（测试服务器无法访问 Docker Hub）
info "2/8 本机构建镜像 ${BRIDGE_IMAGE} (linux/amd64) ..."
docker build --platform linux/amd64 -t "${BRIDGE_IMAGE}" .

# 3. 导出镜像
info "3/8 导出应用镜像 ..."
docker save "${BRIDGE_IMAGE}" | gzip > "${IMAGE_TAR}"

info "4/8 导出 Redis 镜像 (linux/amd64) ..."
docker pull --platform linux/amd64 "${REDIS_IMAGE}" >/dev/null
docker save "${REDIS_IMAGE}" | gzip > "${REDIS_IMAGE_TAR}"

# 4. 停服
info "5/8 停服 ..."
ssh "$DEPLOY_HOST" "if [[ -f ${DEPLOY_DIR}/scripts/stop.sh ]]; then cd ${DEPLOY_DIR} && ./scripts/stop.sh; else echo '跳过停服（首次部署或未安装）'; fi"

# 5. 上传
info "6/8 上传 ..."
ssh "$DEPLOY_HOST" "mkdir -p ${DEPLOY_DIR}"
scp "$ARCHIVE" "${DEPLOY_HOST}:${DEPLOY_DIR}/${ARCHIVE_NAME}"
scp "${IMAGE_TAR}" "${DEPLOY_HOST}:${DEPLOY_DIR}/bridge-image.tar.gz"
scp "${REDIS_IMAGE_TAR}" "${DEPLOY_HOST}:${DEPLOY_DIR}/redis-image.tar.gz"

# 6. 解压并加载镜像
info "7/8 解压并加载镜像 ..."
ssh "$DEPLOY_HOST" "cd ${DEPLOY_DIR} && \
  tar -xzf ${ARCHIVE_NAME} && \
  chmod +x scripts/*.sh && \
  find . -name '._*' -delete && find . -name '.DS_Store' -delete && \
  gunzip -c bridge-image.tar.gz | docker load && \
  gunzip -c redis-image.tar.gz | docker load && \
  if grep -q '^BRIDGE_IMAGE=' .env 2>/dev/null; then \
    sed -i 's#^BRIDGE_IMAGE=.*#BRIDGE_IMAGE=${BRIDGE_IMAGE}#' .env; \
  else \
    echo 'BRIDGE_IMAGE=${BRIDGE_IMAGE}' >> .env; \
  fi && \
  if grep -q '^REDIS_IMAGE=' .env 2>/dev/null; then \
    sed -i 's#^REDIS_IMAGE=.*#REDIS_IMAGE=${REDIS_IMAGE}#' .env; \
  else \
    echo 'REDIS_IMAGE=${REDIS_IMAGE}' >> .env; \
  fi && \
  if grep -q '^NGINX_IMAGE=' .env 2>/dev/null; then \
    sed -i 's#^NGINX_IMAGE=.*#NGINX_IMAGE=${NGINX_IMAGE}#' .env; \
  else \
    echo 'NGINX_IMAGE=${NGINX_IMAGE}' >> .env; \
  fi && \
  docker image inspect '${NGINX_IMAGE}' >/dev/null || { echo '服务器缺少 Nginx 镜像: ${NGINX_IMAGE}'; exit 3; }"

# 7. 启动（不在服务器构建，使用 .env 中已有镜像）
info "8/8 启动 ..."
ssh "$DEPLOY_HOST" "cd ${DEPLOY_DIR} && \
  if [[ ! -f .env || ! -f config/config.yaml ]]; then \
    echo '缺少 .env 或 config/config.yaml'; exit 1; \
  fi && \
  if grep -q '请' .env config/config.yaml 2>/dev/null; then \
    echo '配置仍为占位符，请编辑 .env 与 config/config.yaml'; exit 2; \
  fi && \
  bash scripts/sync-dsn.sh && \
  bash scripts/sync-redis.sh && \
  DOCKER_API_VERSION=1.43 ./scripts/start.sh"

# 8. 健康检查
info "健康检查 ..."
ssh "$DEPLOY_HOST" "curl -sf http://127.0.0.1/health/ready >/dev/null && echo ready OK" \
  || die "健康检查失败: ssh ${DEPLOY_HOST} 'cd ${DEPLOY_DIR} && docker compose logs bridge nginx'"

info "部署完成"
info "  管理后台: http://10.41.61.41/ （自动进入 /admin/）"
info "  健康检查: http://10.41.61.41/health/ready"
