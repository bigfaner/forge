---
title: "Tech Design Skill 改进：重命名与决策归档"
date: 2026-04-22
status: draft
---

## Problem

当前 `design-tech` skill 存在两个问题：

1. **命名不一致**：skill 名称 `design-tech` 与输出文件 `tech-design.md` 方向相反，容易混淆。流程中其他 skill 命名与输出一致（`write-prd` → `prd-spec.md`，`ui-design` → `ui-design.md`）。
2. **决策分散**：技术设计中的关键决策（架构选型、数据模型、错误处理策略等）仅存在于各 feature 的 `tech-design.md` 中，无法跨 feature 追溯。当多个 feature 涉及相同领域时，难以快速定位历史决策依据。

**紧迫性**：随着 feature 数量增长，决策追溯成本持续上升。越早建立集中归档，历史决策越完整。

## Solution

### 1. 重命名 `design-tech` → `tech-design`

- 将 skill 目录从 `plugins/zcode/skills/design-tech/` 重命名为 `plugins/zcode/skills/tech-design/`
- 更新 SKILL.md frontmatter 中的 `name` 字段
- 更新所有引用该名称的文件（hooks guide、exploration examples、plugin 文档等）

调用方式从 `/zcode:design-tech` 变为 `/zcode:tech-design`。

### 2. 决策归档到 `docs/decisions/`

#### 目录结构

```
docs/decisions/
├── manifest.md              # 决策索引，类似 feature manifest
├── architecture.md          # 架构决策
├── interface.md             # 接口设计决策
├── data-model.md            # 数据模型决策
├── dependencies.md          # 依赖选择决策
├── error-handling.md        # 错误处理决策
├── testing.md               # 测试策略决策
├── security.md              # 安全考量决策
└── local-dev-deployment.md  # 本地开发与线上部署决策
```

#### manifest.md

决策目录的单一入口索引文件，由 skill 自动维护：

```markdown
---
updated: "{{DATE}}"
---

# Decisions Index

## Categories

| Category | File | Decisions | Last Updated |
|----------|------|-----------|-------------|
| Architecture | architecture.md | 0 | - |
| Interface | interface.md | 0 | - |
| Data Model | data-model.md | 0 | - |
| Dependencies | dependencies.md | 0 | - |
| Error Handling | error-handling.md | 0 | - |
| Testing | testing.md | 0 | - |
| Security | security.md | 0 | - |
| Local Dev & Deployment | local-dev-deployment.md | 0 | - |

## Recent Decisions

| Date | Feature | Category | Decision | Source |
|------|---------|----------|----------|--------|
```

每次归档决策时，同步更新 manifest.md 的 Decisions 计数、Last Updated 时间和 Recent Decisions 表。

#### 单条决策记录格式

追加到对应类型文件末尾的 markdown 表格行：

```markdown
| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
| 2026-04-22 | feat-xxx | 采用事件驱动架构 | 需要解耦模块间通信，支持异步处理 | feat-xxx/tech-design.md §Architecture |
```

### 3. 决策归档为可选步骤

tech-design 流程中，决策归档不是必须步骤：

- **有决策时**：tech-design 文档获用户批准后，展示编号的候选决策列表，用户输入逗号分隔的编号选择要归档的条目（或 `all` / `none`），确认后写入 `docs/decisions/`
- **无决策时**：跳过归档步骤，直接进入 manifest 更新

**归档确认交互示例**：

```
以下决策被标记为关键决策，建议归档：

  [1] 采用事件驱动架构（Architecture）
  [2] 使用 SQLite 作为本地缓存存储（Data Model）
  [3] 选择 Vitest 而非 Jest 作为测试框架（Dependencies）

输入要归档的编号（逗号分隔），或 all / none：1,3

→ 已归档 2 条决策到 architecture.md 和 dependencies.md
```

用户也可在确认时编辑：输入 `edit:<编号>` 可重新编辑该条决策的 Decision 或 Rationale 字段后再归档。

- **判定标准**：由 AI 在 draft design 阶段识别，标记为"关键决策"的条目才进入归档候选

### 4. 决策记录逻辑独立为 reference

将决策提取和记录的逻辑从 tech-design 的 SKILL.md 中分离，放在独立文件：

```
plugins/zcode/skills/tech-design/
├── SKILL.md
├── references/
│   └── decision-logging.md    # 决策提取与归档的完整流程说明
├── templates/
│   ├── tech-design.md
│   ├── api-handbook.md
│   ├── manifest-update-design.md
│   └── decision-entry.md      # 单条决策的模板
└── examples/
    ├── ask-question.md
    └── exploration.md
```

