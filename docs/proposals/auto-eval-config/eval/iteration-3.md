# Eval Report: Iteration 3 (FINAL)

**Reviewer**: CTO Adversary (Rubric Scoring)
**Document**: `docs/proposals/auto-eval-config/proposal.md`
**Date**: 2026-05-26

---

## Iteration-2 Issue Resolution Tracker

| # | Iter-2 Attack | Resolution Status | Assessment |
|---|---------------|-------------------|------------|
| 1 | Urgency for generic routing rewrite is ungrounded | **Unchanged** | Urgency still asserts "当前需求已到达现有路由的天花板" without establishing why NOW vs. later. The "3-4 future fields" projection remains uncommitted. |
| 2 | Mixed-type intermediate node get behavior undefined | **Resolved** | Lines 40-50 now fully specify `forge config get auto` output with mixed types: ModeToggle, bool, nested struct with indent. |
| 3 | Skill-side config check implementation undefined | **Resolved** | Lines 121-141 provide full bash template with fallback logic, error handling, and CLI unavailability path. |
| 4 | Nested ModeToggle set behavior unclear | **Resolved** | Line 59 explicitly states ModeToggle fields are non-leaf for set, returning error with guidance. |
| 5 | dotnet reference is decorative | **Resolved** | dotnet removed. Replaced with koanf (Go-native, closer analog). |
| 6 | No Go-native reflection-based config library reference | **Resolved** | koanf added (line 157) with specific differentiation on CLI get/set routing. |
| 7 | Minimal viable alternative not honestly presented | **Resolved** | Line 158 now includes "最小可行替代方案": 4 flat fields, 20 minutes, honest trade-offs. |
| 8 | yaml:",inline" tag on CoverageConfig.ByType unaddressed | **Resolved** | Addressed in 5 locations: solution (line 37), feasibility (line 170), risk table (line 241), SC (line 262), scope (line 200). |
| 9 | Mode detection CLI context mechanism undefined | **Resolved** | Line 120: CLI parses pwd to extract feature slug from `.forge/features/<slug>` pattern. |
| 10 | SurfacesMap custom marshaling interaction with reflection | **Resolved** | Risk table line 242: fallback to existing hardcoded path for `yaml.Unmarshaler` types. |
| 11 | Timeline is optimistic | **Partially resolved** | Timeline revised to 4-5 hours (from 3-4) with additional 1h for edge cases. Still tight for scope. |
| 12 | Missing risk: yaml:",inline" tag handling | **Resolved** | Risk table line 241 with specific test case `TestGetByPath_InlineMap`. |
| 13 | Missing risk: SurfacesMap custom type | **Resolved** | Risk table line 242 with fallback mitigation strategy. |
| 14 | Skill behavior SCs are not automatable | **Partially resolved** | Line 282 adds "Skill 行为 SC 代理验证" using code review as proxy method. |
| 15 | Missing SC for forge config get auto.eval intermediate node | **Resolved** | SC line 255 specifies intermediate node output format. |
| 16 | Problem framing mismatch | **Unchanged** | Problem still conflates "消除交互摩擦" with "make interaction configurable". |
| 17 | parseAutoRaw generalization regression gap | **Resolved** | SC line 265: regression test for existing auto fields after generalization. |
| 18 | Assumptions Challenged "Refuted" is opinion | **Unchanged** | "Refuted" for design opinions remains. |

**Summary**: 18 iteration-2 attacks tracked. 11 resolved, 3 partially resolved, 3 unchanged, 1 improved. The proposal has matured substantially since iteration 1 (733 -> 790 -> current).

---

## Dimension Scoring

### 1. Problem Definition (92/110)

**Problem stated clearly (36/40)**:
The dual-problem framing is well-structured. Eval interaction friction is a concrete UX issue with specific code references (AskUserQuestion in 3 of 4 skills, unconditional auto-eval in ui-design). Config routing extensibility bottleneck is a concrete architectural issue with specific code references (`strings.SplitN`, `autoModeField` switch, `parseAutoRaw` hardcoded list). The causal chain between the two problems (eval config needs 3-level nesting, which triggers the routing ceiling) is clearly stated.

