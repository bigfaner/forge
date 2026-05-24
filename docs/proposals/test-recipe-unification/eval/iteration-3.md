---
iteration: 3
scorer: CTO adversary
model: proposal (post iteration-2 revision) vs baseline-snapshot
date: 2026-05-24
---

# Iteration 3 Score Report

## Iteration-2 Attack Points Resolution Audit

| # | Attack Point | Status | Evidence |
|---|-------------|--------|----------|
| 1 | Rollback Strategy batch 执行顺序与编译依赖矛盾 | **Resolved** | Rollback Strategy 表格已修正：Batch 1 现为 Tier 1 + Tier 4 合并（"源代码与测试同步修改"），消除了测试引用未定义符号的编译失败问题。line 321-323 |
| 2 | `just test` 语义倒置 UX 影响经三轮评审仍未承认 | **Resolved** | Key Risks 表格新增首行："`just test` 语义倒置...用户肌肉记忆失效 | H | M | 接受此变更——`unit-test`/`test` 命名对齐业界惯例..."，承认了 UX 影响并提供了接受理由和缓解措施（init-justfile 注释）。line 292 |
| 3 | Out of Scope 与 Tier 5 对"历史 proposals"矛盾 | **Resolved** | Tier 5 中"proposals/"条目已移除。Out of Scope 保留了"历史 lessons/proposals 中的 e2eTest 引用（不影响功能）"，Tier 5 仅保留"lessons/ (~9 files)"。line 222-227 vs line 278 |
| 4 | `handleGateFailure` guide/label map 在 Scope 中列出但无对应 Success Criteria | **Resolved** | Success Criteria 新增第 12 条："`handleGateFailure` guide/label map 中无 `"e2e-test"` 键，已全部迁移为 `"test"`"。line 312 |
| 5 | Prompt 模板和 Skill Markdown 更新无 Success Criteria | **Resolved** | Success Criteria 新增第 13 条："Go 源码、prompt 模板、skill markdown 中无残留 `e2e-test`/`e2e-setup`/`e2e-verify` 引用（`grep -r 'e2e-test\|e2e-setup\|e2e-verify' --include='*.go' --include='*.md'` 返回空）"。全局 grep 验证覆盖了所有文件类型。line 313 |
| 6 | "改动面大导致遗漏"仍标为 M/M | **Resolved** | Risk 表中"改动面大（~42 files）导致遗漏"已升级为 H/M。line 297 |
| 7 | 时间线估算仍缺失 | **Resolved** | Resource & Timeline 段落新增："预估 10-12 个 coding task，按 1 人全职执行估算 1.5–2 周（Batch 1 约 3 天，Batch 2–3 各约 1 天，Batch 4 约 1 天，含验证和修复时间）"。line 139 |
| 8 | 为什么两层而非三层的论证缺失 | **Resolved** | Industry Solutions 段落新增"两层 vs 三层决策分析"子段，从三个角度论证：(a) 门禁时机只有两个、(b) Surface-agnostic 原则、(c) 复杂度边界。line 102-106 |
| 9 | Rollback Strategy 中"每批提交后运行 go test 验证"与 batch 执行顺序矛盾 | **Resolved** | Batch 1 现为 Tier 1 + Tier 4 合并，且显式说明"源代码与测试同步修改，确保每步可编译通过"。line 323 |
| 10 | `parseAutoRaw()` 迁移检测逻辑的移除条件未定义 | **Resolved** | NFR 段落补充移除条件："移除条件：`parseAutoRaw()` 中的旧键名检测逻辑在 v3.1.0 中移除——即 v3.0.x 为过渡期，v3.1.0 起遇到旧键名 `e2eTest` 将直接报错而非映射"。line 79 |

**Resolution Summary**: 10/10 attack points resolved. Iteration-2 revision addressed every single issue raised in iteration-2. This is notable — iteration-2 revision achieved a clean sweep of all outstanding issues.

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

**Problem -> Solution**: 三层错位 -> 两层 recipe 模型。论证链成立且经过三轮 revision 持续强化。两条测试调用路径的区分消除了最大的逻辑矛盾。

