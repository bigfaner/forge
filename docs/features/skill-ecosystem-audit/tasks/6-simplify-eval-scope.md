---
id: "6"
title: "Simplify eval validate-ux inline + fix scope heuristic"
priority: "P2"
estimated_time: "2h"
dependencies: []
scope: "all"
breaking: false
type: "refactor"
mainSession: false
---

# 6: Simplify eval validate-ux inline + fix scope heuristic

## Description

Two simplification targets:
1. `eval/SKILL.md` contains a 90-line `validate-ux` sub-pipeline inline (lines 134-225) with project-type detection, PRD-to-operation translation, snapshot format, and 7 impact verification types. Extract this to a separate rubric file under `eval/`.
2. `breakdown-tasks/SKILL.md` Scope Assignment heuristic (lines 290-304) classifies `src/` as frontend, which is wrong for Go/Rust backend projects.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — P2 findings #8 and #9

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rubrics/validate-ux-pipeline.md` | Extracted validate-ux sub-pipeline instructions |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Replace 90-line inline validate-ux (lines 134-225) with a reference to the extracted rubric file |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Fix Scope Assignment heuristic to handle Go/Rust `src/` correctly |

## Acceptance Criteria
- `wc -l plugins/forge/skills/eval/SKILL.md` returns under 280 lines (was 366)
- `test -f plugins/forge/skills/eval/rubrics/validate-ux-pipeline.md` passes
- eval scoring behavior unchanged — run eval on 3 existing test cases and compare scores (must be identical)
- Scope Assignment correctly classifies `src/` as backend for Go/Rust projects (verify with a Go project fixture)

## Hard Rules
- The validate-ux extraction stays WITHIN `eval/` directory as a rubric file — NOT a new skill
- eval scoring output must be byte-identical before and after extraction

## Implementation Notes
- For eval simplification: Move the PRD-to-Operation Translation table, Operation Impact Verification (7 types), and ux-snapshot.md format template from eval/SKILL.md lines 134-225 into `eval/rubrics/validate-ux-pipeline.md`. Replace the inline content with: "For validate-ux type, follow the validate-ux pipeline defined in `rubrics/validate-ux-pipeline.md`."
- For scope heuristic: The current logic says `src/ → frontend`. Fix to: detect language from project config (go.mod → backend, Cargo.toml → backend, package.json → check framework) or add `internal/` as an additional backend indicator.
