---
id: "2"
title: "Inline knowledge-extraction.md into consuming skills and commands"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 2: Inline knowledge-extraction.md into consuming skills and commands

## Description
Replace `${CLAUDE_SKILL_DIR}/../../references/shared/knowledge-extraction.md` (skills) and `${CLAUDE_SKILL_DIR}/../references/shared/knowledge-extraction.md` (commands) read instructions with full inlined extraction routine. Used by write-prd, tech-design, run-tasks, and fix-bug.

## Reference Files
- `docs/proposals/remove-references-dir/proposal.md` — Source proposal
- `plugins/forge/references/shared/knowledge-extraction.md` — Content to inline

> **Note:** Line numbers are approximate and may drift. Search for `references/shared/knowledge-extraction` to locate exact reference sites.

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/SKILL.md` | Replace line 259 reference with inlined extraction flow |
| `plugins/forge/skills/tech-design/SKILL.md` | Replace line 228 reference with inlined extraction flow |
| `plugins/forge/commands/run-tasks.md` | Replace line 126 reference with inlined extraction flow |
| `plugins/forge/commands/fix-bug.md` | Replace line 212 reference with inlined extraction flow |

## Acceptance Criteria
- [ ] No occurrence of `references/shared/knowledge-extraction` in any of the 4 modified files
- [ ] Each file contains the full knowledge extraction routine verbatim
- [ ] Path references using `${CLAUDE_SKILL_DIR}` to knowledge-extraction.md are fully removed

## Hard Rules
- Inline the full extraction routine, not a summary
- Do not modify the extraction logic — structural copy-paste only

## Implementation Notes
- Skills (write-prd, tech-design) use path prefix `../../references/shared/` (2 levels up from skills/<name>/)
- Commands (run-tasks, fix-bug) use path prefix `../references/shared/` (1 level up from commands/)
- Each consumer passes different parameters to the extraction flow — preserve the parameter context around the inlined content
