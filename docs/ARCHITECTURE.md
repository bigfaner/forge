# Forge 核心架构与流程

> 本文档描述 Forge 插件的内部架构、工作流管道、Agent 协作机制和关键协议。

---

## 目录

- [架构概览](#架构概览)
- [工作流管道](#工作流管道)
- [Agent 架构](#agent-架构)
- [Quality Gate 协议](#quality-gate-协议)
- [测试生命周期](#测试生命周期)
- [对抗式评估循环](#对抗式评估循环)
- [Hooks 系统](#hooks-系统)
- [Manifest 生命周期](#manifest-生命周期)
- [目录约定](#目录约定)

---

## 架构概览

Forge 由三个核心子系统组成：

```
┌─────────────────────────────────────────────────────────────┐
│                      Forge Plugin                          │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│  │  Skills   │  │ Commands │  │  Agents  │  │  Hooks    │  │
│  │  (22)     │  │  (12)    │  │  (4)     │  │  (5)      │  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └─────┬─────┘  │
│       │              │              │               │        │
│       └──────────────┴──────┬───────┴───────────────┘        │
│                             │                                │
│                    ┌────────▼────────┐                       │
│                    │   task-cli      │                       │
│                    │   (Go binary)   │                       │
│                    └─────────────────┘                       │
└─────────────────────────────────────────────────────────────┘
```

| 子系统 | 位置 | 职责 |
|--------|------|------|
| **Skills** | `plugins/forge/skills/` | 文档生成、评估、测试生命周期等技能 |
| **Commands** | `plugins/forge/commands/` | 可直接调用的 slash commands |
| **Agents** | `plugins/forge/agents/` | 自主执行的 subagent 定义 |
| **Hooks** | `plugins/forge/hooks/hooks.json` | 生命周期事件的自动触发器 |
| **task-cli** | `task-cli/` | 任务状态管理 CLI（Go 实现） |

**数据流向**：Skills/Commands → 调用 task-cli 管理状态 → Agents 执行实际开发工作 → Hooks 自动验证和清理。

---

## 工作流管道

### 完整模式（Full Mode）

适用于复杂功能（>2h, >4 任务），经过完整的文档化流程：

```
/brainstorm ──→ /write-prd ──→ /tech-design ──→ /breakdown-tasks ──→ /run-tasks
     ↓               ↓              ↓    ↘            ↓                  ↓
 proposal.md    prd/*.{3}      design/*.{2}  ui/   tasks/*.md       自动执行
              + manifest.md  + manifest.md      + manifest.md
                                /ui-design ↗
```

详细阶段：

| 阶段 | 命令 | 产出 | 可选评估 |
|------|------|------|----------|
| 探索 | `/brainstorm` | `proposal.md` | `/eval-proposal` |
| 需求 | `/write-prd` | `prd/*.{3}` + `manifest.md` | `/eval-prd` |
| 设计 | `/tech-design` | `design/*.{2}` + `manifest.md` | `/eval-design` |
| UI | `/ui-design` | `ui/ui-design.md` + `prototype/` | `/eval-ui` |
| 拆分 | `/breakdown-tasks` | `tasks/*.md` + `index.json` + `manifest.md` | — |
| 测试 | T-test-1~7 | 测试用例、脚本、报告 | `/eval-test-cases` |
| 执行 | `/run-tasks` | 代码 + 记录 | — |
| 收尾 | `/graduate-tests` → `/consolidate-specs` | 回归测试 + 项目级规范 | — |

每个阶段执行前检查前置条件（`ls` 验证文件存在），缺失时中止并提示。

### 快速模式（Quick Mode）

适用于小功能（1-2h, 1-4 任务），跳过 PRD 和设计：

```
/quick ──→ /brainstorm ──→ /quick-tasks ──→ /run-tasks
                ↓                ↓               ↓
           proposal.md      tasks/*.md       自动执行
                          + index.json
```

**与完整模式的差异**：
- 无 PRD、无设计、无评估步骤
- `proposal.md` 是唯一输入文档
- 扁平任务列表（无 phase、无 gate、无 phase summary）
- 测试任务使用 `T-quick-1~5`（完整模式 7 个任务的子集）
- 简化 manifest（无 Traceability 表）
- `--no-test` 标志跳过所有测试任务

### 如何选择模式

| 条件 | 完整模式 | 快速模式 |
|------|---------|---------|
| 工作量 >2h | ✓ | |
| 任务数 >4 | ✓ | |
| 需要 PRD 验收标准 | ✓ | |
| 涉及架构决策 | ✓ | |
| 需要 UI 设计 | ✓ | |
| 有多阶段执行 + gate | ✓ | |
| 小改动、bug 修复 | | ✓ |

---

## Agent 架构

4 个专用 Agent，由 dispatcher 或 main session 按需分发：

### task-executor

**触发**：`/run-tasks` dispatcher 循环分发、`/execute-task` 单任务执行

**5 步工作流**：

```
Step 1: Read task definition
        ├── 读取项目知识（docs/business-rules/, docs/conventions/）
        ├── 读取 phase summary（如跨 phase）
        └── 读取任务文件
Step 2: TDD implementation
        ├── RED: 写失败测试
        ├── GREEN: 最小实现通过
        └── REFACTOR: 清理
Step 3: Quality Gate
        └── compile → fmt → lint → test（严格顺序）
Step 4: Record task（必须）
        └── 调用 record-task skill → task record CLI
Step 5: Git commit
        └── 调用 git-commit skill
```

**铁律**：
- 一次调用只执行一个任务，完成后必须 STOP
- `record-task` 是强制的——没有 record = 任务未完成
- commit 必须在 record 之后
- 禁止：`task claim`、读取下一个任务、后台任务
- 最多 3 次 subagent 调用

**错误处理**：

| 场景 | 行为 |
|------|------|
| 编译失败 | 修复后从 compile 重试 |
| fmt 失败 | 标记 `blocked`（工具链问题） |
| lint 失败 | 自修复（1 次重试），仍失败则 `blocked` |
| 测试失败 | 修复后从 compile 重试 |
| 超范围失败 | `task add` 创建 fix-task，当前任务标记 `blocked` |

**动态任务创建**（fix-task 链）：

```
source-task (blocked) → fix-task-A (blocked) → fix-task-B
                                                    ↓ completed
                                              fix-task-A → pending (auto)
                                                    ↓ completed
                                              source-task → pending (auto)
```

- 通过 `--source-task-id` 关联
- `task record` 完成时自动恢复 source task 为 pending
- 最大嵌套深度 3 层
- 已完成的 fix-task 会被自动解析到根 blocked task

### error-fixer (DEPRECATED)

> **DEPRECATED**: This agent is no longer dispatched. Use `type: fix` tasks (routed via `task prompt <id>`) for compilation/test/lint fixes, and `task prompt <id> --fix-record-missed` for record recovery.

**原触发**：task-executor 失败时分发、`/fix-bug` 手动触发

**5 步工作流**：

```
Step 1: Diagnose  — 分析错误类型、影响范围、根因
Step 2: Locate    — 读取失败文件和相关测试
Step 3: Fix       — 最小修复，保留现有功能
Step 4: Verify    — Quality Gate 全流程
Step 5: Commit    — git-commit
```

**约束**：最小改动，不做重构。

### doc-scorer

**触发**：`/eval-*` 系列评估技能

**职责**：按评分标准（rubric）对文档打分（100 分制），输出结构化评分报告和 attack points。

**工作流**：

```
Step 1: 读取文档目录 + rubric + 上轮报告（如有）
Step 2: 独立评分（不看"努力程度"，只看当前内容）
Step 3: 填写报告模板
Step 4: 输出结构化摘要
         ├── SCORE: total/100
         ├── DIMENSIONS: 各维度得分
         └── ATTACKS: 前 3 个具体弱点（含引用）
```

### doc-reviser

**触发**：`/eval-*` 系列评估技能（scorer 评分后）

**职责**：根据 attack points 修订文档，不注水。

**工作流**：

```
Step 1: 读取文档 + rubric + 评分报告
Step 2: 针对性修订（模糊→具体、缺失→补充、不一致→对齐）
Step 3: 写入 + 报告变更
```

**质量检查**：
- 每个 attack point 都已处理
- 不引入新的模糊表述
- 总字数增长不超过 30%

---

## Quality Gate 协议

所有任务执行工作流（`/execute-task`、`task-executor` agent、`/fix-bug`、`type: fix` 任务）在记录完成前必须通过 Quality Gate。

### 验证序列

```
just compile ──→ just fmt ──→ just lint ──→ just test
     │               │            │             │
     ↓ fail          ↓ fail       ↓ fail        ↓ fail
  修复→重试       blocked     自修复→重试     修复→重试
                  (工具链)     (1次机会)     (从compile重试)
```

严格顺序执行，任何一步失败则停止，不继续后续步骤。

### Scope Resolution

每个 `just <verb>` 命令前，根据任务的 `scope` 字段决定是否传递 scope 参数：

```
scope 缺失/空/"all"  → just <verb>
scope = "frontend"/"backend"
  ├── just project-type → 获取项目类型
  ├── 输出不是 frontend/backend/mixed → just <verb>（fallback）
  ├── 输出 = "mixed"                   → just <verb> <scope>
  └── 输出 = "frontend" 或 "backend"   → just <verb>（fallback，scope 冗余）
```

**设计意图**：只有 mixed 项目（前后端共存）才需要 scope 参数，纯前端或纯后端项目无需指定。

---

## 测试生命周期

三层测试架构，各有独立目的和触发时机：

```
Unit Tests ──────→ Feature E2E ──────→ Regression Suite
(per task)         (per feature)        (project-level)
     ↑                   ↑                     ↑
Quality Gate        T-test-3              all-completed hook
强制执行           gen-test-scripts       graduate-tests
```

| 层 | 命令 | 范围 | 触发 | 通过标准 |
|---|---|---|---|---|
| **Unit** | `just test [scope]` | 任务级 | 每任务 Quality Gate | 全部通过 + 覆盖率 >= 80% |
| **Feature E2E** | `just test-e2e --feature <slug>` | 功能级 | T-test-3 | Playwright 报告全绿 |
| **Regression** | `just test-e2e` | 项目级 | all-completed hook | 全部回归用例通过 |

### 测试生成管道（Full Mode T-test-1~7）

```
T-test-1: /gen-sitemap          → sitemap.json（页面元素映射）
T-test-2: /gen-test-cases       → testing/test-cases.md（结构化测试用例）
          /eval-test-cases      → 评分报告（可选）
T-test-3: /gen-test-scripts     → tests/e2e/features/<slug>/*.spec.ts
          /run-e2e-tests        → results/latest.md（执行报告）
T-test-4: /graduate-tests       → tests/e2e/<module>/（迁移到回归套件）
T-test-5: verify-regression     → 回归验证
T-test-6: /consolidate-specs    → 项目级规范提取
T-test-7: (reserved)
```

### run-e2e-tests 内部流程

```
Setup Environment → Verify Scripts (无占位符) → Run Test Specs
     → Collect Results (解析 Playwright JSON) → Generate Report → Teardown
```

- Teardown 必须执行（即使测试失败）
- 禁止修改测试脚本或跳过失败用例
- UI 失败用例通过截图分析诊断

### graduate-tests 迁移流程

```
检查标记 → 读取源脚本 → 分析现有结构 → 决定分类策略
     → 执行迁移（可能 split/merge）→ TypeScript 编译验证 → Playwright list 验证
     → 创建毕业标记 → 清理源目录
```

- 按功能模块分类（非按测试类型）
- 迁移前创建备份，验证失败则回滚
- 自动重写 `helpers.js` import 路径

---

## 对抗式评估循环

`/eval-*` 系列技能使用 doc-scorer 和 doc-reviser 进行迭代改进，直到达到目标分数：

```
                    ┌──────────────────────────────┐
                    │                              │
                    ▼                              │
doc-scorer 评分 ──→ 达标？── 是 ──→ 输出最终报告   │
                    │                              │
                    否                             │
                    │                              │
                    ▼                              │
            提取 top 3 attack points               │
                    │                              │
                    ▼                              │
            doc-reviser 修订文档                    │
                    │                              │
                    └──────────────────────────────┘
```

**可用评估技能**：

| 技能 | 评估对象 | 评分维度来源 |
|------|---------|-------------|
| `/eval-prd` | PRD 三件套 | `skills/eval-prd/templates/rubric.md` |
| `/eval-design` | 技术设计文档 | `skills/eval-design/templates/rubric.md` |
| `/eval-ui` | UI 设计文档 | 4 个利益相关者视角（用户/设计师/开发/PM） |
| `/eval-proposal` | 提案文档 | `skills/eval-proposal/templates/rubric.md` |
| `/eval-test-cases` | 测试用例文档 | `skills/eval-test-cases/templates/rubric.md` |
| `/eval-consistency` | 跨文档一致性 | PRD 作为 source of truth |
| `/eval-harness` | 测试基础设施 | 基于 OpenAI harness 工程实践 |

**主 session 职责**：orchestrator 负责：
1. 调用 doc-scorer subagent
2. 解析评分结果
3. 判断是否达标
4. 未达标时提取 attack points 并调用 doc-reviser subagent
5. 循环直到达标或达到最大迭代次数

---

## Hooks 系统

Hooks 在关键生命周期事件自动触发，确保状态一致性：

| 事件 | 触发条件 | 执行命令 | 作用 |
|------|---------|---------|------|
| **SessionStart** | 启动/清除/压缩 | `session-start` hook | 加载 forge 上下文 |
| **SubagentStart** | subagent 启动 | `session-start` hook | 为 subagent 加载上下文 |
| **PostToolUse** | Edit/Write 工具调用后 | `validate-index.sh` | 自动验证 index.json 格式 |
| **SessionEnd** | 会话结束 | `task cleanup` | 清理运行时状态 |
| **SubagentStop** | subagent 停止 | `task cleanup` | 清理 subagent 状态 |
| **Stop** | Claude 停止响应时 | `task all-completed` | 全部完成后的最终验证 |

### all-completed Hook

当所有任务完成，Claude 停止响应时作为最终安全网触发：

```
1. Quality Gate（项目级，无 scope）：
   just compile → just fmt → just lint
2. 项目级测试：
   just test
3. E2E 回归：
   just e2e-setup → just probe → just test-e2e
```

任何一步失败都会报告问题。

---

## Manifest 生命周期

`manifest.md` 是 Feature 的单一入口，由各 skill 自动维护。

### 状态流转

```
(none) ──→ prd ──→ design ──→ tasks ──→ in-progress ──→ completed
  │          │         │         │           │
  │      /write-prd    │    /breakdown-  首次 task claim
  │      完成          │    tasks 完成      (或 /execute-task)
  │                    │
  │              /tech-design +
  │              /ui-design 完成
  │              （后者完成时设置）
```

| 状态 | 含义 | 设置者 |
|------|------|--------|
| (none) | 无文档 | — |
| `prd` | PRD 就绪 | `/write-prd` |
| `design` | 设计就绪 | `/tech-design` 或 `/ui-design`（后完成者） |
| `tasks` | 任务已拆分 | `/breakdown-tasks` |
| `in-progress` | 执行中 | 首次 `task claim` |
| `completed` | 全部完成 | all tasks done |

### Manifest 内容结构

```markdown
# Feature: <name>

## Status
<current status>

## Documents
| Document | Path | Summary |
|----------|------|---------|
| PRD      | prd/prd-spec.md | ... |

## Traceability（Full Mode only）
| PRD Requirement | Design Section | Task ID |
|-----------------|----------------|---------|
| ...             | ...            | ...     |
```

---

## 目录约定

### Feature 工作区

```
docs/features/<slug>/
├── manifest.md             # 单一入口（自动维护）
├── prd/
│   ├── prd-spec.md         # 需求规格
│   ├── prd-user-stories.md # 用户故事
│   └── prd-ui-functions.md # UI 功能（可选）
├── design/
│   ├── tech-design.md      # 技术设计
│   └── api-handbook.md     # API 文档（可选）
├── ui/
│   ├── ui-design.md        # UI 规格（可选）
│   ├── DESIGN.md           # 自定义风格（可选）
│   └── prototype/          # HTML 原型（可选）
├── testing/
│   └── test-cases.md       # 测试用例
├── tasks/
│   ├── index.json          # 任务定义（核心）
│   ├── *.md                # 任务详情
│   ├── process/            # 运行时状态（不提交）
│   │   ├── state.json      #   当前任务状态
│   │   └── record.json     #   进行中的记录
│   ├── records/            # 执行记录（task record 生成）
│   └── specs/              # 规范提取预览（consolidate-specs）
└── eval/                   # 评估报告（可选）
```

### 项目级文档

```
docs/
├── ARCHITECTURE.md          # 系统架构（本文档）
├── business-rules/          # 跨功能业务规则（/consolidate-specs 写入）
│   ├── auth.md              #   按领域组织
│   └── ...
├── conventions/             # 技术规范（/consolidate-specs 写入）
│   ├── api.md               #   API 约定
│   ├── testing.md           #   测试约定
│   └── ...
├── decisions/               # 技术决策（/record-decision）
├── lessons/                 # 经验教训（/learn-lesson）
└── sitemap/sitemap.json     # 页面元素映射（/gen-sitemap）
```

### 测试目录

```
tests/e2e/
├── playwright.config.ts     # Playwright 配置
├── helpers.ts               # 共享工具函数
├── features/<slug>/         # 功能级测试（staging）
│   ├── *.spec.ts            #   测试脚本
│   └── results/             #   执行结果
├── .graduated/              # 毕业标记
│   └── <slug>               #   标记该功能已毕业
└── <module>/                # 回归测试（毕业后）
    └── *.spec.ts            #   按功能模块组织
```

### 生成规则

| 目录 | 生成方式 | 是否提交 |
|------|---------|---------|
| `process/` | task-cli 运行时 | **否**（gitignore） |
| `testing/` | `/gen-test-cases` | 是 |
| `tests/e2e/features/` | `/gen-test-scripts` | 是 |
| `tests/e2e/` (根级) | `/graduate-tests` | 是 |
| `records/` | `task record` | 是 |
| `specs/` | `/consolidate-specs` | 是（用户确认后） |
