# 审批管理端系统 API 接口文档

> 自动生成自 `xrxs-work-flow-springboot` 项目 Controller 层代码
> 基础路径: `/approve/service`
> 统一响应包装: `AjaxResult` { code, message, lanDic, success, data }
> 鉴权注解: `@ApproveAuth` (所有主业务 Controller 类级别)

---

## 一、审批操作 (FlowController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-pass-approve-step | 管理员通过一条审批 | POST /approve/service/ajax-pass-approve-step | POST | sid(string/是), flowStepId(string/是), comment(string/是), confirmDate(string/否/""), atEmployeeIds(string/否/""), sign(string/否/"") |
| 2 | ajax-rejected-approve-step | 管理员驳回一条审批 | POST /approve/service/ajax-rejected-approve-step | POST | sid(string/是), flowStepId(string/是), comment(string/是), atEmployeeIds(string/否/"") |
| 3 | ajax-close-approve-step | 管理员撤回一条审批 | POST /approve/service/ajax-close-approve-step | POST | sid(string/是), remark(string/是) |
| 4 | ajax-change-approve-step | 管理员转发一条审批 | POST /approve/service/ajax-change-approve-step | POST | sid(string/是), flowStepId(string/是), remark(string/是), approverId(string/是), nextApproverId(string/是) |
| 5 | ajax-batch-operation-approve-step | 管理员批量操作审批（通过/驳回） | POST /approve/service/ajax-batch-operation-approve-step | POST | data(string/是), type(int/是) |
| 6 | ajax-batch-change-approve-step | 管理员批量转发操作 | POST /approve/service/ajax-batch-change-approve-step | POST | data(string/是) |
| 7 | ajax-increase-process-step | 操作加签节点 | POST /approve/service/ajax-increase-process-step | POST | sid(string/是), stepNodeId(string/是), increaseType(int/是), employeeIds(string/是), flowPassType(int/是), remark(string/是) |
| 8 | ajax-change-prove-status | 证明审批-管理设置证明是否领取 | GET /approve/service/ajax-change-prove-status | GET | sid(string/是), proveStatus(int/是) |
| 9 | ajax-calculate-expression | 根据前端传过来的计算式返回计算结果 | POST /approve/service/ajax-calculate-expression | POST | customFields(string/是), sid(string/否/"") |
| 10 | ajax-get-back-process-step | 获取退回节点 | POST /approve/service/ajax-get-back-process-step | POST | sid(string/是) |
| 11 | ajax-back-process-step | 操作退回节点 | POST /approve/service/ajax-back-process-step | POST | sid(string/是), stepNodeId(string/是), remark(string/是), backStepNodeId(string/是), reApproval(int/是) |

---

## 二、审批列表与搜索 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-approve-list | 获取审批列表 | POST /approve/service/ajax-get-approve-list | POST | monthFilter(int/否/null), start(int/否/0), search[value](string/否/""), speedy[value](string/否/""), speedySearch(bool/否/false), order[0][dir](string/否/"desc"), order[0][field](string/否/"addDate") |
| 2 | ajax-get-admin-approve-list | 管理员获取审批列表(关联审批单) | POST /approve/service/ajax-get-admin-approve-list | POST | start(int/否/0), search[value](string/否/""), pageSize(int/否/自动) |
| 3 | ajax-get-admin-approve-list-all | 获取审批列表(审批中搜索全部审批) | POST /approve/service/ajax-get-admin-approve-list-all | POST | start(int/否/0), search[value](string/否/""), speedy[value](string/否/""), speedySearch(bool/否/false) |
| 4 | ajax-get-filter-num | 获取审批中/本月通过/本月驳回数量 | GET /approve/service/ajax-get-filter-num | GET | 无 |
| 5 | ajax-get-approve-speedy-search-type | 获取审批列表猜你想搜类型 | POST /approve/service/ajax-get-approve-speedy-search-type | POST | employeeName(string/是), employeeId(string/是) |
| 6 | ajax-get-related-approve-type-list | 获取权限范围内审批类型 | GET /approve/service/ajax-get-related-approve-type-list | GET | 无 |
| 7 | ajax-get-flow-list-with-permissions | 获取管理员权限筛选后的审批列表 | POST /approve/service/ajax-get-flow-list-with-permissions | POST | 无 |

