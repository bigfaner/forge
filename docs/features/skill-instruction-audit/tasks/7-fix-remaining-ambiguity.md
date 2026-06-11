---
id: "7"
title: "Fix ambiguity and logic issues in remaining skills/commands"
priority: "P2"
estimated_time: "1.5h"
dependencies: [6]
type: "doc"
mainSession: false
---

# 7: Fix ambiguity and logic issues in remaining skills/commands

## Description

Fix remaining ambiguous instructions and logic issues in non-pipeline skills/commands.

## Reference Files
- `plugins/forge/skills/consolidate-specs/SKILL.md`: Steps 9-11 scope unclear, non-interactive safety assertion
- `plugins/forge/skills/gen-sitemap/SKILL.md`: @latest vs pinning contradiction, "this command" reference ambiguous
- `plugins/forge/skills/clean-code/SKILL.md`: "unless configs" exception vague, default branch detection offline
- `plugins/forge/skills/deep-research/SKILL.md`: Q4 dimension selection unclear, missing "wait for review" step
- `plugins/forge/skills/ui-design/SKILL.md`: web+mobile output collision, eval-skip label unclear
- `plugins/forge/commands/simplify-skill.md`: Paths use .claude/ instead of plugin paths
- `plugins/forge/skills/forensic/SKILL.md`: project-hash unresolvable, "skills parent directory" unclear

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Specify scan directories for Steps 9-11; qualify non-interactive safety assertion |
| `plugins/forge/skills/gen-sitemap/SKILL.md` | Unify @latest recommendation; clarify "this command" references |
| `plugins/forge/skills/clean-code/SKILL.md` | Clarify "configs" exception; add local default branch detection |
| `plugins/forge/skills/deep-research/SKILL.md` | Clarify Q4 dimension selection; add explicit review wait step |
| `plugins/forge/skills/ui-design/SKILL.md` | Add web+mobile output rule; improve eval-skip option labels |
| `plugins/forge/commands/simplify-skill.md` | Clarify scope (user skills vs plugin skills) |
| `plugins/forge/skills/forensic/SKILL.md` | Explain project-hash source; clarify path resolution |

## Acceptance Criteria

- [ ] `consolidate-specs` Steps 9-11 specify scan directories (docs/business-rules/, docs/conventions/)
- [ ] `gen-sitemap` has no @latest vs pinning contradiction
- [ ] `clean-code` has clear config file exception or none
- [ ] `deep-research` has explicit "wait for user review" between report presentation and proposal conversion ask
- [ ] `ui-design` has web+mobile output rule; eval-skip option is clearly labeled
- [ ] `simplify-skill` states whether it targets user skills or plugin skills
- [ ] `forensic` explains how to obtain project-hash path

## Hard Rules

- 仅修改上述 7 个文件
