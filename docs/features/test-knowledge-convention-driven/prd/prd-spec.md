---
feature: "test-knowledge-convention-driven"
status: Draft
db-schema: "no"
---

# test-knowledge-convention-driven — PRD Spec

> PRD Spec: defines WHAT the feature is and why it exists.

## Background

### Why (Reason)

Forge's Profile system (6 language profiles + embed + auto-detection) coexists with the newer Journey-Contract test model, creating two problems:

1. **Conceptual conflict**: Contract layer defers precise matching to generation time (semantic descriptors), but Profile layer hard-codes technical decisions upfront (generate.md dictates assertion libraries, import patterns, test style before the LLM sees project code).

2. **Responsibility confusion**: Profile simultaneously handles 3 independent concerns: (a) project tech stack detection, (b) framework-specific code generation knowledge, (c) test organization axis. The Journey-Contract model already solved (c), but (a) and (b) remain as Profile residue.

3. **User impact**: Profile only works correctly for projects using default frameworks (e.g., Go + go-testing). Non-default frameworks (ginkgo, vitest, pytest) get wrong results from auto-detection, and users cannot fix this without a Forge release cycle. 2 framework support requests are currently blocked.

### What (Target)

Replace the Profile system with user-editable Convention files (`docs/conventions/`) that provide framework knowledge to LLM-driven test generation. Add a `/forge:test-guide` slash command for guided Convention creation.

### Who (Users)

Forge CLI users — developers who use forge to generate and manage e2e tests. Two segments:

- **Default framework users**: Profile works for them but they get no benefit from Profile over LLM defaults
- **Non-default framework users**: Profile fails for them (ginkgo, vitest, etc.) and they have no recourse

## Goals

| Goal | Metric | Notes |
|------|--------|-------|
| Enable non-default framework support | Users can use ginkgo/vitest/pytest without Forge code changes | Currently blocked on Forge release |
| First-pass compile rate >= 85% | gen-test-scripts output passes `just e2e-compile` on first attempt | Measured on forge-cli's 126+ existing tests; baseline: current Profile-based generation achieves ~90% on these same tests (forge-cli uses Go testing + testify — the default profile). The 85% target is a floor, not a ceiling — Convention mode should match Profile quality for default frameworks and exceed it for non-default ones where Profile fails entirely |
| Semantic correctness | Generated test code references the correct Contract test steps and validates the declared outcomes | Compile gate verifies syntax; separate Contract-traceability check verifies each generated test function maps to a Contract step with matching assertions |
| Generated code diff equivalence | Framework core patterns (import, assertion, tag) identical to Profile-generated output | Style differences (variable naming, comments) allowed |
| Full Profile removal | Zero imports of `forge-cli/pkg/profile`; verified by `grep -r "pkg/profile" forge-cli/` returning no results | 19+ files rewritten; import audit is an explicit pre-deletion gate |
| Convention bootstrap in < 5 min | `/forge:test-guide` generates a minimal Convention file | User confirms detected patterns |
| Backward compatibility | All 126+ existing e2e tests pass without modification | forge-cli project as validation baseline |

## Scope

### In Scope

- [ ] Remove `pkg/profile/` directory entirely (6 language profiles, embed, detect, framework mapping)
- [ ] Clean config.yaml — remove `languages`, `interfaces`, `test-framework`, `project-type` fields; retain only `auto.*`, `test-command`, `worktree`
- [ ] Rewrite gen-test-scripts skill — Convention loading + Code Reconnaissance + compile gate
- [ ] Rewrite run-e2e-tests skill — just execution + Convention Result Format-driven parsing
- [ ] Rewrite init-justfile skill — Convention + Reconnaissance-driven recipe generation
- [ ] Rewrite Go packages: `pkg/journey/`, `pkg/e2e/`, `pkg/just/`, `pkg/task/`
- [ ] Rewrite CLI commands: `forge config init`, `forge init`, `forge task index`, `forge task add`
- [ ] Remove CLI commands: `forge test detect`, `forge test get`, `forge test interfaces`, `forge test framework`
- [ ] New `/forge:test-guide` slash command for guided Convention file creation
- [ ] Convention file fixed structure definition (Framework / Assertion / Tags / Result Format)
- [ ] Integrate Convention files into consolidate-specs management

