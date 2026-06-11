---
id: "1"
title: "task-executor.md agent 协议对齐"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: task-executor.md agent 协议对齐

## Description

修复 `plugins/forge/agents/task-executor.md` 中 agent Execution Protocol 与 prompt 模板、guide、submit-task skill 之间的 5 处语义断裂（Issues 1, 2, 3, 8, 9）。

核心问题是多 feature 独立演进时 agent 协议未同步更新，导致 agent 可能执行与 skill 语义矛盾的操作。

## Reference Files
- `docs/proposals/task-executor-prompt-congruence/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/agents/task-executor.md` | step 8→9 加 blocked 复查; 统一 retry 策略; 修正 submit-task 注释; 补 blocked DONE 格式; stop→Complex Error Pause Flow 语义 |

## Acceptance Criteria

- [ ] step 8→9 之间新增 step 8.5：复查 submit-task 结果，若 auto-downgrade 为 blocked 则跳过 git-commit，直接输出 blocked 格式 DONE
- [ ] retry 策略统一为 ~3 attempts（与 Complex Error Pause Flow 一致），模板中 "max 1 retry then stop" 语义不再与 agent 定义的阈值矛盾
- [ ] 所有 "stop" 指令明确包含 "eval Complex Error Pause Flow" 语义：stop 不是立即终止，而是先评估是否需要创建 fix task 并 block source
- [ ] step 8 注释中 "via just test" 的错误描述被修正为准确描述 submit-task 的实际行为
- [ ] blocked 任务有明确的 DONE 输出格式：`DONE: <task-id> | blocked | <summary>`（无 commit-hash）
- [ ] submit-task auto-downgrade 为 blocked 时，agent 不执行 git-commit

## Hard Rules

- 不改变 agent 的整体 11 步结构，仅在 step 8→9 之间插入 step 8.5
- 不修改 prompt.go 或 prompt 模板文件

## Implementation Notes

- Issue 1 (P0) 最紧急：blocked task 违规 commit 是 CI 中的实际缺陷
- "stop" 语义对齐需要同时考虑 agent 定义（task-executor.md）和模板侧（Task 2），但本任务只改 agent 侧
- retry 策略统一方向为 agent 定义的 ~3 attempts（比模板的 1 retry 更宽容），在 agent 协议中明确表述阈值
- blocked DONE 格式：参考 run-tasks dispatcher 的 DONE 解析逻辑，确保 blocked 格式不含 commit-hash
