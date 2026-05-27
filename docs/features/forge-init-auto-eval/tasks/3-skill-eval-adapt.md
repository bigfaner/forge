---
id: "3"
title: "SKILL.md eval 检查统一适配"
priority: "P0"
estimated_time: "1h"
complexity: "low"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: SKILL.md eval 检查统一适配

## Description
将 4 个 SKILL.md 中的 eval 自动运行检查从 `$MODE` 依赖模式改为直接 bool 查询。当前模式先查 mode 再拼 `auto.eval.<type>.$MODE`，brainstorm 时 mode="none" 导致 key 不存在 → 永远 FALLBACK_ASK。

## Reference Files
- `plugins/forge/skills/brainstorm/SKILL.md`: line 122-144 eval 检查去掉 MODE 查询 (source: proposal.md#Part-3)
- `plugins/forge/skills/write-prd/SKILL.md`: line 226 eval 检查去掉 $MODE 后缀 (source: proposal.md#Part-3)
- `plugins/forge/skills/ui-design/SKILL.md`: line 147 eval 检查去掉 $MODE 后缀 (source: proposal.md#Part-3)
- `plugins/forge/skills/tech-design/SKILL.md`: line 181 eval 检查去掉 $MODE 后缀 (source: proposal.md#Part-3)

## Acceptance Criteria
- [ ] 4 个 SKILL.md 的 eval 检查不再查询 `$MODE`，直接 `forge config get auto.eval.<type>`
- [ ] 检查模式统一为：true→AUTO_RUN, false→SKIP, 其他→FALLBACK_ASK
- [ ] `grep -r 'auto\.eval\..*\$MODE' plugins/forge/skills/` 无残留匹配

## Implementation Notes
统一模式：
```bash
EVAL_ENABLED=$(forge config get auto.eval.<type> 2>/dev/null)
if [ "$EVAL_ENABLED" = "true" ]; then echo "AUTO_RUN"
elif [ "$EVAL_ENABLED" = "false" ]; then echo "SKIP"
else echo "FALLBACK_ASK"
fi
```
4 个文件改动完全对称，仅 `<type>` 不同（proposal/prd/uiDesign/techDesign）。
