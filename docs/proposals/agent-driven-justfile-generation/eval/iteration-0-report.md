---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Iteration 0: Pre-Revision Report

## ATTACK_POINTS

1. **[high]** 混合项目的端口冲突处理在提案中完全缺失 | quote: "提案没有说明 `server-lifecycle.md` 是否会包含这种 per-service 的端口感知启动模式" | improvement: 在 Scope > In Scope 中补充对 multi-service lifecycle 的覆盖说明，或在 server-lifecycle.md 范围中明确包含 per-service PID 文件和端口感知启动模式
   - Triage: factual correction (verifiable gap — section is absent from proposal)

2. **[high]** 提案声称可一次 session 完成，但实际涉及 5 个 surface rule 重写+新 rule+SKILL.md 重写+不可逆架构迁移，scope 被低估 | quote: "改动集中在 `init-justfile` skill 目录内...预计单个 skill 改动，可在一次 session 内完成" | improvement: 调整 Resource & Timeline 评估，承认实际影响链，添加不可逆迁移的回退策略
   - Triage: factual correction (verifiable — actual deliverables exceed claimed scope)

3. **[high]** surface rule 文件同时服务 init-justfile 和 run-tests 两个消费者，新增 "Recipe Generation Requirements" section 与现有 Recipe Invocation Contract 缺乏一致性验证 | quote: "如果新增的 'Recipe Generation Requirements' section 与 Recipe Invocation Contract 之间存在不一致...run-tests 会在运行时遇到缺失 recipe 的失败" | improvement: 在 Scope > In Scope 中添加 post-generation 一致性验证步骤
   - Triage: structural/architectural — identifies verifiable internal inconsistency between dual-consumer contract (partially accepted: add post-generation check mention)

4. **[high]** 提案对一致性的定义低估了命令差异对下游消费者的连锁影响——"细微差异"可能导致进程未正确后台化、PID 文件未写入等 | quote: "一致性：相同项目多次运行生成的 justfile 结构一致（recipe 名称、分组、边界标记不变；具体命令可能因 LLM 变化而有细微差异）" | improvement: 明确定义"一致"为结构级（recipe 名称、分组、边界标记、退出码语义），在 NFR 中澄清命令体的允许差异范围
   - Triage: structural/architectural — identifies internal contradiction between "consistent" and "LLM variation" claims

5. **[high]** verification step 的覆盖范围不足以捕获 server lifecycle 的边界情况 | quote: "真正的边界情况（PID 文件残留、进程被外部 kill 后 PID 被回收、Windows 上的 `\r` 污染）在 verification 中完全无法覆盖" | improvement: 在 Key Risks 的 mitigation 中承认 verification step 对 server lifecycle 边界条件的覆盖限制，补充 server-lifecycle.md 提供可执行参考代码（带插槽）作为额外缓解
   - Triage: structural/architectural — identifies verifiable gap in mitigation completeness

6. **[high]** 移除 project-detection.md 后混合语言多 surface 项目的语言检测信号源消失，提案未定义替代检测策略 | quote: "提案移除了项目类型分类和检测算法，但没有说明 agent 在多 surface 项目中如何确定：每个 surface 使用什么语言？" | improvement: 在 Proposed Solution 中补充 agent 的语言/框架检测策略说明（如：扫描 marker files → 确定语言 → 查 Convention），或在 Scope 中保留 marker file → 语言映射作为确定性参考
   - Triage: structural/architectural — identifies verifiable gap in solution (partially accepted: add detection strategy description)

## BORDERLINE_FINDINGS

7. **[high]** server lifecycle bash 提取为描述性 rule 文件后，LLM 生成的 PID 追踪代码可能丢失边界条件处理，这是确定性代码→概率性生成的本质降级 | quote: "这不是'轻微不确定性'——这是将确定性代码替换为概率性生成的本质降级" | improvement: 将 server-lifecycle.md 设计为带插槽的可执行代码模板，而非纯描述性文档
   - Triage: borderline — the concern is legitimate (deterministic→probabilistic transition risk) but the proposed solution (code template with slots) conflicts with the proposal's core thesis (zero-template). Deferring to Scorer cycle for resolution.

## SKIPPED_FINDINGS (subjective preference — not actionable for pre-revision)

8. **[low]** 建议将 server-lifecycle.md 设计为带插槽的可执行代码模板而非描述性文档 | Triage: subjective preference — implementation approach suggestion
9. **[low]** 建议保留 project-detection.md 的检测信号映射并重构为 rules/language-detection.md | Triage: subjective preference — implementation approach suggestion
10. **[low]** 建议为 surface rule 文件引入结构化验证约束确保双消费者一致性 | Triage: subjective preference — implementation approach suggestion
11. **[low]** 建议在混合项目 surface rule 文件中增加 multi-service lifecycle 专门指导 | Triage: subjective preference — implementation approach suggestion
12. **[low]** 建议增加黄金回归测试验证 agent 生成 justfile 与当前模板输出在结构关键点上匹配 | Triage: subjective preference — implementation approach suggestion
13. **[low]** 建议增加回退策略显式文档 | Triage: subjective preference — implementation approach suggestion

## Classification Audit

- 6 findings classified as actionable (2 factual correction + 4 structural with verifiable inconsistency)
- 1 finding classified as borderline (legitimate concern, proposed solution conflicts with proposal thesis)
- 6 findings classified as subjective preference (implementation approach suggestions, deferred to Scorer cycle)
- Triage rate: 100% (all findings classified)
- Accepted + partially accepted: 6/13 (46%)

## Rubric

All dimensions: N/A (freeform findings, no rubric scoring)
