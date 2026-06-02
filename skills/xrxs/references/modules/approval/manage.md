# 审批类型管理 (approval manage)

模块入口：`xrxs approval manage`

## 意图判断

用户说"看一下有哪些审批分组/审批分了哪些组/分组列表":
- 列分组 → `list-groups`

用户说"看看某分组下有哪些审批类型/这个分组有哪些审批":
- 查分组下类型 → `list --group-id <id>`

用户说"查看审批类型配置/这个审批类型的表单怎么配的/流程怎么设置的":
- 查类型配置 → `get --setting-id <id>`

用户说"新增一个审批类型/创建审批":
- 创建 → `create --name <名称>`

用户说"修改审批类型/更新审批配置/编辑表单/编辑流程":
- 更新 → `update --setting-id <id> --sub <basic|form|flow|advanced>`

用户说"删除审批类型/删掉这个审批":
- 删除 → `delete --setting-id <id>`

用户说"启用/停用审批类型/开启/关闭审批":
- 开关 → `toggle --setting-id <id> --status <on|off>`

用户说"预览审批表单/预览审批流程":
- 预览 → `preview --setting-id <id>`

用户说"复制审批流程/把流程从A复制到B":
- 复制流程 → `copy-flow --cover-setting-id <srcId> --setting-id <dstId>`

用户说"验证审批流程/检查流程配置":
- 验证 → `verify --flow-type <int> --setting-id <id>`

用户说"审批发起范围/谁能发起审批":
- 查范围 → `flow-switch --flow-type <int>`
- 保存范围 → `save-flow-switch --flow-type <int>`

用户说"审批分支条件/请假类型/报销类型":
- 分支条件 → `charge-setting`

用户说"新增/重命名分组":
- 分组操作 → `save-group --name <名称>` 或 `save-group --name <名称> --group-id <id>`

用户说"删除分组":
- 删除分组 → `remove-group --group-id <id>`

**用户只说"列出审批类型/有哪些审批类型"（未提分组）**:
- 用 `list-all-types` 一键拉取所有分组及其类型（Pipeline 命令，自动并发）
- 或手动：先调 `list-groups` 拿所有分组 ID，再逐个 `list --group-id <id>`，最后汇总展示

关键区分: `manage list`(查分组下的审批类型，需 groupId) vs `manage list-groups`(查分组元信息，无参数)
关键区分: `manage get`(查审批类型配置模板，参数 --setting-id) vs `detail get`(查审批记录实例，参数 --sid)

## 命令总览

### 列出所有审批分组

```
Usage:
  xrxs approval manage list-groups [flags]
Example:
  xrxs approval manage list-groups
  xrxs approval manage list-groups -f table
  xrxs approval manage list-groups --jq '.[] | {groupId, groupName: ._groupName}'
Flags:
  无必填参数
```

返回字段：

| 字段 | 说明 |
|------|------|
| `groupId` | 分组 ID |
| `_groupName` | 分组中文名（如"员工审批"、"考勤审批"、"自定义审批"） |
| `groupName` | 分组名称 JSON（含 showVal / defVal） |
| `isFixed` | 是否固定分组（1=系统固定，0=可删除） |
| `orderNum` | 排序号 |

### 一键列出所有分组及其审批类型（推荐）

> Pipeline 命令：自动获取所有分组，并发拉取每个分组下的审批类型，合并输出。一条命令替代 "list-groups + 逐个 list"。

```
Usage:
  xrxs approval manage list-all-types [flags]
Example:
  xrxs approval manage list-all-types -f table
  xrxs approval manage list-all-types --format json
Flags:
      --group-id string   审批分组 ID（可选，不传则拉取所有分组）
```

返回格式：`[{group: {groupId, groupName, ...}, types: [{flowSettingId, name, ...}, ...]}, ...]`

并发度默认为 4，某个分组拉取失败自动跳过（不影响其他分组结果）。

### 按分组列出审批类型

```
Usage:
  xrxs approval manage list [flags]
Example:
  xrxs approval manage list --group-id 2
  xrxs approval manage list --group-id 0 -f table
  xrxs approval manage list --group-id 2 --fields name,openStatus,flowSettingId -f table
Flags:
      --group-id string   审批分组 ID（必传）
```

返回字段：

| 字段 | 说明 |
|------|------|
| `flowSettingId` | 审批类型 ID（后续命令的 --setting-id） |
| `flowType` | 流程类型编号 |
| `name` | 审批类型名称 |
| `openStatus` | 启用状态（true/false） |
| `type` | 0=系统预设, 1=自定义 |
| `complete` | 配置完成度（1=完整） |

### 获取审批类型配置详情

```
Usage:
  xrxs approval manage get [flags]
Example:
  xrxs approval manage get --setting-id 8483932
  xrxs approval manage get --setting-id 8483932 --format json
Flags:
      --setting-id string   审批类型 ID（必传）
```

### 创建审批类型

