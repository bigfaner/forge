---
created: 2026-05-20
author: "faner + Claude"
status: Draft
---

# Proposal: task-executor 与 prompt 模板契合度修复

## Problem

task-executor agent 的 Execution Protocol（11 步）与 prompt 模板生成的执行策略之间存在 10 处语义断裂，导致 agent 可能执行与 skill 语义矛盾的操作（如 blocked 任务仍 commit）、忽略已编码的配置规则（如 scope resolution）、或在错误处理路径上产生歧义。

### Evidence

通过逐一阅读 task-executor.md、全部 22 个 prompt 模板、submit-task SKILL.md、git-commit 命令、hook-injected guide 和合成引擎 prompt.go，发现以下具体矛盾：

| # | 矛盾 | 位置 | 严重度 |
|---|------|------|--------|
| 1 | submit-task 可 auto-downgrade 为 blocked，但 agent 不会在 step 8→9 之间复查 | task-executor.md step 8→9 | P0 |
| 2 | 错误重试策略冲突：模板说 max 1 retry then stop，agent 定义说 ~3 attempts | coding-fix.md 等 vs task-executor.md Complex Error Pause Flow | P1 |
| 3 | "Stop" 语义模糊：模板说 stop，agent 定义期望先创建 fix task 再 stop | 8 个 coding/validation/gate 模板 vs task-executor.md | P1 |
| 4 | coding-refactor Pre-check 失败后任务停留在 in_progress 无恢复路径 | coding-refactor.md Pre-check | P1 |
| 5 | Coverage 注入与 "no new tests" 指令直接矛盾 | coding-cleanup.md, coding-refactor.md + prompt.go resolveCoverage() | P1 |
| 6 | guide 说 fmt → WARNING (non-blocking)，模板说 fmt → Stop | guide vs 8 个 coding/validation/gate 模板 | P1 |
| 7 | guide 定义了 scope resolution（检查 project-type），但模板直接盲传 {{SCOPE}} | guide vs prompt.go + 模板 | P1 |
| 8 | task-executor.md 注释称 submit-task "via just test" 做指标收集，实际不是 | task-executor.md step 8 注释 | P1 |
| 9 | blocked 路径无替代 DONE 输出格式 | task-executor.md step 10 | P2 |
| 10 | resolveCoverage() 返回中文，注入英文模板，语言混杂 | prompt.go resolveCoverage() | P2 |

### Urgency

P0 问题意味着：当 `forge task submit` 发现 testsFailed > 0 并 auto-downgrade 为 blocked 时，agent 仍会执行 git-commit。这违反 submit-task SKILL.md 的 "do NOT commit" 约束，在 CI 中产生应该被阻断的 commit。

P1 问题导致：(1) agent 在重试策略上左右为难，(2) 模板 "stop" 后 agent 不知道该不该创建 fix task，(3) Pre-check 失败的任务永远卡在 in_progress，(4) cleanup/refactor 任务收到矛盾的 coverage 指令。

## Proposed Solution

逐一修复 10 处矛盾，使 task-executor agent 的 Execution Protocol 与 prompt 模板、guide、submit-task skill 语义完全对齐。

### Innovation Highlights

无创新，纯一致性修复。问题根因是多个 feature（typed-task-dispatch、scope-resolution、coverage-strategy）独立演进时，agent 协议未同步更新，且模板与 agent 定义各自维护了独立的错误处理策略。

## Requirements Analysis

### Key Scenarios

- **场景 1**（P0）：coding task 的 targeted tests 通过，但 `forge task submit` 的全量测试发现 testsFailed > 0 → auto-downgrade → 不应 commit → 但当前会 commit
- **场景 2**（P1）：lint 失败后，模板说 "max 1 retry then stop"，agent 定义说 "~3 attempts then fix task" → agent 不知道遵循哪个阈值
- **场景 3**（P1）：模板 step 3 static checks 失败说 "stop"，agent 定义的 Complex Error Pause Flow 期望先创建 fix task → 模板设计者和 agent 定义者对 "stop" 的预期行为不同
- **场景 4**（P1）：coding-refactor Pre-check 发现 git status 不干净 → "stop and report" → 任务永远 in_progress → run-tasks dispatcher 下次 claim 不会重新分配
- **场景 5**（P1）：coding-cleanup 任务 config 指定 percentage 策略 80% → 合成 prompt 说 "达到 80% 覆盖率. No new tests." → 直接矛盾
- **场景 6**（P1）：`just fmt` 在 toolchain 异常时失败 → agent Stop → 任务永久卡在 in_progress → 实际是 toolchain 问题而非代码问题
- **场景 7**（P1）：backend 项目的 task scope = "frontend" → 模板生成 `just compile frontend` → 应该 fallback 到 `just compile`
- **场景 8**（P1）：agent 读到 submit-task "via just test" 的注释 → 假设不需要自己跑 targeted tests → 但实际需要
- **场景 9**（P2）：blocked 任务输出 `DONE: ... | ✅ | <commit-hash> | ...` 但无 commit-hash → run-tasks dispatcher 解析出错
- **场景 10**（P2）：英文模板中出现中文 coverage 指令，影响 LLM 理解一致性

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
| **最小定点修复** | 本 proposal | 改动量小，风险可控 | 不解决结构性问题 | **Selected: 最小改动修复 10 处具体矛盾** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有矛盾点都有明确的文件和行号，改动是增量和确定性的。

