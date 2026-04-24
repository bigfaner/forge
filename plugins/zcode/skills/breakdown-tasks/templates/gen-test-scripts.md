---
id: "T-test-2"
title: "生成 e2e 测试脚本"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["T-test-1"]
status: pending
---

# T-test-2: 生成 e2e 测试脚本

## Description

调用 `/gen-test-scripts` skill，从测试用例生成可执行的 TypeScript e2e 测试脚本。

生成的脚本使用：
- UI 测试：`agent-browser` CLI
- API 测试：Node.js 内置 `fetch`
- CLI 测试：`child_process.execSync`
- 测试框架：`node:test` + `node:assert`

## Reference Files

- `testing/test-cases.md` — 测试用例文档（由 T-test-1 生成）

## Acceptance Criteria

- [ ] `testing/scripts/package.json` 已创建
- [ ] `testing/scripts/helpers.ts` 已创建
- [ ] 至少一个 spec 文件已生成（ui.spec.ts / api.spec.ts / cli.spec.ts）
- [ ] 每个 test() 包含追溯注释 `// Traceability: TC-NNN → {PRD Source}`

## User Stories

No direct user story mapping. This is a standard test generation task.

## Implementation Notes

1. 运行 `/gen-test-scripts` skill
2. 验证 `testing/scripts/package.json` 存在
3. 如果 T-test-1 被跳过，此任务同步标记为 skipped
