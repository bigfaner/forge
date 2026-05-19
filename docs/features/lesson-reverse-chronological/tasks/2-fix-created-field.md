---
id: "2"
title: "Fix Go parser to support created frontmatter field"
priority: "P2"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 2: Fix Go parser to support created frontmatter field

## Description
The Go `Metadata` struct uses `yaml:"date"` to parse frontmatter, but lesson templates use the `created` field name. About 40 files use `created`, 20+ have no frontmatter, and only 1 uses `date`. The parser misses the `created` field, causing the `Created` display value to fall back to file modification time unnecessarily.

## Reference Files
- `docs/proposals/lesson-reverse-chronological/proposal.md` — Source proposal
- `forge-cli/pkg/lesson/lesson.go` — Metadata struct and Discover() function

## Acceptance Criteria
- [ ] Metadata struct parses both `created` and `date` frontmatter fields
- [ ] `created` field takes priority over `date` when both are present
- [ ] Existing tests continue to pass
- [ ] New parsing logic has unit test coverage

## Hard Rules
- Support both `created` and `date` fields for backward compatibility

## Implementation Notes
- Add a `CreatedRaw string yaml:"created"` field (or use inline YAML unmarshaling) to handle both field names
- In `Discover()`, check the `created` field first, then fall back to `date`, then to file modification time
- Alternative: use a custom `UnmarshalYAML` method on Metadata to normalize the field names
