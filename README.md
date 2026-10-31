# Snake Go

Game modern Snake (Golang) dalam dua versi: **desktop GUI** (berbasis Ebitengine) dan **terminal**. Dilengkapi sistem penyimpanan skor terstruktur (Top 10 Leaderboard), icon aplikasi kustom multi-platform, sistem versioning SemVer, dan script instalasi otomatis untuk Windows, Linux, dan macOS.

---

## 🚀 Menjalankan Cepat

### 1. Versi Desktop (GUI Window)
```bash
cd gui
go run .
# atau build binary
go build -o snake-go-gui .
./snake-go-gui
```
> Window 632×728, area permainan persegi 25×25, tema gelap cyberpunk/arcade, efek partikel, animasi lidah & mata ular, bonus bintang emas (+50 poin), dan suara sintetis (beep).

### 2. Versi Terminal
```bash
go run .
# Cek info versi
go run . --version
```

---

## 🎮 Kontrol Permainan

| Tombol | Aksi |
|---|---|
| `WASD` / Panah | Gerak (di desktop juga untuk mulai permainan) |
| `P` / Spasi | Pause / lanjut |
| `L` | **Buka / Tutup Papan Skor (Top 10 Leaderboard)** |
| `R` / Enter | Restart setelah game over (`R` di desktop bisa restart kapan saja) |
| `M` | Ganti mode **Tembok** (Dinding) / **Wrap** (Tembus) |
| `Q` / `Esc` / `Ctrl-C` | Keluar |

---

## 🏆 Fitur Penyimpanan Skor & Statistik

Data permainan disimpan secara otomatis dan persisten dalam format JSON terstruktur dengan pengamanan **Atomic Save** (`.tmp` -> rename) untuk mencegah korupsi data:
- **Top 10 Leaderboard**: Mencatat ranking, skor, panjang ular, level kecepatan, mode (Tembok/Wrap), dan tanggal pencapaian.
- **Statistik Akumulatif**: Menghitung total permainan yang dimainkan, total apel dimakan, dan total bintang emas yang didapat.
- **Migrasi Otomatis**: Jika file lama `~/.ular_best` ditemukan, skor terbaik otomatis diimpor ke format baru.

### Lokasi Berkas Penyimpanan per OS:
- **Windows**: `%APPDATA%\SnakeGo\scores.json`
- **Linux**: `~/.config/snake-go/scores.json` (XDG standard)
- **macOS**: `~/Library/Application Support/SnakeGo/scores.json`

---

## 📦 Sistem Instalasi di 3 Platform

Folder `installers/` menyediakan sistem instalasi native untuk masing-masing sistem operasi:

### 1. Windows (PowerShell)
Instalasi per-user tanpa memerlukan hak administrator:
```powershell
# Jalankan dari folder project di PowerShell
.\installers\windows\install.ps1

# Opsi: tambahkan ke PATH pengguna
.\installers\windows\install.ps1 -AddToPath
```
- Terpasang ke `%LOCALAPPDATA%\Programs\SnakeGo\`.
- Otomatis membuat shortcut **Start Menu** dan **Desktop** dengan nama `Snake Go.lnk` dan icon kustom `icon.ico`.
- Terdaftar resmi di menu **Windows Settings / Installed Apps** (bisa dicopot lewat menu Settings atau `uninstall.ps1`).
- Tersedia juga skrip Inno Setup `installers/windows/snake-go.iss` untuk membuat installer wizard `.exe`.

### 2. Linux (Bash)
```bash
# Mode pengguna (ke ~/.local/bin, tanpa butuh sudo)
./installers/linux/install.sh

# Mode sistem (ke /usr/local/bin, butuh sudo)
sudo ./installers/linux/install.sh --system
```
- Memasang binary `snake-go-gui` dan `snake-go`.
- Mendaftarkan desktop entry `snake-go.desktop` ke launcher aplikasi (GNOME, KDE Plasma, XFCE).
- Memasang icon FreeDesktop hicolor resolusi 16, 48, 256, dan SVG scalable.
- Untuk mencopot: `./installers/linux/uninstall.sh`.

### 3. macOS (App Bundle & DMG)
```bash
# Buat bundle Snake-Go.app & disk image .dmg
./installers/macos/package_macos.sh

# Pasang ke /Applications
./installers/macos/install.sh
```
- Membentuk bundle native `Snake-Go.app` dengan metadata `Info.plist` dan icon Apple `.icns`.
- Melakukan ad-hoc codesigning agar Gatekeeper mengizinkan aplikasi dibuka.
- Menghasilkan file DMG siap bagi pengguna yang ingin drag-and-drop ke `/Applications`.

---

## 🏷️ Versioning & Otomasi GitHub Releases

Versi aplikasi dikelola terpusat pada file `VERSION` (SemVer). Saat build dijalankan, informasi versi, git commit hash, dan tanggal kompilasi diinjeksikan otomatis via `-ldflags`:

```bash
./build-all.sh   # Hasil kompilasi seluruh platform di folder dist/
```

### Publikasi Otomatis ke GitHub Releases
Workflow GitHub Actions (`.github/workflows/build.yml`) otomatis membangun seluruh platform dan mempublikasikan binary ke halaman **Releases** repositori Anda:
1. **Otomatis via Git Tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
   GitHub Actions akan otomatis membangun semua binary dan membuat rilis baru di tab **Releases** lengkap dengan file binary `.exe`, `.dmg`, Linux binary, dan paket instalasi zip.
2. **Manual via Tab Actions**:
   Buka tab **Actions** di GitHub → pilih workflow **build** → klik tombol **Run workflow** → centang opsi *"Publikasikan ke GitHub Releases?"*.

---

## 📁 Struktur Berkas

- `VERSION` — Berkas versi semantik terpusat (`1.0.0`).
- `pkg/version/` — Modul informasi versi, git commit, dan tanggal build.
- `pkg/scores/` — Modul penyimpanan skor JSON, atomic write, dan migrasi.
- `main.go` — Game versi terminal (loop raw via `stty`, render ANSI, flag `--version`).
- `gui/main.go` — Game versi desktop modern Ebitengine dengan Leaderboard UI.
- `gui/assets/` — Aset grafis icon (SVG, PNG 16–512, ICO multi-layer).
- `gui/snake-go.rc` — Windows PE resource definition (Icon + `VERSIONINFO`).
- `snake-go.desktop` — Desktop Entry untuk launcher menu Linux.
- `installers/` — Skrip instalasi lengkap untuk Windows, Linux, dan macOS.
- `tools/generate_icons.py` — Script generator icon resolusi penuh.
