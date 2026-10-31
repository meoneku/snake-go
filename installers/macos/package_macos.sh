#!/usr/bin/env bash
# Script pembuatan macOS App Bundle (Snake-Go.app) & Disk Image (.dmg)
# Dijalankan di macOS (atau runner CI macOS)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

VERSION="1.0.0"
if [ -f "$PROJECT_ROOT/VERSION" ]; then
  VERSION="$(cat "$PROJECT_ROOT/VERSION" | tr -d '[:space:]')"
fi

echo "=============================================="
echo "    Pembuatan macOS App Bundle: Snake Go      "
echo "    Versi: $VERSION                           "
echo "=============================================="

DIST_DIR="$PROJECT_ROOT/dist"
mkdir -p "$DIST_DIR"

APP_NAME="Snake-Go.app"
APP_DIR="$DIST_DIR/$APP_NAME"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

rm -rf "$APP_DIR"
mkdir -p "$MACOS_DIR" "$RESOURCES_DIR"

# 1. Salin Info.plist
cp "$SCRIPT_DIR/Info.plist" "$CONTENTS_DIR/Info.plist"
sed -i '' "s/1.0.0/$VERSION/g" "$CONTENTS_DIR/Info.plist" 2>/dev/null || true

# 2. Siapkan Icon .icns untuk macOS
ICONSET_DIR="$DIST_DIR/icon.iconset"
mkdir -p "$ICONSET_DIR"
ASSETS_DIR="$PROJECT_ROOT/gui/assets"

if command -v iconutil >/dev/null 2>&1; then
  echo "[*] Membuat format icon Apple .icns menggunakan iconutil..."
  [ -f "$ASSETS_DIR/icon_16.png" ] && cp "$ASSETS_DIR/icon_16.png" "$ICONSET_DIR/icon_16x16.png"
  [ -f "$ASSETS_DIR/icon_32.png" ] && cp "$ASSETS_DIR/icon_32.png" "$ICONSET_DIR/icon_16x16@2x.png"
  [ -f "$ASSETS_DIR/icon_32.png" ] && cp "$ASSETS_DIR/icon_32.png" "$ICONSET_DIR/icon_32x32.png"
  [ -f "$ASSETS_DIR/icon_64.png" ] && cp "$ASSETS_DIR/icon_64.png" "$ICONSET_DIR/icon_32x32@2x.png"
  [ -f "$ASSETS_DIR/icon_128.png" ] && cp "$ASSETS_DIR/icon_128.png" "$ICONSET_DIR/icon_128x128.png"
  [ -f "$ASSETS_DIR/icon_256.png" ] && cp "$ASSETS_DIR/icon_256.png" "$ICONSET_DIR/icon_128x128@2x.png"
  [ -f "$ASSETS_DIR/icon_256.png" ] && cp "$ASSETS_DIR/icon_256.png" "$ICONSET_DIR/icon_256x256.png"
  [ -f "$ASSETS_DIR/icon_512.png" ] && cp "$ASSETS_DIR/icon_512.png" "$ICONSET_DIR/icon_256x256@2x.png"
  [ -f "$ASSETS_DIR/icon_512.png" ] && cp "$ASSETS_DIR/icon_512.png" "$ICONSET_DIR/icon_512x512.png"
  iconutil -c icns "$ICONSET_DIR" -o "$RESOURCES_DIR/icon.icns"
  rm -rf "$ICONSET_DIR"
  echo "[+] icon.icns berhasil dibuat."
else
  # Fallback: jika iconutil tidak ada (misal di Linux), salin icon PNG
  [ -f "$ASSETS_DIR/icon_512.png" ] && cp "$ASSETS_DIR/icon_512.png" "$RESOURCES_DIR/icon.png"
fi

# 3. Cari atau Build Executable macOS
GUI_BIN=""
for b in "$DIST_DIR/snake-go-gui" \
         "$DIST_DIR/ular-gui" \
         "$PROJECT_ROOT/snake-go-gui" \
         "$PROJECT_ROOT/ular-gui" \
         "$DIST_DIR/snake-go-gui-darwin-arm64" \
         "$DIST_DIR/snake-go-gui-darwin-amd64" \
         "$PROJECT_ROOT/gui/snake-go-gui"; do
  if [ -f "$b" ]; then
    GUI_BIN="$b"
    break
  fi
done

if [ -n "$GUI_BIN" ]; then
  cp "$GUI_BIN" "$MACOS_DIR/snake-go-gui"
  chmod +x "$MACOS_DIR/snake-go-gui"
  echo "[+] Executable disalin ke $MACOS_DIR/snake-go-gui"
else
  echo "[!] Peringatan: Binary GUI macOS belum ada. Silakan build dengan 'cd gui && go build -o snake-go-gui .'"
fi

# 4. Ad-hoc Codesign (Menghindari blokir Gatekeeper macOS)
if command -v codesign >/dev/null 2>&1; then
  echo "[*] Melakukan ad-hoc codesigning bundle..."
  codesign -s - --force --deep "$APP_DIR" || true
  echo "[+] Codesign selesai."
fi

# 5. Buat File Disk Image (.dmg) jika hdiutil tersedia
if command -v hdiutil >/dev/null 2>&1; then
  DMG_NAME="Snake-Go-${VERSION}-macOS.dmg"
  DMG_PATH="$DIST_DIR/$DMG_NAME"
  DMG_TMP="$DIST_DIR/dmg_tmp"

  echo "[*] Mengemas file DMG..."
  rm -rf "$DMG_TMP" "$DMG_PATH"
  mkdir -p "$DMG_TMP"
  cp -R "$APP_DIR" "$DMG_TMP/"
  ln -s /Applications "$DMG_TMP/Applications"

  hdiutil create -volname "Snake Go" -srcfolder "$DMG_TMP" -ov -format UDZO "$DMG_PATH"
  rm -rf "$DMG_TMP"
  echo "[+] DMG siap didistribusikan: $DMG_PATH"
fi

echo "=============================================="
echo "   Selesai: $APP_DIR siap dijalankan!         "
echo "=============================================="
