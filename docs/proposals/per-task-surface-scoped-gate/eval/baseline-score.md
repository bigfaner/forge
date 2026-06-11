---
iteration: baseline
scorer: cto
date: 2026-06-07
---

# Proposal Evaluation Report — Baseline

**Document**: Per-Task Quality Gate Surface Scoping
**Rubric**: proposal (1000 pts)
**Scorer Persona**: CTO

---

## Dimension Breakdown

### 1. Problem Definition — 100/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 37/40 | The core problem is unambiguous: `validateQualityGate()` runs full compile/fmt/lint/unit-test regardless of surface-key, blocking backend tasks with frontend failures. Quote: "backend 任务被前端 lint 失败阻塞". One reader could not misinterpret this. Deduction: the problem statement opens with the generic claim "对所有 surface 运行全量 compile/fmt/lint/unit-test，不区分任务所属的 surface-key" which is correct, but does not explicitly state the unit of blocking (the per-task `submit` command) until the Evidence section. Minor ambiguity. |
| Evidence provided | 38/40 | Strong code-level evidence: traces `validateQualityGate()` → `just.RunGate()` → `ResolveScope()` → parameter-mode probe → `mixed.just` compile recipe rejects args → scope never takes effect. Real project incident (pm-work-tracker, task 2.3) cited with specific failure cause (`@stylistic/eslint-plugin-ts` npm timeout). Deduction: no quantitative data (e.g., "this blocked N tasks over M days") — evidence is qualitative/anecdotal. |
| Urgency justified | 25/30 | "已造成实际阻塞（任务 stuck、fix-task 循环）。每个多 surface 项目都会遇到。" Concrete and credible. Deduction: "随着多 surface 项目增多" is a forward-looking claim without data on how many multi-surface projects exist or are planned. "What's the cost of delay?" is implied but not quantified (e.g., "X engineer-hours wasted per week"). |

### 2. Solution Clarity — 90/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 32/40 | Two changes clearly named: (1) `RunGate()` prefixed recipe resolution, (2) surface rule stub recipe addition. Quote: "当任务有 surface-key（如 backend）时，优先尝试 prefixed recipe（backend-compile、backend-lint），找不到则回退到通用 recipe". A reader can explain back what will be built. Deduction: no pseudo-code or control flow for the `RunGate()` change. The statement "复用 resolveRecipe() 的模式" is suggestive but not definitive — resolveRecipe iterates all types while RunGate picks one key. The dispatching semantics differ, and without pseudo-code the reader cannot verify the approach handles this correctly. |
| User-facing behavior described | 40/45 | Key Scenarios section provides concrete observable behavior: scenario 1 shows `surface-key: backend` → `just backend-compile` → `just backend-lint`. Scenario 3 shows empty key → generic recipe. Scenario 4 shows no prefixed recipe → fallback. This is good user-facing description. Deduction: failure-mode UX is not described — when a prefixed recipe fails, the error message will show `backend-compile` instead of `compile`, which changes user-visible output. Not discussed. |
| Technical direction clear | 18/35 | The general approach (prefixed recipe lookup + fallback) is clear, but the implementation details are insufficient for a developer to start coding. Missing: (1) where in `RunGate()` the prefixed resolution is inserted (before or after `ResolveScope()`), (2) whether it replaces or augments the existing scope mechanism, (3) the exact guard condition for feature-level gate (quote: "需确保 feature-level gate 传空 scope 时不触发 prefixed 解析" — stated as a risk, not as a design decision). The reference to `resolveRecipe()` pattern is a hint but not a specification. |

### 3. Industry Benchmarking — 62/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 18/40 | Single sentence: "多模块项目通常按模块执行验证（Maven 的 `-pl module`，Nx 的 `--projects=`），而非全量验证。" This is a one-liner mention, not a referenced solution. No open-source project links, no published patterns (e.g., monorepo tooling patterns from Bazel, Turborepo, Gradle, Lerna), no architectural patterns (e.g., affected-project detection, dependency graph traversal). The two tools mentioned (Maven, Nx) are the minimum viable reference. |
| At least 3 meaningful alternatives | 16/30 | The comparison table lists 4 rows: (1) Do nothing, (2) RunGate-only change, (3) RunGate + surface rules (selected), (4) Feature-level gate filtering. Row (1) is "do nothing" — required. Row (2) is a subset of (3), not a genuinely different approach. Row (4) is a different layer attack, which is genuinely different. Only 2 genuinely different alternatives (do nothing, different layer), plus one straw-man (RunGate-only is explicitly "不完整" and presented primarily to be rejected). Missing: industry-validated alternatives such as dependency-graph-based affected detection (Nx/Turborepo style), configuration-based gate profiles, or ignore patterns. |
| Honest trade-off comparison | 14/25 | Pros/cons are honest for the selected approach: "改动涉及 Go 代码 + 5 个 rule 文件". Row (4) correctly identifies it attacks the wrong layer. Deduction: Row (2) "Cons: 非 mixed 项目无 prefixed recipe 可用" is a genuine con, but the Verdict "Rejected: 依赖不存在的 recipe" is circular — it rejects the alternative for not solving the whole problem, which is a completeness argument, not a trade-off argument. |
| Chosen approach justified against benchmarks | 14/25 | The selected approach is justified internally (complete solution vs. partial), but not against industry benchmarks. Quote from comparison: "根源修复，所有多 surface 项目受益" — this justifies completeness, not why this approach over Nx-style dependency graph analysis or Turborepo's affected detection. No rationale for why prefixed recipe lookup is better/worse than industry alternatives. |

