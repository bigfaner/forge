---
name: CLI Query Knowledge in Guide
slug: cli-query-knowledge-in-guide
status: Approved
created: 2026-05-24
author: faner
---

# CLI Query Knowledge in Guide

## Problem

When a user sends a message like `/quick-tasks <slug>` or mentions a proposal/feature slug in conversation, the agent doesn't know that `forge proposal <slug>` and `forge feature status <slug>` CLI commands exist. Instead of efficiently querying via CLI, the agent fumbles through file system with Glob/Read to piece together proposal and feature status — slow and error-prone.

The CLI commands already provide structured, reliable output for exactly this purpose, but the agent has zero awareness of them.

## Solution

Add a concise **CLI Query Commands** section to `plugins/forge/hooks/guide.md` — the hook file automatically loaded into every forge session. This gives the agent immediate knowledge of these commands without changing any skill's data access patterns.

Two commands only:

- **`forge proposal <slug>`** — query proposal info (slug, status, created date, associated PRD/feature)
- **`forge feature status <slug>`** — query feature status (status, task breakdown by state, artifact scores)

## Alternatives

| # | Approach | Trade-off |
|---|----------|-----------|
| 1 | **Add to guide.md (recommended)** | Zero disruption to existing skills; agent gains CLI knowledge via session hook |
| 2 | Do nothing | Agent continues slow file spelunking; user frustration persists |
| 3 | Rewrite skills to use CLI instead of file reads | Massive refactor; contradicts `forge-cli-reference.md` convention; over-engineered for this need |

## Scope

### In Scope

- Add "CLI Query Commands" section to `plugins/forge/hooks/guide.md` documenting `forge proposal <slug>` and `forge feature status <slug>`
- Include brief description of output format so agent knows what to expect

### Out of Scope

- Changing how skills read proposal/feature data (they continue to use direct file reads)
- Adding other CLI commands beyond these two
- Modifying `forge-cli-reference.md` or any convention documents

## Risks

| Risk | Mitigation |
|------|------------|
| Agent uses CLI when direct file read would be more efficient (e.g., skill already has file path) | Guide text should clarify: CLI commands are for ad-hoc queries; skills should continue using direct file reads |
| Output format changes in future CLI versions | Impact is low — agent adapts to any text output; no structured parsing dependency |

## Success Criteria

- [ ] Agent uses `forge proposal <slug>` when user mentions a proposal slug in conversation
- [ ] Agent uses `forge feature status <slug>` when user mentions a feature slug in conversation
- [ ] Existing skill data access patterns remain unchanged
