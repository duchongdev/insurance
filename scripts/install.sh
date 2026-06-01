#!/usr/bin/env bash
# 首次部署：检查环境、复制配置模板、同步数据库连接串。
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

info() { echo "[install] $*"; }
warn() { echo "[install] 警告: $*" >&2; }
die() { echo "[install] 错误: $*" >&2; exit 1; }

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "未找到命令: $1"
}

need_cmd docker
if docker compose version >/dev/null 2>&1; then
  COMPOSE="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE="docker-compose"
else
  die "请先安装 Docker Compose（docker compose 或 docker-compose）"
fi

# ---------- 配置文件 ----------
if [[ ! -f .env ]]; then
  cp .env.example .env
  info "已生成 .env（请编辑后再启动）"
  CREATED_ENV=1
else
  info "已存在 .env"
fi

if [[ ! -f config/config.yaml ]]; then
  cp config/config.yaml.example config/config.yaml
  info "已生成 config/config.yaml（请编辑后再启动）"
  CREATED_CFG=1
else
  info "已存在 config/config.yaml"
fi

mkdir -p logs

bash scripts/sync-dsn.sh
info "已根据 MYSQL_* 写入 BRIDGE_DATABASE_DSN"

# ---------- 基本校验 ----------
check_placeholder() {
  local file=$1 val=$2 name=$3
  if echo "$val" | grep -q '请'; then
    warn "${file} 中 ${name} 仍为占位符，启动前请务必修改"
    NEED_EDIT=1
  fi
}

NEED_EDIT=0
# shellcheck disable=SC1091
source .env
check_placeholder ".env" "${MYSQL_ROOT_PASSWORD:-}" "MYSQL_ROOT_PASSWORD"
check_placeholder ".env" "${MYSQL_PASSWORD:-}" "MYSQL_PASSWORD"
check_placeholder ".env" "${BRIDGE_SECURITY_DATA_ENCRYPTION_KEY:-}" "BRIDGE_SECURITY_DATA_ENCRYPTION_KEY"
check_placeholder ".env" "${BRIDGE_SECURITY_JWT_SECRET:-}" "BRIDGE_SECURITY_JWT_SECRET"
check_placeholder ".env" "${BRIDGE_ADMIN_DEFAULT_PASSWORD:-}" "BRIDGE_ADMIN_DEFAULT_PASSWORD"

if [[ "${#BRIDGE_SECURITY_DATA_ENCRYPTION_KEY}" -ne 32 ]]; then
  warn "BRIDGE_SECURITY_DATA_ENCRYPTION_KEY 必须为 32 字节（当前长度 ${#BRIDGE_SECURITY_DATA_ENCRYPTION_KEY}）"
  NEED_EDIT=1
fi

if [[ -n "${CREATED_ENV:-}" || -n "${CREATED_CFG:-}" || "$NEED_EDIT" -eq 1 ]]; then
  echo
  info "请编辑以下文件后执行: ./scripts/start.sh"
  echo "  - .env"
  echo "  - config/config.yaml"
  echo
  info "详细说明见 DEPLOY.md"
  exit 0
fi

info "配置检查通过，可执行 ./scripts/start.sh 启动服务"