**Solution -> Evidence**: 行业参考经过两轮升级，从泛泛而谈到 Bazel + GitHub Actions 具体产品引用，再到"两层 vs 三层决策分析"子段。论证从"引用原则"升级为"分析适用性"。

**Evidence -> Success Criteria**: 13 条 Success Criteria 经过两轮 revision 从 8 条扩展到 13 条，覆盖了所有 iteration-1 和 iteration-2 指出的 deliverable 缺口。parseAutoRaw 移除条件已定义。

**Self-contradiction check (post iteration-2 revision)**:
1. ~~Rollback Strategy batch 执行顺序~~ — 已解决（Batch 1 = Tier 1 + Tier 4 合并）
2. ~~`just test` 语义倒置未承认~~ — 已解决（Risk 表首行）
3. ~~Out of Scope vs Tier 5 矛盾~~ — 已解决（Tier 5 移除 proposals 条目）
4. ~~handleGateFailure 无 Success Criteria~~ — 已解决（新增第 12 条）
5. ~~Prompt/Skill 无 Success Criteria~~ — 已解决（新增第 13 条全局 grep）
6. ~~"改动面大" M/M~~ — 已解决（升级为 H/M）
7. ~~时间线缺失~~ — 已解决（1.5-2 周估算 + 分批时间）
8. ~~两层 vs 三层论证~~ — 已解决（新增决策分析子段）
9. ~~parseAutoRaw 移除条件~~ — 已解决（v3.1.0 移除）

**New observations post iteration-2 revision**:
- Rollback Strategy Batch 2（Tier 3 Justfile Templates）引用了 `test` recipe，但 Batch 2 中的 `templates/*.just` 文件可能引用 `unit-test` recipe（如在 shell 脚本中调用 `just unit-test`）。这不是编译依赖问题（justfile 是文本文件），但 Batch 2 提交后如果用户只运行 `just test` 而不运行 `just unit-test`，可能遇到模板逻辑不一致。这是一个极低严重性的观察。
- Success Criteria 第 13 条使用了 `grep -r 'e2e-test\|e2e-setup\|e2e-verify'`，但未包含 `e2eTest`（驼峰命名）的 grep。而 Success Criteria 第 8 条仅覆盖 `run-e2e-tests` 任务 key。配置文件中的 `e2eTest` YAML 键已被第 7 条（`auto.e2eTest` 完全移除）覆盖，但 grep 模式未统一。这是一个覆盖缺口。
- Resource & Timeline 段落给出了时间估算但标注了 Batch 编号（Batch 1 约 3 天，Batch 2–3 各约 1 天，Batch 4 约 1 天），然而 Batch 编号在 Rollback Strategy 表格中为 1/2/3/4，Resource & Timeline 中的 Batch 引用与 Rollback Strategy 的 Batch 定义一致。无矛盾。

---

## Phase 2: Rubric Scoring with Verification Stance

### 1. Problem Definition (110 pts)

**Problem stated clearly (40/40)**: 三层错位描述清晰，经过三轮 revision 保持稳定。Evidence 中的 5 条证据从行为层面到代码层面层层递进。iteration-2 revision 未修改问题定义，结构完整。

**Evidence provided (39/40)**: 5 条证据每条引用具体代码。~90s 量化数据具体。扣 1 分因为"10 个 breaking 任务累计 15 分钟"仍出现在 Urgency 段落而非 Evidence 段落——这是 iteration-2 已指出但未修改的结构性问题。

> "Breaking 任务 submit 运行 `just test`（即 `go test -race ./...`），单次 ~90s"

Urgency 段落的量化汇总应与 Evidence 段落交叉引用或合并。

**Urgency justified (29/30)**: 量化了延迟成本，v3.0.0 窗口期论据合理。扣 1 分因为"避免后续返工"仍缺乏具体性——iteration-2 指出此问题但未修改。

> "v3.0.0 分支正在进行测试能力重构，现在对齐可避免后续返工。"

"后续返工"指什么具体工作？新增功能会基于旧命名吗？文档会继续扩散 `e2eTest` 吗？

**Subtotal: 108/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete (40/40)**: 两层 recipe 模型完整。两条路径段落清晰。Gate sequence 表格三个 sequence 精确定义。Recipe 参数签名约定表格完整。经过三轮 revision，方案描述已达到可实施水平。