tech-design 的 SKILL.md 在决策归档步骤中引用 `references/decision-logging.md`，而非内联所有逻辑。新增的 `/zcode:record-decision` skill 也引用同一份 reference。

### 5. 新增 `/zcode:record-decision` 命令

提供独立的 slash command，允许用户在任意阶段主动记录技术决策：

**调用方式**：`/zcode:record-decision`

**流程**：
1. 通过 AskUserQuestion 逐项收集决策信息（4 轮交互）：
   - **类型**：展示编号列表（1=Architecture, 2=Interface, …, 8=Local Dev & Deployment），用户输入编号 1-8
   - **决策描述**：单行文本，用户输入决策内容
   - **决策理由**：单行文本，用户输入理由
   - **关联 feature**：用户输入 feature 名称（如 `feat-log-decisions`）
2. 自动补充 Date 和 Source
3. 写入对应 `docs/decisions/<type>.md`
4. 更新 `docs/decisions/manifest.md`

**交互示例**：

```
=== 记录技术决策 ===

选择决策类型：
  [1] Architecture  [2] Interface      [3] Data Model
  [4] Dependencies   [5] Error Handling [6] Testing
  [7] Security       [8] Local Dev & Deployment

输入编号 (1-8)：1

决策描述（一句话）：采用事件驱动架构解耦模块间通信

决策理由（一句话）：需要支持异步处理和模块独立演进

关联 feature：feat-log-decisions

→ 已写入 architecture.md，manifest.md 已更新
```

**适用场景**：
- 实现过程中发现需要记录的非 design 阶段决策
- 对历史决策的补充或修正
- brainstorm / PRD 阶段产生的重要技术约束决策

### 6. 移除 DECISIONS.md 引用

更新 hooks guide（`plugins/zcode/hooks/guide.md`）和 exploration 示例中对 `docs/DECISIONS.md` 的引用，改为 `docs/decisions/` 目录。

## Alternatives

### A. 模板化决策系统（未选择）

为每条决策创建独立 `.md` 文件（类似 ADR），包含 Context/Decision/Consequences/Alternatives 结构化字段，并在 tech-design 流程中增加决策评审步骤。

- 优点：单条决策可独立引用（如 `decisions/adr-001.md`），字段完整，适合合规要求高的项目
- 代价：每条决策需创建文件 + 填写 4-5 个字段，单次归档交互时间约 2 分钟；10 个 feature 后产生 30-50 个文件，目录膨胀快

### B. 仅重命名（未选择）

只做重命名，不加决策归档。现有 tech-design.md 的 "Alternatives Considered" 部分已覆盖部分需求。

- 优点：零新增复杂度，改动局限于目录名和引用
- 代价：决策分散在各 feature 目录，无法跨 feature 追溯；搜索某技术选型需遍历所有 feature 的 tech-design.md

### C. 不做任何变更（未选择）

- 代价：命名混乱持续，决策追溯成本随 feature 增长而上升

### Chosen Approach Rationale

本方案在 **粒度**（per-decision 文件 vs. 单文件 vs. 分类文件）和 **流程开销**（结构化模板 vs. 表格行）两个轴上取中间点：

1. **为什么 8 个分类文件而非 per-decision 文件（方案 A）**：per-decision 文件适合正式 ADR 流程，但本项目是个人/小团队工具链，决策频率高且以实用性为主。8 个分类文件将同类决策聚合，单文件内用表格行组织，检索时打开一个文件即可看到该领域全部历史决策，无需在数十个文件间跳转。

2. **为什么分类文件而非单一 DECISIONS.md（方案 B 的延伸）**：单一文件在 20+ 条决策后检索效率下降。按类型拆分使每个文件预期控制在 10-30 条，同时 manifest.md 提供跨类型的全局视图。

3. **为什么表格行而非结构化模板**：表格行只需 5 个字段（Date, Feature, Decision, Rationale, Source），单条记录约 1-2 行 markdown。结构化模板（Context/Decision/Consequences/Alternatives）字段更多，适合可审计场景，但对本项目的决策频率和团队规模来说增加了不必要的填写成本。

## Scope

### In Scope

