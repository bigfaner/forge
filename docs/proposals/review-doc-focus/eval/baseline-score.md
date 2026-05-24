# Baseline Score: Review-Doc Execution Scope Focus

**Evaluator**: CTO Adversary (CTO Rubric)
**Date**: 2026-05-24
**Document**: `docs/proposals/review-doc-focus/proposal.md`
**Informational Baseline** — score before any revision

---

## Total Score: 560/1000

---

## Dimension Scores

### 1. Problem Definition: 85/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 32/40 | Core problem is identifiable: agent reads irrelevant files during review-doc execution, causing token waste, quality degradation, and accidental modifications. The three sub-problems are enumerated but the first two ("token waste" and "quality degradation") overlap conceptually -- quality degradation is partly a consequence of token waste, not an independent problem. Still, the statement is specific enough to act on. |
| Evidence provided | 30/40 | Cites concrete file paths (`prompt/data/doc-review.md` Step 1, Step 2, `task/data/doc-review.md` Discovery Strategy) and references "actual execution records." However, no quantitative evidence: how much token waste? How many execution records were examined? What percentage of reviewed files were irrelevant? The evidence is anecdotal, not measured. Quote: "实际执行记录显示 Referenced Documents 包含任务文件和记录文件" -- which records? How many? This is a hand-wave, not data. |
| Urgency justified | 23/30 | States the pipeline is merged and in active use, so every execution suffers. Valid but thin: no estimate of frequency (daily? weekly?) or magnitude of impact (100 tokens wasted? 10,000?). The urgency claim rests on "持续浪费 token 且审校质量受影响" without quantification. |

**Deductions**:
- Vague language without quantification: "实际执行记录显示" (-20 pts) -- no specific records cited, no numbers.

### 2. Solution Clarity: 65/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 25/40 | Two-pronged approach named (build-time AC embedding + filter-based discovery), but the "AC embedding" mechanism is described at hand-wave level. Quote: "从所有 doc 任务的 `.md` 文件中提取 Acceptance Criteria，汇总写入 `review-doc.md` 的专门区域" -- what is the "专门区域"? What is the format? How are per-task ACs distinguished? An implementer cannot explain back the exact data flow. |
| User-facing behavior described | 20/45 | The end user (agent operator) will see: no change in invocation, but the agent will no longer read task files or records. This is implied, never explicitly stated. There is no "before/after" comparison of what the agent reads. The user-facing behavior section is entirely absent -- the proposal describes internal mechanics, not the user experience. |
| Technical direction clear | 20/35 | Names three files to modify (`autogen.go`, `build.go`, templates). However, the freeform review already identified that the proposal misidentifies the modification surface: `GetReviewDocTask()` does not need modification -- the change is to `BodyContext` construction and `renderBody()`. Quote: "修改 `autogen.go`：review-doc 任务定义支持嵌入 AC 数据" -- but `AutoGenTaskDef` has no AC field. The technical direction points at the wrong functions. |

**Deductions**:
- Vague language: "汇总写入 `review-doc.md` 的专门区域" (-20 pts) -- no format specified for the "dedicated area."

### 3. Industry Benchmarking: 35/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 10/40 | No real-world solutions or patterns cited. No reference to how other agent frameworks (LangGraph, CrewAI, AutoGen) handle context scoping or task isolation. No reference to RAG filtering patterns, sandboxing approaches, or agent constraint mechanisms in the industry. |
| At least 3 meaningful alternatives | 15/30 | Three alternatives listed, but two are weak. "Do nothing" is present (good). "Pure prompt constraints" is present but treated as a straw-man: the sole argument against it is "LLM frequently ignores constraints" which is asserted, not evidenced. The proposal itself then relies on prompt-based filtering as half its solution, undermining its own rejection argument. The third alternative is the selected approach. There are no genuinely different alternatives explored (e.g., execution-time extraction, separate AC manifest file, agent sandboxing via filesystem permissions). |
| Honest trade-off comparison | 5/25 | The trade-off table has one row per alternative with vague pros/cons. Quote for selected approach: "从根源消除问题，信息流清晰" vs "build 阶段需解析 AC，复杂度略增" -- "复杂度略增" understates the actual scope (new extraction pipeline, new BodyContext field, template rewrite, temporal coupling). The selected approach's con is minimized. |
| Chosen approach justified against benchmarks | 5/25 | No benchmarks to justify against. The rejection of "pure prompt" is self-contradictory given the proposal's own filtering mechanism is prompt-based. |

