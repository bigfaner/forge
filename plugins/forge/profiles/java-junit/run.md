# JUnit 5 Run Strategy

Profile-specific execution and result-parsing rules for the `run-e2e-tests` skill.

## Execution

| Item | Value |
|------|-------|
| Command | `mvn test -pl tests/e2e -Dtest=*E2E` |
| Invoked via | `just test-e2e --feature <slug>` |
| Setup | `mvn compile` (verify project compiles) |
| Teardown | JVM exits after test run; no persistent processes |
| Report format | Maven Surefire XML |

## Result Format

Maven Surefire XML reports in `tests/e2e/target/surefire-reports/`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="com.example.UserE2E" tests="3" failures="1" errors="0" skipped="0" time="2.5">
  <testcase name="testLoginWithValidCredentials" classname="com.example.UserE2E" time="0.8">
    <!-- passed: no child elements -->
  </testcase>
  <testcase name="testCreateResource" classname="com.example.UserE2E" time="1.2">
    <failure message="expected: 201 but was: 500" type="org.opentest4j.AssertionFailedError">
      Stack trace here
    </failure>
  </testcase>
  <testcase name="testDeleteResource" classname="com.example.UserE2E" time="0.5">
    <skipped message=" precondition failed"/>
  </testcase>
</testsuite>
```

## Result Parsing Rules

### Suite Traversal

Glob `tests/e2e/target/surefire-reports/TEST-*.xml`, parse each `<testsuite>`.

### Field Mapping

| Data | Source | Notes |
|------|--------|-------|
| TC ID | `testcase/@name` via regex `TC-\d+` or method name `testTC_NNN_*` | Extract from `@DisplayName` if present, else method name |
| Status | Child element of `<testcase>`: none = passed, `<failure>` = failed, `<error>` = error, `<skipped>` = skipped | |
| Duration | `testcase/@time` | Seconds (float) |
| Error message | `failure/@message` or `error/@message` | |
| Stack trace | Text content of `<failure>` or `<error>` | |

## TC ID Extraction

Two strategies, in priority order:

1. **DisplayName**: Parse `@DisplayName("TC-NNN: ...")` from source or surefire output
2. **Method name**: Extract from `testTC_NNN_description` pattern

## Error Handling

| Condition | Action |
|-----------|--------|
| Maven not installed | Report error with install instructions, abort |
| Compilation failure | Report compile error, skip all tests |
| Test dependency missing | Report missing dependency, suggest `mvn dependency:resolve` |
| Test timeout | Surefire reports `<error>` with timeout message; mark FAIL |
| No test classes found | Report discovery issue, suggest checking file naming |

## Failure Diagnosis

| Failure rate | Response |
|--------------|----------|
| >30% tests fail | Infrastructure problem — check compilation, dependencies, service availability |
| 10-30% | Spot-check error messages for patterns |
| <10% | Per-test fix tasks |
