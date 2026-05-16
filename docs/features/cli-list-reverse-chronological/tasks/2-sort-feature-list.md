---
id: "2"
title: "Sort forge feature list output by manifest mtime descending"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Sort forge feature list output by manifest mtime descending

## Description
`forge feature list` lists features in filesystem (lexical) order. Add descending sort by `manifest.md` modification time.

## Reference Files
- `docs/proposals/cli-list-reverse-chronological/proposal.md` — Source proposal
- `forge-cli/internal/cmd/feature.go` — `runFeatureList()` and `discoverFeatures()` functions
- `forge-cli/internal/cmd/feature_test.go` — Existing tests

## Acceptance Criteria
- [ ] `runFeatureList()` sorts features by `manifest.md` mtime descending (newest first)
- [ ] Features with missing/unreadable manifest sort to the end
- [ ] Existing tests continue to pass
- [ ] New test verifies sort order

## Hard Rules
- Use `sort.Slice()` or `slices.SortFunc()` — no external sort library

## Implementation Notes
- `discoverFeatures()` already reads `manifest.md` content. Add `os.Stat()` to capture mtime into a new field on the feature discovery struct.
- Manifest mtime reflects "most recently active" which is the desired semantics for features.
