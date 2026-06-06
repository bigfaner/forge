---
created: "2026-06-06"
author: "faner"
status: Draft
intent: "cleanup"
---

# Proposal: Forge CLI 可读性全面清理（第二轮）

## Problem

Forge CLI 生产代码中存在多个超大函数和文件，严重降低人类可读性。最大的函数 `BuildIndex` 达 390 行，8 个文件超过 500 行目标上限，多处嵌套达 7 层。在最近的功能开发中，开发者需要在 390 行的函数中反复滚动才能理解上下文，显著拖慢了开发和 review 效率。

### Evidence

| 文件 | 行数 | 最大函数 | 最大函数行数 | 最大嵌套 |
|------|------|----------|-------------|----------|
| `pkg/task/build.go` | 682 | `BuildIndex` | **390** | 5 |
| `internal/cmd/forensic/extract.go` | 321+ | `runExtract` | **304** | 7+ |
| `pkg/forgeconfig/config.go` | 1365 | `setByPath` | 62 | 7 |
| `pkg/task/pipeline.go` | 1103 | `GenerateTestTasks` | 79 | 7 |
| `pkg/forgeconfig/detect_surface.go` | 963 | `DetectSurfacesWithConflicts` | 78 | 7 |
| `internal/cmd/task/list.go` | 454 | `runList` | **217** | 5 |
| `internal/cmd/task/submit.go` | 407 | `doSubmit` | 131 | 4 |
| `internal/cmd/task/validate.go` | 573 | `validateGateIntegrity` | 66 | **7** |

- `config.go` 混合了 3 种职责：配置读写、reflect 路径遍历、AutoConfig 默认值
- `pipeline.go` 445 行非函数代码（var 块、类型定义）打断阅读流
- `detect_surface.go` 前 150 行全是信号映射表，推断函数按生态重复模式
- `quality_gate.go` 含 4 处 `os.Exit(0)` 导致函数不可测试

### Urgency

Forge CLI 正在活跃开发中，可读性债务会随功能增加持续累积。`BuildIndex` 的 390 行上帝函数每次修改都需要完整理解 9 个步骤的上下文，边际成本递增。现在清理的投资回报率最高。

## Proposed Solution

对所有超标文件执行系统性分解重构，遵循 4 条硬约束：函数 ≤ 80 行、文件 ≤ 500 行、嵌套 ≤ 4 层、每文件单一职责。同时清除已确认的死代码并修复 `os.Exit` 反模式。

### Innovation Highlights

这是一次标准的 Go 代码健康度重构，无特殊创新。采用的手段包括：提取命名函数替代步骤注释、文件拆分按职责边界、early return / guard clause 平坦化嵌套、`os.Exit` 改为 error return。遵循 Go 社区通用的代码组织实践。

## Requirements Analysis

### Key Scenarios

- 开发者打开任意文件，无需上下滚动即可看到完整函数体
- 开发者阅读函数时，嵌套不超过 4 层，控制流线性清晰
- 开发者定位某职责时，文件名即可指向正确位置
- CI 运行 `go test ./...` 全绿，行为零变更

### Non-Functional Requirements

- **零行为变更**：所有 CLI 命令的输入输出保持不变
- **向后兼容**：`pkg/` 层导出 API 签名不变

### Constraints & Dependencies

- 所有现有测试必须通过，不新增测试文件
- 遵循 `cmd -> internal -> pkg` 依赖方向（CLAUDE.md 约束）
- 遵循 Go 标准 `cmd/internal/pkg` 项目布局

## Alternatives & Industry Benchmarking

### Industry Solutions

Go 社区普遍采用文件拆分（同包多文件）和函数提取来控制代码复杂度。`golangci-lint` 的 `gocyclo`、`funlen`、`nestif` linter 分别检测圈复杂度、函数长度和嵌套深度。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 可读性债务持续累积，开发效率下降 | Rejected: 边际成本已显著 |
| 仅拆分最大 3 个文件 | 局部优化 | 改动小，风险低 | 其他 5 个文件仍超标 | Rejected: 用户要求全面清理 |
| 引入 gocyclo/funlen lint 门禁 | golangci-lint | 自动化预防 | 不解决存量问题 | 可作为后续跟进 |
| **全面分解重构** | Go 社区标准实践 | 彻底解决，一次性收益 | 改动范围大，需仔细验证 | **Selected: 用户确认全面清理** |

