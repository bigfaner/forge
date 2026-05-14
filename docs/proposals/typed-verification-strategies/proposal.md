---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: Interface-Type-Specific Verification Strategies

## Problem

Forge 的测试生成 pipeline 不区分 interface 类型的验证策略。gen-test-cases 将所有 UI 测试用例归为同一类型，生成相同的验证条件；gen-test-scripts 使用同一套模板逻辑生成测试代码。但不同 interface 类型（TUI、web-ui、mobile-ui、API、CLI）的"正确性"定义完全不同：

- TUI 的正确性是**视觉渲染正确性**（对齐、溢出、宽度），不是 DOM 状态
- web-ui 的正确性是**DOM 交互正确性**（元素可见、可点击、状态变化）
- mobile-ui 的正确性是**设备适配正确性**（触摸、方向、平台差异）
- API 的正确性是**契约正确性**（请求/响应符合 spec、错误路径完整）
- CLI 的正确性是**输出正确性**（退出码、输出格式、参数组合）

### Evidence

**TUI lesson**（docs/lessons/lesson-tui-visual-verify.md）：deep-drill-analytics feature 的 11 个 bug 全部通过了"编译 + 测试"verify gate，因为测试只检查了逻辑正确性（函数返回值），没有检查视觉渲染正确性（行数、宽度、对齐）。根因是 gen-test-cases 不为 TUI 类型生成 golden file 测试用例和维度检查用例。

**API/CLI 盲区**：当前 profile 声明了 `api` 和 `cli` capability，但 gen-test-cases 只生成功能性测试用例，不生成契约测试、边界值测试、参数组合测试等集成层面的测试。

### Urgency

每次 TUI feature 都需要手动在 task verify criteria 中追加 golden file 和维度检查条件，依赖人工记忆。lesson 已经证明这条路不可靠——漏掉的 bug 比发现的多。

## Proposed Solution

### 核心机制

1. **Profile 策略文件**：每个 profile 新增 `verification-strategies.md`，定义其各 capability 的验证策略（包含验证维度、边界场景、测试数据要求、测试级别标记）
2. **类型化用例生成**：gen-test-cases 读取 profile 策略，按 interface 类型生成不同的验证条件，并自动标记测试级别（e2e / integration）
3. **级别化脚本生成**：gen-test-scripts 根据测试级别和 interface 类型选择不同的代码生成策略

### 类型验证策略矩阵

| Interface Type | Test Level | 核心验证维度 | 边界场景 |
|---------------|------------|-------------|---------|
| **TUI** | e2e | Golden file 对比、维度检查（行数=高度/宽度<=终端宽度）、ANSI 色码一致性 | CJK 字符、长路径(>50)、多位数字(>9)、空字段、窄终端(80x24)、宽终端(140x40) |
| **web-ui** | e2e | DOM 交互、视觉回归（截图）、响应式布局、可访问性 | 空状态、加载态、错误态、边界数据量、不同视口尺寸 |
| **mobile-ui** | e2e | 设备渲染、触摸交互、屏幕方向、平台差异 | 横竖屏切换、低网速、推送通知中断、小屏设备 |
| **API** | integration | 契约验证（请求/响应对 spec）、错误路径、边界值、真实依赖 | 无效输入、认证失败、超限、并发、空响应、大数据量 |
| **CLI** | integration | 输出 golden file、退出码、参数组合、管道兼容 | 无效参数、--help、参数互斥、空输入、管道+重定向 |

### 测试级别定义

| Level | 含义 | 触发条件 | 验证方式 |
|-------|------|---------|---------|
| **e2e** | 端到端，验证用户可见的完整行为 | interface 类型为 visual/interactive（TUI、web-ui、mobile-ui） | 渲染输出、DOM 状态、截图、设备行为 |
| **integration** | 集成测试，验证组件间协作行为 | interface 类型为 non-visual（API、CLI） | HTTP 契约、退出码、输出格式、参数解析 |

### 创新点

**策略与 profile 解耦但就近定义**：策略文件在 profile 目录内（每个 profile 独立定义），但 gen-test-cases 通过统一的 capability key 查询策略，不需要硬编码类型映射。新增 profile 只需加策略文件，不需要改 skill 逻辑。

## Requirements Analysis

### Key Scenarios

1. **TUI feature 测试生成**：gen-test-cases 检测到 `tui` capability → 自动生成 golden file 测试用例 + 维度检查用例 + 边界场景用例，标记为 e2e 级别
2. **web-ui feature 测试生成**：gen-test-cases 检测到 `web-ui` capability → 生成 DOM 交互用例 + 视觉回归用例 + 响应式用例，标记为 e2e 级别
3. **API 集成测试生成**：gen-test-cases 检测到 `api` capability → 生成契约测试用例 + 错误路径用例 + 边界值用例，标记为 integration 级别
4. **CLI 集成测试生成**：gen-test-cases 检测到 `cli` capability → 生成输出 golden file 用例 + 退出码用例 + 参数组合用例，标记为 integration 级别
5. **Mixed profile**：go-test profile 同时有 `tui` + `api` + `cli` → 生成 3 类测试用例，TUI 标记 e2e，API/CLI 标记 integration

