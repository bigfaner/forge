---
created: "2026-05-26"
author: "fanhuifeng"
status: Draft
---

# Proposal: Surface-Specific Test Type Model（Surface 测试类型模型）

## Problem

Forge 把所有生成的测试统称为 "e2e 测试"，但不同 Surface 的测试在语义、执行模型和验证维度上存在本质差异。CLI 测试启动子进程验证输出，API 测试发送 HTTP 请求验证响应，Web 测试驱动浏览器验证交互——它们不是同一类测试，却共享同一个标签。

### Evidence

1. **目录结构**：测试代码统一放入 `tests/<journey>/`（部分项目用 `tests/e2e/`），不区分 surface 类型
2. **Justfile recipe**：使用 `just test` 或 `just test-e2e` 执行所有测试，无 surface 区分
3. **Task 类型名**：`test.gen-scripts`、`test.run` 等任务类型不携带 surface 信息
4. **文档混用**：ARCHITECTURE.md 中 "高级测试" 和 "e2e 测试" 交替使用，概念边界模糊
5. **质量门**：`just unit-test`（开发者自写）和 `just test`（Forge 生成）两层划分，但第二层缺乏测试类型语义
6. **gen-test-scripts 的 types/ 目录**：已有 CLI/API/TUI/Web/Mobile 五种生成策略，但生成的测试代码在命名和标签上不体现这一差异

### Urgency

随着 Forge 支持的项目类型增多，"e2e" 标签的不精确性会导致：
- 用户误解测试覆盖范围（以为 CLI 测试是端到端的）
- 新 skill/规则文件编写时术语不一致
- 与外部工具集成时无法准确描述测试类型

## Proposed Solution

取消 "e2e 测试" 作为统一标签，建立 **Surface → Test Type** 的映射模型。每种 surface 有自己的测试类型名称、语义定义和验证维度。

### Test Type Mapping

| Surface | Test Type（EN） | 测试类型（CN） | 验证维度 | 执行模型 |
|---------|-----------------|---------------|---------|---------|
| `cli` | CLI Integration Test | CLI 集成测试 | 退出码 + stdout/stderr | 子进程执行 |
| `tui` | Terminal Integration Test | 终端集成测试 | 终端输出 + 交互行为 | 子进程 + stdin pipe |
| `api` | API Contract Test | API 契约测试 | HTTP 状态码 + 响应体 + Header | HTTP 客户端 |
| `web` | Web E2E Test | Web 端到端测试 | UI 渲染 + 用户交互 + 状态流转 | 浏览器自动化 |
| `mobile` | Mobile UI Test | 移动端 UI 测试 | UI 渲染 + 用户交互 | Maestro YAML / 手动验证 |

### 语义定义

- **CLI 集成测试**：编译独立二进制，通过子进程调用，验证命令行参数解析、输出格式、退出码、错误处理。黑盒视角，不测试内部函数。
- **终端集成测试**：编译独立二进制，通过 stdin pipe 模拟用户输入，验证终端渲染输出（ANSI 序列处理、布局、异步 Cmd 响应）。半黑盒视角。
- **API 契约测试**：启动 HTTP 服务器（或使用测试服务器），发送请求，验证响应符合 Contract 定义的六个维度。黑盒视角，Contract 即契约。
- **Web 端到端测试**：启动 dev server，通过浏览器自动化（Playwright）模拟用户操作，验证 UI 渲染、交互逻辑、跨页面状态流转。真正的端到端。
- **移动端 UI 测试**：通过 Maestro YAML 定义操作序列，验证移动端 UI 行为。Best-effort 模式，部分场景标记为 manual-only。

### Innovation Highlights

这不是创造新概念，而是**精确命名已有实践**。Forge 的 gen-test-scripts 已经按 surface 类型分化了生成策略（`types/cli.md`、`types/api.md` 等），run-tests 已经分化了编排序列，但这些分化一直没有上升到概念层。本提案将这些隐含的分化显式化、命名化。

## Requirements Analysis

### Key Scenarios

1. **概念查询**：用户查阅文档时，能快速找到自己项目 surface 对应的测试类型定义
2. **测试生成**：gen-test-scripts 生成的测试代码文件名/注释中体现测试类型（而非统一的 "e2e"）
3. **测试执行**：justfile recipe 按测试类型命名（如 `test-cli-integration`、`test-api-contract`），而非统一的 `test-e2e`
4. **任务追踪**：index.json 中 test 任务的类型名携带 surface 信息（如 `test.gen-scripts.cli`）
5. **质量门**：质量门报告区分不同测试类型的执行结果

### Constraints & Dependencies

- 现有 skill 的 `types/` 和 `rules/surfaces/` 文件已按 surface 分化，是本提案的基础
- `forge surfaces` CLI 已提供 surface 检测能力
- Justfile recipe 命名需与 init-justfile skill 的 surface 规则同步更新
- 任务类型名的变更需与 task-lifecycle business rule 中的保留类型列表协调

## Alternatives & Industry Benchmarking

### Industry Solutions

