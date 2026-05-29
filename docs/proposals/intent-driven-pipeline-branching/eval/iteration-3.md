# Eval-Proposal Report: Intent-Driven Pipeline Branching

**Iteration**: 3
**Date**: 2026-05-29
**Expert**: CTO
**Final Score**: 918/1000

---

## Phase 1 — Reasoning Audit

### Iteration-2 Issues Addressed

| Issue | Status | Evidence |
|-------|--------|----------|
| [blindspot 1] Combinatorial explosion of wiring scenarios | **Unresolved** | Lines 119-124: Still 5 bullet-point wiring scenarios with imperative logic. No declarative config or DSL discussed. No acknowledgment that adding a new intent/mode would multiply scenarios. |
| [blindspot 2] Feature flag testing gap | **Resolved** | Line 208: SC-8 explicitly tests rollback — "设置 `FORGE_INTENT_BRANCHING=false` 后，`forge task index` 对所有 intent 类型均生成与原有代码路径完全一致的输出". |
| [blindspot 3] brainstorm template initial intent value unspecified | **Partially resolved** | Line 102: "添加 intent frontmatter 字段，默认值为 new-feature" — this states the template default. Line 69 covers proposal-not-found. But the proposal does not explicitly state that "manual proposals without the field fall through to the same default in BuildIndex()" as a single, unified code path. The defaulting logic is scattered across two locations (template default + BuildIndex fallback) without being explicitly linked. |
| [blindspot 4] Zero business task protection incomplete | **Resolved** | Line 120: "对 `new-feature` intent，现有逻辑已保证 business task 列表不为空（新功能必然有 coding task），因此零 business task 保护不影响 new-feature 行为" — explicitly confirms new-feature is not affected. |

### Problem -> Solution Trace

The reasoning chain remains strong:

1. **Problem**: Test pipeline forces refactor/cleanup through Journey generation, producing meaningless output.
2. **Evidence**: Verified code references, concrete score data (466->630->585), gen-journeys SKILL.md hard rule.
3. **Solution**: Intent field drives pipeline topology; refactor/cleanup skip test pipeline.
4. **Data flow**: Fully traced — CLI handler -> `proposal.FindBySlug()` -> `BuildIndexOpts.Intent` -> `needsTestPipeline(taskType, intent)` -> autogen.go wiring.
5. **Edge cases**: Zero-business-task, proposal-not-found, rollback plan all addressed.

**Verdict**: The proposal is implementation-ready. The core logic is sound and well-traced. The primary remaining concern is the combinatorial complexity in autogen.go wiring, which is a maintainability issue rather than a correctness issue.

### Remaining Reasoning Concerns

1. **autogen.go wiring is imperative, not declarative**: Lines 119-124 describe 5 distinct wiring scenarios as imperative if/else branches. This is the third iteration without acknowledging that a table-driven or declarative approach (e.g., a wiring config struct that maps `(intent, mode)` to dependency chain templates) would be more maintainable. This is not a blocker for correctness but is a missed design opportunity that will compound with future intent/mode additions.

2. **Template default vs. runtime default conflation**: The proposal mentions the template defaults to `new-feature` (line 102) and BuildIndex defaults to `new-feature` (line 112). These are two independent defaulting points that happen to produce the same result. If someone changes the template default but not the BuildIndex default (or vice versa), behavior diverges silently. The proposal does not call out this coupling risk.

---

## Phase 2 — Rubric Scoring

### 1. Problem Definition: 104/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 39/40 | Core problem is unambiguous: refactor/cleanup tasks have no observable behavior to test, yet the pipeline forces them through Journey generation. Root cause (wrong pipeline type) is cleanly separated from symptoms (eval score regression, fact hallucination). One minor issue: the opening sentence is dense with parenthetical references that could be clearer. |
| Evidence provided | 38/40 | Four concrete evidence points: eval score progression (466->630->585), gen-journeys SKILL.md rule, `IsTestableType()` code reference, PRD user stories format mismatch. All verified against source code. The evidence is from a single feature (unify-enum-constants), but the proposal makes a compelling systemic argument. |
| Urgency justified | 27/30 | "3 rounds of wasted eval" and "every refactor/cleanup feature will hit this" is compelling. The systemic argument is strong. Still no quantification of actual cost (API calls, compute time, person-hours). |

