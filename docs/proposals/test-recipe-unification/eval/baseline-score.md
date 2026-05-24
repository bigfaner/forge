---
iteration: 0
scorer: CTO adversary
model: baseline-snapshot proposal + freeform-review
date: 2026-05-24
---

# Baseline Score Report

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem -> Solution**: 问题链清晰——三层错位（submit 门禁过慢、模板硬编码 Playwright、配置命名过时）。两层 recipe 模型（unit-test + test）直接解决了三层问题。论证链成立。

**Solution -> Evidence**: 方案引用了具体代码文件和行号，证据充分。但 freeform review 发现 `addFixTask` 中 step=="unit-test" -> "just test" 的映射问题未被提案捕捉，说明对现有代码的审计有遗漏。

**Evidence -> Success Criteria**: 成功标准与三层问题对应，但存在内部矛盾——"无 Fallback"声明与 `RunProjectTests` 现有的五级探测链冲突。

**Self-contradiction check**:
1. "无 Fallback" vs `RunProjectTests` 的 fallback chain
2. `DefaultGateSequence()` 改名 vs 新增 `UnitGateSequence()` 之间的语义重叠
3. Impact Analysis 遗漏了 `internal/cmd/test/test.go` 和 `internal/cmd/config_test.go`
4. `quality_gate.go` 中的 `runE2ERegression` 整个函数（约 50 行）在 Impact Analysis 中仅一笔带过

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

**Problem stated clearly (34/40)**: 三层错位的描述清晰——submit 门禁跑全量测试浪费时间、模板硬编码 Playwright、配置项命名过时。但"三层错位"这一概括有轻微模糊——第一层（submit 跑全量）是性能问题，第二层（模板硬编码 Playwright）是功能问题，第三层（命名过时）是维护性问题，它们的优先级和关联性未明确。扣 6 分。

**Evidence provided (36/40)**: 5 条证据，每条都有具体代码引用和量化数据（~90s）。但"10 个 breaking 任务累计 15 分钟"是基于假设的线性外推，未提供实际项目中的真实数据支撑。扣 4 分。

**Urgency justified (28/30)**: 量化了延迟成本（90s/次），说明了 v3.0.0 窗口期。但"避免后续返工"的论据较模糊——具体哪些后续工作会被阻塞未说明。扣 2 分。

**Subtotal: 98/110**

### 2. Solution Clarity (120 pts)

**Approach is concrete (36/40)**: 两层 recipe 模型定义清晰——`unit-test` vs `test` 的职责和调用时机都明确。扣 4 分因为 `DefaultGateSequence` 和 `UnitGateSequence` 的关系未在方案中精确定义（freeform review 已指出）。

> "DefaultGateSequence() 的 test -> unit-test；新增 UnitGateSequence()"

两个函数并存且职责不清。

**User-facing behavior described (40/45)**: submit 门禁和 all-completed 的行为描述清楚。扣 5 分因为用户手动执行 `just test` 的语义发生了倒置（从单元测试变为高级测试），这一 UX 影响仅在 freeform review 中被提及，proposal 本身未警告用户。

**Technical direction clear (32/35)**: Impact Analysis 五层分级详尽，文件级改动清楚。扣 3 分因为 `RunProjectTests` 探测链如何适配两层模型的细节缺失。

**Subtotal: 108/120**

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (28/40)**: 仅引用了 Test Pyramid（Kent Beck/Martin Fowler），这是分层测试的基础概念。缺少对具体 CI/CD 系统如何实现分层门禁的引用——例如 GitHub Actions 的 job-level gate、GitLab CI 的 stage-level gate、Bazel 的 test size classes（small/medium/large）。扣 12 分。

**At least 3 meaningful alternatives (22/30)**: 4 个选项包括"do nothing"。但"统一 + 缓存 + 并行"是稻草人——它被定义为"过设计"然后被拒绝，无任何实际方案细节，只是为了衬托两层模型。扣 8 分。

> "统一 + 缓存 + 并行 | — | 最大性能提升 | 复杂度过高，缓存机制引入新 bug 面 | Rejected: 过度设计"

**Honest trade-off comparison (18/25)**: 选中方案的 Cons 仅写"需更新多个组件"，这是实现成本而非架构权衡。真正的 trade-off——`just test` 语义倒置导致用户肌肉记忆失效——未被承认。扣 7 分。

**Chosen approach justified against benchmarks (20/25)**: "直击痛点，复杂度可控"是结论性陈述而非论证。为什么 Test Pyramid 的三层映射到两层而非三层？为什么不是 unit-test + integration-test + e2e-test？扣 5 分。

**Subtotal: 88/120**

### 4. Requirements Completeness (110 pts)