### Resource & Timeline

10 处修改，预计 2-3 个 coding task 即可完成（按修改区域分组）。

### Dependency Readiness

无外部依赖。forgeconfig 和 prompt.go 的 API 已经就绪。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| task-executor 的 step 7 检查足够保证 blocked 安全 | Stress Test | Overturned: submit-task 的 auto-downgrade 发生在 step 7 之后 |
| 模板的错误处理与 agent 定义的错误处理一致 | 直接对比 | Overturned: 模板说 1 retry/stop，agent 说 ~3 attempts/fix task |
| 模板 "stop" 与 agent "PAUSE" 语义等价 | 语义分析 | Overturned: stop 无副作用，PAUSE 创建 fix task 并 block source |
| coverage 注入对所有 coding.* 类型都合理 | Assumption Flip | Overturned: cleanup/refactor 的 "no new tests" 与 percentage 策略矛盾 |
| guide 和 template 的 fmt 行为一致 | 直接对比 | Overturned: guide 说 WARNING，template 说 Stop |
| scope resolution 已在自动化路径中实现 | Occam's Razor | Overturned: 仅在 guide 自然语言中描述，未编码 |
| submit-task 注释准确描述了其行为 | 直接验证 | Overturned: submit-task 不运行 just test |
| Pre-check 失败有明确的状态转换路径 | 逻辑检查 | Overturned: coding-refactor Pre-check 只说 "stop and report"，不设 blocked |

## Scope

### In Scope

1. **task-executor.md**: 加 step 8.5（blocked 复查）+ 修正 submit-task 注释 + 补 blocked 输出格式 + 统一 retry 策略表述
2. **prompt/data/*.md**: 8 个模板的 fmt 行为从 Stop 改为 WARNING + coding-cleanup/coding-refactor 覆盖策略与 "no new tests" 对齐 + coding-refactor Pre-check 补 blocked 状态设置
3. **prompt.go**: 扩展 scope resolution 逻辑（检查 project-type）+ resolveCoverage() 返回英文文本
4. **错误处理对齐**: 统一模板 "stop" 与 agent Complex Error Pause Flow 的交互语义

### Out of Scope

- template 之间 convention 加载覆盖不一致的问题（P3，与 agent-prompt 契合无关）
- run-tasks.md 和 execute-task.md 的 dispatcher 层面改动
- guide 本身的修改
- 提交消息语言统一（与 agent-prompt 契合无关）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| prompt.go scope resolution 引入新 bug | L | M | 已有 scope_resolution_cli_test.go 覆盖 |
| fmt WARNING 导致格式问题未被发现 | L | L | All-Completed Hook 的 quality gate 仍会运行 fmt |
| blocked DONE 格式改变影响 run-tasks 解析 | L | M | run-tasks 不解析 DONE 输出（仅检查 forge task status） |
| retry 策略统一后 agent 行为变化 | L | L | 统一方向为 agent 定义的 ~3 attempts，比模板的 1 retry 更宽容 |
| coverage 逻辑改动影响已有测试 | L | M | prompt_test.go 已有 resolveCoverage 覆盖 |

## Success Criteria

- [ ] submit-task auto-downgrade 为 blocked 时，agent 不执行 git-commit（可构造 testsFailed > 0 场景验证）
- [ ] 模板的 retry 策略与 agent 定义的 Complex Error Pause Flow 阈值一致（统一为 ~3 attempts）
- [ ] 模板中的 "stop" 明确包含 "eval Complex Error Pause Flow" 语义
- [ ] coding-refactor Pre-check 失败时设置 blocked 状态（非 in_progress）
- [ ] coding-cleanup/coding-refactor 的 coverage 注入与 "no new tests" 不矛盾
- [ ] `just fmt` 失败时 agent 输出 warning 并继续，不 Stop
- [ ] backend 项目的 task scope="frontend" 时，模板生成 `just compile`（无 scope 参数）
- [ ] task-executor.md 中不再包含 "via just test" 的错误注释
- [ ] blocked 任务输出 `DONE: ... | blocked | ...` 格式（无 commit-hash）
- [ ] resolveCoverage() 返回英文文本，模板内无语言混杂

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
