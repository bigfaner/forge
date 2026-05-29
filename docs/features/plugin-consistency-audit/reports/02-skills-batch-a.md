# Report 02: Skills Deep Audit — Batch A (Layer 2-3)

**Baseline commit**: `e43323a16cfd7d5265d9bed2247a582e7f087445`
**Date**: 2026-05-30
**Scope**: 7 skills (eval, gen-test-scripts, run-tests, tech-design, write-prd, brainstorm, breakdown-tasks)
**Audit order**: run-tests → write-prd → brainstorm → gen-test-scripts → breakdown-tasks → tech-design → eval (randomized)
**Method**: Layer 2 (Instruction Consistency) + Layer 3 (Timing & Flow)

---

## 1. run-tests

### Structured Summary (SKILL.md)

- **Steps**: 0 (Stale State Recovery) → 1 (Detect Surface) → 1.5 (Discover Journeys) → 2 (Load Orchestration Rules) → 3 (Env Check) → 4 (Execute Sequence) → 5 (Parse Results) → 6 (Generate Report)
- **Key constraints**: HARD-GATE (no test modification, no skipping failed tests, no retrying probe after failure); HARD-RULE (env detection no auto-fix, teardown mandatory, >30% batch failure gate)
- **References**: `rules/env-check.md`, `rules/result-parsing.md`, `rules/confidence.md`, `rules/failure-diagnosis.md`, `rules/test-isolation.md` (cross-skill from gen-test-scripts), `rules/surfaces/{web,api,cli,tui,mobile}.md`, `templates/test-report.md`
- **Field names**: surface-type, JOURNEYS, test-state.json

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| RT-01 | `rules/env-check.md` L49 | 2 | CONFLICT | P1 | Web surface environment check #3 hardcodes `npx playwright install` in the "How to Verify" column ("Run framework install command (e.g., `npx playwright install --dry-run`). Exit code 0 = pass") and Repair Suggestion ("Run `npx playwright install`"). SKILL.md is fully profile-agnostic (no Playwright reference anywhere) and v3.0.0 test profile system replaced Playwright hardcoding with pluggable Convention-based frameworks. The rule file's repair suggestion directly contradicts SKILL.md's design intent of framework-agnosticism. | Replace hardcoded Playwright references with Convention-derived framework commands. Read the loaded Convention file's `framework` section to determine the actual browser automation framework and its install command. Use generic phrasing: "Run the browser automation framework install command (per Convention file). Verify via the framework's dry-run or version check." | high |
| RT-02 | `rules/env-check.md` L13 | 2 | CONFLICT | P2 | env-check.md references "Look up the corresponding surface rule file in **gen-journeys** skill's `rules/surface-<type>.md`" for Environment Readiness Checks. However, the run-tests skill has its own `rules/surfaces/` directory (web.md, api.md, cli.md, tui.md, mobile.md). The env-check file points to gen-journeys for surface rules, but SKILL.md Step 2 loads from run-tests' own `rules/surfaces/<type>.md`. The two paths resolve to different skill directories. | Clarify: the env-check "How It Works" section should reference the run-tests skill's own `rules/surfaces/<type>.md` for orchestration rules, not gen-journeys. Or clarify that env-check reads the "Environment Readiness Checks" table from gen-journeys' surface rules (which contain env check items), while the orchestration sequence comes from run-tests' own surface rules. | high |
| RT-03 | `rules/env-check.md` L98-104 | 2 | CONFLICT | P2 | The "Integration with SKILL.md" section describes a workflow step sequence ("between Step 3 Setup and Step 4 Pre-check, or as a new early step") that does not match the actual SKILL.md workflow. SKILL.md places env-check as Step 3 (explicit: "Read the rule file `rules/env-check.md` for the detection framework. Execute the environment checks defined for the detected surface type."). The env-check file's "Integration" section uses incorrect step references. | Update the "Integration with SKILL.md" section to reference the correct SKILL.md Step 3. | medium |
| RT-04 | `rules/env-check.md` L70 | 2 | CONFLICT | P2 | Mobile checks are declared "best-effort -- all items are non-blocking" via HARD-RULE. However, check #1 "Maestro CLI installed" has `Blocking: No` and check #2 "Emulator/simulator available" has `Blocking: No`. The HARD-RULE and the table columns are consistent, but the table's "Repair Suggestion" for check #1 provides an install command as if installation is expected, which contradicts the "best-effort" stance. | No action needed — the consistency between HARD-RULE and table is correct. The repair suggestion is informational (what the user *can* do, not what they *must* do). Flag as reviewed, no issue. | low |
| RT-05 | `rules/confidence.md` L21-22 | 2 | INCOMPLETE | P2 | confidence.md references `.eval-status.json` with `"status": "eval-skipped"` and `"status": "eval-bypassed"` for forced downgrade detection. These file paths (`testing/<journey>/.eval-status.json` and `testing/<journey>/contracts/.eval-status.json`) are not documented in SKILL.md or any other run-tests rule file. The eval-status file is produced by the eval skill but the run-tests skill's SKILL.md does not mention it as a prerequisite or artifact to check. | Add a note in SKILL.md Step 5 or Step 6 that the confidence calculation may read `.eval-status.json` files from the testing directory. | medium |
| RT-06 | `templates/test-report.md` | 2 | INCOMPLETE | P3 | The template uses `{{SCREENSHOTS_SECTION}}` placeholder but SKILL.md Step 6 does not mention screenshots in the report generation. Only surface tests (web/mobile) would produce screenshots; the template has no conditional logic for when screenshots are absent. | Add conditional rendering logic or remove the placeholder for non-visual surfaces. | low |

