---
description: Pull the latest source branch, create a new branch from it, and switch to it. Automatically derives branch name from active feature/proposal. Optionally accepts a parameter to specify the source branch (defaults to main).
argument-hints:
  - name: source-branch
    description: Source branch to pull from (e.g. main, develop). Defaults to main.
    required: false
---

Pull the latest source branch, create a new branch from it, and switch to it. Optionally accepts a parameter $ARGUMENTS to specify the source branch (defaults to main).

Steps:

1. Determine the source branch:
   - If $ARGUMENTS is provided, use it as the source branch.
   - Otherwise, default to `main`.

2. Derive the new branch name automatically:
   - Try `task feature` CLI to get the active feature slug.
   - If unavailable, scan `docs/features/` and `docs/proposals/` for a single active directory.
   - If the conversation is clearly about a specific feature or proposal, extract the slug from context.
   - If a slug is found, suggest `feat/<slug>` as the branch name (use `fix/<slug>` or `chore/<slug>` if the work is a bugfix or chore).
   - If no slug can be determined, ask the user for the branch name.
   - In all cases, present the suggested name and let the user confirm or override before proceeding.

3. Run the following git commands:
   - `git fetch origin`
   - `git checkout <source-branch>`
   - `git pull origin <source-branch>`
   - `git checkout -b <new-branch-name>`

4. Confirm the current branch to the user.
