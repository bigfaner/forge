# Freeform Expert Review: Pipeline Integration Stitch (v2)

**Reviewer**: Go Pipeline Integration & Type System Engineer
**Date**: 2026-05-24
**Document**: `docs/proposals/pipeline-integration-stitch/proposal.md`

---

## Section 1: Background Assessment

### Problem Claim

The proposal identifies a set of integration gaps left by the `auto-gen-journeys-contracts` proposal. These gaps are stratified into three severity tiers: P0 (4 missing prompt templates causing guaranteed pipeline failure in `Synthesize()`), P1 (eval type misclassification into `CategoryCoding`, stale quick-mode matching in `findFirstTestTaskIdx`, dependency injection order coupling in `build.go`), and P2 (comprehensive dead-code removal of `test.gen-and-run` across production code, tests, and documentation). The core thesis is that the type registration system is complete but the execution-layer scaffolding -- prompt templates, category dispatch, record rendering, and validation -- was not updated in lockstep.

### Core Technical Approach

Three workstreams executed in dependency order: (1) create 4 prompt template files under `prompt/data/` to resolve P0, (2) add `CategoryEval` with full submit-validation and record-rendering branches plus dependency injection consolidation to resolve P1, and (3) mechanically remove all `test.gen-and-run` references from production code, tests, and active documentation to resolve P2. The auto-discovery mechanism from the prior task (naming convention `strings.ReplaceAll(typeName, ".", "-") + ".md"`) is assumed to handle template-to-type mapping without manual updates.

### Assumptions

1. The naming convention auto-discovery (Task 1) has already eliminated the need for hand-written type-to-filename maps in `prompt.go`.
2. The 4 new prompt templates are execution-phase prompts (agent instructions at runtime), distinct from the autogen planning-phase templates already present in `task/data/`.
3. `eval.journey` and `eval.contract` are semantically review/assessment tasks (scoring rubrics, spawning subagents), not test-generation tasks, and therefore deserve a dedicated category.
4. Removing `test.gen-and-run` is safe because the staged pipeline has fully replaced it; old index.json files can be handled by a migration error message.

---

## Section 2: Key Risk Identification

### P0 Confirmation: Prompt Templates Are Indeed Missing

I verified that `forge-cli/pkg/prompt/data/` contains 19 template files. The files `test-gen-journeys.md`, `test-gen-contracts.md`, `eval-journey.md`, and `eval-contract.md` are absent. Meanwhile, the corresponding types (`TypeTestGenJourneys`, `TypeTestGenContracts`, `TypeEvalJourney`, `TypeEvalContract`) are registered in `ValidTypes` (types.go:111-118). When `Synthesize()` calls `templatePath(t.Type)` and then `templateFS.ReadFile()`, these 4 types will fail at runtime. The P0 claim is accurate.

### Template Content Quality Gap

风险：The proposal states: "每个模板包含：任务上下文说明、输入格式（从 index.json task configuration 读取）、期望输出格式、质量标准"

Looking at the existing prompt templates in `prompt/data/`, the execution-phase pattern is well-established: `test-gen-scripts.md`, `test-run.md`, `test-verify-regression.md` all follow a consistent structure with `TASK_ID`, `TASK_FILE`, `SCOPE` placeholders, a task-constraints section, and a numbered workflow. However, the proposal does not specify the actual content of the 4 new templates. For `test-gen-journeys.md` and `test-gen-contracts.md`, the agent needs to invoke skills (`forge:gen-journeys`, `forge:gen-contracts`) -- the proposal should specify which skills and what the Record Fields section should list. For `eval-journey.md` and `eval-contract.md`, the situation is different: these are `MainSession: true` tasks that spawn subagents. Their prompt templates must instruct the agent to use the eval rubric system, which is architecturally distinct from the standard skill-invocation pattern.

问题：The proposal says: "参考现有执行阶段模板结构（code-quality-simplify.md）" in the risk mitigation table. But `code-quality-simplify.md` is a coding-category task with coverage injection and quality gate workflow -- structurally irrelevant to eval tasks. The eval templates need a fundamentally different pattern: load rubric, score artifact, iterate until threshold, write eval report. The proposal should specify the eval-specific template structure explicitly.

### CategoryEval: Correct but Incomplete Specification

问题：The proposal correctly identifies that `eval.*` types fall into `CategoryCoding` via the default branch in `CategoryForType()` (category.go:30-31). I confirmed this: `CategoryForType("eval.journey")` matches no prefix case and returns `CategoryCoding`. The proposal's solution -- adding `CategoryEval = "eval"` and a `strings.HasPrefix(typ, "eval.")` branch -- is architecturally sound.

