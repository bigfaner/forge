---
id: "1.gate"
title: "Phase 1 Gate: CLI 基础能力验证"
priority: "P0"
estimated_time: "30min"
dependencies: ["1.summary"]
status: pending
breaking: true
type: "gate"
---

# 1.gate: Phase 1 Gate — CLI 基础能力验证

## Description

Exit verification gate for Phase 1. Confirms that all CLI commands are implemented, tested, and ready before proceeding to Phase 2 (Schema & Templates).

## Verification Checklist

- [ ] `task prompt <id>` 对 11 种 type 各输出正确 prompt（无 `{{` 残留）
- [ ] `task prompt <id>` 在 type 缺失/未知时 exit 1，stdout 为空
- [ ] `task prompt <id> --fix-record-missed` 使用 fix-record-missed 模板
- [ ] `task migrate` 对所有已知 ID 模式推断正确 type
- [ ] `task migrate` 在 in_progress 任务存在时报错，index.json 不修改
- [ ] `task validate` 对缺失/非法 type 报 error
- [ ] `task claim` 输出包含 TYPE 字段
- [ ] `go test ./...` 通过，`pkg/prompt` 覆盖率 ≥ 80%
- [ ] `golangci-lint run ./...` 无新增 lint 错误

## Acceptance Criteria

- [ ] 所有上述检查项通过
- [ ] `task validate docs/features/typed-task-dispatch/tasks/index.json` 无报错
