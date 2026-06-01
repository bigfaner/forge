---
status: "completed"
started: "2026-06-01 21:32"
completed: "2026-06-01 21:36"
time_spent: "~4m"
---

# Task Record: 5 gen-test-scripts/types/ui: Audit existing surface guard for consistency

## Summary
Audited and upgraded gen-test-scripts/types/ui.md Sitemap Resolution guard: replaced declarative 'web-ui interface' phrasing with explicit forge surfaces --json detection pattern consistent with Task 1-4.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/types/ui.md

### Key Decisions
无

## Document Metrics
consistency: 100 (all 7 files now use forge surfaces --json + web surface type pattern)

## Referenced Documents
- docs/proposals/sitemap-surface-guard/proposal.md
- plugins/forge/skills/gen-web-sitemap/SKILL.md
- plugins/forge/skills/write-prd/rules/ui-functions.md
- plugins/forge/skills/write-prd/rules/self-check.md
- plugins/forge/skills/breakdown-tasks/rules/ui-placement.md
- plugins/forge/skills/eval/rules/validate-ux-pipeline.md

## Review Status
final

## Acceptance Criteria
- [x] 审查已有守卫是否与 forge surfaces --json 检测方式一致
- [x] 记录审查结论：守卫不充分，已实施修改并记录变更内容
- [x] 守卫措辞与 Task 1-4 中建立的模式保持一致

## Notes
Two gaps found in original guard: (1) 'web-ui interface' phrasing inconsistent with Task 1-4's 'web' surface type; (2) lacked explicit forge surfaces --json CLI call and PASS/SKIP decision logic. Both resolved by replacing with the established pattern.