### Constraints & Dependencies

- 依赖现有 profile 的 capability 声明（manifest.yaml）
- 策略文件是 profile 目录的一部分，随 profile 版本更新
- 不引入新的外部依赖

## Alternatives & Industry Benchmarking

### Industry Solutions

**Testing Trophy / Testing Quadrant**：业界通用的测试分层模型（Kent C. Dodds 的 Testing Trophy、Martin Fowler 的 Testing Quadrant）将测试按粒度和目的分层。本提案将 interface 类型映射到测试级别，本质上是测试象限的 forge 化实现。

**Playwright 的 test annotation**：Playwright 支持 `test.configure({ mode: 'serial' })` 和自定义 tag（`@fast`、`@slow`），但没有自动化的类型-级别映射。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | TUI lesson 已证明不可靠 | Rejected: bug 捕获率为 0/11 |
| 共享策略文档 | 本提案早期选项 | 所有 profile 复用 | 不同 framework 的同类型测试策略不同（go-test 的 TUI vs rust-test 的 TUI） | Rejected: 粒度太粗 |
| 硬编码在 skill 中 | 传统做法 | 简单直接 | 新增 profile 需要改 skill 代码 | Rejected: 耦合 |
| **Profile 策略文件** | 本提案 | 按框架定制、新增 profile 零改动 skill | 6 个 profile 各写一份 | **Selected: 灵活性最佳** |

## Feasibility Assessment

### Technical Feasibility

所有改动都在 skill 文件层面（SKILL.md、模板文件、profile 目录内的策略文件），不涉及编译代码。profile 已有 capability 声明，策略文件只是补充验证维度信息。

### Resource & Timeline

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | 6 个 profile 各写 verification-strategies.md | 2h |
| 2 | test-cases.md 模板更新（Level 字段 + 类型 section） | 1h |
| 3 | gen-test-cases SKILL.md 增强策略读取和类型化生成 | 2h |
| 4 | gen-test-scripts SKILL.md 增强级别化代码生成 | 1.5h |
| 5 | eval-test-cases rubric 更新（类型化验证完整度维度） | 1h |
| 6 | 端到端验证 | 1h |
| **Total** | | **~8.5h** |

## Scope

### In Scope

- 6 个 profile 各新增 `verification-strategies.md`（go-test、web-playwright、maestro、pytest、rust-test、java-junit）
- gen-test-cases SKILL.md 增强：读取 profile 策略 → 按 interface 类型生成不同验证条件 → 自动标记 e2e/integration 级别
- gen-test-scripts SKILL.md 增强：按测试级别和 interface 类型选择代码生成策略
- test-cases.md 模板更新：新增 Level 字段、类型专属验证 section
- eval-test-cases rubric 更新：新增"类型化验证完整度"评分维度

### Out of Scope

- breakdown-tasks 改动（verify 模板注入）
- execute-task quality gate 改动
- eval-design 改动
- 新 profile 或新测试框架
- 视觉回归基础设施（截图对比服务）
- run-e2e-tests / graduate-tests 改动
- CI/CD 管道改动

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 策略文件定义的验证维度与实际 PRD 不匹配 | Medium | Medium — 生成的测试用例覆盖不全 | 策略文件列出的是最小验证维度集，agent 可根据 PRD 追加额外维度 |
| 6 个 profile 策略文件维护成本 | Low | Low — 策略相对稳定 | 策略变化频率低（随 profile 大版本更新） |
| gen-test-cases 读取策略后 token 开销增加 | Medium | Low — 策略文件精简 | 策略文件控制在 200 行以内，只读当前 profile 的策略 |
| TUI golden file 测试在 CI 环境中的 terminal 模拟差异 | Medium | High — CI 中 golden test 不稳定 | 策略文件中明确要求 golden test 使用固定 terminal 尺寸（如 80x24），不依赖环境 |

## Success Criteria

- [ ] 6 个 profile 各有 verification-strategies.md，定义了各 capability 的验证维度和边界场景
- [ ] gen-test-cases 生成的 test-cases.md 包含 Level 字段（e2e/integration），且 interface 类型与策略文件一致
- [ ] TUI 类型的 test-cases 自动包含 golden file + 维度检查 + 边界场景用例
- [ ] API 类型的 test-cases 自动包含契约测试 + 错误路径 + 边界值用例
- [ ] gen-test-scripts 按 Level 字段选择不同的代码生成策略
- [ ] eval-test-cases rubric 包含"类型化验证完整度"维度

## Next Steps

- Proceed to `/write-prd` to formalize requirements
