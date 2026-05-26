# Eval Report — Iteration 2 (PM)

**Score**: 872/1000
**Iteration**: 2
**Role**: Senior PM (Adversary)
**Date**: 2026-05-25

## Dimension Breakdown

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 89 | 100 |
| Flow Diagrams | 130 | 150 |
| Flow Completeness | 180 | 200 |
| User Stories | 178 | 200 |
| Scenario Completeness | 127 | 150 |
| Edge Case Coverage | 79 | 100 |
| Scope Clarity | 89 | 100 |

## Iteration 1 Issue Resolution

| # | Attack Point | Status | Evidence |
|---|-------------|--------|----------|
| 1 | 跨组件数据流文档缺失 | RESOLVED | 新增"跨组件数据流"7步传递链表格 + Fallback 链描述 |
| 2 | Exit code 与 HARD-GATE 语义不一致 | RESOLVED | 新增"Exit Code 语义与 HARD-GATE 定义"节，对齐 BIZ-error-reporting-001/002 |
| 3 | 错误路径覆盖严重不足 | RESOLVED | 新增"Error Handling Paths"表格，覆盖5种错误场景 |
| 4 | run-tests fallback 缺失 | RESOLVED | Flow 步骤1定义了两来源均失败时的行为（exit 2 + 恢复提示） |
| 5 | Surface-key 自定义 vs 5 固定类型张力 | RESOLVED | Goals 表格明确"用户自定义的是 key 名称而非 type 枚举" |
| 6 | forge surfaces CLI 前置依赖未声明 | RESOLVED | In Scope 中明确"当前状态：需新建" |
| 7 | 混合项目缺少端到端流程 | RESOLVED | 新增"混合项目生成与编排流程"完整节 + 代码示例 |
| 8 | Probe 重试参数分散 | RESOLVED | 新增"Probe 重试规格"单一表格定义 |

**净改进**: +44 分（从 828 提升至 872）

## Detailed Scoring

### 1. Background & Goals — 89/100

**Background 三要素 (25/30):**
Why/What/Who 三要素齐全。扣分：Who 部分"Forge 用户（项目开发者）"的定义偏模糊，描述的是使用场景而非用户画像——缺少用户的技术背景假设、使用频率、典型项目规模等。

**Goals 量化 (28/30):**
5 项 Goals 均有量化指标（层数、类型数、组件数、字段数）。扣分："7+ 组件"中的"+"不够精确。

**逻辑一致性 (36/40):**
Background 和 Goals 逻辑一致。扣分："75% 的实际示例已通过 just 命令调用"——此数据支撑消除 test.execution 的决策，但来源不可独立验证。

### 2. Flow Diagrams — 130/150

**Mermaid 图存在 (45/50):**
涵盖 init-justfile 和 run-tests 两个主要流程。扣分：混合项目流程没有独立的 Mermaid 图，只有文字描述 + 代码示例。

**主路径完整 (45/50):**
两个主流程 start→end 完整。扣分：cli/tui 分支的 Test2→RunEnd 没有错误处理分支（exit 1/2），与 Flow Description "每步检查退出码"不一致。

**决策点 + 错误分支 (40/50):**
含 4 个决策点、probe exit 分支。扣分：(1) run-tests "获取 surface 信息失败"分支在文字中描述但 Mermaid 图未体现；(2) 加载规则文件失败分支缺失；(3) init-justfile 中 just 版本检查步骤在图中缺失。

### 3. Flow Completeness — 180/200

**流程步骤 (63/70):**
init-justfile 5步、run-tests 4步、混合项目 4步均完整。Surface 编排模式表覆盖 5 种类型。扣分：`forge surfaces` CLI 的内部行为（longest-prefix-match 算法、性能特征）无接口契约描述。

**数据流文档 (65/70):**
跨组件数据流表格覆盖 7 步传递链，格式清晰。扣分：步骤 7 "surface-key 列表及其类型"缺少具体 JSON schema；步骤 1 "文件读取 + CLI 查询"不够精确。

