---
id: "2"
title: "Rewrite scorer-prompt.md: 4-phase adversarial scoring methodology"
priority: "P0"
estimated_time: "1-2h"
dependencies: [1]
type: "documentation"
mainSession: false
---

# 2: Rewrite scorer-prompt.md: 4-phase adversarial scoring methodology

## Description
Replace the current flat checklist scorer prompt with a 4-phase adversarial process: (1) build workflow graph, (2) per-node adversarial testing, (3) per-file precision review, (4) basic integrity. The current scorer does plane checks and scores 965/1000 while missing all runtime issues. The new scorer must construct a workflow graph, identify bypass vectors, find instruction conflicts, and detect redundancies.

## Reference Files
- `docs/proposals/eval-forge-runtime-audit/proposal.md` — Source proposal (Section 2: Scorer Methodology)
- `.claude/skills/eval-forge/templates/scorer-prompt.md` — Current scorer prompt (to be rewritten)
- `.claude/skills/eval-forge/templates/rubric.md` — New rubric (rewritten by Task 1)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `.claude/skills/eval-forge/templates/scorer-prompt.md` | Complete rewrite: 4-phase methodology, new output format |

## Acceptance Criteria
- [ ] Scorer prompt implements 4-phase process: (1) Build workflow graph → (2) Per-node adversarial testing → (3) Per-file precision review → (4) Basic integrity
- [ ] Phase 1 instructs scorer to read rubric's embedded workflow specs, scan actual skills/commands/agents, compare specs vs actual, find breakpoints/dead-ends/unreachable states
- [ ] Phase 2 instructs scorer to list every gate/confirm point, assume "lazy agent" perspective, check HARD-RULE enforcement, eval loop independence, quality gate CLI enforcement
- [ ] Phase 3 instructs scorer to check instruction conflicts first (highest priority), then step ambiguity, then incomplete conditionals, then undefined variables, then content redundancy (A/B/C categories)
- [ ] Phase 4 instructs scorer to check reference integrity, frontmatter, eval templates, name alignment
- [ ] Output format updated to 6-dimension scorecard (matching new report.md from Task 4)
- [ ] Scorer reads rubric at `.claude/skills/eval-forge/templates/rubric.md` and report template at `.claude/skills/eval-forge/templates/report.md`
- [ ] Input section specifies all files to scan: `plugins/forge/skills/*/SKILL.md`, `plugins/forge/commands/*.md`, `plugins/forge/agents/*.md`, `hooks/hooks.json`, `hooks/guide.md`, task CLI source code
- [ ] EXTREMELY-IMPORTANT block preserved: adversarial stance, exact file paths, no full marks unless perfect
- [ ] Dimension 7 task CLI source code reading section preserved and adapted (still reads Go source for D5 reference integrity)

## Hard Rules
- Scorer MUST read the rubric as its first action — dimensions and criteria come from rubric, not hardcoded in the prompt
- The 4 phases must be sequential (Phase N+1 depends on Phase N findings)
- Scorer output format must match the new report.md template structure

## Implementation Notes
- The proposal's Section 2 defines the 4-phase process with specific steps — use it as the methodology spec
- Dimension 7 (Task CLI Alignment) in the old scorer was the closest to runtime checking — its source code reading should be preserved for the new D5 (Reference Integrity) and D1 (Workflow Completeness) checks
- The scorer should output structured ATTACKS section with dimension-specific prefixes (D1, D2, D3, D4, D5, D6)
