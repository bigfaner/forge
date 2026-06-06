---
id: "10"
title: "Update hooks/guide.md and submit-task record format"
priority: "P1"
estimated_time: "30min"
dependencies: [1]
type: "doc"
mainSession: false
---

# 10: Update hooks/guide.md and submit-task record format

## Description
hooks/guide.md 中有测试目录约定声明，submit-task 的 record-format-test.md 中有示例路径，两者都需要更新以反映新的目录结构。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Proposed Solution, Scope

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| plugins/forge/hooks/guide.md | 测试目录约定声明更新 |
| plugins/forge/skills/submit-task/data/record-format-test.md | 示例路径更新 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] hooks/guide.md 测试目录约定声明反映 surface-key 自适应规则
- [ ] submit-task/data/record-format-test.md 示例路径与新目录结构一致

## Implementation Notes
hooks/guide.md 在每个会话开始时注入上下文，是 agent 的主要参考之一。确保其中的测试目录描述准确。
