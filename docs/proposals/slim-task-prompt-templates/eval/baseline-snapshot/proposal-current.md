---
created: "2026-05-27"
author: "forge-brainstorm"
status: Draft
---

# Proposal: 精简任务 Prompt 模板

## Problem

Forge 的 15 个任务 prompt 模板包含大量非指令内容——注释、解释性描述、冗长的角色定义——这些内容不指导 agent 行为，只增加 token 消耗并稀释指令清晰度。同时 task-executor agent 的 Execution Protocol 存在步骤冗长、逻辑重叠的问题。

### Evidence

模板中冗余内容的量化分析：

| 冗余类别 | 出现文件数 | 每处损失 | 总冗余 |
|---------|-----------|---------|-------|
| HTML 注释 | 1 | 4 行 | 4 行 |
| Step 2 解释性描述 | 5 (gen-* / test-run / verify) | 1-2 行 | ~7 行 |
| 冗长角色描述 | 10 (coding.*, gate, doc) | 1 行 | ~10 行 |
| CODING_PRINCIPLES 解释性冗余 | 5 (coding.*) | 每原则 2-5 行 | ~50 行 |
| AC 验证块冗余 | 9 (coding.*, gate, doc) | 每处 ~12 行可缩至 ~4 行 | ~70 行 |
| Record Fields 描述性文字 | 9 | ~3 行可缩至 ~1 行 | ~20 行 |
| task-executor Execution Protocol | 1 | 11 步可合并为 8 步 | ~30 行 |

总计：约 **190 行** 非指令冗余。注意并非每个任务都加载全部 200 行——单个 coding.* task 的模板包含约 80-100 行冗余（AC 验证块为主体），gate/doc/test 类更少（约 20-40 行）。200 行为全模板集上界，操作以单模板实际冗余为准。

**Token 估算**（以 Claude Sonnet tokenizer 为参考）：每行 markdown 平均 ~15 tokens，coding.* 模板 ~80-100 行冗余 ≈ 1200-1500 tokens/task。按 daily task 量估算（以团队每日 10-20 个 coding.* task 计），每日 token 节省约 12K-30K tokens（输入侧），月度约 250K-600K tokens。数值为近似，精确测算在精简后执行。


### Urgency

- 每个 task 执行都在消耗这些冗余 token，日积月累规模可观
- 清晰的 prompt 减少 agent 误解和执行偏差
- Prompt 精简是持续优化的一部分，目前已有 prompt-template-audit 等基础，可以在此基础上推进

## Proposed Solution

**就地精简**：保持现有模板独立，在每个文件内部删除非指令内容，将模糊描述改为清晰指令。

不抽取公共模块，不改变现有分类体系。

### Innovation Highlights

本方案不是技术创新，而是对现有 prompt 的"清理"。核心原则是"prompt 是指令，不是文档"——删掉所有不能直接指导 agent 行动的文字。

**行业参照：** 本方案的设计哲学与以下行业实践一致：
- **LangChain Prompt Templates** 在模板中区分"指令（instructions）"与"上下文（context）"，推荐仅将直接影响模型行为的文本保留在系统 prompt 中，解释性描述移至外部文档。
- **Anthropic Prompt Engineering Guide** 强调"show, don't just tell"——通过示例约束行为而非通过自然语言角色描述；本方案中的 AC 验证块精简（保留 REQUIRED 指令、删除展开说明）遵循同一原则。
- **OpenAI GPTs Instructions** 模式的演变方向也是删除冗余的系统 prompt 装饰，改用精确的祈使句指令。

## Requirements Analysis

### Key Scenarios

1. **coding-feature / coding-enhancement / coding-fix / coding-cleanup / coding-refactor** 五个核心模板：
   - 角色描述从自然语言改为祈使句
   - CODING_PRINCIPLES 去掉举例和解释，保留核心约束
   - AC 验证块从 ~12 行精简到 ~4 行
   - Step 2 的实现说明保留，只去修饰性语言

2. **gate / doc** 模板：
   - 角色描述精简
   - AC 验证块精简

3. **test-run / test-gen-scripts / test-gen-contracts / test-gen-journeys / test-verify-regression** 模板：
   - Step 2 中的 "This generates X from Y..." 解释性描述删除
   - 角色描述精简

