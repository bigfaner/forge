---
id: "8"
title: "Add auto-extract trigger to write-prd"
priority: "P2"
estimated_time: "30m"
dependencies: ["3"]
type: "enhancement"
mainSession: false
---

# 8: Add auto-extract trigger to write-prd

## Description

Add knowledge auto-extraction to the `write-prd` skill. After PRD completion, scan PRD content for notable business rules and user-facing constraints. If found, present extracted knowledge for user confirmation.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal (Part 2)
- `plugins/forge/skills/write-prd/SKILL.md` — Current write-prd skill
- `plugins/forge/references/shared/knowledge-extraction.md` — Shared extraction routine (created in Task 3)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/SKILL.md` | Add knowledge review step at PRD completion |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] After PRD completion, a knowledge review step runs
- [ ] Step reads `plugins/forge/references/shared/knowledge-extraction.md` for extraction logic
- [ ] Scans PRD content for new business rules and user-facing constraints
- [ ] Silent when PRD contains no cross-cutting knowledge
- [ ] Presents extracted knowledge via AskUserQuestion for user confirmation
- [ ] Writes confirmed knowledge to appropriate directories using shared formats

## Hard Rules
- Must include the shared extraction routine by reference, not by copying its content
- Do not modify existing PRD generation steps — only add a post-completion review step

## Implementation Notes
- PRD is a natural source of business rules and constraints that should be captured
- The trigger should focus on business rules that apply across features, not feature-specific logic
