---
id: "3"
title: "Update quality gate steps and failure handling"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["1", "2"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Update quality gate steps and failure handling

## Description

更新 `quality_gate.go` 中的 gate 步骤和失败处理：Step 2 使用 `unit-test`，Step 3 使用 `test`。重构 `addFixTask` 移除硬编码 recipe 名映射，改为通用规则 `step → "just " + step`。更新 `handleGateFailure` 中 guide/label map 的 `"e2e-test"` 键为 `"test"`。迁移 `runE2ERegression` 函数。

## Reference Files
- `proposal.md#Proposed-Solution` — defines gate sequence steps (compile → fmt → lint → unit-test → test → probe) and two invocation paths
- `proposal.md#Constraints-&-Dependencies` — addFixTask generic rule, handleGateFailure guide/label map migration
- `proposal.md#runE2ERegression-迁移要点` — detailed 5-step migration plan for runE2ERegression function
- `proposal.md#Success-Criteria` — criteria for addFixTask generic rule, handleGateFailure migration, and no-recipe error handling

## Acceptance Criteria
- Step 2 (unit test step) runs `just unit-test`
- Step 3 (advanced test step) runs `just test`
- `addFixTask` uses generic rule `step → "just " + step`, no hardcoded recipe name mapping
- `handleGateFailure` guide/label map has `"test"` key instead of `"e2e-test"`
- `runE2ERegression` migrated: `e2e-setup` → `test-setup`, `e2e-test` → `test`, function renamed to `runTestRegression` or inlined
- When `unit-test` recipe missing, quality gate reports error suggesting run `init-justfile` (no fallback)

## Hard Rules
- Gate Sequence 路径无 fallback——recipe 不存在时 gate 报错，不回落到其他 recipe

## Implementation Notes
- `addFixTask` 中移除硬编码 `step=="unit-test" → "just test"` 映射后，改为直接使用 step 名：`step → "just " + step`
- `runE2ERegression` 函数内的 `just dev` 调用保持不变（非重命名范围）
- 日志/错误信息中的 `e2e` 引用需更新为 `test`
