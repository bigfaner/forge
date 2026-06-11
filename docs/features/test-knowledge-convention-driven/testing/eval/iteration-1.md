# Evaluation Report: CLI Test Cases — Iteration 1

**Document**: `cli-test-cases.md`
**Rubric**: `cli-test-cases.md` (1000 pts)
**Scorer**: Senior QA Engineer (adversarial)
**Date**: 2026-05-20

---

## Overall Score: 698 / 1000

---

## Dimension Scores

### D1: PRD Traceability — 170 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC-to-AC mapping exists | 60/70 | Most TCs trace to specific AC level (e.g., "Story 1 / AC-1"). Some Sources lack AC granularity: TC-015 "Spec FS-7 / Related Changes" points to a section, not an acceptance criterion. TC-011 "Spec FS-1 / Validation Rules" is section-level, not AC-level. |
| Traceability table complete | 65/70 | All 20 TCs listed. Missing traceability for FS-4 compile gate retry mechanism intermediate behavior (retry + error feedback to LLM). Only TC-010 covers the exhausted-retry terminal state. |
| Reverse coverage | 45/60 | All Story ACs covered. Gaps: (1) FS-4 retry mechanism (max 2 retries + feed error back to LLM) has no independent TC — only the exhaustion path is tested. (2) Related Changes #1 (pkg/journey/) and #4 (pkg/just/) lack dedicated TCs. (3) `run-e2e-tests` skill rewrite (listed in PRD Scope) has zero TCs. |

### D2: Step Actionability — 165 / 250 [BLOCKED — below 200 threshold]

| Criterion | Score | Notes |
|-----------|-------|-------|
| Steps are concrete actions | 55/90 | Multiple TCs use non-CLI steps: TC-001 step 5 "Inspect the generated test file for import statements" is file content inspection, not CLI invocation. TC-003 steps 3-4 ("Verify the skill scans..." / "Confirm the extracted patterns") are interactive skill operations, not concrete CLI calls. TC-008 steps 4-5 "Inspect which Convention file was loaded (via logs or generated code)" is ambiguous. TC-012 step 1 `grep -r "pkg/profile"` is a shell command, not a forge CLI command. TC-018 step 1 "Record the baseline generation time" is vague. |
| Expected results are verifiable | 60/90 | TC-007 "Exit code non-zero" — should specify exit code 1 per BIZ-error-reporting-001. TC-010 "Exit code indicates failure/blocked status" — unspecified concrete value. TC-003/004 "Skill output lists detected patterns/candidates" — no concrete output text or regex. TC-008 "Only testing-go.md Convention is loaded" — no observable CLI output specified. TC-009 "Last-loaded Convention wins" — no deterministic way to verify which was last-loaded. |
| Preconditions are explicit | 50/70 | TC-001 "A Go project exists using ginkgo" — no fixture creation method described. TC-006 "Forge has been upgraded to the Convention-based version" — no mechanism to simulate upgrade in test. TC-010 "The compile gate will fail consistently" — no explanation of how to guarantee deterministic failure. TC-018 "Baseline generation time with Profile is recorded" — contradicts TC-012's pre-condition that Profile is removed. |

**Verdict**: D2 = 165 < 200. Downstream gen-test-scripts is BLOCKED.

### D3: Command Coverage Accuracy — 95 / 150