**User-facing behavior described (44/45)**: iteration-2 revision 在 Risk 表中承认了 `just test` 语义倒置的 UX 影响，提供了接受理由（对齐业界惯例）和缓解措施（init-justfile 注释）。这是三轮评审中首次回应此问题。扣 1 分因为：

- Risk 表中承认了 UX 影响，但 Solution 正文和 NFR 段落未提及。UX 影响仅存在于 Risk 表中，如果读者跳过 Risk 表，不会意识到这一语义变更。建议在 Solution 段落中也简要提及命名变更的影响。

**Technical direction clear (35/35)**: Gate sequence 表格、addFixTask 通用规则、runE2ERegression 迁移要点、Rollback Strategy batch 合并（Tier 1 + Tier 4）共同提供了充分且一致的技术方向。Rollback Strategy 的执行顺序逻辑问题已通过 batch 合并解决。

**Subtotal: 119/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (38/40)**: Bazel test size classes 和 GitHub Actions job gate 两个具体产品引用质量高。iteration-2 revision 新增的"两层 vs 三层决策分析"子段深入分析了 Forge 与 Bazel/GitHub Actions 的适用性差异。扣 2 分因为：

- 缺少对同类本地构建工具（如 Taskfile、Make-based test targets）的分层测试实践引用。所有引用要么是构建系统（Bazel）要么是 CI 平台（GitHub Actions），但 Forge 定位是本地优先的 CLI 工具，缺少同类工具的参考。

> "Bazel 将测试标注为 small（<1min）、medium（<5min）、large（无限制）"
> "GitHub Actions 通过 jobs.<id>.needs 声明 job 依赖"

两者都不是 Forge 的直接同类。

**At least 3 meaningful alternatives (28/30)**: 5 个选项包括 do nothing。Bazel 和 GitHub Actions 替代了 iteration-0 的稻草人。扣 2 分因为：

- "仅加 `test-quick` recipe" 仍是无 Source 的内部方案。iteration-1 和 iteration-2 均指出此问题，未修改。不过作为"最小改动"的对照方案，它有存在的价值——展示完全不做重命名时的问题。
- Bazel 替代方案的 Cons 列为"需测试框架配合 size 标注，Forge 无构建系统"——iteration-2 指出这接近稻草人但有一定论证价值（说明为什么不选）。可接受。

**Honest trade-off comparison (24/25)**: iteration-2 revision 新增的"两层 vs 三层决策分析"显著改善了选中方案的论证。三个角度（门禁时机、Surface-agnostic、复杂度边界）提供了诚实的 trade-off 分析。`just test` 语义倒置已在 Risk 表中承认。扣 1 分因为：

- 选中方案 Cons 列仍为"需更新多个组件；~42 文件迁移风险"，未包含 UX 影响（语义倒置）。虽然 Risk 表已承认，但 Comparison Table 的 Cons 应反映完整的 trade-off。

**Chosen approach justified against benchmarks (24/25)**: "两层 vs 三层决策分析"子段直接回答了 iteration-2 提出的"为什么两层"问题。映射到 submit vs all-completed 时机清晰。扣 1 分因为：

> "本方案将上述原则映射到 Forge 的 submit（仅 unit-test）vs all-completed（unit-test + test + probe）时机。"

这句话在 Industry Solutions 段落末尾，是将分析结论与方案选择关联的关键陈述。但它没有显式说明"因此两层是最优选择"——读者需要自行推导。

**Subtotal: 114/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage (38/40)**: 四个主场景（Breaking submit、非 Breaking submit、all-completed、journey isolation）覆盖完整。Recipe 参数签名约定和 Gate Sequence 无 Fallback 段落明确了边界条件。扣 2 分因为：

- 无 justfile 的项目（纯 go mod）submit 时的行为仍未定义。UnitGateSequence 要求 `unit-test` recipe，但如果项目没有 justfile（例如纯 Go 库），`HasRecipe` 会返回 false。提案声明"无 fallback——如果 recipe 不存在，gate 报错并提示运行 `init-justfile`"，但对于纯 go mod 项目，`init-justfile` 可能不适用。

