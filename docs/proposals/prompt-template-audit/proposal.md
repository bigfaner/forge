---
name: prompt-template-audit
status: Draft
created: 2026-05-20
---

# Prompt 模板体系全链路审查 + 优化方案

## Problem

Forge 的 typed-task-dispatch 系统通过 22 个嵌入模板 + 1 个 agent 定义 + 1 个合成引擎将任务分派给 task-executor 子代理。这些模板是 agent 执行任务的唯一指令来源，但目前缺乏系统性的质量审查。具体问题包括：

- **结构不统一**：各模板之间的 header 声明、步骤命名、提交处理方式不一致，增加维护成本
- **潜在双重提交**：6 个模板内嵌了 submit-task 步骤，但 task-executor agent 也会在策略执行完毕后自动调用 submit-task
- **上下文利用率低**：PHASE_SUMMARY、SCOPE 等占位符在不同模板中的使用和引用方式不一致
- **跨模板一致性**：各模板独立但需要保持同步的通用部分（如 static checks 流程、失败处理表）有不一致
- **缺失关键指令**：部分模板缺少 conventions 加载、Hard Rules 处理等通用步骤

## Solution

对 22 个 prompt 模板 + task-executor.md + prompt.go 进行逐模板四维深度审查（结构清晰度、步骤合理性、指令精确度、冗余/缺失），产出评估报告和具体优化方案。

## Scope

### In Scope

- `forge-cli/pkg/prompt/data/` 下 22 个嵌入模板的逐模板审查
- `plugins/forge/agents/task-executor.md` agent 定义的审查
- `forge-cli/pkg/prompt/prompt.go` 合成逻辑的审查
- 模板间一致性分析
- 每个问题的具体修复建议

### Out of Scope

- task 文件模板（task.md/task-doc.md）—— 属于任务生成体系，非 prompt 分派体系
- fix-task/cleanup-task 模板 —— 属于运行时任务创建，非 agent 执行指令
- 其他 skill 定义文件的内容审查
- 实际修改模板文件（本提案仅产出方案，执行需另建任务）
- 模板合并或共享片段抽取（违反独立模板原则）

### Design Constraints

1. **每个模板独立自包含**：agent 只收到一份完整 prompt，不做跨文件引用
2. **不做抽取合并**：即使有重复内容，也保持每个模板完整独立，便于单独理解和维护
3. **简单可靠优先**：改动应降低复杂度而非引入新的抽象层

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| 优化建议引入新的不一致 | 模板改动后行为回归 | 方案附带一致性检查清单 |
| 过度统一导致特殊场景丢失 | refactor/validation 等特殊类型执行失败 | 保留类型特定的差异，仅统一通用部分 |
| 双重提交问题的实际影响不确定 | 可能是 benign（幂等）也可能是 bug | 方案中明确标注需验证 |

## Success Criteria

- [ ] 每个模板有四维评估结论（结构/步骤/指令/上下文）
- [ ] 跨模板一致性问题全部列出并分级（P0-P3）
- [ ] 双重提交问题有明确的根因分析和修复建议
- [ ] 优化方案包含具体的 diff 级别修改指导
- [ ] 所有建议附带风险评估（改动影响范围）

---

# 评估报告

## 1. 全局发现（跨模板系统性问题）

### 1.1 [P0] 双重提交：模板内嵌 submit 步骤 vs task-executor 自动提交

**影响模板**：doc.md、doc-eval.md、doc-summary.md、doc-consolidate.md、doc-drift.md、clean-code.md

**问题**：这 6 个模板在步骤末尾显式要求调用 `Skill(skill="forge:submit-task")`，但 task-executor agent 定义的第 8 步也总是调用 submit-task。结果：submit-task 被调用两次。

**根因**：模板设计时未考虑 task-executor 的自动提交机制。doc 类模板最早编写时 submit 尚未集成到 agent 中，后来 agent 加入了自动提交但模板未同步清理。

**风险评估**：如果 submit-task 是幂等的（覆盖同一条记录），则影响为 benign；如果会创建重复记录或产生 side-effect，则是一个 bug。

**修复建议**：
- **方案 A（推荐）**：从这 6 个模板中移除显式 submit 步骤，统一依赖 task-executor 的自动提交。需验证 submit-task 的幂等性。
- **方案 B**：保留模板内的 submit，但在 task-executor 中添加判断——如果策略最后一步已经是 submit，则跳过自动提交。更复杂，但兼容性更好。

---

### 1.2 [P1] Header 变量声明不一致

**影响模板**：全部 22 个

各模板的 header 变量声明存在以下不一致：

| 变量 | 声明模板数 | 缺失模板 |
|------|-----------|---------|
| TASK_ID | 22/22 | — |
| TASK_FILE | 22/22 | — |
| SCOPE | 20/22 | doc.md, doc-eval.md |
| FEATURE_SLUG | 4/22 | 仅 doc.md, doc-eval.md, doc-summary.md, doc-eval.md |
| PHASE_SUMMARY | 16/22 | 缺失：coding-fix, doc-eval, gate, fix-record-missed, test-eval-cases, clean-code |
| COVERAGE_STRATEGY/TARGET | 5/22 | 仅 coding.* 模板 |
| TEST_TYPE_ARG | 3/22 | 仅 test-gen-scripts, test-gen-and-run 相关 |

**分析**：
- SCOPE 缺失的 2 个 doc 模板确实不需要 scope（文档任务不跑 just 命令）→ **合理**
- FEATURE_SLUG 仅在 doc-summary 中实际使用（构建 records 路径），doc.md 和 doc-eval.md 声明了但未使用 → **冗余声明**
- PHASE_SUMMARY 缺失的模板中：
  - coding-fix（修复任务）—— 可能需要上一阶段的上下文，**缺失**
  - doc-eval（文档评估）—— 不需要阶段上下文，**合理**
  - gate（阶段门检查）—— 正在检查的阶段上下文由 gate 任务文件本身提供，**合理**
  - fix-record-missed（恢复任务）—— 一次性验证，**合理**
  - test-eval-cases —— 测试评估独立于阶段，**合理**
  - clean-code —— 委托给 skill，**合理**
- COVERAGE_STRATEGY 仅用于 coding.* 类型 → prompt.go 的 `IsTestableType` 逻辑已正确处理，非 testable 类型替换为空串后被 cleanTemplateOutput 移除 → **合理**

