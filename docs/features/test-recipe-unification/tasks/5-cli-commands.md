---
id: "5"
title: "Update CLI command text and help messages"
priority: "P1"
estimated_time: "30min-1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 5: Update CLI command text and help messages

## Description

更新 `init.go` 中 init wizard 的提示文案，以及 `test/test.go` 中所有 `e2e-test` 引用和帮助文案，使其反映新的 recipe 命名模型。

## Reference Files
- `proposal.md#Proposed-Solution` — defines recipe naming: unit-test, test, test-setup, probe; retires e2e-test, e2e-setup, e2e-verify
- `proposal.md#Constraints-&-Dependencies` — init wizard prompt text requirements
- `proposal.md#Success-Criteria` — all e2e-test references in test.go updated to test, help text synchronized

## Acceptance Criteria
- `init.go` wizard prompts reference `unit-test` and `test` (not `e2e-test`)
- `test/test.go` all `e2e-test` references updated to `test`
- Help text accurately describes the two-layer test model

## Implementation Notes
- These are text-only changes in CLI command source files
- Ensure help text explains the distinction between `unit-test` (language-level) and `test` (surface-level advanced tests)