**Non-functional requirements (39/40)**: 性能量化目标清晰。NFR 措辞"无持久兼容层"与 parseAutoRaw() 迁移检测逻辑一致。iteration-2 revision 补充了移除条件（v3.1.0 移除）。扣 1 分因为：

- 迁移提示的精确格式仍未完全定义。Success Criteria 说"输出迁移提示到 stderr"，Constraints 说提示内容为 `"config key 'auto.e2eTest' is renamed to 'auto.test' in v3.0.0; please update your config.yaml"`。但提示是否影响退出码（exit code）、是否在日志中记录、是否支持 i18n 都未提及。对于 Forge 这样的 CLI 工具，stderr 输出的格式规范是重要的 NFR。

**Constraints & dependencies (29/30)**: init-justfile 模板、journey_isolation.go、forgeconfig.Config、parseAutoRaw() 约束均明确。iteration-2 revision 补充了 parseAutoRaw() 迁移提示的具体内容。扣 1 分因为：

- Rollback Strategy 中"按反向顺序 revert（4→3→2→1）"暗示 Batch 4（Documentation）可以独立于 Batch 1（Go Source）revert。但文档中引用了新的 recipe 名（如 `unit-test`、`test`），如果 Batch 4 revert 到旧命名而 Batch 1 未 revert，文档将与实际代码不一致。这不是编译问题但影响文档准确性。

**Subtotal: 106/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (22/40)**: Proposal 声明"非创新性方案"。Surface-agnostic 抽象有一定价值但本质是复杂度下推到 recipe 层。经过三轮 revision 未改变此维度。这是诚实声明，不应扣分但不应给高分。

**Cross-domain inspiration (15/35)**: 未引用跨域灵感。Bazel 和 GitHub Actions 属于同一领域（构建/CI 系统的分层测试）。两层 vs 三层决策分析是良好的工程判断但非跨域创新。

**Simplicity of insight (24/25)**: 两层 recipe 模型将三层错位归约为简洁方案。两条路径的区分使简单更诚实。Rollback Strategy batch 合并（Tier 1 + Tier 4）解决了 iteration-2 的复杂度问题。扣 1 分因为 parseAutoRaw() 迁移检测逻辑引入了比表面看起来更多的复杂度（需要维护旧键名检测、输出格式、v3.1.0 移除时间线），但这些复杂度已被妥善管理。

**Subtotal: 61/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility (40/40)**: Rollback Strategy batch 合并解决了 iteration-2 指出的最大可行性问题。源代码与测试同步修改确保每步可编译通过。Go 代码改动、模板替换、配置重命名都是确定性工作。runE2ERegression 迁移要点提供了函数级重构方案。Full score。

**Resource & timeline feasibility (28/30)**: iteration-2 revision 补充了时间估算："按 1 人全职执行估算 1.5–2 周（Batch 1 约 3 天，Batch 2–3 各约 1 天，Batch 4 约 1 天，含验证和修复时间）"。估算合理且分批估算增加了可信度。扣 2 分因为：

> "涉及 ~42 个文件，影响范围评估如下。预估 10-12 个 coding task，按 1 人全职执行估算 1.5–2 周"

估算基于"1 人全职"，但未说明是否考虑了：(a) review 时间，（b）集成测试运行时间（Go 测试 + justfile 集成测试），（c）文档更新的实际工作量。Tier 5 列出了 ~15 文件（4 CLI docs + 1 ARCHITECTURE + 1 quality-gate + 1 testing/go + 1 forge-distribution + ~9 lessons），但 Batch 4 仅分配了 1 天。9 个 lessons 文件即使只是文本替换也需要时间。

**Dependency readiness (29/30)**: 无外部依赖。扣 1 分因为 Forge 自身 justfile 也需要修改，但 iteration-2 revision 未说明 Forge 自身的 justfile 修改是否包含在 Batch 2（Justfile Templates）中——如果 Forge 自己的 justfile 在 Batch 1 之前就需要 `unit-test` recipe 来运行测试验证，存在鸡生蛋问题。

**Subtotal: 97/100**

---

### 7. Scope Definition (80 pts)

**In-scope items are concrete (29/30)**: 每个文件级改动都有具体变更内容。Tier 1-5 分层清晰。iteration-2 revision 解决了 Out of Scope 与 Tier 5 的矛盾。扣 1 分因为 Impact Analysis Tier 1 和 Scope In Scope Go 代码条目对 `quality_gate.go` 的描述仍然完全重复——iteration-2 指出但未修改。这是维护性风险而非功能性错误。

