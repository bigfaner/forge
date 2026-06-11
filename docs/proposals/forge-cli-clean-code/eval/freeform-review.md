# Freeform Review: Forge CLI Clean Code Proposal

**Reviewer**: Go Codebase Health Engineer
**Date**: 2026-05-24
**Document**: `docs/proposals/forge-cli-clean-code/proposal.md`

---

## Section 1: Background Assessment

This proposal addresses accumulated technical debt in forge-cli, a Go codebase with 92 source files. The audit identified 18 concrete issues across four categories: dead code (4), duplicated logic (5), oversized files (2), and anti-patterns (4). The proposed solution is a four-phase bottom-up refactoring executed in strict sequence: dead code removal, duplicate consolidation, file splitting, and anti-pattern repair.

The proposal is candid about its nature — it explicitly states "无创新" (no innovation), positioning itself as standard Go community hygiene practice. The self-described scope is pure refactoring: zero behavioral change, zero new dependencies, zero new features. This framing is appropriate and sets clear boundaries.

The core assumption is that phase ordering by increasing blast radius — from safe deletions (Phase 1) through structural reorganization (Phase 3) to behavioral refactoring (Phase 4) — minimizes risk at each step. Each phase is gated on a full test suite pass. The codebase is at v3.0.0-rc.19, so the urgency argument (pre-release cleanup window) is sound.

However, after cross-referencing the proposal claims against the actual codebase, several assertions require scrutiny, and the phase ordering contains a subtle but significant inversion that could compound errors across phases.

---

## Section 2: Key Risk Identification

### Phase Ordering Inversion

`问题：` The proposal states the phase order as "死代码消除 → 重复逻辑合并 → 超大文件拆分 → 反模式修复" and describes it as "按自底向上顺序执行四阶段". However, Phase 2's re-export layer cleanup and Phase 3's file splitting are in the wrong relative order. The re-export layer in `errors.go` and `output.go` currently serves as the import target for sibling packages within `internal/cmd/`. If Phase 2 removes these re-exports first, every file that currently uses `Exit()`, `PrintField()`, `CalcSlugColWidth()`, `Debugf()`, etc. must be updated to use `base.Exit()`, `base.PrintField()`, etc. Then Phase 3 splits `forensic.go` and `worktree.go`, requiring yet another round of import path adjustments. This double-churn on the same call sites is avoidable by reversing Phase 2 and Phase 3: split files first (mechanical, same package), then clean up re-exports (requires import path changes across packages).

`风险：` The proposal says "清理重导出层（`errors.go`、`output.go`）" in Phase 2. But the actual usage data shows that `lesson.go`, `proposal.go`, `research.go`, `version.go`, and `quality_gate.go` — all within `internal/cmd/` — call `Exit()`, `PrintField()`, `CalcSlugColWidth()`, `Debugf()` through the re-export layer, not through `base.` directly. That is a minimum of 5 files with dozens of call sites that must all be updated. Meanwhile, `forensic/forensic.go` (already in its own sub-package) calls `base.Exit()` and `base.PrintField()` directly, as does `feature/feature.go`. The codebase is already inconsistent. Removing the re-export layer means either (a) updating all `internal/cmd/*.go` files to use `base.` directly, or (b) keeping the re-export. Option (a) is the right long-term choice but is not a simple cleanup — it is a cross-cutting import refactor that belongs late in the sequence, not in Phase 2.

### validateRecordData and os.Exit — Mischaracterized Scope

