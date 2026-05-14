---
date: "2026-05-14"
doc_dir: "docs/proposals/clean-code-task/"
iteration: 2
target_score: "800"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 830/1000** (target: 800)

```
+--------------------------------------------------------------------------+
|                     PROPOSAL QUALITY SCORECARD (1000 pts)                |
+-------------------------------------+----------+----------+--------------+
| Dimension                           | Score    | Max      | Status       |
+-------------------------------------+----------+----------+--------------+
| 1. Problem Definition               |   92     |  110     | +            |
|    Problem clarity                  |   36/40  |          |              |
|    Evidence provided                |   30/40  |          |              |
|    Urgency justified                |   26/30  |          |              |
+-------------------------------------+----------+----------+--------------+
| 2. Solution Clarity                 |  108     |  120     | +            |
|    Approach concrete                |   36/40  |          |              |
|    User-facing behavior             |   38/45  |          |              |
|    Technical direction              |   34/35  |          |              |
+-------------------------------------+----------+----------+--------------+
| 3. Industry Benchmarking            |   98     |  120     | ~            |
|    Industry solutions referenced    |   32/40  |          |              |
|    3+ meaningful alternatives       |   24/30  |          |              |
|    Honest trade-off comparison      |   20/25  |          |              |
|    Chosen approach justified        |   22/25  |          |              |
+-------------------------------------+----------+----------+--------------+
| 4. Requirements Completeness        |   94     |  110     | +            |
|    Scenario coverage                |   36/40  |          |              |
|    Non-functional requirements      |   32/40  |          |              |
|    Constraints & dependencies       |   26/30  |          |              |
+-------------------------------------+----------+----------+--------------+
| 5. Solution Creativity              |   80     |  100     | ~            |
|    Novelty over industry baseline   |   30/40  |          |              |
|    Cross-domain inspiration         |   25/35  |          |              |
|    Simplicity of insight            |   25/25  |          |              |
+-------------------------------------+----------+----------+--------------+
| 6. Feasibility                      |   95     |  100     | +            |
|    Technical feasibility            |   38/40  |          |              |
|    Resource & timeline feasibility  |   28/30  |          |              |
|    Dependency readiness             |   29/30  |          |              |
+-------------------------------------+----------+----------+--------------+
| 7. Scope Definition                 |   72     |   80     | +            |
|    In-scope concrete                |   27/30  |          |              |
|    Out-of-scope explicit            |   22/25  |          |              |
|    Scope bounded                    |   23/25  |          |              |
+-------------------------------------+----------+----------+--------------+
| 8. Risk Assessment                  |   78     |   90     | +            |
|    Risks identified (>=3)           |   26/30  |          |              |
|    Likelihood + impact rated        |   26/30  |          |              |
|    Mitigations actionable           |   26/30  |          |              |
+-------------------------------------+----------+----------+--------------+
| 9. Success Criteria                 |   62     |   80     | ~            |
|    Measurable and testable          |   42/55  |          |              |
|    Coverage complete                |   20/25  |          |              |
+-------------------------------------+----------+----------+--------------+
| 10. Logical Consistency             |   51     |   90     | ~            |
|     Solution <-> Problem            |   22/35  |          |              |
|     Scope <-> Solution <-> Criteria |   18/30  |          |              |
|     Requirements <-> Solution       |   11/25  |          |              |
+-------------------------------------+----------+----------+--------------+
| TOTAL                               |  830     | 1000     |              |
+-------------------------------------+----------+----------+--------------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem:Evidence | "常见累积问题：dead imports、commented-out code、重复逻辑、命名不一致" -- lists problem types but still no quantification (how often? what % of features? any specific feature where this was observed?) | -10 pts (Evidence) |
| Problem:Urgency | "随着 forge pipeline 自动化程度提高，缺乏自动清理环节成为质量闭环的缺口" -- still a qualitative assertion without data or a concrete incident that triggered this proposal | -4 pts (Urgency) |
| Solution:User-facing | "Each change is written to disk and logged in the task record" -- no specification of log format or where the developer reads this record | -4 pts (User-facing behavior) |
| Solution:User-facing | "the agent runs `just test`" -- assumes `just test` exists; no fallback described if project has no justfile or uses a different test runner | -3 pts (User-facing behavior) |
| Benchmarking | "SonarQube (v10.x)" and "ESLint `--fix` in CI (v9.x)" -- version numbers given but no links to documentation, no citation of specific features or published benchmarks | -4 pts (Industry references) |
| Benchmarking | "Pre-commit hook (husky + lint-staged)" listed as an alternative but dismissed for "adds Node.js dependency to Go project" -- this is project-specific cherry-picking that does not apply to non-Go projects; the trade-off is not generalizable | -5 pts (Trade-off comparison) |
| Benchmarking:Selection | "these are semantic problems requiring code understanding, which is why an LLM-based cleanup is appropriate" -- the leap from "semantic cleanup" to "therefore LLM" is under-justified; no discussion of why AST-based or pattern-based semantic analysis (e.g., Semgrep rules, tree-sitter queries) would not suffice | -3 pts (Justification) |
| Requirements:NFR | "总执行时间不超过 5 分钟" -- the 5-minute target is arbitrary; no basis for this number (how many files? what LLM latency? what quality gate duration?) | -4 pts (NFR) |
| Requirements:NFR | "clean-code prompt 模板必须包含 post-cleanup quality gate 步骤" -- this is listed as NFR but is actually a functional requirement; also "hard requirement" phrasing in NFR section is redundant with the risk mitigation table | -4 pts (NFR) |
| Consistency | Solution says "record-driven scope" (line 42) but constraints list 4 manual Go-file registrations (line 67-68) -- these are two different scopes conflated: the runtime scope (which files to clean) vs. the build-time scope (where to register the type). The proposal does not disambiguate. | -8 pts (Solution <-> Problem) |
| Consistency | Success criteria include "覆盖率不低于改动前基线" (line 147) but no criterion validates that clean-code actually performs cleanup -- only that it runs without breaking tests. A clean-code task that does nothing would pass all success criteria. | -7 pts (Scope <-> Criteria) |
| Consistency | Scenario 5 ("Mid-execution crash") says "变更保留在 working tree 供开发者检查，不自动 revert" but scenario 6 ("Regression introduced") says the same thing. Both leave dirty state; the proposal never explains why auto-revert is not applied despite quality gate failure, creating an inconsistency with the stated "quality gate is hard requirement" principle. | -6 pts (Requirements <-> Solution) |

---

## Attack Points

### Attack 1: Success Criteria -- no validation that clean-code actually cleans code

**Where**: Success Criteria lists 8 checkboxes (lines 142-149). None verify that cleanup occurs.

**Why it's weak**: Every criterion tests that the infrastructure works (task appears, dependencies correct, prompt returns, tests pass, docs pass). None test the actual purpose: that the clean-code task identifies and removes dead imports, duplicated logic, commented-out code, or inconsistent naming. A clean-code task whose prompt template says "output 'no cleanup needed'" for every input would pass all 8 criteria. The proposal's own Problem section says the gap is "缺少一个自动化代码清理环节" -- but the success criteria never measure whether cleaning happened, only whether the pipeline step exists.

**What must improve**: Add at least one criterion that validates cleanup effectiveness, e.g., "Given a test fixture with known dead imports and duplicated logic, clean-code task produces a diff that removes at least N of the planted issues." Or: "clean-code task record's cleanup summary lists at least one file modified when run against a feature with intentionally planted cleanup candidates."

### Attack 2: Logical Consistency -- quality gate principle contradicts failure-handling behavior

**Where**: Scenario 6 (line 54-55) states "变更保留在 working tree 供开发者检查，不自动 revert" while NFR (line 62) states "quality gate 失败 = 任务失败，这是硬性要求" and Risk table (line 135) says "quality gate 失败则任务标记 failed."

**Why it's weak**: The proposal establishes quality gate as a "hard requirement" and positions it as the mechanism that prevents regressions from reaching the codebase. But when the quality gate fails (the very case it is designed to catch), the proposal's failure handling leaves the broken changes in the working tree and marks the task as failed. This means a failed clean-code task leaves the codebase in a broken state (regression present) until a human intervenes. The "hard requirement" framing implies automatic enforcement, but the actual behavior is "mark failed and wait for a human." This is inconsistent: either the quality gate should auto-revert on failure (making it a true gate), or the failure handling should acknowledge that the quality gate is advisory and the human is the actual safety net.

**What must improve**: Either (a) add auto-revert on quality gate failure ("if quality gate fails, git checkout the files modified by clean-code before marking the task failed"), or (b) reframe the quality gate as "best-effort validation" and acknowledge in the NFR section that the developer is the final safety net. Pick one framing and make the entire document consistent with it.

### Attack 3: Industry Benchmarking -- the "LLM therefore" leap is under-justified

**Where**: Selection Rationale states "These are semantic problems requiring code understanding, which is why an LLM-based cleanup is appropriate" (line 94).

**Why it's weak**: The proposal establishes that industry tools (ESLint, Ruff, Biome) handle syntactic cleanup but not semantic cleanup. Fair. But the jump from "semantic cleanup needed" to "LLM is the right tool" skips an entire category of non-LLM semantic analysis: Semgrep rules can detect duplicated logic patterns, tree-sitter queries can find dead code across files, AST diffing tools can detect structural clones. The proposal never explains why these deterministic, fast, reproducible alternatives were rejected in favor of a non-deterministic LLM that requires a quality gate precisely because it is unreliable. The comparison table's "LLM-based (non-deterministic)" in the Cons column is acknowledged but never resolved -- it is listed as a con and then ignored.

**What must improve**: Add one paragraph explaining why deterministic semantic analysis tools (Semgrep, tree-sitter based dead-code detectors) were considered and rejected, or why the LLM approach is preferred despite non-determinism. If the answer is "Forge already uses LLMs for all pipeline tasks so consistency dictates LLM here too," state that explicitly -- it is a valid rationale that is currently implicit.

---

## Previous Issues Check

| Previous Attack (Iteration 1) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: Industry Benchmarking -- superficial and straw-man-filled alternatives | Partially | Three real industry tools now cited (SonarQube, ESLint/Ruff, Biome) with version numbers. Straw-man `/simplify` removed. However, no documentation links, and the "pre-commit hook" alternative is still a somewhat weak entry that is dismissed for project-specific reasons. Also missing: Semgrep, tree-sitter, or other deterministic semantic tools as alternatives. |
| Attack 2: Solution Clarity -- no user-facing behavior described | Yes | New "User-Facing Behavior" subsection (lines 28-38) with 5 numbered items covering what the developer sees before, during, and after execution; failure handling; and empty-output case. This was the single biggest improvement. |
| Attack 3: Requirements -- missing error scenarios and performance NFRs | Yes | 7 scenarios now (from 4), including mid-execution crash (scenario 5), regression introduced (scenario 6), and empty LLM output (scenario 7). Performance NFR added (5-minute / 15-minute targets). Reliability NFR added (quality gate as hard requirement). Observability NFR added (cleanup record). |

---

## Verdict

- **Score**: 830/1000
- **Target**: 800/1000
- **Gap**: 0 points (target exceeded by 30)
- **Action**: Target reached. Primary remaining weaknesses are Success Criteria (no effectiveness validation) and Logical Consistency (quality gate framing vs. failure behavior). These are quality gaps but not blockers for proceeding to implementation.
