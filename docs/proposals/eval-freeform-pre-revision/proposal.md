---
created: 2026-05-24
author: faner
status: Draft
---

# Proposal: Freeform Pre-Revision——让自由专家发现直接修订 proposal

## Problem

Eval-proposal 的 Phase 0 自由专家评审产出的 findings 当前通过 Scorer 间接传递给 Reviser。Scorer 作为中间层有权映射、标记 `[beyond-rubric]`、或忽略 freeform findings，导致信息在传递过程中被压缩或丢失。

### Evidence

- `docs/lessons/arch-freeform-findings-indirect-influence.md` 记录了完整的信息传递链分析和根因
- **spec-authority-enforcement eval**：8 条 freeform findings 中 3 条标 `[beyond-rubric]`（37.5%），2 条被降级为 low severity（25%），核心发现"标记稀释效应"在 Reviser 修订中未被触及
- **unify-surfaces eval**：7 条 freeform findings 中，1 条核心发现"Surface 与 Interface 应合并"被 Scorer 映射为 medium-severity scope attack（语义压缩：结构合并问题被降级为措辞问题），1 条关于路径解析一致性的发现被标 `[beyond-rubric]`（完全丢失）
- 跨两次 eval 共 15 条 findings，7 条（47%）信息受损或丢失（spec-authority: 3 beyond-rubric + 2 severity 降级 = 5 条；unify-surfaces: 1 语义压缩 + 1 完全丢失 = 2 条）——rubric 维度与专家视角的结构性不匹配

### Urgency

2 个活跃 proposal 受此影响。Freeform review 已是 eval-proposal 标准阶段，47% 信息损失率在每个 proposal eval 中复现。

### 信息传递链（当前）

```
Freeform Review → 提取 findings → 注入 Scorer Prompt
                                        ↓
                                  Scorer 映射/过滤 → ATTACKS
                                        ↓
                                  Reviser 修订（看不到原始 findings）
```

### 信息损失的三种路径

1. **映射到 rubric 维度** → 成为标准 attack point，但语义可能被压缩
2. **标为 `[beyond-rubric]`** → Reviser 可能优先级低于 rubric 维度内的攻击点
3. **Scorer 忽略** → 完全丢失（当前管道中无审计机制）

## Industry Context

多阶段评审中如何整合专家输入是共性问题。以下三个行业实践为本提案提供设计参照：

**ACM SIGPLAN meta-reviewer 规则**：SIGPLAN 要求 meta-reviewer 必须原样传递 reviewer 的 core findings，不得按自己的 rubric 维度过滤或重映射。映射到 Forge：Scorer 当前扮演 meta-reviewer 角色，将 freeform findings 映射到 rubric 维度——正是 SIGPLAN 明确禁止的行为。**因果论证**：SIGPLAN 禁止重映射的根因是 meta-reviewer 的 rubric 维度与 reviewer 的评审视角存在结构性不匹配——这与 Forge 的问题完全同构（47% findings 受损正是因为 Scorer 的 rubric 维度无法覆盖专家视角）。本提案采用相同分离原则：findings 直达 Reviser，Scorer 退回纯 rubric 评分角色。类比成立：47% 受损率支持了维度不匹配前提，且分离原则不依赖评审者是人类还是 LLM。证伪条件：若后续 eval 的 findings 受损率降至 15% 以下（即 Scorer 映射层能保留绝大多数 findings），则维度不匹配前提不成立，分离策略的必要性需重新评估。

**Gerrit Code-Review +2 机制**：Gerrit 将功能性评审（代码逻辑）和风格评审（格式规范）分属不同 reviewer，各自独立给出 +1/-1，两者都达到 +2 才能合入。映射到 Forge：Pre-revision 等价于功能性评审——专家 findings 指出实质问题，Reviser 执行修订。Scorer 标注盲审等价于风格评审——独立按 rubric 标准打分，同时通过变更标记获得必要的上下文。两道关卡独立通过，避免单一 reviewer 同时承担两种职责时的认知偏差。**因果论证**：Gerrit 双 reviewer 的核心收益是消除角色冲突——易量化的风格问题挤占需深度理解的功能问题。Forge 的 Scorer 同时处理 rubric 评分和 findings 映射时同样存在此冲突（结构化的 rubric 维度挤占非结构化 findings），Gerrit 分离策略在此适用。

