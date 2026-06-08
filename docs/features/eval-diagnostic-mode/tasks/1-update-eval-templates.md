---
id: "1"
title: "Update eval templates to non-blocking diagnostic model"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Update eval templates to non-blocking diagnostic model

## Description
Convert eval-journey and eval-contract task templates from blocking quality gates to non-blocking diagnostic evaluations. Currently both templates hardcode 850/1000 as acceptance criteria and abort on low scores, which causes the dispatcher's 3-strikes rule to kill converging improvement loops. The fix: replace `/eval --type` calls with dedicated slash commands, remove hardcoded thresholds, and change AC to require only report generation.

## Reference Files
- `docs/proposals/eval-diagnostic-mode/proposal.md` — Problem, Proposed Solution, Scope > In Scope, Success Criteria, Key Risks
- `forge-cli/pkg/task/templates/eval-journey.md` — Current journey template with hardcoded 850 threshold (ref: Scope > In Scope)
- `forge-cli/pkg/task/templates/eval-contract.md` — Current contract template with hardcoded 850 threshold (ref: Scope > In Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/templates/eval-journey.md` | Replace `/eval --type journey` with `/eval-journey`, remove hardcoded 850/1000, change AC to report-generation, remove blocking guardrail |
| `forge-cli/pkg/task/templates/eval-contract.md` | Replace `/eval --type contract` with `/eval-contract`, remove hardcoded 850/1000, change AC to report-generation, remove blocking guardrail |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `eval-journey.md` calls `/eval-journey` (not `/eval --type journey`) and contains zero references to hardcoded `850`
- [ ] `eval-contract.md` calls `/eval-contract` (not `/eval --type contract`) and contains zero references to hardcoded `850`
- [ ] Both templates' AC sections require only eval report generation, not score thresholds
- [ ] `grep -r ">= 850" forge-cli/pkg/task/templates/eval-*.md` returns zero results

## Implementation Notes
- Remove the "abort if score below threshold" guardrail paragraph from both templates
- The eval skill's scorer-gate-revise loop already handles iterative improvement — the task template's job is only to trigger the measurement
- Target score and iterations are resolved from `forge config` by the slash commands themselves — no need to hardcode in templates
- Low-quality risk: eval report captures scores vs configured target for manual review (ref: Key Risks)