### Out of Scope

- gen-journeys skill (no Profile dependency)
- gen-contracts skill (no Profile dependency)
- `forge test verify` / `forge test promote` (no Profile dependency)
- Existing test migration or rewriting
- Anti-pattern documentation generation
- Convention file auto-sync / watch mode
- Unit test coverage targets for rewritten Go code

## Flow Description

### Business Flow Description

**Main flow — existing project with tests:**

1. User runs gen-test-scripts for a Journey
2. Skill loads Convention files matching `domains: [testing, ...]` from `docs/conventions/`
3. Skill runs Code Reconnaissance on existing test files (imports, tags, naming patterns)
4. Skill generates test code using Convention + Reconnaissance + Skill methodology
5. Compile gate: `just e2e-compile` validates generated code
6. If compile fails (max 2 retries) → feed error back to LLM → regenerate
7. If no Convention found → output hint, proceed with LLM defaults + Reconnaissance

**Bootstrap flow — new project (cold start):**

1. User runs init-justfile → generates justfile with e2e-compile/e2e-test/e2e-setup recipes
   - If justfile already exists with e2e recipes → skip, use existing
   - If justfile missing or has no e2e-compile recipe → output warning with manual instructions (fallback: `go test -c ./...` / `python -m pytest --collect-only`)
2. User runs gen-test-scripts → LLM generates with defaults → compile gate validates
   - If `just e2e-compile` not available → gen-test-scripts outputs actionable error: "Missing justfile e2e-compile recipe. Run `/forge:init-justfile` first, or add a recipe manually."
3. User runs `/forge:test-guide` → detects framework signals → generates Convention file
   - If file signals ambiguous (e.g., both go.mod and package.json) → list candidates, user selects
   - If no recognizable signals → ask user to specify language and framework manually
4. Subsequent runs use Convention for consistent output

**Convention creation flow (test-guide):**

1. User invokes `/forge:test-guide`
2. Skill scans project for file signals (go.mod, package.json, etc.) and existing test files
3. If test files exist → extract patterns (imports, tags, naming) → present to user for confirmation
4. If no test files → present candidate frameworks for user selection
5. User confirms → skill writes `docs/conventions/testing-<scope>.md` with minimal structure
   - If write fails (permission, directory missing) → output error with path and recovery hint
   - If existing Convention file found → show diff, ask user to confirm update or keep

### Business Flow Diagram

```mermaid
flowchart TD
    Start([User invokes gen-test-scripts]) --> LoadConv{Convention files\nfound?}
    LoadConv -->|Yes| ValidateConv{Convention\nvalid?}
    LoadConv -->|No| Hint[Output hint:\nNo Convention found]
    ValidateConv -->|Yes| LoadContent[Load Convention content\nmatching domains]
    ValidateConv -->|No| ConvError[Output warning:\nConvention malformed.\nList missing sections.\nProceed with LLM defaults]
    ConvError --> Recon
    LoadContent --> Recon[Code Reconnaissance:\nscan existing test files]
    Hint --> Recon
    Recon --> ReconResult{Reconnaissance\nfound useful signals?}
    ReconResult -->|Yes| Generate[LLM generates test code\nusing Convention + Reconnaissance]
    ReconResult -->|No| GenerateBare[LLM generates with\ndefaults only]
    GenerateBare --> Compile
    Generate --> Compile{just e2e-compile\navailable?}
    Compile -->|Yes| RunCompile[just e2e-compile]
    Compile -->|No| NoJust[Error: missing e2e-compile recipe.\nRun /init-justfile first]
    RunCompile -->|Pass| Done([Success])
    RunCompile -->|Fail| Retry{Retries < 2?}
    Retry -->|Yes| FeedError[Feed compile error\nback to LLM]
    FeedError --> Generate
    Retry -->|No| Blocked([Blocked: compile gate failed.\nRecovery: 1) Check Convention\n2) Run /test-guide\n3) Edit Convention\n4) Re-run gen-test-scripts])
```

