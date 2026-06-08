---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["2", "3"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the test-pipeline-interleaved feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-harden-testrun-template

- [ ] 模板 AcceptanceCriteria 部分包含 AC：所有测试用例必须通过（不能有 skip 或预期失败）
- [ ] 模板 AcceptanceCriteria 部分包含 AC：必须是真实测试（验证实际功能行为），不能是占位符或 always-pass 假测试


### 3-harden-testrun-prompt

- [ ] 模板 TASK-CONSTRAINTS 包含指令：确认是正式代码 bug 才能修改正式代码；测试脚本本身的 bug 可以修，但不能为了通过测试而篡改测试断言/逻辑
- [ ] 模板 TASK-CONSTRAINTS 包含指令：问题多时通过 `forge task add` 追加 fix 任务（而非在一个任务中修复所有问题），与 task-executor Pause Protocol 协调


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/test-pipeline-interleaved/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/test-pipeline-interleaved/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All acceptance criteria met
