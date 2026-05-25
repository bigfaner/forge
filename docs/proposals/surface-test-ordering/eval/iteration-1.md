# Eval Iteration 1: Surface Test Ordering & Journey Unification

**Evaluator**: CTO Adversary
**Date**: 2026-05-25
**Proposal**: docs/proposals/surface-test-ordering/proposal.md
**Rubric**: plugins/forge/skills/eval/rubrics/proposal.md
**Previous Score**: 680/1000 (iteration-0)

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem**: Two issues — (A) no cross-surface test ordering for fail-fast feedback, (B) gen-journeys per-surface split contradicts Journey's cross-surface semantic definition.

**Solution**: (A) run-tests split into per-surface-key serial tasks with `execution-order` config, (B) gen-journeys merged to single task.

**Evidence**: Code references (autogen.go, SKILL.md, config.go), semantic argument about Journey definition, token budget reasoning.

**Success Criteria**: 6 testable checklist items.

### Self-Contradiction Check

1. The proposal claims "minimal change: only modify autogen.go dependency resolution and config.go, not gen-journeys/gen-contracts/gen-scripts SKILL.md core logic" (Non-Functional Requirements). But the Key Risks table lists "gen-journeys SKILL.md must adapt multi-surface internal traversal" with M/M likelihood/impact. The in-scope section also mentions `renderBody` adaptation for empty `TestType`. These are SKILL.md / template-level changes that directly contradict the "minimal change" NFR assertion.

2. The Evidence section claims "gen-journeys is pure narrative extraction (reads PRD + writes MD), doesn't read code, token budget small (20-30min), parallel benefit ~0." But the proposal simultaneously requires the merged gen-journeys to "internally traverse all surface types loading corresponding rules." If surface rules influence generation output, the task is not "pure narrative" — it depends on surface-specific rules. The "pure narrative" claim is used to justify merging but is undermined by the implementation requirement.

3. The Evidence section references "SKILL.md HARD-RATE" but the actual SKILL.md uses `HARD-RULE` (confirmed via codebase verification). This is a minor factual error but undermines confidence in the evidence quality — did the author verify the source or rely on memory?

4. The Single Surface scenario says "degrades to suffix-less `T-test-run` (not `T-test-run-api`), gen-journeys similarly degrades to suffix-less single task, behavior completely identical to before." But the In Scope section says "InferType function: change T-test-run exact match to typeSuffixedID prefix match." If InferType now uses prefix matching, it would match `T-test-run` differently than before (previously it was an exact match case). This creates a subtle behavioral change even in the single-surface case — the matching semantics change, even if the result is the same for now. Future task IDs starting with `T-test-run-` could create ambiguity.

---

## Phase 2: Dimension Scoring

### D1. Problem Definition (110 pts)

**Problem stated clearly**: 38/40
Both problems are concrete and unambiguous. The two-problem bundling concern from iteration-0 remains (no explicit coupling justification), but the problems share the same codebase locus (autogen.go surface handling) and the same config schema (surfaces map), which provides implicit coupling. Pre-revision did not address this gap but the problems themselves are clearly stated.
> Deduction: -2 for bundled problems without explicit coupling argument.

**Evidence provided**: 32/40
Evidence references real code paths. The SKILL.md "HARD-RATE" typo (should be "HARD-RULE") is a minor accuracy issue. More substantively:
- "gen-journeys is pure narrative extraction... parallel benefit ~0" remains an assertion without profiling data.
- The "two independent gen-journeys tasks extracting from same PRD may produce inconsistent boundaries" evidence is strong and concrete.
- Pre-revision added evidence about `GetBreakdownTestTasks` and `GetQuickTestTasks` both having per-type loops, which is verifiable and correct.
> Deduction: -4 for unsubstantiated performance claims; -4 for "HARD-RATE" misquote suggesting insufficient source verification.

