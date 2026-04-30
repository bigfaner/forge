---
date: "2026-04-29"
doc_dir: "docs/proposals/justfile-standard-vocabulary/"
iteration: "2"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 85/100** (target: 90)

```
+---------------------------------------------------------------+
|                    PROPOSAL QUALITY SCORECARD                  |
+------------------------------+----------+----------+----------+
| Dimension                    | Score    | Max      | Status   |
+------------------------------+----------+----------+----------+
| 1. Problem Definition        |  18      |  20      | :warning:|
|    Problem clarity           |   7/7    |          |          |
|    Evidence provided         |   6/7    |          |          |
|    Urgency justified         |   5/6    |          |          |
+------------------------------+----------+----------+----------+
| 2. Solution Clarity          |  17      |  20      | :warning:|
|    Approach concrete         |   7/7    |          |          |
|    User-facing behavior      |   5/7    |          |          |
|    Differentiated            |   5/6    |          |          |
+------------------------------+----------+----------+----------+
| 3. Alternatives Analysis     |  14      |  15      | :white_check_mark:|
|    Alternatives listed (>=2) |   5/5    |          |          |
|    Pros/cons honest          |   5/5    |          |          |
|    Rationale justified       |   4/5    |          |          |
+------------------------------+----------+----------+----------+
| 4. Scope Definition          |  14      |  15      | :white_check_mark:|
|    In-scope concrete         |   5/5    |          |          |
|    Out-of-scope explicit     |   5/5    |          |          |
|    Scope bounded             |   4/5    |          |          |
+------------------------------+----------+----------+----------+
| 5. Risk Assessment           |  11      |  15      | :warning:|
|    Risks identified (>=3)    |   5/5    |          |          |
|    Likelihood + impact rated |   3/5    |          |          |
|    Mitigations actionable    |   3/5    |          |          |
+------------------------------+----------+----------+----------+
| 6. Success Criteria          |  11      |  15      | :warning:|
|    Measurable                |   4/5    |          |          |
|    Coverage complete         |   4/5    |          |          |
|    Testable                  |   3/5    |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  85      |  100     |          |
+------------------------------+----------+----------+----------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem section | Evidence cites "多个项目（如 pm-work-tracker）" but does not quantify how many projects or skills are affected -- "多个" is vague for a countable metric | -1 pt |
| Problem section | Urgency implied but never explicit -- no "what breaks if we delay" or deadline | -1 pt |
| Solution section | User-facing behavior described primarily through skill-internal logic (the 3-step scope resolution algorithm); the end-user experience (what the developer sees when running a skill) is never stated | -2 pts |
| Solution section | Differentiation from alternatives is argued inside the Alternatives section rather than in-place; the Solution section itself does not explain why this approach is distinct | -1 pt |
| Alternatives section | Rationale for choosing A is strong on paper but uses circular logic: "方案 A 需更新 10 个 skill..." then concludes "一次性迁移的成本是可控的确定性工作" -- this asserts controllability without evidence (e.g., prior migration precedent or time estimate) | -1 pt |
| Scope section | Scope item 7 ("更新 forge 项目 justfile 作为参考实现") has no corresponding success criterion that verifies the reference implementation actually contains all 16 commands and passes them -- partially covered by criterion 5 but 5 only says "包含全部 16 个标准命令" without execution verification | -1 pt |
| Risks section | Likelihood ratings are present but internally inconsistent: Risk 2 rates Skill behavior change as "High" likelihood and "High" impact, yet the mitigation is "分两阶段部署" which is a standard deployment strategy, not a risk-specific mitigation -- a High/High risk should have a more aggressive mitigation than phased rollout | -2 pts |
| Risks section | Risk 4 mitigation ("init-justfile 检测已有 justfile 时提示用户确认") is partially contradicted by the same risk's boundary marker proposal ("# --- forge standard recipes ---") -- if the tool uses boundary markers to preserve user content, why does it also need to prompt? The two mitigations solve different sub-problems but this is not acknowledged | -2 pts |
| Success Criteria | Criterion 3 says "不再包含原始 shell 命令（如 `go test`、`npm run build`）" but the parenthetical examples are incomplete -- what about `cargo test`, `python -m pytest`, `npx serve`, or other toolchain-specific commands? The boundary of "raw shell command" is defined by example rather than by rule | -1 pt |
| Success Criteria | Criterion 7 is excellent for breakdown-tasks scope field but mixes two distinct verifiable claims: (a) scope field exists with valid values, and (b) scope values are correctly assigned. A single criterion should verify one thing | -1 pt |

---

## Attack Points

### Attack 1: Success Criteria — Criterion 3 boundary is still fuzzy despite iteration 1 fix

**Where**: Criterion 3: "所有 skill/agent/command 文件中不再包含原始 shell 命令（如 `go test`、`npm run build`），统一通过 just 命令调用"
**Why it's weak**: The criterion uses a negative test ("不再包含" / "no longer contain") with an incomplete example list. This creates an ambiguity: does `go test -race -coverprofile=coverage.out ./...` count as a "raw shell command"? What about `go test` inside a comment? What about `npm run build` inside a markdown code block that documents expected behavior rather than prescribing execution? The parenthetical "(如 `go test`、`npm run build`)" is illustrative, not exhaustive, and no rule is given for what constitutes a forbidden raw command. A team could pass this criterion by removing the two cited examples but leaving other raw commands intact. The criterion should either enumerate the complete list of forbidden command patterns, or define a classification rule (e.g., "any shell command that executes a build/test/lint/run/dev action without going through `just`").
**What must improve**: Replace the illustrative examples with a verifiable rule: "Any shell command in a skill/agent/command file that invokes a build, test, lint, format, compile, run, or dev action must use `just <verb>` syntax. Commands in documentation comments, example blocks, or non-executable prose are exempt." This gives a reviewer an unambiguous pass/fail test.

### Attack 2: Risk Assessment — High/High risk has only a generic mitigation

**Where**: Risk 2: "Skill 行为变更可能导致现有工作流中断" rated at Likelihood=High, Impact=High, Mitigation="分两阶段部署：(1) 先更新 justfile 确保 recipe 存在，(2) 再更新 skill 引用；每个 skill 更新后立即运行其对应的 e2e 测试验证"
**Why it's weak**: A risk rated High/High is the top-priority risk in the entire proposal. Its mitigation is "phased deployment" -- a standard practice that would apply to any deployment of any change. This is not a risk-specific mitigation; it is a general deployment hygiene practice. For a High/High risk, the proposal should identify the specific failure mode (which skill breaks first? what test catches it?) and propose a targeted safety net. The mitigation mentions "运行其对应的 e2e 测试" but does not specify which e2e tests exist today, how many skills have e2e coverage, or what happens for skills that lack e2e tests. If only 3 of 8 skills have e2e tests, the mitigation has a 62% coverage gap that goes unmentioned.
**What must improve**: Enumerate current e2e test coverage per skill. Identify which skills lack e2e tests and state what manual verification will substitute. For the highest-risk skill (likely `execute-task` or `error-fixer` since they chain compile+test), propose a specific smoke test or rollback procedure.

### Attack 3: Solution Clarity — Scope resolution algorithm has an unreachable branch

**Where**: Skill integration section, step 2: "若任务无 scope 或 scope=all: 执行 `just project-type` / if output != 'mixed': 直接 `just build`（无 scope）/ else: `just build`（全部构建）"
**Why it's weak**: The else branch for mixed projects says `just build`（全部构建）, which is identical to the non-mixed branch's behavior (`just build`). This means for a mixed project where the task has no specific scope, the skill runs `just build` with no scope -- and the justfile is supposed to build everything. But the proposal also defines scope parameter behavior for mixed projects elsewhere: `just build frontend` / `just build backend` / `just build`. So step 2's mixed branch is correct but redundant -- it produces the same command as the non-mixed branch. More importantly, step 3 says "if project-type == 'mixed': just build frontend 或 just build backend" but does not explain what happens when the task scope does not match a valid project partition. What if a task is marked `scope=frontend` in a pure-backend project? The algorithm says "just build (无 scope)" which silently ignores the scope -- but this means the task intended frontend work on a backend-only project, which is likely an error in task breakdown, not a normal case. The algorithm has no error handling for this mismatch.
**What must improve**: Add an explicit handling for scope/project-type mismatch: "If task scope is `frontend` but project-type is `backend`, log a warning and proceed with `just build` (no scope), as the scope annotation is inapplicable." Remove the redundant else branch in step 2 or clarify that it is intentional fallback behavior. Add a note about how this mismatch might indicate a breakdown-tasks bug.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (iter 1): Success Criteria missing coverage for breakdown-tasks scope field | YES | New Criterion 7 added: "breakdown-tasks 生成的 index.json 中每个任务包含 scope 字段...正确分配 scope=frontend 和 scope=backend；跨端任务默认 scope=all" |
| Attack 2 (iter 1): Risk Assessment has no likelihood/impact ratings | PARTIAL | Likelihood and Impact columns now present with rated values (Medium/High). However, internal consistency is weak: Risk 2 is High/High with only a generic phased-deployment mitigation. Risk 4 was replaced (good) but the new Risk 4 (custom justfile overwrite) has two mitigations that partially overlap without acknowledging the overlap. |
| Attack 3 (iter 1): No explicit rationale for chosen alternative | YES | "Decision: 选择方案 A" paragraph added with specific failure mode analysis for Alternative B ("execute-task 调用 just compile, 而同一次任务流中 fix-bug 仍直接调用 go test") and cost breakdown for A ("16 处修改, 每处修改模式一致"). |

---

## Verdict

- **Score**: 85/100
- **Target**: 90/100
- **Gap**: 5 points
- **Action**: Continue to iteration 3 -- remaining gaps are: (1) Success Criteria criterion 3 needs a classification rule instead of illustrative examples, (2) Risk 2 (High/High) needs a targeted mitigation beyond generic phased deployment with current e2e coverage enumerated, (3) scope resolution algorithm needs error handling for scope/project-type mismatch and removal of redundant branch
