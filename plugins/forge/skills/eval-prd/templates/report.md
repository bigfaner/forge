---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
scoring_mode: "{{MODE_A_or_MODE_B}}"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}}, mode: {{MODE_A_or_MODE_B}})

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  ___     │  15      │ ✅/⚠️/❌    │
│    Three elements            │  ___/5   │          │            │
│    Goals quantified          │  ___/4   │          │            │
│    Logical consistency       │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  ___     │  20      │ ✅/⚠️/❌    │
│    Mermaid diagram exists    │  ___/7   │          │            │
│    Main path complete        │  ___/7   │          │            │
│    Decision + error branches │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3a. Functional Specs (A)     │  ___     │  20      │ ✅/⚠️/❌    │
│  OR 3b. Flow Completeness(B) │          │          │            │
│    Sub-criterion 1           │  ___/7   │          │            │
│    Sub-criterion 2           │  ___/7   │          │            │
│    Sub-criterion 3           │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  ___     │  30      │ ✅/⚠️/❌    │
│    Coverage per user type    │  ___/7   │          │            │
│    Format correct            │  ___/7   │          │            │
│    AC per story (G/W/T)      │  ___/6   │          │            │
│    AC verifiability          │  ___/10  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  ___     │  15      │ ✅/⚠️/❌    │
│    In-scope concrete         │  ___/5   │          │            │
│    Out-of-scope explicit     │  ___/4   │          │            │
│    Consistent with specs     │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode A** (prd-ui-functions.md present): Dimension 3 evaluates Functional Specs from prd-ui-functions.md.
> Sub-criteria: Placement & Interaction completeness /7, Data Requirements & States clarity /7, Validation Rules explicit /6.
>
> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /7, Data flow documented /7, Exception handling and edge cases /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| <!-- section:line --> | <!-- vague language / missing content / cross-file inconsistency --> | <!-- -N pts --> |

---

## Attack Points

### Attack 1: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 2: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 3: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| <!-- attack from iter N-1 --> | ✅/❌ | <!-- what changed / what didn't --> |

---

## Verdict

- **Score**: {{SCORE}}/100
- **Target**: {{TARGET}}/100
- **Gap**: {{GAP}} points
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
