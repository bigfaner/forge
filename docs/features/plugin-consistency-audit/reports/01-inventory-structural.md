# Report 01: Component Inventory and Layer 1 Structural Scan

**Baseline commit**: `08327e1598253ec6fe28a587fb9f0ad19b999cfa`
**Date**: 2026-05-30
**Scope**: 21 skills + 18 commands + 1 agent + hooks/guide.md

---

## 1. Component Inventory

### 1.1 Skills (21)

| # | Skill | templates | rules | data | examples | types | experts | rubrics | Total Files |
|---|-------|-----------|-------|------|----------|-------|---------|---------|-------------|
| 1 | brainstorm | 1 | 2 | - | - | - | - | - | 3 |
| 2 | breakdown-tasks | 3 | 4 | - | - | - | - | - | 7 |
| 3 | clean-code | 1 | - | - | - | - | - | - | 1 |
| 4 | consolidate-specs | 6 | 6 | - | - | - | - | - | 12 |
| 5 | deep-research | 1 | 1 | - | - | - | - | - | 2 |
| 6 | eval | - | 11 | - | - | - | 15 | 12 | 38* |
| 7 | extract-design-md | 3 | 3 | - | - | - | - | - | 6 |
| 8 | forensic | 1 | 1 | - | - | - | - | - | 2 |
| 9 | gen-contracts | 2 | 6 | - | - | - | - | - | 8 |
| 10 | gen-journeys | 1 | 6 | - | - | - | - | - | 7* |
| 11 | gen-sitemap | 1 | 3 | - | - | 1 (.yaml) | - | - | 4 |
| 12 | gen-test-scripts | - | 5 | - | - | 6 | - | - | 11 |
| 13 | init-justfile | 6 (.just) | 2 + 5 surfaces | - | - | - | - | - | 13 |
| 14 | learn | 3 | - | - | - | - | - | - | 3 |
| 15 | quick-tasks | 3 | - | - | - | - | - | - | 3 |
| 16 | run-tests | 1 | 5 + 5 surfaces | - | - | - | - | - | 11 |
| 17 | submit-task | - | - | 6 | - | - | - | - | 6 |
| 18 | tech-design | 6 | 3 | - | 2 | - | - | - | 11 |
| 19 | test-guide | 1 | 4 | - | - | - | - | - | 5 |
| 20 | ui-design | 6 + 3 platforms + 7 styles | 2 | - | - | - | - | - | 18 |
| 21 | write-prd | 4 | 3 | - | 3 | - | - | - | 10 |

\* eval counts experts/ and rubrics/ subdirectories; gen-journeys includes cross-skill reference to gen-contracts/rules/journey-contract-model.md

**Total skill files**: 21 SKILL.md + 187 supporting files = **208** (excluding SKILL.md itself)

#### Detailed File Lists Per Skill

**brainstorm** (3 files)
- `templates/proposal.md`
- `rules/challenge-protocol.md`
- `rules/sc-consistency.md`

**breakdown-tasks** (7 files)
- `templates/manifest-update-tasks.md`
- `templates/task.md`
- `templates/task-doc.md`
- `rules/db-schema.md`
- `rules/existing-code-split.md`
- `rules/phase-detection.md`
- `rules/ui-placement.md`

**clean-code** (1 file)
- `templates/summary.md`

**consolidate-specs** (12 files)
- `templates/biz-specs.md`
- `templates/commit-messages.md`
- `templates/markers.md`
- `templates/review-choices.md`
- `templates/tech-specs.md`
- `templates/vocabulary-index.md`
- `rules/constraints.md`
- `rules/domain-frontmatter.md`
- `rules/drift-detection.md`
- `rules/overlap-detection.md`
- `rules/spec-classification.md`
- `rules/vocabulary-generation.md`

**deep-research** (2 files)
- `templates/research-report.md`
- `rules/research-dimensions.md`

