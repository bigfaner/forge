---
title: "Forge Plugin Distribution Model"
domains: [plugin, distribution, path-resolution, skills, hooks, references, CLAUDE_PLUGIN_ROOT]
---

# Forge Plugin 架构与分发机制

Forge 是一个分发到用户环境的 Claude Code plugin，不是只在开发源码中运行的工具。理解分发模型对编写 skill 至关重要。

## 1. 分发模型

### 安装位置

```
~/.claude/plugins/cache/forge/forge/<version>/
```

例如 Windows: `C:\Users\<user>\.claude\plugins\cache\forge\forge\3.0.0-rc.23\`

### 分发包内容（源码 → 安装映射）

源码仓库 `plugins/forge/` 下的以下目录随 plugin 包分发到用户环境：

```
~/.claude/plugins/cache/forge/forge/<version>/
├── .claude-plugin/plugin.json     # plugin 元数据（name, version, description）
├── agents/                        # 分发 — subagent 定义
│   └── task-executor.md
├── commands/                      # 分发 — 斜杠命令入口
├── hooks/                         # 分发 — 生命周期钩子
│   ├── hooks.json                 #   hook 注册表（使用 ${CLAUDE_PLUGIN_ROOT}）
│   ├── session-start              #   注入 guide.md 到上下文
│   └── guide.md                   #   forge 规范指南（随 SessionStart 加载）
└── skills/                        # 分发 — skill 定义 + templates
    └── <skill-name>/
        ├── SKILL.md
        ├── templates/
        └── experts/               # eval skill 专用：协议 + 专家角色
            ├── protocol/
            │   ├── scorer-protocol.md
            │   └── reviser-protocol.md
            ├── freeform/          # freeform 评估专用：动态专家生成
            │   ├── expert-inference.md
            │   ├── expert-template.md
            │   ├── extraction-prompt.md
            │   ├── freeform-reviewer.md
            │   └── freeform-review-protocol.md
            └── scorer/
                ├── architect.md
                ├── code-reviewer.md
                ├── cto.md
                ├── editor.md
                ├── pm.md
                ├── qa.md
                ├── ux-auditor.md
                └── ux-engineer.md
```

**不分发的内容**：源码仓库中的 `docs/`、`.git/`、Go 源码、测试文件等。

## 2. 组件说明

| 组件 | 作用 | 分发 | 用户的交互方式 |
|------|------|------|---------------|
| **skills/** | Skill 定义文件（SKILL.md + templates + experts） | 是 | `/skill-name` 斜杠命令 |
| **commands/** | 轻量命令入口（单 .md 文件） | 是 | `/command-name` 斜杠命令 |
| **agents/** | Subagent 定义 | 是 | Skill/Command 通过 Agent 工具调用 |
| **hooks/** | 生命周期钩子 + forge 规范指南 | 是 | 自动触发（SessionStart/SubagentStart 注入 guide.md，Stop 触发 quality-gate 和 feature complete） |

## 3. 核心依赖

### `just` — 构建任务运行器

- **必须依赖**，不是可选的
- `forge init` 引导用户安装 `just`
- `just` 抽象不同语言的常用命令（compile, fmt, lint, unit-test, test, probe 等）
- Skill 中引用 `just` 是预期行为

### Forge CLI — 命令行工具

- `forge task claim/submit/status/index` 等命令是 skill 的操作接口
- Skill 中引用 `forge` CLI 命令是预期行为

### `skills/eval/experts/` — 评估协议与专家角色

- 随 plugin 分发在 `skills/eval/experts/` 目录中
- `protocol/` 包含通用评分和修订工作流（不含领域知识）
- `freeform/` 包含 freeform 评估的动态专家生成模板（领域推理、专家模板、提取提示词、评审协议）
- `scorer/` 包含各领域专家角色描述（角色 + 领域失败模式）
- Eval skill 在运行时组合 protocol + expert，通过 `general-purpose` agent 执行
- Skill 中引用 `experts/` 下的文件是预期行为（相对路径，与 `rubrics/` 同级约定）

## 4. 用户项目目录规范

以下目录由 forge skill 在**用户项目**中生成，用户必须遵循：

```
<user-project>/
├── docs/
│   ├── ARCHITECTURE.md       # 系统架构（用户维护）
│   ├── business-rules/       # 跨 feature 业务规则（consolidate-specs 生成）
│   ├── conventions/          # 技术规范：编码标准、API 约定、命名规则
│   ├── reference/            # 系统规范：环境、部署、技术栈
│   ├── decisions/            # 技术决策（/record-decision 生成）
│   ├── lessons/              # 经验教训（/learn-lesson 生成）
│   ├── proposals/            # 改进提案（/brainstorm 生成）
│   └── sitemap/sitemap.json  # 页面元素地图（/gen-sitemap 生成）
├── docs/features/<slug>/     # Feature 工作区
│   ├── manifest.md           # Feature 入口
│   ├── tasks/                # 任务文件
│   └── testing/              # 测试脚本（运行时生成）
├── tests/                    # 项目级 E2E 测试（各 feature 测试套件）
├── forge-cli/tests/          # forge-cli 自身 E2E 测试
└── .forge/
    └── config.yaml           # Forge 项目配置（test profile 等）
