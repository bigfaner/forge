# Eval-Proposal: Iteration 1

**Score: 732/1000** (target: 900)

## DIMENSIONS

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 92 | 110 |
| 2. Solution Clarity | 82 | 120 |
| 3. Industry Benchmarking | 38 | 120 |
| 4. Requirements Completeness | 78 | 110 |
| 5. Solution Creativity | 50 | 100 |
| 6. Feasibility | 78 | 100 |
| 7. Scope Definition | 68 | 80 |
| 8. Risk Assessment | 68 | 90 |
| 9. Success Criteria | 60 | 80 |
| 10. Logical Consistency | 78 | 90 |

## ATTACKS

1. **[Dim 3 - Industry Benchmarking] Missing industry references entirely.**
   The "Alternatives" table (lines 152-159) lists only internal approaches (W1-only, W2-only, rewrite-from-scratch). There are zero references to how other plugin ecosystems (VS Code extensions, JetBrains plugins, npm packages, Claude Code community plugins) handle path resolution, asset distribution, or deduplication of template files. The section needs at least 3 external references with honest trade-off analysis.
   - *Section*: "Alternatives" (lines 152-159)
   - *Fix*: Research and cite how at least 3 other plugin/extension systems solve path resolution and template deduplication. Add a proper "Industry Benchmarking" section before "Alternatives."

2. **[Dim 3 - Industry Benchmarking] Straw-man alternatives.**
   The "Rewrite from scratch" row is a straw man -- it is included solely to make the selected approach look good by comparison. No reasonable person would propose rewriting all 18 skills to fix path references and dead code.
   - *Section*: "Alternatives" table, row 4 (line 159)
   - *Fix*: Replace with a realistic alternative (e.g., "Symlink-based resolution," "Build-time path rewriting," "Convention-based relative paths"). Each alternative should be a genuinely defensible choice.

3. **[Dim 1 - Problem / Urgency] Urgency is implied but not justified.**
   The document catalogs bugs and issues but never explains *why now*. Is there a release imminent? Are users actively hitting these failures? The urgency justification for "why this audit must happen before other work" is absent.
   - *Section*: "Problem" (lines 9-17)
   - *Fix*: Add a paragraph explaining the triggering event or deadline (e.g., upcoming v3.0.0 release, user reports of failures, plugin store submission).

4. **[Dim 2 - Solution Clarity] Path resolution approach is undecided.**
   W2 lists three options for fixing `plugins/forge/` paths (lines 117-119): `${CLAUDE_PLUGIN_ROOT}`, relative paths, or a `<FORGE_ROOT>` placeholder. The proposal does not commit to one. This is the single most critical technical decision in the entire document, and it is left as an open question.
   - *Section*: "W2: Distributability Fixes" (lines 113-123)
   - *Fix*: Evaluate each option, select one, and justify the selection. Remove the "needs verification" hedge.

5. **[Dim 2 - Solution / User-facing behavior] User-facing behavior is incomplete.**
   The "Key Scenarios" section (lines 139-144) lists 5 scenarios but skips the most common ones: user runs `/brainstorm`, `/write-prd`, `/tech-design`, `/breakdown-tasks` -- do they all still work? What does the user see if a path resolution fails? There is no description of error behavior or degraded-mode experience.
   - *Section*: "Key Scenarios" (lines 139-144)
   - *Fix*: Add the happy-path for at least 3 major pipeline skills. Add an error-path scenario: what happens when a path reference cannot be resolved?

6. **[Dim 4 - Requirements / Non-functional] Non-functional requirements are absent.**
   There is no mention of performance constraints (should plugin install size be under X MB?), no security considerations (are hardcoded paths a security surface?), no accessibility or localization requirements, and no backward-compatibility testing strategy.
   - *Section*: "Requirements Analysis" (lines 135-150)
   - *Fix*: Add a "Non-functional Requirements" subsection covering: install size target, backward compatibility verification method, security implications of path conventions, and any performance constraints.

7. **[Dim 4 - Requirements / Scenario coverage] Missing failure scenarios.**
   The scenarios cover happy-path behaviors but not failure modes: what if `${CLAUDE_PLUGIN_ROOT}` is undefined? What if the user's project has no `docs/conventions/`? What if two skills share a template and one modifies it?
   - *Section*: "Key Scenarios" (lines 139-144)
   - *Fix*: Add at least 2 failure/recovery scenarios.

