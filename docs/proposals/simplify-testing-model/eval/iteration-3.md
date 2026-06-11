---
iteration: 3
date: 2026-05-17
score: 835
target: 900
previous_score: 775
---

# Proposal Evaluation Report — Iteration 3

**Document**: `docs/proposals/simplify-testing-model/proposal.md`
**Score**: 835/1000
**Target**: 900
**Outcome**: NOT MET — 65 points short
**Improvement from Iteration 2**: +60 points

## DIMENSIONS

### 1. Problem Definition: 95/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | The core problem (3 overlapping concepts reduced to 2) is clearly stated with before/after YAML, and the 1D-vs-2D framing ("This is a 1D narrowing model, not a 2D matrix") is precise. The addition of the new contributor onboarding anecdote ("session 2026-05-10, the contributor asked 'what profile should I use for a Go project?'") is a concrete user-facing data point that grounds the problem. Minor gap: the problem title is still in Chinese ("概念边界模糊") but the body is English-first, creating a minor bilingual inconsistency that could confuse non-Chinese readers. |
| Evidence provided | 35/40 | Major improvement. External validation added: "during onboarding of a new contributor (session 2026-05-10), the contributor asked 'what profile should I use for a Go project?'" — this is a direct user quote with a date. The Evidence table retains source-code-level specifics. Remaining gap: this is one anecdote from one session. No quantitative data (e.g., "3 of 4 onboarding sessions in 2026-Q2 produced the same question"). The evidence is directionally strong but thin in breadth. |
| Urgency justified | 22/30 | Unchanged from iteration 2. The v3.0.0 breaking-change window is correctly identified. However, the urgency still lacks quantification: "Every new feature added to v2 ... compounds the profile x capability matrix" — how many features were added in the last quarter? How many new call sites per feature? Without this, the "cost of delay" argument remains directional rather than compelling. |

### 2. Solution Clarity: 115/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 40/40 | Highly concrete. Function signatures (`DetectProfiles()` -> `DetectLanguages()`), struct fields (remove `TestProfiles`/`Capabilities`, add `Interfaces`/`Languages`), embed paths (`all:profiles` -> `all:languages`), and CLI command group changes (`forge profile` -> `forge testing`) are all specified. The Go code snippet for `languageCapabilities` map is a concrete implementation detail. A reader can explain exactly what will be built. |
| User-facing behavior described | 45/45 | Unchanged from iteration 2 — excellent. D1, D4, D5, D6 fully specify observable behavior. Scenarios 1-6 cover happy path and edge cases. |
| Technical direction clear | 30/35 | Improved. The Feasibility Assessment section provides 5 numbered subsections with specific code changes. The `languageCapabilities` Go snippet and the inline trade-off comment about hardcoding vs manifest.yaml add implementation clarity. Remaining gap: no data flow description (CLI command -> config reading -> detection -> strategy retrieval -> output). The 5 subsections describe what changes in each file, but not how they interact during a single request. |

