---
name: config-schema
description: Schema definition for test.execution configuration in .forge/config.yaml
---

# Test Execution Config Schema

Schema for the `test.execution` node in `.forge/config.yaml`. This config drives the `/run-tests` skill's execution behavior.

## Full Example

```yaml
# .forge/config.yaml
test:
  execution:
    run: "just test {slug}"                          # Required
    setup: "just test-setup"                         # Optional
    pre-check: "just probe"                          # Optional
    teardown: "just test-teardown"                   # Optional
    results-dir: "tests/{journey}/results"            # Optional
    timeout: 300                                      # Optional
```

## Field Reference

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `run` | string | **Yes** | -- | Command template to execute tests. Supports template variables. |
| `setup` | string | No | -- | Pre-execution setup command (e.g., install dependencies, start servers). Runs once before tests. |
| `pre-check` | string | No | -- | Validation command before test execution (e.g., check for unresolved markers). Failure aborts the run. |
| `teardown` | string | No | -- | Post-execution cleanup command. Runs even if tests fail. Guaranteed via state file recovery. |
| `results-dir` | string | No | `tests/{journey}/results` | Directory path for test result output. Supports template variables. |
| `timeout` | integer | No | `600` | Maximum execution time in seconds. On timeout, all tests are marked FAIL(timeout). |

## Template Variables

Variables are resolved in command strings before execution.

| Variable | Source | Required | Default if missing |
|----------|--------|----------|-------------------|
| `{slug}` | `forge feature` | Yes | **Error** -- abort with `forge feature <slug>` prompt |
| `{journey}` | Convention or directory scan | No | `e2e` |
| `{test-dir}` | Convention Framework section | No | `tests` |
| `{results-dir}` | `test.execution.results-dir` config | No | `tests/{journey}/results` |

**Escape rule**: `{{var}}` resolves to literal `{var}`. Use this when a command needs literal curly-brace syntax.

## Config Examples by Project Type

### Playwright E2E (just)

```yaml
test:
  execution:
    run: "just test {slug}"
    setup: "just test-setup"
    pre-check: "just probe"
    teardown: "just test-teardown"
    results-dir: "tests/e2e/results"
    timeout: 300
```

### Go Unit Tests (go test)

```yaml
test:
  execution:
    run: "go test -json -v ./..."
    results-dir: "test-results"
    timeout: 120
```

### Vitest Integration Tests (npx)

```yaml
test:
  execution:
    run: "npx vitest run --reporter=json --outputFile={results-dir}/report.json"
    setup: "npm run build"
    results-dir: "tests/integration/results"
    timeout: 180
```

### Makefile-based Projects

```yaml
test:
  execution:
    run: "make test FEATURE={slug}"
    setup: "make test-setup"
    teardown: "make test-clean"
    timeout: 600
```
