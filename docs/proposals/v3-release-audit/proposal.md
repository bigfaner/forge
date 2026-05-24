---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: v3.0.0 Release Audit — Documentation-Implementation Drift Remediation

## Problem

Forge v3.0.0 的核心文档（README.md、ARCHITECTURE.md）与实际实现存在系统性偏差——过时的计数、不存在的组件、断裂的交叉引用和死代码。skills/ 下的 SKILL.md 亦存在超标（+50% 行数）、断裂 CLI 引用和孤儿 rules 文件。这些文档是用户和 agent 运行时的第一入口，当前状态会严重误导使用者和 agent。根因：v2→v3 快速迭代中缺乏文档同步机制。

### Evidence

系统性审计覆盖 5 维度，发现 **50 个偏差项**（事实性声明=版本号、组件计数、文件路径、命令名及参数、API 签名、配置键名；freeform review 额外发现幽灵命令和子系统遗漏）。测量方法：文档自检→实现-比对→契约验证→规范合规→依赖图。

严重等级：**Critical**=运行时阻断；**Major**=人工误导；**Minor**=外观格式；**Advisory**=风格建议。

| 维度 | Critical | Major | Minor | Advisory |
|------|----------|-------|-------|----------|
| README.md 事实性声明 | 7 | 6 | 1 | 0 |
| ARCHITECTURE.md 事实性声明 | 6 | 0 | 3 | 0 |
| CLI Reference 准确性 | 0 | 2 | 6 | 0 |
| Skill-CLI 交叉引用 | 2 | 2 | 1 | 0 |
| 架构健康度 | 2 | 3 | 4 | 5 |
| **合计** | **17** | **13** | **15** | **5** |

### Urgency

v3.0.0 是主版本发布，已有 3 个 beta 用户。当前至少 5 个 open issues 直接源于文档-实现偏差（错误 CLI 命令导致运行时失败）。README 版本号仍为 2.16.1（实际 3.0.0-rc.24），任务类型使用过时命名——系统性文档技术债。

## Proposed Solution

按优先级分层修复：Critical→Major→Minor/Advisory。文档更新、死代码清理和 SKILL.md 重组，不修改 Go 运行时代码（SKILL.md 拆分除外，见 P0.4）。

### Target State

- README.md 3.0.0：计数与 `forge --help` 一致，dot-notation 新命名，目录匹配。结构：标题→安装（Go 版本）→快速开始→命令速查（与 `forge --help` 一一对应）→技能列表（计数匹配 `ls skills/`）→任务类型表（dot-notation 全覆盖）→架构→贡献→License
- ARCHITECTURE.md：1 agent，无 PostToolUse，1000 分制
- 零断裂 CLI 引用

### Innovation Highlights

"5维度交叉审计"借鉴财务审计重要性阈值（materiality：按金额分级，低于阈值可忽略）和供应链可追溯性（有向图标记孤儿），形成可复用框架——按运行时影响设定严重等级，<Major 降级为发布后处理。

## Requirements Analysis

### Key Scenarios

- 用户阅读 README → 错误版本号、技能数、任务类型
- 贡献者阅读 ARCHITECTURE.md → 寻找不存在组件（4 agents、PostToolUse）
- Agent 调用 CLI → gen-test-scripts/run-tests 调用不存在命令
- Agent 加载 rules/ → 15 个 rules 未被 SKILL.md 引用

### Non-Functional Requirements

- 文档变更不引入新错误
- 死代码清理不影响分发路径解析
- SKILL.md 拆分不破坏 rule/template 引用链
- **可维护性**：漂移预防（如 README 计数与 `ls skills/ | wc -l` 断言）
- **可回归性**：P0 后端到端回归（agent 加载、eval/gen-test-scripts、CLI 调用链）

### Constraints & Dependencies

