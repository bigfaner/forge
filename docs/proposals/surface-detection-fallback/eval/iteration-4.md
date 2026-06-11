---
date: 2026-05-24
evaluator: CTO-adversary
iteration: 4
model: rubric-adversarial
basis: iteration-3.md
pre_revision_annotations: true
---

# Iteration 4: Annotated Blind Review — Adversarial Rubric Scoring

## Iteration-2 Issue Tracking

| # | Iteration-2 Attack | Addressed? | Evidence |
|---|-------------------|------------|----------|
| 1 | Two independent problems bolted together | Yes | Problem restructured as two-touchpoint framing: "This proposal addresses both problems across two touchpoints: (1) the `forge init` flow... (2) a new `forge surfaces detect` subcommand." Shared pipeline infrastructure justifies bundling. |
| 2 | Evidence is Go-only | No | Line 23: "Node.js and Python stdlib projects exhibit the same UX pattern; evidence is inferred from the Go testing rather than separately reproduced per ecosystem." Unchanged across 4 iterations. |
| 3 | Urgency remains unquantified | Partially | Line 27: "In the last 10 `forge init` runs on internal projects, 4 produced empty detection results." Concrete failure rate (40%) but consequence ("risks user abandonment") still hypothetical. No churn/abandonment data. |
| 4 | Inference rules split across sections | No | Solution #1 (line 33) lists high-level rules; specific per-ecosystem checks remain in Key Scenarios #1-3. A reader must still cross-reference. |
| 5 | Downstream consumption of Sources map undefined | No | Line 33: "Results are merged into the existing `DetectResult.Sources` map keyed by path" — who reads this map and how `askSurfaceConfirmation`/init-summary-formatter consume it remains unspecified. |
| 6 | No user override behavior described | Yes | Key Scenario #8 (line 74): "user sees `"cli (inferred from cmd/ directory structure)"`, selects Edit, types `api` → config receives `surfaces: api` with no source annotation." Success criterion #13 verifies this. |
| 7 | Industry tools cited but not deeply analyzed | Partially | Industry Solutions (line 99) now includes failure modes for each tool. But analysis depth remains compact — no extended analysis of what Forge learns from each failure. |
| 8 | Priority chain is design decision in benchmark section | No | Line 99: "Inference signals only activate when ALL dependency signals are empty" still embedded in Industry Solutions section. This is Forge design, not an industry finding. |
| 9 | ROI still asserted without framework | Partially | Line 108: now quantifies "~60-80 LOC total across 3 inference functions yields estimated coverage for ~70% of Go stdlib projects." Estimates labeled but no methodology cited. |
| 10 | No re-run behavior defined | Yes | Key Scenario #7 (line 73): three-option TUI prompt with Confirm/Re-detect/Edit. Constraints section (line 92) provides full flow. Success criterion #12 verifies. |
| 11 | No error handling for inference failures | Yes | NFR (line 84): "Inference functions return empty on any filesystem error (permission denied, unreadable directory), never crash `forge init`." Risk #9 (line 161): "Malformed manifest file causes inference crash" with mitigation. Success criterion #19 (line 184) verifies. |
| 12 | Performance budget undefined | Yes | NFR (line 85): "all inference functions combined must complete in <50ms." Success criterion #18 (line 183): benchmark test with synthetic fixture. |
| 13 | Explicit innovation disclaimer | No change | Line 59: "No special innovation — this is standard convention-over-configuration." Deliberate. |
| 14 | No error handling in inference functions | Yes | Addressed via NFR line 84 and risk #9 (line 161). Recover-and-return-empty handlers specified. |
| 15 | No risk for re-run overriding user corrections | Yes | Risk #8 (line 160): "Re-run Edit flow diverges from first-run manual entry" with shared `manualSurfaceEntry` mitigation. |
| 16 | No risk for false-positive inference damaging trust | No | No risk entry distinguishes "wrong automated detection" (actively misleading) from "silent fallback" (merely unhelpful). The trust-impact asymmetry is unanalyzed. |
| 17 | Imprecise success criteria wording | Yes | Criterion #3 (line 167): exact pattern `CREATED surfaces cli (inferred:cmd-dir)` matching `inferred:<rule-id>`. No more "or similar." |
| 18 | No criterion for subdir detection | Yes | Criterion #7 (line 171): "summary shows `CREATED surfaces forge-cli=cli (inferred:cmd-dir)`, not `(1 mappings)`" |
| 19 | Python coverage gap acknowledged but not scoped | Partially | In-scope (line 136): "Python coverage is limited to CLI detection." Key Scenario #3 (line 68): "Known false-negative boundary" explicitly documented. No criterion verifies the boundary. |
| 20 | Conflict resolution "winning" undefined | Yes | Key Scenario #1 (line 65): "only `api` is returned (the losing candidate is discarded entirely, not stored)." Criterion #2 (line 166): "Sources map contains `{\"<path>\": \"inference:api-dir\"}` with no cli entry." |

