# Report 03: Skills Deep Audit — Batch B (Layer 2-3)

**Baseline commit**: `7fd5aab0c2c2d7e6f0c8a67e11afefbf9edc3d45`
**Date**: 2026-05-30
**Scope**: 7 skills (quick-tasks, consolidate-specs, clean-code, deep-research, forensic, ui-design, learn)
**Audit order**: quick-tasks → consolidate-specs → clean-code → deep-research → forensic → ui-design → learn (randomized)
**Method**: Layer 2 (Instruction Consistency) + Layer 3 (Timing & Flow)

---

## 1. quick-tasks

### Structured Summary (SKILL.md)

- **Steps**: 0 (Resolve Language) → 1 (Read Proposal) → 2 (Derive Tasks) → 3 (Create Task Files) → 4 (Test Tasks, auto-generated) → 5 (Generate index.json) → 6 (Validate) → 7 (Create Manifest) → 8 (Commit Planning Artifacts)
- **Key constraints**: HARD-GATE (max 6 AC per task); HARD-RULE (no silent default language, naming conventions, stage only planning artifacts, `.md` files are non-compilable regardless of directory)
- **References**: `templates/task.md`, `templates/task-doc.md`, `templates/manifest-quick.md`
- **Docs-Only Fast Path**: skip Step 0 and Step 4 when all tasks are type: "doc"
- **Intent Propagation**: `proposal.md` frontmatter `intent` → per-task type override
- **Field names**: ID, TITLE, PRIORITY, ESTIMATED_TIME, DEPENDENCIES, SLUG, SURFACE_KEY, SURFACE_TYPE, COMPLEXITY, breaking

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| QT-01 | SKILL.md Step 0 vs `templates/task.md` L9 | 2 | INCOMPLETE | P2 | SKILL.md Step 0 says "Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.)." This is the same `domains` filtering pattern flagged in batch A (GTS-01, BT-01, TD-05). quick-tasks uses Convention loading for language detection, but the `domains` approach conflicts with gen-test-scripts' `rules/convention-guide.md` HARD-RULE against domains filtering. | Same as batch A finding GTS-01: verify whether quick-tasks should use domains filtering or the index.md-based approach. If different, document the intentional difference. | medium |
| QT-02 | `templates/task.md` L7 vs SKILL.md Step 2 | 2 | CONFLICT | P3 | `templates/task.md` hardcodes `complexity: "medium"` in frontmatter. SKILL.md Step 2 defines a complexity heuristic with three levels (low/medium/high) and LLM override. The template always produces "medium" regardless of the derived complexity. The agent is instructed to "Set the `complexity` field in task frontmatter accordingly" but the template placeholder is not `{{COMPLEXITY}}` — it is a hardcoded string. | Replace the hardcoded `complexity: "medium"` with `complexity: "{{COMPLEXITY}}"` in the template to align with SKILL.md's complexity heuristic. | high |
| QT-03 | `templates/task.md` L11 vs SKILL.md Step 3 Type Assignment | 2 | CONFLICT | P3 | `templates/task.md` hardcodes `type: "coding.feature"` in frontmatter. SKILL.md Step 3 Type Assignment defines 8 possible types (coding.feature, coding.enhancement, coding.cleanup, coding.refactor, coding.fix, doc, doc.consolidate, doc.drift) with a fallback to `coding.feature`. The template always produces "coding.feature" regardless of the actual type. Additionally, `templates/task-doc.md` correctly hardcodes `type: "doc"`, but the coding template lacks a `{{TYPE}}` placeholder. | Add a `{{TYPE}}` placeholder to `templates/task.md` and instruct the agent to set it per the Type Assignment table. Default to "coding.feature" if not overridden. | high |
| QT-04 | `templates/task.md` L9 vs SKILL.md Step 3 | 2 | INCOMPLETE | P3 | `templates/task.md` includes `surface-key` and `surface-type` fields. SKILL.md Step 2 Surface-Key/Type Inference says these apply to "coding tasks only" and the template placeholders reference `{{SURFACE_KEY}}` and `{{SURFACE_TYPE}}`. However, SKILL.md Step 3 "Task Template Placeholders" table lists these placeholders as "coding tasks only" — consistent with the template. No issue here, but the template's default values are empty strings, which could cause validation issues if `forge task index` expects non-empty surface-key for coding tasks. | No action needed — the template correctly includes these fields for coding tasks only. Verify that `forge task index` handles empty surface-key gracefully. | low |
| QT-05 | SKILL.md Step 3 "Breaking Task Test Impact Assessment" vs `templates/task.md` L10 | 2 | INCOMPLETE | P2 | SKILL.md Step 3 "Breaking Task Test Impact Assessment" section says "When setting `breaking: true` on a task, the task description MUST include a test impact assessment." The template has `breaking: false` hardcoded (no `{{BREAKING}}` placeholder). Additionally, the template lacks placeholders for test impact fields (Affected test suite(s), Expected fixture changes, Risk level). The agent must add these manually to `## Implementation Notes`, but the template does not prompt for them. | Add a conditional section to `templates/task.md` or add an instruction note that when `breaking: true`, the Implementation Notes must include the Test Impact subsection. | medium |
| QT-06 | `templates/manifest-quick.md` L6 vs SKILL.md Step 7 | 2 | INCOMPLETE | P3 | SKILL.md Step 7 says "Write to `docs/features/<slug>/manifest.md`" and lists placeholder `{{TASK_ROWS}}` with format `\| <ID> \| <title> \| pending \| <ID>-<slug>.md \|`. The template at `templates/manifest-quick.md` has `{{TASK_ROWS}}` but the manifest template does not include the `mode: quick` field that SKILL.md Output Checklist expects ("manifest.md written with mode: quick"). The template does include `mode: quick` in frontmatter (line 6), so this is consistent. | No action needed — verified consistent. | low |

