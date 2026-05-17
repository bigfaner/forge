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

Replace the vague "Read relevant project knowledge files" instruction in all 5 prompt templates with the discovery instruction that uses `domains` frontmatter for relevance matching.

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
- [ ] No file contains the text "Read relevant project knowledge files"
- [ ] The discovery instruction matches the proposal's specification:
  ```
  Check `docs/conventions/` and `docs/business-rules/` for project-specific knowledge relevant to this task.
  Read each file's YAML frontmatter `domains` field to determine relevance.
  Load files whose domains overlap with the task context.
  If no files match, skip — no matching convention files for this task.
  ```

## Hard Rules

- The discovery instruction must replace the existing "Read relevant project knowledge files" line exactly — do not add it as a separate paragraph
- Do not change any other content in the prompt templates

## Implementation Notes

- All 5 files have the same pattern at the same location (line 12 or 20) — batch replacement
- The instruction goes where the current vague instruction is, maintaining the same context flow