```
Usage:
  xrxs approval manage create [flags]
Example:
  xrxs approval manage create --name "自定义审批"
  xrxs approval manage create --name "采购审批" --name-eng "Purchase" --description "采购流程"
Flags:
      --name string         审批类型名称（必传）
      --name-eng string     英文名称
      --description string  描述
```

### 更新审批类型

```
Usage:
  xrxs approval manage update [flags]
Example:
  xrxs approval manage update --setting-id 8483958 --name "新名称"
  xrxs approval manage update --setting-id 8483958 --sub form
  xrxs approval manage update --setting-id 8483958 --sub flow
Flags:
      --setting-id string   审批类型 ID（必传）
      --name string         新名称
      --sub string          子模块: basic(基本信息) / form(表单设计) / flow(流程设置) / advanced(高级设置)
```

### 删除审批类型

> **CAUTION:** 不可逆操作，执行前必须向用户确认。

```
Usage:
  xrxs approval manage delete [flags]
Example:
  xrxs approval manage delete --setting-id 8961118
  xrxs approval manage delete --setting-id 8483958 --flow-type 1
Flags:
      --setting-id string   审批类型 ID（必传）
      --flow-type int       流程类型（自定义审批时传入，默认 1）
```

### 启停审批类型

```
Usage:
  xrxs approval manage toggle [flags]
Example:
  xrxs approval manage toggle --setting-id 8483942 --status on
  xrxs approval manage toggle --setting-id 8483942 --status off
Flags:
      --setting-id string   审批类型 ID（必传）
      --status string       开关: on(启用) / off(停用)
```

### 预览审批

```
Usage:
  xrxs approval manage preview [flags]
Example:
  xrxs approval manage preview --setting-id 8483935
  xrxs approval manage preview --setting-id 8483935 --employee-id <id>
Flags:
      --setting-id string    审批类型 ID（必传）
      --employee-id string   员工 ID（模拟该员工视角预览）
```

### 新增/重命名审批分组

```
Usage:
  xrxs approval manage save-group [flags]
Example:
  xrxs approval manage save-group --name "新分组"
  xrxs approval manage save-group --name "重命名" --group-id 5
Flags:
      --name string      分组名称（必传）
      --group-id string  分组 ID（重命名时传入）
```

### 删除审批分组

```
Usage:
  xrxs approval manage remove-group [flags]
Example:
  xrxs approval manage remove-group --group-id 5
Flags:
      --group-id string   分组 ID（必传）
```

### 获取新审批类型 ID 列表

```
Usage:
  xrxs approval manage type-ids [flags]
Example:
  xrxs approval manage type-ids
  xrxs approval manage type-ids --format json
Flags:
  无必填参数
```

### 获取新建审批初始化数据

```
Usage:
  xrxs approval manage init-create [flags]
Example:
  xrxs approval manage init-create
  xrxs approval manage init-create --format json
Flags:
  无必填参数
```

### 检查审批名称重复

```
Usage:
  xrxs approval manage check-name [flags]
Example:
  xrxs approval manage check-name --setting-id 8483958 --flow-name "自定义审批"
Flags:
      --setting-id int     设置 ID（必传）
      --flow-name string   审批名称（必传）
      --flow-name-eng string  英文名称
```

### 验证审批流程

```
Usage:
  xrxs approval manage verify [flags]
Example:
  xrxs approval manage verify --flow-type 1 --setting-id 8483958
Flags:
      --flow-type int     流程类型（必传）
      --setting-id int    审批类型 ID（必传）
```

### 复制审批流到目标类型

> 将一个审批流配置完整复制到另一个类型。

```
Usage:
  xrxs approval manage copy-flow [flags]
Example:
  xrxs approval manage copy-flow --cover-setting-id 8483935 --setting-id 8483958
Flags:
      --cover-setting-id int   源类型 ID（必传）
      --setting-id int         目标类型 ID（必传）
```

### 获取审批发起范围

```
Usage:
  xrxs approval manage flow-switch [flags]
Example:
  xrxs approval manage flow-switch --flow-type 1
  xrxs approval manage flow-switch --flow-type 1 --format json
Flags:
      --flow-type int   审批类型（必传）
```

返回：管理员可发起(adminIsOpen)、主管可发起(leaderIsOpen)、员工可发起(employeeIsOpen)。

### 保存审批发起范围

```
Usage:
  xrxs approval manage save-flow-switch [flags]
Example:
  xrxs approval manage save-flow-switch --flow-type 1 --admin-is-open 1 --leader-is-open 1 --employee-is-open 0
Flags:
      --flow-type int          审批类型（必传）
      --admin-is-open int      管理员可发起: 1=是, 0=否（默认 1）
      --leader-is-open int     主管可发起: 1=是, 0=否（默认 1）
      --employee-is-open int   员工可发起: 1=是, 0=否（默认 1）
```

### 获取分支条件