### Layer 3 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| RT-T01 | `rules/env-check.md` vs SKILL.md | 3 | TIMING | P3 | env-check.md "Integration with SKILL.md" section references a conceptual step ordering that predates the current SKILL.md workflow. No actual timing issue in the SKILL.md workflow itself — Step 3 correctly loads env-check.md after surface detection. The stale integration section may mislead maintainers. | Remove or update the "Integration with SKILL.md" section in env-check.md to match the current step numbers. | low |

---

## 2. write-prd

### Structured Summary (SKILL.md)

- **Steps**: 1 (Explore Context) → 2 (Assess Scope) → 3 (Ask Questions) → 4 (Propose Approaches) → 5 (Present PRD Sections) → 6 (Write PRD Spec) → 7 (Write User Stories) → 7A (Write Spec-Only PRD) → 8 (Write UI Functions) → 9 (Create Manifest) → 10 (Self-Check) → 11 (Review & Commit) → 12 (Adversarial Eval) → 13 (Knowledge Review)
- **Intent branches**: new-feature (default), refactor, cleanup
- **Key constraints**: HARD-GATE (no code until PRD approved); HARD-RULE (no tech selection); HARD-RULE (no auto-commit)
- **References**: `templates/prd-spec.md`, `templates/prd-user-stories.md`, `templates/prd-ui-functions.md`, `templates/manifest.md`, `rules/self-check.md`, `rules/ui-functions.md`, `rules/knowledge-extraction.md`, `examples/ask-questions.md`, `examples/propose-approaches.md`, `examples/user-stories.md`

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| WP-01 | `rules/self-check.md` | 2 | INCOMPLETE | P2 | self-check.md does not include any check for the refactor/cleanup intent branch. It assumes new-feature intent throughout (checks for "User stories" with Given/When/Then AC, "Flow diagram" with Mermaid, "Placement completeness", "Page Composition valid"). For refactor/cleanup intent, user stories, flow diagrams, and UI functions are not generated (per SKILL.md Step 7A and Step 8 intent gates). The self-check rules have no conditional logic for these cases. | Add intent-aware conditional checks. When intent is refactor/cleanup: verify the three mandatory fields (Change Scope, Constraints, Verification Criteria) are present; skip user stories, flow diagram, and UI-related checks. | high |
| WP-02 | `rules/knowledge-extraction.md` vs SKILL.md Step 13 | 2 | REDUNDANT | P3 | SKILL.md Step 13 says "After Step 12 (Adversarial Eval Prompt) completes, run knowledge auto-extraction from the PRD. Full extraction flow, knowledge type definitions, notable-knowledge heuristics, and deduplication rules — see `rules/knowledge-extraction.md`." The rules file is a complete standalone reference that re-describes the full extraction flow. This is a legitimate expansion (rules add implementation detail), not pure redundancy. | No action needed — rules provide necessary implementation detail that SKILL.md intentionally defers. | low |
| WP-03 | SKILL.md Step 8 vs `rules/ui-functions.md` | 2 | CONFLICT | P2 | SKILL.md Step 8 heading says "Write UI Functions (mandatory for UI features)" and states "This step is **mandatory** when the feature has any UI surface. Skip only for backend/API/CLI-only features with no UI surface." However, `rules/ui-functions.md` uses HARD-RULE level language ("Every UI Function MUST have a Placement section. Missing Placement → error, do not proceed"). The rules file adds a stricter constraint (Placement mandatory) that SKILL.md does not mention. This is an additive constraint in rules, not a conflict per se, but SKILL.md should at least reference the Placement requirement. | Add a note in SKILL.md Step 8 that Placement sections are mandatory for every UI Function, per `rules/ui-functions.md`. | medium |
| WP-04 | SKILL.md Step 3 vs `examples/ask-questions.md` | 2 | INCOMPLETE | P3 | SKILL.md Step 3 references `examples/ask-questions.md` but this file was not audited in this batch (not in the 7 skill directories). Cross-skill reference validation is out of scope per audit proposal. | N/A (cross-skill reference, noted for completeness). | low |

