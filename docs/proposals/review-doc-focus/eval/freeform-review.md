# Freeform Expert Review: Review-Doc Execution Scope Focus

**Reviewer**: Prompt Compliance Architect
**Date**: 2026-05-24
**Document**: `docs/proposals/review-doc-focus/proposal.md`

---

## Section 1: Background Assessment

This proposal addresses a behavioral drift problem in the `T-review-doc` agent execution flow. The core observation is sound: the current prompt templates (`prompt/data/doc-review.md` and `task/data/doc-review.md`) instruct the agent to perform two operations that are architecturally unnecessary and actively harmful -- scanning the entire tasks directory for `.md` files, and reading acceptance criteria from those files at runtime.

The proposal's two-pronged approach is architecturally correct in principle:

1. **Build-time AC extraction**: Move AC from per-task `.md` files into the review-doc task definition at `forge task index` time, eliminating the need for runtime file scanning.
2. **Filter-based discovery strategy**: Constrain the agent's file discovery to content documents only, excluding `tasks/`, `records/`, `manifest.md`, and `index.json`.

The proposal correctly identifies the problem as an information flow design issue rather than a prompt engineering issue, explicitly noting that "LLM frequently ignores prompt constraints" as the rationale for rejecting a pure-prompt approach. This is the right instinct -- structural constraints that remove the capability to deviate are always more reliable than behavioral instructions that merely discourage deviation.

However, the proposal operates at a high level of abstraction. After reading the actual source code, I find several significant gaps between the proposal's description of what needs to change and what the implementation actually requires.

---

## Section 2: Key Risk Identification

风险：The proposal's AC extraction design relies on a single shared `BodyContext.AcceptanceCriteria` field that currently serves the PRD-level AC (breakdown mode). The proposal says "forge task index generates the review-doc task, extracts Acceptance Criteria from all doc tasks' `.md` files, and aggregates them into a dedicated area in `review-doc.md`." But the existing `extractBodyContext()` function in `extract.go` only extracts AC from the PRD (in breakdown mode) via `extractAcceptanceCriteria(content)` reading `prd-spec.md`. The proposal requires an entirely new extraction path: iterating over all doc-task `.md` files in the tasks directory, parsing each one's `## Acceptance Criteria` section, and aggregating them. This is a fundamentally different operation from the existing single-file AC extraction. The proposal treats it as a minor extension of `BuildIndex()` but it is actually a new data pipeline with its own parsing, aggregation, and error handling requirements. The current `BodyContext` struct has no field for per-task AC mapping; the `AcceptanceCriteria []string` field is a flat list of PRD-level criteria, not a structured per-task aggregation. The proposal does not address this schema gap.

> Direct quote: "修改 `build.go`：`BuildIndex()` 从 doc 任务文件提取 AC 并传入 review-doc 生成"

问题：The proposal states the change is to `build.go` `BuildIndex()`, but the actual extraction will need to happen after the task scanning loop (Step 5 in `BuildIndex`) has already identified all doc-type tasks, and the AC data must flow through `BodyContext` into `renderBody()` which then substitutes `{{ACCEPTANCE_CRITERIA}}`. The current `renderBody` handles `{{ACCEPTANCE_CRITERIA}}` as a flat checkbox list. To support per-task AC aggregation, either: (a) the `BodyContext` needs a new field like `DocTaskCriteria map[string][]string` with a new template placeholder, or (b) the existing flat list must be repurposed with a naming convention like "Task X: criterion Y." Neither approach is described.

> Direct quote: "修改 `autogen.go`：review-doc 任务定义支持嵌入 AC 数据"

风险：The `GetReviewDocTask()` function currently takes no parameters and returns a static `AutoGenTaskDef`. The proposal says to modify it to "accept AC data," but `AutoGenTaskDef` has no AC-related field. AC data flows through `BodyContext` into `renderBody()` and then into the template via `{{ACCEPTANCE_CRITERIA}}`. This means the change is not to `GetReviewDocTask()` at all -- it is to the `BodyContext` construction in `BuildIndex()` and to the template file `task/data/doc-review.md`. The proposal misidentifies the modification surface. An implementer following these instructions literally would modify the wrong function.

> Direct quote: "agent prompt 的文档发现从'扫描全部'改为'只扫描内容文档'，明确排除 `tasks/`、`tasks/records/`、`manifest.md`、`index.json`"

问题：Looking at the actual autogen template `task/data/doc-review.md`, the Discovery Strategy currently says:

```
Scan these directories for ALL documents created or modified by this feature:
- docs/features/{{FEATURE_SLUG}}/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/{{FEATURE_SLUG}}/
```

This Discovery Strategy lives in the autogen task template, which is the task `.md` file content. But the agent's actual file-reading behavior is governed by `prompt/data/doc-review.md`, which has its own Step 1 and Step 2 instructions. The proposal's filtering must be applied consistently to both templates, but it only mentions the autogen template (`task/data/doc-review.md`) for the Discovery Strategy change and the prompt template (`prompt/data/doc-review.md`) for removing task scanning. These are two separate behavioral anchors that must be kept in sync, and the proposal does not acknowledge this coupling.

