---
iteration: 1
scorer: CTO adversary
model: proposal (post pre-revision) vs baseline-snapshot
date: 2026-05-24
---

# Iteration 1 Score Report

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

**Problem -> Solution**: 三层错位 -> 两层 recipe 模型。论证链成立且清晰。Pre-revision 未改变问题定义。

**Solution -> Evidence**: 方案引用了具体代码文件。Pre-revision 新增了"两条测试调用路径"段落（line 37-43），明确区分 RunGate 和 RunProjectTests 的独立行为。这是对 iteration-0 中最大矛盾点的直接回应。新增的 gate sequence 表格（line 113-121）明确了三个 sequence 的精确内容。Evidence 链得到显著加强。

**Evidence -> Success Criteria**: 成功标准未在 pre-revision 中修改，但新增的 Requirements 段落（Recipe 参数签名约定、Gate Sequence 无 Fallback 说明）使部分成功标准的含义更明确。

**Self-contradiction check**:
1. ~~"无 Fallback" vs RunProjectTests 的 fallback chain~~ -- 已解决：新增段落（line 73-74）明确区分两条路径
2. ~~DefaultGateSequence 与 UnitGateSequence 职责不清~~ -- 已解决：新增表格（line 113-121）+ 重命名为 FullGateSequence
3. ~~Impact Analysis 遗漏 internal/cmd/test/test.go~~ -- 已解决：新增到 Tier 1（line 160）
4. ~~addFixTask 映射问题~~ -- 已解决：新增到 Impact Analysis（line 154）

### Residual Contradictions (Post Pre-Revision)

1. Risk 表中 `journey_isolation.go` 的 mitigation 仍写 "just test 需支持 --feature 参数"（line 281），但 Recipe 参数签名约定段落（line 62-70）已明确定义为 positional argument `journey=''`。Mitigation 与 Requirements 段落不一致。
2. Non-Functional Requirements 仍声明"无向后兼容负担"（line 79），但 Constraints 段落（line 89）已引入 `parseAutoRaw()` 检测旧键名的迁移提示逻辑。NFR 的"直接重命名，不保留旧键名/旧 recipe 的兼容逻辑"与 Constraints 中"检测到旧键名时，将其值映射到新字段"语义冲突——映射旧字段就是一种兼容逻辑。

---

## Phase 2: Rubric Scoring with Verification Stance

### 1. Problem Definition (110 pts)

**Problem stated clearly (38/40)**: 三层错位的描述清晰。Pre-revision 未改变问题定义，保持稳定。扣 2 分因为"三层错位"中的第二层（模板硬编码 Playwright）和第三层（配置命名过时）在 Evidence 中的位置暗示它们是独立问题，但实际上它们都是同一次重构遗漏的不同表现——问题本质是"justfile recipe 模型未同步适配去 Profile 化重构"的一体多面。这一结构可让读者误解为三个需要分别解决的问题。

> "submit 门禁跑全量测试浪费时间、模板仍硬编码 Playwright、配置项命名过时"

三个并列描述暗示独立性，但实质是同一根因的不同症状。

**Evidence provided (38/40)**: 5 条证据，每条引用具体代码。~90s 量化数据具体。扣 2 分因为 Evidence 中"10 个 breaking 任务累计 15 分钟纯等待"出现在 Urgency 段落而非 Evidence 段落，但实际被当作 Evidence 使用——这是结构性问题。

**Urgency justified (28/30)**: 量化了延迟成本（90s/次），v3.0.0 窗口期论据合理。扣 2 分因为"避免后续返工"仍然模糊——具体什么后续工作？如果 v3.0.0 已经在做测试重构，为什么不在这轮重构中一并完成此方案，而是在"现在做"和"以后做"之间选择了"现在"？

> "v3.0.0 分支正在进行测试能力重构，现在对齐可避免后续返工。"

**Subtotal: 104/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete (39/40)**: 两层 recipe 模型 + 新增的两条路径段落 + gate sequence 表格。Pre-revision 显著改善了具体性。扣 1 分因为 Innovation Highlights 段落（line 45-47）在 pre-revision 后变成了非连续段落——它紧跟"两条测试调用路径"段落之后，但标题层级（###）暗示它是 Proposed Solution 的平级子节，读起来像是把 Innovation Highlights 从"核心方案"移到了"补充说明"的位置，逻辑流不够顺畅。

