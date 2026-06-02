---
created: "2026-06-02"
author: "fanhuifeng"
status: Draft
intent: "cleanup"
---

# Proposal: 全局文档-代码一致性审计与知识库清理

## Problem

Forge 项目的用户文档、规范文档和知识库与实际代码实现之间存在未量化的不一致：文档描述的行为与代码实际行为矛盾、规范文档引用了已废弃的实现路径、知识库中积累了大量可能过时或重复的条目。这些不一致会误导 AI 代理执行错误操作，增加新成员上手成本。

### Evidence

- 已有 5 个局部审计提案（plugin-consistency-audit、skill-ecosystem-audit、skill-instruction-audit、prompt-template-audit、test-pipeline-consistency-audit）发现了不同层面的不一致，但均未执行
- `test-pipeline-consistency-audit` 已确认 Go 代码仍使用旧术语（`tests/e2e/`、`E2E`、`graduation`、`staging`），而文档层已迁移到新术语（`tests/<journey>/`、tag-based promotion）
- docs/lessons/ 下积累了 133 条经验教训，未经过系统性有效性审查
- docs/conventions/ 下 16+ 份规范文档，部分可能描述了已不存在的代码结构

### Urgency

项目正处于 v3.0.0 分支开发阶段，文档-代码不一致会随版本迭代持续恶化。越晚清理，积累的错误越多，AI 代理基于过时文档执行的风险越高。

## Proposed Solution

执行一次性的三层系统性审计，产出结构化问题报告并转化为可执行 Task：

1. **L1 用户文档层**：审计 README.md、ARCHITECTURE.md、DESIGN.md、docs/user-guide/、docs/official-references/ 与代码实际行为的一致性
2. **L2 规范文档层**：审计 docs/business-rules/、docs/conventions/、docs/reference/ 与实际实现的一致性
3. **L3 知识库层**：审查 docs/lessons/（133条）和 docs/decisions/（10条）的有效性，标记过时、重复、无参考价值的条目

每层审计检查三个维度：过时/错误（文档与代码矛盾）、缺失（代码做了但文档没说，或反过来）、冗余（重复或无价值的内容）。

### Innovation Highlights

无特殊创新——这是标准的文档审计实践。亮点在于利用 AI 代理的代码理解能力进行自动化交叉比对，减少人工审查成本。

## Requirements Analysis

### Key Scenarios

- **S1**：AI 代理读取 ARCHITECTURE.md 了解系统架构时，文档描述的流程与实际 hook 执行顺序一致
- **S2**：AI 代理读取 docs/conventions/naming.md 时，命名规则与 CLI 代码中的实际常量一致
- **S3**：新成员查阅 docs/lessons/ 时，所有条目仍有参考价值，不存在指向已删除代码的教训
- **S4**：清理后，知识库条目数量显著减少，每条剩余条目都有明确的适用范围

### Non-Functional Requirements

- 审计结果必须包含文件路径和行号，可直接定位问题
- 每个问题必须标注严重级别（P0/P1/P2/P3），便于优先级排序
- 生成的 Task 必须可由 task-executor 独立执行

### Constraints & Dependencies

- 审计基于 v3.0.0 分支当前代码状态
- 不修改任何代码或文档，只生成报告和 Task
- 知识库清理需人工确认，不可自动删除

## Alternatives & Industry Benchmarking

### Industry Solutions

文档-代码一致性是经典问题。常见做法：
- **Doc-as-code**：文档与代码同仓库，通过 CI 检查一致性
- **Automated linting**：用工具（如 liche、markdown-link-check）检查链接有效性
- **Periodic audit**：定期人工或半自动审查

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 不一致持续恶化，AI 代理风险增加 | Rejected: v3.0.0 发布前必须清理 |
| 执行现有5个提案的86个task | 现有提案 | 已有任务定义 | 覆盖不完整（缺知识库、缺用户文档层），部分任务可能过时 | Rejected: 不够全面 |
| 增强现有工具 | /consolidate-specs | 可重复执行 | 开发成本高，当前需快速解决 | Rejected: 本质是一次性工作 |
| **分层系统性审计** | 本提案 | 覆盖完整，可直接执行 | 工作量较大 | **Selected: 最符合当前需求** |

## Feasibility Assessment

### Technical Feasibility

完全可行。AI 代理已具备代码阅读和交叉比对能力，无需额外工具开发。

### Resource & Timeline

- L1 用户文档层：约 15-20 文件，预计 2-3 个 Task
- L2 规范文档层：约 21 文件，预计 3-5 个 Task
- L3 知识库层：约 143 条目，预计 3-5 个 Task
- 总计约 8-13 个 Task

### Dependency Readiness

无外部依赖。所有审计目标文件均在项目仓库内。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "现有5个提案可以覆盖需求" | XY Detection | Overturned: 现有提案只覆盖 plugin 层和 test pipeline，不覆盖用户文档、知识库 |
| "133条 lessons 都有参考价值" | Assumption Flip | Refined: 大量 lessons 可能已过时，需要逐条验证 |
| "审计需要开发新工具" | Occam's Razor | Overturned: AI 代理直接审计即可，无需开发自动化工具 |

## Scope

### In Scope

- L1: README.md、ARCHITECTURE.md、DESIGN.md、docs/user-guide/、docs/official-references/ 与代码的一致性审计
- L2: docs/business-rules/、docs/conventions/、docs/reference/ 与实现的一致性审计
- L3: docs/lessons/（133条）和 docs/decisions/（10条）的有效性审查
- 审计产出：结构化问题报告（文件路径+行号+严重级别+建议动作）
- 将问题报告转化为可执行 Task

### Out of Scope

- docs/features/（183个 feature 目录）的清理
- docs/proposals/（204个 proposal）的清理
- Plugin skill/command 内部一致性（已有 plugin-consistency-audit 覆盖）
- CLI Go 代码内部质量改进
- 测试代码审查
- 自动修复或自动删除

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 审计范围过大，Task 数量失控 | M | M | 每层独立审计，严格控制每条 Task 的粒度 |
| 知识库条目有效性判断主观 | M | L | 对 lessons 和 decisions 只标记建议，由人工最终确认 |
| 审计期间代码持续变化导致结果过时 | L | L | 审计基于 v3.0.0 分支快照，生成 Task 时标注审计基准 commit |
| 误删有价值的知识库条目 | L | H | 知识库清理 Task 必须标记为需人工确认，不自动执行 |

## Success Criteria

- [ ] L1 层：所有用户文档文件完成审计，发现的不一致问题 100% 记录在报告中
- [ ] L2 层：所有规范文档文件完成审计，发现的不一致问题 100% 记录在报告中
- [ ] L3 层：133 条 lessons 和 10 条 decisions 完成逐条有效性审查，每条标记为"有效"/"过时"/"重复"/"需更新"
- [ ] 审计报告中的每个问题包含：文件路径、行号范围、严重级别（P0-P3）、建议动作
- [ ] 所有问题已转化为可执行 Task，每个 Task 可由 task-executor 独立执行
- [ ] 知识库清理相关的 Task 均标注为"需人工确认"

## Next Steps

- Proceed to `/quick-tasks` to generate tasks directly from this proposal (skip PRD and design for cleanup work)