**修复建议**：
- coding-fix 添加 PHASE_SUMMARY 支持
- doc.md 和 doc-eval.md 移除未使用的 FEATURE_SLUG 声明
- 其余保持现状

---

### 1.3 [P1] Step 1 "加载项目规范" 步骤不一致

**影响模板**：全部 22 个

模板在 Step 1 中加载 `docs/conventions/` 和 `docs/business-rules/` 的情况：

| 分类 | 模板 | 有规范加载 |
|------|------|-----------|
| coding.* | feature, enhancement, cleanup, refactor | 有 |
| coding.fix | fix | 有（但顺序不同） |
| doc* | doc, doc-eval, doc-summary, doc-consolidate, doc-drift | **无** |
| test.* | 全部 8 个 | **无** |
| validation.* | code, ux | 有 |
| gate | gate | **无** |
| clean-code | clean-code | **无**（委托给 skill） |
| special | fix-record-missed | **无** |

**分析**：
- coding.* 模板一致加载规范 → 合理，编码任务需要项目约定
- doc-consolidate/doc-drift 不需要 → 它们本身就是处理规范的，合理
- test.* 模板委托给具体 skill → skill 内部可能加载规范，但模板自身不指导 agent 加载 → **可接受**
- **gate 不加载规范** → gate 需要验证阶段交付物是否符合标准，应该加载规范 → **缺失**
- **doc/doc-eval 不加载规范** → 文档任务可能需要遵循项目文档规范 → **建议添加**

**修复建议**：
- 为 gate.md 添加规范加载步骤
- doc.md 和 doc-eval.md 可选择性添加（低优先级）

---

### 1.4 [P1] Hard Rules 处理不一致

各模板处理 task 文件中 Hard Rules 的方式：

| 模板 | 处理方式 | 语义 |
|------|---------|------|
| coding-feature, enhancement, cleanup | `<IMPORTANT>` + "Follow them exactly during the entire TDD cycle" | 编码期间遵守 |
| coding-refactor | `<IMPORTANT>` + "Follow them exactly throughout the entire workflow" | 全流程遵守 |
| coding-fix | `<IMPORTANT>` + "Respect file scope restrictions... Respect command restrictions..." | **不同的语义焦点** |
| validation-code, ux | `<IMPORTANT>` + "Treat every MUST as pass/fail criterion" | **验证视角** |
| gate | `<IMPORTANT>` + "Treat every MUST as pass/fail criterion" | **验证视角** |
| doc*, test*, clean-code, fix-record-missed | 无 Hard Rules 处理 | — |

**分析**：
- 编码模板的 "Follow exactly" 和验证模板的 "Treat as pass/fail" 是正确的语义区分
- coding-fix 的独特表述（聚焦文件范围限制）是有意的——修复任务有更严格的边界约束
- test.* 模板使用 `<HARD-RULE>` 标签定义了任务自身的硬规则（如 MUST invoke skill），这与 task 文件的 Hard Rules 是不同概念 → **命名冲突**
- doc.* 模板不处理 Hard Rules → 文档任务也可能有约束（如格式要求），但实际意义不大 → **可接受**

