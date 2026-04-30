---
created: 2026-04-30
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Justfile 标准命令词汇表

## Overview

将 justfile 从 6 个标准目标扩展至 15 个，通过三类变更实现：

1. **init-justfile 契约扩展**：更新命令定义文件，增加自适应生成逻辑（3 种项目类型 × 15 个命令模板）
2. **task-cli schema 扩展**：在 Go `Task` struct 中添加 `scope` 字段
3. **skill/agent/command 迁移**：14 处机械替换 + 1 处 prompt 工程修改（breakdown-tasks scope 标注）

所有变更限定在 forge 插件层和 task-cli 层，不涉及运行时服务或数据库。

## Skill/Agent/Command Migration

### Implementation Tasks（5 项需代码变更）

| # | 文件 | 当前命令 | 迁移后命令 | 操作 |
|---|------|---------|-----------|------|
| 1 | `run-e2e-tests/SKILL.md` | `npx serve` | `just run` | sed 替换 `npx serve` → `just run` |
| 2 | `execute-task.md` (agent) | `just build && just test` | `just compile && just test` | sed 替换 `just build` → `just compile` |
| 3 | `task-executor.md` (agent) | `just build && just test` | `just compile && just test` | sed 替换 `just build` → `just compile` |
| 4 | `error-fixer.md` (agent) | `just build && just test` | `just compile && just test` | sed 替换 `just build` → `just compile` |
| 5 | `breakdown-tasks/SKILL.md` | — | 添加 scope 标注逻辑（见下方） | prompt 工程：插入 Scope Assignment 段落 |

### Validation Checklist（10 项确认无需改动）

| # | 文件 | 当前命令 | 验证方法 |
|---|------|---------|---------|
| 6 | `gen-test-scripts/SKILL.md` | `just e2e-setup`、`just e2e-verify` | `grep -c 'just e2e-setup\|just e2e-verify' plugins/forge/skills/gen-test-scripts/SKILL.md` 返回匹配数 ≥ 1 |
| 7 | `record-task/SKILL.md` | `just test` | `grep -c 'just test' plugins/forge/skills/record-task/SKILL.md` 返回匹配数 ≥ 1 |
| 8 | `fix-bug.md` | `just test`、`just test-e2e` | `grep -c 'just test' plugins/forge/agents/fix-bug.md` 返回匹配数 ≥ 1 |
| 9 | `run-tasks.md` | `just test` | `grep -c 'just test' plugins/forge/agents/run-tasks.md` 返回匹配数 ≥ 1 |
| 10 | `improve-harness.md` | `just test` | `grep -c 'just test' plugins/forge/skills/improve-harness/SKILL.md` 返回匹配数 ≥ 1 |
| 11 | `graduate-tests/SKILL.md` | `just test-e2e` | `grep -c 'just test-e2e' plugins/forge/skills/graduate-tests/SKILL.md` 返回匹配数 ≥ 1 |
| 12 | task 模板 1 | `just test-e2e`、`just e2e-verify` | `grep -c 'just test-e2e\|just e2e-verify' <template-path-1>` 返回匹配数 ≥ 1 |
| 13 | task 模板 2 | `just test-e2e`、`just e2e-verify` | `grep -c 'just test-e2e\|just e2e-verify' <template-path-2>` 返回匹配数 ≥ 1 |
| 14 | task 模板 3 | `just test-e2e`、`just e2e-verify` | `grep -c 'just test-e2e\|just e2e-verify' <template-path-3>` 返回匹配数 ≥ 1 |
| 15 | task 模板 4 | `just test-e2e`、`just e2e-verify` | `grep -c 'just test-e2e\|just e2e-verify' <template-path-4>` 返回匹配数 ≥ 1 |

### breakdown-tasks Scope 标注（第 15 项：prompt 工程修改）

在 `breakdown-tasks/SKILL.md` 的 Step 4a（Business Tasks）之后插入以下 scope 标注指令：

**新增到 SKILL.md 的文本**：

