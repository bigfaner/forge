---
date: "2026-05-13"
doc_dir: "docs/features/forge-cli-v3/design/"
iteration: 2
target_score: 900
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 893/1000** (target: 900)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Architecture Clarity      │  185     │  200     │ ⚠️         │
│    Layer placement explicit  │  70/70   │          │            │
│    Component diagram present │  65/70   │          │            │
│    Dependencies listed       │  50/60   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Interface & Model Defs    │  180     │  200     │ ⚠️         │
│    Interface signatures typed│  65/70   │          │            │
│    Models concrete           │  60/70   │          │            │
│    Directly implementable    │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Error Handling            │  135     │  150     │ ⚠️         │
│    Error types defined       │  45/50   │          │            │
│    Propagation strategy clear│  50/50   │          │            │
│    HTTP status codes mapped  │  40/50   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Testing Strategy          │  118     │  150     │ ⚠️         │
│    Per-layer test plan       │  40/50   │          │            │
│    Coverage target numeric   │  50/50   │          │            │
│    Test tooling named        │  28/50   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Breakdown-Readiness ★     │  185     │  200     │ ⚠️         │
│    Components enumerable     │  65/70   │          │            │
│    Tasks derivable           │  65/70   │          │            │
│    PRD AC coverage           │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Security Considerations   │  90      │  100     │ ⚠️         │
│    Threat model present      │  45/50   │          │            │
│    Mitigations concrete      │  45/50   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  893     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness >= 180/200 — can proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Section 1 (Architecture) | External dependency versions still not pinned — `github.com/spf13/cobra` and `gopkg.in/yaml.v3` listed without version constraints | -10 pts from Dependencies listed |
| Section 2 (Interfaces) | `submit.go` flags specified as comment `// Flags: --data <path>, --force` rather than typed Cobra `Flags()` code block; `add.go` and `index.go` flag specs are comment-only too | -5 pts from Interface signatures typed |
| Section 2 (Interfaces) | `addFixTask` signature has inline prose comment: "Signature mirrors current all_completed.go but adds cap check before the existing task.AddTask → task.CreateTaskMarkdown → feature.EnsureForgeState sequence" — not a pure typed signature | included in above |
| Section 2 (Models) | `RunOpts` struct fields lack validation constraints (e.g., `ProjectRoot` must be absolute path, `Feature` must match slug regex) | -10 pts from Models concrete |
| Section 2 (Directly implementable) | `countActiveFixTasks` implementation provided but `addFixTask` body is prose description, not code — developer must infer the cap-check insertion point | -5 pts from Directly implementable |
| Section 3 (Error) | `ErrBadProfile` and `ErrFeatureNotFound` described as "wrapped with value" but no wrapped-error construction code shown | -5 pts from Error types defined |
| Section 3 (Error) | HTTP status codes not applicable to CLI, exit codes mapped instead — partial credit (same as iter 1, deduction rule applies) | -10 pts from HTTP status codes mapped |
| Section 4 (Testing) | No test plan for Phase 4 reference updates (24 files) — the largest bulk-work phase has no per-layer test row | -10 pts from Per-layer test plan |
| Section 4 (Testing) | Test tooling section says `ExecRunner interface mock` but never names a specific mock library or pattern (e.g., `gomock`, `testify/mock`, or hand-rolled struct) | -12 pts from Test tooling named |
| Section 4 (Testing) | "Migration Equivalence Tests" section describes the process but does not specify the comparison framework or assertion library | -10 pts from Test tooling named |
| Section 5 (Breakdown) | PRD Goal "AI agent command selection >= 9/10" still has no design verification mechanism — how will this metric be validated post-implementation? | -5 pts from PRD AC coverage |
| Section 1 (Architecture) | Component diagram row for Top-level has alignment inconsistencies — `quality-gate` and `verify-task-done` names overflow into adjacent boxes | -5 pts from Component diagram |
| Section 6 (Security) | File-lock design is solid but Windows fallback mentions "goroutine+timer" without specifying the exact syscall or error path on Windows | -5 pts from Threat model present |
| Section 6 (Security) | Lock file cleanup strategy missing — if a process crashes without releasing, the lock auto-releases on fd close, but orphaned lock files are not addressed | -5 pts from Mitigations concrete |

---

## Attack Points

### Attack 1: [Testing Strategy — test tooling remains unnamed despite revision]

