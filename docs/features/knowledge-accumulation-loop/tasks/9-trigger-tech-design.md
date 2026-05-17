---
id: "9"
title: "Add auto-extract trigger to tech-design"
priority: "P2"
estimated_time: "30m"
dependencies: ["3"]
type: "enhancement"
mainSession: false
---

# 9: Add auto-extract trigger to tech-design

## Description

Add knowledge auto-extraction to the `tech-design` skill. After tech-design completion, scan the design document for architecture decisions, dependency choices, and data model decisions. If found, present extracted knowledge for user confirmation.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal (Part 2)
- `plugins/forge/skills/tech-design/SKILL.md` — Current tech-design skill
- `plugins/forge/references/shared/knowledge-extraction.md` — Shared extraction routine (created in Task 3)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/tech-design/SKILL.md` | Add knowledge review step at design completion |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] After tech-design completion, a knowledge review step runs
- [ ] Step reads `plugins/forge/references/shared/knowledge-extraction.md` for extraction logic
- [ ] Scans design document for architecture decisions, dependency choices, data model decisions
- [ ] Silent when design contains no notable architectural knowledge
- [ ] Presents extracted knowledge via AskUserQuestion for user confirmation
- [ ] Writes confirmed knowledge to appropriate directories using shared formats

## Hard Rules
- Must include the shared extraction routine by reference, not by copying its content
- Do not modify existing tech-design generation steps — only add a post-completion review step
- The existing decision archiving flow (decision-logging.md Section 2) already archives key decisions from tech-design — the auto-extract trigger should complement this, not duplicate it

## Implementation Notes
- tech-design already has a decision archiving step via decision-logging.md Section 2. The auto-extract trigger should focus on lessons and conventions that the existing archiving misses, or on decisions that weren't marked as "key decisions" but are still notable.
- Consider coordinating with the existing archiving flow to avoid duplicate entries
