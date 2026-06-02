# xrxs-cli 功能列表

## 全局能力

| 能力 | 说明 |
|------|------|
| 多输出格式 | `-f json` / `-f table` / `-f raw` |
| 字段筛选 | `--fields employeeName,title,statusName` |
| jq 过滤 | `--jq '.[] \| select(.status == 0)'` |
| 详细日志 | `-v, --verbose` |
| 干运行 | `--dry-run` |
| 跳过确认 | `-y, --yes` |
| 参数命名容错 | camelCase / snake_case / kebab-case 自动归一化（如 `--flowStepId` = `--flow-step-id`） |

---

## xrxs auth — 认证

| 命令 | 说明 |
|------|------|
| `auth login` | 登录薪人薪事，自动保存 session |
| 自动续期 | 从 macOS Keychain / Linux secret-tool 读取密码解密 cookies |

---

## xrxs approval list — 审批列表查询与操作

### 查询

| 命令 | 类型 | 说明 |
|------|------|------|
| `search` | 单 API | 搜索审批列表，支持 keyword / status / flowGroupId / 分页 |
| `admin-search` | 单 API | 管理/关联审批列表 |
| `admin-search-all` | 单 API | 审批全局搜索（跨所有分组） |
| `speedy-search-type` | 单 API | 猜你想搜下拉选项 |
| `filter-num` | 单 API | 各状态审批数量统计 |
| `open-status` | 单 API | 审批开启状态 |
| `search-employee` | 单 API | 搜索部门/员工（@功能） |

### 单条操作

| 命令 | 类型 | 说明 |
|------|------|------|
| `approve` | 单 API | 通过审批 |
| `reject` | 单 API | 驳回审批 |
| `cancel` | 单 API | 撤销审批 |
| `forward` | 单 API | 转交审批 |
| `urge` | 单 API | 催办审批 |
| **`safe-urge`** | **Pipeline** | **安全催办：test-urge → urge（先检查再执行）** |
| `test-urge` | 单 API | 检查是否可催办 |

### 批量操作

| 命令 | 类型 | 说明 |
|------|------|------|
| `batch-approve` | 单 API | 批量通过/驳回 |
| `batch-forward` | 单 API | 批量转交 |
| `batch-urge` | 单 API | 批量催办 |
| **`safe-batch-urge`** | **Pipeline** | **安全批量催办：test-batch-urge → batch-urge** |
| `test-batch-urge` | 单 API | 检查是否可批量催办 |

### 导出

| 命令 | 类型 | 说明 |
|------|------|------|
| `export` | 单 API | 导出审批列表 |
| `export-check` | 单 API | 导出前检查权限 |

### 其他

| 命令 | 类型 | 说明 |
|------|------|------|
| `get-last-sign` | 单 API | 获取管理员上次签名 |
| `change-prove-status` | 单 API | 领取/更改证明状态 |
| `cancel-handover` | 单 API | 撤销离职交接 |
| `get-confirm-date` | 单 API | 获取离职确认日期 |
| `transfer-approval` | 单 API | 转交离职交接审批 |

---

## xrxs approval detail — 审批详情

| 命令 | 类型 | 说明 |
|------|------|------|
| `get` | 单 API | 获取审批详情（支持 type=1 详情 / type=2 打印） |
| `path` | 单 API | 获取审批流程路径/历史 |
| `print` | 单 API | 获取审批打印数据 |
| `edit-history` | 单 API | 获取审批编辑历史 |
| `save-page-edit` | 单 API | 保存审批页面编辑 |
| `overtime-hour` | 单 API | 获取加班时间信息 |
| `back-process-step` | 单 API | 获取可退回的审批节点列表 |
| `save-back-process-step` | 单 API | 执行审批退回操作 |
| `increase-process-step` | 单 API | 执行加签操作 |
| `open-mask-field` | 单 API | 掩码字段获取明文 |
| `check-form-required` | 单 API | 审批通过前校验必填项 |
| `check-biz-field` | 单 API | 审批业务字段校验 |
| `handover-detail` | 单 API | 获取离职交接审批详情 |

