#!/bin/sh
# xrxs CLI 一键安装脚本 (含 AI Agent Skill)
# 适用于 macOS / Linux
#
# 用法:
#   curl -fsSL https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh
#
# 环境变量:
#   XRXS_VERSION     — 指定版本 (默认 latest)
#   XRXS_INSTALL_DIR — 安装目录 (默认 ~/.local/bin)
#   XRXS_NO_SKILLS   — 设为 1 跳过 Skill 安装

set -eu

REPO="LucyHeres/xrxs-cli"
BIN_NAME="xrxs"
DEFAULT_DIR="$HOME/.local/bin"
INSTALL_DIR="${XRXS_INSTALL_DIR:-$DEFAULT_DIR}"
VERSION="${XRXS_VERSION:-latest}"
BASE_URL="${XRXS_BASE_URL:-https://github.com/${REPO}/releases}"
NO_SKILLS="${XRXS_NO_SKILLS:-0}"

say()  { printf '  %s\n' "$@"; }
err()  { printf '  ❌ %s\n' "$@" >&2; exit 1; }
warn() { printf '  ⚠️  %s\n' "$@" >&2; }
need_cmd() { command -v "$1" >/dev/null 2>&1; }

download() {
  url="$1"; dest="$2"
  if need_cmd curl; then
    curl -fsSL ${XRXS_INSECURE:+-k} "$url" -o "$dest"
  elif need_cmd wget; then
    wget -qO "$dest" "$url"
  else
    err "请安装 curl 或 wget 后重试"
  fi
}

detect_platform() {
  arch="$(uname -m)"
  case "$(uname -s)" in
    Linux)  os="linux" ;;
    Darwin) os="darwin" ;;
    *)      err "不支持的操作系统: $(uname -s)" ;;
  esac
  case "$arch" in
    x86_64|amd64)  arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *)             err "不支持的 CPU 架构: $arch" ;;
  esac
  echo "${os}-${arch}"
}

echo ""
echo "  ╔══════════════════════════════════════╗"
echo "  ║   薪人薪事 CLI (xrxs) 安装程序       ║"
echo "  ╚══════════════════════════════════════╝"
echo ""

PLATFORM="$(detect_platform)"
say "检测到系统: $PLATFORM"

if [ "$VERSION" = "latest" ]; then
  say "获取最新版本..."
  if need_cmd curl; then
    VERSION=$(curl -fsSL ${XRXS_INSECURE:+-k} "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": "\(.*\)".*/\1/')
  fi
  if [ -z "$VERSION" ]; then
    VERSION="dev"
    warn "无法获取最新版本号"
  fi
fi
say "版本: $VERSION"

ARCHIVE="xrxs_${VERSION}_${PLATFORM}.tar.gz"
DOWNLOAD_URL="${BASE_URL}/download/${VERSION}/${ARCHIVE}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

say "下载 $DOWNLOAD_URL ..."
download "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE"

say "解压..."
tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"

mkdir -p "$INSTALL_DIR"
cp "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"
say "已安装到: $INSTALL_DIR/$BIN_NAME"

case ":$PATH:" in
  *:"$INSTALL_DIR":*) ;;
  *)
    warn "$INSTALL_DIR 不在 PATH 中"
    SHELL_RC=""
    case "$SHELL" in
      */zsh)  SHELL_RC="$HOME/.zshrc" ;;
      */bash) SHELL_RC="$HOME/.bashrc" ;;
      */fish) SHELL_RC="$HOME/.config/fish/config.fish" ;;
    esac
    if [ -n "$SHELL_RC" ]; then
      echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
      say "已添加 PATH 到 $SHELL_RC (运行 source $SHELL_RC 生效)"
    fi
    ;;
esac

# 安装 AI Agent Skills（覆盖所有主流 AI 编辑器）
if [ "$NO_SKILLS" != "1" ] && [ -f "$TMP_DIR/skills/xrxs/SKILL.md" ]; then
  echo ""
  say "安装 AI Agent Skills..."

  SKILL_SRC="$TMP_DIR/skills/xrxs/SKILL.md"
  INSTALLED=0

  for agent_dir in \
    ".agents/skills" \
    ".claude/skills" \
    ".cursor/skills" \
    ".gemini/skills" \
    ".codex/skills" \
    ".github/skills" \
    ".windsurf/skills" \
    ".augment/skills" \
    ".cline/skills" \
    ".amp/skills" \
    ".kiro/skills" \
    ".trae/skills" \
    ".openclaw/skills" \
    ".hermes/skills" \
    ".qoder/skills" \
    ".opencode/skills"
  do
    dest="$HOME/$agent_dir/xrxs/SKILL.md"
    mkdir -p "$(dirname "$dest")"
    if cp "$SKILL_SRC" "$dest" 2>/dev/null; then
      say "  ✅ ~/$agent_dir/xrxs"
      INSTALLED=$((INSTALLED + 1))
    fi
  done

  if [ "$INSTALLED" -gt 0 ]; then
    say "Skills 已安装到 $INSTALLED 个 AI 编辑器 (Claude Code / Cursor / Trae / Codex / Windsurf ...)"
  fi
fi

echo ""
say "安装完成!"
say ""
say "  xrxs auth login --base-url https://your-company.example.com"
say "  xrxs approval list search --status 0"
say ""
