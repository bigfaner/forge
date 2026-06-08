---
status: "completed"
started: "2026-06-08 12:49"
completed: "2026-06-08 12:51"
time_spent: "~2m"
---

# Task Record: 1 Update eval templates to non-blocking diagnostic model

## Summary
Updated eval-journey.md and eval-contract.md templates from blocking quality gates to non-blocking diagnostic evaluations: replaced /eval --type calls with dedicated slash commands, removed hardcoded 850/1000 thresholds, changed AC to report-generation only, removed blocking guardrail text

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/templates/eval-journey.md
- forge-cli/pkg/task/templates/eval-contract.md

### Key Decisions
无

## Document Metrics
2 template files updated, 4 AC items validated

## Referenced Documents
- docs/proposals/eval-diagnostic-mode/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] eval-journey.md calls /eval-journey with zero references to hardcoded 850
- [x] eval-contract.md calls /eval-contract with zero references to hardcoded 850
- [x] Both templates AC sections require only eval report generation, not score thresholds
- [x] grep -r '>= 850' forge-cli/pkg/task/templates/eval-*.md returns zero results

## Notes
Followed proposal.md Scope > In Scope as authoritative spec. Removed blocking guardrail paragraphs and hardcoded thresholds from both templates.
