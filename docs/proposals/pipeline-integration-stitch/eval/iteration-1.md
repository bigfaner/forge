# Evaluation Report: Pipeline Integration Stitch — Iteration 1

**Evaluator**: CTO (Adversarial)
**Date**: 2026-05-24
**Document**: `docs/proposals/pipeline-integration-stitch/proposal.md`
**Rubric**: `plugins/forge/skills/eval/rubrics/proposal.md`

---

## Phase 1: Pre-Score Anchors (Reasoning Audit)

### Problem -> Solution Trace
The problem is real and well-scoped: 14 integration gaps between two completed proposals. The auto-discovery solution directly addresses the root cause (hand-written maps missed in prompt.go). This is a strong problem-solution fit.

### Solution -> Evidence Trace
Evidence is concrete: specific files, functions, and failure modes are cited. However, the evidence is entirely internal — no user reports, no incident data, no quantified impact ("how many features failed?"). For an internal tool proposal, this is acceptable but not exemplary.

### Evidence -> Success Criteria Trace
Success criteria cover most P0/P1 items but have notable gaps:
- No criterion for verifying `clean-code.md` rename doesn't break references
- No criterion for verifying `eval.*` in `CategoryTest` produces correct submit-task behavior
- No criterion for verifying `test.gen-and-run` removal produces helpful error messages

### Self-Contradiction Check
- The proposal claims to fix `eval.*` categorization but introduces a new semantic mismatch by putting eval tasks into `CategoryTest` (which expects test-generation fields at submit time)
- Claims "向后兼容" as NFR but does not specify how backward-compatible errors differ from generic errors
- Claims auto-discovery eliminates the bug class, but the failure mode is isomorphic (missed file vs. missed map entry)

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | The 14 integration gaps are enumerated with precise file/function references. Two readers would agree on what's broken. Minor deduction: the distinction between P0/P1/P2 could be more explicit about user-visible impact vs. developer maintenance cost. |
| Evidence provided | 32/40 | Evidence is entirely code-structural analysis — specific files, functions, failure modes. Strong for an internal tool. However, no user-facing incident data, no frequency estimate ("X features have been blocked"), no reproduction steps. The evidence is developer analysis, not empirical data. |
| Urgency justified | 28/30 | P0 means "pipeline 必定失败" — urgency is self-evident. Clear statement: "任何使用 Forge 执行 test pipeline 的 feature 都会失败." The cost of delay is obvious. |

**Subtotal: 98/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 36/40 | Three-pronged approach is clear. Auto-discovery formula (`strings.ReplaceAll(typeName, ".", "-") + ".md"`) is precise. The override for `clean-code.md` is mentioned. A reader could explain back what will be built. |
| User-facing behavior described | 35/45 | The proposal describes internal behavior well but user-facing behavior is limited. What does the developer see when a template is missing? What error message appears for old `test.gen-and-run` tasks? The NFR says "明确错误" but never specifies the error content. The success criteria describe verification commands but not the developer experience. |
| Technical direction clear | 32/35 | Sufficient technical hints: naming convention, embed.FS ReadFile, override map. File-level changes are listed. Missing: no pseudo-code for the auto-discovery function, no mention of `init()` validation, no detail on how `resolveMixedFeatureDeps` would be restructured. |

**Subtotal: 103/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 10/40 | No industry solutions, open-source projects, or published patterns are cited. This is entirely self-invented. No reference to how other task pipeline systems (Airflow, Temporal, GitHub Actions) handle type-to-template mapping or integration testing between independently developed modules. |
| At least 3 meaningful alternatives | 18/30 | Three alternatives are listed: do nothing, manual map completion, auto-discovery + full fix. "Do nothing" is a straw man (P0 = guaranteed failure). "Manual map completion" is the closest to a genuine alternative. No industry-validated solution is included. Missing: init-time validation, code generation from templates, convention-over-configuration frameworks. |
| Honest trade-off comparison | 15/25 | The comparison table is honest about "变更量较大" for the selected approach. But it understates the risk: auto-discovery replaces one runtime failure mode with another (file-not-found vs. map-miss), and the `eval.*` -> CategoryTest change introduces a new semantic mismatch that the trade-off analysis ignores. |
| Chosen approach justified | 15/25 | The justification is essentially "it's more thorough." The argument for auto-discovery over manual map completion is reasonable (eliminates bug class) but overstates the benefit — the developer still must create the template file. The proposal does not justify why auto-discovery is better than init-time validation (which would catch missing templates at startup). |

