# Baseline Evaluation Report: Submit-Task Record Regression Proposal

**Evaluator:** CTO Persona (Adversarial Scoring)
**Date:** 2026-05-24
**Iteration:** 1 (Baseline)
**Document:** `docs/proposals/submit-task-record-regression/proposal.md`

---

## Total Score: 643/1000

---

## Dimension Scores

### 1. Problem Definition: 84/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 32/40 | Core problem ("缺乏系统性验证手段") is clear, but scope is ambiguous: the title and problem statement imply testing submit-task's type dispatch logic, while the actual test target is the Go CLI rendering pipeline. Two readers could legitimately disagree on what is being tested. |
| Evidence provided | 30/40 | Four concrete evidence items with specific feature names and task type counts. However, the core risk stated in evidence -- "record-format 模板中的示例 JSON 可能与 Go 端实际接受的 schema 存在偏差" (line 17) -- cannot be detected by the proposed golden dataset test, since historical records already passed Go validation. |
| Urgency justified | 22/30 | Reasonable urgency argument: "随着 task type 持续增加" + "成本最低——数据已存在" (line 22). But lacks quantification of the cost of inaction -- what happens when a type dispatch bug surfaces? How expensive is it to fix post-hoc? |

### 2. Solution Clarity: 92/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 36/40 | "从已完成的 feature 提取历史 record 作为 golden dataset，按 task type 分组建立 Go CLI 层的回归测试：validateRecordData 校验 + markdown 模板渲染对比" (line 26) -- specific and reproducible. Feature Sources table and Coverage Matrix provide detailed data provenance. |
| User-facing behavior described | 28/45 | As internal test infrastructure, the "user" is a developer. But developer experience is underspecified: missing test file structure, fixture file format, test function signatures. The success criterion "新增 fixture 时只需复制文件 + 加一行测试用例" (line 155) is inaccurate per freeform review analysis. |
| Technical direction clear | 28/35 | Mentions `validateRecordData()`, `RenderRecord()`, table-driven by task type, `-update` flag. Direction is clear. But fixture format is ambiguous: is it JSON input + expected markdown output pairs? Or just one of them? |

### 3. Industry Benchmarking: 64/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 20/40 | Only generic reference to "Go testing `-update` flag" and "jsonschema". No links, no articles, no specific library names beyond the standard library. For a proposal centered on golden/snapshot testing, this section should reference established patterns from projects like Hugo, Terraform, or other Go CLI tools that use golden tests. |
| At least 3 meaningful alternatives | 18/30 | Three alternatives listed (Do nothing, JSON Schema, Golden dataset). "Do nothing" is a valid baseline. But "JSON Schema 验证" is presented as a straw man -- its con is "偏离 golden dataset 对比的目标" (line 65), which uses the conclusion to dismiss the alternative rather than genuinely evaluating it. Missing: property-based testing, contract testing, template unit testing as alternatives. |
| Honest trade-off comparison | 12/25 | Golden dataset cons listed as only "fixture 需随格式演进更新" (line 66). This understates the real costs: (1) the `-update` flag silently overwrites fixtures with no diff review mechanism, (2) fixture maintenance involves judgment calls when historical records don't match current templates (freeform review steps 4-5), (3) the approach cannot detect the specific schema drift it claims to address. |
| Chosen approach justified against benchmarks | 14/25 | "与目标最匹配" (line 66) is insufficient justification. No explanation of why JSON Schema's "indirect validation" is inadequate -- if the goal is detecting schema drift between templates and Go structs (line 17), JSON Schema validation might actually be more direct. |

### 4. Requirements Completeness: 76/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 28/40 | Six scenarios listed covering happy path, schema drift, template rendering, missing fields, legacy types. However, "Schema 偏差" (line 38) and "缺失字段" (line 40) are phantom scenarios -- the proposed golden dataset test cannot detect template-to-Go schema mismatches because it feeds historical (already-valid) data into Go functions. These scenarios create a false sense of coverage. |
| Non-functional requirements | 24/40 | Three NFRs listed: CI runtime < 30s, extensibility, task-type grouping. Missing: fixture storage size, CI environment requirements (Go version?), test stability/flakiness considerations, fixture update frequency expectations. The 30s CI budget (line 45) is asserted without decomposition -- freeform review estimates the realistic time at under 5 seconds. |
| Constraints & dependencies | 24/30 | Two constraints clearly stated: depends on existing Go CLI functions, will not modify validation logic. "不修改 Go 端校验逻辑（只测试，不改行为）" (line 52) is an important constraint well-articulated. |

