---
id: "2"
title: "Add sc-consistency rule reference to brainstorm SKILL.md Step 5"
priority: "P1"
estimated_time: "30m"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Add sc-consistency rule reference to brainstorm SKILL.md Step 5

## Description

Modify `plugins/forge/skills/brainstorm/SKILL.md` Step 5 (Write Proposal) to add a mandatory reference to `rules/sc-consistency.md`. This ensures the brainstorm agent executes the SC consistency check after writing Success Criteria and In Scope sections.

## Reference Files

- `proposal.md#Proposed-Solution` — Layer 1 (brainstorm prevention): the rule is applied in Step 5 after SC and InScope are written
- `proposal.md#Key-Risks` — risk of agent ignoring rules; mitigation includes listing consistency check as a mandatory step in SKILL.md

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/brainstorm/SKILL.md` | Add mandatory sc-consistency check step in Step 5 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] SKILL.md Step 5 contains an explicit reference to `rules/sc-consistency.md`
- [ ] The reference is positioned after SC and InScope writing, before the quality standards table
- [ ] Consistency check is described as a mandatory step (not optional), aligning with the "hard protection" strategy from the proposal

## Hard Rules

- Follow `docs/conventions/forge-distribution.md` — use relative path `rules/sc-consistency.md` (skill-internal reference)

## Implementation Notes

- The reference should be concise — a single paragraph or bullet pointing to the rule file, similar to how `rules/challenge-protocol.md` is referenced in Step 2
- The consistency check produces a structured result that becomes part of the proposal writing workflow
