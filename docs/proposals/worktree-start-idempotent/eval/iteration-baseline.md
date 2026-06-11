# Proposal Evaluation Report — Baseline

**Document**: `proposal.md` — Make `forge worktree start` Idempotent & Rename `copy-files` → `includes`
**Date**: 2026-06-09
**Eval Type**: Baseline (informational, not a gate)
**Evaluator Role**: CTO Adversary

---

## Phase 1: Reasoning Audit

### 1. Problem → Solution: Does the proposed solution actually address the stated problem?

**Problem 1**: `start` errors when worktree exists; no way to enter existing worktree with a fresh session.

Solution: Make `start` idempotent — skip creation if exists, launch fresh session. **Directly addresses the problem.** The mapping is clean and unambiguous.

**Problem 2**: `copy-files` naming describes implementation, not intent.

Solution: Rename to `includes`. **Directly addresses the problem.** However, the proposal bundles two orthogonal changes without justification (see Logical Consistency).

### 2. Solution → Evidence: Does evidence support the solution?

Evidence is experiential (the author describes current behavior), not data-driven. There are no user surveys, issue tracker counts, or frequency metrics. For an internal tool enhancement of this scope, experiential evidence is adequate but not rigorous.

### 3. Evidence → Success Criteria: Do the SC test what matters?

The Success Criteria cover the key behavioral changes (idempotent start, branch warning, includes skip, backward compatibility, config rename, no legacy code). However, there is a notable gap: no SC verifies that the user actually sees distinguishable output for the two paths (new vs. existing worktree) — SC #2 says "明确的提示信息" but is not testable against a specific output format. Also missing: no SC for the `--no-launch` scenario listed in Key Scenarios #5.

### 4. Self-contradiction check

**Contradiction found**: In "Constraints & Dependencies", the proposal says `CopyFiles → Includes + 兼容读取` (compatibility reading). But in "In Scope", it says `直接替换，不保留旧字段` (direct replacement, no backward compatibility). And in Success Criteria #7: `代码中不存在任何 copy-files / CopyFiles 兼容逻辑`. These three statements are mutually contradictory. The proposal cannot simultaneously do "compatibility reading" and "no legacy code at all."

---

## Phase 2: Rubric Scoring

### Dimension 1: Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 36/40 | Two problems clearly stated with concrete context. Minor ambiguity: Problem 2 (rename) is tangential to Problem 1 (idempotency) and bundling them without explaining why they must be solved together slightly muddies the problem statement. |
| Evidence provided | 28/40 | Evidence is experiential only — four bullet points describing current behavior. No quantitative data: how many users hit this? How often? Are there GitHub issues or user complaints? The evidence is concrete but entirely anecdotal. Quote: "start 在 .forge/worktrees/<slug> 已存在时直接报错" — this is a factual observation, not user-impact evidence. |
| Urgency justified | 24/30 | "日常高频操作" is claimed but not substantiated. "Cost of delay：持续的体验摩擦" is a reasonable argument but lacks specificity. What is the actual frequency? How many users are affected? |

**Dimension 1 Total: 88/110**

### Dimension 2: Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | The idempotent behavior is precisely specified with two clear branches (exists/does not exist). A developer can implement this without further clarification. |
| User-facing behavior described | 40/45 | Good coverage of the two main paths and several edge cases (--source-branch, --no-launch, --interactive). Deduction: the actual output the user sees is underspecified. Quote: "输出区分性日志" — what format? What level? What exact text? |
| Technical direction clear | 30/35 | Files to modify are named, and the approach (modify error branch to skip+launch) is clear. Deduction: "约 40 行变更" is a claim not backed by analysis. The contradictory mention of "兼容读取" vs "直接替换" muddies the technical direction for the config rename part. |

**Dimension 2 Total: 108/120**

### Dimension 3: Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 28/40 | Only `kubectl apply`, `terraform plan`, `mkdir -p`, and `touch` are mentioned — all in a single sentence. No depth of analysis. Quote: "幂等 CLI 命令是常见模式：kubectl apply、terraform plan 都采用「已存在则跳过」的策略" — this is a passing mention, not a benchmark. How do these tools handle the equivalent of --source-branch conflicts? What about their deprecation strategies for config renames? Not discussed. |
| At least 3 meaningful alternatives | 24/30 | Four alternatives listed including "do nothing". The alternatives are reasonable but lack depth. "新增 open 子命令" and "给 resume 加 --fresh flag" are presented as straw men with single-sentence cons. Quote: "语义矛盾（resume + fresh 自相矛盾）" — this is a straw-man dismissal; `git checkout --force` also has "contradictory" semantics but works fine. |
| Honest trade-off comparison | 18/25 | The comparison table is present but shallow. Pros/cons are one-liners without quantitative analysis. The "Cons" for the selected approach is "轻微改变 start 的现有语义" — this understates the risk of breaking user scripts that depend on the error behavior. |
| Chosen approach justified against benchmarks | 18/25 | Justified by "最小惊讶原则" and `mkdir -p` analogy, which is reasonable. But no analysis of how other tools handle the equivalent of the config rename half of the proposal. The benchmarking completely ignores the rename aspect. |

