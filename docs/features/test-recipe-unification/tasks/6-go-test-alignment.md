---
id: "6"
title: "Align Go tests with new recipe names and config fields"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["1", "2", "3", "4", "5"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: Align Go tests with new recipe names and config fields

## Description

更新所有 Go 测试文件中的断言和 fixture，使其与源码中的新 recipe 名和 config 字段对齐。修复已失败的集成测试 TC_005/TC_015/TC_016。

## Reference Files
- `proposal.md#Impact-Analysis` — Tier 4 lists all 11+ Go test files with specific expected changes per file
- `proposal.md#Success-Criteria` — all Go tests pass (`go test -race ./...`), TC_005/TC_015/TC_016 fixed, no residual e2e-test/e2eTest references

## Acceptance Criteria
- `pkg/just/just_test.go`: assertions for `test` → `unit-test`
- `pkg/testrunner/testrunner_test.go`: justfile fixture `test:` → `unit-test:`
- `internal/cmd/quality_gate_test.go`: `HasRecipe(dir, "test")` → `HasRecipe(dir, "unit-test")`
- `forgeconfig/config_test.go`: all `E2eTest` assertions → `Test`
- `task/autoconfig_test.go`: `auto.E2eTest` → `auto.Test`
- `task/autogen_test.go`: `E2eTest` fixture + `run-e2e-tests` key → `run-test`
- `task/submit_test.go`: `run-e2e-tests` fixtures → `run-test`
- `task/status_test.go`: `run-e2e-tests` fixtures → `run-test`
- `tests/justfile-integration/mixed_cli_test.go`: `just test` → `just unit-test`; fix TC_005/TC_015/TC_016
- `tests/justfile-integration/forge_detection_test.go`: recipe list `test` → `unit-test`
- `tests/task-type-system/task_types_dispatch_test.go`: `just test` → `just unit-test`
- All tests pass with `go test -race ./...`

## Hard Rules
- TC_005/TC_015/TC_016 修复需对齐实际文件内容，而非简单跳过

## Implementation Notes
- 大部分变更是机械性的 find-and-replace（字段名、recipe 名、key 名）
- TC_005/TC_015/TC_016 在 proposal 前已在失败，需分析根本原因再修复
- 此任务应在所有源码任务（1-5）完成后执行，确保引用的符号都已更新
