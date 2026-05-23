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
- spec-authority-enforcement eval 中，自由评审发现"标记稀释效应"和"Agent 层职责混淆"两个核心问题，Scorer 将两者均标为 `[beyond-rubric]`
- Scorer 的映射决策权是架构权衡（非 bug），但对高价值领域专家评审而言，信息损失不可接受

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

## Proposed Solution

在 Phase 0（freeform review + extract findings）完成后、Scorer 循环启动前，增加一个 **Pre-Revision 阶段**。Freeform findings 直接格式化为 ATTACK_POINTS，喂给现有 Reviser 执行修订。Scorer 对修订后的版本做盲审（不注入 freeform findings），独立验证修订质量。

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

**选项**：A. 绕过 Scorer | B. 双通道并行 | C. 前置修订

**选择 C**。先让 freeform findings 修正 proposal 的明显缺陷，然后 Scorer 再基于修订后的版本打分。信息不在 Scorer 映射中损失，而 Scorer 仍然作为质量关卡。

### Decision 2: Scorer 盲审——不看 freeform findings

Pre-revision 后 Scorer 只看到修订后的 proposal，不注入 freeform findings。

理由：Scorer 的价值在于用 rubric 独立评估当前版本质量。如果 pre-reviser 修好了问题，Scorer 自然给高分；如果修坏了，Scorer 自然扣分。带溯源可能引入确认偏误。

### Decision 3: 处理全部 findings（不限 severity）

所有 high/medium/low findings 都交给 pre-reviser。Pre-reviser 的指令中要求只接受可以从文档中验证的事实性修正，对不确定的建议保持 conservative。

理由：人为筛选 severity 会引入偏差。freeform 专家的核心价值是独立视角。

### Decision 4: 复用现有 Reviser，零新 protocol

Freeform findings 格式化为 `- **[severity]** summary | 原文引用: "quote"`（复用 `freeform-injection.md` 中已有的格式），直接作为 ATTACK_POINTS 喂给现有 Reviser。

不新建 "expert reviser" 角色，不注入专家 profile。Severity + summary + quote 已提供足够的权重信号和溯源。

### Decision 5: Pre-revision 占用 iteration 0，总预算不变

Pre-revision 计入 MAX_ITERATIONS。例如 `--iterations 3` 表示 iteration 0 = pre-revision + iteration 1-2 = Scorer 循环。

理由：总资源消耗可预测，用户不需要调整迭代习惯。

### Decision 6: 废弃 `freeform-injection.md`

Scorer 不再注入 freeform findings（盲审），因此 `rules/freeform-injection.md` 整个废弃。`scorer-composition.md` 中移除 `<injected-freeform-findings>` 块。

## Scope

### 改动文件

| 文件 | 改动类型 | 说明 |
|------|----------|------|
| `plugins/forge/skills/eval/SKILL.md` | 修改 | Phase 0 后新增 P0.5 Pre-Revision 步骤；P0.5（旧）改为设 `FREEFORM_INJECTION = false` |
| `plugins/forge/skills/eval/rules/freeform-injection.md` | 废弃 | Scorer 不再注入 freeform findings |
| `plugins/forge/skills/eval/rules/scorer-composition.md` | 修改 | 移除 `<injected-freeform-findings>` 组合步骤 |

### 不改动的文件

| 文件 | 理由 |
|------|------|
| `experts/protocol/reviser-protocol.md` | Reviser 本身不变，输入格式兼容 |
| `rules/reviser-composition.md` | 复用现有逻辑 |
| `experts/freeform/freeform-review-protocol.md` | Phase 0 前半段不变 |
| `experts/freeform/extraction-prompt.md` | Findings 提取逻辑不变 |
| `rules/freeform-expert-persistence.md` | 专家持久化不变 |

### 仅影响 proposal 类型

Freeform review 当前仅对 `type == proposal` 生效。Pre-revision 同样仅对 proposal 生效，不影响 prd/design/ui 等类型的 eval 流程。

## Key Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Pre-reviser 机械回应 findings，不理解专家深层意图 | 修订可能不够深入 | 盲审 Scorer 独立验证修订质量；findings 有 quote 溯源，Reviser 编辑有据可查 |
| Pre-revision 修坏 proposal | Scorer 给更低分 | Scorer 盲审自然扣分；后续 Reviser 可修复；rollback 机制兜底（INITIAL_SCORE 在 iteration 1 记录，对比 final score） |
| 某个 finding 本身就是错的（专家也会犯错） | 错误直接进了 proposal | 盲审 Scorer 作为独立安全网；pre-reviser 指令要求 conservative 修改 |
| 废弃 freeform-injection.md 是单向门 | 无法回退 | 重新创建该 rule 文件的成本极低（纯 prompt 组合规则，无逻辑代码） |

## Success Criteria

1. Phase 0 完成后，在 Scorer 启动前，findings 自动转换为 ATTACK_POINTS 并触发 Reviser 修订
2. Scorer 的 composed prompt 中不包含 freeform findings（盲审验证）
3. Pre-revision 的修改记录写入 iteration-0 报告
4. 最终 eval report 中体现 pre-revision 阶段的存在
5. 现有 degradation 路径不受影响（Phase 0 失败时跳过 pre-revision，直接进 Scorer）