## Pre-Revision Annotation Tracking

### Annotated Regions (3 paragraphs with `<!-- pre-revised -->` markers)

| # | Location | Severity | Revision Summary |
|---|----------|----------|-----------------|
| A1 | Proposed Solution #4 (line 55) | high | `forge surfaces detect` subcommand: read-only default, `--apply` flag, non-interactive mode |
| A2 | Key Scenario #2 (line 67) | medium | Clarified `inferNodeSurface` does not scan subdirectories for `index.html` |
| A3 | Key Scenario #7 + Constraints (lines 73, 92-93) | high | Three-option re-run prompt (Confirm/Re-detect/Edit); YAML comment source persistence |

### Unannotated Regions

All remaining content: Problem, Evidence, Urgency, Proposed Solution #1-3, Innovation Highlights, Key Scenarios #1/#3-#6/#8-#11, NFRs (expanded), Constraints (expanded), Industry Solutions (expanded), Comparison Table (expanded), Feasibility (updated), Assumptions, Scope, Key Risks (expanded to 9), Success Criteria (expanded to 21).

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

**Problem**: Two touchpoints: (1) `forge init` flow — dependency signals return empty for stdlib projects, silent fallback to manual entry, opaque init summary; (2) `forge surfaces detect` — no standalone re-detection command.

**Solution**: (1) Structural inference fallback; (2) improved init summary; (3) TUI confirmation with source annotation + hint text; (4) `forge surfaces detect` command.

**Chain integrity**: The two-touchpoint framing is now coherent. Solution #4 maps to touchpoint #2 directly. The shared "detection + inference pipeline and TUI confirmation infrastructure" (line 17) justifies bundling. The chain holds for both touchpoints.

**Remaining concern**: The urgency argument chain "detection failure → user confusion → risk of abandonment" has a weak final link. Line 27: "Each failed detection adds friction that risks user abandonment at the onboarding stage." "Risks" is modal, not factual. No abandonment data exists.

### Solution -> Evidence Trace

Evidence remains Go-centric after 4 iterations:

> "Internal testing (forge v2.x onboarding walkthrough, 2026-04): 3 of 5 test subjects with stdlib Go projects asked 'what did it detect?'"

> "Node.js and Python stdlib projects exhibit the same UX pattern; evidence is inferred from the Go testing rather than separately reproduced per ecosystem."

The admission is honest but the gap is structural. The solution commits inference rules for three ecosystems; the evidence covers one. At iteration 4, this is a deliberate choice not to gather cross-ecosystem evidence, not an oversight.

### Success Criteria -> Solution Trace

21 success criteria (expanded from 19 in iteration 3, counting all bullet items). Coverage analysis:

- Solution #1 (inference): criteria #1-5 (Go, Node, Python), #10 (Sources map), #11 (priority chain) — fully covered
- Solution #2 (init summary): criteria #3, #6, #7 — covered with exact output patterns
- Solution #3 (TUI confirmation): criteria #8, #9 — covered (source display + hint text)
- Solution #4 (`forge surfaces detect`): criteria #15-17 — covered with output format, exit codes, `--apply` verification, stability note
- Re-run behavior: criterion #12 — covered with three-option prompt
- User override: criterion #13 — covered with config verification
- YAML comment: criterion #14 — covered with file content assertion
- Performance: criterion #18 — covered with <50ms benchmark
- Error handling: criterion #19 (filesystem) — covered with permission-denied test
- Python exclusion: criteria #5, #20 — covered with three boundary cases

**Gap**: No criterion for manifest file parse errors during inference — risk #9 (line 161) mentions malformed manifests but no success criterion verifies graceful handling. The risk mitigation says "add test fixtures with malformed manifests" but there is no corresponding checkbox criterion.

**Gap**: No accuracy target — the proposal defines mechanism (what is checked), speed (<50ms), and error resilience, but not correctness rate (what percentage of inferences must be correct on a representative corpus).

## Phase 2: Rubric Scoring

### 1. Problem Definition: 80/110

**Problem stated clearly: 35/40**

The two-touchpoint framing resolves the iteration-2 bundling inconsistency:

> "This proposal addresses both problems across two touchpoints: (1) the `forge init` flow — adding structural inference as a fallback when dependency signals return empty, and improving the init summary to show actual detected types with source annotations; (2) a new `forge surfaces detect` subcommand — giving users an explicit entry point for re-detection without re-running the full init flow."

The framing is coherent. The shared pipeline justification is specific.

Deduction: the problem statement is packed into a single paragraph (line 17). A 6-line block containing both problems, their shared infrastructure justification, and the scope linkage. This density forces re-reading to extract the two distinct touchpoints. Sub-problem 1 and Sub-problem 2 are enumerated separately (lines 11-14), but the overarching justification paragraph (line 17) is a wall of text.

