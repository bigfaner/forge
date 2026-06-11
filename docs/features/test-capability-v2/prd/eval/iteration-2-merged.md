# Merged Attack Points — Iteration 2

**PM Score**: 955/1000
**QA Score**: 773/1000
**Average**: 864/1000

## Merged Attacks (14 unique)

1. [Background & Goals]: Goal "降低 Mobile 接入成本"的 metric 是交付物列表而非量化指标 — 补充对照度量

2. [Flow Completeness]: PRD 前置检查错误路径未纳入 Pipeline Exit Codes 表和 Mermaid 图 — "若 PRD 不存在或质量前置检查未通过，管线在步骤 1 报错" — 将此错误路径补充到 Exit Codes 表和流程图

3. [Flow Completeness]: Data Flow Table 缺少 eval → revise 数据传递路径和 GEN_SCRIPTS → RUN_TESTS 测试代码传递 — 添加 eval 结果到 revise 的数据格式行和测试代码传递行

4. [Flow Completeness]: PAUSE_J/PAUSE_C 后的恢复路径未定义 — "由用户决定后续操作" — 定义 2-3 种恢复路径（跳过门禁继续、放弃管线、修改后重跑）

5. [Scenario Completeness]: 质量门禁与 BIZ-quality-gate-001 集成关系未定义 — "将现有单一门禁替换为多阶段门禁" — 需说明新 eval 门禁与已有 compile→fmt→lint→unit/integration→e2e 管线的关系

6. [Scenario Completeness]: eval 退出码归因矛盾 — Exit Codes 表说"rubric 配置错误"但 Flow Description 说"LLM 输出无法解析" — 统一归因描述

7. [User Stories]: Story 7 "不受影响"的验证标准模糊 — "测试生成结果不受影响（回归验证）" — 定义可接受的差异范围（如：测试函数数量变化 ≤ 5%，无新增编译错误）

8. [Edge Case Coverage]: gen-test-scripts 输出无质量验证 — GEN_SCRIPTS 直接进入下游 — 添加生成后语法/可执行性检查步骤

9. [Edge Case Coverage]: gen-contracts 合约 schema 验证失败未覆盖 — 定义 schema 验证失败的处理

10. [blindspot]: eval-contract "事实依据"维度与衍生引擎架构矛盾 — LLM 衍生的边界 Outcome 不在 Fact Table 中 — 在 eval rubric 中区分"事实依据"和"合理推理"两类声明

11. [blindspot]: Fact Table runtime 覆盖 static 策略可能与覆盖率公式冲突 — runtime 替换后 confidence 可能不是 confirmed — 明确替换时保留 static 作为 fallback

12. [blindspot]: Delivery Phasing 阶段 1 门禁未量化 — "2+ 个已有项目跑完整管线无报错" — 定义项目特征和通过标准

13. [blindspot]: Story 7 缺少新场景类型评测门禁要求 — eval rubric 需动态适配新场景类型的 required_outcomes — 增加 AC

14. [blindspot]: 场景检测规则 CLI/TUI 仅覆盖 Go 项目 — 补充非 Go 语言的 CLI/TUI 检测规则或明确声明覆盖范围
