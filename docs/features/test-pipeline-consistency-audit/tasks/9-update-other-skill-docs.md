---
id: "9"
title: "更新其他 Skill/Command 文档中的旧术语和路径"
priority: "P1"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 9: 更新其他 Skill/Command 文档中的旧术语和路径

## Description
批量更新多个 Skill/Command 文档中的旧路径引用和术语：`fix-bug.md` 路径更新；`run-tasks.md` 中 `T-test-verify-regression` 和 "e2e verification" 引用清理；test-guide 术语修正（含 build tag 表格）；`submit-task` record template 中 `test.verify-regression` 类型删除和路径更新；`gen-sitemap` 配置文件重命名（`e2e-config.yaml` → `test-config.yaml`）及 SKILL.md 引用更新；`consolidate-specs/SKILL.md` 术语修正；`init-justfile` 路径更新和 6 个 justfile 模板中 `tests/e2e/` 路径更新；`run-tests/rules/test-isolation.md` 路径引用更新。

## Reference Files
- `proposal.md#Layer-2-Skill-文档层术语统一` — 第 11 项定义了所有需修改的文档列表
- `proposal.md#Success-Criteria` — 验证条件：build tag 对齐、术语替换完成

## Acceptance Criteria
- [ ] `commands/fix-bug.md` 中 `tests/e2e/features/` → `tests/<journey>/`
- [ ] `commands/run-tasks.md` 中 `T-test-verify-regression` 和 "e2e verification" 引用已清理
- [ ] test-guide（含 `rules/draft-generation.md` 和 `rules/pattern-extraction.md`）术语修正
- [ ] `submit-task/data/record-format-test.md` 删除 `test.verify-regression` 类型，更新 `tests/e2e/` 示例路径
- [ ] `gen-sitemap` 配置文件从 `e2e-config.yaml` 重命名为 `test-config.yaml`，SKILL.md 路径引用更新
- [ ] `consolidate-specs/SKILL.md` 第 22 行 "e2e tests are promoted" → "all tests pass"
- [ ] `init-justfile/SKILL.md` 中 `tests/e2e/` 示例路径更新
- [ ] `init-justfile/templates/` 下 6 个 justfile 模板中 `tests/e2e/` 路径更新
- [ ] `run-tests/rules/test-isolation.md` 中 4 处 `tests/e2e/` 路径引用更新
- [ ] Convention 文件中 build tag 与 surface 类型对齐

## Implementation Notes
- 范围广但每处修改简单（术语/路径替换）
- build tag 表格需与 task 5（Go build tag 重命名）对齐

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | 路径更新 |
| `plugins/forge/commands/run-tasks.md` | 旧类型和术语清理 |
| `plugins/forge/skills/run-tests/rules/test-isolation.md` | 4 处路径引用更新 |
| `plugins/forge/skills/submit-task/data/record-format-test.md` | 类型删除 + 路径更新 |
| `plugins/forge/skills/gen-sitemap/SKILL.md` | 配置文件路径更新 |
| `plugins/forge/skills/consolidate-specs/SKILL.md` | 术语修正 |
| `plugins/forge/skills/init-justfile/SKILL.md` | 示例路径更新 |
| `plugins/forge/skills/init-justfile/templates/python.just` | `tests/e2e/` 路径更新 |
| `plugins/forge/skills/init-justfile/templates/rust.just` | `tests/e2e/` 路径更新 |
| `plugins/forge/skills/init-justfile/templates/node.just` | `tests/e2e/` 路径更新 |
| `plugins/forge/skills/init-justfile/templates/go.just` | `tests/e2e/` 路径更新 |
| `plugins/forge/skills/init-justfile/templates/mixed.just` | `tests/e2e/` 路径更新 |
| `plugins/forge/skills/init-justfile/templates/generic.just` | `tests/e2e/` 路径更新 |

### Delete
| File | Reason |
|------|--------|
