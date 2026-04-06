# Zcode Plugin & Marketplace 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 创建 zcode 插件及其 marketplace，包含任务管理和流程辅助 skills/commands。

**Architecture:** 在 zcode 仓库内创建 `.claude-plugin/marketplace.json` 和 `plugins/zcode/` 目录，将 claude-code-go 中的 skills/commands 复制并改造为使用 `task` CLI 替代 Go 脚本。

**Tech Stack:** Claude Code Plugin System, Markdown, Shell Scripts

---

## File Structure

```
zcode/
├── .claude-plugin/
│   └── marketplace.json           # Marketplace 定义
├── plugins/
│   └── zcode/
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
│           ├── init-zcode.md
│           ├── simplify-skill.md
│           ├── execute-task.md
│           └── run-tasks.md
```

---

### Task 1: 创建目录结构

**Files:**
- Create: `.claude-plugin/`
- Create: `plugins/zcode/.claude-plugin/`
- Create: `plugins/zcode/skills/`
- Create: `plugins/zcode/commands/`

- [ ] **Step 1: 创建 marketplace 目录**

```bash
mkdir -p .claude-plugin
```

- [ ] **Step 2: 创建插件目录结构**

```bash
mkdir -p plugins/zcode/.claude-plugin
mkdir -p plugins/zcode/skills/{claim-task,record-task,set-task-status,breakdown-tasks,write-prd,design-tech,learn-lesson,git-commit}
mkdir -p plugins/zcode/commands
```

- [ ] **Step 3: 验证目录结构**

```bash
ls -la .claude-plugin plugins/zcode/
```

Expected: 目录存在

---

### Task 2: 创建 marketplace.json

**Files:**
- Create: `.claude-plugin/marketplace.json`

- [ ] **Step 1: 创建 marketplace.json**

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
      "description": "Task management and workflow helper tools for Claude Code",
      "version": "1.0.0",
      "keywords": ["task", "workflow", "productivity", "prd", "git"]
    }
  ]
}
```

- [ ] **Step 2: 验证 JSON 格式**

```bash
cat .claude-plugin/marketplace.json | python -m json.tool
```

Expected: JSON 格式正确

---

### Task 3: 创建 plugin.json

**Files:**
- Create: `plugins/zcode/.claude-plugin/plugin.json`

- [ ] **Step 1: 创建 plugin.json**

```json
{
  "name": "zcode",
  "version": "1.0.0",
  "description": "Task management and workflow helper tools for Claude Code",
  "keywords": ["task", "workflow", "productivity", "prd", "git"]
}
```

- [ ] **Step 2: 验证 JSON 格式**

```bash
cat plugins/zcode/.claude-plugin/plugin.json | python -m json.tool
```

Expected: JSON 格式正确

---

### Task 4: 创建 init-zcode 命令

**Files:**
- Create: `plugins/zcode/commands/init-zcode.md`

- [ ] **Step 1: 创建 init-zcode.md**

```markdown
---
name: init-zcode
description: 自动编译并安装 claude-task-cli 工具
---

# /init-zcode

自动编译安装 claude-task-cli 工具。

## 流程

1. 检测操作系统（Windows/Linux/macOS）
2. 定位 claude-task-cli 仓库路径（与 zcode 平级目录）
3. 调用对应安装脚本编译并安装
4. 提示用户重新打开终端

## 执行步骤

### Step 1: 定位 claude-task-cli

```bash
# 检查是否存在 claude-task-cli
TASK_CLI_DIR="$(dirname "${CLAUDE_PROJECT_ROOT}")/claude-task-cli"
if [ ! -d "$TASK_CLI_DIR" ]; then
  echo "ERROR: claude-task-cli not found at $TASK_CLI_DIR"
  echo "Please clone claude-task-cli to the parent directory of zcode"
  exit 1
fi
echo "Found claude-task-cli at: $TASK_CLI_DIR"
```

### Step 2: 检测操作系统并安装

**Windows (PowerShell):**
```powershell
cd $TASK_CLI_DIR
powershell -ExecutionPolicy Bypass -File scripts/install-local.ps1
```

**Linux/macOS:**
```bash
cd "$TASK_CLI_DIR" && bash scripts/install-local.sh
```

### Step 3: 验证安装

```bash
task --version
```

### Step 4: 提示用户

安装完成后，输出：

