# CTO Adversarial Evaluation Report

**Document**: `docs/proposals/skill-instruction-audit/proposal.md`
**Reviewer**: CTO (Proposal Expert)
**Iteration**: 1 (baseline)
**Date**: 2026-05-28

---

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

The proposal identifies three categories of systemic defects in Forge plugin skill instruction files: CLI behavior descriptions (22 instances), redundancy (33 instances), and clarity/self-consistency issues (40 instances). The solution directly addresses all three categories with deletion, simplification, and correction respectively. The mapping is straightforward and the solution is not solving a different problem.

However, the core solution mechanism -- "delete CLI behavioral descriptions, keep imperative instructions" -- relies on a binary distinction that does not exist in practice. The proposal itself acknowledges the boundary is fuzzy in its risk mitigation: "保留 exit code 契约（0=成功/1=失败）和输出字段名列表，只删除语义解释". But "output field name lists" and "behavioral descriptions" often coexist in the same paragraph. The proposal provides no classification rule, leaving each task executor to make independent judgment calls across 40 files.

### Evidence -> Solution Trace

The evidence section provides convincing pattern-level examples for all three defect categories. However, individual instance evidence is uneven. Some claims lack sufficient specificity for verification:
- "gen-contracts/SKILL.md 引用不存在的 Section 编号" -- no section number, no line reference, no quote.
- "quick.md 配置读取失败时 fallback 为跳过确认门（逻辑倒置）" -- the source file contains an explicit design justification ("This preserves quick mode's streamlined nature") that the proposal does not acknowledge or refute.

### Evidence -> Success Criteria Trace

Success criteria are predominantly structural (grep checks, step number verification). The grep pattern "What .* Does" for SC-1 only catches sections with that specific title format. Inline behavioral descriptions within step instructions -- the harder and more numerous category identified by the proposal's own evidence (execute-task.md Step 1.5 describing MAIN_SESSION routing, quick-tasks/SKILL.md explaining what forge task index generates) -- have no verification method.

### Self-Contradiction Check

1. **Cross-file duplication as evidence vs. out-of-scope**: The proposal cites "quick-tasks <-> breakdown-tasks 之间 12 处近乎逐字重复" as evidence of the redundancy problem, then explicitly scopes out cross-file deduplication. If the 12 duplications are evidence of a defect, leaving them unfixed means the defect persists after the proposal is "complete."

2. **SC-5 vs. documented design intent**: SC-5 proposes reversing the quick.md fallback from "skip gate" to "show gate." The source file documents this as an intentional design choice: "This preserves quick mode's streamlined nature." The proposal treats it as an unambiguous bug without acknowledging the existing design rationale.

3. **SC-1 verification gap**: SC-1 verifies CLI description deletion via `grep 无 "What .* Does" section`. But In Scope item 1 claims 22 instances of CLI behavioral descriptions, and the proposal's own evidence cites inline descriptions that would not be caught by this grep. The success criterion tests a subset of what the scope promises.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

**Problem stated clearly (35/40)**: The three-category classification (CLI behavior descriptions, redundancy, clarity/self-consistency) is clear and unambiguous. Two readers would likely converge on the same interpretation. However, the boundary between "CLI behavioral description" and "output contract" is not defined, which creates ambiguity in what exactly constitutes the problem.

**Evidence provided (30/40)**: Pattern-level evidence is strong -- specific file names and specific patterns are cited. But individual instance evidence is uneven. The gen-contracts section reference lacks specificity. The quick.md fallback characterization omits the source file's documented design justification. The "约 22 处", "约 33 处", "约 40 处" counts suggest precision but the "约" qualifier undermines it.

**Urgency justified (25/30)**: The v3.0.0 release context provides a concrete deadline. The cost of delay is quantified in terms of agent execution cycles (30-minute subagent timeout per failure). However, no data is provided on actual failure rates -- how often do agents actually fail due to these defects? Without this, the urgency claim rests on assumption rather than measurement.

**Subtotal: 90/110**

### 2. Solution Clarity (120 pts)

**Approach is concrete (35/40)**: The three-batch approach (delete CLI descriptions, remove redundancy, fix clarity issues) is concrete enough to explain back. A reader can describe what will be built. The per-batch description is specific. Deduction for not defining the deletion boundary between "behavioral description" and "output contract."

**User-facing behavior described (35/45)**: The proposal correctly states that external behavior must not change ("修改不能改变任何 skill 的外部行为"). But it does not describe what the agent's reading experience will be after the changes -- how will the instruction files read differently? A "before/after" example for at least one file would make this concrete. The proposal describes what gets removed but not what the result looks like.

