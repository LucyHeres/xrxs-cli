# 证明与打印 (approval proof)

模块入口：`xrxs approval proof`

## 意图判断

用户说"查看证明信息/证明详情":
- 信息 → `info --sid <sid>`

用户说"证明设置/打印设置":
- 设置 → `setting --sid <sid>`

## 命令总览

### 查看证明信息

```
Usage:
  xrxs approval proof info [flags]
Example:
  xrxs approval proof info --sid <sid>
  xrxs approval proof info --sid <sid> --format json
Flags:
      --sid string   审批记录 ID（必传）
```

### 证明设置

```
Usage:
  xrxs approval proof setting [flags]
Example:
  xrxs approval proof setting --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

## 核心工作流

### 工作流: 打印审批证明

```bash
# 1. 查看证明信息
xrxs approval proof info --sid <sid> --format json

# 2. 配置打印设置
xrxs approval proof setting --sid <sid>
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `info` | 证明内容 | 向用户展示或打印 |

## 注意事项

- `--sid` 来自 `approval list search` 返回的 sid 字段
- 自动发现更多命令：`xrxs approval proof --help`

## 相关模块

- [approval-list](./approval-list.md) — 审批列表查询
- [approval-detail](./approval-detail.md) — 审批详情（含打印视图）
