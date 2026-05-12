---
id: "4.gate"
title: "Phase 4 Gate: 清理验证"
priority: "P0"
estimated_time: "30min"
dependencies: ["4.summary"]
status: pending
breaking: true
type: "gate"
---

# 4.gate: Phase 4 Gate — 清理验证

## Description

Exit verification gate for Phase 4. Final gate before T-test tasks. Confirms the full feature is complete and consistent.

## Verification Checklist

- [ ] `grep -r "forge:error-fixer" plugins/forge/` 无结果
- [ ] `grep -r "error-fixer" plugins/forge/commands/` 无结果
- [ ] error-fixer.md 顶部含 deprecated 标注
- [ ] ARCHITECTURE.md error-fixer 描述已更新
- [ ] `task validate docs/features/typed-task-dispatch/tasks/index.json` 无报错
- [ ] `go build ./...` 通过（task-cli）
- [ ] `go test ./...` 通过，`pkg/prompt` 覆盖率 ≥ 80%
- [ ] `golangci-lint run ./...` 无新增 lint 错误

## Acceptance Criteria

- [ ] 所有上述检查项通过
- [ ] 整体功能可端到端验证：`task prompt <id>` → `Agent(forge:task-executor)` 链路正常
