---
name: clean-code
description: Simplify and clean up code within the current feature branch scope. Applies five cleanup principles with optional quality gate.
allowed_tools: ["Bash", "Read", "Edit", "Write", "Glob", "Grep"]
---

# Clean Code

Scoped code cleanup for the current feature branch. Applies five refinement principles to code within `git diff` scope, with an optional quality gate.

**Core principle**: Only modify code within the feature diff scope. Never change what the code does — only how it does it.

## When to Use

**Trigger conditions:**
- Invoked via `/forge:clean-code` (standalone)
- Invoked via pipeline task `T-clean-code-1` (when `auto.cleanCode` is enabled)
- User explicitly requests code cleanup

**Skip when:**
- No files changed in diff scope
- All changed files are documentation-only

## Workflow

```
Step 1: Scope Detection → Step 2: Code Cleanup → Step 3: Quality Gate (optional) → Step 4: Cleanup Summary
```

## Step 1: Scope Detection

Determine which files to clean up.

```bash
git diff --name-only main
```

If the base branch is not `main`, determine it from context (e.g., the feature branch's merge base).

<HARD-RULE>
**Only modify files within the git diff scope.** Never touch files outside the diff. If a file is not in the diff output, do not edit it — even if it has obvious cleanup opportunities.
</HARD-RULE>

If the diff is empty or contains no code files (only `.md`, `.txt`, etc.), output:

```
No code files in scope. Nothing to clean up.
```

And stop here.

For large diffs (50+ files), process files in batches of 10-15 to avoid context overflow.

Output: `Step 1/4: Scope detection... DONE (N files in scope)`

## Step 2: Code Cleanup

Read each file in scope and apply the five refinement principles. Only edit files where cleanup opportunities exist — if a file is already clean, skip it.

### The Five Principles

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

5. **Focus Scope**: Only refine code that is within the git diff scope. Do not touch adjacent code, even if it has obvious issues. Do not refactor things that are not broken. Every changed line should trace directly to the diff scope.

### What to Clean Up

- Dead code (unused imports, unreachable branches, commented-out code)
- Unnecessary complexity (nested conditionals that can be flattened, redundant checks)
- Poor naming (single-letter variables in non-trivial contexts, misleading names)
- Code duplication within the scope (extract shared logic)
- Unnecessary abstractions (single-use wrapper functions, trivial indirection)
- Overly verbose patterns (boilerplate that adds no clarity)

### What NOT to Clean Up

- Code outside the diff scope
- Pre-existing code that you did not change (even if adjacent)
- Working abstractions that serve a purpose
- Comments explaining *why* (domain knowledge, non-obvious constraints)
- Error handling for real edge cases

<HARD-RULE>
**Every edit must correspond to a file in the diff scope.** If you cannot trace a changed line back to the diff scope, do not change it.
</HARD-RULE>

Output: `Step 2/4: Code cleanup... DONE (M files modified, K files skipped)`

## Step 3: Quality Gate (Optional)

After cleanup, verify no regressions were introduced.

Check if the project has a test infrastructure:

```bash
just --evaluate 2>/dev/null && grep -q "^test" justfile 2>/dev/null
```

**If `just test` is available**: Run it.

```bash
just test
```

If tests fail, the cleanup introduced a regression:
1. Report the failure
2. Revert the changes that caused the failure
3. Re-run tests to confirm the revert fixes the issue
4. Continue to summary with a note about the reverted changes

**If `just test` is not available**: Skip the quality gate.

Output one of:
- `Step 3/4: Quality gate... DONE (tests passed)`
- `Step 3/4: Quality gate... SKIPPED (no just test available)`
- `Step 3/4: Quality gate... DONE (N regressions reverted)`

## Step 4: Cleanup Summary

Output a summary of what was done:

```
## Cleanup Summary

**Scope**: N files (git diff against main)
**Modified**: M files
**Skipped**: K files (already clean)
**Quality gate**: passed / skipped / N regressions reverted

### Changes by Type
- Dead code removed: N instances
- Complexity reduced: N instances
- Naming improved: N instances
- Duplication eliminated: N instances
- Other: N instances

### Files Modified
- path/to/file1.ext — <brief description of changes>
- path/to/file2.ext — <brief description of changes>
```

If invoked as a standalone command (not via pipeline task), the summary is the final output.

If invoked via pipeline task (`T-clean-code-1`), invoke the skill after the summary:

```
Skill(skill="forge:submit-task")
```
