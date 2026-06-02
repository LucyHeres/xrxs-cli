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

用户说"上个月/本月/某天的审批/最近一周的":
- 按提交时间 → `search --status <status> --addtime 4 --addtime-start <YYYY-MM-DD> --addtime-end <YYYY-MM-DD>`
- 按审批时间 → `search --status <status> --modtime 4 --modtime-start <YYYY-MM-DD> --modtime-end <YYYY-MM-DD>`

用户说"某人的审批/张三提交的":
- 按发起人 → `search --employee-id <id>`

用户说"某部门的审批":
- 按部门 → `search --department-id <id>`

用户说"包含加班/请假内容":
- 按表单内容 → `search --form-content <关键词>`

用户说"上个月审批完成的所有审批":
- 组合条件 → `search --status 1 --modtime 4 --modtime-start <YYYY-MM-DD> --modtime-end <YYYY-MM-DD>`
（注意：上个月→审批时间范围，"完成"→status=1）

用户说"同意/通过这个审批":
- 通过 → `approve --sid <sid> --comment "<审批意见>"`（意见从用户输入提取，用户未提供时追问）

用户说"驳回/拒绝这个审批":
- 驳回 → `reject --sid <sid> --comment "<驳回原因>"`（原因从用户输入提取，用户未提供时必须追问）

用户说"撤销/撤回这个审批":
- 撤销 → `cancel --sid <sid> --remark "<撤销原因>"`（原因从用户输入提取）

用户说"转交/转给某人审批":
- 转交 → `forward --sid <sid> --approver-id <id> --remark "<转交说明>"`

用户说"催一下/催办":
- 催办 → `urge --sid <sid>` 或 `safe-urge --sid <sid>`（推荐，先检查再催办）
- 批量催办 → `batch-urge` 或 `safe-batch-urge`（推荐）

用户说"加签/加签给某人":
- 加签 → `add-sign --sid <sid> --step-node-id <id> --employee-ids <ids> --remark "<原因>"`

用户说"退回/退回到某节点/打回重审":
- 先查可退回节点 → `back-nodes --sid <sid>`
- 再执行退回 → `send-back --sid <sid> --step-node-id <id> --back-node-id <id> --remark "<原因>"`

用户说"导出审批/下载审批列表":
- 导出 → `export`

用户说"看看各状态各有多少审批":
- 统计 → `filter-num`

用户说"看看掩码字段/展示明文":
- 展示掩码 → `unmask --sid <sid> --label-name <字段名>`

**用户只说"帮我查一个审批"（未给任何条件）**:
- 默认 → `search --status 0`（待审批列表）

## 命令总览

### 搜索审批列表

```
Usage:
  xrxs approval list search [flags]
Example:
  # 基础查询
  xrxs approval list search --status 0 -f table
  xrxs approval list search --status 0 --page-size 50 --format json
  xrxs approval list search --keyword "张三" --status 0 -f table
  xrxs approval list search --flow-group-id 2 --status 0 -f table
  # 时间范围查询
  xrxs approval list search --status 1 --modtime 4 --modtime-start "2026-05-01" --modtime-end "2026-05-31" -f table
  xrxs approval list search --status 1 --addtime 4 --addtime-start "2026-05-01" --addtime-end "2026-05-31" -f table
  # 按人/部门/表单内容查询
  xrxs approval list search --employee-id "uid123" --status 0 -f table
  xrxs approval list search --department-id "dept456" --status 0 -f table
  xrxs approval list search --form-content "加班" --status 0 -f table
  # 字段筛选
  xrxs approval list search --status 0 --fields employeeName,title,addDate -f table
Flags:
      --keyword string          搜索关键词
      --status string           状态: 0=审批中, 1=已通过, 2=已驳回, 3=已撤销
      --flow-group-id string    审批分组 ID
      --flow-type string        审批类型 ID
      --employee-id string      发起人 ID
      --department-id string    部门 ID
      --approve-ids string      审批记录 ID（多个逗号分隔）
      --addtime string          提交时间过滤类型: 空=不限, 4=自定义范围
      --addtime-start string    提交时间起始 (YYYY-MM-DD)，addtime=4 时有效
      --addtime-end string      提交时间结束 (YYYY-MM-DD)，addtime=4 时有效
      --modtime string          审批时间过滤类型: 空=不限, 4=自定义范围
      --modtime-start string    审批时间起始 (YYYY-MM-DD)，modtime=4 时有效
      --modtime-end string      审批时间结束 (YYYY-MM-DD)，modtime=4 时有效
      --form-content string     表单内容关键词
      --page int                页码 (默认 1)
      --page-size int           每页条数 (默认 20)
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

### 安全催办（推荐）

> Pipeline 命令：先检查是否可催办，可催办则自动执行。避免盲目催办导致的 API 错误提示。

```
Usage:
  xrxs approval list safe-urge [flags]
