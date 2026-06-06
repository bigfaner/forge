---
id: "9"
title: "Update test-guide rules"
priority: "P1"
estimated_time: "1h"
dependencies: [8]
type: "doc"
mainSession: false
---

# 9: Update test-guide rules

## Description
test-guide 的 4 个 rule 文件控制约定文件的生成和验证逻辑，需确保它们与新目录结构一致。

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
| plugins/forge/skills/test-guide/rules/convention-structure.md | 目录结构定义反映 surface-key 路径 |
| plugins/forge/skills/test-guide/rules/draft-generation.md | 生成的约定内容使用正确路径 |
| plugins/forge/skills/test-guide/rules/pattern-extraction.md | 检查并更新路径模式 |
| plugins/forge/skills/test-guide/rules/signal-detection.md | 检查并更新路径信号 |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] convention-structure.md 目录结构定义包含 surface-key 分区规则
- [ ] draft-generation.md 生成的约定内容路径正确
- [ ] pattern-extraction.md 和 signal-detection.md 中的路径模式已更新（如有引用）

## Implementation Notes
需先检查这 4 个 rule 文件是否实际引用了 `tests/<journey>/` 路径。如果未引用，此任务可能只需要确认兼容性而无需修改。
