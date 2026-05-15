---
id: "4"
title: "Update quick test task generation for per-type tasks"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 4: Update quick test task generation for per-type tasks

## Description

Modify `GetQuickTestTasks()` in `testgen.go` to create per-type gen-scripts tasks instead of a single T-quick-2 per profile. Same pattern as Task 3 but for quick mode: for each profile, iterate its capabilities and create a separate task per type (e.g., `T-quick-2-tui`, `T-quick-2-api`, `T-quick-2-cli`).

Update `resolveQuickDeps()` so that T-quick-3 depends on ALL per-type T-quick-2-* tasks for its profile.

## Reference Files
- `docs/proposals/test-scripts-per-type/proposal.md` — Source proposal
- `forge-cli/pkg/task/testgen.go` — `GetQuickTestTasks()`, `resolveQuickDeps()`
- `forge-cli/pkg/task/build.go` — `generateTestTasks()`, `BuildIndex()`

## Acceptance Criteria
- [ ] `GetQuickTestTasks()` creates separate tasks per type: e.g., `T-quick-2-tui`, `T-quick-2-api`, `T-quick-2-cli` for go-test profile
- [ ] Only types with test cases in test-cases.md get tasks (no empty tasks)
- [ ] `T-quick-3` depends on ALL per-type `T-quick-2-*` tasks for its profile
- [ ] Multi-profile works: `T-quick-2a-tui`, `T-quick-2a-api`, `T-quick-2b-tui`, `T-quick-2b-api`
- [ ] Task keys and file names include type suffix: `quick-gen-scripts-go-test-tui`, `quick-gen-scripts-go-test-api`
- [ ] `GenerateTestTaskMD()` produces correct `.md` content with profile + type info
- [ ] Existing quick test tasks (T-quick-1, T-quick-3, T-quick-4, T-quick-5) are unchanged

## Hard Rules
- MUST read test-cases.md to detect which types have cases — do NOT create tasks for types with zero cases
- Consistent with Task 3's approach for type detection and task naming

## Implementation Notes
- Share the type-detection logic with Task 3 (breakdown mode) — extract a common helper if both tasks modify the same file
- The `resolveQuickDeps()` function currently uses fixed index arithmetic (`i*4 + offset`). With per-type tasks, the number of gen-scripts tasks per profile is variable, so the arithmetic must be adapted
- Consider using a dynamic approach: track gen-scripts task IDs per profile, then set T-quick-3's dependencies to all of them
- Key risk: changing the fixed-index arithmetic in `resolveQuickDeps()` — test thoroughly with single and multi-profile scenarios
