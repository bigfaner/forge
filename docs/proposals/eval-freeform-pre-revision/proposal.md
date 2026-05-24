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
- **unify-surfaces eval**：专家发现"Surface 与 Interface 应合并"，Scorer 映射为 medium-severity scope attack，Reviser 仅做措辞调整，未触及结构性问题
- 跨两次 eval 共 15 条 findings，7 条（47%）信息受损或丢失——rubric 维度与专家视角的结构性不匹配

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
3. **Scorer 忽略** → 完全丢失

## Industry Context

多阶段评审中如何整合专家输入是共性问题。以下三个行业实践为本提案提供设计参照：

**ACM SIGPLAN meta-reviewer 规则**：SIGPLAN 要求 meta-reviewer 必须原样传递 reviewer 的 core findings，不得按自己的 rubric 维度过滤或重映射。映射到 Forge：Scorer 当前扮演 meta-reviewer 角色，将 freeform findings 映射到 rubric 维度——正是 SIGPLAN 明确禁止的行为。本提案让 findings 直接到达 Reviser（等同于 reviewer findings 直达 author），Scorer 退回纯 rubric 评分角色（等同于 meta-reviewer 只做最终质量判定，不介入中间传递）。

**Gerrit Code-Review +2 机制**：Gerrit 将功能性评审（代码逻辑）和风格评审（格式规范）分属不同 reviewer，各自独立给出 +1/-1，两者都达到 +2 才能合入。映射到 Forge：Pre-revision 等价于功能性评审——专家 findings 指出实质问题，Reviser 执行修订。Scorer 盲审等价于风格评审——独立按 rubric 标准打分，不受专家视角影响。两道关卡独立通过，避免单一 reviewer 同时承担两种职责时的认知偏差。

**MT-Bench multi-judge 评估**：MT-Bench 使用多个独立 judge 分别评分再聚合，避免单个 judge 的系统性偏差。映射到 Forge：盲审 Scorer 作为独立 judge，不接触专家 findings，其评分代表 rubric 维度的纯质量评估。Pre-reviser 作为另一个 judge（专家视角），处理 findings 维度的问题。两个视角独立评估，最终通过 Gate 机制综合判定。

这些实践指向同一设计原则：**专家输入不应经过中间映射层才到达修订者**。本提案是该原则在 Forge eval 管道中的具体应用。

## Proposed Solution

在 Phase 0（freeform review + extract findings）完成后、Scorer 循环启动前，增加一个 **Pre-Revision 阶段**。Freeform findings 直接格式化为 ATTACK_POINTS，喂给现有 Reviser 执行修订。Scorer 对修订后的版本做盲审（不注入 freeform findings），独立验证修订质量。

**创新性定位**：本提案的核心思路——"将专家输入直接路由给修订者，评分者独立盲审"——是 peer review 领域的标准分离模式（edit/review separation），非原创洞察。提案的价值在于将此模式适配到 Forge eval 管道的具体约束中：合成 eval report 解决 Reviser protocol 依赖、三层 findings 策略平衡全面处理与保守修改、Scorer 盲审消除确认偏误的同时承认信息代价。这些实现层面的适配是提案的实际贡献。

### 新信息传递链