### Layer 3 Findings

No TIMING issues found. The step sequence (0 → 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8) has correct dependency ordering. The Docs-Only Fast Path correctly skips steps that are unnecessary for doc-only features.

---

## 2. consolidate-specs

### Structured Summary (SKILL.md)

- **Steps**: 1 (Check Idempotency) → 2 (Read Feature Documents) → 3 (Extract biz-specs) → 4 (Extract tech-specs) → 5 (Generate Preview + Detect Overlaps) → 6 (Present to User) → 7 (Integrate Approved Specs) → 8 (Write Integration Marker + Update Manifest) → 9 (Detect Drift) → 10 (Auto-Fix Drifted Specs) → 11 (Commit Spec Changes) → 12 (Generate Vocabulary Index) → 13 (Record Task)
- **Key constraints**: HARD-GATE (no integration without user confirmation, no overwrite, no re-run if marker exists, no rule inference)
- **References**: `templates/biz-specs.md`, `templates/tech-specs.md`, `templates/markers.md`, `templates/review-choices.md`, `templates/commit-messages.md`, `templates/vocabulary-index.md`, `rules/constraints.md`, `rules/domain-frontmatter.md`, `rules/drift-detection.md`, `rules/overlap-detection.md`, `rules/spec-classification.md`, `rules/vocabulary-generation.md`
- **Skip Conditions**: no extractable rules, all items LOCAL, non-interactive session
- **Drift-Only Mode**: skip Steps 1-8 when no PRD/design exists
- **Field names**: CROSS/LOCAL classification, project-global IDs (BIZ-<domain>-NNN, TECH-<topic>-NNN), domains frontmatter

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| CS-01 | SKILL.md Step 6 vs `rules/constraints.md` L13-14 | 2 | REDUNDANT | P3 | SKILL.md Step 6 "Non-Interactive Mode" and `rules/constraints.md` L11-14 both describe non-interactive mode behavior in detail. SKILL.md says "Auto-integrate all [CROSS] items without blocking... Auto-write review-choices.md... Commit with [auto-specs] tag." constraints.md says "Non-interactive mode auto-integrates all [CROSS] items without blocking -- commit includes [auto-specs] tag... [auto-specs] commits must be separate from code change commits." The constraints file adds the "separate commit" rule that SKILL.md Step 6 does not mention. This is additive detail, not pure redundancy, but the overlap in non-interactive mode description could be consolidated. | Consider having SKILL.md Step 6 reference constraints.md for non-interactive mode details instead of inlining them. | low |
| CS-02 | SKILL.md Step 9 vs `rules/drift-detection.md` | 2 | REDUNDANT | P2 | SKILL.md Step 9 says "validate each rule against the current codebase, and classify as `current`, `drifted`, or `orphaned`" with a brief description. `rules/drift-detection.md` repeats the same classification taxonomy (`current`, `drifted`, `orphaned`) and expands with the full procedure. The classification taxonomy is defined identically in both places. While the rules file provides legitimate expansion (search methodology, output format), the classification names should be defined once. | Define the classification taxonomy (current/drifted/orphaned) only in drift-detection.md and have SKILL.md Step 9 reference it without repeating the names. | low |
| CS-03 | SKILL.md Step 5 "Early exit" vs `templates/markers.md` | 2 | CONFLICT | P3 | SKILL.md Step 5 "Early exit" says "Write `docs/features/<slug>/specs/.integrated` for the early-exit case using the early-exit marker template from `templates/markers.md`." The markers.md template for "Early-Exit Marker (all LOCAL)" includes `status: "skipped: all local"`. However, SKILL.md Step 1 says "If the marker exists, this feature's specs have already been integrated. Read the marker to confirm, then skip with status `completed`." The early-exit marker's `status: "skipped: all local"` would trigger the idempotency check to skip, which is the correct behavior — a re-run after early exit should not re-extract. But the semantic difference between "integrated" (standard marker) and "skipped: all local" (early-exit marker) is not documented in Step 1's idempotency check logic. Step 1 says "read the marker to confirm" but does not specify what happens if the marker says "skipped: all local" — should it skip or re-run extraction? | Add logic to SKILL.md Step 1: if the marker's status is "skipped: all local", check whether the feature documents have changed since the marker was written. If unchanged, skip; if changed (new rules extractable), re-run. Or clarify that "skipped: all local" is treated the same as "integrated" for idempotency. | medium |
| CS-04 | `rules/spec-classification.md` "Preview ID Numbering" vs `templates/biz-specs.md`/`templates/tech-specs.md` | 2 | INCOMPLETE | P3 | `rules/spec-classification.md` says "Feature-local IDs in preview files use sequential 3-digit numbering starting at 001" with format `BIZ-001`, `TECH-001`. The templates (`biz-specs.md`, `tech-specs.md`) use `BIZ-NNN` and `TECH-NNN` as placeholders, which is consistent. However, the templates do not explicitly state that numbering starts at 001 or is 3-digit — the "NNN" convention is only documented in spec-classification.md. | No action needed — the convention is documented in the rules file, which is the appropriate location. Flagged for awareness. | low |
| CS-05 | `templates/vocabulary-index.md` output path vs SKILL.md Step 12 | 2 | INCOMPLETE | P3 | `templates/vocabulary-index.md` says "Write the vocabulary index to `docs/.vocabulary.md`". SKILL.md Step 12 says "Scan all four knowledge directories and produce a vocabulary index for use by `/learn` and auto-extract triggers. This step runs unconditionally. See `rules/vocabulary-generation.md` for scan targets, base vocabulary, aggregation rules, output format, and usage by other skills." SKILL.md does not mention the output path `docs/.vocabulary.md` explicitly — it defers to vocabulary-generation.md, which says "Use the output template from `templates/vocabulary-index.md`" and the template specifies the path. The path is discoverable but not in SKILL.md. | Add the output path `docs/.vocabulary.md` to SKILL.md Step 12 for direct discoverability. | low |
| CS-06 | `rules/overlap-detection.md` L6 "Domain-to-decision-file mapping" vs `rules/spec-classification.md` | 2 | INCOMPLETE | P3 | overlap-detection.md defines a mapping from "spec domain keywords" to "decision file" using the 8-category vocabulary. The mapping table includes entries like "naming, conventions, coding standards → architecture.md" and "validation, state transitions, calculation rules → closest match or architecture.md". However, spec-classification.md does not reference this mapping, and the two files use different terminology for the same concept (spec-classification says "target filename" while overlap-detection says "decision file"). The mapping is defined only in overlap-detection.md. | No action needed — the mapping serves the overlap detection step specifically. | low |
| CS-07 | `rules/constraints.md` L16 vs `rules/domain-frontmatter.md` L19 | 2 | REDUNDANT | P3 | constraints.md L16 says "Domain overlap >50% between files triggers a warning during the user confirmation step (Step 6)". domain-frontmatter.md L24-26 says "Threshold: If overlap ratio > 50%, flag as a potential duplicate/merge candidate during the user confirmation step (Step 6)". Both describe the same threshold and trigger point. domain-frontmatter.md adds the formula, constraints.md just states the rule. | No action needed — constraints.md provides the rule summary, domain-frontmatter.md provides the formula. Legitimate summary/detail split. | low |
| CS-08 | SKILL.md HARD-GATE "Do NOT infer rules not explicitly stated in source documents" vs `rules/constraints.md` L2 | 2 | REDUNDANT | P3 | SKILL.md HARD-GATE says "Do NOT infer rules not explicitly stated in source documents." constraints.md L2 says "Only extract rules that are explicitly stated in source documents -- do not infer". These express the same constraint with nearly identical wording. | Remove one instance — preferably keep the HARD-GATE in SKILL.md and have constraints.md reference it ("Follow the HARD-GATE: do not infer rules"). | low |

