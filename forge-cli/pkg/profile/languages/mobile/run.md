# Maestro Run Strategy

Profile-specific execution and result-parsing rules for the `run-e2e-tests` skill.

## Execution

| Item | Value |
|------|-------|
| Command | `maestro test <flow-file-or-dir> --format junit --output tests/e2e/results/` |
| Invoked via | `just test-e2e --feature <slug>` |
| Setup | Verify maestro CLI installed, device connected |
| Teardown | Maestro handles app state reset between flows |

## Result Format

JUnit XML output (Maestro writes JUnit format):

```xml
<testsuite name="maestro" tests="3" failures="1" errors="0">
  <testcase name="TC-001 Description" classname="e2e" time="12.5"/>
  <testcase name="TC-002 Failed" classname="e2e" time="5.2">
    <failure message="Assertion failed">...</failure>
  </testcase>
</testsuite>
```

## Result Parsing Rules

| Data | Source | Notes |
|------|--------|-------|
| TC ID | `testcase name` via regex `TC-\d+` | From flow name |
| Status | `testcase` with/without `<failure>` | No `<failure>` = passed |
| Duration | `time` attribute | Seconds |
| Error | `<failure message>` | Failure description |

## Test Type Classification

| Type | Indicators |
|------|------------|
| Mobile UI | Uses `tapOn`, `swipe`, `assertVisible` |
| API | Limited — Maestro primarily for mobile UI |

## Error Handling

| Condition | Action |
|-----------|--------|
| Maestro CLI not found | Report error, prompt to install |
| No device connected | Report error, prompt to start emulator |
| App not installed | Report error, prompt to build and install |
