---
id: "2b"
title: "Verify and finalize gate/doc/validation template edits"
priority: "P0"
estimated_time: "20min"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2b: Verify and finalize gate/doc/validation template edits

## Description
The interrupted Task 2 already applied content slimming edits to all 21 prompt templates. This task verifies and finalizes the gate, doc, doc-review, validation-code, and validation-ux templates. These templates have AC verification blocks and Record Fields but no CODING_PRINCIPLES.

## Reference Files
- forge-cli/pkg/prompt/templates/gate.md: AC block + Record Fields verification (source: proposal.md#Key-Scenarios-2)
- forge-cli/pkg/prompt/templates/doc.md: AC block + Record Fields verification (source: proposal.md#Key-Scenarios-2)
- forge-cli/pkg/prompt/templates/doc-review.md: AC block + Record Fields verification (source: proposal.md#Key-Scenarios-2)
- forge-cli/pkg/prompt/templates/validation-code.md: AC block + Record Fields + role verification (source: proposal.md#Key-Scenarios-4)
- forge-cli/pkg/prompt/templates/validation-ux.md: AC block + Record Fields + role verification (source: proposal.md#Key-Scenarios-4)

## Acceptance Criteria
- [ ] SC1: All instruction/constraint nodes from functional snapshot retained in all 5 templates
- [ ] AC verification blocks compressed consistently (same pattern as coding-* templates)
- [ ] Record Fields field names and value structures preserved
- [ ] Role descriptions converted to imperative sentences
- [ ] Consistency check: gate/doc templates follow the same slimming pattern

## Hard Rules
- **Only modify these files**: gate.md, doc.md, doc-review.md, validation-code.md, validation-ux.md

## Implementation Notes
- These templates are simpler than coding-* — no CODING_PRINCIPLES to verify
- Review `git diff` for each file against committed baseline