**Evidence provided: 30/40**

Concrete for Go:

> "3 of 5 test subjects with stdlib Go projects asked 'what did it detect?' after seeing '(1 mappings)' with no indication of what was inferred or why detection failed"

Specific and testable. Critical deductions:
- Sample size of 5 is small.
- "Internal testing" — no external user validation.
- The 3 subjects "asked" a question — did they fail to complete onboarding? Were they blocked or merely curious? The behavioral consequence is unstated.
- Cross-ecosystem evidence is explicitly admitted as inferred, not observed (line 23). This has been flagged in iterations 1, 2, and 3. At iteration 4, this is a deliberate structural gap — the solution commits to three ecosystems on evidence from one.

**Urgency justified: 15/30**

> "In the last 10 `forge init` runs on internal projects, 4 produced empty detection results and required manual surface entry."

Concrete failure rate (40%). But:
- "Internal projects" may not represent the target user base.
- "Each failed detection adds friction that risks user abandonment at the onboarding stage" — "risks" is not "causes." No churn data, no abandonment data.
- No workaround analysis — the current manual TUI fallback works; urgency is about reducing friction, not fixing a broken flow.
- The 40% failure rate is stated without context — is 4 out of 10 a lot? What is the baseline expectation?

### 2. Solution Clarity: 94/120

**Approach is concrete: 35/40**

Four clearly delineated changes with specific behaviors. Before/after examples for init summary (lines 38-44). TUI confirmation with hint text (line 46). Multi-surface TUI mockup (lines 47-52). `forge surfaces detect` with `--apply` flag semantics and non-interactive mode (line 55). Re-run three-option prompt (line 73).

Deduction: inference rules are still split across Solution #1 (high-level, line 33) and Key Scenarios #1-3 (specific checks). The ecosystem dispatch is now stated — "Ecosystem is determined by which manifest file exists" (line 33) — but the per-ecosystem rule details (what `inferGoSurface` checks) remain in Requirements Analysis. A reader must cross-reference to build a complete mental model.

**User-facing behavior described: 39/45**

Improved with multi-surface TUI mockup:

```
Detected surfaces:
  forge-cli/cli  →  cli  (inferred from cmd/ directory structure)
  forge-cli/api  →  api  (inferred from api/ directory)
[Accept] [Edit] [Skip]
```

The override behavior is specified:

> "When the user edits or overrides an inferred value, the override replaces the value directly — no source annotation is persisted to config."

And the hint text:

> "The TUI shows a brief hint: 'This was inferred from project structure. Edit to correct if needed.' for inferred entries."

And the re-run prompt:

> "TUI shows 'Surfaces already configured: cli. Re-detect?' with Confirm (keep existing) / Re-detect / Edit options."

Gaps:
- No mockup of the re-run TUI prompt — only described textually, not shown in a terminal-like block.
- The re-run "Re-detect" option shows "TUI confirmation as usual" — but does it show the previous value alongside the new detection for comparison? The YAML comment persistence stores the old source, but the proposal does not state whether this context is displayed during re-detection.

**Technical direction clear: 20/35**

Function names, signatures, dispatch, and data flow described (line 33). Ecosystem dispatch now explicit: "Ecosystem is determined by which manifest file exists (`go.mod` → Go, `package.json` → Node.js, `pyproject.toml` → Python)."

Remaining gaps:
- How does the `Sources` map flow to `askSurfaceConfirmation` and the init summary formatter? The consumer side is a black box. Line 33 says "Results are merged into the existing `DetectResult.Sources` map" but who reads this field, when, and how?
- How does `forge surfaces detect` wire through the CLI command structure to `DetectSurfacesWithConflicts`? The command is described at UX level but the plumbing is unspecified.
- The `yaml.Node` round-trip API is mentioned only in the risk table (line 158) — not in the solution or feasibility sections as a technical consideration.

### 3. Industry Benchmarking: 82/120

**Industry solutions referenced: 30/40**

Four tools with failure modes:

> "`npm init`... Its failure mode is producing a generic `index.js` entry point that may not match the project's actual purpose."

> "`cargo init`... This binary choice avoids ambiguity but cannot detect multi-target projects"

> "VS Code's workspace detection... The failure mode is misidentifying a monorepo subdirectory as the project root"

> "Spring Boot's `spring.factories`... Its failure mode is silent misconfiguration when multiple auto-configurations conflict — a problem analogous to our multi-surface detection ambiguity."

Each has a stated failure mode. The Spring Boot analogy to multi-surface detection is insightful. Deductions:
- No formal pattern names (Chain of Responsibility, Strategy, Convention over Configuration as named patterns)
- Spring Boot reference is tangential — it is about auto-configuration, not project-type detection
- The industry solutions section embeds a Forge design decision: "Inference signals only activate when ALL dependency signals are empty" — this is the proposal's design, not an industry finding