- 修改 `plugins/forge/` 需遵守 forge-distribution.md
- SKILL.md 需遵守 skill-structure.md 350 行限制
- 跨技能依赖需遵守 skill-self-containment.md

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Kubernetes** 自动从 API 生成参考文档，CI 验证
- **Rust** RFC 要求功能变更附带文档更新
- **Spring Boot** CI 运行 docusaurus validated links
- **remark-lint** 等 CI 工具检测断链/过时引用

Forge 的特殊性：文档是 agent 运行时输入，偏差=功能 bug。

**P2 Follow-up: 自动化文档生成**（Kubernetes/Hugo 模式）：Go 源码自动提取 CLI/flag/版本号生成 Markdown 片段。~4-8h。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | 多数小项目 | 零成本 | Forge 文档=运行时输入，偏差=bug | Rejected |
| 仅修复 README | 标准做法 | 快速 | ARCHITECTURE.md、交叉引用被遗留 | Rejected: 问题系统性 |
| 仅修复 Critical | 优先级驱动 | 聚焦 | Major 断裂引用导致运行时错误 | Rejected: 断裂 CLI 也是运行时问题 |
| CI 文档验证自动化 | K8s/Spring Boot | 长期预防漂移 | 需 ~4-8h 编写脚本（out of scope），无 CI 基础设施 | Deferred: 纳入 NFR 漂移预防 |
| **分层全量修复** | 本提案 | 完整覆盖 | 工作量较大 | **Selected: 依赖关系需一次性对齐** |

## Feasibility Assessment

### Technical Feasibility

所有修复为文档编辑和文件清理，无运行时代码变更。SKILL.md 拆分最复杂。

### Resource & Timeline

- 文档更新（README + ARCHITECTURE + CLI）：~2h
- SKILL.md 拆分：~3h（eval 含 7 rules + Mermaid + 三路分支）
- CLI 交叉引用修复（5 处）：~30min
- Rules 补全 + 死代码清理：~1h
- ARCHITECTURE.md 子系统概述（P1.12）：~2h
- Mermaid 同步更新（P0.4 附带）：~30min
- 总计：~9.5h

### Dependency Readiness

无外部依赖。信息已在审计中收集。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| README.md 基本准确 | Evidence Audit | Overturned: 版本号/计数/任务类型/Pipeline ID/目录全部过时；freeform 发现幽灵命令 |
| ARCHITECTURE.md 反映架构 | Evidence Audit | Overturned: 4 agents→1, PostToolUse 不存在；freeform 发现子系统遗漏 |
| CLI Reference 完整 | Code Comparison | Refined: flags 遗漏 8 处 |
| Skills-CLI 契约一致 | Cross-reference | Overturned: 5 断裂 CLI + 4 孤儿 rules 引用 |
| SKILL.md ≤350 行 | Compliance Check | Overturned: gen-test-scripts 527 行, eval 488 行 |

## Scope

### In Scope

**P0 — 发布阻塞（Critical，17 项）。顺序：P0.4（SKILL.md 拆分）→ P0.5（harness 决策）→ P0.2（ARCHITECTURE.md）→ P0.3（CLI 引用）→ P0.1（README 重写最后，依赖前面计数）。通用回滚原则：feature branch + go/no-go checkpoint。P0.4 回滚：git stash，回归失败则回滚并降级 P1。**

1. README.md 全面重写：版本号、计数、任务类型表（dot-notation）、Pipeline ID、Go 版本、移除幽灵命令/web 引用、修正路径
2. ARCHITECTURE.md 修正（P0=已有事实错误；P1=缺失新增）：Agent、Hook、Eval、路径修正
3. CLI 引用修复：`forge config get surface` → `forge surfaces`（4处）、`test.execution` → 正确命令
4. SKILL.md 超标拆分（**最先执行**）：gen-test-scripts 提取 Step 0.5/1 到 rules/、eval 提取 freeform pipeline 到 rules/。**改变 agent 行为，但超标（+50%/+39%）违反强制约束。**风险：(a) 需显式 Load；(b) 三路分支一致性；(c) Mermaid 同步
5. 缺失 rubric 文件决策（**第二执行**）：(a) 创建 harness.md（超范围）；(b) 异常处理（运行时变更）；(c) 从 Prerequisites 表移除 harness 类型，保留文件 → **推荐 (c)**

