---
id: "2"
title: "Remove run-tasks knowledge review section"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Remove run-tasks knowledge review section

## Description

Remove the entire Knowledge Review section (and its sub-sections: Parameters, Artifact Scanning Scope, Knowledge Types, Extraction Flow, Notable Knowledge Heuristics, Deduplication, Rules) from `plugins/forge/commands/run-tasks.md`. run-tasks is a task dispatcher — real knowledge extraction is handled by `doc.consolidate` tasks.

Also remove the "Commit Remaining Artifacts" section that depends on knowledge extraction output, since it exists solely to commit knowledge files extracted by the removed section.

## Reference Files
- `docs/proposals/auto-knowledge-save/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Remove Knowledge Review section (lines ~101-246) and Commit Remaining Artifacts section (lines ~248-273) |

## Acceptance Criteria
- [ ] `run-tasks.md` no longer contains the Knowledge Review section
- [ ] `run-tasks.md` no longer contains the Commit Remaining Artifacts section
- [ ] Post-Completion section ends at the e2e suggestion paragraph
- [ ] No dangling references to removed sections elsewhere in the file

## Implementation Notes
- The Knowledge Review section starts at the "### Knowledge Review" heading (after Post-Completion).
- The Commit Remaining Artifacts section follows immediately after the Knowledge Review Rules.
- Both sections should be removed cleanly, leaving the Post-Completion section as the last content before the end of the file.
