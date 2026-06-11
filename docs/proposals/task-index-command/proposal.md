---
created: 2026-05-12
author: faner
status: Draft
---

# Proposal: CLI 任务提取逻辑 — `task index` 命令

## Problem

当前 `index.json` 完全由 AI agent 通过 skill prompt 生成（`breakdown-tasks`、`quick-tasks`）。AI 同时编写 task .md 文件和 index.json，充当"编译器"角色。这带来几个问题：

1. **AI 生成结构化 JSON 不稳定** — AI 容易在 index.json 中产生结构性错误（字段遗漏、依赖引用错误、type 赋值不一致），需要 `task validate` 兜底
2. **两份数据重复维护** — .md frontmatter 和 index.json 存在大量冗余字段，AI 必须保持两者同步
3. **skill prompt 膨胀** — breakdown-tasks 和 quick-tasks 的 prompt 中包含大量 JSON 模板填充指令，增加 AI 认知负担
4. **模板漂移** — 不同 feature 的测试任务内容因 AI 自由发挥而产生不一致（已观察到结构漂移现象）

### Evidence

- `task validate` 命令存在 10+ 条校验规则，多数针对 AI 生成 index.json 时易犯的错误
- 跨 feature 对比发现测试任务 .md 内容存在结构漂移（路径格式、fix-task 语法、描述详略不一致）
- .md frontmatter 已包含 id/title/priority/dependencies 等结构化字段，与 index.json 高度冗余

## Proposed Solution

新增 `task index --feature <slug> [--no-test]` 命令，将 index.json 的生成从 AI prompt 转移到 CLI 代码。

**核心原则：task .md 文件是 single source of truth，index.json 是纯派生产物。**

### 职责划分

```
AI (skills) 负责：
  - 编写 business task .md 文件（含 YAML frontmatter）
  - 编写 gate .md 文件
  - 编写 phase summary .md 文件
  - 更新 manifest.md

CLI (task index) 负责：
  - 解析 tasks/*.md 的 YAML frontmatter
  - 推断 type（InferType）、status、file、record
  - 检测模式（breakdown / quick）并生成测试任务 .md + index 条目
  - 合并到 index.json（幂等，保留运行时状态）
  - 自动校验生成的 index.json
```

### Frontmatter 最小集

AI 只需在 .md frontmatter 中写入：

| 字段 | 必填 | 说明 |
|------|------|------|
| `id` | 是 | 任务 ID，如 `1.1`、`3.gate` |
| `title` | 是 | 任务标题 |
| `priority` | 是 | P0 / P1 / P2 |
| `dependencies` | 是 | 依赖的任务 ID 列表 |
| `scope` | 是 | frontend / backend / all |
| `estimated_time` | 否 | 工时估算 |
| `breaking` | 否 | 是否触发全量测试 |
| `mainSession` | 否 | 是否必须在主会话执行 |
| `noTest` | 否 | 跳过质量门禁 |

### CLI 推断字段

| 字段 | 推断逻辑 |
|------|----------|
| `type` | InferType()，按 ID 模式匹配（如 `.gate` → `gate`，`T-test-*` → 对应 type） |
| `status` | 新任务 → `pending`；已存在 → 保留 index.json 中的值 |
| `file` | .md 文件名本身 |
| `record` | `records/<filename>` |

### Feature 元数据推断

CLI 从 `--feature <slug>` 参数 + 文件约定自动推断：

| 字段 | 推断规则 |
|------|----------|
| `feature` | `--feature` 参数值 |
| `prd` | 若 `prd/prd-spec.md` 存在则填入 |
| `design` | 若 `design/tech-design.md` 存在则填入 |
| `proposal` | 若 `docs/proposals/<slug>/proposal.md` 存在则填入 |
| `created` | 新 feature → 今天日期；已存在 → 保留原值 |

### 模式检测

| 条件 | 模式 | 测试任务集 |
|------|------|-----------|
| `prd/prd-spec.md` 存在 | breakdown | T-test-1 到 T-test-5（7 个） |
| `proposal.md` 存在 | quick | T-quick-1 到 T-quick-5（5 个） |
| `--no-test` 标志 | 任一 | 跳过测试任务生成 |

### 测试任务生成

CLI 从嵌入模板生成测试任务的 .md 文件和 index 条目，进行机械替换：

- `<slug>` → feature slug
- 依赖链：
  - breakdown: T-test-1 → `[最高阶段 gate]`，后续测试任务依赖前一个
  - quick: T-quick-1 → `[最大业务任务 ID]`，后续依赖前一个

### 合并语义（幂等）

`task index` 对已存在的 index.json 执行合并，而非覆盖：

| 字段来源 | 字段 |
|----------|------|
| .md frontmatter 覆盖 | id, title, priority, dependencies, scope, breaking, mainSession, noTest, estimatedTime |
| CLI 推断覆盖 | type, file, record |
| index.json 保留 | status, sourceTaskID, blockedReason |
| 新 .md 文件 | 作为新任务加入，status = pending |
| index.json 有但无 .md | 输出 WARNING，保留条目不变 |

### 架构集成

```
pkg/task/sync.go        — 共享合并逻辑（BuildIndex 函数）
internal/cmd/index.go   — task index 命令
internal/cmd/add.go     — task add 调用 BuildIndex 共享逻辑
pkg/task/testgen.go     — 测试任务模板嵌入 + 生成
```

- `task add` 创建 .md 文件后调用 `BuildIndex()`，保持单一写入路径
- `task index` 批量扫描后调用 `BuildIndex()`
- 自动在写入后执行 validate 逻辑，失败则输出错误

### 输出

命令执行后输出：
1. 合并结果摘要（新增 N 个、更新 M 个、保留 K 个）
2. 警告（孤立项、验证失败）
3. 提示 agent 更新 manifest.md

## Scope

### In Scope

- `task index` 命令实现
- 共享 `BuildIndex()` 合并逻辑
- 测试任务模板嵌入 + 机械替换
- `task add` 重构为调用共享逻辑
- 修改 `breakdown-tasks` 和 `quick-tasks` skill，移除 index.json 生成职责
- 自动 validate 集成

### Out of Scope

- manifest.md 自动更新（agent 职责）
- 现有 feature 的 index.json 迁移（增量兼容）
- InferType() 逻辑变更（保持现有规则）

## Impact

- **breakdown-tasks skill** — 移除 Step 5（assemble index.json）和 Step 4d（test task 生成），简化为"写 .md + 调用 `task index`"
- **quick-tasks skill** — 同上
- **task add** — 内部重构，行为不变
- **task validate** — 仍可独立调用，但 build-index 后自动执行
