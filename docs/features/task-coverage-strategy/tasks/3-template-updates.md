---
id: "3"
title: "Update coding-type prompt templates with coverage directives"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Update coding-type prompt templates with coverage directives

## Description

更新所有 coding 类型的 prompt 模板，使用 Task 2 注入的 `{{COVERAGE_TARGET}}` 和 `{{COVERAGE_STRATEGY}}` 占位符。各模板按分档策略呈现覆盖率要求。同时为 `coding-refactor` 模板增加增量编译检查策略。

## Reference Files
- `docs/proposals/task-coverage-strategy/proposal.md` — Source proposal
- `forge-cli/pkg/prompt/data/coding-feature.md` — Feature 模板
- `forge-cli/pkg/prompt/data/coding-enhancement.md` — Enhancement 模板
- `forge-cli/pkg/prompt/data/coding-fix.md` — Fix 模板
- `forge-cli/pkg/prompt/data/coding-refactor.md` — Refactor 模板
- `forge-cli/pkg/prompt/data/coding-cleanup.md` — Cleanup 模板

## Acceptance Criteria

- `coding-feature.md` 模板：测试步骤包含 `{{COVERAGE_TARGET}}` 指令，agent 知道达到目标后停止补充测试
- `coding-enhancement.md` 模板：同 feature，使用 `{{COVERAGE_TARGET}}`
- `coding-fix.md` 模板：测试步骤包含 `{{COVERAGE_TARGET}}`，agent 写针对性修复测试而非追求高覆盖
- `coding-refactor.md` 模板：使用 maintain 策略指令 + **增量编译策略**（每修改一个文件后立即 `just compile`，确认无级联错误再继续）
- `coding-cleanup.md` 模板：使用 maintain 策略指令，不新增测试
- 所有模板更新后，`forge prompt get-by-task-id` 能正确渲染覆盖率指令
- 现有测试通过

## Hard Rules

- 模板中覆盖率指令放在测试步骤的显著位置，确保 agent 优先看到
- refactor 模板的增量编译策略格式为：修改文件 → `just compile` → 通过则继续 → 失败则立即修复当前文件，不要继续修改其他文件
- 保持模板整体结构不变，只在测试相关段落插入覆盖率指令

## Implementation Notes

- `coding-fix.md` 当前没有 coverage 输出格式，需要新增测试步骤中的覆盖率指令段落
- `coding-refactor.md` 已有 "Coverage is informational... drops >2%, investigate" 的描述，需要更新为 `{{COVERAGE_TARGET}}` 的 maintain 模式输出
- 各模板的测试步骤位置不完全相同，需要逐一确认插入点
