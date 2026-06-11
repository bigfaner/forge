---
created: "2026-05-16"
tags: [testing, architecture]
---

# CLI Project E2E Lifecycle: Lightweight by Design

## Problem

When running the test pipeline (T-quick-4 graduate, T-quick-5 verify-regression) on a CLI-only project (forge), it was unclear whether `just e2e-setup` and `just e2e-discover` are necessary overhead or meaningful steps.

## Root Cause

Confusion arises because the e2e lifecycle vocabulary (setup, discover, probe) originates from web-app testing where `e2e-setup` starts a browser/server and `probe` crawls routes. For CLI projects, these steps appear unnecessary at first glance.

## Solution

The go-test profile already adapts these recipes to CLI-appropriate operations:
- `e2e-setup` = `go mod download` — ensures test dependencies are available
- `e2e-discover` = `go test -list '.*' ./...` — lists available e2e test names
- No `probe` recipe exists — web crawling is irrelevant for CLI

Both are lightweight Go idioms, not browser/server overhead.

## Reusable Pattern

When adding a new test profile for non-web projects (CLI, TUI, API-only), ensure `e2e-setup` and `e2e-discover` recipes perform project-appropriate lightweight operations rather than no-ops. The key distinction:

| Lifecycle step | Web profile (playwright) | CLI profile (go-test) |
|----------------|--------------------------|----------------------|
| `e2e-setup` | Start browser/server | `go mod download` |
| `e2e-discover` | Crawl routes, capture DOM | `go test -list` |
| `probe` | Accessibility tree scan | Not applicable — skip |

## Example

```bash
# justfile for go-test profile — CLI-appropriate e2e lifecycle
e2e-setup:
    cd tests/e2e && go mod download

e2e-discover:
    cd tests/e2e && go test -tags=e2e -list '.*' ./...
```

## Related Files

- `justfile` — e2e-setup and e2e-discover recipes
- `.forge/config.yaml` — active profile configuration
- `forge-cli/pkg/profile/` — profile definitions

## References

- Test profile system proposal: `docs/proposals/test-profile-system/proposal.md`
