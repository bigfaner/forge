---
id: "4"
title: "Inline config.yaml and sitemap.json examples into gen-sitemap command"
priority: "P1"
estimated_time: "15m"
dependencies: []
type: "documentation"
mainSession: false
---

# 4: Inline config.yaml and sitemap.json examples into gen-sitemap command

## Description
Replace `${CLAUDE_SKILL_DIR}/../references/shared/config.yaml` and `${CLAUDE_SKILL_DIR}/../references/shared/sitemap.json` references in gen-sitemap.md with inline examples of these files.

## Reference Files
- `docs/proposals/remove-references-dir/proposal.md` — Source proposal
- `plugins/forge/references/shared/config.yaml` — Config template to inline
- `plugins/forge/references/shared/sitemap.json` — Sitemap example to inline

> **Note:** Line numbers are approximate and may drift. Search for `references/shared/config.yaml` or `references/shared/sitemap.json` to locate exact reference sites.

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/gen-sitemap.md` | Replace line 42 (config.yaml template reference) and line 72 (sitemap.json example reference) with inline content |

## Acceptance Criteria
- [ ] No occurrence of `references/shared/config.yaml` in gen-sitemap.md
- [ ] No occurrence of `references/shared/sitemap.json` in gen-sitemap.md
- [ ] Config template and sitemap example are fully inline as code blocks

## Hard Rules
- Use fenced code blocks with appropriate language tags (yaml, json) for the inlined examples

## Implementation Notes
- Line 42 references config.yaml as a "copy from template" source — inline the template content
- Line 72 references sitemap.json as a "full example" — inline the example content