```
╔════════════════════════════════════════════════════════════════╗
║  ✅ claude-task-cli 安装成功                                    ║
╠════════════════════════════════════════════════════════════════╣
║  请重新打开终端以刷新环境变量，然后运行:                           ║
║  task --version                                                ║
╚════════════════════════════════════════════════════════════════╝
```

## 错误处理

| 错误 | 解决方案 |
|------|----------|
| claude-task-cli 未找到 | 克隆仓库到 zcode 平级目录 |
| 编译失败 | 检查 Go 环境 |
| 权限错误 | 检查安装目录写入权限 |
```

- [ ] **Step 2: 验证文件创建**

```bash
cat plugins/zcode/commands/init-zcode.md
```

Expected: 文件内容正确

---

### Task 5: 复制并改造 Skills (Part 1: 任务管理)

**Files:**
- Create: `plugins/zcode/skills/claim-task/SKILL.md`
- Create: `plugins/zcode/skills/record-task/SKILL.md`
- Create: `plugins/zcode/skills/set-task-status/SKILL.md`
- Create: `plugins/zcode/skills/breakdown-tasks/SKILL.md`

- [ ] **Step 1: 创建 claim-task/SKILL.md**

```markdown
---
name: claim-task
description: Use when you need to claim and start working on the next available task from the project task list. Claims the highest priority task with all dependencies met.
---

# Claim Task

Claim the next available task from `docs/tasks/index.json`.

## Usage

```bash
task claim
```

## Output

```
---
KEY: <task-key>
ID: <task-id>
FILE: <task-file>
---
```

## After Claiming

1. Read task file: `docs/tasks/<FILE>`
2. Implement following TDD (RED → GREEN → REFACTOR)
3. Update record: `/record-task <TASK_ID>`
4. Mark complete: `/set-task-status <ID> completed`

## Related

- `/set-task-status` - Update task status
- `/run-tasks` - Auto-execute all tasks
- `/execute-task` - Manual single task workflow
```

- [ ] **Step 2: 创建 record-task/SKILL.md**

```markdown
---
name: record-task
description: Use after completing a task to create its execution record and update task status.
---

# Record Task

## Overview

任务完成后的收尾操作：创建执行记录 + 更新任务状态。

## Usage

```bash
# 使用 JSON 文件
echo '{"summary":"...","filesCreated":[...],"filesModified":[...]}' > .claude/process/record.json
task record <TASK_ID> -data .claude/process/record.json
```

## JSON Data Format

```json
{
  "status": "completed",
  "summary": "实现了什么",
  "filesCreated": ["internal/app/foo.go"],
  "filesModified": ["internal/app/bar.go"],
  "keyDecisions": ["决策 1"],
  "testsPassed": 12,
  "testsFailed": 0,
  "coverage": 85.6,
  "acceptanceCriteria": [
    {"criterion": "验收标准 1", "met": true}
  ]
}
```

## Fields

| 字段 | 类型 | 说明 |
|------|------|------|
| `status` | string | 任务状态，默认 `completed` |
| `summary` | string | 实现摘要 |
| `filesCreated` | array | 新建文件列表 |
| `filesModified` | array | 修改文件列表 |
| `keyDecisions` | array | 关键设计决策 |
| `testsPassed` | int | 通过测试数 |
| `testsFailed` | int | 失败测试数 |
| `coverage` | float | 覆盖率 |
| `acceptanceCriteria` | array | `{criterion, met}` 对象 |

## Related

- `/claim-task` - Claim next available task
- `/set-task-status` - Direct status update only
```

- [ ] **Step 3: 创建 set-task-status/SKILL.md**

```markdown
---
name: set-task-status
description: Use when you need to update the status of a task. Supports pending, in_progress, completed, blocked, skipped statuses.
---

# Set Task Status

## Overview

Update the status of a task in `docs/tasks/index.json`.

## Usage

```bash
task status <task-id-or-key> <status>
```

## Valid Statuses

| Status | When to Use |
|--------|-------------|
| `pending` | Task not started (default) |
| `in_progress` | Currently working on it |
| `completed` | Task finished successfully |
| `blocked` | Cannot proceed due to external issue |
| `skipped` | Task not needed |

## Workflow

```
pending → in_progress → completed
                 ↓
              blocked → in_progress
                 ↓
              skipped
```

## Related

- `/claim-task` - Claim next available task
- `/run-tasks` - Auto-execute all tasks
```

- [ ] **Step 4: 创建 breakdown-tasks/SKILL.md**

```markdown
---
name: breakdown-tasks
description: Use when design.md is finalized to break down into executable tasks. Creates task files based on technical design.
---

