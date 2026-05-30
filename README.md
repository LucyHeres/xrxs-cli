# xrxs CLI

薪人薪事 HR SaaS 系统命令行工具。

## 安装

```bash
# macOS / Linux
curl -fsSL https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh

# Windows PowerShell
irm https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.ps1 | iex
```

## 使用

```bash
xrxs auth login --base-url <地址>     # 登录
xrxs approval list search --status 0   # 查询审批列表
xrxs approval detail get --sid <sid>   # 查看审批详情
```

## License

MIT
