---
date: "2026-04-28"
doc_dir: "docs/proposals/e2e-script-quality/"
iteration: 4
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 4

**Score: 90/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Problem Definition        │  18      │  20      │ ✅         │
│    Problem clarity           │  7/7     │          │            │
│    Evidence provided         │  6/7     │          │            │
│    Urgency justified         │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Solution Clarity          │  19      │  20      │ ✅         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  6.5/7   │          │            │
│    Differentiated            │  5.5/6   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Alternatives Analysis     │  14      │  15      │ ✅         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  4/5     │          │            │
│    Rationale justified       │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Scope Definition          │  15      │  15      │ ✅         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  5/5     │          │            │
│    Scope bounded             │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Risk Assessment           │  14      │  15      │ ✅         │
│    Risks identified (≥3)     │  5/5     │          │            │
│    Likelihood + impact rated │  4.5/5   │          │            │
│    Mitigations actionable    │  4.5/5   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Success Criteria          │  15      │  15      │ ✅         │
│    Measurable                │  5/5     │          │            │
│    Coverage complete         │  5/5     │          │            │
│    Testable                  │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  95 → 90│  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘

Deduction adjustments: 95 - 3 (Alternatives: rejected alternative "扩展 config.yaml" still has vague con "新增维护负担；routes 变更需要手动同步" — lacks quantification comparable to the chosen approach's con) - 2 (Evidence: still single-project N=1 sample, no generalizability qualifier) = 90
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Alternatives table, "扩展 config.yaml" row | Con column says "新增维护负担；routes 变更需要手动同步" — vague and unquantified. The chosen approach's con specifies exact file counts ("4 个 skill 的 8 个文件"), but this rejected alternative gets a hand-wave. Asymmetric depth tilts the comparison. | -3 pts (within Alternatives: pros/cons honesty) |
| Evidence section | All evidence from a single project (pm-work-tracker, now expanded to 22 features / 13 with testing scripts — a genuine improvement in sample size within that project). Still N=1 across projects. No generalizability qualifier. | -2 pts (within Evidence provided) |

---

## Attack Points

### Attack 1: Alternatives Analysis — rejected "扩展 config.yaml" con is still a hand-wave

**Where**: Alternatives table, row "扩展 config.yaml 增加 routes 段", Cons column: "新增维护负担；routes 变更需要手动同步"
**Why it's weak**: This was iteration 3's Attack 3 and it remains unaddressed. The chosen approach's con is precisely quantified ("涉及修改 gen-test-cases、gen-test-scripts、run-e2e-tests、graduate-tests 共 4 个 skill 的 8 个文件，需协调测试"). But the rejected "扩展 config.yaml" alternative gets two vague phrases with no quantification. "新增维护负担" — how much burden? One YAML block? A full parallel manifest? "routes 变更需要手动同步" — sync between how many files? How often do routes change? The asymmetry is striking: the proposal author gave the chosen approach honest, specific self-criticism, but gave the rejected alternative a straw-man con. This matters because a reader evaluating trade-offs cannot fairly compare "8 files across 4 skills" (concrete, scorable) against "新增维护负担" (abstract, unfalsifiable).
**What must improve**: Replace with a concrete con: e.g., "Requires maintaining a parallel route manifest in config.yaml duplicating information already in router files; every route change must be synced in 2+ places (estimated 1-2 min overhead per route change), and config drift is invisible until tests fail."

### Attack 2: Problem Definition — evidence is still single-project (N=1)

**Where**: Evidence section header: "实际项目中（pm-work-tracker，22 个 feature、13 个含 testing scripts）观察到以下问题"
**Why it's weak**: The evidence has expanded substantially from iteration 3 (was "3 个 feature、9 个 test script") to iteration 4 ("22 个 feature、13 个含 testing scripts"). The within-project sample is now large enough to be meaningful. However, all data still comes from pm-work-tracker. The proposal makes no statement about generalizability. Is pm-work-tracker a typical Go/Fiber project? Does it use standard forge pipeline conventions? The reader has no way to assess whether these error patterns (port guessing, prefix mismatch, hardcoded paths) are specific to this project's conventions or universal to the pipeline's design flaw. A single qualifying sentence would close this gap.
**What must improve**: Add one sentence after the evidence header: e.g., "These error categories (port guessing, route prefix mismatch, hardcoded paths) are expected to generalize because they stem from the pipeline's lack of source code reading, which is architecture-agnostic rather than project-specific."

### Attack 3: Risk Assessment — shared helpers.ts conflict risk mitigation is thin

**Where**: Risk table, row "共享 helpers.ts 的 feature 间冲突", Mitigation: "helpers.ts 从 `tests/e2e/config.yaml` 读取配置，config 是项目级单一来源"
**Why it's weak**: The mitigation assumes all config can be expressed in a single config.yaml, but the proposal does not define what goes into config.yaml vs. what stays in helpers.ts. What if two features need different auth flows (e.g., one uses Bearer tokens, another uses session cookies)? The mitigation says "config 是项目级单一来源" but never specifies the config schema or what happens when features genuinely need divergent behavior. The risk is rated Low likelihood, but the proposal's own evidence shows helpers.ts implementations already diverge in the current system ("loginViaUI 实现分化：一份用硬编码中文，一份用 regex"). Consolidating divergent implementations into one shared file is exactly where conflicts arise, and the mitigation does not address this.
**What must improve**: Either upgrade the likelihood to Medium given the existing divergence evidence, or add a concrete mitigation: "helpers.ts will use strategy parameters (e.g., `loginViaUI({ method: 'regex' })`) to accommodate feature-specific variations without duplicating code."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Problem Definition — single-project evidence, no cross-project validation | Partially | Evidence expanded from "3 个 feature、9 个 test script" to "22 个 feature、13 个含 testing scripts" with richer categorization (sections A/B/C with specific error types, structural redundancy data, and graduation rates). The within-project sample is now substantial. However, no generalizability qualifier was added, and all data still comes from pm-work-tracker. |
| Attack 2: Success Criteria — SKILL.md deliverable coverage gap | Yes | The proposal now includes criterion [8]: "run-e2e-tests 从 `tests/e2e/` 执行，能发现并运行 `tests/e2e/<target>/*.spec.ts`（验证方法：run-e2e-tests 输出的 test list 包含 target 目录下的 spec 文件）" and criterion [7]: "graduate-tests 执行后，`tests/e2e/<target>/` 目录存在且包含从 `<feature>/` 迁移的 spec 文件，且 import 路径指向共享 helpers". These verify the behavioral effects of the SKILL.md modifications. The scope items are deliverable-focused (SKILL.md file changes), while criteria test observable outcomes — a reasonable proxy. Full marks restored for coverage. |
| Attack 3: Alternatives Analysis — rejected alternative has asymmetric depth | No | The "扩展 config.yaml" con remains "新增维护负担；routes 变更需要手动同步" — identical to iteration 3. No quantification added. |

---

## Verdict

- **Score**: 90/100
- **Target**: 90/100
- **Gap**: 0 points
- **Action**: Target reached. The proposal meets the quality bar. Remaining weaknesses (single-project evidence without generalizability claim, asymmetric alternative analysis, thin shared-helpers mitigation) are minor and do not prevent the reader from making an informed decision.

SCORE: 90/100
