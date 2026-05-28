---
status: "completed"
started: "2026-05-28 00:51"
completed: "2026-05-28 01:02"
time_spent: "~11m"
---

# Task Record: 2 迁移 21 个 prompt 模板占位符语法并添加条件块

## Summary
Migrated all 21 prompt template files from {{X}} placeholder syntax to {{.X}} dot-notation for text/template compatibility. Added {{if}} conditional blocks for PhaseSummary (18 templates, two independent positions), CoverageStrategy (5 coding templates), SurfaceKey label line and just commands (18 templates), and Complexity (4 coding templates replacing <!-- IF NOT_LOW --> markers). Added SURFACE_KEY header to 4 doc templates. Updated 2 test assertions to match new dot-notation format.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/data/code-quality-simplify.md
- forge-cli/pkg/prompt/data/doc.md
- forge-cli/pkg/prompt/data/doc-consolidate.md
- forge-cli/pkg/prompt/data/doc-drift.md
- forge-cli/pkg/prompt/data/doc-review.md
- forge-cli/pkg/prompt/data/doc-summary.md
- forge-cli/pkg/prompt/data/eval-contract.md
- forge-cli/pkg/prompt/data/eval-journey.md
- forge-cli/pkg/prompt/data/fix-record-missed.md
- forge-cli/pkg/prompt/data/gate.md
- forge-cli/pkg/prompt/data/test-gen-contracts.md
- forge-cli/pkg/prompt/data/test-gen-journeys.md
- forge-cli/pkg/prompt/data/test-gen-scripts.md
- forge-cli/pkg/prompt/data/test-run.md
- forge-cli/pkg/prompt/data/validation-code.md
- forge-cli/pkg/prompt/data/validation-ux.md
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- PhaseSummary label line uses {{if .PhaseSummary}}{{.PhaseSummary}}{{end}} because Go code already prepends 'PHASE_SUMMARY: ' prefix to the value
- SurfaceKey label line uses {{if .SurfaceKey}}SURFACE_KEY: {{.SurfaceKey}}{{end}} to eliminate both label-only and trailing-space issues
- just commands use just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}} pattern for trailing-space-free rendering
- CoverageStrategy uses {{if .CoverageStrategy}}...{{end}} wrapping the entire IMPORTANT block, Go code handles three-state logic (empty/maintain/percentage)
- 4 coding templates use {{if ne .Complexity "low"}}...{{end}} replacing <!-- IF NOT_LOW -->...<!-- END_IF --> markers, cleanTemplateOutput still handles backward compatibility

## Test Results
- **Tests Executed**: Yes
- **Passed**: 57
- **Failed**: 0
- **Coverage**: 80.9%

## Acceptance Criteria
- [x] 21 templates use {{.X}} format, no {{PLACEHOLDER}} residue
- [x] {{if .PhaseSummary}} blocks wrap two independent positions in 18 templates
- [x] {{if .CoverageStrategy}} handles three-state correctly
- [x] 4 doc templates have SURFACE_KEY header; 4 coding templates use {{if ne .Complexity "low"}} replacing <!-- IF NOT_LOW -->

## Notes
Go code placeholderReplacer bridge remains functional as safety net. Template files now use dot-notation natively, making bridge effectively a no-op. cleanTemplateOutput still handles edge cases during transition period.
