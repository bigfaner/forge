---
created: "2026-05-28"
author: "faner"
status: Draft
---

# Proposal: Autogen Test Task Path References

## Problem

Auto-generated test pipeline tasks (`T-test-run`, `T-test-gen-scripts`, etc.) 生成内容过于简略，缺少 feature 级路径上下文。Subagent 执行任务时（Step 1: Read Task Definition），无法从 task .md 文件中获知 journeys 的目录位置，需要额外探索才能定位测试产物。

### Evidence

三层架构中路径信息的当前分布：

| 层 | 路径信息 | `FeatureSlug` 来源 |
|---|---------|-------------------|
| **Task embed 模板** | gen-journeys/gen-contracts 较完整；test-run/test-gen-scripts 几乎为空 | CLI 从 `docs/features/<slug>/` 目录路径填充 |
| **Prompt 模板** | 无路径信息；`FeatureSlug` 在 context 中声明但**未渲染到输出** | Dispatcher 从 `index.json` 的 `feature` 字段传入 |
| **Skill** | 唯一包含完整发现逻辑的位置 | 运行时从 state/CLI/路径解析获取 |

Prompt 模板当前输出：

```
TASK_ID: T-test-run
TASK_FILE: /path/to/docs/features/my-feature/tasks/run-test.md
SURFACE_KEY: api

缺失: FEATURE_SLUG  ← 声明了但未渲染
```

Agent 无法直接获取 slug，需要从 TASK_FILE 路径反推。run-tests skill Step 1.5 列出 3 种 slug 获取源（模板变量 / `forge feature status` / 路径解析），恰恰说明这个问题。

### Urgency

不影响正确性（skill 能动态发现路径），但降低了 agent 执行效率。三层之间缺乏联动，路径发现逻辑仅在 skill 中，task file 和 prompt 未提供有效上下文。

## Proposed Solution

三层联动：embed 模板补充发现命令 + prompt 模板渲染 FeatureSlug + skill 保留完整逻辑。

### 改动 1：Embed 模板补充 `## Feature Paths`

6 个测试流水线 embed 模板统一添加两类发现命令：

**发现 journeys**（与 `run-tests` skill Step 1.5 一致）：

```bash
ls docs/features/{{.FeatureSlug}}/testing/
```

**发现 contracts**（与 `gen-test-scripts` skill 一致）：

```bash
ls docs/features/{{.FeatureSlug}}/testing/<journey>/contracts/
```

> `<journey>` 由上一步 `ls` 结果获得。

格式统一为：

```
## Feature Paths

Discover journeys:

```bash
ls docs/features/{{.FeatureSlug}}/testing/
```

Discover contracts (per journey):

```bash
ls docs/features/{{.FeatureSlug}}/testing/<journey>/contracts/
```
```

`{{.FeatureSlug}}` 由 CLI 在 `forge task index` 时从 `docs/features/<slug>/` 目录路径填充，不依赖 `.forge/state.json` 或分支名，生成时即确定。

**FeatureSlug 来源一致性**：Embed 模板和 Prompt 模板虽然取值路径不同（前者从 CLI 目录路径填充，后者从 `index.json` 的 `feature` 字段传入），但两者溯源到同一个上游——`forge task index` 命令。该命令扫描 `docs/features/<slug>/` 目录结构，将 `<slug>` 写入 `index.json` 的 `feature` 字段，同时将其传入 embed 模板变量。因此 `index.json` 是 CLI 目录路径的镜像，不存在独立演变的可能：目录重命名后需重新执行 `forge task index`，embed 和 `index.json` 同步更新。

**薄模板与富模板的差异化策略**：6 个模板中，`test-gen-journeys` 和 `test-gen-contracts` 属于富模板（已包含较完整的路径上下文），而 `test-run`、`test-gen-scripts`、`eval-journey`、`eval-contract` 属于薄模板（路径信息极少）。对于富模板，若已有等价的路径引用则不重复添加；对于薄模板，统一补充 `## Feature Paths` 区域。实施时逐一检查每个模板的现有内容，仅在路径上下文不足的模板中添加。

**Discovery 命令定位**：`## Feature Paths` 区域的 `ls` 命令是**信息参考**，供 agent 了解目录布局，而非要求 agent 在 Step 1 执行。agent 应在 skill 被调用后、由 skill 的 Step 1.5 执行实际路径发现。Task file 中的 discovery 命令仅为 agent 提供路径结构预期，避免盲目探索。

### 改动 2：Prompt 模板渲染 `FEATURE_SLUG`

6 个测试流水线 prompt 模板在 `TASK_FILE` 行之后添加：

```
FEATURE_SLUG: {{.FeatureSlug}}
```

