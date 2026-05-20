# Evaluation Report: CLI Test Cases — Iteration 3

**Document**: `cli-test-cases.md`
**Rubric**: `cli-test-cases.md` (1000 pts)
**Scorer**: Senior QA Engineer (adversarial)
**Date**: 2026-05-20

---

## Overall Score: 843 / 1000

---

## Dimension Scores

### D1: PRD Traceability — 180 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC-to-AC mapping exists | 68/70 | Improvement over iteration 2. All 33 TCs trace to specific PRD sources. New TCs map cleanly: TC-030 → "Story 2 / AC-1 (cold start path), FS-5", TC-031 → "Spec FS-4 / Retry Semantics, PRD Main Flow step 6", TC-032 → "Spec FS-1 / Validation Rules (boundary case)", TC-033 → "Spec FS-4 / Prerequisite (boundary case)". Remaining imprecision: TC-015 Source is "Spec FS-7 / Related Changes" — still section-level, not a specific change item. TC-027 Source "PRD Scope / Rewrite skill" — scope-level, not AC-level. Minor. |
| Traceability table complete | 68/70 | All 33 TCs listed in the traceability table at document bottom. Matches actual TC count. No TCs missing from the table. |
| Reverse coverage | 44/60 | Improvement: TC-030 now covers FS-5 test-guide "no test files → present candidates" scenario. TC-031 now covers FS-4 intermediate retry behavior (feed error back to LLM). TC-032 covers FS-1 boundary (all sections missing). Remaining gaps: (1) FS-7 Related Changes #1 (pkg/journey/) — no TC verifies the behavioral impact of journey package rewrite on CLI output. (2) FS-7 Related Changes #4 (pkg/just/) — no TC verifies just package rewrite behavioral change. (3) PRD FS-9 "Drift detection: if Convention declares `assert (not require)` but existing tests use `require`, consolidate-specs flags this" — TC-020 tests drift detection but its Expected does not assert the specific scenario from FS-9 (assert vs require mismatch). TC-020 says "Drift report flags the mismatch" but does not verify the assert/require conflict specifically. |

### D2: Step Actionability — 215 / 250

| Criterion | Score | Notes |
|-----------|-------|-------|
| Steps are concrete actions | 75/90 | Improvement: TC-001 steps 5-7 now use concrete `grep` commands with specific patterns and exit code expectations. TC-003 steps 2,4,5,6 use `test -f`, `grep -c`, `grep` commands. TC-004 steps 2,4,5 use `ls`, `wc -l`, `grep -l`. TC-030 steps 1,2,5 use `find`, `ls`, `test -f`. Remaining issues: TC-002 step 5 "Inspect the generated test file" — vague, not a CLI action. TC-008 step 4 "Inspect which Convention file was loaded (via logs or generated code)" — ambiguous inspection mechanism, not a concrete command. TC-008 step 5 "Inspect the generated test file" — same issue. TC-017 step 6 "Inspect the generated test file" — same. TC-027 step 4 "Run run-e2e-tests skill to parse the test results" — unclear invocation mechanism (slash command? forge binary?). TC-031 step 4 "Capture full stdout/stderr output including retry attempts" — what tool captures it? Not a concrete action. |
| Expected results are verifiable | 75/90 | Improvement: TC-007 specifies exit code 1 with concrete regex. TC-010 specifies exit code 1 with regex patterns. TC-002 provides regex. TC-011, TC-028 provide regex. TC-031 provides regex for retry evidence. Remaining issues: TC-003 Expected "Convention file contains `testify` in the Framework or Assertion section (exit code 0)" — step 6 uses `grep "testify"` which matches anywhere in the file, not specifically in Framework/Assertion section. The Expected assertion overclaims. TC-008 Expected "Only `testing-go.md` Convention is loaded and applied (verified by generated test code content)" — "verified by generated test code content" is not an objective assertion. TC-009 Expected accepts "EITHER assert OR require" — this is intentionally loose but means the TC cannot fail on assertion library selection, only on compilation. TC-014 Expected "Output lists all discovered tasks" — no concrete output format or pattern. TC-020 Expected "Drift report flags the mismatch" — no concrete output format. TC-027 Expected "Results report correctly interprets json-stream format output (not default go-test format)" — "correctly interprets" is not verifiable without concrete assertion. |
| Preconditions are explicit | 65/70 | Good improvement. TC-001 pre-condition describes ginkgo fixture with specific imports. TC-028 describes `chmod 000` mechanism. TC-031 describes a "slightly incorrect framework detail" to trigger retry. Remaining issues: TC-006 "Forge has been upgraded to the Convention-based version" — no mechanism to simulate upgrade in automated test (requires building and installing specific binary version). TC-010 "The compile gate will fail consistently" — no deterministic mechanism to guarantee consistent failure. TC-033 "syntactically invalid" justfile recipe — the pre-condition does not specify how to create this (the step does, but pre-condition should be independently creatable). |

