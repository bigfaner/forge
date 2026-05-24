---
iteration: 2
scorer: CTO adversary
model: proposal (post iteration-1 revision) vs baseline-snapshot
date: 2026-05-24
---

# Iteration 2 Score Report

## Iteration-1 Attack Points Resolution Audit

| # | Attack Point | Status | Evidence |
|---|-------------|--------|----------|
| 1 | 稻草人替代方案"统一 + 缓存 + 并行" | **Resolved** | Comparison Table 替换为 Bazel test size classes 和 GitHub Actions job gate，Source 列有具体产品名 |
| 2 | 行业方案引用深度不足 | **Resolved** | 新增 Bazel（small/medium/large size classes）和 GitHub Actions（job dependency gate）两个具体实现参考 |
| 3 | Comparison Table Pros "向后兼容" 与 NFR 矛盾 | **Resolved** | Pros 改为"无旧接口兼容负担"，NFR 改为"无持久兼容层：...仅提供一次性迁移辅助"，表述一致 |
| 4 | Risk 表 `--feature` 与 Recipe 参数签名矛盾 | **Resolved** | Risk 表改为"positional argument `journey=''`"，与 Recipe 参数签名约定一致 |
| 5 | NFR "无向后兼容负担" 与 parseAutoRaw() 迁移检测矛盾 | **Resolved** | NFR 重写为"无持久兼容层...仅提供一次性迁移辅助"，承认了迁移逻辑的存在但限定为一次性 |
| 6 | 新增 in-scope deliverable 缺少可验证标准 | **Resolved** | 新增 3 条 Success Criteria 覆盖 parseAutoRaw、addFixTask、test.go |
| 7 | "耗时 <30s" 缺测量条件 | **Resolved** | 补充"排除编译缓存后的首次运行，稳态 CI 环境" |
| 8 | NFR 绝对化"无兼容"与受控兼容层矛盾 | **Resolved** | NFR 措辞已从绝对化声明改为限定量化的"无持久兼容层" |
| 9 | handleGateFailure guide/label map 遗漏 | **Resolved** | Tier 1 Impact Analysis 和 Scope 均列出 `handleGateFailure` 中 `"e2e-test"` → `"test"` 迁移 |
| 10 | 缺少 rollback plan | **Resolved** | 新增"Rollback Strategy"段，按 Tier 分批提交，每批可独立 revert，反向顺序回退 |
| 11 | runE2ERegression 函数迁移方案缺失 | **Resolved** | 新增"runE2ERegression 迁移要点"段，列出 5 条迁移要点 |
| 12 | Pre-revision 内容与未修订内容的交叉不一致 | **Partially Resolved** | 主要矛盾已修复（NFR、Comparison Table、Risk 表），但存在新的残留问题（见 Phase 3） |

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

**Problem -> Solution**: 三层错位 -> 两层 recipe 模型。论证链成立。"两条测试调用路径"段落准确区分了 RunGate 和 RunProjectTests 的独立行为，消除了 iteration-0 中最大的逻辑矛盾。

**Solution -> Evidence**: 行业参考从泛泛而谈升级为 Bazel 和 GitHub Actions 两个具体产品。Comparison Table 的 Cons 列从"需更新多个组件"升级为"需更新多个组件；~42 文件迁移风险"，更诚实。但行业参考仍偏浅——Bazel 的引用仅说明了"分阶段门禁"这一共同原则，未深入分析 Forge 与 Bazel 在构建系统层面的本质差异是否使该参考不完全适用。

**Evidence -> Success Criteria**: 11 条 Success Criteria 覆盖了 iteration-1 指出的所有新增 deliverable。parseAutoRaw、addFixTask、test.go 都有了对应标准。测量条件已补充。

