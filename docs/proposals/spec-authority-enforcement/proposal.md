---
created: 2026-05-23
author: faner
status: Draft
---

# Proposal: 规范权威性强制执行——从模板和任务生成层面防止 Agent 偏离

## Problem

Agent 执行任务时以现有代码为参照而非以权威规范文档（tech-design.md 等）为参照，导致大规模偏离。在 test-capability-v2 特性中产生了 43 处与 tech-design.md 不一致的偏差。

### Evidence

- `docs/lessons/gotcha-spec-authority-drift.md` 记录了完整的根因分析和 5 级溯源
- Level 0 根因：coding.* 模板的 Step 1 只说"read the task file"，未要求读 Reference Files
- Level 3 过程原因：缺少"自顶向下验证"步骤——反应式局部修复而非主动式全局审计
- 探索发现：quick-tasks 模板硬编码 Reference Files 为单一 proposal.md；breakdown-tasks 对非 UI 任务没有 Reference Files 填充指引

### Urgency

这是系统性缺陷——每次涉及规范驱动修改的任务都可能重蹈覆辙。教训已提取，修复成本极低（纯文档修改），拖延无正当理由。

## Proposed Solution

两层防护：

1. **Agent 层（task-executor.md）**：在执行协议中增加两个强制步骤——"声明 Reference Files 为权威来源，按需加载"和"AC 逐条验收"
2. **任务生成层（quick-tasks / breakdown-tasks）**：改进 Reference Files 生成质量，要求精确到 section 或行号范围

### Innovation Highlights

非创新性改进，而是将已验证的工程实践制度化。关键洞察来自教训文档的 Level 4 分析：LLM agent 天然倾向局部一致性而非全局一致性，因此强制措施必须嵌入执行流程本身，而非依赖 agent 的自觉性。

## Requirements Analysis

### Key Scenarios

- **场景 1：coding.* 任务执行**：agent 收到合成 prompt → 声明 Reference Files 列出的文档为权威来源 → 按需加载 tech-design.md 中与当前实现相关的 section → 以规范为权威而非现有代码 → 实现完成后对照 AC 逐条验收
- **场景 2：doc.* 任务执行**：同上，但验收重点是文档结构合规而非路径命名
- **场景 3：quick-tasks 生成任务**：从 proposal 提取关键技术约束，将相关 design 文档的精确 section 写入 Reference Files
- **场景 4：breakdown-tasks 生成任务**：从 tech-design.md 的架构决策中提取每个任务相关的 section，精确引用而非笼统指向整个文件

### Non-Functional Requirements

- Reference Files 精确引用不应导致 prompt 过长——每个任务引用 2-5 个 section，而非整个文件
- 模板改动不改变现有任务的执行流程结构，只是在已有步骤间插入新步骤

### Constraints & Dependencies

- 修改 `plugins/forge/` 下的文件前必须遵循 `docs/conventions/forge-distribution.md`
- 纯 Markdown 文档修改，不涉及 Go 代码变更

## Alternatives & Industry Benchmarking

### Industry Solutions

业界解决 LLM agent 规范遵循问题的常见手段：RAG 检索、prompt 模板强制引用、structured output 约束。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 问题必复发 | Rejected: 教训已证明代价高昂 |
| CLI 层 auto-inline | 自研 | 最可靠，agent 无法跳过 | 需改 Go 代码，增加 prompt 长度 | Rejected: 用户选择纯文档方案 |
| Prompt 模板 + Agent 协议（本方案） | Prompt Engineering 最佳实践 | 零代码改动，覆盖所有任务类型 | 依赖 agent 遵守 `<EXTREMELY-IMPORTANT>` 标记 | **Selected: 最小有效改动** |

## Feasibility Assessment

### Technical Feasibility

纯 Markdown 编辑，无技术风险。

### Resource & Timeline

预计 4-6 个文件修改（task-executor.md + 审计后确认的模板 + quick-tasks/breakdown-tasks skill），1-2 小时完成。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "模板层改 4 个文件就够了" | 事实核查（Explore agent 审计全部 19 个模板） | Refined: 需审计全部模板，实际需修改数量待定 |
| "quick-tasks/breakdown-tasks 已经在 Reference Files 中填入设计文档" | 事实核查（读取模板和 SKILL.md） | Overturned: quick-tasks 硬编码只有 proposal.md；breakdown-tasks 对非 UI 任务无填充指引 |
| "Reference Files 只需列文件路径" | XY Detection — 用户实际需要的是精确到 section 的引用 | Refined: 精确引用减少 agent 注意力分散，提高规范遵循率 |

## Scope

### In Scope

- 修改 `plugins/forge/agents/task-executor.md`：增加"声明 Reference Files 权威性并按需加载"步骤和"AC 逐条验收"步骤
- 审计全部 19 个 `forge-cli/pkg/prompt/data/*.md` 模板，确定哪些需要 Reference Files 权威性声明
- 对需要强化的模板，在 Step 1 加 `<EXTREMELY-IMPORTANT>` Reference Files 声明
- 改进 `plugins/forge/skills/quick-tasks/`：Reference Files 从 proposal 和相关文档中提取精确 section 引用
- 改进 `plugins/forge/skills/breakdown-tasks/`：为非 UI 任务增加 Reference Files 填充指引，要求精确到 section

### Out of Scope

- 修改 `forge-cli` Go 代码（auto-inline Reference Files）
- 修改任务文件格式（index.json schema）
- 添加 hooks 或运行时验证机制
- 改变 `forge prompt get-by-task-id` 的合成逻辑

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent 忽略 `<EXTREMELY-IMPORTANT>` 标记 | M | H | 标记放在执行流程的关键位置（Step 1 和 submit 前），并声明"按需加载"降低遵从成本 |
| Reference Files section 引用过时（design 文档更新后行号偏移） | M | M | 引用 section 标题而非行号，行号仅作为辅助定位 |
| breakdown-tasks 生成任务时 Reference Files 填充不完整 | L | M | 在 SKILL.md 中加显式 checklist：每个任务必须有 ≥1 个 design-level Reference File |

## Success Criteria

- [ ] task-executor.md 的执行协议包含"声明 Reference Files 权威性并按需加载"步骤（Step 5 前）和"AC 逐条验收"步骤（Step 8 前）
- [ ] 全部 19 个模板完成审计，输出审计报告（哪些模板需要强化 + 理由）
- [ ] 需强化的模板 Step 1 包含 `<EXTREMELY-IMPORTANT>` Reference Files 权威性声明
- [ ] quick-tasks 生成的 coding 任务 Reference Files 包含 ≥1 个精确 section 引用（非仅 proposal.md）
- [ ] breakdown-tasks SKILL.md 包含非 UI 任务的 Reference Files 填充规则
- [ ] 所有 Reference Files 条目格式统一：`path/to/file.md#section — 简要说明该 section 定义了什么`

## Next Steps

- Proceed to `/quick-tasks` to generate and execute tasks
