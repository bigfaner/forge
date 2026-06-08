# Baseline Score Report — Behavioral Test Accuracy

**Iteration**: 0 (baseline)
**Date**: 2026-06-08
**Scorer**: CTO Adversary

---

## Phase 1: Reasoning Audit

### 1. Problem → Solution
The problem is that generated tests are structural (CRUD on empty containers, HTTP 200 assertions) rather than behavioral. The proposed solution targets three pipeline stages: (a) gen-journeys forces Golden Path journeys, (b) gen-contracts adds Fixture Specification, (c) gen-test-scripts enforces deep assertions and rich fixtures. The chain holds — each layer addresses one root cause level (L1, L2, L3). This is well-structured.

### 2. Solution → Evidence
Evidence is a single project case (pm-work-tracker). The solution extrapolates from one data point to a systemic claim. No quantitative evidence of how widespread this problem is across other Forge-generated projects is provided.

### 3. Evidence → Success Criteria
Success criteria measure process outputs (journeys contain golden paths, contracts have fixture specs, 80% behavioral assertions). They do NOT measure the actual goal: "tests can discover real bugs." The chain from evidence to success criteria has a gap — passing all SCs does not guarantee the problem is solved.

### 4. Self-contradiction Check
The proposal claims to solve "structural tests passing on empty containers" but SC-3 allows 20% structural assertions. This is explicitly acknowledged as a mitigation for legitimate cases (logging, health check). Acceptable trade-off, not a true contradiction.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 36/40 | Problem is unambiguous: tests are structural not behavioral. One interpretation only. Minor deduction: "structural" vs "behavioral" binary could be more nuanced — there's a spectrum. |
| Evidence provided | 34/40 | pm-work-tracker case is concrete and real. L1/L2/L3 root cause analysis is strong. Deduction: single project evidence, no data on prevalence across other projects. "这是管线架构的系统性缺陷" is asserted but only one example supports it. |
| Urgency justified | 28/30 | "pm-work-tracker 项目已因此产生了虚假的'全部通过'结果" — concrete cost of delay. Good. Minor deduction: no mention of how many other projects are affected or upcoming deadlines. |

**Dimension Total: 98/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 36/40 | Three clear pipeline stages with specific modifications. Reader can explain back exactly what changes. Deduction: "断言深度规则" and "seed data 丰富度规则" are named but their exact content is unspecified — what constitutes a "business result" assertion vs "structural" assertion? |
| User-facing behavior described | 40/45 | The happy path scenario is well-described. Deduction: what does the user (developer running the pipeline) SEE differently? Error messages? New output files? The observable UX change is implied but not explicitly described. |
| Technical direction clear | 32/35 | Clear enough: modify skill rules, templates, rubrics. Deduction: "从 Contract 读取 fixture spec" — how does gen-test-scripts consume this? By parsing markdown? By structured YAML frontmatter? The integration mechanism is vague. |

**Dimension Total: 108/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | Playwright/Cypress, TestNG/JUnit, Pact are cited. Deduction: these are testing frameworks, not AI-generated test pipeline solutions. The actual industry benchmark — AI code generation tools (Cursor, Copilot test generation, Codeium) and their test quality approaches — is completely absent. The benchmarks cited are for human-written tests, which is the baseline Forge is already trying to automate. |
| At least 3 meaningful alternatives | 25/30 | 4 alternatives listed including "do nothing." Deduction: "仅增强 eval rubric" is attributed to "eval-test-cases 提案" but no detail given — is this a real existing proposal or a straw reference? "轻量级规则增强" has no source attribution and reads like a straw man: "上下游信息链断裂" dismisses it without evidence. |
| Honest trade-off comparison | 20/25 | The selected approach lists "改动涉及 3 skill + 2 rubric" as its con, which is honest. Deduction: pros column says "从源头保证行为性描述" — this is the goal, not a demonstrated advantage. The con of "改动涉及 3 skill + 2 rubric" understates the coordination risk. |
| Chosen approach justified against benchmarks | 20/25 | Justification is that it's the only approach addressing all three root cause levels. This is internally consistent. Deduction: no evidence that the "全链路" approach is superior to a more targeted fix with iteration. The binary framing (do nothing / only eval / only scripts / full chain) excludes incremental approaches. |

**Dimension Total: 95/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | Happy path, edge cases (simple feature, single-entity CRUD), and error scenarios (missing PRD workflow, unverifiable entity relations) are all covered. Deduction: missing scenario for when fixture specification is present but the test execution environment doesn't support creating the declared fixtures (e.g., missing API endpoints, insufficient permissions). |
| Non-functional requirements | 24/40 | No NFRs mentioned. Performance impact of richer fixtures (longer test execution time), security implications of more realistic test data, backward compatibility of existing projects — all absent. The "Constraints & Dependencies" section lists dependencies but these are functional dependencies, not NFRs. |
| Constraints & dependencies | 28/30 | Clear: 3 existing skills, 2 eval rubrics, no external dependencies, no new pipeline stages. Good. Minor deduction: does not mention dependency on the quality of PRD/Design input — if the PRD itself lacks workflow descriptions, the entire chain fails (acknowledged in error scenarios but not as a hard dependency). |

