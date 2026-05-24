---
status: "completed"
started: "2026-05-25 00:16"
completed: "2026-05-25 00:21"
time_spent: "~5m"
---

# Task Record: 6 Full README.md rewrite for v3.0.0

## Summary
Full README.md rewrite for v3.0.0: corrected version badge (2.16.1->5.6.0), skill count (17->21), command count (17->18), agent count (3->1), task types (13 legacy->21 dot-notation), removed ghost web/raycast references, fixed Go version requirement (1.26.1+->1.25+), added complete command reference table aligned with forge --help, corrected all directory paths (task-cli->forge-cli, removed web/).

## Changes

### Files Created
无

### Files Modified
- README.md

### Key Decisions
无

## Document Metrics
1 file rewritten, 7 Critical errors fixed, 6 Major errors fixed, 1 Minor error fixed

## Referenced Documents
- docs/proposals/v3-release-audit/proposal.md
- docs/features/v3-release-audit/tasks/6-rewrite-readme.md
- docs/ARCHITECTURE.md
- docs/features/v3-release-audit/tasks/records/4-fix-architecture-md.md
- docs/features/v3-release-audit/tasks/records/5-fix-cli-references.md

## Review Status
final

## Acceptance Criteria
- [x] 版本号 = cat forge-cli/scripts/version.txt
- [x] 技能计数 = ls plugins/forge/skills/ | wc -l (21)
- [x] 任务类型表覆盖所有 21 种 dot-notation 类型
- [x] 命令速查与 forge --help 一一对应
- [x] 无幽灵命令（web/raycast 引用已移除）
- [x] 安装步骤指向正确 Go 版本要求 (1.25+)
- [x] 路径引用与实际目录匹配

## Notes
Task dependencies (4, 5) both completed. All factual claims cross-verified against live codebase. Version source is forge-cli/scripts/version.txt (plugins/forge/scripts/version.txt does not exist). Web/ directory reference removed as it does not exist in this worktree. Added forge completion and forge help commands to table for complete forge --help alignment.
