---
name: run-tests
description: Execute test scripts and generate a results report. Pure executor: reads commands from .forge/config.yaml test.execution node. Convention-driven result parsing.
---

# Run Tests

Execute test scripts and generate a test results report.

**Core principle**: A pure executor that reads execution commands from `.forge/config.yaml` and runs them. Does three things: execute configured commands, parse results, generate report.

<HARD-GATE>
This skill only executes configured test commands and reports results. Forbidden:
- Modifying test script content
- Skipping failed tests (must report faithfully)
- "Fixing" tests during execution to make them pass
- Using any hardcoded command names -- all commands come from config
</HARD-GATE>

## When to Use

**Trigger:**
- User asks to "run tests"
- User provides `/run-tests` command
- After `/gen-test-scripts` has generated test scripts

**Skip:**
- `tests/<journey>/` doesn't exist (run `/gen-test-scripts` first)

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `.forge/config.yaml` with `test.execution.run` | See "Missing config" error below |
| `tests/<journey>/` directory (at least one) | Run `/gen-test-scripts` first |

<PRINCIPLE>
**Shared infrastructure first.** Before executing any test actions, verify that shared dependencies are complete and functional. Shared file names are defined in the Convention's Framework section. If shared files are missing symbols imported by test files, all tests will fail at the import stage. When inconsistencies are found, go back to `/gen-test-scripts` to fix shared dependencies before running tests.
</PRINCIPLE>

## Workflow

```
0. Load Convention → 1. Load Config → 2. Validate Output Flags → 3. Setup (optional) → 4. Pre-check (optional) → 5. Run → 6. Parse Results → 7. Generate Report → 8. Teardown (optional)
```

### Step 0: Load Convention Result Format

Scan `docs/conventions/` for files with `domains` containing `testing` and load the **Result Format** section from each matched file.

For each convention file whose frontmatter `domains` includes `testing`, read the file and extract:

1. **format-type**: One of `json-stream`, `json-report`, `text-verbose`
2. **output-flags**: The flags passed to the test runner (e.g., `-json -v`, `--reporter=json`)

<HARD-RULE>
Parsing logic must be driven by Convention Result Format section, not framework name. The format-type determines the parsing strategy; never branch on language or framework identity.
Convention Result Format only provides `format-type` and `Output flags` for parsing, never for execution.
</HARD-RULE>

**Convention merge semantics**: When multiple Convention files match, merge at section level. If two files both declare a Result Format section, the later file's values win. Log a note about the overlap for user awareness.

**Fallback -- no Convention found**: If no Convention files exist in `docs/conventions/`, proceed with `text-verbose` as the default format-type. Use generic text-based parsing: scan output lines for PASS/FAIL/SKIP patterns and extract test names from leading markers.

### Step 1: Load Config

Read test execution configuration from `.forge/config.yaml`:

```bash
forge config get test.execution
```

**Required field**: `test.execution.run` -- the command template to execute tests.

If `test.execution` or `test.execution.run` is missing, abort with:

> **Missing test execution config.** Add the following to `.forge/config.yaml`:
>
> ```yaml
> test:
>   execution:
>     run: "just e2e-test --feature {slug}"   # Required: command template
>     # setup: "just e2e-setup"               # Optional: pre-execution setup
>     # pre-check: "just e2e-verify --feature {slug}"  # Optional: validation before run
>     # teardown: "just e2e-teardown"         # Optional: post-execution cleanup
>     # results-dir: "tests/{journey}/results"  # Optional: results directory
>     # timeout: 300                           # Optional: timeout in seconds (default 600)
> ```

### Template Variables

Resolve template variables in command strings before execution:

| Variable | Source | Default if missing |
|----------|--------|-------------------|
| `{slug}` | `forge feature` | **Error** -- abort with message below |
| `{journey}` | Convention or directory scan | `e2e` |
| `{test-dir}` | Convention Framework | `tests` |
| `{results-dir}` | `test.execution.results-dir` config | `tests/{journey}/results` |

**Escape rule**: `{{var}}` resolves to literal `{var}`.

**Variable resolution order**:
1. Replace `{{` with a temporary sentinel (preserves literal braces)
2. Replace `{slug}`, `{journey}`, `{test-dir}`, `{results-dir}` with resolved values
3. Replace sentinel back to `{`

**Missing slug** error:

> **No active feature slug.** Run `forge feature <slug>` to set the active feature, then retry.

### Step 2: Validate Output Flags

Before executing any commands, verify consistency between Convention and config:

