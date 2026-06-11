# Merged Attack Points — Iteration 1

**PM Score**: 666/1000
**QA Score**: 913/1000
**Average**: 790/1000

## Merged Attacks (16 unique)

1. [Background & Goals]: Mobile 接入成本降低与 Background 三大缺陷无逻辑联系 — "降低 Mobile 接入成本" 作为 Goal 无对应 Background 问题陈述 — 在 Background Why 中补充 Mobile 问题或标注为独立演进方向

2. [Flow Diagrams]: FIX_DECIDE 缺少 Contract 语义错误回退路径 — Flow Description 步骤 14 文字描述了两种回退但图中只有一条回到 GEN_SCRIPTS — 在 Mermaid 图中为 FIX_DECIDE 添加区分"脚本问题"和"Contract 语义错误"的决策分支

3. [Flow Diagrams]: Run-to-Learn 骨架测试执行失败分支在图中缺失 — R2L 节点只有成功路径到 ENV_CHECK — 为 R2L 节点添加骨架测试失败的降级路径（使用已有静态信息继续）

4. [Flow Diagrams]: eval LLM 解析失败分支在图中缺失 — EVAL_J/C 只有通过/不通过两个出口 — 添加 eval 解析失败的第三分支

5. [Flow Completeness]: 数据流中 Convention、Fact Table、置信度评级的数据传递路径未文档化 — 多个步骤隐式消费这些数据但传递机制不明 — 添加数据流表或在每个步骤明确输入/输出

6. [Flow Completeness]: test-guide 拒绝重试路径在 Flow 中缺失 — Story 5 描述"基于用户反馈重新生成草稿，最多重试 2 次"但图中只有"用户审核确认"单一出口 — 添加 TEST_GUIDE 拒绝重试回路

7. [User Stories]: Story 2 AC 的 1.5x 比较在只有高风险 Journey 时不可验证 — "同一功能的 Journey 比较高风险 vs 低风险 variant" — 添加退化为绝对值时的备选 AC

8. [User Stories]: Story 5 缺少 When 条件 — "Given 内置 Convention 库 / Then 包含 pytest、JUnit、Rust" 无 When — 补充 When 条件或合并到前一个 GWT 块

9. [Scenario Completeness]: 风险驱动密度表下界数据与 Goal 1.5x 指标数学矛盾 — High 下界 10 vs Low 上界 8，10/8 = 1.25 < 1.5 — 调整密度表下界使 High 下界 ≥ Low 上界 × 1.5

10. [Scenario Completeness]: 质量门禁更新无细节 — In Scope "质量门禁更新以反映新管线"与 BIZ-quality-gate-001 多阶段管线如何融合未说明 — 补充 quality-gate 与新管线集成的方式

11. [Edge Case Coverage]: eval-skipped 降级后果和恢复路径不完整 — "由用户手动审核"但审核内容、审核后操作均未定义 — 为 eval-skipped 定义自动降级策略（如自动标记为 LOW 置信度）

12. [Edge Case Coverage]: gen-journeys 提取失败（PRD 无用户故事）未覆盖 — "项目必须已有 PRD"但未定义 PRD 质量要求 — 添加 PRD 质量前置检查

13. [Scope Clarity]: "LLM prompt 增强策略"作为 in-scope 是实现手段而非功能 — 替换为功能描述如"边界/异常场景自动衍生引擎"

14. [blindspot]: Story 4 AC "评分 ≥ 850/1000" vs Other Notes 维度最低阈值 — 总分 ≥ 850 但某维度不达标的情况在 AC 中不会被捕获 — 在 Story 4 AC 中引用维度最低阈值

15. [blindspot]: "不得降低断言严格度"不可验证 — 步骤 14 无检测机制描述 — 定义可机器检查的约束或标注为 human-verified

16. [blindspot]: 风险密度数字范围重叠 — High 10-20, Medium 7-13, Low 4-8 — 10-13 和 7-8 区域重叠 — 收紧范围消除重叠区间
