---
id: "4"
title: "P1: test.* 模板 HARD-RULE 标签重命名为 TASK-CONSTRAINTS"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 4: P1: test.* 模板 HARD-RULE 标签重命名为 TASK-CONSTRAINTS

## Description
7 个 test.* 模板使用 `<HARD-RULE>` 标签定义任务约束（如必须调用 skill），但这与 task 文件中的 "Hard Rules" 概念混淆。将标签改名为 `<TASK-CONSTRAINTS>` 消除歧义。同时从 test-eval-cases.md 中移除不应由模板控制的 MAIN_SESSION 声明。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (Sections 1.4, 2.11-2.17, P2 #13)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/test-gen-cases.md` | `<HARD-RULE>` → `<TASK-CONSTRAINTS>` |
| `forge-cli/pkg/prompt/data/test-eval-cases.md` | `<HARD-RULE>` → `<TASK-CONSTRAINTS>`；移除 MAIN_SESSION 声明 |
| `forge-cli/pkg/prompt/data/test-gen-scripts.md` | `<HARD-RULE>` → `<TASK-CONSTRAINTS>` |
| `forge-cli/pkg/prompt/data/test-run.md` | `<HARD-RULE>` → `<TASK-CONSTRAINTS>` |
| `forge-cli/pkg/prompt/data/test-gen-and-run.md` | `<HARD-RULE>` → `<TASK-CONSTRAINTS>` |
| `forge-cli/pkg/prompt/data/test-graduate.md` | `<HARD-RULE>` → `<TASK-CONSTRAINTS>` |
| `forge-cli/pkg/prompt/data/test-verify-regression.md` | `<HARD-RULE>` → `<TASK-CONSTRAINTS>` |

## Acceptance Criteria
- [ ] 所有 7 个 test.* 模板中 `<HARD-RULE>` 和 `</HARD-RULE>` 替换为 `<TASK-CONSTRAINTS>` 和 `</TASK-CONSTRAINTS>`
- [ ] test-eval-cases.md 中 MAIN_SESSION 相关的 `<EXTREMELY-IMPORTANT>` 声明已移除
- [ ] 标签内的约束内容不变，仅标签名变更

## Implementation Notes
- 仅替换标签名，不修改标签内容
- MAIN_SESSION 由 dispatcher（execute-task/run-tasks）在分派时决定，不应硬编码在模板中
