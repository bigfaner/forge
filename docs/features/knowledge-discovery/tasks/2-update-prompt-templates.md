---
id: "2"
title: "Update prompt templates with discovery instruction"
priority: "P0"
estimated_time: "30m"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Update prompt templates with discovery instruction

## Description

Ensure all 5 prompt templates contain the discovery instruction that uses `domains` frontmatter for relevance matching. The templates already contain a partial knowledge-check instruction ("Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge") but need the `domains` frontmatter filtering logic added.

## Reference Files
- `docs/proposals/knowledge-discovery/proposal.md` — Source proposal (discovery instruction text)

## Affected Files

### Create

| File | Description |
|------|-------------|
| (none) | |

### Modify

| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/fix.md` | Replace vague instruction with discovery instruction |
| `forge-cli/pkg/prompt/data/cleanup.md` | Replace vague instruction with discovery instruction |
| `forge-cli/pkg/prompt/data/enhancement.md` | Replace vague instruction with discovery instruction |
| `forge-cli/pkg/prompt/data/feature.md` | Replace vague instruction with discovery instruction |
| `forge-cli/pkg/prompt/data/refactor.md` | Replace vague instruction with discovery instruction |

### Delete

| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] All 5 files contain the discovery instruction
- [ ] All 5 files include the `domains` frontmatter filtering logic ("Read each file's YAML frontmatter `domains` field to determine relevance")
- [ ] The discovery instruction matches the proposal's specification:
  ```
  Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge relevant to this task.
  Read each file's YAML frontmatter `domains` field to determine relevance.
  Load files whose domains overlap with the task context.
  If no files match, skip — no matching convention files for this task.
  ```

## Hard Rules

- The discovery instruction must include both the knowledge-check directive and the `domains` frontmatter filtering logic
- Do not change any other content in the prompt templates

## Implementation Notes

- All 5 files have the same pattern at the same location (lines 12-15 or 20-23) — batch replacement
- The instruction goes where the current knowledge-check instruction is, maintaining the same context flow
