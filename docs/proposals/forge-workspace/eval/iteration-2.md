---
iteration: 2
evaluator: CTO Persona (Rubric-Based)
date: 2026-06-07
previous_iteration: 1
document_sha: pending
---

# Iteration 2: Rubric Evaluation Report

## Overview

Evaluation of `docs/proposals/forge-workspace/proposal.md` against the 1000-point Proposal Evaluation Rubric. This is iteration 2. The proposal has been revised to address the 5 attack points from iteration 1.

## Iteration 1 Attack Resolution Check

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | **[D3] No real industry solutions cited** | **Resolved** | Alternatives table now references npm/Yarn workspaces (line 156), Cargo workspaces (line 157), Turborepo/Nx (line 158), VS Code multi-root workspaces (line 159). Each includes Source, Pros, Cons, and differential Verdict. Innovation Highlights section (lines 109-117) adds a formal comparison table. |
| 2 | **[D9] SC coverage gaps — close, register, cache, field inheritance** | **Resolved** | SC#7 (line 277) covers `register`. SC#8 (line 283) covers `close` with preconditions. SC#11 (line 286) covers cache behavior. SC#7 (line 282) now lists specific inherited fields (title, intent, status=Draft, source). |
| 3 | **[D5] No cross-domain inspiration** | **Resolved** | Innovation Highlights section (lines 109-117) explicitly maps cross-domain borrowing with a 4-row table (Cargo, VS Code, Turborepo) and summarizes core differentiation. |
| 4 | **[D6] No timeline / work breakdown** | **Resolved** | Resource & Timeline section (lines 175-185) now has a 6-row work breakdown table with person-day estimates per scope item, totaling 8-13d. |
| 5 | **[D10] Scope-SC alignment gaps** | **Resolved** | All scope items now have corresponding SC entries. Field mapping table is cross-referenced by SC#7 (line 282) which lists specific fields. |

**Unresolved from iteration-1 "Should Fix"**:
- Git-less workspace risk: Not acknowledged. No mitigation proposed.
- Failure scenarios for corrupted `.forge-workspace.yaml` and overlapping workspaces: Not addressed.

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem -> Solution**: Sound. Four pain points; 1-3 scoped to process documents, 4 excluded. Workspace as process document management layer maps directly.

**Solution -> Evidence**: Improved. The Innovation Highlights section now provides comparative evidence through the cross-domain table (lines 109-117). However, evidence remains theoretical — no prototype, no spike. Acceptable for proposal stage.

