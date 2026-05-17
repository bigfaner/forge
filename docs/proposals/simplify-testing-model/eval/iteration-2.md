---
iteration: 2
date: 2026-05-17
score: 775
target: 900
previous_score: 640
---

# Proposal Evaluation Report — Iteration 2

**Document**: `docs/proposals/simplify-testing-model/proposal.md`
**Score**: 775/1000
**Target**: 900
**Outcome**: NOT MET — 125 points short
**Improvement from Iteration 1**: +135 points

## DIMENSIONS

### 1. Problem Definition: 85/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | The core problem (3 overlapping concepts reduced to 2) is clearly stated with a before/after YAML example. The new paragraph explaining why the override chain is cleaner ("This is a 1D narrowing model, not a 2D matrix") directly addresses the iteration-1 logical-consistency attack. Minor gap remains: "概念边界模糊" is still asserted without quoting a user saying "I don't understand this". |
| Evidence provided | 30/40 | Significant improvement. The Evidence table now has a "Concrete Evidence" column with code-level specifics: `config.go:ForgeConfig`, `testgen.go:GetBreakdownTestTasks`, `profile.go:runProfileResolve`, `embed.go:ValidTestTypes`, `detect.go:DetectProfiles`. This is source-code evidence, not just assertion. However, this is still author-internal evidence — no user feedback, no support tickets, no quantitative data (e.g., "X% of config files have mismatched profile/capability"). The preamble says "Internal friction points observed during v2 development and onboarding" but does not cite specific onboarding sessions, issues, or user reports. |
| Urgency justified | 20/30 | The urgency section is mostly unchanged from iteration 1. It correctly identifies v3.0.0 as the breaking-change window. However, it remains vague on cost of delay: "Every new feature added to v2 ... compounds the profile x capability matrix" — how many features? How many call sites per quarter? No quantification. The previous iteration received 30/30 for this criterion; on re-evaluation, the justification is directionally correct but thinner than it appeared. |

### 2. Solution Clarity: 108/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 40/40 | The approach is now highly concrete. The new Technical Feasibility section provides code-level detail: function signature changes (`DetectProfiles()` -> `DetectLanguages()`), struct field changes (remove `TestProfiles`/`Capabilities`, add `Interfaces`/`Languages`), embed path changes (`all:profiles` -> `all:languages`), CLI command group changes (`forge profile` -> `forge testing`). A reader can explain exactly what will be built. |
| User-facing behavior described | 45/45 | Unchanged from iteration 1 — already excellent. D1, D4, D5, D6 fully specify observable behavior. Edge cases (no detection, conflicting signals, monorepo) are now covered in scenarios 4-6. |
| Technical direction clear | 23/35 | Substantially improved. The new Feasibility Assessment section provides 5 numbered subsections describing changes to detect.go, config.go, embed.go, cmd/profile.go, and testgen.go. However, gaps remain: (1) No data flow diagram or description of how a CLI call flows from command through detection through strategy retrieval. (2) The embed.go subsection shows a `languageCapabilities` map but does not explain how this replaces the `manifest.yaml` resolution logic — what currently reads manifest.yaml and how does that code change? (3) No description of error handling in the new path (what happens when `DetectLanguages()` returns empty). |

