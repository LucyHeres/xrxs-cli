# 审批详情 (approval detail)

模块入口：`xrxs approval detail`

## 意图判断

用户说"查看审批详情/审批内容/看看这个审批":
- 详情 → `get --sid <sid>`

用户说"看看审批流转/审批经过谁了/审批路径":
- 流转路径 → `path --sid <sid>`

用户说"打印审批/打印视图":
- 打印 → `print --sid <sid>`

用户说"审批的编辑历史/谁改过":
- 编辑历史 → `edit-history --sid <sid>`

用户说"确认日期/离职确认日期":
- 确认日期 → `confirm-date --sid <sid>`

用户说"补打卡次数/这个月补打卡了几次":
- 补打卡 → `bdk-count --sid <sid>`

用户说"加班时长/加班了多少小时":
- 加班时长 → `overtime-hour --sid <sid>`

用户说"保存编辑/提交修改":
- 保存编辑 → `save-edit --sid <sid> --data '<JSON>'`

用户说"校验表单/检查必填字段":
- 校验表单 → `check-form --step-node-id <id>`
- 校验业务字段 → `check-biz --sid <sid> --data '<JSON>'`

用户说"查看增编/编制信息/增编明细":
- 增编详情 → `headcount --headcount-department-id <id> --apply-for-expansion <n> --start-value <YYYY-MM> --end-value <YYYY-MM>`
- 增编明细 → `headcount-detail --headcount-department-id <id> --sid <sid>`

用户说"查一下这个审批/看一下这个审批"（给了 sid）:
- 默认 → `get --sid <sid>`

关键区分: `detail get`(查审批**记录实例**，参数 --sid) vs `manage get`(查审批**类型配置**，参数 --setting-id)

## 命令总览

### 获取审批详情

```
Usage:
  xrxs approval detail get [flags]
Example:
  xrxs approval detail get --sid <sid>
  xrxs approval detail get --sid <sid> --type 1 --format json
  xrxs approval detail get --sid <sid> -f table
Flags:
      --sid string    审批记录 ID（必传）
      --type int      详情类型: 1=基础详情, 2=完整详情（默认 1）
```

返回包含：提交人信息、审批内容/表单数据、审批意见记录、当前审批步骤。

### 获取审批流转路径

```
Usage:
  xrxs approval detail path [flags]
Example:
  xrxs approval detail path --sid <sid>
  xrxs approval detail path --sid <sid> --format json
Flags:
      --sid string   审批记录 ID（必传）
```

返回各步骤审批人、审批状态、审批时间。

### 打印视图

```
Usage:
  xrxs approval detail print [flags]
Example:
  xrxs approval detail print --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

### 编辑历史

```
Usage:
  xrxs approval detail edit-history [flags]
Example:
  xrxs approval detail edit-history --sid <sid>
  xrxs approval detail edit-history --sid <sid> --format json
Flags:
      --sid string   审批记录 ID（必传）
```

### 获取确认日期

```
Usage:
  xrxs approval detail confirm-date [flags]
Example:
  xrxs approval detail confirm-date --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

### 获取补打卡次数

```
Usage:
  xrxs approval detail bdk-count [flags]
Example:
  xrxs approval detail bdk-count --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

### 获取加班时长

```
Usage:
  xrxs approval detail overtime-hour [flags]
Example:
  xrxs approval detail overtime-hour --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

返回加班时长和加班日期。

### 保存编辑

```
Usage:
  xrxs approval detail save-edit [flags]
Example:
  xrxs approval detail save-edit --sid <sid> --data '<JSON>'
Flags:
      --sid string    审批记录 ID（必传）
      --data string   编辑数据 JSON（必传）
```

### 校验表单

```
Usage:
  xrxs approval detail check-form [flags]
Example:
  xrxs approval detail check-form --step-node-id <id>
Flags:
      --step-node-id string   步骤节点 ID（必传）
```

### 校验业务字段

```
Usage:
  xrxs approval detail check-biz [flags]
Example:
  xrxs approval detail check-biz --sid <sid> --data '<JSON>'
Flags:
      --sid string    审批记录 ID（必传）
      --data string   校验数据 JSON（必传）
```

### 增编详情（编辑模式）

```
Usage:
  xrxs approval detail headcount [flags]
Example:
  xrxs approval detail headcount --headcount-department-id <id> --apply-for-expansion 2 --start-value "2026-01" --end-value "2026-06"
Flags:
      --headcount-department-id string   编制部门 ID（必传）
      --apply-for-expansion int          申请编制数（必传）
      --start-value string               开始月份 YYYY-MM（必传）
      --end-value string                 结束月份 YYYY-MM（必传）
```

### 增编明细（详情模式）

```
Usage:
  xrxs approval detail headcount-detail [flags]
Example:
  xrxs approval detail headcount-detail --headcount-department-id <id> --sid <sid>
Flags:
      --headcount-department-id string   编制部门 ID（必传）
      --sid string                       审批记录 ID（必传）
```

## 核心工作流

### 工作流 1: 查看并理解审批流程

```bash
# 1. 从列表拿到 sid 后，先看详情
xrxs approval detail get --sid <sid> --format json

# 2. 如果想知道经过了谁，看流转路径
xrxs approval detail path --sid <sid> --format json

# 3. 如果需要纸质留档，打印
xrxs approval detail print --sid <sid>
```

### 工作流 2: 处理前确认审批内容

```bash
# 1. 查看详情（确认表单数据、审批意见）
xrxs approval detail get --sid <sid> -f table

# 2. 确认后操作（通过/驳回/转交）
xrxs approval list approve --sid <sid> --comment "<审批意见>"
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `get` | 审批内容/表单数据 | 判断是否 approve / reject |
| `get` | 当前审批节点 | 判断自己是否当前审批人 |
| `path` | 各步骤审批人/状态 | 了解流转情况 |
| `edit-history` | 修改记录 | 追溯修改历史 |
| `confirm-date` | 确认日期 | 离职审批补充信息 |
| `bdk-count` | 补打卡次数 | 判断补打卡余额 |
| `overtime-hour` | 加班时长/日期 | 加班审批补充信息 |
| `check-form` | 校验结果 | 表单提交前验证 |
| `headcount` / `headcount-detail` | 编制数据 | 增编审批查看 |

## 注意事项

- `detail get` 查的是审批**记录实例**（参数 `--sid`），**不要**和 `manage get`（查审批类型配置，参数 `--setting-id`）混淆
- `--sid` 来自 `approval list search` 返回的 sid 字段，不是来自 `manage list` 的 flowSettingId
- `--type` 不同值返回的字段和深度不同，需要完整数据时传 `--type 2`
- `save-edit` 提交的是 JSON 数据，结构需与表单定义一致
- 自动发现更多命令：`xrxs approval detail --help`

## 相关模块

- [list](./list.md) — 审批列表查询与操作（操作前先看详情）
- [manage](./manage.md) — 审批类型配置管理（注意 manage get 和 detail get 的区别）
- [proof](./proof.md) — 证明与打印
- [dismiss-handover](./dismiss-handover.md) — 离职交接管理
