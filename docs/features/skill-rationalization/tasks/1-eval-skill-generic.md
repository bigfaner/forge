---
id: "1"
title: "Create generic eval skill and extract rubric files"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Create generic eval skill and extract rubric files

## Description
Create the single generic `skills/eval/SKILL.md` that replaces all 7 eval skills, and move existing rubric templates into `skills/eval/rubrics/`.

The generic skill implements the scorer→gate→revise loop that is currently duplicated across eval-proposal, eval-prd, eval-design, eval-ui, eval-test-cases, eval-consistency, and eval-harness. It parameterizes on: rubric path, target score, max iterations, and scale (100-point vs 1000-point).

## Reference Files
- `docs/proposals/skill-rationalization/proposal.md` — Source proposal
- `plugins/forge/skills/eval-proposal/SKILL.md` — Reference for scorer→gate→revise loop
- `plugins/forge/skills/eval-harness/SKILL.md` — Reference for 100-point scale variant (no reviser)
- `plugins/forge/skills/eval-ui/SKILL.md` — Reference for multi-platform rubric selection
- `plugins/forge/skills/eval-*/templates/rubric*.md` — Existing rubric templates to migrate

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/SKILL.md` | Generic eval skill with scorer→gate→revise loop |
| `plugins/forge/skills/eval/rubrics/proposal.md` | Migrated from eval-proposal/templates/rubric.md |
| `plugins/forge/skills/eval/rubrics/prd.md` | Migrated from eval-prd/templates/rubric.md |
| `plugins/forge/skills/eval/rubrics/design.md` | Migrated from eval-design/templates/rubric.md |
| `plugins/forge/skills/eval/rubrics/ui-web.md` | Migrated from eval-ui/templates/rubric-web.md |
| `plugins/forge/skills/eval/rubrics/ui-mobile.md` | Migrated from eval-ui/templates/rubric-mobile.md |
| `plugins/forge/skills/eval/rubrics/ui-tui.md` | Migrated from eval-ui/templates/rubric-tui.md |
| `plugins/forge/skills/eval/rubrics/test-cases.md` | Migrated from eval-test-cases/templates/rubric.md |
| `plugins/forge/skills/eval/rubrics/consistency.md` | Migrated from eval-consistency/templates/rubric.md |
| `plugins/forge/skills/eval/rubrics/harness.md` | Migrated from eval-harness/templates/rubric.md |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none — old files deleted in Task 3) | |

## Acceptance Criteria
- [ ] `skills/eval/SKILL.md` contains a single generic scorer→gate→revise loop that reads rubric from `rubrics/<type>.md`
- [ ] Each rubric file is self-contained with frontmatter declaring `scale`, `target`, `iterations`, and `type`
- [ ] `eval-harness` rubric declares `scale: 100` (not 1000) — the generic skill must handle both
- [ ] `eval-ui` rubric resolution: generic skill detects UI platform from manifest/config and selects ui-web, ui-mobile, or ui-tui
- [ ] No eval-specific orchestration logic outside the generic `skills/eval/SKILL.md`

## Hard Rules
- Rubric frontmatter must include: `scale` (100 or 1000), `target` (default score threshold), `iterations` (max reviser rounds), `type` (eval type identifier)
- The generic skill must use `doc-scorer` and `doc-reviser` subagents exactly as current eval skills do
- `eval-harness` uses only `doc-scorer` (no reviser loop) — the generic skill must detect `iterations: 0` or `iterations: 1` and skip the reviser phase

## Implementation Notes
- Read `eval-proposal/SKILL.md` as the canonical reference for the scorer→gate→revise loop structure
- Read `eval-harness/SKILL.md` for the 100-point scale variant — note it does NOT use doc-reviser
- Read `eval-ui/SKILL.md` for multi-platform rubric selection logic
- Each rubric file should include a `target` and `iterations` override in frontmatter (e.g., `eval-test-cases` defaults to 900 target / 6 iterations per the forge guide)
- The rubric content itself is copied verbatim from existing templates — no content changes
