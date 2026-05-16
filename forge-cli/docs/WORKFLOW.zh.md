# forge CLI 关键流程

## 1. Feature 识别流程

```
┌─────────────────────────────────────────────────────────────────┐
│                   GetCurrentFeature()                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 获取 Git 上下文 │
                    │ (worktree/branch)│
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ Git 上下文存在   │             │ 无 Git 上下文    │
    │ 检查 feature    │             │ 扫描 process/   │
    │ 目录是否存在    │             │ 目录            │
    └────────┬────────┘             └────────┬────────┘
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ 存在: 返回      │             │ 有 task-state:  │
    │ 不存在: 创建    │             │ 返回该 feature  │
    │ 并返回          │             │                 │
    └─────────────────┘             └─────────────────┘

Feature 识别优先级:
1. Git Worktree 名称 (如: feature-auth-login)
2. Git 分支名称 (提取 feature/xxx 中的 xxx)
3. Features 目录中有 tasks/process/state.json 的 feature
4. Features 目录中唯一有 index.json 的 feature
```

### Git 分支 → Feature 映射

```
分支名称                    → Feature Slug
─────────────────────────────────────────────
feature/auth-login         → auth-login
feat/user-registration     → user-registration
fix/null-pointer           → null-pointer
bugfix/memory-leak         → memory-leak
hotfix/security-issue      → security-issue
chore/update-deps          → update-deps
main/master/HEAD           → (忽略，回退到目录扫描)
custom-branch              → custom-branch
```

---

## 2. 任务声明流程 (forge task claim)

```
┌─────────────────────────────────────────────────────────────────┐
│                     forge task claim                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 调用            │
                    │ GetCurrentFeature│
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 加载 task-state │
                    │ 检查进行中任务   │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ 有进行中任务     │             │ 无进行中任务     │
    │ 直接返回该任务   │             │ 搜索下一个任务   │
    └─────────────────┘             └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ 加载 index.json │
                                    │ 获取所有任务     │
                                    └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ 过滤 pending    │
                                    │ 状态任务        │
                                    └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ 排除依赖未满足   │
                                    │ 的任务          │
                                    └────────┬────────┘
                                              │
                                              ▼
                              ┌─────────────────────────┐
                              │ 按 Phase → Priority     │
                              │ 排序                    │
                              └────────────┬────────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ 选择排名第一    │
                                    │ 更新状态        │
                                    └────────┬────────┘
                                              │
                                              ▼
                                    ┌─────────────────┐
                                    │ 保存 state.json │
                                    │ 到 tasks/process│
                                    └─────────────────┘
```

### 依赖检查逻辑

```
检查任务 T 的依赖是否满足:

for each dep in T.Dependencies:
    if dep 包含 ".x":           # 通配符依赖 (如 "1.x")
        phase = 提取 phase 编号
        if 该 phase 下所有任务都已完成或已跳过:
            依赖满足
        else:
            依赖不满足
    else:                        # 精确依赖 (如 "1.1")
        if dep 任务状态 == completed OR skipped:
            依赖满足
        else:
            依赖不满足
```

---

## 3. 任务记录生成流程 (forge task submit)

```
┌─────────────────────────────────────────────────────────────────┐
│          forge task submit <task-id> --data <path>              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 从 stdin 读取   │
                    │ JSON 数据       │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 解析 RecordData │
                    │ 验证必填字段    │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 从模板生成      │
                    │ Markdown 内容   │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 写入 records/   │
                    │ <task-id>.md    │
                    └─────────────────┘
```

### RecordData 结构

```json
{
    "taskId": "1.1",
    "status": "completed",
    "summary": "实现摘要",
    "filesCreated": ["path/to/new/file.go"],
    "filesModified": ["path/to/modified/file.go"],
    "keyDecisions": ["决策1", "决策2"],
    "testsPassed": 5,
    "testsFailed": 0,
    "coverage": 85.5,
    "acceptanceCriteria": [
        {"criterion": "功能正常", "met": true}
    ],
    "notes": "可选备注",
    "typeReclassification": {
        "originalType": "fix",
        "actualType": "cleanup",
        "reason": "flaky test, not introduced by this feature"
    }
}
```

---

## 4. verify-task-done 流程

