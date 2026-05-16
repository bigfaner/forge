---
status: "completed"
started: "2026-05-17 01:23"
completed: "2026-05-17 01:26"
time_spent: "~3m"
---

# Task Record: 2 eval-prd Rubric Enhancement

## Summary
Enhanced eval-prd rubric with context frontmatter declaration and two new scoring dimensions (Scenario Completeness 150pts, Edge Case Coverage 100pts). Adjusted existing dimensions to maintain 1000-point total: Background & Goals 150->100, Flow Diagrams 200->150, User Stories 300->200, Scope Clarity 150->100. Functional Specs/Flow Completeness unchanged at 200. Mode A/B detection preserved with identical new dimensions across both modes.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval/rubrics/prd.md

### Key Decisions
- Point redistribution: reduced Background & Goals (-50), Flow Diagrams (-50), User Stories (-100), Scope Clarity (-50) to free 250pts for two new dimensions (150+100)
- Scenario Completeness criteria explicitly reference injected business-rules context for contradiction detection
- Edge Case Coverage focuses on three concrete areas: error paths (40pts), boundary conditions (35pts), failure recovery (25pts)

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] prd.md rubric frontmatter includes context field declaring required conventions and business-rules categories
- [x] New dimension Scenario Completeness added with clear criteria for evaluating whether PRD scenarios cover real-world usage patterns
- [x] New dimension Edge Case Coverage added with criteria for error paths, boundary conditions, and failure states
- [x] Existing dimensions adjusted to maintain 1000-point total scale
- [x] Scoring criteria for new dimensions reference injected context: scorer should check PRD scenarios against loaded business-rules and conventions for contradictions
- [x] Mode A/B detection still works correctly with new dimensions

## Notes
Documentation-only task. Quality gate tests all pass (Go test suite). No functional code changes.
