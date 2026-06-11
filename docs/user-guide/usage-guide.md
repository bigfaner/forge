# 使用指南

> 最后更新：2026-05-30 | 对应版本：v3.0.0

本文档提供 Forge 两种工作模式的端到端实战示例、常用单命令场景以及常见问题排错指引。如果你还没有完成安装和环境配置，请先阅读 [环境准备](environment-setup.md) 和 [初始化项目](initialization.md)。

---

## 目录

- [Full Mode 端到端实战](#full-mode-端到端实战)
- [Quick Mode 端到端实战](#quick-mode-端到端实战)
- [单命令场景](#单命令场景)
- [常见问题与排错](#常见问题与排错)

---

## Full Mode 端到端实战

Full Mode 适合复杂功能开发（预计 >10 个任务、涉及架构决策或需要 PRD 验收标准）。以下是一个从零开始的完整示例。

### 场景：为 Web 应用添加用户通知系统

假设你正在开发一个 Web 应用，需要添加一套用户通知系统（站内消息 + 邮件推送）。预计涉及前后端多个模块，任务数量 >10。

#### 第 1 步：探索需求

```
/brainstorm
```

与 AI 对话，讨论通知系统的功能范围。完成后生成 `docs/features/user-notification/proposal.md`，包含功能边界、技术约束和优先级建议。

可选评估：

```
/eval-proposal
```

对 proposal 进行 1000 分制评分，未达标时自动迭代修订。

#### 第 2 步：编写 PRD

```
/write-prd
```

基于 proposal.md 生成 PRD 三件套：

- `prd/prd-spec.md` — 需求规格
- `prd/prd-user-stories.md` — 用户故事
- `prd/prd-ui-functions.md` — UI 功能（可选）

同时更新 `manifest.md` 状态为 `prd`。

可选评估：

```
/eval-prd
```

#### 第 3 步：技术设计

```
/tech-design
```

基于 PRD 生成技术设计文档 `design/tech-design.md`（和可选的 `design/api-handbook.md`）。

如果涉及 UI 交互，先执行 UI 设计：

```
/ui-design
```

生成 `ui/ui-design.md` 和可选的 HTML 原型 `ui/prototype/`。

可选评估：

```
/eval-design
/eval-ui
```

#### 第 4 步：拆分任务

```
/breakdown-tasks
```

基于设计文档自动拆分为多个任务文件，生成 `tasks/*.md` 和 `tasks/index.json`。每个任务文件包含标题、描述、验收标准和依赖关系。

#### 第 5 步：执行任务

```
/run-tasks
```

`/run-tasks` 启动自动分发循环，按拓扑排序依次将任务分发给 task-executor agent 执行。你也可以手动执行单个任务：

```
/execute-task
```

执行时提供任务 ID，agent 会按照 TDD 流程完成实现。

每个任务执行完成后，agent 会自动：
1. 运行 Quality Gate（compile -> fmt -> lint -> unit-test）
2. 调用 `/submit-task` 记录执行结果
3. 创建 git commit

#### 第 6 步：收尾

全部任务完成后，Hook 自动触发项目级 Quality Gate（FullGateSequence：compile -> fmt -> lint -> unit-test -> test -> probe）。你也可以手动运行：

```
forge quality-gate
```

最后，将 feature 文档中的规范沉淀到项目级：

```
/consolidate-specs
```

---

## Quick Mode 端到端实战

Quick Mode 适合小功能、bug 修复或配置调整（预计 1-10 个任务），跳过 PRD 和设计阶段。

### 场景：修复用户登录超时 Bug

假设用户反馈登录时偶尔出现超时错误，需要定位并修复。

#### 第 1 步：启动快速模式

```
/quick
```

Forge 会引导你进入快速模式流程。

#### 第 2 步：快速探索

```
/brainstorm
```

简要描述问题现象和修复思路。完成后生成 `proposal.md`。

#### 第 3 步：生成任务

```
/quick-tasks
```

基于 proposal.md 直接生成任务文件和 index.json。跳过 PRD 和设计阶段。

#### 第 4 步：执行任务

```
/run-tasks
```

自动分发执行。如果只有一个任务，也可以用：

```
/execute-task
```

每个任务完成后自动通过 Quality Gate、提交记录和 git commit。

#### 第 5 步：验证

全部完成后，Hook 自动触发最终 Quality Gate。如果需要手动检查：

```bash
forge task list          # 确认所有任务状态为 completed
forge quality-gate       # 运行项目级质量门禁
```

---

## 单命令场景

除了完整工作流，以下场景可以直接使用单个 Skill 完成。

### 场景 1：记录技术决策和经验

在开发过程中，遇到重要的技术决策或踩了坑，用 `/learn` 知识积累：

```
/learn
```

`/learn` 是统一知识积累入口，可以记录：
- **决策**（Decision）：为什么选择方案 A 而不是方案 B
- **经验**（Lesson）：踩过的坑和解决方案
- **惯例**（Convention）：团队约定的编码规范
- **业务规则**（Business Rule）：从需求中提取的业务逻辑

记录会写入 `docs/decisions/` 和 `docs/lessons/`，后续 agent 执行任务时会自动读取这些知识。

### 场景 2：规范沉淀与漂移检测

完成一个功能后，将 feature 文档中的业务规则和技术规范沉淀到项目级目录：

```
/consolidate-specs
```

该命令会：
1. 从 feature 文档中提取业务规则和技术规范到预览文件
2. 检测与 `docs/business-rules/` 和 `docs/conventions/` 中已有内容的重叠
3. 用户确认后集成到项目级目录
4. 自动检测并修复规范与代码的漂移（spec drift）

### 场景 3：Bug 修复

定位到一个具体 bug，直接使用 TDD 流程修复：

```
/fix-bug
```

该命令会引导你：复现问题 -> 编写失败测试 -> 修复代码 -> 验证测试通过。

### 场景 4：代码质量精炼

在不改变行为的前提下优化代码表达：

```
/clean-code
```

在限定 scope 内应用精炼原则，可选附带 Quality Gate 验证安全性。

### 场景 5：深度调研

需要对比多个技术方案或深入了解某项技术：

```
/deep-research
```

产出结构化研究报告，支持单技术深度分析和多候选方案对比两种模式。

---

## 常见问题与排错

### 1. 安装失败：forge upgrade 报错

**症状**：执行 `forge upgrade` 时报错。

**排查步骤**：

```bash
# 检查 forge CLI 是否已安装
forge version

# 检查网络连接
curl -I https://github.com

# 重新安装 CLI
curl -fsSL https://github.com/bigfaner/forge/releases/latest/download/install.sh | bash

# 刷新终端后重试
source ~/.zshrc    # zsh 用户
source ~/.bashrc   # bash 用户
forge upgrade
```

**常见原因**：
- CLI 未安装：先执行 curl 安装脚本
- 网络问题：检查代理设置
- PATH 未更新：刷新终端或重新打开

### 2. forge init 失败或配置错误

**症状**：运行 `forge init` 后 `.forge/config.yaml` 未生成，或配置有误。

**排查步骤**：

```bash
# 确认插件已安装
/plugin list

# 确认 forge 二进制在 PATH 中
which forge

# 手动初始化
forge init

# 检查配置文件
cat .forge/config.yaml
```

**常见原因**：
- 插件未安装：先执行 `forge upgrade`
- 权限问题：检查项目目录的读写权限

### 3. 工作流中断：Skill 报错"前置条件不满足"

**症状**：执行某个 Skill 时提示"缺少前置文件"或"前置条件未满足"。

**排查步骤**：

每个 Skill 执行前会检查前置文件是否存在。按以下顺序确认：

| Skill | 前置文件 |
|-------|---------|
| `/write-prd` | `proposal.md` |
| `/tech-design` | `prd/prd-spec.md` |
| `/ui-design` | `prd/prd-spec.md`（以及可选的 `prd/prd-ui-functions.md`） |
| `/breakdown-tasks` | `design/tech-design.md` |
| `/quick-tasks` | `proposal.md` |
| `/run-tasks` | `tasks/index.json` |

```bash
# 检查 feature 目录下有哪些文件
ls docs/features/<your-feature>/

# 确认 manifest.md 的状态
cat docs/features/<your-feature>/manifest.md
```

如果中间步骤的文件缺失，需要先执行对应的上游 Skill。

### 4. 任务状态 blocked：如何解除

**症状**：`forge task list` 显示某个任务状态为 `blocked`，任务无法被认领执行。

**排查步骤**：

```bash
# 查看任务详情（含依赖关系）
forge task query <task-id> -v

# 查看哪些依赖未完成
forge task list
```

**blocked 的常见原因**：

| 原因 | 解除方式 |
|------|---------|
| 依赖任务未完成 | 等待依赖任务完成，或先完成依赖任务 |
| 依赖任务被 rejected | `forge task reopen <dep-id>` 恢复依赖任务 |
| Quality Gate 中 fmt/lint 失败 | 检查工具链配置，修复格式问题 |
| 执行中出现超范围错误 | agent 已自动创建 fix-task，完成 fix-task 后源任务自动恢复为 pending |

手动解除 blocked（仅在确认问题已解决时使用）：

```bash
# 手动切换状态（需要提供 reason）
forge task transition <task-id> --status pending --reason "依赖任务已完成"
```

**注意**：`completed` 状态不可逆，无法通过任何命令恢复。`rejected` 和 `skipped` 可通过 `forge task reopen` 恢复为 `pending`。

### 5. 测试失败：Quality Gate 报错

**症状**：任务执行时 Quality Gate 在某个步骤失败。

**排查步骤**：

根据失败步骤针对性排查：

| 失败步骤 | 可能原因 | 排查方式 |
|---------|---------|---------|
| compile | 语法错误或类型不匹配 | 查看 compile 输出的错误信息，修复后重试 |
| fmt | 代码格式不符合规范 | 运行 `just fmt` 查看差异，自动格式化 |
| lint | 代码质量问题 | 查看 lint 报告，按提示修复 |
| unit-test | 单元测试失败 | 查看测试输出，定位失败的测试用例 |
| test | Surface 级测试失败 | 查看 `results/latest.md` 测试报告 |

```bash
# 手动运行单个 gate 步骤定位问题
just compile
just fmt
just lint
just unit-test
just test

# 查看任务记录
forge task query <task-id> -v
```

**Fix-task 机制**：如果 agent 在执行中无法自行修复，会自动创建 fix-task（类型 `coding.fix`）。完成 fix-task 后，源任务会自动恢复为 `pending` 状态并重新认领执行。

### 6. /run-tasks 无任务可执行

**症状**：执行 `/run-tasks` 后提示没有可用任务。

**排查步骤**：

```bash
# 确认当前 feature 上下文
forge feature

# 查看所有任务状态
forge task list

# 检查 index.json 是否存在
ls docs/features/<slug>/tasks/index.json

# 重建 index.json（如果文件损坏）
forge task index --feature <slug>
```

**常见原因**：
- 未设置 feature 上下文：`forge feature <slug>` 设置
- 所有任务已处于终态（completed/skipped/rejected）
- index.json 缺失或损坏：用 `forge task index --feature <slug>` 重建

### 7. Worktree 创建失败

**症状**：`forge worktree start` 报错。

**排查步骤**：

```bash
# 查看现有 worktree
forge worktree list

# 检查远程分支状态
git fetch --prune

# 使用交互模式选择
forge worktree start -i
```

**常见原因**：
- 远程追踪引用过期：先 `git fetch --prune` 清理过期引用
- 同名分支已存在：先 `forge worktree remove <name>` 或使用不同名称
- 有未推送的 commit：forge 会阻止删除包含未推送 commit 的 worktree，先 push 或确认删除

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [环境准备](environment-setup.md) | 安装 Forge 和配置开发环境 |
| [初始化项目](initialization.md) | 在项目中初始化 Forge |
| [架构概览](architecture-overview.md) | Forge 的架构设计和核心概念 |