> Direct quote: "过滤基于明确路径模式（tasks/、records/、manifest.md、index.json），不影响 docs/ 下的内容文件"

风险：The proposal claims the filter "does not affect content files under docs/," but this assumes all content documents live under `docs/features/{{FEATURE_SLUG}}/` or `docs/proposals/{{FEATURE_SLUG}}/`. In the Forge pipeline, task `.md` files themselves live in the `tasks/` directory within the feature folder. But what about edge cases where a doc task's deliverable is itself a task-related file? For example, if someone creates a doc task whose deliverable is in `docs/features/SLUG/tasks/some-guide.md` (a path that contains "tasks" but is actually a content document), the naive string-based exclusion could filter it out. The proposal does not specify whether the filtering uses path prefix matching, path component matching, or exact directory matching. This ambiguity creates a correctness risk.

> Direct quote: "若 section 不存在则标记为空"

问题：The proposal's risk mitigation for "AC extraction parsing instability" says "若 section 不存在则标记为空." But what happens at runtime when the agent receives an empty AC for a doc task? The proposal's Step 2 workflow says "核对 AC -> 修复", but if AC is empty, the agent has no criteria to check. This creates a silent degradation path: doc tasks without AC sections will pass review without any substantive quality check. The proposal should define what "标记为空" means behaviorally -- does the agent skip that task entirely, or does it attempt a freeform review? The current prompt has no handling for this case.

> Direct quote: "review-doc 依赖最后业务任务完成后才执行，手动修改发生在执行前，构建时状态即为最新"

风险：This assumes that `forge task index` is only called once, after all business tasks are complete. But `BuildIndex()` is designed to be idempotent and can be called at any time. If a user runs `forge task index` before completing all business tasks, the AC embedded in `review-doc.md` will be stale or partial. The dependency chain ensures `T-review-doc` waits for business tasks to complete, but the AC data was baked in at `BuildIndex()` time, not at execution time. This is a temporal coupling that the proposal does not acknowledge. The AC could be stale if tasks are added or modified between index time and review execution time.

> Direct quote: "修改 `prompt/data/doc-review.md`：简化 agent prompt——移除任务扫描指令，改用内嵌 AC，加入文档过滤规则，禁止修改任务文件和记录"

问题：This is the prompt engineering heart of the proposal, and it is the thinnest section. The proposal says "改用内嵌 AC" and "加入文档过滤规则" and "禁止修改任务文件" as if these are trivial prompt changes. But the current `prompt/data/doc-review.md` is a 96-line prompt template with a carefully structured 5-step workflow. Removing Step 1's "scanning tasks directory" instruction and Step 2's "Read the task's acceptance criteria from its .md file" instruction requires restructuring the entire Step 1-2 flow. The agent needs to know: (1) where to find the embedded AC (in the task file itself? in a new template section?), (2) how to map AC items to deliverable documents without scanning task files, and (3) what the new file discovery mechanism looks like. The proposal provides no prompt template mockup, no step-by-step restructuring plan, and no example of the modified prompt. This is the highest-risk modification -- prompt template changes directly control agent behavior -- and it receives the least specification.

> Direct quote: "纯 prompt 约束被否决。Reason: LLM 经常忽略 prompt 约束，问题根源需从信息流层面解决"

风险：The proposal rejects pure prompt constraints as unreliable, then proposes... a prompt modification as half of its solution. The "filter-based discovery strategy" is fundamentally a prompt instruction that says "only scan these directories, exclude these paths." If the agent can ignore "scan tasks directory," it can also ignore "do not scan tasks directory." The build-time AC embedding is the genuinely structural change -- it removes the *need* to scan. But the filtering is still a behavioral instruction. The proposal's own logic undermines the reliability claim for the filtering half of the solution.

> Direct quote: "AC 覆盖率不降低（所有 doc 任务的 AC 均被检查）"

问题：This success criterion is stated as a boolean, but there is no verification mechanism proposed. How will the implementation verify that every doc task's AC was actually extracted and embedded? The current `extractCheckboxItems` function in `extract.go` is a line-by-line parser that matches `## Acceptance Criteria` headings. If a doc task formats its AC differently (e.g., `## 验收标准` in Chinese, or `## Acceptance Criteria & Quality Gates` with extra text in the heading), the extraction will silently fail and return nil. The proposal's risk table acknowledges "doc 任务的 AC 格式不完全统一" but rates the likelihood as Medium and the mitigation as "tolerate content format differences." Tolerance of format differences in heading matching is not the same as tolerance of content differences -- the heading must match exactly or the entire section is missed.

> Direct quote: "预估 2-3 个 coding 任务"

