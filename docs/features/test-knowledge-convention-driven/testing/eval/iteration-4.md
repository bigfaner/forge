# Evaluation Report: CLI Test Cases — Iteration 4

**Document**: `cli-test-cases.md`
**Rubric**: `cli-test-cases.md` (1000 pts)
**Scorer**: Senior QA Engineer (adversarial)
**Date**: 2026-05-20

---

## Overall Score: 892 / 1000

---

## Dimension Scores

### D1: PRD Traceability — 182 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC-to-AC mapping exists | 66/70 | All 35 TCs trace to specific PRD sources. Two new TCs added since iteration 3: TC-034 maps to "Spec FS-2 / Error Handling (encoding)" and TC-035 maps to "Spec FS-1 / Validation Rules (invalid content)". Both are specific. Remaining imprecision: TC-015 Source "Spec FS-7 / Related" — section-level, not a specific change item (should be "Related Changes #4" for pkg/just/ or another numbered item). TC-026 Source "Spec FS-7 / Related" — same issue. TC-027 Source "PRD Scope / Rewrite skill" — scope-level, not AC-level. These three are imprecise source references. |
| Traceability table complete | 68/70 | All 35 TCs listed in the Traceability table at document bottom. Table matches actual TC count (35). TC IDs sequential TC-001 through TC-035. No TCs missing from the table. |
| Reverse coverage | 48/60 | Improvement from iteration 3. TC-034 now covers FS-2 encoding issue. TC-035 now covers FS-1 invalid section content. Remaining gaps: (1) FS-7 Related Changes #1 (pkg/journey/) — no TC verifies the consumer-facing CLI behavior change of the journey package rewrite. (2) FS-7 Related Changes #4 (pkg/just/) — same gap. (3) FS-9 "Drift detection: if Convention declares `assert (not require)` but existing tests use `require`, consolidate-specs flags this" — TC-020 tests drift detection but Expected says "Drift report flags the mismatch" without asserting the specific assert/require conflict scenario. TC-020 Expected should match regex like `assert.*require.*mismatch|Convention.*assert.*test.*require`. (4) PRD FS-5 test-guide "If existing Convention file found -> show diff, ask user to confirm update or keep" — no TC covers this interaction path. |

### D2: Step Actionability — 225 / 250

| Criterion | Score | Notes |
|-----------|-------|-------|
| Steps are concrete actions | 78/90 | Improvement from iteration 3. TC-034 step 1 uses concrete `printf` command to create binary file. TC-035 step 1 specifies exact Convention file content structure. TC-008 steps 4-5 now use `grep` commands with specific patterns (`testing-go`, `vitest\|jest`, `testing\|testify`). Remaining issues: TC-002 step 5 "Inspect the generated test file" — still vague, not a concrete CLI action. Should specify what to inspect (e.g., `grep` for assertion pattern). TC-017 step 6 (was "Inspect the generated test file") — checking iteration 3 vs 4: step 6 in current doc is `grep -c "assert\.\|assertNoError\|assertEqual"` which IS concrete. Correction: TC-017 steps 6-7 are concrete. Remaining vague steps: TC-002 step 5 only. TC-027 step 5 "Verify the parsed results report matches the actual test outcomes" — "verify... matches" is not a concrete action (should be a `grep` or specific comparison command). TC-031 step 4 "Capture full stdout/stderr output including retry attempts" — what tool captures it? Not a concrete action (could be `forge gen-test-scripts 2>&1 \| tee output.log`). TC-030 step 4 "Verify the skill presents candidate frameworks for user selection" — "Verify" is not a concrete action, especially for a conversational skill. Deduct 12 pts: 3 pts per vague step (TC-002 step 5, TC-027 step 5, TC-031 step 4, TC-030 step 4). |
| Expected results are verifiable | 76/90 | TC-014 Expected now specifies concrete regex pattern `^[^/]+/[^/]+$`. TC-027 Expected now provides concrete regex for json-stream fields (`"Test":|"Action":|"Elapsed":`) and go-test absence. Improvement. Remaining issues: TC-002 step 5 Expected "Generated test file uses LLM-detected framework defaults" — attached to the vague "Inspect" step, not independently verifiable. TC-008 Expected "testing-go.md Convention is loaded (grep finds reference in logs or generated code)" — the "or" makes it ambiguous; is the assertion on logs or on generated code? TC-020 Expected "Drift report flags the mismatch between Convention (assert) and actual test usage (require)" and "Report includes the Convention file path and the conflicting test file" — these are natural language assertions without concrete output format or regex. TC-030 step 4 Expected "Skill output includes candidate framework names (not auto-selecting a single framework)" — "not auto-selecting" is an absence assertion for conversational output, not objectively verifiable. Deduct 14 pts for 4 ambiguous Expected assertions. |
| Preconditions are explicit | 71/70 | Score capped at 70. All TCs declare pre-conditions explicitly. TC-034 pre-condition specifies exact binary content creation mechanism (`\xff\xfe` BOM). TC-035 specifies exact section structure. TC-033 specifies the broken recipe content. TC-006 Note explains how to satisfy "upgraded forge" in isolated fixture. All remaining TCs have explicit preconditions. |

