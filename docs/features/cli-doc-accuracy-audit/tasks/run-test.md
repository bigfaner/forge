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

Execute staged test scripts for the cli-doc-accuracy-audit feature.

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/cli-doc-accuracy-audit/testing/                                 # journeys
ls docs/features/cli-doc-accuracy-audit/testing/<journey>/contracts/              # contracts
```

## Feature Context
- Scope: .

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.

Type: **cli**

## Acceptance Criteria
- [ ] 所有 CLI functional test 脚本执行完成且全部通过
- [ ] 测试失败时已定位根因并应用最小修复
