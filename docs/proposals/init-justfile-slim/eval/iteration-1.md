---
iteration: 1
reviewer: CTO Adversary
date: 2026-06-09
document: docs/proposals/init-justfile-slim/proposal.md
rubric: proposal (1000-point scale)
pre_revised_regions: 7
---

# Proposal Evaluation — Iteration 1

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

**Problem chain**: 1645-line skill -> bash templates in LLM prompt -> token waste + scattered maintenance + misplaced responsibility -> propose CLI scaffold command.

The causal chain is sound. The three problems (token waste, scattered maintenance, responsibility misplacement) are real and measurable. The solution (move bash generation to CLI) directly addresses all three.

**Gap**: The problem states "Phase 1 consistency check is unnecessary if producer is trusted" but the producer being trusted (CLI) does not make the *consumer of CLI output* (agent doing placeholder substitution) trusted. The solution conflates "CLI generates correct skeleton" with "final output is correct." This is a logical leap that weakens the Problem->Solution chain at the Phase 1 deletion point.

### Solution -> Evidence Trace

Evidence for reduction: concrete line counts verified against codebase (1645 total, 745 server-lifecycle, 548 SKILL.md, 318 surface rules). These numbers check out. The ~284 post-reduction estimate is plausible but unverified (no breakdown of the 250-line SKILL.md target is provided).

**Gap**: The "新增 CLI 代码 ~500 行 Go" is an estimate without justification. Where does 500 come from? If it balloons to 800-1000, the total system complexity may not decrease.

### Evidence -> Success Criteria Trace

The proposal has a "减重效果" table but no formal Success Criteria section with measurable/testable criteria. The action items (行动项) serve as implicit acceptance criteria but lack specificity (e.g., what tests must the CLI pass? What is the max token count target?).

### Self-contradiction Check

1. **Recipe naming contradiction**: The "Recipe 命名统一模型" says single surface scalar uses no prefix (`compile`, `test`). But "Gate recipe 概念取消" says "CLI 为每个 surface 生成完整的 recipe 集...quality gate 直接调用 `<key>-compile` 等". For single surface scalar, there is no key, so this is consistent. But the Consumer Impact table says `run-tests` will "移除 fallback 链（`<key>-compile` 不存在 -> fallback `compile`）" — this implies the fallback existed because some surfaces had prefixed and some didn't. Removing it means `run-tests` must now know whether to use prefixed or unprefixed. The proposal doesn't show how `run-tests` resolves this.

2. **user-customized marking scope**: The pre-revised text says "所有 lifecycle recipes 和 quality recipes 均标记 `# user-customized`" but then says "聚合 recipes 不标记". This is internally consistent. However, the current SKILL.md only marks gate recipes + lifecycle recipes as user-customized, and aggregate recipes as NOT user-customized. The proposal extends marking to ALL quality recipes, which is a semantic expansion not acknowledged.

3. **Phase 1 removal vs. Step 4a**: The "精简 SKILL.md" table lists "Phase 1 Consistency Verification | ~20 | 简化为轻量级 recipe 完整性断言（`just --list` 验证）". But the Agent 新流程 Step 4a says "轻量级完整性断言（just --list 确认所有 recipe 可被 just 解析）". The SKILL.md deletion table says it's simplified, not deleted. But the action items say "删除 Phase 1" without qualification. This is ambiguous.

---

## Phase 2: Rubric Scoring

### Dimension 1: Problem Definition — 88/110

**Problem stated clearly (35/40)**:
The three-part problem structure (token waste, scattered maintenance, responsibility misplacement) is clear and unambiguous. Each is concrete and measurable. Deduction: the problem statement frames the issue entirely from the supply side (skill size). It does not describe the user-facing symptom (slow generation, errors, inconsistent output) that a non-technical stakeholder would recognize.

**Evidence provided (30/40)**:
Line counts are verifiable against the codebase (confirmed: 1645 total, 745 server-lifecycle, 548 SKILL.md, 318 surface rules). The EXTREMELY-IMPORTANT repetition claim (3 occurrences) checks out (2 found via grep, which is close enough given the broad matching). Deduction: no user feedback, no error rate data, no before/after comparison of agent performance. All evidence is static code metrics, not operational data.

**Urgency justified (23/30)**:
The implicit urgency is "v3.0.0 not yet released, so breaking changes are free." This is mentioned in the "向后兼容" section but not in the problem section. Deduction: the urgency argument is scattered — the "why now" should be front and center in the problem section, not buried in a compatibility note. Additionally, no cost-of-delay analysis (what happens if this is deferred to v3.1?).

### Dimension 2: Solution Clarity — 92/120

