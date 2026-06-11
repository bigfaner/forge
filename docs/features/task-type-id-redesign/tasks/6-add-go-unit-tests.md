---
id: "6"
title: "Add Go unit tests for type/ID/prefix changes"
priority: "P1"
estimated_time: "40min"
dependencies: ["1", "2", "3"]
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 6: Add Go unit tests for type/ID/prefix changes

## Description
Add Go unit tests covering the new type constants, prefix-based checks, ID inference, and validation task generation introduced in tasks 1-3.

### Test coverage areas
1. **Type constants** — verify all new constants have correct string values
2. **IsTestableType** — prefix-based: `coding.*` → true, `doc*` → false, `test.*` → false, `validation.*` → false
3. **isDocsOnlyType** — prefix-based: `doc`, `doc.eval`, `doc.consolidate` → true; `coding.feature` → false
4. **InferType** — all new IDs map to correct types (`T-test-gen-cases` → `test.gen-cases`, `T-validate-code` → `validation.code`, etc.)
5. **isAutoGenTaskID/isTestTaskID** — recognize new ID prefixes
6. **Validation task generation** — `T-validate-code` and `T-validate-ux` generated when config enabled

## Reference Files
- `docs/proposals/task-type-id-redesign/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — Type constants
- `forge-cli/pkg/task/build.go` — IsTestableType, isDocsOnlyType, isAutoGenTaskID
- `forge-cli/pkg/task/infer.go` — InferType
- `forge-cli/pkg/task/testgen.go` — Task generation
- `forge-cli/pkg/task/*_test.go` — Existing test files

## Acceptance Criteria
- [ ] Test for each new type constant value
- [ ] `IsTestableType` tested with `coding.feature`, `coding.enhancement`, `doc`, `test.gen-cases`, `validation.code`
- [ ] `isDocsOnlyType` tested with `doc`, `doc.eval`, `coding.feature`
- [ ] `InferType` tested with all new IDs from Part B mapping
- [ ] `isAutoGenTaskID` tested with new prefixes
- [ ] Validation task generation tested (enabled/disabled)

## Hard Rules
- Tests must pass without any existing features' index.json files present
- Use table-driven tests for IsTestableType, isDocsOnlyType, InferType

## Implementation Notes
- Focus on the success criteria from the proposal — these are the concrete acceptance tests
- Existing test files in `pkg/task/` may need updating for renamed constants/IDs