### Layer 3 Findings

No TIMING issues found. Step numbering is consistent between SKILL.md and all referenced rule files.

---

## 3. brainstorm

### Structured Summary (SKILL.md)

- **Steps**: 1 (Analyze Context) → 2 (Walk Design Tree) → 3 (Propose Approaches) → 4 (Define Scope) → 4.5 (Infer Feature Intent) → 5 (Write Proposal) → 6 (Commit) → 7 (Adversarial Eval)
- **Key constraints**: HARD-GATE (no code); HARD-RULE (no tech selection); HARD-RULE (no auto-commit)
- **References**: `templates/proposal.md`, `rules/challenge-protocol.md`, `rules/sc-consistency.md`
- **Challenge tools**: 5 Whys, XY Detection, Assumption Flip, Stress Test (embedded in Decision Clusters)

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| BS-01 | `templates/proposal.md` L102 | 2 | INCOMPLETE | P3 | The proposal template's "Next Steps" section says "Proceed to `/write-prd` to formalize requirements" but does not mention the eval step. SKILL.md Step 7 adds adversarial eval (eval-proposal) between commit and write-prd. The template's Next Steps skips the eval step. | Update the template's Next Steps to mention eval-proposal as an intermediate step: "Proceed to `/eval-proposal` for adversarial evaluation, then `/write-prd` to formalize requirements." | medium |
| BS-02 | `rules/sc-consistency.md` L101-102 | 2 | INCOMPLETE | P2 | The sc-consistency rule produces a `consistency_check_result` field that must be included in the proposal's SC section. However, `templates/proposal.md` has no placeholder or mention of this field. The template does not guide the agent to include the consistency check result. | Add a `consistency_check_result` field placeholder to the template's Success Criteria section, or add an instruction note. | medium |
| BS-03 | `rules/challenge-protocol.md` | 2 | REDUNDANT | P3 | The challenge protocol Need Gate section (lines 9-48) duplicates SKILL.md Step 2's "Need Gate" description almost verbatim. Both describe the same three checks (Simpler Alternative, Real Need, Timing) with the same override logic. The rules file adds search strategy details and tone guidelines, making it a legitimate expansion, but the overlap in the gate description itself is pure redundancy. | Remove the Need Gate description from SKILL.md Step 2 and replace with a reference: "Apply the Need Gate checks defined in `rules/challenge-protocol.md`." Keep only the trigger condition in SKILL.md. | low |

### Layer 3 Findings

No TIMING issues found. Step 4.5 (Infer Intent) is correctly positioned after Define Scope and before Write Proposal.

---

## 4. gen-test-scripts

### Structured Summary (SKILL.md)

