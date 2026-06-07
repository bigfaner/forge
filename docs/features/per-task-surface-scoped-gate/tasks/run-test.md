---
id: "T-test-run"
title: "Run CLI Functional Test"
priority: "P1"
estimated_time: "30min-1h"
dependencies: ["T-test-gen-journeys"]
type: "test.run"
surface-key: "."
surface-type: "cli"
---

Execute staged test scripts for the per-task-surface-scoped-gate feature.

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/per-task-surface-scoped-gate/testing/                                 # journeys
ls docs/features/per-task-surface-scoped-gate/testing/<journey>/contracts/              # contracts
```

## Feature Context
- Scope: .

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.

Type: **cli**

## Acceptance Criteria
- [ ] 所有 staged test scripts 执行完毕并返回通过结果
- [ ] 失败的测试已定位根因并修复（minimal fix），修复后重跑通过
