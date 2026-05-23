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
2. **深度增强**：合约规范支持边界/异常场景自动衍生；基于风险等级差异化测试密度；支持集成/E2E 层级测试生成（**不含单元测试** — 开发者在 feature 开发中已编写）
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
- **评测门禁**: journey 和 contract 各有独立评测技能，低于阈值自动迭代修正，最多迭代 3 轮；3 轮后仍未达阈值则暂停管线，输出当前评分和未通过项明细，由用户决定手动修正后继续或接受当前质量

### Non-Functional Requirements

- Convention 文件扩展不应破坏已有 Convention 的加载逻辑
- 新增评测技能需复用现有 eval 框架（scorer-gate-revise 循环）
- 退休 gen-test-cases 不能影响已依赖该技能的功能（需确认无外部依赖）

### Constraints & Dependencies

- 必须遵循 Forge 的 skill/agent/command 分发模型（见 docs/conventions/forge-distribution.md）
- Convention 文件 schema 需保持向后兼容
- Journey-Contract 模型已在 `testing-journey-contract.md` 中定义，增强需基于此模型

### Per-Scenario Strategy Summary

**定位**：Forge 测试管线只生成开发者手动编写成本高的复杂测试（集成测试、E2E 测试）。单元测试由开发者在 feature 开发过程中自行编写，不在管线范围内。

复杂测试映射为两个层级：
- **Contract 测试**（集成层）：验证单步 I/O 行为是否符合六维度契约
- **Journey 烟测试**（E2E 层）：端到端执行完整用户工作流，验证步骤间数据传递和状态一致性

各场景类型的支持级别和策略：

| 维度 | CLI | TUI | WebUI | Mobile | API |
|------|-----|-----|-------|--------|-----|
| **支持级别** | 核心 | 核心 | 核心 | **尽力而为** | 核心 |
| **AI Agent 适用性** | 4.0/5 | 3.0/5 | 3.0/5 | **2.0/5** | 3.5/5 |
| **交互模型** | 进程请求-响应 | 事件循环+键盘流 | DOM 事件驱动 | 触控+生命周期 | HTTP 请求-响应 |
| **执行方式** | subprocess 断言 | stdin pipe + ANSI 清洗 | 浏览器自动化(Playwright) | 设备/模拟器(Maestro) | HTTP 客户端断言 |
| **AI 优先侧重** | Contract 为主(80%) | Contract 为主(80%) | 平衡(50/50) | Journey 骨架 + deep link | 平衡(50/50) |
| **Contract 核心价值** | 单命令 I/O 验证 | 单步交互+异步 Cmd 验证 | 单页面行为验证 | 单屏幕交互(deep link 入口) | 单端点 I/O 契约验证 |
| **Journey 烟测试核心价值** | 命令链状态传播 | 连续键盘流累积效果 | 跨页面完整用户体验 | 跨屏幕移动端流程 | API 链路完整业务流程 |
| **必须衍生的边界 Outcome** | `not-found` + `already-exists` | `timeout`(每个异步 Cmd) | `validation-error` + `session-expired` | — | `unauthorized`(每个认证端点) |
| **最大 flakiness 来源** | 环境变量传染 | 异步渲染时序 | 网络/渲染时序 | App state 泄漏 | 服务启动时序 |
| **测试隔离难度** | 低(temp dir) | 低(stdin pipe) | 中(browser context) | 高(app lifecycle) | 中(API data cleanup) |
| **环境依赖复杂度** | 最低(编译工具链) | 低(编译+终端) | 中(浏览器+dev server) | 最高(模拟器/真机) | 中(API server+DB) |
| **最适合的生成格式** | Go test / pytest | 框架特定(teatest) | Playwright / Cypress | Maestro YAML | Go test / supertest |

**Mobile 场景的"尽力而为"策略**：
- 只生成 Maestro YAML 骨架（app lifecycle + navigation flow）和 deep link 测试
- 复杂场景（手势操作、屏幕内容验证、平台差异、权限矩阵）标记 `manual-only`
- 不追求高覆盖率，目标是降低开发者手动编写骨架的启动成本
- 降级理由：AI Agent 适用性 2/5（信息缺口最大、执行环境最复杂、flakiness 最高）

