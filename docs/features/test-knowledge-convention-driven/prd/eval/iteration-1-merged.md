# PRD Eval — Iteration 1 Merged Report

**Type**: PRD (Mode B — no UI)
**Experts**: PM + QA
**PM Score**: 735/1000
**QA Score**: 756/1000
**Gate Score**: 746/1000 (target: 900)

## Dimension Breakdown

| Dimension | PM | QA | Avg |
|-----------|-----|-----|-----|
| Background & Goals | 75/100 | 93/100 | 84/100 |
| Flow Diagrams | 120/150 | 125/150 | 123/150 |
| Flow Completeness | 145/200 | 160/200 | 153/200 |
| User Stories | 155/200 | 160/200 | 158/200 |
| Scenario Completeness | 105/150 | 90/150 | 98/150 |
| Edge Case Coverage | 55/100 | 45/100 | 50/100 |
| Scope Clarity | 80/100 | 83/100 | 82/100 |

## Merged Attack Points

1. **Background & Goals**: 85% compile-rate target lacks baseline evidence — "First-pass compile rate >= 85%... Measured on forge-cli's 126+ existing tests" — must provide baseline data from current Profile-based generation and justify why LLM-default + Convention can match or exceed it
2. **Flow Completeness**: Only main flow has error handling; bootstrap and test-guide flows document zero failure modes — "If compile fails (max 2 retries) → feed error back to LLM → regenerate" is the only error path; bootstrap flow (4 steps) and test-guide flow (5 steps) have no error branches
3. **Flow Completeness**: Missing error paths for malformed Convention files — "Convention files use fixed sections: Framework, Assertion, Tags, Result Format" (FS-1) — no error handling when a Convention file has missing/empty required sections, invalid frontmatter, or corrupted content
4. **User Stories**: Default framework user segment has no dedicated story — "Default framework users: Profile works for them but they get no benefit from Profile over LLM defaults" identifies the segment but no story covers backward compatibility for this group
5. **User Stories**: No ACs cover error/failure scenarios — all 5 stories only have happy-path ACs — must add acceptance criteria for: Convention load failure, compile gate exhaustion, and Reconnaissance signal conflicts
6. **User Stories**: Story 3 AC is a population metric, not per-invocation verifiable — "just e2e-compile passes with >= 85% first-pass success rate" cannot be verified from a single user run
7. **Edge Case Coverage**: Convention file validation and structure schema entirely absent — "Convention files use fixed sections: Framework, Assertion, Tags, Result Format (minimum set)" (FS-1) defines section names but no field-level schema, no required vs optional constraints, no validation rules
8. **Edge Case Coverage**: No recovery guidance after terminal compile failure — "Blocked: compile gate failed" (Flow Diagram) — must describe what the user should do next: edit Convention? re-run test-guide? manually fix code?
9. **Scenario Completeness**: consolidate-specs integration in scope but zero scenario/user story coverage — "Integrate Convention files into consolidate-specs management" (In Scope item 10) has no functional spec detail, no flow, no user story
10. **Scenario Completeness**: Compile gate retry count conflicts with business rules — "max 2 retries" (Flow Description) vs BIZ-quality-gate-001 "retry-once policy" — must clarify whether compile gate is part of the quality-gate pipeline and reconcile the retry count
11. **Scenario Completeness**: Unstated prerequisite that just e2e-compile always exists — "Compile gate: just e2e-compile validates generated code" (Flow step 5) — must document what happens when no justfile exists or the e2e-compile recipe is missing
12. **[blindspot]**: No go/no-go criteria for phase transitions — "Point of no return: Phase 3 start" defines the boundary but not what metric triggers the decision — must define explicit criteria that gate each phase
13. **[blindspot]**: Code Reconnaissance underspecified core capability — "Fact Table (runtime LLM notes, not persisted) extended to collect test framework info: file patterns, import analysis, build tag analysis, function signature patterns" (FS-3) — must define reliability expectations, failure modes when reconnaissance finds nothing useful, and how it degrades
14. **[blindspot]**: Compile-rate measures compilation not correctness — "gen-test-scripts output passes just e2e-compile on first attempt" — must add a semantic correctness check (e.g., generated test actually exercises the Contract's test steps) or explicitly acknowledge this gap
15. **[blindspot]**: Cold-start 85% claim contradicts feature motivation — "Given a new project with no Convention files and no existing test files ... just e2e-compile passes with >= 85% first-pass success rate" (Story 3 AC) — if LLM defaults achieve 85% without Convention, the feature's premise that Convention is needed is undermined
16. **[blindspot]**: No audit mechanism for Profile consumer completeness — "Delete pkg/profile/ entirely. Rewrite all 19+ consumers." (FS-7) — "19+" indicates uncertain consumer count; must include an import-audit step to ensure no untracked consumers exist before deletion