---

## 三、审批详情 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-approve-detail | 获取审批详情 | POST /approve/service/ajax-get-approve-detail | POST | sid(string/是), type(int/否/1) |
| 2 | ajax-get-approve-path | 获取审批流程 | GET /approve/service/ajax-get-approve-path | GET | sid(string/是) |
| 3 | ajax-get-confirm-date | 获取确认日期 | GET /approve/service/ajax-get-confirm-date | GET | sid(string/是) |
| 4 | ajax-get-bdk-count | 获取员工补打卡次数 | GET /approve/service/ajax-get-bdk-count | GET | sid(long/是) |
| 5 | ajax-get-approve-over-timehour | 获取员工加班时长(管理员) | POST /approve/service/ajax-get-approve-over-timehour | POST | sid(string/是) |

---

## 四、审批催办 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-flow-urge | 审批催办 | GET /approve/service/ajax-flow-urge | GET | sid(string/是) |
| 2 | ajax-test-urge | 催办验证 | GET /approve/service/ajax-test-urge | GET | sid(string/是) |
| 3 | ajax-test-batch-urge | 批量催办验证 | GET /approve/service/ajax-test-batch-urge | GET | 无 |
| 4 | ajax-batch-flow-urge | 批量催办 | GET /approve/service/ajax-batch-flow-urge | GET | snowIds(string/是) |

---

## 五、审批设置 (FlowSettingController + FlowOldController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-init-flow-setting-new | 获取新建审批基础设置 | GET /approve/service/ajax-get-init-flow-setting-new | GET | 无 |
| 2 | ajax-get-flow-setting-new-list | 获取公司新审批类型id | GET /approve/service/ajax-get-flow-setting-new-list | GET | 无 |
| 3 | ajax-get-flow-setting | 获取组内的审批设置列表 | GET /approve/service/ajax-get-flow-setting | GET | groupId(int/是) |
| 4 | ajax-get-flow-setting-new | 获取审批基础设置 | GET /approve/service/ajax-get-flow-setting-new | GET | flowSettingId(int/是) |
| 5 | ajax-get-approve-flow-detail | 获取一个公司某项审批的详情设置 | GET /approve/service/ajax-get-approve-flow-detail | GET | type(int/否/null) |
| 6 | ajax-get-approve-flow-detail-new | 获取公司审批设置(新) | GET /approve/service/ajax-get-approve-flow-detail-new | GET | type(int/否/null), flowSettingId(int/否/null) |
| 7 | ajax-save-approve-flow-detail-new | 保存审批流设置接口(新) | POST /approve/service/ajax-save-approve-flow-detail-new | POST | data(string/是) |
| 8 | ajax-get-approve-setting-new-list | 根据管理员权限获取(新)审批设置 | GET /approve/service/ajax-get-approve-setting-new-list | GET | settingId(int/是) |
| 9 | ajax-get-setting-detail | 获取审批高级设置 | GET /approve/service/ajax-get-setting-detail | GET | flowSettingId(int/是) |
| 10 | ajax-save-setting-detail | 保存审批高级设置 | POST /approve/service/ajax-save-setting-detail | POST | data(string/是) |
| 11 | ajax-add-flow-setting | 新建审批 | POST /approve/service/ajax-add-flow-setting | POST | data(string/是) |
| 12 | ajax-check-flow-name | 判断审批名称是否重复 | GET /approve/service/ajax-check-flow-name | GET | flowSettingId(int/是), flowName(string/是), flowNameEng(string/否/"") |
| 13 | ajax-operate-switch | 操作审批设置开关 | POST /approve/service/ajax-operate-switch | POST | flowType(int/是), openStatus(bool/是), flowSettingId(int/是), isOld(int/是) |
| 14 | ajax-get-flow-setting-by-suiteType | 获取某一个套件类型审批列表 [已废弃] | GET /approve/service/ajax-get-flow-setting-by-suiteType | GET | suiteType(int/是) |
| 15 | ajax-get-approve-setting | 获取一个公司的审批设置列表 [已废弃] | GET /approve/service/ajax-get-approve-setting | GET | 无 |
| 16 | ajax-save-approve-flow-detail | 保存审批设置(旧) | POST /approve/service/ajax-save-approve-flow-detail | POST | data(string/是) |