```markdown
### Scope Assignment

For each task, determine the `scope` field for `index.json`:

**Algorithm**: inspect the task's affected file paths (listed in the task's "Files Created/Modified" section derived from the tech-design).

1. Classify each file path:
   - `frontend`: path starts with `ui/`, `src/`, `components/`, `pages/`, `styles/`, `public/`, or any directory containing `package.json` with no `go.mod`/`Cargo.toml` at the same level
   - `backend`: path starts with `cmd/`, `internal/`, `pkg/`, `api/`, or any directory containing `go.mod`/`Cargo.toml`/`pyproject.toml` with no `package.json` at the same level
   - `undetermined`: path does not match either pattern (e.g., `docs/`, root config files, `justfile`)

2. Compute scope:
   - If ALL paths are `frontend` → `scope: "frontend"`
   - If ALL paths are `backend` → `scope: "backend"`
   - Otherwise (mixed paths, `undetermined` paths, or no file paths) → `scope: "all"`

3. Write `scope` into the task entry in `index.json`.

**Non-mixed projects**: when `init-justfile` detects a pure frontend or backend project, all tasks receive `scope: "all"` (scope distinction is irrelevant when `just project-type` does not return `"mixed"`).
```

**行为示例**：

| 任务涉及文件路径 | scope 值 | 原因 |
|----------------|---------|------|
| `ui/components/Button.tsx`, `src/styles.css` | `frontend` | 全部前端路径 |
| `cmd/server/main.go`, `pkg/handler/api.go` | `backend` | 全部后端路径 |
| `ui/App.tsx`, `cmd/server/main.go` | `all` | 前后端混合 |
| `docs/WORKFLOW.md`, `justfile` | `all` | 无法确定归属 |
| 纯后端项目中任何任务 | `all` | 非混合项目，scope 区分无意义 |

## Architecture

### Layer Placement

```
┌─────────────────────────────────────────────────┐
│  Claude Code Agent (执行层)                      │
│  ┌───────────┐  ┌───────────┐  ┌──────────────┐ │
│  │ execute-  │  │ run-tasks │  │ error-fixer  │ │
│  │ task      │  │ (dispatch)│  │ (fix loop)   │ │
│  └─────┬─────┘  └─────┬─────┘  └──────┬───────┘ │
│        │              │                │         │
│        ▼              ▼                ▼         │
│  ┌─────────────────────────────────────────┐    │
│  │  justfile 命令适配层                      │    │
│  │  compile / build / test / run / lint ... │    │
│  └─────────────┬───────────────────────────┘    │
│                │                                 │
│        ┌───────┴───────┐                        │
│        ▼               ▼                        │
│  ┌──────────┐   ┌──────────┐                   │
│  │ Go CLI   │   │ Node.js  │                   │
│  │ (backend)│   │(frontend)│                   │
│  └──────────┘   └──────────┘                   │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  Forge Plugin (配置层)                           │
│  ┌──────────────┐  ┌────────────────────────┐   │
│  │ init-justfile │  │ breakdown-tasks        │   │
│  │ (生成 justfile)│  │ (生成 index.json+scope)│   │
│  └──────────────┘  └────────────────────────┘   │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│  task-cli (数据层)                               │
│  ┌──────────────────────────────────────────┐   │
│  │ Task struct: id, title, scope, status...  │   │
│  │ TaskIndex: tasks map + metadata           │   │
│  └──────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

### Component Diagram

```
init-justfile ──generates──▶ justfile (project root, merge within boundary markers)
                                │
justfile ◀──calls────────── skill/agent/command files
   │
   ├── project-type ──outputs──▶ "frontend" / "backend" / "mixed"
   │
   └── compile/build/test/run/... ──dispatch──▶ language toolchain

breakdown-tasks ──generates──▶ index.json
                                  │
                                  └── tasks[].scope ──reads──▶ skill execution
```

### Dependencies

| 依赖 | 类型 | 说明 |
|------|------|------|
| just >= 1.50.0 | 外部 | `[arg("feature", long)]` 语法、位置参数默认值 |
| task-cli (Go) | 内部 | `Task` struct 需添加 `scope` 字段 |
| `task-cli/pkg/task/types.go` | 内部 | `Task` 和 `TaskState` struct 定义所在文件，需添加 `Scope` 字段 |
| `index.schema.json` | 内部 | JSON Schema 验证文件，需添加 `scope` 属性定义（Model 3） |
| forge plugin skills | 内部 | 8 个 skill 文件 + 2 个 agent 文件 + 4 个 task 模板 |

## Interfaces

### Interface 1: justfile Recipe Contract

每个 recipe 遵循统一的签名模式：

```
// 纯前端/纯后端项目（无 scope）
recipe_name:
    <language-specific-command>

