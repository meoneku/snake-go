#!/usr/bin/env bash
# Build Ular-ularan untuk Linux, Windows, dan macOS.
# Hasil ada di dist/
set -e
cd "$(dirname "$0")"
mkdir -p dist

echo "== terminal (murni stdlib, tanpa CGO) =="
GOOS=linux   GOARCH=amd64 go build -o dist/ular-terminal-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o dist/ular-terminal-windows-amd64.exe .
GOOS=darwin  GOARCH=amd64 go build -o dist/ular-terminal-darwin-amd64 .
GOOS=darwin  GOARCH=arm64 go build -o dist/ular-terminal-darwin-arm64 .

echo "== desktop GUI (Ebitengine, butuh CGO) =="
if [ "$(go env GOOS)" = "linux" ]; then
  (cd gui && go build -o ../dist/ular-gui-linux-amd64 .)
fi
if command -v x86_64-w64-mingw32-gcc >/dev/null; then
  # -H windowsgui = tanpa jendela console saat .exe diklik
  (cd gui && CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 \
    go build -ldflags="-H windowsgui" -o ../dist/ular-gui-windows-amd64.exe .)
else
  echo "!! lewati GUI Windows: butuh mingw-w64 (Arch: pacman -S mingw-w64-gcc)"
fi
echo "!! GUI macOS: build di Mac langsung (cd gui && go build .), cross-compile dari Linux butuh osxcross"

echo "== selesai =="
ls -la dist/
