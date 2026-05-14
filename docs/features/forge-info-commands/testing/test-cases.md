---
feature: "forge-info-commands"
sources:
  - docs/proposals/forge-info-commands/proposal.md
  - docs/features/forge-info-commands/tasks/1-config-commands.md
  - docs/features/forge-info-commands/tasks/2-info-commands.md
  - docs/features/forge-info-commands/tasks/3-init-command.md
  - docs/features/forge-info-commands/tasks/4-migration-resolve-scope.md
generated: "2026-05-14"
---

# Test Cases: forge-info-commands

## Summary

| Type | Count |
|------|-------|
| CLI  | 32    |
| **Total** | **32** |

> **Note**: This feature is a CLI-only project (Go/cobra binary). No UI or HTTP API interfaces exist. All acceptance criteria map to CLI test cases. Profile capabilities `tui` and `api` are not applicable as product interfaces for this feature.

---

## CLI Test Cases

### Config Commands (Task 1)

## TC-001: forge config init — interactive setup
- **Source**: Task 1 / AC-1, AC-2
- **Type**: CLI
- **Target**: cli/config-init
- **Test ID**: cli/config-init/interactive-setup
- **Pre-conditions**: `.forge/config.yaml` does NOT exist; working directory is a forge project
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge config init`
  2. Select `backend` for project-type
  3. Select `go-test` for test-profiles
  4. Confirm capabilities `[tui, api, cli]`
- **Expected**: `.forge/config.yaml` is created with `project-type: backend`, `test-profiles: [go-test]`, `capabilities: [tui, api, cli]`
- **Priority**: P0

## TC-002: forge config init — reconfigure when config exists
- **Source**: Task 1 / AC-3
- **Type**: CLI
- **Target**: cli/config-init
- **Test ID**: cli/config-init/reconfigure-prompt
- **Pre-conditions**: `.forge/config.yaml` already exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge config init`
  2. Verify prompt "Config already exists. Reconfigure? [y/N]" is shown
  3. Enter `n`
- **Expected**: Command exits without modifying existing config
- **Priority**: P0

## TC-003: forge config init — reconfigure accepted
- **Source**: Task 1 / AC-3
- **Type**: CLI
- **Target**: cli/config-init
- **Test ID**: cli/config-init/reconfigure-accepted
- **Pre-conditions**: `.forge/config.yaml` already exists with `project-type: backend`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge config init`
  2. Enter `y` at reconfigure prompt
  3. Select `frontend` for project-type
  4. Complete remaining prompts
- **Expected**: `.forge/config.yaml` is overwritten with new `project-type: frontend`
- **Priority**: P1

## TC-004: forge config get project-type
- **Source**: Task 1 / AC-4
- **Type**: CLI
- **Target**: cli/config-get
- **Test ID**: cli/config-get/project-type
- **Pre-conditions**: `.forge/config.yaml` exists with `project-type: backend`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge config get project-type`
- **Expected**: Output is plain text `backend` (no formatting blocks, no quotes)
- **Priority**: P0

## TC-005: forge config get capabilities — array output
- **Source**: Task 1 / AC-5
- **Type**: CLI
- **Target**: cli/config-get
- **Test ID**: cli/config-get/capabilities-array
- **Pre-conditions**: `.forge/config.yaml` exists with `capabilities: [tui, api, cli]`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge config get capabilities`
- **Expected**: Output is three lines: `tui`, `api`, `cli` (one per line, no formatting)
- **Priority**: P0

## TC-006: forge config get — missing key
- **Source**: Task 1 / AC-6
- **Type**: CLI
- **Target**: cli/config-get
- **Test ID**: cli/config-get/missing-key
- **Pre-conditions**: `.forge/config.yaml` exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge config get nonexistent-key`
- **Expected**: No output to stdout; exit code is 1
- **Priority**: P0

## TC-007: ForgeConfig struct fields
- **Source**: Task 1 / AC-7
- **Type**: CLI
- **Target**: cli/config-struct
- **Test ID**: cli/config-struct/field-validation
- **Pre-conditions**: Go code compiles
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `go vet ./forge-cli/pkg/profile/...`
  2. Verify `ForgeConfig` struct has `ProjectType string`, `TestProfiles []string`, `Capabilities []string` fields
- **Expected**: Struct fields exist with correct types
- **Priority**: P1

### Proposal Commands (Task 2)

