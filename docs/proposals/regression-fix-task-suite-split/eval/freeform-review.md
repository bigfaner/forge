---
created: "2026-05-28"
reviewer: domain-expert/ci-fix-task-scoping
document: proposal.md
---

# Freeform Review: Regression Fix Task Suite Split

## Section 1: Background Assessment

This proposal addresses a concrete production failure: when the quality-gate hook detects regression test failures spanning 20+ failures across 4 test suites, it creates a single `fix` task whose scope overwhelms the LLM agent, causing it to stall until manual intervention. The evidence is documented in `docs/lessons/gotcha-fix-task-broad-scope.md`, which traces the root cause to a granularity mismatch -- the fix task scope is "整次 test 运行" (an entire test run) rather than being bounded to a single failing unit.

The core technical approach is to introduce `addRegressionFixTasks` in `quality_gate.go`, which reuses the existing `extractSourceFiles` function to pull file paths from test output, applies an `isTestFile` naming convention matcher to identify test files, then groups failures by test file -- one fix task per test file. When naming conventions are unrecognized, the system falls back to the current directory-based `groupFilesByDir` behavior. The proposal also removes the `maxFixTasksPerStep` cap (currently set to 3), arguing that with per-file granularity the cap is no longer needed.

The approach explicitly avoids framework-specific parsing (the original v1 used Go `FAIL <package>` parsing) in favor of naming convention detection, claiming language-agnostic coverage for Go (`*_test.go`), Python (`test_*.py`), JS/TS (`*.test.ts`, `*.spec.ts`), Java (`*Test.java`), and Ruby (`*_test.rb`). The proposal positions this as a minimal-change solution that reuses existing infrastructure.

## Section 2: Key Risk Identification

问题：输出行关联算法的准确性边界未被精确界定。

The proposal states in the Scope section: "每个 fix task 只包含该测试文件相关的输出行（从 output 中提取包含该文件路径的行及上下文）". This is the most technically demanding part of the entire proposal, yet it receives no algorithmic specification. The existing `extractSourceFiles` function uses `sourceFileRe` (a regex matching `file.ext:line` or `file.ext:line:col` patterns) to extract file paths from output. But extracting paths and filtering lines by path are fundamentally different operations.

Consider the real test output formats this must handle: Go's `--- FAIL: TestName (0.00s)` headers followed by indented `file_test.go:42: message` lines; Python's pytest `FAILED tests/test_foo.py::TestClass::test_method` format; TypeScript's Playwright or Jest output with stack traces. The "上下文" (context) notion is left completely undefined. Does it mean N lines before and after each match? Does it mean capturing the full `--- FAIL` block for Go tests? Does it handle multi-line stack traces where the file reference only appears on one line? The proposal's own Key Risks table acknowledges "输出行关联测试文件的准确性" with likelihood M, but the mitigation -- "多匹配一些行（宁可多包含）比漏掉好，agent 可自行过滤" -- is hand-waving. Over-inclusion can be just as problematic as under-inclusion: if the filtering is too loose, two adjacent test files' outputs bleed into each other's fix tasks, partially recreating the original scope problem at a smaller scale.

风险：移除 `maxFixTasksPerStep` cap 的论据建立在理想假设上。

The proposal asserts: "移除 `maxFixTasksPerStep` cap 限制——拆分后每个 fix task scope 已收窄到单文件，cap 不再必要". The Assumptions Challenged table reinforces this: "cap 是因为单个 fix task scope 过大导致循环". But the cap was introduced for a different reason than scope -- it was a safety valve against runaway task creation loops, as documented in the existing forensics reports (`docs/forensics/fix-task-loop/report.md`). The cap prevents a scenario where quality-gate repeatedly creates fix tasks that each fail and trigger more fix tasks. With the cap removed, if 30 test files each have one failure, the system creates 30 fix tasks simultaneously. Each fix task, when claimed and failed, could trigger another quality-gate cycle that creates another batch. The proposal's risk table rates "移除 cap 后大量 fix task 并发" as likelihood L, but this assessment lacks supporting evidence. The existing `countFixTasks` function and its comprehensive test coverage (including `TestAddFixTask_CapEnforced`, `TestAddFixTask_CrossStepIndependence`) represent deliberate design intent that should not be discarded without a more rigorous argument.

问题：`extractSourceFiles` 复用假设未经验证 -- 该函数返回的是所有源文件，不区分测试文件和产品文件。

The proposal states: "复用现有 `extractSourceFiles` 提取文件路径". But examining the actual implementation in `quality_gate.go` (lines 584-607), `extractSourceFiles` extracts ALL source files from the output -- both test files and production files. In a typical test failure, the output contains references to both: `handler_test.go:42: Expected 200` (test file) and `handler.go:108: actual return` (production code). The proposal's `addRegressionFixTasks` would need to first filter `extractSourceFiles` results to only test files via `isTestFile`, then group by test file. But this means the production file references -- which are often the actual bug locations -- are discarded during grouping.

This creates an asymmetry: the fix task is grouped by test file but the agent needs to read and modify production code. If two test files fail because of the same production code bug, the proposal creates two fix tasks that will attempt to fix the same root cause independently, potentially conflicting. The proposal's comparison table dismisses "按目录分组" as "未改进", but directory-based grouping naturally co-locates test and production code, avoiding this exact problem.