---

## xrxs approval manage — 审批类型管理

### 分组管理

| 命令 | 类型 | 说明 |
|------|------|------|
| `list-groups` | 单 API | 列出所有审批分组 |
| **`list-all-types`** | **Pipeline + Fan-out** | **一键拉取所有分组及其审批类型，并发度 4，单分组失败跳过** |
| `list` | 单 API | 按分组列出审批类型 |
| `save-group` | 单 API | 新增/重命名审批分组 |
| `remove-group` | 单 API | 删除审批分组 |

### 类型配置

| 命令 | 类型 | 说明 |
|------|------|------|
| `get` | 单 API | 获取审批类型配置详情 |
| `get-init` | 单 API | 获取新增审批类型时的初始数据 |
| `create` | 单 API | 创建审批类型 |
| `update` | 单 API | 更新审批类型配置 |
| `delete` | 单 API | 删除审批类型 |
| `toggle` | 单 API | 启用/停用审批类型 |
| `check-name` | 单 API | 检查审批类型名称是否重复 |

### 流程与表单

| 命令 | 类型 | 说明 |
|------|------|------|
| `get-flow-detail` | 单 API | 获取审批流详情 |
| `preview` | 单 API | 预览审批表单 |
| `copy-flow` | 单 API | 复制审批流 |
| `verify-flow` | 单 API | 校验审批流重复节点 |

### 关联与权限

| 命令 | 类型 | 说明 |
|------|------|------|
| `list-permission` | 单 API | 获取有权限的审批类型列表 |
| `list-related` | 单 API | 获取关联审批类型列表 |
| `get-charge-setting` | 单 API | 获取审批流分支条件的请假类型 |
| `get-print-hide-fields` | 单 API | 获取审批打印可隐藏字段 |

---

## Pipeline 能力

| 特性 | 说明 |
|------|------|
| 多步骤串联 | 按 Steps 顺序执行，前一步输出可被后一步引用 |
| 跨步骤数据传递 | `{{.steps.<stepId>.<field>}}` 模板语法 |
| 条件执行 | `condition` 字段，Go template 表达式 |
| Fan-out 并发 | `fanOut.source` 指定数组源，`concurrency` 控制并发度 |
| Fan-out 错误处理 | `onError: "skip"` 跳过失败项 / `"fail"` 终止 Pipeline |
| Fan-out 输出定制 | `itemKey` / `resultKey` 自定义输出结构 |
| 响应解包 | `unwrap` 字段按路径提取嵌套数据 |

### Pipeline 命令一览

| 命令 | 步骤 | 特性 |
|------|------|------|
| `manage list-all-types` | fetchGroups → fan-out fetchTypes | 并发 4，skip 错误 |
| `list safe-urge` | test-urge → urge | 自然错误传播，无需显式 condition |
| `list safe-batch-urge` | test-batch-urge → batch-urge | 同上 |

---

## 系统命令

| 命令 | 说明 |
|------|------|
| `xrxs version` | 显示版本信息 |
| `xrxs upgrade` | 自动升级到最新版本 |
| `xrxs uninstall` | 卸载 xrxs CLI |
| `xrxs --help` | 查看全局帮助 |

---

## 技术架构

| 特性 | 说明 |
|------|------|
| 命令生成 | Schema JSON 驱动，`cli.Builder` 自动生成 Cobra 命令 |
| Schema 加载 | 环境变量 > 嵌入二进制 > `config/schemas/` 回退 |
| HTTP 客户端 | Session + CSRF Token 自动管理 |
| 输出格式化 | `output.WriteCommandPayload` 统一处理 JSON/Table/Raw + fields + jq |
| 参数命名 | camelCase API 参数 → kebab-case CLI flag 自动转换 |
