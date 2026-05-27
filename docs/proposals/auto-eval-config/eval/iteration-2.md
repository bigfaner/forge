# Eval Report: Iteration 2

**Reviewer**: CTO Adversary (Rubric Scoring)
**Document**: `docs/proposals/auto-eval-config/proposal.md` (revised after iteration-1)
**Date**: 2026-05-26

---

## Iteration-1 Issue Resolution Tracker

| # | Iter-1 Attack | Resolution Status | Assessment |
|---|---------------|-------------------|------------|
| 1 | Urgency is asserted, not demonstrated | **Partially resolved** — line 24 now quantifies "~60-80s 交互时间/每条流水线运行" with per-interaction breakdown ("约 15-20s 延迟"). But measurement methodology still unstated (timed? estimated?). |
| 2 | Set path YAML Node confusion | **Resolved** — lines 41-46 completely rewritten to "Go struct marshal" with explicit rationale for choosing struct marshal over YAML Node. No residual ambiguity. |
| 3 | Error messages undefined | **Resolved** — line 46 defines error behavior for 4 cases: nonexistent field, non-leaf set, excess depth, type mismatch. Key Scenarios (lines 91-92) and SC (lines 213-216) cover error paths. |
| 4 | Viper reference superficial | **Improved** — line 120 now discusses Viper's flat key map approach, nil pointer handling gap, custom type formatting gap, and raw tracking gap. 3x more depth than baseline. |
| 5 | "仅扩展路由支持三层" straw-man | **Improved** — line 116 now quantifies: "~30 行/字段 vs 200 行一次性投入" with breakeven analysis (~7 fields). Still somewhat hand-wavy on the 200-line estimate but dramatically better. |
| 6 | Breakeven analysis missing | **Resolved** — line 116 provides explicit breakeven: "当前已规划 eval（4 字段）+ 未来知识管理和验证自动化预计 3-4 字段" = ~7-8 fields approaching breakeven. |
| 7 | Performance NFR unsubstantiated | **Resolved** — line 98 quantifies: "GetConfigValue 端到端延迟 <1ms（当前硬编码分发器 <0.1ms）" with explanation "反射调用次数等于 key 深度，每次 <0.1ms". Concrete bound provided. |
| 8 | Mode detection outside feature directory | **Resolved** — line 90 specifies: "在 feature 目录外调用 `forge config get mode` 返回 `"none"`". Line 105 adds CLI binary fallback behavior. |
| 9 | Creativity floor acknowledged | **Unchanged** — "复用已有 ModeToggle 模式，无新概念" at line 75. No improvement attempted. Author's choice. |
| 10 | Out-of-scope self-contradiction | **Resolved** — line 179 now says "forge guide 文档更新（后续迭代）" without "随代码一起完成". Clean out-of-scope. |
| 11 | Mode detection API unscoped mini-feature | **Improved** — line 171 fully specifies mode detection API with 3 return values and their conditions. SC lines 224-226 cover all 3 cases. Scope item is now detailed. |
| 12 | parseAutoRaw risk underrated | **Resolved** — Risk table line 202 now rates parseAutoRaw at M/M with detailed mitigation including specific test case name. |
| 13 | Skill-to-CLI subprocess dependency risk | **Resolved** — Risk table line 204 covers this explicitly: "CLI binary 未构建或版本过旧" with mitigation "CLI 退出码非零时回退到 AskUserQuestion；CI 中在 skill 测试前强制 `go install`". |
| 14 | No SC for parseAutoRaw | **Resolved** — SC line 219: "parseAutoRaw 对 `auto.eval.*` 字段生成正确的 flat-path raw map（如 `map["eval.proposal"]["quick"]=true`）" with specific test case name. |
| 15 | No SC for mode detection | **Resolved** — SC lines 224-226 cover all 3 mode detection scenarios. |
| 16 | SC formatting inconsistency | **Partially resolved** — line 212 defines intermediate node format explicitly: "每行一个字段，格式为 `<fieldName>: quick:<bool> full:<bool>`". But the separator style still differs from leaf node output (space after colon inconsistency). |
| 17 | Scope-Solution-SC misalignment for mode detection | **Resolved** — Scope (line 171), Solution (line 76, 89-90), and SC (lines 224-226) all cover mode detection consistently. |
| 18 | SC not mapped to PRs | **Resolved** — SC section now split into "PR-1: Generic Config Key Resolution" (lines 210-220) and "PR-2: Auto-Eval Configuration" (lines 222-235). Clean mapping. |

**Summary**: 18 iteration-1 attacks tracked. 12 resolved, 4 improved, 1 partially resolved, 1 unchanged (creativity). Revision quality is high.

---

## Dimension Scoring

### 1. Problem Definition (95/110)

