---
date: "2026-05-13"
doc_dir: "docs/features/forge-cli-v3/prd/"
iteration: "4"
target_score: "900"
scoring_mode: "MODE_B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 4

**Score: 925/1000** (target: 900, mode: MODE_B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 1. Background & Goals        │  132     │  150     │ ⚠️         │
│    Three elements            │  43/50   │          │            │
│    Goals quantified          │  34/40   │          │            │
│    Logical consistency       │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 2. Flow Diagrams             │  185     │  200     │ ⚠️         │
│    Mermaid diagram exists    │  70/70   │          │            │
│    Main path complete        │  63/70   │          │            │
│    Decision + error branches │  52/60   │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 3b. Flow Completeness (B)    │  188     │  200     │ ⚠️         │
│    Complete business process │  66/70   │          │            │
│    Data flow documented      │  65/70   │          │            │
│    Exception & edge cases    │  57/60   │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 4. User Stories              │  282     │  300     │ ⚠️         │
│    Coverage per user type    │  67/70   │          │            │
│    Format correct            │  68/70   │          │            │
│    AC per story (G/W/T)      │  55/60   │          │            │
│    AC verifiability          │  92/100  │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 5. Scope Clarity             │  138     │  150     │ ⚠️         │
│    In-scope concrete         │  46/50   │          │            │
│    Out-of-scope explicit     │  37/40   │          │            │
│    Consistent with specs     │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ TOTAL                        │  925     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /70, Data flow documented /70, Exception handling and edge cases /60.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:Users table:27 | "Hooks/CI（自动化）" persists as a user persona across 4 iterations. It is a system trigger, not a human actor with intent. PRD best practice requires user personas to be people who make decisions. | -7 |
| prd-spec.md:Goals:33 | Goal 4 "品牌一致性" metric is "所有对外命令统一使用 forge 前缀" — a binary pass/fail assertion with no measurement procedure or verification script, not a quantified metric. | -6 |
| prd-spec.md:Goals:33 | Goal 1 scenario selection methodology unspecified — "10 个任务场景" is a count but the criteria for selecting those 10 scenarios (which command groups? difficulty distribution?) remains absent after 4 iterations. | -5 |
| prd-spec.md:Developer+CI Flow:84-95 | Developer flow (lines 84-89) and CI/Hook flow (lines 91-95) remain narrative-only without Mermaid diagrams. Only Agent flow has a diagram. Two of three documented usage flows lack visual decision/error branch coverage. | -7 |
| prd-spec.md:Error Table:119-120 | Error table expanded with `forge profile set` and `forge profile get` invalid-profile rows (confirmed improvement). However, `forge profile detect` still has no error-table entry — what happens when config.yaml is absent or no test framework config files are detected? `forge feature` (set/display) also lacks any error-table coverage. | -3 |
| prd-user-stories.md:Story 8:166 | Story 7 `forge forensic extract` AC says "退出码为 1，stderr 输出 'file not found: /nonexistent/path.jsonl'" — the AC hardcodes a user-supplied path in the expected error message. This is a valid pattern but `forge forensic subagents` (line 162-163) has no error-path AC at all (e.g., nonexistent session-dir-path). | -3 |
| prd-user-stories.md:Story 7 | `forge forensic search` AC (line 155) says "输出匹配的会话列表（包含 session ID、时间戳、skill 名称）" but no AC covers the case where history.jsonl is empty or contains zero matches — only `records/ 目录不存在` is in the error table, not in Story 7 ACs. | -3 |
| prd-user-stories.md | `forge feature`, `forge probe`, and `forge version` remain without user stories. These are lower-priority commands but `forge feature` (set/display feature context) is used by "Agent, 开发者" per the command spec table (prd-spec.md line 227) and was flagged in iteration 3. | -3 |
| prd-user-stories.md:Story 6:142 | "stdout 输出空列表（0 行类型记录），退出码为 0" — improved from "或" ambiguity (confirmed fix). However the word "列表" is still vague: is it literally zero stdout bytes, or a header row with no data rows? A test cannot distinguish "empty list" from "no output" without a precise definition. | -2 |
| prd-spec.md:In Scope:53 | "更新 23 个 skills 中的命令引用" — "23 个 skills" remains a count without enumeration after 4 iterations. An explicit list or at least a glob pattern would eliminate ambiguity. | -4 |
| prd-spec.md:Scope | `forge feature`, `forge probe`, and `forge version` appear in the command structure spec (lines 226-232) but are absent from both In Scope and Out of Scope lists — implicit scope items without explicit categorization. | -2 |

---

## Attack Points

### Attack 1: Dimension 4 (User Stories) — Three commands still lack stories despite iteration-3 flagging (67/70 coverage, -3)

**Where**: prd-user-stories.md now has 8 stories. prd-spec.md command structure spec (lines 226-232) defines `forge feature`, `forge probe`, and `forge version` as top-level commands. None of these have user stories. Iteration 3 explicitly flagged "forge feature" and "forge probe" in Attack 1.

**Why it's weak**: `forge feature` is listed as used by "Agent, 开发者" (prd-spec.md line 227) — it is not a trivial command. An agent or developer calling `forge feature` with no feature context set, or with an invalid feature slug, has no defined acceptance criteria. The forensic and profile gaps were fixed (Stories 7-8 confirmed), but feature/probe/version were left uncovered despite being part of the same attack. Partial fix = still a gap.

**What must improve**: Add a Story 9 for `forge feature` covering: set with valid slug, set with invalid/missing slug, display current context when unset. `forge probe` and `forge version` are minor enough to omit but should be explicitly noted as out-of-scope for stories in a comment, or given trivial 1-AC stories.

### Attack 2: Dimension 1 (Background & Goals) — Hooks/CI persona and unquantified metrics persist after 4 iterations (132/150)

**Where**: prd-spec.md Background Users table (line 27): "Hooks/CI（自动化）| 每次会话/提交 | forge cleanup, forge quality-gate, forge verify-task-done". Goals table (line 33): Goal 4 "品牌一致性" metric "所有对外命令统一使用 forge 前缀". Goal 1 metric "10 个任务场景".

**Why it's weak**: These three issues have been flagged in every iteration (1, 2, 3, and now 4) and never addressed. "Hooks/CI（自动化）" is a system trigger, not a persona — it does not "want" anything or make decisions. The brand-consistency goal has no test procedure (how do you verify "all" commands? grep? CI check?). The 10-scenario selection criteria are undefined. Three iterations of non-response suggests these are being treated as acceptable, but they represent -18 cumulative points that prevent the document from reaching higher quality.

**What must improve**: (1) Replace "Hooks/CI（自动化）" with a note: "Hooks and CI are automation triggers; their behavior is covered in Story 4." Keep the table to human personas only (AI agent, developer). (2) Change Goal 4 metric to: "CI target `just check-stale-refs` finds zero `task` command references after migration (grep -r 'task ' across hooks/, skills/, docs/)." (3) Add one sentence to Goal 1: "10 scenarios selected to cover each command group (task/e2e/forensic/profile/prompt) with 2 scenarios per group."

### Attack 3: Dimension 2 (Flow Diagrams) — Developer and CI/Hook flows still have no Mermaid diagrams (185/200)

**Where**: prd-spec.md lines 84-95 describe Developer and CI/Hook usage flows as narrative text only. The only Mermaid diagrams are "Business Flow Diagram" (lines 144-178, covering migration phases) and "Agent Task Execution Flow" (lines 180-205).

**Why it's weak**: The rubric criterion "Decision points + error branches covered" requires diamond nodes and error/exception branches. The Developer flow has a conditional (profile detection → choose test suite) and the CI/Hook flow has three distinct trigger paths with different failure modes (cleanup no-op, quality-gate failure creating fix-task, verify-task-done blocking). None of these decision points or error branches are captured in a diagram. The Agent flow diagram covers its path well but it is only 1 of 3 documented usage flows.

**What must improve**: Add two Mermaid diagrams: (1) Developer Flow covering `forge e2e run` → profile detection → suite selection → pass/fail, plus `forge task list-types` and `forge forensic search` branches. (2) CI/Hook Flow covering SessionEnd → cleanup, Stop → quality-gate → fix-task creation, PreToolUse → verify-task-done → block/allow. Each diagram should have at least one diamond decision node and one error branch.

---

## Previous Issues Check

| Previous Attack (Iter 3) | Addressed? | Evidence |
|--------------------------|------------|----------|
| Attack 1: Missing user stories for forensic/profile/feature/probe (4 command groups, -12 pts) | Partially | Stories 7 (forensic) and 8 (profile) added with G/W/T ACs including error paths. `forge feature`, `forge probe`, `forge version` still uncovered. Gap reduced from 10 uncovered commands to 3. |
| Attack 2: Three subjective ACs — dedup logic, index.json integrity, "或" ambiguity | Partially | Story 4 dedup AC (line 94) rewritten with precise title format "fix-compile-3" and grep count verification — confirmed fix. Story 6 "或" ambiguity resolved — now single deterministic output "stdout 输出空列表（0 行类型记录）" (line 142). Story 3 "index.json 数据完整无损坏" — the phrase was removed from the concurrent-submit AC and replaced with "index.json 可被 `jq .` 正常解析且 JSON 语法合法" (line 66) — confirmed fix. Minor residual: "0 行类型记录" could be more precise. |
| Attack 3: Error table gaps for profile and submit-index-missing | Mostly | Error table now includes `forge profile set` invalid profile (line 119), `forge profile get` invalid profile (line 120), `forge e2e run` feature-not-found (line 121), `forge task submit` index-missing (line 118). `forge profile detect` (no config.yaml / no frameworks detected) still absent. `forge feature` error cases still absent. |
| Persistent across all iterations: Hooks/CI as user persona | No | "Hooks/CI（自动化）" remains in prd-spec.md Background Users table (line 27). Unchanged across 4 iterations. |

---

## Verdict

- **Score**: 925/1000
- **Target**: 900/1000
- **Gap**: -25 (target exceeded by 25 points)
- **Action**: Target reached. Score improved from 895 (iter 3) to 925 (iter 4), a +30 point gain. The addition of Stories 7-8 and the error-table expansion are genuine improvements with substantive ACs. Remaining imperfections are minor: 3 uncovered commands (feature/probe/version), persistent Hooks/CI persona issue, and missing Mermaid diagrams for developer/CI flows. These are below the threshold for further mandatory iteration.

SCORE: 925/1000
