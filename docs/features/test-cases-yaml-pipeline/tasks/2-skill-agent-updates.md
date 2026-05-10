---
id: "2"
title: "Update skill and agent instructions for YAML pipeline"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
status: pending
breaking: false
noTest: false
mainSession: false
---

# 2: Update skill and agent instructions for YAML pipeline

## Description

Update all skill SKILL.md files and agent markdown files to support the new YAML format and three-phase revision pipeline. This task covers the instruction/prompt layer — no new scripts needed (those are in Task 1).

**Files to modify**:
1. `gen-test-cases/SKILL.md` — Output format from markdown tables to YAML; template reference from `test-cases.md` to `test-cases.yaml`
2. `eval-test-cases/SKILL.md` — Step 4 revision pipeline restructured: preprocess → model generate → apply
3. `agents/doc-reviser.md` — New value-map output mode for test-cases revision (no Edit tool calls)
4. `eval-test-cases/templates/rubric.md` — Adjust "Required Sections" format descriptions from markdown table to YAML list

## Reference Files
- `docs/proposals/test-cases-yaml-pipeline/proposal.md` — Sections 3 (pipeline), 4 (format), 5.3 (rubric), 5.4 (gen-test-scripts compat)
- `plugins/forge/skills/eval-test-cases/SKILL.md` — Current eval-test-cases skill
- `plugins/forge/skills/gen-test-cases/SKILL.md` — Current gen-test-cases skill
- `plugins/forge/agents/doc-reviser.md` — Current doc-reviser agent
- `plugins/forge/skills/eval-test-cases/templates/rubric.md` — Current rubric

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-cases/SKILL.md` | Output format instructions: markdown table → YAML; template path; Step 4 output path `.yaml` not `.md`; Step 3 example format |
| `plugins/forge/skills/eval-test-cases/SKILL.md` | Step 4 revision: three-phase pipeline (preprocess.js → model value-map → apply-values.js); prerequisite checks look for `.yaml`; Step 1 path resolution |
| `plugins/forge/agents/doc-reviser.md` | Add value-map output mode: when revising test-cases, output YAML value-map instead of Edit calls; controlled by input parameter |
| `plugins/forge/skills/eval-test-cases/templates/rubric.md` | Required Sections: update format descriptions from markdown table to YAML list structure |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `gen-test-cases/SKILL.md` references `test-cases.yaml` template and writes `testing/test-cases.yaml`; Step 3 example shows YAML format instead of markdown table; Step 4 output path is `.yaml`
- [ ] `eval-test-cases/SKILL.md` Step 4 describes three-phase revision: (1) run `preprocess.js`, (2) spawn doc-reviser in value-map mode to generate model-output.yaml, (3) run `apply-values.js`; Step 1 looks for `test-cases.yaml` (fallback to `test-cases.md` for compat)
- [ ] `doc-reviser.md` has a conditional value-map output mode: when `OUTPUT_MODE=value-map` input is set, agent reads the preprocessed YAML + attack points and outputs a YAML value-map file instead of using Edit tool; value-map schema matches proposal Section 3.3
- [ ] `rubric.md` Required Sections checklist describes YAML structure (frontmatter with `feature`/`sources`/`generated`, summary comment, grouped TC lists, traceability YAML list, route validation YAML list)
- [ ] All field names in instructions remain unchanged (id, test_id, target, type, source, priority, preconditions, steps, expected, route, element) — only the structural format changes

## Implementation Notes

1. **gen-test-cases/SKILL.md changes**:
   - Step 3 example: replace markdown table template with YAML list template from `test-cases.yaml`
   - Step 4: `Write to docs/features/<slug>/testing/test-cases.yaml` (not `.md`)
   - Template reference: `plugins/forge/skills/gen-test-cases/templates/test-cases.yaml`
   - Step 3.5 route validation: output as YAML list section, not markdown table
   - Keep all field names identical — only structural format changes

2. **eval-test-cases/SKILL.md changes**:
   - Step 1: Look for `test-cases.yaml` first, fallback to `test-cases.md`
   - Step 4 (the critical change): Replace single reviser call with three sub-steps:
     - **4a**: Run `node plugins/forge/skills/eval-test-cases/bin/preprocess.js --input <yaml> --out <preprocessed-yaml>`
     - **4b**: Spawn doc-reviser with `OUTPUT_MODE=value-map` + `PREPROCESSED_PATH=<preprocessed-yaml>` → generates `model-output.yaml`
     - **4c**: Run `node plugins/forge/skills/eval-test-cases/bin/apply-values.js --base <preprocessed-yaml> --values model-output.yaml --out <yaml>`
   - Keep Steps 2, 3, 5, 6 mostly unchanged (scorer reads YAML, gate logic same, report format same, next step same)

3. **doc-reviser.md changes**:
   - Add `OUTPUT_MODE` input (optional, default: `edit`)
   - When `OUTPUT_MODE=value-map`:
     - Skip Edit tool entirely
     - Read `PREPROCESSED_PATH` YAML file
     - Read eval report + rubric as usual
     - Output a single YAML value-map file with schema: `{tc_values: {<TC_ID>: {route, element}}, expected_rewrites: {<TC_ID>: <text>}, new_test_cases: [...], route_validation: [...]}`
     - Write to `{{DOC_DIR}}/model-output.yaml`
   - Keep existing edit mode behavior unchanged (backward compat)

4. **rubric.md changes**:
   - Required Sections: change "Traceability table" to "Traceability YAML list", "Route Validation table" to "Route Validation YAML list"
   - Structure & ID Integrity dimension: update "Summary table matches actual" description — counts in summary comment match actual TC counts per section
   - All dimension checks remain by field name — format-agnostic