- **Deduction -4**: The two problems have different scopes and natural solutions. The eval problem is solvable with 4 flat fields + 4 switch cases (~20 minutes, as the proposal now honestly acknowledges in line 158). The routing problem is an architectural investment that may or may not pay off. Bundling them creates a single proposal that is harder to evaluate on its merits. The reader must accept both problems as equally urgent, which they are not.

**Evidence provided (34/40)**:
Code-level evidence is strong: function names, line references, specific code patterns. `AskUserQuestion` in brainstorm/write-prd/tech-design verified. `ui-design` unconditional auto-eval verified. `SplitN(rest, ".", 2)` and `autoModeField` switch verified. `parseAutoRaw` hardcoded `modeFields` verified.

- **Deduction -3**: The "~60-80s 交互时间" estimate (line 24) remains an estimation without methodology. The breakdown ("LLM 生成提问 ~5-8s + 用户阅读选择 ~5-10s + LLM 解析回答 ~3-5s") is plausible but unverifiable. No timing logs, no user session data, no empirical measurement.
- **Deduction -3**: No user feedback evidence. Has anyone complained about the manual eval confirmation? Is there a support ticket, a Discord message, a user interview note? The problem is entirely developer-inferred.

**Urgency justified (22/30)**:
"eval 的手动确认成为剩余的主要交互摩擦点" -- the "剩余" framing (last remaining friction after other automations) is effective. The generic routing urgency is grounded in the specific technical ceiling: "三层嵌套意味着硬编码路由需要新引入 eval 子分发器（~30 行/字段）".

- **Deduction -4**: The routing urgency is still assumed rather than demonstrated. Line 152 quantifies the breakeven at "~5 个新字段" with "当前已规划 eval（4 字段）+ 未来知识管理和验证自动化预计 3-4 字段". "预计" is not "committed". The 3-4 future fields are aspirational, not planned. If they don't materialize, the routing rewrite was over-engineering.
- **Deduction -4**: No cost-of-delay. What happens if this proposal is deferred by 2 weeks? 1 month? The answer is: nothing catastrophic. The manual eval confirmation works. The hardcoded routing works. This is a quality-of-life improvement, not a time-critical fix.

### 2. Solution Clarity (108/120)

**Approach is concrete (38/40)**:
Part 1 (generic routing) and Part 2 (eval config) are cleanly separated. Function signatures provided (`getStructValueByPath`, `setStructValueByPath`). YAML tag matching priority, nil pointer handling, inline tag detection, leaf type formatting -- all specified. The "消除的硬编码" list (lines 61-65) gives clear deletion scope with the important nuance of "保留函数体作为 fallback" for custom types.

- **Deduction -2**: The proposal specifies non-leaf node get behavior with a detailed `auto` example (lines 40-50) but the map-of-struct formatting for `CoverageConfig.ByType` (inline map with `CoverageStrategy` values) is scattered across risk table (line 243) rather than specified in the solution design. The Get path design should include map-of-struct as a type case alongside bool, string, ModeToggle.

**User-facing behavior described (38/45)**:
CLI get/set behavior is thoroughly specified with 11 key scenarios including error paths. Error messages are concrete: "cannot set non-leaf key, use auto.eval.<field>.<subfield>", "invalid value \"maybe\" for bool field auto.eval.proposal.quick: expected true or false". The bash config check template (lines 121-141) with fallback logic is a significant improvement from earlier iterations.

- **Deduction -4**: The bash template (lines 122-141) uses `AskUserQuestion` as the fallback, but this function is a Claude Code MCP tool -- it requires the skill to be running inside Claude Code. What if the skill is invoked in a different context? The template assumes Claude Code as the execution environment without stating this constraint.
- **Deduction -3**: The `forge config get auto.eval` intermediate output format (line 255) specifies "字段按 Go struct 定义顺序排列" -- but Go struct field order is a compile-time detail that is not documented in any user-facing contract. If someone reorders the struct fields, the output changes silently. This is an implicit coupling between implementation detail and user-facing output.

**Technical direction clear (32/35)**:
Get path (reflection with YAML tag matching) and Set path (reflection + struct marshal) are clear and internally consistent. The inline tag handling design is specified (line 37). Custom YAML type fallback is specified (line 170). `parseAutoRaw` generalization with flat-path keys is described (line 203).