**User-facing behavior described (42/45)**: submit 门禁和 all-completed 行为清晰。Pre-revision 新增了 Recipe 参数签名约定表格（line 59-69），用户可观测行为更加明确。扣 3 分因为：
- 用户手动运行 `just test` 时从"单元测试"变为"高级测试"这一语义倒置的 UX 影响仍未在提案中承认（iteration-0 已指出，pre-revision 未回应）
- `parseAutoRaw()` 检测旧键名时的迁移提示体验未描述——用户看到什么格式的提示？stderr? stdout? 退出码是否改变？

> "输出迁移提示：`"config key 'auto.e2eTest' is renamed to 'auto.test' in v3.0.0; please update your config.yaml"`"

提示内容定义了，但输出渠道和运行时行为未说明。

**Technical direction clear (34/35)**: gate sequence 表格和 addFixTask 通用规则让技术方向明确。Pre-revision 补充了关键细节。扣 1 分因为 `quality_gate.go` 中 `runE2ERegression` 函数的完整迁移方案仍然过于简化——Impact Analysis 中仅一行"Step 3 用 test"，但该函数约 50 行代码涉及多个 recipe 调用（e2e-setup、dev、e2e-test），完整重构需要更多技术说明。

**Subtotal: 115/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (28/40)**: 仅引用 Test Pyramid 概念。Pre-revision 未改善此维度。缺少具体产品/项目的分层测试门禁实现参考（GitHub Actions job gate、GitLab CI stage gate、Bazel test size classes）。扣 12 分。

> "分层测试（Test Pyramid）是业界标准：单元测试（快速、大量）→ 集成测试（中速、适量）→ E2E 测试（慢、少量）。CI/CD 通常在 commit 时跑单元测试，merge/merge request 时跑集成+E2E。"

泛泛而谈 CI/CD 的通用做法，没有引用任何具体实现或产品。

**At least 3 meaningful alternatives (22/30)**: 4 个选项包括 do nothing。但"统一 + 缓存 + 并行"仍然是稻草人——Pre-revision 未替换此选项。无实际方案细节，无 Source 引用，仅为了衬托选中方案。

> "统一 + 缓存 + 并行 | — | 最大性能提升 | 复杂度过高，缓存机制引入新 bug 面 | Rejected: 过度设计"

Source 列为 "—"，无实质内容。

**Honest trade-off comparison (19/25)**: Pre-revision 未改善。选中方案的 Cons 仍仅写"需更新多个组件"。

> "结构清晰；Forge surface-agnostic；向后兼容 | 需更新多个组件 | Selected"

Pros 列出"向后兼容"，但 NFR 明确说"无向后兼容负担"。自相矛盾。真正的 trade-off（语义倒置的 UX 影响、~42 files 的迁移风险）未被承认。

**Chosen approach justified against benchmarks (20/25)**: "直击痛点，复杂度可控"仍是结论性陈述。Pre-revision 未改善。为什么映射到两层而非三层？为什么不是 unit-test + integration-test + e2e-test？扣 5 分。

**Subtotal: 89/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage (36/40)**: Pre-revision 新增了 Recipe 参数签名约定段落和 Gate Sequence 无 Fallback 段落，场景覆盖改善。扣 4 分因为仍缺少：
- 无 justfile 的项目（纯 go mod）submit 时的行为——UnitGateSequence 要求 `unit-test` recipe，但无 justfile 的项目无法满足
- `unit-test` recipe 执行超时的处理

**Non-functional requirements (34/40)**: 性能量化目标清晰。Pre-revision 在 Constraints 中新增了迁移提示逻辑，改善了迁移 UX。扣 6 分因为：
- NFR 中"无向后兼容负担"与 Constraints 中 `parseAutoRaw()` 迁移提示逻辑直接矛盾——如果做了旧键名检测和映射，那就是一种兼容逻辑
- 缺少迁移后错误提示的 UX 规范（输出格式、渠道）

> "无向后兼容负担：v3.0.0 大版本重构，直接重命名，不保留旧键名/旧 recipe 的兼容逻辑"

vs

> "parseAutoRaw() 需检测旧键名 e2eTest，输出迁移提示...检测到旧键名时，将其值映射到新字段而非静默忽略"

后者本质上就是向后兼容层。

**Constraints & dependencies (27/30)**: Pre-revision 新增了 `parseAutoRaw()` 约束和 `internal/cmd/test/test.go` 约束。扣 3 分因为 `quality_gate.go` 中 `runE2ERegression` 整个函数（约 50 行）的迁移复杂度仍被低估。

