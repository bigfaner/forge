# Baseline Evaluation: Skill & Command Independence Audit

**Reviewer**: CTO Adversary (Baseline)
**Date**: 2026-06-03
**Document**: docs/proposals/skill-command-independence-audit/proposal.md

---

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

The problem statement identifies three classes of issues: (1) cross-skill internal file references, (2) redundant descriptions, (3) low-value Related/Integration/Reference sections. The proposed solution addresses each with a corresponding action: inline, compress, delete. This trace is clean on the surface.

However, the coupling graph is incomplete. The Evidence section lists "gen-contracts: 内联 gen-jours Surface Detection 相关知识" in Scope, but the Evidence section does not explicitly list gen-contracts -> gen-journeys/SKILL.md as a cross-skill reference. The reference is confirmed at gen-contracts/SKILL.md line 58: `See \`gen-journeys/SKILL.md\` "Surface Detection" section`. The Scope section implicitly covers it under "gen-contracts: 内联 gen-jours Surface Detection 相关知识", but the bidirectional nature of gen-journeys <-> gen-contracts coupling is not made explicit. The freeform review (R1) flagged this correctly.

### Solution -> Evidence Chain

The solution claims 6 cross-skill references. Verified:
1. gen-journeys -> gen-contracts/rules/journey-contract-model.md (confirmed, 3 references at lines 20, 305, 386)
2. gen-test-scripts -> run-tests/rules/test-isolation.md (confirmed, line 227)
3. extract-design-md -> ui-design/templates/styles/ (confirmed, line 124)
4. init-justfile -> test-guide/references/test-type-model.md (confirmed, line 490)
5. gen-contracts -> gen-journeys/SKILL.md (confirmed, line 58) -- NOT listed in Evidence section
6. fix-bug command -> learn/templates/ and consolidate-specs/rules/ (confirmed, lines 241-260)

The Evidence section lists 5 skill-to-skill references but misses gen-contracts -> gen-journeys/SKILL.md. It appears only in the Scope section as an action item. This is a gap in traceability.

### SC Clustering

Five success criteria. Three are binary presence/absence checks (0 references, 0 Related sections) -- these are strong. The 15% line reduction target is quantified but the freeform review (S7) correctly demonstrates that inline operations will ADD lines, making the net 15% target tight. The "functional equivalence" SC is unmeasurable as stated.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (72/110)

**Problem stated clearly (32/40)**: Three classes of issues clearly enumerated with specific counts. Deduction: The coupling graph is incomplete -- gen-contracts -> gen-journeys/SKILL.md reference is not in Evidence, only implicitly in Scope. Quote: "gen-contracts 引用 gen-journeys/SKILL.md" is absent from the Evidence bullet list.

**Evidence provided (30/40)**: Specific file paths and line references are given. The "total ~6000 行中有约 30% 可精简" claim is precise enough (verified: 6011 skill lines + 1127 command lines = 7138). Deduction: The 30% claim is asserted without per-file breakdown. The Evidence section also says "6 处跨 skill 内部文件引用" but actually lists only 5 skill-to-skill references in its bullet (gen-contracts -> gen-journeys/SKILL.md is missing from the list).

**Urgency justified (10/30)**: "v3.0.0 开发阶段是清理文档债务的窗口期" is a timing argument, but no evidence is provided that waiting would cause concrete harm. The claim "随着 skill 数量增长，维护成本将持续上升" is speculative -- no historical example of a maintenance incident caused by cross-skill coupling is cited.

### 2. Solution Clarity (85/120)

**Approach concrete (35/40)**: Three clear actions: inline, compress, delete. Each skill/command is listed with its specific action. Deduction: The phrase "内联 XXX 所需内容" is used uniformly without defining what "所需内容" means per case. journey-contract-model.md is 184 lines; not all of it is needed by gen-journeys.

**User-facing behavior described (35/45)**: "AI agent 加载单个 skill 时，无需读取其他 skill 的内部文件即可完整理解并执行" is clear. Deduction: No before/after example is provided. A single concrete example of how a skill file changes (e.g., showing the gen-journeys reference before and after inlining) would significantly improve clarity.

**Technical direction clear (15/35)**: "将引用的外部知识内联到引用方" is the entire technical specification. There is no guidance on: (a) how to handle the extract-design-md case where the reference is to runtime data files (7 style templates, 894 lines total), not knowledge; (b) how to ensure inline fidelity; (c) what to do when inlining would bloat a skill significantly.