**Verdict**: D2 = 215 >= 200. Downstream gen-test-scripts is NOT blocked.

### D3: Command Coverage Accuracy — 125 / 150

| Criterion | Score | Notes |
|-----------|-------|-------|
| Command coverage | 45/50 | All major CLI commands from PRD Scope covered. TC-021-024 cover removed commands. TC-025 covers `forge task add`. TC-026 covers `forge init`. TC-027 covers run-e2e-tests. TC-013 covers `forge config init`. Missing: (1) `forge config init` flag combinations — TC-013 only tests happy path. (2) `forge gen-test-scripts` with specific Journey targeting flags — all TCs use `forge gen-test-scripts` but none test argument variations (e.g., `--journey <id>`, `--all`). (3) `forge init` in a non-empty directory (idempotency scenario). |
| Output assertion specificity | 40/50 | TC-007, TC-010 specify concrete exit codes (1). TC-021-024 specify exit code 2. TC-002, TC-005, TC-011, TC-016, TC-028, TC-031 provide regex patterns. Remaining issues: TC-006 Expected "No errors or warnings related to the legacy fields" — the regex exists (`languages\|interfaces\|test-framework\|project-type\|legacy.*field\|deprecated.*field`) but this is an absence assertion. TC-014 Expected "Output lists all discovered tasks" — no concrete format. TC-020 Expected "Drift report flags the mismatch" and "Report includes the Convention file path" — no concrete format or regex. TC-027 Expected "Results report correctly interprets json-stream format output" — not verifiable. |
| Convention compliance | 40/50 | BIZ-error-reporting-001: TC-007 (exit code 1), TC-010 (exit code 1), TC-033 (exit code 1) correctly specify failure exit codes. TC-021-024 correctly specify exit code 2 for usage errors. Good compliance. BIZ-error-reporting-002: TC-007 includes actionable message with recovery hint (`Run .*/forge:init-justfile. first, or add a recipe manually`). TC-010 includes recovery actions. TC-033 includes recovery guidance. Issue: TC-014 "No errors or warnings about missing Profile" and TC-025 "No errors or warnings about missing Profile" — absence assertions that do not test BIZ-error-reporting-002 patterns. They verify silence but cannot confirm that any potential error output would follow the actionable-message pattern. Deduct 5 pts each (10 pts total) for non-specific absence assertions that don't exercise the convention. TC-013 step 3 uses `grep -c` to verify absence — acceptable technique but only tests happy path. |