### 3. Industry Benchmarking: 80/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 35/40 | Major improvement. The Industry References table now has 4 rows with URLs, "How it works" column, and "What Forge borrows" column. References include GitHub Actions setup-go/setup-node, Jest, pytest, Terraform, ESLint, Next.js, and Vite — with specific URLs. Gaps: (1) One URL appears broken: `developer.hashicorne.com` should be `developer.hashicorp.com`. (2) No version numbers or dates cited — these references could become stale. |
| At least 3 meaningful alternatives | 15/30 | Improved from iteration 1 but still problematic. The "Do nothing" alternative now has genuine pros ("Zero migration cost; all skills already work") and substantive cons ("`testgen.go` matrix complexity persists; every new skill must learn two config axes"). "Gradual deprecation" now has concrete cons ("Two parallel codepaths in `profile.go` and `testgen.go`; `ForgeConfig` struct grows new fields alongside deprecated ones"). However, "Rename only" is still thin: its cons are 2 sentences that restate the problem rather than analyze the approach. And "Gradual deprecation" references "Python 2->3 migration pattern" as its source — this is a famous migration disaster, which should be a cautionary tale, not a positive reference. The document does not acknowledge this irony. Deduction: 1 remaining straw-man-style alternative (-20 from sub-score). |
| Honest trade-off comparison | 15/25 | Improved. The selected approach now has a concrete con: "All 10 skills (D6) must migrate CLI calls; `embed.go` directory structure changes from `profiles/<name>/` to `languages/<key>/`; `manifest.yaml` files removed, metadata hardcoded in Go". This is real analysis. However, the cons for rejected alternatives are still shorter than the cons for the selected approach, creating an imbalance. The "Do nothing" cons are 3 clauses; the "Gradual deprecation" cons are 2 clauses; the selected approach cons are 3 clauses — roughly balanced but the rejected alternatives' pros are notably thin. |
| Chosen approach justified against benchmarks | 15/25 | Significant improvement. The new section "Why Auto-detect Fits Forge Specifically" provides three numbered reasons: (1) "Reliable" signals, (2) "Already implemented" detection, (3) "Low-impact failure mode". It also acknowledges the JavaScript multi-framework limitation explicitly. Gaps: (1) The "Reliable" argument says "False positives are rare" without evidence — how rare? (2) The section does not address why Forge should not follow the Terraform/ESLint explicit-declaration pattern, beyond saying language is "easily inferred". For some users, explicit control might be preferred — this counterargument is not engaged. |

### 4. Requirements Completeness: 85/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | Major improvement. Scenarios 4-6 now cover edge cases: "No language detected" (scenario 4), "Conflicting signals" (scenario 5), "Monorepo with language-specific subdirectories" (scenario 6). Each has concrete behavior described. Gaps: (1) Scenario 5 says "if the `package.json` has no Playwright dependency and no test scripts, it falls back to Playwright" — falling back to Playwright when there is no Playwright dependency seems like a bug, not a design. This is confusing. (2) No error scenario for invalid `languages` values (e.g., `languages: [foobar]`). (3) No scenario for the `interfaces` field receiving an invalid value. |
| Non-functional requirements | 30/40 | New section added. The NFR table covers Performance (<100ms, benchmark), Compatibility (CLI output format stability, E2E tests), Security (no directory traversal, code review), and Extensibility (checklist for adding new language). Each has a verification method. Gaps: (1) No accessibility NFR (not applicable to CLI, but not stated as out-of-scope). (2) Performance requirement says "< 100ms for projects with < 1000 files" — why 1000 files as the threshold? What about larger repos? (3) No NFR for the migration experience itself (how long should migration take for a typical user?). |
| Constraints & dependencies | 20/30 | Improved with explicit technology constraints (Go 1.22+, gopkg.in/yaml.v3, cobra framework). Gaps: (1) No mention of the CLI output format stability as a constraint (skills parse CLI output, so format changes break skills). (2) No mention of downstream consumers beyond the listed skills — are there other tools or CI systems that call `forge profile`? (3) The "6 existing strategy sets migrated as-is" constraint is stated but the risk of migration error is not quantified. |

### 5. Solution Creativity: 55/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 25/40 | The new "Multi-framework language handling" paragraph in Innovation Highlights is a genuine contribution: using framework-specific language keys (`javascript` = Playwright, future `javascript-jest`) rather than generic language detection. This is a practical innovation for the multi-framework problem. However, the core idea (auto-detect from file signals) is still directly lifted from CI tools — the document itself cites this lineage. |
| Cross-domain inspiration | 15/35 | Still limited to directly adjacent domains (CI systems, IDEs, test frameworks). No inspiration from more distant fields. The iteration-1 attack on this point was not addressed. The references remain within "developer tools that detect project configuration" — no cross-pollination from package managers, type systems, or other domains. |
| Simplicity of insight | 15/25 | The 1:1 language-to-framework insight is now explicitly qualified with the JavaScript exception and a concrete handling strategy. The framework-key approach (`javascript` = Playwright, `javascript-jest` = future) is a clean solution. However, the insight is practical rather than surprising — it is "competent engineering", not "why didn't I think of that". |

