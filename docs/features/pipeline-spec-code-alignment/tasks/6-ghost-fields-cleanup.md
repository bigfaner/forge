---
id: "6"
title: "Clean up ghost fields and stale references in skill docs"
priority: "P0"
estimated_time: "2h"
dependencies: [5]
type: "doc"
mainSession: false
---

# 6: Clean up ghost fields and stale references in skill docs

## Description

Remove deprecated files and fix all stale references to non-existent config fields, commands, and files across skill documentation. This covers Cluster 1 from the proposal (issues A1-A10):

1. **Delete** `plugins/forge/skills/breakdown-tasks/rules/scope-assignment.md` — marked deprecated but still present, SKILL.md doesn't instruct to ignore it
2. **Fix `interfaces` → `surfaces`**: quick-tasks/SKILL.md (lines 54, 169) and breakdown-tasks/SKILL.md (line 173) reference `interfaces` config field which doesn't exist; actual field is `surfaces`
3. **Fix `SCOPE` → `SURFACE_KEY`/`SURFACE_TYPE`**: run-tasks.md (line 54) and execute-task.md (line 25) extract `SCOPE` field from `forge task claim` output, but Go code outputs `SURFACE_KEY`/`SURFACE_TYPE`
4. **Fix `decision-logging.md` → `decision-entry.md`**: fix-bug.md (line 233), write-prd/rules/knowledge-extraction.md (line 22), tech-design/rules/knowledge-extraction.md (line 22) reference non-existent `decision-logging.md`; actual file is `learn/templates/decision-entry.md`
5. **Fix scope → surface in rules**: db-schema.md (line 36) uses `scope: "backend"`, existing-code-split.md (lines 24, 32) references deprecated scope assignment algorithm
6. **Fix ui-placement.md** (line 9): uses undefined `<HAS_UI>`/`<UI_ONLY>` condition macros — define them or remove
7. **Fix scope-to-surface-key.md** (line 54): corrects `forge version` reference
8. **Remove `test-template-dir`**: gen-test-scripts/SKILL.md (line 283) references non-existent `test-template-dir` config field
9. **Fix existing-code-split.md** (line 45): maintenance note references non-existent "Step 5: Task Dependencies"

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence A1-A10 (all ghost field and stale reference issues)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 1 description
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Scope` — Cluster 1 In Scope bullets
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for interfaces, SCOPE, decision-logging fixes

## Affected Files

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/breakdown-tasks/rules/scope-assignment.md` | Deprecated, replaced by surface-key/surface-type |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | `interfaces` → `surfaces` config field reference |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | `interfaces` → `surfaces` config field reference |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Remove `test-template-dir` reference |
| `plugins/forge/commands/run-tasks.md` | `SCOPE` → `SURFACE_KEY`/`SURFACE_TYPE` extraction |
| `plugins/forge/commands/execute-task.md` | `SCOPE` → `SURFACE_KEY`/`SURFACE_TYPE` extraction |
| `plugins/forge/commands/fix-bug.md` | `decision-logging.md` → `decision-entry.md` |
| `plugins/forge/skills/write-prd/rules/knowledge-extraction.md` | `decision-logging.md` → `decision-entry.md` |
| `plugins/forge/skills/tech-design/rules/knowledge-extraction.md` | `decision-logging.md` → `decision-entry.md` |
| `plugins/forge/skills/breakdown-tasks/rules/db-schema.md` | `scope: "backend"` → surface-type fields |
| `plugins/forge/skills/breakdown-tasks/rules/existing-code-split.md` | scope assignment → surface inference, fix Step 5 ref |
| `plugins/forge/skills/breakdown-tasks/rules/ui-placement.md` | Define or remove `<HAS_UI>`/`<UI_ONLY>` macros |
| `plugins/forge/skills/breakdown-tasks/rules/scope-to-surface-key.md` | Fix `forge version` reference |

## Acceptance Criteria
- [ ] `scope-assignment.md` deleted
- [ ] No skill doc references `interfaces` config field (grep confirms zero hits)
- [ ] No skill doc extracts `SCOPE` from claim output (use `SURFACE_KEY`/`SURFACE_TYPE`)
- [ ] No skill doc references `decision-logging.md` (use `decision-entry.md`)
- [ ] No skill doc references `test-template-dir` config field
- [ ] `db-schema.md` uses surface-type fields instead of scope
- [ ] `existing-code-split.md` references surface inference, not scope assignment

## Hard Rules
- Do not change the actual behavior of any skill — only fix documentation references
- Preserve all existing prose structure and formatting

## Implementation Notes
- For ui-placement.md macros: if no other rule file defines similar macros, remove them and use plain conditional prose instead
- For existing-code-split.md Step 5 reference: check if the note refers to content that was moved elsewhere
