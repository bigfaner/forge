# PRD Evaluation Report — Iteration 1

**Feature:** test-knowledge-convention-driven
**Date:** 2026-05-19
**Mode:** B (no UI)
**Scorer:** Senior QA Engineer (adversarial)

---

## Phase 1: Reasoning Audit — Pre-Score Anchors

1. **Problem-Solution fit is sound.** The Profile system's conceptual conflict with Journey-Contract model is real, and Convention files + Code Reconnaissance directly address both the rigidity and the responsibility-overlap problems.

2. **Flow diagram coverage is incomplete.** The PRD describes 3 flows (main generation, bootstrap, convention creation) but only the main generation flow has a Mermaid diagram.

3. **Cold-start success rate claim is self-contradictory.** Story 3 claims >= 85% first-pass compile rate on a project with no Convention files and no existing tests. In that scenario, Code Reconnaissance has zero signals to collect, so the 85% target rests entirely on "LLM defaults" — which is exactly what the Background section says fails for non-default frameworks.

4. **Convention file structure is named but never defined.** FS-1 lists section names but provides no schema, no field types, no required vs optional constraints beyond a brief mention. This is the core artifact of the entire feature and it remains underspecified.

5. **"Silent migration" vs "remove fields" tension.** Scope says "remove fields from config.yaml" while Story 4 and FS-8 say "silently ignored." These describe different behaviors: schema removal vs parser tolerance.

---

## Phase 2: Rubric Scoring

### Dimension 1: Background & Goals — 93/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Three elements (Reason/Target/Users) | 28/30 | All three present. Reason has 3 specific problems. Target is clear. Users have 2 named segments. Docked 2: user segments are defined only by framework choice; no mention of mixed-language projects (which Story 5 addresses), suggesting the user analysis is incomplete. |
| Goals quantified | 30/30 | Six goals, all with numeric or countable targets (85% compile rate, < 5 min bootstrap, zero imports, 126+ tests pass, 20% generation time ceiling). |
| Logical consistency | 35/40 | Goals follow from problem. One tension: Background says Profile handles "(a) project tech stack detection, (b) framework-specific code generation knowledge, (c) test organization axis" and that (c) is already solved by Journey-Contract. But the scope includes rewriting `pkg/journey/`, implying deeper coupling than Background acknowledges. |

### Dimension 2: Flow Diagrams — 125/150

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Mermaid diagram exists | 50/50 | One valid Mermaid flowchart with `flowchart TD`. |
| Main path complete | 35/50 | Main generation flow is covered start-to-end. However, the PRD text describes 3 flows (main, bootstrap, convention creation) and only the main flow has a diagram. The bootstrap flow and convention creation flow have no visual representation, forcing readers to parse text-only descriptions for those critical paths. |
| Decision points + error branches | 40/50 | Two decision diamonds (`Convention files found?`, `Retries < 2?`) and error branches (`No Convention found` hint, compile fail → retry, blocked terminal state). Adequate for the main flow but missing decision points for the undocumented flows. |

### Dimension 3: Flow Completeness (Mode B) — 160/200

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Complete business process | 55/70 | Three flows described with numbered steps. State transitions are implicit but not formalized. The main flow covers trigger → Convention loading → Reconnaissance → generation → compile gate → success/blocked. However, state transitions for Convention files themselves (create → update → version) are entirely absent. What happens when a user re-runs test-guide after updating a Convention file is not described. |
| Data flow | 70/70 | "N/A — single-system feature." Auto full-score per rubric. |
| Exception handling | 35/60 | Compile failure with retry is documented. "No Convention found" path is documented. Missing: (1) Convention file exists but is malformed or has invalid frontmatter — no error path. (2) `just e2e-compile` does not exist (no justfile) — no error path. (3) Reconnaissance finds conflicting signals (e.g., mixed framework indicators in test files) — no error path. (4) Convention file has empty or incomplete sections — no handling. These are real failure modes for a CLI tool that an agent must execute. |

### Dimension 4: User Stories — 160/200

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Coverage per user type | 35/50 | Background defines 2 user segments. Story 1 covers non-default users. Story 3 covers cold-start users (overlaps both segments). Story 5 introduces "mixed-language project user" — a user type not present in the Background section. This is a cross-section inconsistency: the Background's user analysis is incomplete relative to the stories. |
| Format correct | 45/50 | All 5 stories use As a / I want / So that format. Actions are concrete and specific. Minor deduction: Story 4's "I want to continue using my existing project without any migration steps" is slightly vague about what "continue using" entails operationally. |
| AC per story (Given/When/Then) | 45/50 | All stories have at least one AC in Given/When/Then format. Story 2's AC uses Given/When/Then + And, which is acceptable. Minor: Story 2 has only one AC scenario; no negative path (e.g., what if user rejects detected patterns?). |
| AC verifiability & boundaries | 35/50 | Most ACs are objectively testable. Issues: (1) Story 1 "uses ginkgo imports, assertions, and style conventions as declared in the Convention file" — "style conventions" is subjective without a precise definition of what counts as matching. (2) Story 3 ">= 85% first-pass success rate" — measurement methodology not specified (which 126+ tests? all of them? a subset?). (3) No AC covers error paths: Convention load failure, Reconnaissance failure, compile gate failure after retries. Only happy-path and "silently ignored" scenarios are covered. |

### Dimension 5: Scenario Completeness — 90/150

