# Forge Plugin 架构与分发机制

Forge 是一个分发到用户环境的 Claude Code plugin，不是只在开发源码中运行的工具。理解分发模型对编写 skill 至关重要。

## 1. 分发模型

### 安装位置

```
~/.claude/plugins/cache/forge/forge/<version>/
```

例如 Windows: `C:\Users\<user>\.claude\plugins\cache\forge\forge\2.18.0\`

### 分发包内容（源码 → 安装映射）

源码仓库 `plugins/forge/` 下的以下目录随 plugin 包分发到用户环境：

```
~/.claude/plugins/cache/forge/forge/<version>/
├── .claude-plugin/manifest.json   # plugin 元数据（name, version, description）
├── agents/                        # 分发 — subagent 定义
│   ├── doc-scorer.md
│   ├── doc-reviser.md
│   └── task-executor.md
├── commands/                      # 分发 — 斜杠命令入口
├── hooks/                         # 分发 — 生命周期钩子
│   ├── hooks.json                 #   hook 注册表（使用 ${CLAUDE_PLUGIN_ROOT}）
│   ├── session-start              #   注入 guide.md 到上下文
│   └── guide.md                   #   forge 规范指南（随 SessionStart 加载）
├── references/                    # 分发 — 共享参考文档
│   └── shared/
│       ├── config.yaml
│       ├── decision-logging.md
│       └── sitemap.json
├── scripts/                       # 分发 — 辅助脚本
│   └── validate-index.sh
└── skills/                        # 分发 — skill 定义 + templates
    └── <skill-name>/
        ├── SKILL.md
        └── templates/
```

**不分发的内容**：源码仓库中的 `docs/`、`.git/`、Go 源码、测试文件等。

## 2. 组件说明

| 组件 | 作用 | 分发 | 用户的交互方式 |
|------|------|------|---------------|
| **skills/** | Skill 定义文件（SKILL.md + templates） | 是 | `/skill-name` 斜杠命令 |
| **commands/** | 轻量命令入口（单 .md 文件） | 是 | `/command-name` 斜杠命令 |
| **agents/** | Subagent 定义 | 是 | Skill/Command 通过 Agent 工具调用 |
| **hooks/** | 生命周期钩子 + forge 规范指南 | 是 | 自动触发（SessionStart 注入 guide.md） |
| **references/** | 共享参考文档 | 是 | Skill 内部 Read 引用 |
| **scripts/** | 辅助 shell 脚本 | 是 | hooks 或 skill 通过 Bash 调用 |

## 3. 核心依赖

### `just` — 构建任务运行器

- **必须依赖**，不是可选的
- `forge init` 引导用户安装 `just`
- `just` 抽象不同语言的常用命令（compile, fmt, lint, test, e2e-* 等）
- Skill 中引用 `just` 是预期行为

### Forge CLI — 命令行工具

- `forge task claim/submit/status/index` 等命令是 skill 的操作接口
- Skill 中引用 `forge` CLI 命令是预期行为

### `doc-scorer` / `doc-reviser` — 内置 Agent

- 随 plugin 分发在 `agents/` 目录中
- Skill 中引用这些 agent 是预期行为

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
├── tests/e2e/                # E2E 回归测试（仅通过 /graduate-tests 添加）
└── .forge/
    └── config.yaml           # Forge 项目配置（test profile 等）
```

**关键约束**：
- `docs/conventions/` 和 `docs/business-rules/` 由 `/consolidate-specs` 从 feature 文档中提取，agent 在任务执行时读取这些规范
- `tests/e2e/` 只能通过 `/graduate-tests` 添加，不能手动写入
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

### Skills — Markdown 指令，无自动路径解析

Skill 文件（SKILL.md）是 Claude AI 读取的 Markdown 指令。当 skill 中写 `Read plugins/forge/skills/eval/rubrics/harness.md` 时：

- AI 尝试从**用户项目根目录**读取 `plugins/forge/skills/eval/rubrics/harness.md`
- 实际文件在 `~/.claude/plugins/cache/forge/forge/<version>/skills/eval/rubrics/harness.md`
- **路径不匹配** — 这是当前已知的分发正确性问题

### 正确的路径引用方式

- 引用 **plugin 内部文件**（templates、rubrics、references）：需要确保路径在安装后可解析
- 引用 **用户项目文件**（docs/、features/、.forge/）：使用项目相对路径即可
- 引用 **forge CLI**：`forge <command>` 即可，CLI 已安装到 PATH

## 6. 两条 Pipeline

### Full Pipeline

```
/brainstorm → /write-prd → /ui-design? → /tech-design → /breakdown-tasks → 执行 → /submit-task
                  ↓              ↓              ↓
              /eval-prd     /eval-ui      /eval-design
```

### Quick Pipeline（1-10 个任务）

```
/brainstorm → /quick-tasks → 执行 → /submit-task
```

### 测试 Pipeline

```
/gen-test-cases → /gen-test-scripts → /run-e2e-tests → /graduate-tests
```

### 辅助 Skill（任意阶段可用）

- `/consolidate-specs` — 提取业务规则和技术规范
- `/learn-lesson` — 记录经验教训
- `/record-decision` — 记录架构决策
- `/forensic` — 分析历史会话偏差原因
- `/init-justfile` — 初始化项目 Justfile
- `/eval` — 通用文档评估（16 种 rubric 类型）
