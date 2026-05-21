---
id: "3"
title: "Delete dead templates, generalize report, update references"
priority: "P1"
estimated_time: "30m"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: Delete dead templates, generalize report, update references

## Description
Delete `gen-test-scripts/templates/` directory (27MB dead code). Generalize report template from e2e-specific to test-type-agnostic. Update all 5 `run-e2e-tests` references to `run-tests`.

## Reference Files
- `docs/proposals/run-tests-decouple/proposal.md` — Source proposal
- `plugins/forge/skills/run-e2e-tests/templates/e2e-report.md` — Current report template
- `plugins/forge/commands/run-tasks.md` — Has run-e2e-tests reference at line 97
- `plugins/forge/skills/gen-contracts/SKILL.md` — Has run-e2e-tests reference at line 176
- `plugins/forge/skills/gen-test-cases/SKILL.md` — Has run-e2e-tests reference at line 141
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Has run-e2e-tests reference at line 372

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/run-tests/templates/test-report.md` | Generalized report template |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Update `run-e2e-tests` → `run-tests` at line 97 |
| `plugins/forge/skills/gen-contracts/SKILL.md` | Update `run-e2e-tests` → `run-tests` at line 176 |
| `plugins/forge/skills/gen-test-cases/SKILL.md` | Update `run-e2e-tests` → `run-tests` at line 141 |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Update `run-e2e-tests` → `run-tests` at line 372 |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/gen-test-scripts/templates/` (entire directory) | Dead code: validate-specs.mjs + test + fixtures + 27MB node_modules, SKILL.md zero references |
| `plugins/forge/skills/run-e2e-tests/templates/e2e-report.md` | Replaced by test-report.md |

## Acceptance Criteria
- [ ] `gen-test-scripts/templates/` directory deleted (validate-specs.mjs, validate-specs.test.mjs, __test_fixtures__/, node_modules/)
- [ ] Report template renamed: `e2e-report.md` → `test-report.md`
- [ ] Report template contains zero "E2E" or "e2e" references (use "Test" instead)
- [ ] Screenshots section uses conditional rendering (only shown when screenshots discovered)
- [ ] Summary table test types (UI/API/CLI) dynamically generated from parsed results, not hardcoded
- [ ] All 5 files updated: `run-e2e-tests` → `run-tests` in references
- [ ] No remaining `run-e2e-tests` references in `plugins/forge/` (verify with grep)

## Hard Rules
- Verify with `grep -r "run-e2e-tests" plugins/forge/` that zero references remain after update
- The report template must work for non-e2e test types (unit, integration, etc.)

## Implementation Notes
- For the report template: replace `E2E Test Results` with `Test Results`, make Screenshots a conditional block
- Test type classification in report summary: extract types from parsed result data rather than hardcoding UI/API/CLI rows
- When deleting gen-test-scripts/templates/, also check if `.gitignore` or any other config references it