**风险驱动测试密度（复杂测试层级，3-5 步 Journey）**：

| 风险等级 | Contract 测试(每 Step) | Journey 烟测试 | 总测试数估算 |
|---------|----------------------|---------------|------------|
| High | 3-5 个 Outcome(含必须边界) | 1 个 happy path + 1 个失败路径 | 10-20 |
| Medium | 2-3 个 Outcome | 1 个 happy path | 7-13 |
| Low | 1-2 个 Outcome | 1 个 happy path | 4-8 |

**管线关键瓶颈**（优先级排序）：
1. **语义描述符→regex 转换断裂**：Fact Table 覆盖率直接决定测试质量（CLI/API 70-80%，WebUI/Mobile 40-50%）
2. **执行环境准备缺乏自动化**：完全依赖手动配置，无环境就绪检测
3. **Fact Table 覆盖率不足**：grep 只能获取静态信息，动态渲染/i18n/运行时输出无法获取
4. **失败诊断缺乏场景特定策略**
5. **测试数据管理缺乏场景特定策略**

**三个关键设计建议**：
1. **Run-to-Learn 机制**：生成骨架测试→运行捕获实际输出→丰富 Fact Table→重新生成精确测试，解决信息缺口根本问题
2. **场景特定执行环境就绪检测**：CLI 检查二进制、WebUI 探测 dev server、API 检查服务+DB 连通性
3. **置信度评级系统**：HIGH/MEDIUM/LOW 替代二元 pass/fail，区分可信赖测试和需审查测试

**gen-contracts 差异化建议**：基于场景类型差异化 Outcome 衍生策略，在 type 文件中增加 `required_outcomes` 声明。

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Cucumber/Gherkin BDD (v10+)**: 通过自然语言描述行为，自动映射到测试步骤。Forge 的 Journey 模型类似但更结构化
- **Postman/Newman API testing (Newman v6+, Postman Collections v2.1 schema)**: 通过 Collection Schema 定义 HTTP 请求契约（method/headers/body/status code），支持 pre-request scripts 和 test assertions 进行边界测试。局限：Contract 仅覆盖 HTTP 语义层，无法描述 CLI subprocess I/O、TUI 键盘流交互、Mobile deep link 等非 HTTP 场景。Forge 的 Contract 规范以六维度模型（Input/Output/Precondition/Postcondition/SideEffect/Constraint）描述单步行为，与传输协议无关，因此可统一覆盖 CLI/API/TUI/WebUI/Mobile 五种场景类型 — 这是 Newman 仅支持 HTTP 所不具备的
- **Playwright v1.40+ / TestProject (已停止维护, 2024)**: 端到端浏览器自动化。Forge 的 gen-test-scripts 已支持但仅作为模板层

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 双路径混乱持续，深度和通用性不足 | Rejected: v3.0.0 是最佳重构窗口 |
| 增量分期（先统一管线，再分阶段增强深度和通用性） | — | 改动可控，每期交付可验证；降低并行开发风险 | 管线统一后仍需深度增强和通用扩展才能体现完整价值，分期交付拉长整体周期；且管线统一会触及 gen-test-scripts 核心路径，后续深度增强仍需修改同一区域，存在二次变更成本 | Rejected: 可行但次优 — 管线统一（退休 gen-test-cases）和深度增强（边界衍生、风险驱动）在 gen-contracts/gen-test-scripts 中紧密耦合，拆期会导致同一文件二次重构。一次性解决避免了跨期返工 |
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
| "管线应该也生成单元测试" | XY Detection | Rejected: 单元测试由开发者在 feature 开发中编写（质量门禁已覆盖）。AI 管线的价值在于生成开发者手动编写成本高的复杂测试（集成/E2E），而非重复开发者已做的工作 |
| "Mobile 应该和其他场景同等投入" | ROI 分析 | Rejected: Mobile AI Agent 适用性 2/5（信息缺口最大、执行环境最复杂、flakiness 最高），ROI 不足以支撑核心级别投入。降级为"尽力而为" |