- **Deduction -2**: The `parseAutoRaw` generalization changes the raw map structure from field names to flat-path keys (`map["eval.proposal"]["quick"]`). The proposal says `applyDefaults` changes to flat-path calls -- but this is in a scope item (line 203), not in the Solution section. The solution design should trace how `applyDefaults` is modified to consume the new format for both new and existing fields.
- **Deduction -1**: The `WorktreeConfig.CopyFiles` ([]string) formatting is specified once (line 39) as "换行连接" but doesn't appear in the feasibility section's edge case list or in any SC.

### 3. Industry Benchmarking (96/120)

**Industry solutions referenced (33/40)**:
Three references: Viper, koanf, and a minimal viable alternative. The Viper reference (line 156) is substantive: discusses flat key map approach, nil pointer handling gap, custom type formatting gap, raw tracking gap, and the key distinction that Viper's `Unmarshal` is one-directional while this proposal needs bidirectional struct-to-CLI routing. The koanf reference (line 157) correctly identifies the closest analog: "koanf 侧重多源合并，不提供 CLI get/set 路由能力".

- **Deduction -4**: The Viper comparison still conflates `Get()` (flat map) with `Unmarshal()` (reflection-based struct mapping). The proposal says "Viper 的 `Get()` 用 flat map" -- correct -- but doesn't discuss Viper's `Unmarshal()` which uses reflection to map to structs. The distinction matters: Viper's reflection path is relevant to this proposal's approach, but Viper doesn't expose it as a CLI get/set API. The proposal reaches the right conclusion but the analytical path could be sharper.
- **Deduction -3**: No reference to other Go reflection-based CLI config systems. `go-flags`, `urfave/cli`, or Kubernetes' `client-go` flag binding all use reflection for struct-to-CLI mapping. These are closer analogs than Viper (a config file library) for the "struct IS the CLI API" pattern.

**At least 3 meaningful alternatives (26/30)**:
Four alternatives plus a minimal viable alternative (line 158). "Do nothing" is baseline. "扁平命名空间" is reasonable. "泛化路由 + 嵌套 eval 配置" is selected. "仅扩展路由支持三层" has quantification ("~30 行/字段"). The minimal viable alternative (20 minutes, 4 flat fields) is now honestly presented with "如果仅关注 eval 需求，此方案性价比最高".

- **Deduction -4**: The "仅扩展路由支持三层" breakeven analysis (line 152) says "~5 个新字段" but then notes "当前已规划 eval（4 字段）+ 未来知识管理和验证自动化预计 3-4 字段". The 4 eval fields are current need, not future. So the real breakeven comparison is: 4 fields needed NOW at 15 lines each = 60 lines incremental, vs 200 lines rewrite NOW. That is a 3.3x premium for the architectural investment, not a breakeven calculation. The "50 行首字段 + 15 行后续" for the incremental approach also deserves scrutiny -- the proposal itself (line 152) notes "~50-80 行：autoModeField 新增 eval case + getEvalKeyValue + setEvalConfigValue".

**Honest trade-off comparison (20/25)**:
The cons for the selected approach (line 151) now includes three honest drawbacks: "影响所有现有 forge config get/set 调用路径", "反射遍历 bug 会同时影响所有配置操作", "反射代码调试难度高于显式分发器". This is a significant improvement from earlier iterations.

- **Deduction -3**: The debuggability con is mentioned but understated. Reflection stack traces are significantly harder to interpret than direct function calls. In production debugging, a panic inside `reflect.Value.FieldByName` gives no indication of which config key was being resolved. The proposal says "stack trace 通过 reflect.Value 间接调用，不如函数调用直观" but doesn't propose any mitigation (e.g., wrapping reflection errors with key path context).
- **Deduction -2**: The "扁平命名空间" alternative (line 150) is dismissed as "不解决路由瓶颈，命名不直观". But this alternative solves the user-facing problem (eval config) perfectly, and the routing bottleneck is a developer concern. The dismissal undervalues a pragmatic solution that addresses the primary user need.

**Chosen approach justified (17/25)**:
Multi-faceted justification: zero marginal cost (line 87), platform thinking (line 91), breakeven analysis (line 152), honest minimal alternative (line 158).

