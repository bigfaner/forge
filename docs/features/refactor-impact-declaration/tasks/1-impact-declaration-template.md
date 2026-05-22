---
id: "1"
title: "Add Impact Declaration to coding-refactor.md template"
priority: "P1"
estimated_time: "45m"
dependencies: []
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 1: Add Impact Declaration to coding-refactor.md template

## Description

在重构执行模板 `forge-cli/pkg/prompt/data/coding-refactor.md` 中新增 Impact Declaration 机制。执行器在动手修改代码前，先分析受影响的测试并分类为 PRESERVE（行为不变）或 EVOLVE（行为预期变化），输出结构化声明。重构过程中 EVOLVE 测试的失败被视为预期变化，直接更新测试断言；未声明或 PRESERVE 分类的测试失败仍触发暂停。

## Reference Files
- `docs/proposals/refactor-impact-declaration/proposal.md` — Source proposal
- `docs/lessons/gotcha-characterization-test-vs-refactoring.md` — Root cause lesson

## Acceptance Criteria

- [ ] Step 2 (Impact Mapping) 末尾新增 Impact Declaration 子步骤，要求执行器在动手前输出 `IMPACT_DECLARATION` 结构化声明
- [ ] 声明格式包含每个受影响测试的：test name、classification (PRESERVE/EVOLVE)、reason（必填）、expected_change（EVOLVE 时必填）
- [ ] EVOLVE 条目缺少 reason 或 expected_change 时，模板要求执行器重新分类为 PRESERVE
- [ ] Step 3 (Refactor) 的 Universal constraints 更新：EVOLVE 测试失败时不输出 BEHAVIOR_CHANGE_DETECTED，而是更新测试断言
- [ ] Step 4 (Static Checks) 的 targeted test 失败处理更新：区分 EVOLVE（更新）、PRESERVE/未声明（暂停）
- [ ] 声明输出在执行记录中可见（作为 Step 2 输出的一部分）
- [ ] `go build ./...` 通过（模板变更不影响编译）

## Hard Rules

- 保留现有 BEHAVIOR_CHANGE_DETECTED 作为 PRESERVE/未声明测试的失败处理机制
- 不修改 Step 2 的 Impact Mapping 逻辑本身，只在末尾追加 Declaration 子步骤
- 声明格式必须包含具体示例，降低执行器的遵循门槛

## Implementation Notes

当前模板的 Step 2 已有 Impact Mapping（文件级/符号级分析），新增的 Declaration 子步骤是在此基础上做测试级分类。两者是互补关系，不是替代。

Step 3 当前的 Universal constraints 第 90 行："If test assertions need changes, the refactor is changing behavior — output `BEHAVIOR_CHANGE_DETECTED: <description>` and skip that specific change" — 这条需要修改为：先检查 Declaration，如果是 EVOLVE 分类则更新测试，否则仍输出 BEHAVIOR_CHANGE_DETECTED。

Step 4 第 167 行的 targeted test 失败处理："assertion changes → BEHAVIOR_CHANGE_DETECTED + skip" — 同样需要区分 EVOLVE vs PRESERVE。

风险控制：执行器可能过度声明 EVOLVE。通过要求 EVOLVE 条目必须有 reason 和 expected_change 来约束——空理由的 EVOLVE 视为无效。
