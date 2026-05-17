---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
scoring_mode: "{{MODE_docs_or_MODE_full}}"
evaluator: Claude (automated, adversarial)
---

# Consistency Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}}, scope: {{MODE_docs_or_MODE_full}})

```
┌─────────────────────────────────────────────────────────────────┐
│                  CONSISTENCY QUALITY SCORECARD                    │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD-Design Alignment      │  ___     │  25      │ ✅/⚠️/❌    │
│    Functional coverage       │  ___/10  │          │            │
│    No orphan design          │  ___/8   │          │            │
│    Error handling alignment  │  ___/7   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. PRD-UI Consistency        │  ___     │  15      │ ✅/⚠️/❌    │
│    UI component coverage     │  ___/6   │          │            │
│    Data binding alignment    │  ___/5   │          │            │
│    Validation rules          │  ___/4   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Design-Task Coverage      │  ___     │  20      │ ✅/⚠️/❌    │
│    Design module coverage    │  ___/8   │          │            │
│    PRD AC coverage           │  ___/7   │          │            │
│    Task source references    │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Terminology Consistency   │  ___     │  15      │ ✅/⚠️/❌    │
│    Entity naming             │  ___/6   │          │            │
│    Field naming              │  ___/5   │          │            │
│    Status/enum consistency   │  ___/4   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Data Model Consistency    │  ___     │  15      │ ✅/⚠️/❌    │
│    ER matches PRD data       │  ___/6   │          │            │
│    API params match models   │  ___/5   │          │            │
│    Schema matches design     │  ___/4   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Traceability Completeness │  ___     │  10      │ ✅/⚠️/❌    │
│    Table completeness        │  ___/5   │          │            │
│    No orphan references      │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

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

- **Score**: {{SCORE}}/100
- **Target**: {{TARGET}}/100
- **Gap**: {{GAP}} points
- **Inconsistencies Found**: {{COUNT}}
- **Inconsistencies Fixed**: {{FIXED_COUNT}}
- **Action**: {{Fix applied, re-scoring / Target reached / Iterations exhausted}}
