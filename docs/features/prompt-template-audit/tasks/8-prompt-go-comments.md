---
id: "8"
title: "P3: prompt.go 添加占位符 escaping 警告注释"
priority: "P2"
estimated_time: "15m"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 8: P3: prompt.go 添加占位符 escaping 警告注释

## Description
prompt.go 的 `renderTemplate` 使用 `strings.ReplaceAll` 进行模板替换，如果模板内容中恰好包含 `{{TASK_ID}}` 等字符串（如代码示例中），会被意外替换。目前没有 escaping 机制。添加警告注释和 genScriptBases 映射说明注释。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (Section 4, P3 #16)

## Acceptance Criteria
- [ ] `renderTemplate` 函数附近添加警告注释：说明 `{{...}}` 占位符无 escaping 机制，模板内容中不能出现裸占位符字符串
- [ ] `extractTestTypeArg` 的 `genScriptBases` 列表附近添加注释：说明此列表与 task ID 格式的对应关系

## Implementation Notes
- 仅添加注释，不修改任何逻辑
- 这是防御性文档——未来如果模板中需要使用 `{{...}}` 字面量，需要先实现 escaping 机制
