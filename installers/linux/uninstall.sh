#!/usr/bin/env bash
# Uninstaller Snake Go untuk Linux
#
# Penggunaan:
#   ./uninstall.sh           # User mode
#   sudo ./uninstall.sh --system   # System mode

set -e

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

echo "=============================================="
echo "       Pencopotan (Uninstall) Snake Go        "
echo "=============================================="

rm -f "$PREFIX/bin/snake-go"
rm -f "$PREFIX/bin/snake-go-gui"
rm -f "$PREFIX/bin/snake"
rm -f "$PREFIX/bin/snake-gui"
rm -f "$DATADIR/applications/snake-go.desktop"
rm -f "$DATADIR/applications/snake.desktop"
rm -f "$DATADIR/icons/hicolor/scalable/apps/snake-go.svg"
rm -f "$DATADIR/icons/hicolor/256x256/apps/snake-go.png"
rm -f "$DATADIR/icons/hicolor/48x48/apps/snake-go.png"
rm -f "$DATADIR/icons/hicolor/16x16/apps/snake-go.png"
rm -f "$DATADIR/icons/hicolor/scalable/apps/snake.svg"
rm -f "$DATADIR/icons/hicolor/256x256/apps/snake.png"
rm -f "$DATADIR/icons/hicolor/48x48/apps/snake.png"
rm -f "$DATADIR/icons/hicolor/16x16/apps/snake.png"

if command -v update-desktop-database >/dev/null 2>&1; then
  update-desktop-database "$DATADIR/applications" 2>/dev/null || true
fi
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
  gtk-update-icon-cache -f -t "$DATADIR/icons/hicolor" 2>/dev/null || true
fi

echo "[+] Berkas instalasi berhasil dicopot dari sistem."
echo "[*] Catatan: Data skor di ~/.config/snake-go/scores.json tetap disimpan."
echo "=============================================="
echo "       Snake Go telah berhasil dihapus.       "
echo "=============================================="