---

## 六、审批表单 (FlowSettingController + FlowOldController + FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-flow-form-new | 获取审批字段表单new | GET /approve/service/ajax-get-flow-form-new | GET | flowSettingId(int/是), flowType(int/是) |
| 2 | ajax-save-custom-form-new | 保存审批表单设置new | POST /approve/service/ajax-save-custom-form-new | POST | data(string/是) |
| 3 | ajax-delete-custom-flow-new | 删除自定义审批 | POST /approve/service/ajax-delete-custom-flow-new | POST | flowType(int/是), flowSettingId(int/是), isOld(int/是) |
| 4 | ajax-get-custom-form | 获取审批自定义字段表单 [已废弃] | GET /approve/service/ajax-get-custom-form | GET | flowType(int/是) |
| 5 | ajax-get-flow-form | 获取审批字段表单 | GET /approve/service/ajax-get-flow-form | GET | flowType(int/是) |
| 6 | ajax-save-custom-form | 保存审批自定义字段表单 | POST /approve/service/ajax-save-custom-form | POST | data(string/是) |
| 7 | ajax-delete-custom-flow | 删除审批自定义字段表单 | POST /approve/service/ajax-delete-custom-flow | POST | flowType(int/是) |
| 8 | ajax-get-flow-preview-form | 获取审批预览的表单接口 | POST /approve/service/ajax-get-flow-preview-form | POST | flowType(int/是), flowSettingId(int/是), employeeId(string/否/null) |
| 9 | ajax-get-preview-relation-form | 根据不同选项获取审批预览字段 | POST /approve/service/ajax-get-preview-relation-form | POST | flowType(int/是), flowSettingId(int/是), holidayType(string/是), data(string/否/null) |
| 10 | ajax-launch-flow-preview | 发起审批预览 | POST /approve/service/ajax-launch-flow-preview | POST | data(string/是) |
| 11 | ajax-get-source-fields | 获取系统资源字段 | GET /approve/service/ajax-get-source-fields | GET | 无 |
| 12 | ajax-get-source-fields-new | 获取系统资源字段新 | POST /approve/service/ajax-get-source-fields-new | POST | 无 |
| 13 | ajax-get-business-field | 获取表单业务字段 | POST /approve/service/ajax-get-business-field | POST | flowSettingId(int/是), suiteType(int/是) |
| 14 | ajax-get-abstract-title-field | 获取摘要和标题可选字段 | POST /approve/service/ajax-get-abstract-title-field | POST | settingId(int/是) |
| 15 | ajax-check-biz-field | 校验表单业务字段 | POST /approve/service/ajax-check-biz-field | POST | fieldId/fieldName/fieldValue等(SingleFieldInfo绑定) |
| 16 | ajax-check-form-required | 校验表单必填字段 | POST /approve/service/ajax-check-form-required | POST | stepNodeId(string/是) |
| 17 | ajax-get-reference-field-setting | 获取管理字段勾选接口 | POST /approve/service/ajax-get-reference-field-setting | POST | 无 |
| 18 | ajax-set-reference-field-setting | 保存管理字段勾选接口 | POST /approve/service/ajax-set-reference-field-setting | POST | data(string/是) |
| 19 | ajax-fixed-add-field | 获取固定组可以添加的字段 | POST /approve/service/ajax-fixed-add-field | POST | settingId(int/是) |
| 20 | ajax-get-flow-relation-fields | 获取关联字段 | POST /approve/service/ajax-get-flow-relation-fields | POST | flowType(int/否/null), customField(string/否/null), fieldId(string/否/null), choose(string/否/null), sid(long/否/null), headcountDepartmentId(string/否/null) |