8. **[Dim 5 - Solution Creativity] No novelty beyond straightforward cleanup.**
   The proposal is essentially: delete dead files, fix paths, extract shared templates. There is no creative insight -- no proposal for a generic skill asset resolution protocol, no template pre-processor, no build-time validation that catches path regressions automatically. The "simplification of insight" criterion is partially met (the workstream prioritization is clean), but there is no novelty or cross-domain inspiration.
   - *Section*: "Proposed Solution" (lines 101-133)
   - *Fix*: Propose at least one preventive mechanism (e.g., a CI check that scans for hardcoded `plugins/forge/` paths, a linter for skill assets, or a canonical path-registry pattern).

9. **[Dim 6 - Feasibility / Resource] No effort estimates or timeline.**
   The workstreams are labeled "Low/Medium effort" but no concrete hour-count, day-count, or sprint allocation is provided. W2 is "Medium effort, Critical correctness" -- medium could mean 2 days or 2 weeks.
   - *Section*: "Proposed Solution" (lines 101-133)
   - *Fix*: Add estimated hours or days per workstream. State who will do the work and when.

10. **[Dim 6 - Feasibility / Dependency readiness] Path resolution mechanism is unverified.**
    The document itself says `${CLAUDE_PLUGIN_ROOT}` availability "needs verification" (line 117). The entire W2 workstream depends on choosing a path resolution strategy, and the feasibility of each option is unknown.
    - *Section*: "W2: Distributability Fixes" (line 117)
    - *Fix*: Before this proposal can be approved, someone must verify whether skills have access to `${CLAUDE_PLUGIN_ROOT}` (or equivalent). Add a prerequisite spike or investigation task.

11. **[Dim 7 - Scope / Scope boundedness] No sequencing or acceptance criteria per workstream.**
    The scope lists 9 items as in-scope but does not indicate whether all 3 workstreams must be completed in one release or can be phased. Without phasing, the scope is unbounded in time.
    - *Section*: "Scope" (lines 161-184)
    - *Fix*: State the phasing strategy: W1 in release X, W2 in release Y, W3 in release Z (or all in one). Add a "not doing in this iteration" line.

12. **[Dim 8 - Risk / Mitigations are vague] Mitigations lack specificity.**
    Risk mitigations say things like "Test with multiple Claude Code versions; provide fallback" -- which fallback? What does the fallback look like? "Verify gen-test-scripts output is identical before/after" -- using what verification method? Manual inspection? Automated diff? The mitigations are aspirations, not actionable steps.
    - *Section*: "Key Risks" (lines 188-193)
    - *Fix*: Make each mitigation concrete. For example: "Write a bash script that invokes gen-test-scripts with a Playwright profile and diffs output against a golden fixture." Name the specific test or tool.

13. **[Dim 8 - Risk / Missing risks] No regression risk from shared template extraction.**
    When breakdown-tasks and quick-tasks share templates via reference files, a change to the shared template affects both skills. This tight coupling risk is not listed.
    - *Section*: "Key Risks" (lines 188-193)
    - *Fix*: Add a risk for "Shared template changes have unintended cross-skill effects" with likelihood/impact/mitigation.

14. **[Dim 9 - Success Criteria] Criteria are not all testable.**
    "breakdown-tasks and quick-tasks share common logic via reference files (not copy-paste)" -- how do you test "not copy-paste"? Is there a metric? "All slash commands produce identical behavior" -- identical to what baseline? There is no mention of a regression test suite or golden fixtures.
    - *Section*: "Success Criteria" (lines 197-205)
    - *Fix*: Make each criterion falsifiable. For example: "grep -r 'record-task' in agents/, skills/, templates/ returns 0 hits." "diff <(gen-test-scripts --profile web-playwright) <(golden-output/) is empty."

15. **[Dim 9 - Success Criteria] Missing coverage for W3.**
    Success criteria cover W1 (items 1-2), W2 (items 3-6), but W3 (dedup and eval simplification) only has one criterion (item 7). There is no criterion for "eval validate-ux extracted" or "scope heuristic fixed for backend projects."
    - *Section*: "Success Criteria" (lines 197-205)
    - *Fix*: Add: "eval/SKILL.md is under N lines (was 366)" and "breakdown-tasks scope heuristic correctly classifies Go/Rust `src/` as backend."