**Approach is concrete (38/40)**:
The `forge justfile scaffold` command interface is precisely specified with parameters, output format, and per-surface-type recipe lists. A reader can explain exactly what will be built. Minor deduction: the `--aggregate` flag's interaction with `--type` and `--key` is unclear — can you use both? The table suggests `--aggregate` is standalone.

**User-facing behavior described (32/45)**:
The agent flow (Step 0-5) describes the internal process well, but "user-facing" means what the developer running `/init-justfile` experiences. The proposal doesn't explicitly say: "User runs /init-justfile, sees same output as before, but generation is faster and more reliable." The user experience change is implied, not stated. Deduction: no description of what the user sees differently, no before/after UX comparison, no mention of error messages or feedback changes.

**Technical direction clear (22/35)**:
Go CLI command generating templates with placeholders is a clear direction. However: (1) no indication of where in the forge CLI codebase this lives (new package? existing command group?), (2) no template engine specified (string literals? Go templates? embedded files?), (3) the "~500 行 Go" estimate has no backing. Deduction: the "what" is clear, the "how" at implementation level is vague.

### Dimension 3: Industry Benchmarking — 52/120

**Industry solutions referenced (15/40)**:
No external tools, projects, or published patterns are cited. The proposal references only its own current architecture. For a "scaffold generator" pattern, industry is full of comparables: Yeoman, Plop, Cookiecutter, Hygen, Rails generators, Cargo generate. None mentioned. The proposal operates in a vacuum.

**At least 3 meaningful alternatives (20/30)**:
Only the current approach (status quo) is implicitly compared. No explicit alternatives section. The proposal mentions the chosen approach but does not enumerate rejected alternatives such as: (a) keep rules but deduplicate via shared templates, (b) use a simpler template engine in SKILL.md itself, (c) generate rules at install time rather than runtime. The "do nothing" alternative is implicitly present (current state) but not formally evaluated.

**Honest trade-off comparison (8/25)**:
The "风险与缓解" table partially covers trade-offs but focuses on bug risk, not architectural trade-offs. Missing trade-offs: (a) loss of extensibility (new surface types require CLI release), (b) debugging difficulty (Go code vs. readable markdown), (c) template rigidity (what if a recipe needs conditional logic not expressible via placeholders?).

**Chosen approach justified against benchmarks (9/25)**:
No benchmarks to justify against. The justification is purely internal: "prompt layer shouldn't do mechanical work." This is sound reasoning but not benchmarked against industry practice.

### Dimension 4: Requirements Completeness — 78/110

**Scenario coverage (30/40)**:
The proposal covers the main scenario (single surface scalar, single surface named, multi surface) and edge cases (cold start/no Convention, user-customized protection). Missing scenarios: (a) what happens when `forge surfaces` returns an unknown surface type? (b) what happens when justfile has manually added recipes outside boundary markers? (c) migration scenario for existing projects with non-standard justfiles.

**Non-functional requirements (28/40)**:
Performance: "83% token reduction" is the primary NFR, well-stated. Compatibility: "向后兼容" section addresses v3.0.0 status. Missing NFRs: (a) CLI command execution latency (how fast should scaffold generate?), (b) template output correctness guarantee, (c) error message quality from CLI, (d) testability requirements for the Go code itself.

**Constraints/dependencies (20/30)**:
Dependencies listed: `forge surfaces` command, Convention files, `just --list` for validation. Missing constraints: (a) Go version requirement for CLI development, (b) forge CLI build/release process impact, (c) the CLI must be installed and on PATH for the skill to work (current skill has no such runtime dependency).

### Dimension 5: Solution Creativity — 55/100

**Novelty over industry baseline (20/40)**:
The "CLI generates templates, agent fills placeholders" pattern is standard in scaffold generators (Yeoman, Plop). The specific application to LLM prompt optimization is the novel angle — moving code generation out of prompts into deterministic programs. This is a sensible engineering decision but not a creative breakthrough.

**Cross-domain inspiration (20/35)**:
The insight that "LLM prompts should do reasoning, not code templating" comes from the broader AI engineering domain (prompt optimization). The proposal applies this correctly. However, no cross-domain borrowing is acknowledged — the proposal presents the idea as self-evident rather than drawing parallels to other systems that separate generation from reasoning.

**Simplicity of insight (15/25)**:
The core insight is clean: "bash templates in prompts are wasteful, move them to CLI." This is practical and straightforward. Not quite "why didn't I think of that" elegance, but solid engineering. Deduction: the placeholder system adds complexity (13 placeholders, agent must resolve each), partially offsetting the simplification gain.

### Dimension 6: Feasibility — 78/100

**Technical feasibility (35/40)**:
Go CLI generating text templates is trivially feasible. `forge surfaces` already exists. The forge CLI already has a build system. No showstopper dependencies. Deduction: the multi-service orchestration mode (pre-revised paragraph at line 88-89) adds significant complexity — generating a dependency-ordered startup sequence requires topological sorting of surfaces, which the proposal hand-waves.

