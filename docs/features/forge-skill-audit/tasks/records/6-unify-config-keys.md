---
status: "blocked"
started: "2026-06-10 19:28"
completed: "N/A"
time_spent: ""
---

# Task Record: 6 Unify auto.eval config key naming (M-1)

## Summary
M-1 config key naming: verified Go config reader does not support kebab-case queries (findFieldByYAMLTag uses exact yaml tag match). Per Hard Rule, did not rename keys. brainstorm/write-prd already kebab-compatible. Unified ui-design TODO marker from HTML comment to bash-script comment format and aligned eval auto-run to consistent bash script template. Blocked on AC-1: cannot rename uiDesign/techDesign until Go config reader adds alias support.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/ui-design/SKILL.md

### Key Decisions
无

## Document Metrics
1 file modified; 2 files verified as already kebab-compatible; 2 files with TODO(M-1) markers confirmed

## Referenced Documents
- docs/proposals/forge-skill-audit/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [ ] All auto.eval config keys unified to kebab-case in skill markdown
- [x] Implementation notes mark Go config reader alias compat as follow-up task

## Notes
Blocked: AC-1 cannot be satisfied because Go config reader (findFieldByYAMLTag in config_reflect.go) matches yaml tags exactly. Hard Rule requires verifying Go reader support before renaming — reader does not support kebab-case. ui-design eval auto-run was upgraded from natural-language to bash script template as improvement. Needs Go alias support (out of scope for this doc-only task) to unblock.
