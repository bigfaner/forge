---
created: "2026-06-04"
topic: "Claude Code Dynamic Workflows"
mode: "single-tech-deep-dive"
candidates: []
dimensions: [overview-positioning, architecture-core-concepts, learning-curve, ecosystem-community]
---

# Claude Code Dynamic Workflows 研究报告

## Overview

Claude Code 的 "dynamic workflows" 不是单一功能，而是一个由 6 个互补机制组成的**多层次动态编排体系**，覆盖从简单的定时轮询到大规模多代理脚本化编排的完整光谱。

**Research mode:** 单技术深潜

**Key question:** Claude Code 的动态工作流体系包含哪些核心机制？它们如何协同运作？对 Forge 项目有何借鉴价值？

## Research Background & Objectives

Claude Code (Anthropic 的 CLI 工具) 在 2025-2026 年间逐步引入了一系列动态编排能力。这些能力分散在不同的文档页面中，没有一个统一的入口解释它们之间的关系。本次调研旨在系统性地梳理这些机制的架构、运作原理和适用场景。

### Research Scope

| Dimension | Value |
|---|---|
| Topic | Claude Code Dynamic Workflows |
| Mode | 单技术深潜 |
| Dimensions covered | 概述与定位、架构与核心概念、学习曲线、生态与社区 |
| Candidates | N/A |
| Project adaptation | 是 |

---

## 深潜：Claude Code Dynamic Workflows

### 概述与定位

Claude Code 的 "dynamic workflows" 是一个**伞形概念**，涵盖 6 个独立的编排机制，从轻量到重量级排列：

| 机制 | 定位 | 编排者 | 适用规模 |
|---|---|---|---|
| **Hooks** | 生命周期事件触发器 | 事件驱动 | 单次工具调用 |
| **Skills** | 可复用的 prompt 扩展 | Claude (turn-by-turn) | 单会话内的步骤序列 |
| **Subagents** | 专业化工作代理 | Claude (turn-by-turn) | 几个委派任务 |
| **/loop & Cron** | 定时任务调度 | 时间驱动 (cron/wakeup) | 持续轮询/监控 |
| **Dynamic Workflows** | 脚本化多代理编排 | JavaScript 脚本 (runtime) | 几十到数百个代理 |
| **Agent Teams** | 多实例协作团队 | Team Lead (peer-to-peer) | 3-5 个并行工作流 |

**核心设计哲学差异：** 在 Skills、Subagents 中，Claude 自己是编排者——它逐轮决定下一步做什么，所有中间结果都落在 context window 中。而 Dynamic Workflows **将编排计划移入代码**——JavaScript 脚本持有循环、分支和中间结果，Claude 的 context 只持有最终答案。

> "A workflow moves the plan into code. With subagents, skills, and agent teams, Claude is the orchestrator: it decides turn by turn what to spawn or assign next, and every result lands in a context window. A workflow script holds the loop, the branching, and the intermediate results itself."
> — Claude Code 官方文档

---

### 架构与核心概念

#### 1. Hooks：事件驱动的自动化骨架

Hooks 是 Claude Code 最底层的自动化机制——用户定义的 shell 命令、HTTP 端点或 LLM prompt，在 Claude Code 生命周期的特定节点自动执行。

**完整的 Hook 事件列表（30+ 种），按生命周期排列：**

