# Eval Report: Iteration 1

**Reviewer**: CTO Adversary (Rubric Scoring)
**Document**: `docs/proposals/auto-eval-config/proposal.md` (pre-revised version)
**Date**: 2026-05-26

---

## Annotated Region Analysis

### Pre-revised regions (5 annotations)

| Marker | Region Summary | Pre-revision Direction | Revision Quality |
|--------|---------------|----------------------|------------------|
| `<!-- pre-revised: medium -->` Get path | Changed from bare "反射" to detailed reflection rules (YAML tag priority, nil pointer handling, leaf type formatting) | Freeform finding #5: YAML tag matching + nil pointer undefined | Revision adequately addresses the finding. No new issues introduced. |
| `<!-- pre-revised: medium -->` Set path | Changed from "YAML Node" to "Go struct marshal (反射 set)" with rationale | Freeform finding #4: YAML Node vs struct marshal ambiguity | Revision resolves the ambiguity. Explicit rationale for choosing struct marshal is sound. |
| `<!-- pre-revised: high -->` ModeToggle mode detection | Changed from "manifest 文件检测" to "CLI 级别 mode 检测 (forge config get mode)" | Freeform finding #3: manifest mode field doesn't exist | Critical fix. Revision introduces new dependency: mode detection API must be built. |
| `<!-- pre-revised: medium -->` parseAutoRaw | Added flat path key definition for raw data structure | Freeform finding #6: raw map key format undefined | Partial fix. Defines key format but doesn't show how `applyDefaults` will consume it. |
| `<!-- pre-revised: high -->` Delivery Strategy | New section: recommends splitting into 2 PRs | Freeform finding #1: coupled delivery risk | Good addition. Addresses the core delivery risk directly. |

### Attack density

- **Annotated regions**: 3 attacks (see dimensions below) — 0.6 attacks per annotation
- **Unannotated regions**: 15 attacks — denser criticism on unrevised content
- **Bias assessment**: No systematic over-attack of pre-revised regions. Pre-revised sections generally show higher quality; attacks focus on residual gaps and unrevised areas.

---

## Dimension Scoring

### 1. Problem Definition (88/110)

**Problem stated clearly (35/40)**:
The dual-problem framing is effective: interaction friction from eval skills + config routing extensibility bottleneck. Both problems are real and verified against code. Minor deduction: the two problems have different scopes and urgency levels, but the proposal presents them as equally weighted. The eval interaction friction is a UX convenience; the routing bottleneck is an architectural issue. Conflating them risks obscuring the actual value proposition.

**Evidence provided (35/40)**:
Evidence is concrete and verifiable. `strings.SplitN(rest, ".", 2)` at line 598 of `config.go`, `autoModeField` switch/case at lines 515-531, `parseAutoRaw` hardcoded `modeFields` at line 383 — all confirmed by code review. The `AskUserQuestion` pattern in 3 of 4 skills is verified. Minor gap: no quantification of the actual "interaction cost" — how many seconds/clicks does each `AskUserQuestion` cost? Without this, the problem severity is qualitative only.

**Urgency justified (18/30)**:
The urgency argument is the weakest part of this section. "eval 的手动确认成为剩余的主要交互摩擦点" — this is asserted without measurement. "一次投资长期收益的基础设施改进" — the generic routing rewrite would be valuable regardless of eval configuration; tying urgency to eval is circular. The routing bottleneck has existed for as long as the hardcoded dispatchers have been there — why is *now* the right time? No trigger event is cited.

### 2. Solution Clarity (98/120)

**Approach is concrete (38/40)**:
The two-part structure is clear. Part 1 (generic routing) and Part 2 (eval config) are well-separated. The pre-revised Get/Set paths with detailed reflection rules are concrete enough to implement from. Function signatures are provided. The "消除的硬编码" list gives a clear deletion scope. Minor gap: the transition from current code to new code isn't traced — which functions are deleted, which are rewritten in-place, which are new?

**User-facing behavior described (38/45)**:
CLI behavior is well-described: `forge config get auto.eval.proposal` returns `quick:true full:true`, `forge config set auto.eval.prd.full true` writes nested config. The 7 key scenarios cover the main use cases. Gap: error messages are undefined. What does `forge config get auto.eval.nonexistent` return? What about `forge config set auto.eval.proposal maybe`? The success criteria partially address this (intermediate node behavior for get/set), but edge cases around invalid inputs are missing. The mode detection API (`forge config get mode`) behavior when not in a feature directory is unspecified.

