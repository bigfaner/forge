# Directory Structure Redesign for AI Programming

**Date**: 2026-04-09
**Status**: Approved
**Supersedes**: The directory structure section of `docs/zcode-redesign-plan.md`. This spec's structure takes precedence — the parent plan must be updated to match.

---

## Problem

The current redesign plan's directory structure has two issues for AI agent workflows:

1. **Context inefficiency**: Generic filenames (`overview.md`) in multiple subdirs confuse AI agents. No single entry point means agents must discover and read multiple files to understand a feature.
2. **Weak traceability**: No explicit linkage between PRD sections, design sections, and tasks. Agents can't trace requirements forward or backward through the pipeline.

## Design Goals

- AI agents read one file (`manifest.md`) to understand the full feature context
- Source-prefixed filenames prevent naming collisions across directories
- Flat subdirs at feature level (no sub-subdirs beyond `tasks/`)
- Explicit traceability map from PRD → design → tasks

## Directory Structure

```
docs/
  proposals/
    <slug>/
      proposal.md                    # /brainstorm output
  features/
    <slug>/
      manifest.md                    # Feature index & traceability map (auto-generated)
      prd/
        prd-spec.md                  # PRD Spec (was prd.md)
        prd-user-stories.md          # User Stories
        prd-ui-functions.md          # UI Functions (NEW)
      design/
        tech-design.md               # Tech Design (was design.md)
        api-handbook.md              # API Handbook (NEW)
      ui/
        ui-design.md                 # UI Design (NEW)
        *.pen                        # External design tool artifacts
      tasks/
        index.json                   # Task definitions
        <N.N>-<slug>.md              # Task detail files
        process/                     # Runtime state
        records/                     # Execution records
```

### Key Decisions

| Decision | Rationale |
|----------|-----------|
| `prd-spec.md` — PRD Spec | Carries background, goals, key business flows. Not just an "overview" — it's the authoritative source for what the feature is and why it exists. |
| `ui/` at feature level (NOT nested under `design/`) | Parallel to `design/`, mirrors the `/ui-design` skill's independence. Overrides parent plan's `design/ui/`. |
| `manifest.md` at feature root | Single entry point for AI context; auto-generated and maintained by skills |
| No `tech/` dir | Renamed to `design/` for consistency with skill naming (`/design-tech`) |
| `proposals/` (plural) | Convention: directory names are plural. Overrides parent plan's `proposal/` (singular). |
| `proposals/` stays separate from `features/` | Brainstorm is pre-feature; merging it into feature dirs would pollute the feature namespace |
| `prd-ui-functions.md` (requirements) vs `ui/ui-design.md` (spec) | `prd-ui-functions.md` defines WHAT the UI must do (requirements layer). `ui-design.md` defines HOW it looks and behaves (design layer). Clear separation of concerns. |

## Manifest File Format

```markdown
# Feature: <slug>

## Status

<!-- Auto-updated by skills. Do not edit manually. -->
prd -> design -> tasks -> in-progress -> done

Current: <status>

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | <1-line auto-generated summary> |
| User Stories | prd/prd-user-stories.md | <1-line auto-generated summary> |
| UI Functions | prd/prd-ui-functions.md | <1-line auto-generated summary> |
| Tech Design | design/tech-design.md | <1-line auto-generated summary> |
| API Handbook | design/api-handbook.md | <1-line auto-generated summary> |
| UI Design | ui/ui-design.md | <1-line auto-generated summary> |

## Traceability

| PRD Section | Design Section | Tasks |
|-------------|----------------|-------|
| "Functional Requirements > User Auth" (prd-spec §3) | "Architecture > Auth Middleware" (tech-design §2) | 1.2, 1.3 |
| "UI Functions > Login Form" | "UI Design > Login Component" | 2.1 |
```

References use **section headings** from the source documents. AI agents parse the heading text to locate the corresponding section.

**UI-not-applicable features**: For backend/API/CLI features with no UI surface, the `ui/` directory and its manifest entries are omitted entirely. The status advances from `prd` to `design` when `/design-tech` completes, without waiting for `/ui-design`.

### Status State Machine

