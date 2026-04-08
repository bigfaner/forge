# ZCode Plugin Redesign Plan

## Context

根据 `docs/proposal/new-workflow/proposal.md`，将 zcode 插件从扁平文件结构升级为嵌套目录结构，新增 brainstorm 和 ui-design 两个 skill，重构现有 skill 的输入输出路径。这是一个 breaking change（v1.0.5 → v2.0.0）。

**目标目录结构：**
```
docs/
  proposal/<slug>/proposal.md
  features/<slug>/
    prd/
      overview.md            # PRD 总览
      user-stories.md        # 用户故事
      ui-functions.md        # UI 功能要点（需求层）
    design/
      overview.md            # 技术设计总览
      api.md                 # API 文档
      ui/                    # UI 设计产出
        *.md                 # Markdown 组件规格
        *.pen                # 外部设计工具产出（如 pencil.dev）
    tasks/                   # 不变
```

**新工作流：**
```
/brainstorm → /write-prd → /eval-prd → /design-tech → /eval-design → /breakdown-tasks
     ↓            ↓            ↓            ↓              ↓                ↓
proposal.md  prd/*.{3}   prd-eval.md  design/*.{2+}  design-eval.md   tasks/*
                                   ↗
                        /ui-design
```

---

## Phase 0: Task-CLI 路径常量更新

**目的**：更新 Go task-cli 的路径常量，使其识别新的嵌套目录结构。

### 修改文件

1. **`task-cli/pkg/feature/constants.go`**
   - `PRDFileName`: `"prd.md"` → `"overview.md"`
   - `DesignFileName`: `"design.md"` → `"overview.md"`
   - 新增常量：
     ```go
     PRDDirName         = "prd"
     DesignDirName      = "design"
     UserStoriesFile    = "user-stories.md"
     UIFunctionsFile    = "ui-functions.md"
     APIDesignFile      = "api.md"
     UIDesignDir        = "ui"
     ProposalBaseDir    = "docs/proposal"
     ProposalFileName   = "proposal.md"
     ```

2. **`task-cli/pkg/feature/paths.go`**
   - `GetFeaturePRDFile`: 改为 `filepath.Join(FeaturesDir, feature, PRDDirName, PRDFileName)`
   - `GetFeatureDesignFile`: 改为 `filepath.Join(FeaturesDir, feature, DesignDirName, DesignFileName)`
   - 新增函数：`GetFeaturePRDDir`, `GetFeatureDesignDir`, `GetFeatureUserStoriesFile`, `GetFeatureUIFunctionsFile`, `GetFeatureAPIDesignFile`, `GetFeatureUIDesignDir`, `GetProposalDir`, `GetProposalFile`

3. **`task-cli/pkg/feature/feature.go`**
   - `EnsureFeatureDir`: 增加 `prd/`、`design/`、`design/ui/` 目录创建

4. **测试文件更新**（路径从 `"prd.md"` → `"prd/overview.md"` 等）：
   - `task-cli/pkg/feature/paths_test.go`
   - `task-cli/pkg/feature/feature_test.go`
   - `task-cli/internal/cmd/check_test.go`
   - `task-cli/internal/cmd/validate_test.go`
   - `task-cli/internal/cmd/runners_test.go`
   - `task-cli/internal/cmd/feature_test.go`
   - `task-cli/internal/cmd/claim_test.go`

5. **`task-cli/internal/cmd/validate.go`** — 更新 warning 文案

6. **`task-cli/docs/OVERVIEW.md`** — 更新目录结构图

**验证**：`cd task-cli && go build ./... && go test -race -cover ./...`

---

## Phase 1: 新建模板文件

**目的**：创建所有新模板，为 skill 改造做准备。

### 新建文件

| 文件 | 说明 |
|------|------|
| `plugins/zcode/skills/brainstorm/templates/proposal.md` | proposal 模板 |
| `plugins/zcode/skills/ui-design/templates/ui-design.md` | UI 组件规格模板 |
| `plugins/zcode/skills/write-prd/templates/ui-functions.md` | UI 功能要点模板 |
| `plugins/zcode/skills/design-tech/templates/api.md` | API 文档模板 |

### 重命名/创建文件

| 现有文件 | 新文件 | 说明 |
|---------|--------|------|
| `write-prd/templates/prd.md` | `write-prd/templates/prd-overview.md` | 模板内容基本不变，重命名以明确语义 |
| `design-tech/templates/design.md` | `design-tech/templates/design-overview.md` | 同上 |

### 修改文件

- **`plugins/zcode/skills/breakdown-tasks/templates/index.json`** — `prd` 和 `design` 字段更新为新路径

---

## Phase 2: 新建 brainstorm Skill

**目的**：创建从模糊想法到结构化提案的 skill。

### 新建文件

1. **`plugins/zcode/skills/brainstorm/SKILL.md`**
   - 从用户模糊想法出发，通过协作对话产出 `docs/proposal/<slug>/proposal.md`
   - HARD-GATE：不写代码，只产出提案
   - 流程：探索上下文 → 讨论愿景 → 识别约束 → 提出范围 → 写提案 → 提交
   - 输入：用户口头/文字描述
   - 输出：`docs/proposal/<slug>/proposal.md`
   - 衔接：`/write-prd` 可选读取 proposal.md 作为输入

