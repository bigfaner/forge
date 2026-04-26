---
description: Pull the latest source branch, create a new branch from it, and switch to it. Optionally accepts a parameter to specify the source branch (defaults to main).
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

2. Ask the user for the new branch name. Suggest a convention like `feat/<slug>`, `fix/<slug>`, or `chore/<slug>` based on the nature of the work.

3. Run the following git commands:
   - `git fetch origin`
   - `git checkout <source-branch>`
   - `git pull origin <source-branch>`
   - `git checkout -b <new-branch-name>`

4. Confirm the current branch to the user.
