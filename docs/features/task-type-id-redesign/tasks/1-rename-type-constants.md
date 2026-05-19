---
id: "1"
title: "Rename type constants and implement prefix checks"
priority: "P1"
estimated_time: "50min"
dependencies: []
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 1: Rename type constants and implement prefix checks

## Description
Rename all `Type*` constants in `types.go` to prefix format (`coding.*`, `doc*`, `test.*`, `validation.*`), add `TypeValidationCode` and `TypeValidationUx`, update Registry and ValidTypes. In `build.go`, replace `IsTestableType` and `isDocsOnlyType` with prefix-based checks, update `isAutoGenTaskID`/`isTestTaskID` to cover new ID prefixes, add validation auto config field.

### Mapping (Part A)
| New constant | Old constant |
|---|---|
| `TypeCodingFeature` | `TypeFeature` |
| `TypeCodingEnhancement` | `TypeEnhancement` |
| `TypeCodingCleanup` | `TypeCleanup` |
| `TypeCodingRefactor` | `TypeRefactor` |
| `TypeCodingFix` | `TypeFix` |
| `TypeCodingClean` | _(new)_ |
| `TypeDoc` | `TypeDocumentation` |
| `TypeDocEval` | _(new)_ |
| `TypeDocSummary` | _(new)_ |
| `TypeDocConsolidate` | _(new)_ |
| `TypeDocDrift` | _(new)_ |
| `TypeTestGenCases` | `TypeTestPipelineGenCases` |
| `TypeTestEvalCases` | `TypeTestPipelineEvalCases` |
| `TypeTestGenScripts` | `TypeTestPipelineGenScripts` |
| `TypeTestRun` | `TypeTestPipelineRun` |
| `TypeTestGenAndRun` | `TypeTestPipelineGenAndRun` |
| `TypeTestGraduate` | `TypeTestPipelineGraduate` |
| `TypeTestVerifyRegression` | `TypeTestPipelineVerifyRegression` |
| `TypeValidationCode` | _(new)_ |
| `TypeValidationUx` | _(new)_ |

### Prefix checks
```go
func IsTestableType(typ string) bool {
    return strings.HasPrefix(typ, "coding.")
}

func isDocsOnlyType(typ string) bool {
    return strings.HasPrefix(typ, "doc")
}
```

### ID prefix coverage
Update `isAutoGenTaskID` and `isTestTaskID` to recognize new prefixes: `T-test-`, `T-quick-`, `T-specs-`, `T-clean-`, `T-validate-`, `T-eval-`.

## Reference Files
- `docs/proposals/task-type-id-redesign/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — Type constants, Registry, ValidTypes
- `forge-cli/pkg/task/build.go` — IsTestableType, isDocsOnlyType, isAutoGenTaskID, isTestTaskID

## Acceptance Criteria
- [ ] All type constants renamed to prefix format per mapping table
- [ ] `TypeValidationCode` and `TypeValidationUx` constants added
- [ ] Registry and ValidTypes updated with all new constants
- [ ] `IsTestableType("coding.feature")` → true
- [ ] `IsTestableType("doc")` → false
- [ ] `IsTestableType("test.gen-cases")` → false
- [ ] `isDocsOnlyType("doc")` → true
- [ ] `isDocsOnlyType("doc.eval")` → true
- [ ] `isDocsOnlyType("coding.feature")` → false
- [ ] `isAutoGenTaskID` and `isTestTaskID` recognize all new ID prefixes

## Hard Rules
- All callers of old constant names must be updated in the same commit — no partial migration
- `IsTestableType` must use prefix matching, not a hardcoded set

## Implementation Notes
- The `testableTypes` map in types.go becomes unnecessary — replace with `IsTestableType` prefix check
- Search all Go files for references to old constant names (`TypeFeature`, `TypeEnhancement`, etc.) and update
- The `Registry` map values should use the new constant names as keys
