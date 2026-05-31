# 错误码说明

审批模块错误参考 + 排查流程。Agent 遇到错误时查阅此文档。

## 错误三层分类

错误发生在三个层级，需区分处理：

| 层级 | 来源 | 典型输出 | Agent 行为 |
|------|------|----------|------------|
| CLI 参数层 | Cobra flag 校验失败 | `required flag(s) "group-id" not set` | 补传参数后重试 |
| 网络传输层 | HTTP 连接/超时 | `request: connection refused` | 至多重试 1 次，仍失败则报告用户 |
| API 业务层 | 服务端返回错误码 | `API 错误 [code]: message` | 按下方分类处理 |

---

## 一、CLI 参数层错误

Cobra 在命令执行前自动校验 flag，错误直接输出到 stderr：

```
Error: required flag(s) "group-id" not set
Usage:
  xrxs approval manage list [flags]
```

**Agent 行为**: 这是参数缺失，不是 API 错误。读取 `--help` 确认必填 flag，补传后立即重试。

**常见场景**:

| 命令 | 必填 flag | 缺失时的提示 |
|------|-----------|-------------|
| `manage list` | `--group-id` | "manage list 需要指定分组ID，我先查一下有哪些分组..." → 先调 `list-groups` |
| `manage get` | `--setting-id` | "manage get 需要指定审批类型ID，请提供 setting-id" |
| `detail get` | `--sid` | "detail get 需要指定审批记录 SID，请提供 sid" |
| `list approve/reject/cancel/forward` | `--sid` | "需要指定要操作的审批记录 SID" |

---

## 二、网络传输层错误

```
Error: request: dial tcp: connect: connection refused
Error: request: context deadline exceeded (Client.Timeout exceeded)
```

**Agent 行为**:
1. 加 `--verbose` 重试一次（排除临时抖动）
2. 仍失败 → 报告用户，建议检查：`--base-url` 是否正确、网络是否可达、服务是否正常运行
3. 可用 `xrxs auth status` 快速验证连通性

---

## 三、API 业务层错误

### 返回格式

```json
{"code": "AUTH_TOKEN_EXPIRED", "status": false, "message": "登录已过期", "data": null}
```

终端输出: `Error: API 错误 [AUTH_TOKEN_EXPIRED]: 登录已过期`

### 3.1 可自行修复（Agent 处理后重试）

| 错误特征 | 原因 | Agent 行为 |
|----------|------|------------|
| `required flag(s) "X" not set` | Cobra 层参数缺失 | 补传 flag 重试 |
| `InvalidParameter` + flag 名 | 参数名格式错误（如用了 camelCase） | 改用 kebab-case 重试 |
| `InvalidParameter` + 字段不存在 | `--fields` 指定的字段名不对 | 先用 `-f json` 看实际字段名 |
| `--jq` 表达式报错 | jq 语法错误 | 检查 jq 表达式后重试 |

### 3.2 需追问用户后重试

| 错误特征 | 原因 | Agent 行为 |
|----------|------|------------|
| `InvalidParameter` + SID/ID 无效 | 用户提供的 ID 不存在 | 告知用户该 ID 查不到，建议列出可选记录让用户重选 |
| 参数值语义错误（如审批意见为空） | 缺少用户主观输入 | 追问用户补充（comment、驳回原因等） |

追问示例：
> 找不到 SID 为 `abc123` 的审批记录。要我列出你当前的待审批记录吗？

### 3.3 需用户介入（Agent 不重试）

| 错误码 | 终端输出示例 | Agent 行为 |
|--------|-------------|------------|
| `AUTH_TOKEN_EXPIRED` | `API 错误 [AUTH_TOKEN_EXPIRED]: 登录已过期` | 告知用户登录过期，给出重登命令 |
| `PermissionDenied` | `API 错误 [PermissionDenied]: 无权限` | 告知用户无权限，解释具体是哪个操作被拒 |
| 资源不存在 | `API 错误 [...]: XXX 不存在` | 告知用户资源不存在，确认是否 ID 写错 |

**认证过期时的提示模板**:
> 登录已过期。请重新登录：
> ```bash
> xrxs auth login --base-url <你的服务地址>
> ```
> 登录完成后告诉我，我继续帮你操作。

**权限不足时的提示模板**:
> 当前账号无权执行此操作（[具体操作]）。请检查审批权限配置。

---

## 通用排查流程

```
1. 看错误来源
   ├─ "required flag(s)" → CLI 参数层 → 补传参数重试
   ├─ "request:" / "connection" / "timeout" → 网络层 → 重试 1 次，仍失败报告用户
   └─ "API 错误 [code]" → API 业务层 → 按下表判断
2. API 业务层判断
   ├─ AUTH_TOKEN_EXPIRED → 引导用户重新登录
   ├─ PermissionDenied → 告知权限不足
   ├─ InvalidParameter + 参数格式 → 修正后重试
   ├─ InvalidParameter + 值语义 → 追问用户
   └─ 其他 → 报告用户，不要自行编造解决方案
3. 加 --verbose 重试后仍不确定时 → 报告完整错误信息给用户
```
