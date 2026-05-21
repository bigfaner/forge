---
id: "1"
title: "Rewrite run-e2e-tests → run-tests as pure executor"
priority: "P0"
estimated_time: "90m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Rewrite run-e2e-tests → run-tests as pure executor

## Description
Rename `run-e2e-tests` skill directory to `run-tests` and rewrite SKILL.md as a pure executor. The new skill reads execution commands from `.forge/config.yaml` `test.execution` node instead of hardcoded `just e2e-*` commands. Skill only does three things: execute configured commands → parse results → generate report.

## Reference Files
- `docs/proposals/run-tests-decouple/proposal.md` — Source proposal
- `plugins/forge/skills/run-e2e-tests/SKILL.md` — Current skill to rewrite
- `plugins/forge/skills/run-e2e-tests/rules/result-parsing.md` — Parsing logic (keep)
- `plugins/forge/skills/run-e2e-tests/rules/failure-diagnosis.md` — Diagnosis logic (keep)
- `plugins/forge/commands/forge/lib/config-schema.yaml` — Config schema to extend

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/run-tests/SKILL.md` | New pure executor skill |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/forge/lib/config-schema.yaml` | Add `test.execution` schema (run required, setup/pre-check/teardown/results-dir/timeout optional) |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/run-e2e-tests/` | Renamed to run-tests |

## Acceptance Criteria
- [ ] Skill directory renamed: `run-e2e-tests` → `run-tests`
- [ ] SKILL.md frontmatter name changed to `run-tests`
- [ ] SKILL.md contains zero hardcoded `just` or `e2e` commands
- [ ] SKILL.md reads `test.execution` from `.forge/config.yaml` via `forge config get test.execution`
- [ ] Template variables defined: `{slug}` (required), `{journey}`, `{test-dir}`, `{results-dir}` (optional with defaults)
- [ ] Escape rule: `{{var}}` → literal `{var}`
- [ ] Output-flags consistency validation step exists: checks Convention format-type against run command flags before execution
- [ ] Missing `test.execution.run` config produces clear error with config example
- [ ] Missing `{slug}` variable produces clear error prompting `forge feature <slug>`
- [ ] Workflow: load Convention → load config → validate output-flags → setup (optional) → pre-check (optional) → run → parse → report → teardown (optional)
- [ ] Teardown uses state file `.forge/test-state.json` for cross-session reliability
- [ ] `result-parsing.md` and `failure-diagnosis.md` preserved unchanged
- [ ] `test.execution` schema added to config-schema.yaml with all fields documented

## Hard Rules
- Do NOT modify `result-parsing.md` or `failure-diagnosis.md` — these are parsing logic, unchanged
- Do NOT add any hardcoded command names — all commands come from config
- Convention Result Format only provides `format-type` and `Output flags` for parsing, never for execution

## Implementation Notes
- The config schema should document each field with description, type, required/optional, and default values
- Template variable resolution: use sed/envsubst or simple string replacement in bash
- State file format: `{"teardown": "<command>", "timestamp": "<ISO8601>"}` — written before run, deleted after teardown
- For output-flags validation: extract format-type from Convention, extract expected flags, grep run command for presence
