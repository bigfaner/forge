# Proposal Evaluation Report — Iteration 1

**Document**: `docs/proposals/auto-gen-journeys-contracts/proposal.md`
**Date**: 2026-05-23
**Scorer**: CTO-level adversarial review

---

## Pre-Score Anchors (Phase 1 — Independent Reasoning Audit)

1. **Problem → Solution gap**: The problem is "gen-journeys and gen-contracts require manual invocation, breaking the automated pipeline." The solution (make them auto-generated tasks) directly addresses this. The problem-solution link is sound.

2. **Evidence → Solution alignment**: Evidence cites specific code (`autogen.go`, `GetBreakdownTestTasks()`) and user-visible behavior (manual `/gen-journeys` calls). Concrete and verifiable.

3. **Self-contradiction — hardcoded indices**: The proposal claims `resolveBreakdownDeps` uses `findTaskIndex` by ID and doesn't rely on hardcoded indices. Code review reveals this is **partially false**: lines 426-427 of `autogen.go` hardcode `evalJourneyIdx := 0` and `evalContractIdx := 1`. Only later tasks (validate, specs-consolidate, clean-code) use `findTaskIndex`. When gen-journeys and gen-contracts are inserted as new tasks at positions 0 and 1, these hardcoded indices (`0`, `1`, `2`) will point to wrong tasks. The proposal's risk mitigation is factually incorrect.

4. **Quick mode architectural change understated**: Replacing `gen-and-run` with split tasks fundamentally changes the Quick pipeline topology. `resolveQuickDeps` currently uses index-based arithmetic (`verifyIdx := nTypes`) assuming gen-and-run tasks occupy positions 0..nTypes-1. With split tasks (gen-journeys + gen-contracts + gen-scripts per type), the index math is entirely different. The proposal does not discuss this.

5. **SKILL.md semantic conflicts unaddressed**: gen-journeys SKILL.md has a HARD-RULE requiring user approval before commit. gen-contracts SKILL.md requires eval-journey to have passed. Neither is discussed as needing adaptation for non-interactive auto-task execution.

---

## Rubric Scoring (Phase 2)

