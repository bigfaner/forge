# forge CLI 功能概览

> 基于 features 目录结构的任务管理 CLI 工具，为 Claude Code 工作流提供智能任务声明与依赖管理。

## 核心功能

### 1. 智能任务声明 (`forge task claim`)

基于多维度策略自动选择下一个可用任务：

| 维度 | 优先级规则 |
|------|-----------|
| Priority | P0 > P1 > P2 |
| Dependencies | 仅声明依赖已满足的任务 |
| In-Progress | 自动恢复进行中的任务 |

**依赖语法支持：**
- 精确匹配：`1.1`, `1.2`
- 通配符匹配：`1.x`（前缀级别依赖）

### 2. 任务记录生成 (`forge task submit`)

从 JSON 输入生成结构化 markdown 执行记录，包含：

- 任务摘要与状态
- 时间追踪
- 创建/修改的文件清单
- 关键决策
- 测试结果
- 验收标准确认

**验证规则（hard validation）：**

| Condition | Error | Fix |
|-----------|-------|-----|
| `status=completed` + `testsPassed=0` + `testsFailed=0` + `coverage >= 0` | No test evidence | Run tests, or set `coverage: -1.0` |
| `status=completed` + any `acceptanceCriteria.met=false` | Unmet AC | Fix issue, or set `status: "blocked"` |
| `summary` empty or whitespace | Missing summary | Provide a summary |

Override with `--force`: `forge task submit <id> --data record.json --force`

### 3. 状态管理

| 命令 | 功能 |
|------|------|
| `forge task status <id>` | 查询任务状态 |
| `forge task status <id> <status>` | 更新任务状态 |
| `forge task query <id>` | 查询任务详情 |
| `forge feature [slug]` | 设置/显示当前 feature |

**状态值：** `pending`, `in_progress`, `completed`, `blocked`, `skipped`

### 4. 校验与验证

| 命令 | 功能 |
|------|------|
| `forge task validate-index [file]` | 验证 index.json 结构 |
| `forge task check-deps` | 检查所有任务依赖 |

**验证规则：**
- JSON 语法检查
- 必填字段验证
- 依赖引用有效性
- 循环依赖检测
- 文件存在性检查

### 5. Claude Code 集成命令

| 命令 | 用途 | 功能 |
|------|------|------|
| `forge verify-task-done` | PreToolUse (git commit) | 验证任务完成状态，阻止未完成任务提交 |
| `forge cleanup` | Stop | 清理已完成任务的状态文件 |
| `forge quality-gate` | Stop hook | 检查所有任务是否完成，若完成则自动运行测试 |

**all-completed 行为：**
- 所有任务均为 `completed` 或 `skipped` → 运行项目级测试 + e2e 回归，exit 0
- 任意任务为 `pending`/`in_progress`/`blocked` → 静默退出，exit 0
- 无 feature 或无 project root → 静默退出，exit 0

**e2e 测试失败恢复：**
- 当回归测试 (`just test-e2e`) 失败时，保存原始输出到 `tests/results/raw-output.txt`
- 阻止 Stop hook，指示 agent 分析失败原因并使用 `forge task add` 创建修复任务
- agent 读取原始输出，确定根因，动态添加修复任务

**feature e2e 测试（不由本 hook 运行）：**
- Feature e2e 执行由 T-test-3（`run-e2e-tests` 任务）负责
- 若 `tests/e2e/features/<feature>/` 存在但无毕业标记，hook 打印 WARNING 引导迁移

**e2e 测试标签晋升模型：**
- 晋升通过 `/run-tests` 完成——测试通过后自动替换 `@feature` 为 `@regression`
- 标签生命周期：`@feature`（新生成，验证中）-> `@regression`（已验证，回归测试）

**测试命令自动检测顺序（项目级）：**
1. `index.json` 中的 `testCommand` 字段（显式配置）
2. `justfile`/`Justfile` 含 `test` recipe → `just test`
3. `Makefile` 含 `test:` target → `make test`
4. `go.mod` 存在 → `go test ./...`
5. `package.json` 含 `scripts.test` → `npm test`
6. `pytest.ini` / `pyproject.toml` 存在 → `pytest`

**e2e 测试检测顺序：**
1. `justfile`/`Justfile` 含 `test-e2e` recipe → `just test-e2e`

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
│       ├── testing/
│       │   ├── test-cases.md      # 测试用例（含 target 字段）
│       │   └── results/
│       │       └── latest.md      # e2e 测试结果报告
│       └── tasks/
│           ├── index.json          # 任务定义
│           ├── process/            # 运行时状态
│           │   ├── state.json
│           │   └── record.json
│           ├── 1.1-<title>.md     # 任务详情
│           └── records/            # 执行记录
├── tests/
│   └── e2e/                       # 回归测试套件（通过标签晋升）
│       ├── .graduated/            # 遗留毕业标记文件
│       │   └── <slug>             # YAML 标记（schema_version, status, timestamp, source, targets, modules, testCount）
│       ├── ui/<page>/             # UI 测试（按页面聚合）
│       │   └── ui.spec.ts
│       ├── api/<resource>/        # API 测试（按资源聚合）
│       │   └── api.spec.ts
│       └── cli/<command>/         # CLI 测试（按命令聚合）
│           └── cli.spec.ts
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
    Breaking      bool     `json:"breaking,omitempty"`      // 全局性变更标记，完成后触发全量测试
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
    TestCommand  string          `json:"testCommand,omitempty"`
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
forge task claim              # 声明下一个任务
forge task submit 1.1         # 生成任务记录
forge task submit 1.1 --force # 生成任务记录（跳过验证）
forge task add --title "Fix: ..." --priority P0 --breaking  # 动态添加新任务
forge task status 1.1         # 查询任务状态
forge task status 1.1 done    # 更新状态
forge task query 1.1          # 查询任务详情
forge feature auth            # 切换 feature
forge task check-deps         # 依赖检查
forge task validate-index     # 验证 index.json
forge verify-task-done        # 验证任务完成（git commit hook）
forge cleanup                 # 清理已完成任务状态（stop hook）
```
