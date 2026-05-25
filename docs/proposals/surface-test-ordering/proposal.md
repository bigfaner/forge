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

<!-- pre-revised: high -->
1. **gen-journeys 合并为单任务**：从 per-surface 并行任务改为单任务，内部遍历所有 surface type 加载对应规则，保持 Journey 跨 surface 的工作流完整性。此改动同时影响 breakdown 和 quick 两种模式（`autogen.go` 中 `GetBreakdownTestTasks` 和 `GetQuickTestTasks` 均包含 per-type gen-journeys 循环）。

2. **run-tests 拆分为 per-surface-key 有序任务**：将 `T-test-run` 拆分为 `T-test-run-{surface-key}`，通过 `execution-order` 配置实现串行依赖。默认优先级约定覆盖常见场景，同类型冲突时要求显式配置。

<!-- pre-revised: high -->
> **命名策略**：使用 surface-key（YAML map 的 key，如 `backend`、`frontend`）作为任务后缀，而非 surface-type（如 `api`、`web`）。理由：surface-key 在项目中唯一标识一个 surface，surface-type 不唯一（如多个 api surface）。`execution-order` 引用的也是 surface-key。

<!-- pre-revised: medium -->
#### 3-Surface 依赖链示例

配置 `surfaces: { auth-service: api, admin: web, cli: cli }`，`execution-order: [auth-service, admin, cli]`：

**Breakdown 模式：**
```
T-test-gen-journeys
    └─ T-test-gen-contracts
        └─ T-test-gen-scripts-auth-service
        └─ T-test-gen-scripts-admin
        └─ T-test-gen-scripts-cli
            └─ T-test-run-auth-service
                └─ T-test-run-admin
                    └─ T-test-run-cli
                        └─ T-test-verify-regression
```

**Quick 模式：**
```
T-test-gen-journeys
    └─ T-test-run-auth-service
        └─ T-test-run-admin
            └─ T-test-run-cli
                └─ T-test-verify-regression
```

### Innovation Assessment

本方案的核心机制——约定优先 + 显式覆盖——是业界成熟的 convention-over-configuration 模式（Rails, 2005），并非创新。真正的设计价值在于两点工程决策：(1) gen-journeys 合并，利用 Journey 跨 surface 的语义本质消除了不必要的并行拆分；(2) 单 surface 退化规则（scalar 形式不添加后缀），在零配置场景下保持向后兼容。

## Requirements Analysis

### Key Scenarios

- **典型 fullstack**：`surfaces: { frontend: web, backend: api }`，无显式配置 → 默认 api 先于 web 执行 run-test
- **多 api surface**：`surfaces: { auth-service: api, payment-service: api, admin: web }` → 检测到同类型冲突，提示用户配置 `execution-order`
<!-- pre-revised: high -->
- **单 surface（scalar 形式）**：`surfaces: api` → 退化为无后缀 `T-test-run`（非 `T-test-run-api`），gen-journeys 同理退化为无后缀单任务，行为与改动前完全一致
- **上游失败传播**：`T-test-run-backend` 失败 → `T-test-run-frontend` 状态变为 blocked，跳过执行

### Non-Functional Requirements

- **向后兼容**：单 surface 项目的任务结构和依赖链不变
- **改动范围**：涉及 `autogen.go`（任务生成、依赖链、迁移）、`config.go`（ExecutionOrder 字段、surface-key 校验）、`infer.go`（InferType 前缀匹配）及 `renderBody` 模板（空 TestType 适配）。gen-contracts 和 gen-scripts 的 SKILL.md 核心逻辑不动
- **命名一致性说明**：run-tests 使用 surface-key 后缀（`T-test-run-{key}`），gen-scripts 继续使用 surface-type 后缀（`T-test-gen-scripts-{type}`）。两套命名共存是因为 gen-scripts 按 type 并行生成（同类型 surface 共享生成规则），而 run-tests 需按 key 独立执行（同类型 surface 的测试不可合并）。gen-journeys 合并后无后缀。用户可见的任务列表中三种命名方案同时出现，这是当前的设计取舍

### Constraints & Dependencies

- 需与 surface-aware-justfile 提案协调 task struct 改动（SurfaceKey、SurfaceType 字段）
- gen-contracts 保持单任务不变，依赖链不涉及
- gen-scripts 保持 per-surface 并行不变
<!-- pre-revised: high -->
- **surface-key 合法性约束**：surface-key 必须匹配 `[a-z][a-z0-9-]*`（小写字母开头，仅含小写字母、数字、连字符）。YAML map 中不符合规则的 key 在 config load time 归一化（空格和特殊字符替换为 `-`，大写转小写），归一化后仍不合法则报错
- **验证时机**：所有配置校验（`execution-order` 引用校验、surface-key 合法性检查、同类型冲突检测）均在 config load time 执行（fail fast），不推迟到 build time

## Alternatives & Industry Benchmarking

### Industry Solutions