**Problem stated clearly (37/40)**:
The dual-problem framing is well-executed. The eval interaction friction and config routing extensibility bottleneck are both real and verified against code. The relationship between the two problems (eval config is the *trigger* for needing deeper nesting, which motivates the generic routing rewrite) is clearer in this revision.

- **Deduction -3**: The two problems remain conflated in urgency. The routing bottleneck is a legitimate architectural concern regardless of eval configuration. The eval interaction friction could be solved without the routing rewrite (just add 4 more fields to `autoModeField`). The proposal acknowledges this in the alternatives table (line 116: "~30 行/字段") but the Problem section doesn't clearly establish *why* solving both simultaneously is necessary rather than sequential.

**Evidence provided (35/40)**:
Evidence is concrete and code-verifiable. `strings.SplitN(rest, ".", 2)` at line 598 of `config.go` — confirmed at line 598. `autoModeField` switch/case at lines 515-531 — confirmed. `parseAutoRaw` hardcoded `modeFields` at line 383 — confirmed. The `AskUserQuestion` pattern in brainstorm (line 117-122), write-prd (line 215-220), tech-design (line 170-176) — all confirmed. ui-design unconditional auto-eval at Step 7 — confirmed.

- **Deduction -3**: The urgency section (line 24) now quantifies "~60-80s 交互时间/每条流水线运行" with "约 15-20s 延迟" per AskUserQuestion. But this is an estimate, not a measurement. No methodology is stated. Is this from user observation? Log analysis? Developer self-timing? The number is plausible but unverifiable.
- **Deduction -2**: No user feedback or frustration evidence. The problem is developer-inferred, not user-reported.

**Urgency justified (23/30)**:
"eval 的手动确认成为剩余的主要交互摩擦点" — the "剩余" framing is effective, showing this is the last remaining friction point after other automations.

- **Deduction -4**: The urgency for the generic routing rewrite is weaker than for eval config. "泛化 config key resolution 是一次投资长期收益的基础设施改进" — this is a *general* infrastructure argument that could be made at any time. Why now? The trigger is eval config needing 3-level nesting, but the alternatives table shows this can be done with ~30 lines of hardcoded extension. The urgency for the architectural improvement is assumed, not established.
- **Deduction -3**: No cost-of-delay calculation. If the proposal cited "we expect 3 more config fields in the next 2 months" with specific planned features, urgency would be grounded. The line 116 mention of "未来知识管理和验证自动化预计 3-4 字段" is vague future planning.

### 2. Solution Clarity (103/120)

**Approach is concrete (36/40)**:
The two-part structure is clear and well-separated. Part 1 (generic routing via reflection) and Part 2 (eval config) are independently understandable. Function signatures (`getStructValueByPath`, `setStructValueByPath`) are provided. The "消除的硬编码" list (lines 48-52) gives clear deletion scope.

- **Deduction -2**: The transition plan from current code to new code is incomplete. The proposal lists functions to delete but doesn't trace the call graph: `GetConfigValue` currently dispatches to `getAutoKeyValue` / `getWorktreeKeyValue` / `getCoverageKeyValue` — the new version dispatches to a single `getByPath`. But `getAutoKeyValue` also handles `auto.gitPush` (a bool, not ModeToggle). How does the generic reflect-based router handle `auto.gitPush`? It's a `bool` field on `AutoConfig`, so reflection should work — but this transition detail isn't explicitly addressed.
- **Deduction -2**: The `Config` struct has `Surfaces SurfacesMap` (a custom type with `UnmarshalYAML`/`MarshalYAML`). If reflection encounters this custom type, does it call the custom marshaler? The proposal doesn't address how custom YAML types interact with reflection-based routing.

**User-facing behavior described (36/45)**:
CLI behavior is well-described. 10 key scenarios cover the main paths. Error behavior is now specified (lines 46, 91-92).

- **Deduction -4**: The intermediate node get behavior (line 212) says "每行一个字段，格式为 `<fieldName>: quick:<bool> full:<bool>`，字段按 Go struct 定义顺序排列". This format specification is for `auto.eval`, which has 4 ModeToggle children. But what about `auto` itself? `forge config get auto` would need to enumerate *all* auto sub-fields including `Test`, `ConsolidateSpecs`, `GitPush` (bool), and the new `Eval` (nested struct). The formatting rules for mixed-type intermediate nodes (some ModeToggle, some bool, some nested struct) are undefined.
- **Deduction -3**: Skill-side behavior is specified as "config check" (line 170) but the actual skill markdown modification pattern isn't shown. What does the config check look like in a skill file? A bash `if` block? A rules file directive? The proposal says "forge config get + mode 检测" but doesn't show the skill-level pseudocode. This makes the PR-2 implementation ambiguous.
- **Deduction -2**: The `auto.eval` intermediate node set rejection (line 213) says error message is "cannot set non-leaf key, use auto.eval.<field>.<subfield>". But what about setting `auto.eval.proposal` (a ModeToggle, which is a non-leaf in the sense that it has sub-fields)? Can you set it to a boolean (setting both quick and full)? The current `setAutoConfigValue` supports `auto.{field}` for ModeToggle (setting both). The proposal doesn't clarify this for nested ModeToggles.

