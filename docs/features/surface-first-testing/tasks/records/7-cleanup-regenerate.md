---
status: "completed"
started: "2026-06-02 22:13"
completed: "2026-06-02 22:15"
time_spent: "~2m"
---

# Task Record: 7 清理旧文件并重新生成 Forge 项目 conventions

## Summary
Deleted 6 old framework convention files (ginkgo/go/junit/pytest/rust/vitest) + old test-type-model.md, regenerated surface-first convention files using test-guide skill: docs/conventions/testing/index.md, docs/conventions/testing/cli/index.md, docs/conventions/testing/cli/core.md

## Changes

### Files Created
- docs/conventions/testing/index.md
- docs/conventions/testing/cli/index.md
- docs/conventions/testing/cli/core.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
3 files created (~90 lines), 8 files deleted, core.md contains all 7 mandatory sections + assertion preference table

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] docs/conventions/testing/ 下 6 个旧框架文件已删除
- [x] docs/reference/test-type-model.md 已删除
- [x] docs/conventions/testing/cli/ 已用新 test-guide 重新生成（含 index.md + core.md）
- [x] 顶层 docs/conventions/testing/index.md 已重新生成

## Notes
Framework detected: Go testing + testify (high confidence). Assertion preference table filled with testify/assert row. Content of old test-type-model.md verified as fully migrated to plugins/forge/skills/test-guide/references/test-type-model.md before deletion.
