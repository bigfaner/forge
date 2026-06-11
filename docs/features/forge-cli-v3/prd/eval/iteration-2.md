---
date: "2026-05-13"
doc_dir: "docs/features/forge-cli-v3/prd/"
iteration: "2"
target_score: "900"
scoring_mode: "MODE_B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 863/1000** (target: 900, mode: MODE_B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  128     │  150     │ ⚠️         │
│    Three elements            │  45/50   │          │            │
│    Goals quantified          │  33/40   │          │            │
│    Logical consistency       │  50/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  165     │  200     │ ⚠️         │
│    Mermaid diagram exists    │  70/70   │          │            │
│    Main path complete        │  55/70   │          │            │
│    Decision + error branches │  40/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3b. Flow Completeness (B)    │  160     │  200     │ ⚠️         │
│    Complete business process │  60/70   │          │            │
│    Data flow documented      │  60/70   │          │            │
│    Exception & edge cases    │  40/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  265     │  300     │ ⚠️         │
│    Coverage per user type    │  65/70   │          │            │
│    Format correct            │  65/70   │          │            │
│    AC per story (G/W/T)      │  55/60   │          │            │
│    AC verifiability          │  80/100  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  145     │  150     │ ✅         │
│    In-scope concrete         │  48/50   │          │            │
│    Out-of-scope explicit     │  38/40   │          │            │
│    Consistent with specs     │  59/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  863     │  1000    │ ⚠️         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /70, Data flow documented /70, Exception handling and edge cases /60.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:Users table:27 | "Hooks/CI（自动化）" listed as user persona — system actor, not a human user. Same issue as iter 1, unresolved. | -5 |
| prd-spec.md:Goals:33 | Goal 4 "品牌一致性" metric is "所有对外命令统一使用 forge 前缀" — binary pass/fail, not quantified with a numeric target or measurement procedure. | -7 |
| prd-spec.md:Goals:33 | Goal 1 "LLM 命令选择测试" still underspecified — no mention of which LLM, how many test runs, confidence interval, or test prompt construction methodology. "10 个任务场景" is stated but the scenario selection criteria are missing. | -8 |
| prd-spec.md:Business Flow Diagram:129-166 | Agent Task Execution Flow diagram shows "Quality Gate: just compile → fmt → lint → test" as inline steps, but the quality-gate command is `forge quality-gate` per the spec. Inconsistent naming between diagram and command spec. | -5 |
| prd-spec.md:Agent Flow Diagram:179-183 | Diagram shows `fix → quality` loop but no terminal condition for infinite loop (what if fix keeps failing?). No max-retry or escalation branch. | -8 |
| prd-spec.md:Agent Flow Diagram:175 | Diagram starts with `forge task claim` but the Background's primary scenario (line 25) says `forge prompt get-by-task-id` is the typical entry point. Ordering inconsistency. | -7 |
| prd-spec.md:Error Handling:112-118 | Task state transition table shows "blocked → in_progress" as allowed with error "task already in terminal state" — contradictory: blocked is NOT a terminal state in the table header, yet the error message says it is. The table's "允许的目标状态" column for `blocked` shows `in_progress` but the error column says "task already in terminal state". | -8 |
| prd-spec.md:Error Handling:99-108 | Error table covers 8 failure scenarios across 5 commands, but many commands are missing: `forge task claim` (what if all tasks claimed?), `forge task check-deps` (what if deps not met?), `forge task validate-index` (what if index corrupt?), `forge forensic search` (what if no results?), `forge task status` (what if no task found?). | -12 |
| prd-user-stories.md:Story 4:69 | "As a CI/Hook 自动化流程" — system actor, inconsistent with Background user types which should be human roles per PRD conventions. Same as iter 1. | -5 |
| prd-user-stories.md:Story 1:17 | "每个子命令描述包含'命令名+动词+宾语'三要素" — "三要素" is defined but the AC then says "描述长度 <= 80 字符". The first assertion is partially subjective (what counts as "动词+宾语"?), though the length constraint helps. | -3 |
| prd-user-stories.md:Story 4:89-90 | Fix-task recursive failure AC: "创建新的 P0 fix-task（不覆盖已有 fix-task）" — no AC covers the scenario where fix-tasks accumulate without bound. No de-duplication or max-count constraint specified. | -5 |
| prd-user-stories.md:Story 3 | No AC for concurrent submit race condition (two agents submit same task simultaneously). Given that index.json.lock is mentioned in Security Requirements, this is a notable gap. | -8 |
| prd-user-stories.md:Story 5 | No AC for `forge e2e run --feature <name>` when the named feature does not exist. The e2e command references a feature but no error AC covers invalid feature name. | -4 |
| prd-spec.md:In Scope:53 | "更新 23 个 skills 中的命令引用" — still says "23 个 skills" generically without naming them. Related Changes table also uses "23 个 skills". | -2 |

---

## Attack Points

### Attack 1: Dimension 3 (Flow Completeness) — Exception handling table incomplete and internally contradictory (40/60)

