---
id: "5"
title: "Add INLINE cross-skill version markers (M-9)"
priority: "P1"
estimated_time: "30m"
dependencies: [2, 3, 4]
type: "doc"
mainSession: false
---

# 5: Add INLINE cross-skill version markers (M-9)

## Description

4 处 INLINE 跨 skill 引用缺乏同步标记。如果源文件更新但内联副本未同步，LLM 会基于过时的模型定义生成 journeys/contracts，产生不会报错的语义间隙。需要为每处 INLINE 引用添加源文件版本号标记。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — M-9: INLINE 跨 skill 引用同步风险（4 处，由 L-4 升级）, Proposed Solution
- `plugins/forge/skills/gen-journeys/SKILL.md` or rules: Add INLINE marker for journey-contract-model.md (ref: M-9: INLINE 跨 skill 引用同步风险（4 处，由 L-4 升级）)
- `plugins/forge/skills/gen-contracts/SKILL.md` or rules: Add INLINE marker for Surface Detection logic (ref: M-9: INLINE 跨 skill 引用同步风险（4 处，由 L-4 升级）)
- `plugins/forge/skills/gen-test-scripts/SKILL.md` or rules: Add INLINE marker for test-isolation.md (ref: M-9: INLINE 跨 skill 引用同步风险（4 处，由 L-4 升级）)
- `plugins/forge/skills/init-justfile/SKILL.md` or rules: Add INLINE marker for test-type-model.md (ref: M-9: INLINE 跨 skill 引用同步风险（4 处，由 L-4 升级）)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/` (file TBD by grep) | Add `<!-- INLINE from ... @ v3.0.0-rc.53 -->` marker |
| `plugins/forge/skills/gen-contracts/` (file TBD by grep) | Add INLINE version marker |
| `plugins/forge/skills/gen-test-scripts/` (file TBD by grep) | Add INLINE version marker |
| `plugins/forge/skills/init-justfile/` (file TBD by grep) | Add INLINE version marker |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 4 处 INLINE 引用均标注源文件路径和版本号标记（格式：`<!-- INLINE from <source-path> @ <version> -->`）
- [ ] `grep -r "INLINE" plugins/forge/skills/` 显示所有内联引用均有版本号标记

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes
- 先用 `grep -rn "INLINE"` 定位 4 处引用的精确文件和行号，再逐一添加标记

## Implementation Notes
- 4 处 INLINE 引用：(1) gen-journeys 内联 gen-contracts 的 journey-contract-model.md；(2) gen-contracts 内联 gen-journeys 的 Surface Detection 逻辑（双向依赖）；(3) gen-test-scripts 内联 run-tests 的 test-isolation.md；(4) init-justfile 内联 test-guide 的 test-type-model.md
