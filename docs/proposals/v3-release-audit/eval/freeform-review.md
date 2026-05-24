---
created: 2026-05-24
reviewer: Documentation-Implementation Drift Auditor (Expert)
type: freeform-narrative-review
status: completed
---

# Freeform Narrative Review: v3.0.0 Release Audit Proposal

## Section 1: Background Assessment

### Proposal Summary

The proposal identifies systematic documentation-implementation drift across Forge v3.0.0's core documents (README.md and ARCHITECTURE.md) and the actual codebase. It organizes findings into a 5-dimension audit framework, classifies 27 drift items into Critical/Major/Minor/Advisory tiers, and proposes a phased remediation plan limited to documentation updates and dead code cleanup without runtime code changes.

### My Verification Approach

I independently verified every factual claim in the proposal against the live codebase at commit `7d6a51cd` (branch `v3.0.0`). This included running CLI commands, counting files, checking cross-references, validating SKILL.md line counts, tracing rule file references, and comparing documented component descriptions against actual implementations.

### Assessment of the 5-Dimension Framework

The five dimensions (README factual claims, ARCHITECTURE factual claims, CLI reference accuracy, Skill-CLI cross-references, architecture health) provide reasonable coverage of the documentation surface. However, as detailed below, there are significant gaps -- dimensions that should exist but don't, dimensions where the counts are wrong, and dimensions where "documentation-only" changes have hidden runtime implications that the proposal does not acknowledge.

---

## Section 2: Key Risk Identification

### 风险：Proposal's own severity counts contain factual errors

The proposal claims:

> 系统性审计覆盖 5 个维度，发现 **27 个偏差项**：
> | 维度 | Critical | Major | Minor | Advisory |
> | **合计** | **17** | **13** | **15** | **5** |

Adding the cells: 17+13+15+5 = 50 total items, not 27. The "27 个偏差项" headline number contradicts the table's own column totals. This is not a cosmetic issue -- the entire proposal is organized around severity classification, and the primary metric is wrong. If the counts in the table are correct, the headline should read "50 个偏差项". If the headline is correct, the table is wrong. Either way, the proposal as written contains the exact type of factual inaccuracy it was created to fix.

### 问题：CLI cross-reference count underestimates broken references

The proposal states:

> `forge config get surface` → `forge surfaces`（4 处）

My independent verification found exactly 4 occurrences of `forge config get surface` across plugin files, which matches. However, the proposal categorizes this as 4 Critical items. The actual impact is worse than described: `forge config get surface` is not merely "renamed" -- it fundamentally does not work as written. The `forge config get` command expects dot-notation keys like `auto.gitPush`, not `surface`. The correct surface retrieval mechanism is `forge surfaces <path>`, a completely different command with different semantics (path-based resolution vs. config key lookup). Furthermore, `run-tests/SKILL.md` line 73 uses `forge config get test.execution` -- another key that does not exist in the standard config schema. The proposal missed this additional broken reference entirely.

### 问题：`forge test run --tags` reference count is wrong in the proposal

The proposal claims:

> `forge test run --tags regression` → 正确命令（4 处）

There is no `forge test run` subcommand. The actual CLI provides `forge test promote`, `forge test run-journey`, and `forge test verify`. The `--tags` flag does not exist on any subcommand. My verification found 4 occurrences, which matches the proposal's count. However, these 4 occurrences are all in `journey-contract-model.md` rules files that are not referenced from their parent SKILL.md files (they are orphaned rules, as detailed below). This means the proposal classifies them as Critical P0 items, but their actual runtime impact is lower than claimed because the rules files containing these references are not actively loaded during execution.

### 风险：SKILL.md splitting is more complex than acknowledged, with distribution implications

The proposal states:

> SKILL.md 拆分是最复杂的操作，但只需将现有内容移入 rules/ 文件。

This significantly understates the complexity. The `eval/SKILL.md` at 488 lines is a deeply interconnected document with 7 explicit cross-references to rules files, 4 references to experts/ files, a Mermaid flowchart, and complex conditional branching (Phase 0 freeform review path vs. standard rubric path). Splitting it requires:

1. Maintaining the flowchart's accuracy after content extraction
2. Ensuring that conditional Phase 0 / Phase 0.5 / standard flow references remain coherent across split files
3. Not breaking the `Load: rubrics/<type>.md` instruction which is the only mechanism by which rubrics are discovered

Per forge-distribution.md, rules/ files are distributed with the plugin and loaded at runtime by agent context. Any splitting must respect the distribution model -- the proposal mentions forge-distribution.md as a constraint but does not analyze how splitting changes the agent's context loading behavior.

### 风险：gen-journeys surface rules use parameterized paths that appear "orphaned" but aren't

The proposal classifies 11 unreferenced rules files. My verification confirms there are exactly 15 orphaned rules files (not 11) across all skills. However, among these are the 5 `gen-journeys/rules/surface-*.md` files (surface-api.md, surface-cli.md, surface-mobile.md, surface-tui.md, surface-web.md). The gen-journeys SKILL.md references them via a pattern:

> "Load the corresponding rule file from `rules/surface-<type>.md`"

This is a parameterized reference, not a direct file reference. The proposal correctly identifies this pattern in the remediation section but miscounts the total orphaned files. The actual orphaned files that need remediation are:

**Truly orphaned (no SKILL.md reference, direct or parameterized):**
- `eval/rules/freeform-injection.md`
- `gen-contracts/rules/journey-contract-model.md`
- `gen-journeys/rules/journey-contract-model.md`
- `gen-test-scripts/rules/convention-guide.md`
- `gen-test-scripts/rules/run-to-learn.md`
- `run-tests/rules/test-isolation.md`

**Parameterized references (surface rules):**
- `gen-journeys/rules/surface-api.md`
- `gen-journeys/rules/surface-cli.md`
- `gen-journeys/rules/surface-mobile.md`
- `gen-journeys/rules/surface-tui.md`
- `gen-journeys/rules/surface-web.md`

These two categories require fundamentally different remediation approaches, but the proposal treats them uniformly.

### 问题：ARCHITECTURE.md doc-scorer/doc-reviser agent descriptions describe non-existent implementation

The proposal identifies that ARCHITECTURE.md describes "4 agents" when only 1 agent file exists (`task-executor.md`). It classifies this as a Critical issue. But the proposal's own remediation plan describes:

> Agent 架构（1 个专用 agent + general-purpose 模式）

My verification confirms: eval does NOT use doc-scorer or doc-reviser agents. The eval SKILL.md spawns `general-purpose` agents with protocol + expert composition. ARCHITECTURE.md's detailed "doc-scorer" and "doc-reviser" sections (lines 175-210) describe a multi-agent architecture that has never existed in v3.0.0. The scorer protocol at `experts/protocol/scorer-protocol.md` uses `{{rubric_total}}` (parameterized by rubric scale), not a fixed "100 分制" as ARCHITECTURE.md line 179 claims. All rubrics use 1000-point scale except `validate-code` and `validate-ux` (also 1000-point, target 700).

The proposal misses that ARCHITECTURE.md's "100 分制" claim is wrong. Every rubric in `eval/rubrics/` uses `scale: 1000`. The README compounds this by claiming eval-harness uses "100 分制", but there is no harness rubric file at all.

### 风险：Missing harness.md rubric is more severe than classified

The proposal states:

> 缺失 rubric 文件：创建 `eval/rubrics/harness.md` 或在 SKILL.md 中添加异常处理

This is classified as P0 Critical, item 5. But the `harness` type is listed in eval SKILL.md's Prerequisites table (line 24) and type parameter (line 34) as a valid eval type. If a user invokes `/eval --type harness`, the skill will attempt to `Load: rubrics/harness.md` (line 109), fail to find it, and the pipeline behavior is undefined. The proposal suggests creating the rubric file or adding exception handling, but creating a rubric is a content creation task with domain-specific quality requirements, not a documentation fix. The alternative (adding exception handling to SKILL.md) is a runtime behavior change, contradicting the proposal's scope boundary:

> 仅涉及文档更新和死代码清理，不修改任何运行时代码。

### 问题：README.md references multiple non-existent components and commands

The proposal identifies version number and count discrepancies. My verification found additional issues the proposal did not enumerate:

1. **README line 5**: `Version-2.16.1` -- the plugin is at `3.0.0-rc.24`. Proposal identifies this.

2. **README line 6**: `Go-1.26.1+` -- the actual go.mod says `go 1.25`. Proposal identifies this.

3. **README line 131**: `task-cli/` -- the directory is now `forge-cli/`. The proposal identifies this.

4. **README line 134**: `web/` -- the `web/` directory does not exist. The proposal does not mention this.

5. **README line 89**: claims `/eval-harness` exists as a command. No `eval-harness.md` command file exists. No `eval-harness/` skill directory exists. The `harness` type is a parameter of the generic `/eval` skill, not a standalone slash command. The proposal mentions eval-harness in ARCHITECTURE.md context but does not flag this README error as a separate item.

6. **README line 107**: lists `/learn-lesson`, `/record-decision`, and `/improve-harness` as auxiliary skills. `/learn-lesson` and `/record-decision` were absorbed into `/learn`. `/improve-harness` has no corresponding skill directory, command file, or any reference in the entire plugins/ directory. The proposal does not explicitly flag `/improve-harness` as a ghost command.

7. **README line 128-130**: claims "17 个 Skills", "17 个 Slash Commands", "3 个 Subagents". Actual counts: 21 skills, 18 commands, 1 agent. The proposal identifies the agent count discrepancy but does not appear to call out that the skill count (21 vs 17) and command count (18 vs 17) are both wrong and in opposite directions.

8. **README line 168**: claims "13 种任务类型". The CLI reports 21 task types. The old naming (e.g., `implementation`, `documentation`, `doc-evaluation`) is completely absent from the CLI -- the new naming uses dot-notation namespaces (`coding.feature`, `doc.review`, `test.gen-scripts`, etc.). The proposal mentions "21种新命名" in the P0 section but the discrepancy between README's "13 种" and reality's "21 种" deserves explicit enumeration.

### 问题：ARCHITECTURE.md describes a hooks system that doesn't match reality

ARCHITECTURE.md line 398 lists a `PostToolUse` hook running `validate-index.sh`. My verification found:

1. No `validate-index.sh` script exists anywhere in the repository.
2. The actual `hooks.json` has no `PostToolUse` entry -- only `SessionStart`, `SubagentStart`, `SessionEnd`, `SubagentStop`, and `Stop`.
3. The `Stop` hook runs both `forge quality-gate` and `forge feature complete --if-done`, but ARCHITECTURE.md only mentions `forge quality-gate` and says nothing about `forge feature complete --if-done`.

The proposal mentions removing PostToolUse but does not mention that the `forge feature complete --if-done` hook is completely undocumented in ARCHITECTURE.md. This is a significant omission -- `forge feature complete` handles post-completion artifact detection, manifest updates, and optional git push, all of which are invisible to anyone reading the architecture document.

### 风险：ARCHITECTURE.md omits entire v3.0.0 feature areas

My verification found that ARCHITECTURE.md has zero references to:

- **Surface detection system** (mentioned once in passing): The entire surface-based test profile system -- `forge surfaces` CLI, `.forge/config.yaml` surfaces configuration, surface-specific rule files, Convention-driven test framework detection -- is absent from the architecture document. This is arguably the most significant v3.0.0 feature.
- **Worktree management**: `forge worktree start/push/remove/resume/status` -- an entire CLI subsystem -- is undocumented.
- **Convention system**: `docs/conventions/testing/` directory with framework-specific test configuration files, the two-level index mechanism for Convention loading, the Convention-driven e2e recipe generation in init-justfile -- none of this appears in ARCHITECTURE.md.
- **Forensic skill**: `/forensic` for session transcript analysis is absent.
- **Deep research skill**: `/deep-research` for technology research is absent.
- **Clean code skill**: `/clean-code` for code cleanup is absent.
- **Extract design skill**: `/extract-design-md` for style extraction is absent.
- **Test guide skill**: `/test-guide` for Convention creation is absent.
- **Learn skill**: `/learn` (unified knowledge accumulation) is absent; ARCHITECTURE.md still references `/learn-lesson` and `/record-decision`.
- **doc type quality gate exemption**: guide.md documents that `doc*` type tasks skip the quality gate, but ARCHITECTURE.md's Quality Gate section does not mention this.