Example:
  xrxs approval list safe-urge --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

执行流程：`test-urge` →（可催办）→ `urge`；如不可催办，Pipeline 在第一步终止并返回原因。

### 安全批量催办（推荐）

> Pipeline 命令：先检查是否可批量催办，可催办则自动执行。

```
Usage:
  xrxs approval list safe-batch-urge [flags]
Example:
  xrxs approval list safe-batch-urge --data '<JSON>'
Flags:
      --data string   批量审批 ID JSON（必传）
```

执行流程：`test-batch-urge` →（可催办）→ `batch-urge`。

### 加签

> 将审批加签给其他人，增加审批节点。

```
Usage:
  xrxs approval list add-sign [flags]
Example:
  xrxs approval list add-sign --sid <sid> --step-node-id <id> --employee-ids "id1,id2" --remark "<加签原因>"
Flags:
      --sid string            审批记录 ID（必传）
      --step-node-id string   当前步骤节点 ID（必传）
      --increase-type int     加签类型（必传）
      --employee-ids string   加签员工 ID，多个用逗号分隔（必传）
      --flow-pass-type int    通过类型（必传）
      --remark string         加签原因（必传）
```

### 退回

> **⚠️ 不可逆，将审批退回到之前的节点重新审批。**

```
Usage:
  xrxs approval list send-back [flags]
Example:
  # 先查可退回的节点
  xrxs approval list back-nodes --sid <sid>
  # 再执行退回
  xrxs approval list send-back --sid <sid> --step-node-id <id> --back-node-id <id> --remark "<退回原因>"
Flags:
      --sid string             审批记录 ID（必传）
      --step-node-id string    当前步骤节点 ID（必传）
      --back-node-id string    退回到的节点 ID（必传）
      --remark string          退回原因（必传）
      --re-approval int        是否重新审批（必传）
```

### 获取可退回节点

```
Usage:
  xrxs approval list back-nodes [flags]
Example:
  xrxs approval list back-nodes --sid <sid>
  xrxs approval list back-nodes --sid <sid> -f table
Flags:
      --sid string   审批记录 ID（必传）
```

返回可退回到的历史节点列表（stepNodeId、stepNodeName、isStartNode）。

### 掩码字段展示

> 展示被掩码隐藏的敏感字段明文。

```
Usage:
  xrxs approval list unmask [flags]
Example:
  xrxs approval list unmask --sid <sid> --label-name "手机号" --is-old-value 0 --group-index 0 --is-fixed 1
Flags:
      --sid string         审批记录 ID（必传）
      --label-name string  字段标签名（必传）
      --is-old-value int   是否旧值（必传）
      --group-index int    分组索引（必传）
      --is-fixed int       是否固定字段（必传）
```

### 批量通过/驳回

> **⚠️ 不可逆，确认前须列出所有目标审批。**

```
Usage:
  xrxs approval list batch-approve [flags]
Example:
  xrxs approval list batch-approve --data '<JSON>' --type 1
  xrxs approval list batch-approve --data '<JSON>' --type 2
Flags:
      --data string   审批数据列表 JSON（必传）
      --type int      操作类型: 1=通过, 2=驳回（必传）
```

### 批量转交

> **⚠️ 不可逆，确认前须列出所有目标审批。**

```
Usage:
  xrxs approval list batch-forward [flags]
Example:
  xrxs approval list batch-forward --data '<JSON>'
Flags:
      --data string   转发数据列表 JSON（必传）
```

### 导出审批列表

