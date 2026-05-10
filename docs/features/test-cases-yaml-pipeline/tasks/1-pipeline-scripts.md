---
id: "1"
title: "Create pipeline scripts and YAML template"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
status: pending
breaking: false
noTest: false
mainSession: false
---

# 1: Create pipeline scripts and YAML template

## Description

Create the three Node.js scripts that form the new eval-test-cases revision pipeline infrastructure, plus a new YAML template to replace the markdown table template for gen-test-cases.

**Context**: The current eval-test-cases revision pipeline requires doc-reviser to make 83 Edit tool calls on a 32KB file, causing 30+ minute stalls. The new pipeline separates concerns: scripts handle mechanical transformations, model only generates structured values.

The three scripts are:
1. **preprocess.js** — Parses YAML test-cases, auto-extracts derivable values (CLI/API routes from Steps, traceability table, route validation skeleton), outputs preprocessed YAML
2. **apply-values.js** — Reads preprocessed YAML + model value-map, merges by TC ID, writes final test-cases.yaml
3. **migrate-format.js** — One-time migration tool converting markdown table format to YAML

## Reference Files
- `docs/proposals/test-cases-yaml-pipeline/proposal.md` — Source proposal (Section 3: pipeline architecture, Section 4: format spec)
- `plugins/forge/skills/gen-test-cases/templates/test-cases.md` — Current markdown template (to understand field structure)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval-test-cases/bin/preprocess.js` | YAML parser + auto-extract CLI/API routes from Steps + generate traceability + route validation skeletons |
| `plugins/forge/skills/eval-test-cases/bin/apply-values.js` | Merge model value-map into preprocessed YAML by TC ID; replace Expected fields; insert new TCs |
| `plugins/forge/skills/eval-test-cases/bin/migrate-format.js` | Convert markdown table test-cases.md to YAML format |
| `plugins/forge/skills/gen-test-cases/templates/test-cases.yaml` | YAML template replacing test-cases.md |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | test-cases.md template kept for backward compat until migration complete |

## Acceptance Criteria

- [ ] `preprocess.js` accepts `--input <yaml> --out <yaml>` flags; parses YAML TCs; auto-derives CLI routes from Steps (commands like `agent-forensic --lang en`), API routes from Steps (function calls); generates traceability list and route validation skeleton sections
- [ ] `apply-values.js` accepts `--base <yaml> --values <yaml> --out <yaml>` flags; reads both files; merges `tc_values` by TC ID (route, element fields); rewrites `expected` from `expected_rewrites`; appends `new_test_cases`; regenerates traceability and route validation sections
- [ ] `migrate-format.js` accepts `--input <md> --output <yaml>` flags; converts markdown table format to YAML list-of-objects; preserves all fields; creates backup of original .md file
- [ ] `test-cases.yaml` template mirrors all fields from current `test-cases.md` template but in YAML format (see proposal Section 4.2 for target format)
- [ ] All scripts use `js-yaml` or built-in YAML parsing (no exotic dependencies)
- [ ] All scripts include `--help` flag with usage instructions
- [ ] All scripts validate inputs and exit with non-zero code on error

## Implementation Notes

1. **preprocess.js auto-extraction logic**:
   - CLI route: regex match `Run \`<command>\`` from Steps → extract command as route value
   - API route: regex match function call patterns from Steps → extract as route value
   - Only fill route if currently empty/missing/`???`
   - Generate traceability section by iterating all TCs, extracting `{tc_id, source, type, target, priority}`
   - Generate route validation section by grouping TCs by route

2. **apply-values.js merge logic**:
   - Match `tc_values` entries by TC ID to base YAML TCs
   - For each match: merge `route` and `element` fields (overwrite if present)
   - For `expected_rewrites`: match by TC ID, replace entire `expected` field
   - For `new_test_cases`: append to the appropriate type section
   - After all merges: regenerate traceability and route validation sections from the merged data
   - Include schema validation: fail gracefully if value-map has invalid TC IDs

3. **migrate-format.js conversion**:
   - Parse markdown table rows (Field/Value pairs) into key-value pairs per TC
   - Convert `Steps` numbered list (`1. Step...| 2. ...`) to YAML list
   - Convert `Pre-conditions` to `preconditions` list
   - Convert traceability markdown table to YAML list
   - Preserve frontmatter (convert `source` string to `sources` list)
   - Write `.md.bak` backup before creating .yaml

4. **test-cases.yaml template**: Use the proposal Section 4.2 format. Include:
   - Frontmatter with `feature`, `sources` (list), `generated`
   - Summary section (comment, not YAML — for human readability)
   - TC entries as YAML list with all standard fields
   - Traceability and Route Validation as separate YAML list sections