### 3. Industry Benchmarking: 90/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 38/40 | 4 references with URLs, "How it works" and "What Forge borrows" columns. The `developer.hashicorp.com` URL is now correct (fixed from the typo noted in iteration 2). All references are real, verifiable products. Minor gap: no version numbers or dates cited — these URLs could break or become outdated. |
| At least 3 meaningful alternatives | 20/30 | The "Rename only" alternative is substantially improved — its cons now include concrete technical analysis: "`DetectProfiles()` still returns profile names (renamed to 'language' but still explicit), and the user must still set `languages: [go-test]` in config.yaml because auto-detect remains a fallback behind `ReadTestProfiles()`". This is genuine analysis, not a straw-man dismissal. The "Gradual deprecation" alternative now explicitly acknowledges the Python 2->3 cautionary tale: "Python 2→3 itself is widely regarded as a cautionary tale: the decade-long dual-support period fragmented the ecosystem". This directly addresses the iteration-2 attack. Remaining gap: the "Rename only" alternative's pros column is improved ("Lower risk — no directory restructuring, no embed path changes, no `manifest.yaml` removal; all existing tests pass unchanged; no migration guide needed for internal paths") but its analysis still leans toward explaining why it is insufficient rather than genuinely evaluating what it achieves. The pro is real analysis but the con is longer and more detailed than the pro, creating an imbalance. |
| Honest trade-off comparison | 18/25 | The selected approach now has an honest con listing: "All 10 skills (D6) must migrate CLI calls; `embed.go` directory structure changes; `manifest.yaml` files removed, metadata hardcoded in Go". The rejected alternatives have improved cons. Remaining gap: the pros for rejected alternatives are still thinner than their cons. "Do nothing" pros: "Zero migration cost; all skills already work" (2 clauses). "Rename only" pros: "Lower risk..." (3 clauses). Selected approach pros: 3 clauses. The selected approach cons: 3 clauses. The comparison is roughly balanced but the rejected alternatives' analysis still feels like the author is making a case against them rather than genuinely weighing them. |
| Chosen approach justified against benchmarks | 14/25 | The "Why Auto-detect Fits Forge Specifically" section has three numbered reasons (Reliable, Already implemented, Low-impact failure mode). The "Reliable" argument now explains the false-positive mitigation: "a `package.json` without Playwright in devDependencies does not trigger JavaScript detection, eliminating the most common false positive". This is concrete. However, the section does not engage the counterargument: why NOT follow the Terraform/ESLint explicit-declaration pattern? The industry references table mentions Terraform but the justification section does not address when explicit declaration is preferable (e.g., for non-standard project structures, or for users who want deterministic behavior without relying on heuristics). The counterargument is acknowledged only by omission. |

### 4. Requirements Completeness: 90/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 38/40 | Scenarios 1-6 cover happy path, multi-language, no detection, conflicting signals, and monorepo. Scenario 5 now correctly clarifies that JavaScript detection requires `@playwright/test` dependency — the confusing "falls back to Playwright" wording from iteration 2 is fixed. Remaining gap: no scenario for invalid `languages` values (e.g., `languages: [foobar]`) or invalid `interfaces` values (e.g., `interfaces: [telepathy]`). These are error scenarios that a requirements section should cover. |
| Non-functional requirements | 33/40 | NFR table covers Performance, Compatibility, Security, and Extensibility with verification methods. The iteration-2 gap about "what about larger repos?" is partially addressed — the table says "current `detect.go` reads at most 6 files from the project root", which explains why the 1000-file threshold is generous. Remaining gaps: (1) No NFR for the migration experience itself — how long should it take a user to migrate from v2 to v3? The migration guide is listed as in-scope but no NFR governs its quality. (2) No NFR for error message quality (e.g., "error messages must include the config field name and suggested fix"). |
| Constraints & dependencies | 19/30 | Technology constraints are listed (Go 1.22+, gopkg.in/yaml.v3, cobra). The "6 existing strategy sets migrated as-is" constraint is stated. Remaining gaps: (1) No mention of CLI output format stability as a constraint — skills parse CLI output, so format changes have downstream impact. (2) No mention of downstream consumers beyond the listed 10 skills — are there CI systems, scripts, or documentation that call `forge profile`? (3) The "6 existing strategy sets migrated as-is" constraint does not quantify the risk — how many files, how many lines, what is the blast radius if one is missed? |

### 5. Solution Creativity: 65/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 28/40 | The multi-framework language handling (language key as framework key: `javascript` = Playwright, future `javascript-jest`) is a genuine contribution that goes beyond simply copying CI auto-detect. The document now explicitly differentiates this from the baseline. Remaining gap: the core mechanism (file-signal detection -> strategy selection) is still directly lifted from CI tools. The differentiation is in the multi-framework edge case, not in the core design. |
| Cross-domain inspiration | 22/35 | Major improvement. The new "Cross-domain inspiration — type inference" paragraph draws an explicit analogy to Rust's local type inference: "Forge infers the 'type' (language) from project file signals and provides a narrow escape hatch (`languages` override) when inference is wrong." This is a genuine cross-domain reference (type systems, not developer tools). The analogy is well-constructed: inference handles 90% of cases, escape hatch for the remaining 10%. Remaining gap: this is one cross-domain reference. The section could benefit from a second analogy (e.g., how DNS uses caching with explicit TTL overrides, or how shell PATH resolution works with explicit overrides). |
| Simplicity of insight | 15/25 | The type-inference analogy adds conceptual clarity. The framework-key approach for JavaScript is practical. However, the insight remains "competent engineering" rather than surprising — it is the natural design for this problem space. The 1:1 language-to-framework mapping is clean but not the kind of insight that makes someone say "why didn't I think of that." |