## Scope

### In Scope

- 退休 gen-test-cases 技能及相关评测能力（eval-test-cases 命令、test-cases 评测 rubric、类型子 rubric）
- 删除 test.graduate 任务类型和相关任务文件（从未执行、无技能实现、与 Journey-Contract 管线的按旅程组织模型不兼容）
- 新增 eval-journey 评测技能（含 rubric）
- 新增 eval-contract 评测技能（含 rubric）
- 合约规范增强：支持边界/异常场景自动衍生描述（排除纯函数逻辑分支等单元测试级别边界）。技术路线：在 gen-contracts 技能中采用 LLM prompt 增强策略 — 根据 Contract 的 Input 维度字段类型（string/number/enum）注入类型特定边界提示词（如 string → empty/overflow/unicode；number → zero/min/max；enum → invalid_value），结合项目 Convention 中声明的场景类型 `required_outcomes` 规则（如 API 场景必须衍生 `unauthorized`），由 LLM 在生成合约时一并产出边界/异常 Outcome，而非依赖后置模板匹配或单独的规则推导引擎
- 风险驱动测试密度：高风险旅程生成更密集的测试矩阵。风险标记机制：复用 testing-journey-contract.md 中已定义的 Risk 字段（journey 级别的 risk_level: high/medium/low），由用户在编写 journey 文档时手动标记（或通过 gen-journeys 根据 PRD 中的安全/合规关键词自动建议 risk_level），管线在 gen-contracts 阶段读取该字段并按密度表差异化衍生 Outcome 数量
- Contract 测试（集成层）+ Journey 烟测试（E2E 层）生成，按场景类型差异化侧重比例
- 场景差异化：CLI/TUI/WebUI/API 核心支持 + Mobile 尽力而为，按场景类型定义策略差异、层级侧重、和必须衍生的边界 Outcome
- 内置 Convention 文件扩充（pytest、JUnit、Rust/cargo test）
- test-guide 增强：自动扫描项目信号检测测试框架并生成 Convention 草稿
- gen-test-scripts 适配增强后的合约规范和场景差异化
- Run-to-Learn 机制：生成骨架测试 → 运行捕获实际输出 → 丰富 Fact Table → 重新生成精确测试，作为管线内置的迭代增强环节
- 场景特定执行环境就绪检测：CLI 场景检查目标二进制存在性和可执行权限；WebUI 场景探测 dev server 端口响应；API 场景验证服务进程存活 + 数据库连通性。检测不通过时输出具体缺失项和修复建议，而非静默失败
- 置信度评级系统：为每个生成的测试赋予 HIGH/MEDIUM/LOW 置信度标签 — HIGH 表示 Fact Table 覆盖率充足且执行环境验证通过；MEDIUM 表示 Fact Table 部分缺失但核心路径可验证；LOW 表示信息缺口大或环境未就绪，结果需人工审查。评级结果随测试报告一同输出，供用户快速筛选可信结果
- 质量门禁更新以反映新管线

### Out of Scope

