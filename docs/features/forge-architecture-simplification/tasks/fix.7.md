---
id: "fix.7"
title: "Update characterization tests after SourceTaskID sentinel elimination"
priority: "P0"
dependencies: []
status: 
type: "coding.fix"
breaking: true
---

# fix.7: Update characterization tests after SourceTaskID sentinel elimination

## Problem

Task 2.8 eliminated SourceTaskID sentinel and changed `countFixTasks` to active-only. This changed behavior of `--block-source` and fix-task dependency tracking. Characterization tests from task 0.1 that documented the old behavior now fail:

1. `TestAddCmd_WithTemplateAndVars` — source task 1.1 should have fix-task as dependency, got []
2. `TestAddCmd_BlockSource` — source 1.1 should be blocked, got ""
3. `TestAdd_BlockSource_CurrentBehavior_AllowsCompletedToBlocked` — source task not found in index
4. `TestBuildIndex_Orphan_WarningOnly` — likely index behavior change

## Scope
- `forge-cli/internal/cmd/add_cmd_test.go` — update TestAddCmd_WithTemplateAndVars, TestAddCmd_BlockSource
- `forge-cli/internal/cmd/characterization_test.go` — update TestAdd_BlockSource_CurrentBehavior
- `forge-cli/internal/cmd/characterization_test.go` — update TestBuildIndex_Orphan_WarningOnly

## Acceptance Criteria
- [ ] `go test ./internal/cmd/ -count=1` passes (0 failures)
