---
created: "2026-05-15"
tags: [testing, architecture]
---

# Forge Skill File Changes Wrongly Trigger Test Pipeline

## Problem

A feature that only modifies forge skill files (SKILL.md, templates, rubrics, config) generated a full test pipeline (T-quick-1 through T-quick-5) with `go-test` profile. The test tasks are meaningless — there is no Go code to test, only markdown file changes.

Example: `tui-ui-design` feature adds TUI platform support to `ui-design` skill. All 6 tasks modify markdown files under `plugins/forge/skills/`, but all were typed as `"implementation"`, causing `forge task index` to generate 5 test pipeline tasks.

## Root Cause

Causal chain (5 levels):

1. **Symptom**: T-quick-1~5 generated for a docs-only feature
2. **Direct cause**: All 6 business tasks have `type: "implementation"` in frontmatter
3. **`isDocsOnlyFeature()` logic** (`forge-cli/pkg/task/build.go`): returns `false` if ANY business task has type `implementation` or `fix` — so test tasks get generated
4. **SKILL.md guidance gap**: quick-tasks SKILL.md says "set business tasks to `implementation`, documentation tasks to `documentation`" but provides no decision criteria for distinguishing them
5. **Domain mismatch**: forge skill files (SKILL.md, templates, rubrics) are "implementation" in the forge domain (they implement pipeline behavior) but "documentation" in the test pipeline domain (e2e test infrastructure cannot test markdown changes)

## Solution

For features whose scope is entirely within `plugins/forge/skills/` (SKILL.md, templates, rubrics, configs), set business task `type: "documentation"` so `isDocsOnlyFeature()` returns `true`, generating `T-eval-doc` instead of T-quick-1~5.

Decision rule:
- Tasks modifying files under `plugins/forge/skills/` → `type: "documentation"`
- Tasks modifying Go/JS/TS source code → `type: "implementation"`
- Tasks fixing bugs in code → `type: "fix"`

## Reusable Pattern

**When creating tasks for forge-internal features (modifying skill files, templates, rubrics), use `type: "documentation"` on all business tasks.** This triggers the docs-only path in `forge task index`, which generates `T-eval-doc` (documentation quality evaluation) instead of a full test pipeline that cannot meaningfully test markdown files.

Quick check: if all affected files in the task's "Affected Files" section are markdown/YAML files under `plugins/forge/` or `docs/`, the task should be `type: "documentation"`.

## Example

```yaml
# WRONG — triggers test pipeline for markdown-only changes
type: "implementation"

# CORRECT — triggers docs-only path (T-eval-doc)
type: "documentation"
```

Affected path patterns that indicate documentation type:
- `plugins/forge/skills/*/SKILL.md`
- `plugins/forge/skills/*/templates/*.md`
- `plugins/forge/skills/*/eval-ui/templates/rubric-*.md`
- `docs/conventions/`
- `docs/business-rules/`

## Related Files

- `forge-cli/pkg/task/build.go` — `isDocsOnlyFeature()` function (lines 400-410)
- `plugins/forge/skills/quick-tasks/SKILL.md` — Step 3 "Type Assignment" (line 91-93)
- `docs/features/tui-ui-design/tasks/index.json` — example of the problem

## References

- Forge Guide > Quick Mode > "Docs-only features auto-detected: no test tasks, generates T-eval-doc instead"
