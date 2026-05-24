---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Iteration 0: Pre-Revision Report

## ATTACK_POINTS (Actionable Findings)

### Factual Corrections

- **[high]** 提案核心指标"27个偏差项"与表格汇总(17+13+15+5=50)矛盾 | quote: "发现 **27 个偏差项**" | improvement: 修正标题数字为50，或在Evidence表格中标注各维度去重后的独立问题数
- **[medium]** CLI交叉引用漏报了run-tests/SKILL.md第73行的`forge config get test.execution`断裂引用 | quote: "断裂 CLI 引用修复：`forge config get surface` → `forge surfaces`（4 处）、`forge test run --tags regression` → 正确命令（4 处）" | improvement: 将第6处断裂引用`forge config get test.execution`加入P0清单
- **[medium]** `forge test run --tags`的4处引用位于孤儿rules文件中（未从SKILL.md加载），运行时影响为nil，不应为P0 | quote: "`forge test run --tags regression` → 正确命令（4 处）" | improvement: 将此4处从P0降级为P1，并注明运行时影响评估
- **[medium]** 孤儿rules文件计数为15而非11 | quote: "11 个未引用 rules 文件" | improvement: 修正计数，区分真正孤儿文件（6个）和参数化引用文件（5个surface rules）
- **[high]** ARCHITECTURE.md声称eval使用"100分制"是错误的，所有rubric使用1000分制 | quote: 无（提案未提及此问题） | improvement: 在P0第2项中补充"修正ARCHITECTURE.md的100分制描述为1000分制"
- **[medium]** 提案未捕获README引用web/目录（不存在）的错误 | quote: 无（提案遗漏） | improvement: 在P0第1项的README重写清单中补充"移除web/目录引用"
- **[medium]** 提案未捕获README列出/improve-harness作为辅助技能（不存在）的错误 | quote: 无（提案遗漏） | improvement: 在P0第1项的README重写清单中补充"移除/improve-harness引用"
- **[medium]** 提案未提及ARCHITECTURE.md完全遗漏了Stop hook中的`forge feature complete --if-done`命令 | quote: 无（提案遗漏） | improvement: 在P0第2项中补充"补充forge feature complete --if-done到hooks表"
- **[medium]** ARCHITECTURE.md描述的tests/e2e/目录结构和playwright.config.ts不存在，混淆了Forge内部测试与用户项目测试 | quote: 无（提案遗漏） | improvement: 在P0第2项中补充"修正tests/目录描述，区分Forge内部Go测试与用户项目e2e测试"
- **[low]** ARCHITECTURE.md存在"forge forge task claim"重复拼写错误（第147行和第443行） | quote: 无（提案未提及） | improvement: 在P0第2项中补充"修正forge forge双重拼写错误"
- **[low]** ARCHITECTURE.md使用"all-completed Hook"命名但实际是Stop hook事件，命名不匹配 | quote: 无（提案未提及） | improvement: 在P0第2项中补充"统一all-completed Hook与Stop event的命名描述"
- **[low]** 任务类型变更不仅是数量差异（13→21），而是完整的命名体系变更（旧命名→dot-notation新命名），零重叠 | quote: "21种新命名" | improvement: 在P0第1项中明确"这不是数量修正，而是命名体系的完全替换"

### Structural/Architectural Suggestions

