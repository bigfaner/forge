---
date: "2026-05-13"
doc_dir: "docs/features/forge-cli-v3/design/"
iteration: 3
target_score: 900
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 3

**Score: 966/1000** (target: 900)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Architecture Clarity      │  190     │  200     │ ⚠️         │
│    Layer placement explicit  │  70/70   │          │            │
│    Component diagram present │  70/70   │          │            │
│    Dependencies listed       │  50/60   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 2. Interface & Model Defs    │  193     │  200     │ ⚠️         │
│    Interface signatures typed│  68/70   │          │            │
│    Models concrete           │  70/70   │          │            │
│    Directly implementable    │  55/60   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3. Error Handling            │  140     │  150     │ ⚠️         │
│    Error types defined       │  50/50   │          │            │
│    Propagation strategy clear│  50/50   │          │            │
│    HTTP status codes mapped  │  40/50   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. Testing Strategy          │  148     │  150     │ ✅         │
│    Per-layer test plan       │  50/50   │          │            │
│    Coverage target numeric   │  50/50   │          │            │
│    Test tooling named        │  48/50   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 5. Breakdown-Readiness ★     │  195     │  200     │ ✅         │
│    Components enumerable     │  70/70   │          │            │
│    Tasks derivable           │  70/70   │          │            │
│    PRD AC coverage           │  55/60   │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 6. Security Considerations   │  100     │  100     │ N/A        │
│    Threat model present      │  N/A     │          │            │
│    Mitigations concrete      │  N/A     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │  966     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness >= 180/200 — can proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture/Dependencies | External dependency versions not pinned — `github.com/spf13/cobra` and `gopkg.in/yaml.v3` listed without version constraints (e.g., no `v1.8.0` or commit hash) | -10 pts from Dependencies listed |
| Section 7, line 667 | `addFixTask` body is prose description ("Signature mirrors current all_completed.go but adds cap check before the existing task.AddTask → task.CreateTaskMarkdown → feature.EnsureForgeState sequence") — not directly implementable code; developer must infer the cap-check insertion point and call ordering | -5 pts from Directly implementable |
| Section 6, lines 631-653 | `stubExec` mock struct is well-specified with typed code, but the `// Usage in tests:` block is commented-out pseudocode rather than a real `TestXxx(t *testing.T)` function — a developer must still write the test function scaffolding, the `runner = s` assignment pattern (package-level var mutation), and the restore/cleanup logic after each test | -2 pts from Test tooling named |
| Section 3, HTTP status codes | CLI tool — exit codes mapped instead of HTTP status codes. The rubric allocates 50 pts for "HTTP status codes mapped" with criterion "If API: are error types mapped to HTTP status codes?". This is a CLI refactoring, not an API, so full credit is not applicable. Exit code mapping serves the equivalent purpose and is thorough. | -10 pts from HTTP status codes mapped |
| PRD Coverage Map | PRD Goal "AI agent command selection >= 9/10" (line 33 of PRD) has no design verification mechanism — no test scenario set, no evaluation protocol, no "10 task scenarios" test harness described in the design. The PRD Coverage Map addresses all user stories (S1-S8) but the Goals table's primary success metric is untouched by any design component or test plan entry. | -5 pts from PRD AC coverage |
| Section 6, RunOpts struct | `RunOpts` struct fields lack validation constraints — `ProjectRoot` should be an absolute path, `Feature` should match a slug regex. No `Validate() error` method or constraint annotations shown. | -2 pts from Interface signatures typed |

---

## Attack Points

### Attack 1: Breakdown-Readiness — PRD primary success metric has no verification design

**Where**: PRD line 33: "AI agent command selection >= 9/10（旧 <= 7/10），通过 LLM 命令选择测试度量" and the design's PRD Coverage Map has no row addressing this goal.
**Why it's weak**: This is the PRD's top-line success metric — the entire motivation for the CLI rename is improving AI agent command selection accuracy. The design exhaustively maps every command structure, flag, and group, but never describes how this metric will be measured. No "10 task scenarios" test set, no evaluation protocol, no LLM test harness. The PRD Coverage Map rows address all user stories (S1-S8) but the Goals table metric is untouched. A task breakdown derived from this design would produce no verification task for the primary success metric.
**What must improve**: Add a row in the PRD Coverage Map or Testing Strategy specifying: (a) the 10 test scenarios (task descriptions) used to measure selection accuracy, (b) the evaluation method (present old vs new command list to an LLM, measure correct selection rate), (c) a test/verification task for running this evaluation post-implementation.

