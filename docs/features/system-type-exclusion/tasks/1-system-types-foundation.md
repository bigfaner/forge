---
id: "1"
title: "Add SystemTypes set, IsSystemType(), and remove TypeCodingClean dead code"
priority: "P1"
estimated_time: "45min"
dependencies: []
scope: "backend"
breaking: true
type: "coding.feature"
mainSession: false
---

# 1: Add SystemTypes set, IsSystemType(), and remove TypeCodingClean dead code

## Description

在 `types.go` 中新增 `SystemTypes` 黑名单集合（13 种自动生成的系统类型）和 `IsSystemType()` 查询函数。同时移除 `TypeCodingClean`（`coding.clean`）死代码——该常量无任何生产引用。

系统类型列表（13 种）：
- `gate`
- `test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.graduate`, `test.verify-regression`（7 种测试类型）
- `validation.code`, `validation.ux`（2 种验证类型）
- `doc.eval`, `doc.summary`
- `code-quality.simplify`

排除（双重身份，可做业务任务）：`doc.consolidate`、`doc.drift`

## Reference Files
- `docs/proposals/system-type-exclusion/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — 类型常量、ValidTypes、TaskTypeRegistry
- `forge-cli/pkg/task/types_test.go` — 类型相关测试
- `forge-cli/pkg/task/testgen.go` — 交叉验证自动生成类型
- `forge-cli/pkg/task/stage_gates.go` — 交叉验证自动生成类型

## Acceptance Criteria

- [ ] `SystemTypes` map 包含恰好 13 个条目
- [ ] `IsSystemType()` 对所有 13 种系统类型返回 true
- [ ] `IsSystemType()` 对业务类型和双重身份类型（doc.consolidate、doc.drift）返回 false
- [ ] `TypeCodingClean` 常量从 types.go 中移除
- [ ] `TypeCodingClean` 从 TaskTypeRegistry 中移除
- [ ] `TypeCodingClean` 从 ValidTypes 中移除
- [ ] types_test.go 中 TypeCodingClean 相关测试移除或更新
- [ ] `go test ./forge-cli/...` 通过

## Hard Rules

- ValidTypes 总数从 22 减少到 21（移除 coding.clean，不新增——SystemTypes 是独立的 map）
- SystemTypes 是独立的 `map[string]bool`，与 ValidTypes 无关

## Implementation Notes

- 交叉验证：从 testgen.go 和 stage_gates.go 提取所有自动生成任务的类型，确保 13 种完整无遗漏
- `coding.clean` 已确认无生产引用，仅在 types.go 声明、注册表、测试文件中出现
- Key Risk：系统类型集合遗漏某种自动生成类型 → 用 testgen.go、stage_gates.go、infer.go 三文件交叉验证
