---
id: "2"
title: "Fix forensic hardcoded developer paths"
priority: "P0"
estimated_time: "1h"
dependencies: []
scope: "all"
breaking: false
type: "cleanup"
mainSession: false
---

# 2: Fix forensic hardcoded developer paths

## Description

`forensic/SKILL.md` contains hardcoded paths specific to the original developer's machine (`~/.claude/projects/-Users-fanhuifeng-...`) and a stale build command referencing `~/.zcode-forge-cli/task`. These break for every other user.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal, P0 finding #3

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/forensic/SKILL.md` | Replace hardcoded paths, fix stale build command |

## Acceptance Criteria
- `grep -r 'fanhuifeng' plugins/forge/skills/forensic/` returns 0 hits
- `grep -r 'zcode-forge-cli' plugins/forge/skills/forensic/` returns 0 hits
- All example paths in forensic/SKILL.md use generic placeholders or instruct users to derive paths from `forge forensic search` output

## Hard Rules
- Example paths must work on any user's machine — use `<SESSION_ID>` or `<PROJECT_PATH>` placeholders

## Implementation Notes
1. Lines 89, 93, 100: Replace `~/.claude/projects/-Users-fanhuifeng-Projects-ai-coding-harness-forge/<SESSION_ID>.jsonl` with generic instruction: "Use `forge forensic search` to find session paths for the target project."
2. Line 35: Replace stale build command `cd forge-cli && go build -o ~/.zcode-forge-cli/task ./cmd/task/` with current command: `forge init` or remove entirely if forge CLI is pre-built.