**Out-of-scope explicitly listed (25/25)**: 9 条排除项清晰。iteration-2 revision 从 Tier 5 移除了 proposals 条目，消除了与 Out of Scope 的矛盾。Full score。

**Scope is bounded (23/25)**: 10-12 个 coding task，分 4 个 Batch，时间估算 1.5-2 周。iteration-2 revision 的时间补充使范围边界更清晰。扣 2 分因为：

- 时间估算仍基于"1 人全职"的假设，如果团队资源受限，范围边界会模糊。
- Tier 5 中"lessons/ (~9 files)"标为 Low priority，但 Low priority 意味着"可做可不做"，这与 In Scope 的语义冲突——如果它真的在 scope 内，应该有明确的完成标准。

**Subtotal: 77/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified (29/30)**: 6 个风险，包括 iteration-2 新增的 `just test` 语义倒置。风险覆盖了命名变更、迁移、测试失败、遗漏。扣 1 分因为：

- 缺少"用户忽略 parseAutoRaw() 迁移提示"的风险。NFR 声明"一次性迁移辅助"并定义了 v3.1.0 移除时间线，但如果用户持续忽略 stderr 提示（stderr 输出常被用户忽略），v3.1.0 移除检测逻辑后这些用户将直接报错，体验比 v3.0.x 更差。此风险未在 Risk 表中列出。

**Likelihood + impact rated (28/30)**: iteration-2 revision 将"改动面大导致遗漏"升级为 H/M，更诚实。`just test` 语义倒置标为 H/M 合理。扣 2 分因为：

- "auto.e2eTest → auto.test"标为 H/L（line 294），但 parseAutoRaw() 检测逻辑说明影响比 L 大。iteration-2 指出此问题但 Impact 评级未调整。如果用户不更新 config.yaml：v3.0.x 输出提示但仍工作，v3.1.0 直接报错。这意味着 H/L 中的 L 低估了长期影响（v3.1.0 后变为 H）。

**Mitigations are actionable (28/30)**: 语义倒置的缓解措施（业界惯例对齐 + init-justfile 注释）是可操作的。Rollback Strategy batch 合并解决了中间态不可测试问题。全局 grep 确认无残留。扣 2 分因为：

- "v3.0.0 要求重新运行 init-justfile 生成新 justfile"（line 293）仍缺少主动通知机制。用户如何知道需要重新运行 init-justfile？是在 submit 报错时看到提示？还是有主动的 upgrade guide？Risk 表未说明。
- parseAutoRaw() 的迁移提示输出到 stderr，但 stderr 在许多 CI 环境中被管道/重定向忽略，用户可能看不到。

**Subtotal: 85/90**

---

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (52/55)**: 13 条标准中 12 条可直接验证。iteration-2 revision 新增的第 12 条（handleGateFailure）和第 13 条（全局 grep）补全了覆盖缺口。第 13 条提供了具体的 grep 命令。扣 3 分因为：

- 第 13 条 grep 模式为 `'e2e-test\|e2e-setup\|e2e-verify'`，但不包含 `e2eTest`（驼峰命名）。Go 源码中的 `E2eTest` 字段名、`run-e2e-tests` 任务 key 等使用驼峰/kebab 命名，但 grep 模式只覆盖了 kebab 命名（`e2e-test`）。第 7 条和第 8 条分别覆盖了 `auto.e2eTest` 和 `run-e2e-tests`，但第 13 条的全局 grep 声称覆盖"Go 源码、prompt 模板、skill markdown"，如果只 grep `e2e-test` 模式，会遗漏 `e2eTest` 相关引用。

> "Go 源码、prompt 模板、skill markdown 中无残留 `e2e-test`/`e2e-setup`/`e2e-verify` 引用（`grep -r 'e2e-test\|e2e-setup\|e2e-verify' --include='*.go' --include='*.md'` 返回空）"

grep 模式与声明覆盖范围不完全匹配。