风险：The proposal underestimates the implementation surface. The changes touch: (1) a new AC aggregation pipeline in `BuildIndex()`, (2) a new `BodyContext` field or template placeholder for per-task AC, (3) the `task/data/doc-review.md` template rewrite, (4) the `prompt/data/doc-review.md` template rewrite with restructured workflow steps, (5) the `renderBody()` function for the new placeholder, (6) tests for the extraction pipeline, and (7) tests for the modified prompt template rendering. This is not 2-3 tasks; it is at minimum 4-5 tasks, and the prompt template changes alone carry enough behavioral risk to warrant their own dedicated validation task.

> Direct quote: "autogen 模板增加 AC 汇总区域，更新 Discovery Strategy 排除 tasks/ 和 records/"

问题：The proposal says "增加 AC 汇总区域" in the autogen template, but does not define the format of this aggregation. Should it be a `## Aggregated Acceptance Criteria` section with sub-headings per task? A flat list with task-name prefixes? A table? The format matters because the agent prompt (`prompt/data/doc-review.md`) must instruct the agent how to parse and use this section. Without a defined format, the prompt template cannot be written correctly, and the implementation tasks cannot be specified. This is a specification gap that blocks implementation.

---

## Section 3: Improvement Suggestions

建议：Define a concrete AC aggregation schema before implementation. Specify the exact markdown format for the embedded AC section in `review-doc.md`. For example:

```markdown
## Aggregated Acceptance Criteria

### Task: write-prd (T-write-prd)
- [ ] PRD contains functional requirements section
- [ ] PRD defines success metrics

### Task: write-design (T-write-design)
- [ ] Design document covers architecture decisions
- [ ] Design document includes data model
```

This gives the prompt template a stable parse target and makes the mapping between AC items and deliverable documents explicit. Without this, the prompt template author (or the implementation task executor) will have to invent a format, creating an implicit coupling between the Go extraction code and the prompt template that is invisible until runtime.

This addresses the risk of "AC 汇总区域" being underspecified, and the problem of the prompt template needing to know the data format it consumes.

建议：Add a build-time validation step after AC extraction that logs a warning when any doc task has an empty or missing AC section. The existing `BuildIndex()` already accumulates `result.Warnings` for various issues. Add a warning like:

```
WARNING: doc task "T-write-prd" has no "## Acceptance Criteria" section — review-doc will have no criteria for this task
```

This makes the "silent degradation" path visible. Without it, a missing AC section is invisible until review time, where the agent simply skips that task's quality check. This directly addresses the risk of empty AC being "marked as empty" with no behavioral consequence.

建议：Acknowledge the temporal coupling between `BuildIndex()` and `T-review-doc` execution explicitly in the proposal's Constraints section. Add a note that `BuildIndex()` must be re-run if doc tasks are added or modified after the initial index build. Better yet, modify the proposal to extract AC at execution time rather than at index time -- the task-executor could pass AC data as a parameter when dispatching `T-review-doc`, ensuring freshness. If execution-time extraction is too architecturally invasive, at minimum document that `forge task index` should be the last command before task execution begins.

This addresses the risk of stale AC data caused by the build-time vs. execution-time gap.

建议：Separate the prompt template changes into their own implementation task with explicit behavioral acceptance criteria. The prompt template modification is the highest-risk change because it directly controls agent behavior, and the proposal currently lumps it into "模板文件修改" alongside the autogen template. The prompt task should include: (1) a mockup of the modified `prompt/data/doc-review.md` with the restructured workflow, (2) a clear specification of what replaces Step 1's task scanning and Step 2's AC reading, and (3) a concrete example of the embedded AC section format from the agent's perspective. The success criteria for this task should include a test execution showing that the agent no longer reads `tasks/` or `records/`.

This addresses the problem of the prompt modification being the most critical change yet receiving the least specification.

建议：Replace the "过滤式发现策略" with an explicit allowlist approach in the prompt template, rather than a denylist. Instead of "exclude tasks/, records/, manifest.md, index.json," specify:

```
## Target Documents
You may ONLY read and modify files in these directories:
- docs/features/{{FEATURE_SLUG}}/prd/
- docs/features/{{FEATURE_SLUG}}/design/
- docs/features/{{FEATURE_SLUG}}/testing/
- docs/proposals/{{FEATURE_SLUG}}/
```

An allowlist is structurally stronger than a denylist for LLM behavioral control. A denylist says "don't go here" which the agent can ignore or work around; an allowlist says "you can only go here" which creates a tighter mental constraint. This directly addresses the risk that the filtering half of the solution is still a behavioral instruction and therefore inherently unreliable.

建议：Add a degradation path specification to the proposal. What should the agent do when the embedded AC section is empty or malformed? Options include: (a) fail with a clear error message rather than proceeding with no criteria, (b) fall back to scanning task files as a degraded mode, or (c) skip the affected task and report it as "unable to review." The proposal should choose one and document it. Without this, the agent's behavior on extraction failure is undefined and will be determined by whatever the implementer decides at coding time.

This addresses the risk of undefined behavior when the build-time extraction fails, which the proposal's own risk table acknowledges as a Medium likelihood event.
