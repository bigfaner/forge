---
date: 2026-05-24
evaluator: CTO-adversary
iteration: 1
model: rubric-adversarial
basis: iteration-0-report.md
pre_revision_annotations: true
---

# Iteration 1: Adversarial Rubric Scoring

## Pre-Revision Annotation Tracking

### Annotated Regions (3 paragraphs)

| # | Location | Severity | Revision Summary |
|---|----------|----------|-----------------|
| A1 | Proposed Solution #4 (line 42) | high | Added `forge surfaces detect` subcommand with `--apply` flag and non-interactive mode |
| A2 | Key Scenario #2 (line 54) | medium | Clarified `inferNodeSurface` does not scan subdirectories for `index.html` |
| A3 | Key Scenario #7 (line 60) | high | Expanded re-run prompt to three options: Confirm / Re-detect / Edit |

### Unannotated Regions

All remaining content (~25 paragraphs): Problem, Evidence, Urgency, Proposed Solution #1-3, Innovation Highlights, Key Scenarios #1/#3-#6/#8-#11, NFRs, Constraints, Industry Solutions, Comparison Table, Feasibility, Assumptions, Scope, Key Risks, Success Criteria, Next Steps.

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

**Problem**: Two issues stated: (1) dependency-signal detection returns empty for stdlib projects with silent fallback to manual entry; (2) init summary shows opaque "N mappings".

**Solution**: (1) Structural inference fallback; (2) improved init summary format; (3) TUI confirmation showing inference source; (4) `forge surfaces detect` subcommand.

**Chain integrity**: The primary problem (stdlib detection gap) maps cleanly to Solution #1. The secondary problem (opaque summary) maps to Solution #2. Solutions #3 and #4 are UX enhancements that go beyond the stated problem — they address implicit pain but expand scope. The chain holds for the core problem but has scope creep toward the edges.

Notable: The problem statement says "both affect the same `forge init` code path" as justification for bundling, but Solution #4 (`forge surfaces detect`) is explicitly independent of `forge init`. This contradicts the bundling rationale.

### Solution -> Evidence Trace

Three evidence points: (1) forge-cli itself — detection works but summary is opaque; (2) pure stdlib Go project — detection returns empty; (3) user testing — 3/5 subjects confused.

The evidence supports Go detection issues. For Node.js and Python, the proposal admits: "evidence is inferred from the Go testing rather than separately reproduced per ecosystem." This is a significant gap — the solution covers three ecosystems but evidence covers one.

### Success Criteria -> Solution Trace

17 success criteria listed. Coverage analysis:
- Solution #1 (inference fallback): criteria for Go (#1-#3), Node.js (#4), Python (#5), priority chain (#10) — covered
- Solution #2 (init summary): criteria #6, #7 — covered
- Solution #3 (TUI confirmation): criterion #8 — covered
- Solution #4 (`forge surfaces detect`): criteria #13-#15 — covered
- Re-run behavior: criteria #11, #16 — covered
- YAML comment persistence: criterion #12 — covered
- User override: criterion #16 — covered

This is more complete than the iteration-0 assessment suggested. However, gaps remain:
- No criterion for performance budget (<50ms NFR)
- No criterion for filesystem error handling (NFR: "return empty on any filesystem error")
- No criterion for non-interactive terminal output format

## Phase 2: Rubric Scoring

### 1. Problem Definition: 76/110

**Problem stated clearly: 33/40**

The two-part problem is clearly articulated. Deduction: the framing bundles two distinct problems without rigorous prioritization.

> "This proposal addresses two distinct issues: (1) the detection gap for stdlib projects is the primary driver; (2) the init summary display opacity is a secondary UX problem"

The text does label (1) as primary and (2) as secondary. However, the bundling justification — "both affect the same `forge init` code path" — is undermined by Solution #4, which adds a separate `forge surfaces detect` command entirely independent of `forge init`. If the code-path argument justified bundling, adding a command outside that code path invalidates the argument.

**Evidence provided: 28/40**

Three concrete examples, but with notable gaps:

> "Internal testing (forge v2.x onboarding walkthrough, 2026-04): 3 of 5 test subjects with stdlib Go projects asked 'what did it detect?'"