5. **code-quality-simplify / validation-code / validation-ux** 模板（共 3 个，约 30-50 行/个）：
   - 角色描述精简（同 coding-* 模式）
   - 无 AC 验证块和 CODING_PRINCIPLES——冗余集中在角色描述和框架性说明行
   - code-quality-simplify：~35 行，含 5 行角色描述 + 3 行 Record Fields 说明 → 可精简 ~6 行
   - validation-code：~45 行，含 4 行角色描述 + 4 行 AC 验证说明 + 3 行 Record Fields 说明 → 可精简 ~8 行
   - validation-ux：~50 行，同上模式 → 可精简 ~8 行
   - 三者合计精简约 20 行

4. **task-executor agent**：
   - Execution Protocol 步骤合并（步骤 4/5/6 处理 prompt 获取逻辑可合并为 1 步）
   - Retry Strategy 与 Complex Error Pause Flow 去重合并
   - 输出格式合并为紧凑格式

   **步骤 4/5/6 合并前错误恢复分析：** Step 4 读取模板（失败不可恢复，终止），Step 5 替换变量（失败可降级继续），Step 6 组装（无独立失败场景）。Steps 4-6 构成严格顺序链，合并后错误恢复路径不变，合并安全。

   **Retry 与 Error Pause 正交性分析：** Retry 操作单次 LLM 调用的临时错误，Error Pause 操作整个 task 的持久错误，层级不同，正交可合并。

AC 验证块逐行分析：

| 行类型 | 典型数量 | 处理策略 | 压缩后行数 |
|--------|---------|---------|-----------|
| AC:REQUIRED 指令 | 3-5 行 | 保留 | 3-5 行 |
| 指令展开说明 | 3-5 行 | 合并至指令行 | 0 |
| 场景举例 | 0-2 行 | 删除 | 0 |
| 格式装饰 | 3-4 行 | 保留必要空行 | 1-2 行 |

~12 行 → ~4 行（66%），功能行保留。

CODING_PRINCIPLES 逐原则分析：

| 原则条目 | 行数 | 功能判定 | 处理策略 |
|---------|------|---------|---------|
| 原则 1: 纯指令行 | 1 行 | 核心约束 | 保留 |
| 原则 1: 行为示例/边界说明 | 2-5 行 | 约束边界演示——非核心指令，但可能作为 few-shot 约束模型行为 | 压缩为 1 行边界概括。注意：这是表示变化而非纯压缩——若 SC2 轨迹一致性 < 90%，每原则保留 1 个代表性示例 |
| 原则 2: 纯指令行 | 1 行 | 核心约束 | 保留 |
| 原则 2: 反例/边界说明 | 2-5 行 | 约束边界演示 | 压缩为 1 行边界概括。同上——SC2 检测偏离时回退保留示例 |
| 超原则通用说明（如作用域声明） | 1 行 | 元指令 | 保留 |

~50 行 → ~20 行（60%），每原则保留 1 行指令 + 1 行边界概括。

Record Fields 逐字段分析：

| 行类型 | 典型数量 | 处理策略 |
|--------|---------|---------|
| 字段名 + 值（如 `## Output\n{...}`） | 1 行 | 保留 |
| 字段用途说明（如 "This field describes..."） | 1-2 行 | 删除——字段名自解释 |
| 格式示例/占位符展开 | 1-2 行 | 删除——嵌入实际值即可 |

~3 行 → ~1 行（66%），字段名和值保留。

### Non-Functional Requirements

- 精简后模板的指令覆盖必须与精简前等价（不能遗漏 agents 需要知道的信息）
- 所有 task-executor 的行为不发生变化

### Constraints & Dependencies

- 模板文件位于 `forge-cli/pkg/prompt/data/*.md`，由 `prompt.go` 通过 embed FS 加载
- 修改模板不影响 Go 代码，只需修改 .md 文件
- task-executor agent 位于 `plugins/forge/agents/task-executor.md`

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| 分层模板组合 | LangChain PromptTemplate / Vercel AI SDK | 语义分离（instruction/tool/context 分层），单层修改不影响其他，改一处影响所有 | 需重构模板分类体系并修改 `prompt.go` 加载逻辑，与"不改后端代码"约束冲突；对 15 个文件引入抽象层，改动面大于收益 | Rejected: 架构约束否决 |
| 引入 DSL 生成 | 模板引擎模式 | 声明式模板定义，通过编译生成最终 prompt，压缩逻辑集中在 DSL 层 | 需要增加 DSL 定义文件、解析器、编译管线，对 15 个小模板引入完整工具链成本过高——模板改动频次低（月级而非天级），DSL 抽象层在小规模场景下维护负担超过收益 | Rejected: 模板规模小、变更频次低，DSL 工具链成本不合理 |
| 什么都不做 | — | 零风险 | token 持续浪费、指令不够清晰 | Rejected: 成本太低 |
| 抽取公共模块 | DRY 模式 | 修改一处同步所有模板 | 需要改 `prompt.go` 逻辑，且被用户否决 | Rejected: 不满足就地要求 |
| **就地精简** | Forge 现有风格 | 零架构变更，每模板独立修改，风险隔离 | 每个文件都要改 | **Selected: 简单直接** |

