---
name: auto-task-main-instructions
status: Approved
created: 2026-05-22
---

# Proposal: Auto-Generated Tasks Missing Main Instructions

## Problem

CLI auto-generated tasks (fix tasks, cleanup tasks, test pipeline tasks, and non-template tasks created via `forge task add`) lack a `## Main Instructions` section. The prompt templates (`forge prompt get-by-task-id`) provide generic workflow guidance (e.g., "Read task → Fix → Verify"), but the task files themselves contain no task-specific workflow instructions. Agents must infer execution steps from sparse descriptions like "Execute this test pipeline task."

## Solution

Add `## Main Instructions` section to all CLI auto-generated task files, with content dynamically filled based on task type and context. Three code paths are affected:

1. **Template-based tasks** (`fix-task.md`, `cleanup-task.md`): Add `{{MAIN_INSTRUCTIONS}}` placeholder to templates; callers pass specific instructions
2. **Non-template tasks** (`buildTaskMarkdown()`): Generate default Main Instructions from Description field
3. **Auto-gen test tasks** (`GenerateTestTaskMD()`): Generate type-specific workflow guidance based on StrategyKind and Type

## Alternatives

1. **Enhance prompt templates only** — No task file changes. Cannot personalize per task instance. Rejected because user explicitly wants task-level customization.
2. **Read from convention files** — Autogen tasks load instructions from `docs/conventions/`. Adds coupling to convention file structure. Rejected as over-engineering for this scope.

## Scope

### In Scope

- Add `## Main Instructions` with `{{MAIN_INSTRUCTIONS}}` placeholder to `forge-cli/pkg/template/data/fix-task.md`
- Add `## Main Instructions` with `{{MAIN_INSTRUCTIONS}}` placeholder to `forge-cli/pkg/template/data/cleanup-task.md`
- Add `## Main Instructions` generation to `buildTaskMarkdown()` in `forge-cli/pkg/task/add.go`
- Add `## Main Instructions` generation to `GenerateTestTaskMD()` in `forge-cli/pkg/task/autogen.go`
- Update callers (`addFixTask()` in quality_gate.go) to provide MAIN_INSTRUCTIONS variable

### Out of Scope

- Prompt templates (`forge-cli/pkg/prompt/data/`) — already have workflow guidance
- Skill-level templates (`plugins/forge/skills/quick-tasks/templates/`, `plugins/forge/skills/breakdown-tasks/templates/`) — managed by skills, not CLI
- Main Session Instructions section — separate concept for `mainSession: true` tasks

## Risks

| Risk | Mitigation |
|------|-----------|
| Placeholder unfilled → empty section in task file | ApplyVars already errors on unfilled placeholders; add MAIN_INSTRUCTIONS to built-in vars with sensible defaults |
| Duplicate guidance between prompt template and Main Instructions | Main Instructions provides task-SPECIFIC guidance; prompt template provides GENERIC workflow — they complement, not duplicate |
| Test task instructions too generic | Use StrategyKind and Type to generate targeted instructions per task type |

## Success Criteria

- [ ] All auto-generated task files (fix, cleanup, test pipeline, non-template) contain `## Main Instructions` section
- [ ] Main Instructions content is non-empty and provides actionable workflow guidance
- [ ] Existing tests pass without modification (backward compatible)
- [ ] Template placeholder substitution preserves existing behavior when MAIN_INSTRUCTIONS is provided