- **Deduction -5**: The breakeven analysis still relies on "未来知识管理和验证自动化预计 3-4 字段" -- uncommitted future work. If those fields don't materialize, the 200-line investment is a net loss compared to the 60-80-line incremental approach. The proposal should either commit to a roadmap with those fields, or acknowledge the routing rewrite is an architectural preference, not an economic optimization.
- **Deduction -3**: The "platform thinking" framing (line 91) remains self-congratulatory: "反射路由是平台思维而非功能思维". This is a standard Go technique presented as strategic thinking. Every infrastructure rewrite can be called "platform thinking." The 6-field config system doesn't yet warrant platform-grade routing.

### 4. Requirements Completeness (97/110)

**Scenario coverage (35/40)**:
11 key scenarios cover happy paths (4-level get/set, mode detection), error paths (invalid key depth, invalid key path, type mismatch), edge cases (config file anomaly with type mismatch, mode detection outside feature directory), and behavior changes (ui-design unification).

- **Deduction -3**: Missing scenario: `forge config set auto.eval.prd.full true` when `auto.eval` section doesn't exist in config.yaml yet. Does the reflection-based set create the entire intermediate structure (auto -> eval -> proposal -> quick)? The Set path says "指针 -> nil 时自动初始化（reflect.New + Set）" in the feasibility section (line 176), but this scenario is not in the Key Scenarios.
- **Deduction -2**: Missing scenario: `forge config set` with a value that contains spaces or special characters. The bash template (line 129) uses `EVAL_ENABLED=$(forge config get ...)` -- what if the config value contains whitespace? The set path doesn't specify quoting or escaping rules.

**Non-functional requirements (36/40)**:
Three NFRs: backward compatibility, error message quality, config hot-reload, performance. Performance is now quantified with explicit measurement conditions (line 114): "GetConfigValue 端到端延迟 <1ms" with "CLI 进程启动本身约 20-50ms，反射遍历 <1ms 是相对于进程启动时间的可接受开销". Error message quality is now an explicit NFR (line 112) with specific examples.

- **Deduction -2**: The "配置热生效" NFR (line 113) remains trivially satisfied -- CLI tools are invoked per-operation, there is no persistent process to "reload". This is a non-requirement taking up space.
- **Deduction -2**: The performance NFR measurement conditions (line 114) state "CLI 正常运行态（非首次冷启动）" but then admit "CLI 每次调用是新进程，首次即唯一次". If first == only, there is no "非首次" state. The parenthetical tries to address this but ends up confusing. The measurement condition should simply state "每次调用" without the cold/warm distinction.

**Constraints & dependencies (26/30)**:
Dependencies well-identified: Go reflect (stdlib), yaml.v3 (existing), mode detection API (new), CLI subprocess dependency for skills. The bash template with fallback (lines 121-141) addresses the subprocess dependency. The `yaml:",inline"` tag constraint is now explicitly called out (line 119).

- **Deduction -2**: The skill -> CLI subprocess dependency (line 141) says "选择 CLI 子进程而非直接读取 config.yaml（via yq）的理由：CLI 封装了默认值应用逻辑". This is a valid reason, but the 50-100ms subprocess latency per eval check is a real cost. The comparison to "分钟级" pipeline time is misleading -- the subprocess overhead should be compared to the decision it enables (a boolean check that takes microseconds in-process), not the entire pipeline.
- **Deduction -2**: The mode detection API (line 120) depends on parsing `pwd` to extract feature slug. This assumes the user is running forge commands from within the feature directory. But forge can be invoked from the project root or other contexts. The pwd-based detection may not work for all invocation patterns.

### 5. Solution Creativity (32/100)

**Novelty over industry baseline (10/40)**:
The proposal explicitly states "eval 配置复用已有 ModeToggle 模式，无新概念" (line 88). The generic routing is standard Go practice -- Viper's `Unmarshal`, koanf's struct mapping, and dozens of other libraries use reflection for config traversal. The combination (generic routing + eval config in one proposal) is practical but not creative.

