---
date: 2026-05-24
evaluator: CTO-adversary
iteration: 3
model: rubric-adversarial
basis: iteration-2.md
pre_revision_annotations: true
---

# Iteration 3: Annotated Blind Review — Adversarial Rubric Scoring

## Iteration-1 Issue Tracking

| # | Iteration-1 Attack | Addressed? | Evidence |
|---|-------------------|------------|----------|
| 1 | Evidence gap for Node.js/Python | Partially | Proposal admits "evidence is inferred from the Go testing rather than separately reproduced per ecosystem" — acknowledged but not resolved. No new evidence added. |
| 2 | Circular urgency | Partially | Urgency now has concrete data: "In the last 10 `forge init` runs on internal projects, 4 produced empty detection results." But still no churn/abandonment data. |
| 3 | Bundling justification contradicts scope | Yes | Problem statement restructured as two-touchpoint framing: "This proposal improves surface detection across two touchpoints: (1) the `forge init` flow... (2) a new `forge surfaces detect` subcommand." No longer a bolted-on addition. |
| 4 | Architecture gap | Partially | Solution #1 describes dispatch: "calls ecosystem-specific inference functions only after all dependency signal maps return empty." But how ecosystem type is determined from manifest files is not stated. |
| 5 | Non-interactive output format undefined | Yes | Success criterion specifies: "Stdout format: one line per surface, `<path>=<type> (<source>)`, where `<source>` is `detected:<signal>` or `inferred:<rule-id>`" with exit codes 0/1. |
| 6 | Superficial benchmark analysis | Partially | Industry solutions now include failure modes for each tool (npm: generic entry point, cargo: binary choice, VS Code: monorepo misidentification, Spring Boot: conflicting auto-config). More depth, but still compact. |
| 7 | ROI claim without analysis | Yes | Comparison table verdict now quantifies: "~60-80 LOC total across 3 inference functions yields coverage for ~70% of Go stdlib projects, ~50% of Node stdlib, ~40% of Python stdlib" with explicit comparison to regex scanning LOC. |
| 8 | Python exclusion condition untested | Yes | Success criterion added: three specific test cases for Python exclusion boundary. |
| 9 | Explicit innovation disclaimer limits creativity | No change | Still disclaims: "No special innovation — this is standard convention-over-configuration." |
| 10 | Task dependencies unstated | Yes | Resource & Timeline now states: "Task dependencies: tasks 1-4 can be parallelized within the same PR; task 5 depends on tasks 1-3. Critical path: tasks 1-3 -> task 5, with task 4 in parallel." |
| 11 | YAML comment preservation risk missing | Yes | New risk: "YAML comment stripping during config round-trip | M | L" with specific mitigation referencing `gopkg.in/yaml.v3` `yaml.Node` round-trip API and integration test. |
| 12 | Inference suppression masks diagnostic gaps | Yes | New risk: "Inference suppression masks valid surfaces | L | M" with explicit trade-off analysis. |
| 13 | [conflict-with-pre-revision] Re-run Edit flow risk missing | Yes | New risk: "Re-run Edit flow diverges from first-run manual entry | L | L" with shared `manualSurfaceEntry` function mitigation. |
| 14 | Performance budget untested | Yes | New criterion: "Performance budget: all inference functions combined complete in <50ms on a project with 50+ directories (benchmark test with synthetic fixture)" |
| 15 | Filesystem error handling untested | Yes | New criterion: "Filesystem error resilience: inference function called on a directory with permission denied returns empty result without panicking" |
| 16 | Scope justification contradiction | Yes | Problem statement restructured to frame two touchpoints explicitly. |
| 17 | Python scope caveat absent from criteria | Yes | Python exclusion boundary criterion added with three specific test cases. |

## Pre-Revision Annotation Tracking

### Annotated Regions (3 paragraphs)

| # | Location | Severity | Revision Summary |
|---|----------|----------|-----------------|
| A1 | Proposed Solution #4 (line 42-43) | high | Added `forge surfaces detect` subcommand with `--apply` flag, non-interactive mode, explicit read-only default |
| A2 | Key Scenario #2 (line 54-55) | medium | Clarified `inferNodeSurface` does not scan subdirectories for `index.html` |
| A3 | Key Scenario #7 + Constraints (line 60-61, 80) | high | Expanded re-run behavior to three-option prompt (Confirm/Re-detect/Edit); added YAML comment source persistence |

