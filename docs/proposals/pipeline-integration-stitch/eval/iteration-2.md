# Evaluation Report: Pipeline Integration Stitch — Iteration 2

**Evaluator**: CTO (Adversarial)
**Date**: 2026-05-24
**Document**: `docs/proposals/pipeline-integration-stitch/proposal.md`
**Rubric**: `plugins/forge/skills/eval/rubrics/proposal.md`
**Previous Score**: 745/1000 (Iteration 1)

---

## Iteration-1 Attack Resolution Audit

Before scoring, I verify which of the 14 iteration-1 attacks were addressed:

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | No industry solutions cited | **Fixed** | New "Industry Patterns" section cites Rails (CoC), Spring Boot (auto-config), ASP.NET Core (controller discovery), GitHub Actions (reusable workflows), Airflow (DAG auto-discovery). |
| 2 | Auto-discovery isomorphic failure mode overstated | **Fixed** | New "Auto-Discovery 的等价故障模式" section explicitly acknowledges the isomorphic failure and adds init-time validation to compensate. Claims "故障检测时机优于原始手写 map" with justification. |
| 3 | eval/CategoryTest mismatch | **Fixed** | Proposal now creates dedicated `CategoryEval` instead of putting eval.* into CategoryTest. CategoryEval requires review fields (summary/findings/severity), not test evidence. |
| 4 | NFR contradiction (compile-time vs runtime) | **Fixed** | NFR now reads "部署时安全" with init-time validation, not "编译时安全". The contradiction is resolved. |
| 5 | Missing risks (4 items) | **Partially Fixed** | clean-code.md rename risk added. findFirstTestTaskIdx regression risk added. Template semantic confusion risk added. eval/CategoryEval field mismatch risk added. However, risk table now has 8 items; coverage is good but one risk's mitigation is still generic. |
| 6 | Misleading template reference to autogen | **Fixed** | Proposal now explicitly states: "这些是执行阶段模板（agent 运行时指令），非 autogen 规划阶段模板...结构参考现有执行阶段模板（如 clean-code.md / code-quality-simplify.md），而非 autogen.go 的规划阶段模板". Risk item #4 also clarifies template requirements. |
| 7 | Missing success criterion for clean-code rename | **Fixed** | Added: `grep -r "clean-code" forge-cli/ plugins/` 仅匹配重命名后的 code-quality-simplify 文件名，无残留 clean-code.md 引用. |
| 8 | Missing success criterion for eval submit-task | **Fixed** | Added: `forge submit-task` 对 eval 任务接受 review 类字段（summary/findings），不要求测试证据（testsPassed/coverage）. |
| 9 | Missing success criterion for validate_index migration | **Fixed** | Added: validate_index.go 对引用 test.gen-and-run 的旧 index.json 返回包含 "deprecated" 和 "regenerate" 的迁移指引错误信息. |
| 10 | eval classification undermines type consistency | **Fixed** | CategoryEval created. Key scenarios section states: "eval 任务归入 CategoryEval 而非 CategoryTest，保持语义正确性". |
| 11 | isTestTaskID scope undecided | **Fixed** | Now reads: "isTestTaskID: 扩展语义覆盖 T-review-doc，并更新函数文档注释说明包含的任务类型范围". Direction is committed. |
| 12 | [beyond-rubric] findFirstTestTaskIdx regression | **Fixed** | Added as explicit risk + success criterion: "findFirstTestTaskIdx 对 quick-mode pipeline 正确返回首个 test task 索引（非 -1）". Also in scope item: "build.go: 更新 findFirstTestTaskIdx 中 quick-mode fallback". |
| 13 | [beyond-rubric] Template phase confusion | **Fixed** | Addressed in scope and risk items. Scope says templates are "执行阶段模板" with 4 required components. Risk item #4 specifies content requirements. |
| 14 | [beyond-rubric] Incomplete removal checklist | **Fixed** | New "gen-and-run Removal Checklist" enumerates 10 specific files/areas. |

