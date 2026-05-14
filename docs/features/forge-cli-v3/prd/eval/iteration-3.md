---
date: "2026-05-13"
doc_dir: "docs/features/forge-cli-v3/prd/"
iteration: "3"
target_score: "900"
scoring_mode: "MODE_B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 3

**Score: 895/1000** (target: 900, mode: MODE_B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  132     │  150     │ ⚠️         │
│    Three elements            │  43/50   │          │            │
│    Goals quantified          │  34/40   │          │            │
│    Logical consistency       │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  182     │  200     │ ⚠️         │
│    Mermaid diagram exists    │  70/70   │          │            │
│    Main path complete        │  60/70   │          │            │
│    Decision + error branches │  52/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3b. Flow Completeness (B)    │  180     │  200     │ ⚠️         │
│    Complete business process │  63/70   │          │            │
│    Data flow documented      │  65/70   │          │            │
│    Exception & edge cases    │  52/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  262     │  300     │ ⚠️         │
│    Coverage per user type    │  65/70   │          │            │
│    Format correct            │  65/70   │          │            │
│    AC per story (G/W/T)      │  52/60   │          │            │
│    AC verifiability          │  80/100  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  139     │  150     │ ⚠️         │
│    In-scope concrete         │  47/50   │          │            │
│    Out-of-scope explicit     │  37/40   │          │            │
│    Consistent with specs     │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  895     │  1000    │ ⚠️         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /70, Data flow documented /70, Exception handling and edge cases /60.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:Users table:27 | "Hooks/CI（自动化）" persists as a user persona across 3 iterations. It is a system trigger, not a human actor. PRD best practice requires user personas to be people who make decisions. | -7 |
| prd-spec.md:Goals:33 | Goal 4 "品牌一致性" metric is "所有对外命令统一使用 forge 前缀" — binary pass/fail assertion with no measurement procedure, not a quantified metric. | -6 |
| prd-spec.md:Goals:33 | Goal 1 scenario selection methodology still unspecified — "10 个任务场景" is a count but the criteria for selecting those 10 scenarios (coverage of command groups? difficulty levels?) is absent. | -5 |
| prd-spec.md:Agent Flow Diagram:178-200 | Agent Task Execution Flow is the only usage flow with a diagram. Developer flow (lines 84-89) and CI/Hook flow (lines 91-95) have narrative text but no Mermaid diagrams, leaving their decision/error branches undocumented visually. | -10 |
| prd-spec.md:Error Table:99-117 | Error table expanded to 17 rows (improvement confirmed), but `forge task submit` still lacks an entry for "index.json not found" — only covers terminal-state, missing-flag, and concurrency. `forge profile` subcommands (set/detect/get) have zero error coverage despite being listed in the command structure spec. | -8 |
| prd-user-stories.md:Story 4:94 | "基于失败步骤名称去重" is a dedup strategy described in natural language but the matching logic is undefined — exact match on title substring? Step name equality? This makes the AC partially subjective. | -4 |
| prd-user-stories.md:Story 3:66 | "index.json 数据完整无损坏" — "完整无损坏" is not objectively verifiable. A testable AC would specify "index.json remains valid JSON parseable by `forge task status <id>` and the winning submission is reflected." | -4 |
| prd-user-stories.md:Story 6:142 | "stdout 输出空列表或 'no task types defined' 提示" — two acceptable outputs joined by "或" creates ambiguity. Which one should the test assert? | -3 |
| prd-user-stories.md | No user stories cover `forge forensic`, `forge profile`, `forge feature`, or `forge probe` despite all four being documented in the command structure spec (prd-spec.md lines 216-228) and mentioned in Background scenarios. | -12 |
| prd-spec.md:In Scope:53 | "更新 23 个 skills 中的命令引用" — "23 个 skills" is a count without enumeration. After three iterations this remains generic. | -3 |
| prd-spec.md:Out of Scope:60-65 | `forge forensic` commands appear in the command structure spec (line 216) but are absent from both In Scope and Out of Scope lists — an implicit scope item without explicit categorization. | -3 |

---

## Attack Points

### Attack 1: Dimension 4 (User Stories) — Four command groups have zero user stories, creating a coverage gap (65/70 coverage, 80/100 verifiability)

**Where**: prd-user-stories.md contains 6 stories covering: help/discovery (Story 1), prompt (Story 2), task submit (Story 3), hooks/cleanup/quality-gate (Story 4), e2e run (Story 5), task list-types (Story 6). prd-spec.md command structure spec (lines 211-228) defines 5 command groups + 5 top-level commands = 28 total command entry points.

**Why it's weak**: The following commands/groups have zero user story coverage:
- `forge forensic` group (3 subcommands: search, extract, subagents) — mentioned in Background developer scenario but no story.
- `forge profile` group (3 subcommands: set, detect, get) — listed in command spec, zero stories.
- `forge feature` top-level (set/display feature context) — listed in command spec and Background, zero stories.
- `forge probe` top-level (HTTP health check) — listed in command spec, zero stories.
- `forge version` — minor, but also absent.

This means 10 of 28 command entry points (36%) have no user story. While not every command needs a full story (e.g., `forge version` is trivial), `forensic search` and `profile set/detect/get` are substantive features with real user interactions that warrant stories with acceptance criteria. This is the single largest scoring gap at -12 points for coverage and contributes to the cross-section consistency deduction in Dimension 5.

**What must improve**: Add at minimum: (1) a Story 7 for `forge forensic search` covering normal search, no-results, and missing records directory; (2) a Story 8 for `forge profile` covering set, detect, and get with invalid-config edge cases; (3) optionally a Story 9 for `forge feature` covering set/display behavior. These additions would close the coverage gap from 36% uncovered to under 10%.

### Attack 2: Dimension 4 (User Stories) — AC verifiability still has subjectivity in 3 acceptance criteria (80/100)

**Where**: Three specific ACs use imprecise language that cannot be verified without subjective judgment:
- prd-user-stories.md line 94: "基于失败步骤名称去重" — the dedup matching algorithm is unspecified.
- prd-user-stories.md line 66: "index.json 数据完整无损坏" — no definition of "完整" or "无损坏" in testable terms.
- prd-user-stories.md line 142: "stdout 输出空列表或 'no task types defined' 提示" — two acceptable outputs joined by "或" means a test cannot deterministically assert one outcome.

**Why it's weak**: Each of these ACs requires an implementer or tester to make a judgment call. "基于失败步骤名称去重" could mean exact match on step name, substring match, or fuzzy match — each produces different behavior. "数据完整无损坏" could mean "valid JSON", "valid JSON with correct schema", or "valid JSON with correct schema and the winning agent's data persisted". The "或" in Story 6 means two implementations are equally valid, but a test suite must pick one.

**What must improve**: (1) Change "基于失败步骤名称去重" to "fix-task title contains the failing step name as extracted from the quality-gate error output; dedup checks for exact match on this title substring". (2) Change "index.json 数据完整无损坏" to "index.json is valid JSON, parseable by `forge task status T-impl-1`, and T-impl-1.status reflects the winning agent's submitted result". (3) Change "空列表或 'no task types defined'" to a single deterministic output — pick one.

### Attack 3: Dimension 3 (Flow Completeness) — Error table gaps for `forge profile` and `forge task submit` index-missing scenarios (52/60)

**Where**: prd-spec.md Error Handling table (lines 99-117) covers 17 failure scenarios. However: (1) `forge task submit` has no entry for when `index.json` does not exist — the table covers terminal-state conflict, missing flag, and concurrency, but not the basic prerequisite failure. (2) The entire `forge profile` group (set/detect/get — 3 subcommands) has zero error-table entries. (3) `forge feature` has no error-table coverage for invalid or unset feature context.

**Why it's weak**: The error table claims to document "命令级失败行为" (command-level failure behavior) but omits failure modes for 4 commands. For `forge profile set`, what happens if the profile value is unsupported? For `forge profile detect`, what happens if no config.yaml exists? For `forge task submit`, what happens if the task's index.json file is absent? These are not exotic edge cases — they are the most common failure modes for a CLI tool operating on the filesystem. The expansion from 8 to 17 rows (confirmed improvement from iteration 2) addressed many gaps but stopped short of full command coverage.

**What must improve**: Add error-table rows for: (1) `forge task submit` when index.json is missing or unreadable; (2) `forge profile set` with invalid profile value; (3) `forge profile detect` when config.yaml is absent; (4) `forge feature` when feature context is not set. This would bring the table to ~21 rows and achieve near-complete command failure coverage.

---

## Previous Issues Check

| Previous Attack (Iter 2) | Addressed? | Evidence |
|--------------------------|------------|----------|
| Attack 1: Blocked state contradiction + incomplete error table (17 rows needed) | Yes | prd-spec.md line 129 now explicitly states "blocked 不是终态" with note. Error table expanded from 8 to 17 rows including claim, check-deps, validate-index, status, forensic search, verify-task-done, concurrent submit. |
| Attack 2: Missing concurrency AC, fix-task bound AC, invalid-feature AC | Yes | Story 3 lines 64-66 add concurrent submit AC. Story 4 lines 96-98 add fix-task max-count bound. Story 5 lines 122-124 add invalid feature name AC. |
| Attack 3: Diagram naming inconsistency + no retry termination | Yes | Agent Flow Diagram line 189 now shows "forge quality-gate" instead of "just compile...". Lines 191-194 add "RetryCount < 3?" diamond with escalation to "blocked" terminal state. |
| Unresolved from Iter 1-2: Hooks/CI as user persona | No | "Hooks/CI（自动化）" persists in Background Users table and Story 4 role. Not fixed across 3 iterations. |

---

## Verdict

- **Score**: 895/1000
- **Target**: 900/1000
- **Gap**: 5 points
- **Action**: Iterations exhausted. Gap narrowed from 257 (iter 1) to 37 (iter 2) to 5 (iter 3). All three iteration-2 attacks were substantively addressed. Remaining gap is driven by: (1) missing user stories for forensic/profile/feature/probe command groups (-12 coverage), (2) residual AC subjectivity in 3 instances (-11 verifiability), (3) Hooks/CI persona issue persisted across all 3 iterations (-7). To reach 900+, add Stories 7-8 for forensic and profile commands and tighten the 3 subjective ACs.

SCORE: 895/1000
