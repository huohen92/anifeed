#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
VERSION="${1:?usage: scripts/set-asset-version.sh 0.03.14.1}"
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "asset version must look like 0.03.14.1" >&2
  exit 1
fi
perl -0pi -e "s/(\\?v=)[0-9.]+/\${1}${VERSION}/g" webui/assets/index.html
echo "asset version: $VERSION"
