---
created: 2026-05-16
author: faner
status: Draft
---

# Proposal: gen-test-cases Per-Type Dispatch

## Problem

gen-test-cases uses a 271-line monolithic SKILL.md that loads ALL interface type instructions (UI, TUI, Mobile, API, CLI) into agent context simultaneously, regardless of which types the project actually uses. A CLI-only project still processes the full UI/TUI/Mobile generation rules, wasting context and diluting focus.

The recent overhaul (6125c64) improved type handling *within* the monolithic framework — TUI/Mobile are now first-class types with dedicated rubric criteria, Interface Accuracy uses percentage-based point splitting, and Antipattern Prevention replaced the old Test Code Quality dimension. However, the fundamental architecture remains: one SKILL.md, one template, one rubric, one output file for all types.

### Evidence

- **Context waste**: A project with only `cli`+`api` capabilities still loads 271 lines of SKILL.md containing full UI/TUI/Mobile generation instructions (Steps 3-4, target derivation rules, classification indicators, integration TC generation) into agent context
- **Diluted evaluation**: The Interface Accuracy dimension (150 pts) now uses percentage-based splitting with type-specific sub-criteria, but this is a *scoring workaround* — it still evaluates all types from a single monolithic rubric with conditional branches rather than giving each type its own dedicated rubric
- **Single-file bottleneck**: All types funnel into one `testing/test-cases.md` — gen-test-scripts already supports `--type` filtering, but must parse the entire file and group by type, then discard non-matching groups. Per-type files would eliminate this grouping step
- **Template rigidity**: The single template has placeholder sections for all 5 types — a CLI-only project generates `test-cases.md` with 4 empty sections (UI, TUI, Mobile, API) plus the CLI section

### Urgency

The monolithic architecture limits the quality gains from per-type specialization. Even with percentage-based rubric splitting, the scorer must navigate a single rubric with conditional branches for each type rather than evaluating a focused, type-specific rubric. This adds cognitive overhead to the scoring loop and reduces scoring precision for niche types.

## Proposed Solution

Split gen-test-cases into a **dispatcher + per-type instruction files** architecture. SKILL.md handles shared setup (profile resolution, PRD reading, AC extraction, interface detection), then loops through each active type loading its dedicated instruction file. Each type gets its own template and the eval skill gets its own rubric.

Output changes from a single `test-cases.md` to per-type files (`ui-test-cases.md`, `api-test-cases.md`, etc.) plus a `testing/manifest.md` aggregator.

Both `gen-test-cases` and `gen-test-scripts` will **on-demand load** conventions from `docs/conventions/` based on the active interface type (see Convention Loading Design below). This ensures generated artifacts conform to project-specific coding standards without baking conventions into the skill instructions.

The `eval-test-cases` command dispatches per-type evaluation via the existing `eval` skill with custom per-type rubrics (`--type test-cases-ui`, `--type test-cases-api`, etc.). The eval skill already executes its scorer→gate→revise loop in the main session (per its Iron Laws), so no architectural bypass is needed — just new rubric files, 5 new location entries in eval's prerequisite table, and a thin dispatcher loop in the eval-test-cases command.

### Convention Loading Design

**Loading point**: The SKILL.md dispatcher loads conventions immediately after Step 2.5 (interface detection), before entering the per-type loop. Each per-type instruction file declares which convention files it requires via a frontmatter `conventions` field. gen-test-scripts follows the same pattern.

**Type-to-file mapping** (by naming convention, not hardcoded paths):

Each per-type instruction file declares its convention dependencies in frontmatter:

```yaml
# types/ui.md frontmatter example
conventions:
  - testing-ui.md        # type-specific testing conventions
  - frontend.md          # domain conventions (optional)
```

