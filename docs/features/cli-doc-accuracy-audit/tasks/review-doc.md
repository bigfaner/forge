---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the cli-doc-accuracy-audit feature (quick mode).

## Acceptance Criteria
- [ ] Task 1 的 4 条 AC 已通过人工 review 确认（guide.md 命令名和描述准确）
- [ ] Task 2 的 4 条 AC 已通过人工 review 确认（新增命令与 CLI --help 输出一致）
- [ ] review 结果记录到任务执行记录中

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-update-guide-commands
- [ ] guide.md 中 `forge task validate-index` 替换为 `forge task validate [file]`，且 `forge task validate-index` 在 CLI 中返回 "unknown command" 错误
- [ ] guide.md 中 `forge quality-gate` 描述准确反映实际行为（含 fix task 自动创建、retry-once、docs-only 跳过）
- [ ] guide.md 中 `forge cleanup` 描述从 "clean stale artifacts" 改为具体行为说明（包含 blocked/suspended/rejected 状态的清理）
- [ ] guide.md 中 `forge task submit` 描述补充 `--quiet` 标志


### 2-add-guide-entries
- [ ] guide.md 新增 `forge task query <id-or-key>` 命令描述，与 `forge task query --help` 输出一致
- [ ] guide.md 新增 `forge task check-deps` 命令描述，与 `forge task check-deps --help` 输出一致
- [ ] guide.md 新增 `forge feature list` 命令描述，与 `forge feature list --help` 输出一致
- [ ] guide.md 中 `forge task list` 描述补充 `--tree` 标志


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/cli-doc-accuracy-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/cli-doc-accuracy-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
