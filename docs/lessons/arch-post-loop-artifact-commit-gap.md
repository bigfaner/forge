---
created: "2026-05-20"
tags: [architecture, testing]
---

# Post-loop artifacts left uncommitted after run-tasks completes

## Problem

After `/run-tasks` completes all tasks and finishes knowledge extraction, uncommitted artifacts remain in the working tree: knowledge entries (docs/decisions/, docs/business-rules/), fix task records (records/*.md), and updated index.json. The user had to manually run `/git-commit` to capture these leftovers.

## Root Cause

1. **Surface**: After the run-tasks loop ends and knowledge extraction writes entries, no commit step follows.
2. **Mechanism**: The `/run-tasks` skill's post-completion flow does knowledge extraction → presents to user → ends. There is no "commit all generated artifacts" step after knowledge is confirmed. Fix task records created by `forge task submit` are also not auto-committed.
3. **Structural**: The pipeline assumes each sub-agent commits its own task artifacts (via `forge task submit`), but knowledge extraction and fix task record updates are done in the main session without a commit gate.

## Solution

After knowledge extraction completes and the user confirms entries, commit all remaining artifacts:

```bash
git add docs/decisions/ docs/business-rules/ docs/lessons/ docs/conventions/ docs/features/<slug>/tasks/index.json docs/features/<slug>/tasks/records/
git commit -m "chore(<slug>): commit post-loop artifacts (knowledge + records)"
```

Alternatively, add a "commit artifacts" step to the `/run-tasks` post-completion flow.

## Reusable Pattern

When the run-tasks loop ends (including fix tasks triggered by stop hooks), check for uncommitted artifacts before reporting completion. Common sources:
- Knowledge entries from the extraction step (decisions, lessons, conventions, business rules)
- Fix task records (records/*.md) and updated index.json from `forge task submit`
- Manifest status updates

Run `git status` after knowledge extraction to catch anything missed.

## Related Files

- plugins/forge/skills/run-tasks/SKILL.md (post-completion flow)
- plugins/forge/skills/submit-task/SKILL.md (submit behavior)
