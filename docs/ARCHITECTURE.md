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
- [v3.0.0 子系统](#v3.0.0-子系统)

---

## 架构概览

Forge 由三个核心子系统组成：

```
┌─────────────────────────────────────────────────────────────┐
│                      Forge Plugin                          │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│  │  Skills   │  │ Commands │  │  Agents  │  │  Hooks    │  │
│  │  (21)     │  │  (16)    │  │  (1)     │  │  (5)      │  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └─────┬─────┘  │
│       │              │              │               │        │
│       └──────────────┴──────┬───────┴───────────────┘        │
│                             │                                │
│                    ┌────────▼────────┐                       │
│                    │   forge CLI     │                       │
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

**数据流向**：Skills/Commands → 调用 forge CLI 管理状态 → Agents 执行实际开发工作 → Hooks 自动验证和清理。forge CLI 是统一的 Go 二进制，集成了任务状态管理、feature 管理、surface 检测等功能。

---

## 工作流管道

### 完整模式（Full Mode）

适用于复杂功能（>10 任务），经过完整的文档化流程：

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
| 测试 | T-test-gen-journeys / T-eval-journey / T-test-gen-contracts / T-eval-contract / T-test-gen-scripts / T-test-run | Journey、Contract、测试脚本、报告 | `/eval-journey` → `/eval-contract` |
| 执行 | `/run-tasks` | 代码 + 记录 | — |
| 收尾 | `/consolidate-specs` | 项目级规范 | — |

每个阶段执行前检查前置条件（`ls` 验证文件存在），缺失时中止并提示。

### 快速模式（Quick Mode）

适用于小功能（1-10 任务），跳过 PRD 和设计：

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
- 测试任务使用 `T-test-gen-journeys` + `T-test-run`（无 eval 质量关卡的精简串行链路）
- 简化 manifest（无 Traceability 表）
- 纯文档 feature 自动跳过测试任务，生成 T-eval-doc 替代

### 如何选择模式

| 条件 | 完整模式 | 快速模式 |
|------|---------|---------|
| 任务数 >10 | ✓ | |
| 需要 PRD 验收标准 | ✓ | |
| 涉及架构决策 | ✓ | |
| 需要 UI 设计 | ✓ | |
| 有多阶段执行 + gate | ✓ | |
| 小改动、bug 修复 | | ✓ |

---

## Agent 架构

1 个专用 Agent（`task-executor`），由 dispatcher 或 main session 按需分发。评估技能（`/eval-*`）内部通过 scorer/reviser 协议文件实现迭代改进，不使用独立 agent。

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
        └── 调用 submit-task skill → forge task submit CLI
Step 5: Git commit
        └── 调用 git-commit skill
```

**铁律**：
- 一次调用只执行一个任务，完成后必须 STOP
- `submit-task` 是强制的——没有 record = 任务未完成
- commit 必须在 record 之后
- 禁止：`forge task claim`、读取下一个任务、后台任务
- 最多 3 次 subagent 调用

**错误处理**：

| 场景 | 行为 |
|------|------|
| 编译失败 | 修复后从 compile 重试 |
| fmt 失败 | 标记 `blocked`（工具链问题） |
| lint 失败 | 自修复（1 次重试），仍失败则 `blocked` |
| 测试失败 | 修复后从 compile 重试 |
| 超范围失败 | `forge task add` 创建 fix-task，当前任务标记 `blocked` |

**动态任务创建**（fix-task 链）：

```
source-task (blocked) → fix-task-A (blocked) → fix-task-B
                                                    ↓ completed
                                              fix-task-A → pending (auto)
                                                    ↓ completed
                                              source-task → pending (auto)
```

- 通过 `--source-task-id` 关联
- `forge task submit` 完成时自动恢复 source task 为 pending
- 最大嵌套深度 3 层
- 已完成的 fix-task 会被自动解析到根 blocked task

---

## Quality Gate 协议

所有任务执行工作流（`/execute-task`、`task-executor` agent、`/fix-bug`、`type: fix` 任务）在记录完成前必须通过 Quality Gate。

### 验证序列

Forge 使用三种 Gate Sequence，根据任务类型和触发场景选择：

| Gate Sequence | Steps | 适用场景 |
|---|---|---|
| **FullGateSequence** | `compile → fmt → lint → unit-test → test → probe` | all-completed hook（项目级全量验证） |
| **UnitGateSequence** | `compile → fmt → lint → unit-test` | Breaking 任务 submit gate（快速反馈） |
| **NonBreakingGateSequence** | `compile → fmt → lint` | Non-breaking 任务 submit gate（仅静态检查） |

```
UnitGateSequence（Breaking 任务 submit）：
just compile ──→ just fmt ──→ just lint ──→ just unit-test
     │               │            │               │
     ↓ fail          ↓ fail       ↓ fail          ↓ fail
  修复→重试       blocked     自修复→重试      修复→重试
                  (工具链)     (1次机会)       (从compile重试)

FullGateSequence（all-completed hook）：
just compile ──→ just fmt ──→ just lint ──→ just unit-test ──→ just test ──→ just probe
```

严格顺序执行，任何一步失败则停止，不继续后续步骤。

### Scope Resolution

每个 `just <verb>` 命令前，根据任务的 `scope` 字段决定是否传递 scope 参数（实现：`forge-cli/pkg/just/just.go` `ResolveScope`）：

```
scope 缺失/空/"all"  → just <verb>（不传 scope）
scope = "frontend"/"backend"/其他非空值
  ├── just --dry-run compile <scope> 成功 → just <verb> <scope>
  └── just --dry-run compile <scope> 失败 → just <verb>（fallback）
```

**设计意图**：通过 `just --dry-run` 探测 justfile 的 compile recipe 是否接受 scope 参数。只有 justfile 定义了带 scope 参数的 recipe 时才传递 scope，否则静默降级。

---

## 测试生命周期

两层测试 recipe 模型，解耦语言级单元测试与 surface 级测试（按 Surface 类型区分，详见 [test-type-model.md](../../plugins/forge/skills/test-guide/references/test-type-model.md)）：

```
unit-test (语言级) ──────→ test (Surface 级，按 surface type 区分)
(per task submit)          (all-completed hook)
     ↑                           ↑
UnitGateSequence          FullGateSequence
compile→fmt→lint→unit-test  compile→fmt→lint→unit-test→test→probe
```

| 层 | 命令 | 范围 | 触发 | 通过标准 |
|---|---|---|---|---|
| **Unit** | `just unit-test [scope]` | 任务级 | Breaking 任务 submit gate | 全部通过 + 覆盖率 >= 80% |
| **Surface Test** | `just test [journey]` | 功能/项目级 | all-completed hook | 全部 surface 测试通过（CLI 功能测试 / API 功能测试 / Web 端到端测试等） |

### 测试生成管道

#### Breakdown 模式（Full Mode）

含 eval 质量关卡的完整链路，每个阶段自动生成对应任务：

```
gen-journeys ──→ eval-journey ──→ gen-contracts ──→ eval-contract ──→ gen-scripts ──→ run ──→ verify
     │                │                  │                  │               │            │
     │                │                  │                  │               │            └─ /run-tests (tag promotion)
     │                │                  │                  │               └─ /run-tests
     │                │                  │                  └─ /eval-contract（6 维度门禁）
     │                │                  └─ /gen-contracts（6 维度合约 + 边界衍生）
     │                └─ /eval-journey（6 维度 1000 分制，总分 ≥850）
     └─ /gen-journeys（Journey 文档 + surface 检测 + 风险分级）
```

**任务映射**：

| 任务 | Skill | 产出 | 自动生成 |
|------|-------|------|---------|
| T-test-gen-journeys | `/gen-journeys` | `testing/<journey>/journey.md` | 是 |
| T-test-eval-journey | `/eval-journey` | 评分报告 | 是 |
| T-test-gen-contracts | `/gen-contracts` | `testing/<journey>/contracts/step-*.md` | 是 |
| T-test-eval-contract | `/eval-contract` | 评分报告 | 是 |
| T-test-gen-scripts | `/gen-test-scripts` | `tests/<journey>/*` | 是 |
| T-test-run | `/run-tests` | `results/latest.md` | 是 |

前置任务：`/gen-web-sitemap`（生成 `sitemap.json` 页面元素映射）。

#### Quick 模式

跳过 eval 质量关卡的精简链路，采用 **staged across types** 拓扑：

```
┌──────────────────────────────┐
│ gen-journeys                 │  ← 单一任务覆盖所有 surface type
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│ run-test (各 surface 串行)    │  ← 无 gen-contracts / gen-scripts，串行执行
│   ├─ key: api                │
│   ├─ key: web                │
│   └─ key: cli                │
└──────────────┬───────────────┘
               │
               ▼
          run → verify
```

**与 Breakdown 模式的差异**：
- 无 eval-journey / eval-contract 质量关卡
- gen-journeys 以 `proposal.md` 为输入（非 PRD user stories）
- gen-journeys 通过 `AUTO_COMMIT=true` 跳过人工审批
- 若 gen-journeys 产出零 Journey（proposal.md 信息不足），任务 abort 并输出诊断信息

#### 依赖解析

- **Breakdown 模式**：基于 `findTaskIndexByPrefix` 的 ID 查找（非硬编码索引）
- **Quick 模式**：staged across types 策略，gen-journeys 单一任务后接各 surface 的 run-test 串行链（无 gen-contracts / gen-scripts）

### run-tests 内部流程

```
Setup Environment → Verify Scripts (无占位符) → Run Test Specs
     → Collect Results (解析测试输出) → Generate Report → Teardown
```

- Teardown 必须执行（即使测试失败）
- 禁止修改测试脚本或跳过失败用例
- UI 失败用例通过截图分析诊断

---

## 对抗式评估循环

`/eval-*` 系列技能使用 scorer/reviser 协议进行迭代改进，直到达到目标分数。scorer 和 reviser 是 eval skill 内部的协议文件（`skills/eval/`），不是独立 agent：

```
                    ┌──────────────────────────────┐
                    │                              │
                    ▼                              │
scorer 评分 ──→ 达标？── 是 ──→ 输出最终报告   │
                    │                              │
                    否                             │
                    │                              │
                    ▼                              │
            提取 top 3 attack points               │
                    │                              │
                    ▼                              │
            reviser 修订文档                    │
                    │                              │
                    └──────────────────────────────┘
```

**可用评估技能**：

| 技能 | 评估对象 | 评分维度来源 |
|------|---------|-------------|
| `/eval-prd` | PRD 三件套 | `skills/eval/rubrics/prd.md` |
| `/eval-design` | 技术设计文档 | `skills/eval/rubrics/design.md` |
| `/eval-ui` | UI 设计文档 | `skills/eval/rubrics/ui-web.md` / `ui-mobile.md` / `ui-tui.md` |
| `/eval-proposal` | 提案文档 | `skills/eval/rubrics/proposal.md` |
| `/eval-journey` | Journey 文档 | `skills/eval/rubrics/journey.md` |
| `/eval-contract` | Contract 文档 | `skills/eval/rubrics/contract.md` |

**主 session 职责**：orchestrator 负责：
1. 调用 scorer 协议
2. 解析评分结果
3. 判断是否达标
4. 未达标时提取 attack points 并调用 reviser 协议
5. 循环直到达标或达到最大迭代次数

---

## Hooks 系统

Hooks 在关键生命周期事件自动触发，确保状态一致性：

| 事件 | 触发条件 | 执行命令 | 作用 |
|------|---------|---------|------|
| **SessionStart** | 启动/清除/压缩 | `session-start` hook | 加载 forge 上下文 |
| **SubagentStart** | subagent 启动 | `session-start` hook | 为 subagent 加载上下文 |
| **SessionEnd** | 会话结束 | `forge cleanup` | 清理运行时状态 |
| **SubagentStop** | subagent 停止 | `forge cleanup` | 清理 subagent 状态 |
| **Stop** | Claude 停止响应时 | `forge quality-gate` + `forge feature complete --if-done` | 全部完成后的最终验证 |

### all-completed Hook

当所有任务完成，Claude 停止响应时作为最终安全网触发：

```
1. Quality Gate（three-phase pipeline，项目级，无 scope）：
   Phase 1 (NonBreakingGateSequence): just compile → just fmt → just lint
   Phase 2: just unit-test (retry-once)
   Phase 3: surface-aware test regression (dev → probe → test → teardown)
```

任何一步失败都会报告问题。

---

## Manifest 生命周期

`manifest.md` 是 Feature 的单一入口，由各 skill 自动维护。

### 状态流转

```
(none) ──→ prd ──→ design ──→ tasks ──→ in-progress ──→ completed
  │          │         │         │           │
  │      /write-prd    │    /breakdown-  首次 forge task claim
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
| `in-progress` | 执行中 | 首次 `forge task claim` |
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
│   ├── records/            # 执行记录（forge task submit 生成）
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
├── decisions/               # 技术决策（/learn）
├── lessons/                 # 经验教训（/learn）
└── sitemap/sitemap.json     # 页面元素映射（/gen-web-sitemap）
```

### 测试目录

```
tests/<surface-key>/
├── helpers.ts               # 共享工具函数
├── features/<slug>/         # 功能级测试（staging）
│   ├── *.spec.ts            #   测试脚本
│   └── results/             #   执行结果
└── <module>/                # 回归测试（通过 /run-tests 标签晋升）
    └── *.spec.ts            #   按功能模块组织
```

> Web surface 的测试目录下额外包含 `playwright.config.ts`（Playwright 配置）。测试类型命名遵循 Surface → Test Type 映射（见 [test-type-model.md](../../plugins/forge/skills/test-guide/references/test-type-model.md)）。

### 生成规则

| 目录 | 生成方式 | 是否提交 |
|------|---------|---------|
| `process/` | forge CLI 运行时 | **否**（gitignore） |
| `testing/` | `/gen-journeys` + `/gen-contracts` | 是 |
| `tests/<surface-key>/features/` | `/gen-test-scripts` | 是 |
| `tests/<surface-key>/` (根级) | `/run-tests` (tag promotion) | 是 |
| `records/` | `forge task submit` | 是 |
| `specs/` | `/consolidate-specs` | 是（用户确认后） |

---

## v3.0.0 子系统

v3.0.0 新增的辅助子系统，扩展了 Forge 的环境感知、知识管理和质量保证能力。

### Surface Detection

项目 surface（api/web/cli/tui/mobile）自动检测机制。`forge surfaces detect` 扫描项目目录结构和依赖文件，识别测试 surface 类型，结果用于 gen-journeys 和 gen-test-scripts 的 Convention 路由。实现在 forge CLI (`forge-cli/pkg/forgeconfig/detect_surface.go`)，非 skill 组件。

- 相关规则：[surface-api.md](../../plugins/forge/skills/gen-journeys/rules/surface-api.md) | [surface-web.md](../../plugins/forge/skills/gen-journeys/rules/surface-web.md) | [surface-cli.md](../../plugins/forge/skills/gen-journeys/rules/surface-cli.md) | [surface-tui.md](../../plugins/forge/skills/gen-journeys/rules/surface-tui.md) | [surface-mobile.md](../../plugins/forge/skills/gen-journeys/rules/surface-mobile.md)

### Worktree

Git worktree 隔离开发环境管理。`forge worktree` 命令组提供 start/list/remove/resume/push/status 子命令，支持功能分支的物理隔离——每个 worktree 拥有独立工作目录和分支，避免多任务间的文件冲突。实现在 forge CLI (`forge-cli/internal/cmd/worktree/`)，非 skill 组件。

### Convention

测试框架 Convention 文件生成系统。`/test-guide` 驱动项目测试框架自动检测（文件信号 + 依赖分析），生成 `docs/conventions/testing/<scope>.md` Convention 文件。Convention 文件定义测试发现、结构、断言模式和标签规范，供 gen-test-scripts 和 run-tests 消费，解耦 Forge 与具体测试框架。

- SKILL.md: [test-guide](../../plugins/forge/skills/test-guide/SKILL.md)

### Forensic

Agent 偏差溯源分析。搜索 JSONL 会话历史，提取思维链和工具调用序列，与 SKILL.md 定义的行为规范比对，定位 agent 偏离预期的根因。适用于多 session 重复偏差诊断，不用于单 session 事后分析（使用 `/learn` 替代）。

- SKILL.md: [forensic](../../plugins/forge/skills/forensic/SKILL.md)

### Deep Research

技术/产品深度调研。从主题名到结构化研究报告——自适应多源调查、交叉引用、上下文关联，产出可执行洞察。支持单技术深度分析和多候选方案对比两种模式。纯文档产出，不执行任何代码变更。

- SKILL.md: [deep-research](../../plugins/forge/skills/deep-research/SKILL.md)

### Clean Code

代码质量精炼。在限定 scope（git diff / 指定路径 / 全功能范围）内应用五项精炼原则，仅改变代码表达方式，不改变行为。可附带 Quality Gate（compile + fmt + lint）验证修改安全性。支持 standalone 调用和 pipeline 任务 `T-clean-code-1`。

- SKILL.md: [clean-code](../../plugins/forge/skills/clean-code/SKILL.md)

### Extract Design MD

视觉风格提取。从 web/mobile/tui 应用中自动提取视觉语言，生成 forge 兼容的 `DESIGN.md` 供 `/ui-design` skill 消费。产出设计令牌（颜色、字体、间距）和组件模式，桥接现有产品到 Forge UI 设计流程。

- SKILL.md: [extract-design-md](../../plugins/forge/skills/extract-design-md/SKILL.md)

### Learn

统一知识积累入口。合并 `/record-decision` 和 `/learn-lesson` 的功能，从单一入口捕获决策、经验、惯例和业务规则。核心原则：先写入后审核——立即持久化知识，用户在最终报告中审核和修正。

- SKILL.md: [learn](../../plugins/forge/skills/learn/SKILL.md)

### Test Guide

测试 Convention 文件引导生成。见上方 Convention 子系统。MANUAL-ONLY 技能，用户显式调用时引导完成框架检测、Convention 草稿生成、审核反馈循环和文件写入的完整流程。

- SKILL.md: [test-guide](../../plugins/forge/skills/test-guide/SKILL.md)

### Init Justfile

项目 Justfile 初始化。为用户项目生成标准化的 `justfile`，包含 compile、fmt、lint、unit-test、test、probe、dev、teardown 等 Forge Quality Gate 所需的 recipe。检测项目语言和构建工具链，生成对应的 recipe 模板。

- SKILL.md: [init-justfile](../../plugins/forge/skills/init-justfile/SKILL.md)
