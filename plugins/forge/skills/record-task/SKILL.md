---
name: record-task
description: Record task execution results and update task status. Delegates to submit-task for the full record workflow.
---

# Record Task

Record task completion metrics and update status. This skill is the metrics-collection front-end for the submit-task workflow.

## Metrics Collection (MANDATORY)

Before recording any task, collect real metrics from the project's test runner:

```
just test [scope]
```

<HARD-RULE>
All numeric fields (`coverage`, `testsPassed`, `testsFailed`) must come from actual `just test` output, never guessed or defaulted. The `just test` command dispatches to the correct toolchain automatically.
</HARD-RULE>

## Workflow

1. Run `just test [scope]` to collect real test metrics
2. Parse output for: tests passed, tests failed, coverage percentage
3. Write metrics to `process/record.json`
4. Run `forge task submit <TASK_ID> --data process/record.json`

See `plugins/forge/skills/submit-task/SKILL.md` for the full record submission workflow and validation rules.
