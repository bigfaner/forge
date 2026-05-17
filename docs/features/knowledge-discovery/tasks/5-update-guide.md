---
id: "5"
title: "Update guide.md to reference domains frontmatter"
priority: "P2"
estimated_time: "15m"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 5: Update guide.md to reference domains frontmatter

## Description

Update the project knowledge note in `guide.md` to inform agents that convention and business-rule files use `domains` frontmatter for self-description.

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
| `plugins/forge/hooks/guide.md` | Update project knowledge note |

### Delete

| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] The note on line 30 mentions that convention/business-rule files have `domains` frontmatter
- [ ] The note explains that domains are used for relevance matching during task execution
- [ ] The note is concise (1-2 sentences) — guide.md is session-injected, every line costs tokens

## Hard Rules

- Keep the change minimal — only update the existing note, do not add new sections

## Implementation Notes

- Current text: `> Agents read docs/business-rules/ and docs/conventions/ during task execution for domain constraints and coding standards. These are populated by /consolidate-specs, which also performs drift verification to keep specs in sync with code.`
- Add a brief mention of `domains` frontmatter and its purpose
