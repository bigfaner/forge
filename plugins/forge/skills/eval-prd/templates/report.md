---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
scoring_mode: "{{MODE_A_or_MODE_B}}"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/1000** (target: {{TARGET}}, mode: {{MODE_A_or_MODE_B}})

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  ___     │  150      │ ✅/⚠️/❌    │
│    Three elements            │  ___/50  │          │            │
│    Goals quantified          │  ___/40  │          │            │
│    Logical consistency       │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  ___     │  200      │ ✅/⚠️/❌    │
│    Mermaid diagram exists    │  ___/70  │          │            │
│    Main path complete        │  ___/70  │          │            │
│    Decision + error branches │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3a. Functional Specs (A)     │  ___     │  200      │ ✅/⚠️/❌    │
│  OR 3b. Flow Completeness(B) │          │          │            │
│    Sub-criterion 1           │  ___/70  │          │            │
│    Sub-criterion 2           │  ___/70  │          │            │
│    Sub-criterion 3           │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  ___     │  300      │ ✅/⚠️/❌    │
│    Coverage per user type    │  ___/70  │          │            │
│    Format correct            │  ___/70  │          │            │
│    AC per story (G/W/T)      │  ___/60  │          │            │
│    AC verifiability          │  ___/100 │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  ___     │  150      │ ✅/⚠️/❌    │
│    In-scope concrete         │  ___/50  │          │            │
│    Out-of-scope explicit     │  ___/40  │          │            │
│    Consistent with specs     │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  1000     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode A** (prd-ui-functions.md present): Dimension 3 evaluates Functional Specs from prd-ui-functions.md.
> Sub-criteria: Placement & Interaction completeness /70, Data Requirements & States clarity /70, Validation Rules explicit /60.
>
> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /70, Data flow documented /70, Exception handling and edge cases /60.

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

- **Score**: {{SCORE}}/1000
- **Target**: {{TARGET}}/1000
- **Gap**: {{GAP}} points
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
