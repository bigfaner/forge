---
id: "4"
title: "Fix CLI help text for init/forensic/worktree/fact commands"
priority: "P1"
estimated_time: "1h"
dependencies: []
complexity: "medium"
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 4: Fix CLI help text for init/forensic/worktree/fact commands

## Description
修复 5 个 CLI 命令的 cobra Long/Short 描述，使其准确反映代码实际行为（C1, C7-C10）。仅修改 cobra Command 定义的字符串常量，不涉及逻辑变更。

## Reference Files
- `docs/proposals/cli-doc-accuracy-audit/proposal.md` — In Scope (CLI Help Text 修复 C1, C7-C10), Success Criteria
- `forge-cli/internal/cmd/init.go`: init Long 补充 surface detection 步骤 (ref: In Scope)
- `forge-cli/internal/cmd/forensic/search.go`: forensic search 新增 Long 描述 (ref: In Scope)
- `forge-cli/internal/cmd/forensic/subagents.go`: forensic subagents 新增 Long 描述 (ref: In Scope)
- `forge-cli/internal/cmd/worktree/cmd_status.go`: worktree status Long 补充 UNPUSHED 字段 (ref: In Scope)

## Acceptance Criteria
- [ ] `forge init --help` 的 Long 描述补充 surface detection 步骤说明
- [ ] `forge forensic search --help` 新增 Long 描述（当前为空），描述搜索行为和输出格式
- [ ] `forge forensic subagents --help` 新增 Long 描述（当前为空），描述子代理追踪行为
- [ ] `forge worktree status --help` 的 Long 描述补充 UNPUSHED 字段含义
- [ ] `forge fact summary --help` 的 Long 描述补充 COVERAGE 指标说明

## Implementation Notes
- C7/C8（forensic）完全没有 Long 描述，需先阅读对应 RunE 函数理解完整行为后再撰写
- C1（init）需阅读 init.go 和 init_surfaces.go 了解 surface detection 的完整流程
- Complexity override: 规则判定 AC=5 → high，但所有修改遵循相同模式（更新 cobra Long 字符串），实际复杂度 medium
