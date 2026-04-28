# E2E 生命周期纳入任务链（Phase 3）

## 问题

### 当前状态

Phase 1 和 Phase 2 完成后，任务链末尾是：

```
T-test-1: gen-test-cases
T-test-2: gen-test-scripts
```

e2e 测试执行和脚本毕业不在任务链里：

- **e2e 执行**：由 `all-completed` Stop hook 自动触发
- **脚本毕业**：`run-tasks` 建议用户手动运行 `/graduate-tests`

### 问题 1：e2e 执行在 hook 里，没有执行记录

hook 触发的 e2e 测试没有 task record，无法追溯"哪次运行通过了、失败了什么"。失败时的修复任务（`task add` 动态追加）有记录，但触发它的测试运行本身没有。

### 问题 2：毕业没有强制保证

`run-tasks` 只是建议用户手动运行 `/graduate-tests`，容易被遗忘。没有任务就没有 AC，没有 AC 就没有完成标准。

### 问题 3：all_completed.go 职责混乱

当前 hook 承担了两类职责：

```
职责 A（project-level，属于 hook）：
  - 检查所有任务完成
  - 运行 project-wide unit/integration tests
  - 运行 e2e regression（tests/e2e/）

职责 B（feature-level，不属于 hook）：
  - 运行 feature e2e 测试（testing/scripts/）
  - 分析失败、保存 raw output
  - block Stop hook，指示 agent 追加 fix 任务
```

职责 B 在 hook 里是因为之前没有任务链来承载它。

---

## 设计原则

| 原则 | 说明 |
|------|------|
| hook = project 不变式守卫 | 只跑 unit/integration tests + regression，快速可靠 |
| 任务链 = feature 工作流 | e2e 执行和毕业是 feature-level 工作，属于任务链 |
| 失败响应需要 agent | e2e 失败分析和修复需要理解力，不能在 hook 里硬编码 |
| regression 属于 hook | regression 是跨 feature 的 project-level 检查 |

---

## 两类 e2e 的本质区别

```
testing/scripts/*.spec.ts   ← feature e2e（开发中，T-test-3 的职责）
tests/e2e/**/*.spec.ts      ← regression（已毕业，稳定套件）
```

| 维度 | feature e2e | regression |
|------|-------------|------------|
| 内容 | 本次 feature 的新测试 | 所有已毕业 feature 的测试 |
| 状态 | 开发中，可能不稳定 | 已验证通过，相对稳定 |
| 失败含义 | 新功能有 bug | 改动破坏了已有功能 |
| 应该由谁跑 | T-test-3（任务） | all-completed hook |

---

## 方案

### 任务链扩展

在 `breakdown-tasks` 末尾追加两个标准任务：

```
T-test-1: gen-test-cases      (depends: last gate)
T-test-2: gen-test-scripts    (depends: T-test-1)
T-test-3: run-e2e-tests       (depends: T-test-2)   ← 新增
T-test-4: graduate-tests      (depends: T-test-3)   ← 新增
```

### T-test-3：run-e2e-tests

**职责**：运行 `testing/scripts/` 下的 feature e2e 测试，生成结果报告。

**执行**：调用 `/run-e2e-tests` skill。

**失败处理**：task-executor 分析 raw output，对每个失败：

```bash
task add --title "Fix: <description>" --priority P0 --breaking \
         --description "<root cause context>"
```

T-test-3 标记 `completed`（它完成了"运行测试并生成报告"的职责）。fix 任务（P0）优先级高于 T-test-4（P1），会在 T-test-4 之前被 claim。

**通过**：`testing/results/latest.md` 写入 PASS，T-test-4 可 claim。

### T-test-4：graduate-tests

**职责**：验证 e2e 通过，将脚本迁移到 `tests/e2e/` 回归套件。

**前置条件检查**：

```
1. 读 testing/results/latest.md
2. status = PASS → 执行 /graduate-tests → completed
3. status = FAIL → 标记 blocked，输出：
   "e2e tests are still failing. Wait for fix tasks to complete,
    then: task status T-test-4 pending"
```

**执行**：调用 `/graduate-tests` skill（agent 驱动的智能归类，Phase 2 实现）。

**不负责 regression**：regression 由 hook 负责。

### all_completed.go 简化

移除 feature e2e 执行逻辑，保留 project-level 检查：

