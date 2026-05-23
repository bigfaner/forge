---
score: 920
total: 1000
mode: B
evaluator: qa-adversary
date: 2026-05-23
---

# Final QA Evaluation Report — test-capability-v2 PRD

## Mode

**Mode B (No UI)** — `prd-ui-functions.md` 不存在，适用 Mode B 评分维度。

## Scoring Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 96 | 100 |
| 2. Flow Diagrams | 148 | 150 |
| 3. Flow Completeness | 188 | 200 |
| 4. User Stories | 198 | 200 |
| 5. Scenario Completeness | 135 | 150 |
| 6. Edge Case Coverage | 90 | 100 |
| 7. Scope Clarity | 65 | 100 |
| **Total** | **920** | **1000** |

## Dimension-by-Dimension Analysis

### 1. Background & Goals — 96/100

**Background 三要素 (28/30):**
- Reason: 三大结构性缺陷（双路径、深度不足、通用性有限）明确且具体。
- Target: 5 个升级方向清晰。
- Users: Forge 用户和 Forge 维护者两种角色。
- **-2**: Users 中"Forge 用户（项目开发者）"未进一步区分不同场景类型用户（CLI/WebUI/Mobile）的差异化痛点，而后续文档大量涉及场景差异化。

**Goals 量化 (30/30):**
- 6 个目标全部包含量化指标：测试数阈值（≥13）、风险倍率（≥1.5×）、覆盖率提升（≥20pp）、Convention 数量（≥3）、eval 评分（≥850/1000）、Mobile 接入时间（≤30 分钟）。

**逻辑一致性 (38/40):**
- Background Reason → Goals 映射完整：双路径→退休、深度不足→风险驱动+边界衍生、通用性有限→Convention 扩充。
- **-2**: Run-to-Learn 机制在 Goals 和 Flow 中占据大量篇幅，但 Background Reason 未将其作为独立缺陷明确提出，存在轻微的 Background-Goals 不对称。

### 2. Flow Diagrams — 148/150

**Mermaid 图存在 (50/50):**
- 一个约 40+ 节点的大型 flowchart，覆盖完整管线。

**主路径完整 (50/50):**
- START → PRD_CHECK → SCENE_DETECT → DETECT → GEN_JOURNEY → EVAL_J → GEN_CONTRACT → EVAL_C → GEN_SCRIPTS → R2L_CHOICE → ENV_CHECK → CONFIDENCE → RUN_TESTS → REPORT → END，完整无断裂。

**决策点 + 错误分支 (48/50):**
- 4 个菱形决策节点（PRD_CHECK, DETECT/R2L_CHOICE/FIX_DECIDE），7+ 个错误/异常分支。
- **-2**: GEN_CONTRACT 的 schema 验证失败分支在 Flow Description 文本中有描述（"验证失败则记录不符合项明细，自动重新生成一次"），但 Mermaid 图中 GEN_CONTRACT 和 EVAL_C 之间缺少该验证失败节点。

### 3. Flow Completeness — 188/200

**流程步骤完整 (65/70):**
- 四阶段 14 步骤完整覆盖，Data Flow Table 9 行数据传递路径详细。
- **-5**: (1) gen-test-scripts 验证失败的传递路径未出现在 Data Flow Table 中；(2) Run-to-Learn "覆盖率达标"退出条件中的"达标"阈值未在 Flow Description 中量化（Goals 提到"提升 ≥ 20 个百分点"，但这是从初始值算起还是绝对值？）。

**数据流文档 (68/70):**
- Data Flow Table 完整覆盖关键数据流，包含传递方式和存储位置。
- **-2**: 置信度评级的消费步骤缺少 eval 门禁反馈——eval-skipped 降级为 LOW 后该信息是否反馈到 eval_result 字段，流向不明确。

**异常处理和边界情况 (55/60):**
- eval 3 轮用尽恢复路径（3 选项）、eval-skipped 降级策略（4 步骤）、Run-to-Learn 4 种失败场景、Pipeline Exit Codes 表格均覆盖。
- **-5**: (1) Run-to-Learn "覆盖率达标"阈值缺失——若覆盖率从 40% 提升到 55%（未达 60% 但 3 轮已用完）的行为未描述；(2) 骨架测试超时保护在 Story 6 AC 中提到但 Flow Description 未指定具体超时时间。

### 4. User Stories — 198/200

**覆盖度 (50/50):**
- Story 1-3, 5-6 覆盖 Forge 用户，Story 4, 7 覆盖 Forge 维护者。

**格式正确 (50/50):**
- 全部 7 个 Story 使用 As a / I want / So that 格式，动作具体。

