---
id: "1"
title: "Audit all 19 prompt templates for Reference Files strengthening"
priority: "P0"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Audit all 19 prompt templates for Reference Files strengthening

## Description

Audit all 19 prompt template files under `forge-cli/pkg/prompt/data/*.md` to determine which templates need the Reference Files authority declaration and AC per-item validation inserted. This audit drives Task 2's scope.

## Reference Files
- `docs/proposals/spec-authority-enforcement/proposal.md#Scope` — In Scope audit criteria and strengthening conditions
- `docs/proposals/spec-authority-enforcement/proposal.md#Requirements-Analysis` — Key Scenarios and Edge Cases

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none — this task produces analysis only) | |

### Modify
| File | Changes |
|------|---------|
| (none — read-only audit) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] All 19 template files in `forge-cli/pkg/prompt/data/` have been read and analyzed
- [ ] Each template is classified as "needs strengthening" or "skip" with explicit reasoning
- [ ] Audit criteria applied: template needs strengthening if ANY of (a) template is for coding or doc task type, (b) template's Step 1 contains "read task file" step, (c) template involves spec-driven implementation/modification tasks
- [ ] Output a structured audit table: Template Name | Needs Strengthening | Reason | Current Step 1 Location | Current Verify Step Location

## Implementation Notes

- Templates already known to need strengthening: coding-feature.md, coding-enhancement.md, coding-refactor.md, coding-fix.md, coding-cleanup.md (5 coding.* templates)
- doc.md already reads reference files in Step 1 but lacks authority declaration and AC validation — determine if it needs strengthening
- Delegate-only templates (clean-code.md, test-*.md, doc-consolidate.md, doc-drift.md, doc-summary.md) likely don't need strengthening since they delegate to skills
- gate.md, validation-*.md templates need individual assessment
- The audit result will determine Task 2's exact file list