## Feasibility Assessment

### Technical Feasibility

纯重构，不涉及架构变更或外部依赖。Go 的同包多文件机制天然支持文件拆分（无需改包名或导入路径）。

### Resource & Timeline

10 个改动点，每个平均涉及 1-2 个文件的拆分或重组。估计 1-2 天完成。

### Dependency Readiness

无外部依赖。所有工具（`go test`、`golangci-lint`）已就绪。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 未使用的导出符号应清除 | XY Detection | Overridden: 用户明确保留所有导出符号（公共 API 设计意图） |
| 测试文件需要同步更新 | 5 Whys | Refined: 分解重构后现有测试应自然通过；若 `os.Exit` 改造导致测试失败，需调整实现策略确保测试通过 |
| `Scope` 字段是死代码应清除 | XY Detection | Refined: 保留用于迁移兼容，不在本次范围内 |

## Scope

### In Scope

1. `pkg/task/build.go` — 将 390 行 `BuildIndex` 拆分为 ~9 个命名步骤函数
2. `internal/cmd/forensic/extract.go` — 将 304 行 `runExtract` 拆分为解析/聚合/输出阶段
3. `pkg/forgeconfig/config.go` — 提取 reflect 路径遍历机点到 `config_reflect.go`
4. `pkg/task/pipeline.go` — 提取校验逻辑到 `pipeline_validate.go`，重组 var 块
5. `pkg/forgeconfig/detect_surface.go` — 提取信号表到 `detect_surface_signals.go`，统一推断模式
6. `internal/cmd/qualitygate/quality_gate.go` — `os.Exit(0)` 改为返回 error，由调用方处理
7. `internal/cmd/task/list.go` — 拆分 217 行 `runList`
8. `internal/cmd/task/submit.go` — 拆分 131 行 `doSubmit`
9. `internal/cmd/task/validate.go` — 平坦化嵌套过深的 validator 方法
10. 删除死代码：`requireSurfaceInference`（quality_gate.go）、`extractScope`（extract.go），同步删除对应测试用例

### Out of Scope

- 未使用的导出符号（合理公共 API 设计保留）
- `Scope` 字段（迁移兼容保留）
- 任何 CLI 行为变更
- 新增测试用例
- `gocyclo`/`funlen`/`nestif` lint 门禁配置（可后续跟进）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 文件拆分引入 import cycle | L | M | 严格遵循 `cmd -> internal -> pkg` 方向；同包拆分不涉及跨包引用 |
| `os.Exit` 改造破坏现有测试 | M | H | 先分析测试结构，采用 error return + 顶层 exit 策略；若测试直接调用函数则返回值兼容 |
| 大范围改动引入回归 | M | H | 每个文件独立重构+测试验证；`go test ./...` 全绿后才提交 |
| 拆分过度导致文件碎片化 | L | L | 按职责边界拆分，不按函数数量机械拆分 |

## Success Criteria

consistency_check_result:
  status: pass
  pairs_checked: 28
  conflicts_found: 1
  resolved: SC-4 vs InScope-10 — 放宽 SC-4 允许同步删除被清理函数的测试用例

- [ ] SC-1: 所有生产 .go 函数 ≤ 80 行（`golangci-lint funlen` 或人工验证）
- [ ] SC-2: 所有生产 .go 文件 ≤ 500 行
- [ ] SC-3: 所有函数嵌套 ≤ 4 层（`golangci-lint nestif` 或人工验证）
- [ ] SC-4: `go test ./...` 全绿；仅允许删除被清理函数对应的测试用例，不新增测试
- [ ] SC-5: 零行为变更（CLI 输出与重构前一致）
- [ ] SC-6: 死函数已删除：`requireSurfaceInference`、`extractScope`
- [ ] SC-7: `os.Exit` 仅存在于 `cmd/forge/` 入口和 `base.Exit` 统一出口，`quality_gate.go` 无直接 `os.Exit` 调用

## Next Steps

- Proceed to `/breakdown-tasks` to create task breakdown (skip PRD/tech-design for cleanup)
