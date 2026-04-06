# Zcode Plugin & Marketplace 设计

## 概述

将 `claude-code-go/.claude` 中的 skills、agents、commands 抽取为 Claude Code 插件，通过插件市场分发。

## 目标

- **目标用户**：开源社区
- **插件范围**：任务管理系统 + 流程辅助
- **依赖**：claude-task-cli（通过 `/init-zcode` 命令自动编译安装）

## 插件内容

### Skills

| Skill | 用途 |
|-------|------|
| claim-task | 认领任务 |
| record-task | 记录任务 |
| set-task-status | 设置任务状态 |
| breakdown-tasks | 分解任务 |
| write-prd | 编写 PRD |
| design-tech | 技术设计 |
| learn-lesson | 记录经验教训 |
| git-commit | Git 提交规范 |

### Commands

| Command | 用途 |
|---------|------|
| /init-zcode | 自动编译安装 claude-task-cli |
| /simplify-skill | 重构 skill 文件，拆分模板/示例到独立文件 |
| /execute-task | 执行单个任务（TDD 工作流） |
| /run-tasks | 自动分发任务到 subagents 执行 |

## 目录结构

```
zcode/
├── .claude-plugin/
│   └── marketplace.json           # Marketplace 定义
├── plugins/
│   └── zcode/                     # 插件目录
│       ├── .claude-plugin/
│       │   └── plugin.json        # 插件清单
│       ├── skills/
│       │   ├── claim-task/
│       │   │   └── SKILL.md
│       │   ├── record-task/
│       │   │   └── SKILL.md
│       │   ├── set-task-status/
│       │   │   └── SKILL.md
│       │   ├── breakdown-tasks/
│       │   │   └── SKILL.md
│       │   ├── write-prd/
│       │   │   └── SKILL.md
│       │   ├── design-tech/
│       │   │   └── SKILL.md
│       │   ├── learn-lesson/
│       │   │   └── SKILL.md
│       │   └── git-commit/
│       │       └── SKILL.md
│       └── commands/
│           ├── init-zcode.md      # 自动安装命令
│           ├── simplify-skill.md  # 重构 skill
│           ├── execute-task.md    # 执行单个任务
│           └── run-tasks.md       # 批量执行任务
└── ...其他文件
```

## 核心组件设计

### 1. marketplace.json

```json
{
  "name": "zcode-marketplace",
  "owner": {
    "name": "zcode"
  },
  "metadata": {
    "description": "Claude Code productivity tools for task management and workflow"
  },
  "plugins": [
    {
      "name": "zcode",
      "source": "./plugins/zcode",
      "description": "Task management and workflow helper tools",
      "version": "1.0.0"
    }
  ]
}
```

### 2. plugin.json

```json
{
  "name": "zcode",
  "version": "1.0.0",
  "description": "Task management and workflow helper tools for Claude Code",
  "keywords": ["task", "workflow", "productivity", "prd", "git"]
}
```

### 3. /init-zcode 命令

**功能**：自动检测操作系统，调用 claude-task-cli 编译脚本，安装到用户 PATH。

**流程**：
1. 检测操作系统（Windows/Linux/macOS）
2. 定位 claude-task-cli 仓库路径
3. 调用对应脚本：
   - Windows: `scripts/install-local.ps1`
   - Linux/macOS: `scripts/install-local.sh`
4. 编译并安装到 `~/.local/bin` (Unix) 或用户目录 (Windows)
5. **提示用户重新打开终端**以刷新环境变量

**SKILL.md 内容**：
```markdown
---
description: 自动编译并安装 claude-task-cli 工具
---

自动编译安装 claude-task-cli：

1. 检测操作系统
2. 调用 claude-task-cli/scripts 中的安装脚本
3. 编译并安装二进制文件

安装完成后，**请重新打开终端**以使环境变量生效。
```

### 4. Skills 改造

原 skills 中的 Go 脚本调用改为使用 `task` CLI 命令：

**改造前**：
```
go run ${CLAUDE_PROJECT_ROOT}/.claude/skills/claim-task/scripts/claim.go
```

**改造后**：
```
task claim <task-id>
```

## 用户安装流程

```bash
# 1. 添加 marketplace
/plugin marketplace add <github-repo-url>

# 2. 安装插件
/plugin install zcode@zcode-marketplace

# 3. 初始化（安装依赖）
/init-zcode

# 4. 重新打开终端

# 5. 开始使用
/claim-task <task-id>
/record-task <task-id>
/execute-task           # 执行单个任务（TDD 工作流）
/run-tasks              # 批量执行任务
/simplify-skill <name>  # 重构 skill 文件
```

## 实施步骤

1. 创建目录结构和配置文件
2. 从 claude-code-go 复制并改造 skills（8个）
3. 从 claude-code-go 复制 commands（4个）
4. 创建 init-zcode 命令
5. 更新 skills/commands 中的脚本调用为 task CLI
6. 测试安装流程
7. 发布到 GitHub

## 版本策略

- 遵循语义版本 (SemVer)
- 初始版本: 1.0.0
- 通过 plugin.json 和 marketplace.json 同步版本

## 后续扩展

- 若需拆分功能，可从 zcode 派生独立插件
- 可添加更多 workflow 相关 skills
- 可添加 MCP servers 集成外部工具