**Technical direction clear (31/35)**:
Get path (reflection) and Set path (struct marshal) are clear and internally consistent. YAML tag matching priority, nil pointer handling, leaf type formatting — all specified.

- **Deduction -2**: The Get path says "map → 按 key 查找后递归格式化 value" (line 37). The `CoverageConfig.ByType` is `map[string]CoverageStrategy` where `CoverageStrategy` is a struct with `Type` and `Percentage` fields. "递归格式化 value" for a struct means... what? Each field on its own line? A single-line representation? The proposal addresses this in the risk table (line 200) with "注册类型感知格式化器" but this is a risk mitigation, not a solution design. The formatting behavior for map-of-struct values is under-specified.
- **Deduction -2**: The `parseAutoRaw` generalization (line 166) changes the raw data structure from field names to flat-path keys (`map["eval.proposal"]["quick"]`). This requires `applyDefaults` to also change — it currently uses `applyModeDefault(&a.Test, a.raw, "test", d.Test)` with hardcoded field names. With nested eval config, it would need `applyModeDefault(&a.Eval.Proposal, a.raw, "eval.proposal", d.Eval.Proposal)`. The proposal mentions this in the scope (line 166) but the `applyDefaults` change isn't traced in the solution design.

### 3. Industry Benchmarking (88/120)

**Industry solutions referenced (30/40)**:
Two references: Viper and dotnet Microsoft.Extensions.Configuration. The Viper reference (line 120) is now substantive — discusses flat key map vs reflection, nil pointer handling, custom type formatting, raw tracking. Good improvement from baseline.

- **Deduction -5**: The dotnet reference remains thin: "基于 `:` 分隔的路径遍历，支持任意深度" and "类似 Viper 的 flat key 存储，同样缺乏 struct 类型感知". Two sentences, no depth. What does dotnet do about type coercion? Configuration binding? Environment variable override hierarchy? The comparison adds no actionable insight.
- **Deduction -5**: Missing reference to Go-native configuration libraries that use reflection for struct traversal. For example, `go-structconf`, `koanf`, or even `json-iterator`'s reflection-based approach. These are closer analogs than Viper (which deliberately avoids reflection for config access). The proposal uses Viper as the primary comparison but Viper's design philosophy (flat map) is fundamentally different from the proposed approach (reflective struct traversal).

**At least 3 meaningful alternatives (23/30)**:
Four alternatives listed. "Do nothing" is baseline. "扁平命名空间" is reasonable. "泛化路由 + 嵌套 eval 配置" is selected. "仅扩展路由支持三层" is the closest competitor.

- **Deduction -4**: "仅扩展路由支持三层" (line 116) now has quantification ("~30 行/字段") but the 30-line figure is unsourced. Looking at the existing code, adding a new field to `autoModeField` (1 case in switch) + a new `getEvalKeyValue` function (copy of `getAutoKeyValue` with eval-specific field mapping) + a new `setEvalConfigValue` function (copy of `setAutoConfigValue` with eval-specific handling) would be closer to 50-80 lines for the first eval field, not 30. The proposal's estimate may be optimistic, weakening the breakeven argument.
- **Deduction -3**: No "minimal viable alternative" — what's the smallest change that solves the eval problem without the routing rewrite? The "仅扩展路由" alternative comes close but is presented as dismissable tech debt. A genuine minimal alternative would be: just add `auto.evalProposal`/`auto.evalPrd`/etc. as flat ModeToggle fields in `AutoConfig`, extend `autoModeField` with 4 new cases, and call it done. This would take ~20 minutes. The proposal doesn't honestly present this as an option.

**Honest trade-off comparison (18/25)**:
The cons for the selected approach (line 115) now includes "影响所有现有 `forge config get/set` 调用路径，任何反射遍历 bug 会同时影响所有配置操作". This is a significant improvement from baseline.

- **Deduction -4**: The cons still understates the debuggability cost. Reflection-based code is significantly harder to debug than explicit dispatchers. When a config get returns unexpected results, the debugging path goes through reflect.Value indirection rather than direct function calls. Stack traces are less informative. This is a real operational cost not mentioned.
- **Deduction -3**: The "扁平命名空间" alternative (line 114) is dismissed as "不解决路由瓶颈，命名不直观". But this alternative *does* solve the eval problem (the actual user-facing problem) — it just doesn't solve the routing architecture problem (the developer-facing problem). The dismissal conflates the two problems and undervalues an alternative that addresses the primary user need.