**Technical direction clear (30/35)**: The general approach (text deletion, simplification, renumbering) is clear. The innovation highlight section provides the guiding principle ("imperative, not descriptive"). However, the absence of a classification rule for the CLI deletion boundary means the technical direction is clear at the macro level but ambiguous at the execution level.

**Subtotal: 100/120**

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (10/40)**: No industry solutions, open-source projects, or published patterns are cited. The proposal references only self-generated alternatives. For a proposal about AI agent instruction design, there is relevant published work on prompt engineering, instruction tuning, and LLM-optimized documentation that could have been referenced.

**At least 3 meaningful alternatives (20/30)**: Three alternatives are presented (do nothing, batch by issue type, batch by file). "Do nothing" is the required baseline. However, both "batch by issue type" and "batch by file" are nearly identical approaches -- both are about task organization strategy for the same set of changes. They are not genuinely different solution approaches to the underlying problem. A genuinely different alternative would be: automated linting/validation for instruction files, a schema-based approach to enforce imperative-only instructions, or a template system that separates output contracts from behavioral descriptions.

**Honest trade-off comparison (10/25)**: The trade-off analysis is shallow. "跨文件上下文切换多" for the selected approach and "同类问题在不同文件中修复方式可能不一致" for the rejected approach are both minor operational concerns. No comparison of actual solution quality differences between approaches.

**Chosen approach justified against benchmarks (5/25)**: No industry benchmarks are cited, so no justification against benchmarks is possible. The selected approach is justified only against the two self-generated alternatives.

**Subtotal: 45/120**

### 4. Requirements Completeness (110 pts)

**Scenario coverage (30/40)**: Happy path (agent reads and executes correctly after fix) is covered. Edge cases are partially addressed: the risk table identifies the edge case of over-deletion (losing necessary context). However, the scenario of "agent encounters a modified file and interprets the remaining text differently than intended because a behavioral description that provided implicit context was removed" is not explicitly identified. Error scenarios (what if a task executor misclassifies a boundary case?) are not addressed.

**Non-functional requirements (25/40)**: The constraint "修改不能改变任何 skill 的外部行为（输入/输出/副作用）" is a correctness requirement, but no verification method is proposed for it. No performance, compatibility, or accessibility NFRs are discussed. For an instruction-layer change, backward compatibility with existing agent workflows is a relevant NFR that is not addressed.

**Constraints & dependencies (25/30)**: Four explicit constraints are listed, all concrete. The forge-distribution.md path convention dependency is correctly identified. However, the constraint on "quick-tasks 和 breakdown-tasks 必须各自独立自洽" conflates structural independence with textual duplication, as noted in the reasoning audit.

**Subtotal: 80/110**

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (20/40)**: The core insight -- "AI agent instructions should be imperative, not descriptive" -- is not novel in the prompt engineering domain. This is a well-established principle in prompt engineering literature. The proposal applies it to a specific codebase (Forge plugin skills), which is valuable work but not a creative leap beyond the industry baseline.

**Cross-domain inspiration (10/35)**: No cross-domain inspiration is evident. The proposal does not reference how other domains (API documentation, technical writing for LLMs, prompt engineering research) handle the tension between descriptive context and imperative instructions.

**Simplicity of insight (20/25)**: The insight is genuinely simple and elegant: "for AI, repetition increases inconsistency risk rather than reinforcing memory." This is well-articulated and has the "why didn't I think of that" quality. The inversion of human documentation best practices for AI consumers is cleanly stated.

**Subtotal: 50/100**

### 6. Feasibility (100 pts)

**Technical feasibility (35/40)**: The assessment is accurate -- these are text modifications with no code changes. The claim "纯文本修改，无代码变更，无依赖风险" is verifiably true. However, the assessment understates the difficulty of the judgment calls required: distinguishing "behavioral description" from "output contract" across 22 instances requires domain knowledge of each skill's purpose, not just text editing.

**Resource & timeline feasibility (25/30)**: "约 95 处修改分布在 ~40 个文件中。预计 8-12 个 coding task" is reasonable for pure text changes. However, no estimate is provided for the review effort required to verify that the judgment calls (especially the CLI deletion boundary) were applied correctly across all files.

**Dependency readiness (25/30)**: No external dependencies are claimed, which is correct. The only dependency is forge-distribution.md conventions, which already exist. Deduction for not mentioning that the proposal's success depends on each task executor having sufficient context about each skill's purpose to make correct deletion/simplification decisions.

