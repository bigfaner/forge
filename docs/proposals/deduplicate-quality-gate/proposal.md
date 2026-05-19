---
created: 2026-05-19
author: "fanhuifeng"
status: Draft
---

# Proposal: Tiered Test Execution Model

## Problem

`just test` (project-wide) runs redundantly across multiple layers per task. Four independent layers execute tests without coordination, compounding token and wall-clock cost. Meanwhile, agents lack instructions for targeted tests — the most cost-effective feedback loop. Two additional debris items persist: the `noTest` flag (100% redundant with `IsTestableType()`) and submit-task SKILL.md validation rules (fully duplicated in `submit.go`).

### Evidence

**Current state (before):**

| Task Type | `just test` Runs | Where |
|-----------|-------------------|--------|
| Normal (coding) | 2x | Template verification + CLI submit |
| BREAKING task | 3x | Template verification + Dispatcher gate + CLI submit |
| Gate task | 1x | Agent verification (necessary) |

**Proposed state (after):**

| Task Type | `just test` Runs | Where |
|-----------|-------------------|--------|
| Normal (coding) | 0x | Targeted tests only (framework-native, no `just test`) |
| Breaking task | 1x | CLI submit (full gate) |
| Gate task | 1x | Agent verification (unchanged) |
| All-completed | 1x | Quality gate hook (unchanged) |

### Root Cause

Test execution evolved organically. Each layer was added to fill a gap, but no layer was removed when a later one subsumed it:

| Layer | File | What it does | Origin |
|-------|------|-------------|--------|
| Task type template | `coding-*.md` | Agent runs `just test` before submitting | Pre-unification |
| submit-task skill | `SKILL.md` | Agent runs `just test` for metrics collection | Pre-unification |
| CLI quality gate | `submit.go` | `validateQualityGate()` runs compile→fmt→lint→test for all `coding.*` | 2026-05-03 unification |
| Dispatcher breaking gate | `execute-task.md` Step 3 | Dispatcher runs `just test` for BREAKING tasks | Pre-unification |

Additionally, no layer provides **targeted tests** (framework-native commands on changed packages) — the fastest feedback loop for individual tasks.

## Proposed Solution

Replace the ad-hoc multi-layer model with a **tiered test execution model** that balances quality with efficiency:

| Tier | When | What | Scope |
|------|------|------|-------|
| **1. Targeted tests** | Each coding task (during development) | Framework-native commands on changed code (e.g. `go test ./pkg/foo/...`) | Changed packages/files only |
| **2. Gate verification** | Gate task (stage gate) | All unit tests (`just test`) | Project-wide |
| **3. CLI submit gate** | Task submission (`forge task submit`) | compile→fmt→lint→test (breaking) or compile→fmt→lint (non-breaking) | Task scope |
| **4. Quality gate** | All-completed hook | `just test` + E2E regression | Project-wide |

The CLI submit gate reads the `breaking` flag from the task file's frontmatter (not from claim output). Rules:
- **breaking=true**: full gate (compile→fmt→lint→test)
- **Non-breaking coding tasks**: static gate (compile→fmt→lint)
- **Non-`coding.*` types** (doc, gate, summary, etc.): skip entirely

### Innovation Highlights

The tiered model mirrors the testing pyramid principle applied to CI orchestration: fast targeted feedback at the base, broad verification at the top. The key insight is that `just test` (project-wide) is expensive and should only run at verification checkpoints (gate task, CLI submit for breaking tasks, all-completed hook), not during individual task development.

## Requirements Analysis

### Key Scenarios

- **Normal coding task**: Agent develops feature → runs targeted tests on changed code → submits → CLI runs static gate (compile+fmt+lint) → passes
- **Breaking task (fix, cleanup)**: Agent develops fix → runs targeted tests → submits → CLI runs full gate (compile+fmt+lint+test) → passes
- **Gate task**: Agent verifies all previous phase work → runs `just test` as verification → submits → CLI skips gate (type is not `coding.*`)
- **All-completed**: Hook fires → runs compile+fmt+lint → `just test` project-wide → E2E regression
- **Non-coding task (doc, summary)**: Agent completes doc work → submits → CLI skips gate entirely

### Non-Functional Requirements

- No regression in quality gate pass/fail accuracy
- Reduced wall-clock time per task submission (no redundant `just test`)
- Reduced token consumption (dispatcher no longer parses BREAKING output, fewer steps to execute)

### Constraints & Dependencies

- `forge task claim` output format change is backward-compatible (removing a field)
- Test profile system (`just test` scope resolution) remains unchanged
- E2E test pipeline is a separate system and unaffected

## Alternatives & Industry Benchmarking

### Industry Solutions