However, the proposal lists specific implementation items: "submit.go: validateRecordData 为 CategoryEval 添加验证分支（接受 review 字段 summary/findings/severity，拒绝纯测试字段）" and "types.go: RecordData 结构添加 eval 特有字段（evalScore、evalFindings、evalSeverity、evalPassed）". But looking at `RecordData` (types.go:280-312), there are no `evalScore`, `evalFindings`, `evalSeverity`, or `evalPassed` fields. The proposal introduces 4 new struct fields without specifying their JSON tag names, omitempty behavior, or how they interact with existing fields like `Summary` and `Notes`.

风险：Adding `evalScore`, `evalFindings`, `evalSeverity`, `evalPassed` as new `RecordData` fields creates a field-proliferation pattern. Currently `RecordData` uses category-agnostic fields with category-specific interpretation (e.g., `CasesGenerated` is only meaningful for `CategoryTest`). Adding 4 eval-specific fields to the same flat struct follows this pattern but makes `RecordData` increasingly incoherent. The `RecordTemplateData` struct in `record.go` (lines 15-52) would need corresponding formatting fields, and a new `record-eval.md` Go template would need to be created.

### RenderRecord Dispatch Gap

问题：The proposal states: "record.go: 新增 record-eval.md Go 模板 + RenderEvalRecord 函数 + RenderRecord switch 添加 CategoryEval case"

Looking at `RenderRecord()` (record.go:234-247), the switch dispatches on `CategoryForType(t.Type)`. Currently there is no `CategoryEval` case. Without this case, eval tasks would fall into the `default` branch and use `RenderCodingRecord`, which renders coding-specific fields like `TestsPassed`, `TestsFailed`, and `Coverage`. This is semantically wrong for eval tasks. The proposal correctly identifies this, but the scope entry should also mention updating `RecordTemplateData` to include eval-specific formatted fields and `NewRecordTemplateData` to populate them.

### findFirstTestTaskIdx: Verified Stale Match

问题：The proposal states: "findFirstTestTaskIdx quick-mode 分支匹配废弃类型（build.go:492-494）"

I confirmed this. `findFirstTestTaskIdx()` (build.go:485-503) has a quick-mode branch that matches `T-quick-gen-and-run*`. Looking at `GetQuickTestTasks()` (autogen.go:217-306), quick mode now generates `T-test-gen-journeys-<type>` tasks (line 229-231), not `T-quick-gen-and-run*`. The quick-mode branch in `findFirstTestTaskIdx` will never match. The function currently relies on the `return 0` fallback (line 499-500), which works only because `T-test-gen-journeys-<type>` tasks happen to be first in the generated task list. This is fragile: if task ordering changes in `GetQuickTestTasks`, the fallback will return the wrong index.

风险：The proposal's fix -- "更新为匹配 T-test-gen-journeys 前缀" -- is correct but incomplete. Looking at `GetQuickTestTasks`, the first tasks are `T-test-gen-journeys-<type>`, which already use the prefix `T-test-gen-journeys`. The `findFirstTestTaskIdx` function should use `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")` for quick mode, consistent with how `ResolveFirstTestDep` already uses `findTaskIndexByPrefixOrPanic(tasks, "T-test-gen-journeys")` (autogen.go:678, 694). The proposal should ensure both functions use the same discovery mechanism.

### Dependency Injection Consolidation

问题：The proposal states: "build.go: 将 ResolveFirstTestDep + T-review-doc prepend 合并为单步操作"

Looking at the current code in `BuildIndex()` (build.go:326-347), the two-step process is: (1) `ResolveFirstTestDep(testTasks, index.TasksMap(), mode)` sets the first test task's dependencies, then (2) `testTasks[firstTestIdx].Dependencies = append([]string{"T-review-doc"}, testTasks[firstTestIdx].Dependencies...)` prepends T-review-doc. The proposal correctly identifies that this ordering is hard-coded and fragile. The consolidation into a single function is the right approach.

However, the current code has a subtle correctness guarantee: `ResolveFirstTestDep` (autogen.go:658-702) uses `findTaskIndexByPrefixOrPanic(tasks, "T-test-gen-journeys")` which finds the first task matching that prefix. After T-review-doc is prepended, the dependency chain becomes `[T-review-doc, original-dep]`. The merged function must preserve this exact ordering -- T-review-doc must come first because it depends on the last business task, and the original dep (gate or T-clean-code) is a transitive dependency through T-review-doc.

风险：The merged function must handle the case where `needsEval` is false (no T-review-doc). The current code guards with `if needsEval && firstTestIdx >= 0` (build.go:335). The merged function must accept a `needsEval` parameter or derive it from context. The proposal does not specify the function signature.

### gen-and-run Dead Code: Exhaustiveness Concern

