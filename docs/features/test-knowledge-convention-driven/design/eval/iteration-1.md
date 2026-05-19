---
iteration: 1
date: 2026-05-19
score: 810
total: 1000
status: revise
---

# Design Evaluation: test-knowledge-convention-driven — Iteration 1

**Score: 810 / 1000**

## Pre-Score Anchors (Phase 1 — Reasoning Audit)

### Problem → Solution Alignment

The design correctly addresses the stated problem: Profile system couples framework knowledge to Forge binary releases, blocking non-default framework users. The Convention file approach genuinely decouples framework knowledge from the CLI. The solution is sound — no XY problem detected.

### Solution → Evidence

The design cites concrete metrics (126+ existing tests, 85% first-pass compile rate, 19+ consumer files) and maps them to phase gates. This is well-grounded. However, the POC gate (70% on 3 frameworks) is oddly low compared to the Phase 2 gate (85%). If POC only achieves 70%, the approach may be fundamentally weaker than Profile on default frameworks — yet the design promises backward compatibility with zero regression.

### Evidence → Success Criteria

The compile-rate metric is a good proxy but incomplete. The PRD also requires "semantic correctness" (Contract-traceability) and "diff equivalence" (framework core patterns identical). The testing strategy table mentions diff equivalence in the Validation phase but the design has no mechanism to verify semantic correctness — it relies entirely on the compile gate, which cannot catch wrong assertions or missing test steps.

### Self-Contradiction Check

1. The design claims "No new infrastructure" but proposes a new `pkg/forgeconfig/` package. This is a minor terminology issue — it is new Go code but no new tooling/infrastructure. Acceptable.

2. The design says `forgeconfig.ReadAutoConfig` returns `*AutoConfig` (pointer), but the current `profile.ReadAutoConfig` returns `AutoConfig` (value). This is a concrete type mismatch that would break callers if not caught.

3. The design omits `Validation ModeToggle` from the proposed `AutoConfig` struct, but this field exists in the current code and is actively used in `pkg/task/testgen.go` (lines 91-92, 175-176). Dropping it silently would break validation task generation.

---

## Dimension Scores (Phase 2 — Rubric Scoring)

### 1. Architecture Clarity: 150 / 170

**Layer placement explicit (55/60):**
The three-layer diagram (Plugin Skills / Go CLI / Convention Files) clearly shows where each component lives. Each interface (I-1 through I-6) is tagged with its package. Minor deduction: the "Convention Files" layer is described as "in user project docs/conventions/" but the Plugin Skills layer says "distributed to user environment" — the boundary between where Skills read from (user project vs distributed plugin) could be sharper for someone unfamiliar with Forge's distribution model.

**Component diagram present (55/60):**
The ASCII component diagram shows all four skills, all five Go packages, and the Convention file source. Arrows show data flow direction. Minor deduction: the diagram shows `forgeconfig → config.yaml` but does not show the eliminated `pkg/profile/` as a crossed-out or annotated component — the reader has to infer what was removed vs what remains.

**Dependencies listed (40/50):**
Internal dependencies are listed (6 packages). External dependencies: "No new external dependencies" — correct per `go.mod` verification (cobra, testify, yaml.v3). Deduction: the design mentions the existing `consolidate-specs` skill as integration but does not list it as a dependency, nor does it list `pkg/contract/` (which is mentioned in I-2 as "unchanged") in the dependency section. The "unchanged" packages are dependencies of the changed packages and should be enumerated.

---

### 2. Interface & Model Definitions: 140 / 170 (db-schema: "no" sub-rubric)

**Interface signatures typed (50/60):**
All Go interfaces (I-1 through I-5) have typed params and return values. Skill interfaces (I-6) use prose steps — acceptable since skills are LLM prompts, not Go code. Deductions:
- I-1 proposes `ReadAutoConfig(projectRoot string) (*AutoConfig, error)` returning a pointer, but the actual `profile.ReadAutoConfig` returns `AutoConfig` (value type). This is a signature inaccuracy.
- I-5 proposes `ResolveScope(projectRoot string, verb string, scope string) []string` with a new `verb` parameter and `[]string` return, but the actual `ResolveScope(projectRoot, scope string) string` has 2 params and returns a string. This is a significant signature change from current reality that could mislead implementers.

