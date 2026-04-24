# claude-task-cli 功能概览

> 基于 features 目录结构的任务管理 CLI 工具，为 Claude Code 工作流提供智能任务声明与依赖管理。

## 核心功能

### 1. 智能任务声明 (`task claim`)

基于多维度策略自动选择下一个可用任务：

| 维度 | 优先级规则 |
|------|-----------|
| Priority | P0 > P1 > P2 |
| Dependencies | 仅声明依赖已满足的任务 |
| In-Progress | 自动恢复进行中的任务 |

**依赖语法支持：**
- 精确匹配：`1.1`, `1.2`
- 通配符匹配：`1.x`（前缀级别依赖）

### 2. 任务记录生成 (`task record`)

从 JSON 输入生成结构化 markdown 执行记录，包含：

- 任务摘要与状态
- 时间追踪
- 创建/修改的文件清单
- 关键决策
- 测试结果
- 验收标准确认

### 3. 状态管理

| 命令 | 功能 |
|------|------|
| `task status <id>` | 查询任务状态 |
| `task status <id> <status>` | 更新任务状态 |
| `task query <id>` | 查询任务详情 |
| `task feature [slug]` | 设置/显示当前 feature |

**状态值：** `pending`, `in_progress`, `completed`, `blocked`, `skipped`

### 4. 校验与验证

| 命令 | 功能 |
|------|------|
| `task validate [file]` | 验证 index.json 结构 |
| `task check` | 检查所有任务依赖 |

**验证规则：**
- JSON 语法检查
- 必填字段验证
- 依赖引用有效性
- 循环依赖检测
- 文件存在性检查

### 5. Claude Code 集成命令

| 命令 | 用途 | 功能 |
|------|------|------|
| `task verifyCompletion` | PreToolUse (git commit) | 验证任务完成状态，阻止未完成任务提交 |
| `task cleanup` | Stop | 清理已完成任务的状态文件 |
| `task all-completed` | Stop hook | 检查所有任务是否完成，若完成则自动运行测试 |

**all-completed 行为：**
- 所有任务均为 `completed` 或 `skipped` → 运行 feature e2e 测试 + 项目级测试，exit 0
- 任意任务为 `pending`/`in_progress`/`blocked` → 静默退出，exit 1
- 无 feature 或无 project root → 静默退出，exit 1

**测试命令自动检测顺序（项目级）：**
1. `index.json` 中的 `testCommand` 字段（显式配置）
2. `justfile`/`Justfile` 含 `test` recipe → `just test`
3. `Makefile` 含 `test:` target → `make test`
4. `go.mod` 存在 → `go test ./...`
5. `package.json` 含 `scripts.test` → `npm test`
6. `pytest.ini` / `pyproject.toml` 存在 → `pytest`

**e2e 测试检测顺序：**
1. `justfile`/`Justfile` 含 `test-e2e` recipe → `just test-e2e`
2. `Makefile` 含 `test-e2e:` target → `make test-e2e`
3. `testing/scripts/package.json` 存在 → `npm run test:all --if-present`

---

## 目录结构约定

```
project-root/
├── docs/
│   ├── proposals/<slug>/           # /brainstorm 产出
│   │   └── proposal.md
│   └── features/<slug>/            # Feature 工作区
│       ├── manifest.md             # Feature 索引 & 可追溯性映射
│       ├── prd/
│       │   ├── prd-spec.md         # PRD Spec
│       │   ├── prd-user-stories.md # 用户故事
│       │   └── prd-ui-functions.md # UI 功能要点（可选）
│       ├── design/
│       │   ├── tech-design.md      # 技术设计
│       │   └── api-handbook.md     # API 文档
│       ├── ui/
│       │   └── ui-design.md        # UI 设计规格（可选）
│       └── tasks/
│           ├── index.json          # 任务定义
│           ├── process/            # 运行时状态
│           │   ├── state.json
│           │   └── record.json
│           ├── 1.1-<title>.md     # 任务详情
│           └── records/            # 执行记录
```

### 项目根目录检测

工具自动检测项目根目录，支持多种项目类型和 monorepo 结构：

**检测优先级**（从高到低）：
1. 环境变量：`CLAUDE_PROJECT_DIR` > `PROJECT_ROOT`
2. Workspace 标记：`go.work`, `pnpm-workspace.yaml`, `lerna.json`, `turbo.json`, `nx.json`, `WORKSPACE`, `settings.gradle*`
3. Project 标记：`go.mod`, `package.json`, `Cargo.toml`, `pyproject.toml`, `pom.xml`, `build.gradle*`
4. VCS 边界：`.git`, `.hg`

**支持的项目类型**：
- Go (`go.mod`, `go.work`)
- Node.js (`package.json`)
- Rust (`Cargo.toml`)
- Python (`pyproject.toml`, `setup.py`)
- Java/Maven (`pom.xml`)
- Java/Gradle (`build.gradle`, `settings.gradle`)
- Bazel (`WORKSPACE`)
- 通用 Git 仓库 (`.git`)

**Feature 自动识别**：
- Git worktree 名称 → feature slug
- Git 分支名称 (如 `feature/auth-login`) → auth-login
- 目录扫描（有 `tasks/process/state.json` 的 feature 优先）

**状态隔离**：每个 feature 的运行时状态存放在独立的 `docs/features/<slug>/tasks/process/` 目录下，避免多个 feature 状态冲突。

---

## 数据模型

### Task

```go
type Task struct {
    ID            string   `json:"id"`                      // 任务ID (如 "1.1")
    Title         string   `json:"title"`                   // 任务标题
    Priority      string   `json:"priority"`                // P0/P1/P2
    EstimatedTime string   `json:"estimatedTime,omitempty"` // 预估时间
    Dependencies  []string `json:"dependencies,omitempty"`  // 依赖任务ID列表
    Status        string   `json:"status"`                  // pending/in_progress/completed/blocked/skipped
    File          string   `json:"file"`                    // 任务文件
    Record        string   `json:"record"`                  // 记录文件
}
```

### TaskIndex

```go
type TaskIndex struct {
    Feature      string          `json:"feature"`
    PRD          string          `json:"prd,omitempty"`
    Design       string          `json:"design,omitempty"`
    Created      string          `json:"created,omitempty"`
    Status       string          `json:"status,omitempty"`
    Tasks        map[string]Task `json:"tasks"`
    StatusEnum   []string        `json:"statusEnum,omitempty"`
    PriorityEnum []string        `json:"priorityEnum,omitempty"`
}
```

---

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.21 |
| CLI 框架 | github.com/spf13/cobra |
| 外部依赖 | 仅 cobra (极轻量) |

---

## 架构约束

```
依赖方向: cmd → internal → pkg (严禁反向)
模块交互: 通过接口/类型定义，不直接依赖内部实现
```

## 命令速查

```bash
task claim              # 声明下一个任务
task record 1.1         # 生成任务记录
task status 1.1         # 查询任务状态
task status 1.1 done    # 更新状态
task query 1.1          # 查询任务详情
task feature auth       # 切换 feature
task check              # 依赖检查
task validate           # 验证 index.json
task verifyCompletion   # 验证任务完成（git commit hook）
task cleanup            # 清理已完成任务状态（stop hook）
```
