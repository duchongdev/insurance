#!/usr/bin/env bash
# 生成交付压缩包：dist/insurance-bridge-<version>-<date>.tar.gz
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

VERSION="$(grep '^go ' go.mod | awk '{print $2}')"
DATE="$(date +%Y%m%d)"
NAME="insurance-bridge-${VERSION}-${DATE}"
OUT_DIR="dist"
ARCHIVE="${OUT_DIR}/${NAME}.tar.gz"

mkdir -p "$OUT_DIR"

echo "[package] 生成交付包: ${ARCHIVE}"

# macOS：不打包扩展属性 / 资源分叉，避免 ._xxx 等垃圾文件
export COPYFILE_DISABLE=1
TAR_EXTRA=()
if [[ "$(uname -s)" == "Darwin" ]]; then
  TAR_EXTRA=(--no-xattrs)
fi

tar -czf "$ARCHIVE" \
  ${TAR_EXTRA+"${TAR_EXTRA[@]}"} \
  --exclude='.git' \
  --exclude='.cursor' \
  --exclude='.idea' \
  --exclude='.vscode' \
  --exclude='logs' \
  --exclude='*.log' \
  --exclude='bin' \
  --exclude='bridge' \
  --exclude='dist' \
  --exclude='.env' \
  --exclude='config/config.yaml' \
  --exclude='.DS_Store' \
  --exclude='**/.DS_Store' \
  --exclude='__MACOSX' \
  --exclude='**/__MACOSX' \
  --exclude='**/__MACOSX/**' \
  --exclude='._*' \
  --exclude='**/._*' \
  --exclude='.AppleDouble' \
  --exclude='**/.AppleDouble' \
  --exclude='.LSOverride' \
  --exclude='.Spotlight-V100' \
  --exclude='.Trashes' \
  --exclude='.fseventsd' \
  --exclude='.TemporaryItems' \
  --exclude='.DocumentRevisions-V100' \
  -C "$ROOT" \
  .

echo "[package] 完成: ${ARCHIVE} ($(du -h "$ARCHIVE" | awk '{print $1}'))"
echo "[package] 交付后请让用户: 解压 → 编辑 .env / config → ./scripts/start.sh"
