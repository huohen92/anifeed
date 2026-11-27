#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
VERSION="${VERSION:-0.03.18}"
if ! grep -q "?v=${VERSION}\\." webui/assets/index.html; then
  echo "index.html is not using a ${VERSION}.* asset version" >&2
  exit 1
fi
for asset in app.min.js theme.min.js styles.min.css personalization.min.css settings-layout.min.css; do
  if ! grep -q "/${asset}?v=${VERSION}\\." webui/assets/index.html; then
    echo "${asset} is missing a versioned reference" >&2
    exit 1
  fi
done
for module in app-settings.min.js app-ai.min.js; do
  if [ ! -f "webui/assets/${module}" ]; then
    echo "${module} is missing" >&2
    exit 1
  fi
done