**异常处理 (52/60):**
Error Handling Paths 表覆盖 5 种场景，Exit Code 语义表完整，Probe 重试规格详细。扣分：(1) teardown 本身失败时的行为未定义；(2) dev server 在 probe 成功后、test 执行期间崩溃的行为未定义；(3) `.forge/test-state.json` 恢复流程只有一句话，缺少具体步骤。

### 4. User Stories — 178/200

**覆盖 (45/50):**
覆盖 Forge 用户（2 stories）和 Forge 插件开发者（2 stories）。扣分：缺少混合项目的独立用户故事——混合项目在 Flow Description 中有详尽描述，但无对应 story。

**格式 (48/50):**
4 个 Story 均使用 As a / I want / So that 格式。规范。

**AC 格式 (45/50):**
每个 Story 有 3-4 条 Given/When/Then AC。扣分：Story 2 第3条 AC "probe 失败后禁止在同一编排周期内重试 probe 或重启 dev（HARD-GATE）"不是 Given/When/Then 格式。

**AC 可验证性 (40/50):**
大部分 AC 可通过检查生成文件或 exit code 验证。扣分：(1) Story 4 "旧任务的 scope 字段通过 GetSurfaceKey() 兼容访问"——"兼容访问"含义模糊（返回旧值？映射到新字段？）；(2) Story 4 "forge task add 从源任务继承"缺少"无源任务时"的边界条件；(3) Story 2 未覆盖 teardown 执行失败场景。

### 5. Scenario Completeness — 127/150

**端到端覆盖 (52/60):**
覆盖 5 种 surface 的 init 和 run-tests、混合项目生成与编排。扣分：缺少"从零配置 surfaces 到运行第一次测试"的完整用户旅程走查——各片段分散在不同节中。

**隐含假设 (32/40):**
主要假设已暴露（CLI 前置、版本要求、零回归）。扣分：(1) "编排级配方"的定义未明确——Flow 中"Surface 规则覆盖语言模板的编排级配方"但哪些配方属于"编排级"未列清单；(2) 各 surface 规则文件独立更新的兼容性假设未暴露。

**业务规则一致性 (43/50):**
Exit code 和 error message 与 BIZ 规范对齐。扣分：Mermaid 图中 cli/tui 分支的 test 失败没有错误处理，但 Flow Description 步骤 4 说"每步检查退出码"——存在不一致。

### 6. Edge Case Coverage — 79/100

**错误路径 (33/40):**
Error Handling Paths 表 + Exit Code 表 + Probe 重试规格。扣分：(1) init-justfile 某个 surface 规则文件损坏/格式错误的处理缺失；(2) teardown 配方执行失败的处理未定义。

**边界条件 (28/35):**
覆盖 just 版本、无 surface 项目、双源失败、混合项目、旧任务兼容。扣分：(1) surface-key 名称含特殊字符（空格、非 ASCII）时影响 just 配方名（`dev-<surface-key>`）；(2) longest-prefix-match 等长匹配冲突未定义。

**失败恢复 (18/25):**
提及 test-state.json 和 git revert。扣分："会话中断后可通过 `.forge/test-state.json` 恢复清理"只有一句话——谁写这个文件？何时写？恢复时谁来读？读后做什么？均未定义。

### 7. Scope Clarity — 89/100

**In-scope 具体交付物 (32/35):**
checkbox 格式交付物列表。扣分：(1) "SKILL.md 新增 surface 检测步骤和 surface 感知配方生成流程"——"流程"作为交付物偏模糊；(2) "16 个 prompt 模板"未列出具体模板清单。

**Out-of-scope (27/30):**
6 项明确排除。扣分："新增 forge CLI 命令"在 Out of Scope 中，但 In Scope 中 `forge surfaces` CLI 标记为"需新建"——自相矛盾。

**Scope 一致性 (30/35):**
4 个 Story 与 Scope 对应，Related Changes 表覆盖完整。扣分："移除 test.execution 节点文档"没有对应的用户故事，也未描述对现有用户配置的影响。

## New Attack Points (Iteration 2)

