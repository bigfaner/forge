# Eval-Proposal: Iteration 2

**Score: 830/1000** (target: 900)

## DIMENSIONS

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 100 | 110 |
| 2. Solution Clarity | 110 | 120 |
| 3. Industry Benchmarking | 85 | 120 |
| 4. Requirements Completeness | 95 | 110 |
| 5. Solution Creativity | 50 | 100 |
| 6. Feasibility | 90 | 100 |
| 7. Scope Definition | 75 | 80 |
| 8. Risk Assessment | 70 | 90 |
| 9. Success Criteria | 75 | 80 |
| 10. Logical Consistency | 80 | 90 |

## ATTACKS

**[Dimension 1] Section "Urgency" -- urgency remains qualitative, not quantified with concrete cost-of-delay data.**
Quote: "Restructuring now, while the skill is still evolving, prevents accumulating architectural debt."
This is still a generic "do it before it gets worse" argument. How many new profiles or types are planned in the next quarter? How many times has the monolith caused a bug or delayed a change? There is no metric like "each new type adds ~80 lines to the monolith" or "the convention bug (line 50) took X hours to diagnose." The urgency is plausible but not quantified with time, dollar, or incident data. (-10 pts)

**[Dimension 3] Section "Alternatives & Industry Benchmarking" -- trade-off comparison is shallow on quantitative analysis.**
Quote for Config-driven dispatch Cons: "Adds a new file format that the agent must parse; YAML frontmatter in SKILL.md already exists but a separate dispatch config introduces another loading step."
This con is asserted without evidence. How much harder is it for the agent to parse a YAML file vs. a markdown loop? Is there precedent for YAML parsing failures in the agent runtime? The pros/cons are qualitative opinions, not grounded in measurements or past project experience. The comparison would be stronger with concrete examples: "when we introduced YAML config for X, it caused Y issues." (-10 pts)

**[Dimension 3] Section "Alternatives & Industry Benchmarking" -- chosen approach not explicitly justified against the industry benchmarks.**
Three industry patterns are cited (pytest, Jest, Cucumber) and all use the Strategy pattern. The chosen approach (full dispatcher) is justified by "architectural consistency with gen-test-cases" -- an internal precedent. But the document never explicitly states why the full dispatcher matches the industry patterns better than config-driven dispatch. The link between "pytest's conftest.py model" and "full dispatcher with type files" is implied but never explicitly drawn: "our type files are analogous to pytest's conftest.py because...". The reader must infer the mapping. (-10 pts)

**[Dimension 3] Section "Alternatives & Industry Benchmarking" -- only 2 genuine alternatives beyond do-nothing (3 required).**
The three alternatives are: (1) Full dispatcher, (2) Config-driven dispatch, (3) Do nothing. The rubric requires "at least 3 meaningful alternatives" including do-nothing, meaning at least 2 genuine alternatives beyond do-nothing. Config-driven dispatch is a genuine alternative, so this meets the minimum. However, both alternatives are dispatch-based. Missing: an alternative like "hybrid approach -- keep monolith but add convention section per type" or "use profile generate.md files for type-specific logic instead of new type files." The alternatives explore only one axis (how to dispatch). (-15 pts)

