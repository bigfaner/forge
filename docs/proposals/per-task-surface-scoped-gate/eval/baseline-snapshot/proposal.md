---
created: "2026-06-07"
author: "faner"
status: Draft
intent: "enhancement"
---

# Proposal: Per-Task Quality Gate Surface Scoping

## Problem

Per-task quality gate（`forge task submit` 中的 `validateQualityGate()`）对所有 surface 运行全量 compile/fmt/lint/unit-test，不区分任务所属的 surface-key，导致 backend 任务被前端 lint 失败阻塞。

### Evidence

pm-work-tracker 项目（`backend=api` + `frontend=web`）中，task 2.3（纯 backend 任务，`surface-key: backend`）执行后被 quality gate 阻塞。阻塞原因不是 backend 问题，而是 frontend lint 无法安装 `@stylistic/eslint-plugin-ts`（npm 网络超时）。后续生成的 fix-task 也是 backend 任务，同样被前端 lint gate 卡住。

代码层面验证：`validateQualityGate()` 调用 `just.RunGate()`，传入 `scope`（即 surfaceKey），但 `RunGate()` 通过 `ResolveScope()` 探测 `just --dry-run compile backend`（参数模式）。`mixed.just` 的 `compile` recipe 不接受参数 → scope 永远不生效 → 始终执行全量 compile/fmt/lint/unit-test。虽然 `mixed.just` 已有 `backend-compile`/`backend-lint` 等 prefixed recipe，但 `RunGate()` 从未使用它们。

### Urgency

已造成实际阻塞（任务 stuck、fix-task 循环）。每个多 surface 项目都会遇到。随着多 surface 项目增多，问题持续复现。

## Proposed Solution

两处改动，使 per-task gate 按 surface-key 验证：

1. **`RunGate()` 增加 prefixed recipe 解析**：当任务有 surface-key（如 `backend`）时，优先尝试 prefixed recipe（`backend-compile`、`backend-lint`），找不到则回退到通用 recipe（`compile`、`lint`）。复用 `quality_gate_lifecycle.go` 中 `resolveRecipe()` 的模式，但改用 surface **key**（非 type）作为前缀。

2. **Surface rules 增加 compile/fmt/lint/unit-test recipe 模板**：api.md、web.md、cli.md、tui.md、mobile.md 各增加 `<key>-compile`/`<key>-fmt`/`<key>-lint`/`<key>-unit-test` 的 stub recipe，使非 mixed 多 surface 项目也能生成 prefixed recipe。

### Innovation Highlights

直接采用已有的 `resolveRecipe()` fallback 模式（surface-specific 优先 → generic 兜底），但将维度从 surface type 切换到 surface key。这是关键区分：lifecycle 测试（dev/probe/test/teardown）按 type 分组是合理的（同 type 共享服务），但 compile/fmt/lint/unit-test 按 key 分组是必要的（每个 surface 有独立代码）。

## Requirements Analysis

### Key Scenarios

1. **多 surface 项目，backend 任务提交**：任务 `surface-key: backend` → per-task gate 执行 `just backend-compile` → `just backend-fmt` → `just backend-lint` → `just backend-unit-test`，跳过 frontend 验证
2. **多 surface 项目，frontend 任务提交**：同理，只验证 frontend surface
3. **单 surface 项目**：任务无 surface-key 或 key 为空 → 回退到通用 recipe（`just compile` → `just lint`），行为与当前一致
4. **多 surface 项目，无 prefixed recipe**：justfile 未生成 surface-specific recipe → 回退到通用 recipe，向后兼容
5. **非 coding 任务提交**：doc/gate/summary 类型 → `IsTestableType()` 返回 false → 不执行 quality gate（已有逻辑，无变化）

### Non-Functional Requirements

- **向后兼容**：单 surface 项目、无 surface-specific recipe 的项目行为不变
- **零配置**：自动从任务的 `surface-key` 推导，无需用户手动配置

### Constraints & Dependencies

- Surface rule recipe 模板使用 surface **key** 作为前缀（与 SKILL.md line 259-263 一致），而非 surface type
- `mixed.just` 模板已硬编码 `backend-compile`/`frontend-lint` 等 recipe，不需要修改
- 单语言模板（go.just, node.just 等）不生成 prefixed recipe — 因为单 surface 项目不需要

## Alternatives & Industry Benchmarking

### Industry Solutions