1. **[Scope Clarity]** `forge surfaces` CLI 在 In Scope（"需新建"）与 Out of Scope（"新增 forge CLI 命令"）之间存在矛盾 — In Scope: "前置依赖：`forge surfaces` CLI 命令...当前状态：需新建"；Out of Scope: "新增 forge CLI 命令" — 必须明确：要么将 forge surfaces CLI 从 Out of Scope 列表中排除并注明"本特性包含此前置依赖"，要么将 In Scope 中的实现项改为"假设 forge surfaces CLI 已就绪"。

2. **[Flow Completeness]** teardown 执行失败的行为完全未定义 — Flow 中多处提到"执行 teardown"作为失败恢复动作（probe 失败、test 失败），但 teardown 本身失败时的行为未定义。Reliability 节说"teardown 幂等（PID 不存在时跳过）"但这只覆盖了一种情况 — 必须定义 teardown 失败（如权限不足、磁盘满）的退出码和后续行为。

3. **[Scenario Completeness]** Mermaid 图中 cli/tui 分支缺少错误处理路径 — Flow Description 步骤 4 说"每步检查退出码：exit 0 继续；exit 1/2 执行 teardown"，但 Mermaid 图中 Test2（cli/tui）直接到 RunEnd 没有任何 exit code 分支 — 要么在 Mermaid 图中为 cli/tui 添加错误分支（与 web/api 保持一致），要么在文档中明确说明 cli/tui 的 test 为什么不产生非零退出码。

4. **[Edge Case Coverage]** surface-key 名称作为 just 配方名一部分时的字符限制未定义 — 混合项目节说"配方名带 surface-key 前缀：`dev-<surface-key>`"，just 配方名不允许空格和部分特殊字符 — 必须定义 surface-key 的命名约束（允许的字符集、长度限制），或在 init-justfile 生成时做清洗/校验。

5. **[blindspot]** `.forge/test-state.json` 恢复机制是空壳 — "会话中断后可通过 `.forge/test-state.json` 恢复清理"只此一句，无任何具体流程：谁负责写入、何时写入、文件格式、恢复时的读取逻辑和清理动作 — 要么补充完整的恢复流程描述（至少包含写入时机、文件格式、恢复步骤），要么从文档中移除此承诺以避免误导实现者。

6. **[blindspot]** "编排级配方"概念缺少定义清单 — Flow 步骤 3 说"Surface 规则覆盖语言模板的编排级配方（test/dev/run/probe）"——括号中列了 4 种，但 Scope 中 init-justfile 列了"test/dev/run/probe/test-setup"5 种，且 mobile 的 test-setup 也属于编排级。概念边界不清会导致实现者对"哪些配方参与仲裁"产生歧义 — 必须在 Flow Description 中明确定义"编排级配方"的完整清单和判定标准。

7. **[blindspot]** Config schema 移除 test.execution 对现有用户的影响未评估 — "移除 `test.execution` 节点文档"在 Scope 中但无用户故事描述影响。现有用户如果已有 test.execution 配置，升级后这些配置会被静默忽略还是报错？— 必须添加迁移说明或用户通知策略。

## Comparison: Iteration 1 vs Iteration 2

| Dimension | Iter 1 | Iter 2 | Delta |
|-----------|--------|--------|-------|
| Background & Goals | 88 | 89 | +1 |
| Flow Diagrams | 135 | 130 | -5 |
| Flow Completeness | 165 | 180 | +15 |
| User Stories | 170 | 178 | +8 |
| Scenario Completeness | 115 | 127 | +12 |
| Edge Case Coverage | 75 | 79 | +4 |
| Scope Clarity | 80 | 89 | +9 |
| **Total** | **828** | **872** | **+44** |

**Remarks**: 第一轮的 8 个攻击点全部得到实质性修复。Flow Completeness（+15）、Scenario Completeness（+12）、Scope Clarity（+9）改进最大。Flow Diagrams 微降 5 分是因为第二轮评审标准更严格——cli/tui 分支的错误处理不一致在第一轮未被捕获。第二轮新发现 7 个攻击点，主要集中在：(1) Scope 内部矛盾（forge surfaces CLI 归属），(2) teardown 失败恢复的完整性，(3) Mermaid 图与文字描述的不一致。