**Resolution: 14/14 attacks addressed. 13 fully resolved, 1 partially (risk mitigations still contain some generic items).**

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Problem -> Solution Trace
14 integration gaps with P0/P1/P2 prioritization. Solution is three-pronged: auto-discovery for P0, type fixes for P1, cleanup for P2. Each P0 item maps to a concrete fix. Strong traceability.

### Solution -> Evidence Trace
Evidence is code-structural (files, functions, failure modes). No user-facing incident data, but for an internal toolchain proposal this is acceptable. The new "Assumptions Challenged" table provides structured reasoning.

### Evidence -> Success Criteria Trace
Success criteria now cover all P0 items, most P1 items, and key P2 items. 12 criteria total, up from 6 in iteration 1. Remaining gap: no explicit criterion for `record-format-test.md` update verification (though covered implicitly by the gen-and-run grep criterion).

### Self-Contradiction Check
- NFR "部署时安全" now correctly matches init-time validation (no contradiction).
- CategoryEval is semantically consistent with eval task requirements.
- Auto-discovery failure mode is acknowledged as isomorphic but justified by improved detection timing.
- One remaining tension: the proposal claims "O(1) 约定" but still requires manual template file creation. The O(N) to O(1) framing is about map maintenance, not file creation. This is accurate but could be clearer.

---

## Phase 2: Rubric Scoring with Verification Stance

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | 14 gaps enumerated with precise file/function references, P0/P1/P2 prioritization. Root cause identified (hand-written maps in prompt.go missed during development). Two readers would agree. Minor deduction: P2 items could distinguish "developer inconvenience" from "latent bug" more clearly. |
| Evidence provided | 33/40 | Concrete code-level evidence. Root cause analysis traces the bug to specific developer behavior (updated types.go + autogen.go but missed prompt.go). "Assumptions Challenged" table adds structured reasoning. No empirical data (no user reports, no incident count), but acceptable for internal tool. |
| Urgency justified | 28/30 | P0 = pipeline guaranteed failure. Clear statement of blast radius. Cost of delay is self-evident. |

**Subtotal: 99/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Three-pronged approach with specific formula (`strings.ReplaceAll(typeName, ".", "-") + ".md"`). Override for clean-code.md documented. Auto-discovery mechanism is precisely described. Removal checklist enumerates 10 items. Minor gap: no pseudo-code for the auto-discovery function itself. |
| User-facing behavior described | 40/45 | Significantly improved. Error messages specified ("task type 'test.gen-and-run' is deprecated, regenerate index.json via forge quick-tasks"). Developer experience for missing templates described (init-time fatal with type name). CategoryEval validation fields enumerated (summary/findings/severity). Gap: no description of what the developer sees when running `forge prompt get-by-task-id` for a valid new type — what does success look like beyond "returns valid prompt"? |
| Technical direction clear | 33/35 | embed.FS ReadFile, init-time validation loop, override map for clean-code. File-level changes listed. Template content requirements specified (4 components each). Good. |