**Models concrete (50/60):**
Convention File data model is detailed with field names, types, and required flags. Config model matches current reality (minus the `Validation` field issue). Deductions:
- `AutoConfig` is missing the `Validation ModeToggle` field that exists in the current `profile.AutoConfig` and is actively used by `pkg/task/testgen.go`. This is not a "removed field" — the design's FS-6 and I-3 section never mention removing `Validation`. This is a gap that would cause compilation failure.
- `BuildIndexOpts` in I-3 shows `AutoConfig *forgeconfig.AutoConfig` (pointer), but the current `task.BuildIndexOpts` has `AutoConfig profile.AutoConfig` (value). Inconsistent without explanation.
- `TestGenerationOpts.Facts` is typed as `string` — this is a large blob (Fact Table content). The design does not explain its format or how Skills should parse it.

**Directly implementable (40/50):**
A developer could implement the Go interfaces from this design with one important exception: the `AutoConfig` struct is incomplete (missing `Validation`). The Convention File model is well-specified enough for skill prompt construction. The skill interfaces (I-6) are steps, not APIs — a developer would need to read existing skill files to know how to implement them, which is fine for skills but the design does not reference the existing skill methodology format.

---

### 3. Error Handling: 115 / 130

**Error types defined (40/45):**
The error handling table defines 8 error contexts with specific behaviors and user-facing messages. No custom error types or error codes are defined — the design states "No new error types needed" which is consistent with the simplified approach. Deduction: the design does not specify Go-level error wrapping or sentinel errors for the `forgeconfig` package (e.g., `ErrConfigNotFound`, `ErrInvalidConfig`). The current `profile` package has error paths that callers handle; the replacement should enumerate its error conditions.

**Propagation strategy clear (40/45):**
Two-layer propagation is clearly stated: Skills report text messages; Go CLI uses exit code 1 + stderr. This aligns with BIZ-error-reporting-001 (exit code 1 = failure) and TECH-error-handling-001 (stderr for errors). Deduction: the design does not address what happens when `forgeconfig.ReadConfig` fails at the Go level — does it return a wrapped error? A sentinel error? The current `profile.ReadConfig` has nuanced error handling (file not found vs malformed YAML vs unknown fields). The replacement should specify.

**HTTP status codes mapped (35/40):**
N/A for CLI tool. However, the design does not map CLI exit codes to specific error types. BIZ-error-reporting-001 specifies exit code 2 for usage errors; the design should confirm that the simplified commands (init, config, task index, task add) still respect this convention after their rewrite.

---

### 4. Testing Strategy: 110 / 130

**Per-layer test plan (40/45):**
The phase-based test table covers: POC compile rate, Phase 1 regression (`go build`, `go test`), Phase 2 end-to-end pipeline, Phase 3 functional, and Validation diff equivalence. Deduction: the design explicitly states "No new unit tests for the removed code" and "the testing strategy is 'everything that worked before must still work.'" This is reasonable for Phase 1 (Profile removal) but Phase 2 (Skill rewrites) introduces significant new behavior (Convention loading, Code Reconnaissance, compile gate) that deserves unit-level testing of the skill prompt changes. The design does not explain how to test the Skill layer (LLM prompts) at a unit level.

**Coverage target numeric (35/45):**
85% first-pass compile rate is a clear numeric target for the end-to-end pipeline. Deduction: there is no unit test coverage target for the Go CLI layer. The forge-cli CLAUDE.md mandates 80%+ coverage, but the design explicitly excludes "Unit test coverage targets for rewritten Go code" from scope. This contradicts the project's own TDD mandate. The design should at minimum acknowledge this tension and explain why rewritten packages are exempt.

**Test tooling named (35/40):**
`just e2e-compile`, `go build ./...`, `go test ./...`, `diff` are named. The Key Test Scenarios section references `grep -r "pkg/profile"` for the import audit gate. Deduction: no specific assertion library or test framework is named for the Go unit test level (the project uses `testify` per go.mod). The diff equivalence testing approach is described conceptually but no tool is specified for automated pattern comparison.

---

### 5. Breakdown-Readiness: 165 / 180

