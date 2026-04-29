---
date: "2026-04-30"
doc_dir: "docs/features/justfile-standard-vocabulary/design/"
iteration: "3"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 3

**Score: 92/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Architecture Clarity      │  20      │  20      │ ✅         │
│    Layer placement explicit  │  7/7     │          │            │
│    Component diagram present │  7/7     │          │            │
│    Dependencies listed       │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Interface & Model Defs    │  18      │  20      │ ✅         │
│    Interface signatures typed│  7/7     │          │            │
│    Models concrete           │  7/7     │          │            │
│    Directly implementable    │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Error Handling            │  14      │  15      │ ✅         │
│    Error types defined       │  5/5     │          │            │
│    Propagation strategy clear│  5/5     │          │            │
│    HTTP status codes mapped  │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Testing Strategy          │  13      │  15      │ ✅         │
│    Per-layer test plan       │  5/5     │          │            │
│    Coverage target numeric   │  3/5     │          │            │
│    Test tooling named        │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Breakdown-Readiness ★     │  18      │  20      │ ✅         │
│    Components enumerable     │  7/7     │          │            │
│    Tasks derivable           │  6/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Security Considerations   │  9       │  10      │ ✅         │
│    Threat model present      │  4/5     │          │            │
│    Mitigations concrete      │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  92      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 18/20 — PASSED (above 12/20 gate)

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Model 5 validation checklist (lines 42-44) | Task templates referenced as `<template-path-1>` through `<template-path-4>` — these are placeholder paths that a developer cannot resolve into actual file locations. This makes the validation grep commands non-executable as written. | -2 pts (Directly Implementable) |
| Error Types table (lines 399-403) vs Interface 4 (lines 219-238) | The Error Types table lists `exit != 0` generically for `just project-type`, but Interface 4's Scope Resolution Protocol distinguishes exit 127 (command not found) as a specific failure mode. These two sections are inconsistent — the error table should enumerate the same exit code categories the protocol handles. | -1 pt (Exit Codes Mapped) |
| Testing Strategy (lines 426-443) | "100% PRD 覆盖" is claimed (line 443) without a test-to-AC traceability matrix. The PRD has 12+ acceptance criteria across 5 stories; the 5 "Key Test Scenarios" (lines 435-439) do not map 1:1 to ACs. For example, Story 3 has 4 ACs but no single test scenario covers "非混合项目所有任务 scope=all". Story 4's sequential execution AC ("agent 连续执行 `just install`、`just compile`、`just test`") has no corresponding test scenario. | -2 pts (Coverage Target) |
| Breakdown-Readiness — validation checklist (lines 33-44) | The 10 validation items specify `grep -c` commands as ad-hoc verification, but these are not formalized as test tasks with pass/fail criteria or integrated into the testing strategy. They exist in a separate section from the Testing Strategy, creating a gap: the Testing Strategy's "skill migration: 14/14 文件" target does not reference the validation checklist's grep commands as the means of verification. | -1 pt (Tasks Derivable) |
| PRD AC Coverage (lines 464-482) | Story 4 AC: "agent 连续执行 `just install`、`just compile`、`just test`，全部成功，每步退出码均为 0，agent 无需人工介入" — sequential command chaining with partial failure handling is not addressed. The PRD Coverage Map row for Story 4 AC maps to "Error Handling" but the error handling section only covers individual recipe errors, not sequential execution semantics. This AC was flagged in iteration 2 and remains unresolved. | -1 pt (PRD AC Coverage) |
| Security (lines 449-452) | Threat model contains only 2 threats. For a feature that modifies 14+ files, introduces a scope dispatch mechanism with bash case evaluation, and generates justfiles that execute arbitrary toolchain commands, additional threats warrant consideration: (a) malicious project-type output from a compromised justfile, (b) template injection in init-justfile if project signals contain unexpected characters. | -1 pt (Threat Model) |

---

## Attack Points

### Attack 1: Testing Strategy — no test-to-AC traceability matrix despite claiming "100% PRD coverage"

**Where**: Lines 435-443: "Key Test Scenarios" lists 5 scenarios; line 443 states "e2e 覆盖所有 15 个命令 + 3 种项目类型 + scope 解析流程 = 100% PRD 覆盖"
**Why it's weak**: The design claims 100% PRD coverage but provides no mapping from test scenarios to acceptance criteria. The PRD user stories contain at least 12 distinct ACs:
- Story 1: 3 ACs (standard verb, pure-backend `just test`, mixed `just build frontend`)
- Story 2: 3 ACs (package.json → frontend, go.mod → backend, both → mixed)
- Story 3: 4 ACs (index.json scope field, frontend-path → scope=frontend, cross-scope → scope=all, non-mixed → all)
- Story 4: 4 ACs (exit 0 on success, exit != 0 on failure, compile error to stderr, sequential execution)
- Story 5: 2 ACs (scope mismatch warning, mixed scope success)

The 5 test scenarios cannot cover 12+ ACs without explicit mapping. "100% PRD 覆盖" is an assertion, not evidence. This was flagged in iteration 1 AND iteration 2 — it persists as the longest-unresolved deduction.
**What must improve**: Add a test-to-AC traceability table. For each AC, specify which test scenario (or new scenario) validates it. If "100% PRD coverage" means only that every PRD requirement has at least one test, prove it with a mapping.

