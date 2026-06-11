---
id: "1"
title: "Rename validate-index to validate and add AC count validation"
priority: "P0"
estimated_time: "1.5h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 1: Rename validate-index to validate and add AC count validation

## Description
Rename the `forge task validate-index` command to `forge task validate`, and add programmatic AC (Acceptance Criteria) count validation to enforce task sizing rules. The validator must parse each task .md file's `## Acceptance Criteria` section, count `- [ ]` prefixed lines, and reject tasks with AC > 6 or AC = 0.

This addresses the root cause identified in the proposal: LLM agents ignore written sizing rules, so programmatic enforcement is needed.

## Reference Files
- `docs/proposals/task-sizing-gate/proposal.md` — Proposed Solution, Constraints & Dependencies, Scope > In Scope
- `forge-cli/internal/cmd/task/validate_index.go`: 重命名命令 Use 字段 + 新增 validateACCount 校验函数，解析 ## Acceptance Criteria 下的 `- [ ]` 行计数 (ref: Proposed Solution)
- `forge-cli/internal/cmd/task/register.go`: 命令注册名需更新 (ref: Proposed Solution)
- `forge-cli/internal/cmd/task/testbridge.go`: ExportRunValidateIndex 导出名需同步更新 (ref: Proposed Solution)
- `forge-cli/internal/cmd/root_test.go`: expectedTaskSubs 中 "validate-index" 需更新为 "validate" (ref: Proposed Solution)

## Acceptance Criteria
- [ ] `forge task validate docs/features/<slug>/tasks/index.json` 校验 index.json 结构 + 所有 task 文件的 AC 数量（向后兼容原有校验）
- [ ] AC > 6 时返回 exit 1 + 错误信息（包含 task 文件名和 AC 数量）
- [ ] AC = 0 时返回 exit 1 + 错误信息
- [ ] `forge task validate-index` 不再存在（直接替换，breaking change）

## Implementation Notes
- AC 解析逻辑：读取 task .md 文件，定位 `## Acceptance Criteria` section，统计 `- [ ]` 前缀的行数。格式由模板保证一致性。
- 文件重命名：建议将 `validate_index.go` 重命名为 `validate.go`，内部函数名 `runValidateIndex` → `runValidate`，`validateIndexCmd` → `validateCmd`。testbridge 导出名同步更新。
- 需更新所有测试文件中的 "validate-index" 字符串引用（约 6 个文件），使用 `grep -r "validate-index" forge-cli/ tests/` 确认完整列表。
- Key Risk: AC 解析误判（非标准格式）— 风险低，因格式由模板保证。

### Test Impact
- Affected test suite(s): `forge-cli/internal/cmd/`, `forge-cli/tests/task-type-system/`, `tests/task-type-system/`
- Expected fixture changes: 无（字符串替换即可）
- Risk level: low