**Chosen approach justified (17/25)**:
The justification is now multi-faceted: zero marginal cost per new field (line 74), platform thinking vs feature thinking (line 77), breakeven analysis (line 116).

- **Deduction -5**: The breakeven analysis (line 116) assumes "未来知识管理和验证自动化预计 3-4 字段" but provides no roadmap or commitment. If those fields don't materialize, the generic routing is over-engineering for the current need. The breakeven is *projected* not *guaranteed*.
- **Deduction -3**: The "platform thinking" argument (line 77) is valid but self-serving. Every infrastructure rewrite can be justified as "platform thinking." The question is whether the current 6-field config system warrants a platform-grade routing solution. The proposal doesn't acknowledge that the system works fine for its current scale.

### 4. Requirements Completeness (92/110)

**Scenario coverage (32/40)**:
10 key scenarios cover happy paths and error paths (invalid key depth, invalid key path). Mode detection with fallback specified.

- **Deduction -4**: Missing scenario: concurrent config modification. Two `forge config set` calls running simultaneously (e.g., two terminal windows) both read the config, modify it, and write it back. The last write wins, losing the first write's changes. The current code has this same issue, but the proposal is extending the system — should it address this? At minimum, the scenario should be acknowledged as out-of-scope.
- **Deduction -4**: Missing scenario: config file with unexpected YAML structure. What if the YAML has `auto.eval.proposal: "string instead of mapping"`? The reflection path would encounter a string where it expects a struct. Error handling for type mismatches at intermediate nodes is undefined.

**Non-functional requirements (35/40)**:
Three NFRs: backward compatibility, config hot-reload, performance. Performance NFR is now quantified: "GetConfigValue 端到端延迟 <1ms" (line 98). Hot-reload is correctly identified as trivially satisfied (per-invocation re-read).

- **Deduction -3**: The performance NFR says "<1ms" but doesn't specify measurement conditions. Is this on a cold start (first reflection call) or warm (subsequent calls)? Reflection's first call for a type involves slower path lookups. On a typical CLI invocation, the first `forge config get` may be slower than <1ms due to type initialization overhead.
- **Deduction -2**: Missing NFR: error message quality. When `forge config set auto.eval.proposal.maybe true` fails, what does the user see? The proposal says "类型错误" but not the exact message. For a user-facing CLI, error message quality is a usability NFR.

**Constraints & dependencies (25/30)**:
Dependencies correctly identified: Go reflect (stdlib), yaml.v3 (existing), mode detection API (new from `feature_complete.go`), CLI subprocess dependency for skills.

- **Deduction -3**: The skill-to-CLI subprocess dependency (line 105) is a significant architectural constraint. Each skill spawns `forge config get` as a subprocess (~50-100ms). In a quick pipeline with 4 eval-triggering skills, that's 200-400ms of subprocess overhead for config checks alone. This latency is acknowledged ("约 50-100ms 延迟") but not treated as a constraint — it's a design choice with performance implications that could be avoided (e.g., skills reading config.yaml directly via Bash/yq).
- **Deduction -2**: The mode detection depends on `feature_complete.go`'s logic (checking proposal.md existence in feature directory). This logic is 4 lines (lines 106-109 in `feature_complete.go`). But the proposal says `forge config get mode` — meaning the mode detection needs to know which feature directory the user is in. How does the CLI command determine the current feature? The current `feature_complete.go` receives `featureSlug` as a parameter. The proposal doesn't specify how the CLI `get mode` command determines the feature slug context.

### 5. Solution Creativity (35/100)

**Novelty over industry baseline (12/40)**:
The proposal explicitly states "eval 配置复用已有 ModeToggle 模式，无新概念" (line 75). The generic routing via reflection is standard Go practice — Viper, koanf, and dozens of other libraries do similar things. The combination (generic routing + eval config in one proposal) is mildly novel but not creative.

- **Deduction -8**: The reflection-based routing is a standard Go technique presented as if it's innovative ("反射路由是平台思维而非功能思维" — line 77). It's good engineering, not creativity. The "platform thinking" framing is self-congratulatory — every infrastructure improvement can be called "platform thinking." The proposal should acknowledge this is straightforward application of well-known techniques.
- **Deduction -20**: Zero novelty in the eval configuration design. It's a direct copy of `auto.runTasks` pattern applied to eval.

**Cross-domain inspiration (8/35)**:
No cross-domain inspiration is evident. ModeToggle from same codebase. Reflection from Go stdlib. Generic routing from standard config library patterns.

