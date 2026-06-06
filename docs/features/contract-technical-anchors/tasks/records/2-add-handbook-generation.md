---
status: "completed"
started: "2026-06-05 17:25"
completed: "2026-06-05 17:28"
time_spent: "~3m"
---

# Task Record: 2 tech-design 增加全 surface handbook 自动生成

## Summary
扩展 tech-design skill 支持全 surface handbook 自动生成：新增 cli-handbook.md、page-map.md、screen-map.md 三个模板，修改 SKILL.md 增加 Pipeline Configuration 列、Step 8 Artifact 表和 Surface Detection 段落、Step 9 Manifest 更新逻辑

## Changes

### Files Created
- plugins/forge/skills/tech-design/templates/cli-handbook.md
- plugins/forge/skills/tech-design/templates/page-map.md
- plugins/forge/skills/tech-design/templates/screen-map.md

### Files Modified
- plugins/forge/skills/tech-design/SKILL.md

### Key Decisions
无

## Document Metrics
3 templates created (~45 lines each), 1 SKILL.md updated (Pipeline table +3 cols, Step 8 +4 artifact rows + Surface Detection section, Step 9 +6 manifest rows)

## Referenced Documents
- docs/proposals/contract-technical-anchors/proposal.md
- plugins/forge/skills/tech-design/SKILL.md
- plugins/forge/skills/tech-design/templates/api-handbook.md

## Review Status
final

## Acceptance Criteria
- [x] tech-design skill 根据 surface 配置自动选择并生成对应 handbook
- [x] cli-handbook 模板覆盖 CLI/TUI surface 的命令、子命令、参数、别名
- [x] page-map 模板覆盖 Web surface 的页面名称、路由、布局、认证要求
- [x] screen-map 模板覆盖 Mobile surface 的屏幕名称、导航路径、deeplink、平台
- [x] 新 handbook 模板格式与 api-handbook 保持一致的 frontmatter 和 section 结构

## Notes
模板格式严格对齐 api-handbook 的 frontmatter (created, related) 和 section 结构模式。SKILL.md Pipeline Configuration 表从 6 列扩展到 9 列，新增 CLI Handbook / Page Map / Screen Map。Step 8 新增 Surface Detection for Handbook Generation 段落，包含 forge surfaces 检测命令和 surface-to-template mapping 表。
