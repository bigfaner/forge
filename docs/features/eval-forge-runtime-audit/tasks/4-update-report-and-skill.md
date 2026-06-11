---
id: "4"
title: "Update report.md scorecard and SKILL.md dimension table"
priority: "P1"
estimated_time: "30m"
dependencies: [1, 2, 3]
type: "documentation"
mainSession: false
---

# 4: Update report.md scorecard and SKILL.md dimension table

## Description
Update the report template scorecard from 12 dimensions to 6 dimensions, and update SKILL.md's Final Report section to match. These are structural changes that depend on the rubric, scorer, and reviser being finalized.

## Reference Files
- `docs/proposals/eval-forge-runtime-audit/proposal.md` — Source proposal
- `.claude/skills/eval-forge/templates/report.md` — Current report template (12-dimension scorecard)
- `.claude/skills/eval-forge/SKILL.md` — Current skill definition (12-dimension Final Report)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `.claude/skills/eval-forge/templates/report.md` | Replace 12-dimension scorecard with 6-dimension scorecard |
| `.claude/skills/eval-forge/SKILL.md` | Update Final Report dimension table from 12 to 6 dimensions |

## Acceptance Criteria
- [ ] report.md scorecard has exactly 6 dimensions: Workflow Completeness (250), Bypass Resistance (250), Instruction Precision (200), Cross-file Dedup (150), Reference Integrity (100), Structural Convention (50)
- [ ] Each dimension in scorecard shows sub-criteria rows matching the rubric (e.g., D1: 1a/80, 1b/40, 1c/50, 1d/30, 1e/50)
- [ ] Score total line shows ___/1000
- [ ] SKILL.md Final Report section (Step 6) dimension table shows 6 rows matching new dimensions with correct max scores
- [ ] SKILL.md Parameters section unchanged (target: 950, iterations: 3)
- [ ] SKILL.md Architecture diagram unchanged (still 6-step loop)
- [ ] SKILL.md Steps 1-5 unchanged (scan, score, gate, classify, fix)
- [ ] Frontmatter preserved (name: eval-forge, description unchanged)

## Hard Rules
- Do NOT modify Steps 1-5 logic in SKILL.md — only update the Final Report dimension table
- Do NOT modify the Architecture diagram or Orchestrator Iron Laws
- Parameters must remain target: 950, iterations: 3

## Implementation Notes
- Straightforward structural update — the rubric (Task 1) defines the dimension names and scores, just mirror them into report.md and SKILL.md
- report.md sub-criteria rows should exactly match rubric's per-dimension criteria tables