// 混合项目（有 scope）
recipe_name scope="":
    #!/usr/bin/env bash
    case "{{scope}}" in
      frontend) <frontend-command> ;;
      backend)  <backend-command> ;;
      "")       <frontend-command> && <backend-command> ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac
```

**约束**：
- scope 参数类型：`string`，合法值为 `frontend`、`backend`、`""`（空=全部）
- 非法 scope：退出码 1，错误信息到 stderr
- 所有 recipe 使用 `#!/usr/bin/env bash` shebang（确保跨平台兼容）
- 正常输出到 stdout，错误输出到 stderr

### Interface 2: project-type Recipe

```
project-type:
    @echo "frontend"  // 或 "backend" 或 "mixed"
```

**约束**：
- 输出恰好一个单词 + 换行，无额外文本
- 退出码始终为 0
- 无副作用，可重复调用
- 由 init-justfile 自动生成，非人工维护

### Interface 3: scope 字段扩展（index.json）

```
Task struct 扩展:
  Scope: string  // json:"scope,omitempty", enum: frontend/backend/all, default: all
```

**消费者协议**：
- skill 读取 `tasks[id].scope`
- 若字段缺失或为空，视为 `"all"`
- 仅在 `just project-type` 返回 `"mixed"` 时使用 scope 参数

### Interface 4: Skill Scope Resolution Protocol

skill 文件中粘贴以下标准 scope 解析块（作为 prompt 文本嵌入 SKILL.md）：

```markdown
## Scope Resolution

Before executing any `just <verb>` command, resolve scope:

1. Read the current task's `scope` field from index.json (or task claim output).
2. If `scope` is missing, empty, or `"all"` → execute `just <verb>` (no scope argument). Done.
3. If `scope` is `"frontend"` or `"backend"`:
   a. Run `just project-type` and capture stdout (trimmed) and exit code.
   b. If exit code != 0, or stdout is not one of `frontend`/`backend`/`mixed`:
      - Log: `[forge] just project-type failed (exit N); falling back to just <verb>`
      - Execute `just <verb>` (no scope). Done.
   c. If stdout == `"mixed"`:
      - Execute `just <verb> <scope>` (e.g., `just build frontend`). Done.
   d. If stdout is `"frontend"` or `"backend"` (not mixed):
      - Log: `[forge] scope=<scope> but project-type=<type>; falling back to just <verb>`
      - Execute `just <verb>` (no scope). Done.
```

**Worked examples**:

| Task scope | `just project-type` output | Action |
|-----------|---------------------------|--------|
| `frontend` | `mixed` (exit 0) | `just build frontend` |
| `backend` | `mixed` (exit 0) | `just test backend` |
| `all` | (not called) | `just compile` |
| `frontend` | `backend` (exit 0) | `[forge] scope=frontend but project-type=backend; falling back` → `just build` |
| `frontend` | exit 127 (command not found) | `[forge] just project-type failed (exit 127); falling back` → `just build` |
| `backend` | `unknown` (exit 0) | `[forge] just project-type failed; falling back` → `just test` |

## Data Models

### Model 1: Task struct（Go）

```go
// file: task-cli/pkg/task/types.go
Task = {
    ID:            string   // json:"id"
    Title:         string   // json:"title"
    Priority:      string   // json:"priority"
    EstimatedTime: string   // json:"estimatedTime,omitempty"
    Dependencies:  []string // json:"dependencies,omitempty"
    Status:        string   // json:"status"
    File:          string   // json:"file"
    Record:        string   // json:"record"
    Breaking:      bool     // json:"breaking,omitempty"
    Scope:         string   // json:"scope,omitempty"   // NEW: frontend/backend/all
}
```

### Model 2: TaskState struct（Go）

