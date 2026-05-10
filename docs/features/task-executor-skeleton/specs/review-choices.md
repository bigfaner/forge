---
feature: "task-executor-skeleton"
reviewed: "2026-05-10"
status: pending-review
---

# Review Choices

## Approved for Integration

(none yet — awaiting user review)

## CROSS Items Pending Review

### Business Rules

- BIZ-003 (Terminal States) -> `docs/business-rules/task-lifecycle.md`
  Rule: Task execution has exactly two terminal states: `completed` and `failed`. No `partial` state.
- BIZ-007 (Backward Compatibility Fallback) -> `docs/business-rules/task-lifecycle.md`
  Rule: Task files without `## Execution Workflow` fallback to TDD + Quality Gate via default template.
- BIZ-009 (Agent Timeout) -> `docs/business-rules/task-lifecycle.md`
  Rule: Agent timeout triggers `failed` status with timeout information recorded.
- BIZ-010 (Quality Gate Safety Net) -> `docs/business-rules/task-lifecycle.md`
  Rule: `task record` runs quality gate pre-check for ALL tasks uniformly; no bypass.

### Technical Specs

- TECH-003 (Workflow Content Model W1-W5) -> `docs/conventions/task-templates.md`
  Requirement: Every `## Execution Workflow` must satisfy W1-W5 (numbered, concrete, success criteria, failure handling, stop condition).
- TECH-008 (Structured Error Handoff) -> `docs/conventions/agent-cli-contracts.md`
  Requirement: Agent-to-task-cli errors flow via structured CLI args (`task record --status`), not freeform text parsing.
- TECH-009 (Error Propagation Channels) -> `docs/conventions/agent-cli-contracts.md`
  Requirement: Go errors use `*AIError` structs; agent prompts use step output strings. No error codes cross the boundary.

## Skipped

(none)

## Related Existing Entries

No overlaps detected. Decision files (error-handling.md, architecture.md, interface.md, testing.md, data-model.md) are empty. No lesson tags match the extracted domains.