**修复建议**：
- 将 test.* 模板中的 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`，避免与 task 文件的 Hard Rules 概念混淆
- 其余保持现状

---

### 1.5 [P1] CODING_PRINCIPLES 数量不一致

| 模板 | CODING_PRINCIPLES 条目 |
|------|----------------------|
| coding-feature, enhancement, fix | 4 条（Think Before Coding + Simplicity First + Surgical Changes + Goal-Driven） |
| coding-cleanup | 2 条（Simplicity First + Surgical Changes） |
| coding-refactor | 1 条（Surgical Changes，标题格式不同——用了 ### 前缀） |

**分析**：
- cleanup 省略 "Think Before Coding" → 清理任务目标明确，不需要先分析 → **合理**
- cleanup 省略 "Goal-Driven" → 清理任务的成功条件由现有测试定义 → **合理**
- refactor 使用了 ### 标题格式而非 - 列表格式 → **格式不一致**

**修复建议**：
- coding-refactor 的 CODING_PRINCIPLES 统一使用 - 列表格式

---

### 1.6 [P2] "Static Checks + Targeted Tests" 步骤高度重复

**影响模板**：coding-feature, coding-enhancement, coding-cleanup, coding-fix（4 个模板共享近乎相同的 Step 3/4）

**问题**：这 4 个模板的静态检查步骤有约 15 行完全相同的内容（just 命令序列、失败处理表、targeted test 指令）。唯一差异在 coding-refactor 中有更复杂的 fmt/lint 处理逻辑。

**分析**：当前设计是每个模板自包含——agent 收到一个完整的 prompt，不需要跨文件引用。这种重复是故意的，确保 agent 在任何类型任务中都有完整指令。每个模板保持独立是正确的设计选择。

**修复建议**：
- **维持现状**——每个模板独立完整，不做抽取合并
- 确保修改任一模板的 static checks 步骤时，同步检查其他模板

---

## 2. 逐模板深度审查

### 2.1 coding-feature.md

**结构**：3 步流程（Read → TDD → Verify），清晰简洁。
**步骤合理性**：✅ 良好。TDD 周期 RED→GREEN→REFACTOR 是标准实践。
**指令精确度**：
- ⚠️ Step 2 的 TDD 描述过于笼统——"Follow the TDD cycle for each requirement"，但没有指导如何从 task 文件中提取 requirement
- ⚠️ `go test -race -cover ./changed/package/...` 是 Go 特定的示例，但模板应该语言无关——注释说"framework-native test commands"但示例会误导 agent

**上下文利用**：✅ PHASE_SUMMARY、SCOPE、COVERAGE 均正确使用。

**具体优化建议**：
1. Step 2 添加引导：先从 task 文件的 Acceptance Criteria 中提取测试需求，再逐个执行 TDD
2. 将 Go 示例改为更通用的描述："运行受影响包/模块的测试命令"

---

### 2.2 coding-enhancement.md

**结构**：3 步流程，与 coding-feature 几乎相同。
**步骤合理性**：✅ 良好。
**指令精确度**：
- 与 coding-feature 有 90% 内容重复，差异仅在：
  1. 角色描述："enhancing an existing feature" vs "implementing a new feature"
  2. TDD 描述："captures the desired behavior improvement" vs "Write failing test first"
  3. 额外一行："Review existing tests for the code being enhanced"

**冗余分析**：与 coding-feature 高度重复。设计原则：每个模板保持独立自包含，不做抽取合并。

**具体优化建议**：
1. 维持独立模板，在两者头部添加同步维护注释，提醒修改 coding-feature 时同步检查 coding-enhancement

---

### 2.3 coding-cleanup.md

**结构**：3 步流程，清晰。与 coding-feature 结构一致但 CODING_PRINCIPLES 更少。
**步骤合理性**：✅ 不写新测试、依赖现有测试的策略对清理任务合理。
**指令精确度**：✅ 明确列出了清理类型（dead code、unused declarations、fixing existing tests）。

**具体优化建议**：无重大问题。

---

### 2.4 coding-refactor.md

**结构**：4 步流程（Read → Impact Map → Refactor → Verify），是最复杂的模板（185 行 vs 其他 ~80 行）。
**步骤合理性**：✅ 非常好。Add→Migrate→Remove 阶段化策略安全且可回滚。Pre-check 确保干净起点。
**指令精确度**：
- ✅ Impact Mapping 分类（Structural vs Behavioral）精确且实用
- ✅ 动态耦合扫描（reflection、string-based type checks）是独有且高价值的指令
- ⚠️ Step 4 的 fmt 失败处理过于复杂（`git stash && just fmt && git diff --name-only && git stash pop`）——可能让 agent 执行出错
- ⚠️ Batch sizing 规则（≤10 全部、>10 分 15-20 或 3-5）是硬编码数值，不适合所有项目规模

**上下文利用**：✅ PHASE_SUMMARY、SCOPE、COVERAGE 均正确使用。增量编译策略利用了 SCOPE。

**具体优化建议**：
1. 简化 fmt 失败处理：改为"如果 fmt 只影响 refactor 涉及的文件，修复；如果影响的是无关文件，继续"
2. 考虑让 batch sizing 参考项目的文件规模（但模板无法获取这一信息 → 维持现状）

---

### 2.5 coding-fix.md

**结构**：4 步流程（Read → Locate → Fix → Verify），比其他 coding 模板多一个 Locate 步骤。
**步骤合理性**：✅ "先定位再修复"的流程对 bug fix 合理。
**指令精确度**：
- ⚠️ Step 1 中 conventions 加载在 task 文件读取之后（与其他模板相反）
- ✅ E2E test failure 的特殊处理（不启动 dev server）是精确且有价值的
- ⚠️ 缺少 PHASE_SUMMARY（其他 coding 模板都有）

**上下文利用**：⚠️ 缺少 PHASE_SUMMARY。

**具体优化建议**：
1. 将 Step 1 中的 conventions 加载移到 task 文件读取之前，与其他 coding 模板一致
2. 添加 PHASE_SUMMARY header 声明和 Step 1 的条件加载

---

### 2.6 doc.md

**结构**：4 步流程（Read → Execute → Self-Check → Submit），比其他 doc 模板多一个 Self-Check。
**步骤合理性**：✅ Self-Check 是文档任务独有的好实践。
**指令精确度**：
- ⚠️ Step 2 "Execute Document Work" 过于模糊——"Perform the documentation work described in the task file" 没有给出具体指导
- ⚠️ 声明了 FEATURE_SLUG 但未使用

**冗余分析**：⚠️ Step 4 的 submit 步骤与 task-executor 自动提交重复（P0 问题 1.1）。

**具体优化建议**：
1. Step 2 添加引导：先识别任务类型（创建/修改/删除文档），然后按类型执行
2. 移除未使用的 FEATURE_SLUG header
3. 移除 Step 4 的显式 submit 调用

---

### 2.7 doc-eval.md

**结构**：3 步流程（Read → Evaluate → Submit），自带完整的 8 维度 rubric。
**步骤合理性**：✅ 3 轮迭代 + 900 分阈值是合理的质量标准。
**指令精确度**：
- ⚠️ 8 维度 rubric 全部权重相同（各 125 分），但某些维度（如 Accuracy、Completeness）在实际评估中应更重要
- ⚠️ 声明了 FEATURE_SLUG 但未使用
- ⚠️ 缺少 PHASE_SUMMARY（文档评估通常不需要 → 合理但不一致）

**冗余分析**：⚠️ Step 3 的 submit 步骤与 task-executor 自动提交重复。

**具体优化建议**：
1. 移除未使用的 FEATURE_SLUG header
2. 移除 Step 3 的显式 submit 调用
3. 考虑为 rubric 维度添加权重差异（低优先级——当前等权设计简单明确）

---

### 2.8 doc-summary.md

**结构**：3 步流程（Read → Generate → Submit），结构清晰。
**步骤合理性**：✅ 5 段式总结结构是经过设计的好模板。
**指令精确度**：✅ 精确指定了记录读取路径和输出格式。

**上下文利用**：✅ 正确使用了 FEATURE_SLUG（构建 records 路径）和 PHASE_SUMMARY。

**冗余分析**：⚠️ Step 3 的 submit 步骤与 task-executor 自动提交重复。

**具体优化建议**：
1. 移除 Step 3 的显式 submit 调用

---

### 2.9 doc-consolidate.md

**结构**：3 步流程（Read → Consolidate → Submit），委托给 consolidate-specs skill。
**步骤合理性**：✅ 非交互模式说明明确。
**指令精确度**：✅ 清晰说明了 auto-integrate 行为和 `[auto-specs]` commit tag。

**冗余分析**：
- ⚠️ Step 3 的 submit 步骤与 task-executor 自动提交重复
- ⚠️ SCOPE 声明在 header 但未使用（consolidate 不跑 just 命令）

**具体优化建议**：
1. 移除 Step 3 的显式 submit 调用
2. 移除未使用的 SCOPE header

---

### 2.10 doc-drift.md

**结构**：3 步流程（Read → Detect → Submit），与 doc-consolidate 结构相同。
**步骤合理性**：✅ drift-only 模式说明清晰。
**指令精确度**：✅ 清楚区分了 drift-only 和 full consolidation 模式。

**冗余分析**：
- ⚠️ Step 3 的 submit 步骤与 task-executor 自动提交重复
- ⚠️ SCOPE 声明在 header 但未使用

**具体优化建议**：
1. 移除 Step 3 的显式 submit 调用
2. 移除未使用的 SCOPE header

---

### 2.11 test-gen-cases.md

**结构**：2 步流程（Read → Generate），简洁。
**步骤合理性**：✅ 委托给 skill，模板本身是轻量 wrapper。
**指令精确度**：✅ HARD-RULE 明确要求必须使用 skill。

**命名问题**：⚠️ `<HARD-RULE>` 标签与 task 文件的 "Hard Rules" 概念混淆。

**具体优化建议**：
1. 将 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`

