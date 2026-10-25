#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if ! command -v fnpack >/dev/null 2>&1; then
  echo "fnpack is required: https://developer.fnnas.com/docs/cli/fnpack" >&2
  exit 1
fi

VERSION="$(sed -n 's/^version[[:space:]]*=[[:space:]]*//p' fpk/manifest | head -n 1 | tr -d '[:space:]')"
[ -n "$VERSION" ] || { echo "cannot read version from fpk/manifest" >&2; exit 1; }

mkdir -p dist
fnpack build --directory fpk
mv -f anifeed.fpk "dist/anifeed-${VERSION}-x86.fpk"
echo "built: dist/anifeed-${VERSION}-x86.fpk"
