# Contract: task-lifecycle / Step 3: Quality Gate

## Outcome "gate-pass"
- Preconditions: "all tasks completed, project compiles, fmt, lint, and tests pass"
- Input: `forge quality-gate` with CLAUDE_PROJECT_DIR set to project root
- Output: "quality gate passed confirmation, exit code 0"
- State: "no new tasks created, feature marked as quality-gate-passed"
- Side-effect: "runs compile, fmt, lint, and test commands in sequence"

## Outcome "gate-fail-creates-fix-task"
- Preconditions: "quality gate fails, existing fix task count for the failing step is below max (3)"
- Input: `forge quality-gate` with failing quality checks
- Output: "quality gate failed message with fix task creation confirmation, exit code non-zero"
- State: "new fix task created with type coding.fix and sourceTaskID referencing the failing step"
- Side-effect: "fix task added to index.json"

## Outcome "gate-fail-max-fix-tasks"
- Preconditions: "quality gate fails, 3 existing fix tasks for the same step already exist"
- Input: `forge quality-gate` with failing quality checks and max fix tasks reached
- Output: "error message indicating max fix tasks reached for the step, exit code non-zero"
- State: "no new fix task created"
- Side-effect: none

## Outcome "cleanup-terminal-tasks"
- Preconditions: "feature has tasks in terminal states (completed, blocked, skipped, rejected)"
- Input: `forge cleanup` with CLAUDE_PROJECT_DIR set
- Output: "cleanup confirmation listing removed state files, exit code 0"
- State: "state.json files for terminal tasks removed"
- Side-effect: none

## Outcome "cleanup-no-terminal-tasks"
- Preconditions: "feature has no tasks in terminal states"
- Input: `forge cleanup` with CLAUDE_PROJECT_DIR set
- Output: "message indicating no terminal tasks to clean up, exit code 0"
- State: "no state changes"
- Side-effect: none

## Outcome "stage-gates-generated"
- Preconditions: "feature has phases with 2+ business tasks, stage gates not yet generated"
- Input: `forge task index --feature <slug>` in project with qualifying phases
- Output: "INDEX_BUILT action with task count including gate/summary tasks"
- State: "summary.md and gate.md files created for qualifying phases, index.json updated"
- Side-effect: none

## Outcome "stage-gates-idempotent"
- Preconditions: "stage gates already generated for a feature"
- Input: `forge task index --feature <slug>` re-run
- Output: "identical output to previous run, exit code 0"
- State: "existing gate/summary files unchanged (content and modification time)"
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
- task_id stable once assigned
- index.json remains valid JSON throughout
