# Freeform Expert Review

**Reviewer Profile**: Surface-Aware Dispatcher & Test Orchestration Architect
**Document**: Per-Task Quality Gate Surface Scoping
**Date**: 2026-06-07

---

## Background Assessment

This proposal addresses a concrete, production-verified problem: in multi-surface projects (e.g. `backend=api` + `frontend=web`), the per-task quality gate invoked by `forge task submit` runs the full compile/fmt/lint/unit-test suite against all surfaces, regardless of which surface the submitted task belongs to. The root cause is clearly diagnosed — `RunGate()` in `just.go` receives a `scope` parameter (the task's `surfaceKey`), but `ResolveScope()` probes for parameterized recipes (`compile <scope>`), not prefixed recipes (`backend-compile`). Since `mixed.just`'s `compile` recipe takes no arguments, scope resolution always returns empty, and the gate runs unscoped.

The proposed solution is two-fold: (1) extend `RunGate()` to resolve prefixed recipes (`<key>-<recipe>`) before falling back to generic recipes, and (2) add prefixed gate recipe stubs (`<key>-compile`, `<key>-fmt`, `<key>-lint`, `<key>-unit-test`) to all five surface rule files (api.md, web.md, cli.md, tui.md, mobile.md) under `init-justfile/rules/surfaces/`. The proposal explicitly reuses the `resolveRecipe()` pattern already proven in `quality_gate_lifecycle.go` for lifecycle orchestration, but shifts the prefix dimension from surface type to surface key — a distinction the proposal argues is necessary because lifecycle steps (dev/probe/test/teardown) group by type (shared services), while gate steps (compile/fmt/lint/unit-test) must group by key (independent codebases).

The proposal rests on several assumptions: that `HasRecipe()` reliably detects prefixed recipes cross-platform, that the `RunGate()` change can be safely scoped to per-task invocations only, that surface rule stubs will be correctly filled by the LLM during justfile generation, and that the two changes (Go code + surface rules) can be deployed atomically without intermediate breakage.

---

## Key Risks

The core technical idea — prefixed recipe resolution in `RunGate()` — is sound and already proven in the lifecycle layer via `resolveRecipe()`. However, the proposal contains several risks at the boundary between the two changes and in the deployment model.

风险：`RunGate()` 被 per-task gate 和 feature-level gate 共用，提案虽然声明了 awareness，但未给出具体的 scope 守卫机制设计。

> "需确保 feature-level gate 传空 scope 时不触发 prefixed 解析" — Key Risks 表格

查看 `quality_gate.go` 第 161 行，feature-level gate 调用 `just.RunGate(result.ProjectRoot, "", gateSteps, ...)`, scope 为空字符串。当前的 `ResolveScope()` 对空 scope 直接返回空，所以现有逻辑不会触发 prefixed 解析。但提案要改变的是 `RunGate()` 本身——如果新的 prefixed resolution 逻辑被放在 `ResolveScope()` 之后、recipe 执行之前，那么当 `resolvedScope == ""` 但步骤名如 `compile` 存在、同时 `backend-compile` 也存在时，是否仍会触发 prefixed 路径？提案没有给出伪代码或控制流设计，这是关键的实现空白。如果 prefixed resolution 仅在 `scope != ""` 时触发（即与 `ResolveScope()` 解耦），那 feature-level gate 自然安全。但如果 prefixed resolution 是无条件的（对每个 step name 尝试 `<key>-<step>` 前缀），那 feature-level gate 在多 surface 项目中会意外只验证第一个 surface，造成回归盲区。

问题：提案声称要"复用 `quality_gate_lifecycle.go` 中 `resolveRecipe()` 的模式"，但两个 dispatcher 的上下文语义完全不同。

> "复用 `quality_gate_lifecycle.go` 中 `resolveRecipe()` 的模式，但改用 surface key（非 type）作为前缀" — Proposed Solution

`resolveRecipe()` 在 lifecycle 层的输入是 `surfaceType`（如 "web"、"api"），且调用方已经按 type 遍历（`runTestRegressionSurface` 中 `for _, surfaceType := range surfaceTypes`）。而 `RunGate()` 的输入是单个 `scope`（任务级别的 surfaceKey，如 "backend"）。这意味着 `RunGate()` 不会遍历所有 surface——它只验证任务所属的那一个 surface。这是正确的行为，但提案没有明确指出：**prefixed recipe resolution 在 `RunGate()` 中是单次选择（pick one），而 lifecycle 层是全量遍历（iterate all）**。这个 dispatcher isomorphism 差异需要在设计文档中显式声明，否则实现者可能错误地将遍历逻辑引入 `RunGate()`。

风险：`HasRecipe()` 基于 `just --dry-run` 子进程探测，在 Windows 上的性能和可靠性未充分分析。

> "`just.HasRecipe()` 可探测 prefixed recipe 是否存在" — Dependency Readiness

当前 `HasRecipe()` 实现（`just.go` 第 63-67 行）每次调用都 fork 一个 `just --dry-run <recipe>` 子进程。在 `RunGate()` 的 per-task 路径中，对于 UnitGateSequence（4 steps），如果每个 step 先探测 prefixed recipe、再探测 generic recipe，那一次 submit 就会触发 8 次 `just --dry-run` 子进程（4 prefixed + 4 generic，最坏情况）。在 Windows 上，进程创建开销比 Linux 大一个数量级。提案没有分析这个 N×2 探测的性能影响，也没有考虑缓存策略。`resolveRecipe()` 在 lifecycle 层也有同样的 N×2 探测问题，但 lifecycle 是 per-feature（频率低），而 per-task gate 是 per-submit（频率高）。

风险：Surface rule stub recipe 的 LLM 填充正确性依赖隐式知识传递链。

> "Surface rule 的 stub recipe 未被 LLM 正确填充" — Key Risks 表格，Likelihood=M, Impact=M

提案的缓解措施是"现有 lifecycle recipe 已有相同的 stub 模式，LLM 已能正确处理"。但 lifecycle stub 和 gate stub 的语义不同：lifecycle recipe（`api-dev`、`api-probe`）是运行时行为（启动服务、健康检查），LLM 填充时可以从 Convention knowledge 推导。而 gate recipe（`api-compile`、`api-lint`）是构建时行为，且需要与已有的 generic `compile`/`lint` recipe 保持一致。如果 `init-justfile` skill 为多 surface 项目同时生成了 generic recipe（如 `compile` 扫全量）和 prefixed recipe（如 `api-compile` 只扫描 api 目录），两者必须语义一致。提案没有讨论这个一致性保证机制——是否 generic recipe 应该变成所有 prefixed recipe 的聚合？还是 generic recipe 只在 prefixed 不存在时使用？

问题：提案没有定义 prefixed recipe 执行失败时的错误信息中是否包含 surface 上下文。

> "执行 `just backend-compile` → `just backend-fmt` → `just backend-lint` → `just backend-unit-test`" — Key Scenarios 1

当前 `RunGate()` 的 `onFail` 回调接收 `(step, output)`，其中 `step` 是 recipe 名。如果使用 prefixed recipe，`step` 将是 `"backend-compile"` 而非 `"compile"`。`submit.go` 中的 `validateQualityGate()` 将 `"Quality gate failed at step: just %s"` 格式化给用户。如果 prefixed recipe 命名（如 `backend-compile`）本身就包含了 surface 信息，那错误信息对用户是清晰的。但如果回退到 generic recipe，用户看到的是 `"Quality gate failed at step: just compile"`——无法区分这是 scoped 还是 unscoped 失败。提案没有讨论这个 UX 一致性问题。

风险：两个改动的部署原子性依赖于用户重新运行 `init-justfile`。

> "Surface rules 增加 compile/fmt/lint/unit-test recipe 模板" — Proposed Solution

`RunGate()` 的改动在 CLI 升级后立即生效，但 surface rule 的改动需要用户重新运行 `init-justfile` 才能生成 prefixed recipe。在 CLI 已升级但 justfile 未重新生成的中间态，`RunGate()` 会探测 prefixed recipe，找不到，回退到 generic recipe。这是安全的（向后兼容），但意味着修复在中间态不生效——用户的 backend 任务仍会被 frontend lint 阻塞。提案将此列为 "Likelihood=L, Impact=L"，但我认为 Impact 应该是 M：如果用户不知道需要重新生成 justfile，问题看起来就像"升级了但没修复"。

风险：提案未考虑 surface key 与 prefixed recipe 名之间的命名冲突。

在 `mixed.just` 中，prefixed recipe 使用 `frontend-compile`、`backend-compile`。但 surface key 是用户在 `.forge/config.yaml` 中自定义的字符串（经过 `NormalizeSurfaceKey` 规范化）。如果用户的 surface key 是 `common` 或 `shared`，而 justfile 中恰好有一个 `common-compile` recipe 做其他用途（比如编译共享库），`HasRecipe()` 会返回 true，导致错误匹配。`NormalizeSurfaceKey` 做了哪些字符限制？提案没有讨论这个边界条件。在 lifecycle 层，surface type 是枚举值（web/api/cli/tui/mobile），不存在冲突。但 gate 层用 surface key（用户自定义字符串），命名空间更大，冲突概率虽然低但非零。

---

## Improvement Suggestions

建议：在提案中增加 `RunGate()` 改动的伪代码，明确 prefixed resolution 的触发条件和控制流。

Addresses: "RunGate() 共用风险" 和 "dispatcher isomorphism 差异"

具体来说，伪代码应展示：(1) prefixed resolution 仅在 `scope != ""` 时触发（不是无条件的）；(2) 对每个 gate step，先构造 `<scope>-<step.Name>` 并调用 `HasRecipe()`，找到则替换 step name；(3) 找不到则保持原 step name 不变。这与 `resolveRecipe()` 的模式一致，但关键区别是触发条件——`resolveRecipe()` 总是尝试 prefixed（因为调用方已经按 type 遍历），而 `RunGate()` 应仅在 scope 非空时尝试。伪代码加上这个守卫条件后，feature-level gate 传空 scope 时自然跳过 prefixed 路径，消除共用风险。

建议：增加 `HasRecipe()` 探测结果的进程内缓存（TTL=0，per-call），或改为先解析 justfile 文本再探测。

Addresses: "HasRecipe() Windows 性能风险"

`RunGate()` 的 per-task 路径中，一次 submit 最多触发 8 次 `just --dry-run`（4 steps × 2 probes）。在 Windows 上每次 fork 约 50-100ms，总计 400-800ms 延迟。建议在 `RunGate()` 入口处一次性收集所有需要的 recipe 名（包括 prefixed 变体），批量探测。或者更简单地，将 `HasRecipe()` 的调用次数从 N×2 减少到 N——如果 `HasRecipe(scope+"-"+step.Name)` 返回 true，就不需要再调用 `HasRecipe(step.Name)`。提案的 `resolveRecipe()` 模式已经是这样的（先探 specific，再探 generic），所以实际上最坏情况是 2N 次（所有 prefixed 都不存在），平均情况 < N+2（大部分 prefixed 存在）。这个优化足够了，但提案应该显式分析并记录。

建议：在 surface rule stub recipe 中增加注释，说明 prefixed gate recipe 与 generic recipe 的关系。

Addresses: "LLM 填充正确性" 和 "generic/prefixed 一致性"

当前的 lifecycle stub recipe（如 `api-dev`）是独立的——它们没有对应的 generic 版本。但 gate stub recipe（如 `api-compile`）有对应的 generic 版本（`compile`）。建议在 surface rule 模板中增加约束注释：`# This recipe should compile ONLY the <key> surface code. The generic "compile" recipe compiles ALL surfaces.` 这样 LLM 在填充 stub 时能理解两者的区别，避免生成语义不一致的 recipe。

建议：明确文档化中间态的预期行为和升级路径。

Addresses: "部署原子性风险"

提案应增加一段 "Migration Guide"：CLI 升级后，已有 justfile 不受影响（回退到 generic recipe，行为不变）。要激活 per-task scoped gate，用户需要运行 `forge init-justfile` 重新生成 justfile。这个中间态是安全的但功能不生效，建议在 CLI 的 release notes 中明确说明。也可以考虑在 `RunGate()` 检测到 scope 非空但 prefixed recipe 不存在时，打印一次 info-level 提示："scope=backend but no backend-compile recipe found; falling back to generic compile. Re-run `forge init-justfile` to generate surface-scoped recipes."

建议：增加 `NormalizeSurfaceKey` 输出字符集的文档化验证，确保 key 不与 justfile recipe 命名空间冲突。

Addresses: "命名冲突风险"

查看 `execution_order.go`，`NormalizeSurfaceKey` 对 key 做了哪些规范化？如果它已经限制了字符集为 `[a-z][a-z0-9-]*`，那与 just recipe 命名规则一致，冲突概率极低（因为 prefixed recipe 的前缀就是 surface key）。但提案应显式引用 `NormalizeSurfaceKey` 的约束，作为 prefixed recipe 命名安全性的论据。如果 `NormalizeSurfaceKey` 不限制字符集，那这是一个必须修复的前置条件。

建议：在提案的 Key Scenarios 中增加 "prefixed recipe 存在但执行失败" 的失败模式分析。

Addresses: "失败模式分析空白"

当前 5 个 Key Scenarios 只覆盖了 happy path（prefixed 存在并成功）和 backward compat（prefixed 不存在回退）。缺少：(1) prefixed recipe 存在但执行失败——此时 `onFail` 收到的 step 是 `backend-compile`，错误信息是否足够清晰？(2) 同一个 gate sequence 中，`backend-compile` 成功但 `backend-lint` 失败——部分成功是否正确处理（当前逻辑是 blocking=true 直接返回 false，但 fmt 是 blocking=false 会继续）？这些失败模式在 `RunGate()` 现有逻辑中已经正确处理，但提案应显式验证，因为 prefixed recipe 的引入改变了 `step.Name` 的值，可能影响 downstream 的错误分类逻辑。
