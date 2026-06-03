#!/bin/sh
# xrxs CLI 一键安装脚本 (含 AI Agent Skill)
# 适用于 macOS / Linux
#
# 用法:
#   curl -fsSL https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh
#
# 国内网络可给安装脚本与下载走镜像（示例）:
#   curl -fsSL "https://gh-proxy.org/https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh" | sh
#   XRXS_GITHUB_PROXY=https://gh-proxy.org/ curl -fsSL .../install.sh | sh
#
# 环境变量:
#   XRXS_VERSION         — 指定版本 (默认 latest)
#   XRXS_GITHUB_PROXY    — GitHub 镜像前缀，如 https://gh-proxy.org/
#   XRXS_BASE_URL        — Release 下载根地址（默认经 gh_url 处理）
#   XRXS_INSTALL_DIR     — 安装目录
#   XRXS_NO_SKILLS       — 设为 1 跳过 Skill 安装

set -eu

REPO="LucyHeres/xrxs-cli"
BIN_NAME="xrxs"
VERSION="${XRXS_VERSION:-latest}"
GITHUB_PROXY="${XRXS_GITHUB_PROXY:-}"
NO_SKILLS="${XRXS_NO_SKILLS:-0}"

# 国内访问 GitHub 较慢，默认放宽超时；可用环境变量覆盖
CURL_CONNECT_TIMEOUT="${XRXS_CURL_CONNECT_TIMEOUT:-30}"
CURL_MAX_TIME="${XRXS_CURL_MAX_TIME:-300}"
CURL="curl -fsSL --connect-timeout ${CURL_CONNECT_TIMEOUT} --max-time ${CURL_MAX_TIME} --retry 2 --retry-delay 2"

say()  { printf '  %s\n' "$@"; }
err()  { printf '  \033[31m%s\033[0m\n' "$@" >&2; exit 1; }
need_cmd() { command -v "$1" >/dev/null 2>&1; }

gh_url() {
  case "$1" in
    "${GITHUB_PROXY}"*) printf '%s\n' "$1" ;;
    *) printf '%s%s\n' "$GITHUB_PROXY" "$1" ;;
  esac
}

curl_fetch() {
  url="$1"
  dest="${2:-}"
  if [ -n "$dest" ]; then
    $CURL ${XRXS_INSECURE:+-k} "$url" -o "$dest"
  else
    $CURL ${XRXS_INSECURE:+-k} "$url"
  fi
}

resolve_latest_version() {
  # 1. 跟随 releases/latest 重定向（不依赖 api.github.com，可走镜像）
  location="$($CURL ${XRXS_INSECURE:+-k} -D - -o /dev/null \
    "$(gh_url "https://github.com/${REPO}/releases/latest")" 2>/dev/null \
    | grep -i '^location:' | tail -1 | tr -d '\r')"
  tag="$(printf '%s\n' "$location" | sed -n 's|.*/tag/\([^/]*\).*|\1|p')"
  if [ -n "$tag" ]; then
    printf '%s\n' "$tag"
    return 0
  fi

  # 2. 回退 GitHub API（仅直连，部分镜像不支持 api.github.com）
  if [ -z "$GITHUB_PROXY" ]; then
    tag="$($CURL ${XRXS_INSECURE:+-k} "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null \
      | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": "\(.*\)".*/\1/')"
    if [ -n "$tag" ]; then
      printf '%s\n' "$tag"
      return 0
    fi
  fi

  return 1
}

install_binary() {
  src="$1"

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

  mkdir -p /usr/local/bin 2>/dev/null || true
  if [ -w /usr/local/bin ]; then
    cp "$src" /usr/local/bin/$BIN_NAME
    chmod +x /usr/local/bin/$BIN_NAME
    echo "/usr/local/bin"
    return
  fi

  if command -v sudo >/dev/null 2>&1; then
    sudo cp "$src" /usr/local/bin/$BIN_NAME
    sudo chmod +x /usr/local/bin/$BIN_NAME
    echo "/usr/local/bin"
    return
  fi

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

download_archive() {
  url="$1"
  dest="$2"

  if curl_fetch "$url" "$dest" 2>/dev/null; then
    return 0
  fi

  # 直连失败且未配置镜像时，自动尝试 gh-proxy
  if [ -z "$GITHUB_PROXY" ]; then
    GITHUB_PROXY="https://gh-proxy.org/"
    if curl_fetch "$(gh_url "$url")" "$dest" 2>/dev/null; then
      return 0
    fi
    GITHUB_PROXY=""
  fi

  return 1
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

need_cmd curl || err "安装失败: 系统缺少 curl"

PLATFORM="$(detect_platform)"

if [ "$VERSION" = "latest" ]; then
  VERSION="$(resolve_latest_version)" || err "获取版本信息失败，请设置 XRXS_VERSION 或 XRXS_GITHUB_PROXY 后重试"
fi

VER="${VERSION#v}"
ARCHIVE="xrxs_${VER}_${PLATFORM}.tar.gz"
DEFAULT_BASE="https://github.com/${REPO}/releases"
BASE_URL="${XRXS_BASE_URL:-$(gh_url "$DEFAULT_BASE")}"
DOWNLOAD_URL="${BASE_URL}/download/${VERSION}/${ARCHIVE}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

download_archive "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE" \
  || err "下载失败: $ARCHIVE（可尝试: XRXS_GITHUB_PROXY=https://gh-proxy.org/ 或 XRXS_VERSION=${VERSION}）"

tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR"
INSTALL_DIR="$(install_binary "$TMP_DIR/$BIN_NAME")"

if [ "$NO_SKILLS" != "1" ] && [ -d "$TMP_DIR/skills/xrxs" ]; then
  SKILL_SRC_DIR="$TMP_DIR/skills/xrxs"
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
    dest_dir="$HOME/$agent_dir/xrxs"
    rm -rf "$dest_dir" 2>/dev/null || true
    mkdir -p "$(dirname "$dest_dir")" 2>/dev/null || continue
    cp -R "$SKILL_SRC_DIR" "$dest_dir" 2>/dev/null || continue
  done
fi

echo ""
say "薪人薪事CLI 安装完成！安装位置：$INSTALL_DIR/$BIN_NAME"
say ""
say "登录："
say "  xrxs auth login --base-url https://s122.devtest.vip"