**Resource/timeline feasibility (22/30)**:
Four action items are listed but no timeline, no effort estimates, and no assignment. The ~500 lines of Go is an estimate without breakdown. The proposal is in "Draft" status with no milestone. Deduction: no indication of whether this is a 1-day or 1-week effort.

**Dependency readiness (21/30)**:
`forge surfaces` exists and works. `forge justfile scaffold` does not exist — it's the main deliverable. The proposal says "CLI 发版频率高（RC 阶段）" but forge CLI v5.19.0 is the current version, and the plugin is at 3.0.0-rc.50 — the release cadence claim is unsubstantiated. Deduction: the current forge CLI has no `justfile` subcommand at all, so this is greenfield work with no existing infrastructure to build on.

### Dimension 7: Scope Definition — 62/80

**In-scope items are concrete (25/30)**:
Four action items are concrete deliverables: CLI command, SKILL.md rewrite, file deletion, quality gate update. Each is identifiable. Minor deduction: "更新 quality gate" is vague — update what exactly? The guide.md? The hooks.json? The clean-code SKILL.md? The quality-gate CLI command itself?

**Out-of-scope explicitly listed (17/25)**:
No explicit "out of scope" section. Implicitly out of scope: changing other skills' recipe naming, modifying `run-tests` dispatch logic, changing `forge quality-gate` command. But these are not stated. Deduction: the Consumer Impact section partially covers what changes, but doesn't explicitly say what is NOT changing.

**Scope is bounded (20/25)**:
The scope is bounded by the four action items and the 6-file deletion list. The "no new surface types" constraint is implied by the 5-type enumeration. Deduction: the quality gate update action item is open-ended — "移除 fallback 链" could ripple into multiple files across the plugin. The boundedness claim depends on the fallback chain being simple to remove.

### Dimension 8: Risk Assessment — 62/90