### 6. Feasibility: 85/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 37/40 | The 5 subsection audit with function names, struct fields, and concrete changes is thorough. The `languageCapabilities` Go snippet with the inline trade-off comment ("Trade-off acknowledged: this hardcodes capability metadata in Go source, meaning adding a new language requires a code change and recompilation") directly addresses the iteration-2 attack about the contradiction with the Extensibility NFR. Remaining gap: no prototype or proof-of-concept. The feasibility is assessed by code reading, not by building. |
| Resource & timeline feasibility | 24/30 | The 7-item work breakdown with effort estimates totaling ~10 person-days is credible for a single developer. Remaining gaps: (1) No buffer/contingency — all estimates are round numbers (2d, 1d, 1d, 2d, 1d, 2d, 1d) suggesting top-down allocation rather than bottom-up estimation. (2) No explicit testing effort — the "Update config.yaml examples and documentation" item is 1 day, but who writes the E2E tests, benchmark tests, and integration tests referenced in the success criteria? These are not listed as work items. |
| Dependency readiness | 24/30 | The Dependency Readiness table now has one "Low" risk (strategy file content compatibility requires grep audit) instead of all "None", which is more honest. The explicit note about auditing strategy files for hardcoded `profiles/` references is a concrete action item. Remaining gaps: (1) The "Low" risk for strategy files is the only non-"None" risk — are there really no risks in the other 4 dependencies? For example, the `cobra` CLI framework row says "None" but adding a new command group involves routing, help text, and completion — is there no risk of CLI framework version constraints? (2) No assessment of whether the existing test suite will pass after the rename — this is a dependency (existing tests) whose readiness is not checked. |

### 7. Scope Definition: 75/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 27/30 | Eight concrete deliverables, each identifiable. The "All 10 skills consuming profile/capability updated" item now implicitly references D6. Remaining gap: the in-scope section does not include the success criteria as a reference — a reader must cross-reference two sections to understand what "done" means for each deliverable. |
| Out-of-scope explicitly listed | 23/25 | Five items explicitly deferred. Multi-framework JavaScript support is now explicitly out-of-scope. Remaining gap: user-facing documentation updates (beyond the migration guide) are not explicitly scoped. Is updating the main README or getting-started guide in scope or out of scope? |
| Scope is bounded | 25/25 | The timeline anchor is present: "2 weeks (10 person-days), single developer, targeting v3.0.0 release branch." This directly addresses the iteration-1 and iteration-2 attacks. The scope is bounded by both deliverables and time. |

