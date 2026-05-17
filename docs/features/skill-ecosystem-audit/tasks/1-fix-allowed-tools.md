---
id: "1"
title: "Fix allowed-tools field name and format across all skills/commands"
priority: "P0"
estimated_time: "45m"
dependencies: []
type: "cleanup"
scope: "all"
breaking: false
mainSession: false
---

# 1: Fix allowed-tools field name and format across all skills/commands

## Description

The official Claude Code Frontmatter reference documents the field as `allowed-tools` (hyphen), but 15 files use `allowed_tools` (underscore). Additionally, all 17 files with this field use JSON array syntax (`["Bash", "Read"]`) instead of the documented space-separated string format (`Bash Read`).

If the underscore form is not recognized by Claude Code, permission grants silently fail — tools that should be pre-approved will prompt the user each time.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal (W1, items 1-2)

## Affected Files

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/commands/clean-code.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Edit", "Write", "Glob", "Grep"]` → `Bash Read Edit Write Glob Grep` |
| `plugins/forge/commands/execute-task.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Agent", "TaskOutput", "Skill"]` → `Bash Read Agent TaskOutput Skill` |
| `plugins/forge/commands/extract-design-md.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Write", "WebFetch"]` → `Bash Read Write WebFetch` |
| `plugins/forge/commands/fix-bug.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]` → `Bash Read Write Edit Grep Glob Agent LSP` |
| `plugins/forge/commands/gen-sitemap.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Write", "Grep", "Glob"]` → `Bash Read Write Grep Glob` |
| `plugins/forge/commands/git-checkout.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read"]` → `Bash Read` |
| `plugins/forge/commands/git-commit.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read"]` → `Bash Read` |
| `plugins/forge/commands/init-forge.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read"]` → `Bash Read` |
| `plugins/forge/commands/quick.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "Skill", "AskUserQuestion"]` → `Bash Read Write Edit Grep Glob Agent Skill AskUserQuestion` |
| `plugins/forge/commands/run-tasks.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Agent", "TaskOutput", "Skill"]` → `Bash Read Agent TaskOutput Skill` |
| `plugins/forge/commands/simplify-skill.md` | `allowed_tools` → `allowed-tools`, `["Read", "Write", "Edit", "AskUserQuestion"]` → `Read Write Edit AskUserQuestion` |
| `plugins/forge/skills/clean-code/SKILL.md` | `allowed_tools` → `allowed-tools`, `["Bash", "Read", "Edit", "Write", "Glob", "Grep"]` → `Bash Read Edit Write Glob Grep` |
| `plugins/forge/skills/init-justfile/SKILL.md` | `['Bash', 'Read', 'Write', 'Edit']` → `Bash Read Write Edit` |

## Acceptance Criteria

- [ ] Zero files contain `allowed_tools` (underscore)
  `grep -rn 'allowed_tools' plugins/forge/skills/ plugins/forge/commands/` returns 0 hits
- [ ] Zero files use JSON array format for `allowed-tools`
  `grep -rn 'allowed-tools: \[' plugins/forge/skills/ plugins/forge/commands/` returns 0 hits
- [ ] All `allowed-tools` values use space-separated format
  `grep -rn 'allowed-tools:' plugins/forge/skills/ plugins/forge/commands/` shows only space-separated strings

## Hard Rules

- Only modify frontmatter `allowed_tools`/`allowed-tools` lines. Do not change any other content.
- Preserve the exact list of tools — only change the field name and format syntax.

## Implementation Notes

- `init-justfile/SKILL.md` uses single-quote array syntax `['Bash', ...]` — also convert to space-separated.
- `improve-harness/SKILL.md` and `eval/SKILL.md` do NOT have `allowed-tools` — skip them.