# Breakdown Tasks

## Overview

从技术设计文档拆解成可执行的任务。

**核心原则**：任务粒度适中（1-4 小时），依赖关系明确，验收标准可测试。

## Position in Workflow

```
/write-prd → /design-tech → /breakdown-tasks
     ↓              ↓              ↓
   prd.md      design.md      tasks/*.md
```

## Directory Structure

```
docs/features/<feature-slug>/
├── prd.md              # PRD document
├── design.md           # Technical design (input)
├── index.json          # Task index for this feature
├── tasks/              # Task definitions
└── records/            # Execution records
```

## When to Use

**Trigger conditions:**
- Design document exists at `docs/features/<slug>/design.md`
- User asks to "break down" or "split" a design into tasks

**Skip when:**
- No design.md exists (use `/design-tech` first)
- Tasks already exist for the feature

## Workflow

```
1. Read Design → 2. Map interfaces → 3. Define order → 4. Create task files → 5. Create index.json → 6. Validate
```

## Step 1: Read Design

Read `docs/features/<slug>/design.md`:
- Understand architecture and component structure
- Map interfaces to implementation tasks
- Identify data models and their tasks

## Step 2: Map Interfaces to Tasks

| Design Element | Task Type |
|----------------|-----------|
| Interface definition | Interface task |
| Data model | Model task |
| Component | Implementation task |
| Error type | Error handling task |

## Step 3: Define Task Order

```
1.x Interfaces → 2.x Models → 3.x Implementation → 4.x Integration → 5.x Tests
```

## Step 4: Create Task Files

**Naming convention:**
```
<sequence>.<sub-sequence>-<slug>.md
```

## Step 5: Create index.json

Create `docs/features/<slug>/index.json` with task definitions.

## Step 6: Validate

```bash
task validate -file docs/features/<slug>/index.json
```

## Integration

Works well with:
- `/design-tech` - Creates the design.md input
- `/claim-task` - Starts working on tasks
- `/record-task` - Records task completion
```

- [ ] **Step 5: 验证任务管理 skills 创建**

```bash
ls plugins/zcode/skills/
```

Expected: claim-task, record-task, set-task-status, breakdown-tasks 目录存在

---

### Task 6: 复制并改造 Skills (Part 2: 流程辅助)

**Files:**
- Create: `plugins/zcode/skills/write-prd/SKILL.md`
- Create: `plugins/zcode/skills/design-tech/SKILL.md`
- Create: `plugins/zcode/skills/learn-lesson/SKILL.md`
- Create: `plugins/zcode/skills/git-commit/SKILL.md`

- [ ] **Step 1: 创建 write-prd/SKILL.md**

```markdown
---
name: write-prd
description: Use when user provides requirements that need to be formalized into a PRD document through collaborative dialogue.
---

# Write PRD

## Overview

从模糊需求产出清晰的 PRD（产品需求文档）。

<HARD-GATE>
Do NOT write any code until the PRD is finalized and approved.
</HARD-GATE>

## When to Use

**Trigger conditions:**
- User describes a feature without clear specifications
- User says "I want to..." or "We need..." without details
- Starting a new phase or major feature

**Skip when:**
- Clear task definitions already exist
- Simple bug fix or small tweak

## Process Flow

```
Explore context → Assess scope → Ask questions → Propose approaches → Present PRD → Write doc → Commit
```

## Checklist

1. **Explore project context** — check files, docs, recent commits
2. **Assess scope** — determine if request needs decomposition
3. **Ask clarifying questions** — one at a time via AskUserQuestion
4. **Propose 2-3 approaches** — with trade-offs and recommendation
5. **Present PRD sections** — get approval after each section
6. **Write PRD document** — save to `docs/features/<feature-slug>/prd.md`
7. **Commit**

## Step 1: Explore Context

- Read `docs/ARCHITECTURE.md`
- Read `docs/DECISIONS.md`
- Check `docs/tasks/index.json`
- Review recent git commits

## Step 2: Assess Scope

- Multiple independent subsystems → **Decompose first**
- Single focused feature → **Proceed with questions**

## Step 3: Ask Questions

**CRITICAL**: Use `AskUserQuestion` tool. One question at a time.

## Step 4: Propose Approaches

Present 2-3 approaches with trade-offs. Lead with recommendation.

## Step 5: Present PRD

Present incrementally, get approval after each section:

| Section | Content |
|---------|---------|
| Background | Problem statement, context |
| Goals | Primary goals, success metrics |
| Scope | In/out of scope items |
| Requirements | Functional requirements |
| Acceptance Criteria | Testable conditions |

## Step 6: Write Document

```
docs/features/<feature-slug>/
├── prd.md
├── tasks/
└── records/
```

## Integration

Works well with:
- `/breakdown-tasks` - Break PRD into tasks
- `/design-tech` - Create technical design
```

- [ ] **Step 2: 创建 design-tech/SKILL.md**

```markdown
---
name: design-tech
description: Use after PRD is finalized to create technical design with architecture and implementation details.
---

# Design Tech

## Overview

从 PRD 产出技术设计文档。

<HARD-GATE>
Do NOT write any implementation code until design.md is approved.
</HARD-GATE>

## Position in Workflow

```
/write-prd → /design-tech → /breakdown-tasks
     ↓              ↓              ↓
   prd.md      design.md      tasks/*.md
```

## When to Use

**Trigger conditions:**
- PRD document exists at `docs/features/<slug>/prd.md`
- PRD is approved and ready for technical design

**Skip when:**
- No PRD exists (use `/write-prd` first)
- Design already exists for the feature

## Process Flow

```
1. Read PRD → 2. Explore context → 3. Identify decisions → 4. Ask questions → 5. Draft design → 6. Review → 7. Finalize
```

## Step 1: Read PRD

Read `docs/features/<slug>/prd.md`:
- Understand requirements
- Note non-functional requirements
- Identify acceptance criteria

## Step 2: Explore Context

| Source | What to Look For |
|--------|------------------|
| `docs/ARCHITECTURE.md` | Layer constraints |
| `docs/DECISIONS.md` | Existing decisions |
| `go.mod` | Current dependencies |
| `internal/` | Existing patterns |

## Step 3: Identify Decisions

| Decision Type | Example Questions |
|---------------|-------------------|
| Architecture | Where does this fit? |
| Interface | What interfaces needed? |
| Data Model | What structures needed? |
| Dependencies | New dependencies? |
| Error Handling | How to handle errors? |
| Testing | Test strategy? |
| Security | Security considerations? |

## Step 4: Ask Questions

Use `AskUserQuestion` for ALL uncertain areas.

## Step 5: Draft Design

Present incrementally, section by section:

| Section | Content |
|---------|---------|
| Overview | High-level approach |
| Architecture | Component diagram |
| Interfaces | Interface definitions |
| Data Models | Struct definitions |
| Error Handling | Error strategy |
| Testing | Test strategy |
| Security | Security considerations |

## Step 6: Get Approval

For each section, wait for user approval.

## Step 7: Write design.md

Save to `docs/features/<slug>/design.md`

## Integration

Works well with:
- `/write-prd` - Creates PRD input
- `/breakdown-tasks` - Uses design.md to create tasks
```

- [ ] **Step 3: 创建 learn-lesson/SKILL.md**

```markdown
---
name: learn-lesson
description: Use when you have solved an error or discovered a useful pattern. Extracts reusable knowledge from the current session.
---

# Learn Lesson

## Overview

从当前会话中提取可复用的知识点，记录到 `docs/lessons/`。

**核心原则**：记录"下次遇到类似问题可以怎么处理"，不是"我做了什么"。

## When to Use

**Trigger conditions:**
- Solved a non-trivial error with reusable insights
- Discovered a pattern/technique worth remembering
- User explicitly requests `/learn-lesson`

**Skip when:**
- Trivial typo fixes
- One-off environment issues
- Information already documented elsewhere

## Workflow

```
1. Identify lesson → 2. Classify category → 3. Write doc → 4. Commit
```

## Step 1: Identify Lesson

- 遇到的问题（症状）
- 根本原因
- 解决方案
- **可复用的知识点**

## Step 2: Classify Category

| Category | Prefix | Example |
|----------|--------|---------|
| Debugging | `debug-` | `debug-race-condition.md` |
| Architecture | `arch-` | `arch-dependency-direction.md` |
| Tooling | `tool-` | `tool-go-test-coverage.md` |
| Pattern | `pattern-` | `pattern-error-wrapping.md` |
| Gotcha | `gotcha-` | `gotcha-context-cancellation.md` |

## Step 3: Write Document

Output: `docs/lessons/<category-prefix><slug>.md`

## Step 4: Commit

