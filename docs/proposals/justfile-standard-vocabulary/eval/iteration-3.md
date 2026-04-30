---
date: "2026-04-29"
doc_dir: "docs/proposals/justfile-standard-vocabulary/"
iteration: "3"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 87/100** (target: 90)

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
| 5. Risk Assessment           |  13      |  15      | :white_check_mark:|
|    Risks identified (>=3)    |   5/5    |          |          |
|    Likelihood + impact rated |   4/5    |          |          |
|    Mitigations actionable    |   4/5    |          |          |
+------------------------------+----------+----------+----------+
| 6. Success Criteria          |  11      |  15      | :warning:|
|    Measurable                |   4/5    |          |          |
|    Coverage complete         |   4/5    |          |          |
|    Testable                  |   3/5    |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  87      |  100     |          |
+------------------------------+----------+----------+----------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem section | Evidence item "多个项目（如 pm-work-tracker）前后端混合" -- "多个" is a vague quantifier for a countable metric. How many projects? 2? 5? The reader cannot assess prevalence from this. | -1 pt |
| Problem section | Urgency is implied (skills break, commands don't exist) but never made explicit. No statement of "what breaks if we delay 3 months" or any deadline. The reader must infer urgency from evidence rather than being told. | -1 pt |
| Solution section | User-facing behavior is described as skill-internal logic (the 3-step scope resolution pseudocode). The proposal never states what the developer running a skill actually sees differently: same output? new output? different command invocations logged? The "user" here is the skill author, not the end developer, and that distinction is blurred. | -2 pts |
| Solution section | The Solution section does not independently justify why this approach is distinct from alternatives. The reader must cross-reference the Alternatives section to understand differentiators. The Solution section should stand on its own. | -1 pt |
| Alternatives section | Decision paragraph claims "每处修改模式一致，无架构复杂度" but this is asserted without evidence. Scope item 6 (adding scope field to breakdown-tasks index.json) is a fundamentally different kind of change than scope item 3 (replacing command strings in skill files). The claim of uniformity is an oversimplification. | -1 pt |
| Scope section | Scope item 6 (breakdown-tasks scope annotation) is a feature addition requiring prompt engineering / instruction modification, not a command replacement. It is qualitatively different from the other 6 scope items but is presented as equivalent work. The scope estimation treats all items as uniform, which understates the risk of the hardest item. | -1 pt |
| Risks section | Risk 4 likelihood ("Medium") is debatable -- any widely-adopted tool will eventually encounter users with custom justfiles who re-run init-justfile. The rating feels optimistic without justification for why it would only happen occasionally. | -1 pt |
| Risks section | Risk 4 has two mitigations (user confirmation prompt + boundary markers) but does not specify which is primary or how they interact. If the user declines the confirmation prompt, what happens? Do standard recipes stay stale? The mitigation is incomplete -- it handles "overwrite" but not "decline to overwrite." | -1 pt |
| Success Criteria | Criterion 1 says "生成匹配的 justfile" but "匹配" is not defined. Does "匹配" mean the justfile has all 16 commands? That scope parameters are present only for mixed projects? That the commands actually execute correctly? Criterion 5 partially covers this but only for the forge reference implementation, not for init-justfile's output in general. | -1 pt |
| Success Criteria | Criterion 7 bundles two distinct verifiable claims into one: (a) scope field exists with valid values, and (b) scope values are correctly assigned given specific inputs. These should be separate criteria so a partial failure can be identified precisely. | -1 pt |
| Success Criteria | Criterion 1 is the least testable: "能根据项目结构自动判断类型，生成匹配的 justfile" -- no specific test case or input/output pair is given. What project structure inputs should be tested? What does "matching" mean in terms of justfile content? A tester cannot write a test from this criterion alone. | -2 pts |

---

## Attack Points

### Attack 1: Success Criteria -- Criterion 1 is untestable as written

**Where**: Criterion 1: "init-justfile 能根据项目结构自动判断类型，生成匹配的 justfile（纯前端/纯后端/混合）"
**Why it's weak**: This criterion uses "匹配" (matching) without defining what a "matching" justfile contains. The word "匹配" could mean: (a) the justfile has the correct project-type recipe, (b) the justfile includes all 16 standard commands, (c) the justfile uses scope parameters only for mixed projects, (d) all generated recipes actually execute and return exit code 0, or any combination of these. A tester reading this criterion cannot write a pass/fail test because the expected output is undefined. Criterion 5 says the forge justfile must "包含全部 16 个标准命令" and be verifiable via `just project-type -> mixed`, but Criterion 1 provides no equivalent specificity for the general case of init-justfile output. The criterion also says "自动判断类型" but does not define what inputs trigger each type -- what directory structure makes a project "frontend" vs "mixed"? This is critical because the entire scope-resolution algorithm depends on `project-type` being correct.
**What must improve**: Replace Criterion 1 with a concrete test case: "Given a project with `package.json` and no `go.mod`, init-justfile generates a frontend-style justfile (no scope parameters, `just project-type` outputs `frontend`). Given a project with `go.mod` and `package.json` in a subdirectory, init-justfile generates a mixed-style justfile (scope parameters, `just project-type` outputs `mixed`)." Define "matching" explicitly by listing the verifiable properties of each justfile variant.

### Attack 2: Solution Clarity -- User-facing behavior is never described from the developer's perspective

**Where**: Solution section, Skill integration subsection (lines 87-100): the 3-step algorithm describes skill-internal decision logic, not what the developer experiences
**Why it's weak**: The entire Solution section describes internal mechanisms (scope resolution algorithm, project-type probing, scope field in index.json) rather than the observable developer experience. A developer using forge skills after this change would never see the scope resolution algorithm -- they would see commands being executed in their terminal. The proposal never answers: What does the developer see that is different from today? Do they see `just build frontend` instead of `go build`? Do they see a warning when scope mismatches? Does `just project-type` output appear in their logs? The pseudocode block (lines 87-100) is written for the skill author, not the end user. This means a stakeholder who is not a skill implementer cannot understand what this proposal delivers to them. The one user-facing element (the "scope 与项目类型不匹配" warning on line 97-99) is buried inside algorithm pseudocode and described as "记录警告" with no specification of where the warning appears or who reads it.
**What must improve**: Add a "Developer Experience" subsection that describes what a developer sees before and after this proposal. For example: "Before: A developer runs execute-task and sees `go test ./...` in the terminal. After: The developer sees `just test backend` (for a backend-scoped task in a mixed project) or `just test` (for a single-type project). If the task scope mismatches the project type, the developer sees a warning: `[forge] scope=frontend but project-type=backend; falling back to just build`."

### Attack 3: Alternatives Analysis -- Uniformity claim masks heterogeneous work items

**Where**: Decision paragraph: "方案 A 需更新 10 个 skill/agent/command 文件 + 4 个 task 模板 + 1 个 breakdown-tasks skill + 1 个 init-justfile 契约 = 共 16 处修改。每处修改模式一致（将原始命令替换为 `just <verb>`），无架构复杂度"
**Why it's weak**: The claim "每处修改模式一致" is factually inaccurate. The 16 modifications fall into at least three distinct categories: (1) Command string replacement in skill/agent/command files (items 3-5 in scope) -- mechanical find-and-replace. (2) Adding adaptive generation logic to init-justfile (item 1) -- this requires new code to detect project structure and generate different justfile variants, which is non-trivial. (3) Modifying breakdown-tasks to add scope annotations to index.json (item 6) -- this requires changing the AI agent's prompt/instructions to teach it how to classify tasks by scope, which is prompt engineering work with inherently uncertain outcomes. The proposal treats these three categories as equivalent "16 处修改" and concludes "无架构复杂度." But item 6 (breakdown-tasks scope annotation) and item 1 (adaptive justfile generation) are qualitatively harder and riskier than the other 14 items. By flattening all 16 items into a single cost estimate, the proposal underestimates the actual complexity and makes the "controllable deterministic work" claim feel like an oversimplification rather than an honest assessment.
**What must improve**: Break the 16 modifications into categories by difficulty/risk. Acknowledge that items 1 and 6 are not simple command replacements. Provide a rough time estimate or complexity assessment per category. For example: "14 items are mechanical command replacements (~15 min each). Item 1 (init-justfile adaptive generation) requires new detection logic and template system (~1-2 days). Item 6 (breakdown-tasks scope annotation) requires prompt engineering and validation (~1 day)." This gives the decision-maker an honest cost picture rather than a misleading uniformity claim.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (iter 2): Success Criteria criterion 3 boundary is still fuzzy despite iteration 1 fix | YES | Criterion 3 now reads: "任何可执行指令（Bash tool call 或代码块）若执行 build、test、lint、format、compile、run、dev、clean、install 操作，必须通过 `just <verb>` 调用。直接调用语言工具链（`go`、`npm`、`cargo`、`pytest`、`npx` 等）为违规。文档注释、示例说明、非可执行文本中的命令不受此规则约束". This replaces the illustrative examples with a classification rule and explicit exemption scope. |
| Attack 2 (iter 2): High/High risk has only generic mitigation | YES | Risk 2 mitigation now enumerates current e2e coverage gap explicitly: "当前 e2e 覆盖：`justfile-e2e-integration` 测试套件验证 skill 文件内容匹配（20/20 通过），但本次修改的 8 个 skill 文件和 4 个 task 模板均无独立的 per-skill e2e 测试". It provides a specific manual verification procedure: "对无 e2e 覆盖的 skill，逐个手动验证：执行更新后的 skill 文件中引用的每个 `just <verb>` 命令，确认 justfile 中对应 recipe 存在且返回码为 0；若执行失败，回滚该 skill 文件至上一版本并记录失败 recipe 名称". This is a targeted, actionable mitigation. |
| Attack 3 (iter 2): Scope resolution algorithm has an unreachable branch and no error handling for mismatch | PARTIAL | The algorithm now includes mismatch handling at line 97: "scope 与项目类型不匹配（如 scope=frontend 但项目为纯 backend）记录警告：可能是 breakdown-tasks scope 分配错误 → 回退执行 `just build`（无 scope 参数）". This addresses the error handling gap. However, the redundant branch in step 2 is still present: for both mixed and non-mixed projects where scope=all or scope is absent, the algorithm produces the identical command `just build` with no scope parameter. The behavioral difference (mixed project builds everything, single-type project builds everything) is identical, making the mixed-specific branch in step 2 redundant. This was not cleaned up. |

---

## Verdict

- **Score**: 87/100
- **Target**: 90/100
- **Gap**: 3 points
- **Action**: Continue to iteration 4 -- remaining gaps are: (1) Criterion 1 needs concrete test cases defining what "matching" means with specific project-structure inputs and expected justfile outputs, (2) Solution section needs a Developer Experience subsection describing what the end user observes before/after, (3) Alternatives analysis must acknowledge that the 16 modifications are not uniform in complexity and break them into honest difficulty categories
