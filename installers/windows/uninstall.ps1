# Uninstaller PowerShell untuk Snake Go di Windows
#
# Penggunaan:
#   .\uninstall.ps1

param (
    [switch]$KeepSaveData = $true
)

$ErrorActionPreference = "Continue"

Write-Host "==============================================" -ForegroundColor Yellow
Write-Host "        Pencopotan (Uninstall) Snake Go       " -ForegroundColor Yellow
Write-Host "==============================================" -ForegroundColor Yellow

# 1. Hapus Shortcut Start Menu
$StartMenuDir = Join-Path $env:APPDATA "Microsoft\Windows\Start Menu\Programs\SnakeGo"
if (Test-Path $StartMenuDir) {
    Remove-Item -Path $StartMenuDir -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "[+] Shortcut Start Menu dihapus." -ForegroundColor Gray
}

# 2. Hapus Shortcut Desktop
$DesktopDir = [Environment]::GetFolderPath('Desktop')
$DesktopShortcut = Join-Path $DesktopDir "Snake Go.lnk"
if (Test-Path $DesktopShortcut) {
    Remove-Item -Path $DesktopShortcut -Force -ErrorAction SilentlyContinue
    Write-Host "[+] Shortcut Desktop dihapus." -ForegroundColor Gray
}

# 3. Hapus Registrasi Windows Settings (Uninstall Key)
$RegPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\SnakeGo"
if (Test-Path $RegPath) {
    Remove-Item -Path $RegPath -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "[+] Registrasi Windows Settings dihapus." -ForegroundColor Gray
}

# 4. Hapus Berkas Program
$InstallDir = Split-Path -Parent $MyInvocation.MyCommand.Path
if (-not $InstallDir -or $InstallDir -eq "" -or (-not (Test-Path (Join-Path $InstallDir "snake-go-gui.exe")))) {
    $InstallDir = "$env:LOCALAPPDATA\Programs\SnakeGo"
}

if (Test-Path $InstallDir) {
    # Hapus file selain uninstaller itu sendiri jika sedang berjalan
    Get-ChildItem -Path $InstallDir -Exclude "uninstall.ps1" | Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "[+] Berkas instalasi dihapus dari $InstallDir." -ForegroundColor Gray

    # Self-delete helper via cmd untuk menghapus folder instalasi setelah script selesai
    Start-Process -FilePath "cmd.exe" -ArgumentList "/c timeout /t 1 /nobreak >nul & rd /s /q `"$InstallDir`"" -WindowStyle Hidden
}

Write-Host ""
Write-Host "[*] Catatan: Data rekor skor (%APPDATA%\SnakeGo\scores.json) tetap disimpan secara aman." -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Green
Write-Host "     Snake Go telah berhasil dicopot. Bye!    " -ForegroundColor Green
Write-Host "==============================================" -ForegroundColor Green