**Components enumerable (60/65):**
All components can be counted: 1 new package (`pkg/forgeconfig/`), 5 simplified packages, 4 skills (1 new), 1 deleted package (`pkg/profile/`). Deduction: the PRD FS-7 mentions "19+ consumers" that need rewriting, but the design does not enumerate them. The Related Changes table lists 5 change points, which is coarser than the 19+ file-level consumers. For task breakdown, knowing which specific files in `internal/cmd/` need changes would be valuable.

**Tasks derivable (60/65):**
Each interface maps to at least one implementation task: I-1 (forgeconfig package creation), I-2 (journey simplification), I-3 (task simplification), I-4 (e2e simplification), I-5 (just simplification), I-6 (4 skill rewrites/creations). The phase structure provides a natural ordering. Deduction: the `forge test promote` change (regex-based tag discovery) is mentioned in I-2 but has no dedicated interface or task. It is a behavioral change to an existing command that is easy to overlook during breakdown.

**PRD AC coverage (45/50):**
The PRD Coverage Map table maps all 9 functional specs and all 6 user stories to design components. Spot-checking:
- FS-1 → Convention File data model: covered
- FS-4 → Compile gate in I-6: covered
- FS-7 → Profile removal: covered
- Story 6 (Compile gate recovery) → Error handling table: covered
- Story 5 (Multi-framework) → "Convention domains filtering" in Skill prompt: covered at prose level

Deduction: Story 5's acceptance criterion "both files are loaded and merged; last-loaded wins for conflicting fields; the skill logs a note about the domain overlap" is addressed only at the prose level in the error handling table row "Convention vs Reconnaissance conflict." The merging logic and "last-loaded wins" strategy for multiple Convention files with overlapping domains is not modeled in any interface or data structure — it lives entirely in skill prompt behavior. This is borderline acceptable for LLM-driven skills but risks inconsistency across runs.

---

### 6. Security Considerations: 70 / 80

**Threat model present (35/40):**
One threat identified: "Convention file injection" — malicious Convention content instructing the LLM to generate harmful code. The threat is specific and relevant. Deduction: the threat model is narrow. It does not consider: (a) what "harmful code" means in the context of test generation — could a Convention file cause the generated test to exfiltrate environment variables, access files outside the test scope, or create network connections? (b) the compile gate is presented as a mitigation but compilation success does not prevent runtime malicious behavior.

**Mitigations concrete (35/40):**
Three mitigations listed: user-controlled files, compile gate, no new attack surface vs Profile. The comparison to Profile is effective for risk framing. Deduction: the compile gate is described as a "sandbox" but it is not a sandbox — it verifies compilation, not runtime behavior. A test file that compiles and contains `os.Readfile("/etc/passwd")` would pass the compile gate. This is a minor concern for an internal tool but the characterization is inaccurate.

---

### 7. Implementation Feasibility: 110 / 140

**Dependencies available (35/50):**
The design correctly states "No new external dependencies" which matches `go.mod` verification (cobra, testify, yaml.v3). However:
- (-15 pts) The design proposes `pkg/forgeconfig/` as a new package, but this package does not exist yet. The design does not specify whether it will be created by extracting code from `pkg/profile/config.go` or written fresh. The extraction path has a pitfall: `profile.ReadAutoConfig` returns `AutoConfig` (value) with a `raw` field for tracking which YAML keys were present — the design's `forgeconfig.ReadAutoConfig` returns `*AutoConfig` (pointer) without the `raw` field. This suggests the design author may not have examined the actual implementation closely enough.
- (-0 pts) No convention violations detected — the design correctly uses existing patterns.

**Architecture fits project structure (40/50):**
The proposed architecture fits the existing `cmd → internal → pkg` dependency direction. The layer separation (Go CLI vs Plugin Skills) matches the forge-distribution model. Deductions:
- (-10 pts) The `pkg/just/` change (I-5) proposes a try-based resolution: "attempt `just <verb> <scope>` — if recipe not found, fall back to `just <verb>` without scope argument." This introduces implicit shell-level error handling that was previously explicit in Go code. The current `ResolveScope` uses `profile.ReadConfig` to make a deterministic decision; the proposed replacement relies on `just` execution failure as a control flow mechanism. This is a pattern change that contradicts the project's preference for explicit logic.
- (-0 pts) The overall package restructuring is sound.

