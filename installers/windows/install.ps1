# Installer PowerShell untuk Snake Go di Windows
# Dapat dijalankan tanpa hak administrator (instalasi per-user)
#
# Penggunaan:
#   .\install.ps1
#   .\install.ps1 -AddToPath
#   .\install.ps1 -InstallDir "D:\Games\SnakeGo"

param (
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\SnakeGo",
    [switch]$AddToPath = $false,
    [switch]$NoDesktopShortcut = $false
)

$ErrorActionPreference = "Stop"

Write-Host "==============================================" -ForegroundColor Cyan
Write-Host "       Instalasi Game Snake Go (Windows)      " -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Cyan

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Resolve-Path (Join-Path $ScriptDir "..\..")

# 1. Baca Versi
$VersionFile = Join-Path $ProjectRoot "VERSION"
$Version = "1.0.0"
if (Test-Path $VersionFile) {
    $Version = (Get-Content $VersionFile -Raw).Trim()
}
Write-Host "[*] Versi aplikasi: $Version" -ForegroundColor Gray

# 2. Cari Binary Executable
$GuiSource = $null
$PossibleGui = @(
    (Join-Path $ProjectRoot "dist\snake-go-gui-windows-amd64.exe"),
    (Join-Path $ProjectRoot "dist\ular-gui-windows-amd64.exe"),
    (Join-Path $ProjectRoot "gui\snake-go-gui.exe"),
    (Join-Path $ProjectRoot "gui\snake-gui.exe"),
    (Join-Path $ProjectRoot "snake-go-gui.exe")
)
foreach ($p in $PossibleGui) {
    if (Test-Path $p) {
        $GuiSource = $p
        break
    }
}

$TermSource = $null
$PossibleTerm = @(
    (Join-Path $ProjectRoot "dist\snake-go-terminal-windows-amd64.exe"),
    (Join-Path $ProjectRoot "dist\ular-terminal-windows-amd64.exe"),
    (Join-Path $ProjectRoot "snake-go-terminal.exe"),
    (Join-Path $ProjectRoot "snake-go.exe"),
    (Join-Path $ProjectRoot "snake.exe")
)
foreach ($p in $PossibleTerm) {
    if (Test-Path $p) {
        $TermSource = $p
        break
    }
}

$IconSource = Join-Path $ProjectRoot "gui\icon.ico"

# 3. Buat Folder Instalasi
Write-Host "[*] Menyiapkan direktori instalasi: $InstallDir" -ForegroundColor Gray
if (-not (Test-Path $InstallDir)) {
    New-Item -Path $InstallDir -ItemType Directory -Force | Out-Null
}

# 4. Salin Berkas
if ($GuiSource -and (Test-Path $GuiSource)) {
    Copy-Item -Path $GuiSource -Destination (Join-Path $InstallDir "snake-go-gui.exe") -Force
    Unblock-File -Path (Join-Path $InstallDir "snake-go-gui.exe") -ErrorAction SilentlyContinue
    Write-Host "[+] Berhasil memasang snake-go-gui.exe" -ForegroundColor Green
} else {
    Write-Host "[!] Catatan: Binary snake-go-gui.exe belum di-build. Anda bisa mengompilasinya nanti." -ForegroundColor Yellow
}

if ($TermSource -and (Test-Path $TermSource)) {
    Copy-Item -Path $TermSource -Destination (Join-Path $InstallDir "snake-go.exe") -Force
    Unblock-File -Path (Join-Path $InstallDir "snake-go.exe") -ErrorAction SilentlyContinue
    Write-Host "[+] Berhasil memasang snake-go.exe (Terminal)" -ForegroundColor Green
}

if (Test-Path $IconSource) {
    Copy-Item -Path $IconSource -Destination (Join-Path $InstallDir "icon.ico") -Force
    Write-Host "[+] Berhasil memasang icon.ico" -ForegroundColor Green
}