**MT-Bench multi-judge 评估**：MT-Bench 使用多个独立 judge 分别评分再聚合，避免单个 judge 的系统性偏差。映射到 Forge：盲审 Scorer 作为独立 judge，不接触专家 findings，其评分代表 rubric 维度的纯质量评估。Pre-reviser 作为另一个 judge（专家视角），处理 findings 维度的问题。两个视角独立评估，最终通过 Gate 机制综合判定。**因果论证**：MT-Bench 多 judge 的根因是单 judge 的系统性偏差。Forge 的 Scorer 对 rubric 维度敏感但对 rubric 外发现系统性忽略（47% 受损率），引入独立 Pre-reviser 覆盖此盲区。维度正交性由标注盲审设计保证。

**类比失效点：LLM vs 人类专家可信度**。上述三个行业实践均假设评审者是可信的人类专家。Forge 的 freeform expert 是 LLM——可能产生幻觉、过度泛化或遗漏关键约束。直接路由未经 Scorer 过滤的 LLM 发现给 Reviser，与路由人类专家发现是性质不同的命题。设计应对：Pre-reviser 的三层分类策略（Decision 3）本质上是 LLM 输出的可信度分层——事实性修正（可验证）优先执行，主观偏好（不可验证）跳过。Scorer 标注盲审（Decision 2）作为第二道防线：即使 pre-reviser 基于 LLM 幻觉做了错误修改，Scorer 对 `<!-- pre-revised -->` 标记区域的检查也能捕获异常。

## Proposed Solution

在 Phase 0（freeform review + extract findings）完成后、Scorer 循环启动前，增加一个 **Pre-Revision 阶段**。Freeform findings 直接格式化为 ATTACK_POINTS，喂给现有 Reviser 执行修订。修订后的段落以 HTML 注释标注变更区域（`<!-- pre-revised: {severity} -->`）。Scorer 对修订后的版本做标注盲审——不注入 freeform findings，但可见变更区域标记，独立验证修订质量。

**术语约定**：本文中 "Reviser" 指使用现有 Reviser protocol（`experts/protocol/reviser-protocol.md`）的 agent 实例。"Pre-revision 阶段的 Reviser" 是同一个 Reviser 在 Phase 0.5 中的调用，非独立角色或新 protocol。为行文简洁，部分段落使用 "Reviser（Phase 0.5）" 或简称 "Pre-revision Reviser" 指代此调用上下文。

**创新性定位**：本提案的核心思路——"将专家输入直接路由给修订者，评分者独立盲审"——是 peer review 领域的标准分离模式（edit/review separation），非原创洞察。提案的实际贡献在于适配到 Forge eval 管道的具体约束中：合成 eval report 解决 Reviser protocol 依赖、三层 findings 策略平衡全面处理与保守修改、**标注盲审**在消除确认偏误的同时保留变更上下文（避免完全盲审的信息丢失）。

### 新信息传递链

```
Phase 0: Freeform Review → Extract Findings
                ↓
Phase 0.5 (NEW): Pre-Revision
  findings 格式化 → ATTACK_POINTS → Reviser 修订
  修订后的段落标注 <!-- pre-revised: {severity} -->
  占用 iteration 0（从总预算中扣除）
                ↓
Iteration 1..(MAX-1): Scorer(标注盲审) → Gate → Reviser
  Scorer 不注入 freeform findings，但可见变更区域标记
  Scorer 对标记区域检查修订质量，对未标记区域正常评估
  INITIAL_SCORE 在 iteration 1 记录
  Rollback: 循环内回滚到 pre-revised checkpoint；整体回滚到 Phase 0 baseline snapshot
                ↓
Final Report + Rollback（不变）
```

## Design Decisions

### Decision 1: Pre-Revision 发生在 Scorer 前（而非替代或并行）

**选项**：D. 保持现状 | A. 绕过 Scorer | B. 双通道并行 | C. 前置修订

- **D. 保持现状**：零改动，但 47% findings 受损且无法通过调优 Scorer prompt 解决（结构性不匹配）
- **A. 绕过 Scorer**：最大保留专家输入，但失去 rubric 标准化质量关卡
- **B. 双通道并行**：不损失 Scorer 循环，但映射问题未解决且合并逻辑复杂
- **C. 前置修订**：Scorer 不再是信息瓶颈，盲审保证独立性，代价是消耗一个 iteration 预算