CI 系统通常通过 job dependency graph 表达跨服务测试顺序（如 GitHub Actions 的 `needs` 字段、GitLab CI 的 `dependencies`）。这验证了任务级依赖是业界标准做法。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 手动改 index.json，每次重新生成覆盖 | Rejected: 不可持续 |
| 执行级排序（run-test 内部） | 内部方案 | 改动最小 | 对调度器不可见，无法可视化 | Rejected: 用户已否决 |
| Post-gen 依赖注入 | 内部方案 | 只改依赖解析 | gen-scripts 排序无实际意义 | Rejected: 语义不清晰 |
| **任务级依赖 + 混合排序** | GitHub Actions `needs` | 可视化、调度器原生、失败传播 | 详见下方 Cons 明细 | **Selected: 技术分析见下** |

**Selected Approach 技术论证**：执行级排序（方案 2）在代码改动量上更小，但对调度器不透明——调度器只能看到单个 `T-test-run` 任务，无法在任务列表中展示各 surface 的独立状态，也无法在任一 surface 失败时自动 block 其余 surface。任务级依赖让每个 surface 的测试成为独立可见的调度单元，失败传播由调度器原生处理（blocked 状态），无需在 run-test 内部实现条件逻辑。代价是改动面更广（见下方 Cons 明细），但换来的是调度语义的正确性和用户体验的一致性。

**Selected Approach Cons 明细**：
- `autogen.go`：`GetBreakdownTestTasks` 和 `GetQuickTestTasks` 函数签名需从 `capabilities []string` 改为接收 surfaces map 或 surface-key 列表，所有调用方需适配
- `infer.go`：`InferType` 函数从 `T-test-run` 精确匹配改为 `typeSuffixedID` 前缀匹配，需处理歧义风险（未来新增以 `T-test-run-` 开头的非 surface 任务 ID）
- `autogen.go`：新增 index.json 迁移逻辑（在 BuildIndex 阶段，非 GetBreakdownTestTasks），将 `T-test-run` 的 fix-tasks 重映射到 `T-test-run-{surface-key}`
- 命名方案：gen-scripts 继续使用 type 后缀（`-api`、`-web`），run-tests 使用 key 后缀（`-backend`、`-frontend`），两套命名共存（详见 Key Risks）

## Feasibility Assessment

### Technical Feasibility

- `autogen.go` 已有 per-surface 循环逻辑，改为单任务只需移除循环
- `config.go` 新增 `ExecutionOrder []string` 字段是增量改动
- 依赖解析函数已有成熟的模式，新增串行链路无技术障碍
- **函数签名变更**：`GetBreakdownTestTasks` 和 `GetQuickTestTasks` 当前接收 `capabilities []string`（去重后的 type 列表如 `["api", "web"]`），需改为接收 surfaces map 或 surface-key 列表以支持 per-surface-key 任务生成。所有调用方（含 `BuildIndex` 及相关入口函数）需适配新签名
- **InferType 变更**（`infer.go`）：`T-test-run` 的精确匹配改为 `typeSuffixedID` 前缀匹配，含对应测试用例更新
- **迁移逻辑位置**：index.json 中已有 `SourceTaskID: "T-test-run"` 的 fix-tasks 重映射逻辑应放在 `BuildIndex` 阶段（非 `GetBreakdownTestTasks`），因为 `BuildIndex` 拥有 index.json 读写权限和旧任务状态上下文
- **受影响代码路径清单**：`autogen.go`（gen-journeys 循环移除、run-tests per-key 生成、verify-regression 依赖链尾）、`config.go`（ExecutionOrder 字段、surface-key 校验）、`infer.go`（InferType 前缀匹配）、`renderBody` 相关模板（空 TestType 适配）

### Resource & Timeline

预计 5-7 个 coding task，属于 small-to-medium feature 范畴。增加的工作量主要来自函数签名变更（影响所有调用方）、InferType 前缀匹配、index.json 迁移逻辑和 surface-key 校验归一化。

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

