---
created: "2026-05-28"
author: fanhuifeng
status: Draft
---

# Proposal: Regression Fix Task 按 Test Suite 拆分

## Problem

Quality-gate hook 检测到 `just test` 失败后创建单个 fix task 覆盖所有 test suite 的全部失败，当失败数量多（20+）且跨多个不相关 suite 时，agent 执行时卡住。

### Evidence

实际发生：quality-gate 创建的 fix-1 包含 4 个 test suite 共 20+ 失败，agent 执行后长时间无响应被用户手动中断。Concise error 只展示输出尾部（test-suite-health），agent 需要读 raw-output.txt 才能看到完整失败列表，且无法区分"我引入的"和"预存在的"失败。详见 `docs/lessons/gotcha-fix-task-broad-scope.md`。

### Urgency

高。每次 regression 测试出现多 suite 失败都会触发此问题，agent 卡死浪费时间和 token。这是一个可机械解决的问题（纯 shell 解析，无需 LLM 分析）。

## Proposed Solution

在 `quality_gate.go` 中新建 `addRegressionFixTasks` 函数，从 test output 中按 `FAIL` 行提取失败 suite 列表，为每个 suite 独立创建 fix task，每个 task 包含该 suite 的完整失败详情（全部 `--- FAIL` 行和错误消息）。

`runTestRegression` 失败时调用新函数替代原有 `addFixTask`。解析失败时 fallback 到现有单 task 行为。

同时移除 `maxFixTasksPerStep` cap 限制——拆分后的每个 fix task scope 已收窄，cap 不再必要。

### Innovation Highlights

无创新。标准的错误分组策略：按 test suite 边界机械拆分，不引入 LLM 分析。与 CI 系统中按 test module 分组失败报告的常规做法一致。

## Requirements Analysis

### Key Scenarios

- **Happy path**：4 个 suite 各有 2-5 个失败 → 创建 4 个 fix task，每个包含对应 suite 的完整失败详情
- **单 suite 失败**：1 个 suite 有 3 个失败 → 创建 1 个 fix task，行为与当前一致
- **输出格式不符**：无法解析出 suite 信息 → fallback 到现有单 task 行为
- **全部通过**：无失败 → 不创建任何 fix task（现有行为不变）

### Non-Functional Requirements

- **性能**：解析 test output 的时间应可忽略（纯字符串匹配，output 通常 < 1MB）
- **兼容性**：不影响 compile/fmt/lint/unit-test 步骤的现有逻辑

### Constraints & Dependencies

- 依赖 Go test 的 `-v` 输出格式：`FAIL <package-path>` 行标识失败 suite，`--- FAIL: <test-name>` 行标识失败测试
- 当前 `just test` 已使用 `-v` 标志，输出格式稳定

## Alternatives & Industry Benchmarking

### Industry Solutions

CI 系统普遍按 test module/suite 分组报告失败（GitHub Actions 的 test grouping、JUnit XML 的 testsuite 元素）。本质相同：按结构性边界拆分错误报告。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 卡死问题持续 | Rejected: 用户体验差，实际已发生 |
| 在 addFixTask 内分支处理 | 内部方案 | 改动集中 | addFixTask 职责膨胀 | Rejected: 违反单一职责 |
| **独立函数 addRegressionFixTasks** | 标准分组策略 | 职责清晰，现有逻辑零风险 | 多一个函数 | **Selected: 最小风险，最大清晰度** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go 标准库 `bufio.Scanner` + `strings.HasPrefix` 即可完成解析，无需外部依赖。当前 `addFixTask` 已有类似的字符串解析逻辑（`extractSourceFiles`），新增解析代码符合现有模式。

### Resource & Timeline

预计 1-2 小时实现 + 测试。改动集中在 `quality_gate.go`，影响面小。

### Dependency Readiness

无外部依赖。`just test` 输出格式由 Go test toolchain 保证稳定性。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| fix task cap (maxFixTasksPerStep=3) 是必要的防护 | 5 Whys: cap 是因为单个 fix task scope 过大导致循环，拆分后 scope 已收窄 | Overridden: 拆分后每个 task scope 窄，cap 不再必要，用户明确要求移除 |

## Scope

### In Scope

- 新建 `addRegressionFixTasks` 函数，解析 test output 提取失败 suite
- 每个失败 suite 独立创建 fix task，包含该 suite 的完整失败详情
- 移除 `maxFixTasksPerStep` cap 限制（变量、计数逻辑、错误返回）
- `runTestRegression` 失败时调用新函数替代 `addFixTask`
- 解析失败时 fallback 到现有单 task 行为
- 更新或新增单元测试覆盖新逻辑

### Out of Scope

- 基线对比过滤预存在失败（lesson 第二层改进）
- unit-test / compile / lint 步骤的改动
- Surface inference 改进
- `coding.fix.md` 模板变更（模板内容足够，无需修改）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Go test 输出格式变化导致解析失败 | L | L | fallback 到现有单 task 行为，零功能损失 |
| 移除 cap 后大量 fix task 并发 | L | M | 每个 fix task scope 已收窄到单 suite，并发执行风险可控 |
| Surface inference 对 test suite 路径匹配不准 | M | L | 保持现有 inferSurface 逻辑，soft-failure 容错 |

## Success Criteria

- [ ] 4 个 suite 各有失败时，创建 4 个独立 fix task（而非 1 个综合 task）
- [ ] 每个 fix task 的 description 包含该 suite 的全部 `--- FAIL` 行和错误消息
- [ ] `maxFixTasksPerStep` 变量和 `countFixTasks` 函数被移除
- [ ] 解析失败时 fallback 创建 1 个 fix task，行为与改动前一致
- [ ] 现有 compile/fmt/lint/unit-test 步骤的 fix task 创建不受影响

## Next Steps

- Proceed to `/write-prd` to formalize requirements