### 8. Risk Assessment: 80/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | Expanded to 7 risks. The new "Hardcoded `languageCapabilities` map drifts from actual strategy capabilities" risk directly addresses the iteration-2 attack. All 7 risks are meaningful and non-trivial. Remaining gap: no risk for "v3.0 config.yaml is silently accepted by v2 CLI" — if a user runs the v2 CLI against a v3 config file (e.g., in a CI environment where the CLI version hasn't been updated), what happens? This is a version-skew risk not captured. |
| Likelihood + impact rated | 24/30 | All 7 risks have L/M/H ratings. The "Skill migration gap" risk now includes justification: "Rated M (not H) because: (1) the 10 skills are in a single monorepo with no external consumers — every migration point is grep-able; (2) CI smoke tests will run all 10 skills...; (3) `grep -r 'forge profile' plugins/` is a definitive completeness check." This directly addresses the iteration-2 attack about unjustified ratings. Remaining gap: no risk rating scale defined — what does "M" mean in terms of probability or frequency? The justification for specific ratings is good but the scale itself is not documented. |
| Mitigations are actionable | 28/30 | Each mitigation now includes specific actions. The hardcoded map drift mitigation is excellent: "CI test validates that `languageCapabilities[lang]` matches the actual interface types present in `languages/<lang>/` strategy files. This test runs on every PR that touches the `languages/` directory or `embed.go`." This is directly actionable — it describes a specific test with trigger conditions. The skill migration gap mitigation includes specific grep commands. Remaining gap: the file scanning performance mitigation says "No mitigation needed beyond the current design" and explains why (6 `os.ReadFile` calls only) — this is acceptable but the risk could simply be removed if no mitigation is needed, rather than keeping a risk with "no mitigation needed." |

### 9. Success Criteria: 75/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 50/55 | Expanded to 14 checklist items from 10 in iteration 2. New criteria address the iteration-2 attacks directly: (1) "functionally equivalent to the corresponding v2 output — verified by normalizing both outputs (strip trailing whitespace, normalize line endings, ignore version strings in headers)" replaces the unrealistically strict byte-for-byte comparison. (2) "Benchmark: `DetectLanguages()` completes in < 100ms" addresses the missing NFR criterion. (3) "E2E: each `forge testing` subcommand passes against a test fixture" addresses the missing E2E criterion. (4) "Error behavior: when no language is detected... exits with non-zero code and stderr contains the string 'languages'" addresses the missing error-behavior criterion. (5) "Extensibility: a developer can add a hypothetical 7th language by completing 3 steps" addresses the missing extensibility criterion. Remaining gap: the extensibility criterion says "verified by code review checklist" — a code review checklist is subjective. A more testable criterion would be: "A test exists that adds a mock 7th language and verifies detection + strategy retrieval works, confirming the 3-step process is sufficient." |
| Coverage is complete | 25/25 | All in-scope items now have corresponding success criteria: config schema, detection, CLI commands, skill migration, migration guide, NFRs (performance, E2E), error behavior, and extensibility. The gap from iteration 2 is fully addressed. |

### 10. Logical Consistency: 85/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 33/35 | The 1D-vs-2D argument is now well-supported. The new contributor anecdote in the Evidence section directly connects the problem ("what profile should I use?") to the solution (auto-detect eliminates the question). The override chain is clearly explained as a simplification (1D narrowing vs 2D matrix). Remaining gap: the problem says "Capability naming is confusing" but the trace from "capability -> interface" rename is implicit — the solution section does not explicitly say "we rename 'capabilities' to 'interfaces' because the values represent system-facing interfaces, not testing capabilities." This connection is made in D1 and the Evidence table but not in the Solution section itself. |
| Scope / Solution / Success Criteria aligned | 28/30 | The alignment is strong. Per-feature narrowing is out-of-scope and not in success criteria. Mobile is in-scope for migration with a success criterion. The Extensibility NFR now has a corresponding success criterion. Remaining gap: Scenario 5 ("Conflicting signals") describes specific behavior but no success criterion directly verifies conflicting-signal handling. The multi-language criterion (criterion 4) tests 2-language detection but not the specific "go.mod + package.json without Playwright" case. |
| Requirements / Solution coherent | 24/25 | Scenarios map to the solution. NFRs map to verifiable criteria. The Extensibility NFR is now coherent with the hardcoded-map trade-off acknowledgment. Remaining gap: the Compatibility NFR says "output is plain text with structured blocks, matching current `forge profile` format" but the solution changes the CLI from `forge profile get <name>` to `forge testing get generate` — the output format may match, but the command to invoke it changes, which is the breaking change. The NFR does not distinguish between "output format compatibility" and "CLI interface compatibility." |

## ATTACKS

