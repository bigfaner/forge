---
id: "8"
title: "Mark forge-init-config-sync proposal as Superseded"
priority: "P2"
estimated_time: "15min"
dependencies: []
type: "doc"
mainSession: false
---

# 8: Mark forge-init-config-sync proposal as Superseded

## Description

Update the `forge-init-config-sync` proposal status to `Superseded` since the unify-surfaces proposal covers its scope (config schema changes + init integration).

## Reference Files
- `proposal.md#Constraints-Dependencies` — states forge-init-config-sync is superseded-by this proposal
- `proposal.md#Dependency-Readiness` — implementation must mark it as Superseded-by

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/proposals/forge-init-config-sync/proposal.md` | Update frontmatter `status: Superseded`, add `superseded-by: unify-surfaces` |

## Acceptance Criteria

- [ ] `docs/proposals/forge-init-config-sync/proposal.md` frontmatter updated: `status: Superseded`
- [ ] Frontmatter includes `superseded-by: unify-surfaces` field
- [ ] Proposal body unchanged (preserved for historical reference)

## Hard Rules

- Do NOT delete the proposal file — preserve for historical reference

## Implementation Notes

- Simple frontmatter edit only
