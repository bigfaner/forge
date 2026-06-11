---
status: "completed"
started: "2026-06-02 22:47"
completed: "2026-06-02 22:49"
time_spent: "~2m"
---

# Task Record: 9 Fix stale test-type-model references in gen-journeys + run-tests + init-justfile

## Summary
Fixed 16 stale references to deleted docs/reference/test-type-model.md across 3 skills (gen-journeys, run-tests, init-justfile). Replaced with self-contained inline test type descriptions per Hard Rules. Also fixed 'UI tests' to 'Web tests' in result-parsing.md.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-journeys/rules/surface-cli.md
- plugins/forge/skills/gen-journeys/rules/surface-api.md
- plugins/forge/skills/gen-journeys/rules/surface-tui.md
- plugins/forge/skills/gen-journeys/rules/surface-web.md
- plugins/forge/skills/gen-journeys/rules/surface-mobile.md
- plugins/forge/skills/run-tests/rules/surfaces/cli.md
- plugins/forge/skills/run-tests/rules/surfaces/api.md
- plugins/forge/skills/run-tests/rules/surfaces/tui.md
- plugins/forge/skills/run-tests/rules/surfaces/web.md
- plugins/forge/skills/run-tests/rules/surfaces/mobile.md
- plugins/forge/skills/run-tests/rules/result-parsing.md
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md

### Key Decisions
无

## Document Metrics
16 files modified, 16 stale references removed, 0 cross-skill references remaining

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] gen-journeys 5 个 surface rule 文件中 test-type-model.md 引用已移除，替换为自包含 test type 声明
- [x] run-tests 5 个 surface rule 文件中旧路径引用已移除
- [x] run-tests result-parsing.md 中 'UI tests' 改为 'Web tests'
- [x] init-justfile 5 个 surface rule 文件中旧路径引用已移除
- [x] 所有修改后的文件不包含任何对其它 skill 内部文件的引用

## Notes
gen-journeys: replaced external reference with inline self-contained test type description per surface. run-tests/init-justfile: removed reference line entirely (files already had complete surface-specific info). result-parsing.md: fixed UI tests -> Web tests terminology.
