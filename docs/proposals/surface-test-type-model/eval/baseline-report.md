# Baseline Evaluation Report: Surface-Specific Test Type Model

**Iteration**: 0 (Baseline)
**Date**: 2026-05-26
**Evaluator**: CTO Adversary

---

## Phase 1: Reasoning Audit

### Problem -> Solution
The problem is clearly stated: "e2e" is a misleading umbrella label for tests that differ fundamentally by surface. The proposed solution (Surface -> Test Type mapping) directly addresses this. **Chain holds.**

### Solution -> Evidence
Evidence (6 items from directory structure, justfile, task types, documentation, gen-test-scripts types/) shows the infrastructure already differentiates by surface. This supports feasibility but does not validate the *naming choices* in the mapping (e.g., why "integration" for CLI but "contract" for API). **Chain partially holds — evidence supports "something should change" but not the specific taxonomy.**

### Evidence -> Success Criteria
Success criteria are concrete (document written, e2e label eliminated, skills updated). They test the *execution* of the rename, not whether the new taxonomy is correct. **Chain holds for implementation verification, not for taxonomy quality.**

### Self-contradiction Check
**Flagged:** The proposal claims to bring precision but introduces imprecise naming ("integration" for CLI lacks definition, "contract" for API conflicts with industry meaning, "e2e" reserved for Web but Mobile is logically identical). The freeform expert review identified this correctly. **Partial self-contradiction: precision is the goal but not consistently achieved.**

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 78/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 34/40 | The core problem is unambiguous: "e2e" is a misleading label for heterogeneous test types. Two readers would interpret the problem the same way. Minor deduction: the problem scope is narrow (terminology) but framed as if it were architecturally significant. |
| Evidence provided | 24/40 | Six evidence items are listed, all pointing to code/directory structure. However, none are *user-facing impact* evidence. No user complaints, no misdiagnosed bugs caused by the terminology, no concrete instance where the "e2e" label led to a wrong decision. The evidence shows the *state* of the codebase, not the *pain* it causes. Quote: "目录结构：测试代码统一放入 tests/<journey/>" — this is a factual observation, not impact evidence. |
| Urgency justified | 20/30 | Quote: "随着 Forge 支持的项目类型增多" — the urgency is forward-looking ("will get worse") rather than demonstrated by present pain. No concrete cost-of-delay quantified. What happens if this is done in 3 months instead of now? Never answered. |

### 2. Solution Clarity: 82/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 34/40 | The mapping table is specific: 5 surfaces, 5 test types, with verification dimensions and execution models. A reader can explain back what will be built. Deduction: the "semantic definitions" section uses terms ("black-box", "semi-black-box") that are themselves undefined, creating a circular reference. |
| User-facing behavior described | 35/45 | Quote: "justfile recipe 按测试类型命名（如 test-cli-integration、test-api-contract）" — the user-visible change is clear (recipe names, task type names, documentation). Deduction: what does the end user *experience* differently when running tests? The output format, error messages, and test reports are not described. Only the names change. |
| Technical direction clear | 13/35 | Quote: "纯文档 + 命名变更，不涉及核心逻辑重构" — the proposal explicitly says there is no technical implementation. The "technical direction" is renaming files and updating strings. No mention of how the justfile alias backward-compatibility works mechanically, no migration path for existing projects, no technical specification for how `test.gen-scripts.cli` differs from `test.gen-scripts` in the task system. |

### 3. Industry Benchmarking: 62/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 28/40 | Four references: Go build tags, Spring Boot @Tag, Playwright project configs, Postman/Newman. These are genuine industry patterns. Deduction: none are analyzed in depth — just one-line mentions. No discussion of *why* these approaches work in their contexts or how they map to Forge's constraints. |
| At least 3 meaningful alternatives | 18/30 | Four alternatives listed including "do nothing". However, the comparison table is biased: alternatives 2 and 3 ("统一改名高级测试" and "引入标准测试分层") are presented as straw men. Quote: "仍然是一个笼统概念，不解决类型错配问题" — this rejects an option by assertion, not by analysis. The "行业标准" option is dismissed with: "不适合 Forge 的场景——CLI 测试不是传统意义上的 integration test" — but then the proposal *itself* calls CLI tests "集成测试" (integration test), which is exactly the industry term it just rejected. This is a logical contradiction. |
| Honest trade-off comparison | 8/25 | The "Cons" column for the selected approach says only "需要更新多个文件和概念" — this is trivially true for any rename. No real trade-offs identified: what about the learning curve for new users? What about documentation migration cost? What about the semantic conflict with "contract testing" (Pact)? The freeform expert review identified these, but the proposal itself does not. |
| Chosen approach justified against benchmarks | 8/25 | Quote: "最小惊讶原则——名称匹配实际行为" — this is the only justification principle cited. It is a valid principle but not sufficient: "API 契约测试" violates the principle of least surprise for anyone familiar with Pact. The proposal does not acknowledge this conflict. |