- **Go 社区**：使用 build tags 区分测试类型（`//go:build integration`、`//go:build e2e`）
- **Spring Boot**：`@Tag("integration")`、`@Tag("e2e")` 注解区分测试类型
- **Playwright**：按项目配置不同测试套件，每个套件有独立的测试类型语义
- **Postman/Newman**：API 测试直接称为 "collection run"，不套用 "e2e" 标签

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 术语混乱持续恶化，阻碍后续重构 | Rejected: 问题会随 Forge 支持的 surface 类型增多而加剧 |
| 统一改名 "高级测试" | ARCHITECTURE.md 现有术语 | 改动小 | 仍然是一个笼统概念，不解决类型错配问题 | Rejected: 只是换了一个模糊标签 |
| 引入标准测试分层（unit/integration/e2e） | 行业标准 | 概念通用 | 不适合 Forge 的场景——CLI 测试不是传统意义上的 integration test，API 测试也不是传统意义上的 contract test | Rejected: 行业分层模型与 Forge 的 surface 模型不对齐 |
| **Surface → Test Type 映射** | Forge 自身实践 | 精确、与已有分化一致、可扩展 | 需要更新多个文件和概念 | **Selected: 最小惊讶原则——名称匹配实际行为** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Forge 已有 surface 分化的基础设施：
- `gen-test-scripts/types/` 下 5 种生成策略
- `run-tests/rules/surfaces/` 下 5 种编排规则
- `init-justfile/rules/surfaces/` 下 5 种 justfile 模板
- `gen-journeys/rules/surface-*.md` 下 5 种 surface 规则

本提案将这些已有分化在**概念层**和**命名层**统一表达。

### Resource & Timeline

纯文档 + 命名变更，不涉及核心逻辑重构。预计工作量：
- 概念参考文档：1 个 doc 任务
- 术语更新（skill 文件、文档）：若干 doc 任务
- 命名变更（task type、justfile recipe）：若干 coding 任务

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "所有 Forge 生成的测试都是 e2e 测试" | Assumption Flip：如果 CLI 测试是 e2e，那它端到端测试了什么？答案是一个子进程调用，这不是 "端到端" | Overturned: 只有 Web surface 的测试真正端到端 |
| "用户不需要知道测试类型差异" | Stress Test：当 CLI 项目的测试报告说 "e2e 覆盖率 100%" 时，用户会以为所有功能都端到端验证了，但实际上只验证了子进程层面 | Confirmed: 误导性命名影响用户判断 |
| "统一叫 e2e 可以简化概念" | Occam's Razor：简化了命名但增加了认知负担——用户需要自行区分 "这个 e2e 测试实际做了什么" | Refined: 统一名称 ≠ 简化概念，精确命名才是真正的简化 |

## Scope

### In Scope

- 定义 Surface → Test Type 的映射模型和语义
- 梳理当前所有使用 "e2e" 术语的文件和位置
- 编写测试类型概念参考文档
- **更新 guide.md（Terminology 部分），补充 Surface Type → Test Type 的简要说明**，使所有 agent 在任务执行时能正确使用测试类型术语
- 更新 ARCHITECTURE.md 中的测试相关章节
- 更新 skill SKILL.md 文件中的测试类型术语
- 更新 justfile recipe 命名（从 `test`/`test-e2e` 到 surface-specific 名称）
- 更新 task type 命名（携带 surface 信息）
- 更新 business rules 文档中的测试相关术语
- 更新 gen-test-scripts 输出的测试代码中的注释/标签

### Out of Scope

- 测试管线流程的重构（gen-journeys → gen-contracts → gen-test-scripts → run-tests 的流程不变）
- 质量门逻辑的改动（两层门结构不变，只是第二层的命名更精确）
- eval 管线的改动（eval-journey、eval-contract 的评分维度不变）
- 测试目录结构的重新组织（保持 `tests/<journey>/` 或按 surface-key 分目录可作为后续优化）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 术语变更导致已有文档/教程失效 | M | M | 在变更点添加术语映射表（旧术语 → 新术语），方便迁移 |
| justfile recipe 重命名影响已有项目的 CI 流程 | M | H | 提供向后兼容的 alias recipe（旧名 → 新名），设置过渡期 |
| 新概念增加用户学习成本 | L | L | 概念文档以一页纸为限，映射表一目了然 |

## Success Criteria

- [ ] 概念参考文档完成，包含 5 种 surface 的测试类型定义、语义、验证维度和执行模型
- [ ] guide.md Terminology 部分包含 Surface → Test Type 映射的简要说明，agent 可据此正确使用测试类型术语
- [ ] Forge 代码库中不再有将所有生成测试统称为 "e2e" 的地方（搜索 "e2e" 只出现在 Web surface 的上下文中）
- [ ] 所有涉及测试类型的 skill 文件使用 surface-specific 测试类型名称
- [ ] 概念文档被至少 3 个现有 skill 的 rules 文件引用

## Next Steps

- 如变更范围 > 15 个 coding 任务，转入 full pipeline（/write-prd → /tech-design → /breakdown-tasks）
- 如变更范围可控，使用 /quick-tasks 直接生成任务
