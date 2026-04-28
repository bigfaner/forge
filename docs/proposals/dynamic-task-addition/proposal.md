# 动态任务追加 + 测试结果分析重构

## 问题

### 问题 1：无法动态追加任务

当前 task 系统假设"一次性规划完所有任务"。动态追加仅存在于 `all_completed.go` 的 `appendFixTask`——硬编码为 e2e 修复专用（`fix-e2e-<round>-<idx>`，最多 3 轮）。

实际开发中，agent 执行测试时会发现未预知的问题，需要动态追加修复任务。但没有任何 CLI 入口支持。

### 问题 2：CLI 测试结果解析不准确

`test_results.go` 用正则猜测测试输出，试图覆盖所有框架：

```go
// 匹配规则覆盖 node:test spec、go test、pytest、generic
// 但 node:test TAP 格式的 "not ok N - name" 完全匹配不到
// countPassingTests 误匹配非测试行
// extractErrorContext 位置依赖，经常截断或遗漏
```

`all_completed.go` 支持的测试运行器至少 6 种（just test / make test / go test / npm test / pytest / npx tsx），每种输出格式不同。**用正则覆盖所有框架是不可能的。**

### 根因

CLI 承担了需要智能的工作。违反"CLI 是工具，Agent 是大脑"原则。

---

## 设计原则

| 原则 | 说明 |
|------|------|
| CLI 是工具 | 校验 + 写入，不做决策 |
| Agent 是大脑 | 分析 + 决策 + 调用 CLI |
| Exit code 是唯一通用信号 | 跨所有测试框架，不依赖输出格式 |
| Raw output 交给 Agent | Agent 理解任意框架的输出，正则做不到 |

---

## 方案概览

两个改动是一个整体：

1. **`task add` 命令**：让 agent 能安全追加任务到 index.json
2. **移除 CLI 测试解析**：`all-completed` 只用 exit code 判断成败，raw output 交给 agent 分析

---

## 涉及的文件清单

### 新建

| 文件 | 说明 |
|------|------|
| `task-cli/pkg/task/add.go` | `AddTaskOpts`、`AddTask()`、`CreateTaskMarkdown()`、`generateDiscID()` |
| `task-cli/pkg/task/add_test.go` | 单元测试 |
| `task-cli/internal/cmd/add.go` | Cobra 命令：flags、`runAdd`、`executeAdd` |
| `task-cli/internal/cmd/add_test.go` | 命令集成测试 |

### 修改

| 文件 | 变更 |
|------|------|
| `task-cli/internal/cmd/root.go` | `init()` 注册 `addCmd` |
| `task-cli/internal/cmd/all_completed.go` | 移除 `appendFixTask`/`createFixTaskFile`/解析逻辑，简化为 exit code + raw output |
| `task-cli/internal/cmd/test_results.go` | 删除 `parseTestFailures`/`countPassingTests`/`extractErrorContext`/`extractRelevantOutput`/`matchTestCaseID`/`sanitizeTestName`/`writeFailureFiles`。保留 `TestFailure`/`TestStats` 类型和 `writeLatestMd`（简化） |
| `task-cli/internal/cmd/errors.go` | 新增 `ErrTaskIDConflict`、`ErrInvalidDependency` |
| `task-cli/docs/OVERVIEW.md` | 新增 `task add` 命令文档 |
| `task-cli/docs/WORKFLOW.md` | 新增 Dynamic Task Addition workflow + 更新 All-Completed workflow |
| `plugins/forge/agents/task-executor.md` | Error Handling 表新增 `task add` 路径 |
| `plugins/forge/commands/run-tasks.md` | Breaking Gate 新增 `task add` 选项 |

---

## 改动 1：`task add` 命令

### `pkg/task/add.go` — 共享逻辑

```go
type AddTaskOpts struct {
    ID            string   // 空 = 自动生成 disc-N
    Title         string   // 必填
    Priority      string   // 默认 P1
    EstimatedTime string
    Dependencies  []string
    Breaking      bool
    Description   string   // task markdown body
    Status        string   // 默认 pending
}
```

