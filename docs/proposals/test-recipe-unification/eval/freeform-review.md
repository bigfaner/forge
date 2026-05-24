# Freeform Expert Review

## Background Assessment

This proposal addresses a three-layer misalignment in Forge's test infrastructure following a profile-removal refactoring exercise. The core problem: the justfile recipe model was never updated to reflect the new convention-driven, config.yaml-based architecture. This manifests as three concrete pain points -- breaking tasks running full test suites (~90s) on every submit, init-justfile templates hardcoding Playwright for all languages, and config field names (`auto.e2eTest`) that describe only a subset of what they control.

The proposed solution introduces a two-tier test recipe model: `just unit-test` for language-level unit tests (fast, per-submit gate) and `just test` for surface-level advanced tests (slow, all-completed gate). The `e2e-test`/`e2e-setup`/`e2e-verify` naming is retired in favor of `test` + auxiliary recipes. The config key `auto.e2eTest` is renamed to `auto.test`, and the task key `run-e2e-tests` becomes `run-test`.

The proposal frames this as a non-innovative, industry-standard test pyramid mapping -- which is honest and appropriate. The core design decision is that Forge remains surface-agnostic: it calls `just test` without knowing whether that invokes Playwright, pytest, or curl-based API integration tests. This is a sound architectural choice that avoids coupling Forge to any specific test framework.

The proposal claims ~42 files affected across five tiers (Go source, prompts/skills, justfile templates, tests, documentation), with an estimated 10-12 coding tasks. It leverages the v3.0.0 major version as a clean break -- no backward compatibility logic is planned, with the expectation that users re-run `init-justfile` to regenerate their justfiles.

## Key Risks

After tracing the data flow from config schema through CLI detection into task generation and gate execution, I have identified several risks that warrant attention before implementation.

风险：`RunProjectTests` 探测链与新 recipe 模型的语义冲突

The proposal states that `testrunner.go` should detect and call `unit-test` instead of `test`. However, examining the actual `RunProjectTests` function in `forge-cli/pkg/testrunner/testrunner.go`, its current fallback chain is: `just test` -> `make test` -> `go test ./...` -> `npm test` -> `pytest`. The proposal's own quote says the function's "探测和调用 `unit-test` 替代 `test`"，but `RunProjectTests` is only called from `quality_gate.go` Step 2 (the unit/integration test step), not from `submit.go`. The submit flow uses `DefaultGateSequence()` which runs `just test` as a gate recipe. This means there are two distinct code paths that invoke tests, and the proposal conflates them. If `RunProjectTests` is changed to probe `unit-test` first but the fallback chain still eventually hits `go test ./...` (which is a unit test), the behavior is consistent. But if someone changes the justfile to have `unit-test` run something other than unit tests, the fallback's semantic meaning diverges. The proposal does not address how the fallback chain in `RunProjectTests` should adapt to the new two-tier model.

> "pkg/testrunner/testrunner.go | RunProjectTests() 探测和调用 unit-test 替代 test"

The consequence: the `RunProjectTests` fallback chain (go test, npm test, pytest) has unit-test semantics baked in, but the justfile recipe it probes first changes name from `test` to `unit-test`. If a project has no justfile, the fallback correctly runs unit tests. But if a project has a justfile with `test` (now meaning "advanced tests"), the function would skip it and fall through -- which is actually the desired behavior. The risk is that the proposal does not explicitly document this fallback interaction, leaving implementers to infer the correct behavior.

问题：`quality_gate.go` 中 `runE2ERegression` 的 `e2e-test` 到 `test` 迁移存在语义歧义

The current `runE2ERegression` function in `quality_gate.go` (line 204) explicitly checks `just.HasRecipe(projectRoot, "e2e-test")` and then runs `just e2e-test`. The proposal says this should become `just test`, but `quality_gate.go` Step 2 already runs `RunProjectTests` which currently calls `just test`. After the migration, Step 2 should run `just unit-test` (via RunProjectTests probing `unit-test`) and Step 3 should run `just test`. The proposal's impact table correctly identifies this split, but there is a gap: the current `addFixTask` function (line 430-431) has a hardcoded mapping where `step == "unit-test"` maps to `testScript = "just test"`. This mapping becomes incorrect after the rename -- if the step name is `unit-test`, the fix task should reference `just unit-test`, not `just test`. The proposal does not call out this specific line change.

