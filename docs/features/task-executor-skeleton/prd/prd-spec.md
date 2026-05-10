---
feature: "task-executor-skeleton"
status: Draft
db-schema: "no"
---

# task-executor-skeleton — PRD Spec

> PRD Spec: defines WHAT the feature is and why it exists.

## Background

### Why (Reason)

task-executor 硬编码了 TDD 工作流（Step 2-3），所有非 MAIN_SESSION 任务走同一条路径。执行型任务（如 T-test-3 "Run e2e Tests"）被强行塞入 TDD 循环，导致 14 分钟无效重试。`noTest` 标志作为权宜之计混入模板 frontmatter，语义模糊——它不表示"不需要测试"，而是"不走 TDD 循环"，导致新模板配置频繁出错。

### What (Target)

引入 `## Execution Workflow` 机制：每个任务模板声明自己的执行步骤，task-executor 读取并执行。同时完全移除 `noTest`，消除歧义。task-executor 变成纯骨架：读任务 → 执行 workflow → 记录 → 提交。

### Who (Users)

- **Template Author** — 编写任务模板的开发者，需要在模板中添加 Execution Workflow 段落
- **Task-executor Agent** — 读取并执行 workflow 的 AI agent，从硬编码 TDD 变为通用骨架
- **Forge Maintainer** — 维护 harness 的开发者，需要从代码中清理 noTest/NO_TEST

## Goals

| Goal | Metric | Notes |
|------|--------|-------|
| 消除执行型任务的无效重试 | T-test-3 执行时间 < 5 min（当前 ~14 min） | 核心收益 |
| 消除 noTest 歧义 | grep `noTest` / `NO_TEST` 零匹配 + 语义审查无残留分支 | 完全移除 |
| 统一任务执行机制 | 所有 16 个模板包含 `## Execution Workflow` | 10 breakdown + 6 quick |
| 向后兼容 | 无 Execution Workflow 的旧任务回退到 TDD，行为不变 | fallback 保障 |

## Scope

### In Scope

- [ ] `agents/task-executor.md`: Step 2-3 合并为"执行 workflow"（从任务文件读取），删除 NO_TEST input
- [ ] 所有 breakdown-tasks 模板（10 个，不含 manifest-update-tasks.md 和 eval-test-cases.md）: 添加 `## Execution Workflow`，删除 `noTest`
- [ ] 所有 quick-tasks 模板（6 个，不含 manifest-quick.md）: 添加 `## Execution Workflow`，删除 `noTest`
- [ ] `index.schema.json` (breakdown + quick): 删除 `noTest` 字段定义
- [ ] `commands/run-tasks.md`: 删除 NO_TEST 从 claim 解析和 dispatch
- [ ] `commands/execute-task.md`: 删除 NO_TEST 相关内容
- [ ] `task-cli`: 删除 noTest 字段和相关逻辑（types.go, record.go, claim 输出）
- [ ] `skills/record-task/SKILL.md`: 删除 noTest 引用
- [ ] `skills/quick-tasks/SKILL.md`: 删除 `--no-test` 标志
- [ ] `skills/consolidate-specs/SKILL.md`: 删除 noTest 引用

### Out of Scope

- 模板内容重写（仅添加 Execution Workflow 段落，不改动现有 Implementation Notes）
- mainSession 任务的路由逻辑（已由 dispatcher 处理，不涉及 task-executor）
- `## Execution Workflow` 的模板化/标准化（未来优化，将常用 workflow 提取为可引用片段）
- `task add --template` 自动注入 Execution Workflow（未来优化）

## Flow Description

### Business Flow Description

task-executor 收到任务后，读取任务文件并检测是否包含 `## Execution Workflow` 标题。若存在且正文非空，将其注入到 agent prompt 的 Step 2 指令区域，替换硬编码的 TDD 步骤。若不存在，回退到当前 TDD + Quality Gate 步骤。若标题存在但正文为空，视为配置错误，记录警告并回退到 TDD。

**成功路径**：执行完成后统一进入 Step 3（记录 + 提交），任务状态设为 `completed`。