**函数：**

| 函数 | 职责 |
|------|------|
| `AddTask(indexPath string, opts AddTaskOpts) (string, error)` | LoadIndex → 校验 → 自动生成 ID → 插入 index.Tasks → SaveIndex → 返回 ID |
| `CreateTaskMarkdown(tasksDir, filename string, opts AddTaskOpts) error` | YAML frontmatter + markdown body → 写入文件 |
| `generateDiscID(index *TaskIndex) string` | gap-filling：扫描 `disc-*`，找最小空位 |

**校验规则：**
- ID 唯一性（与 index.Tasks 已有 key 不冲突）
- Title 非空
- Priority ∈ {P0, P1, P2}
- Dependencies 中的每个 ID 必须存在于 index.Tasks

**`disc-1` 与 claim 排序兼容性：**
`claim.go:271` 的 `compareVersionIDs` 按 `.` 分段。`disc-1` 整体为一段，`parseSegment` 归为 `default: 0, false`，与 `fix-e2e-*` 同排序。由于 fix 类 task 通常 P0，会被优先 claim。无需改 claim 排序。

### `internal/cmd/add.go` — CLI 命令

**Flags：**

| Flag | 必填 | 默认 | 说明 |
|------|------|------|------|
| `--title` | Yes | - | 任务标题 |
| `--id` | No | auto `disc-N` | 自定义 ID |
| `--priority` | No | P1 | P0/P1/P2 |
| `--depends-on` | No | none | 逗号分隔 |
| `--estimated-time` | No | - | 时间估算 |
| `--breaking` | No | false | 触发全量测试 |
| `--description` | No | - | 写入 markdown body |

**`executeAdd()` 流程：**

```
FindProjectRoot → RequireFeature → 校验 title → 解析 flags
  → task.AddTask(indexPath, opts) → ID
  → task.CreateTaskMarkdown(tasksDir, id+".md", opts, opts.Description)
  → feature.EnsureForgeState(projectRoot, featureSlug)  // 重置 allCompleted=false
  → PrintFields(ACTION:ADDED, KEY, ID, TITLE, PRIORITY, STATUS, FILE, RECORD, BREAKING, DEPENDENCIES)
```

**Forge state 重置：**
`task add` 添加 pending task 后调用 `EnsureForgeState` 设 `allCompleted=false`。确保：
- `task claim` 能继续领取新 task
- Stop hook 不会因旧 `allCompleted=true` 提前触发

### `internal/cmd/errors.go` — 新增错误

```go
func ErrTaskIDConflict(id string) *AIError        // CONFLICT
func ErrInvalidDependency(deps []string) *AIError  // VALIDATION_ERROR
```

---

## 改动 2：移除 CLI 测试解析，简化 `all-completed`

### 当前流程

```
all-completed:
  1. 运行 e2e 测试 → 捕获 output string
  2. parseTestFailures(output) → []TestFailure    ← 正则猜测，不准确
  3. countPassingTests(output) → int               ← 误匹配
  4. matchTestCaseID → 映射到 test case ID         ← 基于不准确的 test name
  5. writeLatestMd + writeFailureFiles             ← 基于不准确的解析
  6. appendFixTask(projectRoot, slug, failures)    ← 硬编码 fix-e2e-*
  7. printHookJSON({"decision":"block"})
```

### 改造后流程

```
all-completed:
  1. 运行 e2e 测试 → 捕获 output string + exit code
  2. exit code == 0?
     → YES: graduate + project tests（不变）
     → NO:
       a. 保存 raw output → testing/results/raw-output.txt
       b. writeLatestMd（简化版：status=FAIL + raw output 路径）
       c. printHookJSON({
            "decision": "block",
            "reason": "e2e tests failed. Read testing/results/raw-output.txt, analyze failures, then use task add to create fix tasks."
          })
  3. Agent 读 raw-output.txt → 分析 → task add → claim
```

