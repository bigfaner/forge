---
id: "15"
title: "Fix: gen-contracts Fact Table format inconsistency"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 15: Fix: gen-contracts Fact Table format inconsistency

## Description
gen-contracts SKILL.md Step 2 描述 Fact Table 输出格式为 JSON（含 `source`, `confidence`, `updated_at` 等字段），但 `rules/code-reconnaissance.md` 描述 Fact Table 输出格式为 Markdown 表格（`| Key | Value | Source |`）。两种格式不一致，SKILL.md 引用的权威格式在 `forge-cli/pkg/facttable/facttable.go`。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-10`: P1 CONFLICT, Fact Table 格式不一致 (source: Report 04)
- `plugins/forge/skills/gen-contracts/rules/code-reconnaissance.md`: 需修改的 rules 文件 (source: audit finding)
- `plugins/forge/skills/gen-contracts/SKILL.md`: Step 2 的 JSON 格式作为权威参考 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/rules/code-reconnaissance.md` | 将 Fact Table 格式从 Markdown 表格对齐为 JSON 格式 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `rules/code-reconnaissance.md` 中的 Fact Table 格式使用 JSON（与 SKILL.md 一致）
- [ ] 或明确说明 Markdown 是中间展示格式，最终写入 `.forge/fact-table.json` 时转换为 JSON

## Hard Rules
- 仅修改 `plugins/forge/skills/gen-contracts/rules/code-reconnaissance.md`

## Implementation Notes
- 以 SKILL.md 的 JSON 格式为权威，更新 rules 文件
- 如果 Markdown 格式是中间步骤（AI 内部推理用），保留但标注清楚
