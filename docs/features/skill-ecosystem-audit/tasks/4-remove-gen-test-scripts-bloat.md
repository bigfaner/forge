---
id: "4"
title: "Remove gen-test-scripts dead weight (templates + 27MB node_modules)"
priority: "P1"
estimated_time: "2h"
dependencies: [3]
scope: "all"
breaking: false
type: "cleanup"
mainSession: false
---

# 4: Remove gen-test-scripts dead weight (templates + 27MB node_modules)

## Description

`gen-test-scripts/templates/` contains 14 Playwright-specific template files that are bit-identical copies of the `web-playwright` profile's templates, plus a 27MB `node_modules/` directory needed only by `validate-specs.mjs`. Remove the duplicates and relocate the validation tooling.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — P0 finding #1, P1 finding #6

## Affected Files

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/gen-test-scripts/templates/node_modules/` (27MB, 31 packages) | Needed only by validate-specs.mjs; should be a project-level devDependency |
| `plugins/forge/skills/gen-test-scripts/templates/api.spec.ts` | Duplicate of web-playwright profile template |
| `plugins/forge/skills/gen-test-scripts/templates/auth-setup.ts` | Duplicate of web-playwright profile template |
| `plugins/forge/skills/gen-test-scripts/templates/cli.spec.ts` | Duplicate of web-playwright profile template |
| `plugins/forge/skills/gen-test-scripts/templates/helpers.ts` | Duplicate of web-playwright profile template |
| `plugins/forge/skills/gen-test-scripts/templates/playwright-ui.spec.ts` | Duplicate of web-playwright profile template |
| `plugins/forge/skills/gen-test-scripts/templates/playwright.config.ts` | Duplicate of web-playwright profile template |
| `plugins/forge/skills/gen-test-scripts/templates/package.json` | Duplicate of web-playwright profile template |
| `plugins/forge/skills/gen-test-scripts/templates/tsconfig.json` | Duplicate of web-playwright profile template |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Remove/redirect references to deleted templates; delegate fully to profile system |

### Relocate
| File | Destination |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/templates/validate-specs.mjs` | Move to `forge-cli/pkg/profile/profiles/web-playwright/` as profile-owned tooling |
| `plugins/forge/skills/gen-test-scripts/templates/validate-specs.test.mjs` | Move alongside validate-specs.mjs |
| `plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/` | Move alongside validate-specs.mjs |

## Acceptance Criteria
- `test -d plugins/forge/skills/gen-test-scripts/templates/node_modules && echo FAIL || echo OK` → OK
- `ls plugins/forge/skills/gen-test-scripts/templates/*.ts 2>/dev/null | wc -l` returns 0
- gen-test-scripts SKILL.md references `{profile-templates-dir}` for all template paths, no local template references
- validate-specs.mjs relocated to web-playwright profile directory
- Plugin install size reduced by ~27MB

## Hard Rules
- Verify the web-playwright profile templates are bit-identical BEFORE deleting the gen-test-scripts copies (use `diff -r`)
- Do NOT delete validate-specs.mjs — relocate it
- Extract Playwright-specific auth instructions (SKILL.md lines 359-371) to the profile's generate.md

## Implementation Notes
- The 14 template files are confirmed identical to `forge-cli/pkg/profile/profiles/web-playwright/templates/`. Verify with diff before deletion.
- `validate-specs.mjs` uses ts-morph for AST validation — it's Playwright-specific and belongs in the profile.
- After cleanup, the gen-test-scripts/templates/ directory should only contain shared (non-profile-specific) files if any remain.
