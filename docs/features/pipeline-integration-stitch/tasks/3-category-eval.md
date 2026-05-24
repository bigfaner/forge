---
id: "3"
title: "新增 CategoryEval + submit-task 验证分支 + category 测试"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 3: CategoryEval 类型分类

## Description

在 `category.go` 中新增 `CategoryEval` 常量，将 `eval.` 前缀映射到 `CategoryEval`（非 CategoryTest）。eval 任务是质量门控（review），不是测试生成器，提交时需要 review 结论字段（summary/findings/severity）而非测试证据（testsPassed/coverage）。在 submit-task 验证逻辑中为 CategoryEval 添加专用验证分支。

## Reference Files
- `docs/proposals/pipeline-integration-stitch/proposal.md` — Source proposal
- `forge-cli/pkg/task/category.go` — 当前分类逻辑
- `forge-cli/pkg/task/category_test.go` — 现有测试

## Acceptance Criteria
- [ ] `category.go` 新增 `CategoryEval = "eval"` 常量
- [ ] `CategoryForType("eval.journey")` 和 `CategoryForType("eval.contract")` 返回 `CategoryEval`
- [ ] submit-task 验证逻辑中为 CategoryEval 添加验证分支，接受 review 类字段（summary/findings/severity）
- [ ] submit-task 对 eval 任务**拒绝**仅含测试字段（testsPassed/coverage 无 summary/findings）的提交
- [ ] `category_test.go` 新增 CategoryEval 正向用例、负向用例、边界用例
- [ ] 遍历所有 switch-on-category 分支确认 CategoryEval 不落入 default/else
- [ ] 所有现有测试通过

## Hard Rules
- CategoryEval 为新增类别，不得修改 CategoryTest 现有行为
- 负向用例验证必须存在：eval 任务仅提交测试字段时被拒绝

## Implementation Notes
- eval 任务提交时需要的证据字段是 review 结论（summary/findings/severity），不是测试结果
- 注意 submit-task 验证逻辑中 switch-on-category 的所有分支都需覆盖
