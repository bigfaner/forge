---
status: "completed"
started: "2026-05-30 20:59"
completed: "2026-05-30 21:02"
time_spent: "~3m"
---

# Task Record: 2 编写初始化文档 initialization.md

## Summary
创建 docs/user-guide/initialization.md 初始化指南文档，涵盖 forge init 完整流程（6步）、config.yaml 全字段参考（28个配置项）、Surface 检测机制说明和端到端示例

## Changes

### Files Created
- docs/user-guide/initialization.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
~350 lines, 28 config fields, 5 AC criteria met

## Referenced Documents
- .forge/config.yaml
- plugins/forge/commands/init-forge.md
- docs/ARCHITECTURE.md
- docs/reference/test-type-model.md
- README.md
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/surfaces_detect.go
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/forgeconfig/detect_surface.go

## Review Status
final

## Acceptance Criteria
- [x] 包含 forge init 的完整流程说明（从命令执行到项目就绪）
- [x] 包含 config.yaml 全字段表格，至少 8 个配置项，每个字段有名称、类型、默认值、说明
- [x] 包含 Surface 检测机制说明（forge surfaces detect 的使用和结果解读）
- [x] 包含首个项目设置的端到端示例（从 init 到可以开始使用 Forge）
- [x] 所有代码示例可直接复制执行，无需额外修改

## Notes
文档使用中文，顶部标注最后更新日期 2026-05-30 和版本 v3.0.0。config.yaml 字段表格从实际源码（forgeconfig/config.go）提取，覆盖 version、project-type、auto（含 13 个子字段）、worktree、coverage、surfaces、execution-order、test-framework、languages 共 28 个配置项。Surface 检测说明覆盖 forge surfaces detect 和 forge surfaces detect --apply 两种用法。