- **Deduction -27**: No borrowing from other domains. The proposal could have drawn from: database query planning (cost-based routing), network routing tables (longest prefix match), compiler symbol tables (scoped name resolution), or even IDE configuration systems (layered settings). None of these are referenced or leveraged.

**Simplicity of insight (15/25)**:
The insight "use reflection to replace hardcoded dispatchers" is the obvious answer to "how do we make config routing extensible?" The observation that eval config and routing generalization should be bundled is practical.

- **Deduction -10**: The insight is straightforward, not deep. A competent Go developer would arrive at the same solution within minutes of seeing the problem. The proposal correctly identifies the bottleneck but the solution is standard, not insightful.

### 6. Feasibility (80/100)

**Technical feasibility (32/40)**:
Go reflect can do everything described. The Get path with YAML tag priority and nil pointer handling is well-designed. The Set path using struct marshal is consistent with existing code. Function signatures are correct.

- **Deduction -4**: The `SurfacesMap` custom type (with `UnmarshalYAML`/`MarshalYAML`) is a potential issue. If `getStructValueByPath` encounters a `SurfacesMap` field via reflection, it sees a `map[string]string` — but the YAML representation is custom (can be a string or a mapping). The reflection router would need to handle this custom type specially, or skip it. The proposal doesn't address how the generic router handles fields with custom YAML marshaling.
- **Deduction -4**: The `CoverageConfig.ByType` uses `yaml:",inline"` tag (line 99), which means the map entries are direct children of the `coverage` node in YAML, not nested under a `byType` key. The reflection router would need to understand the `,inline` tag semantics to correctly resolve `coverage.coding.feature` as "find the ByType map (which is inlined) and look up 'coding.feature' key". This is a non-trivial reflection edge case.

**Resource & timeline feasibility (24/30)**:
"预计 3-4 小时" — the breakdown: 2h for generic routing rewrite + 1h for eval config + 1h for testing.

- **Deduction -4**: The generic routing rewrite touches: `GetConfigValue`, `SetConfigValue`, `getAutoKeyValue`, `getWorktreeKeyValue`, `getCoverageKeyValue`, `setAutoConfigValue`, `setWorktreeConfigValue`, `setCoverageConfigValue`, `autoModeField`, `parseAutoRaw`, `applyDefaults`, `Config` struct, `AutoConfig` struct — that's 12+ functions. Plus handling edge cases for `SurfacesMap`, `CoverageConfig.ByType` with `,inline` tag, `WorktreeConfig.CopyFiles` ([]string type), and the new `EvalConfig` nested struct. In 2 hours, that's 10 minutes per function. Optimistic.
- **Deduction -2**: Testing estimates 1 hour for all tests including: reflection-based get/set for all existing paths (worktree, coverage, auto), new eval paths, mode detection API, parseAutoRaw generalization, applyDefaults with nested structs. This is tight.

**Dependency readiness (24/30)**:
Go reflect (stdlib) and yaml.v3 (existing) are ready. Mode detection API needs extraction from `feature_complete.go`.

- **Deduction -3**: The mode detection API extraction is described as "提取 `feature_complete.go` 中的 quick mode 判定逻辑为独立函数" (line 104). But `feature_complete.go`'s logic (lines 106-109) checks for `proposal.md` existence — it doesn't know the current feature slug. The `forge config get mode` CLI command needs context: which feature directory is the user in? The proposal doesn't specify how the CLI determines the feature context (environment variable? working directory detection? command argument?). This is an unsolved design problem.
- **Deduction -3**: The `yaml:",inline"` tag on `CoverageConfig.ByType` (line 99 of config.go) means the map is serialized inline, not under a `byType` key. The current `getCoverageKeyValue` handles this with custom logic. The generic reflection router would need special handling for `,inline` tagged fields — but the proposal doesn't mention this at all. This is a hidden dependency on understanding yaml.v3's inline semantics.

### 7. Scope Definition (72/80)

**In-scope items are concrete (27/30)**:
In-scope items name specific functions and code changes. The `parseAutoRaw` generalization now includes the flat-path key format (line 166). Mode detection API is fully specified with return values (line 171).

- **Deduction -3**: "4 个 skill 增加 config check" (line 170) is still coarse. Which specific steps in each skill are modified? The brainstorm skill has 7 steps — does the config check go between step 6 (commit) and step 7 (eval prompt)? Does it replace step 7 entirely? The exact modification points in the skill markdown files aren't specified.

**Out-of-scope explicitly listed (22/25)**:
Seven out-of-scope items. The "forge guide 文档更新" now correctly says "后续迭代" (line 179) without the self-contradictory "随代码一起完成".