### Unannotated Regions

All remaining content (~28 paragraphs): Problem, Evidence, Urgency, Proposed Solution #1-3, Innovation Highlights, Key Scenarios #1/#3-#6/#8-#11, NFRs (expanded), Constraints (expanded), Industry Solutions (expanded), Comparison Table (expanded), Feasibility (updated), Assumptions, Scope, Key Risks (expanded from 5 to 8), Success Criteria (expanded from 8 to 19), Next Steps.

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

**Problem**: The proposal frames the problem as two touchpoints sharing a pipeline: (1) `forge init` flow — dependency signals return empty for stdlib projects, falling back to manual entry silently; init summary shows opaque "N mappings"; (2) `forge surfaces detect` subcommand — users need an explicit entry point for re-detection.

**Solution**: (1) Structural inference fallback; (2) improved init summary; (3) TUI confirmation with source annotation; (4) `forge surfaces detect` command.

**Chain integrity**: The two-touchpoint framing resolves the iteration-2 bundling inconsistency. Solution #4 maps directly to touchpoint #2. The shared "detection + inference pipeline and TUI confirmation infrastructure" justifies bundling them in one proposal. The chain holds cleanly for both touchpoints.

**Remaining concern**: The urgency argument uses internal-run data (4 of 10 empty detections = 40% failure rate) but does not quantify the cost of those failures. The chain from "detection failure -> user confusion -> risk of abandonment" has a weak final link — "risk of abandonment" is hypothetical.

### Solution -> Evidence Trace

Evidence remains Go-centric:

> "Internal testing (forge v2.x onboarding walkthrough, 2026-04): 3 of 5 test subjects with stdlib Go projects asked 'what did it detect?'"

> "Node.js and Python stdlib projects exhibit the same UX pattern; evidence is inferred from the Go testing rather than separately reproduced per ecosystem."

Honest but inadequate. The solution covers three ecosystems; the evidence covers one. This has been flagged in iteration 1 and iteration 2 and remains unaddressed. At this point, the absence of cross-ecosystem evidence is a deliberate choice, not an oversight.

### Success Criteria -> Solution Trace

19 success criteria — a substantial expansion from iteration 2's 8. Coverage is now near-complete:

- Solution #1 (inference): criteria #1-#5 (Go, Node, Python), #9 (Sources map), #10 (priority chain) — fully covered
- Solution #2 (init summary): criteria #3, #6, #7 — covered with exact output patterns
- Solution #3 (TUI confirmation): criterion #8 — covered
- Solution #4 (`forge surfaces detect`): criteria #15-#17 — covered with output format and exit codes
- Re-run behavior: criterion #12 — covered with three-option prompt
- User override: criterion #13 — covered with config verification
- YAML comment: criterion #14 — covered with file content assertion
- Performance: criterion #18 — covered with <50ms benchmark
- Error handling: criterion #19 (filesystem) — covered with permission-denied test
- Python exclusion: criterion #5, #19 — covered with three boundary cases

**Remaining gap**: No criterion for manifest file parse errors during inference (malformed `pyproject.toml`). The filesystem error criterion covers permission denied but not parse failures.

## Phase 2: Rubric Scoring

### 1. Problem Definition: 82/110

**Problem stated clearly: 35/40**

The two-touchpoint framing resolves the bundling issue from iteration 2:

> "This proposal improves surface detection across two touchpoints: (1) the `forge init` flow — adding structural inference as a fallback when dependency signals return empty, and improving the init summary to show actual detected types with source annotations; (2) a new `forge surfaces detect` subcommand — giving users an explicit entry point for re-detection without re-running the full init flow."

This is coherent. The shared pipeline justification ("Both touchpoints share the same detection + inference pipeline and TUI confirmation infrastructure") is specific and defensible.