**AC 格式 (50/50):**
- 全部 Story 包含 Given/When/Then 格式的 AC。

**AC 可验证性 (48/50):**
- 大部分 AC 可客观验证：搜索关键词、task list 命令、评分阈值、覆盖率公式。
- **-2**: Story 5 的 "修改行数 ≤ 20%" 和 Story 6 的 "边界/异常 Outcome 占比 ≥ 30%" 标记为 `human-verified`，AC 未提供自动化验证手段或替代指标。

### 5. Scenario Completeness — 135/150

**端到端场景覆盖 (55/60):**
- 四阶段完整流程 + 5 种场景差异化策略 + Run-to-Learn 循环 + eval 门禁循环。
- **-5**: "新项目首次使用"场景的衔接不够明确——TEST_GUIDE 节点中"用户审核确认"是同步还是异步？用户是否必须立即确认才能继续管线？Flow Description 步骤 3 只说"供用户审核"，与 Mermaid 图的"用户审核确认"箭头语义不同。

**隐含假设 (32/40):**
- PRD 前置条件、场景检测规则、风险分级规则均明确。
- **-4**: (1) LLM API 可用性对 eval 迭代的影响——服务中断时 3 轮全部超时的处理未描述（与 eval-skipped 不同，那是解析失败）；(2) Mermaid 图中 R2L 在 ENV_CHECK 之前执行，但骨架测试依赖被测系统可编译，这意味着 R2L 可能先于环境检测运行；（3）R2L_CHOICE 依赖 `.forge/config.yaml` 的 `run_to_learn: true` 配置，但 Flow Description 文本未提及此配置步骤。

**业务规则一致性 (48/50):**
- BIZ-quality-gate-001: 文档明确区分两类门禁，串行执行独立判定。✓
- BIZ-error-reporting-001: Exit Codes 表格完整遵循 0/1/2 语义。✓
- BIZ-task-lifecycle-003: 提到删除 test.graduate 和 test.gen-cases。✓
- **-2**: 删除 `test.gen-cases` 和 `test.graduate` 系统类型后，未说明是否需要同步修改 BIZ-task-lifecycle-003 规则定义本身。

### 6. Edge Case Coverage — 90/100

**错误路径文档化 (38/40):**
- 覆盖 Run-to-Learn 4 种失败、eval-skipped、schema 验证失败、环境不就绪、FIX 耗尽等。
- **-2**: Convention 草稿被用户拒绝 2 次后仍不满意的场景——Story 5 AC 提到"最多重试 2 次"，但 Flow Description 未描述 2 次用尽后用户仍不满意的行为（管线暂停？使用默认 Convention？）。

**边界条件 (30/35):**
- 覆盖测试数范围、迭代次数上限、修复次数上限。
- **-5**: (1) 场景检测匹配到多个类型时用户选择与实际信号冲突的后果未描述；(2) Fact Table 覆盖率边界值（= 0% 或 = 59%）的行为未描述；(3) 超大项目（数千测试文件）的性能边界未提及。

**失败恢复 (22/25):**
- PAUSE_J/PAUSE_C 三选项恢复、eval-skipped 手动清除、Run-to-Learn 降级、FIX 回退上游。
- **-3**: "修改后重跑 eval"不计入自动迭代轮次，但未说明是否有总次数限制——理论上用户可无限重跑。

### 7. Scope Clarity — 65/100

**In-scope 具体可交付物 (34/35):**
- 每项都是具体的功能模块或文件，可追溯。
- **-1**: "质量门禁更新以反映新管线"作为 in-scope 条目不够原子化，实际在后续详细描述中拆分为了多阶段门禁。

**Out-of-scope 延迟项 (28/30):**
- 10 项明确列出，附带解释理由。
- **-2**: "执行环境自动准备与配置（仅做就绪检测）"与 in-scope 的 ENV_CHECK 边界模糊——"就绪检测"vs"自动准备"的区分依赖主观判断。

**Scope 与 Specs/User Stories 一致 (3/35 after cross-section deduction):**
- **Cross-section inconsistency (-30)**: Story 7 描述了一个完整的"可扩展场景类型系统"——通过向 `scenarios/` 目录添加配置文件来新增场景类型（Desktop/Electron、嵌入式、gRPC）。但 In Scope 列表中没有对应的"场景类型配置系统"或"可扩展场景类型"条目。虽然"场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为"部分覆盖了场景支持，但 Story 7 描述的**可扩展性机制**（新场景类型无需修改管线代码即可接入）是一个独立的能力维度，在 Scope 中既未列入 in-scope 也未列入 out-of-scope。这导致下游实现者无法从 Scope 节判断 Story 7 是否在本版本范围内。
- **Base score: 33/35, after -30 deduction: 3/35**

