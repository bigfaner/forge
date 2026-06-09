---
id: "3"
title: "scaffold 命令单元测试"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [2]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: scaffold 命令单元测试

## Description

为 `forge justfile scaffold` 命令编写全面的单元测试，覆盖单 surface 模式（5 种 surface type）和 aggregate 模式，以及参数校验的边界场景。满足成功标准 #3：每个 surface type 至少 1 个 test case + 聚合模式 test case。

## Reference Files
- `docs/proposals/init-justfile-slim/proposal.md` — 成功标准 (#3)
- `forge-cli/internal/cmd/scaffold/` — 被测代码（Task 1、2 产出）
- `forge-cli/internal/cmd/qualitygate/quality_gate_test.go` — 现有测试风格参考

## Acceptance Criteria
- [ ] 每个 surface type（cli/tui/api/web/mobile）有 ≥1 个 test case 验证生成的 recipe 集结构完整
- [ ] 聚合模式有 test case 验证 install/ci/clean 生成正确性
- [ ] 边界场景：unknown surface type 报错、scalar surface 传 --key 报错、named surface 缺 --key 报错
- [ ] 验证所有生成 recipe 的占位符使用 `<<...>>` 语法，不包含 `{{...}}`
- [ ] 所有测试通过 `go test -race -cover ./forge-cli/internal/cmd/scaffold/...`

## Implementation Notes
- 遵循 forge-cli CLAUDE.md 的 TDD 规范：table-driven tests，coverage target 80%+
- 参考 `quality_gate_test.go` 中的测试风格
- 建议使用 `golden file` 或字符串断言验证 recipe 输出
