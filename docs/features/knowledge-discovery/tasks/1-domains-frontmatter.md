---
id: "1"
title: "Add domains frontmatter to convention and business-rule files"
priority: "P0"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Add domains frontmatter to convention and business-rule files

## Description

Add a `domains` field to the YAML frontmatter of all 6 existing convention and business-rule files. Each file self-describes the domains it covers via 3-7 specific keywords derived from the file's own content.

This is the foundational task — all other tasks depend on the frontmatter schema being defined and populated.

## Reference Files
- `docs/proposals/knowledge-discovery/proposal.md` — Source proposal

## Affected Files

### Create

| File | Description |
|------|-------------|
| (none) | |

### Modify

| File | Changes |
|------|---------|
| `docs/conventions/error-handling.md` | Add `domains` frontmatter |
| `docs/conventions/profile-system.md` | Add `domains` frontmatter |
| `docs/conventions/testing-isolation.md` | Add `domains` frontmatter |
| `docs/business-rules/error-reporting.md` | Add `domains` frontmatter |
| `docs/business-rules/quality-gate.md` | Add `domains` frontmatter |
| `docs/business-rules/task-lifecycle.md` | Add `domains` frontmatter |

### Delete

| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] All 6 files have a `domains` field in YAML frontmatter
- [ ] Each file has 3-7 domain keywords
- [ ] Each keyword appears in the file's own content at least once (source code identifier, file path, or spec term)
- [ ] Frontmatter retains existing `title` field unchanged
- [ ] Schema: `domains: [keyword1, keyword2, ...]` — YAML list

## Hard Rules

- Read each file's full content before assigning domains — keywords must reflect actual content, not just the title
- Do NOT modify file content beyond the frontmatter

## Implementation Notes

- Use the file's title, section headings, and key terms as domain keyword candidates
- Prefer specific terms over generic ones (e.g., "exit-code" not "error" if the file specifically covers exit codes)
