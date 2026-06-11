---
id: "1"
title: "新增 surface 策略模板和 test-type-model 参考"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: 新增 surface 策略模板和 test-type-model 参考

## Description
为 test-guide skill 创建 5 个 per-surface 策略模板（cli/api/web/tui/mobile）和 1 个完整的 test-type-model 参考文档。这些是 test-guide skill 重写的前置依赖——skill 将从这些模板生成用户项目的 convention 文件。

每个模板必须包含语言无关的 surface 测试策略，覆盖隔离模型、断言重点、超时策略等 7 个必要段落，以及断言偏好表（per-framework 一行）。test-type-model 需从 `docs/reference/test-type-model.md` 迁移完整内容。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution (核心变更 #1, #3), Non-Functional Requirements, Success Criteria
- `docs/reference/test-type-model.md`: 完整分类标准和语义定义需迁移至 plugin 层 (ref: Proposed Solution)
- `plugins/forge/skills/test-guide/templates/convention-template.md`: 现有模板结构作为新模板参考 (ref: Proposed Solution)
- `plugins/forge/skills/gen-test-scripts/types/cli.md`: CLI 测试策略权威来源，core.md 模板内容需与之一致 (ref: Out of Scope)
- `plugins/forge/skills/gen-test-scripts/types/_shared.md`: 跨 surface 共享策略参考 (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/test-guide/templates/surfaces/cli.md` | CLI 测试策略模板 |
| `plugins/forge/skills/test-guide/templates/surfaces/api.md` | API 测试策略模板 |
| `plugins/forge/skills/test-guide/templates/surfaces/web.md` | Web E2E 策略模板 |
| `plugins/forge/skills/test-guide/templates/surfaces/tui.md` | TUI 测试策略模板 |
| `plugins/forge/skills/test-guide/templates/surfaces/mobile.md` | Mobile E2E 策略模板 |
| `plugins/forge/skills/test-guide/references/test-type-model.md` | 完整 test-type-model 参考文档 |

## Acceptance Criteria
- [ ] `templates/surfaces/cli.md` 包含 7 个必要段落：文件位置、隔离模型、断言重点、超时策略、生命周期、Contract/Journey 比例、反模式，以及断言偏好表
- [ ] `templates/surfaces/api.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `templates/surfaces/web.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `templates/surfaces/tui.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `templates/surfaces/mobile.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `references/test-type-model.md` 包含完整分类标准、Surface → Test Type 映射表、e2e 术语约束和语义定义

## Hard Rules
- 每个模板只允许 7 个固定段落 + 断言偏好表。断言偏好表的列固定为：断言库、mock 机制、fixture 模式。新增列需通过提案评审
- 断言偏好表中的内容以 gen-test-scripts `types/*.md` 为首要权威来源

## Implementation Notes
- 先读取 `docs/reference/test-type-model.md` 获取完整的分类体系和语义定义，用于生成 `references/test-type-model.md`
- 读取 gen-test-scripts 的 `types/cli.md`、`types/api.md` 等文件提取各 surface 的策略信息，确保模板内容与 type rules 一致
- 模板是 Jinja 风格的占位符模板，由 test-guide skill 在运行时填充项目特有信息
- 注意内容膨胀防护：不要在模板中添加超出 7 段 + 偏好表的内容