**Subtotal: 85/100**

### 7. Scope Definition (80 pts)

**In-scope items are concrete (25/30)**: All five in-scope items are specific and deliverable. "删除 22 处 CLI 行为描述" is more concrete than "improve instruction quality." However, "确保 quick-tasks 和 breakdown-tasks 各自内部自洽" is somewhat vague -- what constitutes "self-consistency" and how is it verified?

**Out-of-scope explicitly listed (20/25)**: Six out-of-scope items are named. The most important exclusion (cross-file deduplication) is explicitly listed with justification. However, the justification creates an internal tension: the proposal claims cross-file duplication as evidence of the problem, then excludes fixing it. The "why" for this exclusion could be stronger.

**Scope is bounded (20/25)**: The scope is bounded by the three defect categories and the 95-modification estimate. The "8-12 coding tasks" estimate provides a timeframe proxy. However, the scope is open-ended in one dimension: "修复 40 处清晰度/自洽性问题" -- without specific line references for all 40 instances, the actual scope of this item could expand during execution.

**Subtotal: 65/80**

### 8. Risk Assessment (90 pts)

**Risks identified (20/30)**: Three risks are identified. The first two (losing necessary context, missing constraints after E-I simplification) are meaningful. The third (renumbering errors) is trivial by comparison -- it is a mechanical check with minimal consequence. Missing risks: the risk of inconsistent judgment calls across task executors on the CLI deletion boundary; the risk that the quick.md fallback change breaks existing workflows that depend on the current behavior; the risk that the 40 clarity fixes introduce new clarity issues if not reviewed holistically.

**Likelihood + impact rated (20/30)**: The ratings are reasonable but lean toward understatement. "删除 CLI 描述后 agent 丢失必要上下文" is rated L/M -- given that the boundary between "behavioral description" and "output contract" is undefined, a more honest assessment would be M/H. All three risks have relatively safe ratings; none are rated H/H.

**Mitigations are actionable (20/30)**: The mitigations are partially actionable. "保留 exit code 契约... 和输出字段名列表" is concrete. "逐一与正文步骤比对，确保无遗漏" is a process description but lacks a specific verification method (what constitutes "no omission"?). "重排后验证 Process Flow 与实际步骤编号一一对应" is actionable.

**Subtotal: 60/90**

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (20/30)**: SC-1 (grep check), SC-3 (step sequence), SC-4 (absence of specific string), SC-5 (fallback behavior), SC-6 (frontmatter check) are measurable. SC-2 ("每个 E-I 块条目在正文中无对应") is ambiguous -- "对应" is undefined. SC-7 ("agent 只读其一即可执行") is not measurable without defining what "complete independent instruction set" means operationally.

**Coverage is complete (18/25)**: SC-1 through SC-7 cover all three defect categories. However, the 40 clarity/self-consistency fixes are covered by only 2 specific SCs (SC-3, SC-4) plus a generic SC-7. The remaining ~38 clarity issues have no specific verification criteria. The E-I deduplication (33 instances) is covered by a single SC-2 with an ambiguous verification method.

**SC internal consistency (15/25)**:

SC-1 (grep 无 "What .* Does" section) vs. InScope-1 (删除 22 处 CLI 行为描述): SC-1 only verifies a subset of InScope-1. Satisfying SC-1 does not guarantee InScope-1 is fully met, because inline behavioral descriptions are not caught by the grep pattern. This is a coverage gap, not a contradiction, but it means SC-1 is a weak proxy for InScope-1 completion.

SC-2 (E-I items have no body correspondence) vs. Risk-2 mitigation (verify no constraints lost): These are in tension. SC-2 says "remove if body has correspondence." Risk-2 mitigation says "ensure no constraints are lost." An E-I item that corresponds to a body entry at a lower enforcement level would be removed by SC-2 but should be preserved by Risk-2 mitigation. The direction of satisfying SC-2 could violate the Risk-2 mitigation goal.

**Subtotal: 53/80**

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (30/35)**: The solution directly addresses all three defect categories. However, the "quick.md fallback" fix does not address a defect in the instruction layer -- it addresses a behavioral design choice. The proposal classifies it under "clarity/self-consistency" but it is actually a design reversal, which is a different category of change than what the problem statement describes.