### 5. Solution Creativity: 28/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 8/40 | The document itself admits "Golden dataset 回归测试是标准的软件工程实践，无特别创新" (line 30). The claimed differentiator -- "端到端正确性" -- is inaccurate; only the Go CLI rendering pipeline is tested, not the submit-task skill's template selection. |
| Cross-domain inspiration | 5/35 | None. Pure application of standard Go testing practices with no borrowing from other domains. |
| Simplicity of insight | 15/25 | The core insight is elegantly simple: extract existing historical records as test fixtures. The "data already exists" observation is the genuine insight here. But this simplicity is undermined by the scope-title mismatch. |

### 6. Feasibility: 88/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 38/40 | Highly feasible. Functions exist, data exists, Go testing framework is mature. No technical blockers identified. |
| Resource & timeline feasibility | 22/30 | "预计 5-8 个 coding task" (line 76) is reasonable for the mechanical extraction and test writing. However, freeform review correctly notes that steps 4-5 of fixture creation (verifying rendered markdown matches historical records, handling mismatches due to template iteration) involve judgment work that scales linearly with template changes. |
| Dependency readiness | 28/30 | "无外部依赖。所有数据在本地仓库" (line 80) -- accurate and unambiguous. |

### 7. Scope Definition: 64/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | Four concrete deliverables: extract fixtures, write Go tests, fix template issues, cover 12 task types. Each is a verifiable deliverable. |
| Out-of-scope explicitly listed | 22/25 | Four explicit exclusions: test/validation categories, validation rule modifications, new CLI commands, LLM determinism testing. Clear and justified. |
| Scope is bounded | 16/25 | Quantified boundary: 12 types x 2-3 records = ~30 fixtures, 5-8 coding tasks. But the fundamental scope mismatch -- the title promises testing submit-task type dispatch, the implementation tests only Go CLI rendering -- is a scope definition failure that propagates throughout the document. |

### 8. Risk Assessment: 58/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 20/30 | Four risks listed. Missing critical risks: (1) the `-update` flag silently overwrites fixtures with no diff review gate, (2) the approach fundamentally cannot detect the schema drift it claims to address (template JSON examples vs Go struct). These are not edge cases -- they go to the core value proposition. |
| Likelihood + impact rated | 22/30 | Uses M/L matrix consistently. But "fixture 数量大导致测试维护成本高" rated L/L is optimistic given freeform review's analysis of the 5-step fixture maintenance process. |
| Mitigations are actionable | 16/30 | "确认是有意变更则更新 fixture" (line 145) is a statement of intent, not an actionable mechanism. "明确区分两种 type 的 fixture，验证各自走正确的模板" (line 148) is a mitigation for a non-problem -- both `fix` and `coding.fix` dispatch to the same template. Only "从已提交的 markdown records 反向验证" (line 146) is a genuinely actionable mitigation. |

### 9. Success Criteria: 46/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 30/55 | "≥12 种 task type 的 golden dataset fixture 建立完成" (line 152) -- measurable (count files). "Go CLI 回归测试通过 go test 运行" (line 153) -- testable. But "record-format 模板中发现的问题全部修复" (line 154) -- "全部修复" is unmeasurable when the number of problems is unknown. "新增 fixture 时只需复制文件 + 加一行测试用例" (line 155) -- inaccurate per freeform review analysis. "submit-task 为每种 task type 选择的模板经过 golden dataset 验证" (line 156) -- promises something the Go-level test cannot deliver. |
| Coverage is complete | 16/25 | Criteria cover fixture creation, test execution, template fixes, extensibility, and type dispatch verification. But the last criterion is unachievable given the chosen approach, and the template-fix criterion is vague. |

