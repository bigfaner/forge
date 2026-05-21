---
created: 2026-05-21
author: faner
status: Approved
---

# Proposal: gen-test-scripts 内建类型黄金守则

## Problem

gen-test-scripts 的 `types/` 目录中已有 5 个类型文件（cli/tui/ui/mobile/api），但这些文件存在三个层面的问题，导致它们无法有效发挥黄金守则的作用：

1. **未接入流程**：SKILL.md 没有任何加载 types/ 的逻辑，5 个文件形同虚设
2. **原则纯度不足**：文件中混合了黄金守则（WHAT）和框架实现细节（HOW），定位模糊。最严重的是 mobile.md 几乎完全绑定 Maestro YAML 语法
3. **守则覆盖缺失**：跨类型通用原则（超时保护、确定性、隔离）零覆盖；各类型缺少关键守则（CLI 缺超时、TUI 缺终端尺寸、UI 缺会话复用、Mobile 缺状态重置）

### Evidence

- SKILL.md（326 行）无一处引用 `types/` 目录
- 所有 5 个 type 文件的 Output 路径写 `tests/e2e/features/<feature>/`，SKILL.md Step 3 写 `tests/<journey>/`——路径不一致
- 所有 5 个 type 文件引用 "SKILL.md Step 1.5"，但 SKILL.md 只有 Step 1.1-1.4——步骤编号不存在
- mobile.md 的 Generation Patterns 章节直接输出 Maestro YAML 语法，违反框架无关原则
- cli.md 侦察策略包含 Go/Node.js/Python 特定 grep 命令，超出原则层范畴
- 0 个文件提到超时保护、确定性原则、或跨类型共享反模式

### Urgency

forge-architecture-simplification 包含测试生成任务，这些任务即将执行。如果 types/ 文件不修正好并接入 SKILL.md，生成的测试代码将缺少超时保护（CI 悬挂）、Binary 隔离（测试不稳定）、会话复用（性能浪费）等基本保障。

## Proposed Solution

分三步修正：

1. **新增 `types/_shared.md`**：定义跨类型通用黄金守则（隔离、确定性、超时保护、幂等性、资源清理），5 个类型文件引用而非重复。三层模型：`_shared.md`（抽象原则）→ 类型文件 Golden Rules（类型特定约束）→ Convention（框架实现）
2. **修正 5 个类型文件**：每个文件分为 `## Golden Rules`（框架无关原则）和 `## Reconnaissance Hints`（发现辅助，不指导生成）两个清晰区域；补充专家评估指出的缺失守则；修正 Output 路径和步骤编号引用
3. **修改 SKILL.md**：新增类型规则加载步骤，明确 types/（原则）与 Convention（实现）的分层关系

### Innovation Highlights

- **原则/实现分层**：types/ 定义跨框架通用原则（"CLI 测试必须使用编译产物"），Convention 定义框架特定实现（"Go 用 exec.Command"）。types/ 的 Golden Rules 章节是声明式约束，Reconnaissance Hints 是辅助发现手段，不直接指导代码生成
- **跨类型共享层**：通过 `_shared.md` 统一隔离、确定性、超时等跨类型原则，避免 5 个文件重复定义且不一致
- **与 gen-test-cases 对称**：gen-test-cases 的 types/ 管"测试用例描述规则"，gen-test-scripts 的 types/ 管"测试代码生成原则"

## Requirements Analysis

### Key Scenarios

- **Convention 存在**：types/ 黄金守则（原则）+ Convention（框架实现）共同指导生成
- **Convention 缺失**：types/ 黄金守则作为原则源，LLM defaults 填充框架实现细节
- **Convention 与 types/ 重叠**：types/ 的 Golden Rules 不可覆盖；Reconnaissance Hints 可被 Convention 的框架声明覆盖
- **混合类型 Journey**：一个 Journey 的不同 Step 可能涉及不同接口类型。加载策略：读取 Contract 后提取所有涉及的接口类型，加载 `_shared.md` + 所有匹配类型的 type 文件。按 Step 各自类型应用对应 type 规则。如涉及超过 3 个类型，警告 token 预算风险

### Non-Functional Requirements

- 每个 type 文件的 Golden Rules 章节保持框架无关
- Reconnaissance Hints 标注为辅助发现手段，不指导代码生成。Hints 中发现的信息应转化为 Fact Table 值，而非直接进入生成代码
- 按需加载：仅加载与检测到的接口类型匹配的 type 文件

### Constraints & Dependencies

- types/ 文件随 gen-test-scripts skill 一起分发（遵循 forge-distribution.md）
- 不引入跨 skill 依赖

## Alternatives & Industry Benchmarking

### Industry Solutions

业界测试框架（Playwright、Cypress、Detox）通常在框架文档中内建最佳实践指南，而非完全依赖用户配置。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 现有文件未接入、质量参差 | Rejected: 问题已确认 |
| 跨 skill 共享规则目录 | 自研 | DRY | 违背分发模型 | Rejected: 分发约束 |
| 仅修复 SKILL.md 加载逻辑 | 自研 | 改动小 | 现有文件质量问题未解决 | Rejected: 治标不治本 |
| **修正 types/ + 接入 SKILL.md** | 对标 gen-test-cases | 解决根因 | 6 个文件修改 | **Selected: 全面修正** |

