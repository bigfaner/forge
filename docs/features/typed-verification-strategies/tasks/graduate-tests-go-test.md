---
id: "T-test-4"
title: "Graduate Test Scripts (go-test)"
priority: "P1"
estimated_time: "30min"
dependencies: ["T-test-3"]
type: "test-pipeline.graduate"
scope: "all"
profile: "go-test"
---

# Graduate Test Scripts (go-test)

Profile: **go-test**

# Go Test Graduate Strategy

Profile-specific graduation rules for the `graduate-tests` skill.

## Source File Discovery

| Item | Value |
|------|-------|
| File extension | `_test.go` |
| Source directory | `tests/e2e/features/<slug>/` (staging) |
| Target directory | `tests/e2e/` (regression) |

## Import Rewrite

**None required.** Go uses module paths (defined in `go.mod`) rather than relative file paths. All imports resolve identically regardless of file location within the module.

## Validation

| Check | Command | Failure action |
|-------|---------|----------------|
| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |
| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |
| Test discovery | `just e2e-discover` | Rollback via migration manifest |

## Merge Procedure

Go file-level merge when a target file already exists at the graduation destination:

1. Read both source and target test files
2. Backup target file (only if no backup exists -- prevents overwriting original on re-run)
3. Combine imports, deduplicate by import path
4. Match test functions by name -- `func TestTC_NNN_*`
5. Deduplicate: identical function names keep source version; different names with same TC ID prefix keep both
6. Append new test functions that don't exist in target
7. Write the merged file

Merge strategy is `package` -- all test functions reside in the same `package e2e` declaration.

## Shared Infrastructure

These files already exist at `tests/e2e/` and must NOT be copied or modified during graduation:

- `main_test.go` (TestMain setup/teardown)
- `helpers_test.go` (shared test helpers)
- `testdata/` (golden files, fixtures)

## Compilation Check

After migration, verify all packages compile correctly with the new test files in place:

```bash
just e2e-compile
```

## Test Discovery

Verify all expected tests are discoverable:

```bash
just e2e-discover
```

Output is a plain list of test function names, one per line. Compare against expected TC IDs.

## Graduation Marker

Written only after validation passes (atomic -- no marker = not graduated):

```yaml
schema_version: 1
status: completed
timestamp: <UTC ISO timestamp>
source: tests/e2e/features/<slug>/
targets:
  - tests/e2e/<test-file>
modules:
  - <module-name>
testCount: <N>
```
