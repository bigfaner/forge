---
status: "completed"
started: "2026-05-27 00:28"
completed: "2026-05-27 00:35"
time_spent: "~7m"
---

# Task Record: 5 Update 4 skill files with config check bash template

## Summary
Updated 3 skill files (write-prd, tech-design, ui-design) with unified config check bash template replacing AskUserQuestion eval trigger. brainstorm was already updated in a prior task.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md
- plugins/forge/skills/ui-design/SKILL.md

### Key Decisions
无

## Document Metrics
3 files updated, 1 already done (brainstorm). 4/4 skills now use identical config check template.

## Referenced Documents
- docs/proposals/auto-eval-config/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
completed

## Acceptance Criteria
- [x] brainstorm auto.eval.proposal=true skips AskUserQuestion, runs eval-proposal
- [x] brainstorm auto.eval.proposal=false keeps AskUserQuestion
- [x] write-prd auto.eval.prd=true skips AskUserQuestion, runs eval-prd
- [x] write-prd auto.eval.prd=false keeps AskUserQuestion
- [x] ui-design auto.eval.uiDesign=true skips AskUserQuestion, runs eval-ui
- [x] ui-design auto.eval.uiDesign=false keeps AskUserQuestion
- [x] tech-design auto.eval.techDesign=true skips AskUserQuestion, runs eval-design
- [x] tech-design auto.eval.techDesign=false keeps AskUserQuestion
- [x] 4 skills use identical config check template (code review verified)
- [x] CLI unavailable (non-zero exit) falls back to AskUserQuestion

## Notes
brainstorm/SKILL.md was already updated in task 1 (verified in context). ui-design changed from unconditional auto-run to config-driven; default values (quick:true full:true) preserve original behavior. All templates use EXTREMELY-IMPORTANT annotation per Hard Rules.