**Self-contradiction check (post iteration-1 revision)**:
1. ~~Comparison Table "向后兼容" vs NFR "无向后兼容负担"~~ — 已解决
2. ~~Risk 表 `--feature` vs Recipe 参数签名 `journey=''`~~ — 已解决
3. ~~NFR "无兼容" vs parseAutoRaw 迁移映射~~ — 已解决（措辞改为"无持久兼容层"）
4. ~~handleGateFailure 遗漏~~ — 已解决
5. ~~Rollback plan 缺失~~ — 已解决
6. **新发现**: Impact Analysis Tier 1 中 `quality_gate.go` 描述（line 162）同时列出了 `addFixTask` 和 `handleGateFailure` 的改动，但 Scope In Scope 的 Go 代码条目（line 243）中同一条目也列出了相同内容，两处描述完全重复。这不是错误但造成了维护负担——如果未来修改其中一处但忘记另一处，将产生不一致。
7. **新发现**: Rollback Strategy 表中 Batch 1 为 Tier 4（Go Tests），但 Tier 4 依赖 Tier 1（Go Source）中定义的函数签名。如果先改测试文件中的断言（Batch 1）再改源代码（Batch 3），测试将引用不存在的函数名/字段名，导致中间状态不可通过测试。Rollback 策略的执行顺序存在逻辑问题。

### Residual Contradictions (Post Iteration-1 Revision)

1. **Rollback Strategy 执行顺序与 Tier 依赖关系矛盾**。Batch 1（Tier 4 Go Tests）的测试断言引用 Tier 1（Go Source）中定义的函数名和字段名。如果按 Batch 1 → 2 → 3 顺序执行，测试将引用 `FullGateSequence`、`UnitGateSequence`、`auto.Test` 等尚不存在的符号，`go test` 必然失败。提案要求"每批提交后运行 `go test -race ./...` 验证"（line 321），但 Batch 1 提交后测试无法通过。

2. **Success Criteria 第 10 条"所有 Go 测试通过"（line 304）与 Rollback Strategy 的分批验证循环依赖**。按 Rollback Strategy，每批提交后需运行全量测试验证。但如果 Batch 1（测试文件）先于 Batch 3（源代码）提交，测试必然失败——因为测试断言引用的符号尚未存在。这意味着要么：(a) Rollback Strategy 的执行顺序不正确，要么 (b) Success Criteria 的"所有 Go 测试通过"不适用于中间 batch。

---

## Phase 2: Rubric Scoring with Verification Stance

### 1. Problem Definition (110 pts)

**Problem stated clearly (39/40)**: 三层错位的描述清晰。iteration-1 revision 未改变问题定义，保持稳定。扣 1 分因为 Evidence 中 5 条证据的排列暗示它们是独立问题，但实际上它们都是"justfile recipe 模型未同步适配去 Profile 化重构"的不同症状——第一条（submit 门禁跑全量）是行为问题，第二到五条是代码层面的具体错位。读者可能误判为需要多个独立修复。

> "submit 门禁跑全量测试浪费时间、模板仍硬编码 Playwright、配置项命名过时"

并列列举暗示独立性，但本质是同一根因。

**Evidence provided (39/40)**: 5 条证据，每条引用具体代码。~90s 量化数据具体。iteration-1 revision 未修改 Evidence 段落。扣 1 分因为"10 个 breaking 任务累计 15 分钟"出现在 Urgency 段落而非 Evidence 段落，这是结构性问题——如果有人仅阅读 Evidence 部分来验证问题严重性，会遗漏这个量化依据。

**Urgency justified (29/30)**: 量化了延迟成本（90s/次），v3.0.0 窗口期论据合理。扣 1 分因为"避免后续返工"仍然不够具体——具体什么后续工作？是文档中已出现的 `e2eTest` 引用会继续扩散？还是新增功能会继续基于旧命名？

> "v3.0.0 分支正在进行测试能力重构，现在对齐可避免后续返工。"

**Subtotal: 107/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete (40/40)**: 两层 recipe 模型完整。两条路径段落清晰区分 RunGate 和 RunProjectTests。Gate sequence 表格明确了三个 sequence 的精确内容和调用方。Recipe 参数签名约定表格定义了每个 recipe 的签名。Full score。

