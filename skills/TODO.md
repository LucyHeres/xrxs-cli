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
