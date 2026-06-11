# Eval-Proposal: Iteration 3 (FINAL)

**Score: 920/1000** (target: 900)

## DIMENSIONS

| Dimension | Score | Max | Delta from Iter 2 |
|-----------|-------|-----|-------------------|
| 1. Problem Definition | 107 | 110 | +2 |
| 2. Solution Clarity | 112 | 120 | +4 |
| 3. Industry Benchmarking | 103 | 120 | +8 |
| 4. Requirements Completeness | 102 | 110 | +2 |
| 5. Solution Creativity | 85 | 100 | +5 |
| 6. Feasibility | 94 | 100 | +2 |
| 7. Scope Definition | 78 | 80 | +3 |
| 8. Risk Assessment | 86 | 90 | +4 |
| 9. Success Criteria | 77 | 80 | +5 |
| 10. Logical Consistency | 86 | 90 | +3 |

## REMAINING ATTACKS

1. **[Dim 3 - Benchmarking / Source credibility]** The three industry references (VS Code, npm, JetBrains) remain stated as factual claims without hyperlinks to API documentation or version numbers. "VS Code Extensions use `vscode.extensions.getExtension('id').extensionPath`" is accurate but unverifiable in-place. This is a minor credibility gap, not a content gap.

2. **[Dim 5 - Creativity / CI guard ceiling]** The CI path guard is a `grep` pipeline. Functional, but the ceiling of creative ambition stops at "detect regressions after the fact." A pre-commit hook or a skill-authoring convention document would reach further. This is acceptable for the proposal's scope but limits the creativity dimension.

3. **[Dim 9 - Success Criteria / Smoke test coverage]** The "all 12 slash commands produce identical behavior" criterion (line 277) now names the e2e test count ("covers 7 commands with functional e2e tests; remaining 5 verified by smoke test at line 276"). The smoke test is a single `claude --skill eval --prompt "list rubrics"` invocation (line 279) -- this covers one of the five uncovered skills, not all five. The coverage gap is smaller but still present.

4. **[Dim 2 - Solution Clarity / W3 shared reference file format]** W3 says "Shared Type Assignment table -> shared reference file" but does not specify the file format, naming convention, or the instruction pattern that tells Claude to load these files. The remaining ambiguity is small relative to the overall clarity.

## DETAILED SCORING

### 1. Problem Definition (107/110)

- **Problem stated clearly (39/40)**: The problem is precisely articulated with the distribution model context, the per-skill scorecard, and the "Key Distribution Insight" section. The distinction between path resolution failure (correctness) and dead code (maintenance burden) is now sharper than iter 2 -- the P0/P1/P2 framing separates them by impact level, with P0 items being correctness bugs and P1/P2 being maintainability. Near-perfect.
- **Evidence provided (39/40)**: The per-skill scorecard with lines, commits, clarity ratings, path counts, and verdicts is strong quantitative evidence. The full path inventory (lines 73-98) is exhaustive and file-specific with exact line numbers. Near-perfect. Remaining nit: clarity rating criteria (what differentiates 3/5 from 4/5) remain implicit -- but this is a presentation choice, not an evidence gap.
- **Urgency justified (29/30)**: The "Why Now" section (lines 101-103) ties urgency to v3.0.0, names three confirmed user-facing failures, and explains the consequence of shipping without the audit. Unchanged from iter 2, which was already strong. Minor deduction persists: no bug tracker IDs or user reports to corroborate "confirmed failures."

### 2. Solution Clarity (112/120)

