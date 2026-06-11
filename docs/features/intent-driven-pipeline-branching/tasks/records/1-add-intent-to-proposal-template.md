---
status: "completed"
started: "2026-05-29 16:19"
completed: "2026-05-29 16:21"
time_spent: "~2m"
---

# Task Record: 1 Add intent field to proposal template and brainstorm skill

## Summary
Added intent field to proposal template frontmatter and intent inference step to brainstorm SKILL.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/brainstorm/templates/proposal.md
- plugins/forge/skills/brainstorm/SKILL.md

### Key Decisions
无

## Document Metrics
template: +1 frontmatter field; SKILL.md: +24 lines (Step 4.5 + process flow update)

## Referenced Documents
- docs/proposals/intent-driven-pipeline-branching/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] proposal.md frontmatter contains intent field after status
- [x] brainstorm SKILL.md contains intent inference step with task type mapping rules
- [x] brainstorm uses AskUserQuestion to confirm intent, user can override
- [x] coding.fix intent heuristic: new user-observable behavior -> new-feature, internal only -> refactor
- [x] Mixed content proposal handled by primary intent assessment with user override

## Notes
Process flow updated to include Infer intent step. Step 4.5 inserted between Define Scope and Write Proposal to ensure intent is confirmed before proposal generation.