**Subtotal: 111/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 32/40 | Major improvement. Rails (CoC for routing/table names), Spring Boot (auto-configuration), ASP.NET Core (controller discovery), GitHub Actions (reusable workflow discovery), Airflow (DAG auto-discovery) all cited. Deduction: no specific reference to how these systems handle the init-time validation pattern the proposal relies on, and no link to published CoC pattern documentation (e.g., Martin Fowler's articles). |
| At least 3 meaningful alternatives | 24/30 | Four alternatives in comparison table: manual map completion, init-time validation + hand-written map, code-gen from types.go, auto-discovery + init-time validation. "Manual map completion" is genuine (not a straw man). "Code-gen from types.go" is a real alternative. Good. Deduction: the comparison table doesn't include "do nothing" explicitly (it was removed), but the absence is fine since P0 precludes it. |
| Honest trade-off comparison | 22/25 | Honest about "变更量较大" and "运行时文件查找". New "等价故障模式" section is intellectually honest about isomorphic failure. Deduction: the "故障检测时机优于原始手写 map" claim is partially valid but should note that a simple init-time validation on the existing hand-written map would achieve the same detection timing without auto-discovery. |
| Chosen approach justified | 20/25 | Auto-discovery + init-time validation justified as eliminating the entire bug class. The "O(N) to O(1)" framing is about maintenance burden, not runtime. Reasonable justification. Deduction: the proposal does not justify why auto-discovery is better than simply adding init-time validation to the existing hand-written map (which would be strictly less change with the same safety improvement). |

**Subtotal: 98/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 36/40 | Happy path (auto-discovery works), edge cases (clean-code override, re-index idempotency, mixed features), error scenarios (missing template → init-time fatal, old index.json → migration error). CategoryEval submit validation analyzed. findFirstTestTaskIdx quick-mode scenario covered. Deduction: no explicit scenario for "what happens when a type name contains characters that produce the same filename as another type" (collision scenario), though the proposal notes all existing types use only `.` and letters. |
| Non-functional requirements | 35/40 | Three NFRs: backward compatibility (migration-aware errors), deployment-time safety (init-time validation), submit validation semantic correctness (CategoryEval with review fields). All three are specific and verifiable. Deduction: no performance NFR for init-time FS reads (trivial but not stated). No NFR for test coverage requirements on new CategoryEval validation logic. |
| Constraints & dependencies | 25/30 | Upstream proposals complete, changes span Go + plugin data + docs. Good. Missing: no mention of Go version constraint for embed.FS features used, no mention of test infrastructure needed for verification. |

**Subtotal: 96/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 28/40 | Auto-discovery via naming convention is standard CoC. The combination with init-time validation is sound but not novel — it's the standard production-ready CoC pattern. The proposal applies it competently. The "O(N) to O(1)" maintenance framing is a useful re-characterization but not a creative leap. |
| Cross-domain inspiration | 20/35 | Rails, Spring Boot, Airflow, GitHub Actions all cited as inspiration. Good attribution. Deduction: no inspiration from domains outside web/CI frameworks (e.g., how plugin architectures like VSCode extensions, ESLint rules, or Babel plugins handle type-to-resource mapping). |
| Simplicity of insight | 20/25 | The `strings.ReplaceAll` one-liner is genuinely elegant. The insight that the existing naming convention is already deterministic is clean. Deduction: the proposal could be simpler — init-time validation on the existing map would achieve the same safety without auto-discovery, and the clean-code override complicates the "zero exception" claim. |

**Subtotal: 68/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | `strings.ReplaceAll` + embed.FS ReadFile is trivial. CategoryEval is a straightforward enum addition with new validation branch. Removal is mechanical with the 10-item checklist. Risk: 4 new prompt templates require domain expertise to write correctly, and the proposal only provides structural guidance, not content guidance. |
| Resource & timeline | 26/30 | "2 coding tasks + doc tasks" is realistic. The removal checklist and specific file enumeration make estimation credible. Deduction: writing 4 semantically correct execution-phase prompt templates may require more iteration than estimated. |
| Dependency readiness | 28/30 | Upstream proposals complete. No external dependencies. The 10-item checklist makes it hard to miss files. Good. |

**Subtotal: 90/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | Three priority tiers with specific file-level changes. isTestTaskID direction now committed (not "or"). CategoryEval direction committed. Removal checklist provides 10 concrete items. Good. |
| Out-of-scope explicitly listed | 23/25 | Four items out of scope with rationale. "旧 index.json 迁移工具" is correctly out of scope with the migration-aware error as the mitigation. Good. |
| Scope is bounded | 22/25 | The 10-item removal checklist bounds P2. P0 and P1 are well-scoped. Deduction: "更新引用废弃类型的测试文件（8+ 文件，需逐一 grep 确认）" is still slightly open-ended — the proposal could enumerate these files to fully close scope. |

**Subtotal: 73/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 27/30 | 8 risks identified, up from 4 in iteration 1. Now covers: auto-discovery edge cases, old index.json migration, category modification side effects, prompt template accuracy, clean-code rename external references, findFirstTestTaskIdx regression, incomplete gen-and-run removal, eval/CategoryEval field mismatch. Deduction: no risk for "new CategoryEval constant causes dispatch logic errors in untested code paths" (the freeform [low] finding about adding compile-time or init-time verification of category coverage). |
| Likelihood + impact rated | 26/30 | Ratings are honest. "自动发现遗漏 edge case" at L/H is justified by the naming convention stability. "findFirstTestTaskIdx quick-mode 回退失效" at H/M is a refreshingly honest high-likelihood rating. Deduction: "category 修改影响质量门控现有行为" rated M/M — the impact could be H if CategoryEval validation is wrong (eval tasks would be unsubmitable), similar to the original P1. |
| Mitigations are actionable | 24/30 | Most mitigations are now specific: grep commands for clean-code rename, specific error message for old index.json, unit tests for CategoryEval. Deduction: "运行全量测试验证" for category modification is still generic. "约定已稳定，所有现有类型名仅含 `.` 和字母" for auto-discovery edge cases is an assertion, not a mitigation — what happens when someone adds a type with a `-` in it? |

**Subtotal: 77/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 48/55 | 12 criteria, up from 6. Most use concrete grep commands or specific behavior checks. "forge prompt get-by-task-id" returns valid prompt is testable. "forge submit-task 对 eval 任务接受 review 类字段" is testable. "init-time 校验：缺失模板文件时 CLI 启动失败并报告缺失的类型名" is testable. Deduction: "mixed feature re-index 幂等：T-review-doc 依赖不丢失" is testable but could be more specific (what does "不丢失" mean in terms of the index.json content?). "所有现有测试通过" is necessary but not sufficient — it doesn't verify the new behaviors. |
| Coverage is complete | 22/25 | P0: template mapping covered (criterion 1-2). P1: CategoryEval covered (4-5), clean-code rename covered (7), re-index covered (9). P2: gen-and-run removal covered (6), validate_index migration covered (8), findFirstTestTaskIdx covered (10). Deduction: no criterion explicitly verifying record-format-doc.md or record-format-test.md content correctness (though covered implicitly by grep criteria). |

**Subtotal: 70/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 33/35 | Auto-discovery fixes P0 root cause. CategoryEval fixes P1 misclassification. Removal checklist addresses P2. All 14 gaps have corresponding solution items. Strong alignment. Deduction: the "O(N) to O(1)" framing in Innovation Highlights is about map entries, but the Scope section still lists the clean-code override as a special case — this is O(1) with an exception, not purely O(1). |
| Scope <-> Solution <-> Success Criteria aligned | 26/30 | Significantly improved. Most scope items have corresponding success criteria. isTestTaskID direction committed. findFirstTestTaskIdx covered. validate_index covered. Deduction: scope item "补充 category_test.go 测试用例" has no direct success criterion (covered only by "所有现有测试通过"). Scope item "task-lifecycle.md: 更新系统类型列表" has no direct success criterion. |
| Requirements <-> Solution coherent | 22/25 | NFRs map cleanly: backward compatibility → migration error, deployment-time safety → init-time validation, submit validation → CategoryEval. No contradictions. Deduction: the requirement "提交验证语义正确性" (NFR) and the solution "CategoryEval 验证分支" are well-aligned, but the success criterion only checks that submit-task "accepts" review fields — it doesn't verify that it "rejects" test-only fields for eval tasks (negative test). |

**Subtotal: 81/90**

---

## Phase 3: Blindspot Hunt

### [beyond-rubric] Issues

1. **Auto-discovery filename collision risk unaddressed**: The convention `typeName "." → fileName "-"` means two type names like `foo.bar-baz` and `foo.bar.baz` would map to the same file `foo-bar-baz.md`. The proposal states "所有现有类型名仅含 `.` 和字母" but does not prevent future types from containing `-`. No guard or validation for this is proposed. The init-time check verifies file existence but not uniqueness of the mapping.

2. **CategoryEval dispatch side effects**: The proposal creates a new `CategoryEval` constant and adds a validation branch. However, it does not audit all code paths that switch on category (e.g., any dispatch logic, logging, metrics, dashboard rendering). If any code uses a default/else branch, CategoryEval tasks could fall into unintended behavior paths. The freeform [low] finding about compile-time or init-time verification of category coverage is relevant here.

3. **Removal checklist assumes no external consumers**: The 10-item removal checklist covers internal code but does not address whether any external tools, scripts, or CI pipelines reference `test.gen-and-run` or `T-quick-gen-and-run` as string literals. The clean-code rename risk addresses this for `clean-code.md` but not for the removed type.

4. **Init-time validation timing**: The proposal mentions "init 阶段" or "CLI 入口" for validation but does not specify which. If it runs in `init()`, it executes even during `go test`, potentially slowing test runs. If it runs at CLI entry, `go test` won't catch missing templates. The choice matters and should be specified.

---

## Score Summary

| Dimension | Score | Max | Delta from Iter 1 |
|-----------|-------|-----|-------------------|
| Problem Definition | 99 | 110 | +1 |
| Solution Clarity | 111 | 120 | +8 |
| Industry Benchmarking | 98 | 120 | +40 |
| Requirements Completeness | 96 | 110 | +11 |
| Solution Creativity | 68 | 100 | +15 |
| Feasibility | 90 | 100 | +2 |
| Scope Definition | 73 | 80 | +5 |
| Risk Assessment | 77 | 90 | +13 |
| Success Criteria | 70 | 80 | +12 |
| Logical Consistency | 81 | 90 | +11 |
| **Total** | **863** | **1000** | **+118** |

---

## ATTACKS

1. [Industry Benchmarking]: Trade-off comparison does not justify auto-discovery over "init-time validation + hand-written map" — quote: "init-time 校验 + 手写 map" is listed as Rejected with con "map 仍需手动维护；O(N) 条目随类型增长", but the same safety improvement (early detection) is achieved with strictly less change. The rejection rationale focuses on maintenance burden, not safety, but safety was the original P0 concern. Must justify why maintenance convenience outweighs the risk of introducing auto-discovery bugs.

2. [Solution Creativity]: The "O(1) 约定" framing is undermined by the clean-code override — quote: "唯一例外：TypeCleanCode = 'code-quality.simplify' 在 prompt.go 映射到 clean-code.md（非约定名），用小 override map 处理". The override map is itself a maintenance surface. The innovation is less clean than claimed. Must either eliminate the override (by renaming the file) or qualify the "O(1)" claim.

3. [Requirements Completeness]: No NFR for test coverage of the new CategoryEval validation logic — quote: "CategoryEval 为新增类别，不修改 CategoryTest 现有行为" in risk table, but no NFR requires testing the new category's validation behavior. Without a test coverage NFR, CategoryEval could ship with untested validation paths. Must add an NFR or success criterion for CategoryEval unit test coverage.

4. [Risk Assessment]: Risk "category 修改影响质量门控现有行为" mitigation is "运行全量测试验证" — this is generic and does not specify what "全量" means. If the test suite has no tests for CategoryEval submission, "全量测试通过" gives false confidence. Must specify that new CategoryEval-specific tests are required as part of the mitigation.

5. [Logical Consistency]: Success criterion "forge submit-task 对 eval 任务接受 review 类字段（summary/findings），不要求测试证据" only verifies the positive case (accepts review fields) but not the negative case (rejects test-only fields for eval tasks). A valid submission with only testsPassed/coverage could still be accepted if the validation is additive rather than exclusive. Must add a criterion verifying eval tasks are rejected when only test fields are provided.

6. [Scope Definition]: Scope item "更新引用废弃类型的测试文件（8+ 文件，需逐一 grep 确认）" is open-ended — the removal checklist partially addresses this with item 10, but "8+" is not a bound. The proposal could enumerate these files to fully close scope. Must enumerate the 8+ files or provide the grep command that would produce the definitive list.

7. [beyond-rubric]: Auto-discovery filename collision risk — the convention `typeName "." → fileName "-"` could produce collisions if future type names contain `-`. The init-time check verifies file existence but not mapping uniqueness. A `test.gen-journeys` and a hypothetical `test.gen-jour-neys` would map to the same file. Should add a validation that all derived filenames are unique across registered types.

8. [beyond-rubric]: Init-time validation timing unspecified — the proposal mentions "init() 或 CLI 入口" but these have different implications. `init()` runs during tests; CLI entry does not. Must choose one and justify.