### `all_completed.go` 具体变更

**删除的函数/逻辑：**

| 函数 | 行号 | 原因 |
|------|------|------|
| `parseTestFailures` 调用 | L149-155 | 正则猜测不准确 |
| `matchTestCaseID` 调用 | L153-155 | 依赖不准确的 test name |
| `TestStats` 计算 | L158-164 | 依赖不准确的 countPassingTests |
| `writeFailureFiles` 调用 | L168-170 | 依赖不准确的解析结果 |
| `appendFixTask` 全部 | L224-284 | 硬编码 fix-e2e，由 agent + task add 取代 |
| `createFixTaskFile` 全部 | L287-326 | appendFixTask 的子函数 |
| `errFixLimitExceeded` | L51 | 3 轮上限由 agent 决定 |

**`test_results.go` 删除的函数：**

| 函数 | 原因 |
|------|------|
| `parseTestFailures` | 正则猜测，由 agent 分析取代 |
| `extractErrorContext` | 位置依赖，不准确 |
| `extractRelevantOutput` | 截取不精确 |
| `countPassingTests` | 误匹配非测试行 |
| `matchTestCaseID` | 基于不准确的 test name |
| `sanitizeTestName` | matchTestCaseID 的辅助 |
| `writeFailureFiles` | 基于不准确的解析 |

**`test_results.go` 保留：**

| 内容 | 原因 |
|------|------|
| `TestFailure` struct | 类型定义仍有用（latest.md 引用） |
| `TestStats` struct | 类型定义 |
| `writeLatestMd` | 简化：只写 status + raw output 路径 |

**新增逻辑：**

```go
// 在 runAllCompleted 中，e2e 测试失败时：
if !e2eSuccess {
    // 保存 raw output
    resultsDir := filepath.Join(projectRoot, feature.GetFeatureTestingResultsDir(result.FeatureSlug))
    os.MkdirAll(resultsDir, 0755)
    rawPath := filepath.Join(resultsDir, "raw-output.txt")
    os.WriteFile(rawPath, []byte(e2eOutput), 0644)

    // 简化的 latest.md
    writeLatestMd(projectRoot, result.FeatureSlug, TestStats{Fail: 1}, nil)
    // latest.md 内容：status=FAIL + "See raw-output.txt for details"

    // 通知 agent
    printHookJSON(map[string]any{
        "decision": "block",
        "reason":   fmt.Sprintf("e2e tests failed. Read testing/results/raw-output.txt, analyze failures, then use `task add` to create fix tasks for each failure."),
    })
    os.Exit(0)
}
```

**`writeLatestMd` 简化版：**

```go
func writeLatestMd(projectRoot, featureSlug string, stats TestStats, _ []TestFailure) error {
    // 只写：status + raw output 路径 + timestamp
    // 不再写 per-failure 详情（那是 agent 的工作）
}
```

### 对 `run-e2e-tests` skill 的影响

**无影响。** `run-e2e-tests` 是独立的 skill，由 agent 直接调用。它自己执行测试、解析输出、生成报告。`all-completed` 是 Stop hook 中的自动流程，两者独立。

---

## Agent/Skill 提示词更新

### `task-executor.md`

Error Handling 表新增行：

```markdown
| Test failures beyond scope | Use `task add --title "Fix: ..." --priority P0 --breaking` to create fix task |
```

新增段落 "Dynamic Task Addition"：

```markdown
### Adding Tasks Dynamically

When discovering issues beyond the current task's scope (pre-existing bugs, environment issues,
unrelated module failures):

1. Run `task add --title "Fix: <concise description>" --priority P0 --breaking --description "Context"`
2. Continue current task — the new task will be picked up by next `task claim`
```

### `run-tasks.md`

Step 5 Breaking Gate 新增选项：

```markdown
If tests fail:
  Option A: Dispatch error-fixer (existing)
  Option B: `task add --title "Fix: <failure>" --priority P0 --breaking` → continue loop
```

