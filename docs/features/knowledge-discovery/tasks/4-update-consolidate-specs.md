---
id: "4"
title: "Update consolidate-specs to manage domains frontmatter"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 4: Update consolidate-specs to manage domains frontmatter

## Description

Update the consolidate-specs SKILL.md so that it generates and maintains `domains` frontmatter when creating new convention/business-rule files and updates it during drift detection.

## Reference Files
- `docs/proposals/knowledge-discovery/proposal.md` — Source proposal
- `plugins/forge/skills/consolidate-specs/SKILL.md` — Current skill definition

## Affected Files

### Create

| File | Description |
|------|-------------|
| (none) | |

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Add domains frontmatter generation and drift-update logic |

### Delete

| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] SKILL.md instructs the agent to write `domains` frontmatter when creating new files in `docs/conventions/` or `docs/business-rules/`
- [ ] Domains are derived from the spec content (ID keywords, source keywords) — not freeform
- [ ] During drift detection (Steps 9-11), `domains` are re-derived when file content changes substantially
- [ ] Domain overlap >50% between files triggers a warning during the user confirmation step
- [ ] Each file gets 3-7 specific keywords
- [ ] The existing `title` frontmatter behavior is unchanged

## Hard Rules

- Read `docs/conventions/forge-distribution.md` before modifying files under `plugins/forge/`
- Domains must be derived programmatically (from spec ID and source keywords), not invented by the agent

## Implementation Notes

- The skill already has a domain-to-decision-file mapping table (lines 183-188) — this task adds `domains` as a NEW concern orthogonal to that table
- Drift detection already validates file content against codebase — adding `domains` re-derivation is an extension of that step
- The domain overlap check should use keyword intersection: if two files share >50% of their domain keywords, flag during confirmation
