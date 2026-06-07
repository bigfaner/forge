---
status: "completed"
started: "2026-06-07 23:48"
completed: "2026-06-07 23:51"
time_spent: "~3m"
---

# Task Record: 2 Surface rules 增加 compile/fmt/lint/unit-test gate recipe 模板

## Summary
Added compile/fmt/lint/unit-test gate recipe stubs and Recipe Invocation Contract entries to all 5 surface rule files (api.md, web.md, cli.md, tui.md, mobile.md)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md

### Key Decisions
无

## Document Metrics
5 files modified, 4 gate recipes per file (20 total stubs + 20 contract entries)

## Referenced Documents
- docs/proposals/per-task-surface-scoped-gate/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] 5 surface rule files contain <key>-compile/<key>-fmt/<key>-lint/<key>-unit-test stub recipe definitions
- [x] 5 surface rule files contain compile/fmt/lint/unit-test Recipe Invocation Contract entries

## Notes
Each gate recipe stub includes dual-platform variants (linux+windows) and surface isolation constraint comments. Implementation notes fully observed: key-based prefixing, constraint annotations, consistent template pattern across all files including mobile.md.