```go
// file: task-cli/pkg/task/types.go
TaskState = {
    TaskID:        string   // json:"task_id"
    Key:           string   // json:"key"
    Title:         string   // json:"title"
    Priority:      string   // json:"priority"
    EstimatedTime: string   // json:"estimatedTime,omitempty"
    Dependencies:  []string // json:"dependencies,omitempty"
    File:          string   // json:"file"
    Record:        string   // json:"record"
    StartedTime:   string   // json:"startedTime"
    Breaking:      bool     // json:"breaking,omitempty"
    Scope:         string   // json:"scope,omitempty"   // NEW: mirrors Task.Scope
}
```

### Model 3: index.schema.json 扩展

```json
{
  "tasks": {
    "additionalProperties": {
      "properties": {
        "scope": {
          "type": "string",
          "enum": ["frontend", "backend", "all"],
          "default": "all",
          "description": "Task scope: frontend/backend/all"
        }
      }
    }
  }
}
```

### Model 4: init-justfile 项目类型探测映射

```
ProjectDetection = {
    signals: {
        "package.json":   "frontend",
        "go.mod":         "backend",
        "Cargo.toml":     "backend",
        "pyproject.toml": "backend"
    }
    classify(detected: Set<string>) => "frontend" | "backend" | "mixed"
    // mixed = detected contains both frontend and backend signals
    // frontend = only frontend signals
    // backend = only backend signals
}
```

### Model 5: Recipe Template Examples

模板以字符串字面量形式内嵌于 `init-justfile` 命令文件中，按项目类型选择性拼接生成。生成时包裹于 `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` 边界标记内，已有 justfile 仅替换边界标记间内容。

**混合项目 `build` recipe（带 scope）**：

```just
# --- forge standard recipes ---
build scope="":
    #!/usr/bin/env bash
    set -euo pipefail
    case "{{scope}}" in
      frontend) npm run build ;;
      backend)  go build ./... ;;
      "")       npm run build && go build ./... ;;
      *)        echo "[forge] invalid scope '{{scope}}'; expected frontend/backend" >&2; exit 1 ;;
    esac
```

**纯后端项目 `build` recipe（无 scope）**：

```just
# --- forge standard recipes ---
build:
    go build ./...
```

**`project-type` recipe（三种变体）**：

```just
# 混合项目
project-type:
    @echo "mixed"

# 纯前端项目
project-type:
    @echo "frontend"

# 纯后端项目
project-type:
    @echo "backend"
```

**各项目类型生成的 recipe 清单**：

| Recipe | 混合项目 (mixed) | 纯前端 (frontend) | 纯后端 (backend) |
|--------|------------------|-------------------|------------------|
| `project-type` | `@echo "mixed"` | `@echo "frontend"` | `@echo "backend"` |
| `compile` | bash case: `tsc --noEmit` / `go vet ./...` | `npm run compile` | `go vet ./...` |
| `build` | bash case: `npm run build` / `go build ./...` | `npm run build` | `go build ./...` |
| `run` | bash case: `npm start` / `go run .` | `npm start` | `go run .` |
| `dev` | bash case: `npm run dev` / `go run . --dev` | `npm run dev` | `go run . --dev` |
| `test` | bash case: `npm test` / `go test -race ./...` | `npm test` | `go test -race ./...` |
| `test-e2e` | 同一模板（所有类型共用） | 同左 | 同左 |
| `lint` | bash case: `npm run lint` / `golangci-lint run ./...` | `npm run lint` | `golangci-lint run ./...` |
| `fmt` | bash case: `npx prettier --write .` / `gofmt -w .` | `npx prettier --write .` | `gofmt -w .` |
| `check` | bash case: `npm run lint && npx tsc --noEmit` / `golangci-lint run ./...` | `npm run lint && npx tsc --noEmit` | `golangci-lint run ./...` |
| `clean` | bash case: `rm -rf dist` / `go clean ./...` | `rm -rf dist` | `go clean ./...` |
| `install` | bash case: `npm install` / `go mod download` | `npm install` | `go mod download` |
| `ci` | 无 scope（固定流水线） | 同左 | 同左 |
| `e2e-setup` | 同一模板（所有类型共用） | 同左 | 同左 |
| `e2e-verify` | 同一模板（所有类型共用） | 同左 | 同左 |

