---
id: "4"
title: "Update guide.md to reflect spec drift detection flow"
priority: "P2"
estimated_time: "30min-1h"
dependencies: ["1", "2", "3"]
type: "documentation"
mainSession: false
---

# 4: Update guide.md to reflect spec drift detection flow

## Description

Update `plugins/forge/hooks/guide.md` to document the spec drift detection capability: the consolidate-specs skill now performs drift audit + auto-fix, and quick mode includes T-quick-6 for drift detection.

## Reference Files
- `docs/proposals/spec-drift-detection/proposal.md` — Source proposal
- `plugins/forge/hooks/guide.md` — File to modify

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Update Skill Workflow diagram, Quick Mode section, and consolidate-specs references |

## Acceptance Criteria

- [ ] Skill Workflow mermaid diagram updated: T-test-5 node shows "consolidate-specs + drift audit"
- [ ] Quick Mode section updated: T-quick-1~6 (was T-quick-1~5), mentions drift detection as the final test step
- [ ] `specs/` rule in Directory Conventions updated to mention drift detection
- [ ] Agent note about `docs/business-rules/` and `docs/conventions/` updated to mention drift verification

## Hard Rules

- Only modify `plugins/forge/hooks/guide.md`
- Keep changes minimal and additive — don't restructure existing sections

## Implementation Notes

- The guide.md is injected via session-start hook, so changes take effect on next session
- The mermaid diagram update should be minimal — just extend the existing T-test-5 label to include drift audit
