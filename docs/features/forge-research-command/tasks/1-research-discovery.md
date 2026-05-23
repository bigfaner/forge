---
id: "1"
title: "Implement research discovery package"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Implement research discovery package

## Description

Create `pkg/research/research.go` with Discover and FindBySlug functions that parse research report frontmatter from `docs/research/<slug>.md` files. Follow the established pattern from `pkg/proposal/` and `pkg/lesson/`.

Research reports have this frontmatter:
```yaml
created: "YYYY-MM-DD"
topic: "Research Topic"
mode: "deep-dive" | "comparison"
dimensions: [...]
candidates: [...]
```

Include comprehensive unit tests in `pkg/research/research_test.go` using TDD (RED → GREEN → REFACTOR).

## Reference Files
- `docs/proposals/forge-research-command/proposal.md` — Source proposal
- `forge-cli/pkg/proposal/proposal.go` — Pattern reference for Discover/FindBySlug
- `forge-cli/pkg/lesson/lesson.go` — Pattern reference for frontmatter parsing + mtime fallback
- `forge-cli/pkg/proposal/proposal_test.go` — Test pattern reference
- `forge-cli/pkg/lesson/lesson_test.go` — Test pattern reference (table-driven, edge cases)

## Acceptance Criteria

- [ ] `Discover(projectRoot)` walks `docs/research/*.md` and returns `[]Report` with parsed frontmatter (slug, created, topic, mode, dimensions, candidates, filePath)
- [ ] `FindBySlug(projectRoot, slug)` returns a single `*Report` by slug or error if not found
- [ ] Reports sorted by created date descending (newest first), mtime as fallback
- [ ] Graceful handling: missing `docs/research/` directory returns empty slice, no error
- [ ] Graceful handling: empty `docs/research/` directory returns empty slice
- [ ] Graceful handling: files with malformed or missing frontmatter are skipped
- [ ] Created date falls back to mtime when frontmatter `created` is missing
- [ ] Unit tests cover: empty dir, no dir, single report, multiple reports, no frontmatter, malformed frontmatter, FindBySlug found/not found, sorting order

## Hard Rules

- Follow dependency direction: `pkg/` has no dependency on `internal/cmd/`
- Use `gopkg.in/yaml.v3` for frontmatter parsing (consistent with proposal/lesson packages)
- Use `testify/assert` (NOT `require`) in tests

## Implementation Notes

- Research report structure differs from proposal/lesson: flat files in `docs/research/<slug>.md` (like lessons), not nested directories (like proposals)
- Dimensions and candidates are optional fields — handle gracefully when absent
- Slug is derived from filename (strip `.md` extension)
