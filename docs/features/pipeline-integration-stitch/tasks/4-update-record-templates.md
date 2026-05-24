---
id: "4"
title: "更新 record-format-doc.md 和 record-format-test.md"
priority: "P1"
estimated_time: "30min"
dependencies: ["3"]
type: "doc"
mainSession: false
---

# 4: 更新 record-format 模板

## Description

更新 submit-task 的记录格式模板文件：`record-format-doc.md` 中 `doc.eval` → `doc.review`；`record-format-test.md` 移除已废弃类型（`test.gen-cases`、`test.eval-cases`、`test.gen-and-run`），添加新类型（`test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract`）。

## Reference Files
- `docs/proposals/pipeline-integration-stitch/proposal.md` — Source proposal
- `plugins/forge/skills/submit-task/data/record-format-doc.md`
- `plugins/forge/skills/submit-task/data/record-format-test.md`

## Acceptance Criteria
- [ ] `record-format-doc.md` 中 `doc.eval` 替换为 `doc.review`
- [ ] `record-format-test.md` 移除 `test.gen-cases`、`test.eval-cases`、`test.gen-and-run`
- [ ] `record-format-test.md` 添加 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract`
- [ ] `grep -r "doc.eval" plugins/forge/skills/submit-task/` 返回零结果