**选择 C**。选项 D 的维持成本随 proposal 数量累积，C 是一次性结算。

### Decision 2: Scorer 标注盲审——看变更区域，不看 findings

**标注盲审 (annotated blind review)**：Scorer 可见文档中被 `<!-- pre-revised: {severity} -->` HTML 注释标记的变更区域，但不可见触发修改的原始 freeform findings 内容。这是介于完全盲审（不看任何变更信息）和全量溯源（看完整 findings）之间的折中方案。

Pre-revision 后 Scorer 不注入 freeform findings，但 proposal 文档中 pre-revision 修改过的段落会被标注 `<!-- pre-revised: {finding_severity} -->` HTML 注释标记。

**为什么不是完全盲审**：完全盲审让 Scorer 无法区分"原文就有的缺陷"和"pre-revision 引入的新问题"。如果 pre-reviser 删除了一个有争议段落，盲审 Scorer 会基于段落缺失生成新攻击点，浪费已减少的 iteration 预算。这是用信息丢失解决信息丢失——旧的损失在 Scorer 映射层，新的损失在 Scorer 盲审层。

**为什么不是全量溯源**：全量溯源（将 findings 原文注入 Scorer）会引入确认偏误——Scorer 倾向于在 findings 指出的区域寻找问题，即使 pre-revision 已经修正。

**标注盲审的折中设计**：
- Scorer 知道**哪些区域被改过**（`<!-- pre-revised: high -->`），但不知道**为什么改**（原始 findings 内容不暴露）
- Scorer 对标注区域的评估策略：检查修订是否引入了新问题，而非重新评估原始问题
- severity 标记帮助 Scorer 分配注意力权重：high 区域值得更仔细检查，low 区域快速扫过
- HTML 注释对最终报告不可见（不污染文档正文），仅在 Scorer 读取文档时可见

**标注偏误检测机制**：severity 标记引入了一个隐式信息通道——`<!-- pre-revised: high -->` 等价于告知 Scorer"专家认为此区域有高严重性问题并已修复"，这可能在方向上复现提案本想消除的确认偏误。Prompt 指令约束（"severity 标记供注意力分配参考，不影响评分标准"）不足以消除此风险。因此增加实证检测：eval report 要求 Scorer 对每个标注区域单独记录 attack density（单位文本长度内的攻击点数量），并与未标注区域的 attack density 比较。若标注区域 density 系统性偏高（连续 >= 2 次 eval 中偏差超过 30%），触发"标注偏误告警"，需调整标注策略或移除 severity 标记。此机制为标注盲审的折中设计提供实证反馈闭环。

**标记生命周期**：`<!-- pre-revised: {severity} -->` 标记的生命周期分为三个阶段：
1. **生成**：Phase 0.5 Pre-Revision 完成后，由 SKILL.md 编排层在 Reviser 编辑的段落插入标记。Scorer-Reviser 循环（iteration 1..N）中 Reviser 的后续修改不追加标记（这些修改由 rubric 攻击点驱动，已有 iteration report 追踪）
2. **存活**：标记在 Scorer 评审期间有效。Scorer 每次读取文档时均可见标记
3. **清除**：eval 流程结束后（Final Report 生成 + Rollback 决定完成），SKILL.md Step 5 清理动作一次性剥离文档中所有 `<!-- pre-revised: ... -->` 注释。同时，Step 1.4（pre-processing）在每次 eval 启动时扫描并清除残留标记，防止旧标记误导新的 Scorer

**Scorer prompt 补充指令**（约 5 行，加入 scorer-composition.md）：
```
<!-- pre-revised: {severity} --> 标记表示该段落经过 Pre-Revision 修改。
对标记区域：关注修订是否引入了新问题或遗漏，而非重新评估已修正的原始问题。
severity 标记供注意力分配参考，不影响评分标准。
在 eval report 中分别记录标注区域与未标注区域的 attack density，供偏误检测。
当 Scorer 的 rubric 判断与 pre-revision 修改方向矛盾时（如 Scorer 认为某段应删除但 pre-revision 刚添加），以 rubric 标准为准生成 attack point，但在 attack point 中标注 `conflict-with-pre-revision` 供审查。
```