16. **[Dim 10 - Logical Consistency] Solution W3 partially contradicts scope exclusion.**
    Scope excludes "New skill creation" (line 181), but W3 says "Extract 90-line validate-ux inline sub-pipeline to a separate skill or rubric" (line 132). Creating a new skill is out of scope but W3 proposes exactly that as one option.
    - *Section*: "W3" (line 132) vs "Out of Scope" (line 181)
    - *Fix*: Either remove "separate skill" as an option, or adjust the out-of-scope list to say "New skill creation (except validate-ux extraction)."

17. **[Dim 10 - Requirements ↔ Solution] Constraints mention backward compatibility but solution does not address testing it.**
    Constraint: "All slash commands must work identically after changes" (line 148). Solution: no mention of a test plan, regression suite, or verification protocol. Success criteria say "identical behavior" but the proposal provides no mechanism to verify this.
    - *Section*: "Constraints & Dependencies" (line 148) and "Proposed Solution"
    - *Fix*: Add a verification plan section describing how backward compatibility will be confirmed (e.g., existing e2e tests, manual walkthrough of each skill, golden-output comparison).

## DETAILED SCORING

### 1. Problem Definition (92/110)

- **Problem stated clearly (35/40)**: The problem is well-articulated with the distribution model insight. The per-skill scorecard is an excellent evidence-based framing. The core thesis -- "paths resolve to user project root instead of plugin cache" -- is precise. Deduction: the problem statement mixes two distinct issues (path resolution and dead code) without clearly prioritizing which is the primary problem and which is secondary.
- **Evidence provided (38/40)**: The per-skill scorecard with lines, commits, clarity ratings, and path counts is strong quantitative evidence. The full path inventory (lines 73-98) is thorough. Minor deduction: some clarity ratings (e.g., "4/5") lack explanation of the scoring criteria.
- **Urgency justified (19/30)**: The P0/P1/P2 prioritization implies urgency but never states *why this must happen now*. No release deadline, no user complaints cited, no plugin store submission blocking. The urgency is structural (broken paths) but the timing urgency is unstated.

### 2. Solution Clarity (82/120)

- **Approach is concrete (28/40)**: The three-workstream structure is clear and well-prioritized. W1 is fully concrete. W2 is partially concrete but the critical path resolution decision is left open with three options and no selection. W3 is the most vague -- "extract shared logic" and "simplify eval" describe goals, not approaches.
- **User-facing behavior described (30/45)**: The Key Scenarios section provides 5 scenarios but they are thin. Missing: what does the user see on success? What happens on failure? How do error messages change? No description of the before/after user experience.
- **Technical direction clear (24/35)**: W1 is technically precise (which file, which line, what change). W2's technical direction is ambiguous -- the path resolution question is the crux and it is unresolved. W3's extraction targets are named but the mechanism (shared reference files) is not detailed.

### 3. Industry Benchmarking (38/120)

- **Industry solutions referenced (5/40)**: Zero external references. No mention of how VS Code extensions, JetBrains plugins, Homebrew formulae, or any other plugin ecosystem handles asset path resolution or template deduplication. The document operates entirely within its own context.
- **At least 3 meaningful alternatives (8/30)**: The Alternatives table lists 4 options but 3 are internal variants (W1-only, W2-only, full rewrite). Only the selected approach represents a meaningful alternative to itself.
- **Honest trade-off comparison (15/25)**: The comparison is honest but shallow. "Medium total effort" without quantification. "High risk" for rewrite without elaboration. The selected approach has no cons listed beyond "Medium total effort."
- **Chosen approach justified against benchmarks (10/25)**: Justification is minimal -- the three-workstream approach is reasonable but not benchmarked against anything. The decision is internally consistent but not externally validated.

### 4. Requirements Completeness (78/110)

- **Scenario coverage (30/40)**: Five scenarios cover the main use cases but omit failure scenarios, edge cases (user with no `docs/conventions/`, user with custom test profile), and the full pipeline happy path.
- **Non-functional requirements (18/40)**: Backward compatibility is mentioned as a constraint. The 27MB install size reduction is implicit. But there are no explicit NFRs for: target install size, performance, security, accessibility, or reliability. The constraint about "backward compatibility" is stated but not operationalized.
- **Constraints & dependencies (30/30)**: Constraints are well-stated: backward compatibility, plugin distribution model, test profile system, and the existing skill-rationalization proposal are all identified. Dependencies are realistic.