# Salin uninstaller
$UninstallScript = Join-Path $ScriptDir "uninstall.ps1"
if (Test-Path $UninstallScript) {
    Copy-Item -Path $UninstallScript -Destination (Join-Path $InstallDir "uninstall.ps1") -Force
}

# Simpan metadata versi
Set-Content -Path (Join-Path $InstallDir "VERSION") -Value $Version -Force

# 5. Buat Shortcut Windows (.lnk)
$WshShell = New-Object -ComObject WScript.Shell
$TargetExe = Join-Path $InstallDir "snake-go-gui.exe"
$InstalledIcon = Join-Path $InstallDir "icon.ico"

# A. Start Menu Shortcut
$StartMenuDir = Join-Path $env:APPDATA "Microsoft\Windows\Start Menu\Programs\SnakeGo"
if (-not (Test-Path $StartMenuDir)) {
    New-Item -Path $StartMenuDir -ItemType Directory -Force | Out-Null
}
$StartShortcut = $WshShell.CreateShortcut((Join-Path $StartMenuDir "Snake Go.lnk"))
$StartShortcut.TargetPath = $TargetExe
$StartShortcut.WorkingDirectory = $InstallDir
$StartShortcut.Description = "Game Snake Go v$Version"
if (Test-Path $InstalledIcon) {
    $StartShortcut.IconLocation = "$InstalledIcon,0"
}
$StartShortcut.Save()
Write-Host "[+] Shortcut Start Menu berhasil dibuat: $StartMenuDir" -ForegroundColor Green

# B. Desktop Shortcut
if (-not $NoDesktopShortcut) {
    $DesktopDir = [Environment]::GetFolderPath('Desktop')
    $DesktopShortcut = $WshShell.CreateShortcut((Join-Path $DesktopDir "Snake Go.lnk"))
    $DesktopShortcut.TargetPath = $TargetExe
    $DesktopShortcut.WorkingDirectory = $InstallDir
    $DesktopShortcut.Description = "Game Snake Go v$Version"
    if (Test-Path $InstalledIcon) {
        $DesktopShortcut.IconLocation = "$InstalledIcon,0"
    }
    $DesktopShortcut.Save()
    Write-Host "[+] Shortcut Desktop berhasil dibuat." -ForegroundColor Green
}

# 6. Registrasi ke Windows Settings / Installed Apps (Add/Remove Programs)
$RegPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\SnakeGo"
if (-not (Test-Path $RegPath)) {
    New-Item -Path $RegPath -Force | Out-Null
}
Set-ItemProperty -Path $RegPath -Name "DisplayName" -Value "Snake Go"
Set-ItemProperty -Path $RegPath -Name "DisplayVersion" -Value $Version
Set-ItemProperty -Path $RegPath -Name "Publisher" -Value "Antigravity Games"
Set-ItemProperty -Path $RegPath -Name "InstallLocation" -Value $InstallDir
if (Test-Path $InstalledIcon) {
    Set-ItemProperty -Path $RegPath -Name "DisplayIcon" -Value $InstalledIcon
}
$UninstallCmd = "powershell.exe -ExecutionPolicy Bypass -NoProfile -File `"$InstallDir\uninstall.ps1`""
Set-ItemProperty -Path $RegPath -Name "UninstallString" -Value $UninstallCmd
Write-Host "[+] Terdaftar resmi di Windows Settings (Add/Remove Programs)." -ForegroundColor Green

# 7. Tambah ke User PATH (Opsional)
if ($AddToPath) {
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
        Write-Host "[+] $InstallDir telah ditambahkan ke User PATH." -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "==============================================" -ForegroundColor Green
Write-Host "       Instalasi Snake Go v$Version Selesai!   " -ForegroundColor Green
Write-Host "   Buka dari Start Menu atau Desktop. Enjoy!  " -ForegroundColor Green
Write-Host "==============================================" -ForegroundColor Green