**Subtotal: 58/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 32/40 | Happy path (auto-discovery works), edge cases (clean-code override, re-index idempotency), and error scenarios (missing template, old index.json) are identified. Gaps: (1) `eval.*` in CategoryTest — what happens at submit-task validation? The proposal doesn't analyze the submit-time requirements for eval tasks. (2) `findFirstTestTaskIdx` quick-mode fallback references `T-quick-gen-and-run*` which would break after removal — not mentioned as a scenario. |
| Non-functional requirements | 28/40 | Two NFRs stated: backward compatibility and compile-time safety. Both are vague. "明确错误而非静默失败" — what error? "运行时验证文件存在" — this is the same failure mode as before, just triggered differently. Missing NFRs: performance (auto-discovery adds FS reads), security (template injection via type names), testability (how to verify auto-discovery without creating real template files). |
| Constraints & dependencies | 25/30 | Clear: upstream proposals must be complete, changes span Go source + plugin data + docs. Good. Missing: no mention of Go version constraints for embed.FS features, no dependency on test infrastructure for verification. |

**Subtotal: 85/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 25/40 | Auto-discovery via naming convention is a standard pattern (convention over configuration — Rails, Spring, etc.). The proposal applies it competently but does not innovate beyond the baseline. The "O(N) to O(1)" framing is overstated — the developer still creates N files, just not N map entries. |
| Cross-domain inspiration | 10/35 | No cross-domain inspiration cited. Convention-over-configuration is well-known but not attributed. No borrowing from build systems, package managers, or plugin architectures that face similar registration problems. |
| Simplicity of insight | 18/25 | The insight is genuinely simple: if naming is deterministic, map entries are redundant. The `strings.ReplaceAll` one-liner is elegant. But the proposal fails to see the equally simple alternative: an `init()` validation loop that checks all registered types against the FS, which would be both simpler and safer. |

**Subtotal: 53/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | Highly feasible. `strings.ReplaceAll` + embed.FS ReadFile is trivial. The clean-code rename is straightforward. Removing dead code is mechanical. Risks: (1) the four new prompt templates require understanding the execution-phase semantics, which the proposal under-specifies; (2) removing all `test.gen-and-run` references requires a complete audit that the proposal does not enumerate. |
| Resource & timeline | 25/30 | "2 coding tasks + doc tasks" is realistic for the scope. The proposal does not overcommit. However, the 4 new prompt templates may require more iteration than estimated, since they serve a different purpose than the autogen templates they reference. |
| Dependency readiness | 28/30 | Upstream proposals are complete. No external dependencies. The only risk is that the existing codebase may have more references to `test.gen-and-run` than identified (the proposal says "8+ files" without being exhaustive). |

**Subtotal: 88/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | Three priority tiers with specific file-level changes. Each item is actionable. Good. Minor gap: "更新 isTestTaskID 语义或文档" is ambiguous — which one? The proposal doesn't decide. |
| Out-of-scope explicitly listed | 22/25 | Four items explicitly out of scope. Clear. Good inclusion of "历史 feature 文件中的 doc.eval 类型" as out of scope with rationale. |
| Scope is bounded | 20/25 | The scope is well-bounded for the P0/P1 fixes. The P2 cleanup ("彻底移除 test.gen-and-run") is potentially open-ended — the proposal says "8+ 文件" without enumerating them all. The scope item "更新引用废弃类型的测试文件" is vague: which files? |

**Subtotal: 68/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 22/30 | Four risks identified. Missing risks: (1) `eval.*` in CategoryTest causing submit-task validation errors — the freeform review identifies this but the proposal's risk table does not; (2) clean-code.md rename breaking external references; (3) `findFirstTestTaskIdx` quick-mode fallback breaking; (4) incomplete removal of `test.gen-and-run` references causing partial compilation failure. |
| Likelihood + impact rated | 22/30 | Ratings are reasonable. "自动发现遗漏 edge case" rated L/H is honest about low likelihood. However, "移除 gen-and-run 后旧 index.json 报错" rated L/M understates the impact — if a team has active features with old index.json files, this blocks their workflow until they regenerate. |
| Mitigations are actionable | 20/30 | "运行全量测试验证" is generic. "参考 autogen.go 中已有的对应模板结构" is insufficient guidance for execution-phase templates. "约定已稳定" is not a mitigation — it's an assertion. "明确错误提示，引导用户重新生成" is the most actionable but doesn't specify the error message. |

