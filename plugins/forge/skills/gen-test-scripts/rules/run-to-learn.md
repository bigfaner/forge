# Run-to-Learn (R2L) Mechanism

Run-to-Learn generates skeleton tests, runs them to capture runtime behavior, enriches the Fact Table with confirmed facts, then regenerates more precise tests. R2L is optional and always degrades gracefully.

## Enabling R2L

R2L is disabled by default. Enable via either:

1. `.forge/config.yaml` field: `enabled: true` under `run_to_learn` block
2. CLI flag (when invoking gen-test-scripts): `--run-to-learn`

If neither is set, skip R2L entirely and generate tests using static Fact Table data only.

### Configuration

```yaml
# .forge/config.yaml
run_to_learn:
  enabled: true                     # enable/disable
  max_iterations: 3                 # hard cap (default: 3)
  coverage_threshold: 0.80          # exit threshold (default: 0.80)
  timeout_per_test: "60s"           # per-skeleton timeout (default: 60s)
  skip_on_env_failure: true         # skip when environment not ready (default: true)
```

<HARD-RULE>
If R2L is enabled but environment is not ready (see Environment Check below), skip R2L entirely. Do NOT block the pipeline. Proceed with static Fact Table data.
</HARD-RULE>

## Environment Check

Before entering the R2L loop, verify the execution environment is ready:

1. Read the project's surface type from `.forge/config.yaml` (`forge config get surface`).
2. Check surface-specific prerequisites:
   - **CLI/TUI**: Binary can be compiled (`go build`, `cargo build`, etc.). Build command derived from Fact Table or Convention.
   - **API**: Service is reachable at the configured base URL. Database (if applicable) is accessible.
   - **WebUI**: Dev server is running and responds to health check.
   - **Mobile**: N/A -- Mobile skeleton tests are Maestro YAML and do not execute code.
3. If any prerequisite fails:
   - Log: "R2L environment not ready: <specific missing prerequisite>. Skipping Run-to-Learn."
   - Skip R2L entirely.
   - Proceed with static Fact Table data for test generation.

<HARD-RULE>
R2L environment readiness is a separate concern from the run-tests env-check (task 2.8). R2L checks only what is needed to run skeleton tests. If `skip_on_env_failure` is true (default), skip R2L on any failure. If false, abort with actionable error message.
</HARD-RULE>

## R2L Loop

```
iteration = 0
while iteration < max_iterations AND coverage < coverage_threshold:
    1. Generate Skeleton Tests
    2. Run Skeleton Tests
    3. Process Results -> Write to Fact Table
    4. Compute Coverage
    5. iteration++
```

### Step 1: Generate Skeleton Tests

For each Contract Outcome that lacks a confirmed/runtime fact:

1. Generate a **skeleton test** -- a test function that contains:
   - **Setup**: Create isolated environment (`t.TempDir()` or equivalent), set up required preconditions from the Contract.
   - **Execution**: Invoke the system under test (subprocess, HTTP call, UI interaction).
   - **Output capture**: Capture stdout, stderr, exit code (or HTTP status/response body for API).
   - **No assertions**: The skeleton test does NOT assert on expected values. It only captures actual runtime output.

2. Write skeleton tests to a temporary directory (NOT to `tests/<journey>/`). Use `os.TempDir()` + journey-specific subdirectory.

3. Skeleton tests must respect the Convention's framework patterns (imports, test runner, file naming) so they compile and execute correctly.

<HARD-RULE>
Skeleton tests run in a temporary directory, never in the project's `tests/` directory. They are ephemeral -- created for one R2L iteration, discarded after data extraction.
</HARD-RULE>

#### Timeout Protection

Each skeleton test has a timeout. The timeout value comes from `.forge/config.yaml` `timeout_per_test` (default: 60s).

Apply timeout per the Convention's framework mechanism:
- Go: `ctx, cancel := context.WithTimeout(context.Background(), timeout)` + pass ctx to subprocess
- Python: `@pytest.mark.timeout(seconds)` or `subprocess.run(..., timeout=seconds)`
- JavaScript: `jest.setTimeout(ms)` or `execFile(..., { timeout: ms })`

