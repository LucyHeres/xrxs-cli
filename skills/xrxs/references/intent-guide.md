# 意图路由指南

当用户请求难以判断归属哪个模块时，参考本指南。

## 易混淆场景快速对照表

| 用户说... | 真实意图 | 应该用 | 不要用 | 理由 |
|-----------|----------|--------|--------|------|
| "查看有哪些审批分组" | 查分组结构 | `approval manage list-groups` | `approval manage list` | 查分组用 list-groups；list 是查分组下的类型 |


---

## 典型场景详解

### 1. 审批记录 vs 审批类型配置 — 实例 vs 模板

这是最容易混淆的一对。审批系统有两个核心概念：

**审批记录（实例）** — 某个人提交的某条具体审批，有 sid
- 查询 → `approval list search`
- 详情 → `approval detail get --sid <sid>`
- 操作（通过/驳回/转交/撤销）→ `approval list approve/reject/forward/cancel`

**审批类型（配置/模板）** — 审批的表单设计、流程步骤设置，有 flowSettingId
- 查分组 → `approval manage list-groups`
- 查类型列表 → `approval manage list --group-id <id>`
- 查类型配置 → `approval manage get --setting-id <id>`
- 增删改启停 → `approval manage create/delete/update/toggle`

**判断关键**：用户给了 sid（审批编号）→ 实例操作；提到"配置/表单/流程/类型/分组" → 类型管理

### 2. 审批角色 vs 审批委托

**用 `approval role` 的场景**：
- "查看审批角色列表" — 角色管理
- 用户提到"角色"、"审批人角色"、"部门负责人"

**用 `approval entrust` 的场景**：
- "查看审批委托" — 委托/代理管理
- "张三委托李四审批" — 委托关系
- 用户提到"委托"、"代理"、"代为审批"

**判断关键**：提到"角色" → role；提到"委托/代理" → entrust

---

## 跨模块工作流路由

以下场景需要多个模块配合完成，注意上下文传递顺序。

### 查看并处理待审批（approval list → approval detail → approval list）

用户说"帮我看看有什么要批的"：

```bash
# 1. 查看待审批列表 — 提取 sid
xrxs approval list search --status 0 -f table

# 2. 查看某条审批详情（确认内容后决定操作） — 提取审批内容
xrxs approval detail get --sid <sid> --format json

# 3. 操作（意见从用户输入提取）
xrxs approval list approve --sid <sid> --comment "<审批意见>"
```

### 浏览全部分组和审批类型（approval manage → approval manage）

用户说"看看系统里有哪些审批类型"：

```bash
# 1. 先列全部分组 — 提取 groupId
xrxs approval manage list-groups --format json

# 2. 逐个分组查类型 — 提取 flowSettingId
xrxs approval manage list --group-id 2 --format json

# 3. 查看某类型配置详情
xrxs approval manage get --setting-id <setting_id> --format json
```
