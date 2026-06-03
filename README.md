# xrxs CLI

薪人薪事 HR SaaS 系统命令行工具。安装后，在 Claude Code 中通过 `/xrxs` Skill 直接对话操作审批，无需记忆命令。

## 安装

**macOS / Linux**（推荐 Release 安装包，不要用 `raw.githubusercontent.com`，国内易超时）：

```bash
curl -fsSL --connect-timeout 30 --max-time 300 \
  https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh
```

**国内网络较慢时**（安装脚本与二进制均走镜像）：

```bash
curl -fsSL --connect-timeout 30 --max-time 300 \
  "https://gh-proxy.org/https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh" | sh
```

或已下载脚本后：

```bash
XRXS_GITHUB_PROXY=https://gh-proxy.org/ sh install.sh
```

**Windows PowerShell**

```powershell
irm https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.ps1 | iex
```

安装完成后，打开新终端窗口即可使用。

## 登录

```bash
xrxs auth login --base-url https://s122.devtest.vip
```

按提示输入账号和密码。登录成功后，在 Claude Code 中输入 `/xrxs` 即可通过自然语言操作审批。

## 升级

```bash
xrxs upgrade
```

自动检查并更新到最新版本。

## 卸载

```bash
xrxs uninstall
```

或手动删除：

```bash
rm -f ~/.local/bin/xrxs /usr/local/bin/xrxs
rm -rf ~/.xrxs
```

## License

MIT