- "所有 Go 测试通过"（line 310）是基线要求，任何 PR 都需满足。作为提案成功标准价值有限——不会有人提交测试不过的 PR。iteration-2 指出但未修改。
- "parseAutoRaw() 检测到旧键名 e2eTest 时输出迁移提示到 stderr"（line 302）——验证需要捕获 stderr 的测试用例，但提案未要求编写此测试。

**Coverage is complete (24/25)**: 13 条 Success Criteria 经过两轮 revision 从 8 条扩展到 13 条。handleGateFailure、prompt/skill 文件更新都有对应标准。全局 grep 覆盖了所有文件类型。扣 1 分因为：

- Success Criteria 未覆盖 Forge 自身 justfile 的更新。Scope In Scope 中列出了"项目根 `justfile`"（line 265），但 Success Criteria 中无验证项确认 Forge 自身 justfile 已包含 `unit-test` recipe。

**Subtotal: 76/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (35/35)**: 两层模型解决了三层错位。两条路径区分消除了最大逻辑矛盾。NFR 与 Constraints 一致。Rollback Strategy batch 合并解决了执行顺序问题。Full score。

**Scope <-> Solution <-> Success Criteria aligned (28/30)**: iteration-2 revision 显著改善了对齐。handleGateFailure 有 Scope + Impact Analysis + Success Criteria。Prompt/skill 文件有 Scope + Success Criteria（全局 grep）。Out of Scope 与 Tier 5 不再矛盾。扣 2 分因为：

- Forge 自身 justfile 在 Scope 中列出（line 265: "项目根 `justfile`"）但 Success Criteria 中无对应验证项。
- Tier 5 中"lessons/ (~9 files)"标为 Low priority，In Scope 中也列出了"文档（6+ files）"包括"lessons/ (~9 files)"。Low priority 暗示"可跳过"，但 In Scope 暗示"必须做"。两者对同一 deliverable 的完成要求不一致。

**Requirements <-> Solution coherent (24/25)**: Recipe 参数签名约定、Gate Sequence 无 Fallback、两条路径区分、parseAutoRaw() 移除条件共同使 Requirements 和 Solution 高度一致。扣 1 分因为：

- Requirements 中"Gate Sequence 无 Fallback"声明适用于 RunGate 路径，而 RunProjectTests 保留 fallback。但 Requirements 中未说明为什么两条路径的行为差异是合理的。iteration-2 revision 在 Solution 正文中解释了两条路径的差异，但 Requirements 段落中缺少对这一设计决策的引用。

**Subtotal: 87/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] Success Criteria 第 13 条 grep 模式不完整

> "Go 源码、prompt 模板、skill markdown 中无残留 `e2e-test`/`e2e-setup`/`e2e-verify` 引用（`grep -r 'e2e-test\|e2e-setup\|e2e-verify' --include='*.go' --include='*.md'` 返回空）"

grep 模式仅覆盖 kebab-case（`e2e-test`、`e2e-setup`、`e2e-verify`），但遗漏了 camelCase 形式的 `e2eTest`。Go 源码中可能存在 `E2eTest` 字段名（已在第 7 条覆盖），但也可能存在注释、日志字符串、error message 中的 `e2eTest` 引用。如果第 13 条声称"Go 源码中无残留"但 grep 模式不包含 `e2eTest`，此验证标准存在覆盖缺口。

建议修改 grep 命令为 `grep -r 'e2e-test\|e2e-setup\|e2e-verify\|e2eTest' --include='*.go' --include='*.md'`。

**Severity**: Medium. 验证标准存在覆盖缺口，可能导致 `e2eTest` 残留通过验证。

### [blindspot] "用户忽略 parseAutoRaw() 迁移提示"风险未评估

> NFR: "无持久兼容层：v3.0.0 大版本重构，不引入永久兼容逻辑。仅提供一次性迁移辅助（`parseAutoRaw()` 检测旧键名输出提示并映射到新字段），迁移完成后移除。"
> Constraints: "`parseAutoRaw()` 需检测旧键名 `e2eTest`，输出迁移提示：`"config key 'auto.e2eTest' is renamed to 'auto.test' in v3.0.0; please update your config.yaml"`"