- **Deduction -5**: "Go struct IS the schema" (line 89) is presented as a novel insight but is standard practice in Go config libraries. Every Go struct-tag-based config system treats the struct as the schema. The "meta-programming" framing adds buzzwords without adding novelty.
- **Deduction -15**: The Innovation Highlights section (lines 86-91) lists 5 items, none of which are genuinely innovative. (1) Generic routing -- standard technique. (2) ModeToggle reuse -- explicitly conservative. (3) "Go struct IS the schema" -- standard practice. (4) Mode detection via CLI -- straightforward. (5) "Platform thinking" -- a framing choice, not a technical innovation.
- **Deduction -10**: The reflection-based CLI routing is described as "这种 meta-programming 模式在 Go 配置库中不常见" (line 91) but the proposal itself admits Viper and koanf use similar mechanisms. The claimed novelty is "反射与 CLI get/set 路由结合" -- but this combination is the minimum viable approach, not a creative leap.

**Cross-domain inspiration (7/35)**:
No cross-domain inspiration. ModeToggle from same codebase. Reflection from Go stdlib. Generic routing from standard config library patterns.

- **Deduction -28**: Zero borrowing from other domains. Database query planning, network routing (longest prefix match), compiler symbol tables (scoped name resolution), IDE configuration (layered settings), filesystem path resolution -- none referenced.

**Simplicity of insight (15/25)**:
The core insight -- "use reflection to replace hardcoded dispatchers" -- is the obvious answer to "how do we make config routing extensible?" The observation that eval config and routing generalization should be bundled is practical.

- **Deduction -10**: The insight is straightforward. A competent Go developer would arrive at the same solution within minutes of seeing the hardcoded switch/case dispatchers.

### 6. Feasibility (82/100)

**Technical feasibility (34/40)**:
Go reflect can do everything described. Function signatures are correct. The inline tag handling (line 170) and SurfacesMap fallback (line 170, 242) are now addressed. The `parseAutoRaw` generalization with flat-path keys (line 203) is a well-defined approach.

- **Deduction -3**: The `CoverageConfig.ByType` inline map has values of type `CoverageStrategy` (a struct with `Type` and `Percentage`). The reflection router must: (1) detect the `,inline` tag, (2) skip the field name level, (3) look up the key in the map, (4) recursively format the struct value. This 4-step chain is specified across multiple locations but not as a cohesive design. The "类型感知格式化器" concept (line 243) is introduced in the risk table, not the solution design.
- **Deduction -3**: The `WorktreeConfig.CopyFiles` ([]string) and `SurfacesMap` (custom YAML marshaling) are edge cases that each require special handling. The proposal acknowledges them but doesn't show how the reflection router's type switch handles these non-standard types uniformly.

**Resource & timeline feasibility (25/30)**:
4-5 hours revised estimate: 2.5h routing rewrite + 1h eval config + 1-1.5h testing + 1h edge cases. More realistic than previous 3-4h.

- **Deduction -3**: 2.5h for the routing rewrite covers: `getStructValueByPath`, `setStructValueByPath`, YAML tag matching, inline tag handling, SurfacesMap fallback, `parseAutoRaw` generalization, `applyDefaults` modification, deleting old dispatchers, and updating all call sites. Approximately 15 functions in 2.5h = 10 min/function. Possible for an expert who wrote the original code, but no margin for unexpected edge cases.
- **Deduction -2**: The "额外 1 小时用于处理边缘情况" (line 180) implicitly acknowledges the base estimate is tight. The edge cases listed (inline tag, SurfacesMap, CopyFiles) are non-trivial and could each take 30+ minutes.

**Dependency readiness (23/30)**:
Go reflect (stdlib) and yaml.v3 (existing) are ready. Mode detection API needs new implementation. `yaml:",inline"` tag handling adds complexity.

- **Deduction -4**: The mode detection API (line 120) requires: parsing pwd to extract feature slug, checking for proposal.md existence, registering a new CLI command. This is a mini-feature with its own edge cases (symlinks, nested feature directories, Windows path separators). The proposal treats it as a simple extraction from `feature_complete.go` but the CLI command registration and pwd-based slug detection are new work.
- **Deduction -3**: The `yaml:",inline"` tag handling requires the reflection router to inspect struct tags and modify traversal behavior. This interacts with yaml.v3's serialization semantics and must be tested against actual YAML output. The implementation complexity is non-trivial.

### 7. Scope Definition (76/80)

**In-scope items are concrete (28/30)**:
In-scope items name specific functions and code changes. `parseAutoRaw` generalization includes flat-path key format. Skill modifications now specify exact insertion points (line 207: "brainstorm 步骤 6-7 之间、write-prd 步骤 6-7 之间、ui-design 步骤 6-7 之间、tech-design 步骤 5-6 之间"). Mode detection API is fully specified.