**Deductions**:
- Straw-man alternative: "纯 prompt 约束" rejected on a single assertion without evidence, while the proposal itself uses prompt constraints (-20 pts).

### 4. Requirements Completeness: 65/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 25/40 | Three scenarios listed (pure-doc, mixed, no-doc). Missing scenarios: (1) doc task with no `## Acceptance Criteria` section, (2) doc task with malformed AC section, (3) feature with only one doc task (trivial case), (4) concurrent modification of tasks between index and execution. The "no doc tasks" scenario correctly notes existing logic is unchanged, but the edge cases around empty/missing AC are unaddressed. |
| Non-functional requirements | 25/40 | Lists token reduction, AC coverage preservation, and backward compatibility. Missing: no performance requirement for AC extraction (how fast must parsing be for N tasks?), no security consideration (can malicious AC content in task files inject agent instructions?), no accessibility requirement. The NFRs are stated as aspirations ("降低", "不降低") without measurable thresholds. Quote: "review-doc 执行 token 消耗降低" -- lower by how much? 10%? 50%? |
| Constraints & dependencies | 15/30 | Lists three dependencies. Missing: temporal coupling constraint (BuildIndex vs execution time), the constraint that doc tasks must have `## Acceptance Criteria` as exact section header, the constraint that the agent framework (Claude) respects allowlist instructions. |

**Deductions**:
- Vague language: "token 消耗降低" without quantification (-20 pts).