### Decision 3: 处理全部 findings，分策略响应

所有 findings 都交给 pre-reviser，按可验证性分三层：
- **事实性修正**（可定位原文缺陷）：直接编辑
- **结构/架构建议**：仅当存在直接矛盾时修改，否则标注 "deferred to Scorer cycle"
- **主观偏好**：标注 "not actionable"，不编辑，但写入 iteration-0 报告的"分类审计"章节，附分类理由和原始 finding 摘要

**分类边界审查**：当 finding 不明确属于某一层时（如领域专业知识密集型 finding），Pre-reviser 必须标注为 "borderline" 并保留在 ATTACK_POINTS 中（降级为 deferred 而非直接标记 not actionable），避免因领域知识不足导致的静默误分类。Borderline findings 在 iteration-0 报告中单独列出，用户可据此判断是否需要手动干预。

每条 finding 都被处理——"处理"指"分诊"（triage）：每条 finding 都被评估、分类、记录，但不一定都导致编辑。"not actionable" findings 不会从管道中消失：iteration-0 报告列出所有 findings 的分类结果（accepted / partially-accepted / deferred / skipped），skipped 项附分类理由和原始 finding 摘要，用户可在报告中追溯任何 finding 的处置路径。

### Decision 4: 复用现有 Reviser，受控 protocol 扩展

Freeform findings 格式化为 `- **[severity]** summary | 原文引用: "quote" | 期望改进方向: <动词短语>`，作为 ATTACK_POINTS 喂给现有 Reviser。"期望改进方向"字段为 Reviser 提供方向锚点，帮助它在 fix strategy 表格中选择更合适的策略，但不回退到注入完整 narrative。不注入完整专家 profile，但 findings 格式中包含专家视角的浓缩信息（severity + 改进方向）。

**Protocol 适配**：Reviser protocol Step 1 要求 `EVAL_REPORT_PATH`，但 iteration 0 不存在 eval report。SKILL.md 编排层构造合成 eval report（`iteration: 0` + ATTACK_POINTS + 空 rubric）。Reviser 核心逻辑不变，但 SKILL.md 需新增约 40 行编排代码（含格式化 findings、构造合成 eval report、保存 Phase 0 快照、4 种失败场景的分支处理）。

**空 rubric 兼容性论证**：Reviser protocol 的标题即声明 "Generic attack-point-driven revision workflow"——修订逻辑以 ATTACK_POINTS 为唯一驱动，rubric 分数不参与 Reviser 的执行路径。验证：
- **Protocol Step 1 (Read Inputs)**：读取 DOC_DIR 中所有 markdown 文件 + eval report。合成 report 包含 `iteration: 0`、ATTACK_POINTS 列表、rubric 全维度 N/A——满足 Step 1 的最小输入格式要求。Step 1 的 `<HARD-RULE>` 要求 "Do NOT skip reading the eval report"，合成 report 满足此约束
- **Protocol Step 2 (Edit by Attack Point)**：以 ATTACK_POINTS 列表为输入循环处理。合成 report 的 ATTACK_POINTS 格式与正常 eval report 完全一致（`- [dimension]: [weakness] — [quote] — [improve]`），Reviser 正常执行编辑。Rubric 分数在 Reviser 执行路径中无任何消费者——Reviser 不读分数、不比分数、不用分数决策
- **Protocol Step 3 (Report)**：输出修改摘要，不含 rubric 操作
- **Fallback**：若 Reviser 遇到无法处理的 report 格式，Phase 0.5 失败处理捕获异常，降级为跳过 pre-revision

### Decision 5: Pre-revision 占用 iteration 0，总预算不变

Pre-revision 计入 MAX_ITERATIONS。例如 `--iterations 3` 表示 iteration 0 = pre-revision + iteration 1-2 = Scorer 循环。

理由：总资源消耗可预测，用户不需要调整迭代习惯。

**用户可见变化**：
- Iteration-0 报告标题 "Pre-Revision (Freeform Findings)"，含每条 finding 处理状态及编辑摘要
- Iteration-1 为首次 Scorer 评分，INITIAL_SCORE 在此记录
- 最终 eval report summary 增加一行：`Pre-Revision: N findings triaged (M accepted, D deferred, K skipped)`