2. **`plugins/zcode/skills/brainstorm/examples/`** — 示例文件

---

## Phase 3: 新建 ui-design Skill

**目的**：创建从 PRD 的 ui-functions.md 到 UI 设计产出的 skill。

### 新建文件

1. **`plugins/zcode/skills/ui-design/SKILL.md`**
   - 读取 `prd/ui-functions.md`，产出 UI 设计到 `design/ui/`
   - 产出格式：
     - Markdown `.md` 文件：组件规格（布局结构、交互状态、组件层级、数据绑定）
     - 外部工具产出（如 `.pen` 文件）：引用或存放
   - 位置在流程中与 `/design-tech` 并行
   - 跳过条件：纯后端/API/CLI 项目无 UI 表面

2. **`plugins/zcode/skills/ui-design/examples/`** — 示例文件

---

## Phase 4: 重构现有 Skill

### 4.1 write-prd

**文件**：`plugins/zcode/skills/write-prd/SKILL.md`

改动点：
- Step 1：增加 `docs/proposal/<slug>/proposal.md` 可选输入检测
- Step 6：输出路径 `prd.md` → `prd/overview.md`
- Step 7：输出路径 `user-stories.md` → `prd/user-stories.md`
- 新增 Step 8：输出 `prd/ui-functions.md`（UI 功能要点，非设计稿）
- 更新 Output Documents 表（3 个文件）
- 更新目录结构图
- 模板引用更新

### 4.2 design-tech

**文件**：`plugins/zcode/skills/design-tech/SKILL.md`

改动点：
- Step 1：读取路径 `prd.md` → `prd/overview.md`
- Step 7：输出拆分为 `design/overview.md` + `design/api.md`
- 说明 `design/ui/` 由 `/ui-design` skill 填充
- 更新目录结构图和 Integration 节

### 4.3 eval-prd

**文件**：`plugins/zcode/skills/eval-prd/SKILL.md` + `templates/report.md`

改动点：
- 定位文档路径全部更新（`prd.md` → `prd/overview.md` 等）
- 增加 `prd/ui-functions.md` 可选检查维度
- 更新 agent prompt 中的路径替换

### 4.4 eval-design

**文件**：`plugins/zcode/skills/eval-design/SKILL.md` + `templates/report.md`

改动点：
- 定位文档路径全部更新
- 增加 `design/api.md` 和 `design/ui/` 存在性检查
- 更新 agent prompt 中的路径替换

### 4.5 breakdown-tasks

**文件**：`plugins/zcode/skills/breakdown-tasks/SKILL.md`

改动点：
- Step 1 扩展：读取所有可用文档（`prd/overview.md`、`prd/user-stories.md`、`prd/ui-functions.md`、`design/overview.md`、`design/api.md`、`design/ui/*`）
- 每个任务标注来源文档段落（可追溯性）
- 更新目录结构图和 trigger 条件

---

## Phase 5: 更新 Guide 和 Hooks

### 修改文件

1. **`plugins/zcode/hooks/guide.md`** — 替换 Document Index 为新的嵌套目录结构，增加 workflow 说明
2. **`plugins/zcode/.claude-plugin/plugin.json`** — version `1.0.5` → `2.0.0`，keywords 增加 `"brainstorm"`, `"ui-design"`

### 不变文件

- `hooks/hooks.json` — PostToolUse validator 对 `index.json` 不受影响
- `commands/*.md` — 所有 task 路径不变
- `agents/*.md` — 所有 task 路径不变
- 其余 7 个 skill（claim-task, set-task-status, record-task, git-commit, learn-lesson, eval-harness, improve-harness）— 无需改动

---

## Phase 6: 编译验证

1. 编译 task-cli：`cd task-cli && go build ./... && go test -race -cover ./...`
2. 安装新版 task-cli：运行 `/init-zcode`
3. 端到端验证：
   - `/brainstorm` → 产出 proposal
   - `/write-prd` → 产出 prd/ 目录下 3 个文件
   - `/eval-prd` → 能定位并评估新路径文件
   - `/design-tech` → 产出 design/ 目录
   - `/ui-design` → 产出 design/ui/ 内容
   - `/eval-design` → 能定位并评估新路径文件
   - `/breakdown-tasks` → 读取所有输入文档，产出 tasks/
   - `/claim-task` → 正常认领任务

---

## 风险与缓解

| 风险 | 缓解 |
|------|------|
| 旧 flat-file feature 不兼容 | v2.0.0 breaking change，guide.md 中提供迁移说明 |
| breakdown-tasks 上下文消耗增加（读 5-6 个文档） | 明确优先级：prd/overview.md 定义 WHAT，design/overview.md 定义 HOW，其余为补充 |
| 外部设计文件（.pen 等）的二进制格式 | skill 只负责引用和存放路径管理，不解析二进制内容 |
