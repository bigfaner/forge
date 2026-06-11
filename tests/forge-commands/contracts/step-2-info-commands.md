# Contract: forge-commands / Step 2: Info Commands

## Outcome "config-get-project-type"
- Preconditions: "forge CLI with valid .forge/config.yaml containing project-type: backend"
- Input: `forge config get project-type`
- Output: "exit code 0, plain text 'backend' without formatting"
- State: "no state changes"
- Side-effect: none

## Outcome "config-get-interfaces"
- Preconditions: "forge CLI with valid .forge/config.yaml containing interfaces array"
- Input: `forge config get interfaces`
- Output: "exit code 0, at least 3 lines, each line is plain interface name without quotes"
- State: "no state changes"
- Side-effect: none

## Outcome "config-get-missing-key"
- Preconditions: "forge CLI with valid .forge/config.yaml"
- Input: `forge config get nonexistent-key`
- Output: "exit code 1, no stdout output"
- State: "no state changes"
- Side-effect: none

## Outcome "proposal-list"
- Preconditions: "forge CLI with docs/proposals/ containing at least one proposal"
- Input: `forge proposal`
- Output: "exit code 0, table with SLUG, CREATED, STATUS, PRD, FEATURE columns"
- State: "no state changes"
- Side-effect: none

## Outcome "proposal-detail"
- Preconditions: "forge CLI with docs/proposals/forge-info-commands/"
- Input: `forge proposal forge-info-commands`
- Output: "exit code 0, detail view with SLUG, CREATED, STATUS, FILE fields"
- State: "no state changes"
- Side-effect: none

## Outcome "feature-list"
- Preconditions: "forge CLI with at least one feature in docs/features/"
- Input: `forge feature list`
- Output: "exit code 0, table with SLUG, STATUS, PROGRESS columns"
- State: "no state changes"
- Side-effect: none

## Outcome "feature-status-detail"
- Preconditions: "forge CLI with feature forge-info-commands"
- Input: `forge feature status forge-info-commands`
- Output: "exit code 0, detail view with SLUG, STATUS, TASKS fields"
- State: "no state changes"
- Side-effect: none

## Outcome "lesson-list"
- Preconditions: "forge CLI with docs/lessons/ containing at least one lesson"
- Input: `forge lesson`
- Output: "exit code 0, table with NAME, CREATED, CATEGORY, TAGS columns"
- State: "no state changes"
- Side-effect: none

## Outcome "lesson-detail"
- Preconditions: "forge CLI with valid lesson in docs/lessons/"
- Input: `forge lesson <name>`
- Output: "exit code 0, detail view with NAME, FILE fields; no full markdown content"
- State: "no state changes"
- Side-effect: none

## Outcome "init-creates-forge-dir"
- Preconditions: "clean project state (no .forge/ directory)"
- Input: `forge init`
- Output: "exit code 0, .forge/ directory created"
- State: ".forge/ directory and config files created"
- Side-effect: filesystem modifications

## Outcome "justfile-no-project-type-recipe"
- Preconditions: "project with justfile"
- Input: "read justfile content"
- Output: "no 'project-type:' recipe found in justfile"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- forge binary path consistent across all steps
- all commands use built binary, not system-installed