---

## 七、审批选项关联 (FlowSettingController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-save-option-associated | 保存选项关联关系 | POST /approve/service/ajax-save-option-associated | POST | flowTypeId(int/是), relations(string/是) |
| 2 | ajax-get-option-associated-data | 获取选项关联数据 | GET /approve/service/ajax-get-option-associated-data | GET | flowTypeId(int/是) |
| 3 | ajax-get-preview-option-associated | 预览根据不同选项展示的字段 | POST /approve/service/ajax-get-preview-option-associated | POST | flowTypeId(int/是), flowSettingId(int/是), fieldId(string/是), choose(string/是), employeeId(string/是), holidayType(string/否/null), data(string/否/null) |

---

## 八、审批分组 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-flow-group | 获取审批分组 | GET /approve/service/ajax-get-flow-group | GET | 无 |
| 2 | ajax-add-flow-group | 添加审批分组 | POST /approve/service/ajax-add-flow-group | POST | flowGroupName(string/是) |
| 3 | ajax-update-flow-group | 更新审批分组 | POST /approve/service/ajax-update-flow-group | POST | flowGroupName(string/是), groupId(int/是) |
| 4 | ajax-save-flow-group | 保存审批分组 | POST /approve/service/ajax-save-flow-group | POST | flowGroupName/groupId等(FlowGroupSettingModel绑定) |
| 5 | ajax-remove-flow-group | 删除审批分组 | POST /approve/service/ajax-remove-flow-group | POST | groupId(int/是) |

---

## 九、审批节点规则 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-all-setting-node-rule-by-setting-id | 获取所有节点规则 | GET /approve/service/ajax-get-all-setting-node-rule-by-setting-id | GET | settingId(int/是) |
| 2 | ajax-overwrite-all-setting-node-rule-by-setting-id | 覆盖所有节点规则 | POST /approve/service/ajax-overwrite-all-setting-node-rule-by-setting-id | POST | settingId(int/是), rules(@RequestBody List/是) |
| 3 | ajax-get-setting-node-rule | 获取节点规则 | GET /approve/service/ajax-get-setting-node-rule | GET | settingId/flowType等(SettingNodeRuleParam绑定) |
| 4 | ajax-save-setting-node-rule | 保存节点规则 | POST /approve/service/ajax-save-setting-node-rule | POST | settingId/flowType等(SettingNodeRuleParam绑定) |

---

## 十、审批角色 (FlowRuleController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-save-flow-role | 新增/编辑角色 | POST /approve/service/ajax-save-flow-role | POST | data(string/是) |
| 2 | ajax-get-flow-role-detail | 获取角色详情 | GET /approve/service/ajax-get-flow-role-detail | GET | roleId(string/是) |
| 3 | ajax-get-flow-role-list | 获取角色列表 | GET /approve/service/ajax-get-flow-role-list | GET | keyWord(string/否/null) |
| 4 | ajax-get-flow-role-used | 获取角色使用情况 | GET /approve/service/ajax-get-flow-role-used | GET | roleId(string/是) |
| 5 | ajax-get-all-role | 获取所有角色 | GET /approve/service/ajax-get-all-role | GET | 无 |
| 6 | ajax-delete-flow-role | 删除角色 | POST /approve/service/ajax-delete-flow-role | POST | roleId(string/是) |

---

