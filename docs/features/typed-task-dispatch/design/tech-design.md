---
created: 2026-05-11
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: typed-task-dispatch

## Overview

将策略逻辑从 agent markdown 下沉到 task-cli，通过新增 `task prompt <id>` 命令合成类型专属 agent prompt。task-executor 变为薄执行器（只保留执行约束层）。

**两层模型**：

| 层 | 位置 | 内容 | 生效时机 |
|----|------|------|---------|
| 执行约束层 | `task-executor.md`（agent 定义） | ONE TASK、record-task 强制、无后台任务、最多 3 次 subagent、STOP 规则 | 始终生效 |
| 任务策略层 | `task prompt <id>` stdout（作为 Agent prompt 参数传入） | 类型专属执行步骤（TDD / 文档生成 / fix 诊断等） | 随 type 变化 |

两层不合并：约束层内嵌在 agent 定义里，策略层通过 prompt 参数注入。

## Architecture

### Layer Placement

纯 CLI + agent markdown 改造，无 UI、无数据库、无网络调用。

- **task-cli**（Go binary）：新增 `pkg/prompt` 包 + 3 个新命令
- **plugins/forge/agents/**：task-executor.md 精简
- **plugins/forge/commands/**：run-tasks.md、execute-task.md 路由更新
- **plugins/forge/skills/**：breakdown-tasks、quick-tasks 生成逻辑更新
- **plugins/forge/skills/breakdown-tasks/templates/**：任务模板 frontmatter 更新

### Component Diagram

```
run-tasks / execute-task
    │
    ├─ task claim ──────────────────────────────────────────────────────────┐
    │   └─ outputs: TASK_ID, TYPE, SCOPE, FEATURE, MAIN_SESSION             │
    │                                                                        │
    ├─ task prompt <id>  ◄──── new CLI command                              │
    │   ├─ read .forge/state.json → feature slug                            │
    │   ├─ read index.json → task.Type, task.Scope, task.Phase              │
    │   ├─ PhaseDetect() → phaseSummaryPath                                 │
    │   ├─ select embed template by type (pkg/prompt/data/*.md)             │
    │   └─ string substitution → stdout (pure text, no prefix/suffix)       │
    │                                                                        │
    ├─ if type == eval-cases: execute prompt in main session                 │
    └─ else: Agent(forge:task-executor, prompt=<stdout>)                    │
                │                                                            │
                └─ task-executor.md (constraint layer only, ~40 lines)  ────┘

task migrate  ◄──── new CLI command (one-time migration)
    └─ InferType(task.ID) → fills task.Type for all tasks

task validate  ◄──── extended: validates task.Type field
```

### Dependencies

変更なし（task-cli は既に `github.com/spf13/cobra` のみ依存）。新規外部依存なし。

## Interfaces

### Interface 1: `pkg/prompt.Synthesize`

```go
// SynthesizeOpts holds inputs for prompt synthesis.
type SynthesizeOpts struct {
    ProjectRoot     string // absolute path to project root
    FeatureSlug     string // e.g. "auth-refresh"
    TaskID          string // e.g. "2.1"
    FixRecordMissed bool   // true → use fix-record-missed template
}

// Synthesize returns the synthesized agent prompt for the given task.
// On success: returns non-empty string, nil error.
// On failure: returns empty string, non-nil error (caller writes to stderr + exits 1).
func Synthesize(opts SynthesizeOpts) (string, error)
```

### Interface 2: `pkg/prompt.InferType`

```go
// InferType infers the task type from the task ID using the migration rules.
// Always returns a non-empty string (falls back to TypeImplementation).
func InferType(id string) string
```

### Interface 3: `task prompt` CLI contract

```
Input:  task ID (positional arg), --fix-record-missed flag
Output: synthesized prompt → stdout (UTF-8, no prefix/suffix markers)
Errors: → stderr only; stdout empty on error
Exit:   0 on success, 1 on any error
Timing: < 500ms (local file reads + string substitution, no network)
```

### Interface 4: `task migrate` CLI contract

```
Input:  none (reads current feature's index.json)
Output: summary line → stdout ("Migrated N tasks. Run task validate to verify.")
Errors: → stderr; index.json unchanged on error
Exit:   0 on success, 1 if in_progress tasks exist or file I/O fails
```

## Data Models

### Task struct additions (`pkg/task/types.go`)

```go
// Type is the task execution type (e.g. "implementation", "fix", "gate").
// Required for all tasks after migration; validated by task validate.
// omitempty allows existing index.json files to load without error.
Type string `json:"type,omitempty"`

// BlockedReason records why a task entered blocked state.
// Written by run-tasks when task prompt exits non-zero.
BlockedReason string `json:"blockedReason,omitempty"`
```

### TaskState struct addition (`pkg/task/types.go`)

```go
// Type mirrors Task.Type for the claimed task (same pattern as MainSession, NoTest).
Type string `json:"type,omitempty"`
```

### Type enum constants (new block in `pkg/task/types.go`)

```go
const (
    TypeImplementation               = "implementation"
    TypeDocGenerationSummary         = "doc-generation.summary"
    TypeDocGenerationConsolidate     = "doc-generation.consolidate"
    TypeTestPipelineGenCases         = "test-pipeline.gen-cases"
    TypeTestPipelineEvalCases        = "test-pipeline.eval-cases"
    TypeTestPipelineGenScripts       = "test-pipeline.gen-scripts"
    TypeTestPipelineRun              = "test-pipeline.run"
    TypeTestPipelineGraduate         = "test-pipeline.graduate"
    TypeTestPipelineVerifyRegression = "test-pipeline.verify-regression"
    TypeFix                          = "fix"
    TypeGate                         = "gate"
)

// ValidTypes is the complete set of valid task type values.
var ValidTypes = map[string]bool{
    TypeImplementation:               true,
    TypeDocGenerationSummary:         true,
    TypeDocGenerationConsolidate:     true,
    TypeTestPipelineGenCases:         true,
    TypeTestPipelineEvalCases:        true,
    TypeTestPipelineGenScripts:       true,
    TypeTestPipelineRun:              true,
    TypeTestPipelineGraduate:         true,
    TypeTestPipelineVerifyRegression: true,
    TypeFix:                          true,
    TypeGate:                         true,
}
```

### index.schema.json additions

在 task properties 中新增：

```json
"type": {
  "type": "string",
  "enum": [
    "implementation",
    "doc-generation.summary",
    "doc-generation.consolidate",
    "test-pipeline.gen-cases",
    "test-pipeline.eval-cases",
    "test-pipeline.gen-scripts",
    "test-pipeline.run",
    "test-pipeline.graduate",
    "test-pipeline.verify-regression",
    "fix",
    "gate"
  ],
  "description": "Task execution type. Required for all tasks."
},
"blockedReason": {
  "type": "string",
  "description": "Why this task entered blocked state. Written by run-tasks when task prompt fails."
}
```

`noTest` 字段保留在 schema 中（向后兼容），但标记为 deprecated，不再由 breakdown-tasks 生成。

## New Package: `pkg/prompt`

```
task-cli/pkg/prompt/
├── prompt.go          # Synthesize(), PhaseDetect(), InferType(), typeToTemplate
├── prompt_test.go
└── data/              # Go embed (//go:embed data/*.md)
    ├── implementation.md
    ├── doc-generation-summary.md
    ├── doc-generation-consolidate.md
    ├── test-pipeline-gen-cases.md
    ├── test-pipeline-eval-cases.md
    ├── test-pipeline-gen-scripts.md
    ├── test-pipeline-run.md
    ├── test-pipeline-graduate.md
    ├── test-pipeline-verify-regression.md
    ├── fix.md
    ├── gate.md
    └── fix-record-missed.md
```

**Type → template filename mapping**（dots replaced with dashes）：

```go
var typeToTemplate = map[string]string{
    TypeImplementation:               "implementation.md",
    TypeDocGenerationSummary:         "doc-generation-summary.md",
    TypeDocGenerationConsolidate:     "doc-generation-consolidate.md",
    TypeTestPipelineGenCases:         "test-pipeline-gen-cases.md",
    TypeTestPipelineEvalCases:        "test-pipeline-eval-cases.md",
    TypeTestPipelineGenScripts:       "test-pipeline-gen-scripts.md",
    TypeTestPipelineRun:              "test-pipeline-run.md",
    TypeTestPipelineGraduate:         "test-pipeline-graduate.md",
    TypeTestPipelineVerifyRegression: "test-pipeline-verify-regression.md",
    TypeFix:                          "fix.md",
    TypeGate:                         "gate.md",
}
// fix-record-missed uses a separate key, not in ValidTypes
const templateFixRecordMissed = "fix-record-missed.md"
```

**Placeholder format**（与现有 `pkg/template/data/fix-task.md` 保持一致，使用 `{{PLACEHOLDER}}`）：

| Placeholder | 值来源 |
|-------------|--------|
| `{{TASK_ID}}` | task.ID |
| `{{SCOPE}}` | task.Scope（默认 "all"） |
| `{{FEATURE_SLUG}}` | feature slug |
| `{{TASK_FILE}}` | 任务定义文件绝对路径 |
| `{{RECORD_FILE}}` | record 文件绝对路径 |
| `{{PHASE_SUMMARY_PATH}}` | 前一 phase summary 路径（空字符串表示不注入） |

### Phase Boundary Detection

```
PhaseDetect(projectRoot, featureSlug, taskID string) string

算法：
1. 加载 index.json 中所有任务
2. maxCompletedPhase = max{ phase(t.ID) | t.Status == "completed" AND isBusinessTask(t.ID) }
   - phase(id) = id 第一段的整数（"2.1" → 2），非整数前缀返回 -1
   - isBusinessTask: 不以 ".gate" 或 ".summary" 结尾，不以 "T-" 开头
3. currentPhase = phase(taskID)
4. if currentPhase > maxCompletedPhase AND currentPhase > 1:
   - path = "docs/features/<slug>/tasks/records/<currentPhase-1>-summary.md"
   - if file exists: return path
   - if file missing: return ""（非致命，模板处理空值）
5. else: return ""
```

## New Commands

### `task prompt` (`internal/cmd/prompt.go`)

```
task prompt <id> [--fix-record-missed]

执行流程：
1. RequireFeature() → featureSlug
2. LoadIndex() → 按 ID 查找任务
3. if task.Type == "" → stderr "type field missing for task <id>" → exit 1
4. if !ValidTypes[task.Type] → stderr "unknown type: <type>" → exit 1
5. PhaseDetect() → phaseSummaryPath
6. Synthesize() → prompt string
7. fmt.Print(prompt) → exit 0

任何步骤失败：fmt.Fprintf(os.Stderr, ...) → exit 1，stdout 为空
```

### `task migrate` (`internal/cmd/migrate.go`)

```
task migrate

执行流程：
1. RequireFeature() → featureSlug
2. LoadIndex()
3. if any task.Status == "in_progress" → stderr "complete or manually mark in-progress tasks first" → exit 1
4. for each task: task.Type = InferType(task.ID)
5. SaveIndex()
6. stdout: "Migrated N tasks. Run task validate to verify."

InferType 推断规则（按优先级）：
  strings.HasSuffix(id, ".summary")    → TypeDocGenerationSummary
  strings.HasSuffix(id, ".gate")       → TypeGate
  id == "T-test-1"                     → TypeTestPipelineGenCases
  id == "T-test-1b"                    → TypeTestPipelineEvalCases
  id == "T-test-2"                     → TypeTestPipelineGenScripts
  id == "T-test-3"                     → TypeTestPipelineRun
  id == "T-test-4"                     → TypeTestPipelineGraduate
  id == "T-test-4.5"                   → TypeTestPipelineVerifyRegression
  id == "T-test-5"                     → TypeDocGenerationConsolidate
  strings.HasPrefix(id, "fix-") ||
  strings.HasPrefix(id, "disc-")       → TypeFix
  default                              → TypeImplementation
```

### `task validate` extension (`internal/cmd/validate.go`)

在 `validateTasks()` 中新增：

```go
// Type validation
if t.Type == "" {
    v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'type' field", key))
} else if !task.ValidTypes[t.Type] {
    v.errors = append(v.errors, fmt.Sprintf("Task '%s': invalid type '%s'", key, t.Type))
}
// mainSession consistency check
if t.MainSession && t.Type != task.TypeTestPipelineEvalCases {
    v.warnings = append(v.warnings,
        fmt.Sprintf("Task '%s': mainSession=true but type is '%s' (only eval-cases should use mainSession)", key, t.Type))
}
```

### `task status <id> pending`（blocked 状态恢复）

不新增 `task unblock` 命令。使用现有 `task status <id> pending` 命令重置 blocked 任务状态。blocked 状态恢复流程：

```
1. 读取 task list 中 blocked 任务的 blockedReason 字段定位根因
2. 修复模板文件或 type 注册
3. task status <id> pending  ← 现有命令，重置状态
4. 重新运行 run-tasks
```

## Agent Changes

### `task-executor.md` — 精简为约束层

移除所有策略块（TDD 步骤、文档生成步骤、quality gate 步骤、NO_TEST 分支）。保留：

```markdown
---
name: task-executor
description: "Thin executor: follow the steps in your prompt. Hard constraints always active."
model: sonnet
color: green
memory: project
---

## Hard Constraints

<EXTREMELY-IMPORTANT>
1. ONE TASK PER INVOCATION — after completing, STOP immediately, no exceptions
2. record-task IS MANDATORY — task is NOT done without it
3. NO BACKGROUND TASKS — all commands run synchronously
4. Maximum 3 subagent calls per task
5. FORBIDDEN: run "task claim", read index.json, or start any subsequent task
</EXTREMELY-IMPORTANT>

Execute the task described in your prompt. The prompt contains all steps and context.
Call forge:record-task when done. Then STOP.
```

inputs 字段移除（prompt 参数携带所有上下文）。预计从 259 行缩减至 ~40 行。

### `run-tasks.md` — 路由更新

**Step 2 dispatch 替换**：

```
旧：
  Agent(forge:task-executor, prompt="TASK_KEY: ... TASK_FILE: ... NO_TEST: ...")

新：
  1. 运行: task prompt <TASK_ID>
     - exit != 0 → 将 stderr 写入 blockedReason，task status <KEY> blocked，continue loop
     - type == test-pipeline.eval-cases → 在主会话中直接按 prompt 执行（不 dispatch）
     - 其他 → Agent(forge:task-executor, prompt=<stdout>)

  2. record 缺失恢复：
     旧：Agent(forge:error-fixer, ...)
     新：task prompt <TASK_ID> --fix-record-missed → Agent(forge:task-executor, prompt=<stdout>)
```

移除 error-fixer dispatch 及相关 Error Handling 条目。

**claim output 新增字段**：run-tasks 从 claim 输出中提取 `TYPE` 字段（`task claim` 输出新增 `TYPE` 行，与 `MAIN_SESSION`、`NO_TEST` 同格式）。

### `execute-task.md` — 路由同步

与 run-tasks 保持相同路由逻辑：
- 调用 `task prompt <id>` 合成 prompt
- eval-cases 在主会话执行
- record 缺失时调用 `task prompt <id> --fix-record-missed`
- 不再使用 TASK_FILE + NO_TEST 参数组合

## Skill Layer Changes

### `breakdown-tasks` SKILL.md

在 Step 3（Create Task Files）中新增规则：

```
Type Assignment（每个任务必须设置 type 字段）：

| 任务类型 | type 值 |
|---------|---------|
| 业务实现任务（implementation, interface, model, error handling） | "implementation" |
| Phase summary 任务（ID 以 .summary 结尾） | "doc-generation.summary" |
| Gate 任务（ID 以 .gate 结尾） | "gate" |
| T-test-1 | "test-pipeline.gen-cases" |
| T-test-1b | "test-pipeline.eval-cases" |
| T-test-2 | "test-pipeline.gen-scripts" |
| T-test-3 | "test-pipeline.run" |
| T-test-4 | "test-pipeline.graduate" |
| T-test-4.5 | "test-pipeline.verify-regression" |
| T-test-5 | "doc-generation.consolidate" |
| fix-task（ID 以 fix- 或 disc- 开头） | "fix" |

无法匹配任何规则时：type = "implementation"，stderr 输出警告（含任务 ID 和原因），不中断生成流程。
```

### `quick-tasks` SKILL.md

在 Step 3（Create Task Files）中新增相同的 Type Assignment 规则（quick mode 任务 ID 为简单整数，均映射到 "implementation"；T-quick-* 任务按 T-test-* 规则映射）。

### 任务模板 frontmatter 更新

`breakdown-tasks/templates/` 下所有模板文件：

| 模板文件 | 变更 |
|---------|------|
| `task.md` | 新增 `type: "{{TYPE}}"` 字段，移除 `noTest: false` |
| `gate-task.md` | 新增 `type: "gate"` 固定值，移除 `noTest: false` |
| `phase-summary-task.md` | 新增 `type: "doc-generation.summary"` 固定值，移除 `noTest: true` |
| `gen-test-cases.md` | 新增 `type: "test-pipeline.gen-cases"` 固定值 |
| 其他 T-test-* 模板 | 按类型添加对应固定 type 值 |

## Error Handling

| 失败场景 | 行为 |
|---------|------|
| `task prompt` exit 非零 | run-tasks 将 stderr 写入 `task.BlockedReason`，`task status <KEY> blocked`，continue loop |
| 模板文件缺失 | `pkg/prompt.Synthesize` 返回 error → stderr + exit 1 |
| 模板占位符格式错误 | `strings.ReplaceAll` 不报错（占位符原样保留）；通过单元测试在开发期捕获 |
| `.forge/state.json` 缺失 | `RequireFeature()` 降级到 git context → state.json → 单 feature 目录扫描 |
| task ID 不在 index.json 中 | stderr "task <id> not found" → exit 1 |
| `task migrate` 存在 in_progress 任务 | stderr 报错，index.json 不修改 |
| `task prompt --fix-record-missed` 失败 | run-tasks 将 stderr 写入 `task.BlockedReason`，`task status <KEY> blocked` |

## Cross-Layer Data Map

Single-layer feature（CLI + agent markdown 改造，无跨层数据流）。

唯一跨边界数据：`task prompt` stdout → `Agent(forge:task-executor)` prompt 参数。

| 数据 | task prompt 侧 | run-tasks 侧 | task-executor 侧 |
|------|---------------|-------------|-----------------|
| 合成 prompt | `string`（stdout） | 捕获为变量 | `prompt` 参数（string） |
| 错误信息 | `string`（stderr） | 写入 `blockedReason` | 不可见 |

## Testing Strategy

### `pkg/prompt` 单元测试（目标覆盖率 ≥ 80%）

| 测试 | 覆盖内容 |
|------|---------|
| `TestSynthesize_AllTypes` | Table-driven：11 种 type × 验证必要占位符已替换（无 `{{` 残留） |
| `TestSynthesize_FixRecordMissed` | `--fix-record-missed` 选择正确模板 |
| `TestSynthesize_TypeMissing` | task.Type 为空时返回 error |
| `TestSynthesize_UnknownType` | 未注册 type 返回 error |
| `TestPhaseDetect_NewPhase` | currentPhase > maxCompleted → 注入 summary 路径 |
| `TestPhaseDetect_SamePhase` | currentPhase == maxCompleted → 不注入 |
| `TestPhaseDetect_FirstPhase` | phase 1 任务 → 不注入 |
| `TestPhaseDetect_SummaryFileMissing` | summary 文件不存在 → 返回空字符串（非 error） |
| `TestInferType_AllRules` | Table-driven：所有 ID 模式 → 预期 type |

### `internal/cmd` 集成测试

| 测试 | 覆盖内容 |
|------|---------|
| `TestPromptCmd_OutputToStdout` | stdout 非空，stderr 为空，exit 0 |
| `TestPromptCmd_MissingType_ExitsNonZero` | exit 1，stdout 为空 |
| `TestPromptCmd_UnknownType_ExitsNonZero` | exit 1，stdout 为空 |
| `TestMigrateCmd_InProgress_Blocked` | 存在 in_progress 任务时报错，index.json 不变 |
| `TestMigrateCmd_AllTypes_Inferred` | 迁移后所有任务均有 type 字段，task validate 通过 |
| `TestValidateCmd_MissingType_Error` | type 缺失时 validate 报错 |
| `TestValidateCmd_InvalidType_Error` | 非法 type 值时 validate 报错 |

### Version Bump

`1.11.0` → `1.12.0`（minor：新增 `prompt`、`migrate` 命令；扩展 `validate`）

## PRD Coverage Map

| User Story | 设计元素 |
|------------|---------|
| Story 1：非编码任务获得正确流程 | task-executor.md 约束层；策略来自 `task prompt` 模板 |
| Story 2：新增类型只需一个模板文件 | `pkg/prompt/data/` embed；`ValidTypes` map 注册 |
| Story 3：`task prompt <id>` 独立检查 | `task prompt` 命令 → stdout |
| Story 4：`task migrate` 迁移旧 index.json | `task migrate` 命令 + `InferType()` |
| Story 5：breakdown-tasks 自动设置 type | breakdown-tasks / quick-tasks SKILL.md Type Assignment 规则 |
| Story 6：execute-task 路由与 run-tasks 一致 | execute-task.md 同步路由更新 |
| Story 7：error-fixer 废弃后等价覆盖 | run-tasks.md 移除 error-fixer dispatch；fix 模板承接全部能力 |