**Dimension 3 Total: 88/120**

### Dimension 4: Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 32/40 | Six scenarios identified, which is good coverage. Missing scenarios: (a) worktree exists but is in a corrupted/broken state (e.g., after failed git operations); (b) concurrent `start` calls on the same slug; (c) worktree exists but its backing branch has been deleted remotely. |
| Non-functional requirements | 24/40 | Only two NFRs listed: backward compatibility and performance. Missing NFRs that are critical for this change: (a) observability/logging requirements (only mentioned tangentially in In Scope); (b) migration/upgrade path for existing configs; (c) documentation requirements; (d) error message quality requirements. Quote: "性能：无影响（只是跳过了创建步骤）" — this is trivially true and suggests the NFR section was treated as a checkbox. |
| Constraints & dependencies | 20/30 | Files to modify are named, and the dependency on existing worktree validation logic is noted. But the section contains the "兼容读取" contradiction noted above, undermining its reliability. No mention of documentation updates, changelog requirements, or test coverage needs. |

**Dimension 4 Total: 76/110**

### Dimension 5: Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 20/40 | The proposal explicitly states this is not novel: "幂等设计，不是业界首创". The `mkdir -p` analogy is apt but there is no innovation beyond applying a well-known pattern. The config rename is a straightforward renaming exercise. |
| Cross-domain inspiration | 12/35 | Only CLI tools are referenced. No inspiration from other domains (e.g., how IDEs handle "open project" when project is already open, how container runtimes handle idempotent container start, how package managers handle already-installed packages). |
| Simplicity of insight | 20/25 | The core insight ("start should mean 'start working', not 'create worktree'") is genuinely elegant and simple. This is the proposal's strongest creative contribution. |

**Dimension 5 Total: 52/100**

### Dimension 6: Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | Clearly feasible. The change is well-scoped technically and relies on existing patterns in the codebase. Minor concern: the claim of "40 lines" may underestimate the config rename's blast radius across the codebase. |
| Resource & timeline feasibility | 24/30 | "1 小时内完成" for 2 changes that touch config schema, command logic, and require full search-replace of `CopyFiles` across the codebase is optimistic. The estimate does not account for testing, documentation updates, or handling edge cases discovered during implementation. |
| Dependency readiness | 30/30 | No external dependencies. Correctly assessed. |

**Dimension 6 Total: 90/100**

### Dimension 7: Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | Five concrete deliverables. Each is a specific behavioral change. Good. |
| Out-of-scope explicitly listed | 22/25 | Six items explicitly deferred. Good coverage. The glob/directory support and `.worktreeinclude` file format deferrals show awareness of scope creep risks. |
| Scope is bounded | 16/25 | The scope claims "40 lines" but the actual impact of the config rename is not bounded. The proposal does not identify how many files reference `CopyFiles` or `copy-files` across the codebase. Additionally, bundling two orthogonal changes in one scope makes the boundary less crisp than it should be. |

**Dimension 7 Total: 64/80**

### Dimension 8: Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 22/30 | Three risks identified. Missing risks: (a) silent config breakage from rename (the biggest risk, identified in the freeform review but not in the proposal itself); (b) `resume` command discoverability reduction (a UX risk); (c) corrupted/stale worktree state being masked by idempotent behavior; (d) the risk of bundling two orthogonal changes. |
| Likelihood + impact rated | 20/30 | Ratings are present but questionable. The "用户脚本依赖报错行为" risk is rated L likelihood — but any scripting user will hit this immediately upon upgrade. The "includes 被跳过" risk is rated L/L, but combined with the silent config breakage risk, this is actually higher impact than acknowledged. |
| Mitigations are actionable | 22/30 | Mitigations are partially actionable. "在 release notes 中标注" is actionable. "输出 warning" is actionable but the warning content is not specified. "首次创建时已复制" is not a mitigation — it is a justification for why the risk is acceptable, which is different from an action to take. |

**Dimension 8 Total: 64/90**

### Dimension 9: Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 22/30 | Most criteria are testable (e.g., "不再报错", "输出 warning", "被跳过"). But SC #2 is vague: "明确的提示信息" — what constitutes "clear"? SC #7 ("不存在任何 copy-files / CopyFiles 兼容逻辑") is testable via grep but sets a problematic bar if the implementation needs any migration logic. Missing: no SC for `--no-launch` scenario (Key Scenario #5), no SC for `--interactive` scenario (Key Scenario #6). |
| Coverage is complete | 16/25 | Gaps identified: (a) no SC for `--no-launch` with existing worktree; (b) no SC for `--interactive` with existing worktree; (c) no SC verifying that the "区分性日志" has specific, parseable format; (d) no SC for upgrade/migration experience. The In Scope item about "输出区分性日志" is only weakly covered by SC #2. |
| SC internal consistency | 18/25 | SC #7 ("不存在任何兼容逻辑") directly contradicts the Constraints section's mention of "兼容读取". If the implementation truly has zero compatibility code, then the Constraints section is wrong; if it has compatibility reading, then SC #7 cannot be satisfied. This is a genuine internal contradiction within the SC set. Additionally, SC #5 ("worktree 不存在时行为完全一致") and the overall idempotent change are compatible, but the proposal does not define what "完全一致" means for test verification — byte-exact output matching? Behavioral matching? |