### 3. Industry Benchmarking (55/120)

**Industry solutions referenced (20/40)**: "标准做法是'模块自包含'" is a one-sentence mention. No specific methodology, standard, or toolchain is cited (e.g., DRY vs WISP in documentation, docs-as-code patterns, Sphinx/Jekyll cross-reference handling, Kubernetes docs module independence).

**At least 3 alternatives (15/30)**: Three alternatives are listed in the comparison table (do nothing, shared layer, inline+compress). This meets the count requirement but the analysis is shallow -- each gets one sentence of pros/cons.

**Honest trade-offs (10/25)**: The table admits "知识多份存在，可能漂移" as a con. However, the Risk section marks drift as M likelihood / L impact with "可接受" -- the freeform review (R8) correctly identifies this as under-analyzed with no quantitative backing.

**Chosen approach justified (10/25)**: "符合 Forge 分发模型" is the sole justification. This is reasonable but not rigorously argued -- no analysis of WHY Forge's distribution model makes shared layers infeasible, only an assertion.

### 4. Requirements Completeness (68/110)

**Scenario coverage (30/40)**: Three scenarios listed, all realistic and covering the main use cases. Deduction: Missing scenario -- what happens when a shared concept (like Journey/Step definitions) needs to evolve? How does the inlined-copy model handle concept evolution? This is a scenario the solution must address.

**NFRs (18/40)**: Two NFRs: functional equivalence and line reduction. Both are important but incomplete. Missing NFRs: (a) no regression in AI agent task completion rate; (b) no increase in individual skill file size beyond a threshold; (c) maintainability of inlined copies over time; (d) readability for new developers.

**Constraints (20/30)**: Two constraints listed (forensic exemption, docs-only). The forensic exemption is correct and well-justified. Deduction: Missing constraint -- extract-design-md's reference to ui-design/templates/styles/ is a data dependency, not a knowledge dependency. This should be called out as a constraint on the inlining approach.

### 5. Solution Creativity (30/100)

**Novelty (10/40)**: The proposal explicitly states "无创新，标准文档清理". Honest. No novelty penalty per se, but the score reflects the absence of creative problem-solving.

**Cross-domain inspiration (10/35)**: No cross-domain references. The "每个 skill 文件是一个独立的知识单元" principle is standard modular documentation, not inspired by any external domain.

**Simplicity of insight (10/25)**: The insight that Related Skills sections are redundant because "pipeline 上下游关系已在正文流程中体现" is a useful observation. However, the freeform review (R6) demonstrated that gen-contracts and gen-journeys Reference sections contain concept definitions not found elsewhere in the skill body -- so this insight is partially wrong.

### 6. Feasibility (70/100)

**Technical feasibility (35/40)**: "纯文档编辑，无技术风险" is accurate. Deduction: The extract-design-md case involves inlining style matching logic from 7 files totaling 894 lines, which is more complex than "纯文档编辑" and the proposal does not address this.

**Resource/timeline (25/30)**: "预计 1 个 session 可完成" is reasonable for the scope described. Deduction: The freeform review (S7) shows inline operations will ADD 150-250 lines while needing to cut 1200-1400 lines elsewhere for the 15% target. This compression workload in one session is tight.

**Dependency readiness (10/30)**: "无外部依赖" is correct for the editing task itself. However, the proposal depends on the assumption that all cross-references have been identified. The missing gen-contracts -> gen-journeys reference suggests the dependency analysis is incomplete.

### 7. Scope Definition (60/80)

**In-scope concrete (25/30)**: Detailed lists per skill/command with specific actions. Deduction: The gen-contracts entry says "内联 gen-jours Surface Detection 相关知识" -- "gen-jours" is a typo for "gen-journeys", and "Surface Detection 相关知识" is vague about what exactly gets inlined.

**Out-of-scope listed (20/25)**: Four clear out-of-scope items. Deduction: Does not explicitly exclude the "Reference" sections that contain concept definitions (gen-contracts Reference, gen-journeys Reference). These are lumped with Related Skills but serve a different purpose.

**Scope bounded (15/25)**: The scope is bounded by file list. However, the boundary between "redundant description" and "behaviorally important context" is not defined, leaving the implementer to make judgment calls on every line.

### 8. Risk Assessment (48/90)

**Risks identified (18/30)**: Three risks listed. Missing risks: (a) extract-design-md style template drift after inlining; (b) incorrect identification of "redundant" content that AI agents actually depend on; (c) the freeform review's R5 finding that quick-tasks' `## Reference Files` is a template usage guide, not pipeline info -- deleting it would be a bug.

