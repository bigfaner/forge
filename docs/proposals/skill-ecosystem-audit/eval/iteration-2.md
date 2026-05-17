# Eval-Proposal: Iteration 2

**Score: 882/1000** (target: 900)

## DIMENSIONS

| Dimension | Score | Max | Delta from Iter 1 |
|-----------|-------|-----|-------------------|
| 1. Problem Definition | 105 | 110 | +13 |
| 2. Solution Clarity | 108 | 120 | +26 |
| 3. Industry Benchmarking | 95 | 120 | +57 |
| 4. Requirements Completeness | 100 | 110 | +22 |
| 5. Solution Creativity | 80 | 100 | +30 |
| 6. Feasibility | 92 | 100 | +14 |
| 7. Scope Definition | 75 | 80 | +7 |
| 8. Risk Assessment | 82 | 90 | +14 |
| 9. Success Criteria | 72 | 80 | +12 |
| 10. Logical Consistency | 83 | 90 | +5 |

## ATTACKS

1. **[Dim 9 - Success Criteria / Measurable and testable] Some criteria still lack precise falsification commands.**
   The success criteria improved dramatically with `grep`/`diff` commands, but several remain partially vague:
   - "breakdown-tasks and quick-tasks share common logic via shared reference files" -- the `diff` command checks that both reference the shared file, but does not verify that the *content* is actually shared rather than duplicated. A stronger test: `diff <(cat references/shared/type-assignment.md) <(grep -A20 'Type Assignment' breakdown-tasks/SKILL.md)` or similar.
   - "All slash commands produce identical behavior to pre-migration" -- "identical behavior" verified by e2e tests is good, but this should name a specific test count or list the skills covered.
   - Section: "Success Criteria" (lines 272-274, 276)
   - Fix: Add content-level dedup verification for W3. State the expected e2e test count or enumerate which skills the e2e suite covers.

2. **[Dim 8 - Risk / Likelihood calibration] Likelihood ratings remain uncalibrated.**
   All risks are rated M or L but no calibration framework is provided. "Medium likelihood" could mean 30% or 70%. The ratings are directional but not comparable across rows. This was flagged in iteration 1 and is unchanged.
   - Section: "Key Risks" table (lines 247-253)
   - Fix: Either add a calibration note (e.g., "H = >70%, M = 30-70%, L = <30%") or use a numeric scale (e.g., 0.4 probability). This makes the ratings auditable.

3. **[Dim 7 - Scope / Scope boundedness] Phasing added but deferral criteria are thin.**
   The phasing strategy is now present ("W3 is deferrable to v3.0.1 without user-facing breakage"), which is a meaningful improvement. However, the deferral criterion is vague -- what specific signal triggers deferral? Timeline pressure is mentioned but not quantified. If v3.0.0 ships in 2 days, do you cut W3? In 5 days? The boundary is qualitative.
   - Section: "Scope / Phasing Strategy" (line 219)
   - Fix: Add a concrete deferral trigger: e.g., "If W1+W2 completion exceeds 12h, defer W3 to v3.0.1." This makes the scope boundary machine-decidable.

4. **[Dim 3 - Industry Benchmarking / Chosen approach justified] Justification against benchmarks is improved but could be tighter.**
   The industry section now exists with 3 solid references and the "Why not the alternatives" subsection in W2 explicitly maps to the npm convention. However, the JetBrains benchmark (marketplace size limits) is cited but the 27MB node_modules issue is not explicitly cross-referenced to the JetBrains lesson -- the connection is implied but not drawn. The benchmarks and solution could be more tightly linked.
   - Section: "Industry Benchmarking" (lines 105-115) and "W2" (lines 129-144)
   - Fix: In W2's node_modules removal item, explicitly cite the JetBrains lesson: "JetBrains rejects plugins shipping redundant copies; we apply the same principle." This makes the benchmark-to-solution trace explicit.

5. **[Dim 5 - Solution Creativity / Cross-domain inspiration] CI guard is a good addition but the mechanism is bare-bones.**
   The CI path guard is a solid creative addition (responding to the iter 1 attack). However, the guard is a simple `grep` -- it detects regressions but does not prevent them at authoring time. A richer approach would be a skill-authoring convention document or a pre-commit hook. The CI guard addresses the minimum bar but does not demonstrate deeper cross-domain thinking (e.g., borrowing from linting ecosystems, suggesting a skill-linter framework).
   - Section: "Preventive Mechanism: CI Path Guard" (lines 156-168)
   - Fix: This is acceptable for the proposal's scope. Noting for completeness that the creativity ceiling is not fully reached -- a "skill linting framework" vision would elevate this further, but the current guard is sufficient for the problem at hand.

