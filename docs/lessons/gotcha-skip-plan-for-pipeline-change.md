---
created: "2026-05-14"
tags: [architecture, local-dev-deployment]
---

# Don't Skip Plan Mode for Forge Pipeline Changes

## Problem

When asked to modify the forge plugin's own pipeline (e.g., adding a task to `/quick`), I jumped straight to editing Go source files without first entering plan mode. This led to:

1. Starting code edits before confirming the design approach
2. No opportunity for the user to review the implementation plan
3. Wasted effort if the approach was wrong

## Root Cause

**Causal chain** (3 levels):

1. **Surface symptom**: Skipped planning phase and went directly to code edits (types.go, testgen.go, etc.)
2. **Direct cause**: Misclassified the request as a "simple code change" rather than a multi-file feature implementation
3. **Root cause**: A one-sentence request ("在 /quick 末尾增加 clean code 的任务") triggered a "just do it" response pattern, bypassing the `EnterPlanMode` guard that should activate for non-trivial, multi-file changes
4. **Trigger condition**: The request was concise and mentioned `/quick`, which I interpreted as a small tweak rather than a pipeline modification spanning Go code, prompt templates, tests, and documentation

## Solution

Before editing code for any request that touches multiple files or modifies forge's own pipeline:

1. **Enter plan mode** (`EnterPlanMode`) to design the approach
2. **Present the plan** for user review before any code changes
3. **Get explicit approval** on scope, file list, and design decisions

## Reusable Pattern

**Rule**: Any request that modifies forge's own skills, commands, agents, or task pipeline is a *feature implementation*, not a quick fix. Apply the same discipline as any multi-file change:

- **3+ files** → plan mode first
- **Cross-module** (pkg + internal + plugins) → plan mode first
- **Pipeline changes** (testgen, infer, prompt templates) → plan mode first

**Why**: Forge pipeline changes have a blast radius across task generation, type inference, prompt synthesis, and documentation. A wrong approach cascades through all these layers. Planning upfront is cheap; reverting scattered edits is expensive.

**How to apply**: When a request mentions `/quick`, `/breakdown-tasks`, or any forge skill/command/agent name, treat it as a feature request and enter plan mode before writing code.

## Example

```
# BAD: Jump to code
User: "在 /quick 末尾增加 clean code 的任务"
→ Start editing types.go, testgen.go, infer.go...

# GOOD: Plan first
User: "在 /quick 末尾增加 clean code 的任务"
→ EnterPlanMode
→ Explore codebase (testgen.go, infer.go, prompt.go, quick-tasks SKILL.md)
→ Present plan: which files to change, new type constant, task placement, template content
→ User approves
→ Implement
```

## Related Files

- `forge-cli/pkg/task/types.go` — Task type constants
- `forge-cli/pkg/task/testgen.go` — Test task generation
- `forge-cli/pkg/task/infer.go` — Type inference from task IDs
- `forge-cli/pkg/prompt/prompt.go` — Prompt template routing
- `plugins/forge/skills/quick-tasks/SKILL.md` — Quick-tasks skill definition

## References

- CLAUDE.md: "直接执行 + 深度交互" — should have challenged the approach before executing
- `EnterPlanMode` tool description: "For UI or frontend changes, start the dev server... For multi-file changes, get user approval on approach"
