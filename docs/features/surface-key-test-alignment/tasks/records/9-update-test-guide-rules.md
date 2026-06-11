---
status: "completed"
started: "2026-06-06 13:54"
completed: "2026-06-06 13:55"
time_spent: "~1m"
---

# Task Record: 9 Update test-guide rules

## Summary
Updated test-guide rule files to reflect surface-key directory structure: draft-generation.md Top-Level Index paths now use surfaceKey-aware rule; pattern-extraction.md search priority now includes multi-surface paths. convention-structure.md was already correct; signal-detection.md does not reference test directory paths.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/test-guide/rules/draft-generation.md
- plugins/forge/skills/test-guide/rules/pattern-extraction.md

### Key Decisions
无

## Document Metrics
2 files modified, 2 files verified unchanged, 0 files created

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] convention-structure.md 目录结构定义包含 surface-key 分区规则
- [x] draft-generation.md 生成的约定内容路径正确
- [x] pattern-extraction.md 和 signal-detection.md 中的路径模式已更新（如有引用）

## Notes
convention-structure.md already had correct surface-key paths (no change needed). signal-detection.md does not reference test directory paths (no change needed). Only draft-generation.md and pattern-extraction.md required updates.
