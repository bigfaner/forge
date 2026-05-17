---
id: "3"
title: "Fix all plugins/forge/ hardcoded path references"
priority: "P1"
estimated_time: "4h"
dependencies: [1]
scope: "all"
breaking: false
type: "refactor"
mainSession: false
---

# 3: Fix all plugins/forge/ hardcoded path references

## Description

22+ references to `plugins/forge/skills/...` paths across 12+ skill files resolve to the user's project root instead of the plugin cache directory. Implement relative-from-self resolution: skills reference their own assets using relative paths from the SKILL.md file location, following the npm `__dirname` convention.

Also fix the gen-test-cases `types/` vs `templates/` path mismatch (2 skills reference `types/` but files live in `templates/`).

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — P1 finding #4, P0 finding #2, Full Path Inventory section

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Lines 85, 347: replace `plugins/forge/skills/eval/rubrics/` with `rubrics/` |
| `plugins/forge/skills/tech-design/SKILL.md` | Line 191: replace path to decision-logging.md |
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Line 177: replace path to decision-logging.md |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Lines 50, 143: fix types/ → templates/ and sitemap.json path |
| `plugins/forge/skills/gen-test-cases/SKILL.md` | Line 101: fix `types/` → `templates/` |
| `plugins/forge/skills/run-e2e-tests/SKILL.md` | Line 132: fix template path |
| `plugins/forge/skills/graduate-tests/SKILL.md` | Lines 120, 122, 154: fix 3 template paths |
| `plugins/forge/skills/improve-harness/SKILL.md` | Line 31: fix rubric path |
| `plugins/forge/skills/improve-harness/templates/improvements.md` | Line 5: fix rubric path |
| `plugins/forge/skills/forensic/SKILL.md` | Lines 118, 144, 183: fix skill/template paths |
| `plugins/forge/skills/init-justfile/SKILL.md` | Lines 183-198: fix 12 template paths |
| `plugins/forge/skills/breakdown-tasks/templates/validate-ux-task.md` | Line 48: fix rubric path |
| `plugins/forge/skills/breakdown-tasks/templates/validate-code-task.md` | Line 45: fix rubric path |
| `plugins/forge/skills/quick-tasks/templates/validate-ux-task.md` | Line 48: fix rubric path |
| `plugins/forge/skills/quick-tasks/templates/validate-code-task.md` | Line 45: fix rubric path |

### Create
| File | Description |
|------|-------------|
| `plugins/forge/references/shared/decision-logging.md` | If not loadable via relative path from skill, document the resolution convention |

## Acceptance Criteria
- `grep -rn 'plugins/forge/' plugins/forge/skills/ plugins/forge/commands/` returns 0 hits (excluding eval reports)
- `grep -rn 'plugins/forge/' plugins/forge/agents/` returns 0 hits
- All path references use relative-from-self convention (e.g., `rubrics/harness.md` instead of `plugins/forge/skills/eval/rubrics/harness.md`)
- Cross-platform: all relative paths use forward slashes
- gen-test-cases SKILL.md references `templates/{type}.md` (not `types/{type}.md`)

## Hard Rules
- Path convention must work when forge is installed as a plugin (not just in source repo)
- Forward slashes only — no backslashes even on Windows
- Do NOT modify init-justfile template references if they already work correctly in the installed plugin (verify first)

## Implementation Notes
- Prerequisite spike (1h): Place a test skill file in the plugin cache directory, invoke it via Claude, confirm that Claude can Read a relative path from the skill's own directory location.
- Convention: When a SKILL.md says "Read `rubrics/harness.md`", Claude resolves this relative to the skill directory (same directory as SKILL.md).
- For references to OTHER skills' files (e.g., gen-test-scripts referencing gen-test-cases), use a path relative to the plugin root that Claude can discover, or reference via CLI command output.
- The 3 command files (gen-sitemap, record-decision, extract-design-md) also have hardcoded paths — fix those too.
