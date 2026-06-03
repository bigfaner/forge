---
reviewer: Modular Documentation Architect
date: "2026-06-03"
document: docs/proposals/skill-command-independence-audit/proposal.md
---

# Freeform Review: Skill & Command Independence Audit

## 1. Background Assessment

本提案目标是消除 Forge plugin 21 个 skill 和 16 个 command 中的文档耦合，实现"每个文件是独立知识单元"的原则。具体操作分为三维度：(1) 将跨 skill 内部文件引用的知识内联到引用方；(2) 精简冗余描述；(3) 删除 Related Skills / Integration / References 章节及相关变体。

提案正确识别了 6 个有跨 skill 引用的 skill（gen-journeys、gen-contracts、gen-test-scripts、extract-design-md、init-justfile）和 1 个有跨 skill 引用的 command（fix-bug）。将 forensic 的动态 SKILL.md 加载排除在 scope 之外是合理的——forensic 的核心功能就是运行时加载其他 skill 定义并与实际行为做对比，这是一种设计意图而非耦合问题。

提案选择了"内容内联+精简"而非"提升到共享层"，这与 Forge 的分发模型一致——Forge plugin 作为独立单元分发到用户环境，skill 之间不存在共享层的路径约定。

---

## 2. Key Risks

### R1. 耦合图不完整——遗漏了 gen-contracts 到 gen-journeys 的反向引用

**风险**：提案声明修复"gen-contracts: 内联 gen-jours Surface Detection 相关知识"，但实际代码中 gen-contracts 的 SKILL.md 第 58 行写着：

> `Detect the project's surface types via forge surfaces. See gen-journeys/SKILL.md "Surface Detection" section for the full detection flow, exit code contract, and detection flow steps.`

这是一个跨 skill 引用 gen-journeys 的 SKILL.md（不是 rules 子文件），但提案在 Scope 中仅列了"gen-journeys: 内联 gen-contracts/rules/journey-contract-model.md 所需内容"作为单向修复。实际上耦合是双向的：gen-journeys 引用 gen-contracts/rules/journey-contract-model.md（3次），gen-contracts 又引用 gen-journeys 的 SKILL.md Surface Detection 章节。提案的 Evidence 部分没有列出这个反向引用，导致修复方案可能遗漏。

### R2. 耦合图不完整——init-justfile 引用路径描述不精确

**问题**：提案 Evidence 称"init-justfile 引用 test-guide/references/test-type-model.md"。实际代码中 init-justfile SKILL.md 第 490 行写的是：

> `Test type terminology follows the [Surface Test Type Model](../test-guide/references/test-type-model.md).`

这是一个 Markdown 链接引用，且路径是 `../test-guide/...`（相对路径），跨 skill 引用。但提案 Scope 中将此项列为"内联 test-guide/references/test-type-model.md 所需内容"。该文件有 51 行，包含 Surface 到 Test Type 的完整映射表和分类标准。内联后需要确保 init-justfile 的上下文中这些映射关系仍然清晰——提案未说明需要内联该文件的全部内容还是摘要。

### R3. 内联保真度未定义——缺少"所需内容"的边界说明

**问题**：提案在 Scope 中使用了统一措辞"内联 XXX 所需内容"（如"gen-journeys: 内联 gen-contracts/rules/journey-contract-model.md 所需内容"），但没有定义"所需内容"的边界。journey-contract-model.md 有 184 行，其中包含：

- Core Concepts（Journey、Step、Contract、Outcome 定义）——gen-journeys 必须知道
- Semantic Descriptors 规则——gen-journeys 需要生成语义描述
- Contract File Format 和 Template——gen-journeys 不需要（这是 gen-contracts 的输出格式）
- Tag-Based Promotion——gen-journeys 不需要（这是 run-tests 的职责）
- Migration Guide（旧模型到新模型）——gen-journeys 不需要

如果全量内联 184 行，会增加 gen-journeys 约 47%（当前 392 行）的体积，且包含大量无关内容。如果选择性内联，提案缺少选择标准。

