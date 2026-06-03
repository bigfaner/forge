---
status: "completed"
started: "2026-06-03 20:34"
completed: "2026-06-03 20:43"
time_spent: "~9m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift: constants.md and enum-constants.md claimed defaultHealthPath was extracted to pkg/serverprobe/constants.go, but the constant did not exist and the code used inline "/health". Fixed by adding the constant and updating serverprobe.go to reference it. All other 20 spec files verified as current against the codebase.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/serverprobe/constants.go
- forge-cli/pkg/serverprobe/serverprobe.go

### Key Decisions
无

## Document Metrics
21 spec files checked, 1 drift found and fixed, 0 orphaned, 0 implicit new rules

## Referenced Documents
- docs/conventions/constants.md
- docs/conventions/enum-constants.md
- docs/business-rules/error-reporting.md
- docs/business-rules/quality-gate.md
- docs/business-rules/surface-orchestration.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/code-structure.md
- docs/conventions/dead-code.md
- docs/conventions/dispatcher-quality.md
- docs/conventions/error-handling.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/forge-distribution.md
- docs/conventions/naming.md
- docs/conventions/package-organization.md
- docs/conventions/prompt-template-hierarchy.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md
- docs/conventions/surface-cli.md
- docs/conventions/surface-rules.md
- docs/conventions/testing/index.md
- docs/conventions/testing/cli/core.md

## Review Status
final

## Acceptance Criteria
- [x] Changed files identified via git diff --name-only main...HEAD
- [x] Spec files with overlapping domains verified against current code
- [x] Drifted specs auto-fixed and committed with [auto-specs] tag, or no drift found

## Notes
Drift found in constants.md P4 and enum-constants.md TECH-const-001: defaultHealthPath constant was documented as extracted but did not exist. Code fix applied. All other rules validated as current.
