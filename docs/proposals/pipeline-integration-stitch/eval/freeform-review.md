# Freeform Expert Review: Pipeline Integration Stitch

**Reviewer**: Test Pipeline Architect
**Date**: 2026-05-24
**Document**: `docs/proposals/pipeline-integration-stitch/proposal.md`

---

## Section 1: Background Assessment

### Problem Claim

The proposal identifies 14 integration gaps between two previously completed proposals (`review-doc-pipeline` and `auto-gen-journeys-contracts`). Two are P0 (pipeline guaranteed failure), four are P1 (scenario-specific failure), and eight are P2 (maintenance risk). The core claim is that `prompt.go` and `autogen.go` both use hand-written maps for type-to-template-filename mapping, and new types were added to `autogen.go` but not `prompt.go`, causing runtime failures when the pipeline tries to look up prompt templates for `test.gen-journeys`, `test.gen-contracts`, `eval.journey`, and `eval.contract`.

### Core Technical Approach

Three-pronged: (1) replace hand-written maps in `prompt.go` and `autogen.go` with naming-convention-based auto-discovery (`strings.ReplaceAll(typeName, ".", "-") + ".md"`), (2) fix category classification for `eval.*` types and update stale references in data files, (3) remove all vestiges of the deprecated `test.gen-and-run` type.

### Assumptions

1. The naming convention (`"."` -> `"-"` + `.md`) is universal across all current and future types, with exactly one exception (`code-quality.simplify` -> `clean-code.md` in prompt.go).
2. Renaming `prompt/data/clean-code.md` to `code-quality-simplify.md` has no external consumers.
3. The four new prompt template files can be written correctly by referencing the existing autogen.go templates.
4. Removing `test.gen-and-run` is safe because old index.json files should be migrated, and an explicit error message is sufficient for backward compatibility.

---

## Section 2: Key Risk Identification

### Auto-Discovery Fragility

风险：The proposal states: "自动发现实现简单：`strings.ReplaceAll(typeName, ".", "-") + ".md"`，加上 embed.FS 的 ReadFile 验证文件存在。唯一 override 是 `clean-code.md`。"

The auto-discovery mechanism replaces compile-time map key errors with a runtime file-not-found error. Today, if a developer adds a new type constant but forgets the template file, the failure mode shifts from "map lookup returns false" to "embed.FS.ReadFile returns error." Both are runtime failures in `Synthesize()`. The proposal frames this as eliminating an entire bug class, but the actual failure mode is isomorphic: one map was missed, now one file is missed. The real improvement is that the file naming convention is a single source of truth, but the claim that this is a qualitative shift from O(N) to O(1) is overstated. The developer still must create the template file; the map entry was the thing eliminated, not the template file itself.

问题：The proposal states: "约定已稳定，所有现有类型名仅含 `.` 和字母；重命名 clean-code.md 消除唯一例外". However, `code-quality.simplify` contains a hyphen in its first segment (`code-quality`). The naming convention replaces `.` with `-`, producing `code-quality-simplify.md`. This already works for `autogen.go`. But the proposal also plans to rename `prompt/data/clean-code.md` to `code-quality-simplify.md`. This rename is inside the plugin's embedded FS. If any external tool, script, or documentation references `clean-code.md` by name, this is a silent break. The proposal does not audit for external consumers of this filename.

### eval.* Category Misclassification

风险：The proposal identifies: "`eval.` 前缀不匹配任何规则，落入 default → CategoryCoding" and plans to add `eval.` prefix to `CategoryTest`.

Looking at `category.go`, the current logic is prefix-based (`strings.HasPrefix`). Adding `eval.` to `CategoryTest` means `eval.journey` and `eval.contract` would map to `CategoryTest`. But consider `renderTemplate` in `prompt.go` (line 139): it checks `task.IsTestableType(t.Type)` to decide whether to inject coverage targets. `IsTestableType` returns true only for `coding.*` and `TypeCleanCode`. Since `eval.*` would now be `CategoryTest` (not `CategoryCoding`), the coverage injection would correctly be skipped.

