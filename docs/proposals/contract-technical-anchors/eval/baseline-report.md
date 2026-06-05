# Eval-Proposal Baseline Report

**Final Score**: 678/1000 (target: 900)
**Iterations Used**: 0/3 (baseline)

---

## Reasoning Audit

### Problem → Solution Trace
Proposal states: Contract specs lack technical anchors, causing gen-test-scripts to rely on LLM inference, which produces mismatches invisible to all three test layers.

Solution: Build information chain from design docs → Contract anchors → test code with cross-validation.

**Verdict**: The solution directly addresses the stated problem. The chain is logically sound. However, the solution introduces a new authority-source assumption (design docs are always correct) that partially reintroduces the problem it claims to solve — replacing "LLM guesses" with "potentially stale handbook values" is the same class of issue.

### Solution → Evidence Trace
Evidence: One concrete incident (pm-work-tracker POST vs PUT mismatch).

**Verdict**: Single data point. Strong as a proof-of-existence, insufficient as proof-of-scale. No data on how frequently this class of bug occurs across projects.

### Evidence → Success Criteria Trace
SC1-SC2 test anchor presence. SC3 tests the specific incident case. SC4 tests backward compatibility. SC5 tests stale design doc detection.

**Verdict**: SC covers the demonstrated scenario and backward compat. Missing: SC for multi-surface scenarios (CLI/Web/Mobile), SC for handbook staleness, SC for auto-fix accuracy.

### Self-Contradiction Check
The proposal claims "自动修复" (auto-fix) but also lists as a risk that "设计文档本身有误". The mitigation ("保存原始值到注释") contradicts the spirit of automation — it's a manual recovery mechanism for an automated process. The proposal claims "无额外网络或 IO 开销" for cross-validation but cross-validation must read and parse handbook files, which is IO.

---

## Dimension Breakdown

### 1. Problem Definition: 82/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | Core problem is unambiguous: Contract lacks technical anchors, LLM infers wrong details, tests miss it. Two readers would interpret this the same way. Minor deduction: "三层测试均无法捕获" could be read as "all three layers failed" or "the three-layer system as a whole failed to catch it" — the actual incident shows the handbook caught it, just not the Contract. |
| Evidence provided | 22/40 | One concrete incident cited (pm-work-tracker). Strong as existence proof. However: no data on frequency ("此类问题会在每个有 API 或 CLI surface 的项目中重复出现" is a claim without supporting data), no user feedback, no metrics on how many Contracts currently lack anchors. Quote: "此类问题会在每个有 API 或 CLI surface 的项目中重复出现" — assertion without quantitative backing. |
| Urgency justified | 25/30 | "每个有 API 或 CLI surface 的项目中重复出现" gives urgency. "三层测试全部漏掉，直到生产环境返回 422" shows consequence of delay. Deduction: cost of delay not quantified. |

### 2. Solution Clarity: 92/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 35/40 | Four-step solution clearly laid out. A reader can explain back what will be built. Deduction: "交叉验证" logic underspecified — what counts as "match"? Exact string match? Normalized path comparison? |
| User-facing behavior described | 35/45 | The observable behavior of the pipeline is described (auto-fix Contract, flag code bugs). Deduction: What does the user SEE when auto-fix triggers? A log message? A diff? A prompt? The UX of the fix process is undefined. Quote: "不匹配时以设计文档为准自动修复 Contract" — silent auto-modification of user files without describing user experience. |
| Technical direction clear | 22/35 | Direction is clear (add fields to frontmatter, read from handbook, compare with Fact Table). Deduction: The new handbook formats (cli-handbook, page-map, screen-map) are mentioned but their structure is completely undefined. "复用 api-handbook 的成熟格式模式" is not a technical specification. |

