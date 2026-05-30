---
name: xrxs
description: 薪人薪事 HR SaaS 系统命令行工具。当用户需要查询审批列表、查看审批详情、审批操作（通过/驳回/转交/催办）、管理审批类型、审批角色、审批委托时使用。
cli_version: ">=0.1.0"
---

# 薪人薪事 CLI (xrxs)

通过 `xrxs` 命令管理薪人薪事 HR SaaS 系统。

## 安装

```bash
# macOS / Linux
curl -fsSL https://code.qijiayoudao.net/liuxin/xrxs-cli/-/releases/latest/downloads/install.sh | sh

# Windows PowerShell
irm https://code.qijiayoudao.net/liuxin/xrxs-cli/-/releases/latest/downloads/install.ps1 | iex
```

## 全局参数

所有命令继承以下全局参数：

| Flag | 默认值 | 说明 |
|------|--------|------|
| `--base-url` | 环境变量 `XRXS_BASE_URL` 或配置文件 | API 服务地址 |
| `-f, --format` | `json` | 输出格式: `json` / `table` / `raw` |
| `--jq` | - | jq 表达式过滤输出 |
| `--fields` | - | 只输出指定字段，逗号分隔 |
| `-v, --verbose` | `false` | 显示详细请求日志 |
| `--dry-run` | `false` | 预览操作但不执行 |
| `-y, --yes` | `false` | 跳过确认提示 |

## 严格禁止 (NEVER DO)

- 不要使用 xrxs 命令以外的方式操作（禁止 curl、HTTP API、浏览器）
- 不要编造 ID，必须从命令返回中提取
- 不要猜测参数值，操作前必须先查询确认

## 严格要求 (MUST DO)

- 大规模输出时必须加 `--fields` 或 `--jq` 减少 token 消耗
- 所有命令默认使用 JSON 输出（程序化解析），展示给用户时使用 `-f table`
- 危险操作（删除、批量审批）必须先向用户确认

## 命令列表

### 认证

```bash
xrxs auth login --base-url <地址>     # 登录
xrxs auth logout                       # 退出
xrxs auth status                       # 查看登录状态
```

### 审批列表

```bash
# 查询列表
xrxs approval list search [flags]

# 审批操作
xrxs approval list approve  --sid <sid> [--flow-step-id <id>] [--comment <意见>]
xrxs approval list reject   --sid <sid> [--flow-step-id <id>] [--comment <意见>]
xrxs approval list cancel   --sid <sid> [--remark <说明>]
xrxs approval list forward  --sid <sid> --approver-id <id> [--remark <说明>]
xrxs approval list urge     --sid <sid>
xrxs approval list export   [--keyword <关键词>] [--name <文件名>]
```

**search 参数：**

| Flag | 类型 | 说明 |
|------|------|------|
| `--keyword` | string | 搜索关键词 |
| `--status` | string | 0=审批中, 1=已通过, 2=已驳回, 3=已撤销 |
| `--flow-group-id` | string | 审批分组 ID |
| `--page` | int | 页码 (默认 1) |
| `--page-size` | int | 每页条数 (默认 20) |

### 审批详情

```bash
xrxs approval detail get          --sid <sid> [--type 1|2]
xrxs approval detail path         --sid <sid>
xrxs approval detail print        --sid <sid>
xrxs approval detail edit-history --sid <sid>
```

### 审批管理

```bash
xrxs approval manage list    [--group-id <id>]
xrxs approval manage get     --setting-id <id>
xrxs approval manage create  --name <名称> [--name-eng <英文名>] [--description <描述>]
xrxs approval manage update  --setting-id <id> [--name <名称>] [--sub basic|form|flow|advanced]
xrxs approval manage delete  --setting-id <id> [--flow-type 1]
xrxs approval manage toggle  --setting-id <id> [--status on|off]
xrxs approval manage preview --setting-id <id> [--employee-id <id>]
```

### 审批角色

```bash
xrxs approval role list   [--keyword <关键词>]
xrxs approval role get    --role-id <id>
xrxs approval role create --name <名称> [--members <JSON>]
xrxs approval role update --role-id <id> [--name <名称>] [--members <JSON>]
xrxs approval role delete --role-id <id>
xrxs approval role usage  --role-id <id>
```

### 审批委托

```bash
xrxs approval entrust list    [--page 1] [--page-size 20]
xrxs approval entrust create  --agent-id <id> [--scope 0|1] [--start-time <YYYY-MM-DD>] [--end-time <YYYY-MM-DD>] [--reason <原因>]
xrxs approval entrust cancel  --id <id>
xrxs approval entrust setting [--status 0|1]
```

## 典型用法

```bash
# 查看待审批列表（表格形式）
xrxs approval list search --status 0 -f table

# 只看关键字段
xrxs approval list search --status 0 --fields employeeName,title,addDate -f table

# 批量查看所有待审批
xrxs approval list search --status 0 --page-size 100 -f table

# 查看详情
xrxs approval detail get --sid <sid>

# 通过审批
xrxs approval list approve --sid <sid> --comment "同意"
```

## 扩展

未来会支持考勤、人事、薪酬等模块。命令通过 Schema JSON 配置自动生成。