Deduction: the problem statement is dense — the entire framing is packed into a single paragraph. A reader must parse a 6-line block to extract the two distinct touchpoints. The density obscures the core insight.

**Evidence provided: 32/40**

Concrete evidence for Go:

> "3 of 5 test subjects with stdlib Go projects asked 'what did it detect?' after seeing '(1 mappings)' with no indication of what was inferred or why detection failed"

Specific and testable. But:
- Sample size of 5 is small.
- "Internal testing" — no external user validation.
- The 3 subjects "asked" a question — did they fail to complete onboarding? Were they blocked or merely curious?
- Cross-ecosystem evidence is explicitly admitted as inferred, not observed. Flagged in iterations 1 and 2; still unresolved.

**Urgency justified: 15/30**

> "In the last 10 `forge init` runs on internal projects, 4 produced empty detection results and required manual surface entry."

Concrete data point (40% failure rate). But:
- "Internal projects" may not represent the target user base.
- "Each failed detection adds friction that risks user abandonment" — "risks" is not "causes." No actual abandonment data.
- No workaround analysis — the current manual TUI fallback works; the urgency is about reducing friction, not fixing a broken flow.

### 2. Solution Clarity: 95/120

**Approach is concrete: 36/40**

Four clearly delineated changes with specific behaviors. Before/after examples for init summary. TUI confirmation with hint text. `forge surfaces detect` with `--apply` flag semantics and non-interactive mode. Re-run three-option prompt.

Deduction: the inference rules are still split across Solution #1 (high-level) and Requirements Analysis > Key Scenarios (specific checks). A reader must cross-reference to understand what `inferGoSurface` actually checks.

**User-facing behavior described: 39/45**

Improved significantly. The TUI confirmation now includes the override behavior:

> "When the user edits or overrides an inferred value, the override replaces the value directly — no source annotation is persisted to config."

And the hint text:

> "The TUI shows a brief hint: 'This was inferred from project structure. Edit to correct if needed.' for inferred entries."

And the re-run prompt:

> "TUI shows 'Surfaces already configured: cli. Re-detect?' with Confirm (keep existing) / Re-detect / Edit (enter `manualSurfaceEntry` flow to modify manually) options."

Gaps:
- No mockup of the re-run TUI prompt with three options — only described textually.
- No mockup of multi-surface TUI display (what does the user see when both `cli` and `api` are detected?).

**Technical direction clear: 20/35**

Function names, signatures, dispatch, and data flow described:

> "`detectSurfaceAtDirWithConflicts` calls ecosystem-specific inference functions (e.g., `inferGoSurface`, `inferNodeSurface`, `inferPythonSurface`) only after all dependency signal maps return empty. Each inference function returns `(surfaceType, sourceAnnotation)` and is stateless — pure filesystem reads with no side effects."

Remaining gaps:
- How is the ecosystem type determined from manifest files? The dispatch mechanism is implied but not stated.
- How does the `Sources` map flow to `askSurfaceConfirmation` and the init summary formatter? The consumer side is a black box.
- How does `forge surfaces detect` wire through the CLI command structure to `DetectSurfacesWithConflicts`?

### 3. Industry Benchmarking: 85/120

**Industry solutions referenced: 32/40**

Four tools with failure modes:

> "`npm init` infers project type from existing `package.json` fields... Its failure mode is producing a generic `index.js` entry point"

> "`cargo init`... This binary choice avoids ambiguity but cannot detect multi-target projects"

> "VS Code's workspace detection... The failure mode is misidentifying a monorepo subdirectory as the project root"

> "Spring Boot's `spring.factories`... Its failure mode is silent misconfiguration when multiple auto-configurations conflict — a problem analogous to our multi-surface detection ambiguity"

Each tool has a stated failure mode. The Spring Boot analogy to multi-surface detection is insightful. Deductions:
- No published pattern names (Chain of Responsibility, Convention over Configuration as formal patterns)
- The Spring Boot reference is tangential to project-type detection (it is about auto-configuration)
- No reference to language server protocols which also perform project-type detection

**At least 3 meaningful alternatives: 22/30**

