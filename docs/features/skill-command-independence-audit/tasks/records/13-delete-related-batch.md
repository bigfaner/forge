---
status: "completed"
started: "2026-06-04 01:16"
completed: "2026-06-04 01:17"
time_spent: "~1m"
---

# Task Record: 13 Delete Related sections in consolidate-specs, run-tests, ui-design

## Summary
Deleted Related Skills/Integration sections from consolidate-specs, run-tests, and ui-design SKILL.md files. All deleted content was purely pipeline upstream/downstream information already inferable from each document's body.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/consolidate-specs/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/ui-design/SKILL.md

### Key Decisions
无

## Document Metrics
3 sections removed (~25 lines total), 0 unique information lost

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] 三个 skill 的 Related Skills、Integration、References 章节已删除
- [x] 删除的内容均可从正文中隐含推断，无独有信息丢失
- [x] 仅修改以下文件：consolidate-specs/SKILL.md、run-tests/SKILL.md、ui-design/SKILL.md

## Notes
consolidate-specs: removed ## Related Skills (4 entries). run-tests: removed ## Related Skills (2 entries + tag promotion note). ui-design: removed ## Integration (3 entries). All relationships are already expressed in Prerequisites, When to Use, and step-level invocations.
