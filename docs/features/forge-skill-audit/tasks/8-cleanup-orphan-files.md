---
id: "8"
title: "Clean up test-guide orphan files (M-4)"
priority: "P2"
estimated_time: "15m"
dependencies: [5]
type: "doc"
mainSession: false
---

# 8: Clean up test-guide orphan files (M-4)

## Description

test-guide 存在 2 个孤儿规则文件（draft-generation.md 和 pattern-extraction.md），SKILL.md 和其他 rules 文件均未引用。内容完整，暗示曾是流程一部分后被重构移除但文件残留。需移至 `_deprecated/` 目录或补充引用。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — M-4: test-guide 存在 2 个孤儿规则文件
- `plugins/forge/skills/test-guide/rules/draft-generation.md`: Orphan file to move (ref: M-4: test-guide 存在 2 个孤儿规则文件)
- `plugins/forge/skills/test-guide/rules/pattern-extraction.md`: Orphan file to move (ref: M-4: test-guide 存在 2 个孤儿规则文件)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/test-guide/rules/_deprecated/draft-generation.md` | Moved from rules/ |
| `plugins/forge/skills/test-guide/rules/_deprecated/pattern-extraction.md` | Moved from rules/ |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/test-guide/rules/draft-generation.md` | Moved to _deprecated/ |
| `plugins/forge/skills/test-guide/rules/pattern-extraction.md` | Moved to _deprecated/ |

## Acceptance Criteria
- [ ] `draft-generation.md` 和 `pattern-extraction.md` 已移至 `skills/test-guide/rules/_deprecated/` 目录
- [ ] SKILL.md 和其他 rules 文件无引用断裂（本就无引用）

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes

## Implementation Notes
- 与 eval 系统的 `_deprecated/` 惯例一致（eval/rules/_deprecated/freeform-injection.md）
- 确认 `_deprecated/` 目录是否存在，不存在则创建