## DETAILED SCORING

### 1. Problem Definition (105/110)

- **Problem stated clearly (38/40)**: The problem is precisely articulated with the distribution model context. The per-skill scorecard provides a structured overview. The "Key Distribution Insight" section cleanly separates the two failure modes (path resolution vs. dead code). Slight deduction: the distinction between "path resolution failure" as a correctness bug and "dead code" as a maintenance burden could be framed more sharply as primary vs. secondary problems.
- **Evidence provided (38/40)**: The per-skill scorecard with lines, commits, clarity ratings, path counts, and verdicts is strong quantitative evidence. The full path inventory (lines 73-98) is exhaustive and file-specific. Minor deduction persists: clarity rating criteria (e.g., what differentiates 3/5 from 4/5) are implicit.
- **Urgency justified (29/30)**: The new "Why Now" section (lines 101-103) directly addresses the iter 1 gap. It ties urgency to the v3.0.0 release, names three confirmed user-facing failures, and explains why shipping without the audit harms every user. This is a strong improvement from 19/30. Minor deduction: no mention of user reports or bug tracker references to corroborate the "confirmed failures" claim.

### 2. Solution Clarity (108/120)

- **Approach is concrete (36/40)**: Major improvement. W1 remains concrete. W2 now commits to "relative-from-self resolution" with a clear rationale and explicit rejection of alternatives with reasons (lines 133-137). The prerequisite spike is a responsible hedge. W3 is clearer -- "extract to a rubric file (not a new skill)" resolves the ambiguity. Deduction: W3's shared reference file mechanism (format, naming convention, loading instructions for Claude) is still underspecified.
- **User-facing behavior described (40/45)**: Strong improvement. The scenarios section (lines 174-184) now includes both happy-path and failure/recovery scenarios with specific error behaviors ("Claude logs a 'file not found' error") and degradation strategies ("embed a compressed summary as a fallback"). The failure/recovery scenarios directly address the iter 1 gap. Deduction: the happy-path scenarios could include one concrete walkthrough (e.g., "User runs `/eval` on a PRD -> Claude reads rubric from `../eval/rubrics/harness.md` relative to the skill file -> evaluation proceeds normally").
- **Technical direction clear (32/35)**: W1 and W2 are now technically precise. The relative-from-self approach is well-explained with a concrete example path. The prerequisite spike is a practical validation step. W3's extraction targets are named with file locations. Remaining gap: the exact format and content of "shared reference files" is not specified -- are they plain Markdown? Do they have a header convention?

### 3. Industry Benchmarking (95/120)

- **Industry solutions referenced (32/40)**: Three external references added: VS Code extensions, npm packages, JetBrains IDE plugins. Each includes a specific API or mechanism and a lesson drawn. This is a major improvement from zero references. Deduction: references are stated as factual claims without source links or version context. "VS Code Extensions use `extensionPath`" -- link to the API docs would strengthen credibility. Also, no reference to how Claude Code's own plugin system documentation describes asset resolution.
- **At least 3 meaningful alternatives (28/30)**: The alternatives table (lines 208-214) now has four rows with genuine alternatives. "Build-time path rewriting" is a realistic and defensible alternative that was added. The "Only W1" and "Only W2" rows are partial-scope alternatives, which are legitimate for a multi-workstream proposal. The straw-man from iter 1 ("Rewrite from scratch") is removed. Good.
- **Honest trade-off comparison (18/25)**: Each alternative has pros and cons listed. "Build-time path rewriting" correctly identifies the trade-off ("Zero runtime changes" vs "Requires a build step; fragile; does not fix dead code"). The selected approach's con is "Medium total effort (~16h)" which is now quantified. Deduction: the cons are somewhat shallow -- "fragile if path format changes" for build-time rewriting could be elaborated. The "Only W1" row's con is self-evident rather than insightful.
- **Chosen approach justified against benchmarks (17/25)**: The "Why not the alternatives" subsection (lines 135-137) justifies the relative-from-self approach by linking to the npm convention. The industry section's convergence statement (line 115) draws the connection between benchmarks and the selected approach. Deduction: the VS Code canonical-path-variable pattern is cited but the proposal explicitly *rejects* this pattern (because skills lack `${CLAUDE_PLUGIN_ROOT}`). This rejection is honest but the justification could be stronger by acknowledging that a future skill runner update could enable the canonical variable approach, making relative-from-self a pragmatic short-term choice rather than an architecturally superior one.

