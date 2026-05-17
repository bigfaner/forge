# Eval-Proposal: Iteration 3

**Score: 891/1000** (target: 900)

## DIMENSIONS

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 105 | 110 |
| 2. Solution Clarity | 115 | 120 |
| 3. Industry Benchmarking | 110 | 120 |
| 4. Requirements Completeness | 105 | 110 |
| 5. Solution Creativity | 75 | 100 |
| 6. Feasibility | 95 | 100 |
| 7. Scope Definition | 75 | 80 |
| 8. Risk Assessment | 85 | 90 |
| 9. Success Criteria | 80 | 80 |
| 10. Logical Consistency | 46 | 90 |

## ATTACKS

**[Dimension 1] Section "Urgency" -- the "~80 lines per new profile" figure is an extrapolation from aggregate data, not a measured average.**
Quote: "each adding ~80 lines of type-specific branches to the monolith (extrapolated from the 530-line current state divided across 6 profiles * 5 types)"
The parenthetical reveals this is `530 / (6 * 5) ≈ 17.7 lines`, not 80. The math as written does not produce 80. Either the denominator is wrong or the 80 is unjustified. This weakens the quantified urgency claim because the growth rate appears inflated. (-5 pts)

**[Dimension 2] Section "User-Facing Behavior" -- the "After" description conflates operator experience with agent mechanics.**
Quote: "For `--type cli`, the agent loads only `types/cli.md` -- no sitemap/locator content, no UI conventions."
This describes what the agent loads (internal mechanics), not what the operator experiences. The rubric requires "what does the end user experience." The preceding sentence ("The operator sees no behavioral change") is the actual user-facing claim, but it is stated in a single sentence, while the rest of the paragraph describes agent internals. The user-facing behavior section still tilts heavily toward internals. (-5 pts)

**[Dimension 3] Section "Alternatives & Industry Benchmarking" -- the config-driven dispatch cons cite a specific incident but without verifiable detail.**
Quote: "which caused 2 misreads in the prior `config.yaml` loading step (the agent treated a commented-out line as active)"
This is a strong, specific claim. However, the incident is not traceable -- there is no commit hash, session log reference, or date. If this is the "convention loading bug" from line 50, the reader cannot confirm it. If it is a separate incident, the reader has no way to verify. The claim bolsters the rejection argument but is not independently verifiable. (-5 pts)

