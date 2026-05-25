---
status: "completed"
started: "2026-05-25 01:00"
completed: "2026-05-25 01:03"
time_spent: "~3m"
---

# Task Record: 13 Hook Unix parameter validation and misc file relocations

## Summary
Three P2 fixes: (1) added empty-parameter guard to run-hook.cmd Unix section; (2) moved validate-ux-pipeline.md from rubrics/ to rules/ and updated pre-processing.md reference; (3) evaluated all 18 commands — all retained with documented rationale

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/run-hook.cmd
- plugins/forge/skills/eval/rules/pre-processing.md

### Key Decisions
无

## Document Metrics
1 file moved (rubrics→rules), 1 reference updated, 1 guard added, 18 commands evaluated

## Referenced Documents
- docs/proposals/v3-release-audit/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] Hook 参数格式已审查，问题已记录或修复
- [x] validate-ux-pipeline.md 位于 rules/ 目录
- [x] 未暴露 skill 的 command 入口已评估，决策已记录

## Notes
Hook audit: session-start and debug scripts already had set -euo pipefail and proper error handling. Only run-hook.cmd Unix section lacked $1 empty check. Command evaluation: 7 eval-* commands delegate to unified eval skill (intentional), 8 independent commands each serve clear purpose. Conclusion: all commands retained. validate-ux-pipeline.md was a pipeline description (not a rubric), correctly moved to rules/.
