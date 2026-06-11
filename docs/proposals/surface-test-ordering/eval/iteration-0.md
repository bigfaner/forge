# Eval Iteration 0: Surface Test Ordering & Journey Unification

**Evaluator**: CTO Adversary
**Date**: 2026-05-25
**Proposal**: docs/proposals/surface-test-ordering/proposal.md
**Rubric**: plugins/forge/skills/eval/rubrics/proposal.md

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem**: Two issues identified — (A) no cross-surface test ordering for fail-fast, (B) gen-journeys per-surface split contradicts Journey semantics.

**Solution**: (A) run-tests split into per-surface-key serial tasks, (B) gen-journeys merged to single task.

**Evidence**: Code references to autogen.go, SKILL.md HARD-RATE, token budget claims.

**Success Criteria**: 6 testable checklist items.

### Self-Contradiction Check

1. The proposal claims "minimal change: only modify autogen.go dependency resolution and config.go" but then lists 8 in-scope items including dependency chain rewrites in both resolveBreakdownDeps AND resolveQuickDeps, conflict detection, default priority logic, and regression dependency fan-in. This understates scope.

2. The proposal says gen-journeys is "pure narrative extraction" that "doesn't read code" as evidence for merging, but simultaneously says the merged task must "internally traverse all surface types loading corresponding rules." If surface rules affect generation, the "pure narrative" claim is weakened.

3. The Innovation Highlights section calls "convention + override" a hybrid ordering strategy as if this is novel, but the Comparison Table itself cites GitHub Actions `needs` as the source. This is industry-standard, not innovative.

4. The proposal rejects "execution-level ordering (internal to run-test)" because it is "invisible to scheduler," but the freeform review (Recommendation 8) argues this avoids massive architectural cost for the same fail-fast behavior. The rejection rationale does not justify the complexity introduced.

---

## Phase 2: Dimension Scoring

### D1. Problem Definition (110 pts)

**Problem stated clearly**: 35/40
The two problems (test ordering + gen-journeys semantic inconsistency) are stated concretely. However, they are bundled into one proposal without justifying why they must be solved together. Are they coupled, or is this scope creep? The proposal never explains the dependency between these two changes.
> Deduction: -5 for bundling two problems without justifying coupling.

**Evidence provided**: 30/40
Four evidence items reference actual code (autogen.go, SKILL.md). However:
- "gen-journeys is pure narrative extraction" is an assertion, not evidence. No token usage data or profiling is cited.
- "parallel benefit ~0" claims subagent startup overhead "may offset" gains but provides no measurements.
- The SKILL.md HARD-RATE quote is strong evidence for the semantic inconsistency.
> Deduction: -10 for unsubstantiated performance claims presented as evidence.

**Urgency justified**: 20/30
> "surface-aware-justfile proposal already approved, involving task struct and run-tests refactoring. This proposal should land before justfile proposal implementation to avoid duplicate implementation and dependency chain conflicts."

This creates a false urgency: if justfile is already approved, why not integrate ordering INTO that proposal? The proposal assumes sequential dependency but doesn't explore co-design. Also, "cost of delay" is stated vaguely — no quantification of what "duplicate implementation" means in effort terms.
> Deduction: -10 for unexamined alternative of merging into the approved proposal.

**D1 Total: 85/110**

---

### D2. Solution Clarity (120 pts)

**Approach is concrete**: 30/40
The two changes are described at a high level. However, critical details are missing:
- What is the task ID for the merged gen-journeys? `T-test-gen-journeys` (no suffix)? Not specified.
- What is the dependency chain for a 3-surface project? Only 2-surface examples are given.
- How does the serial ordering work when execution-order has 2 items but there are 3 surfaces? The proposal says "same-type conflict requires explicit config" but doesn't show the error recovery path.
> Deduction: -10 for missing concrete details on task IDs and multi-surface topology.

**User-facing behavior described**: 35/45
Key scenarios cover the main paths. However:
- Error messages are not shown (what does the "conflict detected" message look like?).
- The "blocked" status behavior is described but not how users discover WHY a task was blocked.
- No description of what happens when a user adds a new surface to an existing config — do they need to update execution-order?
> Deduction: -10 for incomplete error/user experience description.