This is the strongest evidence: specific test context, sample size, outcome. But:
- 3/5 subjects with Go projects — what about the 2 who didn't ask? What did they do?
- "Node.js and Python stdlib projects exhibit the same UX pattern; evidence is inferred from the Go testing rather than separately reproduced per ecosystem" — this is an explicit admission of missing evidence for 2 of 3 covered ecosystems. The solution commits to inference functions for all three, but evidence only supports one.

> "In the last 10 `forge init` runs on internal projects, 4 produced empty detection results"

Concrete data point. But: "internal projects" may not represent the target user base (internal projects likely skew toward stdlib/CLI patterns). No context on what those projects were.

**Urgency justified: 15/30**

> "A confusing or silent detection step undermines trust in the tool."

This is an assertion, not a cost analysis. "Undermines trust" is plausible but unquantified.

> "Each failed detection adds friction that risks user abandonment at the onboarding stage."

"Risks user abandonment" is stated without evidence of actual abandonment. No churn data, no support tickets, no user complaints citing detection as a reason for dropping the tool. The urgency argument rests on hypothetical harm, not demonstrated harm.

No workaround described — can the user simply type the surface type manually (which is the current fallback)? If so, the cost of delay is "users must type surface types manually until this ships," which is inconvenience, not a blocker.

### 2. Solution Clarity: 90/120

**Approach is concrete: 36/40**

Four clearly delineated changes, each with specific behavior described. A reader can explain back what will be built. Deduction: the inference rules themselves are defined in Requirements Analysis, not in the Solution section — a reader must cross-reference to understand what "structural inference" actually checks.

> "Each inference function returns `(surfaceType, sourceAnnotation)` and is stateless — pure filesystem reads with no side effects."

Good: signature, purity constraint, and side-effect guarantee stated. But no description of how inference functions are dispatched (by ecosystem? by manifest file type?).

**User-facing behavior described: 37/45**

Concrete before/after example for init summary:

> Before: `CREATED surfaces (1 mappings)`
> After: `CREATED surfaces cli (inferred from cmd/ directory structure)`

TUI confirmation flow described with specific display text:

> `"cli (inferred from cmd/ directory structure)"` vs `"cli (detected from cobra dependency)"`

Re-run TUI prompt described with three options:

> `"Surfaces already configured: cli. Re-detect?"` with Confirm / Re-detect / Edit

Non-interactive mode described:

> "Non-interactive terminals print results to stdout and exit (no TUI, no config write regardless of flags)"

Gaps:
- No terminal output mockup for `forge surfaces detect` command
- No description of what the human-readable stdout format looks like in non-interactive mode
- No description of error output (what if detection fails entirely?)

**Technical direction clear: 17/35**

> "Inference functions live in `detect_surface.go` alongside existing signal maps"

Location, not architecture.

> "`detectSurfaceAtDirWithConflicts` calls ecosystem-specific inference functions (e.g., `inferGoSurface`, `inferNodeSurface`, `inferPythonSurface`) only after all dependency signal maps return empty"

This is the clearest architecture statement, but it is buried in Solution #1. Missing:
- How does dispatch work? Does `detectSurfaceAtDirWithConflicts` detect the ecosystem from the manifest file and call the appropriate inference function?
- How are results merged into `DetectResult.Sources`? Is it additive or does inference only populate Sources when dependency signals are empty?
- What is the call chain from `forge surfaces detect` through to `DetectSurfacesWithConflicts`?

### 3. Industry Benchmarking: 62/120

**Industry solutions referenced: 22/40**

Four tools referenced: `npm init`, `cargo init`, VS Code workspace detection, Spring Boot `spring.factories`. These are real tools with real detection strategies.

> "VS Code's workspace detection reads `.vscode/settings.json` then falls back to probing directory conventions (e.g., `manage.py` -> Django)"

This is a concrete, relevant example. But:
- No depth — one sentence per tool. No analysis of how each tool's layered detection actually works.
- No reference to any published pattern (e.g., "convention over configuration" as a named pattern, Chain of Responsibility pattern).
- The Spring Boot reference is tangential — auto-configuration is a different problem domain than project-type detection.

**At least 3 meaningful alternatives: 15/30**

Four alternatives listed: "Do nothing", "Expand signal maps only", "Regex-based import scanning", "Structural inference + display fix".

"Do nothing" is genuinely different. "Expand signal maps only" is a subset of the status quo, not a genuinely different approach — it's "do the same thing but more," which is incremental, not alternative.