### 4. Requirements Completeness — 88/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 32/40 | 5 scenarios identified: (1) backend task with prefixed recipe, (2) frontend task with prefixed recipe, (3) single surface / no key → generic, (4) multi surface but no prefixed recipe → fallback, (5) non-coding task → skip gate. Good happy-path + backward-compat coverage. Deduction: missing failure scenarios — prefixed recipe exists but fails at runtime (what does the user see?), partial gate sequence success (backend-compile passes but backend-lint fails — does the error report include surface context?), and the case where `HasRecipe()` itself fails (subprocess error). |
| Non-functional requirements | 38/40 | Backward compatibility and zero-configuration are clearly stated. Quote: "向后兼容：单 surface 项目、无 surface-specific recipe 的项目行为不变" and "零配置：自动从任务的 surface-key 推导". Deduction: performance NFR not mentioned — the `HasRecipe()` probing adds subprocess calls per gate step, which is a latency concern on Windows. Not acknowledged as an NFR. |
| Constraints & dependencies | 18/30 | Three constraints listed: surface key (not type) as prefix, mixed.just already has prefixed recipes, single-language templates don't generate prefixed recipes. Three dependencies: `surface-key` field exists, `HasRecipe()` available, `mixed.just` has prefixed recipes. Deduction: missing constraints — (1) `NormalizeSurfaceKey` output character set (what characters are valid in recipe names?), (2) the deployment coupling between Go code change and surface rule change (users must re-run `init-justfile`), (3) `RunGate()` is shared between per-task and feature-level gates (stated in risks but not as a constraint). |

### 5. Solution Creativity — 72/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 30/40 | The key insight — lifecycle steps group by type (shared services) while gate steps must group by key (independent codebases) — is a genuine differentiation. Quote: "lifecycle 测试（dev/probe/test/teardown）按 type 分组是合理的（同 type 共享服务），但 compile/fmt/lint/unit-test 按 key 分组是必要的（每个 surface 有独立代码）". This is not a standard industry pattern. Deduction: the implementation (prefixed recipe lookup) is straightforward — the creativity is in the type-vs-key distinction, not in the mechanism. |
| Cross-domain inspiration | 22/35 | The proposal reuses an existing internal pattern (`resolveRecipe()` from lifecycle layer) and adapts it to a new context (gate layer with key-based prefixing instead of type-based). This is internal reuse, not cross-domain inspiration. No borrowing from other domains (e.g., dependency injection patterns, plugin architectures, feature flag systems). The Maven/Nx reference is a brief mention without showing how their patterns informed the design. |
| Simplicity of insight | 20/25 | The type-vs-key distinction is elegant and the kind of insight that seems obvious in retrospect. The solution is minimal (two changes, no new abstractions). Quote from Innovation Highlights captures this well. Deduction: slightly forced — the proposal claims "这是关键区分" but doesn't prove that type-based prefixing for gates would fail (it just states it). A concrete counterexample (what goes wrong with type-based gate prefixing?) would make the insight more compelling. |