**Technical direction clear**: 25/35
Points to autogen.go and config.go as targets. Mentions resolveBreakdownDeps and resolveQuickDeps. But:
- Does NOT mention InferType changes (critical — exact match contract must change for T-test-run-{key}).
- Does NOT mention that function signatures must change from `capabilities []string` to include the full surfaces map.
- Does NOT address the surface-key vs surface-type naming inconsistency (existing system uses type suffixes, proposal introduces key suffixes).
> Deduction: -10 for missing three critical technical changes identified in freeform review.

**D2 Total: 90/120**

---

### D3. Industry Benchmarking (120 pts)

**Industry solutions referenced**: 25/40
Only GitHub Actions `needs` and GitLab CI `dependencies` are mentioned in passing. No specific implementation details from either system. No open-source projects or published patterns are cited. The reference is superficial — a single sentence.
> Quote: "CI systems usually express cross-service test ordering through job dependency graphs (e.g., GitHub Actions' `needs` field, GitLab CI's `dependencies`). This validates that task-level dependencies are an industry standard approach."
> This is a hand-waving reference, not a benchmark. No analysis of HOW these systems handle ordering, conflict resolution, or default priority.
> Deduction: -15 for superficial references.

**At least 3 meaningful alternatives**: 25/30
Four alternatives listed including "do nothing." However:
- "Execution-level ordering (internal to run-test)" is labeled "Rejected: user has vetoed" — this is a straw-man if the user was presented with a false dichotomy. The freeform review shows this alternative avoids massive complexity.
- "Post-gen dependency injection" is dismissed as "semantically unclear" without explaining WHY it's unclear.
- Only one alternative is industry-validated (GitHub Actions `needs`).
> Deduction: -5 for one straw-man alternative with weak dismissal.

**Honest trade-off comparison**: 15/25
The comparison table is skewed:
- Selected approach cons: only "needs to modify autogen.go" — vastly understated. The freeform review identifies 7+ additional impacts (InferType, naming schemes, migration, SourceTaskID, etc.).
- "Execution-level ordering" cons: "invisible to scheduler, cannot visualize" — but the proposal has no visualization in scope either! The out-of-scope section explicitly lists "visualization of cross-surface dependencies."
> Deduction: -10 for asymmetric trade-off analysis that understates selected approach's costs.

**Chosen approach justified against benchmarks**: 15/25
> Quote: "User confirmed" is the justification for the selected approach. This is appeal to authority, not technical justification. The proposal doesn't explain WHY task-level dependencies are better than execution-level ordering in the specific context of Forge's architecture. The freeform review's Recommendation 8 makes a compelling case that execution-level ordering avoids all the naming/migration complexity for the same user-visible behavior.
> Deduction: -10 for justification by user approval rather than technical analysis.

**D3 Total: 80/120**

---

### D4. Requirements Completeness (110 pts)

**Scenario coverage**: 30/40
Four scenarios covered (typical fullstack, multi-api, single-surface, failure propagation). Missing scenarios:
- What happens when a surface is REMOVED from config? Do old per-surface-key tasks get orphaned?
- What happens when execution-order is partially specified (2 of 3 surfaces)?
- What about the scalar form `surfaces: api` which internally becomes `{".": "api"}` — the dot-key sentinel?
- Quick mode: same changes needed but not explicitly listed as a scenario.
> Deduction: -10 for missing critical edge cases (dot-key sentinel, removal, partial config).

**Non-functional requirements**: 30/40
Two NFRs listed: backward compatibility and minimal change. Missing:
- **Migration strategy**: No plan for existing index.json with old task IDs. The freeform review identifies this as HIGH severity.
- **Performance**: No analysis of whether serial execution significantly increases total pipeline time vs parallel.
- **Security**: Surface-key values could contain path traversal characters; no validation mentioned.
> Deduction: -10 for missing migration and performance NFRs.

**Constraints & dependencies**: 25/30
Dependencies on surface-aware-justfile are stated. However:
- Does not mention the existing InferType contract as a constraint.
- Does not mention that autogen.go function signatures need to change (capabilities-only input).
- Does not surface the SourceTaskID migration constraint.
> Deduction: -5 for missing technical constraints.

**D4 Total: 85/110**

---

### D5. Solution Creativity (100 pts)

**Novelty over industry baseline**: 15/40
The selected approach IS the industry baseline (task-level dependencies, a la GitHub Actions `needs`). The "convention + override" pattern is standard configuration everywhere (convention-over-configuration is a well-established pattern from Ruby on Rails, 2005). The Innovation Highlights section claims novelty but provides none:
> Quote: "Hybrid ordering strategy (convention + override) balances zero-config experience and flexibility."
> This is not innovation. This is standard software engineering practice.
> Deduction: -25 for claiming novelty where none exists.

