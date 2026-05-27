---
id: "3"
title: "更新 quality_gate 路径和 mobile-test-setup 集成"
priority: "P0"
estimated_time: "2h"
dependencies: [1, 2]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 3: 更新 quality_gate 路径和 mobile-test-setup 集成

## Description
修复 `quality_gate.go` 中残留的旧路径引用（`tests/e2e/results/raw-output.txt`），替换所有 `GetE2EStagingDir`/`GetE2EGraduatedMarker` 调用，将 `mobile-test-setup` 集成到 `runSurfaceLifecycle()` 流程中，并更新过时注释。

## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 第 4 项定义了 quality_gate 的具体修复点
- `proposal.md#Scope` — In Scope 第 6 项要求 mobile-test-setup 集成
- `proposal.md#Success-Criteria` — 验证条件：quality_gate.go 包含 mobile-test-setup 调用

## Acceptance Criteria
- [ ] `runTestRegressionSurface` 中 `tests/e2e/results/raw-output.txt` 路径已替换为新路径
- [ ] 所有 `GetE2EStagingDir` / `GetE2EGraduatedMarker` 调用已替换
- [ ] `runSurfaceLifecycle()` 包含 `mobile-test-setup` 调用（当 surface type 为 mobile 时）
- [ ] 注释 "promoted scripts in tests/e2e/" 已更新
- [ ] `grep -rn "graduated\|graduation" forge-cli/pkg/ forge-cli/internal/ --include="*.go"`（排除 _test.go、config 迁移代码）返回 0 结果
- [ ] `go build ./...` 通过

## Implementation Notes
- MR #173 只修了 legacy 路径，surface 路径遗漏 — 本次修复 surface 路径
- `mobile-test-setup` 仅在 surface type 为 mobile 时调用，需条件判断
- `grep -rn "staging" forge-cli/pkg/ forge-cli/internal/ --include="*.go"` 也应返回 0 结果

### Integration Test Impact
- Affected test suite(s): `forge-cli/internal/`, `forge-cli/pkg/feature/`
- Expected fixture changes: 涉及 quality_gate 行为的测试 fixture
- Risk level: medium
