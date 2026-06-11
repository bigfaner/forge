---
id: "5"
title: "Migration logic and proposal intent field"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 5: Migration logic and proposal intent field

## Description

Two changes: (1) Update migration logic so that `implementation` type in old index.json files maps to `feature` (conservative default — feature is the broadest code-change type). (2) Add `intent` field to proposal frontmatter metadata, so proposals can declare their dominant intent (feature/enhancement/fix/cleanup/refactor) which tasks inherit by default.

## Reference Files
- `docs/proposals/task-type-refinement/proposal.md` — Source proposal (D1, migration risk)
- `forge-cli/internal/cmd/migrate.go` — Migration command (lines 29-71)
- `forge-cli/pkg/proposal/proposal.go` — `Metadata` struct (lines 116-133)

## Acceptance Criteria
- [ ] `migrate.go` maps `TypeImplementation` → `TypeFeature` as the conservative default fallback
- [ ] `forge task migrate` on an old index.json with `type: "implementation"` produces `type: "feature"`
- [ ] `proposal.Metadata` struct has `Intent string \`yaml:"intent"\`` field
- [ ] `parseFrontmatter()` correctly parses `intent` from proposal frontmatter
- [ ] Empty/missing `intent` field does not break existing proposals

## Hard Rules
- Migration must be idempotent: running it twice produces the same result.
- `feature` is the conservative default because it's the broadest — any `implementation` task could plausibly be a feature. More specific types (enhancement/cleanup/refactor) should be set by the user or skill, not migration.

## Implementation Notes
- In `migrate.go` line 56, the current fallback is `task.TypeImplementation`. Change to `task.TypeFeature`.
- The `intent` field in proposal Metadata is consumed by `/quick-tasks` and `/breakdown-tasks` skills (task 6). The CLI only needs to parse and expose it — the skill logic for using it is in the markdown.
