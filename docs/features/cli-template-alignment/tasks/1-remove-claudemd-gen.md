---
id: "1"
title: "Remove CLAUDE.md generation from forge init"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 1: Remove CLAUDE.md generation from forge init

## Description
`forge init` 生成的 CLAUDE.md 模板存在三层冗余（与 prompt 模板重复、与 hook 注入的 guide.md 重复、无独特价值）。移除 init 流程中的 CLAUDE.md 生成步骤，同时更新 init 命令描述以反映变更。

## Reference Files
- `forge-cli/internal/embedded/claudemd_template.md`: 删除此嵌入模板文件 (source: proposal.md#Proposed-Solution item 1)
- `forge-cli/internal/embedded/claudemd.go`: 删除 embed 声明和导出变量 CLAUDEmdTemplate (source: proposal.md#Proposed-Solution item 2)
- `forge-cli/internal/embedded/claudemd_test.go`: 删除关联测试文件
- `forge-cli/internal/cmd/init.go`: 移除 createCLAUDEmd() 函数及其调用，更新 initCmd Long 描述 (source: proposal.md#Proposed-Solution items 3-4)

## Acceptance Criteria
- [ ] `claudemd_template.md`、`claudemd.go`、`claudemd_test.go` 已删除
- [ ] `init.go` 中 `createCLAUDEmd()` 函数及其调用已移除
- [ ] 全局搜索确认无残留的 `CLAUDEmdTemplate` 和 `claudemd` 引用（除已删除文件外）
- [ ] `initCmd` 的 Long 描述已移除 CLAUDE.md 相关语句

## Implementation Notes
- 提案已确认行为准则由 prompt 模板覆盖（`coding-*.md`），Forge 上下文由 hook 注入的 `<forge-guide>` 覆盖，模板无独特价值
- 删除文件前需确认 `claudemd.go` 的导出变量在项目内无其他引用（提案 Key Risk 提及 `CLAUDEmdTemplate` 和 `claudemd` 全局搜索）
- `claudemd_test.go` 是 `claudemd.go` 的测试文件，需一并删除
