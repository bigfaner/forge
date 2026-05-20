# Evaluation Report: CLI Test Cases — Iteration 2

**Document**: `cli-test-cases.md`
**Rubric**: `cli-test-cases.md` (1000 pts)
**Scorer**: Senior QA Engineer (adversarial)
**Date**: 2026-05-20

---

## Overall Score: 828 / 1000

---

## Dimension Scores

### D1: PRD Traceability — 175 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC-to-AC mapping exists | 65/70 | Significant improvement over iteration 1. All 29 TCs now trace to specific PRD sources at AC or section level. TC-029 traces to "PRD / Main flow" — this is flow-level, not AC-level, but the PRD does not define separate ACs for the main flow, so this is acceptable. TC-015 Source is "Spec FS-7 / Related Changes" — points to a section, not a specific change item. TC-027 Source is "PRD Scope / Rewrite skill" — scope-level, not AC-level. Minor imprecision, not a gap. |
| Traceability table complete | 65/70 | All 29 TCs listed in the traceability table. Missing: TC-027 traces to "PRD Scope / Rewrite run-e2e-tests skill" but the PRD Scope section does not define separate ACs for this item — the TC source is inherently imprecise. TC-029 traces to "PRD / Main flow" which is a flow description, not an AC. |
| Reverse coverage | 45/60 | All 6 Stories and their ACs have at least one TC. Gaps remain: (1) FS-4 compile gate retry mechanism intermediate behavior — TC-010 only covers the exhausted-retry terminal state, not the "feed error back to LLM and regenerate" intermediate step. (2) FS-5 test-guide "no test files" scenario — PRD says "If no test files → present candidate frameworks for user selection" but neither TC-003 nor TC-004 covers the zero-test-file cold start path for test-guide. (3) FS-7 Related Changes #1 (pkg/journey/) and #4 (pkg/just/) still lack dedicated TCs — TC-014 covers task-index (#5), TC-025 covers task-add (#5), but journey and just package rewrites are untested. |

### D2: Step Actionability — 205 / 250

| Criterion | Score | Notes |
|-----------|-------|-------|
| Steps are concrete actions | 70/90 | Improvement over iteration 1. TC-003 steps now use concrete `test -f`, `grep -c`, and `grep` commands instead of vague "verify skill scans" language. TC-004 similarly improved with `ls`, `grep`, `wc -l`. Remaining issues: TC-001 step 5 "Inspect the generated test file for import statements and assertion patterns" is file content inspection, not CLI invocation — the Expected section compensates with grep-like assertions but the step itself is not actionable. TC-005 step 4 "Capture stdout/stderr for hint output" is vague — what tool captures it? TC-008 step 4 "Inspect which Convention file was loaded (via logs or generated code)" — ambiguous inspection mechanism. TC-012 steps 1-3 use `grep -r`, `go build`, `go vet` which are shell/Go commands, not `forge` CLI commands. TC-013 step 2 "Inspect the generated `.forge/config.yaml`" — inspection, not CLI action. TC-027 step 4 "Run run-e2e-tests skill to parse the test results" — unclear how this is invoked (slash command? forge binary?). TC-029 step 8 "Verify the Convention-based generation produces consistent output: `grep -c "testify" <generated-test-file>`" — `<generated-test-file>` is a placeholder, not a concrete path. |
| Expected results are verifiable | 70/90 | TC-007 now correctly specifies exit code 1 per BIZ-error-reporting-001. TC-010 now correctly specifies exit code 1. TC-002 provides regex patterns. Remaining issues: TC-003 Expected "Convention file contains `testify` in the Framework or Assertion section (exit code 0)" — the step says `grep "testify"` which matches anywhere in the file, not specifically in Framework/Assertion section. TC-006 Expected "No errors or warnings related to the legacy fields" — how is "related to" objectively verified? No regex or pattern. TC-008 Expected "Only `testing-go.md` Convention is loaded and applied (verified by generated test code content)" — content-based verification with no concrete assertion pattern. TC-009 Expected "Last-loaded Convention wins for conflicting fields" — which file is "last-loaded" is non-deterministic (filesystem ordering varies). TC-027 Expected "Results report correctly interprets json-stream format output (not default go-test format)" — no concrete assertion for what "correctly interprets" means. TC-020 Expected "Drift report flags the mismatch" and "Report includes the Convention file path and the conflicting test file" — no concrete output format or regex. |
| Preconditions are explicit | 65/70 | Improvement. TC-001 pre-condition now describes ginkgo fixture with specific imports. Remaining issues: TC-006 "Forge has been upgraded to the Convention-based version" — no mechanism to simulate upgrade. TC-010 "The compile gate will fail consistently" — no explanation of how to guarantee deterministic failure in automated test. TC-027 "Generated tests produce json-stream output" — circular dependency (tests must already produce correct output to test result parsing). TC-029 "A standard Go project exists using go-testing + testify" — no fixture creation method. |