**Cross-domain inspiration**: 5/35
No evidence of cross-domain thinking. No references to how other domains (build systems like Make/Bazel, package managers, workflow engines) handle ordering. The proposal stays entirely within CI/CD conventions.
> Deduction: -30 for zero cross-domain inspiration.

**Simplicity of insight**: 15/25
The gen-journeys merge is a genuinely simple insight ("why were we splitting by surface when journeys are cross-surface by definition?"). The run-test split, however, introduces significant complexity (as documented in the freeform review) for questionable marginal benefit over execution-level ordering.
> Deduction: -10 for the run-test split being forced rather than elegant.

**D5 Total: 35/100**

---

### D6. Feasibility (100 pts)

**Technical feasibility**: 25/40
The proposal claims:
> "autogen.go already has per-surface loop logic, changing to single task only requires removing the loop"
> "config.go adding ExecutionOrder []string field is an incremental change"
> "dependency resolution function already has mature patterns"

This severely understates the technical changes required. The freeform review identifies:
- InferType exact-match contract must change (critical)
- Function signatures must change from capabilities to surfaces map
- Surface-key normalization for arbitrary YAML keys
- Dot-key sentinel handling for scalar surfaces
- SourceTaskID migration for existing fix-tasks

The proposal's technical feasibility assessment is incomplete to the point of being misleading.
> Deduction: -15 for incomplete technical feasibility analysis.

**Resource & timeline feasibility**: 20/30
> Quote: "Estimated 3-5 coding tasks, belongs to small feature scope."

Given the freeform review's findings (11 issues, 3 critical, 3 high), 3-5 tasks is likely an underestimate. The InferType change alone, with its test implications and the naming scheme decision, could be 1-2 tasks. Migration handling is another task. The freeform review's Recommendation 2 (align with surface-type not surface-key) would be a fundamental design decision requiring rework of the proposal itself.
> Deduction: -10 for underestimated scope.

**Dependency readiness**: 25/30
The surface-aware-justfile proposal dependency is acknowledged. The proposal correctly notes it can introduce SurfaceKey fields first. However, it doesn't verify whether the justfile proposal already has conflicting plans for run-test restructuring.
> Deduction: -5 for unverified overlap with justfile proposal.

**D6 Total: 70/100**

---

### D7. Scope Definition (80 pts)

**In-scope items are concrete**: 25/30
Eight in-scope items are listed. Most are concrete deliverables. However:
- "Update resolveBreakdownDeps and resolveQuickDeps dependency chains" — this is vague. What SPECIFIC changes? A loop becomes what?
- "Failure propagation: upstream run-test failure -> downstream blocked" — this is a behavior, not a deliverable. What code change implements this?
> Deduction: -5 for two vague in-scope items.

**Out-of-scope explicitly listed**: 20/25
Six out-of-scope items listed. However:
- "Partial continue running (one surface fails, continue running the rest)" is listed as out-of-scope, but the proposal's fail-fast design makes this impossible in the future without architectural changes. This should be noted as a design decision that forecloses a future option.
- InferType changes are NOT listed as in-scope or out-of-scope — they are simply not mentioned at all.
> Deduction: -5 for missing silent scope items.

**Scope is bounded**: 20/25
The "3-5 coding tasks" estimate provides a bound. However, the freeform review's findings suggest the actual scope is larger. The proposal's scope boundary is drawn based on incomplete understanding of affected code paths.
> Deduction: -5 for scope boundary based on incomplete impact analysis.

**D7 Total: 65/80**

---

### D8. Risk Assessment (90 pts)

**Risks identified**: 20/30
Four risks listed. Missing risks:
- InferType contract breakage (critical, identified in freeform review)
- Index.json migration / orphaned tasks (high, identified in freeform review)
- SourceTaskID migration for existing fix-tasks (high, identified in freeform review)
- Surface-key character validation (high, identified in freeform review)
- Scalar form dot-key sentinel producing invalid task IDs (critical, identified in freeform review)

Five significant risks are unaccounted for. The four risks that ARE listed are the less severe ones.
> Deduction: -10 for missing 5 significant risks.