- **Steps**: 0 (Load Convention) → 0.5 (Surface Detection) → 1 (Code Reconnaissance) → 2 (Read Contracts) → 2.5 (Load Type Rules) → 3 (Generate Test Code) → 4 (Compile Gate)
- **Key constraints**: HARD-RULE (no domains filtering, no silent default, batch generation per journey, single journey per invocation, types/ golden rules precedence, _shared.md always loaded, no staging area, every journey needs smoke test, no hardcoded secrets, gen-failed files not deleted, at most 1 auto-retry)
- **References**: `rules/convention-guide.md`, `rules/quality-gates.md`, `rules/run-to-learn.md`, `rules/step-0.5-validation.md`, `rules/step-1-contract-loading.md`, `types/_shared.md`, `types/{cli,tui,ui,mobile,api}.md`

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| GTS-01 | `rules/convention-guide.md` L33 | 2 | CONFLICT | P2 | convention-guide.md HARD-RULE says "Do NOT use `domains` frontmatter filtering." breakdown-tasks SKILL.md Step 0 says "Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.)." tech-design SKILL.md Step 0 also says "Load Convention files from `docs/conventions/` by `domains` frontmatter." gen-test-scripts SKILL.md Step 0 says "Load Convention files from `docs/conventions/` by `domains` frontmatter" BUT then its own rules/convention-guide.md explicitly forbids domains filtering. The SKILL.md itself contradicts its own rules file. | Update gen-test-scripts SKILL.md Step 0 to match convention-guide.md: selection is based on index.md descriptions and project context, not domains frontmatter filtering. | high |
| GTS-02 | `rules/run-to-learn.md` L46 | 2 | CONFLICT | P2 | run-to-learn.md HARD-RULE says "R2L environment readiness is a separate concern from the run-tests env-check (task 2.8)." However, run-tests skill does not have a "task 2.8" — the reference is to a task number in a different context (likely the gen-test-scripts pipeline step numbering, but run-tests is a separate skill). The cross-reference is misleading. | Replace the cross-reference with a clear statement: "R2L environment readiness is a separate concern from the run-tests skill's Step 3 (Environment Readiness Check)." | medium |
| GTS-03 | SKILL.md Step 0 vs `rules/convention-guide.md` L36-75 | 2 | INCOMPLETE | P2 | SKILL.md Step 0.3 says "check required sections: `framework`, `discovery`, `structure`, `assertions`" but does not mention the optional sections (Import Patterns, Code Style, Anti-patterns, Helpers) that convention-guide.md describes. Later in Step 3, SKILL.md says "Use the Convention's Code Style and Helpers sections as the template" — referencing sections that were never loaded or validated. | Add a note in SKILL.md Step 0.3 that optional sections (Import Patterns, Code Style, Anti-patterns, Helpers) should also be read if present, as they are used in Step 3 generation. | medium |
| GTS-04 | SKILL.md Step 4.4 vs `rules/quality-gates.md` | 2 | REDUNDANT | P3 | SKILL.md Step 4.4 says "Antipattern Guard & Duplicate Name Check: See `rules/quality-gates.md`." The quality-gates rule file contains the full antipattern table and error handling table. SKILL.md does not duplicate the content — it properly defers. This is correct reference-based design. | No action needed. Noted as reviewed. | low |

### Layer 3 Findings

No TIMING issues found. The pipeline steps (0 → 0.5 → 1 → 2 → 2.5 → 3 → 4) are correctly ordered with proper dependencies.

---

## 5. breakdown-tasks

### Structured Summary (SKILL.md)

