# Merged Attack Points — Iteration 3

**PM Score**: 864/1000
**QA Score**: 868/1000
**Average**: 866/1000

## Score Progression

| Iteration | PM | QA | Average | Target |
|-----------|------|------|---------|--------|
| 1 | 578 | 698 | 638 | 900 |
| 2 | 815 | 835 | 825 | 900 |
| 3 | 864 | 868 | 866 | 900 |

## Merged Attacks (12 unique)

1. [Background & Goals]: 风险比值无绝对下限 — "高风险旅程测试数 ≥ 低风险旅程 × 1.5"——增加绝对下限（如高风险 Journey 平均 ≥ 8 个测试），避免低基数时差异微乎其微

2. [Flow Completeness]: eval rubric 评分维度框架只在 Scope 中列出，Flow Description 未覆盖——6 维度表格是 eval 核心交付物但 Flow 步骤 5/7 只说"评估质量"——在 Flow Description 中引用评分维度或在 eval 步骤处注明评分框架细节

3. [Flow Completeness]: FIX_DECIDE 回路粒度不对且无安全边界——回到 GEN_SCRIPTS 无法修复 Contract 语义错误，且无触发条件分类（脚本问题 vs 系统问题）和安全约束（如不降低断言严格度）——需区分失败类型并设计不同回退层级

4. [Scenario Completeness]: 检测信号表只覆盖 Go 和 Node.js——7 行信号全部基于 Go/Node.js，但 Convention 扩充计划包含 pytest(Python)/JUnit(Java)/Rust——为 Python/Java/Rust 项目补充信号检测规则

5. [Scenario Completeness]: PRD 存在性前置条件未声明——步骤 4"从 PRD 用户故事提取"但流程从"用户运行测试生成"开始，跳过 PRD 检查——在流程图 SCENE_DETECT 之前增加 PRD 存在性检查节点

6. [User Stories]: Story 3 统计分母未定义——"Contract 测试占比 ≥ 80%"的分母定义不明确——明确分母定义（Outcome 数量？测试函数数量？）并与 Per-Scenario Strategy 表格的"AI 优先侧重"概念对齐

7. [User Stories]: Story 1 否定验证 AC 不可穷举——"没有任何 gen-test-cases 相关的技能、命令、或 rubric 文件残留"——需提供明确删除清单替代否定表述

8. [Edge Case Coverage]: eval 评分系统自身失败未覆盖——eval-journey/eval-contract 依赖 LLM 评分，LLM 返回无法解析时管线行为未定义——需定义评分失败的 fallback 行为（如重试、跳过门禁、暂停）

9. [Edge Case Coverage]: FIX_DECIDE 自动修复无安全边界——需定义触发条件分类（脚本问题 vs 系统问题）和安全约束（如不降低断言严格度）

10. [blindspot]: Eval Rubric 维度级阈值缺失，850/1000 可被严重缺陷维度绕过——需定义每个维度的最低通过阈值

11. [blindspot]: 管线终止点退出码未与 BIZ-error-reporting-001 对齐——PAUSE_J/PAUSE_C、FIX_DECIDE -> END 节点需为每个非正常终止点定义退出码

12. [blindspot]: gold standard 校准方法论缺失——"基于 gold standard 评分集校准"需定义集合大小、场景覆盖、标注者、校准方法

## Outcome

**FAIL** — 平均分 866/1000，未达目标 900/1000。文档已回滚到迭代 1 修正前的备份状态。