## Feasibility Assessment

### Technical Feasibility

gen-test-cases 已验证 types/ 模式可行。现有 types/ 文件已有良好基础（侦察策略、Fact Table 门禁、反模式防护），主要需要结构调整和守则补充。

### Resource & Timeline

- 1 个 doc task：`_shared.md` 新建（隔离、确定性、超时保护、幂等性、资源清理五大原则）
- 1 个 coding task：cli.md + tui.md + api.md 结构修正 + 守则补充 + 路径/引用修复
- 1 个 coding task：ui.md 修正 + mobile.md 重写（框架无关化）
- 1 个 coding task：SKILL.md 加载逻辑 + 优先级声明（类型规则加载在 Step 2 读取 Contract 之后、Step 3 生成代码之前）
- 总计 4 个 task

### Dependency Readiness

无外部依赖。所有修改限定在 gen-test-scripts skill 目录内。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "types/ 文件已存在，只需接入 SKILL.md" | 5 Whys | Overturned: 文件存在但质量参差——mobile.md 绑定 Maestro、Output 路径不一致、步骤引用不存在、缺少超时/确定性/隔离通用原则 |
| "types/ 的侦察策略属于原则层" | Assumption Flip | Refined: 侦察策略是辅助发现手段（HOW），不是黄金守则（WHAT）。应明确标注为 Hints 而非 Rules |
| "每个类型文件独立定义反模式即可" | Stress Test | Overturned: 5 个文件的反模式有大量重叠（Sleep、硬编码、空断言），提取到 `_shared.md` 更一致且易维护 |

## Scope

### In Scope

**新建**：
- `plugins/forge/skills/gen-test-scripts/types/_shared.md` — 跨类型通用黄金守则

**修正（结构 + 内容）**：
- `plugins/forge/skills/gen-test-scripts/types/cli.md` — 补充超时保护、Binary 编译隔离、环境变量密封性；重构为 Golden Rules + Reconnaissance Hints 双区域
- `plugins/forge/skills/gen-test-scripts/types/tui.md` — 补充终端尺寸契约、ANSI 序列处理、稳定态检测；同上重构
- `plugins/forge/skills/gen-test-scripts/types/ui.md` — 补充会话复用、网络拦截策略、视口管理；同上重构
- `plugins/forge/skills/gen-test-scripts/types/mobile.md` — 重写为框架无关原则层，补充 App 状态重置、权限处理、Deep Link 模式；将 Maestro 特定语法移至 Reconnaissance Hints
- `plugins/forge/skills/gen-test-scripts/types/api.md` — 补充幂等性验证、请求超时、Content-Type 验证；同上重构

**集成**：
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — 新增类型规则加载步骤 + types/ vs Convention 优先级声明

**修正（全局）**：
- 所有 5 个 type 文件的 Output 路径统一为 `tests/<journey>/`
- 所有 5 个 type 文件的步骤编号引用修正

### Out of Scope

- 修改 gen-test-cases 的 types/（定位不同，管测试用例描述）
- 修改 gen-journeys、gen-contracts（上游 skill）
- 修改 run-e2e-tests（执行层）
- 具体框架代码模板（仍由 Convention 提供）
- 任务模板强化（可作为 follow-up）
- templates/node_modules 清理（独立问题，影响 token 预算但非本提案范围）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| types/ 重构后内容过多，影响 token 预算 | M | M | 按需加载：仅加载匹配类型的 type 文件 + `_shared.md`，不全量加载 |
| Reconnaissance Hints 中的语言特定 grep 仍可能误导 LLM 生成特定框架代码 | M | M | 在 Hints 区域加入注释标注"discovery hints, not generation instructions"；Golden Rules 区域完全框架无关 |
| mobile.md 重写为框架无关后实用性下降 | L | H | 保留 Maestro 作为最常见框架的参考示例，但标注为"示例"而非"规则" |
| `_shared.md` 与 5 个类型文件的守则重复 | L | L | `_shared.md` 只定义共享原则，类型文件只定义类型特有规则，引用而非重复 |

## Success Criteria

- [ ] SKILL.md 包含 types/ 加载逻辑，按需加载匹配类型的 type 文件
- [ ] SKILL.md 包含 `<HARD-RULE>` 声明 types/（原则）与 Convention（实现）的分层关系
- [ ] `_shared.md` 定义跨类型通用原则：隔离、确定性、超时保护、幂等性、资源清理
- [ ] 每个 type 文件分为 `## Golden Rules` 和 `## Reconnaissance Hints` 两个清晰区域
- [ ] Golden Rules 区域完全框架无关，不含语言特定代码或命令
- [ ] 所有 Output 路径统一为 `tests/<journey>/`
- [ ] 所有步骤编号引用与 SKILL.md 一致
- [ ] CLI 补充：超时保护（区分进程级超时和测试函数级超时）、Binary 编译隔离、环境变量密封性
- [ ] TUI 补充：终端尺寸契约、ANSI 序列处理、稳定态检测
- [ ] UI 补充：会话复用、网络拦截策略、视口管理
- [ ] Mobile 补充：App 状态重置、权限处理、Deep Link 模式
- [ ] API 补充：幂等性验证、请求超时、Content-Type 验证

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