## Adversarial Attacks (Blindspot Hunt)

1. **[Scope Clarity — Cross-section inconsistency]** Story 7 描述"通过添加一个场景类型配置文件...让新场景类型无缝接入测试管线"并提到 `scenarios/` 目录和回归验证标准，但 In Scope 列表中只有"场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为"，**没有**"可扩展场景类型系统"条目。Out of Scope 中也未列出此项。引用原文：Story 7 AC "将配置文件放入 `scenarios/` 目录...无需修改任何管线技能代码" vs In Scope "场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为"。**改进要求**: 在 In Scope 中新增"可扩展场景类型系统"条目，或将 Story 7 移至 Out of Scope 并标注为未来迭代。

2. **[Scenario Completeness — R2L 执行顺序与环境检测时序]** Mermaid 图显示 R2L_CHOICE 在 ENV_CHECK 之前（步骤 9 vs 步骤 10），但 Run-to-Learn 需要实际运行骨架测试（依赖被测系统可编译可执行）。引用原文 Flow Description 步骤 9: "运行骨架测试捕获实际输出" vs 步骤 10: "场景特定环境就绪检测：验证执行环境是否准备好"。**改进要求**: 明确说明 R2L 是否需要环境已就绪作为前置条件，如果是则调整 Mermaid 图中的执行顺序或增加 ENV_CHECK → R2L 的路径。

3. **[Flow Completeness — Run-to-Learn 覆盖率达标阈值缺失]** Flow Description 步骤 9 提到 "≤ 3 轮或覆盖率达标" 作为退出条件，但"覆盖率达标"的绝对阈值未在 Flow 中定义。引用原文: "经过 ≤ 3 轮迭代后，Fact Table 覆盖率从初始值（gen-contracts 静态侦察结果）提升 ≥ 20 个百分点"——这是 Goals 中的描述，但 Goals 说的是"提升 ≥ 20pp"而非"达到某个绝对值"。如果初始覆盖率是 50%，提升 20pp 后是 70%，这算"达标"吗？如果初始是 80%，提升 20pp 后是 100%（不可能），怎么算？**改进要求**: 在 Flow Description 中明确"覆盖率达标"的判定标准（是绝对值还是相对提升量）。

4. **[Edge Case Coverage — Convention 草稿拒绝耗尽]** Story 5 AC: "最多重试 2 次"。Flow Description TEST_GUIDE 节点: "用户拒绝, 重试 ≤ 2"。但 2 次重试用尽后用户仍不满意的行为未描述。引用原文: "test-guide 基于用户反馈重新生成草稿（保留用户认可的部分，仅修正被指出的错误），最多重试 2 次"。**改进要求**: 补充 2 次重试用尽后的处理策略（暂停管线供用户手动编写 Convention？使用最后一次草稿并标记为 unverified？）。

5. **[Background & Goals — Users 角色粒度不足]** Background Users 只列了两种角色，但文档核心特性之一是"场景差异化"（CLI/TUI/WebUI/Mobile/API），不同场景用户面临截然不同的接入路径和痛点。引用原文: "Forge 用户（项目开发者）：使用 Forge 的 /quick 或 full pipeline 创建功能，期望管线自动生成高质量测试"——这个描述对 Mobile 项目开发者和 CLI 项目开发者一视同仁，但 Mobile 是"尽力而为"策略。**改进要求**: 在 Users 中至少区分"核心场景用户"和"尽力而为场景用户"（或等效区分），使 Background 与后续场景差异化策略形成对照。

6. **[Edge Case Coverage — 用户手动修改后无限重跑 eval]** Flow Description PAUSE_J/PAUSE_C 恢复路径选项 c: "用户手动修改 Journey/Contract 文档后，重新进入对应 eval 步骤（不计入自动迭代轮次）"。引用原文: "不计入自动迭代轮次"——这意味着用户可以无限次修改后重跑 eval，没有总次数上限。**改进要求**: 增加手动重跑的总次数上限（如 ≤ 5 次），或在 Scope 中说明"无限重跑是预期行为"。

7. **[Scenario Completeness — LLM API 可用性假设]** eval 门禁依赖 LLM API 调用进行评分，但文档假设 LLM 始终可用且响应及时。eval-skipped 处理的是"LLM 输出无法解析"（格式问题），但不是"LLM 服务不可用"（连接问题）。引用原文: "eval 评分因 LLM 输出无法解析而失败，记录错误日志并重试评分一次"。**改进要求**: 区分"LLM 输出格式异常"和"LLM 服务不可用"两种失败场景，后者应有独立的超时和重试策略。
