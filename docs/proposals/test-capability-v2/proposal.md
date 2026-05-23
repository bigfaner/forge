---
created: 2026-05-23
author: "faner"
status: Draft
---

# Proposal: Test Capability 2.0 — Unified Pipeline + Deep Testing + Generality

## Problem

Forge 的测试管线存在三大结构性缺陷：两条并行的测试生成路径（gen-test-cases vs Journey-Contract）造成用户困惑和维护负担；生成的测试偏重 happy path，缺乏边界/异常/集成层级的深度覆盖；Convention 文件只内置了 Go/Vitest/Ginkgo，新框架接入成本高。

### Evidence

- `gen-test-cases` 与 `gen-journeys → gen-contracts → gen-test-scripts` 两条路径并存，都汇入 `gen-test-scripts`，用户不知道该走哪条
- 当前合约规范和测试脚本生成主要覆盖 happy path，边界值、异常输入、错误恢复、集成交互等场景需要手动补充
- Convention 文件仅覆盖 3 个框架（Go testing、Vitest、Ginkgo），Python/Java/Rust 等主流生态无内置支持
- Journey-Contract 路径上没有质量评测技能（journey 和 contract 阶段缺少 eval 门禁）
- gen-test-scripts/types/ 下各场景类型（CLI/TUI/UI/Mobile/API）的差异化仅停留在代码模板层面，策略和测试层级未做区分

### Urgency

v3.0.0 是重构测试架构的最佳窗口 — 正处于大版本分支上，可以做大范围变更而不破坏已发布版本。延迟的成本是：后续每新增一个框架或场景类型，都要在两条并行路径上分别适配，维护成本翻倍。

## Proposed Solution

将 Forge 测试能力升级为 2.0 架构：

1. **管线统一**：退休 gen-test-cases 及相关评测能力（eval-test-cases 命令、类型评测规则），Journey-Contract 成为唯一测试生成路径
2. **深度增强**：合约规范支持边界/异常场景自动衍生；基于风险等级差异化测试密度；支持集成/E2E 层级测试生成
3. **通用扩展**：内置主流框架 Convention 文件（pytest、JUnit、Rust/cargo test）；test-guide 智能化自动扫描生成 Convention 草稿
4. **评测补全**：新增 eval-journey 和 eval-contract 两个评测技能，在管线关键节点建立质量门禁
5. **场景差异化**：按场景类型（CLI/TUI/WebUI/Mobile UI/API）定义不同的测试策略和测试层级组合。每种场景独立测试，不搞跨场景组合编排

### Innovation Highlights

- **风险驱动密度**：高风险旅程自动生成更密集的测试矩阵（更多边界值、更多异常路径），低风险旅程只生成核心覆盖 — 区别于业界常见的"一刀切"测试密度
- **双维度场景差异化**：同时在策略层面（怎么测）和层级层面（测多深）做场景类型区分，而非仅停留在代码模板层面
- **Convention 自动生成**：test-guide 从项目文件信号自动推导 Convention 文件草稿，将新框架接入时间从"手动编写"降到"审核微调"

## Requirements Analysis

### Key Scenarios

- **Happy path**: 用户用 /quick 或 full pipeline 创建功能，管线自动从 PRD 生成 journey → contract → test scripts → 执行 → 报告
- **新项目冷启动**: 用户在无 Convention 文件的项目中首次运行测试生成，test-guide 自动检测框架并生成 Convention 草稿
- **高风险功能测试**: 用户为安全相关功能标记高风险，管线自动产出更密集的测试矩阵
- **多场景类型项目**: monorepo 包含 CLI + WebUI 等多种场景类型时，按场景类型独立生成测试，不做跨场景组合编排
- **评测门禁**: journey 和 contract 各有独立评测技能，低于阈值自动迭代修正

### Non-Functional Requirements