**Urgency justified**: 22/30
The urgency argument ("must land before justfile proposal") remains unchanged from iteration-0. Pre-revision did not address the unexamined alternative of integrating into the justfile proposal. The "cost of delay" remains vague — no quantification of "duplicate implementation" in effort terms.
> Deduction: -5 for unexamined merge-into-justfile alternative; -3 for vague cost-of-delay.

**D1 Total: 92/110**

---

### D2. Solution Clarity (120 pts)

**Approach is concrete**: 36/40
Pre-revision added the 3-surface dependency chain example (lines 39-63), which significantly improves clarity. The breakdown and quick mode diagrams show the exact task topology. Task IDs are now explicit (`T-test-run-{surface-key}`). However:
- The 3-surface example uses surface-keys as task suffixes but does not show how `execution-order` resolves when a surface-key contains characters that are invalid in task IDs (e.g., hyphens, which are already used as delimiters in `T-test-run-auth-service`). The proposal later adds a surface-key regex constraint (`[a-z][a-z0-9-]*`), but this constraint makes the example `auth-service` valid yet does not address how to distinguish `T-test-run-auth-service` as "surface-key auth-service" vs "some future task type auth" — the parsing is ambiguous.
> Deduction: -4 for task ID parsing ambiguity in the surface-key naming scheme.

**User-facing behavior described**: 38/45
Key scenarios expanded by pre-revision to include single-surface scalar form degradation. The upstream failure propagation scenario is described. However:
- Error message content is still not shown. The proposal says "reports error, prompts user to configure execution-order" but no actual error message text.
- No description of what happens when a user adds a new surface to an existing config — do they need to update execution-order? If yes, what's the migration UX?
- The blocked status behavior is described ("T-test-run-frontend status becomes blocked, skips execution") but not how users diagnose WHY a task was blocked.
> Deduction: -7 for incomplete error UX and missing add-surface migration scenario.

**Technical direction clear**: 28/35
Pre-revision significantly improved this dimension by adding:
- `InferType` change to prefix matching (line 144)
- `T-test-verify-regression` dependency on chain tail (line 145)
- Surface-key regex constraint (line 90)
- `renderBody` empty TestType adaptation (line 137)

However:
- The proposal does not mention that `GetBreakdownTestTasks` and `GetQuickTestTasks` function signatures must change from `capabilities []string` (deduplicated types) to include the full surfaces map or surface-keys. Currently, `capabilities` is `[]string` of deduplicated types (e.g., `["api", "web"]`), but the proposal needs per-surface-key iteration, which requires the full `map[string]string` with keys. This is a fundamental API change not mentioned.
- The `SourceTaskID` migration (line 148) is described but the technical mechanism ("detect index.json entries with SourceTaskID: T-test-run, auto-remap") is vague. How does autogen.go read index.json during task generation? This is a BuildIndex concern, not a GetBreakdownTestTasks concern.
> Deduction: -7 for missing function signature change and vague migration mechanism.

**D2 Total: 102/120**

---

### D3. Industry Benchmarking (120 pts)

**Industry solutions referenced**: 25/40
Unchanged from iteration-0. Only GitHub Actions `needs` and GitLab CI `dependencies` are mentioned in a single sentence. No analysis of HOW these systems handle ordering, default priority, or conflict resolution. No open-source projects, no published patterns, no specifications.
> Quote: "CI systems usually express cross-service test ordering through job dependency graphs (e.g., GitHub Actions' `needs` field, GitLab CI's `dependencies`). This validates that task-level dependencies are an industry standard approach."
> This is name-dropping, not benchmarking.
> Deduction: -15 for superficial references.

