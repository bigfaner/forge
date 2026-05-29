---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

## ATTACK_POINTS

- **[high]** IsTestableType()修改范围被低估，影响stage-gate生成 | quote: "forge-cli `build.go`：`IsTestableType()` 区分行为变更 vs 行为保持" — improvement: 明确讨论 refactor/cleanup 下 stage-gate 的处理策略
- **[high]** refactor Breakdown管道依赖链断裂，lastRunID为空导致validate-code无依赖 | quote: "下游任务（validate-code, clean-code, consolidate-specs）的依赖链需重新挂载" — improvement: 给出具体的接线逻辑——挂载到最后一个 business task、gate 还是 summary
- **[high]** coding.fix类型被忽略，可能导致intent与任务类型判断不一致 | quote: "只区分了 new-feature、refactor、cleanup 三种 intent" — improvement: 明确定义 coding.fix 在各 intent 下的管道行为
- **[high]** 混合类型特征下refactor intent会错误跳过新增功能的测试覆盖 | quote: "混合 intent 支持——不支持，按主要意图归类" — improvement: 定义"主要意图"的判断标准，或改用 per-task 级别的 testable 判断
- **[medium]** Quick+refactor管道中T-clean-code和T-validate-code依赖来源不明确 | quote: "proposal -> quick-tasks -> quality-gate -> done" — improvement: 补充 Quick+refactor/cleanup 的完整任务链和依赖接线
- **[medium]** Intent字段注入时机与detectMode和setFeatureMetadata存在依赖耦合 | quote: "`BuildIndex()` 在 `detectMode()` 之后读取 `proposal.md` frontmatter 的 `intent` 字段" — improvement: 将 intent 作为 BuildIndexOpts 字段传入，或在 Feasibility 中明确时序
- **[medium]** spec-only PRD格式可能与下游tech-design skill不兼容 | quote: "write-prd 分支确保 spec 格式包含 tech-design 需要的字段" — improvement: 明确 spec-only PRD 包含哪些字段以满足 tech-design 输入需求
- **[medium]** Breakdown+cleanup组合缺少防止用户触发的代码级别防护 | quote: "cleanup 不走 Breakdown 模式" — improvement: 添加代码级别防护或在 Key Risks 中讨论此隐式无效路径
- **[medium]** quality-gate hook作为唯一验证手段存在覆盖盲区 | quote: "已有 quality-gate hook（compile + fmt + lint + test）" — improvement: 在 Key Risks 中承认此覆盖盲区，说明接受此风险的权衡

## BORDERLINE_FINDINGS

- (none)

## SKIPPED_FINDINGS

- 建议：将intent作为BuildIndexOpts字段传入 → classified as structural suggestion, defer to Scorer cycle for design evaluation
- 建议：为refactor/cleanup下游任务提供显式接线函数 → classified as implementation detail, already covered by ATTACK_POINT #2
- 建议：per-task级别testable判断 → classified as significant design change beyond pre-revision scope, defer to Scorer cycle
- 建议：为Breakdown+cleanup添加代码级别防护 → classified as implementation detail, already covered by ATTACK_POINT #8
- 建议：定义coding.fix的intent映射 → already covered by ATTACK_POINT #3

## Classification Audit

Total findings by triage layer:
- Factual correction: 4 (findings 2, 3, 4, 7)
- Structural suggestion: 5 (findings 1, 5, 6, 8, 9)
- Subjective preference: 0
- Implementation suggestions: 5 (deferred — overlap with factual/structural findings)
