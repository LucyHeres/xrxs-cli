# TODO — 暂缓事项

以下事项经架构分析后认为暂时不需要实施，待条件成熟后再做。

---

## 1. `xrxs skill setup` CLI 命令

**原因**: 当前通过 `install.sh` 安装 skill 已满足需求（已修复为递归复制整个 `skills/xrxs/` 目录）。DWS 的 `dws skill setup` 支持交互式选择、mono/multi 模式切换、多源路径解析等能力，但 XRXS 目前只有一个产品模块（审批），不需要这些复杂度。

**触发条件**: 当产品模块达到 3 个以上，或用户反馈需要更新/修复 skill 而不想重跑完整安装脚本时，再实现此命令。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/internal/app/skill_setup.go`

---

## 2. `scripts/` 脚本目录

**原因**: DWS 的 `scripts/` 封装了需要多次 CLI 调用的复杂操作（如批量导入、导出消息历史、自动预约会议等）。XRXS 审批模块的大部分操作本身就是单步 CLI 调用，不需要脚本编排。

**触发条件**: 当出现需要 3 步以上 CLI 调用才能完成的常用场景时（如"导出某时间段全部审批记录并生成报表"），再添加 Python 脚本。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/skills/mono/scripts/`

---

## 3. `references/field-rules.md`

**原因**: DWS 的 AI 表格（aitable）有复杂的字段类型规则（单选/多选/公式/附件等），需要专门的字段规则文档。XRXS 审批模块的字段主要是文本和枚举，类型简单。

**触发条件**: 当审批模块出现复杂自定义字段（如级联选择、动态表单、附件上传）且 AI agent 频繁因字段格式错误导致命令失败时，再编写此文件。

---

## 4. 多产品意图路由指南

**原因**: 当前只有审批一个模块，`intent-guide.md` 中的成对区分已足够。DWS 有 20 个产品，需要关键词→产品决策树。XRXS 目前不需要。

**触发条件**: 当考勤、人事、薪酬等新模块上线后，在 `intent-guide.md` 中增加产品级路由表。

---

## 5. Windows 安装脚本 (install.ps1) 增加 Skill 安装

**原因**: `install.ps1` 目前只安装二进制文件，没有安装 skill 到 AI agent 目录。Windows 用户无法使用 skill 功能。

**触发条件**: 有 Windows 用户反馈 skill 不生效时修复。修复方式参照已更新的 `install.sh`，递归复制 `skills/xrxs/` 到各 agent 目录。

---

## 6. Mono/Multi 双模式支持

**原因**: DWS 支持单 skill（mono）和多 skill（multi）两种部署模式。XRXS 只有一个产品，不需要 multi 模式。

**触发条件**: 当产品模块达到 5 个以上且用户希望按需安装部分模块时，实现 multi 模式。

---

## 7. 结构化错误体系 + 退出码

**原因**: 当前所有错误统一 `exit code 1`，无法区分是登录过期、参数错误还是网络超时。DWS 有 5 种错误类别（API/Auth/Validation/Discovery/Internal），每种固定退出码，支持机器可读 JSON + 人类可读两种输出格式，附带 Hint/Actions 等恢复建议。

**触发条件**: 当 CLI 被脚本/CI/Agent 自动化调用，需要根据失败原因走不同处理分支（重试/重新登录/报错）时实施。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/internal/errors/errors.go`

---

## 8. Schema 远程发现 + 本地缓存

**原因**: 当前 Schema 通过本地 JSON 文件加载（开发）或嵌入二进制（生产），需要手动维护两份副本。DWS 通过 Market API 远程拉取 Schema，本地缓存到 `~/.dws/cache/`，支持 TTL 过期、原子写入、离线降级兜底。

**触发条件**: Schema 迁移到服务端、支持多个产品模块后实施。届时删掉 `internal/schema/embed.go` 和 symlink，替换为 HTTP loader。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/internal/discovery/service.go`, `internal/cache/store.go`

---

## 9. Flag → 参数 Normalizer 闭包管线

**原因**: 当前 `buildFormParams` 只做简单的 flag → URL value 映射。DWS 的 `buildOverrideBindings` 返回一个 8 步 Normalizer 管线：默认值注入（按类型强制）→ 环境变量回退 → Runtime Defaults → 转换函数 → MapsTo 路由 → OmitWhen → 嵌套结构 → Body 包装。可以替代很多 body template 中的重复逻辑。

**触发条件**: 当 body template 越来越复杂、出现大量重复的 `{{if .xxx}}...{{end}}` 条件判断时实施。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/internal/compat/dynamic_commands.go:564-890`

---

## 10. 输出自动检测（findDataList）

**原因**: 当前需要手动在 Schema 中配置 `"unwrap": "data.result.data"` 来定位响应中的数组。DWS 的 `findDataList` 自动按优先级搜索常见 key（value/items/results/data/list/records 等），table/csv 格式化器共享同一检测逻辑。

**触发条件**: 当 API 响应格式不统一、unwrap 配置过多且易出错时实施。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/internal/output/filter.go:89-121`

---

## 11. 测试 Fixture 模式

**原因**: 当前 `builder_test.go` 是手动 `httptest.NewServer` + 硬编码 handler，每个测试都要重写完整的 mock 逻辑。DWS 的 `mock_mcp.Server` + `Fixture` 结构体支持声明式定义 mock 行为（Registry/CLI/Detail/MCP 四层 mock），测试代码更简洁。

**触发条件**: 当测试文件超过 500 行、新增测试越来越繁琐时实施。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/test/mock_mcp/server.go`, `fixture.go`

---

## 12. Edition 扩展点系统

**原因**: 当前所有行为硬编码在 `root.go` 中。DWS 的 `edition.Hooks` 提供覆盖点：自定义配置目录、自定义 Token 持久化、静态服务器列表、自定义 HTTP Header、注册额外命令等。企业私有化部署通过外部注入，不修改核心代码。

**触发条件**: 出现私有化部署需求、需要支持不同客户的定制行为时实施。

**参考**: `/Users/liuxin/code/dingtalk-workspace-cli/pkg/edition/edition.go`