**Verdict**: D2 = 205 >= 200. Downstream gen-test-scripts is NOT blocked.

### D3: Command Coverage Accuracy — 120 / 150

| Criterion | Score | Notes |
|-----------|-------|-------|
| Command coverage | 45/50 | Major improvement. TC-021-024 cover removed commands (forge test detect/get/interfaces/framework). TC-025 covers `forge task add`. TC-026 covers `forge init`. TC-027 covers run-e2e-tests. Missing: (1) `forge config init` flag combinations — TC-013 only tests happy path, no flag variations. (2) `forge gen-test-scripts` argument variations — no TC tests gen-test-scripts with specific Journey targeting flags or different argument combinations. (3) `forge init` only tested in empty directory — no TC for `forge init` in an existing project (idempotency). |
| Output assertion specificity | 40/50 | TC-007, TC-010 now specify concrete exit codes (1). TC-002, TC-005, TC-011 provide regex patterns. Remaining issues: TC-006 Expected "No errors or warnings related to the legacy fields" — no regex or concrete assertion. TC-014 Expected "Output lists all discovered tasks" — no concrete output format. TC-020 Expected "Drift report flags the mismatch" — no concrete output. TC-027 Expected "Results report includes per-Journey pass/fail status" — no concrete output format or regex. TC-009 "Last-loaded Convention wins" — non-deterministic. |
| Convention compliance | 35/50 | BIZ-error-reporting-001: TC-007 (exit code 1) and TC-010 (exit code 1) correctly specify failure exit codes. TC-021-024 correctly specify exit code 2 for usage errors. Good compliance. BIZ-error-reporting-002: TC-007 Expected includes actionable message with recovery hint. TC-010 Expected includes recovery actions. Issue: TC-014 "No errors or warnings about missing Profile" — this is an absence assertion, not testing that error messages follow the actionable pattern. TC-025 "No errors or warnings about missing Profile" — same issue. These TCs assert absence of errors but don't test that any error output (if produced) would follow BIZ-error-reporting-002 patterns. Deduct 5 pts for non-specific absence assertions. TC-013 step 3 "Verify no legacy fields... are present" uses `grep -c` which is fine for verification but the Expected doesn't specify what exit code `forge config init` returns — it just says "Exit code 0" without being explicit about stdout/stderr content. |

