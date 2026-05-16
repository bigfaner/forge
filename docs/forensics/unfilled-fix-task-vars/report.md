---
created: "2026-05-16"
sessions: [47054d72-1a05-4201-842a-57506a648623]
skillsInvolved: [quick-tasks, run-tasks, execute-task]
severity: high
---

# fix-1.md Template Variables Not Filled — Test Pollution from Unit Tests

## Executive Summary

`fix-1.md` was created with unfilled template variables (`{{SOURCE_FILES}}`, `{{TEST_SCRIPT}}`, `{{TEST_RESULTS}}`) because `TestAddCmd_WithTemplateAndVars` leaked into the real project directory during a quality-gate `just test` run. The CLI's `ApplyVars` function silently leaves unfilled placeholders — no validation rejects incomplete templates. Fix: (1) add placeholder detection in `ApplyVars`, (2) ensure test isolation for `FindProjectRoot`.

## Investigation Scope

| Dimension | Value |
|-----------|-------|
| Sessions analyzed | 1 main + 8 subagents |
| Time range | 2026-05-16 03:24 → 04:58 (1.6h) |
| Skills involved | quick-tasks, run-tasks |
| Trigger | User noticed fix-1.md has `{{SOURCE_FILES}}` etc. unfilled |

## Timing Overview

| Session | Duration | Tool Time | Idle | Top Bottleneck |
|---------|----------|-----------|------|---------------|
| 47054d72 (main) | 1.6h | 4234.7s | ~1600s | Agent (1189.4s on T-quick-5) |

| Tool | Calls | Total | Avg | Max |
|------|-------|-------|-----|-----|
| Agent | 8x | 3685.4s | 460.7s | 1189.4s |
| Bash | 90x | 526.5s | 5.9s | 293.2s |
| Read | 17x | 4.2s | 248ms | 1.6s |

## Findings

### Finding 1: Test Pollution Creates Spurious fix-1.md

**Category:** `pipeline-gap`

**Affected sessions:** 47054d72 (quality-gate hook subprocess)

**Symptom:**
`docs/features/fix-task-claim-priority/tasks/fix-1.md` was created with unfilled template variables:
- `{{SOURCE_FILES}}` — not replaced
- `{{TEST_SCRIPT}}` — not replaced
- `{{TEST_RESULTS}}` — not replaced

The title "Fix: login bug" and `sourceTaskID: "1.1"` match EXACTLY the test data in `TestAddCmd_WithTemplateAndVars` (`add_cmd_test.go:14-88`):
```go
rootCmd.SetArgs([]string{
    "task", "add",
    "--title", "Fix: login bug",
    "--template", "fix-task",
    "--source-task-id", "1.1",
    "--description", "Selector not found",
})
```
No `--var` flags are passed in this test, so `{{SOURCE_FILES}}`, `{{TEST_SCRIPT}}`, `{{TEST_RESULTS}}` remain unfilled.