However, the submit-task validation logic uses `CategoryForType` to determine what fields are required. The proposal says eval tasks would be classified as `CategoryTest`, which means they would be expected to provide test-category fields like `casesGenerated`, `scriptsCreated`, etc. But eval tasks evaluate quality (scoring rubrics), they do not generate test cases or scripts. This is a semantic mismatch: `eval.*` tasks are conceptually quality-gate tasks, not test-generation tasks.

问题：The proposal says: "category 影响提交验证和记录渲染" but does not detail what `CategoryTest` requires from eval tasks at submission time. The `record-format-test.md` lists fields like `casesGenerated`, `casesEvaluated`, `scriptsCreated`, `testResults` -- none of which are relevant to `eval.journey` or `eval.contract`. Putting eval tasks into `CategoryTest` may cause submit-task to expect inappropriate fields.

### Mixed Feature Dependency Injection Order

风险：The proposal identifies: "ResolveFirstTestDep 设置依赖后，才 prepend T-review-doc，re-index 幂等性风险"

Looking at `build.go` lines 329-347, `ResolveFirstTestDep` is called first (line 329), which sets the first test task's dependency to the highest gate or T-clean-code. Then T-review-doc is prepended (lines 335-338) as a dependency of the first test task. On re-index, `PreserveRuntimeFields` (lines 316-319) preserves existing deps from the previous index. But `ResolveFirstTestDep` (line 329) re-computes deps from scratch on the in-memory `testTasks` slice, not from the preserved index. So the prepended T-review-doc dep from the previous run is on the *index* task, not on the `testTasks` slice. The `testTasks[firstTestIdx].Dependencies` starts fresh.

This means on re-index: (1) ResolveFirstTestDep sets `testTasks[first].Dependencies = [dep]`, (2) then T-review-doc is prepended: `testTasks[first].Dependencies = ["T-review-doc", dep]`, (3) this is written back to the index. This is idempotent because the testTasks slice is regenerated from scratch each time. The proposal's concern about "re-index idempotency" appears to be a non-issue given the current code structure, because auto-gen tasks are always rebuilt from scratch.

问题：The proposal lists "mixed feature re-index 幂等性：T-review-doc 依赖不丢失" as a success criterion, but the code already handles this correctly by regenerating test tasks fresh on every BuildIndex call. The real risk is not idempotency but the interaction between `ResolveFirstTestDep` and the T-review-doc prepend: if `ResolveFirstTestDep` sets the first test task's dep to `T-clean-code`, and then T-review-doc is prepended, the chain becomes `T-clean-code -> T-review-doc -> first-test-task`. But T-review-doc has no dependency on T-clean-code -- it depends on the last business task. So the first test task would wait for both T-review-doc AND T-clean-code, which is correct. However, if `ResolveFirstTestDep` sets dep to a business gate (not T-clean-code), the chain becomes `gate -> T-review-doc -> first-test-task`, which means the first test task waits for T-review-doc but not the gate directly. T-review-doc depends on the last business task, which may or may not be the gate. This is semantically correct but the proposal does not analyze this chain.

### clean-code.md Rename Side Effects

风险：The proposal says: "重命名 `prompt/data/clean-code.md` -> `code-quality-simplify.md`（统一命名约定）"

This file is embedded via `//go:embed data/*.md` in `prompt.go`. Renaming it affects only the binary. However, if the clean-code template content references its own filename or if any documentation references `clean-code.md`, this is a silent break. The proposal does not audit for such references. Additionally, the proposal's success criterion `grep -r "doc.eval" forge-cli/ plugins/` does not cover `clean-code` references that would need updating after the rename.

### Backward Compatibility with test.gen-and-run

风险：The proposal states: "旧 index.json 应迁移到新 pipeline" and "明确错误提示，引导用户重新生成"

