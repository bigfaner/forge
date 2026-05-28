---
id: "2"
title: "autoBehaviorPrompts 新增 4 个 eval 提示"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: autoBehaviorPrompts 新增 4 个 eval 提示

## Description
在 `autoBehaviorPrompts()` 中新增 4 个 eval 相关的交互提示，让用户通过 `forge init`（TUI）和 `forge config init`（stdin）设置 eval 自动运行策略。当前 init 流程覆盖 6 类自动行为但缺少 eval 配置入口。

## Reference Files
- `forge-cli/internal/cmd/init.go`: autoBehaviorPrompts() (line 257-338) 新增 4 个 eval prompt，插入到 gitPush 之前 (source: proposal.md#Part-2)

## Acceptance Criteria
- [ ] `forge init` 交互流程包含 4 个 eval 提示（proposal/prd/uiDesign/techDesign）
- [ ] 提示默认值与 `AutoConfigDefaults()` 一致：proposal:true, prd:false, uiDesign:true, techDesign:false
- [ ] 4 个 eval 提示位于 gitPush 提示之前
- [ ] `forge config init`（stdin 版本）通过共享 `askAutoBehavior()` 同步覆盖

## Implementation Notes
与 gitPush 的单开关风格一致，无 quick/full 前缀。提示文本参考 proposal Part 2 表格：
1. "Auto-evaluate proposals?" → auto.eval.proposal, default true
2. "Auto-evaluate PRD documents?" → auto.eval.prd, default false
3. "Auto-evaluate UI designs?" → auto.eval.uiDesign, default true
4. "Auto-evaluate tech designs?" → auto.eval.techDesign, default false