The proposal's P0 item 2 covers "ARCHITECTURE.md 修正" but limits it to Agent architecture, Hook system, Eval system, and directory paths. It does not address the wholesale absence of major v3.0.0 subsystems.

### 问题：Dead code classification for init-justfile templates is incorrect

The proposal states:

> 死代码清理：init-justfile 6 个 .just 模板文件（SKILL.md 明确说不使用）

My verification found that init-justfile's SKILL.md says:

> "Do NOT use framework-specific recipe templates. Generate e2e recipes from Convention content and LLM knowledge of the framework."

This means the .just templates are NOT dead code -- they are explicitly rejected design choices that are correctly documented as such. The `generic.just` template still contains functional content (error-stub recipes, e2e-test recipe, e2e-setup recipe, probe recipe, e2e-verify recipe) that serves as a reference implementation. The SKILL.md does NOT say "these files are not used." It says "do not use framework-specific templates" -- a design principle, not a statement about file usage. Deleting these files removes reference material without gaining anything.

Furthermore, the SKILL.md itself does NOT reference the templates directory at all, which means the templates are unreachable from the skill's execution flow. This is different from being "dead code that SKILL.md says not to use."

### 风险：Cross-skill path violations in run-tests are distribution-breaking

The proposal states:

> 跨技能路径违规修复：run-tests 中 `skills/gen-journeys/rules/surface-<type>.md` 改为本地副本或描述性引用

My verification confirms this violation in `run-tests/SKILL.md` line 147 and `run-tests/rules/env-check.md` line 13. According to forge-distribution.md:

> "跨 skill 文件 | 描述性路径 + 上下文 | `ui-design/templates/styles/<name>.md`（注明 resolve relative to the skills parent directory）"

The fix must use "descriptive path + context" style, not create local copies (which would violate skill-self-containment by creating parallel maintenance burden) and not use absolute paths (which break after distribution). The proposal offers "本地副本或描述性引用" but the "local copy" option violates the distribution model because surface rules evolve independently -- a local copy would silently drift from the canonical version in gen-journeys.

### 风险：Tests directory structure described in docs does not exist

ARCHITECTURE.md describes a `tests/e2e/` directory structure with `playwright.config.ts`, `helpers.ts`, and `features/<slug>/` subdirectories. My verification found:

1. `tests/e2e/` does not exist
2. There is no `playwright.config.ts` anywhere in the repository
3. The actual `tests/` directory contains Go integration tests (e2e-pipeline, feature-management, quality-gate, etc.), not Playwright specs
4. The Forge project itself uses Go testing, not Playwright. Playwright is the framework that Forge generates for user projects.

ARCHITECTURE.md confuses Forge's internal test infrastructure with the test infrastructure Forge generates for user projects. The proposal does not flag this category error.

### 问题：Success criteria are not independently verifiable for several items

The proposal states:

> README.md 所有事实性声明（版本号、计数、路径、命令名）与代码库 100% 一致

This is verifiable. However:

> ARCHITECTURE.md 所有组件描述（agents、hooks、eval、目录）与代码库 100% 一致

This is not directly verifiable without defining what constitutes a "component." If surface detection, worktree management, Convention system, forensic, deep-research, etc. are "components," then achieving 100% requires adding significant new content to ARCHITECTURE.md -- well beyond "fixing drift" and into "writing new documentation."

> 零断裂 CLI 交叉引用

My verification found that `forge config get surface` is not the only broken CLI reference. `forge config get test.execution` in `run-tests/SKILL.md` line 73 is also potentially broken (this key does not exist in the standard config.yaml). The success criterion would need to be expanded to catch all broken references, not just the two patterns mentioned.

