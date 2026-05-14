---
status: "completed"
started: "2026-05-14 16:21"
completed: "2026-05-14 16:27"
time_spent: "~6m"
---

# Task Record: 2 Info commands (proposal, feature, lesson)

## Summary
Verified and validated the three info commands (forge proposal, forge feature list/status, forge lesson) with full test coverage. Implementation was already in place: proposal.go/lesson.go cmd files, pkg/proposal and pkg/lesson packages, and feature.go extended with list/status subcommands. All 33 tests pass across the three packages with coverage above 80%.

## Changes

### Files Created
无

### Files Modified
无

### Key Decisions
- Implementation already complete -- no code changes needed, only verification
- pkg/proposal uses feature package constants for path construction (ProposalBaseDir, ProposalFileName)
- pkg/lesson infers category from filename prefix (gotcha-/arch-/pattern-/tool-/lesson-/hook-)
- feature list/status subcommands registered via init() in feature.go, preserving existing no-arg and slug-arg behavior

## Test Results
- **Tests Executed**: Yes
- **Passed**: 33
- **Failed**: 0
- **Coverage**: 80.6%

## Acceptance Criteria
- [x] forge proposal lists all proposals in table format: Slug | Created | Status | PRD | Feature
- [x] forge proposal <slug> shows detail: metadata, content summary, linked artifacts, file path
- [x] Created date reads from frontmatter created field, falls back to file birth time
- [x] PRD column checks docs/features/{slug}/prd/prd-spec.md existence
- [x] Feature column reads docs/features/{slug}/manifest.md status field
- [x] forge feature list lists all features: Slug | Status | Progress | PRD(score) | Design(score) | UI(score) | Tests(score)
- [x] Progress shows completed/total from tasks/index.json
- [x] Scores read from frontmatter score field; show em-dash when missing
- [x] forge feature status <slug> shows manifest summary, task counts by status, artifacts with scores
- [x] forge lesson lists all lessons: Name | Created | Tags | Category
- [x] Category inferred from file prefix (gotcha-/arch-/pattern-/tool-/lesson-/hook-)
- [x] forge lesson <name> shows metadata and file path (not full content)
- [x] Test coverage >= 80% for new and modified code

## Notes
All implementation was pre-existing. Task verified full acceptance criteria compliance and quality gate passage.