### 3. Industry Benchmarking: 62/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 20/40 | Pact and OpenAPI/Swagger are mentioned in passing: "行业中常见的做法是 Contract Testing（Pact）或 OpenAPI spec 驱动的测试生成". This is a one-sentence mention with no analysis of how they solve the anchor problem, what their limitations are, or how Forge's approach compares in detail. No product names beyond these two, no open-source projects, no published patterns cited. |
| At least 3 meaningful alternatives | 17/30 | Four alternatives listed including "do nothing". However: (1) "仅增强 Fact Table" is a straw-man — it's described as "治标不治本" (treating symptoms not root cause) without explaining why it can't address design-implementation consistency; (2) "OpenAPI spec 驱动" is dismissed as "架构不匹配" without elaboration. Quote: "与 Forge 的语义 Contract 模型不兼容" — why incompatible? No explanation. |
| Honest trade-off comparison | 10/25 | Trade-offs are biased toward the selected approach. The selected approach's cons are "需要扩展 tech-design 和 gen-test-scripts" which understates the scope (three new handbook formats to design, auto-fix safety, cross-surface testing). Quote: "需要扩展 tech-design 和 gen-test-scripts" — omits the non-trivial work of designing cli-handbook/page-map/screen-map formats. |
| Chosen approach justified against benchmarks | 15/25 | "最小改动覆盖最大范围" is the justification. This is a conclusion, not an argument. No analysis of why adding anchor fields to Contract is less effort than enhancing Fact Table reconnaissance. |

### 4. Requirements Completeness: 68/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 28/40 | Four key scenarios identified including the critical "设计文档与代码不一致" case and the "Handbook 不存在" backward-compat case. Deductions: (1) Missing: Contract manually edited after anchor fill — what happens? (2) Missing: Multiple Contracts sharing the same endpoint — does cross-validation handle N:1 mappings? (3) Missing: Partial handbook coverage (e.g., project has api-handbook but not cli-handbook). |
| Non-functional requirements | 22/40 | Backward compatibility mentioned. Performance claim: "无额外网络或 IO 开销" — this is **factually incorrect**. Cross-validation requires reading and parsing handbook files (file IO) plus comparison logic. Security not addressed (auto-modifying files has security implications). Accessibility irrelevant. |
| Constraints & dependencies | 18/30 | Three dependencies listed. Missing: (1) Fact Table code reconnaissance accuracy as a constraint — cross-validation's reliability is bounded by reconnaissance coverage. (2) The assumption that tech-design will be re-run when code changes (no staleness detection mechanism). |

### 5. Solution Creativity: 65/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 28/40 | The "design doc as authority source" with auto-fix and code-bug flagging is a meaningful extension over standard contract testing. The insight that design-implementation mismatches are code bugs, not test bugs, is genuinely useful. Deduction: The core mechanism (populate fields from a spec, compare with code) is standard practice in OpenAPI-driven tooling. |
| Cross-domain inspiration | 15/35 | No evidence of borrowing from other domains. The proposal stays within the testing/contract domain. No references to how other ecosystems (e.g., type systems, compiler error recovery, database schema migration strategies) handle similar authority-source conflicts. |
| Simplicity of insight | 22/25 | "Contract 是设计意图的规格说明，技术锚点应该来自设计阶段" — this is an elegant and clear insight. The "why didn't I think of that" quality is present. |

### 6. Feasibility: 72/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 32/40 | Three existing skills need modification — plausible. api-handbook pattern exists as reference. Deduction: "完全可行" is overstated. The proposal requires designing three new handbook formats (cli-handbook, page-map, screen-map) whose complexity is underspecified. CLI command identification alone has significant edge cases (nested subcommands, aliases, parameter variants). |
| Resource & timeline feasibility | 20/30 | "单项 enhancement，改动点明确，无外部依赖". This is too optimistic. The scope touches 4+ skills (tech-design, gen-contracts, gen-test-scripts, eval-contract) and requires designing 3 new document formats. No timeline estimate provided. |
| Dependency readiness | 20/30 | api-handbook is ready. "cli-handbook / page-map / screen-map 是新增文档类型，无前置依赖" — calling them "no dependencies" conflates "no blocking dependency" with "ready to build". These formats need to be designed, validated, and stabilized before gen-contracts can reliably fill anchors from them. |