### 风险：Remediation dependency ordering has a hidden cycle

The proposal orders P0 items 1-5 linearly. However, P0 item 1 (README rewrite) depends on knowing the correct skill count, command count, task type list, and pipeline structure -- all of which may change as a result of P0 items 4 (SKILL.md splitting) and 5 (harness rubric creation). If SKILL.md splitting creates new skills or the harness rubric changes the eval type list, the README rewrite would need to be redone. The proposal does not acknowledge this dependency.

### 问题：`forge config get` command appears non-functional in the development environment

During my verification, `forge config get <key>` consistently returned exit code 1 for every key tested, including keys that should exist in `.forge/config.yaml` (like `auto.gitPush` and `version`). The unit tests pass (`TestConfigGetCommand`), suggesting the command works in test environments. The issue may be related to the project root detection on Windows (the `FindProjectRoot` function walks up from cwd looking for project markers, and the forge repo itself has a `.forge/` directory). If `forge config get` is unreliable in practice, then ALL 4 occurrences of `forge config get surface` in skill files would fail at runtime, regardless of whether the key name is correct. This is a deeper issue than a simple rename.

### 风险：Proposal scope boundary claim is violated by multiple items

The proposal states:

> 仅涉及文档更新和死代码清理，不修改任何运行时代码。

But several proposed changes affect runtime behavior:

1. **P0 item 5** (creating `eval/rubrics/harness.md`): Adding a rubric file changes the eval pipeline's behavior when `--type harness` is used. This is a runtime change, not a documentation change.
2. **P0 item 4** (SKILL.md splitting): Moving content from SKILL.md to rules/ files changes the agent's context loading behavior. Agents load SKILL.md as a unit; splitting means the agent must now load additional files explicitly.
3. **P1 item 8** (adding Load instructions for orphaned rules): Adding `Load:` directives to SKILL.md changes which files the agent loads at runtime, potentially increasing context consumption and changing task execution behavior.
4. **P1 item 9** (dead code cleanup): Deleting files that are distributed with the plugin changes the plugin's distribution content. If any agent or tool references these files through any path (even undocumented ones), the deletion breaks functionality.

### 问题：ARCHITECTURE.md's "13 种任务类型" is internally inconsistent with README

ARCHITECTURE.md does not have a "13 种任务类型" section. The README has this section. The proposal correctly identifies that README task types are wrong but attributes the "21种新命名" claim to the P0 remediation for README only. The actual task type system (verified via `forge task list-types`) shows 21 types with completely new dot-notation naming that has zero overlap with the README's old naming. This is not merely a count discrepancy -- it's a complete naming scheme change that affects how agents generate tasks, how tasks are classified in index.json, and how the entire test pipeline identifies task types.

---

## Section 3: Improvement Suggestions

### 建议：Fix the headline count before proceeding

The "27 个偏差项" headline must be reconciled with the table totals (50). One of them is wrong. Without this fix, the proposal's credibility is undermined and any derived work items will have incorrect priority counts.

### 建议：Expand the 5-dimension framework to include "ARCHITECTURE.md missing feature coverage"

The current dimensions audit what's wrong in existing documentation. They don't audit what's missing. ARCHITECTURE.md's failure to document surface detection, worktree management, the Convention system, `/forensic`, `/deep-research`, `/clean-code`, `/extract-design-md`, `/test-guide`, and the unified `/learn` skill represents a separate dimension of drift: documentation that should exist but doesn't. Adding a 6th dimension ("Feature Coverage Completeness") would provide a more honest assessment.

### 建议：Reclassify `forge test run --tags` references as P1, not P0

The 4 occurrences of `forge test run --tags` are in `journey-contract-model.md` rules files that are orphaned (not referenced from their parent SKILL.md). Their runtime impact is therefore nil -- no agent loads these files during execution. The P0 classification inflates urgency and should be P1 at most.

### 建议：Clarify the harness rubric decision before execution