## Feasibility Assessment

### Technical Feasibility

纯文本编辑，无技术风险。

### Resource & Timeline

10-15 个文件的文字精简，1 次编码任务即可完成。

### Dependency Readiness

前置条件：本次 brainstorm 输出的 proposal 通过。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "我需要保留角色描述让 agent 理解上下文" | Assumption Flip | 角色描述中的自然语言（"You are a focused..."）对 LLM 行为的影响可由祈使句替代——但这仍是假设而非定论。该领域存在争议（部分研究表明系统角色有效，亦有研究显示模型更遵循后续指令），需通过实施后的行为等价验证确认。**此假设与 NFR #2（"所有 task-executor 的行为不发生变化"）存在固有张力：前者承认不确定性，后者要求确定性。该张力通过 SC2 trial run 协议化解——若 SC2 检测到轨迹偏移（一致性 < 90%），回退角色描述修改部分，保留其他精简项。即角色描述精简为有回退机制的实验性变更，而非无条件承诺不变。** |
| "每个模板独立意味着不需要关注跨模板一致性" | XY Detection | 用户确认了「核心流程重复是允许的」，所以跨模板一致性不是问题，不需要抽取公共模块。 |

## Scope

### In Scope
- 修改 `forge-cli/pkg/prompt/data/` 下全部 15 个模板文件：
  - coding-feature.md, coding-fix.md, coding-enhancement.md, coding-cleanup.md, coding-refactor.md
  - gate.md, doc.md
  - test-run.md, test-gen-scripts.md, test-gen-contracts.md, test-gen-journeys.md, test-verify-regression.md
  - code-quality-simplify.md, validation-code.md, validation-ux.md
- 修改 `plugins/forge/agents/task-executor.md`
- 删除 HTML 注释
- 删除 Step 2 解释性描述（"This generates X from Y"）
- 精简角色描述（自然语言 → 祈使句）
- 精简 AC 验证块（~12 行 → ~4 行）
- 精简化 Record Fields（去掉引导性描述，保留字段名和值）
- 精简 CODING_PRINCIPLES（去掉举例和解释）
- 精简 task-executor Execution Protocol（合并步骤）

### Out of Scope
- 不抽取公共模块文件
- 不修改 `prompt.go` 代码逻辑
- 不新增/删除模板文件
- 不改动模板占位符（`{{TASK_ID}}` 等）
- 不改动 Spec Authority Enforcement 逻辑结构（保留现有约束块内容，不增不减）
- 不改动 Hard Rules / CRITICAL 块的逻辑

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 精简过度导致 agent 遗漏关键行为 | Low | High | **制品：** 每模板建立"功能快照清单"——JSON 格式节点台账，每个节点包含：`{id, category, type, content_snippet, role}`。category 为 instruction/constraint/example/format 之一，type 进一步细分（如 hard-rule/critical/ac-required/ac-explanation/role-desc/record-field），role 标明保留/删除/合并。**创建时机：** 修改前由修改者按模板逐行标注，reviewer 签署确认。**存储：** 仓库 `scripts/function-snapshots/<template-name>.json`，版本化管理支持 PR diff。**流程：** (1) 修改前签署清单；(2) 修改后逐项比对，每项标记 pass/fail；(3) 全部 pass 方可合并。**判断标准：** 清单中任一项 fail → 回滚修改并重新调整。|
| 多个模板同步修改，跨模板不一致 | Medium | Medium | **基线文件：** 以 `coding-feature.md` 作为分类基准模板，其余 coding-* 模板修改后与其做 diff，确保同样角色结构、同样原则格式、同样 AC 精简模式。**操作者：** 修改者执行 diff，reviewer 确认 diff 差异为合理的特定语义差异（非格式漂移）。**pass/fail：** 非语义性结构差异 > 3 处 → fail。|
| 现有测试基础设施无法检测 prompt 层行为漂移——当前 forge 测试覆盖 Go 代码逻辑和 task 执行结果，但无机制检测因 prompt 修改导致的 agent 行为差异（如指令理解偏差、约束优先级变化） | Medium | High | **治理措施：** 上述 Risk 1 功能快照清单覆盖指令/约束/示例/格式全部节点，修改前后逐项比对。**补充手段：** (1) 修改后对每个模板对应的典型测试场景执行 2 次 trial run（配合 SC2 自动化轨迹 diff 脚本），对比输出一致性；(2) 非删除项语义一致性检查，100% 覆盖，否则回滚。**CI 化：** (a) 功能快照清单存储为版本化 JSON（`scripts/function-snapshots/` 目录），PR 自动 diff 检查节点不可被意外删除；(b) 轨迹对比脚本作为可选 PR check（不阻塞合并但报告差异）|
| prompt 变更为有状态修改，合入后发现影响需要回滚但无标准化流程 | Medium | Medium | **回滚流程：** (1) 每批模板修改独立提交，禁止单 commit 修改全部 16 个文件——分 3 批提交（coding-* 为 1 批、gate/doc 为 1 批、test-* 为 1 批），批间间隔至少 1 个 CI 周期以确保隔离。(2) 合入后观察期：合入后运行一轮完整 journey 测试（`just test-e2e`），若任一 journey 出现与 baseline 不同的行为且确认为 prompt 修改导致，立即 `git revert` 对应批次的 commit。(3) 无需 feature flag——模板为静态文件，revert 即恢复行为。(4) 应急备选：若 revert 冲突（如中间插入了其他 commit），从 baseline snapshot（`eval/baseline-snapshot/`）复制回原始模板重新提交。|

