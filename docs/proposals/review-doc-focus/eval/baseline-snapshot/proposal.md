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
- 实际执行记录显示 Referenced Documents 包含任务文件和记录文件

### Urgency

review-doc-pipeline 已合并并投入使用，每次执行都产生上述问题。持续浪费 token 且审校质量受影响。

## Proposed Solution

两处核心改动：

1. **构建时嵌入 AC**：`forge task index` 生成 review-doc 任务时，从所有 doc 任务的 `.md` 文件中提取 Acceptance Criteria，汇总写入 `review-doc.md` 的专门区域。agent 执行时无需再读取其他任务文件。
2. **过滤式发现策略**：agent prompt 的文档发现从"扫描全部"改为"只扫描内容文档"，明确排除 `tasks/`、`tasks/records/`、`manifest.md`、`index.json`。agent 只能读取和修改目标交付文档。

### Innovation Highlights

无特殊创新。核心思想是**关注点分离**——审校任务只关注交付物质量，不应涉及任务管理文件和执行记录。

Challenge Override: Need Gate 的"更简方案"（纯 prompt 约束）被用户否决。Reason: LLM 经常忽略 prompt 约束，问题根源需从信息流层面解决。

## Requirements Analysis

### Key Scenarios

- **纯文档特性**：3 个 doc 任务 → `forge task index` 提取 3 组 AC 汇总到 review-doc.md → agent 只读 review-doc.md + 目标文档 → 核对 AC → 修复
- **混合特性**：doc + coding 任务 → 同上，review-doc 在测试流水线前执行
- **无 doc 任务的特性**：不生成 review-doc（现有逻辑不变）

### Non-Functional Requirements

- review-doc 执行 token 消耗降低（减少无关文件读取）
- AC 覆盖率不降低（所有 doc 任务的 AC 均被检查）
- 向后兼容：已有的 index.json 不受影响（仅影响 autogen 逻辑）

### Constraints & Dependencies

- 依赖 `forge task index` 的 auto-generation 逻辑
- AC 提取需能解析 doc 任务 `.md` 文件中的 `## Acceptance Criteria` section
- 不改变 task-executor 的调度机制

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **构建时嵌入 AC + 过滤式发现** | 从根源消除问题，信息流清晰 | build 阶段需解析 AC，复杂度略增 | **Selected** |
| 纯 prompt 约束 | 改动最小 | LLM 常忽略约束，根源未解 | Rejected: 效果不确定 |
| 不做改动 | 零成本 | 三大问题持续 | Rejected |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动集中在三处：
- Go CLI `autogen.go`：`GetReviewDocTask()` 接受 AC 数据
- Go CLI `build.go`：`BuildIndex()` 提取 AC 并传入 review-doc 生成
- 模板文件：`task/data/doc-review.md` 和 `prompt/data/doc-review.md` 简化

### Resource & Timeline

- 预估 2-3 个 coding 任务
- 涉及 Go CLI 修改 + 模板文件修改

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| review-doc 需要运行时发现所有 doc 任务 | Assumption Flip | Overturned: doc 任务在 `forge task index` 时已确定，AC 可在构建时静态提取，无需运行时扫描 |
| agent 需要读取任务文件才能了解审校标准 | 5 Whys | Root cause: 当前 AC 存储在任务文件中，是因为没有其他传递机制。构建时嵌入消除了这个依赖 |

## Scope

### In Scope

- 修改 `autogen.go`：review-doc 任务定义支持嵌入 AC 数据
- 修改 `build.go`：`BuildIndex()` 从 doc 任务文件提取 AC 并传入 review-doc 生成
- 修改 `task/data/doc-review.md`：autogen 模板增加 AC 汇总区域，更新 Discovery Strategy 排除 tasks/ 和 records/
- 修改 `prompt/data/doc-review.md`：简化 agent prompt——移除任务扫描指令，改用内嵌 AC，加入文档过滤规则，禁止修改任务文件和记录

### Out of Scope

- 其他任务类型的改动
- AC 在 doc 任务中的编写格式变更
- eval pipeline 相关改动
- record 模板改动（record 结构不变）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| AC 提取解析不稳定（doc 任务的 AC 格式不完全统一） | M | M | 提取逻辑基于 section header（`## Acceptance Criteria`），容忍内容格式差异；若 section 不存在则标记为空 |
| 构建时提取的 AC 在执行时已过时（任务被手动修改） | L | L | review-doc 依赖最后业务任务完成后才执行，手动修改发生在执行前，构建时状态即为最新 |
| 过滤规则过于激进，排除合法内容文档 | L | M | 过滤基于明确路径模式（tasks/、records/、manifest.md、index.json），不影响 docs/ 下的内容文件 |

## Success Criteria

- [ ] review-doc agent 执行时不再读取 tasks/ 目录中的其他任务文件
- [ ] review-doc agent 执行时不再读取 records/ 目录中的记录文件
- [ ] review-doc agent 只读取和修改目标交付文档
- [ ] AC 覆盖率不降低（所有 doc 任务的 AC 均在 review-doc.md 中可查）
- [ ] autogen 模板包含 AC 汇总区域和过滤式 Discovery Strategy

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
