---
id: "4"
title: "Fix record-format type coverage (H-4)"
priority: "P1"
estimated_time: "30m"
dependencies: [1]
type: "doc"
mainSession: false
---

# 4: Fix record-format type coverage (H-4)

## Description

record-format-coding.md 列出 `doc.fix` 作为 coding category 任务类型，但 CLI 的 `CategoryForType` 将 `doc.fix` 归类为 `doc` category（`strings.HasPrefix("doc")`）。导致记录格式文档与实际运行时行为矛盾。需从 coding 格式移除 `doc.fix` 并在 doc 格式中补充。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — H-4: record-format-coding.md 错误列出 doc.fix, Proposed Solution, Success Criteria
- `plugins/forge/skills/submit-task/data/record-format-coding.md`: Remove doc.fix (ref: H-4: record-format-coding.md 错误列出 doc.fix)
- `plugins/forge/skills/submit-task/data/record-format-doc.md`: Add doc.fix coverage (ref: H-4: record-format-coding.md 错误列出 doc.fix)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/data/record-format-coding.md` | Remove `doc.fix` from listed types |
| `plugins/forge/skills/submit-task/data/record-format-doc.md` | Add `doc.fix` to covered types |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] record-format-coding.md 不再列出 `doc.fix`
- [ ] record-format-doc.md 包含 `doc.fix` 覆盖

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes

## Implementation Notes
- 修复后运行回归验证：`grep -r "doc.fix" plugins/forge/skills/submit-task/data/` 确认仅出现在 record-format-doc.md 中
- `code-quality.simplify` 有类似的硬编码特殊映射问题（proposal 中记录为脆弱性分析），但不在本修复 scope 内