- Convention 文件扩展不应破坏已有 Convention 的加载逻辑
- 新增评测技能需复用现有 eval 框架（scorer-gate-revise 循环）
- 退休 gen-test-cases 不能影响已依赖该技能的功能（需确认无外部依赖）

### Constraints & Dependencies

- 必须遵循 Forge 的 skill/agent/command 分发模型（见 docs/conventions/forge-distribution.md）
- Convention 文件 schema 需保持向后兼容
- Journey-Contract 模型已在 `testing-journey-contract.md` 中定义，增强需基于此模型

### Per-Scenario Strategy Summary

基于测试专家深度分析，各场景类型的差异化策略如下：

| 维度 | CLI | TUI | WebUI | Mobile | API |
|------|-----|-----|-------|--------|-----|
| **交互模型** | 进程请求-响应 | 事件循环+键盘流 | DOM 事件驱动 | 触控+生命周期 | HTTP 请求-响应 |
| **执行方式** | subprocess 断言 | stdin pipe + ANSI 清洗 | 浏览器自动化(Playwright) | 设备/模拟器(Maestro) | HTTP 客户端断言 |
| **推荐入口层级** | E2E 优先 | 单元优先 | 均衡(集成+E2E) | 均衡 | 集成优先 |
| **层级占比** | 20/40/40 | 40/30/30 | 20/40/40 | 30/30/40 | 30/50/20 |
| **AI 生成难度** | 低 | 高 | 中 | 中 | 低 |
| **最大挑战** | 环境密闭性 | 终端仿真差异 | 选择器稳定性 | 设备/模拟器依赖 | 契约漂移 |
| **最适合的生成格式** | Go test / pytest | 框架特定(teatest) | Playwright / Cypress | Maestro YAML | Go test / supertest |
| **边界衍生能力** | 强(参数+前置条件) | 中(按键+异步) | 强(表单+路由+网络) | 中(权限+生命周期) | 极强(参数+认证+幂等) |

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Cucumber/Gherkin BDD**: 通过自然语言描述行为，自动映射到测试步骤。Forge 的 Journey 模型类似但更结构化
- **Postman/Newman API testing**: 契约测试模式。Forge 的 Contract 规范更通用，不限于 HTTP API
- **Playwright/TestProject E2E**: 端到端浏览器自动化。Forge 的 gen-test-scripts 已支持但仅作为模板层

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 双路径混乱持续，深度和通用性不足 | Rejected: v3.0.0 是最佳重构窗口 |
| 增量修补（只统一管线） | — | 最小改动范围 | 深度和通用性问题延后，补丁式方案 | Rejected: 三个维度互相依赖，拆开效果打折 |
| 完整 2.0（选定的方案） | — | 一次性解决结构性问题 | 工作量大，需要 full pipeline | **Selected: 结构性问题需要结构性解决方案** |

## Feasibility Assessment

### Technical Feasibility

- Journey-Contract 模型已建立，增强是在已有基础上扩展
- eval 框架已成熟（scorer-gate-revise），新增 eval-journey/eval-contract 是复用模式
- Convention 文件系统已有 schema 和验证机制，扩充内置库是增量的
- gen-test-scripts/types/ 已有分区结构，深化差异化是增强而非重建

### Resource & Timeline

- 需要 full pipeline（PRD → Design → Tasks），预计 20+ coding tasks
- 主要涉及 skill 文件编写和 skill 内部规则/reference 文件更新

### Dependency Readiness