- **Steps**: 0 (Resolve Language) → 1 (Read Documents) → 2 (Map → Tasks) → 3 (Derive Phases & Dependencies) → 4 (Create Task Files) → 5 (Generate index.json) → 6 (Validate) → 7 (Update Manifest) → 8 (Commit)
- **Condition-Rule Matrix**: phase-detection, ui-placement, db-schema, existing-code-split
- **Key constraints**: HARD-GATE (max 6 AC per task); HARD-RULE (no silent default language, naming conventions, stage only planning artifacts)
- **References**: `templates/task.md`, `templates/task-doc.md`, `templates/manifest-update-tasks.md`, `rules/phase-detection.md`, `rules/ui-placement.md`, `rules/db-schema.md`, `rules/existing-code-split.md`

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| BT-01 | SKILL.md Step 0 vs gen-test-scripts `rules/convention-guide.md` | 2 | CONFLICT | P2 | breakdown-tasks SKILL.md Step 0 says "Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.)." This is the same pattern as gen-test-scripts SKILL.md Step 0, and both conflict with the gen-test-scripts `rules/convention-guide.md` HARD-RULE that forbids domains filtering. While this is a cross-skill issue, the breakdown-tasks skill is loading conventions for language detection (a different purpose than test generation), so the `domains` approach may be intentional here. | Verify whether breakdown-tasks should use domains filtering or the same index.md-based approach as gen-test-scripts. If different, document the difference. If same, align with convention-guide.md's approach. | medium |
| BT-02 | `rules/ui-placement.md` L56 | 2 | INCOMPLETE | P3 | The "Placement format note" describes the canonical form as `<mode>:<target-page-value>` but this combined format is only defined in this rule file — the PRD template (`write-prd/templates/prd-ui-functions.md`) stores them as separate fields. The conversion logic is documented but the mapping from separate fields to canonical form could be clearer. | No action needed — the note is sufficient. Flagging for awareness. | low |
| BT-03 | SKILL.md "Docs-Only Fast Path" | 2 | INCOMPLETE | P2 | SKILL.md says "When all tasks are `type: "doc"` (non-compilable, non-runnable output only), skip Step 0." However, the Type Assignment table does not include a "detect docs-only" step — it only lists individual type values. The fast path detection logic ("Step 1 scans artifacts → every element targets non-compilable files → docs-only") is described in prose but not formalized as a check. | Add a formal detection sub-step in Step 1: "Scan all Affected Files from design elements. If all entries target non-compilable extensions (.md, .yaml, .json under docs/), set docs-only flag and skip Step 0." | medium |

### Layer 3 Findings

No TIMING issues found. The Condition-Rule Matrix correctly gates rule loading based on artifact presence, and the step sequence (0 → 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8) has proper dependency ordering.

---

## 6. tech-design

### Structured Summary (SKILL.md)

- **Steps**: 0 (Detect Language) → 1 (Read PRD) → 2 (Explore Context) → 3 (Identify Decisions) → 4 (Ask Questions) → 5 (Draft Design) → 6 (Get Approval) → 7 (Archive Decisions) → 8 (Write Design Documents) → 9 (Update Manifest) → 10 (Adversarial Eval) → 11 (Auto-Extract Knowledge)
- **Intent branches**: new-feature (default), refactor, cleanup
- **Key constraints**: HARD-GATE (no code until approved); HARD-RULE (no silent default language)
- **References**: `templates/tech-design.md`, `templates/api-handbook.md`, `templates/er-diagram.md`, `templates/schema.sql`, `templates/decision-entry.md`, `templates/manifest-update-design.md`, `rules/design-quality-checks.md`, `rules/decision-archiving.md`, `rules/knowledge-extraction.md`, `examples/ask-question.md`, `examples/exploration.md`

