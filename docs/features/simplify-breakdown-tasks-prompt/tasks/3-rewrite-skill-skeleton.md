---
id: "3"
title: "Rewrite SKILL.md skeleton with condition-rule matrix"
priority: "P1"
estimated_time: "2h"
dependencies: ["1", "2"]
type: "doc"
mainSession: false
---

# 3: Rewrite SKILL.md skeleton with condition-rule matrix

## Description

Rewrite `SKILL.md` to contain only always-needed rules with a condition-rule matrix for on-demand rule file loading. Remove all 6 conditional tag blocks (`<HAS_UI>`, `<NO_UI>`, `<UI_ONLY>`, `<HAS_PLACEMENT>`, `<RULE>`, `<HAS_DB>`). Replace the tag-based conditional system with a flat dispatch table where each row maps one file-existence check to one rule file.

The skeleton must be complete for all features without any rule files — rule files are additive, not required.

## Reference Files
- `docs/proposals/simplify-breakdown-tasks-prompt/proposal.md` — Source proposal (Condition-Rule Matrix, Load Model, Token Savings)
- `plugins/forge/skills/breakdown-tasks/rules/ui-placement.md` — Created in task 1
- `plugins/forge/skills/breakdown-tasks/rules/phase-detection.md` — Created in task 2
- `plugins/forge/skills/breakdown-tasks/rules/db-schema.md` — Created in task 2
- `plugins/forge/skills/breakdown-tasks/rules/existing-code-split.md` — Created in task 2

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Rewrite from 421 lines / ~23KB to ~160 lines / ~8KB skeleton |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] SKILL.md reduced from ~23KB to ≤8KB (target: ~160 lines)
- [ ] All 6 conditional tags removed: `<HAS_UI>`, `<NO_UI>`, `<UI_ONLY>`, `<HAS_PLACEMENT>`, `<RULE>`, `<HAS_DB>`
- [ ] Condition-rule matrix added as the **first instruction block** after Prerequisites section
- [ ] Matrix contains exactly 4 rows, each with a single file-existence check:
  - `rules/phase-detection.md` — IF PRD has phase/gate structure
  - `rules/ui-placement.md` — IF `ui/ui-design.md` exists OR `prd/prd-ui-functions.md` exists
  - `rules/db-schema.md` — IF `design/er-diagram.md` exists
  - `rules/existing-code-split.md` — IF tech-design modifies existing shared code
- [ ] Each step that needs a rule file re-prints the load instruction inline as safety net
- [ ] Skeleton never contains rule content inline — fail-safe design (missing rules produce simpler output, not wrong output)
- [ ] Always-needed rules preserved inline:
  - Process flow (Steps 0-7 overview)
  - Element mapping table (non-UI rows only: Interface, Data Model, Backend, Error Type, PRD Flow Gate)
  - Scope assignment algorithm
  - Type assignment table and classification rules
  - Intent propagation rules
  - Template selection rules
  - PRD coverage verification
  - Task granularity guidelines (1-4h)
  - Dependency principles (base rules)
  - Step 5 (forge task index), Step 6 (validate), Step 7 (update manifest)
  - Output checklist (adapted for conditional structure)
  - Docs-Only Fast Path section
- [ ] Prerequisites section updated: remove conditional tag documentation, keep artifact table
- [ ] Step 1 (Read Documents) updated: read manifest.md + artifacts, then "IF applicable rule file loaded, apply its read instructions"
- [ ] Step 2 (Map → Tasks) updated: base element mapping inline, conditional rows via rule files
- [ ] Step 3 (Derive Phases) updated: "Apply rules/phase-detection.md if loaded, else artifact-driven decomposition"
- [ ] Step 4a (Create Task Files) updated: scope/type/template always inline, conditional rules via rule files
- [ ] Steps 5-7 unchanged (CLI-driven, not affected by refactor)

## Hard Rules

- All file paths in the skeleton must use skill-relative references (`rules/X.md`), compatible with forge distribution model under `${CLAUDE_SKILL_DIR}`
- The condition-rule matrix must be the FIRST instruction block the LLM encounters after Prerequisites — this is the core reliability mechanism
- Each condition must be a single file-existence check with no boolean expressions or nested logic
- The skeleton must be structurally complete without any rule files — a feature with no phases, no UI, no DB, and no existing code modifications must produce valid output using only the skeleton

## Implementation Notes

The Condition-Rule Matrix from the proposal (lines 239-254) provides the target structure:

```
Step 2: Map → Tasks
├─ Read rules/phase-detection.md   IF PRD has phase/gate structure
├─ Add UI mapping rows             IF rules/ui-placement.md loaded
└─ Add DB mapping rows             IF rules/db-schema.md loaded

Step 3: Derive Phases
└─ Apply rules/phase-detection.md  IF loaded (else: artifact-driven decomposition)

Step 4a: Create Task Files
├─ Apply rules/existing-code-split.md  IF loaded (shared code modifications detected)
├─ Apply rules/db-schema.md            IF loaded (DB schema tasks)
├─ Apply rules/ui-placement.md          IF loaded (UI task chains + reference files)
└─ Apply scope/type/template rules      (always — inline in skeleton)
```

Key risk from proposal: "LLM loads rule files unconditionally" (Medium likelihood, High impact). Mitigation: matrix is first section, each step re-checks inline, skeleton never contains rule content.
