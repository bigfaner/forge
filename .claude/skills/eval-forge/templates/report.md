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
│                  PLUGIN CONSISTENCY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Directory-Name Alignment  │  ___     │  40      │ ✅/⚠️/❌    │
│    Skill name matches dir    │  ___/25  │          │            │
│    Command name matches file │  ___/15  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Agent Reference Integrity │  ___     │  100     │ ✅/⚠️/❌    │
│    Referenced agents exist   │  ___/70  │          │            │
│    No orphan agents          │  ___/30  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Reference Integrity       │  ___     │  80      │ ✅/⚠️/❌    │
│    Template refs valid       │  ___/25  │          │            │
│    Cross-skill refs valid    │  ___/30  │          │            │
│    No orphan templates       │  ___/25  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Frontmatter Completeness  │  ___     │  110     │ ✅/⚠️/❌    │
│    Skill frontmatter         │  ___/45  │          │            │
│    Command frontmatter       │  ___/35  │          │            │
│    Agent frontmatter         │  ___/30  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Eval Template Convention  │  ___     │  100     │ ✅/⚠️/❌    │
│    rubric.md exists          │  ___/30  │          │            │
│    report.md exists          │  ___/30  │          │            │
│    Rubric→report chain valid │  ___/20  │          │            │
│    Rubric totals correct     │  ___/20  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Orchestrator Convention   │  ___     │  40      │ ✅/⚠️/❌    │
│    Iron Laws present         │  ___/25  │          │            │
│    Hard Gate present         │  ___/15  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 7. Task CLI Alignment        │  ___     │  240     │ ✅/⚠️/❌    │
│    Command existence         │  ___/25  │          │            │
│    Flag correctness          │  ___/25  │          │            │
│    Output field parsing      │  ___/15  │          │            │
│    Status machine align      │  ___/35  │          │            │
│    Claim scheduling align    │  ___/35  │          │            │
│    Record validation align   │  ___/35  │          │            │
│    Dynamic task add align    │  ___/25  │          │            │
│    Schema-code alignment     │  ___/20  │          │            │
│    All-completed hook align  │  ___/10  │          │            │
│    Template existence        │  ___/10  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 8. Hook Wiring Integrity     │  ___     │  70      │ ✅/⚠️/❌    │
│    hooks.json valid JSON     │  ___/10  │          │            │
│    Hook scripts exist        │  ___/25  │          │            │
│    Hook CLI commands valid   │  ___/15  │          │            │
│    Hook event names valid    │  ___/20  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 9. Guide Coverage            │  ___     │  70      │ ✅/⚠️/❌    │
│    Guide references valid    │  ___/30  │          │            │
│    Core skills documented    │  ___/40  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 10. Command Metadata         │  ___     │  60      │ ✅/⚠️/❌    │
│    allowed_tools declared    │  ___/35  │          │            │
│    argument-hints declared   │  ___/25  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 11. Plugin Metadata          │  ___     │  40      │ ✅/⚠️/❌    │
│    keywords coverage         │  ___/25  │          │            │
│    description accurate      │  ___/15  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 12. Safety Marker Consist.   │  ___     │  50      │ ✅/⚠️/❌    │
│    Command/agent markers     │  ___/30  │          │            │
│    Dispatch cmd coverage     │  ___/20  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  1000    │            │
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