- **Deduction -3**: Missing from out-of-scope: rollback plan for Part 1. If the generic routing rewrite has a regression, what's the rollback strategy? Revert to hardcoded dispatchers? The Delivery Strategy mentions splitting into 2 PRs for independent rollback, but doesn't specify the rollback procedure for PR-1 itself.

**Scope is bounded (23/25)**:
Scope is well-bounded with the Delivery Strategy providing clear separation (PR-1 + PR-2).

- **Deduction -2**: The mode detection API (line 171) adds significant scope beyond pure config routing. It requires: feature directory detection logic, proposal.md existence check, CLI command registration, and 3 skill integrations. This is effectively a mini-feature embedded in the scope. The proposal acknowledges this but doesn't treat it as a separate concern requiring its own design.

### 8. Risk Assessment (75/90)

**Risks identified (25/30)**:
6 risks identified. The risk table now includes map type serialization complexity (line 200), parseAutoRaw raw tracking (line 202), and skill-to-CLI subprocess dependency (line 204). Good coverage of the key risks.

- **Deduction -3**: Missing risk: `yaml:",inline"` tag handling. `CoverageConfig.ByType` uses inline tag — the reflection router must handle this correctly. If it doesn't, `coverage.coding.feature` lookups break silently. This is a real implementation risk not identified.
- **Deduction -2**: Missing risk: `SurfacesMap` custom marshaling interaction with reflection. `SurfacesMap` has `UnmarshalYAML`/`MarshalYAML` — if the reflection router encounters this type, it may not handle the custom representation correctly.

**Likelihood + impact rated (25/30)**:
Ratings are reasonable. Part 1 regression at M/H is appropriate. Map type serialization at M/H is correct.

- **Deduction -3":** "4 个 skill 的 config check 逻辑漂移" at L/M — this seems underrated. The skills are markdown files executed by AI agents. There's no compile-time or runtime enforcement of consistency. The "EXTREMELY-IMPORTANT 标注" mitigation is effectively a comment — no stronger than a TODO. This should be M/M at minimum.
- **Deduction -2**: "parseAutoRaw 泛化后 raw tracking 精度变化" at M/M is appropriate (improved from L/M in baseline). But the mitigation "新增测试用例 `TestParseAutoRaw_EvalConfig`" names a test but doesn't specify what it verifies. Is it: (a) raw map has correct keys, (b) applyDefaults only supplements missing fields, (c) defaults don't override explicit settings? All three?

**Mitigations are actionable (25/30)**:
The 2-PR delivery strategy is concrete. Map type formatter registration is specific. CLI fallback with `go install` in CI is actionable.

- **Deduction -3**: "为 map[string]struct 类型注册类型感知格式化器" (line 200) introduces a new concept ("type-aware formatters") that isn't part of the generic routing design. This is scope expansion hiding in a risk mitigation. How are formatters registered? A global map? A registry pattern? This design decision is deferred to implementation.
- **Deduction -2**: The YAML comment preservation mitigation (line 201) says "当前所有 set 路径已有此特性，不是新增风险" — this is correct but the proposal should acknowledge that the reflection-based set path will *also* lose comments, and this is a known limitation that users may encounter when using `forge config set` for the first time with the new system.

### 9. Success Criteria (72/80)

**Criteria are measurable and testable (26/30)**:
PR-1 has 11 SCs covering get, set, error paths, and regression. PR-2 has 12 SCs covering mode detection, skill behavior, and defaults. Most are directly testable via CLI commands.

- **Deduction -2**: "brainstorm 在 `auto.eval.proposal` 对应模式为 true 时跳过 AskUserQuestion" (line 227) — this is a behavior test on a markdown skill file. How is it tested? The skill is markdown executed by an AI agent; there's no automated test for "skip AskUserQuestion". This criterion is manually verifiable only.
- **Deduction -2**: "parseAutoRaw 对 `auto.eval.*` 字段生成正确的 flat-path raw map" (line 219) — "正确" is undefined. What does a "correct" flat-path raw map look like for `auto.eval.proposal` with both quick and full set? The example `map["eval.proposal"]["quick"]=true` shows one entry but doesn't show the full expected output for a complete eval config.

**Coverage is complete (23/25)**:
SCs cover: get (3, 4 levels), set (happy/error), error paths (excess depth, nonexistent key), regression (worktree, coverage), parseAutoRaw, mode detection (3 cases), skill behavior (8 scenarios: 4 skills x 2 modes), defaults.

- **Deduction -2**: Missing SC for intermediate node get format. Line 212 describes `forge config get auto.eval` output format ("每行一个字段") but there's no SC verifying this specific output format. The SC at line 210 says `forge config get auto.eval.proposal` returns `quick:true full:true` but `forge config get auto.eval` has no SC.