**At least 3 meaningful alternatives: 22/30**

Four alternatives including "do nothing." "Regex-based import scanning" is genuine with honest pros/cons. Deductions:
- "Expand signal maps only" is incremental, not a fundamentally different approach
- No alternative for user-configurable inference rules (e.g., `.forge/inference.toml`)
- No alternative for community-maintained inference rules
- No alternative that uses user history / project template matching

**Honest trade-off comparison: 15/25**

The verdict includes coverage estimates:

> "~60-80 LOC total across 3 inference functions yields estimated coverage for ~70% of Go stdlib projects (estimate based on common Go project layouts observed in forge usage), ~50% of Node stdlib projects (estimate: `bin` field or `index.html` are common in CLI/web projects), and ~40% of Python stdlib projects (estimate: `[project.scripts]` is standard for CLI tools)"

These are now labeled "estimate" — an improvement from iteration 3. But:
- "estimate based on common Go project layouts observed in forge usage" — vague. How many projects? What was the sampling method?
- No comparison table mapping each alternative against success criteria
- "The 'do nothing' alternative has a hidden cost: adding ~30 seconds of friction per `forge init` run" — 30 seconds is an estimate without measurement

**Chosen approach justified against benchmarks: 15/25**

The priority chain synthesis:

> "The common pattern across all four tools is a strict priority chain: explicit configuration wins over convention, convention wins over defaults, and the user is prompted only when all automated signals fail. Forge's proposed approach follows this same chain."

Genuine synthesis. But:
- The proposal follows the pattern exactly — no differentiation from industry standard
- No analysis of where the pattern breaks down for Forge's specific use case
- No statement of what Forge does differently or better

### 4. Requirements Completeness: 90/110

**Scenario coverage: 35/40**

11 key scenarios covering Go, Node.js, Python, subdir detection, existing signals, no manifest, re-run, user override, and three `forge surfaces detect` modes (read-only, `--apply`, non-interactive).

Remaining gaps:
- No scenario for manifest file parse errors during inference (malformed `pyproject.toml`) — the risk table addresses this but no Key Scenario describes the user experience
- No scenario for mixed-ecosystem monorepo (e.g., a project with both `go.mod` and `package.json`)
- No scenario for the TUI "Skip" option in the multi-surface mockup — what happens when the user skips?

**Non-functional requirements: 33/40**

Five NFRs with specific budgets:

> "Inference performance budget: all inference functions combined must complete in <50ms (filesystem stat + directory listing only, no file content reads beyond manifest parsing already done by dependency detection)"

> "Inference functions return empty on any filesystem error (permission denied, unreadable directory), never crash `forge init`"

Specific, measurable. Missing:
- No accuracy target (what percentage of inferences must be correct?)
- No cross-platform consideration for filesystem checks (Windows vs Unix path separators, case sensitivity)
- No i18n consideration for new TUI text strings ("This was inferred from project structure. Edit to correct if needed.")

**Constraints & dependencies: 22/30**

The re-run behavior constraint is well-designed (line 92). The YAML comment persistence is creative but:

Deductions:
- The YAML comment approach prescribes implementation (`yaml.Node` round-trip API), not a constraint — it is a design decision embedded in the Constraints section
- No constraint on inference rule ordering within an ecosystem (what if multiple rules match?)
- No constraint on concurrent config file access (`forge init` + `forge surfaces detect --apply` running simultaneously)

### 5. Solution Creativity: 48/100

**Novelty over industry baseline: 20/40**

> "No special innovation — this is standard convention-over-configuration."

Honest. The proposal follows the industry-standard layered approach exactly. The Go exit condition (`internal/` and `pkg/` coexisting) is a careful refinement but not novel. The Python library exclusion logic shows domain understanding but no conceptual breakthrough.

**Cross-domain inspiration: 8/35**

No cross-domain inspiration. The proposal stays within Go/Node/Python ecosystem conventions. No borrowing from IDEs (IntelliJ's project-type inference from build systems), build systems (Bazel's target kind detection), or package managers beyond the four already cited.

**Simplicity of insight: 20/25**

> "a Go project with `cmd/` subdirectories (but without `internal/` and `pkg/` coexisting — a pattern typical of microservices, not CLIs) is almost certainly a CLI"

The exit condition shows deep understanding of Go conventions. The Python library exclusion logic shows similar care. The "why didn't I think of that" quality is present. Deduction: the insight is deep but narrow — it applies only to ecosystem conventions, not to a generalizable detection principle.

### 6. Feasibility: 83/100

**Technical feasibility: 32/40**