### Layer 2 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| TD-01 | `rules/decision-archiving.md` L88 | 2 | CONFLICT | P3 | decision-archiving.md references `examples/ask-question.md` and `examples/exploration.md` for question formatting and context exploration. Task 1 inventory (Report 01) flagged both examples as ORPHAN (O-06, O-07) — they are not referenced in SKILL.md directly. The rule file does reference them, confirming they are not truly orphaned (second-level reference). | Add references to these examples in SKILL.md for discoverability, per Task 1 recommendation. | low |
| TD-02 | `rules/knowledge-extraction.md` vs write-prd `rules/knowledge-extraction.md` | 2 | REDUNDANT | P2 | tech-design's knowledge-extraction.md and write-prd's knowledge-extraction.md share 80%+ identical content (same Knowledge Types table, same Notable Knowledge Heuristics, same Deduplication section, same auto-save configuration check flow). The differences are: trigger parameter (`tech-design` vs `write-prd`), artifact scope (design doc vs PRD docs), and Step 7 coordination note (tech-design only). This is cross-skill redundancy within the batch. | Consider extracting the shared logic into a common rules file (e.g., `shared/knowledge-extraction-common.md`) and having each skill reference it with skill-specific parameters. | medium |
| TD-03 | `rules/design-quality-checks.md` L9 | 2 | INCOMPLETE | P2 | The 5.1 PRD Coverage Verification says "For each AC from `prd-user-stories.md`" but does not account for the refactor/cleanup intent where prd-user-stories.md is not generated. For refactor/cleanup, SKILL.md Step 1 says to extract AC from the PRD spec's "Verification Criteria" section instead. The quality check rule does not have this conditional logic. | Add intent-aware conditional: "For `new-feature` intent: verify AC from `prd-user-stories.md`. For `refactor`/`cleanup` intent: verify AC from the PRD spec's Verification Criteria section." | high |
| TD-04 | `rules/design-quality-checks.md` L54-58 | 2 | CONFLICT | P2 | The 5.5 DB Schema Branch rule says "When `db-schema: 'no'`: scan content for keywords: TABLE, REFERENCES, FOREIGN KEY, CREATE TABLE, ALTER TABLE, migration, schema." If found, it prompts to switch to db-schema "yes" path. However, for `refactor`/`cleanup` intent, SKILL.md explicitly says ER Diagram and Schema are skipped (not generated). The quality check does not gate the db-schema scan behind intent check, potentially prompting the user to generate ER diagrams for refactoring tasks where they are not applicable. | Add intent gate: "Skip this check when intent is `refactor` or `cleanup` — db-schema artifacts are not generated for these intents." | high |
| TD-05 | SKILL.md Step 0 vs gen-test-scripts `rules/convention-guide.md` | 2 | CONFLICT | P2 | Same issue as BT-01/GTS-01: tech-design SKILL.md Step 0 uses `domains` frontmatter filtering for Convention loading, which contradicts gen-test-scripts' convention-guide.md HARD-RULE against domains filtering. | Same fix as GTS-01: align with the convention-guide.md approach if conventions are shared, or document the intentional difference. | medium |

### Layer 3 Findings

No TIMING issues found. The step sequence is correctly ordered with proper intent branching.

---

## 7. eval

### Structured Summary (SKILL.md)

- **Architecture**: Mermaid flowchart defining scorer→gate→revise loop
- **Steps**: 1 (Resolve Type & Load Rubric) → Phase 0 (Freeform Expert Review, proposal only) → Expert Dispatch → Iteration Init → 2 (Invoke Scorer) → 3a/3b (Decision Gate) → 4 (Invoke Reviser) → 5 (Final Report) → 5.5 (Cleanup) → 6 (Next Step)
- **Sub-phases**: P0.1 (Expert Reuse), P0.2 (Expert Inference), P0.3 (Freeform Review), P0.4 (Extract Findings), P0.5a-g (Pre-Revision)
- **Key constraints**: EXTREMELY-IMPORTANT (main session owns the loop, scorer/reviser ALWAYS via Agent tool)
- **References**: 11 rule files, 15 expert files, 12 rubric files

### Layer 2 Findings (by sub-directory group)

#### rules/ group

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| EV-01 | `rules/freeform-injection.md` | 2 | CONFLICT | P1 | The file is marked `status: deprecated` with `deprecated-by: eval-freeform-pre-revision`. However, it is still listed as "loaded by freeform-pipeline.md at runtime" in SKILL.md Phase 0 section: "`rules/freeform-injection.md` — **DEPRECATED**: legacy injection logic inlined into `rules/scorer-composition.md` (kept for historical reference only)". SKILL.md correctly marks it as deprecated, but the file still exists in the rules/ directory and could be accidentally loaded. The deprecation frontmatter provides restore instructions, which is good. | The current state (deprecated + SKILL.md note) is sufficient. Consider moving to a `rules/_deprecated/` subdirectory or prefixing with `_` to further reduce accidental loading risk. | high |
| EV-02 | `rules/scorer-composition.md` L4 | 2 | CONFLICT | P3 | The file title is "Expert Dispatch Table" but the file actually contains both the dispatch table AND scorer prompt composition rules AND freeform integration rules. The title underrepresents the file's scope. | Rename the file's top-level heading to "Expert Dispatch, Scorer Composition & Freeform Integration". | low |
| EV-03 | `rules/reviser-composition.md` L29-31 | 2 | INCOMPLETE | P2 | The "Reviser Type-Specific Constraints" section is incomplete — `journey` and `contract` entries have incomplete sentences: "After reviser completes:" with no continuation. The `consistency` entry also has a dangling instruction. These appear to be cut off or placeholder text. | Complete the reviser type-specific constraints for journey, contract, and consistency types. Specify what happens after the reviser completes for each type. | high |
| EV-04 | `rules/rubric-reference.md` | 2 | INCOMPLETE | P3 | The rubric reference table lists all rubrics but does not include a `harness` entry in the rubric file column, despite the report-format.md mentioning `harness` as a type with "priority improvement table". The dispatch table in scorer-composition.md lists `harness` as "(uses generic inline fallback)". | Add a `harness` row to the rubric reference table with a note that it uses generic inline fallback (no dedicated rubric file). | low |
| EV-05 | `rules/pre-processing.md` | 2 | INCOMPLETE | P3 | The pre-processing table mentions `harness` type in report-format.md but not in pre-processing.md. If `harness` type has no pre-processing, this is acceptable, but it should be explicitly stated. | Add `harness` to the pre-processing table with "No special pre-processing" or a dash, for completeness. | low |
| EV-06 | `rules/report-format.md` | 2 | INCOMPLETE | P3 | The report format does not include `proposal` type-specific additions. Given that proposal type has the most complex flow (Phase 0 freeform, pre-revision, baseline score, rollback), the report format should mention proposal-specific sections (Pre-Revision Findings Triage, Baseline Score Comparison, Baseline Drift Alert) — though these are already documented in SKILL.md Steps 5.1-5.4. | Add `proposal` to the type-specific additions list referencing SKILL.md Steps 5.1-5.4. | low |