同理，run-tests/rules/test-isolation.md 有 138 行（含 4 条编号规则 TEST-isolation-000 到 TEST-isolation-004），ui-design/templates/styles/ 下有 7 个风格文件（每个 98-138 行）。提案未说明从 test-isolation.md 需要内联哪些规则，以及 extract-design-md 是否需要内联全部 7 个风格文件。

**建议**：为每个内联操作明确标注"内联哪些段落"或"内联至多 N 行"。

### R4. extract-design-md 的跨 skill 引用性质与其他不同

**问题**：extract-design-md SKILL.md 第 124 行写的是：

> `If "match built-in" is chosen, match against TUI theme characteristics per rules/platform-routing.md and read the corresponding style file from ui-design/templates/styles/<name>.md.`

这里的引用不是"加载知识以理解流程"——而是一个运行时动作指令：当用户选择匹配内置风格时，需要去读取 ui-design 的模板文件。这与其他 skill 的引用性质不同。其他 skill 引用外部文件是为了"理解概念"（如 gen-journeys 引用 journey-contract-model 以理解 Journey/Step 概念），而 extract-design-md 是"使用外部文件作为数据源"。

将风格匹配逻辑内联到 extract-design-md 意味着需要在 extract-design-md 中复制 7 个风格文件的内容（或至少匹配特征）。这不仅是知识内联，而是数据复制。每次 ui-design 新增风格模板时，都需要同步更新 extract-design-md——这正是提案所承认的"多份拷贝的漂移风险"的最坏场景。

### R5. Related Skills/Integration/References 章节清单不准确

**问题**：提案称"9 个 skill 的 Related Skills / Integration / References 章节"需要删除，并列出了"consolidate-specs, gen-contracts, gen-journeys, gen-test-scripts, run-tests, quick-tasks, tech-design, ui-design, write-prd"。但实际扫描发现：

1. **consolidate-specs** 有 `## Related Skills`（第 265 行）——符合
2. **gen-contracts** 有 `## Related Skills`（第 259 行）和 `## Reference`（第 267 行）——符合
3. **gen-journeys** 有 `## Related Skills`（第 376 行）和 `## Reference`（第 384 行）——符合
4. **gen-test-scripts** 有 `## Related Skills`（第 363 行）——符合
5. **run-tests** 有 `## Related Skills`（第 311 行）——符合
6. **quick-tasks** 有 `## Reference Files`（第 122 行，这是任务模板的引用说明，不是 pipeline 上下游信息）和 `## Integration`（第 340 行）——`## Reference Files` 不应删除，它是模板使用指南
7. **tech-design** 有 `## Integration`（第 438 行）——符合
8. **ui-design** 有 `## Integration`（第 229 行）——符合
9. **write-prd** 有 `## Integration`（第 446 行）——符合

quick-tasks 的 `## Reference Files` 是模板占位符使用说明（`{{REFERENCE_FILES}}` 的替换规则），不是 pipeline 上下游关系。提案将其归类为"Related 无用信息"会导致模板使用指南被误删。

### R6. 删除 Related/Integration/Reference 会丢失部分信息

**风险**：提案声称这些章节"内容均可从正文中隐含推断"。检查实际内容：

- gen-contracts 的 `## Reference`（第 267-276 行）列出了 Contract、Outcome、Semantic Descriptors、TUI Await Semantics、State Verification Levels、Batch Processing 的简明定义。这些概念在 SKILL.md 正文中被使用但未集中定义。删除后，阅读 gen-contracts 的开发者需要从其他 skill 的文件中找到这些定义——这恰恰违反了"独立加载"原则。

- gen-journeys 的 `## Reference`（第 384-392 行）同样如此：它集中定义了 Journey、Step、Risk Classification、Journey Invariants、Semantic Descriptors 五个核心概念。这些概念在 gen-journeys 正文中被使用但没有集中定义段落。在"内联 journey-contract-model.md 所需内容"之后，这些定义应来自内联内容，但提案没有明确说明内联后 Reference 段落的处理方式。