多模块项目通常按模块执行验证（Maven 的 `-pl module`，Nx 的 `--projects=`），而非全量验证。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每个多 surface 项目都遇到误阻 | Rejected: 已造成实际阻塞 |
| 仅改 RunGate() | — | 最小改动 | 非 mixed 项目无 prefixed recipe 可用，fix 不完整 | Rejected: 依赖不存在的 recipe |
| **RunGate() + surface rules 双改** | — | 根源修复，所有多 surface 项目受益 | 改动涉及 Go 代码 + 5 个 rule 文件 | **Selected: 完整解决方案** |
| Feature-level gate 过滤 | lesson 原文 | 解决 stop hook 问题 | 不解决 per-task gate 的误阻（实际触发点） | Rejected: 攻错了层级 |

## Feasibility Assessment

### Technical Feasibility

完全可行。`RunGate()` 改动局部（增加 prefixed recipe 检测逻辑），`resolveRecipe()` 模式已在 `quality_gate_lifecycle.go` 中验证。Surface rule 改动为文本模板追加，无技术风险。

### Resource & Timeline

预计 1-2 小时完成。Go 代码改动集中，surface rule 改动为机械性模板追加。

### Dependency Readiness

所有前置条件已就绪：
- 任务 `surface-key` 字段已存在于 `TaskState`、index.json、task YAML frontmatter
- `just.HasRecipe()` 可探测 prefixed recipe 是否存在
- `mixed.just` 已有 prefixed recipe 可验证端到端行为

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| `RunGate()` 的 scope 机制已按 surface 过滤 | Codebase Check | Overturned: `ResolveScope()` 探测参数模式，`mixed.just` 的 `compile` 不接受参数，scope 永远不生效 |
| Surface rules 已覆盖所有需要的 recipe | Codebase Check | Overturned: Surface rules 只生成 lifecycle recipe（dev/probe/test/teardown），不生成 compile/fmt/lint/unit-test |
| Feature-level gate 是误阻的触发点 | 5 Whys | Refined: 实际触发点是 per-task gate（`forge task submit`）。Feature-level gate 全量验证是其设计职责 |
| Lifecycle recipe 的 type 前缀模式适用于 gate recipe | Assumption Flip | Overturned: Lifecycle 按 type 分组合理（同 type 共享服务），gate recipe 必须按 key 分组（每个 surface 有独立代码） |

## Scope

### In Scope

1. **`RunGate()` 增加 prefixed recipe 解析**：当 `scope`（surfaceKey）非空时，优先探测 `<key>-<recipe>` 是否存在，存在则用 prefixed recipe 替代 generic recipe
2. **Surface rules 增加 gate recipe 模板**：api.md、web.md、cli.md、tui.md、mobile.md 各增加 `<key>-compile`、`<key>-fmt`、`<key>-lint`、`<key>-unit-test` 的 stub recipe 定义和 Recipe Invocation Contract 条目
3. **向后兼容验证**：单 surface 项目、无 prefixed recipe 的项目回退到 generic recipe

### Out of Scope

- Feature-level gate（`forge quality-gate`）— 全量回归是其设计职责，不改动
- Test pipeline 任务（`test.run`）— 已按 surface 生成
- `mixed.just` 模板 — 已有 prefixed recipe，不需要修改
- 单语言模板（go.just、node.just 等）— 单 surface 项目无需 prefixed recipe
- Fix-task surface 推断（`inferSurface()`）— 已有基于文件路径的机制
- `compile`/`fmt`/`lint`/`unit-test` 之外的 recipe

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Prefixed recipe 不存在时回退到 generic，但 generic recipe 仍跑全量 | L | L | 这是预期行为（向后兼容）；新增的 surface rules 会逐步覆盖更多项目类型 |
| Surface rule 的 stub recipe 未被 LLM 正确填充 | M | M | 现有 lifecycle recipe 已有相同的 stub 模式，LLM 已能正确处理 |
| `RunGate()` 改动影响 feature-level gate | L | H | `RunGate()` 被 per-task gate 和 feature-level gate 共用；需确保 feature-level gate 传空 scope 时不触发 prefixed 解析 |

## Success Criteria

- [ ] `mixed.just` 项目中，`surface-key: backend` 的任务提交时执行 `just backend-compile`/`just backend-lint`（而非 `just compile`/`just lint`），prefixed recipe 不存在时回退到 generic
- [ ] 单 surface 项目中，任务提交行为与改动前完全一致（执行 `just compile`/`just lint`）
- [ ] `RunGate()` 的 `scope=""` 调用路径（feature-level gate）行为不变
- [ ] 5 个 surface rule 文件均包含 compile/fmt/lint/unit-test 的 Recipe Invocation Contract 条目和 stub recipe
- [ ] 通过 `go test ./internal/cmd/task/...` 和 `go test ./pkg/just/...` 全部测试

## Next Steps

- Proceed to `/write-prd` to formalize requirements
