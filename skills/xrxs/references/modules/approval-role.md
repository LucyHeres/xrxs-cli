# 审批角色 (approval role)

模块入口：`xrxs approval role`

## 意图判断

用户说"查看审批角色/有哪些角色/角色列表":
- 列表 → `list`

用户说"搜一下某角色/查找角色":
- 搜索 → `list --keyword <关键词>`

用户说"查看角色详情/这个角色有哪些成员":
- 详情 → `get --role-id <id>`

用户说"创建角色/新增审批角色":
- 创建 → `create --name <名称>`

用户说"修改角色/更新角色/给角色加成员":
- 更新 → `update --role-id <id>`

用户说"删除角色/移除角色":
- 删除 → `delete --role-id <id>`

用户说"这个角色在哪些审批里用到了":
- 使用情况 → `usage --role-id <id>`

关键区分: `role`(审批角色，如部门负责人/HR) vs `entrust`(审批委托/代理，A 委托 B 代为审批)

## 命令总览

### 列出审批角色

```
Usage:
  xrxs approval role list [flags]
Example:
  xrxs approval role list
  xrxs approval role list --keyword "HR" -f table
  xrxs approval role list --format json
Flags:
      --keyword string  搜索关键词
```

### 获取角色详情

```
Usage:
  xrxs approval role get [flags]
Example:
  xrxs approval role get --role-id <id>
  xrxs approval role get --role-id <id> --format json
Flags:
      --role-id string   角色 ID（必传）
```

### 创建角色

```
Usage:
  xrxs approval role create [flags]
Example:
  xrxs approval role create --name "部门主管"
  xrxs approval role create --name "HR 审批" --members '{"userIds":["id1","id2"]}'
Flags:
      --name string     角色名称（必传）
      --members string   角色成员 JSON
```

### 更新角色

```
Usage:
  xrxs approval role update [flags]
Example:
  xrxs approval role update --role-id <id> --name "新角色名"
  xrxs approval role update --role-id <id> --members '{"userIds":["id1","id2"]}'
Flags:
      --role-id string   角色 ID（必传）
      --name string      新名称
      --members string   新成员 JSON
```

### 删除角色

> **⚠️ 不可逆，执行前须向用户确认。**

```
Usage:
  xrxs approval role delete [flags]
Example:
  xrxs approval role delete --role-id <id>
Flags:
      --role-id string   角色 ID（必传）
```

### 查看角色使用情况

```
Usage:
  xrxs approval role usage [flags]
Example:
  xrxs approval role usage --role-id <id>
  xrxs approval role usage --role-id <id> --format json
Flags:
      --role-id string   角色 ID（必传）
```

## 核心工作流

### 工作流 1: 创建并配置审批角色

```bash
# 1. 创建角色 — 提取 roleId
xrxs approval role create --name "新角色" --format json

# 2. 添加成员
xrxs approval role update --role-id <id> --members '{"userIds":["uid1","uid2"]}'

# 3. 验证 — 查详情
xrxs approval role get --role-id <id> --format json
```

### 工作流 2: 查看角色在哪些审批中生效

```bash
# 1. 查使用情况
xrxs approval role usage --role-id <id> --format json
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `list` | `roleId` | get / update / delete / usage 的 --role-id |
| `list` | 角色名称 | 向用户展示，确认操作目标 |
| `create` | `roleId` | update 追加成员 |
| `get` | 成员列表 | 判断是否需要更新 |
| `usage` | 关联的审批类型 | 判断删除是否会影响现有审批 |

## 注意事项

- `delete` 删除的是审批角色，执行前先用 `usage` 确认该角色是否还在被审批流程引用
- `--members` 是 JSON 字符串，格式需正确
- 区别于 `entrust`(审批委托/代理)，两者命名相近但职责不同
- 自动发现更多命令：`xrxs approval role --help`

## 相关模块

- [approval-entrust](./approval-entrust.md) — 审批委托管理（注意区分角色 vs 委托）
- [approval-manage](./approval-manage.md) — 审批类型管理（角色在审批流程设置中使用）
