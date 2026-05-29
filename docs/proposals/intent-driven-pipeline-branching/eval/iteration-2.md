# Eval-Proposal Report: Intent-Driven Pipeline Branching

**Iteration**: 2
**Date**: 2026-05-29
**Expert**: CTO
**Final Score**: 896/1000

---

## Phase 1 — Reasoning Audit

### Iteration-1 Issues Addressed

| Issue | Status | Evidence |
|-------|--------|----------|
| [blindspot 1] Data flow gap — intent not reachable from BuildIndex() | **Resolved** | Lines 109-114: Complete data flow described — CLI handler calls `proposal.FindBySlug(slug)`, assigns `Proposal.Intent` to `BuildIndexOpts.Intent`, passes to `BuildIndex(opts)`. Explicitly states "BuildIndex() 内部不再重复解析 frontmatter，直接使用 opts.Intent" |
| [blindspot 2] autogen.go assumes business tasks exist | **Resolved** | Line 120: "零 business task 保护" paragraph added — "若 intent 为 refactor/cleanup 但 business task 列表为空...则不生成 validate-code/clean-code 等下游任务" |
| [blindspot 3] No rollback plan | **Resolved** | Lines 180-186: Full "Rollback Plan" section with feature flag (`FORGE_INTENT_BRANCHING=false`), autogen.go fallback, skill-layer safety, and verification steps |
| [blindspot 4] Stage-gate generation coupling with needsTest | **Resolved** | Line 168: Explicit scope item "将 stage-gate 生成逻辑从 needsTest 条件中独立出来，改为由 autogen.go 的依赖接线逻辑统一控制" |
| [blindspot 5] proposal.md may not exist at index time | **Resolved** | Line 69: "若 proposal.md 不存在...proposal.FindBySlug() 返回空 Proposal，此时 CLI handler 将 opts.Intent 设为默认值 new-feature，行为与当前一致" |
| [Self-contradiction] brainstorm template intent default unclear | **Partially resolved** | The proposal clarifies that brainstorm's AI infers intent then user confirms, but the template's initial state (blank? `new-feature`?) is still unspecified |
| [Self-contradiction] readFrontmatter() vs proposal.FindBySlug() duplicate parsing | **Resolved** | Lines 109-114: Eliminates readFrontmatter() entirely — reuses `proposal.FindBySlug()` which already parses intent from frontmatter |

### Problem -> Solution Trace

The reasoning chain remains coherent after revision:
1. **Problem**: Test pipeline treats all coding tasks identically; refactor/cleanup produce meaningless journeys.
2. **Evidence**: Verified code references (`IsTestableType()`, gen-journeys SKILL.md rule), concrete score data (466->630->585).
3. **Solution**: Intent field drives pipeline topology; refactor/cleanup skip test pipeline entirely.
4. **Data flow**: Now fully traced — CLI handler -> `proposal.FindBySlug()` -> `BuildIndexOpts.Intent` -> `needsTestPipeline(taskType, intent)` -> autogen.go dependency wiring.
5. **Edge cases**: Zero-business-task scenario handled; proposal-not-found scenario handled; rollback plan exists.

**Verdict**: All five critical blindspots from iteration-1 have been addressed. The proposal is now significantly more implementation-ready.

### Remaining Reasoning Concerns

1. **brainstorm template initial state**: The proposal says brainstorm adds an "intent frontmatter field" to the template, and AI infers intent during brainstorm. But the template itself — before any brainstorm runs — what does it contain? If a user creates a proposal manually without brainstorm, the field is absent, which triggers the default-new-feature fallback (acceptable). But if a user copies the template and runs brainstorm, does brainstorm overwrite a blank field or a pre-filled `new-feature`? This is a minor ambiguity.

2. **autogen.go complexity budget**: The proposal now describes 6 distinct wiring scenarios (new-feature Breakdown, new-feature Quick, refactor Breakdown, refactor Quick, cleanup Quick, zero-business-task). Each has subtly different dependency chains. The proposal treats this as manageable, but the combinatorial growth is real — adding a new intent or mode would multiply scenarios. The proposal does not acknowledge this scaling concern.

---

## Phase 2 — Rubric Scoring

### 1. Problem Definition: 102/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | Core problem is unambiguous: refactor/cleanup tasks have no observable behavior, yet the test pipeline forces them through Journey generation. The statement separates root cause (wrong pipeline type) from symptoms (eval scores, fact hallucination). Minor deduction: the opening sentence is dense with parenthetical references — a cleaner narrative flow would help. |
| Evidence provided | 37/40 | Four concrete evidence points: eval score progression (466->630->585), gen-journeys SKILL.md hard rule, `IsTestableType()` code reference, PRD user stories format mismatch. All verifiable against source. The `build.go` reference is accurate. Minor deduction: the eval scores are from a single feature — while the proposal argues the problem is systematic, one more example would strengthen this. |
| Urgency justified | 27/30 | "3 rounds of wasted eval" plus "every refactor/cleanup feature will hit this" is compelling. The systemic argument ("每多一个此类特征，就多一轮无效的 journey 生成") is strong. Minor deduction: still no quantification of actual cost (API calls, compute time, person-hours). |

### 2. Solution Clarity: 108/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 39/40 | The 2D matrix with exact pipeline paths per cell is highly concrete. Each cell lists exact skill names. The data flow (CLI -> proposal.FindBySlug -> BuildIndexOpts.Intent) is now fully traced. The autogen.go wiring scenarios are specified with exact `depends_on` assignments. |
| User-facing behavior described | 36/45 | The Key Scenarios describe the user flow from proposal creation through pipeline execution. The CLI output is implied (task list changes) but not explicitly described — what does the user see differently when running `forge task index` on a refactor vs new-feature proposal? The brainstorm interaction (AI proposes intent, user confirms via AskUserQuestion) is well described. |
| Technical direction clear | 33/35 | Significantly improved from iteration-1. The data flow is explicit, the function signatures are clear (`needsTestPipeline(taskType, intent)`), the autogen.go wiring logic is detailed with exact scenarios. The only gap: the proposal mentions `BuildIndexOpts.Intent` but doesn't show the struct definition change. |

### 3. Industry Benchmarking: 88/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 28/40 | GitHub Actions `paths`, Bazel `test_suite`/`query`, Maven/Gradle lifecycle phases are referenced. The comparison goes deeper than iteration-1 — each tool's mechanism is named (paths filter, dependency graph pruning, declarative phase skipping) and contrasted with Forge's approach. However, no specific open-source projects beyond tool names are cited, and no published patterns/papers are referenced. |
| At least 3 meaningful alternatives | 24/30 | Four alternatives listed: do nothing, gen-journeys internal skip, fix gen-journeys verification, and intent-driven branching. Each has pros/cons with verdicts. The "fix gen-journeys verification" alternative is still somewhat dismissed as "方向错误" without fully exploring whether improved gen-journeys + selective pipeline could be a hybrid approach. |
| Honest trade-off comparison | 18/25 | Comparison table with pros/cons exists. The "需要修改 forge-cli + 多个 skill" con for the chosen approach is still somewhat understated — the actual change surface is 7+ files with complex dependency rewiring in autogen.go. The proposal calls it "中等规模" but the 6 wiring scenarios in autogen.go suggest higher complexity. |
| Chosen approach justified | 18/25 | "最彻底，且复杂度可控" remains the core justification. Improved from iteration-1 with the addition of specific comparison against each alternative. However, still no quantitative comparison (estimated LOC, person-days, risk scores). |

### 4. Requirements Completeness: 98/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 36/40 | Six key scenarios covering all matrix cells plus the ambiguous-intent edge case. The iteration-2 addition of explicit handling for proposal.md-not-found and zero-business-task scenarios strengthens coverage. Remaining gap: what happens if a user manually edits proposal.md to change intent after some tasks are already generated? Does `forge task index` regenerate the full task list? |
| Non-functional requirements | 36/40 | Backward compatibility (default new-feature) and minimal invasiveness are well stated. The proposal now explicitly identifies the constraint "不改变 task 类型定义和 status 状态机". Missing: no mention of performance impact (does reading intent from proposal.md add latency to `forge task index`?). |
| Constraints & dependencies | 26/30 | Significantly improved. proposal.md-not-found edge case now documented. Quality-gate hook confirmed ready. Intent Propagation reuse identified. Remaining gap: the proposal says "brainstorm 阶段的 AI 推断完成，用户确认后写入 proposal frontmatter" but doesn't specify what happens if the user runs `forge task index` before brainstorm completes — the proposal.md exists but has no intent field. This is partially covered by the default-new-feature fallback, but the brainstorm-incomplete scenario is distinct from the proposal-not-found scenario. |

### 5. Solution Creativity: 80/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 34/40 | "Intent as first-class pipeline citizen" is a meaningful abstraction beyond what CI/CD tools do. The 2D matrix (Mode x Intent) is cleaner than file-path-based triggers. "Regression validation replaces journey validation" is a genuine insight. The distinction between "semantic signal" (intent) and "structural signal" (file paths/dependency graphs) is well articulated. |
| Cross-domain inspiration | 24/35 | CI/CD patterns are referenced. The proposal could have drawn from compiler optimization levels (different passes for different code patterns), test impact analysis (selective test execution), or feature flag systems (conditional execution paths). The "Assumptions Challenged" table is a nice touch that shows structured thinking. |
| Simplicity of insight | 22/25 | "Different work types need different validation strategies" is elegant and obvious in retrospect. The 2D matrix is a clean mental model. Minor deduction: the 6 wiring scenarios in autogen.go suggest the implementation is less simple than the insight. |

### 6. Feasibility: 85/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 34/40 | Feasible with clear implementation plan. The skill-layer changes are straightforward. The CLI changes are now well-specified with data flow, function signatures, and wiring logic. Remaining concern: the 6 distinct wiring scenarios in autogen.go create a combinatorial testing burden — the proposal identifies this risk but underestimates the testing matrix size. |
| Resource & timeline feasibility | 26/30 | "中等规模改动" is fair — 7+ files across skill docs and Go code. No timeline estimate provided (same as iteration-1). The scope is realistic for a focused iteration. |
| Dependency readiness | 25/30 | Quality-gate hook confirmed ready. Intent Propagation in breakdown-tasks/quick-tasks confirmed. The proposal now correctly notes that Intent Propagation is being "extended to pipeline level" rather than merely "reused" — this is more honest than iteration-1. |

### 7. Scope Definition: 70/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 27/30 | Each item specifies file and change. The iteration-2 additions (stage-gate decoupling, data flow specification, autogen.go wiring) make the scope more concrete. The in-scope list now includes the critical "stage-gate 生成逻辑从 needsTest 条件中独立出来" item that was missing in iteration-1. |
| Out-of-scope explicitly listed | 22/25 | Six out-of-scope items listed, including gen-journeys code verification, mixed-intent support, and intent inference accuracy optimization. The addition of "new-feature 场景下的事实幻觉问题" as explicitly out-of-scope is good — it acknowledges the limitation without scope-creeping. Missing: migrating existing in-flight features. |
| Scope is bounded | 21/25 | Bounded to specific files and functions. The proposal now identifies the cascade (build.go -> needsTestPipeline -> autogen.go -> downstream tasks) more explicitly. However, the actual blast radius of the autogen.go changes — particularly the "stage-gate decoupling from needsTest" — could affect any code path that calls `needsTestPipeline()` or checks its return value. The proposal doesn't list all callers of `needsTestPipeline()`. |

### 8. Risk Assessment: 82/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 27/30 | Six risks identified (up from four in iteration-1). New additions: "stage-gate 生成与 needsTest 解耦引入回归" and "autogen.go 重构引入依赖接线 bug". These directly address iteration-1 blindspots 4 and the autogen.go complexity concern. The risk table is now comprehensive. |
| Likelihood + impact rated | 27/30 | Ratings are honest. "Intent 推断错误导致新功能跳过 journey" is L/H (correct). "autogen.go 重构引入依赖接线 bug" is M/H (correct). The ratings are internally consistent. |
| Mitigations are actionable | 28/30 | Mitigations are concrete: "AI 推断后用户确认", "默认 new-feature", "通过回归测试验证: 对 intent: new-feature 生成完整管道, 对 intent: refactor 生成跳过测试但保留 stage-gate 的新依赖链". The rollback plan with feature flag (`FORGE_INTENT_BRANCHING=false`) is directly actionable. The mitigation for autogen.go wiring bugs specifies "修改前为 autogen.go 添加单元测试覆盖所有 4 种接线场景" — this is testable and specific. |

### 9. Success Criteria: 72/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 27/30 | Nine success criteria, most testable: "不生成 gen-journeys 任务" (check task index), "quality-gate hook 执行 compile+fmt+lint+test" (run hook), "spec-only PRD（无 prd-user-stories.md 文件）" (check file absence). The addition of criterion 9 (stage-gate dependency chain verification) directly addresses iteration-1's concern about stage-gate coupling. |
| Coverage is complete | 23/25 | Criteria cover all three intents, both modes, backward compatibility, and brainstorm inference. The addition of brainstorm inference as criterion 1 addresses iteration-1's gap. The stage-gate dependency chain criterion (SC-9) ensures the autogen.go changes are verified. Remaining gap: no criterion for the feature flag rollback mechanism — can you verify that setting `FORGE_INTENT_BRANCHING=false` restores original behavior? |
| SC internal consistency | 22/25 | The SC set is internally consistent. SC-1 (brainstorm inference) + SC-2 (refactor PRD) + SC-6 (no gen-journeys tasks) form a coherent chain. SC-9 (stage-gate dependencies) ensures the autogen.go wiring is correct. One concern: SC-8 (backward compatibility) says "行为不变" but the autogen.go changes modify the stage-gate generation path — the "不变" claim depends on intent defaulting to new-feature correctly in all code paths. The proposal addresses this but the SC could be more specific about what "不变" means (same task list? same dependency chain?). |

### 10. Logical Consistency: 83/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 33/35 | The solution directly solves the problem: refactor/cleanup skip the journey pipeline. The spec-only PRD addresses the empty-user-stories sub-problem. The data flow is now fully traced. The iteration-1 gap about "fact hallucination not addressed for refactor" is acknowledged in out-of-scope ("new-feature 场景下的事实幻觉问题 — 本 proposal 通过跳过 refactor/cleanup 的 journey 管道绕过此问题"). This is honest scoping. |
| Scope <-> Solution <-> SC aligned | 26/30 | Significantly improved. The scope now includes the stage-gate decoupling that iteration-1 flagged as missing. The SC set covers all scope items. One remaining misalignment: the scope lists "autogen.go: GetBreakdownTestTasks/GetQuickTestTasks 根据 intent 跳过" but the solution text (line 167) says these functions "本身不需要修改 — 它们由 needsTestPipeline() 的返回值控制是否被调用". This is correct but the scope item's wording is slightly misleading. |
| Requirements <-> Solution coherent | 24/25 | Key scenarios map cleanly to solution components. The coding.fix mapping rule addresses the refactor/fix boundary. No orphan requirements. The rollback plan adds a new dimension not in original requirements but is a net positive. |

---

## Phase 3 — Blindspot Hunt

### [blindspot] 1: Scaling concern — combinatorial explosion of wiring scenarios

The proposal defines 6 wiring scenarios in autogen.go (new-feature Breakdown/Quick, refactor Breakdown/Quick, cleanup Quick, zero-business-task). Each has different dependency chains. The proposal treats this as manageable, but does not acknowledge that adding a new intent type (e.g., `hotfix`, `experiment`) or a new mode would multiply scenarios. Quote: "具体接线逻辑：" followed by 4 bullet points — this is a maintenance burden that will grow.

**What must improve**: Acknowledge the combinatorial growth explicitly. Consider whether a declarative dependency config (e.g., a wiring table or DSL) would be more maintainable than imperative if/else branches in autogen.go.

### [blindspot] 2: Feature flag testing gap

The rollback plan introduces `FORGE_INTENT_BRANCHING` feature flag, but no success criterion verifies that the flag works. Quote from Rollback Plan: "CLI handler 通过环境变量 FORGE_INTENT_BRANCHING=false 禁用 intent 分支". The success criteria section does not include a testable criterion for this flag.

**What must improve**: Add a success criterion: "Setting FORGE_INTENT_BRANCHING=false causes forge task index to produce identical output as the pre-intent code path for all intent types."

### [blindspot] 3: brainstorm template initial intent value unspecified

The proposal says "brainstorm template: 添加 intent frontmatter 字段" (line 102) and "AI 推断 intent 的步骤，用 AskUserQuestion 确认" (line 103). But the template's initial state is unspecified — does it ship with `intent: new-feature` as default, or `intent: ""` (blank)? If blank, what happens if a user creates a proposal manually and runs `forge task index` without brainstorm? The proposal.md-not-found case is handled (default new-feature), but the proposal-exists-without-intent case relies on the same fallback without explicitly stating this is the same code path.

**What must improve**: State explicitly: "The brainstorm template includes `intent: new-feature` as the default value. Manual proposals without the field fall through to the same default in BuildIndex()."

### [blindspot] 4: autogen.go wiring "zero business task" protection is incomplete

Line 120: "若 intent 为 refactor/cleanup 但 business task 列表为空（例如纯文档类型特征），则不生成 validate-code/clean-code 等下游任务". This protection is only described for refactor/cleanup — but what about new-feature with zero business tasks? The current code presumably handles this (or doesn't generate business tasks for doc-only features), but the proposal doesn't confirm that the existing behavior is preserved under the new logic.

**What must improve**: Confirm that the zero-business-task protection does not change existing behavior for new-feature intent, or explicitly state that zero business tasks is impossible for new-feature features.

---

## Score Summary

| Dimension | Score | Max | Delta from Iter-1 |
|-----------|-------|-----|--------------------|
| Problem Definition | 102 | 110 | +7 |
| Solution Clarity | 108 | 120 | +10 |
| Industry Benchmarking | 88 | 120 | +6 |
| Requirements Completeness | 98 | 110 | +10 |
| Solution Creativity | 80 | 100 | +5 |
| Feasibility | 85 | 100 | +7 |
| Scope Definition | 70 | 80 | +8 |
| Risk Assessment | 82 | 90 | +12 |
| Success Criteria | 72 | 80 | +14 |
| Logical Consistency | 83 | 90 | +9 |
| **Total** | **896** | **1000** | **+148** |

**Target**: 900
**Gap**: 4 points
