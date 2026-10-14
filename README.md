# Snake (Golang)

Game snake dalam dua versi: terminal dan desktop.

## Versi terminal (tanpa dependensi)

```bash
go run .
# atau
go build -o snake .
./snake
```

## Versi desktop — Ular-ularan (window GUI)

Berbasis [Ebitengine](https://ebitengine.org/), window asli 632×728 dengan area main persegi **25×25** (rasio 1:1).

Tampilan modern: tema gelap gradien, HUD pil (skor / terbaik / panjang / level), ular berkarakter (mata mengikuti arah, lidah berkedip, badan gradien + gerak halus terinterpolasi), apel + bonus bintang emas (+50, berwaktu), ledakan partikel, teks skor melayang, screen shake + flash merah saat mati, panel start / pause / game over, dan penanda rekor baru.

```bash
cd gui
go run .
# atau
go build -o snake-gui .
./snake-gui
```

Butuh compiler C + lib X11/OpenGL (di Arch/CachyOS biasanya sudah ada: `libx11 libxrandr libxi libxcursor libxinerama mesa base-devel`).

## Build untuk Linux, Windows, macOS

```bash
./build-all.sh   # hasil di dist/
```

Versi **terminal** murni stdlib sehingga cross-compile langsung bisa dari Linux (sudah diverifikasi: ELF, PE32+ `.exe`, Mach-O amd64+arm64). Catatan jujur: mode raw terminal memakai `stty` — didukung di Linux (`-F`) dan macOS (`-f`, otomatis dipilih), tapi **tidak di Windows** (tidak ada `stty`), jadi binary Windows-nya jalan tapi kontrolnya tidak responsif. Perlu `golang.org/x/term` agar Windows penuh — bisa jadi follow-up.

Versi **GUI** butuh CGO: Linux build langsung bisa; Windows butuh `mingw-w64-gcc`; macOS paling gampang build di Mac langsung (`cd gui && go build .`).

Alternatif tanpa Mac fisik: push repo ke GitHub — workflow di `.github/workflows/build.yml` otomatis membuild semuanya di runner Linux/Windows/macOS, lalu binary per OS bisa diunduh dari tab Actions.

## Kontrol (kedua versi)

| Tombol | Aksi |
|---|---|
| `WASD` / Panah | Gerak (di desktop juga untuk mulai) |
| `P` / Spasi | Pause / lanjut (di desktop Spasi juga untuk mulai) |
| `R` / Enter | Restart setelah game over (`R` desktop restart kapan saja) |
| `M` | Ganti mode Tembok / Wrap |
| `Q` / `Esc` / `Ctrl-C` | Keluar |

## Aturan

- Terminal: makan `★` = +10 poin. Desktop: makan apel = +10, bintang emas bonus = +50 (muncul tiap 5 apel, kedaluwarsa ~8 detik).
- Tiap makan: ular memanjang dan kecepatan naik (ada indikator LEVEL di HUD desktop).
- Menabrak dinding atau badan sendiri = game over.
- Putar balik 180° ditolak otomatis.
- Rintangan acak muncul di papan; menabraknya = game over.
- Mode **Wrap** (tekan `M`): ular tembus dinding, muncul di sisi berlawanan (tanpa game over dinding).
- Skor terbaik disimpan permanen ke file (`~/.ular_best`).
- Versi desktop memiliki efek suara sederhana (beep).

## File

- `main.go` — seluruh game terminal (loop, input raw via `stty`, render ANSI).
- `main_test.go` — test logika: jalan, makan, tabrakan, anti-putar-balik (`go test -v`).
- `gui/main.go` — versi desktop modern (modul terpisah `banguncandi-gui`, tanpa mengubah versi terminal).

## Catatan terminal

Game memakai mode `cbreak` via `stty -F /dev/tty` dan kode ANSI. Jalankan di terminal sungguhan (bukan output IDE yang non-TTY). State terminal dikembalikan otomatis saat keluar.
