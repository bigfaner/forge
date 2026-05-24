---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: v3.0.0 Release Audit — Documentation-Implementation Drift Remediation

## Problem

Forge v3.0.0 的核心文档（README.md、ARCHITECTURE.md）与实际实现存在系统性偏差——包括过时的计数、不存在的组件、过时的命名方案、断裂的交叉引用和死代码。这些文档是用户和贡献者理解 Forge 的第一入口，当前状态会严重误导使用者。

### Evidence

系统性审计覆盖 5 个维度，发现 **27 个偏差项**：

| 维度 | Critical | Major | Minor | Advisory |
|------|----------|-------|-------|----------|
| README.md 事实性声明 | 7 | 6 | 1 | 0 |
| ARCHITECTURE.md 事实性声明 | 6 | 0 | 3 | 0 |
| CLI Reference 准确性 | 0 | 2 | 6 | 0 |
| Skill-CLI 交叉引用 | 2 | 2 | 1 | 0 |
| 架构健康度 | 2 | 3 | 4 | 5 |
| **合计** | **17** | **13** | **15** | **5** |

### Urgency

v3.0.0 是主版本发布，文档-实现一致性是发布质量的基本门槛。当前 README 版本号仍停留在 2.16.1，任务类型使用完全过时的命名——这些不是小修小补，是系统性的文档技术债。延迟修复意味着第一批 v3.0.0 用户将获得错误的项目理解。

## Proposed Solution

按优先级分层修复：Critical（发布阻塞）→ Major（发布前建议修复）→ Minor/Advisory（发布后迭代）。仅涉及文档更新和死代码清理，不修改任何运行时代码。

### Innovation Highlights

这是一次标准的技术债清理，无特殊创新。方法论上采用了"5维度交叉审计"——文档自检、文档-实现比对、接口契约验证、架构规范合规、依赖图分析——确保审计结果的完整性和可追溯性。

## Requirements Analysis

### Key Scenarios

- 用户阅读 README.md 快速理解 Forge 能力 → 当前会获得错误信息（版本号、技能数量、任务类型名称）
- 新贡献者阅读 ARCHITECTURE.md 理解系统架构 → 当前会寻找不存在的组件（4个 agent、PostToolUse hook、eval-harness）
- Agent 执行 skill 调用 CLI 命令 → 当前会在 gen-test-scripts 和 run-tests 中调用不存在的 CLI 命令
- Agent 加载 rules/ 文件 → 当前有 11 个 rules 文件因未在 SKILL.md 中引用而无法被发现

### Non-Functional Requirements

- 文档变更不能引入新的错误描述
- 死代码清理不能影响分发后的路径解析
- SKILL.md 拆分不能破坏已有的 rule/template 引用链

### Constraints & Dependencies

- 修改 `plugins/forge/` 下文件前必须遵守 forge-distribution.md 路径约束
- SKILL.md 拆分需遵守 skill-structure.md 的 350 行限制
- 跨技能依赖需遵守 skill-self-containment.md 的自包含原则

## Alternatives & Industry Benchmarking

### Industry Solutions

开源项目的发布审计通常通过 CHANGELOG + Breaking Changes 文档完成。Forge 的特殊性在于文档本身就是系统的一部分（被 agent 运行时读取），因此文档偏差等同于功能 bug。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 用户和 agent 都会被误导 | Rejected: v3.0.0 主版本必须保证文档质量 |
| 仅修复 README | 标准做法 | 快速 | ARCHITECTURE.md、交叉引用、架构问题被遗留 | Rejected: 问题系统性存在 |
| 仅修复 Critical 级别 | 优先级驱动 | 聚焦 | Major 级别的断裂引用会导致运行时错误 | Rejected: 断裂 CLI 引用也是运行时问题 |
| **分层全量修复** | 本提案 | 完整覆盖 | 工作量较大 | **Selected: 问题间存在依赖，分批修复不如一次性对齐** |

## Feasibility Assessment

### Technical Feasibility

所有修复均为文档编辑和文件清理，无运行时代码变更，风险可控。SKILL.md 拆分是最复杂的操作，但只需将现有内容移入 rules/ 文件。

### Resource & Timeline

- 文档更新（README + ARCHITECTURE + CLI reference）：~2h
- SKILL.md 拆分（gen-test-scripts + eval）：~1h
- CLI 交叉引用修复（6处断裂引用）：~30min
- Rules 引用补全 + 死代码清理：~1h
- 总计：~4.5h

### Dependency Readiness