- **Deduction -2**: The "4 个 skill 增加 config check" scope item (line 207) specifies step insertion points but doesn't clarify that the config check and AskUserQuestion coexist (AskUserQuestion is the fallback path). The scope should state "替换 AskUserQuestion 步骤为 config check + fallback 到 AskUserQuestion".

**Out-of-scope explicitly listed (24/25)**:
11 out-of-scope items. Comprehensive list including: eval skill behavior, new eval types, config validation, concurrent modification, rollback automation, forge guide updates, debuggability documentation.

- **Deduction -1**: "PR-1 rollback 方案的自动化测试" is out-of-scope (line 222) but Delivery Strategy says rollback is "revert commit". If the rewrite deletes old dispatchers, a revert must restore them. The revert procedure is stated but not verified.

**Scope is bounded (24/25)**:
Delivery Strategy with 2 PRs provides clear separation. PR-1 is infrastructure only. PR-2 is feature only. Scope is well-bounded.

- **Deduction -1**: The mode detection API (lines 208, 120) remains a significant embedded feature. Feature directory detection, slug extraction, proposal.md checking, and CLI command registration are each non-trivial. The proposal treats this as a single scope item but it could be a separate PR.

### 8. Risk Assessment (80/90)

**Risks identified (27/30)**:
8 risks identified. Risk table now covers: Part 1 regression, inline tag handling, SurfacesMap custom types, map-of-struct formatting, YAML comment preservation, parseAutoRaw generalization, skill config check drift, CLI subprocess dependency. Comprehensive.

- **Deduction -2**: Missing risk: `applyDefaults` interaction with flat-path raw format for existing fields. The `parseAutoRaw` generalization changes raw keys from `"test"` to flat-path format. If existing field keys change, `applyDefaults` calls may break silently.
- **Deduction -1**: Missing risk: bash template divergence over time. Markdown skills are edited by AI agents. The "EXTREMELY-IMPORTANT" annotation and "共享 snippet" mitigation (line 246) are aspirational -- the shared snippet is described as something that "can" be done, not something committed to.

**Likelihood + impact rated (26/30)**:
Ratings are well-calibrated. Part 1 regression at M/H. Inline tag handling at M/H. SurfacesMap at L/M (reasonable with fallback).

- **Deduction -2**: "4 个 skill 的 config check 逻辑漂移" at M/M -- likelihood should be H. Markdown skills are edited by AI agents with no compile-time enforcement. The bash template is 15+ lines with 3 branches. Over time, agents will modify or simplify the template. M/M is optimistic.
- **Deduction -2**: "YAML 注释和格式在 set 操作后的保真度" at L/L -- impact should be M. When users first use the new routing to set a value, all comments disappear. This generates support requests. "当前所有 set 路径已有此特性" is correct but doesn't reduce user-facing impact.

**Mitigations are actionable (27/30)**:
Mitigations are specific. PR splitting, inline tag detection with test, SurfacesMap fallback, parseAutoRaw with test cases, CLI fallback with `go install` in CI, bash template with annotation.

- **Deduction -2**: "为 map[string]struct 类型注册类型感知格式化器" (line 243) introduces a formatter registry not part of the core design. How are formatters registered? Global map? Interface? Function table? The mitigation adds design complexity without specifying the approach.
- **Deduction -1**: "将模板固化为 skill 的共享 snippet" (line 246) is aspirational ("可" = "can be"). Either commit to the shared snippet or explicitly defer it.

### 9. Success Criteria (75/80)

**Criteria are measurable and testable (27/30)**:
PR-1 has 14 SCs covering get (multiple depths and types), set (happy, non-leaf rejection, ModeToggle rejection), error paths, regression, and parseAutoRaw. PR-2 has 13 SCs covering mode detection (3 cases), skill behavior (4 skills x 2 modes), defaults, and a proxy verification method. Most are directly testable.

- **Deduction -2**: "forge config get auto.eval" (line 255) says "字段按 Go struct 定义顺序排列" -- testable in principle but fragile. If struct field order changes in a future refactor, this SC breaks silently. A more robust SC would specify a deterministic ordering rule independent of implementation.
- **Deduction -1**: "Skill 行为 SC 代理验证" (line 282) says "手动验证至少 1 个 skill" without specifying which skill or what constitutes "验证". A more specific SC would name the skill and the expected observable behavior.