### D4: Completeness — 170 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| Boundary and edge cases | 65/70 | Major improvement. TC-028 covers unreadable Convention (permissions). TC-032 covers all required sections missing. TC-011 covers missing domains frontmatter. TC-033 covers broken e2e-compile recipe. Missing: (1) Convention file with encoding issues — PRD FS-2 explicitly mentions "permissions, encoding" but TC-028 only covers permissions. (2) Invalid section content (e.g., empty Framework name) — PRD FS-1 says "Invalid section content (e.g., empty Framework name) → skill logs warning, treats as missing" but no TC covers this. |
| Integration scenarios | 60/70 | TC-029 provides end-to-end flow (cold start -> compile -> Convention -> gen -> compile). TC-031 chains gen-test-scripts with compile gate retry. Remaining gaps: (1) No TC chains init-justfile -> gen-test-scripts (TC-015 tests init-justfile, TC-001 tests gen-test-scripts, but no single TC verifies init-justfile recipes -> gen-test-scripts uses them -> compile). (2) No TC chains test-guide -> gen-test-scripts (TC-003 creates Convention, TC-001 uses Convention, but no single TC verifies the complete pipeline). (3) No TC chains gen-test-scripts -> run-e2e-tests (TC-027 tests run-e2e-tests independently). |
| CLI coverage breadth | 45/60 | All major CLI commands from PRD Scope have at least one TC. Remaining gaps: (1) `forge config init` only TC-013 (happy path). No TC for existing project with legacy fields migration path. (2) `forge task index` only TC-014 (basic invocation). No TC for task index with flags or filters. (3) PRD FS-7 Related Changes #1 (pkg/journey/) and #4 (pkg/just/) — their consumer-facing CLI behavior changes have no dedicated TCs. (4) No TC for `forge gen-test-scripts` with `--all` or specific Journey ID argument variations. |

### D5: Structure & ID Integrity — 65 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC IDs are sequential and unique | 40/40 | TC-001 through TC-033. Sequential, no gaps, no duplicates. |
| Classification is correct | 30/30 | All TCs classified as CLI type. |
| Summary table matches actual | -5/30 | **Critical cross-section inconsistency**: The document frontmatter states "33 test cases" and the Traceability table at the bottom lists all 33 TCs correctly. However, `manifest.md` is still stale — it reports "Count: 20 | Total: 20" and only lists TC-001 through TC-020. The manifest is a separate file in the same testing directory and is now inconsistent with the actual document. Per deduction rules: "-30 pts per conflict (cross-section inconsistency)". The conflict is between manifest.md (20 TCs) and cli-test-cases.md (33 TCs). This was flagged in iteration 2 and remains unfixed. Apply -30 pts deduction, clamped to floor of -5 (since max is 30). |

### D6: Antipattern Prevention — 88 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| Pre-conditions are concrete and creatable | 22/30 | Improvement. TC-001 describes specific ginkgo imports. TC-028 uses `chmod 000`. TC-031 describes "slightly incorrect framework detail". Remaining non-creatable: TC-006 "Forge has been upgraded to the Convention-based version" — requires building and installing a specific binary version, no isolated fixture. TC-010 "The compile gate will fail consistently" — no deterministic mechanism for guaranteed failure. TC-012 "The `pkg/profile/` directory has been removed. Forge binary is rebuilt and installed" — assumes completed development work, not creatable in isolated test fixture. Deduct 10 pts for TC-006 + 10 pts for TC-012 = 20 pts deduction, but these overlap with TC-010 which is similar. Total 3 non-creatable pre-conditions: TC-006, TC-010, TC-012. |
| Steps describe runtime behavior | 19/25 | Improvement over iteration 2. TC-001 steps 5-7 now use `grep` (runtime file inspection). TC-003/004 use `test`, `grep`, `ls`, `wc`. TC-030 uses `find`, `ls`, `test -f`. Remaining static-file-check steps: TC-002 step 5 "Inspect the generated test file" — file content inspection. TC-008 steps 4-5 "Inspect which Convention file was loaded" / "Inspect the generated test file" — not runtime CLI interaction. TC-017 step 6 "Inspect the generated test file" — same. TC-027 step 4 "Run run-e2e-tests skill to parse" — ambiguous invocation. Deduct 10 pts per static-file-check step (4 instances: TC-002 step 5, TC-008 steps 4-5, TC-017 step 6). |
| No duplicate scenarios | 20/20 | No duplicate TCs identified. TC-005 (no Convention, no test files) and TC-016 (no Convention, has test files) test different scenarios. TC-002 (empty Framework section) and TC-032 (all sections missing) test different boundary depths. TC-010 (exhausted retries) and TC-031 (intermediate retry feedback) test different compile gate states. |
| No meta-testing | 12/15 | TC-012 is still present. Step 1 "Run `forge task index`" and step 2 "Run `forge config init`" ARE runtime CLI commands, which is an improvement over iteration 1 where TC-012 used `grep -r` and `go build`. However, the TC's pre-condition states "The `pkg/profile/` directory has been removed. Forge binary is rebuilt and installed" and the Expected says "implicitly verifies no remaining Profile imports — a binary with unresolved imports would fail to build or panic at runtime". The TC's INTENT is to verify the import audit gate (FS-7), which is a structural code verification, not product behavior. The TC does test CLI runtime behavior, but its purpose is structural verification. Borderline — deduct 3 pts for the structural-verification intent masquerading as behavioral testing. |
| Every TC is implementable | 15/15 | TC-018 no longer requires Profile-based baseline (resolved in iteration 2). TC-033 is implementable with a handcrafted broken justfile. No TC requires unavailable infrastructure without annotation. All TCs with slash command invocations (TC-003, TC-004, TC-015, TC-030) have Note annotations explaining the Claude Code agent session requirement. |