**`--iterations 2` 的设计权衡**：`--iterations 2` 是当前最低有用配置，pre-revision 将 Scorer 循环从 2 次减为 1 次（减少 50%）。这意味着仅一次 Scorer 评估 + 一次 Reviser 修订——如果 pre-revision 引入了问题，没有第二次 Scorer 循环来捕获和修复。质量退化预期：pre-revision 可能误编辑某段落（如基于 LLM 幻觉的分类错误），单次 Scorer 循环需同时完成"检测 pre-revision 问题"和"检测原始 proposal 问题"双重任务，攻击点密度翻倍但修复机会减半。推荐最低配置为 `--iterations 3`（1 次 pre-revision + 2 次 Scorer 循环），`--iterations 2` 仅适用于低风险 proposal（如纯格式/措辞类 proposal，无架构决策变更）。

应对策略：SKILL.md 在 `--iterations 2` 且 freeform review 生效时，输出 warning："Pre-revision 占用 1 个 iteration，Scorer 仅执行 1 轮评估。建议使用 `--iterations 3` 保证 Scorer 有修正机会。" 此 warning 不阻止执行，但明确告知用户行为变化。长期方案：评估是否将 pre-revision 不计入 iteration 预算（增加总 LLM 调用成本，但保持 Scorer 循环数不变）。

### Decision 6: 条件性废弃 `freeform-injection.md`

Scorer 不再注入 freeform findings（盲审），但采用条件性废弃而非物理删除以降低单向门风险。`rules/freeform-injection.md` 保留，在文件头部添加 `status: deprecated` frontmatter 标记，保留完整的注入语义定义。`scorer-composition.md` 中的 freeform injection 块改为条件分支：当 `FREEFORM_INJECTION = false`（即 pre-revision 模式生效时）跳过注入，否则走原路径。条件分支在 SKILL.md 编排层判断（非每次 Scorer 调用时执行），不增加运行时复杂度。若未来其他 eval 类型引入 freeform review，可通过配置启用注入通道，无需重建。

## Scope

### 改动文件

| 文件 | 改动类型 | 说明 |
|------|----------|------|
| `plugins/forge/skills/eval/SKILL.md` | 修改 | Phase 0 后新增 P0.5 Pre-Revision 步骤；原 P0.5 (Inject Findings) 改为设置 `FREEFORM_INJECTION = false`；Iteration Initialization 改为 `ITERATION = 0`，pre-revision 后递增为 1，Scorer 循环从 iteration 1 开始；新增 BASELINE_SCORE 单次 Scorer 评估（SC #6 baseline） |
| `plugins/forge/skills/eval/rules/freeform-injection.md` | 条件性废弃 | 添加 `status: deprecated` frontmatter，保留完整语义定义供未来恢复 |
| `plugins/forge/skills/eval/rules/scorer-composition.md` | 修改 | 移除 `<injected-freeform-findings>` 组合步骤，新增 `<!-- pre-revised -->` 标注盲审指令 + 偏误检测 report template（attack density 分区域输出格式） |
| `plugins/forge/skills/eval/SKILL.md`（rollback + cleanup 段） | 修改 | 两级 rollback 设计：(1) Scorer 循环内 rollback 恢复到 pre-revised checkpoint；(2) 整体流程 rollback 恢复到 Phase 0 原始快照。Baseline snapshot 在 Phase 0.5 启动前保存到 `<DOC_DIR>/eval/baseline-snapshot/`。Step 5 新增标记清理：eval 结束后剥离所有 `<!-- pre-revised -->` 注释；Step 1.4 新增启动时残留标记清除 |

### Non-Functional Requirements

- **延迟**：增加一次 LLM 调用（~30-60s），与单次 Reviser iteration 相当。总 iteration 数不变
- **兼容性**：仅影响 `type == proposal`。`--iterations 1` 时跳过 pre-revision，行为不变
- **安全**：修改权限与现有 Reviser 一致，无外部输入注入风险

### 不改动的文件

