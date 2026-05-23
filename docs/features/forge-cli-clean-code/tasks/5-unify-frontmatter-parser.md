---
id: "5"
title: "Extend ParseFrontmatter as shared YAML parser"
priority: "P1"
estimated_time: "1.5h"
dependencies: [4]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 5: Extend ParseFrontmatter as shared YAML parser

## Description
YAML frontmatter parsing is independently implemented in 4+ files with varying signatures. Extend the existing `pkg/task.ParseFrontmatter()` with a two-layer API: return raw YAML bytes as a second return value, so callers can unmarshal into their own structs. This eliminates all duplicate parsing implementations. Phase 3 (duplicate logic consolidation).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/pkg/task/frontmatter.go` — Existing shared parser
- Files with duplicate parsing: search for `frontmatter` or `yaml` parsing across the codebase

## Acceptance Criteria
- [ ] `ParseFrontmatter()` returns `(metadata, rawYAMLBytes, error)` or equivalent two-layer API
- [ ] All 4+ duplicate frontmatter parsing sites replaced with calls to the shared parser
- [ ] Zero duplicate frontmatter parsing implementations remain
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- The shared API must be backward-compatible with existing callers of `ParseFrontmatter()`
- Each caller's specific struct unmarshaling remains at the call site

## Implementation Notes
- First grep for all frontmatter/yaml parsing implementations to get exact count and signatures
- The two-layer API: `ParseFrontmatter() ([]byte, []byte, error)` where first return is the parsed generic result and second is raw YAML bytes
- Risk: different callers may use different YAML parsers — ensure the shared function uses the most common one