Standard CI/CD pipelines use tiered testing: unit tests per commit, integration per stage, full regression before merge. This proposal applies the same tiering to AI agent orchestration.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No implementation cost | 2-3x redundant test runs per task; no targeted tests | Rejected: wasteful |
| CLI gate as sole authority (original proposal) | Prior design | Simple: remove `just test` everywhere except CLI submit | No targeted tests; all tasks pay full test cost at submit | Rejected: misses fast-feedback opportunity |
| **Tiered model** | Testing pyramid | Targeted feedback + tiered verification; breaking tasks get full gate; minimal redundancy | More nuanced; requires template-level changes | **Selected: balances quality and efficiency** |

## Feasibility Assessment

### Technical Feasibility

All changes are within existing Go CLI and skill template files. No new external dependencies. The `just.RunGate` function already supports configurable gate sequences — adding a static gate variant (compile→fmt→lint) is straightforward.

### Resource & Timeline

Small scope: ~10 files changed (3 Go, 5 skill templates, 1 SKILL.md, 1 testgen.go). Single developer, 1-2 sessions.

### Dependency Readiness

No external dependencies. All tools (`just`, `forge task claim`, `forge task submit`) already exist.

## Scope

### In Scope

1. **CLI submit gate** (`submit.go`): Add static gate path (compile→fmt→lint) for non-breaking coding tasks. Breaking tasks keep full gate. `breaking` flag read from task frontmatter (not claim output).
2. **Submit-task SKILL.md cleanup** (2 changes in one file):
   - Remove `just test [scope]` from metrics collection. Agent records metrics from its targeted test run (if available).
   - Remove CLI-enforced validation rules (quality gate pre-check, data validation table, "what submit does" description). Keep only agent-unique instructions: metrics collection, workflow steps, type reclassification, recovery.
3. **Task claim** (`claim.go`): Remove `BREAKING: true/false` from output.
4. **run-tasks.md / execute-task.md**: Remove Step 3 breaking gate entirely.
5. **Coding task type templates** (5 files): Replace `just test` with targeted test instructions (framework-native commands).
6. **Fix-task/cleanup-task templates** (2 files): Update test instructions to targeted tests.
7. **validateRecordData()** (`submit.go`): Accept zero metrics for non-breaking tasks (agent ran targeted tests, not `just test`).
8. **Remove `noTest` flag** (full removal): Delete `NoTest` from `Task`, `TaskState`, `FrontmatterData` structs. Remove from `testgen.go` (all auto-generated tasks use non-`coding.*` types, already handled by `IsTestableType()`). Remove from `submit.go`, `claim.go`, `formatTestsExecuted()`. Remove from SKILL.md references.

### Out of Scope

- CLI metrics parsing from test output (separate proposal)
- Gate task template changes (already correctly handled via type check)
- E2E test pipeline changes
- Profile system changes
- Changing quality gate sequence (compile→fmt→lint→test order)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent skips targeted tests, submits code with no verification | M | M | Gate task (Tier 2) catches it. Quality gate (Tier 4) is final safety net. |
| Non-breaking task introduces cross-package breakage missed by targeted tests | M | M | Gate task runs `just test` project-wide. Quality gate catches at session end. |
| Records show zero metrics for non-breaking tasks | H | L | Targeted tests produce framework-native output; agent can capture pass/fail counts manually. Zero metrics acceptable when no relevant tests exist. |
| Removing BREAKING output breaks external tooling | L | L | `forge task claim` output is consumed only by run-tasks/execute-task dispatchers within forge. No external consumers. |
| Removing `noTest` breaks edge-case task types | L | L | All auto-generated tasks using `noTest` have non-`coding.*` types. `IsTestableType()` already handles them. Verified by audit. |

## Success Criteria

- [ ] Normal coding tasks: `just test` runs 0x per task (targeted tests only)
- [ ] Breaking tasks: `just test` runs 1x at CLI submit (full gate)
- [ ] Gate task: `just test` runs 1x at verification (unchanged)
- [ ] All-completed: `just test` runs 1x at quality gate (unchanged)
- [ ] Non-breaking coding tasks: CLI submit runs static gate (compile→fmt→lint, no test)
- [ ] Agent runs targeted tests during task development (framework-native commands)
- [ ] `forge task claim` no longer outputs BREAKING field; CLI submit reads `breaking` from frontmatter
- [ ] run-tasks/execute-task have no Step 3 breaking gate
- [ ] No regression in quality gate pass/fail accuracy
- [ ] `noTest` flag fully removed from all structs, templates, and logic; `IsTestableType()` is sole authority
- [ ] submit-task SKILL.md contains only agent-unique instructions (no CLI-enforced validation rules)
- [ ] validateRecordData() accepts zero metrics for non-breaking tasks

## Next Steps

- Proceed to `/write-prd` to formalize requirements