**Scenario coverage (32/40)**: 6 个关键场景覆盖了主路径。但缺少边界场景：
- 无 justfile 的项目（纯 go mod）submit 时怎么办？`RunProjectTests` 的 fallback 是否仍然有效？
- `unit-test` recipe 不存在但 `test` recipe 存在时的行为
- `unit-test` recipe 执行超时的处理

扣 8 分。

**Non-functional requirements (32/40)**: 性能量化目标清晰（90s -> 20s），向后兼容策略明确（v3.0.0 直接断）。扣 8 分因为缺少：
- 迁移后的错误提示 UX（用户运行旧命令时看到什么？）
- 无 justfile 项目的性能基准

**Constraints & dependencies (25/30)**: 4 条约束都具体。扣 5 分因为 `quality_gate.go` 中 `runE2ERegression` 整个函数（约 50 行代码）的迁移约束未被提及——该函数硬编码了 `just e2e-test`、`just e2e-setup`、`just dev` 等多个 recipe 调用。

**Subtotal: 89/110**

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (20/40)**: Proposal 自己声明"非创新性方案"，直接映射 Test Pyramid。不加分的"创新"是 surface-agnostic 抽象（Forge 不区分 e2e 还是集成测试），但这本质上是将复杂度下推到 recipe 层，不是真正的创新。扣 20 分。

**Cross-domain inspiration (15/35)**: 未引用任何跨域灵感。两层模型纯粹来自测试领域的常识。扣 20 分。

**Simplicity of insight (22/25)**: 将三层错位归约为两层 recipe，方案简洁。扣 3 分因为"无 Fallback"声明实际上引入了隐含的复杂性（`RunProjectTests` 仍然有 fallback，需要区分两条路径的不同策略）。

**Subtotal: 57/100**

### 6. Feasibility (100 pts)

**Technical feasibility (36/40)**: Go 代码改动和模板替换都是确定性工作。扣 4 分因为 `quality_gate.go` 中 `runE2ERegression` 函数的迁移复杂度被低估——它不仅涉及 `e2e-test` -> `test` 的简单重命名，还涉及 `e2e-setup` 和 `dev` recipe 的协调。

**Resource & timeline feasibility (24/30)**: ~42 文件 / 10-12 个 coding task 是合理的估算。扣 6 分因为没有提供时间线估算（天数/周数），仅给出了任务数量。

> "涉及 ~42 个文件，影响范围评估如下。预估 10-12 个 coding task。"

**Dependency readiness (28/30)**: 无外部依赖是正确的。扣 2 分因为 Forge 自身的 justfile 也需要修改，这是"自己给自己做手术"——如果在迁移过程中 Forge 自己的测试 recipe 暂时不可用，开发流程会受影响。

**Subtotal: 88/100**

### 7. Scope Definition (80 pts)

**In-scope items are concrete (28/30)**: 每个文件级改动都列出了具体变更内容。扣 2 分因为 Impact Analysis Tier 1 列出 8 个 Go 文件，但实际代码中还有 `internal/cmd/test/test.go`、`internal/cmd/config_test.go`、`internal/cmd/config_schema_test.go` 未被计入。

**Out-of-scope explicitly listed (22/25)**: 9 条排除项，涵盖了并行/缓存等独立优化。扣 3 分因为"历史 lessons/proposals 中的 e2eTest 引用"被标为 Out of Scope 但 Low priority 文档可能仍会造成混淆。

**Scope is bounded (22/25)**: 任务数量有上限（10-12），但缺少时间边界。扣 3 分。

**Subtotal: 72/80**

### 8. Risk Assessment (90 pts)

**Risks identified (24/30)**: 5 个风险。扣 6 分因为遗漏了以下风险：
- `just test` 语义倒置对用户肌肉记忆的影响
- `quality_gate.go` 中 `handleGateFailure` 的 `guide`/`label` map 仍硬编码 `"e2e-test"` 键
- `internal/cmd/test/test.go` 中 `ExecuteJourneyInIsolation` 的调用链遗漏

**Likelihood + impact rated (22/30)**: 评估大体诚实。扣 8 分因为：
- "改动面大（~42 files）导致遗漏"标为 M/M，但 Impact Analysis 本身已经遗漏了文件（`test.go`、`config_test.go`），说明遗漏概率实际为 H
- "auto.e2eTest -> auto.test" 标为 H/L 但用户静默失败的实际影响可能是 H（CI pipeline 坏掉无人注意）

**Mitigations are actionable (22/30)**: 扣 8 分因为：
- "v3.0.0 要求重新运行 init-justfile 生成新 justfile"——这不是 mitigation 而是 requirement，它没有说明如何确保用户知道要这样做
- "直接重命名，用户运行 forge init 或手动更新即可"——假设用户会注意到，缺少主动通知机制
- "just test 需支持 --feature 参数，模板统一生成"——与实际代码使用 positional argument 不匹配

