---
id: "6"
title: "更新 quality gate 下游 consumer"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [2]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 6: 更新 quality gate 下游 consumer

## Description

根据提案的 Recipe 命名统一模型，更新 quality gate 的两个下游 consumer：
1. `run-tests` skill：移除 fallback 链（`<key>-compile` 不存在时 fallback 到 `compile`），直接调用 `<key>-compile` 等 surface recipe
2. `forge quality-gate` Go binary：将硬编码的 recipe 调用改为按 surface key 拼接 recipe 名（最小改动，不做完整动态解析重构）

## Reference Files
- `docs/proposals/init-justfile-slim/proposal.md` — Consumer Impact, 行动项 (#4, #5)
- `forge-cli/pkg/just/just.go` — ResolvePrefixedRecipe 函数（包含需移除的 fallback 逻辑）
- `forge-cli/internal/cmd/qualitygate/quality_gate_lifecycle.go` — resolveRecipe 函数（同含 fallback 逻辑）
- `plugins/forge/skills/run-tests/SKILL.md` — run-tests skill（需移除 fallback 链描述）

## Acceptance Criteria
- [ ] `forge-cli/pkg/just/just.go` 中 `ResolvePrefixedRecipe` 移除 generic fallback：有 scope 时仅返回 `<scope>-<recipe>`，不回退到无前缀 recipe
- [ ] `forge-cli/internal/cmd/qualitygate/quality_gate_lifecycle.go` 中 `resolveRecipe` 移除 generic fallback，行为与 ResolvePrefixedRecipe 保持一致
- [ ] `run-tests` SKILL.md 中移除 fallback 链描述（如 `<key>-compile` 不存在 → fallback `compile`），改为直接调用 surface recipe
- [ ] 单 surface scalar 项目（无 scope）仍使用无前缀 recipe（compile、unit-test 等），功能不变

## Hard Rules
- 修改 `plugins/forge/` 下的文件前，必须先读 `docs/conventions/forge-distribution.md` 了解分发模型
- 范围外：不做完整的 recipe 名动态解析机制重构，本次仅做最小改动（移除 fallback）

## Implementation Notes
- 当前 `ResolvePrefixedRecipe` 逻辑：有 scope 时先尝试 `<scope>-<recipe>`，不存在则 fallback 到 `recipe`；需改为有 scope 时仅返回 `<scope>-<recipe>` 或空字符串
- 当前 `resolveRecipe` 逻辑同上，需保持两个函数行为一致
- scalar surface（scope 为空）走原有无前缀路径，不受影响

### Test Impact
- Affected test suite(s): `forge-cli/pkg/just/`, `forge-cli/internal/cmd/qualitygate/`
- Expected fixture changes: fallback 相关的 test case 需更新期望行为
- Risk level: medium
