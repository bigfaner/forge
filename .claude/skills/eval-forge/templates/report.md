---
date: "{{DATE}}"
plugin_version: "{{VERSION}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (structural audit)
---

# Forge Plugin Audit — Iteration {{ITERATION}}

**Score: {{SCORE}}/1000** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│               RUNTIME RELIABILITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Workflow Completeness     │  ___     │  250     │ ✅/⚠️/❌    │
│    1a. Full mode chain       │  ___/80  │          │            │
│    1b. Quick mode chain      │  ___/40  │          │            │
│    1c. Conditional branching │  ___/50  │          │            │
│    1d. Manifest status       │  ___/30  │          │            │
│    1e. Test lifecycle        │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Bypass Resistance         │  ___     │  250     │ ✅/⚠️/❌    │
│    2a. Quality gates         │  ___/70  │          │            │
│    2b. Eval integrity        │  ___/70  │          │            │
│    2c. User interaction      │  ___/45  │          │            │
│    2d. Required steps        │  ___/35  │          │            │
│    2e. Prohibition           │  ___/30  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Instruction Precision     │  ___     │  200     │ ✅/⚠️/❌    │
│    3a. Instruction conflicts │  ___/80  │          │            │
│    3b. Step ambiguity        │  ___/50  │          │            │
│    3c. Incomplete condition. │  ___/40  │          │            │
│    3d. Variable clarity      │  ___/30  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Cross-file Dedup          │  ___     │  150     │ ✅/⚠️/❌    │
│    4a. Content copy          │  ___/60  │          │            │
│    4b. guide.md overlap      │  ___/50  │          │            │
│    4c. Unreasonable inline   │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Reference Integrity       │  ___     │  100     │ ✅/⚠️/❌    │
│    5a. Agent refs            │  ___/30  │          │            │
│    5b. Template refs         │  ___/25  │          │            │
│    5c. Cross-skill refs      │  ___/25  │          │            │
│    5d. Hook refs             │  ___/20  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Structural Convention     │  ___     │  50      │ ✅/⚠️/❌    │
│    6a. Frontmatter           │  ___/25  │          │            │
│    6b. Eval templates        │  ___/15  │          │            │
│    6c. Name alignment        │  ___/10  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___/1000│          │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| # | Check | File | Issue | Penalty |
|---|-------|------|-------|---------|
| | <!-- check name --> | <!-- path --> | <!-- description --> | <!-- -N --> |

---

## Attack Points

### Attack 1: [dimension — specific issue]

**Where**: <!-- file path and location -->
**What's wrong**: <!-- concrete description -->
**How to fix**: <!-- actionable fix -->

### Attack 2: [dimension — specific issue]

**Where**: <!-- file path and location -->
**What's wrong**: <!-- concrete description -->
**How to fix**: <!-- actionable fix -->

### Attack 3: [dimension — specific issue]

**Where**: <!-- file path and location -->
**What's wrong**: <!-- concrete description -->
**How to fix**: <!-- actionable fix -->

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| <!-- attack from iter N-1 --> | ✅/❌ | <!-- what changed --> |

---

## Fix Summary

<!-- Only for iteration > 1 -->

| File Changed | What Changed |
|-------------|--------------|
| <!-- path --> | <!-- summary --> |

---

## Verdict

- **Score**: {{SCORE}}/1000
- **Target**: {{TARGET}}/1000
- **Gap**: {{GAP}} points
- **Action**: {{Fix applied, re-auditing / Target reached / Iterations exhausted}}