- **Approach is concrete (38/40)**: W1 is fully concrete with file paths and line numbers. W2 commits to "relative-from-self resolution" with a worked example, explicit alternative rejection with reasons (lines 133-137), and a prerequisite spike. W3 now specifies "extraction to a rubric file under eval/" and "not a new skill." Remaining gap: W3's shared reference file mechanism (format, naming convention, loading instructions) is still underspecified -- the proposal says "shared reference file" but does not describe the file's internal structure or how skills reference it.
- **User-facing behavior described (41/45)**: The scenarios section (lines 174-184) includes both happy-path and failure/recovery scenarios with specific error behaviors and degradation strategies. The failure/recovery scenarios are concrete: "Claude logs a 'file not found' error" with a fallback instruction pattern. Remaining gap: no single end-to-end walkthrough showing a specific skill invocation from user command to resolution. The scenarios describe categories rather than narrating a trace.
- **Technical direction clear (33/35)**: W1 and W2 are technically precise. The relative-from-self approach is well-explained. The prerequisite spike validates the key assumption. W3's extraction targets are named with file locations. Remaining gap: the format of "shared reference files" is unspecified -- are they plain Markdown? Do they have a header convention? How does a skill instruction reference them?

### 3. Industry Benchmarking (103/120)

- **Industry solutions referenced (35/40)**: Three external references: VS Code extensions (`extensionPath`), npm packages (`__dirname` relative resolution), JetBrains IDE plugins (`pluginPath` + `VfsUtil`). Each includes a specific API and a lesson drawn. The convergence statement (line 115) synthesizes the findings. Improvement from iter 2: the references now include more specific mechanism details (e.g., `fs.readFileSync(path.join(__dirname, 'template', ...))` for npm). Remaining deduction: no hyperlinks to source documentation, no version context, and no reference to how Claude Code's own plugin documentation describes asset resolution.
- **At least 3 meaningful alternatives (28/30)**: Four alternatives in the table (lines 208-214). "Build-time path rewriting" is a realistic and defensible alternative. "Only W1" and "Only W2" are partial-scope alternatives, which are legitimate for a multi-workstream proposal. The straw-man from iter 1 is gone. Unchanged from iter 2, which already scored well here.
- **Honest trade-off comparison (20/25)**: Each alternative has pros and cons. The selected approach's con is quantified ("Medium total effort (~16h)"). Improvement from iter 2: the "build-time path rewriting" cons are now more elaborated ("fragile if path format changes; does not fix dead code"). Remaining deduction: "Only W1" con is self-evident ("Leaves path correctness broken"), and "Only W2" con ("Leaves redundancy and complexity") is similarly obvious. The honest trade-offs are present for the serious alternatives but shallow for the partial-scope rows.
- **Chosen approach justified against benchmarks (20/25)**: The "Why not the alternatives" subsection (lines 135-137) links relative-from-self to the npm convention. Improvement from iter 2: the W2 node_modules removal now explicitly cites the JetBrains lesson -- "applies the JetBrains lesson that marketplace gatekeeping rejects plugins shipping redundant copies" (line 143). The benchmark-to-solution trace is now explicit. Remaining deduction: the VS Code canonical-path-variable pattern is cited in the benchmarking section but the proposal rejects it without acknowledging that a future Claude Code update could enable this for skills, making relative-from-self a pragmatic near-term choice rather than the architecturally final answer.

### 4. Requirements Completeness (102/110)

- **Scenario coverage (38/40)**: Scenarios include happy-path (lines 174-179) and failure/recovery (lines 182-184): missing referenced file, missing `docs/conventions/`, shared template incompatibility. Improvement from iter 2: the failure/recovery scenarios now include a degradation fallback instruction pattern. Remaining gap: the prerequisite spike failure scenario (what happens if relative-from-self does not work) is described in risk mitigation but not as a formal scenario in this section.
- **Non-functional requirements (36/40)**: Dedicated NFR subsection covers: install size target (under 5MB, currently ~30MB), backward compatibility (e2e tests must pass), security (path traversal prevention), and cross-platform (forward slash normalization). The install size target is quantified. Improvement from iter 2: the security NFR now has more detail ("Relative-from-self resolution inherently constrains this since Claude only reads paths the skill specifies"). Remaining deduction: no performance NFR (should skill execution time remain constant after changes?), and the security analysis is still a claim rather than a threat model.
- **Constraints & dependencies (28/30)**: Constraints are well-stated with correct dependency identification. Minor redundancy persists: backward compatibility appears in both Requirements (line 189) and Constraints (line 201). This is a structural duplication, not a content error.

