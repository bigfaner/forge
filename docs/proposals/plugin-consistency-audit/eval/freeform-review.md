# Freeform Expert Review: Forge Plugin 内部一致性审计 Proposal

**Reviewer**: Documentation-Implementation Drift Auditor
**Date**: 2026-05-29
**Document**: `docs/proposals/plugin-consistency-audit/proposal.md`

---

## Section 1: Background Assessment

This proposal addresses a concrete and well-motivated problem: the Forge plugin has undergone two significant architectural refactorings (test profile system and intent-driven pipeline branching), and the resulting changes were applied locally to individual files rather than validated globally. The proposal argues that this creates a high risk of internal inconsistency within each component — where a skill's SKILL.md says one thing while its templates/rules/data files assume another — and proposes a systematic audit of all components to produce a structured issue report.

The approach is scoped deliberately: audit individual components for internal self-consistency rather than cross-component coordination. This "single-component self-consistency" framing is the proposal's central design decision. The audit targets 22 skills, 18 commands, and 1 agent, using four classification categories (CONFLICT, REDUNDANT, TIMING, REFERENCE), and produces a structured report with severity levels (P0–P3) without performing actual fixes.

The proposal positions itself against three alternatives (do nothing, manual review, automated schema validation) and selects AI-assisted layered auditing as the most cost-effective approach for the current file count and urgency.

---

## Section 2: Key Risk Identification

### 2.1 Factual Accuracy: Component Count Discrepancy

问题：The proposal states "22 个 skill" as a factual claim in multiple locations, but the current codebase contains exactly 21 skill directories under `plugins/forge/skills/`. This is not an approximate count — it appears as a specific, auditable number used in success criteria.

Direct quote from the proposal's Non-Functional Requirements:
> "审计覆盖率: 100% 的 skill（22个）、command（18个）、agent（1个）"

And again in Success Criteria:
> "22 个 skill 100% 覆盖审计，每个 skill 的 SKILL.md 与其 templates/rules/data 逐一对比"

Consequence: If this proposal were executed as written, the auditor would be looking for a 22nd skill that does not exist. This either means (a) the proposal was drafted against a different codebase snapshot where a 22nd skill existed and has since been removed, or (b) the count was never verified against the actual directory listing. Either way, the success criterion "22 个 skill 100% 覆盖审计" is not verifiable against the current codebase. This is precisely the type of factual drift that this proposal aims to detect in the plugin itself — the proposal has its own factual drift problem.

### 2.2 File Count Claim Unsupported

问题：The proposal claims "172+ 个 .md 文件通过手动维护交叉引用" but the actual count is 208 total `.md` files in the plugin directory, or 182 when excluding eval-specific specialist files (experts/scorer, experts/freeform, experts/protocol, rubrics). Neither number is "172+."

Direct quote from Evidence section:
> "172+ 个 .md 文件通过手动维护交叉引用，重构过程中依赖局部修改而非全局验证"

Consequence: The "172+" figure may have been accurate at an earlier commit, but it is presented as a current-state claim. More importantly, this number is used to argue urgency and scale — if the number is wrong, the urgency argument is weakened or the scope is misestimated. An audit proposal that cannot accurately count the files it proposes to audit undermines its own credibility.

### 2.3 Incomplete Subdirectory Taxonomy

问题：The proposal describes each skill's supporting files as "templates/rules/data" but two additional subdirectory types exist in the codebase that are not mentioned: `examples/` (in tech-design and write-prd skills) and `types/` (in gen-test-scripts). The proposal also omits the `experts/` directory under the eval skill.

Direct quote from Proposed Solution:
> "逐一检查每个 skill 的 SKILL.md 与其 templates/rules/data 之间"

And from the Scope section:
> "22 个 skill 的 SKILL.md 与其各自的 templates/rules/data 之间的逻辑自洽性"

Consequence: If the audit proceeds with this taxonomy, files in `examples/`, `types/`, and `experts/` directories may be missed entirely. Specifically: tech-design has 2 example files that guide SKILL.md behavior, gen-test-scripts has 6 type-specific files that define per-surface-type behavior, and eval has 15 expert files plus 11 rubric files. These are not edge cases — the gen-test-scripts types are central to its test profile system, which is one of the two refactorings that motivated this proposal.

### 2.4 Out-of-Scope Boundary Creates Blind Spot for eval Skill