> "just test 需支持 --feature 参数，模板统一生成"

实际 `ExecuteJourneyInIsolation` 使用 `exec.Command("just", "e2e-test", journeyName)` 传递 positional argument，不是 `--feature` flag。

**Subtotal: 68/90**

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (42/55)**: 8 条标准中 6 条可验证（门禁运行特定 recipe、配置字段迁移、测试通过）。扣 13 分因为：
- "耗时 <30s（Go 项目）"——未指定是首次运行还是稳态运行，是否排除编译缓存
- "完整覆盖，无 e2e-test 调用"——需要全代码库 grep 验证，但未提供验证命令
- "所有 Go 测试通过"——这是基线要求，不是本提案的成功标准

**Coverage is complete (20/25)**: 覆盖了主要 in-scope 项。扣 5 分因为缺少以下 in-scope 项的成功标准：
- init-justfile 模板为各语言生成 `test` recipe 的验证
- `journey_isolation.go` 迁移后 journey 测试仍通过的验证
- Prompt 模板和 skill markdown 更新的验证

**Subtotal: 62/80**

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (30/35)**: 两层模型解决了三层错位中的全部三个问题。扣 5 分因为 `RunProjectTests` 探测链的"无 Fallback"声明与实际代码行为矛盾——代码中 `RunProjectTests` 有完整的 fallback chain，而成功标准要求"无 unit-test recipe 时 quality gate 报错"。

> "无 Fallback：v3.0.0 直接要求 unit-test recipe，不回落到 test"
> "无 unit-test recipe 时 quality gate 报错提示运行 init-justfile（不 fallback）"

**Scope <-> Solution <-> Success Criteria aligned (24/30)**: 基本对齐。扣 6 分因为：
- Scope 列出 `quality_gate.go` 的 `addFixTask` 映射更新，但未在 Solution 中详细说明
- Success Criteria 未覆盖 Scope 中的 prompt 模板和 skill markdown 更新

**Requirements <-> Solution coherent (20/25)**: 扣 5 分因为 `addFixTask` 中 `step == "unit-test"` 映射到 `"just test"` 的现有代码需要修改，这一需求在 Requirements Analysis 中未提及，仅在 Impact Analysis Tier 1 中隐含。Requirements 和 Solution 之间存在缺口。

**Subtotal: 74/90**

## Phase 3: Blindspot Hunt

### [blindspot] `quality_gate.go` 中 `handleGateFailure` 硬编码的 guide/label map 遗漏

当前 `handleGateFailure` 函数（quality_gate.go 第 262-273 行）有硬编码的 map：

```go
guide := map[string]string{
    "compile":   "fix compilation errors",
    "lint":      "fix lint errors",
    "unit-test": "fix failing tests",
    "e2e-test":  "fix failing e2e tests",
}
```

迁移后 `"e2e-test"` 键需改为 `"test"`。这一改动未在 Impact Analysis 的任何 Tier 中提及。如果遗漏，all-completed 门禁的高级测试失败时将显示 "fix the issue" 而非具体指导。

### [blindspot] `quality_gate.go` 中 `runE2ERegression` 整个函数需要重构，而非简单重命名

当前 `runE2ERegression` 函数（quality_gate.go 第 198-251 行）包含约 50 行代码，涉及：
- `just.HasRecipe(projectRoot, "e2e-test")` 探测
- `just e2e-setup` 调用
- `just dev` 启动检查
- `just e2e-test` 执行
- `addFixTask(... "e2e-test" ...)` 调用

这不是简单的字符串替换，而是需要重构整个函数。Impact Analysis 中仅列了一行 `Step 3 用 test`，严重低估了复杂度。

### [blindspot] `internal/cmd/config_test.go` 和 `internal/cmd/config_schema_test.go` 遗漏

这两个文件包含大量 `e2eTest` 断言（config_test.go 约 6 处，config_schema_test.go 约 4 处），但均未出现在 Impact Analysis 的任何 Tier 中。

### [blindspot] `quality_gate_test.go` 中 `e2e-test` 相关测试用例的全面遗漏

测试文件中的测试用例引用了 `"e2e-test"` step（如第 539、572、725、947、1313 行），但 Impact Analysis Tier 4 的 `quality_gate_test.go` 条目仅写了 `HasRecipe(dir, "test") -> HasRecipe(dir, "unit-test")`，未涵盖 `e2e-test` 相关测试的重命名。

### [blindspot] 缺少 rollback plan