迁移提示输出到 stderr。在以下场景中用户可能忽略此提示：(a) CI 环境中 stderr 被重定向到 /dev/null，(b) 用户不习惯查看 stderr，(c) 提示信息淹没在其他输出中。v3.1.0 移除检测逻辑后，这些用户将遇到直接报错而非迁移提示，从"有警告但仍可用"变为"直接失败"。这一升级路径的 UX 退化未在 Risk 表中评估。

**Severity**: Medium. v3.1.0 移除检测逻辑的后果（直接报错 vs 迁移提示）未纳入风险评估。

### [blindspot] Forge 自身 justfile 的鸡生蛋问题

> Rollback Strategy: "每批提交后运行 `go test -race ./...` 验证"（line 328）
> Scope: "项目根 `justfile`" 在 Tier 3 Justfile Templates 中（line 265）

Forge 自身的 justfile 被归入 Batch 2（Tier 3 Justfile Templates）。但 Batch 1（Tier 1 + Tier 4 Go Source + Tests）提交后需要运行 `go test -race ./...` 验证。如果 Forge 的集成测试中有任何测试调用 `just test`（如 `tests/justfile-integration/` 中的测试），Batch 1 提交后这些测试将调用旧的 `test` recipe（此时语义已从"单元测试"变为"高级测试"？不对——justfile 还未更新）。

等等，Batch 1 只改 Go 源码和测试代码，不改 justfile。所以 Batch 1 提交后 Forge 的 justfile 仍是旧的——`test` recipe 仍跑旧逻辑。但 Go 测试中的断言已改为调用 `unit-test`，而 justfile 中尚无 `unit-test` recipe。这意味着 Batch 1 提交后运行 `go test` 时，如果有测试调用 `HasRecipe(dir, "unit-test")` 且测试 fixture（testdata 目录中的 justfile）已更新为包含 `unit-test` recipe，那么测试可以通过。但如果测试依赖 Forge 自身的 justfile（而非 testdata fixture），则会失败。

Rollback Strategy 声明"源代码与测试同步修改，确保每步可编译通过"——这确保了编译通过，但不确保测试通过（因为测试可能依赖 justfile 中尚不存在的 recipe）。提案未区分"编译通过"和"测试通过"。

**Severity**: Low-Medium. 取决于 Forge 的测试是否依赖 Forge 自身 justfile（而非 testdata fixture）。如果所有测试使用 testdata fixture，则无问题。

---

## Bias Detection Report

**Annotated regions (pre-revised)**: Markers at lines 36-43, 58-74, 88-89, 112-121, 148-160, 228-241, 292. Approximately 14 annotated paragraphs.

**Unannotated regions**: All remaining paragraphs (~32 paragraphs).

**Attack point analysis by region:**

Annotated regions attacks:
1. Success Criteria grep 模式不覆盖 camelCase (line 313, pre-revised: medium)
2. Risk 表 parseAutoRaw 迁移提示 stderr 可见性 (line 292, pre-revised: high)

Unannotated regions attacks:
1. Forge 自身 justfile 鸡生蛋问题 (line 321-328)
2. Tier 5 "Low priority" 与 In Scope 语义冲突 (line 222-227)
3. "所有 Go 测试通过"基线标准 (line 310)
4. Impact Analysis Tier 1 与 Scope 重复描述 (line 162 vs line 243)
5. Success Criteria 缺少 Forge 自身 justfile 验证 (line 265)
6. "Urgency"量化数据在 Urgency 段而非 Evidence 段 (line 23)
7. "auto.e2eTest → auto.test" H/L 低估长期影响 (line 294)

- Annotated regions: 2 attack points / 14 annotated paragraphs = density 0.14
- Unannotated regions: 7 attack points / ~32 unannotated paragraphs = density 0.22
- Ratio (annotated/unannotated): 0.64

**Bias assessment**: 当前比率为 0.64（annotated 攻击密度低于 unannotated），说明 pre-revised 区域的质量略高于 unrevised 区域。这与 iteration-2 的 1.00 比率相比有所下降，原因是 pre-revised 区域的 revision 质量很高（解决了所有 iteration-2 问题），而 unrevised 区域中一些长期存在的小问题（Urgency 结构、基线标准、Tier 5 语义冲突）仍未修改。总体而言，pre-revision 对质量的提升是实质性的。

---

## Score Summary

