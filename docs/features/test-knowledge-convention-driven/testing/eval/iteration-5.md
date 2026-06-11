# Evaluation Report: CLI Test Cases — Iteration 5

**Document**: `cli-test-cases.md`
**Rubric**: `cli-test-cases.md` (1000 pts)
**Scorer**: Senior QA Engineer (adversarial)
**Date**: 2026-05-20

---

## Overall Score: 901 / 1000

---

## Dimension Scores

### D1: PRD Traceability — 186 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC-to-AC mapping exists | 67/70 | All 36 TCs trace to specific PRD sources. TC-036 added since iteration 4 maps to "PRD / Main flow + Bootstrap flow (FS-5 -> FS-2 -> FS-4)" — a valid multi-spec source. Remaining imprecision: TC-015 Source "Spec FS-7 / Related #4" — section-level but specific enough (references pkg/just/). TC-026 Source "Spec FS-7 / Related #5" — same pattern, references internal/cmd/. TC-027 Source "PRD Scope / Rewrite skill" — scope-level, not AC-level (should be "PRD Scope / In Scope — Rewrite run-e2e-tests skill"). These three are imprecise but functional. Deduct 3 pts. |
| Traceability table complete | 70/70 | All 36 TCs listed in the Traceability table. Table matches actual TC count (36). TC IDs sequential TC-001 through TC-036. No TCs missing from the table. |
| Reverse coverage | 49/60 | Improvement from iteration 4: TC-036 now covers the test-guide -> gen-test-scripts pipeline (FS-5 -> FS-2 -> FS-4), addressing the integration gap. Remaining gaps: (1) FS-5 "If existing Convention file found -> show diff, ask user to confirm update or keep" — no TC covers this update interaction path. TC-003 covers creation, TC-030 covers cold start, but neither covers the update path. (2) FS-7 Related Changes #1 (pkg/journey/ — "Remove FrameworkInfo dependency; Convention-driven tags") — no TC directly verifies journey package CLI behavior change. (3) FS-9 drift detection — TC-020 tests it but Expected regex `Convention.*assert.*conflict.*test.*require|drift.*Convention:.*assert.*actual:.*require` is well-specified now; however TC-020 step 3 still uses `forge consolidate-specs` (with parenthetical "(or the equivalent CLI invocation for drift detection mode)") which is hedged rather than definitive. (4) PRD "Rewrite gen-test-scripts skill" in Scope list — covered by TC-001, TC-005, TC-008, TC-019 but no TC explicitly verifies the "Rewrite" aspect (Convention + Code Reconnaissance + compile gate working together as a complete unit) beyond TC-029 (which is close but focuses on cold start). Deduct 11 pts for 2 uncovered AC paths and 2 partially-covered spec items. |

### D2: Step Actionability — 228 / 250

