#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
ESBUILD="${ESBUILD:-$(command -v esbuild || true)}"
if [ -z "$ESBUILD" ]; then
  echo "esbuild is required (npm install -g esbuild)" >&2
  exit 1
fi
JS=(app app-settings app-ai theme)
CSS=(styles personalization settings-layout)
for name in "${JS[@]}"; do
  "$ESBUILD" "webui/assets/${name}.js" --minify --target=es2020 --legal-comments=none --outfile="webui/assets/${name}.min.js"
done
for name in "${CSS[@]}"; do
  "$ESBUILD" "webui/assets/${name}.css" --minify --target=es2020 --legal-comments=none --outfile="webui/assets/${name}.min.css"
done
