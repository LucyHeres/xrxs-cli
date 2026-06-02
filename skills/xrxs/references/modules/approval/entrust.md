# 审批委托 (approval entrust)

模块入口：`xrxs approval entrust`

## 意图判断

用户说"查看委托列表/当前的审批委托/委托代理":
- 列表 → `list`

用户说"设置委托/委托某人审批/我不在的时候谁代审批":
- 创建 → `create --agent-id <id> --reason "<委托原因>"`（原因从用户输入提取，用户未提供时追问）

用户说"取消委托/撤销委托/不用代审批了":
- 取消 → `cancel --id <id>`

用户说"查看委托设置/委托开关":
- 设置 → `setting`

用户说"委托时可以选哪些审批类型/委托可选审批":
- 审批类型列表 → `type-list`

关键区分: `entrust`(审批委托/代理，A 委托 B 代为审批) vs `role`(审批角色，如部门负责人/HR)

## 命令总览

### 列出审批委托

```
Usage:
  xrxs approval entrust list [flags]
Example:
  xrxs approval entrust list
  xrxs approval entrust list --page 1 --page-size 20 -f table
  xrxs approval entrust list --format json
Flags:
      --page int       页码 (默认 1)
      --page-size int  每页条数 (默认 20)
```

### 创建委托

```
Usage:
  xrxs approval entrust create [flags]
Example:
  xrxs approval entrust create --agent-id <id> --start-time "<YYYY-MM-DD>" --end-time "<YYYY-MM-DD>" --reason "<委托原因>"
  xrxs approval entrust create --agent-id <id> --scope 0
Flags:
      --agent-id string    被委托人 ID（必传）
      --scope int          委托范围: 0=全部审批, 1=指定审批类型
      --start-time string  开始日期 (YYYY-MM-DD)
      --end-time string    结束日期 (YYYY-MM-DD)
      --reason string      委托原因
```

### 取消委托

> 取消后委托立即失效，执行前须向用户确认。

```
Usage:
  xrxs approval entrust cancel [flags]
Example:
  xrxs approval entrust cancel --id <id>
Flags:
      --id string   委托 ID（必传）
```

### 委托设置

```
Usage:
  xrxs approval entrust setting [flags]
Example:
  xrxs approval entrust setting
  xrxs approval entrust setting --status 1
Flags:
      --status int   状态: 0=关闭, 1=开启
```

### 获取审批类型列表（委托用）

```
Usage:
  xrxs approval entrust type-list [flags]
Example:
  xrxs approval entrust type-list
  xrxs approval entrust type-list --format json
Flags:
  无必填参数
```

返回 key/value 对象，key 为设置 ID，value 为审批名称。

## 核心工作流

### 工作流 1: 设置出差期间的审批委托

```bash
# 1. 创建委托
xrxs approval entrust create \
  --agent-id <委托人ID> \
  --start-time "2026-06-01" \
  --end-time "2026-06-07" \
  --reason "<委托原因>" \
  --format json

# 2. 验证 — 查看委托列表
xrxs approval entrust list -f table
```

### 工作流 2: 取消委托

```bash
# 1. 查当前委托 — 提取 id
xrxs approval entrust list --format json

# 2. 确认后取消
xrxs approval entrust cancel --id <id>
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `list` | `id` | `cancel --id` |
| `list` | `fromName` / `toName` | 向用户展示委托人和被委托人 |
| `list` | `flowName` | 委托关联的审批类型 |
| `list` | `startDate` / `endDate` / `status` | 判断委托是否在有效期 |
| `type-list` | `key` / `value` | `create` 时选择审批类型 |

## 注意事项

- 委托生效期间，被委托人可以代为审批，结束后权限自动回收
- `cancel` 会立即终止委托，执行前确认
- `--start-time` 和 `--end-time` 格式为 `YYYY-MM-DD`
- 区别于 `role`(审批角色)，两者命名相近但职责完全不同
- 自动发现更多命令：`xrxs approval entrust --help`

## 相关模块

- [role](./role.md) — 审批角色管理（注意区分角色 vs 委托）
- [list](./list.md) — 审批列表查询与操作