**Likelihood + impact rated**: 20/30
The ratings use M/L and L/L. Notably:
- "gen-journeys single task loading multi-surface rules increases context noise" is rated M likelihood / L impact. But the proposal claims gen-journeys is "pure narrative extraction" — if that's true, why is noise even a risk? This contradiction suggests the M likelihood is understated or the L impact is wrong.
- "Task struct change conflict with surface-aware-justfile proposal" is rated M/M. But the proposal's urgency argument says this MUST land first — if there's a real M/M conflict risk, the urgency argument is undermined.
> Deduction: -10 for ratings that contradict other claims in the proposal.

**Mitigations are actionable**: 20/30
- "gen-journeys only does information reference, no surface-specific generation, noise impact limited" — this is not a mitigation, it's a dismissal.
- "This proposal introduces SurfaceKey field first, justfile proposal reuses" — this is a plan but not a mitigation for CONFLICT. What if justfile proposal needs SurfaceKey in a different form?
- "Users can explicitly override via execution-order" — this is a workaround, not a mitigation.
> Deduction: -10 for non-actionable mitigations.

**D8 Total: 60/90**

---

### D9. Success Criteria (80 pts)

**Criteria are measurable and testable**: 25/30
Most criteria are testable. However:
- "gen-journeys generates a single T-test-gen-journeys task, internally loads all surface type rules" — how do you test "internally loads all surface type rules"? This is an implementation detail, not an observable behavior. What's the observable output difference?
- "Single surface project task structure and dependency chain completely identical to before" — "completely identical" is vague. Which fields must match? What about task IDs that change?
> Deduction: -5 for two criteria with ambiguity.

**Coverage is complete**: 15/25
Gaps in SC coverage:
- No SC for quick mode behavior (resolveQuickDeps changes are in-scope but untested).
- No SC for execution-order validation error message content.
- No SC for surface-key normalization (what characters are valid?).
- No SC for backward compatibility with existing index.json (migration).
- No SC for InferType correctly identifying new task IDs.
- No SC for SourceTaskID migration.
> Deduction: -10 for 6 uncovered in-scope items.

**SC internal consistency**: 15/25
Cluster analysis:

**gen-journeys cluster**: SC3 says "generates single T-test-gen-journeys task." In-scope says "autogen.go change." This is consistent IF the task ID is T-test-gen-journeys. But what about eval-journey's dependency on gen-journeys? SC doesn't test this dependency correctly resolves.

