# 全局参考

## 认证

```bash
# 首次登录 (SSO 用户名+密码)
xrxs auth login --base-url https://your-instance.example.com

# 查看状态
xrxs auth status

# 退出
xrxs auth logout
```

登录后凭证加密存储在 `~/.xrxs/cookies.enc`，日常使用无需重复登录。

### 认证失败处理

- 命令返回 `AUTH_TOKEN_EXPIRED` / "登录已过期" / "未登录" → 执行 `xrxs auth login --base-url <URL>` 重新登录
- 如果 `xrxs auth status` 显示已登录但命令仍然报认证错误 → 可能是 base-url 不一致，检查 `~/.xrxs/config.json` 中的配置

### CI/CD / Headless 环境

先交互式登录一次（凭证持久化到 `~/.xrxs/cookies.enc`），后续非交互式调用自动使用已保存的凭证。

## 全局标志

所有命令继承以下全局标志：

| 标志 | 短名 | 说明 | 默认 |
|------|:---:|------|------|
| `--base-url` | | API 服务地址 | 环境变量 `XRXS_BASE_URL` 或配置文件 |
| `--format` | `-f` | 输出格式: `json` / `table` / `raw` | `json` |
| `--jq` | | jq 表达式过滤输出 | 无 |
| `--fields` | | 筛选取字段，逗号分隔 | 无 |
| `--verbose` | `-v` | 显示详细请求日志 | `false` |
| `--dry-run` | | 预览操作但不执行 | `false` |
| `--yes` | `-y` | 跳过确认提示 | `false` |

## 输出格式

### `--format json` (默认，机器可读)

JSON 结构化输出，适合 Agent 解析。

### `--format table` (人类可读)

表格形式，适合直接呈现给用户。中文列名自动映射。

### 减少 Token 消耗

- 大规模输出时用 `--fields` 只取需要的字段
- 需要过滤/变换时用 `--jq` 表达式
- 不确定有哪些字段时先用 `-f json` 看完整结构

## 环境变量

| 变量 | 说明 |
|------|------|
| `XRXS_BASE_URL` | 默认 API 服务地址 |
| `XRXS_SCHEMA_DIR` | 开发用：从文件系统加载 schema（覆盖内嵌 schema） |
| `XRXS_CONFIG_DIR` | 覆盖默认配置目录（默认 `~/.xrxs`） |
