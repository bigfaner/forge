---
status: "completed"
started: "2026-06-06 14:39"
completed: "2026-06-06 14:46"
time_spent: "~7m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed 3 spec drifts in project-level spec files caused by surface-key-test-alignment feature's switch from per-surface-type to per-surface-key expansion

## Changes

### Files Created
无

### Files Modified
- docs/conventions/testing/index.md
- docs/conventions/testing/cli/core.md
- docs/ARCHITECTURE.md

### Key Decisions
无

## Document Metrics
3 drifted specs detected, 0 orphaned, 0 implicit new rules; all 3 fixed and committed

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md
- docs/business-rules/surface-orchestration.md
- docs/conventions/surface-rules.md
- docs/conventions/surface-cli.md
- plugins/forge/skills/test-guide/SKILL.md
- plugins/forge/skills/gen-test-scripts/SKILL.md
- forge-cli/pkg/task/pipeline.go

## Review Status
final

## Acceptance Criteria
- [x] Drift detection completed for specs overlapping with changed files
- [x] Drifted specs auto-fixed and committed (if any drift found)

## Notes
Drift-only mode (no PRD/design docs for quick-mode feature). Used git diff to narrow scope to 5 candidate spec files, found drift in 3: testing/index.md (file location column), testing/cli/core.md (directory description + build tag reference), ARCHITECTURE.md (gen-test-scripts output path). quality-gate.md per-surface-type reference confirmed correct (Phase 3 orchestrates by type, not key). Commit 5a17d485 with [auto-specs] tag.