```
┌─────────────────────────────────────────────────────────────────┐
│                   forge verify-task-done                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 检查 task-state │
                    │ 是否存在        │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ 无 task-state   │             │ 有 task-state   │
    │ 返回成功(0)     │             │ 检查任务状态    │
    └─────────────────┘             └────────┬────────┘
                                              │
                              ┌───────────────┴───────────────┐
                              │                               │
                              ▼                               ▼
                    ┌─────────────────┐             ┌─────────────────┐
                    │ 任务已完成      │             │ 任务未完成      │
                    │ 检查记录文件    │             │ 返回失败(2)     │
                    └────────┬────────┘             └─────────────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ 有记录文件      │             │ 无记录文件      │
    │ 返回成功(0)     │             │ 返回失败(2)     │
    └─────────────────┘             └─────────────────┘

注意: verify-task-done 只验证状态，不删除任何文件。
```

---

## 5. cleanup 流程

```
┌─────────────────────────────────────────────────────────────────┐
│                       forge cleanup                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 检查 task-state │
                    │ 是否存在        │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ 无 task-state   │             │ 有 task-state   │
    │ 直接退出(0)     │             │ 检查任务状态    │
    └─────────────────┘             └────────┬────────┘
                                              │
                              ┌───────────────┴───────────────┐
                              │                               │
                              ▼                               ▼
                    ┌─────────────────┐             ┌─────────────────┐
                    │ 任务已完成      │             │ 任务未完成      │
                    │ 删除状态文件    │             │ 保留状态文件    │
                    └────────┬────────┘             └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 删除:           │
                    │ - state.json    │
                    │ - record.json   │
                    │   (如存在)      │
                    └─────────────────┘
```

---

## 6. quality-gate 流程

```
┌─────────────────────────────────────────────────────────────────┐
│                      forge quality-gate                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 加载 index.json │
                    │ 获取所有任务    │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ 全部 completed  │             │ 有未完成任务     │
    │ 或 skipped      │             │ 静默退出(0)     │
    └────────┬────────┘             └─────────────────┘
              │
              ▼
    ┌─────────────────┐
    │ 若 e2e 脚本存在 │
    │ 但未毕业则警告  │
    └────────┬────────┘
              │
              ▼
    ┌─────────────────┐
    │ 运行项目级      │
    │ 单元/集成测试   │
    └────────┬────────┘
              │
              ▼
    ┌─────────────────┐
    │ E2e 回归测试    │
    │ (just test-e2e) │
    │ 若可用          │
    └────────┬────────┘
              │
    ┌─────────┴──────────┐
    │                    │
    ▼                    ▼
┌──────────┐      ┌──────────────────────────┐
│ 通过:    │      │ 失败: 保存原始输出       │
│ exit 0   │      │ 阻止 hook → Agent 读取   │
└──────────┘      │ 原始输出 → forge task add  │
                  │ → 声明修复任务           │
                  └──────────────────────────┘
```

**注意**: Feature e2e 测试不由此 hook 运行。
由 T-test-3（`run-e2e-tests` 任务）负责。
此 hook 是项目健康门禁：单元/集成测试 + 回归套件。

**测试命令检测顺序：**
1. `testCommand` 字段在 `index.json` 中
2. `justfile`/`Justfile` 含 `test` recipe → `just test`
3. `Makefile` (含 test: target) → `make test`
4. `go.mod` → `go test ./...`
5. `package.json` (含 scripts.test) → `npm test`
6. `pytest.ini` / `pyproject.toml` → `pytest`

---

## 7. 验证流程 (forge task validate-index)

```
┌─────────────────────────────────────────────────────────────────┐
│                forge task validate-index [file]                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 加载 index.json │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ JSON 语法验证   │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 必填字段检查    │
                    │ (id, title)     │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 依赖引用验证    │
                    │ (引用存在的ID)  │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 循环依赖检测    │
                    │ (DFS 拓扑排序)  │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 文件存在性检查  │
                    │ (tasks/*.md)    │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 输出验证结果    │
                    └─────────────────┘
```

---

## 8. 循环依赖检测算法

```go
// 深度优先搜索检测循环
func detectCycle(tasks map[string]Task) []string {
    visited := make(map[string]bool)
    recStack := make(map[string]bool)

    var cycle []string

    var dfs func(id string) bool
    dfs = func(id string) bool {
        visited[id] = true
        recStack[id] = true

        for _, dep := range tasks[id].Dependencies {
            if !visited[dep] {
                if dfs(dep) {
                    cycle = append(cycle, dep)
                    return true
                }
            } else if recStack[dep] {
                cycle = append(cycle, dep)
                return true
            }
        }

        recStack[id] = false
        return false
    }

    for id := range tasks {
        if !visited[id] {
            dfs(id)
        }
    }

    return cycle
}
```

---

## 9. 典型开发工作流