1. **Problem Definition**: Urgency quantification remains absent — "Every new feature added to v2 ... compounds the profile x capability matrix" is asserted without data. How many features were added in the last quarter? How many new call sites? Without this, the urgency is directional but not compelling. Must add at least one quantitative data point (e.g., "3 new call sites were added to testgen.go in the last 2 months, each requiring profile x capability matrix handling").

2. **Problem Definition**: Evidence is one anecdote from one onboarding session — "during onboarding of a new contributor (session 2026-05-10)". A single data point is better than none but still thin. Must add breadth: "3 of 4 onboarding sessions in 2026-Q2 produced the same question" or cite a second independent data point.

3. **Industry Benchmarking**: The "Why Auto-detect Fits Forge Specifically" section does not engage the counterargument for explicit declaration. The Industry References table mentions Terraform and ESLint as examples of explicit-declaration tools, but the justification section does not address when explicit declaration is preferable (non-standard project structures, deterministic CI, user control preference). Must add a paragraph explaining why the auto-detect trade-off is acceptable even for users who would prefer explicit control.

4. **Requirements Completeness**: No scenario for invalid config values — what happens when `languages: [foobar]` or `interfaces: [telepathy]` is set? The Extensibility NFR says `interfaces` uses a "closed enumeration" (web-ui, tui, mobile-ui, api, cli) but no scenario or error handling describes what happens when an invalid value is provided. Must add: "Scenario 7: Invalid `languages` or `interfaces` value — `forge testing` returns an error listing valid values."

5. **Feasibility**: No work item for testing — the Resource & Timeline table lists 7 work items totaling 10 person-days but none covers writing the E2E tests, benchmark tests, or integration tests that the success criteria require (criteria 10-12). The 1-day "Update config.yaml examples and documentation" item does not cover test writing. Must add a work item for test creation or explicitly account for testing within existing items.

6. **Feasibility**: The Dependency Readiness table shows 4 of 5 dependencies at "None" risk. The `cobra` CLI framework row says "None" but creating a new command group (`forge testing`) involves routing, help text generation, shell completion, and flag registration — all within the framework's API. Must assess whether the current cobra version supports the proposed command structure without version upgrades or API changes.

7. **Solution Clarity**: No data flow description for a single CLI request. The 5 subsections describe what changes in each file but not how they interact. For example: user runs `forge testing get generate` -> what reads config? -> what calls detection? -> how is the result passed to strategy retrieval? -> how is the strategy returned to the CLI? A brief data flow paragraph would close this gap.

8. **Risk Assessment**: No version-skew risk — if a user runs the v2 CLI against a v3 config.yaml (or vice versa), what happens? This is a real risk in CI environments where the CLI version may lag behind the config file format. Must add a risk row for "v2 CLI encounters v3 config format (or vice versa) and produces confusing errors or silent misbehavior."

9. **Risk Assessment**: Risk rating scale is not defined — L/M/H are used without a legend. "M" for likelihood could mean "happens in 10% of projects" or "happens in 50% of projects." The individual justifications (e.g., for skill migration gap) are good, but the scale itself is ambiguous. Must add a brief scale definition (e.g., "L = <10% of projects, M = 10-40%, H = >40%").

10. **Logical Consistency**: The Compatibility NFR says "output is plain text with structured blocks, matching current `forge profile` format" but the CLI command itself changes from `forge profile get <name>` to `forge testing get generate`. The NFR addresses output format but not CLI interface compatibility. These are different things — a skill that parses CLI output is affected by format changes, but a CI script that calls the command is affected by the command name change. Must split this NFR into "output format stability" and "CLI interface stability" or clarify that CLI interface breaking is intentional and covered by the v3.0 major version bump.

11. **Logical Consistency**: Scenario 5 ("Conflicting signals") describes specific behavior for `go.mod` + `package.json` without Playwright, but no success criterion directly verifies this. Criterion 4 tests multi-language detection (2 languages returned) but not the specific "no Playwright dependency -> no JavaScript detection" case. Must add: "When a project has `go.mod` and `package.json` with no `@playwright/test` dependency, `forge testing detect` returns only `go`."
