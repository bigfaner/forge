---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the surface-scalar-dot-fix feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-switch-to-text-mode
- [ ] init-justfile 从 `forge surfaces --json` 切换到 `forge surfaces`（文本模式），scalar 形式（文本输出无 `=`）生成无 prefix recipe：`test`、`build`、`dev`、`teardown`
- [ ] run-tests 从 `forge surfaces --json` 切换到 `forge surfaces`（文本模式），scalar 形式调用 `just test` 而非 `just <key>-test`
- [ ] test-guide 从直接读取 config.yaml 切换到 `forge surfaces`（文本模式），统一数据源与其他 skill 一致
- [ ] breakdown-tasks 和 quick-tasks 的 Surface-Key/Type Inference 从 `--json` 切换到 `forge surfaces`（文本模式），scalar 形式下 surface-key 留空、surface-type 为 type 值
- [ ] Named key 形式（如 `app=tui`）下，recipe 名为 `<key>-<verb>`（如 `app-test`），行为与当前一致
- [ ] 所有 skill 使用统一解析规则：每行按 `=` 分割，无 `=` 则为 scalar（只有 type，无 key）；有 `=` 则左侧为 key，右侧为 type


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/surface-scalar-dot-fix/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/surface-scalar-dot-fix/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria
- [ ] All doc task deliverables reviewed against acceptance criteria summary above
- [ ] Review findings reported