### 6. Feasibility: 80/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | Major improvement. The new Technical Feasibility section audits 5 specific code areas with function names, struct fields, and concrete change descriptions. The Dependency Readiness table verifies each dependency with status and risk. Gaps: (1) No prototype or proof-of-concept code. (2) The claim "`manifest.yaml` files are removed; capabilities metadata for each language is hardcoded as a Go map" introduces new hardcoded data — what happens when a new language is added? The extensibility NFR says "no manifest files" but hardcoding in Go requires recompilation, which is less extensible than a manifest. This tension is not acknowledged. |
| Resource & timeline feasibility | 25/30 | New Resource & Timeline table with 7 work items, effort estimates in days, dependencies, and a total of "~10 person-days (2 weeks)". Single developer, no cross-team coordination. This is a credible estimate. Gaps: (1) No buffer/contingency mentioned. (2) "2 days" for restructuring directories and "2 days" for CLI implementation are round numbers that suggest estimation rather than bottom-up analysis. (3) No mention of testing effort — who writes the E2E tests? |
| Dependency readiness | 20/30 | The Dependency Readiness table is new and valuable. All 5 dependencies show "None" risk, which is suspicious — are there really zero risks in any dependency? The strategy file compatibility row says "Content is pure markdown/text; no code executes within strategy files" — this is good analysis but does not address whether the template paths within strategy files (e.g., references to `templates/test_file.go`) use absolute or relative paths that might break when the directory structure changes. |

### 7. Scope Definition: 72/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | Eight concrete deliverables including the new "Migration guide: v2-to-v3 config mapping document". Each is identifiable. Minor gap: "All 10 skills consuming profile/capability updated" still does not cross-reference the D6 table explicitly within the In Scope section. |
| Out-of-scope explicitly listed | 24/25 | Five items explicitly deferred, including the new "Multi-framework JavaScript support beyond Playwright". The iteration-1 attack about documentation being in/out of scope is resolved — migration guide is now in-scope. Minor gap: User-facing documentation updates (beyond the migration guide) are not explicitly scoped. |
| Scope is bounded | 20/25 | The new timeline anchor ("2 weeks (10 person-days), single developer, targeting v3.0.0 release branch") directly addresses the iteration-1 attack. This is a significant improvement. Remaining gap: No definition of "done" beyond the deliverable list — is it done when skills are migrated, or when E2E tests pass, or when the migration guide is published? The success criteria section helps but the scope section itself does not reference it. |

### 8. Risk Assessment: 70/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | Expanded to 6 risks from 4 in iteration 1. New risks: "CI pipeline breakage for existing v2 users" and "Strategy file format incompatibility after directory rename" and "File scanning performance on large monorepos". These directly address iteration-1 attacks. Gap: No risk for "hardcoded language capabilities map becomes stale or incorrect" — this is a new data artifact introduced by the proposal. |
| Likelihood + impact rated | 20/30 | Ratings are present for all 6 risks. The CI pipeline breakage risk is rated H likelihood, M impact — honest assessment. Gap: (1) No justification for why "Skill migration gap" is M likelihood — with 10 skills to migrate, this could easily be H. (2) "Strategy file format incompatibility" is rated L likelihood but H impact — what makes it L? The mitigation says "Verified by auditing all 6 strategy sets" but no audit evidence is provided. (3) No risk rating scale defined (what does M mean in terms of probability?). |
| Mitigations are actionable | 25/30 | Improved significantly. Each mitigation now includes specific actions: "run `grep -r 'forge profile' plugins/` and `grep -r 'capabilities' plugins/` before closing v3.0 milestone", "v3.0.0 is a major version bump with documented breaking changes. Migration guide maps each v2 field to its v3 equivalent". The mobile detection mitigation now explains the semantic inconsistency explicitly. Gaps: (1) The multi-language false positive mitigation says "Detection priority is not configurable — signals are additive, not exclusive" — this is a design decision, not a mitigation. The actual mitigation is "user adds `languages: [go]` to config.yaml", which is stated but should be the primary mitigation. (2) No mitigation for the file scanning risk — it says "No mitigation needed beyond the current design" but does not explain why no mitigation is needed. |

