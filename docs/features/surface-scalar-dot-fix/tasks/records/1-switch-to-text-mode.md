---
status: "completed"
started: "2026-06-03 22:42"
completed: "2026-06-03 22:49"
time_spent: "~7m"
---

# Task Record: 1 统一所有 skill 使用 forge surfaces 文本模式

## Summary
Unified all 5 skills to use forge surfaces text mode instead of --json or config.yaml parsing, with consistent parsing rules for scalar and named surfaces

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/test-guide/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
5 files modified, 6 AC all met, unified parsing rule across all skills

## Referenced Documents
- docs/proposals/surface-scalar-dot-fix/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] init-justfile switches from --json to text mode, scalar generates prefix-less recipes (test, build, dev, teardown)
- [x] run-tests switches from --json to text mode, scalar calls just test not just <key>-test
- [x] test-guide switches from config.yaml to forge surfaces text mode, unified data source
- [x] breakdown-tasks and quick-tasks Surface-Key/Type Inference switches to text mode, scalar: surface-key empty, surface-type = type value
- [x] Named key form (e.g. app=tui) produces <key>-<verb> recipes (e.g. app-test), behavior unchanged
- [x] All skills use unified parsing rule: split on =, no = means scalar, = means named

## Notes
Hard Rules respected: only 5 SKILL.md files modified, no Go CLI code changed, no forge surfaces output format changes. Rule files under skills/*/rules/ were intentionally left unchanged as they are outside scope.