`问题：` The proposal lists "修复 `validateRecordData()` 中的 `os.Exit` → 返回 error" in Phase 4. After inspecting the code, `validateRecordData()` at `internal/cmd/task/submit.go:314` does not call `os.Exit()` directly — it calls `base.Exit()`, which is a wrapper that calls `os.Exit()`. More critically, this function is called from `runSubmit()` at line 142, which is a cobra `RunE` function. Changing `validateRecordData` to return an `error` instead of calling `base.Exit()` means its caller (`runSubmit`) must handle that error, and the behavioral contract changes: currently `base.Exit()` terminates immediately; returning an error means the error will propagate up through cobra's error handling, which may produce different output formatting. The proposal acknowledges this risk in the Key Risks table ("os.Exit 移除改变错误处理流程") but underestimates it: the test file `submit_test.go` calls `validateRecordData` directly in 30+ test cases, all of which assume the function either succeeds or calls `os.Exit`. Converting to error-return will require restructuring every one of those test cases.

`风险：` The proposal's mitigation states "只改内部函数，顶层 RunE 保持 Exit 行为". But `validateRecordData` is already an internal function — the issue is that its 30+ direct test callers depend on its current exit-based semantics. This is not a low-impact change; it is a test architecture refactor disguised as a production code fix.

### testbridge.go — File Name vs. Build Tag Discrepancy

`问题：` The proposal states "testbridge.go 导出 40+ 内部符号供测试使用，但文件在正式构建中". After inspecting `internal/cmd/task/testbridge.go`, the file comment explains: "Named export_for_test.go (not export_test.go) so it's compiled into the regular package (not just the test binary)." This is intentional — the file exists specifically because `cmd/integration_test.go` (in a different package) needs to import these symbols. If this file is cleaned up or its exports are removed, `integration_test.go` loses access to all these symbols, and there is no `_test.go` file in the `task` package that could provide them (since `integration_test.go` is in package `cmd`, not package `task`). The proposal acknowledges "保持导出接口不变，仅重新组织" but then also says "清理 testbridge.go 模式 + 迁移 `getTaskPhase()` 到 `pkg/task/`" — these two goals are in tension. "Cleaning up" the pattern while "keeping the export interface unchanged" is a contradiction if the cleanup involves removing the bridge file or reorganizing its exports.

`风险：` The testbridge exports 37 symbols (not "40+" as stated — a minor factual inaccuracy). Any reorganization that changes the exported names or moves the underlying functions to different packages will break `integration_test.go` at 25+ call sites (e.g., `taskpkg.ExportRunSubmit`, `taskpkg.ExportValidateRecordData`, etc.).

### Frontmatter Parsing Duplication — Already Partially Consolidated

`问题：` The proposal states "YAML frontmatter 解析在 3 个文件中独立实现". The actual codebase shows a more nuanced picture. There is already a canonical implementation at `pkg/task/frontmatter.go:ParseFrontmatter()` that parses into a typed `FrontmatterData` struct. Additionally, there is `pkg/infocmd/infocmd.go:ParseFrontmatter()` with a different signature (generic target), `internal/cmd/feature/feature.go:parseYAMLFrontmatter()` (private, different signature), and `internal/cmd/worktree/worktree.go:parseFrontmatter()` (private, different signature). These are not three identical implementations — they have different signatures, different return types, and different error handling. Merging them into `pkg/frontmatter/` as proposed is not a simple deduplication; it requires designing a unified API that satisfies all four callers, which may mean a generic function that accepts `any` or a typed struct. This is design work, not mechanical refactoring, and it blurs the line between Phase 2 (duplicate elimination) and feature work.

`风险：` Consolidating four different frontmatter parsing APIs into one package introduces the risk of over-generalization. The `pkg/task.ParseFrontmatter()` returns `(FrontmatterData, []byte, error)` — a task-specific type. The `infocmd.ParseFrontmatter()` targets `any`. If the unified API uses generics, it adds complexity. If it uses `any`, it loses type safety for the task-specific use case.

### SetFeature() Migration — Underestimated Call Site Count