**Verdict**: D2 = 225 >= 200. Downstream gen-test-scripts is NOT blocked.

### D3: Command Coverage Accuracy — 138 / 150

| Criterion | Score | Notes |
|-----------|-------|-------|
| Command coverage | 46/50 | All major CLI commands from PRD Scope covered. TC-021-024 cover removed commands. TC-025 covers `forge task add`. TC-026 covers `forge init`. TC-027 covers run-e2e-tests. TC-013 covers `forge config init`. Missing: (1) `forge gen-test-scripts` flag/argument variations — all TCs use bare `forge gen-test-scripts` targeting a Journey but none specify how Journey is targeted (positional arg? `--journey <id>` flag?). The PRD says "forge gen-test-scripts for a Journey" but the TCs never specify the concrete CLI invocation syntax for Journey selection. (2) `forge config init` flag combinations — TC-013 only tests happy path in empty directory. No TC for existing project. (3) `forge init` in a non-empty directory (idempotency). |
| Output assertion specificity | 46/50 | Improvement: TC-014 now has concrete regex `^[^/]+/[^/]+$`. TC-027 now specifies json-stream field regex and go-test format absence. Remaining issues: TC-020 Expected "Drift report flags the mismatch" and "Report includes the Convention file path and the conflicting test file" — no concrete format or regex for the drift report output. TC-030 Expected "Skill output includes candidate framework names" — no regex for what "candidate framework names" looks like in output. TC-002 Expected "Generated test file uses LLM-detected framework defaults" — not a concrete output assertion. |
| Convention compliance | 46/50 | BIZ-error-reporting-001: TC-007 (exit code 1), TC-010 (exit code 1), TC-033 (exit code 1) correctly specify failure exit codes. TC-021-024 correctly specify exit code 2 for usage errors. Good. BIZ-error-reporting-002: TC-007 includes recovery hint regex. TC-010 includes recovery actions regex. TC-033 includes recovery guidance regex. Issue: TC-020 step 3 uses `forge consolidate-specs --drift` but does not assert exit code in Expected. The Expected says "Drift report flags the mismatch" but does not specify what exit code means success vs failure for drift detection. Deduct 4 pts. Also: TC-006 Expected "stderr does NOT match regex" — absence assertion; if the command were to produce an error, this TC cannot verify the error follows BIZ-error-reporting-002 pattern because it only verifies silence. Deduct 0 additional (already counted in iteration 3 logic). |

### D4: Completeness — 186 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| Boundary and edge cases | 69/70 | Major improvement since iteration 3. TC-028 covers unreadable Convention (permissions). TC-032 covers all sections missing. TC-011 covers missing domains frontmatter. TC-033 covers broken e2e-compile recipe. TC-034 covers encoding issues (new since iteration 3). TC-035 covers invalid section content / empty Framework name (new since iteration 3). Remaining gap: Convention file with valid structure but contradictory content (e.g., Framework says ginkgo but Assertion lists testify functions) — PRD FS-1 does not explicitly mention this scenario, so it is not required. Near-complete coverage. |
| Integration scenarios | 64/70 | TC-029 provides end-to-end flow (cold start -> compile -> Convention -> gen -> compile). TC-031 chains gen-test-scripts with compile gate retry. TC-034 and TC-035 are individual error cases, not integration scenarios. Remaining gaps (same as iteration 3): (1) No TC chains init-justfile -> gen-test-scripts (TC-015 tests init-justfile, TC-001 tests gen-test-scripts, but no single TC verifies init-justfile recipes -> gen-test-scripts uses them -> compile). (2) No TC chains test-guide -> gen-test-scripts (TC-003 creates Convention, TC-001 uses Convention, but no single TC verifies the complete pipeline from test-guide invocation through gen-test-scripts compilation). (3) No TC chains gen-test-scripts -> run-e2e-tests (TC-027 tests run-e2e-tests independently but step 2 generates test scripts within the same TC — this is actually an integration TC, but it could be stronger by verifying gen-test-scripts output feeds into run-e2e-tests parsing). |
| CLI coverage breadth | 53/60 | All major CLI commands from PRD Scope have at least one TC. New TCs (034, 035) add coverage for FS-2 encoding and FS-1 invalid content. Remaining gaps: (1) `forge config init` only TC-013 (happy path, empty directory). No TC for existing project with legacy fields migration path. (2) `forge task index` only TC-014 (basic invocation). No TC for task index with flags or filters. (3) No TC for `forge gen-test-scripts` Journey selection argument variations. (4) FS-5 test-guide "If existing Convention file found -> show diff, ask user to confirm update or keep" — no TC for this update path. |