| Dimension | Score | Max | Delta vs Iteration-2 |
|-----------|-------|-----|---------------------|
| 1. Problem Definition | 108 | 110 | +1 |
| 2. Solution Clarity | 119 | 120 | +2 |
| 3. Industry Benchmarking | 114 | 120 | +10 |
| 4. Requirements Completeness | 106 | 110 | +4 |
| 5. Solution Creativity | 61 | 100 | +1 |
| 6. Feasibility | 97 | 100 | +6 |
| 7. Scope Definition | 77 | 80 | +3 |
| 8. Risk Assessment | 85 | 90 | +8 |
| 9. Success Criteria | 76 | 80 | +4 |
| 10. Logical Consistency | 87 | 90 | +3 |
| **Total** | **930** | **1000** | **+42** |

## ATTACK_POINTS

1. **[Success Criteria]** grep 模式不覆盖 camelCase 命名 — "Go 源码、prompt 模板、skill markdown 中无残留 `e2e-test`/`e2e-setup`/`e2e-verify` 引用（`grep -r 'e2e-test\|e2e-setup\|e2e-verify'`）" — grep 模式遗漏 `e2eTest`（camelCase），可能导致 Go 源码中注释、日志、error message 的 `e2eTest` 残留通过验证。必须扩展 grep 模式为 `'e2e-test\|e2e-setup\|e2e-verify\|e2eTest'`。

2. **[Risk Assessment]** "用户忽略 parseAutoRaw() 迁移提示"风险未评估 — NFR 声明 v3.1.0 移除检测逻辑，但 stderr 输出在 CI 环境中常被忽略。v3.1.0 移除后用户从"有警告但仍可用"升级为"直接报错"。必须新增 Risk 条目评估此升级路径的 UX 退化。

3. **[blindspot]** Forge 自身 justfile 鸡生蛋问题 — Rollback Strategy "每批提交后运行 `go test -race ./...` 验证"但 Batch 1 只改 Go 源码不改 justfile，如果集成测试依赖 Forge 自身 justfile 中的 `unit-test` recipe（而非 testdata fixture），Batch 1 后测试将失败。必须在 Rollback Strategy 中区分"编译通过"和"测试通过"，并说明 Forge 自身 justfile 在 Batch 1 中的处理策略。

4. **[Scope Definition]** Tier 5 "Low priority" 与 In Scope 语义冲突 — "lessons/ (~9 files)" 在 Tier 5 标为 "Low（历史文档，不强制更新）" 但同时列在 In Scope 的"文档（6+ files）"中。Low priority 暗示"可跳过"，In Scope 暗示"必须做"。必须统一：要么从 In Scope 移除并归入 Out of Scope，要么在 In Scope 中明确标注"可选"。

5. **[Feasibility]** 时间线未考虑 review 和集成测试运行时间 — "按 1 人全职执行估算 1.5–2 周" 但未说明是否包含 code review 时间、集成测试运行时间（justfile 集成测试可能耗时较长）、文档更新工作量。Batch 4 分配 1 天但 Tier 5 列出 ~15 文件。必须补充估算假设或增加缓冲。

6. **[Problem Definition]** "Urgency"中量化汇总应在 Evidence 中交叉引用 — "10 个 breaking 任务累计 15 分钟纯等待" 出现在 Urgency 段落但不在 Evidence 段落。如果读者仅阅读 Evidence 验证问题严重性，会遗漏此量化依据。必须在 Evidence 中增加交叉引用或将汇总数据移入 Evidence。

7. **[Logical Consistency]** Forge 自身 justfile 在 Scope 中但无 Success Criteria — Scope 列出"项目根 `justfile`"（line 265）但 13 条 Success Criteria 中无验证项确认 Forge 自身 justfile 已包含 `unit-test` recipe。必须新增对应验证标准。

8. **[Risk Assessment]** "auto.e2eTest → auto.test" H/L 低估长期影响 — "H | L" 但 v3.1.0 移除 parseAutoRaw() 检测后影响升级为"直接报错"（从 L 变为 M-H）。必须将 Impact 评级改为 "L（v3.0.x）/ M（v3.1.0+）" 或在 Mitigation 中说明 v3.1.0 前的主动通知策略。
