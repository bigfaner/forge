---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
scoring_mode: "{{MODE_docs_or_MODE_full}}"
evaluator: Claude (automated, adversarial)
---

# Consistency Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/1000** (target: {{TARGET}}, scope: {{MODE_docs_or_MODE_full}})

<!-- Scorecard for Mode: docs (default). When scope=full, replace with the full-mode scorecard below. -->

```
┌─────────────────────────────────────────────────────────────────┐
│          CONSISTENCY QUALITY SCORECARD — Mode: docs              │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD-Design Alignment      │  ___     │  250     │ ✅/⚠️/❌    │
│    Functional coverage       │  ___/100 │          │            │
│    No orphan design          │  ___/80  │          │            │
│    Error handling alignment  │  ___/70  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. PRD-UI Consistency        │  ___     │  150     │ ✅/⚠️/❌    │
│    UI component coverage     │  ___/60  │          │            │
│    Data binding alignment    │  ___/50  │          │            │
│    Validation rules          │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Design-Task Coverage      │  ___     │  200     │ ✅/⚠️/❌    │
│    Design module coverage    │  ___/80  │          │            │
│    PRD AC coverage           │  ___/70  │          │            │
│    Task source references    │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Terminology Consistency   │  ___     │  150     │ ✅/⚠️/❌    │
│    Entity naming             │  ___/60  │          │            │
│    Field naming              │  ___/50  │          │            │
│    Status/enum consistency   │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Data Model Consistency    │  ___     │  150     │ ✅/⚠️/❌    │
│    ER matches PRD data       │  ___/60  │          │            │
│    API params match models   │  ___/50  │          │            │
│    Schema matches design     │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Traceability Completeness │  ___     │  100     │ ✅/⚠️/❌    │
│    Table completeness        │  ___/50  │          │            │
│    No orphan references      │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

<!-- Scorecard for Mode: full. Delete the docs scorecard above and use this one when scope=full.

```
┌─────────────────────────────────────────────────────────────────┐
│          CONSISTENCY QUALITY SCORECARD — Mode: full              │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD-Design Alignment      │  ___     │  150     │ ✅/⚠️/❌    │
│    Functional coverage       │  ___/60  │          │            │
│    No orphan design          │  ___/50  │          │            │
│    Error handling alignment  │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. PRD-UI Consistency        │  ___     │  100     │ ✅/⚠️/❌    │
│    UI component coverage     │  ___/40  │          │            │
│    Data binding alignment    │  ___/30  │          │            │
│    Validation rules          │  ___/30  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Design-Task Coverage      │  ___     │  150     │ ✅/⚠️/❌    │
│    Design module coverage    │  ___/60  │          │            │
│    PRD AC coverage           │  ___/50  │          │            │
│    Task source references    │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Terminology Consistency   │  ___     │  150     │ ✅/⚠️/❌    │
│    Entity naming             │  ___/60  │          │            │
│    Field naming              │  ___/50  │          │            │
│    Status/enum consistency   │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Data Model Consistency    │  ___     │  100     │ ✅/⚠️/❌    │
│    ER matches PRD data       │  ___/40  │          │            │
│    API params match models   │  ___/40  │          │            │
│    Schema matches design     │  ___/20  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Traceability Completeness │  ___     │  100     │ ✅/⚠️/❌    │
│    Table completeness        │  ___/50  │          │            │
│    No orphan references      │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 7. Interface-Code Alignment  │  ___     │  150     │ ✅/⚠️/❌    │
│    Function signatures match │  ___/60  │          │            │
│    Return types consistent   │  ___/50  │          │            │
│    Error types consistent    │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 8. Data Model-Code Alignment │  ___     │  100     │ ✅/⚠️/❌    │
│    Struct/class fields match │  ___/60  │          │            │
│    Validation constraints    │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
-->

> **N/A Dimensions**: If a dimension's source documents don't exist, it gets full points. List any N/A dimensions here.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| <!-- doc:section --> | <!-- inconsistency description --> | <!-- -N pts --> |

---

## Attack Points

### Attack 1: [dimension — specific inconsistency]

**Where**: <!-- quote from document A vs quote from document B -->
**Why it's inconsistent**: <!-- concrete description of the mismatch -->
**How to fix**: <!-- actionable fix: which document to change and what to change -->

### Attack 2: [dimension — specific inconsistency]

**Where**: <!-- quote from document A vs quote from document B -->
**Why it's inconsistent**: <!-- concrete description of the mismatch -->
**How to fix**: <!-- actionable fix: which document to change and what to change -->

### Attack 3: [dimension — specific inconsistency]

**Where**: <!-- quote from document A vs quote from document B -->
**Why it's inconsistent**: <!-- concrete description of the mismatch -->
**How to fix**: <!-- actionable fix: which document to change and what to change -->

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| <!-- attack from iter N-1 --> | ✅/❌ | <!-- what changed / what didn't --> |

---

## Fix Summary

<!-- Only for iteration > 1: summarize what the reviser changed -->

| File Changed | What Changed |
|-------------|--------------|
| <!-- path --> | <!-- summary of modifications --> |

---

## Verdict

- **Score**: {{SCORE}}/1000
- **Target**: {{TARGET}}/1000
- **Gap**: {{GAP}} points
- **Inconsistencies Found**: {{COUNT}}
- **Inconsistencies Fixed**: {{FIXED_COUNT}}
- **Action**: {{Fix applied, re-scoring / Target reached / Iterations exhausted}}