**run-test cluster**: SC1 and SC2 test ordering and failure propagation for a 2-surface setup. SC4 tests conflict detection for 2-api setup. SC5 tests single-surface degradation. But SC5 says "completely identical to before" while SC1 introduces NEW task IDs (T-test-run-backend, T-test-run-frontend). For single-surface, would the task be T-test-run-api or T-test-run? If T-test-run, then single-surface is consistent. If T-test-run-api (because there's one surface), it's NOT identical. The proposal doesn't resolve this ambiguity.

SC6 (execution-order validation) and SC4 (same-type conflict) are consistent with each other.

SC1 and SC2 are satisfiable together. SC3 is independently satisfiable. SC5 potentially contradicts the naming scheme if dot-key sentinel is not handled.

> Deduction: -10 for ambiguous SC5 interpretation and missing eval-journey dependency SC.

**D9 Total: 55/80**

---

### D10. Logical Consistency (90 pts)

**Solution addresses the stated problem**: 25/35
- Problem A (no cross-surface ordering): Solution (per-surface-key serial tasks) addresses this. However, the freeform review shows that execution-level ordering within a single task would also solve this with far less complexity. The solution works but is not the simplest solution to the stated problem.
- Problem B (gen-journeys semantic inconsistency): Solution (merge to single task) directly addresses this. Clean alignment.
- But the proposal bundles these as if they're coupled. The gen-journeys merge could ship independently with zero relation to run-test ordering. The logical coupling is asserted, not demonstrated.
> Deduction: -10 for insufficient justification of problem coupling and for choosing a complex solution when a simpler one exists.

**Scope <-> Solution <-> Success Criteria aligned**: 15/30
Significant misalignment:
- In-scope lists "update resolveQuickDeps" but no SC tests quick mode.
- In-scope lists "default priority convention: api > web > cli > tui > mobile" but no SC tests this priority order specifically (only the 2-surface api+web case).
- In-scope lists "conflict detection" and SC4 covers it, but the error message format is unspecified.
- Solution says "execution-order config" but no SC tests the full config loading/validation path.
- The freeform review identifies that InferType must change — this is not in scope, not in solution, and not in SC, yet it is a prerequisite for the solution to work.
> Deduction: -15 for multiple scope/solution/SC gaps.

**Requirements <-> Solution coherent**: 15/25
- NFR "minimal change: only modify autogen.go dependency resolution and config.go" is incoherent with the solution's actual impact (InferType changes, naming scheme changes, migration, function signature changes).
- Scenario "multi-api surface" requires conflict detection, which is covered. But the solution provides no graceful degradation — the user MUST configure. This is a hard requirement masquerading as a scenario.
- The "backward compatibility" NFR is asserted but the solution introduces task ID changes that break compatibility with existing index.json. The freeform review's finding #4 (orphan risk) directly contradicts this NFR.
> Deduction: -10 for NFR/solution incoherence.

**D10 Total: 55/90**

---

## Phase 3: Blindspot Hunt

**[blindspot-1] No rollback plan.** The proposal has no mechanism to revert if the serial ordering causes longer pipeline times. If API tests take 2 hours and web tests take 10 minutes, users are forced to wait 2 hours before web tests even start. Previously they ran in parallel. The proposal optimizes for fail-fast but degrades the happy path. No measurement plan to validate that the trade-off is worthwhile.

**[blindspot-2] Hidden coupling to surface-aware-justfile proposal.** The urgency argument says "must land before justfile proposal." But if justfile proposal changes the task struct in incompatible ways after this proposal lands, there's no recourse. The proposal creates a one-way dependency without a fallback.

**[blindspot-3] gen-journeys merge is actually two changes, not one.** Merging gen-journeys requires: (a) changing the task generation loop, AND (b) changing the SKILL.md to handle multi-surface traversal. The proposal treats this as one item but the SKILL.md change is a behavior change that affects output quality. No validation criteria for output quality after merge.

**[blindspot-4] Serial execution is a performance regression for happy path.** The proposal assumes fail-fast is always desirable. But in many projects, API and web tests are independent — a web test failure doesn't mean API tests are invalid. The proposal removes the option to run them in parallel without providing evidence that serial is better in the majority case.

**[blindspot-5] The "user has vetoed" rejection of execution-level ordering is hearsay.** The proposal rejects the simpler alternative by appealing to a user's veto, but provides no context for WHY the user vetoed it or whether the user understood the complexity trade-off. This is a design decision outsourced to an unspecified conversation.

**[blindspot-6] No integration test strategy.** The proposal modifies the core autogen pipeline but lists no integration test plan. The success criteria are all unit-level (single config assertions). No end-to-end scenario tests are specified.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1. Problem Definition | 85 | 110 |
| D2. Solution Clarity | 90 | 120 |
| D3. Industry Benchmarking | 80 | 120 |
| D4. Requirements Completeness | 85 | 110 |
| D5. Solution Creativity | 35 | 100 |
| D6. Feasibility | 70 | 100 |
| D7. Scope Definition | 65 | 80 |
| D8. Risk Assessment | 60 | 90 |
| D9. Success Criteria | 55 | 80 |
| D10. Logical Consistency | 55 | 90 |
| **Total** | **680** | **1000** |

---

## Top Attacks (Priority Order)

1. **D6 Feasibility**: Technical feasibility assessment omits 3 critical changes (InferType contract, function signatures, dot-key sentinel) — "autogen.go already has per-surface loop logic, changing to single task only requires removing the loop" — must enumerate ALL affected code paths before claiming feasibility.

2. **D8 Risk Assessment**: 5 significant risks unaccounted for (index.json migration, InferType breakage, SourceTaskID migration, surface-key normalization, dot-key sentinel) — only 4 lower-severity risks listed — must conduct a full impact analysis of the codebase before claiming completeness.

3. **D10 Logical Consistency**: The "backward compatibility" NFR directly contradicts the solution's task ID renaming — "single surface project task structure and dependency chain completely identical to before" cannot be satisfied if task IDs change — must either prove ID stability or remove the NFR.

4. **D3 Industry Benchmarking**: Trade-off analysis is asymmetric — selected approach cons listed as only "needs to modify autogen.go" but freeform review identifies 11 additional impacts — comparison table must reflect actual implementation cost honestly.

5. **D5 Solution Creativity**: Zero novelty beyond industry baseline — "hybrid ordering strategy (convention + override)" is standard convention-over-configuration from 2005 — must either acknowledge this is straightforward adoption of industry patterns or find genuine innovation.

6. **D9 Success Criteria**: 6 in-scope items lack corresponding testable criteria (quick mode, InferType, migration, normalization, error messages, eval-journey dependency) — every in-scope item needs at least one success criterion.
