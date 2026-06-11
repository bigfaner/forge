---
id: "2c"
title: "Verify and finalize test/code-quality template edits"
priority: "P1"
estimated_time: "20min"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2c: Verify and finalize test/code-quality template edits

## Description
The interrupted Task 2 already applied content slimming edits to all 21 prompt templates. This task verifies and finalizes the test-* (test-run, test-gen-scripts, test-gen-contracts, test-gen-journeys) and code-quality-simplify templates. These are the simplest templates — mainly role description and Step 2 description slimming.

## Reference Files
- forge-cli/pkg/prompt/templates/test-run.md: Step 2 desc + role + Record Fields verification (source: proposal.md#Key-Scenarios-3)
- forge-cli/pkg/prompt/templates/test-gen-scripts.md: Same pattern (source: proposal.md#Key-Scenarios-3)
- forge-cli/pkg/prompt/templates/test-gen-contracts.md: Same pattern (source: proposal.md#Key-Scenarios-3)
- forge-cli/pkg/prompt/templates/test-gen-journeys.md: Same pattern (source: proposal.md#Key-Scenarios-3)
- forge-cli/pkg/prompt/templates/code-quality-simplify.md: Role description verification only (source: proposal.md#Key-Scenarios-4)

## Acceptance Criteria
- [ ] SC5: All Step 2 explanatory descriptions ("This generates X from Y...") fully removed — grep confirms zero residuals
- [ ] Role descriptions converted to imperative sentences
- [ ] Record Fields field names and value structures preserved
- [ ] Consistency check: all test-* templates follow the same slimming pattern

## Hard Rules
- **Only modify these files**: test-run.md, test-gen-scripts.md, test-gen-contracts.md, test-gen-journeys.md, code-quality-simplify.md

## Implementation Notes
- Simplest template group — mainly role + Step 2 desc + Record Fields
- Verify Step 2 descriptions are gone via grep
