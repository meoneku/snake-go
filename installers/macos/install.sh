#!/usr/bin/env bash
# Installer Snake Go ke /Applications di macOS
#
# Penggunaan:
#   ./install.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DIST_DIR="$PROJECT_ROOT/dist"
APP_SOURCE="$DIST_DIR/Snake-Go.app"

echo "=============================================="
echo "      Instalasi Snake Go ke macOS             "
echo "=============================================="

# Jika belum ada bundle .app, jalankan script packager
if [ ! -d "$APP_SOURCE" ]; then
  echo "[*] Bundle belum ada. Menjalankan packaging..."
  bash "$SCRIPT_DIR/package_macos.sh"
fi

TARGET_APP="/Applications/Snake-Go.app"
echo "[*] Memasang ke $TARGET_APP..."
rm -rf "$TARGET_APP"
cp -R "$APP_SOURCE" "/Applications/"

# Pasang binary terminal ke /usr/local/bin jika tersedia
TERM_BIN=""
for b in "$DIST_DIR/snake-go-terminal-darwin-arm64" \
         "$DIST_DIR/snake-go-terminal-darwin-amd64" \
         "$DIST_DIR/snake-terminal-darwin" \
         "$PROJECT_ROOT/snake-go"; do
  if [ -f "$b" ]; then
    TERM_BIN="$b"
    break
  fi
done

if [ -n "$TERM_BIN" ]; then
  if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
    cp "$TERM_BIN" "/usr/local/bin/snake-go"
    chmod +x "/usr/local/bin/snake-go"
    echo "[+] Binary terminal terpasang di /usr/local/bin/snake-go"
  fi
fi

echo "=============================================="
echo "    Instalasi Berhasil!                       "
echo "    Buka 'Snake Go' dari Launchpad            "
echo "    atau folder /Applications.                "
echo "=============================================="
