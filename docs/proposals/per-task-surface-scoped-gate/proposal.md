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

1. **`RunGate()` 增加 prefixed recipe 解析**：当任务有 surface-key（如 `backend`）时，优先尝试 prefixed recipe（`backend-compile`、`backend-lint`），找不到则回退到通用 recipe（`compile`、`lint`）。复用 `quality_gate_lifecycle.go` 中 `resolveRecipe()` 的模式，但改用 surface **key**（非 type）作为前缀。核心伪代码：

   ```go
   // RunGate prefixed resolution — scope guard is critical
   func resolvePrefixedRecipe(scope string, recipe string) string {
       if scope != "" {
           prefixed := scope + "-" + recipe  // e.g. "backend-compile"
           if just.HasRecipe(prefixed) {
               return prefixed
           }
       }
       return recipe  // fallback to generic
   }
   ```

   关键守卫条件：`scope != ""` 确保 feature-level gate（传入空 scope）时不会误触发 prefixed 解析，保持全量验证行为不变。

   **与 ResolveScope() 的集成**：`resolvePrefixedRecipe()` 替代 `ResolveScope()` 在 `RunGate()` 中的作用。当前 `ResolveScope()` 探测参数模式（`compile <scope>`），在新设计中不再使用——prefixed recipe 模式（`<scope>-compile`）直接替代。`RunGate()` 的调用方无需改动，scope 传入方式不变，只是内部解析策略从参数探测切换为前缀匹配。

2. **Surface rules 增加 compile/fmt/lint/unit-test recipe 模板**：api.md、web.md、cli.md、tui.md、mobile.md 各增加 `<key>-compile`/`<key>-fmt`/`<key>-lint`/`<key>-unit-test` 的 stub recipe，使非 mixed 多 surface 项目也能生成 prefixed recipe。

### Innovation Highlights

直接采用已有的 `resolveRecipe()` fallback 模式（surface-specific 优先 → generic 兜底），但将维度从 surface type 切换到 surface key。这是关键区分：lifecycle 测试（dev/probe/test/teardown）按 type 分组是合理的（同 type 共享服务），但 compile/fmt/lint/unit-test 按 key 分组是必要的（每个 surface 有独立代码）。

**Dispatcher isomorphism 差异**：lifecycle 层的 `resolveRecipe()` 对每个 type 执行全量遍历（iterate all by type — 注册了哪些 type 就全部启动），而 `RunGate()` 的 prefixed resolution 是单次选择（pick one by key — 只选择匹配当前 surface-key 的 recipe）。两者虽然共享 fallback 模式，但 dispatcher 语义本质不同：前者是扇出（fan-out），后者是选择（select）。

## Requirements Analysis

### Key Scenarios

1. **多 surface 项目，backend 任务提交**：任务 `surface-key: backend` → per-task gate 执行 `just backend-compile` → `just backend-fmt` → `just backend-lint` → `just backend-unit-test`，跳过 frontend 验证
2. **多 surface 项目，frontend 任务提交**：同理，只验证 frontend surface
3. **单 surface 项目**：任务无 surface-key 或 key 为空 → 回退到通用 recipe（`just compile` → `just lint`），行为与当前一致
4. **多 surface 项目，无 prefixed recipe**：justfile 未生成 surface-specific recipe → 回退到通用 recipe，向后兼容
5. **非 coding 任务提交**：doc/gate/summary 类型 → `IsTestableType()` 返回 false → 不执行 quality gate（已有逻辑，无变化）
6. **Prefixed recipe 存在但执行失败**：`surface-key: backend` → 探测 `just backend-compile` 存在 → 执行失败 → `onFail` 回调收到的 step name 为 `backend-compile`（含 surface 上下文），可直接映射到对应 surface 的错误处理路径，无需额外解析

### Non-Functional Requirements

- **向后兼容**：单 surface 项目、无 surface-specific recipe 的项目行为不变
- **零配置**：自动从任务的 `surface-key` 推导，无需用户手动配置

### Constraints & Dependencies

- Surface rule recipe 模板使用 surface **key** 作为前缀（与 SKILL.md line 259-263 一致），而非 surface type
- `mixed.just` 模板已硬编码 `backend-compile`/`frontend-lint` 等 recipe，不需要修改
- 单语言模板（go.just, node.just 等）不生成 prefixed recipe — 因为单 surface 项目不需要
- **NormalizeSurfaceKey 字符集约束**：`NormalizeSurfaceKey()` 的输出必须限定为 `[a-z][a-z0-9-]*`（小写字母开头，仅含小写字母、数字、连字符）。这保证 `<key>-<recipe>` 形式的 prefixed recipe 名称在 justfile 中合法且无歧义（连字符分隔 key 与 recipe name，key 内部的连字符不与分隔符冲突）

## Alternatives & Industry Benchmarking

### Industry Solutions

多模块项目通常按模块执行验证而非全量验证：

