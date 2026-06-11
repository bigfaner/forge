---
id: "2"
title: "Remove Convention Execution command and update init-justfile"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Remove Convention Execution command and update init-justfile

## Description
Remove `Execution command` field from Convention Result Format sections (it's now in `.forge/config.yaml`). Update `init-justfile` Step 3a to read execution commands from Config instead of Convention.

## Reference Files
- `docs/proposals/run-tests-decouple/proposal.md` — Source proposal (Convention vs Config section)
- `docs/conventions/testing-conventions.md` — Master spec with Result Format schema
- `docs/conventions/testing-go.md` — Has Execution command to remove
- `docs/conventions/testing-vitest.md` — Has Execution command to remove
- `docs/conventions/testing-ginkgo.md` — Has Execution command to remove
- `plugins/forge/skills/init-justfile/SKILL.md` — Step 3a needs data source change

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/conventions/testing-conventions.md` | Remove `Execution command` row from Result Format table |
| `docs/conventions/testing-go.md` | Remove `Execution command` line from Result Format section |
| `docs/conventions/testing-vitest.md` | Remove `Execution command` line from Result Format section |
| `docs/conventions/testing-ginkgo.md` | Remove `Execution command` line from Result Format section |
| `plugins/forge/skills/init-justfile/SKILL.md` | Step 3a: read from Config `test.execution.run` instead of Convention `Execution command` |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] `testing-conventions.md` Result Format table no longer has `Execution command` row
- [ ] `testing-go.md`, `testing-vitest.md`, `testing-ginkgo.md` no longer have `Execution command` lines
- [ ] Convention Result Format retains `Output flags` and `Format type` fields (unchanged)
- [ ] `init-justfile` SKILL.md Step 3a references Config `test.execution.run` as the source for e2e-test recipe generation
- [ ] `init-justfile` Step 3a fallback behavior documented: if Config missing, prompt user to configure test.execution

## Hard Rules
- Do NOT modify Convention sections other than Result Format (Framework, Assertion, Tags are untouched)
- Do NOT change init-justfile's recipe names (e2e-setup, e2e-test, e2e-compile remain)
- Only the data source for recipe content changes, not the recipe naming convention

## Implementation Notes
- In `testing-conventions.md`, the Result Format table currently has 3 rows. After removing `Execution command`, it will have 2 rows (`Output flags` and `Format type`)
- In the complete examples (Go, Python, JavaScript), remove the `- **Execution command**: ...` lines from Result Format sections
- In `init-justfile`, Step 3a currently says "Convention provides execution command" — change to "Config provides execution command via test.execution.run"
