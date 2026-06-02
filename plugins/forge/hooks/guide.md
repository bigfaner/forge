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
- `forge task status <id>` — query current task status and record info
- `forge task add --type <type> --title "..." [--source-task-id <id>] [--block-source] [--var KEY=VALUE] --description "..."` — create a new task (fix-tasks use `--source-task-id --block-source`)
- `forge task transition <id> <status> --reason "..."` — manually transition (unblock, skip, reject)
- `forge task reopen <id>` — re-activate a rejected/skipped task
- `forge task submit <id> --data <record-path>` — submit task execution record

### Feature Management

- `forge feature set <slug>` — set the active feature context for the session
- `forge feature complete --if-done` — mark feature completed if all tasks done (stop hook)

### Configuration

Forge uses `forge config get <key>` to read project-level automation settings. Common keys:

| Key | Purpose | Used by |
|-----|---------|---------|
| `auto.runTasks` | Auto-run task dispatch after quick pipeline | `/quick` |
| `auto.knowledgeSave` | Auto-save extracted knowledge after bug fix | `/fix-bug` |
| `auto.eval.proposal` | Auto-evaluate proposals after brainstorm | `/brainstorm` |
| `auto.eval.prd` | Auto-evaluate PRD after write-prd | `/write-prd` |
| `auto.eval.techDesign` | Auto-evaluate tech design | `/tech-design` |
| `auto.eval.uiDesign` | Auto-evaluate UI design | `/ui-design` |

All keys return `true` or empty. Check with: `forge config get <key>`; if output contains `true`, the behavior is enabled.

### Pipeline Utilities

- `forge prompt get-by-task-id <id>` — retrieve task execution prompt (dispatcher/agent entry point)
- `forge quality-gate` — run quality gate checks for the current feature
- `forge surfaces detect` — auto-detect surface types for configured surface keys
- `forge task index --feature <slug>` — regenerate task index from task files
- `forge task validate-index <path>` — validate index.json structure
- `forge cleanup` — clean stale artifacts

## Testing

| Surface | Test Type | Execution |
|---------|-----------|-----------|
| `cli` | CLI Functional Test | subprocess (exit code + stdout) |
| `tui` | Terminal Functional Test | subprocess + stdin pipe |
| `api` | API Functional Test | HTTP client |
| `web` | Web E2E Test | browser automation |
| `mobile` | Mobile E2E Test | Maestro YAML / manual |

> **"e2e" is reserved for Web and Mobile only.** CLI/TUI/API tests use "Functional Test" — their validation is protocol-level, not device-level automation.

Test file locations: `tests/{surface}/` for cli/tui/api, `tests/e2e/` for web/mobile. Run `/test-guide` for full per-surface strategy.

## Terminology

- **Surface**: a testable system entry point managed by Forge (e.g. a web app, an API server, a CLI binary). Each Surface is identified by a user-defined **Surface Key** (alphanumeric + `-_`) configured in `.forge/config.yaml`.
- **Surface Type**: the kind of surface — one of `web`, `api`, `cli`, `tui`, `mobile`. Determines the orchestration strategy for build/dev/test. `web`/`api` require probe + teardown; `mobile` adds test-setup before the same lifecycle (test-setup → dev → probe → test → teardown); `cli`/`tui` use build → dev → test. Auto-detected via `forge surfaces detect`.
- **Test Type**: the test classification derived from Surface Type. Each surface maps to a specific test type: `cli` → CLI Functional Test, `tui` → Terminal Functional Test, `api` → API Functional Test, `web` → Web E2E Test, `mobile` → Mobile E2E Test. "e2e" is used exclusively for Web/Mobile surfaces. See [test-type-model.md](../../docs/reference/test-type-model.md) for the full mapping and classification rules.
