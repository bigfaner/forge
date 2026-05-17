---
id: "3"
title: "Fix forensic hardcoded paths and stale record-task reference"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "cleanup"
scope: "all"
breaking: false
mainSession: false
---

# 3: Fix forensic hardcoded paths and stale record-task reference

## Description

Forensic SKILL.md contains 4 hardcoded developer-specific paths (`~/.claude/projects/-Users-fanhuifeng-...`) that only work for the original developer. Additionally, `breakdown-tasks/templates/consolidate-specs.md` references `/record-task` which was replaced by `/submit-task`.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal (W2)

## Affected Files

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/skills/forensic/SKILL.md` | Replace 4 hardcoded paths at lines 89, 93, 100 with generic instructions using `${CLAUDE_SESSION_ID}` or `forge forensic` CLI patterns |
| `plugins/forge/skills/breakdown-tasks/templates/consolidate-specs.md` | Line 88: change `/record-task` to `/submit-task` |

## Acceptance Criteria

- [ ] Zero hardcoded user-specific paths in forensic
  `grep -rn '~/.claude/projects/' plugins/forge/skills/forensic/SKILL.md` returns 0 hits
- [ ] Zero `record-task` references (excluding submit-task)
  `grep -rn 'record-task' plugins/forge/ | grep -v 'submit-task'` returns 0 hits
- [ ] Forensic SKILL.md uses generic path patterns or CLI commands that work for any user

## Hard Rules

- Do not remove the forensic functionality — only replace hardcoded paths with generic equivalents.
- Use `${CLAUDE_SESSION_ID}` variable for session-specific paths (documented in Claude Code frontmatter reference).

## Implementation Notes

- Lines 89, 93, 100 in forensic/SKILL.md contain paths like `~/.claude/projects/-Users-fanhuifeng-Projects-ai-coding-harness-forge/<sessionId>.jsonl`. Replace with pattern `~/.claude/projects/<project-hash>/${CLAUDE_SESSION_ID}.jsonl` or instruct Claude to discover the path dynamically.
- The `record-task` → `submit-task` change is a simple text replacement in a template file.
