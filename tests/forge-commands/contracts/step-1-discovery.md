# Contract: forge-commands / Step 1: Discovery

## Outcome "help-output"
- Preconditions: "forge CLI binary built from current source"
- Input: `forge --help`
- Output: "exit code 0, output contains command groups: task, e2e, forensic, test, prompt and top-level commands: cleanup, feature, probe, quality-gate, verify-task-done"
- State: "no state changes"
- Side-effect: none

## Outcome "task-subcommand-help"
- Preconditions: "forge CLI binary built from current source"
- Input: `forge task --help`
- Output: "exit code 0, output contains subcommands: claim, submit, status, query, check-deps, validate-index, add, index, migrate, list-types"
- State: "no state changes"
- Side-effect: none

## Outcome "unknown-command-error"
- Preconditions: "forge CLI binary built from current source"
- Input: `forge taks` (misspelled command)
- Output: "exit code 1, output contains 'unknown' text"
- State: "no state changes"
- Side-effect: none

## Outcome "unknown-task-subcommand-help"
- Preconditions: "forge CLI binary built from current source"
- Input: `forge task nonexistent-sub`
- Output: "exit code 0 (cobra shows help), output lists valid subcommands"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- forge binary path consistent across all steps
- all commands use built binary, not system-installed
