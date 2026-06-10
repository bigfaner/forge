---
id: "17"
title: "Fix knowledgeSave format description consistency (MINOR-H2)"
priority: "P2"
estimated_time: "15m"
dependencies: [10]
type: "doc"
mainSession: false
---

# 17: Fix knowledgeSave format description consistency

## Description

auto.knowledgeSave 的输出格式描述在不同文件中措辞不一致：fix-bug.md 描述为 "plain text key:value pairs"，knowledge-extraction.md 描述为 "quick:<val> full:<val>"。虽然语义相同，但表述不一致。

## Reference Files
- `plugins/forge/commands/fix-bug.md`: Standardize format description
- `plugins/forge/skills/tech-design/rules/knowledge-extraction.md`: Standardize format description
- `plugins/forge/skills/write-prd/rules/knowledge-extraction.md`: Standardize format description

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | Unify knowledgeSave format description wording |
| `plugins/forge/skills/tech-design/rules/knowledge-extraction.md` | Unify knowledgeSave format description wording |
| `plugins/forge/skills/write-prd/rules/knowledge-extraction.md` | Unify knowledgeSave format description wording |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 3 个文件中 auto.knowledgeSave 的输出格式描述使用统一的措辞

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`

## Implementation Notes
- 统一为具体格式示例（如 `Output format: key:value pairs (e.g., "quick:true full:false")`），既描述格式又给出示例