<!-- pre-revised: medium -->
- gen-journeys 从 per-surface 并行改为单任务（`autogen.go` 改动），合并后 TestType 字段留空（原 per-type 任务的 TestType 值不再适用），`renderBody` 函数适配空 TestType 场景
- run-tests 从单任务拆分为 per-surface-key 串行任务（`autogen.go` 改动）
- 新增 `execution-order` 配置字段（`config.go`）
- 默认优先级约定：api > web > cli > tui > mobile。未覆盖的组合（如 tui + cli）按 config 中的 surface 声明顺序排列（YAML map 的 key 顺序）
- 同类型冲突检测：多个同类型 surface 时报错提示显式配置
- 更新 `resolveBreakdownDeps` 和 `resolveQuickDeps` 依赖链
<!-- pre-revised: high -->
- 更新 `InferType` 函数：将 `T-test-run` 的精确匹配改为 `typeSuffixedID` 前缀匹配，以识别 `T-test-run-{surface-key}`（含 `infer_test.go` 测试用例更新）
- `T-test-verify-regression` 依赖 execution-order 中最后一个 run-test 子任务（即链尾），而非全部
- 失败传播：上游 run-test 失败 → 下游 blocked
<!-- pre-revised: high -->
- 迁移步骤：在 `BuildIndex` 阶段（非 task 生成函数），检测 `index.json` 已有 `SourceTaskID: "T-test-run"` 的 fix-tasks，自动重映射到对应的 `T-test-run-{surface-key}`（单 surface 场景保持 `T-test-run` 不变，无迁移成本）

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
| gen-journeys 单任务加载多 surface 规则增加 context 噪音 | M | L | gen-journeys 的规则加载仅为参考信息（mandatory outcomes、test ratio），不直接影响生成输出。SKILL.md 中增加 `## Multi-Surface Rules Loading` 段落，明确按 surface type 分节加载，避免规则混淆 |
| 与 surface-aware-justfile 提案的 task struct 改动冲突 | M | M | 本提案先行引入 SurfaceKey 字段，justfile 提案复用。需在实现前对齐 surface-key 定义：justfile 提案使用 `/` 到 `+` 转换，本提案使用原始 YAML map key——统一为 config load time 归一化后的值（`/` 转 `-`），两提案共用同一归一化函数 |
| 默认优先级不覆盖所有场景（如 tui + cli 项目） | L | L | 用户可通过 execution-order 显式覆盖 |
| gen-journeys SKILL.md 需适配多 surface 内部遍历 | M | M | SKILL.md 已有 surface 检测逻辑，扩展为多 surface 遍历。验收标准：生成输出的 Journey 文件覆盖所有配置 surface，由 SC3 覆盖 |
<!-- pre-revised: high -->
| index.json 已有 `T-test-run` 条目变为孤儿 | M | M | 迁移时将旧 `T-test-run` 的 status/blocked-reason 复制到 `T-test-run-{first-surface-key}`；单 surface 项目无此问题；多 surface 项目需按 execution-order 首个 surface 继承状态。迁移逻辑放在 BuildIndex 阶段 |
| InferType 前缀匹配引入歧义 | M | M | 前缀匹配 `T-test-run-` 后的片段作为 surface-key 查找 surfaces map；若未命中任何已知 key，回退到原有精确匹配逻辑。新增 InferType 单元测试覆盖：已知 key、未知 key、单 surface 退化三种场景 |
| 函数签名变更影响所有调用方 | M | L | `capabilities []string` → `surfaceKeys []string` 或新增 `surfaces map[string]string` 参数。逐个函数修改，编译器类型检查确保无遗漏调用方 |
| gen-scripts type 后缀与 run-tests key 后缀并存造成用户困惑 | L | M | 文档中明确说明命名差异的原因（gen-scripts 按 type 并行共享规则 vs run-tests 按 key 独立执行），任务列表中任务标题包含 surface-key 以便区分 |

## Success Criteria

- [ ] 配置 `surfaces: { frontend: web, backend: api }` 且无 `execution-order` 时，`T-test-run-backend` 的依赖链排在 `T-test-run-frontend` 之前
- [ ] `T-test-run-backend` 失败时，`T-test-run-frontend` 状态为 blocked，不执行
- [ ] gen-journeys 生成单个 `T-test-gen-journeys` 任务，输出覆盖所有配置 surface 的 Journey 文件（非仅单个 surface）
- [ ] 同类型冲突场景（2 个 api surface）在 config load time 报错，提示配置 `execution-order`
<!-- pre-revised: high -->
- [ ] 单 surface 项目（`surfaces: api`）退化为无后缀 `T-test-run`（非 `T-test-run-api`），任务 ID 和依赖列表与改动前一致
- [ ] `execution-order` 配置验证：引用不存在的 surface-key 时在 config load time 报错
- [ ] Quick 模式：`surfaces: { frontend: web, backend: api }` 无 `execution-order` 时，`T-test-run-backend` 仍排在 `T-test-run-frontend` 之前，且 `T-test-gen-journeys` 为直接上游
- [ ] `InferType("T-test-run-backend")` 返回正确的 surface type（`api`），含前缀匹配而非精确匹配
- [ ] 迁移正确性：多 surface 项目 index.json 中已有 `SourceTaskID: "T-test-run"` 的 fix-tasks 在 BuildIndex 阶段自动重映射到 `T-test-run-{execution-order 首个 surface-key}`
- [ ] Surface-key 校验：`surfaces: { "ADMIN PANEL": web }` 归一化为 `admin-panel` 通过；`surfaces: { "123bad": web }` 在 config load time 报错
- [ ] 默认优先级：`surfaces: { mobile: mobile, cli: cli, web: web, api: api }` 无 `execution-order` 时，执行顺序为 api → web → cli → mobile

## Next Steps

- Proceed to `/write-prd` to formalize requirements
