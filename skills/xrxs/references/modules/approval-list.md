# 审批列表查询与操作 (approval list)

模块入口：`xrxs approval list`

## 意图判断

用户说"查看待审批/待我审批/审批列表/看看有哪些审批":
- 查审批中 → `search --status 0`

用户说"已通过的审批/已经批完的":
- 查已通过 → `search --status 1`

用户说"被驳回的/被打回来的":
- 查已驳回 → `search --status 2`

用户说"我撤销的/已撤销的":
- 查已撤销 → `search --status 3`

用户说"搜一下某人的审批/包含某关键词的审批":
- 关键词搜索 → `search --keyword <关键词>`

用户说"请假审批/报销审批/某分组的审批":
- 按分组筛选 → `search --flow-group-id <id>`

用户说"同意/通过这个审批":
- 通过 → `approve --sid <sid> --comment "<审批意见>"`（意见从用户输入提取，用户未提供时追问）

用户说"驳回/拒绝这个审批":
- 驳回 → `reject --sid <sid> --comment "<驳回原因>"`（原因从用户输入提取，用户未提供时必须追问）

用户说"撤销/撤回这个审批":
- 撤销 → `cancel --sid <sid> --remark "<撤销原因>"`（原因从用户输入提取）

用户说"转交/转给某人审批":
- 转交 → `forward --sid <sid> --approver-id <id> --remark "<转交说明>"`

用户说"催一下/催办":
- 催办 → `urge --sid <sid>`

用户说"导出审批/下载审批列表":
- 导出 → `export`

用户说"看看各状态各有多少审批":
- 统计 → `filter-num`

**用户只说"帮我查一个审批"（未给任何条件）**:
- 默认 → `search --status 0`（待审批列表）

## 命令总览

### 搜索审批列表

```
Usage:
  xrxs approval list search [flags]
Example:
  xrxs approval list search --status 0 -f table
  xrxs approval list search --status 0 --page-size 50 --format json
  xrxs approval list search --keyword "张三" --status 0 -f table
  xrxs approval list search --flow-group-id 2 --status 0 -f table
  xrxs approval list search --status 0 --fields employeeName,title,addDate -f table
Flags:
      --keyword string        搜索关键词
      --status string         状态: 0=审批中, 1=已通过, 2=已驳回, 3=已撤销
      --flow-group-id string  审批分组 ID
      --page int              页码 (默认 1)
      --page-size int         每页条数 (默认 20)
```

### 通过审批

> **⚠️ 不可逆，执行前须向用户展示审批摘要并获得确认。**

```
Usage:
  xrxs approval list approve [flags]
Example:
  xrxs approval list approve --sid <sid>
  xrxs approval list approve --sid <sid> --comment "<审批意见>"
  xrxs approval list approve --sid <sid> --flow-step-id <id> --comment "<审批意见>"
Flags:
      --sid string           审批记录 ID（必传）
      --flow-step-id string  审批步骤 ID
      --comment string       审批意见
```

### 驳回审批

> **⚠️ 不可逆，执行前须向用户展示审批摘要并获得确认。**

```
Usage:
  xrxs approval list reject [flags]
Example:
  xrxs approval list reject --sid <sid>
  xrxs approval list reject --sid <sid> --comment "<驳回原因>"
  xrxs approval list reject --sid <sid> --flow-step-id <id> --comment "<驳回原因>"
Flags:
      --sid string           审批记录 ID（必传）
      --flow-step-id string  审批步骤 ID
      --comment string       驳回意见
```

### 撤销审批

> **⚠️ 不可逆，执行前须向用户展示审批摘要并获得确认。**

```
Usage:
  xrxs approval list cancel [flags]
Example:
  xrxs approval list cancel --sid <sid>
  xrxs approval list cancel --sid <sid>
  xrxs approval list cancel --sid <sid> --remark "<撤销原因>"
Flags:
      --sid string    审批记录 ID（必传）
      --remark string  撤销说明
```

### 转交审批

> **⚠️ 不可逆，执行前须向用户确认转交目标人。**

```
Usage:
  xrxs approval list forward [flags]
Example:
  xrxs approval list forward --sid <sid> --approver-id <id>
  xrxs approval list forward --sid <sid> --approver-id <id>
  xrxs approval list forward --sid <sid> --approver-id <id> --remark "<转交说明>"
Flags:
      --sid string         审批记录 ID（必传）
      --approver-id string  转交目标人 ID（必传）
      --remark string       转交说明
```