**Dimension Total: 87/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 35/40 | "Contract-level declarative Fixture Specification" is a genuine innovation — lifting fixture declarations from code to spec level for AI-generated pipelines. The differentiation from human-written test patterns is clearly articulated in the Innovation Highlights section. Deduction: the concept itself (declarative fixtures) exists in many frameworks — the novelty is specifically in applying it to an AI generation pipeline, which narrows the creative contribution. |
| Cross-domain inspiration | 25/35 | Contract Testing (Pact) concepts applied to functional completeness rather than API compatibility. Deduction: only one cross-domain reference (Pact). Could have drawn from type systems (dependent types for data relationships), database migration systems (declarative schema + seed data), or property-based testing (Hypothesis-style data generation strategies). |
| Simplicity of insight | 22/25 | The core insight — "declare fixture needs at the spec level, not guess at generation time" — is elegant and obvious-in-retrospect. Deduction: the three-layer solution (journeys + contracts + test scripts) is somewhat overengineered for the insight; the same principle could potentially be achieved with fewer moving parts. |

**Dimension Total: 82/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | All changes are in existing skill files (rules, templates, rubrics). No architecture changes. Feasible. Deduction: "从 Contract 读取 fixture spec" requires gen-test-scripts to parse a new Contract dimension — the parsing mechanism is unstated and could be non-trivial depending on how Contracts are structured. |
| Resource & timeline feasibility | 20/30 | States "3 skill modifications + 2 rubric updates, each independent, can be parallelized." No timeline estimate, no resource allocation, no effort sizing beyond "修改量" being described as manageable. The claim that modifications are "独立的" is questionable — gen-test-scripts depends on gen-contracts output format, which depends on gen-journeys output format. They are sequentially dependent. |
| Dependency readiness | 28/30 | No external dependencies, all files are internal. Good. Minor deduction: assumes current skill implementations are modifiable without regression — no mention of existing test coverage for these skills. |

**Dimension Total: 84/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 27/30 | Four concrete deliverables: gen-journeys rules, gen-contracts templates/rules, gen-test-scripts rules, eval rubrics. Each maps to a specific file type. Deduction: "断言深度规则 + seed data 丰富度规则" — these are rule names, not deliverables with acceptance criteria. |
| Out-of-scope explicitly listed | 23/25 | Five items explicitly out of scope, including rationale (eval-diagnostic-mode covers pipeline reliability, run-tests needs no change). Good. Minor deduction: no mention of documentation updates or migration guide for existing projects. |
| Scope is bounded | 22/25 | Bounded by "不改变管线阶段顺序或新增阶段." Good constraint. Deduction: the "两端 eval rubric 同步新增评估维度" in the solution section implies rubric changes, but the rubric scoring weight adjustments (SC-5: 150 pts, SC-6: 100 pts) could cascade to existing eval logic — this scope impact is not acknowledged. |

**Dimension Total: 72/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 22/30 | Three risks identified. Deduction: missing risks: (1) backward compatibility — existing projects with generated tests may break when regenerated with new rules; (2) fixture specification complexity explosion — complex domain models could produce unwieldy fixture specs; (3) false positives — overly strict behavioral assertions could generate flaky tests. The identified risks focus on "too much" but ignore "wrong direction" risks. |
| Likelihood + impact rated | 25/30 | Ratings are provided and varied (M/L, M/M, L/M). Honest assessment. Deduction: "Golden Path 强制要求导致简单 feature 的 Journey 过度膨胀" rated M likelihood — no basis for this rating provided. Could be higher if most Forge features are simple CRUD. |
| Mitigations are actionable | 24/30 | Mitigations are specific: "Golden Path 可以是简洁的多步 CRUD 循环", "Fixture Specification 作为 Preconditions 的子维度", "80% 阈值允许 20% 结构性断言." Good. Deduction: mitigations describe what to do but not how to verify they work — no rollback plan if a mitigation fails. |