> "internal/cmd/quality_gate.go | Step 2 用 unit-test，Step 3 用 test；addFixTask 映射更新"

The consequence: without updating this specific mapping logic, fix tasks generated from unit-test failures would instruct the agent to run `just test` (the advanced test suite) instead of `just unit-test`, leading to incorrect error reproduction.

风险：`DefaultGateSequence` 的 recipe 名称变更与 `submit.go` 的 breaking 判断逻辑耦合不清

The proposal introduces a new `UnitGateSequence()` for breaking tasks, and states that `DefaultGateSequence()` should change `test` to `unit-test`. But `submit.go` (line 366-369) uses this logic: `if breaking { steps = just.DefaultGateSequence() } else { steps = just.LintGateSequence() }`. Currently `DefaultGateSequence` returns `compile -> fmt -> lint -> test`. After the rename, if `DefaultGateSequence` returns `compile -> fmt -> lint -> unit-test`, then breaking tasks on submit run `unit-test` -- which is correct. But the proposal also says a new `UnitGateSequence()` should be created. This raises the question: does `DefaultGateSequence` become `UnitGateSequence`, or are they separate? The proposal lists both changes: "DefaultGateSequence() 的 test -> unit-test；新增 UnitGateSequence()". If both exist, what is `DefaultGateSequence` used for after the migration? Is it renamed to `UnitGateSequence` and a new `DefaultGateSequence` is created for all-completed? The proposal is ambiguous about the exact semantics of these three sequence functions after migration.

> "pkg/just/just.go | DefaultGateSequence() 的 test -> unit-test；新增 UnitGateSequence()"

The consequence: an implementer could create three sequences (Default, Unit, Lint) where Default and Unit are identical, or where Default still has `test` and Unit has `unit-test`, leading to confusion about which to use where. The proposal needs to define exactly what each sequence contains post-migration.

问题：`auto.e2eTest` 到 `auto.test` 的重命名与 `parseAutoRaw` 硬编码字段列表不一致

Looking at the actual `config.go`, the `parseAutoRaw` function (line 302) has a hardcoded list of mode fields: `"e2eTest", "consolidateSpecs", "cleanCode", "validation", "runTasks", "knowledgeSave"`. Similarly, `autoModeField` (line 429) has a case for `"e2eTest"`. The proposal states the YAML tag changes from `e2eTest` to `test` and there is "直接替换无兼容层". However, the YAML tag `test` is potentially ambiguous -- it could collide with a hypothetical future top-level `test` key or with `test-framework`. More critically, the `autoModeField` function dispatches on the Go field name string, and changing `"e2eTest"` to `"test"` means `SetConfigValue("auto.test", ...)` would work, but `SetConfigValue("auto.e2eTest", ...)` from any existing script, hook, or documentation would silently fail with "unknown config key" -- there is no deprecation warning. For a v3.0.0 breaking change this is acceptable, but the proposal claims "用户运行 forge init 或手动更新即可" which understates the failure mode. Users who have scripts or CI pipelines that call `forge config set auto.e2eTest.full false` will get an opaque error.

> "pkg/forgeconfig/config.go | E2eTest -> Test，YAML tag e2eTest -> test，直接替换无兼容层"
> "auto.e2eTest -> auto.test 需更新现有 config.yaml | H | L | 直接重命名，用户运行 forge init 或手动更新即可"

风险：`journey_isolation.go` 的 `just e2e-test` 到 `just test` 迁移缺少参数传递约定

The current `ExecuteJourneyInIsolation` function calls `exec.Command("just", "e2e-test", journeyName)`. The proposal says this should become `just test`. But the proposal also defines `just test` as "Surface 级高级测试（Web UI -> e2e，API -> 集成测试）" without mentioning a journey name parameter. The proposal's own Key Scenarios section says: "Journey isolation：just test <journeyName> 运行单个 journey 的高级测试". This means `just test` must accept an optional positional argument. However, the justfile templates listed in Tier 3 do not show this signature. The proposal does not define the `test` recipe's parameter contract -- specifically, whether `test` accepts an optional journey filter argument, and if so, what happens when the justfile is generated for a project that does not support journey-level filtering.