If a skeleton test exceeds the timeout, treat it as a runtime crash (see Failure Handling below).

### Step 2: Run Skeleton Tests

Execute the skeleton tests in the temporary directory:

```bash
# Framework-specific test runner, scoped to temp directory
just test <skeleton-dir>  # or Convention-specific runner command
```

Capture for each skeleton test:
- **stdout**: The captured output of the system under test
- **stderr**: Error output
- **exit code**: Process exit code (or HTTP status code for API)
- **timing**: Execution duration (for diagnostic purposes)

<HARD-RULE>
For API skeleton tests: only send GET requests. Write operations (POST/PUT/DELETE) MUST NOT be executed. For write operations, generate rollback SQL/statements and record in Fact Table with `side_effect: requires_cleanup`.
</HARD-RULE>

### Step 3: Process Results and Write to Fact Table

For each skeleton test result, write a runtime fact to the Fact Table using `forge fact` CLI or direct Fact Table manipulation:

#### Successful Capture

When the skeleton test captures valid output:

```json
{
  "fact_id": "cli.forge-task-claim-output_format-{nonce}",
  "source": "runtime",
  "subject": "cli.forge task claim",
  "kind": "output_format",
  "value": {
    "stdout": "claimed task 2.7",
    "exit_code": 0,
    "pattern": "claimed task (?P<task_id>[\\w.]+)"
  },
  "confidence": "confirmed",
  "updated_at": "2026-05-23T10:30:00Z"
}
```

#### Runtime-to-Static Replacement Semantics

<HARD-RULE>
Runtime fact replaces the static fact with the same `subject` + `kind`. However, if the runtime fact's `confidence` is NOT `confirmed`, the static fact is preserved as a fallback entry (both entries coexist in the Fact Table, distinguished by `source` field).
</HARD-RULE>

Implementation:

1. For each captured result, construct a `FactEntry` with `source: "runtime"`.
2. If the capture was clean and complete, set `confidence: "confirmed"`.
3. If the capture was partial or uncertain, set `confidence: "inferred"` and keep the existing static entry.
4. Write the entry to `.forge/fact-table.json` via the `facttable` package or `forge fact` CLI.

### Step 4: Compute Coverage

Coverage formula:

```
coverage = (confirmed runtime facts covering Outcomes) / (total Outcomes) * 100%
```

Only `source: runtime` AND `confidence: confirmed` facts count toward the numerator. Static and inferred facts do not.

**Computing coverage**:

1. Count total Outcomes across all Contract files for the current Journey.
2. For each Outcome, check if there exists a runtime+confirmed fact whose `subject` matches the Outcome's action and whose `kind` is relevant (output_format, error_code, side_effect, precondition, signature).
3. An Outcome is "covered" if at least one runtime+confirmed fact matches its action.
4. `coverage = covered_outcomes / total_outcomes * 100`

If `coverage >= coverage_threshold` (default 80%), exit the loop early.

### Step 5: Regenerate Tests

After the R2L loop completes (either by reaching coverage threshold or max iterations):

1. Read the enriched Fact Table (now containing both static and runtime facts).
2. Use `EffectiveEntries()` semantics: runtime confirmed replaces static; runtime non-confirmed keeps static as fallback.
3. Regenerate test code using the enriched Fact Table, following the normal gen-test-scripts pipeline (Steps 1-4 of SKILL.md).

The regenerated tests now have:
- Precise regex patterns derived from actual runtime output
- Confirmed error codes and messages
- Validated state transitions
- Confirmed side effects

## Failure Handling

R2L failures NEVER block the pipeline. Every failure degrades gracefully to static Fact Table data.

### Compilation Failure

**Detection**: Test runner returns non-zero exit code during compilation phase.

**Action**:
1. Record the compilation error to the Fact Table:
   ```json
   {
     "fact_id": "cli.forge-task-claim-compilation_error-{nonce}",
     "source": "runtime",
     "subject": "cli.forge task claim",
     "kind": "compilation_error",
     "value": { "error": "undefined: assertEqual", "file": "skeleton_test.go:42" },
     "confidence": "assumed",
     "updated_at": "..."
   }
   ```
