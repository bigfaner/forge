---
id: "4"
title: "Update test runner probe chain and journey isolation"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 4: Update test runner probe chain and journey isolation

## Description

更新 `testrunner.go` 的探测链优先级为 `unit-test → test → go test`，替换原有的 `test → go test`。将 `journey_isolation.go` 中的 `just e2e-test` 调用迁移为 `just test`（支持 journey 参数）。

## Reference Files
- `proposal.md#Proposed-Solution` — RunProjectTests probe chain (unit-test → test → go test) with fallback; gate sequence without fallback (independent paths)
- `proposal.md#Requirements-Analysis` — Recipe 参数签名约定: `test journey=''` optional positional argument
- `proposal.md#Constraints-&-Dependencies` — journey_isolation.go must migrate from `just e2e-test` to `just test` with journey parameter
- `proposal.md#Key-Risks` — risk of journey_isolation migration affecting existing journey tests; mitigation via test recipe supporting positional argument

## Acceptance Criteria
- `RunProjectTests()` probe chain: `unit-test` → `test` → `go test` (按优先级依次 `HasRecipe()`)
- `journey_isolation.go` calls `just test <journeyName>` instead of `just e2e-test`
- `RunProjectTests` retains fallback mechanism (first matching recipe wins)
- Gate sequence path and RunProjectTests path behave independently

## Hard Rules
- `RunProjectTests` 探测链保留 fallback 机制——与 gate sequence 的无 fallback 策略独立
- `test` recipe 调用需传入 journey 参数（positional argument）

## Implementation Notes
- `RunProjectTests` 用于非 gate 场景（如 `forge run-tests` CLI 命令），保留 fallback 是合理的
- `journey_isolation.go` 的 journey 参数传递方式需匹配 justfile 模板中的 `test journey=''` 签名