### 6. Feasibility — 80/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 32/40 | All prerequisites are in place: `surface-key` field exists, `HasRecipe()` available, `mixed.just` already has prefixed recipes for e2e verification. The `resolveRecipe()` pattern is proven in the lifecycle layer. Quote: "完全可行。RunGate() 改动局部（增加 prefixed recipe 检测逻辑）". Deduction: the `RunGate()` change is described as "局部" but the function is shared between per-task and feature-level gates — the proposal acknowledges this risk but does not provide the guard condition design, which is a feasibility gap (can we change a shared function safely without a clear separation mechanism?). |
| Resource & timeline feasibility | 28/30 | "预计 1-2 小时完成。Go 代码改动集中，surface rule 改动为机械性模板追加。" This is realistic for the described changes. The scope is small and well-defined. Deduction: the estimate does not account for testing — writing tests for the `RunGate()` prefixed resolution (including the feature-level gate guard) could add 1-2 hours. |
| Dependency readiness | 20/30 | All listed dependencies are verifiable and present. Quote: "surface-key 字段已存在于 TaskState、index.json、task YAML frontmatter". Deduction: `HasRecipe()` is listed as ready but its performance characteristics on Windows are not analyzed (each call forks a subprocess). For per-task gate (high frequency), this could be a hidden dependency issue. Also, the proposal does not verify that `NormalizeSurfaceKey` produces recipe-safe names. |

### 7. Scope Definition — 72/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | Three concrete deliverables: (1) `RunGate()` prefixed recipe resolution, (2) 5 surface rule files with gate recipe stubs, (3) backward compatibility verification. Each is a specific artifact. Deduction: "向后兼容验证" is a test activity, not a deliverable — it should specify what verification (unit tests? integration tests? manual test matrix?). |
| Out-of-scope explicitly listed | 25/25 | Six items explicitly out of scope: feature-level gate, test pipeline tasks, mixed.just template, single-language templates, fix-task surface inference, recipes beyond compile/fmt/lint/unit-test. Each is named and reasoned. Excellent. |
| Scope is bounded | 19/25 | The scope is small and executable within the stated timeline. Deduction: the boundary between the two changes (Go code vs surface rules) is clear, but the deployment dependency between them is not bounded — surface rule changes require users to re-run `init-justfile`, which is not in scope. This means the "complete solution" is only complete after user action, but that user action is out of scope. |

### 8. Risk Assessment — 66/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 22/30 | Three risks: (1) generic recipe runs full suite when prefixed not found, (2) LLM stub filling correctness, (3) `RunGate()` shared between per-task and feature-level gates. These are meaningful. Deduction: missing risks — (1) `HasRecipe()` performance on Windows (N×2 subprocess forks per submit), (2) surface key naming collision with existing justfile recipes, (3) deployment atomicity — the fix only activates after users re-run `init-justfile`, which is not obvious. The pre-revision freeform review identified all three of these. |
| Likelihood + impact rated | 22/30 | Ratings are: (L, L), (M, M), (L, H). Not all low — honest distribution. Deduction: Risk 1 "generic recipe runs full" is rated (L, L) but the proposal's own scenario 4 shows this is the expected behavior for non-mixed projects — so Likelihood should be H (many projects won't have prefixed recipes initially). Risk 3 "RunGate() shared" is rated (L, H) — the Likelihood is arguably M because the shared function is the default code path, not an edge case. |
| Mitigations are actionable | 22/30 | Risk 1 mitigation: "预期行为（向后兼容）" — this is a justification, not an actionable mitigation. Risk 2: "现有 lifecycle recipe 已有相同的 stub 模式，LLM 已能正确处理" — cites precedent, which is actionable context. Risk 3: "需确保 feature-level gate 传空 scope 时不触发 prefixed 解析" — this is a requirement for the implementation, not a mitigation plan. A mitigation would specify the mechanism (e.g., "add guard clause `if scope == ""` before prefixed resolution"). |

### 9. Success Criteria — 62/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 24/30 | SC-1: "执行 just backend-compile/just backend-lint" — testable via `just --dry-run` or by mocking `RunGate()`. SC-2: "行为与改动前完全一致" — testable via regression tests. SC-3: "scope="" 调用路径行为不变" — testable with unit test. SC-4: "5 个 surface rule 文件均包含..." — testable via grep. SC-5: "通过 go test ./internal/cmd/task/... 和 go test ./pkg/just/..." — directly testable. Deduction: SC-2 "行为与改动前完全一致" is ambiguous — which specific behaviors? Compile, lint, fmt, unit-test all run? Same output? Same exit codes? "完全一致" is broad. |
| Coverage is complete | 20/25 | SC covers: prefixed execution (SC-1), backward compat (SC-2), feature-level gate safety (SC-3), surface rule completeness (SC-4), test passing (SC-5). In-scope item 1 (RunGate() change) → SC-1, SC-3. In-scope item 2 (surface rule stubs) → SC-4. In-scope item 3 (backward compat) → SC-2. Good coverage. Deduction: SC-4 says "5 个 surface rule 文件均包含... Recipe Invocation Contract 条目和 stub recipe" but does not specify what the stub recipe should compile/lint — just that it exists. A stub recipe that compiles nothing would pass SC-4. |
| SC internal consistency | 18/25 | Clustering by affected area: **RunGate() cluster** (SC-1, SC-3): If SC-1 is satisfied (prefixed recipes run for `scope=backend`), SC-3 must still hold (empty scope → no prefixed resolution). These are compatible: the guard condition `scope != ""` satisfies both bidirectionally. **Surface rules cluster** (SC-4): Independent, satisfiable. **Test cluster** (SC-5): Independent, satisfiable. **Cross-cluster**: SC-2 (single surface unchanged) and SC-1 (multi-surface uses prefixed) — these target different project types, no conflict. Deduction: SC-2 "行为与改动前完全一致" and SC-1 "prefixed recipe 不存在时回退到 generic" — in a single-surface project, SC-1's fallback behavior should match SC-2's "unchanged" behavior. But if a single-surface project has surface rules that generate prefixed recipes (e.g., `api-compile`), SC-1 would use the prefixed recipe instead of generic, which changes behavior. This is ambiguous — does SC-2 mean "for projects with no surface-key on tasks" or "for single-surface projects regardless of recipe availability"? The logical relationship between SC-1's fallback and SC-2's unchanged behavior is unclear when a single-surface project has prefixed recipes. Marked as: **ambiguous — requires author clarification**. |