### 催办

```
Usage:
  xrxs approval list urge [flags]
Example:
  xrxs approval list urge --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

### 导出审批列表

```
Usage:
  xrxs approval list export [flags]
Example:
  xrxs approval list export
  xrxs approval list export --keyword "请假" --name "请假审批导出"
Flags:
      --keyword string  搜索关键词
      --name string     导出文件名
```

### 各状态审批数量统计

```
Usage:
  xrxs approval list filter-num
Example:
  xrxs approval list filter-num --format json
Flags:
  无必填参数
```

### 其他子命令

```bash
xrxs approval list admin-search         # 管理/关联审批列表
xrxs approval list admin-search-all     # 审批全局搜索
xrxs approval list batch-approve        # 批量通过/驳回（须确认）
xrxs approval list batch-forward        # 批量转交（须确认）
xrxs approval list batch-urge           # 批量催办
xrxs approval list cancel-handover      # 撤销离职交接
xrxs approval list change-prove-status  # 领取/更改证明状态
xrxs approval list export-check         # 导出前检查权限
xrxs approval list get-confirm-date     # 离职确认日期
xrxs approval list get-last-sign        # 管理员上次签名
xrxs approval list open-status          # 审批开启状态
xrxs approval list search-employee      # 搜索部门/员工（@功能）
xrxs approval list speedy-search-type   # 猜你想搜下拉选项
xrxs approval list test-batch-urge      # 检查是否可批量催办
xrxs approval list test-urge            # 检查是否可催办（单个）
xrxs approval list transfer-approval    # 转交离职交接审批
```

## 核心工作流

### 工作流 1: 查看并处理待审批

```bash
# 1. 查看待审批列表 — 提取 sid
xrxs approval list search --status 0 -f table

# 2. 查看某条审批详情（确认内容后决定操作）
xrxs approval detail get --sid <sid> --format json

# 3. 操作（从用户输入提取意见）
xrxs approval list approve --sid <sid> --comment "<审批意见>"
xrxs approval list reject --sid <sid> --comment "<驳回原因>"
xrxs approval list forward --sid <sid> --approver-id <id> --remark "<转交说明>"
xrxs approval list cancel --sid <sid> --remark "<撤销原因>"
```

### 工作流 2: 搜索特定审批

```bash
# 1. 按关键词搜索
xrxs approval list search --keyword "张三" --status 0 -f table

# 2. 按分组搜索
xrxs approval list search --flow-group-id 2 --status 0 -f table

# 3. 组合搜索
xrxs approval list search --keyword "请假" --status 1 --page-size 50 -f table
```

### 工作流 3: 批量处理

```bash
# 1. 批量通过
xrxs approval list batch-approve  # 执行前确认

# 2. 批量转交
xrxs approval list batch-forward  # 执行前确认

# 3. 批量催办
xrxs approval list batch-urge
```

### 工作流 4: 导出审批数据

```bash
# 1. 检查导出权限
xrxs approval list export-check

# 2. 导出
xrxs approval list export --keyword "请假" --name "请假审批_2026Q2"
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `search` | `sid` | approve / reject / cancel / forward / urge 的 --sid |
| `search` | `title` / `employeeName` | 向用户展示审批列表，确认操作目标 |
| `search` | `status` | 判断是否需要操作（已通过/已驳回的无需再操作） |
| `filter-num` | 各状态计数 | 让用户知道待办量 |
| `detail get` | 审批内容 | 确认后再决定 approve / reject |

## 注意事项

- `approve` / `reject` / `cancel` / `forward` 为**不可逆操作**，执行前必须向用户展示摘要并获得明确确认
- `batch-approve` / `batch-forward` 影响多条记录，确认前应列出所有目标审批
- `--status` 参数值为数字字符串: `0`=审批中, `1`=已通过, `2`=已驳回, `3`=已撤销
- `--sid` 是审批记录 ID，不是审批类型 ID（`--flow-setting-id`），不要从 `manage list` 返回中取
- `--flow-group-id` 来自 `approval manage list-groups` 返回的 groupId
- 自动发现更多命令：`xrxs approval list --help`

## 相关模块

- [approval-detail](./approval-detail.md) — 查看审批详情（操作前先看内容）
- [approval-manage](./approval-manage.md) — 审批类型管理（查分组、查类型配置）
- [approval-proof](./approval-proof.md) — 证明与打印
