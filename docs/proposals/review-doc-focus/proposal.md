---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Review-Doc 执行范围聚焦

## Problem

review-doc（`T-review-doc`）在执行过程中，agent 会读取并处理大量与审校无关的文件（任务 `.md` 文件、`records/` 目录、`index.json` 等），导致三个问题：

1. **Token 浪费**：任务文件和记录文件消耗 token 但不增加审校价值
2. **审校质量下降**：agent 注意力被非目标文档分散，对目标文档的审校深度不足
3. **Agent 误操作**：agent 有时修改任务文件或记录文件，而非只修改目标交付文档

### Evidence

- agent prompt（`prompt/data/doc-review.md` Step 1）明确要求 "scan tasks directory" 读取所有 doc 任务的 `.md` 文件
- Step 2 要求 "For each doc task found: Read the task's acceptance criteria from its .md file"
- autogen 模板（`task/data/doc-review.md`）的 Discovery Strategy 要求扫描整个 feature 目录，未排除 tasks/ 和 records/
- autogen 模板 Discovery Strategy 要求扫描 feature 全目录，按设计就会引用 tasks/ 和 records/ 下的文件（这是 prompt 逻辑推导的必然结果，非偶发现象）。注：当前缺少实际执行日志的量化数据（如某次 review-doc 引用的无关文件数、额外 token 消耗量），上述证据均为源码结构层面的推导，非运行时实测

### Urgency

review-doc-pipeline 已合并并投入使用，每次执行都产生上述问题。持续浪费 token 且审校质量受影响。

## Proposed Solution

两处核心改动：

1. **构建时嵌入 AC**：`forge task index` 生成 review-doc 任务时，从所有 doc 任务的 `.md` 文件中提取 Acceptance Criteria，汇总写入 `review-doc.md` 的专门区域。agent 执行时无需再读取其他任务文件。
2. **Allowlist 文档发现**：agent prompt 的文档发现从"扫描全部"改为 allowlist 策略——仅扫描 `docs/` 目录下的 `.md` 文件，从根本上排除 `tasks/`、`tasks/records/`、`manifest.md`、`index.json`。agent 只能读取和修改目标交付文档。

### Innovation Highlights

无特殊创新。核心思想是**关注点分离**——审校任务只关注交付物质量，不应涉及任务管理文件和执行记录。需要诚实指出：下文的跨领域类比（RAG pre-filtering、编译器 dead-code elimination）是方案确定后的事后类比验证（post-hoc justification），并非设计阶段的灵感来源。方案的真正灵感来源是 forge 自身的构建时管线架构——`forge task index` 已具备静态遍历任务文件的能力，扩展 AC 提取是自然延伸。

Challenge Override: Need Gate 的"更简方案"（纯 prompt 约束）在方案选择中被评为"Partial"而非完全否决——单独使用效果不确定，但作为结构层的配合手段是有效的。本方案的本质是**双层防护**：构建时从信息流层面移除无关数据（AC 嵌入），同时 prompt 层面用 allowlist 定义文档发现范围来增强可靠性。两者互补，缺一则不完整。**诚实声明**：Problem 3（agent 误修改任务文件）仅通过 prompt 层 allowlist 缓解，结构层无法阻止 agent 读写无关文件。该问题被部分解决（降低概率），非完全消除。彻底解决需要 task-executor 层的文件系统写权限沙箱，超出本方案范围。

## Requirements Analysis

### Key Scenarios

- **纯文档特性**：3 个 doc 任务 → `forge task index` 提取 3 组 AC 汇总到 review-doc.md → agent 只读 review-doc.md + 目标文档 → 核对 AC → 修复
- **混合特性**：doc + coding 任务 → 同上，review-doc 在测试流水线前执行
- **无 doc 任务的特性**：不生成 review-doc（现有逻辑不变）

### Non-Functional Requirements

- review-doc 执行 Referenced Documents 数量减少 50%+（以含 3+ doc 任务的典型特性为基准：当前扫描 tasks/ 全部文件，优化后仅读 review-doc.md + docs/ 下交付文档）
- AC 覆盖率不降低（所有 doc 任务的 AC 均被检查）
- 向后兼容：已有的 index.json 不受影响（仅影响 autogen 逻辑）

### Constraints & Dependencies

