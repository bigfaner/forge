# Forge Guide

## Directory Conventions

### Project-Level Documents

Non-skill documents shared across features:

```
docs/
  ARCHITECTURE.md       — System architecture
  business-rules/       — Cross-feature business rules (by domain, e.g. auth.md)
  conventions/          — Technical specs (coding standards, API conventions, naming rules)
  reference/            — System specs (environment, deployment, tech stack)
  decisions/            — Technical decisions (/learn)
  lessons/              — Lessons learned (/learn)
  proposals/            — Improvement proposals (docs/proposals/{slug}/proposal.md, via /brainstorm or ad-hoc)
  sitemap/sitemap.json  — Page element map (project-level, /gen-sitemap)
```

> Agents read `docs/business-rules/` and `docs/conventions/` during task execution for domain constraints and coding standards. Each file carries a `domains` frontmatter field (auto-managed by `/consolidate-specs`) with topic keywords — agents use it to load only files relevant to the current task, skipping the rest. `/consolidate-specs` also performs drift verification to keep specs in sync with code.

### Manifest

`manifest.md` is the single entry point for a Feature. An AI agent reads this file to understand the full context:
- **Documents** table: lists all document paths and auto-generated summaries
- **Tasks** table: task ID, title, status, and file path for each task
- **Status** (feature-level): prd → design → tasks → in-progress → completed
  - Not to be confused with task-level statuses in index.json: pending, in_progress, completed, blocked, suspended, skipped, rejected

## Forge CLI

Run `forge -h` or `forge [command] -h` for full reference.

**Query commands** — use these for ad-hoc lookups when a user mentions a slug:
- `forge proposal <slug>` — proposal summary (slug, status, created, associated PRD/feature)
- `forge feature status <slug>` — feature status (phase, task breakdown, artifact scores)

**Task lifecycle** — error recovery:
- `forge task transition <id> <status> --reason "..."` — manually transition (unblock, skip, reject)
- `forge task reopen <id>` — re-activate a rejected/skipped task

## Terminology

- **Surface**: a testable system entry point managed by Forge (e.g. a web app, an API server, a CLI binary). Each Surface is identified by a user-defined **Surface Key** (alphanumeric + `-_`) configured in `.forge/config.yaml`.
- **Surface Type**: the kind of surface — one of `web`, `api`, `cli`, `tui`, `mobile`. Determines the orchestration strategy for build/dev/test (e.g. `web`/`api` require probe + teardown; `cli`/`tui` use build → dev → test). Auto-detected via `forge surfaces detect`.