| 事件 | 触发时机 | 可阻止 | 典型用途 |
|---|---|---|---|
| `SessionStart` | 会话开始或恢复 | 否 | 注入上下文、环境初始化 |
| `Setup` | `--init-only` 或 `-p --init` 启动 | 否 | 环境准备 |
| `InstructionsLoaded` | CLAUDE.md / rules 加载后 | 否 | 动态指令补充 |
| `UserPromptSubmit` | 用户提交 prompt 时 | 是 | 拦截/修改用户输入 |
| `UserPromptExpansion` | 斜杠命令展开前 | 是 | 修改命令参数 |
| `PreToolUse` | 工具执行前 | 是 | 拦截不安全操作 |
| `PermissionRequest` | 权限对话框出现时 | 是 | 自动审批/拒绝 |
| `PermissionDenied` | 分类器拒绝工具调用 | 否 | 日志/告警 |
| `PostToolUse` | 工具成功完成后 | 是 | 自动格式化、lint |
| `PostToolUseFailure` | 工具调用失败后 | 否 | 错误日志 |
| `PostToolBatch` | 一批并行工具调用完成后 | 是 | 批量后处理 |
| `MessageDisplay` | 消息流式显示到屏幕时 | 否 | 自定义输出 |
| `Notification` | 通知事件 | 否 | 自定义通知行为 |
| `SubagentStart` | 子代理创建时 | 否 | 代理初始化 |
| `SubagentStop` | 子代理完成时 | 是 | 代理清理/质量检查 |
| `TaskCreated` | 任务创建 | 是 | 验证任务合法性 |
| `TaskCompleted` | 任务完成 | 是 | 质量门禁 |
| `TeammateIdle` | 团队成员即将空闲 | 是 | 强制继续工作 |
| `Stop` | 主代理完成响应 | 是 | 质量门禁、feature 完成 |
| `StopFailure` | API 错误导致回合结束 | 否 | 错误恢复 |
| `PreCompact` | 上下文压缩前 | 是 | 保存关键信息 |
| `PostCompact` | 上下文压缩后 | 否 | 压缩后处理 |
| `ConfigChange` | 配置文件会话中变更 | 是 | 配置热更新 |
| `CwdChanged` | 工作目录变更 | 否 | 目录切换响应 |
| `FileChanged` | 监视文件在磁盘上变更 | 否 | 文件变更响应 |
| `WorktreeCreate` | 工作树创建 | 是 | 工作树初始化 |
| `WorktreeRemove` | 工作树移除 | 否 | 工作树清理 |
| `Elicitation` | MCP 请求用户输入时 | 是 | MCP 交互控制 |
| `ElicitationResult` | MCP 响应发回前 | 是 | 响应过滤 |
| `SessionEnd` | 会话终止 | 否 | 清理、状态保存 |

**5 种 Hook Handler 类型：**
1. `command` — Shell 命令（stdin JSON / stdout JSON / exit code）
2. `http` — POST 到 HTTP 端点
3. `mcp_tool` — 调用 MCP 服务器工具
4. `prompt` — 单轮 LLM 评估（yes/no 决策）
5. `agent` — 生成子代理验证后返回决策

**关键设计原则：**
- Hook 通过 `exit code 2` 阻止操作执行，stderr 反馈给 Claude
- 同一事件的所有匹配 hooks **并行运行**，结果合并（deny > defer > ask > allow）
- Matcher 支持正则匹配工具名（如 `"Edit|Write"`）
- Hooks **先于权限模式检查**，`PreToolUse` 返回 `deny` 即使 `bypassPermissions` 模式也有效
- Hooks 只能收紧不能放松限制（返回 `allow` 不能覆盖已有的 `deny`）
- Stop hook 连续阻止上限 8 次（可通过 `CLAUDE_CODE_STOP_HOOK_BLOCK_CAP` 调整）
- Hooks 可定义在 subagent frontmatter（代理级）、settings.json（项目/用户级）或 managed policy（组织级）

#### 2. Subagents：专业化工作代理

每个 subagent 运行在独立的 context window 中，拥有自定义系统 prompt、特定工具访问权限和独立权限。

**内置代理：**

| 代理 | 模型 | 用途 |
|---|---|---|
| Explore | Haiku | 快速只读搜索和分析代码库 |
| Plan | 继承 | plan mode 下的代码库调研 |
| General-purpose | 继承 | 需要探索和修改的复杂多步任务 |
| claude-code-guide | Haiku | 回答关于 Claude Code 功能的问题 |
| statusline-setup | Sonnet | 配置状态栏 |

**自定义代理的关键 frontmatter 字段：**
- `name` / `description`（必需）
- `tools` / `disallowedTools`（工具白名单/黑名单）
- `model`（sonnet / opus / haiku / inherit）
- `permissionMode`（权限模式）
- `maxTurns`（最大轮次）
- `skills`（预加载的技能）
- `mcpServers`（MCP 服务器配置）
- `hooks`（代理级 hooks）
- `memory`（持久化记忆：user / project / local）
- `isolation`（worktree 隔离）
- `background`（后台运行）

**Fork 模式（`CLAUDE_CODE_FORK_SUBAGENT=1`）：**
Fork 继承父会话的完整对话历史，而非从零开始。因为 fork 的系统 prompt 和工具定义与父级完全相同，其首个请求可复用父级的 prompt cache，比生成全新 subagent 更经济。

#### 3. /loop & Cron：时间驱动的持续执行

