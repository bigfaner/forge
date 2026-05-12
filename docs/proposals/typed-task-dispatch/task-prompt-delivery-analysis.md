---
created: 2026-05-11
author: faner + claude
status: Draft
---

# Task Prompt 投递方式分析

本文档记录 `typed-task-dispatch` 提案中 "策略指令如何传递给 task-executor subagent" 这一架构决策的深度分析过程与结论。

## 背景

proposal.md 定义了 `task prompt <id>` 命令，CLI 合成类型专属 agent prompt。核心问题：**合成的 prompt 如何投递给 task-executor subagent？**

proposal 原始方案：mainSession 调用 `task prompt <id>`，捕获 stdout，作为 `Agent()` 的 prompt 参数传入。

本分析对三种投递方式进行了对比。

## Claude Code Agent Prompt 层级

理解投递方式的前提是明确 Claude Code 的 prompt 层级架构。

当调用 `Agent(subagent_type="forge:task-executor", prompt="...")` 时，subagent 看到的层级：

| 层级 | 来源 | 角色 |
|------|------|------|
| **System prompt** | agent 定义文件的 body（task-executor.md） | 身份、约束、行为规则。subagent 不继承 mainSession 的系统提示词 |
| **First user message** | Agent() 的 prompt 参数 | 任务指令 |
| **Tool output** | subagent 内部的 Bash/Read 等工具调用结果 | 信息来源，通常为二级参考 |

关键约束：subagent 只看到自己的 agent 定义 + prompt 参数，与 mainSession 完全隔离。

## 三种投递方式

### 方式 A：stdout 传递（proposal 原始方案）

mainSession 调用 `task prompt`，捕获完整输出，作为 prompt 参数传入 subagent。

```
run-tasks:
  1. Bash("task prompt T-impl-1")          → stdout 留在 mainSession 上下文
  2. Agent(prompt=<stdout>)                 → 同一文本进入 subagent 上下文
```

**task-executor.md**：仅包含硬约束（~40 行）。

```
System: [硬约束]
User:   [task prompt 完整输出 — ~500 tokens]
```

### 方式 B：文件投递

`task prompt` 写入文件，prompt 参数仅传递文件路径。

```
run-tasks:
  1. Bash("task prompt T-impl-1")          → 写入 .forge/task-prompts/T-impl-1.md，无 stdout
  2. Agent(prompt="Execute T-impl-1. Read .forge/task-prompts/T-impl-1.md")
```

**task-executor.md**：硬约束 + 文件读取指令（~45 行）。

```
System: [硬约束 + "读取策略文件并严格执行"]
User:   "Execute task T-impl-1. Strategy: .forge/task-prompts/T-impl-1.md" (~60 tokens)
Agent:  Read(.forge/task-prompts/T-impl-1.md) → [完整策略]
```

### 方式 C：系统指令委托

prompt 参数仅传 task ID，subagent 内部调用 `task prompt` 获取策略。

```
run-tasks:
  1. Agent(prompt="Execute task T-impl-1")  → ~20 tokens
```

**task-executor.md**：硬约束 + 执行协议（~50 行）。

```
System: [硬约束 + Execution Protocol: 调用 task prompt，严格遵循输出]
User:   "Execute task T-impl-1" (~20 tokens)
Agent:  Bash("task prompt T-impl-1") → [完整策略]
Agent:  按策略执行
```

## 对比分析

### Token 效率

策略文本 ~500 tokens / 任务，10 个任务的累计成本：

| 维度 | 方式 A（stdout 传递） | 方式 B（文件投递） | 方式 C（系统指令委托） |
|------|---------------------|------------------|---------------------|
| mainSession 增量 | ~5000 tokens（task prompt stdout） | ~600 tokens（路径引用） | ~200 tokens（task ID） |
| subagent 增量 | ~500 tokens（prompt 参数） | ~560 tokens（路径 + 文件内容） | ~500 tokens（task prompt stdout） |
| 策略出现次数 | **2 次**（mainSession + subagent） | 1 次（subagent 中的文件内容） | **1 次**（subagent 中的 task prompt stdout） |
| 总 token 消耗 | 最高 | 中等 | **最低** |

方式 A 中策略文本出现两次是结构性浪费：mainSession 的 stdout 留在上下文中无法被复用，仅作为传递媒介。

### run-tasks 复杂度

| 维度 | 方式 A | 方式 B | 方式 C |
|------|--------|--------|--------|
| dispatch 步骤 | 2 步（调用 task prompt + 传递输出） | 2 步（调用 task prompt 写文件 + 传递路径） | **1 步（只传 task ID）** |
| 错误处理 | mainSession 处理 task prompt 失败 | mainSession 处理 task prompt 失败 | subagent 内部处理（record as blocked） |
| 文件 I/O | 无 | 需要读写文件 | **无** |

### 策略指令的层级与权威性

这是三种方式最核心的架构差异。

**方式 A**：策略在 prompt 参数（user-level），约束在 agent 定义（system-level）。标准层级，约束天然覆盖策略。

**方式 B**：策略在文件内容（tool output），被 system 指令 "读取并执行" 提升权威性。层级略模糊。

**方式 C**：策略在 task prompt stdout（tool output），被 system 指令 "严格遵循输出" 提升为事实上的 system-level 指令。

方式 C 的隐含风险：如果策略指令与硬约束冲突，两者都被 system 背书，优先级不明确。

**缓解措施**：task prompt 模板由人工审核，确保不输出与硬约束冲突的指令。这是可执行的约束——模板数量固定（12 个），变更时通过 code review 验证。

### 上下文压缩风险

