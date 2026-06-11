---
status: "completed"
started: "2026-05-30 21:03"
completed: "2026-05-30 21:06"
time_spent: "~3m"
---

# Task Record: 3 编写架构概览文档 architecture-overview.md

## Summary
创建 docs/user-guide/architecture-overview.md，以用户视角介绍 Forge 插件机制、四大组件角色（Skill/Command/Agent/Hook）、数据流向（三层 ASCII 图）、状态管理（Feature 和任务状态流转）、目录约定（完整项目结构）以及工作模式（完整模式 vs 快速模式）。所有内容基于用户视角，不包含 Go 包结构等开发者内部实现细节。

## Changes

### Files Created
- docs/user-guide/architecture-overview.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
~368 lines, 7 sections, 4 tables, 6 ASCII diagrams

## Referenced Documents
- docs/ARCHITECTURE.md
- README.md
- docs/conventions/forge-distribution.md
- .forge/config.yaml

## Review Status
final

## Acceptance Criteria
- [x] 包含插件机制说明（Claude Code 插件加载方式和 Forge 的定位）
- [x] 包含四大组件角色表格（skill、command、agent、hook），每个有名称、用途、触发方式
- [x] 包含数据流向图解（从用户输入 → Forge 处理 → 文件系统变更的可视化说明）
- [x] 包含目录约定说明（.forge/ 目录结构、docs/features/ 结构、manifest.md 作用）
- [x] 不包含 Go 包结构、CLI 内部命令注册、ResolveScope 等开发者内部实现细节

## Notes
文档使用中文，包含 Mermaid 兼容的 ASCII 图解，顶部标注最后更新日期和版本 v3.0.0
