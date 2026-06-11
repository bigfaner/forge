---
status: "completed"
started: "2026-05-29 16:21"
completed: "2026-05-29 16:39"
time_spent: "~18m"
---

# Task Record: 4 build.go — intent-driven pipeline routing

## Summary
Implemented intent-driven pipeline routing in build.go: added Intent field to BuildIndexOpts, updated needsTestPipeline() to short-circuit for refactor/cleanup intents, updated detectMode() to force Quick mode for cleanup intent, and wired CLI handlers (index.go, add.go) to read Proposal.Intent via proposal.FindBySlug(). IsTestableType() left unchanged per Hard Rules.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/extract_test.go
- forge-cli/internal/cmd/task/index.go
- forge-cli/internal/cmd/task/add.go

### Key Decisions
- Intent defaults to 'new-feature' inside BuildIndex() when opts.Intent is empty, ensuring backward compatibility without requiring CLI callers to set a default
- needsTestPipeline() short-circuits before iterating tasks when intent is refactor/cleanup, keeping IsTestableType() untouched per Hard Rules
- detectMode() checks cleanup intent before document existence, forcing Quick mode regardless of PRD presence
- CLI handlers silently ignore proposal.FindBySlug() errors (e.g. no proposal file), letting empty intent flow to BuildIndex() which defaults to new-feature

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 87.8%

## Acceptance Criteria
- [x] BuildIndexOpts has Intent string field, CLI handler reads Proposal.Intent, empty defaults to new-feature
- [x] needsTestPipeline() signature includes intent param, refactor/cleanup returns false immediately
- [x] detectMode() forces Quick for cleanup intent, ignores PRD existence
- [x] IsTestableType() not modified, pure type judgment preserved
- [x] Backward compatible: missing intent field defaults to new-feature pipeline, behavior unchanged

## Notes
All existing tests pass with updated signatures. 6 new test cases added covering intent short-circuit (refactor, cleanup, empty), detectMode cleanup override, and backward compatibility (empty vs explicit new-feature produce identical output).