## TC-008: forge proposal — list all proposals
- **Source**: Task 2 / AC-1
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/all-proposals
- **Pre-conditions**: `docs/proposals/` directory contains at least one proposal with valid frontmatter
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge proposal`
- **Expected**: Output contains a table with columns: Slug | Created | Status | PRD | Feature. Each row matches a proposal in `docs/proposals/`
- **Priority**: P0

## TC-009: forge proposal — slug detail view
- **Source**: Task 2 / AC-2
- **Type**: CLI
- **Target**: cli/proposal-detail
- **Test ID**: cli/proposal-detail/slug-detail
- **Pre-conditions**: A proposal with slug `forge-info-commands` exists in `docs/proposals/`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge proposal forge-info-commands`
- **Expected**: Output shows metadata (created, author, status), content summary (Problem + Proposed Solution), linked artifacts (PRD/Feature/Task status), and file path
- **Priority**: P0

## TC-010: forge proposal — created date from frontmatter
- **Source**: Task 2 / AC-3
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/created-date-frontmatter
- **Pre-conditions**: A proposal has `created: 2026-05-14` in frontmatter
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge proposal`
  2. Check the Created column for that proposal
- **Expected**: Created column shows `2026-05-14` (from frontmatter, not file system time)
- **Priority**: P1

## TC-011: forge proposal — PRD column checks prd-spec.md
- **Source**: Task 2 / AC-4
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/prd-column
- **Pre-conditions**: Proposal `forge-info-commands` has no `docs/features/forge-info-commands/prd/prd-spec.md`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge proposal`
  2. Check PRD column for `forge-info-commands`
- **Expected**: PRD column shows absence indicator (e.g. `-` or `No`) for proposals without prd-spec.md
- **Priority**: P1

## TC-012: forge proposal — Feature column reads manifest status
- **Source**: Task 2 / AC-5
- **Type**: CLI
- **Target**: cli/proposal-list
- **Test ID**: cli/proposal-list/feature-status
- **Pre-conditions**: Proposal `forge-info-commands` has `docs/features/forge-info-commands/manifest.md` with `status: tasks`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge proposal`
  2. Check Feature column for `forge-info-commands`
- **Expected**: Feature column shows `tasks` (from manifest.md status field)
- **Priority**: P1

### Feature Commands (Task 2)

## TC-013: forge feature list — lists all features
- **Source**: Task 2 / AC-6
- **Type**: CLI
- **Target**: cli/feature-list
- **Test ID**: cli/feature-list/all-features
- **Pre-conditions**: `docs/features/` contains at least one feature with `manifest.md`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge feature list`
- **Expected**: Output contains table with columns: Slug | Status | Progress | PRD(score) | Design(score) | UI(score) | Tests(score). Each row matches a feature in `docs/features/`
- **Priority**: P0

## TC-014: forge feature list — progress from index.json
- **Source**: Task 2 / AC-7
- **Type**: CLI
- **Target**: cli/feature-list
- **Test ID**: cli/feature-list/progress-counts
- **Pre-conditions**: Feature has `tasks/index.json` with 3 completed and 2 pending tasks
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge feature list`
  2. Check Progress column for that feature
- **Expected**: Progress column shows `3/5` (completed/total)
- **Priority**: P1

## TC-015: forge feature list — scores from frontmatter
- **Source**: Task 2 / AC-8
- **Type**: CLI
- **Target**: cli/feature-list
- **Test ID**: cli/feature-list/score-display
- **Pre-conditions**: Feature has no eval scores written to frontmatter
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge feature list`
  2. Check score columns
- **Expected**: Score columns show `—` when no `score` field exists in frontmatter
- **Priority**: P1

## TC-016: forge feature status — detailed view
- **Source**: Task 2 / AC-9
- **Type**: CLI
- **Target**: cli/feature-status
- **Test ID**: cli/feature-status/detail-view
- **Pre-conditions**: Feature `forge-info-commands` exists with manifest, tasks, and artifacts
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge feature status forge-info-commands`
- **Expected**: Output shows manifest summary, task counts by status (pending/in_progress/completed/blocked/skipped/rejected), artifacts with scores, and current in-progress task info
- **Priority**: P0

## TC-017: forge feature — no args keeps existing behavior
- **Source**: Task 2 / Hard Rules
- **Type**: CLI
- **Target**: cli/feature
- **Test ID**: cli/feature/no-args-existing-behavior
- **Pre-conditions**: A current feature is set
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge feature` (no arguments)
- **Expected**: Shows current feature (existing behavior preserved)
- **Priority**: P0

### Lesson Commands (Task 2)