- Journey-Contract 模型已稳定（testing-journey-contract.md 已发布）
- eval 框架和 rubric 模式已在 gen-test-cases eval 中验证
- test-guide 技能文件已存在，增强有基础

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "gen-test-cases 有独特价值不能删" | 5 Whys | Overturned: Journey-Contract 管线完全覆盖 gen-test-cases 的能力，且产出更结构化的中间产物（journey + contract） |
| "需要更多合约维度才能实现深度测试" | XY Detection | Overturned: 用户确认当前 6 维度足够，深度瓶颈在下游生成环节（边界/异常衍生、风险驱动密度）而非合约定义 |
| "所有场景类型用同一套测试策略就行" | Assumption Flip | Overturned: 不同场景类型天然需要不同的测试层级组合 — CLI 侧重 subprocess 断言，WebUI 侧重浏览器自动化，API 侧重 HTTP 契约验证 |
| "test.graduate 毕业机制需要保留" | 5 Whys | Overturned: 从未执行、无技能实现、有已知缺陷（毕业后不重跑测试），且与 Journey-Contract 管线的 `tests/<journey>/` 组织模型不兼容 — 新管线中测试按旅程存放，不需要"暂存→回归"晋升 |
| "需要跨场景组合编排" | Occam's Razor | Rejected: 用户明确每个项目独立场景类型，monorepo 按场景分别测试。跨场景编排增加复杂度但价值不明确 |

## Scope

### In Scope

- 退休 gen-test-cases 技能及相关评测能力（eval-test-cases 命令、test-cases 评测 rubric、类型子 rubric）
- 删除 test.graduate 任务类型和相关任务文件（从未执行、无技能实现、与 Journey-Contract 管线的按旅程组织模型不兼容）
- 新增 eval-journey 评测技能（含 rubric）
- 新增 eval-contract 评测技能（含 rubric）
- 合约规范增强：支持边界/异常场景自动衍生描述
- 风险驱动测试密度：高风险旅程生成更密集的测试矩阵
- 集成/E2E 层级测试生成支持
- 场景差异化：按 CLI/TUI/WebUI/Mobile/API 定义策略差异和层级差异
- 内置 Convention 文件扩充（pytest、JUnit、Rust/cargo test）
- test-guide 增强：自动扫描项目信号生成 Convention 草稿
- gen-test-scripts 适配增强后的合约规范和场景差异化
- 质量门禁更新以反映新管线

### Out of Scope

- 性能/负载测试
- 安全测试
- 视觉回归测试
- CI/CD 集成或测试环境管理
- 合约 6 维度模型修改
- 已使用 gen-test-cases 项目的迁移工具
- gen-test-scripts 的编译/lint 执行器变更（维持现有机制）
- 跨场景组合编排（每个项目/功能独立场景类型，monorepo 按场景分别测试）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 退休 gen-test-cases 影响已有工作流 | M | H | 先全面搜索确认无外部技能/agent 依赖 gen-test-cases 输出格式 |
| Convention 自动生成准确率不足 | M | M | test-guide 生成草稿后要求用户审核确认，不自动应用 |
| 场景差异化定义过于复杂导致维护负担 | M | M | 每种场景类型的策略和层级定义收敛到 1 个文件，避免碎片化 |
| 风险驱动密度的阈值难以定义 | H | M | 初始版本用简单的三级分类（高/中/低），后续可迭代精细化 |

## Success Criteria

- [ ] gen-test-cases 及所有相关文件（技能、命令、rubric、模板）已完全删除
- [ ] test.graduate 任务类型、任务文件、run-tasks/run-tests 中的引用已清理
- [ ] eval-journey 技能可用，对 journey 文档评分准确率 ≥ 850/1000
- [ ] eval-contract 技能可用，对 contract 文档评分准确率 ≥ 850/1000
- [ ] 高风险旅程自动生成的测试用例数量比低风险旅程多 ≥ 50%
- [ ] 新增 ≥ 3 个内置 Convention 文件（pytest、JUnit、Rust）
- [ ] test-guide 能从项目文件信号自动检测 ≥ 5 种测试框架
- [ ] CLI/TUI/WebUI/Mobile/API 各场景类型有独立的测试策略和层级定义文档
- [ ] 已有 Go/Vitest/Ginkgo Convention 文件的功能不受影响

## Next Steps

- Proceed to `/write-prd` to formalize requirements