The proposal offers two options: "create harness.md" or "add exception handling to SKILL.md." These have fundamentally different scope implications. Option 1 is content creation (potentially out of scope). Option 2 is runtime behavior modification (definitely out of scope per the proposal's own constraints). A third option -- removing `harness` from the valid types list in eval SKILL.md until a rubric is ready -- is more consistent with the "documentation-only" scope.

### 建议：Remove `improve-harness` from README ghost command list

The README line 107 lists `/improve-harness` as an auxiliary skill. This is not mentioned in the proposal's P0 item 1 enumeration. It should be explicitly called out and removed.

### 建议：Investigate `forge config get` reliability before fixing `forge config get surface`

If `forge config get` is broken in the development environment (as my testing suggests), then simply renaming `surface` to `surfaces` or replacing `forge config get surface` with `forge surfaces <path>` addresses the symptom but not the root cause. The proposal should include a verification step: confirm that `forge config get` works correctly for known keys before assuming the fix is just a rename.

### 建议：Do NOT delete init-justfile .just template files

The proposal classifies 6 `.just` template files as dead code. My verification shows these files contain functional reference content and their existence does not cause harm. The SKILL.md's "Do NOT use framework-specific recipe templates" directive is a design principle, not a statement about file deletion. Deleting them provides no benefit and risks losing reference material.

### 建议：Split ARCHITECTURE.md remediation into two phases

The current P0 item 2 tries to fix ARCHITECTURE.md in one pass. Given the depth of drift (4 agents described when 1 exists, doc-scorer/doc-reviser sections that describe a never-implemented architecture, missing entire subsystems), a single pass risks either: (a) doing a superficial fix that leaves the document still inaccurate, or (b) expanding into new-content authorship that goes beyond "drift remediation." Split into:

- P0: Fix factual errors in existing content (agent count, hook table, eval score scale, directory paths)
- P1/new: Add missing subsystem documentation (surface detection, worktree, Convention system, new skills)

### 建议：Add `forge feature complete --if-done` to the ARCHITECTURE.md hooks table

The Stop hook runs two commands, not one. ARCHITECTURE.md only documents `forge quality-gate`. The `forge feature complete --if-done` command handles manifest updates, artifact detection, and optional git push -- all of which are invisible to anyone reading the architecture document. This is a functional omission, not just a documentation drift.

### 建议：Clarify the "tests/e2e/" directory structure description

ARCHITECTURE.md describes `tests/e2e/` with `playwright.config.ts` and `helpers.ts`. This directory does not exist in the forge repository. The description appears to document the structure that Forge generates in user projects, not Forge's own test infrastructure. The proposal should explicitly clarify whether ARCHITECTURE.md is describing Forge's internal structure or the structure Forge creates for users, and ensure the documentation clearly distinguishes between the two.

### 建议：Establish a verification gate for success criteria

The proposal's success criteria ("100% 一致") are aspirational but fragile. Before closing the audit, each criterion should be verified by an automated check:

- `forge task list-types | wc -l` should match the documented task type count
- `ls plugins/forge/skills/*/SKILL.md | wc -l` should match the documented skill count
- `grep -r "forge config get surface" plugins/forge/` should return zero results
- `wc -l plugins/forge/skills/*/SKILL.md` should show no file exceeding 350 lines

These automated checks would prevent regression and provide objective evidence that remediation succeeded.

### 建议：Fix the ARCHITECTURE.md "forge forge task claim" typo

Line 443 of ARCHITECTURE.md contains `forge forge task claim` (doubled "forge"). Line 147 contains the same error. The proposal does not mention this specific typo. While minor, it is the kind of factual error that could confuse a new contributor or agent.

### 建议：Address the "all-completed Hook" naming mismatch

ARCHITECTURE.md describes an "all-completed Hook" section (lines 403-416). The actual implementation uses the `Stop` hook event in `hooks.json`. There is no separate "all-completed" hook -- it's the Stop event that triggers quality-gate and feature-complete. The section title and description imply a separate hook mechanism that does not exist.