`{{.FeatureSlug}}` 由 `run-tasks` dispatcher 从 `index.json` 的 `feature` 字段传入，同样是 `forge task index` 时写入的，确定性强。

### 改动 3：Skill 不变

Skill 保留自己的发现逻辑（支持用户独立调用 `/run-tests`、`/gen-test-scripts` 等，不经过 task pipeline）。

### 三层联动效果

```
Prompt:   FEATURE_SLUG: my-feature              ← dispatcher 从 index.json 传入
Task .md: ## Feature Paths                       ← CLI 从目录路径填充 {{.FeatureSlug}}
            ls docs/features/my-feature/testing/
            ls docs/features/my-feature/testing/<journey>/contracts/
Skill:    Step 1.5 Discover Journeys             ← 保留完整逻辑，支持独立调用
```

三层各有清晰职责，不重复但互相补充：

| 层 | 职责 | 何时可用 |
|---|------|---------|
| Prompt | 给出 slug，agent 无需路径解析 | 任务执行时（dispatcher 传入） |
| Task .md | 给出发现命令，agent 可预探索 | 任务执行时（读 task file） |
| Skill | 完整发现 + 执行逻辑 | 任务执行时 + 独立调用时 |

### Innovation Highlights

无创新，纯信息补全。将已有的 `FeatureSlug` 变量从"声明但不渲染"变为"渲染到输出"，同时将 skill 的发现命令同步到 embed 模板。

## Requirements Analysis

### Key Scenarios

- Subagent 收到 prompt，直接看到 `FEATURE_SLUG: my-feature`，无需从路径解析
- Subagent 读 task file，看到 discovery 命令，可预先定位 journeys 和 contracts
- 用户独立调用 `/run-tests`（不经过 pipeline），skill 的发现逻辑仍然工作

### Constraints & Dependencies

- Embed 模板：`{{.FeatureSlug}}` 已在 identity 组声明，不需新增变量
- Prompt 模板：`{{.FeatureSlug}}` 已在 context 组声明且 `promptTemplateData` 已有此字段，不需修改 Go 结构体
- Skill：无改动

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 无改动 | agent 效率低，三层断裂 | Rejected |
| 只改 embed 模板 | 最小改动 | prompt 仍未渲染 slug，agent 仍需路径解析 | Rejected: 联动不完整 |
| **三层联动（embed + prompt + skill 不变）** | 完整联动，职责清晰，改动量小 | 无 | **Selected** |
| 改 embed + prompt + 简化 skill | 减少重复 | skill 需支持独立调用，不能简化 | Rejected: skill 不能简化 |

## Scope

### In Scope

- 修改 6 个测试流水线 **embed 模板**（`forge-cli/pkg/task/templates/`），统一添加 `## Feature Paths` 区域：
  - `test-gen-journeys.md`
  - `eval-journey.md`
  - `test-gen-contracts.md`
  - `eval-contract.md`
  - `test-gen-scripts.md`
  - `test-run.md`
- 修改 6 个测试流水线 **prompt 模板**（`forge-cli/pkg/prompt/templates/`），添加 `FEATURE_SLUG: {{.FeatureSlug}}`：
  - `test-gen-journeys.md`
  - `eval-journey.md`
  - `test-gen-contracts.md`
  - `eval-contract.md`
  - `test-gen-scripts.md`
  - `test-run.md`

### Out of Scope

- Skill 无改动
- `autogenTemplateData` / `promptTemplateData` Go 结构体不变（`FeatureSlug` 字段已存在）
- 验证 / 文档 / 清理类模板不变
- 不修改任何 Go 代码（仅改模板 .md 文件）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 路径与实际目录不一致 | L | M | 路径基于 forge 硬编码目录约定；embed discovery 命令中的路径与 skill 路径共享同一个权威来源（`docs/features/<slug>/testing/` 目录约定），不存在独立演进，无需额外同步机制 |
| FeatureSlug 渲染为空 | L | **M** | 空 slug 会导致 embed 模板生成无效路径 `ls docs/features//testing/`，agent 执行时将得到空结果或错误；Mitigation: `forge task index` 在扫描时保证 slug 非空（目录名即 slug，空目录名不存在）；prompt 模板中 FeatureSlug 已在 context 声明，dispatcher 始终传入 |

## Success Criteria

- [ ] 6 个 embed 模板均包含 `## Feature Paths` 区域，含 journeys 和 contracts 两个 discovery 命令
- [ ] 6 个 prompt 模板输出 `FEATURE_SLUG: <slug>` 行
- [ ] `go build ./...` 通过
- [ ] `go test ./...` 通过（模板验证测试）
