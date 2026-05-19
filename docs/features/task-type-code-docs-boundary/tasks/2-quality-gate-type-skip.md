---
id: "2"
title: "Add type-based quality-gate skip to guide.md and submit-task"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Add type-based quality-gate skip to guide.md and submit-task

## Description
Update two execution-layer documents to enforce the type-based quality-gate skip rule. When `type: "documentation"`, the quality-gate (compile + fmt + lint + test) should be skipped entirely — equivalent to `noTest: true`.

**guide.md**: The quality-gate protocol section currently only mentions `noTest: true` as the skip condition. Add `type: "documentation"` as an equivalent skip trigger.

**submit-task SKILL.md**: The Quality Gate Pre-check section currently documents `noTest: true` as the only skip condition. Update to also document type-based skip, matching the Go code behavior from task 4.

## Reference Files
- `docs/proposals/task-type-code-docs-boundary/proposal.md` — Source proposal
- `plugins/forge/references/shared/type-assignment.md` — Classification rule (task 1 output)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Quality Gate Protocol: add `type: "documentation"` skip rule |
| `plugins/forge/skills/submit-task/SKILL.md` | Quality Gate Pre-check: add type-based skip documentation |

## Acceptance Criteria
- [ ] `guide.md` Quality Gate Protocol explicitly states `type: "documentation"` skips quality-gate (same as `noTest: true`)
- [ ] `submit-task SKILL.md` Quality Gate Pre-check documents that documentation tasks skip quality-gate
- [ ] `noTest: true` mentioned as retained for edge-case override, not deprecated

## Hard Rules
- Must load `docs/conventions/forge-distribution.md` before modifying plugin files
- guide.md uses `${CLAUDE_PLUGIN_ROOT}` paths; submit-task uses `${CLAUDE_SKILL_DIR}` paths — do not mix conventions

## Implementation Notes
- guide.md is loaded via SessionStart hook — changes affect all agent sessions immediately
- submit-task SKILL.md documents Go code behavior; must match what submit.go actually does (task 4 will implement the Go side)