- 依赖 `forge task index` 的 auto-generation 逻辑
- AC 提取需能解析 doc 任务 `.md` 文件中的 `## Acceptance Criteria` section
- 标题匹配容错策略：精确匹配 `## Acceptance Criteria`，同时尝试 `## Acceptance criteria`（大小写差异）和 `## 验收标准`（中文别名）。匹配失败时输出构建警告，不静默跳过
- 不改变 task-executor 的调度机制
- **迁移约束**：`BuildIndex()` 仅在 `review-doc.md` 不存在时生成新文件。已存在的旧格式 review-doc.md 不会自动更新。用户需手动删除旧的 review-doc.md 后重新执行 `forge task index` 以获得含 AC Summary 的新格式

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **构建时嵌入 AC + Allowlist 发现（双层防护）** | 结构层（AC 嵌入）从信息流层面移除无关数据，prompt 层（allowlist）提供二次保障；即使 prompt 约束被 LLM 部分忽略，结构层已确保核心信息就位 | build 阶段需解析 AC，两个模板强耦合需同步修改 | **Selected** |
| 仅结构层（嵌入 AC，无 allowlist） | 消除 AC 扫描需求，改动集中 | agent 仍可能扫描 tasks/ 目录读取无关内容，误操作风险未解决 | Partial：结构层是必要非充分条件 |
| 仅 prompt 层（allowlist，无 AC 嵌入） | 改动最小，纯模板变更 | LLM 可能忽略 allowlist 约束；无法解决 AC 需运行时扫描的根本问题 | Partial：prompt 层单独使用效果不确定，需结构层配合 |
| **执行时 AC 提取 + 文件系统沙箱** | agent 调度时动态提取 AC（不依赖构建时），task-executor 通过 OS 级写权限限制 agent 只能写 `docs/` 目录 | 运行时提取仍需扫描任务文件；沙箱依赖 task-executor 实现，需改动调度器；跨平台权限管理复杂度高 | Rejected：改动范围超出 forge plugin 边界，且不解决 token 浪费（仍需读取任务文件提取 AC） |
| 不做改动 | 零成本 | 三大问题持续 | Rejected |

### Industry Context

本方案的核心思想与业界 agent context management 模式一致：

- **RAG context windowing**：检索增强生成（RAG）系统中，通过 pre-filtering 和 relevance scoring 在注入 context 前过滤无关文档。具体实现：LangChain v0.1+ 的 `VectorStoreRetriever`（基于 `similarity_search` top-k 截断）和 `ContextualCompressionRetriever`（基于 `llm_chain_extract` 过滤）均在 retrieval 阶段裁剪 context，而非将全部文档塞入 prompt。本方案的"构建时嵌入 AC"等价于 RAG 的 pre-filtering——在注入前就已排除了无关数据。
- **Anthropic context window optimization**：Anthropic 的 prompt engineering 指南强调"只给模型完成任务所需的最小上下文"（minimal sufficient context），避免 attention dilution。本方案将 allowlist 与结构化数据注入结合，正是这一原则的具体实践。
- **编译时优化 vs 运行时优化**：构建时提取 AC 的思路类比编译器的 dead-code elimination——在构建阶段确定哪些信息是可达的（reachable），运行时只携带可达信息。这比运行时动态过滤（denylist）更可靠。

## Feasibility Assessment

### Technical Feasibility

完全可行。改动集中在四条主线：

1. **BodyContext schema 扩展**（`build.go`）：新增 `DocTaskCriteria map[string]string` 字段（key=任务名称, value=AC markdown）。`renderBody()` 使用现有 `strings.ReplaceAll` 机制，在 Go 代码中将 map 序列化为 markdown 后替换新的 `{{DOC_TASK_AC}}` 占位符——无需切换到 Go `text/template`，保持与现有 12+ 模板一致的渲染策略。

2. **AC 提取管线**（`build.go`）：`BuildIndex()` 遍历 doc 任务文件，按 header 提取 AC section 内容，填充 BodyContext。

3. **autogen 模板**（`task/data/doc-review.md`）：新增 AC 汇总区域（格式见下文），更新 Discovery Strategy。

4. **agent prompt 模板**（`prompt/data/doc-review.md`）：结构化重构（规格见下文）。

### Resource & Timeline

- 预估 4-5 个 coding 任务：
  1. AC 提取管线（BodyContext schema 扩展 + `extract.go` 新增提取函数 + BuildIndex 集成）
  2. autogen 模板 AC 汇总区域 + Discovery Strategy 更新
  3. agent prompt 模板结构化重构（AC 加载 + allowlist 文档发现）
  4. 构建时 AC 验证（缺失检测 + 格式容错）
  5. 集成测试（端到端 review-doc 执行验证）

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| review-doc 需要运行时发现所有 doc 任务 | Assumption Flip | Overturned: doc 任务在 `forge task index` 时已确定，AC 可在构建时静态提取，无需运行时扫描 |
| agent 需要读取任务文件才能了解审校标准 | 5 Whys | Root cause: 当前 AC 存储在任务文件中，是因为没有其他传递机制。构建时嵌入消除了这个依赖 |