**Subtotal: 97/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (22/40)**: Proposal 声明"非创新性方案"。surface-agnostic 抽象有一定价值但本质是复杂度下推。Pre-revision 未改变此维度。扣 18 分。

**Cross-domain inspiration (15/35)**: 未引用跨域灵感。Pre-revision 未改善。扣 20 分。

**Simplicity of insight (23/25)**: 将三层错位归约为两层 recipe，方案简洁。Pre-revision 通过新增两条路径的明确区分，使"简单"更诚实。扣 2 分因为引入 `parseAutoRaw()` 迁移检测增加了方案的隐含复杂度——一个声称"直接重命名"的方案实际包含了旧键名检测和映射逻辑。

**Subtotal: 60/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility (37/40)**: Go 代码改动和模板替换都是确定性工作。Pre-revision 新增了 gate sequence 精确定义和 addFixTask 通用规则，可行性论证更充分。扣 3 分因为 `runE2ERegression` 函数的完整重构复杂度仍被低估。

**Resource & timeline feasibility (24/30)**: ~42 文件 / 10-12 个 coding task 估算合理。Pre-revision 将 Tier 1 Go 文件从 8 增加到 9（新增 `test.go`），估算更准确。扣 6 分因为仍无时间线估算。

> "涉及 ~42 个文件，影响范围评估如下。预估 10-12 个 coding task。"

**Dependency readiness (28/30)**: 无外部依赖。扣 2 分原因同前。

**Subtotal: 89/100**

---

### 7. Scope Definition (80 pts)

**In-scope items are concrete (29/30)**: Pre-revision 补充了 `test.go`、`parseAutoRaw()` 迁移逻辑、addFixTask 通用规则。每个文件级改动都有具体变更内容。扣 1 分因为 `handleGateFailure` 中 guide/label map 的 `"e2e-test"` -> `"test"` 迁移仍未在任何 Tier 中列出（iteration-0 blindspot 已指出）。

**Out-of-scope explicitly listed (22/25)**: 9 条排除项。Pre-revision 未改变。扣 3 分因为"历史 lessons/proposals 中的 e2eTest 引用"被标为 Out of Scope 但这些文档中的过时引用可能在新用户阅读时造成混淆。

**Scope is bounded (22/25)**: 任务数量有上限（10-12），Tier 分批执行策略清晰。扣 3 分因为仍无时间边界。

**Subtotal: 73/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified (25/30)**: 5 个风险。Pre-revision 未修改 Risk 表。扣 5 分因为遗漏：
- `handleGateFailure` guide/label map 中 `"e2e-test"` 键的迁移（iteration-0 blindspot 已指出）
- `just test` 语义倒置对用户肌肉记忆的影响（iteration-0 已指出）
- NFR 声称"无向后兼容"但实际引入了 `parseAutoRaw()` 迁移检测——方案内部的自相矛盾本身是风险

**Likelihood + impact rated (23/30)**: 评估大体诚实。扣 7 分因为：
- "改动面大（~42 files）导致遗漏"标为 M/M（line 283），但 iteration-0 已证明 `handleGateFailure` 的遗漏，实际遗漏概率为 H
- "auto.e2eTest -> auto.test" 标为 H/L（line 280），但提案新增了迁移提示逻辑本身就说明影响比原先评估的要大

> "改动面大（~42 files）导致遗漏 | M | M"

iteration-0 发现了至少 3 处遗漏，说明 M 评级不准确。

**Mitigations are actionable (24/30)**: Pre-revision 未修改 Risk 表中的 mitigations。扣 6 分因为：
- "just test 需支持 --feature 参数，模板统一生成"（line 281）仍使用 `--feature` flag 描述，但 Recipe 参数签名约定段落已定义使用 positional argument `journey=''`。Mitigation 与 Requirements 段落不一致。

> "just test 需支持 --feature 参数，模板统一生成"

> "justfile 模板中的 test recipe 必须接受可选的第一个参数 journey（just 语法：test journey=''）"

同一提案内对同一技术点的描述自相矛盾。

**Subtotal: 72/90**