**User-facing behavior described (43/45)**: submit 门禁和 all-completed 行为清晰。Recipe 参数签名约定表格明确了用户可观测的调用方式。parseAutoRaw() 迁移提示在 Success Criteria 中明确了输出到 stderr。扣 2 分因为：

- `just test` 从"单元测试"变为"高级测试"的语义倒置 UX 影响仍未在提案正文中承认。iteration-0 和 iteration-1 都指出了此问题，两轮 revision 均未回应。这是一个对现有用户的 UX 破坏性变更，提案应至少承认这一影响并说明为什么选择接受它。

**Technical direction clear (34/35)**: gate sequence 表格、addFixTask 通用规则、runE2ERegression 迁移要点共同提供了充分的技术方向。扣 1 分因为 Rollback Strategy 的执行顺序存在逻辑问题（Batch 1 测试引用 Batch 3 源代码符号），说明对实现顺序的思考不够深入。

**Subtotal: 117/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (35/40)**: iteration-1 revision 新增了 Bazel test size classes 和 GitHub Actions job gate 两个具体产品参考。引用质量从"泛泛而谈"升级为有实质内容。扣 5 分因为：

- Bazel 引用说明了"分阶段门禁"原则但未分析 Forge 与 Bazel 在构建系统层面的本质差异（Bazel 依赖构建图，Forge 依赖 justfile 约定）是否使该参考不完全适用
- GitHub Actions 引用是 CI 平台层面的模式，Forge 是本地优先的工具，两者的门禁执行环境差异未被讨论
- 缺少对同类 CLI 工具（如 Taskfile、Make-based）的分层测试实践引用

> "Bazel 将测试标注为 small（<1min）、medium（<5min）、large（无限制），CI 按大小分阶段门禁"

引用了原则但未分析适用性边界。

**At least 3 meaningful alternatives (26/30)**: 4 个选项包括 do nothing。iteration-1 revision 将稻草人"统一 + 缓存 + 并行"替换为 Bazel 和 GitHub Actions 两个真实行业方案。扣 4 分因为：

- Bazel test size classes 被作为替代方案呈现但其 Cons 写"需测试框架配合 size 标注，Forge 无构建系统"——这个 Cons 本质上是在说"这个方案不适用于 Forge"，这使它接近于稻草人。如果真要做适配分析，应讨论 Bazel 的 size class 概念是否可以在不引入 Bazel 构建系统的情况下被 Forge 借鉴（例如通过 justfile 注释或环境变量标注 size）。
- "仅加 test-quick recipe" 仍是无 Source 的内部方案，不是行业验证过的替代方案。

**Honest trade-off comparison (21/25)**: iteration-1 revision 改善了选中方案的 Cons——从"需更新多个组件"升级为"需更新多个组件；~42 文件迁移风险"。NFR 措辞改为"无持久兼容层"后与 Comparison Table 的 Pros "无旧接口兼容负担"一致。扣 4 分因为：

- 选中方案的真正 trade-off 仍未充分承认：(a) `just test` 语义倒置对用户肌肉记忆的影响，(b) 两层模型与三层测试金字塔的映射损失（集成测试被折叠进 `test`，无法独立门禁）
- Cons 列中"~42 文件迁移风险"是量化风险描述，但 Impact Analysis 中部分 Tier 的文件列表仍可能不完整（见 blindspot）

**Chosen approach justified against benchmarks (22/25)**: "直击痛点，复杂度可控"仍是主要理由，但 iteration-1 revision 补充了"本方案将上述原则映射到 Forge 的 submit vs all-completed 时机"，提供了更多论证。扣 3 分因为：

- 为什么映射到两层而非三层的问题仍未回答。Bazel 用三层（small/medium/large），GitHub Actions 用多层（lint → unit → integration → e2e），但 Forge 只用两层。提案应解释为什么三层不适用于 Forge。
- "直击痛点，复杂度可控"仍是结论性陈述，不是论证。

**Subtotal: 104/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage (37/40)**: iteration-1 revision 新增了 Recipe 参数签名约定和 Gate Sequence 无 Fallback 段落。覆盖了 Breaking 任务 submit、非 Breaking submit、all-completed、journey isolation 四个主场景。扣 3 分因为仍缺少：