**Dimension Total: 71/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 24/30 | SC-1 through SC-4 are measurable (counts, percentages). SC-5 and SC-6 are less measurable — "评分标准包含 Golden Path 存在性和多步覆盖度" describes the rubric content, not the pass/fail threshold. Deduction: SC-3 "≥80% 的断言验证业务结果" — who classifies what counts as a "business result" assertion? This requires a classifier or human judgment, making it less objectively testable than it appears. |
| Coverage is complete | 20/25 | All four in-scope items have corresponding SCs. Deduction: no SC for the eval rubric quality itself — SC-5 and SC-6 say the rubric "新增维度" but don't specify that the new dimensions must achieve a minimum score on themselves (meta-quality). Also no SC for the "两端 eval rubric 同步新增评估维度" promise — only journey and contract rubrics are mentioned, but what about the eval consistency between them? |
| SC internal consistency | 20/25 | SCs are mostly consistent as a set. Deduction: SC-1 requires golden paths for features "包含父子实体关系" — but SC-3 applies "每个" test's assertions at ≥80%. This means ALL features must have ≥80% behavioral assertions, but only parent-child features need golden paths. For simple features without parent-child relationships, the 80% threshold may be inappropriately high if most legitimate assertions are structural (e.g., CRUD on a single entity). The relationship between SC-1's scope (parent-child only) and SC-3's scope (all features) creates an implicit tension. |

**Dimension Total: 64/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 32/35 | Solution addresses all three root cause levels. Strong alignment. Deduction: the stated problem is "测试通过但核心功能可能完全不可用." The success criteria measure process outputs (golden paths exist, fixture specs exist) but do NOT measure whether tests actually catch real bugs. A project could satisfy all SCs and still have tests that pass on broken functionality — the criteria proxy the solution, not the problem. |
| Scope ↔ Solution ↔ Success Criteria aligned | 25/30 | Generally aligned. Four scope items map to solution components and SCs. Deduction: the solution mentions "两端 eval rubric 同步新增评估维度" but scope item 4 only says "eval rubrics (journey.md, contract.md): 新增对应评估维度" — the "同步" (synchronized) aspect is in the solution but not in scope or SCs. |
| Requirements ↔ Solution coherent | 22/25 | Requirements map to solution. Happy path, edge cases, and error scenarios all have solution components. Deduction: the error scenario "PRD 中缺少工作流描述：gen-journeys 应从用户故事中推断，或要求补充" maps to a solution behavior ("应从...推断") but has no corresponding SC or scope item. This is an orphan requirement — stated but not tracked for delivery. |

**Dimension Total: 79/90**

---

## Phase 3: Blindspot Hunt

### Blindspot 1: No validation against real projects
The proposal is motivated by pm-work-tracker failure but proposes no validation step. How will we know the changes work? The SCs measure that the pipeline generates different artifacts, not that those artifacts catch more bugs. A proper validation would include: regenerate tests for pm-work-tracker with new pipeline and confirm the milestone map bug is caught.

**Quote**: "每个包含父子实体关系的 feature，gen-journeys 至少生成 1 个 Golden Path Journey" — this measures generation, not bug-catching ability.

### Blindspot 2: "Business result" assertion classification is undefined
SC-3's core metric depends on classifying assertions as "business result" vs "structural." This classification is subjective and the proposal provides no rubric, examples, or decision tree for making this classification. Without it, SC-3 is unimplementable as a verifiable criterion.

**Quote**: "≥80% 的断言验证业务结果（实体存在、状态正确、关系完整），而非仅 HTTP 状态码" — "实体存在" could be verified by checking HTTP 200 + response contains entity, which is arguably structural.

### Blindspot 3: Sequential dependency disguised as parallelism
"每个 skill 的修改是独立的，可以并行" — this is false. gen-test-scripts reads fixture specs from gen-contracts output. gen-contracts reads workflow descriptions from gen-journeys output. These are sequentially dependent by definition. The proposal contradicts itself on this point.

**Quote**: "每个 skill 的修改是独立的，可以并行" vs "从 Contract 读取 fixture spec，生成丰富 fixture" — gen-test-scripts depends on gen-contracts output, making them not independent.

### Blindspot 4: No migration path for existing projects
Projects that already have generated tests (like pm-work-tracker) would need regeneration. The proposal does not address: (a) whether existing tests become invalid, (b) whether regenerating breaks existing test baselines, (c) whether projects need manual intervention to adopt the new fixture spec format.

### Blindspot 5: Fixture Specification may not be expressible for all domains
The proposal assumes fixture needs can always be declaratively specified. But some fixture needs are dynamic (e.g., "create a user who has performed at least 5 actions in the last 24 hours"). The Fixture Specification model may not handle temporal or state-dependent fixtures.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 98 | 110 |
| Solution Clarity | 108 | 120 |
| Industry Benchmarking | 95 | 120 |
| Requirements Completeness | 87 | 110 |
| Solution Creativity | 82 | 100 |
| Feasibility | 84 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 71 | 90 |
| Success Criteria | 64 | 80 |
| Logical Consistency | 79 | 90 |
| **Total** | **840** | **1000** |