### Dimension 1: Problem Definition — 82/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | Core problem is clear: pipeline gap between business tasks and test generation. Minor ambiguity: "automated pipeline" could mean CI/CD or Forge's autogen — the document doesn't disambiguate until the Solution section. |
| Evidence provided | 22/40 | Only one form of evidence: code structure analysis (autogen.go references). No user feedback, no issue tracking references, no concrete failure scenarios with reproduction steps. The evidence is "we looked at the code" rather than "this broke in production" or "users reported friction." |
| Urgency justified | 25/30 | v3.0.0 dependency is stated. Cost of delay is implied (pipeline won't be end-to-end) but not quantified. What specifically breaks in the release without this? |

**Deductions**: -10 for evidence limited to code inspection with no user-facing impact data; -18 for lack of concrete failure reports or user friction evidence.

### Dimension 2: Solution Clarity — 85/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 35/40 | Pipeline diagrams are clear. Task type names, template files, and mode-specific behavior are specified. A reader can explain back what will be built. |
| User-facing behavior described | 35/45 | The observable behavior (automatic task generation in `forge task index`) is described. However, the Quick mode user experience change (gen-and-run disappears, replaced by split tasks) is mentioned but not elaborated — what does the user see differently in `task index` output? How does task count change? |
| Technical direction clear | 15/35 | High-level direction (add task types to types.go, add templates to embed.FS) is stated. But the critical implementation detail — how dependency resolution indices change when new tasks are inserted — is factually wrong (see blindspot). The ResolveFirstTestDep adaptation is listed in scope but not explained. |

**Deductions**: -10 for incomplete user-facing behavior in Quick mode; -20 for technical direction containing a factual error about index-based dependency resolution; -5 for not explaining how ResolveFirstTestDep must change for Quick mode's new first task.

### Dimension 3: Industry Benchmarking — 45/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 15/40 | One vague sentence: "CI/CD pipelines, test generation is usually manual or semi-automatic." No specific products, open-source projects, or published patterns cited. No links, no names (e.g., Cypress, Playwright codegen, Selenium IDE, Testim, Mabl). |
| At least 3 meaningful alternatives | 15/30 | Three alternatives listed, but "merge into gen-and-run" is borderline straw-man (presented mainly to be rejected — "user explicitly requested split"). "Do nothing" is genuine. "Independent auto-tasks" is the proposal itself. None are industry-validated solutions. |
| Honest trade-off comparison | 10/25 | Pros/cons exist but are superficial. "Slightly more development work" is vague — how much? The comparison table lacks quantitative depth. |
| Justified against benchmarks | 5/25 | The proposal states "Forge's differentiation" but doesn't justify why this approach is better than any industry pattern. No benchmark to justify against. |

**Deductions**: This is the weakest dimension. The benchmarking section reads as an afterthought rather than a genuine analysis.

### Dimension 4: Requirements Completeness — 70/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 25/40 | Four key scenarios listed. Missing edge cases: (1) What happens when gen-journeys produces zero journeys? (2) What if proposal.md has minimal content? (3) Error scenarios during auto-task execution. (4) What if test profile system is not configured? |
| Non-functional requirements | 25/40 | NFRs are minimal: MainSession=false for both tasks, embed.FS consistency. Missing: execution time impact (gen-journeys can take 10-30 minutes), token/resource consumption for auto-generated tasks, error handling behavior when tasks fail in the pipeline. |
| Constraints & dependencies | 20/30 | Key dependencies listed (PRD, proposal.md, test profile system). Missing: gen-contracts SKILL.md requires eval-journey to have passed (Blocker condition) — this constraint conflicts with Quick mode skipping eval. |

**Deductions**: -15 for missing edge cases; -15 for thin NFRs; -10 for undiscussed constraint conflicts.

### Dimension 5: Solution Creativity — 30/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 15/40 | The proposal itself acknowledges: "No special innovation. Natural extension of existing framework." Honest but scores low — it's connecting two existing systems, not creating anything new. |
| Cross-domain inspiration | 5/35 | No cross-domain references at all. The proposal is purely inward-looking at Forge's own architecture. |
| Simplicity of insight | 10/25 | The insight is simple (add two task types to an existing framework), which is good. But the implementation isn't as simple as presented — the index-based dependency resolution and SKILL.md semantic conflicts reveal hidden complexity. |

**Deductions**: The self-assessment of "no innovation" is honest. Score reflects the straightforward nature of the work.

### Dimension 6: Feasibility — 60/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 25/40 | Stated as "fully feasible" but the factual error about `resolveBreakdownDeps` using `findTaskIndex` (when it actually hardcodes indices 0 and 1) means the author may not have fully understood the code they need to modify. This introduces execution risk. |
| Resource & timeline | 20/30 | "5-8 coding tasks" is stated but not justified. Given the hidden complexity (resolveBreakdownDeps rewrite, resolveQuickDeps rewrite, ResolveFirstTestDep adaptation, SKILL.md modifications), this estimate may be optimistic. |
| Dependency readiness | 15/30 | Test profile system stated as "merged." gen-journeys/gen-contracts SKILL.md stated as "stable." But the SKILL.md files contain HARD-RULEs that conflict with non-interactive execution, meaning they are NOT ready for this use case without modification. |

**Deductions**: -15 for factual error in feasibility assessment; -10 for understated scope; -15 for missed SKILL.md readiness gap.

### Dimension 7: Scope Definition — 62/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | Nine concrete deliverables listed. Most are specific (type names, file names). "Update ARCHITECTURE.md and related documents" is vague — what related documents? |
| Out-of-scope explicitly listed | 22/25 | Four out-of-scope items clearly stated. Good boundary setting. |
| Scope is bounded | 15/25 | "5-8 coding tasks" is a bounded estimate, but the actual scope is larger than presented due to undocumented changes needed (resolveQuickDeps rewrite, ResolveFirstTestDep Quick branch, SKILL.md HARD-RULE adaptations). |

**Deductions**: -5 for vague documentation item; -3 for incomplete out-of-scope (SKILL.md HARD-RULE changes are in neither in nor out); -10 for scope boundary that underestimates actual work.

### Dimension 8: Risk Assessment — 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 18/30 | Three risks listed. Missing risks: (1) SKILL.md HARD-RULE conflicts with non-interactive execution. (2) resolveQuickDeps complete rewrite required. (3) ResolveFirstTestDep Quick branch points to wrong task ID. (4) gen-contracts eval-journey Blocker condition in Quick mode. |
| Likelihood + impact rated | 17/30 | Ratings exist. The "resolveBreakdownDeps index" risk is rated M/H with a mitigation that is factually incorrect — the actual likelihood should be "certain" if implemented as described. |
| Mitigations are actionable | 20/30 | First mitigation (downgrade strategy) is actionable. Second mitigation (findTaskIndex claim) is incorrect. Third mitigation (keep type definition) is actionable. |

**Deductions**: -12 for missing critical risks; -13 for incorrect risk mitigation; -10 for unactionable/incorrect mitigation.

### Dimension 9: Success Criteria — 58/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Measurable and testable | 40/55 | Six criteria, mostly testable. Issue: `test.test-gen-contracts` has a duplicate `test.` prefix (typo: should be `test.gen-contracts`). The criterion "all existing tests pass" is testable but doesn't cover the new pipeline's functional correctness — only regression. |
| Coverage is complete | 18/25 | Covers main scenarios but misses: (1) SKILL.md HARD-RULE adaptation verification. (2) Error handling when gen-journeys produces no output. (3) Quick mode pipeline produces different task set than Breakdown. |

**Deductions**: -15 for typo in criterion and lack of functional correctness tests; -7 for coverage gaps.

### Dimension 10: Logical Consistency — 50/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 25/35 | Yes, making gen-journeys/gen-contracts auto-generated directly addresses the pipeline gap. But the Quick mode change (replacing gen-and-run) is scope expansion beyond the stated problem, which is about automation, not restructuring. |
| Scope ↔ Solution ↔ Success Criteria aligned | 10/30 | Multiple misalignments: (1) Scope says "update ResolveFirstTestDep" but solution doesn't explain how. (2) Success criteria don't verify dependency chain correctness. (3) Scope includes "modify gen-journeys SKILL.md for proposal.md input" but no success criterion verifies this works. (4) The `test.test-gen-contracts` typo means success criteria may be measuring the wrong thing. |
| Requirements ↔ Solution coherent | 15/25 | Quick mode proposal.md input is a requirement with no detailed solution (how does gen-journeys detect and handle proposal.md vs PRD?). The "backward compatibility" requirement has no corresponding success criterion. |

**Deductions**: -10 for scope expansion beyond problem statement; -20 for scope/solution/criteria misalignment; -10 for orphan requirements.

---

## Cross-Dimension Coherence Check

The most significant cross-cutting issue is the **factual error about resolveBreakdownDeps**. This affects:
- Dimension 2 (Solution Clarity): Technical direction is wrong
- Dimension 6 (Feasibility): Author may not understand the code sufficiently
- Dimension 8 (Risk Assessment): Mitigation is incorrect
- Dimension 10 (Logical Consistency): Solution cannot work as described

This single issue cascades across four dimensions and represents the document's most critical weakness.

---

## Phase 3 — Blindspot Hunt

### [blindspot-1] resolveBreakdownDeps uses hardcoded indices 0 and 1, not findTaskIndex
The proposal states: "resolveBreakdownDeps uses findTaskIndex by ID, not hardcoded indexes." Code at lines 426-427 shows `evalJourneyIdx := 0` and `evalContractIdx := 1`. When gen-journeys and gen-contracts are inserted at the beginning of the task array, all subsequent index arithmetic (`genStart := 2`, `runIdx := genStart + nTypes`, `verifyIdx := runIdx + 1`) will be wrong. The entire `resolveBreakdownDeps` function needs to be rewritten, not just extended. This is a **critical** oversight.

### [blindspot-2] resolveQuickDeps requires complete rewrite for new topology
Current Quick mode generates `gen-and-run` tasks at positions 0..nTypes-1 with `verifyIdx := nTypes`. The proposal replaces gen-and-run with independent gen-journeys + gen-contracts + gen-scripts + run tasks. The index arithmetic is entirely different. The proposal does not acknowledge this function needs rewriting.

### [blindspot-3] gen-contracts SKILL.md eval-journey Blocker conflicts with Quick mode
gen-contracts SKILL.md states: "Run /eval --type journey first. Blocker: do not proceed if any Journey scored below target." Quick mode explicitly skips eval. This means gen-contracts SKILL.md MUST be modified for Quick mode, but this modification is listed as out-of-scope.

### [blindspot-4] gen-journeys SKILL.md HARD-RULE conflicts with non-interactive execution
gen-journeys SKILL.md states: "Do NOT commit documents automatically. Present all generated Journey files to the user for review and wait for explicit approval before committing." When gen-journeys runs as an auto-generated task, there is no interactive user session. The proposal does not address how to adapt this behavior.

### [blindspot-5] ResolveFirstTestDep Quick branch references T-quick-gen-and-run
Code at line 571 shows `firstTestIdx := findTaskIndexByPrefix(tasks, "T-quick-gen-and-run")`. After replacing gen-and-run, this prefix will never match, and the Quick mode first test task will have no dependency wired. The proposal mentions updating ResolveFirstTestDep but doesn't specify this critical change.

### [blindspot-6] Success Criteria typo: test.test-gen-contracts
The criterion states `test.test-gen-contracts` which has a duplicate `test.` prefix. This should be `test.gen-contracts`. If implemented literally, the CLI type would be nonsensical.

### [blindspot-7] data/test-gen-and-run.md template orphaned
The embed template file `data/test-gen-and-run.md` is kept in the codebase (per out-of-scope) but with gen-and-run tasks no longer generated, this template becomes dead code. No disposal strategy is specified beyond "mark deprecated."

---

## Freeform Finding Integration

| Finding | Rubric Dimension | Resolution |
|---------|-----------------|------------|
| [high] resolveBreakdownDeps hardcoded indices | Dim 2, 6, 8, 10 | **Confirmed**: proposal's claim is factually wrong. Coded as critical deductions across 4 dimensions. |
| [high] resolveQuickDeps rewrite needed | Dim 2, 6 | **Confirmed**: architectural change not discussed. Deducted from Solution Clarity and Feasibility. |
| [high] gen-contracts eval-journey Blocker conflict | Dim 4, 7, 10 | **Confirmed**: constraint conflict makes scope incomplete. Listed as blindspot-3. |
| [medium] gen-journeys HARD-RULE interactive conflict | Dim 4, 6 | **Confirmed**: SKILL.md not ready for non-interactive use. Listed as blindspot-4. |
| [medium] Quick mode proposal.md input quality | Dim 4, 6 | **Confirmed**: downgrade strategy is vague. Deducted from Requirements Completeness. |
| [medium] ResolveFirstTestDep Quick branch | Dim 2, 6 | **Confirmed**: listed in scope but not explained. Listed as blindspot-5. |
| [low] test.test-gen-contracts typo | Dim 9 | **Confirmed**: success criterion contains error. Listed as blindspot-6. |
| [low] data/test-gen-and-run.md orphaned | Dim 7 | **Confirmed**: dead code with no disposal plan. Listed as blindspot-7. |

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 82 | 110 |
| 2. Solution Clarity | 85 | 120 |
| 3. Industry Benchmarking | 45 | 120 |
| 4. Requirements Completeness | 70 | 110 |
| 5. Solution Creativity | 30 | 100 |
| 6. Feasibility | 60 | 100 |
| 7. Scope Definition | 62 | 80 |
| 8. Risk Assessment | 55 | 90 |
| 9. Success Criteria | 58 | 80 |
| 10. Logical Consistency | 50 | 90 |
| **Total** | **597** | **1000** |
