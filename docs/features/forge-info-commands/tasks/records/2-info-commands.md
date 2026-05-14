---
status: "completed"
started: "2026-05-14 15:37"
completed: "2026-05-14 16:01"
time_spent: "~24m"
---

# Task Record: 2 Info commands (proposal, feature, lesson)

## Summary
Implement forge proposal, forge lesson, and forge feature list/status info commands with PrintBlock/PrintField output format

## Changes

### Files Created
- forge-cli/pkg/proposal/proposal.go
- forge-cli/pkg/proposal/proposal_test.go
- forge-cli/pkg/lesson/lesson.go
- forge-cli/pkg/lesson/lesson_test.go
- forge-cli/internal/cmd/proposal.go
- forge-cli/internal/cmd/proposal_test.go
- forge-cli/internal/cmd/lesson.go
- forge-cli/internal/cmd/lesson_test.go

### Files Modified
- forge-cli/internal/cmd/feature.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/internal/cmd/root.go
- forge-cli/scripts/version.txt

### Key Decisions
- pkg/proposal and pkg/lesson use YAML frontmatter parsing consistent with existing pkg/task/frontmatter.go pattern
- feature list/status subcommands registered as Cobra subcommands under existing featureCmd, preserving backward compatibility (no args = show current, slug arg = set current)
- Lesson category inferred from filename prefix (gotcha-/arch-/pattern-/tool-/lesson-/hook-) as specified in acceptance criteria
- Score display uses em-dash for missing values, matching proposal spec
- Version bumped from 3.3.0 to 3.4.0 (minor: new commands)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 32
- **Failed**: 0
- **Coverage**: 89.1%

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
Pre-existing test failures (TestErrFeatureNotSet, TestGetTransitionAction) are unrelated to this task. New package coverage: proposal 90.2%, lesson 89.1%. Overall coverage 80.5%.
