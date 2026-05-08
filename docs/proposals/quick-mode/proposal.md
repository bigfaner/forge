---
created: 2026-05-09
author: faner
status: Draft
---

# Proposal: Forge 快速模式

## Problem

Forge 的完整工作流是一条重量级流水线：

```
brainstorm → write-prd → eval-prd → [ui-design → eval-ui] → tech-design → eval-design → breakdown-tasks → execute → T-test 1-5
```

对于 1-2 小时的短平快开发场景，这条流水线存在严重的投入产出失衡：

1. **文档开销远超实现成本**——一个半小时的改动，可能要花同样甚至更多时间在 PRD、设计、评分上
2. **环节串行等待**——每个 eval 环节至少一轮 adversarial iteration，用户必须多次确认
3. **测试流水线过重**——T-test 1-5 完整流水线（test-cases → eval → scripts → run → graduate → regression → consolidate）对小改动来说是过度工程

### Evidence

实际使用中，以下场景频繁出现但被迫走完整流程：

- 新增一个 CLI 命令
- 修改现有 skill 的模板格式
- 调整 task-cli 的输出格式
- 为现有功能增加一个配置项

这些任务的核心特征：**需求明确、范围有限、1-4 个可执行任务即可完成**。

## Proposed Solution

新增 `/quick` 命令，提供精简的快速模式流水线。

### 流程

```
/quick [--no-test]
  ├─ 1. /brainstorm → proposal.md（完整交互，无变化）
  ├─ 2. 用户确认 proposal ✋（暂停等待用户确认）
  ├─ 3. /quick-tasks → 读取 proposal.md
  │     ├─ 任务文件 + index.json
  │     ├─ T-quick-1 生成 test-cases.md    ← --no-test 时跳过
  │     ├─ T-quick-2 生成冒烟测试脚本      ← --no-test 时跳过
  │     ├─ T-quick-3 执行冒烟测试          ← --no-test 时跳过
  │     ├─ T-quick-4 graduate-tests        ← --no-test 时跳过
  │     ├─ T-quick-5 verify-regression     ← --no-test 时跳过
  │     └─ 简化版 manifest.md
  └─ 4. /run-tasks → 执行所有任务
```

### 与完整模式的对比

| 维度 | 完整模式 | 快速模式 |
|------|---------|---------|
| 入口 | 分步调用各 skill | `/quick` 一键 |
| 需求文档 | proposal → PRD → 设计 | proposal only |
| 评分环节 | eval-prd, eval-ui, eval-design, eval-test-cases | 无 |
| 任务拆分 | `/breakdown-tasks`（从 PRD+设计） | `/quick-tasks`（从 proposal） |
| 测试流水线 | T-test-1~5（gen-sitemap + gen-test-cases → eval → scripts → run → graduate → regression → consolidate） | T-quick-1~5（gen-test-cases → scripts → run → graduate → regression） |
| 质量门 | 完整 | 完整（一致） |
| manifest | 完整（Documents + Traceability） | 简化版（proposal + 任务列表） |
| 目录结构 | prd/ + ui/ + design/ + testing/ + tasks/ + specs/ | manifest.md + tasks/ + testing/ |
| 测试格式 | Playwright TypeScript, `tests/e2e/features/<slug>/` | 相同 |
| 测试毕业 | 支持 | 支持（流程一致） |
| 升级到完整模式 | — | 不支持 |

### 核心设计决策

#### 1. brainstorm 保持原样

`/quick` 内部调用完整 `/brainstorm`，输出标准 proposal.md。快速模式的用户脑子里已有大致想法，但 proposal 作为唯一的"需求源头"（后续没有 PRD 和设计来纠错），保留完整结构确保信息充分。

#### 2. 用户确认点

brainstorm 完成后暂停，展示 proposal 摘要让用户确认。因为 proposal 是整个快速模式唯一的输入文档，方向偏了后续任务全白费。确认成本几秒，避免几十分钟返工。

#### 3. 新增 `/quick-tasks` skill

从 proposal.md 直接拆分可执行任务，跳过 PRD 和设计环节。

**输入**：proposal.md
**输出**：
- 任务文件（`tasks/*.md`）+ `index.json`（格式兼容 `/run-tasks`）
- 每个任务文件包含涉及的文件清单（新建 / 修改 / 删除）
- T-quick-1~5 测试任务（除非 `--no-test`）
- 简化版 `manifest.md`

**与 `/breakdown-tasks` 的区别**：

| | `/breakdown-tasks` | `/quick-tasks` |
|---|---|---|
| 输入 | PRD + 设计文档 | proposal.md |
| 任务粒度 | 精细，phase 分组 + gate | 适中，扁平列表 |
| 条件标签 | `<HAS_UI>` `<NO_UI>` `<HAS_DB>` 等 | 无 |
| phase/gate | 有 | 无 |
| 测试任务 | T-test-1~5（含 gen-sitemap + eval） | T-quick-1~5（跳过 gen-sitemap + eval） |
| index.json | 完整格式（含 prd + design 引用） | 简化但兼容（proposal 引用替代 prd/design） |

