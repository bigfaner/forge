---
id: "3"
title: "Migrate existing files and verify pipeline end-to-end"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2"]
status: pending
breaking: false
noTest: false
mainSession: false
---

# 3: Migrate existing files and verify pipeline end-to-end

## Description

Run the migration tool on any existing test-cases.md files in the forge project, then verify the complete three-phase pipeline works end-to-end.

This is the integration verification step: ensure preprocess.js → model value-map → apply-values.js produce correct results, and that migrate-format.js correctly converts existing files.

## Reference Files
- `docs/proposals/test-cases-yaml-pipeline/proposal.md` — Section 4.3 (migration), Section 3.5 (performance targets)
- `plugins/forge/skills/eval-test-cases/bin/migrate-format.js` — Created in Task 1
- `plugins/forge/skills/eval-test-cases/bin/preprocess.js` — Created in Task 1
- `plugins/forge/skills/eval-test-cases/bin/apply-values.js` — Created in Task 1

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| Any `docs/features/*/testing/test-cases.md` | Converted to `test-cases.yaml` via migrate-format.js (backup preserved) |

### Delete
| File | Reason |
|------|--------|
| (none) | `.md.bak` backup kept |

## Acceptance Criteria

- [ ] All existing `test-cases.md` files under `docs/features/*/testing/` converted to `test-cases.yaml` with `.md.bak` backup preserved
- [ ] Each converted YAML file is valid YAML (passes `node -e "require('js-yaml').load(require('fs').readFileSync('<file>','utf8'))"`)
- [ ] No data loss: every TC from the original `.md` appears in the `.yaml` with all fields preserved
- [ ] `preprocess.js` runs successfully on a converted YAML file and produces correct auto-extracted values
- [ ] `apply-values.js` runs successfully with a sample value-map and produces a valid merged YAML file
- [ ] Round-trip fidelity: convert md→yaml, preprocess, apply empty value-map → output matches input (no data corruption)

## Implementation Notes

1. **Find existing test-cases files**: `find docs/features -name "test-cases.md" -path "*/testing/*"`
2. **Run migration**: For each file found, run `node plugins/forge/skills/eval-test-cases/bin/migrate-format.js --input <path> --output <path-with-.yaml-ext>`
3. **Validate conversion**: For each converted file, load with js-yaml and verify all TC IDs present
4. **Test preprocess**: Pick one converted file, run `preprocess.js`, verify auto-extracted routes are correct
5. **Test apply-values**: Create a minimal value-map file with 1-2 TC value overrides, run `apply-values.js`, verify only specified fields changed
6. **Test round-trip**: Apply empty value-map to preprocessed file → output should equal preprocessed file
7. If no existing test-cases.md files are found in this project, create a small synthetic test-cases.md (3 TCs) to verify migration, then clean up