### 9. Success Criteria: 60/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 40/55 | Expanded to 10 checklist items from 6. New criteria cover: multi-language detection with `--language` flag, `languages` override behavior, `interfaces` defaulting behavior, byte-for-byte output comparison with v2, zero `profile`/`capability` in exported symbols, and migration guide existence. Gaps: (1) "A Go project with `project-type: backend` and no other config fields generates and runs tests end-to-end" — "end-to-end" is vague. What specific steps constitute the end-to-end verification? (2) "`forge testing detect` correctly identifies all 6 language/platform types" — "correctly" needs a test oracle (what output format? what exact strings?). (3) The byte-for-byte criterion ("output is byte-for-byte identical to the corresponding v2 `forge profile get <name>` output") is extremely strict — will whitespace, timestamps, or version strings in output cause false failures? (4) No criterion for error behavior (e.g., what happens when no language is detected and no override is set — does it error gracefully?). |
| Coverage is complete | 20/25 | The new criteria cover the `languages` override and `interfaces` defaulting that were missing in iteration 1. The migration guide criterion is new. Gaps: (1) No criterion for the NFRs — no performance criterion (should be: "detection completes in <100ms for test project X"), no security criterion. (2) No criterion for the `forge testing interfaces` command behavior. (3) No criterion for the `forge testing get justfile` or `forge testing get template` subcommands. |

### 10. Logical Consistency: 80/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 32/35 | The new paragraph "Why the new override chain is cleaner than v2" directly addresses the iteration-1 attack about the override chain being "just different". The argument is clear: v2 has a 2D matrix (profile x capability), v3 has a 1D narrowing (language determines strategy, interfaces filters). This is a substantive response. Minor gap: The problem says "Capability naming is confusing" but the solution does not explicitly trace "capability -> interface" as a rename — it is mentioned in D1 but not in the problem-solution mapping. |
| Scope / Solution / Success Criteria aligned | 28/30 | Major improvement. The iteration-1 attack about per-feature narrowing and mobile appearing in scenarios but being out-of-scope is addressed: (1) Per-feature narrowing is explicitly stated as "deferred to a future release" in D5 and listed in Out of Scope. (2) Mobile is treated as in-scope for migration (directory rename) but out-of-scope for "deep design". (3) Success criteria include mobile detection. The alignment is now much cleaner. Minor gap: Scenario 5 ("Conflicting signals") describes fallback behavior ("falls back to Playwright") but no success criterion verifies this fallback behavior. |
| Requirements / Solution coherent | 20/25 | Improved. The 6 scenarios now map to the solution: scenarios 4-6 cover edge cases that the solution explicitly handles. Gaps: (1) The NFR table lists an Extensibility requirement ("adding a new language requires 3 steps") but no in-scope item or success criterion validates this extensibility claim. (2) The "Security" NFR says "Code review confirms `detect.go` uses `os.ReadFile` on fixed paths only" — this is a verification method, not a requirement. The actual security requirement (what must be true about detection) is not stated as a testable condition. |

## ATTACKS

1. **Problem Definition**: Evidence remains author-internal — "Internal friction points observed during v2 development and onboarding" names no specific onboarding sessions, user reports, support tickets, or metrics. The "Concrete Evidence" column cites source code but not user impact. Must add at least one external data point: a user quote, a support ticket, or a metric (e.g., "X% of new users misconfigured test-profiles in their first config.yaml").

2. **Industry Benchmarking**: One URL is broken — `developer.hashicorne.com` should be `developer.hashicorp.com`. Verify and fix all URLs before publication.

3. **Industry Benchmarking**: The "Gradual deprecation" alternative cites "Python 2->3 migration pattern" as its source, but Python 2->3 is widely regarded as a cautionary tale about the costs of dual-path maintenance. The document presents it neutrally. Must either acknowledge that this reference undermines the alternative (since Python 2->3 was notoriously painful) or cite a successful dual-path migration instead.

