# Eval Report: Iteration 3 (FINAL)

**Iteration**: 3
**Date**: 2026-05-26
**Evaluator**: CTO Adversary

---

## Iteration 2 Issue Resolution Check

| # | Iteration 2 Issue | Status | Evidence |
|---|-------------------|--------|----------|
| 1 | No user-facing impact evidence | **Resolved** (carried from iter 1) | Line 22: eval-design scoring instance; line 67: before/after user-facing experience |
| 2 | No NFRs | **Resolved** (carried from iter 1) | Lines 96-103: six NFRs including backward compat, migration perf, discoverability, execution perf, CI integration, transition tracking |
| 3 | No edge cases | **Resolved** (carried from iter 1) | Lines 105-109: three edge cases |
| 4 | Blast radius unquantified | **Resolved** (carried from iter 1) | Lines 152-162: 27 files enumerated |
| 5 | Resource estimate vague | **Resolved** (carried from iter 1) | Lines 166-171: 9 tasks quantified |
| 6 | Task type parser feasibility unverified | **Resolved** (carried from iter 1) | Lines 71, 150: explicit parser verification |
| 7 | Trade-off analysis shallow | **Partially Resolved** | Lines 129-136: four trade-offs remain with one-sentence mitigations each. No cost quantification added. |
| 8 | SC5 threshold arbitrary | **Resolved** (carried from iter 1) | Line 222: now references specific skill rules files |
| 9 | No SC for backward compatibility | **Resolved** (carried from iter 1) | Line 225: SC7 (alias still executable) |
| 10 | No SC for quality gate output | **Resolved** (carried from iter 1) | Line 226: SC8 (suite names + quality gate classification) |
| 11 | `forge surfaces` output key not verified | **Unresolved** | Line 92: "forge surfaces CLI 已提供 surface 检测能力" -- still an assertion without evidence. No code reference, no output format example. |
| 12 | "Assumptions Challenged" uses rhetorical loading | **Unresolved** | Line 177: "Assumption Flip" still selectively defines "e2e" to mean "covers full stack with device automation" and then proves CLI/API don't match. The reasoning is circular: define e2e narrowly, then show CLI/API don't qualify. |
| 13 | Technical direction underspecified for justfile alias | **Resolved** | Lines 73: detailed implementation -- `alias old_name := new_name`, DEPRECATED comment, version-based removal, `init-justfile` template mechanism |
| 14 | Missing risks (classification extensibility, partial adoption, terminology regression) | **Resolved** | Lines 212-214: three new risks added -- classification binary extensibility (L/H), terminology regression (M/M), partial adoption during transition (M/M) |
| 15 | SC coverage gaps (task type names in index.json, business rules doc, transition period) | **Resolved** | Lines 223-224: SC6 covers index.json task type names; SC7 covers business rules doc update |
| 16 | Trade-off costs unquantified | **Unresolved** | Trade-offs at lines 133-136 still use one-sentence mitigations. No hour estimates, no cost projections. Blast radius data (27 files) exists in the proposal but is not referenced in trade-off analysis. |
| 17 | "端到端" Chinese naming creates SC3 gap | **Partially Resolved** | SC3 (line 220) now explicitly states: "搜索英文 'e2e' 只出现在 Web/Mobile surface 的端到端测试上下文中；搜索中文 '端到端' 仅出现在 Web/Mobile surface 的测试类型名称和定义中" -- both English and Chinese are now covered. However, this creates a new issue: "端到端" is allowed only in Web/Mobile test type names and definitions, but the classification standard declaration (line 51) uses "端到端测试" as a category label in the general sense. SC3 may conflict with the classification standard's usage. |
| 18 | SC4/SC5 overlap/ambiguity | **Unresolved** | SC4 (line 221: "所有涉及测试类型的 skill 文件使用 surface-specific 测试类型名称") and SC5 (line 222: "所有包含测试规则的 skill rules 文件...引用概念文档中的测试类型定义") still overlap heavily. |

**Summary of iteration 2 -> 3 changes**: The proposal added three new risks (classification extensibility, terminology regression, partial adoption), detailed justfile alias implementation, expanded NFRs from 3 to 6, added SC6-SC8 to cover previously identified gaps, and clarified SC3 for Chinese/English coverage. Key unresolved: trade-off quantification, `forge surfaces` verification, SC4/SC5 overlap, rhetorical loading in Assumptions.

---

## Phase 1: Reasoning Audit