**eval** (38 files)
- Rules (11): `freeform-expert-persistence.md`, `freeform-injection.md`, `freeform-pipeline.md`, `pre-processing.md`, `report-format.md`, `reviser-composition.md`, `rubric-context.md`, `rubric-reference.md`, `scorer-composition.md`, `validate-ux-pipeline.md`
- Experts (15): `experts/freeform/expert-inference.md`, `experts/freeform/expert-template.md`, `experts/freeform/extraction-prompt.md`, `experts/freeform/freeform-reviewer.md`, `experts/freeform/freeform-review-protocol.md`, `experts/protocol/reviser-protocol.md`, `experts/protocol/scorer-protocol.md`, `experts/scorer/architect.md`, `experts/scorer/code-reviewer.md`, `experts/scorer/cto.md`, `experts/scorer/editor.md`, `experts/scorer/pm.md`, `experts/scorer/qa.md`, `experts/scorer/ux-auditor.md`, `experts/scorer/ux-engineer.md`
- Rubrics (12): `consistency.md`, `contract.md`, `design.md`, `journey.md`, `prd.md`, `proposal.md`, `ui-mobile.md`, `ui-tui.md`, `ui-web.md`, `validate-code.md`, `validate-ux.md`

**extract-design-md** (6 files)
- `templates/design-mobile.md`
- `templates/design-tui.md`
- `templates/design-web.md`
- `rules/extraction-layers.md`
- `rules/match-strategy.md`
- `rules/platform-routing.md`

**forensic** (2 files)
- `templates/report.md`
- `rules/deviation-categories.md`

**gen-contracts** (8 files)
- `templates/contract.md`
- `templates/outcome-block.md`
- `rules/code-reconnaissance.md`
- `rules/dimension-rules.md`
- `rules/journey-contract-model.md`
- `rules/risk-density.md`
- `rules/tui-async.md`
- `rules/validation.md`

**gen-journeys** (6 files)
- `templates/journey.md`
- `rules/surface-api.md`
- `rules/surface-cli.md`
- `rules/surface-mobile.md`
- `rules/surface-tui.md`
- `rules/surface-web.md`

**gen-sitemap** (4 files)
- `templates/test-config.yaml`
- `rules/merge-validation.md`
- `rules/page-exploration.md`
- `rules/schema.md`

**gen-test-scripts** (11 files)
- `rules/convention-guide.md`
- `rules/quality-gates.md`
- `rules/run-to-learn.md`
- `rules/step-0.5-validation.md`
- `rules/step-1-contract-loading.md`
- `types/_shared.md`
- `types/api.md`
- `types/cli.md`
- `types/mobile.md`
- `types/tui.md`
- `types/ui.md`

**init-justfile** (13 files)
- `templates/generic.just`
- `templates/go.just`
- `templates/mixed.just`
- `templates/node.just`
- `templates/python.just`
- `templates/rust.just`
- `rules/project-detection.md`
- `rules/self-correction.md`
- `rules/surfaces/api.md`
- `rules/surfaces/cli.md`
- `rules/surfaces/mobile.md`
- `rules/surfaces/tui.md`
- `rules/surfaces/web.md`

**learn** (3 files)
- `templates/convention-entry.md`
- `templates/decision-entry.md`
- `templates/lesson-entry.md`

**quick-tasks** (3 files)
- `templates/manifest-quick.md`
- `templates/task.md`
- `templates/task-doc.md`

**run-tests** (11 files)
- `templates/test-report.md`
- `rules/confidence.md`
- `rules/env-check.md`
- `rules/failure-diagnosis.md`
- `rules/result-parsing.md`
- `rules/test-isolation.md`
- `rules/surfaces/api.md`
- `rules/surfaces/cli.md`
- `rules/surfaces/mobile.md`
- `rules/surfaces/tui.md`
- `rules/surfaces/web.md`

**submit-task** (6 files)
- `data/record-format-coding.md`
- `data/record-format-doc.md`
- `data/record-format-eval.md`
- `data/record-format-gate.md`
- `data/record-format-test.md`
- `data/record-format-validation.md`

**tech-design** (11 files)
- `templates/api-handbook.md`
- `templates/decision-entry.md`
- `templates/er-diagram.md`
- `templates/manifest-update-design.md`
- `templates/schema.sql`
- `templates/tech-design.md`
- `rules/decision-archiving.md`
- `rules/design-quality-checks.md`
- `rules/knowledge-extraction.md`
- `examples/ask-question.md`
- `examples/exploration.md`

**test-guide** (5 files)
- `templates/convention-template.md`
- `rules/convention-structure.md`
- `rules/draft-generation.md`
- `rules/pattern-extraction.md`
- `rules/signal-detection.md`

