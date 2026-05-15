---
id: "3"
title: "Add docs-only exceptions to guide.md"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 3: Add docs-only exceptions to guide.md

## Description

Update `guide.md` to add explicit docs-only exceptions in the Quality Gate Protocol and All-Completed Hook sections. Currently these sections state "All task-executing workflows MUST pass the quality gate" with no exception for documentation tasks (`noTest: true`).

## Reference Files
- `docs/proposals/docs-only-fast-path/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Add docs-only exceptions to Quality Gate Protocol and All-Completed Hook sections |

## Acceptance Criteria

- [ ] Quality Gate Protocol section explicitly states that documentation tasks (`noTest: true`) skip the quality gate
- [ ] All-Completed Hook section explicitly states that `forge quality-gate` already skips docs-only features
- [ ] Both exceptions are concise (1-2 sentences each), not verbose
- [ ] An agent reading only guide.md can determine that docs-only features skip both quality gate and all-completed hook test steps

## Implementation Notes

- The runtime already handles docs-only correctly — `forge quality-gate` skips docs-only features and `noTest: true` is set on doc tasks. These additions make the existing behavior explicit in the agent-facing guide.
- Keep additions minimal — append to existing sections rather than restructuring
