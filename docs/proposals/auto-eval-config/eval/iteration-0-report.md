---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision (Freeform Findings)

## ATTACK_POINTS

- **[high]** 嵌套命名空间破坏现有 YAML 解析和 dot-notation 路由逻辑 | quote: "当前 AutoConfig 是纯 flat 结构…引入 auto.eval.proposal 意味着 auto.eval 是一个中间结构体而非 ModeToggle，现有的 autoModeField(\"eval\") 返回 nil" | improvement: 重新评估使用扁平命名空间（auto.evalProposal 等）替代嵌套结构体，或详细说明嵌套解析的改造方案
- **[high]** Go struct 嵌套与 YAML 序列化的零值陷阱，无法区分显式 false 与未配置 | quote: "ModeToggle 的零值是 {false, false}。当 uiDesign 的默认值恰好等于零值时，AutoConfig.WithDefaults() 和 applyDefaults() 无法区分\"用户显式设置为 false\"和\"字段未配置\"" | improvement: 说明 raw 追踪逻辑如何扩展到嵌套结构体以支持差异化默认值
- **[medium]** quick/full 区分在 eval 场景缺乏语义基础，skill 无法判断当前模式 | quote: "eval 的触发点是 skill 内部（brainstorm、write-prd、tech-design、ui-design），这些 skill 本身并不区分 quick/full 流水线。proposal 没有说明 skill 如何判断当前处于 quick 还是 full 模式。" | improvement: 补充 skill 如何感知当前 quick/full 上下文的机制说明
- **[high]** ui-design 从无条件自动运行降级为默认询问，构成功能回退 | quote: "将 \"无条件自动\" 降级为 \"默认询问\" 是一个功能回退。ui-design 是唯一一个当前无条件自动评估的 skill，这个行为可能恰恰是因为 ui 设计质量更需要自动化验证而有意为之。" | improvement: 为 ui-design 默认行为变更提供迁移路径，或将 uiDesign 默认设为 true/true
- **[medium]** skill 中 config check 逻辑的实现机制未说明 | quote: "当前 skill 是 markdown 文件（SKILL.md），它们通过自然语言指令指导 Claude Code agent 的行为。提案没有说明 config check 的实现机制。" | improvement: 明确 skill 中 config check 的实现方式（如 forge config get 调用 + 条件分支）
- **[medium]** 四个 skill 的 config check 一致性在 markdown 语境下无运行时保证 | quote: "在四种不同的控制流中插入统一的 config check 模板，需要修改每个 skill 的特定步骤，而 \"统一模板\" 在 skill markdown 语境下只是一个 copy-paste 的文本块，没有运行时保证。" | improvement: 承认一致性问题并在风险/缓解措施中说明
- **[medium]** forge config get auto.eval 无子字段时的行为未定义 | quote: "如果 auto.eval 是一个嵌套结构体，forge config get auto.eval 应该返回什么？" | improvement: 明确定义 auto.eval 路径的行为边界
- **[medium]** 1-2 小时估算未包含嵌套配置系统的工程复杂度 | quote: "这个估算假设改动是纯增量的。但 parseAutoRaw 需要支持嵌套解析…这些不是四个独立的简单改动，而是对配置系统核心路由逻辑的修改。" | improvement: 修订时间估算为更现实的数值（如 3-5 小时）

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

- (subjective-preference) 建议在 proposal 中显式列出 AutoConfig Go struct 的变更伪代码 — 属于实现细节，不在 proposal 层面强制要求
- (subjective-preference) 建议补充 skill 中 config check 的具体实现规格 — 部分已由 ATTACK_POINTS 中的相关项覆盖
- (subjective-preference) 建议增加 forge config get auto.eval 无子字段时的行为定义 — 已作为 ATTACK_POINT 列出