```bash
git add docs/lessons/<filename>.md
git commit -m "docs(lessons): add <title>"
```

## Common Mistakes

| Mistake | Correction |
|---------|------------|
| Recording "what I did" | Focus on "what to do next time" |
| Too specific | Generalize to reusable pattern |
| Missing root cause | Always include why |
```

- [ ] **Step 4: 创建 git-commit/SKILL.md**

```markdown
---
name: git-commit
description: Use when creating git commits. Ensures commit messages follow Conventional Commits format.
---

# Git Commit

## Overview

Follow Conventional Commits specification with project-defined type and scope rules.

## Atomic Commits

**Core principle**: group high-related changes; split unrelated changes.

| Practice | Description |
|----------|-------------|
| **Group together** | feature + its tests + its docs in one commit |
| **Split apart** | unrelated features, fixes, independent refactor |

## Format

```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

## Allowed Types

| Type | When to Use |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `test` | Adding/modifying tests |
| `refactor` | Code refactoring |
| `chore` | Maintenance, tooling, deps |

## Subject Rules

1. **Lowercase first letter** - `add` not `Add`
2. **No trailing period**
3. **Imperative mood** - `add` not `added`
4. **Max 72 characters**

## Examples

```bash
# Good
feat(api): add streaming support
fix(parser): handle empty input
docs(readme): update install steps

# Bad
Update(api): add streaming    # Wrong type
feat(api): Added support.     # Past tense, period
```

## Quick Checklist

- [ ] Type is one of: feat / fix / docs / test / refactor / chore
- [ ] Scope matches affected module
- [ ] Subject starts with lowercase
- [ ] Subject has no trailing period
- [ ] Subject is imperative mood
```

- [ ] **Step 5: 验证流程辅助 skills 创建**

```bash
ls plugins/zcode/skills/
```

Expected: 所有 8 个 skill 目录存在

---

### Task 7: 复制 Commands

**Files:**
- Create: `plugins/zcode/commands/simplify-skill.md`
- Create: `plugins/zcode/commands/execute-task.md`
- Create: `plugins/zcode/commands/run-tasks.md`

- [ ] **Step 1: 创建 simplify-skill.md**

```markdown
---
name: simplify-skill
description: Refactor skill files by extracting templates/examples to separate files.
argument-hints: skill name
allowed_tools: ["Read", "AskUserQuestion"]
---

# /simplify-skill

重构 skill 文件：**拆分非核心内容**，保持主流程清晰。

## Core Principle

```
MAIN FILE = WORKFLOW ONLY

skill.md     → 流程步骤、决策点
templates/   → JSON 模板、输出格式示例
examples/    → 完整用例、边界情况
```

## Workflow

```
1. Identify Skill → 2. Analyze Content → 3. Ask Approval → 4. Extract
```

## Phase 1: Identify Target

If no argument provided, ask user which skill to refactor.

Target locations:
- Skills: `.claude/skills/<name>/SKILL.md`
- Commands: `.claude/commands/<name>.md`

## Phase 2: Analyze Extractables

| Category | Indicators | Extract To |
|----------|-----------|------------|
| Templates | JSON blocks, output formats | `templates/` |
| Examples | Multi-line code samples | `examples/` |
| Reference tables | Field definitions | `reference.md` |
| Verbose context | Background explanations | `context.md` |

## Phase 3: Ask for Approval

Use `AskUserQuestion` with multiSelect for which content to extract.

## Phase 4: Execute Extraction

1. Create directory structure
2. Extract content to new files
3. Replace in skill.md with reference
4. Keep workflow steps intact

## Iron Law

```
NEVER extract without user approval
NEVER remove content, only relocate
ALWAYS add file references
KEEP workflow steps intact
```
```

- [ ] **Step 2: 创建 execute-task.md**

```markdown
---
name: execute-task
description: Execute single task with focused TDD workflow.
allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]
---

# /execute-task

Execute a single task with streamlined TDD workflow.

## Workflow (5 Steps)

```
Step 1: Read task definition
Step 2: TDD (RED → GREEN → REFACTOR)
Step 3: Full verification
Step 4: Record task (MANDATORY)
Step 5: Git commit
```

## Step 1: Claim & Read

```bash
task claim
```

Parse output for KEY, ID, FILE. Read task file.

## Step 2: TDD Implementation

```
RED      → Write failing test first
GREEN    → Implement minimal code to pass
REFACTOR → Clean up while keeping tests green
```

## Step 3: Full Verification

Run project-specific verification commands.

## Step 4: Record Task (MANDATORY)

```bash
echo '{"summary":"..."}' > .claude/process/record.json
task record <TASK_ID> -data .claude/process/record.json
```

## Step 5: Commit

```
Skill(skill="git-commit")
```

## Rules

- **record-task is mandatory** - No completion without it
- **All verifications must pass**
- **Commit only after record**

## Related Commands

| Command | Usage |
|---------|-------|
| `/run-tasks` | Auto-execute all tasks |
| `/claim-task` | Claim task only |
| `/record-task` | Create record + update status |
```