> "pkg/testrunner/journey_isolation.go | just e2e-test -> just test"
> "Journey isolation：just test <journeyName> 运行单个 journey 的高级测试"

The consequence: if the `test` recipe in generated justfiles does not accept a journey name argument, `ExecuteJourneyInIsolation` will fail with a just error ("recipe `test` does not accept argument"), silently breaking all journey isolation tests. This is flagged as "Likelihood: L, Impact: H" in the proposal's own risk table but the mitigation ("just test 需支持 --feature 参数") mentions `--feature` while the code passes a positional argument, not a flag. There is a parameter-style mismatch between the proposed mitigation and the actual call site.

问题：影响范围评估遗漏了 `internal/cmd/test/test.go` 中的 `e2e-test` 引用

Grep results show `internal/cmd/test/test.go` line 36 contains the string "Runs just e2e-test from the project root with the journey name as filter." This file is not listed in any of the five tiers of the Impact Analysis. If this is a CLI command that directly invokes `just e2e-test`, it must also be migrated. Its absence from the impact analysis suggests either the file was overlooked or the string is only in a doc comment -- but even a doc comment would need updating.

> The proposal's Tier 1 lists exactly 8 Go files; `internal/cmd/test/test.go` is not among them.

The consequence: if this file contains functional code (not just comments) that calls `just e2e-test`, the migration would leave a dangling reference. Even if it is only a comment, the inconsistency in documentation reflects incomplete impact analysis.

风险：`test` recipe 的语义过载 -- 同一名称在不同上下文含义不同

Before the migration, `just test` runs the full gate test (currently `go test -race ./...`). After migration, `just test` runs surface-level advanced tests (e2e, integration, etc.), while `just unit-test` runs language-level unit tests. This is a semantic inversion of the existing `test` recipe's meaning. Users who have internalized `just test` as "run unit tests" will now get slow integration/e2e tests when they type `just test`. The proposal acknowledges this is intentional (Test Pyramid alignment), but underestimates the muscle-memory risk. The success criteria state "Breaking 任务 submit 门禁运行 just unit-test（非 just test）" which is correct for the submit flow, but a developer manually running `just test` during development would now trigger the wrong test tier.

> "just test：Surface 级高级测试（Web UI -> e2e，API -> 集成测试），用于 all-completed 门禁"

The consequence: developers accustomed to typing `just test` for quick feedback will instead get slow e2e tests, and must relearn to type `just unit-test` for fast feedback. This is a UX regression in the short term, even if the new naming is architecturally correct long-term.

问题：`RunProjectTests` 的探测顺序暗示 `just unit-test` 是可选的，与 "无 Fallback" 声明矛盾

The proposal's Key Scenarios states: "无 Fallback：v3.0.0 直接要求 unit-test recipe，不回落到 test". But `RunProjectTests` in `testrunner.go` has a five-level fallback chain. If this function is supposed to probe `unit-test` first, what happens when `unit-test` is not found? Does it fall through to `go test ./...`? That would be a fallback to a language-native test command, which is semantically a unit test -- so it is a fallback, just not to `just test`. The "无 Fallback" claim is misleading. What the proposal likely means is "no fallback from unit-test to test (the advanced recipe)", but the actual fallback to language-native test commands still exists and is valuable. This distinction needs to be made explicit, because the success criterion "无 unit-test recipe 时 quality gate 报错提示运行 init-justfile（不 fallback）" directly contradicts the existing `RunProjectTests` behavior which would fall through to `go test` or `npm test`.

> "无 Fallback：v3.0.0 直接要求 unit-test recipe，不回落到 test"
> "无 unit-test recipe 时 quality gate 报错提示运行 init-justfile（不 fallback）"

The consequence: the submit flow via `DefaultGateSequence` uses `HasRecipe` to check for `unit-test` before running it. If `unit-test` is not found and is non-optional, the gate prints a warning and skips. But the success criterion says it should error. The proposal does not reconcile this with the existing `RunGate` behavior where missing non-optional recipes produce a warning but do not fail the gate.

## Improvement Suggestions

建议：明确 `DefaultGateSequence` / `UnitGateSequence` / `LintGateSequence` 迁移后的精确内容
Addresses: `DefaultGateSequence` 语义歧义风险