## TC-018: forge lesson — list all lessons
- **Source**: Task 2 / AC-10
- **Type**: CLI
- **Target**: cli/lesson-list
- **Test ID**: cli/lesson-list/all-lessons
- **Pre-conditions**: `docs/lessons/` contains lesson markdown files
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge lesson`
- **Expected**: Output contains table with columns: Name | Created | Tags | Category
- **Priority**: P0

## TC-019: forge lesson — category from filename prefix
- **Source**: Task 2 / AC-11
- **Type**: CLI
- **Target**: cli/lesson-list
- **Test ID**: cli/lesson-list/category-prefix
- **Pre-conditions**: `docs/lessons/` contains files like `gotcha-something.md`, `arch-design.md`, `pattern-retry.md`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge lesson`
  2. Check Category column for each lesson
- **Expected**: `gotcha-something.md` shows `gotcha`, `arch-design.md` shows `arch`, `pattern-retry.md` shows `pattern`
- **Priority**: P1

## TC-020: forge lesson name — detail view
- **Source**: Task 2 / AC-12
- **Type**: CLI
- **Target**: cli/lesson-detail
- **Test ID**: cli/lesson-detail/name-detail
- **Pre-conditions**: A lesson file exists in `docs/lessons/`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge lesson <name>`
- **Expected**: Output shows metadata (created, tags) and file path. Full lesson content is NOT printed
- **Priority**: P0

### Init Command (Task 3)

## TC-021: forge init — creates .forge/ directory
- **Source**: Task 3 / AC-1
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/create-forge-dir
- **Pre-conditions**: `.forge/` directory does NOT exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
  2. Check for `.forge/` directory
- **Expected**: `.forge/` directory exists; output shows `CREATED .forge/`
- **Priority**: P0

## TC-022: forge init — generates CLAUDE.md from template
- **Source**: Task 3 / AC-2
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/generate-claudemd
- **Pre-conditions**: `CLAUDE.md` does NOT exist in project root
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
  2. Check for `CLAUDE.md` in project root
- **Expected**: `CLAUDE.md` exists with content from embedded template; output shows `CREATED CLAUDE.md (from template)`
- **Priority**: P0

## TC-023: forge init — appends to .gitignore with dedup
- **Source**: Task 3 / AC-3
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/gitignore-dedup
- **Pre-conditions**: `.gitignore` exists but does not contain forge runtime entries
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
  2. Read `.gitignore`
- **Expected**: `.gitignore` contains all 5 forge runtime entries (`# Forge runtime`, `docs/features/*/tasks/process/`, `.forge/state.json`, `tests/results/.last-run.json`, etc.); output shows `APPENDED .gitignore (5 entries)`
- **Priority**: P0

## TC-024: forge init — .gitignore dedup skips existing entries
- **Source**: Task 3 / AC-3, Hard Rules
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/gitignore-dedup-existing
- **Pre-conditions**: `.gitignore` already contains `.forge/state.json`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
  2. Read `.gitignore`
- **Expected**: `.forge/state.json` line appears only once (not duplicated)
- **Priority**: P1

## TC-025: forge init — appends justfile recipes with dedup
- **Source**: Task 3 / AC-4
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/justfile-recipes
- **Pre-conditions**: `justfile` exists but has no `claude` or `claude-c` recipes
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
  2. Read `justfile`
- **Expected**: `justfile` contains `claude:` and `claude-c:` recipes; output shows `APPENDED justfile (2 recipes: claude, claude-c)`
- **Priority**: P0

## TC-026: forge init — justfile dedup skips existing recipes
- **Source**: Task 3 / AC-4, Hard Rules
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/justfile-dedup-existing
- **Pre-conditions**: `justfile` already contains a `claude:` recipe
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
  2. Read `justfile`
- **Expected**: `claude:` recipe appears only once (not duplicated)
- **Priority**: P1

## TC-027: forge init — runs config init when no config
- **Source**: Task 3 / AC-5
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/config-init-trigger
- **Pre-conditions**: `.forge/config.yaml` does NOT exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init` (complete all init steps including interactive config)
- **Expected**: `.forge/config.yaml` is created via config init; output shows `CREATED .forge/config.yaml (interactive)`
- **Priority**: P0

## TC-028: forge init — skips existing files
- **Source**: Task 3 / AC-6, Hard Rules
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/skip-existing-files
- **Pre-conditions**: `.forge/`, `CLAUDE.md`, `.forge/config.yaml` all already exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
- **Expected**: Output shows `SKIPPED` for each existing item; no files are overwritten
- **Priority**: P0

## TC-029: forge init — execution result report format
- **Source**: Task 3 / AC-7
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/result-report-format
- **Pre-conditions**: Clean project (no .forge, no CLAUDE.md)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init`
  2. Capture output