## Success Criteria

**主要指标（保留率）和次要指标（行数）的双层结构：** 保留率为首要校验门禁，行数压缩为次要效率指标——当保留率不达标时禁止合并，行数压缩不达标可接受。

### 功能保留（首要门禁）

- [SC1] 功能约束保留率 **100%**——每个模板修改后，对照功能快照清单逐项比对，所有指令/约束/格式节点保留率为 100% 方可合并。节点分类包括：(1) Hard Rules、(2) CRITICAL 块、(3) Spec Authority Enforcement、(4) CODING_PRINCIPLES 各原则的指令行（不含边界说明——边界说明允许按 SC3 压缩）、(5) Record Fields 字段名与值结构、(6) AC:REQUIRED 指令。检测方法：修改者逐节点标注 pass/fail，reviewer 签署确认。

### 行为等价性

- [SC2] 模板精简后，agent 执行相同 task 的行为无可见差异。**检测协议：** (1) 选取典型 task——覆盖规则：每个 template 选取 1 个典型 task，该 task 必须至少覆盖该 template 功能快照清单中 80% 的 instruction/constraint 类别节点（如 coding-feature 的典型 task 必须触发 AC 验证、compile/fmt/lint/test 全部 4 个步骤）。典型 task 由修改者从现有 task 库中选取，reviewer 确认覆盖率。(2) 分别在修改前/后模板上执行该 task 各 2 次（共 4 次 run）；(3) 对比同一 task 的 agent 执行轨迹（关键步骤序列、工具调用参数、最终输出结构）；(4) 轨迹一致性 ≥ 90%（容差：步骤顺序因 LLM 生成随机性导致的非功能性差异）视为通过。**自动化：** 轨迹对比通过脚本自动完成——提取每次 run 的步骤名序列、工具调用名、输出结构 key 集合，生成结构化 diff 报告，仅差异判定环节需人工介入。该脚本纳入仓库 `scripts/` 目录，作为 PR check 的可选验证（不阻塞合并但报告差异）。**注意：** 此 SC 为熊市测试——通过不一定保证完全等价，不通过则明确失败。

### 结构验证

- [SC3] CODING_PRINCIPLES 在 5 个 coding-* 模板中保留全部核心约束指令（每原则至少保留 1 行指令 + 1 行边界概括），通过 diff 确认对比精简前后原则覆盖无遗漏。
- [SC4] Record Fields 在所有出现模板中保留字段名和值结构，字段用途说明可删除。通过 diff 确认字段名节点未被误删。
- [SC5] Step 2 解释性描述（"This generates X from Y..."）在 5 个 test-* / verify 模板中全部删除，通过 grep 确认无残留。

### 效率指标

- [SC6] 15 个模板文件 + task-executor 共减少 **≥150 行**（去除注释、解释性描述、冗长定义）。次要指标——保留率门禁未通过时，此行数指标不构成合并理由。
- [SC7] task-executor 的 Execution Protocol 步骤数从 11 步减少到 ≤8 步。

## Next Steps

- Proceed to `/write-prd` to formalize requirements