---

### 2.12 test-eval-cases.md

**结构**：2 步流程（Read → Evaluate）。
**步骤合理性**：✅ 委托给 eval skill。
**指令精确度**：
- ⚠️ 包含 `<EXTREMELY-IMPORTANT>` 声明 "This task runs in the MAIN SESSION"——这是唯一的 test 模板有此声明
- 这个信息应该由 dispatcher（execute-task/run-tasks）在分派时决定，而不是硬编码在模板中

**具体优化建议**：
1. 移除 MAIN_SESSION 声明（由 dispatcher 控制）
2. 将 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`

---

### 2.13 test-gen-scripts.md

**结构**：2 步流程（Read → Generate），简洁。
**步骤合理性**：✅ 清晰。
**指令精确度**：✅ `{{TEST_TYPE_ARG}}` 的条件注入设计合理。

**具体优化建议**：
1. 将 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`

---

### 2.14 test-run.md

**结构**：2 步流程（Read → Run），简洁。
**步骤合理性**：✅ 失败重试逻辑（max 3 attempts）合理。
**指令精确度**：✅ 明确禁止直接运行测试命令。

**具体优化建议**：
1. 将 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`

---

### 2.15 test-gen-and-run.md

**结构**：3 步流程（Read → Generate → Run），组合了 gen-scripts 和 run。
**步骤合理性**：✅ 两阶段执行逻辑清晰。
**指令精确度**：✅ Hard Rules 覆盖了两个阶段的约束。

**具体优化建议**：
1. 将 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`

---

### 2.16 test-graduate.md

**结构**：2 步流程（Read → Graduate），简洁。
**步骤合理性**：✅ 清晰。
**指令精确度**：✅ 明确禁止手动移动文件。