### Layer 3 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| CS-T01 | SKILL.md Step 9 vs Step 10 | 3 | TIMING | P3 | Step 9 says "If no drift is found, skip Steps 10-11 and proceed to Step 12." Step 10 says "For each rule classified as drifted or orphaned in Step 9." The gating is correct: Step 10 depends on Step 9's output. Step 11 (Commit) depends on Step 10's output. Step 12 (Vocabulary) runs unconditionally. The flow is correctly ordered. | No action needed. Verified correct. | low |

---

## 3. clean-code

### Structured Summary (SKILL.md)

- **Steps**: 1 (Scope Detection) → 2 (Code Cleanup) → 3 (Quality Gate, optional) → 4 (Cleanup Summary)
- **Key constraints**: HARD-RULE (stage only planning artifact paths — in Step 4 context only); "Core principle: Only modify code within the determined scope. Never change what the code does — only how it does it."
- **References**: `templates/summary.md`
- **Scope Resolution Priority**: 1. User-specified paths → 2. Git diff → 3. Feature context
- **Field names**: SCOPE_COUNT, SCOPE_SOURCE, MODIFIED_COUNT, SKIPPED_COUNT, GATE_RESULT, FILE_CHANGES, DEAD_CODE_COUNT, COMPLEXITY_COUNT, NAMING_COUNT, DUPLICATION_COUNT, OTHER_COUNT

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| CC-01 | `templates/summary.md` vs SKILL.md Step 4 | 2 | INCOMPLETE | P2 | SKILL.md Step 4 "Template Fields" table lists 6 fields: SCOPE_COUNT, SCOPE_SOURCE, MODIFIED_COUNT, SKIPPED_COUNT, GATE_RESULT, FILE_CHANGES. However, `templates/summary.md` includes additional fields not listed in the SKILL.md table: DEAD_CODE_COUNT, COMPLEXITY_COUNT, NAMING_COUNT, DUPLICATION_COUNT, OTHER_COUNT (under "Changes by Type" section). The template has 11 placeholders total, but SKILL.md only documents 6. The 5 additional "count by type" fields are not documented in the SKILL.md Step 4 placeholder table. | Add the 5 "Changes by Type" count fields (DEAD_CODE_COUNT, COMPLEXITY_COUNT, NAMING_COUNT, DUPLICATION_COUNT, OTHER_COUNT) to the SKILL.md Step 4 Template Fields table. | high |
| CC-02 | SKILL.md Step 2 "What to Clean Up" vs `templates/summary.md` "Changes by Type" | 2 | INCOMPLETE | P3 | SKILL.md Step 2 lists 6 cleanup categories: Dead code, Unnecessary complexity, Poor naming, Code duplication, Unnecessary abstractions, Overly verbose patterns. The template's "Changes by Type" section has 5 categories: Dead code removed, Complexity reduced, Naming improved, Duplication eliminated, Other. "Unnecessary abstractions" and "Overly verbose patterns" from SKILL.md are collapsed into "Other" in the template. The mapping from SKILL.md categories to template categories is implicit and undocumented. | Either align the template's categories with SKILL.md's 6 categories, or document the mapping in SKILL.md Step 4 (e.g., "abstractions and verbose patterns are counted under 'Other'"). | medium |
| CC-03 | SKILL.md "HARD-RULE" in Step 8 (commit) | 2 | INCOMPLETE | P3 | SKILL.md has a HARD-RULE at Step 8: "Stage only planning artifact paths — never use `git add -A` or `git add .`". However, clean-code is a code cleanup skill that modifies source files, not planning artifacts. The HARD-RULE text about "planning artifact paths" appears to be copied from quick-tasks Step 8 and is semantically incorrect for clean-code. The intent is clear (only stage files you modified), but the wording "planning artifacts" is misleading. | Rephrase the HARD-RULE to: "Stage only files modified during cleanup — never use `git add -A` or `git add .`." | high |
| CC-04 | SKILL.md "When to Use" vs actual skill structure | 2 | INCOMPLETE | P3 | The skill is referenced as both a standalone command (`/forge:clean-code`) and a pipeline task (`T-clean-code-1`). However, SKILL.md does not describe when `T-clean-code-1` is triggered in the task pipeline — which task types generate it, what config field controls it, or how it differs from standalone invocation. The only hint is "when `auto.cleanCode` is enabled" in the trigger conditions, but the pipeline context is not elaborated. | Add a brief note about the pipeline trigger: "When `auto.cleanCode` is enabled in `.forge/config.yaml`, the `forge task index` command auto-generates a `T-clean-code-1` task for features with code tasks." | low |