### 10. Logical Consistency: 43/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 15/35 | Critical logical break. Problem: "无法确认 submit-task 为每种 task type 选择的模板、填充的字段、渲染的 markdown 是否都正确" (line 11). Solution: Go CLI-level tests on `validateRecordData` + `RenderRecord`. The submit-task skill's template selection logic (TypeScript/YAML code in `plugins/forge/skills/submit-task/`) is completely outside the test surface. The solution addresses roughly half of the stated problem. |
| Scope <-> Solution <-> Success Criteria aligned | 16/30 | Scope says "覆盖 12 种 task type" -- solution can deliver. Success criteria says "submit-task 为每种 task type 选择的模板经过 golden dataset 验证" (line 156) -- solution cannot deliver. This is a direct alignment failure between success criteria and solution. |
| Requirements <-> Solution coherent | 12/25 | "Schema 偏差" (line 38) and "缺失字段" (line 40) are listed as requirements but the golden dataset approach cannot detect them -- historical records already pass validation, so they cannot reveal template-to-Go schema mismatches. These are orphan requirements with no corresponding solution capability. |

---

## Blindspots

1. **[blindspot] Title-solution systematic mismatch.** The title "Submit-Task 按任务类型构建 Record 的回归验证" and problem statement promise testing the submit-task skill's type dispatch logic, but the implementation only tests the Go CLI rendering pipeline. The submit-task skill (TypeScript/YAML) is not under test. The document never acknowledges this limitation. Freeform review diagnosed this precisely: "A Go-level test cannot verify whether submit-task selected the correct record-format template for a given task type."

2. **[blindspot] "端到端" claim is inaccurate.** Innovation Highlights claims "验证 submit-task 的类型分发逻辑（task type → template → fields → rendered markdown）的端到端正确性" (line 30), but the actual test starts from `validateRecordData` with fixed JSON input, skipping the "task type → template" step entirely. This is not end-to-end testing.

3. **[blindspot] Coverage Matrix references features not in the Feature Sources table.** `doc.eval` lists "eval-freeform-expert, run-tasks-git-status, enforce-forge-task-add, run-tests-decouple" (line 120) and `doc.drift` lists "worktree-unpushed, auto-task-main, forge-research, list-tasks, cli-restructure, refactor-impact" (line 122) -- none of these appear in the Feature Sources table (lines 95-106). These appear to be task slugs or feature names from a different naming scheme, making the coverage data unverifiable.

4. **[blindspot] "原提案 4 个 feature" reference lacks context.** The Assumptions Challenged table references "原提案 4 个 feature 足够覆盖" (line 89), implying a prior proposal version that is neither linked nor explained. A reader cannot evaluate this assumption without access to the original version.

5. **[blindspot] The proposal's own evidence undermines its core value proposition.** Evidence item "record-format 模板中的示例 JSON 可能与 Go 端实际接受的 schema 存在偏差" (line 17) identifies the most valuable thing to detect, but the golden dataset approach (feeding historical records through Go functions) cannot detect it. Historical records were generated by agents following the templates and then validated by Go -- if a template was wrong, the agent adapted, and the resulting JSON passed validation. The test would only catch future Go-side regressions, not current template-documentation-to-Go-struct gaps.

---

## Summary of Required Revisions

1. **Resolve the title/solution mismatch.** Either (a) rename and re-scope to accurately reflect that only the Go CLI rendering pipeline is under test, or (b) add a test layer that actually validates submit-task's template selection logic.

2. **Remove or reframe phantom scenarios and success criteria.** "Schema 偏差" and "缺失字段" scenarios, and the success criterion about submit-task template verification, cannot be delivered by the proposed approach. Either remove them or add a complementary test mechanism.

3. **Fix unverifiable coverage data.** Cross-reference all Feature Sources in the Coverage Matrix against the Feature Sources table, or explain the naming discrepancy.

4. **Strengthen risk mitigations.** The `-update` flag mitigation needs a concrete mechanism (e.g., diff gating), not a social convention.

5. **Deepen industry benchmarking.** Reference specific projects, libraries, or articles that use golden/snapshot testing. Present alternatives honestly rather than as straw men.
