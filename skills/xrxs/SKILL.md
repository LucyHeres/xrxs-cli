---
name: xrxs
description: 薪人薪事 HR SaaS 系统命令行工具。当用户需要查询审批列表、查看审批详情、审批操作（通过/驳回/转交/催办）、管理审批类型、审批角色、审批委托时使用。
cli_version: ">=0.1.0"
---

# 薪人薪事 CLI (xrxs)

通过 `xrxs` 命令管理薪人薪事 HR SaaS 系统。

## 安装

```bash
# macOS / Linux
curl -fsSL https://gh-proxy.org/https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.sh | sh

# Windows PowerShell
irm https://github.com/LucyHeres/xrxs-cli/releases/latest/download/install.ps1 | iex
```

## 全局参数

所有命令继承以下全局参数：

| Flag | 默认值 | 说明 |
|------|--------|------|
| `--base-url` | 环境变量 `XRXS_BASE_URL` 或配置文件 | API 服务地址 |
| `-f, --format` | `json` | 输出格式: `json` / `table` / `raw` |
| `--jq` | - | jq 表达式过滤输出 |
| `--fields` | - | 只输出指定字段，逗号分隔 |
| `-v, --verbose` | `false` | 显示详细请求日志 |
| `--dry-run` | `false` | 预览操作但不执行 |
| `-y, --yes` | `false` | 跳过确认提示 |

## 严格禁止 (NEVER DO)

- 不要使用 xrxs 命令以外的方式操作（禁止 curl、HTTP API、浏览器）
- 不要编造 ID，必须从命令返回中提取
- 不要猜测参数值，操作前必须先查询确认
- 不要线性扫描命令列表到第一个匹配就停，必须全量对比候选后再选

## 严格要求 (MUST DO)

- 大规模输出时必须加 `--fields` 或 `--jq` 减少 token 消耗
- 所有命令默认使用 JSON 输出（程序化解析），展示给用户时使用 `-f table`
- 危险操作（删除、批量审批）必须先向用户确认
- 审批意见（`--comment`）、驳回原因、委托原因等用户主观内容，必须从用户输入中提取；用户未提供时**必须追问**，严禁编造默认值

## 核心流程

作为智能助手，在选择命令前必须严格遵循以下流程：