| Criterion | Score | Notes |
|-----------|-------|-------|
| Steps are concrete actions | 82/90 | Improvement: TC-036 steps are mostly concrete CLI actions (ls, test, grep, forge gen-test-scripts, just e2e-compile). Step 2 "Invoke test-guide skill and confirm framework detection (or simulate by creating the Convention file...)" is acceptable with the "or simulate" alternative providing a concrete action. Remaining issues from iteration 4 that persist: TC-002 step 5 "Run `grep -c "testing\|testify\|assert" <generated-test-path>` to verify generated test file uses LLM-detected defaults" — this IS concrete (uses grep), so the previous iteration's complaint about "Inspect" was partially resolved. Correction: TC-002 step 5 is now concrete. TC-027 step 5 "Verify the parsed results report matches the actual test outcomes" — still vague, "verify...matches" is not a concrete CLI action. TC-030 step 4 "Verify the skill presents candidate frameworks for user selection (e.g., 'go-testing', 'ginkgo')" — "Verify" is not a concrete action. TC-031 step 4 "Capture full stdout/stderr output including retry attempts" — capture mechanism unspecified. Deduct 8 pts for 4 vague steps (2 pts each). |
| Expected results are verifiable | 78/90 | TC-036 Expected is verifiable: specific exit codes, grep counts, concrete assertions. Improvement from iteration 4. Remaining issues: TC-008 Expected step 4 "testing-go.md Convention is loaded (grep finds reference in logs or generated code)" — the "or" between logs vs generated code creates ambiguity. TC-020 Expected "Drift report flags the assert/require mismatch" and "Report includes the Convention file path and the conflicting test file" — now has concrete regex in Expected: `Convention.*assert.*conflict.*test.*require|drift.*Convention:.*assert.*actual:.*require` — this is improved. Checking actual TC-020 Expected text: "Output matches regex: `Convention.*assert.*conflict.*test.*require|drift.*Convention:.*assert.*actual:.*require`" — yes, this IS concrete. Correction: TC-020 Expected is now verifiable. Remaining: TC-030 step 4 Expected "Skill output includes candidate framework names (not auto-selecting a single framework)" — "not auto-selecting" is an absence assertion, not objectively verifiable. TC-002 Expected step 5 "grep count >= 1 (generated test file uses Go testing + testify defaults)" — this IS verifiable. Correction: TC-002 is fine. Deduct 12 pts for TC-008 "or" ambiguity (-3), TC-030 absence assertion (-3), TC-027 step 5 "verify...matches" (-3), TC-031 step 4 capture unspecified (-3). |
| Preconditions are explicit | 68/70 | All 36 TCs declare pre-conditions. TC-036 pre-condition is explicit: "standard Go project using go-testing + testify with existing test files... No Convention files exist... valid Journey... justfile contains e2e-compile. Forge binary is built from the Convention-based source." Slight issue: TC-036 pre-condition "Forge binary is built from the Convention-based source" — while descriptive, it assumes a build step that is not reproducible from the TC itself (same issue as TC-006, TC-012). Deduct 2 pts for implicit build assumption. |

**Verdict**: D2 = 228 >= 200. Downstream gen-test-scripts is NOT blocked.

### D3: Command Coverage Accuracy — 140 / 150

| Criterion | Score | Notes |
|-----------|-------|-------|
| Command coverage | 47/50 | All major CLI commands from PRD Scope covered. TC-036 adds integration coverage. Missing: (1) `forge gen-test-scripts` Journey targeting syntax — all TCs say "Run `forge gen-test-scripts` targeting the Journey" but never specify the concrete CLI invocation syntax for Journey selection (positional arg? `--journey <id>` flag?). This was flagged in iteration 4 and remains unfixed. (2) `forge config init` in an existing project (TC-013 only tests empty directory). (3) `forge init` idempotency in a non-empty directory. Deduct 3 pts. |
| Output assertion specificity | 46/50 | TC-020 now has concrete regex for drift output — improvement. TC-036 has concrete grep counts. Remaining: TC-008 Expected step 4 "or" ambiguity (logs vs generated code). TC-030 Expected "not auto-selecting" absence assertion. TC-020 step 3 `forge consolidate-specs` invocation is hedged. Deduct 4 pts for 2 ambiguous output assertions and 1 hedged command. |
| Convention compliance | 47/50 | BIZ-error-reporting-001: TC-007 (exit 1), TC-010 (exit 1), TC-033 (exit 1) correctly specify failure exit codes. TC-021-024 correctly specify exit code 2 for usage errors. BIZ-error-reporting-002: TC-007, TC-010, TC-033 include recovery hints. TC-036 does not specify error behavior (no negative path). Issue: TC-020 does not assert exit code — only output regex. Is drift detection success (exit 0) or failure (exit 1) when drift is found? The TC does not say. Deduct 3 pts for missing exit code assertion on TC-020. |

### D4: Completeness — 188 / 200

