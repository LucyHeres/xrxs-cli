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

Write-Host ""
Write-Host "  ╔══════════════════════════════════════╗"
Write-Host "  ║   欢迎使用薪人薪事 CLI                ║"
Write-Host "  ╚══════════════════════════════════════╝"
Write-Host ""

$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$Platform = "windows-$Arch"

if ($Version -eq "latest") {
    Write-Host "  获取最新版本..."
    try {
        $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $Version = $Release.tag_name
    } catch {
        $Version = "dev"
        Write-Host "  ⚠️  无法获取最新版本号，使用 dev 版本"
    }
}
Write-Host "  版本: $Version"

$Ver = $Version -replace '^v',''
$Archive = "xrxs_${Ver}_${Platform}.zip"
$DownloadUrl = "$BaseUrl/download/$Version/$Archive"
$TempDir = Join-Path $env:TEMP "xrxs-install"
New-Item -ItemType Directory -Force -Path $TempDir | Out-Null
$ZipPath = Join-Path $TempDir $Archive

Write-Host "  下载 $DownloadUrl ..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $ZipPath

Write-Host "  解压..."
Expand-Archive -Path $ZipPath -DestinationPath $TempDir -Force

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$ExePath = Join-Path $InstallDir $BinName
Copy-Item -Path (Join-Path $TempDir $BinName) -Destination $ExePath -Force
Write-Host "  已安装到: $ExePath"

$CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
    Write-Host "  已添加到用户 PATH"
}

Remove-Item -Recurse -Force $TempDir

Write-Host ""
Write-Host "  安装完成!"
Write-Host ""
Write-Host "  登录:"
Write-Host "    xrxs auth login --base-url https://s122.devtest.vip"
Write-Host ""
Write-Host "  登录后在 Claude Code 中输入 /xrxs 即可通过对话操作审批。"
Write-Host ""