### Layer 3 Findings

No TIMING issues found. The 4-step flow (Scope → Cleanup → Quality Gate → Summary) is correctly ordered. The Quality Gate correctly depends on Step 2's output (modified files). The git diff scope detection (Step 1 Priority 2) correctly gates the commit step (Step 4 context).

---

## 4. deep-research

### Structured Summary (SKILL.md)

- **Steps**: Phase 1 (Clarify Needs, 2 rounds AskUserQuestion) → Phase 2 (Execute Research, adaptive) → Phase 3 (Output Report)
- **Key constraints**: HARD-GATE (no code or implementation action); HARD-RULE (AskUserQuestion max 2 rounds / 4 questions per round); HARD-RULE (convergence rule: stop when 2 consecutive rounds yield no new info); HARD-RULE (no auto-commit)
- **References**: `templates/research-report.md`, `rules/research-dimensions.md`
- **Parameters**: topic, --compare, --focus
- **Field names**: TOPIC, MODE, CANDIDATES, DIMENSIONS, KEY_QUESTION, RESEARCH_MODE, CONFIDENCE_LEVEL

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| DR-01 | SKILL.md Phase 1 Q4 vs `rules/research-dimensions.md` | 2 | INCOMPLETE | P3 | SKILL.md Phase 1 Round 2 Q4 says "Select research dimensions to cover (multiSelect) — present the dimension set matching the research mode from `rules/research-dimensions.md` (use `single-tech` dimensions for deep dive mode, `comparison` dimensions for comparison mode). Core dimensions (marked with `*` in `rules/research-dimensions.md`) are pre-selected." However, `rules/research-dimensions.md` marks core dimensions as "Core" (not `*`) in the Core/Optional column. The `*` notation referenced in SKILL.md does not exist in the rules file. | Replace the `*` reference in SKILL.md Phase 1 Q4 with "Core" to match the actual column values in research-dimensions.md. | high |
| DR-02 | `templates/research-report.md` vs SKILL.md Phase 3 | 2 | INCOMPLETE | P2 | SKILL.md Phase 3 says "Write the report to `docs/research/<slug>.md` where `<slug>` is a kebab-case derivation of the topic." The template `templates/research-report.md` does not include a placeholder for `{{SLUG}}` or guidance on slug derivation. Additionally, the template uses `{{TOPIC}}` for both the title and the slug-related fields. The template's frontmatter has `topic: "{{TOPIC}}"` but no `slug` field. The slug is only used in the output file path, not in the template content. This means the slug derivation logic is only in SKILL.md, not in the template. | No action needed — slug is a file-naming concern, not a template content concern. The template correctly uses TOPIC for display purposes. | low |
| DR-03 | `rules/research-dimensions.md` "Usage in AskUserQuestion" L36-43 | 2 | INCOMPLETE | P3 | The "Usage in AskUserQuestion" section says "Allow user to deselect core dimensions if they explicitly don't want them" and "Accept custom dimensions the user proposes — these override the predefined list." These instructions are specific to AskUserQuestion presentation, but they are not reflected in SKILL.md's Phase 1 Round 2 Q4 description. SKILL.md only says "present the dimension set" and "Core dimensions are pre-selected" — it does not mention that users can deselect core dimensions or propose custom ones. | Add to SKILL.md Phase 1 Q4: "Users may deselect core dimensions or propose custom dimensions that override the predefined list." | medium |
| DR-04 | SKILL.md Phase 2 "Convergence rule" vs `rules/research-dimensions.md` | 2 | INCOMPLETE | P3 | SKILL.md Phase 2 has a HARD-RULE: "Convergence rule: Stop researching a topic when 2 consecutive search rounds yield no substantively new information." `rules/research-dimensions.md` does not mention the convergence rule or any search stopping criteria. The convergence rule is entirely in SKILL.md with no supporting detail in the rules file. This is not a problem per se — the rule is self-contained — but it means the rules file provides no guidance on what constitutes "substantively new information" or how to detect convergence. | Consider adding a brief note in research-dimensions.md about convergence detection signals (e.g., "repeated citations, identical findings from different sources"). | low |
| DR-05 | `templates/research-report.md` "Sources" section vs SKILL.md Phase 2 | 2 | INCOMPLETE | P3 | The template has a "Sources" section at the bottom with a table format (Source, URL, Used for). SKILL.md Phase 2 says "Record findings with source URLs for citation." However, SKILL.md Phase 3 does not explicitly instruct the agent to populate the Sources table — it only says "Read `templates/research-report.md` for the report structure." The agent must infer that the Sources section needs to be filled from the template structure alone. | Add an explicit instruction in SKILL.md Phase 3: "Populate the Sources table with all URLs cited during research." | low |

