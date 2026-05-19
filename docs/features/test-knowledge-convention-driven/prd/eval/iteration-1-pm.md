# PRD Evaluation Report — test-knowledge-convention-driven

**Iteration**: 1
**Evaluator**: Senior PM (Mode B — no UI)
**Date**: 2026-05-19
**Total Score**: 735 / 1000

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Anchor 1: Compile-rate target is an unsubstantiated assertion
The document sets "First-pass compile rate >= 85%" as a goal (prd-spec.md line 39) but provides zero evidence that removing hardcoded Profile knowledge and replacing it with LLM defaults + Convention files will achieve this. The current Profile system presumably achieves a high compile rate precisely because it hardcodes correct patterns. The 85% figure appears aspirational rather than derived from any measurement or pilot.

### Anchor 2: Success criteria test compile, not correctness
The measurable success criteria focus entirely on compilation (just e2e-compile passes). A generated test file that compiles but uses wrong assertions, missing test steps, or incorrect test structure would count as "success." The document acknowledges this indirectly by noting "Style differences (variable naming, comments) allowed" but does not define the boundary between acceptable style variance and semantic divergence.

### Anchor 3: "Silent migration" contradicts "full Profile removal"
The document claims "no migration needed" (line 153, FS-8) and that old config fields are "silently ignored" (line 146, FS-6). But removing 19+ files, 4 CLI commands, and multiple config fields IS a migration — it's just one where the tool does not help the user. This framing hides the real user impact.

### Anchor 4: Auto-detection remains but mechanism is unspecified
The scope retains `auto.*` fields in config.yaml (line 50) and the test-guide command does "project signal detection" (line 95). This means auto-detection survives the Profile removal, but it is unclear whether `auto.*` is re-implemented or whether it's a separate mechanism. The document does not explain how auto-detection works without Profile.

---

## Phase 2: Dimension Scoring

### 1. Background & Goals — 75/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Background has three elements | 25/30 | Reason, Target, and Users are all present. Reason identifies 3 concrete problems with specific examples (ginkgo, vitest). Target describes the Convention file replacement. Users identify two segments. Deduction: the "Default framework users" segment is described as getting "no benefit from Profile over LLM defaults" — this is a claim about the current system's value that is not substantiated, making the user segmentation feel assumed rather than validated. |
| Goals are quantified | 25/30 | Four of six goals have numeric targets: 85% compile rate, 0 imports, <5 min bootstrap, 126+ tests pass. However, "Users can use ginkgo/vitest/pytest without Forge code changes" is not quantified — it is a capability claim, not a measurable metric. "Framework core patterns identical" is quantified only loosely by "identical" with no tolerance defined. |
| Background and goals logically consistent | 25/40 | The logical chain from "Profile hardcodes decisions upfront" to "replace with user-editable Convention files" is sound. However, the goal "First-pass compile rate >= 85%" does not clearly follow from the stated problem. The problem is about non-default frameworks failing; the metric measures all generation. A project using default frameworks might see compile rate go DOWN after Profile removal, offsetting gains from non-default. The document does not address this risk. Also, the goal "Generated code diff equivalence" contradicts the claim that Profile "hardcodes technical decisions" — if the old output was wrong for non-default frameworks, why measure equivalence? |

### 2. Flow Diagrams — 120/150

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Mermaid diagram exists | 45/50 | A Mermaid flowchart is present (lines 103-115) with proper start/end nodes, decision diamonds, and labeled edges. Clear and readable. Minor deduction: only one diagram covers the main gen-test-scripts flow; the bootstrap flow and test-guide flow are described in text but lack their own diagrams. |
| Main path complete | 40/50 | The diagram covers the main happy path: start -> load convention -> reconnaissance -> generate -> compile -> success. However, step 7 from the text ("If no Convention found -> output hint, proceed with LLM defaults") IS in the diagram (LoadConv -> No -> Hint -> Recon) but the "hint" node does not specify what happens next — it merges back into Recon, which is correct but the diagram does not show that the hint is informational only. |
| Decision points + error branches | 35/50 | Two decision diamonds exist: "Convention files found?" and "Retries < 2?". One error branch: "Blocked: compile gate failed." Missing: what happens when Code Reconnaissance finds no existing test files (the cold-start scenario). The diagram assumes Recon always succeeds. Missing: the test-guide flow has no diagram at all despite being a distinct interaction flow. |