### Problem -> Solution
The problem: "e2e" as a universal label misrepresents what different surface tests do, causing eval-design scoring errors and user confusion. The solution: Surface -> Test Type mapping with a two-tier classification (surface as primary key, functional/e2e as secondary attribute based on verification mechanism). **Chain holds.** The mapping directly addresses terminological imprecision.

**Reservation (carried from iter 2, still valid)**: The solution introduces a classification taxonomy (functional vs. e2e binary with a full classification standard declaration) that goes beyond what the problem demands. The problem is "the name 'e2e' is wrong for CLI/API." The solution could be "rename them to match their behavior" without building a two-axis taxonomic model. The taxonomy is defensible but represents scope inflation relative to the stated problem.

### Solution -> Evidence
Evidence (6 structural observations + eval-design scoring impact) supports surface differentiation existing in the codebase. **Chain holds for problem validation.** However, evidence does not validate the *new naming*. No data showing that "CLI Functional Test" is less ambiguous than "e2e" for the target audience. No A/B test, no user survey, no feedback from a pilot. The evidence proves the problem exists, not that this specific solution is the best one.

### Evidence -> Success Criteria
Eight SCs now cover documentation (SC1-2), terminology elimination (SC3), skill adoption (SC4-5), task type naming (SC6), business rules sync (SC7), backward compatibility (SC8), and test output (SC9). **Chain holds for implementation verification.** Gap: no SC validates the taxonomy's clarity or correctness for users -- all SCs are implementation completeness checks, none verify that users actually understand the new taxonomy.

### Self-contradiction Check

1. **Classification standard vs. semantic definitions -- tension persists but is now explicitly acknowledged.** Lines 51-53 now contain a "关键区分" paragraph that explicitly addresses the concern: "CLI 测试可以遍历完整技术栈...它们在技术栈覆盖上可能是'端到端'的，但验证机制是在协议边界上的单次调用观测." This is a deliberate definitional choice, not an oversight. **Conflict resolved by explicit acknowledgment.** However, the claim "验证机制是在协议边界上的单次调用观测" may not hold for complex CLI tools that involve multi-step interactive workflows. A CLI tool that prompts for input, processes data, and displays results is more than a "single call observation."

2. **Mobile "Best-effort" E2E testing -- still tensions.** Line 63: "Best-effort 模式，部分场景标记为 manual-only." If some scenarios are manual-only, the "端到端测试" label overclaims for those scenarios. The semantic definition does not specify what percentage of scenarios can be manual-only before the "E2E" label becomes misleading. Is 50% manual-only still "E2E"? 80%?

3. **"功能测试" semantic conflict -- mitigation improved but still leaky.** Line 210: mitigation says "每次使用都带 surface 前缀限定范围." Line 135: "每次使用时都带 surface 前缀（如 'CLI 功能测试'），限定范围." But justfile recipe names (line 85: `test-cli-functional`) use "functional" as a standalone word after the hyphen. A user reading `just test-cli-functional` sees "functional" without the definitional context. The surface prefix ("cli-") acts as a namespace, not as a semantic qualifier for "functional."

4. **SC3 scope vs. classification standard usage of "端到端".** SC3 (line 220) allows "端到端" "仅出现在 Web/Mobile surface 的测试类型名称和定义中." But the classification standard declaration (line 51) uses "端到端测试" as a general category label contrasting with "功能测试" -- this is not within a specific surface's test type name or definition. It is a meta-level categorical use. SC3's restriction may inadvertently make the classification standard declaration itself non-compliant, since it uses "端到端" in a general (not Web/Mobile-specific) context.

### SC Consistency Deep-Dive

Cluster by affected area:

**Cluster A: Documentation** (SC1, SC2)
- SC1: concept doc with 5 types. SC2: guide.md terminology. Both satisfiable independently. No conflict.

**Cluster B: Skill files** (SC4, SC5)
- SC4: "所有涉及测试类型的 skill 文件使用 surface-specific 测试类型名称." SC5: "所有包含测试规则的 skill rules 文件...引用概念文档中的测试类型定义." Overlap: if all skill files use the new names (SC4), and the names are defined in the concept doc, SC5 is trivially satisfied. SC5 adds value only if "引用概念文档" means explicit text citation (e.g., "see test-type-model.md for definition") rather than mere terminological consistency. The ambiguity in "引用" makes the relationship between SC4 and SC5 unclear -- are they independent requirements or is SC5 a subset of SC4?

**Cluster C: Task type / business rules** (SC6, SC7)
- SC6: index.json carries surface info in task type names. SC7: business rules doc (task-lifecycle.md) reserved types list includes new names. These are independent and consistent.

