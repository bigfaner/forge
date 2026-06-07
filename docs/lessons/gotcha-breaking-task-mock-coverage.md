---
created: "2026-06-07"
tags: [testing, interface, go]
---

# breaking: true tasks must verify ALL downstream mock implementations compile

## Problem
Task 2.3 (MainItem milestone repo methods) completed successfully with `breaking: true` — it added 4 methods to the shared `MainItemRepo` interface. The task executor agent claimed it "Updated 7 mock/stub implementations" but missed `mockViewMainItemRepo` in `report_service_test.go`. This caused compilation errors discovered only by the dispatcher's next claim cycle, requiring a fix-1 task.

## Root Cause
1. Agent updated the most obvious mock implementations (those in the same package or directly referenced) but missed a mock in an unrelated test file (`report_service_test.go`)
2. The agent's quality gate used `go build ./...` — but **`go build` 不编译 `_test.go` 文件**（这是 Go 的设计，不是意外），所以 mock 遗漏必然无法被发现
3. The `breaking: true` flag on the task warned about this risk, but the agent didn't exhaustively verify ALL implementations of the modified interface
4. Dispatcher had no mechanism to pre-validate compilation (including test files) before marking a breaking task as complete

## Solution
When a task has `breaking: true` and modifies a Go interface:
1. After implementation, run `go test -run=^$ ./...`（编译所有测试代码但不执行测试，是检查 interface satisfaction 最可靠的方式）
2. Optionally use `grep -rl "InterfaceName" --include="*_test.go"` as a supplementary check to locate mock files for manual review
3. Do NOT rely on `go build ./...` — it skips `_test.go` files by design and will miss mock compilation errors

## Reusable Pattern
For any `breaking: true` task that modifies a shared Go interface:
```bash
# 1. Compile all test code (catches missing interface methods in mocks)
cd backend && go test -run=^$ ./...

# 2. (Optional) Locate mock files for manual review
grep -rl "repository.MainItemRepo" backend/ --include="*_test.go"
```

## Related Files
- `backend/internal/repository/main_item_repo.go` — the modified interface
- `backend/internal/service/report_service_test.go` — the missed mock
