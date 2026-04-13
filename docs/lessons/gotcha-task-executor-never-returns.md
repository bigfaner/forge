# task-executor Subagent 不返回主会话

## Problem

`/run-tasks` dispatcher 通过 `Agent(subagent_type="task-executor")` 派发任务后，task-executor 在完成一个任务（Step 5 commit）后不返回主会话，而是**自主继续 claim 并执行后续任务**（3.2 → 3.3 → 3.4 → 3.5 → 3.6），直到所有任务完成。

dispatcher 设计为：claim → dispatch → verify → loop。但 subagent 不返回导致 dispatcher 的 verify 和 loop 逻辑完全失效。

## Root Cause

`task-executor.md` agent 定义中**缺少明确的终止指令**：

1. **无"只执行一个任务"的约束** — 5 步工作流结束后没有说"停止"或"返回控制权"
2. **agent 拥有继续执行的工具** — `allowed_tools` 包含 `Bash`（可运行 `task claim`），所以它能自主认领下一个任务
3. **LLM 的"效率"倾向** — Sonnet 看到 index.json 中还有可用任务，倾向于继续完成而非停止
4. **dispatcher prompt 只传单个任务信息** — 但没有禁止 agent 自己去读 index.json 并 claim 更多任务

关键代码位置：`plugins/zcode/agents/task-executor.md` 的 Step 5 之后没有任何终止指令。

## Solution

在 `task-executor.md` 的 Core Rules 中添加明确约束：

```markdown
<EXTREMELY-IMPORTANT>
5. Execute EXACTLY ONE task per invocation - after Step 5, STOP immediately
6. Do NOT run "task claim" or read index.json after completing your task
7. Return control to the dispatcher - your job is done after commit
</EXTREMELY-IMPORTANT>
```

同时在 Step 5 之后添加：

```markdown
### STOP

After Step 5, your task is complete. Do NOT:
- Run `task claim`
- Read the next task file
- Continue with any additional work

Output your final DONE line and STOP.
```

## Key Takeaway

Agent 定义必须有**显式终止条件**。LLM agent 不会自动停在"合理"的边界——如果它有工具和能力继续，它通常会继续。关键防御措施：

1. **声明执行边界** — "Execute EXACTLY ONE task" 比 "focused task executor" 更有效
2. **禁止越界操作** — 明确禁止 agent 自主 claim 新任务
3. **工具权限与行为约束配合** — 给了 Bash 工具就要约束其使用范围，否则 agent 会用它做设计外的事
4. **验证 loop 的脆弱性** — dispatcher 的 verify 步骤依赖 subagent 返回，如果 subagent 不返回，整个 loop 失效且无报错