### D4: Completeness — 155 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| Boundary and edge cases | 55/70 | TC-028 added for Convention file unreadable (permissions) — addresses FS-2 error handling. TC-011 covers missing domains frontmatter. TC-002 covers empty Framework section. Missing: (1) Convention file with ALL required sections missing (not just Framework) — FS-1 validation rules explicitly cover this scenario. TC-002 only tests Framework missing, not the case where Framework + Assertion + Tags + Result Format are all missing. (2) Convention file with encoding issues (FS-2 mentions "permissions, encoding") — TC-028 covers permissions but not encoding. (3) `justfile` with e2e-compile recipe that fails (not missing, but broken) — PRD FS-4 says compile gate runs `just e2e-compile`, but no TC tests a justfile where the recipe exists but is syntactically invalid. |
| Integration scenarios | 55/70 | TC-029 added — end-to-end flow covering cold start gen -> compile -> Convention creation -> gen again -> compile. This is a strong integration test. Missing: (1) No TC covers the gen-test-scripts -> run-e2e-tests pipeline (TC-027 tests run-e2e-tests but doesn't chain from gen-test-scripts output). (2) No TC covers init-justfile -> gen-test-scripts integration (TC-015 tests init-justfile alone, TC-001 tests gen-test-scripts alone, but no TC verifies the full chain: init-justfile creates recipes -> gen-test-scripts uses those recipes -> compile passes). (3) No TC covers test-guide -> Convention -> gen-test-scripts chain (TC-003 creates Convention, TC-001 uses Convention, but no single TC verifies the complete pipeline). |
| CLI coverage breadth | 45/60 | Major improvement with TC-021-029. All major CLI commands from PRD Scope now have at least one TC. Remaining gaps: (1) `forge config init` — only TC-013 (happy path). No TC for running `forge config init` on an existing project with legacy fields (migration path). (2) No TC for `forge init` in a non-empty directory. (3) No TC for `forge task index` with flags or filters (TC-014 only tests basic invocation). (4) PRD FS-7 Related Changes #1 (pkg/journey/) and #4 (pkg/just/) — these are Go package rewrites, not CLI commands, but their consumer-facing behavior (journey-related CLI output, just recipe generation) should have TCs verifying the behavioral change. |

### D5: Structure & ID Integrity — 95 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC IDs are sequential and unique | 35/40 | TC-001 through TC-029. Sequential, no gaps, no duplicates. However, the frontmatter states "29 test cases" and the Traceability table lists 29 entries but the manifest.md was not updated — manifest still reports "Count: 20 | Total: 20" and only lists TC-001 through TC-020. The manifest is now stale and inconsistent with the actual document. This is a cross-section inconsistency per deduction rules. |
| Classification is correct | 30/30 | All TCs classified as CLI type. |
| Summary table matches actual | 30/30 | Traceability table at the bottom of cli-test-cases.md lists all 29 TCs. Matches actual count. |

### D6: Antipattern Prevention — 78 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| Pre-conditions are concrete and creatable | 20/30 | TC-001 step 1 "Set up a Go project fixture with ginkgo imports" — now describes HOW to create the fixture (specific imports). TC-003 step 1 similarly concrete. TC-028 step 2 "Set file permissions to unreadable: `chmod 000`" — concrete creation method. Remaining non-creatable: TC-006 "Forge has been upgraded to the Convention-based version" — no fixture mechanism to simulate upgrade (requires building and installing a specific binary version). TC-010 "The compile gate will fail consistently" — no deterministic mechanism. TC-012 "The `pkg/profile/` directory has been removed. All consumer rewrites are complete" — assumes completed development work, not creatable in isolated test. TC-027 "Generated tests produce json-stream output" — circular pre-condition. |
| Steps describe runtime behavior | 17/25 | Improvement. TC-003/004 now use concrete shell commands (`test -f`, `grep`, `ls`, `wc -l`) instead of vague "verify skill scans". However: TC-001 step 5 "Inspect the generated test file for import statements" — file content inspection. TC-002 step 5 "Inspect the generated test file" — same. TC-005 step 4 "Capture stdout/stderr" — ambiguous capture mechanism. TC-008 step 4 "Inspect which Convention file was loaded (via logs or generated code)" — ambiguous inspection. TC-012 steps 1-3 (`grep -r`, `go build`, `go vet`) — source code audit, not runtime CLI behavior. TC-013 step 2 "Inspect the generated `.forge/config.yaml`" — file inspection. TC-027 step 4 "Run run-e2e-tests skill to parse" — unclear invocation mechanism. |
| No duplicate scenarios | 20/20 | No duplicate TCs identified. TC-005 (no Convention, no test files) and TC-016 (no Convention, has test files) test different scenarios. TC-019 (standard Go project, first-pass compile) and TC-029 (end-to-end flow) cover different scopes. |
| No meta-testing | 11/15 | TC-012 is still a build/import audit — it verifies source code structure, not runtime product behavior. `grep -r "pkg/profile" forge-cli/` is source code search, `go build ./...` is a build check, `go vet ./...` is a static analysis check. This is a meta-test verifying that code was correctly rewritten, not testing that the forge CLI behaves correctly. Deduct per meta-test: TC-012 is 1 meta-test TC. |
| Every TC is implementable | 10/10 | TC-018 no longer requires a Profile-based baseline — the pre-condition states "No Profile package exists (removed per FS-7)" and uses an absolute budget (60 seconds). The contradiction from iteration 1 is resolved. TC-019 uses a standard Go project, implementable. No TC requires unavailable infrastructure without annotation. |

---

## Blindspot Attacks

1. **[blindspot] TC-009 "last-loaded wins" is non-deterministic and cannot produce a reliable pass/fail**: TC-009 Expected states "Last-loaded Convention wins for conflicting fields (e.g., `require` assertions if file B is loaded after file A)". Filesystem directory listing order is not guaranteed — `os.ReadDir` order varies across OS and filesystem types. The TC has no mechanism to control which file is loaded last. The verification step `grep -c "require\." <generated-test-file>` expects `>= 1` but if file A loads last, the generated file would use `assert` instead. This TC will produce flaky results — passing or failing depending on filesystem ordering. The TC needs either: (a) a deterministic naming/ordering mechanism, or (b) an Expected that accepts either `assert` or `require` with a logged overlap note.

2. **[blindspot] TC-015 invokes a slash command, not a forge CLI binary command**: TC-015 step 2 says "Run `/forge:init-justfile`". This is a Claude Code skill invocation (slash command), not a `forge` CLI binary command. Other TCs consistently use `forge gen-test-scripts`, `forge task index`, `forge init`, etc. The test suite mixes two fundamentally different invocation mechanisms. Slash commands execute within Claude Code's agent context; forge CLI commands are standalone Go binaries. This matters because `gen-test-scripts` downstream generates test code that executes forge binary commands, not Claude Code skill invocations. TC-003 and TC-004 correctly annotate this distinction (their Note fields), but TC-015 has no such annotation.

3. **[blindspot] TC-027 has a circular pre-condition**: TC-027 pre-condition states "Generated tests produce json-stream output" and "A valid Journey with generated test scripts exists". But generating test scripts is itself the behavior under test (via `forge gen-test-scripts` in TC-001, TC-019, etc.). The TC pre-condition assumes the system under test already works correctly. Step 2 says "Generate test scripts using `forge gen-test-scripts` for a Journey" — so the TC is both assuming generated tests exist (pre-condition) AND generating them (step 2). This contradiction makes the pre-condition un-creatable independently.

---

## Required Improvements

1. **D1 Reverse Coverage**: Add TC for FS-5 test-guide "no test files" scenario — PRD says "If no test files → present candidate frameworks for user selection" but no TC covers this cold-start test-guide path. Add TC for FS-4 intermediate retry behavior (feed error back to LLM, regenerate).

2. **D2 Steps Actionability**: Replace "Inspect" steps with concrete `grep` commands. TC-001 step 5 should be `grep "ginkgo/v2" <path>` not "Inspect the generated test file". TC-008 step 4 should specify the log/output mechanism for verifying which Convention was loaded. TC-013 step 2 should use `cat .forge/config.yaml` or `grep` instead of "Inspect". TC-029 step 8 should specify concrete file path, not `<generated-test-file>`.

3. **D2/D3 Expected Results**: TC-006 "No errors or warnings related to legacy fields" needs concrete regex or assertion. TC-014 "Output lists all discovered tasks" needs concrete format. TC-020 needs concrete drift report format. TC-027 needs concrete assertion for "correctly interprets json-stream".

4. **D3 Convention Compliance**: TC-014 and TC-025 "No errors or warnings about missing Profile" are absence assertions that don't verify BIZ-error-reporting-002 patterns. Consider adding a positive assertion about what stdout DOES contain.

5. **D4 Boundary Cases**: Add TC for Convention file with all required sections missing. Add TC for justfile with broken e2e-compile recipe (exists but syntactically invalid).

6. **D4 Integration**: Add TC chaining init-justfile -> gen-test-scripts -> compile. Add TC chaining test-guide -> gen-test-scripts.

7. **D5 Manifest Update**: Update manifest.md to reflect all 29 TCs. Current manifest is stale (shows 20).

8. **D6 Meta-testing**: Convert TC-012 to a runtime behavioral test — instead of `grep -r "pkg/profile"`, test that a forge command that previously depended on Profile still works correctly (which implicitly verifies removal).

9. **Blindspot TC-009 Non-determinism**: Add a deterministic mechanism for "last-loaded wins" (e.g., alphabetical ordering assumption documented, or test both outcomes as acceptable).

10. **Blindspot TC-015 Annotation**: Add a Note to TC-015 clarifying it invokes a Claude Code slash command, similar to TC-003/004 annotations.

11. **Blindspot TC-027 Circular Pre-condition**: Remove "Generated tests produce json-stream output" from pre-conditions — it's established by step 2.
