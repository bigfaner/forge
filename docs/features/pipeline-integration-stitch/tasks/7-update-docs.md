---
id: "7"
title: "更新文档：README、ARCHITECTURE、task-lifecycle"
priority: "P2"
estimated_time: "30min"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 7: 更新过时文档引用

## Description

更新 README.md 和 ARCHITECTURE.md 中过时的 `T-eval-doc` 引用为 `T-review-doc`。更新 `docs/business-rules/task-lifecycle.md` 中系统类型列表：移除 `test.gen-cases`、`test.eval-cases`、`doc.eval`，添加 `test.gen-journeys`、`test.gen-contracts`、`doc.review`、`eval.journey`、`eval.contract`。

## Reference Files
- `docs/proposals/pipeline-integration-stitch/proposal.md` — Source proposal
- `README.md` — T-eval-doc 引用
- `docs/ARCHITECTURE.md` — T-eval-doc 引用
- `docs/business-rules/task-lifecycle.md` — 系统类型列表

## Acceptance Criteria
- [ ] `README.md` 中 `T-eval-doc` 替换为 `T-review-doc`
- [ ] `ARCHITECTURE.md` 中 `T-eval-doc` 替换为 `T-review-doc`
- [ ] `task-lifecycle.md` 系统类型列表反映当前实际类型（含 `doc.review`、`test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract`，无 `doc.eval`、`test.gen-cases`、`test.eval-cases`）
- [ ] `grep -r "T-eval-doc\|doc\.eval" README.md docs/ARCHITECTURE.md docs/business-rules/` 返回零结果