```
Phase 0: Freeform Review → Extract Findings
                ↓
Phase 0.5 (NEW): Pre-Revision
  findings 格式化 → ATTACK_POINTS → Reviser 修订
  占用 iteration 0（从总预算中扣除）
                ↓
Iteration 1..(MAX-1): Scorer(blind) → Gate → Reviser
  Scorer 不注入 freeform findings，独立按 rubric 打分
  INITIAL_SCORE 在 iteration 1 记录
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

### Decision 2: Scorer 盲审——不看 freeform findings

Pre-revision 后 Scorer 只看到修订后的 proposal，不注入 freeform findings。

理由：Scorer 的价值在于用 rubric 独立评估当前版本质量。如果 pre-reviser 修好了问题，Scorer 自然给高分；如果修坏了，Scorer 自然扣分。带溯源可能引入确认偏误。

**盲审的信息代价分析**：盲审确实让 Scorer 丧失了 pre-revision 变更上下文。具体而言，Scorer 无法区分"原文就有的缺陷"和"pre-revision 引入的新问题"——它只能评估当前版本的 rubric 维度质量。这是有意的设计取舍：Scorer 的职责是"当前版本是否达标"，而非"相比之前版本是改善还是恶化"。变更溯源的职责交给 iteration-0 报告（记录了每条 finding 的处理状态和编辑摘要），用户可据此追溯。若 Scorer 发现分数低于预期，后续 Reviser 循环仍有修复机会（`--iterations >= 3` 时）。

### Decision 3: 处理全部 findings，分策略响应

所有 findings 都交给 pre-reviser，按可验证性分三层：
- **事实性修正**（可定位原文缺陷）：直接编辑
- **结构/架构建议**：仅当存在直接矛盾时修改，否则标注 "deferred to Scorer cycle"
- **主观偏好**：标注 "not actionable"，不编辑

每条 finding 都被处理（计入处理率），但不一定都导致编辑。

### Decision 4: 复用现有 Reviser，最小 protocol 适配

Freeform findings 格式化为 `- **[severity]** summary | 原文引用: "quote"`，作为 ATTACK_POINTS 喂给现有 Reviser。不注入专家 profile。

**Protocol 适配**：Reviser protocol Step 1 要求 `EVAL_REPORT_PATH`，但 iteration 0 不存在 eval report。SKILL.md 编排层构造合成 eval report（`iteration: 0` + ATTACK_POINTS + 空 rubric）。Reviser 核心逻辑不变，但 SKILL.md 需新增约 20 行编排代码——这不是完全的"零新"，而是最小适配。

**空 rubric 兼容性验证**：Reviser protocol 的修订逻辑以 ATTACK_POINTS 为驱动（"Process each attack point"），rubric 分数仅用于 Gate 判定和 report 格式化。合成 report 中 rubric 所有维度标记为 N/A（分数 0 + 注释 "pre-revision: not scored"），Reviser 会跳过 rubric-grounded 分数解读，仅执行 ATTACK_POINTS 驱动的编辑。此行为已对照 Reviser protocol 验证：protocol 的 Step 2 明确以 ATTACK_POINTS 列表为输入循环处理，Step 3 的 rubric 对比仅在分数非 N/A 时触发。

### Decision 5: Pre-revision 占用 iteration 0，总预算不变

Pre-revision 计入 MAX_ITERATIONS。例如 `--iterations 3` 表示 iteration 0 = pre-revision + iteration 1-2 = Scorer 循环。

理由：总资源消耗可预测，用户不需要调整迭代习惯。

**用户可见变化**：
- Iteration-0 报告标题 "Pre-Revision (Freeform Findings)"，含每条 finding 处理状态及编辑摘要
- Iteration-1 为首次 Scorer 评分，INITIAL_SCORE 在此记录
- 最终 eval report summary 增加一行：`Pre-Revision: N findings processed (M accepted, K skipped)`

**`--iterations 2` 的设计权衡**：`--iterations 2` 是当前最低有用配置，pre-revision 将 Scorer 循环从 2 次减为 1 次（减少 50%）。这意味着仅一次 Scorer 评估 + 一次 Reviser 修订——如果 pre-revision 引入了问题，没有第二次 Scorer 循环来捕获和修复。这不是简单的用户选择，而是真实的质量退化风险。

应对策略：SKILL.md 在 `--iterations 2` 且 freeform review 生效时，输出 warning："Pre-revision 占用 1 个 iteration，Scorer 仅执行 1 轮评估。建议使用 `--iterations 3` 保证 Scorer 有修正机会。" 此 warning 不阻止执行，但明确告知用户行为变化。长期方案：评估是否将 pre-revision 不计入 iteration 预算（增加总 LLM 调用成本，但保持 Scorer 循环数不变）。

### Decision 6: 废弃 `freeform-injection.md`

Scorer 不再注入 freeform findings（盲审），因此 `rules/freeform-injection.md` 整个废弃。`scorer-composition.md` 中移除 `<injected-freeform-findings>` 块。

## Scope

### 改动文件

| 文件 | 改动类型 | 说明 |
|------|----------|------|
| `plugins/forge/skills/eval/SKILL.md` | 修改 | Phase 0 后新增 P0.5 Pre-Revision 步骤；P0.5（旧）改为设 `FREEFORM_INJECTION = false` |
| `plugins/forge/skills/eval/rules/freeform-injection.md` | 废弃 | Scorer 不再注入 freeform findings |
| `plugins/forge/skills/eval/rules/scorer-composition.md` | 修改 | 移除 `<injected-freeform-findings>` 组合步骤 |
| `plugins/forge/skills/eval/SKILL.md`（rollback 段） | 修改 | Rollback 对比点从 INITIAL_SCORE 改为 Phase 0 原始快照（见 Decision 5 baseline 保存） |

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

注：SKILL.md 的 rollback 对比点需修改（Phase 0 原始快照替代 INITIAL_SCORE），已列入改动文件表。Reviser protocol 和 composition 的核心逻辑不受影响。

### Phase 0.5 失败处理

| 失败场景 | 处理 |
|----------|------|
| Findings 格式化失败 | 跳过 pre-revision，直接进 Scorer |
| Pre-reviser 返回错误 | 跳过，记录 warning |
| 产出空报告 | 记录 iteration-0（"no changes"），Scorer 正常开始 |
| 产出格式异常 | 丢弃 pre-revision 结果，用原始 proposal 进 Scorer |

降级原则：任何 Phase 0.5 异常都回退到 Scorer 直接评估模式——即跳过 pre-revision，Scorer 不注入 freeform findings 直接开始 rubric 评分循环。注意：这不是回退到废弃前的"注入 findings"模式，而是回退到无注入的 Scorer 评估。此降级路径始终可用，不受 freeform-injection.md 废弃影响。

### 仅影响 proposal 类型

Freeform review 当前仅对 `type == proposal` 生效。Pre-revision 同样仅对 proposal 生效，不影响 prd/design/ui 等类型的 eval 流程。

**架构承诺**：废弃 `freeform-injection.md` 移除了 freeform findings 注入通道。若其他 eval 类型未来引入 freeform review，需重建注入机制。当前判断：废弃的短期收益高于重建成本——恢复需同步修改 freeform-injection.md + scorer-composition.md + SKILL.md P0.5 逻辑共 3 个文件，但每处修改均为 prompt 组合编排，无复杂逻辑代码。

## Implementation Estimate

| 步骤 | 估算 | 依据 |
|------|------|------|
| SKILL.md 新增 P0.5 编排（含合成 report 构造、baseline 保存、错误处理） | ~40 行代码 | 包含格式化 findings、构造合成 eval report、保存 Phase 0 快照、4 种失败场景的分支处理，原估算 20 行仅覆盖主路径 |
| SKILL.md rollback 对比点修改 | ~5 行 | 将 INITIAL_SCORE 引用改为 Phase 0 快照引用 |
| scorer-composition.md 移除注入块 | ~3 行删除 | 移除 `<injected-freeform-findings>` 占位符及关联逻辑 |
| 废弃 freeform-injection.md | 0 行（删除文件） | — |
| 端到端测试 | 1 个 proposal eval 全流程 | 验证 P0.5→iteration-1→rollback 全链路 |
| **总计** | **约 1 天工作量** | 含实现、自测、全流程验证 |

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Pre-reviser 机械回应 findings，不理解专家深层意图 | Medium | Medium | 预防：findings 强制包含原文 quote（`"quote"` 段），Reviser 必须引用原文上下文才能编辑，抑制脱离语境的表面修补。检测：盲审 Scorer 独立验证修订质量，对机械性修改自然扣分。Pre-reviser 指令显式标注："若 finding 的深层意图不明确，标注 deferred 而非猜测性编辑" |
| Pre-revision 修坏 proposal | Low | Medium | Scorer 盲审自然扣分；后续 Reviser 可修复；rollback 机制兜底 |
| 某个 finding 本身就是错的（专家也会犯错） | Medium | Low | 盲审 Scorer 作为独立安全网；pre-reviser 指令要求 conservative 修改 |
| 废弃 freeform-injection.md 是单向门 | Low | Medium | 恢复成本不仅是重建单个 rule 文件，还需同步恢复 scorer-composition.md 中的 `<injected-freeform-findings>` 占位符、SKILL.md 中 P0.5 的注入编排逻辑、以及 extraction-prompt.md 到 Scorer 的数据管道。依赖链横跨 3 个文件，恢复需全链路回归测试 |
| INITIAL_SCORE 基线漂移 | High | Medium | Rollback 对比点改为 Phase 0 原始快照（非 INITIAL_SCORE），pre-revision 前保存 baseline。Scope 已包含 SKILL.md rollback 段修改 |
| Reviser 处理 freeform findings 兼容性假设不成立 | Low | High | 合成 report 标注 `source: freeform`，兼容性问题出现时降级跳过 pre-revision |

## Success Criteria

1. Phase 0 完成后，在 Scorer 启动前，findings 自动转换为 ATTACK_POINTS 并触发 Reviser 修订
2. Scorer 的 composed prompt 中不包含 freeform findings（盲审验证）
3. Pre-revision 的修改记录写入 iteration-0 报告
4. 最终 eval report 包含 "Pre-Revision" 独立章节，列出每条 finding 的处理状态（accepted/partially-accepted/skipped）及对应编辑摘要
5. 现有 degradation 路径不受影响（Phase 0 失败或 Phase 0.5 异常时跳过 pre-revision，直接进 Scorer）
6. Pre-revision 的质量可度量：iteration-1 Scorer 盲审评分 >= 同一 proposal 无 pre-revision 的 Scorer 盲审评分。Baseline 获取方式：对同一 proposal 先运行一次 `--iterations 1 --no-freeform`（跳过 freeform review），记录 Scorer 盲审分数作为 baseline；再运行完整 pre-revision 流程，取 iteration-1 的 Scorer 盲审分数进行对比。两次运行均使用盲审模式，消除评审模式差异的混淆因素。已知局限：baseline 运行无 freeform review（故无 findings 注入），pre-revision 运行有 freeform review 但 findings 已被 pre-revision 消费——两次运行的 Scorer 均不看 findings，对比公平。但 pre-revision 版本的文本长度和结构可能系统性差异（如新增 Industry Context 段落），Scorer 可能因此给出不同分数。此混淆因素通过控制同一 proposal 内容来缓解
7. Freeform findings 中标为 high-severity 的条目，pre-revision 处理率 >= 80%（accepted + partially-accepted）