### Attack 2: Breakdown-Readiness — validation checklist is disconnected from testing strategy

**Where**: Lines 33-44 (Validation Checklist with grep commands) vs lines 426-443 (Testing Strategy section)
**Why it's weak**: The design contains two separate verification mechanisms that do not reference each other. The Validation Checklist proposes 10 `grep -c` commands to verify skill files contain the correct commands. The Testing Strategy proposes "skill 迁移 | e2e 测试 | skill 文件内容检查（无原始命令）| 14/14 文件". The grep commands in the checklist are the obvious implementation of the e2e test, but this connection is never made explicit. A developer implementing tests from this design must decide: do I write the grep-based checks as standalone scripts? As part of the e2e suite? As CI checks? The design does not say. Furthermore, the scope annotation prompt engineering (item 5) has no validation step — how do you verify that the breakdown-tasks scope assignment prompt produces correct scope values?
**What must improve**: Merge the validation checklist into the testing strategy. Each grep command should become a named test case in the e2e suite. Add a test scenario for "scope assignment correctness" that validates the breakdown-tasks prompt output.

### Attack 3: PRD Coverage — sequential command execution AC remains unaddressed for third iteration

**Where**: PRD Story 4 AC: "Given agent 连续执行 `just install`、`just compile`、`just test`，When 全部成功时，Then 每步退出码均为 0，agent 无需人工介入". PRD Coverage Map (line 476) maps this to "recipe error handling" in Error Handling section (lines 394-413).
**Why it's weak**: The PRD Coverage Map claims this AC is addressed by "Error Handling", but the Error Handling section only defines individual recipe error behavior. No design element addresses sequential execution: what is the contract between `just install` and `just compile`? If `just install` succeeds but `just compile` fails, what state is the project in? Does the agent retry from the beginning or from the failed step? This AC requires the design to specify that (a) each recipe is independently executable, (b) the agent can chain them sequentially, and (c) partial failure leaves the project in a known state. The design implicitly assumes (a) and (b) but never states them, and completely omits (c). This was flagged in iteration 2 and has not been addressed.
**What must improve**: Add a "Sequential Execution Contract" subsection stating: (1) each recipe is independently callable and idempotent where applicable; (2) the standard agent execution sequence is `install → compile → test`; (3) on partial failure, the agent may re-invoke from the failed step. Or explicitly mark this AC as "design-scope: N/A — sequential execution is an agent-level concern, not a justfile-level concern" and justify why it does not belong in the design.

---

## Previous Issues Check

| Previous Attack (Iteration 2) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: recipe template implementation tasks are still implicit | ✅ Yes | New "Task Decomposition: Recipe Templates" section (lines 381-392) explicitly lists 4 tasks (A, B, C, D) with content, recipe count, key features, and dependency ordering. |
| Attack 2: migration checklist has no task mapping — implementation and validation mixed | ✅ Yes | Migration section is now split into "Implementation Tasks（5 项需代码变更）" (lines 21-29) and "Validation Checklist（10 项确认无需改动）" (lines 33-44), with grep-based verification methods specified for each validation item. |
| Attack 3: init-justfile overwrite threat has unresolved agent contradiction | ✅ Yes | Security section (lines 456-461) now specifies a three-layer mechanism: (1) boundary marker merge preserves user recipes, (2) `--force` flag for agent invocation, (3) interactive confirmation only for human users without `--force`. The contradiction with "所有 recipe 不使用交互式输入" is resolved — agents use `--force`, humans get interactive confirmation. |

**Other iteration-2 deductions addressed:**

| Iteration 2 Deduction | Addressed? | Evidence |
|----------------------|------------|----------|
| Recipe table uses shorthand "bash case:" notation | Partially | Two full recipe examples remain (lines 319-336), but the main table (lines 360-377) still uses shorthand. However, the pattern is now well-established with the examples, making the shorthand acceptable for an experienced developer. |
| Error table does not distinguish exit codes (1 vs 127 vs 2) | Partially | Interface 4 Scope Resolution Protocol (lines 219-238) handles exit 127 explicitly. The Error Types table (lines 399-403) still uses generic "exit != 0". The protocol section compensates but the inconsistency remains. |
| "100% PRD coverage" with no traceability matrix | ❌ No | Still claims "100% PRD 覆盖" (line 443) without a test-to-AC mapping table. Flagged in iterations 1, 2, and now 3. |
| Sequential command execution (Story 4 AC) not addressed | ❌ No | Still not addressed. No design for how agent chains `just install && just compile && just test` with partial failure handling. Flagged in iterations 2 and now 3. |

---

## Verdict

- **Score**: 92/100
- **Target**: 90/100
- **Gap**: 0 points (target exceeded)
- **Breakdown-Readiness**: 18/20 — can proceed to `/breakdown-tasks` (well above 12/20 gate)
- **Action**: Target reached. The two persistent issues (test-to-AC traceability, sequential execution AC) are minor and do not block implementation. They should be addressed during implementation but are not design-blocking.