风险：The proposal states: "grep -r 'gen-and-run\|quick-gen-and-run\|T-quick-gen' forge-cli/ plugins/ 返回零结果" as a success criterion.

I found 20 files in `forge-cli/` containing these patterns. The proposal identifies key files (types.go, infer.go, prompt.go, validate_index.go, build.go) but the grep results reveal a wider surface area: `autogen.go`, `autogen_test.go`, `stage_gates_test.go`, `quality_gate_test.go`, `output_test.go`, `list_test.go`, and `docs/OVERVIEW.md`. The proposal says "14 个测试文件中 ~95 处引用" but does not enumerate them. The risk of partial removal is real: if `infer.go:32-33` (the `T-quick-gen-and-run` case) is removed but `types.go:55` (`TypeTestGenAndRun`) is not, the code will not compile because the case references a deleted constant.

问题：The proposal should provide a complete file-by-file checklist with line numbers for the gen-and-run removal. The current "5 files ~15 places" for production code undercounts: `autogen.go` is not listed (it does not directly reference gen-and-run in current code, but tests may), and `prompt.go:297,304` references `T-quick-gen-and-run` in `genScriptBases` which must be removed.

### validate_index.go Migration Error

问题：The proposal states: "validate_index.go 对引用 test.gen-and-run 的旧 index.json 返回迁移指引错误信息"

Looking at `validate_index.go:138-139`, the current validation uses `!task.ValidTypes[t.Type]` which produces the generic error "invalid type 'test.gen-and-run'". The proposal wants a migration-aware error like "Task type 'test.gen-and-run' has been removed. Run `forge task build-index` to regenerate with the new pipeline." This requires adding a specific check before the generic `ValidTypes` check. The proposal correctly identifies this need but does not specify where in the validation order this check should appear.

Additionally, `validate_index.go:224-226` has a `T-quick-gen-and-run-` prefix check in `validateFilesExist` that checks for unresolved template placeholders. This dead-code branch must also be removed. The proposal mentions this location but should confirm that removing it does not affect the validation of staged pipeline tasks.

### RecordData Field Design for Eval

风险：The proposal specifies: "types.go: RecordData 结构添加 eval 特有字段（evalScore、evalFindings、evalSeverity、evalPassed）"

The field names use an `eval` prefix (`evalScore`, `evalFindings`, `evalSeverity`, `evalPassed`). But the existing category-specific fields in `RecordData` do NOT use a category prefix: `CasesGenerated` (not `testCasesGenerated`), `ValidationPassed` (not `validationPassed`), `GatePassed` (not `gatePassed`). The naming convention is inconsistent. Either all category-specific fields should be unprefixed (following the existing pattern) or all should be prefixed (for clarity). The proposal introduces a new convention that contradicts the established pattern.

### record-format-test.md Stale References

问题：The proposal states: "record-format-test.md 列出已废弃类型（test.gen-cases/test.eval-cases/test.gen-and-run），缺少新类型"

I verified this. `plugins/forge/skills/submit-task/data/record-format-test.md` line 3 lists: `test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.verify-regression`. Of these, `test.gen-cases` and `test.eval-cases` are not even registered in `ValidTypes` -- they appear to be phantom types that were never implemented. The current valid test types are `test.gen-journeys`, `test.gen-contracts`, `test.gen-scripts`, `test.run`, `test.verify-regression`. The proposal should update this file to reflect only the currently valid types.

### Missing record-format-eval.md

风险：The proposal identifies: "缺少 record-format-eval.md：agent 执行 eval 任务时无 JSON 字段参考"

This is a real gap. When the agent executes an eval task (via the `forge:submit-task` skill), it looks up the record format reference document to know what JSON fields to provide. Without `record-format-eval.md`, the agent has no guidance on eval-specific fields. The proposal correctly identifies this as in-scope, but should specify that this document must be placed at `plugins/forge/skills/submit-task/data/record-format-eval.md` (alongside the other record-format files).

### CategoryForType Default Branch Hazard

问题：The proposal does not address the `default` branch in `CategoryForType` (category.go:30-31): `default: return CategoryCoding`. After adding the `eval.` prefix branch, any future type that lacks a matching prefix will still silently fall into `CategoryCoding`. This is how the eval types originally got misclassified. The root cause -- an overly permissive default -- remains unaddressed.

建议：Consider changing the default to return an error or a `CategoryUnknown` sentinel, with callers handling the unknown case explicitly. This prevents future types from being silently misclassified. At minimum, add a log warning in the default branch so that misclassifications are detectable at runtime.

---

## Section 3: Improvement Suggestions