| 文件 | 理由 |
|------|------|
| `experts/protocol/reviser-protocol.md` | Reviser 核心逻辑不变；SKILL.md 编排层构造合成 eval report 满足其 `EVAL_REPORT_PATH` 依赖（见 Decision 4） |
| `rules/reviser-composition.md` | 复用现有逻辑 |
| `experts/freeform/freeform-review-protocol.md` | Phase 0 前半段不变 |
| `experts/freeform/extraction-prompt.md` | Findings 提取逻辑不变 |
| `rules/freeform-expert-persistence.md` | 专家持久化不变 |

注：SKILL.md 的 rollback 已升级为两级设计（Scorer 循环内回滚到 pre-revised checkpoint，整体回滚到 Phase 0 baseline snapshot），已列入改动文件表。Reviser protocol 和 composition 的核心逻辑不受影响。

### Phase 0.5 失败处理

| 失败场景 | 处理 |
|----------|------|
| Findings 格式化失败 | 跳过 pre-revision，直接进 Scorer |
| Pre-reviser 返回错误 | 跳过，记录 warning |
| 产出空报告 | 记录 iteration-0（"no changes"），Scorer 正常开始 |
| 产出格式异常 | 丢弃 pre-revision 结果，用原始 proposal 进 Scorer |

降级原则：任何 Phase 0.5 异常都回退到 Scorer 直接评估模式——即跳过 pre-revision，Scorer 不注入 freeform findings 也不看标注，直接开始 rubric 评分循环。此降级路径始终可用，不受 freeform-injection.md 废弃影响。

### 仅影响 proposal 类型

Freeform review 当前仅对 `type == proposal` 生效。Pre-revision 同样仅对 proposal 生效，不影响 prd/design/ui 等类型的 eval 流程。

**架构承诺**：条件性废弃 `freeform-injection.md` 保留了 freeform findings 注入通道的语义定义，但默认禁用。若其他 eval 类型未来引入 freeform review，可通过 scorer-composition.md 中的条件分支重新启用注入机制，无需重建语义定义。恢复路径：移除 `status: deprecated` frontmatter + 移除 scorer-composition.md 条件分支 = 2 处配置变更。

## Implementation Estimate

| 步骤 | 估算 | 依据 |
|------|------|------|
| SKILL.md 新增 P0.5 编排（含合成 report 构造、baseline 保存、错误处理） | ~40 行代码 | 包含格式化 findings、构造合成 eval report、保存 Phase 0 快照、4 种失败场景的分支处理，原估算 20 行仅覆盖主路径 |
| SKILL.md rollback 对比点修改 | ~15 行 | 两级 rollback 语义设计：保存 Phase 0 baseline snapshot、pre-revised checkpoint；Scorer 循环内回滚到 pre-revised，整体回滚到 baseline |
| scorer-composition.md 替换注入块为标注盲审指令 + 偏误检测 report template | ~10 行替换 | 移除 `<injected-freeform-findings>` 占位符，新增 `<!-- pre-revised -->` 标注解读指令（~5 行）+ 新增标注/未标注区域 attack density 输出格式（~5 行） |
| 条件性废弃 freeform-injection.md | ~3 行 | 添加 `status: deprecated` frontmatter + scorer-composition.md 条件分支 |
| 端到端测试 | 1 个 proposal eval 全流程 | 验证 P0.5→iteration-1→rollback 全链路 |
| **总计** | **约 1 天工作量** | 含实现、自测、全流程验证 |

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Pre-reviser 机械回应 findings，不理解专家深层意图 | Medium | Medium | 预防：findings 强制包含原文 quote（`"quote"` 段），Reviser 必须引用原文上下文才能编辑，抑制脱离语境的表面修补。检测：标注盲审 Scorer 对 `<!-- pre-revised -->` 区域检查修订质量，对机械性修改自然扣分。Pre-reviser 指令显式标注："若 finding 的深层意图不明确，标注 deferred 而非猜测性编辑" |
| Pre-revision 修坏 proposal | Low | Medium | 标注盲审 Scorer 对修改区域重点检查；后续 Reviser 可修复；rollback 机制兜底 |
| 某个 finding 本身就是错的（LLM 专家幻觉） | Medium | Low | 标注盲审 Scorer 作为独立安全网；pre-reviser 三层分类策略过滤不可验证项；`<!-- pre-revised -->` 标记让 Scorer 知道哪些是新增内容 |
| 废弃 freeform-injection.md 是单向门 | Low | Low | 采用条件性废弃（`status: deprecated` frontmatter + scorer-composition.md 条件分支），保留低成本回退路径。恢复仅需修改 frontmatter 和移除条件分支，无需重建语义定义 |
| INITIAL_SCORE 基线漂移 | High | Medium | 两级 rollback：(1) Scorer 循环内回滚到 pre-revised checkpoint；(2) 整体流程回滚到 Phase 0 baseline snapshot。Baseline 在 Phase 0.5 启动前保存。漂移检测：若 iteration-1 的 `INITIAL_SCORE` 低于 Phase 0.5 前的 `BASELINE_SCORE` 超过 50 分（1000 分制下 5% 偏差，超过典型 LLM 评分方差），在 eval report 中标注"基线漂移告警"供人工审查，但不自动触发 rollback（rollback 由最终分数与 INITIAL_SCORE 的对比决定） |
| 标注盲审假阳性——Scorer 对 pre-revised 区域过度审查 | Medium | Medium | Scorer prompt 明确指令：对标记区域关注"修订是否引入新问题"，而非"重新评估原始问题"。severity 标记仅作注意力分配参考，不影响评分标准。若实测假阳性率高，可调整为仅标记 high-severity 区域 |
| Reviser 处理 freeform findings 兼容性假设不成立 | Low | High | 合成 report 标注 `source: freeform`，兼容性问题出现时降级跳过 pre-revision |
| Borderline 分类元认知失败——Pre-Reviser 过度自信，应标 borderline 的 finding 被错误归入确定分类 | Medium | Medium | 分类审计章节列出所有分类结果供人工审查；当 borderline 率异常低（如 0 borderline / >10 findings）时触发告警，提示分类可能过于自信 |

