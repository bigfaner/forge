---
# Template placeholders:
#   COMPLEXITY — low | medium | high (default: medium)
#   TYPE — coding.feature | coding.enhancement | coding.cleanup | coding.refactor | coding.fix | doc | doc.consolidate | doc.drift (default: coding.feature)
id: "{{ID}}"
title: "{{TITLE}}"
priority: "{{PRIORITY}}"
estimated_time: "{{ESTIMATED_TIME}}"
complexity: "{{COMPLEXITY}}"
dependencies: [{{DEPENDENCIES}}]
surface-key: "{{SURFACE_KEY}}"
surface-type: "{{SURFACE_TYPE}}"
breaking: false
type: "{{TYPE}}"
mainSession: false
---

# {{ID}}: {{TITLE}}

## Description
{{DESCRIPTION}}

## Reference Files
- `docs/proposals/{{SLUG}}/proposal.md` — Source proposal

## Acceptance Criteria
{{ACCEPTANCE_CRITERIA}}

## Hard Rules
{{HARD_RULES}}

## Implementation Notes
{{NOTES}}
