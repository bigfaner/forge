---
id: "7"
title: "Add auto-extract trigger to fix-bug"
priority: "P1"
estimated_time: "30m"
dependencies: ["3"]
type: "enhancement"
mainSession: false
---

# 7: Add auto-extract trigger to fix-bug

## Description

Add knowledge auto-extraction to the `fix-bug` command. After a bug fix is completed and committed, scan the root cause analysis and fix approach for notable knowledge. If found, present extracted knowledge for user confirmation.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal (Part 2)
- `plugins/forge/commands/fix-bug.md` — Current fix-bug command
- `plugins/forge/references/shared/knowledge-extraction.md` — Shared extraction routine (created in Task 3)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | Add knowledge review step after Step 6 (Commit) |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] After Step 6 (Commit), a new knowledge review step runs before Output Summary
- [ ] Step reads `plugins/forge/references/shared/knowledge-extraction.md` for extraction logic
- [ ] Scans root cause analysis and fix approach for notable knowledge
- [ ] Looks for: non-obvious root causes, debugging patterns, gotchas
- [ ] Silent when the fix was trivial (typo, simple config change)
- [ ] Presents extracted knowledge via AskUserQuestion for user confirmation
- [ ] Writes confirmed knowledge to appropriate directories using shared formats
- [ ] Does not interfere with existing fix-bug workflow (Steps 1-6)

## Hard Rules
- Do not modify Steps 1-5 of the existing fix-bug workflow
- Only add after Step 6 (Commit) and before Output Summary
- Must include the shared extraction routine by reference, not by copying its content

## Implementation Notes
- The root cause note from Step 4 ("Root cause: <why>") is a key input for knowledge extraction
- Trivial fixes (typos, simple config) should be filtered out by the "notable knowledge" heuristics
