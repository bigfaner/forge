# Directory Structure Redesign for AI Programming

**Date**: 2026-04-09
**Status**: Approved
**Affects**: zcode plugin v2.0.0 redesign plan (`docs/zcode-redesign-plan.md`)

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
      manifest.md                    # Feature index & linkage map (auto-generated)
      prd/
        prd-overview.md              # PRD overview (was prd.md)
        prd-user-stories.md          # User stories
        prd-ui-functions.md          # UI function requirements (NEW)
      design/
        design-overview.md           # Technical design (was design.md)
        design-api.md                # API documentation (NEW)
      ui/
        ui-design.md                 # UI component specs (NEW)
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
| `prd-overview.md` not `overview.md` | Source-prefixed names prevent grep ambiguity; AI can glob `prd-*` or `design-*` |
| `ui/` at feature level | Parallel to `design/`, not nested under it; mirrors the `/ui-design` skill's independence |
| `manifest.md` at feature root | Single entry point for AI context; auto-generated and maintained by skills |
| No `tech/` dir | Renamed to `design/` for consistency with skill naming (`/design-tech`) |
| `proposals/` stays separate | Brainstorm is pre-feature; merging it into feature dirs would pollute the feature namespace |

## Manifest File Format

```markdown
# Feature: <slug>

## Status: prd | design | tasks | in-progress | done

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Overview | prd/prd-overview.md | <1-line auto-generated summary> |
| User Stories | prd/prd-user-stories.md | <1-line auto-generated summary> |
| UI Functions | prd/prd-ui-functions.md | <1-line auto-generated summary> |
| Design Overview | design/design-overview.md | <1-line auto-generated summary> |
| API Design | design/design-api.md | <1-line auto-generated summary> |
| UI Design | ui/ui-design.md | <1-line auto-generated summary> |

## Traceability

| PRD Section → | Design Section → | Tasks |
|---------------|-------------------|-------|
| <prd ref> | <design ref> | <task IDs> |
```

### Auto-generation Rules

| Skill | Action on manifest.md |
|-------|----------------------|
| `/write-prd` | Create manifest with PRD entries + summaries |
| `/design-tech` | Add design entries + traceability links to PRD sections |
| `/ui-design` | Add UI entry + traceability links to PRD UI functions |
| `/breakdown-tasks` | Update Tasks column with task IDs linked to design sections |
| Status field | Updated automatically: exists when documents are created |

## Workflow Mapping

```
/brainstorm → /write-prd → /eval-prd → /design-tech → /eval-design → /breakdown-tasks
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
| `/design-tech` | `manifest.md` → `prd/prd-overview.md` | `design/design-*.md` (2 files), `manifest.md` |
| `/ui-design` | `manifest.md` → `prd/prd-ui-functions.md` | `ui/ui-design.md`, `manifest.md` |
| `/eval-design` | `manifest.md` → `design/design-*.md`, `ui/ui-design.md` | eval report |
| `/breakdown-tasks` | `manifest.md` → all docs | `tasks/<N.N>-*.md`, `index.json`, `manifest.md` |

Every skill reads `manifest.md` first, then loads only the specific documents it needs.

## Impact on Redesign Plan (docs/zcode-redesign-plan.md)

### Phase 0: Task-CLI Constants

- `PRDFileName`: `"overview.md"` → `"prd-overview.md"`
- `DesignFileName`: `"overview.md"` → `"design-overview.md"`
- New constants:
  ```go
  PRDOverviewFile   = "prd-overview.md"
  PRDUserStoriesFile = "prd-user-stories.md"
  PRDUIFunctionsFile = "prd-ui-functions.md"
  DesignOverviewFile = "design-overview.md"
  DesignAPIFile      = "design-api.md"
  UIDesignDir        = "ui"
  UIDesignFile       = "ui-design.md"
  ManifestFileName   = "manifest.md"
  ProposalBaseDir    = "docs/proposals"
  ProposalFileName   = "proposal.md"
  ```
- New path functions: `GetFeatureManifest`, `GetFeatureUIDesignDir`, `GetFeatureUIDesignFile`, `GetProposalDir`, `GetProposalFile`
- `EnsureFeatureDir`: create `prd/`, `design/`, `ui/`, `tasks/` dirs

### Phase 1: Templates

- `write-prd/templates/prd-overview.md` (renamed from prd.md, content updated for new format)
- `write-prd/templates/prd-user-stories.md` (unchanged)
- `write-prd/templates/prd-ui-functions.md` (NEW)
- `write-prd/templates/manifest.md` (NEW — manifest template with PRD section)
- `design-tech/templates/design-overview.md` (renamed from design.md)
- `design-tech/templates/design-api.md` (NEW)
- `design-tech/templates/manifest-update-design.md` (NEW — manifest update snippet)
- `brainstorm/templates/proposal.md` (NEW)
- `ui-design/templates/ui-design.md` (NEW)
- `ui-design/templates/manifest-update-ui.md` (NEW — manifest update snippet)
- `breakdown-tasks/templates/index.json` — update `prd` and `design` paths
- `breakdown-tasks/templates/manifest-update-tasks.md` (NEW — manifest update snippet)

### Phase 4: Skill Refactoring

All skill modifications use the new filenames. Key changes:

- **write-prd**: Output to `prd/prd-overview.md`, `prd/prd-user-stories.md`, `prd/prd-ui-functions.md`; create `manifest.md`
- **design-tech**: Read from `prd/prd-overview.md` (via manifest); output to `design/design-overview.md`, `design/design-api.md`; update `manifest.md`
- **ui-design**: Read from `prd/prd-ui-functions.md` (via manifest); output to `ui/ui-design.md`; update `manifest.md`
- **eval-prd**: Locate via `manifest.md`; evaluate `prd/prd-*.md`
- **eval-design**: Locate via `manifest.md`; evaluate `design/design-*.md`, `ui/ui-design.md`
- **breakdown-tasks**: Read `manifest.md` → all docs; output tasks; update `manifest.md` traceability section

### Phase 5: Guide and Hooks

- `guide.md`: Update Document Index to describe new structure and manifest pattern
- `plugin.json`: version `2.0.0`, add keywords `"brainstorm"`, `"ui-design"`, `"manifest"`