- **[high]** SKILL.md拆分复杂度被严重低估——eval/SKILL.md含7个rules交叉引用、4个experts引用、Mermaid流程图、Phase 0/0.5/标准流三路条件分支，且拆分会改变agent上下文加载行为（违反"仅文档变更"范围） | quote: "SKILL.md 拆分是最复杂的操作，但只需将现有内容移入 rules/ 文件" | improvement: 补充拆分风险评估：(1)拆分后agent必须显式加载额外rules文件 (2)Phase 0条件分支跨文件一致性 (3)Mermaid流程图需同步更新
- **[high]** harness rubric修复方案违反提案自身范围边界——创建rubric是内容创作（超出范围），添加异常处理是运行时行为变更（违反约束） | quote: "创建 `eval/rubrics/harness.md` 或在 SKILL.md 中添加异常处理" | improvement: 提供第三个选项：从eval SKILL.md的有效类型列表中移除harness直到rubric准备好，此选项是纯文档删除操作
- **[high]** ARCHITECTURE.md遗漏整个v3.0.0子系统（surface detection、worktree、Convention、forensic、deep-research、clean-code、extract-design-md、test-guide、learn），P0第2项仅涵盖"修正已有内容"而非"补充缺失内容" | quote: "ARCHITECTURE.md 修正：Agent 架构...Hook 系统...Eval 系统...目录路径" | improvement: 将ARCHITECTURE.md修复分为两阶段：P0修已有内容的事实错误，P1补充缺失子系统文档
- **[high]** 提案范围边界被至少4个P0/P1项违反——创建rubric文件、SKILL.md拆分、添加Load指令、删除文件均影响运行时分发内容 | quote: "仅涉及文档更新和死代码清理，不修改任何运行时代码" | improvement: 扩展范围声明为"仅涉及文档更新、死代码清理和插件内容调整"，或从P0/P1中移除运行时影响项
- **[high]** P0修复项存在隐藏依赖环——README重写依赖SKILL.md拆分和harness rubric的结果（计数可能变化），应先完成拆分再重写README | quote: P0 items 1-5 线性排列 | improvement: 调整执行顺序为：P0.4(SKILL.md拆分) → P0.5(harness rubric决策) → P0.2(ARCHITECTURE.md) → P0.3(CLI引用修复) → P0.1(README重写)
- **[medium]** 成功标准"ARCHITECTURE.md所有组件描述与代码库100%一致"不可独立验证——"组件"未定义，100%可能需要添加大量新内容超出"drift修复"范围 | quote: "ARCHITECTURE.md 所有组件描述（agents、hooks、eval、目录）与代码库 100% 一致" | improvement: 将成功标准改为"ARCHITECTURE.md中已有内容的所有事实性声明与代码库100%一致；缺失子系统文档列为P1后续任务"
- **[medium]** init-justfile的.just模板被错误分类为死代码——SKILL.md说的是"不要使用框架特定模板"的设计原则，不是"这些文件无用"的声明，模板包含功能性参考内容 | quote: "死代码清理：init-justfile 6 个 .just 模板文件（SKILL.md 明确说不使用）" | improvement: 从P1第9项中移除init-justfile .just模板的删除，改为"评估是否保留为参考实现"
- **[medium]** 跨技能路径修复不应使用本地副本方案——会导致与gen-journeys中canonical版本的静默漂移 | quote: "改为本地副本或描述性引用" | improvement: 移除"本地副本"选项，仅保留"描述性路径+上下文"方案

## BORDERLINE_FINDINGS

- forge config get在开发环境返回exit code 1，问题可能比key重命名更深层 → 建议先调查再修复，但此为运行时代码问题，属于Out of Scope

## SKIPPED_FINDINGS (Subjective Preference)

- 建议：扩展5维度为6维度（"Feature Coverage Completeness"）→ 值得考虑但非提案事实错误
- 建议：harness rubric决策应从valid types列表移除 → 部分重叠于结构建议中的第三个选项
- 建议：建立自动化验证门禁 → 合理建议但超出提案scope
- 建议：Split ARCHITECTURE.md修复为两阶段 → 已纳入结构建议
- 建议：补充forge feature complete --if-done到hooks表 → 已纳入事实修正
- 建议：明确tests/e2e/目录描述区分 → 已纳入事实修正
- 建议：修复"forge forge task claim"拼写错误 → 已纳入事实修正
- 建议：不要删除init-justfile .just模板 → 已纳入结构建议
- 建议：将forge test run --tags从P0降为P1 → 已纳入事实修正
- 建议：移除README中的improve-harness引用 → 已纳入事实修正

## Classification Audit

| Triage Layer | Count | Percentage |
|---|---|---|
| Factual correction | 12 | 40% |
| Structural/architectural suggestion | 8 | 27% |
| Subjective preference | 10 | 33% |

## Rubric

(All dimensions): N/A

## Triage Metrics

- Total findings: 30
- Accepted (factual): 12
- Accepted (structural): 8
- Deferred (borderline): 1
- Skipped (subjective): 9 (overlapping with accepted)
- Triage rate (accepted + deferred): 70%
- Accepted rate: 67%
