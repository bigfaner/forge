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

这是第 {{round}} 轮修复尝试。修复步骤：

1. 读取 `testing/results/latest.md` 查看失败概览
2. 读取 `testing/results/failures/failure-{{test_case_id}}.md` 了解具体失败详情
3. 定位根本原因（代码逻辑 / 测试脚本 / 环境配置）
4. 修复并验证

## Reference Files

- `testing/results/latest.md` — 测试结果概览
- `testing/results/failures/failure-{{test_case_id}}.md` — 失败详情
- `testing/test-cases.md` — 测试用例文档
- `testing/scripts/` — 测试脚本目录

## Acceptance Criteria

- [ ] 已定位失败的根本原因
- [ ] 已修复代码或测试脚本
- [ ] 单元测试全部通过