**P1 — 发布前建议（Major 13 项）**

6. CLI flags 补全：worktree、feature -v、task query -v 等 8 处
7. 跨技能路径违规：run-tests 中硬编码路径改为描述性引用
8. 孤儿 rules 修复：6 个真孤儿添加 Load；5 个参数化 surface rules 标注引用；4 处 `forge test run --tags` 降级归入此项
9. 死代码清理：gen-sitemap sitemap-example.json；init-justfile 6 个 .just 模板评估
10. guide.md 过时 Pipeline 描述更新
11. forge-distribution.md Pipeline 描述对齐
12. ARCHITECTURE.md 补充 9 个 v3.0.0 子系统概述（概述+架构角色+SKILL.md 链接，≤180 行）：surface detection、worktree、Convention、forensic、deep-research、clean-code、extract-design-md、test-guide、learn

**P2 — 发布后迭代（Minor+Advisory 20 项）**

13. 模板变量命名统一
14. 模板 frontmatter 一致性
15. Hook Unix 端参数校验
16. validate-ux-pipeline.md 从 rubrics/ 移至 rules/
17. 未暴露 skill 的 command 入口评估
18. consolidate-specs SKILL.md 提取 ≥50 行
19. guide.md UTF-8 字符处理评估
20. CLI 测试中过时 skill 路径引用更新
21. README 补充 v3.0.0 新特性

### Out of Scope

- 运行时代码重构
- 新功能开发
- eval rubric 质量评估
- 性能优化
- 国际化

## Key Risks

| Risk | L | I | Mitigation |
|------|---|---|------------|
| README 重写引入新错误 | M | M | 逐条交叉验证，不改未审计部分 |
| SKILL.md 拆分破坏引用链 | H | H | 拆分前 grep 确认；git stash 回滚 |
| 死代码误删 | L | M | 通过 grep 全仓库确认无引用后删除 |
| harness rubric 不合规 | L | L | 方案 (c) 纯删除 |
| P0 串行依赖级联失败 | H | H | P0.4 回滚已定义；P0.2/0.3 可并行 |
| P1.12 范围蔓延 | M | M | 每子系统 ≤20 行，≤180 行 |
| 后续迭代漂移 | H | M | NFR 检查点 + CI 验证（deferred） |

## Success Criteria

**P0 验收（发布阻塞门禁）：**
- [ ] README 事实性声明 100% 一致：版本号/计数/路径/命令名与代码库比对零差异
- [ ] ARCHITECTURE.md 已有内容 100% 一致
- [ ] 零断裂 CLI 引用（grep 无 `forge config get surface` / `test.execution` 结果）
- [ ] 所有 SKILL.md ≤ 350 行
- [ ] 所有 rules/ 被父 SKILL.md 引用（入度 ≥ 1）
- [ ] P0.4 agent 回归：eval/gen-test-scripts 无报错
- [ ] `forge init` 路径完整性通过
- [ ] P0.5 harness 类型已从 Prerequisites 移除

**P1 验收（发布前建议）：**
- [ ] CLI flags 补全，`--help` 与文档一致
- [ ] 零跨技能路径违规
- [ ] 零孤儿 rules（入度 = 0 为 0）
- [ ] init-justfile 模板决策已记录
- [ ] ARCHITECTURE.md 含 9 个子系统概述
- [ ] P1.10 guide.md Pipeline 描述与实现一致
- [ ] P1.11 forge-distribution.md Pipeline 描述与实现一致
- [ ] 漂移预防：≥3 自动化断言（skill 计数、task type 计数、CLI 命令覆盖）

## Next Steps

- Proceed to `/quick-tasks` 生成修复任务
