---
name: clean-code
description: Simplify and clean up code. Supports scoped cleanup (git diff, files, directories) with optional quality gate.
allowed-tools: Bash Read Edit Write Glob Grep
---

# Clean Code

Code cleanup applying five refinement principles. Supports three scope modes: user-specified paths, git diff, or full feature scope.

**Core principle**: Only modify code within the determined scope. Never change what the code does — only how it does it. Every edit must trace back to the scope file list.

## When to Use

**Trigger conditions:**
- Invoked via `/forge:clean-code` (standalone) — optionally with paths as arguments
- Invoked via pipeline task `T-clean-code-1` (when `auto.cleanCode` is enabled)
- User explicitly requests code cleanup

**Skip when:**
- No files in scope
- All files in scope are documentation-only

## Process Flow

```
Step 1: Scope Detection → Step 2: Code Cleanup → Step 3: Quality Gate (optional) → Step 4: Cleanup Summary
```

## Step 1: Scope Detection

**Output**: a concrete file list. All subsequent steps operate on this list only.

Resolve scope from the first applicable source:

| Priority | Source | When |
|----------|--------|------|
| 1 | User arguments | `/forge:clean-code path/to/file.go pkg/service/` |
| 2 | Git diff | On a feature branch with changes vs base |
| 3 | Current feature | Pipeline task (`T-clean-code-1`) with feature context |

### Priority 1: User-Specified Paths

If the user provided file or directory paths as arguments, use those directly. For directories, list code files within, excluding vendor/dependency directories.

Example for a Go project:
```bash
find <path> -type f -name "*.go" ! -path "*/vendor/*"
```

### Priority 2: Git Diff

If no arguments and on a feature branch:

```bash
git diff --name-only main
```

If the base branch is not `main`, detect it. Prefer offline detection first; fall back to remote query only if offline methods fail:

```bash
# Offline: check common base branch names
(git rev-parse --verify main >/dev/null 2>&1 && echo main) || \
(git rev-parse --verify master >/dev/null 2>&1 && echo master) || \
# Fallback: remote query (requires network)
git remote show origin | grep 'HEAD branch' | awk '{print $NF}'
```

### Priority 3: Feature Context

If invoked as a pipeline task, the feature's changed files are already in the working tree. Use git diff against the base branch (same as Priority 2), or read the feature's task records to collect changed files:

```bash
cat docs/features/<slug>/tasks/index.json | grep -o '"file":"[^"]*"' | cut -d'"' -f4
```

### Scope Validation

Filter out non-code files (`.md`, `.txt`, `.json` unless they are config files — i.e. files that configure tooling, CI, or build behavior such as `tsconfig.json`, `.eslintrc.json`, `jest.config.json`, `justfile`, `Makefile`, `.github/workflows/*.yml`). If scope is empty after filtering:

```
No code files in scope. Nothing to clean up.
```

And stop here.

For large scopes (50+ files), process files in batches of 10-15 to avoid context overflow.

Output: `Step 1/4: Scope detection... DONE (N files in scope, source: <user-specified|git-diff|feature-context>)`

## Step 2: Code Cleanup

Read each file in scope and apply the five refinement principles. Only edit files where cleanup opportunities exist — if a file is already clean, skip it.

### The Principles

1. **Preserve Functionality**: Never change what the code does — only how it does it. All original features, outputs, behaviors, and side effects must remain intact. If unsure whether a change preserves behavior, do not make it.

2. **Apply Project Standards**: Follow the coding standards from CLAUDE.md and project conventions. Match existing code style even if you would do it differently. Respect naming conventions, import ordering, error handling patterns, and structural norms already established in the codebase.

3. **Enhance Clarity**: Simplify code structure by:
   - Reducing unnecessary complexity and nesting
   - Eliminating redundant code and abstractions
   - Improving readability through clear variable and function names
   - Consolidating related logic
   - Removing comments that describe obvious code (keep comments that explain *why*, not *what*)

4. **Maintain Balance**: Avoid over-simplification that could:
   - Reduce code clarity or maintainability
   - Create overly clever solutions that are hard to understand
   - Combine too many concerns into single functions
   - Remove helpful abstractions that improve code organization
   - Prioritize fewer lines over readability (e.g., dense one-liners, nested ternaries)

### What to Clean Up

- Dead code (unused imports, unreachable branches, commented-out code)
- Unnecessary complexity (nested conditionals that can be flattened, redundant checks)
- Poor naming (single-letter variables in non-trivial contexts, misleading names)
- Code duplication within the scope (extract shared logic)
- Unnecessary abstractions (single-use wrapper functions, trivial indirection)
- Overly verbose patterns (boilerplate that adds no clarity)

### What NOT to Clean Up

- Code outside the resolved scope
- Pre-existing code that you did not change (even if adjacent)
- Working abstractions that serve a purpose
- Comments explaining *why* (domain knowledge, non-obvious constraints)
- Error handling for real edge cases

Output: `Step 2/4: Code cleanup... DONE (M files modified, K files skipped)`

## Step 3: Quality Gate (Optional)

After cleanup, verify no regressions were introduced.

### Surface-Aware Recipe Resolution

Determine the unit-test recipe using surface-key aware fallback:

1. **Detect surface-key**: If running in a task execution context, read the current task's `surface-key` from its YAML frontmatter (via `forge prompt get-by-task-id` or the task file). If no task context or no `surface-key`, treat as empty.

2. **Resolve recipe**: When `surface-key` is non-empty, probe for a prefixed recipe first:

```bash
just --dry-run <key>-unit-test 2>/dev/null
```

- **If `<key>-unit-test` exists** (exit code 0): use `<key>-unit-test` as the test recipe
- **If `<key>-unit-test` does not exist**: fall back to the generic `unit-test` recipe

When `surface-key` is empty, always use the generic `unit-test` recipe directly (no probing).

3. **Execute the resolved recipe**:

```bash
# Prefixed (surface-scoped)
just <key>-unit-test

# Or generic (no surface-key, or prefixed not available)
just unit-test
```

If tests fail, the cleanup introduced a regression:
1. Report the failure
2. Revert the changes that caused the failure
3. Re-run tests to confirm the revert fixes the issue
4. Continue to summary with a note about the reverted changes

**If `just unit-test` is not available** (and no prefixed recipe exists): Skip the quality gate.

Output one of:
- `Step 3/4: Quality gate... DONE (tests passed)`
- `Step 3/4: Quality gate... SKIPPED (no just unit-test available)`
- `Step 3/4: Quality gate... DONE (N regressions reverted)`

## Step 4: Cleanup Summary

Output a summary using the template at `templates/summary.md`. Fill in the placeholders with actual counts.

### Template Fields

| Field | Value |
|-------|-------|
| `{{SCOPE_COUNT}}` | Total files in scope |
| `{{SCOPE_SOURCE}}` | Scope resolution source (user-specified / git-diff / feature-context) |
| `{{MODIFIED_COUNT}}` | Files where changes were made |
| `{{SKIPPED_COUNT}}` | Files already clean |
| `{{GATE_RESULT}}` | `passed` / `skipped` / `N regressions reverted` |
| `{{FILE_CHANGES}}` | One `- path/to/file — <description>` line per modified file |

If invoked as a standalone command (not via pipeline task), the summary is the final output.

If invoked via pipeline task (`T-clean-code-1`), invoke the skill after the summary:

```
Skill(skill="forge:submit-task")
```