| Criterion | Score | Notes |
|-----------|-------|-------|
| Command coverage | 30/50 | Missing TCs for: (1) `forge test detect/get/interfaces/framework` removal verification — PRD Scope explicitly lists these as commands to remove. (2) `forge task add` (PRD FS-7 Related Changes #5). (3) `forge init` (PRD FS-7). (4) `run-e2e-tests` skill (PRD Scope: "Rewrite run-e2e-tests skill"). |
| Output assertion specificity | 35/50 | TC-007: "Exit code non-zero" should be exit code 1. TC-010: "Exit code indicates failure/blocked status" unspecified. TC-002: "Output contains a warning listing the missing Framework section" — no concrete text or regex. TC-003: "Skill output lists detected patterns" — no concrete output. TC-008: "Only testing-go.md is loaded" — no observable output format. TC-011: "Output contains warning" — no concrete text. |
| Convention compliance | 30/50 | BIZ-error-reporting-001 violations: TC-007 and TC-010 do not specify concrete exit codes (should be 1 for failure). Deduct 10 pts each. BIZ-error-reporting-002 compliance is good for TC-007 and TC-010 (both include failure reason + recovery hint). |

### D4: Completeness — 120 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| Boundary and edge cases | 50/70 | Missing: (1) Convention file unreadable (permissions/encoding) — FS-2 Error Handling explicitly states this scenario. (2) Convention file with all required sections missing (not just Framework) — FS-1 validation rules cover this. (3) Cold start with file signals present but unrecognized patterns — FS-3 Failure Mode describes this. |
| Integration scenarios | 45/70 | No end-to-end flow TC (test-guide -> Convention -> gen-test-scripts -> compile -> run-e2e-tests). PRD "Main flow" is a key business flow with no integrated TC. No TC for gen-test-scripts -> run-e2e-tests integration (run-e2e-tests rewrite is in PRD Scope). |
| CLI coverage breadth | 25/60 | Missing commands: `forge task add`, `forge init`, `forge test detect/get/interfaces/framework` (removal verification), `run-e2e-tests` skill. |

### D5: Structure & ID Integrity — 100 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC IDs are sequential and unique | 40/40 | TC-001 through TC-020. Sequential, no gaps, no duplicates. |
| Classification is correct | 30/30 | All TCs classified as CLI type. |
| Summary table matches actual | 30/30 | Manifest reports 20 CLI TCs, actual count is 20. Traceability table lists all 20. |

### D6: Antipattern Prevention — 48 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| Pre-conditions are concrete and creatable | 15/30 | TC-006 "Forge upgraded" — no fixture creation method. TC-010 "compile gate will fail consistently" — no deterministic mechanism. TC-012 "pkg/profile/ removed" — assumes completed refactoring. TC-018 "Profile baseline recorded" — contradicts Profile removal. TC-019 "126+ Journeys" — no fixture scaling strategy. |
| Steps describe runtime behavior | 11/25 | TC-001 step 5, TC-008 steps 4-5, TC-009 step 6, TC-017 step 6: file content inspection, not runtime CLI behavior. TC-012 step 1: `grep` is source code search. TC-015 step 3: "Inspect the generated justfile". TC-003 step 3: "Verify the skill scans test files". |
| No duplicate scenarios | 20/20 | TC-005 and TC-016 have similar Expected but different preconditions (no test files vs has test files). Not duplicates. |
| No meta-testing | 0/15 | TC-012 step 3 "Run existing test suite" verifies test infrastructure, not product behavior. TC-019 "Run gen-test-scripts for all 126+ Journeys" is a benchmark/statistical test, not a functional test case. Both are meta-tests. |
| Every TC is implementable | 2/10 | TC-018 requires Profile-based baseline but Profile is removed (per TC-012). TC-019 requires 126+ Journey project scale. Neither is annotated as requiring special infrastructure. |

---

## Blindspot Attacks

1. **[blindspot] Timeline contradiction between TC-012 and TC-018**: TC-012 pre-condition states "The `pkg/profile/` directory has been removed". TC-018 step 1 requires "Record the baseline generation time for a representative set of Journeys using Profile-based generation". These two TCs assume mutually exclusive system states. If Profile is removed (TC-012, P0), the baseline measurement in TC-018 cannot be performed. The test suite contains an internal logical inconsistency that makes TC-018 impossible to execute in the same test run as TC-012.

2. **[blindspot] TC-003 and TC-004 test interactive multi-turn skills as CLI commands**: TC-003 step 4 says "Confirm the extracted patterns when presented". TC-004 step 4 says "Select one or more languages when prompted". PRD FS-5 defines test-guide as a "Multi-turn conversational skill". These TCs describe interactive human-in-the-loop actions, not automatable CLI invocations. Downstream gen-test-scripts cannot generate executable test scripts for interactive skill sessions. These TCs are functional acceptance tests masquerading as CLI test cases.

3. **[blindspot] TC-003/004 use slash commands, not forge CLI binary commands**: TC-003 step 2 "Run `/forge:test-guide` interactively" invokes a Claude Code slash command (skill dispatch), not the `forge` Go binary. Other TCs use `forge gen-test-scripts`, `forge task index`, etc. The test suite mixes two fundamentally different invocation mechanisms. Slash commands run within Claude Code's agent context; forge CLI commands are standalone executables. This distinction matters because gen-test-scripts generates test code that executes forge CLI binary commands, not Claude Code skill invocations.

---

## Required Improvements

1. **D2 Step Actionability (BLOCKING)**: Rewrite TC-003 and TC-004 to either (a) convert to functional acceptance tests in a separate document, or (b) define concrete CLI-equivalent steps that can verify test-guide behavior through observable file system changes (Convention file creation, output artifacts) rather than interactive prompts.

2. **D2 Expected Results**: Specify concrete exit codes for TC-007 (exit code 1) and TC-010 (exit code 1). Add concrete output text or regex patterns for TC-002, TC-003, TC-004, TC-008, TC-009, TC-011 warnings/hints.

3. **D2/D6 TC-018 Contradiction**: Resolve the timeline contradiction with TC-012. Either (a) make TC-018 a pre-migration benchmark that runs before Profile removal, or (b) use an archived baseline value as a hardcoded expected threshold, or (c) remove TC-018 and track performance as a separate non-test-case metric.

4. **D3 Command Coverage**: Add TCs for: `forge test detect` removal verification (should return error/exit 2), `forge test get` removal, `forge test interfaces` removal, `forge test framework` removal, `forge task add`, `forge init`, and `run-e2e-tests` skill behavior.

5. **D4 Boundary Cases**: Add TC for Convention file unreadable (permissions/encoding) per FS-2. Add TC for Convention file with all required sections missing per FS-1.

6. **D4 Integration**: Add at least one end-to-end flow TC covering: test-guide creates Convention -> gen-test-scripts uses Convention -> compile passes -> run-e2e-tests parses results.

7. **D6 Meta-testing**: Remove or restructure TC-012 step 3 ("Run existing test suite") and TC-019 — these verify test infrastructure quality, not product behavior. TC-019 should be a separate benchmark report, not a test case.

8. **D3 Convention Compliance**: Align all error-expecting TCs with BIZ-error-reporting-001 — specify exit code 1 for failure, exit code 2 for usage errors.