#### 4. 任务模板对齐 breakdown-tasks

quick-tasks 的任务文件模板基于 breakdown-tasks 的 `templates/task.md`，保持一致的结构：

```markdown
---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "{{PRIORITY}}"
estimated_time: "{{ESTIMATED_TIME}}"
dependencies: [{{DEPENDENCIES}}]
status: pending
breaking: false
---

# {{ID}}: {{TITLE}}

## Description
{{DESCRIPTION}}

## Affected Files

### Create
| File | Description |
|------|-------------|
| {{NEW_FILES}} |

### Modify
| File | Changes |
|------|---------|
| {{MODIFIED_FILES}} |

### Delete
| File | Reason |
|------|--------|
| {{DELETED_FILES}} |

## Reference Files
{{REFERENCE_FILES}}

## Acceptance Criteria
{{ACCEPTANCE_CRITERIA}}

## User Stories
{{USER_STORIES}}

## Implementation Notes
{{NOTES}}
```

**与 breakdown-tasks `task.md` 模板的差异**：

| 差异点 | breakdown-tasks | quick-tasks |
|--------|----------------|-------------|
| 新增 `## Affected Files` 章节 | 无 | 有（Create/Modify/Delete 三类） |
| `prd`/`design` 字段 | index.json 引用 PRD + 设计路径 | index.json 引用 proposal 路径 |

quick-tasks 从 proposal 的 scope 和 solution 描述中推导文件清单，写入 Affected Files 章节。这帮助 agent 执行时精确定位文件范围，也方便用户快速评估任务影响面。

#### 5. 测试任务拆分

快速模式的测试任务对应完整模式 T-test 的子集，跳过 gen-sitemap 和 eval-test-cases，但保留 test-cases 中间产物以保证脚本生成质量：

| 快速模式 | 对应完整模式 | 内容 |
|---------|------------|------|
| T-quick-1 | T-test-1 的一半 | 从 proposal 生成 test-cases.md（无 gen-sitemap） |
| T-quick-2 | T-test-2 | 从 test-cases 生成冒烟测试脚本（Playwright TypeScript） |
| T-quick-3 | T-test-3 | 执行冒烟测试 |
| T-quick-4 | T-test-4 | graduate-tests |
| T-quick-5 | T-test-4.5 | verify-regression |

跳过的环节：

| 跳过 | 原因 |
|------|------|
| T-test-1 中的 gen-sitemap | 快速模式无需 sitemap 依赖 |
| T-test-1b (eval-test-cases) | test-cases 是过渡产物，不需要评分 |
| T-test-5 (consolidate-specs) | 快速模式不提取 project-level specs |

**为什么保留 test-cases 中间产物**：

T-quick-2 直接从 proposal 跳到脚本生成会降低测试质量——agent 缺少"测什么"的结构化思考过程。将 test-cases 拆为独立任务，让 agent 先从 proposal 提取测试场景（结构化），再基于 test-cases 生成脚本（有依据），生成质量接近完整模式的 T-test-2。eval-test-cases 被跳过，因为 test-cases 是过渡产物而非长期文档，快速模式接受这个权衡。

#### 6. 简化版 manifest

```markdown
---
name: <feature-name>
status: quick-mode
---

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/<slug>/proposal.md |
| Test Cases | testing/test-cases.md |

## Tasks

| ID | Title | Status |
|----|-------|--------|
| 1.1 | ... | pending |
| T-quick-1 | Generate test cases | pending |
| T-quick-2 | Generate smoke test scripts | pending |
| T-quick-3 | Run smoke tests | pending |
| T-quick-4 | Graduate tests | pending |
| T-quick-5 | Verify regression | pending |
```

不包含 Documents Traceability 表（无 PRD/设计可追踪）。

#### 7. 目录结构

```
docs/features/<slug>/
  manifest.md              # 简化版
  testing/
    test-cases.md          # T-quick-1 生成
  tasks/
    index.json
    1.1-*.md               # 实现任务
    T-quick-1.md           # 生成 test-cases
    T-quick-2.md           # 生成冒烟测试脚本
    T-quick-3.md           # 执行冒烟测试
    T-quick-4.md           # 毕业
    T-quick-5.md           # regression 验证
    records/
    process/

# proposal 保持在原位（brainstorm 已有行为）
docs/proposals/<slug>/proposal.md

# 测试脚本（与完整模式一致）
tests/e2e/features/<slug>/
  smoke.spec.ts
```

#### 8. `/quick-tasks` 独立可用

注册在 SKILLS.md，用户可单独调用。场景：已有 proposal，只想快速拆任务执行，不需要重跑 brainstorm。

#### 9. 质量门保持一致

每个任务执行后照跑 `just compile → just fmt → just lint`。全部任务完成后跑 `just test`。与完整模式完全一致，零额外学习成本。

