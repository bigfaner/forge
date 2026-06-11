# Baseline Evaluation Report

**Proposal**: Agent-Driven Justfile Generation
**Evaluator**: CTO-level Adversarial Review
**Date**: 2026-06-08
**Rubric**: `plugins/forge/skills/eval/rubrics/proposal.md` (1000-point scale)

---

## Phase 1 -- Reasoning Audit

### Argument Chain Trace

| Link | Claim | Verdict |
|------|-------|---------|
| Problem | Language templates are a flexibility bottleneck; non-standard projects must manually edit generated justfiles | Valid. Verified: 6 templates, 811 lines, mixed.just = 229 lines. The three-layer process (template -> Convention -> LLM) does create the described conflict where template defaults get entirely replaced. |
| Solution | Remove templates; let agent generate commands from surface rules + project detection | Plausible but underspecified. The solution asserts agent knowledge is sufficient but provides no evidence of this claim beyond assertion. |
| Evidence | 6 templates ~800 lines; mixed.just ~230 lines; template/LLM conflict | Verified accurate against codebase. |
| Success Criteria | Checklist of 7 items covering deletion, creation, and verification | Adequate for a refactor of this scope, but has gaps (see D9). |

### Self-Contradiction Check

1. **NFR "consistency" vs LLM generation**: The proposal states "相同项目多次运行生成的 justfile 结构一致（recipe 名称、分组、边界标记不变；具体命令可能因 LLM 变化而有细微差异）". This admits non-determinism in command generation while claiming "向后兼容" with current output. These two statements are in tension -- backward compatibility demands identical output, but LLM non-determinism guarantees deviation. The proposal handwaves this as "细微差异" without quantifying what level of deviation is acceptable.

2. **"dependency readiness" vs actual work needed**: Claims "surface rule 文件已定义完整的 recipe 契约（已有 5 个）" but the proposal also lists in scope "简化 5 个 surface rule 文件：...替换 TODO stub 模板为 Recipe Generation Requirements section". If the rules are "ready" as dependencies, why do they need modification as deliverables? This is a circular dependency claim.

---

## Phase 2 -- Rubric Scoring

### D1. Problem Definition (110 pts)