| Criterion | Score | Notes |
|-----------|-------|-------|
| Boundary and edge cases | 68/70 | Comprehensive coverage: TC-011 (missing domains frontmatter), TC-028 (permissions), TC-032 (all sections missing), TC-034 (invalid encoding), TC-035 (invalid section content), TC-033 (broken recipe). Remaining gap: Convention file with valid structure but contradictory content across sections (e.g., Framework says ginkgo but Assertion lists testify-specific functions). PRD FS-1 validation rules do not explicitly mention cross-section consistency validation, so this is not required. Near-complete. Deduct 2 pts for minor gap. |
| Integration scenarios | 67/70 | Major improvement: TC-036 now chains test-guide -> gen-test-scripts -> compile in a single pipeline, directly addressing the iteration 4 gap. TC-029 covers cold start -> compile -> Convention -> gen -> compile. TC-031 chains gen-test-scripts with compile gate retry feedback. Remaining gaps: (1) No TC chains init-justfile -> gen-test-scripts (init-justfile creates justfile, gen-test-scripts uses it). TC-015 tests init-justfile and TC-001 tests gen-test-scripts, but no single TC verifies the chain. (2) No TC chains gen-test-scripts -> run-e2e-tests (TC-027 generates scripts within its steps but does not verify that gen-test-scripts output is directly consumable by run-e2e-tests). Deduct 3 pts for 2 missing chain scenarios. |
| CLI coverage breadth | 53/60 | All major CLI commands from PRD Scope have at least one TC. TC-036 adds integration breadth. Remaining gaps: (1) `forge config init` only TC-013 (happy path, empty directory). No TC for existing project migration. (2) `forge task index` only TC-014 (basic invocation). No TC for flags or filters. (3) No TC for `forge gen-test-scripts` Journey selection argument syntax. (4) FS-5 test-guide "If existing Convention file found -> show diff, ask user to confirm update or keep" — no TC for this update path. Deduct 7 pts. |

### D5: Structure & ID Integrity — 92 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| TC IDs are sequential and unique | 40/40 | TC-001 through TC-036. Sequential, no gaps, no duplicates. |
| Classification is correct | 28/30 | All TCs classified as CLI type. However: TC-003, TC-004, TC-030 are marked as CLI type but their Note fields acknowledge they "require a Claude Code agent session" and "cannot be automated via forge binary alone." These are conversational skill interactions, not pure CLI commands. While the rubric says "No UI, API, TUI, or Mobile TCs," these TCs are CLI-adjacent but not purely CLI (they test agent-driven behavior). This is a classification borderline — the TCs verify file-system artifacts produced by CLI-adjacent skills. Deduct 2 pts for borderline classification without a "hybrid" or "agent-assisted" label. |
| Summary table matches actual | 24/30 | **Manifest is now updated**: manifest.md reports "36 | Draft" which matches the actual 36 TCs. The frontmatter description says "36 test cases" — matches. The Traceability table lists 36 entries — matches. **However, frontmatter is missing the `sources` field**: rubric requires "Frontmatter with `feature`, `sources`, `generated`" but the actual frontmatter has only `feature`, `type`, `generated` — no `sources` field. This is a structural requirement gap. Deduct 6 pts for missing required frontmatter field. |

### D6: Antipattern Prevention — 67 / 100