### 10. Logical Consistency — 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 32/35 | The problem is "backend tasks blocked by frontend lint". The solution makes `RunGate()` resolve prefixed recipes by surface key, so backend tasks only run backend gate steps. This directly addresses the problem. The root cause (scope not taking effect) is fixed by the prefixed recipe mechanism. Deduction: the solution fixes the symptom (wrong recipes executed) but preserves the architecture that caused it — `ResolveScope()` still exists and still doesn't work for parameter-mode recipes. The proposal adds a parallel resolution path rather than fixing `ResolveScope()`. This is pragmatically correct (backward compat) but architecturally adds a second code path where one broken one exists. |
| Scope ↔ Solution ↔ Success Criteria aligned | 22/30 | In-scope items map to solution changes and to success criteria (see D9 coverage analysis). Deduction: In-scope item 1 says "优先探测 <key>-<recipe> 是否存在，存在则用 prefixed recipe 替代 generic recipe" but SC-1 only tests the "存在" case. The "优先" implies ordering (try prefixed first), but no SC tests that generic is tried when prefixed doesn't exist — only scenario 4 describes this, not as a success criterion. Also, Scope lists "向后兼容验证" as in-scope but the verification mechanism is unspecified — how will backward compat be verified beyond SC-2's vague "完全一致"? |
| Requirements ↔ Solution coherent | 18/25 | Requirements (5 scenarios) map to the solution: scenarios 1-4 are covered by the RunGate() change + surface rules, scenario 5 is existing behavior. NFRs (backward compat, zero config) are addressed. Deduction: Constraint "Surface rule recipe 模板使用 surface key 作为前缀" is stated but no requirement specifies what happens when surface key contains characters invalid for just recipe names (e.g., hyphens, spaces). The constraint assumes `NormalizeSurfaceKey` produces recipe-safe names, but this is not stated as a requirement. Also, the solution introduces a new implicit requirement — "users must re-run `init-justfile` to activate the fix" — that is not listed in Requirements Analysis. |

---

## Cross-Dimension Coherence Check

- **D2 (Solution Clarity) vs D10 (Logical Consistency)**: The solution says "复用 resolveRecipe() 的模式" but D10 identifies that the dispatcher semantics differ (single-pick vs iterate-all). This is a gap in D2 (technical direction not clear enough to reveal this distinction) and D10 (the reuse claim is not fully validated).
- **D4 (Requirements) vs D8 (Risks)**: D4 does not include the `HasRecipe()` performance concern, but D8 also misses it. This is a gap in both dimensions.
- **D7 (Scope) vs D6 (Feasibility)**: Scope says "surface rules 增加 gate recipe 模板" but Feasibility does not assess the LLM filling correctness risk with sufficient depth (only rates it M/M).

---

## Blindspot Attacks

1. **[blindspot]** `RunGate()` 共用风险的缓解设计缺失 — quote: "需确保 feature-level gate 传空 scope 时不触发 prefixed 解析" — 提案将此作为风险识别（D8 得分部分），但未提供设计层面的守卫机制。这不是单一维度的问题——它横跨 D2（Solution Clarity 缺少伪代码）、D6（Feasibility 未验证共享函数的安全性）、D8（Risk mitigation 不可操作）。根本问题是：一个被两个调用者共享的函数需要行为变更，但提案没有定义变更的边界条件（guard clause），使得任何单一维度都无法完整捕捉这个缺陷。Reasoning audit flagged this independently of dimension scoring.