1. Read Convention Result Format's `format-type` and `output-flags`
2. Check `test.execution.run` command for presence of expected output flags
3. If flags are required by format-type but missing from run command, abort:

> Convention declares format-type `json-stream` which requires output flags like `-json`, but `test.execution.run` does not include these flags. Either add the flags to your run command in config, or change Convention's format-type to `text-verbose`.

### Step 3: Setup (Optional)

If `test.execution.setup` is configured, execute it:

```bash
# Template: test.execution.setup (after variable resolution)
# Example: "just e2e-setup"
```

Ensure results directory exists:

```bash
mkdir -p "${results_dir}"
```

### Step 4: Pre-check (Optional)

If `test.execution.pre-check` is configured, execute it:

```bash
# Template: test.execution.pre-check (after variable resolution)
# Example: "just e2e-verify --feature {slug}"
```

If pre-check fails (non-zero exit), abort and report:

> Pre-check command failed. This usually means test scripts have unresolved markers or missing dependencies. Return to `/gen-test-scripts` to resolve issues.

### Step 5: Run Tests

Execute the run command:

```bash
# Template: test.execution.run (after variable resolution)
# Example: "just e2e-test --feature {slug}"
```

Capture the full stdout/stderr output for result parsing in Step 6.

**Timeout**: If `test.execution.timeout` is configured, wrap execution with a timeout. Default timeout is 600 seconds. On timeout, terminate the process and mark all tests as FAIL(timeout).

**State file**: Before execution, write teardown state to `.forge/test-state.json`:

```json
{"teardown": "<resolved teardown command>", "timestamp": "<ISO8601>"}
```

This enables cleanup recovery if the session is interrupted.

### Step 6: Parse Results

Parse test results based on the **format-type** loaded from Convention in Step 0.

**Guard**: Before parsing, verify result output exists and is valid. If result output is missing or empty: report the error with the test runner's console output as evidence, and abort report generation. Do NOT attempt to parse missing/malformed output.

For detailed parsing strategies per format-type, see `rules/result-parsing.md`.

### Step 7: Generate Report

Read the template at `templates/test-report.md`. Fill in:
- Summary statistics (total/pass/fail/skip per type)
- Per-test-case results with evidence
- Failed test details with error messages
- Screenshot paths (for UI tests only)

**Screenshots**: Use `glob ${results_dir}/**/*.png` to discover screenshots. When available, use the `mcp__zai-mcp-server__analyze_image` tool to examine screenshots and add diagnostic notes. Include screenshots section only when screenshot files are found.

Write to: `${results_dir}/latest.md`

### Step 8: Teardown (Optional)

<HARD-RULE>
**Teardown is mandatory when configured**, even if tests fail.
</HARD-RULE>

If `test.execution.teardown` is configured, execute it:

```bash
# Template: test.execution.teardown (after variable resolution)
```

After successful teardown, delete the state file:

```bash
rm -f .forge/test-state.json
```

**Stale state recovery**: On skill startup, check for `.forge/test-state.json`. If it exists from a previous interrupted session, execute the stored teardown command before proceeding with the current run.

## Output

After completion, report to the user:

```
Test Results: X/Y passed (Z failed)

Failed tests:
- TC-NNN: {failure reason}
- TC-NNN: {failure reason}

Report: tests/<journey>/results/latest.md
```

If all tests pass:

```
Test Results: X/X passed
Report: tests/<journey>/results/latest.md
```

## Error Handling

| Situation | Action | Retries |
|-----------|--------|---------|
| `test.execution.run` not configured | Abort with config example | 0 |
| No active feature slug | Abort with `forge feature` prompt | 0 |
| Output flags mismatch (Convention vs config) | Abort with mismatch details | 0 |
| Pre-check command fails | Abort, suggest returning to `/gen-test-scripts` | 0 |
| Setup command fails | Report error, abort | 0 |
| Test timeout | Mark as FAIL with timeout reason | 0 |
| Test file doesn't compile | Report compilation error, skip that file | 0 |
| Convention file has no Result Format section | Fallback to text-verbose parsing, log a note | 0 |
| Result output missing or empty | Report error with console output, abort report generation | 0 |
| Teardown command fails | Log error, leave state file for recovery | 0 |
| Stale state file detected on startup | Execute stored teardown, then proceed | 0 |

## Failure Diagnosis

When tests fail, follow the diagnostic flow in `rules/failure-diagnosis.md`. Key gate:

<HARD-RULE>
When **>30% of tests fail simultaneously**, do NOT proceed to individual test fix tasks. Run app health diagnostics first.
</HARD-RULE>

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-scripts` | Generate executable test scripts |
