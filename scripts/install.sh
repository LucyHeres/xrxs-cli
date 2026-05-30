#!/bin/sh
# xrxs CLI 一键安装脚本 (含 AI Agent Skill)
# 适用于 macOS / Linux
#
# 用法:
#   curl -fsSL https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh
#
# 环境变量:
#   XRXS_VERSION     — 指定版本 (默认 latest)
#   XRXS_INSTALL_DIR — 安装目录 (默认 /usr/local/bin, 不可写时回退到 ~/.local/bin)
#   XRXS_NO_SKILLS   — 设为 1 跳过 Skill 安装

set -eu

REPO="LucyHeres/xrxs-cli"
BIN_NAME="xrxs"
VERSION="${XRXS_VERSION:-latest}"
BASE_URL="${XRXS_BASE_URL:-https://github.com/${REPO}/releases}"
NO_SKILLS="${XRXS_NO_SKILLS:-0}"

say()  { printf '  %s\n' "$@"; }
err()  { printf '  \033[31m%s\033[0m\n' "$@" >&2; exit 1; }
need_cmd() { command -v "$1" >/dev/null 2>&1; }

# 选择安装目录：优先 /usr/local/bin (macOS/Linux 默认 PATH 内)
pick_install_dir() {
  if [ -n "${XRXS_INSTALL_DIR:-}" ]; then
    mkdir -p "$XRXS_INSTALL_DIR" 2>/dev/null || true
    if [ -w "$XRXS_INSTALL_DIR" ]; then
      echo "$XRXS_INSTALL_DIR"
      return
    fi
  fi

  # 尝试 /usr/local/bin
  if [ ! -d /usr/local/bin ]; then
    mkdir -p /usr/local/bin 2>/dev/null || true
  fi
  if [ -w /usr/local/bin ]; then
    echo "/usr/local/bin"
    return
  fi

  # 回退到 ~/.local/bin
  mkdir -p "$HOME/.local/bin" 2>/dev/null || true
  echo "$HOME/.local/bin"
}

download() {
  url="$1"; dest="$2"
  if need_cmd curl; then
    curl -fsSL ${XRXS_INSECURE:+-k} "$url" -o "$dest"
  elif need_cmd wget; then
    wget -qO "$dest" "$url"
  else
    err "安装失败: 系统缺少 curl 或 wget"
  fi
}

detect_platform() {
  arch="$(uname -m)"
  case "$(uname -s)" in
    Linux)  os="linux" ;;
    Darwin) os="darwin" ;;
    *)      err "暂不支持当前操作系统" ;;
  esac
  case "$arch" in
    x86_64|amd64)  arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *)             err "暂不支持当前 CPU 架构" ;;
  esac
  echo "${os}-${arch}"
}

# 确保 INSTALL_DIR 在 PATH 中（仅在回退到 ~/.local/bin 时需要）
ensure_path() {
  dir="$1"
  # /usr/local/bin 已在系统 PATH，无需处理
  case "$dir" in /usr/local/bin) return ;; esac

  case ":$PATH:" in
    *:"$dir":*) return ;;
  esac

  # 写入 shell 配置文件，新终端自动生效
  for rc in "$HOME/.zshenv" "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
    if [ -f "$rc" ] || [ "$rc" = "$HOME/.zshenv" ] || [ "$rc" = "$HOME/.profile" ]; then
      if ! grep -q "$dir" "$rc" 2>/dev/null; then
        echo "export PATH=\"\$PATH:$dir\"" >> "$rc"
      fi
    fi
  done
}

echo ""
echo "  ╔══════════════════════════════════════╗"
echo "  ║   薪人薪事 CLI (xrxs) 安装程序       ║"
echo "  ╚══════════════════════════════════════╝"
echo ""

PLATFORM="$(detect_platform)"
say "系统: $PLATFORM"

if [ "$VERSION" = "latest" ]; then
  if need_cmd curl; then
    VERSION=$(curl -fsSL ${XRXS_INSECURE:+-k} "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": "\(.*\)".*/\1/')
  fi
  if [ -z "$VERSION" ]; then
    err "获取版本信息失败，请检查网络连接"
  fi
fi

INSTALL_DIR="$(pick_install_dir)"
# 去掉 VERSION 开头的 v (API 返回 v0.1.2, 文件名是 0.1.2)
VER="${VERSION#v}"
ARCHIVE="xrxs_${VER}_${PLATFORM}.tar.gz"
DOWNLOAD_URL="${BASE_URL}/download/${VERSION}/${ARCHIVE}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

say "正在下载 xrxs ${VERSION} ..."
download "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE"

say "正在安装..."
tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"
cp "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

ensure_path "$INSTALL_DIR"

say "已安装: $INSTALL_DIR/$BIN_NAME"

# 安装 AI Agent Skills
if [ "$NO_SKILLS" != "1" ] && [ -f "$TMP_DIR/skills/xrxs/SKILL.md" ]; then
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
    mkdir -p "$(dirname "$dest")" 2>/dev/null || continue
    cp "$SKILL_SRC" "$dest" 2>/dev/null || continue
    INSTALLED=$((INSTALLED + 1))
  done
fi

echo ""
say "安装完成！"
say ""
say "  使用方法:"
say "    xrxs auth login --base-url https://your-company.example.com"
say "    xrxs approval list search --status 0"
say ""
