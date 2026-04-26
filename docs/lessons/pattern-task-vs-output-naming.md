# Standard Test Task File Naming Convention

## Problem

`/breakdown-tasks` generated standard test task files as `T-test-1.md`, `T-test-2.md` — using task IDs instead of task title slugs. This breaks the naming convention used by all other task files.

## Root Cause

The `/breakdown-tasks` skill auto-appends standard test tasks with IDs like `T-test-1`, `T-test-2`. The file naming used these IDs directly instead of converting the task title to a kebab-case slug.

**Causal chain:**
- Symptom: files named `T-test-1.md` instead of descriptive names
- Direct cause: `/breakdown-tasks` used task ID as filename
- Root cause: the skill's template for standard test tasks didn't apply the `{titleSlug}.md` naming convention

## Solution

Task files in `tasks/` must be named by **task title slug**, not task ID.

| Task ID | Task Title | Correct Filename | Wrong Filename |
|---------|-----------|-----------------|----------------|
| `T-test-1` | Generate e2e Test Cases | `gen-test-cases.md` | ~~`T-test-1.md`~~ |
| `T-test-2` | Generate e2e Test Scripts | `gen-test-scripts.md` | ~~`T-test-2.md`~~ |

Compare with existing tasks that follow the convention:
- ID `1.1` + title "Fix Assign column name" → `1.1-fix-assign-column.md`
- ID `4.1` + title "Table rename to pmw_ prefix" → `4.1-table-rename-pmw-prefix.md`

## Key Takeaway

**Task definition files must use `{titleSlug}.md` naming.** When `/breakdown-tasks` generates standard test tasks, rename the files from `T-test-{N}.md` to a slug of the task title (e.g., `gen-test-cases.md`, `gen-test-scripts.md`). The task ID stays in `index.json` and the file's frontmatter `id` field — the filename itself should be human-readable.