### Layer 3 Findings

No TIMING issues found. The three-phase flow (Clarify → Research → Report) has correct ordering. Phase 2's adaptive research loop (search → cross-reference → record) is internally consistent. The mid-research clarification (max 2 questions) correctly interrupts Phase 2 without breaking the flow.

---

## 5. forensic

### Structured Summary (SKILL.md)

- **Steps**: 1 (Locate Target Sessions) → 2 (Extract Evidence) → 3 (Load Expected Behavior) → 4 (Analyze Deviations) → 5 (Generate Report)
- **Key constraints**: HARD-RULE (report step elapsed time); HARD-RULE (evidence files via --out only, never raw JSONL); HARD-RULE (if skillsUsed empty, ask user); HARD-RULE (causal chain 3 levels deep); EXTREMELY-IMPORTANT (proactively display timing data)
- **References**: `templates/report.md`, `rules/deviation-categories.md`
- **Parameters**: --keyword, --session, --skill, --last, --target
- **Field names**: DEVIATION_CATEGORY, SESSION_IDS, SKILL_NAMES, SEVERITY, EVIDENCE

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| FN-01 | `rules/deviation-categories.md` vs SKILL.md Step 4 | 2 | REDUNDANT | P3 | SKILL.md Step 4 says "Classify each finding using the deviation categories defined in `rules/deviation-categories.md` (instruction-gap, context-starvation, trust-without-verify, wrong-priority, scope-creep, pipeline-gap)." The rules file lists the same 6 categories with descriptions and examples. SKILL.md lists the category names inline, while the rules file provides the full definitions. The category names are duplicated in both files. | Remove the inline category list from SKILL.md Step 4 and just reference the rules file: "Classify each finding per `rules/deviation-categories.md`." | low |
| FN-02 | SKILL.md Step 3 "Resolving skills parent directory" vs actual skill structure | 2 | INCOMPLETE | P3 | Step 3 says to resolve skills by checking "1. `plugins/forge/skills/<skill-name>/SKILL.md` — plugin-distributed skills; 2. `.claude/skills/<skill-name>/SKILL.md` — user-authored project skills." This resolution logic is specific to the forensic skill and not documented in any shared convention file. Other skills that reference files (e.g., eval skill loading rubrics) use "resolve relative to the skills parent directory" without this two-location check. The forensic skill's resolution is more thorough, but the inconsistency with other skills' resolution patterns could cause confusion. | No action needed — forensic's resolution logic is appropriate for its use case (reading skill definitions from any source). The inconsistency is intentional. | low |
| FN-03 | `templates/report.md` "Cross-Session Patterns" vs SKILL.md | 2 | INCOMPLETE | P2 | The report template has a "Cross-Session Patterns" section with a table (Pattern, Sessions, Category). SKILL.md Step 4 describes analyzing individual sessions but does not mention cross-session pattern analysis. Step 5 says "Write the forensic report using the template" without adding instructions for the Cross-Session Patterns section. The agent must infer from the template that it should identify patterns across sessions. | Add a note in SKILL.md Step 4 or Step 5: "After analyzing individual sessions, identify patterns that appear across multiple sessions (e.g., same deviation category, same skill, same trigger condition). Populate the Cross-Session Patterns section of the report." | medium |
| FN-04 | SKILL.md "Prerequisites" section vs Step 1 CLI usage | 2 | INCOMPLETE | P3 | The Prerequisites section says "`forge` CLI must be installed with forensic subcommand (v2.15.0+). Verify: `forge forensic --help`." Step 1 then uses `forge forensic search`. Step 2 uses `forge forensic extract` and `forge forensic subagents`. These CLI commands are assumed to exist but the prerequisite check only verifies the subcommand exists (via --help), not that the specific sub-subcommands (search, extract, subagents) are available. If the CLI version is >= 2.15.0 but the sub-subcommands have different names, the skill would fail at runtime. | The prerequisite check is sufficient for practical purposes — if `forge forensic --help` succeeds, the sub-subcommands should be available. Flagging for awareness. | low |

### Layer 3 Findings

No TIMING issues found. The 5-step sequence (Search → Extract → Load → Analyze → Report) has correct dependency ordering. Step 3 (Load Expected Behavior) depends on Step 2's `skillsUsed` field. Step 4 (Analyze) depends on Step 2's evidence and Step 3's skill definitions. Step 5 (Report) depends on Step 4's findings.

---

## 6. ui-design

### Structured Summary (SKILL.md)

