#!/bin/sh
# xrxs CLI 一键安装脚本 (含 AI Agent Skill)
# 适用于 macOS / Linux
#
# 用法:
#   curl -fsSL https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh
#
# 环境变量:
#   XRXS_VERSION     — 指定版本 (默认 latest)
#   XRXS_INSTALL_DIR — 安装目录 (默认 /usr/local/bin, 不可写时尝试 sudo, 最终回退到 ~/.local/bin)
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

install_binary() {
  src="$1"

  # 1. 用户指定了安装目录
  if [ -n "${XRXS_INSTALL_DIR:-}" ]; then
    mkdir -p "$XRXS_INSTALL_DIR" 2>/dev/null || true
    if [ -w "$XRXS_INSTALL_DIR" ]; then
      cp "$src" "$XRXS_INSTALL_DIR/$BIN_NAME"
      chmod +x "$XRXS_INSTALL_DIR/$BIN_NAME"
      echo "$XRXS_INSTALL_DIR"
      return
    fi
    err "安装目录 $XRXS_INSTALL_DIR 不可写"
  fi

  # 2. 优先 /usr/local/bin (macOS/Linux 系统默认 PATH)
  mkdir -p /usr/local/bin 2>/dev/null || true
  if [ -w /usr/local/bin ]; then
    cp "$src" /usr/local/bin/$BIN_NAME
    chmod +x /usr/local/bin/$BIN_NAME
    echo "/usr/local/bin"
    return
  fi

  # 3. 不可写时尝试 sudo
  if command -v sudo >/dev/null 2>&1; then
    sudo cp "$src" /usr/local/bin/$BIN_NAME
    sudo chmod +x /usr/local/bin/$BIN_NAME
    echo "/usr/local/bin"
    return
  fi

  # 4. 最终回退: ~/.local/bin
  mkdir -p "$HOME/.local/bin" 2>/dev/null || true
  cp "$src" "$HOME/.local/bin/$BIN_NAME"
  chmod +x "$HOME/.local/bin/$BIN_NAME"

  case ":$PATH:" in
    *:"$HOME/.local/bin":*) ;;
    *)
      for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.profile"; do
        if [ -w "$rc" ] 2>/dev/null || [ ! -f "$rc" ]; then
          grep -q "$HOME/.local/bin" "$rc" 2>/dev/null && continue
          echo "export PATH=\"\$PATH:\$HOME/.local/bin\"" >> "$rc" 2>/dev/null || true
        fi
      done
      ;;
  esac

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


echo ""
echo "  ╔══════════════════════════════════════╗"
echo "  ║   欢迎使用薪人薪事 CLI                ║"
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

VER="${VERSION#v}"
ARCHIVE="xrxs_${VER}_${PLATFORM}.tar.gz"
DOWNLOAD_URL="${BASE_URL}/download/${VERSION}/${ARCHIVE}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

say "正在下载 xrxs ${VERSION} ..."
download "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE"

say "正在安装..."
tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"
INSTALL_DIR="$(install_binary "$TMP_DIR/$BIN_NAME")"

say "已安装: $INSTALL_DIR/$BIN_NAME"

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
say "  登录:"
say "    xrxs auth login --base-url https://s122.devtest.vip"
say ""
say "  登录后在 Claude Code 中输入 /xrxs 即可通过对话操作审批。"
say ""