**[Dimension 4] Section "Non-Functional Requirements" -- agent execution time claim lacks measurement methodology.**
Quote: "Measured impact must be under 2 seconds added latency compared to monolithic SKILL.md loading."
The 2-second threshold is arbitrary -- why 2 seconds and not 1 or 5? And the parenthetical explains that reading a 200-line + 100-line file is "strictly smaller" than 530 lines, which contradicts the need for a latency threshold (if it's strictly smaller, latency cannot increase). The NFR is internally inconsistent: if the new format is smaller, why set a threshold at all? Either the NFR should measure end-to-end agent execution time (which includes agent reasoning time, not just file loading) or it should acknowledge that file size reduction is the primary metric. (-10 pts)

**[Dimension 4] Section "Error Scenarios" -- convention conflict resolution is underspecified.**
Quote: "If both declare the same key (e.g., `test-id-attribute`), profile generate.md wins because it is profile-authoritative."
This rule is stated but no requirement specifies how the agent detects the conflict. Does the dispatcher check for key overlap before proceeding? Or is this a documentation guideline that type file authors must follow manually? If the agent doesn't enforce it, the "resolution" is aspirational, not operational. (-5 pts)

**[Dimension 5] Section "Innovation Highlights" -- the proposal explicitly claims zero innovation.**
Quote: "This is explicitly a pattern-reuse refactor, not a novel design."
The document preemptively surrenders on creativity. While honesty is appreciated, the rubric asks whether the proposal innovates *beyond* the industry baseline. The one claimed "non-obvious design decision" (reconnaissance strategy placement) is the obvious consequence of splitting by type. The pytest/Cucumber analogies are mentioned but the proposal does not push beyond them. There is no cross-domain inspiration beyond citing these testing frameworks, which are in-domain. The proposal earns credit for honesty but loses heavily on the creativity dimension itself. (-30 pts on novelty, -15 pts on cross-domain, -5 pts on simplicity of insight)

**[Dimension 5] Section "Innovation Highlights" -- cross-domain inspiration is in-domain only.**
The cited references (pytest, Jest, Cucumber) are all testing frameworks -- the same domain as gen-test-scripts. Cross-domain inspiration would borrow from a different domain (e.g., compiler passes, router middleware chains, plugin systems from editors like VS Code extensions, or game engine entity-component systems). The document does not look outside the testing world. (-15 pts on cross-domain)

**[Dimension 8] Section "Key Risks" -- only 4 risks identified, and one is trivial (Low/Low).**
Quote: "Mobile type file is speculative (no mobile profile usage yet)" with Likelihood L, Impact L.
A Low/Low risk contributes negligible analytical value. Removing it leaves only 3 meaningful risks -- the bare minimum. The proposal does not consider risks like: "Type file authors may duplicate logic across type files, recreating the monolith problem in distributed form" or "The dispatcher loop in markdown may be too abstract for the agent to follow reliably, leading to skipped type files." (-10 pts)

**[Dimension 8] Section "Key Risks" -- mitigations are partially non-actionable or describe design decisions, not recovery actions.**
Quote for Risk 2: "Type files describe *strategy* (what to look for), generate.md describes *syntax* (how to write it) -- complementary, not competing."
This is an assertion that the risk won't materialize, not a mitigation. If divergence does occur, what does the team do? A mitigation would be: "Add a CI check that greps type files for syntax-level keywords that belong in generate.md" or "During code review, verify type files contain no generate.md-overlapping keys."
Quote for Risk 3: "Hard gate in main SKILL.md: Step 3.5 always runs regardless of type."
This is a design choice baked into the dispatcher, not a response to the risk of accidental logic migration. The risk is about human error during authoring; the "mitigation" assumes the author follows the instruction perfectly. (-10 pts)

**[Dimension 9] Section "Success Criteria" -- "at most 250 lines" criterion is arbitrary without justification.**
Quote: "Main SKILL.md line count is at most 250 lines (hard ceiling; target is ~200)."
Why 250 and not 200 or 300? The "~200 lines" estimate in the Solution section is not derived from any analysis. If the dispatcher has Steps 0-1, 1.5 generic, 3.5, Step 4 dispatch loop, how many lines does each step consume? Without a breakdown, the 250-line ceiling is a round number, not a reasoned constraint. (-5 pts)

**[Dimension 9] Section "Success Criteria" -- multi-type generation scenario has no dedicated criterion.**
Quote: The "Key Scenarios" section lists "Multi-type generation (no filter): Dispatcher iterates detected types, loading each type's instruction file sequentially."
But no success criterion explicitly tests multi-type behavior. There is no criterion like: "When `--type` is not provided and the profile has both CLI and API capabilities, the dispatcher generates test scripts for both types, and both output files are present and correct." The closest is the regression gate (byte-identical output), but that does not isolate multi-type dispatch correctness. (-0 pts: covered implicitly by regression gate, but worth noting for improvement)

**[Dimension 10] Logical Consistency -- success criteria do not cover the "convention loading" fix from the Problem section.**
The Problem section identifies Inconsistency #2: "gen-test-scripts loads conventions from gen-test-cases type files rather than its own." The Solution addresses this by having type files declare their own conventions. But no success criterion explicitly tests that gen-test-scripts no longer reads gen-test-cases type files for conventions. The criterion "Each type file declares conventions in frontmatter" verifies the new behavior but does not verify the old (incorrect) behavior is gone. (-5 pts)

**[Dimension 10] Logical Consistency -- Requirements section does not explicitly require the convention inconsistency fix.**
The Problem section calls out the convention loading bug (line 50 reads gen-test-cases conventions). The Error Scenarios section covers convention conflict. But no requirement explicitly states "gen-test-scripts must load conventions from its own type files, not from gen-test-cases type files." The solution fixes it implicitly, but the requirements-to-solution mapping has an orphan problem statement. (-5 pts)

**[Dimension 7] Section "Scope" -- in-scope items mix deliverables with attributes.**
Items like "Per-type reconnaissance strategies" and "Per-type Fact Table required keys" are attributes of the type files, not independent deliverables. The in-scope list conflates "what we build" (5 type files, refactored SKILL.md) with "what each file contains" (reconnaissance, Fact Table keys, verification methods). This is a clarity issue, not a functional gap. (-5 pts)

## STRENGTHS

**[Dimension 1] Problem Definition -- evidence is concrete and verifiable.**
Line numbers, file paths, and comparative metrics (530 lines vs. ~150 lines + 5 files) are provided. A reader can open the referenced files and verify each claim. The convention loading bug is pinned to line 50. This is strong empirical grounding.

**[Dimension 2] Solution Clarity -- user-facing behavior is now explicitly described.**
The new "User-Facing Behavior" subsection describes the before/after operator experience concretely: context window consumption, error message specificity, and behavioral parity. This directly addresses the iteration-1 gap.

**[Dimension 2] Solution Clarity -- technical direction is precise.**
The dispatcher flow (Steps 0-1, 1.5, 3.5, Step 4 dispatch) with specific responsibilities per step is clear. The decision to move Steps 2-3 to ui.md is stated with rationale (UI-exclusive logic).

**[Dimension 3] Industry Benchmarking -- real industry patterns are now cited with depth.**
pytest (conftest.py), Jest (custom runners), and Cucumber (formatters) are described with specific mechanism details, not just names. The "relevance" paragraph draws an explicit analogy to the knowledge boundary concept. This is a major improvement from iteration 1.

**[Dimension 4] Requirements Completeness -- error scenarios are now comprehensive.**
Four error scenarios cover missing type file, unknown type, malformed frontmatter, and convention conflict. Each specifies the expected agent behavior (halt, reject, error message format). This fully addresses the iteration-1 gap.

**[Dimension 4] Requirements Completeness -- non-functional requirements are now present.**
Three NFRs are identified: agent execution time (2s threshold), type file loading reliability (static files, no network), and backward compatibility (byte-identical output). This addresses the iteration-1 gap where NFRs were entirely absent.

**[Dimension 4] Requirements Completeness -- the "Existing Per-Type Task Infrastructure" section provides excellent boundary context.**
The architecture overview, task generation code, dependency chain, type inference rules, and prompt synthesis flow demonstrate that the orchestration layer is complete and orthogonal. The two-layer table (orchestration vs. skill internals) is a clean separation.

**[Dimension 6] Feasibility -- the new Feasibility Assessment section is thorough.**
Technical feasibility is grounded in the proven gen-test-cases precedent. Resource estimate (1-2 days, broken down by task) is specific and realistic. Dependency readiness table names each dependency with status and risk. This fully addresses the iteration-1 gap.

**[Dimension 9] Success Criteria -- most criteria are now measurable and testable.**
Criteria use concrete verification commands (grep, diff) and hard thresholds (250 lines, 20% overlap, byte-identical). The "adding a new type file requires zero edits" criterion is an excellent extensibility test. This is a significant improvement.

**[Dimension 10] Logical Consistency -- problem-solution alignment is strong.**
Each of the three stated problems (maintenance burden, inconsistency, cognitive load) maps directly to a solution element. The "Existing Per-Type Task Infrastructure" section prevents confusion about what this proposal touches vs. what already exists.