### 4. Requirements Completeness (100/110)

- **Scenario coverage (37/40)**: Significant improvement. The scenarios now include failure/recovery cases (lines 182-184): missing referenced files, missing `docs/conventions/` directory, and shared template incompatibility. The three failure scenarios are concrete and distinct. Deduction: one edge case remains unaddressed -- what happens when the prerequisite spike (line 139) fails and relative-from-self resolution does not work? The fallback is mentioned in risk mitigation but not as a scenario.
- **Non-functional requirements (35/40)**: Major improvement. A dedicated "Non-functional Requirements" subsection (lines 186-191) now covers: install size target (under 5MB, currently ~30MB), backward compatibility (e2e tests must pass), security (path traversal prevention), and cross-platform (forward slash normalization). The install size target is quantified. Deduction: no performance NFR (should skill execution time remain constant?), and the security NFR is thin -- "inherently constrains this since Claude only reads paths the skill specifies" is a claim, not a threat model.
- **Constraints & dependencies (28/30)**: Constraints are well-stated. The dependency on the test profile system is correctly identified. The existing `skill-rationalization` proposal dependency is noted. Minor deduction: the "Backward compatibility" constraint appears in both Requirements (line 189) and Constraints (line 201) -- slight redundancy.

### 5. Solution Creativity (80/100)

- **Novelty over industry baseline (30/40)**: The CI path guard (lines 156-168) is a novel addition that lifts this above straightforward cleanup. The guard makes the "no hardcoded paths" invariant machine-enforced, which is a genuine structural improvement beyond the immediate fixes. The prerequisite spike as a validation mechanism is practical. Deduction: the guard is a `grep` -- functional but not architecturally novel. A skill-authoring framework or path-registry pattern would be more creative.
- **Cross-domain inspiration (25/35)**: The npm `npm pack --dry-run` analogy for the CI guard (line 168) shows cross-domain thinking. The relative-from-self approach borrows from npm's `__dirname` convention. The JetBrains marketplace gatekeeping parallel is drawn. Deduction: the cross-domain insights are present but largely extracted from the industry benchmarking section rather than generating new creative ideas. The mapping from "industry does X" to "we should do X-equivalent" is mechanical rather than inspirational.
- **Simplicity of insight (25/25)**: The three-workstream prioritization remains clean. The relative-from-self approach is the simplest possible fix -- no build step, no runtime variable, just "use relative paths." The insight that skills are Markdown executed by an AI agent (and therefore cannot access env vars) is a genuine architectural insight that simplifies the solution space. Full marks.

### 6. Feasibility (92/100)