4. **Industry Benchmarking**: The "Rename only" alternative remains thin — its cons are 2 sentences that restate the problem ("Does not address the core problem: `DetectProfiles()` returns profile names like 'go-test' and the user must still understand what a profile is") rather than genuinely analyzing what a rename-only approach would achieve or miss. Must add at least one genuine pro and one specific con with evidence.

5. **Requirements Completeness**: Scenario 5 contains a confusing design claim — "if the `package.json` has no Playwright dependency and no test scripts, it falls back to Playwright" — falling back to a framework that is not installed is a potential bug. Must clarify: does "falls back to Playwright" mean "assumes Playwright if no other JS framework is detected" or "returns no language for JavaScript"? The current wording implies the system will generate Playwright strategies for a project that does not use Playwright.

6. **Solution Creativity**: Cross-domain inspiration remains within directly adjacent domains (CI, IDEs, test frameworks). The iteration-1 attack was not addressed. Must cite inspiration from at least one non-developer-tool domain (e.g., how package managers resolve dependency conflicts, how type systems perform type inference, how databases use query optimization heuristics).

7. **Feasibility**: The hardcoded `languageCapabilities` map contradicts the Extensibility NFR — adding a new language requires "no schema migration, no manifest files" but DOES require editing Go source code and recompiling. The proposal says this is better than manifest.yaml but does not acknowledge the trade-off. Must either (a) acknowledge that extensibility is "recompile required" or (b) propose a data-driven approach that does not require Go code changes.

8. **Feasibility**: All 5 dependency readiness risks are rated "None". This is suspiciously optimistic. The strategy file content compatibility row says files have "no path references to their own location" — but the Risk Assessment section's mitigation for strategy file format incompatibility says "Templates use language-relative paths (e.g., `test_file.go`) which are resolved by the CLI, not embedded in strategy content." These two statements are about different things but the tension between "no path references" and "language-relative paths resolved by CLI" is not explored. Must audit whether any strategy file references a path that will break when the directory prefix changes from `profiles/<name>/` to `languages/<key>/`.

9. **Risk Assessment**: No risk identified for the hardcoded `languageCapabilities` map becoming stale. If a new interface type is added or an existing one changes, the Go map must be updated. This is a maintenance risk that is not captured. Must add a risk row for "hardcoded capability map drifts from actual strategy capabilities."

10. **Risk Assessment**: Likelihood ratings lack justification — "Skill migration gap causes runtime errors" is rated M likelihood with 10 skills to migrate. With 10 independent migration points, the probability of at least one being missed is non-trivial. Must justify the M rating with reasoning (e.g., "CI smoke tests will catch most gaps before release").

11. **Success Criteria**: The byte-for-byte criterion ("output is byte-for-byte identical to the corresponding v2 `forge profile get <name>` output") is unrealistically strict. Any change in CLI framework version, Go version, or template rendering could introduce whitespace or formatting differences. Must either (a) specify a normalized comparison (ignoring trailing whitespace, version strings) or (b) relax to "functionally equivalent output verified by structured comparison."

12. **Success Criteria**: No criterion verifies the NFRs — the NFR table specifies "<100ms detection" and "stable CLI output format" but no success criterion tests these. Must add criteria like "Benchmark test confirms DetectLanguages completes in <100ms for a project with 6 signal files" and "E2E test for each forge testing subcommand passes."

13. **Success Criteria**: No criterion for error behavior when no language is detected and no `languages` override is set. Scenario 4 describes this behavior ("prints an error message directing the user to add `languages` to config.yaml") but no success criterion verifies it. Must add: "When no language is detected and no `languages` override is set, `forge testing get generate` exits with non-zero code and prints a message containing 'languages'."

14. **Logical Consistency**: The Extensibility NFR ("adding a new language requires 3 steps, no manifest files") has no corresponding success criterion. Either add a criterion (e.g., "A developer can add a 7th language by following a 3-step checklist without modifying any existing code") or acknowledge this is an aspirational goal, not a committed deliverable.