风险：Rust 等"无特殊测试文件命名"语言的 fallback 路径实际回退到了提案声称要改进的方案本身。

The proposal acknowledges: "Rust 等无特殊测试文件命名的语言" with likelihood L and fallback to directory-based grouping. But Rust is not edge-case -- it is a top-10 language in the `sourceExts` whitelist (`.rs` is already supported). In Rust, tests live in the same file as production code (using `#[cfg(test)]` modules), or in a `tests/` directory with arbitrary filenames. There is no `*_test.rs` convention. The `isTestFile` function would not match any Rust file, so ALL Rust failures would fall through to the current `addFixTask` with `groupFilesByDir` behavior -- which is exactly the behavior the proposal exists to improve. Similarly, C/C++ test files often follow no standard convention (Google Test files can be named anything). The proposal's success criterion "支持至少 5 种语言的测试文件命名约定（Go、Python、JS/TS、Java、Ruby）" implicitly acknowledges that other languages are second-class citizens.

问题：与 compile/fmt/lint/unit-test 步骤的隔离保证不足。

The Scope section declares: "unit-test / compile / lint 步骤的改动" as Out of Scope, and the proposal only targets `runTestRegression` failures. But examining the actual call sites in `quality_gate.go`, `addFixTask` is called from three locations: the gate step callback (line 163), `runUnitTestStep` (line 516), and both `runTestRegressionLegacy` (line 259) and `runTestRegressionSurface` (line 288). The proposal introduces `addRegressionFixTasks` as a replacement only for the last two call sites. The `addFixTask` function itself remains unchanged and still contains the cap check with `maxFixTasksPerStep`. If the proposal removes the cap constant and `countFixTasks`, these remaining call sites will break at compile time. If the proposal does not remove them, then there are two parallel fix-task creation paths with different cap policies -- one capped, one uncapped -- which is confusing and error-prone. The isolation boundary between old and new code paths is not cleanly defined.

## Section 3: Improvement Suggestions

建议：为输出行关联算法增加精确的伪代码规范。

The proposal should specify the line-association algorithm before implementation. A minimal specification would be: (1) split the output into lines, (2) for each test file `T`, collect all lines where `T`'s path appears as a substring, (3) for each matched line, also include N preceding and N following lines as context (specify N, e.g., 2), (4) deduplicate overlapping context windows, (5) for lines that match multiple test files (impossible by definition since paths are unique, but stack traces may reference shared utility files), include them in all matching groups. This addresses the first risk by making the "上下文" concept concrete and testable. The success criterion should include a specific test case: given a realistic multi-file test output (e.g., the actual 20+ failure scenario from the gotcha document), the algorithm produces the expected line assignments.

建议：将 cap 移除改为 cap 提升，保留安全阀。

Instead of removing `maxFixTasksPerStep` entirely, raise it to a value that accommodates realistic per-file splitting (e.g., 10 or 15) while still preventing runaway loops. The proposal should also document the loop-breaker mechanism explicitly: what prevents a fix task from failing and triggering another quality-gate cycle that creates more fix tasks? If the existing loop-breaker (documented in `docs/proposals/quality-gate-fix-task-loop-breaker/`) relies on the cap as part of its defense, removing the cap undermines it. The updated proposal should state: "raise `maxFixTasksPerStep` from 3 to N where N is derived from [analysis], and retain `countFixTasks` for the loop-breaker." This preserves the safety properties while allowing the finer granularity.

建议：明确 `extractSourceFiles` 的使用方式 -- 是复用还是封装。

The proposal should clarify whether `addRegressionFixTasks` calls `extractSourceFiles` directly and then filters the result with `isTestFile`, or whether it introduces a new function `extractTestFiles` that combines extraction and filtering. The latter is preferable because it avoids the semantic mismatch: `extractSourceFiles` is designed for "what files are mentioned in this error" (used by surface inference), while the new use case needs "which test files have failures in this output." These are different questions. A dedicated `extractTestFiles` function could also preserve the association between each test file and its specific output lines, which `extractSourceFiles` deliberately discards (it returns a flat comma-separated string). This addresses the third risk by making the data flow explicit and avoiding the loss of file-to-line mapping information.

建议：为 fallback 场景提供定量评估而非定性声明。

Rather than listing Rust in a risk table with likelihood L, the proposal should state: "In the current Forge codebase, X% of regression failures come from Rust code" (likely 100% since Forge CLI is written in Go, but Forge users' projects may use Rust). If the Forge project itself is the primary user of this feature, and its tests are Go, then the Rust fallback is truly low-impact. But if the feature is meant to be general-purpose for any Forge user's project, then Rust support matters. The proposal should either (a) scope the feature explicitly to "Go, Python, JS/TS, Java, Ruby projects with standard test naming conventions" and acknowledge that other projects get no improvement, or (b) provide a secondary grouping strategy for fallback cases -- for instance, when `isTestFile` matches zero files, fall back to grouping by top-level directory rather than immediate parent directory, which would give coarser but still-better-than-nothing splitting for Rust monorepo layouts.
