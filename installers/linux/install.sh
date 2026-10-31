#!/usr/bin/env bash
# Installer Snake Go untuk Linux
#
# Penggunaan:
#   ./install.sh           # User mode (ke ~/.local, tanpa sudo)
#   sudo ./install.sh --system   # System mode (ke /usr/local)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Mode instalasi
if [ "$1" = "--system" ]; then
  if [ "$(id -u)" -ne 0 ]; then
    echo "Error: Opsi --system membutuhkan hak akses root (sudo)." >&2
    exit 1
  fi
  PREFIX="/usr/local"
  DATADIR="/usr/local/share"
else
  PREFIX="$HOME/.local"
  DATADIR="$HOME/.local/share"
fi

VERSION="1.0.0"
if [ -f "$PROJECT_ROOT/VERSION" ]; then
  VERSION="$(cat "$PROJECT_ROOT/VERSION" | tr -d '[:space:]')"
fi

echo "=============================================="
echo "      Instalasi Snake Go v$VERSION (Linux)    "
echo "=============================================="
echo "[*] Lokasi prefix: $PREFIX"

# 1. Cari binary GUI & Terminal
GUI_BIN=""
for b in "$PROJECT_ROOT/dist/snake-go-gui-linux-amd64" \
         "$PROJECT_ROOT/dist/ular-gui-linux-amd64" \
         "$PROJECT_ROOT/gui/snake-go-gui" \
         "$PROJECT_ROOT/gui/snake-gui" \
         "$PROJECT_ROOT/snake-go-gui" \
         "$PROJECT_ROOT/ular-gui"; do
  if [ -f "$b" ]; then
    GUI_BIN="$b"
    break
  fi
done

TERM_BIN=""
for b in "$PROJECT_ROOT/dist/snake-go-terminal-linux-amd64" \
         "$PROJECT_ROOT/dist/ular-terminal-linux-amd64" \
         "$PROJECT_ROOT/snake-go" \
         "$PROJECT_ROOT/snake" \
         "$PROJECT_ROOT/ular-terminal"; do
  if [ -f "$b" ]; then
    TERM_BIN="$b"
    break
  fi
done

# 2. Siapkan direktori
mkdir -p "$PREFIX/bin"
mkdir -p "$DATADIR/applications"
mkdir -p "$DATADIR/icons/hicolor/scalable/apps"
mkdir -p "$DATADIR/icons/hicolor/256x256/apps"
mkdir -p "$DATADIR/icons/hicolor/48x48/apps"
mkdir -p "$DATADIR/icons/hicolor/16x16/apps"

# 3. Pasang Executable
if [ -n "$GUI_BIN" ]; then
  cp "$GUI_BIN" "$PREFIX/bin/snake-go-gui"
  chmod +x "$PREFIX/bin/snake-go-gui"
  echo "[+] Terpasang: $PREFIX/bin/snake-go-gui"
else
  echo "[!] Peringatan: Binary GUI belum di-build. Bangun dulu via 'cd gui && go build -o snake-go-gui .'"
fi

if [ -n "$TERM_BIN" ]; then
  cp "$TERM_BIN" "$PREFIX/bin/snake-go"
  chmod +x "$PREFIX/bin/snake-go"
  echo "[+] Terpasang: $PREFIX/bin/snake-go (Terminal)"
fi

# 4. Pasang Ikon Hicolor FreeDesktop
ASSETS_DIR="$PROJECT_ROOT/gui/assets"
if [ -d "$ASSETS_DIR" ]; then
  [ -f "$ASSETS_DIR/icon.svg" ] && cp "$ASSETS_DIR/icon.svg" "$DATADIR/icons/hicolor/scalable/apps/snake-go.svg"
  [ -f "$ASSETS_DIR/icon_256.png" ] && cp "$ASSETS_DIR/icon_256.png" "$DATADIR/icons/hicolor/256x256/apps/snake-go.png"
  [ -f "$ASSETS_DIR/icon_48.png" ] && cp "$ASSETS_DIR/icon_48.png" "$DATADIR/icons/hicolor/48x48/apps/snake-go.png"
  [ -f "$ASSETS_DIR/icon_16.png" ] && cp "$ASSETS_DIR/icon_16.png" "$DATADIR/icons/hicolor/16x16/apps/snake-go.png"
  echo "[+] Aset ikon FreeDesktop hicolor berhasil disalin."
fi

# 5. Pasang Desktop Entry
DESKTOP_SRC="$PROJECT_ROOT/snake-go.desktop"
if [ -f "$DESKTOP_SRC" ]; then
  cp "$DESKTOP_SRC" "$DATADIR/applications/snake-go.desktop"
  chmod +r "$DATADIR/applications/snake-go.desktop"
  echo "[+] Desktop entry terpasang di $DATADIR/applications/snake-go.desktop"
fi

# 6. Perbarui Cache Sistem
if command -v update-desktop-database >/dev/null 2>&1; then
  update-desktop-database "$DATADIR/applications" 2>/dev/null || true
fi
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
  gtk-update-icon-cache -f -t "$DATADIR/icons/hicolor" 2>/dev/null || true
fi

echo "=============================================="
echo "    Instalasi Berhasil!                       "
echo "    Buka 'Snake Go' dari menu aplikasi,       "
echo "    atau jalankan perintah 'snake-go-gui'.    "
echo "=============================================="