**失败路径**：
- **任务文件不可读**：Step 1 阶段若任务文件缺失或格式损坏（无法解析 frontmatter），记录错误日志，任务状态设为 `failed`，不进入 Step 2。跳过执行，直接提交一条错误记录。
- **Agent 执行超时**：若 agent 在执行 workflow 步骤时超过 task-executor 内置超时阈值，强制终止当前执行，任务状态设为 `failed`，记录超时信息。
- **Workflow 执行失败**：若 agent 在执行过程中遇到错误（命令执行失败、测试框架崩溃、外部依赖不可达），任务状态设为 `failed`。若 workflow 本身声明了失败处理指令（如 T-test-3 的"创建 fix task"），agent 按指令执行后停止；若无显式失败指令，agent 记录失败原因并停止，不进入 TDD 循环重试。
- **部分完成**：若 workflow 包含多个步骤且中途失败，agent 记录已完成步骤和失败点，任务状态设为 `failed`，失败记录包含已执行步骤摘要和失败原因。

**终端状态**：`completed`（执行成功，进入 Step 3 提交）、`failed`（上述任一失败路径触发）。无 `partial` 状态——部分完成归入 `failed` 并附带已执行步骤信息。

### Business Flow Diagram

```mermaid
flowchart TD
    Start([task-executor 收到任务]) --> S1[Step 1: 读取任务文件]
    S1 --> FileOK{任务文件可读<br>且 frontmatter 合法?}
    FileOK -->|No| FailRead[记录错误日志<br>状态 = failed]
    FailRead --> EndFail([终端: failed])
    FileOK -->|Yes| HasWF{任务文件包含<br>## Execution Workflow ?}
    HasWF -->|Yes| EmptyWF{Workflow 正文为空?}
    EmptyWF -->|Yes| Warn[记录配置错误警告<br>回退到 TDD]
    EmptyWF -->|No| Inject[注入 workflow 正文到<br>Step 2 指令区域]
    Inject --> Execute[Agent 执行 workflow 步骤]
    HasWF -->|No| Fallback[回退到 TDD + Quality Gate]
    Fallback --> Execute
    Warn --> Execute
    Execute --> ExecOK{执行结果?}
    ExecOK -->|成功| S3[Step 3: 记录结果 + 提交<br>状态 = completed]
    S3 --> EndOK([终端: completed])
    ExecOK -->|失败/超时| FailExec[记录失败原因 + 已执行步骤摘要<br>状态 = failed]
    FailExec --> EndFail
```

## Functional Specs

> This feature has no UI surface. All changes are internal to the forge harness (agent prompts, task templates, CLI code, skill docs).

### Related Changes

| # | Component | Change Point | Updated Logic |
|---|-----------|-------------|---------------|
| 1 | task-executor.md | Step 2-3 硬编码 TDD → 读取 Execution Workflow | 段落检测 + 注入 + fallback |
| 2 | task-cli (Go) | noTest 字段及相关逻辑 | types.go 删除字段、record.go 删除逻辑、claim 输出移除 NO_TEST |
| 3 | run-tasks.md / execute-task.md | NO_TEST 解析和传递 | 删除 claim 输出解析和 dispatch prompt 中的 NO_TEST |
| 4 | 16 任务模板 | 添加 Execution Workflow + 删除 noTest frontmatter | 每个模板声明自己的执行步骤 |
| 5 | index.schema.json × 2 | noTest 字段定义 | 删除字段定义，验证模板合规 |
| 6 | 3 个 skill 文档 | noTest / --no-test 引用 | record-task, quick-tasks, consolidate-specs |

## Other Notes

### Compatibility Requirements

- 无 `## Execution Workflow` 的旧任务文件自动 fallback 到 TDD + Quality Gate，行为不变
- 回退逻辑是 task-executor 的内置行为，无需模板声明

### Quality Requirements

- task-executor 须严格遵循 workflow 指令，使用 `<EXTREMELY-IMPORTANT>` / `<HARD-RULE>` 标签保障（与 "ONE TASK PER INVOCATION" 规则机制一致，该规则自部署以来零违规）
- noTest 移除后须 grep 确认零匹配 + 逐文件审查条件分支无残留
- 每个 Execution Workflow 模板须包含明确的结束条件，不允许开放式指令

---

## Quality Checklist

- [x] Is the requirement title accurate and descriptive
- [x] Does the background include all three elements: reason, target, users
- [x] Are the goals quantified
- [x] Is the flow description complete
- [x] Does the business flow diagram exist (Mermaid format)
- [x] Is there any ambiguous or vague wording — none
- [x] Is the spec actionable and verifiable