**具体优化建议**：
1. 将 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`

---

### 2.17 test-verify-regression.md

**结构**：2 步流程（Read → Run Regression）。
**步骤合理性**：✅ 直接运行 just test-e2e 是正确的验证方式。
**指令精确度**：
- ✅ 明确禁止手动启动 dev server
- ✅ 明确限制修复范围（source code or test selectors only）

**具体优化建议**：
1. 将 `<HARD-RULE>` 改名为 `<TASK-CONSTRAINTS>`

---

### 2.18 validation-code.md

**结构**：2 步流程（Read → Validate），包含 criteria 检查 + quality gate。
**步骤合理性**：✅ 先检查 criteria 再跑 gate，逻辑清晰。
**指令精确度**：
- ✅ 失败处理有 trivial/non-trivial 区分
- ✅ 质量门包含完整 4 步（compile + fmt + lint + test）

**具体优化建议**：无重大问题。

---

### 2.19 validation-ux.md

**结构**：2 步流程（Read → Validate），与 validation-code 结构相同但无 quality gate。
**步骤合理性**：✅ UX 验证不需要编译/测试。
**指令精确度**：✅ 与 validation-code 一致的失败处理逻辑。

**具体优化建议**：无重大问题。

---

### 2.20 gate.md

**结构**：2 步流程（Read → Verify），包含 criteria 检查 + quality gate + mermaid 流程图。
**步骤合理性**：✅ 阶段门检查逻辑完整。
**指令精确度**：
- ✅ mermaid 流程图提供了清晰的决策可视化
- ⚠️ 缺少 conventions 加载——gate 检查应了解项目标准
- ⚠️ 缺少 PHASE_SUMMARY——gate 可能需要了解前序阶段的决策上下文

**上下文利用**：⚠️ 缺少 PHASE_SUMMARY。

**具体优化建议**：
1. 添加 conventions 加载步骤
2. 添加 PHASE_SUMMARY header 声明和条件加载

---

### 2.21 clean-code.md

**结构**：3 步流程（Read → Clean → Submit），委托给 clean-code skill。
**步骤合理性**：✅ 纯委托型模板。
**指令精确度**：✅ 清晰。

**冗余分析**：
- ⚠️ Step 3 的 submit 步骤与 task-executor 自动提交重复
- ⚠️ 缺少 PHASE_SUMMARY（但委托给 skill，由 skill 决定是否需要）

**具体优化建议**：
1. 移除 Step 3 的显式 submit 调用

---

### 2.22 fix-record-missed.md

**结构**：1 步流程（Verify Only），特殊模板。
**步骤合理性**：✅ verify-only 约束严格且适当。
**指令精确度**：
- ✅ 所有 quality gate 失败均 → blocked，符合 verify-only 语义
- ✅ 明确禁止重新实现

**上下文利用**：✅ 不需要 PHASE_SUMMARY（一次性验证任务）。

**具体优化建议**：无重大问题。

---

### 3. task-executor.md Agent 定义审查

**结构**：11 步执行协议，结构清晰。
**指令精确度**：
- ✅ Hard Constraints 使用 EXTREMELY-IMPORTANT 标签，层级正确
- ✅ "Follow every step in the synthesized strategy exactly" 明确了执行语义
- ✅ Step 6 "If you lose track..." 的策略恢复机制是好的容错设计

**潜在问题**：
- ⚠️ Step 8-9 的自动 submit + commit 与模板内嵌的 submit 步骤冲突（P0 问题 1.1）
- ⚠️ Step 7 检查 blocked 状态需要额外的 CLI 调用，增加延迟但保证了正确性

**具体优化建议**：
1. 在 Step 8 前添加注释说明"如果策略的最后一步已经是 submit-task，则跳过此步"（方案 B of 1.1）
2. 或者移除所有模板中的 submit 步骤，让 agent 统一管理（方案 A of 1.1）

---

### 4. prompt.go 合成逻辑审查

**结构**：Synthesize → renderTemplate → cleanTemplateOutput，流程清晰。
**指令精确度**：
- ✅ cleanTemplateOutput 处理了空值标签、条件语句清理、尾部空白等边缘情况
- ✅ Coverage 策略的优先级链设计合理（task frontmatter > config per-type > built-in default）
- ✅ PhaseDetect 的条件判断逻辑正确

**潜在问题**：
- ⚠️ `renderTemplate` 使用 `strings.ReplaceAll` 进行模板替换——如果模板内容中恰好包含 `{{TASK_ID}}` 等字符串（如在代码示例中），会被意外替换。目前没有 escaping 机制。
- ⚠️ `extractTestTypeArg` 只处理了 2 个 genScriptBases，但 type-to-template 映射有更多 test 类型 → 如果新增带 type suffix 的任务，需要同步更新此列表

**具体优化建议**：
1. 添加注释警告 `{{...}}` 占位符的 escaping 问题
2. 考虑将 genScriptBases 做成可配置项或从 task 常量中派生

---

## 4. 优化方案汇总（按优先级）

### P0 — 必须修复

| # | 问题 | 修复 | 影响文件 |
|---|------|------|---------|
| 1 | 双重提交 | 从 6 个模板中移除显式 submit 步骤 | doc.md, doc-eval.md, doc-summary.md, doc-consolidate.md, doc-drift.md, clean-code.md |

### P1 — 应该修复

| # | 问题 | 修复 | 影响文件 |
|---|------|------|---------|
| 2 | coding-fix 缺 PHASE_SUMMARY | 添加 header + Step 1 条件加载 | coding-fix.md |
| 3 | coding-fix conventions 加载顺序 | 移到 task 文件读取之前 | coding-fix.md |
| 4 | gate 缺 conventions 加载 | 添加 Step 1 规范加载 | gate.md |
| 5 | gate 缺 PHASE_SUMMARY | 添加 header + Step 1 条件加载 | gate.md |
| 6 | test.* HARD-RULE 标签命名 | 改为 TASK-CONSTRAINTS | test-gen-cases.md, test-eval-cases.md, test-gen-scripts.md, test-run.md, test-gen-and-run.md, test-graduate.md, test-verify-regression.md |
| 7 | 未使用的 FEATURE_SLUG | 从 doc.md, doc-eval.md 移除 | doc.md, doc-eval.md |
| 8 | 未使用的 SCOPE | 从 doc-consolidate.md, doc-drift.md 移除 | doc-consolidate.md, doc-drift.md |
| 9 | coding-refactor CODING_PRINCIPLES 格式 | 统一为 - 列表格式 | coding-refactor.md |

### P2 — 建议修复

| # | 问题 | 修复 | 影响文件 |
|---|------|------|---------|
| 10 | coding-feature TDD 引导不足 | Step 2 添加从 Acceptance Criteria 提取需求的引导 | coding-feature.md |
| 11 | Go 特定测试示例 | 改为语言无关描述 | coding-feature.md, coding-enhancement.md, coding-cleanup.md, coding-fix.md, coding-refactor.md |
| 12 | doc.md Step 2 过于模糊 | 添加按任务类型执行的引导 | doc.md |
| 13 | test-eval-cases MAIN_SESSION 声明 | 移除，由 dispatcher 控制 | test-eval-cases.md |

### P3 — 可选优化

| # | 问题 | 修复 | 影响文件 |
|---|------|------|---------|
| 14 | coding-refactor fmt 失败处理过于复杂 | 简化 git stash diff 为更简单的描述 | coding-refactor.md |
| 15 | coding-enhancement 与 coding-feature 高度重复 | 维持独立，添加同步维护注释 | coding-enhancement.md |
| 16 | prompt.go 占位符无 escaping | 添加注释警告 | prompt.go |

---

# 补充分析 A：每模板量化评分

## 评分标准（每维度 10 分）

| 维度 | 10 分 | 7 分 | 5 分 | 3 分 |
|------|-------|------|------|------|
| 结构清晰度 | 步骤逻辑清晰、层次分明、无歧义 | 结构清晰但有小瑕疵 | 结构可理解但不统一 | 结构混乱或有矛盾 |
| 步骤合理性 | 步骤顺序最优、无多余/缺失步骤 | 合理但有改进空间 | 基本可行但有明显问题 | 步骤有逻辑缺陷 |
| 指令精确度 | 每条指令具体、可执行、无二义 | 大部分精确，少量模糊 | 有多条模糊或二义指令 | 指令过于笼统 |
| 上下文利用 | 充分利用所有可用占位符 | 利用大部分，有小遗漏 | 利用不足 | 未利用或错误利用 |

## 评分表

### Coding 模板组

| 模板 | 结构 | 步骤 | 指令 | 上下文 | 总分 | 关键扣分原因 |
|------|------|------|------|--------|------|------------|
| coding-feature.md | 9 | 8 | 7 | 9 | 33/40 | 指令-2：TDD 引导笼统；Go 示例不通用 |
| coding-enhancement.md | 9 | 8 | 7 | 9 | 33/40 | 指令-2：与 feature 90% 重复，模糊相同；多出一行 review 指令但缺乏具体方法 |
| coding-cleanup.md | 9 | 9 | 8 | 9 | 35/40 | 指令-2：Go 示例不通用；结构简洁合理 |
| coding-refactor.md | 9 | 9 | 7 | 10 | 35/40 | 指令-2：fmt 失败处理复杂；CODING_PRINCIPLES 格式不一致；上下文利用最充分（增量编译策略） |
| coding-fix.md | 8 | 8 | 7 | 7 | 30/40 | 结构-1：conventions 加载顺序不一致；指令-2：E2E 处理精确但缺 TDD 引导；上下文-2：缺 PHASE_SUMMARY |

### Doc 模板组

| 模板 | 结构 | 步骤 | 指令 | 上下文 | 总分 | 关键扣分原因 |
|------|------|------|------|--------|------|------------|
| doc.md | 8 | 7 | 6 | 7 | 28/40 | 步骤-2：Self-Check 好但 Execute 过于模糊；指令-3：无具体执行引导；上下文-2：未用 FEATURE_SLUG、双重提交 |
| doc-eval.md | 9 | 9 | 8 | 7 | 33/40 | 指令-1：rubric 等权可能不够精确；上下文-2：未用 FEATURE_SLUG、双重提交 |
| doc-summary.md | 9 | 9 | 9 | 10 | 37/40 | 最佳 doc 模板——5 段结构精确，FEATURE_SLUG 正确使用 |
| doc-consolidate.md | 9 | 9 | 8 | 8 | 34/40 | 指令-1：非交互模式说明精确；上下文-1：声明未使用的 SCOPE；双重提交 |
| doc-drift.md | 9 | 9 | 8 | 8 | 34/40 | 与 consolidate 相同结构；上下文-1：声明未使用的 SCOPE；双重提交 |

### Test 模板组

| 模板 | 结构 | 步骤 | 指令 | 上下文 | 总分 | 关键扣分原因 |
|------|------|------|------|--------|------|------------|
| test-gen-cases.md | 9 | 9 | 9 | 9 | 36/40 | 命名冲突（HARD-RULE 标签）是唯一问题 |
| test-eval-cases.md | 8 | 8 | 7 | 9 | 32/40 | 结构-1：MAIN_SESSION 硬编码不应在模板中；指令-2：HARD-RULE 标签混淆 |
| test-gen-scripts.md | 9 | 9 | 9 | 9 | 36/40 | 仅命名冲突问题 |
| test-run.md | 9 | 9 | 9 | 9 | 36/40 | 仅命名冲突问题 |
| test-gen-and-run.md | 9 | 9 | 9 | 9 | 36/40 | 仅命名冲突问题；两阶段组合设计良好 |
| test-graduate.md | 9 | 9 | 9 | 9 | 36/40 | 仅命名冲突问题 |
| test-verify-regression.md | 9 | 9 | 9 | 9 | 36/40 | 禁止启动 dev server 和限制修复范围是精确的约束 |

### Validation + Gate + Special 模板组

| 模板 | 结构 | 步骤 | 指令 | 上下文 | 总分 | 关键扣分原因 |
|------|------|------|------|--------|------|------------|
| validation-code.md | 9 | 9 | 9 | 9 | 36/40 | 无重大问题；trivial/non-trivial 区分精确 |
| validation-ux.md | 9 | 9 | 9 | 9 | 36/40 | 无重大问题；合理省略 quality gate |
| gate.md | 8 | 8 | 8 | 7 | 31/40 | 上下文-2：缺 PHASE_SUMMARY；步骤-1：缺 conventions 加载；mermaid 流程图加分 |
| clean-code.md | 9 | 9 | 9 | 9 | 36/40 | 双重提交问题；其他良好 |
| fix-record-missed.md | 9 | 9 | 9 | 9 | 36/40 | 最佳特殊模板——verify-only 约束严格 |

### 基础设施

| 组件 | 结构 | 步骤 | 指令 | 上下文 | 总分 | 关键扣分原因 |
|------|------|------|------|--------|------|------------|
| task-executor.md | 9 | 9 | 8 | 9 | 35/40 | 指令-1：与模板内嵌 submit 冲突；策略恢复机制加分 |
| prompt.go | 9 | 9 | 8 | 9 | 35/40 | 指令-1：占位符无 escaping；cleanTemplateOutput 设计加分 |

## 评分分布

```
36-37 分（优秀）: test-gen-cases, test-gen-scripts, test-run, test-gen-and-run,
                  test-graduate, test-verify-regression, validation-code, validation-ux,
                  clean-code, fix-record-missed, doc-summary                    (11 个)