#### experts/ group

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| EV-07 | `experts/` (all 15 files) | 2 | INCOMPLETE | P3 | Task 1 (Report 01) identified all 15 expert files as ORPHAN (O-02) — not directly referenced in SKILL.md. They are loaded indirectly via rules files (freeform-pipeline.md, scorer-composition.md). This is second-level reference pattern, not a true orphan. | Add a summary reference in SKILL.md: "Expert files: `experts/freeform/` (inference, template, review), `experts/protocol/` (scorer, reviser), `experts/scorer/` (domain experts for each eval type). See rules files for loading logic." | low |

#### rubrics/ group

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| EV-08 | `rubrics/` (12 files) | 2 | INCOMPLETE | P3 | Task 1 identified rubrics as pattern-referenced (`rubrics/<type>.md`). SKILL.md Step 1.1 correctly uses this pattern. No actual issue — the pattern reference resolves correctly for all types. | No action needed. Noted as verified. | low |

### Cross-Group Findings (汇总轮)

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| EV-09 | `rules/pre-processing.md` vs `rules/scorer-composition.md` | 2 | INCOMPLETE | P2 | pre-processing.md defines context injection (loading conventions/business-rules files based on rubric `context` frontmatter). scorer-composition.md also describes context injection ("Context Injection: If `CONTEXT_CONTENT` was loaded in Step 1.4, append..."). The two files describe the same mechanism from different perspectives — pre-processing.md describes the loading, scorer-composition.md describes the injection into the prompt. This is a legitimate separation of concerns, not redundancy. However, reviser-composition.md also has a "Context Injection" section. All three describe the same `CONTEXT_CONTENT` variable but only pre-processing.md actually loads it. The consistency of the variable reference is correct. | No action needed. Noted as verified — proper separation of loading (pre-processing) and injection (scorer/reviser composition). | low |
| EV-10 | SKILL.md Architecture flowchart vs `rules/freeform-pipeline.md` | 2 | INCOMPLETE | P2 | The SKILL.md flowchart labels Phase 0 steps using short names (P0A-P0G) but the freeform-pipeline.md uses P0.1-P0.5g labels. The mapping is not documented — an agent executing the flowchart would need to cross-reference two different naming schemes. | Add a label mapping table in SKILL.md Phase 0 section, or align the flowchart labels with freeform-pipeline.md labels. | medium |

### Layer 3 Findings

| # | file_path | layer | category | severity | description | fix_suggestion | confidence |
|---|-----------|-------|----------|----------|-------------|----------------|------------|
| EV-T01 | SKILL.md flowchart vs `rules/freeform-pipeline.md` P0.5 | 3 | TIMING | P2 | The flowchart shows the P0.5 pre-revision sequence as a linear chain: BASELINE_SCORE → SAVE_SNAP → FORMAT → FORMAT_OK → SYNTH → PREREV → PREREV_OK → TAG → INCR → DISPATCH. However, freeform-pipeline.md P0.5e (Invoke Reviser) and P0.5f (Tag Modified Paragraphs) have a dependency: tagging depends on the reviser completing and producing output. The flowchart correctly represents this with sequential arrows. No actual timing issue found — the flowchart is consistent with the pipeline. | No action needed. Verified correct. | low |

