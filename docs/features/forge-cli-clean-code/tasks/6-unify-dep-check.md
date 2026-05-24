---
id: "6"
title: "Unify dependency check logic"
priority: "P1"
estimated_time: "1h"
dependencies: [4]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: Unify dependency check logic

## Description
Dependency checking logic (including `.x` wildcard pattern) is duplicated in 4 files. Consolidate into a single function in `pkg/task/`. Phase 3 (duplicate logic consolidation).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/pkg/task/` — Target location for unified function

## Acceptance Criteria
- [ ] Single dependency check function created in `pkg/task/`
- [ ] All 4 duplicate implementations replaced with calls to the unified function
- [ ] `.x` wildcard handling preserved in the unified implementation
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- The unified function must handle all existing patterns including `.x` wildcard
- Must be exported from `pkg/task/` for cross-package use

## Implementation Notes
- First grep for dependency check patterns across the codebase to identify all 4 sites
- Compare implementations to find the most complete version as the base
- The `.x` wildcard is a key semantic — ensure test coverage exists
