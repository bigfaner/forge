---
id: "10"
title: "Fix guide.md stale reference + unify test directory strategy"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 10: Fix guide.md stale reference + unify test directory strategy

## Description
修复 guide.md hook 中指向已删除的 `docs/reference/test-type-model.md` 的链接，改为自包含描述。同时统一测试目录策略：当前 test-guide 模板和 guide.md 告知 web/mobile 测试放 `tests/e2e/`，但 gen-test-scripts 统一生成到 `tests/<journey>/` 并明确禁止 `tests/e2e/`。需要将 test-guide 模板和 guide.md 中的目录说明统一为 `tests/<journey>/`。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution, Non-Functional Requirements
- `plugins/forge/hooks/guide.md:97`: 指向已删除文件的链接 (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/templates/surfaces/web.md:12`: `tests/e2e/` 目录 (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/templates/surfaces/mobile.md:12`: `tests/e2e/` 目录 (ref: Proposed Solution)
- `plugins/forge/skills/gen-test-scripts/SKILL.md:231,251`: 明确禁止 `tests/e2e/` (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | 移除旧路径链接，改为自包含描述；web/mobile 测试目录从 `tests/e2e/` 改为 `tests/<journey>/` |
| `plugins/forge/skills/test-guide/templates/surfaces/web.md` | 文件位置规则从 `tests/e2e/` 改为 `tests/<journey>/` |
| `plugins/forge/skills/test-guide/templates/surfaces/mobile.md` | 文件位置规则从 `tests/e2e/` 改为 `tests/<journey>/` |

## Acceptance Criteria
- [ ] guide.md 不再包含 `docs/reference/test-type-model.md` 链接，映射表和 e2e 约束已自包含
- [ ] guide.md 中 web/mobile 测试目录说明从 `tests/e2e/` 改为 `tests/<journey>/`，与 gen-test-scripts 一致
- [ ] `test-guide/templates/surfaces/web.md` 的文件位置段落不再引用 `tests/e2e/`
- [ ] `test-guide/templates/surfaces/mobile.md` 的文件位置段落不再引用 `tests/e2e/`
- [ ] guide.md Testing section 总行数仍 <= 20 行

## Hard Rules
- guide.md 增量 <= 20 行
- 不引用其它 skill 的内部文件

## Implementation Notes
- 测试目录统一规则：所有 surface 的测试代码都生成到 `tests/<journey>/`，不区分 surface type。journey 名称由 gen-journeys 生成，不是 surface type 名称
- guide.md 中的目录说明应简化为通用规则，如 "测试代码生成到 tests/<journey>/ 目录"
