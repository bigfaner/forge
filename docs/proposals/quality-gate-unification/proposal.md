---
title: Quality Gate Unification & Guide Refinement
status: approved
date: 2026-05-03
---

# Quality Gate Unification & Guide Refinement

## 目标

1. 将 justfile 标准命令（compile→fmt→lint→test）作为强制性质量门禁嵌入 forge 工作流
2. 精炼 guide.md：目录规范、Testing Lifecycle、Quality Gate Protocol
3. 新增 consolidate-specs 技能，解决知识分散问题

## 决策记录

### D1. 门禁位置：任务级

| 执行入口 | 文件 | 门禁卡点 |
|---------|------|---------|
| `/execute-task` | `commands/execute-task.md` | verify 步骤 |
| `task-executor` agent | `agents/task-executor.md` | verify 步骤 |
| `/fix-bug` | `commands/fix-bug.md` | verify 步骤 |

### D2. 门禁序列

```
just compile [scope] → just fmt [scope] → just lint [scope] → just test [scope]
```

- 严格顺序执行，前一步失败则中止
- scope-aware：遵循 Scope Resolution Protocol
- `just fmt` 用自动修复模式（非 --check）

### D3. Agent 层失败处理

| 失败环节 | 行为 |
|---------|------|
| `compile` 失败 | 触发 error-fixer 修复编译错误 |
| `fmt` 失败 | 标记任务为 `blocked`（fmt 自动修复仍失败说明工具链/模板有问题） |
| `lint` 失败 | task-executor 自行修复，最多重试 **1 次**，仍失败则 `blocked` |
| `test` 失败 | 触发 error-fixer |

### D4. 两层防御

| 层 | 位置 | 职责 |
|----|------|------|
| **Agent 层**（verify 步骤） | 早期反馈 + 自动修复 | Agent 跑门禁，发现问题当场修 |
| **CLI 层**（task record） | 硬执行兜底 | validation 阶段执行门禁，任一失败拒绝 record |

CLI 层实现：`task record` 的 validation 阶段检测 justfile 和 recipe 存在性，执行 `just compile && just fmt && just lint && just test`。失败时输出精要错误（哪步失败 + 关键错误信息），可通过 `--force` 跳过。

### D5. all-completed hook 兜底

`all-completed`（Stop hook）补上 `fmt` 和 `lint`，作为全局兜底：

```
just compile → just fmt → just lint → just test → just e2e-setup → just test-e2e
```

### D6. 目录规范按工作流阶段组织

guide.md 的目录规范改为按 skill 输出归类，agent 可直接根据当前执行上下文确定文件存放位置。

非 skill 产出文档归位规则：

```
docs/
  business-rules/       — 跨功能业务规则（由 consolidate-specs 整合）
  conventions/          — 技术规范（编码规范、API 规范、命名规范）
  reference/            — 系统规范（环境要求、部署流程、技术选型）
  decisions/            — 技术决策（已有）
  lessons/              — 经验教训（已有）
  ARCHITECTURE.md       — 架构文档（已有）
```

业务规则按领域组织（如 `auth.md`、`payment.md`）。

### D7. Testing Lifecycle

三层测试体系：

```
Unit Tests (just test)     — 任务级，每次 verify 必跑
E2E Tests (just test-e2e)  — 功能级，breakdown-tasks 末尾标准任务
Regression Suite           — 毕业，tests/e2e/ 目录
```

### D8. consolidate-specs 技能

新增标准任务 T-test-6，排在 graduate-tests 之后：

```
T-test-6: /consolidate-specs
  1. 扫描 docs/features/<slug>/prd/ 提取业务规则 → specs/biz-specs.md
  2. 扫描 docs/features/<slug>/design/tech-design.md 提取技术规格 → specs/tech-specs.md
  3. 暂停，展示预览给用户
  4. 用户确认后：
     biz-specs → docs/business-rules/<domain>.md
     tech-specs → docs/conventions/<topic>.md
```

文件结构：
```
docs/features/<slug>/
  specs/
    biz-specs.md     — 从 PRD 提取的业务规则
    tech-specs.md    — 从 design 提取的技术规格
```

## 变更范围

| # | 文件 | 类型 | 变更 |
|---|------|------|------|
| 1 | `plugins/forge/hooks/guide.md` | 重构 | 目录规范 + Testing Lifecycle + Quality Gate Protocol |
| 2 | `plugins/forge/commands/execute-task.md` | 修改 | verify 步骤加门禁序列 |
| 3 | `plugins/forge/agents/task-executor.md` | 修改 | verify 步骤加门禁序列 + lint 自修逻辑 |
| 4 | `plugins/forge/commands/fix-bug.md` | 修改 | verify 步骤加门禁序列 |
| 5 | `task-cli/internal/cmd/all_completed.go` | 修改 | 加 fmt/lint 到执行序列 |
| 6 | `task-cli/internal/cmd/record.go` | 修改 | validation 加门禁 pre-check + 精要错误输出 |
| 7 | `plugins/forge/skills/breakdown-tasks/SKILL.md` | 修改 | 追加 T-test-6 + specs/ 目录规则 |
| 8 | `plugins/forge/skills/consolidate-specs/SKILL.md` | 新增 | 提取 → 预览 → 用户确认 → 整合 |
| 9 | `plugins/forge/SKILLS.md` | 修改 | 注册新 skill |
