# Surface: cli — CLI Functional Test Orchestration

This rule file defines the CLI functional test orchestration sequence for the cli surface in the run-tests skill. The consumer is the SKILL.md dispatcher.

## Orchestration Sequence

| Step | just Recipe | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|------------|-------------|-------------|-------------|-------------|
| test | `just <recipe-prefix>-test <journey>` | CLI functional tests passed | CLI functional tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | `just <recipe-prefix>-teardown` | Cleanup complete | Cleanup failed | — | End |

Notes:
- **No dev step**: CLI surface does not start a persistent service
- **No probe step**: CLI tools do not require HTTP health checks
- **No aggregate recipe**: CLI surface does not execute a `just cli` aggregate recipe

## Failure Handling

### test failure

- Exit code 1: Execute teardown, exit with exit code 1
- Exit code 2 (retryable): Execute teardown, prompt the user "Test environment error, consider retrying", exit with exit code 2

### teardown failure

When teardown fails, log the error and preserve `.forge/test-state.json` for recovery. Exit with the current step's exit code.

## Suite Name

Test report suite names use the `cli-functional/<journey-name>` format.

## Journey Filter

| Tag | Match Rule |
|-----|-----------|
| `@cli` | Exact match |

## Per-Journey Execution

The test step for CLI surface executes per journey. Use the `recipe-prefix` determined in SKILL.md Step 1 (for single-surface projects, the surface-type "cli"; for multi-surface projects, the surface-key) to construct recipe names:

```
for each journey in JOURNEYS:
    just <recipe-prefix>-test <journey>
    record results
    on failure: just <recipe-prefix>-teardown, exit
just <recipe-prefix>-teardown
```

The test recipe invocation format is `just <recipe-prefix>-test <journey>`, where `<journey>` is a directory name discovered from `docs/features/<slug>/testing/`. `<recipe-prefix>` is "cli" for single-surface projects, or the corresponding surface-key for multi-surface projects.
