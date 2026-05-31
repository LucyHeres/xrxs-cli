你怎么又在改claude的skill？我# 审批类型管理 (approval manage)

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

**用户只说"列出审批类型/有哪些审批类型"（未提分组）**:
- 先调 `list-groups` 拿所有分组 ID，再逐个 `list --group-id <id>`，最后汇总展示

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

## 核心工作流

### 工作流 1: 浏览全部分组和审批类型

```bash
# 1. 列全部分组
xrxs approval manage list-groups --format json

# 2. 从返回中提取所有 groupId，逐个查类型
xrxs approval manage list --group-id 0 --format json   # 员工审批
xrxs approval manage list --group-id 1 --format json   # 工资社保审批
xrxs approval manage list --group-id 2 --format json   # 考勤审批
# ... 依此类推，按 list-groups 返回的 groupId 遍历
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

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `list-groups` | `groupId` | `list --group-id` |
| `list` | `flowSettingId` | `get / update / delete / toggle / preview --setting-id` |
| `list` | `name` | 向用户展示审批类型名称，确认操作目标 |
| `list` | `openStatus` | 判断是否需要启停 |
| `create` | `flowSettingId` | `update / delete / toggle / preview --setting-id` |
| `get` | 表单/流程配置 | `update --sub form / --sub flow` |

## 注意事项

- `list-groups` **无参数**，返回全部分组；`list` **必须传 --group-id**，返回某分组下的审批类型。两者不可互换
- `manage get` 查的是审批**类型配置**（模板/表单/流程），参数 `--setting-id`；**不要**和 `detail get`（查审批记录实例，参数 `--sid`）混淆
- `delete` 删除的是审批类型配置，不是审批记录，执行前务必确认
- `update --sub` 修改子模块时，不同 sub 对应的操作逻辑不同，不确定时先用 `get` 看当前配置
- `preview` 可用于验证配置效果，建议在启用前先预览
- 系统预设类型（`type: 0`）部分字段可能不可编辑或不可删除
- 自动发现更多命令：`xrxs approval manage --help`

## 相关模块

- [approval-list](./approval-list.md) — 审批记录查询与操作（search / approve / reject / cancel）
- [approval-detail](./approval-detail.md) — 审批记录实例详情（查的是实例，不是类型配置）
- [approval-role](./approval-role.md) — 审批角色管理
- [approval-entrust](./approval-entrust.md) — 审批委托/代理管理
