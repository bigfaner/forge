---
id: "6"
title: "Add auto-extract trigger to run-tasks"
priority: "P1"
estimated_time: "30m"
dependencies: ["3"]
type: "enhancement"
mainSession: false
---

# 6: Add auto-extract trigger to run-tasks

## Description

Add knowledge auto-extraction to the `run-tasks` command. After all tasks complete, scan task outcomes, code changes, and the manifest for notable knowledge. If found, present extracted knowledge for user confirmation.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal (Part 2)
- `plugins/forge/commands/run-tasks.md` — Current run-tasks command
- `plugins/forge/references/shared/knowledge-extraction.md` — Shared extraction routine (created in Task 3)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Add knowledge review step after loop completion |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] After the main task loop ends (Post-Completion section), a new knowledge review step runs
- [ ] Step reads `plugins/forge/references/shared/knowledge-extraction.md` for extraction logic
- [ ] Scans task outcomes, code changes, and manifest for notable knowledge
- [ ] Looks for: architectural decisions, novel patterns, gotchas, business rules
- [ ] Silent when no notable knowledge detected (routine tasks)
- [ ] Presents extracted knowledge via AskUserQuestion for user confirmation
- [ ] Writes confirmed knowledge to appropriate directories using shared formats
- [ ] Does not interfere with existing Post-Completion flow (test/e2e suggestions)

## Hard Rules
- Do not modify the existing task dispatch loop (Steps 0-4)
- Only add to the Post-Completion section
- Must include the shared extraction routine by reference, not by copying its content

## Implementation Notes
- Add a new step between Post-Completion and the final summary output
- The trigger context for run-tasks: task outcomes include task records, code diffs, and any architectural decisions made during execution
