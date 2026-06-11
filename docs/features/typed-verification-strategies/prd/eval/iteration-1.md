---
date: "2026-05-14"
doc_dir: "docs/features/typed-verification-strategies/prd/"
iteration: 1
target_score: "1000"
scoring_mode: "Mode B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 873/1000** (target: 1000, mode: Mode B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┼──────────┼──────────┬────────────┤
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
│ 3. Flow Completeness (B)     │  185     │  200     │ ✅         │
│    Complete business process │  65/70   │          │            │
│    Data flow documented      │  70/70   │          │            │
│    Exception handling        │  50/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  230     │  300     │ ⚠️         │
│    Coverage per user type    │  65/70   │          │            │
│    Format correct            │  70/70   │          │            │
│    AC per story (G/W/T)      │  55/60   │          │            │
│    AC verifiability          │  70/100  │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  125     │  150     │ ⚠️         │
│    In-scope concrete         │  50/50   │          │            │
│    Out-of-scope explicit     │  40/40   │          │            │
│    Consistent with specs     │  35/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ TOTAL                        │  873     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /70, Data flow documented /70, Exception handling and edge cases /60.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:25 | "Forge 用户（开发者）" is broad — does not specify what kind of projects these developers build | -2 pts |
| prd-spec.md:15 | Background focuses heavily on TUI bugs (11 examples) but goals span all interface types; the API/CLI motivation is comparatively thin | -5 pts |
| prd-spec.md:72-96 | Mermaid diagram: TagE2E/TagInt both merge into Output, but the convergence point is visually ambiguous — GenGeneric also flows to Output but via a separate path that could be clearer | -5 pts |
| prd-spec.md:72-96 | Golden file staleness exception described in text (line 67) and exception table (line 175) but NOT represented in the Mermaid diagram | -5 pts |
| prd-spec.md:57 | Flow step "遍历每个 capability" does not describe loop semantics (sequential vs parallel, what happens on per-capability failure) | -5 pts |
| prd-spec.md:67-68 | "策略文件解析失败" vs "策略文件格式错误" boundary unclear — what distinguishes a parse failure from a format error? | -10 pts |
| prd-user-stories.md:42 | Story 3 AC references "CI lint 失败" but Scope (prd-spec.md:51) explicitly lists "CI/CD 管道改动" as Out of Scope — cross-section inconsistency | -30 pts |
| prd-user-stories.md:16-17 | Story 1 AC bundles TUI + API + CLI verification into a single Then clause; should be separate ACs for independent verifiability | -5 pts |
| prd-user-stories.md | No user story covers gen-test-scripts level-based code generation despite it being a major in-scope deliverable | -15 pts |
| prd-user-stories.md | No user story covers eval-test-cases rubric update or test-cases.md template change despite both being in scope | -10 pts |
| prd-user-stories.md:29 | Story 2 AC does not cover edge case where a capability's interface type mapping is ambiguous or dual-typed | -5 pts |
| prd-spec.md:41 | Scope includes "gen-test-scripts SKILL.md 增强" but no corresponding user story exists | -15 pts (consistency) |
| prd-spec.md:44 | Scope includes "eval-test-cases rubric 更新" but no corresponding user story exists | -10 pts (consistency) |
| prd-spec.md:51 + prd-user-stories.md:42 | CI/CD out of scope but Story 3 AC requires CI lint — repeated cross-section inconsistency | -10 pts |

---

## Attack Points

### Attack 1: User Stories — Three in-scope deliverables have zero story coverage

**Where**: prd-spec.md Scope (lines 41-44) lists 5 in-scope items, but prd-user-stories.md only has stories covering 2 of them (gen-test-cases enhancement and profile strategy files). Missing stories for: gen-test-scripts SKILL.md enhancement (line 42), test-cases.md template update (line 43), eval-test-cases rubric update (line 44).

**Why it's weak**: A PRD that defines deliverables without corresponding user stories leaves implementers without acceptance criteria for 60% of the scope. The gen-test-scripts enhancement alone constitutes a major behavioral change (level-based code generation strategy selection) and deserves at least one story with Given/When/Then ACs. Without stories, there is no verifiable contract for what "different structures" means in practice for e2e vs integration test code.

**What must improve**: Add at least one user story each for gen-test-scripts level-based generation, test-cases.md template changes, and eval-test-cases rubric update. Each story must have concrete ACs. For gen-test-scripts specifically, the ACs should define what "different structures" means (import lists, assertion libraries, directory layout) with examples that can be objectively verified.

### Attack 2: User Stories — Story 3 AC contradicts Scope Out-of-Scope list

**Where**: prd-user-stories.md Story 3 AC (line 42): "否则 CI lint 失败" — but prd-spec.md Scope Out-of-Scope (line 51): "CI/CD 管道改动"

**Why it's weak**: The acceptance criterion requires CI lint enforcement (a CI/CD pipeline change), yet the scope explicitly excludes CI/CD pipeline changes. An implementer cannot simultaneously satisfy both constraints. Either the AC must be rewritten to use a non-CI enforcement mechanism (e.g., "gen-test-cases rejects the strategy file at runtime"), or the scope must be updated to include the necessary CI lint rule. This is not ambiguous — it is a direct logical contradiction that makes the PRD unimplementable as written for Story 3.

**What must improve**: Either (a) change Story 3 AC to validate at gen-test-cases runtime instead of CI lint: "Then gen-test-cases rejects the profile and outputs a validation error listing missing sections", or (b) move "verification-strategies.md CI lint rule" from Out of Scope to In Scope with a specific deliverable. Option (a) is cleaner — it avoids CI scope creep and keeps validation at the point of use.

### Attack 3: Scope Clarity — Scope-to-stories alignment gap leaves acceptance criteria undefined for majority of deliverables

**Where**: prd-spec.md Scope In Scope (lines 40-44) vs prd-user-stories.md (5 stories). Three of five in-scope deliverables (gen-test-scripts enhancement, test-cases.md template, eval-test-cases rubric) have no user stories. prd-spec.md Functional Specs section "级别化代码生成" (lines 131-139) describes gen-test-scripts behavior in detail but provides no acceptance-level verification criteria.

**Why it's weak**: The Functional Specs describe WHAT gen-test-scripts should do ("引用 os/exec + golden file 读取" for e2e, "引用 net/http + assert 库" for integration), but without a user story and Given/When/Then ACs, there is no objective pass/fail test. The phrase "三方面均必须至少存在一处差异" (line 139) is a requirement masquerading as a spec detail — it belongs in an AC. Similarly, "test-cases.md 模板更新：新增 Level 字段" is listed as in-scope but has no story defining what the updated template must contain or how to verify it.

**What must improve**: For each in-scope deliverable without a story, either (a) add a user story with full ACs, or (b) if the change is mechanical (template field addition), fold it into an existing story's AC. For gen-test-scripts, a new story is non-negotiable — the behavioral change is too significant to leave without acceptance criteria.

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 873/1000
- **Target**: 1000/1000
- **Gap**: 127 points
- **Action**: Continue to iteration 2 — primary gaps are in User Stories (missing coverage for 3/5 deliverables and a cross-section inconsistency with CI/CD scope) and Scope Clarity (scope-to-stories misalignment). Background & Goals and Flow Diagrams are strong.
