---
id: "3"
title: "Implement copy-files mechanism for worktree creation"
priority: "P0"
estimated_time: "45m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 3: Implement copy-files mechanism for worktree creation

## Description
Implement the copy-files logic that copies configured files from the project root to the new worktree after `git worktree add`. Add path validation to reject absolute paths and `..` traversals. Pre-validate all copy-files exist before creating the worktree.

## Reference Files
- `docs/proposals/worktree-source-branch/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — `runWorktreeStart` where copy logic integrates
- `forge-cli/pkg/profile/config.go` — ForgeConfig with WorktreeConfig.CopyFiles

## Acceptance Criteria
- [ ] `worktree.copy-files: [.env]` in config copies `.env` from project root to worktree after creation
- [ ] Multiple files in `copy-files` all get copied
- [ ] Worktree creation aborts with clear error if any copy-file is missing from project root
- [ ] Copy-files rejects absolute paths (e.g., `/etc/passwd`)
- [ ] Copy-files rejects `..` traversals (e.g., `../../etc/passwd`)
- [ ] If a file already exists in the worktree checkout, it is overwritten with the project root version
- [ ] No copy operation happens when `worktree` config is absent or `copy-files` is empty

## Hard Rules
- Path validation MUST happen BEFORE `git worktree add` to avoid leaving orphan worktrees
- Use `filepath.IsAbs()` and `strings.Contains(path, "..")` for validation
- Copy via `os.ReadFile` / `os.WriteFile` (not `os.Link`) for cross-platform support

## Implementation Notes
- Copy-files apply regardless of whether `--source-branch` is used — the two features are independent
- If copy fails after worktree creation, log the error but do NOT remove the worktree (it's still usable, just missing env files)
- Consider adding a `validateCopyFilePath(relPath string) error` helper for the path validation logic