- **Expected**: Output contains a report block between `>>>` and `<<<` markers, with CREATED/APPENDED/SKIPPED status for each step
- **Priority**: P1

### Migration (Task 4)

## TC-030: ResolveScope reads config.yaml directly
- **Source**: Task 4 / AC-1, Hard Rules
- **Type**: CLI
- **Target**: cli/resolve-scope
- **Test ID**: cli/resolve-scope/config-read
- **Pre-conditions**: `.forge/config.yaml` exists with `project-type: backend`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Call `ResolveScope()` (via test or indirectly)
  2. Verify no subprocess is spawned
  3. Verify result matches config value
- **Expected**: `ResolveScope()` returns `backend` without calling `just project-type` subprocess
- **Priority**: P0

## TC-031: ResolveScope — missing config returns empty
- **Source**: Task 4 / AC-6, Hard Rules
- **Type**: CLI
- **Target**: cli/resolve-scope
- **Test ID**: cli/resolve-scope/missing-config
- **Pre-conditions**: `.forge/config.yaml` does NOT exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Call `ResolveScope()`
- **Expected**: Returns empty string (skip scope) without error; no subprocess call
- **Priority**: P0

## TC-032: justfile has no project-type recipe
- **Source**: Task 4 / AC-3
- **Type**: CLI
- **Target**: cli/migration
- **Test ID**: cli/migration/no-project-type-recipe
- **Pre-conditions**: Migration task is complete
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `justfile`
  2. Search for `project-type:` recipe
- **Expected**: No `project-type:` recipe exists in justfile
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Task 1 / AC-1, AC-2 | CLI | cli/config-init | P0 |
| TC-002 | Task 1 / AC-3 | CLI | cli/config-init | P0 |
| TC-003 | Task 1 / AC-3 | CLI | cli/config-init | P1 |
| TC-004 | Task 1 / AC-4 | CLI | cli/config-get | P0 |
| TC-005 | Task 1 / AC-5 | CLI | cli/config-get | P0 |
| TC-006 | Task 1 / AC-6 | CLI | cli/config-get | P0 |
| TC-007 | Task 1 / AC-7 | CLI | cli/config-struct | P1 |
| TC-008 | Task 2 / AC-1 | CLI | cli/proposal-list | P0 |
| TC-009 | Task 2 / AC-2 | CLI | cli/proposal-detail | P0 |
| TC-010 | Task 2 / AC-3 | CLI | cli/proposal-list | P1 |
| TC-011 | Task 2 / AC-4 | CLI | cli/proposal-list | P1 |
| TC-012 | Task 2 / AC-5 | CLI | cli/proposal-list | P1 |
| TC-013 | Task 2 / AC-6 | CLI | cli/feature-list | P0 |
| TC-014 | Task 2 / AC-7 | CLI | cli/feature-list | P1 |
| TC-015 | Task 2 / AC-8 | CLI | cli/feature-list | P1 |
| TC-016 | Task 2 / AC-9 | CLI | cli/feature-status | P0 |
| TC-017 | Task 2 / Hard Rules | CLI | cli/feature | P0 |
| TC-018 | Task 2 / AC-10 | CLI | cli/lesson-list | P0 |
| TC-019 | Task 2 / AC-11 | CLI | cli/lesson-list | P1 |
| TC-020 | Task 2 / AC-12 | CLI | cli/lesson-detail | P0 |
| TC-021 | Task 3 / AC-1 | CLI | cli/init | P0 |
| TC-022 | Task 3 / AC-2 | CLI | cli/init | P0 |
| TC-023 | Task 3 / AC-3 | CLI | cli/init | P0 |
| TC-024 | Task 3 / AC-3 | CLI | cli/init | P1 |
| TC-025 | Task 3 / AC-4 | CLI | cli/init | P0 |
| TC-026 | Task 3 / AC-4 | CLI | cli/init | P1 |
| TC-027 | Task 3 / AC-5 | CLI | cli/init | P0 |
| TC-028 | Task 3 / AC-6 | CLI | cli/init | P0 |
| TC-029 | Task 3 / AC-7 | CLI | cli/init | P1 |
| TC-030 | Task 4 / AC-1 | CLI | cli/resolve-scope | P0 |
| TC-031 | Task 4 / AC-6 | CLI | cli/resolve-scope | P0 |
| TC-032 | Task 4 / AC-3 | CLI | cli/migration | P1 |
