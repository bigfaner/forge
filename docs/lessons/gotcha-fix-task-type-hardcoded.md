---
created: "2026-05-28"
tags: [architecture, error-handling]
---

# Fix task type hardcoded to coding.fix regardless of source task type

## Problem

When the `/run-tasks` dispatcher encounters a blocked task (review AC failure, agent timeout, etc.), it creates fix tasks with `--type coding.fix` regardless of whether the source task is a `doc` or `coding` type. For example, a doc review task (T-review-doc) with AC failures in SKILL.md files spawned `coding.fix` tasks instead of `doc.fix` tasks.

This causes cascading failures:
1. Doc file fixes get miscategorized as coding fixes
2. The quality gate for `coding.fix` runs code-level checks (golangci-lint, go test) that are irrelevant for markdown edits
3. Pre-existing code test failures block doc-only submissions, creating an infinite chain of fix tasks

## Root Cause

- Level 1: All `forge task add` calls in `/run-tasks` error handling hardcode `--type coding.fix`
- Level 2: The dispatcher protocol has no mechanism to derive fix type from the source task's type/category
- Level 3: The skill instructions template was designed with only coding tasks in mind — doc tasks were added later but error recovery wasn't updated accordingly

## Solution

The dispatcher should derive the fix task type from the source task:
- Source task type `doc` or `doc.review` → `--type doc.fix`
- Source task type `coding` → `--type coding.fix`
- Default to source task type if no clear mapping exists

The `forge task add` command in the error handling sections of `/run-tasks` should use the claimed task's `TYPE` and `TASK_CATEGORY` fields (already available from `forge task claim` output) to determine the correct fix type.

## Reusable Pattern

When designing error recovery paths, always propagate the source task's type/category to the fix task. A one-size-fits-all fix type creates false coupling between unrelated quality gates.

## Related Files

- plugins/forge/skills/run-tasks/SKILL.md (dispatcher instructions)
- plugins/forge/commands/run-tasks/COMMAND.md (command-level instructions)
