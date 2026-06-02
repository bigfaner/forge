---
id: "8"
title: "Fix gen-test-scripts: stale references + rename ui.md → web.md"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 8: Fix gen-test-scripts: stale references + rename ui.md → web.md

## Description
修复 gen-test-scripts skill 中的三类问题：(1) 所有引用已删除的 `docs/reference/test-type-model.md` 的地方改为自包含内联描述（路径 A）；(2) 将 `types/ui.md` 重命名为 `types/web.md`，文件内容和所有引用处的 Surface 名称从 "UI" 改为 "Web"；(3) 确保所有 surfaceType 使用统一术语 web/api/cli/tui/mobile。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution, Success Criteria, Non-Functional Requirements
- `plugins/forge/skills/gen-test-scripts/SKILL.md:3,103,199,225`: 引用旧路径 + UI→Web 映射 (ref: Proposed Solution)
- `plugins/forge/skills/gen-test-scripts/types/ui.md`: 需重命名为 web.md，type 字段和内容中 "UI" → "Web" (ref: Proposed Solution)
- `plugins/forge/skills/gen-test-scripts/types/_shared.md:8,12`: interface type 列表中 "UI" → "Web" + 旧路径引用 (ref: Proposed Solution)
- `plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md:20,69`: 引用 `types/ui.md` → `types/web.md` (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/types/web.md` | 从 ui.md 重命名，内容中 "UI" → "Web" |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | 移除旧路径引用，内联映射表；`UI` → `types/web.md` |
| `plugins/forge/skills/gen-test-scripts/types/_shared.md` | "UI" → "Web" in type list；移除旧路径引用 |
| `plugins/forge/skills/gen-test-scripts/types/cli.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-test-scripts/types/api.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-test-scripts/types/tui.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-test-scripts/types/mobile.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md` | `types/ui.md` → `types/web.md` |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/gen-test-scripts/types/ui.md` | 已重命名为 web.md |

## Acceptance Criteria
- [ ] SKILL.md 不再引用 `docs/reference/test-type-model.md`，映射表已内联到文档中
- [ ] `types/ui.md` 已重命名为 `types/web.md`，文件内 `type: ui` 改为 `type: web`，所有 "UI" 引用改为 "Web"
- [ ] `types/_shared.md` 的 interface type 列表中 "UI" 改为 "Web"，旧路径引用已移除
- [ ] 所有 `types/*.md` 文件中 `docs/reference/test-type-model.md` 引用已替换为自包含的 test type 定义
- [ ] `rules/step-0.5-validation.md` 中 `types/ui.md` 引用已改为 `types/web.md`
- [ ] 所有 surfaceType 统一使用 web/api/cli/tui/mobile

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md` 了解 plugin 分发模型
- 不引用其它 skill 的内部文件，所有知识自包含
- 重命名 ui.md 时必须用 `git mv` 保留 git 历史

## Implementation Notes
- 内联映射表格式（供 SKILL.md 使用）：
  ```
  Surface → Test Type: cli → CLI Functional Test, api → API Functional Test, tui → Terminal Functional Test, web → Web E2E Test, mobile → Mobile E2E Test.
  e2e 术语仅用于 Web 和 Mobile surface。
  ```
- 每个 types/*.md 中的旧路径引用替换为一句话声明（如 "Test type: CLI 功能测试（CLI Functional Test），通过子进程执行验证进程退出码 + stdout/stderr。e2e 标签不适用于 CLI 测试。"）