### 5. Solution Creativity (50/100)

- **Novelty over industry baseline (15/40)**: The proposal is straightforward maintenance work: delete dead files, fix paths, extract shared templates. There is no novel mechanism, pattern, or abstraction proposed. No path resolution protocol, no template pre-processor, no CI guard.
- **Cross-domain inspiration (10/35)**: No evidence of borrowing ideas from other domains. The per-skill scorecard format is nice but is a presentation technique, not a creative solution component.
- **Simplicity of insight (25/25)**: The three-workstream prioritization is clean and practical. The recognition that W1 (cleanup) should precede W2 (correctness) which should precede W3 (dedup) shows good judgment. The insight that `plugins/forge/` paths resolve to the wrong location is a genuine and valuable finding.

### 6. Feasibility (78/100)

- **Technical feasibility (30/40)**: W1 is clearly feasible. W2's feasibility depends on the unresolved path resolution question. W3 (shared reference files) is feasible but the proposal does not describe the file format, location, or loading mechanism for shared references.
- **Resource & timeline feasibility (18/30)**: No effort estimates, no timeline, no resource allocation. "Low/Medium effort" labels are subjective and unquantified.
- **Dependency readiness (30/30)**: Dependencies are ready: the test profile system exists, the plugin distribution mechanism is understood, and no external tools or libraries are needed.

### 7. Scope Definition (68/80)

- **In-scope items are concrete (28/30)**: All 9 in-scope items are specific and actionable. Each can be mapped to files and changes.
- **Out-of-scope explicitly listed (23/25)**: 10 out-of-scope items are listed. This is thorough. Minor issue: "Naming changes (`gen-` prefix)" is listed as out-of-scope but finding #10 flagged it -- this creates ambiguity about whether the finding is accepted or deferred.
- **Scope is bounded (17/25)**: The scope is bounded in *what* but not in *when*. There is no phasing plan, no "v1 will include W1+W2, v2 will add W3" statement. The scope could expand indefinitely within the 9 items.

### 8. Risk Assessment (68/90)

- **Risks identified (22/30)**: Four risks are identified. Missing: regression risk from shared templates, risk of path resolution inconsistency across OS platforms (Windows vs macOS), risk of the eval extraction changing scoring behavior.
- **Likelihood + impact rated (22/30)**: Ratings are present (M/H, M/H, M/H, L/M) but the scale is not defined. What does "Medium likelihood" mean? 30%? 50%? No calibration is provided.
- **Mitigations are actionable (24/30)**: Mitigations are directional but not fully actionable. "Test with multiple Claude Code versions" -- which versions? "Verify gen-test-scripts output is identical" -- using what method? The mitigations need to be specific enough that someone can execute them without further clarification.

### 9. Success Criteria (60/80)

- **Criteria are measurable and testable (38/55)**: Some criteria are measurable ("zero stale references," "27MB reduction"). Others are not testable as stated: "share common logic via reference files (not copy-paste)" -- what test confirms this? "All slash commands produce identical behavior" -- identical by what measurement? No regression test suite or golden fixture is referenced.
- **Coverage is complete (22/25)**: Coverage is strong for W1 and W2 but thin for W3. The eval simplification and scope-heuristic fix have no dedicated success criteria. The "Plugin installs and all skills work" criterion is a catch-all but not decomposed.

### 10. Logical Consistency (78/90)

- **Solution addresses the stated problem (33/35)**: The three workstreams directly address the P0/P1/P2 findings. The mapping from problem to solution is tight and traceable.
- **Scope <-> Solution <-> Success Criteria aligned (23/30)**: Mostly aligned, but W3's eval extraction appears in scope and solution but has weak success criteria coverage. The out-of-scope exclusion of "new skill creation" partially contradicts W3's option to extract validate-ux "to a separate skill."
- **Requirements <-> Solution coherent (22/25)**: Requirements and solution are coherent. The backward compatibility requirement is well-understood but the solution does not describe a testing protocol to verify it. The dependency on the test profile system is correctly identified and W2's template cleanup leverages it appropriately.