- 无 justfile 的项目（纯 go mod）submit 时的行为——UnitGateSequence 要求 `unit-test` recipe，但无 justfile 的项目无法满足。如果 `HasRecipe` 探测到 justfile 不存在，gate 如何处理？
- `unit-test` recipe 执行超时的处理——提案提到了性能目标（<30s），但如果 `unit-test` recipe 挂起或超时，gate 的行为未定义。

**Non-functional requirements (37/40)**: 性能量化目标清晰。iteration-1 revision 将 NFR 重写为"无持久兼容层：...仅提供一次性迁移辅助"，消除了与 Constraints 段落 parseAutoRaw() 的矛盾。扣 3 分因为：

- 迁移完成后如何移除 parseAutoRaw() 中的检测逻辑？是随下一个 minor version 清理？还是需要在成功标准中增加"parseAutoRaw 无旧键名检测代码"？
- 缺少迁移提示的 UX 规范——Success Criteria 说输出到 stderr，但提示的精确格式、是否影响退出码、是否在日志中记录都未在 NFR 中定义。

**Constraints & dependencies (28/30)**: iteration-1 revision 新增了 parseAutoRaw() 约束和 internal/cmd/test/test.go 约束。扣 2 分因为：

- Rollback Strategy 的 Tier 执行顺序（测试先于源代码）与实际编译依赖之间存在未声明的约束——修改测试断言前必须先修改源代码中的函数/字段定义。

**Subtotal: 102/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (22/40)**: Proposal 声明"非创新性方案"（line 47）。surface-agnostic 抽象有一定价值但本质是复杂度下推到 recipe 层。iteration-1 revision 未改变此维度。扣 18 分。

**Cross-domain inspiration (15/35)**: 未引用跨域灵感。iteration-1 revision 新增了 Bazel 和 GitHub Actions 参考，但它们都属于同一领域（构建/CI 系统的分层测试），不是跨域。扣 20 分。

**Simplicity of insight (23/25)**: 将三层错位归约为两层 recipe，方案简洁。iteration-1 revision 通过两条路径的区分使"简单"更诚实。扣 2 分因为方案实际引入的隐含复杂度（parseAutoRaw 迁移检测、Rollback Strategy 的 batch 依赖管理）比表面看起来更多。

**Subtotal: 60/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility (38/40)**: Go 代码改动和模板替换都是确定性工作。Gate sequence 表格和 addFixTask 通用规则让可行性论证充分。runE2ERegression 迁移要点补全了函数级重构方案。扣 2 分因为 Rollback Strategy 的 batch 执行顺序存在技术可行性问题——Batch 1（测试）先于 Batch 3（源代码）提交会导致编译失败。

**Resource & timeline feasibility (24/30)**: ~42 文件 / 10-12 个 coding task 估算合理。扣 6 分因为仍无时间线估算。

> "涉及 ~42 个文件，影响范围评估如下。预估 10-12 个 coding task。"

**Dependency readiness (29/30)**: 无外部依赖。扣 1 分因为 Forge 自身 justfile 也需要修改，迁移过程中 Forge 自己的测试 recipe 可能暂时不可用。

**Subtotal: 91/100**

---

### 7. Scope Definition (80 pts)

**In-scope items are concrete (29/30)**: 每个文件级改动都有具体变更内容。handleGateFailure guide/label map 已纳入 Tier 1 和 Scope。扣 1 分因为 Impact Analysis Tier 1（line 162）和 Scope In Scope Go 代码条目（line 243）对 `quality_gate.go` 的描述完全重复——两处同时列出 `addFixTask` 和 `handleGateFailure` 的改动。如果未来修改一处但忘记另一处，将产生不一致。

**Out-of-scope explicitly listed (23/25)**: 9 条排除项。扣 2 分因为"历史 lessons/proposals 中的 e2eTest 引用"被标为 Out of Scope（line 279）但 Impact Analysis Tier 5 中又列出了"proposals/（历史提案中的 e2eTest 引用）"标为 Low priority。同一条目同时出现在 Out of Scope 和 Tier 5 中，产生歧义——它到底做不做？

