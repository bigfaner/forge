---
id: "3"
title: "Update breakdown test task generation for per-type tasks"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Update breakdown test task generation for per-type tasks

## Description

Modify `GetBreakdownTestTasks()` in `testgen.go` to create per-type gen-scripts tasks instead of a single T-test-2 per profile. For each profile, iterate its capabilities and create a separate task per type (e.g., `T-test-2-tui`, `T-test-2-api`, `T-test-2-cli`). Only create tasks for types that have matching test cases in test-cases.md.

Update `resolveBreakdownDeps()` so that T-test-3 depends on ALL per-type T-test-2-* tasks for its profile, not just one T-test-2.

## Reference Files
- `docs/proposals/test-scripts-per-type/proposal.md` — Source proposal
- `forge-cli/pkg/task/testgen.go` — `GetBreakdownTestTasks()`, `resolveBreakdownDeps()`
- `forge-cli/pkg/task/build.go` — `generateTestTasks()`, `BuildIndex()`

## Acceptance Criteria
- [ ] `GetBreakdownTestTasks()` creates separate tasks per type: e.g., `T-test-2-tui`, `T-test-2-api`, `T-test-2-cli` for go-test profile
- [ ] Only types with test cases in test-cases.md get tasks (no empty tasks for types without cases)
- [ ] `T-test-3` depends on ALL per-type `T-test-2-*` tasks for its profile
- [ ] Multi-profile works: `T-test-2a-tui`, `T-test-2a-api`, `T-test-2b-tui`, `T-test-2b-api`
- [ ] Task keys and file names include type suffix: `gen-test-scripts-go-test-tui`, `gen-test-scripts-go-test-api`
- [ ] `GenerateTestTaskMD()` produces correct `.md` content with profile + type info
- [ ] Existing breakdown test tasks (T-test-1, T-test-1b, T-test-3, T-test-4, T-test-4.5, T-test-5) are unchanged

## Hard Rules
- MUST read test-cases.md to detect which types have cases — do NOT create tasks for types with zero cases
- Profile capability detection determines which types are possible; test-cases.md determines which types are present

## Implementation Notes
- The function needs access to test-cases.md to detect present types. Currently `GetBreakdownTestTasks()` doesn't read test-cases.md — it may need a new parameter (e.g., `detectedTypes []string`) passed from `BuildIndex()` which can read the file
- Alternative: create tasks for all profile capabilities and let the task executor skip types with no cases (simpler but creates unnecessary tasks). The proposal prefers "only create when test-cases.md contains cases of that type"
- Key risk from proposal: tasks generated for types with no test cases — mitigated by reading test-cases.md first
- Dependency graph change: T-test-3 currently depends on one T-test-2 per profile; after change, depends on N T-test-2-* tasks per profile where N = number of types with cases