1. **意图分类** — 判断用户核心动词/动作：是想查记录实例？管理审批类型？查看分组结构？审批操作？
2. **歧义处理** — 当意图模糊或涉及易混淆命令时，**严禁猜测**。先读 [intent-guide.md](references/intent-guide.md) 对照意图路由表和成对区分
3. **精准命令映射** — 意图明确后，参考下方模块总览和 [详细参考](#详细参考按需读取) 中对应产品的 reference 文件选择正确命令
4. **执行** — 命令不确定时先用 `--help` 验证参数和用法

## 命令发现

参数名称和用法以 `--help` 输出为准：

```bash
# 查看子命令用法
xrxs approval list search --help
xrxs approval manage get --help

# 查看所有可用子命令
xrxs approval --help
```

本文档中的 flag 列表仅作参考，如果与 `--help` 输出冲突，以 `--help` 为准。

## 危险操作确认

以下操作为不可逆或高影响操作，执行前**必须先向用户展示操作摘要并获得明确同意**，同意后才加 `--yes` 执行。

| 模块 | 命令 | 说明 |
|------|------|------|
| 审批列表 | `approve` / `reject` / `cancel` | 审批操作不可撤销，需明确 comment |
| 审批列表 | 批量 approve/reject/cancel | 影响多条审批记录 |
| 审批管理 | `manage delete` | 删除审批类型配置，关联数据受影响 |
| 审批管理 | `manage toggle` | 启用/停用审批类型 |
| 审批角色 | `role delete` | 删除角色配置 |
| 审批委托 | `entrust cancel` | 取消委托关系 |

### 确认流程

```
Step 1 → 展示操作摘要（操作类型 + 目标对象 + 影响范围）
Step 2 → 用户明确回复确认（如 "确认" / "好的"）
Step 3 → 加 --yes 执行命令
```

## 错误处理

错误发生在三个层级，需区分处理：

1. **CLI 参数层** — stderr 包含 `required flag(s)` → 补传缺失的 flag 重试
2. **网络层** — stderr 包含 `connection` / `timeout` → 加 `--verbose` 重试 1 次，仍失败则报告用户并建议检查 `--base-url` 和服务连通性
3. **API 业务层** — stderr 包含 `API 错误 [code]`：
   - 认证过期 → 引导用户执行 `xrxs auth login --base-url <URL>`
   - 权限不足 → 报告用户，禁止重试
   - 参数格式错误 → 修正后重试
   - 参数值语义错误（如 ID 不存在）→ 追问用户

如果仍无法解决错误或者定位不到错误原因，报告完整错误信息给用户，禁止自行尝试替代方案
详细错误码和排查流程见 [error-codes.md](references/error-codes.md)，全局配置见 [global-reference.md](references/global-reference.md)。

## 模块总览

| 模块 | 命令入口 | 用途 |
|------|----------|------|
| 认证 | `xrxs auth` | 登录/退出/状态 |
| 审批列表 | `xrxs approval list` | 查询审批记录、通过/驳回/转交/撤销/催办/导出 |
| 审批详情 | `xrxs approval detail` | 查看审批实例详情、流转路径、打印、编辑历史 |
| 审批管理 | `xrxs approval manage` | 管理审批类型配置、审批分组 |
| 审批角色 | `xrxs approval role` | 审批角色 CRUD |
| 审批委托 | `xrxs approval entrust` | 审批委托/代理管理 |
| 证明打印 | `xrxs approval proof` | 证明与打印 |

## 典型用法

```bash
# 查看待审批列表（表格形式）
xrxs approval list search --status 0 -f table

# 只看关键字段
xrxs approval list search --status 0 --fields employeeName,title,addDate -f table

# 批量查看所有待审批
xrxs approval list search --status 0 --page-size 100 -f table

# 一键查看所有分组及其审批类型（Pipeline 命令，推荐）
xrxs approval manage list-all-types -f table

# 查看具体某个分组下的审批类型
xrxs approval manage list --group-id 2 -f table

# 查看审批详情
xrxs approval detail get --sid <sid>

# 催办(先检查再执行，推荐)
xrxs approval list safe-urge --sid <sid>

# 通过审批
xrxs approval list approve --sid <sid> --comment "同意"
```

### Pipeline 命令

部分命令通过 Pipeline 机制将多个 API 调用串联为一条命令：

| 命令 | 步骤 | 说明 |
|------|------|------|
| `manage list-all-types` | `fetchGroups` → fan-out `fetchTypes` | 一键拉取所有分组+类型，并发度 4，单分组失败自动跳过 |
| `list safe-urge` | `test-urge` → `urge` | 先检查是否可催办，可则自动执行；不可则终止并返回原因 |
| `list safe-batch-urge` | `test-batch-urge` → `batch-urge` | 同上，批量场景 |

## 详细参考（按需读取）

- [references/global-reference.md](references/global-reference.md) — 全局参考（认证、环境变量、全局 flag、输出格式）
- [references/intent-guide.md](references/intent-guide.md) — 意图路由指南（常见意图路由表 + 易混淆命令成对区分 + 典型混淆场景详解）
- [references/error-codes.md](references/error-codes.md) — 错误码说明与排查流程
- [references/modules/approval-list.md](references/modules/approval-list.md) — 审批列表查询与操作
- [references/modules/approval-detail.md](references/modules/approval-detail.md) — 审批详情
- [references/modules/approval-manage.md](references/modules/approval-manage.md) — 审批类型管理（含 list vs list-groups 区分）
- [references/modules/approval-role.md](references/modules/approval-role.md) — 审批角色
- [references/modules/approval-entrust.md](references/modules/approval-entrust.md) — 审批委托
- [references/modules/approval-proof.md](references/modules/approval-proof.md) — 证明与打印

## 扩展

未来会支持考勤、人事、薪酬等模块。命令通过 Schema JSON 配置自动生成。
