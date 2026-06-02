# Surface: tui — Terminal Functional Test Orchestration

This rule file defines the terminal functional test orchestration sequence for the tui surface in the run-tests skill. The consumer is the SKILL.md dispatcher.

## Orchestration Sequence

| Step | just Recipe | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|------------|-------------|-------------|-------------|-------------|
| test | `just <recipe-prefix>-test <journey>` | Terminal functional tests passed | Terminal functional tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | `just <recipe-prefix>-teardown` | Cleanup complete | Cleanup failed | — | End |

Notes:
- **No dev step**: TUI surface does not start a persistent service
- **No probe step**: TUI applications do not require HTTP health checks
- **No aggregate recipe**: TUI surface does not execute a `just tui` aggregate recipe

## Failure Handling

### test failure

- Exit code 1: Execute teardown, exit with exit code 1
- Exit code 2 (retryable): Execute teardown, prompt the user "Test environment error, consider retrying", exit with exit code 2

### teardown failure

When teardown fails, log the error and preserve `.forge/test-state.json` for recovery. Exit with the current step's exit code.

## Suite Name

Test report suite names use the `tui-functional/<journey-name>` format.

## Journey Filter

| Tag | Match Rule |
|-----|-----------|
| `@tui` | Exact match |

## Per-Journey Execution

The test step for TUI surface executes per journey. Use the `recipe-prefix` determined in SKILL.md Step 1 (for single-surface projects, the surface-type "tui"; for multi-surface projects, the surface-key) to construct recipe names:

```
for each journey in JOURNEYS:
    just <recipe-prefix>-test <journey>
    record results
    on failure: just <recipe-prefix>-teardown, exit
just <recipe-prefix>-teardown
```

The test recipe invocation format is `just <recipe-prefix>-test <journey>`, where `<journey>` is a directory name discovered from `docs/features/<slug>/testing/`. `<recipe-prefix>` is "tui" for single-surface projects, or the corresponding surface-key for multi-surface projects.
