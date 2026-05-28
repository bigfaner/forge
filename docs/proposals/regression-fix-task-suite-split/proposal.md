---
created: "2026-05-28"
author: fanhuifeng
status: Draft
---

# Proposal: Regression Fix Task 按 Test File 拆分（语言无关）

## Problem

Quality-gate hook 检测到 regression test 失败后创建单个 fix task 覆盖所有失败，当失败数量多（20+）且跨多个不相关文件时，agent 执行时卡住。

### Evidence

实际发生：quality-gate 创建的 fix-1 包含 4 个 test suite 共 20+ 失败，agent 执行后长时间无响应被用户手动中断。Concise error 只展示输出尾部，agent 需要读 raw-output.txt 才能看到完整失败列表。详见 `docs/lessons/gotcha-fix-task-broad-scope.md`。

### Urgency

高。每次 regression 测试出现多文件失败都会触发此问题，agent 卡死浪费时间和 token。

## Proposed Solution

在 `quality_gate.go` 中新建 `addRegressionFixTasks` 函数，复用现有 `extractSourceFiles` 提取文件路径，通过命名约定识别测试文件，按测试文件分组创建独立 fix task。每个 task 只包含该测试文件相关的输出行。

当无法识别测试文件时，fallback 到现有按目录分组的 `addFixTask` 行为。

移除 `maxFixTasksPerStep` cap 限制——拆分后每个 fix task scope 已收窄到单文件，cap 不再必要。

### Innovation Highlights

无创新。原提案使用 Go 专属的 `FAIL <package>` 解析，改进为复用现有的语言无关 `extractSourceFiles` + 测试文件命名约定识别。与 CI 系统中按 failing module 分组报告的常规做法一致。

## Requirements Analysis

### Key Scenarios

- **Happy path**：4 个测试文件各有 2-5 个失败 → 创建 4 个 fix task，每个包含对应文件的输出行
- **单文件失败**：1 个测试文件有 3 个失败 → 创建 1 个 fix task，行为与当前一致
- **非标准命名**：测试文件名不匹配任何命名约定 → fallback 到按目录分组
- **全部通过**：无失败 → 不创建任何 fix task（现有行为不变）

### Non-Functional Requirements

- **性能**：解析 test output 的时间可忽略（现有 `extractSourceFiles` 已优化）
- **兼容性**：不影响 compile/fmt/lint/unit-test 步骤的现有逻辑

### Constraints & Dependencies

- 依赖测试文件命名约定（`*_test.go`、`test_*.py`、`*.test.ts` 等），无法识别的测试文件 fallback 到按目录分组
- 依赖现有 `extractSourceFiles` 函数的文件路径提取能力（已支持 15+ 种扩展名）

## Alternatives & Industry Benchmarking

### Industry Solutions

CI 系统普遍按 test module/suite 分组报告失败（GitHub Actions test grouping、JUnit XML testsuite 元素）。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 卡死问题持续 | Rejected: 用户体验差 |
| Go 专属 suite 解析（原提案 v1） | 提案初版 | 实现简单 | 仅适配 Go test 输出 | Rejected: 不适配其他语言 |
| 按目录分组（现有行为） | 当前实现 | 语言无关 | 同一目录多文件仍 scope 过大 | Rejected: 未改进 |
| **按测试文件分组** | 通用分组策略 | 语言无关、scope 最窄、复用现有代码 | 依赖命名约定识别 | **Selected: 最小改动，最大通用性** |

## Feasibility Assessment

### Technical Feasibility

完全可行。复用现有 `extractSourceFiles`（已支持 15+ 扩展名），新增 `isTestFile` 命名约定匹配函数。`extractSourceFiles` 和 `groupFilesByDir` 已有完善的单元测试覆盖。

### Resource & Timeline

预计 2-3 小时实现 + 测试。改动集中在 `quality_gate.go`，新增 `isTestFile` 和 `addRegressionFixTasks`。

### Dependency Readiness

无外部依赖。`extractSourceFiles` 已稳定运行。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 需要按 Go suite 解析 FAIL 行 | XY Detection: 真正的需求是缩小 fix task scope，不是解析 suite | Overturned: 复用现有 extractSourceFiles + 测试文件命名约定即可，无需框架专属解析 |
| maxFixTasksPerStep cap 是必要的 | 5 Whys: cap 是因为单个 fix task scope 过大导致循环 | Overridden: 拆分后每个 task scope 窄，cap 不再必要 |
| 需要接口抽象层支持多框架 | Occam's Razor: 只有一个真实场景 | Overturned: 函数封装 + 命名约定已足够，无需接口 |

## Scope

### In Scope

- 新建 `isTestFile` 函数，识别 Go/Python/JS/TS/Java/Ruby 等常见测试文件命名约定
- 新建 `addRegressionFixTasks` 函数，按测试文件分组创建 fix task
- 每个 fix task 只包含该测试文件相关的输出行（从 output 中提取包含该文件路径的行及上下文）
- 移除 `maxFixTasksPerStep` 变量和 `countFixTasks` 函数
- `runTestRegression` 失败时调用新函数替代 `addFixTask`
- 无法识别测试文件时 fallback 到现有按目录分组的 `addFixTask`
- 更新或新增单元测试覆盖新逻辑

### Out of Scope

- 基线对比过滤预存在失败（lesson 第二层改进）
- unit-test / compile / lint 步骤的改动
- Surface inference 改进
- `coding.fix.md` 模板变更

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 测试文件命名不规范导致识别失败 | M | L | fallback 到按目录分组，零功能损失 |
| 移除 cap 后大量 fix task 并发 | L | M | 每个 fix task scope 已收窄到单文件，并发执行风险可控 |
| 输出行关联测试文件的准确性 | M | L | 多匹配一些行（宁可多包含）比漏掉好，agent 可自行过滤 |
| Rust 等无特殊测试文件命名的语言 | L | L | fallback 到按目录分组 |

## Success Criteria

- [ ] 4 个测试文件各有失败时，创建 4 个独立 fix task（而非 1 个综合 task）
- [ ] 每个 fix task 的 description 包含该测试文件相关的输出行（包含该文件路径的行及上下文）
- [ ] `maxFixTasksPerStep` 变量和 `countFixTasks` 函数被移除
- [ ] 无法识别测试文件时 fallback 创建按目录分组的 fix task，行为与改动前一致
- [ ] 现有 compile/fmt/lint/unit-test 步骤的 fix task 创建不受影响
- [ ] 支持至少 5 种语言的测试文件命名约定（Go、Python、JS/TS、Java、Ruby）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