**ui-design** (18 files)
- `templates/manifest-update-ui.md`
- `templates/prototype.md`
- `templates/ui-design.md`
- `templates/platforms/mobile.md`
- `templates/platforms/tui.md`
- `templates/platforms/web.md`
- `templates/styles/apple.md`
- `templates/styles/minimal-ascii-tui.md`
- `templates/styles/modern-dark-tui.md`
- `templates/styles/shadcn.md`
- `templates/styles/stripe.md`
- `templates/styles/tailwind-ui.md`
- `templates/styles/vercel.md`
- `rules/style-selection.md`
- `rules/tui-panel-requirements.md`

**write-prd** (10 files)
- `templates/manifest.md`
- `templates/prd-spec.md`
- `templates/prd-ui-functions.md`
- `templates/prd-user-stories.md`
- `rules/knowledge-extraction.md`
- `rules/self-check.md`
- `rules/ui-functions.md`
- `examples/ask-questions.md`
- `examples/propose-approaches.md`
- `examples/user-stories.md`

### 1.2 Commands (18)

| # | Command | File |
|---|---------|------|
| 1 | clean-code | `commands/clean-code.md` |
| 2 | eval-consistency | `commands/eval-consistency.md` |
| 3 | eval-contract | `commands/eval-contract.md` |
| 4 | eval-design | `commands/eval-design.md` |
| 5 | eval-journey | `commands/eval-journey.md` |
| 6 | eval-prd | `commands/eval-prd.md` |
| 7 | eval-proposal | `commands/eval-proposal.md` |
| 8 | eval-ui | `commands/eval-ui.md` |
| 9 | execute-task | `commands/execute-task.md` |
| 10 | extract-design-md | `commands/extract-design-md.md` |
| 11 | fix-bug | `commands/fix-bug.md` |
| 12 | gen-sitemap | `commands/gen-sitemap.md` |
| 13 | git-checkout | `commands/git-checkout.md` |
| 14 | git-commit | `commands/git-commit.md` |
| 15 | init-forge | `commands/init-forge.md` |
| 16 | quick | `commands/quick.md` |
| 17 | run-tasks | `commands/run-tasks.md` |
| 18 | simplify-skill | `commands/simplify-skill.md` |

Commands are single `.md` files with YAML frontmatter (name, description, allowed-tools, argument-hint). Internal file references are cross-skill references only (e.g., `fix-bug.md` references `learn/templates/decision-entry.md`).

### 1.3 Agent (1)

| Agent | File |
|-------|------|
| task-executor | `agents/task-executor.md` |

The agent definition contains YAML frontmatter (name, description, model, color, memory, inputs) and inline execution protocol. No external file references within the plugin directory.

### 1.4 Hooks

| File | Type | Description |
|------|------|-------------|
| `hooks/guide.md` | Documentation | Forge directory conventions, CLI reference, terminology |
| `hooks/hooks.json` | Configuration | Hook event definitions (SessionStart, SubagentStart, SessionEnd) |
| `hooks/run-hook.cmd` | Script | Hook dispatcher script |
| `hooks/session-start` | Script (bash) | Session start hook handler |
| `hooks/debug` | Script (bash) | Debug hook handler |

**Note**: `guide.md` does not reference `hooks.json`, `run-hook.cmd`, `session-start`, or `debug` scripts. It functions as a standalone reference document for directory conventions and terminology, not as a hook script index.

---

## 2. Layer 1: Structural Integrity Check

### Methodology

1. Extract all file path references from each SKILL.md using regex pattern matching for `templates/`, `rules/`, `data/`, `examples/`, `types/`, `experts/`, `rubrics/` prefixed paths
2. Compare referenced paths against actual file system listing
3. Classify discrepancies as:
   - **REFERENCE**: Path referenced in SKILL.md but file does not exist on disk
   - **ORPHAN**: File exists on disk but is not referenced in SKILL.md

### 2.1 REFERENCE Issues