2. Skip this round of R2L for the affected test(s).
3. Continue with remaining skeleton tests.
4. Use existing static facts for the skipped test's Outcomes.

### Runtime Crash

**Detection**: Process killed by signal (SIGSEGV, SIGKILL, etc.) or timeout exceeded.

**Action**:
1. Record the crash to the Fact Table:
   ```json
   {
     "fact_id": "cli.forge-task-claim-runtime_crash-{nonce}",
     "source": "runtime",
     "subject": "cli.forge task claim",
     "kind": "runtime_crash",
     "value": { "signal": "SIGSEGV", "stderr": "panic: ..." },
     "confidence": "assumed",
     "updated_at": "..."
   }
   ```
2. Mark related Outcomes as LOW confidence.
3. Do NOT retry the crashed skeleton test in subsequent iterations.
4. Use existing static facts as fallback.

### Dirty Data Output

**Detection**: Captured output format does not match expected schema (e.g., garbled binary output when text was expected, malformed JSON when structured output was expected).

**Action**:
1. Discard all runtime data from this round.
2. Preserve the Fact Table state from the previous round (or initial static state if first round).
3. Log the mismatch details: expected schema vs actual output.
4. Continue to next iteration (if iterations remain).

### API Write Operation Side Effects

**Detection**: A skeleton test would need to execute a write operation (POST/PUT/DELETE) to capture runtime behavior.

**Action**:
1. Do NOT execute the write operation.
2. Instead, generate a rollback statement for the write operation:
   ```
   -- Rollback for POST /tasks: DELETE /tasks/{id} where id is from response
   ```
3. Record in Fact Table with `side_effect: requires_cleanup`:
   ```json
   {
     "fact_id": "api.POST-tasks-side_effect-{nonce}",
     "source": "runtime",
     "subject": "api.POST /tasks",
     "kind": "side_effect",
     "value": { "method": "POST", "path": "/tasks", "rollback": "DELETE /tasks/{id}", "requires_cleanup": true },
     "confidence": "inferred",
     "updated_at": "..."
   }
   ```
4. Use existing static facts for this endpoint's assertions.

<HARD-RULE>
The fallback principle: any R2L failure results in graceful degradation. The pipeline MUST NOT be blocked. Tests are generated using whatever static information is available. Failed R2L items are logged and flagged in the final test report.
</HARD-RULE>

## Output and Reporting

After R2L completes (or is skipped), include in the test generation output:

1. **R2L Status**: Enabled/Disabled/Skipped
2. **Iterations used**: N of max_iterations
3. **Final coverage**: X% (Y confirmed runtime facts covering Z of W Outcomes)
4. **Failed items**: List of skeleton tests that failed (compilation error, runtime crash, dirty data), with failure type and affected Outcomes
5. **Degraded outcomes**: Outcomes using static fallback (LOW confidence) due to R2L failures

This information is included as comments in the generated test files:

```go
// R2L: enabled, 2/3 iterations, coverage 85%
// R2L-FAILED: cli.forge task submit (runtime_crash) -> using static fallback
```

## Integration with gen-test-scripts Pipeline

R2L integrates into the existing gen-test-scripts SKILL.md pipeline as an optional step between Step 1 (Code Reconnaissance) and Step 3 (Generate Test Code):

```
Step 0: Load Convention
Step 0.5: Surface Detection
Step 1: Code Reconnaissance -> Static Fact Table
Step 1.5: [R2L] Run-to-Learn (if enabled) -> Enriched Fact Table  <-- THIS RULE
Step 2: Read Contract Specifications
Step 2.5: Load Type Rules
Step 3: Generate Test Code (using enriched Fact Table)
Step 4: Compile Gate
```

When R2L is enabled and succeeds, Step 3 uses the enriched Fact Table (runtime facts take precedence over static via EffectiveEntries). When R2L is disabled or fails, Step 3 uses the static Fact Table unchanged.
