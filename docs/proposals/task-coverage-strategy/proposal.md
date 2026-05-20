---
created: 2026-05-20
author: "fanhuifeng"
status: Draft
---

# Proposal: Task Coverage Strategy

## Problem

Task executor agents 循环于测试修复和级联编译错误中无法收敛，尤其在 bug fix 和深度重构场景下，反复尝试导致 token 浪费和超时。

### Evidence

- `coding.fix` 任务：agent 写了测试但反复修改无法通过，缺少明确的"够好了"的退出信号
- `coding.refactor` 任务：修改一个接口引发级联编译错误，agent 逐个修复但越修越多，缺少增量编译策略
- `coding.cleanup` 任务：删除死代码后仍然被要求写新测试达到 80% 覆盖率，纯属浪费

### Urgency

当前覆盖率目标（80%）仅作为 CLAUDE.md 中的指导性文字存在，对 agent 行为无实际约束。agent 在不确定何时停止时倾向于"再多试一次"，直接导致 token 浪费和执行超时。每次失败都会消耗大量时间和成本。

## Proposed Solution

按任务类型分三档覆盖率策略，通过配置驱动 agent 行为，给 agent 明确的"达标即停"信号：

| 档位 | 任务类型 | 策略 |
|------|---------|------|
| 高覆盖 | `coding.feature` | 80% 覆盖率目标 |
| 中覆盖 | `coding.enhancement`, `coding.fix` | 60% 覆盖率目标 |
| 保持 | `coding.refactor`, `coding.cleanup`, `coding.clean` | 不新增测试，保持现有覆盖率，下降不超过 2% |

配置存放于 `.forge/config.yaml`，支持全局默认 + 任务 frontmatter 覆盖。`coding.refactor` prompt 模板增加增量编译检查策略（每改一个文件后立即 compile）。

### Innovation Highlights

业界标准做法是 CI 中强制覆盖率门（如 JaCoCo、Istanbul 的 coverage threshold），但这对 AI agent 执行场景不够：agent 需要在执行过程中就知道目标，而不是事后发现不达标。

本方案的创新在于：
- **配置驱动 prompt 注入**：覆盖率目标不是 CI 门，而是直接写入 agent 的执行指令中，影响 agent 的行为策略
- **maintain 策略**：对重构/清理类任务不设数字目标，而是要求"不下降"，避免 agent 为了数字而写无意义测试
- **增量编译策略**：在 refactor prompt 中注入"改一个文件即编译"的策略，从根源解决级联错误问题

## Requirements Analysis

### Key Scenarios

- Agent 执行 `coding.feature` 任务，写完实现后写测试，覆盖率达到 80% 即停止补充测试
- Agent 执行 `coding.fix` 任务，写针对性修复测试，覆盖率达到 60% 即停止
- Agent 执行 `coding.refactor` 任务，每改一个文件立即 compile，测试只需通过且覆盖率不下降超 2%
- 用户在任务 frontmatter 中设置 `coverage: 90`，覆盖全局默认的 80%
- 用户在 `.forge/config.yaml` 中自定义各类型默认值

### Non-Functional Requirements

- 配置缺失时平滑降级到合理的内置默认值，不报错不阻塞
- 覆盖率目标传递到 prompt 不增加显著延迟
- 任务 frontmatter 覆盖优先级高于全局配置

### Constraints & Dependencies

- 依赖现有的 Go CLI prompt 合成机制（`forge-cli/pkg/prompt/`）
- 覆盖率值仅作为 agent 指令，不在 submit 时强制校验（保持现有灵活性）
- 配置格式需兼容现有 `.forge/config.yaml` schema

## Alternatives & Industry Benchmarking

### Industry Solutions

- **CI 覆盖率门**（JaCoCo、Istanbul、go test -cover）：事后强制，适合 CI 流水线但不适合 agent 执行时引导
- **Aider/Testpilot 等 AI 编码工具**：通常不区分任务类型，统一测试策略

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | agent 继续陷入循环，每次浪费 token 和时间 | Rejected: 成本持续累积 |
| CI 强制覆盖率门 | JaCoCo/Istanbul | 业界标准，可靠 | 事后检查，不改变 agent 行为，治标不治本 | Rejected: 无法预防死循环 |
| 统一降低覆盖率目标 | — | 简单 | 对 feature 任务降低了质量标准 | Rejected: 一刀切不可取 |
| **按类型分档 + 配置驱动** | 原创 | 精准匹配任务复杂度，配置灵活，从根源引导行为 | 需要修改 config schema、prompt 模板和 CLI 代码 | **Selected: 最精准且可扩展** |

## Feasibility Assessment

### Technical Feasibility

完全可行。需要修改的组件：
1. `.forge/config.yaml` schema — 新增 `coverage` 段
2. Go CLI config 解析 — 读取覆盖率配置
3. Prompt 合成 — 将覆盖率目标注入 `{{COVERAGE_TARGET}}` 占位符
4. Prompt 模板 — 各 coding 类型模板使用覆盖率指令
5. Task frontmatter 解析 — 读取可选 `coverage` 字段

### Resource & Timeline

规模小，预计 3-5 个任务即可完成。

### Dependency Readiness

所有依赖（CLI config 系统、prompt 合成、task frontmatter 解析）均已就绪。

## Scope

### In Scope

- `.forge/config.yaml` 新增 `coverage` 配置段（按任务类型设置策略）
- Go CLI config 解析新增 coverage 配置读取
- Task frontmatter 新增可选 `coverage` 字段（覆盖全局默认）
- Prompt 合成注入覆盖率目标到各 coding 类型模板
- `coding.refactor` prompt 增加增量编译检查策略
- 内置合理默认值，无配置时平滑降级

### Out of Scope

- Submit 时强制校验覆盖率百分比（保持柔性）
- 测试修复次数限制机制
- 质量门变更（仍为 pass/fail）
- 覆盖率趋势追踪和报告

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 覆盖率目标过低导致质量下降 | M | M | 默认值经过仔细考量；用户可通过 frontmatter 调高 |
| maintain 策略下 agent 误删测试导致覆盖率下降 | L | M | prompt 中明确"覆盖率下降超 2% 需调查" |
| 配置格式与未来需求不兼容 | L | L | 使用可扩展的 map 结构而非固定字段 |

## Success Criteria

- [ ] `coding.feature` 任务执行时 prompt 包含 80% 覆盖率目标指令
- [ ] `coding.fix` 任务执行时 prompt 包含 60% 覆盖率目标指令
- [ ] `coding.refactor` 任务执行时 prompt 包含"保持覆盖率 + 增量编译"策略
- [ ] `.forge/config.yaml` 中配置的自定义覆盖率值能正确传递到 prompt
- [ ] 任务 frontmatter 中的 `coverage` 字段能覆盖全局默认值
- [ ] 无配置时使用内置默认值，agent 行为正常

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