---

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (44/55)**: 8 条标准中 6 条可验证。Pre-revision 未修改 Success Criteria。扣 11 分因为：
- "耗时 <30s（Go 项目）"（line 287）未指定测量条件（冷启动？排除编译缓存？首次 vs 稳态？）
- "完整覆盖，无 e2e-test 调用"（line 288）需要全代码库 grep 验证但未提供验证命令
- "所有 Go 测试通过"（line 293）是基线要求而非本提案的成功标准
- 缺少 `parseAutoRaw()` 迁移提示的成功标准——旧键名检测和映射的行为应可验证

**Coverage is complete (21/25)**: 覆盖了主要 in-scope 项。Pre-revision 新增了多个 in-scope 项（parseAutoRaw 检测、addFixTask 通用规则、test.go），但 Success Criteria 未相应扩展。扣 4 分因为缺少：
- `parseAutoRaw()` 检测旧键名并输出迁移提示的验证标准
- `addFixTask` 中 `step -> "just " + step` 通用规则的验证标准
- Prompt 模板和 skill markdown 更新的验证标准

**Subtotal: 65/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (32/35)**: 两层模型解决了三层错位。Pre-revision 新增的两条路径段落消除了 iteration-0 最大的逻辑矛盾。扣 3 分因为：

NFR 声称"无向后兼容负担"但 Constraints 引入了迁移提示逻辑。这不是"负担"的程度问题，而是语义矛盾——映射旧字段值到新字段就是兼容逻辑。

> "无向后兼容负担：v3.0.0 大版本重构，直接重命名，不保留旧键名/旧 recipe 的兼容逻辑"

> "检测到旧键名时，将其值映射到新字段而非静默忽略"

**Scope <-> Solution <-> Success Criteria aligned (26/30)**: Pre-revision 显著改善了 Solution 和 Scope 的对齐（新增 parseAutoRaw、addFixTask、test.go）。扣 4 分因为：
- Scope 新增了 `parseAutoRaw()` 迁移提示逻辑，但 Success Criteria 未包含对应验证项
- Scope 新增了 `addFixTask` 通用规则，但 Success Criteria 未包含对应验证项
- Risk 表中 mitigation 的 `--feature` 与 Solution 中 positional argument 矛盾

**Requirements <-> Solution coherent (22/25)**: Pre-revision 新增了 Recipe 参数签名约定和 Gate Sequence 无 Fallback 段落，改善了 Requirements 和 Solution 的对齐。扣 3 分因为：
- Requirements 中"无 Fallback"声明（line 74）在 pre-revision 后精确限定了"仅适用于 gate sequence 的 RunGate，不适用于 RunProjectTests"，但这一限定与 Risk 表中"无 unit-test recipe 时 quality gate 报错"的成功标准仍存在微妙的执行路径问题——如果 RunProjectTests 被其他场景调用（如 CLI `forge run-tests`），而项目没有 `unit-test` recipe，fallback 到 `test`（高级测试）用于"运行项目测试"的语义是否正确？Requirements 未定义此场景。

**Subtotal: 80/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] Risk 表 mitigation 中 `--feature` flag 与 Recipe 参数签名约定的 positional argument 矛盾

> Risk 表: "just test 需支持 --feature 参数，模板统一生成"
> Recipe 参数签名约定: "just test [journey]...test journey=''"

Pre-revision 新增了精确的参数签名定义但未同步更新 Risk 表。这是 pre-revision 引入的新不一致。

**Tag: conflict-with-pre-revision** — Pre-revision 在 Requirements 中修正了参数签名但未传播到 Risk 表。

### [blindspot] NFR "无向后兼容负担" 与 Constraints 中 `parseAutoRaw()` 迁移检测的语义矛盾

> NFR: "v3.0.0 大版本重构，直接重命名，不保留旧键名/旧 recipe 的兼容逻辑"
> Constraints: "parseAutoRaw() 需检测旧键名 e2eTest，输出迁移提示...检测到旧键名时，将其值映射到新字段而非静默忽略"

`parseAutoRaw()` 的旧键名检测和映射本质上是一个受控的向后兼容层。NFR 的绝对化声明与实际方案不一致。Pre-revision 引入了 Constraints 中的迁移逻辑但未同步修改 NFR 的措辞。

**Tag: conflict-with-pre-revision** — Pre-revision 在 Constraints 中新增了迁移逻辑但未更新 NFR 以反映这一设计变化。

### [blindspot] Success Criteria 未覆盖 pre-revision 新增的 in-scope 项