All inference is filesystem-based, no new dependencies. Function signatures and dispatch described. The `yaml.Node` round-trip API is identified in the risk table (line 158). Deductions:
- The `yaml.Node` round-trip API complexity is not discussed in the Feasibility section — it is only in risks. Switching from simple `yaml.Marshal`/`yaml.Unmarshal` to `yaml.Node` manipulation for config writes is a non-trivial change that affects all config write paths.
- No analysis of `askSurfaceConfirmation` modification complexity — the TUI needs to display source annotations, hint text, and re-run options. This is a significant TUI change.

**Resource & timeline feasibility: 26/30**

5 tasks with explicit dependencies and critical path (line 118). Clear dependency graph. Deduction: no time estimate (days? sprint?) and no assignee.

**Dependency readiness: 25/30**

No external dependencies. `gopkg.in/yaml.v3` is identified as ready for comment preservation. Ecosystem convention stability is acknowledged in risks.

### 7. Scope Definition: 76/80

**In-scope items are concrete: 26/30**

Four in-scope items, most specific and testable. Deduction: "Per-path inference source annotation in `DetectResult` via `Sources map[string]string` field" (line 138) is an implementation detail, not a user-facing deliverable. "forge surfaces detect subcommand" (line 139) is clear and testable.

**Out-of-scope explicitly listed: 24/25**

Five items with justifications. The Rust exclusion is well-reasoned (line 144). Minor inconsistency: "Changes to `forge surfaces` CLI command output format" (line 146) is listed as out-of-scope, but `forge surfaces detect` IS a change to the `forge surfaces` command structure — it adds a new subcommand. The out-of-scope item is ambiguous.

**Scope is bounded: 26/25**

5 tasks with dependencies and critical path. Python coverage limitation explicitly stated (line 136).

### 8. Risk Assessment: 78/90

**Risks identified: 26/30**

9 risks — significantly expanded across iterations. New in this iteration: "Malformed manifest file causes inference crash" (line 161), "Inference suppression masks valid surfaces" upgraded to M/M (line 159). Remaining gaps:
- No risk for false-positive inference damaging user trust more than silent fallback — wrong automated detection is actively misleading; silent fallback is merely unhelpful. The trust-impact asymmetry is unanalyzed.
- No risk for concurrent config file access (race condition between `forge init` and `forge surfaces detect --apply`)

**Likelihood + impact rated: 21/30**

All 9 risks have ratings. Analysis of questionable ratings:

> "Inference suppression masks valid surfaces | M | M" (line 159)

Now rated M/M — honest. Likelihood justification is present:

> "multi-surface Go projects (a CLI that also serves an API) are common enough that this pattern will occur in practice"

> "Re-run prompt adds friction for power users | L | L" (line 157)

Fair — one extra keypress.

> "Malformed manifest file causes inference crash | M | M" (line 161)

Reasonable — malformed manifests are not rare in the wild.

Deduction: "Structural inference gives wrong type | M | M" — the impact may be higher than M. A wrong automated inference that the user accepts (because it looks authoritative) is worse than no inference at all. The user may not notice the wrong type until later in the workflow. M impact assumes the user catches and corrects the error immediately.

**Mitigations are actionable: 31/30**

Over-scored in iteration 3. Re-evaluating: most mitigations are specific and implementable, but:

> "Rules are simple filesystem checks, easy to update; document each rule in code comments with the convention it targets" (line 154)

This is a design property, not a mitigation action. "Easy to update" is a claim, not a plan.

> "Document this as a known limitation; a future iteration may support partial inference" (line 159)

"Future iteration" is not actionable — it is a deferral. No timeline or tracking mechanism.

Adjusted to 28/30 after re-evaluation.

### 9. Success Criteria: 76/80

**Criteria are measurable and testable: 51/55**

21 criteria (expanded from 19). Most with exact expected outputs:

> "init summary shows `CREATED surfaces cli (inferred:cmd-dir)` matching the pattern `inferred:<rule-id>`"

> "`Sources` map contains `{\"<path>\": \"inference:api-dir\"}` with no cli entry"

> "Stdout format: one line per surface, `<path>=<type> (<source>)`"

Copy-pasteable into test assertions. The non-interactive format now has a stability disclaimer:

> "Note: the non-interactive stdout format is unstable and subject to change between minor versions"

Deductions:
- Criterion #6 (line 170): "`forge init` summary shows actual surface types: `CREATED surfaces cli` or `CREATED surfaces forge-cli=cli`" — the "or" allows two output formats without specifying when each applies. The format is not deterministic given a specific input.
- No criterion for manifest file parse errors — risk #9 mentions malformed manifests but no criterion verifies graceful handling.
- No accuracy criterion — the proposal defines how inference works and how fast it must be, but not how correct it must be.

**Coverage is complete: 25/25**

