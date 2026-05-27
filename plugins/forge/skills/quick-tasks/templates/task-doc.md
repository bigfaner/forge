---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "{{PRIORITY}}"
estimated_time: "{{ESTIMATED_TIME}}"
dependencies: [{{DEPENDENCIES}}]
type: "doc"
mainSession: false
# Note: surface-key and surface-type fields are intentionally absent from doc tasks.
# Doc tasks produce non-compilable output (markdown, specs, templates) and do not
# interact with the quality gate or test pipeline, so surface routing is unnecessary.
---

# {{ID}}: {{TITLE}}

## Description
{{DESCRIPTION}}

## Reference Files
- `docs/proposals/{{SLUG}}/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| {{NEW_FILES}} |

### Modify
| File | Changes |
|------|---------|
| {{MODIFIED_FILES}} |

### Delete
| File | Reason |
|------|--------|
| {{DELETED_FILES}} |

## Acceptance Criteria
{{ACCEPTANCE_CRITERIA}}

## Hard Rules
{{HARD_RULES}}

## Implementation Notes
{{NOTES}}