**Subtotal: 64/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 40/55 | Six criteria, five use `grep` commands (measurable). "所有现有测试通过" is measurable. Gaps: (1) No criterion for `clean-code.md` rename verification — after renaming, no grep checks for dangling `clean-code` references. (2) No criterion for `eval.*` submit-task behavior — does `forge submit-task` accept eval task submissions with appropriate fields? (3) "mixed feature re-index 幂等性：T-review-doc 依赖不丢失" is testable but vague about how to verify. |
| Coverage is complete | 18/25 | Covers P0 (template mapping), P1 (category, references, re-index), and P2 (gen-and-run removal). Missing: no criterion for P1 item "record-format-doc.md" update verification, no criterion for `validate_index.go` providing migration-aware error messages, no criterion for `findFirstTestTaskIdx` working correctly after `gen-and-run` removal. |

**Subtotal: 58/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 30/35 | Auto-discovery directly fixes P0 (missing template mappings). Category fix addresses P1 (eval misclassification). Cleanup addresses P2. Strong alignment. Deduction: the category fix introduces a new semantic mismatch (eval tasks in CategoryTest expecting test-generation fields) that partially undermines the stated goal of "类型系统一致性". |
| Scope <-> Solution <-> Success Criteria aligned | 22/30 | Scope lists clean-code rename, but success criteria don't verify it. Scope lists "更新 isTestTaskID 语义或文档", but no success criterion checks this. Scope lists "validate_index.go: 移除 T-quick-gen-and-run- 前缀检查", but the success criterion only checks grep for "gen-and-run", which is broader but may miss the specific validation logic. |
| Requirements <-> Solution coherent | 18/25 | Most requirements map to solution items. Gap: NFR "向后兼容" (backward compatibility) maps to "明确错误提示" in the solution, but the solution does not specify what makes this error "backward-compatible" vs. generic. The requirement "编译时安全" (compile-time safety) maps to "运行时验证文件存在" — which is runtime, not compile-time. This is a direct contradiction: the NFR asks for compile-time safety, the solution provides runtime safety. |

**Subtotal: 70/90**

### Cross-Dimension Coherence Check

1. **Problem Definition vs. Success Criteria**: The problem lists 14 items but success criteria cover approximately 8. Items 4 (record-format-doc.md), 9 (validate_index.go), 10 (isTestTaskID gap), 11 (category_test.go), 14 (README/ARCHITECTURE references) have no direct success criterion.

2. **Solution Clarity vs. Feasibility**: The solution claims "2 coding tasks" but the P2 cleanup involves touching 8+ test files, updating multiple documentation files, renaming files, and creating 4 new prompt templates. This may be undercounted.

3. **Risk Assessment vs. Freeform Findings**: The freeform review identified at least 5 additional risks not in the proposal's risk table: CategoryTest/eval semantic mismatch, clean-code rename external references, findFirstTestTaskIdx quick-mode fallback, incomplete gen-and-run reference removal, and execution-phase vs. planning-phase template confusion.

---

## Phase 3: Blindspot Hunt

### Beyond-Rubric Issues

1. **Execution model conflation**: `eval.*` tasks spawn subagents (MainSession: true) while most test tasks do not. Adding `eval.*` to `CategoryTest` without distinguishing this execution model difference creates a latent bug for any future category-based dispatch logic.

2. **`findFirstTestTaskIdx` regression**: `build.go` line ~494 has a quick-mode fallback checking for `T-quick-gen-and-run*`. After removing `test.gen-and-run`, this fallback will never match, potentially causing `findFirstTestTaskIdx` to return -1 or the wrong index for quick-mode tasks. The proposal does not address this.

3. **Incomplete removal checklist**: The proposal says "彻底移除 test.gen-and-run" but does not enumerate all files requiring changes. The freeform review identified at least: `types.go` (constant), `infer.go` (InferType switch case, line 32-33), `prompt.go` (genScriptBases, line 294), `autogen.go` (if any), `validate_index.go` (lines 224-226), and 8+ test files. Without an exhaustive list, partial removal is likely.