| Criterion | Score | Justification |
|-----------|-------|---------------|
| End-to-end scenario coverage | 40/60 | Three flows described, but coverage gaps exist. Missing scenarios: (1) Convention file update after initial creation — no flow for re-running test-guide or editing Convention. (2) Multiple Convention files for same language domain (conflict resolution). (3) Convention file with domains that don't match any Journey — no handling. (4) gen-test-scripts when no justfile exists — no scenario. |
| Implicit assumptions surfaced | 20/40 | Several unstated assumptions: (1) `just e2e-compile` recipe always exists — never stated as prerequisite. (2) `docs/conventions/` directory exists or is auto-created — never specified. (3) LLM can reliably generate correct test code from Convention + Reconnaissance — this is the core hypothesis and is never examined. (4) Code Reconnaissance works across all supported languages — what signals for each language? (5) Convention file content is always valid markdown with correct frontmatter — assumed but never enforced. |
| Business-rules consistency | 30/50 | BIZ-quality-gate-001 specifies "retry-once policy" for tests. PRD's compile gate says "max 2 retries." If the compile gate is part of the quality-gate pipeline, this is a conflict. The PRD also involves rewriting `pkg/task/` (scope item) but no AC verifies that task state transitions (BIZ-task-lifecycle-001) remain correct after the rewrite. BIZ-error-reporting-001/002: the compile gate failure path feeds errors to LLM but does not specify user-facing error format or exit codes for the gen-test-scripts skill. |

### Dimension 6: Edge Case Coverage — 45/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Error paths documented | 20/40 | Compile failure → retry is documented. "No Convention found" is documented. Missing: (1) Malformed Convention file — no error handling. (2) Convention file with empty required sections — no handling. (3) Reconnaissance finds zero signals (no test files, no framework indicators) — Story 3 covers this as success, not as a potential problem. (4) `just e2e-compile` not found — no error path. (5) Permission denied reading Convention files — no error path. |
| Boundary conditions covered | 15/35 | Multiple Convention files (Story 5) and zero Convention files (Story 3) are covered. Missing: (1) Multiple Convention files matching the same domain (conflict). (2) Very large Convention file or very many Convention files. (3) Convention file with all optional sections vs only minimum set. (4) Concurrent gen-test-scripts invocations reading the same Convention. (5) Unicode or special characters in Convention file content. |
| Failure recovery described | 10/25 | Compile gate retry mechanism provides partial recovery. Missing: (1) After "Blocked: compile gate failed" — what does the user do? No recovery steps. (2) Convention file write failure during test-guide — no recovery. (3) User rejects detected patterns during test-guide — no recovery path (just implied "try again"). |

### Dimension 7: Scope Clarity — 83/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Concrete deliverables | 30/35 | 10 in-scope items, most are specific (named packages, named CLI commands, named skill). "Convention file fixed structure definition" is ambiguous — is it a document, a Go struct, a JSON schema? |
| Out-of-scope items named | 28/30 | 7 explicitly named out-of-scope items. Clear and specific. |
| Scope consistent with stories | 25/35 | Most scope items map to stories. Gap: "Integrate Convention files into consolidate-specs management" is in scope but has no corresponding user story or AC. This is a deliverable with no verification. |

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Cold-start 85% claim is self-defeating

**Quote:** "Given a new project with no Convention files and no existing test files ... Then ... `just e2e-compile` passes with >= 85% first-pass success rate" (Story 3 AC)

**Issue:** This AC is the feature's biggest risk and it rests on a logical contradiction. The Background states that Profile "only works correctly for projects using default frameworks" and that "non-default frameworks get wrong results." The entire motivation for Convention files is that LLM defaults are insufficient. Yet Story 3 asserts that in the worst case (no Convention, no existing tests to scan), the system achieves 85% — relying entirely on the LLM defaults that the feature was created to supersede. If LLM defaults can hit 85% on cold start, the Convention system's value proposition is undermined. If they cannot, the AC is unachievable. The document never acknowledges this tension.

### [blindspot-2] Convention file structure is underspecified for its centrality

**Quote:** "Convention files use fixed sections: Framework, Assertion, Tags, Result Format (minimum set). Optional sections: Helpers, Import Patterns, Code Style, Anti-patterns." (FS-1)

**Issue:** The Convention file is the single most important artifact in this feature — it replaces an entire Profile system. Yet FS-1 provides only a list of section names with no field-level specification. What goes in the Framework section? What is the required format? How are Tags specified (YAML list? prose? code blocks?)? What validation rules apply? Without this specification, downstream agents cannot verify Convention file correctness, and test-guide cannot be validated against a schema. This is the equivalent of designing a database migration without a schema.

### [blindspot-3] No mechanism to detect missed Profile consumers

**Quote:** "Delete `pkg/profile/` entirely. Rewrite all 19+ consumers." (FS-7)

**Issue:** The "19+" count suggests uncertainty about the actual number of consumers. The scope lists 5 packages and CLI commands to rewrite, but Forge is a plugin-based system. Other plugins, hooks, external scripts, or user project configurations may import or depend on `pkg/profile/` functionality. The PRD provides no audit step (e.g., `grep -r "pkg/profile" .` or go reference analysis) to ensure complete consumer identification. Deleting `pkg/profile/` with an incomplete consumer list would cause build failures in unlisted files.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Background & Goals | 93 | 100 |
| 2. Flow Diagrams | 125 | 150 |
| 3. Flow Completeness | 160 | 200 |
| 4. User Stories | 160 | 200 |
| 5. Scenario Completeness | 90 | 150 |
| 6. Edge Case Coverage | 45 | 100 |
| 7. Scope Clarity | 83 | 100 |
| **Total** | **756** | **1000** |

**Passes gate (900)?** No. Needs revision.