## Scope

### In Scope

- 修改 `build.go`：扩展 BodyContext schema（新增 `DocTaskCriteria map[string]string` 字段），`BuildIndex()` 调用新提取函数填充该字段，传入 review-doc 生成
- 修改 `extract.go`：新增 `extractDocTaskCriteria(taskDir string) map[string]string` 函数，遍历 doc 任务文件，按 header 提取 AC section 内容（含标题匹配容错），返回任务名称到 AC markdown 的映射。解析算法：逐行扫描 `.md` 文件，找到匹配 `## Acceptance Criteria`（含容错变体）的行后，收集该行之后、下一个 `## ` 开头行之前的所有行（含子列表、代码块、嵌套结构），作为 AC 原始内容返回。该算法与现有 `extractCheckboxItems()` 的结构化解析完全不同——后者仅匹配 `- [ ]` 行，新函数需保留 section 间的完整 markdown
- 修改 `autogen.go`：`renderBody()` 新增 `{{DOC_TASK_AC}}` 占位符，将 `DocTaskCriteria` 序列化为 markdown 后通过 `strings.ReplaceAll` 注入模板
- 修改 `task/data/doc-review.md`：autogen 模板增加 AC 汇总区域（格式见 AC 汇总区域格式 subsection），更新 Discovery Strategy 排除 tasks/ 和 records/
- 修改 `prompt/data/doc-review.md`：结构化重构（规格见 Prompt 模板重构规格 subsection）

### Out of Scope

- 其他任务类型的改动
- AC 在 doc 任务中的编写格式变更
- eval pipeline 相关改动
- record 模板改动（record 结构不变）

### AC 汇总区域格式

autogen 模板 `task/data/doc-review.md` 新增的 AC 汇总区域遵循以下 markdown schema：

```markdown
## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### [task-name-1]
<raw AC content from task .md file, preserved verbatim>

### [task-name-2]
<raw AC content from task .md file, preserved verbatim>
```

规则：
- 每个 doc 任务对应一个 `###` 子 section，标题为任务名称（即 `.md` 文件名去掉扩展名）
- section 内容为 `## Acceptance Criteria` 以下至下一个 `##` header 之间的原始 markdown
- 若任务文件中不存在 `## Acceptance Criteria` section，显示 `> No acceptance criteria defined.` 并输出构建警告
- 两个模板（autogen + agent prompt）的 AC 引用必须同步更新（见耦合约束）

### Prompt 模板重构规格

`prompt/data/doc-review.md` 重构为三步：

**Step 1: Load Pre-extracted AC** — 从 autogen 模板的 AC 汇总区域直接读取，移除原有的 "scan tasks directory" 指令。

**Step 2: Discover Target Documents** — allowlist 策略：仅遍历 `docs/` 子树下的 `.md` 文件，不扫描 `tasks/` 目录。

**Step 3: Review & Fix** — 对照 Step 1 的 AC 审校目标文档，仅修改 `docs/` 下文档。显式禁止修改 `tasks/` 和 `records/`。

重构步骤：删除 Step 1 扫描指令 → 新增 AC 加载指令 → 重写 Step 2 为 allowlist → Step 3 增加禁止修改约束。

## Coupling Constraints