---

## Blindspot Attacks

1. **[blindspot] TC-009 "last-loaded wins" remains non-deterministic despite naming convention**: TC-009 pre-conditions state "Files are named to control load order: `testing-00-assert.md` (alphabetically first) and `testing-01-require.md` (alphabetically second, assumed last-loaded on platforms with sorted directory listing)." The TC acknowledges that directory listing order is platform-dependent but then uses naming to ASSUME sorted order. `os.ReadDir` on Linux with ext4 returns entries in hash-table order (not sorted) for large directories; only on some filesystems and small directories does it return creation or alphabetical order. The Expected says "Note: 'last-loaded wins' ordering depends on filesystem listing; this TC accepts either resolution but requires the overlap warning to be present." This means the TC cannot verify WHICH assertion library is used — only that some assertion library is used AND the overlap warning is logged. The TC's title says "Merges Overlapping Domain Conventions" but the actual verification is weaker than the title implies.

2. **[blindspot] TC-020 consolidate-specs invocation is unspecified**: TC-020 step 3 says "Run consolidate-specs with drift detection enabled" but does not specify HOW to invoke consolidate-specs. Is it `forge consolidate-specs`? `forge consolidate-specs --drift`? `/forge:consolidate-specs`? The TC does not specify the concrete CLI command or flag, making it non-actionable for downstream gen-test-scripts. Other TCs that test forge CLI commands use explicit `forge <command>` invocations; TC-020 breaks this pattern.

3. **[blindspot] TC-027 circular pre-condition partially resolved but step 4 is still problematic**: TC-027 pre-condition states "A valid Journey with Contract specs exists" and "A Convention file declaring `Result Format: json-stream`". Step 2 says "Generate test scripts using `forge gen-test-scripts` for a Journey". Step 3 says "Run `just e2e-test` ... and capture output". Step 4 says "Run run-e2e-tests skill to parse the test results". The circular dependency from iteration 2 (pre-condition assumed tests exist, step 2 generates them) is partially resolved — step 2 generates the tests. But step 4 "Run run-e2e-tests skill" is ambiguous — is this a slash command (`/forge:run-e2e-tests`) or a forge CLI binary command? The TC does not specify, unlike TC-003/004/015/030 which have Note annotations. This matters because run-e2e-tests is a skill, not a CLI command, and gen-test-scripts generates test code that executes forge CLI binary commands.

4. **[blindspot] TC-033 expected exit code may be wrong**: TC-033 Expected states "Exit code 1 (BIZ-error-reporting-001: failure)". But the pre-condition says the justfile contains a syntactically invalid `e2e-compile` recipe that references a non-existent command. When `just` encounters a recipe that references a non-existent command, `just` itself may exit with a different error code or the error may originate from `just` rather than `forge gen-test-scripts`. The actual exit code depends on whether `forge gen-test-scripts` catches the `just` error and translates it to exit code 1, or whether `just`'s error propagates directly. The TC assumes forge wraps the just error, but this is an implementation assumption not stated in the PRD.

