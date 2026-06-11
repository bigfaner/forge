---
id: "3"
title: "Fix CLI help text for task/feature/quality-gate commands"
priority: "P0"
estimated_time: "1h"
dependencies: []
complexity: "medium"
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 3: Fix CLI help text for task/feature/quality-gate commands

## Description
修复 6 个 CLI 命令的 cobra Long/Short 描述，使其准确反映代码实际行为（C2-C6, C11）。仅修改 cobra Command 定义的字符串常量，不涉及逻辑变更。

## Reference Files
- `docs/proposals/cli-doc-accuracy-audit/proposal.md` — In Scope (CLI Help Text 修复 C2-C6, C11), Success Criteria, Key Risks
- `forge-cli/internal/cmd/cleanup.go`: cleanup Short/Long 补充 blocked/suspended/rejected 状态 (ref: In Scope)
- `forge-cli/internal/cmd/task/validate.go`: validate Long 补充全部 12+ 项验证 (ref: In Scope)
- `forge-cli/internal/cmd/task/add.go`: add Long 补充使用概述 (ref: In Scope)
- `forge-cli/internal/cmd/qualitygate/quality_gate.go`: quality-gate Long 补充副作用描述 (ref: In Scope)

## Acceptance Criteria
- [ ] `forge cleanup --help` 的 Long 描述包含 blocked/suspended/rejected 状态的清理行为，与 cleanup RunE 实现一致
- [ ] `forge task claim --help` 的 Long 描述包含 auto-unblock 行为说明，与 claim RunE 实现一致
- [ ] `forge task validate --help` 的 Long 描述列出全部 12+ 项验证步骤，覆盖实际执行的验证逻辑
- [ ] `forge task add --help` 的 Long 描述包含使用概述（fix-task vs 普通任务的参数差异）
- [ ] `forge feature --help` 的 Long 描述包含 `set` 子命令和行为差异说明
- [ ] `forge quality-gate --help` 的 Long 描述包含 fix task 自动创建、retry-once、docs-only 跳过等副作用

## Implementation Notes
- 修改前需阅读对应命令的 RunE 函数，确保 Long 描述中提及的每一项功能都有代码支撑
- C4（validate）修改量最大：需在 Long 中列出全部 12+ 项验证，先审计 validate.go 中所有 validate* 函数
- C11（quality-gate）需阅读 quality_gate_fix_task.go 了解 fix task 自动创建逻辑
- Complexity override: 规则判定 AC=6 → high，但所有修改遵循相同模式（更新 cobra Long/Short 字符串），实际复杂度 medium
