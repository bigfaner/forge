# Freeform Expert Review: Forge Plugin Skill System Consistency Fix

**Reviewer**: Senior Quality Engineer — plugin-architecture-audit, cross-reference-verification, template-placeholder-consistency, rubric-scale-validation, pipeline-contract-alignment

**Date**: 2026-06-10

**Document reviewed**: `proposal.md` (current state, iteration 5+)

---

## Section 1: Background Assessment

This proposal addresses consistency drift in the Forge Plugin v3.0.0-rc.53 codebase after large-scale refactoring. The plugin comprises 21 skills, 16 commands, and 5 hook events. The audit claims to have found 23 issues across 8 dimensions, categorized as 4 HIGH, 8 MEDIUM, and 11 MINOR severity.

The core technical approach is text-level fixes to markdown configuration files with no Go code changes. This is a narrow-scope remediation focused on correcting stale references between files that function as a distributed configuration system: rubric frontmatter acts as the primary truth source, `rubric-reference.md` acts as a secondary cache, command `argument-hint` and `description` fields act as tertiary projections, and template files act as code generation inputs.

The proposal rests on several assumptions: (a) the rubric frontmatter values are the single source of truth, (b) the fix scope is limited to markdown files only, (c) the current Go runtime behavior is correct and only the documentation/LLM-facing files are wrong, and (d) the audit was comprehensive enough to catch all significant inconsistencies.

---

## Section 2: Key Risk Identification

### Quantitative Claims vs. Source Verification

风险：The proposal's Problem section states "22 个 skill", but the actual `plugins/forge/skills/` directory contains 21 subdirectories. This is a minor factual error, but it undermines confidence in the audit's completeness claim. If the auditor miscounted the most basic unit of analysis, what else was miscounted? The text reads: "Forge Plugin v3.0.0-rc.53（22 个 skill、16 个命令、5 个 hook）经历大规模重构后".

问题：The Summary Statistics table for dimension C states "约40 模板文件" as the audit scope, but the actual `plugins/forge/skills/*/templates/*.md` path contains 56 files. A 40% undercount of the audit scope raises questions about whether template-level analysis was thorough. If the auditor only examined roughly 40 of 56 templates, 16 templates were excluded from scrutiny.

问题：M-9 states "INLINE 跨 skill 引用同步风险（3 处）", listing three INLINE locations. Actual grep of `<!-- INLINE:origin=` across all skill SKILL.md files reveals 4 INLINE references. The missing one is `gen-contracts/SKILL.md` line 58: `<!-- INLINE:origin=gen-journeys/SKILL.md#Surface Detection -->`. This creates a bidirectional INLINE dependency: gen-contracts inlines from gen-journeys (Surface Detection section), and gen-journeys inlines from gen-contracts (`journey-contract-model.md`). The proposal lists each direction independently under M-9 but does not recognize or flag the mutual dependency cycle. If `journey-contract-model.md` changes in gen-contracts AND Surface Detection changes in gen-journeys, each skill carries a stale copy of the other's content, and neither update propagates correctly.

### Missed Audit Findings in Eval Command Descriptions

风险：The proposal identifies the `1000-point` description in `eval-journey.md` and `eval-contract.md` as stale (H-1), but misses a more fundamental description inaccuracy: both command descriptions list the rubric dimensions incompletely. `eval-journey.md` description reads "Scores completeness, semantic purity, precondition exclusivity, fact alignment, surface fitness, and internal consistency" — that is 6 dimensions. The actual journey rubric at `skills/eval/rubrics/journey.md` has 7 dimensions, with "Workflow Coverage (工作流覆盖度)" at 150 points being the omitted dimension. Similarly, `eval-contract.md` description reads "Scores six-dimension structural integrity" and lists 6 evaluation areas, but the contract rubric at `skills/eval/rubrics/contract.md` has 8 dimensions (Anchor Integrity and Fixture Specification are missing from the description). These description inaccuracies are not cosmetic — they are read by the LLM as context during eval execution and may cause the scorer to skip dimensions it does not know exist.