**Dimension 9 Total: 56/80**

### Dimension 10: Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 30/35 | The idempotent `start` cleanly addresses Problem 1. The rename cleanly addresses Problem 2. However, the proposal never justifies why these two changes must be bundled together. They are orthogonal: idempotency does not require a rename, and a rename does not require idempotency. The bundling introduces unnecessary coupling. |
| Scope ↔ Solution ↔ Success Criteria aligned | 18/30 | The contradiction between "兼容读取" (Constraints) and "不存在任何兼容逻辑" (SC #7) is a direct Scope ↔ SC misalignment. Additionally, the Scope includes "输出区分性日志" as a deliverable but the SC only vaguely addresses it. Key Scenarios #5 and #6 (`--no-launch`, `--interactive`) appear in requirements but not in Success Criteria — a coverage gap. |
| Requirements ↔ Solution coherent | 18/25 | The requirements and solution are mostly coherent, but the config rename has orphan implications: it is a solution feature (rename config) that addresses a naming concern, but the "non-functional requirement" of upgrade/migration is absent from Requirements. The solution does more (breaking change) than the requirements section acknowledges. |

**Dimension 10 Total: 66/90**

---

## Phase 3: Blindspot Hunt

**[blindspot] Silent breaking change on config rename**: The rubric does not have a dimension that specifically penalizes silent breaking changes. The rename from `copy-files` to `includes` with no migration path means existing user configs will silently stop working. YAML parsers do not error on unknown fields — the old `copy-files` key will simply be ignored, and `includes` will be absent, resulting in no files being copied. This is the most dangerous aspect of the proposal and the rubric's Risk Assessment dimension partially catches it, but there is no dimension that evaluates "breaking change management" as a first-class concern.

**[blindspot] Two-change coupling**: The rubric evaluates the proposal as a monolith. There is no dimension that rewards or penalizes proposals for change coupling — i.e., bundling orthogonal changes together. In this proposal, idempotency and config rename are independently valuable and independently shippable. Bundling them increases rollback risk and makes it harder to attribute issues. No rubric dimension captures this.

**[blindspot] Missing user migration narrative**: The rubric's Requirements Completeness dimension checks for "constraints & dependencies" but does not explicitly require a migration/upgrade strategy for existing users. For any change that alters config schema or CLI behavior, the migration path should be a first-class requirement.

**[blindspot] Testability of "完全一致"**: SC #5 says behavior for non-existing worktree should be "与当前完全一致" (completely identical to current). This sets a very high bar — does it mean byte-exact CLI output matching? Exit code matching? Behavioral matching only? The rubric's "measurable and testable" criterion flags vagueness but does not specifically address temporal regression claims like "behavior unchanged."

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 88 | 110 |
| Solution Clarity | 108 | 120 |
| Industry Benchmarking | 88 | 120 |
| Requirements Completeness | 76 | 110 |
| Solution Creativity | 52 | 100 |
| Feasibility | 90 | 100 |
| Scope Definition | 64 | 80 |
| Risk Assessment | 64 | 90 |
| Success Criteria | 56 | 80 |
| Logical Consistency | 66 | 90 |
| **Total** | **752** | **1000** |

**Target: 900/1000 — NOT MET (gap: 148 points)**

---

## Top 5 Actionable Improvements (Priority Order)

1. **Resolve the "兼容读取" vs "直接替换" contradiction.** Choose one strategy, update all three locations (Constraints, In Scope, Success Criteria) to be consistent. This alone would improve Logical Consistency (+10-15) and Success Criteria (+5-8).

2. **Add a migration/upgrade strategy for the config rename.** At minimum: detect old `copy-files` and emit a loud deprecation warning. Ideally: support both keys for one version cycle. This would improve Risk Assessment (+10-15) and Requirements Completeness (+8-10).

3. **Deepen industry benchmarking.** Analyze how 2-3 specific tools handle the analogous scenarios (not just name-drop them). How does `kubectl` handle idempotent apply when parameters differ? How does Terraform handle renamed config keys? This would improve Industry Benchmarking (+15-20).

4. **Add missing Success Criteria for Key Scenarios #5 and #6.** `--no-launch` and `--interactive` scenarios appear in requirements but not in SC. Also make SC #2 ("明确的提示信息") testable by specifying the expected output format. This would improve Success Criteria (+10-15).

5. **Justify or separate the two-change coupling.** Either argue why these changes must ship together, or split into two proposals. This would improve Logical Consistency (+5-8) and Scope Definition (+5-8).
