---
status: "completed"
started: "2026-05-24 00:29"
completed: "2026-05-24 00:32"
time_spent: "~3m"
---

# Task Record: 6 gen-journeys SKILL.md 适配：proposal.md 输入 + AUTO_COMMIT 模式

## Summary
Modified gen-journeys SKILL.md to support proposal.md as alternative input (Proposal Mode), AUTO_COMMIT non-interactive mode, and minimum information checks for proposal.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-journeys/SKILL.md

### Key Decisions
无

## Document Metrics
1 file modified, Prerequisites section expanded with dual-mode support, Step 1/2/4/6 updated for Proposal Mode and AUTO_COMMIT

## Referenced Documents
- docs/proposals/auto-gen-journeys-contracts/proposal.md
- plugins/forge/skills/gen-journeys/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] Prerequisites distinguishes PRD Mode and Proposal Mode
- [x] Proposal Mode scope + success criteria are mandatory (abort on missing)
- [x] Key Scenarios missing degrades to smoke-level quality=low Journey
- [x] Step 6 AUTO_COMMIT conditional path added
- [x] Manual /gen-journeys behavior unchanged (Interactive Mode default)
- [x] proposal.md upgraded from optional to conditionally required

## Notes
Core Journey generation logic (Steps 1-5 workflow, risk classification, file generation, validation) preserved as required by Hard Rules. AUTO_COMMIT path still mandates Step 5 validation.