> Out of Scope: "历史 lessons/proposals 中的 e2eTest 引用（不影响功能）"
> Tier 5: "proposals/（历史提案中的 e2eTest 引用）| ~5 | Low（历史文档，不强制更新）"

**Scope is bounded (22/25)**: 任务数量有上限（10-12），Tier 分批执行策略清晰。扣 3 分因为仍无时间边界，且 Rollback Strategy 的 batch 执行顺序未考虑编译依赖。

**Subtotal: 74/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified (27/30)**: 5 个风险。iteration-1 revision 修正了 Risk 表中 `--feature` 为 `positional argument journey=''`。扣 3 分因为仍遗漏：

- `just test` 语义倒置对用户肌肉记忆的影响（iteration-0 和 iteration-1 均指出，仍未回应）
- Rollback Strategy batch 执行顺序导致的中间态不可测试问题

**Likelihood + impact rated (24/30)**: 评估大体诚实。扣 6 分因为：

- "改动面大（~42 files）导致遗漏"仍标为 M/M（line 291），但 iteration-0 和 iteration-1 共发现了至少 5 处 Impact Analysis 遗漏（test.go、config_test.go、config_schema_test.go、handleGateFailure、quality_gate_test.go 中 e2e-test 用例）。虽然这些遗漏已在后续 revision 中修复，但它们证明了在 ~42 文件的大规模重命名中遗漏概率确实是 H。
- "auto.e2eTest → auto.test" 标为 H/L（line 288），但 parseAutoRaw() 迁移检测逻辑本身就说明影响比 L 大——如果用户不更新 config.yaml，Forge 会输出提示但仍然工作。这个风险的实际影响是"用户可能忽略提示继续使用旧配置"，长期后果是 parseAutoRaw() 中的检测逻辑永远不会被移除，与 NFR 的"一次性迁移辅助"承诺矛盾。

**Mitigations are actionable (26/30)**: iteration-1 revision 修正了 `--feature` 为 positional argument，消除了与 Recipe 参数签名的矛盾。Rollback Strategy 的补充使大部分风险有了可操作的缓解方案。扣 4 分因为：

- "v3.0.0 要求重新运行 init-justfile 生成新 justfile"（line 287）是 requirement 而非 mitigation——它没有说明如何确保用户知道要这样做。无主动通知机制。
- Rollback Strategy 中"每批提交后运行 go test -race ./... 验证"（line 321）与 batch 执行顺序矛盾——Batch 1 提交后测试无法通过。

**Subtotal: 77/90**

---

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (49/55)**: 11 条标准中 9 条可验证。iteration-1 revision 新增了 3 条标准覆盖 parseAutoRaw、addFixTask、test.go。测量条件已补充（"排除编译缓存后的首次运行，稳态 CI 环境"）。扣 6 分因为：

- "完整覆盖，无 e2e-test 调用"（line 299）需要全代码库 grep 验证但未提供验证命令（如 `grep -r 'e2e-test' --include='*.go'`）。验证方式隐含但不显式。
- "所有 Go 测试通过"（line 304）是基线要求而非本提案的独特成功标准——任何 PR 都需要这个条件。
- "parseAutoRaw() 检测到旧键名 e2eTest 时输出迁移提示到 stderr"（line 296）——如何验证"输出到 stderr"？需要测试用例捕获 stderr。提案未要求为此编写测试。

**Coverage is complete (23/25)**: iteration-1 revision 显著改善了覆盖——新增的 3 条 Success Criteria 覆盖了 iteration-1 指出的新增 deliverable。扣 2 分因为缺少：

- `handleGateFailure` guide/label map 中 `"e2e-test"` → `"test"` 迁移的验证标准——此改动已在 Scope 和 Impact Analysis 中列出但无对应 Success Criteria
- Prompt 模板和 skill markdown 更新的验证标准（iteration-1 指出但未修复）

