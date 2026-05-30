---
id: "12"
title: "Fix: resolve init-justfile template vs SKILL.md design contradiction"
priority: "P1"
estimated_time: "45min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 12: Fix: resolve init-justfile template vs SKILL.md design contradiction

## Description
init-justfile SKILL.md Step 0 HARD-RULE 说 "Do NOT use framework-specific recipe templates. Generate from Convention content and LLM knowledge"。但 `templates/` 目录下存在 6 个 .just 模板文件（generic.just, go.just, node.just, python.just, rust.just, mixed.just），且模板含硬编码命令。SKILL.md 的 "generate from LLM knowledge" 指令与 .just 模板的存在直接矛盾。需决定模板的角色并统一。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-22`: P1 CONFLICT, SKILL.md HARD-RULE vs 模板存在 (source: Report 04)
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-23`: P1 CONFLICT, LLM-generation vs 硬编码模板 (source: Report 04)
- `plugins/forge/skills/init-justfile/SKILL.md`: 需修改的 SKILL.md Step 0 和 Step 3a (source: audit finding)
- `plugins/forge/skills/init-justfile/templates/*.just`: 6 个 .just 模板文件 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | 更新 Step 0 和 Step 3a 以明确模板角色 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] SKILL.md Step 0 的 HARD-RULE 与 templates/ 目录内容一致（不再矛盾）
- [ ] SKILL.md 明确说明模板的用途：是作为 LLM 生成的起点/参考，还是不应使用
- [ ] 若模板为活跃组件：Step 3a 应引用模板而非仅说 "generate from LLM knowledge"
- [ ] 若模板为遗留文件：删除 .just 模板或移至 `_deprecated/`

## Hard Rules
- 仅修改 `plugins/forge/skills/init-justfile/SKILL.md`，决策结果可能影响 templates/ 目录

## Implementation Notes
- 这是设计决策任务——需要先确定模板的角色（活跃 vs 遗留），然后更新 SKILL.md
- 推荐方案：将模板作为"起点模板"（starting point），SKILL.md Step 3a 先加载对应语言的 .just 模板，再根据 Convention 定制
- 此决策影响 Task 7（go.just 修复）的方向