Removing `TypeTestGenAndRun` from `ValidTypes` means that any existing index.json containing tasks with `type: "test.gen-and-run"` will fail `validateTasks` in `validate_index.go` (line 139: `!task.ValidTypes[t.Type]`). This will cause `forge task validate-index` to error. But the proposal does not specify what error message should appear, nor does it mention updating `validate_index.go` to provide a migration-aware error message instead of the generic "invalid type" error.

问题：The proposal lists "已有 index.json 引用 `test.gen-and-run` 时给出明确错误而非静默失败" as a non-functional requirement, but does not include any implementation detail about how this error will be differentiated from a generic "invalid type" error. Without a specific migration error path, users will see "invalid type 'test.gen-and-run'" and may not understand that they need to regenerate the index.

### genScriptBases Dead Code and Test Task ID Extraction

风险：The proposal identifies: "`prompt.go` genScriptBases 包含死代码 `T-quick-gen-and-run`" (item 8).

Looking at `prompt.go` lines 292-295, `genScriptBases` lists `"T-quick-gen-and-run"`. The `extractTestTypeArg` function uses this to extract `--type` arguments from task IDs. If this entry is removed without also updating the test task generation in Quick mode, any Quick-mode tasks that were previously `T-quick-gen-and-run-cli` format would lose their `--type` argument extraction. But since Quick mode now uses `T-test-gen-scripts-<type>` (staged pipeline), this is indeed dead code. The risk is minimal but the proposal should confirm that no other code path generates `T-quick-gen-and-run-*` task IDs.

### isTestTaskID vs IsAutoGenTaskID Coverage Gap

问题：The proposal identifies: "`isTestTaskID` 与 `IsAutoGenTaskID` 覆盖范围不一致" (item 10) but does not propose a specific fix.

Looking at the code: `isTestTaskID` (build.go line 430) checks prefixes `T-test-`, `T-quick-`, `T-specs-`, `T-clean-`, `T-validate-`, `T-eval-`. `IsAutoGenTaskID` (build.go line 508) calls `isTestTaskID` plus checks `T-review-doc` and gate/summary suffixes. But `isTestTaskID` does NOT check `T-review-doc`, while `IsAutoGenTaskID` does. This means `needsTestPipeline` (which uses `IsAutoGenTaskID` to skip) would correctly skip T-review-doc, but `needsReviewDoc` uses `CategoryForType` directly. The proposal flags this inconsistency but does not specify whether the fix is to add `T-review-doc` to `isTestTaskID` or to leave it as-is and update documentation.

### Four New Prompt Templates

风险：The proposal says: "创建 4 个缺失的 prompt 模板文件：test-gen-journeys.md、test-gen-contracts.md、eval-journey.md、eval-contract.md"

These are prompt templates for the `Synthesize()` function in `prompt.go`, which generates the agent prompt at runtime. They are distinct from the autogen templates in `autogen.go` which generate .md task files. The proposal says to "参考 autogen.go 中已有的对应模板结构" but these serve different purposes: autogen templates create task description files, while prompt templates create agent execution instructions. Copying the structure of one to build the other could lead to semantically incorrect prompts. The proposal should specify what each of the four prompt templates should contain and how they differ from the autogen templates.

### Missing Success Criterion for clean-code Rename

问题：The success criteria include: "`grep -r "gen-and-run" forge-cli/ plugins/` 返回零结果" but do not include a corresponding criterion for verifying that `clean-code.md` references (in prompts, tests, or documentation) have been updated to `code-quality-simplify.md`. The rename could leave dangling references that pass compilation but cause confusion.

### infer.go Stale Entry

问题：The proposal lists removing `test.gen-and-run` from `types.go`, `infer.go`, `prompt.go`, `autogen.go` but does not explicitly mention removing the `T-quick-gen-and-run` case from `InferType()` in `infer.go`. Line 32-33: `case id == "T-quick-gen-and-run", typeSuffixedID(id, "T-quick-gen-and-run"): return TypeTestGenAndRun`. Since `TypeTestGenAndRun` constant would be removed, this case would fail to compile. But the proposal should confirm that all references are removed together to avoid partial compilation failures.