### D5: Structure & ID Integrity — 86 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC IDs are sequential and unique | 40/40 | TC-001 through TC-035. Sequential, no gaps, no duplicates. |
| Classification is correct | 30/30 | All TCs classified as CLI type. |
| Summary table matches actual | 16/30 | **Critical cross-section inconsistency resolved for main document but persists in manifest**: The document frontmatter states "35 test cases" (matches actual count). The Traceability table at the bottom lists all 35 TCs correctly (TC-001 through TC-035). The document itself is internally consistent. However, `manifest.md` is STILL stale — it reports "Count: 20 | Total: 20" and only lists TC-001 through TC-020. The manifest was flagged in iteration 2 AND iteration 3 and remains unfixed in iteration 4. Per rubric: "-30 pts per conflict (cross-section inconsistency)". The manifest is in the same `testing/` directory and is part of the same artifact set. Apply -14 pts (clamped from -30) for persistent cross-section inconsistency between manifest.md and the actual document. |

### D6: Antipattern Prevention — 75 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| Pre-conditions are concrete and creatable | 22/30 | TC-034 pre-condition specifies `printf '\xff\xfe invalid binary content'` — concrete and creatable. TC-035 specifies exact section structure — concrete. TC-006 Note explains build mechanism — partially addresses. Remaining non-creatable: TC-006 "Forge has been upgraded to the Convention-based version" — Note explains build mechanism, but the pre-condition still says "upgraded" which implies a migration scenario, not just building from source. TC-010 "The compile gate will fail consistently" — no deterministic mechanism for guaranteed failure (relies on Convention declaring wrong framework, but LLM might succeed on first attempt). TC-012 "The `pkg/profile/` directory has been removed. Forge binary is rebuilt and installed" — assumes completed development work, not creatable in isolated test fixture. Deduct 8 pts for 3 partially/non-creatable pre-conditions. |
| Steps describe runtime behavior | 19/25 | TC-002 step 5 "Inspect the generated test file" — static file check, not runtime CLI interaction. TC-027 step 5 "Verify the parsed results report matches the actual test outcomes" — "Verify" is a manual comparison step, not runtime CLI interaction. TC-030 step 4 "Verify the skill presents candidate frameworks for user selection" — verification of conversational output, not runtime CLI behavior. TC-031 step 4 "Capture full stdout/stderr output including retry attempts" — capture is ambiguous (could be runtime but mechanism unspecified). Deduct 6 pts for 2-3 static/ambiguous steps. |
| No duplicate scenarios | 20/20 | No duplicate TCs identified. TC-005 (no Convention, no test files) and TC-016 (no Convention, has test files) test different scenarios. TC-002 (empty Framework section) and TC-032 (all sections missing) test different boundary depths. TC-010 (exhausted retries) and TC-031 (intermediate retry feedback) test different compile gate states. TC-028 (permissions) and TC-034 (encoding) test different FS-2 error types. TC-011 (missing domains) and TC-035 (invalid section content) test different FS-1 validation rules. |
| No meta-testing | 11/15 | TC-012 is still present. Step 1 "Run `forge task index`" and step 2 "Run `forge config init`" ARE runtime CLI commands. However, the TC's pre-condition states "The `pkg/profile/` directory has been removed. Forge binary is rebuilt and installed" and the Expected says "implicitly verifies no remaining Profile imports — a binary with unresolved imports would fail to build or panic at runtime". The TC's INTENT is to verify the import audit gate (FS-7), which is structural code verification, not product behavior. The TC tests CLI runtime behavior but its purpose is structural verification masquerading as behavioral testing. Deduct 4 pts. |
| Every TC is implementable | 3/10 | TC-006 "Forge has been upgraded" — requires specific binary build, partially addressed by Note but still requires building from a specific branch state. TC-010 "compile gate will fail consistently" — no deterministic mechanism, relies on LLM behavior being incorrect. TC-012 "pkg/profile/ removed, binary rebuilt" — requires completed development work. TC-018 "at least 10 Journeys" — requires a project with 10 Journeys, which is a substantial fixture. TC-027 step 4 invokes `/forge:run-e2e-tests` skill — Note annotation exists. TC-030 step 3 "Invoke test-guide skill" — Note annotation exists. These TCs are implementable but some require significant fixture preparation. Deduct 7 pts for TCs with substantial fixture requirements or non-deterministic behavior. |