```
prd ──(/write-prd completes)──→ design ──(/design-tech + /ui-design complete)──→ tasks ──(/breakdown-tasks completes)──→ in-progress ──(first task claimed)──→ done ──(all tasks completed)──→
```

| Current Status | Advances When | Set By |
|---------------|----------------|--------|
| (none) | `/write-prd` completes | `/write-prd` |
| `prd` | `/design-tech` completes (and `/ui-design` if applicable) | `/design-tech` or `/ui-design` (last to complete) |
| `design` | `/breakdown-tasks` completes | `/breakdown-tasks` |
| `tasks` | First task is claimed via `/claim-task` | `/claim-task` |
| `in-progress` | All tasks in `index.json` have status `done` | `/set-task-status` (on last task) |

### Auto-generation Rules

| Skill | Action on manifest.md |
|-------|----------------------|
| `/brainstorm` | Does not touch manifest (writes to `proposals/` only) |
| `/write-prd` | Create manifest with PRD entries + summaries; set status to `prd` |
| `/eval-prd` | Reads manifest only; no manifest update. Eval report is returned in-memory to user, not written to disk. |
| `/design-tech` | Add design entries + traceability links to PRD sections; advance status to `design` if `/ui-design` already completed (or if UI is not applicable) |
| `/ui-design` | Add UI entry + traceability links to PRD UI functions; advance status to `design` if `/design-tech` already completed |
| `/eval-design` | Reads manifest only; no manifest update. Eval report is returned in-memory to user. |
| `/breakdown-tasks` | Update Tasks column with task IDs linked to design sections; advance status to `tasks` |
| `/claim-task` | Advance status to `in-progress` on first claim |
| `/set-task-status` | Advance status to `done` when last task completes |

### Manifest Update Semantics

- **Idempotent**: Re-running a skill merges new content into existing manifest sections. Existing entries are updated (summary text), not duplicated.
- **Manual edits**: The `Documents` and `Traceability` sections are auto-generated. Adding comments (`<!-- -->`) is safe. Manual restructuring may be overwritten on next skill run.
- **Missing files**: If a document referenced in manifest is deleted, the next skill run removes the entry.

## Workflow Mapping

```
/brainstorm → /write-prd → /eval-prd → /design-tech ─→ /eval-design → /breakdown-tasks
     ↓            ↓            ↓            ↓              ↓                ↓
proposal.md  prd/*.{3}   eval report  design/*.{2}  eval report      tasks/*
              manifest.md              manifest.md                   manifest.md
                                                     ↗
                                           /ui-design
                                               ↓
                                         ui/ui-design.md
                                         manifest.md
```

### Per-Skill Context Loading

| Skill | Reads | Writes |
|-------|-------|--------|
| `/brainstorm` | project codebase | `proposals/<slug>/proposal.md` |
| `/write-prd` | proposal.md (optional) | `prd/prd-*.md` (3 files), `manifest.md` |
| `/eval-prd` | `manifest.md` → `prd/prd-*.md` | eval report (in-memory) |
| `/design-tech` | `manifest.md` → `prd/prd-spec.md` | `design/tech-design.md`, `design/api-handbook.md`, `manifest.md` |
| `/ui-design` | `manifest.md` → `prd/prd-ui-functions.md` | `ui/ui-design.md`, `manifest.md` |
| `/eval-design` | `manifest.md` → `design/tech-design.md`, `design/api-handbook.md`, `ui/ui-design.md` | eval report (in-memory) |
| `/breakdown-tasks` | `manifest.md` → all docs | `tasks/<N.N>-*.md`, `index.json`, `manifest.md` |

Every skill reads `manifest.md` first, then loads only the specific documents it needs.

## Impact on Redesign Plan (docs/zcode-redesign-plan.md)

**This spec supersedes the directory structure defined in `zcode-redesign-plan.md`.** The parent plan must be updated to reflect these changes:

### Changes to Parent Plan