### 7. Scope Definition: 58/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 22/30 | Five concrete items listed. Most are deliverables. Deduction: "eval-contract 评分规则增加技术锚点完整性检查" — what completeness means is undefined (all Contracts must have anchors? only when handbook exists?). |
| Out-of-scope explicitly listed | 18/25 | Four items listed. Good. Missing: "Contract manual editing consistency" and "handbook staleness detection" are not in-scope or out-of-scope — they're unaddressed gaps. |
| Scope is bounded | 18/25 | In-scope covers 4 surfaces (API + CLI/TUI + Web + Mobile) in one batch, but Risk Table says "可分批实现". These contradict. Quote from Scope: "tech-design 增加 CLI/TUI cli-handbook、Web page-map、Mobile screen-map 自动生成" vs Risk Table: "每种 surface 的锚点字段独立、互不影响，可分批实现". Which is it? |

### 8. Risk Assessment: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 20/30 | Three risks listed. Missing critical risks: (1) Fact Table reconnaissance inaccuracy causing false positives in cross-validation; (2) handbook staleness (design doc updated without handbook regeneration); (3) auto-fix creating cascading errors in downstream test generation. |
| Likelihood + impact rated | 15/30 | All three risks use L/M/H ratings. However: "自动修复覆盖了正确的 Contract" is rated L likelihood — given the freeform reviewer's analysis that design docs are often stale in fast-moving projects, this should be at least M. The rating appears optimistic rather than honest. |
| Mitigations are actionable | 20/30 | "复用 api-handbook 的成熟格式模式" — actionable. "修复前保存原始值到 Contract 的注释中，可回溯" — partially actionable (saves original, but no mechanism to detect corruption or auto-revert). "可分批实现" — not a mitigation, it's scope reduction. |

### 9. Success Criteria: 52/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 22/30 | SC1/SC2 are measurable (100% coverage, conditional on handbook existence). SC3 is testable (specific POST vs PUT scenario). SC4 is testable (pipeline runs without error). SC5 is vague — "明确的代码 bug 标记报告" — what makes a report "明确的"? What format? How is it verified? |
| Coverage is complete | 15/25 | Gaps: (1) No SC for CLI/TUI/Web/Mobile surfaces — only API is tested in SC3. (2) No SC for auto-fix accuracy (what % of auto-fixes are correct?). (3) No SC for handbook staleness detection. (4) No SC for the eval-contract completeness check that's in scope. |
| SC internal consistency | 15/25 | SC1/SC2 depend on handbook existence (conditional). SC4 ensures pipeline works without handbook. These are consistent. However: SC3 assumes api-handbook defines PUT, but SC1 only requires `endpoint` field presence, not correctness. A Contract with `endpoint: "POST /move"` passes SC1 but fails SC3. The gap between "anchor exists" and "anchor is correct" is unaddressed in SC. The `consistency_check_result` block at the top shows `status: pass, pairs_checked: 15, conflicts_found: 0` — this appears to be a pre-computed result embedded in the proposal, which is unusual and raises questions about its rigor. |

### 10. Logical Consistency: 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 30/35 | Yes — anchor chain directly addresses the LLM inference problem. Deduction: The solution replaces one unreliable source (LLM inference) with another potentially unreliable source (stale handbook values) for the cross-validation scenario. The class of problem (wrong authority source) persists, just shifted from LLM to handbook. |
| Scope ↔ Solution ↔ Success Criteria aligned | 18/30 | In-scope includes "eval-contract 评分规则增加技术锚点完整性检查" but no SC validates this. In-scope includes all 4 surfaces but SC only tests API surface. Solution says "自动修复" but no SC measures auto-fix correctness. Quote: Scope says "Web page-map、Mobile screen-map 自动生成" but no SC verifies these are generated correctly or used. |
| Requirements ↔ Solution coherent | 24/25 | Requirements map cleanly to solution. The four key scenarios each map to a solution component. Minor gap: "Handbook 不存在" scenario maps to "降级为 Fact Table 推断" which is the current behavior — this is backward compat, not a solution improvement for that scenario. |