`问题：` The proposal says "完成 `SetFeature()` 迁移并删除废弃函数" in Phase 1 (dead code elimination). However, `SetFeature()` is called from at least 7 distinct files: `feature.go`, `integration_test.go` (5 call sites), `prompt_test.go`, and `testing_helpers_test.go`. The deprecated function is a one-line wrapper: `return EnsureFeatureDir(projectRoot, featureSlug)`. Migration is trivial (search-and-replace `SetFeature` → `EnsureFeatureDir`), but this is not "dead code" — it is deprecated-but-active code. Classifying it as Phase 1 dead code elimination is semantically wrong. If the migration is incomplete or a call site is missed, the compiler will catch it (the function still exists), but if the function is deleted before all call sites are migrated, compilation fails.

`风险：` Phase 1 should handle only truly dead code (zero call sites). `SetFeature()` has 7+ call sites and should be in Phase 2 or a standalone mini-phase. Mixing it into Phase 1 dilutes the safety guarantee that Phase 1 changes are trivially safe (delete unused code).

### quality_gate.go os.Exit — Missing from Scope

`问题：` The proposal identifies `validateRecordData()` as having an `os.Exit` anti-pattern, but `quality_gate.go` contains two additional `os.Exit(0)` calls at lines 140 and 280 that are not mentioned. Line 140 (`os.Exit(0)` for docs-only features) and line 280 (after quality gate completion) are in the `runQualityGate` function — a cobra `RunE` function that should return an error, not call `os.Exit` directly. These are the same class of anti-pattern as the `validateRecordData` issue but are not listed in the proposal's evidence or scope.

`风险：` The proposal's success criterion states "0 处非顶层函数中的 `os.Exit` 调用". But `runQualityGate` is a cobra `RunE` handler — arguably a "top-level" function in the CLI sense. The distinction between "top-level" and "internal" is not precisely defined, which creates ambiguity in what the success criterion actually measures. The two `os.Exit(0)` calls in `quality_gate.go` should either be explicitly scoped in or explicitly out of scope.

### Dead Code Claim Verification

`问题：` The proposal claims "`run.go` 中 `GetVersion()`、`GetName()`、`IsTestMode()` 无任何调用者". After verification, these functions are in `cmd/forge/run.go` (package `main`). Since they are in the `main` package, they cannot be imported by other packages, and grep confirms no usage within the `main` package itself. This claim is verified and safe.

However, the proposal also claims "`cmd/output.go` 重复定义了 `base.Debugf`，且该重导出未被使用". The file `internal/cmd/output.go` defines a `Debugf` function that differs from `base.Debugf` — it is inlined with a local `fmt.Fprintf` implementation rather than calling through to `base`. It is used extensively in `quality_gate.go` (10+ call sites). The claim that this re-export is "unused" is incorrect.

`风险：` Deleting the `Debugf` from `output.go` as "unused dead code" would break `quality_gate.go` and any other file calling `cmd.Debugf()`. This is a factual error in the evidence that could lead to a destructive action if not caught during implementation.

### mapXxxToSlugLens — Generic Refactoring Scope

`问题：` The proposal suggests "用泛型替代 4 个 `mapXxxToSlugLens` 函数". I found only 3 such functions: `mapFeaturesToSlugLens` (in `feature/feature.go`), `mapProposalsToSlugLens` (in `proposal.go`), and `mapReportsToSlugLens` (in `research.go`). There is also `mapLessonsToNameLens` in `lesson.go` (mentioned indirectly). Each of these functions extracts a single integer field from a different struct type. A generic version would look like `func mapToLens[T any](items []T, extract func(T) int) []int`. This is a valid refactoring, but the proposal overcounts the duplication (claims 4, actual is 3 with 1 variant). More importantly, these functions are 3 lines each. The reduction in total line count is minimal (save ~9 lines, add ~5 for the generic + 4 wrapper calls). This is low-value refactoring that adds cognitive complexity for negligible gain.

`风险：` Introducing generics for a 3-line helper that is used once per file adds an abstraction layer without meaningful benefit. The original code is immediately readable; the generic version requires understanding Go type constraints. For a codebase in pre-release stabilization, this is the kind of change that increases review burden without proportional value.

---

