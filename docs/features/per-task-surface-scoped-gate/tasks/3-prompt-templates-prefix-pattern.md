---
id: "3"
title: "Prompt templates gate recipe 从参数模式切换为前缀模式"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: Prompt templates gate recipe 从参数模式切换为前缀模式

## Description
Task 1 在 Go CLI 路径（`RunGate()`）实现了 prefixed recipe 解析，但 Agent 路径（LLM 读 prompt template 后直接执行 `just` 命令）仍使用参数模式 `just compile backend`。`mixed.just` 的 `compile` recipe 不接受参数，导致 Agent 在任务执行中的预检 compile 会失败。需要将所有 prompt template 中的 gate recipe 调用从参数模式切换为前缀模式。

## Reference Files
- `docs/proposals/per-task-surface-scoped-gate/proposal.md` — Proposed Solution, Success Criteria
- `forge-cli/pkg/prompt/templates/coding-feature.md`: 3-step gate 参数→前缀 (ref: Proposed Solution)
- `forge-cli/pkg/prompt/templates/gate.md`: 4-step gate 参数→前缀 (ref: Proposed Solution)
- `forge-cli/pkg/prompt/templates/coding-refactor.md`: 增量 compile + final gate 共 9 处 (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/templates/coding-feature.md` | gate recipe 从参数模式 `just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` 切换为前缀模式 `just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}compile` |
| `forge-cli/pkg/prompt/templates/coding-enhancement.md` | 同上（3 处） |
| `forge-cli/pkg/prompt/templates/coding-fix.md` | 同上（3 处） |
| `forge-cli/pkg/prompt/templates/coding-cleanup.md` | 同上（3 处） |
| `forge-cli/pkg/prompt/templates/coding-refactor.md` | 同上（9 处：6 处增量 compile + 3 处 final gate） |
| `forge-cli/pkg/prompt/templates/gate.md` | 同上（4 处，含 unit-test） |
| `forge-cli/pkg/prompt/templates/validation-code.md` | 同上（4 处） |
| `forge-cli/pkg/prompt/templates/fix-record-missed.md` | 同上（4 处） |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 8 个 prompt template 文件中所有 33 处 gate recipe 调用从参数模式切换为前缀模式：`just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}<recipe>`
- [ ] SurfaceKey 为空时渲染结果仍为 `just compile`（无前缀），与改动前一致
- [ ] 通过 grep 验证所有模板中不再包含旧参数模式 `{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` 出现在 recipe 调用上下文中

## Implementation Notes
- 替换模式：`just <recipe>{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` → `just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}<recipe>`
- 4 个 recipe 名称：`compile`、`fmt`、`lint`、`unit-test`，每个文件涉及的 recipe 子集不同
- `coding-refactor.md` 有 9 处引用：lines 189, 212, 218, 230, 237, 244（增量 compile）+ lines 262, 263, 264（final gate 的 compile/fmt/lint）
- 此修改与 Task 1 的 `RunGate()` prefixed resolution 形成互补：Go CLI 路径和 Agent 路径使用相同的前缀模式
- 注意：旧参数模式在非 `just` 上下文中的 `{{if .SurfaceKey}}` 引用不要误改（如 identity 段落的 SurfaceKey 声明）
