---
created: 2026-05-20
author: "faner + Claude"
status: Draft
---

# Proposal: task-executor 与 prompt 模板契合度修复

## Problem

task-executor agent 的 Execution Protocol（11 步）与 prompt 模板生成的执行策略之间存在 5 处语义断裂，导致 agent 可能执行与 skill 语义矛盾的操作（如 blocked 任务仍 commit）或忽略已编码的配置规则（如 scope resolution）。

### Evidence

通过逐一阅读 task-executor.md、全部 22 个 prompt 模板、submit-task SKILL.md、git-commit 命令和 hook-injected guide，发现以下具体矛盾：

| # | 矛盾 | 位置 | 严重度 |
|---|------|------|--------|
| 1 | submit-task 可 auto-downgrade 为 blocked，但 agent 不会在 step 8→9 之间复查 | task-executor.md step 8→9 | P0 |
| 5 | guide 说 fmt → WARNING (non-blocking)，模板说 fmt → Stop | guide vs 8 个 coding/validation/gate 模板 | P1 |
| 6 | guide 定义了 scope resolution（检查 project-type），但模板直接盲传 {{SCOPE}} | guide vs prompt.go + 模板 | P1 |
| 7 | task-executor.md 注释称 submit-task "via just test" 做指标收集，实际不是 | task-executor.md step 8 注释 | P1 |
| 9 | blocked 路径无替代 DONE 输出格式 | task-executor.md step 10 | P3 |

### Urgency

P0 问题意味着：当 `forge task submit` 发现 testsFailed > 0 并 auto-downgrade 为 blocked 时，agent 仍会执行 git-commit。这违反 submit-task SKILL.md 的 "do NOT commit" 约束，在 CI 中产生应该被阻断的 commit。

## Proposed Solution

逐一修复 5 处矛盾，使 task-executor agent 的 Execution Protocol 与 prompt 模板、guide、submit-task skill 语义完全对齐。

### Innovation Highlights

无创新，纯一致性修复。问题根因是 typed-task-dispatch 和 scope-resolution 两个 feature 独立演进时，agent 协议未同步更新。

## Requirements Analysis

### Key Scenarios

- **场景 1**（P0）：coding task 的 targeted tests 通过，但 `forge task submit` 的全量测试发现 testsFailed > 0 → auto-downgrade → 不应 commit → 但当前会 commit
- **场景 2**（P1）：`just fmt` 在 toolchain 异常时失败 → agent Stop → 任务永久卡在 in_progress → 实际是 toolchain 问题而非代码问题
- **场景 3**（P1）：backend 项目的 task scope = "frontend" → 模板生成 `just compile frontend` → 应该 fallback 到 `just compile`
- **场景 4**（P1）：agent 读到 submit-task "via just test" 的注释 → 假设不需要自己跑 targeted tests → 但实际需要
- **场景 5**（P3）：blocked 任务输出 `DONE: ... | ✅ | <commit-hash> | ...` 但无 commit-hash → run-tasks dispatcher 解析出错

### Non-Functional Requirements

- 所有改动限于 markdown 文件和 prompt.go，不涉及 forge-cli 核心逻辑
- 不改变现有 task 类型定义或 index schema
- 向后兼容：不影响已生成的 task 文件

### Constraints & Dependencies

- `forge-cli/pkg/forgeconfig` 已支持 `forge config get project-type`
- prompt.go 的 `renderTemplate` 已有 scope 空值处理，需扩展为完整 resolution
- 改动需通过 `forge-cli` 的现有测试套件

## Alternatives & Industry Benchmarking

### Industry Solutions

这是内部 agent 协议一致性问题，无行业标准可参照。类似问题在 multi-prompt orchestration 系统中通过 protocol state machine 或 message contract 解决。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | P0 缺陷持续存在 | Rejected: blocked task 违规 commit |
| 全面重写 agent 协议 | — | 彻底 | 改动面大，引入新风险 | Rejected: 过度工程 |
| **最小定点修复** | 本 proposal | 改动量小（~20 行），风险可控 | 不解决结构性问题 | **Selected: 最小改动修复 5 处具体矛盾** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有矛盾点都有明确的文件和行号，改动是增量和确定性的。

### Resource & Timeline

5 处修改，预计 1-2 个 coding task 即可完成。

### Dependency Readiness

无外部依赖。forgeconfig 和 prompt.go 的 API 已经就绪。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| task-executor 的 step 7 检查足够保证 blocked 安全 | Stress Test | Overturned: submit-task 的 auto-downgrade 发生在 step 7 之后 |
| guide 和 template 的 fmt 行为一致 | 直接对比 | Overturned: guide 说 WARNING，template 说 Stop |
| scope resolution 已在自动化路径中实现 | Occam's Razor | Overturned: 仅在 guide 自然语言中描述，未编码 |
| submit-task 注释准确描述了其行为 | 直接验证 | Overturned: submit-task 不运行 just test |

## Scope

### In Scope

- task-executor.md: 加 step 8.5（blocked 复查）+ 修正 submit-task 注释 + 补 blocked 输出格式
- prompt/data/*.md: 8 个模板的 fmt 行为从 Stop 改为 WARNING
- prompt.go: 扩展 scope resolution 逻辑（检查 project-type）

### Out of Scope

- template 之间 convention 加载覆盖不一致的问题（P3，与 agent-prompt 契合无关）
- run-tasks.md 和 execute-task.md 的 dispatcher 层面改动
- guide 本身的修改

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| prompt.go scope resolution 引入新 bug | L | M | 已有 scope_resolution_cli_test.go 覆盖 |
| fmt WARNING 导致格式问题未被发现 | L | L | All-Completed Hook 的 quality gate 仍会运行 fmt |
| blocked DONE 格式改变影响 run-tasks 解析 | L | M | run-tasks 不解析 DONE 输出（仅检查 forge task status） |

## Success Criteria

- [ ] submit-task auto-downgrade 为 blocked 时，agent 不执行 git-commit（可构造 testsFailed > 0 场景验证）
- [ ] `just fmt` 失败时 agent 输出 warning 并继续，不 Stop
- [ ] backend 项目的 task scope="frontend" 时，模板生成 `just compile`（无 scope 参数）
- [ ] task-executor.md 中不再包含 "via just test" 的错误注释
- [ ] blocked 任务输出 `DONE: ... | ❌ | blocked | ...` 格式

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