- **Steps**: 1 (Read Manifest) → 2 (Read UI Functions) → 2.5 (Extract Nav Architecture) → 3 (Select Design Style) → 4 (Draft UI Design) → 5 (Write UI Design) → 6 (Update Manifest) → 7 (Auto Eval UI Design) → 8 (Generate Prototype) → 9 (Human Review Gate) → 10 (Reconcile PRD, Optional) → 11 (Next Step)
- **Key constraints**: HARD-GATE (no implementation code except HTML prototypes); HARD-GATE (no proceeding to tech-design until user approves prototype); EXTREMELY-IMPORTANT (eval auto-run check — do NOT use AskUserQuestion when config enables auto-run)
- **References**: `templates/ui-design.md`, `templates/prototype.md`, `templates/manifest-update-ui.md`, `templates/platforms/{web,mobile,tui}.md`, `templates/styles/{vercel,shadcn,tailwind-ui,stripe,apple}.md`, `templates/styles/{modern-dark-tui,minimal-ascii-tui}.md`, `rules/style-selection.md`, `rules/tui-panel-requirements.md`
- **Platform routing**: web → ui-design.md, tui → ui-design-tui.md, mobile → ui-design-mobile.md, multi → per-platform files
- **Field names**: FEATURE_NAME, Component Name, Design System, SLUG

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| UD-01 | `rules/tui-panel-requirements.md` vs `templates/platforms/tui.md` "Structural Requirements" | 2 | REDUNDANT | P2 | `rules/tui-panel-requirements.md` defines 5 mandatory structural requirements (ASCII Layout Mockup, Dimensions, Character Palette, Color Mapping, Edge Cases) with detailed descriptions. `templates/platforms/tui.md` has an identical "Structural Requirements" section listing the same 5 requirements with the same descriptions. The two files duplicate the structural requirements verbatim. Additionally, `templates/ui-design.md` also includes these 5 requirements as subsections of the "TUI Component" template section. The same requirement set is defined in THREE places. | Consolidate: keep the authoritative definition in `rules/tui-panel-requirements.md`, have `templates/platforms/tui.md` reference the rules file instead of inlining the requirements, and have the ui-design.md template reference the rules file for the mandatory check list. | medium |
| UD-02 | SKILL.md Step 5 "Output file naming" vs `rules/style-selection.md` "Multi-Platform Features" | 2 | CONFLICT | P3 | SKILL.md Step 5 defines output file naming: "Web only → ui/ui-design.md; TUI only → ui/ui-design-tui.md; Mobile only → ui/ui-design-mobile.md; Web + mobile → ui/ui-design.md + ui/ui-design-mobile.md; Multi-platform (web + tui) → ui/ui-design-web.md + ui/ui-design-tui.md." Note: "Web only" uses `ui-design.md` but "Multi-platform (web + tui)" uses `ui-design-web.md`. `rules/style-selection.md` "Multi-Platform Features" says: "Single platform (web): ui/ui-design.md; Single platform (tui): ui/ui-design-tui.md; Multi-platform: ui/ui-design-web.md + ui/ui-design-tui.md." The rules file and SKILL.md agree on the naming, but the inconsistency within SKILL.md itself (web-only = `ui-design.md` vs multi-platform web = `ui-design-web.md`) could confuse agents. | No action needed — the naming convention is context-dependent (single vs multi-platform) and is correctly described in both SKILL.md and the rules file. The distinction is intentional. | low |
| UD-03 | `templates/prototype.md` "Navigation Contract" vs SKILL.md Step 2.5 | 2 | INCOMPLETE | P2 | `templates/prototype.md` "Navigation Contract" section says "Before generating any page, load the platform-specific navigation rules: 1. Read the Navigation Architecture section from prd-ui-functions.md; 2. Identify the target platform; 3. Read the corresponding platform file and apply its navigation patterns." This duplicates the work that SKILL.md Step 2.5 already performs (extract navigation architecture and read platform rules). The prototype template assumes the navigation architecture has been extracted, but it re-describes the extraction process as a prerequisite. If the agent follows both SKILL.md Step 2.5 and the prototype template's Navigation Contract, it would read the same files twice. | Add a note in `templates/prototype.md` Navigation Contract: "If Step 2.5 already extracted the Navigation Architecture, reuse those results here. Skip re-reading prd-ui-functions.md." | medium |
| UD-04 | SKILL.md Step 4 "TUI Panel Design Requirements" vs `rules/tui-panel-requirements.md` L13 HARD-RULE | 2 | CONFLICT | P2 | SKILL.md Step 4 says "When platform=tui, each panel MUST include all mandatory structural requirements per `rules/tui-panel-requirements.md`: 5 mandatory structural items: ASCII mockup, dimensions, character palette, color mapping, edge cases; Additional per-panel specs: states, key bindings, data binding." The rules file says "These 5 structural requirements are MANDATORY for every TUI panel. Skipping any item is a spec defect." Both use MANDATORY/MUST language — consistent. However, SKILL.md lists "Additional per-panel specs: states, key bindings, data binding" which are NOT part of the 5 mandatory structural requirements in the rules file. The rules file does list these as additional requirements (L15-18: "In addition to the 5 structural requirements, each TUI panel must also specify: States, Key Bindings, Data Binding"), but they are not covered by the HARD-RULE enforcement. The rules file treats them as softer requirements ("must also specify" without HARD-RULE), while SKILL.md groups them under the same "MUST" language as the 5 mandatory items. | Align the language: either make states/key bindings/data binding part of the HARD-RULE in tui-panel-requirements.md (making 8 mandatory items), or separate them clearly in SKILL.md as "5 mandatory + 3 additional required." | medium |
| UD-05 | `templates/manifest-update-ui.md` "Frontmatter" vs SKILL.md Step 6 | 2 | INCOMPLETE | P3 | SKILL.md Step 6 says "Advance status to `design`" and "Update `manifest.md`". The template `templates/manifest-update-ui.md` says "Update `status` to `design`." in the Frontmatter section. Consistent. However, SKILL.md Step 6 also says "Add traceability links from UI Functions to UI Design sections" which the template covers with the "Traceability (updated)" section. The template is consistent with SKILL.md. | No action needed. Verified consistent. | low |
| UD-06 | SKILL.md Step 8 "TUI Prototype" vs `templates/prototype.md` "TUI Prototype Rules" | 2 | REDUNDANT | P2 | SKILL.md Step 8 "For TUI Platform" describes TUI prototype rules: "HTML simulates a terminal window: dark monospace background, box-drawing characters rendered via CSS; All panels rendered in a single index.html within a black terminal-window div; Bottom area simulates key buttons ([Tab], [1], [q], [:command]) to switch panels." The `templates/prototype.md` "TUI Prototype Rules" section describes the exact same rules with much more detail (terminal CSS, simulated key buttons, panel rendering, panel layout matching). SKILL.md's Step 8 description is a summary that adds no information beyond what the template provides. The partial duplication creates a maintenance risk: changes to TUI prototype rules would need to be synchronized. | Replace SKILL.md Step 8 TUI section with a reference: "For TUI prototypes, follow the detailed TUI Prototype Rules in `templates/prototype.md` — single-file structure with terminal-window container, simulated key buttons, and panel rendering from ASCII mockups." | medium |
| UD-07 | `templates/prototype.md` L39-41 "Platform Reference Files" vs `rules/style-selection.md` | 2 | INCOMPLETE | P3 | `templates/prototype.md` Navigation Contract references platform files: "Web → templates/platforms/web.md, Mobile → templates/platforms/mobile.md, TUI → templates/platforms/tui.md". `rules/style-selection.md` does not reference these platform files at all — it focuses only on style/theme selection. The platform files contain navigation rules, not style rules. This is a correct separation (style-selection handles visuals, platform files handle navigation), but the relationship between the two is not documented. | No action needed — the separation is correct. The platform files are referenced by the prototype template, which is the appropriate consumer. | low |

