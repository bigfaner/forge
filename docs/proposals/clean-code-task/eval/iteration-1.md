---
date: "2026-05-14"
doc_dir: "docs/proposals/clean-code-task/"
iteration: 1
target_score: "800"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 1

**Score: 640/1000** (target: 800)

```
+--------------------------------------------------------------------------+
|                     PROPOSAL QUALITY SCORECARD (1000 pts)                |
+-------------------------------------+----------+----------+-------------+
| Dimension                           | Score    | Max      | Status      |
+-------------------------------------+----------+----------+-------------+
| 1. Problem Definition               |   82     |  110     | +           |
|    Problem clarity                  |   35/40  |          |             |
|    Evidence provided                |   28/40  |          |             |
|    Urgency justified                |   19/30  |          |             |
+-------------------------------------+----------+----------+-------------+
| 2. Solution Clarity                 |   82     |  120     | +           |
|    Approach concrete                |   30/40  |          |             |
|    User-facing behavior             |   25/45  |          |             |
|    Technical direction              |   27/35  |          |             |
+-------------------------------------+----------+----------+-------------+
| 3. Industry Benchmarking            |   50     |  120     | x           |
|    Industry solutions referenced    |   15/40  |          |             |
|    3+ meaningful alternatives       |   12/30  |          |             |
|    Honest trade-off comparison      |   10/25  |          |             |
|    Chosen approach justified        |   13/25  |          |             |
+-------------------------------------+----------+----------+-------------+
| 4. Requirements Completeness        |   62     |  110     | x           |
|    Scenario coverage                |   30/40  |          |             |
|    Non-functional requirements      |   18/40  |          |             |
|    Constraints & dependencies       |   14/30  |          |             |
+-------------------------------------+----------+----------+-------------+
| 5. Solution Creativity              |   65     |  100     | ~           |
|    Novelty over industry baseline   |   25/40  |          |             |
|    Cross-domain inspiration         |   15/35  |          |             |
|    Simplicity of insight            |   25/25  |          |             |
+-------------------------------------+----------+----------+-------------+
| 6. Feasibility                      |   85     |  100     | +           |
|    Technical feasibility            |   35/40  |          |             |
|    Resource & timeline feasibility  |   25/30  |          |             |
|    Dependency readiness             |   25/30  |          |             |
+-------------------------------------+----------+----------+-------------+
| 7. Scope Definition                 |   68     |   80     | +           |
|    In-scope concrete                |   28/30  |          |             |
|    Out-of-scope explicit            |   20/25  |          |             |
|    Scope bounded                    |   20/25  |          |             |
+-------------------------------------+----------+----------+-------------+
| 8. Risk Assessment                  |   60     |   90     | ~           |
|    Risks identified (>=3)           |   20/30  |          |             |
|    Likelihood + impact rated        |   20/30  |          |             |
|    Mitigations actionable           |   20/30  |          |             |
+-------------------------------------+----------+----------+-------------+
| 9. Success Criteria                 |   48     |   80     | x           |
|    Measurable and testable          |   35/55  |          |             |
|    Coverage complete                |   13/25  |          |             |
+-------------------------------------+----------+----------+-------------+
| 10. Logical Consistency             |   38     |   90     | x           |
|     Solution <-> Problem            |   15/35  |          |             |
|     Scope <-> Solution <-> Criteria |   13/30  |          |             |
|     Requirements <-> Solution       |   10/25  |          |             |
+-------------------------------------+----------+----------+-------------+
| TOTAL                               |  640     | 1000     |             |
+-------------------------------------+----------+----------+-------------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem:line 16 | "常见累积问题" -- no quantification, no frequency data, no user reports | -12 pts (Evidence) |
| Problem:line 21 | "随着 forge pipeline 自动化程度提高" -- vague future projection with no data | -8 pts (Urgency) |
| Solution:line 25 | "审查本 feature 所有已变更文件并执行代码清理" -- what does "代码清理" concretely entail? No specification of actions | -10 pts (Approach) |
| Solution:line 31 | "只清理明确的问题（dead code、duplication）" -- no definition of what qualifies as "明确" | -8 pts (Approach) |
| Solution | No user-facing behavior section -- what does the developer see when clean-code runs? stdout? logs? progress? | -20 pts (User-facing behavior) |
| Benchmarking:line 56 | "大多数 CI/CD pipeline" -- no specific CI/CD product named, no links, no references | -15 pts (Industry references) |
| Benchmarking:line 56 | "GitHub Actions 的 super-linter" -- one named tool, not enough for robust benchmarking | -10 pts (Industry references) |
| Benchmarking:line 64 | "手动 `/simplify`" -- straw-man alternative; presented only to be dismissed for not being "pipeline 级别" | -20 pts (Alternatives, per deduction rules) |
| Benchmarking:line 56 | "IDE 的 Reformat Code" -- not a pipeline-level solution, barely comparable | -8 pts (Alternatives) |
| Benchmarking:line 65 | "最符合 forge 自动化理念" -- circular justification, not benchmarked against alternatives | -12 pts (Justification) |
| Requirements:line 44 | "向后兼容: 新增 task type 不影响现有 index.json（omitempty 已保证）" -- no explanation of HOW backward compatibility is verified | -10 pts (NFR) |
| Requirements | No error scenario coverage -- what if clean-code fails mid-execution? What if it introduces a regression despite conservative mode? | -12 pts (Scenario coverage) |
| Requirements | No performance NFR -- how long should clean-code take? What if the feature has 100 changed files? | -10 pts (NFR) |
| Requirements:line 49 | "Prompt 模板需嵌入 pkg/prompt/data/" -- dependency without version/status | -8 pts (Constraints) |
| Risks:line 104 | "模板强调 conservative approach" -- mitigation is "the prompt says be careful", not an engineering guard | -10 pts (Mitigations) |
| Risks | Only 3 risks, one is cosmetic (T-test-4.6 numbering "看起来不直观"), not a meaningful risk | -10 pts (Risks identified) |
| Success:line 113 | "所有现有测试通过（go test -race -cover ./...）" -- no coverage threshold, just "pass" | -8 pts (Measurable) |
| Success | No criterion verifying that clean-code actually cleans anything -- only that it exists and runs | -12 pts (Coverage) |
| Consistency | Solution claims "record-driven scope" but requirements list 4 manual file registrations -- contradictory framing | -12 pts (Solution <-> Problem) |
| Consistency | Scope lists "测试覆盖：更新 testgen_test、infer_test" but success criteria has no test for clean-code task type itself | -10 pts (Scope <-> Criteria) |

---

## Attack Points

### Attack 1: Industry Benchmarking -- superficial and straw-man-filled alternatives section

**Where**: "大多数 CI/CD pipeline 在 test 之后有 lint/format 步骤（如 GitHub Actions 的 super-linter）。IDE 的 Reformat Code 是最接近的手动操作。" and "手动 `/simplify` | forge skill | 灵活可控 | 依赖人工触发，容易遗忘 | Rejected: 不是 pipeline 级别"

**Why it's weak**: The industry benchmarking names exactly one external tool (super-linter) and one generic concept (IDE reformat). The `/simplify` alternative is a textbook straw man -- it is listed only to be rejected for "not being pipeline level," which is the premise of the proposal itself, making the rejection circular. No industry-validated automated cleanup tools are evaluated: SonarQube cleanup rules, ESLint --fix with CI enforcement, Ruff's automated fix mode, Prettier in CI, Biome's lint-and-fix, or even language-agnostic tools like Semgrep. The comparison table has no source URLs, no version numbers, and no concrete feature comparison.

**What must improve**: Cite at least 3 specific industry tools with version numbers and explain what they do. Replace the straw-man `/simplify` with a genuine alternative (e.g., "run lint --fix as a CI step after tests" which IS an industry-validated approach). Provide honest trade-off analysis explaining why an in-pipeline task is superior to a CI lint step, not just "it fits our philosophy."

### Attack 2: Solution Clarity -- no user-facing behavior described

**Where**: The entire Proposed Solution section describes internal mechanics ("插入一个新的 clean-code 任务类型", "由 forge task index 自动生成") but never says what the developer experiences.

**Why it's weak**: A developer using forge cannot determine from this proposal: (1) Will they see a new task appear in `forge task list`? (2) What output does clean-code produce -- a diff? a summary? (3) Does it auto-apply changes or propose them for review? (4) What happens if the cleanup breaks something -- is there a rollback? The proposal describes the plumbing but not the faucet. The "Approach is concrete" criterion asks whether "a reader can explain back what will be built" -- a reader cannot explain the observable behavior because it is never stated.

**What must improve**: Add a "User-Facing Behavior" subsection describing: what the developer sees before, during, and after clean-code execution; whether changes are auto-applied or require approval; what output/log format is produced; and what the developer does if the result is unsatisfactory.

### Attack 3: Requirements Completeness -- missing error scenarios and performance requirements

**Where**: Requirements Analysis only covers 4 "happy/flag" scenarios and 2 minimal NFRs. No error paths. No performance bounds.

**Why it's weak**: The scenario coverage lists only happy-path variants (all tasks done, multi-profile, --no-test, no changes needed). Missing scenarios: (a) clean-code fails mid-execution -- what state is the codebase left in? (b) clean-code introduces a regression despite "conservative" mode -- what is the recovery path? (c) the feature has a very large diff (hundreds of files) -- does clean-code time out? (d) clean-code runs but the LLM produces no useful changes -- is this a success or failure? The NFRs mention backward compatibility and idempotency but ignore performance (how long is acceptable?), reliability (what if the LLM returns invalid code?), and observability (how do we know clean-code did its job?). This is a significant gap for a pipeline step that modifies code automatically.

**What must improve**: Add at least 3 error/failure scenarios to scenario coverage. Add performance NFR (max execution time). Add reliability NFR (what validates the cleanup did not break things -- the quality gate mentioned in risks should be a requirement, not just a risk mitigation). Add observability NFR (what metric or log confirms cleanup occurred and was effective).

---

## Previous Issues Check

*Not applicable -- iteration 1.*

---

## Verdict

- **Score**: 640/1000
- **Target**: 800/1000
- **Gap**: 160 points
- **Action**: Continue to iteration 2 -- primary improvement areas are Industry Benchmarking (+70 potential), Solution Clarity (+38 potential), Requirements Completeness (+48 potential), and Success Criteria (+32 potential)