---

## 8. Validity Verification

### run-tests `rules/env-check.md` Playwright Hardcoding (Known Issue)

**Finding**: RT-01 confirms the known issue from the proposal. `rules/env-check.md` line 49 hardcodes `npx playwright install` in the Web surface environment check, directly contradicting SKILL.md's profile-agnostic design.

**Severity**: P1 (per proposal definition — "会导致行为偏差但不会完全阻断")

**Verdict**: **REPRODUCED**. The audit successfully identifies this as a P1 CONFLICT between SKILL.md and rules/env-check.md.

---

## 9. Summary Statistics

| Metric | Value |
|--------|-------|
| Skills audited | 7 |
| Total SKILL.md files read | 7 |
| Total associated files read | 45+ |
| Layer 2 findings | 26 |
| Layer 3 findings | 2 (1 confirmed, 1 verified-correct) |
| CONFLICT | 9 (1 P1, 0 P0) |
| INCOMPLETE | 11 |
| REDUNDANT | 3 |
| TIMING (confirmed) | 1 (P3) |
| High confidence | 13 |
| Medium confidence | 10 |
| Low confidence | 8 |

### Severity Distribution

| Severity | Count | Finding IDs |
|----------|-------|-------------|
| P1 | 1 | RT-01 |
| P2 | 13 | RT-02, RT-03, RT-05, WP-01, WP-03, BS-02, GTS-01, GTS-02, GTS-03, BT-01, BT-03, TD-02, TD-03, TD-04, TD-05, EV-03, EV-10 |
| P3 | 14 | RT-04, RT-06, RT-T01, WP-02, WP-04, BS-01, BS-03, GTS-04, BT-02, TD-01, EV-02, EV-04, EV-05, EV-06, EV-07, EV-08, EV-09, EV-T01 |

### Category Distribution

| Category | Count |
|----------|-------|
| CONFLICT | 9 |
| INCOMPLETE | 11 |
| REDUNDANT | 3 |
| TIMING | 2 |

### Per-Skill Finding Count

| Skill | CONFLICT | INCOMPLETE | REDUNDANT | TIMING | Total |
|-------|----------|------------|-----------|--------|-------|
| run-tests | 2 | 2 | 0 | 1 | 5 (+ 1 low-flagged no-issue) |
| write-prd | 1 | 2 | 1 | 0 | 4 |
| brainstorm | 0 | 2 | 1 | 0 | 3 |
| gen-test-scripts | 1 | 2 | 1 | 0 | 4 |
| breakdown-tasks | 1 | 2 | 0 | 0 | 3 |
| tech-design | 2 | 2 | 1 | 0 | 5 |
| eval | 1 | 6 | 0 | 1 | 8 |

---

## 10. Cross-Skill Patterns

### Pattern 1: Convention Loading Method Inconsistency (GTS-01, BT-01, TD-05)

Three skills (gen-test-scripts, breakdown-tasks, tech-design) use `domains` frontmatter filtering for Convention file loading in their SKILL.md Step 0. However, gen-test-scripts' own `rules/convention-guide.md` explicitly forbids this approach with a HARD-RULE. This is a systematic cross-skill inconsistency likely introduced when the Convention system was designed (gen-test-scripts rules were updated but the SKILL.md files of all three skills were not).

### Pattern 2: Intent-Aware Checks Missing (WP-01, TD-03, TD-04)

Multiple rule files for write-prd and tech-design assume `new-feature` intent and do not have conditional logic for `refactor`/`cleanup` intent branches. The SKILL.md files correctly branch on intent, but the rules they reference (self-check, design-quality-checks) do not.

### Pattern 3: Knowledge Extraction Duplication (TD-02)

write-prd and tech-design have nearly identical knowledge-extraction.md files. This is a maintainability concern — changes to the extraction logic need to be synchronized across two files.