**Technical claims grounded (35/40):**
The 85% first-pass compile rate claim is grounded in the 126+ existing test baseline. The POC 70% gate is appropriately conservative. The "generation time should not increase by more than 20%" performance requirement is stated but not grounded — no measurement methodology or baseline is provided. Deduction: the design does not state the current generation time, making the 20% cap unverifiable.

---

## Phase 3 — Blindspot Hunt

### [blindspot-1] Missing `Validation` field in `AutoConfig` — a real implementation blocker

The design's I-1 `AutoConfig` struct contains 4 fields (`E2eTest`, `ConsolidateSpecs`, `CleanCode`, `GitPush`) but the actual `profile.AutoConfig` has 5 fields — the design omits `Validation ModeToggle`. This field is actively used:
- `pkg/task/testgen.go:92`: `if auto.Validation.Full` — gates validation task generation
- `pkg/task/testgen.go:176`: `if auto.Validation.Quick` — gates quick validation tasks

The PRD's FS-6 states "Remove `languages`, `interfaces`, `test-framework`, `project-type`" but does NOT mention removing `Validation`. The design's FS-6 section says "Retain `auto.*`, `test-command`, `worktree`" — `Validation` is under `auto.*` and should be retained. This is a concrete omission that would cause `go build` failure in Phase 1.

### [blindspot-2] `ReadAutoConfig` return type discrepancy signals shallow code analysis

The design proposes `ReadAutoConfig(projectRoot string) (*AutoConfig, error)` returning a pointer. The actual implementation returns `AutoConfig` (value) because the `AutoConfig` struct contains a `raw map[string]map[string]bool` field used by `applyDefaults` to distinguish "explicitly set to false" from "not present in YAML." The design's `AutoConfig` omits this `raw` field, which means the default-application logic would also be lost. This pattern (value-type return with internal state tracking) is a deliberate design choice in the current code that the proposed replacement discards without discussion.

### [blindspot-3] `BuildIndexOpts` redesign is incomplete

The design's I-3 `BuildIndexOpts` drops `Languages`, `TestInterfaces`, `ResolveStrategy` and keeps only `AutoConfig`. But the current `BuildIndexOpts` (verified in `pkg/task/build.go:19-28`) has 7 fields including `FeatureSlug`, `ProjectRoot`, `TasksDir`, `IndexPath`. The design only shows `AutoConfig` — it does not address whether these other fields are retained, removed, or restructured. This is a partial interface definition that an implementer cannot code from.

### [blindspot-4] `ResolveScope` signature change is underspecified

The current `ResolveScope(projectRoot, scope string) string` reads `profile.ReadConfig` to determine project type and returns either the scope or an empty string. The design proposes `ResolveScope(projectRoot string, verb string, scope string) []string` with try-based resolution. This is a significant behavioral change:
1. New `verb` parameter — what values does it take? How does it affect resolution?
2. Return type changes from `string` to `[]string` — what does each element represent? Multiple scopes?
3. "Try-based resolution" means calling `just` as a side effect — this is a fundamentally different approach from the current deterministic config-based resolution. The design should explain the rationale for this behavioral change.

### [blindspot-5] Skill testing strategy gap

The design states "No new unit tests for the removed code" and excludes unit test coverage for rewritten Go code from scope. But the forge-cli CLAUDE.md mandates "TDD (mandatory)" and "Coverage target 80%+." The design explicitly contradicts a project-level convention without acknowledging the tension or providing justification beyond "everything that worked before must still work." The rewritten packages (`pkg/forgeconfig/`, simplified `pkg/task/`, etc.) will contain new logic that deserves new tests.

### [blindspot-6] Convention file merging semantics are undefined

PRD Story 5 states: "both files are loaded and merged; last-loaded wins for conflicting fields." The design's error handling table mentions "Convention vs Reconnaissance conflict" but does not address Convention vs Convention conflict. What is the merge strategy? Field-level? Section-level? What about partial overlaps (e.g., two files both declaring an Assertion section with different values)? This is a behavioral specification that lives entirely in the LLM prompt and has no concrete definition.