### 3. Flow Completeness — 145/200

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Flow steps describe complete business process | 50/70 | Three flows are described: main generation, bootstrap, and test-guide creation. The main flow is thorough (7 steps from trigger to end state with state transitions). However, the bootstrap flow (lines 87-90) is skeletal — 4 steps without any error handling. What if init-justfile fails? What if the first gen-test-scripts fails in the cold start? What if test-guide mis-detects the framework? These are not addressed. |
| Data flow documented | 70/70 | Single-system feature (forge-cli only). Data Flow Description explicitly states N/A with justification. Auto-full-score per rubric. |
| Exception handling and edge cases | 25/60 | The main flow documents compile failure with retry logic (max 2 retries). But this is the ONLY error path documented across all three flows. Missing error paths: (1) Convention file malformed or has invalid frontmatter — no handling described. (2) Code Reconnaissance fails to find any test patterns in existing files (partial match, not cold start). (3) test-guide user rejects detected patterns — what happens? (4) just e2e-compile command itself is not installed or fails to run (infra error, not compile error). (5) Multiple Convention files with conflicting guidance — which wins? |

### 4. User Stories — 155/200

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Coverage: one story per target user | 40/50 | Two user segments identified: "Default framework users" and "Non-default framework users." Stories cover: non-default framework (Story 1), new project/cold start (Story 2, 3), upgrading user (Story 4), multi-framework user (Story 5). Missing: no story for "Default framework users" specifically — the segment identified in the background. Story 3 (cold start) partially covers them but focuses on "no Convention, no tests" rather than "default framework, Profile works for me, don't break my workflow." Story 4 covers upgrade but from config perspective, not from "my generated tests should still be correct" perspective. |
| Format correct | 40/50 | All 5 stories use "As a / I want / So that" format. Actions are concrete: "generate e2e tests that use my project's actual framework," "create a test Convention file," "generate e2e tests that compile." However, Story 5's "So that" clause ("test generation uses the correct framework knowledge for each Journey's interface type") describes HOW the system works internally, not user value. |
| AC per story (Given/When/Then) | 40/50 | Every story has at least one AC in Given/When/Then format. Stories 1, 3, and 5 have multi-step ACs with "And" continuations. Deduction: Story 4's AC is weak — "those fields are silently ignored (no errors, no warnings, no migration prompts)" is three negations packed into one "Then." This should be split into separate ACs for clarity. |
| AC verifiability & boundary coverage | 35/50 | Verifiability: Story 1's AC is well-verifiable (check imports, run compile). Story 3's ">= 85% first-pass success rate" is verifiable but is a population metric, not a per-invocation criterion — you cannot verify it from a single run. Story 5's AC about "only the Go Convention file is loaded" requires internal inspection, not observable behavior. Boundary coverage: No AC covers error cases. No AC covers the case where Convention file exists but is incomplete or wrong. No AC covers concurrent Convention file access. |

### 5. Scenario Completeness — 105/150

| Criterion | Score | Justification |
|-----------|-------|---------------|
| End-to-end scenario coverage | 45/60 | The three flows (main, bootstrap, test-guide) cover the primary user-facing scenarios. Each describes a lifecycle, though the bootstrap flow is thin. Missing scenario: what happens when a user edits a Convention file after initial creation — does gen-test-scripts pick up changes? Is there validation? Also missing: the scenario where a user runs gen-test-scripts for a Journey that spans multiple Convention domains (e.g., a Journey with both Go and TypeScript interfaces). |
| Implicit assumptions surfaced | 25/40 | Several assumptions are not surfaced: (1) Convention files are written by hand or test-guide — no validation on file structure is described. (2) The document assumes `just` is installed and available, but does not state this as a prerequisite. (3) The "Code Reconnaissance" step assumes existing test files follow recognizable patterns — what if they are unconventional? (4) The document assumes LLM defaults will produce compilable code >= 85% of the time without explaining what "defaults" means when Profile is gone. |
| Business-rules consistency | 35/50 | Checking against injected context: BIZ-quality-gate-001 defines a multi-phase pipeline (compile -> tests -> e2e). The PRD's "compile gate" (FS-4) only describes the compile step, not the full quality gate pipeline. The PRD mentions "just e2e-compile" as the gate but the quality gate rule says compile is just phase 1. This is a gap: the PRD should clarify how the Convention system interacts with the full quality gate pipeline, not just compilation. BIZ-task-lifecycle-001/002 are not violated. No contradictions with error reporting rules. |