**建议**：对 gen-contracts 和 gen-journeys 的 Reference 段落，应将其内容合并到内联的知识中（作为定义段落），而非简单删除。

### R7. execute-task 与 run-tasks 的"重叠"判断需要重新审视

**问题**：提案称"execute-task 与 run-tasks command 60-70% 结构重叠"需要精简。实际对比发现：

- execute-task（149 行）是单任务执行器：claim → dispatch/execute → stop
- run-tasks（164 行）是循环调度器：set feature → claim loop → dispatch/verify → continue loop

两者共享 claim 输出解析格式和 fix-type derivation 表（约 20-30 行），但核心逻辑完全不同：execute-task 有 HARD-RULE 强制"ONE TASK PER INVOCATION"，run-tasks 有 Iron Laws 强制"claim → dispatch → continue loop"。它们不是"重叠"，而是"接口契约复用"——claim 输出格式和 fix-type 映射表是 CLI binary 定义的外部接口。

如果强行精简，可能破坏 command 的自洽性。一个 command 被单独加载时，需要看到完整的 claim 输出格式和错误处理表。

### R8. 漂移风险接受缺少量化依据

**风险**：提案在 Key Risks 中承认"多份拷贝的知识在未来修改时未同步更新"风险，并将 Mitigation 标注为"可接受——独立性带来的维护简化大于同步成本"。在 Assumption Challenged 表中也标注"Overturned: 对 AI agent 而言，独立加载更可靠；多份拷贝的漂移风险低于跨 skill 耦合的维护负担"。

但这个判断缺少量化依据。以 journey-contract-model.md 为例：

- 当前：1 份权威文件，被 gen-journeys、gen-contracts、gen-test-scripts 三个 skill 引用
- 内联后：知识分布在 3 个 skill 的 SKILL.md 中，如果修改 Journey/Step 的定义或 Tag-Based Promotion 规则，需要同步修改 3 处

当前 Forge 有 21 个 skill，如果每个内联操作平均产生 2 份拷贝，则有约 12 份额外拷贝需要维护同步。这不是一个可以无脑接受的 trade-off——需要考虑是否为关键模型定义（如 journey-contract-model）建立某种"权威标记"机制。

### R9. 成功标准的"功能等价"难以验证

**问题**：提案 Success Criteria 包含"所有 skill/command 修改后功能等价（无行为变更）"。但 skill 文件是自然语言指导 AI agent 行为的文档，不是可执行的程序。验证"功能等价"只能通过人工审读或 A/B 测试 agent 输出。

提案未提供验证方法。在"精简描述"操作中，"保留所有硬规则和决策表，只压缩描述性文字"的原则是合理的，但区分"描述性文字"和"行为指导性文字"的边界并不总是清晰的——特别是当描述性文字包含 agent 决策所需的上下文信息时。

---

## 3. Improvement Suggestions

### S1. 补全耦合图为双向图

**建议**：提案 Evidence 应构建完整的引用关系图，而非单向列表。当前遗漏了 gen-contracts → gen-journeys/SKILL.md 的反向引用。建议格式：

```
gen-journeys ←→ gen-contracts（双向）
  gen-journeys → gen-contracts/rules/journey-contract-model.md（概念引用，3 处）
  gen-contracts → gen-journeys/SKILL.md "Surface Detection"（流程引用，1 处）
```

完整耦合图能避免修复时的遗漏，也便于后续维护者理解原始耦合结构。

### S2. 为每个内联操作定义精确的压缩边界

**建议**：替换"内联 XXX 所需内容"为具体段落清单。例如：

```
gen-journeys 内联 journey-contract-model.md:
  - INJECT: Core Concepts（Journey/Step/Contract/Outcome 定义，约 40 行）
  - INJECT: Semantic Descriptors 规则（约 5 行）
  - INJECT: Directory Convention（约 15 行）
  - INJECT: Tag-Based Promotion 摘要（约 5 行，仅 @feature→@regression 流程）
  - SKIP: Contract File Format/Template（gen-contracts 职责）
  - SKIP: Migration Guide（一次性迁移参考）
  预计净增: ~65 行（替换当前 ~8 行引用说明）
```