问题：The `eval/SKILL.md` description field reads "Supports 100-point and 1000-point scales" but no rubric in the system uses a 100-point scale (all are 1000, 1100, or 1150). Additionally, the claim of "1000-point scales" is incomplete because journey (1150) and contract (1100) rubrics exceed 1000. This description is read by the LLM to determine valid score ranges. If the scorer believes the scale ceiling is 1000, it may cap scores for journey/contract evaluations. The proposal's H-1 fix updates `rubric-reference.md` and command-level `argument-hint`/`description`, but does not address the eval SKILL.md description itself.

### Fix Completeness Risks

风险：The H-1 fix proposes updating `rubric-reference.md`, `eval-journey.md` argument-hint and description, and `eval-contract.md` argument-hint and description. However, the H-1 analysis itself identifies "5 处" that need synchronization: "rubric frontmatter、rubric-reference.md、命令 argument-hint、命令 description、config 键默认值". The proposed fix addresses 4 of the 5 but does not mention updating config key defaults. If `eval.journey.target` or `eval.contract.target` config keys have hardcoded default values in Go code, these defaults remain at the stale values. The proposal marks Go code changes as out-of-scope, which is defensible, but the H-1 fix description should explicitly acknowledge this gap rather than implying the fix is complete.

风险：The M-9 fix proposes adding version stamps like `<!-- INLINE from ... @ v3.0.0-rc.53 -->` to INLINE references. This approach has a fundamental weakness: version stamps are only checked when someone remembers to grep for them. There is no mechanism to detect drift between versions. If the next release is v3.1.0 and someone updates `journey-contract-model.md` without checking all INLINE consumers, the version stamp at the consumer site still says `v3.0.0-rc.53`, and the grep check only works if someone runs it. The bidirectional dependency between gen-contracts and gen-journeys makes this particularly dangerous because changes in either direction can create silent semantic gaps.

问题：The M-1 fix proposes renaming config keys from camelCase to kebab-case in markdown files only, with Go alias compatibility marked as a follow-up task. The proposal's Risk table acknowledges: "如果 Go config reader 不支持 kebab-case，用户按新 key 配置后 config 读取失败回退到默认值". However, the Success Criteria for M-1 includes the precondition: "验证 Go config reader 是否支持 kebab-case 查询；如不支持，将 M-1 与 Go alias 绑定为原子操作，推迟至 Go 代码变更窗口执行". This creates a conditional success criterion that may cause M-1 to be skipped entirely if the Go config reader check fails. The proposal should clarify whether M-1 is being committed to in this cycle or conditionally deferred.

### H-2 Dead Path Analysis Completeness

问题：The H-2 fix proposes removing the dead path `docs/features/<slug>/proposal.md` from `tech-design/SKILL.md` line 47. However, the same file at line 24 references `docs/features/<slug>/prd/prd-spec.md` as a prerequisite check, and line 35 references `docs/features/<slug>/manifest.md`. These `docs/features/` paths are used extensively across 20+ locations in gen-contracts, gen-journeys, gen-test-scripts, ui-design, and other skills. The proposal's fix for H-2 says "需先搜索确认 proposal 文件是否可能存在于该路径" but does not address the broader question: is `docs/features/` the canonical output path for the write-prd/tech-design pipeline, or is `docs/proposals/` the canonical path? The `docs/features/` path appears to be the standard output directory for the feature pipeline (manifest.md, prd/, design/, testing/ all live there), while `docs/proposals/` is the brainstorm output directory. The proposal correctly identifies that no skill creates `docs/features/<slug>/proposal.md`, but the broader path topology question deserves a clearer analysis.

### Eval Journey Score Threshold Calculation