### Data Flow Description

N/A — single-system feature (forge-cli only). No cross-system data flows.

## Functional Specs

### FS-1: Convention File Structure

Convention files use fixed sections: `Framework`, `Assertion`, `Tags`, `Result Format` (minimum set). Optional sections: `Helpers`, `Import Patterns`, `Code Style`, `Anti-patterns`. Files are markdown with `domains` frontmatter for selective loading.

**Required section schema:**

| Section | Required Fields | Format | Example |
|---------|----------------|--------|---------|
| Framework | name, file-pattern, package | Bullet list | `Go testing package + testify/assert`, `File pattern: *_test.go`, `Package: e2e` |
| Assertion | library, key functions | Bullet list | `assert (not require)`, `assert.NoError / assert.Contains` |
| Tags | build-tag/marker syntax | Bullet list | `//go:build e2e` |
| Result Format | output-flags, format-type | Bullet list | `Output command flags: -json`, `Format type: json-stream` |

**Validation rules:**
- Missing `domains` frontmatter → skill treats file as non-loadable, output warning
- Missing required section → skill logs warning listing missing sections, proceeds with LLM defaults for that section's area
- Invalid section content (e.g., empty Framework name) → skill logs warning, treats as missing
- Multiple Convention files with overlapping domains → skill merges, last-loaded wins for conflicting fields; log a note about the overlap

### FS-2: Convention Loading Mechanism

Skills load Convention files by listing `docs/conventions/` and filtering by `domains` frontmatter. This uses the existing guide.md convention-loading mechanism — no new infrastructure.

**Error handling:**
- No Convention files found → output hint: "No test Convention files found in docs/conventions/. Generation will use LLM defaults. Run `/forge:test-guide` to create one."
- Convention file unreadable (permissions, encoding) → skip file, log warning with file path
- Convention file parseable but missing required sections → see FS-1 validation rules

### FS-3: Code Reconnaissance Extension

Fact Table (runtime LLM notes, not persisted) extended to collect test framework info: file patterns, import analysis, build tag analysis, function signature patterns. Supplements Convention for project-specific details.

**Reliability expectations:**
- Reconnaissance is best-effort — it may find nothing (no test files, unrecognized patterns)
- When Reconnaissance finds signals that conflict with Convention → Convention wins (user-edited knowledge overrides auto-detected patterns), but skill logs the conflict for user awareness
- When Reconnaissance finds nothing useful → skill proceeds with Convention alone, or LLM defaults if no Convention

**Failure mode:** If no test files exist and no file signals are recognizable, Reconnaissance produces an empty Fact Table. Skill proceeds with LLM defaults only. This is the expected cold-start behavior — not an error condition.

### FS-4: Compile Gate

`just e2e-compile` as the quality gate after generation. Max 2 retries on failure with compile error feedback. No fallback to Profile.

**Prerequisite:** `just e2e-compile` must be available. If the recipe is missing:
- gen-test-scripts outputs actionable error: "Missing justfile e2e-compile recipe. Run `/forge:init-justfile` first, or add a recipe manually."
- Generation is blocked until the recipe exists

**Note on retry semantics:** The compile gate retry (max 2) is a generation-time mechanism within gen-test-scripts — it is NOT part of the quality-gate pipeline defined in BIZ-quality-gate-001. The quality-gate pipeline operates after all tasks complete; the compile gate operates during individual test script generation.

**Recovery on exhaustion (all retries fail):**
1. Output the compile error to the user with the generated file path
2. Suggest recovery actions: (a) check Convention file for incorrect framework/assertion declarations, (b) run `/forge:test-guide` to regenerate Convention from project analysis, (c) manually edit the generated test file to fix compilation
3. Do not auto-delete the generated file — leave it for user inspection

### FS-5: /forge:test-guide Command

Multi-turn conversational skill. Detects project signals (file signatures + existing test patterns). Cold start: presents framework candidates for user selection. Warm start: extracts patterns from existing tests for user confirmation. Writes minimal Convention file.

### FS-6: config.yaml Cleanup