The proposal should define the post-migration gate sequences as an explicit table:

| Sequence | Steps | Caller |
|----------|-------|--------|
| `UnitGateSequence` | compile -> fmt -> lint -> unit-test | `submit.go` (breaking=true) |
| `DefaultGateSequence` | compile -> fmt -> lint -> unit-test -> test -> probe | `quality_gate.go` all-completed |
| `LintGateSequence` | compile -> fmt -> lint | `submit.go` (breaking=false) |

If `DefaultGateSequence` is being repurposed to mean "full pipeline" instead of "default submit gate", rename it to `FullGateSequence` to avoid the ambiguity of "default" implying "most commonly used." The submit path would use `UnitGateSequence`, not `DefaultGateSequence`. This makes the data flow self-documenting.

建议：为 `test` recipe 定义标准参数签名约定
Addresses: `journey_isolation.go` 参数传递风险

Define in the justfile template contract that `test` accepts an optional positional parameter for journey filtering:

```
test journey="":
    # if journey is set, filter to that journey; otherwise run all
```

The proposal should specify this in the `skills/init-justfile/SKILL.md` "Standard Target Contract" section, and the template generation logic must ensure all generated `test` recipes accept this optional parameter. Additionally, the mitigation for `journey_isolation.go` should specify a positional argument (matching the current `e2e-test` call pattern), not a `--feature` flag.

建议：补充 `RunProjectTests` 的迁移后探测链行为规范
Addresses: `RunProjectTests` 探测链与新 recipe 模型语义冲突 + "无 Fallback" 矛盾

Specify the post-migration `RunProjectTests` probe chain explicitly:

```
1. just HasRecipe("unit-test")? -> just unit-test
2. Has Makefile + make test? -> make test
3. Has go.mod? -> go test ./...
4. Has package.json + test script? -> npm test
5. Has pytest.ini/pyproject.toml? -> pytest
6. WARNING: no unit-test command found
```

Clarify that the "无 Fallback" success criterion applies only to the gate sequence (`RunGate`), not to `RunProjectTests`. The gate sequence should hard-fail when `unit-test` is missing (exit with error + prompt to run `init-justfile`), while `RunProjectTests` retains its fallback chain for projects without justfiles. This requires changing `RunGate` behavior for non-optional recipes from "print warning and skip" to "print error and fail" -- a behavioral change the proposal does not currently call out.

建议：在 `addFixTask` 中修正 `unit-test` 到 `just unit-test` 的映射
Addresses: `addFixTask` 硬编码映射风险

The current code at `quality_gate.go` line 430-431 has:

```go
testScript := "just " + step
if step == "unit-test" {
    testScript = "just test"
}
```

This special case was presumably added because the gate step name `unit-test` was expected to map to the actual command `just test`. After the migration, this mapping should be removed -- `step == "unit-test"` should naturally map to `just unit-test` via the default `testScript = "just " + step`. Similarly, the step `"e2e-test"` case in `fixTypeFromStep` (line 398) should become `"test"`. The proposal should explicitly list these one-line corrections.

建议：补充 `internal/cmd/test/test.go` 到 Tier 1 影响范围
Addresses: 影响范围遗漏问题

Add `internal/cmd/test/test.go` to the Tier 1 Go source code impact table. At minimum, the doc comment on line 36 referencing `just e2e-test` must be updated. If the command itself invokes `just e2e-test` programmatically, it must be migrated to `just test` with the correct argument pattern. This file should be audited during the `/tech-design` phase to determine the full extent of changes needed.

建议：为 `auto.test` 重命名提供迁移提示而非静默失败
Addresses: `auto.e2eTest` 到 `auto.test` 重命名的 CLI 错误 UX

Even in a v3.0.0 breaking change, when a user provides `auto.e2eTest` in config.yaml, the system should not silently ignore it (which is what happens if the YAML field no longer exists). Instead, add a deprecation warning during config parsing: if the raw YAML contains `e2eTest` under `auto`, print a warning to stderr directing the user to rename it to `test`. This costs ~10 lines of code in `parseAutoRaw` and prevents a class of confusing silent-config-ignored bugs. The proposal's claim of "直接重命名，用户运行 forge init 或手动更新即可" assumes users will notice their config is not taking effect, which is optimistic.