无外部依赖。所有信息已在审计中收集完毕。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| README.md 基本准确 | Evidence Audit | Overturned: 版本号、技能数、Agent 数、任务类型、Pipeline ID、目录路径全部过时 |
| ARCHITECTURE.md 反映当前架构 | Evidence Audit | Overturned: 4 agents → 1, PostToolUse hook 不存在, doc-scorer/doc-reviser 不存在 |
| CLI Reference 完整准确 | Code Comparison | Refined: 命令结构完整，但 flags 细节有 8 处遗漏 |
| Skills 与 CLI 契约一致 | Cross-reference Audit | Overturned: 6 处断裂引用会导致运行时错误 |
| SKILL.md 遵守 350 行限制 | Compliance Check | Overturned: gen-test-scripts 527 行 (+50%), eval 488 行 (+39%) |

## Scope

### In Scope

**P0 — 发布阻塞（Critical，共 17 项）**

1. README.md 全面重写：版本号、技能/命令/Agent 计数、任务类型表（21种新命名）、Pipeline ID（语义化）、Go 版本、移除幽灵命令、修正目录路径
2. ARCHITECTURE.md 修正：Agent 架构（1 个专用 agent + general-purpose 模式）、Hook 系统（移除 PostToolUse）、Eval 系统（7 个 eval 命令）、目录路径（forge-cli/ 而非 task-cli/）
3. 断裂 CLI 引用修复：`forge config get surface` → `forge surfaces`（4 处）、`forge test run --tags regression` → 正确命令（4 处）
4. SKILL.md 超标拆分：gen-test-scripts 提取 Step 0.5/Step 1 到 rules/、eval 提取 freeform review pipeline 到 rules/
5. 缺失 rubric 文件：创建 `eval/rubrics/harness.md` 或在 SKILL.md 中添加异常处理

**P1 — 发布前建议修复（Major，共 13 项）**

6. CLI reference flags 补全：worktree push [slug]、worktree start --source-branch/--no-launch、feature -v、task query -v 等 8 处
7. 跨技能路径违规修复：run-tests 中 `skills/gen-journeys/rules/surface-<type>.md` 改为本地副本或描述性引用
8. 11 个未引用 rules 文件：在对应 SKILL.md 步骤中添加 Load 指令
9. 死代码清理：init-justfile 6 个 .just 模板文件（SKILL.md 明确说不使用）、gen-sitemap sitemap-example.json
10. guide.md 中过时的 Pipeline 描述更新（if applicable）
11. forge-distribution.md 中 Pipeline 描述与实际对齐

**P2 — 发布后迭代（Minor + Advisory，共 20 项）**

12. 模板变量命名风格统一（{{VAR}} vs <VAR>）
13. 模板 frontmatter 一致性
14. Hook Unix 端参数校验
15. validate-ux-pipeline.md 从 rubrics/ 移至 rules/
16. 未暴露 skill 的 command 入口评估
17. consolidate-specs SKILL.md 预留安全余量（346 行，距限制仅 4 行）
18. guide.md UTF-8 字符处理评估
19. CLI 测试中的过时 skill 路径引用更新
20. README 补充 v3.0.0 新特性：surface-based test profiles、worktree 管理、forensic 工具

### Out of Scope

- 运行时代码重构（Go CLI 代码变更）
- 新功能开发
- eval rubric 内容质量评估
- 性能优化
- 国际化

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| README 重写引入新的不准确描述 | M | M | 逐条与代码交叉验证，不改写未审计的部分 |
| SKILL.md 拆分破坏已有引用链 | L | H | 拆分前用 grep 确认所有引用，拆分后验证 |
| 死代码清理误删有用文件 | L | M | 仅清理已明确确认为死代码的文件 |
| harness rubric 创建不符合 eval 协议预期 | M | L | 参考现有 rubric 模板格式，或改为 SKILL.md 异常处理 |

## Success Criteria

- [ ] README.md 所有事实性声明（版本号、计数、路径、命令名）与代码库 100% 一致
- [ ] ARCHITECTURE.md 所有组件描述（agents、hooks、eval、目录）与代码库 100% 一致
- [ ] 零断裂 CLI 交叉引用（grep `forge config get surface` 和 `forge test run --tags` 在 skills/ 下无结果）
- [ ] 所有 SKILL.md 行数 ≤ 350
- [ ] 所有 rules/ 文件至少被其父 SKILL.md 引用一次
- [ ] init-justfile/templates/ 下的 .just 死代码文件已清理

## Next Steps

- Proceed to `/quick-tasks` to generate remediation tasks from this proposal