`/loop` 是 Claude Code 内置的定时任务能力，有三种运行模式：

| 提供内容 | 示例 | 行为 |
|---|---|---|
| 间隔 + prompt | `/loop 5m check the deploy` | 固定间隔执行 |
| 仅 prompt | `/loop check the deploy` | **动态模式**——Claude 自选间隔 |
| 无或仅间隔 | `/loop` / `/loop 15m` | 内置维护 prompt |

**动态模式的核心机制——ScheduleWakeup：**
- 当 `/loop` 省略间隔时，Claude 在每次迭代后自主选择 1 分钟到 1 小时的延迟
- 选择基于当前观察：构建快完成时短等，没有待处理事项时长等
- 底层使用 `ScheduleWakeup` 工具（v2.1.101 引入）
- 延迟和原因在每次迭代结束时打印

**缓存考量：** Anthropic 的 prompt cache TTL 为 5 分钟。睡眠超过 300 秒意味着下次唤醒时读取完整上下文（无缓存命中）。因此：
- < 270s：缓存保持热态
- 300s-3600s：承受缓存 miss
- 建议空闲时默认 1200s-1800s（20-30 分钟）

**自定义 loop.md：**
- 项目级：`.claude/loop.md`
- 用户级：`~/.claude/loop.md`
- 替换内置维护 prompt，编辑后下一次迭代即生效

**Cron 工具三件套：**
- `CronCreate`：创建定时任务（5 字段 cron 表达式）
- `CronList`：列出所有定时任务
- `CronDelete`：取消定时任务
- 7 天自动过期，最多 50 个并发任务

**调度选项对比：**

| | Cloud (Routines) | Desktop Scheduled | `/loop` |
|---|---|---|---|
| 需要机器在线 | 否 | 是 | 是 |
| 需要打开会话 | 否 | 否 | 是 |
| 跨重启持久化 | 是 | 是 | `--resume` 恢复未过期的 |
| 访问本地文件 | 否（fresh clone） | 是 | 是 |

#### 4. Dynamic Workflows：脚本化大规模编排

这是 Claude Code 中"dynamic workflows"最精确的含义——**一个 JavaScript 脚本编排大规模 subagents**。

**核心特征：**
- Claude 为你的任务编写 JavaScript 编排脚本
- 运行时在隔离环境中执行脚本，与对话分离
- 中间结果保存在脚本变量中，不占用 Claude 的 context
- 会话保持响应，代理在后台工作

**触发方式：**
1. 在 prompt 中包含关键词 `ultracode`（或说"用 workflow"）
2. 设置 `/effort ultracode`——Claude 自动为每个实质性任务规划 workflow
3. 运行保存的 workflow 命令（如 `/deep-research`）

**运行时约束：**

| 约束 | 值 | 原因 |
|---|---|---|
| 并发代理上限 | 16（CPU 有限时更少） | 限制本地资源 |
| 总代理上限 | 1,000/run | 防止无限循环 |
| 运行中用户输入 | 不支持 | 仅代理权限提示可暂停 |
| 脚本的文件系统访问 | 无 | 代理读写运行命令，脚本只协调代理 |

**保存和复用：**
- 运行成功后按 `s` 保存为命令
- 保存到 `.claude/workflows/`（项目级）或 `~/.claude/workflows/`（用户级）
- 通过 `/<name>` 重新运行，支持 `args` 参数传递

**监控：**
- `/workflows` 列出运行中和已完成的 workflow
- 进度视图显示每个 phase 的代理数量、token 总量和耗时
- 支持暂停、恢复、停止单个代理或整个 workflow

**JavaScript Workflow Script API 深潜：**

Workflow 脚本是纯 JavaScript，运行在 Node `vm` 沙箱中。脚本必须以 `export const meta` 开头，随后是使用全局 API 函数的脚本体。

**1. 脚本结构模板：**

```javascript
export const meta = {
  name: 'workflow-name',           // 必需，kebab-case
  description: '一句话描述',       // 必需
  whenToUse: 'Claude 何时使用此 workflow', // 可选
  phases: [                        // 可选，进度视图的阶段大纲
    { title: 'Phase 1', detail: '做什么' },
    { title: 'Phase 2', detail: '做什么' },
  ],
}

// 脚本体：使用 agent()、parallel()、pipeline()、phase() 编排
```

**2. 全局 API 函数：**

