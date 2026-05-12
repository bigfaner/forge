---
id: "T-test-5"
title: "Consolidate Specs"
priority: "P2"
estimated_time: "20min"
dependencies: ["T-test-4.5"]
status: pending
type: "doc-generation.consolidate"
---

# T-test-5: Consolidate Specs

## Description

Call `/consolidate-specs` skill to extract business rules and technical specifications from this feature into `specs/` directory for user review before integrating to project-level shared directories.

## Reference Files

- `prd/prd-spec.md` — Source for business rules
- `design/tech-design.md` — Source for technical specs

## Acceptance Criteria

- [ ] `specs/` directory created with extracted specs
- [ ] User reviews and confirms integration to project-level directories

## Implementation Notes

1. Run `/consolidate-specs` skill
2. Present preview to user for review