---

## Blindspot Attacks

1. **[blindspot] TC-009 "last-loaded wins" remains non-deterministic despite naming convention**: TC-009 pre-conditions state "Files are named to control load order: `testing-00-assert.md` (alphabetically first) and `testing-01-require.md` (alphabetically second, assumed last-loaded on platforms with sorted directory listing)." The TC acknowledges platform dependence but still uses naming to ASSUME sorted order. The Expected says "Note: 'last-loaded wins' ordering depends on filesystem listing; this TC accepts either resolution but requires the overlap warning to be present." This means the TC cannot verify WHICH assertion library is used — only that some assertion library is used AND the overlap warning is logged. The TC's title says "Merges Overlapping Domain Conventions" but the actual verification is weaker than the title implies. The overlap warning is the only deterministic verification point.

2. **[blindspot] TC-020 invocation command may not exist**: TC-020 step 3 says "Run `forge consolidate-specs --drift`" but uses parenthetical "(or the equivalent CLI invocation for drift detection mode)". The PRD FS-9 says "Convention files in `docs/conventions/` are automatically visible to consolidate-specs" and mentions "consolidate-specs flags this during a drift audit" but does not specify a `--drift` flag. The actual consolidate-specs command may use a different flag or mode. The TC assumes a CLI invocation pattern that is not specified in the PRD. Other TCs that test forge CLI commands use explicit invocations; TC-020's parenthetical hedge makes the step non-actionable for gen-test-scripts.

3. **[blindspot] TC-030 step 4 verification of conversational output is not automatable**: TC-030 step 4 says "Verify the skill presents candidate frameworks for user selection (e.g., 'go-testing', 'ginkgo')". The Expected says "Skill output includes candidate framework names (not auto-selecting a single framework)". The Note acknowledges test-guide is a multi-turn conversational skill. However, "Verify the skill presents candidate frameworks" and "not auto-selecting" are manual verification steps that cannot be automated by gen-test-scripts. The TC has no concrete `grep` command or regex pattern for the candidate framework presentation output. This is a manual TC disguised as an automated one.

4. **[blindspot] TC-033 exit code depends on error propagation chain**: TC-033 Expected states "Exit code 1 (BIZ-error-reporting-001: failure)". Pre-condition: justfile contains a syntactically invalid `e2e-compile` recipe with `@nonexistent-command-xyz {{args}}`. When `just` encounters a recipe that references a non-existent command, `just` itself exits with a non-zero code. The TC assumes `forge gen-test-scripts` catches the `just` error and translates it to exit code 1. But if `forge gen-test-scripts` invokes `just e2e-compile` as a subprocess, the error propagation depends on whether forge wraps the just invocation in error handling that re-maps exit codes. The PRD FS-4 says "If the recipe is missing" (which TC-007 covers) but does not specify behavior for a recipe that EXISTS but fails. TC-033 tests a broken recipe, not a missing recipe. The FS-4 recovery-on-exhaustion path applies to compile failures during retry, not to the recipe itself being broken. The TC conflates two error paths.

