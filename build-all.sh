#!/usr/bin/env bash
# Build Snake Go untuk Linux, Windows, dan macOS dengan versioning terpusat.
# Hasil ada di dist/
set -e
cd "$(dirname "$0")"
mkdir -p dist

VERSION="1.0.1"
if [ -f "VERSION" ]; then
  VERSION="$(cat VERSION | tr -d '[:space:]')"
fi
GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "none")"
BUILD_DATE="$(date -u +%Y-%m-%d)"

LDFLAGS="-X banguncandi/pkg/version.Version=${VERSION} -X banguncandi/pkg/version.GitCommit=${GIT_COMMIT} -X banguncandi/pkg/version.BuildDate=${BUILD_DATE}"
LDFLAGS_WINGUI="-H windowsgui $LDFLAGS"

echo "=============================================="
echo "    Build Snake Go v$VERSION ($GIT_COMMIT)    "
echo "=============================================="

echo "== terminal (murni stdlib, tanpa CGO) =="
GOOS=linux   GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/snake-go-terminal-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/snake-go-terminal-windows-amd64.exe .
GOOS=darwin  GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/snake-go-terminal-darwin-amd64 .
GOOS=darwin  GOARCH=arm64 go build -ldflags="$LDFLAGS" -o dist/snake-go-terminal-darwin-arm64 .

echo "== desktop GUI (Ebitengine, butuh CGO) =="
if [ "$(go env GOOS)" = "linux" ]; then
  (cd gui && go build -ldflags="$LDFLAGS" -o ../dist/snake-go-gui-linux-amd64 .)
fi

if command -v x86_64-w64-mingw32-gcc >/dev/null; then
  # Kompilasi icon resource Windows jika windres tersedia
  if command -v x86_64-w64-mingw32-windres >/dev/null; then
    (cd gui && x86_64-w64-mingw32-windres -i snake-go.rc -O coff -o icon_windows_amd64.syso)
  fi
  # -H windowsgui = tanpa jendela console saat .exe diklik
  (cd gui && CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
    go build -ldflags="$LDFLAGS_WINGUI" -o ../dist/snake-go-gui-windows-amd64.exe .)
else
  echo "!! lewati GUI Windows: butuh mingw-w64 (Arch: pacman -S mingw-w64-gcc)"
fi
echo "!! GUI macOS: build di Mac langsung (cd gui && go build .), cross-compile dari Linux butuh osxcross"

echo "== selesai =="
ls -la dist/