### Stop hook 后的 agent 行为（隐含，不需改文件）

当 `task all-completed` 返回 `{"decision":"block", "reason":"..."}` 后，Claude 主 session：

1. 读 `testing/results/raw-output.txt`
2. 分析：哪个测试失败了？根因是什么？是代码 bug / 测试脚本问题 / 环境问题？
3. 对每个需要独立修复的问题：`task add --title "Fix: ..." --priority P0 --breaking --description "..."`
4. `task claim` → 开始修复

---

## 文档更新

### `task-cli/docs/OVERVIEW.md`

- Command Quick Reference 新增 `task add`
- 新增 "Dynamic Task Addition" 小节

### `task-cli/docs/WORKFLOW.md`

- 新增 "Dynamic Task Addition Workflow" 节
- 更新 "All-Completed Workflow"：移除 fix-e2e 追加，改为 raw output + agent 分析

---

## 测试策略

### 新增测试

| 文件 | 测试用例 |
|------|----------|
| `pkg/task/add_test.go` | `TestAddTask_Basic`, `TestAddTask_AutoID`, `TestAddTask_AutoID_Sequential`, `TestAddTask_AutoID_GapFill`, `TestAddTask_DuplicateID`, `TestAddTask_InvalidPriority`, `TestAddTask_EmptyTitle`, `TestAddTask_DependencyNotFound`, `TestCreateTaskMarkdown_Basic`, `TestCreateTaskMarkdown_WithBody` |
| `internal/cmd/add_test.go` | `TestExecuteAdd_Basic`, `TestExecuteAdd_NoTitle`, `TestExecuteAdd_ForgeStateReset` |

### 需要更新的测试

| 文件 | 变更 |
|------|------|
| `all_completed_test.go` | `TestAppendFixTask` 删除（函数已移除），替换为 `TestAllCompleted_WritesRawOutput` |

### 回归验证

```bash
go test -race -cover ./...
go vet ./...
golangci-lint run ./...
bash ../claude-code-go/scripts/lint-arch.sh
```

---

## 实施顺序

1. `pkg/task/add.go` + `add_test.go` — TDD
2. `internal/cmd/errors.go` — 新增错误构造函数
3. `internal/cmd/add.go` + `add_test.go` — TDD
4. `internal/cmd/root.go` — 注册 addCmd
5. `test_results.go` — 删除解析函数，简化 `writeLatestMd`
6. `all_completed.go` — 移除 `appendFixTask`/`createFixTaskFile`，改用 raw output 流程
7. `all_completed_test.go` — 更新测试
8. 全量测试 + lint
9. `docs/OVERVIEW.md` + `WORKFLOW.md`
10. `agents/task-executor.md` + `commands/run-tasks.md`

---

## Phase 2：Agent 驱动的测试脚本毕业

### 问题

当前 `graduateTestScripts` 在 `all-completed` hook 中自动执行，用正则从 `test-cases.md` 扫描 `**Target**` 行决定归类路径，然后硬拷贝 spec 文件到 `tests/e2e/<type>/<target>/`。

**缺陷：**

| 问题 | 说明 |
|------|------|
| 归类依赖 test-cases.md 格式 | gen-test-cases 没写 Target 行时，归类失败或跳过 |
| 一个 spec 只能归一个目录 | `ui.spec.ts` 可能测试多个页面，硬拷贝无法拆分 |
| Agent 不参与决策 | Hook 自动拷贝，没有机会让 agent 审查、合并、拆分、归类 |
| "unknown" placeholder | 解析失败时生成无用任务（见 lessons/tool-fix-e2e-unknown-placeholder.md） |

**根因：毕业是需要理解力的决策，不是文件操作。**

### 设计：从 hook 剥离毕业，交给 Agent

#### 毕业需要回答的问题（每一个都需要理解力）

