---
created: "2026-05-19"
tags: [architecture, testing]
---

# Subagent gets stuck on large constant rename tasks

## Problem
Subagent dispatched for task 1 (rename ~20 type constants + update all callers) got stuck after completing the core files (types.go, build.go) but before updating 17 other files with 52 remaining references to old constant names. The agent fell into a whack-a-mole loop: fix one compilation error → check → find another file → fix → repeat.

## Root Cause
1. Task scope was "rename constants AND update all callers in one commit" — 17+ files affected
2. After renaming constants in types.go, every other Go file referencing them produced compilation errors
3. The subagent tried to fix errors file-by-file but each fix surfaced new errors in other files
4. Context window filled with diagnostic noise (20+ undefined errors) before the subagent could complete all updates
5. No pre-computed file list was given — the subagent had to discover affected files via compilation errors one at a time

## Solution
For large mechanical renames across many files:
1. **Pre-compute the full file list**: `grep -rl OLD_CONSTANT forge-cli/ --include="*.go"` before starting
2. **Use bulk replace**: `sed -i` or editor replace_all for deterministic rename mapping
3. **Verify once at the end**: compile/test only after ALL replacements are done, not after each file
4. **Consider splitting**: if >5 files need changes, split into "core rename" + "caller update" tasks

## Reusable Pattern
When a task requires renaming shared constants across many files, do NOT let the subagent discover affected files via compilation errors. Instead:
- Provide the full affected file list in the task description
- Or use a scripted bulk-replace approach
- The task definition should say "use sed/find-replace on these specific files" rather than "update all callers"
- Never compile/test between individual file edits in a bulk rename — only verify at the end

## Example
```bash
# Pre-compute affected files
grep -rl "TypeFeature\|TypeEnhancement\|..." forge-cli/ --include="*.go"

# Bulk replace with sed
sed -i '' 's/TypeFeature/TypeCodingFeature/g' file1.go file2.go ...
# Verify once
go build ./...
```

## Related Files
- `docs/features/task-type-id-redesign/tasks/1-rename-type-constants.md`
- `forge-cli/pkg/task/types.go`
- `forge-cli/pkg/task/build.go`