| # | Item | Effort | Depends on | Phase |
|---|------|--------|------------|-------|
| 1 | 重命名 skill 目录 `design-tech/` → `tech-design/`，更新注册名和所有引用 | S | — | 1 |
| 2 | 创建 `docs/decisions/manifest.md` 索引文件和 8 个类型模板文件 | S | — | 1 |
| 3 | 创建 `references/decision-logging.md` 独立决策记录逻辑 | M | 2 | 2 |
| 4 | 创建 `templates/decision-entry.md` 决策条目模板 | S | 2 | 2 |
| 5 | tech-design 流程新增可选"决策归档"步骤（引用 reference 3） | M | 1, 3, 4 | 2 |
| 6 | 新增 `/zcode:record-decision` slash command skill | M | 3, 4 | 2 |
| 7 | 更新 hooks guide (`plugins/zcode/hooks/guide.md`) 和 exploration 示例 (`plugins/zcode/skills/tech-design/examples/exploration.md`) 中的引用 | S | 1 | 3 |
| 8 | 更新 `zcode/CLAUDE.md` 中的 skill 列表和 `plugins/zcode/SKILLS.md`（如存在）中的 skill 注册信息 | S | 5, 6 | 3 |

**Effort scale**: S = < 30 min, M = 30 min–2 h

**Phases**:
- **Phase 1**（基础设施）：目录重命名 + 模板文件创建，无依赖，可并行
- **Phase 2**（核心逻辑）：决策记录 reference、流程集成、新命令，依赖 Phase 1 的文件结构
- **Phase 3**（文档同步）：更新外部引用和文档列表，依赖 Phase 2 确定最终 skill 名称和路径

### Out of Scope

- `docs/DECISIONS.md` 的创建
- 下游 skill（eval-design、breakdown-tasks）的变更
- tech-design.md 模板的 "Alternatives Considered" 部分
- 自动合并冲突决策的检测机制
- 决策版本管理（撤销、修订历史）

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 归档表格随时间膨胀，难以检索 | Medium | M — 单文件超过 200 行后可读性下降 | 按类型分 8 个文件控制单文件规模；当任意类型文件超过 50 条时，reference 流程按年份自动分割为 `architecture-2026.md` 子文件 |
| manifest.md 与类型文件不同步 | Medium | M — 索引计数错误，Recent Decisions 表缺少条目 | 新增 `validate-manifest` CI 脚本：解析所有类型文件的实际条目数，与 manifest 计数比对，不一致则报错；脚本随项目 CI 运行 |
| `/zcode:record-decision` 使用频率低 | High | H — 用户仍依赖手动记录或跳过归档 | tech-design 流程结束后自动提示归档（当前设计）；若连续 10 个 feature 未产生任何归档记录，则将归档步骤从"可选"提升为"默认执行"，用户需主动输入 `skip` 跳过 |
| 重命名导致现有引用断裂 | Low | H — 旧文档中的 `/zcode:design-tech` 调用失效 | 全局搜索 `design-tech` 并更新所有引用（预估 5 处）；重命名后运行 `grep -r "design-tech" plugins/ docs/` 验证无残留引用 |
| 决策记录与 tech-design.md 内容重复 | Low | L — 同一决策在两处维护，更新时遗漏一处 | 表格记录为摘要级（Decision + Rationale 各一句话），tech-design.md 保留完整分析；字段设计上避免复制长文本 |

## Success Criteria

1. `/zcode:tech-design` 调用后，在 feature 目录下生成 `tech-design.md`，且流程中产生的关键决策能写入对应 `docs/decisions/<type>.md` 文件
2. tech-design 审批后：有决策时展示编号候选列表，用户选择后归档；无决策时跳过归档步骤——两种分支均可通过构造有/无决策的测试 feature 验证
3. `/zcode:record-decision` 独立调用后，输入 4 轮信息（类型编号、描述、理由、feature），对应类型文件新增 1 条表格行，manifest.md 计数 +1
4. `validate-manifest` CI 脚本通过：脚本解析所有类型文件的实际条目数，与 manifest.md 中 Decisions 列的计数一致，且 Recent Decisions 表包含最近 5 条记录
5. 决策记录包含 Date、Feature、Decision、Rationale、Source 五个字段
6. hooks guide 和 exploration 示例不再引用 `docs/DECISIONS.md`
7. 旧的 `/zcode:design-tech` 调用方式不再可用，`grep -r "design-tech" plugins/ docs/` 无残留结果
8. 决策记录逻辑在 `references/decision-logging.md` 中独立维护，tech-design 和 record-decision 共享
9. `zcode/CLAUDE.md` 的 skill 索引表包含 `tech-design` 条目；若 `plugins/zcode/SKILLS.md` 存在，则注册了 `tech-design` 和 `record-decision` 两个 skill
