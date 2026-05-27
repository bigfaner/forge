---
id: "1"
title: "Fix prompt template + add journey isolation to run-tests skill"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Fix prompt template + add journey isolation to run-tests skill

## Description

Two related changes to the run-tests skill and its prompt template:

1. **P0 Bug Fix**: `forge-cli/pkg/prompt/data/test-run.md` references `Skill(skill="forge:run-e2e-tests")` which doesn't exist. Must be changed to `Skill(skill="forge:run-tests")`.

2. **P1 Journey Isolation**: The run-tests skill currently executes `just test` (no journey parameter), running all surface tests. Must be changed to discover current feature's journeys via `ls docs/features/<slug>/testing/` and execute `just test <journey>` for each discovered journey.

3. **Surface orchestration combo**: For web/api surfaces, the dev→probe lifecycle wraps all journeys (start once, loop per-journey test, teardown once). For cli/tui surfaces, the simplified test→teardown sequence applies per the existing surface rules.

## Reference Files
- `proposal.md#Proposed-Solution` — defines the four-layer fix and journey discovery mechanism
- `proposal.md#Requirements-Analysis` — Key Scenarios (happy path, no journey edge case, surface orchestration combo)
- `proposal.md#Key-Risks` — risk of ls directory scan failure on non-standard structures
- `proposal.md#Constraints-&-Dependencies` — justfile already supports `just test <journey>` parameter

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/test-run.md` | Fix `forge:run-e2e-tests` → `forge:run-tests` (lines 11, 31) |
| `plugins/forge/skills/run-tests/SKILL.md` | Add journey discovery step (Step 1.5), modify Step 4 to loop per-journey, update output format |
| `plugins/forge/skills/run-tests/rules/surfaces/cli.md` | Update Journey 过滤 section to reflect per-journey execution |
| `plugins/forge/skills/run-tests/rules/surfaces/web.md` | Update Journey 过滤 section to reflect per-journey execution |
| `plugins/forge/skills/run-tests/rules/surfaces/api.md` | Update Journey 过滤 section if exists |

## Acceptance Criteria

- [ ] `forge-cli/pkg/prompt/data/test-run.md` references `forge:run-tests` (not `forge:run-e2e-tests`)
- [ ] SKILL.md includes a journey discovery step that runs `ls docs/features/<slug>/testing/` before test execution
- [ ] SKILL.md specifies per-journey execution: `just test <journey>` for each discovered journey directory
- [ ] SKILL.md specifies dev/probe execute once, per-journey loop for test, teardown once (for web/api/mobile)
- [ ] SKILL.md handles the "no journey" edge case: `docs/features/<slug>/testing/` missing or empty → error message suggesting run gen-journeys first
- [ ] Surface rule files updated to reflect per-journey test execution pattern

## Hard Rules

- Journey discovery MUST use `ls docs/features/<slug>/testing/` — no forge CLI command, no new abstraction
- `just test <journey>` is the only test invocation pattern (justfile already supports the journey parameter)
- Surface orchestration sequence (dev/probe/test/teardown) behavior is unchanged; only test step changes from single `just test` to loop `just test <journey>`

## Implementation Notes

- The prompt template at `forge-cli/pkg/prompt/data/test-run.md` is embedded in Go via `//go:embed`. After modifying, the forge binary must be rebuilt for changes to take effect at runtime.
- The SKILL.md workflow is: Step 0 (Stale State Recovery) → Step 1 (Detect Surface) → **NEW Step 1.5 (Discover Journeys)** → Step 2 (Load Rules) → Step 3 (Env Check) → Step 4 (Execute Sequence, now per-journey) → Step 5 (Parse Results) → Step 6 (Generate Report)
- Feature slug comes from the task's `{{FEATURE_SLUG}}` template variable or from `forge feature status` CLI
- For the surface orchestration combo (web/api): the existing dev→probe→test→teardown sequence wraps the journey loop. Dev and probe run once. Then for each journey: execute `just test <journey>`. Finally teardown runs once.
