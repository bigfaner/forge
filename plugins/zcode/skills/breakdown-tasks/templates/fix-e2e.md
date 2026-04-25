---
id: "fix-e2e-{{round}}-{{index}}"
title: "修复 e2e 测试失败: {{test_name}}"
priority: "P0"
estimated_time: "30min-2h"
dependencies: []
status: pending
---

# fix-e2e-{{round}}-{{index}}: 修复 e2e 测试失败

## Description

e2e 测试失败，需要分析失败原因并修复代码。

## Reference Files

- `testing/results/latest.md` — 测试结果概览
- `testing/results/failures/failure-{{test_case_id}}.md` — 本任务对应的失败详情
- `testing/test-cases.md` — 测试用例文档
- `testing/scripts/` — 测试脚本目录

## Acceptance Criteria

- [ ] 已读取 `testing/results/latest.md` 了解失败概览
- [ ] 已读取 `testing/results/failures/failure-{{test_case_id}}.md` 了解具体失败详情
- [ ] 已定位失败的根本原因
- [ ] 已修复代码或测试脚本
- [ ] 本地验证测试通过（可选）
- [ ] `task all-completed` 再次运行时测试通过

## User Stories

No direct user story mapping. This is a test fix task.

## Implementation Notes

1. 读取 `testing/results/latest.md` 查看失败概览
2. 读取 `testing/results/failures/failure-{{test_case_id}}.md` 了解具体失败详情
3. 分析失败原因：
   - 代码逻辑错误？
   - 测试脚本问题？
   - 环境配置问题？
4. 修复问题
5. 如果需要，可从 `testing/scripts/` 目录运行 `npm run test:all` 本地验证
6. 完成后执行 `task record` 记录修复内容

## Context

这是第 {{round}} 轮修复尝试（共 3 轮）。