## 十一、审批表单规则 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-relation-field-rule | 获取审批请假表单设置规则 | POST /approve/service/ajax-get-relation-field-rule | POST | flowType(int/是) |
| 2 | ajax-save-relation-field-rule | 保存审批表单设置规则 | POST /approve/service/ajax-save-relation-field-rule | POST | flowType(int/是), data(string/是) |
| 3 | ajax-validate-expression | 验证公式 | POST /approve/service/ajax-validate-expression | POST | settingId(int/是), expression(string/是), flowFieldModels(string/是) |
| 4 | ajax-verification-setting | 验证审批流 | GET /approve/service/ajax-verification-setting | GET | flowType(int/是), settingId(int/是) |
| 5 | ajax-verification-setting-branch-repeat | 审批流校验重复分支 | POST /approve/service/ajax-verification-setting-branch-repeat | POST | data(string/是) |
| 6 | ajax-save-checked-condition-fields | 保存勾选的分支条件字段 | POST /approve/service/ajax-save-checked-condition-fields | POST | data(string/是), settingId(int/是), conditionFields(string/是) |
| 7 | ajax-optional-employee | 分支条件可选审批人 | GET /approve/service/ajax-optional-employee | GET | flowType(int/是), settingId(int/是) |
| 8 | ajax-check-write-rule | 检查写入规则 | POST /approve/service/ajax-check-write-rule | POST | data(string/是) |
| 9 | ajax-alter-write-rule | 修改写入规则 | POST /approve/service/ajax-alter-write-rule | POST | data(string/是) |

---

## 十二、审批导出 (FlowExportController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-export-flow-check | 获取审批列表前置检查 | POST /approve/service/ajax-export-flow-check | POST | search[value](string/否/""), monthFilter(int/否/-1), speedy[value](string/否/""), speedySearch(bool/否/false) |
| 2 | ajax-export-flow | 导出审批excel | POST /approve/service/ajax-export-flow | POST | name(string/是), search[value](string/否/""), monthFilter(int/否/-1), speedy[value](string/否/""), speedySearch(bool/否/false), verifyKey(string/否/""), isExportProcess(int/否/0), isExportFile(int/否/0), isExportMask(int/否/1) |

---

## 十三、审批委托 (FlowEntrustController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-flow-entrust-list | 获取审批委托列表 | GET /approve/service/ajax-get-flow-entrust-list | GET | pageNum(int/否/1), pageSize(int/否/50) |
| 2 | ajax-append-entrust | 审批委托 | POST /approve/service/ajax-append-entrust | POST | data(string/是) |
| 3 | ajax-cancel-entrust | 取消委托 | POST /approve/service/ajax-cancel-entrust | POST | id(int/是) |
| 4 | ajax-flow-entrust-setting | 设置员工发起委托 | GET /approve/service/ajax-flow-entrust-setting | GET | status(int/是) |

---

## 十四、证明类审批 (ProofFlowController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-add-proof-record | 添加证明开具记录 | POST /approve/service/ajax-add-proof-record | POST | sid(string/是), recordStatus(int/是) |
| 2 | ajax-search-proof-record | 获取证明开具记录 | POST /approve/service/ajax-search-proof-record | POST | sid(string/是) |
| 3 | ajax-get-proof-snap | 获取证明快照 | POST /approve/service/ajax-get-proof-snap | POST | sid(string/是) |

---

## 十五、离职交接 (DismissHandoverController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-handover-scope-filter | 获取离职交接方案适用范围筛选 | GET /approve/service/ajax-get-handover-scope-filter | GET | 无 |
| 2 | ajax-get-dimission-setting-list | 获取公司离职交接方案集合 | GET /approve/service/ajax-get-dimission-setting-list | GET | 无 |
| 3 | ajax-get-dimission-setting | 获取离职交接的设置 | GET /approve/service/ajax-get-dimission-setting | GET | handoverBasicId(int/是) |
| 4 | ajax-update-dimission-setting | 保存离职交接的设置 | POST /approve/service/ajax-update-dimission-setting | POST | dimissionSetting(string/是) |
| 5 | ajax-delete-dimission-setting | 删除离职交接方案 | POST /approve/service/ajax-delete-dimission-setting | POST | handoverBasicId(int/是) |
| 6 | ajax-copy-dimission-setting | 复制离职交接的设置 | POST /approve/service/ajax-copy-dimission-setting | POST | dimissionSetting(string/是) |
| 7 | ajax-get-demission-basic-message-option | 获取基本信息可选字段 | GET /approve/service/ajax-get-demission-basic-message-option | GET | handoverBasicId(int/否/null) |
| 8 | ajax-get-handover-approve-detail | 获取审批的离职交接详情 | GET /approve/service/ajax-get-handover-approve-detail | GET | flowProcessId(int/是), isPrint(int/否/0) |
| 9 | ajax-get-transfer-approval-list | 离职交接审批转交列表 [已废弃] | GET /approve/service/ajax-get-transfer-approval-list | GET | flowProcessId(int/是) |
| 10 | ajax-transfer-approval | 离职交接审批转交操作 | POST /approve/service/ajax-transfer-approval | POST | approvalId(int/是), directorEmployeeId(string/是), transferReason(string/是) |
| 11 | ajax-cancel-handover-approval | 离职交接管理员撤销 | POST /approve/service/ajax-cancel-handover-approval | POST | approvalId(int/是), cancelReason(string/是) |
| 12 | ajax-urge | 离职交接催办 | POST /approve/service/ajax-urge | POST | approvalId(int/是) |

