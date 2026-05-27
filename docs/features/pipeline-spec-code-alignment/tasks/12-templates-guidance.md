---
id: "12"
title: "Complete templates and guidance docs"
priority: "P2"
estimated_time: "1.5h"
dependencies: [11]
type: "doc"
mainSession: false
---

# 12: Complete templates and guidance docs

## Description

Fix remaining template issues, placeholder documentation, and guidance gaps. Covers Cluster 8 (issues E1-E14, F1-F16):

1. **manifest-quick.md**: Dual slug placeholder (`{{FEATURE_SLUG}}` vs `{{SLUG}}`) — unify to single placeholder. Also references non-existent `testing/test-cases.md` — fix or remove the reference.

2. **quick-tasks/SKILL.md**: (a) 6 template placeholders have no mapping explanation in SKILL.md — add mapping docs, (b) Output Checklist mentions stage-gate files which are never generated in quick mode — fix, (c) Step 5 claims auto-generates stage-gate files (untrue for quick mode) — fix, (d) `{{REFERENCE_FILES}}` described as template placeholder but no template contains it — clarify.

3. **breakdown-tasks/SKILL.md**: Missing Commit step (Step 8 equivalent). Add commit instructions. Also phase-inventory.json write check is conditional but the code checks unconditionally — fix docs to match.

4. **gen-journeys/SKILL.md** (~line 59): CLI error message example doesn't match actual output. Fix to match.

5. **run-tests/SKILL.md** (lines 84, 228): References `BIZ-error-reporting-001` without providing a resolution path. Add path to the business rule file. Also remove remaining Chinese text if any persists after Task 8.

6. **forge-distribution.md** (~lines 105-106): References `/record-decision` and `/learn-lesson` which have been absorbed into `/learn`. Update references.

7. **prompt-template-hierarchy.md**: Claims "three-level" tag system but `<HARD-RULE>` is a fourth undocumented level. Document it.

8. **journey-contract-model.md**: Check if gen-journeys and gen-contracts each have independent copies. If so, unify to single canonical copy.

9. **execute-task.md** (issue F1): Suspected dead code — verify if called by any automation pipeline. If not, add a note that it's a manual entry point.

10. **clean-code/SKILL.md** (issue F9): References `just test` — should be `just unit-test` (if not already fixed in Task 11).

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence E1-E14, F1-F16
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 8 description
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for manifest-quick.md, quick-tasks SKILL.md

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/templates/manifest-quick.md` | Unify slug placeholder, fix test-cases.md ref |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Placeholder mapping, stage-gate fix, REFERENCE_FILES clarification |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Add Commit step, fix phase-inventory condition |
| `plugins/forge/skills/gen-journeys/SKILL.md` | Error message alignment |
| `plugins/forge/skills/run-tests/SKILL.md` | BIZ-error-reporting path, remaining Chinese text |
| `docs/conventions/forge-distribution.md` | Update `/record-decision`/`/learn-lesson` → `/learn` |
| `docs/conventions/prompt-template-hierarchy.md` | Document `<HARD-RULE>` fourth level |
| `plugins/forge/commands/execute-task.md` | Add manual entry point note if needed |

### Create or Deduplicate
| File | Changes |
|------|---------|
| `docs/reference/journey-contract-model.md` | Verify single canonical copy; deduplicate if needed |

## Acceptance Criteria
- [ ] manifest-quick.md uses single unified slug placeholder
- [ ] manifest-quick.md does not reference non-existent `testing/test-cases.md`
- [ ] quick-tasks/SKILL.md documents all template placeholder mappings
- [ ] quick-tasks/SKILL.md Output Checklist is accurate for quick mode (no stage-gate claims)
- [ ] breakdown-tasks/SKILL.md has a Commit step
- [ ] gen-journeys/SKILL.md error messages match actual CLI output
- [ ] run-tests/SKILL.md `BIZ-error-reporting-001` has resolvable path
- [ ] forge-distribution.md references `/learn` not `/record-decision`/`/learn-lesson`
- [ ] prompt-template-hierarchy.md documents `<HARD-RULE>` as fourth level
- [ ] journey-contract-model.md has single canonical copy (not duplicated)

## Hard Rules
- Do not add new conventions or rules — only fix existing documentation
- For journey-contract-model.md deduplication: keep the copy in the canonical location, update references in both gen-journeys and gen-contracts

## Implementation Notes
- For manifest-quick.md: check which placeholder name the task generation code expects and align the template
- For prompt-template-hierarchy.md: the four levels are likely `<IMPORTANT>`, `<MANDATORY>`, `<HARD-RULE>`, and some base level — document the hierarchy and when each is used
- Low-priority items (F1-F16) should be addressed opportunistically if they're in files already being modified
