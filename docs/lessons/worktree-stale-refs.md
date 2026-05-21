---
name: worktree-stale-ref-guard
description: forge worktree commands should guard against stale remote tracking refs and respect --source-branch when reusing branches
metadata:
  type: feedback
---

## Lesson: `forge worktree start` ignores `--source-branch` when stale remote tracking refs exist

**What happened:** Deleted a remote branch on GitHub + local branch, then ran `forge worktree start <slug> --source-branch v3.0.0`. The worktree was created from a stale `origin/<slug>` ref (pointing to an old commit) instead of v3.0.0. The `--source-branch` flag was silently ignored because the command prioritized reusing existing refs.

**Root cause:** Git does not auto-prune remote tracking refs when branches are deleted remotely. `forge worktree start` checks for existing refs (including stale `origin/` refs) before honoring `--source-branch`. Stale refs are zombies — locally cached but pointing to deleted remote branches.

**Why:** Two separate issues compound:
1. Git's design: `git push --delete` removes the remote ref but leaves `refs/remotes/origin/<branch>` locally
2. `forge worktree start` branch resolution: prefers existing refs over `--source-branch` without checking if the ref is stale or if the user explicitly requested a different source

### How to apply: forge worktree 命令改进建议

1. **`forge worktree start`** — 当用户提供 `--source-branch` 时，应始终使用该分支作为 base，忽略任何已有的同名追踪引用。若用户未指定 `--source-branch` 且发现同名 ref，提示用户选择：复用已有分支 或 从最新 HEAD 创建。

2. **`forge worktree remove --hard`** — 已有 `--hard` flag 用于删除本地分支。改进点：`--hard` 删除本地分支后，应自动执行 `git remote prune origin` 清理关联的 stale tracking refs（或在输出中提示用户手动 prune）。

3. **`forge worktree start`** — 创建 worktree 后，验证 HEAD 是否与 `--source-branch`（或默认 base）一致。若不一致（如被 stale ref 覆盖），输出 warning 并提供 `git reset --hard <source-branch>` 修复建议。

4. **防御性提示** — 当检测到同名 stale ref 存在时，输出 `"Warning: stale remote tracking ref origin/<slug> found. Run 'git remote prune origin' to clean up."` 这样用户至少知道问题存在。

### Workaround (当前)

删除 worktree 和分支后，手动执行 `git remote prune origin`，再创建新 worktree。若已创建但指向错误 commit，在 worktree 内 `git reset --hard v3.0.0`。
