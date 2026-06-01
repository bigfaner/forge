# 初始化指南

> 最后更新：2026-05-30 | 版本：v3.0.0

本文档说明如何从零开始初始化一个 Forge 项目，包括 `forge init` 的完整流程、配置文件全字段参考和 Surface 检测机制。

---

## 目录

- [快速开始](#快速开始)
- [forge init 完整流程](#forge-init-完整流程)
- [config.yaml 全字段参考](#configyaml-全字段参考)
- [Surface 检测机制](#surface-检测机制)
- [端到端示例](#端到端示例)
- [常见问题](#常见问题)

---

## 快速开始

```bash
# 1. 安装 forge CLI
curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash

# 2. 安装 forge Plugin（CLI + Plugin 一步到位）
forge upgrade

# 3. 在项目中初始化
cd my-project && forge init

# 4. 验证
forge --version
forge surfaces   # 查看已检测的项目 surface 类型
```

---

## forge init 完整流程

`forge init` 是一站式的项目环境初始化命令，执行以下六个步骤：

```
forge init
  |
  ├── Step 1: 创建 .forge/ 目录
  ├── Step 2: 更新 .gitignore
  ├── Step 3: 检查 just 安装
  ├── Step 4: 交互式配置 config.yaml
  ├── Step 5: Surface 检测与确认
  └── Step 6: 输出摘要报告
```

### Step 1: 创建 .forge/ 目录

在项目根目录创建 `.forge/` 目录。如果已存在则跳过。

```
CREATED    .forge/
```

### Step 2: 更新 .gitignore

将以下运行时条目追加到 `.gitignore`（已存在的条目不会重复添加）：

```gitignore
# Forge
.forge/state.json
.forge/test-state.json
.forge/worktrees/
docs/features/*/tasks/process/
docs/features/*/tasks/index.json.lock
docs/features/*/testing/results/
tests/results/
```

### Step 3: 检查 just 安装

确保 `just`（命令运行器）已安装。如果未安装，`forge init` 会尝试自动安装。安装失败不会阻塞初始化流程，但会输出 WARNING。

```bash
# 跳过 just 安装检查
forge init --skip-just
```

### Step 4: 交互式配置 config.yaml

在交互式终端中，`forge init` 会逐项询问自动行为配置。如果 `.forge/config.yaml` 已存在，会询问是否重新配置。

配置项按 Quick/Full 两种模式分别设置，覆盖以下功能：

| 配置项 | Quick 模式默认值 | Full 模式默认值 |
|--------|-----------------|----------------|
| 自动运行高级测试 | `false` | `true` |
| 自动合并规范 | `true` | `true` |
| 自动代码清理 | `false` | `false` |
| 自动验证检查 | `false` | `false` |
| 自动运行任务 | `true` | `false` |
| 自动保存知识 | `true` | `false` |
| 自动评估 Proposal | `true` | — |
| 自动评估 PRD | `false` | — |
| 自动评估 UI 设计 | `true` | — |
| 自动评估技术设计 | `false` | — |
| 自动 git push | `false` | — |

此外，可选择配置 worktree 源分支和需要复制到 worktree 的文件。

### Step 5: Surface 检测与确认

自动扫描项目目录结构和依赖文件，检测项目的 surface 类型，并通过 TUI 确认后写入配置。详见 [Surface 检测机制](#surface-检测机制)。

### Step 6: 输出摘要报告

所有步骤完成后，输出格式化的摘要报告：

```
>>>
CREATED    .forge/
APPENDED   .gitignore (7 entries)
INSTALLED  just installation (just 1.40.0 already available)
CREATED    .forge/config.yaml (interactive)
CREATED    surfaces (cli)
<<<
```

### 命令行参数

| 参数 | 说明 |
|------|------|
| `--project-root <path>` | 指定项目根目录（默认自动检测） |
| `--skip-just` | 跳过 just 安装检查 |

---

## config.yaml 全字段参考

`.forge/config.yaml` 是 Forge 项目的核心配置文件，位于 `.forge/config.yaml`。

### 完整示例

```yaml
version: '1'
project-type: fullstack
auto:
  test:
    quick: false
    full: true
  consolidateSpecs:
    quick: true
    full: true
  cleanCode:
    quick: false
    full: true
  validation:
    quick: false
    full: true
  runTasks:
    quick: true
    full: false
  gitPush: true
  knowledgeSave:
    quick: true
    full: true
  eval:
    proposal: true
    prd: false
    uiDesign: true
    techDesign: false
worktree:
  source-branch: main
  copy-files:
    - .env
    - .env.local
coverage:
  coding.feature: 80
  coding.enhancement: 60
  coding.fix: 60
  coding.refactor: maintain
  coding.cleanup: maintain
test-framework: jest
languages:
  - typescript
  - go
surfaces: cli
execution-order:
  - frontend
  - backend
```

### 字段表格

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `version` | string | `"1"` | 配置文件格式版本。自动填充，通常无需修改 |
| `project-type` | string | — | 项目类型。可选值：`fullstack`、`mobile`、`library`、`mixed` |
| `auto.test.quick` | bool | `false` | Quick 模式下是否自动运行 surface 级高级测试 |
| `auto.test.full` | bool | `true` | Full 模式下是否自动运行 surface 级高级测试 |
| `auto.consolidateSpecs.quick` | bool | `true` | Quick 模式下是否自动提取并合并规范 |
| `auto.consolidateSpecs.full` | bool | `true` | Full 模式下是否自动提取并合并规范 |
| `auto.cleanCode.quick` | bool | `false` | Quick 模式下是否自动精炼代码 |
| `auto.cleanCode.full` | bool | `false` | Full 模式下是否自动精炼代码 |
| `auto.validation.quick` | bool | `false` | Quick 模式下是否自动运行验证检查 |
| `auto.validation.full` | bool | `false` | Full 模式下是否自动运行验证检查 |
| `auto.runTasks.quick` | bool | `true` | Quick 模式下是否自动认领并执行任务 |
| `auto.runTasks.full` | bool | `false` | Full 模式下是否自动认领并执行任务 |
| `auto.gitPush` | bool | `false` | 所有任务完成后是否自动推送到远程仓库 |
| `auto.knowledgeSave.quick` | bool | `true` | Quick 模式下是否自动保存知识 |
| `auto.knowledgeSave.full` | bool | `false` | Full 模式下是否自动保存知识 |
| `auto.eval.proposal` | bool | `true` | 生成 Proposal 后是否自动评估 |
| `auto.eval.prd` | bool | `false` | 生成 PRD 后是否自动评估 |
| `auto.eval.uiDesign` | bool | `true` | 生成 UI 设计后是否自动评估 |
| `auto.eval.techDesign` | bool | `false` | 生成技术设计后是否自动评估 |
| `worktree.source-branch` | string | — | 创建 worktree 时使用的源分支（如 `main`、`develop`） |
| `worktree.copy-files` | string[] | `[]` | 创建 worktree 时从源分支复制的文件列表（如 `.env`） |
| `surfaces` | string / map | — | 项目的 surface 类型。标量形式（如 `surfaces: cli`）用于单一类型；映射形式（如 `surfaces: {frontend: web, backend: api}`）用于多类型项目 |
| `execution-order` | string[] | — | Surface 的执行顺序。当多个 surface 类型相同时必须指定 |
| `test-framework` | string | — | 测试框架名称（如 `jest`、`pytest`），由 test-guide skill 检测生成 |
| `languages` | string[] | `[]` | 项目使用的编程语言列表 |
| `coverage` | map | 内置默认值 | 按任务类型配置覆盖率策略。值为整数（百分比）或 `maintain`（保持现有覆盖率） |

### 覆盖率默认值

| 任务类型 | 默认策略 |
|---------|---------|
| `coding.feature` | 80% |
| `coding.enhancement` | 60% |
| `coding.fix` | 60% |
| `coding.refactor` | maintain（保持现有覆盖率） |
| `coding.cleanup` | maintain（保持现有覆盖率） |

### Surface 值

`surfaces` 字段支持以下类型：

| Surface | 说明 |
|---------|------|
| `cli` | 命令行应用 |
| `tui` | 终端用户界面应用 |
| `api` | HTTP API 服务 |
| `web` | Web 前端应用 |
| `mobile` | 移动端应用 |

### 通过命令行管理配置

除了 `forge init` 的交互式配置，还可以使用 `forge config` 命令读写单个配置项：

```bash
# 查看完整配置
forge config

# 读取单个字段
forge config surfaces

# 设置单个字段
forge config set surfaces cli
forge config set auto.gitPush true
forge config set worktree.source-branch main
```

---

## Surface 检测机制

Surface 检测是 Forge v3.0.0 的核心能力，用于自动识别项目的测试 surface 类型（cli/tui/api/web/mobile）。检测结果直接影响测试生成和执行策略。

### 检测原理

`forge surfaces detect` 扫描项目目录结构和依赖文件，通过两级策略识别 surface 类型：

**第一级：依赖信号检测**

扫描项目的依赖管理文件，匹配已知框架/库：

| 依赖文件 | 检测的信号 |
|---------|-----------|
| `package.json` | react、vue、express、commander、ink 等 |
| `go.mod` | gin、cobra、bubbletea、tview 等 |
| `Cargo.toml` | actix-web、clap、ratatui 等 |
| `pyproject.toml` / `setup.py` | flask、click、django 等 |
| `AndroidManifest.xml` / `*.xcodeproj` / `pubspec.yaml` | Android / iOS / Flutter |

**第二级：结构推断**

当依赖信号为空时，通过目录结构推断 surface 类型：

| 生态系统 | 规则 | 推断结果 |
|---------|------|---------|
| Go | `cmd/` 下有子目录 | `cli` (inference:cmd-dir) |
| Go | `api/` 目录存在 | `api` (inference:api-dir) |
| Go | `handler/` 目录存在 | `api` (inference:handler-dir) |
| Node.js | `package.json` 中有 `bin` 字段 | `cli` (inference:bin-field) |
| Node.js | 根目录有 `index.html` | `web` (inference:index-html) |
| Python | `pyproject.toml` 中有 `[project.scripts]` | `cli` (inference:py-scripts) |
| Python | 根目录有 `app.py` 或 `main.py`（非库项目） | `cli` (inference:py-main) |

### 信号冲突解决

当同一目录检测到多个 surface 信号时，按以下优先级选择（数值越小优先级越高）：

| 优先级 | Surface |
|--------|---------|
| 1（最高） | web |
| 2 | mobile |
| 3 | api |
| 4 | cli |
| 5（最低） | tui |

### Workspace 模式

对于使用 workspace/monorepo 结构的项目（检测 `pnpm-workspace.yaml` 或 `package.json` 中的 `workspaces` 字段），检测模式有所不同：

- **非 workspace**：扫描根目录和子目录，如果只检测到一个 surface 类型，使用标量形式（`surfaces: cli`）
- **Workspace**：跳过根目录依赖，只扫描子目录，每个子目录独立检测，使用映射形式（`surfaces: {frontend: web, backend: api}`）

### forge surfaces detect

```bash
# 只读模式：检测结果输出到终端，不修改配置
forge surfaces detect

# 输出示例：
# cli (detected:cobra)
# 或
# backend=api (detected:gin)
# frontend=web (detected:react)
```

### forge surfaces detect --apply

```bash
# 交互模式：检测结果通过 TUI 确认后写入 .forge/config.yaml
forge surfaces detect --apply
```

`--apply` 模式的流程：

1. 运行检测，展示结果和来源标注
2. TUI 交互确认（支持确认、重新检测、手动编辑）
3. 确认后写入 `.forge/config.yaml` 的 `surfaces` 字段

如果配置中已有 surface 设置，TUI 会提供三个选项：

| 选项 | 行为 |
|------|------|
| **Confirm** | 保留现有配置 |
| **Re-detect** | 重新运行完整检测 |
| **Edit** | 手动输入 surface 配置 |

### 检测深度控制

通过环境变量 `FORGE_DETECT_DEPTH` 控制目录扫描深度：

```bash
# 默认深度：3 层
# 自定义深度（1-10）
FORGE_DETECT_DEPTH=5 forge surfaces detect
```

注意：`FORGE_DETECT_DEPTH=0` 或负值会产生错误。

---

## 端到端示例

以下是一个完整的从零开始设置 Forge 项目的流程，以一个 Go CLI 项目为例。

### 1. 安装 Forge CLI

```bash
# 安装 forge CLI binary
curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash

# 刷新终端环境
source ~/.zshrc    # zsh 用户
source ~/.bashrc   # bash 用户

# 安装 forge Plugin
forge upgrade

# 验证安装
forge --version
```

### 2. 初始化项目

```bash
cd /path/to/your/project
forge init
```

交互式配置过程：

```
? Quick mode: auto-run advanced tests? No
? Full mode: auto-run advanced tests? Yes
? Quick mode: auto-consolidate specs? Yes
? Full mode: auto-consolidate specs? Yes
? Quick mode: auto code cleanup? No
? Full mode: auto code cleanup? No
? Quick mode: auto validation? No
? Full mode: auto validation? No
? Quick mode: auto-run tasks? Yes
? Full mode: auto-run tasks? No
? Quick mode: auto knowledge save? Yes
? Full mode: auto knowledge save? No
? Auto-evaluate proposals? Yes
? Auto-evaluate PRD documents? No
? Auto-evaluate UI designs? Yes
? Auto-evaluate tech designs? No
? Auto git push after all tasks complete? No
? Worktree source branch (leave empty to skip) main
? Files to copy into worktrees [ ] (select with space)
```

Surface 检测结果：

```
Detected surfaces:
  cli (detected:cobra)
? Confirm detected surfaces? Yes
```

摘要报告：

```
>>>
CREATED    .forge/
APPENDED   .gitignore (7 entries)
SKIPPED    just installation (just 1.40.0 already available)
CREATED    .forge/config.yaml (interactive)
CREATED    surfaces (cli)
<<<
```

### 3. 验证配置

```bash
# 查看完整配置
forge config

# 查看 surface 类型
forge surfaces
# 输出：cli

# 查看所有 surface 类型
forge surfaces --types
# 输出：cli
```

### 4. 开始使用

```bash
# 快速模式：直接启动功能开发
/quick

# 完整模式：逐步走完流程
/brainstorm -> /write-prd -> /tech-design -> /breakdown-tasks -> /run-tasks
```

### 5. 查看任务状态

```bash
# 设置当前 feature
forge feature my-feature

# 查看任务列表
forge task list

# 认领并执行下一个任务
forge task claim
```

---

## 常见问题

### forge init 在非交互终端中的行为

在非交互终端（如 CI 环境或管道）中，`forge init` 会跳过交互式配置和 surface 检测：

- Step 4（config.yaml）：显示 "SKIPPED (non-interactive terminal)"
- Step 5（surfaces）：显示 "SKIPPED (non-interactive terminal)"

此时需要手动创建 `.forge/config.yaml`：

```bash
mkdir -p .forge
cat > .forge/config.yaml << 'EOF'
version: '1'
surfaces: cli
EOF
```

### 如何重新配置

直接再次运行 `forge init`。如果 `.forge/config.yaml` 已存在，会询问是否重新配置：

```bash
? Config already exists. Reconfigure? Yes
```

如果 surfaces 已配置，TUI 会提供 Confirm / Re-detect / Edit 三个选项。

### 如何单独更新 Surface 配置

无需重新运行完整 `forge init`，直接使用 `forge surfaces detect --apply`：

```bash
# 重新检测并确认写入
forge surfaces detect --apply
```

或通过 `forge config set` 直接修改：

```bash
forge config set surfaces api
```

### Workspace/Monorepo 项目的 Surface 配置

对于 monorepo 项目，`surfaces` 使用映射形式：

```yaml
surfaces:
  frontend: web
  backend: api
```

如果多个 surface 类型相同（如两个 `api`），必须指定 `execution-order`：

```yaml
surfaces:
  user-service: api
  order-service: api
execution-order:
  - user-service
  - order-service
```

### init 失败后如何处理

| 错误 | 解决方案 |
|------|---------|
| `.forge/` 创建失败 | 检查目录写权限 |
| `.gitignore` 更新失败 | 检查文件权限 |
| just 安装失败 | 手动安装：`brew install just`（macOS）或参考 [just 官方文档](https://github.com/casey/just) |
| Surface 检测无结果 | 确认项目有依赖文件（`go.mod`、`package.json` 等）或手动设置：`forge config set surfaces cli` |
| 配置验证失败 | 检查 surface key 格式（仅允许小写字母、数字、连字符）和 execution-order 引用 |
