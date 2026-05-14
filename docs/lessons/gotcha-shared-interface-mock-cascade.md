---
created: "2026-05-13"
tags: [architecture, interface, testing]
---

# Shared Interface Modification Causes Mock Cascade Stalls

## Problem

Three consecutive task-executor agents stalled (600s watchdog kill) when implementing features that required adding methods to the shared `MainItemRepo` interface (17 methods). The agents spent their entire budget fixing cascading test mock failures instead of implementing the feature.

**Stall pattern:**
1. Agent adds method to shared interface (e.g., `CalcCompletionByMap`)
2. 9 mock types across 6 test files immediately break
3. Each mock needs 17+ method stubs to satisfy the interface
4. Agent spends 10+ minutes fixing mocks, enters extended thinking, stalls

## Root Cause

Causal chain (3 levels):

1. **Symptom:** Agent stalled with no progress for 600s
2. **Direct cause:** Agent was fixing cascading mock failures (9 mock types x 17 methods = 153+ method stubs to maintain)
3. **Root cause:** `MainItemRepo` is a "fat interface" — 17 methods consumed by 9+ unrelated test files. Any method addition triggers O(consumers) mock updates. TDD amplifies this: the agent must fix ALL mocks before it can write its first test.
4. **Trigger condition:** Feature requires cross-cutting changes (adding methods to shared repo interfaces)

**Why TDD makes it worse but isn't the root cause:**
- TDD requires: write failing test → make it pass → refactor
- Writing a failing test requires creating a mock of ALL interface methods
- When the interface was just modified, the mock doesn't compile until ALL 17 methods are updated
- This is O(consumers) work before ANY feature logic can be tested
- The task splitting pattern (infra/CRUD/complex) solved scope overload but NOT interface cascade

## Solution

**Immediate fix:** Fix all mocks manually in the main session (as done for task 2.1c), then dispatch the feature agent on a clean codebase.

**Structural fix (future):** When a task requires adding methods to a shared interface:

1. **Pre-task: Interface update task** — Add the new method(s) to the interface AND all existing mocks in one focused task. No business logic, just interface + mock updates. This is mechanical, fast, and testable.
2. **Feature task** — Implement the actual feature logic using the already-updated interface. Mocks already compile, agent can focus on business logic.

## Reusable Pattern

**Rule:** If a task adds methods to a shared interface with >3 consumers, split into two tasks:

1. Interface + mock update only (breaking change reconciliation)
2. Feature implementation (using the updated interface)

**Detection heuristic:** Before dispatching, check if the task modifies:
- `main_item_repo.go` (MainItemRepo — 9 consumers)
- `sub_item_repo.go` (SubItemRepo — multiple consumers)
- Any interface with >5 methods

**Alternative pattern (better long-term):** Use Go interface embedding to create narrow interfaces. Instead of one fat `MainItemRepo`, embed `MainItemReader`, `MainItemWriter`, `MainItemMilestoneQueries` — consumers only mock what they use.

## Example

```go
// Before: fat interface, any change cascades to all 9 mock types
type MainItemRepo interface {
    Create(...) error
    FindByID(...) (*model.MainItem, error)
    // ... 17 methods total
    CalcCompletionByMap(...) (float64, error)  // NEW → breaks 9 mocks
}

// After (narrow interfaces): only milestone-aware consumers need updating
type MainItemReader interface {
    FindByID(...) (*model.MainItem, error)
    FindByBizKey(...) (*model.MainItem, error)
    List(...) (*dto.PageResult[model.MainItem], error)
}

type MainItemMilestoneQueries interface {
    CalcCompletionByMap(...) (float64, error)
    CountByMap(...) (int64, error)
    UnbindByMap(...) error
}

type MainItemRepo interface {
    MainItemReader
    MainItemWriter
    MainItemMilestoneQueries
}
```

## Related Files

- `backend/internal/repository/main_item_repo.go` — The fat interface (17 methods)
- `backend/internal/service/main_item_service_test.go` — 1 mock type
- `backend/internal/service/item_pool_service_test.go` — 1 mock type
- `backend/internal/service/view_service_test.go` — 1 mock type
- `backend/internal/service/milestone_service_test.go` — 1 mock type
- `backend/internal/service/milestone_map_service_test.go` — 1 mock type
- `backend/internal/handler/item_pool_handler_test.go` — 2 mock types
- `backend/internal/handler/router_test_stubs.go` — 1 mock type
- `docs/forensics/task-2.1-stall/report.md` — Forensic analysis of first stall

## References

- `docs/lessons/pattern-task-split-by-layer.md` — Layer splitting pattern (solved scope overload, not mock cascade)
- ISP (Interface Segregation Principle) from SOLID — the architectural principle being violated