Four alternatives including "do nothing." "Regex-based import scanning" is a genuine alternative with honest pros/cons. Deductions:
- "Expand signal maps only" is incremental, not a fundamentally different approach
- No alternative for user-configurable inference rules (e.g., `.forge/inference.toml`)
- No alternative for crowdsourced/community-maintained inference rules

**Honest trade-off comparison: 16/25**

The verdict now includes quantified coverage estimates:

> "~60-80 LOC total across 3 inference functions yields coverage for ~70% of Go stdlib projects (most have `cmd/` or `api/` layout), ~50% of Node stdlib projects (`bin` field or `index.html`), and ~40% of Python stdlib projects"

These percentages appear data-driven but have no cited source. Deductions:
- The coverage percentages are stated without methodology — are they from corpus analysis? Expert judgment? Gut feel?
- No comparison table mapping each alternative against the success criteria
- "The 'do nothing' alternative has a hidden cost: adding ~30 seconds of friction per `forge init` run" — 30 seconds is an estimate without measurement

**Chosen approach justified against benchmarks: 15/25**

The priority chain is mapped to industry:

> "The common pattern across all four tools is a strict priority chain: explicit configuration wins over convention, convention wins over defaults, and the user is prompted only when all automated signals fail."

Genuine synthesis. But:
- The proposal follows the pattern exactly — no differentiation from industry standard
- No analysis of where the pattern breaks down for Forge's specific use case

### 4. Requirements Completeness: 92/110

**Scenario coverage: 36/40**

11 key scenarios covering Go, Node.js, Python, subdir detection, existing signals, no manifest, re-run, user override, and three `forge surfaces detect` modes (read-only, `--apply`, non-interactive). The re-run scenario is now detailed with three options.

Remaining gaps:
- No scenario for manifest file parse errors (malformed `pyproject.toml`)
- No scenario for multi-surface detection across subdirectories
- No scenario for mixed-ecosystem monorepo

**Non-functional requirements: 34/40**

Five NFRs, improved with specific budgets:

> "Inference performance budget: all inference functions combined must complete in <50ms (filesystem stat + directory listing only, no file content reads beyond manifest parsing already done by dependency detection)"

Specific, measurable, with implementation constraints.

> "Inference functions return empty on any filesystem error (permission denied, unreadable directory), never crash `forge init`"

Concrete error handling constraint.

Missing:
- No accuracy target (what percentage of inferences must be correct on a test corpus?)
- No i18n consideration for new TUI text strings
- No cross-platform (Windows vs Unix) consideration for filesystem checks

**Constraints & dependencies: 22/30**

The re-run behavior constraint is well-designed:

> "`runSurfaceConfig` detects existing surfaces in config, shows TUI prompt with Confirm / Re-detect / Edit options. Confirm skips detection and returns `SKIPPED surfaces (already configured)`. Re-detect runs the full detection + inference pipeline. Edit enters `manualSurfaceEntry` flow. This eliminates the need for a `--force` flag."

This is a specific, actionable design. The YAML comment persistence constraint is creative but prescribes implementation rather than stating a constraint.

Deductions:
- The YAML comment approach is a design decision, not a constraint — it prescribes using YAML comments for source persistence
- No constraint on inference rule ordering (what if multiple rules match for the same ecosystem?)

### 5. Solution Creativity: 48/100

**Novelty over industry baseline: 20/40**

> "No special innovation — this is standard convention-over-configuration"

Honest. The proposal follows the industry-standard layered approach exactly. The Go exit condition (`internal/` and `pkg/` coexisting) is a careful refinement but not novel.

**Cross-domain inspiration: 8/35**

No cross-domain inspiration. The proposal stays within Go/Node/Python ecosystem conventions. No borrowing from IDEs, build systems, or package managers beyond the four already cited.

**Simplicity of insight: 20/25**

> "a Go project with `cmd/` subdirectories (but without `internal/` and `pkg/` coexisting — a pattern typical of microservices, not CLIs) is almost certainly a CLI"

The exit condition shows deep understanding of Go conventions. The Python library exclusion logic shows similar care. The "why didn't I think of that" quality is present.

### 6. Feasibility: 85/100

**Technical feasibility: 33/40**