### Attack 2: Interface Definitions — `addFixTask` body is prose, not implementable code

**Where**: Section 7, line 665-668: "Signature mirrors current all_completed.go but adds cap check before the existing task.AddTask → task.CreateTaskMarkdown → feature.EnsureForgeState sequence."
**Why it's weak**: The document provides full typed Go code for `countActiveFixTasks` (lines 691-703), including the loop, filter criteria, and return value. But `addFixTask` — the function that actually creates fix tasks and enforces the cap — is described only with a prose comment explaining what it does "before the existing... sequence." A developer must: (1) read the current `all_completed.go` source to understand the existing sequence, (2) determine where to insert the cap check (before or after `task.AddTask`?), (3) write the actual Go code. This is the only behavioral change to quality-gate (the PRD's most complex new behavior), yet its implementation is left to inference rather than specification. The `countActiveFixTasks` code proves the document CAN show full implementations — `addFixTask` chose not to.
**What must improve**: Show the `addFixTask` function body with the cap check logic, the conditional `ErrMaxFixTasks` return path, and the delegation to `task.AddTask`. Even 10-15 lines of pseudocode with the branching logic would make this directly implementable.

### Attack 3: Architecture — external dependency versions unpinned

**Where**: Architecture/Dependencies section: `github.com/spf13/cobra` and `gopkg.in/yaml.v3` listed without version constraints.
**Why it's weak**: A `go.mod` change is a key artifact of this refactoring (module rename from `task-cli` to `forge-cli`). The design specifies `module forge-cli` but treats the dependency versions as "unchanged" without stating what they are. During the rename, `go.mod` will be regenerated — if the current codebase pins `cobra` at `v1.7.0` and the new module pulls `v1.8.4` (latest), behavioral regressions are possible (Cobra has had breaking changes in flag handling between minor versions). The "unchanged" label is also inconsistent with the File Structure Changes section which meticulously tracks every other file-level change — `go.mod` gets a new module path but its dependency block is handwaved.
**What must improve**: Pin exact versions for both external dependencies (e.g., `github.com/spf13/cobra v1.7.0`), or explicitly state that the `go.mod` dependency block is copied verbatim from the current module with only the module path changed.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Test tooling unnamed — no mock framework or assertion library specified | ✅ Yes | Section 6 now includes the full `stubExec` hand-rolled mock struct definition (lines 631-653) with `execResponse` type, `Run` method implementation, and usage example. The text explicitly states "hand-rolled mock matching project convention" and explains the project uses "only Go's standard `testing` package (`t.Fatalf`, `t.Errorf`) for assertions, with no `testify` or `gomock` dependency (verified in `go.mod`)." Migration Equivalence Tests now specify "using `t.Fatalf` and `strings.Contains` (standard library only — no assertion library)." |
| Attack 2: Flag specs as comments not typed code — Cobra Flags() registration missing | ✅ Yes | Flag specs now have explicit `init()` function code blocks with `cmd.Flags().StringP(...)`, `cmd.Flags().Bool(...)`, `cmd.Flags().String(...)`, `cmd.Flags().StringSlice(...)` calls (submit.go, prompt_get.go, add.go, index.go, e2e_run.go, e2e_setup.go, e2e_verify.go, forensicExtractCmd). Each flag-bearing command also has a Flag table with columns: Flag Name | Shorthand | Type | Required | Default. The `add.go` Use field no longer embeds flag names — flags are cleanly separated into `init()` with `.MarkFlagRequired("title")`. |
| Attack 3: Phase 4 reference updates absent from test plan | ✅ Yes | Per-layer test plan table now has a "Phase 4 ref updates" row (line 899): "Verification \| `go test` + `just check-stale-refs` \| Every file in Phase 4 map has all old refs replaced; no stale `task <cmd>` patterns remain; markdown files parse without errors \| 100% of mapped files". Key Test Scenarios #9 (line 911) specifies the exact `grep -rE` regex for detecting stale references and verifies `just check-stale-refs` CI target passes. |

---

## Verdict

- **Score**: 966/1000
- **Target**: 900/1000
- **Gap**: -66 points (surplus: target reached)
- **Breakdown-Readiness**: 195/200 — **can proceed to /breakdown-tasks** (above 180 threshold)
- **Action**: Target reached. All 3 previous attack points have been substantively addressed with concrete code, flag tables, and test plan rows. Remaining deductions are minor: dependency version pinning, one prose-only function body, and a PRD metric verification gap. None block progression to task breakdown.