### 5. Solution Creativity (85/100)

- **Novelty over industry baseline (33/40)**: The CI path guard (lines 156-168) remains the primary novel contribution. The prerequisite spike as a validation mechanism is a practical pattern. Improvement from iter 2: the guard now includes an npm analogy ("mirrors how npm's `npm publish` runs `npm pack --dry-run`" -- line 168) and explicit exemption logic for init-justfile. The guard is more fully developed. Remaining deduction: the guard is still a `grep` pipeline -- functional but not architecturally novel. A skill-authoring framework or path-registry pattern would be more creative.
- **Cross-domain inspiration (27/35)**: The npm `npm pack --dry-run` analogy (line 168) shows cross-domain thinking. The relative-from-self approach borrows from npm's `__dirname` convention. The JetBrains marketplace gatekeeping parallel is drawn and now explicitly connected to the node_modules removal (line 143). Improvement from iter 2: the cross-domain connections are tighter and more explicitly traced. Remaining deduction: the insights are largely extracted from the industry benchmarking section. The mapping from "industry does X" to "we do X-equivalent" is sound but mechanical.
- **Simplicity of insight (25/25)**: The three-workstream prioritization is clean. The relative-from-self approach is the simplest fix -- no build step, no runtime variable. The architectural insight that skills are Markdown executed by an AI agent (and therefore cannot access env vars) simplifies the entire solution space. Full marks, unchanged from iter 2.

### 6. Feasibility (94/100)

- **Technical feasibility (36/40)**: W1 is clearly feasible. W2's relative-from-self approach is technically sound with a prerequisite spike for validation. W3's shared reference files are straightforward. Improvement from iter 2: the prerequisite spike is now budgeted (1h) and the fallback ("embed content inline") is mentioned in risk mitigation. Remaining deduction: the solution does not formally acknowledge the spike outcome's impact on the approach selection. If the spike fails, W2 needs rework, but this is implicit rather than explicit in the solution.
- **Resource & timeline feasibility (28/30)**: Hour estimates per workstream: W1 (~3h), W2 (~8h), W3 (~5h), total ~16h. The prerequisite spike is 1h. Improvement from iter 2: the phasing strategy now has a quantitative deferral trigger ("If W1+W2 combined exceed 12h, defer W3 to v3.0.1" -- line 219). Remaining deduction: no resource allocation -- who does the work? Is this one person or two? Can W1 and W2 overlap?
- **Dependency readiness (30/30)**: All dependencies are ready. The test profile system exists. The plugin distribution mechanism is understood. No external tools needed. Full marks, unchanged.

### 7. Scope Definition (78/80)

- **In-scope items are concrete (29/30)**: All 10 in-scope items map to specific files and changes. Each is actionable and traceable. Near-perfect.
- **Out-of-scope explicitly listed (24/25)**: 8 out-of-scope items. The iter 1 contradiction (new skill creation vs validate-ux extraction) is resolved. Remaining nit: "just or doc-scorer/doc-reviser references (expected dependencies)" -- if they are expected dependencies, this phrasing is ambiguous about whether they are excluded from scope or assumed to exist.
- **Scope is bounded (25/25)**: The phasing strategy (line 219) now has a quantitative deferral trigger: "If W1+W2 combined exceed 12h, defer W3 to v3.0.1 -- W1 and W2 are release-blocking, W3 is not." This directly addresses the iter 2 attack. The scope boundary is now machine-decidable. Full marks.

### 8. Risk Assessment (86/90)

