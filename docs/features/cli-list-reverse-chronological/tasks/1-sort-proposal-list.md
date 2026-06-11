---
id: "1"
title: "Sort forge proposal output by created date descending"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Sort forge proposal output by created date descending

## Description
`forge proposal` lists proposals in filesystem (lexical) order. Add descending sort by the `Created` field (already parsed from frontmatter by `proposal.Discover()`).

## Reference Files
- `docs/proposals/cli-list-reverse-chronological/proposal.md` — Source proposal
- `forge-cli/internal/cmd/proposal.go` — `runProposalList()` function
- `forge-cli/pkg/proposal/proposal.go` — `Proposal` struct and `Discover()` function
- `forge-cli/internal/cmd/proposal_test.go` — Existing tests
- `forge-cli/pkg/proposal/proposal_test.go` — Package tests

## Acceptance Criteria
- [ ] `runProposalList()` sorts proposals by `Created` date descending (newest first)
- [ ] Proposals without `created` frontmatter still sort correctly (fallback mtime)
- [ ] Existing tests continue to pass
- [ ] New test verifies sort order

## Hard Rules
- Use `sort.Slice()` or `slices.SortFunc()` — no external sort library

## Implementation Notes
- The `Proposal` struct already has a `Created` field (string `YYYY-MM-DD` format). Parse to `time.Time` for comparison.
- `Discover()` already handles mtime fallback for missing `created`.