**Likelihood+impact (15/30)**: L/M/M ratings are provided but not justified. The "精简过度导致 AI agent 行为偏差" risk is rated L likelihood -- given that the proposal intends to cut ~1200 lines, this seems optimistic.

**Mitigations actionable (15/30)**: "内联后对比原文确保无遗漏" is a principle, not a mitigation. "可接受" is not a mitigation. The only semi-actionable mitigation is "保留所有硬规则和决策表，只压缩描述性文字" -- but distinguishing "描述性" from "行为指导性" text is the hard part.

### 9. Success Criteria (48/80)

**Measurable/testable (20/30)**: Three binary checks (0 references, 0 Related sections) are measurable. The 15% line reduction is measurable. "功能等价" is not measurable as stated -- the freeform review (R9) correctly identifies this.

**Coverage complete (15/25)**: Covers the three problem classes. Missing: no SC for "no regression in AI agent task completion quality" and no SC for "inline fidelity" (ensuring inlined content accurately represents the original).

**SC consistency (13/25)**: The 15% line reduction target may conflict with the "功能等价" SC. The freeform review (S7) shows that achieving 15% net reduction requires aggressive compression that may conflict with preserving behavioral guidance. This tension is not acknowledged.

### 10. Logical Consistency (60/90)

**Solution addresses problem (28/35)**: The three-pronged solution maps cleanly to the three problem classes. The extract-design-md exception is the one gap -- the proposal treats it identically to knowledge references, but it is a runtime data dependency.

**Scope <-> Solution <-> SC aligned (15/30)**: The Scope lists specific files but the Solution says "对全部 21 个 skill 和 16 个 command". Not all 21 skills and 16 commands are in scope -- only ~15 skills and ~3 commands have actions. The wording creates a misleading impression of blanket coverage.

**Requirements <-> Solution coherent (17/25)**: The requirement "每个 skill 文件是独立知识单元" is coherent with the inlining solution. However, the solution's treatment of Reference sections (delete) conflicts with the independence requirement when those sections contain unique concept definitions (gen-contracts Reference at lines 267-276, gen-journeys Reference at lines 384-392). Deleting these would make the skills LESS self-contained because the definitions would exist nowhere in the skill file.

---

## Phase 3: Blindspot Hunt

### B1. No verification methodology for "functional equivalence"

The proposal's most critical SC -- "所有 skill/command 修改后功能等价" -- has no verification plan. The freeform review (S6) suggests a concrete checklist (grep for HARD-RULE/HARD-GATE/etc. tags), but the proposal itself is silent. For a document that guides AI agent behavior, "functional equivalence" can only be verified by running the agent on representative tasks and comparing outputs. The proposal should acknowledge this gap and propose at minimum a structural checklist.

### B2. extract-design-md is a qualitatively different case

Quote: "extract-design-md 引用 ui-design/templates/styles/" and "内联 ui-design/styles 匹配逻辑". The reference at line 124 of extract-design-md/SKILL.md is: `read the corresponding style file from \`ui-design/templates/styles/<name>.md\``. This is an instruction to read runtime data files (7 files, 894 lines), not to understand a concept. Treating this the same as "inline the journey-contract-model concepts into gen-journeys" is a category error. The proposal does not distinguish between "knowledge references" (concepts needed to understand the skill) and "data references" (files the skill processes at runtime).

### B3. The 15% target may incentivize over-compression

Quote: "总行数减少 >= 15%". With 7138 total lines, this means cutting ~1071 lines. The freeform review calculates that inline operations add ~150-250 lines, requiring ~1200-1400 lines of compression. Given that the proposal specifies "保留所有硬规则和决策表，只压缩描述性文字", the question is whether there are 1200-1400 lines of pure description that can be cut without affecting agent behavior. The proposal provides no per-file compression estimate to support this target.

---

## Dimension Scores Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 72 | 110 |
| Solution Clarity | 85 | 120 |
| Industry Benchmarking | 55 | 120 |
| Requirements Completeness | 68 | 110 |
| Solution Creativity | 30 | 100 |
| Feasibility | 70 | 100 |
| Scope Definition | 60 | 80 |
| Risk Assessment | 48 | 90 |
| Success Criteria | 48 | 80 |
| Logical Consistency | 60 | 90 |
| **Total** | **596** | **1000** |