**Resolution algorithm**:
1. Read the active type's instruction file frontmatter → extract `conventions` list
2. For each filename, check `docs/conventions/{filename}` exists
3. If exists → load into context via Read tool
4. If missing → skip silently (no warning, no abort)
5. Always load project-wide conventions declared in the dispatcher's own `conventions` frontmatter field (e.g., `conventions: [testing-isolation.md]` in SKILL.md frontmatter)

**Current state**: Existing convention files (`error-handling.md`, `profile-system.md`, `testing-isolation.md`) are project-wide. Type-specific conventions (`testing-ui.md`, `testing-cli.md`, etc.) may not exist yet — the graceful skip handles this. The naming convention establishes the contract for future convention files created via `/consolidate-specs`.

### manifest.md Schema

Generated by the SKILL.md dispatcher after the per-type loop completes. Structure:

```yaml
---
feature: "{{FEATURE_SLUG}}"
types: [ui, api]  # active types that were generated
generated: "{{DATE}}"
---

# Test Cases Manifest: {{FEATURE_SLUG}}

## Summary

| Type | File | Count |
|------|------|-------|
| UI   | testing/ui-test-cases.md  | {{UI_COUNT}}   |
| API  | testing/api-test-cases.md | {{API_COUNT}}  |
| **Total** | | **{{TOTAL}}** |

## Cross-Type Traceability

| TC ID | Source | Type | Target | Priority | File |
|-------|--------|------|--------|----------|------|
| TC-001 | Story 1 / AC-1 | UI | ui/login | P0 | ui-test-cases.md |
| TC-005 | Spec 3.2 | API | api/auth | P1 | api-test-cases.md |
```

This file serves as the single entry point for downstream skills to discover all per-type test case files.

### Innovation Highlights

This follows the **strategy pattern** — the dispatcher selects the appropriate generation strategy at runtime based on detected interface types. The key insight is that the split boundary aligns with the existing Step 2.5/3 boundary: Steps 0-2.5 produce type-agnostic data (profile, PRD content, AC list, detected types), while Steps 3-4 are already type-specific in the current code. The decomposition is mechanical, not architectural.

The recent overhaul's improvements (percentage-based splitting, TUI/Mobile first-class criteria, Antipattern Prevention) provide the foundation for clean extraction — each type's sub-criteria in the current rubric map directly to its dedicated per-type rubric.

## Requirements Analysis

### Key Scenarios

- **Single type**: Project is a CLI tool. SKILL.md detects `{CLI}`, loads `types/cli.md`, generates `cli-test-cases.md`
- **Multiple types**: Web app with UI + API. SKILL.md detects `{UI, API}`, loads `types/ui.md` then `types/api.md` sequentially, generates `ui-test-cases.md` + `api-test-cases.md` + `manifest.md`
- **Mixed + TUI**: Terminal tool with TUI + CLI. SKILL.md detects `{TUI, CLI}`, loads `types/tui.md` then `types/cli.md`
- **Profile missing**: SKILL.md aborts with profile resolution prompt (unchanged behavior)
- **Type not in profile**: Detected type absent from profile capabilities → skip generation for that type (unchanged behavior)

### Non-Functional Requirements

- **Token efficiency**: Per-type instructions should replace the equivalent type-specific sections in the current monolithic SKILL.md, reducing loaded context for single-type projects
- **Backward compatibility**: Existing features with `testing/test-cases.md` should continue to work (gen-test-scripts reads both old and new formats)
- **Extensibility**: Adding a new interface type requires only a new `types/{type}.md` + template + rubric, no changes to the dispatcher
- **Convention awareness**: Per-type instructions and gen-test-scripts load relevant `docs/conventions/` files on demand, grounding output in project standards

### Constraints & Dependencies

- `gen-test-scripts` must be updated: prerequisite check to accept either per-type files or legacy `test-cases.md`, input discovery to read both formats, Step Actionability gate path updated
- `eval-test-cases` must be updated to dispatch per-type rubrics
- `eval` skill's prerequisite/location table must be extended with 5 new `test-cases-*` type entries (each mapping to `testing/` directory, same as existing `test-cases`)
- The existing `test-cases.md` rubric must be decomposed into 5 type-specific rubrics with dimension adaptation guidance
- `testing/manifest.md` is a new artifact that downstream skills must understand

