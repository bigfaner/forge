---
status: "completed"
started: "2026-05-27 20:13"
completed: "2026-05-27 20:22"
time_spent: "~9m"
---

# Task Record: 11 同步更新 OVERVIEW 和 WORKFLOW 文档（含中文版）

## Summary
替换 OVERVIEW.md、WORKFLOW.md 及其中文版中所有 'e2e' 泛用为 surface-specific 术语，更新 graduation/staging 描述为 tag-based promotion，移除 'profile' 旧术语引用

## Changes

### Files Created
无

### Files Modified
- forge-cli/docs/OVERVIEW.md
- forge-cli/docs/WORKFLOW.md
- forge-cli/docs/OVERVIEW.zh.md
- forge-cli/docs/WORKFLOW.zh.md

### Key Decisions
无

## Document Metrics
4 files modified, ~25 substitutions across all files

## Referenced Documents
- docs/proposals/test-pipeline-consistency-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] OVERVIEW.md 和 OVERVIEW.zh.md 中 'e2e' 泛用替换为 surface-specific 术语
- [x] WORKFLOW.md 和 WORKFLOW.zh.md 中 'e2e' 泛用替换为 surface-specific 术语
- [x] graduation/staging 描述已更新为 tag-based promotion
- [x] 'profile' 旧术语引用已移除
- [x] grep for 'tests/e2e' 在中文版文档中返回 0 结果

## Notes
Spec-code conflict detected: 'forge e2e validate-specs' command no longer exists in code; quality gate flow simplified to compile -> test -> regression (no separate e2e-setup/probe/test-e2e steps); 'Profile detection signals' section replaced with 'Surface type detection signals'; 'embedded profiles' index build step replaced with 'detected surfaces'; --test-profiles flag removed