| # | Component | Referenced Path | Layer | Category | Severity | Description | Fix Suggestion |
|---|-----------|----------------|-------|----------|----------|-------------|----------------|
| R-01 | gen-journeys | `rules/journey-contract-model.md` | 1 | REFERENCE | P2 | SKILL.md line 20 references `gen-contracts/rules/journey-contract-model.md` via cross-skill path. The file exists in gen-contracts but NOT in gen-journeys/rules/. The SKILL.md text says "resolve relative to the skills parent directory", which makes this a valid cross-skill reference — however, it uses a `rules/` prefix that does not exist locally, which is confusing. | Clarify: the SKILL.md explicitly states "resolve relative to the skills parent directory", so this resolves correctly to gen-contracts/rules/journey-contract-model.md. No file missing, but the indirection pattern is fragile. Consider adding a note in the File Index section. |
| R-02 | gen-test-scripts | `rules/test-isolation.md` | 1 | REFERENCE | P2 | SKILL.md line 214 references `rules/test-isolation.md (located in the run-tests skill directory, resolve relative to the skills parent directory)`. File does not exist in gen-test-scripts/rules/; it exists at run-tests/rules/test-isolation.md. Cross-skill reference is explicit. | Same as R-01: explicit cross-skill reference resolves correctly, but the `rules/` prefix is misleading. |
| R-03 | run-tests | `rules/error-reporting.md` | 1 | REFERENCE | **Not a real issue** | Initial regex extraction flagged this, but context review shows the reference is to `docs/business-rules/error-reporting.md` (a user-project path, not a skill-internal path). No file missing. | N/A — false positive from regex extraction. |
| R-04 | quick-tasks | `data/coding-enhancement.md` | 1 | REFERENCE | **Not a real issue** | Initial regex extraction flagged this, but context review shows the reference is to `forge-cli/pkg/prompt/data/coding-enhancement.md` (a Go source code path, not a skill-internal path). No data/ directory exists for quick-tasks, which is correct. | N/A — false positive from regex extraction. |

**True REFERENCE issues**: 0 confirmed (all initial hits resolved to either cross-skill references with explicit resolution instructions or user-project/codebase paths, not skill-internal paths).

### 2.2 ORPHAN Issues

| # | Component | File Path | Layer | Category | Severity | Description | Fix Suggestion |
|---|-----------|-----------|-------|----------|----------|-------------|----------------|
| O-01 | eval | `rules/validate-ux-pipeline.md` | 1 | ORPHAN | P2 | File exists in rules/ but is not referenced in SKILL.md directly. It is referenced indirectly via `rules/pre-processing.md` line 9 ("Full sub-pipeline: `rules/validate-ux-pipeline.md`"). SKILL.md references pre-processing.md, which chains to this file. | Not a true orphan — it is a second-level reference (SKILL.md -> pre-processing.md -> validate-ux-pipeline.md). Consider adding to SKILL.md file index for discoverability. |
| O-02 | eval | All 15 `experts/` files | 1 | ORPHAN | P2 | 15 expert files exist under experts/ but none are referenced directly in SKILL.md. They are referenced indirectly via rules files (freeform-pipeline.md, scorer-composition.md, freeform-expert-persistence.md). | Not a true orphan — they are second/third-level references chained through rules/. Consider adding a summary reference in SKILL.md for discoverability. |
| O-03 | eval | All 12 `rubrics/` files | 1 | ORPHAN | P3 | 12 rubric files exist but SKILL.md uses a pattern reference `rubrics/<type>.md` rather than listing each file individually. The pattern resolves correctly at runtime. | Not a true orphan — pattern reference covers all files. Low severity. |
| O-04 | consolidate-specs | `templates/vocabulary-index.md` | 1 | ORPHAN | P2 | File exists in templates/ but is not referenced in SKILL.md. It is referenced indirectly via `rules/vocabulary-generation.md` line 39 ("Use the output template from `templates/vocabulary-index.md`"). | Not a true orphan — second-level reference. Consider adding to SKILL.md file index. |
| O-05 | init-justfile | All 6 `templates/*.just` files | 1 | ORPHAN | P1 | 6 `.just` template files exist in templates/ but SKILL.md never references them explicitly. The SKILL.md says "generate from Convention knowledge and LLM understanding" (line 191), which implies the .just files may be dead code, OR they are used by a mechanism not documented in SKILL.md. | Investigate: if these templates are actively loaded by the skill at runtime, add explicit reference in SKILL.md. If they are dead code from a previous design, remove them. This is P1 because the disconnect between "generate from LLM knowledge" and "6 template files exist" creates ambiguity about the skill's actual behavior. |
| O-06 | tech-design | `examples/ask-question.md` | 1 | ORPHAN | P2 | File exists in examples/ but is not referenced anywhere in SKILL.md. | Determine if this is an active example that should be referenced or dead code. If active, add reference. |
| O-07 | tech-design | `examples/exploration.md` | 1 | ORPHAN | P2 | File exists in examples/ but is not referenced anywhere in SKILL.md. | Same as O-06. |
| O-08 | test-guide | `templates/convention-template.md` | 1 | ORPHAN | P2 | File exists in templates/ but is not referenced in SKILL.md. | Determine if this template is actively used and add reference if so. |
| O-09 | ui-design | All 7 `templates/styles/*.md` files | 1 | ORPHAN | P3 | 7 style template files exist but are not referenced in SKILL.md directly. They are referenced via `rules/style-selection.md` (lines 32, 69). | Not a true orphan — second-level reference via rules/style-selection.md. Consider adding summary reference in SKILL.md. |
| O-10 | ui-design | `templates/platforms/{web,mobile,tui}.md` | 1 | ORPHAN | P3 | Platform files are referenced in SKILL.md line 74 using brace expansion syntax `templates/platforms/{web,mobile,tui}.md`. The regex extraction captured this correctly. | Not an orphan — brace expansion pattern resolves to existing files. |