### 6. Edge Case Coverage — 55/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Error paths documented | 20/40 | Only one error path is explicitly described: compile failure with retry (FS-4). Missing: malformed Convention file, missing docs/conventions/ directory, test-guide detection failure, just not installed, Permission denied on Convention file write, network/LLM failure during generation. |
| Boundary conditions covered | 20/35 | Some boundaries mentioned: "max 2 retries," "minimal structure" for Convention. Missing: empty Convention file (all sections blank), Convention file with conflicting directives (e.g., two Framework declarations), very large Convention file, Convention file referencing a framework that does not exist in the project, zero existing test files (partially covered by cold start but not as a boundary of Code Reconnaissance). |
| Failure recovery described | 15/25 | Compile failure recovery is described (retry with error feedback). "Blocked: compile gate failed" is the terminal failure state. But what does the user DO when blocked? The document does not say. For test-guide: if user rejects detected patterns, what next? No recovery described. |

### 7. Scope Clarity — 80/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete deliverables | 30/35 | All 11 in-scope items are specific: "Remove pkg/profile/ directory entirely," "Rewrite gen-test-scripts skill," "New /forge:test-guide slash command," etc. Each names a specific module or command. Minor deduction: "Convention file fixed structure definition" is slightly vague — a "definition" is a document/artifact, not a deliverable feature. |
| Out-of-scope explicitly lists deferred items | 25/30 | Six items explicitly listed as out-of-scope: gen-journeys, gen-contracts, forge test verify/promote, existing test migration, anti-pattern docs, auto-sync/watch mode, unit test coverage. These are named, not implied. Minor deduction: "Unit test coverage targets for rewritten Go code" is listed as out-of-scope, but 19+ files being rewritten with no test coverage requirement is a risk not acknowledged. |
| Scope consistent with functional specs and user stories | 25/35 | The in-scope "Rewrite Go packages: pkg/journey/, pkg/e2e/, pkg/just/, pkg/task/" and "Rewrite CLI commands" match FS-7 and the Related Changes table. However, the user stories do not cover the full scope: FS-6 (config.yaml cleanup) and FS-7 (Profile removal) are infrastructure changes with no corresponding user story that describes the developer experience of these removals from a usage perspective (Story 4 covers config upgrade but not the full removal impact). Also, "Integrate Convention files into consolidate-specs management" is in scope but has zero coverage in user stories and no functional spec. |

---

## Phase 3: Blindspot Hunt

### [blindspot-1] No rollback path for compile gate failures at scale
The document defines "Point of no return: Phase 3 start" (line 189) as the rollback boundary. But the 85% compile rate target is a population metric measured on 126+ tests. If Phase 3 starts and the compile rate is only 70%, there is no documented rollback strategy. The document says "Phase 1 and 2 are independently revertible" but does not define what triggers a revert or who makes that decision. This is a go/no-go criteria gap that no rubric dimension captures.

**Quote**: "Point of no return: Phase 3 start (test-guide depends on Profile-free environment)." (prd-spec.md line 189)

### [blindspot-2] Convention file authoring UX is undefined
The document describes Convention file structure (FS-1) and a creation tool (FS-5), but never addresses the editing lifecycle. Convention files are markdown — what happens when they contain syntax errors? What happens when the `domains` frontmatter is misspelled? What happens when a Convention file references a framework that does not match the project? The entire feature assumes Convention files are always correct, but provides no validation, no error messages, and no debugging path. This is a user experience gap that falls between Scope Clarity and Edge Case Coverage.

**Quote**: "Convention files use fixed sections: Framework, Assertion, Tags, Result Format (minimum set)." (prd-spec.md line 126) — Fixed structure defined but no validation of conformance.

### [blindspot-3] "Code Reconnaissance" is an undefined capability
FS-3 mentions extending the "Fact Table (runtime LLM notes, not persisted)" to collect test framework info. But this is described as a one-line bullet: "file patterns, import analysis, build tag analysis, function signature patterns." How reliable is this? What if the project has 0 test files? What if test files use multiple frameworks? What if imports are aliased? The entire cold-start and fallback path depends on this capability, but it is treated as trivial. No rubric dimension captures the risk of an under-specified core capability.

**Quote**: "Fact Table (runtime LLM notes, not persisted) extended to collect test framework info: file patterns, import analysis, build tag analysis, function signature patterns." (prd-spec.md line 134)

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 75 | 100 |
| 2. Flow Diagrams | 120 | 150 |
| 3. Flow Completeness | 145 | 200 |
| 4. User Stories | 155 | 200 |
| 5. Scenario Completeness | 105 | 150 |
| 6. Edge Case Coverage | 55 | 100 |
| 7. Scope Clarity | 80 | 100 |
| **Total** | **735** | **1000** |