> "Regex-based import scanning | grep/sed-style scan of source files | Covers more cases than filesystem checks (e.g., detects `import http` usage) | Requires reading source files, fragile across languages, false positives from unused imports"

This is a meaningful alternative. The cons are honestly stated. However, the verdict "Rejected: couples init to source code parsing; breaks 'no AST' constraint" — the "no AST" constraint is self-imposed and circular when used to reject alternatives. Regex scanning is not AST parsing.

**Selected: best ROI** — "best ROI" is a label without analysis. What is the investment? What is the return? No numbers.

**Honest trade-off comparison: 10/25**

The comparison table has Pros/Cons/Verdict columns but:
- No quantified trade-offs
- No comparison against success criteria
- "Covers stdlib case, improves all detection UX" for the selected approach lists only pros in the Cons column: "Inference is probabilistic (annotated in output)" — this is the only con, and it is immediately dismissed by annotation

**Chosen approach justified against benchmarks: 15/25**

> "**Selected: best ROI** -- 3-4 inference functions at ~20 LOC each yield coverage for the most common stdlib cases (Go `cmd/`, Node `bin`, Python `[project.scripts]`) vs. regex scanning which requires per-language import syntax handling (~50 LOC per language) with higher false-positive rates"

This is the first quantified comparison (LOC estimate). It is a genuine attempt at justification. However:
- LOC is not the only cost dimension (maintenance, false-positive rate, test coverage)
- The ROI claim is relative to regex scanning only, not to "expand signal maps" or "do nothing"
- No user-impact analysis (how many users does each approach serve?)

### 4. Requirements Completeness: 82/110

**Scenario coverage: 32/40**

