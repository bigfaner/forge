---
id: "5"
title: "Update 4 skill files with config check bash template"
priority: "P1"
estimated_time: "1h"
dependencies: ["1", "2", "3"]
type: "doc"
mainSession: false
---

# 5: Update 4 skill files with config check bash template

## Description
在 4 个 skill 文件（brainstorm、write-prd、tech-design、ui-design）中，用统一的 bash config check 模板替换各 skill 中 eval 触发点的 `AskUserQuestion` 步骤。模板读取 mode 和 eval 配置，自动运行或跳过 eval，CLI 不可用时回退到 AskUserQuestion。

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/brainstorm/SKILL.md` | Step 7（eval-proposal AskUserQuestion）→ config check bash 模板 |
| `plugins/forge/skills/write-prd/SKILL.md` | Step 11（eval-prd AskUserQuestion）→ config check bash 模板 |
| `plugins/forge/skills/tech-design/SKILL.md` | Step 10（eval-design AskUserQuestion）→ config check bash 模板 |
| `plugins/forge/skills/ui-design/SKILL.md` | Step 7（无条件 eval-ui）→ config check bash 模板 |

## Reference Files
- `docs/proposals/auto-eval-config/proposal.md#Constraints-Dependencies` — bash config check 模板代码（含 CLI fallback 逻辑）
- `docs/proposals/auto-eval-config/proposal.md#Key-Scenarios` — 自动运行/手动确认/无 feature 上下文 3 种场景
- `docs/proposals/auto-eval-config/proposal.md#Scope` — 4 个 skill 修改的具体步骤位置

## Acceptance Criteria
- [ ] brainstorm 在 `auto.eval.proposal` 为 true 时跳过 AskUserQuestion，直接运行 eval-proposal
- [ ] brainstorm 在 `auto.eval.proposal` 为 false 时保持 AskUserQuestion
- [ ] write-prd 在 `auto.eval.prd` 为 true 时跳过 AskUserQuestion，直接运行 eval-prd
- [ ] write-prd 在 `auto.eval.prd` 为 false 时保持 AskUserQuestion
- [ ] ui-design 在 `auto.eval.uiDesign` 为 true 时跳过 AskUserQuestion，直接运行 eval-ui
- [ ] ui-design 在 `auto.eval.uiDesign` 为 false 时保持 AskUserQuestion
- [ ] tech-design 在 `auto.eval.techDesign` 为 true 时跳过 AskUserQuestion，直接运行 eval-design
- [ ] tech-design 在 `auto.eval.techDesign` 为 false 时保持 AskUserQuestion
- [ ] 4 个 skill 使用相同的 config check 模板（代码审查验证一致性）
- [ ] CLI 不可用时（退出码非零）回退到 AskUserQuestion

## Hard Rules
- 加载 `docs/conventions/forge-distribution.md` 再修改 plugin 文件
- 使用 EXTREMELY-IMPORTANT 标注 config check 逻辑块
- bash 模板中 `{skillKey}` 替换为各 skill 对应的 eval 配置名：brainstorm→proposal、write-prd→prd、ui-design→uiDesign、tech-design→techDesign

## Implementation Notes
- ui-design 当前是无条件运行 eval-ui，改为 config check 后默认值保持 ON（`quick:true full:true`），行为不变但实现路径统一
- 模板中 `forge config get auto.eval.{skillKey}.$MODE` 使用 mode 检测结果动态选择 quick/full 子键
- `forge config get mode` 返回 `"none"` 时回退到 AskUserQuestion
