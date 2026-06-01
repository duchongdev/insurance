#!/usr/bin/env bash
# 构建管理后台 Vue 前端 → web/admin/dist/
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN_DIR="$ROOT/web/admin"

if ! command -v npm >/dev/null 2>&1; then
  echo "[build-admin] 错误: 需要 Node.js 18+ 与 npm"
  exit 1
fi

cd "$ADMIN_DIR"
echo "[build-admin] 安装依赖 ..."
if [[ -f package-lock.json ]]; then
  npm ci
else
  npm install
fi
echo "[build-admin] 构建 ..."
npm run build
echo "[build-admin] 完成: web/admin/dist/"

# macOS bsdtar 会跳过 .gitignore 中的 dist，复制到未忽略目录供打包与 Nginx 挂载
rm -rf "$ROOT/deploy/admin-dist"
cp -a dist "$ROOT/deploy/admin-dist"
echo "[build-admin] 已同步至 deploy/admin-dist/"