All inference is filesystem-based, no new dependencies. Function signatures and dispatch described. `Sources map[string]string` backward compatibility analyzed. The YAML comment approach identifies the specific library mechanism:

> "Forge uses `gopkg.in/yaml.v3` which preserves comments when using `yaml.Node` round-trip API"

Deductions:
- The `yaml.Node` round-trip API is more complex than simple `yaml.Marshal`/`yaml.Unmarshal` — this complexity is not discussed in feasibility
- No analysis of `askSurfaceConfirmation` modification complexity

**Resource & timeline feasibility: 27/30**

5 tasks with explicit dependencies and critical path:

> "Task dependencies: tasks 1-4 can be parallelized within the same PR (shared `detect_surface.go` file but independent functions); task 5 depends on tasks 1-3 (reuses detection + inference pipeline and TUI infrastructure). Critical path: tasks 1-3 -> task 5, with task 4 in parallel."

Clear dependency graph. Deduction: no time estimate (days? sprint?) and no assignee.

**Dependency readiness: 25/30**

No external dependencies. `gopkg.in/yaml.v3` is identified as ready for comment preservation. Ecosystem convention stability is acknowledged in risks.

### 7. Scope Definition: 78/80

**In-scope items are concrete: 27/30**

Four in-scope items, most specific and testable. Deduction: "Per-path inference source annotation in `DetectResult` via `Sources map[string]string` field" is an implementation detail, not a user-facing deliverable.

**Out-of-scope explicitly listed: 24/25**

Five items with justifications. The Rust exclusion is well-reasoned. Minor gap: "Changes to `forge surfaces` CLI command output format" is ambiguous given that `forge surfaces detect` IS a change to the `forge surfaces` command structure.

**Scope is bounded: 27/25**

5 tasks with dependencies and critical path. Python coverage limitation explicitly stated.

### 8. Risk Assessment: 76/90

**Risks identified: 24/30**

8 risks — significantly expanded. New risks address iteration-2 blindspots: YAML comment stripping, inference suppression masking valid surfaces, re-run Edit flow divergence. Remaining gaps:
- No risk for manifest file parse errors (malformed `pyproject.toml`)
- No risk for inference rule ordering ambiguity
- No risk for false-negative from overly conservative Python exclusion (projects with both `[project.scripts]` and `[project.packages]`)

**Likelihood + impact rated: 22/30**

All 8 risks have ratings. Analysis of key ratings:

> "Inference suppression masks valid surfaces | L | M"

The L likelihood is questionable. Projects with both CLI and API components are common in Go (e.g., a cobra CLI that also serves an API). The proposal itself acknowledges the scenario. L may be optimistic.

> "YAML comment stripping during config round-trip | M | L"

Fair — M likelihood because YAML comment preservation is fragile; L impact because comments are informational only.

**Mitigations are actionable: 30/30**

All 8 mitigations are specific and implementable:
- "add integration test suite covering each inference rule with known-correct fixtures"
- "ship with verbose logging (`FORGE_DEBUG=1`)"
- "grep test in CI verifies no programmatic parser exists"
- "Forge uses `gopkg.in/yaml.v3` which preserves comments when using `yaml.Node` round-trip API... Verification: integration test reads config file after `forge init` and asserts comment is present"
- "Document this as a known limitation; a future iteration may support partial inference"
- "shared function ensures single code path; integration test verifies Edit produces same output as first-run manual entry"

Full marks — every mitigation includes a concrete verification mechanism.

### 9. Success Criteria: 76/80

**Criteria are measurable and testable: 52/55**

19 criteria — most with exact expected outputs:

> "init summary shows `CREATED surfaces cli (inferred:cmd-dir)` matching the pattern `inferred:<rule-id>`"

> "`Sources` map contains `{"<path>": "inference:api-dir"}` with no cli entry"

> "Stdout format: one line per surface, `<path>=<type> (<source>)`"

Copy-pasteable into test assertions. Deductions:
- Criterion #7: "`forge init` summary shows actual surface types: `CREATED surfaces cli` or `CREATED surfaces forge-cli=cli`" — the "or" allows two output formats without specifying when each applies. Is the format deterministic given the input?
- No criterion for manifest file parse errors during inference
- No criterion for inference accuracy (what percentage of inferences must be correct?)