**Problem stated clearly (36/40)**: The core problem -- template rigidity -- is unambiguous. However, the framing is narrow: it describes a mechanical issue (templates don't fit non-standard projects) without quantifying the user pain. How many projects hit this? What percentage of runs require manual edits? Quote: "非标准项目生成后必须手动编辑 justfile" -- "must" is strong but no data backs this frequency claim.

**Evidence provided (32/40)**: Concrete codebase metrics (6 templates, ~800 lines, mixed.just ~230 lines) are strong. The three-layer conflict observation ("模板给出的默认命令与实际项目不匹配，LLM 需要全部替换") is a real architectural insight. Deduction: no user feedback, no issue tracker references, no frequency data. Evidence is entirely structural/code-level, not user-validated.

**Urgency justified (22/30)**: Quote: "移除模板后，agent 直接针对实际项目结构生成正确命令，消除这层摩擦。" This argues benefits but not urgency -- why now? What happens if delayed? The cost of delay is implicit (template maintenance continues) but never quantified. No deadline, no competing priority, no user complaint escalation.

**D1 Total: 90/110**

---

### D2. Solution Clarity (120 pts)

**Approach is concrete (32/40)**: The 3-step process (enumerate recipes from surface rules -> fill content from detection + knowledge -> extract server lifecycle patterns) is clear. A reader can explain what will be built. Deduction: Step 2 "agent 根据检测到的语言/框架 + Convention 知识 + 自身知识生成每个 recipe 的具体命令" is the core mechanism and it's vague -- what does "自身知识" mean operationally? How is the prompt structured? What context is fed to the LLM?

**User-facing behavior described (35/45)**: The proposal describes what changes for the developer (no templates, agent generates directly), but user-facing behavior is described only by omission (what goes away). The observable behavior of running `/init-justfile` is not described -- does the UX change? Is there a new step? Is the output identical? Quote: "向后兼容：生成的 justfile 结构...与当前模板输出一致" -- this claims identical UX but the NFR section admits LLM variability.

**Technical direction clear (28/35)**: The surface-rule-driven approach is clear. Server lifecycle extraction to a rule file is concrete. Deduction: the proposal doesn't address how agent knowledge is operationalized. Current SKILL.md has detailed HARD-RULE blocks with three-layer generation; the proposal says "agent 根据检测到的语言/框架 + Convention 知识 + 自身知识" but doesn't specify what replaces the HARD-RULE template-loading mechanism.

**D2 Total: 95/120**

---

### D3. Industry Benchmarking (120 pts)

**Industry solutions referenced (28/40)**: Quote: "Make/Just 脚手架工具普遍使用模板（template -> fill -> output）。更现代的工具（如 Earthly、Taskfile）通过 DSL 抽象部分命令". Earthly and Taskfile are named, but only in a single sentence with no analysis of how they solve this problem. Cookiecutter, yeoman, and hygen are mentioned in the Innovation section but not analyzed here. No URLs, no version references, no published patterns cited. This is name-dropping, not benchmarking.

**At least 3 meaningful alternatives (22/30)**: Three alternatives presented: (1) do nothing, (2) parameterized templates, (3) agent-driven. "Do nothing" is required. "Parameterized templates" is a genuinely different approach. "Agent-driven" is the proposal itself. However, "parameterized templates" is treated as a straw man -- the cons column says "仍需维护模板，复杂度转移为参数爆炸" which dismisses it without analysis. A genuinely different third alternative (e.g., DSL-based generation, hybrid template+agent, community-contributed templates) is missing.

**Honest trade-off comparison (18/25)**: The comparison table has "LLM 生成有轻微不确定性" as the sole con for the selected approach. This undersells the risk. No mention of: increased token cost, longer generation time, testing difficulty, regression risk from non-determinism. The "parameterized templates" row has "复杂度转移为参数爆炸" without evidence.

**Chosen approach justified against benchmarks (18/25)**: Quote: "彻底解决灵活性问题". The justification is that it eliminates template maintenance. But no benchmark was analyzed in depth -- the industry tools (Earthly, Taskfile) were mentioned in passing, not compared. The rationale is internal (our surface rules + agent knowledge), not grounded in external validation.

**D3 Total: 86/120**

---

### D4. Requirements Completeness (110 pts)

**Scenario coverage (30/40)**: Five key scenarios listed: single surface, multi-surface named, mixed language, non-standard structure, empty surfaces. Good coverage of happy paths and one edge case (empty surfaces). Missing scenarios: (1) existing justfile with user customizations -- the current SKILL.md has extensive handling for this but the proposal doesn't mention it as a scenario; (2) migration from template-generated justfile to agent-generated -- what happens to existing projects? (3) cold start with no Convention AND no standard language detection.

**Non-functional requirements (30/40)**: Two NFRs listed: consistency and backward compatibility. Both are relevant. Deduction: missing NFRs -- (1) performance/latency: removing templates means every generation is an LLM call, how much slower? (2) determinism: the consistency NFR admits variability but doesn't set an acceptable threshold. (3) token cost: agent-driven generation consumes more tokens than template filling.

**Constraints & dependencies (25/30)**: Four dependencies listed, all validated (forge surfaces command, surface rules, Convention mechanism, just >= 1.50.0). Good. Deduction: the proposal claims "所有依赖已就绪" but also lists surface rule modification as in-scope work -- if rules need modification, they're not ready as-is.

**D4 Total: 85/110**

---

### D5. Solution Creativity (100 pts)

**Novelty over industry baseline (32/40)**: The proposal explicitly positions itself against template-driven generation. Quote: "传统脚手架工具（cookiecutter、yeoman、hygen）使用模板驱动生成。本方案用 LLM 驱动 + 结构约束 替代模板". The differentiation is clear and honest: LLM-driven with structural constraints vs template-driven. This is a genuine innovation in the context of build scaffolding tools.

**Cross-domain inspiration (20/35)**: The idea of "contract-driven LLM generation" (surface rules define "what", agent decides "how") borrows from API contract patterns, but the proposal doesn't articulate this cross-domain connection. No reference to similar patterns in other domains (e.g., OpenAPI codegen, GraphQL schema-driven resolvers, infrastructure-as-code). The insight is good but the inspiration trail is missing.

**Simplicity of insight (20/25)**: The core insight -- "agent already knows how to generate commands, templates are an unnecessary indirection" -- is elegant. It follows the "why didn't I think of that" pattern. The Occam's Razor argument in the Assumptions Challenged table is well-made. Deduction: the elegance is slightly undercut by the complexity of what replaces templates (surface rules + Convention + agent knowledge synthesis) -- it's simpler in code but not simpler in mental model.

**D5 Total: 72/100**

---

### D6. Feasibility (100 pts)

**Technical feasibility (32/40)**: Quote: "Agent（Claude/GPT）已具备主流语言构建命令的知识". This is a reasonable claim for mainstream languages (Go, Node, Python, Rust). However, the assertion is untested in this specific context -- has anyone tried generating full justfile recipes via agent without templates? The current SKILL.md shows a complex multi-step process with careful HARD-RULE blocks; replacing this with "agent knowledge" is a significant trust leap. The proposal doesn't cite any proof-of-concept.

**Resource & timeline feasibility (25/30)**: Quote: "改动集中在 init-justfile skill 目录内...预计单个 skill 改动，可在一次 session 内完成". The scope is bounded and realistic for a single skill refactor. The file list (delete 6+1, add 1, rewrite 1, modify 5) is concrete. Deduction: "一次 session" is optimistic -- rewriting SKILL.md (currently 490+ lines) while maintaining all the existing behavior guarantees is substantial.

**Dependency readiness (22/30)**: Claims all dependencies are ready. The `forge surfaces` command is stable. But the proposal says surface rule files need modification (in scope item: "简化 5 个 surface rule 文件"), which means they're NOT ready as-is -- they need to be changed to replace "TODO stub 模板" with "Recipe Generation Requirements" sections. This is a deliverable, not a dependency.

**D6 Total: 79/100**

---

### D7. Scope Definition (80 pts)

**In-scope items are concrete (26/30)**: Five concrete deliverables: delete 6 templates, delete 1 rule, add 1 rule, rewrite SKILL.md, simplify 5 surface rules. Each is a file-level action. Good. Deduction: "简化 5 个 surface rule 文件：保留编排序列/recipe 契约/journey 策略，替换 TODO stub 模板为 Recipe Generation Requirements section" -- "Recipe Generation Requirements section" is an abstraction, not a deliverable. What does this section contain?

**Out-of-scope explicitly listed (22/25)**: Four items explicitly excluded: CLI surfaces detect command changes, new surface types, Convention loading changes, other skill updates. Good. Deduction: missing from out-of-scope -- what about testing strategy? Is updating/creating tests in scope? What about documentation?

**Scope is bounded (20/25)**: The scope is contained within a single skill directory, which is strongly bounded. "一次 session" timeframe is clear. Deduction: "简化 5 个 surface rule 文件" is potentially unbounded -- "simplify" could mean anything from a few line changes to complete rewrites.

**D7 Total: 68/80**

---

### D8. Risk Assessment (90 pts)

**Risks identified (24/30)**: Four risks identified: LLM inconsistency, server lifecycle complexity, cold start quality, rare language failure. These are meaningful risks. Deduction: missing risks -- (1) regression risk: existing projects with template-generated justfiles that re-run init-justfile get different output; (2) token cost / latency increase from LLM-driven generation; (3) prompt engineering risk: the agent needs specific context to generate correct commands -- what context window is required?; (4) observability: how do you debug a bad generation?

**Likelihood + impact rated (22/30)**: Ratings seem honest: two M/M, one L/M, one L/L. Not everything is "low likelihood, high impact". Deduction: "LLM 生成命令不一致" is rated M likelihood, L impact. But the NFR says "具体命令可能因 LLM 变化而有细微差异" -- if this is a known behavior with M likelihood, the impact on CI/CD pipelines (where exact commands matter) should be rated more explicitly. The "轻微不确定性" framing in the comparison table understates what's rated M/L here.

**Mitigations are actionable (20/30)**: Mitigations reference existing mechanisms: "Surface rule 的 recipe 契约 + Convention 约束", "提取为独立 rule 文件", "agent 具备主流语言的默认知识". These are somewhat actionable. Deduction: "verification step (dry-run + actual) 捕获错误" is the primary mitigation for the top risk, but this is an existing mechanism, not a new mitigation. The proposal doesn't add any new safeguards beyond what already exists. Quote for cold start: "agent 具备主流语言的默认知识；verification step 捕获错误并自修正" -- this is "the LLM will handle it", not an actionable mitigation.

**D8 Total: 66/90**

---

### D9. Success Criteria (80 pts)

**Criteria are measurable and testable (22/30)**: Seven criteria, all checkable via file existence, command execution, or structural comparison. "agent 生成的 recipe 通过 verification step（dry-run + actual execution）" is testable. Deduction: "生成的 justfile 结构...与当前输出一致" is vague -- what does "一致" mean given the admitted LLM variability? Is it structural identity (same recipe names, groups, markers) or byte-level identity? This is not fully measurable without a threshold.

**Coverage is complete (18/25)**: Criteria cover: template deletion (SC2), rule creation (SC2-3), surface rule modification (SC1), empty surface handling (SC4), language-specific generation (SC5), multi-surface generation (SC6), structural compatibility (SC7). Missing from coverage: (1) no criterion for the SKILL.md rewrite -- the most complex deliverable; (2) no criterion for `rules/project-detection.md` deletion; (3) no performance criterion (generation should not take significantly longer).

**SC internal consistency (18/25)**: Clustering SC by affected area:
- **Templates**: SC2 (delete) -- satisfiable, no conflict.
- **Surface rules**: SC1 (modify 5 files) -- satisfiable.
- **SKILL.md**: SC4 (empty surface prompt) -- but no SC for the SKILL.md rewrite itself.
- **Server lifecycle**: SC3 (create rule file) -- satisfiable.
- **Output quality**: SC5 (Go/Node/Python/Rust pass verification), SC6 (mixed multi-surface), SC7 (structural compatibility).

SC5 and SC7 are in mild tension: SC5 requires recipes that "pass verification" (which may produce different commands than templates), while SC7 requires output "与当前输出一致" (which is template output). If the agent generates a valid but different command for, say, `go vet` (e.g., `go build ./...` instead), SC5 passes but SC7 may fail. The ambiguity requires author clarification.

**D9 Total: 58/80**

---

### D10. Logical Consistency (90 pts)

**Solution addresses the stated problem (30/35)**: The problem is template rigidity; the solution removes templates. Direct alignment. Deduction: the problem states "agent 的三层生成流程...第一步（模板）和第三步（LLM 微调）经常冲突" -- the solution eliminates the first layer entirely, which directly addresses the conflict. However, the problem also mentions "新增语言需从零编写模板" -- the solution doesn't address whether agent knowledge covers all the same languages, or whether adding a new language still requires work (presumably adding a surface rule + convention).

**Scope <-> Solution <-> Success Criteria aligned (22/30)**: In-scope items map to success criteria: template deletion -> SC2, server lifecycle rule -> SC3, surface rule modification -> SC1, SKILL.md rewrite -> no direct SC (gap). The alignment is mostly clean but the SKILL.md rewrite, which is the largest deliverable, has no corresponding success criterion beyond the indirect SC4 (empty surfaces prompt).

**Requirements <-> Solution coherent (20/25)**: The five key scenarios map to solution steps and success criteria: single surface -> SC7, multi-surface -> SC6, mixed language -> SC5, non-standard -> SC5 (implicit), empty surfaces -> SC4. NFRs (consistency, backward compatibility) are acknowledged but the solution's LLM non-determinism undermines both. The orphan is: the proposal doesn't explain how "一致性" as an NFR is achieved with LLM-driven generation beyond "recipe 名称、分组、边界标记不变".

**D10 Total: 72/90**

---

## Phase 3 -- Blindspot Hunt

### [blindspot] B1: No proof-of-concept or pilot evidence

Quote: "Agent（Claude/GPT）已具备主流语言构建命令的知识。Surface rule 文件已定义完整的 recipe 契约。Convention 机制已提供框架特定知识的注入通道。唯一的技术挑战是 server lifecycle bash 代码的可靠性"

The entire proposal rests on the assumption that an LLM can generate correct justfile recipes from surface rules alone, without templates. This is an empirical claim that should be validated before committing to a refactor. No PoC, no A/B test, no single-language pilot is proposed. For a refactor that deletes 811 lines of working code, this is a significant gap.

### [blindspot] B2: Regression testing strategy absent

The proposal lists no testing approach. The current system has 6 templates with 811 lines of tested, working code. The refactor replaces this with LLM-generated content. Quote from Success Criteria: "对 Go/Node/Python/Rust 项目，agent 生成的 recipe 通过 verification step（dry-run + actual execution）". But this is a runtime check, not a regression test. How do you prevent regressions across releases? How do you test that the agent generates correct output for all 5 surface types x 4 languages x 2 platforms = 40 combinations?

### [blindspot] B3: SKILL.md rewrite complexity underestimated

The current SKILL.md is 490+ lines with multiple HARD-RULE blocks, detailed process flows, and extensive verification steps. The proposal says "重写 SKILL.md 流程" as one in-scope item, but doesn't acknowledge the complexity of rewriting a 490-line instruction document that agents must follow precisely. The risk of the rewrite itself introducing behavioral regressions (agent misinterpreting the new instructions) is not addressed.

### [blindspot] B4: Prompt engineering gap

Quote: "agent 根据检测到的语言/框架 + Convention 知识 + 自身知识生成每个 recipe 的具体命令"

This is the core mechanism, yet there's no discussion of prompt design. What context window does the agent need? What's the prompt structure? How are surface rules, Convention data, and project signals combined into a generation prompt? The current HARD-RULE in SKILL.md provides detailed template-loading instructions; the proposal needs equivalent detailed instructions for the new approach, but none are specified.

### [blindspot] B5: Token cost and latency impact

Replacing template loading (file read + string substitution) with LLM-driven generation (full context injection + LLM inference) will increase both token consumption and latency. For a tool that runs interactively (MANUAL-ONLY), this matters. No estimate of the cost/latency impact is provided.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1. Problem Definition | 90 | 110 |
| D2. Solution Clarity | 95 | 120 |
| D3. Industry Benchmarking | 86 | 120 |
| D4. Requirements Completeness | 85 | 110 |
| D5. Solution Creativity | 72 | 100 |
| D6. Feasibility | 79 | 100 |
| D7. Scope Definition | 68 | 80 |
| D8. Risk Assessment | 66 | 90 |
| D9. Success Criteria | 58 | 80 |
| D10. Logical Consistency | 72 | 90 |
| **Total** | **771** | **1000** |

**Pass threshold**: 900/1000
**Result**: FAIL (129 points below threshold)