```go
func runAllCompleted(...) {
    result := checkAllCompleted(...)
    if result == nil { os.Exit(0) }

    // 1. Unit + integration tests（project-wide）
    runProjectTests(result.ProjectRoot, result.TestCommand)

    // 2. E2E regression（tests/e2e/，已毕业套件）
    if hasJustfile(...) && hasJustRecipe(..., "test-e2e") {
        output, ok := runCmdCapture(..., "just", "test-e2e")
        if !ok {
            writeRawOutput(projectRoot, featureSlug, output)
            printHookJSON(map[string]any{
                "decision": "block",
                "reason":   "regression failed. Read testing/results/raw-output.txt, analyze failures, task add fix tasks.",
            })
            os.Exit(0)
        }
    }
}
```

**regression 失败响应**：hook block → agent 分析 → `task add` fix 任务 → fix 完成 → hook 再次触发 → regression 重跑。循环由 hook 维持。

**向后兼容**：没有 T-test-3/T-test-4 的旧项目，hook 不再自动跑 feature e2e。打印 WARNING 引导迁移：

```
WARNING: feature e2e scripts exist but haven't been run or graduated.
Add T-test-3 (run-e2e-tests) and T-test-4 (graduate-tests) to your task index,
or run /run-e2e-tests and /graduate-tests manually.
```

---

## 职责边界（最终）

```
T-test-3（任务）
  = "feature e2e 测试通过了吗？"
  → 运行 testing/scripts/，生成 latest.md
  → 失败：task add fix 任务（P0, breaking）
  → 通过：completed

T-test-4（任务）
  = "测试脚本归入回归套件了吗？"
  → 检查 latest.md PASS
  → /graduate-tests（智能归类迁移）
  → completed

all_completed hook
  = "project 整体健康吗？"
  → unit + integration tests
  → regression（tests/e2e/）
  → 失败：block → agent task add
```

---

## 涉及的文件变更

### 新增

| 文件 | 说明 |
|------|------|
| `breakdown-tasks/templates/run-e2e-tests.md` | T-test-3 任务模板 |
| `breakdown-tasks/templates/graduate-tests.md` | T-test-4 任务模板 |

### 修改

| 文件 | 变更 |
|------|------|
| `breakdown-tasks/SKILL.md` | Step 4d 追加 T-test-3、T-test-4；Output Checklist 更新 |
| `breakdown-tasks/templates/index.json` | 追加 T-test-3、T-test-4 示例条目 |
| `task-cli/internal/cmd/all_completed.go` | 移除 feature e2e 执行逻辑；regression 失败改为 block hook；添加向后兼容 WARNING |
| `task-cli/internal/cmd/test_results.go` | `writeRawOutput`/`writeLatestMd` 保留（regression 失败时 hook 仍需写入） |
| `task-cli/internal/cmd/all_completed_test.go` | 更新测试：移除 feature e2e 路径，新增 regression block 路径 |
| `task-cli/docs/OVERVIEW.md` | 更新 all-completed 行为描述 |
| `task-cli/docs/WORKFLOW.md` | 更新 All-Completed Workflow 图 |
| `plugins/forge/commands/run-tasks.md` | 移除 post-completion 建议（已由任务链覆盖） |

### 删除

| 代码 | 原因 |
|------|------|
| feature e2e 执行 switch（all_completed.go ~20 行） | 移到 T-test-3 |
| `runSpecsIndividually` 函数 | 仅用于 feature e2e |
| `countTestLines` 函数 | 仅用于 feature e2e stats |
| `writeRawOutput`/`writeLatestMd` 调用（feature e2e 路径） | 移到 `/run-e2e-tests` skill |
| block hook JSON（feature e2e 失败路径） | 移到 task-executor |

---

## 测试策略

### 新增测试

| 文件 | 测试用例 |
|------|----------|
| `all_completed_test.go` | `TestAllCompleted_RegressionFail_BlocksHook`、`TestAllCompleted_NoFeatureE2E_Warning` |

### 需要更新的测试

| 文件 | 变更 |
|------|------|
| `all_completed_test.go` | 移除 feature e2e 执行路径的测试 |
| `integration_test.go` | 移除 `runSpecsIndividually` 相关测试 |

---

## 实施顺序

1. `breakdown-tasks/templates/run-e2e-tests.md` — T-test-3 模板
2. `breakdown-tasks/templates/graduate-tests.md` — T-test-4 模板
3. `breakdown-tasks/templates/index.json` — 追加 T-test-3、T-test-4
4. `breakdown-tasks/SKILL.md` — 更新 Step 4d 和 Output Checklist
5. `all_completed.go` — 移除 feature e2e，regression 改为 block hook，添加 WARNING
6. `all_completed_test.go` / `integration_test.go` — 更新测试
7. 全量测试 + lint
8. `docs/OVERVIEW.md` + `WORKFLOW.md`
9. `plugins/forge/commands/run-tasks.md`