### 4. Requirements Completeness: 62/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 28/40 | Five key scenarios listed (concept lookup, test generation, test execution, task tracking, quality gate). These cover the happy path. Deduction: no edge cases identified. What happens when a project has multiple surfaces? What about projects with no surface detected? What about the transition period where some files use old terminology and others use new? |
| Non-functional requirements | 14/40 | No NFRs mentioned. Zero. No performance considerations (does per-surface test execution change runtime?), no backward compatibility requirements quantified, no accessibility, no security. The "backward-compatible alias" for justfile is mentioned in risk mitigation but not as a requirement. |
| Constraints & dependencies | 20/30 | Four constraints listed. Quote: "现有 skill 的 types/ 和 rules/surfaces/ 文件已按 surface 分化" — this is correct and specific. Deduction: no mention of dependency on the `forge surfaces` CLI output format, no mention of how the task-lifecycle business rule's reserved type list constrains the naming, and no timeline or version-bound constraint. |

### 5. Solution Creativity: 45/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 15/40 | The proposal itself acknowledges: "这不是创造新概念，而是精确命名已有实践" — it explicitly disclaims novelty. The mapping is a 1:1 rename of existing surface names to test type names. There is no innovation beyond the existing `types/` directory structure. |
| Cross-domain inspiration | 10/35 | No cross-domain ideas. The proposal looks only at Forge's internal structure. No consideration of how other test frameworks (Jest, pytest, RSpec) handle multi-surface test classification. No borrowing from taxonomy or ontology design principles. |
| Simplicity of insight | 20/25 | The insight is genuinely simple and the proposal admits it: the differentiation already exists, just name it. This is elegant in its simplicity. Deduction: the execution is not as simple as the insight, because the naming choices introduce new ambiguities (as the expert review demonstrated). |

### 6. Feasibility: 78/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | Strong: the proposal lists 4 existing infrastructure elements that already differentiate by surface. The change is primarily renaming, which is technically trivial. Deduction: no analysis of the task type system's constraints on naming — can `test.gen-scripts.cli` be introduced without breaking the task lifecycle parser? |
| Resource & timeline feasibility | 22/30 | Quote: "纯文档 + 命名变更... 预计工作量：概念参考文档：1 个 doc 任务... 若干 doc 任务... 若干 coding 任务" — "若干" (several) is not an estimate. No timeline given. The next steps section mentions two paths depending on scope size, but does not commit to either. |
| Dependency readiness | 20/30 | The `forge surfaces` CLI is mentioned as existing. But no analysis of whether it outputs the exact surface keys used in the mapping table, or whether it needs updates. No mention of whether the justfile generation skill supports alias recipes. |

### 7. Scope Definition: 58/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 24/30 | Ten in-scope items, most are specific deliverables ("更新 justfile recipe 命名"). Deduction: "编写测试类型概念参考文档" — what goes in this document? The proposal defines the content (5 test types, semantics, dimensions, execution model) but does not specify the document's audience, format, or location. |
| Out-of-scope explicitly listed | 18/25 | Four items explicitly out of scope. Good. Deduction: "测试目录结构的重新组织" is listed as out of scope with a parenthetical "可作为后续优化" — this is scope creep-by-hint. Either it is in scope or it is not. The parenthetical suggests the author is not fully committed to excluding it. |
| Scope is bounded | 16/25 | The scope is bounded by surface types (5 fixed). However, the number of files affected is not quantified. Quote from next steps: "如变更范围 > 15 个 coding 任务" — the author does not know the scope yet, which means the scope section is aspirational rather than definitive. |

### 8. Risk Assessment: 52/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 18/30 | Three risks listed. Only one (justfile recipe rename affecting CI) is operationally meaningful. The other two are low-impact ("术语变更导致文档失效" and "新概念增加用户学习成本"). Missing risks identified by the expert review: semantic conflict with "contract testing" (Pact), inconsistency between Web and Mobile classification, undefined "semi-black-box" concept. |
| Likelihood + impact rated | 14/30 | All three risks have medium/low ratings. The CI disruption risk is rated M likelihood / H impact — fair. But the "术语映射表" risk is rated M/M without justification — how many existing tutorials/documents exist? What is the blast radius? No data. |
| Mitigations are actionable | 20/30 | Quote: "提供向后兼容的 alias recipe（旧名 → 新名），设置过渡期" — this is actionable. Quote: "在变更点添加术语映射表" — also actionable. Quote: "概念文档以一页纸为限" — this is a constraint, not a mitigation for learning cost. The mitigation for learning cost should be about onboarding, not page count. |