### 方式一：使用 Git 分支（推荐）

```bash
# 1. 创建 feature 分支
$ git checkout -b feature/auth-login

# 2. 领取任务（自动识别 feature: auth-login）
$ forge task claim
> Claimed task 1.1: 实现用户认证

# 3. 开发任务
# ... 编写代码、测试 ...

# 4. 生成记录
$ forge task submit 1.1 --data record.json

# 5. 更新状态
$ forge task status 1.1 completed

# 6. 提交代码（verify-task-done 自动验证）
$ git commit -m "feat(auth): implement login"
> verify-task-done: 任务已完成且有记录 → 允许提交

# 7. 循环
$ forge task claim
> Claimed task 1.2: 实现权限检查
```

### 方式二：使用 Git Worktree

```bash
# 1. 创建 worktree（自动识别 feature）
$ git worktree add ../auth-login feature/auth-login

# 2. 在 worktree 中工作
$ cd ../auth-login
$ forge task claim
> Claimed task 1.1: 实现用户认证

# 3. 开发、记录、提交 ...
```

### 方式三：手动设置 Feature

```bash
# 1. 手动设置 feature
$ forge feature auth-login

# 2. 领取任务
$ forge task claim

# 3. 开发、记录、提交 ...
```

---

## 10. 错误处理流程

```
错误类型              处理方式
─────────────────────────────────────────────────
Feature 不存在        返回错误，提示运行: forge feature <slug>
多个活跃 Feature      返回错误，列出活跃 feature，提示切换
Task-state 损坏       返回错误，建议手动删除
index.json 语法错误   返回详细错误位置
依赖不存在            返回错误，列出无效依赖
循环依赖              返回错误，显示循环路径
文件不存在            返回警告，不阻止操作
```

---

## 11. Feature 状态管理

### 设置 Feature

```bash
$ forge feature <slug>
```

创建 `docs/features/<slug>/tasks/process/` 目录作为 feature 的运行时状态存储。

### 显示当前 Feature

```bash
$ forge feature
> Current feature: auth-login
```

### Feature 识别优先级

```
优先级    来源                              示例
─────────────────────────────────────────────────────────────────
1        Git Worktree                      worktrees/auth-login → auth-login
2        Git 分支名称                       feature/auth-login → auth-login
3        State 文件                         docs/features/auth-login/tasks/process/state.json
4        唯一 feature 目录                  只有一个 feature 有 index.json 时使用

```

### 从 Git 推断 Feature 的规则

```
分支前缀           → 移除前缀
───────────────────────────────────
feature/           → 移除
feat/              → 移除
fix/               → 移除
bugfix/            → 移除
hotfix/            → 移除
chore/             → 移除
main/master/HEAD   → 忽略，使用目录扫描
其他               → 替换 / 为 -
```

示例：
- `feature/user-auth` → `user-auth`
- `custom/branch/name` → `custom-branch-name`
- `main` → 使用目录扫描
```

---

## 12. 动态任务添加流程

```
┌─────────────────────────────────────────────────────────────────┐
│                       forge task add                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ FindProjectRoot │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ RequireFeature  │
                    └────────┬────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ 验证 title      │
                    │ 非空            │
                    └────────┬────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐
    │ 提供 --id:      │             │ 省略 --id:      │
    │ 验证唯一性      │             │ 自动生成        │
    └────────┬────────┘             │ disc-N          │
             │                      └────────┬────────┘
             └──────────┬────────────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ 验证依赖        │
              │ 存在于 index 中 │
              └────────┬────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ 添加到 index    │
              │ 创建 .md 文件   │
              └────────┬────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ 重置 forge      │
              │ state           │
              │ (allCompleted=  │
              │  false)         │
              └────────┬────────┘
                        │
                        ▼
              ┌─────────────────┐
              │ 输出 ADDED 块   │
              └─────────────────┘
```

**自动 ID 生成（填充间隙）：**
- 扫描已有任务中的 `disc-*` 键
- 找到最小未使用整数 N（从 1 开始）
- 返回 `disc-{N}`

**Flags:**

| Flag | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `--title` | 是 | - | 任务标题 |
| `--id` | 否 | 自动 `disc-N` | 自定义任务 ID |
| `--priority` | 否 | P1 | P0/P1/P2 |
| `--depends-on` | 否 | 无 | 逗号分隔的任务 ID |
| `--estimated-time` | 否 | - | 时间估算 |
| `--breaking` | 否 | false | 触发全量测试 |
| `--description` | 否 | - | 任务正文内容 |