**Technical direction clear (22/35)**:
The pre-revised Get path (reflection) and Set path (struct marshal) are clear and internally consistent. However, the proposal describes Set path as "反射 set" in the heading, then says "使用现有的 `readOrCreateConfig` + 反射 set + `writeConfig`" — but also mentions "YAML Node 操作仅在 get 路径中用于精确定位节点位置". This is confusing: the Set path uses reflection to navigate Go structs and then `yaml.Marshal` to write back, never touching YAML Nodes. Why mention YAML Node in the Set context at all? The baseline version said Set was YAML Node; the revision changed it to struct marshal but left a confusing YAML Node reference. Deduct for this residual ambiguity.

Additionally, the `getStructValueByPath` signature says it returns `(string, error)` — but for intermediate nodes (e.g., `auto.eval`), the success criteria say it should return a multi-line summary. This means the return type is effectively "string that might be multi-line with a specific format", which is a contract that should be explicit in the function design.

### 3. Industry Benchmarking (78/120)

**Industry solutions referenced (25/40)**:
Two references: Viper (Go) and dotnet Microsoft.Extensions.Configuration. Both are relevant. However, the references are shallow — "Viper 支持 `viper.Get("a.b.c")`" is a one-liner without discussing Viper's actual implementation approach (Viper uses its own flat key map, not reflection), its handling of nested structs, or how it deals with the same challenges (nil pointers, map types, type coercion). The dotnet reference is even thinner. No comparison of how these systems handle the specific challenges the proposal faces (ModeToggle formatting, CoverageStrategy map serialization, raw tracking for defaults).

**At least 3 meaningful alternatives (25/30)**:
Four alternatives are listed in the comparison table. "Do nothing" is a legitimate baseline. "扁平命名空间" is a reasonable alternative. "泛化路由 + 嵌套 eval 配置" is the selected approach. "仅扩展路由支持三层" is the closest to a straw-man — it's described as "渐进式技术债" without quantifying the actual cost. How many lines of code would it take to extend `autoModeField` with an `eval` case? The proposal doesn't say. Without this quantification, the reader can't judge whether the "技术债" is real or hypothetical.

**Honest trade-off comparison (18/25)**:
The cons column for the selected approach says "路由层重写的一次性投入" — this understates the risk. A routing layer rewrite affects *all* existing config operations. The trade-off should mention: (a) every existing `forge config get/set` call path changes, (b) any bug in the generic router breaks all config operations simultaneously, (c) the reflection-based approach is harder to debug than the explicit dispatchers. The Delivery Strategy section partially addresses this, but the comparison table itself doesn't reflect the full risk.

**Chosen approach justified (10/25)**:
The justification is essentially "it's the best long-term option" — which is true but unconvincing without addressing *when* the extensibility benefit would be realized. The proposal says "未来新增任何 `auto.*` 配置只需加 Go struct 字段 + 默认值" — but how often does Forge add new auto config fields? Looking at the history, the current 6 fields (`test`, `consolidateSpecs`, `cleanCode`, `validation`, `runTasks`, `knowledgeSave`) were added over several months. Adding one more (eval) would take maybe 15 minutes with the current hardcoded approach. The generic router saves 15 minutes per new field, at a cost of a 2-hour rewrite + regression risk. The breakeven is at ~8 new fields. Is there a roadmap for 8+ new auto config fields? The proposal doesn't say.

### 4. Requirements Completeness (85/110)

**Scenario coverage (32/40)**:
7 key scenarios cover the happy paths. Missing scenarios:
- What happens when `forge config set` is called with an invalid key depth (e.g., `auto.eval.proposal.quick.extra`)?
- What happens when a skill is invoked outside a feature directory (no mode detection possible)?
- What happens when the config file is corrupted or contains unexpected types?
- What happens when two concurrent processes try to set config values simultaneously?

**Non-functional requirements (28/40)**:
Three NFRs listed: backward compatibility, hot reload, performance. "配置热生效：修改配置后无需重启即可生效" — this is trivially true because Forge CLI is a command-line tool invoked per operation; there's no long-running process to "restart". This NFR is a non-requirement dressed up as a requirement. The performance NFR says "反射开销在配置读取场景可忽略" — this is asserted without measurement. How many reflect calls per `GetConfigValue` invocation? What's the expected latency increase in milliseconds?