风险：The proposal explicitly excludes "rules/rubrics/experts 的功能性质量审查" from scope, but the eval skill has 10 rules files, 11 rubric files, and 15 expert files — making it by far the most file-rich skill in the plugin (36 of 188 skill `.md` files, or ~19%). Declaring these out of scope for "functional quality review" is reasonable, but the proposal conflates "functional quality" with "internal consistency."

Direct quote from Out of Scope:
> "rules/rubrics/experts 的功能性质量审查"

Consequence: The eval skill's SKILL.md references specific rubric file paths and expert roles. If those references are stale (e.g., a rubric was renamed or an expert role was split), that is an internal consistency issue (REFERENCE category), not a functional quality issue. By excluding rubrics/experts from scope without distinguishing between consistency checking and quality evaluation, the proposal may leave ~19% of skill files un-audited despite claiming 100% skill coverage.

### 2.5 hooks/guide.md Auditability Undefined

问题：The proposal includes "hooks/guide.md 的内部一致性" in scope, but guide.md is a single standalone file. The proposal does not define what "internal consistency" means for a single file with no sub-components.

Direct quote from In Scope:
> "hooks/guide.md 的内部一致性"

Consequence: guide.md contains references to CLI commands (`forge proposal`, `forge feature status`, `forge task transition`, `forge task reopen`), file paths (`docs/business-rules/`, `docs/conventions/`), and skill names (`/consolidate-specs`, `/learn`). "Internal consistency" for this file should mean verifying that these references are accurate — but that is cross-reference validation, which the proposal explicitly scopes out for skills ("跨 skill 之间的冗余内容" is out of scope). The proposal needs to clarify whether guide.md references to CLI commands and file paths will be validated, or whether "internal consistency" for guide.md is a null operation.

### 2.6 Classification Scheme Coverage Gap

问题：The four-category classification (CONFLICT, REDUNDANT, TIMING, REFERENCE) does not cover all failure modes that are likely to be found in a prompt-based plugin system. Specifically missing:

- **DEPRECATED/OBSOLETE**: Instructions that reference removed features or capabilities (e.g., references to `test.execution` from `.forge/config.yaml` that was deliberately removed)
- **AMBIGUITY**: Instructions that are internally unclear or have multiple valid interpretations within a single component
- **INCOMPLETENESS**: SKILL.md describes a step that has no corresponding template or rule to support it

Direct quote from the classification:
> "问题分类: 矛盾(CONFLICT)、冗余(REDUNDANT)、时序(TIMING)、引用(REFERENCE)"

Consequence: When the auditor encounters an issue that doesn't fit cleanly into these four categories, it will either be forced into the nearest match (distorting the classification) or left unclassified. The INCOMPLETENESS failure mode is particularly likely given the refactoring history — if a step was removed from a SKILL.md but its supporting template was not, or vice versa, that's not a CONFLICT or REDUNDANT, it's a structural gap.

### 2.7 Severity Framework Underdefined

风险：The proposal mentions P0–P3 severity levels in success criteria but provides no definition of what each level means. In a plugin system where documentation is runtime-consumed by agents, severity classification must account for blast radius — a broken cross-reference in a SKILL.md that causes an agent to read the wrong template at runtime is functionally a P0 bug, not a P2 documentation issue.

Direct quote from Success Criteria:
> "每个问题包含: 文件路径、问题描述、严重等级(P0-P3)、修复建议"

Consequence: Without severity definitions, different auditors (or the same auditor on different days) will classify identical issues at different severity levels. This makes the report non-reproducible and undermines the prioritization goal stated in the risk table: "每个问题标注严重等级（P0-P3）" is listed as mitigation for "报告问题过多导致修复优先级不清."

### 2.8 Assumption That Single-Component Auditing Is Sufficient

风险：The proposal's core design decision — auditing only single-component self-consistency — assumes that cross-component inconsistencies are "设计层面的合理重复" and therefore low-risk. But in the Forge plugin, cross-component references are pervasive: quick-tasks references brainstorm templates, run-tests references init-justfile rules, submit-task references task types defined in breakdown-tasks. The line between "internal" and "cross-component" is not as clean as the proposal assumes.

Direct quote from Innovation Highlights:
> "审计按'单一组件自洽'而非'跨组件协调'组织——这降低了审计复杂度，同时覆盖了最可能出问题的维度（组件内部重构后的残留不一致）。跨组件冗余是设计层面的合理重复，不在此次审计范围内。"

