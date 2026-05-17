---
id: "1"
title: "Context Injection Infrastructure"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Context Injection Infrastructure

## Description

Implement the foundation for context injection: define the `context` frontmatter field spec for rubrics, and update eval SKILL.md pre-processing to read the declaration, filter conventions/business-rules files, and inject the content into the scorer prompt.

This is Batch 1 from the proposal — all downstream batches depend on this infrastructure.

## Reference Files
- `docs/proposals/eval-reality-validation/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Add context parsing in Step 1.4, context injection in Step 2 |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] Eval SKILL.md defines the `context` frontmatter field spec with two sub-fields: `conventions` (list of category strings) and `business-rules` (`auto` or list of filenames)
- [ ] Step 1 (Resolve Type) reads the rubric frontmatter `context` field and stores the declaration
- [ ] Step 1.4 (Pre-Processing by Type) gains a new generic entry: "All types: if rubric has `context` frontmatter, load filtered conventions and business-rules"
- [ ] Step 2 (Invoke Scorer Subagent) injects filtered context content as an additional prompt section when context is declared
- [ ] Context filtering logic: `conventions` field matches files in `docs/conventions/` by filename prefix or category; `business-rules: auto` loads all files in `docs/business-rules/`
- [ ] Missing convention/business-rule files are skipped silently (no error, no abort)
- [ ] doc-scorer.md and doc-reviser.md are NOT modified

## Hard Rules

- Do NOT modify doc-scorer.md or doc-reviser.md — context injection happens at the eval skill orchestrator level via prompt construction
- Do NOT create new files outside of the eval SKILL.md changes — this task is only the infrastructure, not rubric updates

## Implementation Notes

- The `context` field is additive to existing rubric frontmatter (`scale`, `target`, `iterations`, `type`). Rubrics without `context` continue to work unchanged.
- Conventions filtering by category: the `conventions` list contains strings like `[api, naming, ux]`. Match against filenames in `docs/conventions/` — e.g., `api` matches `api*.md`, `ux` matches `ux*.md`. If a string doesn't match any file, skip silently.
- Business-rules `auto`: load all `.md` files from `docs/business-rules/`. Future enhancement can be domain-based filtering.
- The injection format in the scorer prompt should be a clearly demarcated section (e.g., `<injected-context>...content...</injected-context>`) so the scorer knows this is external reference material, not part of the document being evaluated.