**Coverage is complete (24/25)**:
SCs cover: get (multiple depths/types), set (happy and error), error paths, regression (worktree, coverage, parseAutoRaw), mode detection (3 cases), skill behavior (8 scenarios), defaults, intermediate node format, inline tag handling.

- **Deduction -1**: Missing SC for `forge config set auto.eval.prd.full true` when the `auto.eval` section doesn't exist in config.yaml. Does set auto-create intermediate structure? SC line 259 says "正确写入嵌套 config" but doesn't specify auto-creation of intermediate path.

**SC internal consistency (24/25)**:
SCs split into PR-1 and PR-2. Internal logic consistent. ModeToggle get/set asymmetry explicitly addressed.

- **Deduction -1**: SC line 256 specifies `forge config get auto` with "嵌套 struct（eval）缩进 2 空格后递归展开" but the "缩进 2 空格" phrasing is ambiguous -- is the indent applied to the entire nested output, or only to the field labels?

### 10. Logical Consistency (80/90)

**Solution addresses stated problem (31/35)**:
The solution directly addresses both problems. Generic routing solves the extensibility bottleneck. Eval configuration solves the interaction friction. Mode detection enables quick/full differentiation.

- **Deduction -3**: The problem framing ("消除交互摩擦") doesn't match the solution ("make interaction configurable"). Default config requires manual confirmation for prd and techDesign. The urgency argument (60-80s wasted per pipeline run) only fully materializes if users change defaults. The proposal makes the safe choice (conservative defaults) but the problem statement overpromises the benefit.
- **Deduction -1**: The Assumptions Challenged table (line 194) marks "proposal 默认应询问用户" as "Challenged" via "5 Whys" but the finding is a single design opinion ("自动评估可尽早发现问题"), not a logical refutation or root cause analysis. No evidence that auto-evaluating proposals is net-positive.

**Scope <-> Solution <-> SC aligned (26/30)**:
Strong alignment. Scope includes mode detection API -> Solution describes it -> SC has 3 criteria. Scope includes parseAutoRaw generalization -> Solution describes flat-path keys -> SC has eval-specific and regression criteria. PRs mapped to SCs.

- **Deduction -2**: Scope includes "applyDefaults 相应改为 applyModeDefault(&a.Eval.Proposal, a.raw, 'eval.proposal', d.Eval.Proposal) 的 flat-path 调用方式" -- does this mean existing `applyModeDefault(&a.Test, a.raw, "test", d.Test)` calls are also changed? If so, the SC should verify existing fields still work. If not, the code has two calling conventions.
- **Deduction -2**: Delivery Strategy (line 231) says "PR-1 合并前需通过 acceptance test" but doesn't define what passing means. All PR-1 SCs? Full regression suite? The acceptance criteria for PR-1 should be explicit.

**Requirements <-> Solution coherent (23/25)**:
Requirements and solution well-aligned. 11 key scenarios map to solution components. Default values match backward-compatibility requirements. NFRs constrain the implementation.

- **Deduction -2**: The "配置热生效" NFR (line 113) is trivially satisfied and doesn't constrain the solution. It is a non-requirement adding noise. More useful NFRs (e.g., "config get output is parseable by scripts", "error messages are grep-friendly") would add real value.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 92 | 110 |
| Solution Clarity | 108 | 120 |
| Industry Benchmarking | 96 | 120 |
| Requirements Completeness | 97 | 110 |
| Solution Creativity | 32 | 100 |
| Feasibility | 82 | 100 |
| Scope Definition | 76 | 80 |
| Risk Assessment | 80 | 90 |
| Success Criteria | 75 | 80 |
| Logical Consistency | 80 | 90 |
| **Total** | **818** | **1000** |

## Deductions Applied

- Vague language without quantification (2 instances): -40 pts
  1. "未来知识管理和验证自动化预计 3-4 字段" -- uncommitted projection treated as economic justification (Industry Benchmarking)
  2. "反射路由是平台思维而非功能思维" -- standard technique presented as strategic innovation (Solution Creativity)
