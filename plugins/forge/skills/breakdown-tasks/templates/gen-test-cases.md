---
id: "T-test-1"
title: "生成 e2e 测试用例"
priority: "P1"
estimated_time: "1-2h"
dependencies: [{{LAST_BUSINESS_TASK_ID}}]
status: pending
---

# T-test-1: 生成 e2e 测试用例

## Description

调用 `/gen-test-cases` skill，从 PRD 验收标准生成结构化测试用例文档。

每个测试用例包含：
- Source: PRD 中的具体验收标准
- Type: UI / API / CLI
- Target: 测试目标路径（如 ui/login, api/auth）
- Test ID: 唯一标识符（如 ui/login/login-with-valid-credentials）
- Pre-conditions, Steps, Expected, Priority

## Reference Files

- `prd/prd-spec.md` — PRD 规格文档
- `prd/prd-user-stories.md` — 用户故事（含 Given/When/Then 验收标准）
- `prd/prd-ui-functions.md` — UI 功能要点（可选）

## Acceptance Criteria

- [ ] `testing/test-cases.md` 文件已创建
- [ ] 每个测试用例包含 Target 和 Test ID 字段
- [ ] 所有测试用例可追溯到 PRD 验收标准
- [ ] 测试用例按类型分组（UI → API → CLI）

## User Stories

No direct user story mapping. This is a standard test generation task.

## Implementation Notes

1. 运行 `/gen-test-cases` skill
2. 验证生成的 `testing/test-cases.md` 包含 Target 和 Test ID 字段
3. 如果 PRD 无 UI/API/CLI 需求，将任务标记为 skipped 并说明原因
