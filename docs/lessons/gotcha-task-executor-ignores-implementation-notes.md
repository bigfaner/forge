---
created: "2026-05-13"
tags: [testing]
---

# task-executor 忽略 Implementation Notes 中的明确 bash 指令

## Problem

T-test-3（run-e2e-tests）的 task-executor subagent 使用 `npx playwright test features/milestone-map/` 运行测试，而非任务文件明确指定的 `just test-e2e --feature milestone-map`。

这导致缺少 server 启动、环境变量（`E2E_FEATURE=1`）设置、testIgnore 绕过等 `just test-e2e` 提供的保障。

## Root Cause

**症状**：subagent 用 `npx playwright test` 直接跑测试

**直接原因**：subagent 读了任务文件和 justfile，确认 `just test-e2e --feature milestone-map` 指令存在且 recipe 可用，但仍然选择自己构造 `npx playwright test` 命令

**根本原因**：task-executor 对 Implementation Notes 的遵从优先级低于其自身行为模式。它的默认模式是"读任务 → 自己决定怎么执行 → just test → 记录"，当任务文件建议的命令与其直觉冲突时，它按自己直觉走。

**证据**（from subagent transcript agent-aa5c6f0904a4c4072）：
1. Line 4-5: 读任务文件，看到 Implementation Notes `just test-e2e --feature milestone-map`
2. Line 18: 读 justfile，grep `test-e2e` 找到 recipe
3. Line 39-40: 手动 curl 检查 backend/frontend 是否在运行
4. Line 56: 选择运行 `cd tests/e2e && npx playwright test features/milestone-map/ --reporter=list`

## Solution

将关键命令从 `## Implementation Notes` 提升为 `## Hard Rules`：

```markdown
## Hard Rules
- MUST use `just e2e-setup` then `just test-e2e --feature <slug>` to run tests
- MUST NOT use `npx playwright test` directly
```

## Reusable Pattern

**当任务需要 task-executor 执行特定的 bash 命令（而非自己决定如何执行）时：**

- 将命令写在 `## Hard Rules` 节，不要写在 `## Implementation Notes`
- task-executor 对 HARD-RULE 的遵从度 >> Implementation Notes
- 这个优先级差异是 task-executor agent prompt 的设计决策，短期内不会改变

**适用场景**：
- 必须通过 justfile recipe 执行（而非直接调用底层工具）
- 命令带有环境变量、server lifecycle 等隐含依赖
- agent 自行构造命令会导致测试环境不完整

**不适用场景**：
- 普通代码任务的 `just test`（agent 默认行为已足够）
- agent 可以自行决定执行策略的灵活任务

## Example

```markdown
# ❌ 不可靠 — agent 可能忽略
## Implementation Notes
1. Run `just test-e2e --feature milestone-map`

# ✅ 可靠 — agent 遵从 Hard Rules
## Hard Rules
- MUST use `just test-e2e --feature milestone-map` to run tests
- MUST NOT use `npx playwright test` directly
```

## Related Files

- `docs/features/milestone-map/tasks/run-e2e-tests.md` (已修复为 Hard Rules)
- `justfile` (`test-e2e` recipe)
- `tests/e2e/playwright.config.ts`
- `docs/lessons/gotcha-split-task-missing-shared-setup.md` (级联上游问题)