```
Usage:
  xrxs approval manage charge-setting [flags]
Example:
  xrxs approval manage charge-setting
  xrxs approval manage charge-setting --format json
Flags:
  无必填参数
```

返回审批流分支条件（如请假类型列表）。

### 其他子命令

```bash
xrxs approval manage get-flow-detail-new     # 获取公司审批流详情（新）
xrxs approval manage get-flow-detail         # 获取公司审批流详情（旧）
xrxs approval manage save-flow-detail-new    # 保存审批流设置（新）
xrxs approval manage get-advanced            # 获取审批高级设置
xrxs approval manage save-advanced           # 保存审批高级设置
xrxs approval manage get-all-roles           # 获取所有可用审批角色（用于流程配置时选角色）
```

## 核心工作流

### 工作流 1: 浏览全部分组和审批类型

```bash
# 推荐：一键拉取所有分组及其类型（Pipeline 命令，自动并发）
xrxs approval manage list-all-types -f table
xrxs approval manage list-all-types --format json

# 或手动遍历（仅当需要自定义筛选时使用）
xrxs approval manage list-groups --format json         # 1. 列全部分组
xrxs approval manage list --group-id 0 --format json   # 2. 逐个查类型
xrxs approval manage list --group-id 1 --format json
# ...
```

### 工作流 2: 查看审批类型配置

```bash
# 1. 先看分组下有哪些类型 — 提取 flowSettingId
xrxs approval manage list --group-id 2 --format json

# 2. 获取类型配置详情
xrxs approval manage get --setting-id 8483935 --format json
```

### 工作流 3: 创建自定义审批类型

```bash
# 1. 创建审批类型 — 提取 flowSettingId
xrxs approval manage create --name "新审批" --name-eng "NewFlow" --format json

# 2. 编辑基本信息
xrxs approval manage update --setting-id <SETTING_ID> --sub basic

# 3. 设计表单
xrxs approval manage update --setting-id <SETTING_ID> --sub form

# 4. 设置流程
xrxs approval manage update --setting-id <SETTING_ID> --sub flow

# 5. 启用
xrxs approval manage toggle --setting-id <SETTING_ID> --status on
```

### 工作流 4: 停用并删除审批类型

> 停用和删除均为不可逆操作，**必须先向用户确认**。

```bash
# 1. 停用
xrxs approval manage toggle --setting-id <SETTING_ID> --status off

# 2. (可选) 确认后删除
xrxs approval manage delete --setting-id <SETTING_ID>
```

### 工作流 5: 复制审批流配置

```bash
# 将已有审批流的配置复制到新类型
xrxs approval manage copy-flow --cover-setting-id <源ID> --setting-id <目标ID>
```

### 工作流 6: 配置发起范围

```bash
# 1. 查看当前发起范围
xrxs approval manage flow-switch --flow-type 1

# 2. 调整范围（例：仅管理员和主管可发起）
xrxs approval manage save-flow-switch --flow-type 1 --admin-is-open 1 --leader-is-open 1 --employee-is-open 0
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `list-groups` | `groupId` | `list / save-group / remove-group --group-id` |
| `list` | `flowSettingId` | `get / update / delete / toggle / preview / verify / copy-flow --setting-id` |
| `list` | `name` | 向用户展示审批类型名称，确认操作目标 |
| `list` | `openStatus` | 判断是否需要启停 |
| `list` | `flowType` | `flow-switch / verify / copy-flow --flow-type` |
| `create` / `init-create` | `flowSettingId` | `update / delete / toggle / preview --setting-id` |
| `get` | 表单/流程配置 | `update --sub form / --sub flow` |
| `flow-switch` | 发起范围设置 | `save-flow-switch` 参数参考 |
| `verify` | 问题节点列表 | 修复流程配置 |

## 注意事项

- `list-groups` **无参数**，返回全部分组；`list` **必须传 --group-id**，返回某分组下的审批类型。两者不可互换
- `manage get` 查的是审批**类型配置**（模板/表单/流程），参数 `--setting-id`；**不要**和 `detail get`（查审批记录实例，参数 `--sid`）混淆
- `delete` 删除的是审批类型配置，不是审批记录，执行前务必确认
- `update --sub` 修改子模块时，不同 sub 对应的操作逻辑不同，不确定时先用 `get` 看当前配置
- `preview` 可用于验证配置效果，建议在启用前先预览
- `copy-flow` 会完整覆盖目标类型的审批流配置，执行前确认目标类型
- 系统预设类型（`type: 0`）部分字段可能不可编辑或不可删除
- 自动发现更多命令：`xrxs approval manage --help`

## 相关模块

- [list](./list.md) — 审批记录查询与操作（search / approve / reject / cancel）
- [detail](./detail.md) — 审批记录实例详情（查的是实例，不是类型配置）
- [role](./role.md) — 审批角色管理
- [entrust](./entrust.md) — 审批委托/代理管理
- [dismiss-handover](./dismiss-handover.md) — 离职交接方案配置