**Evidence -> Success Criteria**: Bidirectional mapping is now complete. Every in-scope item has at least one SC entry. Field inheritance is explicitly testable (SC#7 line 282). Cache behavior has a concrete SC (line 286).

**Self-contradiction**: No new contradictions found. The "项目不变" principle is correctly bounded by the explicit brainstorm adaptation scope (line 147).

### SC Consistency Deep-Dive

**Cluster A: Project Discovery & Registration**
- In Scope #1 (init, register, .forge-workspace.yaml)
- SC#1 (init discovers), SC#2 (register), SC#9 (unhealthy marking), SC#10 (zero-discovery diagnostics)
- **SC <-> SC**: No contradiction.
- **SC <-> InScope**: Bidirectional. init and register both testable.
- **Gap**: None. Resolved from iteration 1.

**Cluster B: Status & Feature Aggregation**
- In Scope #2 (status), #3 (features)
- SC#3 (status table + perf), SC#4 (status <project>), SC#5 (features), SC#11 (cache), SC#12 (v0 schema)
- **SC <-> SC**: No contradiction.
- **SC <-> InScope**: Complete coverage.
- **Gap**: None.

**Cluster C: Workspace-level Proposals**
- In Scope #4 (propose, assign, close)
- SC#6 (propose), SC#7 (assign with field list), SC#8 (close with preconditions)
- **SC <-> SC**: No contradiction. SC#7 requires proposal in Approved state; SC#8 requires feature in Completed/Closed — consistent state machine.
- **SC <-> InScope**: Bidirectional. Field mapping table in Scope (lines 235-248) is testable via SC#7 which lists specific fields. State transitions in Scope (line 246) are verifiable through SC#7 (Approved prerequisite) and SC#8 (Completed/Closed prerequisite).
- **Gap**: None. Resolved from iteration 1.

**Cluster D: Module Boundaries**
- In Scope #5, #6 (internal API, isolation)
- No SC entries.
- **Verdict**: Acceptable — explicitly marked "v1 内部 API 契约，非公开接口".

---

## Phase 2: Rubric Scoring

### D1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 36/40 | Four concrete pain points. The scoping note ("1-3 是过程文档管理问题，4 是知识管理问题") is unambiguous. Minor deduction: pain point #3 ("多项目切换摩擦 - 丢失上下文，需要反复重新加载状态") remains vaguely described. What specific context is lost? What does "重新加载状态" entail? A concrete example ("I forget which phase project X is in") would strengthen this. |
| Evidence provided | 30/40 | Four evidence points. All self-reported but concrete (4-8 projects, existing proposals in separate docs/, two existing Draft proposals). "已有两个 Draft 提案" is the strongest evidence — it shows the need has already been recognized. Deduction: no frequency data (how many context switches per day?), no time cost measurement. |
| Urgency justified | 25/30 | "每多一个项目，过程文档管理摩擦非线性增长" is intuitive but still unquantified. However, "4-8 个项目的规模已影响日常效率" grounds it in current state. The "推迟意味着持续在项目间来回切换" is a concrete cost of delay. Deduction: "非线性增长" remains an unsupported claim. |

**D1 Total: 91/110**

### D2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 39/40 | CLI commands, file paths, directory structure, field mapping table, state machine — all specified. The architecture diagram is clear. The only minor gap: the relationship between `.forge-workspace.yaml` (registry) and `.forge-workspace/config.yaml` (module config) is mentioned but the lifecycle of the latter (when is it created?) is implicit. |
| User-facing behavior described | 42/45 | Six key scenarios with commands and expected outputs. `status` table columns specified. `close` preconditions explicit. Deduction: `features` command output format is still unspecified — what columns? What does "阶段" look like? Is it the feature manifest phase? |
| Technical direction clear | 34/35 | File system scanning, manifest reading, mtime caching with specific parameters (5-minute threshold, per-project mtime fingerprint). The cache strategy is well-specified. Minor gap: workspace directory detection logic (how does a command know it's in a workspace vs a project?) is implied but not stated. |

**D2 Total: 115/120**

### D3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 34/40 | npm/Yarn workspaces, Cargo workspaces, Turborepo/Nx, VS Code multi-root workspaces — four specific, named industry solutions with sources. Each has Pros/Cons derived from actual characteristics. Strong improvement from iteration 1. Minor deduction: no reference to CLI-based project tracking tools (Taskwarrior, todo.txt ecosystem) or knowledge/document management patterns (Obsidian vaults, Notion databases). The benchmarking is monorepo-heavy; the proposal's problem domain (process document lifecycle) has parallels in other areas. |
| At least 3 meaningful alternatives | 26/30 | Eight alternatives including "do nothing". Four are industry-validated (npm/Yarn, Cargo, Turborepo/Nx, VS Code). Two are internal (forge-dashboard, forge-wiki). All are genuine alternatives, not straw men — each has honest Pros listed. Minor deduction: "forge-dashboard" and "forge-wiki" are internal proposals, not industry-validated; the rubric asks for "at least one industry-validated solution", which is met (four of them), but the internal proposals dilute the benchmarking density. |
| Honest trade-off comparison | 20/25 | Pros/cons are based on actual product characteristics. The selected approach's cons now include "叠加层模式未经大规模验证" — this is a genuine, honest limitation. Improvement from iteration 1. Deduction: the Cons for industry alternatives focus on what they DON'T do for Forge (manage process documents), which is true but frames them unfavorably. A more balanced analysis would acknowledge that these tools solve their target problems well — Forge's problem is simply different. |
| Chosen approach justified against benchmarks | 20/25 | The Innovation Highlights section (lines 109-117) now provides explicit differentiation: "上述工具都面向**代码产物**... Forge Workspace 面向**过程文档**". This is a clear comparative argument. The "Innovation Highlights" table maps specific borrowings. Improvement from iteration 1. Deduction: the justification explains HOW Forge Workspace differs, but not WHY this difference is important enough to warrant a new tool rather than adapting an existing pattern. Why not layer process document management on top of Cargo-style workspaces? |

**D3 Total: 100/120**

### D4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 34/40 | Six key scenarios cover happy paths. Edge cases addressed: zero-discovery diagnostics, unhealthy projects, manifest schema versioning, state machine preconditions for assign/close. Improvement from iteration 1. Missing error scenarios: (1) `assign` target project not registered in workspace; (2) running workspace commands outside a workspace directory; (3) `.forge-workspace.yaml` corrupted; (4) a project registered in two workspaces. These are operational failure modes for a multi-project tool. |
| Non-functional requirements | 33/40 | Performance target (<2s cold, <0.5s cached) quantified. Caching strategy with mtime fingerprint well-specified. "项目发现兼容任意子目录命名" present. Schema versioning strategy described (v0 default, visual `[v0]` tag). Improvement from iteration 1. Missing: (1) Disk space impact of workspace-level proposals and cache. (2) Concurrency — two terminal sessions running `status` simultaneously. (3) Windows path handling for `register <path>` — relevant given the project runs on Windows. |
| Constraints & dependencies | 26/30 | Dependencies on existing Forge CLI capabilities listed. Skill adaptation scope explicitly bounded (lines 142-147). Improvement from iteration 1. Missing: (1) Workspace directory is not a git repo — workspace-level proposals have no version control. Is this acceptable? (2) The flat-directory constraint (projects must be siblings) is stated but not discussed as a limitation for users who nest projects. |

**D4 Total: 93/110**

### D5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 32/40 | The "process document vs knowledge" orthogonal split is a genuine contribution. The Innovation Highlights table (lines 109-117) now articulates differentiation clearly: existing tools aggregate code artifacts; Forge Workspace aggregates process documents with lifecycle and state machines. The cross-project proposal staging area is novel. Deduction: the core mechanics (file scanning, aggregation, caching) remain standard. The novelty is in the domain application, not in technical innovation. |
| Cross-domain inspiration | 28/35 | Major improvement from iteration 1. Explicit borrowing table maps: Cargo workspaces (workspace.members aggregation concept), VS Code multi-root (projects as first-class citizens in a container), Turborepo (cross-project caching with fingerprints). The "核心差异" summary (line 117) synthesizes the cross-domain analysis. Deduction: no borrowing from non-technical domains (project management methodologies, knowledge management systems, document lifecycle standards). |
| Simplicity of insight | 21/25 | "Project parent directory is the natural workspace boundary" is elegant. The three-module decomposition is clean. The "overlay, don't migrate" principle avoids a common trap. Deduction: the machinery (caching, schema versioning, field inheritance, state machine, mtime fingerprints) is non-trivial for what could be simpler in v1. The 8-13d estimate confirms this is not a trivial feature. |

**D5 Total: 81/100**

### D6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 37/40 | All operations are file system reads/writes. No external services. Forge CLI has required infrastructure. Brainstorm adaptation scoped to single skill. Sound. Minor deduction: no spike or proof-of-concept to validate the workspace directory detection logic or the cache invalidation strategy. |
| Resource & timeline feasibility | 25/30 | Major improvement from iteration 1. Work breakdown table (lines 177-185) estimates 8-13 person-days across 6 items. This is realistic for file-system-based operations. Deduction: (1) The range is wide (8-13d = 63% variance). Which end is more likely? (2) No buffer for unknowns or integration testing. (3) The brainstorm skill adaptation estimate (1-2d) is honest but unvalidated — it depends on brainstorm's current architecture. |
| Dependency readiness | 28/30 | All dependencies are existing and available. No external APIs. No upstream blockers. Minor deduction: the proposal does not confirm brainstorm's architecture supports workspace context detection as a small change vs a major refactor. |

**D6 Total: 90/100**

### D7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | Six numbered scope items with specific commands, behaviors, and data structures. The `assign` field mapping table is model-level concrete. The `close` preconditions are explicit. Improvement from iteration 1. Minor gap: Items #5 and #6 are design principles, not deliverables — they describe constraints, not things to build. |
| Out-of-scope explicitly listed | 22/25 | Five items explicitly listed with "→ Wiki/Dashboard 模块独立提案" notation. Improvement from iteration 1 — the three-module architecture makes boundaries clearer. Missing: (1) Is `forge workspace unregister` in or out? (2) Is workspace-level PRD/design creation in scope? The proposal mentions brainstorm in workspace context but not the full pipeline (write-prd, tech-design, breakdown-tasks). Line 146 says "assign 后在目标项目内运行，无需适配" which implies they're NOT in workspace scope, but this should be explicit in Out of Scope. |
| Scope is bounded | 22/25 | Bounded by "v1" label, module boundaries, and the 8-13d estimate. Improvement from iteration 1 — the timeline estimate makes scope assessable. Deduction: (1) The 8-13d range is wide. (2) The brainstorm adaptation (1-2d) could expand if brainstorm's architecture is complex. (3) No explicit MVP definition — what is the minimum viable subset? |

**D7 Total: 72/80**

### D8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | Five risks covering: discovery failure, context loss on assign, module coupling, format inconsistency, cognitive confusion. Good coverage of solution's own risks. Missing from iteration 1 and still missing: (1) Risk of `.forge-workspace.yaml` corruption or accidental deletion — no recovery path. (2) Risk of workspace directory not being under version control — all workspace-level proposals are untracked by git. This was flagged as a "Should Fix" in iteration 1 and remains unaddressed. (3) Risk that the flat-directory constraint is violated by users with nested project structures. |
| Likelihood + impact rated | 25/30 | Ratings present and varied (L/M, M/M, L/H, M/M, M/L). The L/H for module coupling is honest and important. No rating feels dishonest. Minor deduction: "项目结构变化导致发现失败" rated L is debatable — if projects are actively developed, their internal structure evolves. "认知混淆" rated M/L seems right. |
| Mitigations are actionable | 25/30 | Most mitigations are concrete: "unhealthy" marking, field inheritance table, schema versioning, CLI namespace separation. The module coupling mitigation ("v1 不定义跨模块公开接口") is a design decision that prevents the risk but doesn't address eventual module interaction needs. The "close" command with preconditions is a concrete mitigation for zombie proposals. Improvement from iteration 1. Deduction: the mitigation for format inconsistency (schema versioning with v0 default) is good but the `[v0]` visual tag is a UX choice, not a technical mitigation for parsing failures. |

**D8 Total: 75/90**

### D9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 27/30 | Major improvement from iteration 1. Most criteria are specific commands with measurable outputs. SC#7 (line 282) now lists explicit fields (title, intent, status=Draft, source 链接). SC#8 (line 283) has clear preconditions. SC#11 (line 286) specifies cache behavior with concrete parameters (mtime fingerprints, 5-minute threshold). SC#12 (line 287) specifies v0 schema behavior with `[v0]` tag. Minor deduction: SC#6 (line 281) "brainstorm 技能可在 workspace 上下文运行" — what constitutes "可运行"? Creating a file in the correct directory? Using the correct template? The criterion is testable but the success condition could be more specific. |
| Coverage is complete | 22/25 | Major improvement from iteration 1. 12 SC entries now cover all four scope clusters. `register` (SC#2), `close` (SC#8), cache (SC#11), v0 schema (SC#12) — all gaps from iteration 1 are addressed. Remaining gap: (1) No SC for the `assign` field inheritance table completeness — SC#7 lists specific fields but doesn't verify the full table (e.g., frontmatter custom fields passthrough). (2) No SC for the state machine transition enforcement (assign requires Approved status). SC#7 tests the outcome but not the precondition rejection. |
| SC internal consistency | 22/25 | No contradictions within the SC set. State machine consistency verified: SC#7 requires proposal in Approved state (implicitly, per Scope line 246); SC#8 requires feature in Completed/Closed. These are consistent with the Draft → Approved → Assigned → Done state machine. Minor gap: SC#3 specifies both cold start <2s and cached <0.5s — the cache SC (SC#11) specifies the mechanism but there's no SC verifying the <0.5s warm performance target. The <0.5s is in SC#3 but the cache behavior is in SC#11; these are consistent but could be cross-referenced more explicitly. |

**D9 Total: 71/80**

### D10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 33/35 | Direct mapping: pain point 1 (scattered documents) -> status/features; pain point 2 (no global view) -> status table; pain point 3 (switching friction) -> unified workspace entry point; pain point 4 (knowledge) -> excluded, deferred to Wiki. Clean. Minor gap from iteration 1 remains: pain point 3 ("丢失上下文") — the workspace provides a global view but does not restore per-project agent context when switching. The solution addresses visibility (seeing all states), not context preservation (restoring session state). This is a gap between problem statement and solution scope. |
| Scope <-> Solution <-> SC aligned | 26/30 | Major improvement from iteration 1. All scope items now have corresponding SC entries. Field mapping table in Scope is testable via SC#7. State machine transitions in Scope are verifiable through SC#7 and SC#8 preconditions. Cache strategy in NFR has SC#11. Remaining tensions: (1) Scope describes `assign` field inheritance including "frontmatter 自定义字段" passthrough (line 242) but SC#7 only lists four specific fields — custom fields passthrough is unverified. (2) Scope describes state machine enforcement ("assign 要求 proposal 处于 Approved 状态，否则拒绝并提示", line 246) but no SC verifies the rejection behavior (only the happy path). |
| Requirements <-> Solution coherent | 21/25 | Key scenarios map to solution commands. Constraints section correctly identifies dependencies. One orphan from iteration 1: `.forge-workspace/config.yaml` for display preferences (aliases, etc.) is mentioned in Scope (line 219) but has no corresponding scenario, requirement, or SC entry. What is it for? When is it created? If it's v1 infrastructure, it needs a scenario. If it's future, it should be in Out of Scope. |

**D10 Total: 80/90**

---

## Phase 3: Blindspot Hunt

Findings the rubric dimensions did not surface:

1. **[blindspot] Git-less workspace persists**: Flagged in iteration 1 as a "Should Fix". Still unaddressed. The workspace directory (parent of projects) is not a git repository. All workspace-level proposals, cache files, and config live outside version control. If the developer's machine fails, workspace-level proposals are lost. This is particularly concerning for the "创意暂存区" (creative staging area) — the very documents the proposal emphasizes as valuable. The proposal does not discuss this at all.

2. **[blindspot] `assign` rejection path untested**: The Scope defines a state machine precondition ("assign 要求 proposal 处于 Approved 状态，否则拒绝并提示"), and a single-assignment constraint ("重复 assign 报错"). These are explicit failure behaviors. But no SC verifies that `assign` correctly rejects when: (a) proposal is in Draft status, (b) proposal is already assigned, (c) target project is not registered. The SC set only tests the happy path.

3. **[blindspot] Workspace-level proposal eval flow incomplete**: The Scope states "brainstorm/eval 技能在 workspace 上下文运行" (line 233) but the Constraints section (lines 142-147) only addresses brainstorm adaptation. What about eval? If workspace-level proposals go through eval-proposal, does that skill also need adaptation? The proposal implies eval works but does not declare the modification scope. This was flagged in iteration 1 and remains unaddressed.

4. **[blindspot] `.forge-workspace/config.yaml` zombie scope item**: Line 219 mentions `.forge-workspace/config.yaml` — "Workspace 模块自身配置（别名、显示偏好等）" — but this config file has no corresponding command, scenario, requirement, or SC entry. It appears in Scope but is never built, tested, or used. It should either be removed from v1 scope or given a concrete purpose.

5. **[blindspot] Workspace initialization idempotency still undefined**: Flagged in iteration 1. What happens when `forge workspace init` is run a second time? Does it re-scan and update? Does it preserve manual `register` entries? Does it warn about removed projects? The proposal describes `init` as a one-time operation but real usage requires idempotent behavior.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1. Problem Definition | 91 | 110 |
| D2. Solution Clarity | 115 | 120 |
| D3. Industry Benchmarking | 100 | 120 |
| D4. Requirements Completeness | 93 | 110 |
| D5. Solution Creativity | 81 | 100 |
| D6. Feasibility | 90 | 100 |
| D7. Scope Definition | 72 | 80 |
| D8. Risk Assessment | 75 | 90 |
| D9. Success Criteria | 71 | 80 |
| D10. Logical Consistency | 80 | 90 |
| **Total** | **868** | **1000** |

## Top 5 Attacks (Priority Order)

1. **D3 (Industry Benchmarking)**: Benchmarking is improved but still incomplete — 100/120 leaves 20 points. No CLI-based project tracking tools cited (Taskwarrior, todo.txt), no document lifecycle management patterns (Obsidian vaults, Notion databases). The Cons for industry alternatives are framed as "what they don't do for Forge" rather than honest assessment of their strengths in their own domain. The justification explains HOW Forge differs but not WHY this difference merits a new tool rather than adapting an existing pattern.

2. **D4 (Requirements Completeness)**: Error scenarios still incomplete — 93/110. No coverage for: corrupted `.forge-workspace.yaml`, `assign` to unregistered project, running workspace commands outside a workspace, overlapping workspaces. The git-less workspace constraint is unacknowledged. Windows path handling for `register <path>` is unmentioned despite the project running on Windows.

3. **D9 (Success Criteria)**: Happy-path bias — 71/80. The `assign` command has explicit rejection behaviors (wrong state, already assigned, unregistered project) defined in Scope, but no SC verifies these rejection paths. The field inheritance table includes "frontmatter 自定义字段透传" but SC#7 only lists four specific fields. Cache warm performance (<0.5s) is in SC#3 but the cache behavior SC#11 doesn't cross-reference it.

4. **D8 (Risk Assessment)**: Git-less workspace risk persists — 75/90. Flagged in iteration 1 as "Should Fix", not addressed. Workspace-level proposals have no version control. No risk identified for `.forge-workspace.yaml` corruption. No mitigation for flat-directory constraint violation.

5. **D10 (Logical Consistency)**: Orphan scope item — 80/90. `.forge-workspace/config.yaml` appears in Scope (line 219) with a purpose ("别名、显示偏好等") but has no scenario, requirement, or SC. If it's infrastructure, it needs a scenario. If it's future, move to Out of Scope. Pain point 3 ("上下文丢失") still maps imperfectly — workspace provides visibility, not context restoration.

## Revision Guidance

### Must Fix (blocking approval)

1. **Acknowledge the git-less workspace risk** — either add a mitigation (e.g., `init` initializes a git repo in workspace root for proposal versioning) or explicitly accept the risk in Key Risks with rationale. This is a significant operational gap for a tool managing "创意暂存区" documents.

2. **Add SC entries for `assign` rejection paths** — at minimum one SC verifying that `assign` rejects when proposal is not in Approved status, and one verifying rejection when already assigned. These are explicit behaviors promised in Scope.

### Should Fix (significant quality improvement)

3. **Resolve `.forge-workspace/config.yaml` zombie** — either give it a concrete v1 purpose (with scenario and SC) or remove it from v1 Scope. Currently it's a scope item with no justification.

4. **Add 1-2 failure scenarios** to Requirements Analysis: corrupted config, running outside workspace, overlapping workspaces.

5. **Specify `init` idempotency behavior** — what happens on re-run? This affects the implementation design.

### Nice to Have

6. Add eval skill to the workspace context adaptation scope (or explicitly exclude it with rationale).
7. Cross-reference cache warm performance (<0.5s in SC#3) with cache behavior SC#11.
8. Add one CLI-based project tracking tool (Taskwarrior, etc.) to the Industry Benchmarking table.
