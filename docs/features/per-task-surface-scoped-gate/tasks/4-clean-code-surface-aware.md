---
id: "4"
title: "clean-code skill 增加 surface-aware gate recipe 指引"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 4: clean-code skill 增加 surface-aware gate recipe 指引

## Description
`clean-code` skill在 Step 3 运行 `just unit-test` 验证 cleanup 无回归。在多 surface 项目中，generic `unit-test` 可能不存在（仅有 `<key>-unit-test` prefixed recipe），或运行全量测试导致不相关的 surface 失败阻塞当前 cleanup。需要更新 SKILL.md 使其 surface-aware：优先使用当前任务的 surface-key 对应的 prefixed recipe。

## Reference Files
- `docs/proposals/per-task-surface-scoped-gate/proposal.md` — Proposed Solution, Key Scenarios
- `plugins/forge/skills/clean-code/SKILL.md`: Step 3 Quality gate 的 unit-test 指引需增加 surface-aware 逻辑 (ref: Key Scenarios)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/clean-code/SKILL.md` | Step 3 Quality gate 段落增加 surface-aware prefixed recipe 检测逻辑 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] SKILL.md Step 3 指引 agent 检测当前任务是否有 surface-key，有则优先运行 `just <key>-unit-test`，不存在时回退到 `just unit-test`
- [ ] 无 surface-key 时行为与改动前一致（运行 `just unit-test`，不存在则 skip）
- [ ] 修改前必须先加载 `docs/conventions/forge-distribution.md`（CLAUDE.md MANDATORY 规则）

## Hard Rules
- 修改 `plugins/forge/` 下的文件前必须先加载 `docs/conventions/forge-distribution.md`

## Implementation Notes
- 当前 SKILL.md lines 138-161 定义了 Step 3 Quality gate 逻辑
- 检测逻辑可参考 Task 1 中 `resolvePrefixedRecipe()` 的 fallback 模式：先探测 prefixed recipe 存在性，不存在则回退 generic
- clean-code 在任务执行上下文中被调用，agent 可从任务 frontmatter 获取 surface-key
- `gen-test-scripts` skill 和 `run-tests` skill 中也有类似的 `just compile` / `just unit-test` 引用，但它们属于测试生成/运行流程，不在本次 scope 内
