---
id: "4"
title: "Update forge guide and references for new eval structure"
priority: "P1"
estimated_time: "1h"
dependencies: ["3"]
type: "documentation"
mainSession: false
---

# 4: Update forge guide and references for new eval structure

## Description
Update the forge guide (system prompt section in CLAUDE.md), the eval-forge audit skill, and any other references that mention the old eval skill structure.

## Reference Files
- `docs/proposals/skill-rationalization/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/hooks.json` | Update eval skill references if listed |
| `forge-cli/pkg/prompt/data/*.md` | Update CLI prompt templates that reference old eval skill paths |
| `.claude/skills/eval-forge/SKILL.md` | Update audit to validate new eval structure |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] Forge guide mentions generic `eval` skill instead of 7 separate eval skills
- [ ] `eval-forge` audit skill correctly validates the new eval structure (skills/eval/ + rubrics/ + command wrappers)
- [ ] CLI prompt templates reference correct paths for the new eval skill
- [ ] Skill count in documentation reflects 17 (down from 24)

## Hard Rules
- Only update references to eval skills that were consolidated — do not modify unrelated sections
- Preserve all non-eval content in modified files

## Implementation Notes
- The forge guide in CLAUDE.md system prompt lists skill names — update the eval entries
- Check `.claude/skills/eval-forge/SKILL.md` for hardcoded references to eval-proposal, eval-prd, etc.
- Check `forge-cli/pkg/prompt/data/` for prompt templates that embed eval skill paths
