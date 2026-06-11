---
created: "2026-06-08"
author: "fanhuifeng"
status: Draft
intent: "enhancement"
---

# Proposal: Make Eval Pipeline Tasks Non-Blocking Diagnostic

## Problem

Eval task templates (`eval-journey.md`, `eval-contract.md`) hardcode `850/1000` as an acceptance criterion and declare "All [artifacts] scored >= 850/1000" as a blocking AC. When the eval skill finishes its scorer-gate-revise loop without meeting the threshold, the submit-task skill marks the AC as `met: false`, sets `status: "blocked"`, which triggers the dispatcher to create a fix task and increment `consecutive_failures`. After 3 such cycles, the dispatcher stops -- killing a converging improvement loop (e.g., 687->798->820).

### Evidence

Production incident in pm-work-tracker: T-eval-journey scored 687 (iteration 1) -> 798 (iteration 2) with 2/7 journeys passing. The dispatcher's 3-strikes rule stopped the pipeline despite clear score convergence. See `docs/lessons/gotcha-eval-iterative-vs-failure.md`.

### Urgency

This breaks any feature pipeline where eval tasks don't hit the target score on first pass. The eval skill's scorer-gate-revise loop converges over iterations, but the dispatcher's binary failure model prevents that convergence from completing.

## Proposed Solution

Convert eval pipeline tasks from **blocking quality gates** to **diagnostic evaluations**:

1. Task templates call `/eval-journey` and `/eval-contract` slash commands (which already resolve target/iterations from `forge config`)
2. AC changes from "All artifacts scored >= 850/1000" to "Eval report generated"
3. Remove "abort if score below threshold" guardrail text

The task completes as long as the eval report is written. The score is recorded in the report and task record for informational purposes, but does not block the pipeline.

### Innovation Highlights

This follows the pattern of **separation of evaluation from enforcement**: the eval skill measures quality (with its own iterative loop and max-iterations bound), while the pipeline task's only responsibility is ensuring the measurement happened. The quality threshold enforcement belongs in the eval skill's scorer-gate-revise loop, not in the task lifecycle.

## Requirements Analysis

### Key Scenarios

- **Happy path**: Eval runs, score meets target from config, report generated, task completes
- **Below threshold but converging**: Eval runs max iterations, score < target, report generated with score progression, task completes. Score available for manual review
- **Config-driven parameters**: User sets `eval.journey.target: 900` in `.forge/config.yaml`, eval uses that target instead of hardcoded 850
- **No config set**: Eval falls back to rubric default (eval skill handles this)

### Non-Functional Requirements

- Backward compatible: existing `.forge/config.yaml` `eval.*` keys already exist and work
- No Go code changes: only task template files change

### Constraints & Dependencies

- `/eval-journey` and `/eval-contract` commands must exist (they do, at `plugins/forge/commands/`)
- `forge config get eval.{journey,contract}.{target,iterations}` must be functional (it is)
- Eval task templates are Go templates at `forge-cli/pkg/task/templates/`

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD systems typically separate "quality report" stages from "quality gate" stages. A test report stage always succeeds (produces artifacts), while a separate gate stage can fail the pipeline based on thresholds.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | -- | No effort | Dispatcher kills converging eval loops; hardcoded thresholds ignore config | Rejected: breaks real pipelines |
| Patch consecutive_failures counter | lesson option A/B | Minimal dispatcher change | Still treats eval as "pass/fail"; doesn't address root cause (blocking AC) | Rejected: symptomatic fix |
| **Non-blocking diagnostic tasks** | CI/CD report pattern | Eliminates root cause; reuses existing config resolution; no Go code changes | Loses automatic "abort on low quality" for downstream tasks | **Selected: solves root cause with minimal change** |

## Feasibility Assessment

### Technical Feasibility

Fully feasible. The `/eval-journey` and `/eval-contract` commands already implement config resolution (`forge config get eval.journey.target` etc.). Only the Go template files need updating.

### Resource & Timeline

Single developer, < 30 minutes. Two template files to edit.

### Dependency Readiness

All dependencies exist and are stable: slash commands, config keys, eval skill.

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Eval tasks must block when score is below threshold | 5 Whys | Overturned: eval is a measurement, not a gate. The eval skill's scorer-gate-revise loop already handles iterative improvement. The task lifecycle should only care whether the measurement happened |
| Hardcoded 850/1000 is reasonable default | XY Detection | Overturned: the config system already supports per-type targets. Hardcoding in templates is redundant and prevents user customization |
| MainSession evals bypassing failure counter is sufficient | Assumption Flip | Refined: bypassing the counter treats the symptom. Making eval non-blocking eliminates the problem at the source -- no blocked status means no fix task means no counter increment |

Challenge Override: User chose MainSession-exempt approach initially, then refined to non-blocking diagnostic model after examining concrete task template content.

## Scope

### In Scope

- `forge-cli/pkg/task/templates/eval-journey.md`: Replace `/eval --type journey` with `/eval-journey` slash command, remove hardcoded 850/1000, change AC to report-generation, remove blocking guardrail
- `forge-cli/pkg/task/templates/eval-contract.md`: Same changes for contract evaluation

### Out of Scope

- `run-tasks.md` dispatcher consecutive_failures logic (problem eliminated at source)
- `forge-cli/internal/cmd/task/submit.go` (no Go code changes)
- `plugins/forge/skills/eval/` (eval skill unchanged)
- `plugins/forge/commands/eval-journey.md` and `eval-contract.md` (already have config resolution)
- `plugins/forge/commands/run-tasks.md` (no dispatcher changes needed)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Low-quality artifacts proceed to next pipeline stage | M | M | Eval report captures scores; manual review possible. The eval skill's scorer-gate-revise loop still attempts to improve quality within its iteration budget |
| Users unaware that eval tasks no longer enforce thresholds | L | L | Eval report explicitly shows scores vs configured target. Score progression visible in report |

## Success Criteria

- [ ] SC-1: `eval-journey.md` template calls `/eval-journey` (not `/eval --type journey`) and contains zero references to hardcoded `850`
- [ ] SC-2: `eval-contract.md` template calls `/eval-contract` (not `/eval --type contract`) and contains zero references to hardcoded `850`
- [ ] SC-3: Both templates' AC sections require only eval report generation, not score thresholds
- [ ] SC-4: `grep -r ">= 850" forge-cli/pkg/task/templates/eval-*.md` returns zero results

consistency_check_result:
  status: pass
  pairs_checked: 4
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements (or skip directly to implementation given the narrow scope)