| 函数 | 签名 | 说明 |
|---|---|---|
| `agent(prompt, opts?)` | `string × {label?, schema?, phase?} → Promise<object\|string>` | 生成一个隔离 subagent，返回其最终文本或（带 `schema` 时）经过验证的结构化对象 |
| `parallel(thunks)` | `Array<() => Promise<T>> → Promise<Array<T>>` | 并发运行 thunk 数组，结果按输入顺序返回 |
| `pipeline(items, ...stages)` | `Array<T> × ...((prev, original, index) => Promise<U>) → Promise<Array<?>>` | 将每个 item 通过连续阶段传递，同时 item 间扇出并行 |
| `phase(title)` | `string → void` | 标记当前阶段，用于进度视图分组 |
| `log(message)` | `string → void` | 追加 workflow 级别日志行 |

**3. 全局变量：**

| 变量 | 类型 | 说明 |
|---|---|---|
| `args` | `any` | 用户传入的参数（`Workflow({name, args})`） |
| `cwd` | `string` | 当前工作目录 |
| `budget` | `{total, spent(), remaining()}` | Token 预算追踪器 |

**4. `agent()` 选项详解：**

```javascript
const result = await agent('Prompt text for the subagent', {
  label: 'human-readable label',   // 进度视图中显示
  phase: 'Phase Name',             // 指定归属阶段
  schema: {                        // JSON Schema → 强制结构化输出
    type: 'object',
    required: ['field1'],
    properties: {
      field1: { type: 'string' },
      field2: { type: 'array', items: { type: 'string' } },
    },
  },
})
// result 是经过 schema 验证的对象
// 不传 schema 时 result 是字符串
```

**5. `parallel()` 模式：**

```javascript
// 并行运行多个代理
const results = await parallel([
  () => agent('Search angle 1', { label: 'search-1', schema: SEARCH_SCHEMA }),
  () => agent('Search angle 2', { label: 'search-2', schema: SEARCH_SCHEMA }),
  () => agent('Search angle 3', { label: 'search-3', schema: SEARCH_SCHEMA }),
])
// results = [result1, result2, result3]，按输入顺序
```

**6. `pipeline()` 模式（搜索→去重→抓取链）：**

```javascript
const output = await pipeline(
  items,                                    // 输入数组
  (item) => agent('Stage 1 prompt', {...}), // 阶段 1：处理每个 item
  (stage1Result) => {                       // 阶段 2：基于阶段 1 结果
    return parallel(                        //   可以在阶段内再并行
      stage1Result.subItems.map(sub => () =>
        agent('Stage 2 prompt for ' + sub, {...})
      )
    )
  },
)
```

pipeline 的关键特性：**item 间并行、阶段间串行**。每个 item 独立通过所有阶段，不同 item 的阶段 1 可以并行运行。

**7. `/deep-research` 实战脚本分析（官方内置 workflow）：**

`/deep-research` 是官方内置的 workflow，其脚本展示了完整的 API 使用模式：

```
Phase 1 (Scope):  1 agent → 分解问题为 5 个搜索角度
Phase 2 (Search+Fetch):  pipeline(
    5 angles → 5 并行 search agents
    → URL 去重 + 预算管理
    → 每个新 URL 并行 fetch agents（≤15）
  )
Phase 3 (Verify):  每个声明 3 个并行 adversarial agents（3-vote）
  → ≥2/3 refutations → 杀死声明
Phase 4 (Synthesize): 1 agent → 合并、去重、排序、引用
```

关键设计模式：
- **JSON Schema 强制结构化输出**：每个 agent 调用都有 schema，确保阶段间数据传递类型安全
- **URL 去重 + 预算管理**：在 pipeline 阶段间维护 `seen` Map 和 `fetchSlots` 计数器
- **Adversarial verification**：3 个独立代理尝试反驳同一声明，≥2 个反驳则杀死
- **优雅降级**：验证全部失败时返回原始声明而非空结果；合成失败时返回原始验证数据

**8. 确定性规则（sandbox 限制）：**

Workflow 脚本运行在 Node `vm` 沙箱中，以下不可用：
- `Date.now()`、`new Date()` — 时间不确定性
- `Math.random()` — 随机性
- `require`、`import`、`fs`、网络 API — 无副作用
- `meta` 内不支持 spreads、computed keys、template interpolation、function calls

这确保 `meta` 可被静态解析，脚本运行可复现。

**与其他编排方式的对比：**

