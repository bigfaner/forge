---
id: "11"
title: "Update forge-cli test files for surface-key naming"
priority: "P1"
estimated_time: "1-2h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 11: Update forge-cli test files for surface-key naming

## Description
forge-cli 中引用 `gen-test-scripts-{type}` 命名模式的测试用例需要更新为 `{key}` 命名。同时验证所有相关测试通过。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Evidence, Scope, Success Criteria
- `forge-cli/pkg/task/pipeline_test.go` — pipeline expansion 测试 (ref: Evidence)
- `forge-cli/pkg/task/autogen_test.go` — 自动生成任务测试 (ref: Scope)

## Acceptance Criteria
- [ ] 所有引用 `gen-test-scripts-{type}` 命名的测试用例更新为 `{key}` 命名
- [ ] 测试 fixture 和期望值更新为 per-surface-key expansion
- [ ] `go test ./...` 全部通过

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/`, `forge-cli/internal/cmd/`, `forge-cli/tests/`
- Expected fixture changes: pipeline expansion 期望值、task key/id 模式
- Risk level: medium

需 grep `gen-test-scripts` 和 `gen-scripts` 确认所有引用位置，重点关注 expansion 模式的期望输出。