问题：The proposal's urgency section states: "H-1 意味着当用户参照 argument-hint 手动传入 --target 850 时，eval-journey 会以 850/1150=73.9% 作为通过标准，而非正确的 975/1150=84.8%，通过门槛被降低了 11 个百分点". This analysis is correct for the manual `--target` case. However, the proposal should also identify which downstream consumers read `rubric-reference.md` directly. The eval command's Config Resolution logic (eval-journey.md lines 26-28) shows that when no `--target` is passed and no config value exists, the argument is omitted entirely and "eval skill uses rubric default" (line 19). The eval SKILL.md confirms: "CLI --target/--iterations override frontmatter". So the default path correctly reads target=975 from rubric frontmatter. But the `rubric-reference.md` showing target=850 means any process or LLM context that reads rubric-reference.md instead of rubric frontmatter would use the wrong value. The proposal should enumerate these consumers explicitly.

---

## Section 3: Improvement Suggestions

建议：Expand the H-1 fix scope to include the `eval/SKILL.md` description field. The current text reads "Supports 100-point and 1000-point scales" and should be updated to reflect the actual scale range used by rubrics in the system. A more accurate description would be "Supports configurable scales defined per rubric (e.g., 1000, 1100, 1150)". This ensures the eval skill's self-description does not mislead LLM scorer subagents about valid score ceilings.

建议：Expand the H-1 fix to also update the dimension descriptions in `eval-journey.md` and `eval-contract.md`. The journey command description should include "workflow coverage" as the 7th dimension. The contract command description should be updated from "six-dimension structural integrity" to reflect the actual 8 dimensions including anchor integrity and fixture specification. These descriptions are read by the LLM during eval execution and directly influence scoring completeness.

建议：Update M-9 to acknowledge all 4 INLINE references, not just 3. The bidirectional dependency between gen-contracts and gen-journeys (gen-contracts inlines Surface Detection from gen-journeys; gen-journeys inlines journey-contract-model from gen-contracts) should be explicitly flagged as a mutual drift risk. The fix should add version stamps to all 4 locations, and the regression verification grep command should be updated to check for 4 matches instead of 3.

建议：Correct the "22 个 skill" count in the Problem section to 21. This is a factual error that should be fixed before the proposal is finalized. Also update the dimension C scope from "约40 模板文件" to a more accurate count (56 files) to ensure the audit scope claim is defensible.

建议：The H-1 fix description should explicitly enumerate all 5 truth-source locations that the proposal identifies and state which ones are being fixed and which are deferred. Currently the fix says "更新 rubric-reference.md 表格 + eval-journey/eval-contract 的 argument-hint 和 description 字段" — this covers 3 of the 5 locations. The config key defaults and eval SKILL.md description should be addressed or explicitly acknowledged as known gaps.

建议：For the M-9 version stamp approach, consider adding a structured comment format that can be machine-validated. Instead of `<!-- INLINE from ... @ v3.0.0-rc.53 -->`, use a format that includes a checksum or line count of the source content at the time of stamping (e.g., `<!-- INLINE:origin=... @ v3.0.0-rc.53 lines=104 -->`). This enables automated drift detection by comparing the current line count of the source against the stamped value, without requiring version string management.

建议：The proposal should add a regression verification step specifically for the eval SKILL.md description. After all fixes, `grep "100-point" plugins/forge/skills/eval/SKILL.md` should return no results, and the description should accurately reflect the supported scale range. Similarly, `grep "six-dimension" plugins/forge/commands/eval-contract.md` should be added to the regression verification to catch stale dimension count references.

建议：Clarify the M-1 execution condition in the Success Criteria. The current wording creates ambiguity about whether M-1 is committed or conditional. Either commit to M-1 with the precondition that the Go config reader must be verified first (and include the verification step in the fix order), or explicitly defer M-1 to a follow-up proposal. The current formulation where M-1 appears in the fix list but may be skipped creates confusion about the proposal's actual scope.