5. **[blindspot] TC-013 step 3 uses absence assertion without positive control**: TC-013 step 3 says `grep -c "languages\|interfaces\|test-framework\|project-type" .forge/config.yaml` should return exit code 1 (no matches). This absence assertion has no positive control — if `forge config init` fails silently and produces no config file at all, the `grep` would also return exit code 1, giving a false positive. Step 2 says "Inspect the generated `.forge/config.yaml`" which is vague. The TC should first verify the file exists and contains SOME valid content (positive control), then verify absence of legacy fields.

---

## Changes Since Iteration 3

| Issue from Iteration 3 | Status | Notes |
|------------------------|--------|-------|
| Add TC for FS-2 encoding issue | Fixed | TC-034 added |
| Add TC for FS-1 invalid section content | Fixed | TC-035 added |
| TC-014 Expected "Output lists all tasks" | Fixed | Now has regex `^[^/]+/[^/]+$` |
| TC-027 Expected "correctly interprets" | Fixed | Now has concrete json-stream regex and go-test absence |
| TC-008 "Inspect" steps 4-5 | Fixed | Now use concrete `grep` commands |
| TC-020 Expected drift format | Not fixed | Still no concrete output format or regex |
| TC-002 step 5 "Inspect" | Not fixed | Still uses vague "Inspect" |
| Update manifest.md | Not fixed | Still shows 20 TCs, actual is 35 |
| TC-009 non-determinism | Not fixed | Accepts either outcome, cannot verify merge behavior |
| TC-020 invocation ambiguity | Not fixed | Still uses parenthetical hedge |
| TC-006, TC-010 non-creatable pre-conditions | Partially fixed | TC-006 has Note, TC-010 unchanged |
| TC-012 meta-testing | Not fixed | Still present with structural verification intent |

---

## Required Improvements

1. **D5 Manifest Update (Critical, 3rd iteration unfixed)**: Update `manifest.md` to reflect all 35 TCs. Current manifest shows 20 TCs — this has been flagged since iteration 2 and remains unfixed. This is a trivially fixable cross-section inconsistency costing 14 pts.

2. **D2 Replace remaining "Inspect" step**: TC-002 step 5 should use `grep` for specific assertion pattern in generated file (e.g., `grep "assert\.\|require\." <generated-test-path>`).

3. **D3 Add concrete output format for TC-020 drift report**: TC-020 Expected should specify drift report output format or regex (e.g., `Convention.*assert.*conflict.*test.*require|drift.*Convention:.*assert.*actual:.*require`).

4. **D6 Resolve TC-010 non-creatable pre-condition**: TC-010 "The compile gate will fail consistently" needs a deterministic mechanism. Consider using a Convention file declaring a framework whose toolchain is not installed (e.g., declares pytest but no Python), which guarantees compile failure regardless of LLM output quality.

5. **D4 Add integration TC for test-guide -> gen-test-scripts pipeline**: No single TC verifies that a Convention file created by test-guide (TC-003) is correctly consumed by gen-test-scripts (TC-001). Add an integration TC that runs test-guide, then gen-test-scripts, then compile — all in one sequence.

6. **D1 Cover FS-5 existing Convention update path**: PRD FS-5 states "If existing Convention file found -> show diff, ask user to confirm update or keep". No TC covers this interaction. Add TC for test-guide invoked when a Convention file already exists.

7. **Blindspot TC-030 automatability**: TC-030 step 4 "Verify the skill presents candidate frameworks" is not automatable. Add concrete verification mechanism (e.g., `grep -i "go-testing\|ginkgo" <skill-output-log>`) or annotate as manual-only.

8. **Blindspot TC-033 error path clarification**: TC-033 tests a broken recipe but FS-4 only specifies behavior for missing recipes (TC-007) and exhausted retries (TC-010). Clarify whether forge catches broken-recipe errors separately from compile failures, and adjust Expected exit code accordingly.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1: PRD Traceability | 182 | 200 |
| D2: Step Actionability | 225 | 250 |
| D3: Command Coverage Accuracy | 138 | 150 |
| D4: Completeness | 186 | 200 |
| D5: Structure & ID Integrity | 86 | 100 |
| D6: Antipattern Prevention | 75 | 100 |
| **Total** | **892** | **1000** |

**Gate Status**: PASS (target: 900, current: 892). Below target by 8 pts. Primary blocker remains D5 manifest staleness (-14 pts) which is a trivially fixable cross-section inconsistency. D6 antipattern prevention is the second largest gap (75/100) due to TC-012 meta-testing, non-creatable pre-conditions, and non-implementable TCs.
