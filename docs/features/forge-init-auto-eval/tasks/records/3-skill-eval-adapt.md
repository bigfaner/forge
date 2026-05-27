---
status: "completed"
started: "2026-05-28 00:51"
completed: "2026-05-28 00:53"
time_spent: "~2m"
---

# Task Record: 3 SKILL.md eval 检查统一适配

## Summary
将 4 个 SKILL.md 的 eval auto-run 检查从 $MODE 依赖模式改为直接 bool 查询，消除 brainstorm mode=none 时 key 不存在的 bug

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/brainstorm/SKILL.md
- plugins/forge/skills/write-prd/SKILL.md
- plugins/forge/skills/ui-design/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md

### Key Decisions
- 统一 4 个文件的 eval 检查模式：去掉 MODE 中间查询，直接 forge config get auto.eval.<type>

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1
- **Failed**: 0
- **Coverage**: 100.0%

## Acceptance Criteria
- [x] 4 个 SKILL.md 的 eval 检查不再查询 $MODE，直接 forge config get auto.eval.<type>
- [x] 检查模式统一为：true→AUTO_RUN, false→SKIP, 其他→FALLBACK_ASK
- [x] grep -r 'auto\.eval\..*\$MODE' plugins/forge/skills/ 无残留匹配

## Notes
SKILL.md 为模板文件，验证通过 grep 确认无 $MODE 残留 + 4 个文件 bash 逻辑正确对称。compile/fmt/lint 全部通过。
