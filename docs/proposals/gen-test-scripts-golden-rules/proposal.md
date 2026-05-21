---
created: 2026-05-21
author: faner
status: Draft
---

# Proposal: gen-test-scripts 内建类型黄金守则

## Problem

gen-test-scripts 完全依赖 Convention 文件提供框架特定规则，自身零 TUI/CLI/UI/Mobile/API 知识。当 Convention 缺失时，测试代码生成的质量完全失控——生成的代码既不遵循 binary 隔离原则，也不考虑会话复用、选择器优先级等跨框架通用最佳实践。

### Evidence

- gen-test-scripts 的 SKILL.md（326 行）是纯流程编排，没有任何界面类型特定的生成规则
- Domain Reconnaissance 表格仅提及侦察目标（"TUI components: Model fields, View output patterns"），但不定义如何利用这些信息生成代码
- 任务模板如 `gen-test-scripts-go-cli.md` 仅 13 行，全靠 skill prompt 兜底，但 skill prompt 里也没有 CLI 知识
- gen-test-cases 已通过 `types/` 目录实现类型分发（cli.md、tui.md、ui.md、mobile.md、api.md），gen-test-scripts 缺少对等机制

### Urgency

forge-architecture-simplification 包含测试生成任务，这些任务即将执行。如果 gen-test-scripts 没有内建黄金守则，生成的测试代码将反复违反基本的测试原则（binary 隔离、会话复用等），导致质量门禁失败和返工。

## Proposed Solution

在 gen-test-scripts 中新增 `types/` 目录，为 CLI、TUI、UI、Mobile、API 五种界面类型各创建一个黄金守则文件。这些守则定义**原则性约束**（WHAT），Convention 文件提供**框架实现细节**（HOW），两者互补不冲突。

### Innovation Highlights

- **原则/实现分层**：types/ 定义跨框架通用原则（"CLI 测试必须使用编译产物"），Convention 定义框架特定实现（"Go 用 exec.Command，JS 用 child_process"）。这种分层让守则在任何框架下都有效，同时不限制框架选择自由度
- **与 gen-test-cases 对称**：gen-test-cases 的 types/ 管"测试用例描述规则"，gen-test-scripts 的 types/ 管"测试代码生成原则"，形成完整覆盖

## Requirements Analysis

### Key Scenarios

- **Convention 存在**：types/ 黄金守则 + Convention 实现细节共同指导生成
- **Convention 缺失**：types/ 黄金守则作为唯一规则源，LLM defaults 仅用于无原则覆盖的细节
- **Convention 与 types/ 冲突**：Convention 可补充但不能覆盖 types/ 的原则性约束

### Non-Functional Requirements

- 每个 type 文件保持框架无关，不包含特定语言的代码模板
- type 文件内容必须足够具体，能直接指导代码生成（不能只是泛泛的原则陈述）

### Constraints & Dependencies

- types/ 文件随 gen-test-scripts skill 一起分发（遵循 forge-distribution.md 的分发模型）
- 不引入跨 skill 依赖（types/ 是 gen-test-scripts 内部文件）

## Alternatives & Industry Benchmarking

### Industry Solutions

业界测试框架（Playwright、Cypress、Detox）通常在框架文档中内建最佳实践指南，而非完全依赖用户配置。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | Convention 缺失时质量失控 | Rejected: 已有证据表明问题存在 |
| 跨 skill 共享规则目录 | 自研 | DRY | 违背 forge 分发模型（skill 间不直接引用文件） | Rejected: 分发约束 |
| 仅强化任务模板 | 自研 | 改动小 | 任务模板只能指引参考方向，不能定义规则本身 | Rejected: 治标不治本 |
| **types/ 黄金守则** | 对标 gen-test-cases | 架构一致、自包含、框架无关 | 5 个新文件 + SKILL.md 修改 | **Selected: 架构一致性最优** |

## Feasibility Assessment

### Technical Feasibility

gen-test-cases 已验证 types/ 模式可行。gen-test-scripts 只需复制这一模式，调整内容为代码生成导向。

### Resource & Timeline

5 个 type 文件（每个约 80-120 行）+ SKILL.md 修改（约 20 行），预计 1 个 coding task 可完成。

### Dependency Readiness

无外部依赖。仅需 gen-test-scripts skill 目录。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "Convention 文件总是存在" | Stress Test | Overturned: 冷启动项目、新用户、test-guide 未运行时均无 Convention，skill 不应假设 Convention 存在 |
| "LLM defaults 足以生成好的测试代码" | Assumption Flip | Overturned: LLM 不知道 binary 隔离、会话复用等原则，会生成 `go run` 式测试、每个 case 都重新登录等反模式 |
| "gen-test-cases 的 types/ 足以覆盖 gen-test-scripts" | XY Detection | Overturned: gen-test-cases 的 types/ 管测试用例描述格式，gen-test-scripts 需要代码生成规则，是不同维度 |

## Scope

### In Scope

- `plugins/forge/skills/gen-test-scripts/types/cli.md` — CLI 测试代码生成黄金守则
- `plugins/forge/skills/gen-test-scripts/types/tui.md` — TUI 测试代码生成黄金守则
- `plugins/forge/skills/gen-test-scripts/types/ui.md` — Web UI 测试代码生成黄金守则
- `plugins/forge/skills/gen-test-scripts/types/mobile.md` — Mobile 测试代码生成黄金守则
- `plugins/forge/skills/gen-test-scripts/types/api.md` — API 测试代码生成黄金守则
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — 新增类型规则加载逻辑

### Out of Scope

- 修改 gen-test-cases 的 types/（已有，定位不同）
- 修改 gen-journeys、gen-contracts（上游 skill，不属于生成层）
- 修改 run-e2e-tests（执行层，不涉及代码生成）
- 具体框架代码模板（仍由 Convention 提供）
- 任务模板强化（可作为 follow-up）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| types/ 规则与 Convention 内容重叠导致执行困惑 | M | M | 明确分层：types/ = 原则，Convention = 实现，SKILL.md 说明优先级 |
| type 文件过于抽象，无法有效指导代码生成 | L | H | 每个 type 文件包含具体的 do/don't 示例模式 |
| 新增 type 文件增加 skill 维护成本 | L | L | type 文件是稳定的黄金守则，不随框架变化 |

## Success Criteria

- [ ] gen-test-scripts 对 CLI/TUI/UI/Mobile/API 五种类型始终遵循对应的黄金守则
- [ ] Convention 只补充框架实现细节，不覆盖 types/ 定义的原则性约束
- [ ] 每个 type 文件覆盖断言原则 + 交互原则 + 策略原则三个维度
- [ ] 所有 type 文件保持框架无关，不包含特定语言或测试框架的代码
- [ ] SKILL.md 明确 types/ 与 Convention 的原则/实现分层关系

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