Consequence: A concrete example: if quick-tasks SKILL.md says "intent is read from proposal frontmatter" but brainstorm templates generate proposals without an `intent` field, that's a cross-component inconsistency that directly causes a runtime failure. The current scope boundary would miss this because each component is internally consistent on its own. The proposal's scope boundary trades completeness for tractability without acknowledging the specific cross-component failure modes that the recent refactorings introduced.

### 2.9 Proposal Status Inconsistency

问题：The proposal's frontmatter declares `status: Draft` and `intent: "refactor"`, but the document is presented as ready for execution. A refactoring intent would skip the test pipeline (per the intent-driven branching rules), yet this proposal is itself a documentation audit that produces no testable code — so the intent classification is irrelevant and potentially confusing.

Direct quote from frontmatter:
```yaml
status: Draft
intent: "refactor"
```

Consequence: Minor but indicative — if the proposal's own metadata is misclassified, it suggests the metadata conventions may not be well-understood by the system that generates them.

---

## Section 3: Improvement Suggestions

建议：Correct the component counts before execution. Replace "22 个 skill" with the actual count (21) throughout the proposal, or identify which skill was removed/merged and document the discrepancy. The success criteria must be verifiable against the actual codebase state at audit time. This addresses the risk in Section 2.1 — a wrong count in an audit proposal that claims 100% coverage is a credibility issue, not just a typo.

建议：Verify the "172+ 个 .md 文件" claim against the current codebase and either update it or clarify what subset it refers to (e.g., "SKILL.md + rules/ + templates/ + data/ files, excluding eval specialists"). The current state is 208 total or 182 core files. This addresses Section 2.2 — the file count is used to argue urgency and must be accurate.

建议：Expand the subdirectory taxonomy in the Scope section to include all directory types that exist in the codebase: `templates/`, `rules/`, `data/`, `types/`, `examples/`, and `experts/` (for the eval skill). The audit methodology should explicitly enumerate which subdirectories are checked for each skill. This addresses Section 2.3 — the current taxonomy misses `examples/`, `types/`, and `experts/` directories that contain files directly referenced by their SKILL.md.

建议：Refine the out-of-scope boundary for eval skill rubrics and experts. Instead of excluding "rules/rubrics/experts 的功能性质量审查," exclude "rubrics/experts 的 prompt engineering 质量" while keeping cross-reference validation in scope. Specifically: verify that SKILL.md and eval rules correctly reference existing rubric and expert file paths, even if the content quality of those files is not evaluated. This addresses Section 2.4 — the current scope excludes ~19% of skill files without distinguishing between consistency checking and quality evaluation.

建议：Define what "hooks/guide.md 的内部一致性" means operationally. At minimum, this should include verifying that CLI commands referenced in guide.md exist in the current forge CLI, and that file paths and skill names referenced are accurate. This is a cross-reference check, not a self-consistency check, and should be explicitly acknowledged as such. This addresses Section 2.5.

建议：Add a fifth classification category — INCOMPLETE — to cover structural gaps where a SKILL.md describes behavior that has no supporting file, or where a supporting file exists but is never referenced by SKILL.md. Optionally add DEPRECATED for references to removed capabilities. This addresses Section 2.6 — the current four categories do not cover all failure modes that a real audit would encounter.

建议：Define severity level criteria explicitly. For a runtime-consumed documentation system, I recommend: P0 = will cause agent to crash or produce wrong output (broken file references, missing templates); P1 = will cause agent to skip important steps or produce incomplete output (missing rules, contradictory instructions); P2 = causes confusion or inefficiency but agent can recover (redundant instructions, vague wording); P3 = cosmetic or stylistic issues. This addresses Section 2.7 — without definitions, severity assignments will be inconsistent and the prioritized report will be unreliable.

建议：Add a scoped cross-component reference check as a second audit pass, limited to the specific interfaces that were modified during the two refactorings. Specifically: verify that quick-tasks/breakdown-tasks correctly read the `intent` field from brainstorm-generated proposals, and that run-tests correctly references init-justfile rules after the test profile system change. This is a targeted exception to the single-component scope, not a blanket cross-component audit. This addresses Section 2.8 — the most dangerous inconsistencies in this codebase are at the seams between refactored components, not within individual components.