| | Subagents | Skills | Agent Teams | Workflows |
|---|---|---|---|---|
| 编排者 | Claude (逐轮) | Claude (逐轮) | Lead agent (逐轮) | **脚本** |
| 中间结果 | Claude 的 context | Claude 的 context | 共享 task list | 脚本变量 |
| 可重复性 | 工作者定义 | 指令 | 团队定义 | **编排本身** |
| 规模 | 几个委派任务 | 同 subagents | 几个长期 peer | **数十到数百个代理** |
| 中断恢复 | 重启 turn | 重启 turn | 队友继续运行 | **同一会话内可恢复** |

#### 5. Agent Teams：多实例对等协作

Agent Teams 是 Claude Code 的多代理协作模式——多个独立的 Claude Code 实例作为一个团队协同工作。

**架构组件：**

| 组件 | 角色 |
|---|---|
| Team Lead | 主 Claude Code 会话，创建团队、生成队友、协调工作 |
| Teammates | 独立的 Claude Code 实例，各自处理分配的任务 |
| Task List | 共享任务列表，队友认领和完成 |
| Mailbox | 代理间通信系统 |

**与 Subagents 的核心区别：**
- Subagents 只能向主代理汇报结果
- Agent Teams 的队友之间可以直接通信
- 用户可以直接与任何队友交互（不通过 lead）
- 每个队友是完整的独立 Claude Code 会话

**显示模式：**
- In-process：所有队友在主终端内运行（Shift+Down 切换）
- Split panes：每个队友独立面板（需 tmux 或 iTerm2）

**启用：** `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`（默认关闭）

**局限：**
- 一次只能管理一个团队
- 不支持嵌套团队
- Lead 固定，不能转让
- Split panes 不支持 VS Code 终端、Windows Terminal

#### 6. 各机制的协作关系

```
用户请求
  │
  ├── 简单任务 → Claude 直接处理（Skills 辅助）
  │
  ├── 需要委派 → Subagent（Explore / Plan / 自定义）
  │     └── 后台运行 (background)
  │
  ├── 大规模任务 → Dynamic Workflow（JS 脚本编排数十-数百代理）
  │     └── 代理并行（≤16 并发）
  │
  ├── 持续任务 → /loop（固定或动态间隔）
  │     └── ScheduleWakeup / CronCreate
  │
  └── 并行协作 → Agent Teams（3-5 队友 + Lead）
        └── 共享 Task List + Mailbox

全生命周期 → Hooks 注入行为
  ├── SessionStart → 上下文注入
  ├── PreToolUse → 拦截不安全操作
  ├── PostToolUse → 自动格式化
  ├── Stop → 质量门禁
  └── SubagentStart/Stop → 代理生命周期管理
```

---

### 学习曲线

Claude Code 的动态工作流体系存在明显的分层学习曲线：

| 层级 | 机制 | 上手难度 | 前置知识 | 达到高效所需时间 |
|---|---|---|---|---|
| L1 | Hooks（基础） | 低 | JSON 配置、shell 脚本 | 1-2 小时 |
| L2 | Skills | 低 | Markdown、prompt 工程 | 2-4 小时 |
| L3 | /loop & Cron | 低-中 | Cron 表达式 | 1-2 小时 |
| L4 | 自定义 Subagents | 中 | YAML frontmatter、工具 API | 半天 |
| L5 | Dynamic Workflows | 中-高 | JavaScript、异步编程 | 1-2 天 |
| L6 | Agent Teams | 高 | 分布式系统思维、任务分解 | 2-3 天 |

**最大的认知跳跃：** 从 L3 到 L4（从"使用内置机制"到"定义自定义代理"），以及从 L4 到 L5（从"Claude 编排"到"脚本编排"）。

**前置知识依赖图：**
```
JSON 配置 + Shell 脚本 → Hooks → Subagent Hooks
                                    ↓
Markdown + Prompt 工程 → Skills → Subagent Skills
                                    ↓
Cron + 异步思维 → /loop → ScheduleWakeup
                                    ↓
JavaScript + 异步编程 → Dynamic Workflows
                                    ↓
分布式系统 + 任务分解 → Agent Teams
```

---

### 生态与社区

**官方生态：**

