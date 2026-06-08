---
created: "2026-06-08"
author: "fanhuifeng"
status: Draft
intent: "enhancement"
---

# Proposal: Test Pipeline Interleaved Execution + Hardened AC

## Problem

当前测试 pipeline 的依赖链是"所有 gen-scripts → 所有 run-tests"，导致：

1. **API bug 延迟暴露**：后端问题要等到所有 surface 的测试脚本都生成完才能发现，浪费了前面脚本的生成成本
2. **E2E 测试缺少 API 反馈**：前端测试脚本基于未经验证的 API 行为生成，容易产生假失败
3. **test-run 任务缺少严格 AC**：没有强制要求测试必须真实通过、不能是假测试，也没有限制随意修改正式代码

### Evidence

以 `backend=api, frontend=web` 为例，当前依赖链：

```
T-test-gen-scripts-backend   ← eval-contract
T-test-gen-scripts-frontend  ← gen-scripts-backend
T-test-run-backend           ← gen-scripts-frontend  (等前端脚本也生成完)
T-test-run-frontend          ← run-tests-backend
```

这意味着后端测试必须等前端脚本生成完才能执行。如果后端有 API bug，前端脚本也会基于错误行为生成。

### Urgency

多 surface 项目中，这种串行生成-串行执行的浪费随 surface 数量线性增长。修复成本低（改依赖链接线+模板），收益明确（更早发现问题、减少无效生成）。

## Proposed Solution

两处改动：

1. **依赖链交错**：将 gen-scripts 和 run-tests 按 surface 配对执行（gen-A → run-A → gen-B → run-B）
2. **模板强化**：在 prompt 模板中写入关键执行指令，在 auto-gen 模板中补充硬性 AC

### Innovation Highlights

业界 CI/CD 的标准做法是"分层测试，逐层卡点"（API 先行，E2E 后行）。此处将其应用到 Forge 的任务 DAG 中，通过依赖链重排实现，无需引入新机制。

## Requirements Analysis

### Key Scenarios

- **多 surface 项目**（backend=api + frontend=web）：API 测试先跑，修复后再生成 Web 测试脚本
- **单 surface 项目**：依赖链自然退化为 gen → run，无影响
- **三 surface 项目**（backend=api + frontend=web + cli=cli）：按 execution_order 依次配对执行

### Non-Functional Requirements

- 向后兼容：单 surface 项目行为不变
- 不影响 task-executor 的现有 Error Handling 协议

### Constraints & Dependencies

- 依赖链改在 Go 代码中（`pipeline.go` 的 `GenerateTestTasks`）
- AC 改在 auto-gen 模板中（`pkg/task/templates/test-run.md`）
- 执行指令改在 prompt 模板中（`pkg/prompt/templates/test-run.md`）

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD pipeline 中标准做法是 stage 串行（build → unit → integration → e2e），每个 stage 内可并行。Forge 的 per-surface 扩展本质上就是 multi-stage。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无成本 | 浪费生成成本、问题延迟暴露 | Rejected: 收益明确且成本低 |
| 交错执行 | CI/CD stage 模式 | 更早发现问题、减少无效生成 | 多 surface 串行执行时间更长（但发现问题的总时间更短） | **Selected: 成本最低、收益最直接** |
| 并行执行所有 surface | 并行 CI | 总时间最短 | 无法利用前序 surface 反馈 | Rejected: 丧失信息传递优势 |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动集中在：
- `pipeline.go` 的 per-surface-key 依赖接线逻辑（~15 行）
- 两个模板文件（纯内容变更）

### Resource & Timeline

单人半天可完成。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 所有 gen-scripts 必须在 run-tests 之前完成 | Assumption Flip: 如果 run-A 在 gen-B 之前执行，gen-B 可以利用 run-A 的反馈 | Confirmed: 交错执行信息效率更高 |
| test-run 的 agent 可以自由修改正式代码 | Stress Test: 如果测试本身写错了呢？ | Refined: 必须先确认是正式代码问题，不能盲目修改 |

## Scope

### In Scope

- 修改 `GenerateTestTasks` 的依赖接线，实现 per-surface gen→run 交错
- 在 auto-gen 模板 `test-run.md` 中补充硬性 AC（测试必须通过、必须是真实测试、确认是正式代码问题才能改正式代码）
- 在 prompt 模板 `test-run.md` 中补充关键执行指令（与 task-executor Error Handling 协调）

### Out of Scope

- task-executor agent 定义本身（不改）
- quality-gate 机制（不改）
- gen-scripts 模板（不改）
- 单 surface 项目的任何行为变化

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 交错执行使多 surface 项目串行执行时间变长（因为不能并行跑不同 surface 的测试） | M | L | execution_order 已由用户配置控制，用户可以选择让 fast surface 先跑 |
| prompt 模板中的新指令与 task-executor 的 Pause Protocol 冲突 | L | M | 新指令作为 TASK-CONSTRAINTS 级别的补充，不覆盖 task-executor 的 EXTREMELY-IMPORTANT 层级 |
| 现有测试需要更新以适应新的依赖链 | M | L | 更新 pipeline 的单元测试用例 |

## Success Criteria

- [ ] 多 surface 项目中，`T-test-run-{surface-N}` 依赖 `T-test-gen-scripts-{surface-N}`，而非 `T-test-gen-scripts-{surface-N+1}`
- [ ] 多 surface 项目中，`T-test-gen-scripts-{surface-N}` (N>0) 依赖 `T-test-run-{surface-N-1}`，而非 `T-test-gen-scripts-{surface-N-1}`
- [ ] 单 surface 项目的依赖链不变（gen → run）
- [ ] test-run 任务模板包含 AC：所有测试用例必须通过
- [ ] test-run 任务模板包含 AC：必须是真实测试，不能是假测试
- [ ] test-run prompt 模板包含指令：确认是正式代码问题才能修改正式代码；测试脚本本身的 bug 可以修，但不能为了通过测试而篡改测试逻辑
- [ ] test-run prompt 模板包含指令：问题多时通过 `forge task add` 追加 fix 任务（与 task-executor Pause Protocol 协调）
- [ ] 所有现有单元测试通过

consistency_check_result:
  status: pass
  pairs_checked: 18
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements
