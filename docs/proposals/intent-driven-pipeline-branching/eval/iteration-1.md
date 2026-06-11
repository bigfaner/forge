# Eval-Proposal Report: Intent-Driven Pipeline Branching

**Iteration**: 1
**Date**: 2026-05-29
**Expert**: CTO
**Final Score**: 748/1000

---

## Phase 1 — Reasoning Audit

### Problem -> Solution Trace

The argument chain is:
1. **Problem**: Forge test pipeline treats all coding tasks identically, forcing refactor/cleanup through Journey generation that cannot produce meaningful results.
2. **Evidence**: unify-enum-constants eval-journey scores 466->630->585 (verified in `tasks/records/eval-journey.md`), gen-journeys has a hard rule forbidding code reading, `IsTestableType()` returns true for all `coding.*`.
3. **Solution**: Intent field in proposal frontmatter drives pipeline topology selection; refactor/cleanup skip the entire test pipeline.
4. **Verification**: Success criteria are measurable (no gen-journeys tasks generated, quality-gate hook still runs).

**Verdict**: The reasoning chain is coherent. The problem is real and verified. The solution directly addresses the root cause (wrong pipeline for wrong intent type).

### Solution -> Evidence Trace

- `IsTestableType()` confirmed: `build.go:519` — `strings.HasPrefix(typ, "coding.")` returns true for all `coding.*`.
- gen-journeys rule confirmed: `SKILL.md:210` — "Do not read source code, test files, or implementation files."
- Score progression confirmed: `eval-journey.md:6` — "466 -> 630 -> 585".
- `proposal.go` already has `Intent` field — foundation exists.
- Intent Propagation already exists in `breakdown-tasks/SKILL.md` and `quick-tasks/SKILL.md`.

**Gap**: The proposal claims `build.go` has `IsTestableType()` that returns true for `coding.refactor` and `coding.cleanup` — confirmed correct. However, the proposal also says `needsTestPipeline()` checks intent — this function currently only checks task types, not intent. The proposed change is to modify it to read intent, which is a new behavior.

### Self-Contradiction Check

1. **Matrix consistency**: The matrix in "Proposed Solution" is internally consistent — each cell defines a clear pipeline path.
2. **Scope vs. Solution**: In-scope items map cleanly to solution components.
3. **Contradiction found**: The proposal says "brainstorm template: add intent frontmatter field" but also says intent is determined by "AI inference during brainstorm" — this means the template must have intent as optional/blank initially, with AI filling it in. The proposal does not clarify whether the template pre-populates a default or leaves it empty.
4. **Minor tension**: The proposal says `build.go` reads intent via `readFrontmatter(proposalPath, "intent")` but `proposal.go` already parses intent. The proposal does not address whether to reuse `proposal.FindBySlug()` or implement a new read function, creating potential for duplicate parsing logic.

---

## Phase 2 — Rubric Scoring

### 1. Problem Definition: 95/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 36/40 | The core problem is unambiguous: "test pipeline treats all coding tasks identically, but refactor/cleanup have no observable behavior to test." Two readers would reach the same interpretation. Minor deduction: the problem statement mixes symptom (eval scores dropping) with root cause (wrong pipeline type). |
| Evidence provided | 36/40 | Concrete data: eval scores (466->630->585), specific code references (`IsTestableType()`, gen-journeys SKILL.md rule), specific feature name (`unify-enum-constants`). All evidence verified against source code. The `build.go` reference is accurate — `IsTestableType()` does use `strings.HasPrefix(typ, "coding.")`. |
| Urgency justified | 23/30 | "3 rounds of wasted eval" and "every refactor/cleanup feature will hit this" are compelling. However, the urgency section does not quantify the actual cost (how many eval API calls, how much compute time, how many person-hours wasted). |

### 2. Solution Clarity: 98/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | The 2D matrix (Mode x Intent) with explicit pipeline paths per cell is highly concrete. Each cell specifies exact skill names. |
| User-facing behavior described | 35/45 | The proposal describes internal pipeline topology, which is developer-tool behavior rather than end-user behavior. Key Scenarios partially compensate — they describe the user flow (create proposal -> AI infers intent -> confirm -> pipeline runs). But what the user *sees* during execution (CLI output, task list changes) is not described. Deduction: the "user-facing behavior" is mostly about what tasks appear in the index, which is observable but not described from the user's perspective. |
| Technical direction clear | 25/35 | Technical direction is clear for skill-layer changes. For CLI changes, the proposal specifies which files to modify (`build.go`, `autogen.go`) and what functions to change. However, the proposal introduces `readFrontmatter()` as a new function but `proposal.go` already has `FindBySlug()` that parses intent — the relationship between these is unclear. Also, `needsTestPipeline()` currently checks task types; the proposal says to "add intent parameter" but doesn't specify the API signature change or how intent flows from `BuildIndex()` to `needsTestPipeline()`. |

### 3. Industry Benchmarking: 82/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 25/40 | GitHub Actions `paths`, Bazel `test_suite`, Maven/Gradle lifecycle phases are mentioned. These are valid analogies but the comparison is superficial — one sentence each with no analysis of how their filtering mechanisms map to Forge's pipeline model. No specific open-source projects or published patterns are cited beyond naming these tools. |
| At least 3 meaningful alternatives | 22/30 | Four alternatives are listed: do nothing, gen-journeys internal skip, fix gen-journeys verification, and the proposed solution. These are genuinely different approaches. However, only the "do nothing" and the proposed solution are meaningfully analyzed. The "fix gen-journeys verification" alternative is labeled "方向错误" (wrong direction) without fully exploring whether a hybrid approach (improve gen-journeys + selective pipeline) might be better. |
| Honest trade-off comparison | 18/25 | The comparison table is present with pros/cons. However, the proposal's own cons ("需要修改 forge-cli + 多个 skill") underestimates the change surface — it's 7+ file modifications across skill docs and Go code, with dependency rewiring in `autogen.go`. The "M" (medium) complexity assessment in the comparison feels understated. |
| Chosen approach justified | 17/25 | "最彻底，且复杂度可控" (most thorough, manageable complexity) is the justification. This is a valid rationale but lacks quantitative comparison — how much more effort vs. the alternatives? What's the estimated LOC change? |

### 4. Requirements Completeness: 88/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 32/40 | Six key scenarios are identified covering all matrix cells plus the ambiguous-intent edge case. Missing scenarios: (1) What happens when a user changes intent mid-pipeline (e.g., realizes a refactor is actually a new feature)? (2) What happens when a Breakdown+cleanup is attempted despite the matrix saying "not applicable"? The matrix says `build.go` forces `mode=Quick` for cleanup, but no scenario describes this override behavior from the user's perspective. |
| Non-functional requirements | 35/40 | Backward compatibility (default to new-feature) and minimal invasiveness are well stated. Missing: performance impact of intent reading, any migration path for in-flight features that might have refactor tasks but no intent field. |
| Constraints & dependencies | 21/30 | Quality-gate hook dependency is noted as "already exists." Intent Propagation in breakdown-tasks/quick-tasks is noted as reusable. Missing constraints: (1) What if the proposal.md doesn't exist yet when `forge task index` runs? (2) The proposal says "intent 持久化在 proposal.md frontmatter" but `BuildIndex()` currently doesn't read proposal frontmatter — it reads feature task files. The data flow from proposal.md -> build.go is not traced. |

### 5. Solution Creativity: 75/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 30/40 | "Intent as first-class pipeline citizen" is a meaningful abstraction. The 2D matrix (Mode x Intent) is cleaner than the CI/CD analogies (which typically use file-path triggers). The innovation of "regression validation replaces journey validation" is a genuine insight. |
| Cross-domain inspiration | 22/35 | The proposal draws from CI/CD patterns (pipeline selection) but doesn't reference other domains. For example: compiler optimization levels (different passes for different code patterns), test impact analysis (selective test execution based on change graphs), or feature flag systems (conditional execution paths). |
| Simplicity of insight | 23/25 | The core insight — "different work types need different validation strategies" — is elegant and obvious in retrospect. The 2D matrix is a clean mental model. |

### 6. Feasibility: 78/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 30/40 | Mostly feasible. The skill-layer changes are straightforward (conditional logic in SKILL.md). The CLI changes are more complex: `autogen.go` dependency rewiring for refactor/cleanup requires careful handling of `lastRunID` lookup removal and direct business-task-to-stage-gate wiring. The proposal acknowledges this complexity but doesn't address the risk of the `resolveBreakdownDeps()` / `resolveQuickDeps()` functions becoming entangled with intent-aware branching. |
| Resource & timeline feasibility | 25/30 | "中等规模改动" (medium-scale change) is a fair assessment. No timeline estimate is provided — the proposal says "中等" but doesn't say "2 days" or "1 sprint." The scope (7+ files) is realistic for a focused iteration. |
| Dependency readiness | 23/30 | Quality-gate hook is confirmed ready. Intent Propagation is confirmed in both skill docs. However, the proposal's claim that "Intent Propagation 机制已在 breakdown-tasks / quick-tasks 中存在" is slightly misleading — the existing Intent Propagation maps intent to task *type*, not to pipeline *topology*. The proposal extends this concept significantly, not merely reuses it. |

### 7. Scope Definition: 62/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | Each in-scope item specifies the file and the change. "build.go: IsTestableType() 区分行为变更 vs 行为保持" is concrete. "autogen.go: GetBreakdownTestTasks() 和 GetQuickTestTasks() 根据 intent 跳过" is concrete. |
| Out-of-scope explicitly listed | 20/25 | Five out-of-scope items are listed, including the gen-journeys code verification and mixed-intent support. Good that deferred items are named. Missing: is there an out-of-scope item for "migrating existing features that already have refactor tasks"? |
| Scope is bounded | 17/25 | The scope is bounded to specific files and functions. However, the cascade of changes — `build.go` reads intent -> `needsTestPipeline()` uses intent -> `autogen.go` wires dependencies differently -> downstream task chain changes — means the actual blast radius may be larger than the 7 listed files. The proposal doesn't identify the full set of functions that call `needsTestPipeline()` or reference its output. |

### 8. Risk Assessment: 70/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 24/30 | Four risks are identified. Missing risks: (1) Risk of `autogen.go` regression — the dependency wiring for refactor is different from new-feature, and a bug in the refactor path could silently drop stage-gates. (2) Risk of intent field schema drift — if `coding.fix` type mapping rules change, the intent inference could become inconsistent. |
| Likelihood + impact rated | 22/30 | Ratings are present and mostly honest. "Intent 推断错误导致新功能跳过 journey" is correctly rated L/H. However, "refactor PRD spec-only 格式与下游 skill 不兼容" is rated M/M — this should be M/H because if tech-design or breakdown-tasks crash on missing user stories, it blocks the entire pipeline. |
| Mitigations are actionable | 24/30 | "AI 推断后用户确认" is actionable. "默认 new-feature" is actionable. "write-prd 分支确保 spec 格式包含三个字段" is actionable. However, the mitigation for "refactor PRD spec-only 格式不兼容" is described at a high level — it doesn't specify how to validate that tech-design handles the spec-only format correctly. |

### 9. Success Criteria: 58/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 22/30 | Most criteria are testable: "不生成 gen-journeys 任务" can be verified by checking the task index. "quality-gate hook 执行 compile+fmt+lint+test" can be verified by running the hook. However, "spec-only PRD（无 prd-user-stories.md 文件）" is a negative criterion (absence of file) — this is testable but the proposal doesn't specify what happens if write-prd accidentally generates user-stories for refactor. |
| Coverage is complete | 18/25 | Success criteria cover all three intents and both modes. Gaps: (1) No success criterion for the brainstorm skill adding intent inference. (2) No success criterion for the `build.go` intent reading mechanism itself. (3) The criterion "refactor/cleanup 特征的 business tasks 完成后，quality-gate hook 执行" overlaps with existing behavior — quality-gate already runs for all coding tasks. What's new here? |
| SC internal consistency | 18/25 | SC items are internally consistent as a set. One concern: SC-1 says "测试管道跳过 journey/contract/script" but SC-6 says "write-prd 对 intent: refactor 生成 spec-only PRD" — these are in different scopes (test pipeline vs. PRD format) but both are required for refactor to work correctly. If write-prd generates user-stories but the test pipeline is skipped, the PRD is inconsistent with the pipeline behavior. The proposal addresses this in the solution but the SC doesn't make this dependency explicit. |