| Criterion | Score | Notes |
|-----------|-------|-------|
| Pre-conditions are concrete and creatable | 20/30 | TC-036 pre-condition "Forge binary is built from the Convention-based source" — not creatable in isolated test fixture, requires building from a specific branch state. Same issue as TC-006 ("Forge has been upgraded") and TC-012 ("pkg/profile/ removed, binary rebuilt"). TC-010 "The compile gate will fail consistently" — relies on Convention declaring wrong framework; non-deterministic because LLM might still produce compilable code. TC-018 "at least 10 Journeys" — requires substantial fixture. TC-036 step 2 "or simulate by creating the Convention file" — the "simulate" alternative IS creatable (create the file manually), which partially addresses the issue. Deduct 10 pts for 3 non-creatable pre-conditions (TC-006, TC-010, TC-012) and 1 partially non-creatable (TC-018). |
| Steps describe runtime behavior | 18/25 | TC-036 steps 1, 3-7 are runtime CLI commands. Step 2 "Invoke test-guide skill and confirm framework detection (or simulate by creating the Convention file...)" — the "simulate" path is a static file creation step, not runtime CLI interaction. The "invoke test-guide" path is runtime but requires agent session. Remaining from iteration 4: TC-027 step 5 "Verify the parsed results report matches the actual test outcomes" — manual comparison. TC-030 step 4 "Verify the skill presents candidate frameworks" — verification of conversational output. TC-031 step 4 "Capture full stdout/stderr output" — mechanism unspecified. Deduct 7 pts for 3-4 static/ambiguous steps. |
| No duplicate scenarios | 20/20 | No duplicate TCs identified. TC-003 and TC-036 both involve test-guide Convention creation, but TC-003 tests Convention creation in isolation while TC-036 tests the full pipeline (creation -> consumption -> compilation). Different scenarios. TC-029 and TC-036 both test end-to-end flows, but TC-029 starts from cold start (no Convention) while TC-036 starts from test-guide creation — different starting states. |
| No meta-testing | 5/15 | TC-012 persists from iteration 4: its intent is structural verification (import audit gate, FS-7), not product behavior. The TC tests CLI runtime commands but its purpose is verifying code structure (pkg/profile removal). TC-036 is a legitimate integration TC (product behavior: test-guide -> gen-test-scripts -> compile). Improvement: no new meta-tests added. Deduct 10 pts for TC-012 meta-testing. |
| Every TC is implementable | 4/10 | TC-006 "Forge has been upgraded" — requires specific binary build. TC-010 "compile gate will fail consistently" — non-deterministic. TC-012 "pkg/profile/ removed, binary rebuilt" — requires completed development work. TC-018 "at least 10 Journeys" — substantial fixture requirement. TC-036 "Forge binary is built from the Convention-based source" — same binary build assumption. These TCs have Note annotations (TC-003, TC-004, TC-015, TC-027, TC-030, TC-036) acknowledging agent-session requirements, which partially addresses implementability. However TC-006, TC-010, TC-012 lack equivalent annotations for their build assumptions. Deduct 6 pts for TCs with substantial fixture/build requirements. |

---

## Blindspot Attacks

1. **[D5] Frontmatter missing `sources` field**: The rubric Required Sections explicitly state "Frontmatter with `feature`, `sources`, `generated`". The actual frontmatter contains `feature: "test-knowledge-convention-driven"`, `type: CLI`, `generated: "2026-05-20"` but NO `sources` field. This is a structural requirement that has been missing since iteration 1 and has not been flagged in previous evaluations. Quote from rubric: "Frontmatter with `feature`, `sources`, `generated`". The frontmatter should contain something like `sources: ["prd-spec.md", "prd-user-stories.md"]`.

2. **[D2] TC-002 step 5 is concrete but TC-002 step 2 "Create Convention file" is vague**: TC-002 step 2 says "Create a valid Journey and Contract spec" without specifying HOW to create them. The pre-conditions say "The project has a valid Journey with Contract specs" which implies they should already exist, but step 2 creates them. This is a pre-condition/steps inconsistency — the pre-condition says "has a valid Journey" but step 2 says "Create a valid Journey." One of them should be authoritative.

3. **[D1] FS-5 update path still uncovered**: PRD FS-5 Convention creation flow step 5 states: "If existing Convention file found -> show diff, ask user to confirm update or keep." This is a distinct interaction path from creation (TC-003) and cold-start (TC-030). No TC covers this update path. TC-036 creates a Convention from scratch (no existing file). The update path requires an existing Convention file and verifies diff/update behavior.

4. **[D4] TC-036 step 2 "or simulate" weakens integration value**: TC-036's value as an integration TC is undermined by the "or simulate" alternative in step 2. If the tester simulates by creating the Convention file manually, the TC becomes functionally identical to TC-001 (gen-test-scripts uses existing Convention) with a preceding file-creation step. The integration value comes from the test-guide -> gen-test-scripts chain, but the "simulate" path bypasses test-guide entirely. The TC should either (a) require test-guide invocation without a simulate alternative, or (b) be annotated as "partial integration" with the test-guide step being optional.