| 组件 | 说明 |
|---|---|
| `/deep-research` | 内置 workflow，多源调研并交叉验证 |
| `loop.md` | 自定义默认循环 prompt |
| `.claude/workflows/` | 项目级可复用 workflow 脚本 |
| Routines | Anthropic 云端调度，独立于会话 |
| Desktop Scheduled Tasks | 本地持久化调度 |
| Agent SDK | 编程接口，CI/CD 集成 |
| Channels | 外部事件推送（CI → 会话） |

**社区生态：**

| 项目/资源 | 说明 |
|---|---|
| [claude-code-hooks-mastery](https://github.com/disler/claude-code-hooks-mastery) | Hook 模式大全，复杂工作流和代理链 |
| [agent-almanac](https://www.jsdelivr.com/package/npm/agent-almanac) | Self-Continuation Loops Playbook |
| [harness-loop](https://www.claudepluginhub.com/skills/chachamaru127-claude-code-harness/harness-loop) | /loop + ScheduleWakeup 编排插件 |
| [claude-quotas](https://mcpmarket.com/server/claude-quotas) | 配额监控 MCP，自动 sleep through reset |
| [claude-mem](https://docs.claude-mem.ai/hooks-architecture) | Hooks 架构文档和模式 |

**标志性案例：** Jarred Sumner 使用 Claude Code Dynamic Workflows 在 6 天内将 75 万行 Bun 代码从 Zig 移植到 Rust。

---

## Project Adaptation Assessment

### 当前 Forge 项目对 Claude Code 动态工作流的使用现状

| Claude Code 机制 | Forge 使用状态 | 具体用法 |
|---|---|---|
| **Hooks** | 已使用 | SessionStart（上下文注入）、SessionEnd（cleanup）、Stop（quality-gate + feature-complete）、SubagentStart/Stop |
| **Skills/Commands** | 已使用 | 16+ commands、20+ skills、项目级 skills |
| **Agents** | 已使用 | `forge:task-executor` 自定义代理 |
| **Agent Memory** | 已使用 | `.claude/agent-memory/forge-task-executor/` |
| **Plugin System** | 已使用 | 完整 Forge 插件（v3.0.0-rc.44） |
| **Worktrees** | 已使用 | `.forge/worktrees/` + `worktree.baseRef` 配置 |
| `/loop` | **未使用** | `/run-tasks` 手动实现了自己的循环 |
| ScheduleWakeup | **未使用** | run-tasks 是单会话阻塞循环 |
| CronCreate/Delete | **未使用** | 无定时健康检查或周期性任务 |
| Dynamic Workflows | **未使用** | 无 JS 脚本编排 |
| Agent Teams | **未使用** | 多代理协作走 Subagent 通道 |
| PreToolUse/PostToolUse | **未使用** | 无工具级拦截 |

### Current Stack Impact

| Area | Impact | Details |
|---|---|---|
| `/run-tasks` 循环模式 | **中** | 当前是手动 claim-dispatch-verify-continue 阻塞循环，可用 `/loop` 动态模式或 ScheduleWakeup 替代，提升中断恢复能力 |
| 定期维护 | **低-中** | `forge cleanup` 仅在 SessionEnd 运行，可用 CronCreate 增加定期健康检查 |
| 代理编排规模 | **低** | 当前 Subagent 模式足够，Dynamic Workflows 仅在需要 10+ 并发代理时有价值 |
| Agent Teams | **低** | 当前 topological task 依赖图通过 Forge CLI 管理，Agent Teams 会引入额外复杂度 |
| PreToolUse Hooks | **中** | 可用于强制执行编码规范或阻止不安全操作 |

### 潜在借鉴方向

1. **用 `/loop` 动态模式替代 `/run-tasks` 的阻塞循环**：利用 ScheduleWakeup 的自步调度能力，让代理根据任务复杂度自主调整节奏，同时获得中断恢复能力
2. **引入 PreToolUse hooks 强制编码规范**：在工具执行前拦截不合规操作，替代当前的 post-hoc 质量检查
3. **利用 CronCreate 增加定期维护**：spec drift 检测、stale task 清理等可从"只在会话结束时运行"提升为"定期自动检查"
4. **评估 Dynamic Workflows 用于大规模编排**：如全代码库审计、批量迁移等场景，可利用 JS 脚本的确定性编排替代 Claude 的 turn-by-turn 决策

---

## Risks & Caveats

| Risk | Severity | Mitigation |
|---|---|---|
| `/loop` 动态模式的 ScheduleWakeup 会重发完整 slash command | 中 | GitHub Issue #54086 已报告，需注意昂贵的用户操作不要放在 /loop prompt 中 |
| Dynamic Workflows token 消耗显著高于单会话 | 中 | 先在小切片上测试，用 `/workflows` 视图监控 token 使用 |
| Agent Teams 仍为实验性功能 | 中 | 不适合生产关键路径，当前 limitation 较多 |
| 7 天自动过期限制长周期任务 | 低 | 需要更长周期时使用 Routines 或 Desktop Scheduled Tasks |
| Subagent 不继承父会话的对话历史 | 低 | Fork 模式可解决，但有额外复杂度 |
| Hook 脚本跨平台兼容性 | 低 | Windows 需 PowerShell 脚本 + `shell: powershell` |

---

## Recommendation

Claude Code 的动态工作流体系是一个精心设计的分层编排系统。对于 Forge 项目，**不建议全面迁移到新机制**，而是**选择性引入**：

1. **短期（立即可做）**：引入 `PreToolUse` hooks 强制编码规范；将 `forge cleanup` 也可通过 CronCreate 定期运行
2. **中期（需评估）**：评估 `/loop` 动态模式是否比 `/run-tasks` 的手动循环更适合长时间任务执行
3. **长期（观察）**：Agent Teams 和 Dynamic Workflows 当前过于重量级，待 API 稳定后再评估

**Confidence level:** 高 — 基于官方文档、多个社区源和 Forge 代码库的交叉验证。

---

## Sources

| Source | URL | Used for |
|---|---|---|
| Claude Code 官方文档 — Dynamic Workflows | https://code.claude.com/docs/en/workflows | Workflows 架构、运行时约束、对比表 |
| Claude Code 官方文档 — Scheduled Tasks | https://code.claude.com/docs/en/scheduled-tasks | /loop 三种模式、ScheduleWakeup、Cron 工具 |
| Claude Code 官方文档 — Subagents | https://code.claude.com/docs/en/subagents | 代理类型、frontmatter、Fork 模式 |
| Claude Code 官方文档 — Agent Teams | https://code.claude.com/docs/en/agent-teams | 团队架构、与 subagent 对比、局限 |
| Claude Code 官方文档 — Hooks | https://code.claude.com/docs/en/hooks | 完整 Hook 事件列表 |
| Medium — Dynamic Workflows 实战指南 | https://medium.com/nginity/dynamic-workflows-in-claude-code-a-practical-setup-and-use-case-guide-0ab54304ab6f | 实际使用体验和踩坑 |
| MindStudio — Workflows vs /goal vs Agent Teams 决策框架 | https://www.mindstudio.ai/blog/claude-code-dynamic-workflows-vs-goal-vs-agent-teams-decision-framework/ | 选择决策框架 |
| Webcoda — 750K 行移植案例 | https://ai-checker.webcoda.com.au/articles/claude-code-dynamic-workflows-parallel-agents-2026 | 标志性使用案例 |
| GitHub — claude-code-hooks-mastery | https://github.com/disler/claude-code-hooks-mastery | Hook 模式和代理链 |
| Reddit — ScheduleWakeup 讨论 | https://www.reddit.com/r/ClaudeAI/comments/1sirzyx/schedulewakeup_loop_dynamic_mode_whats_new_in_cc/ | ScheduleWakeup 引入背景 |
| GitHub Issue #54086 — ScheduleWakeup 重发问题 | https://github.com/anthropics/claude-code/issues/54086 | 已知限制 |
| Blog — Every Claude Code Feature (Shrivu) | https://blog.sshh.io/p/how-i-use-every-claude-code-feature | Hooks 实战经验 |
| azukiazusa.dev — Trying Dynamic Workflow | https://azukiazusa.dev/en/blog/claude-code-dynamic-workflow | `/deep-research` 完整脚本源码和 API 分析 |
| Medium — Build Deterministic Agent Runs (Rezvani) | https://alirezarezvani.medium.com/claude-code-workflows-build-deterministic-agent-runs-eaf2c6ac52d5 | 动态 vs 确定性 workflow 设计哲学 |
| pi.dev — pi-dynamic-workflows | https://pi.dev/packages/pi-dynamic-workflows | 开源 workflow runtime 实现，完整 API 文档 |
| Agent SDK TypeScript Reference | https://code.claude.com/docs/en/agent-sdk/typescript | 官方 `agent()`、`parallel()`、`pipeline()` API 规范 |
