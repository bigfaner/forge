---
id: "1"
title: "Remove noTest from index.json and task frontmatter"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Remove noTest from index.json and task frontmatter

## Description
Batch-remove the deprecated `"noTest": true` field from all index.json files and `noTest: true` lines from task .md frontmatter. The Go runtime already ignores these via `IsTestableType()`, so this is dead data cleanup.

## Reference Files
- `docs/proposals/remove-notest-references/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `docs/features/*/tasks/index.json` (~40 files) | Remove `"noTest": true` entries |
| `docs/features/*/tasks/*.md` (~9 files) | Remove `noTest: true` frontmatter lines |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All `"noTest": true` entries removed from index.json files; JSON remains valid
- [ ] All `noTest: true` frontmatter lines removed from task .md files
- [ ] No other fields in index.json or frontmatter are modified

## Hard Rules
- Do NOT modify files under `docs/proposals/`, `docs/forensics/`, or `docs/self-evolution/`

## Implementation Notes
- Use a script or find-and-replace approach for the ~40 index.json files
- For JSON: remove the `"noTest": true` key-value pair, ensuring trailing commas are handled correctly
- For frontmatter: remove the entire `noTest: true` line
- Verify JSON validity after bulk edits (e.g., `jq . <file>` or `python -m json.tool`)