5. **[D6] TC-012 meta-testing persists**: TC-012's intent is to verify the import audit gate (FS-7) — a structural code verification. The TC title is "Forge Commands Function Correctly After Profile Removal" which sounds behavioral, but the pre-condition "The `pkg/profile/` directory has been removed. Forge binary is rebuilt and installed" and the Expected "No errors or warnings about missing Profile" reveal the structural verification purpose. A downstream agent executing this TC would need to have already completed the Profile removal development work before testing — this is a development gate, not a test case.

6. **[D3] TC-020 exit code unspecified for drift found scenario**: TC-020 Expected does not specify what exit code means "drift was found." Is exit 0 "drift detected and reported" (successful detection) or exit 1 "drift detected" (failure state)? The TC pre-condition sets up a drift scenario (Convention says assert, tests use require) and step 3 runs `forge consolidate-specs`. The Expected says "Output matches regex" but no exit code is specified. Per BIZ-error-reporting-001, the exit code semantics should be explicit.

---

## Changes Since Iteration 4

| Issue from Iteration 4 | Status | Notes |
|------------------------|--------|-------|
| Update manifest.md to reflect 35 TCs | Fixed | manifest.md now shows 36 TCs (updated for TC-036) |
| Add integration TC for test-guide -> gen-test-scripts | Fixed | TC-036 added |
| TC-002 step 5 "Inspect" replaced with grep | Partially fixed | Step 5 now uses `grep -c "testing\|testify\|assert"` but step 2 has pre-condition/steps inconsistency |
| TC-020 Expected drift format | Fixed | Now has concrete regex `Convention.*assert.*conflict.*test.*require\|drift.*Convention:.*assert.*actual:.*require` |
| TC-009 non-determinism | Not fixed | Accepts either outcome, cannot verify merge behavior |
| TC-006, TC-010 non-creatable pre-conditions | Not fixed | Same issues persist |
| TC-012 meta-testing | Not fixed | Still present |
| TC-030 automatability | Not fixed | Step 4 "Verify" still not automatable |
| TC-033 error path clarification | Not fixed | Conflates broken recipe with missing recipe |

---

## Required Improvements

1. **D5 Add `sources` field to frontmatter**: The rubric requires `sources` in frontmatter. Add `sources: ["prd-spec.md", "prd-user-stories.md"]` (or equivalent) to the YAML frontmatter.

2. **D1 Add TC for FS-5 Convention update path**: PRD FS-5 step 5 "If existing Convention file found -> show diff, ask user to confirm update or keep" has no TC. Add a TC where test-guide is invoked on a project that already has a Convention file.

3. **D3 Specify exit code for TC-020**: TC-020 must specify what exit code means "drift detected." Is it success (exit 0 = drift detection worked) or failure (exit 1 = drift found)?

4. **D6 Resolve TC-012 meta-testing**: Either remove TC-012 (import audit is a development gate, not a test case) or restructure it to test a genuine product behavior (e.g., "forge commands produce correct output when Profile-dependent code paths are exercised").

5. **D2 Fix TC-031 step 4**: Replace "Capture full stdout/stderr output including retry attempts" with a concrete mechanism, e.g., `forge gen-test-scripts 2>&1 | tee output.log`.

6. **D6 Add annotations to non-implementable TCs**: TC-006, TC-010, TC-012, TC-018 should have explicit annotations about their build/fixture requirements, similar to how TC-003, TC-015, TC-027 acknowledge agent-session requirements.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1: PRD Traceability | 186 | 200 |
| D2: Step Actionability | 228 | 250 |
| D3: Command Coverage Accuracy | 140 | 150 |
| D4: Completeness | 188 | 200 |
| D5: Structure & ID Integrity | 92 | 100 |
| D6: Antipattern Prevention | 67 | 100 |
| **Total** | **901** | **1000** |

**Gate Status**: PASS (target: 900, current: 901). Above target by 1 pt. Primary improvements since iteration 4: manifest updated, TC-036 integration TC added, TC-020 drift regex made concrete. Remaining weaknesses concentrated in D6 (antipattern prevention — TC-012 meta-testing, non-creatable pre-conditions, non-implementable TCs) and D1 (FS-5 update path uncovered).
