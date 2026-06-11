---
id: "2"
title: "Replace test-type detection with config-driven capabilities and remove legacy code"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "enhancement"
mainSession: false
---

# 2: Replace test-type detection with config-driven capabilities and remove legacy code

## Description

Replace the `DetectTypesFromTestCases()` runtime detection in `BuildIndex` with config-driven capabilities from `.forge/config.yaml`. This is a single coherent change: all parts are interdependent (can't wire capabilities without changing signatures, can't remove legacy branches without the new path). Removes all legacy generic code paths, making `BuildIndex` fully deterministic and config-driven.

## Reference Files
- `docs/proposals/task-index-dedup-legacy-types/proposal.md` — Source proposal
- `forge-cli/pkg/task/build.go` — `BuildIndex`, `BuildIndexOpts`, `generateTestTasks`
- `forge-cli/pkg/task/testgen.go` — `GetQuickTestTasks`, `GetBreakdownTestTasks`, `resolveQuickDeps`, `resolveBreakdownDeps`, `DetectTypesFromTestCases`, `summaryTableRow`
- `forge-cli/internal/cmd/index.go` — Caller of `BuildIndex`
- `forge-cli/internal/cmd/add.go` — Caller of `BuildIndex`
- `forge-cli/pkg/task/build_test.go` — Existing BuildIndex tests
- `forge-cli/pkg/task/testgen_test.go` — Existing testgen tests (if exists)
- `forge-cli/pkg/profile/config.go` — `ReadConfig()`, `ForgeConfig`
- `forge-cli/pkg/profile/embed.go` — `UnionCapabilities()`, `ValidateCapabilities()` (from Task 1)

## Acceptance Criteria

- [ ] `BuildIndexOpts` has `TestCapabilities []string` field
- [ ] `BuildIndex` reads `opts.TestCapabilities` instead of calling `DetectTypesFromTestCases()`
- [ ] Callers (`index.go`, `add.go`) resolve capabilities: `config.yaml` → fallback to `UnionCapabilities(profiles)` → pass to `BuildIndexOpts`
- [ ] `DetectTypesFromTestCases()` function and `summaryTableRow` regex deleted entirely
- [ ] `detectedTypes` parameter removed from `GetQuickTestTasks`, `GetBreakdownTestTasks`, `resolveQuickDeps`, `resolveBreakdownDeps`
- [ ] All legacy `else` branches (generic task creation) removed from all 4 functions
- [ ] `forge task index` produces identical output regardless of whether `test-cases.md` exists
- [ ] All unit tests updated and passing: `go test -race -cover ./forge-cli/...`
- [ ] No orphaned generic tasks in any scenario

## Hard Rules

- Follow TDD: update tests first to reflect new signatures and behavior (RED), then implement changes (GREEN)
- `breaking: true` — this modifies `BuildIndexOpts` struct and multiple public function signatures
- After all code changes, bump version in `scripts/version.txt` (minor bump — new behavior)

## Implementation Notes

### Capability resolution in callers

```
// In index.go and add.go:
cfg, _ := profile.ReadConfig(projectRoot)
caps := cfg.Capabilities
if len(caps) == 0 {
    caps, _ = profile.UnionCapabilities(profiles)
}
if err := profile.ValidateCapabilities(caps); err != nil {
    // warn but don't abort — invalid caps are filtered, not fatal
}
opts.TestCapabilities = caps
```

### Function signature changes

```go
// Before:
func GetQuickTestTasks(profiles []string, detectedTypes []string) []TestTaskDef
func GetBreakdownTestTasks(profiles []string, detectedTypes []string) []TestTaskDef
func resolveQuickDeps(tasks []TestTaskDef, profiles []string, _ bool, detectedTypes []string)
func resolveBreakdownDeps(tasks []TestTaskDef, profiles []string, _ bool, detectedTypes []string)

// After:
func GetQuickTestTasks(profiles []string, capabilities []string) []TestTaskDef
func GetBreakdownTestTasks(profiles []string, capabilities []string) []TestTaskDef
func resolveQuickDeps(tasks []TestTaskDef, profiles []string, _ bool, capabilities []string)
func resolveBreakdownDeps(tasks []TestTaskDef, profiles []string, _ bool, capabilities []string)
```

### What to delete

- `DetectTypesFromTestCases()` function (testgen.go:566-586)
- `summaryTableRow` regex (testgen.go:560-564)
- `build.go:271-276` — the test-cases.md reading block
- All `else` blocks in the 4 functions (generic task creation branches)

### Test updates

- `build_test.go`: all `BuildIndexOpts` constructions need `TestCapabilities` field where profiles are set
- Tests that relied on two-pass behavior (no test-cases.md → generic, then test-cases.md → per-type) must be rewritten to expect deterministic per-type output
- New test: `BuildIndex` with empty `TestCapabilities` produces no per-type test tasks
- New test: `BuildIndex` with `[cli]` produces T-quick-2-cli (not generic T-quick-2)
