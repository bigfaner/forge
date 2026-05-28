---
id: "3"
title: "Slim task-executor Execution Protocol"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Slim task-executor Execution Protocol

## Description
Merge the task-executor agent's Execution Protocol from 11 steps to ≤8 steps by combining logically overlapping steps and compacting the output format. The execution semantics must remain unchanged — only the presentation is simplified.

Merge plan:
- Steps 4/5/6 (prompt fetch/parse/inject) → 1 step "fetch and inject prompt"
- Steps Retry Strategy + Complex Error Pause Flow → 1 step "error handling" (deduplicated shared "pause and notify user" semantics)
- Output format from multi-line description to compact single-line summary

## Reference Files
- plugins/forge/agents/task-executor.md: Execution Protocol section — merge steps 4/5/6 and Retry/Error steps (source: proposal.md#Key-Scenarios-5)
- docs/conventions/forge-distribution.md: task-executor.md is a distributed plugin file (source: proposal.md#Constraints-&-Dependencies)

## Acceptance Criteria
- [ ] Execution Protocol steps reduced from 11 to ≤8
- [ ] Steps 4/5/6 (prompt fetch/parse/inject) merged into single "fetch and inject prompt" step
- [ ] Retry Strategy and Complex Error Pause Flow merged into single "error handling" step
- [ ] Output format changed from multi-line description to compact single-line summary
- [ ] All original behavioral semantics preserved — no functional step lost

## Hard Rules
- **Distribution model**: task-executor.md is distributed via plugin cache — modify source at `plugins/forge/agents/task-executor.md` only

## Implementation Notes
- Expected 8-step structure: Initialize → Task validation → Fetch and inject prompt → Execute task → Output formatting → Result recording → Error handling → Completion notification
- Estimated ~30 line reduction