42 个文件的批量重命名，如果中途发现设计问题需要回退，rollback 策略是什么？提案未提及。v3.0.0 的"无兼容层"决策使得 partial rollback 极其困难——如果已经迁移了 Go 代码但未迁移模板，系统将处于不一致状态。

### [blindspot] `test` recipe 的可选参数签名未在 justfile template contract 中定义

`ExecuteJourneyInIsolation` 传递 positional argument `journeyName` 给 `just e2e-test`。迁移后 `just test <journeyName>` 需要生成的 justfile 中 `test` recipe 接受可选参数。但提案的 Impact Analysis Tier 3 仅写了"增加 test/test-setup"，未定义 `test` recipe 的参数签名。

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 98 | 110 |
| 2. Solution Clarity | 108 | 120 |
| 3. Industry Benchmarking | 88 | 120 |
| 4. Requirements Completeness | 89 | 110 |
| 5. Solution Creativity | 57 | 100 |
| 6. Feasibility | 88 | 100 |
| 7. Scope Definition | 72 | 80 |
| 8. Risk Assessment | 68 | 90 |
| 9. Success Criteria | 62 | 80 |
| 10. Logical Consistency | 74 | 90 |
| **Total** | **804** | **1000** |

## ATTACK_POINTS

1. **[Industry Benchmarking]** 稻草人替代方案 — "统一 + 缓存 + 并行 | — | 最大性能提升 | 复杂度过高，缓存机制引入新 bug 面 | Rejected: 过度设计" — 此替代方案无任何实质内容，仅为了衬托选中方案。必须替换为有意义的替代方案，例如 Bazel test size class model 或 GitHub Actions matrix gate。

2. **[Industry Benchmarking]** 行业方案引用深度不足 — 仅引用 Test Pyramid 概念，未引用具体 CI/CD 系统如何实现分层门禁 — 必须补充至少 2 个具体产品/项目的分层测试门禁实现参考。

3. **[Risk Assessment]** Impact Analysis 自身已证明遗漏概率为 High — "改动面大（~42 files）导致遗漏" 标为 M/M，但 `internal/cmd/test/test.go`、`internal/cmd/config_test.go`、`internal/cmd/config_schema_test.go`、`handleGateFailure` guide map 均未在 Impact Analysis 中列出 — 必须重新评估遗漏概率为 H/M。

4. **[Logical Consistency]** "无 Fallback" 声明与代码行为矛盾 — "无 Fallback：v3.0.0 直接要求 unit-test recipe，不回落到 test" vs `RunProjectTests` 的五级 fallback chain — 必须明确"无 Fallback"仅适用于 gate sequence 的 RunGate，不适用于 RunProjectTests 探测链。

5. **[Logical Consistency]** DefaultGateSequence 与 UnitGateSequence 职责不清 — "DefaultGateSequence() 的 test -> unit-test；新增 UnitGateSequence()" — 必须定义迁移后每个 sequence 函数的精确步骤内容和调用者。

6. **[Risk Assessment]** Mitigation 中参数风格与代码不一致 — "just test 需支持 --feature 参数" 但实际 `ExecuteJourneyInIsolation` 使用 positional argument — 必须统一 mitigation 描述与代码实际行为。

7. **[Requirements Completeness]** `quality_gate.go` 中 `runE2ERegression` 约 50 行代码的迁移复杂度被低估 — Impact Analysis 仅写了 "Step 3 用 test" — 必须将整个 `runE2ERegression` 函数的重构列为独立任务。

8. **[Solution Clarity]** `just test` 语义倒置的 UX 影响未在提案中承认 — 迁移后 `just test` 从"跑单元测试"变为"跑高级测试"，用户肌肉记忆失效 — 必须在 Risk Assessment 或 NFR 中明确承认此 UX 影响。

9. **[blindspot]** `handleGateFailure` 的 guide/label map 中 `"e2e-test"` 键需要迁移为 `"test"`，但未被 Impact Analysis 任何 Tier 覆盖 — 必须补充到 Tier 1。

10. **[blindspot]** 缺少 rollback plan — 42 文件批量重命名无回退策略 — 必须补充 rollback 方案（至少是"按 Tier 分批提交，每批可独立 revert"的策略）。

11. **[blindspot]** `test` recipe 的可选参数签名未在 justfile template contract 中定义 — `ExecuteJourneyInIsolation` 传递 positional argument 但提案未定义 `test` recipe 的参数规范 — 必须在 Requirements 或 Impact Analysis 中明确 `test` recipe 接受可选 positional parameter。

12. **[Success Criteria]** "耗时 <30s（Go 项目）" 缺少测量条件 — 未指定是否排除编译缓存、首次 vs 稳态 — 必须补充测量条件（如"冷启动后第二次运行"或"排除编译缓存后"）。