#### 10. 复用 `/run-tasks`

quick-tasks 产出的 index.json 格式兼容 `/run-tasks`，不创建新的执行引擎。

## Quality Gate Summary

```
/quick
  ├─ brainstorm → proposal.md（标准质量）
  ├─ 用户确认 ✋
  ├─ quick-tasks → 任务拆分 + T-quick-1~5
  └─ run-tasks
       ├─ 每任务后: compile + fmt + lint
       ├─ T-quick-1: 从 proposal 生成 test-cases.md
       ├─ T-quick-2: 基于 test-cases 生成冒烟测试（Playwright TS）
       ├─ T-quick-3: 执行冒烟测试
       ├─ T-quick-4: 毕业 → tests/e2e/<target>/
       ├─ T-quick-5: 全量 regression
       └─ all-completed: just test + e2e regression
```

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| 完整流程不变 | 零改动 | 短任务开销过大，用户绕过流程导致无记录 | Rejected: 不解决根本问题 |
| 只跳过 eval 环节 | 减少 3 轮评分 | PRD + 设计仍然过重 | Rejected: 杯水车薪 |
| 合并 brainstorm + 拆任务为一个 skill | 步骤更少 | 丧失 quick-tasks 独立调用能力；brainstorm 已有稳定逻辑不想改 | Rejected: 灵活性差 |
| T-quick-1 内部两阶段（不拆任务） | 任务数更少 | 失败时无法精确定位是 test-cases 生成问题还是脚本生成问题 | Rejected: 故障隔离差 |
| **T-quick-1~5 拆分 + test-cases 中间产物** | 测试质量接近完整模式；任务级别故障隔离；与完整模式测试流程对齐 | 新增 2 个 skill + 1 个 manifest 模板 + testing/ 目录 | **Selected** |

## Scope

### In Scope

- `/quick` 命令 skill（入口，串联 brainstorm → quick-tasks → run-tasks）
- `/quick-tasks` skill（独立注册，从 proposal 拆任务）
- 简化版 manifest 模板
- `--no-test` 参数支持（`/quick` 和 `/quick-tasks` 均支持）
- SKILLS.md 注册新 skill
- guide.md 更新（新增快速模式流程图）

### Out of Scope

- 快速模式升级为完整模式（不支持）
- 修改现有 `/brainstorm`（完全复用，无变化）
- 修改现有 `/run-tasks`（完全复用，index.json 兼容即可）
- 修改 task-cli（index.json 格式不变）
- 修改 hooks / quality gate（行为一致，无变化）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| proposal 质量不足导致任务拆分方向偏 | Medium | High——没有 PRD/设计纠错，返工成本高 | brainstorm 后强制用户确认点；brainstorm 保持完整交互质量 |
| T-quick-1 生成的 test-cases 质量低于完整模式（缺少 PRD 作为输入） | Medium | Medium——test-cases 覆盖不全面，下游脚本质量受影响 | agent 从 proposal 的目标、范围、约束提取测试场景，覆盖核心路径；test-cases 作为中间产物可在执行时迭代修正 |
| 快速模式滥用——本该走完整模式的大型 feature 也用快速模式 | Medium | Medium——缺乏设计文档导致实现质量下降 | 在 guide.md 中明确适用场景（1-2 小时、1-4 个任务）；快速模式 manifest 标记 `status: quick-mode` 可追溯 |
| index.json 格式不完全兼容 `/run-tasks` | Low | High——执行阶段失败 | quick-tasks 模板对齐 breakdown-tasks 的 index.json schema，只减不增 |

## Success Criteria

- [ ] `/quick` 命令可一键完成 brainstorm → 确认 → 拆任务 → 执行全流程
- [ ] `/quick-tasks` 可独立调用，从 proposal.md 生成任务文件 + index.json
- [ ] 每个非测试任务文件包含 `## Affected Files` 章节（Create/Modify/Delete），路径和说明准确
- [ ] 任务文件模板与 breakdown-tasks `task.md` 保持一致结构（frontmatter + Description + Reference Files + Acceptance Criteria + User Stories + Implementation Notes）
- [ ] `--no-test` 参数生效时不生成 T-quick-1~5 任务
- [ ] 默认模式生成 T-quick-1~5 任务，且 `/run-tasks` 可正常执行
- [ ] 生成的 manifest.md 包含 proposal 路径和任务列表，不包含 Documents Traceability 表
- [ ] 目录只创建 `manifest.md` + `tasks/` + `testing/`（T-quick-1 生成 test-cases 时创建），不创建其他空目录
- [ ] T-quick-1 从 proposal 生成 test-cases.md，位于 `docs/features/<slug>/testing/`
- [ ] T-quick-2 基于 test-cases 生成冒烟测试，为 Playwright TypeScript，位于 `tests/e2e/features/<slug>/`
- [ ] T-quick-4 可成功毕业测试到回归套件
- [ ] 质量门行为与完整模式一致（compile + fmt + lint + test）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