| Parent Plan Location | Current Value | Updated Value |
|---------------------|---------------|---------------|
| Target structure `tech/` (line ~16) | `tech/` | `design/` (tree is inconsistent with Phase 0 constants) |
| `EnsureFeatureDir` (line ~63) | `prd/`, `design/`, `design/ui/` | `prd/`, `design/`, `ui/`, `tasks/` |
| `ProposalBaseDir` (line ~53) | `"docs/proposal"` (singular) | `"docs/proposals"` (plural) |
| Phase 2 brainstorm output (line ~119) | `docs/proposal/<slug>/proposal.md` | `docs/proposals/<slug>/proposal.md` |
| Phase 3 ui-design output (line ~133) | `design/ui/` | `ui/` |
| Phase 4.2 design-tech prose (line ~166) | "`design/ui/` 由 `/ui-design` skill 填充" | "`ui/` 由 `/ui-design` skill 填充" |
| Phase 6 e2e verification (line ~223) | "design/ui/ content" | "ui/ content" |

### Phase 0: Task-CLI Constants

- **Rename** `PRDFileName` → `PRDSpecFile`: value `"prd-spec.md"` (constant renamed + value changed)
- **Rename** `DesignFileName` → `TechDesignFile`: value `"tech-design.md"` (constant renamed + value changed)
- Downstream Go code referencing the old constant names will need updating.
- New constants:
  ```go
  PRDSpecFile        = "prd-spec.md"
  PRDUserStoriesFile = "prd-user-stories.md"
  PRDUIFunctionsFile = "prd-ui-functions.md"
  TechDesignFile     = "tech-design.md"
  APIHandbookFile    = "api-handbook.md"
  UIDesignDir        = "ui"
  UIDesignFile       = "ui-design.md"
  ManifestFileName   = "manifest.md"
  ProposalBaseDir    = "docs/proposals"
  ProposalFileName   = "proposal.md"
  ```
- New path functions: `GetFeatureManifest`, `GetFeatureUIDesignDir`, `GetFeatureUIDesignFile`, `GetProposalDir`, `GetProposalFile`
- `EnsureFeatureDir`: create `prd/`, `design/`, `ui/`, `tasks/` dirs

### Phase 1: Templates

- `write-prd/templates/prd-spec.md` (renamed from prd.md, content updated for new format)
- `write-prd/templates/prd-user-stories.md` (unchanged)
- `write-prd/templates/prd-ui-functions.md` (NEW)
- `write-prd/templates/manifest.md` (NEW — manifest template with PRD section)
- `design-tech/templates/tech-design.md` (renamed from design.md)
- `design-tech/templates/api-handbook.md` (NEW)
- `design-tech/templates/manifest-update-design.md` (NEW — manifest update snippet)
- `brainstorm/templates/proposal.md` (NEW)
- `ui-design/templates/ui-design.md` (NEW)
- `ui-design/templates/manifest-update-ui.md` (NEW — manifest update snippet)
- `breakdown-tasks/templates/index.json` — update `prd` and `design` paths
- `breakdown-tasks/templates/manifest-update-tasks.md` (NEW — manifest update snippet)

### Phase 2: Brainstorm Skill

Create `plugins/zcode/skills/brainstorm/SKILL.md` with output path `docs/proposals/<slug>/proposal.md` (plural).

### Phase 3: UI-Design Skill

Create `plugins/zcode/skills/ui-design/SKILL.md` with output path `ui/ui-design.md` (top-level, not nested under `design/`).

### Phase 4: Skill Refactoring

All skill modifications use the new filenames. Key changes:

- **write-prd**: Output to `prd/prd-spec.md`, `prd/prd-user-stories.md`, `prd/prd-ui-functions.md`; create `manifest.md`
- **design-tech**: Read from `prd/prd-spec.md` (via manifest); output to `design/tech-design.md`, `design/api-handbook.md`; update `manifest.md`
- **ui-design**: Read from `prd/prd-ui-functions.md` (via manifest); output to `ui/ui-design.md`; update `manifest.md`
- **eval-prd**: Locate via `manifest.md`; evaluate `prd/prd-*.md`; no manifest update
- **eval-design**: Locate via `manifest.md`; evaluate `design/tech-design.md`, `design/api-handbook.md`, `ui/ui-design.md`; no manifest update
- **breakdown-tasks**: Read `manifest.md` → all docs; output tasks; update `manifest.md` traceability section

### Phase 5: Guide and Hooks

- `guide.md`: Update Document Index to describe new structure and manifest pattern
- `plugin.json`: version `2.0.0`, add keywords `"brainstorm"`, `"ui-design"`, `"manifest"`