### 9. Success Criteria: 48/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 20/30 | SC1 (document completed) is testable. SC3 (no "e2e" outside Web context) is testable via grep. SC2 (guide.md contains mapping) is testable. SC4 (all skills use surface-specific names) is testable via audit. SC5 ("被至少 3 个现有 skill 的 rules 文件引用") is measurable. Deduction: SC5's threshold of 3 is arbitrary — why 3? There are 5 surface types; shouldn't all 5 reference the concept document? |
| Coverage is complete | 10/25 | The success criteria cover documentation and terminology changes but miss: (1) backward compatibility — no SC for "old justfile recipes still work via alias", (2) task type migration — no SC for "existing tasks with old type names are handled", (3) quality gate output — listed as a key scenario but has no corresponding SC. |
| SC internal consistency | 18/25 | SC1-SC5 are internally consistent as a set — they do not contradict each other. Deduction: SC3 ("不再有将所有生成测试统称为 'e2e' 的地方") and SC5 ("概念文档被至少 3 个现有 skill 的 rules 文件引用") have a tension: if all skills must use surface-specific names (SC4), then all 5+ skill rules files should reference the concept document, making SC5's threshold of 3 trivially satisfied and therefore meaningless as a gate. |

### 10. Logical Consistency: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 25/35 | The mapping model directly addresses the "e2e" umbrella problem. Deduction: the proposal rejects "引入标准测试分层" saying "CLI 测试不是传统意义上的 integration test" (quote from alternatives table), but then names CLI tests "CLI 集成测试" (CLI Integration Test). This is a direct self-contradiction. The proposal rejects a standard term and then adopts that same term. |
| Scope <-> Solution <-> SC aligned | 15/30 | In-scope item "更新 task type 命名（携带 surface 信息）" has no corresponding success criterion. Key scenario "质量门报告区分不同测试类型的执行结果" has no in-scope item for quality gate code changes. The alignment has gaps. |
| Requirements <-> Solution coherent | 15/25 | Key scenario 4 (task tracking with surface info) maps to solution's task type naming change. Key scenario 5 (quality gate differentiation) has no solution component — the quality gate is out of scope but the requirement is in scope. This is an orphan requirement. |

---

## Phase 3: Blindspot Hunt

### [blindspot] No migration strategy for existing projects
Quote: "更新 justfile recipe 命名（从 test/test-e2e 到 surface-specific 名称）" — this is listed as in-scope. But existing Forge-generated projects already have `just test` and `just test-e2e` in their justfiles. The proposal provides a backward-compatible alias mitigation in the risk table but never specifies: (1) who updates the existing project justfiles, (2) whether `forge init` re-generates them, or (3) whether users must manually migrate. This is an operational gap.

### [blindspot] The taxonomy is presented as final but has not been validated
The proposal presents the 5-entry mapping table as the solution. The freeform expert review identified 4 specific naming issues (CLI "integration", API "contract", TUI "semi-black-box", Mobile vs Web inconsistency). The proposal has no validation mechanism — no user testing, no team review checkpoint, no iteration before implementation. A naming proposal should be validated before committing to rename dozens of files.

### [blindspot] No cost quantification
Quote: "需要更新多个文件和概念" — the entire cost of the proposal is summarized as "multiple files." How many files? The proposal says to grep for "e2e" but does not report the count. Without knowing the blast radius, the scope section cannot be considered bounded.

### [blindspot] Assumptions Challenged section is advocacy, not analysis
The "Assumptions Challenged" section uses leading rhetorical questions to steer toward the desired conclusion. Quote: "如果 CLI 测试是 e2e，那它端到端测试了什么？答案是一个子进程调用" — this is a rhetorical device, not an assumption analysis. A genuine assumption challenge would also ask: "Is it possible that 'e2e' in Forge's context means 'end-to-end from the user's perspective' rather than 'end-to-end across system layers'?" The section assumes one interpretation of "e2e" and knocks it down, without considering alternative interpretations.

### [blindspot] No success criterion for the taxonomy's correctness
All 5 success criteria test whether the rename was *executed*, not whether it was *correct*. There is no SC like: "Users can correctly identify the test type for their surface without reading the definition document." A rename that is completed but introduces confusing new names would pass all current SCs.

### [blindspot] Missing consideration of the "no surface detected" case
The proposal assumes every project has exactly one surface. What about libraries, plugins, or multi-surface projects? The mapping table has no entry for these cases, and the requirements analysis does not mention them.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 78 | 110 |
| Solution Clarity | 82 | 120 |
| Industry Benchmarking | 62 | 120 |
| Requirements Completeness | 62 | 110 |
| Solution Creativity | 45 | 100 |
| Feasibility | 78 | 100 |
| Scope Definition | 58 | 80 |
| Risk Assessment | 52 | 90 |
| Success Criteria | 48 | 80 |
| Logical Consistency | 55 | 90 |
| **Total** | **620** | **1000** |