混合项目含 scope 参数的 recipe（10 个）：compile、build、run、dev、test、lint、fmt、check、clean、install。
不含 scope 参数的 recipe（5 个）：project-type、test-e2e、ci、e2e-setup、e2e-verify。

### Task Decomposition: Recipe Templates

将 45 个 recipe 变体分解为 4 个实现任务，每个任务产出 `init-justfile` 命令文件中的一段字符串字面量：

| Task | 内容 | Recipe 数量 | 关键特征 | 依赖 |
|------|------|------------|---------|------|
| **A: backend template** | 纯后端项目 15 个 recipe 字面量 | 15 | 无 scope 参数；直接调用 Go/Rust/Python 工具链 | 无 |
| **B: frontend template** | 纯前端项目 15 个 recipe 字面量 | 15 | 无 scope 参数；直接调用 npm 工具链 | 无 |
| **C: mixed template** | 混合项目 15 个 recipe 字面量 | 15（其中 10 个含 bash case scope 分发） | 10 个 recipe 含 `scope=""` 参数 + bash case；5 个不含 scope | 无 |
| **D: project-type detection** | Model 4 探测逻辑 + project-type recipe 3 变体 | 1 recipe × 3 变体 | 检测 `package.json`/`go.mod`/`Cargo.toml`/`pyproject.toml`，选择 A/B/C 模板拼接 | A、B、C 完成 |

每个 Task 对应 init-justfile 命令文件中一个可独立审查的模板块，Task D 拼接对应模板 + project-type recipe 输出完整 justfile。

## Error Handling

### Error Types & Codes

| 错误码 | 来源 | 描述 | 处理方 |
|--------|------|------|--------|
| exit 1 + stderr | justfile recipe | 非法 scope 参数（如 `just build foo`） | skill: 回退到 `just VERB` |
| exit != 0 | `just project-type` | recipe 缺失或执行失败 | skill: 回退到 `just VERB`（无 scope） |
| 意外输出 | `just project-type` | 返回非 `frontend/backend/mixed` | skill: 视为 error，回退 |
| exit != 0 | `just compile` / `just build` | 编译/构建失败 | agent: 进入 error-fixer 流程 |

### Propagation Strategy

```
justfile recipe error (exit != 0)
  → skill/agent 检测退出码
  → scope 相关错误：回退到无 scope 命令，记录警告
  → 编译/测试错误：进入 error-fixer 流程
  → skill 不重试，不静默忽略错误
```

## Cross-Layer Data Map

| 字段 | task-cli (Go) | index.json (JSON) | skill (prompt text) | justfile (recipe) |
|------|--------------|-------------------|---------------------|-------------------|
| `scope` | `Task.Scope string` | `tasks[id].scope` | scope 解析协议 | recipe 位置参数 `scope=""` |
| `project-type` | — | — | `just project-type` 输出 | `@echo "mixed"` |

## Testing Strategy

### Per-Layer Test Plan

| 层 | 测试类型 | 工具 | 测试内容 | 覆盖目标 |
|----|---------|------|---------|---------|
| task-cli | 单元测试 | `go test` | `Task.Scope` JSON 序列化/反序列化 | 80% |
| init-justfile | e2e 测试 | `node:test` | 3 种项目类型生成正确的 justfile 变体 | 100% |
| justfile recipe | e2e 测试 | `node:test` + `just` | 15 个命令退出码验证、scope 参数验证 | 15/15 命令 |
| skill 迁移 | e2e 测试 | `node:test` | skill 文件内容检查（无原始命令） | 14/14 文件 |

### Key Test Scenarios

1. **init-justfile 自适应生成**：给定纯前端/纯后端/混合项目结构，生成正确变体
2. **scope 参数验证**：`just build frontend` 成功、`just build foo` 退出 1、空 scope 全量执行
3. **project-type 输出**：三种项目类型各返回正确单词
4. **scope 与 project-type 不匹配**：skill 层警告 + 回退
5. **skill 迁移完整性**：所有 skill/agent/command 文件无原始命令