1. **这个脚本测了什么？** → 读代码才能知道
2. **该归到哪个目录？** → 看现有 `tests/e2e/` 结构才能决定
3. **要不要拆分？** → 一个 spec 混合多个功能域怎么办
4. **现有 helpers.ts 是否兼容？** → 可能需要合并而非覆盖
5. **import 路径是否需要调整？** → sitemap 结构可能变了

#### 流程对比

**当前（CLI 自动）：**
```
hook → e2e pass → 正则读 test-cases.md → 硬拷贝 → marker → project tests
```

**改造后（Agent 驱动）：**
```
hook → e2e pass → 写 results (PASS) → project tests → regression → 完成

// 毕业由 agent 在适当时机调用 /graduate-tests skill
agent → 读 scripts/*.spec.ts → 分析内容 → 智能归类 → 迁移 → marker
```

#### all-completed hook 变更

从 `all_completed.go` 中移除 `graduateTestScripts()` 调用。e2e 测试通过后直接跑 project tests + regression，不再自动毕业。

#### 新增 `/graduate-tests` skill

**触发时机：**
- `run-tasks` orchestrator 在所有 task 完成后建议运行
- 或用户手动调用

**Skill 工作流：**

```
Step 1: 读源脚本
  读 testing/scripts/*.spec.ts + helpers.ts

Step 2: 分析现有结构
  读 tests/e2e/ 目录结构（如果存在），理解现有归类规范

Step 3: 智能归类决策
  - 分析每个 spec 的 test 描述和代码
  - 确定归类维度（按类型/功能域/页面路由）
  - 决定是否需要拆分（一个 spec 覆盖多个域）
  - 决定是否需要合并（多个 spec 测同一域）

Step 4: 执行迁移
  - 复制/拆分/合并 spec 文件到 tests/e2e/<category>/
  - 重写 import 路径（./helpers → ../helpers）
  - 合并或保留 helpers.ts（不覆盖已有的）
  - 如果 tests/e2e/package.json 不存在，创建 + npm install

Step 5: 创建 graduation marker
  写 tests/e2e/.graduated/<slug>（时间戳），防止重复毕业
```

**智能归类示例：**

```
# 输入：testing/scripts/
ui.spec.ts    # 包含 login, dashboard, profile 的测试
api.spec.ts   # 全是 auth 相关
cli.spec.ts   # 通用 CLI 命令测试

# Agent 分析后决定：
ui.spec.ts → 拆分为：
  tests/e2e/ui/login/login.spec.ts
  tests/e2e/ui/dashboard/dashboard.spec.ts
  tests/e2e/ui/profile/profile.spec.ts

api.spec.ts → 归入：
  tests/e2e/api/auth/auth.spec.ts

cli.spec.ts → 归入：
  tests/e2e/cli/cli.spec.ts（不拆分，因为都是通用命令）
```

#### 为什么不从 hook 中 block 让 agent 立即毕业

考虑过：`hook → e2e pass → block → agent graduates → continue`

**不采用：**
1. Hook block 是为"需要 agent 修复"设计的，不是为"需要 agent 做决策"
2. 毕业不紧急——可以在任何时间点做
3. 毕业涉及文件组织，可能需要用户确认归类方案
4. 解耦后，毕业可以独立重试（不需要重新跑 e2e）
5. 毕业是幂等的（有 marker 防重），"稍后做"和"立即做"效果相同

#### 涉及的文件变更

| 文件 | 变更 |
|------|------|
| `task-cli/internal/cmd/all_completed.go` | 删除 `graduateTestScripts` 调用 |
| `plugins/forge/skills/graduate-tests/SKILL.md` | 新增：agent 驱动的毕业 skill |
| `plugins/forge/commands/run-tasks.md` | 更新 Post-Completion，建议 `/graduate-tests` |
| `plugins/forge/SKILLS.md` | 注册 `/graduate-tests` |

#### CLI 层面不需要新命令

毕业是纯文件操作（读、写、移、改 import），agent 用 Bash/Write/Edit 完成。Graduation marker 直接 Write 创建，不需要 `task graduate` 命令。