- 单元测试（开发者在 feature 开发中已编写，Forge 管线不生成）
- 性能/负载测试
- 安全测试
- 视觉回归测试
- CI/CD 集成或测试环境管理
- 合约 6 维度模型修改
- 已使用 gen-test-cases 项目的迁移工具
- gen-test-scripts 的编译/lint 执行器变更（维持现有机制）
- 跨场景组合编排（每个项目/功能独立场景类型，monorepo 按场景分别测试）
- 执行环境自动准备与配置（瓶颈 #2 已通过"场景特定执行环境就绪检测"部分缓解 — 检测环境是否就绪并给出修复建议；但全自动配置涉及各场景的依赖安装、服务启停等运维操作，超出测试生成管线职责，延迟至后续版本）
- 失败诊断场景特定策略（瓶颈 #4：需要在各场景类型中积累足够失败模式数据才能定义有效策略，当前版本依赖通用 scorer-gate-revise 循环 + LLM 诊断，延迟至有实际失败数据积累后再场景化）
- 测试数据管理场景特定策略（瓶颈 #5：涉及数据工厂、数据清理、数据隔离等系统性问题，独立于测试生成逻辑，延迟至后续版本专项解决）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 退休 gen-test-cases 影响已有工作流 | M | H | 先全面搜索确认无外部技能/agent 依赖 gen-test-cases 输出格式 |
| Convention 自动生成准确率不足 | M | M | test-guide 生成草稿后要求用户审核确认，不自动应用 |
| 场景差异化定义过于复杂导致维护负担 | M | M | 每种场景类型的策略和层级定义收敛到 1 个文件，避免碎片化 |
| 风险驱动密度的阈值难以定义 | H | M | 初始版本用简单的三级分类（高/中/低），后续可迭代精细化 |
| 20+ tasks 大范围变更导致管线整体退化 | H | H | 分三阶段交付，每阶段设置明确门禁标准：(1) 管线统一 — 门禁：对 2+ 个已有 Journey-Contract 项目跑完整管线，所有 journey/contract/test-script 生成无报错且测试可执行；未通过则修复后重跑，不进入下一阶段；(2) 深度增强 — 门禁：高风险旅程生成测试数 ≥ 低风险旅程 ×1.5，边界 Outcome 中无效比例 < 20%；(3) 通用扩展 — 门禁：新增 Convention 文件在对应框架真实项目上生成可执行测试 ≥ 3 个。每阶段门禁失败：修复时限 2 个工作日，超时则回退阶段变更并升级讨论范围是否裁剪 |
| LLM prompt 增强对边界值推导的准确性不足 | M | H | 边界衍生依赖 LLM 对类型提示词的推理质量，可能产生无效或遗漏的边界 Outcome。缓解：(1) 每个场景类型的 `required_outcomes` 作为硬约束兜底，确保关键边界不遗漏；(2) eval-contract 评测中增加边界 Outcome 准确率维度（抽取 20+ 个边界 Outcome，人工判定有效/无效，准确率目标 ≥ 80%）；(3) 无效边界 Outcome 不阻塞管线，仅标记为 `low-confidence` 供人工审查 |

## Success Criteria

- [ ] gen-test-cases 及所有相关文件（技能、命令、rubric、模板）已完全删除
- [ ] test.graduate 任务类型、任务文件、run-tasks/run-tests 中的引用已清理
- [ ] eval-journey 技能可用，对 journey 文档评分准确率 ≥ 850/1000（度量方法：人工标注 10+ 份 journey 样本作为 gold standard，eval-journey 自动评分与人工评分的 Pearson 相关系数 ≥ 0.85；inter-rater reliability 要求两位评审者对同一样本的评分差异 ≤ 100 分）
- [ ] eval-contract 技能可用，对 contract 文档评分准确率 ≥ 850/1000（度量方法：同上，人工标注 gold standard + Pearson 相关系数 ≥ 0.85）
- [ ] 高风险旅程自动生成的测试用例数量比低风险旅程多 ≥ 50%
- [ ] 新增 ≥ 3 个内置 Convention 文件（pytest、JUnit、Rust）
- [ ] test-guide 能从项目文件信号自动检测 ≥ 5 种测试框架，并能为检测到的框架生成 Convention 草稿且满足：草稿包含全部 4 个必需 section（framework、discovery、structure、assertions）；草稿通过 Convention schema 验证（无 schema 错误）；用户审核微调后（修改量 ≤ 草稿总内容的 20%）即可作为正式 Convention 使用
- [ ] Run-to-Learn 机制可端到端执行：给定一个初始 Fact Table 覆盖率 < 60% 的项目，经过 ≤ 3 轮 Run-to-Learn 迭代后，Fact Table 覆盖率提升 ≥ 20 个百分点，且重新生成的测试中边界/异常 Outcome 占比 ≥ 30%（与不使用 Run-to-Learn 的基线对比）
- [ ] CLI/TUI/WebUI/Mobile/API 各场景类型有独立的测试策略和层级定义文档
- [ ] 已有 Go/Vitest/Ginkgo Convention 文件的功能不受影响

## Next Steps

- Proceed to `/write-prd` to formalize requirements