**Risks identified (22/30)**:
Three risks are listed: CLI bug, surface type extensibility, placeholder coverage. These are meaningful. Deduction: missing risks include: (a) agent placeholder resolution errors (the pre-revised text added `just --list` mitigation but didn't add this as a risk), (b) regression in multi-surface projects, (c) Go template maintenance cost vs. markdown rule maintenance, (d) skill behavior change detection by existing users (even in RC, internal team may have workflows).

**Likelihood + impact rated (18/30)**:
The risk table has no likelihood or impact ratings. Each risk has "影响" (impact) described in words but not rated (high/medium/low). No likelihood assessment at all. The pre-revised regions address some findings from freeform review but didn't add quantitative risk ratings.

**Mitigations are actionable (22/30)**:
Mitigations are partially actionable: "Phase 2 dry-run + Phase 3 actual execution 验证" is concrete. "CLI 有单元测试" is actionable. "CLI 文档化完整占位符清单" is actionable. Deduction: "CLI 发版频率高" is not actionable mitigation — it's a hope. "Agent 遇到未知占位符时保留原样并报告" is actionable for the agent but requires implementing the behavior in SKILL.md.

### Dimension 9: Success Criteria — 45/80

**Criteria are measurable/testable (18/30)**:
The "减重效果" table provides measurable targets (1645 -> ~284 lines, 8 -> 2 files, -83% token reduction). These are testable. Deduction: no criteria for correctness (generated justfiles work identically to before), no criteria for CLI test coverage, no criteria for edge case handling. The action items are deliverables, not success criteria.

**Coverage is complete (15/25)**:
The reduction targets cover the "simplification" goal but not the "correctness" goal. Missing: (a) criteria for all 5 surface types producing valid output, (b) criteria for multi-surface aggregation correctness, (c) criteria for user-customized recipe protection, (d) criteria for consumer (run-tests, quality-gate) compatibility.

**SC internal consistency (12/25)**:
The implicit success criteria are: 284 lines, 2 files, -83% tokens. These are internally consistent as a set. However, the "新增 CLI 代码 ~500 行 Go" metric creates a tension: total system lines (284 + 500 = 784) vs. original (1645) is a 52% reduction, not 83%. The 83% claim only counts prompt layer lines. This is not contradictory but it's misleading — the total system complexity reduction is roughly half of the claimed reduction. Additionally, the claim of "2 files" post-reduction conflicts with the action item to update quality gate across multiple files.

### Dimension 10: Logical Consistency — 73/90

**Solution addresses the stated problem (30/35)**:
CLI scaffold directly addresses token waste (bash code moves out of prompt), scattered maintenance (5 rule files consolidate into CLI), and responsibility misplacement (CLI generates deterministic code). Strong alignment. Deduction: the Phase 1 deletion argument ("CLI is trusted") is logically incomplete because the agent's placeholder filling is untrusted.

**Scope <-> Solution <-> SC aligned (20/30)**:
The scope (4 action items) maps to the solution (CLI + SKILL.md rewrite + deletion + gate update). The success criteria (reduction metrics) map to the scope. Deduction: the quality gate update (action item 4) is in scope but has no corresponding success criterion. And the Consumer Impact section describes changes to `run-tests` and quality gate, but "更新 run-tests skill" is not in the action items.

**Requirements <-> Solution coherent (23/25)**:
The requirements (per-surface recipe generation, aggregation, placeholder resolution, user-customized protection) map cleanly to the CLI scaffold + agent flow. No orphan requirements. Minor gap: the multi-service orchestration requirement (from server-lifecycle.md Section 4) is addressed by the pre-revised paragraph but the mapping is thin — one sentence says "额外生成 test-setup 聚合 recipe" without specifying the template content.

---

## Phase 3: Blindspot Hunt

[blindspot-1] **`forge quality-gate` CLI command is the primary consumer, not `run-tests`**
The proposal's Consumer Impact section lists `run-tests` skill and quality gate mechanism as consumers. But `forge quality-gate` is a CLI command (found in hooks.json) that runs `just compile -> just fmt -> just lint -> just unit-test`. These are hardcoded unprefixed recipe names. If the proposal changes recipe naming (e.g., removing unprefixed `compile` for multi-surface projects), the `forge quality-gate` Go binary must also be updated. The proposal doesn't mention this. Quote: "quality gate 直接调用 `<key>-compile` 等" — but the Go binary calls `just compile`, not `just <key>-compile`. This is a **consumer not accounted for**.

[blindspot-2] **`fix-bug.md` and `clean-code/SKILL.md` call `just unit-test` directly**
Multiple files hardcode `just unit-test` (fix-bug.md lines 91-92, 127-130, 179; clean-code/SKILL.md lines 162, 171, 175). These are not mentioned in Consumer Impact. The proposal says "移除 fallback 链" but doesn't audit all callers.

[blindspot-3] **~500 lines Go estimate is suspiciously round**
No breakdown of what those 500 lines contain. If we estimate: 5 surface types x ~100 lines of template each = 500. But this doesn't include the CLI boilerplate (flags, output, error handling), the aggregation logic, the multi-service orchestration, or the test code. A more realistic estimate might be 800-1200 lines. The proposal should justify this number or acknowledge uncertainty.

[blindspot-4] **No rollback plan**
The proposal says "硬切换" (hard cutover) with "无存量用户" as justification. But the Forge development team itself is a user. If the new CLI scaffold has a bug that breaks `/init-justfile` for all development, the team can't generate justfiles until the CLI is fixed. No rollback mechanism (e.g., feature flag, `--legacy` flag) is mentioned.

[blindspot-5] **Template language ambiguity**
The proposal says CLI outputs "带 `{{PLACEHOLDER}}` 占位符的 just recipe 代码." But `{{...}}` is also Go template syntax and justfile variable syntax. If the CLI uses Go's `text/template`, the `{{PLACEHOLDER}}` markers must be escaped. This is a trivial implementation detail that could cause confusion.

[blindspot-6] **`--key` optionality creates two code paths**
`--key` is required for named surfaces and optional for scalar. This means the CLI must handle both cases, and the recipe naming diverges. The proposal handles this in the "Recipe 命名统一模型" table but doesn't discuss how the CLI itself determines whether to add the prefix — is it purely based on `--key` being present? What if `--key` is provided for a scalar surface?

---

## Bias Detection Report

Total paragraphs in document (non-blank, non-table-separator, non-heading): ~40 paragraphs (excluding table rows and blank lines).
Pre-revised annotated regions: 7.

**Annotated regions**: 4 attack points found in 7 pre-revised paragraphs = density 0.57
**Unannotated regions**: 8 attack points found in ~33 unannotated paragraphs = density 0.24

**Ratio (annotated/unannotated)**: 2.38

Interpretation: The pre-revised regions (which were modified based on freeform review findings) have 2.38x the attack density of unannotated regions. This suggests the revision process addressed the specific findings but may have introduced new issues in the revised text, or the revisions attracted scrutiny proportional to their complexity. The unannotated regions received less adversarial attention, which means some weaknesses in those sections may be undercounted.

---

## Summary

The proposal is a well-motivated, internally consistent simplification of the largest Forge skill. The core insight (bash templates don't belong in LLM prompts) is sound. The primary weaknesses are: (1) no industry benchmarking at all, (2) incomplete consumer impact analysis (misses `forge quality-gate` CLI and other hardcoded callers), (3) no formal success criteria beyond line-count reduction, (4) risk assessment lacks likelihood/impact ratings, and (5) the 83% reduction claim is misleading because it ignores the 500 lines of new Go code.