**At least 3 meaningful alternatives**: 25/30
Four alternatives listed. "Do nothing" and "execution-level ordering" are genuinely different. "Post-gen dependency injection" is a straw-man — dismissed as "semantically unclear" without explanation of WHY the semantics are unclear (gen-scripts ordering IS semantically meaningful: you want to generate scripts for the surface you'll test first). Only one industry-validated alternative.
> Deduction: -5 for one weak dismissal.

**Honest trade-off comparison**: 18/25
Unchanged from iteration-0. Selected approach cons: "needs to modify autogen.go" — still vastly understated given the actual scope (InferType, function signatures, migration, naming scheme). However, pre-revision added surface-key constraint and migration items to scope, partially addressing the scope understatement. The comparison table itself was NOT updated to reflect this fuller scope.
> Deduction: -7 for asymmetric trade-off analysis.

**Chosen approach justified against benchmarks**: 15/25
Justification remains "user confirmed" (appeal to authority, not technical analysis). No comparison of task-level dependencies vs execution-level ordering in Forge's specific architectural context. The freeform review made a strong case for execution-level ordering avoiding complexity; the proposal does not engage with this counter-argument.
> Deduction: -10 for justification by user approval rather than technical analysis.

**D3 Total: 83/120**

---

### D4. Requirements Completeness (110 pts)

**Scenario coverage**: 33/40
Pre-revision added the single-surface scalar form scenario (line 76), addressing iteration-0's gap. Four scenarios now cover: typical fullstack, multi-api conflict, single-surface degradation, and failure propagation. Missing:
- What happens when a surface is REMOVED from config? Old per-surface-key tasks become orphans.
- What happens when execution-order is partially specified (2 of 3 surfaces listed)? The proposal says "config load time validation" but does not specify whether partial specification is an error or ignored.
- Quick mode behavior is not listed as a separate scenario — the same scenarios should be verified for quick mode since the changes affect both `resolveQuickDeps` and `resolveBreakdownDeps`.
> Deduction: -7 for missing removal, partial config, and quick-mode scenarios.

**Non-functional requirements**: 35/40
Pre-revision significantly improved this by adding:
- Surface-key regex constraint (config load time validation)
- Validation timing (fail fast at config load)

The "backward compatibility" NFR is now partially addressed by the single-surface degradation scenario. However:
- The "minimal change" NFR ("only modify autogen.go and config.go") remains inaccurate — InferType changes (infer.go), function signature changes, and renderBody changes span beyond these two files.
- No performance NFR for the serial execution trade-off (happy-path latency increase).
> Deduction: -5 for inaccurate minimal-change claim.

**Constraints & dependencies**: 25/30
Dependencies on surface-aware-justfile are stated. Pre-revision added surface-key constraint details and validation timing. Missing:
- The `InferType` exact-match contract as a constraint.
- The function signature change from `capabilities []string` to per-key iteration as a constraint.
- Whether the surface-aware-justfile proposal has conflicting plans for run-test restructuring (this dependency is unidirectional — only this proposal acknowledges the other, not vice versa).
> Deduction: -5 for missing technical constraints.

**D4 Total: 93/110**

---

### D5. Solution Creativity (100 pts)

**Novelty over industry baseline**: 15/40
Unchanged. The selected approach IS the industry baseline (task-level dependencies). The "convention + override" pattern is standard convention-over-configuration (Rails, 2005). Pre-revision did not address this gap. The Innovation Highlights section still claims:
> "Hybrid ordering strategy (convention + override) balances zero-config experience and flexibility"
This is not innovation. It is standard software engineering practice.
> Deduction: -25 for claiming novelty where none exists.

**Cross-domain inspiration**: 5/35
Unchanged. No references to build systems (Make/Bazel), package managers (npm peerDeps), workflow engines (Airflow DAGs), or any domain beyond CI/CD.
> Deduction: -30 for zero cross-domain inspiration.

**Simplicity of insight**: 18/25
The gen-journeys merge remains a genuinely elegant insight — journeys are cross-surface by definition, splitting by surface was always wrong. The single-surface degradation rule (no suffix when scalar form) added by pre-revision is a nice touch that avoids unnecessary naming verbosity. However, the run-test split remains forced — the complexity it introduces (naming scheme, migration, InferType changes) is disproportionate to the benefit.
> Deduction: -7 for forced run-test split complexity.

**D5 Total: 38/100**

---

### D6. Feasibility (100 pts)

**Technical feasibility**: 30/40
Pre-revision added `InferType` prefix matching (line 144) and surface-key regex constraint (line 90), partially addressing iteration-0's gap. However:
- Function signature change from `capabilities []string` to per-surface-key iteration is still not mentioned. Currently `GetBreakdownTestTasks(capabilities []string, ...)` receives deduplicated types. The proposal needs surface-keys (not types), which requires either changing the function signature or adding a new parameter.
- The `renderBody` adaptation for empty `TestType` is mentioned but the mechanism is not explained. Currently, when `TestType` is empty, the template omits the `{{TEST_TYPE}}` line. But the merged gen-journeys task has no `TestType` — does this mean it loses the surface type information in its task body? How does the downstream agent know which surface rules to load?
- The SourceTaskID migration mechanism (line 148) is described as "autogen.go task generation logic detects index.json entries." But `GetBreakdownTestTasks` and `GetQuickTestTasks` do not read index.json — that happens in `BuildIndex`. The migration belongs in BuildIndex, not in the task definition functions.
> Deduction: -10 for missing function signature change and misplaced migration logic.

**Resource & timeline feasibility**: 22/30
Still estimated at "3-5 coding tasks." Pre-revision added several new in-scope items (InferType update, renderBody adaptation, SourceTaskID migration, verify-regression dependency, surface-key regex). These additions make 3-5 tasks optimistic. A more realistic estimate is 5-7 tasks.
> Deduction: -8 for optimistic task count.

**Dependency readiness**: 25/30
Surface-aware-justfile proposal acknowledged. The proposal claims it can "introduce SurfaceKey fields first" but does not verify whether the justfile proposal's task struct design conflicts. The justfile proposal defines surface-key differently (with `/` to `+` conversion) — this proposal uses the raw YAML map key. This is a potential conflict not examined.
> Deduction: -5 for potential surface-key definition conflict with justfile proposal.

**D6 Total: 77/100**

---

### D7. Scope Definition (80 pts)

**In-scope items are concrete**: 28/30
Pre-revision significantly expanded in-scope items. Most are now concrete deliverables. The InferType update, renderBody adaptation, and verify-regression dependency are specific. Remaining vagueness:
- "Update resolveBreakdownDeps and resolveQuickDeps dependency chains" — the WHAT is clearer now (serial chain), but the HOW is still unspecified. A loop becomes a serial chain — but how is the serial order constructed? From execution-order config? From default priority? Both paths need specification.
> Deduction: -2 for one vague in-scope item.

**Out-of-scope explicitly listed**: 22/25
Six out-of-scope items. However:
- "Partial continue running" is out-of-scope but the serial design makes this architecturally impossible without changes. This should be noted as a design decision that forecloses a future option.
- gen-journeys SKILL.md adaptation IS in the risks table (M/M) but NOT listed as in-scope or out-of-scope. The scope section says "don't touch SKILL.md core logic" but the risk says "SKILL.md must adapt." This is contradictory.
> Deduction: -3 for gen-journeys SKILL.md scope ambiguity.

**Scope is bounded**: 22/25
The expanded in-scope items provide better bounds. However, the scope boundary is still based on incomplete understanding of the function signature change impact. If capabilities → surfaces map is needed, the scope ripples into BuildIndex and all callers of the test task functions.
> Deduction: -3 for incomplete impact boundary.

**D7 Total: 72/80**

---

### D8. Risk Assessment (90 pts)

**Risks identified**: 24/30
Pre-revision added the index.json orphan risk (line 168), partially addressing iteration-0's gap. Now 5 risks listed. Missing:
- InferType contract breakage (the prefix matching change could match unintended task IDs).
- Surface-key definition conflict with justfile proposal (that proposal uses `/` to `+` conversion; this one uses raw map keys).
- Function signature change from capabilities to surfaces map as a risk to all callers.
> Deduction: -6 for 3 missing risks.

**Likelihood + impact rated**: 22/30
The ratings use M/L, M/M, L/L, M/M. The index.json orphan risk is rated M/M which seems reasonable. However:
- "gen-journeys single task loading multi-surface rules increases context noise" (M/L) — but the proposal claims gen-journeys is "pure narrative extraction." If true, noise is irrelevant. This contradiction was flagged in iteration-0 and remains unaddressed.
- "Default priority doesn't cover all scenarios" is rated L/L but the proposal's own 3-surface example uses `cli` type which IS covered. The real gap is types not in the priority list (e.g., `desktop`, `embedded`, or any user-defined type), which is not discussed.
> Deduction: -8 for ratings contradicting other claims.

**Mitigations are actionable**: 22/30
Pre-revision improved the index.json orphan mitigation:
> "Migration copies old T-test-run status/blocked-reason to T-test-run-{first-surface-key}; single surface no issue; multi-surface by execution-order first surface inherits state."

This is more actionable than iteration-0. However:
- "gen-journeys only does information reference, no surface-specific generation, noise impact limited" — still a dismissal, not a mitigation. If noise is real enough to be a risk, the mitigation should describe HOW noise is limited.
- "SKILL.md already has surface detection logic, expand to multi-surface traversal" — this is a plan, not a mitigation. The risk is that the change introduces bugs; the mitigation should describe how to verify correctness.
> Deduction: -8 for non-actionable mitigations.

**D8 Total: 68/90**

---

### D9. Success Criteria (80 pts)

**Criteria are measurable and testable**: 26/30
Pre-revision added two new SCs: single-surface degradation (SC5) and execution-order validation (SC6). Most criteria are testable. Remaining issues:
- SC3: "gen-journeys generates a single T-test-gen-journeys task, internally loads all surface type rules" — "internally loads" is not observable. What's the externally visible output difference? The SC should specify that the output Journey files cover all surface types, not just one.
- SC5: "completely identical to before" — which fields must match? Task ID? Dependencies? Title? "Completely" is vague.
> Deduction: -4 for two criteria with observability/vagueness issues.

**Coverage is complete**: 18/25
Pre-revision added SC5 and SC6, improving coverage. Gaps remain:
- No SC for quick mode behavior (resolveQuickDeps changes are in-scope but no SC tests quick mode dependency chain).
- No SC for surface-key normalization/rejection (the regex constraint is in-scope but no SC tests that invalid keys are rejected).
- No SC for InferType correctly identifying new task IDs (InferType update is in-scope but no SC tests it).
- No SC for SourceTaskID migration correctness (migration is in-scope but no SC verifies it).
- No SC for default priority order (the api > web > cli > tui > mobile convention is in-scope but no SC tests a scenario where it's exercised without explicit execution-order).
> Deduction: -7 for 5 uncovered in-scope items.

**SC internal consistency**: 20/25
Cluster analysis:

**gen-journeys cluster**: SC3 (single gen-journeys task) is satisfiable. In-scope says "TestType field left empty, renderBody adapts." These are consistent. However, the Key Risks table says "gen-journeys SKILL.md must adapt multi-surface traversal" — if SKILL.md changes are required for SC3 to pass, this should be reflected in the success criterion (e.g., "gen-journeys output includes Journey files covering all configured surfaces").

**run-test cluster**: SC1 (ordering), SC2 (failure propagation), SC4 (conflict detection), SC5 (degradation), SC6 (validation). SC1 and SC2 are satisfiable together. SC5 says "completely identical" — but the InferType matching semantics change from exact to prefix. The BEHAVIOR may be identical but the IMPLEMENTATION is different. If "completely identical" means "same task ID and same dependency list," then SC5 is satisfiable. If it means "same code paths," it's not.

SC4 (conflict detection at `forge task build` time) and SC6 (config load time validation) use different validation timing. SC4 says "forge task build time" but the Constraints section says "execution-order reference validation at config load time." This is an inconsistency — same-type conflict detection at build time vs. reference validation at config load time. Why two different times?
> Deduction: -5 for SC4/SC6 validation timing inconsistency and gen-journeys SC lacking output quality verification.

**D9 Total: 64/80**

---

### D10. Logical Consistency (90 pts)

**Solution addresses the stated problem**: 28/35
- Problem A (no cross-surface ordering): per-surface-key serial tasks with execution-order solves this. The solution works.
- Problem B (gen-journeys semantic inconsistency): merging to single task directly addresses this. Clean alignment.
- However: execution-level ordering within a single run-test task would also solve Problem A with less complexity. The proposal acknowledges this alternative existed but rejected it ("user has vetoed"). The solution works but is not the simplest solution to the stated problem.
- The two problems are still bundled without explicit coupling justification. The gen-journeys merge could ship independently.
> Deduction: -7 for suboptimal solution choice and unexplained problem coupling.

**Scope <-> Solution <-> Success Criteria aligned**: 20/30
Pre-revision improved alignment by adding SC5, SC6, InferType update, and migration to scope. Remaining gaps:
- In-scope lists "default priority convention: api > web > cli > tui > mobile" but no SC tests this priority specifically (SC1 only tests api+web, which is the FIRST two in the list, not the full ordering).
- In-scope lists "update resolveQuickDeps" but no SC tests quick mode.
- In-scope lists "surface-key regex constraint" but no SC tests rejection of invalid keys.
- The solution's naming strategy (surface-key, not surface-type) is well-motivated but conflicts with the existing system's type-suffix convention (`T-test-gen-journeys-api` uses type suffix, not key suffix). The proposal addresses gen-journeys (merged, no suffix) and run-tests (key suffix) but does not address gen-scripts, which continues to use type suffixes. This creates an inconsistency: gen-scripts uses type suffixes, run-tests uses key suffixes, for the SAME surfaces.
> Deduction: -10 for missing SC coverage and naming inconsistency with gen-scripts.

**Requirements <-> Solution coherent**: 18/25
- The "minimal change" NFR is incoherent with the expanded scope (InferType changes, function signatures, migration, regex validation). Pre-revision expanded scope but did not update the NFR.
- The "backward compatibility" NFR is mostly addressed by SC5 (single-surface degradation) and the migration plan. However, the multi-surface case (new projects) is not "backward compatible" by definition — there IS no backward state.
- The "conflict-with-pre-revision" tag applies here: pre-revision added significant new scope items (InferType, migration, surface-key constraint) that make the "minimal change" NFR from the original proposal inaccurate. The NFR should have been updated or removed.
> Deduction: -7 for NFR/solution incoherence (conflict-with-pre-revision).

**D10 Total: 66/90**

---

## Phase 3: Blindspot Hunt

**[blindspot-1] No rollback plan.** Unchanged from iteration-0. If serial ordering causes unacceptable pipeline latency increase (API tests take 2h, web tests can't start until they finish), there is no mechanism to revert to parallel execution. The proposal optimizes for fail-fast but may degrade the happy path. No measurement plan.

**[blindspot-2] Surface-key definition conflict with justfile proposal.** The justfile proposal defines surface-key as "surfaces map key with `/` converted to `+`." This proposal uses raw YAML map keys as surface-key. These are two DIFFERENT definitions. If both proposals land, a config with `surfaces: {"admin/panel": web}` would produce surface-key `admin+panel` in justfile but `admin/panel` (or normalized `admin-panel`) in this proposal. The two proposals must agree on surface-key definition.

**[blindspot-3] gen-scripts still uses type suffixes.** While the proposal changes run-tests to use surface-key suffixes and merges gen-journeys (no suffix), gen-scripts continues to use type suffixes (`T-test-gen-scripts-api`). This means the SAME project with `surfaces: {frontend: web, backend: api}` has gen-scripts tasks named by TYPE (`-api`, `-web`) but run-tests tasks named by KEY (`-frontend`, `-backend`). This naming inconsistency is not discussed and could confuse users.

**[blindspot-4] Serial execution is a performance regression for independent surfaces.** If backend API and frontend web are truly independent (no shared state), running them serially is pure waste. The proposal assumes dependencies exist between all surfaces by default, but many fullstack projects have independent test suites. The fail-fast benefit only materializes when tests actually fail — which is the minority case.

**[blindspot-5] No integration test plan.** The proposal modifies the core autogen pipeline (both breakdown and quick mode) but specifies no integration test strategy. All SCs are unit-level assertions on single config scenarios.

---

## Bias Detection Report

Annotated regions analysis:
- 10 pre-revised markers covering approximately 12 paragraphs (some markers annotate adjacent paragraphs treated as one unit).
- Attacks on annotated regions: 8 attack points (naming ambiguity in D2, function signature gap in D2/D6, scope ambiguity in D7, SC4/SC6 timing in D9, migration mechanism in D6, InferType subtle change in Phase 1, surface-key conflict in blindspot-2, NFR incoherence in D10).
- Annotated density: 8/12 = 0.67

Unannotated regions analysis:
- Approximately 25 paragraphs in unannotated sections.
- Attacks on unannotated regions: 10 attack points (SKILL.md HARD-RATE typo in D1, performance assertion in D1, urgency alternatives in D1, benchmarking superficiality in D3, trade-off asymmetry in D3, user-approval justification in D3, quick-mode gap in D4/D9, performance NFR in D4, missing risks in D8, rollback in blindspot-1).
- Unannotated density: 10/25 = 0.40

**Ratio (annotated/unannotated): 1.67**

The higher attack density on annotated regions is expected — pre-revision added substantive new claims (surface-key naming, InferType prefix matching, migration mechanism, scalar degradation) that introduce new surface area for critique. The ratio is within acceptable bounds (under 2.0) and reflects genuine issue density rather than reviewer bias. One `conflict-with-pre-revision` tag was applied (D10 NFR coherence).

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1. Problem Definition | 92 | 110 |
| D2. Solution Clarity | 102 | 120 |
| D3. Industry Benchmarking | 83 | 120 |
| D4. Requirements Completeness | 93 | 110 |
| D5. Solution Creativity | 38 | 100 |
| D6. Feasibility | 77 | 100 |
| D7. Scope Definition | 72 | 80 |
| D8. Risk Assessment | 68 | 90 |
| D9. Success Criteria | 64 | 80 |
| D10. Logical Consistency | 66 | 90 |
| **Total** | **755** | **1000** |

---

## Top Attacks (Priority Order)

1. **D5 Creativity**: Zero novelty beyond industry baseline and zero cross-domain inspiration — "hybrid ordering strategy (convention + override) balances zero-config experience and flexibility" is convention-over-configuration from 2005, not innovation. 62/100 lost. Must either acknowledge straightforward adoption of industry patterns or find genuine differentiation.

2. **D3 Industry Benchmarking**: Superficial references (one sentence mentioning GitHub Actions/GitLab CI), asymmetric trade-off analysis (selected approach cons understated), and justification by user approval rather than technical analysis — "Selected: user confirmed" is appeal to authority. 37/120 lost. Must provide substantive benchmarking and honest trade-offs.

3. **D6 Feasibility**: Function signature change from `capabilities []string` to surfaces-map-based iteration not mentioned; migration logic misplaced in task generation instead of BuildIndex; surface-key definition conflict with justfile proposal unexamined. 23/100 lost. Must enumerate ALL affected code paths including function signatures.

4. **D9 Success Criteria**: 5 in-scope items lack corresponding SCs (quick mode, InferType, migration, normalization, default priority); SC4/SC6 have inconsistent validation timing (build time vs config load time). 16/80 lost. Every in-scope item needs a testable criterion.

5. **D10 Logical Consistency**: gen-scripts continues using type suffixes while run-tests switches to key suffixes — naming inconsistency not discussed; "minimal change" NFR contradicts expanded scope from pre-revision (conflict-with-pre-revision); execution-level ordering rejected without adequate technical justification. 24/90 lost. Must resolve naming consistency and update NFR to match actual scope.

6. **[blindspot-2] Surface-key definition conflict**: The justfile proposal defines surface-key with `/` to `+` conversion; this proposal uses raw YAML map keys. These definitions will collide when both proposals are implemented.