## Alternatives & Industry Benchmarking

### Industry Solutions

Test generation tools like Playwright Codegen, Postman, and Maestro each specialize in a single interface type. The monolithic approach is unusual — most frameworks are type-specific by nature.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero migration cost | Continued context waste, diluted evaluation precision, single-file bottleneck | Rejected: doesn't solve the problem |
| Template-only split | Compromise | Fewer files, simpler downstream | Monolithic instructions still loaded; doesn't address context waste or evaluation precision | Rejected: half-measure |
| **Per-type dispatch** | Strategy pattern | Type-specialized quality, independent evolution, efficient context | ~15-20 new files, downstream updates needed | **Selected: directly addresses remaining gaps** |

## Feasibility Assessment

### Technical Feasibility

The split boundary at Step 2.5/3 is clean — Steps 0-2.5 produce type-agnostic data (profile, PRD content, AC list, detected types). Steps 3-4 are already type-specific in the current code. The decomposition is mechanical, not architectural.

The overhaul's improvements provide a solid foundation: the rubric already has type-specific sub-criteria under Interface Accuracy (web-ui Route Accuracy, tui Output Assertion Accuracy, mobile Interaction Accuracy, api Contract Accuracy, cli Command Coverage), making extraction into dedicated rubrics straightforward.

### Resource & Timeline

Single-contributor change. The main work is extracting type-specific sections from the monolithic SKILL.md into 5 files, decomposing the rubric, and updating 2 downstream skills. No new dependencies.

### Dependency Readiness

All downstream skills (gen-test-scripts, eval-test-cases) exist and are well-understood. The eval skill already supports parameterized rubric loading — adding 5 new type entries to its location table is a mechanical extension. gen-test-scripts already has `--type` filter infrastructure, making per-type file consumption a natural extension.

## Scope

### In Scope

1. Refactor gen-test-cases `SKILL.md` into dispatcher (Steps 0-2.5 + manifest generation) — target under 150 lines
2. Create 5 per-type instruction files: `types/ui.md`, `types/tui.md`, `types/mobile.md`, `types/api.md`, `types/cli.md` — each containing type-specific Steps 3-4 instructions extracted from the current monolithic SKILL.md
3. Create 5 per-type templates: `templates/ui-test-cases.md`, `templates/tui-test-cases.md`, `templates/mobile-test-cases.md`, `templates/api-test-cases.md`, `templates/cli-test-cases.md` — each containing only that type's section (frontmatter + TC placeholders + type-specific traceability)
4. Define `testing/manifest.md` aggregator structure (generated by SKILL.md after per-type loop)
5. Create 5 per-type rubric files: `eval/rubrics/test-cases-ui.md`, `test-cases-tui.md`, `test-cases-mobile.md`, `test-cases-api.md`, `test-cases-cli.md`. Decomposition guidance: keep 1000-point scale; retain PRD Traceability (200), Step Actionability (250), Completeness (200), Structure & ID (100), Antipattern Prevention (100) as shared dimensions; replace Interface Accuracy (150) with type-specific dimensions derived from current sub-criteria (e.g., UI→Visual State Accuracy from web-ui Route Accuracy criteria, API→Contract Accuracy from api Contract Accuracy criteria, CLI→Output Accuracy from cli Output Assertion criteria). Each per-type rubric gets the full 150 pts for its specialized dimension without percentage-based splitting
6. Add 5 new entries to eval skill's prerequisite/location table: `test-cases-ui`, `test-cases-tui`, `test-cases-mobile`, `test-cases-api`, `test-cases-cli` — each mapping to `testing/` directory with `{type}-test-cases.md` as the document file
7. Refactor `eval-test-cases` command into a thin per-type dispatcher loop: for each active type, invoke `eval` skill with `--type test-cases-{type}` pointing to the matching per-type rubric and `{type}-test-cases.md` file. Fallback: if no per-type files exist, invoke with `--type test-cases` for legacy monolithic mode
8. Update `gen-test-scripts` Step 1 prerequisite check to accept per-type files (`{type}-test-cases.md`) OR legacy `test-cases.md`. Discovery logic: glob `testing/*-test-cases.md` first; if empty, fall back to `testing/test-cases.md`. When reading per-type files, skip the type grouping step (file is already single-type). Step Actionability gate path updated accordingly
9. Add convention loading to SKILL.md dispatcher and gen-test-scripts: read per-type instruction frontmatter `conventions` field, load existing files from `docs/conventions/`, skip missing. Dispatcher's own frontmatter declares project-wide conventions (e.g., `conventions: [testing-isolation.md]`)