Pre-revision 在 Scope 中新增了至少 3 个具体 deliverable：
1. `parseAutoRaw()` 检测旧键名输出迁移提示
2. `addFixTask` 移除硬编码映射，改为 `step -> "just " + step` 通用规则
3. `internal/cmd/test/test.go` 的 `e2e-test` -> `test` 迁移

但 8 条 Success Criteria 均未在 pre-revision 中更新以覆盖这些新增项。新增 deliverable 没有 可验证的成功标准。

**Tag: conflict-with-pre-revision** — Pre-revision 扩展了 Scope 但未同步扩展 Success Criteria。

### [blindspot] `handleGateFailure` guide/label map 遗漏仍存在

Iteration-0 blindspot 已指出 `handleGateFailure` 中 `map[string]string` 包含 `"e2e-test"` 键需迁移为 `"test"`。Pre-revision 未回应此问题。Impact Analysis 的任何 Tier 仍不包含此改动。

### [blindspot] `runE2ERegression` 函数的完整迁移方案仍未充分说明

Iteration-0 blindspot 已指出此函数约 50 行代码涉及 `e2e-setup`、`dev`、`e2e-test` 多个 recipe 调用。Pre-revision 新增了 Step 3 的描述但仅限一行，函数级别的重构方案仍不充分。

### [blindspot] 缺少 rollback plan

42 个文件的批量重命名，如果中途发现设计问题需要回退，rollback 策略仍缺失。Pre-revision 未回应。v3.0.0 的"无兼容层"决策使得 partial rollback 极其困难。

### [blindspot] Comparison Table 中"向后兼容"Pros 与 NFR 矛盾

> Comparison Table: "结构清晰；Forge surface-agnostic；向后兼容"
> NFR: "无向后兼容负担：v3.0.0 大版本重构，直接重命名，不保留旧键名/旧 recipe 的兼容逻辑"

选中方案的 Pros 写"向后兼容"，但 NFR 明确说"无向后兼容负担"。同一个提案对"向后兼容"给出了矛盾的判断。

---

## Bias Detection Report

Annotated regions (pre-revised): 13 markers spanning paragraphs at lines 36-43, 58-74, 88-89, 112-121, 148-160, 228-241.
Unannotated regions: All remaining paragraphs.

**Attack point analysis by region:**

Annotated regions attacks:
1. `parseAutoRaw()` 约束引入但 NFR 未同步更新 (line 88-89)
2. Gate sequence 表格内容清晰但 `runE2ERegression` 完整迁移仍缺失 (line 112-121)
3. `addFixTask` 通用规则描述但 `handleGateFailure` 遗漏 (line 153-154)
4. `test.go` 新增到 Scope 但 Success Criteria 未覆盖 (line 159-160, 240-241)
5. NFR vs Constraints 语义矛盾 (line 79 vs line 88-89)
6. Risk 表 `--feature` vs Recipe 参数签名 positional argument (line 281 vs line 62-70)
7. Scope 新增项无 Success Criteria 对应 (line 228-241)

Unannotated regions attacks:
1. Comparison Table "向后兼容" vs NFR "无向后兼容负担" (line 104)
2. Industry benchmarking 深度不足 (line 91-95)
3. Straw-man alternative "统一 + 缓存 + 并行" (line 103)
4. `runE2ERegression` 迁移复杂度低估 (line 110)
5. Success Criteria "耗时 <30s" 缺测量条件 (line 287)
6. Success Criteria 未覆盖 pre-revision 新增项 (line 286-295)
7. Missing rollback plan (entire document)
8. `handleGateFailure` guide/label map 遗漏 (unaddressed from iteration-0)

- Annotated regions: 7 attack points / 13 annotated paragraphs = density 0.54
- Unannotated regions: 8 attack points / ~35 unannotated paragraphs = density 0.23
- Ratio (annotated/unannotated): 2.35

**Bias assessment**: Annotated regions have significantly higher attack density (2.35x). This is expected — pre-revised regions contain new content that introduces new inconsistencies with existing unrevised content (NFR, Risk table, Comparison Table). The attacks on annotated regions are predominantly `conflict-with-pre-revision` tagged, indicating the revisions improved their local context but created cross-section misalignment by not propagating changes to dependent sections.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 104 | 110 |
| 2. Solution Clarity | 115 | 120 |
| 3. Industry Benchmarking | 89 | 120 |
| 4. Requirements Completeness | 97 | 110 |
| 5. Solution Creativity | 60 | 100 |
| 6. Feasibility | 89 | 100 |
| 7. Scope Definition | 73 | 80 |
| 8. Risk Assessment | 72 | 90 |
| 9. Success Criteria | 65 | 80 |
| 10. Logical Consistency | 80 | 90 |
| **Total** | **844** | **1000** |