**SC internal consistency (23/25)**:
SCs are split into PR-1 and PR-2. Internal logic is consistent.

- **Deduction -2**: PR-1 SC line 210 says `forge config get auto.eval.proposal` returns `quick:true full:true` (three levels deep). But `auto.eval.proposal` is a ModeToggle, which has sub-fields (quick, full). The SC treats this as a "leaf" format. However, the Set path SC (line 213) rejects setting `auto.eval` (also non-leaf). The consistency question: is `auto.eval.proposal` a leaf or non-leaf? It's a leaf for set purposes (can't set it to a boolean) but returns a formatted value for get. This asymmetry is functional but should be explicitly noted.

### 10. Logical Consistency (78/90)

**Solution addresses stated problem (30/35)**:
The solution directly addresses both problems. Generic routing solves the extensibility bottleneck. Eval configuration solves the interaction friction. The mode detection API enables quick/full differentiation.

- **Deduction -3":** The problem says "eval 的手动确认成为剩余的主要交互摩擦点" but the solution doesn't eliminate all manual interaction by default. It makes proposal and uiDesign automatic (true, true) while keeping prd and techDesign manual (false, false). The problem framing implies all evals should be configurable (achieved) but the urgency argument implies they should all be automated (not achieved). The gap: the problem is stated as "消除交互摩擦" but the solution is "make interaction configurable" — these are different goals.
- **Deduction -2**: The Assumptions Challenged table (line 157) claims "Refuted" for "proposal 默认应询问用户" with the argument "proposal 是流水线入口，自动评估可尽早发现问题". This is an opinion presented as a finding. Whether auto-evaluating proposals catches more issues than it wastes resources on is an empirical question, not a logical refutation.

**Scope <-> Solution <-> SC aligned (25/30)**:
Scope includes mode detection API — Solution describes it — SC has 3 criteria for it. Scope includes parseAutoRaw generalization — Solution describes flat-path keys — SC has criterion. PRs are mapped to SCs.

- **Deduction -3**: Scope includes "泛化 `parseAutoRaw`" (line 166) with flat-path key definition. The SC (line 219) verifies output for `auto.eval.*` fields. But the scope also says "递归扫描 auto 子树" — this implies the generalization should work for *all* auto fields, not just eval. There's no SC verifying that existing fields (test, consolidateSpecs, etc.) still work correctly with the generalized parseAutoRaw. The regression SCs (lines 217-218) test `get`/`set` but not `parseAutoRaw` for existing fields.
- **Deduction -2**: The Delivery Strategy (line 186-192) recommends 2 PRs but doesn't specify the integration testing strategy between PR-1 and PR-2. PR-2 depends on PR-1's generic routing. How is this dependency tested? Is there a PR-1 acceptance test that verifies the generic routing works for arbitrary nested structs before PR-2 starts?

**Requirements <-> Solution coherent (23/25)**:
Requirements and solution are well-aligned. 10 key scenarios map to solution components. Default values match backward-compatibility requirements.

- **Deduction -2**: The NFR "配置热生效：修改配置后无需重启即可生效" (line 97) is trivially satisfied (CLI tools don't have restart semantics). This NFR doesn't actually constrain the solution — it's a non-requirement taking up space. The performance NFR ("<1ms") is more meaningful but only constrains the implementation, not the design.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 95 | 110 |
| Solution Clarity | 103 | 120 |
| Industry Benchmarking | 88 | 120 |
| Requirements Completeness | 92 | 110 |
| Solution Creativity | 35 | 100 |
| Feasibility | 80 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 75 | 90 |
| Success Criteria | 72 | 80 |
| Logical Consistency | 78 | 90 |
| **Total** | **790** | **1000** |

## Deductions Applied

- Vague language without quantification (3 instances, -20 pts each): -60 pts
  1. "~60-80s 交互时间" — estimation methodology unstated (Problem Definition, urgency)
  2. "未来知识管理和验证自动化预计 3-4 字段" — no roadmap or commitment (Industry Benchmarking, justification)
  3. "~30 行/字段" — unsourced estimate, may be optimistic (Industry Benchmarking, alternatives)
- Self-congratulatory language (1 instance): -20 pts
  1. "反射路由是平台思维而非功能思维" — standard technique presented as innovative (Solution Creativity, novelty)

## Attack List

1. **[Problem Definition]** Urgency for generic routing rewrite is ungrounded — "泛化 config key resolution 是一次投资长期收益的基础设施改进" — this is a generic infrastructure argument that could justify the rewrite at any time. No trigger event explains why *now*. The eval config need can be met with ~30 lines of hardcoded extension. The 7-field breakeven projection has no committed roadmap.

2. **[Solution Clarity]** Mixed-type intermediate node get behavior undefined — `forge config get auto` would need to enumerate `Test` (ModeToggle), `GitPush` (bool), `Eval` (nested struct), etc. Formatting rules for this mixed output are unspecified. Only `auto.eval` (uniform ModeToggle children) is described.

3. **[Solution Clarity]** Skill-side config check implementation undefined — "4 个 skill 增加 config check" is in scope but the actual skill markdown modification pattern isn't shown. What does the config check look like in the skill file? A bash block? A rules directive? This makes PR-2 implementation ambiguous.

4. **[Solution Clarity]** Nested ModeToggle set behavior unclear — can `forge config set auto.eval.proposal true` set both quick and full (like current `auto.{field}` behavior)? The proposal says setting `auto.eval` is rejected (non-leaf), but doesn't address `auto.eval.proposal` (also non-leaf in the struct tree, but a leaf ModeToggle in the current system's mental model).

