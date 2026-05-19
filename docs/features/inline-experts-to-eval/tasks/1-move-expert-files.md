---
id: "1"
title: "Move expert files from agents/ to skills/eval/experts/"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Move expert files from agents/ to skills/eval/experts/

## Description
Move 11 expert/protocol files from `agents/experts/` to `skills/eval/experts/`. These files are scoring templates used by the eval skill at runtime — they are not independent agents. Placing them inside the eval skill directory aligns with their actual purpose and eliminates 11 phantom agents from the `/agents` panel. After the move, delete the now-empty `agents/experts/` directory.

## Reference Files
- `docs/proposals/inline-experts-to-eval/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/experts/protocol/scorer-protocol.md` | Scorer protocol |
| `plugins/forge/skills/eval/experts/protocol/reviser-protocol.md` | Reviser protocol |
| `plugins/forge/skills/eval/experts/scorer/architect.md` | Architect expert |
| `plugins/forge/skills/eval/experts/scorer/code-reviewer.md` | Code reviewer expert |
| `plugins/forge/skills/eval/experts/scorer/cto.md` | CTO expert |
| `plugins/forge/skills/eval/experts/scorer/editor.md` | Editor expert |
| `plugins/forge/skills/eval/experts/scorer/harness-engineer.md` | Harness engineer expert |
| `plugins/forge/skills/eval/experts/scorer/pm.md` | PM expert |
| `plugins/forge/skills/eval/experts/scorer/qa.md` | QA expert |
| `plugins/forge/skills/eval/experts/scorer/ux-auditor.md` | UX auditor expert |
| `plugins/forge/skills/eval/experts/scorer/ux-engineer.md` | UX engineer expert |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/agents/experts/` (entire directory) | Files relocated to skills/eval/experts/ |

## Acceptance Criteria
- [ ] All 11 files exist under `plugins/forge/skills/eval/experts/` with identical content to originals
- [ ] `plugins/forge/agents/experts/` directory no longer exists
- [ ] No `forge:experts:*` agents appear in `/agents` panel after plugin cache refresh

## Hard Rules
- File contents must be identical — no modifications during move
- Directory structure must be preserved exactly: `experts/protocol/` and `experts/scorer/`

## Implementation Notes
- Use `git mv` for clean history tracking
- After move, verify with `ls -R plugins/forge/skills/eval/experts/` that all 11 files are present
- The `agents/` directory may still contain `task-executor.md` — only delete `agents/experts/`
