---
date: "2026-05-14"
doc_dir: "docs/features/typed-verification-strategies/prd/"
iteration: 2
target_score: "1000"
scoring_mode: "Mode B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 920/1000** (target: 1000, mode: Mode B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┼──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  143     │  150     │ ✅         │
│    Three elements            │  48/50   │          │            │
│    Goals quantified          │  40/40   │          │            │
│    Logical consistency       │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  190     │  200     │ ✅         │
│    Mermaid diagram exists    │  70/70   │          │            │
│    Main path complete        │  65/70   │          │            │
│    Decision + error branches │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3b. Flow Completeness (B)    │  190     │  200     │ ✅         │
│    Complete business process │  68/70   │          │            │
│    Data flow documented      │  70/70   │          │            │
│    Exception handling        │  52/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  260     │  300     │ ⚠️         │
│    Coverage per user type    │  70/70   │          │            │
│    Format correct            │  70/70   │          │            │
│    AC per story (G/W/T)      │  60/60   │          │            │
│    AC verifiability          │  60/100  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  137     │  150     │ ⚠️         │
│    In-scope concrete         │  50/50   │          │            │
│    Out-of-scope explicit     │  40/40   │          │            │
│    Consistent with specs     │  47/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  920     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /70, Data flow documented /70, Exception handling and edge cases /60.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:25 | "Forge 用户（开发者）" is broad — does not further specify project types, tech stacks, or experience level | -2 pts |
| prd-spec.md:15 | Background motivates with TUI bug examples (11 instances) but goals span all interface types; the CLI/API motivation is comparatively thin relative to its equal billing in goals | -3 pts |
| prd-spec.md:72-96 | Mermaid diagram: TagE2E/TagInt both merge into Output node, but the convergence is visually compressed — GenGeneric also flows to Output via a separate join, making the output aggregation point hard to trace for a reader | -5 pts |
| prd-spec.md:66 | "策略文件 section 不完整（验证维度 <3 或边界场景 <2）：gen-test-cases 拒绝 + 验证错误，中止执行" — this exception is described in text and the exception table (line 172) but is NOT represented in the Mermaid diagram as a distinct decision node (ValidFormat node covers format but not section completeness counts) | -5 pts |
| prd-spec.md:67 | Golden file staleness exception described in text (line 67) and exception table (line 175) but NOT represented in the Mermaid diagram | -3 pts |
| prd-spec.md:57 | Flow step "遍历每个 capability" does not describe loop semantics — what happens if one capability fails mid-loop? Are already-generated cases preserved or discarded? | -2 pts |
| prd-spec.md:66 vs prd-spec.md:172 | "策略文件格式错误" in the exception table (line 172) vs "section 不完整" in the flow description (line 66) — these are the same condition described with different terminology, creating ambiguity about whether format errors and completeness errors are one check or two | -8 pts |
| prd-spec.md:66 vs prd-spec.md:172 | "策略文件解析失败" is not in the exception table — it appears only in the flow description at line 62 ("策略文件格式是否有效"). The exception table only has "格式错误" — the boundary between parse failure and format error is still unclear | -6 pts |
| prd-user-stories.md:16 | Story 1 AC bundles TUI + API + CLI verification into a single Then clause: "TUI 用例包含... API 用例包含... CLI 用例包含..." — should be separate ACs for each interface type for independent verifiability | -5 pts |
| prd-user-stories.md:29 | Story 2 AC: "Level 字段覆盖率 ≥ 95%" — the 95% metric is not operationally defined. 95% of what total? Total capabilities? Total test cases? How is this measured? | -5 pts |
| prd-user-stories.md:42 | Story 3 AC: "否则 gen-test-cases 拒绝该 profile" — "拒绝" is ambiguous: does it mean abort entirely or skip that specific capability? The Mermaid diagram (prd-spec.md:79-80) shows "Reject → Done" which means full abort, but the AC does not clarify this | -5 pts |
| prd-user-stories.md:81 | Story 6 AC: "两类代码在 import 列表、assertion 方式、目录结构三方面均存在至少一处差异" — "至少一处差异" means only 1 of 3 aspects needs to differ, which is an extremely low bar for a feature claiming "不同结构的测试代码" | -5 pts |
| prd-user-stories.md:94 | Story 7 AC: "Level 和 Interface 字段覆盖率 ≥ 95%" — same measurement ambiguity as Story 2; what is the denominator? | -5 pts |
| prd-user-stories.md:107 | Story 8 AC: "当无策略文件时该维度得分率 ≥ 0.8（不因缺少策略而惩罚）" — this is a scoring floor for a missing-feature case, which is an odd AC. It says the system gives a high score (0.8) when the feature is absent — this contradicts the purpose of the scoring dimension | -8 pts |
| prd-user-stories.md | Story 8 AC: "完整度得分 = 已覆盖检查项数 / 总检查项数 × 权重分值" — "权重分值" is not defined anywhere in the PRD. The AC references a weight that has no specified value | -5 pts |
| prd-spec.md:44 | Scope includes "eval-test-cases rubric 更新" but no user story covers the rubric content or scoring criteria change — Story 8 covers the scoring behavior but not the rubric deliverable itself | -5 pts |
| prd-spec.md:40 | Scope includes "6 个 profile 各新增 verification-strategies.md" — no user story covers the creation of strategy files for the 6 specific profiles; Story 3 covers the format but not the actual 6-profile deliverable | -3 pts |
| prd-spec.md:139 + prd-user-stories.md:81 | "三方面均必须至少存在一处差异" in spec vs "三方面均存在至少一处差异" in story — both use the same weak bar (1 of 3) despite claiming "不同结构" | -4 pts |

---

## Attack Points

### Attack 1: User Stories — AC verifiability undermined by undefined metrics and weak thresholds