### Layer 3 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| UD-T01 | SKILL.md Step 7 vs Step 8 vs Step 9 | 3 | TIMING | P3 | SKILL.md Step 7 (Auto Eval UI Design) can invoke eval-ui which may take multiple iterations. Step 8 (Generate Prototype) depends on the final ui-design.md content from Steps 4-5 and the eval result from Step 7. Step 9 (Human Review Gate) depends on the prototype from Step 8. The flow is correct: eval modifies the design → prototype implements the modified design → human reviews the prototype. However, if eval fails (score < 950 and user declines to continue revising), SKILL.md Step 7 says "proceed to prototype anyway." This means a low-scoring design can still generate a prototype, which is the correct graceful degradation behavior. | No action needed. The flow handles edge cases correctly. | low |

---

## 7. learn

### Structured Summary (SKILL.md)

- **Steps**: 1 (Identify Knowledge Type(s)) → 2 (Classify) → 3 (Write) → 4 (Report)
- **Key constraints**: "Core principle: Write knowledge immediately, report for review after. No confirmation gate before writing."
- **References**: `templates/decision-entry.md`, `templates/lesson-entry.md`, `templates/convention-entry.md`
- **Knowledge Types**: decision, lesson, convention, business-rule
- **Classification**: 8-category vocabulary (architecture, interface, data-model, dependencies, error-handling, testing, security, local-dev-deployment)
- **Directory Bootstrap**: auto-create all target directories if missing
- **Field names**: DATE, FEATURE_SLUG, DECISION, RATIONALE, SOURCE, TAGS, TITLE

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| LN-01 | `templates/lesson-entry.md` "Tag Vocabulary" vs SKILL.md Step 2 "Classify" | 2 | REDUNDANT | P3 | Both files define the same 8-category vocabulary. SKILL.md Step 2 has a table with Category/Tag/Decision Type File columns. `templates/lesson-entry.md` has a table with Tag/Domain columns listing the same 8 tags. The vocabulary is defined identically in both places. | Remove the Tag Vocabulary table from `templates/lesson-entry.md` and reference SKILL.md Step 2 or a shared vocabulary source. | low |
| LN-02 | `templates/convention-entry.md` "Domain Derivation" vs consolidate-specs `rules/domain-frontmatter.md` | 2 | REDUNDANT | P2 | `templates/convention-entry.md` has a "Domain Derivation" section (4 steps: extract tokens, extract nouns, deduplicate, 3-7 keywords) that is nearly identical to `consolidate-specs/rules/domain-frontmatter.md` "Domain Derivation Rules" (4 steps with the same logic). Both define the same algorithm. This is cross-skill redundancy — the learn skill's convention-entry template and consolidate-specs' domain-frontmatter rule describe the same domain derivation algorithm. | Extract the shared domain derivation algorithm into a common reference (e.g., a shared rule file or a documentation section that both skills reference). | medium |
| LN-03 | `templates/convention-entry.md` "Project-Global ID Encoding" vs consolidate-specs `rules/spec-classification.md` | 2 | REDUNDANT | P2 | `templates/convention-entry.md` defines project-global ID encoding (prefix, sequence, format: BIZ-<domain>-NNN / TECH-<topic>-NNN, examples). `consolidate-specs/rules/spec-classification.md` defines the same encoding with identical format and examples. Both files describe the same ID scheme. This is the same cross-skill redundancy pattern as LN-02. | Same as LN-02: extract shared ID encoding rules into a common reference. | medium |
| LN-04 | SKILL.md Step 3 "Business-Rule Entry" L113-119 vs `templates/convention-entry.md` | 2 | CONFLICT | P3 | SKILL.md Step 3 "Business-Rule Entry" says "Read `templates/convention-entry.md` for the entry format and project-global ID encoding." However, `templates/convention-entry.md` is titled "Convention / Business-Rule Entry Template" — it serves dual purpose for both conventions AND business rules. SKILL.md Step 3 "Convention Entry" (L106-111) also says "Read `templates/convention-entry.md`." Both convention and business-rule entry types share the same template file. The template correctly handles both (with separate sections for TECH-* and BIZ-* entries). This is not a true conflict but a shared template pattern that could be clearer. | Rename the template to clarify dual purpose: `templates/spec-entry.md` or add a clear header: "This template is used for BOTH convention and business-rule entries." | low |
| LN-05 | SKILL.md "Auto-Generated Vocabulary" section vs consolidate-specs `rules/vocabulary-generation.md` | 2 | INCOMPLETE | P2 | SKILL.md "Auto-Generated Vocabulary" section says "When available, `/learn` reads the vocabulary index generated by `/consolidate-specs` to refine classification suggestions." It says the vocabulary is at `docs/.vocabulary.md` (from consolidate-specs template). However, SKILL.md Step 2 "Classify" does not reference the auto-generated vocabulary — it only mentions the "shared 8-category vocabulary" table. The Step 2 "Classification behavior" says "When auto-generated vocabulary is available (from `/consolidate-specs`), use it to refine suggestions" but does not specify HOW the auto-generated vocabulary refines suggestions beyond the base 8 categories. | Add a concrete example in SKILL.md Step 2: "When `docs/.vocabulary.md` exists, load its Domains table. Use these domain keywords to suggest more specific tags (e.g., if the vocabulary contains `error-handling` with sub-keywords `status, response, stderr`, suggest those as tag options)." | medium |