**Cluster D: Terminology elimination** (SC3)
- SC3: no generic "e2e" / "端到端" outside Web/Mobile context. Independent. Potential conflict with classification standard declaration (see contradiction #4 above).

**Cluster E: Backward compatibility** (SC8)
- SC8: old recipe names as aliases still executable. Independent. No conflict with other SCs. But: no SC verifies that aliases are *removed* at the end of the transition period. NFR #6 (line 103) says "alias 移除作为一个独立 task 记录在 Forge 版本规划中" but no SC enforces this task is actually created and tracked.

**Cluster F: Test output** (SC9)
- SC9: suite names + quality gate classification. "分类展示" is still ambiguous -- separate sections? Separate labels? A filter option? The testability depends on interpretation.

**InScope <-> SC gap check:**
- InScope "更新 task type 命名（携带 surface 信息）" -> SC6 covers this. **Resolved.**
- InScope "更新 business rules 文档中的测试相关术语" -> SC7 covers this. **Resolved.**
- InScope "更新 run-tests 的测试输出格式" -> SC9 partially covers (suite names + quality gate) but "测试输出格式" could include log format, error messages, progress indicators -- not all covered by SC9.
- InScope "更新 gen-test-scripts 输出的测试代码中的注释/标签" -> no SC explicitly verifies generated test code comments/tags. SC4 says skill files use surface-specific names but does not verify the *output* of those skills.
- No SC for transition period completion (alias removal). NFR #6 mentions it but no SC enforces it.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 90/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | Core problem unambiguous. Two concrete examples (eval-design scoring bias, API test-e2e mislabeling). The three-sentence opening paragraph (line 11) precisely states what each surface type does and why they are different. Deduction (-2): the opening paragraph mixes (a) the imprecision of "e2e" as a label with (b) the claim that they "不是同一类测试" (are not the same type of test). These are different claims -- the first is about naming, the second is about categorization. The problem could be stated more precisely as: "they share a label that implies they are the same type when they are not." |
| Evidence provided | 30/40 | Six structural evidence items (lines 15-20) + eval-design scoring impact (line 22) + API test-e2e mislabeling (line 22). The eval-design impact is the strongest evidence because it is a concrete tool behavior issue. Deduction (-10): still no user complaint, no support ticket, no bug filed, no user study. The "用户困惑实例" described in Assumptions Challenged (line 178: "当 CLI 项目的测试报告说 'e2e 覆盖率 100%' 时，用户会以为所有功能都端到端验证了") is presented as a hypothetical scenario, not a reported incident. The eval-design scoring issue is internal tooling behavior, not user-facing harm. The proposal has had 3 iterations to produce a single piece of external evidence and has not done so. |
| Urgency justified | 22/30 | Lines 26-28: present-tense eval-design impact + 3-month projection (~15 to ~25 files). The projection is now quantified with file counts. Deduction (-8): "~15 个" and "~25 个" remain estimates without backing. The "Forge 预计新增 2-3 种 surface 类型的项目支持" projection has no roadmap reference. The eval-design scoring issue's severity is not quantified -- how many projects affected? How much score deviation? Without this, the urgency remains argumentative rather than empirical. |

### 2. Solution Clarity: 100/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Mapping table (lines 36-42) with five surfaces, each with EN/CN names, verification dimensions, and execution models. Classification standard declaration (lines 46-53) provides the two-axis model with explicit definitional criteria. Deduction (-2): the "验证维度" column still mixes levels of abstraction. CLI lists "进程退出码 + stdout 文本 + stderr 文本" (all concrete observables) while Web lists "DOM 元素可见性 + 用户操作响应 + 页面 URL 变更 + 元素属性值" where "用户操作响应" is a high-level concept, not a directly observable output (it is inferred from DOM state changes). |
| User-facing behavior described | 44/45 | Line 67: three concrete before/after changes -- justfile output labels, suite names, CI dashboard display. The User-Facing Experience section is now specific and measurable. Deduction (-1): no description of what happens when a test *fails* -- does the failure message include the test type label? Does the error output format change? |
| Technical direction clear | 18/35 | Line 71: task type parser compatibility verified. Lines 73: justfile alias implementation now detailed -- `alias old_name := new_name`, DEPRECATED comment, version-based removal. This is a significant improvement from iteration 2. Deduction (-17): (1) the proposal disclaims being a code change ("纯文档 + 命名变更") but InScope includes "更新 justfile recipe 命名" and "更新 task type 命名" which are code changes in skill files and templates. This contradiction is not resolved. (2) No specification for how `test.gen-scripts.cli` type names propagate through the task lifecycle system -- what does the index.json entry look like before and after? (3) No specification for how run-tests skill selects the correct surface-specific test execution rule based on the new task type name. (4) The init-justfile alias mechanism is described but the template modification detail is insufficient -- which template file? What does the before/after template look like? |

### 3. Industry Benchmarking: 74/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 32/40 | Four references (Go build tags, Spring Boot @Tag, Playwright project configs, Postman/Newman) each with a structured analysis: pattern description, rejection reason, and "吸收" (what Forge takes from it). This is improved from iteration 2 -- each reference now has three components instead of one-to-two sentences. Deduction (-8): still no code examples for any reference. The Go build tags reference could show a concrete `//go:build e2e` example. The Playwright reference could show a `projects` config snippet. No documentation URLs, no article citations. |
| At least 3 meaningful alternatives | 22/30 | Four alternatives + "do nothing" in comparison table (lines 122-127). The "引入标准测试分层" alternative is no longer a straw man. Deduction (-8): no alternative proposes a *different* naming taxonomy. All alternatives either keep the status quo, use a vaguer label, use industry terms, or propose the selected approach. The space of "different taxonomic approaches to test classification" is not explored -- e.g., tagging-based (metadata labels without renaming), directory-based convention (tests/cli/, tests/api/), or execution-model-based naming without the functional/e2e binary. |
| Honest trade-off comparison | 8/25 | Four trade-offs (lines 131-136) with one-sentence mitigations each. Deduction (-17): trade-offs are *listed* but not *compared*. The section should answer "what do we give up?" not "what inconveniences will we face?" No quantification of any cost. Trade-off #2 says "所有引用 'e2e' 的文件需逐一更新" -- the proposal already counted 27 files. Why not say "27 files across 4 skill areas need updating"? Trade-off #3 says "每次使用时都带 surface 前缀" -- but the justfile recipe name `test-cli-functional` does NOT include the full definition, just the surface prefix and "functional." The mitigation claims more than it delivers. |
| Chosen approach justified against benchmarks | 12/25 | Line 118: "Newman 的命名原则——按执行方式命名而非按假设的覆盖范围命名——直接验证了本提案的核心主张." Line 127: "最小惊讶原则——名称匹配实际行为." Two benchmark justifications. Deduction (-13): no analysis of why Forge's approach is better than adopting Playwright's project config pattern for multi-surface test classification. No analysis of why "name by execution model" was chosen over "name by coverage scope" beyond the one-line Newman reference. The justification does not address why Forge cannot adopt a simpler approach (e.g., just rename "e2e" to "surface-test" without the functional/e2e taxonomy). |

### 4. Requirements Completeness: 88/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 36/40 | Five key scenarios + three edge cases. Multi-surface edge case well-handled (line 107). Transition period edge case addressed (line 109). Deduction (-4): missing edge case: what happens when a project's surface type changes after tests are already generated (e.g., CLI tool gains an API layer)? The generated test code with old type names would need updating, but no mechanism for this is described. |
| Non-functional requirements | 28/40 | Six NFRs (lines 96-103): backward compatibility, migration performance, discoverability, execution performance, CI integration, transition tracking. Significant improvement from iteration 2. Deduction (-12): (1) NFR #1 says "2 个 Forge 版本" transition period but does not define what triggers version counting -- is it major versions? Minor versions? Release cadence? (2) NFR #4 claims "每个 surface recipe 的执行路径与当前 test-e2e recipe 完全一致" but this is an assertion, not a requirement -- it describes what *will be* true, not what *must be ensured*. A proper NFR would specify a measurable constraint: "per-surface recipe execution time <= current test-e2e execution time + 0ms." (3) No security NFR -- terminology changes in generated test code could expose internal classification logic to end users who see test output. |
| Constraints & dependencies | 24/30 | Four constraints (lines 89-94). Line 92: "forge surfaces CLI 已提供 surface 检测能力" -- still unverified. Line 94: "任务类型名的变更需与 task-lifecycle business rule 中的保留类型列表协调" -- identifies coordination dependency. Deduction (-6): `forge surfaces` CLI output format still not verified. No version constraint (which Forge version will include this?). No dependency on init-justfile's template system being compatible with alias syntax. |

### 5. Solution Creativity: 54/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 20/40 | The two-tier classification (Surface + functional/e2e) is a minor conceptual contribution over industry baselines. The insight that CLI/API tests should not be called "e2e" because they use protocol-level verification rather than device-level automation is the novel claim. However, this is essentially a terminological argument, not a technical innovation. The proposal positions itself as naming existing practice, not creating new concepts. |
| Cross-domain inspiration | 14/35 | Limited cross-domain references. No mention of taxonomy design principles, ontology engineering, controlled vocabularies, or how other multi-product systems (monorepo tools like Nx/Turborepo, polyglot build systems like Bazel) handle type classification. The Newman reference (API tools naming by execution method) is the closest to cross-domain inspiration. The proposal looks predominantly at testing tools within the Forge ecosystem. |
| Simplicity of insight | 20/25 | Core insight: "the differentiation already exists in the codebase, just name it accurately." The two-axis model is clean. Deduction (-5): the classification standard declaration (lines 46-53) is a full paragraph of definitional hedging ("关键区分：CLI 测试可以遍历完整技术栈...但验证机制是在协议边界上的单次调用观测"). If the insight were truly simple, it would not require this much justification. The "functional vs. e2e" binary requires a complex paragraph to explain why it is not about technology stack coverage but about verification mechanism -- this complexity suggests the binary may not be as natural a distinction as claimed. |

### 6. Feasibility: 84/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | Lines 142-146: four existing infrastructure elements. Lines 71, 150: task type naming verified against parser behavior. Lines 73: alias mechanism specified. Deduction (-4): `forge surfaces` CLI output key format still not verified (line 92). The init-justfile template alias mechanism assumes `just` natively supports alias syntax -- this is stated but not demonstrated with a concrete example of a before/after template file. |
| Resource & timeline feasibility | 24/30 | Lines 166-171: 9 tasks (6 doc + 3 coding). Concrete breakdown by skill area. Deduction (-6): no time estimate per task. No team size or availability analysis. No sequencing information (which tasks depend on others?). The "概念参考文档" task is a prerequisite for all skill file updates (SC5 requires referencing the concept doc), but this dependency is not mentioned. |
| Dependency readiness | 24/30 | `forge surfaces` CLI mentioned (line 92) but not verified -- output format unknown. init-justfile alias support detailed (line 73) -- just native syntax confirmed. task-lifecycle parser verified (line 150). Deduction (-6): one of three dependencies (`forge surfaces`) is still assumed without evidence. |

### 7. Scope Definition: 74/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 27/30 | Ten in-scope items (lines 185-195), most specific deliverables. Line 188: "更新 guide.md（Terminology 部分），补充 Surface Type -> Test Type 的简要说明" -- specific location and content. Line 195: "更新 run-tests 的测试输出格式，使 suite 名称和标签使用 surface-specific 测试类型名称" -- specific. Deduction (-3): line 193 "更新 business rules 文档中的测试相关术语" -- which business rules documents? Blast radius section mentions only "task-lifecycle.md = 1 文件" but the in-scope item says "文档" (plural). |
| Out-of-scope explicitly listed | 22/25 | Four items out of scope (lines 199-202). Line 202: "测试目录结构的重新组织（保持 tests/<journey>/ 或按 surface-key 分目录可作为后续优化）" -- the parenthetical "可作为后续优化" hints at scope expansion. Either it is firmly out of scope or it is a candidate. The parenthetical weakens the boundary. |
| Scope is bounded | 25/25 | Bounded by 5 surface types, 27 files, 9 tasks. Blast radius quantified. Resource estimate concrete. |

### 8. Risk Assessment: 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 26/30 | Seven risks (lines 206-214). Three new risks added since iteration 2: classification binary extensibility (line 212), terminology regression (line 213), partial adoption during transition (line 214). The three previously identified risks (document invalidation, CI breakage, semantic conflict) remain. Deduction (-4): missing risk: generated test code regression -- when gen-test-scripts updates its output format (comments, tags, file naming conventions), existing generated test files in user projects will not be automatically updated. Users would need to re-run gen-test-scripts, which may overwrite manual edits to test code. This is a different risk from "过渡期混用" (line 214) which addresses terminology, not code regeneration. |
| Likelihood + impact rated | 22/30 | CI risk: M/H. Document invalidation: M/M. Semantic conflict: M/M. Learning cost: L/L. Binary extensibility: L/H. Terminology regression: M/M. Partial adoption: M/M. Ratings are plausible. Deduction (-8): no backing data for any rating. The blast radius section counts 27 files but the risk section does not reference this count. How many CI pipelines reference `just test-e2e`? How many external documents/tutorials exist? The partial adoption risk is rated M/M but the mitigation (line 214: "所有 skill 文件的术语更新在单次 PR 中完成") reduces likelihood to effectively zero -- so why is it rated M? If the mitigation works, the likelihood should be L. |
| Mitigations are actionable | 24/30 | Alias recipe (line 209): actionable. Terminology mapping table (line 208): actionable. Surface prefix on "functional test" (line 210): partially actionable -- works in documentation but not in justfile recipe names where "functional" appears without qualification. Classification binary extensibility (line 212): actionable ("新增第三分类而非强行归入现有二分法"). Terminology regression (line 213): partially actionable -- "review 流程中检查" relies on human process, not automated enforcement. Partial adoption (line 214): actionable ("单次 PR 中完成，不采用分批更新策略"). Deduction (-6): the terminology regression mitigation relies on "review 流程" (human process) and a concept doc annotation ("e2e 仅用于 Web/Mobile"). No automated enforcement (no CI check, no linting rule, no grep gate) despite this being a proposal *about terminology*. A terminology governance proposal that does not include automated terminology enforcement is self-undermining. |

### 9. Success Criteria: 68/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 26/30 | SC1 (concept doc with 5 types): testable. SC2 (guide.md mapping): testable. SC3 (no "e2e" outside Web/Mobile): testable via grep. SC4 (skills use surface-specific names): testable via audit. SC5 (skills reference concept doc): testable. SC6 (index.json task type format): testable. SC7 (business rules doc updated): testable. SC8 (old recipes work): testable. SC9 (output format): partially testable. Deduction (-4): SC9 "质量门报告中不同测试类型的执行结果分类展示" -- "分类展示" is still ambiguous. Separate sections? Separate labels? A filter option? Testability depends on interpretation. SC3 now covers both English "e2e" and Chinese "端到端" which is good, but creates a potential conflict with the classification standard declaration's use of "端到端" as a general category label (see contradiction #4). |
| Coverage is complete | 22/25 | SCs cover: documentation (SC1-2), terminology elimination (SC3), skill adoption (SC4-5), task type naming (SC6), business rules sync (SC7), backward compat (SC8), output format (SC9). Deduction (-3): (1) no SC for transition period completion (alias removal). NFR #6 (line 103) mentions it but no SC enforces it. (2) no SC for generated test code comments/tags (InScope item line 194). SC4 covers skill files but not the *output* of gen-test-scripts. |
| SC internal consistency | 20/25 | SC4/SC5 overlap is the primary issue. SC4 requires surface-specific names; SC5 requires referencing the concept doc. If SC4 is satisfied, SC5 is nearly automatic unless "引用" means explicit text citation. The relationship is ambiguous -- SC5 could be trivially satisfied by SC4, making it redundant, or it could require a stronger form of referencing (explicit citation), making it a distinct requirement. Without clarification, the SC set has an ambiguous satisfiability relationship. SC3 vs. classification standard usage of "端到端" creates a potential contradiction (see contradiction #4). |

### 10. Logical Consistency: 78/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 32/35 | Mapping model directly addresses "e2e" imprecision. Eval-design scoring issue (line 22) demonstrated as a concrete symptom; surface-specific names would resolve it by giving eval-design the correct test type to evaluate against. Deduction (-3): the classification standard's definitional framework (functional = protocol boundary, e2e = device automation) is a stronger claim than the problem requires. The problem says "they are called the same thing but they are different." The solution says "they are different in this specific taxonomic way." The taxonomy may be correct but it is a stronger commitment than the problem demands. |
| Scope <-> Solution <-> SC aligned | 25/30 | Most in-scope items now have corresponding SCs. SC6 covers task type naming. SC7 covers business rules docs. SC9 covers test output. Deduction (-5): (1) InScope "更新 gen-test-scripts 输出的测试代码中的注释/标签" (line 194) has no explicit SC -- SC4 covers skill files but not their generated output. (2) InScope "更新 run-tests 的测试输出格式" (line 195) is partially covered by SC9 but "测试输出格式" includes more than suite names and labels (log format, error messages, progress indicators). (3) No SC for transition period completion. |
| Requirements <-> Solution coherent | 21/25 | Key scenarios map to solution components. Key Scenario 4 (task tracking in index.json) now has SC6. Key Scenario 5 (quality gate) has SC9. Edge cases addressed. Deduction (-4): Key Scenario 2 (test code generation with type in filename/comments) has no explicit SC. The solution describes updating gen-test-scripts output but no SC verifies the *content* of generated test files. |

---

## Phase 3: Blindspot Hunt

### [blindspot] SC3 creates a compliance trap for the classification standard declaration itself
SC3 (line 220) restricts "端到端" to "仅出现在 Web/Mobile surface 的测试类型名称和定义中." The classification standard declaration (line 51) uses "端到端测试" as a general categorical term: "**端到端测试**通过设备级自动化（浏览器驱动、移动设备自动化）模拟真实用户操作序列." This is a meta-level definitional use, not within a specific Web/Mobile surface's test type name or definition. If SC3 is enforced strictly via grep, the classification standard declaration itself would be non-compliant because it uses "端到端" in a general context, not specifically within "Web/Mobile surface 的测试类型名称和定义中." The SC needs a scoped exception for definitional/meta-level usage.

### [blindspot] No SC for generated test code content
InScope line 194 says "更新 gen-test-scripts 输出的测试代码中的注释/标签." But no SC verifies what the generated test code actually looks like. SC4 verifies that skill files use surface-specific names, and SC9 verifies suite names in test output, but neither checks the content of generated test files. A test file could have the correct suite name but still contain "e2e" in comments, variable names, or test descriptions. The gap between "skill file uses correct terminology" and "generated output uses correct terminology" is not covered.

### [blindspot] The "single PR" mitigation for partial adoption risk is fragile
Risk #7 (line 214) mitigates partial adoption by saying "所有 skill 文件的术语更新在单次 PR 中完成（约 27 个文件），不采用分批更新策略，消除混用窗口期." A single PR touching 27 files across 4 skill areas is a large atomic change. The risk of merge conflicts, review fatigue, and partial rollback (if one file causes issues) increases with PR size. The mitigation trades one risk (partial adoption) for another (large PR fragility). This trade-off is not acknowledged.

### [blindspot] The version-based alias removal mechanism assumes users re-run init-justfile
NFR #6 (line 103) says "alias 移除作为一个独立 task 记录在 Forge 版本规划中，由 init-justfile skill 的版本号硬编码机制自动触发." Line 73 says "用户重新运行 forge init-justfile 时自动更新." But what if a user does not re-run `forge init-justfile`? Their old justfile will continue to work (alias still exists in their local file), but the deprecated aliases will never be removed. There is no mechanism to detect or notify users with stale justfiles. The "automatic" removal is only automatic for users who actively re-run the command.

### [blindspot] The proposal does not define what "引用概念文档" means in SC5
SC5 (line 222) requires "引用概念文档中的测试类型定义." What constitutes a "引用"? A hyperlink? A file path reference? A mention of the document name? Using the same terminology defined in the concept doc? The vagueness of "引用" makes SC5's scope ambiguous. If "引用" means "use the same terminology," SC5 is subsumed by SC4. If "引用" means "explicitly cite the concept doc by name/path," SC5 is a distinct and stronger requirement. Without clarification, SC5's satisfiability is undefined.

### [blindspot] Multi-surface projects may have conflicting test type execution in CI
Edge case #1 (line 107) says multi-surface projects get independent recipes (`test-cli-functional`, `test-api-functional`) plus an aggregate `test` recipe. But the proposal does not specify what the aggregate `test` recipe does -- does it run all surface-specific tests sequentially? In parallel? What is the expected behavior if one surface's tests fail -- do others continue? The aggregate recipe's behavior is unspecified, which matters for CI pipeline configuration.

---

## Bias Detection Report

Annotated (pre-revised) regions:
- Test Type Mapping table (lines 34-42) -- pre-revised: high. 1 attack point (verification dimension granularity inconsistency).
- Classification Standard Declaration (lines 44-53) -- pre-revised: medium. 2 attack points (SC3 compliance trap; binary justification complexity undermines simplicity claim).
- Semantic Definitions (lines 57-63) -- pre-revised: high. 1 attack point (Mobile best-effort threshold undefined).
- Assumptions Challenged row 1 (line 177) -- pre-revised: medium. 0 new attack points (rhetorical loading carried over, not re-attacked in this iteration).

Total annotated: 4 attack points / 4 paragraphs = density 1.00

Unannotated regions:
- Problem section (lines 9-28): 1 attack point (problem statement mixes naming vs. categorization claims)
- User-Facing Experience (lines 65-67): 1 attack point (failure message not described)
- Technical Direction (lines 69-73): 1 attack point ("纯文档 + 命名变更" disclaim contradicts InScope code changes)
- Requirements Analysis (lines 79-109): 2 attack points (edge case: surface type change; NFR: version counting undefined)
- Alternatives & Industry Benchmarking (lines 111-136): 2 attack points (trade-offs still unquantified; no alternative taxonomy explored)
- Feasibility (lines 138-171): 1 attack point (forge surfaces unverified)
- Scope (lines 181-202): 1 attack point (business rules doc plural vague)
- Key Risks (lines 204-214): 2 attack points (generated test code regression risk; single PR fragility)
- Success Criteria (lines 216-226): 2 attack points (no SC for generated test code content; SC5 "引用" ambiguous)

Total unannotated: 13 attack points / 13 paragraphs = density 1.00

Ratio (annotated/unannotated): 1.00

Annotated and unannotated regions have equal attack density (1.00). The pre-revisions did not create a bias toward either weaker or stronger scrutiny. No significant bias detected.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 90 | 110 |
| Solution Clarity | 100 | 120 |
| Industry Benchmarking | 74 | 120 |
| Requirements Completeness | 88 | 110 |
| Solution Creativity | 54 | 100 |
| Feasibility | 84 | 100 |
| Scope Definition | 74 | 80 |
| Risk Assessment | 72 | 90 |
| Success Criteria | 68 | 80 |
| Logical Consistency | 78 | 90 |
| **Total** | **782** | **1000** |

### Iteration-over-Iteration Delta

| Dimension | Iter 1 | Iter 2 | Iter 3 | Delta (2->3) |
|-----------|--------|--------|--------|--------------|
| Problem Definition | 82 | 88 | 90 | +2 |
| Solution Clarity | 90 | 96 | 100 | +4 |
| Industry Benchmarking | 68 | 72 | 74 | +2 |
| Requirements Completeness | 65 | 82 | 88 | +6 |
| Solution Creativity | 50 | 52 | 54 | +2 |
| Feasibility | 80 | 85 | 84 | -1 |
| Scope Definition | 62 | 72 | 74 | +2 |
| Risk Assessment | 56 | 62 | 72 | +10 |
| Success Criteria | 52 | 62 | 68 | +6 |
| Logical Consistency | 68 | 74 | 78 | +4 |
| **Total** | **673** | **745** | **782** | **+37** |

The iteration 3 revision made meaningful improvements in Risk Assessment (+10, from adding 3 new risks and improving mitigation actionability), Requirements Completeness (+6, from expanding NFRs to 6 items), and Success Criteria (+6, from adding SC6-SC8 to cover task type naming, business rules, and backward compatibility). All other dimensions saw marginal improvements.

The proposal remains 118 points below the 900 target. The primary blockers are:
1. **Industry Benchmarking** (-46 from max): trade-offs still unquantified, benchmark justifications shallow, no alternative taxonomies explored
2. **Solution Creativity** (-46 from max): limited novelty beyond renaming, no cross-domain inspiration
3. **Technical Direction** (within Solution Clarity, -17 from max): still underspecified for task type propagation, generated test code, and the "pure documentation" claim contradicts code changes in scope
4. **Risk Assessment** (-18 from max): missing generated test code regression risk, no automated terminology enforcement

### Top 5 Remaining Issues (if another iteration were allowed)

1. **Quantify trade-off costs** (+15 potential): Replace one-sentence mitigations with quantified costs. "27 files across 4 skill areas, estimated 2-3 hours per skill area." "Learning curve: 5 new terms, deterministic mapping table." Reference the blast radius data already in the proposal.

2. **Add automated terminology enforcement** (+10 potential): Add a CI grep gate or linting rule as a mitigation for the terminology regression risk. A proposal about terminology precision that relies solely on human review for ongoing enforcement undermines its own premise.

3. **Specify generated test code before/after** (+10 potential): Add concrete examples of what a generated test file looks like before (with "e2e" labels) and after (with surface-specific labels). This addresses the InScope item with no SC and the missing risk around test code regeneration.

4. **Resolve SC3 vs. classification standard contradiction** (+5 potential): Either add a scoped exception in SC3 for meta-level definitional usage of "端到端," or restructure the classification standard to avoid using "端到端" as a general category label.

5. **Explore alternative taxonomies** (+8 potential): Add at least one alternative that proposes a *different* naming taxonomy (e.g., tagging-based without renaming, execution-model-only naming without the functional/e2e binary). This would strengthen Industry Benchmarking and demonstrate that the selected approach was chosen over genuinely different alternatives, not just over straw men and the status quo.

---

## Verdict

The proposal is well-structured and demonstrates clear understanding of the problem domain. The core insight (name tests by what they actually do) is sound. The primary remaining weaknesses are in *depth* rather than *structure*: trade-off quantification, automated enforcement, generated code handling, and benchmark analysis depth. These are all addressable but collectively represent 118 points below target.

This is the final iteration. The proposal does not meet the 900-point target at 782/1000.