```
Usage:
  xrxs approval list export [flags]
Example:
  xrxs approval list export
  xrxs approval list export --keyword "请假" --name "请假审批导出"
  xrxs approval list export --modtime 4 --modtime-start "2026-05-01" --modtime-end "2026-05-31" --name "5月审批数据"
Flags:
      --keyword string        搜索关键词
      --name string           导出文件名
      --addtime string        提交时间过滤类型
      --addtime-start string  提交时间起始
      --addtime-end string    提交时间结束
      --modtime string        审批时间过滤类型
      --modtime-start string  审批时间起始
      --modtime-end string    审批时间结束
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
xrxs approval list batch-urge           # 批量催办
xrxs approval list change-prove-status  # 领取/更改证明状态
xrxs approval list export-check         # 导出前检查权限
xrxs approval list export-list          # 导出记录列表
xrxs approval list get-confirm-date     # 离职确认日期
xrxs approval list get-last-sign        # 管理员上次签名
xrxs approval list open-status          # 审批开启状态
xrxs approval list search-employee      # 搜索部门/员工（@功能）
xrxs approval list speedy-search-type   # 猜你想搜下拉选项
xrxs approval list test-batch-urge      # 检查是否可批量催办
xrxs approval list test-urge            # 检查是否可催办（单个）
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
# 按关键词搜索
xrxs approval list search --keyword "张三" --status 0 -f table

# 按时间范围搜索（上个月已通过的审批）
xrxs approval list search --status 1 --modtime 4 --modtime-start "2026-05-01" --modtime-end "2026-05-31" -f table

# 按提交时间搜索（本周提交的）
xrxs approval list search --addtime 4 --addtime-start "2026-05-26" --addtime-end "2026-06-01" -f table

# 按发起人搜索
xrxs approval list search --employee-id "uid123" --status 0 -f table

# 按表单内容搜索
xrxs approval list search --form-content "加班" --status 0 -f table

# 组合搜索
xrxs approval list search --keyword "请假" --status 1 --addtime 4 --addtime-start "2026-05-01" --addtime-end "2026-05-31" --page-size 50 -f table
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

### 工作流 5: 加签与退回

```bash
# 加签：直接加签给其他人
xrxs approval list add-sign --sid <sid> --step-node-id <id> \
  --increase-type 1 --employee-ids "uid1,uid2" --flow-pass-type 1 \
  --remark "需要你们确认"

# 退回：先查可退回节点
xrxs approval list back-nodes --sid <sid>

# 选择目标节点后执行退回
xrxs approval list send-back --sid <sid> --step-node-id <id> \
  --back-node-id <targetId> --remark "信息不完整" --re-approval 1
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `search` | `sid` | approve / reject / cancel / forward / urge / add-sign / send-back 的 --sid |
| `search` | `title` / `employeeName` | 向用户展示审批列表，确认操作目标 |
| `search` | `status` | 判断是否需要操作（已通过/已驳回的无需再操作） |
| `filter-num` | 各状态计数 | 让用户知道待办量 |
| `detail get` | 审批内容 | 确认后再决定 approve / reject |
| `back-nodes` | `stepNodeId` | send-back --back-node-id |
| `export-check` | 审批数量 | 确认是否继续导出 |

## 注意事项

- `approve` / `reject` / `cancel` / `forward` / `send-back` 为**不可逆操作**，执行前必须向用户展示摘要并获得明确确认
- `batch-approve` / `batch-forward` 影响多条记录，确认前应列出所有目标审批
- `send-back` 退回前先用 `back-nodes` 查询可退回的节点列表，确保目标节点有效
- `--status` 参数值为数字字符串: `0`=审批中, `1`=已通过, `2`=已驳回, `3`=已撤销
- `--sid` 是审批记录 ID，不是审批类型 ID（`--flow-setting-id`），不要从 `manage list` 返回中取
- `--flow-group-id` 来自 `approval manage list-groups` 返回的 groupId
- `--addtime` / `--modtime` 为时间过滤类型，传 `4` 时需要配合 `--*-start` 和 `--*-end` 指定自定义范围
- `--addtime` 按**提交时间**过滤（用户说"某月提交的/某天申请的"）；`--modtime` 按**审批时间**过滤（用户说"某月审批完的/通过时间"）
- 用户说"上个月/本月"等相对时间时，需根据当前日期自行计算起止日期
- 自动发现更多命令：`xrxs approval list --help`

## 相关模块

- [detail](./detail.md) — 查看审批详情（操作前先看内容）
- [manage](./manage.md) — 审批类型管理（查分组、查类型配置）
- [proof](./proof.md) — 证明与打印
- [dismiss-handover](./dismiss-handover.md) — 离职交接管理（交接审批操作）
