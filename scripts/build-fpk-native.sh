#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if ! command -v fnpack >/dev/null 2>&1; then
  echo "fnpack is required: https://developer.fnnas.com/docs/cli/fnpack" >&2
  exit 1
fi
if ! command -v musl-gcc >/dev/null 2>&1; then
  echo "musl-gcc is required (Debian: apt install musl-tools)" >&2
  exit 1
fi

VERSION="$(sed -n 's/^version[[:space:]]*=[[:space:]]*//p' fpk-native/manifest | head -n 1 | tr -d '[:space:]')"
[ -n "$VERSION" ] || { echo "cannot read version from fpk-native/manifest" >&2; exit 1; }

export CGO_ENABLED=1 GOOS=linux GOARCH=amd64
export CC=musl-gcc
${GO:-go} build -trimpath \
  -ldflags "-s -w -extldflags -static -X github.com/huohen92/anifeed/version.Version=0.02" \
  -o fpk-native/app/bin/anifeed .

mkdir -p dist
fnpack build --directory fpk-native
mv -f anifeed-native.fpk "dist/AniFeed-${VERSION}-x86-native.fpk"
echo "built: dist/AniFeed-${VERSION}-x86-native.fpk"