Missing NFRs: security (config file permissions), reliability (atomic writes, config corruption recovery), usability (error message quality for unknown keys).

**Constraints & dependencies (25/30)**:
Dependencies are correctly identified: Go reflect (stdlib), yaml.v3 (existing), mode detection API (new). The pre-revised constraint correctly identifies that mode detection requires extracting logic from `feature_complete.go`. Minor gap: the dependency on "4 个 skill 文件为 markdown，通过 Bash 调用 `forge config get`" is a significant architectural constraint — it means each skill spawns a subprocess to check config, which has latency implications and error-handling complexity. This isn't flagged as a risk.

### 5. Solution Creativity (35/100)

**Novelty over industry baseline (15/40)**:
The proposal explicitly positions itself as *not* novel: "eval 配置复用已有 ModeToggle 模式，无新概念". The generic routing is a standard technique (Viper does it). The eval configuration is a straightforward application of existing patterns. There's no novel algorithm, data structure, or interaction pattern. The only mild novelty is combining the two changes (generic routing + eval config) into a single proposal that justifies the routing rewrite.

**Cross-domain inspiration (10/35)**:
No cross-domain inspiration is evident. The ModeToggle pattern is from the same codebase (`auto.runTasks`). The generic routing is from standard Go libraries (reflect) and common configuration libraries (Viper). No inspiration from other domains (e.g., database query planning, network routing, UI component trees).

**Simplicity of insight (10/25)**:
The insight "泛化路由 + 嵌套配置" is clean but not particularly deep. It's the obvious answer to "how do we make config routing extensible?" The proposal correctly identifies that the existing hardcoded dispatchers are a bottleneck and proposes the standard solution (reflection-based generic traversal). The insight that eval configuration and routing generalization should be bundled is practical but not creative — it's just good engineering.

### 6. Feasibility (82/100)

**Technical feasibility (32/40)**:
Technically sound. Go reflect can do everything described. The pre-revised Get path with YAML tag priority and nil pointer handling is well-thought-out. The Set path using struct marshal is consistent with existing code. The main concern is the `CoverageConfig.ByType` map handling — the proposal acknowledges this in the risk table but the mitigation ("注册类型感知格式化器") introduces a new concept (type-aware formatters) that isn't part of the generic routing design. This is a small scope expansion hiding in a risk mitigation.

**Resource & timeline feasibility (25/30)**:
"预计 3-4 小时" — the baseline said "3-4 hours" and was revised to "3-4 hours" (no change). The freeform review previously flagged this. Breaking it down: 2h for generic routing rewrite + 1h for eval config + 1h for testing. The routing rewrite touches `GetConfigValue`, `SetConfigValue`, `getAutoKeyValue`, `getWorktreeKeyValue`, `getCoverageKeyValue`, `setAutoConfigValue`, `setWorktreeConfigValue`, `setCoverageConfigValue`, `autoModeField`, `parseAutoRaw`, `applyDefaults`, `Config` struct, and `AutoConfig` struct. That's 12+ functions to rewrite/delete/modify in 2 hours. Add the mode detection API extraction and the map type formatter registration. This is optimistic but not unreasonable for an experienced Go developer who knows the codebase well.

**Dependency readiness (25/30)**:
All dependencies exist except mode detection API. The proposal correctly identifies this as a new implementation need. The extraction from `feature_complete.go` is straightforward (the logic is 4 lines at lines 106-109). However, the proposal says "供 skill 通过 CLI 调用" — meaning skills will call `forge config get mode` as a subprocess. This creates a runtime dependency: skills must have the forge CLI built and installed to function. If the CLI isn't built, mode detection silently fails. This dependency chain isn't discussed.

### 7. Scope Definition (62/80)

**In-scope items are concrete (25/30)**:
In-scope items are specific and actionable. Each item names the affected code (e.g., "`autoModeField` switch/case", "`getAutoKeyValue`/`getWorktreeKeyValue`/`getCoverageKeyValue`"). The pre-revised `parseAutoRaw` item includes the flat path key format definition, which is helpful. Minor gap: "4 个 skill 增加 config check" doesn't specify what the config check looks like in the skill markdown — is it a bash command? A conditional block? A rules file?