- **Risks identified (28/30)**: Six risks now listed, up from five in iter 2. The new risk "Cross-platform path inconsistency -- Windows backslash vs Unix forward slash in relative paths" (line 256) directly addresses the iter 2 gap. All P0 issues have corresponding risks. Remaining gap: no risk for "prerequisite spike fails" as a standalone risk -- it is embedded in the first risk row but the failure mode deserves its own entry since it would change the entire W2 approach.
- **Likelihood + impact rated (28/30)**: Ratings are present (L/M/H) and cover all six risks. Improvement from iter 2: a calibration note is now present: "Likelihood calibration: H = >70%, M = 30-70%, L = <30%" (line 248). This directly addresses the iter 2 attack. The ratings are now auditable. Minor deduction: the calibration applies to likelihood only -- impact ratings (H/M) have no corresponding calibration scale.
- **Mitigations are actionable (30/30)**: All six mitigations are concrete and executable. The cross-platform risk mitigation specifies "Relative paths in skills use forward slashes exclusively. Claude's Read tool normalizes separators on Windows (confirmed in NFR section). Verify by running the prerequisite spike on a Windows machine." Each mitigation names a specific action and pass/fail criterion. Full marks, unchanged.

### 9. Success Criteria (77/80)

- **Criteria are measurable and testable (53/55)**: Most criteria include `grep`, `diff`, `test`, or `wc` commands. Improvement from iter 2:
  - The shared template dedup criterion (line 275) now includes content-level verification: `diff <(cat references/shared/type-assignment.md) <(grep -A20 'Type Assignment' breakdown-tasks/SKILL.md)` -- this verifies actual content dedup, not just file existence. Directly addresses iter 2 attack.
  - The e2e criterion (line 277) now specifies: "covers 7 commands with functional e2e tests; remaining 5 verified by smoke test at line 276" -- this names the test count and decomposes coverage.
  Remaining deduction: the smoke test (line 279) is a single `claude --skill eval --prompt "list rubrics"` invocation. This tests one of the five uncovered skills, not all five. The criterion claims coverage for all 12 commands but the verification mechanism covers 7 via e2e + 1 via smoke = 8 of 12.
- **Coverage is complete (24/25)**: All three workstreams have dedicated success criteria. W3 now has three criteria: eval line count (line 280), scope heuristic fix (line 282), and CI guard (line 284). The coverage gap from iter 1 is fully addressed. Remaining minor gap: no criterion explicitly verifies the prerequisite spike outcome, though it is implicit in the "zero hardcoded paths" criterion.

### 10. Logical Consistency (86/90)

- **Solution addresses the stated problem (34/35)**: The three workstreams map directly to P0/P1/P2 findings with tight traceability. W1 addresses P0 #1 (record-task refs) and P0 #3 (forensic paths). W2 addresses P0 #1 (node_modules), P0 #2 (types/ path), P1 #4 (hardcoded paths), and P1 #6 (dead Playwright templates). W3 addresses P1 #5 (duplication) and P2 #8 (eval inline). The mapping is comprehensive.
- **Scope <-> Solution <-> Success Criteria aligned (27/30)**: Significant improvement from iter 2. In-scope item 8 ("Simplify eval validate-ux inline sub-pipeline -- extraction to a rubric file under eval/") now matches the success criterion with line count target (line 280: "under 280 lines, was 366"). The shared template success criterion now verifies content dedup (line 275), not just file existence. Remaining minor gap: the CI guard is listed as in-scope item 10 and has a success criterion (line 284), but the CI guard is described as a "preventive mechanism" rather than a workstream deliverable -- the categorization is slightly inconsistent.
- **Requirements <-> Solution coherent (25/25)**: Requirements and solution are coherent. The backward compatibility requirement is operationalized through the verification plan (lines 194-197) and e2e test references. The security NFR (path traversal) is addressed by the relative-from-self approach. The install size target (under 5MB) is achievable through W2's 27MB removal. The verification plan maps each workstream to specific checks. Full marks, unchanged.