Remove `languages`, `interfaces`, `test-framework`, `project-type`. Retain `auto.*`, `test-command`, `worktree`. Old fields silently ignored on upgrade.

### FS-7: Profile Removal

Delete `pkg/profile/` entirely. Rewrite all 19+ consumers. Remove `forge test detect/get/interfaces/framework` CLI commands.

**Import audit gate:** Before deleting `pkg/profile/`, run `grep -r "pkg/profile" forge-cli/` (or `go vet` with the package path) to verify zero import references remain. This audit is a required pre-deletion step — it catches any consumers not tracked in the initial 19+ count. If the audit finds additional consumers, those must be rewritten before deletion proceeds.

### FS-8: Silent Migration

Existing config.yaml fields silently ignored. No transition period, no coexistence, no migration wizard.

### FS-9: consolidate-specs Integration

Convention files in `docs/conventions/` are automatically visible to consolidate-specs (which already scans this directory). The integration consists of:
- Convention files use standard `domains` frontmatter, making them eligible for drift detection
- consolidate-specs treats Convention files the same as other convention files — no special handling required
- Drift detection: if Convention declares `assert (not require)` but existing tests use `require`, consolidate-specs flags this during a drift audit

### Related Changes

| # | Module | Change Point |
|------|----------|------------|
| 1 | pkg/journey/ | Remove FrameworkInfo dependency; Convention-driven tags |
| 2 | pkg/task/ | Rewrite test task generation; remove StrategyResolver, Languages, TestInterfaces |
| 3 | pkg/e2e/ | Remove Profile calls; test execution via just |
| 4 | pkg/just/ | Remove project-type dependency for scope resolution |
| 5 | internal/cmd/ | Simplify config init, init, task index, task add; remove test subcommands |

## Other Notes

### Performance Requirements

- First-pass compile rate >= 85% (measured on 126+ existing tests)
- Generation time should not increase by more than 20% vs Profile-based generation

### Data Requirements

- No data migration needed
- Convention files are new artifacts in `docs/conventions/`
- Existing config.yaml fields silently ignored

### Monitoring Requirements

- N/A — internal tooling, no production monitoring

### Security Requirements

- N/A — no user-facing API, no sensitive data handling

### Rollback Strategy

Developed on `v3.0.0` branch. Phase 1 and 2 are independently revertible. Point of no return: Phase 3 start (test-guide depends on Profile-free environment).

### Phase Gates

Each phase has an explicit go/no-go criterion that must be met before proceeding:

| Phase | Gate | Go Criteria | No-Go Action |
|-------|------|-------------|-------------|
| Phase 0: POC | POC complete | First-pass compile rate >= 70% on 3 frameworks (Go testing, ginkgo, vitest) | Reconsider approach; fall back to Alternative B (AST-based detection) |
| Phase 1: Profile removal | All code compiles | `go build ./...` succeeds; all existing tests pass; `grep -r "pkg/profile" forge-cli/` returns zero results | Fix remaining imports/compilation errors before proceeding |
| Phase 2: Skill rewrites | Skills produce valid output | gen-test-scripts generates code that passes `just e2e-compile` at >= 85% first-pass rate on forge-cli's 126+ tests; run-e2e-tests parses results correctly; init-justfile generates valid recipes | Iterate on skill prompts; do not proceed to Phase 3 |
| Phase 3: test-guide | End-to-end validation | `/forge:test-guide` generates valid Convention file in < 5 min; Convention file loaded correctly by gen-test-scripts; all 126+ e2e tests pass | Fix test-guide skill; rollback to Phase 2 endpoint if unfixable |

---

## Quality Checklist

- [x] Is the requirement title accurate and descriptive
- [x] Does the background include all three elements: reason, target, users
- [x] Are the goals quantified
- [x] Is the flow description complete
- [x] Does the business flow diagram exist (Mermaid format)
- [x] Are related changes thoroughly analyzed
- [x] Are non-functional requirements considered (performance / data / monitoring / security)
- [x] Are all tables filled completely
- [x] Is there any ambiguous or vague wording
- [x] Is the spec actionable and verifiable