**Where**: prd-spec.md Error Handling section (lines 99-118). The state transition table at lines 112-118 contains a contradiction: the `blocked` row shows `in_progress` as an allowed target state, yet the "非法转换行为" column says "task already in terminal state" — but `blocked` is not a terminal state (it has a valid transition target listed).

**Why it's weak**: (1) The internal contradiction makes the spec unverifiable — an implementer cannot know whether blocked→in_progress is allowed or forbidden. (2) The error table at lines 99-108 covers only 5 of the 23 commands in the spec. Missing commands include `forge task claim` (no available tasks), `forge task check-deps`, `forge task validate-index`, `forge task status`, `forge forensic search`, and `forge e2e setup`. (3) The migration failure strategy (lines 120-125) describes recovery actions but provides no quantitative criteria for "when to abort the entire migration vs. fix and retry."

**What must improve**: Fix the blocked state contradiction (clarify whether blocked is terminal or not; align the "允许的目标状态" and error message columns). Expand the error table to cover at minimum all commands that can fail with a non-zero exit code. Add quantitative abort criteria to the migration failure strategy (e.g., "after 3 consecutive Phase 3 failures, escalate to manual review").

### Attack 2: Dimension 4 (User Stories) — AC verifiability still has gaps in concurrency and resource-exhaustion scenarios (80/100)

**Where**: prd-user-stories.md Story 3 (lines 49-64) and Story 4 (lines 67-91).

**Why it's weak**: Story 3 (submit) has no AC for concurrent access — the Security Requirements mention index.json.lock but no AC verifies that concurrent submissions to the same task produce correct behavior (one succeeds, one gets "already in terminal state"). Story 4's fix-task recursion AC (lines 89-90) specifies that new fix-tasks are created without overwriting existing ones, but there is no bound on accumulation — a repeatedly failing quality-gate could create arbitrarily many P0 fix-tasks with no de-duplication rule. Story 5 has no AC for invalid feature name input. These are precisely the boundary conditions that distinguish a verifiable spec from a reviewable one.

**What must improve**: Add ACs for: (1) concurrent `forge task submit` on the same task ID — verify lock behavior and "already in terminal state" on the loser; (2) fix-task accumulation — add a de-duplication or max-count AC in Story 4; (3) invalid feature name in Story 5's `forge e2e run --feature <name>`.

### Attack 3: Dimension 2 (Flow Diagrams) — Agent Task Execution Flow diagram has naming inconsistency and no failure termination (55/70 main path, 40/60 error branches)

**Where**: prd-spec.md Agent Task Execution Flow diagram (lines 169-188).

**Why it's weak**: (1) The diagram uses "Quality Gate: just compile → fmt → lint → test" (line 179) but the command spec defines `forge quality-gate` as the CLI command. The diagram should reference the actual CLI command, not the internal justfile recipes. (2) The fix→quality loop (lines 181-182) has no termination condition — no max-retry diamond, no "escalate to human" branch. This means the diagram implicitly allows infinite loops, which is an incomplete specification of the main path. (3) The diagram starts with `forge task claim` but the Background primary scenario (line 25) leads with `forge prompt get-by-task-id` — these two entry points are never reconciled in the diagram.

**What must improve**: Replace "Quality Gate: just compile..." with "forge quality-gate" to match the command spec. Add a max-retry or escalation decision diamond in the fix loop. Either reconcile the claim vs. prompt entry-point ordering in the Background section, or add a second diagram path showing the prompt-first flow.

---

## Previous Issues Check

| Previous Attack (Iter 1) | Addressed? | Evidence |
|--------------------------|------------|----------|
| Attack 1: Zero error-case ACs in all 6 stories | Partially | Stories 1-5 now have error-case ACs (invalid IDs, terminal state, missing profile, unknown profile). Story 6 has edge case (empty registry). But Story 3 lacks concurrency AC, Story 4 lacks fix-task bound, Story 5 lacks invalid-feature AC. |
| Attack 2: Zero exception handling narrative text | Yes | New "Error Handling" subsection (lines 96-125) with command-level failure table, state transition constraints table, and migration failure strategy. Internally contradictory on blocked state. |
| Attack 3: Only agent flow has narrative; developer and CI flows missing | Yes | New "开发者使用流程" (lines 85-89) and "CI/Hook 使用流程" (lines 91-95) added. Both cover trigger→processing→end state. Developer flow could be more detailed on e2e subcommands beyond `run`. |

---

## Verdict

- **Score**: 863/1000
- **Target**: 900/1000
- **Gap**: 37 points
- **Action**: Continue to iteration 3 — focus on: (1) fixing the blocked-state contradiction and expanding error table to all commands (+20 pts in Dim 3), (2) adding concurrency/race-condition ACs and fix-task bound (+15 pts in Dim 4), (3) fixing diagram naming inconsistency and adding retry termination (+15 pts in Dim 2).

SCORE: 863/1000
