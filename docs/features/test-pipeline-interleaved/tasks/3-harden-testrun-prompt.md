---
id: "3"
title: "Add execution instructions to test-run prompt template"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
mainSession: false
---

# 3: Add execution instructions to test-run prompt template

## Description

test-run 的 prompt 模板（`pkg/prompt/templates/test-run.md`）缺少对 agent 修改行为的约束指令。需要添加：(1) 确认是正式代码问题才能修改正式代码，测试脚本 bug 可以修但不能篡改测试逻辑；(2) 问题多时通过 `forge task add` 追加 fix 任务，与 task-executor Pause Protocol 协调。

## Reference Files
- `docs/proposals/test-pipeline-interleaved/proposal.md` — Proposed Solution, Scope > In Scope, Success Criteria, Key Risks
- `forge-cli/pkg/prompt/templates/test-run.md`: add execution constraint instructions (ref: Scope > In Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/templates/test-run.md` | 在 TASK-CONSTRAINTS 部分追加行为约束指令 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 模板 TASK-CONSTRAINTS 包含指令：确认是正式代码 bug 才能修改正式代码；测试脚本本身的 bug 可以修，但不能为了通过测试而篡改测试断言/逻辑
- [ ] 模板 TASK-CONSTRAINTS 包含指令：问题多时通过 `forge task add` 追加 fix 任务（而非在一个任务中修复所有问题），与 task-executor Pause Protocol 协调

## Hard Rules

- 新指令作为 TASK-CONSTRAINTS 级别的补充，不覆盖 task-executor 的 EXTREMELY-IMPORTANT 层级指令

## Implementation Notes

当前模板的 TASK-CONSTRAINTS 已有 3 条规则。追加新规则时保持现有格式（`- MUST ...` / `- MUST NOT ...`）。新指令与 Pause Protocol 的协调点：Pause Protocol 在 agent 遇到无法在当前任务范围内解决的问题时触发暂停；`forge task add` 则是在问题可分解但量多时主动创建子任务。两者不冲突。
