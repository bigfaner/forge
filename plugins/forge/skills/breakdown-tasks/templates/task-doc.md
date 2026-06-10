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
# Note: complexity field is intentionally absent from doc tasks.
# Complexity is used to calibrate quality-gate thresholds for coding tasks;
# doc tasks skip the quality gate entirely, so complexity has no operational effect.
---

# {{ID}}: {{TITLE}}

## Description
{{DESCRIPTION}}

## Reference Files
{{REFERENCE_FILES}}

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
