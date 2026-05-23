---
created: "2026-05-23"
status: Draft
intent: "Add forge research CLI command to list and view deep-research reports"
---

# Proposal: forge research Command

## Problem

The `/deep-research` skill generates structured research reports in `docs/research/<slug>.md` with rich frontmatter metadata (topic, mode, dimensions, candidates). However, there is no CLI command to browse or query these reports. Users must manually navigate the filesystem to discover what research exists or read report metadata.

This is an inconsistency in the CLI surface: `forge proposal` and `forge lesson` both provide structured list/detail views, but research reports — which follow the same frontmatter-convention pattern — lack equivalent CLI access.

## Solution

Add a `forge research [slug]` command following the established pattern used by `forge proposal` and `forge lesson`:

- **No arguments**: list all research reports in a table with slug, created date, topic, and mode
- **With slug argument**: show detailed metadata for a specific report

### Command Behavior

**`forge research`** (list mode):
```
── RESEARCH ──────────────────────────────────
  2 found

  SLUG        CREATED     TOPIC                    MODE
  ----------  ----------  -----------------------  ----------
  gbrain      2026-05-15  GBrain AI Platform       deep-dive
  graphify    2026-05-10  Graphify Knowledge Graph  deep-dive
──────────────────────────────────────────────────
```

**`forge research <slug>`** (detail mode):
```
── RESEARCH ──────────────────────────────────
  SLUG:       gbrain
  TOPIC:      GBrain AI Platform
  CREATED:    2026-05-15
  MODE:       deep-dive
  DIMENSIONS: Overview & Positioning, Architecture, Learning Curve
  FILE:       docs/research/gbrain.md
──────────────────────────────────────────────────
```

## Alternatives

| Alternative | Trade-off |
|-------------|-----------|
| **Do nothing** — users run `ls docs/research/` manually | No structured view, no metadata parsing, inconsistent with proposal/lesson UX |
| **Open in browser** — list + open HTML preview | Over-engineered for markdown files; adds unnecessary complexity |

## Scope

### In Scope
- `pkg/research/research.go` — Discover and FindBySlug functions, parsing report frontmatter
- `pkg/research/research_test.go` — Unit tests with table-driven cases
- `cmd/research.go` — CLI command with list and detail modes
- Registration in `cmd/root.go`

### Out of Scope
- Search/filter functionality (can be added later if needed)
- Report creation or editing (handled by `/deep-research` skill)
- HTML rendering or browser integration
- Subcommand structure (use flat `[slug]` arg pattern, not `list`/`show`)

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| Report frontmatter format diverges from template | Low | Graceful fallback: skip malformed files, use mtime as date fallback |
| No reports exist yet in user project | Low | Print "no research found" message, matching proposal/lesson pattern |

## Success Criteria

- [ ] `forge research` lists all reports in `docs/research/` with parsed frontmatter
- [ ] `forge research <slug>` shows detail view with all metadata fields
- [ ] Graceful handling of missing/empty `docs/research/` directory
- [ ] Graceful handling of reports with malformed frontmatter
- [ ] Unit tests cover Discover, FindBySlug, and edge cases (empty dir, no frontmatter)
- [ ] Version bump in `scripts/version.txt` (minor: new command)