- **Maven** (`-pl module`)：通过 `--projects` 选择器限定反应堆（reactor）范围，只构建/测试指定模块及其依赖。Forge 的 surface-key 等价于 Maven 的 module selector。
- **Gradle** (`includedBuild` + `:subproject:test`)：每个 subproject 有独立的 lint/compile/test task，`gradle :backend:test` 只运行 backend 测试。Gradle 的配置回避避（configuration avoidance）更进一步——只实例化被请求的 subproject。
- **Nx** (`--projects=`)：基于项目图（project graph）的 affected detection，`nx test --projects=backend` 只测试 backend 及其依赖。Nx 还支持 `nx affected --test` 基于 git diff 自动推断受影响的项目。
- **Bazel** (Bazel query + test)：通过 target pattern (`//backend/...`) 精确选择构建/测试目标，Bazel 的 query 引擎支持基于依赖关系的 affected target 发现。
- **Turborepo** (`--filter=`)：`turbo run lint --filter=backend` 只对 backend 执行 lint，与 Forge 的 `<key>-lint` 模式在语义上一致。

共同模式：**用户定义模块边界 → 构建工具按边界选择性执行 → 回退到全量执行（无边界时）**。Forge 的 surface-key 即模块边界，prefixed recipe 即选择性执行机制。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每个多 surface 项目都遇到误阻 | Rejected: 已造成实际阻塞 |
| Git diff affected detection | Nx/Bazel 模式 | 自动推断受影响 surface，零配置 | 需要 git 历史和文件→surface 映射，每次 submit 调用 git diff 增加延迟；对新增文件（无 git 历史）可能误判 | Rejected: 复杂度高，且 Forge 已有 task-level surface-key 不需要推断 |
| 仅改 RunGate() | — | 最小改动 | 非 mixed 项目无 prefixed recipe 可用，fix 不完整 | Rejected: 依赖不存在的 recipe |
| **RunGate() + surface rules 双改** | Turborepo `--filter=` 模式 | 根源修复，与行业 scoped validation 模式一致 | 改动涉及 Go 代码 + 5 个 rule 文件 | **Selected: 与 Turborepo `--filter=` 语义一致，利用已有的 surface-key 边界** |
| Feature-level gate 过滤 | lesson 原文 | 解决 stop hook 问题 | 不解决 per-task gate 的误阻（实际触发点） | Rejected: 攻错了层级 |

## Feasibility Assessment

### Technical Feasibility

完全可行。`RunGate()` 改动局部（增加 prefixed recipe 检测逻辑），`resolveRecipe()` 模式已在 `quality_gate_lifecycle.go` 中验证。Surface rule 改动为文本模板追加，无技术风险。

**HasRecipe() 探测性能分析**：prefixed resolution 每步调用 `just.HasRecipe(prefixed)` 探测 recipe 存在性。标准 gate 共 4 步（compile/fmt/lint/unit-test），每步最多 2 次探测（prefixed + fallback generic），共 8 次 `just --dry-run` fork。在 Windows 上每次 fork 约 50ms，总计 ~400ms 额外开销。但实际路径更短：prefixed recipe 命中时无需探测 generic（直接返回），且无 surface-key 时跳过整个 prefixed 分支（0 probe）。因此最坏 N×2（N=4 steps），平均 < N+2（大多数 recipe 在 prefixed 即命中）。

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
| Prefixed recipe 不存在时回退到 generic，但 generic recipe 仍跑全量 | L | M | 这是预期行为（向后兼容）；新增的 surface rules 会逐步覆盖更多项目类型。中间态（仅改 RunGate()，surface rules 未更新）下 generic 回退仍跑全量，但不会比改动前更差 |
| Surface rule 的 stub recipe 未被 LLM 正确填充 | M | M | 现有 lifecycle recipe 已有相同的 stub 模式，LLM 已能正确处理。stub recipe 模板中须包含约束注释（如 `# This recipe compiles ONLY the <key> surface code`），确保 LLM 理解 recipe 的 surface 隔离语义 |
| `RunGate()` 改动影响 feature-level gate | L | H | `RunGate()` 被 per-task gate 和 feature-level gate 共用；需确保 feature-level gate 传空 scope 时不触发 prefixed 解析 |

## Success Criteria

- [ ] `mixed.just` 项目中，`surface-key: backend` 的任务提交时执行 `just backend-compile`/`just backend-lint`（而非 `just compile`/`just lint`），prefixed recipe 不存在时回退到 generic
- [ ] 单 surface 项目中，任务提交时 `RunGate()` 的 scope 为空 → 跳过 prefixed 分支 → recipe 名称仍为 `compile`/`lint`（非 `backend-compile`），退出码和 stdout 输出与改动前一致
- [ ] `RunGate()` 的 `scope=""` 调用路径（feature-level gate）行为不变
- [ ] 5 个 surface rule 文件均包含 compile/fmt/lint/unit-test 的 Recipe Invocation Contract 条目和 stub recipe
- [ ] 通过 `go test ./internal/cmd/task/...` 和 `go test ./pkg/just/...` 全部测试
- [ ] Prefixed recipe 执行失败时，`onFail` 回调的 step name 为 `<key>-<recipe>`（如 `backend-lint`），错误信息中包含 surface 上下文；回退到 generic recipe 失败时 step name 为原始 recipe 名（如 `lint`），两者可区分

## Next Steps

- Proceed to `/write-prd` to formalize requirements
