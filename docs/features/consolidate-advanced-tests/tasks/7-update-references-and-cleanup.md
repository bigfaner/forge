---
id: "7"
title: "更新活跃引用并删除 forge-cli/tests/"
priority: "P2"
estimated_time: "1h"
complexity: "low"
dependencies: [2, 3, 4, 5, 6]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 7: 更新活跃引用并删除 forge-cli/tests/

## Description
所有 journey 迁移完成后，更新文档和技能规则中指向 `forge-cli/tests/` 的活跃引用，然后删除整个 `forge-cli/tests/` 目录。

## Reference Files
- `docs/proposals/consolidate-advanced-tests/proposal.md` — Scope > In Scope: 活跃引用更新、清理, Success Criteria
- `docs/conventions/forge-distribution.md` — 待更新引用
- `plugins/forge/skills/run-tests/rules/test-isolation.md` — 待更新引用
- `tests/test-suite-health/contracts/step-1-test-suite-health.md` — 待更新路径引用

## Acceptance Criteria
- [ ] `docs/conventions/forge-distribution.md` 中无 `forge-cli/tests/` 路径引用
- [ ] `plugins/forge/skills/run-tests/rules/test-isolation.md` 中无 `forge-cli/tests/` 路径引用
- [ ] `forge-cli/tests/` 目录完全删除，不存在于工作树中
- [ ] `grep -r 'forge-cli/tests' tests/` 返回空（源码中无残留引用）

## Implementation Notes
- 先更新引用，再删除目录，避免删除后无法对比
- `tests/test-suite-health/contracts/step-1-test-suite-health.md` 中如有 `forge-cli/tests/` 引用也需一并更新
- 删除前确认 `forge-cli/tests/` 中无其他代码（非测试文件）被 forge-cli 主代码引用