- [ ] **Step 3: 创建 run-tasks.md**

```markdown
---
name: run-tasks
description: Autonomous task dispatcher that continuously claims tasks and dispatches to subagents.
allowed_tools: ["Bash", "Read", "Agent", "TaskOutput"]
---

# /run-tasks

Auto-dispatch tasks to subagents. Main session only handles dispatching.

## Architecture

```
MAIN SESSION (Dispatcher)
   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
   │ 1. Claim    │───▶│ 2. Dispatch │───▶│ 3. Verify   │
   │    Task     │    │   + Timeout │    │   Record    │
   └─────────────┘    └─────────────┘    └─────────────┘
          ▲                                    │
          └────────────────────────────────────┘
                      LOOP
```

## Dispatcher Iron Laws

```
1. Only 3 actions: claim → dispatch → verify
2. NO code reading, NO code writing
3. NO running tests directly
4. 30-minute timeout per task
5. 3 consecutive failures → STOP
```

## Execution Loop

### Step 1: Claim Task

```bash
task claim
```

**Output parsing**:
- `ACTION: CLAIMED` → New task
- `ACTION: CONTINUE` → Resume existing task
- Error → No available task, end loop

### Step 2: Dispatch with Timeout

```
Agent(
  subagent_type="task-executor",
  prompt="TASK_KEY: {{KEY}}
TASK_ID: {{ID}}
TASK_FILE: {{FILE}}",
  isolation="worktree"
)
```

**Timeout**: 30 minutes

### Step 3: Verify Record

Check if record file exists after agent completes.

### Step 4: Continue Loop

Return to Step 1.

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | End loop, print summary |
| Agent timeout | Mark blocked, continue next |
| Record missing | Dispatch error-fixer |
| 3 consecutive failures | STOP dispatcher |

## Related Commands

| Command | Usage |
|---------|-------|
| `/execute-task` | Manual single task |
| `/claim-task` | Claim task only |
| `/record-task` | Create record + update status |
```

- [ ] **Step 4: 验证 commands 创建**

```bash
ls plugins/zcode/commands/
```

Expected: init-zcode.md, simplify-skill.md, execute-task.md, run-tasks.md

---

### Task 8: 验证插件结构

**Files:**
- Verify: 完整插件结构

- [ ] **Step 1: 验证目录结构**

```bash
find plugins/zcode -type f -name "*.md" -o -name "*.json" | sort
```

Expected 输出:
```
plugins/zcode/.claude-plugin/plugin.json
plugins/zcode/commands/execute-task.md
plugins/zcode/commands/init-zcode.md
plugins/zcode/commands/run-tasks.md
plugins/zcode/commands/simplify-skill.md
plugins/zcode/skills/breakdown-tasks/SKILL.md
plugins/zcode/skills/claim-task/SKILL.md
plugins/zcode/skills/design-tech/SKILL.md
plugins/zcode/skills/git-commit/SKILL.md
plugins/zcode/skills/learn-lesson/SKILL.md
plugins/zcode/skills/record-task/SKILL.md
plugins/zcode/skills/set-task-status/SKILL.md
plugins/zcode/skills/write-prd/SKILL.md
```

- [ ] **Step 2: 验证 marketplace.json**

```bash
cat .claude-plugin/marketplace.json | python -m json.tool
```

Expected: JSON 格式正确

---

## Summary

| Task | Description | Files Created |
|------|-------------|---------------|
| 1 | 创建目录结构 | 4 directories |
| 2 | marketplace.json | 1 file |
| 3 | plugin.json | 1 file |
| 4 | init-zcode 命令 | 1 file |
| 5 | 任务管理 skills | 4 files |
| 6 | 流程辅助 skills | 4 files |
| 7 | Commands | 3 files |
| 8 | 验证 | - |

**Total**: 2 JSON files + 12 Markdown files
