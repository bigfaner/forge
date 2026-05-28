---
id: "2a"
title: "Verify and finalize coding-* template edits"
priority: "P0"
estimated_time: "30min"
complexity: "high"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2a: Verify and finalize coding-* template edits

## Description
The interrupted Task 2 already applied content slimming edits to all 21 prompt templates. This task verifies and finalizes the 5 coding-* templates (coding-feature, coding-enhancement, coding-fix, coding-cleanup, coding-refactor). Check that all slimming techniques were correctly applied, no functional content was lost, and the templates are consistent with each other.

## Reference Files
- forge-cli/pkg/prompt/templates/coding-feature.md: Verify CODING_PRINCIPLES + AC block + Record Fields + role description edits (source: proposal.md#Key-Scenarios-1)
- forge-cli/pkg/prompt/templates/coding-enhancement.md: Same verification pattern (source: proposal.md#Key-Scenarios-1)
- forge-cli/pkg/prompt/templates/coding-fix.md: Same verification pattern (source: proposal.md#Key-Scenarios-1)
- forge-cli/pkg/prompt/templates/coding-cleanup.md: CODING_PRINCIPLES has cleanup-specific principles (source: proposal.md#Key-Scenarios-1)
- forge-cli/pkg/prompt/templates/coding-refactor.md: CODING_PRINCIPLES has refactor-specific principles (source: proposal.md#Key-Scenarios-1)

## Acceptance Criteria
- [ ] SC1: All instruction/constraint/format nodes from functional snapshot retained in all 5 coding-* templates (100% retention rate)
- [ ] SC3: CODING_PRINCIPLES core constraint instructions preserved — 1 instruction line + 1 boundary summary + 1 representative example per principle
- [ ] AC verification blocks compressed to ~4 lines, retaining AC:REQUIRED and AC:STRONGLY markers
- [ ] Record Fields field names and value structures preserved
- [ ] Consistency check: all 5 coding-* templates follow the same slimming pattern (role → CODING_PRINCIPLES → AC block → Record Fields)

## Hard Rules
- **Line format invariance**: TASK_FILE, TASK_ID, SURFACE_KEY lines must retain `KEY: {{.Value}}` format and position
- **Only modify these files**: coding-feature.md, coding-enhancement.md, coding-fix.md, coding-cleanup.md, coding-refactor.md

## Implementation Notes
- Use coding-feature.md as the reference — other coding-* templates should follow the same pattern
- Run `git diff` on each file to review changes against the committed baseline
- Fix any inconsistencies found during verification