5. **[Industry Benchmarking]** dotnet reference is decorative — "基于 `:` 分隔的路径遍历，支持任意深度" adds zero actionable insight. Either deepen the comparison or remove it.

6. **[Industry Benchmarking]** No Go-native reflection-based config library reference — koanf, go-structconf, or even viper's `Unmarshal` (which *does* use reflection) are closer analogs. Viper's `Get()` uses flat map, but its `Unmarshal()` uses reflection — this distinction is lost.

7. **[Industry Benchmarking]** Minimal viable alternative not honestly presented — adding 4 flat ModeToggle fields to `AutoConfig` with ~20 minutes of `autoModeField` extension would solve the eval problem without any routing rewrite. This option isn't presented because the proposal bundles both problems. The user-facing eval problem and the developer-facing routing problem should be independently solvable.

8. **[Requirements Completeness]** `yaml:",inline"` tag on `CoverageConfig.ByType` unaddressed — the map is serialized inline under `coverage:`, not under `coverage.byType:`. The generic reflection router must understand this tag to correctly resolve `coverage.coding.feature`. This is a non-trivial implementation requirement not mentioned anywhere in the proposal.

9. **[Requirements Completeness]** Mode detection CLI context mechanism undefined — `forge config get mode` needs to know the current feature directory. `feature_complete.go` receives `featureSlug` as a parameter. How does the CLI command determine which feature directory the user is in? Working directory parsing? Environment variable? Command argument?

10. **[Feasibility]** `SurfacesMap` custom marshaling interaction with reflection unaddressed — `SurfacesMap` has `UnmarshalYAML`/`MarshalYAML` for custom YAML representation. The reflection router encountering this type would see `map[string]string` but the YAML representation is custom (can be scalar or mapping). How does the generic router handle fields with custom YAML types?

11. **[Feasibility]** Timeline is optimistic — 12+ functions to rewrite/delete/modify in 2 hours for the routing rewrite, plus edge cases for `SurfacesMap`, `CoverageConfig.ByType` (inline), `WorktreeConfig.CopyFiles` ([]string), and nested EvalConfig. 10 minutes per function including testing is tight.

12. **[Risk Assessment]** Missing risk: `yaml:",inline"` tag handling — the reflection router must correctly handle inline-tagged fields. If it doesn't, `coverage.*` lookups break silently. This is a concrete implementation risk not in the risk table.

13. **[Risk Assessment]** Missing risk: `SurfacesMap` custom type interaction — custom marshal/unmarshal types may not behave correctly with reflection-based traversal.

14. **[Success Criteria]** Skill behavior SCs are not automatable — "brainstorm 在 `auto.eval.proposal` 对应模式为 true 时跳过 AskUserQuestion" is a manual test on AI-executed markdown. No automated verification method specified.

15. **[Success Criteria]** Missing SC for `forge config get auto.eval` intermediate node format — line 212 describes the format but no SC verifies it. The closest SC (line 210) only covers `auto.eval.proposal` (3 levels), not `auto.eval` (2 levels).

16. **[Logical Consistency]** Problem framing mismatch — problem says "消除交互摩擦" but solution is "make interaction configurable". Default config still requires manual confirmation for prd and techDesign. The urgency argument (60-80s wasted) only materializes if those defaults change.

17. **[Logical Consistency]** parseAutoRaw generalization regression gap — scope says generalization works for all auto fields, but SC only verifies eval-specific raw tracking. No SC confirms existing fields (test, consolidateSpecs, etc.) produce correct raw data after generalization.

18. **[Logical Consistency]** Assumptions Challenged "Refuted" is opinion — "proposal 默认应询问用户" is marked "Refuted" with "自动评估可尽早发现问题". This is a design choice, not a logical refutation. No evidence that auto-evaluating proposals catches more issues than it wastes.