**Subtotal: 72/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (34/35)**: 两层模型解决了三层错位。两条路径段落消除了 iteration-0 最大的逻辑矛盾。NFR 与 Constraints 的矛盾已通过措辞修改解决。扣 1 分因为 iteration-1 revision 新增的 Rollback Strategy 引入了新的逻辑问题——batch 执行顺序与编译依赖矛盾。

**Scope <-> Solution <-> Success Criteria aligned (27/30)**: iteration-1 revision 改善了对齐。handleGateFailure 已纳入 Scope，新增 3 条 Success Criteria 覆盖了关键 deliverable。扣 3 分因为：

- `handleGateFailure` guide/label map 在 Scope（line 243）中明确列出，但 Success Criteria 中无对应验证项
- Prompt 模板和 skill markdown 在 Scope 中列出（line 251-256）但 Success Criteria 中无对应验证项
- Out of Scope "历史 proposals 中 e2eTest 引用"与 Tier 5 "proposals/（历史提案中的 e2eTest 引用）"存在歧义

**Requirements <-> Solution coherent (23/25)**: iteration-1 revision 新增的 Recipe 参数签名约定和 Gate Sequence 无 Fallback 段落显著改善了 Requirements 和 Solution 的对齐。扣 2 分因为：

- Rollback Strategy 的 batch 执行顺序与 Requirements 中"每批提交后运行 go test 验证"的要求矛盾——如果 Batch 1 的测试断言引用 Batch 3 才定义的符号，测试无法通过
- Requirements 中"Gate Sequence 无 Fallback"的声明适用于 RunGate 路径，但 Rollback Strategy 中如果 Batch 3（源代码）被 revert 而 Batch 1（测试）未 revert，gate sequence 将引用不存在的函数，触发"无 Fallback"的报错——这不是真正的 fallback 问题但行为与预期不一致

**Subtotal: 84/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] Rollback Strategy batch 执行顺序与编译依赖矛盾

> "按 Tier 分批提交，每批可独立 revert"
> "| 1 | Tier 4（Go Tests） | ~11 | 独立 revert，不影响生产代码 |"
> "| 3 | Tier 1（Go Source） | ~9 | revert 后代码恢复旧函数名/字段名 |"

Batch 1 修改测试断言中的函数名（如 `HasRecipe(dir, "test")` → `HasRecipe(dir, "unit-test")`，`E2eTest` → `Test`），但这些函数/字段在 Batch 3 才被定义。按 Batch 1 → 2 → 3 顺序执行，Batch 1 提交后测试必然编译失败。提案要求"每批提交后运行 `go test -race ./...` 验证"（line 321）与执行顺序矛盾。

正确顺序应为 Batch 3（源代码）→ Batch 1（测试），或 Batch 3 和 Batch 1 合并为同一批次。

**Severity**: High. Rollback Strategy 的执行顺序如果被团队直接执行，将在 Batch 1 提交后立即遇到编译失败，导致策略不可行。

### [blindspot] `just test` 语义倒置 UX 影响仍未承认

> iteration-0 指出: "迁移后 just test 从'跑单元测试'变为'跑高级测试'，用户肌肉记忆失效"
> iteration-1 指出: "用户手动运行 just test 时从'单元测试'变为'高级测试'这一语义倒置的 UX 影响仍未在提案中承认"

三轮评审（iteration-0、iteration-1、iteration-2）均指出此问题。提案始终未在正文任何位置（Risk、NFR、Solution）承认这一 UX 影响。这不是遗漏——而是有意忽略。如果一个 UX 破坏性变更经过三轮评审都未被承认，说明提案作者认为这不是问题。但作为一个声称"用户体验"重要的提案，对三轮同一反馈的沉默本身就是问题。

### [blindspot] Out of Scope 与 Tier 5 对同一项目的重复/矛盾定义

> Out of Scope: "历史 lessons/proposals 中的 e2eTest 引用（不影响功能）"（line 279）
> Tier 5: "proposals/（历史提案中的 e2eTest 引用）| ~5 | Low（历史文档，不强制更新）"（line 222）

