---
id: "4"
title: "Remove frontmatter duplication and rule preview redundancy"
priority: "P1"
estimated_time: "1h"
dependencies: [3]
type: "doc"
mainSession: false
---

# 4: Remove frontmatter duplication and rule preview redundancy

## Description

Remove (1) frontmatter description repeating as body opening sentence, (2) skill body paragraphs previewing rule file content instead of referencing them.

## Reference Files
- `docs/proposals/skill-instruction-audit/proposal.md#Success-Criteria`: SC-6
- `plugins/forge/commands/clean-code.md`: frontmatter == body line 7
- `plugins/forge/commands/git-commit.md`: frontmatter == body line 8
- `plugins/forge/commands/git-checkout.md`: frontmatter ≈ body line 8
- `plugins/forge/commands/init-forge.md`: frontmatter == body line 9
- `plugins/forge/skills/learn/SKILL.md`: frontmatter ≈ body line 8-9
- `plugins/forge/skills/gen-contracts/SKILL.md`: frontmatter ≈ Core principle
- `plugins/forge/skills/run-tests/SKILL.md`: Opening ≈ Core principle
- `plugins/forge/skills/extract-design-md/SKILL.md`: Mobile/TUI Overviews duplicate rules/

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/clean-code.md` | Delete opening sentence duplicating frontmatter |
| `plugins/forge/commands/git-commit.md` | Delete opening sentence duplicating frontmatter |
| `plugins/forge/commands/git-checkout.md` | Delete opening sentence duplicating frontmatter |
| `plugins/forge/commands/init-forge.md` | Delete opening sentence duplicating frontmatter |
| `plugins/forge/skills/learn/SKILL.md` | Delete opening paragraph duplicating frontmatter |
| `plugins/forge/skills/gen-contracts/SKILL.md` | Shorten description; Core principle only "Every Step gets a Contract" |
| `plugins/forge/skills/run-tests/SKILL.md` | Merge opening into Core principle |
| `plugins/forge/skills/extract-design-md/SKILL.md` | Delete Mobile/TUI Overviews; keep reference to rules/platform-routing.md |

## Acceptance Criteria

- [ ] 4 commands (clean-code, git-commit, git-checkout, init-forge) body doesn't start with frontmatter sentence
- [ ] `learn/SKILL.md` body doesn't duplicate frontmatter
- [ ] `gen-contracts/SKILL.md` description ≤1 sentence
- [ ] `run-tests/SKILL.md` has single merged opening
- [ ] `extract-design-md/SKILL.md` has no Overview paragraphs; only rules reference

## Hard Rules

- 仅修改上述 8 个文件
