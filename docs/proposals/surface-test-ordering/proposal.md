---
created: "2026-05-25"
author: "faner"
status: Draft
---

# Proposal: Surface Test Ordering & Journey Unification

## Problem

多 surface 项目（如 frontend: web + backend: api）的测试管道无法表达跨 surface 的执行顺序，导致 fail-fast 反馈缺失——API 测试失败后仍继续运行 e2e 测试，浪费 CI 资源并产生噪音。

同时，gen-journeys 按 surface type 拆分为并行任务，与 Journey 的语义定义（用户工作流，天然跨 surface）矛盾，且带来重复提取和并发写入风险。

### Evidence

- `autogen.go` 中 gen-journeys 按 surface type 循环创建并行任务，但 SKILL.md HARD-RATE 声明 "Output is organized by Journey (user workflow), NOT by interface type"
- gen-journeys 是纯叙事提取（读 PRD + 写 MD），不读代码，token 预算小（20-30min），并行收益≈0
- 两个独立 gen-journeys 任务从同一 PRD 提取 Journey，边界划分可能不一致（web 提取 5 个，api 提取 4 个），下游 gen-contracts 合并困难
- run-tests 是单任务，内部无 surface 排序机制，API 和 web 测试无序执行

### Urgency

surface-aware-justfile 提案已 approved，涉及 task struct 和 run-tests 重构。本提案应在 justfile 提案实现前落地，避免重复实现和依赖链冲突。

## Proposed Solution

两项改动：

1. **gen-journeys 合并为单任务**：从 per-surface 并行任务改为单任务，内部遍历所有 surface type 加载对应规则，保持 Journey 跨 surface 的工作流完整性。

2. **run-tests 拆分为 per-surface-key 有序任务**：将 `T-test-run` 拆分为 `T-test-run-{surface-key}`，通过 `execution-order` 配置实现串行依赖。默认优先级约定覆盖常见场景，同类型冲突时要求显式配置。

### Innovation Highlights

混合排序策略（约定 + 覆盖）平衡了零配置体验和灵活性：大多数 fullstack 项目直接受益于默认优先级（api > web），而多 api surface 等边缘场景通过显式配置解决。

## Requirements Analysis

### Key Scenarios

- **典型 fullstack**：`surfaces: { frontend: web, backend: api }`，无显式配置 → 默认 api 先于 web 执行 run-test
- **多 api surface**：`surfaces: { auth-service: api, payment-service: api, admin: web }` → 检测到同类型冲突，提示用户配置 `execution-order`
- **单 surface**：`surfaces: api` → gen-journeys 和 run-test 退化为单任务，行为不变
- **上游失败传播**：`T-test-run-backend` 失败 → `T-test-run-frontend` 状态变为 blocked，跳过执行

### Non-Functional Requirements

- **向后兼容**：单 surface 项目的任务结构和依赖链不变
- **最小改动**：只修改 `autogen.go` 依赖解析和 `config.go`，不动 gen-journeys/gen-contracts/gen-scripts 的 SKILL.md 核心逻辑

### Constraints & Dependencies

- 需与 surface-aware-justfile 提案协调 task struct 改动（SurfaceKey、SurfaceType 字段）
- gen-contracts 保持单任务不变，依赖链不涉及
- gen-scripts 保持 per-surface 并行不变

## Alternatives & Industry Benchmarking

### Industry Solutions

CI 系统通常通过 job dependency graph 表达跨服务测试顺序（如 GitHub Actions 的 `needs` 字段、GitLab CI 的 `dependencies`）。这验证了任务级依赖是业界标准做法。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 手动改 index.json，每次重新生成覆盖 | Rejected: 不可持续 |
| 执行级排序（run-test 内部） | 内部方案 | 改动最小 | 对调度器不可见，无法可视化 | Rejected: 用户已否决 |
| Post-gen 依赖注入 | 内部方案 | 只改依赖解析 | gen-scripts 排序无实际意义 | Rejected: 语义不清晰 |
| **任务级依赖 + 混合排序** | GitHub Actions `needs` | 可视化、调度器原生、失败传播 | 需改 autogen.go | **Selected: 用户确认** |

## Feasibility Assessment

### Technical Feasibility

- `autogen.go` 已有 per-surface 循环逻辑，改为单任务只需移除循环
- `config.go` 新增 `ExecutionOrder []string` 字段是增量改动
- 依赖解析函数已有成熟的模式，新增串行链路无技术障碍

### Resource & Timeline

预计 3-5 个 coding task，属于 small feature 范畴。

### Dependency Readiness

surface-aware-justfile 提案已 approved 但未实现，task struct 改动可由本提案先行引入。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| gen-journeys 按 surface 拆分是因为规则隔离需要 | 5 Whys + 代码分析 | Surface 规则对 gen-journeys 的影响有限，主要是 mandatory outcomes 和 test ratio，这些是给下游阶段用的 |
| 并行 gen-journeys 提升效率 | Occam's Razor | gen-journeys 是 IO 密集型，token 预算小，并行收益≈0，subagent 启动开销可能抵消收益 |
| run-test 必须是单任务 | Assumption Flip | 拆分为 per-surface 任务后，依赖链可表达排序，失败传播自然发生 |

## Scope

### In Scope

- gen-journeys 从 per-surface 并行改为单任务（`autogen.go` 改动）
- run-tests 从单任务拆分为 per-surface-key 串行任务（`autogen.go` 改动）
- 新增 `execution-order` 配置字段（`config.go`）
- 默认优先级约定：api > web > cli > tui > mobile
- 同类型冲突检测：多个同类型 surface 时报错提示显式配置
- 更新 `resolveBreakdownDeps` 和 `resolveQuickDeps` 依赖链
- `T-test-verify-regression` 依赖所有 run-test 子任务
- 失败传播：上游 run-test 失败 → 下游 blocked

### Out of Scope

- gen-contracts 改动（保持单任务）
- gen-scripts 改动（保持 per-surface 并行）
- Surface-aware justfile 集成（独立提案）
- 部分继续运行（一个 surface 失败后继续跑其余的）
- 动态/运行时排序配置
- 可视化跨 surface 依赖

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| gen-journeys 单任务加载多 surface 规则增加 context 噪音 | M | L | gen-journeys 只做信息参考，不做 surface-specific 生成，噪音影响有限 |
| 与 surface-aware-justfile 提案的 task struct 改动冲突 | M | M | 本提案先行引入 SurfaceKey 字段，justfile 提案复用 |
| 默认优先级不覆盖所有场景（如 tui + cli 项目） | L | L | 用户可通过 execution-order 显式覆盖 |
| gen-journeys SKILL.md 需适配多 surface 内部遍历 | M | M | SKILL.md 已有 surface 检测逻辑，扩展为多 surface 遍历 |

## Success Criteria

- [ ] 配置 `surfaces: { frontend: web, backend: api }` 且无 `execution-order` 时，`T-test-run-backend` 的依赖链排在 `T-test-run-frontend` 之前
- [ ] `T-test-run-backend` 失败时，`T-test-run-frontend` 状态为 blocked，不执行
- [ ] gen-journeys 生成单个 `T-test-gen-journeys` 任务，内部加载所有 surface type 的规则
- [ ] 同类型冲突场景（2 个 api surface）在 `forge task build` 时报错，提示配置 `execution-order`
- [ ] 单 surface 项目（`surfaces: api`）的任务结构和依赖链与改动前完全一致
- [ ] `execution-order` 配置验证：引用不存在的 surface-key 时报错

## Next Steps

- Proceed to `/write-prd` to formalize requirements