**Agent reasoning (from thinking block):**
> (No agent reasoning — this was a side effect of the quality-gate hook running `just test`, which spawned `go test` → `TestAddCmd_WithTemplateAndVars` → `executeAdd()` → `project.FindProjectRoot()` resolving to the real project instead of the test's temp dir.)

The agent DID discover the test pollution issue later (evidence line 462-490):
> "`FindRootInfoFrom` first checks `CLAUDE_PROJECT_DIR` env var, and if it's set, it uses that. But the `TestExecuteClaim` test doesn't clear `CLAUDE_PROJECT_DIR`. When we run `just test`, the `CLAUDE_PROJECT_DIR` is set to the real project root from our session, and the tests inherit it."

**Expected behavior (from code):**
`TestAddCmd_WithTemplateAndVars` should create `fix-1.md` in its temp directory, not in the real feature directory. `project.FindProjectRoot()` should resolve to the test's isolated temp dir.

**Gap:**
Two compounding gaps:
1. **Test isolation failure**: `FindProjectRoot()` uses `CLAUDE_PROJECT_DIR` env var and git context, both of which point to the real project when tests inherit the parent process environment
2. **No template validation**: `ApplyVars()` in `add.go:236-256` silently leaves unfilled `{{...}}` placeholders — no error, no warning

**Causal chain:**
1. **Symptom:** fix-1.md has unfilled `{{SOURCE_FILES}}`, `{{TEST_SCRIPT}}`, `{{TEST_RESULTS}}`
2. **Direct cause:** `TestAddCmd_WithTemplateAndVars` called `forge task add` without `--var` flags, and the command operated on the real project directory
3. **Root cause:** `project.FindProjectRoot()` resolved to the real project (via `CLAUDE_PROJECT_DIR` env or git context) instead of the test's temp dir, AND `ApplyVars()` silently accepts incomplete variable substitution

### Finding 2: No Validation for Unfilled Template Placeholders

**Category:** `pipeline-gap`

**Affected sessions:** All sessions using `forge task add --template`

**Symptom:**
`ApplyVars()` in `add.go:236-256` is a simple `strings.ReplaceAll` loop that silently skips missing variables:
```go
for key, val := range vars {
    result = strings.ReplaceAll(result, "{{"+key+"}}", val)
}
```
Any `{{PLACEHOLDER}}` without a matching key in `vars` remains in the output.

**Expected behavior:**
When using a template, all `{{...}}` placeholders should either be filled or an error should be raised before writing the file.

**Gap:**
No post-substitution validation checks for remaining `{{...}}` patterns.

**Causal chain:**
1. **Symptom:** fix-1.md written to disk with literal `{{SOURCE_FILES}}`
2. **Direct cause:** `CreateTaskMarkdown` calls `ApplyVars` and writes result without checking for unfilled placeholders
3. **Root cause:** `ApplyVars` has no error path for incomplete substitution — it was designed as best-effort, but callers assume it's complete

### Finding 3: fix-1 Marked Completed Without Actual Work

**Category:** `trust-without-verify`

**Affected sessions:** 47054d72 (main session, after quality-gate hook)

**Symptom:**
fix-1's record shows:
- Summary: "e2e test submit"
- Files Created: 无
- Files Modified: 无
- Time spent: "" (empty)
- Tests: 1 passed, 0 failed, 100% coverage (impossible with no changes)

**Expected behavior:**
A fix task should contain actual changes that fix a test failure. The record should reflect real work.

**Gap:**
The agent completed fix-1 without making changes because fix-1 was a spurious task created by test pollution. The agent didn't verify that the task was legitimate before marking it complete.

### Finding 4: Systemic Test Isolation Gap — 13 Files at Risk

**Category:** `pipeline-gap` (systemic)

**Affected files:** 13 test files in `forge-cli/internal/cmd/`

**Symptom:**
`setupFullProject` (in `integration_test.go:37-111`) has a `UseEnvVar` flag that controls isolation:
- `UseEnvVar: true` (SAFE): calls `t.Setenv("CLAUDE_PROJECT_DIR", dir)` — forces `FindProjectRoot` to temp dir
- `UseEnvVar: false` (DEFAULT, AT-RISK): creates `go.mod` in temp dir and chdirs — relies on filesystem marker walk, but `CLAUDE_PROJECT_DIR` env var takes priority and bypasses it

When `CLAUDE_PROJECT_DIR` is set in the host environment (e.g., Claude Code session, CI), every test using the default mode resolves `FindProjectRoot` to the real project.

**Full audit results:**

| File | Status | Isolation Mechanism |
|------|--------|-------------------|
| `add_cmd_test.go` | AT-RISK | `setupFullProject` default (no env var) |
| `claim_test.go` | AT-RISK | Manual `go.mod` + chdir (no env var) |
| `submit_test.go` | AT-RISK | Mixed: some `t.Setenv`, most default |
| `quality_gate_test.go` | AT-RISK | Mixed: some `t.Setenv`, most default |
| `integration_test.go` | AT-RISK | `setupFullProject` default |
| `claim_integration_test.go` | AT-RISK | `setupClaimTestProject` (no env var) |
| `migrate_test.go` | AT-RISK | `setupFullProject` default |
| `prompt_test.go` | AT-RISK | `setupFullProject` default |
| `runners_test.go` | AT-RISK | `setupClaimTestProject` (no env var) |
| `feature_test.go` | AT-RISK | Manual setup (no env var) |
| `proposal_test.go` | AT-RISK | Manual setup (no env var) |
| `status_test.go` | AT-RISK | Manual setup (no env var) |
| `lesson_test.go` | AT-RISK | Manual setup (no env var) |
| `cleanup_test.go` | SAFE | `UseEnvVar: true` |
| `verify_task_done_test.go` | SAFE | `t.Setenv("CLAUDE_PROJECT_DIR", dir)` |
| `config_test.go` | SAFE | Uses `--project-root` flag |
| `init_test.go` | SAFE | Uses `--project-root` flag |
| 10 other `*_test.go` | SAFE | Pure unit tests, no `FindProjectRoot` path |

**Leakage vector:**
```
just test (quality-gate hook)
  └─ go test ./...
       └─ TestAddCmd_WithTemplateAndVars
            └─ executeAdd()
                 └─ project.FindProjectRoot()
                      └─ CLAUDE_PROJECT_DIR set → resolves to REAL project
                           └─ forge task add operates on REAL index.json
                                └─ creates fix-1.md in REAL feature dir
```

**Root cause:**
`FindProjectRoot` resolution priority: `CLAUDE_PROJECT_DIR` env > `PROJECT_ROOT` env > filesystem marker walk. Tests using `go.mod` markers assume the walk will stop at temp dir, but the env var shortcut bypasses the walk entirely.

## Cross-Session Patterns

Not applicable — single session investigation.

## Recommendations

### P0 — Immediate Fixes (prevent recurrence)

| # | Action | Target File | Finding |
|---|--------|-------------|---------|
| 1 | Unify `setupFullProject` to always call `t.Setenv("CLAUDE_PROJECT_DIR", dir)` regardless of `UseEnvVar` flag | `forge-cli/internal/cmd/integration_test.go` (setupFullProject) | Finding 4 |
| 2 | Add `t.Setenv("CLAUDE_PROJECT_DIR", dir)` to `setupClaimTestProject` | `forge-cli/internal/cmd/runners_test.go` | Finding 4 |
| 3 | Add `t.Setenv("CLAUDE_PROJECT_DIR", dir)` to manual setup in `claim_test.go` (TestExecuteClaim, TestExecuteClaim_Continue, etc.) | `forge-cli/internal/cmd/claim_test.go` | Finding 4 |
| 4 | Add `t.Setenv("CLAUDE_PROJECT_DIR", dir)` to manual setup in `claim_integration_test.go` | `forge-cli/internal/cmd/claim_integration_test.go` | Finding 4 |
| 5 | Add placeholder validation: after `ApplyVars`, check for remaining `{{...}}` patterns and return error if any found | `forge-cli/pkg/task/add.go` (ApplyVars) | Finding 2 |

### P1 — Follow-up Fixes

| # | Action | Target File | Finding |
|---|--------|-------------|---------|
| 6 | Add `t.Setenv("CLAUDE_PROJECT_DIR", dir)` to manual setup in `submit_test.go` integration tests | `forge-cli/internal/cmd/submit_test.go` | Finding 4 |
| 7 | Add `t.Setenv("CLAUDE_PROJECT_DIR", dir)` to manual setup in `quality_gate_test.go` | `forge-cli/internal/cmd/quality_gate_test.go` | Finding 4 |
| 8 | Add `t.Setenv("CLAUDE_PROJECT_DIR", dir)` to `proposal_test.go`, `feature_test.go`, `status_test.go`, `lesson_test.go` | 4 files | Finding 4 |
| 9 | In `forge task submit`, validate task file doesn't contain unfilled template placeholders before recording completion | `forge-cli/internal/cmd/submit.go` | Finding 3 |

### P2 — Structural Improvements

| # | Action | Target File | Finding |
|---|--------|-------------|---------|
| 10 | Add linter/guard: `TestMain` or `init()` in `internal/cmd_test` that warns if `CLAUDE_PROJECT_DIR` is set during test runs | `forge-cli/internal/cmd/` | Finding 4 |

## Evidence

Evidence files at: `docs/forensics/unfilled-fix-task-vars/evidence/`

| File | Source | Size |
|------|--------|------|
| evidence.json | Main session (47054d72) | ~15 KB |
| agent-a007b57f76e6c73dd.json/evidence.json | Subagent (task 1) | ~8 KB |
| agent-a0db6b181a23f99cf.json/evidence.json | Subagent (fix-2) | ~5 KB |
| agent-a46db8c91faed68c9.json/evidence.json | Subagent (T-quick-5) | ~10 KB |
| + 5 other subagent evidence files | | ~3 KB each |

## Key Evidence Summary

1. `fix-1.md` created at 12:42, same minute as `fix-2.md` (quality-gate hook)
2. fix-1 title "Fix: login bug" + sourceTaskID "1.1" = EXACT match with `TestAddCmd_WithTemplateAndVars` test data
3. fix-2 title "fix unit-test: just test failure in quality gate" = matches `addFixTask()` pattern (properly filled)
4. No `forge task add` command found in any subagent or main session transcript → fix-1 was NOT created by an agent
5. Quality-gate hook runs `just test` → `go test` → `TestAddCmd_WithTemplateAndVars` → `executeAdd()` → creates fix-1 in real project dir