### Out of Scope

- Full rewrite of `gen-test-scripts` SKILL.md (only prerequisite check + input discovery + convention loading changes)
- Changes to eval skill core scorer/revise loop logic (reused as-is; only prerequisite/location table extended with 5 new type entries)
- New interface types beyond the existing 5
- Changes to T-quick or T-test task structure in breakdown-tasks
- Changes to `run-e2e-tests` or `graduate-tests`
- Creating type-specific convention files (only the loading mechanism; file creation is `/consolidate-specs`'s responsibility)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Backward incompatibility with existing `test-cases.md` files | M | H | gen-test-scripts reads both old single-file and new per-type formats; eval-test-cases falls back to monolithic rubric if no per-type files found |
| Per-type instructions drift out of sync with shared steps | M | M | Shared steps stay in SKILL.md; per-type files only cover Steps 3-4 |
| Rubric split loses cross-type coverage (e.g., UI+API integration scenarios) | L | M | manifest.md captures cross-type traceability; integration TCs live in the primary type's file with explicit cross-type references |
| Token overhead from loading manifest.md + per-type file | L | L | manifest.md is small (summary + traceability table); per-type instructions replace equivalent monolithic sections |
| eval-test-cases DOC_DIR reads all per-type files at once | M | H | eval-test-cases dispatcher passes single `{type}-test-cases.md` file path to eval skill per invocation, not the entire `testing/` directory |
| Convention files don't exist at launch | H | L | Graceful skip: missing files are silently ignored. The naming convention establishes a contract for future `/consolidate-specs` output |

## Success Criteria

- [ ] SKILL.md is under 150 lines (dispatcher + shared Steps 0-2.5 + manifest generation)
- [ ] Each per-type instruction file covers type-specific Steps 3-4 completely
- [ ] Each per-type rubric has at least 4 dimensions with type-specific scoring criteria (shared dimensions: PRD Traceability 200, Step Actionability 250, Completeness 200, Structure & ID 100, Antipattern Prevention 100; type-specific dimension replaces Interface Accuracy — full 150 pts without percentage-based splitting)
- [ ] gen-test-scripts prerequisite check accepts per-type files via glob `testing/*-test-cases.md`, falls back to legacy `testing/test-cases.md`
- [ ] eval skill's prerequisite/location table has entries for all 5 `test-cases-*` types, each mapping to `testing/` directory
- [ ] eval-test-cases loops per-type: for each active type, invokes eval skill with `--type test-cases-{type}` targeting the matching rubric and single `{type}-test-cases.md` file; falls back to `--type test-cases` for legacy monolithic mode
- [ ] For a CLI-only project, only `types/cli.md` instruction + `cli-test-cases.md` template are loaded (no UI/TUI/Mobile content in context)
- [ ] Convention loading reads per-type instruction frontmatter `conventions` field, loads existing files from `docs/conventions/`, skips missing files silently
- [ ] `testing/manifest.md` is generated by the dispatcher after the per-type loop, following the schema defined in this proposal (frontmatter + summary table + cross-type traceability with file references)

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
