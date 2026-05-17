---
id: "5"
title: "Update testgen.go to use Language instead of ProfileName"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1", "4"]
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 5: Update testgen.go to use Language instead of ProfileName

## Description
Update `GetBreakdownTestTasks` to accept `languages []Language` and `interfaces []string` instead of `profiles []string` and `capabilities []string`. Replace per-profile-per-type expansion with per-language-per-interface expansion. Update `TestTaskDef` struct: `ProfileName string` → `Language Language`. This eliminates the combinatorial profile × capability matrix.

## Reference Files
- `docs/proposals/simplify-testing-model/proposal.md` — Source proposal (Feasibility section, item 5)
- `forge-cli/pkg/task/testgen.go` — Current task generation logic
- `forge-cli/pkg/task/testgen_test.go` — Existing test generation tests

## Acceptance Criteria
- `GetBreakdownTestTasks` signature uses `[]Language` and `[]string` (interfaces) parameters
- Per-language-per-interface task expansion produces correct test tasks
- `TestTaskDef.ProfileName` field replaced with `TestTaskDef.Language` of type `Language`
- Task IDs and slugs use language keys (not profile names)
- Existing test generation test cases updated and passing
- `go test ./...` passes

## Hard Rules
- Task generation output must be compatible with `forge task index` command
- Do not change the task file template structure — only the generation logic

## Implementation Notes
- The simplification is from 2D (profile × capability) to 1D (language → interfaces subset). Each language has a fixed set of supported interfaces from `languageCapabilities`, and the `interfaces` config field narrows which ones to generate tests for.
- This is the core complexity reduction from the proposal: eliminating the combinatorial expansion that was hard to reason about