2. **[blindspot]** 部署原子性导致的"修复了但没生效"预期管理缺失 — quote: "Surface rules 增加 compile/fmt/lint/unit-test recipe 模板" — `RunGate()` 改动在 CLI 升级后立即生效，但 surface rule 改动需要用户重新运行 `init-justfile`。提案在 Key Risks 中将此评为 (L, L)，但实际 Impact 是 M：用户升级后问题仍存在，会认为修复无效。这不是 D8 单独的问题——它影响 D4（Requirements 未包含 migration/upgrade 需求）、D7（Scope 未包含用户升级路径）、D10（Solution 声称"根源修复"但在中间态不生效，逻辑不一致）。

---

## Deductions Summary

| Rule | Instance | Deduction |
|------|----------|-----------|
| Vague language | SC-2 "行为与改动前完全一致" | Applied in D9 |
| No placeholder text found | — | — |
| No straw-man alternative | Row 2 in comparison table is borderline but has a genuine distinction | — |

---

## Final Score

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 100 | 110 |
| 2. Solution Clarity | 90 | 120 |
| 3. Industry Benchmarking | 62 | 120 |
| 4. Requirements Completeness | 88 | 110 |
| 5. Solution Creativity | 72 | 100 |
| 6. Feasibility | 80 | 100 |
| 7. Scope Definition | 72 | 80 |
| 8. Risk Assessment | 66 | 90 |
| 9. Success Criteria | 62 | 80 |
| 10. Logical Consistency | 72 | 90 |
| **Total** | **764** | **1000** |

---

## Attacks

1. **Industry Benchmarking**: 行业参考过于薄弱 — quote: "Maven 的 -pl module，Nx 的 --projects=" — 仅一行提及两个工具名，无链接、无模式分析、无架构对比。需要补充至少 2-3 个具体的行业方案（如 Bazel affected targets、Turborepo 的依赖图过滤、Gradle 的 incremental builds），并分析其 scoped validation 策略与本提案的异同。

2. **Industry Benchmarking**: 替代方案不够多元化 — quote: "仅改 RunGate()" — 这是所选方案的子集，不是真正独立的替代。需要增加至少一个行业验证过的替代方案（如基于依赖图的 affected detection），并给出诚实的对比。

3. **Solution Clarity**: RunGate() 改动缺少伪代码 — quote: "优先尝试 prefixed recipe（backend-compile、backend-lint），找不到则回退到通用 recipe" — 描述了意图但没有控制流。需要增加伪代码，明确：(1) 触发条件 (`scope != ""`)，(2) 探测顺序，(3) 回退行为，(4) feature-level gate 的安全守卫。

4. **Risk Assessment**: 关键风险遗漏 — 缺少 HasRecipe() Windows 性能风险（per-task gate 高频调用 N×2 子进程探测）和部署原子性风险（CLI 升级后修复不立即生效）的分析。需要补充这两个风险并给出 Likelihood/Impact 评级和可操作的缓解措施。

5. **Success Criteria**: SC-2 "行为与改动前完全一致" 过于模糊 — 需要具体化：哪些行为一致？compile/lint/fmt/unit-test 全部运行？相同输出？相同退出码？同时需要澄清单 surface 项目在有 prefixed recipe 时的行为（与 SC-1 回退机制的交互）。

6. **Success Criteria**: 缺少失败模式的成功标准 — 所有 SC 都是 happy path。需要增加：prefixed recipe 执行失败时错误信息是否包含 surface 上下文、部分 gate sequence 成功时的处理是否正确。

7. **Logical Consistency**: 方案引入隐式部署依赖但未在 Requirements 中声明 — 用户必须重新运行 `init-justfile` 才能激活修复，这是隐式需求。需要在 Requirements Analysis 或 Constraints 中显式声明此升级路径。

8. **[blindspot]**: RunGate() 共用风险的缓解设计缺失 — quote: "需确保 feature-level gate 传空 scope 时不触发 prefixed 解析" — 横跨 Solution Clarity、Feasibility、Risk Assessment 三个维度的结构性缺陷。需要在 Proposed Solution 中增加 guard condition 设计（如伪代码中的 `if scope != ""` 守卫）。

9. **[blindspot]**: 部署原子性导致的"修复了但没生效"预期管理缺失 — quote: "Surface rules 增加 compile/fmt/lint/unit-test recipe 模板" — RunGate() 改动立即生效但 surface rule 需要用户主动操作。需要在提案中增加 Migration Notes 或 Upgrade Guide，并在 Key Risks 中将此风险的 Impact 从 L 修正为 M。
