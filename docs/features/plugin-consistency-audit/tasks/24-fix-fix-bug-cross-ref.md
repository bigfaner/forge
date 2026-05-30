---
id: "24"
title: "Fix: fix-bug cross-reference precision"
priority: "P2"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 24: Fix: fix-bug cross-reference precision

## Description
fix-bug command line 258 引用 "domain-to-file mapping from /consolidate-specs skill Step 5"，但 consolidate-specs Step 5 的实际内容是 "Generate Preview Files + Detect Overlaps"，与 domain-to-file mapping 无关。实际 mapping 在 `rules/overlap-detection.md`。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#CMD-03`: P2 INCOMPLETE, 引用不精确 (source: Report 05)
- `plugins/forge/commands/fix-bug.md`: 需修改的 line 258 附近 (source: audit finding)
- `plugins/forge/skills/consolidate-specs/rules/overlap-detection.md`: 正确的 domain-to-file mapping 位置 (source: cross-reference)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | 修正 cross-reference 为 "domain-to-decision-file mapping from /consolidate-specs rules/overlap-detection.md" |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] fix-bug 中对 domain-to-file mapping 的引用指向正确的文件路径（rules/overlap-detection.md）
- [ ] 不再引用 "Step 5"（因为 Step 5 的内容与 mapping 无关）

## Hard Rules
- 仅修改 `plugins/forge/commands/fix-bug.md` 的 cross-reference 部分