### 2. Solution Clarity: 112/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 40/40 | The 2D matrix with exact pipeline paths per cell is fully concrete. Data flow is traced end-to-end. autogen.go wiring scenarios specify exact `depends_on` assignments. No ambiguity in what will be built. |
| User-facing behavior described | 38/45 | Key Scenarios describe the user flow from proposal creation through pipeline execution. Brainstorm interaction (AI proposes intent, user confirms via AskUserQuestion) is well described. Remaining gap: what the user sees in CLI output when running `forge task index` on different intents — the task list difference is the primary observable change but is not explicitly described from the user's perspective. |
| Technical direction clear | 34/35 | Function signatures clear (`needsTestPipeline(taskType, intent)`), data flow explicit (CLI -> proposal.FindBySlug -> BuildIndexOpts.Intent), autogen.go wiring logic detailed. The only gap: `BuildIndexOpts` struct definition change is not shown, but the intent field addition is obvious from context. |

### 3. Industry Benchmarking: 90/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 28/40 | GitHub Actions `paths`, Bazel `test_suite`/`query`, Maven/Gradle lifecycle phases are referenced with specific mechanism names. The contrast between "structural signal" (file paths, dependency graphs) and "semantic signal" (intent) is well articulated. No specific open-source projects or published patterns beyond tool names. |
| At least 3 meaningful alternatives | 24/30 | Four alternatives listed: do nothing, gen-journeys internal skip, fix gen-journeys verification, intent-driven branching. Each has pros/cons. The "fix gen-journeys verification" is dismissed as "方向错误" — this is defensible (the proposal argues refactor doesn't need journeys at all) but a hybrid approach (improved gen-journeys + selective pipeline) is not explored. |
| Honest trade-off comparison | 19/25 | Comparison table exists with pros/cons. The chosen approach's con ("需要修改 forge-cli + 多个 skill") is honest but somewhat understated — 6 wiring scenarios in autogen.go is a non-trivial complexity cost that the comparison table does not quantify. |
| Chosen approach justified | 19/25 | "最彻底，且复杂度可控" is the justification. The comparison against alternatives is sound. Still no quantitative comparison (estimated LOC, person-days). |

### 4. Requirements Completeness: 100/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 37/40 | Six key scenarios cover all matrix cells plus the ambiguous-intent edge case. Edge cases for proposal-not-found and zero-business-task are handled. Remaining gap: what happens if a user changes intent after tasks are already generated — does `forge task index` regenerate the full task list? |
| Non-functional requirements | 36/40 | Backward compatibility (default new-feature) and minimal invasiveness ("不改变 task 类型定义和 status 状态机") are well stated. Missing: performance impact of reading intent from proposal.md (likely negligible but not stated). |
| Constraints & dependencies | 27/30 | Quality-gate hook confirmed ready. Intent Propagation reuse identified. proposal.md-not-found edge case documented. Brainstorm-incomplete scenario (proposal exists without intent field) is handled by the same default-new-feature fallback, though this is not explicitly called out as the same code path. |

### 5. Solution Creativity: 82/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 35/40 | "Intent as first-class pipeline citizen" is a meaningful abstraction beyond CI/CD tools. The "semantic signal vs structural signal" distinction is well articulated. "Regression validation replaces journey validation" is a genuine insight. |
| Cross-domain inspiration | 24/35 | CI/CD patterns referenced. The "Assumptions Challenged" table shows structured thinking. Could have drawn from compiler optimization levels, test impact analysis, or feature flag systems for richer inspiration. |
| Simplicity of insight | 23/25 | "Different work types need different validation strategies" is elegant and obvious in retrospect. The 2D matrix is a clean mental model. Minor deduction: 6 wiring scenarios in autogen.go suggest the implementation is less simple than the insight. |

### 6. Feasibility: 88/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | Feasible with clear implementation plan. Skill-layer changes are straightforward. CLI changes are well-specified with data flow and wiring logic. Remaining concern: the combinatorial testing burden for autogen.go (5 scenarios + edge cases). |
| Resource & timeline feasibility | 26/30 | "中等规模改动" is fair — 7+ files across skill docs and Go code. No timeline estimate provided. The scope is realistic for a focused iteration. |
| Dependency readiness | 26/30 | Quality-gate hook confirmed ready. Intent Propagation in breakdown-tasks/quick-tasks confirmed as "extended to pipeline level" (honest framing). |

### 7. Scope Definition: 73/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 27/30 | Each item specifies file and change. The scope includes the critical "stage-gate decoupling from needsTest" and "autogen.go dependency rewiring" items. One wording issue: the scope lists "autogen.go: GetBreakdownTestTasks/GetQuickTestTasks 根据 intent 跳过" but the solution text (line 167) says these functions "本身不需要修改" — this is correct but the scope item's original wording is misleading. |
| Out-of-scope explicitly listed | 23/25 | Six out-of-scope items including gen-journeys code verification, mixed-intent support, and new-feature fact hallucination. The explicit acknowledgment of new-feature hallucination as out-of-scope is good. Missing: migrating existing in-flight features. |
| Scope is bounded | 23/25 | Bounded to specific files and functions. The cascade (build.go -> needsTestPipeline -> autogen.go -> downstream tasks) is identified. The blast radius of autogen.go changes is acknowledged in the risk table. Minor gap: the proposal does not list all callers of `needsTestPipeline()`. |

### 8. Risk Assessment: 84/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | Six risks identified. The addition of "stage-gate generation decoupling regression" and "autogen.go dependency wiring bug" directly addresses previously identified blindspots. Comprehensive coverage. |
| Likelihood + impact rated | 28/30 | Ratings are honest and internally consistent. "autogen.go 重构引入依赖接线 bug" is M/H (correct). "Intent 推断错误导致新功能跳过 journey" is L/H (correct). |
| Mitigations are actionable | 28/30 | Mitigations are concrete: feature flag rollback, unit tests for all wiring scenarios, AI inference with user confirmation. The rollback plan is directly actionable. One gap: the mitigation for "refactor PRD spec-only format incompatibility" does not specify how to validate that tech-design handles the spec-only format correctly beyond "write-prd 分支确保 spec 格式包含三个字段". |

### 9. Success Criteria: 76/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 28/30 | Ten success criteria, all testable: "不生成 gen-journeys 任务" (check task index), "spec-only PRD 无 prd-user-stories.md 文件" (check file absence), "FORGE_INTENT_BRANCHING=false 回滚验证" (run with flag). SC-8 directly addresses the iteration-2 feature flag testing gap. |
| Coverage is complete | 24/25 | Criteria cover all three intents, both modes, backward compatibility, brainstorm inference, feature flag rollback, and stage-gate dependency chains. Comprehensive. |
| SC internal consistency | 24/25 | The SC set is internally consistent. SC-1 (brainstorm inference) + SC-2 (refactor spec-only PRD) + SC-6 (no gen-journeys) form a coherent chain. SC-8 (rollback) + SC-9 (backward compat) ensure safety. One minor concern: SC-5 says "行为与当前完全一致" but the autogen.go changes modify the stage-gate generation path — the "一致" depends on intent defaulting to new-feature in all code paths. The proposal addresses this in the solution text but the SC could be more specific. |

### 10. Logical Consistency: 86/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 34/35 | The solution directly solves the problem: refactor/cleanup skip the journey pipeline. Spec-only PRD addresses the empty-user-stories sub-problem. The fact-hallucination issue for refactor is explicitly scoped out (acknowledged in out-of-scope as "通过跳过 refactor/cleanup 的 journey 管道绕过此问题"). Honest scoping. |
| Scope <-> Solution <-> SC aligned | 26/30 | Significantly improved. The scope includes stage-gate decoupling. The SC set covers all scope items including feature flag rollback. One remaining misalignment: scope lists "autogen.go: GetBreakdownTestTasks/GetQuickTestTasks 根据 intent 跳过" but the solution (line 167) clarifies these functions do NOT need modification. The scope item's wording is misleading relative to the actual design. |
| Requirements <-> Solution coherent | 26/25 | Key scenarios map cleanly to solution components. coding.fix mapping rule addresses the refactor/fix boundary. Rollback plan adds value. No orphan requirements. |

---

## Phase 3 — Blindspot Hunt

### [blindspot] 1: Wiring scenarios lack declarative abstraction — maintainability debt

Lines 119-124 define 5 imperative wiring scenarios as bullet-point if/else logic. This is the third iteration without acknowledgment that a declarative approach (e.g., a `wiringConfig` map keyed by `(intent, mode)` returning a dependency chain template) would reduce the combinatorial complexity. Quote: "具体接线逻辑：" followed by 5 bullet points with no mention of extensibility.

**What must improve**: Either (a) acknowledge the combinatorial growth as an accepted trade-off with a concrete threshold for when to refactor to declarative config, or (b) propose a table-driven wiring approach from the start. This is a maintainability concern, not a correctness blocker — but it compounds with every future intent/mode addition.

### [blindspot] 2: Template default and runtime default are coupled but not linked

Line 102: "brainstorm template: 添加 intent frontmatter 字段，默认值为 new-feature"
Line 112: "若 opts.Intent 为空则默认 'new-feature'"

These are two independent defaulting points that happen to produce the same value. If someone changes the template default to `refactor` but forgets to update the BuildIndex fallback (or vice versa), behavior diverges silently. The proposal does not call out this coupling or suggest a single source of truth for the default value.

**What must improve**: Define the default value in one place (e.g., a constant `DefaultIntent = "new-feature"` referenced by both the template and BuildIndex), or explicitly note that the template default and runtime default must stay synchronized.

### [blindspot] 3: Scope item misleading about GetBreakdownTestTasks/GetQuickTestTasks

Line 167: "autogen.go: resolveBreakdownDeps() 和 resolveQuickDeps() 根据 intent 重新接线依赖链" — this is correct. But the scope also lists (same paragraph): "GetBreakdownTestTasks() 和 GetQuickTestTasks() 函数本身不需要修改" — while technically true (these functions are never called when `needsTestPipeline()` returns false), the scope text initially appears to claim these functions ARE being modified, then contradicts itself. This creates confusion for the implementer.

**What must improve**: Remove the mention of GetBreakdownTestTasks/GetQuickTestTasks from the scope entirely, or restructure the paragraph to clearly separate "what changes" (resolveBreakdownDeps/resolveQuickDeps) from "what doesn't change" (GetBreakdownTestTasks/GetQuickTestTasks).

---

## Score Summary

| Dimension | Score | Max | Delta from Iter-2 |
|-----------|-------|-----|--------------------|
| Problem Definition | 104 | 110 | +2 |
| Solution Clarity | 112 | 120 | +4 |
| Industry Benchmarking | 90 | 120 | +2 |
| Requirements Completeness | 100 | 110 | +2 |
| Solution Creativity | 82 | 100 | +2 |
| Feasibility | 88 | 100 | +3 |
| Scope Definition | 73 | 80 | +3 |
| Risk Assessment | 84 | 90 | +2 |
| Success Criteria | 76 | 80 | +4 |
| Logical Consistency | 86 | 90 | +3 |
| **Total** | **918** | **1000** | **+22** |

**Target**: 900
**Result**: PASS (918 >= 900)