**True ORPHAN issues requiring attention**: O-05 (init-justfile templates), O-06/O-07 (tech-design examples), O-08 (test-guide template)

**Second-level references (not true orphans)**: O-01, O-02, O-03, O-04, O-09, O-10

### 2.3 Cross-Skill References

| Source Skill | Target Skill/File | Reference Type | Valid? |
|-------------|-------------------|----------------|--------|
| gen-journeys | `gen-contracts/rules/journey-contract-model.md` | Cross-skill | Yes (explicit resolution) |
| gen-test-scripts | `run-tests/rules/test-isolation.md` | Cross-skill | Yes (explicit resolution) |
| gen-contracts | `gen-journeys/rules/surface-<type>.md` | Cross-skill | Yes (explicit resolution) |
| fix-bug (command) | `learn/templates/decision-entry.md`, `learn/templates/lesson-entry.md` | Cross-skill | Yes |
| hooks/guide.md | `docs/reference/test-type-model.md` | User-project | Yes (file exists) |

---

## 3. Summary Statistics

| Metric | Value |
|--------|-------|
| Total skills | 21 |
| Total commands | 18 |
| Total agents | 1 |
| Hook files | 5 |
| Total auditable files | 220 (208 skill files + 18 commands + 1 agent + 5 hooks - 12 SKILL.md = 208 supporting files) |
| REFERENCE issues (true) | 0 |
| ORPHAN issues (true) | 3 (O-05, O-06/O-07, O-08) |
| ORPHAN issues (second-level) | 6 (O-01 through O-04, O-09, O-10) |
| Cross-skill references | 5 |
| Baseline commit | `08327e1598253ec6fe28a587fb9f0ad19b999cfa` |

---

## 4. Issue Priority Summary

### P1 (High)

| ID | Component | Description |
|----|-----------|-------------|
| O-05 | init-justfile | 6 `.just` template files exist but SKILL.md says to generate recipes from LLM knowledge, creating ambiguity about skill behavior |

### P2 (Medium)

| ID | Component | Description |
|----|-----------|-------------|
| O-01 | eval | `rules/validate-ux-pipeline.md` not in SKILL.md file index (second-level reference) |
| O-02 | eval | 15 `experts/` files not in SKILL.md file index (second-level references) |
| O-04 | consolidate-specs | `templates/vocabulary-index.md` not in SKILL.md file index (second-level reference) |
| O-06 | tech-design | `examples/ask-question.md` not referenced in SKILL.md |
| O-07 | tech-design | `examples/exploration.md` not referenced in SKILL.md |
| O-08 | test-guide | `templates/convention-template.md` not referenced in SKILL.md |

### P3 (Low)

| ID | Component | Description |
|----|-----------|-------------|
| O-03 | eval | 12 rubrics use pattern reference instead of explicit listing |
| O-09 | ui-design | 7 style templates referenced only via rules/style-selection.md |
| O-10 | ui-design | Platform files use brace expansion pattern |