33-35 分（良好）: coding-feature, coding-enhancement, coding-cleanup, coding-refactor,
                  doc-eval, doc-consolidate, doc-drift, task-executor, prompt.go  (9 个)
30-32 分（合格）: coding-fix, test-eval-cases, gate                              (3 个)
28-29 分（待改进）: doc.md                                                       (1 个)
```

**关键洞察**：test.* 和 validation.* 模板得分最高（委托给 skill 的轻量 wrapper 设计成功），doc.md 得分最低（Execute 步骤过于模糊），coding-fix 和 gate 有具体缺失项待修复。

---

# 补充分析 B：指令可执行性分析

## 方法论

对每个模板中的每条关键指令进行三级评估：

- **可执行**：指令足够具体，agent 可以直接执行
- **模糊**：指令需要 agent 自行推断执行方式
- **歧义**：指令可能被 agent 多种解读，导致偏离

## 发现摘要

### B.1 高歧义指令（需修复）

#### B.1.1 doc.md Step 2 "Execute Document Work"

> "Perform the documentation work described in the task file:
> - Create new documents or modify existing ones as specified
> - Follow the project's existing documentation conventions and style
> - Ensure cross-references to other documents are accurate
> - Use consistent terminology throughout"

**歧义分析**：
1. "Follow the project's existing documentation conventions and style" — agent 如何知道项目的文档约定是什么？没有指向具体文件或加载机制
2. "Ensure cross-references are accurate" — 如何验证？是检查链接文件是否存在，还是检查内容一致性？
3. "Use consistent terminology" — 以哪个文档为基准？agent 需要一个术语表或参考文档

**修复建议**：Step 1 已加载 docs/conventions/（按修复建议后），Step 2 应引用这些加载的规范；添加具体验证方法（如 `ls` 检查引用文件是否存在）。

#### B.1.2 coding-feature.md Step 2 "Follow the TDD cycle for each requirement"

> "Follow the TDD cycle for each requirement:
> RED → Write failing test first
> GREEN → Implement minimal code to pass
> REFACTOR → Clean up while keeping tests green"

**歧义分析**：
1. "each requirement" — 从哪里获取 requirement 列表？task 文件可能有多个 acceptance criteria，但没有明确说"从 Acceptance Criteria 中逐个提取"
2. "Write failing test first" — 测试文件应放在哪里？命名规范是什么？是新建文件还是追加到已有测试文件？
3. "REFACTOR" — 在什么粒度上重构？每个 RED-GREEN 之后都重构，还是全部 GREEN 之后再重构？

**修复建议**：添加"从 task 文件的 Acceptance Criteria 中提取每个测试需求"的引导。测试文件位置由 agent 根据 conventions 决定（Step 1 已加载）。

#### B.1.3 coding-refactor.md Step 3 "Phase B failure recovery"

> "If partial migration has broken imports → git checkout the failed batch files and report as 'blocked at batch N'"

**歧义分析**：
1. `git checkout` 是 `git checkout -- <files>` 还是 `git checkout <branch>`？对于 agent 来说，这个命令的二义性可能导致错误操作
2. "the failed batch files" — 如果 batch 中的某些文件已成功迁移、某些失败，是否全部 checkout？

**修复建议**：改为 `git checkout -- <failed-files>`，明确指定文件级恢复。

### B.2 中等模糊指令（建议改进）

#### B.2.1 所有 coding 模板 Step 1 "Check docs/conventions/ and docs/business-rules/"

> "Read each file's YAML frontmatter `domains` field to determine relevance."

**模糊点**：如果 conventions 目录有 20+ 文件，逐个读取 frontmatter 效率低。但 agent 通常会用 Glob + Read 逐步检查，实际影响不大。

**建议**：维持现状，agent 的文件读取能力足以处理。

#### B.2.2 coding-refactor.md Step 2 "Dynamic coupling scan"

> "Reflection or metaprogramming calls referencing symbol names as strings"

**模糊点**：agent 如何系统性地扫描这些模式？是 grep 特定模式（如 `reflect.`、字符串拼接的类型名），还是需要更深入的分析？当前指令列出了要查找的模式但没有给出具体的扫描方法。

**建议**：添加示例 grep 命令（如 `grep -r "reflect\." ./affected/`）来引导 agent。

#### B.2.3 doc-eval.md 8 维度 rubric 评分粒度

> "125 = all present, well-organized; 100 = minor gaps; 75 = some sections missing; 50 = major sections missing; 0 = bare skeleton"

**模糊点**：评分粒度是 0/50/75/100/125（5 级），但 agent 可能给 90 分或 110 分——这些中间值没有定义。grader 是否应严格使用这 5 个值？

**建议**：添加"使用最近的定义值"或扩展为连续评分。

### B.3 可执行性良好的指令示例

以下指令设计精确，值得作为模板标准：

1. **coding-refactor.md Pre-check**：三条验证条件（git status clean、targeted tests pass、main branch warning）每条都有明确的不通过处理 → **标准**
2. **validation-code.md Step 2 trivial/non-trivial 区分**：明确界定 "trivial = missing import, typo" 且给出 max attempts → **标准**
3. **fix-record-missed.md 全部 quality gate 失败 → blocked**：无例外、无歧义 → **标准**
4. **test-verify-regression.md "Do NOT start dev server manually"**：直接禁止 + 解释原因（`just test-e2e` handles server lifecycle）→ **标准**
5. **coding-refactor.md "replacement order: longest identifier first"**：精确到操作顺序 → **标准**

---

# 补充分析 C：占位符生命周期分析

## 概述

prompt.go 的 `renderTemplate` 函数使用 `strings.ReplaceAll` 将 7 种占位符替换为运行时值。本节追踪每个占位符从声明→注入→消费→清理的完整路径。

## C.1 TASK_ID

| 阶段 | 详情 |
|------|------|
| **来源** | `opts.TaskID` → 从 `index.ByID()` 查找得到 |
| **注入** | 所有 22 个模板 header 第一行 `TASK_ID: {{TASK_ID}}` |
| **消费路径** | 1. task-executor agent 从 prompt 提取 → 执行 `forge prompt get-by-task-id <TASK_ID>` → 2. fix-record-missed 模板中 `forge task status {{TASK_ID}} blocked` |
| **清理** | 无需清理（始终有值） |
| **问题** | ⚠️ fix-record-missed.md 中 `{{TASK_ID}}` 出现在正文指令中（非 header），如果 ID 含特殊字符可能导致命令注入。但实际 task ID 格式为 `1.1`、`T-test-gen-cases` 等，风险极低。 |
| **建议** | 无需改动 |

## C.2 TASK_FILE

| 阶段 | 详情 |
|------|------|
| **来源** | `filepath.Join(projectRoot, feature.GetTaskFile(slug, t.File))` → 绝对路径 |
| **注入** | 所有 22 个模板 header `TASK_FILE: {{TASK_FILE}}` |
| **消费路径** | 1. 每个模板 Step 1 `read the task file at {{TASK_FILE}}` → 2. fix-record-missed 中验证 "Files Created/Modified" 部分列出的文件 |
| **清理** | 无需清理（始终有值） |
| **问题** | ✅ 生命周期完整，所有模板都正确引用 |
| **建议** | 无需改动 |

## C.3 SCOPE

| 阶段 | 详情 |
|------|------|
| **来源** | `t.Scope` → 如果为空或 "all" 则设为空串 |
| **注入** | 20 个模板 header `SCOPE: {{SCOPE}}` |
| **消费路径** | 1. 所有 coding.* 和 validation.* 模板的 just 命令 `just compile {{SCOPE}}` 等 → 2. coding-refactor 增量编译策略 `just compile {{SCOPE}}` |
| **清理** | `cleanTemplateOutput` 移除空值标签行 `SCOPE:` + 去除 just 命令尾部空白 `just compile ` → `just compile` |
| **问题** | ⚠️ 当 scope 为空时，`just compile ` 的尾部空格由 cleanTemplateOutput 正确处理。但如果 justfile 不支持空 scope 参数（即 `just compile` 等同于 `just compile ""`），可能导致错误。需确认 justfile 的 scope 参数处理。 |
| **未使用** | doc-consolidate.md, doc-drift.md 声明了 SCOPE 但未消费 |
| **未声明** | doc.md, doc-eval.md 不声明也不使用 → 正确 |
| **建议** | 移除 doc-consolidate.md 和 doc-drift.md 中未使用的 SCOPE 声明 |

## C.4 FEATURE_SLUG

| 阶段 | 详情 |
|------|------|
| **来源** | `opts.FeatureSlug` → 从命令参数传入 |
| **注入** | 3 个模板 header `FEATURE_SLUG: {{FEATURE_SLUG}}`：doc.md, doc-eval.md, doc-summary.md |
| **消费路径** | 仅 doc-summary.md Step 2：`docs/features/{{FEATURE_SLUG}}/tasks/records/` 构建路径 |
| **清理** | 无（始终有值） |
| **未消费** | doc.md, doc-eval.md 声明但从未在模板正文中使用 |
| **问题** | doc.md 和 doc-eval.md 是"声明但未使用"的死代码。虽然不影响运行（cleanTemplateOutput 不会移除有值的行），但增加维护混淆 |
| **建议** | 从 doc.md 和 doc-eval.md 移除 FEATURE_SLUG 声明 |

## C.5 PHASE_SUMMARY

| 阶段 | 详情 |
|------|------|
| **来源** | `PhaseDetect()` 函数 → 检查前序阶段的 summary 文件是否存在，返回相对路径或空串 |
| **注入** | 16 个模板 header `{{PHASE_SUMMARY}}` → 被替换为 `PHASE_SUMMARY: <path>` 或空串 |
| **消费路径** | 每个有 header 声明的模板 Step 1 中：`If {{PHASE_SUMMARY}} is non-empty, read that file for key decisions` |
| **清理** | `cleanTemplateOutput` 移除：1. 空值标签行 `PHASE_SUMMARY:` → 2. 条件句 `If `` is non-empty, ...` |
| **未声明** | coding-fix.md, doc-eval.md, gate.md, fix-record-missed.md, test-eval-cases.md, clean-code.md |
| **问题分析** | |
| | coding-fix：修复任务可能需要前序阶段上下文 → **应添加** |
| | gate：阶段门检查可能需要前序阶段决策 → **应添加** |
| | doc-eval：文档评估独立于阶段 → 不需要，**合理** |
| | fix-record-missed：一次性验证 → 不需要，**合理** |
| | test-eval-cases：测试评估独立 → 不需要，**合理** |
| | clean-code：委托给 skill → 不需要，**合理** |
| **清理正确性验证** | 当 PHASE_SUMMARY 为空时： |
| | 1. header 行 `PHASE_SUMMARY:` 被 `isLabelWithEmptyValue` 检测并移除 ✅ |
| | 2. 条件句 `If `` is non-empty, read that file for key decisions and conventions from the previous phase.` 被 `strings.Contains(trimmed, "If `` is non-empty")` 检测并移除 ✅ |
| | 3. 空行塌陷由 `Collapse consecutive blank lines` 处理 ✅ |
| **建议** | 为 coding-fix.md 和 gate.md 添加 PHASE_SUMMARY 声明 + Step 1 条件加载 |

## C.6 COVERAGE_STRATEGY / COVERAGE_TARGET

| 阶段 | 详情 |
|------|------|
| **来源** | `resolveCoverage()` → 优先级：task frontmatter `coverage` > config per-type > built-in default |
| **注入** | 仅在 `IsTestableType(t.Type)` 为 true 时注入（5 个 coding.* 模板） |
| **消费路径** | 嵌入在 `<IMPORTANT>` 标签中： |
| | - coding-feature/enhancement/fix：`{{COVERAGE_STRATEGY}} — 目标: {{COVERAGE_TARGET}}。达到目标后停止补充测试。` |
| | - coding-cleanup：`{{COVERAGE_STRATEGY}} — {{COVERAGE_TARGET}}。不新增测试，通过现有测试保持覆盖率。` |
| | - coding-refactor：`{{COVERAGE_STRATEGY}} — {{COVERAGE_TARGET}}。不新增测试，不追求高覆盖率。` |
| **清理** | 当类型非 testable 时，两个占位符替换为空串 → `<IMPORTANT>` 标签内可能产生 `覆盖率策略: — 目标: 。...` 格式的空指令行 |
| **关键问题** | ⚠️ cleanTemplateOutput 不移除 `<IMPORTANT>` 标签本身或其内部内容！当 COVERAGE_STRATEGY 和 COVERAGE_TARGET 为空时，IMPORTANT 块变为： |
| | ``` |
| | <IMPORTANT> |
| | 覆盖率策略: — 目标: 。达到目标后停止补充测试。 |
| | </IMPORTANT> |
| | ``` |
| | 这条指令仍会被 agent 看到，但内容为无意义的格式碎片。 |
| **验证**：prompt.go 的 `resolveCoverage` 中，非 testable 类型返回 `("", "")`。对于非 testable 类型，模板中本不应有 COVERAGE 占位符——但实际上，`typeToTemplate` 映射只将 coding.* 类型路由到含 COVERAGE 的模板。**非 coding 模板不含 COVERAGE 占位符**，所以此问题仅在理论上存在。✅ **实际安全** |
| **建议** | 无需改动。但如果未来有模板同时用于 testable 和 non-testable 类型，需要在 cleanTemplateOutput 中添加 IMPORTANT 块的空内容清理逻辑 |

## C.7 TEST_TYPE_ARG

| 阶段 | 详情 |
|------|------|
| **来源** | `extractTestTypeArg(t.ID)` → 从 task ID 的类型后缀提取（如 `T-test-gen-scripts-api` → ` --type api`） |
| **注入** | 2 个模板：test-gen-scripts.md 和 test-gen-and-run.md |
| **消费路径** | `Skill(skill="forge:gen-test-scripts"{{TEST_TYPE_ARG}})` → 替换为 `Skill(skill="forge:gen-test-scripts" --type api)` 或无后缀 |
| **清理** | 当无类型后缀时，替换为空串 → `Skill(skill="forge:gen-test-scripts")` → 正确 |
| **问题** | ⚠️ `extractTestTypeArg` 的 `genScriptBases` 列表硬编码为 `["T-test-gen-scripts", "T-quick-gen-and-run"]`，但 `typeToTemplate` 中 test-gen-and-run 的 key 是 `TypeTestGenAndRun`。需确认 `T-quick-gen-and-run` 和 `TypeTestGenAndRun` 对应的 task ID 格式是否一致 |
| **建议** | 添加注释说明 genScriptBases 与 task ID 格式的对应关系 |

## 占位符健康度总览

| 占位符 | 声明一致性 | 消费完整性 | 清理正确性 | 风险等级 |
|--------|-----------|-----------|-----------|---------|
| TASK_ID | ✅ 22/22 | ✅ 全部消费 | ✅ 无需清理 | 低 |
| TASK_FILE | ✅ 22/22 | ✅ 全部消费 | ✅ 无需清理 | 低 |
| SCOPE | ⚠️ 20/22，2 个未使用 | ✅ 消费正确 | ✅ cleanTemplateOutput 处理 | 低 |
| FEATURE_SLUG | ⚠️ 3/22，2 个未消费 | ⚠️ 1/3 实际使用 | ✅ 无需清理 | 低 |
| PHASE_SUMMARY | ⚠️ 16/22，2 个应添加 | ✅ 消费正确 | ✅ 三层清理（标签+条件句+空行） | 中 |
| COVERAGE_* | ✅ 5/22，仅 testable 类型 | ✅ 消费正确 | ✅ 类型路由隔离 | 低 |
| TEST_TYPE_ARG | ✅ 2/22，按需 | ✅ 消费正确 | ✅ 空串自然消除 | 低 |