## Success Criteria

1. Phase 0 完成后，在 Scorer 启动前，findings 自动转换为 ATTACK_POINTS 并触发 Reviser 修订
2. Scorer 的 composed prompt 中不包含 freeform findings，但包含 `<!-- pre-revised -->` 标注解读指令（标注盲审验证）
3. Pre-revision 的修改记录写入 iteration-0 报告
4. 最终 eval report 包含 "Pre-Revision" 独立章节，列出每条 finding 的处理状态（accepted/partially-accepted/deferred/skipped）及对应编辑摘要。Skipped findings 附分类理由和原始 finding 摘要（审计轨迹）
5. 现有 degradation 路径不受影响（Phase 0 失败或 Phase 0.5 异常时跳过 pre-revision，直接进 Scorer）
6. Pre-revision 的质量可度量（信息指标，非门控标准）：对比同一 proposal 有无 pre-revision 的 Scorer 评分。Baseline 获取方式：在 Phase 0.5 启动前，SKILL.md 编排层对原始 proposal 运行一次仅 Scorer 评估（单次 Scorer subagent 调用，不触发 Reviser，不消耗 iteration 预算），记录分数作为 `BASELINE_SCORE`。Pre-revision 完成后 Scorer 循环的 iteration-1 分数（`INITIAL_SCORE`）与 `BASELINE_SCORE` 对比。**方法论局限与应对**：`BASELINE_SCORE` 基于 pre-revision 前的原始文档，`INITIAL_SCORE` 基于 pre-revision 后的文档 + `<!-- pre-revised -->` 标注，两者存在文本长度和结构差异，不是严格的 A/B 测试。因此本指标降级为信息性参考（informational metric），不作为 pass/fail 门控。当 `INITIAL_SCORE < BASELINE_SCORE` 时，触发人工审查：检查 iteration-0 报告中 edits 的具体内容，判断低分是因为编辑引入了问题（需回滚或修复）还是因为评分方差（可接受）
7. Freeform findings 中标为 high-severity 的条目，pre-revision 分诊率 >= 80%（accepted + partially-accepted + deferred），其中 accepted + partially-accepted >= 60%（确保大多数高严重性 findings 得到实质性处理，而非仅分诊）。"partially-accepted"判定标准：修改触及了 finding 指出的原文位置，且修改方向与 finding 的"期望改进方向"字段一致。当 partially-accepted 比例超过 accepted 时，触发人工抽检：审查 partially-accepted edits 的实际质量，排除虚增