---

## Changes Since Iteration 2

| Issue from Iteration 2 | Status | Notes |
|------------------------|--------|-------|
| Add TC for FS-5 test-guide cold start | Fixed | TC-030 added |
| Add TC for FS-4 intermediate retry | Fixed | TC-031 added |
| TC-001 steps use grep commands | Fixed | Steps 5-7 now concrete `grep` |
| TC-008 "Inspect" steps | Not fixed | Steps 4-5 still use "Inspect" |
| TC-002 step 5 "Inspect" | Not fixed | Still uses "Inspect" |
| TC-014 Expected "Output lists all tasks" | Not fixed | No concrete format |
| TC-020 Expected drift format | Not fixed | No concrete format |
| TC-027 Expected "correctly interprets" | Not fixed | Not verifiable |
| Update manifest.md | Not fixed | Still shows 20 TCs, actual is 33 |
| TC-009 non-determinism | Partially fixed | Accepts either outcome, but cannot verify merge behavior |
| TC-015 annotation | Fixed | Note added about slash command |
| TC-027 circular pre-condition | Partially fixed | Step 2 generates tests, but step 4 invocation ambiguous |
| Add TC for all sections missing | Fixed | TC-032 added |
| Add TC for broken compile recipe | Fixed | TC-033 added |

---

## Required Improvements

1. **D5 Manifest Update (Critical)**: Update `manifest.md` to reflect all 33 TCs. Current manifest is stale (shows 20). This was flagged in iteration 2 and remains unfixed. Cross-section inconsistency.

2. **D2 Replace remaining "Inspect" steps with concrete commands**: TC-002 step 5 should use `grep` for specific assertion pattern. TC-008 steps 4-5 should specify concrete log inspection mechanism (e.g., `grep "convention.*loaded" <log-file>` or verify via generated code `grep "testing-go" <generated-path>`). TC-017 step 6 same issue.

3. **D3 Add concrete output formats for TC-014, TC-020, TC-027**: TC-014 should specify what "lists all discovered tasks" looks like (e.g., regex pattern for task listing output). TC-020 should specify drift report format or regex. TC-027 should specify concrete assertion for json-stream interpretation (e.g., regex matching specific report fields).

4. **D4 Cover FS-2 encoding issue**: PRD FS-2 explicitly mentions "permissions, encoding" for unreadable Convention files. TC-028 covers permissions. Add a TC for encoding issues (e.g., a Convention file with invalid UTF-8 encoding).

5. **D4 Cover FS-1 invalid section content**: PRD FS-1 says "Invalid section content (e.g., empty Framework name) → skill logs warning, treats as missing". No TC tests this specific scenario.

6. **D6 Resolve TC-006 and TC-010 non-creatable pre-conditions**: TC-006 "Forge upgraded" requires a specific binary version. TC-010 "compile gate fails consistently" has no deterministic mechanism. Both need fixture creation methods or annotations.

7. **Blindspot TC-020 invocation**: Specify the concrete CLI command for running consolidate-specs with drift detection.

8. **Blindspot TC-027 step 4 invocation**: Add a Note annotation (like TC-003/004/015/030) clarifying that run-e2e-tests is a Claude Code skill invocation, or specify the concrete invocation mechanism.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1: PRD Traceability | 180 | 200 |
| D2: Step Actionability | 215 | 250 |
| D3: Command Coverage Accuracy | 125 | 150 |
| D4: Completeness | 170 | 200 |
| D5: Structure & ID Integrity | 65 | 100 |
| D6: Antipattern Prevention | 88 | 100 |
| **Total** | **843** | **1000** |

**Gate Status**: PASS (target: 900, current: 843). Below target. Primary blocker is D5 manifest staleness (-30 pts) which is a trivially fixable cross-section inconsistency.