subagent 长时间执行时，上下文可能被自动压缩。

| 方式 | 策略所在位置 | 压缩优先级 | 恢复能力 |
|------|------------|-----------|---------|
| A | 初始 user message | 低（通常保留） | 无法恢复（prompt 参数只出现一次） |
| B | Read tool output | 中 | **可恢复（重新 Read 文件）** |
| C | Bash tool output | 高 | **可恢复（重新调用 task prompt）** |

方式 B 和 C 的策略可恢复，因为数据源持久化（文件 / CLI 命令）。方式 A 虽然压缩风险最低，但一旦被压缩则无法恢复。

方式 C 的 Execution Protocol 可加入恢复指令：

```markdown
## Execution Protocol
1. Extract task ID from the task prompt
2. Run `task prompt <TASK_ID>` to get execution strategy
3. Follow every step exactly
4. If you lose track of your strategy, re-run `task prompt <TASK_ID>`
5. If task prompt fails (non-zero exit), record task as blocked and stop
6. After all steps, call forge:record-task
```

### 错误处理

| 失败场景 | 方式 A / B | 方式 C |
|---------|-----------|--------|
| task prompt 失败 | mainSession 拦截，不 dispatch subagent | subagent 内部 record as blocked |
| subagent 资源消耗 | 零 | 已消耗一次 Agent 调用 |
| 可恢复性 | run-tasks 直接重试 | task claim 重新分配 |

task prompt 运行时失败是极端情况（type 未知已在 validate 拦截，模板缺失在开发期发现），消耗一次 Agent 调用的代价可接受。

### 额外 API 轮次

| 方式 | subagent 首次行为 | 额外轮次 |
|------|-----------------|---------|
| A | 直接执行策略 | 0 |
| B | Read 策略文件 → 执行 | +1 |
| C | Bash(task prompt) → 执行 | +1 |

方式 B 和 C 各多一轮 API 交互。对总执行时间影响微小。

## task prompt 的上下文自足性

方式 C 的前提是 `task prompt <id>` 只需要 task ID 即可完成全部合成。验证了 CLI 实现：

**Synthesize 函数接收参数**（`prompt.go:37`）：

```go
type SynthesizeOpts struct {
    ProjectRoot     string
    FeatureSlug     string
    TaskID          string
    FixRecordMissed bool
}
```

**所有占位符由 CLI 内部自动填充**：

| 占位符 | 来源 |
|--------|------|
| `{{TASK_ID}}` | opts.TaskID |
| `{{TASK_FILE}}` | index.json → task.File → 拼接路径 |
| `{{SCOPE}}` | index.json → task.Scope |
| `{{FEATURE_SLUG}}` | opts.FeatureSlug |
| `{{PHASE_SUMMARY}}` | PhaseDetect() 扫描 index.json + 检查磁盘文件 |

CLI 从 `.forge/state.json` 读取 feature slug，从 index.json 读取 task 元数据。agent 只需提供 task ID，CLI 完成其余所有工作。

**边界条件**：如果未来模板需要 agent 运行时才能确定的动态信息（如"上次失败的错误信息"），CLI 无法预填充，方式 C 会受限。当前 12 个模板全部使用静态上下文，此限制不生效。

## 结论

### 推荐方案：方式 C（系统指令委托）

在模板人工审核的前提下，方式 C 在关键维度上优于 proposal 原始方案：

| 维度 | 方式 C 优势 |
|------|-----------|
| Token 效率 | 策略文本只出现 1 次（方式 A 出现 2 次），mainSession 增量从 ~500 降至 ~20 tokens/任务 |
| run-tasks 复杂度 | dispatch 从 2 步缩减为 1 步，run-tasks 不需要理解 task prompt 输出 |
| 策略可恢复性 | 上下文压缩后可重新调用 task prompt 恢复（方式 A 无法恢复） |
| 可测试性 | 不变——task prompt 仍是 CLI 命令，Go 单元测试覆盖所有 type |
| 文件 I/O | 无（方式 B 需要文件读写） |

### 对 proposal.md 的影响

方式 C 对 proposal 的改动范围有限：

| 改动项 | 说明 |
|--------|------|
| `task prompt` 命令 | **不变**——输出到 stdout，CLI 内部完成所有上下文填充 |
| task-executor.md | 增加 Execution Protocol 块（~10 行），总量从 ~40 行增至 ~50 行 |
| run-tasks | dispatch 逻辑简化：只传 task ID，不再调用 task prompt |
| 类型体系 / 模板 / task migrate / task validate | **不变** |

### task-executor.md 预期结构

```markdown
## Hard Constraints
- ONE TASK PER INVOCATION
- record-task is mandatory before stopping
- No background tasks
- Maximum 3 subagent calls
- STOP after record-task

## Execution Protocol
1. Extract task ID from the task prompt
2. Run `task prompt <TASK_ID>` to get execution strategy
3. Follow every step exactly
4. If you lose track of your strategy, re-run `task prompt <TASK_ID>`
5. If task prompt fails (non-zero exit), record task as blocked and stop
6. After all steps, call forge:record-task
```

### 特殊路由

| 任务类型 | 路由方式 | 说明 |
|---------|---------|------|
| `test-pipeline.eval-cases` | mainSession 直接执行 `task prompt <id>` | subagent 无法再 spawn subagent（平台限制） |
| `--fix-record-missed` | prompt 参数区分：`"Fix record for task <id>"` | Execution Protocol 增加："若任务为 fix record，调用 `task prompt <id> --fix-record-missed`" |
| 其余所有类型 | `Agent(prompt="Execute task <id>")` | 统一 dispatch |