### 10. Logical Consistency: 74/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 30/35 | The solution directly addresses the problem: refactor/cleanup skip the journey pipeline that was producing garbage. The spec-only PRD addresses the "PRD user stories are semantically empty for refactor" sub-problem. One gap: the problem statement mentions "生成器凭文档推测事实（常量名、CLI 命令未验证）" — the proposal's solution does not address the fact-hallucination issue for refactor; it only sidesteps it by skipping the pipeline. If a future refactor does need journeys (e.g., refactoring a public API), the hallucination problem persists. |
| Scope <-> Solution <-> SC aligned | 22/30 | Generally aligned. One misalignment: Scope lists "forge-cli build.go: needsTestPipeline() 读取 intent" and "autogen.go: GetBreakdownTestTasks/GetQuickTestTasks 根据 intent 跳过" as separate items, but the solution text describes `needsTestPipeline()` as the gate that controls whether `GenerateTestTasks()` is called at all. If `needsTestPipeline()` returns false, `GenerateTestTasks()` (which calls `GetBreakdownTestTasks()`/`GetQuickTestTasks()`) is never invoked — so modifying the Get* functions is redundant. The scope double-counts this change. |
| Requirements <-> Solution coherent | 22/25 | Key scenarios map to solution components. The "coding.fix type mapping rule" in the solution addresses the ambiguous boundary between fix and refactor. No orphan requirements found. |

---

## Phase 3 — Blindspot Hunt

### [blindspot] 1: Data flow gap — intent is not reachable from BuildIndex()

The proposal says `BuildIndex()` should read intent from `proposal.md` frontmatter. Currently, `BuildIndex()` takes `BuildIndexOpts` which has `FeatureSlug` but no `Intent` field. The proposal mentions adding intent to `BuildIndexOpts`, but the actual data source is `proposal.md` which lives at `docs/proposals/{slug}/proposal.md`. The `proposal.FindBySlug()` function already parses intent from this file. However, `BuildIndex()` is called from a CLI command handler, and the proposal does not trace the call chain from CLI handler -> `proposal.FindBySlug()` -> extract intent -> pass to `BuildIndexOpts.Intent`. This is a concrete implementation gap that could cause the feature to fail if the caller doesn't provide intent.

### [blindspot] 2: autogen.go dependency wiring assumes business tasks exist

The proposal's refactor wiring says "最后一个 business task 的 taskID 作为 validate-code 的 depends_on." But what if there are zero business tasks (e.g., a cleanup feature where all tasks are doc-type)? The `autogen.go` functions currently look up `lastRunID` as an anchor — removing this anchor without handling the zero-business-task case could cause a nil dependency or a dangling `depends_on` reference.

### [blindspot] 3: No rollback plan

The proposal has no rollback mechanism. If the intent-driven branching causes a regression in existing new-feature pipelines (e.g., the autogen.go changes break the existing dependency wiring), there's no described path back. The proposal mentions "default to new-feature" as backward compatibility, but this doesn't help if the autogen.go refactoring itself introduces a bug.

### [blindspot] 4: Stage-gate generation coupling

The proposal says `IsTestableType()` will return false for refactor/cleanup, and stage-gate generation (step 6.5 in build.go) is gated by `needsTest`. But the proposal also says "stage-gate 任务（validate-code, clean-code）的生成不受 IsTestableType() 影响" — this means stage-gates should still be generated for refactor. The current code generates stage-gates only when `needsTest` is true. The proposal needs to decouple stage-gate generation from `needsTest`, but this change is not explicitly called out in the scope.

### [blindspot] 5: proposal.md may not exist at index time

`forge task index` can be run before brainstorm completes. If the proposal.md file doesn't exist yet (user runs `forge task index` on a manually created feature directory), `readFrontmatter()` will fail. The proposal doesn't address this edge case. The `detectMode()` function handles this by returning empty string — intent reading needs similar resilience.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 95 | 110 |
| Solution Clarity | 98 | 120 |
| Industry Benchmarking | 82 | 120 |
| Requirements Completeness | 88 | 110 |
| Solution Creativity | 75 | 100 |
| Feasibility | 78 | 100 |
| Scope Definition | 62 | 80 |
| Risk Assessment | 70 | 90 |
| Success Criteria | 58 | 80 |
| Logical Consistency | 74 | 90 |
| **Total** | **748** | **1000** |

**Target**: 900
**Gap**: 152 points
