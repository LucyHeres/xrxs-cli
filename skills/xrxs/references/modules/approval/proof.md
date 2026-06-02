# 证明与打印 (approval proof)

模块入口：`xrxs approval proof`

## 意图判断

用户说"查看证明信息/证明快照":
- 快照 → `snap --sid <sid>`

用户说"开具证明/开具记录/证明领取记录":
- 记录列表 → `records --sid <sid>`
- 添加记录 → `add-record --sid <sid> --record-status <int>`

用户说"打印审批/打印编辑":
- 打印编辑 → `edit-print --sid <sid>`
- 执行打印 → `print --sid <sid>`
- 生成打印记录 → `generate-print --sid <sid>`

用户说"打印设置/隐藏字段/打印哪些字段":
- 隐藏字段 → `hide-fields --setting-id <id>`
- 保存打印设置 → `save-print --data '<JSON>'`

关键区分: `proof print`(审批打印，生成打印视图) vs `detail print`(审批详情打印模式)

## 命令总览

### 获取证明快照

```
Usage:
  xrxs approval proof snap [flags]
Example:
  xrxs approval proof snap --sid <sid>
  xrxs approval proof snap --sid <sid> --format json
Flags:
      --sid string   审批记录 ID（必传）
```

返回：文件 URL、文件名、快照内容、模板类型。

### 证明开具记录

```
Usage:
  xrxs approval proof records [flags]
Example:
  xrxs approval proof records --sid <sid>
  xrxs approval proof records --sid <sid> --format json
Flags:
      --sid string   审批记录 ID（必传）
```

返回：recordId、recordStatus、createTime。

### 添加证明开具记录

```
Usage:
  xrxs approval proof add-record [flags]
Example:
  xrxs approval proof add-record --sid <sid> --record-status 1
Flags:
      --sid string           审批记录 ID（必传）
      --record-status int    记录状态（必传）
```

### 打印编辑

```
Usage:
  xrxs approval proof edit-print [flags]
Example:
  xrxs approval proof edit-print --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

### 审批打印

```
Usage:
  xrxs approval proof print [flags]
Example:
  xrxs approval proof print --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

### 生成打印记录

```
Usage:
  xrxs approval proof generate-print [flags]
Example:
  xrxs approval proof generate-print --sid <sid>
Flags:
      --sid string   审批记录 ID（必传）
```

### 打印隐藏字段

```
Usage:
  xrxs approval proof hide-fields [flags]
Example:
  xrxs approval proof hide-fields --setting-id 8483935
  xrxs approval proof hide-fields --setting-id 8483935 --format json
Flags:
      --setting-id int   审批类型设置 ID（必传）
```

### 保存打印设置

```
Usage:
  xrxs approval proof save-print [flags]
Example:
  xrxs approval proof save-print --data '<JSON>'
Flags:
      --data string   打印权限数据 JSON（必传）
```

## 核心工作流

### 工作流 1: 查看并打印证明

```bash
# 1. 查看证明快照
xrxs approval proof snap --sid <sid> --format json

# 2. 查看历史开具记录
xrxs approval proof records --sid <sid> --format json

# 3. 添加开具记录
xrxs approval proof add-record --sid <sid> --record-status 1

# 4. 打印
xrxs approval proof print --sid <sid>
```

### 工作流 2: 配置打印模板

```bash
# 1. 查看当前隐藏字段
xrxs approval proof hide-fields --setting-id <id>

# 2. 编辑打印设置
xrxs approval proof edit-print --sid <sid>

# 3. 保存打印设置
xrxs approval proof save-print --data '<JSON>'

# 4. 生成打印记录
xrxs approval proof generate-print --sid <sid>
```

## 上下文传递表

| 操作 | 从返回中提取 | 用于 |
|------|-------------|------|
| `snap` | `url` / `snapContent` | 向用户展示证明内容 |
| `snap` | `templateType` | 判断证明类型 |
| `records` | `recordId` / `recordStatus` | 跟踪开具状态 |
| `hide-fields` | 隐藏字段列表 | `save-print` 时参考 |
| `edit-print` | 打印编辑数据 | `save-print` 提交 |

## 注意事项

- `proof print` 和 `detail print` 是不同的接口，前者用于审批证明打印，后者用于审批详情打印视图
- `--sid` 来自 `approval list search` 返回的 sid 字段
- `hide-fields` 的 `--setting-id` 来自 `manage list` 的 flowSettingId，不是 sid
- 打印设置（save-print）提交的是 JSON 数据，结构需与 `edit-print` 返回一致
- 自动发现更多命令：`xrxs approval proof --help`

## 相关模块

- [list](./list.md) — 审批列表查询（获取 sid）
- [detail](./detail.md) — 审批详情（含详情打印视图）
- [manage](./manage.md) — 审批类型管理（获取 setting-id）