4. **Template semantic mismatch**: The four new prompt templates are execution-phase prompts (instructions for the agent at runtime) but the proposal says to "参考 autogen.go 中已有的对应模板结构" — autogen templates are planning-phase (generating .md files). Copying structure from one to the other would produce semantically incorrect prompts.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 98 | 110 |
| Solution Clarity | 103 | 120 |
| Industry Benchmarking | 58 | 120 |
| Requirements Completeness | 85 | 110 |
| Solution Creativity | 53 | 100 |
| Feasibility | 88 | 100 |
| Scope Definition | 68 | 80 |
| Risk Assessment | 64 | 90 |
| Success Criteria | 58 | 80 |
| Logical Consistency | 70 | 90 |
| **Total** | **745** | **1000** |

---

## ATTACKS

1. [Industry Benchmarking]: No industry solutions or published patterns cited — the comparison table contains only self-invented alternatives with one straw man ("do nothing" with P0 guaranteed failure). Must reference at least one industry-validated approach to convention-over-configuration or type registration systems.

2. [Industry Benchmarking]: Trade-off comparison ignores the auto-discovery's isomorphic failure mode — replacing compile-time map errors with runtime file-not-found errors is not a qualitative improvement as claimed. Must acknowledge this and consider init-time validation as an alternative.

3. [Requirements Completeness]: `eval.*` in CategoryTest creates submit-task validation mismatch — the proposal does not analyze what fields `CategoryTest` expects at submission time for eval tasks. Must specify submit-task requirements for eval tasks or create a dedicated `CategoryEval`.

4. [Requirements Completeness]: NFR "编译时安全" contradicts solution "运行时验证文件存在" — the NFR asks for compile-time safety but the auto-discovery solution provides only runtime validation. Must either change the NFR to "运行时安全" or add init-time/compile-time validation.

5. [Risk Assessment]: Risk table missing 4 risks identified in expert review: clean-code.md rename breaking external references, findFirstTestTaskIdx quick-mode fallback regression, execution-phase vs planning-phase template semantic confusion, and eval/CategoryTest submit-task field mismatch. Must add these to the risk table with mitigations.

6. [Risk Assessment]: Mitigation "参考 autogen.go 中已有的对应模板结构" is misleading — autogen templates are planning-phase while the 4 new templates are execution-phase. Copying autogen structure would produce semantically incorrect prompts. Must specify what each new template should contain.

7. [Success Criteria]: No criterion verifying clean-code.md rename doesn't leave dangling references. Must add `grep -r "clean-code" forge-cli/ plugins/` returning zero results (excluding the renamed file).

8. [Success Criteria]: No criterion verifying `eval.*` submit-task behavior — CategoryTest may require inappropriate fields (casesGenerated, scriptsCreated) for eval tasks. Must add a criterion for eval task submission acceptance.

9. [Success Criteria]: No criterion verifying `validate_index.go` provides migration-aware error messages for old `test.gen-and-run` tasks. The NFR requires backward compatibility but no success criterion tests it.

10. [Logical Consistency]: `eval.*` -> CategoryTest partially undermines the stated goal of "类型系统一致性" — eval tasks are quality gates, not test generators. The solution reintroduces a classification mismatch while claiming to fix one. Must either justify why CategoryTest is semantically correct for eval or create a dedicated category.

11. [Logical Consistency]: Scope item "更新 isTestTaskID 语义或文档" is undecided — the proposal says "or" without committing to a direction. Must choose: add T-review-doc to isTestTaskID, or rename the function, or update documentation.

12. [beyond-rubric]: `findFirstTestTaskIdx` in build.go has a quick-mode fallback checking `T-quick-gen-and-run*` which will never match after removal, potentially causing incorrect first-test-task selection for quick-mode pipelines. The proposal does not address this function.

13. [beyond-rubric]: The four new prompt templates are execution-phase prompts distinct from autogen planning-phase templates. The proposal's instruction to "参考 autogen.go 中已有的对应模板结构" may lead to semantically incorrect templates. Must specify execution-phase template structure requirements.

14. [beyond-rubric]: `genScriptBases` dead code `T-quick-gen-and-run` removal must be coordinated with all other reference removals. The proposal does not enumerate all files needing changes, risking partial compilation failures.