### Scope Boundary: eval.* in CategoryTest vs New Category

风险：The proposal proposes adding `eval.` prefix to `CategoryTest`. But eval tasks spawn subagents (MainSession: true in autogen.go lines 109, 123). This means eval tasks run in the main session, unlike most test tasks. Putting them in `CategoryTest` conflates two different execution models. A future refactor that adds category-based session routing could accidentally dispatch eval tasks to a task executor that cannot spawn subagents.

---

## Section 3: Improvement Suggestions

建议：Instead of adding `eval.` to `CategoryTest`, consider creating a dedicated `CategoryEval` constant. This avoids semantic contamination of the test category, enables submit-task to use eval-specific validation fields (e.g., `score`, `rubricResults`), and prevents future confusion when category-based routing logic is added. The eval category is architecturally distinct: it performs quality assessment, not test generation or execution.

建议：For the auto-discovery mechanism, add a compile-time or init-time validation that cross-references all type constants in `types.go` against the embedded template FS. A simple `init()` function that iterates `TaskTypeRegistry` and verifies each type has a corresponding template file (using the naming convention) would catch missing templates at startup rather than at runtime. This would make the auto-discovery genuinely safer than the hand-written map.

建议：For backward compatibility with `test.gen-and-run`, instead of simply removing it from `ValidTypes`, add a dedicated migration error path in `validate_index.go`. When a task has `type: "test.gen-and-run"`, emit a specific message like "Task type 'test.gen-and-run' has been removed. Run `forge task build-index` to regenerate with the new pipeline." This is more helpful than the generic "invalid type" error and directly addresses the non-functional requirement stated in the proposal.

建议：The proposal should add an explicit success criterion for the `clean-code.md` rename: something like `grep -r "clean-code.md" forge-cli/ plugins/` returns zero results (excluding the renamed file itself). This ensures no dangling references survive the rename.

建议：The proposal should explicitly list all files that need the `T-quick-gen-and-run` reference removed, including `infer.go` line 32-33, `prompt.go` line 294, and `validate_index.go` lines 224-226. Creating a checklist of all 8+ test files mentioned in item 12 would help prevent partial cleanup.

建议：For the four new prompt templates, the proposal should specify that these are execution-phase prompts (distinct from the autogen planning-phase templates). Each should follow the existing prompt template pattern: define the agent's goal, provide context placeholders (`{{TASK_ID}}`, `{{TASK_FILE}}`, etc.), and specify the expected output format. The eval templates should explicitly instruct the agent to use the scoring rubric and subagent spawning pattern, since they are MainSession tasks.

建议：The mixed feature dependency injection (build.go lines 329-347) should be refactored to compute the final dependency chain in a single pass rather than the current two-step approach (ResolveFirstTestDep then T-review-doc prepend). A `resolveMixedFeatureDeps` function that takes both T-review-doc and test tasks as input and produces the final dependency graph would be more maintainable and easier to reason about than the current sequential mutation pattern.

建议：The `isTestTaskID` function in `build.go` should be updated to include `T-review-doc` for consistency with `IsAutoGenTaskID`. Currently `IsAutoGenTaskID` handles this by checking `isTestTaskID` first then adding `T-review-doc` as a special case. Adding it to `isTestTaskID` (or renaming the function to `isAutoGenID` and consolidating) would eliminate the coverage gap and reduce the chance of future inconsistencies.

建议：The proposal should explicitly address the `findFirstTestTaskIdx` function in `build.go` (line 485). Currently it checks for `T-eval-journey` (breakdown) and `T-quick-gen-and-run*` (quick). After removing `test.gen-and-run`, the quick mode fallback (line 494) would never match. The function needs to be updated to look for `T-test-gen-journeys-*` or `T-test-gen-contracts` as the first quick-mode test task instead.