autogen 模板（`task/data/doc-review.md`）和 agent prompt 模板（`prompt/data/doc-review.md`）存在强耦合：autogen 模板定义 AC 汇总区域的结构和位置，agent prompt 依赖该区域来定位和读取 AC 数据。**两个模板必须在同一个 coding task 中同步修改**，否则 agent 将无法正确获取 AC 或错误扫描无关文件。

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| AC 提取解析不稳定（doc 任务的 AC 格式不完全统一） | M | M | 提取逻辑基于 section header（`## Acceptance Criteria`），容忍内容格式差异。若 section 不存在：构建时输出 warning 日志 `"[WARN] task <name> has no Acceptance Criteria section"`，汇总区域显示 `> No acceptance criteria defined.`。agent 在 AC 为空时执行**自由审校模式**（不依赖 AC 对照，基于通用文档质量标准审校），并在执行记录中标注 `free-review` 标志。备选方案：将缺失 AC 设为硬失败（构建中断，要求补充 AC），但这会降低流程容错性，对快速迭代中的 doc 任务过于严格。当前选择 soft-fail + 警告，保留回旋余地 |
| 构建时提取的 AC 在执行时已过时（任务被手动修改） | L | M | review-doc 依赖最后业务任务完成后才执行，手动修改发生在执行前，构建时状态即为最新。但若用户在 `forge task index` 后、review-doc 执行前手动修改了 doc 任务的 AC，则提取内容与实际不同步。缓解：在 autogen 模板中注明 AC 的提取时间戳，若用户手动修改了 doc 任务需重新运行 `forge task index` |
| 过滤规则过于激进，排除合法内容文档 | L | M | 采用 allowlist 策略（仅扫描 `docs/` 下的 `.md` 文件），避免 denylist 的边缘情况漏判。若未来文档存放路径变化，需同步更新 allowlist 配置 |
| `renderBody()` 渲染策略与 AC 扩展不兼容 | M | H | 当前 `renderBody()` 使用 `strings.ReplaceAll` 而非 Go `text/template`，不支持 `{{range}}` 迭代。方案采用新占位符 `{{DOC_TASK_AC}}` + Go 代码序列化，保持与现有 12+ 模板一致的渲染策略，避免引入 `text/template` 依赖 |
| free-review 退化风险 | M | M | 当所有 doc 任务均缺少 AC 时，agent 完全基于通用质量标准审校，可能不如当前流程（至少能从任务文件获取上下文）。缓解：构建时对零 AC 特性输出 `"[WARN] feature has no AC for any doc task"` 警告，提示用户补充 AC 后重新 `forge task index` |
| AC 内容 prompt injection | L | H | AC 内容"保留原文"嵌入任务文件，若 doc 任务的 AC section 含恶意 agent 指令，将直接影响 review-doc agent 行为。本方案的信任模型假定所有 doc 任务 AC 内容均为同一作者/团队可控的善意内容（doc 任务为项目内部产物，非外部输入）。若未来支持外部贡献者的 doc 任务，需增加内容清洗步骤（过滤 `IGNORE PREVIOUS INSTRUCTIONS` 等模式），当前阶段不做清洗 |
| 生产环境回归无回退路径 | L | H | 若新流程审校质量不如预期，需能快速回退。缓解：AC 嵌入为纯增量变更（新增 `DocTaskCriteria` 字段 + 占位符），旧版 `renderBody()` 对未知占位符不报错（`strings.ReplaceAll` 无匹配时保持原样）。回退步骤（三个组件必须同步回退）：(1) 回退 Go 代码（移除 `DocTaskCriteria` 字段 + 占位符逻辑），(2) 回退 autogen 模板（移除 `{{DOC_TASK_AC}}` 占位符和 AC 汇总区域），(3) 回退 prompt 模板到旧版，(4) 删除 feature 目录下已有的 review-doc.md，(5) 重新 `forge task index`。注意：仅回退部分组件会导致残留占位符或格式不匹配 |

## Success Criteria

- [ ] `forge task index` 生成的 review-doc.md 包含 `## Acceptance Criteria Summary` 区域，每个 doc 任务对应一个 `###` 子 section
- [ ] review-doc agent 执行日志中 Referenced Documents 不包含任何 `tasks/` 路径下的文件
- [ ] review-doc agent 执行日志中 Referenced Documents 不包含 `records/`、`manifest.md`、`index.json`
- [ ] review-doc agent 只修改 `docs/` 下的目标交付文档
- [ ] 所有 doc 任务的 AC 均在 review-doc.md 的汇总区域中可查——通过 `BuildIndex()` 的构建断言验证：生成后立即检查 `DocTaskCriteria` 的 key 集合与 doc 任务列表完全匹配，缺失时输出 warning 并在汇总区域标注
- [ ] AC 为空的 doc 任务在汇总区域显示 `> No acceptance criteria defined.` 且 agent 执行记录标注 `free-review`
- [ ] 重构后的 prompt 模板 `prompt/data/doc-review.md` 满足三项结构要求：Step 1 引用 AC Summary section（无 "scan tasks directory" 指令）、Step 2 使用 docs/ allowlist、Step 3 包含 "仅修改 docs/ 下文件" 的显式禁止约束
- [ ] 审校质量回归验证：选取一个已完成 review-doc 的特性，用重构后的 pipeline 重新执行，对比两次审校记录中 AC 覆盖项的数量（新流程不应少于旧流程的 80%）。测量方法：人工统计旧流程执行记录中提及的 AC 项数（按 `## Acceptance Criteria` 下的条目计数）与新流程汇总区域的 AC 条目数对比；若旧记录格式不规范无法提取计数，则降级为定性比对（人工判定新流程是否遗漏关键 AC 项）

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
