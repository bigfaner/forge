---
id: "1"
title: "Add CLI Query Commands section to guide.md"
priority: "P0"
estimated_time: "15m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Add CLI Query Commands section to guide.md

## Description

Add a concise "CLI Query Commands" section to `plugins/forge/hooks/guide.md` documenting `forge proposal <slug>` and `forge feature status <slug>`. This makes the agent aware of these CLI commands for efficient ad-hoc queries when users mention proposal or feature slugs in conversation.

## Reference Files
- `proposal.md#Solution` — defines the two commands and their placement in guide.md
- `proposal.md#Risks` — mitigation: guide text must clarify CLI is for ad-hoc queries, skills continue using direct file reads

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Add "CLI Query Commands" section after "Task-CLI" section |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `guide.md` contains a new section documenting `forge proposal <slug>` and `forge feature status <slug>`
- [ ] Section includes brief description of each command's output (what info the agent gets)
- [ ] Section clarifies: these commands are for ad-hoc interactive queries; skills should continue using direct file reads for structured data access
- [ ] Existing sections in guide.md remain unchanged

## Hard Rules
- Do NOT modify any skill files or convention documents — only `guide.md`

## Implementation Notes
- Place the new section after "Task-CLI" section, before the end of the file
- Keep it concise — agent reads this every session
- The section should read naturally as guidance, not as a spec