## Section 3: Improvement Suggestions

`建议：` **Reorder phases to put file splitting before re-export cleanup.** The proposal's Phase 2 includes "清理重导出层（`errors.go`、`output.go`）" which requires changing import paths across multiple files. Phase 3 splits `forensic.go` and `worktree.go` — operations that stay within the same package and do not change imports. Reversing these two phases would mean: (1) split files mechanically within packages, (2) then update all import paths in a single pass. This avoids double-touching the same files and reduces the window where imports are in an inconsistent state.

`建议：` **Move SetFeature() migration out of Phase 1.** As identified in the risk analysis, `SetFeature()` has 7+ active call sites. Phase 1 should contain only truly dead code (zero callers): `GetVersion()`, `GetName()`, `IsTestMode()` in `run.go`. Moving `SetFeature()` migration to Phase 2 alongside other duplicate elimination work keeps Phase 1 trivially safe and reviewable.

`建议：` **Correct the evidence before proceeding.** The claim that `cmd/output.go`'s `Debugf` is "未被使用" (unused) is factually wrong — it is used in `quality_gate.go`. The claim of "4 个 mapXxxToSlugLens" is overcounted (3 functions with 1 variant). The testbridge exports 37 symbols, not "40+". These factual inaccuracies, while minor individually, undermine confidence in the audit's thoroughness. I recommend a re-verification pass over all 18 evidence items before task breakdown.

`建议：` **Define "top-level" precisely for the os.Exit success criterion.** The proposal says "0 处非顶层函数中的 `os.Exit` 调用" but does not define "顶层". In a cobra CLI, `RunE` functions are the effective top-level handlers. I recommend defining it as: "any function that is not a cobra `RunE` handler or the `main()` function". Under this definition, `quality_gate.go`'s two `os.Exit(0)` calls (in `runQualityGate`, a `RunE` handler) would be explicitly in-scope or explicitly out-of-scope, but not ambiguous. Add them to the scope or document the exclusion.

`建议：` **Downgrade the mapXxxToSlugLens generic refactoring to optional/nice-to-have.** As discussed, the line savings are minimal (~4 lines net), the functions are trivially readable in their current form, and Go generics add cognitive overhead for future readers. For a pre-release stabilization pass, this change does not carry its weight. If included, it should be the last item in Phase 2 so it can be dropped without affecting other work.

`建议：` **Add a testbridge.go impact analysis as a Phase 4 precondition.** Before modifying `testbridge.go`, produce a mapping of every exported symbol to its external callers (primarily `integration_test.go`). The proposal's mitigation of "保持导出接口不变，仅重新组织" is vague — specify what "reorganize" means. Does it mean moving functions to `pkg/task/` and keeping the bridge as thin aliases? Or removing the bridge entirely and moving tests? These are two very different approaches with very different blast radii. Pick one and document it before implementation begins.

`建议：` **For the validateRecordData os.Exit refactor, treat test migration as a first-class work item.** The current test file `submit_test.go` has 30+ direct calls to `validateRecordData` that assume exit-based semantics. Converting the function to return errors means every test case needs to handle the returned error instead of expecting termination. This is not a side effect of the production code change — it is co-equal work. The proposal should acknowledge this explicitly in Phase 4 scope: "Refactor `validateRecordData()` to return error + migrate 30+ test cases in `submit_test.go`".

`建议：` **For frontmatter consolidation, adopt the existing `pkg/task.ParseFrontmatter()` as canonical and deprecate others incrementally.** Rather than creating a new `pkg/frontmatter/` package (which adds a new package to the module), extend `pkg/task.ParseFrontmatter()` with a generic variant if needed, or have callers convert from the typed result. This avoids introducing a new package and keeps the dependency graph simple. The callers in `feature.go` and `worktree.go` parse into different struct types — they can use `yaml.Unmarshal` directly on the extracted YAML bytes, which `pkg/task.ParseFrontmatter()` already separates out as its second return value.