**Scope <-> Solution <-> Success Criteria aligned (18/30)**: Cross-section alignment issues:
- InScope-1 claims 22 CLI description deletions; SC-1 only verifies "What .* Does" sections, which is a subset. The scope promise and the success criteria are misaligned.
- InScope-2 claims 33 redundancy fixes; SC-2 uses an ambiguous "无对应" criterion that may over-delete or under-delete.
- InScope-3 claims 40 clarity fixes; only 2 specific SCs exist for these.
- The 12 cross-file duplications are cited as evidence but excluded from scope. This means the "redundancy" problem is only partially addressed.

**Requirements <-> Solution coherent (20/25)**: Most requirements map to solution elements. The constraint "不抽取共享 rule 文件" is respected. The "不改变外部行为" constraint is stated but has no corresponding verification in the success criteria. The requirement for "新 skill 编写时，有清晰的范式可遵循" is not backed by any in-scope deliverable -- the proposal fixes existing files but does not produce a paradigm document or template.

**Subtotal: 68/90**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Undefined Deletion Boundary -- The Fundamental Execution Risk

The proposal's entire execution depends on a classification that is never defined: what distinguishes a "CLI behavioral description" (to delete) from an "output contract" (to preserve). The proposal states the principle ("指令性的，不是描述性的") but provides no classification rule. This is not covered by any single rubric dimension -- it is a prerequisite for the document to be executable.

> "删除对 CLI 输出语义、内部实现、分支逻辑的解释" + "保留 exit code 契约（0=成功/1=失败）和输出字段名列表，只删除语义解释"

The phrase "CLI 输出语义" in the solution and "语义解释" in the risk mitigation are contradictory: the solution says delete "输出语义", the mitigation says preserve output field name lists. But "field name + what it contains" IS output semantics. No classification rule resolves this contradiction. Every task executor will apply a different standard.

**Resolution required**: Add a "Deletion Boundary Rules" subsection with three categories: (1) imperative instructions (keep), (2) output contracts with field names, types, and absence semantics (keep), (3) behavioral explanations of internal command logic (delete). Provide 2-3 examples per category from actual files.

### [blindspot-2] quick.md Fallback -- Bug Claim Without Acknowledging Documented Design Intent

The proposal treats the quick.md config-read-failure fallback as an unambiguous bug ("逻辑倒置"). However, the source file documents this as an intentional design choice with explicit reasoning. The proposal does not acknowledge this design justification, let alone argue against it. This is not a dimension-scoring issue -- it is a reasoning flaw where the proposal presents a judgment call as a factual error.

> "quick.md 配置读取失败时 fallback 为跳过确认门（逻辑倒置）" -- vs. source file: "This preserves quick mode's streamlined nature."

**Resolution required**: In the Evidence section, acknowledge the source file's design justification and provide an explicit argument for why it is wrong. Distinguish between config-missing (new user, should confirm) and config-corrupted (should fail-safe). Currently the proposal reverses a documented design choice without engaging with the design rationale.

### [blindspot-3] Cross-file Duplication Evidence vs. Scope Exclusion -- The Internal Tension

The proposal uses 12 cross-file duplications as evidence of the redundancy problem, then excludes fixing them from scope. This means the proposal's own evidence demonstrates a problem that the proposal will not solve. This is a structural misalignment between problem framing and scope that no rubric dimension directly captures.

> Evidence: "quick-tasks <-> breakdown-tasks 之间 12 处近乎逐字重复" + Out of Scope: "跨文件去重（quick-tasks<->breakdown-tasks、execute-task<->run-tasks 保持独立）"

**Resolution required**: Either (a) remove the cross-file duplication from the evidence (it is not a defect this proposal addresses), or (b) add lightweight cross-file deduplication to scope with the "same rule, wording appropriate to context" standard. The current state uses evidence that the scope disclaims.

---

## Dimension Breakdown (final)

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 90 | 110 |
| 2. Solution Clarity | 100 | 120 |
| 3. Industry Benchmarking | 45 | 120 |
| 4. Requirements Completeness | 80 | 110 |
| 5. Solution Creativity | 50 | 100 |
| 6. Feasibility | 85 | 100 |
| 7. Scope Definition | 65 | 80 |
| 8. Risk Assessment | 60 | 90 |
| 9. Success Criteria | 53 | 80 |
| 10. Logical Consistency | 68 | 90 |
| **Total** | **696** | **1000** |

### Outcome

Target NOT reached (696/1000, target: 900). Primary gaps: Industry Benchmarking (severely underspecified), Success Criteria (ambiguous verification methods), Risk Assessment (incomplete risk identification), Logical Consistency (scope-evidence-SC misalignment).
