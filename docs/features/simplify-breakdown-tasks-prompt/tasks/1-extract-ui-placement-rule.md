---
id: "1"
title: "Extract ui-placement rule file"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Extract ui-placement rule file

## Description

Extract all UI-related conditional rules from the current `SKILL.md` into a standalone rule file at `rules/ui-placement.md`. This is the largest extraction (~3KB) and consolidates content currently gated by 5 conditional tags: `<HAS_UI>`, `<NO_UI>`, `<UI_ONLY>`, `<HAS_PLACEMENT>`, and `<RULE>`.

The rule file must be self-contained — a contributor reading only this file should understand how UI tasks are generated without reading the skeleton SKILL.md.

## Reference Files
- `docs/proposals/simplify-breakdown-tasks-prompt/proposal.md` — Source proposal (Rule File Extraction Plan, Condition-Rule Matrix)
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Current skill file (source of extracted content)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/breakdown-tasks/rules/ui-placement.md` | UI placement rule file (~3KB) |

### Modify
| File | Changes |
|------|---------|
| (none) | SKILL.md is modified in task 3 |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `rules/ui-placement.md` created under `plugins/forge/skills/breakdown-tasks/rules/`
- [ ] File contains UI-specific element mapping rows (UI Component, Integration Spec, Page composition) from current `<UI_ONLY>` block
- [ ] File contains placement validation procedure from current `<HAS_PLACEMENT>` block (Step 1, lines 70-77 of SKILL.md)
- [ ] File contains UI task split rules from current `<RULE>` block (new-page vs existing-page task chains, lines 99-125)
- [ ] File contains UI reference file requirements for Build, Integration, and Page Assembly tasks (lines 243-274)
- [ ] File contains UI dependency layer rules: "UI components form natural layer after models/interfaces, can mock backend" (lines 182-194)
- [ ] File contains UI prototype reading instruction (lines 66-68)
- [ ] Load condition documented at top of file: "Load IF `ui/ui-design.md` exists OR `prd/prd-ui-functions.md` exists"
- [ ] File includes guard clause: "If referenced artifact has no parseable content, skip this rule and proceed"
- [ ] File includes maintenance note listing skeleton sections it depends on (Step 2 element mapping, Step 3 derive phases, Step 4a create task files)
- [ ] File is independently understandable — a reviewer can identify (a) the load condition, (b) the rules enforced, (c) expected output behavior when loaded vs absent

## Hard Rules

- All file paths in the rule file must use skill-relative references (compatible with forge distribution model under `${CLAUDE_SKILL_DIR}`)
- Do NOT modify `SKILL.md` — that is task 3

## Implementation Notes

The current SKILL.md uses 5 conditional tags for UI rules with complex interdependencies:
- `<HAS_UI>` activates when `ui/ui-design.md` exists
- `<NO_UI>` activates when `ui/ui-design.md` does NOT exist
- `<UI_ONLY>` is always co-activated with `<HAS_UI>`
- `<HAS_PLACEMENT>` activates when `prd/prd-ui-functions.md` exists (independent of ui-design.md)
- `<RULE>` has no independent activation — always co-activated with `<HAS_PLACEMENT>`

The merged rule file must flatten these into a single load condition and include all UI-specific rules regardless of which tag originally gated them. The placement format note (canonical form `<mode>:<target-page-value>`) must be preserved.