"历史 proposals 中 e2eTest 引用"同时出现在 Out of Scope 和 Tier 5 中。Out of Scope 意味着"不做"，Tier 5 意味着"应当更新"。两者矛盾。应二选一：要么从 Out of Scope 移除并在 Tier 5 中明确为"可选更新"，要么从 Tier 5 移除。

### [blindspot] `handleGateFailure` guide/label map 在 Scope 中列出但无 Success Criteria 对应

> Scope: "`handleGateFailure` guide/label map `"e2e-test"` → `"test"`"（line 243）
> Success Criteria: 无对应条目

iteration-1 revision 将 `handleGateFailure` 纳入了 Scope 和 Impact Analysis，但 Success Criteria 仍无对应验证项。这意味着即使 `handleGateFailure` 的迁移被遗漏，所有 11 条 Success Criteria 仍可通过。

### [blindspot] Prompt 模板和 Skill Markdown 更新无 Success Criteria 对应

> Scope: "Prompt 模板（3 files）"和"Skill/Command Markdown（5+ files）"（line 251-256）
> Success Criteria: 无对应条目

Scope 中列出了 ~8 个 prompt/skill 文件的更新，但 Success Criteria 中无任何针对这些文件的验证标准。如果更新遗漏了某个 prompt 文件中的 `e2e-test` 引用，Success Criteria 不会捕获。

---

## Bias Detection Report

**Annotated regions (pre-revised)**: Markers at lines 36-43, 58-74, 88-89, 112-121, 148-160, 228-241. Approximately 13 annotated paragraphs.

**Unannotated regions**: All remaining paragraphs (~35 paragraphs).

**Attack point analysis by region:**

Annotated regions attacks:
1. Rollback Strategy batch 顺序问题关联到 pre-revised gate sequence 定义 (line 121-128)
2. Comparison Table Bazel 引用适用性分析不足 (line 99-100, pre-revised context)
3. NFR "无持久兼容层" 中"迁移完成后移除"无时间线 (line 79)

Unannotated regions attacks:
1. `just test` 语义倒置 UX 影响未被承认 (entire document, unaddressed since iteration-0)
2. Rollback Strategy batch 执行顺序与编译依赖矛盾 (line 309-321)
3. Out of Scope 与 Tier 5 对 "历史 proposals" 的矛盾 (line 222 vs line 279)
4. `handleGateFailure` 在 Scope 但无 Success Criteria (line 243 vs line 293-305)
5. Prompt/Skill 文件更新无 Success Criteria (line 251-256 vs line 293-305)
6. "改动面大导致遗漏" M/M 评级不准确 (line 291)
7. 时间线估算仍缺失 (line 133)
8. "所有 Go 测试通过" 是基线要求非提案成功标准 (line 304)

- Annotated regions: 3 attack points / 13 annotated paragraphs = density 0.23
- Unannotated regions: 8 attack points / ~35 unannotated paragraphs = density 0.23
- Ratio (annotated/unannotated): 1.00

**Bias assessment**: 与 iteration-1 的 2.35 比率相比，当前比率为 1.00，表明 pre-revision 内容与 unrevised 内容的质量差异已基本消除。iteration-1 revision 成功解决了 pre-revision 引入的跨段落不一致（NFR、Comparison Table、Risk 表）。剩余的 attack points 主要集中在未修订区域（Rollback Strategy 是新增内容、Success Criteria 覆盖缺口、UX 影响承认缺失），说明 revision 的质量传播是有效的。

---

## Score Summary

| Dimension | Score | Max | Delta vs Iteration-1 |
|-----------|-------|-----|---------------------|
| 1. Problem Definition | 107 | 110 | +3 |
| 2. Solution Clarity | 117 | 120 | +2 |
| 3. Industry Benchmarking | 104 | 120 | +15 |
| 4. Requirements Completeness | 102 | 110 | +5 |
| 5. Solution Creativity | 60 | 100 | 0 |
| 6. Feasibility | 91 | 100 | +2 |
| 7. Scope Definition | 74 | 80 | +1 |
| 8. Risk Assessment | 77 | 90 | +5 |
| 9. Success Criteria | 72 | 80 | +7 |
| 10. Logical Consistency | 84 | 90 | +4 |
| **Total** | **888** | **1000** | **+44** |