建议：The proposal should provide explicit content specifications for the 4 new prompt templates. Based on analysis of existing templates (`test-gen-scripts.md`, `test-run.md`), the test-pipeline templates should follow this structure: (1) `TASK_ID/TASK_FILE/SCOPE/PHASE_SUMMARY` header, (2) task constraints section specifying which skills to invoke, (3) numbered workflow steps, (4) record fields section listing what to populate in `forge:submit-task`. For `test-gen-journeys.md`, the skill is `forge:gen-journeys`; for `test-gen-contracts.md`, the skill is `forge:gen-contracts`. For eval templates, the pattern is different: (1) header, (2) rubric loading instruction, (3) scoring workflow (score, check threshold, iterate), (4) eval-specific record fields. The proposal should specify these patterns rather than deferring to a vague "reference existing templates" instruction.

建议：The `RecordData` eval field names should follow the established naming convention. Instead of `evalScore`, `evalFindings`, `evalSeverity`, `evalPassed`, consider unprefixed names like `score` (int), `findings` ([]string), `severity` (string), `passed` (bool). This is consistent with how `CasesGenerated`, `ValidationPassed`, and `GatePassed` are named without category prefixes. The JSON tags can use the same unprefixed names since the record format document (`record-format-eval.md`) will clarify which fields apply to which category.

建议：For the `findFirstTestTaskIdx` fix, use the same discovery mechanism as `ResolveFirstTestDep`: `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")`. This ensures consistency across both functions. The breakdown-mode branch should also use `findTaskIndex(tasks, "T-test-gen-journeys-")` pattern matching or, better yet, the first task in the test pipeline is always `T-test-gen-journeys-*` regardless of mode. The function could be simplified to: find first task matching `T-test-gen-journeys` prefix (works for both breakdown per-type and quick per-type), then fallback to `T-test-gen-contracts` for modes that don't generate per-type journeys.

建议：For the dependency injection consolidation, the merged function signature should be: `ResolveFirstTestDeps(tasks []AutoGenTaskDef, existingTasks map[string]Task, mode string, needsEval bool)`. This encapsulates both `ResolveFirstTestDep` and the T-review-doc prepend logic. When `needsEval` is true, the function prepends `T-review-doc` after setting the base dependency. When false, it behaves identically to the current `ResolveFirstTestDep`. This eliminates the ordering coupling in `BuildIndex()`.

建议：For the gen-and-run removal, create a comprehensive checklist grouped by file, specifying exact changes:

Production code (must compile after each file edit):
- `types.go`: Remove `TypeTestGenAndRun` constant (line 55), remove from `ValidTypes` (line 114), remove from `SystemTypes` (line 135), remove from `TaskTypeRegistry` (line 88)
- `infer.go`: Remove case at lines 32-33
- `prompt.go`: Remove `T-quick-gen-and-run` from `genScriptBases` (line 304)
- `validate_index.go`: Remove `T-quick-gen-and-run-` prefix check at lines 224-226, add migration-aware error
- `build.go`: Update `findFirstTestTaskIdx` quick-mode branch (lines 492-494)
- Delete `prompt/data/test-gen-and-run.md` and `task/data/test-gen-and-run.md`

This ensures compilation can be verified incrementally.

建议：The proposal should add a success criterion verifying that `RenderRecord` dispatches correctly for eval types. Something like: "RenderRecord for a task with type `eval.journey` uses the eval record template (not the coding record template)." This is distinct from the existing criterion about `RenderRecord` using "eval 专用 record 模板" -- it should be explicitly testable by checking that `CategoryForType("eval.journey")` returns `CategoryEval` and that the `RenderRecord` switch includes a `CategoryEval` case.

建议：The proposal should specify that the `clean-code.md` rename (completed in Task 1) means that `prompt/data/` now contains `code-quality-simplify.md` (not `clean-code.md`). This is relevant because the success criteria reference `grep` patterns that should verify the rename was successful. The current proposal mentions the rename in the Innovation Highlights section but does not include a success criterion for it.

建议：For the `record-format-test.md` update, the proposal should specify the exact type list to replace the stale one. The current stale list is: `test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.verify-regression`. The correct list should be: `test.gen-journeys`, `test.gen-contracts`, `test.gen-scripts`, `test.run`, `test.verify-regression`. The proposal should also verify that `test.gen-cases` and `test.eval-cases` are not referenced anywhere else in the codebase.

建议：The proposal's success criterion `grep -r "gen-and-run" forge-cli/ plugins/` returning zero results is a strong correctness guarantee. However, it should be qualified to exclude historical documentation (the proposal already identifies "历史 feature/proposal 文档" as out-of-scope). The grep command should specify file type exclusions (e.g., `--exclude-dir=docs/proposals`) or the proposal should explicitly state that historical docs are excluded from the grep check.
