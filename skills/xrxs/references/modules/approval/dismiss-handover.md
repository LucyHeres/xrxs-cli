# 离职交接管理 (dismiss handover)

模块入口：`xrxs approval dismiss-handover`

## 意图判断

用户说"查看离职交接方案/交接方案列表":
- 方案列表 → `list-settings`

用户说"查看交接方案详情/这个方案怎么配的":
- 方案详情 → `get-setting --handover-basic-id <id>`

用户说"新增/修改交接方案/配置离职交接":
- 保存方案 → `save-setting --data '<JSON>'`

用户说"删除交接方案":
- 删除方案 → `delete-setting --handover-basic-id <id>`

用户说"复制交接方案":
- 复制方案 → `copy-setting --data '<JSON>'`

用户说"交接适用范围/适用范围筛选":
- 范围筛选 → `scope-filter`

用户说"查看离职交接审批详情/交接审批内容":
- 交接审批详情 → `detail --flow-process-id <id>`

用户说"转交离职交接/交接审批换人":
- 转交 → `transfer --approval-id <id> --director-employee-id <id> --transfer-reason "<原因>"`

用户说"撤销离职交接/取消离职交接审批":
- 撤销 → `cancel --approval-id <id> --cancel-reason "<原因>"`

用户说"催办离职交接":
- 催办 → `urge --approval-id <id>`

## 命令总览

### 交接方案列表

```
Usage:
  xrxs approval dismiss-handover list-settings [flags]
Example:
  xrxs approval dismiss-handover list-settings
  xrxs approval dismiss-handover list-settings -f table
  xrxs approval dismiss-handover list-settings --format json
Flags:
  无必填参数
```

返回字段：

| 字段 | 说明 |
|------|------|
| `handoverBasicId` | 方案 ID |
| `handoverName` | 方案名称 |
| `scope` | 适用范围 |

### 交接方案详情

```
Usage:
  xrxs approval dismiss-handover get-setting [flags]
Example:
  xrxs approval dismiss-handover get-setting --handover-basic-id 123
  xrxs approval dismiss-handover get-setting --handover-basic-id 123 --format json
Flags:
      --handover-basic-id int   方案 ID（必传）
```

### 保存交接方案

```
Usage:
  xrxs approval dismiss-handover save-setting [flags]
Example:
  xrxs approval dismiss-handover save-setting --data '<JSON>'
Flags:
      --data string   方案数据 JSON（必传）
```

### 删除交接方案

> **⚠️ 不可逆，执行前须向用户确认。**

```
Usage:
  xrxs approval dismiss-handover delete-setting [flags]
Example:
  xrxs approval dismiss-handover delete-setting --handover-basic-id 123
Flags:
      --handover-basic-id int   方案 ID（必传）
```

### 复制交接方案

```
Usage:
  xrxs approval dismiss-handover copy-setting [flags]
Example:
  xrxs approval dismiss-handover copy-setting --data '<JSON>'
Flags:
      --data string   方案数据 JSON（必传）
```

### 适用范围筛选

```
Usage:
  xrxs approval dismiss-handover scope-filter [flags]
Example:
  xrxs approval dismiss-handover scope-filter
  xrxs approval dismiss-handover scope-filter --format json
Flags:
  无必填参数
```

### 交接审批详情

```
Usage:
  xrxs approval dismiss-handover detail [flags]
Example:
  xrxs approval dismiss-handover detail --flow-process-id 12345
  xrxs approval dismiss-handover detail --flow-process-id 12345 --is-print 0
Flags:
      --flow-process-id int   审批流程 ID（必传）
      --is-print int          是否打印: 0=否, 1=是（默认 0）
```

### 转交离职交接审批

> **⚠️ 不可逆，执行前须向用户确认转交目标人。**

```
Usage:
  xrxs approval dismiss-handover transfer [flags]
Example:
  xrxs approval dismiss-handover transfer --approval-id 123 --director-employee-id "uid123" --transfer-reason "原负责人离职"
Flags:
      --approval-id int            审批 ID（必传）
      --director-employee-id string  转交给的员工 ID（必传）
      --transfer-reason string     转交原因（必传）
```

### 撤销离职交接审批

> **⚠️ 不可逆，执行前须向用户确认。**

```
Usage:
  xrxs approval dismiss-handover cancel [flags]
Example:
  xrxs approval dismiss-handover cancel --approval-id 123 --cancel-reason "信息有误"
Flags:
      --approval-id int     审批 ID（必传）
      --cancel-reason string  撤销原因（必传）
```

### 催办离职交接

```
Usage:
  xrxs approval dismiss-handover urge [flags]
Example:
  xrxs approval dismiss-handover urge --approval-id 123
Flags:
      --approval-id int   审批 ID（必传）
```

## 核心工作流

### 工作流 1: 配置离职交接方案

```bash
# 1. 查看现有方案
xrxs approval dismiss-handover list-settings -f table

# 2. 查看方案详情 — 提取配置数据
xrxs approval dismiss-handover get-setting --handover-basic-id <id> --format json

# 3. 修改后保存
xrxs approval dismiss-handover save-setting --data '<JSON>'

# 4. 或复制已有方案快速创建
xrxs approval dismiss-handover copy-setting --data '<JSON>'
```

### 工作流 2: 处理离职交接审批

```bash
# 1. 查看交接审批详情
xrxs approval dismiss-handover detail --flow-process-id <id> --format json

# 2. 催办
xrxs approval dismiss-handover urge --approval-id <id>

# 3. 转交给其他人
xrxs approval dismiss-handover transfer \
  --approval-id <id> \
  --director-employee-id "<uid>" \
  --transfer-reason "<原因>"

# 4. 撤销
xrxs approval dismiss-handover cancel --approval-id <id> --cancel-reason "<原因>"
```

### 工作流 3: 查看适用范围

```bash
# 查看适用范围筛选项
xrxs approval dismiss-handover scope-filter --format json
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `list-settings` | `handoverBasicId` | `get-setting / delete-setting --handover-basic-id` |
| `list-settings` | `handoverName` | 向用户展示方案名称 |
| `get-setting` | 方案完整配置 | `save-setting / copy-setting` 参考 |
| `detail` | 交接审批内容 | 判断是否 transfer / cancel / urge |
| `scope-filter` | 筛选项列表 | `save-setting` 设置适用范围 |

## 注意事项

- `delete-setting` 删除的是离职交接方案，不是审批记录，执行前务必确认
- `transfer` / `cancel` 为**不可逆操作**，执行前须向用户展示摘要并获得确认
- `detail` 的 `--flow-process-id` 来自离职交接审批流，不是普通审批的 sid
- `save-setting` 和 `copy-setting` 提交的 `--data` 是完整方案 JSON，建议先 `get-setting` 查看现有结构
- 自动发现更多命令：`xrxs approval dismiss-handover --help`

## 相关模块

- [list](./list.md) — 审批列表与操作（含离职交接审批的转交、撤销）
- [detail](./detail.md) — 审批详情
- [manage](./manage.md) — 审批类型管理