### Overall Coverage Target

e2e 覆盖所有 15 个命令 + 3 种项目类型 + scope 解析流程 = 100% PRD 覆盖。

## Security Considerations

### Threat Model

| 威胁 | 风险 |
|------|------|
| scope 参数注入 | 低：scope 仅作为 bash case 分支条件，不拼接执行命令 |
| init-justfile 覆盖 | 中：已有 justfile 被覆盖导致自定义 recipe 丢失 |

### Mitigations

- scope 参数通过 `case` 严格匹配，`*)` 分支直接报错退出
- 所有 recipe 使用 `set -euo pipefail` 确保错误不静默传播
- **init-justfile 覆盖保护**（三层机制）：
  1. **边界标记合并**：init-justfile 使用 `# --- forge standard recipes ---` / `# --- end forge standard recipes ---` 边界标记。检测到已有 justfile 时，仅替换边界标记内的内容，保留用户自定义 recipe
  2. **`--force` 标志**：agent 调用 `init-justfile --force` 跳过确认，直接执行合并（不使用交互式 prompt，符合 PRD "所有 recipe 不使用交互式输入" 要求）
  3. **交互式确认（仅人类用户）**：无 `--force` 标志且边界标记不存在（首次写入）时，向人类用户提示确认；用户拒绝则中止，不覆盖

## PRD Coverage Map

| PRD Requirement / AC | Design Component | Interface / Model |
|----------------------|------------------|-------------------|
| Story 1 AC: skill 通过 `just <verb>` 执行 | justfile recipe contract | Interface 1 |
| Story 1 AC: 纯后端 `just test` → `go test` | init-justfile Go template | Interface 1 |
| Story 1 AC: 混合 `just build frontend` 仅前端 | bash case dispatch | Interface 1 |
| Story 2 AC: `package.json` → frontend justfile | ProjectDetection model | Model 4 |
| Story 2 AC: `go.mod` → backend justfile | ProjectDetection model | Model 4 |
| Story 2 AC: 两者都有 → mixed justfile | ProjectDetection model | Model 4 |
| Story 3 AC: index.json 每个 task 有 scope | Task struct extension | Interface 3, Model 1 |
| Story 3 AC: 前端文件 → scope=frontend | breakdown-tasks prompt | Interface 3 |
| Story 3 AC: 非混合项目 scope=all | scope default="all" | Model 3 |
| Story 4 AC: 退出码 0/非 0 | recipe error handling | Error Handling |
| Story 4 AC: 错误到 stderr | bash stderr redirect | Interface 1 |
| Story 4 AC: 幂等 install/e2e-setup | recipe 幂等设计 | Interface 1 |
| Story 5 AC: scope 不匹配 → 警告 + 回退 | scope resolution protocol | Interface 4 |
| PRD 5.4: 迁移 14 个文件 | skill migration list | N/A (文本替换) |
| PRD 5.5: task-cli scope 字段 | Task/TaskState struct | Model 1, Model 2 |

## Open Questions

- [x] scope 分发方式：bash case（已决定）
- [ ] init-justfile 中混合项目的前端/后端目录路径如何配置（固定约定 vs 用户指定）
- [ ] `just run` 在 e2e 测试中如何管理服务生命周期（后台启动 + PID 追踪 + teardown）

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| 子 recipe + 伞 recipe | 每个 scope 独立 recipe，更清晰 | recipe 数量翻倍（30+），维护成本高 | 混合项目 recipe 总量过多 |
| 固定目录约定（web/ + task-cli/） | 零配置 | 不通用，其他项目目录名不同 | 硬编码目录名不可复用 |

### References

- Proposal: `docs/proposals/justfile-standard-vocabulary/proposal.md` (scored 90/100)
- PRD Spec: `docs/features/justfile-standard-vocabulary/prd/prd-spec.md` (scored 92/100)
- Existing e2e tests: `tests/e2e/justfile-e2e-integration/cli.spec.ts` (20/20 passing)
- init-justfile command: `plugins/forge/commands/init-justfile.md`
- Go Task types: `task-cli/pkg/task/types.go`
