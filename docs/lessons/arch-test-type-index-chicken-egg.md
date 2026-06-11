---
created: "2026-05-16"
tags: [architecture, testing]
---

# BuildIndex Creates Orphan Per-Type Tasks Without Migrating Legacy Tasks

## Problem

After implementing `test-scripts-per-type` (per-type gen-scripts tasks like `T-quick-2-cli`), running a feature through quick mode produces an orphan `.md` file `quick-gen-scripts-go-test-cli.md` alongside the generic `quick-gen-scripts-go-test.md`. The per-type file is not tracked in `index.json` and never committed, creating confusion about which task actually generates test scripts.

## Root Cause

**Causal chain (4 levels):**

1. **Symptom**: Orphan `quick-gen-scripts-go-test-cli.md` appears in `tasks/` directory, untracked and uncommitted.

2. **Direct cause**: A subagent re-ran `forge task index` after T-quick-1 generated `test-cases.md`. The re-run detected CLI type and created the per-type `.md` file, but the subagent only committed its own execution results — not the re-indexed files.

3. **Root cause**: `BuildIndex` in `build.go` lacks migration logic. When `detectedTypes` becomes non-empty on re-run (because `test-cases.md` now exists), it generates new per-type task entries AND their `.md` files — but does NOT remove or supersede the legacy generic gen-scripts task (`T-quick-2` with key `quick-gen-scripts-go-test`). Both coexist in the index.

4. **Trigger condition**: `BuildIndex` calls `DetectTypesFromTestCases()` at runtime (build.go:270-275), which returns `nil` on first run (no `test-cases.md` yet) and `["cli"]` on re-run (after T-quick-1 generates it). This asymmetric detection creates the split behavior.

**Architectural issue**: Chicken-and-egg between index generation and test case generation. The index needs to know test types to create per-type tasks, but test types are only known after T-quick-1 runs — which requires the index to exist first.

## Solution

`BuildIndex` needs to handle the legacy-to-per-type transition explicitly:

1. When `detectedTypes` is non-empty, check if legacy gen-scripts tasks exist for the same profile
2. If they do and are still `pending` (not started), remove them and their `.md` files
3. Replace with per-type tasks in the index

Alternatively, defer per-type splitting to a separate step: `forge task index` always creates generic tasks initially, then a `forge task split-by-type` command runs after T-quick-1 to split the generic task into per-type tasks.

## Reusable Pattern

**When building a two-phase pipeline where phase 2 depends on phase 1's output for its own task structure**, the index/generation step must handle the transition from "no data" to "data available" explicitly. Options:

- **Option A**: Always generate the most specific structure upfront, with a fallback that degrades gracefully when data is absent (e.g., always create per-type tasks based on profile capabilities, not test case content)
- **Option B**: Generate a placeholder in phase 1, then replace it in phase 2 with the real structure, cleaning up the placeholder
- **Option C**: Make the index idempotent on re-run — detect and remove superseded tasks

The current implementation chose none of these, leaving the legacy and per-type tasks to coexist.

## Example

```go
// build.go:269-277 — current code
if len(profiles) > 0 && mode != "" {
    var detectedTypes []string
    testCasesPath := filepath.Join(...)
    if tcData, err := os.ReadFile(testCasesPath); err == nil {
        detectedTypes = DetectTypesFromTestCases(tcData)
    }
    testTasks := generateTestTasks(mode, profiles, detectedTypes)
    // Missing: if detectedTypes is non-empty, remove legacy generic gen-scripts tasks
}
```

## Related Files

- `forge-cli/pkg/task/build.go` — `BuildIndex()`, lines 269-318
- `forge-cli/pkg/task/testgen.go` — `GetQuickTestTasks()`, `DetectTypesFromTestCases()`
- `docs/proposals/test-scripts-per-type/proposal.md` — Source proposal

## References

- Proposal: `docs/proposals/test-scripts-per-type/proposal.md` — "Only types with test cases in test-cases.md get tasks"
- Success criterion from proposal: "quick-tasks creates separate tasks per detected test type instead of one T-quick-2"
