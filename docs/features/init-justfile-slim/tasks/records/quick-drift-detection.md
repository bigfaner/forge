---
status: "completed"
started: "2026-06-09 23:17"
completed: "2026-06-09 23:25"
time_spent: "~8m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in 3 files: quality-gate.md (removed fallback recipe resolution descriptions), forge-cli-reference.md (added missing forge justfile scaffold command group), surface-rules.md (updated from rule file format to scaffold CLI). Regenerated vocabulary index with scaffold domain.

## Changes

### Files Created
无

### Files Modified
- docs/business-rules/quality-gate.md
- docs/conventions/forge-cli-reference.md
- docs/conventions/surface-rules.md
- docs/.vocabulary.md

### Key Decisions
无

## Document Metrics
3 drifted specs fixed (quality-gate fallback removal, CLI reference addition, surface-rules scaffold update); 0 orphaned rules; 1 new domain keyword (scaffold)

## Referenced Documents
- docs/proposals/init-justfile-slim/proposal.md
- forge-cli/internal/cmd/qualitygate/quality_gate_lifecycle.go
- forge-cli/internal/cmd/scaffold/register.go
- forge-cli/pkg/just/just.go

## Review Status
final

## Acceptance Criteria
- [x] All acceptance criteria met

## Notes
Drift detected via git diff narrowing: only specs whose domains overlap with changed files (scaffold/justfile/surface/quality-gate) were checked. Key finding: recipe fallback chain was removed in code (ResolvePrefixedRecipe returns empty string instead of generic recipe) but spec still documented fallback behavior.