---

## Blindspot Hunt

**[blindspot-1] Authority source trust model is binary and naive.** The proposal treats the design document as ground truth and code as the thing that needs to match. In practice, both can be wrong. Production systems use confidence scoring or multi-source reconciliation (e.g., Git blame + test results + runtime data) to determine which source to trust. The proposal's binary model will produce wrong fixes in the (not rare) case where the code is correct and the design doc is stale.

**[blindspot-2] No discussion of incremental/regression behavior.** When gen-test-scripts auto-fixes a Contract, does it re-run the test generation? If the fix changes the endpoint from POST to PUT, all previously generated tests referencing POST are now invalid. The proposal doesn't describe what happens after the fix — does it trigger regeneration? Does it invalidate existing test files?

**[blindspot-3] Fact Table reconnaissance limitations unaddressed.** The proposal assumes Fact Table code reconnaissance will accurately detect route registrations. But static analysis of route registration has known blind spots: dynamically registered routes, plugin-registered routes, framework-specific registration patterns (decorators, annotations, convention-based routing). When the Fact Table is wrong, the cross-validation produces false positives/negatives. This is a fundamental limitation that should be acknowledged and bounded.

**[blindspot-4] No success metric for the handbook formats themselves.** Three new document formats are proposed (cli-handbook, page-map, screen-map). There is no criterion for whether these formats are well-designed, complete, or usable. A poorly designed handbook format could make the entire anchor system unreliable without triggering any of the existing SCs.

**[blindspot-5] "100% 包含" success criteria creates a perverse incentive.** SC1 says "API surface 的 Contract 100% 包含 endpoint 字段（当 api-handbook 存在时）". This could be satisfied by filling in placeholder or incorrect values. The criterion tests presence, not correctness. A Contract with `endpoint: "TODO"` would pass SC1.

**[blindspot-6] Timeline and sequencing dependency not discussed.** The proposal touches tech-design, gen-contracts, gen-test-scripts, and eval-contract. These skills form a pipeline. The proposal doesn't discuss implementation sequencing — can they be changed independently, or must they be updated atomically? If gen-test-scripts gets cross-validation before gen-contracts fills anchors, every Contract will fail cross-validation.

---

## Score Progression

| Iteration | Score | Delta |
|-----------|-------|-------|
| 0 (baseline) | 678 | — |

### Outcome

Target NOT reached — baseline evaluation. 678/1000 vs target 900. Significant gaps in Industry Benchmarking (-58 from max), Requirements Completeness (-42), Risk Assessment (-35), and Success Criteria (-28) must be addressed.

## Priority Improvements (ranked by score impact)

1. **Industry Benchmarking (+58 potential)**: Research and cite actual tools/patterns (Spring Cloud Contract, Dredd, Schemathesis, API Sprout). Analyze their anchor/spec-driven approaches in detail. Replace straw-man alternatives with genuine different approaches.
2. **Requirements Completeness (+42 potential)**: Add missing edge cases (partial handbook coverage, manual Contract edits, N:1 endpoint mappings). Fix the "no extra IO" claim. Acknowledge Fact Table accuracy as a constraint.
3. **Risk Assessment (+35 potential)**: Add Fact Table accuracy risk, handbook staleness risk, cascading auto-fix errors. Re-rate auto-fix-wrong-Contract likelihood honestly. Make mitigations specific and actionable.
4. **Success Criteria (+28 potential)**: Add SCs for non-API surfaces, auto-fix accuracy rate, handbook format quality, and eval-contract completeness check. Replace "anchor exists" with "anchor is correct" criterion.
5. **Scope Definition (+22 potential)**: Resolve batch vs full-surface contradiction. Clarify eval-contract completeness semantics.
6. **Feasibility (+28 potential)**: Provide timeline estimate. Acknowledge handbook format design complexity. Separate "no blocking dependency" from "ready to build."