**Out-of-scope explicitly listed (17/25)**:
Five out-of-scope items listed. "forge guide 文档更新（文档变更随代码一起完成）" — this is contradictory. If documentation changes happen with the code, they're *in scope*. Saying "out of scope but will happen" is confusing. Either the docs are part of this proposal (in scope) or they're deferred to a later proposal (out of scope). Also missing from out-of-scope: performance benchmarking, error message localization, config validation (are there invalid key patterns?), migration guide for users with custom configs.

**Scope is bounded (20/25)**:
The scope is generally well-bounded. The Delivery Strategy section helps by recommending 2 PRs. However, the mode detection API ("新增 `forge config get mode`") is a non-trivial addition that expands scope beyond config routing and eval configuration. It's a new CLI command that needs its own design decisions: what does it return when not in a feature directory? When in a feature directory without a manifest? When in a feature directory with both manifest and proposal? These edge cases make the mode detection a mini-feature of its own.

### 8. Risk Assessment (68/90)

**Risks identified (22/30)**:
5 risks identified. The pre-revised risk table is significantly better than the baseline — it addresses the CoverageConfig map complexity and the YAML Node set fidelity issue. Missing risks:
- Mode detection API failure modes (what if the CLI binary isn't built?)
- `applyDefaults` interaction with the new `EvalConfig` nested struct — the flat-path raw tracking may not integrate cleanly with the existing per-field `applyModeDefault` calls
- Breaking change risk: if external tools parse `forge config get` output, the new formatting for intermediate nodes could break them
- Testing gap: the proposal says "单元测试更新" but doesn't mention integration/E2E testing for the full skill → CLI → config → reflection chain

**Likelihood + impact rated (23/30)**:
Ratings are reasonable. The "Part 1 路由重写的回归风险" at M/H is appropriate. "反射遍历对 map 类型" at M/H is also appropriate after the pre-revision improvement. However, "parseAutoRaw 泛化后 raw tracking 精度变化" at L/M seems underrated — the `raw` data structure is central to the default value mechanism, and changing it could cause subtle bugs where defaults are applied incorrectly (false appearing as "user set" vs "missing"). This is M/M at minimum.

**Mitigations are actionable (23/30)**:
Mitigations are mostly actionable. "拆分为独立 PR" is concrete and already reflected in Delivery Strategy. "为 map[string]struct 类型注册类型感知格式化器" is specific. "递归扫描保持相同的叶子节点追踪粒度，测试验证默认值行为不变" — the first part is an implementation approach, the second part is a test strategy. But *which* tests verify default behavior? The existing `config_test.go` covers current defaults; new tests for the `EvalConfig` defaults with the new raw tracking format would need to be written. The mitigation should name specific test cases.

### 9. Success Criteria (62/80)

**Criteria are measurable and testable (22/30)**:
Most criteria are testable: `forge config get auto.eval.proposal` returns a specific string, `forge config set auto.eval.prd.full true` writes correctly. However:
- "brainstorm 在 `auto.eval.proposal` 对应模式为 true 时跳过 AskUserQuestion" — this is a behavior test on a markdown skill file. How is it tested? The skill is markdown executed by an AI agent; there's no automated test for "skip AskUserQuestion". This criterion is manually verifiable only.
- "write-prd/tech-design/ui-design 同理遵循配置驱动" — "同理" is vague. Each skill should have its own explicit criterion.
- The intermediate node get behavior criterion ("遍历该节点下所有叶子字段，按字段名逐行输出") is described in the SC list but the output format isn't fully specified. What's the separator? What order are fields listed?

**Coverage is complete (20/25)**:
The SC covers the main paths: get, set, skill behavior, defaults, regression. Missing:
- No SC for `parseAutoRaw` generalization — the raw tracking is a core mechanism change but has no dedicated criterion
- No SC for mode detection API behavior (`forge config get mode` returns what exactly?)
- No SC for error cases (setting invalid keys, getting nonexistent paths)
- No SC for the `applyDefaults` integration with EvalConfig

**SC internal consistency (20/25)**:
Generally consistent. One tension: the SC says `forge config get auto.eval.proposal` returns `quick:true full:true` (no spaces after colon), but the current `getAutoKeyValue` implementation returns `quick:%v full:%v` (with `fmt.Sprintf` which would produce `quick:true full:true` with spaces between key-value pairs but not after colons). The SC should match actual formatting convention. Also, the SC for intermediate node get ("proposal:quick:true full:true prd:quick:false full:false ...") uses colons as separators, which conflicts with the existing ModeToggle format "quick:true full:true". The nesting of separators could confuse parsing.

### 10. Logical Consistency (75/90)

**Solution addresses stated problem (30/35)**:
The solution directly addresses both stated problems. Generic routing solves the extensibility bottleneck. Eval configuration with ModeToggle solves the interaction friction. The mode detection API enables the quick/full differentiation. One logical gap: the problem statement says "每次交互都需要用户手动选择，增加了流水线的交互成本" — but the proposed solution doesn't eliminate all manual interaction. It makes *some* evals automatic by default (proposal, uiDesign) while keeping others manual (prd, techDesign). The problem implies all evals should be configurable, which the solution achieves, but the framing of the problem as "消除冗余交互" is misleading — it's more accurately "make eval interaction configurable".

**Scope <-> Solution <-> SC aligned (22/30)**:
Scope includes "mode 检测 API" — Solution Part 2 mentions mode detection — SC has no dedicated criterion for mode detection behavior. This is a gap in alignment.
Scope includes "泛化 `parseAutoRaw`" — Solution Part 1 mentions recursive scanning — SC has no criterion verifying the generalized parseAutoRaw produces correct raw data for nested eval fields.
The Delivery Strategy recommends splitting into 2 PRs, but the SC is written as a monolithic checklist — there's no mapping of which SCs belong to PR-1 vs PR-2.

**Requirements <-> Solution coherent (23/25)**:
Requirements and solution are well-aligned. The 7 key scenarios map directly to solution components. Default values in the solution match backward-compatibility requirements. Minor gap: the NFR "配置热生效" is trivially satisfied and thus doesn't actually constrain the solution in any meaningful way.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 88 | 110 |
| Solution Clarity | 98 | 120 |
| Industry Benchmarking | 78 | 120 |
| Requirements Completeness | 85 | 110 |
| Solution Creativity | 35 | 100 |
| Feasibility | 82 | 100 |
| Scope Definition | 62 | 80 |
| Risk Assessment | 68 | 90 |
| Success Criteria | 62 | 80 |
| Logical Consistency | 75 | 90 |
| **Total** | **733** | **1000** |

## Deductions Applied

- Vague language without quantification (4 instances, -20 pts each): -80 pts
  1. "增加了流水线的交互成本" — no quantification of cost (Problem Definition, urgency)
  2. "反射开销在配置读取场景可忽略" — no measurement (Requirements Completeness, NFR)
  3. "渐进式技术债" — no quantification of debt cost (Industry Benchmarking, alternatives)
  4. "配置热生效：修改配置后无需重启即可生效" — non-requirement stated as NFR (Requirements Completeness)
- Contradictory out-of-scope item (1 instance): -20 pts
  1. "文档变更随代码一起完成" is in-scope behavior listed as out-of-scope (Scope Definition)
- Straw-man alternative (1 instance): -20 pts
  1. "仅扩展路由支持三层" dismissed without quantifying implementation cost (Industry Benchmarking)
- Placeholder-like vagueness: "同理遵循配置驱动" (1 instance): -20 pts
  1. SC #9 uses "同理" instead of specifying each skill's expected behavior (Success Criteria)

## Attack List

1. **[Problem Definition]** Urgency is asserted, not demonstrated — "eval 的手动确认成为剩余的主要交互摩擦点" — provide quantitative evidence of the friction (e.g., "each AskUserQuestion adds ~30s of interaction time across 4 eval points per feature, totaling ~2min of manual clicks per pipeline run")

2. **[Solution Clarity]** Set path description contains residual YAML Node confusion — "YAML Node 操作仅在 get 路径中用于精确定位节点位置" — remove this sentence from the Set path section; it describes Get behavior and creates ambiguity about what Set actually does

3. **[Solution Clarity]** Error messages for invalid operations are undefined — the proposal describes happy-path get/set but never specifies what happens for `forge config set auto.eval.proposal.quick.extra true` (5 levels deep on a bool leaf) or `forge config get auto.nonexistent`

4. **[Industry Benchmarking]** Viper reference is superficial — "Viper 支持 `viper.Get("a.b.c")`" — explain how Viper handles the specific challenges in this proposal (nil pointer navigation, map[string]struct formatting, raw tracking for defaults) or remove the reference as decorative

5. **[Industry Benchmarking]** "仅扩展路由支持三层" alternative is a straw-man — dismissed as "渐进式技术债" without quantifying: how many lines of code to add an `eval` case to `autoModeField` and a new `getEvalKeyValue`? If it's 30 lines vs a 200-line reflection rewrite, the trade-off deserves honest numbers

6. **[Industry Benchmarking]** Breakeven analysis missing for routing rewrite — the generic router saves development time on future config additions, but the proposal doesn't estimate how many new config fields are planned. Without a roadmap showing N new fields, the "long-term benefit" is speculative

7. **[Requirements Completeness]** Performance NFR is unsubstantiated — "反射开销在配置读取场景可忽略" — provide a concrete bound: "GetConfigValue reflection overhead should remain under 1ms, vs current hardcoded dispatch at <0.1ms"

8. **[Requirements Completeness]** Missing scenario: mode detection outside feature directory — skills invoked outside a feature pipeline (e.g., standalone `/brainstorm`) would call `forge config get mode` and get... what? An error? "none"? This affects the config check logic in all 4 skills

9. **[Solution Creativity]** Creativity floor is acknowledged but not adequately addressed — the proposal scores itself at zero novelty ("复用已有 ModeToggle 模式，无新概念"), but the generic routing rewrite *is* the creative contribution. The proposal should better articulate why this approach is elegant rather than just standard

10. **[Scope Definition]** Out-of-scope item contradicts itself — "forge guide 文档更新（文档变更随代码一起完成）" — if docs change with the code, they are in scope. Pick one: either docs are out of scope (deferred) or they're in scope (tracked)

11. **[Scope Definition]** Mode detection API is an unscoped mini-feature — "新增 `forge config get mode`" requires: defining return values for 3+ contexts (in feature/quick, in feature/full, not in feature), error handling, and integration with 4 skills. This is significant enough to be a separate scope item with its own design

12. **[Risk Assessment]** `parseAutoRaw` risk is underrated — "raw tracking 精度变化" at L/M should be M/M. The raw data structure change from `map[string]map[string]bool` (flat field names) to flat-path keys ("eval.proposal") is a fundamental change to the default mechanism. Incorrect raw tracking causes silent wrong defaults

13. **[Risk Assessment]** Missing risk: skill-to-CLI subprocess dependency — 4 skills will call `forge config get` via Bash. If the CLI binary is stale or not built, config checks silently fail. No fallback behavior is defined

14. **[Success Criteria]** No SC for `parseAutoRaw` generalization — the raw tracking mechanism is being fundamentally changed but has no dedicated success criterion. Add: "parseAutoRaw produces correct raw map for auto.eval.* fields with flat-path keys"

15. **[Success Criteria]** No SC for mode detection API — `forge config get mode` is listed in scope but has no success criterion. Add: "`forge config get mode` returns 'quick' when proposal.md exists in feature directory, 'full' otherwise"

16. **[Success Criteria]** SC formatting inconsistency — `forge config get auto.eval.proposal` returns `quick:true full:true` (no spaces) but the intermediate node format uses `proposal:quick:true full:true` (colon-separated field name prefix). Define a consistent output format specification

17. **[Logical Consistency]** Scope-Solution-SC misalignment for mode detection — Scope lists "mode 检测 API", Solution describes it, but SC has no criterion. The three should be aligned: either add SC or remove from scope

18. **[Logical Consistency]** SC not mapped to PRs — Delivery Strategy recommends 2 PRs, but SC is a flat list with no PR attribution. SC items 1-7 (config get/set/regression) belong to PR-1; items 8-12 (skill behavior) belong to PR-2. This mapping should be explicit

## Conflict-with-pre-revision Tags

- **[Solution Clarity] Set path YAML Node reference**: The pre-revision direction was to resolve the YAML Node vs struct marshal ambiguity. The revision chose struct marshal but left a confusing YAML Node reference in the Set path description. This is a *conflict-with-pre-revision*: the revision's intent (clarity) was partially undermined by the leftover text. The YAML Node mention should have been removed entirely from the Set path section.