```

**关键约束**：
- `docs/conventions/` 和 `docs/business-rules/` 由 `/consolidate-specs` 从 feature 文档中提取，agent 在任务执行时读取这些规范
- `tests/` 和 `forge-cli/tests/` 测试通过标签管理（`forge test promote`），不使用文件迁移
- `records/` 由 `forge task submit` 生成，不能手动写入

## 5. 路径解析机制

### Hooks/Scripts — 使用 `${CLAUDE_PLUGIN_ROOT}`

hooks.json 和 shell 脚本可以通过环境变量引用 plugin 安装位置：

```json
{
  "hooks": {
    "SessionStart": [{
      "hooks": [{
        "type": "command",
        "command": "\"${CLAUDE_PLUGIN_ROOT}/hooks/run-hook.cmd\" session-start"
      }]
    }]
  }
}
```

### Skills/Commands — 使用相对路径

Skill 文件（SKILL.md）和 Command 文件（.md）中的路径使用相对于当前文件所在目录的相对路径。Claude 知道 SKILL.md 的文件位置，能够正确解析相对路径引用。

**路径规则：**

| 引用目标 | 路径风格 | 示例 |
|---------|---------|------|
| skill/command 内部文件 | 相对路径 | `rules/platform-routing.md`、`templates/decision-entry.md`、`rubrics/<type>.md` |
| 跨 skill 文件 | 描述性路径 + 上下文 | `ui-design/templates/styles/<name>.md`（注明 resolve relative to the skills parent directory） |
| 用户项目文件 | 项目相对路径 | `docs/decisions/<type>.md` |
| forge CLI | 命令名 | `forge task claim` |

**模板文件（templates/）**：不是 SKILL.md，不经过 skill 内容加载机制。模板中的跨 skill 引用路径由读取该模板的 SKILL.md 在执行时负责解析。如果模板需要引用其他 skill 的资源（如 eval rubrics），SKILL.md 应在指令中提供正确的路径，而非在模板中硬编码。

**禁止：**
- 项目根路径（`plugins/forge/...`）— 分发后路径不存在
- `${CLAUDE_SKILL_DIR}` 变量 — 使用相对路径替代，保持 SKILL.md 内容自描述

## 6. Pipeline

### Full Pipeline

```
/brainstorm → /write-prd → /ui-design? → /tech-design → /breakdown-tasks → /run-tasks → /submit-task
                  ↓              ↓              ↓
              /eval-prd     /eval-ui      /eval-design
```

### Quick Pipeline（1-15 个 coding 任务）

```
/brainstorm → /quick-tasks → /run-tasks → /submit-task
```

### 测试 Pipeline

```
/gen-journeys → /eval-journey → /gen-contracts → /eval-contract → /gen-test-scripts → /run-tests → forge test promote
```

### 辅助 Skill（任意阶段可用）

- `/consolidate-specs` — 提取业务规则和技术规范
- `/learn` — 统一知识积累（吸收了 /learn-lesson 和 /record-decision）
- `/forensic` — 分析历史会话偏差原因
- `/init-justfile` — 初始化项目 Justfile
- `/eval` — 通用文档评估（17 种 rubric 类型）
- `/clean-code` — 代码简化与质量提升
- `/deep-research` — 技术或产品深度调研
- `/extract-design-md` — 从应用提取视觉样式生成 DESIGN.md
- `/gen-sitemap` — 自动生成 Web 应用站点地图
- `/test-guide` — 测试指南（profile-based）