### 5. Solution Creativity: 35/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 10/40 | The proposal explicitly states "无特殊创新" (no special innovation) and frames the approach as basic separation of concerns. This is honest but earns low creativity points. Build-time data extraction for runtime consumption is a well-known pattern (e.g., compiled configurations, ahead-of-time optimization). |
| Cross-domain inspiration | 10/35 | No cross-domain borrowing demonstrated. The proposal does not reference any patterns from other fields (compiler design's dead code elimination, database query optimization's predicate pushdown, information retrieval's relevance filtering, etc.). |
| Simplicity of insight | 15/25 | The core insight -- "don't give the agent files it shouldn't read" -- is indeed simple and elegant. The "Assumptions Challenged" section is the strongest part: correctly identifying that runtime discovery was unnecessary when all data is available at build time. Quote: "doc 任务在 `forge task index` 时已确定，AC 可在构建时静态提取，无需运行时扫描" -- this is a clean insight. |

### 6. Feasibility: 55/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 20/40 | The proposal claims "完全可行" (completely feasible) but the freeform review identified multiple implementation gaps: `GetReviewDocTask()` does not need the changes described, the `BodyContext` schema lacks a field for per-task AC, the `renderBody()` function needs a new placeholder, and the prompt template restructuring is unspecified. The proposal misidentifies the modification surface and underestimates complexity. |
| Resource & timeline feasibility | 20/30 | Estimates "2-3 coding tasks." The freeform review correctly identifies this as an underestimate: the actual surface is at minimum 4-5 tasks (new AC extraction pipeline, BodyContext schema + renderBody modification, autogen template rewrite, prompt template rewrite, testing). The underestimate is not catastrophic but indicates incomplete scoping. |
| Dependency readiness | 15/30 | All dependencies are internal (Go CLI, templates). No external dependencies. However, the proposal does not verify that the `extractCheckboxItems` function in `extract.go` can handle per-task AC extraction (it currently only extracts from PRD). The dependency on a specific section header format (`## Acceptance Criteria`) is stated but not validated against actual task files. |

### 7. Scope Definition: 55/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 22/30 | Four in-scope items, each naming a specific file. However, the items are at the "modify this file" level, not the "deliver this specific change" level. Quote: "修改 `autogen.go`：review-doc 任务定义支持嵌入 AC 数据" -- but as noted, `GetReviewDocTask()` doesn't need modification. The deliverable is misdescribed. |
| Out-of-scope explicitly listed | 20/25 | Four out-of-scope items named. Good: clearly excludes other task types, AC format changes, eval pipeline, record template changes. One notable omission: no mention of whether the `extractBodyContext()` function changes are in or out of scope. |
| Scope is bounded | 13/25 | The "2-3 coding tasks" estimate bounds the work, but as noted, this is likely an underestimate. The scope is bounded in intent but not in accurate sizing. No explicit timeframe given (days? sprint?). |

### 8. Risk Assessment: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 18/30 | Three risks identified. Missing risks: (1) stale AC due to temporal coupling between BuildIndex and execution, (2) prompt template changes causing behavioral regression (agent breaks in new ways), (3) schema gap between `BodyContext.AcceptanceCriteria` (flat PRD list) and per-task AC aggregation, (4) misidentification of modification surface leading to wasted implementation effort. |
| Likelihood + impact rated | 18/30 | Ratings are present (M/M, L/L, L/M) but the first risk's Medium likelihood for "AC format not unified" seems honest. However, the temporal coupling risk (BuildIndex time vs execution time) is rated as L likelihood, but `BuildIndex()` is designed to be called at any time -- this is actually a High likelihood scenario that the proposal dismisses. |
| Mitigations are actionable | 19/30 | First mitigation (section header matching with empty fallback) is actionable but has the undefined-behavior problem identified in the freeform review. Second mitigation (temporal argument) is an assertion, not an action. Third mitigation (path-based filtering) is actionable but ambiguous (prefix vs component matching). |

### 9. Success Criteria: 45/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 30/55 | Five checkbox criteria. Three are testable: "不再读取 tasks/ 目录" (can verify via execution logs), "不再读取 records/" (same), "只读取和修改目标交付文档" (same). Two are softer: "AC 覆盖率不降低" has no verification mechanism (how do you prove all ACs were checked?), "autogen 模板包含 AC 汇总区域" is trivially testable (check file contents) but doesn't verify the AC data is correct. |
| Coverage is complete | 15/25 | Missing success criteria for: (1) prompt template changes (only autogen template mentioned), (2) backward compatibility (listed as NFR but no success criterion), (3) handling of tasks with no AC section, (4) build-time validation of AC extraction completeness. |

**Deductions**:
- Vague language: "AC 覆盖率不降低" without definition of how coverage is measured (-20 pts).

### 10. Logical Consistency: 65/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 25/35 | Build-time AC embedding directly addresses Problem 1 (token waste) and Problem 2 (attention distraction) by removing the need to read task files. Problem 3 (agent modifying wrong files) is addressed by the filtering strategy. However, the filtering is a prompt instruction, which the proposal itself acknowledges is unreliable. The solution partially relies on the mechanism it explicitly rejected as unreliable. |
| Scope <-> Solution <-> Success Criteria aligned | 20/30 | Scope items map to solution components. But success criteria do not cover all scope items: the prompt template modification (`prompt/data/doc-review.md`) is in scope but has no dedicated success criterion. The autogen template has a success criterion but the prompt template does not. |
| Requirements <-> Solution coherent | 20/25 | Requirements (token reduction, AC coverage, backward compat) are coherent with the solution approach. The constraint that "AC 提取需能解析 doc 任务 `.md` 文件中的 `## Acceptance Criteria` section" is clear and maps directly to the extraction logic. The weakest link is the temporal dependency constraint which is listed as "不改变 task-executor 的调度机制" but the proposal actually introduces a new temporal coupling. |

---

## Summary of Deductions Applied

| Deduction Type | Instances | Points |
|---------------|-----------|--------|
| Vague language without quantification | 4 instances ("实际执行记录显示", "汇总写入专门区域", "token 消耗降低", "AC 覆盖率不降低") | -80 pts |
| Straw-man alternative | 1 instance (pure prompt rejected while proposal uses prompt constraints) | -20 pts |

---

## Attacks (Top Weaknesses)

1. **Industry Benchmarking**: Zero industry references. No citation of how other agent frameworks handle context scoping. The "Alternatives" table contains the selected approach, a straw-man, and do-nothing -- not three genuinely different alternatives. No justification beyond internal reasoning.

2. **Solution Clarity**: The proposal misidentifies the modification surface. Quote: "修改 `autogen.go`：review-doc 任务定义支持嵌入 AC 数据" -- but `GetReviewDocTask()` returns a static `AutoGenTaskDef` with no AC-related field. The actual change needed is to `BodyContext` and `renderBody()`. An implementer following these instructions would modify the wrong code.

3. **Solution Clarity**: No format specified for the "AC aggregation area." Quote: "汇总写入 `review-doc.md` 的专门区域" -- what section heading? What markdown structure? How are per-task ACs delineated? This is a specification gap that blocks implementation.

4. **Requirements Completeness**: No edge case for tasks with missing `## Acceptance Criteria`. The risk table says "若 section 不存在则标记为空" but no requirement or success criterion defines what happens behaviorally when AC is empty. The agent could silently skip quality checks.

5. **Feasibility**: The proposal relies on prompt-based filtering as half its solution while explicitly rejecting prompt constraints as unreliable. Quote: "LLM 经常忽略 prompt 约束，问题根源需从信息流层面解决" -- then proposes filtering instructions in the prompt template as a core mechanism. This is internally contradictory.

6. **Risk Assessment**: Temporal coupling between `BuildIndex()` (which embeds AC) and `T-review-doc` execution is dismissed as L likelihood. But `BuildIndex()` is designed to be called at any time, and the proposal introduces no guard against stale embedded AC. The likelihood should be M or H.

7. **Success Criteria**: "AC 覆盖率不降低" is unmeasurable as stated. No verification mechanism is proposed. How does one prove that every doc task's AC was extracted and will be checked?

8. **Feasibility**: "2-3 coding tasks" underestimates the implementation surface. The prompt template rewrite alone (restructuring a 96-line, 5-step workflow) warrants its own task with behavioral validation. The actual count is 4-5 minimum.

9. **Scope Definition**: The in-scope item for `autogen.go` describes a change to `GetReviewDocTask()` that is not actually needed. This wastes implementation effort and indicates the proposal was written without verifying the code against the described changes.

10. **Requirements Completeness**: Non-functional requirement "token 消耗降低" has no threshold. Without a target (e.g., "reduce by 50%"), this NFR cannot be verified as met or not met.

---

## Blindspot Hunt (Beyond the Rubric)

1. **Rollback plan entirely absent**: The proposal has no rollback or revert strategy. If the build-time AC extraction breaks or the filtered prompt causes the agent to miss critical quality checks, there is no documented path back to the current behavior. For a production pipeline ("已合并并投入使用"), this is a serious omission.

2. **No migration consideration**: Existing `review-doc.md` files already generated with the old format will coexist with the new format. The proposal does not address whether `BuildIndex()` handles re-generation of existing `review-doc.md` files or whether manual cleanup is needed.

3. **Hidden coupling between two templates**: The proposal treats `task/data/doc-review.md` and `prompt/data/doc-review.md` as independent modifications, but they must stay synchronized. The autogen template defines the data format the agent consumes; the prompt template defines how the agent processes it. No mechanism ensures consistency.

4. **No acceptance test for behavioral correctness**: The success criteria are structural (agent doesn't read certain files) but do not verify that the review quality is maintained. The proposal's Problem 2 is "审校质量下降" but no success criterion measures review quality before and after the change.

5. **Missing cost-benefit analysis**: The proposal claims token savings but does not estimate the development cost of the new AC extraction pipeline vs. the recurring cost of wasted tokens. For a change described as "no innovation, just separation of concerns," the cost justification should be straightforward but is absent.
