---
id: "1"
title: "Add capability validation in profile package"
priority: "P1"
estimated_time: "30min"
dependencies: []
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 1: Add capability validation in profile package

## Description

Add explicit Go validation for test-type capabilities. Currently these values are defined only implicitly in profile manifest YAMLs — no compile-time safety, no validation at runtime. This task creates the validation primitives that Task 2 will consume.

## Reference Files
- `docs/proposals/task-index-dedup-legacy-types/proposal.md` — Source proposal
- `forge-cli/pkg/profile/embed.go` — Existing `UnionCapabilities()` and `GetProfileCapabilities()`
- `forge-cli/pkg/profile/embed_test.go` — Existing tests for `UnionCapabilities`
- `forge-cli/pkg/profile/profiles/*/manifest.yaml` — Profile manifests defining capability enums

## Acceptance Criteria

- [ ] `ValidTestTypes` constant set defined in `pkg/profile/embed.go` with values: `web-ui`, `tui`, `mobile-ui`, `api`, `cli`
- [ ] `ValidateCapabilities(caps []string) error` function rejects any value not in `ValidTestTypes` with actionable error message listing valid values
- [ ] Unit tests cover: valid single, valid multiple, invalid value, empty input, case sensitivity
- [ ] All existing tests still pass: `go test -race -cover ./forge-cli/pkg/profile/...`

## Hard Rules

- Follow TDD: write tests first (RED), implement (GREEN), refactor if needed
- `ValidateCapabilities` should be case-sensitive (values must match manifest YAML exactly)

## Implementation Notes

- The closed enum `{web-ui, tui, mobile-ui, api, cli}` is sourced from all profile manifests under `pkg/profile/profiles/`. Cross-reference with actual manifests to ensure completeness.
- `ValidTestTypes` should be a `map[string]bool` or `[]string` — whichever matches existing patterns in the codebase.
