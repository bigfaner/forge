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
  - Not to be confused with task-level statuses in index.json: pending, in_progress, completed, blocked, skipped, rejected

## Execution Rules

### Quality Gate Protocol

All task-executing workflows MUST pass the quality gate before recording completion. Tasks with a `doc*` type prefix skip the quality gate; only `coding.*` type prefix tasks are gated.

### All-Completed Hook

After all tasks done, `forge quality-gate` runs as a final safety net (project-wide). It automatically skips docs-only features. On failure, a P0 fix-task is automatically created — run `forge task claim` to pick it up.

### Task-CLI

Task CLI manages task lifecycle. Run `forge -h` or `forge [command] -h` for full reference.

Key commands for error recovery:
- `forge task transition <id> <status> --reason "..."` — manually transition a task (unblock, skip, reject)
- `forge task reopen <id>` — re-activate a rejected/skipped task back to pending
