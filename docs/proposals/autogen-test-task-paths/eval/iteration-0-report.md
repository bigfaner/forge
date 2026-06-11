---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision (Freeform Findings)

## ATTACK_POINTS

- **high** Embed 与 Prompt 模板的 FeatureSlug 来源不同（CLI 目录路径 vs dispatcher index.json），目录重命名后会产生不一致 | quote: "Embed 模板和 Prompt 模板的 `FeatureSlug` 来源不同，但提案未分析两者不一致的场景。" | improvement: 分析两个 FeatureSlug 来源的一致性保障机制，明确说明为何不会不一致
- **medium** FeatureSlug 为空时 embed 模板生成无效命令 `ls docs/features//testing/`，提案将影响评为 Low 偏低 | quote: "Proposal 将 'FeatureSlug 渲染为空' 的 Impact 评为 Low，但未分析空值在 embed 模板中产生的具体后果。" | improvement: 分析空 slug 在 embed 模板中的实际后果，更新风险评估
- **medium** 提案对非统一的模板舰队施加统一改动，薄模板（test-run, test-gen-scripts）与富模板（test-gen-journeys, test-gen-contracts）的价值差异被忽略 | quote: "Proposal 声称改动仅涉及 6 个 embed 模板和 6 个 prompt 模板，但实际 embed 模板舰队中存在结构性差异被忽略。" | improvement: 区分薄模板和富模板的改动策略，对已有路径引用的模板避免冗余
- **medium** Embed 模板 discovery 命令路径与 skill 路径独立演进，缺少同步验证机制 | quote: "Embed 模板的 discovery 命令路径约定 (`docs/features/<slug>/testing/`) 与 skill 中的路径可能不同步。" | improvement: 添加同步验证机制或明确声明路径约定为单一权威来源
- **medium** Agent 可能在 Step 1 执行 discovery 命令与 skill 重复工作 | quote: "Embed 模板添加 `## Feature Paths` 后，agent 可能在 Step 1（读 task file）时执行 discovery 命令，但此时 skill 尚未被调用，agent 缺少执行上下文。" | improvement: 明确 discovery 命令的定位——信息参考而非执行指令

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

- (low) 建议：Prompt 模板中 FEATURE_SLUG 设为条件渲染或明确声明始终非空 — subjective preference
- (low) 建议：建立 embed 模板路径与 skill discovery 路径的同步验证机制 — structural suggestion, defer to Scorer cycle
- (low) 建议：明确 agent 使用 discovery 命令的时机 — structural suggestion, defer to Scorer cycle

## rubric

(all dimensions): N/A