这种格式使内联操作可审计、可回溯。

### S3. 重新审视 extract-design-md 的处理策略

**建议**：extract-design-md 引用 ui-design/templates/styles/ 的性质是"运行时数据读取"而非"知识引用"。考虑两种替代方案：

1. **保持引用但标注为设计意图**：在 extract-design-md SKILL.md 中保留对 ui-design 风格文件的引用路径，但添加注释说明这是跨 skill 的数据依赖（类似于 forensic 的动态加载豁免）。
2. **将匹配逻辑提取到共享的 rules 文件**：在 extract-design-md 内部创建 `rules/style-matching.md`，包含匹配特征摘要（而非完整风格文件内容）。这样 extract-design-md 知道何时触发匹配，实际风格数据仍从 ui-design/templates/styles/ 读取。

方案 2 更符合"独立知识单元"原则——extract-design-md 知道何时和如何匹配风格，但风格定义数据留在 ui-design 中。

### S4. 修正 Related/Integration/Reference 删除清单

**建议**：将删除清单区分为两类：

1. **应删除**：纯粹的 pipeline 上下游关系表（consolidate-specs Related Skills、gen-contracts Related Skills、gen-journeys Related Skills、gen-test-scripts Related Skills、run-tests Related Skills、quick-tasks Integration、tech-design Integration、ui-design Integration、write-prd Integration）
2. **应合并而非删除**：包含概念定义的 Reference 段落（gen-contracts Reference、gen-journeys Reference）——这些内容应在内联操作中被吸收到正文
3. **不应删除**：quick-tasks 的 `## Reference Files`——这是模板使用说明，不是 pipeline 关系信息

### S5. 为漂移风险建立轻量缓解机制

**建议**：对内联后的共享知识段落，在目标文件中使用统一标记约定：

```markdown
<!-- INLINE:origin=gen-contracts/rules/journey-contract-model.md#core-concepts -->
[Journey/Step/Contract 定义]
<!-- /INLINE:origin -->
```

这种标记有两个好处：(1) 后续维护者能识别哪些段落是内联的；(2) 如果将来需要同步更新，可以通过标记定位所有拷贝。这不会引入新的路径约定（标记是纯注释，不影响 agent 行为），但提供了可追溯性。

### S6. 为"功能等价"定义可操作的验证步骤

**建议**：将"功能等价"替换为具体检查清单：

- 所有 `<HARD-RULE>` 块完整保留（逐条对比原文）
- 所有 `<HARD-GATE>` 块完整保留
- 所有 `<EXTREMELY-IMPORTANT>` 块完整保留
- 所有 `<PROHIBITIONS>` 块完整保留
- 所有决策表/条件矩阵完整保留
- 所有 Step 序号和流程步骤完整保留（可重命名或合并，不可删除）
- 精简前后，每个 skill 的 `grep -c "HARD-RULE\|HARD-GATE\|EXTREMELY-IMPORTANT\|PROHIBITIONS"` 计数不变

### S7. 15% 行数减少目标需要重新评估

**建议**：提案声称"总行数减少 >= 15%"。当前 skill 文件总计约 6011 行，command 文件总计约 1127 行，合计 7138 行。15% 减少意味着净减约 1071 行。

- 删除 9 个 Related/Integration/Reference 段落约减少 50-70 行
- 精简冗余描述约减少多少行需要逐文件评估
- 但内联操作会**增加**行数：journey-contract-model.md 约 184 行需要部分内联，test-isolation.md 约 138 行需要部分内联，test-type-model.md 约 51 行需要部分内联，风格匹配逻辑的摘要约 30-50 行

净减少 1071 行意味着精简操作需要压缩掉约 1200-1400 行描述性文字，才能抵消内联增加的约 150-250 行。这对"只压缩描述性文字、保留所有硬规则"的原则构成了压力。建议将目标调整为"净减少 >= 10%"或"总行数不超过 6500 行"，以避免过度压缩。
