---
date: 2026-05-24
evaluator: CTO-adversary
iteration: baseline
model: rubric-adversarial-only
---

# Baseline Score: Surface Detection Fallback & Init Display Improvement

## Methodology

Scored against the proposal evaluation rubric (1000-point scale, 900 target). Every deduction includes a direct quote and a concrete improvement requirement. This evaluation treats the proposal as adversarial territory: ambiguities are counted as flaws, not charitable interpretations.

---

## 1. Problem Definition: 78/110

### Problem stated clearly: 34/40

The two-part problem (stdlib detection gap + opaque init summary) is identifiable. However, the framing packages two distinct issues without a dependency analysis between them.

> "This proposal addresses two distinct issues: (1) the detection gap for stdlib projects is the primary driver; (2) the init summary display opacity is a secondary UX problem that can ship independently but is included here because both affect the same `forge init` code path."

"Can ship independently" directly contradicts coupling them in one proposal. If they are independent, they should be separate proposals with separate evaluations. Combining them inflates scope without justified coupling.

**Must improve**: Either establish a hard dependency (init summary cannot be improved without inference) or split into two proposals.

### Evidence provided: 30/40

Three evidence points are concrete:

> "3 of 5 test subjects with stdlib Go projects asked 'what did it detect?' after seeing '(1 mappings)'"

This is the strongest piece. However:

- No verbatim user quotes
- "Node.js and Python stdlib projects exhibit the same UX pattern" is asserted by inference:

> "evidence is inferred from the Go testing rather than separately reproduced per ecosystem."

This is an honest admission but it means Python/Node evidence is hearsay, not data. Python inference is in scope without a single Python evidence point.

- The "4 of 10 init runs produced empty detection" is a specific metric, but the sample is internal projects (forge-cli itself), not representative of the target user base.

**Must improve**: Add at least one Python or Node.js specific evidence point, or narrow scope to Go-only for inference.

### Urgency justified: 14/30

> "Surface detection is the first interactive experience in `forge init`. A confusing or silent detection step undermines trust in the tool."

"First interactive experience" is generic urgency -- any init UX problem can claim this. Missing:

- What is the cost of delay? How many users are affected per week/month?
- Is there a workaround? (Yes: manual TUI entry already works, it's just not optimal)
- The urgency section conflates severity with urgency -- the problem is real, but why must it be solved *now* rather than next quarter?

> "In the last 10 `forge init` runs on internal projects, 4 produced empty detection results"

This is a frequency claim from a biased sample (internal projects are likely Go-centric). It does not establish urgency for the broader user base.

**Must improve**: Provide user-facing impact metrics (support tickets, onboarding drop-off data, or at minimum a timeline constraint like "blocks v3.0 release").

---

## 2. Solution Clarity: 85/120

### Approach is concrete: 35/40

Four numbered changes are clearly delineated. The structural inference concept is well-defined:

> "ecosystem-specific inference functions (e.g., `inferGoSurface`, `inferNodeSurface`, `inferPythonSurface`) only after all dependency signal maps return empty. Each inference function returns `(surfaceType, sourceAnnotation)` and is stateless"

This is specific enough to implement. Deduction: the inference rules themselves are scattered across Key Scenarios rather than consolidated in the Solution section. A reader must cross-reference scenarios 1-3 to understand what `inferGoSurface` actually checks.

**Must improve**: Add a summary table of inference rules in the Solution section (ecosystem | signal | surface type | exclusion condition).

### User-facing behavior described: 32/45

The before/after for init summary is clear:

> "CREATED surfaces cli (inferred from cmd/ directory structure)"

The TUI confirmation flow is described:

> "askSurfaceConfirmation will display the inference source in the Description field alongside each surface type"

But the description is prose, not a terminal mockup. There is no concrete example of what the TUI looks like when both dependency-detected and inference-detected surfaces appear in the same session. The `forge surfaces detect` command's non-interactive output format is not shown.

**Must improve**: Add terminal mockups for (a) TUI with mixed dependency+inference sources, (b) `forge surfaces detect --dry-run` output, (c) non-interactive stdout format.

### Technical direction clear: 18/35

> "inference functions live in `detect_surface.go` alongside existing signal maps"

Location is stated. Architecture is not. Key gaps:

1. **How inference integrates with `detectSurfaceAtDirWithConflicts`**: The current function collects signals and returns. Where do inference functions get called? After signal collection returns empty? The proposal says "only after all dependency signal maps return empty" but the current code has no hook point for this -- inference would require restructuring `detectSurfaceAtDirWithConflicts`.

2. **`DetectResult.Sources` schema**: The proposal proposes `Sources map[string]string` keyed by path. But the current `DetectResult` has no `Sources` field. Adding it affects all construction sites. The proposal mentions backward compatibility:

> "Zero-value of `map[string]string` is `nil`, so existing `DetectResult` construction sites are backward-compatible without changes"

This is correct for deserialization but ignores the semantic contract: any code that reads `Sources` must handle the case where it is nil (dependency detection) vs populated (inference). This is not a breaking change but it is a behavioral change that requires updating all readers.

3. **Conflict between `Sources` map keys and `Surfaces` map keys**: If `Sources` is keyed by path and `Surfaces` is also keyed by path, but `Sources` only contains entries for inferred surfaces (not dependency-detected ones), the key space is inconsistent. The proposal does not address whether dependency-detected surfaces also get a `Sources` entry.

**Must improve**: Add a sequence diagram or pseudocode for the modified `detectSurfaceAtDirWithConflicts` showing where inference fires. Define `Sources` population rules for both dependency and inference paths.

---

## 3. Industry Benchmarking: 75/120

### Industry solutions referenced: 28/40

The proposal references four industry examples:

> "npm init infers project type from existing `package.json` fields before prompting"
> "cargo init distinguishes `--bin` (CLI) from `--lib` (library) using a single flag"
> "VS Code's workspace detection reads `.vscode/settings.json` then falls back to probing directory conventions (e.g., `manage.py` -> Django)"
> "Spring Boot's `spring.factories` auto-configuration uses a layered classpath scan that degrades gracefully"

Four references is adequate. However, each gets 1-2 sentences of description with no analysis of what Forge can learn from their specific implementation. The VS Code reference is the most relevant (layered detection with graceful degradation) but is mentioned in passing.

**Must improve**: For at least 2 references, describe the specific mechanism and how it maps to Forge's architecture.

### At least 3 meaningful alternatives: 15/30

Four alternatives listed:

1. **Do nothing** -- genuine alternative
2. **Expand signal maps only** -- this is a subset of the existing approach, not a genuinely different strategy
3. **Regex-based import scanning** -- presented only to be rejected:

> "Rejected: couples init to source code parsing; breaks 'no AST' constraint"

This is a straw man. The alternative is defined specifically to violate a stated constraint, making rejection inevitable. A fairer framing would evaluate regex scanning as a legitimate approach that *might* justify relaxing the constraint.

4. **Structural inference + display fix** -- selected

**Deduction: -20 for straw-man alternative per rubric rules.**

**Must improve**: Replace "regex-based import scanning" with a genuine third alternative (e.g., "user telemetry-driven detection" or "community-maintained detection rules") or evaluate import scanning on its merits before citing the constraint.

### Honest trade-off comparison: 17/25

The comparison table includes pros, cons, and verdicts. The selected approach's "cons" column says:

> "Inference is probabilistic (annotated in output)"

This is honest. However, the "Expand signal maps only" alternative is dismissed with:

> "Doesn't solve stdlib case, endless whack-a-mole"

"Endless whack-a-mole" is inflammatory without evidence. Adding 5-10 more signal entries (e.g., `net/http`, `gin`, `fiber` for Go API detection) is a bounded task, not an endless one.

**Must improve**: Provide honest characterization of the signal-maps-only approach. Acknowledge it could solve 80% of remaining cases with minimal effort.

### Chosen approach justified against benchmarks: 15/25

> "**Selected: best ROI** -- 3-4 inference functions at ~20 LOC each yield coverage for the most common stdlib cases"

The ROI is asserted with LOC estimates but no justification for why ~80 LOC of inference code is better ROI than ~20 LOC of expanded signal maps. The comparison table's verdict column is bare assertion.

**Must improve**: Quantify ROI: how many more projects does structural inference cover vs expanded signal maps? What is the overlap?

---

## 4. Requirements Completeness: 82/110

### Scenario coverage: 30/40

11 scenarios listed, which is thorough. Key scenarios covered:

- Go stdlib with `cmd/` -> cli (scenario 1)
- Conflict resolution: `cmd/` + `api/` -> api wins (scenario 1)
- Re-run with existing config (scenario 7)
- User overrides inference (scenario 8)
- Non-interactive terminal (scenario 11)

Missing scenarios:

1. **Workspace mode with mixed sources**: subdir A has dependency signals, subdir B has only structural inference. How does `Sources` map reflect this?
2. **Inference returns empty**: what does the user see? Same as today? The proposal does not describe the "all methods failed" UX.
3. **`forge surfaces detect` with existing config**: does it prompt for overwrite? The proposal says it "writes to config on confirm" but does not describe what happens if config already has surfaces.

**Must improve**: Add scenarios for workspace mixed-source detection and the "all inference failed" fallback path.

### Non-functional requirements: 30/40

NFRs listed:

> "Structural inference must be fast (no external network calls, no AST parsing)"
> "Inference performance budget: all inference functions combined must complete in <50ms"
> "Inference functions return empty on any filesystem error, never crash `forge init`"

The <50ms budget is quantified. Good. Missing:

- **Accuracy target**: no minimum accuracy or confidence threshold stated. The Assumptions section mentions "annotated 80% confidence" but this is informal, not an NFR.
- **No localization requirement** for new user-facing strings (inference source annotations, TUI hints).

**Must improve**: Add accuracy NFR (e.g., "inference must agree with manual classification in >= X% of tested projects") or explicitly state that accuracy is best-effort with annotation.

### Constraints & dependencies: 22/30

Constraints are listed and reasonable:

> "No new dependencies -- pure filesystem checks"
> "Inference rules are per-ecosystem (Go, Node.js, Python) -- each needs its own heuristic set"

The constraint "no AST parsing" is clear and consistently maintained throughout the proposal. The proposal correctly identifies no external dependencies.

Missing:

- The proposal does not acknowledge that `Sources map[string]string` is a schema addition to `DetectResult` that downstream code must handle.
- No mention of how inference interacts with the `FORGE_DETECT_DEPTH` env var (does structural inference respect the depth limit?).

**Must improve**: Add constraint for `DetectResult` schema evolution and clarify depth-limit interaction.

---

## 5. Solution Creativity: 50/100

### Novelty over industry baseline: 20/40

The proposal is explicit:

> "No special innovation -- this is standard convention-over-configuration"

Honest but earns limited novelty credit. The differentiation from the industry baseline (VS Code layered detection, Spring Boot auto-config) is the specificity of the inference rules per ecosystem. However, the proposal does not articulate *how* its approach differs from or improves upon these established patterns.

> "The insight is that **project structure is a reliable secondary signal**"

This insight is correct but not novel -- it is the same insight behind every convention-over-configuration system.

**Must improve**: Either acknowledge this is pure industry-standard application (and accept the score) or identify one specific way Forge's inference differs from the benchmarked tools.

### Cross-domain inspiration: 10/35

No cross-domain inspiration is demonstrated. The proposal stays within Go/Node/Python ecosystem conventions. No borrowing from:

- IDE auto-completion confidence scoring (for inference confidence)
- Spell-checker "did you mean?" patterns (for suggesting corrections)
- Search engine query suggestion patterns (for inferring likely intent from partial signals)

**Must improve**: Identify at least one cross-domain inspiration or explicitly state that this is a straightforward application with no cross-domain borrowing.

### Simplicity of insight: 20/25

The insight is genuinely elegant: "when dependency signals fail, project directory structure is a cheap and reliable secondary signal." The "why didn't I think of that" quality is present. The constraint that inference *only* activates when all dependency signals are empty keeps the design clean.

Deduction: the implementation details (per-path `Sources` map, conflict between inference candidates like `cmd/` vs `api/`, exclusion conditions for Python library projects) add complexity that partially obscures the simplicity.

**Must improve**: The simplicity of the insight is the proposal's strongest creative asset. Lean into it by keeping the implementation rules minimal and composable.

---

## 6. Feasibility: 78/100

### Technical feasibility: 30/40

All inference is filesystem-based -- technically straightforward. No external dependencies. The existing `detectSurfaceAtDirWithConflicts` function provides a clear integration point.

Concern: The proposal says inference functions are called "only after all dependency signal maps return empty." The current code's `detectSurfaceAtDirWithConflicts` returns immediately when `allSignals` is empty:

```go
if len(allSignals) == 0 {
    return "", nil
}
```

Integrating inference requires modifying this return path to call the inference functions before returning empty. This is straightforward but the proposal does not describe this code change.

**Must improve**: Show the modified control flow in `detectSurfaceAtDirWithConflicts` or acknowledge the integration point explicitly.

### Resource & timeline feasibility: 25/30

> "3-4 coding tasks"

Listing 5 tasks (including `forge surfaces detect`):

1. Inference rules per ecosystem
2. Sources map integration
3. TUI confirmation flow update
4. Init summary display + tests
5. `forge surfaces detect` command

Five tasks is realistic but the proposal says "3-4" then lists 5. The inconsistency suggests scoping has not been fully thought through.

**Must improve**: Reconcile the task count. If 5 tasks, say 5. If 3-4, explain which tasks are merged.

### Dependency readiness: 23/30

Correctly states no external dependencies. Uses existing `DetectSurfacesWithConflicts` infrastructure.

Unstated dependency: ecosystem convention stability. If Go community conventions shift away from `cmd/` directory structure (unlikely but possible), the inference rules become wrong. This is an unstated external dependency.

**Must improve**: Acknowledge convention stability as an implicit dependency.

---

## 7. Scope Definition: 65/80

### In-scope items are concrete: 23/30

Four in-scope items:

1. Structural inference fallback for Go, Node.js, Python
2. Init summary format improvement
3. Per-path inference source annotation in `DetectResult` via `Sources map[string]string`
4. `forge surfaces detect` subcommand

Item 3 ("Sources map[string]string field") is an implementation detail, not a user-facing deliverable. A better framing: "Users can see inference source in TUI and init summary."

Item 4 (`forge surfaces detect`) is a new CLI subcommand that significantly expands scope beyond the stated "small-scope improvement." A new subcommand with dry-run, non-interactive mode, and config-write behavior is not trivial.

**Must improve**: Frame in-scope items as deliverables observable by users/maintainers, not implementation details. Acknowledge that `forge surfaces detect` is scope-significant.

### Out-of-scope explicitly listed: 22/25

Five items listed as out of scope:

> "Expanding dependency signal maps (separate effort)"
> "Rust/Cargo.toml structural inference"
> "Import-level analysis"
> "Changes to `forge surfaces` CLI command output format"
> "Mobile/TUI project structural inference (too rare to heuristic)"

The Rust exclusion has a clear rationale:

> "Rust project surface detection is not yet implemented in the existing dependency signal maps; adding inference for an ecosystem with no baseline detection support would be premature"

This is a valid technical argument. However, it creates an inconsistency: Cargo.toml *does* have dependency signal maps in the codebase (`cargoTomlSignals`), so the claim "not yet implemented" is partially incorrect.

**Must improve**: Correct the Rust exclusion rationale -- Cargo.toml has signal maps but Rust structural inference is out of scope because the signal maps already cover the common cases (`clap`, `structopt` for CLI; `actix-web`, `axum` for API).

### Scope is bounded: 20/25

> "Proceed to `/quick` (this is a small-scope improvement, no PRD needed)"

The assertion of small scope is contradicted by the 5-task list and a new CLI subcommand. `/quick` is appropriate for 1-15 coding tasks, so the scope is bounded within the forge process, but calling it "small" when it includes a new subcommand is understatement.

**Must improve**: Either remove the `forge surfaces detect` subcommand from scope or acknowledge it as medium scope.

---

## 8. Risk Assessment: 62/90

### Risks identified: 20/30

Five risks listed:

1. Structural inference gives wrong type (M/M)
2. Inference rules become stale (L/L)
3. Init summary format change breaks scripted consumers (L/M)
4. `Sources map[string]string` schema migration (L/M)
5. Re-run prompt adds friction (L/L)

Missing risks:

- **Mixed-source workspace detection**: when some subdirs have dependency signals and others have inference, the `Sources` map has inconsistent key coverage. No risk identified for this.
- **Inference-conflict precedence**: scenario 1 says `api/` beats `cmd/` for Go, but this precedence is a policy decision with no risk assessment. What if a Go project has both a CLI and an API?
- **`forge surfaces detect` overwrites user config**: the command "writes to config on confirm" but there is no risk assessment for accidental overwrite of manually-configured surfaces.
- **Trust impact of wrong inference**: a wrong automated detection is *worse* than no detection because it confirms a false belief. The proposal rates impact as M but the trust dimension is underrated.

**Must improve**: Add risks for workspace mixed-source handling, inference-conflict precedence, and trust impact of wrong inferences.

### Likelihood + impact rated: 20/30

The ratings are mostly reasonable. One notable issue:

> "Structural inference gives wrong type" -- Likelihood: M, Impact: M

Impact should be higher. When Forge tells a user "your project is a CLI" and the user believes it (because the tool said so), the downstream surface configuration is wrong. This affects test generation, task templates, and the entire workflow. Wrong automated detection actively *reduces* trust compared to asking the user.

The `Sources` schema migration risk is rated L/M but the proposal itself identifies it needs:

> "add compile-time check that all code paths populate `Sources` when inference fires"

This mitigation implies the risk is higher than L/M -- if a compile-time check is needed, the risk of missing a code path is real.

**Must improve**: Upgrade "wrong inference" impact to H. Upgrade "Sources schema migration" likelihood to M.

### Mitigations are actionable: 22/30

Most mitigations are actionable:

- "Annotate source in output" -- actionable
- "Add integration test suite covering each inference rule" -- actionable
- "Ship with verbose logging (`FORGE_DEBUG=1`)" -- actionable
- "Zero-value of `map[string]string` is `nil`" -- this is a property, not a mitigation
- "Rules are simple filesystem checks, easy to update" -- this is not a mitigation, it is a design property

Missing mitigation: no rollback plan. If structural inference causes problems in production, how is it disabled? There is no feature flag, no env var to disable inference, no way to revert behavior without a code release.

**Must improve**: Add a rollback mechanism (e.g., `FORGE_DISABLE_INFERENCE=1` env var or config field).

---

## 9. Success Criteria: 68/80

### Criteria are measurable and testable: 48/55

15 checkbox criteria. Most are testable:

> "Go project with only `cmd/` subdirs -> detected as `cli` without manual entry" -- testable with a fixture project

> "Go project with both `cmd/` and `api/` directories -> only `api` returned" -- testable

> "`forge init` summary shows actual surface types: `CREATED surfaces cli`" -- testable with output capture

> "`DetectResult.Sources` map correctly populated: `{\"forge-cli/cli\": \"inference:cmd-dir\"}`" -- testable

Weak criteria:

> "Existing detection tests (27 tests) still pass unchanged" -- "unchanged" is ambiguous. Does it mean the test code is unchanged, or the test results are unchanged?

> "User override: after editing an inferred `cli` to `api` in TUI, config contains `surfaces: api`" -- testable but requires TUI automation, which the proposal does not discuss.

Missing criteria:

- No criterion for the `forge surfaces detect` command's non-interactive mode exit codes (exit code 1 on no surfaces found is mentioned in scenarios but not in success criteria).
- No criterion for inference performance (<50ms NFR has no corresponding success criterion).
- No criterion for Python-specific inference despite Python being in scope.

**Must improve**: Add success criteria for inference performance budget, Python-specific inference, and non-interactive exit codes.

### Coverage is complete: 20/25

The 15 criteria cover most in-scope items. Gaps:

- Python inference (`[project.scripts]`, `app.py` with library exclusion) has no dedicated criterion despite being in scope.
- The `forge surfaces detect` command's core behavior (TUI confirmation + config write) has a criterion but the non-interactive mode behavior does not have a dedicated testable criterion.
- The "re-run with existing config" behavior has a criterion but it is complex (combining multiple behaviors into one checkbox).

**Must improve**: Add Python-specific success criteria matching the Python inference rules described in Key Scenario 3.

---

## 10. Logical Consistency: 72/90

### Solution addresses the stated problem: 28/35

The init summary fix directly addresses problem (2). The structural inference addresses problem (1) for Go projects effectively. For Python:

> "Python coverage is limited to CLI detection via `[project.scripts]`/`setup.py entry_points`/`app.py` with library exclusion; Python API-type stdlib projects without Flask/FastAPI remain undetectable by structural conventions alone"

This is honest. The proposal acknowledges it does not fully solve stdlib detection for Python. The question is: does partial coverage justify the implementation cost? The proposal argues yes via "annotated 80% confidence beats silent 0%." This is a defensible position.

Gap: the proposal claims to solve "stdlib project detection" as the primary driver but only handles Go CLI/API and Python CLI. Node.js stdlib API projects (e.g., Express without the Express dependency) are not addressed. The problem statement is broader than the solution.

**Must improve**: Narrow the problem statement to match the solution's actual coverage: "Go and Python CLI stdlib detection" rather than "stdlib project detection" as a universal claim.

### Scope <-> Solution <-> Success Criteria aligned: 24/30

- Python is in scope but has no dedicated success criterion -- misalignment.
- `forge surfaces detect` is in scope and has success criteria -- aligned.
- "Inference source annotation in DetectResult via Sources map" is in scope and has a success criterion for `Sources` map population -- aligned.
- The Rust exclusion rationale ("not yet implemented") is contradicted by the existence of `cargoTomlSignals` in the codebase -- factual inconsistency.

**Must improve**: Fix Rust exclusion rationale. Add Python success criteria.

### Requirements <-> Solution coherent: 20/25

Most requirements map cleanly to the solution. One orphan requirement:

> "Inference confidence is lower than dependency signals -- annotate output to distinguish inference vs detection"

The solution provides source annotation in the `Sources` map. However, there is no requirement or solution for *measuring* confidence -- the proposal asserts "80% confidence" without defining how confidence is assessed. This is an orphan assumption masquerading as a requirement.

The priority chain requirement is well-specified:

> "Inference signals **only activate when ALL dependency signals are empty** -- they never mix with dependency results."

This is clean and testable. But it creates a gap in workspace mode: subdir A has dependency signals (no inference), subdir B has only structural signals (inference fires). The `Sources` map would only have entries for subdir B, creating an asymmetry. The proposal does not address this.

**Must improve**: Define workspace-mode behavior for mixed dependency/inference sources. Either Sources contains entries for all paths (including dependency-detected ones) or the TUI must handle the case where some paths have source annotations and others do not.

---

## Score Summary

```
SCORE: 713/1000
DIMENSIONS:
  Problem Definition: 78/110
  Solution Clarity: 85/120
  Industry Benchmarking: 75/120
  Requirements Completeness: 82/110
  Solution Creativity: 50/100
  Feasibility: 78/100
  Scope Definition: 65/80
  Risk Assessment: 62/90
  Success Criteria: 68/80
  Logical Consistency: 72/90
ATTACKS:
1. [Industry Benchmarking]: Straw-man alternative -- "Regex-based import scanning" is defined specifically to violate the "no AST" constraint, making rejection inevitable. A fairer evaluation would assess it on merits before citing the constraint. Deduction: -20 per rubric rules. Quote: "Rejected: couples init to source code parsing; breaks 'no AST' constraint"
2. [Industry Benchmarking]: ROI assertion without quantification -- "best ROI" is claimed but never justified with numbers. How many more projects does structural inference cover vs expanded signal maps? What is the overlap? Quote: "**Selected: best ROI** -- 3-4 inference functions at ~20 LOC each"
3. [Problem Definition]: Scope coupling without dependency -- the proposal admits the init summary fix "can ship independently" but bundles it with structural inference. If they are independent, they should be separate proposals. Quote: "(2) the init summary display opacity is a secondary UX problem that can ship independently"
4. [Problem Definition]: Python/Node evidence is hearsay -- the proposal admits evidence for non-Go ecosystems is "inferred from the Go testing rather than separately reproduced." Python inference is in scope without a single Python-specific evidence point. Quote: "evidence is inferred from the Go testing rather than separately reproduced per ecosystem"
5. [Solution Clarity]: Inference integration point not described -- the current `detectSurfaceAtDirWithConflicts` returns empty immediately when signals are empty; the proposal does not describe how inference hooks into this code path. Quote: "Inference functions live in `detect_surface.go` alongside existing signal maps" -- this tells location, not control flow.
6. [Solution Clarity]: No terminal mockups for TUI behavior -- the TUI is the primary user touchpoint but only described in prose. Mixed dependency+inference sources, `forge surfaces detect --dry-run` output, and non-interactive stdout have no visual examples. Quote: "askSurfaceConfirmation will display the inference source in the Description field"
7. [Solution Clarity]: Sources map key inconsistency -- if `Sources` only contains entries for inferred paths, the key space is inconsistent with `Surfaces` (which contains all paths). Dependency-detected paths have no Sources entry, creating an implicit "no entry means dependency" convention. Quote: "Results are merged into the existing `DetectResult.Sources` map keyed by path"
8. [Requirements Completeness]: Missing workspace mixed-source scenario -- in workspace mode, subdir A may have dependency signals while subdir B only has structural inference. The proposal does not describe this case. Quote: (no mention of workspace mode in scenarios despite `detectWorkspaceMode` existing in codebase)
9. [Requirements Completeness]: No accuracy NFR -- the proposal mentions "80% confidence" informally in Assumptions but never establishes it as a measurable requirement. Quote: "Occam's Razor: annotated 80% confidence beats silent 0%"
10. [Scope Definition]: Task count inconsistency -- proposal says "3-4 coding tasks" then lists 5 tasks. Quote: "3-4 coding tasks: (1) inference rules per ecosystem, (2) Sources map integration, (3) TUI confirmation flow update, (4) init summary display + tests. Additionally: (5) `forge surfaces detect` command"
11. [Scope Definition]: Rust exclusion rationale is factually incorrect -- the proposal says "Rust project surface detection is not yet implemented in the existing dependency signal maps" but the codebase contains `cargoTomlSignals` with entries for actix-web, axum, rocket, clap, ratatui, etc. Quote: "Rust project surface detection is not yet implemented in the existing dependency signal maps; adding inference for an ecosystem with no baseline detection support would be premature"
12. [Risk Assessment]: No rollback mechanism -- if structural inference causes problems in production, there is no feature flag, env var, or config option to disable it without a code release. Quote: (no mention of rollback or disable mechanism anywhere in the proposal)
13. [Risk Assessment]: Wrong inference trust impact underrated -- "Structural inference gives wrong type" rated M impact, but wrong automated detection is worse than no detection because it confirms a false belief and erodes tool trust. Quote: "Structural inference gives wrong type | M | M | Annotate source in output"
14. [Success Criteria]: No Python-specific success criterion -- Python inference is in scope (3 rules: `[project.scripts]`, `setup.py entry_points`, `app.py` with library exclusion) but no success criterion tests any of these rules. Quote: (success criteria cover Go scenarios and Node.js `bin`/`index.html` but not Python)
15. [Success Criteria]: Inference performance has no success criterion -- the <50ms NFR is stated but no checkbox criterion verifies it. Quote: "Inference performance budget: all inference functions combined must complete in <50ms"
16. [Logical Consistency]: Problem statement broader than solution coverage -- the problem claims "stdlib project detection" generally, but the solution only covers Go CLI/API and Python CLI. Node.js stdlib API projects are not addressed. Quote: "Surface detection relies solely on framework dependency signals. When a project's `go.mod`/`package.json`/`pyproject.toml` lacks recognized frameworks, detection returns empty"
17. [Logical Consistency]: Workspace Sources map asymmetry not addressed -- when some subdirs use dependency detection and others use inference, the Sources map has inconsistent coverage. The proposal does not describe this behavior. Quote: "Inference signals **only activate when ALL dependency signals are empty** -- they never mix with dependency results" -- this assertion is about per-call behavior, but workspace mode calls detection per-subdir
```

## Verdict

**713/1000 -- BELOW 900 TARGET.**

The proposal identifies a genuine UX gap and proposes a technically sound approach. The strongest aspects are the specific inference rules with exclusion conditions, the priority chain enforcement (inference only when dependency signals are empty), and the honest acknowledgment of Python coverage limitations. The weakest aspects are the straw-man alternative in industry benchmarking, missing rollback mechanism, factual error in the Rust exclusion rationale, and the gap between problem breadth (all stdlib projects) and solution coverage (Go + Python CLI only).

### Priority Revisions

1. **Add rollback mechanism**: `FORGE_DISABLE_INFERENCE=1` env var or equivalent -- this is a hard requirement for production safety
2. **Fix Rust exclusion rationale**: `cargoTomlSignals` exists in codebase; reframe as "signal maps already cover common Rust patterns"
3. **Replace straw-man alternative**: evaluate regex/import scanning on merits or substitute a genuine third alternative
4. **Add Python success criteria**: three rules described in scenarios, zero tested in success criteria
5. **Add inference performance success criterion**: the <50ms NFR has no corresponding checkbox
6. **Narrow problem statement**: change "stdlib project detection" to "Go and Python CLI stdlib detection" to match actual coverage
7. **Add workspace mixed-source scenario**: describe Sources map behavior when subdirs have different detection methods
8. **Reconcile task count**: the proposal says 3-4 tasks but lists 5
