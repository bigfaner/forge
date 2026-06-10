# Forge Guide

## Directory Conventions

### Project-Level Documents

Non-skill documents shared across features:

```
docs/
  ARCHITECTURE.md       — System architecture
  business-rules/       — Cross-feature business rules (by domain, e.g. auth.md)
  conventions/          — Technical specs (coding standards, API conventions, naming rules)
  reference/            — System specs (environment, deployment, tech stack) [optional]
  decisions/            — Technical decisions (/learn)
  lessons/              — Lessons learned (/learn)
  proposals/            — Improvement proposals (docs/proposals/{slug}/proposal.md, via /brainstorm or ad-hoc)
  sitemap/sitemap.json  — Page element map (project-level, /gen-web-sitemap)
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

### Query commands

Use these for ad-hoc lookups when a user mentions a slug:
- `forge proposal <slug>` — proposal summary (slug, status, created, associated PRD/feature)
- `forge feature status <slug>` — feature status (phase, task breakdown, artifact scores)

### Task Management

- `forge task claim` — claim the next pending task (sets status to in_progress)
- `forge task list [slug]` — list tasks in table format (current feature if slug omitted; `--sort id|topo`; `--tree` for interactive dependency tree TUI)
- `forge task status <id>` — query current task status and record info
- `forge task query <id-or-key>` — query task by ID or key (e.g. "1.2.3" or "phase1-1.1.1-project-init"); `--verbose` to show all fields including related fixes
- `forge task add --type <type> --title "..." [--source-task-id <id>] [--block-source] [--var KEY=VALUE] --description "..."` — create a new task (fix-tasks use `--source-task-id --block-source`)
- `forge task transition <id> <status> --reason "..."` — manually transition (unblock, skip, reject)
- `forge task reopen <id>` — re-activate a rejected/skipped task
- `forge task submit <id> [--data <record-path>] [--quiet]` — submit task execution record (`--quiet` for minimal output)

### Feature Management

- `forge feature set <slug>` — set the active feature context for the session
- `forge feature complete --if-done` — mark feature completed if all tasks done (stop hook)
- `forge feature list` — list all features with status, progress, and scores

### Config Management

- `forge config get <key>` — get a config value from `.forge/config.yaml` (plain text output; exit code 1 if key missing)
- `forge config set <key> <value>` — set a config value in `.forge/config.yaml` (supports dot-notation for nested keys)

Config keys use dot-notation for nested values (e.g. `eval.pr.target`, `auto.runTasks`). Skills use `forge config get` as the standard mechanism to read project-level settings — always redirect stderr (`2>/dev/null`) and check exit code / empty output to handle missing keys gracefully.

### Pipeline Utilities

- `forge prompt get-by-task-id <id>` — retrieve task execution prompt (dispatcher/agent entry point)
- `forge quality-gate` — run quality gate when all tasks completed: compile/fmt/lint → unit tests (with retry-once) → regression; auto-creates fix task on failure; skips compile/test for docs-only features
- `forge surfaces detect` — auto-detect surface types for configured surface keys
- `forge task index --feature <slug>` — regenerate task index from task files
- `forge task validate [file]` — validate index.json structure and task sizing (omit file to validate current feature)
- `forge task check-deps` — validate all task dependencies (references exist, wildcards match at least one task)
- `forge cleanup` — remove state.json and record.json for completed, blocked, suspended, or rejected tasks

## Surfaces

A **Surface** is a testable system entry point (e.g. a web app, API server, CLI binary), identified by a user-defined **Surface Key** (alphanumeric + `-_`) configured in `.forge/config.yaml`. Each surface has a **Surface Type** — one of `web`, `api`, `cli`, `tui`, `mobile` — which determines build/dev/test orchestration and maps to a specific **Test Type**:

| Surface Type | Test Type | Orchestration | Execution |
|--------------|-----------|---------------|-----------|
| `cli` | CLI Functional Test | build → dev → test | subprocess (exit code + stdout) |
| `tui` | Terminal Functional Test | build → dev → test | subprocess + stdin pipe |
| `api` | API Functional Test | dev → probe → test → teardown | HTTP client |
| `web` | Web E2E Test | dev → probe → test → teardown | browser automation |
| `mobile` | Mobile E2E Test | test-setup → dev → probe → test → teardown | Maestro YAML / manual |

> **"e2e" is reserved for Web and Mobile only.** CLI/TUI/API tests use "Functional Test" — their validation is protocol-level, not device-level automation. Test files go to `tests/<surfaceKey>/<journey>/` (multi-surface) or `tests/<journey>/` (single surface), where `surfaceKey` is the key from `forge surfaces` output (e.g. `backend`, `frontend`), not the surface type. Run `/test-guide` for full per-surface strategy.

### Surface Output Parsing

`forge surfaces` (text mode) outputs one surface per line. Skills must use text mode (not `--json`) and apply the unified parsing rule:

```
Per line of forge surfaces output:
  if line contains '=':
    key = part before '='
    type = part after '='
    → named surface (key is set)
  else:
    key = (empty)
    type = line
    → scalar surface (no key)
```

| Config form | Text output | key | type |
|-------------|------------|-----|------|
| Scalar: `surfaces: tui` | `tui` | empty | `tui` |
| Named: `surfaces: [{key: app, type: tui}]` | `app=tui` | `app` | `tui` |
| Multi: two named surfaces | `backend=api` then `frontend=web` | per-line | per-line |
