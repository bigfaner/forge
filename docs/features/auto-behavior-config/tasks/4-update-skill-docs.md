---
id: "4"
title: "Update skill docs for renamed task IDs and auto config"
priority: "P2"
estimated_time: "45m"
dependencies: ["2"]
type: "documentation"
mainSession: false
---

# 4: Update skill docs for renamed task IDs and auto config

## Description

Update all skill documentation and guide files that reference T-test-5 / T-quick-5 to use the new T-specs-1 / T-quick-specs-1 names. Document the new `auto` config block where relevant.

## Reference Files
- `docs/proposals/auto-behavior-config/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | T-test-5 → T-specs-1; add auto config mention |
| `plugins/forge/skills/consolidate-specs/SKILL.md` | T-test-5 → T-specs-1 in "When to Use" and "Related Skills" |
| `plugins/forge/hooks/guide.md` | T-quick-1~5 → updated naming; add Auto-Behavior Configuration section |

## Acceptance Criteria
- [ ] Zero remaining references to T-test-5 or T-quick-5 in `plugins/forge/`
- [ ] All renamed IDs use consistent new names: T-specs-1, T-quick-specs-1
- [ ] Auto-behavior config documented in guide.md

## Hard Rules
- Grep broadly for `T-test-5` and `T-quick-5` across entire `plugins/forge/` to catch all references
- Preserve technical accuracy — only change naming and add config context

## Implementation Notes
- Use grep for `T-test-5`, `T-quick-5`, `quick-1~5` patterns across `plugins/forge/`