## ATTACK_POINTS

1. **[Feasibility]** Rollback Strategy batch 执行顺序与编译依赖矛盾 — "Batch 1 | Tier 4（Go Tests）...Batch 3 | Tier 1（Go Source）" — Batch 1 修改测试断言引用 Batch 3 才定义的函数/字段，提交后 `go test` 必然编译失败。必须调整执行顺序为 Batch 3 → Batch 1，或合并为同一批次。

2. **[Solution Clarity]** `just test` 语义倒置 UX 影响经三轮评审仍未承认 — 迁移后 `just test` 从"单元测试"变为"高级测试"，现有用户的肌肉记忆将失效。iteration-0、iteration-1、iteration-2 均指出此问题，提案正文未在任何位置承认此 UX 影响。必须在 Risk Assessment 或 NFR 中明确承认此影响并解释为何接受。

3. **[Scope Definition]** Out of Scope 与 Tier 5 对"历史 proposals"矛盾 — Out of Scope: "历史 lessons/proposals 中的 e2eTest 引用（不影响功能）" vs Tier 5: "proposals/（历史提案中的 e2eTest 引用）| ~5 | Low" — 同一条目同时出现在"不做"和"应当更新"中。必须二选一。

4. **[Success Criteria]** `handleGateFailure` guide/label map 在 Scope 中列出但无对应 Success Criteria — Scope: "`handleGateFailure` guide/label map `"e2e-test"` → `"test"`" 但 11 条 Success Criteria 中无验证项。必须新增对应验证标准。

5. **[Success Criteria]** Prompt 模板和 Skill Markdown 更新无 Success Criteria — Scope 列出 ~8 个 prompt/skill 文件更新但 Success Criteria 无对应条目。如果更新遗漏某文件中的 `e2e-test` 引用，成功标准不会捕获。必须新增如"代码库中无残留 `e2e-test` 引用（Go 源码、prompt 模板、skill markdown）"的全局验证标准。

6. **[Risk Assessment]** "改动面大导致遗漏"仍标为 M/M — "改动面大（~42 files）导致遗漏 | M | M" — iteration-0 和 iteration-1 共发现 5 处 Impact Analysis 遗漏，证明 ~42 文件重命名中遗漏概率为 H。必须升级为 H/M。

7. **[Feasibility]** 时间线估算仍缺失 — "涉及 ~42 个文件...预估 10-12 个 coding task" — 有任务数量但无时间估算。10-12 个 task 是 1 天？1 周？2 周？团队无法评估是否可并行执行。必须补充时间估算。

8. **[Industry Benchmarking]** 为什么两层而非三层的论证缺失 — Bazel 用三层（small/medium/large），GitHub Actions 用多层（lint → unit → integration → e2e），但 Forge 只用两层。提案未解释为什么三层不适用于 Forge。必须在 Industry Benchmarking 或 Solution 中补充两层 vs 三层的决策分析。

9. **[blindspot]** Rollback Strategy 中"每批提交后运行 go test 验证"与 batch 执行顺序矛盾 — "每批提交后运行 `go test -race ./...` 验证"（line 321）但 Batch 1 提交后测试引用 Batch 3 的符号，无法通过。Rollback Strategy 整体不可行，需要重新设计 batch 分组和执行顺序。

10. **[blindspot]** `parseAutoRaw()` 迁移检测逻辑的移除条件未定义 — NFR 声明"仅提供一次性迁移辅助...迁移完成后移除"（line 79）但未定义"迁移完成"的判定条件。如果用户长期不更新 config.yaml，parseAutoRaw() 中的检测逻辑将永远不会被移除。必须定义移除条件（如"下个 minor version"或"检测逻辑在 v3.1 中移除"）。
