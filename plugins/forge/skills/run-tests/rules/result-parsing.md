---
name: result-parsing
description: Parsing strategies for e2e test result formats (json-stream, json-report, text-verbose)
---

# Result Parsing Strategies

Parsing logic is driven by Convention Result Format section's **format-type**, not framework name.

## Format: json-stream

Each line is an independent JSON object representing a test event. Process line-by-line:

1. **Event types**: `run` (test started), `pass` (test passed), `fail` (test failed), `skip` (test skipped), `output` (captured stdout/stderr)
2. **Grouping**: Group events by test name. Each test starts with a `run` event and ends with a `pass`, `fail`, or `skip` event
3. **Output collection**: Concatenate all `output` events for a test to build its log
4. **Duration**: Use the `Elapsed` field from the terminal event (`pass`/`fail`/`skip`)
5. **TC ID extraction**: Extract from test name using pattern `TC[_-](\d+)`, normalize to `TC-NNN`

**Common fields** (adapt field names to actual JSON structure):
- Test name: `Test` or `name`
- Status: `Action` or `status` (`pass`/`fail`/`skip`)
- Duration: `Elapsed` or `duration`
- Output: `Output` or `message`

## Format: json-report

A single JSON document containing all test results. Parse the complete structure:

1. **Structure**: Typically a tree of suite -> test cases with nested results
2. **Traversal**: Walk the suite hierarchy, collecting each leaf test case
3. **Status mapping**: Map the report's status field to pass/fail/skip
4. **TC ID extraction**: Extract from test name using pattern `TC[_-](\d+)`, normalize to `TC-NNN`

**Common fields** (adapt field names to actual JSON structure):
- Suite container: `suites` or `results`
- Test name: `name` or `title`
- Status: `status` or `outcome` (map to pass/fail/skip)
- Duration: `duration` or `time`
- Error: `error` or `message`

## Format: text-verbose

Plain text output from verbose test runners. Parse using line-by-line scanning:

1. **Test start**: Lines matching patterns like `=== RUN`, `running test:`, `PASS:`, `FAIL:`, `SKIP:`, or `ok ` / `FAIL `
2. **Test end**: Lines with pass/fail/skip indicators
3. **Duration**: Extract from trailing duration patterns like `(0.01s)`, `in 0.01s`, `[0.01s]`
4. **TC ID extraction**: Extract from test name using pattern `TC[_-](\d+)`, normalize to `TC-NNN`
5. **Error collection**: Lines between a FAIL marker and the next test start or summary separator are error output

**Generic fallback pattern**: When Convention is absent and text-verbose is assumed, scan for:
- `PASS` or `ok` lines -> passing test
- `FAIL` or `FAILED` lines -> failing test
- `SKIP` or `SKIPPED` lines -> skipped test

## Format-agnostic rules

Regardless of format-type, these rules always apply:

- **TC ID extraction**: Pattern `TC[_-](\d+)` from test name, normalize separator to hyphen -> `TC-NNN`
- **Test type classification**: Infer from test name or content. UI tests reference pages/elements, API tests reference endpoints, CLI tests reference commands. When uncertain, classify as the dominant type for the project.
- **Error messages**: Always capture the full error text, not just the first line