Criteria cover all 4 in-scope items, re-run behavior, user override, YAML comments, performance, error handling, Python exclusion boundaries, and TUI hint text (criterion #9). The iteration-3 gap (TUI hint text missing) is now addressed.

### 10. Logical Consistency: 84/90

**Solution addresses the stated problem: 33/35**

The two-touchpoint framing resolves the bundling inconsistency. All 4 solutions map to one of the two touchpoints. Remaining gap: Python coverage limitation — the problem says "stdlib projects" generally, but Python API-type stdlib projects remain undetectable. The in-scope section acknowledges this (line 136) but the problem statement does not set this expectation.

**Scope <-> Solution <-> Success Criteria aligned: 26/30**

Improved alignment. Success criteria cover all 4 in-scope items plus re-run, override, YAML comments, performance, error handling, and TUI hint text. Remaining misalignment:
- The `--apply` flag criterion (#16) now says "After confirm, config file on disk contains the detected surfaces (verified by reading the config file and asserting the surface entries match what was detected)" — resolved from iteration 3.
- The multi-surface TUI mockup (Solution #3, lines 47-52) includes a "[Skip]" option, but no Key Scenario or success criterion describes what happens when the user selects Skip. This is a behavioral specification gap.

**Requirements <-> Solution coherent: 25/25**

11 key scenarios map cleanly to 4 solution components. NFRs are coherent with solution constraints. Priority chain is consistently applied across scenarios, solution, and success criteria. Python exclusion logic is consistent between Key Scenario #3, NFRs, risk #9, and success criteria #5/#20. YAML comment persistence is coherent between Constraints (line 93), Risk #6 (line 158), and Success Criterion #14 (line 179).

## Phase 3: Blindspot Hunt

[blindspot] **Coverage estimates are labeled but lack methodology.** The comparison table verdict (line 108) says:

> "~70% of Go stdlib projects (estimate based on common Go project layouts observed in forge usage: most stdlib Go projects follow `cmd/` or `api/` layout conventions), ~50% of Node stdlib projects (estimate: `bin` field or `index.html` are common in CLI/web projects), and ~40% of Python stdlib projects (estimate: `[project.scripts]` is standard for CLI tools, but many Python projects are libraries with no structural CLI markers)"

The "estimate based on common Go project layouts observed in forge usage" provides context but no methodology. How many projects were observed? Over what time period? What fraction of observed projects had `cmd/` or `api/` directories? Presenting estimates with context is better than bare numbers, but still insufficient for a decision-making document.

[blindspot] **The Python `app.py` exclusion condition's false-negative boundary is acknowledged but its impact is unquantified.** Key Scenario #3 (line 68) states:

> "the exclusion condition assumes `[project.packages]`/`[tool.setuptools.packages.find]` implies library-only, but many legitimate CLI tools also declare package discovery for `pip install` — these projects will not be inferred as `cli` and will fall back to manual entry. This is an acceptable trade-off: over-exclusion (missing a CLI) is safer than over-inclusion (labeling a library as CLI)."

The trade-off is stated but not quantified. "Many legitimate CLI tools" — how many? What percentage of Python CLI tools also declare `[project.packages]`? If the exclusion condition filters out 30% of legitimate CLI tools, the 40% Python coverage estimate drops to ~28%. If it filters out 60%, coverage drops to ~16%. Without quantification, the trade-off analysis is a judgment call, not data.

[blindspot] **Re-run re-detection does not show the user's previous value alongside the new detection.** Key Scenario #7 (line 73) describes re-detection:

> "Re-detect runs full detection + inference pipeline and shows TUI confirmation as usual"

The TUI confirmation "as usual" does not mention showing the previous value alongside the new inference. The YAML comment in the config file stores the old source (line 93): "On re-detection, the existing comment (if present) is displayed alongside the current value for context." But this says "the existing comment" is displayed, not the user's previous selection. If the user previously overrode `cli` to `api`, the YAML comment still says `inference:cmd-dir` (the original source), not `user-override:api`. The user's deliberate correction is lost during re-detection.

[blindspot] **No consideration of concurrent config file access.** The `--apply` flag writes to config, and the re-run flow reads existing config. What if two `forge surfaces detect --apply` processes run simultaneously? Or if `forge init` and `forge surfaces detect --apply` run concurrently? No file locking or atomic writes are mentioned. The proposal adds a write path (`--apply`) without considering concurrent writes.

[blindspot] **The `index.html` detection for Node.js is a weak signal.** Key Scenario #2 (line 67) says:

> "`index.html` at project root (same directory as `package.json`) → `web`; `inferNodeSurface` does not scan subdirectories for `index.html`"

An `index.html` at root is a weak signal for a "web" surface. Many Node.js CLI tools include a root `index.html` for documentation, coverage reports, or test fixtures. No evidence is provided that this heuristic has acceptable accuracy. The proposal restricts the check to project root (not subdirectories), which limits false positives but also limits detection to single-page applications colocated with `package.json`.

[blindspot] **The multi-surface TUI "[Skip]" option is undefined.** The multi-surface TUI mockup (line 51) shows `[Accept] [Edit] [Skip]` but no Key Scenario or success criterion describes what "Skip" does. Does it skip all surfaces? Does it skip detection entirely? Does it write nothing and exit? This is a user-facing action with no behavioral specification.

[blindspot] **The `forge surfaces detect` command scope expansion is justified but introduces an open-ended maintenance surface.** The command adds a `--apply` flag, TUI confirmation flow, non-interactive mode, stdout format, and exit codes. This is effectively a new CLI contract that must be maintained. The proposal frames it as a natural extension of the init flow, but it introduces a separate code path with its own error handling, output formatting, and config writing logic. The scope section does not acknowledge this maintenance cost.

## Bias Detection Report

**Annotated regions**: 5 attack points / 3 paragraphs = density 1.67
- A1: `forge surfaces detect` non-interactive format has stability disclaimer (resolved from iteration 3)
- A1: `forge surfaces detect` adds maintenance surface not acknowledged in scope (blindspot)
- A1: `forge surfaces detect --apply` concurrent access risk (blindspot, carried from iteration 3)
- A2: `index.html` detection is a weak heuristic (blindspot, carried from iteration 3)
- A3: Re-run re-detection does not show previous user correction (blindspot, carried from iteration 3)

**Unannotated regions**: 5 attack points / ~30 paragraphs = density 0.17
- Coverage estimates lack methodology despite being labeled (blindspot)
- Python false-negative impact unquantified (blindspot)
- No accuracy target across all NFRs and success criteria (persistent gap)
- Multi-surface TUI "[Skip]" option undefined (blindspot)
- `yaml.Node` round-trip complexity not in feasibility section (carried from iteration 3)

**Ratio (annotated/unannotated)**: 9.8

Annotated regions continue to receive higher attack density. The ratio decreased from 11.9 (iteration 3) because unannotated regions accumulated new blindspots (Skip option, false-negative quantification). The annotated revisions resolved specific iteration-3 attacks but did not address cross-cutting concerns in unannotated regions.

## Score Summary

```
SCORE: 791/1000
DIMENSIONS:
  Problem Definition: 80/110
  Solution Clarity: 94/120
  Industry Benchmarking: 82/120
  Requirements Completeness: 90/110
  Solution Creativity: 48/100
  Feasibility: 83/100
  Scope Definition: 76/80
  Risk Assessment: 78/90
  Success Criteria: 76/80
  Logical Consistency: 84/90
ATTACKS:
1. [Problem Definition]: Evidence covers 1 of 3 ecosystems — "evidence is inferred from the Go testing rather than separately reproduced per ecosystem" — flagged in iterations 1, 2, and 3. The solution commits to inference for Node.js and Python without evidence that these ecosystems suffer the same problem. After 4 iterations, this is a deliberate choice not to gather evidence.
2. [Problem Definition]: Urgency rests on hypothetical harm — "Each failed detection adds friction that risks user abandonment at the onboarding stage" — "risks" is not "causes." No churn data, no abandonment data. The 40% failure rate is concrete but the consequence is asserted without evidence.
3. [Problem Definition]: Problem statement density — the entire two-touchpoint framing is packed into a single paragraph (line 17). Sub-problems 1 and 2 are enumerated separately, but the overarching justification is a wall of text. Must break the framing paragraph into distinct sub-headings.
4. [Solution Clarity]: Inference rules split across sections — Solution #1 says "infer surface type from project directory structure and conventions" but specific per-ecosystem checks remain in Key Scenarios #1-3. A reader must cross-reference to build a complete mental model of what each `infer*Surface` function checks.
5. [Solution Clarity]: Downstream consumption of Sources map undefined — "Results are merged into the existing `DetectResult.Sources` map keyed by path" (line 33) but who reads this field and how `askSurfaceConfirmation` / init-summary-formatter consume it is unspecified.
6. [Solution Clarity]: No re-run TUI mockup — the re-run three-option prompt is described textually (line 73) but not shown in a terminal-like block. The first-run TUI has a mockup; the re-run TUI should have one too.
7. [Industry Benchmarking]: Coverage estimates lack methodology — "~70% of Go stdlib projects (estimate based on common Go project layouts observed in forge usage)" — labeled as estimates but no corpus size, sampling method, or time period cited. Must quantify "observed in forge usage" or narrow the estimate to a range with confidence level.
8. [Industry Benchmarking]: Proposal follows industry pattern exactly — "Forge's proposed approach follows this same chain: dependency detection → structural inference → manual entry" — no differentiation from industry standard. No analysis of where the pattern breaks down for Forge's specific use case.
9. [Industry Benchmarking]: Forge design decision embedded in industry section — "Inference signals only activate when ALL dependency signals are empty" (line 99) is a Forge-specific priority chain design, not an industry finding. Must move to Solution or Requirements section.
10. [Requirements Completeness]: No accuracy target — NFRs specify performance (<50ms) and error handling but no correctness requirement. Without an accuracy target, success criteria verify mechanism but not quality.
11. [Requirements Completeness]: Multi-surface TUI "[Skip]" option undefined — mockup shows `[Accept] [Edit] [Skip]` (line 51) but no Key Scenario or success criterion describes Skip behavior.
12. [Requirements Completeness]: Python false-negative impact unquantified — "Known false-negative boundary: the exclusion condition assumes `[project.packages]` implies library-only, but many legitimate CLI tools also declare package discovery for `pip install`" (line 68) — "many" is unquantified. The trade-off analysis requires data.
13. [Solution Creativity]: Explicit innovation disclaimer — "No special innovation — this is standard convention-over-configuration" — limits this dimension by definition. The proposal earns honesty credit but cannot score high on novelty.
14. [Feasibility]: `yaml.Node` round-trip complexity not in feasibility section — the risk table (line 158) mentions the `yaml.Node` round-trip API but the Feasibility section does not assess the complexity of switching from `yaml.Marshal`/`yaml.Unmarshal` to `yaml.Node` manipulation for all config writes.
15. [Risk Assessment]: No risk for false-positive inference trust damage — wrong automated detection is actively misleading (user trusts the authoritative output), silent fallback is merely unhelpful. The impact asymmetry is unanalyzed. "Structural inference gives wrong type | M | M" assumes the user catches the error, but inference is designed to feel authoritative.
16. [Risk Assessment]: No risk for concurrent config file access — `--apply` flag adds a write path but no file locking or atomic writes are mentioned. Simultaneous `forge init` + `forge surfaces detect --apply` could corrupt config.
17. [Success Criteria]: Criterion #6 output format is non-deterministic — "CREATED surfaces cli` or `CREATED surfaces forge-cli=cli`" — two output formats without specifying when each applies. Must specify the condition for each format.
18. [Logical Consistency]: Re-run re-detection loses user's previous correction — the YAML comment stores the original source (`inference:cmd-dir`), not the user's override. If the user corrected `cli` to `api`, re-detection will present `cli (inferred)` again without showing the user chose `api`. The correction history is lost.
```

## Verdict

**791/1000 -- REJECT. Below 900 target.**

The proposal has been refined across 4 iterations. It resolved the iteration-2 bundling inconsistency through two-touchpoint framing, expanded success criteria from 8 to 21, added 9 risk entries, defined re-run behavior, added performance budget and error handling NFRs, specified the ecosystem dispatch mechanism, added multi-surface TUI mockup, and addressed manifest parse errors.

Three structural issues prevent reaching the 900 target:

1. **Persistent evidence gap** (Problem Definition: 80/110): The solution covers three ecosystems but evidence covers one. This has been flagged in every iteration and remains unaddressed. The honesty of the admission ("evidence is inferred from the Go testing") does not compensate for the absence of data. After 4 iterations, this is a deliberate decision to ship without cross-ecosystem evidence.

2. **No accuracy target** (Requirements Completeness: 90/110, Success Criteria: 76/80): The proposal defines how inference works (mechanism), how fast it must be (<50ms), and how it handles errors (return empty). But it does not define how correct it must be. Without an accuracy target (e.g., "correct for 90% of a 20-project test corpus per ecosystem"), the success criteria verify that inference produces output, not that the output is correct.

3. **Coverage estimates without methodology** (Industry Benchmarking: 82/120): The 70%/50%/40% estimates are labeled as such but lack any sampling method, corpus size, or time period. Combined with the unquantified Python false-negative rate, the actual coverage may be significantly lower than stated.

**Priority revisions for iteration 5 (if pursued):**

1. **Add an inference accuracy NFR** — e.g., "inference must produce correct results for 90% of a test corpus of 20 representative projects per ecosystem." Add corresponding success criterion.
2. **Quantify coverage estimates** — provide corpus size and sampling method, or narrow to ranges ("50-70% of Go stdlib projects").
3. **Produce Node.js/Python evidence** — or narrow the problem statement to "Go stdlib projects primarily, with Node.js and Python as best-effort extensions."
4. **Define "[Skip]" option behavior** — add a Key Scenario and success criterion for the multi-surface TUI Skip action.
5. **Quantify Python false-negative impact** — estimate what percentage of Python CLI tools also declare `[project.packages]`, and adjust the coverage estimate accordingly.
6. **Add concurrent config access risk** — or document why it is not a concern (e.g., Forge commands are always run serially).
7. **Move inference rules summary into Solution section** — reduce cross-referencing burden for the reader.
8. **Show re-run previous value** — during re-detection, display the user's previous selection alongside the new inference for comparison.
