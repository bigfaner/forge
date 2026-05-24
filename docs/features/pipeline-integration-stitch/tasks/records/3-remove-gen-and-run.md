---
status: "completed"
started: "2026-05-24 18:51"
completed: "2026-05-24 19:14"
time_spent: "~23m"
---

# Task Record: 3 移除 test.gen-and-run 废弃代码 + 更新测试文件

## Summary
Removed test.gen-and-run / T-quick-gen-and-run deprecated code from production code, tests, and template files. Added migration error guidance in validate_index.go and prompt.go Synthesize().

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/task/category_test.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/record_test.go
- forge-cli/pkg/task/stage_gates_test.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/internal/cmd/task/validate_index.go
- forge-cli/internal/cmd/task/validate_index_test.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/internal/cmd/task/list_test.go
- forge-cli/internal/cmd/base/output_test.go
- forge-cli/tests/task-lifecycle/task_stage_gates_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Kept migration guard strings in validate_index.go and prompt.go as AC3/AC4 require diagnostic messages for old index.json files
- Replaced removed test assertions with equivalent checks using strings.Contains or IsSystemType to maintain test coverage
- Replaced test.gen-and-run in TypeReclassification test with test.run to keep the test valid

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] grep gen-and-run in forge-cli/ (excluding docs/proposals) returns only migration guard references
- [x] grep gen-and-run in plugins/forge/ returns zero results
- [x] validate_index.go returns migration guidance for test.gen-and-run references
- [x] Synthesize() ReadFile failure outputs migration guidance for gen-and-run filenames
- [x] All existing tests pass

## Notes
Deprecated template files prompt/data/test-gen-and-run.md and task/data/test-gen-and-run.md deleted. SystemTypes count reduced from 14 to 13. Version bumped to 5.4.5 (patch).