- Self-congratulatory framing (1 instance): -20 pts
  1. Innovation Highlights lists 5 items, none genuinely innovative. "Go struct IS the schema" is standard practice presented as meta-programming breakthrough (Solution Creativity)

## Attack List

1. **[Problem Definition]** Urgency for routing rewrite is assumed, not demonstrated -- the proposal says the routing ceiling is reached NOW ("当前需求已到达现有路由的天花板") but line 158 shows the eval problem is solvable in 20 minutes with flat fields. The routing rewrite is an architectural preference, not a necessity. No cost-of-delay calculation: deferring by 2-4 weeks has zero consequence.

2. **[Problem Definition]** No user-reported evidence for eval friction -- the problem is entirely developer-inferred. No user feedback, support tickets, or usage data. The 60-80s estimate is plausible but unverifiable.

3. **[Solution Clarity]** Intermediate node output couples to implementation detail -- "字段按 Go struct 定义顺序排列" (line 255) means reordering Go struct fields changes CLI output. This implicit coupling between implementation detail and user-facing contract will cause silent breakage.

4. **[Industry Benchmarking]** Breakeven analysis misframes the economics -- 4 of the ~5 breakeven fields are needed NOW for eval. The real comparison is: 60-80 lines incremental NOW vs 200 lines rewrite NOW. That is a 2.5-3.3x premium, not a breakeven. The "3-4 future fields" projection remains uncommitted.

5. **[Industry Benchmarking]** "Platform thinking" framing is self-serving -- standard Go reflection applied to a 6-field config system. Every infrastructure rewrite can be called "platform thinking." The system doesn't yet have the scale that warrants platform-grade routing.

6. **[Requirements Completeness]** Missing scenario: set on non-existent intermediate path -- `forge config set auto.eval.prd.full true` when `auto.eval` doesn't exist in config.yaml. The feasibility section mentions nil pointer auto-initialization but this behavior has no scenario or SC.

7. **[Requirements Completeness]** Performance NFR measurement conditions are contradictory -- line 114 says "CLI 正常运行态（非首次冷启动）" but admits "CLI 每次调用是新进程，首次即唯一次". If first == only, there is no "非首次" state. The measurement condition is undefined.

8. **[Solution Creativity]** Zero genuine novelty -- Innovation Highlights lists 5 items that are standard practice, obvious, or self-congratulatory. The proposal explicitly acknowledges "eval 配置复用已有 ModeToggle 模式，无新概念". This dimension scores low because the proposal doesn't attempt creativity, which is an honest choice but costs points.

9. **[Feasibility]** Mode detection API is underdesigned for its complexity -- pwd-based slug extraction (line 120) needs to handle symlinks, Windows path separators, nested feature directories, and invocation from project root. The proposal treats it as a simple extraction but the CLI context determination is a new design problem.

10. **[Feasibility]** 2.5h routing rewrite is optimistic given edge case density -- 12+ functions plus 3 special cases plus parseAutoRaw generalization plus applyDefaults modification. The "额外 1h for edge cases" implicitly acknowledges the base estimate is tight.

11. **[Risk Assessment]** Skill config check drift likelihood is underrated -- M/M should be H/M. Markdown skills are edited by AI agents, the bash template is 15+ lines with 3 branches, no compile-time enforcement. The "EXTREMELY-IMPORTANT" annotation is equivalent to a TODO comment.

12. **[Risk Assessment]** Missing risk: applyDefaults flat-path migration for existing fields -- if `parseAutoRaw` generalization changes raw keys for existing fields (even format, not content), `applyDefaults` may silently stop applying defaults. Not in risk table.

13. **[Success Criteria]** PR-1 acceptance criteria undefined -- Delivery Strategy says "PR-1 合并前需通过 acceptance test" (line 231) but doesn't define passing criteria. All PR-1 SCs? Full regression? Manual smoke test?

14. **[Logical Consistency]** Problem framing overpromises -- problem says "消除交互摩擦" but solution makes interaction configurable with conservative defaults. Two of four evals still require manual confirmation by default. The urgency (60-80s saved) only fully materializes if users change defaults.

15. **[Logical Consistency]** Assumptions Challenged entries are opinions presented as findings -- "proposal 默认应询问用户" is marked "Challenged" via "5 Whys" but the finding is a design opinion, not a logical refutation. No evidence that auto-evaluating proposals is net-positive.
