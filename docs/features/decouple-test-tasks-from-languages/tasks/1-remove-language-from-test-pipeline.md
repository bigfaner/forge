---
id: "1"
title: "Remove language dependency from test pipeline"
priority: "P0"
estimated_time: "2-3h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: true
---

# 1: Remove language dependency from test pipeline

## Description

Refactor the test pipeline task generation to bind tasks to interface types only (api, cli, web-ui, etc.), removing all language coupling. This is the core refactoring that resolves the silent test pipeline skipping issue in monorepo projects.

The change touches 4 source files and their tests atomically (all signature changes are interdependent):
- `testgen.go` — Remove language loops, simplify task keys to interface-only format
- `build.go` — Remove `languages` parameter from `GenerateTestTasks`, add warning for missing interfaces
- `detect.go` — Remove all language detection functions, simplify `ReadInterfaces` to config-only
- `infer.go` — Remove `profileSuffixedID`, simplify type matching

## Reference Files
- `docs/proposals/decouple-test-tasks-from-languages/proposal.md` — Source proposal
- `docs/lessons/gotcha-test-pipeline-no-languages.md` — Root cause analysis

## Acceptance Criteria

- [ ] `GetBreakdownTestTasks` and `GetQuickTestTasks` accept `interfaces []string` only (no `languages` param)
- [ ] Task keys use interface-only format: `gen-test-scripts-api`, `run-e2e-tests`, `graduate-tests` (no language in key)
- [ ] `profileSuffix` and `suffixLetter` functions removed from `testgen.go`
- [ ] `profileSuffixedID` function removed from `infer.go`; `typeSuffixedID` simplified (no optional profile letter)
- [ ] `DetectLanguages`, `ReadLanguages`, `UnionLanguageInterfaces`, `defaultInterfaces` removed from `detect.go`
- [ ] `languageCapabilities`, `KnownLanguages`, `IsKnownLanguage` removed from `detect.go`
- [ ] `detectPytest`, `fileExists`, `dirExists` helper functions removed from `detect.go`
- [ ] `ReadInterfaces` simplified to read `interfaces` from config.yaml only; returns empty slice (no error) when not configured
- [ ] `GenerateTestTasks` signature updated to remove `languages` parameter
- [ ] `BuildIndex` outputs a clear warning when `interfaces` is empty and test pipeline is needed
- [ ] All unit tests updated and passing: `testgen_test.go`, `infer_test.go`, `detect_test.go`, `build_test.go`
- [ ] Integration tests updated: `tests/test-generation/test_scripts_per_type_test.go`
- [ ] `go build ./...` passes
- [ ] `go test -race -cover ./...` passes

## Hard Rules

- ALL changes must compile together — do not leave intermediate states where `build.go` calls deleted functions
- `Languages` field remains in Config struct (backward compatibility with existing config.yaml files)
- Do not change `forge config get` behavior — it already doesn't support `languages`/`interfaces` keys

## Implementation Notes

### testgen.go changes

**Breakdown mode task structure** (was per-language, now single block):

```
gen-test-cases (shared)
eval-test-cases (shared)
gen-test-scripts-{type}  (one per interface type)
run-e2e-tests (single, no language suffix)
graduate-tests (single, no language suffix)
verify-regression (shared)
```

**Quick mode task structure** (was per-language, now single block):

```
quick-test-cases (single)
quick-gen-and-run-{type} (one per interface type)
quick-graduate (single)
quick-verify-regression (shared)
```

**Dependency resolution**: Remove `languages` and `suffix` params from `resolveBreakdownDeps` and `resolveQuickDeps`. Simplify block size calculation — no per-language block iteration needed.

### detect.go changes

Replace `ReadInterfaces` with:

```go
func ReadInterfaces(projectRoot string) ([]string, error) {
    cfg, err := ReadConfig(projectRoot)
    if err != nil || cfg == nil {
        return nil, nil
    }
    return cfg.Interfaces, nil
}
```

Delete everything else in the file except this function.

### build.go changes

```go
// In BuildIndex, replace:
//   languages, _ := forgeconfig.ReadLanguages(opts.ProjectRoot)
//   capabilities, _ := forgeconfig.ReadInterfaces(opts.ProjectRoot)
//   testTasks := GenerateTestTasks(mode, languages, capabilities, opts.AutoConfig)
// With:
//   capabilities, _ := forgeconfig.ReadInterfaces(opts.ProjectRoot)
//   if len(capabilities) == 0 {
//       result.Warnings = append(result.Warnings, "No interfaces configured ...")
//   }
//   testTasks := GenerateTestTasks(mode, capabilities, opts.AutoConfig)
```

### infer.go changes

Remove `profileSuffixedID` entirely. Simplify `InferType` to remove all `profileSuffixedID` calls. Simplify `typeSuffixedID` to remove optional profile letter handling (lines 90-92). Simplify `ExtractTypeSuffix` similarly.

### Test updates

- `testgen_test.go`: Remove multi-language tests, update assertions for new key format
- `infer_test.go`: Remove profile suffix test cases, update type suffix tests
- `detect_test.go`: Delete entirely (all tested functions removed)
- `build_test.go`: Update `GenerateTestTasks` calls to new signature
- `tests/test-generation/test_scripts_per_type_test.go`: Update task key assertions