**[Dimension 3] Section "Alternatives & Industry Benchmarking" -- only 3 alternatives beyond do-nothing (3 required by rubric including do-nothing, so the count is met), but all explore only one axis.**
The four alternatives are: full dispatcher, config-driven dispatch, in-place monolith, do nothing. This meets the minimum count. However, all three non-trivial alternatives are "how to organize the dispatch" -- none explores "whether to dispatch at all." A genuinely different axis would be: "use profile-level overrides instead of type-level dispatch" (the profile's generate.md already has type-specific sections). The alternatives explore one design space exhaustively but do not challenge the premise that dispatch is needed. (-5 pts)

**[Dimension 4] Section "Error Scenarios" -- convention conflict resolution mechanism is still underspecified.**
Quote: "If both declare the same key (e.g., `test-id-attribute`), profile generate.md wins because it is profile-authoritative."
This is a policy, not a mechanism. Who enforces it? The agent? A validation step? The type file author? The proposal says "The type file must not duplicate syntax-level conventions" but this is a constraint on authoring, not an enforced behavior. If a type file author violates it, what happens at runtime? The error scenario describes the resolution rule but not how it is operationalized. (-5 pts)

**[Dimension 5] Section "Innovation Highlights" -- the compiler analogy is well-drawn but the actual design does not exploit it.**
Quote: "In LLVM, each target architecture (x86, ARM, RISC-V) implements a `TargetLowering` class that translates generic IR into machine-specific instructions."
The analogy is instructive, but the proposal does not specify a concrete IR schema. The Fact Table is described as the IR equivalent, but there is no IR validation step -- no schema check that ensures the Fact Table has all required keys before generation begins. In LLVM, the IR is typed and validated; in this proposal, the "IR" is free-form markdown that the agent interprets. The analogy highlights a strength of compilers (typed IR) that this proposal does not replicate. The cross-domain inspiration is real but the application is incomplete. (-15 pts on novelty, -10 pts on cross-domain)

**[Dimension 5] Section "Innovation Highlights" -- the "dispatch correctness without mechanical enforcement" claim is a constraint acknowledgment, not a creative insight.**
Quote: "This constraint -- dispatch correctness without mechanical enforcement -- is the non-obvious design challenge"
Acknowledging a constraint is necessary but not creative. The proposal does not propose a creative mechanism to ensure correctness beyond "hard-coded mapping" and "explicit error messages" and "regression gate" -- all of which are standard engineering. The "distinctive" challenge is identified but the solution to it is conventional. (-0 pts -- the simplicity-of-insight score is salvaged by the clear problem framing, but the novelty score takes the hit above)

**[Dimension 7] Section "Scope" -- in-scope items still mix deliverables with attributes of deliverables.**
Quote: Items like "Per-type reconnaissance strategies (what to grep, what patterns to search)" and "Per-type Fact Table required keys (completeness gate)" are attributes of the type files, not independent deliverables.
This was flagged in iteration 2. The list now includes both "Create 5 instruction files" (deliverable) and "Per-type reconnaissance strategies" (attribute). The items should be restructured as: (1) Create 5 type files, each containing: reconnaissance strategies, Fact Table keys, generation patterns, verification methods, antipattern guards, conventions in frontmatter. (2) Refactor main SKILL.md to dispatcher. (3) Move Steps 2-3 to types/ui.md. This remains a clarity issue. (-5 pts)

**[Dimension 8] Section "Key Risks" -- Risk 2 mitigation conflates prevention with recovery.**
Quote for Risk 2 Prevention: "Add a review checklist item verifying type files contain no `generate.md`-overlapping keys."
Quote for Risk 2 Recovery: "If divergence is detected post-merge, the `convention conflict` error scenario (profile generate.md wins) handles it at runtime"
The prevention is actionable (review checklist). But the recovery relies on the error scenario handling convention conflicts -- which, as noted in Dimension 4, has no enforcement mechanism. If the agent does not check for key overlap at runtime, the "recovery" is aspirational. (-5 pts)

**[Dimension 8] Section "Key Risks" -- no risk addresses the agent misinterpreting the dispatcher loop itself.**
Five risks are listed, but none addresses: "The agent reads the Step 4 dispatch loop but does not follow it correctly -- e.g., reads the first type file and stops, or reads all type files but applies only the last one." Risk 1 covers "loads wrong type file" and Risk 5 covers "fails to follow the loop for multi-type," but neither addresses the single-type case where the agent reads the correct file but ignores the dispatcher's instruction to delegate. This is the core risk of a markdown-based dispatcher and it is split across two risks rather than named directly. (-0 pts -- Risk 5 covers multi-type closely enough that the gap is narrow, but the single-type misinterpretation scenario is not explicitly addressed)

**[Dimension 10] Section "Alternatives & Industry Benchmarking" vs "Success Criteria" -- the in-place monolith alternative claims "zero new files" as a pro, but the success criteria require "5 type instruction files exist" without explaining why this is inherently superior.**
The logical chain is: Problem = monolith is bad → Solution = split into files → Success = 5 files exist. But the in-place monolith alternative offers to restructure with sections instead of files, and the rejection reason is "the agent still receives the full 530+ line file." This rejection is based on cognitive load, but no success criterion measures cognitive load reduction -- only line count reduction. Line count and cognitive load are correlated but not identical: a well-structured 400-line file with clear section headers may impose less cognitive load than a 200-line dispatcher + 5 files requiring cross-file navigation. The success criteria measure a proxy (line count) rather than the actual benefit claimed (cognitive load). (-10 pts on solution-addresses-problem alignment)

**[Dimension 10] Requirements-to-solution coherence -- the "Self-owned convention loading" requirement in Key Scenarios maps to a success criterion, but the requirement text specifies a path-level constraint ("SKILL.md must not reference any file path under `gen-test-cases/types/`") while the success criterion checks for the string "gen-test-cases/types" anywhere in SKILL.md.**
Quote from Requirements: "SKILL.md must not reference any file path under `gen-test-cases/types/` for convention resolution."
Quote from Success Criteria: "`grep -c "gen-test-cases/types" SKILL.md` returns 0."
These are coherent. However, the requirement also says "gen-test-scripts must load per-type conventions from its own `types/{type}.md` frontmatter" -- but no success criterion verifies that conventions are actually loaded from the new type files (only that frontmatter declares them). The loading mechanism (the dispatcher actually reading type file frontmatter) is assumed but not verified by any criterion. (-8 pts)

**[Dimension 10] Solution-to-scope coherence -- the solution says "Move Step 2-3 to `types/ui.md`" but the scope lists "Move Step 2 (Sitemap) and Step 3 (Locators) into `types/ui.md`" as in-scope, while the success criteria say "`grep -c "Sitemap" SKILL.md` returns 0."**
The success criterion checks that "Sitemap" is absent from SKILL.md but does not check that Step 3 (Locators) is absent. The string "Sitemap" is specific to Step 2 but Step 3 content may remain in SKILL.md without the grep catching it (if Step 3 uses different wording like "locator mapping" or "element locators"). The success criterion does not fully cover the in-scope item. (-6 pts)

**[Dimension 10] Scope-to-success-criteria alignment -- the scope lists "Per-type verification methods" and "Per-type antipattern guards" as in-scope, and the success criteria require "Each type file contains a verification method section and an antipattern guard section." However, these criteria only check presence (10 matches total across 5 files), not correctness or completeness.**
Quote: "Each type file contains a verification method section and an antipattern guard section (verifiable: 10 matches total across 5 files)"
A type file could have an empty or trivial "## Verification Methods" section with no content and pass this criterion. Presence is verified but quality is not. For in-scope items that are core deliverables (not attributes), the success criteria should include a quality check -- e.g., "each verification method section lists at least 2 methods specific to the type." (-10 pts)

**[Dimension 10] Problem-to-solution -- the problem identifies "cognitive load" as issue #3 but the solution provides no cognitive load measurement.**
The three problems are: (1) maintenance burden, (2) inconsistency, (3) cognitive load. Solutions address all three. Success criteria address (1) via line count and (2) via convention coupling. But no success criterion addresses (3). Cognitive load is the subjective experience of the agent navigating instructions. The line count reduction is a proxy, but the problem statement describes cognitive load as "Agents must navigate 530 lines of mixed generic and type-specific instructions, increasing the chance of applying wrong patterns to wrong types." The proxy (line count) measures navigation burden but not "applying wrong patterns to wrong types." A better criterion would be: "for each type, the generated test scripts contain zero patterns from other types" -- which would directly measure the stated cognitive load consequence. (-10 pts)

## STRENGTHS

**[Dimension 1] Problem Definition -- evidence is concrete, verifiable, and now quantified.**
Line numbers (530 lines, line 50), file paths (gen-test-scripts SKILL.md), comparative metrics (530 vs. ~150 + 5 files), and now incident data (3-session undetected bug, ~1 hour debugging, 3 recent commits each touching 2-3 unrelated sections). A reader can verify every claim by opening the referenced files and commit history.

**[Dimension 1] Problem Definition -- urgency is now quantified with three concrete cost-of-delay vectors.**
Growth pressure (6 profiles, 2 more planned), bug incidence (3-session undetected bug, 1-hour debug cost), and per-change cost (3 recent commits each touched 2-3 unrelated sections). This is a major improvement from iterations 1-2.

**[Dimension 2] Solution Clarity -- user-facing behavior is now explicitly described with before/after comparison.**
The "User-Facing Behavior" subsection provides a concrete before/after scenario with specific invocation (`forge prompt get-by-task-id T-test-2-cli`), describing context window consumption, error message specificity, and behavioral parity.

**[Dimension 2] Solution Clarity -- technical direction is precise and actionable.**
The dispatcher flow (Steps 0-1, 1.5, 3.5, Step 4) with specific responsibilities per step is clear. The decision to move Steps 2-3 to ui.md is stated with rationale (UI-exclusive logic). The line budget breakdown in Success Criteria provides a verifiable technical target.

**[Dimension 3] Industry Benchmarking -- real industry patterns are cited with depth and specificity.**
pytest (conftest.py + pytest_plugins), Jest (custom runners), and Cucumber (formatters) are described with specific mechanism details. The "Benchmark Mapping to Chosen Approach" section draws an explicit point-by-point analogy to pytest's model, addressing the iteration-2 gap.

**[Dimension 3] Industry Benchmarking -- the config-driven dispatch alternative now cites a specific internal incident.**
The con about config-driven dispatch references "2 misreads in the prior `config.yaml` loading step" with a concrete failure mode (agent treated a commented-out line as active). This grounds the trade-off in project experience.

**[Dimension 3] Industry Benchmarking -- the in-place monolith alternative is now a genuine option with real trade-offs.**
Unlike iteration 1's straw-man "additive type files," the in-place monolith with per-type sections is a legitimate alternative that addresses maintenance burden partially, with honest cons (cognitive load persists, convention isolation fails).

**[Dimension 4] Requirements Completeness -- error scenarios are comprehensive with specific agent behaviors.**
Four error scenarios cover missing type file, unknown type, malformed frontmatter, and convention conflict. Each specifies expected agent behavior (halt, reject, error message format, resolution rule).

**[Dimension 4] Requirements Completeness -- non-functional requirements are well-grounded.**
The NFR section addresses context budget (with line-count proxy), file loading reliability (static files, no network), and backward compatibility (byte-identical output). The decision to drop the 2-second latency threshold (from iteration 2) in favor of the line-count proxy is well-reasoned.

**[Dimension 4] Requirements Completeness -- constraints and dependencies are explicitly listed.**
Backward compatibility with 6 profiles, agent runtime file loading, Step 3.5 remaining shared, and the dependency readiness table with named dependencies and status.

**[Dimension 5] Solution Creativity -- the compiler pass analogy is genuinely cross-domain and well-developed.**
The LLVM TargetLowering analogy draws from compiler design (a different domain from testing) and maps concretely to the proposal's architecture: Fact Table = IR, type files = target backends, dispatcher = compiler driver. The analogy is specific and educational.

**[Dimension 5] Solution Creativity -- the "dispatch correctness without mechanical enforcement" framing is distinctive.**
Identifying the core challenge as "a language model must follow dispatch rules from free-form markdown without a parser" is a genuinely non-obvious insight. This constraint shapes multiple design decisions (hard-coded mapping, explicit errors, regression gate) in a coherent way.

**[Dimension 6] Feasibility -- technical feasibility is grounded in a proven internal precedent.**
The gen-test-cases refactor (5 type files in production, zero dispatcher bugs) provides direct evidence that the pattern works in the same runtime environment.

**[Dimension 6] Feasibility -- resource and timeline estimate is specific and realistic.**
1-2 days, broken down by task (10 hours extraction, 3 hours dispatcher, 2 hours verification). Skills required are named (understanding of all 5 test types, familiarity with gen-test-cases dispatcher).

**[Dimension 8] Risk Assessment -- risks are now meaningful and well-structured.**
Five risks, all with M or H likelihood or impact. No trivial Low/Low risks. The new "Logic duplication across type files recreates monolith in distributed form" risk is particularly insightful -- it addresses a failure mode that would undermine the entire refactoring.

**[Dimension 8] Risk Assessment -- mitigations now include both prevention and recovery actions.**
Each risk has a two-tier mitigation: prevention (design choice or review process) and recovery (what to do if the risk materializes). This directly addresses the iteration-2 criticism that mitigations were assertions rather than actions.

**[Dimension 9] Success Criteria -- criteria are now comprehensive and almost universally measurable.**
12 criteria use concrete verification commands (grep, diff), hard thresholds (250 lines with line budget), and testable scenarios (--type graphql, missing file). The "adding a new type file requires zero edits" criterion is an excellent extensibility test. The "byte-identical before/after" regression gate is the gold standard for refactoring verification.

**[Dimension 9] Success Criteria -- the line budget breakdown is a major improvement.**
The 250-line ceiling is now derived from a per-step budget totaling ~225 lines with 25 lines of margin. This makes the criterion verifiable and the ceiling justified rather than arbitrary.

**[Dimension 9] Success Criteria -- convention coupling bug is explicitly verified.**
The criterion "grep -c `gen-test-cases/types` SKILL.md returns 0" directly tests the elimination of the convention loading bug identified in the Problem section. This closes the problem-to-criterion loop.