---

## 十六、组织架构查询 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-auth-department | 根据权限获取部门树 | GET /approve/service/ajax-get-auth-department | GET | isNeedRoot(int/否/1), needVirtualDepart(int/否/0), departmentId(string/否/"") |
| 2 | ajax-get-complete-department-tree-for-search | 获取部门完整树 | GET /approve/service/ajax-get-complete-department-tree-for-search | GET | 无 |
| 3 | ajax-search-department-or-employee | 模糊查询员工或部门 | GET /approve/service/ajax-search-department-or-employee | GET | keyword(string/是), type(int/是), excludeEmployeeIds(string/否/""), excludeDepartmentIds(string/否/""), filterByAuth(bool/否/true), includeOpenScope(bool/否/false), departmentId(string/否/""), isNeedRoot(int/否/1), isNeedVirtualDepart(int/否/0), needVirtualDepart(int/否/0) |
| 4 | ajax-get-department-employees | 根据部门获取员工列表 | GET /approve/service/ajax-get-department-employees | GET | departmentId(string/是), isContainSub(bool/否/true), filterByAuth(bool/否/true), excludeEmployeeIds(string/否/""), pageNum(int/否/1), pageSize(int/否/50) |
| 5 | ajax-get-department-list | 根据员工获取部门列表 | GET /approve/service/ajax-get-department-list | GET | employeeId(string/是) |
| 6 | ajax-get-department-tips | 根据员工ids获取部门链路 | POST /approve/service/ajax-get-department-tips | POST | employeeIds(string/是) |

---

## 十七、审批编辑与打印 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-page-edit-history | 获取编辑历史 | GET /approve/service/ajax-get-page-edit-history | GET | sid(string/是) |
| 2 | ajax-save-page-edit | 保存页面编辑 | POST /approve/service/ajax-save-page-edit | POST | sid(string/是), data(string/是) |
| 3 | ajax-edit-print | 获取打印编辑页面 | POST /approve/service/ajax-edit-print | POST | sid(string/是) |
| 4 | ajax-get-print-hide-fields | 获取打印隐藏字段 | GET /approve/service/ajax-get-print-hide-fields | GET | settingId(int/是) |
| 5 | ajax-saves-flow-print | 保存打印设置 | POST /approve/service/ajax-saves-flow-print | POST | data(string/是) |
| 6 | ajax-flow-print | 审批打印 | POST /approve/service/ajax-flow-print | POST | sid(string/是) |
| 7 | ajax-print-flow | 审批详情打印记录 | POST /approve/service/ajax-print-flow | POST | sid(string/是) |

---

