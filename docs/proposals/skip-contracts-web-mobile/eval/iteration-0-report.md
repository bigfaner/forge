---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

## ATTACK_POINTS

1. **[HIGH]** Pipeline 层与 Skill 层决策优先级未声明 | quote: "两层修复" | improvement: 明确 Skill 层为权威决策者，Pipeline skip 仅为优化
2. **[HIGH]** 覆盖率自检假成功：错误类型测试未覆盖 | quote: "覆盖率自检硬失败兜底" | improvement: 自检不仅检查存在性，还需验证测试类型与 surface-type 匹配
3. **[MEDIUM]** CondHasProtocolSurfaceTask 对 surface-type 缺失值无处理策略 | quote: "检查 feature 所有业务任务的 surface-type 字段" | improvement: 声明空/缺失字段视为"可能有协议级"（保守策略），只有明确的 web/mobile 才触发跳过
4. **[MEDIUM]** 覆盖率自检需按 surface type 细分计数 | quote: "count(journeys) == count(test-script-sets)" | improvement: 自检改为按 surface type 分别验证每个 journey 的测试覆盖
5. **[MEDIUM]** gen-test-scripts 前置条件修改细节缺失 | quote: "跳过 contract 前置检查" | improvement: 在 Scope 中补充前置条件路由的具体修改点和 Step 2/Step 2.5 的替代数据源
6. **[MEDIUM]** ResolveUpstream 回退行为未验证 | quote: "依赖链调整：gen-contracts 跳过时，gen-scripts 依赖 gen-journeys" | improvement: 声明依赖链回退的验证点

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

- TUI 归类边界条件质疑：lesson 文档已有清晰的执行模型分析矩阵论证 TUI 保持 Contract 路径的合理性，此为设计决策而非缺陷

## Rubric

(all dimensions): N/A