**Where**: "Unit tests inject a mock ExecRunner that returns predetermined outputs/errors. Integration tests use RealExec against test fixtures. No build tags required." (Section 6, pkg/e2e) and "Tool: go test + ExecRunner interface mock" (Testing Strategy table)
**Why it's weak**: The design describes the mock *pattern* (inject an interface) but never names the actual mock framework or assertion library. Is the mock a hand-rolled struct? `gomock`? `testify/mock`? The Migration Equivalence Tests say "Assert exit code and key output markers match" but do not name the assertion mechanism. A developer implementing the test plan must decide on tooling, which defeats the purpose of a complete design. The iteration 1 deduction for this same issue was 15 points; the revision added no new specificity.
**What must improve**: Name the exact mock approach (e.g., "hand-rolled struct implementing ExecRunner with canned responses" or "gomock-generated mock"). For Migration Equivalence Tests, name the comparison/assertion library or show the test scaffold code. If hand-rolling, show the mock struct definition.

### Attack 2: [Interface Definitions — flag specs are comments, not typed code]

**Where**: `// Flags: --data <path>, --force (identical to current record.go flags)` (submit.go), `// Flags: --fix-record-missed (identical to current prompt.go flag)` (prompt_get.go), `// Flags: --feature <slug> (optional, runs specific feature tests)` (e2e_run.go), `// Flags: --force (force reinstall)` (e2e_setup.go), `// Flags: --feature <slug> (required)` (e2e_verify.go)
**Why it's weak**: Cobra command definitions use Go struct syntax for Use/Short/Args/Run but fall back to comment prose for flags. In Cobra, flags are registered via `cmd.Flags().StringP(...)` calls. The design provides typed code for the struct but untyped comments for flags, creating an inconsistency where the developer can copy-paste the struct but must infer the flag registration code. For `add.go`, the Use field itself contains the flag names (`--title`, `--id`, `--priority`) embedded as a string rather than separated into `Flags()` calls, making it even less directly implementable.
**What must improve**: Add explicit `Flags()` registration code blocks for commands with flags, or at minimum provide a Flags table with columns: Flag Name | Shorthand | Type | Required | Default. The `add.go` Use field should show the actual Cobra Use string and list flags separately.

### Attack 3: [Testing Strategy — Phase 4 reference updates have no test plan]

**Where**: Phase 4 is described in the Phase 4 Reference Update Map (24 files, 200+ replacements) but the Testing Strategy per-layer table has no row for "Reference updates" or "Phase 4". The only mention is "just check-stale-refs CI target" in PRD Monitoring Requirements.
**Why it's weak**: Phase 4 is the largest bulk-work phase (24 files, hundreds of string replacements). The testing strategy table covers Cobra commands, pkg/e2e, renamed commands, and quality-gate cap — but omits Phase 4 entirely. The PRD explicitly lists `just check-stale-refs` as a monitoring requirement, but the design does not specify: (a) how to verify each file was updated correctly, (b) whether to test skill files are valid markdown after replacement, (c) how to detect partial replacements (e.g., only `task claim` replaced but `task record` missed in the same file). Without this, a task breakdown for Phase 4 would lack a verification step.
**What must improve**: Add a Phase 4 row to the per-layer test plan table. Specify: tool (e.g., `grep -r "task " plugins/ forge-cli/docs/` or a custom Go test), coverage approach (every file in the map has been updated), and verification command (`just check-stale-refs` passing).

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: PRD AC S3 concurrent write conflict unaddressed | ✅ Yes | Section 9 added with full file-locking design: `pkg/index/lock.go` with `LockFile`/`UnlockFile`/`ErrLockConflict`, integration into `submit.go` flow with 5s timeout, atomic index write via temp+rename, per-feature lock scope, error table, Windows fallback rationale |
| Attack 2: Missing interface signatures for renamed and unchanged commands | ✅ Yes | Section 2b added with full Cobra struct specs for forensic (3 subcommands), profile (3 subcommands), cleanup, feature, version, and all 6 unchanged task subcommands (claim, status, query, add, index, migrate). Each has Use/Short/Args/Run fields typed |
| Attack 3: Phase 4 reference updates lack per-file enumeration | ✅ Yes | "Phase 4 Reference Update Map" appendix added with 6 sub-tables: hooks.json (3 line-level refs), guide.md (3 refs), agent files (3 refs), command files (4 files), skill files (23 files enumerated, 12 with refs, 11 with "No task-command refs"), doc files (4 files with full ref lists). Summary: 24 files total |

---

## Verdict

- **Score**: 893/1000
- **Target**: 900/1000
- **Gap**: 7 points
- **Breakdown-Readiness**: 185/200 — **can proceed to /breakdown-tasks** (above 180 threshold)
- **Action**: Continue to iteration 3. Priority fixes: (1) Name specific mock framework and test tooling, (2) Type-flag flag registration as code or table instead of comments, (3) Add Phase 4 verification row to test plan table. These 3 fixes should close the 7-point gap.