11 key scenarios listed — this is thorough. Coverage includes:
- Go stdlib (scenario #1)
- Node.js minimal (scenario #2)
- Python minimal (scenario #3)
- Subdir detection (#4)
- Detection with signals (#5)
- No manifest (#6)
- Re-run with existing config (#7)
- User override (#8)
- `forge surfaces detect` read-only (#9)
- `forge surfaces detect --apply` (#10)
- Non-interactive terminal (#11)

Missing:
- Scenario for multiple inference rules firing simultaneously (e.g., Go project with both `cmd/` and `api/`) — scenario #1 partially addresses this ("if both `cmd/` and `api/` exist -> `api` wins and only `api` is returned"), but the discard behavior (losing candidate discarded entirely, not stored) is not explained in the scenario
- Scenario for filesystem errors during inference (NFR says "return empty on any filesystem error" but no scenario demonstrates this)
- Scenario for mixed ecosystem project (e.g., monorepo with both Go and Python)

**Non-functional requirements: 30/40**

Five NFRs stated:

> "Structural inference must be fast (no external network calls, no AST parsing)"

Good: clear constraint. But "fast" is vague — the performance budget NFR clarifies (<50ms), but it's listed separately.

> "Inference performance budget: all inference functions combined must complete in <50ms (filesystem stat + directory listing only, no file content reads beyond manifest parsing already done by dependency detection)"

Good: specific, measurable budget. The clarification "no file content reads beyond manifest parsing" is valuable.

> "Inference confidence is lower than dependency signals -- annotate output to distinguish inference vs detection"

Good: addresses the probabilistic nature.

> "Backward compatible: existing configs and workflows unchanged"

Asserted but not verified — the `Sources map[string]string` addition is a schema change.

> "Inference functions return empty on any filesystem error (permission denied, unreadable directory), never crash `forge init`"

Good: error handling constraint.

Missing NFRs:
- Accessibility: new TUI text strings — any i18n consideration?
- Security: inference functions read the filesystem — any path traversal risk from user-controlled directory paths?
- Accuracy target: no stated minimum accuracy for inference rules

**Constraints & dependencies: 20/30**

Five constraints listed. The "Source persistence via YAML comments" constraint is well-thought-out:

> "when writing surfaces config, append a comment line like `surfaces: api  # source: inference:cmd-dir` for human context. This preserves source metadata across re-runs without modifying the SurfacesMap schema"

This is a creative, lightweight approach to source persistence. However:
- YAML round-tripping is fragile — many YAML libraries strip comments. The proposal does not address which YAML library forge uses or whether it preserves comments.
- The constraint "Re-running `forge init` with existing surfaces config" is described in detail but is really a requirement, not a constraint.
- The constraint about "no new dependencies" is well-stated.

### 5. Solution Creativity: 48/100

**Novelty over industry baseline: 20/40**

The proposal explicitly disclaims innovation:

> "No special innovation -- this is standard convention-over-configuration"

Honest. The differentiation from the industry baseline is the specific set of structural heuristics (Go `cmd/` exit conditions, Python `[project.scripts]` detection, Node `bin` field). These are useful but not novel — they are direct applications of ecosystem conventions.

**Cross-domain inspiration: 8/35**

No cross-domain inspiration demonstrated. The proposal stays entirely within Go/Node/Python ecosystem conventions. No borrowing from other domains (e.g., IDE project-type detection, build system conventions, container orchestration health checks as an analogy for detection confidence).

**Simplicity of insight: 20/25**

> "project structure is a reliable secondary signal: a Go project with `cmd/` subdirectories (but without `internal/` and `pkg/` coexisting -- a pattern typical of microservices, not CLIs) is almost certainly a CLI"

The exit condition (`internal/` and `pkg/` coexisting) is a genuinely insightful refinement. It shows understanding of Go project layout conventions beyond surface-level pattern matching.

> "a Python project with `[project.scripts]` in `pyproject.toml` is a CLI tool"

This is less insightful — `[project.scripts]` is literally the definition of a CLI entry point. Not a creative leap.

### 6. Feasibility: 75/100

**Technical feasibility: 30/40**

All inference is filesystem-based. No new dependencies. The `Sources map[string]string` field addition is backward-compatible at the Go level (zero value is nil). However:
- The TUI confirmation flow changes require modifying `askSurfaceConfirmation` — the proposal mentions this function by name, suggesting familiarity with the codebase, but does not describe the modification complexity.
- The `forge surfaces detect` command requires wiring through the CLI command structure, which is not discussed.

**Resource & timeline feasibility: 20/30**

> "3-4 coding tasks: (1) inference rules per ecosystem, (2) Sources map integration into DetectResult, (3) TUI confirmation flow update, (4) init summary display + tests. Additionally: (5) `forge surfaces detect` command with `--apply` flag."

The proposal now lists 5 tasks, up from the "2-3 coding tasks" referenced in scope. This is more honest. However:
- No time estimate per task
- No dependency between tasks stated (can they be parallelized?)
- Task #5 "reuses the detection + TUI infrastructure from tasks 1-4" — this is a dependency, implying sequential execution

**Dependency readiness: 25/30**

No external dependencies. Ecosystem convention stability is an unstated dependency — if a Go community convention shifts away from `cmd/` directories, the inference rule degrades silently. The proposal acknowledges this risk indirectly:

> "Inference rules become stale as ecosystems evolve | L | L"

### 7. Scope Definition: 72/80

**In-scope items are concrete: 26/30**

Four in-scope items. Most are specific deliverables. Deduction:

> "Per-path inference source annotation in `DetectResult` via `Sources map[string]string` field"

This is an implementation detail, not a user-facing deliverable. A better phrasing: "Users can see the source of each surface detection result (inference vs dependency signal) in TUI and init output."

> "`forge surfaces detect` subcommand: runs detection + inference independently of `forge init`, default read-only with `--apply` flag for config writing, and non-interactive mode"

Clear, concrete, testable.

**Out-of-scope explicitly listed: 22/25**

Five out-of-scope items, each with justification:

> "Rust/Cargo.toml structural inference (Rust's `[[bin]]`/`[lib]` targets in Cargo.toml could distinguish CLI vs library, but Rust project surface detection is not yet implemented in the existing dependency signal maps; adding inference for an ecosystem with no baseline detection support would be premature)"

Excellent justification — explains why, not just what.

> "Changes to `forge surfaces` CLI command output format"

This is surprising — `forge surfaces detect` is a new subcommand under `forge surfaces`, but "changes to `forge surfaces` CLI command output format" is out of scope. This seems contradictory: adding a new subcommand is by definition a change to the `forge surfaces` command structure. The intent is likely "no changes to existing `forge surfaces` subcommands," but the wording is ambiguous.

**Scope is bounded: 24/25**

5 tasks listed. The scope is bounded. Minor concern: the Python inference note in scope — "Python coverage is limited to CLI detection... Python API-type stdlib projects without Flask/FastAPI remain undetectable by structural conventions alone" — this caveat suggests the scope may need expansion for full ecosystem coverage.

### 8. Risk Assessment: 58/90

**Risks identified: 20/30**

Five risks listed. This is better than the minimum 3. The risks are:

1. Structural inference gives wrong type
2. Inference rules become stale as ecosystems evolve
3. Init summary format change breaks scripted consumers
4. `Sources map[string]string` schema migration
5. Re-run prompt adds friction for power users

Missing risks:
- YAML comment stripping during config round-trip (the YAML comment persistence approach has no risk entry despite being fragile)
- False negative from overly conservative inference rules (e.g., Python `app.py` exclusion condition filters out legitimate CLI projects)
- Non-interactive terminal output format stability (no contract defined for stdout format in `forge surfaces detect`)

**Likelihood + impact rated: 18/30**

Ratings are provided for all 5 risks. Analysis:

> "Structural inference gives wrong type | M | M"

Reasonable. The revision's conservative Go rules and Python exclusion conditions reduce likelihood. However:

> "`Sources map[string]string` schema migration | L | M"

This is rated L likelihood. But the risk description says "add compile-time check that all code paths populate `Sources` when inference fires" — the existence of a mitigation implies the risk is non-trivial. The L rating may be optimistic.

> "Re-run prompt adds friction for power users | L | L"

Fair. First-run experience unchanged.

**Mitigations are actionable: 20/30**

> "Annotate source in output; add integration test suite covering each inference rule with known-correct fixtures; ship with verbose logging (`FORGE_DEBUG=1`)"

Actionable: integration test suite is concrete. `FORGE_DEBUG=1` is an existing mechanism.

> "Rules are simple filesystem checks, easy to update; document each rule in code comments with the convention it targets"

Partially actionable: "easy to update" is a property, not an action. "Document each rule in code comments" is actionable.

> "grep test in CI verifies no programmatic parser exists in codebase"

Actionable and specific.

> "Zero-value of `map[string]string` is `nil`, so existing `DetectResult` construction sites are backward-compatible without changes"

This is an analysis, not a mitigation. The actual mitigation is "add compile-time check that all code paths populate `Sources` when inference fires" — this is actionable.

### 9. Success Criteria: 58/80

**Criteria are measurable and testable: 40/55**

17 criteria listed. Most are specific and testable:

> "Go project with only `cmd/` subdirs -> detected as `cli` without manual entry"

Testable: set up fixture, run detection, assert result.

> "Priority chain verified: when dependency signals return non-empty, inference functions are never called"

Testable: mock dependency signals, verify inference functions not invoked.

> "Existing detection tests (27 tests) still pass unchanged"

Testable: run test suite.

> "YAML comment persistence: after `forge init` infers `cli` from `cmd/` directory, config file contains `surfaces: cli  # source: inference:cmd-dir`; the comment is present in the file but not parsed by the YAML unmarshaler"

Testable: read config file, assert comment present, assert YAML unmarshaler ignores it. However: "not parsed by the YAML unmarshaler" is an implementation assertion, not a user-facing criterion. And this criterion does not address YAML library comment-preservation behavior — if the library strips comments during write, the test will fail for infrastructure reasons, not logic reasons.

Gaps:
- No criterion for performance budget (<50ms NFR has no corresponding testable criterion)
- No criterion for filesystem error handling (NFR says "return empty on any filesystem error" but no criterion verifies this)
- No criterion for non-interactive terminal output format
- No criterion for accuracy (what percentage of inferences must be correct on a test corpus?)

**Coverage is complete: 18/25**

Criteria cover most in-scope items. Gaps:
- No criterion for the TUI hint text: "This was inferred from project structure. Edit to correct if needed." — described in Solution #3 but not in Success Criteria
- No criterion verifying the `--apply` flag behavior difference (read-only vs write)
- Performance budget and error handling NFRs have no corresponding criteria

### 10. Logical Consistency: 72/90

**Solution addresses the stated problem: 28/35**

The stated problem has two parts:
1. Stdlib detection gap: Solution #1 (structural inference) directly addresses this. Covered for Go and Node.js. Python coverage is partial (only CLI inference, no API inference for stdlib Python projects without Flask/FastAPI).
2. Init summary opacity: Solution #2 (summary format) directly addresses this.

Solutions #3 and #4 go beyond the stated problem. The `forge surfaces detect` command (Solution #4) is explicitly "independent of `forge init`," which contradicts the bundling rationale:

> "both affect the same `forge init` code path"

This is a logical inconsistency in scope justification. The solution is sound but the justification for including it in this proposal is weak.

**Scope <-> Solution <-> Success Criteria aligned: 22/30**

Most in-scope items have corresponding success criteria. Gaps:
- "Per-path inference source annotation in `DetectResult` via `Sources map[string]string`" (in-scope) has criterion #7 verifying Sources map content — aligned.
- "Structural inference fallback for Go, Node.js, Python" (in-scope) has criteria #1-#5 — aligned.
- Python scope caveat ("Python coverage is limited to CLI detection") is not reflected in any criterion as a limitation.
- The YAML comment persistence (in-scope per Constraints section) has criterion #12 — aligned.

**Requirements <-> Solution coherent: 22/25**

The Requirements Analysis scenarios map to the Solution components. Key Scenario #1 (Go stdlib) maps to `inferGoSurface`. Key Scenario #2 (Node.js) maps to `inferNodeSurface`. Key Scenario #3 (Python) maps to `inferPythonSurface`. The coherence is strong.

Minor gap: the "No manifest at all" scenario (#6) states "detection returns empty -> manual TUI entry (unchanged behavior)" — but the proposal does not state what the new init summary shows in this case. Does it still show "(1 mappings)" or something else? The success criteria do not address this scenario.

## Phase 3: Blindspot Hunt

[blindspot] **YAML comment preservation is assumed but not verified.** The proposal's approach to source persistence relies on YAML comments surviving round-trip through the config library:

> "when writing surfaces config, append a comment line like `surfaces: api  # source: inference:cmd-dir` for human context"

Many Go YAML libraries (including `gopkg.in/yaml.v3`) strip comments during marshal/unmarshal cycles. The proposal does not identify which YAML library forge uses, whether it preserves comments, or what the fallback is if comments are stripped. This is a technical assumption treated as a fact, with no risk entry and no mitigation.

[blindspot] **The Python `app.py`/`main.py` exclusion condition is under-specified.** Key Scenario #3 states:

> "`app.py`/`main.py` at root -> `cli` only when no `setup.py` with `name` matching the directory AND no `[project.packages]`/`[tool.setuptools.packages.find]` library-only markers"

This exclusion logic has two conditions joined by AND. The first condition (`setup.py` with `name` matching directory) means that a project named "myapp" in a directory called "myapp" with `setup.py` containing `name="myapp"` will NOT be inferred as CLI, even if it has `[project.scripts]`. This may be correct for library packages, but the interaction between the two exclusion conditions and the `[project.scripts]` positive signal is complex and error-prone. No success criterion tests the exclusion conditions.

[blindspot] **Non-interactive terminal output format is undefined.** Solution #4 states:

> "Non-interactive terminals print results to stdout and exit (no TUI, no config write regardless of flags)"

But what format? JSON? Plain text? What fields? This is a user-facing behavior that is completely undefined. If `forge surfaces detect` is intended for scripting use (which the `--apply` flag and non-interactive mode suggest), the output format is a contract, not an implementation detail.

[blindspot] **The "do nothing" alternative was not honestly evaluated.** The comparison table states:

> "Do nothing | Zero cost | User confusion persists, onboarding friction | Rejected: user-reported pain point"

The "zero cost" claim ignores the ongoing cost of user confusion. The "rejected" verdict is based on the problem's existence, not on whether doing nothing is viable. For a CTO evaluating resource allocation, "do nothing" is always a valid option — the question is whether the pain justifies the investment. The proposal does not quantify the investment required, making the rejection of "do nothing" unjustified.

[blindspot] **Inference-only activation constraint may mask diagnostic gaps.** The priority chain states:

> "Inference signals only activate when ALL dependency signals are empty -- they never mix with dependency results"

This means if a project has a dependency signal for one surface and structural inference would detect another, the inference is suppressed entirely. For example: a Go project with cobra detected (CLI) but also an `api/` directory (API surface) — the API inference is never attempted because dependency signals returned non-empty for CLI. The user gets CLI detection but misses the API surface. This is a design trade-off that is not analyzed.

[conflict-with-pre-revision] **The re-run behavior (A3) creates a new user flow but no corresponding risk.** The revision expanded re-run to three options (Confirm / Re-detect / Edit), which is a new interaction pattern. The risk table has "Re-run prompt adds friction for power users | L | L" but this does not capture the risk of the new Edit flow diverging from the first-run manual entry flow (are they the same `manualSurfaceEntry` function? different UX?). The pre-revision two-option flow was simpler and had no such risk.

## Bias Detection Report

**Annotated regions**: 3 attack points / 3 paragraphs = density 1.00
- A1: `forge surfaces detect` output format undefined (blindspot)
- A2: `index.html` subdirectory scanning clarification is reasonable, no new issue introduced
- A3: Re-run Edit flow creates new interaction pattern without corresponding risk (conflict-with-pre-revision)

**Unannotated regions**: 8 attack points / ~25 paragraphs = density 0.32
- YAML comment preservation assumption (blindspot)
- Python exclusion condition complexity (blindspot)
- "Do nothing" alternative not honestly evaluated (blindspot)
- Inference-only activation masking diagnostic gaps (blindspot)
- Evidence gap for Node.js/Python ecosystems
- Circular urgency argument
- Technical direction lacks architecture description
- Python scope caveat not reflected in success criteria

**Ratio (annotated/unannotated)**: 3.12

The annotated regions receive significantly higher attack density, driven by the expectation that revisions should not introduce new issues. A1 introduced a new command without defining its output format. A3 expanded the interaction model without adding corresponding risk analysis.

## Score Summary

```
SCORE: 699/1000
DIMENSIONS:
  Problem Definition: 76/110
  Solution Clarity: 90/120
  Industry Benchmarking: 62/120
  Requirements Completeness: 82/110
  Solution Creativity: 48/100
  Feasibility: 75/100
  Scope Definition: 72/80
  Risk Assessment: 58/90
  Success Criteria: 58/80
  Logical Consistency: 72/90
ATTACKS:
1. [Problem Definition]: Evidence gap for Node.js/Python -- "evidence is inferred from the Go testing rather than separately reproduced per ecosystem" is an explicit admission that 2 of 3 covered ecosystems lack evidence. The solution commits to inference for all three but evidence supports only one. Must either produce ecosystem-specific evidence or narrow the solution scope to match evidence.
2. [Problem Definition]: Circular urgency -- "Each failed detection adds friction that risks user abandonment at the onboarding stage" states hypothetical harm without demonstrating actual abandonment. No churn data, no support tickets. Must quantify the cost of delay or demonstrate actual user loss.
3. [Problem Definition]: Bundling justification contradicts scope -- "both affect the same `forge init` code path" justifies bundling two problems, but Solution #4 adds `forge surfaces detect` which is "independent of `forge init`". Must either remove Solution #4 from this proposal or revise the bundling rationale.
4. [Solution Clarity]: Architecture gap -- "Inference functions live in `detect_surface.go` alongside existing signal maps" states location but not architecture. Must describe the dispatch mechanism (how ecosystem type is determined), the call chain from CLI entry point to inference, and how Sources map is populated/merged.
5. [Solution Clarity]: Non-interactive output format undefined -- "Non-interactive terminals print results to stdout and exit" describes the behavior but not the output format. Must define the stdout contract (format, fields, exit codes) for `forge surfaces detect` in non-interactive mode.
6. [Industry Benchmarking]: Superficial benchmark analysis -- "VS Code's workspace detection reads `.vscode/settings.json` then falls back to probing directory conventions" is one sentence per tool with no analysis of how each tool's layered detection works, what trade-offs they made, or what lessons apply. Must provide depth on at least 2 benchmarks with analysis of their detection strategy, failure modes, and applicability.
7. [Industry Benchmarking]: ROI claim without analysis -- "**Selected: best ROI** -- 3-4 inference functions at ~20 LOC each" provides LOC estimates but no ROI analysis. ROI requires quantifying both return (how many users benefit, what is the value) and investment (not just LOC but maintenance, testing, false-positive cost). Must provide a structured ROI comparison across alternatives.
8. [Requirements Completeness]: Python exclusion condition untested -- "`app.py`/`main.py` at root -> `cli` only when no `setup.py` with `name` matching the directory AND no `[project.packages]`/`[tool.setuptools.packages.find]`" has complex exclusion logic but no success criterion tests the exclusion conditions. Must add criteria testing: (a) project with `app.py` + `setup.py` with matching `name` -> NOT cli; (b) project with `app.py` + `[tool.setuptools.packages.find]` -> NOT cli; (c) project with `app.py` alone -> cli.
9. [Solution Creativity]: Explicit innovation disclaimer -- "No special innovation -- this is standard convention-over-configuration" limits this dimension. The one insight (Go `cmd/` exit conditions) is a refinement. Must either identify genuine creative differentiation or accept that this dimension will remain low.
10. [Feasibility]: Task dependencies unstated -- "task 5 reuses the detection + TUI infrastructure from tasks 1-4" implies sequential dependency but no task ordering or dependency graph is provided. Must state task dependencies and critical path.
11. [Risk Assessment]: YAML comment preservation risk missing -- the source persistence approach relies on YAML comments surviving round-trip, but many Go YAML libraries strip comments. No risk entry exists for this. Must add a risk entry for YAML comment stripping and verify which library forge uses.
12. [Risk Assessment]: Inference suppression masks diagnostic gaps -- "Inference signals only activate when ALL dependency signals are empty" means partial dependency matches suppress all inference, potentially hiding valid surfaces. Must analyze this trade-off in the risk table.
13. [Risk Assessment]: [conflict-with-pre-revision] Re-run Edit flow risk missing -- A3 expanded re-run to three options including Edit, which creates a new interaction pattern. The existing risk "Re-run prompt adds friction for power users" does not capture the risk of the Edit flow diverging from first-run manual entry. Must add a risk entry for the new Edit flow.
14. [Success Criteria]: Performance budget untested -- "all inference functions combined must complete in <50ms" is an NFR but no success criterion verifies it. Must add a benchmark criterion.
15. [Success Criteria]: Filesystem error handling untested -- "Inference functions return empty on any filesystem error" is an NFR but no success criterion verifies the error-handling behavior. Must add criteria for permission denied and unreadable directory cases.
16. [Logical Consistency]: Scope justification contradiction -- bundling rationale ("same `forge init` code path") does not cover Solution #4 (`forge surfaces detect`, independent of init). Must align scope justification with actual scope.
17. [Logical Consistency]: Python scope caveat absent from criteria -- "Python coverage is limited to CLI detection... Python API-type stdlib projects without Flask/FastAPI remain undetectable" is a scope limitation with no corresponding criterion that acknowledges or tests this limitation boundary.
```

## Verdict

**699/1000 -- REJECT. Below 900 target.**

The proposal has solid fundamentals: the problem is real, the solution is practical, and the scope is bounded. The 17 success criteria provide substantial verification coverage. However, three categories of issues prevent approval:

1. **Evidence and urgency gaps** (Problem Definition: 76/110): The solution covers three ecosystems but evidence covers one. Urgency rests on hypothetical harm, not demonstrated cost.

2. **Unverified technical assumptions** (multiple dimensions): YAML comment preservation, inference suppression trade-offs, Python exclusion condition complexity, and non-interactive output format are all assumptions or undefined behaviors treated as resolved.

3. **Scope-solution justification mismatch** (Logical Consistency: 72/90): The bundling rationale contradicts the inclusion of `forge surfaces detect`, which is independent of `forge init`.

**Priority revisions for iteration 2:**

1. **Define non-interactive output format** -- the stdout contract for `forge surfaces detect` must be specified (format, fields, exit codes)
2. **Verify YAML comment preservation** -- identify which YAML library forge uses, test comment preservation, and add a risk entry
3. **Add success criteria for NFRs** -- performance budget (<50ms) and filesystem error handling need testable criteria
4. **Narrow or justify evidence scope** -- either produce Node.js/Python evidence or narrow the solution to Go-only with other ecosystems as follow-up
5. **Fix scope justification** -- either remove `forge surfaces detect` from this proposal or revise the bundling rationale
6. **Analyze inference suppression trade-off** -- document what surfaces are missed when dependency signals are partially non-empty
7. **Add Python exclusion condition test cases** -- the complex AND-based exclusion logic needs explicit success criteria testing each condition
