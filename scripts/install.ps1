# xrxs CLI 一键安装脚本 (Windows PowerShell)
#
# 用法:
#   irm https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.ps1 | iex
#
# 环境变量:
#   $env:XRXS_VERSION     — 指定版本
#   $env:XRXS_INSTALL_DIR — 安装目录 (默认 $env:LOCALAPPDATA\xrxs)

param()

$ErrorActionPreference = "Stop"
$Repo = "LucyHeres/xrxs-cli"
$BinName = "xrxs.exe"
$DefaultDir = Join-Path $env:LOCALAPPDATA "xrxs"
$InstallDir = if ($env:XRXS_INSTALL_DIR) { $env:XRXS_INSTALL_DIR } else { $DefaultDir }
$Version = if ($env:XRXS_VERSION) { $env:XRXS_VERSION } else { "latest" }
$BaseUrl = "https://github.com/$Repo/releases"

$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$Platform = "windows-$Arch"

if ($Version -eq "latest") {
    try {
        $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $Version = $Release.tag_name
    } catch {
        $Version = "dev"
    }
}

$Ver = $Version -replace '^v',''
$Archive = "xrxs_${Ver}_${Platform}.zip"
$DownloadUrl = "$BaseUrl/download/$Version/$Archive"
$TempDir = Join-Path $env:TEMP "xrxs-install"
New-Item -ItemType Directory -Force -Path $TempDir | Out-Null
$ZipPath = Join-Path $TempDir $Archive

Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath
Expand-Archive -Path $ZipPath -DestinationPath $TempDir -Force

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$ExePath = Join-Path $InstallDir $BinName
Copy-Item -Path (Join-Path $TempDir $BinName) -Destination $ExePath -Force

$CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
}

Remove-Item -Recurse -Force $TempDir

Write-Host ""
Write-Host "  薪人薪事CLI 安装完成！安装位置：$ExePath"
Write-Host ""
Write-Host "  登录："
Write-Host "    xrxs auth login --base-url https://s122.devtest.vip"