**Coverage is complete: 24/25**

Criteria cover all 4 in-scope items, re-run behavior, user override, YAML comments, performance, error handling, and Python exclusion boundaries. Remaining gap: no criterion for the TUI hint text described in Solution #3: "This was inferred from project structure. Edit to correct if needed."

### 10. Logical Consistency: 84/90

**Solution addresses the stated problem: 33/35**

The two-touchpoint framing resolves the iteration-2 bundling inconsistency. All 4 solutions map to one of the two touchpoints. Remaining gap: Python coverage limitation — the problem says "stdlib projects" generally, but Python API-type stdlib projects remain undetectable.

**Scope <-> Solution <-> Success Criteria aligned: 26/30**

Improved alignment. Success criteria cover all 4 in-scope items plus re-run, override, YAML comments, performance, and error handling. Remaining misalignment:
- The TUI hint text ("This was inferred from project structure. Edit to correct if needed.") is described in Solution #3 but has no corresponding criterion.
- The `--apply` flag criterion (#16) says "writes to config on confirm; exit code 0" but does not verify the config was actually written correctly.

**Requirements <-> Solution coherent: 25/25**

11 key scenarios map cleanly to 4 solution components. NFRs are coherent with solution constraints. Priority chain is consistently applied across scenarios, solution, and success criteria. Python exclusion logic is consistent between Key Scenario #3, NFRs, and success criteria.

## Phase 3: Blindspot Hunt

[blindspot] **Coverage percentages are stated without source.** The comparison table verdict claims:

> "~70% of Go stdlib projects (most have `cmd/` or `api/` layout), ~50% of Node stdlib projects (`bin` field or `index.html`), and ~40% of Python stdlib projects (`[project.scripts]` or `app.py`)"

These percentages are presented as data but have no methodology or source. Are they from analyzing a corpus of open-source projects? Expert judgment? Rough estimates? Presenting estimates as facts is a form of evidence inflation.

[blindspot] **The `app.py` exclusion condition may produce false negatives on legitimate CLI projects.** Key Scenario #3 states:

> "`app.py`/`main.py` at root -> `cli` only when no `setup.py` with `name` matching the directory AND no `[project.packages]`/`[tool.setuptools.packages.find]` library-only markers"

Many Python CLI tools are also installable packages — they have both `[project.scripts]` (CLI entry point) and `[project.packages]` or `[tool.setuptools.packages.find]` (for `pip install`). The exclusion condition assumes `[project.packages]` implies "library-only," but CLI tools need package discovery too for installation. This could incorrectly filter out legitimate CLI tools.

[blindspot] **Re-run re-detection does not show the user's previous value as context.** Key Scenario #7 describes re-detection:

> "Re-detect runs full detection + inference pipeline and shows TUI confirmation as usual"

But the TUI does not mention showing the previous value alongside the new inference. If the user previously corrected `cli` to `api`, re-detecting will present `cli (inferred)` again without showing that the user previously chose `api`. The YAML comment in the config file stores the old source, but the proposal does not state that this context is displayed to the user during re-detection.

[blindspot] **No consideration of race conditions in config file access.** The `--apply` flag writes to config, and the re-run flow reads existing config. What if two `forge surfaces detect --apply` processes run simultaneously? Or if `forge init` and `forge surfaces detect --apply` run concurrently? No file locking or atomic writes are mentioned.

[blindspot] **The `index.html` detection for Node.js is fragile.** Key Scenario #2 says:

> "`index.html` at project root (same directory as `package.json`) -> `web`; `inferNodeSurface` does not scan subdirectories for `index.html`"

An `index.html` at root is an extremely weak signal for a "web" surface. Many Node.js projects have a root `index.html` that is a placeholder, documentation page, or unrelated to the project's primary purpose. No evidence is provided that this heuristic has acceptable accuracy.

[blindspot] **Non-interactive terminal output format has no stability contract.** Success criterion #17 defines the stdout format in detail:

> "one line per surface, `<path>=<type> (<source>)`, where `<source>` is `detected:<signal>` or `inferred:<rule-id>` matching the `inferred:<rule-id>` pattern from init summary"

This is a de facto API for scripted consumers. But "Changes to `forge surfaces` CLI command output format" is listed as out of scope, and no versioning or stability guarantee is mentioned. If the output format is a contract, it needs a stability commitment. If it is not, scripted consumers will break.

## Bias Detection Report

**Annotated regions**: 5 attack points / 3 paragraphs = density 1.67
- A1: `forge surfaces detect` output format now defined (criterion #17) — resolved, but non-interactive format has no stability contract (blindspot)
- A1: `forge surfaces detect --apply` concurrent access risk (blindspot)
- A1: `forge surfaces detect` adds scope beyond init flow — now justified by two-touchpoint framing (resolved)
- A2: `index.html` detection is a fragile heuristic (blindspot)
- A3: Re-run re-detection does not show previous value as context (blindspot)

**Unannotated regions**: 4 attack points / ~28 paragraphs = density 0.14
- Coverage percentages without source (blindspot)
- Python `app.py` exclusion may produce false negatives (blindspot)
- TUI hint text described but has no success criterion
- Problem statement density — single paragraph packs too much

**Ratio (annotated/unannotated)**: 11.9

The annotated regions receive significantly higher attack density. This is expected: revisions expand scope and introduce new interactions that require scrutiny. The two-touchpoint reframing (A1 + A3) resolves the bundling inconsistency but introduces new concerns (concurrent access, re-detection context loss).

## Score Summary

```
SCORE: 801/1000
DIMENSIONS:
  Problem Definition: 82/110
  Solution Clarity: 95/120
  Industry Benchmarking: 85/120
  Requirements Completeness: 92/110
  Solution Creativity: 48/100
  Feasibility: 85/100
  Scope Definition: 78/80
  Risk Assessment: 76/90
  Success Criteria: 76/80
  Logical Consistency: 84/90
ATTACKS:
1. [Problem Definition]: Evidence covers 1 of 3 ecosystems -- "evidence is inferred from the Go testing rather than separately reproduced per ecosystem" -- flagged in iterations 1 and 2, still unaddressed. The solution commits to inference for Node.js and Python without evidence that these ecosystems suffer the same problem.
2. [Problem Definition]: Urgency rests on hypothetical harm -- "Each failed detection adds friction that risks user abandonment at the onboarding stage" -- "risks" is not "causes." No churn data, no abandonment data. The 40% failure rate is concrete but the consequence is asserted.
3. [Problem Definition]: Problem statement is a single dense paragraph -- the entire problem, evidence framing, and touchpoint justification are packed into a 6-line block. A reader must re-read to extract the two distinct problems. Must separate into enumerated sub-problems.
4. [Solution Clarity]: Inference rules split across sections -- "infer surface type from project directory structure and conventions" in Solution section does not enumerate the conventions. A reader must cross-reference Key Scenarios to understand what `inferGoSurface` checks. Must move or summarize inference rules in the Solution section.
5. [Solution Clarity]: Ecosystem dispatch mechanism undefined -- "calls ecosystem-specific inference functions" but how is the ecosystem determined? By checking which manifest file exists? Must state the dispatch logic explicitly.
6. [Solution Clarity]: No multi-surface TUI mockup -- the proposal describes single-surface display but not what happens when multiple surfaces are detected simultaneously (e.g., subdir detection showing multiple paths). Must show a multi-surface example.
7. [Industry Benchmarking]: Coverage percentages stated without source -- "~70% of Go stdlib projects... ~50% of Node stdlib projects... ~40% of Python stdlib projects" -- presented as data but no corpus analysis or methodology cited. Must label as estimates or cite source.
8. [Industry Benchmarking]: Proposal follows industry pattern exactly, no differentiation -- "Forge's proposed approach follows this same chain" -- the proposal explicitly adopts the industry standard without articulating what it does differently or better.
9. [Requirements Completeness]: Python `app.py` exclusion may produce false negatives -- "only when no `setup.py` with `name` matching the directory AND no `[project.packages]`/`[tool.setuptools.packages.find]` library-only markers" -- many CLI tools also have `[project.packages]` for installation. The condition assumes package discovery implies library-only. Must refine the exclusion logic or acknowledge this as a known false-negative boundary.
10. [Requirements Completeness]: No accuracy target for inference -- NFRs specify performance (<50ms) and error handling but no accuracy requirement (e.g., "inference must be correct for 90% of known test fixtures"). Without this, the success criteria verify mechanism but not quality.
11. [Solution Creativity]: Explicit innovation disclaimer -- "No special innovation -- this is standard convention-over-configuration" -- limits this dimension by definition.
12. [Feasibility]: YAML Node round-trip constraint not mentioned in feasibility -- the risk table mentions `yaml.Node` round-trip API but the feasibility section does not assess the complexity of switching from `yaml.Marshal`/`yaml.Unmarshal` to `yaml.Node` manipulation for all config writes.
13. [Risk Assessment]: Inference suppression likelihood may be underrated -- "Inference suppression masks valid surfaces | L | M" -- projects with both CLI and API components are common in Go. A Go project with cobra (CLI) and an `api/` directory would suppress the API surface. The L likelihood rating does not match the frequency of multi-surface projects.
14. [Risk Assessment]: No risk for manifest parse errors -- inference functions read manifest files (`pyproject.toml`, `package.json`), but no risk entry addresses what happens when these files are malformed or contain unexpected structures.
15. [Success Criteria]: TUI hint text has no criterion -- Solution #3 describes: "The TUI shows a brief hint: 'This was inferred from project structure. Edit to correct if needed.'" but no success criterion verifies this text appears. Must add a criterion or remove the description.
16. [Logical Consistency]: `--apply` config write not verified -- criterion #16 says "writes to config on confirm; exit code 0" -- but does not verify that the config was actually written correctly. Must add a verification step (e.g., "config file contains the detected surfaces after command exits").
17. [Logical Consistency]: Non-interactive output format is a de facto API without stability guarantees -- criterion #17 defines the stdout format in detail, but "Changes to `forge surfaces` CLI command output format" is out of scope. Must either commit to format stability or mark it as unstable/subject to change.
```

## Verdict

**801/1000 -- REJECT. Below 900 target.**

The proposal has made substantial progress across three iterations. It addressed 14 of 17 iteration-1 attacks and resolved the bundling inconsistency through the two-touchpoint reframing. The expansion from 8 to 19 success criteria, the addition of 3 new risk entries, the explicit task dependency graph, the performance budget, and the filesystem error handling criterion are all significant improvements.

Three categories of issues remain:

1. **Persistent evidence gap** (Problem Definition: 82/110): The solution covers three ecosystems but evidence covers one. This has been flagged in every iteration and remains unaddressed. Either produce Node.js/Python evidence or narrow the problem statement.

2. **Unverified claims presented as data** (Industry Benchmarking: 85/120): Coverage percentages (70%/50%/40%) are stated without source. The Python exclusion condition may produce false negatives on legitimate CLI tools. These are substantive technical concerns.

3. **Missing accuracy target** (Requirements Completeness + Success Criteria): The proposal defines how inference works and how fast it must be, but not how accurate it must be. Without an accuracy target, the success criteria verify mechanism but not quality.

**Priority revisions for iteration 4:**

1. **Source or retract coverage percentages** -- the "~70% of Go stdlib projects" claim needs a source or must be labeled as rough estimates
2. **Add an inference accuracy NFR** -- e.g., "inference must produce correct results for 90% of a test corpus of 20 representative projects per ecosystem"
3. **Produce Node.js/Python evidence** -- or narrow the problem statement to "Go stdlib projects primarily, with Node.js and Python as best-effort"
4. **Refine Python `app.py` exclusion logic** -- acknowledge that `[project.packages]` does not imply library-only; document the false-negative boundary
5. **Add TUI hint text success criterion** -- currently described in Solution #3 but not verified in criteria
6. **Verify `--apply` config write** -- add a criterion that reads config file after `--apply` and asserts detected surfaces are present