**Where**: prd-user-stories.md Stories 2 (line 29), 6 (line 81), 7 (line 94), 8 (line 107). Story 2: "Level 字段覆盖率 ≥ 95%"; Story 6: "至少一处差异"; Story 8: "权重分值" (undefined).

**Why it's weak**: Multiple ACs contain metrics that cannot be objectively verified. The "≥ 95%" coverage targets in Stories 2 and 7 never define the denominator — is it test cases, capabilities, or something else? Story 6's "至少一处差异" threshold means the AC passes if only the directory name differs while imports and assertions are identical — this is a trivially satisfiable bar for a feature claiming structural differentiation. Story 8 references "权重分值" which is never defined in the PRD, making the AC formula uncomputable. These are not edge-case nitpicks; they are the core verification criteria that an implementer or tester would need to check, and they are ambiguous or incomplete.

**What must improve**: (a) Define the coverage denominator explicitly in Stories 2 and 7 (e.g., "95% of all test cases in test-cases.md that have a capability mapping"). (b) Raise Story 6's threshold to "at least 2 of 3 aspects differ" or remove the minimum and require all 3. (c) Define the weight value for Story 8's formula or replace the formula with a concrete pass/fail threshold.

### Attack 2: Flow Completeness — Exception handling terminology inconsistency between flow description and exception table

**Where**: prd-spec.md lines 62-66 (Business Flow Description decision points) vs lines 169-175 (Exception Handling table). Line 62: "策略文件格式是否有效"; line 66: "策略文件 section 不完整（验证维度 <3 或边界场景 <2）"; line 172: "策略文件格式错误（section 缺失验证维度或边界场景）".

**Why it's weak**: The flow description lists "策略文件格式是否有效" as decision point #3 (line 62) and "策略文件 section 不完整" as a separate exception (line 66). The exception table collapses these into a single row: "策略文件格式错误（section 缺失验证维度或边界场景）". An implementer cannot determine whether there are two distinct validation checks (format validity + section completeness) or one combined check. The Mermaid diagram shows only one decision node (ValidFormat), suggesting one check, but the text describes what appear to be two different failure modes. This ambiguity makes the error handling unimplementable without guessing the designer's intent.

**What must improve**: Consolidate into one clearly defined validation step with explicit criteria: "策略文件验证" checks (1) valid markdown section structure, (2) each capability section has ≥3 verification dimensions, (3) each capability section has ≥2 boundary scenarios. All three failures produce the same abort behavior. Update both the flow description and exception table to use identical terminology. Consider adding a sub-decision in the Mermaid diagram if these are genuinely separate checks.

### Attack 3: Scope Clarity — Scope-to-stories alignment still has gaps for profile creation deliverable

**Where**: prd-spec.md Scope In Scope (line 40): "6 个 profile 各新增 verification-strategies.md（go-test、web-playwright、maestro、pytest、rust-test、java-junit）". No user story covers this deliverable.

**Why it's weak**: The scope explicitly lists creating 6 profile strategy files as an in-scope deliverable, complete with specific profile names. Story 3 covers the format specification and validation of strategy files, but no story says "As a Profile author, I want each of the 6 built-in profiles to ship with a verification-strategies.md" with ACs defining what each profile must contain. The scope treats this as a concrete deliverable (6 files with specific names), but the user stories treat it as a format spec. An implementer reading only the user stories would not know they need to create 6 files — they would think they only need to build the validation logic. Similarly, Story 8 covers eval-test-cases scoring behavior but the scope item "eval-test-cases rubric 更新" is the rubric deliverable itself, not just the scoring behavior. These gaps mean the scope promises deliverables that the stories do not verify.

**What must improve**: Either (a) add a user story for the 6-profile deliverable with ACs like "Given the 6 built-in profiles, When verification-strategies.md is checked, Then each profile has a valid strategy file with capability-specific sections", or (b) fold it into Story 3's AC by adding: "Then the 6 built-in profiles (go-test, web-playwright, maestro, pytest, rust-test, java-junit) each ship with a valid verification-strategies.md". For eval-test-cases rubric, clarify in Story 8's AC what the rubric deliverable contains (not just the scoring behavior).

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Three in-scope deliverables have zero story coverage (gen-test-scripts, test-cases.md template, eval-test-cases rubric) | ✅ Partially | Story 6 now covers gen-test-scripts level-based generation. Story 7 covers test-cases.md Level/Interface fields. Story 8 covers eval-test-cases scoring. However, the eval-test-cases *rubric deliverable* and *6 profile creation deliverable* remain without explicit story coverage. |
| Attack 2: Story 3 AC contradicts Scope Out-of-Scope list (CI lint vs CI/CD out of scope) | ✅ | Story 3 AC (line 42) now says "gen-test-cases 拒绝该 profile 并输出验证错误" — validation at runtime, not CI lint. The CI/CD contradiction is resolved. |
| Attack 3: Scope-to-stories alignment gap for majority of deliverables | ✅ Partially | Stories 6, 7, 8 added, covering the three previously missing deliverables. Remaining gap: the 6-profile creation deliverable (line 40) still has no story that verifies the actual files exist for each named profile. |

---

## Verdict

- **Score**: 920/1000
- **Target**: 1000/1000
- **Gap**: 80 points
- **Action**: Continue to iteration 3 — primary gaps are in User Stories AC verifiability (undefined metrics like "权重分值", weak thresholds like "至少一处差异", and undefined coverage denominators) and Scope Clarity (6-profile creation deliverable still lacks story coverage, eval-test-cases rubric deliverable vs behavior gap). Background & Goals and Flow Diagrams remain strong. Flow Completeness has a terminology inconsistency between flow description and exception table that should be resolved.
