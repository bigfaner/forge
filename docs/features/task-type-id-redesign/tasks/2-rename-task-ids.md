---
id: "2"
title: "Rename task IDs and update generation/inference"
priority: "P1"
estimated_time: "55min"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 2: Rename task IDs and update generation/inference

## Description
In `testgen.go`, rename all auto-generated task IDs to readable format, add validation task generation (`T-validate-code`, `T-validate-ux`), update dependency chain references. In `infer.go`, update all ID pattern mappings for `InferType()`, add validation entries, update `profileSuffixedID`/`typeSuffixedID` base strings.

### ID mapping (Part B)
| Old ID | New ID |
|--------|--------|
| `T-test-1` | `T-test-gen-cases` |
| `T-test-1b` | `T-test-eval-cases` |
| `T-test-2` | `T-test-gen-scripts` |
| `T-test-3` | `T-test-run` |
| `T-test-4` | `T-test-graduate` |
| `T-test-4.5` | `T-test-verify-regression` |
| `T-specs-1` | `T-specs-consolidate` |
| `T-quick-1` | `T-quick-gen-cases` |
| `T-quick-2` | `T-quick-gen-and-run` |
| `T-quick-3` | `T-quick-graduate` |
| `T-quick-4` | `T-quick-verify-regression` |
| `T-quick-specs-1` | `T-quick-doc-drift` |
| `T-clean-code-1` | `T-clean-code` |
| _(new)_ | `T-validate-code` |
| _(new)_ | `T-validate-ux` |

### Validation task generation (Part D subset)
- `GetBreakdownTestTasks()` and `GetQuickTestTasks()` add `T-validate-code` and `T-validate-ux`
- `T-validate-code`: after last business task, before test pipeline; `noTest: true`, `mainSession: false`
- `T-validate-ux`: after test pipeline, before all-completed; `noTest: true`, `mainSession: true`
- New `auto.Validation.Full` / `auto.Validation.Quick` config fields, default `false`

## Reference Files
- `docs/proposals/task-type-id-redesign/proposal.md` — Source proposal
- `forge-cli/pkg/task/testgen.go` — Task ID generation, dependency chains
- `forge-cli/pkg/task/infer.go` — InferType, profileSuffixedID, typeSuffixedID
- `forge-cli/pkg/task/build.go` — validation auto config

## Acceptance Criteria
- [ ] All task IDs renamed per mapping table
- [ ] `InferType("T-test-gen-cases")` → `"test.gen-cases"`
- [ ] `InferType("T-validate-code")` → `"validation.code"`
- [ ] `profileSuffixedID`/`typeSuffixedID` base strings updated to new ID prefixes
- [ ] Dependency chain references use new IDs
- [ ] `T-validate-code` generated when `auto.Validation` enabled
- [ ] `T-validate-ux` generated when `auto.Validation` enabled

## Hard Rules
- Full-text search all Go files for old ID strings and replace — no orphaned references
- `InferType` still uses a map (ID suffix ≠ type name), but keys change from numeric to readable strings

## Implementation Notes
- Profile suffix pattern preserved: `T-test-gen-scripts-a` (profile a), `T-test-gen-scripts-api` (type api)
- `genScriptBases` must update to `"T-test-gen-scripts"`, `"T-quick-gen-and-run"`
- Validation tasks are gated behind `auto.Validation` config — only generated when enabled
