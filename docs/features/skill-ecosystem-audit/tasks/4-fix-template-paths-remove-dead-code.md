---
id: "4"
title: "Fix 5 template path references and remove Playwright dead code"
priority: "P0"
estimated_time: "45m"
dependencies: []
type: "cleanup"
scope: "all"
breaking: false
mainSession: false
---

# 4: Fix 5 template path references and remove Playwright dead code

## Description

Five template files still contain hardcoded `plugins/forge/` paths that break in the installed plugin. These are in non-SKILL.md files where `${CLAUDE_SKILL_DIR}` substitution does not apply — they must use relative-from-self paths.

Eight Playwright template files in gen-test-scripts are dead code since the dispatcher refactor (#104) moved test generation to the `types/` dispatch system.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal (W3)

## Affected Files

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/templates/validate-code-task.md` | Line 45: `plugins/forge/skills/eval/rubrics/validate-code.md` → `../../../eval/rubrics/validate-code.md` |
| `plugins/forge/skills/breakdown-tasks/templates/validate-ux-task.md` | Line 48: `plugins/forge/skills/eval/rubrics/validate-ux.md` → `../../../eval/rubrics/validate-ux.md` |
| `plugins/forge/skills/quick-tasks/templates/validate-code-task.md` | Line 45: `plugins/forge/skills/eval/rubrics/validate-code.md` → `../../../eval/rubrics/validate-code.md` |
| `plugins/forge/skills/quick-tasks/templates/validate-ux-task.md` | Line 48: `plugins/forge/skills/eval/rubrics/validate-ux.md` → `../../../eval/rubrics/validate-ux.md` |
| `plugins/forge/skills/improve-harness/templates/improvements.md` | Line 5: `plugins/forge/skills/eval/rubrics/harness.md` → `../../../eval/rubrics/harness.md` |

### Delete

| File | Reason |
|------|--------|
| `plugins/forge/skills/gen-test-scripts/templates/api.spec.ts` | Dead code — types/api.md dispatch handles API test generation |
| `plugins/forge/skills/gen-test-scripts/templates/auth-setup.ts` | Dead code — no longer referenced by dispatcher |
| `plugins/forge/skills/gen-test-scripts/templates/cli.spec.ts` | Dead code — types/cli.md dispatch handles CLI test generation |
| `plugins/forge/skills/gen-test-scripts/templates/helpers.ts` | Dead code — helper functions now in type-specific dispatch |
| `plugins/forge/skills/gen-test-scripts/templates/package.json` | Dead code — generated dynamically by type dispatchers |
| `plugins/forge/skills/gen-test-scripts/templates/playwright.config.ts` | Dead code — generated dynamically by type dispatchers |
| `plugins/forge/skills/gen-test-scripts/templates/playwright-ui.spec.ts` | Dead code — types/ui.md dispatch handles UI test generation |
| `plugins/forge/skills/gen-test-scripts/templates/tsconfig.json` | Dead code — generated dynamically by type dispatchers |

**Keep**: `validate-specs.mjs`, `validate-specs.test.mjs`, `__test_fixtures__/` — used by the validation pipeline.

## Acceptance Criteria

- [ ] Zero `plugins/forge/` hardcoded paths in any source file
  `grep -rn 'plugins/forge/' plugins/forge/skills/ plugins/forge/commands/ plugins/forge/agents/ plugins/forge/references/ | grep -v 'init-justfile'` returns 0 hits
- [ ] `gen-test-scripts/templates/` contains no `.ts` or `.json` files
  `ls plugins/forge/skills/gen-test-scripts/templates/*.ts plugins/forge/skills/gen-test-scripts/templates/*.json 2>/dev/null` returns empty
- [ ] `validate-specs.mjs`, `validate-specs.test.mjs`, `__test_fixtures__/` still exist

## Hard Rules

- Do NOT delete `validate-specs.mjs`, `validate-specs.test.mjs`, or `__test_fixtures__/` — these are used by the validation pipeline.
- Use forward slashes in all relative paths (cross-platform).
- The relative path depth is 3 levels up from template to skill root: `templates/` → `breakdown-tasks/` → `skills/` → `forge/`.

## Implementation Notes

- The `validate-code-task.md` and `validate-ux-task.md` templates exist in both `breakdown-tasks/templates/` and `quick-tasks/templates/` — fix all 4 files identically.
- The `improvements.md` template is in `improve-harness/templates/` — same 3-level-up traversal applies: `improve-harness/templates/` → `improve-harness/` → `skills/` → `forge/`.
