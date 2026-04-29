---
date: "2026-04-28"
doc_dir: "docs/proposals/e2e-script-quality/"
iteration: 2
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 83/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Problem Definition        │  17      │  20      │ ⚠️         │
│    Problem clarity           │  6/7     │          │            │
│    Evidence provided         │  6/7     │          │            │
│    Urgency justified         │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Solution Clarity          │  18      │  20      │ ⚠️         │
│    Approach concrete         │  7/7     │          │            │
│    User-facing behavior      │  6/7     │          │            │
│    Differentiated            │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Alternatives Analysis     │  13      │  15      │ ✅         │
│    Alternatives listed (≥2)  │  5/5     │          │            │
│    Pros/cons honest          │  4/5     │          │            │
│    Rationale justified       │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Scope Definition          │  13      │  15      │ ✅         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  4/5     │          │            │
│    Scope bounded             │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Risk Assessment           │  14      │  15      │ ✅         │
│    Risks identified (≥3)     │  5/5     │          │            │
│    Likelihood + impact rated │  4/5     │          │            │
│    Mitigations actionable    │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Success Criteria          │  12      │  15      │ ⚠️         │
│    Measurable                │  4/5     │          │            │
│    Coverage complete         │  4/5     │          │            │
│    Testable                  │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  87      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘

Deduction adjustments: 87 - 3 (scope/success inconsistency) - 1 (vague language in alternative) = 83
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Alternatives table, Selected row | "改动范围较大" — "较大" is unquantified vague language for the chosen approach's con. No file count, no complexity metric. | -1 pt |
| Scope vs Success mismatch | Scope item "gen-test-scripts 模板（playwright-ui.spec.ts, api.spec.ts）：增加 VERIFY 标记" requires modifying templates to contain VERIFY markers, but success criterion [3] checks for absence of VERIFY in output. No criterion verifies the templates themselves contain the VERIFY scaffolding. | -3 pts |

---

## Attack Points

### Attack 1: Success Criteria — template deliverable has no acceptance test

**Where**: Scope lists "gen-test-scripts 模板（playwright-ui.spec.ts, api.spec.ts）：增加 VERIFY 标记" and "helpers.ts 模板：增加 VERIFY 标记", but the only VERIFY-related success criterion is "生成的脚本中 `grep -r '// VERIFY:' tests/e2e/` 返回 0 行"
**Why it's weak**: There is a logical gap. The scope says templates must *contain* VERIFY markers (so the agent knows where to fill in values). The success criterion checks that generated *output* has no residual VERIFY markers. But nothing verifies that the templates were actually modified to include VERIFY markers in the first place. If someone forgets to add VERIFY markers to the templates, the success criterion still passes — it just passes vacuously because no markers were ever present. This is a coverage hole that undermines the entire VERIFY mechanism.
**What must improve**: Add a success criterion such as: "Template files api.spec.ts, playwright-ui.spec.ts, and helpers.ts each contain at least one `// VERIFY:` comment (verified by `grep -c '// VERIFY:' <template-path>` >= 1 for each)."

### Attack 2: Problem Definition — evidence is single-project, no frequency data

**Where**: Evidence section: "实际项目中（pm-work-tracker）观察到：" followed by 7 bullet points
**Why it's weak**: All seven evidence points come from a single project (pm-work-tracker). This is an N=1 sample. The proposal does not state how many projects were tested, whether these errors appear consistently across projects, or what percentage of generated scripts are affected. Are these 7 errors out of 10 generated features (70% error rate) or 7 errors out of 50 features (14%)? Without frequency data, the reader cannot assess whether this is a systemic problem affecting most generation runs or a one-off issue with a single project's unusual structure.
**What must improve**: Add frequency data: "Across N projects/features, X% of generated scripts contained at least one inference error." Or at minimum: "In pm-work-tracker (M features, N test scripts), all 7 of these error categories were observed." Give the reader a denominator.

### Attack 3: Alternatives Analysis — selected approach's con is a hand-wave, not a real trade-off

**Where**: Alternatives table, selected row: Cons = "改动范围较大"
**Why it's weak**: The single con for the chosen approach is "改动范围较大" (change scope is large). This is not a trade-off analysis — it is an admission with no substance. How large? The scope section lists 8 deliverables. Is 8 files "较大"? Compared to the "only change gen-test-scripts" alternative (presumably 2-3 files), how much larger? Without quantifying the con, the reader cannot weigh whether the benefit is worth the cost. The rejected alternatives have specific, concrete cons ("test-cases.md 的错误仍然传递", "gen-test-scripts 仍会推断端口/行为细节") but the chosen approach gets a vague "it's a lot of changes."
**What must improve**: Replace "改动范围较大" with a concrete con: e.g., "Modifies 8 files across 4 skills (gen-test-cases, gen-test-scripts, run-e2e-tests, graduate-tests), requiring coordinated testing of all dependent skills in a single iteration." This lets the reader judge whether the coordination cost is acceptable.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Success criteria aspirational, not measurable | Partially | Success criteria now include 7 items with explicit verification commands (`grep -r`, `find ... wc -l`, `grep -f`). However, criterion [1] still has an ambiguity (does a WARNING count as "每个 Route 字段标注 ✅ 或 ⚠️"? — yes, but the pass condition is unclear when routes are genuinely missing). |
| Attack 2: Risk Assessment — Impact column is a mitigation | Yes | Impact column now contains actual consequence descriptions: e.g., "test cases/scripts proceed with unvalidated routes, downstream scripts inherit incorrect paths and fail at runtime." Mitigations are in a separate column with actionable steps. |
| Attack 3: Solution Clarity — user-facing behavior underspecified | Yes | Added comprehensive Before/After examples section showing: test-cases.md validation output format, Code Reconnaissance fact table, generated script with source-annotated values, and directory structure migration diagram. |

---

## Verdict

- **Score**: 83/100
- **Target**: 90/100
- **Gap**: 7 points
- **Action**: Continue to iteration 3. Priority fixes: (1) add success criterion verifying templates contain VERIFY markers (closes scope/success coverage gap, recovers -3 penalty), (2) quantify alternative cons concretely (recovers -1 penalty + alternatives score), (3) add frequency data to evidence section.
