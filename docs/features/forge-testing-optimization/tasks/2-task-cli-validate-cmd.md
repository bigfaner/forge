---
id: "2"
title: "Add task-cli validate-specs command"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
status: pending
breaking: false
noTest: false
mainSession: false
---

# 2: Add task-cli validate-specs command

## Description

Add a `validate-specs` command to task-cli that spawns the Node.js validation script created in Task 1. The command discovers spec files in the feature's e2e test directory and runs `validate-specs.mjs` against them, returning structured pass/fail output.

## Reference Files
- `docs/proposals/forge-testing-optimization/proposal.md` — Source proposal (Phase 2, Section 2.2)
- `plugins/forge/skills/gen-test-scripts/templates/validate-specs.mjs` — Validation script from Task 1
- `task-cli/internal/cmd/validate.go` — Existing validate command for pattern reference
- `task-cli/internal/cmd/root.go` — Command registration

## Affected Files

### Create
| File | Description |
|------|-------------|
| `task-cli/internal/cmd/validate_specs.go` | validate-specs command implementation |
| `task-cli/internal/cmd/validate_specs_test.go` | Unit tests for validate-specs command |

### Modify
| File | Changes |
|------|---------|
| `task-cli/internal/cmd/root.go` | Register new `validate-specs` command |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `task validate-specs` command executes the validation script against spec files
- [ ] Command discovers spec files from `tests/e2e/features/<slug>/` based on current feature context
- [ ] Structured output: prints validation results (errors/warnings) to stdout
- [ ] Exit code: 0 if no errors, 1 if errors found, 2 if script fails to run
- [ ] Unit tests cover: spec discovery, output parsing, error handling
- [ ] Works on Windows (path separator handling in `exec.Command`)

## Implementation Notes

1. **Command pattern**: Follow the existing `validate.go` pattern — read feature context, discover files, spawn subprocess
2. **Node resolution**: Use `node` from PATH (don't hardcode absolute path) — per the proposal's Windows mitigation
3. **Spec discovery**: Glob `tests/e2e/features/<slug>/**/*.spec.ts` to find all spec files
4. **test-cases.md path**: Auto-detect from feature context: `docs/features/<slug>/testing/test-cases.md`. Pass to script via `--test-cases` flag
5. **Output parsing**: Parse JSON output from validate-specs.mjs. Print human-readable summary. Return exit code from script
6. **Error handling**: If Node.js or ts-morph is not available, report a WARNING (not ERROR) and exit 0 — graceful degradation per proposal risk mitigation
