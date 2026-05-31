---
id: "8"
title: "Extract magic values to named constants"
priority: "P1"
estimated_time: "3h"
complexity: "high"
dependencies: [4, 6]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 8: Extract magic values to named constants

## Description
Phase 2b：提取所有魔法值为命名常量——路径字符串、颜色值、超时值、哨兵数、八进制权限。统一使用 `0o` 前缀。依据 Task 4 产出的 `constants.md` 规范执行。

## Reference Files
- forge-cli/internal/cmd/quality_gate.go: `"tests/results/raw-output.txt"` 出现 2 次，重试参数 `3` 次、`5*time.Second` (source: proposal.md#Evidence)
- forge-cli/internal/cmd/init.go:217: 颜色值 `#7DCFFF` (source: proposal.md#Evidence)
- forge-cli/internal/cmd/init_surfaces.go:17,20,23: 颜色值 `#FF8700`、`#7DCFFF`、`#9ECE6A` (source: proposal.md#Evidence)
- forge-cli/internal/cmd/task/list.go:442: 哨兵数 `99999` (source: proposal.md#Evidence)
- forge-cli/internal/cmd/task/claim.go:376: 哨兵数 `99999` (source: proposal.md#Evidence)

## Acceptance Criteria
- [ ] `grep -rn '"tests/results/' forge-cli/internal/ forge-cli/pkg/` 返回零结果（SC-1）
- [ ] `grep -rn 'lipgloss.Color("#' forge-cli/internal/ forge-cli/pkg/` 返回零结果（SC-2）
- [ ] `grep -rn '\b99999\b' forge-cli/` 返回零结果（SC-3）
- [ ] `grep -rn '0644\|0755' forge-cli/internal/ forge-cli/pkg/` 返回零结果（SC-4，统一 `0o` 前缀）
- [ ] `go build ./...` 和 `go test ./...` 全部通过（SC-11）

## Implementation Notes
- 路径常量集中在使用它的包内的 `constants.go` 文件中
- 颜色常量集中在 `internal/cmd/` 下的常量文件中（如 `internal/cmd/colors.go`）
- 哨兵常量命名为 `maxDepthSentinel` 或类似语义名
- 八进制权限统一使用 `0o644`、`0o755` 格式

### Test Impact
- Affected test suite(s): `forge-cli/internal/cmd/`, `forge-cli/pkg/task/`
- Expected fixture changes: 测试中的硬编码路径/颜色值需同步更新
- Risk level: medium（批量替换，需逐一验证）
