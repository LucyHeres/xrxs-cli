# xrxs CLI

薪人薪事 HR SaaS 系统命令行工具。在终端中管理审批、查看详情、操作流程，无需打开浏览器。

## 安装

**macOS / Linux**

```bash
curl -fsSL https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh
```

**Windows PowerShell**

```powershell
irm https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.ps1 | iex
```

安装完成后，打开新终端窗口即可使用 `xrxs` 命令。

## 升级

```bash
xrxs upgrade
```

自动检查最新版本并更新。如果从 v0.2.1 或更早版本升级，重新执行安装命令即可。

## 快速开始

```bash
# 1. 登录
xrxs auth login --base-url https://your-company.example.com

# 2. 查看审批中的审批（表格形式）
xrxs approval list search --status 0 -f table

# 3. 只看关键字段
xrxs approval list search --status 0 --fields employeeName,title,addDate -f table

# 4. 查看审批详情
xrxs approval detail get --sid <审批ID>

# 5. 查看审批流程路径
xrxs approval detail path --sid <审批ID>
```

## 审批操作

```bash
# 通过审批
xrxs approval list approve --sid <审批ID> --comment "同意"

# 驳回审批
xrxs approval list reject --sid <审批ID> --comment "请补充材料"

# 转交审批
xrxs approval list forward --sid <审批ID> --next-approver-id <目标人ID>

# 催办
xrxs approval list urge --sid <审批ID>

# 撤销审批
xrxs approval list cancel --sid <审批ID>

# 批量操作
xrxs approval list batch-approve --type 1 --data '[{"sid":"xxx","flowStepId":"yyy"}]'
xrxs approval list batch-forward --type 1 --data '[...]'
```

## 审批类型管理

```bash
# 查看所有分组
xrxs approval manage list-groups -f table

# 查看某分组下的审批类型
xrxs approval manage list --group-id <分组ID> -f table

# 创建审批类型
xrxs approval manage create --name "自定义审批" --group-id <分组ID>

# 启用/停用
xrxs approval manage toggle --setting-id <类型ID> --flow-type <流程类型> --status on

# 删除审批类型
xrxs approval manage delete --setting-id <类型ID> --flow-type <流程类型>
```

## 角色与委托

```bash
# 审批角色管理
xrxs approval role list -f table
xrxs approval role get --role-id <角色ID>
xrxs approval role create --name "财务审批角色" --members '[{"employeeId":"xxx"}]'

# 审批委托
xrxs approval entrust list -f table
xrxs approval entrust create --agent-id <被委托人ID> --start-time 2026-06-01 --end-time 2026-06-30 --reason "休假期间委托"
xrxs approval entrust cancel --id <委托ID>
```

## 证明与打印

```bash
xrxs approval proof get-proof-snap --sid <审批ID>
xrxs approval proof add-proof-record --sid <审批ID> --record-status 1
```

## 命令一览

### 审批列表 (list)
`search` `admin-search` `admin-search-all` `speedy-search-type` `filter-num` `open-status` `approve` `reject` `batch-approve` `batch-forward` `cancel` `forward` `urge` `batch-urge` `test-urge` `test-batch-urge` `export` `export-check` `get-last-sign` `search-employee` `change-prove-status` `cancel-handover` `get-confirm-date` `transfer-approval`

### 审批详情 (detail)
`get` `path` `print` `edit-history` `save-page-edit` `overtime-hour` `back-process-step` `save-back-process-step` `increase-process-step` `open-mask-field` `check-form-required` `check-biz-field` `handover-detail`

### 审批类型管理 (manage)
`list-groups` `save-group` `remove-group` `list` `get` `check-name` `create` `update` `delete` `toggle` `preview` `get-init` `list-permission` `list-related` `get-flow-detail` `save-flow-detail` `verify-flow` `copy-flow` `get-charge-setting` `get-advanced-setting` `save-advanced-setting` `get-all-roles`

### 审批角色 (role)
`list` `get` `create` `update` `delete` `usage` `department-tips` `list-all`

### 审批委托 (entrust)
`list` `create` `cancel` `setting` `list-flow-types`

### 证明与打印 (proof)
`get-proof-snap` `add-proof-record` `search-proof-record` `get-bdk-count` `print-flow` `edit-print` `save-print` `get-print-hide-fields`

## 全局参数

| 参数 | 说明 |
|------|------|
| `--base-url` | API 服务地址 |
| `-f, --format` | 输出格式: `json` / `table` / `raw`（默认 `json`） |
| `--fields` | 只输出指定字段，逗号分隔 |
| `--jq` | jq 表达式过滤 JSON 输出 |
| `-v, --verbose` | 显示详细请求日志 |
| `-y, --yes` | 跳过确认提示 |
| `--dry-run` | 预览操作但不执行 |

## 环境变量

| 变量 | 说明 |
|------|------|
| `XRXS_BASE_URL` | 默认 API 服务地址 |
| `XRXS_SCHEMA_DIR` | 自定义 Schema 目录（开发用） |
| `XRXS_CONFIG_DIR` | 配置文件目录（默认 `~/.xrxs`） |
| `HTTPS_PROXY` / `HTTP_PROXY` | 代理设置 |

## License

MIT