### Layer 3 Findings

No TIMING issues found. The 4-step flow (Identify → Classify → Write → Report) is correctly ordered. Multi-type detection in Step 1 correctly feeds multiple write operations in Step 3. Directory Bootstrap in Step 3 correctly runs before write operations.

---

## 8. Summary Statistics

| Metric | Value |
|--------|-------|
| Skills audited | 7 |
| Total SKILL.md files read | 7 |
| Total associated files read | 40+ |
| Layer 2 findings | 35 |
| Layer 3 findings | 2 (0 confirmed issues, 2 verified-correct) |
| CONFLICT | 5 |
| INCOMPLETE | 18 |
| REDUNDANT | 9 |
| TIMING (confirmed) | 0 |
| High confidence | 6 |
| Medium confidence | 13 |
| Low confidence | 18 |

### Severity Distribution

| Severity | Count | Finding IDs |
|----------|-------|-------------|
| P1 | 0 | — |
| P2 | 14 | QT-01, QT-05, CC-01, DR-02, DR-03, FN-03, CS-02, CS-03, UD-01, UD-03, UD-04, UD-06, LN-02, LN-03, LN-05 |
| P3 | 21 | QT-02, QT-03, QT-04, QT-06, CC-02, CC-03, CC-04, DR-01, DR-04, DR-05, FN-01, FN-02, FN-04, CS-01, CS-04, CS-05, CS-06, CS-07, CS-08, UD-02, UD-05, UD-07, LN-01, LN-04 |

### Category Distribution

| Category | Count |
|----------|-------|
| CONFLICT | 5 |
| INCOMPLETE | 18 |
| REDUNDANT | 9 |
| TIMING | 2 (verified correct) |

### Per-Skill Finding Count

| Skill | CONFLICT | INCOMPLETE | REDUNDANT | TIMING | Total |
|-------|----------|------------|-----------|--------|-------|
| quick-tasks | 0 | 3 | 0 | 0 | 6 |
| consolidate-specs | 1 | 4 | 4 | 1 | 10 |
| clean-code | 0 | 4 | 0 | 0 | 4 |
| deep-research | 0 | 5 | 0 | 0 | 5 |
| forensic | 0 | 3 | 1 | 0 | 4 |
| ui-design | 1 | 3 | 2 | 1 | 7 |
| learn | 1 | 2 | 3 | 0 | 6 |

---

## 9. Cross-Skill Patterns

### Pattern 1: Convention Loading `domains` Filtering (recurring from Batch A)

quick-tasks SKILL.md Step 0 uses `domains` frontmatter filtering for Convention file loading — the same pattern flagged in batch A for gen-test-scripts, breakdown-tasks, and tech-design (GTS-01, BT-01, TD-05). This confirms the pattern is systematic across all skills that load Convention files. Total skills affected: 4 (3 from batch A + quick-tasks from batch B).

### Pattern 2: Template Hardcoded Defaults Override SKILL.md Logic (QT-02, QT-03)

`templates/task.md` hardcodes `complexity: "medium"` and `type: "coding.feature"` in frontmatter, while SKILL.md defines multi-value heuristics for both fields. The templates override the SKILL.md logic unless the agent manually edits the generated content. This pattern does not appear in batch A skills because batch A skills use templates with more placeholders or conditional logic.

### Pattern 3: TUI Requirements Triplication (UD-01)

The 5 mandatory TUI panel structural requirements are defined in THREE files: `rules/tui-panel-requirements.md` (authoritative), `templates/platforms/tui.md` (duplicated), and `templates/ui-design.md` (template inline). This is the most severe redundancy pattern in batch B — any change to TUI requirements must be synchronized across three files.

### Pattern 4: Domain Derivation / Project-Global ID Duplication (LN-02, LN-03)

The `learn` skill's `templates/convention-entry.md` and `consolidate-specs`' `rules/domain-frontmatter.md` and `rules/spec-classification.md` define the same domain derivation algorithm and project-global ID encoding scheme independently. This creates a maintenance burden when either skill's conventions change.

### Pattern 5: SKILL.md Summaries That Duplicate Template/Rules Content (UD-06, FN-01, CS-02, DR-01)

Multiple skills include inline summaries in SKILL.md that duplicate content from templates or rules files. While these summaries improve readability, they create maintenance risk when the referenced files change but the SKILL.md summary is not updated. The pattern appears in: ui-design Step 8 (TUI prototype rules), forensic Step 4 (deviation categories), consolidate-specs Step 9 (drift classification), deep-research Phase 1 Q4 (dimension notation).