## ATTACK_POINTS

1. **[Industry Benchmarking]** 稻草人替代方案 — "统一 + 缓存 + 并行 | — | 最大性能提升 | 复杂度过高，缓存机制引入新 bug 面 | Rejected: 过度设计" — Source 列为 "—"，无任何实质内容或行业引用。必须替换为有意义的替代方案（如 Bazel test size class model 或 GitHub Actions matrix gate）。

2. **[Industry Benchmarking]** 行业方案引用深度不足 — "分层测试（Test Pyramid）是业界标准"仅引用概念，未引用任何具体产品/项目的分层门禁实现。必须补充至少 2 个具体实现参考。

3. **[Industry Benchmarking]** Comparison Table Pros 与 NFR 矛盾 — "结构清晰；Forge surface-agnostic；**向后兼容**" vs NFR "无向后兼容负担"。同一个提案对向后兼容性给出了矛盾的判断。必须统一表述。

4. **[Risk Assessment]** Risk 表 mitigation 与 Recipe 参数签名约定不一致 — Risk 表写"just test 需支持 **--feature** 参数"但 Recipe 参数签名约定定义使用 **positional argument** `journey=''`。Pre-revision 修正了后者但未传播到前者。必须将 Risk 表中的 `--feature` 改为 positional argument 描述。

5. **[Requirements Completeness]** NFR "无向后兼容负担" 与 Constraints 中 `parseAutoRaw()` 迁移检测矛盾 — NFR 说"不保留旧键名/旧 recipe 的兼容逻辑"但 Constraints 新增了"检测到旧键名时，将其值映射到新字段"。映射旧字段值就是兼容逻辑。Pre-revision 引入了约束但未同步更新 NFR。必须修改 NFR 措辞以反映实际设计。

6. **[Success Criteria]** 新增 in-scope deliverable 缺少可验证标准 — Pre-revision 在 Scope 中新增了 `parseAutoRaw()` 迁移提示、`addFixTask` 通用规则、`test.go` 迁移 3 个 deliverable，但 8 条 Success Criteria 未做任何扩展。每个 in-scope deliverable 必须有对应的可验证成功标准。

7. **[Success Criteria]** "耗时 <30s（Go 项目）" 缺少测量条件 — 未指定冷启动 vs 稳态、是否排除编译缓存。必须补充测量条件。

8. **[Logical Consistency]** `conflict-with-pre-revision`: NFR 声称绝对化"无兼容"但实际引入了受控兼容层 — 两处描述自相矛盾。必须选择一种表述并全文统一。

9. **[Scope Definition]** `handleGateFailure` guide/label map 遗漏仍存在 — iteration-0 blindspot 指出 `map[string]string{"e2e-test": "fix failing e2e tests"}` 需迁移为 `"test"`，但 Impact Analysis 任何 Tier 仍未包含。必须补充到 Tier 1。

10. **[blindspot]** 缺少 rollback plan — 42 文件批量重命名无回退策略。v3.0.0 的"无兼容层"决策使 partial rollback 极其困难。必须补充 rollback 方案（至少"按 Tier 分批提交，每批可独立 revert"）。

11. **[blindspot]** `runE2ERegression` 函数约 50 行代码的完整重构方案缺失 — Impact Analysis 中仅一行"Step 3 用 test"，但该函数涉及 `e2e-setup`、`dev`、`e2e-test` 多个 recipe 调用。必须列出函数级别的迁移要点。

12. **[blindspot]** `conflict-with-pre-revision`: Pre-revision 新增的内容与未修订内容产生交叉不一致 — Risk 表 `--feature`（未修订）vs Recipe 参数签名 `journey=''`（已修订）、NFR "无兼容"（未修订）vs Constraints 迁移映射（已修订）、Comparison Table "向后兼容"（未修订）vs NFR（未修订但与 Constraints 矛盾）。Pre-revision 改善了局部质量但引入了跨段落不一致。建议在下一轮 revision 中重点处理 pre-revision 内容向下游段落（Risk、NFR、Comparison Table、Success Criteria）的传播。
