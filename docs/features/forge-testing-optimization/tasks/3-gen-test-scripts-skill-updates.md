---
id: "3"
title: "Update gen-test-scripts SKILL.md with Step 4.5 + Prerequisites gate"
priority: "P1"
estimated_time: "50min"
dependencies: ["2"]
status: pending
breaking: false
noTest: false
mainSession: false
---

# 3: Update gen-test-scripts SKILL.md with Step 4.5 + Prerequisites gate

## Description

Update gen-test-scripts SKILL.md to integrate the validation pipeline: add Step 4.5 structural validation after spec generation, and add a Prerequisites gate that blocks when eval-test-cases Step Actionability score is below 20.

## Reference Files
- `docs/proposals/forge-testing-optimization/proposal.md` — Source proposal (Phase 2, Sections 2.2–2.3)
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Target file to modify
- `plugins/forge/skills/gen-test-scripts/templates/validate-specs.mjs` — Validation script from Task 1
- `task-cli/internal/cmd/validate_specs.go` — Command from Task 2

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Add Step 4.5 structural validation; add Prerequisites Step Actionability gate |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] SKILL.md includes Step 4.5: structural validation using `task validate-specs`
- [ ] Step 4.5 runs after Step 4 (spec file generation) and before Step 5 (TypeScript compilation)
- [ ] Step 4.5 behavior: ERROR → block task, report failures; WARNING → report, continue
- [ ] Prerequisites section includes Step Actionability check: if eval-test-cases report exists and Step Actionability < 20, abort with user message
- [ ] Prerequisites check message instructs user to fix test-cases.md before proceeding

## Implementation Notes

1. **Step 4.5 placement**: Insert between current Step 4 (generate spec files) and Step 5 (TypeScript compilation). The validation catches structural issues before compilation
2. **Step Actionability gate**: Check if `docs/features/<slug>/testing/eval-report.md` exists. If it does, parse Step Actionability score. If < 20, abort with clear message. If eval report doesn't exist, proceed (backward compatible)
3. **Step 4.5 flow**: `task validate-specs` → parse output → if ERROR count > 0, mark task blocked and report → if only WARNINGs, append warnings to task record and continue
4. **Keep SKILL.md concise**: Don't repeat the 8 validation rules in SKILL.md — just reference validate-specs.mjs and describe the integration contract