## 十八、审批范围与开关 (FlowAllController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-approve-open-status | 获取当前公司是否有开启审批 | GET /approve/service/ajax-get-approve-open-status | GET | 无 |
| 2 | ajax-get-flow-switch | 获取公司审批范围 | GET /approve/service/ajax-get-flow-switch | GET | flowType(int/是) |
| 3 | ajax-save-flow-switch | 保存公司审批范围 | POST /approve/service/ajax-save-flow-switch | POST | flowType(int/是), adminIsOpen(int/是), leaderIsOpen(int/是), employeeIsOpen(int/是) |
| 4 | ajax-get-charge-setting | 获取加班扣费设置 | GET /approve/service/ajax-get-charge-setting | GET | 无 |
| 5 | ajax-get-ranks | 获取职级 [已废弃] | GET /approve/service/ajax-get-ranks | GET | 无 |
| 6 | ajax-open-or-close-irreversible | 开启或关闭撤销权限 | POST /approve/service/ajax-open-or-close-irreversible | POST | isOpenOrCloseCancelPermissions(int/是) |
| 7 | ajax-open-or-close-limit-start-time | 开启或关闭限制发起时间 | POST /approve/service/ajax-open-or-close-limit-start-time | POST | (OpenOrCloseLimitGrantTimeModel绑定) |
| 8 | ajax-open-or-close-setting-details | 开启或关闭高级设置 | POST /approve/service/ajax-open-or-close-setting-details | POST | data(string/是) |
| 9 | ajax-batch-update-setting-detail | 批量更新高级设置 | POST /approve/service/ajax-batch-update-setting-detail | POST | data(string/是) |
| 10 | ajax-get-company-advanced-setting | 获取公司高级设置 | GET /approve/service/ajax-get-company-advanced-setting | GET | 无 |
| 11 | ajax-open-mask-fields | 掩码字段展示明文信息 | POST /approve/service/ajax-open-mask-fields | POST | sid(long/是), labelName(string/是), isOldValue(int/是), groupIndex(int/是), isFixed(int/是) |

---

## 十九、其他功能 (FlowAllController + FlowSettingController)

| # | Schema Name | 接口含义 | API 路径 | 方法 | 参数（类型/必填/默认值） |
|---|-------------|----------|----------|------|--------------------------|
| 1 | ajax-get-novice-guide | 获取主系统审批-新手引导 | POST /approve/service/ajax-get-novice-guide | POST | 无 |
| 2 | ajax-save-novice-guide | 保存主系统审批-新手引导 | POST /approve/service/ajax-save-novice-guide | POST | type(string/是) |
| 3 | ajax-get-last-sign | 管理员获取最后一次签名 | GET /approve/service/ajax-get-last-sign | GET | 无 |
| 4 | ajax-cover-approve-detail-new | 覆盖新审批流 | POST /approve/service/ajax-cover-approve-detail-new | POST | coverSettingId(int/是), settingId(int/是) |
| 5 | ajax-get-all-custom-suite-module | 获取所有套件设置数据 | GET /approve/service/ajax-get-all-custom-suite-module | GET | 无 |
| 6 | ajax-get-calc-platform-plan-group-list | 获取计算平台方案分组下拉数据 | GET /approve/service/ajax-get-calc-platform-plan-group-list | GET | 无 |
| 7 | ajax-get-calc-platform-group-field-list | 获取计算平台方案分组下字段列表 | GET /approve/service/ajax-get-calc-platform-group-field-list | GET | planId(string/否/null), groupId(int/否/null) |
| 8 | ajax-get-department-headcount | 获取编制信息 | GET /approve/service/ajax-get-department-headcount | GET | headcountDepartmentId(string/是), applyForExpansion(int/否/0), startValue(string/是), endValue(string/是) |
| 9 | ajax-get-department-headcount-detail | 获取部门指定月份的编制信息 | GET /approve/service/ajax-get-department-headcount-detail | GET | headcountDepartmentId(string/是), sid(long/是) |
| 10 | ajax-get-headcount-dimension-quantity | 获取编制纬度数量 | POST /approve/service/ajax-get-headcount-dimension-quantity.json | POST | (@RequestBody DepartmentHeadcountDetailsParam/是) |
| 11 | ajax-get-field-dropdown-list | 获取字段动态下拉选项数据 | POST /approve/service/ajax-get-field-dropdown-list.json | POST | (@RequestBody DyncDropdownListParam/是) |
| 12 | ajax-get-emp-custom-addr-fields | 获取员工自定义地址字段模型列表 | GET /approve/service/ajax-get-emp-custom-addr-fields | GET | 无 |
| 13 | ajax-email-address | 邮件审批中转地址(无需登录) | GET /approve/service/ajax-email-address | GET | token(string/是), pid(string/是) |
