---
created: 2026-05-20
author: "fanhuifeng"
status: Draft
---

# Proposal: 系统类型排除——防止业务任务使用 CLI 自动生成的任务类型

## Problem

业务任务（由 Skill 创建的 .md 文件）可以标记 `gate`、`test.*`、`validation.*` 等 CLI 自动生成的系统类型，导致 run-tasks 调度器行为异常。

### Evidence

实际发生过业务任务被误标为 gate/test 类型，引发任务调度异常。

### Urgency

调度异常直接影响任务执行管线的正确性，可能导致质量门跳过、依赖链断裂。每次误标都是一次管线故障。

## Proposed Solution

在 CLI 校验层和 Skill 规则层实施双层防护，采用**反向排除**策略：维护一个系统类型黑名单（`SystemTypes`），非自动生成任务的类型命中黑名单则拒绝。新业务类型自动合法，无需维护白名单。

### Innovation Highlights

反向排除策略避免了正向白名单的维护负担。系统类型是 CLI 基础设施，变更频率低、集合稳定；而业务类型可能随需求增长，用黑名单天然适配这种不对称性。

## Requirements Analysis

### Key Scenarios

- **正常流程**：Skill 创建业务任务，类型为 `coding.feature` → 通过校验
- **误标拦截**：Skill 创建任务，类型为 `gate` → BuildIndex / validate-index 报错
- **自动生成豁免**：`forge task index` 生成 `T-test-*`、`*.gate` 等任务 → 不受限制
- **质量门 fix 任务**：`addFixTask()` 创建 `coding.fix` / `coding.cleanup` 任务 → 通过校验
- **新增业务类型**：未来在 `types.go` 新增类型，不在黑名单中 → 自动合法

### Non-Functional Requirements

- 错误信息须明确指出哪些类型是系统保留类型
- 校验在 `BuildIndex` 和 `validate-index` 两处一致

### Constraints & Dependencies

- 须兼容 `isAutoGenTaskID()` 现有逻辑（自动生成任务按 ID 模式豁免）
- 须兼容现有 `IsTestableType()` 逻辑（不改）
- 须兼容 legacy 类型值（`implementation`、`fix` 等）——这些不在 SystemTypes 中也不在 ValidTypes 中，现有 validate-index 已会拒绝

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无改动 | 调度异常持续发生 | Rejected: 已有实际问题 |
| 正向白名单（BusinessTypes） | 常见做法 | 显式控制 | 新增业务类型需更新白名单，维护负担高 | Rejected: 可扩展性差 |
| **反向排除（SystemTypes 黑名单）** | — | 系统类型稳定闭集，新业务类型自动合法 | 系统类型增加时需更新（频率极低） | **Selected: 维护负担最低** |

## Feasibility Assessment

### Technical Feasibility

完全在现有 `types.go` + `build.go` + `validate_index.go` 体系内实现，无外部依赖。

### Resource & Timeline

预计 1-2 小时，涉及 Go 代码 + SKILL.md 文档更新。

### Dependency Readiness

无外部依赖，所有相关代码在项目中。

## Scope

### In Scope

- 在 `types.go` 新增 `SystemTypes` 集合（15 种系统类型）+ `IsSystemType()` 函数
- `BuildIndex()` 拦截：非自动生成任务 + 类型命中 SystemTypes → 报错
- `validate-index` 拦截：同上校验逻辑
- 更新 `quick-tasks/SKILL.md` 类型分配表（移除 `gate` 等系统类型）
- 更新 `breakdown-tasks/SKILL.md` 类型分配表（同上）
- 清理 `coding.clean`（TypeCodingClean）死代码：从 types.go、ValidTypes、TaskTypeRegistry、测试文件中移除
- 补充/更新相关单元测试

### Out of Scope

- 重命名/合并 `code-quality.simplify` 和 `coding.cleanup`（功能不同，不能融合）
- 修改自动生成任务的行为
- 修改 `IsTestableType()` 逻辑
- 修复现有 index.json 中的 legacy 类型值（独立问题，由 task-index-dedup-legacy-types 处理）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| SKILL.md 中 `gate` 等类型引用被遗漏 | L | M | 全面搜索所有 SKILL.md 中的类型分配表 |
| 清理 `coding.clean` 影响其他代码 | L | L | 已确认无任何生产代码引用，仅 types.go 声明和测试 |
| 系统类型集合遗漏某种自动生成类型 | L | M | 从 testgen.go、quality_gate.go、infer.go 交叉验证 |

## Success Criteria

- [ ] `forge task validate-index` 对非自动生成任务使用系统类型时报错，错误信息包含具体类型和系统类型列表
- [ ] `forge task index`（BuildIndex）同样拦截系统类型误用
- [ ] 自动生成任务（`T-test-*`、`*.gate` 等）不受影响
- [ ] 质量门 fix 任务（`coding.fix`、`coding.cleanup`）不受影响
- [ ] `coding.clean` 常量、注册项、测试全部移除
- [ ] quick-tasks 和 breakdown-tasks 的 SKILL.md 类型分配表不含系统类型

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
