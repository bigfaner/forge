---
title: "Dispatcher Quality Gate Conventions"
domains: [dispatcher, compilation, quality-gate, fix-task, diagnostics, run-tasks]
---

# Dispatcher Quality Gate Conventions

Rules for the run-tasks dispatcher to maintain codebase integrity across sequential task execution.

### TECH-dispatcher-quality-001: Monitor Compilation Diagnostics After Task Completion

**Requirement**: After each task completes (Step 2b verify), the dispatcher MUST check for compilation errors in the IDE diagnostics or run the appropriate compile recipe (for projects with a justfile). If compilation errors exist, the dispatcher MUST create a fix task targeting those errors before claiming the next feature task. Fix type is derived from source task category (see derivation table below).

**Compile recipe selection**: When the completed task has a `surface-key`, the dispatcher MUST use the prefixed recipe `just <key>-compile` (e.g., `just backend-compile`). If the prefixed recipe does not exist in the justfile, fall back to `just compile` (generic). When the task has no `surface-key`, use `just compile` directly.

**Scope**: [CROSS]

**Source**: /learn entry 2026-05-21

**Why**: Task executors may report "tests passed" based on intermediate-state test runs, then make additional edits that break compilation before submitting. The `forge task status` returning "completed" only proves the executor called submit, not that the code compiles. The dispatcher, as the orchestrator, has visibility into cross-task state and is the only component positioned to independently verify.

**Priority of diagnostic signals**:
1. **Compilation errors** — must block pipeline, spawn fix task immediately
2. **Test failures** — should spawn fix task but non-blocking
3. **Lint warnings** — informational, no action needed
4. **Style suggestions** — ignore

**Implementation**:
- After Step 2b (`forge task status <ID>` returns completed), collect diagnostics
- Filter to compilation errors only (undefined symbols, redeclared names, wrong arg counts)
- If non-empty: `forge task add --type <derived-type> --title "Fix compilation errors from task X.Y"`
- Fix type derivation: extract `TASK_CATEGORY` from claim output, then map: `doc`/`eval` → `doc.fix`, `coding`/`test`/`validation`/`gate` → `coding.fix`
- For fmt/lint failures (non-breaking): use `coding.cleanup` task type instead
- Fix task gets auto-claimed on next loop iteration (priority over feature tasks)

### TECH-dispatcher-quality-002: Quality Gate Must Run on Final Code State

**Requirement**: When a task has `breaking: true`, the `forge task submit` quality gate MUST execute `just compile` and `just unit-test` against the **final state of all modified files**, not an intermediate snapshot. Test results from earlier in the execution session are stale and must not be used as the quality gate verdict.

**Scope**: [CROSS]

**Source**: /learn entry 2026-05-21

**Why**: Task 2.4 (`breaking: true`) recorded "49 passed, 0 failed" and "go build + go test ./internal/cmd/... pass" in its execution record, yet the same package contained undefined symbols (`submitForce`). The executor ran tests at an intermediate point, then made further edits that broke compilation, and submitted with the stale test results.

**Implementation**: The submit command should run build+test as the final step of the quality gate, after all file modifications are complete, not accept pre-computed test results from the executor's memory.