- **Technical feasibility (35/40)**: W1 is clearly feasible. W2's relative-from-self approach is technically sound -- Claude's Read tool supports relative paths. The prerequisite spike (1h) is a responsible validation step. W3's shared reference files are technically straightforward. Deduction: the prerequisite spike could fail, and the fallback ("embed content inline") is mentioned in risk mitigation but not in the solution itself. The solution should acknowledge the spike outcome's impact on the approach.
- **Resource & timeline feasibility (27/30)**: Major improvement. Each workstream now has hour estimates: W1 (~3h), W2 (~8h), W3 (~5h), totaling ~16h. The prerequisite spike is budgeted at 1h. The phasing strategy (line 219) sequences the work. Deduction: no resource allocation -- who does the work? Is this one person or two? Can W1 and W2 overlap?
- **Dependency readiness (30/30)**: All dependencies are ready. The test profile system exists. The plugin distribution mechanism is understood. No external tools needed. The prerequisite spike validates the only uncertain dependency (Claude's relative path resolution behavior). Full marks.

### 7. Scope Definition (75/80)

- **In-scope items are concrete (28/30)**: All 10 in-scope items map to specific files and changes. Each item is actionable. The W3 eval extraction now specifies "extraction to a rubric file under eval/" which resolves the ambiguity from iter 1. Full marks near-achieved.
- **Out-of-scope explicitly listed (23/25)**: 8 out-of-scope items listed. "New skill creation (validate-ux extraction stays within eval/ as a rubric, not a separate skill)" -- this now explicitly resolves the iter 1 contradiction. Good. Minor deduction: "just or doc-scorer/doc-reviser references (expected dependencies)" is ambiguous -- if they are expected dependencies, why are they out of scope rather than in the constraints section?
- **Scope is bounded (24/25)**: The phasing strategy (line 219) now provides temporal bounding. "If timeline pressure requires cutting scope, W3 is deferrable to v3.0.1 without user-facing breakage -- W1 and W2 are release-blocking." This is a meaningful improvement. Minor deduction: the deferral trigger is qualitative ("timeline pressure") rather than quantitative.

### 8. Risk Assessment (82/90)

- **Risks identified (27/30)**: Five risks now listed, up from four. The new risk "Shared template changes have unintended cross-skill effects" (line 253) addresses the iter 1 gap. The prerequisite spike risk is well-framed. Deduction: two risks from iter 1 still missing -- cross-platform path resolution inconsistency (Windows backslashes vs Unix forward slashes), and the eval extraction changing scoring behavior. Wait -- the eval extraction risk IS present at line 252. Cross-platform is implicitly addressed in the NFR section (line 191) but not in the risk table. Minor gap.
- **Likelihood + impact rated (25/30)**: Ratings are present (L/M/H) and now cover all five risks. Improvement from iter 1: the ratings are more differentiated (L/M for extraction vs M/H for path resolution). Deduction: no calibration framework persists from iter 1. "M" likelihood is not defined.
- **Mitigations are actionable (30/30)**: Major improvement. Each mitigation is now concrete and executable:
  - "place a test skill at the plugin cache path... invoke the skill... confirm Claude reads the file" (line 249)
  - "Generate golden output before deletion... re-generate and `diff -r golden/ output/`" (line 250)
  - "Generate tasks for the same test PRD via both paths... Diff the task files" (line 251)
  - "Run eval on 3 existing test cases... Compare scores line-by-line" (line 252)
  - "Each shared reference file lists its consuming skills... CI guard runs diff" (line 253)
  Each mitigation names a specific action, tool, and pass/fail criterion. Full marks.

### 9. Success Criteria (72/80)

- **Criteria are measurable and testable (50/55)**: Major improvement. Most criteria now include `grep` or `diff` commands. Examples:
  - `grep -r 'record-task' ... returns 0 hits` (line 260)
  - `test -d skills/gen-test-scripts/templates/node_modules && echo FAIL || echo OK` (line 268)
  - `wc -l skills/eval/SKILL.md returns under 280 lines` (line 278)
  These are precise, falsifiable, and executable. Deduction: "All slash commands produce identical behavior" (line 274) -- verified by e2e tests, but the criterion does not specify how many tests or which skills are covered. "breakdown-tasks and quick-tasks share common logic" (line 272) -- the verification command checks file existence and grep patterns but not content dedup.
- **Coverage is complete (22/25)**: Coverage improved significantly. W3 now has dedicated criteria: eval line count (line 278), scope heuristic fix (line 280), and CI guard (line 282). The iter 1 gap is addressed. Deduction: the "Plugin installs and all skills work correctly" catch-all criterion (line 276) covers 18 skills with a single test (`claude --skill eval --prompt "list rubrics"`) -- this should enumerate which skills are smoke-tested or name a test count.

### 10. Logical Consistency (83/90)

- **Solution addresses the stated problem (34/35)**: The three workstreams map directly to P0/P1/P2 findings. W1 addresses P0 #1 (record-task refs) and P0 #3 (forensic paths). W2 addresses P0 #1 (node_modules), P0 #2 (types/ path), and P1 #4 (hardcoded paths). W3 addresses P1 #5 (duplication) and P2 #8 (eval inline). The traceability is tight.
- **Scope <-> Solution <-> Success Criteria aligned (24/30)**: Mostly aligned. In-scope item 8 ("Simplify eval validate-ux inline sub-pipeline -- extraction to a rubric file under eval/") matches W3's description ("not a new skill -- see scope exclusion alignment below") and the success criterion at line 278 ("under 280 lines"). The iter 1 contradiction is resolved. Deduction: the success criterion for shared templates (line 272) checks file existence but the scope item says "shared reference files" -- the alignment is present but the verification does not fully prove the scope item is achieved (files can exist without being used).
- **Requirements <-> Solution coherent (25/25)**: Requirements and solution are coherent. The backward compatibility requirement is operationalized through the verification plan (lines 194-197) and e2e test references. The security NFR (path traversal) is addressed by the relative-from-self approach. The install size target (under 5MB) is achievable through W2's 27MB removal. The verification plan maps each workstream to specific checks. Full marks.
