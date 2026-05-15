---
id: "1"
title: "Rewrite rubric.md: 6-dimension runtime reliability scoring"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Rewrite rubric.md: 6-dimension runtime reliability scoring

## Description
Replace the current 12-dimension structural consistency rubric (1000 pts) with a 6-dimension runtime reliability rubric (1000 pts). The current rubric scores 965/1000 but misses 39 runtime issues found by manual audit. The new rubric shifts weight from file-structure checks (290 pts → 50 pts) to runtime reliability (700 pts).

## Reference Files
- `docs/proposals/eval-forge-runtime-audit/proposal.md` — Source proposal (full dimension specs with scoring criteria)
- `.claude/skills/eval-forge/templates/rubric.md` — Current rubric (to be rewritten)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `.claude/skills/eval-forge/templates/rubric.md` | Complete rewrite: 6 dimensions, new scoring criteria, 4-phase methodology reference |

## Acceptance Criteria
- [ ] Rubric defines exactly 6 dimensions: Workflow Completeness (250), Bypass Resistance (250), Instruction Precision (200), Cross-file Dedup (150), Reference Integrity (100), Structural Convention (50)
- [ ] Each dimension has detailed scoring criteria matching proposal tables (1a-1e, 2a-2e, 3a-3d, 4a-4c, 5a-5d, 6a-6c)
- [ ] Dimension 1 embeds the ground-truth workflow specs (Full Mode, Quick Mode, Manifest Status Machine, Per-Skill Precondition/Output Matrix)
- [ ] Dimension 2 lists the 5 bypass types with point allocations and the 14 known bypass vectors (BV-2.1 through BV-5.2)
- [ ] Dimension 3 lists CLI-filled variables from `prompt.go` (not marked as "undefined")
- [ ] Dimension 4 references known redundancy instances (3 categories: content copy, guide overlap, unreasonable inline)
- [ ] Deduction tiers updated to reflect new severity model (breakpoints = -20/-15/-10/-5 instead of flat -5/-15/-25)
- [ ] Point totals sum to 1000
- [ ] Report template reference updated to `.claude/skills/eval-forge/templates/report.md`

## Hard Rules
- Do NOT change the report template file itself (Task 4 handles that). Only reference it.
- Point totals per dimension must exactly match: 250+250+200+150+100+50 = 1000
- The Per-Skill Precondition/Output Matrix in Dimension 1 must include ALL skills listed in the proposal

## Implementation Notes
- The proposal provides the complete spec for each dimension's scoring criteria — use it as the authoritative source
- Known redundancy instances in D4 are explicitly listed in the proposal (5 instances) — include them verbatim
- Known bypass vectors in D2 are explicitly listed (14 vectors) — include them as a "known vectors to verify" table
