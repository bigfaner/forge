---
id: "2"
title: "Prompt 模板错误处理与行为对齐"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Prompt 模板错误处理与行为对齐

## Description

修复 `forge-cli/pkg/prompt/data/` 下 8 个 prompt 模板的错误处理语义，使其与 task-executor agent 的 Execution Protocol 对齐（Issues 3-template, 4, 5-template, 6）。

模板当前定义了独立的错误处理策略（1 retry/stop），与 agent 定义的 Complex Error Pause Flow（~3 attempts/fix task）不一致。

## Reference Files
- `docs/proposals/task-executor-prompt-congruence/proposal.md` — Source proposal
- `plugins/forge/agents/task-executor.md` — Agent protocol（Task 1 修改后的版本）

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-cleanup.md` | fmt→WARNING; coverage 策略与 "no new tests" 对齐 |
| `forge-cli/pkg/prompt/data/coding-refactor.md` | fmt→WARNING; coverage 策略与 "no new tests" 对齐; Pre-check 失败→blocked |
| `forge-cli/pkg/prompt/data/coding-enhancement.md` | fmt→WARNING; stop→Complex Error Pause Flow |
| `forge-cli/pkg/prompt/data/coding-feature.md` | fmt→WARNING; stop→Complex Error Pause Flow |
| `forge-cli/pkg/prompt/data/coding-fix.md` | fmt→WARNING; stop→Complex Error Pause Flow |
| `forge-cli/pkg/prompt/data/validation-code.md` | fmt→WARNING; stop→Complex Error Pause Flow |
| `forge-cli/pkg/prompt/data/validation-ux.md` | stop→Complex Error Pause Flow |
| `forge-cli/pkg/prompt/data/gate.md` | stop→Complex Error Pause Flow |

## Acceptance Criteria

- [ ] 8 个模板中 `just fmt` 失败的行为从 Stop 改为 WARNING（non-blocking），与 guide 一致
- [ ] 模板中所有 "stop" 指令添加 "eval Complex Error Pause Flow" 语义说明
- [ ] coding-refactor Pre-check 失败时明确设置 blocked 状态（而非停留在 in_progress）
- [ ] coding-cleanup 和 coding-refactor 模板中 coverage 指令与 "no new tests" 不矛盾（要么不注入 coverage 指令，要么明确说明 "maintain existing coverage, no new tests required"）
- [ ] All-Completed Hook 的 quality gate 仍会运行 fmt 作为最终安全网

## Hard Rules

- 不修改 prompt.go（Task 3 负责 resolveCoverage 代码修复）
- 不修改 guide 本身（guide 定义了正确行为，模板需要对齐）
- fmt 改为 WARNING 后仍需保留 fmt 检查步骤，只是不阻断任务

## Implementation Notes

- Issue 6 (fmt WARNING) 涉及最多模板文件，但改动模式一致（Stop→WARNING）
- Issue 4 (Pre-check blocked) 仅影响 coding-refactor.md，需要添加明确的状态转换指令
- Issue 5 的模板侧：coding-cleanup 和 coding-refactor 模板中有 "no new tests" 指令，但 resolveCoverage() 注入的 coverage 策略（如 "achieve 80% coverage"）与之矛盾。模板侧修复是确保模板文本本身不含矛盾，代码侧修复在 Task 3
- 风险：fmt WARNING 可能导致格式问题未被发现，但 All-Completed Hook 的 quality gate 仍会运行 fmt 作为安全网
